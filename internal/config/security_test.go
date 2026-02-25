package config

import (
	"testing"
)

// TestValidateConfigValueSecurity tests security validation for config values
func TestValidateConfigValueSecurity(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		value       string
		shouldError bool
	}{
		// Safe values that should pass
		{
			name:        "safe worktree_dir",
			key:         "worktree_dir",
			value:       "~/worktrees",
			shouldError: false,
		},
		{
			name:        "safe terminal",
			key:         "terminal",
			value:       "iterm2",
			shouldError: false,
		},
		{
			name:        "safe copy_files",
			key:         "copy_files",
			value:       ".env,.env.local,CLAUDE.md",
			shouldError: false,
		},

		// Dangerous values for non-command keys
		{
			name:        "backtick in worktree_dir",
			key:         "worktree_dir",
			value:       "~/worktrees`whoami`",
			shouldError: true,
		},
		{
			name:        "dollar in worktree_dir",
			key:         "worktree_dir",
			value:       "~/worktrees$(whoami)",
			shouldError: true,
		},
		{
			name:        "semicolon in terminal",
			key:         "terminal",
			value:       "iterm2;rm -rf /",
			shouldError: true,
		},
		{
			name:        "pipe in copy_files",
			key:         "copy_files",
			value:       ".env|cat /etc/passwd",
			shouldError: true,
		},
		{
			name:        "ampersand in worktree_dir",
			key:         "worktree_dir",
			value:       "~/worktrees && rm -rf /",
			shouldError: true,
		},

		// Command keys should allow metacharacters
		{
			name:        "backtick in pre_session_cmd",
			key:         "pre_session_cmd",
			value:       "echo `date`",
			shouldError: false,
		},
		{
			name:        "dollar in pre_session_cmd",
			key:         "pre_session_cmd",
			value:       "echo $HOME",
			shouldError: false,
		},
		{
			name:        "semicolon in pre_session_cmd",
			key:         "pre_session_cmd",
			value:       "npm install; npm run build",
			shouldError: false,
		},
		{
			name:        "pipe in tmux_layout",
			key:         "tmux_layout",
			value:       "dev:hx|lazygit",
			shouldError: false,
		},
		{
			name:        "ampersand in pre_session_cmd",
			key:         "pre_session_cmd",
			value:       "npm install && npm run build",
			shouldError: false,
		},
		{
			name:        "complex command in tmux_layout",
			key:         "tmux_layout",
			value:       "server:cd api && bin/server",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfigValue(tt.key, tt.value)
			if tt.shouldError && err == nil {
				t.Errorf("ValidateConfigValue(%q, %q) expected error but got none", tt.key, tt.value)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("ValidateConfigValue(%q, %q) unexpected error = %v", tt.key, tt.value, err)
			}
		})
	}
}

// TestIsCommandKey tests the helper function that identifies command keys
func TestIsCommandKey(t *testing.T) {
	tests := []struct {
		key          string
		isCommandKey bool
	}{
		{"worktree_dir", false},
		{"terminal", false},
		{"copy_files", false},
		{"pre_session_cmd", true},
		{"tmux_layout", true},
		{"unknown_key", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := isCommandKey(tt.key)
			if result != tt.isCommandKey {
				t.Errorf("isCommandKey(%q) = %v, want %v", tt.key, result, tt.isCommandKey)
			}
		})
	}
}

// TestContainsMetacharacters tests detection of shell metacharacters
func TestContainsMetacharacters(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"safe path", "~/worktrees", false},
		{"safe filename", ".env.local", false},
		{"safe command name", "iterm2", false},
		{"with backtick", "test`whoami`", true},
		{"with dollar", "test$(cmd)", true},
		{"with semicolon", "test;cmd", true},
		{"with pipe", "test|cmd", true},
		{"with ampersand", "test && cmd", true},
		{"with single dollar at end", "test$", true},
		{"multiple metacharacters", "test`cmd`; rm -rf /", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsMetacharacters(tt.value)
			if result != tt.expected {
				t.Errorf("containsMetacharacters(%q) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

// TestValidateConfig tests validation of entire config map
func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      map[string]string
		shouldError bool
		errorKey    string
	}{
		{
			name: "all safe values",
			config: map[string]string{
				"worktree_dir": "~/worktrees",
				"terminal":     "iterm2",
				"copy_files":   ".env,.env.local",
			},
			shouldError: false,
		},
		{
			name: "command keys with metacharacters",
			config: map[string]string{
				"pre_session_cmd": "npm install && npm run build",
				"tmux_layout":     "dev:hx|lazygit",
			},
			shouldError: false,
		},
		{
			name: "dangerous value in worktree_dir",
			config: map[string]string{
				"worktree_dir": "~/worktrees`whoami`",
				"terminal":     "iterm2",
			},
			shouldError: true,
			errorKey:    "worktree_dir",
		},
		{
			name: "dangerous value in terminal",
			config: map[string]string{
				"worktree_dir": "~/worktrees",
				"terminal":     "iterm2;rm -rf /",
			},
			shouldError: true,
			errorKey:    "terminal",
		},
		{
			name: "mixed safe and dangerous",
			config: map[string]string{
				"worktree_dir":    "~/worktrees",
				"terminal":        "iterm2",
				"copy_files":      ".env|cat /etc/passwd",
				"pre_session_cmd": "npm install",
			},
			shouldError: true,
			errorKey:    "copy_files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if tt.shouldError {
				if err == nil {
					t.Errorf("ValidateConfig() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateConfig() unexpected error = %v", err)
				}
			}
		})
	}
}
