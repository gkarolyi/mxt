// Package tmux handles tmux session creation and management.
package tmux

import (
	"fmt"
	"os"
	"strings"

	"github.com/gkarolyi/mxt/internal/sandbox"
)

// Window represents a tmux window with a name and list of pane commands.
type Window struct {
	Name  string
	Panes []string
}

// ParseLayout parses a custom tmux layout string into a slice of Windows.
//
// Format: window:pane1|pane2;window2:pane3
// Separators:
//   - ';' or ',' or newline: Separates windows
//   - ':': Separates window name from panes
//   - '|': Separates panes within a window
func ParseLayout(layout string) ([]Window, error) {
	if layout == "" {
		return []Window{}, nil
	}

	// Normalize separators: replace commas and newlines with semicolons
	layout = strings.ReplaceAll(layout, ",", ";")
	layout = strings.ReplaceAll(layout, "\n", ";")

	// Split into window specs
	windowSpecs := strings.Split(layout, ";")

	windows := make([]Window, 0)
	for _, spec := range windowSpecs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue // Skip empty specs
		}

		// Split by first colon to get window name and panes
		parts := strings.SplitN(spec, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid window spec (missing ':'): %s", spec)
		}

		windowName := strings.TrimSpace(parts[0])
		if windowName == "" {
			return nil, fmt.Errorf("empty window name in spec: %s", spec)
		}

		panesSpec := strings.TrimSpace(parts[1])

		// Split panes by pipe
		var panes []string
		if panesSpec == "" {
			// Empty panes spec means one empty pane
			panes = []string{""}
		} else {
			paneParts := strings.Split(panesSpec, "|")
			for _, pane := range paneParts {
				panes = append(panes, strings.TrimSpace(pane))
			}
		}

		windows = append(windows, Window{
			Name:  windowName,
			Panes: panes,
		})
	}

	return windows, nil
}

// SessionConfig contains configuration for creating a tmux session.
type SessionConfig struct {
	SessionName  string   // Name of the tmux session
	WorktreePath string   // Path to the worktree (working directory)
	SandboxTool  string   // Optional sandbox tool prefix
	RunCommand   string   // Optional command to run in agent window
	CustomLayout string   // Optional custom layout string
	WindowNames  []string // Resulting window names (populated after creation)
}

// CreateDefaultLayout creates a tmux session with the default layout (dev + agent windows).
//
// Algorithm:
// 1. Create new detached session with first window
// 2. Rename first window to "dev"
// 3. Create second window named "agent"
// 4. If RunCommand provided, send it to agent window
// 5. Select dev window (make it active)
func CreateDefaultLayout(config *SessionConfig) error {
	// Step 1: Create new detached session
	cmd := sandbox.Command(config.SandboxTool, "tmux", "new-session", "-d", "-s", config.SessionName, "-c", config.WorktreePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tmux session: %w", err)
	}

	// Step 2: Rename first window to "dev"
	cmd = sandbox.Command(config.SandboxTool, "tmux", "rename-window", "-t", config.SessionName+":0", "dev")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to rename window to dev: %w", err)
	}

	// Step 3: Create second window named "agent"
	cmd = sandbox.Command(config.SandboxTool, "tmux", "new-window", "-t", config.SessionName, "-n", "agent", "-c", config.WorktreePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create agent window: %w", err)
	}

	// Step 4: Send command to agent window if provided
	if config.RunCommand != "" {
		cmd = sandbox.Command(config.SandboxTool, "tmux", "send-keys", "-t", config.SessionName+":agent", config.RunCommand, "Enter")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to send command to agent window: %w", err)
		}
	}

	// Step 5: Select dev window (make it active)
	cmd = sandbox.Command(config.SandboxTool, "tmux", "select-window", "-t", config.SessionName+":dev")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to select dev window: %w", err)
	}

	// Populate window names
	config.WindowNames = []string{"dev", "agent"}

	return nil
}

// HasSession checks if a tmux session exists.
// Returns true if the session exists, false otherwise.
func HasSession(sessionName, sandboxTool string) bool {
	cmd := sandbox.Command(sandboxTool, "tmux", "has-session", "-t", sessionName)
	return cmd.Run() == nil
}

// KillSession kills a tmux session if it exists.
// Returns nil if session was killed or didn't exist.
func KillSession(sessionName, sandboxTool string) error {
	if !HasSession(sessionName, sandboxTool) {
		return nil
	}

	cmd := sandbox.Command(sandboxTool, "tmux", "kill-session", "-t", sessionName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to kill session: %w", err)
	}

	return nil
}

// AttachToSession attaches to an existing tmux session in the current terminal.
// If windowName is provided, it selects that window before attaching.
func AttachToSession(sessionName, windowName, sandboxTool string) error {
	// Validate window name if provided
	if windowName != "" && windowName != "dev" && windowName != "agent" {
		return fmt.Errorf("Unknown window: %s (use dev or agent)", windowName)
	}

	// Select window if specified
	if windowName != "" {
		target := fmt.Sprintf("%s:%s", sessionName, windowName)
		cmd := sandbox.Command(sandboxTool, "tmux", "select-window", "-t", target)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to select window '%s': %w", windowName, err)
		}
	}

	// Attach to session
	// This replaces the current process with tmux attach
	cmd := sandbox.Command(sandboxTool, "tmux", "attach", "-t", sessionName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// CreateCustomLayout creates a tmux session with a custom layout defined by the user.
//
// Algorithm:
// 1. Parse layout string
// 2. Create first window with session
// 3. Create additional windows
// 4. For each window, create panes and send commands
// 5. If RunCommand provided and agent window exists, send it
// 6. Select first window
func CreateCustomLayout(config *SessionConfig) error {
	// Step 1: Parse layout string
	windows, err := ParseLayout(config.CustomLayout)
	if err != nil {
		return fmt.Errorf("invalid tmux layout: %w", err)
	}

	if len(windows) == 0 {
		return fmt.Errorf("tmux layout is empty")
	}

	// Step 2: Create first window with session
	firstWindow := windows[0]
	cmd := sandbox.Command(config.SandboxTool, "tmux", "new-session", "-d", "-s", config.SessionName, "-c", config.WorktreePath, "-n", firstWindow.Name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tmux session: %w", err)
	}

	// Step 3: Create additional windows
	for i := 1; i < len(windows); i++ {
		window := windows[i]
		cmd := sandbox.Command(config.SandboxTool, "tmux", "new-window", "-t", config.SessionName, "-n", window.Name, "-c", config.WorktreePath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create window '%s': %w", window.Name, err)
		}
	}

	// Step 4: For each window, create panes and send commands
	for _, window := range windows {
		// First pane already exists (created with window)
		// Send command to first pane if non-empty
		if len(window.Panes) > 0 && window.Panes[0] != "" {
			target := fmt.Sprintf("%s:%s.0", config.SessionName, window.Name)
			cmd := sandbox.Command(config.SandboxTool, "tmux", "send-keys", "-t", target, window.Panes[0], "Enter")
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to send command to pane 0 of window '%s': %w", window.Name, err)
			}
		}

		// Create additional panes (vertical splits)
		for i := 1; i < len(window.Panes); i++ {
			target := fmt.Sprintf("%s:%s", config.SessionName, window.Name)
			cmd := sandbox.Command(config.SandboxTool, "tmux", "split-window", "-h", "-t", target, "-c", config.WorktreePath)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to create pane %d in window '%s': %w", i, window.Name, err)
			}

			// Send command to new pane if non-empty
			if window.Panes[i] != "" {
				// After split, the new pane is the last one, but we can just send to the window
				// and tmux will send to the active pane
				cmd := sandbox.Command(config.SandboxTool, "tmux", "send-keys", "-t", target, window.Panes[i], "Enter")
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("failed to send command to pane %d of window '%s': %w", i, window.Name, err)
				}
			}
		}

		// Apply even-horizontal layout if window has multiple panes
		if len(window.Panes) > 1 {
			target := fmt.Sprintf("%s:%s", config.SessionName, window.Name)
			cmd := sandbox.Command(config.SandboxTool, "tmux", "select-layout", "-t", target, "even-horizontal")
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to apply layout to window '%s': %w", window.Name, err)
			}
		}
	}

	// Step 5: If RunCommand provided, send to agent window if it exists
	if config.RunCommand != "" {
		// Search for agent window
		for _, window := range windows {
			if window.Name == "agent" {
				target := fmt.Sprintf("%s:agent.0", config.SessionName)
				cmd := sandbox.Command(config.SandboxTool, "tmux", "send-keys", "-t", target, config.RunCommand, "Enter")
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("failed to send command to agent window: %w", err)
				}
				break
			}
		}
	}

	// Step 6: Select first window
	target := fmt.Sprintf("%s:%s", config.SessionName, firstWindow.Name)
	cmd = sandbox.Command(config.SandboxTool, "tmux", "select-window", "-t", target)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to select first window: %w", err)
	}

	// Populate window names
	config.WindowNames = make([]string, len(windows))
	for i, window := range windows {
		config.WindowNames[i] = window.Name
	}

	return nil
}
