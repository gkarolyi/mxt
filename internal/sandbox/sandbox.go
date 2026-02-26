package sandbox

import (
	"os/exec"
	"strings"
)

// Command builds an exec.Cmd for the provided command, optionally wrapped by a sandbox tool.
// When sandboxTool is set, the command is executed via `sh -c` to allow sandbox tool arguments.
func Command(sandboxTool string, command string, args ...string) *exec.Cmd {
	if strings.TrimSpace(sandboxTool) == "" {
		return exec.Command(command, args...)
	}

	full := CommandString(sandboxTool, command, args...)
	return exec.Command("sh", "-c", full)
}

// CommandString builds the shell command string for running command+args with an optional sandbox tool.
func CommandString(sandboxTool string, command string, args ...string) string {
	parts := make([]string, 0, len(args)+1)
	parts = append(parts, shellQuote(command))
	for _, arg := range args {
		parts = append(parts, shellQuote(arg))
	}

	commandLine := strings.Join(parts, " ")
	if strings.TrimSpace(sandboxTool) == "" {
		return commandLine
	}

	return strings.TrimSpace(sandboxTool) + " " + commandLine
}

func shellQuote(value string) string {
	if value == "" {
		return "''"
	}

	var sb strings.Builder
	sb.WriteByte('\'')
	for _, r := range value {
		if r == '\'' {
			sb.WriteString("'\\''")
			continue
		}
		sb.WriteRune(r)
	}
	sb.WriteByte('\'')
	return sb.String()
}
