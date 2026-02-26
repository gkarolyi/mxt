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
	"golang.org/x/term"
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
	reader := bufio.NewReader(os.Stdin)

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		ui.Warn(fmt.Sprintf("Config already exists at %s", configPath))
		fmt.Println()
		content, _ := os.ReadFile(configPath)
		fmt.Println(string(content))
		fmt.Println()
		if term.IsTerminal(int(os.Stdin.Fd())) {
			fmt.Print("Overwrite? (y/N) ")
		}
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)
		if response != "y" && response != "Y" {
			return nil
		}
		fmt.Println()
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	defaults, err := config.LoadDefaults()
	if err != nil {
		return fmt.Errorf("failed to load defaults: %w", err)
	}

	worktreeDir := defaults["worktree_dir"]
	terminal := defaults["terminal"]

	fmt.Println()

	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Printf("Worktree base directory [%s]: ", worktreeDir)
	}
	inputDir, _ := reader.ReadString('\n')
	inputDir = strings.TrimSpace(inputDir)
	if inputDir != "" {
		worktreeDir = inputDir
	}

	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Printf("Terminal app (terminal/iterm2/ghostty/current) [%s]: ", terminal)
	}
	inputTerm, _ := reader.ReadString('\n')
	inputTerm = strings.TrimSpace(inputTerm)
	if inputTerm != "" {
		terminal = inputTerm
	}

	fmt.Println()

	// Copy files
	ui.Info("Enter files to copy into new worktrees (relative to repo root).")
	ui.Info("Comma-separated, e.g.: .env,.env.local,CLAUDE.md")
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Print("Files to copy: ")
	}
	copyFiles, _ := reader.ReadString('\n')
	copyFiles = strings.TrimSpace(copyFiles)

	// Pre-session command
	fmt.Println()
	ui.Info("Optional: Command to run after worktree setup, before tmux session.")
	ui.Info("Runs in worktree dir. Good for: bundle install, npm install, db:migrate")
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Print("Pre-session command: ")
	}
	preSessionCmd, _ := reader.ReadString('\n')
	preSessionCmd = strings.TrimSpace(preSessionCmd)

	// Tmux layout
	fmt.Println()
	ui.Info("Optional: Tmux layout - define windows and panes.")
	ui.Info("You can define this now (single line) or edit the config file for multi-line format.")
	ui.Info("Example: dev:hx|lazygit,server:bin/server,agent:")
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Print("Tmux layout: ")
	}
	tmuxLayout, _ := reader.ReadString('\n')
	tmuxLayout = strings.TrimSpace(tmuxLayout)

	// Generate config file content
	content := generateGlobalConfigContent(worktreeDir, terminal, copyFiles, preSessionCmd, tmuxLayout)

	// Write config file
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Display success message
	ui.Success(fmt.Sprintf("Config written to %s", configPath))
	fmt.Println()
	fmt.Print(content)

	return nil
}

func initProjectConfig() error {
	gitRoot, err := config.FindGitRoot(".")
	if err != nil {
		return fmt.Errorf("Not inside a git repository. Run muxtree from within your repo.")
	}

	configPath := filepath.Join(gitRoot, ".muxtree")
	reader := bufio.NewReader(os.Stdin)

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		ui.Warn(fmt.Sprintf("Project config already exists at %s", configPath))
		fmt.Println()
		content, _ := os.ReadFile(configPath)
		fmt.Println(string(content))
		fmt.Println()
		if term.IsTerminal(int(os.Stdin.Fd())) {
			fmt.Print("Overwrite? (y/N) ")
		}
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)
		if response != "y" && response != "Y" {
			return nil
		}
		fmt.Println()
	}

	fmt.Println()

	// Copy files
	ui.Info("Enter files to copy into new worktrees for this project (relative to repo root).")
	ui.Info("Comma-separated, e.g.: .env,.env.local,CLAUDE.md")
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Print("Files to copy: ")
	}
	copyFiles, _ := reader.ReadString('\n')
	copyFiles = strings.TrimSpace(copyFiles)

	// Pre-session command
	fmt.Println()
	ui.Info("Optional: Command to run after worktree setup, before tmux session.")
	ui.Info("Runs in worktree dir. Good for: bundle install, npm install, db:migrate")
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Print("Pre-session command: ")
	}
	preSessionCmd, _ := reader.ReadString('\n')
	preSessionCmd = strings.TrimSpace(preSessionCmd)

	// Tmux layout
	fmt.Println()
	ui.Info("Optional: Tmux layout - define windows and panes.")
	ui.Info("You can define this now (single line) or edit the config file for multi-line format.")
	ui.Info("Example: dev:hx|lazygit,server:bin/server,agent:")
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Print("Tmux layout: ")
	}
	tmuxLayout, _ := reader.ReadString('\n')
	tmuxLayout = strings.TrimSpace(tmuxLayout)

	// Generate config file content
	content := generateProjectConfigContent(copyFiles, preSessionCmd, tmuxLayout)

	// Write config file
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Display success message
	ui.Success(fmt.Sprintf("Project config written to %s", configPath))
	fmt.Println()
	fmt.Print(content)

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
	sb.WriteString("# Multi-line format (more readable):\n")
	sb.WriteString("# tmux_layout=[\n")
	sb.WriteString("#   dev:hx|lazygit\n")
	sb.WriteString("#   server:bin/server\n")
	sb.WriteString("#   agent:\n")
	sb.WriteString("# ]\n")
	sb.WriteString("# Or single line: tmux_layout=dev:hx|lazygit,server:bin/server,agent:\n")
	sb.WriteString("#\n")
	sb.WriteString("# Syntax:\n")
	sb.WriteString("# - ',' or newline separates windows\n")
	sb.WriteString("# - ':' separates window name from panes\n")
	sb.WriteString("# - '|' separates panes (vertical split - side by side)\n")
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
	sb.WriteString("# Multi-line format (more readable):\n")
	sb.WriteString("# tmux_layout=[\n")
	sb.WriteString("#   dev:hx|lazygit\n")
	sb.WriteString("#   server:bin/server\n")
	sb.WriteString("#   agent:\n")
	sb.WriteString("# ]\n")
	sb.WriteString("# Or single line: tmux_layout=dev:hx|lazygit,server:bin/server,agent:\n")
	sb.WriteString("#\n")
	sb.WriteString("# Syntax:\n")
	sb.WriteString("# - ',' or newline separates windows\n")
	sb.WriteString("# - ':' separates window name from panes\n")
	sb.WriteString("# - '|' separates panes (vertical split - side by side)\n")
	sb.WriteString("# - Empty command = shell prompt\n")
	sb.WriteString("# If not set, creates default layout: dev + agent windows\n")
	sb.WriteString(fmt.Sprintf("tmux_layout=%s\n", tmuxLayout))

	return sb.String()
}
