// Package errors provides custom error types and error handling utilities.
package errors

import (
	"fmt"

	"github.com/gkarolyi/mxt/internal/ui"
)

// Common error types for mxt operations

// ErrNotInGitRepo indicates the current directory is not inside a git repository
type ErrNotInGitRepo struct{}

func (e ErrNotInGitRepo) Error() string {
	return "Not inside a git repository. Run mxt from within your repo."
}

// ErrBranchExists indicates the branch already exists
type ErrBranchExists struct {
	Branch string
}

func (e ErrBranchExists) Error() string {
	return fmt.Sprintf("Branch '%s' already exists. Use a different name, or delete it first.", e.Branch)
}

// ErrBranchNotFound indicates the branch does not exist
type ErrBranchNotFound struct {
	Branch string
}

func (e ErrBranchNotFound) Error() string {
	return fmt.Sprintf("Base branch '%s' does not exist.", e.Branch)
}

// ErrWorktreeExists indicates a worktree already exists at the path
type ErrWorktreeExists struct {
	Path string
}

func (e ErrWorktreeExists) Error() string {
	return fmt.Sprintf("Worktree already exists at %s", e.Path)
}

// ErrWorktreeNotFound indicates the worktree was not found
type ErrWorktreeNotFound struct {
	Path string
}

func (e ErrWorktreeNotFound) Error() string {
	return fmt.Sprintf("Worktree not found: %s", e.Path)
}

// ErrSessionExists indicates a tmux session already exists
type ErrSessionExists struct {
	Session string
}

func (e ErrSessionExists) Error() string {
	return fmt.Sprintf("Session %s already exists", e.Session)
}

// ErrSessionNotFound indicates a tmux session was not found
type ErrSessionNotFound struct {
	Session string
}

func (e ErrSessionNotFound) Error() string {
	return fmt.Sprintf("Session not found: %s", e.Session)
}

// ErrInvalidCommand indicates an invalid command was provided
type ErrInvalidCommand struct {
	Command string
	Valid   []string
}

func (e ErrInvalidCommand) Error() string {
	return fmt.Sprintf("Invalid command '%s'. Valid options: %v", e.Command, e.Valid)
}

// ErrConfigNotFound indicates no configuration was found
type ErrConfigNotFound struct{}

func (e ErrConfigNotFound) Error() string {
	return "No config found. Run mxt init to create one."
}

// ErrInvalidConfig indicates configuration is invalid
type ErrInvalidConfig struct {
	Key    string
	Reason string
}

func (e ErrInvalidConfig) Error() string {
	return fmt.Sprintf("Invalid config for '%s': %s", e.Key, e.Reason)
}

// Die prints an error message and exits with code 1
// This is a convenience function that wraps ui.Die
func Die(msg string) {
	ui.Die(msg)
}

// Dief prints a formatted error message and exits with code 1
func Dief(format string, args ...interface{}) {
	ui.Dief(format, args...)
}

// DieIf exits with an error message if the error is not nil
func DieIf(err error) {
	if err != nil {
		ui.Die(err.Error())
	}
}

// DieIfWithMsg exits with a custom message if the error is not nil
func DieIfWithMsg(err error, msg string) {
	if err != nil {
		ui.Die(msg)
	}
}

// WarnIf prints a warning if the error is not nil
func WarnIf(err error) {
	if err != nil {
		ui.Warn(err.Error())
	}
}

// WarnIfWithMsg prints a custom warning if the error is not nil
func WarnIfWithMsg(err error, msg string) {
	if err != nil {
		ui.Warn(msg)
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

// Wrapf wraps an error with formatted context
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	context := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %w", context, err)
}
