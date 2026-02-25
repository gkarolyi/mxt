# Feature Spec: Default Tmux Layout

This document captures the exact tmux session structure created by `muxtree new` when no custom `tmux_layout` is configured.

## Test Case: Default layout structure

**Setup:**
- No `tmux_layout` configured in global or project config
- Inside a git repository named `test-repo`

**Input:**
```bash
muxtree new test-branch --bg
```

**Expected Output:**
```
▸ Creating worktree at /Users/username/Code/worktrees/test-repo/test-branch
Preparing worktree (new branch 'test-branch')
HEAD is now at abc1234 Initial commit
✓ Worktree created (branch test-branch from main)
▸ Creating tmux session...
✓   Created session test-repo_test-branch (windows: dev, agent)

✓ Ready! Worktree: /Users/username/Code/worktrees/test-repo/test-branch
```

**Exit Code:** 0

**Side Effects:**

Tmux session created with name: `test-repo_test-branch`

**Session verification:**
```bash
$ tmux has-session -t test-repo_test-branch
# Exit code: 0 (session exists)
```

**Window structure:**
```bash
$ tmux list-windows -t test-repo_test-branch
0: dev* (1 panes) [...]
1: agent- (1 panes) [...]
```

**Pane structure:**
```bash
$ tmux list-panes -t test-repo_test-branch:dev
0: [...]* (1 panes)

$ tmux list-panes -t test-repo_test-branch:agent
0: [...]* (1 panes)
```

**Working directory:**
- All panes start in the worktree directory: `/Users/username/Code/worktrees/test-repo/test-branch`

**Default window:**
- Window `dev` (index 0) is selected by default
- User sees the `dev` window when attaching to the session

**Session properties:**
- Session is created in detached mode (unless terminal is opened)
- Session name follows pattern: `<repo-name>_<sanitized-branch>`
- Both windows have exactly 1 pane each
- No commands are sent to panes initially

---

## Test Case: Default layout with --run claude

**Setup:**
- No `tmux_layout` configured
- Inside a git repository named `test-repo`

**Input:**
```bash
muxtree new test-branch --run claude --bg
```

**Expected Output:**
```
▸ Creating worktree at /Users/username/Code/worktrees/test-repo/test-branch
Preparing worktree (new branch 'test-branch')
HEAD is now at abc1234 Initial commit
✓ Worktree created (branch test-branch from main)
▸ Creating tmux session...
✓   Created session test-repo_test-branch (windows: dev, agent)

✓ Ready! Worktree: /Users/username/Code/worktrees/test-repo/test-branch
```

**Exit Code:** 0

**Side Effects:**

Tmux session created with command sent to agent window:
```bash
$ tmux capture-pane -t test-repo_test-branch:agent -p | head -1
$ claude
```

The command `claude` is sent to the agent window's pane using:
```bash
tmux send-keys -t test-repo_test-branch:agent 'claude' Enter
```

**Behavior:**
- Command is sent AFTER session is fully created
- Command is sent to window named "agent" (index 1)
- Command is sent to pane 0 of the agent window
- Enter key is sent after the command
- If "agent" window doesn't exist (shouldn't happen with default layout), command is not sent

---

## Test Case: Default layout with --run codex

**Setup:**
- No `tmux_layout` configured
- Inside a git repository named `test-repo`

**Input:**
```bash
muxtree new test-branch --run codex --bg
```

**Expected Output:**
```
▸ Creating worktree at /Users/username/Code/worktrees/test-repo/test-branch
Preparing worktree (new branch 'test-branch')
HEAD is now at abc1234 Initial commit
✓ Worktree created (branch test-branch from main)
▸ Creating tmux session...
✓   Created session test-repo_test-branch (windows: dev, agent)

✓ Ready! Worktree: /Users/username/Code/worktrees/test-repo/test-branch
```

**Exit Code:** 0

**Side Effects:**

Tmux session created with command sent to agent window:
```bash
$ tmux capture-pane -t test-repo_test-branch:agent -p | head -1
$ codex
```

---

## Test Case: Branch name sanitization for session name

**Setup:**
- No `tmux_layout` configured
- Inside a git repository named `test-repo`

**Input:**
```bash
muxtree new "feature/auth-api" --bg
```

**Expected Output:**
```
▸ Creating worktree at /Users/username/Code/worktrees/test-repo/feature-auth-api
Preparing worktree (new branch 'feature/auth-api')
HEAD is now at abc1234 Initial commit
✓ Worktree created (branch feature/auth-api from main)
▸ Creating tmux session...
✓   Created session test-repo_feature-auth-api (windows: dev, agent)

✓ Ready! Worktree: /Users/username/Code/worktrees/test-repo/feature-auth-api
```

**Exit Code:** 0

**Side Effects:**

Session name is sanitized:
- Branch: `feature/auth-api`
- Session: `test-repo_feature-auth-api` (slash replaced with dash)

**Sanitization rules:**
- Replace any character that is NOT alphanumeric, underscore, dash, or dot with dash: `[^a-zA-Z0-9._-]` → `-`
- Strip leading dash if present

**Examples:**
- `feature/auth` → `feature-auth`
- `bug-fix-#123` → `bug-fix--123`
- `user@domain` → `user-domain`

---

## Implementation Algorithm

**Default Layout Creation Steps:**

1. Determine session name: `<repo-name>_<sanitized-branch>`
2. Create new detached session:
   ```bash
   tmux new-session -d -s <session> -c <worktree-path>
   ```
3. Rename first window to "dev":
   ```bash
   tmux rename-window -t <session>:0 dev
   ```
4. Create second window named "agent":
   ```bash
   tmux new-window -t <session> -n agent -c <worktree-path>
   ```
5. If `--run` command provided:
   ```bash
   tmux send-keys -t <session>:agent '<command>' Enter
   ```
6. Select dev window (make it active):
   ```bash
   tmux select-window -t <session>:dev
   ```
7. Print success message:
   ```
   ✓   Created session <session> (windows: dev, agent)
   ```

**Success message format:**
- Green checkmark: `✓`
- Two spaces indentation
- Bold session name
- Window list: comma-separated, space after comma

---

## Color Formatting

- Success message: Green `✓`
- Session name in message: Bold
- Window list: Regular text

---

## Exit Codes

| Scenario | Exit Code |
|----------|-----------|
| Success | 0 |
| Tmux not installed | 1 |
| Tmux command fails | 1 |

---

## Configuration Values

**None used for default layout** - this is the fallback when `tmux_layout` is not configured.

---

## Notes for Implementation

1. **Session creation must be in detached mode** (`-d` flag) so it doesn't attach automatically
2. **Working directory** (`-c` flag) must be set to worktree path for all windows
3. **Window naming** uses `-n` flag when creating windows
4. **Window selection** (`tmux select-window`) ensures dev window is active when user attaches
5. **Error handling**: If any tmux command fails, show clear error and exit with code 1
