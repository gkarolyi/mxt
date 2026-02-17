# muxtree — Tmux Worktree Session Manager

A lightweight CLI for spinning up isolated git worktrees paired with tmux sessions, purpose-built for running parallel [Claude Code](https://docs.anthropic.com/en/docs/claude-code) or [Codex](https://openai.com/index/codex/) sessions on macOS.

Each `muxtree new` call gives you a fresh branch in its own directory with your config files copied in and two terminal windows ready to go — one for viewing code and running your app, one for your AI coding agent.

---

## Install

```bash
# Copy the script somewhere on your PATH
cp muxtree /usr/local/bin/muxtree
chmod +x /usr/local/bin/muxtree

# Or with Homebrew's default bin path
cp muxtree ~/.local/bin/muxtree
```

### Prerequisites

- **git** (with worktree support — any modern version)
- **tmux** (`brew install tmux`)
- **macOS** with Terminal.app or iTerm2

---

## Quick Start

```bash
# 1. Run interactive setup (creates ~/.muxtree/config)
muxtree init

# 2. Navigate to your repo
cd ~/projects/my-app

# 3. Create a new worktree + sessions
muxtree new feature-auth

# 4. Two terminal windows open automatically:
#    • my-app_feature-auth-dev     ← run your app, browse code
#    • my-app_feature-auth-claude  ← run claude/codex here
```

That's it. You're working in an isolated branch with your `.env` and config files already copied over.

---

## Commands

### `muxtree init`

Interactive setup. Creates `~/.muxtree/config` where you specify:

- **Worktree base directory** — where all worktrees live (e.g. `~/worktrees`)
- **Terminal app** — `terminal` (Terminal.app) or `iterm2`
- **Files to copy** — comma-separated list of files to copy from your repo root into each new worktree (e.g. `.env,.env.local,CLAUDE.md`)

```bash
$ muxtree init
Worktree base directory [~/worktrees]: ~/worktrees
Terminal app (terminal/iterm2) [terminal]: iterm2
Files to copy: .env,.env.local,.claude/settings.json
✓ Config written to ~/.muxtree/config
```

### `muxtree new <branch> [options]`

Creates a worktree, copies config files, and launches two tmux sessions in new terminal windows.

```bash
# Branch from main (auto-detected)
muxtree new feature-auth

# Branch from a specific base
muxtree new fix-bug --from develop

# Auto-launch Claude Code in the claude session
muxtree new feature-ai --run claude

# Auto-launch Codex instead
muxtree new feature-ai --run codex
```

**What happens:**

1. `git worktree add -b <branch>` at `<worktree_dir>/<repo>/<branch>/`
2. Copies each file from `copy_files` config into the new worktree
3. Creates two detached tmux sessions
4. Opens each in a new terminal window

### `muxtree list`

Shows all managed worktrees with diff stats and session status.

```
Worktrees for my-app
════════════════════════════════════════════════════════════════

  feature-auth  +42 -7
  ~/worktrees/my-app/feature-auth
  Sessions: ● my-app_feature-auth-dev  ● my-app_feature-auth-claude

  fix-bug  +3 -1
  ~/worktrees/my-app/fix-bug
  Sessions: ○ my-app_fix-bug-dev  ○ my-app_fix-bug-claude
```

- `●` = tmux session is running
- `○` = tmux session is not running
- Diff stats show combined staged + unstaged changes vs HEAD

### `muxtree delete <branch> [--force]`

Removes a worktree, kills its tmux sessions, and deletes the local branch.

```bash
$ muxtree delete feature-auth

  Branch:    feature-auth
  Path:      ~/worktrees/my-app/feature-auth
  Changes:   +42 -7

⚠ This will remove the worktree and delete the local branch.
Are you sure? (y/N) y
✓ Killed session my-app_feature-auth-dev
✓ Killed session my-app_feature-auth-claude
✓ Worktree removed
✓ Branch deleted
```

Use `--force` or `-f` to skip confirmation.

### `muxtree sessions <action> <branch> [options]`

Manage tmux sessions independently of the worktree.

```bash
# Close both sessions for a branch
muxtree sessions close feature-auth

# Reopen them (creates new terminal windows)
muxtree sessions open feature-auth

# Reopen with claude auto-running
muxtree sessions open feature-auth --run claude

# Close + reopen in one step
muxtree sessions relaunch feature-auth --run codex

# Attach to a session in your current terminal
muxtree sessions attach feature-auth dev
muxtree sessions attach feature-auth claude
```

### `muxtree config`

Print the current config file.

### `muxtree help`

Show all commands and usage.

---

## Configuration

Config lives at `~/.muxtree/config` (override with `MUXTREE_CONFIG_DIR`). It's a plain key=value file:

```ini
# muxtree configuration

# Base directory for worktrees
worktree_dir=~/worktrees

# Terminal app: terminal | iterm2
terminal=iterm2

# Files to copy from repo root into new worktrees (comma-separated)
# Supports glob patterns
copy_files=.env,.env.local,CLAUDE.md,.claude/settings.json
```

### Config options

| Key | Default | Description |
|-----|---------|-------------|
| `worktree_dir` | `~/worktrees` | Base directory where worktrees are created. Organized as `<worktree_dir>/<repo>/<branch>/` |
| `terminal` | `terminal` | Which terminal app to open: `terminal` (Terminal.app) or `iterm2` |
| `copy_files` | *(empty)* | Comma-separated list of files/globs to copy from repo root into new worktrees |

### Glob patterns in copy_files

The `copy_files` value supports shell glob patterns:

```ini
# Copy specific files
copy_files=.env,.env.local

# Copy all dotenv files
copy_files=.env*

# Mix of specific files and patterns
copy_files=.env*,CLAUDE.md,config/*.local.json
```

---

## Directory Layout

```
~/worktrees/                  ← worktree_dir from config
  my-app/                     ← repo name (auto-detected)
    feature-auth/             ← branch name (sanitized)
      .env                    ← copied from repo root
      .env.local              ← copied from repo root
      src/                    ← full working tree
      ...
    fix-bug/
      ...
```

---

## Tmux Session Naming

Sessions follow the pattern `<repo>_<branch>-<type>`:

```
my-app_feature-auth-dev       ← for running your app / viewing code
my-app_feature-auth-claude    ← for Claude Code / Codex
```

Slashes and dots in branch names are replaced with dashes for tmux compatibility.

---

## Typical Workflow

```bash
# Start your day — create a fresh workspace
cd ~/projects/my-app
muxtree new feature-user-profiles --run claude

# Two terminal windows pop open:
#   Window 1 (dev):    cd'd into the worktree, run your dev server
#   Window 2 (claude): Claude Code is already running

# Check on all your active branches
muxtree list

# Done with a feature — clean up
muxtree delete feature-user-profiles

# Need to step away but keep the worktree? Just close sessions
muxtree sessions close feature-user-profiles

# Come back later and relaunch
muxtree sessions open feature-user-profiles --run claude
```

---

## Security

muxtree is designed with security in mind:

- **No shell execution of config** — config is parsed as plain key=value pairs, not sourced. Values containing shell metacharacters (`$`, `` ` ``, `;`, `|`, `&`) are rejected.
- **AppleScript injection prevention** — session names are escaped before embedding in osascript.
- **Branch name sanitization** — filesystem paths strip non-alphanumeric characters to prevent traversal.
- **Command validation** — `--run` only accepts `claude` or `codex`.
- **Safe file operations** — `--` separators on `rm`, `cp`, `mkdir` to handle edge-case filenames.

---

## Uninstall

```bash
# Remove the binary
rm /usr/local/bin/muxtree

# Remove config
rm -rf ~/.muxtree

# Optionally clean up any remaining worktrees
rm -rf ~/worktrees  # or wherever you configured them
```

---

## License

MIT
