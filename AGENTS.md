# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Rules

NEVER use cat << EOF or echo to output text, summaries, or reports. Use the chat interface directly.
Use the tk tool to manage work.
Use del instead of rm.

## Project Overview

`mxt` is a Go-based CLI for managing git worktrees paired with tmux sessions. It's designed for running parallel Claude Code or Codex sessions on macOS, giving each session its own isolated git worktree in a separate tmux session.

**Key Concept**: Each `mxt new` command creates a new git worktree (isolated branch + directory), copies config files into it, optionally runs setup commands, and launches a tmux session with customizable window/pane layouts.

## Architecture

### Package Layout

1. **CLI wiring** (`main.go`)
   - Cobra commands, argument parsing, and usage strings.
2. **Command handlers** (`internal/commands`)
   - `init`, `new`, `list`, `delete`, `sessions`, `config`, `help`, `version`.
3. **Configuration system** (`internal/config`)
   - Parses global (`~/.mxt/config`) and project-local (`.mxt`) configs.
   - Expands `~`, validates keys, and merges defaults.
4. **Git integration** (`internal/git`)
   - Repo root detection, branch lookup, worktree path helpers.
5. **Worktree management** (`internal/worktree`)
   - Creates git worktrees and copies configured files.
6. **Tmux integration** (`internal/tmux`)
   - Session creation, custom layout parsing, attach/kill helpers.
7. **Terminal integration** (`internal/terminal`)
   - Launches Terminal.app, iTerm2, Ghostty, or attaches in current terminal.
8. **UI helpers** (`internal/ui`)
   - Consistent messaging and formatting.

### Config File Format

Config files use TOML format with layout helpers:

- **Single-line**: `tmux_layout = "dev:hx|lazygit,server:bin/server,agent:"`
- **Multi-line**: `tmux_layout = """..."""` with windows on separate lines
- **Commas or newlines** separate windows in the layout

### Tmux Layout Syntax

Format: `window:pane1|pane2;window2:pane3`

- `,` or newline = window separator
- `:` = separates window name from panes
- `|` = pane separator (creates side-by-side split with `split-window -h`)
- Empty command = shell prompt

Example:

```toml
tmux_layout = """
  dev:hx|lazygit
  server:cd api && bin/server
  agent:
"""
```

This creates 3 windows:
- `dev` with 2 side-by-side panes (helix left, lazygit right)
- `server` with 1 pane running the server
- `agent` with 1 pane (empty shell for Claude Code)

## Testing and Development

- Prefer running `go test ./...` after code changes.
- Manual validation (when needed):
  1. `mxt new test-branch`
  2. `tmux ls` to verify session creation
  3. Test `terminal` config values (terminal/iterm2/ghostty/current)
  4. Validate config edge cases (missing config, invalid keys, special characters)

## Security Considerations

- Config values are parsed as plain text, not sourced by a shell.
- Shell metacharacters are rejected in most config values.
- `pre_session_cmd` and `tmux_layout` are allowed to contain shell metacharacters.
- Session names are sanitized for filesystem/tmux safety and escaped for AppleScript.

## Code Style

- Follow Go idioms and keep functions small and focused.
- Keep command behavior in `internal/commands`; lower-level operations in `internal/*` packages.
- Use `ui` helpers for user-facing output and consistent formatting.

## Important Implementation Details

### Pre-Session Command Hook

`pre_session_cmd` runs **after** worktree creation but **before** tmux session creation. This allows setup tasks (bundle install, npm install) to run before sessions open, and can abort on failure.

### Session Naming and Sanitization

Branches are sanitized for tmux session names and filesystem paths to avoid invalid characters and traversal risks.

### Custom Layout Parsing

Custom layouts are parsed into windows/panes and created via tmux. The first pane already exists for a new window; additional panes are created with `split-window -h` and then laid out evenly.

## Config File Locations

- **Global**: `~/.mxt/config` (or `$MXT_CONFIG_DIR/config`)
- **Project**: `.mxt` in repo root
- **Priority**: Project config overrides global config
- **Init**: `mxt init` creates global, `mxt init --local` creates project

Both use the same TOML format with the same keys:
- `worktree_dir`: Base directory for worktrees
- `terminal`: Which terminal app to use
- `copy_files`: Comma-separated list of files/globs to copy
- `pre_session_cmd`: Command to run after worktree setup
- `tmux_layout`: Custom window/pane layout

## Version and Help

- Version is hardcoded in `main.go`; update it when releasing.
- Help text lives in `internal/commands/help.go`.
