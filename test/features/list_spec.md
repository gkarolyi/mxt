# Feature Spec: muxtree list

This document captures the exact behavior of `muxtree list` (alias: `muxtree ls`) command for feature parity in the Go reimplementation.

## Test Case: List with multiple worktrees

**Setup:**
- Repository name: `test-repo`
- Two worktrees exist:
  - `feature-auth` (3 insertions, 2 deletions, active session)
  - `fix-bug` (0 changes, no session)

**Input:**
```bash
cd /path/to/repo
muxtree list
```

**Expected Output:**
```
Worktrees for test-repo
════════════════════════════════════════════════════════════════

  feature-auth  +3 -2
  /Users/username/Code/worktrees/test-repo/feature-auth
  Session: ● test-repo_feature-auth

  fix-bug  +0 -0
  /Users/username/Code/worktrees/test-repo/fix-bug
  Session: ○ test-repo_fix-bug
```

**Exit Code:** 0

**Behavior:**
- Header shows repository name
- Double-line separator (64 equals signs)
- Each worktree has blank line before it
- Branch name: Bold, cyan
- Change stats: Green `+<num>`, red `-<num>`
- Path: Dim, second line
- Session line: "Session: " + status indicator + session name
- Active session: Green `●`
- Inactive session: Dim `○`

---

## Test Case: List with no worktrees

**Setup:**
- Repository name: `test-repo`
- No worktrees exist in `$WORKTREE_DIR/test-repo/`
- OR: `$WORKTREE_DIR/test-repo/` directory doesn't exist

**Input:**
```bash
cd /path/to/repo
muxtree list
```

**Expected Output:**
```
Worktrees for test-repo
════════════════════════════════════════════════════════════════
▸ No worktrees found. Use muxtree new <branch> to create one.
```

**Exit Code:** 0

**Behavior:**
- Header still displayed
- Info message (blue `▸`) instead of worktree list
- Helpful hint about creating worktrees

---

## Test Case: List with only non-managed worktrees

**Setup:**
- Repository has worktrees, but none are in `$WORKTREE_DIR/<repo>/`
- Example: Manual worktrees in other locations

**Input:**
```bash
cd /path/to/repo
muxtree list
```

**Expected Output:**
```
Worktrees for test-repo
════════════════════════════════════════════════════════════════
▸ No managed worktrees found. Use muxtree new <branch> to create one.
```

**Exit Code:** 0

**Behavior:**
- Only worktrees under `$WORKTREE_DIR/<repo>/` are shown
- Other worktrees are ignored

---

## Test Case: List using alias

**Setup:**
- Same as first test case

**Input:**
```bash
muxtree ls
```

**Expected Output:**
```
(Same output as muxtree list)
```

**Exit Code:** 0

**Behavior:**
- `ls` is an exact alias for `list`
- No difference in output

---

## Test Case: Not in git repository

**Input:**
```bash
cd /tmp
muxtree list
```

**Expected Output:**
```
✗ Not inside a git repository. Run muxtree from within your repo.
```

**Exit Code:** 1

**Behavior:**
- First validation check
- Clear, actionable error message

---

## Test Case: Branch name sanitization in session names

**Setup:**
- Worktree with branch name: `feature/auth-api`
- Worktree path: `$WORKTREE_DIR/test-repo/feature-auth-api`
- Active session: `test-repo_feature-auth-api`

**Input:**
```bash
muxtree list
```

**Expected Output:**
```
Worktrees for test-repo
════════════════════════════════════════════════════════════════

  feature/auth-api  +5 -1
  /Users/username/Code/worktrees/test-repo/feature-auth-api
  Session: ● test-repo_feature-auth-api
```

**Exit Code:** 0

**Behavior:**
- Branch name displayed as-is (with slash)
- Path uses sanitized name (feature-auth-api)
- Session name uses sanitized name

---

## Test Case: Worktree with only staged changes

**Setup:**
- Worktree `feature-auth` has 5 insertions, 2 deletions staged
- No unstaged changes

**Input:**
```bash
muxtree list
```

**Expected Output:**
```
Worktrees for test-repo
════════════════════════════════════════════════════════════════

  feature-auth  +5 -2
  /Users/username/Code/worktrees/test-repo/feature-auth
  Session: ○ test-repo_feature-auth
```

**Exit Code:** 0

**Behavior:**
- Staged and unstaged changes are summed
- Both `git diff --stat` and `git diff --cached --stat` are checked

---

## Test Case: Worktree with both staged and unstaged changes

**Setup:**
- Worktree `feature-auth` has:
  - Staged: 3 insertions, 1 deletion
  - Unstaged: 2 insertions, 1 deletion

**Input:**
```bash
muxtree list
```

**Expected Output:**
```
Worktrees for test-repo
════════════════════════════════════════════════════════════════

  feature-auth  +5 -2
  /Users/username/Code/worktrees/test-repo/feature-auth
  Session: ○ test-repo_feature-auth
```

**Exit Code:** 0

**Behavior:**
- Insertions: 3 + 2 = 5
- Deletions: 1 + 1 = 2

---

## Change Statistics Algorithm

**For each worktree:**

1. Run `git -C <worktree-path> diff --stat HEAD`
   - Capture stdout
   - Parse last line for: `<n> insertion`, `<m> deletion`
   - Extract numbers (default to 0 if not found)

2. Run `git -C <worktree-path> diff --cached --stat HEAD`
   - Capture stdout
   - Parse last line for: `<n> insertion`, `<m> deletion`
   - Extract numbers (default to 0 if not found)

3. Sum the results:
   - `total_insertions = unstaged_insertions + staged_insertions`
   - `total_deletions = unstaged_deletions + staged_deletions`

4. Display: `+<total_insertions> -<total_deletions>`

**Example git diff output:**
```
 file1.go | 5 +++--
 file2.go | 3 ++-
 2 files changed, 5 insertions(+), 3 deletions(-)
```

**Parsing:**
- Last line: `2 files changed, 5 insertions(+), 3 deletions(-)`
- Extract: insertions=5, deletions=3

---

## Session Status Algorithm

**For each worktree:**

1. Determine session name: `<repo-name>_<sanitized-branch>`
2. Run `tmux has-session -t <session-name>`
   - Exit code 0: Session exists (active)
   - Exit code 1: Session does not exist (inactive)
3. Display:
   - Active: Green `●` followed by session name
   - Inactive: Dim `○` followed by session name

---

## Worktree Parsing Algorithm

1. Run `git worktree list --porcelain`
2. Parse output:
   - Each worktree starts with `worktree <path>`
   - Branch line: `branch refs/heads/<branch-name>`
   - Skip worktrees without branch (detached HEAD)
3. Filter worktrees:
   - Only include if path starts with `$WORKTREE_DIR/<repo-name>/`
   - Skip repo root worktree
4. For each managed worktree:
   - Extract branch name from path (last component)
   - Calculate change statistics
   - Check session status
   - Display worktree info

**Example `git worktree list --porcelain` output:**
```
worktree /Users/username/Code/test-repo
HEAD abc123def456
branch refs/heads/main

worktree /Users/username/Code/worktrees/test-repo/feature-auth
HEAD def456abc123
branch refs/heads/feature-auth

worktree /Users/username/Code/worktrees/test-repo/fix-bug
HEAD 789abc012def
branch refs/heads/fix-bug
```

---

## Color Formatting

- Header: Regular text
- Separator: 64 equals signs (═)
- Branch name: Bold, cyan
- Change stats:
  - `+<num>`: Green
  - `-<num>`: Red
- Path: Dim
- Session status:
  - Active `●`: Green
  - Inactive `○`: Dim
  - Session name: Regular text
- Info message (`▸`): Blue

---

## Configuration Values Used

From config files:
- `worktree_dir`: Base directory for worktrees (determines which worktrees are "managed")

---

## Implementation Notes

**Sorting:**
- Worktrees should be listed in the order returned by `git worktree list`
- No specific sorting required

**Error handling:**
- If `git worktree list` fails: Show error and exit
- If `git diff` fails for a worktree: Show 0 changes (don't fail)
- If `tmux has-session` fails: Assume session doesn't exist

**Performance:**
- Run git diff commands in parallel if possible
- Cache results to avoid redundant git operations

**Edge cases:**
- Empty repository (no commits): Handle gracefully
- Worktree with uncommitted changes: Show stats correctly
- Worktree with untracked files: Untracked files don't affect stats

---

## Exit Code Summary

| Scenario | Exit Code |
|----------|-----------|
| Success (with worktrees) | 0 |
| Success (no worktrees) | 0 |
| Not in git repo | 1 |
| Git command fails | 1 |

---

## Examples

### Empty repository
```
$ muxtree list
Worktrees for new-repo
════════════════════════════════════════════════════════════════
▸ No worktrees found. Use muxtree new <branch> to create one.
```

### Single worktree with changes
```
$ muxtree list
Worktrees for my-app
════════════════════════════════════════════════════════════════

  feature-login  +127 -43
  /Users/username/Code/worktrees/my-app/feature-login
  Session: ● my-app_feature-login
```

### Multiple worktrees, mixed session status
```
$ muxtree list
Worktrees for my-app
════════════════════════════════════════════════════════════════

  feature-login  +127 -43
  /Users/username/Code/worktrees/my-app/feature-login
  Session: ● my-app_feature-login

  feature-auth  +53 -12
  /Users/username/Code/worktrees/my-app/feature-auth
  Session: ○ my-app_feature-auth

  fix-typo  +1 -1
  /Users/username/Code/worktrees/my-app/fix-typo
  Session: ○ my-app_fix-typo
```
