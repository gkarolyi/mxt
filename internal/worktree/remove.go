package worktree

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/gkarolyi/mxt/internal/ui"
)

// Remove deletes a git worktree at the provided path.
// It attempts git worktree removal first, then falls back to manual cleanup.
func Remove(worktreePath string) error {
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
