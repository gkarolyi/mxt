# Feature Spec: sessions close

## Test Case: Basic Close - Kill Existing Session

**Setup:**
- Repository: test-repo
- Existing tmux session: `test-repo_feature-auth`
- Worktree exists at: ~/worktrees/test-repo/feature-auth

**Input:**
```
$ mxt sessions close feature-auth
```

**Expected Output:**
```
✓ Killed session test-repo_feature-auth
```

**Expected Side Effects:**
- Tmux session `test-repo_feature-auth` is terminated
- Any terminal windows attached to the session will close

**Exit Code:** 0

---

## Test Case: Close Non-Existent Session

**Setup:**
- Repository: test-repo
- No tmux session exists for branch "feature-api"
- Worktree may or may not exist (doesn't matter)

**Input:**
```
$ mxt sessions close feature-api
```

**Expected Output:**
```
✓ Killed session test-repo_feature-api
```

**Expected Side Effects:**
- No error (idempotent operation)
- Command succeeds even if session didn't exist

**Exit Code:** 0

---

## Test Case: Error - Not in Git Repository

**Setup:**
- Current directory: Not inside a git repository

**Input:**
```
$ mxt sessions close feature-auth
```

**Expected Output:**
```
✗ Not inside a git repository. Run mxt from within your repo.
```

**Exit Code:** 1

---

## Test Case: Alias - kill

**Setup:**
- Repository: test-repo
- Existing tmux session: `test-repo_feature-auth`

**Input:**
```
$ mxt sessions kill feature-auth
```

**Expected Output:**
```
✓ Killed session test-repo_feature-auth
```

**Expected Side Effects:**
- Session is killed (same as 'close')

**Exit Code:** 0

---

## Test Case: Alias - stop

**Setup:**
- Repository: test-repo
- Existing tmux session: `test-repo_feature-auth`

**Input:**
```
$ mxt sessions stop feature-auth
```

**Expected Output:**
```
✓ Killed session test-repo_feature-auth
```

**Expected Side Effects:**
- Session is killed (same as 'close')

**Exit Code:** 0

---

## Implementation Notes

### Command Behavior
1. Require git repository
2. Load configuration (for consistency)
3. Determine repository name
4. Generate session name: `<repo-name>_<sanitized-branch>`
5. Kill tmux session if it exists (idempotent)
6. Display success message

### Idempotent Operation
- No error if session doesn't exist
- Always returns success (exit code 0)
- Useful for cleanup scripts and automation

### Aliases
- `kill` is an alias for `close`
- `stop` is an alias for `close`
