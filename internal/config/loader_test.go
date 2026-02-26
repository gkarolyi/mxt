package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestLoadDefaults tests loading default configuration
func TestLoadDefaults(t *testing.T) {
	config, err := LoadDefaults()
	if err != nil {
		t.Fatalf("LoadDefaults() error = %v", err)
	}

	// Check default values
	expectedDefaults := map[string]string{
		"worktree_dir":    filepath.Join(os.Getenv("HOME"), "worktrees"),
		"terminal":        "terminal",
		"copy_files":      "",
		"pre_session_cmd": "",
		"tmux_layout":     "",
	}

	for key, expectedVal := range expectedDefaults {
		if config[key] != expectedVal {
			t.Errorf("LoadDefaults()[%q] = %q, want %q", key, config[key], expectedVal)
		}
	}
}

// TestMergeConfigs tests merging two config maps
func TestMergeConfigs(t *testing.T) {
	tests := []struct {
		name     string
		base     map[string]string
		override map[string]string
		expected map[string]string
	}{
		{
			name: "override replaces values",
			base: map[string]string{
				"worktree_dir": "/base/path",
				"terminal":     "terminal",
			},
			override: map[string]string{
				"terminal": "iterm2",
			},
			expected: map[string]string{
				"worktree_dir": "/base/path",
				"terminal":     "iterm2",
			},
		},
		{
			name: "override adds new keys",
			base: map[string]string{
				"worktree_dir": "/base/path",
			},
			override: map[string]string{
				"terminal":   "iterm2",
				"copy_files": ".env",
			},
			expected: map[string]string{
				"worktree_dir": "/base/path",
				"terminal":     "iterm2",
				"copy_files":   ".env",
			},
		},
		{
			name: "empty override keeps base",
			base: map[string]string{
				"worktree_dir": "/base/path",
				"terminal":     "terminal",
			},
			override: map[string]string{},
			expected: map[string]string{
				"worktree_dir": "/base/path",
				"terminal":     "terminal",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeConfigs(tt.base, tt.override)
			for key, expectedVal := range tt.expected {
				if result[key] != expectedVal {
					t.Errorf("MergeConfigs()[%q] = %q, want %q", key, result[key], expectedVal)
				}
			}
		})
	}
}

// TestExpandTilde tests tilde expansion in worktree_dir
func TestExpandTilde(t *testing.T) {
	home := os.Getenv("HOME")
	if home == "" {
		t.Skip("HOME environment variable not set")
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "tilde at start",
			input:    "~/worktrees",
			expected: filepath.Join(home, "worktrees"),
		},
		{
			name:     "absolute path unchanged",
			input:    "/home/user/worktrees",
			expected: "/home/user/worktrees",
		},
		{
			name:     "tilde in middle unchanged",
			input:    "/path/to/~/worktrees",
			expected: "/path/to/~/worktrees",
		},
		{
			name:     "just tilde",
			input:    "~",
			expected: home,
		},
		{
			name:     "empty path",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandTilde(tt.input)
			if result != tt.expected {
				t.Errorf("ExpandTilde(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestLoadConfigFile tests loading a single config file
func TestLoadConfigFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")

	content := `worktree_dir = "~/test-worktrees"
terminal = "iterm2"
copy_files = ".env,.env.local"
`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Load the config file
	config, err := LoadConfigFile(configPath)
	if err != nil {
		t.Fatalf("LoadConfigFile() error = %v", err)
	}

	expected := map[string]string{
		"worktree_dir": "~/test-worktrees",
		"terminal":     "iterm2",
		"copy_files":   ".env,.env.local",
	}

	for key, expectedVal := range expected {
		if config[key] != expectedVal {
			t.Errorf("LoadConfigFile()[%q] = %q, want %q", key, config[key], expectedVal)
		}
	}
}

// TestLoadConfigFileNotExists tests loading non-existent file
func TestLoadConfigFileNotExists(t *testing.T) {
	config, err := LoadConfigFile("/nonexistent/path/config")
	if err != nil {
		t.Fatalf("LoadConfigFile() should not error for missing file, got: %v", err)
	}
	if len(config) != 0 {
		t.Errorf("LoadConfigFile() for missing file should return empty map, got: %v", config)
	}
}

// TestFindGitRoot tests finding the git repository root
func TestFindGitRoot(t *testing.T) {
	// This test requires running in a git repository
	// We'll test the actual implementation by checking current directory
	root, err := FindGitRoot(".")
	// If we're in a git repo, we should get a path
	// If not, we should get an error
	if err != nil {
		// Not in a git repo or git command failed - acceptable
		t.Logf("FindGitRoot() returned error (expected if not in git repo): %v", err)
		return
	}

	if root == "" {
		t.Errorf("FindGitRoot() returned empty string without error")
	}

	// Check if the returned path exists and has .git directory
	gitDir := filepath.Join(root, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		t.Errorf("FindGitRoot() returned %q but .git directory doesn't exist", root)
	}
}

// TestLoadConfig tests the full config loading with priority
func TestLoadConfig(t *testing.T) {
	// Create temporary directories
	tmpHome := t.TempDir()
	tmpRepo := t.TempDir()

	// Set up HOME for this test
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", oldHome)

	// Create global config
	globalConfigDir := filepath.Join(tmpHome, ".mxt")
	if err := os.MkdirAll(globalConfigDir, 0o755); err != nil {
		t.Fatalf("Failed to create global config dir: %v", err)
	}
	globalConfigPath := filepath.Join(globalConfigDir, "config")
	globalContent := `worktree_dir = "~/global-worktrees"
terminal = "terminal"
copy_files = ".env"
`
	if err := os.WriteFile(globalConfigPath, []byte(globalContent), 0o644); err != nil {
		t.Fatalf("Failed to create global config file: %v", err)
	}

	// Initialize git repo in tmpRepo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpRepo
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repo: %v", err)
	}

	// Create project config
	projectConfigPath := filepath.Join(tmpRepo, ".mxt")
	projectContent := `terminal = "iterm2"
copy_files = ".env,.env.local,CLAUDE.md"
`
	if err := os.WriteFile(projectConfigPath, []byte(projectContent), 0o644); err != nil {
		t.Fatalf("Failed to create project config file: %v", err)
	}

	// Test loading with project config
	config, err := LoadConfig(tmpRepo)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Expected: defaults overridden by global, overridden by project
	// worktree_dir: global (~/global-worktrees) with tilde expanded
	// terminal: project (iterm2)
	// copy_files: project (.env,.env.local,CLAUDE.md)
	// pre_session_cmd: default (empty)
	// tmux_layout: default (empty)

	expectedWorktreeDir := filepath.Join(tmpHome, "global-worktrees")
	if config["worktree_dir"] != expectedWorktreeDir {
		t.Errorf("LoadConfig()[worktree_dir] = %q, want %q", config["worktree_dir"], expectedWorktreeDir)
	}

	if config["terminal"] != "iterm2" {
		t.Errorf("LoadConfig()[terminal] = %q, want %q", config["terminal"], "iterm2")
	}

	if config["copy_files"] != ".env,.env.local,CLAUDE.md" {
		t.Errorf("LoadConfig()[copy_files] = %q, want %q", config["copy_files"], ".env,.env.local,CLAUDE.md")
	}

	if config["pre_session_cmd"] != "" {
		t.Errorf("LoadConfig()[pre_session_cmd] = %q, want empty", config["pre_session_cmd"])
	}
}

// TestLoadConfigOnlyDefaults tests loading when no config files exist
func TestLoadConfigOnlyDefaults(t *testing.T) {
	// Create temporary directory that's not a git repo and has no config
	tmpDir := t.TempDir()

	// Temporarily override HOME
	oldHome := os.Getenv("HOME")
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", oldHome)

	config, err := LoadConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Should have defaults with tilde expanded
	expectedWorktreeDir := filepath.Join(tmpHome, "worktrees")
	if config["worktree_dir"] != expectedWorktreeDir {
		t.Errorf("LoadConfig()[worktree_dir] = %q, want %q", config["worktree_dir"], expectedWorktreeDir)
	}

	if config["terminal"] != "terminal" {
		t.Errorf("LoadConfig()[terminal] = %q, want %q", config["terminal"], "terminal")
	}
}
