package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gkarolyi/mxt/internal/config"
	"github.com/gkarolyi/mxt/internal/ui"
)

const logo = `                       _
  _ __ ___  _   ___  _| |_ _ __ ___  ___
 | '_ ` + "`" + ` _ \| | | \ \/ / __| '__/ _ \/ _ \
 | | | | | | |_| |>  <| |_| | |  __/  __/
 |_| |_| |_|\__,_/_/\_\\__|_|  \___|\___|
  Tmux Worktree Session Manager v1.0.0
`

// InitCommand implements the init command
func InitCommand(local bool) error {
	// Display logo
	fmt.Print(logo)
	fmt.Println()

	if local {
		return initProjectConfig()
	}
	return initGlobalConfig()
}

func initGlobalConfig() error {
	configPath := config.GetGlobalConfigPath()
	configDir := filepath.Dir(configPath)

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		// Config exists, prompt for overwrite
		content, _ := os.ReadFile(configPath)
		fmt.Println(string(content))
		fmt.Println()
		fmt.Print("Config already exists. Overwrite? (y/N): ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			ui.Info("Keeping existing config")
			return nil
		}
		fmt.Println()
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Prompt for config values
	reader := bufio.NewReader(os.Stdin)

	// Worktree directory and terminal use defaults (no prompt shown in muxtree)
	worktreeDir := "~/worktrees"
	terminal := "terminal"

	// Copy files
	fmt.Println("▸ Enter files to copy into new worktrees (relative to repo root).")
	fmt.Println("▸ Comma-separated, e.g.: .env,.env.local,CLAUDE.md")
	copyFiles, _ := reader.ReadString('\n')
	copyFiles = strings.TrimSpace(copyFiles)

	// Pre-session command
	fmt.Println()
	fmt.Println("▸ Optional: Command to run after worktree setup, before tmux session.")
	fmt.Println("▸ Runs in worktree dir. Good for: bundle install, npm install, db:migrate")
	preSessionCmd, _ := reader.ReadString('\n')
	preSessionCmd = strings.TrimSpace(preSessionCmd)

	// Tmux layout
	fmt.Println()
	fmt.Println("▸ Optional: Tmux layout - define windows and panes.")
	fmt.Println("▸ Format: window:cmd1|cmd2;window2:cmd3")
	fmt.Println("▸ Example: dev:vim|;server:bin/server;agent:")
	tmuxLayout, _ := reader.ReadString('\n')
	tmuxLayout = strings.TrimSpace(tmuxLayout)

	// Generate config file content
	content := generateGlobalConfigContent(worktreeDir, terminal, copyFiles, preSessionCmd, tmuxLayout)

	// Write config file
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Display success message
	ui.Success(fmt.Sprintf("Config written to %s", configPath))
	fmt.Println()
	fmt.Println(content)

	return nil
}

func initProjectConfig() error {
	// Find git root
	gitRoot, err := config.FindGitRoot(".")
	if err != nil {
		return fmt.Errorf("must be in a git repository to create project config: %w", err)
	}

	configPath := filepath.Join(gitRoot, ".muxtree")

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		// Config exists, prompt for overwrite
		content, _ := os.ReadFile(configPath)
		fmt.Println(string(content))
		fmt.Println()
		fmt.Print("Config already exists. Overwrite? (y/N): ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			ui.Info("Keeping existing config")
			return nil
		}
		fmt.Println()
	}

	// Prompt for config values (project-specific only)
	reader := bufio.NewReader(os.Stdin)

	// Copy files
	fmt.Println("▸ Enter files to copy into new worktrees for this project (relative to repo root).")
	fmt.Println("▸ Comma-separated, e.g.: .env,.env.local,CLAUDE.md")
	copyFiles, _ := reader.ReadString('\n')
	copyFiles = strings.TrimSpace(copyFiles)

	// Pre-session command
	fmt.Println()
	fmt.Println("▸ Optional: Command to run after worktree setup, before tmux session.")
	fmt.Println("▸ Runs in worktree dir. Good for: bundle install, npm install, db:migrate")
	preSessionCmd, _ := reader.ReadString('\n')
	preSessionCmd = strings.TrimSpace(preSessionCmd)

	// Tmux layout
	fmt.Println()
	fmt.Println("▸ Optional: Tmux layout - define windows and panes.")
	fmt.Println("▸ Format: window:cmd1|cmd2;window2:cmd3")
	fmt.Println("▸ Example: dev:vim|;server:bin/server;agent:")
	tmuxLayout, _ := reader.ReadString('\n')
	tmuxLayout = strings.TrimSpace(tmuxLayout)

	// Generate config file content
	content := generateProjectConfigContent(copyFiles, preSessionCmd, tmuxLayout)

	// Write config file
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Display success message
	ui.Success(fmt.Sprintf("Project config written to %s", configPath))
	fmt.Println()
	fmt.Println(content)

	return nil
}

func generateGlobalConfigContent(worktreeDir, terminal, copyFiles, preSessionCmd, tmuxLayout string) string {
	timestamp := time.Now().Format("Mon 02 Jan 2006 15:04:05 MST")

	var sb strings.Builder
	sb.WriteString("# muxtree configuration\n")
	sb.WriteString(fmt.Sprintf("# Generated on %s\n\n", timestamp))

	sb.WriteString("# Base directory for worktrees\n")
	sb.WriteString(fmt.Sprintf("worktree_dir=%s\n\n", worktreeDir))

	sb.WriteString("# Terminal app: terminal | iterm2 | ghostty | current\n")
	sb.WriteString(fmt.Sprintf("terminal=%s\n\n", terminal))

	sb.WriteString("# Files to copy from repo root into new worktrees (comma-separated, relative to repo root)\n")
	sb.WriteString("# Supports glob patterns and directories\n")
	sb.WriteString(fmt.Sprintf("copy_files=%s\n\n", copyFiles))

	sb.WriteString("# Command to run after worktree setup, before tmux session (optional)\n")
	sb.WriteString("# Runs in worktree directory. Use for setup tasks like: bundle install, npm install\n")
	sb.WriteString(fmt.Sprintf("pre_session_cmd=%s\n\n", preSessionCmd))

	sb.WriteString("# Tmux layout - define windows and panes (optional)\n")
	sb.WriteString("# Format: window_name:pane_cmd1|pane_cmd2;next_window:cmd\n")
	sb.WriteString("# Example: dev:vim|;server:bin/server;agent:\n")
	sb.WriteString("# - ';' separates windows\n")
	sb.WriteString("# - ':' separates window name from panes\n")
	sb.WriteString("# - '|' separates panes (horizontal split)\n")
	sb.WriteString("# - Empty command = shell prompt\n")
	sb.WriteString("# If not set, creates default layout: dev + agent windows\n")
	sb.WriteString(fmt.Sprintf("tmux_layout=%s\n", tmuxLayout))

	return sb.String()
}

func generateProjectConfigContent(copyFiles, preSessionCmd, tmuxLayout string) string {
	timestamp := time.Now().Format("Mon 02 Jan 2006 15:04:05 MST")

	var sb strings.Builder
	sb.WriteString("# muxtree project config\n")
	sb.WriteString(fmt.Sprintf("# Generated on %s\n\n", timestamp))

	sb.WriteString("# Files to copy from repo root into new worktrees (comma-separated, relative to repo root)\n")
	sb.WriteString("# Supports glob patterns and directories\n")
	sb.WriteString(fmt.Sprintf("copy_files=%s\n\n", copyFiles))

	sb.WriteString("# Command to run after worktree setup, before tmux session (optional)\n")
	sb.WriteString("# Runs in worktree directory. Use for setup tasks like: bundle install, npm install\n")
	sb.WriteString(fmt.Sprintf("pre_session_cmd=%s\n\n", preSessionCmd))

	sb.WriteString("# Tmux layout - define windows and panes (optional)\n")
	sb.WriteString("# Format: window_name:pane_cmd1|pane_cmd2;next_window:cmd\n")
	sb.WriteString("# Example: dev:vim|;server:bin/server;agent:\n")
	sb.WriteString("# - ';' separates windows\n")
	sb.WriteString("# - ':' separates window name from panes\n")
	sb.WriteString("# - '|' separates panes (horizontal split)\n")
	sb.WriteString("# - Empty command = shell prompt\n")
	sb.WriteString("# If not set, creates default layout: dev + agent windows\n")
	sb.WriteString(fmt.Sprintf("tmux_layout=%s\n", tmuxLayout))

	return sb.String()
}
