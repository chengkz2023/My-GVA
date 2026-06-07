package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/scaffold"
)

func main() {
	var opts scaffold.ModuleOptions
	flag.StringVar(&opts.RootDir, "dir", ".", "server root directory")
	flag.StringVar(&opts.Group, "group", "business", "module group under internal/modules")
	flag.StringVar(&opts.Name, "name", "", "module name, lowercase letters and digits")
	flag.StringVar(&opts.Route, "route", "", "HTTP route under /v2, default is /<name>")
	flag.BoolVar(&opts.Force, "force", false, "overwrite existing files")
	flag.Parse()

	result, err := scaffold.GenerateLightModule(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("created module: %s\n", filepath.ToSlash(result.Dir))
	for _, file := range result.Files {
		fmt.Printf("  %s\n", filepath.ToSlash(file))
	}
	fmt.Println()
	fmt.Println("Next: register the module in internal/modules/modules.go after reviewing the generated files.")
}
