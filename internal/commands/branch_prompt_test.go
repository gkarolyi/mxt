package commands

import (
	"bytes"
	"strings"
	"testing"
)

func TestPromptBranchNameTrimsWhitespace(t *testing.T) {
	reader := strings.NewReader("  feature-auth \n")
	var output bytes.Buffer

	branch, err := PromptBranchName(reader, &output)
	if err != nil {
		t.Fatalf("PromptBranchName() error = %v", err)
	}
	if branch != "feature-auth" {
		t.Fatalf("PromptBranchName() = %q, want %q", branch, "feature-auth")
	}
	if output.String() != "Branch name: " {
		t.Fatalf("expected prompt output to be %q, got %q", "Branch name: ", output.String())
	}
}

func TestPromptBranchNameRejectsEmptyInput(t *testing.T) {
	reader := strings.NewReader("   \n")
	var output bytes.Buffer

	_, err := PromptBranchName(reader, &output)
	if err == nil {
		t.Fatal("expected error for empty input")
	}
	if err.Error() != "Branch name is required." {
		t.Fatalf("expected error message %q, got %q", "Branch name is required.", err.Error())
	}
}
