# mxt Feature Specification

## Overview

This document provides a complete specification for `mxt`, a Go reimplementation of the `muxtree` bash tool. The goal is 100% feature parity with muxtree, maintaining the same command-line interface, inputs, outputs, and behaviors.

**Purpose**: Manage git worktrees paired with tmux sessions, designed for running parallel Claude Code or Codex sessions on macOS. Each session gets its own isolated git worktree in a separate tmux session.

**Version Reference**: muxtree v1.0.0

---

## Development Methodology

### Test-Driven Development (TDD)

**Where appropriate, write unit tests first:**
- Config parsing logic
- Layout string parsing
- Branch name sanitization
- Session name generation
- Path calculation
- Multi-line array parsing
- Security validation (metacharacter detection)

**TDD Process:**
1. Write failing test that specifies the expected behavior
2. Implement minimal code to make test pass
3. Refactor while keeping tests green

**Integration tests can be written after implementation** for:
- Git operations (require actual git repos)
- Tmux operations (require tmux running)
- Terminal integration (require terminal apps)
- File system operations (require actual files)

### Feature Specification Tests

**Before implementing each command, create a feature spec file** that captures the exact terminal output:

**Location**: `test/features/<command>_spec.md`

**Format**:
```markdown
# Feature Spec: <command-name>

## Test Case: <scenario-name>

**Input:**
```
$ mxt <command> <args>
```

**Expected Output:**
```
<exact terminal output including colors/formatting>
```

**Exit Code:** 0

## Test Case: <another-scenario>
...
```

**Required Test Cases per Command:**
- Happy path (success case)
- Error cases (invalid input, missing dependencies, etc.)
- Edge cases (special characters, empty values, etc.)

**Implementation Completion Criteria:**
A command is considered fully implemented when:
1. `mxt <command>` produces **exactly the same output** as `muxtree <command>` for all test cases
2. Exit codes match
3. Side effects match (files created, git operations, tmux sessions, etc.)

**Testing Process:**
1. Run `muxtree <command>` with specific inputs
2. Capture output, exit code, and side effects
3. Document in feature spec file
4. Implement `mxt <command>`
5. Run `mxt <command>` with same inputs
6. Verify output matches feature spec exactly

**Feature Spec Validation:**
Create a test harness that:
- Runs both `muxtree` and `mxt` with identical inputs
- Captures and compares outputs
- Reports differences
- Validates exit codes
- Checks side effects (worktree created, session exists, etc.)

### Project Management with tk

**All implementation work MUST be managed using the tk ticket system.**

**Ticket Structure:**
- Epic: `mux-aofc` - Go reimplementation with 100% feature parity
- Phases: 9 major phases (Phase 0 through Phase 8)
- Subtasks: 61+ atomic implementation tasks

**Workflow:**

1. **View available work:**
   ```bash
   tk ready                  # Show tickets ready to work on (deps resolved)
   tk blocked                # Show tickets blocked by dependencies
   tk list --status=open     # Show all open tickets
   ```

2. **Start working on a ticket:**
   ```bash
   tk start <ticket-id>      # Mark ticket as in_progress
   tk show <ticket-id>       # View ticket details
   ```

3. **Complete a ticket:**
   ```bash
   tk close <ticket-id>      # Mark ticket as closed
   ```

4. **Track progress:**
   ```bash
   tk show mux-aofc          # View epic status
   tk dep tree mux-aofc      # View dependency tree
   tk list --status=closed   # View completed tickets
   ```

5. **Add notes:**
   ```bash
   tk add-note <ticket-id> "Implemented X, found issue with Y"
   echo "Details" | tk add-note <ticket-id>
   ```

**Before starting any implementation task:**
1. Check `tk ready` for available tickets
2. Start with Phase 0 unless dependencies are clear
3. `tk start <ticket-id>` to claim the ticket
4. Follow the ticket's description and acceptance criteria
5. `tk close <ticket-id>` when acceptance criteria are met

**Creating new tickets:**
If you discover additional work not covered by existing tickets:
```bash
tk create "Task description" \
  --parent <phase-ticket-id> \
  --description "Detailed description" \
  --acceptance "Acceptance criteria"
```

**Ticket Dependencies:**
All phases have sequential dependencies:
- Phase 0 → Phase 1 → Phase 2 → ... → Phase 8
- Future enhancement tickets (from `todos` file) are blocked by epic `mux-aofc`
- Don't start Phase N until Phase N-1 is complete

**Viewing the Project State:**
```bash
# See overall progress
tk list --status=closed | wc -l    # Count closed tickets
tk list --status=open | wc -l      # Count open tickets

# See what's blocking you
tk blocked                          # Tickets with unresolved dependencies

# Find specific tickets
tk query '.[] | select(.title | contains("config"))'  # Find config-related tickets
```

**Important:** Always update ticket status as you work. This keeps the project organized and helps others see what's in progress and what's available.

---

## Architecture Requirements

### Language and Dependencies
- **Language**: Go (latest stable version)
- **Binary name**: `mxt`
- **External dependencies**:
  - git (command-line tool)
  - tmux (command-line tool)
  - macOS-specific: osascript (for Terminal.app/iTerm2 integration)
  - macOS-specific: `open` command (for Ghostty integration)

### Project Structure
```
mxt/
├── cmd/
│   └── mxt/
│       └── main.go          # Entry point
├── internal/
│   ├── config/              # Configuration loading and parsing
│   ├── worktree/            # Git worktree operations
│   ├── tmux/                # Tmux session management
│   ├── terminal/            # Terminal app integration
│   └── ui/                  # Color output and formatting
├── go.mod
└── go.sum
```

### Code Quality Requirements
- Well-structured, idiomatic Go code
- Unit tests for core logic (config parsing, layout parsing, sanitization)
- Error handling with clear, user-friendly messages
- Documentation comments for exported functions

---

## Global Constants and Defaults

```
VERSION = "1.0.0"  (or start at 2.0.0 to distinguish from bash version)
DEFAULT_WORKTREE_DIR = "$HOME/worktrees"
DEFAULT_TERMINAL = "terminal"
DEFAULT_COPY_FILES = ""
DEFAULT_PRE_SESSION_CMD = ""
DEFAULT_TMUX_LAYOUT = ""
```

### Environment Variables
- `MUXTREE_CONFIG_DIR`: Override default config directory (default: `~/.muxtree`)

---

## Configuration System

### Config File Locations

1. **Global config**: `~/.muxtree/config` (or `$MUXTREE_CONFIG_DIR/config`)
2. **Project config**: `.muxtree` in git repository root

### Config Loading Priority
1. Load defaults
2. Load global config (if exists) - overrides defaults
3. Detect if running inside a git repo
4. Load project config (if exists) - overrides global config

### Config File Format

**Format**: `key=value` (one per line)
- Comments: Lines starting with `#` (optionally preceded by whitespace)
- Empty lines: Ignored
- Whitespace: Leading/trailing whitespace trimmed from keys and values
- Multi-line arrays: Supported with `key=[...]` syntax

**Example single-line format**:
```
worktree_dir=~/Code/worktrees
terminal=iterm2
copy_files=.env,.env.local,CLAUDE.md
pre_session_cmd=npm install && npm run db:migrate
tmux_layout=dev:hx|lazygit,server:bin/server,agent:
```

**Example multi-line format**:
```
worktree_dir=~/Code/worktrees
terminal=current
copy_files=.env,.env.local,CLAUDE.md

pre_session_cmd=[
  npm install
  npm run db:migrate
]

tmux_layout=[
  dev:hx|lazygit
  server:cd api && bin/server
  agent:
]
```

### Multi-line Array Parsing Rules

When a line matches `key=[`, begin multi-line accumulation:
1. Continue reading lines until finding a line with `]`
2. Concatenate all lines between `[` and `]` with spaces
3. The closing `]` can be on the same line as opening `[` (single-line array)
4. Example: `key=[ value1  value2  value3 ]` → `"value1 value2 value3"`

Special handling for `tmux_layout`:
- After accumulating multi-line value, normalize separators
- Replace commas with semicolons: `,` → `;`
- Replace multiple spaces before a window name with semicolon: `  +([a-z_]+:)` → `;$1`
- Clean up double semicolons: `;;` → `;`
- Trim leading/trailing semicolons

### Config Keys

| Key | Type | Description | Security |
|-----|------|-------------|----------|
| `worktree_dir` | string | Base directory for worktrees | No metacharacters |
| `terminal` | string | Terminal app: `terminal`, `iterm2`, `ghostty`, `current` | No metacharacters |
| `copy_files` | string | Comma-separated files/globs to copy | No metacharacters |
| `pre_session_cmd` | string | Command to run after worktree setup | **Allows metacharacters** |
| `tmux_layout` | string | Custom tmux layout definition | **Allows metacharacters** |

### Security: Shell Metacharacter Validation

**For most config values** (worktree_dir, terminal, copy_files):
- **Reject** if contains: `` ` ``, `$`, `;`, `|`, `&`, `$()`
- If detected, warn: `"Ignoring suspicious value for '<key>' in <file>"`
- Do not set the variable

**For command config values** (pre_session_cmd, tmux_layout):
- **Allow** all shell metacharacters (these are meant to be commands)

### Tilde Expansion

After loading config, expand `~` in `worktree_dir`:
- Replace leading `~` with `$HOME`
- Example: `~/worktrees` → `/Users/username/worktrees`

---

## Command: `init`

### Purpose
Create configuration files (global or project-local) with interactive prompts.

### Usage
```
mxt init [--local|-l]
```

### Flags
- `--local` or `-l`: Create project config (`.muxtree`) instead of global config

### Behavior

#### Global Init (no flags)

1. Display logo and version
2. Check if `~/.muxtree/config` exists
   - If yes: Display current content, prompt "Overwrite? (y/N)"
   - If user doesn't confirm 'y' or 'Y', exit gracefully
3. Create `~/.muxtree/` directory if it doesn't exist
4. Prompt for config values:
   - **Worktree base directory**: Default `~/worktrees`
   - **Terminal app**: Default `terminal` (options: terminal/iterm2/ghostty/current)
   - **Files to copy**: Comma-separated, relative to repo root, supports globs
   - **Pre-session command**: Optional, runs after worktree setup
   - **Tmux layout**: Optional, single-line format
5. Write config file with header comment showing generation timestamp
6. Display success message and show file contents
7. Config template (see Config File Template section below)

#### Project Init (--local flag)

1. Display logo and version
2. **Require git repo**: Exit with error if not inside a git repo
3. Determine repo root: `git rev-parse --show-toplevel`
4. Check if `<repo_root>/.muxtree` exists
   - If yes: Display current content, prompt "Overwrite? (y/N)"
   - If user doesn't confirm 'y' or 'Y', exit gracefully
5. Prompt for config values (subset of global):
   - **Files to copy**: Project-specific files to copy
   - **Pre-session command**: Project-specific setup command
   - **Tmux layout**: Project-specific layout
6. Write config file with header comment
7. Display success message and show file contents
8. Config template (see Config File Template section below)

### Config File Templates

**Global config template**:
```
# muxtree configuration
# Generated on <current_date>

# Base directory for worktrees
worktree_dir=<user_input>

# Terminal app: terminal | iterm2 | ghostty | current
terminal=<user_input>

# Files to copy from repo root into new worktrees (comma-separated, relative to repo root)
# Supports glob patterns and directories
copy_files=<user_input>

# Command to run after worktree setup, before tmux session (optional)
# Runs in worktree directory. Use for setup tasks like: bundle install, npm install
pre_session_cmd=<user_input>

# Tmux layout - define windows and panes (optional)
# Multi-line format (more readable):
# tmux_layout=[
#   dev:hx|lazygit
#   server:bin/server
#   agent:
# ]
# Or single line: tmux_layout=dev:hx|lazygit,server:bin/server,agent:
#
# Syntax:
# - ',' or newline separates windows
# - ':' separates window name from panes
# - '|' separates panes (vertical split - side by side)
# - Empty command = shell prompt
# If not set, creates default layout: dev + agent windows
tmux_layout=<user_input>
```

**Project config template** (similar but without worktree_dir and terminal):
```
# muxtree project config
# Generated on <current_date>

# Files to copy from repo root into new worktrees (comma-separated, relative to repo root)
# Supports glob patterns and directories
copy_files=<user_input>

# Command to run after worktree setup, before tmux session (optional)
# Runs in worktree directory. Use for setup tasks like: bundle install, npm install
pre_session_cmd=<user_input>

# Tmux layout - define windows and panes (optional)
# Multi-line format (more readable):
# tmux_layout=[
#   dev:hx|lazygit
#   server:bin/server
#   agent:
# ]
# Or single line: tmux_layout=dev:hx|lazygit,server:bin/server,agent:
#
# Syntax:
# - ',' or newline separates windows
# - ':' separates window name from panes
# - '|' separates panes (vertical split - side by side)
# - Empty command = shell prompt
# If not set, creates default layout: dev + agent windows
tmux_layout=<user_input>
```

### Output Format

Logo:
```
                       _
  _ __ ___  _   ___  _| |_ _ __ ___  ___
 | '_ ` _ \| | | \ \/ / __| '__/ _ \/ _ \
 | | | | | | |_| |>  <| |_| | |  __/  __/
 |_| |_| |_|\__,_/_/\_\__|_|  \___|\___|

  Tmux Worktree Session Manager v<VERSION>
```

### Exit Codes
- `0`: Success
- `1`: Error (not in git repo for --local, or other error)

---

## Command: `config`

### Purpose
Display current configuration (global and/or project-local).

### Usage
```
mxt config
```

### Behavior

1. Load config (following standard priority: defaults → global → project)
2. Check for global config file
   - If exists: Display header "Global config: <path>", separator line, file contents
3. Check for project config file
   - If exists: Display header "Project config: <path> (active)", separator line, file contents
   - If global exists but no project: Display dim text "No project config. Use **muxtree init --local** to create one."
4. If neither exists: Display warning "No config found. Run **muxtree init** to create one.", exit with code 1

### Output Format

```
Global config: /Users/username/.muxtree/config
─────────────────────────────────
<file contents>

Project config: /path/to/repo/.muxtree (active)
─────────────────────────────────
<file contents>
```

### Color Formatting
- "Global config" / "Project config": Bold
- "(active)": Green
- Separator line: 33 dashes
- "No project config" hint: Dim
- "No config found" warning: Yellow warning icon + message

### Exit Codes
- `0`: Success (at least one config exists)
- `1`: No config found

---

## Command: `new`

### Purpose
Create a new git worktree with a new branch and launch tmux session.

### Usage
```
mxt new <branch-name> [--from <base-branch>] [--run <claude|codex>] [--bg]
```

### Arguments
- `<branch-name>`: **Required**. Name of the new branch to create.

### Flags
- `--from <base-branch>`: Base branch to branch from (default: auto-detect main/master)
- `--run <command>`: Auto-run command in agent window. Valid values: `claude`, `codex`
- `--bg`: Create session in background without opening terminal window

### Behavior

#### Prerequisites
1. **Must be inside a git repository**: Exit with error if not
2. Load configuration (defaults → global → project)

#### Validation Phase

1. **Validate --run command**: If provided, must be `claude` or `codex`, otherwise error
2. **Determine base branch**:
   - If `--from` provided: use that value
   - Otherwise: auto-detect main branch
     - Try: `git symbolic-ref refs/remotes/origin/HEAD` → extract branch name
     - Fallback: check if `main` exists, then `master`
     - Final fallback: use `main`
3. **Validate base branch exists**:
   - Check `git show-ref --verify refs/heads/<base-branch>`
   - Or check `git show-ref --verify refs/remotes/origin/<base-branch>`
   - If neither exists: Error "Base branch '<branch>' does not exist."
4. **Check if new branch already exists**:
   - Run `git show-ref --verify refs/heads/<branch>`
   - If exists: Error "Branch '<branch>' already exists. Use a different name, or delete it first."
5. **Determine worktree path**:
   - Path: `$WORKTREE_DIR/<repo-name>/<sanitized-branch>`
   - Sanitization: See "Branch Name Sanitization" section
6. **Check if worktree path already exists**:
   - If directory exists: Error "Worktree already exists at <path>"

#### Execution Phase

**Step 1: Create Worktree**
- Print: "▸ Creating worktree at <path>" (blue info icon)
- Create parent directory if needed: `mkdir -p $(dirname <path>)`
- Run: `git worktree add -b <branch> <path> <base-branch>`
- Print: "✓ Worktree created (branch <branch> from <base-branch>)" (green checkmark)
  - Branch name in cyan
  - Base branch in dim

**Step 2: Copy Config Files**
- If `COPY_FILES` is non-empty:
  - Print: "▸ Copying config files..."
  - Split `COPY_FILES` by comma
  - For each file pattern (trim whitespace):
    - Expand globs relative to repo root
    - If no matches found: Print "⚠ Not found: <pattern>" (yellow warning, dim pattern)
    - For each matched file:
      - Calculate relative path from repo root
      - Create destination parent directories
      - Copy file preserving attributes: `cp -a <src> <dest>`
      - Print: "✓ Copied <relative-path>" (green checkmark, dim path)

**Step 3: Run Pre-Session Command**
- If `PRE_SESSION_CMD` is non-empty:
  - Print: "▸ Running pre-session command..."
  - Print: "  <command>" (dim, indented)
  - Change to worktree directory
  - Execute command via shell
  - If success (exit code 0):
    - Print: "✓ Pre-session command completed" (green checkmark)
  - If failure (non-zero exit code):
    - Print: "⚠ Pre-session command failed (exit code: <code>)" (yellow warning)
    - Prompt: "Continue anyway? (y/N) "
    - If user doesn't confirm 'y' or 'Y':
      - Error: "Aborted due to pre-session command failure"
      - Exit with code 1

**Step 4: Launch Tmux Session**
- Print: "▸ Creating tmux session..."
- Determine session name: `<repo-name>_<sanitized-branch>`
- If `TMUX_LAYOUT` is configured:
  - Use custom layout (see "Custom Layout Creation")
- Else:
  - Use default layout (see "Default Layout Creation")
- Print: "✓ Created session <session-name> (windows: <window-list>)" (green checkmark, bold session name)
- If `--run` flag provided:
  - Send command to agent window (if it exists)
- If `--bg` flag NOT provided:
  - Open terminal window (see "Terminal Integration")
- Print blank line
- Print: "✓ Ready! Worktree: <path>" (green checkmark, cyan path)

### Branch Name Sanitization

**Purpose**: Ensure branch names are safe for:
1. Filesystem paths (no directory traversal)
2. Tmux session names (limited character set)

**Algorithm**:
1. Replace any character that is NOT alphanumeric, underscore, dash, or dot with dash: `[^a-zA-Z0-9._-]` → `-`
2. Strip leading dash if present

**Examples**:
- `feature/auth` → `feature-auth`
- `bug-fix-#123` → `bug-fix--123`
- `user@domain` → `user-domain`

### Exit Codes
- `0`: Success
- `1`: Error (validation failure, worktree creation failure, user abort)

---

## Command: `list` (alias: `ls`)

### Purpose
List all managed worktrees for the current repository with status information.

### Usage
```
mxt list
mxt ls
```

### Behavior

1. **Require git repo**: Exit with error if not inside git repository
2. Load configuration
3. Determine repository name: `basename $(git rev-parse --show-toplevel)`
4. Determine managed worktree base: `$WORKTREE_DIR/<repo-name>`
5. Display header:
   ```
   Worktrees for <repo-name>
   ════════════════════════════════════════════════════════════════
   ```
6. If managed directory doesn't exist:
   - Print: "▸ No worktrees found. Use **muxtree new <branch>** to create one."
   - Exit with code 0
7. Parse `git worktree list` output:
   - For each worktree:
     - Extract worktree path (first field)
     - Extract branch name (inside `[...]`)
     - Skip if no branch name (detached HEAD)
     - Skip if path doesn't start with managed base
     - Display worktree info (see Output Format below)
8. If no managed worktrees found:
   - Print: "▸ No managed worktrees found. Use **muxtree new <branch>** to create one."

### Output Format

For each worktree:
```

  <branch-name>  +<insertions> -<deletions>
  <worktree-path>
  Session: <status> <session-name>
```

- Branch name: Bold, cyan
- Change stats: Green `+<num>`, red `-<num>`
- Path: Dim
- Session status: Green `●` (active) or dim `○` (inactive)
- Session name: Regular text

### Change Statistics Calculation

For a given worktree directory:
1. Run `git -C <dir> diff --stat HEAD` → get unstaged insertions/deletions
2. Run `git -C <dir> diff --cached --stat HEAD` → get staged insertions/deletions
3. Parse last line of each output: `<n> insertion`, `<m> deletion`
4. Extract numbers, default to 0 if not found
5. Sum: `total_insertions = unstaged + staged`, `total_deletions = unstaged + staged`

### Session Status Check

For each worktree:
- Determine session name: `<repo-name>_<sanitized-branch>`
- Check if session exists: `tmux has-session -t <session-name>` (exit code 0 = exists)
- Display green `●` if exists, dim `○` if not

### Exit Codes
- `0`: Success

---

## Command: `delete` (alias: `rm`)

### Purpose
Delete a worktree, kill its tmux session, and delete the local branch.

### Usage
```
mxt delete <branch-name> [--force|-f]
mxt rm <branch-name> [--force|-f]
```

### Arguments
- `<branch-name>`: **Required**. Name of the branch/worktree to delete.

### Flags
- `--force` or `-f`: Skip confirmation prompt

### Behavior

1. **Require git repo**: Exit with error if not inside git repository
2. Load configuration
3. Determine repository name and worktree path
4. **Validate worktree exists**:
   - Check if directory exists at worktree path
   - If not: Error "Worktree not found: <path>"
5. **Calculate change statistics** (same algorithm as `list` command)
6. **Display summary**:
   ```

     Branch:    <branch-name>
     Path:      <worktree-path>
     Changes:   +<insertions> -<deletions>

   ```
7. **Confirmation prompt** (unless `--force`):
   - Print: "⚠ This will remove the worktree and delete the local branch."
   - Prompt: "Are you sure? (y/N) "
   - If user doesn't confirm 'y' or 'Y':
     - Print: "▸ Cancelled."
     - Exit with code 0
8. **Kill tmux session**:
   - Determine session name: `<repo-name>_<sanitized-branch>`
   - If session exists: `tmux kill-session -t <session-name>`
   - Print: "✓ Killed session <session-name>" (bold session name)
9. **Remove worktree**:
   - Print: "▸ Removing worktree..."
   - Try: `git worktree remove <path> --force`
   - If fails:
     - Print: "⚠ git worktree remove failed, cleaning up manually..."
     - Run: `rm -rf <path>`
     - Run: `git worktree prune`
   - Print: "✓ Worktree removed"
10. **Delete branch**:
    - Print: "▸ Deleting branch <branch>..." (cyan branch name)
    - Run: `git branch -D <branch>`
    - If success: Print "✓ Branch deleted"
    - If fails: Print "⚠ Branch may have already been deleted"
11. **Clean up empty repo directory**:
    - Check if `$WORKTREE_DIR/<repo-name>` is empty
    - If empty: Remove directory with `rmdir`
12. Print blank line
13. Print: "✓ Done."

### Exit Codes
- `0`: Success (or user cancelled)
- `1`: Error (not in git repo, worktree not found)

---

## Command: `sessions` (alias: `s`)

### Purpose
Manage tmux sessions for existing worktrees.

### Usage
```
mxt sessions <action> <branch-name> [options]
mxt s <action> <branch-name> [options]
```

### Actions

#### `open` (aliases: `launch`, `start`)

Create tmux session for an existing worktree and open terminal.

**Usage**: `mxt sessions open <branch> [--run <claude|codex>] [--bg]`

**Behavior**:
1. Require git repo
2. Load configuration
3. Determine repository name and worktree path
4. Validate worktree exists at path, error if not
5. Determine session name
6. Check if session already exists:
   - If exists: Print "⚠ Session <session> already exists", exit
7. Create tmux session:
   - Use custom layout if configured, else default layout
   - If `--run` provided: Send command to agent window
8. Open terminal (unless `--bg`)

**Exit codes**: 0 (success), 1 (error)

#### `close` (aliases: `kill`, `stop`)

Kill tmux session for a worktree.

**Usage**: `mxt sessions close <branch>`

**Behavior**:
1. Require git repo
2. Load configuration
3. Determine session name
4. Kill session if exists (see "Session Killing")
5. Print: "✓ Killed session <session>"

**Exit codes**: 0 (success), 1 (error)

#### `relaunch` (aliases: `restart`)

Kill and recreate tmux session.

**Usage**: `mxt sessions relaunch <branch> [--run <claude|codex>] [--bg]`

**Behavior**:
1. Execute close action
2. Execute open action

**Exit codes**: 0 (success), 1 (error)

#### `attach`

Attach to an existing tmux session in current terminal.

**Usage**: `mxt sessions attach <branch> [dev|agent]`

**Arguments**:
- Optional window name: `dev` or `agent`

**Behavior**:
1. Require git repo
2. Load configuration
3. Determine session name
4. Check if session exists:
   - If not: Error "Session not found: <session>"
5. If window name provided:
   - Validate window name is `dev` or `agent`
   - If invalid: Error "Unknown window: <window> (use dev or agent)"
   - Select window: `tmux select-window -t <session>:<window>`
6. Attach to session: `tmux attach -t <session>`

**Exit codes**: 0 (success), 1 (error)

### Common Options for open/relaunch
- `--run <claude|codex>`: Auto-run command in agent window
- `--bg`: Create session without opening terminal

### Exit Codes
- `0`: Success
- `1`: Error (invalid action, worktree not found, etc.)

---

## Tmux Session Management

### Default Layout Creation

**Windows**:
1. `dev`: First window (default selected)
2. `agent`: Second window

**Algorithm**:
1. Create new detached session: `tmux new-session -d -s <session> -c <worktree-path>`
2. Rename first window: `tmux rename-window -t <session>:0 dev`
3. Create second window: `tmux new-window -t <session> -n agent -c <worktree-path>`
4. If `--run` command provided:
   - Send command to agent window: `tmux send-keys -t <session>:agent '<command>' Enter`
5. Select dev window: `tmux select-window -t <session>:dev`

**Success message**: `"✓ Created session <session> (windows: dev, agent)"`

### Custom Layout Creation

**Input format**: `window:pane1|pane2;window2:pane3`

**Separators**:
- `;` or `,` or newline: Separates windows
- `:`: Separates window name from pane commands
- `|`: Separates panes within a window (creates vertical/side-by-side splits)

**Algorithm**:
1. Parse layout string:
   - Split by `;` to get window specs
   - For each window spec:
     - Trim whitespace
     - Skip if empty
     - Split by first `:` to get window name and panes spec
     - Trim window name
2. Create first window:
   - `tmux new-session -d -s <session> -c <worktree-path> -n <window-name>`
3. For each subsequent window:
   - `tmux new-window -t <session> -n <window-name> -c <worktree-path>`
4. For each window, create panes:
   - Split panes spec by `|`
   - First pane: Already exists, just send command if non-empty
     - `tmux send-keys -t <session>:<window>.0 '<command>' Enter`
   - For each additional pane:
     - Create vertical split: `tmux split-window -h -t <session>:<window> -c <worktree-path>`
     - Send command if non-empty: `tmux send-keys -t <session>:<window> '<command>' Enter`
   - If window has multiple panes:
     - Apply even layout: `tmux select-layout -t <session>:<window> even-horizontal`
5. If `--run` command provided:
   - Search window list for window named "agent"
   - If found: `tmux send-keys -t <session>:agent.0 '<command>' Enter`
6. Select first window: `tmux select-window -t <session>:<first-window-name>`

**Success message**: `"✓ Created session <session> (windows: <space-separated window names>)"`

### Session Naming

**Format**: `<repo-name>_<sanitized-branch>`

**Sanitization**: Same algorithm as branch name sanitization
- Replace `[^a-zA-Z0-9_-]` with `-`
- Strip leading dash

### Session Killing

**Algorithm**:
1. Determine session name: `<repo-name>_<sanitized-branch>`
2. Check if session exists: `tmux has-session -t <session>`
3. If exists:
   - Kill session: `tmux kill-session -t <session>`

---

## Terminal Integration

### Terminal Types

1. **`terminal`** (default): macOS Terminal.app
2. **`iterm2`**: iTerm2
3. **`ghostty`**: Ghostty terminal
4. **`current`**: Attach in currently active terminal

### Terminal Launch Behavior

#### Terminal.app (`terminal`)

**Method**: AppleScript via `osascript`

**Script**:
```applescript
tell application "Terminal"
    activate
    do script "tmux attach -t <session-name>"
end tell
```

**Session name escaping**:
- Escape backslashes: `\` → `\\`
- Escape double quotes: `"` → `\"`

**Error handling**: Suppress stderr, continue on failure

#### iTerm2 (`iterm2`)

**Method**: AppleScript via `osascript`

**Script**:
```applescript
tell application "iTerm"
    activate
    create window with default profile
    tell current session of current window
        write text "tmux attach -t <session-name>"
    end tell
end tell
```

**Session name escaping**: Same as Terminal.app

**Error handling**: Suppress stderr, continue on failure

#### Ghostty (`ghostty`)

**Method**: `open` command

**Command**:
```bash
open -a Ghostty --args -e tmux attach -t <session-name>
```

**Error handling**:
- If command fails:
  - Print: "⚠ Failed to open Ghostty. Ensure Ghostty.app is installed."
  - Print: "⚠ Falling back to current terminal..."
  - Print: "▸ Run: tmux attach -t <session>"

#### Current Terminal (`current`)

**Method**: Direct tmux attach

**Behavior**:
1. Print: "▸ Attaching to session in current terminal: <session>" (bold session name)
2. Run: `tmux attach -t <session>`
3. If fails:
   - Print: "⚠ Could not attach automatically. Run: tmux attach -t <session>"

---

## Command: `help`

### Purpose
Display comprehensive help information.

### Usage
```
mxt help
mxt -h
mxt --help
```

### Output

Display logo, then help text (see Help Text Content below).

### Help Text Content

```
USAGE
    mxt <command> [options]

COMMANDS
    init                              Set up global config (~/.muxtree/config)
        --local                       Create project config (.muxtree in repo root)
    config                            Show current config (global + project)

        Terminal options: terminal, iterm2, ghostty, current
        - terminal: macOS Terminal.app (new window)
        - iterm2: iTerm2 (new window)
        - ghostty: Ghostty (new tab via 'open -a')
        - current: Attach in currently active terminal (any app)

    new <branch> [options]             Create worktree + tmux session
        --from <branch>               Base branch (default: main/master)
        --run <claude|codex>          Auto-run command in agent window
        --bg                          Create session without opening terminal

    list                              List worktrees, diff stats, session status

    delete <branch> [--force]          Delete worktree and branch (with confirmation)

    sessions <action> <branch> [opts]  Manage tmux session for a worktree
        open   <branch> [--run cmd]   Create session & open terminal
        close  <branch>               Kill tmux session
        relaunch <branch> [--run cmd] Close + reopen session
        attach <branch> [dev|agent]   Attach to session (optionally select window)

    help                              Show this help message

EXAMPLES
    mxt init                          # Global setup
    mxt init --local                  # Project-specific copy_files
    mxt new feature-auth              # New worktree from main
    mxt new fix-bug --from develop    # New worktree from develop
    mxt new feature-ai --run claude   # Auto-launch claude code
    mxt new fix-bug --bg              # Create without opening terminals
    mxt list                          # Show all worktrees + status
    mxt sessions close feature-auth   # Kill tmux sessions
    mxt sessions relaunch fix-bug     # Restart sessions
    mxt delete feature-auth           # Remove worktree + branch

CONFIG
    Global:  ~/.muxtree/config
             (worktree_dir, terminal, copy_files, pre_session_cmd, tmux_layout)
    Project: .muxtree in repo root      (overrides global settings)
    Env:     MUXTREE_CONFIG_DIR=/path    (override global config dir)

    Hooks & Layout:
    - pre_session_cmd:  Runs after worktree setup, before tmux session
                        Good for: bundle install, npm install, db:migrate

    - tmux_layout:      Define custom tmux windows and panes
                        Multi-line format (recommended):
                        tmux_layout=[
                          dev:hx|lazygit
                          server:cd api && bin/server|cd ui && yarn start
                          logs:tail -f log/development.log
                          agent:
                        ]

                        Single-line format: dev:hx|lazygit,server:bin/server,agent:

                        Syntax:
                        - ',' or newline separates windows
                        - ':' separates window name from panes
                        - '|' separates panes (vertical split - side by side)
                        - Empty command = shell prompt

                        If not set, creates default: dev + agent windows
```

### Color Formatting
- Section headers (USAGE, COMMANDS, etc.): Bold
- Command names: Cyan
- Session names in examples: Regular

---

## Command: `version`

### Purpose
Display version number.

### Usage
```
mxt version
mxt -v
mxt --version
```

### Output
```
mxt v<VERSION>
```

Example: `mxt v1.0.0`

---

## UI and Output Formatting

### Color Codes

| Color | ANSI Code | Usage |
|-------|-----------|-------|
| Red | `\033[0;31m` | Errors, deletions |
| Green | `\033[0;32m` | Success, insertions, active status |
| Yellow | `\033[0;33m` | Warnings |
| Blue | `\033[0;34m` | Info messages |
| Cyan | `\033[0;36m` | Paths, branch names |
| Bold | `\033[1m` | Emphasis |
| Dim | `\033[2m` | De-emphasis |
| Reset | `\033[0m` | Reset to default |

### No-TTY Handling

If stdout is not a TTY:
- Disable all color codes (empty strings)
- Use plain text output

### Output Symbols

| Symbol | Usage | Example |
|--------|-------|---------|
| `▸` | Info/action | `▸ Creating worktree...` |
| `✓` | Success | `✓ Worktree created` |
| `⚠` | Warning | `⚠ Pre-session command failed` |
| `✗` | Error | `✗ Not inside a git repository` |
| `●` | Active status | Session active |
| `○` | Inactive status | Session inactive |

### Message Formatting Functions

**info(msg)**: Blue `▸` + message
**success(msg)**: Green `✓` + message
**warn(msg)**: Yellow `⚠` + message
**error(msg)**: Red `✗` + message (to stderr)
**die(msg)**: error(msg) + exit(1)

---

## Error Handling

### Git Repository Detection

**Check**: `git rev-parse --is-inside-work-tree`
- If fails: `die("Not inside a git repository. Run mxt from within your repo.")`

### Git Operations

All git commands should:
1. Check exit code
2. Capture stderr
3. Display user-friendly error messages
4. Exit with code 1 on failure

### Tmux Operations

**Session existence check**: `tmux has-session -t <session>`
- Exit code 0: Session exists
- Exit code 1: Session does not exist

**Command failures**: If tmux commands fail, display error and exit

### File System Operations

**Directory creation**: Always check return values
**File copying**: Handle missing source files gracefully (warn but continue)
**Path safety**: Use sanitization to prevent directory traversal

---

## Edge Cases and Special Scenarios

### Branch Names with Special Characters

**Example**: `feature/auth-api`
**Handling**: Sanitize for filesystem and tmux compatibility
**Result**: `feature-auth-api`

### Empty Config Values

**Behavior**: Treat as defaults
- Empty `copy_files`: No files copied
- Empty `pre_session_cmd`: No command run
- Empty `tmux_layout`: Use default layout

### Missing Config Files

**Global config missing**: Use defaults, no error
**Project config missing**: Use global/defaults, no error

### Worktree Already Exists

**Scenario**: User manually created a worktree at expected path
**Handling**: Error and exit, don't overwrite

### Branch Already Exists

**Scenario**: Branch name conflicts with existing branch
**Handling**: Error and exit, suggest using different name or deleting

### Pre-Session Command Failure

**Behavior**: Warn user, prompt for confirmation to continue
**User cancels**: Exit with code 1
**User confirms**: Continue with session creation

### Tmux Not Running

**Behavior**: Sessions created in detached mode, no need for existing tmux server

### Git Worktree Remove Fails

**Fallback**: Manually remove directory with `rm -rf`, then `git worktree prune`

### Empty Worktree Directory After Deletion

**Cleanup**: Remove parent `$WORKTREE_DIR/<repo-name>` if empty

### Window Name Doesn't Exist in Custom Layout

**Scenario**: User specifies `--run claude` but no "agent" window defined
**Handling**: Command not sent, session still created successfully

### Multiple Panes with Empty Commands

**Example**: `dev:|`
**Result**: Two panes with shell prompts, no commands sent

---

## Git Helper Functions

### get_repo_root()
**Command**: `git rev-parse --show-toplevel`
**Returns**: Absolute path to repository root

### get_repo_name()
**Algorithm**: `basename(get_repo_root())`
**Returns**: Repository directory name

### get_main_branch()
**Algorithm**:
1. Try: `git symbolic-ref refs/remotes/origin/HEAD` → extract last component
2. If fails or empty:
   - Check if `main` exists: `git show-ref --verify refs/heads/main`
   - If not, check if `master` exists: `git show-ref --verify refs/heads/master`
3. Fallback: Return `"main"`

**Returns**: Name of main branch (e.g., "main", "master")

---

## File Glob Expansion

### Pattern Support

Support standard glob patterns:
- `*`: Match any characters
- `?`: Match single character
- `**`: Recursive directory match (if supported by Go glob library)
- `[abc]`: Character class

### Expansion Algorithm

For each pattern in `copy_files`:
1. Trim whitespace
2. Resolve pattern relative to repository root
3. Expand glob (use filepath.Glob or equivalent)
4. If no matches: Warn and continue
5. For each match:
   - Calculate relative path from repo root
   - Create destination parent directories
   - Copy file preserving attributes

### Example

**Config**: `copy_files=.env*,config/*.yml`
**Repo structure**:
```
.env
.env.local
.env.test
config/
  database.yml
  redis.yml
```
**Copied files**: `.env`, `.env.local`, `.env.test`, `config/database.yml`, `config/redis.yml`

---

## Testing Requirements

### Unit Tests

**Config Parsing**:
- Single-line values
- Multi-line arrays
- Comments and empty lines
- Whitespace trimming
- Shell metacharacter validation
- Tilde expansion

**Layout Parsing**:
- Single window
- Multiple windows
- Multiple panes
- Empty commands
- Semicolon, comma, and newline separators
- Whitespace handling

**Sanitization**:
- Branch name sanitization
- Special characters
- Edge cases (empty string, all special chars)

**Session Naming**:
- Repo + branch combination
- Sanitization applied

### Integration Tests

**Config Loading**:
- Global only
- Project only
- Both (project overrides global)
- Neither (use defaults)

**Worktree Creation** (requires git):
- Basic creation
- With base branch
- File copying
- Glob expansion
- Pre-session command (success and failure)

**Tmux Session Creation** (requires tmux):
- Default layout
- Custom layout
- Multiple panes
- Command execution in panes

### Manual Testing Checklist

- [ ] All commands with --help flag
- [ ] init (global and project)
- [ ] config display
- [ ] new (various flags)
- [ ] list (with and without worktrees)
- [ ] delete (with and without confirmation)
- [ ] sessions (all actions)
- [ ] Terminal integration (all 4 types)
- [ ] Custom layouts (various configurations)
- [ ] Error cases (invalid branch, missing worktree, etc.)

---

## Implementation Order

**Note**: Each phase follows TDD principles. Write tests before implementation where appropriate. Create feature specs before implementing commands.

**Using tk for project management**: All phases and subtasks are tracked in the tk ticket system. Use `tk ready` to see available work, `tk start <id>` to begin work on a ticket, and `tk close <id>` when done. See "Project Management with tk" section above for details.

**Git Commits**: After completing each phase (when all features are implemented, all tests pass, and the phase ticket is closed), create a git commit with:
- Clear commit message describing the phase and what was implemented
- List of tickets closed
- Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>

Example commit message:
```
Implement Phase 2: Configuration System

- Config file parsing (single-line and multi-line)
- Security validation for shell metacharacters
- Config loading priority (defaults → global → project)
- init and config commands implemented

All tests pass. Feature parity validated with test harness.

Tickets closed: mux-30uf, mux-gaau, mux-2ao5, mux-vpqk, mux-s91d,
mux-szys, mux-ixni, mux-aut0, mux-atft, mux-x5mw, mux-mdgk

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

### Phase 0: Setup & Validation
1. Create `test/features/` directory structure
2. Create test harness script for comparing muxtree vs mxt output
3. Initialize Go project structure
4. Set up basic CI/CD (optional)

### Phase 1: Core Infrastructure
1. Project structure and build system (`go mod init`, directory layout)
2. **Write tests first**: UI/output formatting functions (colors, symbols)
   - Test color code output
   - Test TTY detection
   - Test message formatting (info, success, warn, error)
3. Implement UI/output formatting
4. CLI argument parsing framework (use cobra or similar)
5. Error handling framework

### Phase 2: Configuration (TDD)
1. **Write tests first**: Config file parsing
   - Single-line key=value parsing
   - Comment handling
   - Whitespace trimming
   - Multi-line array syntax with `[...]`
2. Implement config file parsing
3. **Write tests first**: Security validation
   - Metacharacter detection
   - Allow-list for command fields
4. Implement security validation
5. **Write tests first**: Config loading priority (defaults → global → project)
6. Implement config loading
7. **Feature spec**: `init` command
   - Capture muxtree init output (global and --local)
   - Document prompts and generated files
8. Implement `init` command
9. **Feature spec**: `config` command
   - Capture muxtree config output
10. Implement `config` command
11. Validate: mxt init/config match muxtree exactly

### Phase 3: Git Operations (TDD)
1. **Write tests first**: Branch name sanitization
   - Test special character replacement
   - Test leading dash removal
   - Test various branch name formats
2. Implement branch name sanitization
3. **Write tests first**: Session name generation
4. Implement session name generation
5. **Write tests first**: Worktree path calculation
6. Implement worktree path calculation
7. Implement git helper functions (repo detection, branch detection)
   - These can be tested with integration tests

### Phase 4: Worktree Management
1. **Feature spec**: `new` command (basic, no tmux)
   - new <branch>
   - new <branch> --from <base>
   - Error cases (invalid branch, already exists, etc.)
2. Implement basic `new` command structure
3. Implement git worktree creation
4. Implement file copying with glob expansion
   - Integration tests for glob expansion
5. Implement pre-session command execution
6. Validate: mxt new (without tmux) matches muxtree

### Phase 5: Tmux Integration (TDD + Feature Specs)
1. **Write tests first**: Custom layout parsing
   - Parse window:pane|pane;window2:pane
   - Handle semicolon, comma, newline separators
   - Test various layout strings
2. Implement layout parsing
3. **Feature spec**: Default tmux layout
   - Capture tmux session structure after muxtree new
   - Document window names, pane count
4. Implement default layout creation
5. **Feature spec**: Custom tmux layouts
   - Test various tmux_layout configurations
   - Capture resulting tmux structure
6. Implement custom layout creation
7. **Feature spec**: `list` command
   - Capture muxtree list output
   - Test with no worktrees, one worktree, multiple worktrees
   - Test diff stats display
8. Implement `list` command
9. **Feature spec**: Complete `new` command with tmux
   - new <branch> with default layout
   - new <branch> with custom layout
   - new <branch> --run claude
   - new <branch> --bg
10. Complete `new` command implementation
11. Validate: mxt new/list match muxtree exactly

### Phase 6: Session Management
1. **Feature spec**: `sessions open` command
2. Implement sessions open
3. **Feature spec**: `sessions close` command
4. Implement sessions close
5. **Feature spec**: `sessions relaunch` command
6. Implement sessions relaunch
7. **Feature spec**: `sessions attach` command
8. Implement sessions attach
9. **Feature spec**: Terminal integration (all 4 types)
   - Terminal.app behavior
   - iTerm2 behavior
   - Ghostty behavior
   - Current terminal behavior
10. Implement terminal integration
11. Validate: mxt sessions matches muxtree exactly

### Phase 7: Cleanup Commands
1. **Feature spec**: `delete` command
   - Normal delete with confirmation
   - delete --force
   - Error cases
2. Implement `delete` command
3. **Feature spec**: `help` command
   - Capture exact help text
4. Implement `help` command
5. **Feature spec**: `version` command
6. Implement `version` command
7. Validate: All commands match muxtree exactly

### Phase 8: Polish & Validation
1. Run full feature spec suite
2. Fix any discrepancies
3. Shell completion (bash, zsh)
4. Documentation
5. README with installation instructions
6. Final validation: 100% feature parity
7. Tag release v1.0.0 (or v2.0.0)

---

## Compatibility Notes

### Backward Compatibility

**Must maintain**:
- All command names and flags
- Config file format and locations
- Environment variable names
- Terminal integration behavior
- Tmux session naming

### Acceptable Differences

**Internal implementation**:
- Programming language (bash → Go)
- Code structure
- Internal function names

**Performance**:
- Go implementation may be faster
- Startup time differences acceptable

**Error messages**:
- Wording can differ slightly if meaning is clear
- Stack traces not needed (Go specific)

---

## Success Criteria

### Feature Parity Checklist

- [ ] All commands implemented
- [ ] All flags supported
- [ ] Config file format identical
- [ ] Tmux layout parsing identical
- [ ] Terminal integration works for all 4 types
- [ ] Session naming matches exactly
- [ ] File copying with globs works
- [ ] Pre-session command execution works
- [ ] Security validation matches
- [ ] Error messages are clear
- [ ] Colors and formatting match
- [ ] Help text is complete

### Validation Tests

1. Run identical commands with both tools
2. Verify identical config files work
3. Verify identical tmux sessions created
4. Verify identical file copying behavior
5. Verify identical error messages

### Migration Path

1. Install `mxt` alongside `muxtree`
2. Test all workflows with `mxt`
3. Verify 100% compatibility
4. Update shell aliases: `alias muxtree=mxt`
5. Remove original `muxtree` script after confidence period

---

## Future Enhancements (Out of Scope)

These are NOT required for initial v1.0 feature parity:

- Cross-platform support (Linux, Windows)
- Non-macOS terminal support
- Configuration file migration tool
- Interactive branch selection
- Worktree archiving/backup
- Remote worktree support
- Integration with other editors
- Plugin system
- Web UI

---

## Appendix: Example Layouts

### Example 1: Simple Two-Pane Dev

**Config**:
```
tmux_layout=dev:hx|lazygit
```

**Result**:
- Window "dev" with 2 side-by-side panes
- Left pane: `hx` command
- Right pane: `lazygit` command

### Example 2: Multi-Window Setup

**Config**:
```
tmux_layout=dev:hx|lazygit;server:bin/server;agent:
```

**Result**:
- Window "dev": 2 panes (hx | lazygit)
- Window "server": 1 pane (bin/server)
- Window "agent": 1 pane (shell)

### Example 3: Multi-Line Format

**Config**:
```
tmux_layout=[
  dev:hx|lazygit
  server:cd api && bin/server|cd ui && yarn start
  logs:tail -f log/development.log
  agent:
]
```

**Result**:
- Window "dev": 2 panes (hx | lazygit)
- Window "server": 2 panes (api server | ui server)
- Window "logs": 1 pane (tail logs)
- Window "agent": 1 pane (shell)

### Example 4: Complex Layout

**Config**:
```
tmux_layout=[
  dev:hx|lazygit|
  server:bin/server
  db:psql myapp_dev
  agent:
]
```

**Result**:
- Window "dev": 3 panes (hx | lazygit | shell)
- Window "server": 1 pane (server)
- Window "db": 1 pane (psql)
- Window "agent": 1 pane (shell)

---

## Appendix: Logo ASCII Art

```
                       _
  _ __ ___  _   ___  _| |_ _ __ ___  ___
 | '_ ` _ \| | | \ \/ / __| '__/ _ \/ _ \
 | | | | | | |_| |>  <| |_| | |  __/  __/
 |_| |_| |_|\__,_/_/\_\__|_|  \___|\___|

  Tmux Worktree Session Manager v<VERSION>
```

**Formatting**: Logo in regular text, version line in dim

---

## Appendix: Command Exit Codes

| Command | Success | Error | User Cancel |
|---------|---------|-------|-------------|
| init | 0 | 1 | 0 |
| config | 0 | 1 | n/a |
| new | 0 | 1 | 1 |
| list | 0 | 1 | n/a |
| delete | 0 | 1 | 0 |
| sessions | 0 | 1 | n/a |
| help | 0 | n/a | n/a |
| version | 0 | n/a | n/a |

---

## Appendix: Required External Commands

| Command | Purpose | Availability Check |
|---------|---------|-------------------|
| git | All git operations | `git --version` |
| tmux | Session management | `tmux -V` |
| osascript | Terminal.app/iTerm2 | macOS only |
| open | Ghostty integration | macOS only |

**Recommendation**: Check for git and tmux on startup for commands that need them.

---

This specification provides complete implementation details for achieving 100% feature parity with muxtree v1.0.0. Each section can be worked on independently by following the documented behaviors, algorithms, and output formats.
