package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateLightModule(t *testing.T) {
	root := t.TempDir()

	result, err := GenerateLightModule(ModuleOptions{
		RootDir: root,
		Group:   "business",
		Name:    "customer",
	})
	if err != nil {
		t.Fatalf("GenerateLightModule() error = %v", err)
	}

	wantDir := filepath.Join(root, defaultModuleRoot, "business", "customer")
	if result.Dir != wantDir {
		t.Fatalf("dir = %s, want %s", result.Dir, wantDir)
	}
	if len(result.Files) != len(lightModuleFiles) {
		t.Fatalf("files = %d, want %d", len(result.Files), len(lightModuleFiles))
	}

	assertFileContains(t, filepath.Join(wantDir, "module.go"), "func (m *Module) RegisterHTTP(routes apphttp.Routes)")
	assertFileContains(t, filepath.Join(wantDir, "handler.go"), `moduleGroup := group.Group("/customer")`)
	assertFileContains(t, filepath.Join(wantDir, "handler_test.go"), `"/api/customer/info"`)
	assertFileContains(t, filepath.Join(wantDir, "repository.go"), `Name:    "business/customer"`)
}

func TestGenerateLightModuleRejectsInvalidName(t *testing.T) {
	_, err := GenerateLightModule(ModuleOptions{
		RootDir: t.TempDir(),
		Group:   "business",
		Name:    "customer-order",
	})
	if err == nil {
		t.Fatal("GenerateLightModule() error is nil, want invalid name")
	}
}

func TestGenerateLightModuleRejectsInvalidRoute(t *testing.T) {
	_, err := GenerateLightModule(ModuleOptions{
		RootDir: t.TempDir(),
		Group:   "business",
		Name:    "customer",
		Route:   `/customer"`,
	})
	if err == nil {
		t.Fatal("GenerateLightModule() error is nil, want invalid route")
	}
}

func TestGenerateLightModuleRefusesOverwrite(t *testing.T) {
	root := t.TempDir()
	opts := ModuleOptions{
		RootDir: root,
		Group:   "business",
		Name:    "customer",
	}
	if _, err := GenerateLightModule(opts); err != nil {
		t.Fatalf("GenerateLightModule() error = %v", err)
	}
	if _, err := GenerateLightModule(opts); err == nil {
		t.Fatal("GenerateLightModule() error is nil, want overwrite refusal")
	}
}

func assertFileContains(t *testing.T, path string, want string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !strings.Contains(string(content), want) {
		t.Fatalf("%s does not contain %q\n%s", path, want, content)
	}
}
