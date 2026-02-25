package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gkarolyi/mxt/internal/config"
	"github.com/gkarolyi/mxt/internal/git"
	"github.com/gkarolyi/mxt/internal/tmux"
	"github.com/gkarolyi/mxt/internal/ui"
)

// DeleteCommand deletes a worktree, kills its tmux session, and removes the branch.
func DeleteCommand(branch string, force bool) error {
	if !git.IsInsideWorkTree() {
		return fmt.Errorf("Not inside a git repository. Run muxtree from within your repo.")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	repoName, err := git.GetRepoName()
	if err != nil {
		return fmt.Errorf("failed to get repository name: %w", err)
	}

	worktreePath := git.CalculateWorktreePath(cfg.WorktreeDir, repoName, branch)
	if _, err := os.Stat(worktreePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Worktree not found: %s", worktreePath)
		}
		return fmt.Errorf("failed to access worktree: %w", err)
	}

	insertions, deletions := calculateChangeStats(worktreePath)

	fmt.Println()
	fmt.Printf("  Branch:    %s\n", ui.BoldText(branch))
	fmt.Printf("  Path:      %s\n", ui.DimText(worktreePath))
	fmt.Printf("  Changes:   %s %s\n", ui.GreenText(fmt.Sprintf("+%d", insertions)), ui.RedText(fmt.Sprintf("-%d", deletions)))
	fmt.Println()

	if !force {
		ui.Warn("This will remove the worktree and delete the local branch.")
		if !promptDeleteConfirm() {
			ui.Info("Cancelled.")
			return nil
		}
	}

	sessionName := git.GenerateSessionName(repoName, branch)
	if tmux.HasSession(sessionName) {
		if err := tmux.KillSession(sessionName); err != nil {
			return fmt.Errorf("failed to kill session %s: %w", sessionName, err)
		}
		ui.Success(fmt.Sprintf("Killed session %s", ui.BoldText(sessionName)))
	}

	ui.Info("Removing worktree...")
	if err := removeWorktree(worktreePath); err != nil {
		return err
	}
	ui.Success("Worktree removed")

	ui.Info(fmt.Sprintf("Deleting branch %s...", ui.CyanText(branch)))
	if err := deleteBranch(branch); err != nil {
		ui.Warn("Branch may have already been deleted")
	} else {
		ui.Success("Branch deleted")
	}

	cleanupRepoDir(cfg.WorktreeDir, repoName)

	fmt.Println()
	ui.Success("Done.")

	return nil
}

func promptDeleteConfirm() bool {
	fmt.Print("Are you sure? (y/N) ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false
	}
	response = strings.TrimSpace(response)
	return response == "y" || response == "Y"
}

func removeWorktree(worktreePath string) error {
	cmd := exec.Command("git", "worktree", "remove", worktreePath, "--force")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		ui.Warn("git worktree remove failed, cleaning up manually...")
		if err := os.RemoveAll(worktreePath); err != nil {
			return fmt.Errorf("failed to remove worktree: %w", err)
		}
		prune := exec.Command("git", "worktree", "prune")
		prune.Stdout = io.Discard
		prune.Stderr = io.Discard
		_ = prune.Run()
	}
	return nil
}

func deleteBranch(branch string) error {
	cmd := exec.Command("git", "branch", "-D", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.Discard
	return cmd.Run()
}

func cleanupRepoDir(worktreeDir, repoName string) {
	repoDir := filepath.Join(worktreeDir, repoName)
	entries, err := os.ReadDir(repoDir)
	if err != nil {
		return
	}
	if len(entries) == 0 {
		_ = os.Remove(repoDir)
	}
}
