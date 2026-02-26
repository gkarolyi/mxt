---
id: mux-k7a8
status: closed
deps: [mux-aofc]
links: []
created: 2026-02-24T22:35:11Z
type: feature
priority: 2
assignee: gkarolyi
---
# Fail-safe worktree teardown on interrupt

Add cleanup logic to handle Ctrl+C or errors during worktree creation. Should safely rollback partial worktree creation.

## Acceptance Criteria

If muxtree/mxt new is interrupted, partial worktrees are cleaned up automatically. No orphaned directories or git state.


## Notes

**2026-02-26T03:48:42Z**

Added interrupt-safe worktree creation with cleanup helpers and tests.
