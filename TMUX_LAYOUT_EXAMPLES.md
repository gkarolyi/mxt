# Tmux Layout Examples

## Syntax

### Multi-line Format (Recommended)

```bash
tmux_layout=[
  window_name:pane_cmd1|pane_cmd2
  next_window:cmd
]
```

### Single-line Format

```bash
tmux_layout=window_name:pane_cmd1|pane_cmd2,next_window:cmd
```

### Rules

- `,` or newline separates windows
- `:` separates window name from panes
- `|` separates panes within a window (vertical split - side by side)
- Empty command = shell prompt

## Example 1: Rails Development

```bash
# .muxtree (in your Rails repo)
copy_files=.env,.env.local,CLAUDE.md
pre_session_cmd=bundle install && bin/rails db:migrate
tmux_layout=[
  dev:hx|lazygit
  server:bin/server
  console:bin/rails c
  logs:tail -f log/development.log
  agent:
]
```

**Result:**
- Window "dev": 2 panes side-by-side (helix editor | lazygit)
- Window "server": 1 pane running `bin/server`
- Window "console": 1 pane with Rails console
- Window "logs": 1 pane tailing development log
- Window "agent": 1 pane with shell (for claude/codex)

## Example 2: Node.js Project

```bash
tmux_layout=[
  dev:hx|
  app:npm run dev
  test:npm test -- --watch
  agent:
]
```

**Result:**
- Window "dev": 2 panes side-by-side (helix | shell)
- Window "app": 1 pane running dev server
- Window "test": 1 pane running tests in watch mode
- Window "agent": 1 pane for AI agent

## Example 3: Full Stack Development

```bash
tmux_layout=[
  code:hx|lazygit
  api:cd api && bin/server
  ui:cd ui && npm run dev
  db:psql myapp_development
  logs:tail -f api/log/development.log|tail -f ui/.next/trace
  agent:
]
```

**Result:**
- Window "code": 2 panes (editor | git UI)
- Window "api": Backend server
- Window "ui": Frontend dev server
- Window "db": Database console
- Window "logs": 2 panes (API logs | UI logs)
- Window "agent": 1 pane for AI agent

## Example 4: Simple Two-Window Setup

```bash
tmux_layout=[
  work:hx|
  agent:
]
```

**Result:**
- Window "work": 2 panes side-by-side (helix | shell)
- Window "agent": 1 pane

## Default Layout (if not configured)

If `tmux_layout` is not set, mxt creates:
- Window "dev": 1 pane (shell)
- Window "agent": 1 pane (shell, or runs `--run` command)

## Using with --run Flag

The `--run` flag still works! If you have an "agent" window in your layout:

```bash
mxt new feature-auth --run claude
```

This will send the `claude` command to the first pane of the "agent" window.

## Tips

1. **Empty commands**: Use empty string for shell prompt
   ```bash
   tmux_layout=[
     dev:hx|     # helix in first pane, shell in second
   ]
   ```

2. **Single pane**: Just put command after colon
   ```bash
   tmux_layout=[
     server:bin/server    # Single pane window
   ]
   ```

3. **Multiple panes**: Separate with `|` for side-by-side layout
   ```bash
   tmux_layout=[
     dev:hx|lazygit|     # 3 panes: helix, lazygit, shell
   ]
   ```

4. **Window without commands**: Just shell prompts
   ```bash
   tmux_layout=[
     work:|       # 2 pane window, both shells
   ]
   ```

5. **Complex commands**: Commands with pipes, redirects, etc. work fine
   ```bash
   tmux_layout=[
     logs:tail -f log/development.log
     db:psql -d myapp_development
     api:cd api && bin/server
   ]
   ```

6. **Single-line format**: For simple layouts, skip the brackets
   ```bash
   tmux_layout=dev:hx|,server:bin/server,agent:
   ```
