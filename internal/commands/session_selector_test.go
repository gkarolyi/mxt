package commands

import (
	"strings"
	"testing"
)

func TestSelectionCandidatesFiltering(t *testing.T) {
	worktrees := []WorktreeInfo{
		{BranchName: "alpha", SessionActive: true},
		{BranchName: "beta", SessionActive: false},
	}

	candidates, message, err := selectionCandidates("open", worktrees)
	if err != nil {
		t.Fatalf("selectionCandidates(open) error = %v", err)
	}
	if message != "" {
		t.Fatalf("selectionCandidates(open) message = %q, want empty", message)
	}
	if len(candidates) != 2 {
		t.Fatalf("selectionCandidates(open) len = %d, want 2", len(candidates))
	}

	candidates, message, err = selectionCandidates("close", worktrees)
	if err != nil {
		t.Fatalf("selectionCandidates(close) error = %v", err)
	}
	if message != "" {
		t.Fatalf("selectionCandidates(close) message = %q, want empty", message)
	}
	if len(candidates) != 1 || candidates[0].BranchName != "alpha" {
		t.Fatalf("selectionCandidates(close) candidates = %#v, want only alpha", candidates)
	}
}

func TestSelectionCandidatesEmptyMessage(t *testing.T) {
	candidates, message, err := selectionCandidates("open", nil)
	if err != nil {
		t.Fatalf("selectionCandidates(open) error = %v", err)
	}
	if len(candidates) != 0 || message != noWorktreesMessage {
		t.Fatalf("selectionCandidates(open) message = %q, want %q", message, noWorktreesMessage)
	}

	worktrees := []WorktreeInfo{{BranchName: "beta", SessionActive: false}}
	candidates, message, err = selectionCandidates("close", worktrees)
	if err != nil {
		t.Fatalf("selectionCandidates(close) error = %v", err)
	}
	if len(candidates) != 0 || message != noActiveSessionsMessage {
		t.Fatalf("selectionCandidates(close) message = %q, want %q", message, noActiveSessionsMessage)
	}
}

func TestResolveSessionBranchNonTTY(t *testing.T) {
	_, err := resolveSessionBranch("open", "", false, nil)
	if err == nil {
		t.Fatal("resolveSessionBranch() expected error when non-tty")
	}
	if !strings.Contains(err.Error(), "Usage:") {
		t.Fatalf("resolveSessionBranch() error = %q, want usage text", err.Error())
	}
}

func TestResolveSessionBranchCancel(t *testing.T) {
	originalRunner := runFzfSelector
	runFzfSelector = func(items []string) (string, int, error) {
		return "", 130, nil
	}
	defer func() {
		runFzfSelector = originalRunner
	}()

	branch, err := resolveSessionBranch("open", "", true, []WorktreeInfo{{BranchName: "alpha", SessionActive: true}})
	if err != nil {
		t.Fatalf("resolveSessionBranch() error = %v", err)
	}
	if branch != "" {
		t.Fatalf("resolveSessionBranch() = %q, want empty", branch)
	}
}

func TestResolveSessionBranchSelectsBranch(t *testing.T) {
	originalRunner := runFzfSelector
	runFzfSelector = func(items []string) (string, int, error) {
		return "feature\t(active)", 0, nil
	}
	defer func() {
		runFzfSelector = originalRunner
	}()

	branch, err := resolveSessionBranch("open", "", true, []WorktreeInfo{{BranchName: "feature", SessionActive: true}})
	if err != nil {
		t.Fatalf("resolveSessionBranch() error = %v", err)
	}
	if branch != "feature" {
		t.Fatalf("resolveSessionBranch() = %q, want feature", branch)
	}
}
