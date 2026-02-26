---
id: mux-gcgl
status: closed
deps: []
links: []
created: 2026-02-24T22:34:35Z
type: task
priority: 2
assignee: gkarolyi
parent: mux-jto1
---
# Run complete feature spec test suite

Execute all feature specs and identify any discrepancies


## Notes

**2026-02-26T00:15:42Z**

Ran harness for help, version, init (via wrappers), config, new (basic + --run), list, delete, sessions open/close/relaunch. Sessions attach and terminal integration not run due to interactive attach/open requirements; needs manual verification.

**2026-02-26T00:33:08Z**

Added test/run_feature_specs.sh to run full harness suite with manual checkpoints for pre-session failure, sessions attach, and terminal integration.

**2026-02-26T01:00:10Z**

Updated run_feature_specs.sh to normalize init timestamp line and added custom layout window list parity fixes; non-interactive portion now passes until manual prompts.

**2026-02-26T01:17:06Z**

Updated run_feature_specs.sh: auto-runs manual commands, logs stdout/stderr/exit codes, continues on mismatches, and records summary/failure logs.

**2026-02-26T01:45:44Z**

Full run_feature_specs.sh completed. Summary log: /var/folders/8w/ht214zm55gs6dsckv8flbnm80000gp/T/mxt-feature-suite.XXXXXX.WJKp8Z5o3H/logs/summary.log. Failure log empty. Manual checks for pre-session failure, sessions attach, and terminal integration all marked PASS.
