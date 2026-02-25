# Feature Spec: muxtree new

This document captures the exact behavior of `muxtree new` command for feature parity in the Go reimplementation.

**Phase 4 Implementation Scope:** Worktree creation, file copying, and pre-session command execution. Tmux session creation will be added in Phase 5.

## Test Case: Basic worktree creation with --bg flag

**Input:**
```bash
cd /path/to/repo
muxtree new test-feature --bg
```

**Expected Output:**
```
▸ Creating worktree at /Users/username/Code/worktrees/repo-name/test-feature
Preparing worktree (new branch 'test-feature')
HEAD is now at abc1234 Initial commit
✓ Worktree created (branch test-feature from master)
▸ Copying config files...
▸ Creating tmux session...
✓   Created session repo-name_test-feature (windows: dev, agent)

✓ Ready! Worktree: /Users/username/Code/worktrees/repo-name/test-feature
```

**Exit Code:** 0

**Side Effects:**
- New git worktree created at `$WORKTREE_DIR/<repo>/<branch>/`
- New branch `test-feature` created from default base branch (master/main)
- Tmux session created (Phase 5) but NOT opened (--bg flag)

**Phase 4 Note:** Implement everything except "Creating tmux session" and session creation

---

## Test Case: Worktree with --from flag

**Input:**
```bash
cd /path/to/repo
git checkout -b develop
muxtree new feature-x --from develop --bg
```

**Expected Output:**
```
▸ Creating worktree at /Users/username/Code/worktrees/repo-name/feature-x
Preparing worktree (new branch 'feature-x')
HEAD is now at abc1234 Initial commit
✓ Worktree created (branch feature-x from develop)
▸ Copying config files...
▸ Creating tmux session...
✓   Created session repo-name_feature-x (windows: dev, agent)

✓ Ready! Worktree: /Users/username/Code/worktrees/repo-name/feature-x
```

**Exit Code:** 0

**Side Effects:**
- New worktree created from `develop` branch instead of default main

---

## Test Case: With file copying (single file)

**Setup:**
- Project config `.muxtree` contains: `copy_files=README.md`
- File `README.md` exists in repo root

**Input:**
```bash
muxtree new feature-with-files --bg
```

**Expected Output:**
```
▸ Creating worktree at /Users/username/Code/worktrees/repo-name/feature-with-files
Preparing worktree (new branch 'feature-with-files')
HEAD is now at abc1234 Initial commit
✓ Worktree created (branch feature-with-files from master)
▸ Copying config files...
✓   Copied README.md
▸ Creating tmux session...
✓   Created session repo-name_feature-with-files (windows: dev, agent)

✓ Ready! Worktree: /Users/username/Code/worktrees/repo-name/feature-with-files
```

**Exit Code:** 0

**Side Effects:**
- `README.md` file copied from repo root to worktree root
- File permissions/attributes preserved

---

## Test Case: With file copying (glob pattern)

**Setup:**
- Project config `.muxtree` contains: `copy_files=*.md`
- Files `README.md` and `test.md` exist in repo root

**Input:**
```bash
muxtree new feature-glob --bg
```

**Expected Output:**
```
▸ Creating worktree at /Users/username/Code/worktrees/repo-name/feature-glob
Preparing worktree (new branch 'feature-glob')
HEAD is now at abc1234 Initial commit
✓ Worktree created (branch feature-glob from master)
▸ Copying config files...
✓   Copied README.md
✓   Copied test.md
▸ Creating tmux session...
✓   Created session repo-name_feature-glob (windows: dev, agent)

✓ Ready! Worktree: /Users/username/Code/worktrees/repo-name/feature-glob
```

**Exit Code:** 0

**Side Effects:**
- All files matching `*.md` pattern copied to worktree
- Files listed in alphabetical order

---

## Test Case: With file copying (missing file)

**Setup:**
- Project config `.muxtree` contains: `copy_files=missing.txt,README.md`
- File `missing.txt` does NOT exist
- File `README.md` exists

**Input:**
```bash
muxtree new feature-missing --bg
```

**Expected Output:**
```
▸ Creating worktree at /Users/username/Code/worktrees/repo-name/feature-missing
Preparing worktree (new branch 'feature-missing')
HEAD is now at abc1234 Initial commit
✓ Worktree created (branch feature-missing from master)
▸ Copying config files...
⚠   Not found: missing.txt
✓   Copied README.md
▸ Creating tmux session...
✓   Created session repo-name_feature-missing (windows: dev, agent)

✓ Ready! Worktree: /Users/username/Code/worktrees/repo-name/feature-missing
```

**Exit Code:** 0

**Behavior:**
- Missing files generate warning but don't abort
- Other files are still copied
- Command completes successfully

---

## Test Case: With pre-session command (success)

**Setup:**
- Project config `.muxtree` contains: `pre_session_cmd=echo 'Setup complete'`

**Input:**
```bash
muxtree new feature-cmd --bg
```

**Expected Output:**
```
▸ Creating worktree at /Users/username/Code/worktrees/repo-name/feature-cmd
Preparing worktree (new branch 'feature-cmd')
HEAD is now at abc1234 Initial commit
✓ Worktree created (branch feature-cmd from master)
▸ Copying config files...
▸ Running pre-session command...
  echo 'Setup complete'
Setup complete
✓ Pre-session command completed
▸ Creating tmux session...
✓   Created session repo-name_feature-cmd (windows: dev, agent)

✓ Ready! Worktree: /Users/username/Code/worktrees/repo-name/feature-cmd
```

**Exit Code:** 0

**Behavior:**
- Command echoed (indented, dim)
- Command output shown
- Success message displayed
- Session creation continues

---

## Test Case: With pre-session command (failure)

**Setup:**
- Project config `.muxtree` contains: `pre_session_cmd=exit 1`

**Input:**
```bash
muxtree new feature-fail --bg
# User input: N (cancel)
```

**Expected Output:**
```
▸ Creating worktree at /Users/username/Code/worktrees/repo-name/feature-fail
Preparing worktree (new branch 'feature-fail')
HEAD is now at abc1234 Initial commit
✓ Worktree created (branch feature-fail from master)
▸ Copying config files...
▸ Running pre-session command...
  exit 1
⚠ Pre-session command failed (exit code: 1)
Continue anyway? (y/N) N
✗ Aborted due to pre-session command failure
```

**Exit Code:** 1

**Behavior:**
- Command fails
- User prompted for confirmation
- If user enters 'N' or anything other than 'y'/'Y', abort
- Error message displayed
- Exit with code 1
- Worktree remains (not rolled back)

---

## Test Case: Branch already exists

**Setup:**
- Branch `existing-branch` already exists in the repository

**Input:**
```bash
muxtree new existing-branch
```

**Expected Output:**
```
✗ Branch 'existing-branch' already exists. Use a different name, or delete it first.
```

**Exit Code:** 1

**Behavior:**
- Check performed BEFORE creating worktree
- Clear error message
- Suggests alternatives

---

## Test Case: Worktree path already exists

**Setup:**
- Directory `$WORKTREE_DIR/repo-name/branch-name` already exists

**Input:**
```bash
muxtree new branch-name
```

**Expected Output:**
```
✗ Worktree already exists at /Users/username/Code/worktrees/repo-name/branch-name
```

**Exit Code:** 1

**Behavior:**
- Check performed BEFORE creating branch
- Prevents overwriting existing directories
- Shows full path in error

---

## Test Case: Base branch doesn't exist

**Input:**
```bash
muxtree new new-feature --from nonexistent-branch
```

**Expected Output:**
```
✗ Base branch 'nonexistent-branch' does not exist.
```

**Exit Code:** 1

**Behavior:**
- Validation performed before worktree creation
- Clear error message with branch name

---

## Test Case: Not in git repository

**Input:**
```bash
cd /tmp
muxtree new test-branch
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

## Test Case: Invalid --run command

**Input:**
```bash
muxtree new feature --run invalid-command
```

**Expected Output:**
```
✗ Invalid --run command: 'invalid-command'. Use 'claude' or 'codex'.
```

**Exit Code:** 1

**Behavior:**
- Validation performed early
- Shows valid options
- **Phase 5 Note:** --run flag handling is part of tmux integration

---

## Test Case: Branch name sanitization

**Input:**
```bash
muxtree new "feature/auth-api" --bg
```

**Expected Output:**
```
▸ Creating worktree at /Users/username/Code/worktrees/repo-name/feature-auth-api
Preparing worktree (new branch 'feature/auth-api')
HEAD is now at abc1234 Initial commit
✓ Worktree created (branch feature/auth-api from master)
▸ Copying config files...
▸ Creating tmux session...
✓   Created session repo-name_feature-auth-api (windows: dev, agent)

✓ Ready! Worktree: /Users/username/Code/worktrees/repo-name/feature-auth-api
```

**Exit Code:** 0

**Behavior:**
- Branch name `feature/auth-api` used as-is for git
- Sanitized to `feature-auth-api` for:
  - Filesystem path (worktree directory)
  - Tmux session name
- Sanitization rule: Replace non-alphanumeric (except underscore, dash, dot) with dash

---

## Color Formatting

- `▸` (info): Blue
- `✓` (success): Green
- `⚠` (warning): Yellow
- `✗` (error): Red
- Branch names in success messages: Cyan
- Base branch in "from X": Dim
- File paths in "Ready!": Cyan
- Command output (pre_session_cmd): Dim, indented

---

## Implementation Notes for Phase 4

**What to implement in Phase 4:**
1. Command structure and CLI parsing (--from, --run, --bg flags)
2. All validation checks (git repo, branch exists, worktree exists, base branch exists, --run validation)
3. Git worktree creation (`git worktree add -b <branch> <path> <base>`)
4. File copying with glob expansion
5. Pre-session command execution with error handling and user confirmation
6. Branch name sanitization
7. All output formatting and colors (except tmux-related output)

**Stub for Phase 5:**
- Tmux session creation (print message but don't create session)
- Terminal opening (--bg flag can be parsed but doesn't need implementation)
- --run command (validate but don't send to tmux)

**Testing approach:**
- Unit tests for sanitization, validation logic
- Integration tests for git operations, file copying
- Feature spec validation (compare output, excluding tmux lines for Phase 4)

---

## Exit Code Summary

| Scenario | Exit Code |
|----------|-----------|
| Success | 0 |
| Not in git repo | 1 |
| Branch already exists | 1 |
| Worktree already exists | 1 |
| Base branch doesn't exist | 1 |
| Invalid --run command | 1 |
| User cancels after pre-session failure | 1 |
| Git worktree creation fails | 1 |

---

## Configuration Values Used

From config files (global `~/.muxtree/config` or project `.muxtree`):
- `worktree_dir`: Base directory for all worktrees (default: `~/worktrees`)
- `copy_files`: Comma-separated list of files/globs to copy
- `pre_session_cmd`: Command to run after worktree creation
- `tmux_layout`: Custom layout (Phase 5)

Environment variables:
- `MUXTREE_CONFIG_DIR`: Override config directory location

---

## Validation Order

The command should perform validations in this order (fail fast):
1. Check if inside git repository
2. Load configuration
3. Validate --run command (if provided)
4. Determine and validate base branch
5. Check if new branch already exists
6. Check if worktree path already exists
7. Proceed with worktree creation

This ensures errors are caught early before any side effects occur.
