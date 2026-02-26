---
id: mxt-suyi
status: open
deps: []
links: []
created: 2026-02-26T04:22:37Z
type: bug
priority: 2
assignee: gkarolyi
---
# Fix go install setup for mxt

Goal: make "go install" work for installing mxt.\n\nSteps: (1) Attempt to install the mxt package using "go install" (document exact command/output). (2) Diagnose the root cause of any failure (module path, tags, build constraints, release artifacts, etc.). (3) Implement the fix so "go install" succeeds.\n\nPreference: the first installable release should be v1.1.0 (if versioning changes are needed, start there).

## Acceptance Criteria

"go install" succeeds for mxt using the documented command. Root cause identified and fixed. First installable release tagged v1.1.0 if release/version changes are required.

