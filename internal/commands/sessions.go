package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/gkarolyi/mxt/internal/config"
	"github.com/gkarolyi/mxt/internal/git"
	"github.com/gkarolyi/mxt/internal/terminal"
	"github.com/gkarolyi/mxt/internal/tmux"
	"github.com/gkarolyi/mxt/internal/ui"
)

// SessionsCommand handles tmux session management for existing worktrees.
// Supports actions: open, close, relaunch, attach
func SessionsCommand(action string, branchName string, runCmd string, bg bool) error {
	// Normalize action (handle aliases)
	switch action {
	case "launch", "start":
		action = "open"
	case "kill", "stop":
		action = "close"
	case "restart":
		action = "relaunch"
	}

	// Dispatch to appropriate handler
	switch action {
	case "open":
		return sessionsOpen(branchName, runCmd, bg)
	case "close":
		return sessionsClose(branchName)
	case "relaunch":
		return sessionsRelaunch(branchName, runCmd, bg)
	case "attach":
		return sessionsAttach(branchName, runCmd) // runCmd used as window name for attach
	default:
		return fmt.Errorf("Unknown action: %s (use open, close, relaunch, or attach)", action)
	}
}

// sessionsOpen creates a tmux session for an existing worktree and opens terminal.
func sessionsOpen(branchName string, runCmd string, bg bool) error {
	// Step 1: Require git repository
	if !git.IsInsideWorkTree() {
		return fmt.Errorf("Not inside a git repository. Run mxt from within your repo.")
	}

	// Step 2: Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Step 3: Validate --run command
	if runCmd != "" && runCmd != "claude" && runCmd != "codex" {
		return fmt.Errorf("Invalid --run command: '%s'. Allowed: claude, codex", runCmd)
	}

	// Step 4: Determine repository name and worktree path
	repoName, err := git.GetRepoName()
	if err != nil {
		return fmt.Errorf("failed to get repository name: %w", err)
	}

	worktreePath := git.CalculateWorktreePath(cfg.WorktreeDir, repoName, branchName)

	// Step 5: Validate worktree exists
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return fmt.Errorf("Worktree not found: %s", worktreePath)
	}

	// Step 6: Determine session name
	sessionName := git.GenerateSessionName(repoName, branchName)

	// Step 7: Check if session already exists
	if tmux.HasSession(sessionName) {
		ui.Warn(fmt.Sprintf("Session %s already exists", sessionName))
		return nil
	}

	// Step 8: Create tmux session

	sessionConfig := &tmux.SessionConfig{
		SessionName:  sessionName,
		WorktreePath: worktreePath,
		RunCommand:   runCmd,
		CustomLayout: cfg.TmuxLayout,
	}

	// Create session (custom or default layout)
	if cfg.TmuxLayout != "" {
		if err := tmux.CreateCustomLayout(sessionConfig); err != nil {
			return fmt.Errorf("failed to create tmux session: %w", err)
		}
	} else {
		if err := tmux.CreateDefaultLayout(sessionConfig); err != nil {
			return fmt.Errorf("failed to create tmux session: %w", err)
		}
	}

	// Format window list for success message
	separator := ", "
	if cfg.TmuxLayout != "" {
		separator = " "
	}
	windowList := strings.Join(sessionConfig.WindowNames, separator)
	ui.Success(fmt.Sprintf("  Created session %s (windows: %s)", ui.BoldText(sessionName), windowList))

	// Step 9: Open terminal (unless --bg)
	if !bg {
		if err := terminal.Open(cfg.Terminal, sessionName); err != nil {
			ui.Warn(fmt.Sprintf("Failed to open terminal: %v", err))
			ui.Info(fmt.Sprintf("Run: tmux attach -t %s", sessionName))
		}
	}

	return nil
}

// sessionsClose kills a tmux session for a worktree.
func sessionsClose(branchName string) error {
	// Step 1: Require git repository
	if !git.IsInsideWorkTree() {
		return fmt.Errorf("Not inside a git repository. Run mxt from within your repo.")
	}

	// Step 2: Load configuration (needed for session name generation)
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	_ = cfg // Not actually needed for close, but keeping for consistency

	// Step 3: Determine session name
	repoName, err := git.GetRepoName()
	if err != nil {
		return fmt.Errorf("failed to get repository name: %w", err)
	}

	sessionName := git.GenerateSessionName(repoName, branchName)

	if !tmux.HasSession(sessionName) {
		return nil
	}
	// Step 4: Kill session if exists
	if err := tmux.KillSession(sessionName); err != nil {
		return fmt.Errorf("failed to kill session: %w", err)
	}

	ui.Success(fmt.Sprintf("Killed session %s", ui.BoldText(sessionName)))

	return nil
}

// sessionsRelaunch kills and recreates a tmux session.
func sessionsRelaunch(branchName string, runCmd string, bg bool) error {
	// Step 1: Close the session
	if err := sessionsClose(branchName); err != nil {
		return err
	}

	// Step 2: Open the session
	if err := sessionsOpen(branchName, runCmd, bg); err != nil {
		return err
	}

	return nil
}

// sessionsAttach attaches to an existing tmux session in the current terminal.
// The windowName parameter is optional and can be "dev" or "agent".
func sessionsAttach(branchName string, windowName string) error {
	// Step 1: Require git repository
	if !git.IsInsideWorkTree() {
		return fmt.Errorf("Not inside a git repository. Run mxt from within your repo.")
	}

	// Step 2: Load configuration (needed for session name generation)
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	_ = cfg // Not actually needed for attach

	// Step 3: Determine session name
	repoName, err := git.GetRepoName()
	if err != nil {
		return fmt.Errorf("failed to get repository name: %w", err)
	}

	sessionName := git.GenerateSessionName(repoName, branchName)

	// Step 4: Check if session exists
	if !tmux.HasSession(sessionName) {
		return fmt.Errorf("Session not found: %s", sessionName)
	}

	// Step 5: Attach to session (with optional window selection)
	if err := tmux.AttachToSession(sessionName, windowName); err != nil {
		return err
	}

	return nil
}
