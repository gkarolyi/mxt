package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/gkarolyi/mxt/internal/config"
)

func TestMigrateConfigFileLegacy(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")
	legacyContent := "worktree_dir=~/worktrees\nterminal=iterm2\ncopy_files=.env,.env.local\n"
	if err := os.WriteFile(configPath, []byte(legacyContent), 0o644); err != nil {
		t.Fatalf("failed to write legacy config: %v", err)
	}

	migrated, err := migrateConfigFile(configPath)
	if err != nil {
		t.Fatalf("migrateConfigFile() error = %v", err)
	}
	if !migrated {
		t.Fatal("expected migrateConfigFile() to report migration")
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read migrated config: %v", err)
	}

	parsed, err := config.ParseConfig(bytes.NewReader(content))
	if err != nil {
		t.Fatalf("ParseConfig() error = %v", err)
	}

	expected := map[string]string{
		"worktree_dir": "~/worktrees",
		"terminal":     "iterm2",
		"copy_files":   ".env,.env.local",
	}
	for key, expectedVal := range expected {
		if parsed[key] != expectedVal {
			t.Errorf("migrated config[%q] = %q, want %q", key, parsed[key], expectedVal)
		}
	}
}

func TestMigrateConfigFileAlreadyToml(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")
	tomlContent := "worktree_dir = \"~/worktrees\"\n"
	if err := os.WriteFile(configPath, []byte(tomlContent), 0o644); err != nil {
		t.Fatalf("failed to write TOML config: %v", err)
	}

	migrated, err := migrateConfigFile(configPath)
	if err != nil {
		t.Fatalf("migrateConfigFile() error = %v", err)
	}
	if migrated {
		t.Fatal("expected migrateConfigFile() to skip TOML config")
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	if string(content) != tomlContent {
		t.Fatalf("expected TOML config to remain unchanged")
	}
}
