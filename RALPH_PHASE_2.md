# Ralph Loop: Implement Phase 2 for mxt

## Objective

Implement Phase 2 (Configuration System) for the mxt Go reimplementation of muxtree, following TDD principles and the specification in SPEC.md.

## Phase 2 Goals (Ticket: mux-30uf)

**Build the configuration system with TDD:**

1. **Write tests for config file parsing (single-line)** (mux-gaau)
   - Test key=value parsing
   - Test comment handling (lines starting with #)
   - Test whitespace trimming
   - Test empty line handling

2. **Write tests for multi-line array parsing** (mux-167e)
   - Test `key=[...]` syntax
   - Test multi-line accumulation
   - Test separator normalization for tmux_layout
   - Test single-line array format

3. **Implement config file parser** (mux-2ao5)
   - Parse single-line key=value
   - Parse multi-line arrays with `[...]`
   - Handle comments and whitespace
   - Normalize tmux_layout separators

4. **Write tests for security validation** (mux-vpqk)
   - Test metacharacter detection (`, $, ;, |, &)
   - Test rejection for non-command fields
   - Test allowance for command fields (pre_session_cmd, tmux_layout)

5. **Implement security validation** (mux-s91d)
   - Validate config values for shell metacharacters
   - Reject suspicious values for worktree_dir, terminal, copy_files
   - Allow metacharacters in pre_session_cmd and tmux_layout

6. **Write tests for config loading priority** (mux-szys)
   - Test defaults only
   - Test global config override
   - Test project config override
   - Test full priority chain (defaults → global → project)

7. **Implement config loading system** (mux-ixni)
   - Load defaults
   - Load global config (~/.muxtree/config)
   - Detect git repo and load project config (.muxtree)
   - Apply priority correctly
   - Expand tilde in worktree_dir

8. **Create feature spec for init command** (mux-aut0)
   - Run muxtree init and capture exact output
   - Run muxtree init --local and capture output
   - Document prompts, file structure, templates
   - Note colors and formatting

9. **Implement init command** (mux-atft)
   - Display logo and version
   - Prompt for config values
   - Create config files with templates
   - Handle --local flag for project config
   - Match muxtree output exactly

10. **Create feature spec for config command** (mux-x5mw)
    - Run muxtree config and capture output
    - Test with global only, project only, both
    - Document formatting and colors

11. **Implement config command** (mux-mdgk)
    - Display global config if exists
    - Display project config if exists
    - Show appropriate messages if missing
    - Match muxtree output exactly

## Success Criteria

You must output the completion promise when ALL of the following are true:

### All Tests Pass:
- [ ] `go test ./internal/config/...` - all config parsing tests pass
- [ ] Config file parsing (single-line) works correctly
- [ ] Multi-line array parsing works correctly
- [ ] Security validation rejects/allows appropriately
- [ ] Config loading priority works correctly

### Commands Implemented:
- [ ] `./mxt init` matches muxtree output
- [ ] `./mxt init --local` matches muxtree output
- [ ] `./mxt config` matches muxtree output
- [ ] Feature specs created for both commands

### Test Harness Validation:
- [ ] `./test/harness.sh init` shows feature parity (or document differences)
- [ ] `./test/harness.sh config` shows feature parity (or document differences)

### Tickets Closed:
- [ ] `tk list --status=closed | grep mux-gaau` - shows ticket closed
- [ ] `tk list --status=closed | grep mux-167e` - shows ticket closed
- [ ] `tk list --status=closed | grep mux-2ao5` - shows ticket closed
- [ ] `tk list --status=closed | grep mux-vpqk` - shows ticket closed
- [ ] `tk list --status=closed | grep mux-s91d` - shows ticket closed
- [ ] `tk list --status=closed | grep mux-szys` - shows ticket closed
- [ ] `tk list --status=closed | grep mux-ixni` - shows ticket closed
- [ ] `tk list --status=closed | grep mux-aut0` - shows ticket closed
- [ ] `tk list --status=closed | grep mux-atft` - shows ticket closed
- [ ] `tk list --status=closed | grep mux-x5mw` - shows ticket closed
- [ ] `tk list --status=closed | grep mux-mdgk` - shows ticket closed
- [ ] `tk close mux-30uf` succeeds

### Final Validation:
- [ ] `go test ./...` - all tests pass
- [ ] `go build cmd/mxt/main.go` - builds without errors
- [ ] `./mxt init --help` - shows help
- [ ] `./mxt config --help` - shows help
- [ ] `tk show mux-30uf` - shows status: closed

### Git Commit:
- [ ] All changes committed with proper message
- [ ] Commit includes closed ticket list
- [ ] Commit includes Co-Authored-By line

## Completion Promise

When all success criteria are met, output:

```
<promise>PHASE 2 COMPLETE - CONFIGURATION SYSTEM READY</promise>
```

## Important Guidelines

### TDD Approach (CRITICAL)
- **ALWAYS write tests BEFORE implementation**
- Run tests to verify they fail initially
- Implement minimal code to make tests pass
- Refactor while keeping tests green

### Test Organization
```
internal/config/
  config.go           # Implementation
  config_test.go      # Unit tests for parsing
  loader.go           # Config loading logic
  loader_test.go      # Tests for loading priority
  security.go         # Security validation
  security_test.go    # Security validation tests
```

### Ticket Management
- Start tickets before working: `tk start <ticket-id>`
- Add notes during work: `tk add-note <ticket-id> "Progress update"`
- Close tickets when acceptance criteria met: `tk close <ticket-id>`
- Check progress: `tk list --status=closed | wc -l`

### Feature Spec Creation
Before implementing init and config commands:
1. Run `muxtree init` and capture exact output
2. Run `muxtree init --local` and capture output
3. Run `muxtree config` in various scenarios
4. Document in `test/features/init_spec.md` and `test/features/config_spec.md`
5. Use as reference for exact output matching

### Code Quality
- Use idiomatic Go
- Add comments for exported functions
- Handle errors properly
- Keep functions focused and testable
- Follow existing code patterns from Phase 1

### Config File Format
Remember the exact format from SPEC.md:
- `key=value` format
- Comments start with `#`
- Multi-line arrays use `key=[...]`
- Whitespace trimmed from keys and values
- Special handling for tmux_layout separators

### Security Requirements
- Reject `` ` ``, `$`, `;`, `|`, `&` in most config values
- ALLOW them in `pre_session_cmd` and `tmux_layout`
- Log warnings for rejected values

## Key Files to Create/Modify

```
internal/config/config.go          # Config struct and parsing
internal/config/config_test.go     # Parsing tests
internal/config/loader.go          # Config loading logic
internal/config/loader_test.go     # Loading tests
internal/config/security.go        # Security validation
internal/config/security_test.go   # Security tests
test/features/init_spec.md         # Feature spec for init
test/features/config_spec.md       # Feature spec for config
cmd/mxt/main.go                    # Update init and config commands
```

## Reference

- Full specification: `SPEC.md` sections on Configuration System
- Phase 2 ticket: `mux-30uf`
- All subtasks: `tk show mux-30uf`
- Previous phases: RALPH_PHASE_0_1.md (completed)

## Notes

- This is Phase 2 of 9 total phases
- Focus on TDD - tests first, implementation second
- Use test harness to validate feature parity with muxtree
- Phase 3 (Git Operations) depends on completion of Phase 2
- Configuration system is critical - all future commands depend on it
