package commands

import (
	"errors"
	"fmt"

	"github.com/gkarolyi/mxt/internal/git"
	"github.com/gkarolyi/mxt/internal/worktree"
)

var (
	removeWorktree = worktree.Remove
	deleteBranch   = git.DeleteBranch
)

func cleanupWorktree(worktreePath, branchName string) error {
	var errs []error
	if err := removeWorktree(worktreePath); err != nil {
		errs = append(errs, fmt.Errorf("remove worktree: %w", err))
	}
	if err := deleteBranch(branchName); err != nil {
		errs = append(errs, fmt.Errorf("delete branch %s: %w", branchName, err))
	}
	return errors.Join(errs...)
}
