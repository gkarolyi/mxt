```
███╗   ███╗██╗  ██╗████████╗
████╗ ████║╚██╗██╔╝╚══██╔══╝
██╔████╔██║ ╚███╔╝    ██║
██║╚██╔╝██║ ██╔██╗    ██║
██║ ╚═╝ ██║██╔╝ ██╗   ██║
╚═╝     ╚═╝╚═╝  ╚═╝   ╚═╝
```

**Tmux Worktree Session Manager**

A lightweight CLI for spinning up isolated git worktrees paired with tmux sessions, purpose-built for running parallel [Claude Code](https://code.claude.com/docs) or [Codex](https://openai.com/codex/) sessions on macOS.

Each `mxt new` call gives you a fresh branch in its own directory with your config files copied in and a tmux session with two windows ready to go — one for viewing code and running your app, one for your AI coding agent. Switch between them with `Ctrl-b n` / `Ctrl-b p`.

## Acknowledgments

This project is an adaptation of the original [muxtree](https://github.com/b-d055/muxtree) tool. Thanks to the muxtree project for the inspiration and foundation; this version is tailored to my own workflows and preferences.

---

## Install

```bash
# Install latest version
go install github.com/gkarolyi/mxt/cmd/mxt@latest

# Or build from source
go build -o mxt ./cmd/mxt

# Copy it somewhere on your PATH
cp mxt /usr/local/bin/mxt

# Or with Homebrew's default bin path
cp mxt ~/.local/bin/mxt
```

Ensure `$(go env GOPATH)/bin` is on your PATH when using `go install`.

### Prerequisites

- **git** (with worktree support — any modern version)
- **tmux** (`brew install tmux`)
- **Go** (1.24+; required to build mxt from source)
- **macOS** with Terminal.app or iTerm2

### Shell Completion

Tab completion is available for bash and zsh, providing completion for commands, flags, session actions, and managed branch names.

**Bash** — requires [`bash-completion`](https://github.com/scop/bash-completion) (`brew install bash-completion@2`). Add to `~/.bashrc` or `~/.bash_profile`:

```bash
source /path/to/mxt/completions/mxt.bash
```

**Zsh** — add to `~/.zshrc`:

```zsh
source /path/to/mxt/completions/mxt.zsh
```

Replace `/path/to/mxt` with the actual path to your mxt checkout or install location.

---

## Quick Start

```bash
# 1. Run interactive setup (creates ~/.mxt/config)
mxt init

# 2. Navigate to your repo
cd ~/projects/my-app

# 3. Create a new worktree + tmux session (prompts for branch name)
mxt new
# Branch name: feature-auth

# 4. A terminal window opens with a tmux session containing two windows:
#    • dev    ← run your app, browse code
#    • agent  ← run claude/codex here
#    Switch windows with Ctrl-b n / Ctrl-b p
```

That's it. You're working in an isolated branch with your `.env` and config files already copied over.

---

## How It Works

```
  Your repo (~/projects/my-app)
  ├── main branch (your normal working copy)
  │
  │  mxt new feature-auth --run claude
  │  ┌──────────────────────────────────────────────────────────┐
  │  │  1. git worktree add                                     │
  │  │     ~/worktrees/my-app/feature-auth/  (branch: feature-auth)
  │  │                                                          │
  │  │  2. Copy config files (.env, CLAUDE.md, etc.)            │
  │  │                                                          │
  │  │  3. Create tmux session: my-app_feature-auth             │
  │  │     ┌─────────────┐  ┌─────────────┐                    │
  │  │     │  dev window  │  │ agent window│                    │
  │  │     │  (run app,   │  │ (claude is  │                    │
  │  │     │   view code) │  │  running)   │                    │
  │  │     └─────────────┘  └─────────────┘                    │
  │  │     Ctrl-b n / Ctrl-b p to switch                        │
  │  │                                                          │
  │  │  4. Open terminal window attached to session             │
  │  └──────────────────────────────────────────────────────────┘
  │
  │  mxt new fix-bug
  │  └─► ~/worktrees/my-app/fix-bug/  →  tmux: my-app_fix-bug
  │
  │  mxt list         ← see all worktrees + diff stats + session status
  │  mxt delete fix-bug  ← kills session, removes worktree + branch
```

Each worktree is a fully independent working directory — separate branch, separate files, separate tmux session. You can run multiple AI agents in parallel without them stepping on each other.

---

## Commands

### `mxt init`

Interactive setup. Creates `~/.mxt/config` (TOML) where you specify:

- **Worktree base directory** — where all worktrees live (e.g. `~/worktrees`)
- **Terminal app** — `terminal` (Terminal.app), `iterm2`, `ghostty`, or `current`
- **Files to copy** — comma-separated list of files to copy from your repo root into each new worktree (e.g. `.env,.env.local,CLAUDE.md`)

```bash
$ mxt init
Worktree base directory [~/worktrees]: ~/worktrees
Terminal app (terminal/iterm2) [terminal]: iterm2
Files to copy: .env,.env.local,.claude/settings.json
✓ Config written to ~/.mxt/config
```

### `mxt new [branch] [options]`

Creates a worktree, copies config files, and launches a tmux session with two windows (dev + agent) in a new terminal.
If the branch name is omitted, mxt prompts for it when stdin is a TTY.

```bash
# Prompt for branch name when omitted
mxt new

# Branch from main (auto-detected)
mxt new feature-auth

# Branch from a specific base
mxt new fix-bug --from develop

# Auto-launch Claude Code in the claude session
mxt new feature-ai --run claude

# Auto-launch Codex instead
mxt new feature-ai --run codex

# Create worktree + session without opening a terminal window
mxt new fix-bug --bg
```

**What happens:**

1. `git worktree add -b <branch>` at `<worktree_dir>/<repo>/<branch>/`
2. Copies each file from `copy_files` config into the new worktree
3. Creates a detached tmux session with two windows (dev + agent)
4. Opens the session in a new terminal window

### `mxt list`

Shows all managed worktrees with diff stats and session status.

```
Worktrees for my-app
════════════════════════════════════════════════════════════════

  feature-auth  +42 -7
  ~/worktrees/my-app/feature-auth
  Session: ● my-app_feature-auth

  fix-bug  +3 -1
  ~/worktrees/my-app/fix-bug
  Session: ○ my-app_fix-bug
```

- `●` = tmux session is running
- `○` = tmux session is not running
- Diff stats show combined staged + unstaged changes vs HEAD

### `mxt delete <branch> [--force]`

Removes a worktree, kills its tmux sessions, and deletes the local branch.

```bash
$ mxt delete feature-auth

  Branch:    feature-auth
  Path:      ~/worktrees/my-app/feature-auth
  Changes:   +42 -7

⚠ This will remove the worktree and delete the local branch.
Are you sure? (y/N) y
✓ Killed session my-app_feature-auth
✓ Worktree removed
✓ Branch deleted
```

Use `--force` or `-f` to skip confirmation.

### `mxt sessions <action> <branch> [options]`

Manage the tmux session independently of the worktree.
If the branch is omitted and stdin/stdout are TTY, mxt opens an interactive selector via fzf (install fzf or pass a branch name).

```bash
# Close session for a branch
mxt sessions close feature-auth

# Reopen it (creates new terminal window)
mxt sessions open feature-auth

# Reopen with claude auto-running in the agent window
mxt sessions open feature-auth --run claude

# Close + reopen in one step
mxt sessions relaunch feature-auth --run codex

# Attach to session in your current terminal
mxt sessions attach feature-auth

# Attach with a specific window selected
mxt sessions attach feature-auth agent
```

### `mxt config`

Shows both global (`~/.mxt/config`) and project-local (`.mxt`) config files, labeling which one is active. Useful for debugging which settings are in effect. Use `mxt config migrate` to convert legacy key=value configs to TOML.

### `mxt config migrate`

Converts legacy key=value config files to TOML in place (global and project). Safe to run multiple times.

### `mxt version`

Print the version number. Also available as `mxt -v` or `mxt --version`.

### `mxt help`

Show all commands and usage. Also available as `mxt -h` or `mxt --help`.

---

## Configuration

Config lives at `~/.mxt/config` (override with `MXT_CONFIG_DIR`). It's a TOML file:

```toml
# mxt configuration (TOML)

# Base directory for worktrees
worktree_dir = "~/worktrees"

# Terminal app: terminal | iterm2 | ghostty | current
terminal = "iterm2"

# Files to copy from repo root into new worktrees (comma-separated)
# Supports glob patterns
copy_files = ".env,.env.local,CLAUDE.md,.claude/settings.json"
```

Legacy key=value configs can be converted with `mxt config migrate`.

### Config options

| Key | Default | Description |
|-----|---------|-------------|
| `worktree_dir` | `~/worktrees` | Base directory where worktrees are created. Organized as `<worktree_dir>/<repo>/<branch>/` |
| `terminal` | `terminal` | Which terminal app to open: `terminal` (Terminal.app), `iterm2`, `ghostty`, or `current` |
| `copy_files` | *(empty)* | Comma-separated list of files/globs to copy from repo root into new worktrees |

### Project-local config

You can create a `.mxt` file in your repo root to override global settings on a per-project basis. This is useful for setting project-specific `copy_files`.

```bash
# Interactive setup for the current repo
mxt init --local
```

The local config file uses the same TOML format. When present, local values override the global config.

### Glob patterns in copy_files

The `copy_files` value supports shell glob patterns:

```toml
# Copy specific files
copy_files = ".env,.env.local"

# Copy all dotenv files
copy_files = ".env*"

# Mix of specific files and patterns
copy_files = ".env*,CLAUDE.md,config/*.local.json"
```

### Command aliases

| Command | Aliases |
|---------|---------|
| `mxt list` | `mxt ls` |
| `mxt delete` | `mxt rm` |
| `mxt sessions` | `mxt s` |
| `mxt help` | `mxt -h`, `mxt --help` |
| `sessions open` | `sessions launch`, `sessions start` |
| `sessions close` | `sessions kill`, `sessions stop` |
| `sessions relaunch` | `sessions restart` |

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

Each worktree gets a single tmux session named `<repo>_<branch>` with two windows:

```
my-app_feature-auth
  ├── dev     ← for running your app / viewing code
  └── agent   ← for Claude Code / Codex
```

Switch between windows with `Ctrl-b n` (next) and `Ctrl-b p` (previous).

Branch names are sanitized in two ways:

- **Session names**: any character that isn't alphanumeric, underscore, or dash is replaced with a dash (tmux compatibility).
- **Filesystem paths**: any character that isn't alphanumeric, dot, underscore, or dash is replaced with a dash (traversal prevention).

---

## Typical Workflow

```bash
# Start your day — create a fresh workspace
cd ~/projects/my-app
mxt new feature-user-profiles --run claude

# A terminal window opens with a tmux session:
#   dev window:    cd'd into the worktree, run your dev server
#   agent window:  Claude Code is already running
# Switch between them with Ctrl-b n / Ctrl-b p

# Check on all your active branches
mxt list

# Done with a feature — clean up
mxt delete feature-user-profiles

# Need to step away but keep the worktree? Just close the session
mxt sessions close feature-user-profiles

# Come back later and relaunch
mxt sessions open feature-user-profiles --run claude
```

---

## Security

mxt is designed with security in mind:

- **No shell execution of config** — config is parsed as plain key=value pairs, not sourced. Values containing shell metacharacters (`$`, `` ` ``, `;`, `|`, `&`) are ignored with a warning.
- **AppleScript injection prevention** — session names are escaped before embedding in osascript.
- **Branch name sanitization** — filesystem paths strip non-alphanumeric characters to prevent traversal.
- **Command validation** — `--run` only accepts `claude` or `codex`.
- **Safe file operations** — `--` separators on `rm`, `cp`, `mkdir` to handle edge-case filenames.

---

## Uninstall

```bash
# Remove the binary
rm /usr/local/bin/mxt

# Remove config
rm -rf ~/.mxt

# Optionally clean up any remaining worktrees
rm -rf ~/worktrees  # or wherever you configured them
```

---

## License

MIT
