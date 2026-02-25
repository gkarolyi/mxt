# Feature Spec: muxtree init

This document captures the exact behavior of `muxtree init` command for feature parity in the Go reimplementation.

## Global Init: `muxtree init`

### Display

Shows ASCII art logo:
```
                       _
  _ __ ___  _   ___  _| |_ _ __ ___  ___
 | '_ ` _ \| | | \ \/ / __| '__/ _ \/ _ \
 | | | | | | |_| |>  <| |_| | |  __/  __/
 |_| |_| |_|\__,_/_/\_\\__|_|  \___|\___|
  Tmux Worktree Session Manager v1.0.0
```

### Prompts (in order)

**No prompt for worktree_dir** - just blank line with default value acceptance

**No prompt for terminal** - just blank line with default value acceptance

1. **Copy files prompt:**
```
▸ Enter files to copy into new worktrees (relative to repo root).
▸ Comma-separated, e.g.: .env,.env.local,CLAUDE.md
```

2. **Pre-session command prompt:**
```
▸ Optional: Command to run after worktree setup, before tmux session.
▸ Runs in worktree dir. Good for: bundle install, npm install, db:migrate
```

3. **Tmux layout prompt:**
```
▸ Optional: Tmux layout - define windows and panes.
▸ Format: window:cmd1|cmd2;window2:cmd3
▸ Example: dev:vim|;server:bin/server;agent:
```

### Success Output

```
✓ Config written to ~/.muxtree/config
```

Followed by the config file contents displayed to stdout.

### Generated Config File Format

Location: `~/.muxtree/config` (or `$MUXTREE_CONFIG_DIR/config`)

```
# muxtree configuration
# Generated on Wed 25 Feb 2026 10:56:33 AEDT

# Base directory for worktrees
worktree_dir=~/test-worktrees

# Terminal app: terminal | iterm2 | ghostty | current
terminal=iterm2

# Files to copy from repo root into new worktrees (comma-separated, relative to repo root)
# Supports glob patterns and directories
copy_files=.env,CLAUDE.md

# Command to run after worktree setup, before tmux session (optional)
# Runs in worktree directory. Use for setup tasks like: bundle install, npm install
pre_session_cmd=echo "setup complete"

# Tmux layout - define windows and panes (optional)
# Format: window_name:pane_cmd1|pane_cmd2;next_window:cmd
# Example: dev:vim|;server:bin/server;agent:
# - ';' separates windows
# - ':' separates window name from panes
# - '|' separates panes (horizontal split)
# - Empty command = shell prompt
# If not set, creates default layout: dev + agent windows
tmux_layout=dev:hx|lazygit,server:,agent:
```

### Default Values

- `worktree_dir`: `~/worktrees`
- `terminal`: `terminal`
- `copy_files`: empty
- `pre_session_cmd`: empty
- `tmux_layout`: empty

### Behavior Notes

1. Creates `~/.muxtree/` directory if it doesn't exist
2. If config already exists, prompts: "Overwrite? (y/N)"
3. Header comment includes generation timestamp
4. Each field has a descriptive comment above it
5. Empty values are written as blank (no default shown in file)

## Project Init: `muxtree init --local`

### Requirements

- Must be run inside a git repository
- If not in git repo, should error

### Display

Same ASCII logo as global init

### Prompts (in order)

Different prompt text mentioning "this project":

1. **Copy files prompt:**
```
▸ Enter files to copy into new worktrees for this project (relative to repo root).
▸ Comma-separated, e.g.: .env,.env.local,CLAUDE.md
```

2. **Pre-session command prompt:**
```
▸ Optional: Command to run after worktree setup, before tmux session.
▸ Runs in worktree dir. Good for: bundle install, npm install, db:migrate
```

3. **Tmux layout prompt:**
```
▸ Optional: Tmux layout - define windows and panes.
▸ Format: window:cmd1|cmd2;window2:cmd3
▸ Example: dev:vim|;server:bin/server;agent:
```

**Note:** No worktree_dir or terminal prompts for project config

### Success Output

```
✓ Project config written to /path/to/repo/.muxtree
```

Followed by the config file contents.

### Generated Project Config Format

Location: `.muxtree` in git repository root

```
# muxtree project config
# Generated on Wed 25 Feb 2026 10:56:41 AEDT

# Files to copy from repo root into new worktrees (comma-separated, relative to repo root)
# Supports glob patterns and directories
copy_files=.env.local,CLAUDE.md,package.json

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
tmux_layout=dev:hx|lazygit;server:npm start;agent:
```

### Differences from Global Config

1. Header says "muxtree project config" instead of "muxtree configuration"
2. Does NOT include `worktree_dir` or `terminal` fields (project overrides only)
3. Only includes: `copy_files`, `pre_session_cmd`, `tmux_layout`

## Help Text: `muxtree init --help`

Should display usage information for the init command.

## Color Coding

- Prompt lines (▸): Likely colored (cyan/blue)
- Success message (✓): Green
- Comments in config: Gray/dim
- Keys in config: May be colored

## Implementation Notes for mxt

1. ASCII logo should be exact match (or use similar styling)
2. Prompt messages should match exactly
3. Comment format and content should match
4. File paths should expand ~ correctly
5. Date format in "Generated on" should match or be similar
6. Success messages should be identical
7. Config file structure and comments must match
