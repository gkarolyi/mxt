---
id: mxt-k4nn
status: closed
deps: []
links: []
created: 2026-02-26T03:10:50Z
type: feature
priority: 2
assignee: gkarolyi
---
# mxt init --reinit rewrites configuration

Add --reinit to mxt init to overwrite existing config without prompting.

## Behavior
- Add --reinit flag to init command (global and --local).
- When target config exists:
  - without --reinit: keep current behavior (warn, show file, prompt to overwrite).
  - with --reinit: skip confirmation and overwrite.
- If used with --import (see mxt-ydfv), --reinit allows overwriting existing TOML; without it, fail with a clear error.

## Scope
- Extend InitCommand/initGlobalConfig/initProjectConfig to accept the reinit flag.
- Update help/usage and completion text to document --reinit.

## Acceptance Criteria
- `mxt init --reinit` overwrites existing config at the target scope.
- Existing behavior remains unchanged when --reinit is not provided.
- Usage text documents the new flag.

## Testing
- Add tests around overwrite decision logic (existing vs reinit).
## Notes

**2026-02-26T04:01:58Z**

Added --reinit flag, overwrite helper with tests, updated help and completions.
