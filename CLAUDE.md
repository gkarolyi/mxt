# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`muxtree` is a bash-based CLI tool for managing git worktrees paired with tmux sessions. It's designed specifically for running parallel Claude Code or Codex sessions on macOS, giving each session its own isolated git worktree in a separate tmux session.

**Key Concept**: Each `muxtree new` command creates a new git worktree (isolated branch + directory), copies config files into it, optionally runs setup commands, and launches a tmux session with customizable window/pane layouts.

## Architecture

### Single-File Design

The entire tool is implemented in a single bash script (`muxtree`) with no external dependencies beyond git, tmux, and macOS (Terminal.app, iTerm2, or Ghostty).

**Why single file?**
- Easy to install (just copy to PATH)
- No build system needed
- Self-contained with shell completion in separate directory

### Core Components

The script is organized into these logical sections:

1. **Configuration System** (`_parse_config`, `_parse_config_set`, `load_config`)
   - Parses both global (`~/.muxtree/config`) and project-local (`.muxtree`) configs
   - Supports multi-line array syntax with `[...]` brackets
   - Sanitizes values to prevent shell injection
   - Key insight: Commands in `pre_session_cmd` and `tmux_layout` are allowed to contain shell metacharacters, but other config values reject them for security

2. **Worktree Management** (`cmd_new`, `cmd_delete`, `cmd_list`)
   - Creates git worktrees under `$WORKTREE_DIR/<repo>/<branch>/`
   - Copies files from `copy_files` config (supports globs)
   - Runs `pre_session_cmd` hook after worktree creation, before tmux
   - Branch names are sanitized for both filesystem safety and tmux compatibility

3. **Tmux Session Creation** (`_launch_sessions`, `_create_default_layout`, `_create_custom_layout`)
   - Default: 2-window layout (dev + agent)
   - Custom: User-defined layout from `tmux_layout` config
   - Sessions named: `<repo>_<sanitized-branch>`
   - Key insight: `_create_custom_layout` parses the layout string, creates windows, then splits them into panes using `tmux split-window -h` (vertical/side-by-side splits with `|` separator)

4. **Terminal Integration** (`_open_terminal`)
   - Supports Terminal.app, iTerm2, Ghostty, and "current" (attach in active terminal)
   - Uses AppleScript for Terminal.app and iTerm2
   - Session names are escaped before embedding in AppleScript

5. **Session Management** (`cmd_sessions`)
   - Close/open/relaunch/attach operations
   - Can reopen sessions with `--run claude` or `--run codex`
   - Uses `has_tmux_session()` to check if session exists

### Config File Format

Config files use `key=value` format with special handling:

- **Single-line**: `tmux_layout=dev:hx|lazygit,server:bin/server,agent:`
- **Multi-line**: `tmux_layout=[...]` with windows on separate lines
- **Semicolons, commas, or newlines** all work as window separators in layout
- Parser (`_parse_config`) handles multi-line arrays by detecting `[` and accumulating lines until `]`

### Tmux Layout Syntax

Format: `window:pane1|pane2;window2:pane3`

- `;` or `,` or newline = window separator
- `:` = separates window name from panes
- `|` = pane separator (creates vertical/side-by-side split with `split-window -h`)
- Empty command = shell prompt

Example:
```bash
tmux_layout=[
  dev:hx|lazygit
  server:cd api && bin/server
  agent:
]
```

This creates 3 windows:
- `dev` with 2 side-by-side panes (helix left, lazygit right)
- `server` with 1 pane running the server
- `agent` with 1 pane (empty shell for Claude Code)

## Testing and Development

### No Automated Tests

This project has no test suite. Testing is done manually by:

1. Creating test worktrees: `muxtree new test-branch`
2. Verifying tmux sessions: `tmux ls`
3. Testing different terminal apps in config
4. Testing config edge cases (special characters, missing files, etc.)

### Testing Changes

When modifying the script:

1. **Config parsing**: Test both single-line and multi-line formats
2. **Layout creation**: Verify window/pane structure with `tmux list-windows -t <session>` and `tmux list-panes -t <session>`
3. **Terminal launching**: Test with different `terminal` config values
4. **Security**: Ensure shell metacharacters in config values are handled safely
5. **Edge cases**: Test with unusual branch names, missing directories, etc.

### Common Debugging Commands

```bash
# Inspect a tmux session structure
tmux list-windows -t my-app_feature-auth
tmux list-panes -t my-app_feature-auth

# Check if session exists
tmux has-session -t my-app_feature-auth

# Manually attach to inspect
tmux attach -t my-app_feature-auth
```

## Security Considerations

The script is designed to prevent shell injection:

1. **Config values** are parsed as plain text, not sourced
2. **Shell metacharacters** (`` ` ``, `$`, `;`, `|`, `&`) are rejected in most config values
3. **Exception**: `pre_session_cmd` and `tmux_layout` allow metacharacters since they're meant to run commands
4. **AppleScript**: Session names are escaped before embedding in osascript
5. **Filesystem**: Branch names are sanitized to prevent directory traversal

When adding new config options, decide whether they should:
- **Reject metacharacters** (if they're just data like paths)
- **Allow metacharacters** (if they're commands to execute)

## Code Style

Follow these conventions:

1. **Function naming**:
   - `cmd_*` = Command implementations (e.g., `cmd_new`, `cmd_delete`)
   - `_*` = Internal helpers (e.g., `_parse_config`, `_open_terminal`)
   - No prefix = Public utilities (e.g., `session_prefix`, `sanitize_branch_name`)

2. **Error handling**:
   - Use `die "message"` for fatal errors (prints to stderr and exits)
   - Use `warn "message"` for non-fatal warnings
   - Use `info`, `success`, `error` for colored output

3. **Variable quoting**: Always quote variables unless you explicitly want word splitting

4. **Command safety**: Use `--` separator with `rm`, `cp`, `mkdir` to handle edge-case filenames

5. **No external dependencies**: Don't add dependencies beyond git, tmux, and macOS built-ins

## Important Implementation Details

### Pre-Session Command Hook

The `pre_session_cmd` runs **after** worktree creation but **before** tmux session creation. This timing is critical because:
- Files are already copied (can read `.env`, etc.)
- Can run setup that affects the session (e.g., `bundle install`)
- If it fails, user can abort before tmux session is created

Implementation in `cmd_new`:
```bash
if [[ -n "$PRE_SESSION_CMD" ]]; then
    info "Running pre-session command..."
    if (cd "$wt_path" && eval "$PRE_SESSION_CMD"); then
        success "Pre-session command completed"
    else
        warn "Pre-session command failed"
        read -rp "Continue anyway? (y/N) " confirm
        [[ "$confirm" =~ ^[Yy]$ ]] || die "Aborted"
    fi
fi
```

### Session Naming and Sanitization

Branches are sanitized differently for different contexts:
- **Tmux sessions**: Only alphanumeric, underscore, dash allowed (tmux compatibility)
- **Filesystem paths**: Same restrictions (prevent traversal)

Implementation: `sanitize_branch_name()` and `sanitize_session_name()` functions.

### Custom Layout Parsing

The `_create_custom_layout` function is complex because it:
1. Splits the layout string by `;` (or `,`) to get windows
2. For each window, splits by `:` to get name and panes
3. For each pane spec, splits by `|` to get individual panes
4. Creates first pane in window, then uses `split-window -h` for subsequent panes
5. Applies `even-horizontal` layout to distribute panes evenly

Key gotcha: The first pane already exists when you create a window, so you only split for panes 2+.

### Terminal App Support

Four terminal modes:
1. **terminal** (default): Opens new Terminal.app window with AppleScript
2. **iterm2**: Opens new iTerm2 window with AppleScript
3. **ghostty**: Uses `open -a Ghostty --args -e tmux attach`
4. **current**: Attaches in the currently active terminal (works with any app)

The "current" mode is useful when Claude Code is already running in a terminal and you want to stay in that terminal.

## Config File Locations

- **Global**: `~/.muxtree/config` (or `$MUXTREE_CONFIG_DIR/config`)
- **Project**: `.muxtree` in repo root
- **Priority**: Project config overrides global config
- **Init**: `muxtree init` creates global, `muxtree init --local` creates project

Both use the same key=value format with the same keys:
- `worktree_dir`: Base directory for worktrees
- `terminal`: Which terminal app to use
- `copy_files`: Comma-separated list of files/globs to copy
- `pre_session_cmd`: Command to run after worktree setup
- `tmux_layout`: Custom window/pane layout

## Version and Help

- Version is hardcoded in `VERSION` variable at top of script
- Update it when making releases
- Help text is in `print_help()` function - keep it in sync with features
