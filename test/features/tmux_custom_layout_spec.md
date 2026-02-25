# Feature Spec: Custom Tmux Layouts

This document captures the exact tmux session structure created by `muxtree new` when a custom `tmux_layout` is configured.

## Test Case: Simple two-pane dev window

**Setup:**
- Config contains: `tmux_layout=dev:hx|lazygit`
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
✓   Created session test-repo_test-branch (windows: dev)

✓ Ready! Worktree: /Users/username/Code/worktrees/test-repo/test-branch
```

**Exit Code:** 0

**Side Effects:**

Tmux session verification:
```bash
$ tmux list-windows -t test-repo_test-branch
0: dev* (2 panes) [...]

$ tmux list-panes -t test-repo_test-branch:dev
0: [...] (active)
1: [...] (active)
```

**Pane commands:**
- Pane 0: `hx` command sent
- Pane 1: `lazygit` command sent

**Layout:**
- Both panes are side-by-side (vertical split)
- Even-horizontal layout applied

---

## Test Case: Multiple windows

**Setup:**
- Config contains: `tmux_layout=dev:hx|lazygit;server:bin/server;agent:`

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
✓   Created session test-repo_test-branch (windows: dev, server, agent)

✓ Ready! Worktree: /Users/username/Code/worktrees/test-repo/test-branch
```

**Exit Code:** 0

**Side Effects:**

Tmux session verification:
```bash
$ tmux list-windows -t test-repo_test-branch
0: dev* (2 panes) [...]
1: server- (1 panes) [...]
2: agent- (1 panes) [...]
```

**Window details:**
- Window 0 (dev): 2 panes with `hx` and `lazygit` commands
- Window 1 (server): 1 pane with `bin/server` command
- Window 2 (agent): 1 pane with no command (shell prompt)

**Default selection:**
- First window (dev) is selected by default

---

## Test Case: Multi-line format

**Setup:**
- Config contains:
  ```
  tmux_layout=[
    dev:hx|lazygit
    server:cd api && bin/server|cd ui && yarn start
    logs:tail -f log/development.log
    agent:
  ]
  ```

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
✓   Created session test-repo_test-branch (windows: dev, server, logs, agent)

✓ Ready! Worktree: /Users/username/Code/worktrees/test-repo/test-branch
```

**Exit Code:** 0

**Side Effects:**

Tmux session verification:
```bash
$ tmux list-windows -t test-repo_test-branch
0: dev* (2 panes) [...]
1: server- (2 panes) [...]
2: logs- (1 panes) [...]
3: agent- (1 panes) [...]
```

**Window details:**
- Window 0 (dev): 2 panes (hx | lazygit)
- Window 1 (server): 2 panes (api server | ui server)
- Window 2 (logs): 1 pane (tail -f)
- Window 3 (agent): 1 pane (shell)

---

## Test Case: Custom layout with --run command

**Setup:**
- Config contains: `tmux_layout=dev:hx;agent:`

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

The `claude` command is sent to the window named "agent":
```bash
$ tmux capture-pane -t test-repo_test-branch:agent -p | head -1
$ claude
```

**Behavior:**
- --run command is sent ONLY if a window named "agent" exists
- Command is sent to pane 0 of the agent window
- If no "agent" window exists, command is silently not sent (no error)

---

## Test Case: No agent window with --run command

**Setup:**
- Config contains: `tmux_layout=dev:hx;server:bin/server` (no agent window)

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
✓   Created session test-repo_test-branch (windows: dev, server)

✓ Ready! Worktree: /Users/username/Code/worktrees/test-repo/test-branch
```

**Exit Code:** 0

**Behavior:**
- Session created successfully
- No error about missing agent window
- --run command simply not executed

---

## Test Case: Empty panes

**Setup:**
- Config contains: `tmux_layout=dev:||` (three empty panes)

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
✓   Created session test-repo_test-branch (windows: dev)

✓ Ready! Worktree: /Users/username/Code/worktrees/test-repo/test-branch
```

**Exit Code:** 0

**Side Effects:**

Tmux session verification:
```bash
$ tmux list-panes -t test-repo_test-branch:dev
0: [...] (active)
1: [...] (active)
2: [...] (active)
```

All three panes show shell prompts (no commands sent).

---

## Test Case: Complex commands with shell metacharacters

**Setup:**
- Config contains: `tmux_layout=server:cd api && npm install && npm start|cd ui && yarn && yarn dev`

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
✓   Created session test-repo_test-branch (windows: server)

✓ Ready! Worktree: /Users/username/Code/worktrees/test-repo/test-branch
```

**Exit Code:** 0

**Behavior:**
- Shell commands with `&&`, `||`, `;`, `|`, etc. are sent as-is to panes
- Tmux executes them in the shell
- No escaping or quoting required

---

## Implementation Algorithm

**Custom Layout Creation Steps:**

1. Parse layout string using `ParseLayout()` function
2. Determine session name: `<repo-name>_<sanitized-branch>`
3. Create first window:
   ```bash
   tmux new-session -d -s <session> -c <worktree-path> -n <first-window-name>
   ```
4. For each subsequent window:
   ```bash
   tmux new-window -t <session> -n <window-name> -c <worktree-path>
   ```
5. For each window, create panes:
   - First pane already exists (created with window)
   - Send command to first pane if non-empty:
     ```bash
     tmux send-keys -t <session>:<window>.0 '<command>' Enter
     ```
   - For each additional pane:
     - Create vertical split (side-by-side):
       ```bash
       tmux split-window -h -t <session>:<window> -c <worktree-path>
       ```
     - Send command if non-empty:
       ```bash
       tmux send-keys -t <session>:<window> '<command>' Enter
       ```
   - If window has multiple panes, apply even layout:
     ```bash
     tmux select-layout -t <session>:<window> even-horizontal
     ```
6. If `--run` command provided:
   - Search window list for window named "agent"
   - If found:
     ```bash
     tmux send-keys -t <session>:agent.0 '<command>' Enter
     ```
7. Select first window:
   ```bash
   tmux select-window -t <session>:<first-window-name>
   ```
8. Print success message:
   ```
   ✓   Created session <session> (windows: <window1>, <window2>, ...)
   ```

---

## Error Cases

### Invalid layout string

**Setup:**
- Config contains: `tmux_layout=dev hx lazygit` (missing colon)

**Input:**
```bash
muxtree new test-branch --bg
```

**Expected Output:**
```
✗ Invalid tmux layout: invalid window spec (missing ':'): dev hx lazygit
```

**Exit Code:** 1

---

### Empty window name

**Setup:**
- Config contains: `tmux_layout=:hx|lazygit` (empty window name)

**Input:**
```bash
muxtree new test-branch --bg
```

**Expected Output:**
```
✗ Invalid tmux layout: empty window name in spec: :hx|lazygit
```

**Exit Code:** 1

---

## Color Formatting

- Success message: Green `✓`
- Error message: Red `✗`
- Session name in message: Bold
- Window list: Regular text, comma-separated

---

## Configuration Values

**Required config:**
- `tmux_layout`: Custom layout string

**Formats supported:**
1. Single-line: `dev:hx|lazygit,server:bin/server,agent:`
2. Multi-line:
   ```
   tmux_layout=[
     dev:hx|lazygit
     server:bin/server
     agent:
   ]
   ```

**Separator normalization:**
- Commas and newlines are converted to semicolons internally
- All three separators (`;`, `,`, newline) work identically

---

## Notes for Implementation

1. **Pane splitting** uses `-h` flag (horizontal layout = vertical split = side-by-side panes)
2. **Pane indexing** in tmux starts at 0 for each window
3. **Command sending** requires target format `<session>:<window>.<pane>` or `<session>:<window>`
4. **Layout selection** (`even-horizontal`) distributes panes evenly
5. **Error handling**: Invalid layout strings should fail early with clear error messages
6. **Working directory**: All panes must start in worktree path (`-c` flag)
7. **Window selection**: First window in layout should be selected by default
