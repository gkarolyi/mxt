// Package worktree manages git worktree creation and operations.
package worktree

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gkarolyi/mxt/internal/ui"
)

// Create creates a new git worktree with a new branch.
//
// Steps:
//  1. Print info message with worktree path
//  2. Create parent directory if needed
//  3. Run: git worktree add -b <branch> <path> <base-branch>
//  4. Print success message with branch and base branch
//
// The git output (Preparing worktree, HEAD is now at...) is automatically
// printed to stdout by the git command.
func Create(worktreePath, branchName, baseBranch string) error {
	ui.Info(fmt.Sprintf("Creating worktree at %s", worktreePath))

	// Create parent directory if it doesn't exist
	parentDir := filepath.Dir(worktreePath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Create worktree: git worktree add -b <branch> <path> <base-branch>
	cmd := exec.Command("git", "worktree", "add", "-b", branchName, worktreePath, baseBranch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git worktree add failed: %w", err)
	}

	// Success message with colored branch names
	branchColored := ui.CyanText(branchName)
	baseColored := ui.DimText("from " + baseBranch)
	ui.Success(fmt.Sprintf("Worktree created (branch %s %s)", branchColored, baseColored))

	return nil
}

// CopyFiles copies files from source directory to worktree directory.
// The copyFiles parameter is a comma-separated list of file patterns (supports globs).
//
// For each pattern:
//  1. Trim whitespace
//  2. Expand glob relative to source directory
//  3. If no matches, print warning and continue
//  4. For each matched file, copy preserving attributes
//  5. Print success message for each copied file
//
// Returns error only if critical failure occurs. Missing files generate warnings but don't fail.
func CopyFiles(sourceDir, destDir, copyFiles string) error {
	ui.Info("Copying config files...")

	// Split by comma and process each pattern
	patterns := strings.Split(copyFiles, ",")

	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}

		// Expand glob pattern relative to source directory
		fullPattern := filepath.Join(sourceDir, pattern)
		matches, err := filepath.Glob(fullPattern)

		if err != nil {
			ui.Warn(fmt.Sprintf("  Invalid pattern: %s", ui.DimText(pattern)))
			continue
		}

		if len(matches) == 0 {
			ui.Warn(fmt.Sprintf("  Not found: %s", ui.DimText(pattern)))
			continue
		}

		// Copy each matched file
		for _, srcPath := range matches {
			// Calculate relative path from source directory
			relPath, err := filepath.Rel(sourceDir, srcPath)
			if err != nil {
				ui.Warn(fmt.Sprintf("  Failed to calculate relative path for %s", ui.DimText(srcPath)))
				continue
			}

			// Calculate destination path
			dstPath := filepath.Join(destDir, relPath)

			// Create destination parent directory if needed
			dstParent := filepath.Dir(dstPath)
			if err := os.MkdirAll(dstParent, 0755); err != nil {
				ui.Warn(fmt.Sprintf("  Failed to create directory for %s", ui.DimText(relPath)))
				continue
			}

			// Copy file
			if err := copyFile(srcPath, dstPath); err != nil {
				ui.Warn(fmt.Sprintf("  Failed to copy %s: %v", ui.DimText(relPath), err))
				continue
			}

			ui.Success(fmt.Sprintf("  Copied %s", ui.DimText(relPath)))
		}
	}

	return nil
}

// RunPreSessionCommand runs a command in the worktree directory.
//
// Steps:
//  1. Print info message
//  2. Print command (indented, dimmed)
//  3. Change to worktree directory
//  4. Execute command via shell
//  5. If success, print success message
//  6. If failure, return error with exit code
//
// Returns error if command fails. The caller should handle the error by
// prompting the user for confirmation.
func RunPreSessionCommand(worktreePath, command string) error {
	ui.Info("Running pre-session command...")

	// Print the command being run (indented and dimmed)
	fmt.Printf("  %s\n", ui.DimText(command))

	// Execute command in worktree directory using shell
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = worktreePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Extract exit code if possible
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode := exitErr.ExitCode()
			return fmt.Errorf("Pre-session command failed (exit code: %d)", exitCode)
		}
		return fmt.Errorf("Pre-session command failed: %w", err)
	}

	ui.Success("Pre-session command completed")
	return nil
}

// copyFile copies a file from src to dst, preserving file permissions.
// It creates the destination file with the same permissions as the source.
func copyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get source file info for permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create destination file with same permissions
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}
