package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	platformconfig "github.com/chengkz2023/My-GVA/server/internal/platform/config"
	"github.com/chengkz2023/My-GVA/server/internal/platform/database"
	"github.com/chengkz2023/My-GVA/server/internal/platform/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MigrationRecord struct {
	Version   string `gorm:"primaryKey;column:version"`
	AppliedAt string `gorm:"column:applied_at"`
}

func (MigrationRecord) TableName() string {
	return "schema_migrations"
}

type RunOptions struct {
	RootDir string
	Dir     string
	DryRun  bool
	DSN     string
}

func InitDB() (*gorm.DB, error) {
	_, cfg := platformconfig.Load()
	log := logger.New(cfg.Zap)
	zap.ReplaceGlobals(log)
	db := database.Open(cfg.Mysql, log)
	return db, nil
}

func InitDBWithDSN(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return InitDB()
	}
	db, err := gorm.Open(nil, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	return db, nil
}

func MustDB(db *gorm.DB) (*gorm.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}
	if err := db.AutoMigrate(&MigrationRecord{}); err != nil {
		return nil, fmt.Errorf("init migration table: %w", err)
	}
	return db, nil
}

func RunUp(db *gorm.DB, opts RunOptions) error {
	groups, err := List(opts.RootDir, opts.Dir)
	if err != nil {
		return err
	}
	if len(groups) == 0 {
		fmt.Println("no migrations found")
		return nil
	}

	var applied []string
	db.Model(&MigrationRecord{}).Pluck("version", &applied)
	appliedSet := make(map[string]bool, len(applied))
	for _, v := range applied {
		appliedSet[v] = true
	}

	pending := make([]Group, 0)
	for _, g := range groups {
		if !appliedSet[g.Version] {
			pending = append(pending, g)
		}
	}
	if len(pending) == 0 {
		fmt.Println("all migrations are up to date")
		return nil
	}

	sort.Slice(pending, func(i, j int) bool { return pending[i].Version < pending[j].Version })

	if opts.DryRun {
		fmt.Println("[dry-run] would apply:")
		for _, g := range pending {
			fmt.Printf("  %s_%s\n", g.Version, g.Name)
		}
		return nil
	}

	for _, group := range pending {
		content, err := os.ReadFile(group.Up)
		if err != nil {
			return fmt.Errorf("read %s: %w", group.Up, err)
		}
		statements := parseStatements(string(content))
		fmt.Printf("applying %s_%s... ", group.Version, group.Name)

		err = db.Transaction(func(tx *gorm.DB) error {
			if execErr := execStatements(tx, statements); execErr != nil {
				return execErr
			}
			return tx.Create(&MigrationRecord{Version: group.Version}).Error
		})
		if err != nil {
			return fmt.Errorf("run %s: %w", filepath.Base(group.Up), err)
		}
		fmt.Println("ok")
	}
	fmt.Printf("applied %d migration(s)\n", len(pending))
	return nil
}

func RunDown(db *gorm.DB, opts RunOptions) error {
	var last MigrationRecord
	if err := db.Order("version desc").First(&last).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("no migrations to roll back")
			return nil
		}
		return err
	}

	groups, err := List(opts.RootDir, opts.Dir)
	if err != nil {
		return err
	}
	groupMap := make(map[string]Group, len(groups))
	for _, g := range groups {
		groupMap[g.Version] = g
	}

	group, ok := groupMap[last.Version]
	if !ok {
		return fmt.Errorf("migration %s not found on disk", last.Version)
	}

	if opts.DryRun {
		fmt.Printf("[dry-run] would rollback %s_%s\n", group.Version, group.Name)
		return nil
	}

	content, err := os.ReadFile(group.Down)
	if err != nil {
		return fmt.Errorf("read %s: %w", group.Down, err)
	}
	statements := parseStatements(string(content))
	fmt.Printf("rolling back %s_%s... ", group.Version, group.Name)

	err = db.Transaction(func(tx *gorm.DB) error {
		if execErr := execStatements(tx, statements); execErr != nil {
			return execErr
		}
		return tx.Where("version = ?", last.Version).Delete(&MigrationRecord{}).Error
	})
	if err != nil {
		return fmt.Errorf("run %s: %w", filepath.Base(group.Down), err)
	}
	fmt.Println("ok")
	return nil
}

func RunStatus(db *gorm.DB, opts RunOptions) error {
	groups, err := List(opts.RootDir, opts.Dir)
	if err != nil {
		return err
	}
	var applied []string
	db.Model(&MigrationRecord{}).Pluck("version", &applied)
	appliedSet := make(map[string]bool, len(applied))
	for _, v := range applied {
		appliedSet[v] = true
	}

	for _, g := range groups {
		status := "[ ]"
		if appliedSet[g.Version] {
			status = "[x]"
		}
		fmt.Printf("%s %s_%s\n", status, g.Version, g.Name)
	}
	return nil
}

func parseStatements(content string) []string {
	lines := strings.Split(content, "\n")
	var statements []string
	var current strings.Builder
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}
		current.WriteString(line)
		current.WriteString("\n")
		if strings.HasSuffix(trimmed, ";") {
			statements = append(statements, current.String())
			current.Reset()
		}
	}
	remaining := strings.TrimSpace(current.String())
	if remaining != "" {
		statements = append(statements, remaining)
	}
	return statements
}

func execStatements(db *gorm.DB, statements []string) error {
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if err := db.Exec(stmt).Error; err != nil {
			return fmt.Errorf("exec %q: %w", truncateSQL(stmt, 80), err)
		}
	}
	return nil
}

func truncateSQL(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
