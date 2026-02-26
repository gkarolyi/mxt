package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gkarolyi/mxt/internal/config"
	"github.com/gkarolyi/mxt/internal/ui"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/term"
)

const logo = `███╗   ███╗██╗  ██╗████████╗
████╗ ████║╚██╗██╔╝╚══██╔══╝
██╔████╔██║ ╚███╔╝    ██║
██║╚██╔╝██║ ██╔██╗    ██║
██║ ╚═╝ ██║██╔╝ ██╗   ██║
╚═╝     ╚═╝╚═╝  ╚═╝   ╚═╝
  Tmux Worktree Session Manager v1.1.0
`

// InitCommand implements the init command
func InitCommand(local bool, reinit bool, importLegacy bool) error {
	// Display logo
	fmt.Print(logo)
	fmt.Println()

	if local {
		return initProjectConfig(reinit, importLegacy)
	}
	return initGlobalConfig(reinit, importLegacy)
}

func shouldOverwriteConfig(reader *bufio.Reader, writer io.Writer, isTerminal bool, reinit bool) (bool, error) {
	if reinit {
		return true, nil
	}
	if isTerminal {
		fmt.Fprint(writer, "Overwrite? (y/N) ")
	}
	response, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, err
	}
	response = strings.TrimSpace(response)
	return response == "y" || response == "Y", nil
}

func importLegacyConfig(legacyPath string, targetPath string, reinit bool) (string, error) {
	if _, err := os.Stat(legacyPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("legacy config not found at %s", legacyPath)
		}
		return "", fmt.Errorf("failed to read legacy config at %s: %w", legacyPath, err)
	}
	if _, err := os.Stat(targetPath); err == nil {
		if !reinit {
			return "", fmt.Errorf("config already exists at %s (use --reinit to overwrite)", targetPath)
		}
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to check config at %s: %w", targetPath, err)
	}
	content, err := os.ReadFile(legacyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read legacy config at %s: %w", legacyPath, err)
	}
	parsed, err := config.ParseLegacyConfig(strings.NewReader(string(content)))
	if err != nil {
		return "", fmt.Errorf("failed to parse legacy config at %s: %w", legacyPath, err)
	}
	encoded, err := config.EncodeConfig(parsed)
	if err != nil {
		return "", fmt.Errorf("failed to encode TOML config: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}
	if err := os.WriteFile(targetPath, []byte(encoded), 0o644); err != nil {
		return "", fmt.Errorf("failed to write config file: %w", err)
	}
	return encoded, nil
}

func initGlobalConfig(reinit bool, importLegacy bool) error {
	configPath := config.GetGlobalConfigPath()
	if importLegacy {
		content, err := importLegacyConfig(config.GetLegacyGlobalConfigPath(), configPath, reinit)
		if err != nil {
			return err
		}
		ui.Success(fmt.Sprintf("Imported config written to %s", configPath))
		fmt.Println()
		fmt.Print(content)
		return nil
	}
	configDir := filepath.Dir(configPath)
	reader := bufio.NewReader(os.Stdin)

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		ui.Warn(fmt.Sprintf("Config already exists at %s", configPath))
		fmt.Println()
		content, _ := os.ReadFile(configPath)
		fmt.Println(string(content))
		fmt.Println()
		shouldOverwrite, err := shouldOverwriteConfig(reader, os.Stdout, term.IsTerminal(int(os.Stdin.Fd())), reinit)
		if err != nil {
			return err
		}
		if !shouldOverwrite {
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
	sandboxTool := defaults["sandbox_tool"]

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
	displaySandbox := sandboxTool
	if displaySandbox == "" {
		displaySandbox = "none"
	}
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Printf("Sandbox tool (firejail/docker, optional) [%s]: ", displaySandbox)
	}
	inputSandbox, _ := reader.ReadString('\n')
	inputSandbox = strings.TrimSpace(inputSandbox)
	if inputSandbox != "" {
		if inputSandbox == "none" {
			sandboxTool = ""
		} else {
			sandboxTool = inputSandbox
		}
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
	ui.Info("Enter a single line now, or edit the config file to use a TOML multi-line string.")
	ui.Info("Example (single line): dev:hx|lazygit,server:bin/server,agent:")
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Print("Tmux layout: ")
	}
	tmuxLayout, _ := reader.ReadString('\n')
	tmuxLayout = strings.TrimSpace(tmuxLayout)

	// Generate config file content
	content, err := generateGlobalConfigContent(worktreeDir, terminal, sandboxTool, copyFiles, preSessionCmd, tmuxLayout)
	if err != nil {
		return fmt.Errorf("failed to format config: %w", err)
	}

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

func initProjectConfig(reinit bool, importLegacy bool) error {
	gitRoot, err := config.FindGitRoot(".")
	if err != nil {
		return fmt.Errorf("Not inside a git repository. Run mxt from within your repo.")
	}

	configPath := config.GetProjectConfigPath(gitRoot)
	if importLegacy {
		content, err := importLegacyConfig(config.GetLegacyProjectConfigPath(gitRoot), configPath, reinit)
		if err != nil {
			return err
		}
		ui.Success(fmt.Sprintf("Imported project config written to %s", configPath))
		fmt.Println()
		fmt.Print(content)
		return nil
	}
	reader := bufio.NewReader(os.Stdin)

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		ui.Warn(fmt.Sprintf("Project config already exists at %s", configPath))
		fmt.Println()
		content, _ := os.ReadFile(configPath)
		fmt.Println(string(content))
		fmt.Println()
		shouldOverwrite, err := shouldOverwriteConfig(reader, os.Stdout, term.IsTerminal(int(os.Stdin.Fd())), reinit)
		if err != nil {
			return err
		}
		if !shouldOverwrite {
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
	// Sandbox tool
	fmt.Println()
	ui.Info("Optional: Sandbox tool command prefix for tmux sessions.")
	ui.Info("Example: firejail --private, docker run --rm -it ...")
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Print("Sandbox tool: ")
	}
	sandboxTool, _ := reader.ReadString('\n')
	sandboxTool = strings.TrimSpace(sandboxTool)

	// Tmux layout
	fmt.Println()
	ui.Info("Optional: Tmux layout - define windows and panes.")
	ui.Info("Enter a single line now, or edit the config file to use a TOML multi-line string.")
	ui.Info("Example (single line): dev:hx|lazygit,server:bin/server,agent:")
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Print("Tmux layout: ")
	}
	tmuxLayout, _ := reader.ReadString('\n')
	tmuxLayout = strings.TrimSpace(tmuxLayout)

	// Generate config file content
	content, err := generateProjectConfigContent(copyFiles, sandboxTool, preSessionCmd, tmuxLayout)
	if err != nil {
		return fmt.Errorf("failed to format config: %w", err)
	}

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

func generateGlobalConfigContent(worktreeDir, terminal, sandboxTool, copyFiles, preSessionCmd, tmuxLayout string) (string, error) {
	timestamp := time.Now().Format("Mon 02 Jan 2006 15:04:05 MST")

	worktreeValue, err := formatTomlValue(worktreeDir)
	if err != nil {
		return "", fmt.Errorf("failed to encode worktree_dir: %w", err)
	}
	terminalValue, err := formatTomlValue(terminal)
	if err != nil {
		return "", fmt.Errorf("failed to encode terminal: %w", err)
	}
	sandboxValue, err := formatTomlValue(sandboxTool)
	if err != nil {
		return "", fmt.Errorf("failed to encode sandbox_tool: %w", err)
	}
	copyFilesValue, err := formatTomlValue(copyFiles)
	if err != nil {
		return "", fmt.Errorf("failed to encode copy_files: %w", err)
	}
	preSessionValue, err := formatTomlValue(preSessionCmd)
	if err != nil {
		return "", fmt.Errorf("failed to encode pre_session_cmd: %w", err)
	}
	tmuxLayoutValue, err := formatTomlValue(tmuxLayout)
	if err != nil {
		return "", fmt.Errorf("failed to encode tmux_layout: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("# mxt configuration (TOML)\n")
	sb.WriteString(fmt.Sprintf("# Generated on %s\n\n", timestamp))
	sb.WriteString("# Base directory for worktrees\n")
	sb.WriteString(fmt.Sprintf("worktree_dir = %s\n\n", worktreeValue))

	sb.WriteString("# Terminal app: terminal | iterm2 | ghostty | current\n")
	sb.WriteString(fmt.Sprintf("terminal = %s\n\n", terminalValue))
	sb.WriteString("# Optional sandbox tool command prefix for tmux sessions\n")
	sb.WriteString("# Example: firejail --private, docker run --rm -it ...\n")
	sb.WriteString(fmt.Sprintf("sandbox_tool = %s\n\n", sandboxValue))

	sb.WriteString("# Files to copy from repo root into new worktrees (comma-separated, relative to repo root)\n")
	sb.WriteString("# Supports glob patterns and directories\n")
	sb.WriteString(fmt.Sprintf("copy_files = %s\n\n", copyFilesValue))

	sb.WriteString("# Command to run after worktree setup, before tmux session (optional)\n")
	sb.WriteString("# Runs in worktree directory. Use for setup tasks like: bundle install, npm install\n")
	sb.WriteString(fmt.Sprintf("pre_session_cmd = %s\n\n", preSessionValue))

	sb.WriteString("# Tmux layout - define windows and panes (optional)\n")
	sb.WriteString("# Multi-line format (more readable):\n")
	sb.WriteString("# tmux_layout = \"\"\"\n")
	sb.WriteString("#   dev:hx|lazygit\n")
	sb.WriteString("#   server:bin/server\n")
	sb.WriteString("#   agent:\n")
	sb.WriteString("# \"\"\"\n")
	sb.WriteString("# Or single line: tmux_layout = \"dev:hx|lazygit,server:bin/server,agent:\"\n")
	sb.WriteString("#\n")
	sb.WriteString("# Syntax:\n")
	sb.WriteString("# - ',' or newline separates windows\n")
	sb.WriteString("# - ':' separates window name from panes\n")
	sb.WriteString("# - '|' separates panes (vertical split - side by side)\n")
	sb.WriteString("# - Empty command = shell prompt\n")
	sb.WriteString("# If not set, creates default layout: dev + agent windows\n")
	sb.WriteString(fmt.Sprintf("tmux_layout = %s\n", tmuxLayoutValue))

	return sb.String(), nil
}

func generateProjectConfigContent(copyFiles, sandboxTool, preSessionCmd, tmuxLayout string) (string, error) {
	timestamp := time.Now().Format("Mon 02 Jan 2006 15:04:05 MST")

	copyFilesValue, err := formatTomlValue(copyFiles)
	if err != nil {
		return "", fmt.Errorf("failed to encode copy_files: %w", err)
	}
	sandboxValue, err := formatTomlValue(sandboxTool)
	if err != nil {
		return "", fmt.Errorf("failed to encode sandbox_tool: %w", err)
	}
	preSessionValue, err := formatTomlValue(preSessionCmd)
	if err != nil {
		return "", fmt.Errorf("failed to encode pre_session_cmd: %w", err)
	}
	tmuxLayoutValue, err := formatTomlValue(tmuxLayout)
	if err != nil {
		return "", fmt.Errorf("failed to encode tmux_layout: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("# mxt project config (TOML)\n")
	sb.WriteString(fmt.Sprintf("# Generated on %s\n\n", timestamp))
	sb.WriteString("# Files to copy from repo root into new worktrees (comma-separated, relative to repo root)\n")
	sb.WriteString("# Supports glob patterns and directories\n")
	sb.WriteString(fmt.Sprintf("copy_files = %s\n\n", copyFilesValue))
	sb.WriteString("# Optional sandbox tool command prefix for tmux sessions\n")
	sb.WriteString("# Example: firejail --private, docker run --rm -it ...\n")
	sb.WriteString(fmt.Sprintf("sandbox_tool = %s\n\n", sandboxValue))

	sb.WriteString("# Command to run after worktree setup, before tmux session (optional)\n")
	sb.WriteString("# Runs in worktree directory. Use for setup tasks like: bundle install, npm install\n")
	sb.WriteString(fmt.Sprintf("pre_session_cmd = %s\n\n", preSessionValue))

	sb.WriteString("# Tmux layout - define windows and panes (optional)\n")
	sb.WriteString("# Multi-line format (more readable):\n")
	sb.WriteString("# tmux_layout = \"\"\"\n")
	sb.WriteString("#   dev:hx|lazygit\n")
	sb.WriteString("#   server:bin/server\n")
	sb.WriteString("#   agent:\n")
	sb.WriteString("# \"\"\"\n")
	sb.WriteString("# Or single line: tmux_layout = \"dev:hx|lazygit,server:bin/server,agent:\"\n")
	sb.WriteString("#\n")
	sb.WriteString("# Syntax:\n")
	sb.WriteString("# - ',' or newline separates windows\n")
	sb.WriteString("# - ':' separates window name from panes\n")
	sb.WriteString("# - '|' separates panes (vertical split - side by side)\n")
	sb.WriteString("# - Empty command = shell prompt\n")
	sb.WriteString("# If not set, creates default layout: dev + agent windows\n")
	sb.WriteString(fmt.Sprintf("tmux_layout = %s\n", tmuxLayoutValue))

	return sb.String(), nil
}

func formatTomlValue(value string) (string, error) {
	data, err := toml.Marshal(map[string]string{"value": value})
	if err != nil {
		return "", err
	}
	line := strings.TrimSpace(string(data))
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("unexpected TOML encoding for value: %q", line)
	}
	return strings.TrimSpace(parts[1]), nil
}
