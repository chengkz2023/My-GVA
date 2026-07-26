package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/migration"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "create":
		err = create(os.Args[2:])
	case "list":
		err = list(os.Args[2:])
	case "validate":
		err = validate(os.Args[2:])
	case "up":
		err = up(os.Args[2:])
	case "down":
		err = down(os.Args[2:])
	case "status":
		err = status(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runOpts(args []string) (migration.RunOptions, error) {
	fs := flag.NewFlagSet("opts", flag.ExitOnError)
	opts := migration.RunOptions{
		RootDir: ".",
		Dir:     migration.DefaultDir,
	}
	fs.StringVar(&opts.RootDir, "dir", ".", "server root directory")
	fs.StringVar(&opts.Dir, "migrations", migration.DefaultDir, "migration directory")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "print what would be executed")
	fs.BoolVar(&opts.DryRun, "n", false, "alias for --dry-run")
	fs.StringVar(&opts.DSN, "dsn", "", "database DSN override")
	if err := fs.Parse(args); err != nil {
		return opts, err
	}
	return opts, nil
}

func create(args []string) error {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	rootDir := fs.String("dir", ".", "server root directory")
	migrationDir := fs.String("migrations", migration.DefaultDir, "migration directory")
	name := fs.String("name", "", "migration name, e.g. create_users")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *name == "" && fs.NArg() > 0 {
		*name = fs.Arg(0)
	}

	pair, err := migration.Create(migration.CreateOptions{
		RootDir: *rootDir,
		Dir:     *migrationDir,
		Name:    *name,
	})
	if err != nil {
		return err
	}

	fmt.Printf("created migration %s_%s\n", pair.Version, pair.Name)
	fmt.Printf("  up:   %s\n", pair.Up)
	fmt.Printf("  down: %s\n", pair.Down)
	return nil
}

func list(args []string) error {
	opts, err := runOpts(args)
	if err != nil {
		return err
	}
	groups, err := migration.List(opts.RootDir, opts.Dir)
	if err != nil {
		return err
	}
	if len(groups) == 0 {
		fmt.Println("no migrations")
		return nil
	}
	for _, group := range groups {
		fmt.Printf("%s_%s\n", group.Version, group.Name)
	}
	return nil
}

func validate(args []string) error {
	opts, err := runOpts(args)
	if err != nil {
		return err
	}
	if err := migration.Validate(opts.RootDir, opts.Dir); err != nil {
		return err
	}
	fmt.Println("migrations ok")
	return nil
}

func up(args []string) error {
	opts, err := runOpts(args)
	if err != nil {
		return err
	}
	db, err := migration.InitDB()
	if err != nil {
		return err
	}
	db, err = migration.MustDB(db)
	if err != nil {
		return err
	}
	return migration.RunUp(db, opts)
}

func down(args []string) error {
	opts, err := runOpts(args)
	if err != nil {
		return err
	}
	db, err := migration.InitDB()
	if err != nil {
		return err
	}
	db, err = migration.MustDB(db)
	if err != nil {
		return err
	}
	return migration.RunDown(db, opts)
}

func status(args []string) error {
	opts, err := runOpts(args)
	if err != nil {
		return err
	}
	db, err := migration.InitDB()
	if err != nil {
		return err
	}
	db, err = migration.MustDB(db)
	if err != nil {
		return err
	}
	return migration.RunStatus(db, opts)
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage:
  go run ./cmd/migrate create -name create_users
  go run ./cmd/migrate list
  go run ./cmd/migrate validate
  go run ./cmd/migrate up             [--dry-run] [--dir .] [--migrations migrations/mysql]
  go run ./cmd/migrate down           [--dry-run]
  go run ./cmd/migrate status

options:
  --dir DIR         server root directory (default: .)
  --migrations DIR  migration directory relative to root (default: migrations/mysql)
  --dry-run, -n     print what would be executed without running
  --dsn DSN         override database connection string`)
}
