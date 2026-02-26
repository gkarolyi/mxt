package commands

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gkarolyi/mxt/internal/ui"
	"golang.org/x/term"
)

const (
	noWorktreesMessage      = "No managed worktrees found. Use mxt new [branch] to create one."
	noActiveSessionsMessage = "No active tmux sessions found."
	selectionCancelledMsg   = "Selection cancelled."
)

type fzfRunner func(items []string) (string, int, error)

var runFzfSelector fzfRunner = runFzf

// SessionUsage returns the usage string for sessions commands.
func SessionUsage(action string) string {
	switch action {
	case "open", "launch", "start":
		return "Usage: mxt sessions open <branch> [--run claude|codex] [--bg] (omit branch to select interactively)"
	case "close", "kill", "stop":
		return "Usage: mxt sessions close <branch> (omit branch to select interactively)"
	case "relaunch", "restart":
		return "Usage: mxt sessions relaunch <branch> [--run claude|codex] [--bg] (omit branch to select interactively)"
	case "attach":
		return "Usage: mxt sessions attach <branch> [dev|agent] (omit branch to select interactively)"
	default:
		return "Usage: mxt sessions <open|close|relaunch|attach> <branch> [--run claude|codex] [--bg] (omit branch to select interactively)"
	}
}

func isInteractive() bool {
	return term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd()))
}

func resolveSessionBranch(action, branchName string, isTTY bool, worktrees []WorktreeInfo) (string, error) {
	if branchName != "" {
		return branchName, nil
	}
	if !isTTY {
		return "", errors.New(SessionUsage(action))
	}

	candidates, emptyMessage, err := selectionCandidates(action, worktrees)
	if err != nil {
		return "", err
	}
	if len(candidates) == 0 {
		ui.Info(emptyMessage)
		return "", nil
	}

	selectionItems := formatWorktreeSelectionItems(candidates)
	selection, cancelled, err := selectWithFzf(selectionItems)
	if err != nil {
		return "", err
	}
	if cancelled || selection == "" {
		ui.Info(selectionCancelledMsg)
		return "", nil
	}

	return parseSelectedBranch(selection), nil
}

func selectionCandidates(action string, worktrees []WorktreeInfo) ([]WorktreeInfo, string, error) {
	switch action {
	case "open", "relaunch":
		if len(worktrees) == 0 {
			return nil, noWorktreesMessage, nil
		}
		return worktrees, "", nil
	case "close", "attach":
		var active []WorktreeInfo
		for _, wt := range worktrees {
			if wt.SessionActive {
				active = append(active, wt)
			}
		}
		if len(active) == 0 {
			return nil, noActiveSessionsMessage, nil
		}
		return active, "", nil
	default:
		return nil, "", fmt.Errorf("unsupported action: %s", action)
	}
}

func formatWorktreeSelectionItems(worktrees []WorktreeInfo) []string {
	items := make([]string, 0, len(worktrees))
	for _, wt := range worktrees {
		status := "inactive"
		if wt.SessionActive {
			status = "active"
		}
		items = append(items, fmt.Sprintf("%s\t(%s)", wt.BranchName, status))
	}
	return items
}

func parseSelectedBranch(selection string) string {
	parts := strings.SplitN(strings.TrimSpace(selection), "\t", 2)
	if len(parts) == 0 {
		return ""
	}
	return strings.TrimSpace(parts[0])
}

func selectWithFzf(items []string) (string, bool, error) {
	selection, exitCode, err := runFzfSelector(items)
	if err != nil {
		return "", false, err
	}
	selection = strings.TrimSpace(selection)
	if exitCode == 130 {
		return "", true, nil
	}
	if exitCode != 0 && exitCode != 1 {
		return "", false, fmt.Errorf("fzf exited with code %d", exitCode)
	}
	if selection == "" {
		return "", true, nil
	}
	return selection, false, nil
}

func runFzf(items []string) (string, int, error) {
	cmd := exec.Command("fzf")
	cmd.Stdin = strings.NewReader(strings.Join(items, "\n") + "\n")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	output := strings.TrimRight(stdout.String(), "\n")
	if err == nil {
		return output, 0, nil
	}
	if errors.Is(err, exec.ErrNotFound) {
		return "", 0, fmt.Errorf("install fzf or pass a branch name")
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return output, exitErr.ExitCode(), nil
	}
	return "", 0, err
}
