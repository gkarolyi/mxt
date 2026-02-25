package tmux

import (
	"reflect"
	"testing"
)

// TestParseLayout tests the custom layout parsing functionality.
func TestParseLayout(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Window
	}{
		{
			name:  "single window with no panes",
			input: "dev:",
			expected: []Window{
				{Name: "dev", Panes: []string{""}},
			},
		},
		{
			name:  "single window with one command",
			input: "dev:hx",
			expected: []Window{
				{Name: "dev", Panes: []string{"hx"}},
			},
		},
		{
			name:  "single window with two panes",
			input: "dev:hx|lazygit",
			expected: []Window{
				{Name: "dev", Panes: []string{"hx", "lazygit"}},
			},
		},
		{
			name:  "single window with three panes",
			input: "dev:hx|lazygit|",
			expected: []Window{
				{Name: "dev", Panes: []string{"hx", "lazygit", ""}},
			},
		},
		{
			name:  "multiple windows semicolon separator",
			input: "dev:hx|lazygit;server:bin/server;agent:",
			expected: []Window{
				{Name: "dev", Panes: []string{"hx", "lazygit"}},
				{Name: "server", Panes: []string{"bin/server"}},
				{Name: "agent", Panes: []string{""}},
			},
		},
		{
			name:  "multiple windows comma separator",
			input: "dev:hx|lazygit,server:bin/server,agent:",
			expected: []Window{
				{Name: "dev", Panes: []string{"hx", "lazygit"}},
				{Name: "server", Panes: []string{"bin/server"}},
				{Name: "agent", Panes: []string{""}},
			},
		},
		{
			name:  "multiple windows newline separator",
			input: "dev:hx|lazygit\nserver:bin/server\nagent:",
			expected: []Window{
				{Name: "dev", Panes: []string{"hx", "lazygit"}},
				{Name: "server", Panes: []string{"bin/server"}},
				{Name: "agent", Panes: []string{""}},
			},
		},
		{
			name:  "whitespace trimming",
			input: "  dev : hx | lazygit  ;  server : bin/server  ",
			expected: []Window{
				{Name: "dev", Panes: []string{"hx", "lazygit"}},
				{Name: "server", Panes: []string{"bin/server"}},
			},
		},
		{
			name:  "empty panes",
			input: "dev:||",
			expected: []Window{
				{Name: "dev", Panes: []string{"", "", ""}},
			},
		},
		{
			name:  "complex commands with shell metacharacters",
			input: "server:cd api && bin/server|cd ui && yarn start;logs:tail -f log/development.log",
			expected: []Window{
				{Name: "server", Panes: []string{"cd api && bin/server", "cd ui && yarn start"}},
				{Name: "logs", Panes: []string{"tail -f log/development.log"}},
			},
		},
		{
			name:  "empty windows ignored",
			input: "dev:hx;;server:bin/server",
			expected: []Window{
				{Name: "dev", Panes: []string{"hx"}},
				{Name: "server", Panes: []string{"bin/server"}},
			},
		},
		{
			name:  "mixed separators",
			input: "dev:hx|lazygit;server:bin/server\nagent:",
			expected: []Window{
				{Name: "dev", Panes: []string{"hx", "lazygit"}},
				{Name: "server", Panes: []string{"bin/server"}},
				{Name: "agent", Panes: []string{""}},
			},
		},
		{
			name:     "empty input",
			input:    "",
			expected: []Window{},
		},
		{
			name:     "only separators",
			input:    ";;;",
			expected: []Window{},
		},
		{
			name:  "window with colon in command",
			input: "dev:echo 'time: 10:30'",
			expected: []Window{
				{Name: "dev", Panes: []string{"echo 'time: 10:30'"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseLayout(tt.input)
			if err != nil {
				t.Fatalf("ParseLayout() returned error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseLayout() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

// TestParseLayoutErrors tests error cases for layout parsing.
func TestParseLayoutErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "window without name",
			input: ":hx|lazygit",
		},
		{
			name:  "missing colon separator",
			input: "dev hx lazygit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseLayout(tt.input)
			if err == nil {
				t.Errorf("ParseLayout() should return error for invalid input: %s", tt.input)
			}
		})
	}
}
