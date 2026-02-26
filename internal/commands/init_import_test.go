package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gkarolyi/mxt/internal/config"
)

func TestImportLegacyConfigWritesToml(t *testing.T) {
	tmpDir := t.TempDir()
	legacyPath := filepath.Join(tmpDir, "config")
	legacyContent := "worktree_dir=~/worktrees\nterminal=iterm2\ncopy_files=.env,.env.local\npre_session_cmd=echo hi\ntmux_layout=dev:hx|lazygit,server:bin/server,agent:\n"
	if err := os.WriteFile(legacyPath, []byte(legacyContent), 0o644); err != nil {
		t.Fatalf("failed to write legacy config: %v", err)
	}

	targetPath := filepath.Join(tmpDir, "config.toml")
	encoded, err := importLegacyConfig(legacyPath, targetPath, false)
	if err != nil {
		t.Fatalf("importLegacyConfig() error = %v", err)
	}
	if _, err := os.Stat(targetPath); err != nil {
		t.Fatalf("expected target config to exist: %v", err)
	}

	parsed, err := config.ParseConfig(strings.NewReader(encoded))
	if err != nil {
		t.Fatalf("ParseConfig() error = %v", err)
	}

	expected := map[string]string{
		"worktree_dir":    "~/worktrees",
		"terminal":        "iterm2",
		"copy_files":      ".env,.env.local",
		"pre_session_cmd": "echo hi",
		"tmux_layout":     "dev:hx|lazygit;server:bin/server;agent:",
	}
	for key, expectedValue := range expected {
		if parsed[key] != expectedValue {
			t.Errorf("imported config[%s] = %q, want %q", key, parsed[key], expectedValue)
		}
	}
}

func TestImportLegacyConfigMissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	targetPath := filepath.Join(tmpDir, "config.toml")
	_, err := importLegacyConfig(filepath.Join(tmpDir, "missing"), targetPath, false)
	if err == nil {
		t.Fatal("expected error when legacy config is missing")
	}
	if !strings.Contains(err.Error(), "legacy config not found") {
		t.Fatalf("unexpected error: %v", err)
	}
}
