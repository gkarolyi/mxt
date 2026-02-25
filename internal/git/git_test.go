package git

import (
	"path/filepath"
	"testing"
)

// TestSanitizeBranchName tests branch name sanitization for filesystem and tmux compatibility
func TestSanitizeBranchName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic cases - already safe
		{
			name:     "alphanumeric only",
			input:    "feature123",
			expected: "feature123",
		},
		{
			name:     "with underscore",
			input:    "feature_auth",
			expected: "feature_auth",
		},
		{
			name:     "with dash",
			input:    "feature-auth",
			expected: "feature-auth",
		},
		{
			name:     "with dot",
			input:    "v1.0.0",
			expected: "v1.0.0",
		},

		// Common special character replacements
		{
			name:     "slash to dash",
			input:    "feature/auth",
			expected: "feature-auth",
		},
		{
			name:     "multiple slashes",
			input:    "feature/auth/api",
			expected: "feature-auth-api",
		},
		{
			name:     "hash to dash",
			input:    "bug-fix-#123",
			expected: "bug-fix--123",
		},
		{
			name:     "at sign to dash",
			input:    "user@domain",
			expected: "user-domain",
		},

		// Multiple special characters
		{
			name:     "mixed special characters",
			input:    "feature/auth#123@test",
			expected: "feature-auth-123-test",
		},
		{
			name:     "spaces to dashes",
			input:    "feature auth api",
			expected: "feature-auth-api",
		},

		// Leading dash removal
		{
			name:     "leading dash removed",
			input:    "-feature",
			expected: "feature",
		},
		{
			name:     "leading special char creates dash then removed",
			input:    "/feature",
			expected: "feature",
		},
		{
			name:     "multiple leading special chars",
			input:    "//feature",
			expected: "feature",
		},

		// Edge cases
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only special characters",
			input:    "/@#$",
			expected: "",
		},
		{
			name:     "only dashes",
			input:    "---",
			expected: "",
		},
		{
			name:     "trailing dash preserved",
			input:    "feature-",
			expected: "feature-",
		},
		{
			name:     "consecutive special chars",
			input:    "feature///auth",
			expected: "feature---auth",
		},

		// Real-world branch names
		{
			name:     "jira style",
			input:    "PROJ-123/feature",
			expected: "PROJ-123-feature",
		},
		{
			name:     "github pr style",
			input:    "dependabot/npm_and_yarn/lodash-4.17.21",
			expected: "dependabot-npm_and_yarn-lodash-4.17.21",
		},
		{
			name:     "semantic version tag",
			input:    "v1.2.3-beta.1",
			expected: "v1.2.3-beta.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeBranchName(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeBranchName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestGenerateSessionName tests session name generation from repo name and branch
func TestGenerateSessionName(t *testing.T) {
	tests := []struct {
		name       string
		repoName   string
		branchName string
		expected   string
	}{
		// Basic cases
		{
			name:       "simple names",
			repoName:   "myapp",
			branchName: "main",
			expected:   "myapp_main",
		},
		{
			name:       "repo and branch with dashes",
			repoName:   "my-app",
			branchName: "feature-auth",
			expected:   "my-app_feature-auth",
		},

		// Branch name needs sanitization
		{
			name:       "branch with slash",
			repoName:   "myapp",
			branchName: "feature/auth",
			expected:   "myapp_feature-auth",
		},
		{
			name:       "branch with special chars",
			repoName:   "myapp",
			branchName: "bug-fix-#123",
			expected:   "myapp_bug-fix--123",
		},
		{
			name:       "branch with at sign",
			repoName:   "myapp",
			branchName: "user@domain",
			expected:   "myapp_user-domain",
		},

		// Complex real-world examples
		{
			name:       "jira style branch",
			repoName:   "project-api",
			branchName: "PROJ-123/feature",
			expected:   "project-api_PROJ-123-feature",
		},
		{
			name:       "github dependabot",
			repoName:   "my-app",
			branchName: "dependabot/npm_and_yarn/lodash-4.17.21",
			expected:   "my-app_dependabot-npm_and_yarn-lodash-4.17.21",
		},

		// Edge cases
		{
			name:       "empty branch",
			repoName:   "myapp",
			branchName: "",
			expected:   "myapp_",
		},
		{
			name:       "empty repo",
			repoName:   "",
			branchName: "main",
			expected:   "_main",
		},
		{
			name:       "both empty",
			repoName:   "",
			branchName: "",
			expected:   "_",
		},

		// Leading special char in branch (gets sanitized)
		{
			name:       "leading slash in branch",
			repoName:   "myapp",
			branchName: "/feature",
			expected:   "myapp_feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSessionName(tt.repoName, tt.branchName)
			if result != tt.expected {
				t.Errorf("GenerateSessionName(%q, %q) = %q, want %q", tt.repoName, tt.branchName, result, tt.expected)
			}
		})
	}
}

// TestCalculateWorktreePath tests worktree path calculation
func TestCalculateWorktreePath(t *testing.T) {
	tests := []struct {
		name         string
		worktreeDir  string
		repoName     string
		branchName   string
		expected     string
	}{
		// Basic cases
		{
			name:         "simple path",
			worktreeDir:  "/home/user/worktrees",
			repoName:     "myapp",
			branchName:   "main",
			expected:     "/home/user/worktrees/myapp/main",
		},
		{
			name:         "with dashes",
			worktreeDir:  "/home/user/worktrees",
			repoName:     "my-app",
			branchName:   "feature-auth",
			expected:     "/home/user/worktrees/my-app/feature-auth",
		},

		// Branch name needs sanitization
		{
			name:         "branch with slash",
			worktreeDir:  "/home/user/worktrees",
			repoName:     "myapp",
			branchName:   "feature/auth",
			expected:     "/home/user/worktrees/myapp/feature-auth",
		},
		{
			name:         "branch with special chars",
			worktreeDir:  "/home/user/worktrees",
			repoName:     "myapp",
			branchName:   "bug-fix-#123",
			expected:     "/home/user/worktrees/myapp/bug-fix--123",
		},
		{
			name:         "jira style branch",
			worktreeDir:  "/Users/dev/wt",
			repoName:     "project-api",
			branchName:   "PROJ-123/feature",
			expected:     "/Users/dev/wt/project-api/PROJ-123-feature",
		},

		// Edge cases
		{
			name:         "trailing slash in worktree dir",
			worktreeDir:  "/home/user/worktrees/",
			repoName:     "myapp",
			branchName:   "main",
			expected:     "/home/user/worktrees/myapp/main",
		},
		{
			name:         "empty branch name",
			worktreeDir:  "/home/user/worktrees",
			repoName:     "myapp",
			branchName:   "",
			expected:     "/home/user/worktrees/myapp",
		},
		{
			name:         "leading slash in branch creates dash then removed",
			worktreeDir:  "/home/user/worktrees",
			repoName:     "myapp",
			branchName:   "/feature",
			expected:     "/home/user/worktrees/myapp/feature",
		},

		// Tilde expansion should be handled by caller
		{
			name:         "unexpanded tilde",
			worktreeDir:  "~/worktrees",
			repoName:     "myapp",
			branchName:   "main",
			expected:     "~/worktrees/myapp/main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateWorktreePath(tt.worktreeDir, tt.repoName, tt.branchName)
			if result != tt.expected {
				t.Errorf("CalculateWorktreePath(%q, %q, %q) = %q, want %q",
					tt.worktreeDir, tt.repoName, tt.branchName, result, tt.expected)
			}
		})
	}
}

// Integration tests for git helper functions
// These tests require running inside a git repository

// TestGetRepoRoot tests getting the repository root path
func TestGetRepoRoot(t *testing.T) {
	// This test requires being run inside a git repository
	if !IsInsideWorkTree() {
		t.Skip("Not inside a git repository, skipping integration test")
	}

	root, err := GetRepoRoot()
	if err != nil {
		t.Fatalf("GetRepoRoot() error = %v", err)
	}

	if root == "" {
		t.Error("GetRepoRoot() returned empty string")
	}

	// The path should be absolute
	if !filepath.IsAbs(root) {
		t.Errorf("GetRepoRoot() = %q, want absolute path", root)
	}
}

// TestGetRepoName tests getting the repository name
func TestGetRepoName(t *testing.T) {
	// This test requires being run inside a git repository
	if !IsInsideWorkTree() {
		t.Skip("Not inside a git repository, skipping integration test")
	}

	name, err := GetRepoName()
	if err != nil {
		t.Fatalf("GetRepoName() error = %v", err)
	}

	if name == "" {
		t.Error("GetRepoName() returned empty string")
	}

	// Should be just the directory name, not a path
	if filepath.Dir(name) != "." {
		t.Errorf("GetRepoName() = %q, want directory name only (no path separators)", name)
	}
}

// TestGetMainBranch tests detecting the main branch
func TestGetMainBranch(t *testing.T) {
	// This test requires being run inside a git repository
	if !IsInsideWorkTree() {
		t.Skip("Not inside a git repository, skipping integration test")
	}

	branch := GetMainBranch()
	if branch == "" {
		t.Error("GetMainBranch() returned empty string")
	}

	// Should be a valid branch name (common ones)
	validBranches := map[string]bool{
		"main":    true,
		"master":  true,
		"develop": true,
		"trunk":   true,
	}

	if !validBranches[branch] {
		// Not a fatal error, just log it - could be a custom default branch
		t.Logf("GetMainBranch() = %q (uncommon but possibly valid)", branch)
	}
}

// TestIsInsideWorkTree tests checking if we're inside a git work tree
func TestIsInsideWorkTree(t *testing.T) {
	// This test should always pass when run via go test in the project
	result := IsInsideWorkTree()
	if !result {
		t.Error("IsInsideWorkTree() = false, want true (tests should run inside git repo)")
	}
}
