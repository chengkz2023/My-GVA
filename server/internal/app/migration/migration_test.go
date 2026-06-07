package migration

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	root := t.TempDir()
	pair, err := Create(CreateOptions{
		RootDir: root,
		Name:    "Create Users",
		Now:     time.Date(2026, 6, 7, 12, 30, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if pair.Version != "20260607123000" || pair.Name != "create_users" {
		t.Fatalf("pair = %+v, want normalized timestamp/name", pair)
	}
	assertExists(t, filepath.Join(root, DefaultDir, "20260607123000_create_users.up.sql"))
	assertExists(t, filepath.Join(root, DefaultDir, "20260607123000_create_users.down.sql"))
}

func TestValidate(t *testing.T) {
	root := t.TempDir()
	if _, err := Create(CreateOptions{
		RootDir: root,
		Name:    "create users",
		Now:     time.Date(2026, 6, 7, 12, 30, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := Validate(root, ""); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestValidateMissingPair(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, DefaultDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "20260607123000_create_users.up.sql"), []byte("select 1;"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := Validate(root, "")
	if err == nil {
		t.Fatal("Validate() error is nil, want missing down file")
	}
}

func TestValidateRejectsBadFilename(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, DefaultDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "create_users.sql"), []byte("select 1;"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := Validate(root, "")
	if err == nil {
		t.Fatal("Validate() error is nil, want invalid filename")
	}
}

func TestListSortsGroups(t *testing.T) {
	root := t.TempDir()
	for _, item := range []struct {
		name string
		now  time.Time
	}{
		{name: "second", now: time.Date(2026, 6, 7, 12, 31, 0, 0, time.UTC)},
		{name: "first", now: time.Date(2026, 6, 7, 12, 30, 0, 0, time.UTC)},
	} {
		if _, err := Create(CreateOptions{RootDir: root, Name: item.name, Now: item.now}); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	groups, err := List(root, "")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(groups) != 2 || groups[0].Name != "first" || groups[1].Name != "second" {
		t.Fatalf("groups = %+v, want sorted first/second", groups)
	}
}

func assertExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
}
