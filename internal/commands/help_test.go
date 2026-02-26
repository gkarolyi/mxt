package commands

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/gkarolyi/mxt/internal/ui"
)

func TestHelpCommandPrintsLogo(t *testing.T) {
	oldStdout := os.Stdout
	defer func() {
		os.Stdout = oldStdout
		ui.InitColors()
	}()

	r, w, _ := os.Pipe()
	os.Stdout = w
	ui.InitColors()

	HelpCommand("test-version")

	w.Close()
	out, _ := io.ReadAll(r)
	output := string(out)

	if !strings.Contains(output, "███╗   ███╗██╗  ██╗████████╗") {
		t.Fatalf("HelpCommand output missing ASCII logo: %q", output)
	}

	if !strings.Contains(output, "Tmux Worktree Session Manager vtest-version") {
		t.Fatalf("HelpCommand output missing version line: %q", output)
	}
}
