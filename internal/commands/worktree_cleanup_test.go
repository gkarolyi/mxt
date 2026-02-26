package commands

import (
	"errors"
	"testing"
)

func TestCleanupWorktreeRunsAllSteps(t *testing.T) {
	originalRemove := removeWorktree
	originalDelete := deleteBranch
	t.Cleanup(func() {
		removeWorktree = originalRemove
		deleteBranch = originalDelete
	})

	removeCalls := 0
	deleteCalls := 0
	removeWorktree = func(path string) error {
		removeCalls++
		if path != "path" {
			t.Fatalf("expected path to be 'path', got %q", path)
		}
		return nil
	}
	deleteBranch = func(branch string) error {
		deleteCalls++
		if branch != "branch" {
			t.Fatalf("expected branch to be 'branch', got %q", branch)
		}
		return nil
	}

	if err := cleanupWorktree("path", "branch"); err != nil {
		t.Fatalf("cleanupWorktree() unexpected error: %v", err)
	}
	if removeCalls != 1 {
		t.Fatalf("expected removeWorktree to be called once, got %d", removeCalls)
	}
	if deleteCalls != 1 {
		t.Fatalf("expected deleteBranch to be called once, got %d", deleteCalls)
	}
}

func TestCleanupWorktreeJoinsErrors(t *testing.T) {
	originalRemove := removeWorktree
	originalDelete := deleteBranch
	t.Cleanup(func() {
		removeWorktree = originalRemove
		deleteBranch = originalDelete
	})

	removeCalls := 0
	deleteCalls := 0
	removeErr := errors.New("remove failed")
	deleteErr := errors.New("delete failed")
	removeWorktree = func(path string) error {
		removeCalls++
		return removeErr
	}
	deleteBranch = func(branch string) error {
		deleteCalls++
		return deleteErr
	}

	err := cleanupWorktree("path", "branch")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, removeErr) {
		t.Fatalf("expected remove error to be joined, got %v", err)
	}
	if !errors.Is(err, deleteErr) {
		t.Fatalf("expected delete error to be joined, got %v", err)
	}
	if removeCalls != 1 {
		t.Fatalf("expected removeWorktree to be called once, got %d", removeCalls)
	}
	if deleteCalls != 1 {
		t.Fatalf("expected deleteBranch to be called once, got %d", deleteCalls)
	}
}
