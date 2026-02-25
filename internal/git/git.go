// Package git provides git repository operations and helpers.
package git

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// SanitizeBranchName sanitizes a branch name for filesystem and tmux compatibility.
// It replaces any character that is NOT alphanumeric, underscore, dash, or dot with dash,
// then strips any leading dash.
//
// Examples:
//   - "feature/auth" → "feature-auth"
//   - "bug-fix-#123" → "bug-fix--123"
//   - "user@domain" → "user-domain"
func SanitizeBranchName(branch string) string {
	if branch == "" {
		return ""
	}

	// Replace any character that is NOT [a-zA-Z0-9._-] with dash
	re := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	sanitized := re.ReplaceAllString(branch, "-")

	// Strip leading dashes
	sanitized = strings.TrimLeft(sanitized, "-")

	return sanitized
}

// GenerateSessionName generates a tmux session name from repo name and branch name.
// Format: <repo-name>_<sanitized-branch>
//
// Examples:
//   - ("myapp", "feature/auth") → "myapp_feature-auth"
//   - ("project-api", "PROJ-123/feature") → "project-api_PROJ-123-feature"
func GenerateSessionName(repoName, branchName string) string {
	sanitized := SanitizeBranchName(branchName)
	return repoName + "_" + sanitized
}

// CalculateWorktreePath calculates the full path to a worktree directory.
// Format: <worktree-dir>/<repo-name>/<sanitized-branch>
//
// The branch name is sanitized for filesystem safety. Trailing slashes in
// worktreeDir are handled automatically by filepath.Join.
//
// Examples:
//   - ("/home/user/worktrees", "myapp", "feature/auth") → "/home/user/worktrees/myapp/feature-auth"
//   - ("/Users/dev/wt", "project-api", "PROJ-123/feature") → "/Users/dev/wt/project-api/PROJ-123-feature"
func CalculateWorktreePath(worktreeDir, repoName, branchName string) string {
	sanitized := SanitizeBranchName(branchName)
	return filepath.Join(worktreeDir, repoName, sanitized)
}

// GetRepoRoot returns the absolute path to the git repository root.
// Uses: git rev-parse --show-toplevel
func GetRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetRepoName returns the repository directory name.
// It's the basename of the repository root path.
func GetRepoName() (string, error) {
	root, err := GetRepoRoot()
	if err != nil {
		return "", err
	}
	return filepath.Base(root), nil
}

// GetMainBranch detects and returns the name of the main branch.
// Algorithm:
//  1. Try: git symbolic-ref refs/remotes/origin/HEAD → extract last component
//  2. If fails, check if "main" exists
//  3. If not, check if "master" exists
//  4. Fallback: return "main"
func GetMainBranch() string {
	// Try to get the default branch from origin/HEAD
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	output, err := cmd.Output()
	if err == nil {
		// Output format: refs/remotes/origin/main
		ref := strings.TrimSpace(string(output))
		parts := strings.Split(ref, "/")
		if len(parts) > 0 {
			branch := parts[len(parts)-1]
			if branch != "" {
				return branch
			}
		}
	}

	// Check if "main" branch exists
	cmd = exec.Command("git", "show-ref", "--verify", "refs/heads/main")
	if err := cmd.Run(); err == nil {
		return "main"
	}

	// Check if "master" branch exists
	cmd = exec.Command("git", "show-ref", "--verify", "refs/heads/master")
	if err := cmd.Run(); err == nil {
		return "master"
	}

	// Fallback to "main"
	return "main"
}

// IsInsideWorkTree checks if the current directory is inside a git work tree.
// Uses: git rev-parse --is-inside-work-tree
func IsInsideWorkTree() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}
