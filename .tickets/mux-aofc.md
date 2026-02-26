---
id: mux-aofc
status: closed
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


## Notes

**2026-02-26T02:58:44Z**

Go reimplementation complete with parity verified; legacy muxtree artifacts removed; README and config naming updated; final validation closed.
