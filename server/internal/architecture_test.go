package internal

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestArchitectureBoundaries(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		imports, err := fileImports(path)
		if err != nil {
			return err
		}

		assertPlatformDoesNotImportModules(t, rel, imports)
		assertPlatformDoesNotImportLegacyRootPackages(t, rel, imports)
		assertDomainIsFrameworkFree(t, rel, imports)
		assertHandlersDoNotUsePersistence(t, rel, imports)
		assertServicesStayTransportFree(t, rel, imports)
		assertModulesDoNotImportOtherInfrastructure(t, rel, imports)
				return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func fileImports(path string) ([]string, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}
	imports := make([]string, 0, len(file.Imports))
	for _, spec := range file.Imports {
		value, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return nil, err
		}
		imports = append(imports, value)
	}
	return imports, nil
}

func assertPlatformDoesNotImportModules(t *testing.T, rel string, imports []string) {
	t.Helper()
	if !strings.HasPrefix(rel, "platform/") {
		return
	}
	for _, imp := range imports {
		if strings.Contains(imp, "/internal/modules") {
			t.Fatalf("%s: platform package must not import module package %q", rel, imp)
		}
	}
}

// assertPlatformDoesNotImportLegacyRootPackages 防止 platform 依赖顶层遗留包
// （如已被移除的 utils/、task/）。platform 只允许依赖自身内部包与 server/config。
func assertPlatformDoesNotImportLegacyRootPackages(t *testing.T, rel string, imports []string) {
	t.Helper()
	if !strings.HasPrefix(rel, "platform/") {
		return
	}
	for _, imp := range imports {
		if !strings.Contains(imp, "/server/") {
			continue
		}
		if strings.Contains(imp, "/server/config") || strings.Contains(imp, "/server/internal/") {
			continue
		}
		t.Fatalf("%s: platform package must not import legacy top-level package %q", rel, imp)
	}
}

func assertDomainIsFrameworkFree(t *testing.T, rel string, imports []string) {
	t.Helper()
	if !strings.Contains(rel, "/domain/") {
		return
	}
	banned := []string{
		"github.com/gin-gonic/gin",
		"go.uber.org/zap",
		"github.com/spf13/viper",
		"github.com/redis/go-redis",
		"github.com/casbin/casbin",
		"gorm.io/gorm",
		"/transport/",
		"/infrastructure/",
	}
	assertNoImports(t, rel, imports, banned, "domain package must stay framework-free")
}

func assertHandlersDoNotUsePersistence(t *testing.T, rel string, imports []string) {
	t.Helper()
	if !strings.HasSuffix(rel, "handler.go") && !strings.Contains(rel, "/transport/http/") {
		return
	}
	banned := []string{
		"gorm.io/gorm",
		"/infrastructure/",
	}
	assertNoImports(t, rel, imports, banned, "HTTP handler must not access persistence directly")
}

func assertServicesStayTransportFree(t *testing.T, rel string, imports []string) {
	t.Helper()
	if !strings.HasSuffix(rel, "service.go") && !strings.Contains(rel, "/application/") {
		return
	}
	banned := []string{
		"github.com/gin-gonic/gin",
	}
	assertNoImports(t, rel, imports, banned, "service/application package must stay transport-free and avoid global state")
}

func assertModulesDoNotImportOtherInfrastructure(t *testing.T, rel string, imports []string) {
	t.Helper()
	if !strings.HasPrefix(rel, "modules/") {
		return
	}
	moduleRoot := moduleRoot(rel)
	if moduleRoot == "" {
		return
	}
	for _, imp := range imports {
		if !strings.Contains(imp, "/internal/modules/") || !strings.Contains(imp, "/infrastructure/") {
			continue
		}
		if !strings.Contains(imp, "/internal/"+moduleRoot+"/") {
			t.Fatalf("%s: module must not import another module infrastructure package %q", rel, imp)
		}
	}
}

func moduleRoot(rel string) string {
	parts := strings.Split(rel, "/")
	if len(parts) < 3 || parts[0] != "modules" {
		return ""
	}
	if parts[1] == "business" && len(parts) >= 4 {
		return strings.Join(parts[:3], "/")
	}
	if parts[1] == "system" && len(parts) >= 4 {
		return strings.Join(parts[:3], "/")
	}
	return strings.Join(parts[:2], "/")
}

func assertNoImports(t *testing.T, rel string, imports []string, banned []string, message string) {
	t.Helper()
	for _, imp := range imports {
		for _, ban := range banned {
			if strings.Contains(imp, ban) {
				t.Fatalf("%s: %s: %q", rel, message, imp)
			}
		}
	}
}
