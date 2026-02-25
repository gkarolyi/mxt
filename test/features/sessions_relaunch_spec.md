# Feature Spec: sessions relaunch

## Test Case: Basic Relaunch - Kill and Recreate Session

**Setup:**
- Repository: test-repo
- Existing worktree: ~/worktrees/test-repo/feature-auth
- Existing tmux session: `test-repo_feature-auth`

**Input:**
```
$ mxt sessions relaunch feature-auth
```

**Expected Output:**
```
✓ Killed session test-repo_feature-auth
▸ Creating tmux session...
✓ Created session test-repo_feature-auth (windows: dev, agent)
✓ Ready! Worktree: ~/worktrees/test-repo/feature-auth
```

**Expected Side Effects:**
- Old session is killed
- New session is created with fresh state
- Terminal window opened (default behavior)
- Dev window is selected

**Exit Code:** 0

---

## Test Case: Relaunch with --run Command

**Setup:**
- Repository: test-repo
- Existing worktree: ~/worktrees/test-repo/feature-api
- Existing tmux session: `test-repo_feature-api`

**Input:**
```
$ mxt sessions relaunch feature-api --run claude
```

**Expected Output:**
```
✓ Killed session test-repo_feature-api
▸ Creating tmux session...
✓ Created session test-repo_feature-api (windows: dev, agent)
✓ Ready! Worktree: ~/worktrees/test-repo/feature-api
```

**Expected Side Effects:**
- Old session killed
- New session created
- Command "claude" sent to agent window
- Terminal window opened

**Exit Code:** 0

---

## Test Case: Relaunch with --bg Flag

**Setup:**
- Repository: test-repo
- Existing worktree: ~/worktrees/test-repo/fix-bug
- Existing tmux session: `test-repo_fix-bug`

**Input:**
```
$ mxt sessions relaunch fix-bug --bg
```

**Expected Output:**
```
✓ Killed session test-repo_fix-bug
▸ Creating tmux session...
✓ Created session test-repo_fix-bug (windows: dev, agent)
✓ Ready! Worktree: ~/worktrees/test-repo/fix-bug
```

**Expected Side Effects:**
- Old session killed
- New session created in detached mode
- NO terminal window opened

**Exit Code:** 0

---

## Test Case: Relaunch When No Session Exists

**Setup:**
- Repository: test-repo
- Existing worktree: ~/worktrees/test-repo/feature-new
- No tmux session exists

**Input:**
```
$ mxt sessions relaunch feature-new
```

**Expected Output:**
```
✓ Killed session test-repo_feature-new
▸ Creating tmux session...
✓ Created session test-repo_feature-new (windows: dev, agent)
✓ Ready! Worktree: ~/worktrees/test-repo/feature-new
```

**Expected Side Effects:**
- Close operation succeeds (idempotent)
- New session created
- Behaves like 'open' command

**Exit Code:** 0

---

## Test Case: Error - Worktree Not Found

**Setup:**
- Repository: test-repo
- No worktree exists for branch "nonexistent"
- May or may not have an old session

**Input:**
```
$ mxt sessions relaunch nonexistent
```

**Expected Output:**
```
✓ Killed session test-repo_nonexistent
✗ Worktree not found: ~/worktrees/test-repo/nonexistent
```

**Exit Code:** 1

---

## Test Case: Error - Not in Git Repository

**Setup:**
- Current directory: Not inside a git repository

**Input:**
```
$ mxt sessions relaunch feature-auth
```

**Expected Output:**
```
✗ Not inside a git repository. Run mxt from within your repo.
```

**Exit Code:** 1

---

## Test Case: Alias - restart

**Setup:**
- Repository: test-repo
- Existing worktree and session for feature-auth

**Input:**
```
$ mxt sessions restart feature-auth
```

**Expected Output:**
```
✓ Killed session test-repo_feature-auth
▸ Creating tmux session...
✓ Created session test-repo_feature-auth (windows: dev, agent)
✓ Ready! Worktree: ~/worktrees/test-repo/feature-auth
```

**Expected Side Effects:**
- Same as 'relaunch'

**Exit Code:** 0

---

## Implementation Notes

### Command Behavior
Relaunch is a compound operation:
1. Execute `sessions close` for the branch
2. Execute `sessions open` for the branch with same flags

### Use Cases
- Refresh a session after config changes
- Reset a session to clean state
- Recover from corrupted session state
- Apply new tmux layout settings

### Aliases
- `restart` is an alias for `relaunch`

### Error Handling
- Close operation always succeeds (idempotent)
- If open fails (e.g., worktree not found), command fails
- Result: session is killed but not recreated
