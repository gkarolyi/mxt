---
id: mux-k7m8
status: closed
deps: [mux-aofc]
links: []
created: 2026-02-24T22:35:12Z
type: feature
priority: 3
assignee: gkarolyi
---
# Add support for sandbox tool

Add integration with sandbox tools (e.g., Firejail, Docker) to isolate worktree environments.

## Acceptance Criteria

Config option to specify sandbox tool. Worktrees can be created and launched within sandboxed environments.


## Notes

**2026-02-26T05:48:22Z**

Added sandbox_tool config option and sandbox command wrapper. Tmux session creation/attach, terminal opens, and pre-session commands now run through sandbox tool when configured; init prompts and TOML output updated, docs/help/AGENTS refreshed, config parsing/security/tests updated. Tests: go test ./internal/...
