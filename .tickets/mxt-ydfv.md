---
id: mxt-ydfv
status: closed
deps: []
links: []
created: 2026-02-26T03:10:36Z
type: feature
priority: 2
assignee: gkarolyi
---
# Config uses TOML and supports legacy import

Switch config files to TOML (.toml) and drop legacy parsing; add `mxt init --import` to convert legacy configs.

## Context
- Current config format is key=value in ~/.config/mxt/config and <repo>/.mxt.
- Config loading should move to TOML only; old format is supported only for import.

## Scope
- Global config path: ~/.config/mxt/config.toml (or $MXT_CONFIG_DIR/config.toml).
- Project config path: <repo>/.mxt.toml.
- Update config loader/parser, init writer, config command, help text, README.
- Keep existing keys/semantics: worktree_dir, terminal, copy_files, pre_session_cmd, tmux_layout.
- Preserve tmux_layout normalization (comma/semicolon/space handling).

## TOML format
- Use TOML with string values; allow arrays for copy_files and tmux_layout.
- Arrays are converted to the existing string representation:
  - copy_files array -> comma-joined string
  - tmux_layout array -> join with spaces then normalize to semicolon separators
- Example:
```toml
worktree_dir = "~/worktrees"
terminal = "terminal"
copy_files = [".env", ".env.local"]
pre_session_cmd = "bundle install"
tmux_layout = ["dev:hx|lazygit", "server:bin/server", "agent:"]
```

## Legacy import
- `mxt init --import` reads legacy config and writes TOML to the new path.
- `--local --import` imports from legacy .mxt in repo root; global import uses ~/.config/mxt/config.
- Import does not prompt for values; it maps existing keys directly.
- If legacy file is missing, return a clear error.
- If target TOML exists and `--reinit` is not set, refuse to overwrite.

## Implementation Notes
- Use github.com/pelletier/go-toml/v2 for decoding/encoding (add toml tags or raw structs).
- Convert array values to strings before security validation.

## Non-goals
- No dual-loading or fallback to legacy files during normal config load.

## Acceptance Criteria
- Config loader reads only TOML files from the new paths.
- mxt init writes TOML files and `mxt init --import` converts legacy configs.
- Help output and README mention the new file names and TOML format.
- Legacy config files are ignored unless `--import` is used.

## Testing
- Update config loader tests to parse TOML (string + array cases).
- Add tests for legacy import conversion and erroring on missing legacy file.
## Notes

**2026-02-26T05:23:38Z**

Switched config loading to TOML-only with new paths (~/.config/mxt/config.toml, .mxt.toml) and added mxt init --import for legacy key=value conversion; removed config migrate command. Parser now accepts copy_files/tmux_layout arrays via go-toml/v2, init/import paths updated, completions/docs/help/AGENTS refreshed. Tests: go test ./internal/config ./internal/commands.
