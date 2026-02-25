# Feature Spec: muxtree delete

This document captures the exact behavior of `muxtree delete` (alias: `muxtree rm`) for feature parity in the Go reimplementation.

## Test Case: Delete worktree with confirmation

**Setup:**
- Repository: `test-repo`
- Worktree exists at: `/Users/username/Code/worktrees/test-repo/feature-auth`
- Tmux session active: `test-repo_feature-auth`
- No staged or unstaged changes (0 insertions, 0 deletions)

**Input:**
```bash
cd /path/to/repo
muxtree delete feature-auth
# User input: y
```

**Expected Output:**
```

  Branch:    feature-auth
  Path:      /Users/username/Code/worktrees/test-repo/feature-auth
  Changes:   +0 -0

⚠ This will remove the worktree and delete the local branch.
Are you sure? (y/N) y
✓ Killed session test-repo_feature-auth
▸ Removing worktree...
✓ Worktree removed
▸ Deleting branch feature-auth...
Deleted branch feature-auth (was abc1234).
✓ Branch deleted

✓ Done.
```

**Exit Code:** 0

**Side Effects:**
- Tmux session `test-repo_feature-auth` is terminated
- Worktree directory is removed
- Local branch `feature-auth` is deleted
- Parent worktree directory removed if it becomes empty

---

## Test Case: Delete with --force (no prompt)

**Setup:**
- Same as previous test case

**Input:**
```bash
muxtree delete feature-auth --force
```

**Expected Output:**
```

  Branch:    feature-auth
  Path:      /Users/username/Code/worktrees/test-repo/feature-auth
  Changes:   +0 -0

✓ Killed session test-repo_feature-auth
▸ Removing worktree...
✓ Worktree removed
▸ Deleting branch feature-auth...
Deleted branch feature-auth (was abc1234).
✓ Branch deleted

✓ Done.
```

**Exit Code:** 0

**Behavior:**
- Skips confirmation prompt
- Proceeds directly with deletion

---

## Test Case: Cancel deletion

**Setup:**
- Same as previous test case

**Input:**
```bash
muxtree delete feature-auth
# User input: N
```

**Expected Output:**
```

  Branch:    feature-auth
  Path:      /Users/username/Code/worktrees/test-repo/feature-auth
  Changes:   +0 -0

⚠ This will remove the worktree and delete the local branch.
Are you sure? (y/N) N
▸ Cancelled.
```

**Exit Code:** 0

**Side Effects:**
- Worktree remains intact
- Local branch is not deleted
- Tmux session stays active

---

## Test Case: Error - Worktree not found

**Setup:**
- Repository: `test-repo`
- No worktree exists at the expected path

**Input:**
```bash
muxtree delete missing-branch
```

**Expected Output:**
```
✗ Worktree not found: /Users/username/Code/worktrees/test-repo/missing-branch
```

**Exit Code:** 1

---

## Test Case: Error - Not in git repository

**Setup:**
- Current directory is not inside a git repository

**Input:**
```bash
cd /tmp
muxtree delete feature-auth
```

**Expected Output:**
```
✗ Not inside a git repository. Run muxtree from within your repo.
```

**Exit Code:** 1

---

## Test Case: Alias - rm

**Setup:**
- Same as delete with --force

**Input:**
```bash
muxtree rm feature-auth --force
```

**Expected Output:**
```
(Same output as "Delete with --force")
```

**Exit Code:** 0

**Behavior:**
- `rm` is an exact alias for `delete`
