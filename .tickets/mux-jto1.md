---
id: mux-jto1
status: in_progress
deps: [mux-htje]
links: []
created: 2026-02-24T22:33:50Z
type: task
priority: 0
assignee: gkarolyi
parent: mux-aofc
---
# Phase 8: Polish & Validation

Final validation, shell completion, documentation, and release.

## Acceptance Criteria

All feature specs pass. Shell completion works. Documentation complete. Ready for v1.0.0 release.


## Notes

**2026-02-26T00:16:46Z**

Completed parity fixes (init/list spacing, glob-only missing file warnings), added mxt bash/zsh completions, updated README/visual docs. Harness runs: help, version, init, config, new (basic + --run), list, delete, sessions open/close/relaunch. Sessions attach + terminal integration still need manual verification; release prep (mux-kdel) still open.

**2026-02-26T00:33:12Z**

Created test/run_feature_specs.sh to automate harness runs and prompt for manual checks (sessions attach, terminal integration, pre-session failure).
