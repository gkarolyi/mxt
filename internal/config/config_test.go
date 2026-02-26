package config

import (
	"strings"
	"testing"
)

// TestParseSingleLineKeyValue tests basic key=value parsing
func TestParseSingleLineKeyValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:  "simple key-value",
			input: "worktree_dir=/home/user/worktrees",
			expected: map[string]string{
				"worktree_dir": "/home/user/worktrees",
			},
		},
		{
			name:  "multiple key-values",
			input: "worktree_dir=/home/user/worktrees\nterminal=iterm2\ncopy_files=.env,.env.local",
			expected: map[string]string{
				"worktree_dir": "/home/user/worktrees",
				"terminal":     "iterm2",
				"copy_files":   ".env,.env.local",
			},
		},
		{
			name:  "tilde in path",
			input: "worktree_dir=~/worktrees",
			expected: map[string]string{
				"worktree_dir": "~/worktrees",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseLegacyConfig(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("ParseLegacyConfig() error = %v", err)
			}
			for key, expectedVal := range tt.expected {
				if result[key] != expectedVal {
					t.Errorf("ParseLegacyConfig()[%q] = %q, want %q", key, result[key], expectedVal)
				}
			}
		})
	}
}

// TestParseConfigComments tests comment handling
func TestParseConfigComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:     "line starting with hash",
			input:    "# This is a comment\nworktree_dir=/home/user",
			expected: map[string]string{"worktree_dir": "/home/user"},
		},
		{
			name:     "line with leading whitespace and hash",
			input:    "  # Indented comment\nworktree_dir=/home/user",
			expected: map[string]string{"worktree_dir": "/home/user"},
		},
		{
			name:     "multiple comments",
			input:    "# Comment 1\n# Comment 2\nworktree_dir=/home/user\n# Comment 3",
			expected: map[string]string{"worktree_dir": "/home/user"},
		},
		{
			name:     "hash in value not treated as comment",
			input:    "pre_session_cmd=echo 'hello # world'",
			expected: map[string]string{"pre_session_cmd": "echo 'hello # world'"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseLegacyConfig(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("ParseLegacyConfig() error = %v", err)
			}
			for key, expectedVal := range tt.expected {
				if result[key] != expectedVal {
					t.Errorf("ParseLegacyConfig()[%q] = %q, want %q", key, result[key], expectedVal)
				}
			}
		})
	}
}

// TestParseConfigWhitespace tests whitespace trimming
func TestParseConfigWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:     "leading whitespace in key",
			input:    "  worktree_dir=/home/user",
			expected: map[string]string{"worktree_dir": "/home/user"},
		},
		{
			name:     "trailing whitespace in value",
			input:    "worktree_dir=/home/user  ",
			expected: map[string]string{"worktree_dir": "/home/user"},
		},
		{
			name:     "whitespace around equals",
			input:    "worktree_dir = /home/user",
			expected: map[string]string{"worktree_dir": "/home/user"},
		},
		{
			name:     "tabs and spaces mixed",
			input:    "\tworktree_dir\t=\t/home/user\t",
			expected: map[string]string{"worktree_dir": "/home/user"},
		},
		{
			name:     "preserve internal whitespace in value",
			input:    "pre_session_cmd=npm install && npm run build",
			expected: map[string]string{"pre_session_cmd": "npm install && npm run build"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseLegacyConfig(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("ParseLegacyConfig() error = %v", err)
			}
			for key, expectedVal := range tt.expected {
				if result[key] != expectedVal {
					t.Errorf("ParseLegacyConfig()[%q] = %q, want %q", key, result[key], expectedVal)
				}
			}
		})
	}
}

// TestParseConfigEmptyLines tests empty line handling
func TestParseConfigEmptyLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:     "single empty line",
			input:    "worktree_dir=/home/user\n\nterminal=iterm2",
			expected: map[string]string{"worktree_dir": "/home/user", "terminal": "iterm2"},
		},
		{
			name:     "multiple empty lines",
			input:    "\n\nworktree_dir=/home/user\n\n\nterminal=iterm2\n\n",
			expected: map[string]string{"worktree_dir": "/home/user", "terminal": "iterm2"},
		},
		{
			name:     "whitespace-only lines",
			input:    "   \nworktree_dir=/home/user\n  \t  \nterminal=iterm2",
			expected: map[string]string{"worktree_dir": "/home/user", "terminal": "iterm2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseLegacyConfig(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("ParseLegacyConfig() error = %v", err)
			}
			for key, expectedVal := range tt.expected {
				if result[key] != expectedVal {
					t.Errorf("ParseLegacyConfig()[%q] = %q, want %q", key, result[key], expectedVal)
				}
			}
		})
	}
}

// TestParseConfigInvalidFormat tests error handling for malformed input
func TestParseConfigInvalidFormat(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
	}{
		{
			name:        "line without equals",
			input:       "worktree_dir",
			shouldError: true,
		},
		{
			name:        "empty key",
			input:       "=/home/user",
			shouldError: true,
		},
		{
			name:        "multiple equals",
			input:       "key=value=extra",
			shouldError: false, // value should be "value=extra"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseLegacyConfig(strings.NewReader(tt.input))
			if tt.shouldError && err == nil {
				t.Errorf("ParseLegacyConfig() expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("ParseLegacyConfig() unexpected error = %v", err)
			}
		})
	}
}

// TestParseConfigEmpty tests empty input
func TestParseConfigEmpty(t *testing.T) {
	result, err := ParseLegacyConfig(strings.NewReader(""))
	if err != nil {
		t.Fatalf("ParseLegacyConfig() error = %v", err)
	}
	if len(result) != 0 {
		t.Errorf("ParseLegacyConfig() = %v, want empty map", result)
	}
}

// TestParseMultiLineArray tests multi-line array parsing with key=[...]
func TestParseMultiLineArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name: "multi-line array basic",
			input: `pre_session_cmd=[
  npm install
  npm run db:migrate
]`,
			expected: map[string]string{
				"pre_session_cmd": "npm install npm run db:migrate",
			},
		},
		{
			name: "multi-line tmux_layout",
			input: `tmux_layout=[
  dev:hx|lazygit
  server:cd api && bin/server
  agent:
]`,
			expected: map[string]string{
				"tmux_layout": "dev:hx|lazygit;server:cd api && bin/server;agent:",
			},
		},
		{
			name:  "single-line array format",
			input: `copy_files=[.env .env.local CLAUDE.md]`,
			expected: map[string]string{
				"copy_files": ".env .env.local CLAUDE.md",
			},
		},
		{
			name: "multi-line with extra whitespace",
			input: `pre_session_cmd=[
    npm install
    npm run build
  ]`,
			expected: map[string]string{
				"pre_session_cmd": "npm install npm run build",
			},
		},
		{
			name: "multi-line array with comments between",
			input: `pre_session_cmd=[
  npm install
  # This is a comment
  npm run build
]`,
			expected: map[string]string{
				"pre_session_cmd": "npm install npm run build",
			},
		},
		{
			name: "multi-line and single-line mixed",
			input: `worktree_dir=~/worktrees
pre_session_cmd=[
  npm install
  npm run build
]
terminal=iterm2`,
			expected: map[string]string{
				"worktree_dir":    "~/worktrees",
				"pre_session_cmd": "npm install npm run build",
				"terminal":        "iterm2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseLegacyConfig(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("ParseLegacyConfig() error = %v", err)
			}
			for key, expectedVal := range tt.expected {
				if result[key] != expectedVal {
					t.Errorf("ParseLegacyConfig()[%q] = %q, want %q", key, result[key], expectedVal)
				}
			}
		})
	}
}

// TestParseTmuxLayoutNormalization tests separator normalization for tmux_layout
func TestParseTmuxLayoutNormalization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "commas to semicolons",
			input:    `tmux_layout=dev:hx|lazygit,server:bin/server,agent:`,
			expected: "dev:hx|lazygit;server:bin/server;agent:",
		},
		{
			name: "multi-line with spaces becomes semicolons",
			input: `tmux_layout=[
  dev:hx|lazygit
  server:bin/server
  agent:
]`,
			expected: "dev:hx|lazygit;server:bin/server;agent:",
		},
		{
			name:     "single-line array format",
			input:    `tmux_layout=[dev:hx|lazygit server:bin/server agent:]`,
			expected: "dev:hx|lazygit;server:bin/server;agent:",
		},
		{
			name:     "mixed separators normalized",
			input:    `tmux_layout=dev:hx,server:bin/server,agent:`,
			expected: "dev:hx;server:bin/server;agent:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseLegacyConfig(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("ParseLegacyConfig() error = %v", err)
			}
			if result["tmux_layout"] != tt.expected {
				t.Errorf("ParseLegacyConfig()[tmux_layout] = %q, want %q", result["tmux_layout"], tt.expected)
			}
		})
	}
}

// TestParseMultiLineArrayErrors tests error handling for malformed multi-line arrays
func TestParseMultiLineArrayErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
	}{
		{
			name: "unclosed bracket",
			input: `pre_session_cmd=[
  npm install
  npm run build`,
			shouldError: true,
		},
		{
			name:        "opening bracket without key",
			input:       "=[value]",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseLegacyConfig(strings.NewReader(tt.input))
			if tt.shouldError && err == nil {
				t.Errorf("ParseLegacyConfig() expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("ParseLegacyConfig() unexpected error = %v", err)
			}
		})
	}
}
