// Package terminal handles terminal application integration.
package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gkarolyi/mxt/internal/ui"
)

// Open opens a terminal window and attaches to the specified tmux session.
// The terminalType parameter specifies which terminal app to use:
// - "terminal": macOS Terminal.app (default)
// - "iterm2": iTerm2
// - "ghostty": Ghostty terminal
// - "current": Attach in currently active terminal
func Open(terminalType, sessionName string) error {
	switch terminalType {
	case "terminal", "":
		return openTerminalApp(sessionName)
	case "iterm2":
		return openITerm2(sessionName)
	case "ghostty":
		return openGhostty(sessionName)
	case "current":
		return openCurrent(sessionName)
	default:
		return fmt.Errorf("unknown terminal type: %s (use terminal, iterm2, ghostty, or current)", terminalType)
	}
}

// openTerminalApp opens a new Terminal.app window with tmux attached.
// Uses AppleScript via osascript.
func openTerminalApp(sessionName string) error {
	// Escape session name for AppleScript
	escapedSession := escapeForAppleScript(sessionName)

	script := fmt.Sprintf(`tell application "Terminal"
    activate
    do script "tmux attach -t %s"
end tell`, escapedSession)

	cmd := exec.Command("osascript", "-e", script)
	cmd.Stderr = nil // Suppress stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open Terminal.app: %w", err)
	}

	return nil
}

// openITerm2 opens a new iTerm2 window with tmux attached.
// Uses AppleScript via osascript.
func openITerm2(sessionName string) error {
	// Escape session name for AppleScript
	escapedSession := escapeForAppleScript(sessionName)

	script := fmt.Sprintf(`tell application "iTerm"
    activate
    create window with default profile
    tell current session of current window
        write text "tmux attach -t %s"
    end tell
end tell`, escapedSession)

	cmd := exec.Command("osascript", "-e", script)
	cmd.Stderr = nil // Suppress stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open iTerm2: %w", err)
	}

	return nil
}

// openGhostty opens Ghostty terminal with tmux attached.
// Uses the 'open' command with --args flag.
func openGhostty(sessionName string) error {
	cmd := exec.Command("open", "-a", "Ghostty", "--args", "-e", "tmux", "attach", "-t", sessionName)

	if err := cmd.Run(); err != nil {
		// Provide fallback instructions
		ui.Warn("Failed to open Ghostty. Ensure Ghostty.app is installed.")
		ui.Warn("Falling back to current terminal...")
		ui.Info(fmt.Sprintf("Run: tmux attach -t %s", sessionName))
		return err
	}

	return nil
}

// openCurrent attaches to the tmux session in the currently active terminal.
// This replaces the current process with tmux attach.
func openCurrent(sessionName string) error {
	ui.Info(fmt.Sprintf("Attaching to session in current terminal: %s", ui.BoldText(sessionName)))

	cmd := exec.Command("tmux", "attach", "-t", sessionName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		ui.Warn(fmt.Sprintf("Could not attach automatically. Run: tmux attach -t %s", sessionName))
		return err
	}

	return nil
}

// escapeForAppleScript escapes special characters in a string for use in AppleScript.
// Escapes backslashes and double quotes.
func escapeForAppleScript(s string) string {
	// Escape backslashes first (must be done before quotes)
	s = strings.ReplaceAll(s, `\`, `\\`)
	// Escape double quotes
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}
