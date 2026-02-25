# Feature Spec: sessions open

## Test Case: Basic Open - Create Session for Existing Worktree

**Setup:**
- Repository: test-repo
- Existing worktree: ~/worktrees/test-repo/feature-auth (branch: feature-auth)
- No tmux session exists for this worktree
- Configuration: default (Terminal.app)

**Input:**
```
$ mxt sessions open feature-auth
```

**Expected Output:**
```
✓   Created session test-repo_feature-auth (windows: dev, agent)
```

**Expected Side Effects:**
- Tmux session `test-repo_feature-auth` created with 2 windows (dev, agent)
- Terminal.app window opened with tmux attached
- Dev window is selected (active)

**Exit Code:** 0

---

## Test Case: Open with --run Command

**Setup:**
- Repository: test-repo
- Existing worktree: ~/worktrees/test-repo/feature-api (branch: feature-api)
- No tmux session exists

**Input:**
```
$ mxt sessions open feature-api --run claude
```

**Expected Output:**
```
✓   Created session test-repo_feature-api (windows: dev, agent)
```

**Expected Side Effects:**
- Tmux session created
- Command "claude" sent to agent window (with Enter)
- Terminal.app window opened

**Exit Code:** 0

---

## Test Case: Open with --bg Flag

**Setup:**
- Repository: test-repo
- Existing worktree: ~/worktrees/test-repo/fix-bug (branch: fix-bug)
- No tmux session exists

**Input:**
```
$ mxt sessions open fix-bug --bg
```

**Expected Output:**
```
✓   Created session test-repo_fix-bug (windows: dev, agent)
```

**Expected Side Effects:**
- Tmux session created in detached mode
- NO terminal window opened

**Exit Code:** 0

---

## Test Case: Open with Custom Layout

**Setup:**
- Repository: test-repo
- Existing worktree: ~/worktrees/test-repo/feature-ui (branch: feature-ui)
- Configuration includes custom tmux_layout:
  ```
  tmux_layout=dev:hx|lazygit;server:bin/server;agent:
  ```

**Input:**
```
$ mxt sessions open feature-ui
```

**Expected Output:**
```
✓   Created session test-repo_feature-ui (windows: dev, server, agent)
```

**Expected Side Effects:**
- Tmux session created with custom windows (dev, server, agent)
- Dev window has 2 panes (hx | lazygit)
- Server window has 1 pane running "bin/server"
- Agent window has 1 empty pane
- Terminal.app window opened

**Exit Code:** 0

---

## Test Case: Error - Worktree Not Found

**Setup:**
- Repository: test-repo
- No worktree exists for branch "nonexistent"

**Input:**
```
$ mxt sessions open nonexistent
```

**Expected Output:**
```
✗ Worktree not found: ~/worktrees/test-repo/nonexistent
```

**Exit Code:** 1

---

## Test Case: Error - Session Already Exists

**Setup:**
- Repository: test-repo
- Existing worktree: ~/worktrees/test-repo/feature-auth (branch: feature-auth)
- Tmux session `test-repo_feature-auth` already exists

**Input:**
```
$ mxt sessions open feature-auth
```

**Expected Output:**
```
⚠ Session test-repo_feature-auth already exists
```

**Exit Code:** 0

---

## Test Case: Error - Not in Git Repository

**Setup:**
- Current directory: Not inside a git repository

**Input:**
```
$ mxt sessions open feature-auth
```

**Expected Output:**
```
✗ Not inside a git repository. Run muxtree from within your repo.
```

**Exit Code:** 1

---

## Test Case: Open with --run and --bg

**Setup:**
- Repository: test-repo
- Existing worktree: ~/worktrees/test-repo/feature-auth (branch: feature-auth)
- No tmux session exists

**Input:**
```
$ mxt sessions open feature-auth --run codex --bg
```

**Expected Output:**
```
✓   Created session test-repo_feature-auth (windows: dev, agent)
```

**Expected Side Effects:**
- Tmux session created
- Command "codex" sent to agent window
- NO terminal window opened (--bg flag)

**Exit Code:** 0

---

## Test Case: Error - Invalid --run Command

**Setup:**
- Repository: test-repo
- Existing worktree: ~/worktrees/test-repo/feature-auth (branch: feature-auth)

**Input:**
```
$ mxt sessions open feature-auth --run invalid
```

**Expected Output:**
```
✗ Invalid --run command: 'invalid' (must be 'claude' or 'codex')
```

**Exit Code:** 1

---

## Implementation Notes

### Command Behavior
1. Require git repository
2. Load configuration (global + project)
3. Determine repository name and worktree path
4. Validate worktree directory exists
5. Determine session name: `<repo-name>_<sanitized-branch>`
6. Check if session already exists (error if yes)
7. Create tmux session with appropriate layout (default or custom)
8. If `--run` provided, send command to agent window
9. If `--bg` NOT provided, open terminal window
10. Display success message

### Session Creation
- Use `CreateDefaultLayout` if no custom layout configured
- Use `CreateCustomLayout` if tmux_layout configured
- Session name format: `<repo>_<sanitized-branch>`
- Working directory: worktree path

### Terminal Integration
- Terminal type from config: `terminal` (default), `iterm2`, `ghostty`, `current`
- Only open terminal if `--bg` flag is NOT set
- Terminal integration details in separate spec

### Aliases
- `launch` and `start` are aliases for `open`
