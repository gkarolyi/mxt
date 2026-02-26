package config

import (
	"strings"
	"testing"
)

func TestParseConfigTOML(t *testing.T) {
	input := `# Example config
worktree_dir = "/home/user/worktrees"
terminal = "iterm2"
copy_files = ".env,.env.local"
pre_session_cmd = "echo \"hello # world\""
`

	config, err := ParseConfig(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseConfig() error = %v", err)
	}

	expected := map[string]string{
		"worktree_dir":    "/home/user/worktrees",
		"terminal":        "iterm2",
		"copy_files":      ".env,.env.local",
		"pre_session_cmd": "echo \"hello # world\"",
	}

	for key, expectedVal := range expected {
		if config[key] != expectedVal {
			t.Errorf("ParseConfig()[%q] = %q, want %q", key, config[key], expectedVal)
		}
	}
}

func TestParseConfigArrayValues(t *testing.T) {
	input := `copy_files = [".env", ".env.local"]
	tmux_layout = ["dev:hx|lazygit", "server:bin/server", "agent:"]`
	config, err := ParseConfig(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseConfig() error = %v", err)
	}
	if config["copy_files"] != ".env,.env.local" {
		t.Errorf("ParseConfig()[copy_files] = %q, want %q", config["copy_files"], ".env,.env.local")
	}
	expectedLayout := "dev:hx|lazygit;server:bin/server;agent:"
	if config["tmux_layout"] != expectedLayout {
		t.Errorf("ParseConfig()[tmux_layout] = %q, want %q", config["tmux_layout"], expectedLayout)
	}
}

func TestParseConfigUnknownKey(t *testing.T) {
	_, err := ParseConfig(strings.NewReader("unknown = \"value\""))
	if err == nil {
		t.Fatal("ParseConfig() expected error for unknown key")
	}
}

func TestParseConfigTmuxLayoutMultiline(t *testing.T) {
	input := `tmux_layout = """
  dev:hx|lazygit
  server:bin/server
  agent:
  """`

	config, err := ParseConfig(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseConfig() error = %v", err)
	}

	expected := "dev:hx|lazygit;server:bin/server;agent:"
	if config["tmux_layout"] != expected {
		t.Errorf("ParseConfig()[tmux_layout] = %q, want %q", config["tmux_layout"], expected)
	}
}

func TestParseConfigEmptyInput(t *testing.T) {
	config, err := ParseConfig(strings.NewReader(""))
	if err != nil {
		t.Fatalf("ParseConfig() error = %v", err)
	}
	if len(config) != 0 {
		t.Errorf("ParseConfig() = %v, want empty map", config)
	}
}

func TestEncodeConfigRejectsUnknownKey(t *testing.T) {
	_, err := EncodeConfig(map[string]string{"unknown": "value"})
	if err == nil {
		t.Fatal("EncodeConfig() expected error for unknown key")
	}
}
