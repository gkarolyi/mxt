---
id: mux-aofc
status: open
deps: []
links: []
created: 2026-02-24T22:33:20Z
type: epic
priority: 0
assignee: gkarolyi
---
# mxt: Go reimplementation of muxtree with 100% feature parity

Reimplement muxtree in Go as 'mxt' binary. Goal is 100% feature parity with muxtree v1.0.0, maintaining identical CLI, inputs, outputs, and behaviors. Use TDD where appropriate and create feature specs for each command.

## Acceptance Criteria

All commands (init, config, new, list, delete, sessions, help, version) produce identical output to muxtree. All feature spec tests pass. Original muxtree script can be safely removed.

