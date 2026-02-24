# Visual Guide: Tmux Layout

## Pane Orientation

### Vertical Split (Side-by-Side) - New Default

Using `|` creates panes side-by-side:

```
dev:hx|lazygit
```

```
┌─────────────────────────────────────┐
│              Window: dev            │
├─────────────────┬───────────────────┤
│                 │                   │
│     helix       │     lazygit       │
│   (editor)      │   (git UI)        │
│                 │                   │
│                 │                   │
└─────────────────┴───────────────────┘
```

With 3 panes:

```
dev:hx|lazygit|
```

```
┌──────────────────────────────────────────────┐
│                Window: dev                   │
├──────────────┬──────────────┬────────────────┤
│              │              │                │
│    helix     │   lazygit    │     shell      │
│  (editor)    │   (git UI)   │   (commands)   │
│              │              │                │
│              │              │                │
└──────────────┴──────────────┴────────────────┘
```

## Config Format Comparison

### Old Format (Still Works)

```bash
# Hard to read with long layouts
tmux_layout=dev:hx|lazygit;server:cd api && bin/server;ui:cd ui && npm run dev;logs:tail -f api/log/development.log|tail -f ui/.next/trace;db:psql myapp_dev;agent:
```

### New Multi-line Format (Recommended)

```bash
# Much easier to read and maintain
tmux_layout=[
  dev:hx|lazygit
  server:cd api && bin/server
  ui:cd ui && npm run dev
  logs:tail -f api/log/development.log|tail -f ui/.next/trace
  db:psql myapp_dev
  agent:
]
```

### New Single-line with Commas

```bash
# Cleaner than semicolons
tmux_layout=dev:hx|lazygit,server:cd api && bin/server,ui:cd ui && npm run dev,agent:
```

## Complete Workspace Example

### Configuration

```bash
# .muxtree in your Rails repo
copy_files=.env,.env.local,CLAUDE.md
pre_session_cmd=bundle install && bin/rails db:migrate

tmux_layout=[
  code:hx|lazygit
  server:bin/server
  console:bin/rails c
  logs:tail -f log/development.log
  test:bin/rspec --format documentation
  agent:
]
```

### Visual Result

When you run `muxtree new feature-auth`, you get:

```
┌────────────────────────────────────────────────────────────┐
│ Window 1: code                                             │
├─────────────────────────────┬──────────────────────────────┤
│          helix              │         lazygit              │
│     (text editor)           │      (git interface)         │
└─────────────────────────────┴──────────────────────────────┘

┌────────────────────────────────────────────────────────────┐
│ Window 2: server                                           │
├────────────────────────────────────────────────────────────┤
│  bin/server                                                │
│  => Rails 7.2 server starting on http://localhost:3000    │
└────────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────────┐
│ Window 3: console                                          │
├────────────────────────────────────────────────────────────┤
│  irb(main):001:0>                                          │
│  (Rails console ready)                                     │
└────────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────────┐
│ Window 4: logs                                             │
├────────────────────────────────────────────────────────────┤
│  Started GET "/users" for 127.0.0.1                        │
│  Processing by UsersController#index                       │
└────────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────────┐
│ Window 5: test                                             │
├────────────────────────────────────────────────────────────┤
│  RSpec output with documentation format                    │
└────────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────────┐
│ Window 6: agent                                            │
├────────────────────────────────────────────────────────────┤
│  $ (Ready for Claude Code or Cursor)                       │
└────────────────────────────────────────────────────────────┘
```

## Navigation

Switch between windows:
- `Ctrl-b 1` → code window
- `Ctrl-b 2` → server window
- `Ctrl-b 3` → console window
- `Ctrl-b 4` → logs window
- `Ctrl-b 5` → test window
- `Ctrl-b 6` → agent window

Within a window (like "code" with 2 panes):
- `Ctrl-b →` → move to right pane (lazygit)
- `Ctrl-b ←` → move to left pane (helix)

## Tips

1. **Widescreen optimized**: Vertical splits work great on wide monitors
2. **Editor + Git**: Perfect combo is `editor|git_ui` for code + review
3. **Multiple services**: Use windows for different services (api, ui, db)
4. **Log monitoring**: Create panes with `tail -f` for real-time logs
5. **Test runners**: Dedicate a window to continuous testing
