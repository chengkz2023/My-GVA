package database

import (
	"context"
	"strings"

	"github.com/chengkz2023/My-GVA/server/config"
	"github.com/chengkz2023/My-GVA/server/internal/platform/database/internal"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Open(cfg config.Mysql, log *zap.Logger) (*gorm.DB, error) {
	return openMySQL(cfg, log)
}

func Ping(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return gorm.ErrInvalidDB
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func openMySQL(m config.Mysql, log *zap.Logger) (*gorm.DB, error) {
	if m.Dbname == "" {
		return nil, nil
	}
	dsn := m.Dsn()
	if m.Config == "" || !strings.Contains(m.Config, "timeout=") {
		dsn += "&timeout=5s&readTimeout=10s&writeTimeout=10s"
	}
	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         191,
		SkipInitializeWithVersion: false,
	}
	general := m.GeneralDB
	db, err := gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config(general, log))
	if err != nil {
		return nil, err
	}
	if m.Engine != "" {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	return db, nil
}
