package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultsItemNameFromRepoRoot(t *testing.T) {
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	project := filepath.Join(repo, "services", "api")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(project, LoadOptions{})
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.ItemName != filepath.Base(repo) {
		t.Fatalf("expected item name %q, got %q", filepath.Base(repo), cfg.ItemName)
	}
	if cfg.SchemaFile != filepath.Join(project, ".env.example") {
		t.Fatalf("unexpected schema file: %s", cfg.SchemaFile)
	}
}

func TestLoadDefaultsItemNameWithoutGit(t *testing.T) {
	project := t.TempDir()

	cfg, err := Load(project, LoadOptions{})
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.ItemName != filepath.Base(project) {
		t.Fatalf("expected fallback item name %q, got %q", filepath.Base(project), cfg.ItemName)
	}
}

func TestLoadItemNameOverrideFromConfig(t *testing.T) {
	project := t.TempDir()
	if err := os.WriteFile(filepath.Join(project, ".envsync.yaml"), []byte("item_name: shared-dev\nmapping:\n  DATABASE_URL: db_url\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(project, LoadOptions{})
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.ItemName != "shared-dev" {
		t.Fatalf("expected override item name, got %q", cfg.ItemName)
	}
	if got := cfg.Mapping["DATABASE_URL"]; got != "db_url" {
		t.Fatalf("expected field override db_url, got %q", got)
	}
}
