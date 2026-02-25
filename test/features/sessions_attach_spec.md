# Feature Spec: sessions attach

## Test Case: Basic Attach - No Window Specified

**Setup:**
- Repository: test-repo
- Existing tmux session: `test-repo_feature-auth`
- Currently not attached to the session
- Session has windows: dev, agent

**Input:**
```
$ mxt sessions attach feature-auth
```

**Expected Output:**
```
(terminal switches to tmux attached mode)
```

**Expected Side Effects:**
- Current terminal attaches to the tmux session
- Last active window is displayed (whichever was selected)
- User sees tmux interface with session

**Exit Code:** 0

---

## Test Case: Attach with Window Selection - dev

**Setup:**
- Repository: test-repo
- Existing tmux session: `test-repo_feature-auth`
- Currently not attached
- Agent window was previously selected

**Input:**
```
$ mxt sessions attach feature-auth dev
```

**Expected Output:**
```
(terminal switches to tmux attached mode)
```

**Expected Side Effects:**
- Terminal attaches to session
- Dev window is selected (becomes active)
- Even if agent was previously active, dev is now active

**Exit Code:** 0

---

## Test Case: Attach with Window Selection - agent

**Setup:**
- Repository: test-repo
- Existing tmux session: `test-repo_feature-api`
- Currently not attached
- Dev window was previously selected

**Input:**
```
$ mxt sessions attach feature-api agent
```

**Expected Output:**
```
(terminal switches to tmux attached mode)
```

**Expected Side Effects:**
- Terminal attaches to session
- Agent window is selected (becomes active)
- Even if dev was previously active, agent is now active

**Exit Code:** 0

---

## Test Case: Error - Session Not Found

**Setup:**
- Repository: test-repo
- No tmux session exists for branch "nonexistent"

**Input:**
```
$ mxt sessions attach nonexistent
```

**Expected Output:**
```
✗ Session not found: test-repo_nonexistent
```

**Exit Code:** 1

---

## Test Case: Error - Invalid Window Name

**Setup:**
- Repository: test-repo
- Existing tmux session: `test-repo_feature-auth`

**Input:**
```
$ mxt sessions attach feature-auth invalid
```

**Expected Output:**
```
✗ Unknown window: invalid (use dev or agent)
```

**Exit Code:** 1

---

## Test Case: Error - Not in Git Repository

**Setup:**
- Current directory: Not inside a git repository

**Input:**
```
$ mxt sessions attach feature-auth
```

**Expected Output:**
```
✗ Not inside a git repository. Run muxtree from within your repo.
```

**Exit Code:** 1

---

## Test Case: Attach to Custom Layout Session

**Setup:**
- Repository: test-repo
- Existing tmux session: `test-repo_feature-ui`
- Session has custom windows: dev, server, agent

**Input:**
```
$ mxt sessions attach feature-ui dev
```

**Expected Output:**
```
(terminal switches to tmux attached mode)
```

**Expected Side Effects:**
- Attaches to session
- Dev window is selected
- Note: Only 'dev' and 'agent' window names are supported for selection

**Exit Code:** 0

---

## Test Case: Attach When Already Attached

**Setup:**
- Repository: test-repo
- Existing tmux session: `test-repo_feature-auth`
- Already attached to this session in current terminal

**Input:**
```
$ mxt sessions attach feature-auth
```

**Expected Output:**
```
(tmux error message about already being attached)
```

**Expected Side Effects:**
- Tmux native behavior for re-attaching
- May create a nested session or show error depending on tmux config

**Exit Code:** Non-zero (tmux error)

---

## Implementation Notes

### Command Behavior
1. Require git repository
2. Load configuration (for consistency)
3. Determine repository name
4. Generate session name: `<repo-name>_<sanitized-branch>`
5. Check if session exists (error if not)
6. If window name provided:
   - Validate window name is "dev" or "agent"
   - Select that window: `tmux select-window -t <session>:<window>`
7. Attach to session: `tmux attach -t <session>`

### Window Selection
- Only "dev" and "agent" are valid window names
- Custom layout may have other windows, but attach command only supports these two
- Window selection happens before attach
- If invalid window name, error before attempting attach

### Process Replacement
- The attach command replaces the current process
- User's shell is replaced with tmux client
- When user detaches (Ctrl-b d), they return to their shell
- Unlike 'open' command, this doesn't spawn a new terminal window

### Use Cases
- Quickly switch to a specific worktree session
- Jump to a particular window (dev or agent)
- Alternative to 'open' when already in a terminal
- Useful in scripts or tmux sessions
