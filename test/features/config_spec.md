# Feature Spec: muxtree config

This document captures the exact behavior of `muxtree config` command for feature parity in the Go reimplementation.

## Display Format

### With Global Config Only

```
Global config: /path/to/.muxtree/config
─────────────────────────────────
[config file contents]

No project config. Use muxtree init --local to create one.
```

### With Both Global and Project Config

```
Global config: /path/to/.muxtree/config
─────────────────────────────────
[config file contents]

Project config: /path/to/repo/.muxtree (active)
─────────────────────────────────
[config file contents]
```

### With No Global Config

```
No global config. Use muxtree init to create one.

[Project config section if exists]
```

## Behavior Details

### Global Config Section

1. Line 1: `Global config: <absolute-path>`
2. Line 2: Separator line of dashes (33 characters: `─────────────────────────────────`)
3. Lines 3+: Full contents of the config file
4. Blank line after config contents

### Project Config Section

1. Blank line separator from previous section
2. Line: `Project config: <absolute-path> (active)`
   - Note: "(active)" suffix indicates this config is being used
3. Separator line of dashes (33 characters)
4. Full contents of the config file

### Missing Config Messages

**No global config:**
```
No global config. Use muxtree init to create one.
```

**No project config (when in a repo or any directory):**
```
No project config. Use muxtree init --local to create one.
```

**Not in a git repository:**
When not in a git repository, only show global config section and the "No project config" message.

## Example Outputs

### Example 1: Both Configs Present

```
Global config: /Users/username/.muxtree/config
─────────────────────────────────
# muxtree configuration
# Generated on Tue 24 Feb 2026 16:29:07 AEDT

# Base directory for worktrees
worktree_dir=~/Code/worktrees

# Terminal app: terminal | iterm2 | ghostty | current
terminal=current

# Files to copy from repo root into new worktrees (comma-separated, relative to repo root)
# Supports glob patterns and directories
copy_files=.env,.env.local,CLAUDE.md

Project config: /Users/username/projects/myapp/.muxtree (active)
─────────────────────────────────
# muxtree project config
# Generated on Wed 25 Feb 2026 10:56:41 AEDT

# Files to copy from repo root into new worktrees (comma-separated, relative to repo root)
# Supports glob patterns and directories
copy_files=.env.local,CLAUDE.md,package.json

# Command to run after worktree setup, before tmux session (optional)
# Runs in worktree directory. Use for setup tasks like: bundle install, npm install
pre_session_cmd=npm install

# Tmux layout - define windows and panes (optional)
# Format: window_name:pane_cmd1|pane_cmd2;next_window:cmd
# Example: dev:vim|;server:bin/server;agent:
# - ';' separates windows
# - ':' separates window name from panes
# - '|' separates panes (horizontal split)
# - Empty command = shell prompt
# If not set, creates default layout: dev + agent windows
tmux_layout=dev:hx|lazygit;server:npm start;agent:
```

### Example 2: Only Global Config

```
Global config: /Users/username/.muxtree/config
─────────────────────────────────
# muxtree configuration
# Generated on Tue 24 Feb 2026 16:29:07 AEDT

# Base directory for worktrees
worktree_dir=~/Code/worktrees

# Terminal app: terminal | iterm2 | ghostty | current
terminal=terminal

# Files to copy from repo root into new worktrees (comma-separated, relative to repo root)
# Supports glob patterns and directories
copy_files=

# Command to run after worktree setup, before tmux session (optional)
# Runs in worktree directory. Use for setup tasks like: bundle install, npm install
pre_session_cmd=

# Tmux layout - define windows and panes (optional)
# Format: window_name:pane_cmd1|pane_cmd2;next_window:cmd
# Example: dev:vim|;server:bin/server;agent:
# - ';' separates windows
# - ':' separates window name from panes
# - '|' separates panes (horizontal split)
# - Empty command = shell prompt
# If not set, creates default layout: dev + agent windows
tmux_layout=

No project config. Use muxtree init --local to create one.
```

### Example 3: No Configs

```
No global config. Use muxtree init to create one.

No project config. Use muxtree init --local to create one.
```

## Implementation Notes

1. Config paths should be absolute paths
2. Separator line is exactly 33 dash characters: `─────────────────────────────────`
3. "(active)" suffix is only on project config line
4. Blank line between sections
5. Display entire file contents verbatim (including comments and empty lines)
6. Message wording must match exactly
7. No color output in file contents (just plain text display)

## Command Options

### `muxtree config --help`

Should display usage information for the config command.

### Error Cases

- If global config exists but can't be read: Show error message
- If project config exists but can't be read: Show error message
- Path resolution errors should be handled gracefully

## Color Coding (Optional)

- Config paths: May be colored (e.g., cyan)
- Separator lines: May be colored (e.g., dim/gray)
- "(active)" suffix: May be colored (e.g., green)
- Missing config messages: Default color or dim
