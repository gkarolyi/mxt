# Feature Spec: muxtree help

This document captures the exact behavior of `muxtree help` command for feature parity in the Go reimplementation.

## Test Case: Help output

**Input:**
```bash
muxtree help
```

**Expected Output:**
```
                       _
  _ __ ___  _   ___  _| |_ _ __ ___  ___
 | '_ ` _ \| | | \ \/ / __| '__/ _ \/ _ \
 | | | | | | |_| |>  <| |_| | |  __/  __/
 |_| |_| |_|\__,_/_/\_\\__|_|  \___|\___|
  Tmux Worktree Session Manager v1.0.0

USAGE
    muxtree <command> [options]

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
    muxtree init                          # Global setup
    muxtree init --local                  # Project-specific copy_files
    muxtree new feature-auth              # New worktree from main
    muxtree new fix-bug --from develop    # New worktree from develop
    muxtree new feature-ai --run claude   # Auto-launch claude code
    muxtree new fix-bug --bg              # Create without opening terminals
    muxtree list                          # Show all worktrees + status
    muxtree sessions close feature-auth   # Kill tmux sessions
    muxtree sessions relaunch fix-bug     # Restart sessions
    muxtree delete feature-auth           # Remove worktree + branch

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

**Exit Code:** 0

**Behavior:**
- Displays ASCII logo followed by help text
- Output is identical for `muxtree help`, `muxtree -h`, and `muxtree --help`

---

## Test Case: Alias - -h

**Input:**
```bash
muxtree -h
```

**Expected Output:**
```
(Same output as "Help output")
```

**Exit Code:** 0

---

## Test Case: Alias - --help

**Input:**
```bash
muxtree --help
```

**Expected Output:**
```
(Same output as "Help output")
```

**Exit Code:** 0
