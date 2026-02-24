# Changelog

## [Unreleased] - 2026-02-25

### Changed

#### Vertical Splits (Side-by-Side Panes)
- **Breaking Change**: `|` now creates vertical splits (side-by-side) instead of horizontal
- Panes are arranged left-to-right for better widescreen use
- Layout uses `even-horizontal` instead of `even-vertical`
- Example: `dev:hx|lazygit` creates helix on left, lazygit on right

#### Multi-line Config Format
- **New Feature**: Support for readable multi-line tmux_layout configuration
- Use brackets `[]` for multi-line definitions
- Separate windows with commas or newlines
- Much easier to read and maintain complex layouts

**Before:**
```bash
tmux_layout=dev:vim|;server:bin/server;logs:tail -f log/development.log;agent:
```

**After (multi-line):**
```bash
tmux_layout=[
  dev:hx|lazygit
  server:bin/server
  logs:tail -f log/development.log
  agent:
]
```

**After (single-line with commas):**
```bash
tmux_layout=dev:hx|lazygit,server:bin/server,logs:tail -f log/development.log,agent:
```

### Migration Guide

If you have existing configs using semicolons, they still work:
```bash
# Old format still works
tmux_layout=dev:vim|;server:bin/server;agent:

# But we recommend updating to:
tmux_layout=[
  dev:vim|
  server:bin/server
  agent:
]
```

### Benefits

1. **More Readable**: Multi-line format makes complex layouts easy to understand
2. **Better Screen Use**: Vertical splits use widescreen monitors efficiently
3. **Flexible**: Both single-line and multi-line formats supported
4. **Natural**: Side-by-side (vertical) panes match typical terminal usage

### Technical Details

- Parser handles `[]` brackets for multi-line arrays
- Commas and newlines both work as window separators
- Semicolons still work for backward compatibility
- `tmux split-window -h` creates vertical (side-by-side) splits
- `even-horizontal` layout distributes panes evenly left-to-right
