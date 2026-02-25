package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gkarolyi/mxt/internal/config"
	"github.com/gkarolyi/mxt/internal/git"
	"github.com/gkarolyi/mxt/internal/tmux"
	"github.com/gkarolyi/mxt/internal/ui"
	"github.com/gkarolyi/mxt/internal/worktree"
)

// NewCommand creates a new git worktree with a new branch and launches tmux session.
//
// Phase 4: Implements worktree creation, file copying, and pre-session command execution.
// Phase 5: Will add tmux session creation and terminal opening.
func NewCommand(branchName string, fromBranch string, runCmd string, bg bool) error {
	// Step 1: Prerequisite Checks
	if !git.IsInsideWorkTree() {
		return fmt.Errorf("Not inside a git repository. Run mxt from within your repo")
	}

	// Step 2: Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Step 3: Validate --run command
	if runCmd != "" && runCmd != "claude" && runCmd != "codex" {
		return fmt.Errorf("Invalid --run command: '%s'. Use 'claude' or 'codex'", runCmd)
	}

	// Step 4: Determine base branch
	baseBranch := fromBranch
	if baseBranch == "" {
		baseBranch = git.GetMainBranch()
	}

	// Step 5: Validate base branch exists
	if err := validateBranchExists(baseBranch); err != nil {
		return fmt.Errorf("Base branch '%s' does not exist", baseBranch)
	}

	// Step 6: Check if new branch already exists
	if err := validateBranchExists(branchName); err == nil {
		return fmt.Errorf("Branch '%s' already exists. Use a different name, or delete it first", branchName)
	}

	// Step 7: Determine worktree path
	repoName, err := git.GetRepoName()
	if err != nil {
		return fmt.Errorf("failed to get repository name: %w", err)
	}

	worktreePath := git.CalculateWorktreePath(cfg.WorktreeDir, repoName, branchName)

	// Step 8: Check if worktree path already exists
	if _, err := os.Stat(worktreePath); err == nil {
		return fmt.Errorf("Worktree already exists at %s", worktreePath)
	}

	// === Execution Phase ===

	// Step 9: Create worktree
	if err := worktree.Create(worktreePath, branchName, baseBranch); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}

	// Step 10: Copy config files
	if cfg.CopyFiles != "" {
		repoRoot, err := git.GetRepoRoot()
		if err != nil {
			return fmt.Errorf("failed to get repo root: %w", err)
		}

		if err := worktree.CopyFiles(repoRoot, worktreePath, cfg.CopyFiles); err != nil {
			// Non-fatal: log warning but continue
			ui.Warn(fmt.Sprintf("Some files could not be copied: %v", err))
		}
	}

	// Step 11: Run pre-session command
	if cfg.PreSessionCmd != "" {
		if err := worktree.RunPreSessionCommand(worktreePath, cfg.PreSessionCmd); err != nil {
			// Prompt user for confirmation
			ui.Warn(fmt.Sprintf("Pre-session command failed: %v", err))
			if !promptContinue() {
				return fmt.Errorf("Aborted due to pre-session command failure")
			}
		}
	}

	// Step 12: Create tmux session
	ui.Info("Creating tmux session...")
	sessionName := git.GenerateSessionName(repoName, branchName)

	// Prepare session configuration
	sessionConfig := &tmux.SessionConfig{
		SessionName:  sessionName,
		WorktreePath: worktreePath,
		RunCommand:   runCmd,
		CustomLayout: cfg.TmuxLayout,
	}

	// Create session (custom or default layout)
	if cfg.TmuxLayout != "" {
		// Use custom layout
		if err := tmux.CreateCustomLayout(sessionConfig); err != nil {
			return fmt.Errorf("failed to create tmux session: %w", err)
		}
	} else {
		// Use default layout
		if err := tmux.CreateDefaultLayout(sessionConfig); err != nil {
			return fmt.Errorf("failed to create tmux session: %w", err)
		}
	}

	// Format window list for success message
	windowList := strings.Join(sessionConfig.WindowNames, ", ")
	ui.Success(fmt.Sprintf("  Created session %s (windows: %s)", ui.BoldText(sessionName), windowList))

	// Step 13: Open terminal (Phase 6)
	// TODO: Phase 6 - Implement terminal opening unless --bg
	_ = bg // Will be used in Phase 6

	// Step 14: Success message
	fmt.Println()
	ui.Success(fmt.Sprintf("Ready! Worktree: %s", ui.CyanText(worktreePath)))

	return nil
}

// validateBranchExists checks if a git branch exists (local or remote).
// Returns nil if branch exists, error if it doesn't.
func validateBranchExists(branch string) error {
	// Check local branch
	if err := runGitCommand("show-ref", "--verify", fmt.Sprintf("refs/heads/%s", branch)); err == nil {
		return nil
	}

	// Check remote branch
	if err := runGitCommand("show-ref", "--verify", fmt.Sprintf("refs/remotes/origin/%s", branch)); err == nil {
		return nil
	}

	return fmt.Errorf("branch not found")
}

// runGitCommand executes a git command and returns an error if it fails.
func runGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	return cmd.Run()
}

// promptContinue prompts the user to continue after a failure.
// Returns true if user enters 'y' or 'Y', false otherwise.
func promptContinue() bool {
	fmt.Print("Continue anyway? (y/N) ")
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y"
}
