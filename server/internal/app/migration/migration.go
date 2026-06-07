package migration

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	DefaultDir = "migrations/mysql"
	timeLayout = "20060102150405"
)

var (
	filePattern = regexp.MustCompile(`^(\d{14})_([a-z][a-z0-9_]*?)\.(up|down)\.sql$`)
	namePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	nameCleaner = regexp.MustCompile(`[^a-z0-9]+`)
)

type CreateOptions struct {
	RootDir string
	Dir     string
	Name    string
	Now     time.Time
}

type Pair struct {
	Version string
	Name    string
	Up      string
	Down    string
}

type File struct {
	Version   string
	Name      string
	Direction string
	Path      string
}

type Group struct {
	Version string
	Name    string
	Up      string
	Down    string
}

func Create(opts CreateOptions) (Pair, error) {
	opts = normalizeCreateOptions(opts)
	if opts.Name == "" {
		return Pair{}, errors.New("migration name is required")
	}
	name := NormalizeName(opts.Name)
	if !namePattern.MatchString(name) {
		return Pair{}, fmt.Errorf("invalid migration name %q", opts.Name)
	}
	if opts.Now.IsZero() {
		opts.Now = time.Now()
	}

	dir := filepath.Join(opts.RootDir, opts.Dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Pair{}, err
	}

	version := opts.Now.Format(timeLayout)
	pair := Pair{
		Version: version,
		Name:    name,
		Up:      filepath.Join(dir, version+"_"+name+".up.sql"),
		Down:    filepath.Join(dir, version+"_"+name+".down.sql"),
	}
	if err := writeNewFile(pair.Up, upTemplate(name)); err != nil {
		return Pair{}, err
	}
	if err := writeNewFile(pair.Down, downTemplate(name)); err != nil {
		_ = os.Remove(pair.Up)
		return Pair{}, err
	}
	return pair, nil
}

func NormalizeName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = nameCleaner.ReplaceAllString(name, "_")
	name = strings.Trim(name, "_")
	for strings.Contains(name, "__") {
		name = strings.ReplaceAll(name, "__", "_")
	}
	return name
}

func List(rootDir, dir string) ([]Group, error) {
	files, err := parseFiles(filepath.Join(defaultRoot(rootDir), defaultDir(dir)))
	if err != nil {
		return nil, err
	}
	groups := groupFiles(files)
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Version < groups[j].Version
	})
	return groups, nil
}

func Validate(rootDir, dir string) error {
	groups, err := List(rootDir, dir)
	if err != nil {
		return err
	}
	seenVersions := map[string]string{}
	for _, group := range groups {
		if existingName, ok := seenVersions[group.Version]; ok && existingName != group.Name {
			return fmt.Errorf("version %s is used by both %q and %q", group.Version, existingName, group.Name)
		}
		seenVersions[group.Version] = group.Name
		if group.Up == "" {
			return fmt.Errorf("migration %s_%s missing up file", group.Version, group.Name)
		}
		if group.Down == "" {
			return fmt.Errorf("migration %s_%s missing down file", group.Version, group.Name)
		}
	}
	return nil
}

func parseFiles(dir string) ([]File, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	files := make([]File, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == ".gitkeep" || entry.Name() == "README.md" {
			continue
		}
		matches := filePattern.FindStringSubmatch(entry.Name())
		if matches == nil {
			return nil, fmt.Errorf("invalid migration filename %q", entry.Name())
		}
		files = append(files, File{
			Version:   matches[1],
			Name:      matches[2],
			Direction: matches[3],
			Path:      filepath.Join(dir, entry.Name()),
		})
	}
	sort.Slice(files, func(i, j int) bool {
		if files[i].Version == files[j].Version {
			return files[i].Direction < files[j].Direction
		}
		return files[i].Version < files[j].Version
	})
	return files, nil
}

func groupFiles(files []File) []Group {
	index := map[string]int{}
	groups := make([]Group, 0, len(files)/2)
	for _, file := range files {
		key := file.Version + "_" + file.Name
		pos, ok := index[key]
		if !ok {
			groups = append(groups, Group{Version: file.Version, Name: file.Name})
			pos = len(groups) - 1
			index[key] = pos
		}
		switch file.Direction {
		case "up":
			groups[pos].Up = file.Path
		case "down":
			groups[pos].Down = file.Path
		}
	}
	return groups
}

func normalizeCreateOptions(opts CreateOptions) CreateOptions {
	opts.RootDir = defaultRoot(opts.RootDir)
	opts.Dir = defaultDir(opts.Dir)
	return opts
}

func defaultRoot(root string) string {
	if root == "" {
		return "."
	}
	return root
}

func defaultDir(dir string) string {
	if dir == "" {
		return DefaultDir
	}
	return dir
}

func writeNewFile(path string, content string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return err
}

func upTemplate(name string) string {
	return fmt.Sprintf(`-- Migration: %s
-- Direction: up

-- Write forward migration SQL here.
`, name)
}

func downTemplate(name string) string {
	return fmt.Sprintf(`-- Migration: %s
-- Direction: down

-- Write rollback migration SQL here.
`, name)
}
