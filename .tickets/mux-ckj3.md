---
id: mux-ckj3
status: closed
deps: [mux-aofc]
links: []
created: 2026-02-24T22:35:12Z
type: feature
priority: 3
assignee: gkarolyi
---
# Use TOML (or other) for configuration format

Consider migrating from custom key=value format to TOML for better structure and validation. This is a breaking change and should wait until after feature parity.

## Acceptance Criteria

Config files use TOML format. Migration tool provided for existing configs. Documentation updated.


## Notes

**2026-02-26T04:42:43Z**

Implemented TOML config parsing with legacy migration command, updated init/help/README, added TOML + migration tests. Tests: go test ./internal/config ./internal/commands
