package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Default configuration values
const (
	DefaultTerminal      = "terminal"
	DefaultCopyFiles     = ""
	DefaultPreSessionCmd = ""
	DefaultTmuxLayout    = ""
)

// LoadDefaults returns the default configuration
func LoadDefaults() (map[string]string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return nil, fmt.Errorf("HOME environment variable not set")
	}

	return map[string]string{
		"worktree_dir":    filepath.Join(home, "worktrees"),
		"terminal":        DefaultTerminal,
		"copy_files":      DefaultCopyFiles,
		"pre_session_cmd": DefaultPreSessionCmd,
		"tmux_layout":     DefaultTmuxLayout,
	}, nil
}

// MergeConfigs merges two config maps, with override taking precedence
func MergeConfigs(base, override map[string]string) map[string]string {
	result := make(map[string]string)

	// Copy base
	for k, v := range base {
		result[k] = v
	}

	// Override with values from override map
	for k, v := range override {
		result[k] = v
	}

	return result
}

// ExpandTilde expands ~ at the start of a path to the user's home directory
func ExpandTilde(path string) string {
	if path == "" {
		return ""
	}

	if path == "~" {
		return os.Getenv("HOME")
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(os.Getenv("HOME"), path[2:])
	}

	return path
}

// LoadConfigFile loads a config file from the given path.
// Returns an empty map (not an error) if the file doesn't exist.
func LoadConfigFile(path string) (map[string]string, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return make(map[string]string), nil
	}

	// Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %w", path, err)
	}
	defer file.Close()

	// Parse config
	config, err := ParseConfig(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	// Validate security
	if err := ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("config file %s failed security validation: %w", path, err)
	}

	return config, nil
}

// FindGitRoot finds the root directory of the git repository containing the given path.
// Returns an error if the path is not in a git repository.
func FindGitRoot(path string) (string, error) {
	// Run git rev-parse --show-toplevel
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = path

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not in a git repository or git command failed: %w", err)
	}

	root := strings.TrimSpace(string(output))
	return root, nil
}

// GetGlobalConfigPath returns the path to the global config file.
// Respects MXT_CONFIG_DIR environment variable.
func GetGlobalConfigPath() string {
	configDir := os.Getenv("MXT_CONFIG_DIR")
	if configDir == "" {
		configDir = filepath.Join(os.Getenv("HOME"), ".mxt")
	}
	return filepath.Join(configDir, "config")
}

// LoadConfig loads configuration with proper priority:
// 1. Load defaults
// 2. Load global config (if exists) - overrides defaults
// 3. Detect git repo and load project config (if exists) - overrides global
// 4. Expand tilde in worktree_dir
func LoadConfig(workDir string) (map[string]string, error) {
	// Start with defaults
	config, err := LoadDefaults()
	if err != nil {
		return nil, err
	}

	// Load global config
	globalConfigPath := GetGlobalConfigPath()
	globalConfig, err := LoadConfigFile(globalConfigPath)
	if err != nil {
		return nil, err
	}
	config = MergeConfigs(config, globalConfig)

	// Try to find git root and load project config
	gitRoot, err := FindGitRoot(workDir)
	if err == nil {
		// We're in a git repo, try to load project config
		projectConfigPath := filepath.Join(gitRoot, ".mxt")
		projectConfig, err := LoadConfigFile(projectConfigPath)
		if err != nil {
			return nil, err
		}
		config = MergeConfigs(config, projectConfig)
	}
	// If not in a git repo, that's okay - just use global config

	// Expand tilde in worktree_dir
	config["worktree_dir"] = ExpandTilde(config["worktree_dir"])

	return config, nil
}
