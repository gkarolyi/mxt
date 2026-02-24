# Ralph Loop: Implement Phase 0 & 1 for mxt

## Objective

Implement Phase 0 (Setup & Validation) and Phase 1 (Core Infrastructure) for the mxt Go reimplementation of muxtree, following the specification in SPEC.md and managing work through tk tickets.

## Phase 0 Goals (Ticket: mux-jz4d)

**Setup the foundation for feature parity validation:**

1. **Create test infrastructure** (mux-rkrm)
   - Create `test/features/` directory structure
   - Add README explaining feature spec format

2. **Build test harness** (mux-8f7k)
   - Create `test/harness.sh` script that:
     - Runs both muxtree and mxt with identical inputs
     - Captures stdout, stderr, exit codes
     - Compares outputs and reports differences
     - Can be run with: `./test/harness.sh <command> <args>`
   - Example usage: `./test/harness.sh new feature-branch`

3. **Initialize Go project** (mux-605l)
   - Run `go mod init github.com/anthropics/mxt` (or appropriate path)
   - Create directory structure:
     ```
     cmd/mxt/main.go
     internal/ui/
     internal/config/
     internal/git/
     internal/worktree/
     internal/tmux/
     internal/terminal/
     ```
   - Add basic `main.go` that prints version and exits

## Phase 1 Goals (Ticket: mux-dsj1)

**Build foundational infrastructure:**

1. **UI/Output formatting with TDD** (mux-llkf, mux-yzos)
   - Write tests in `internal/ui/ui_test.go`:
     - Test color code output (red, green, yellow, blue, cyan, bold, dim)
     - Test TTY detection (colors disabled when not TTY)
     - Test message formatting: `Info()`, `Success()`, `Warn()`, `Error()`
     - Test symbols: ▸, ✓, ⚠, ✗, ●, ○
   - Implement in `internal/ui/ui.go`:
     - Color constants and TTY detection
     - `Info(msg string)`, `Success(msg string)`, `Warn(msg string)`, `Error(msg string)` functions
     - `Die(msg string)` function (error + exit 1)

2. **CLI framework integration** (mux-gxpj)
   - Add cobra dependency: `go get github.com/spf13/cobra@latest`
   - Set up root command in `cmd/mxt/main.go`
   - Create basic command structure (commands can be stubs for now)
   - Test with: `go run cmd/mxt/main.go --help`

3. **Error handling framework** (mux-7lif)
   - Create `internal/errors/errors.go`:
     - Custom error types for common scenarios
     - Error wrapping utilities
     - Integration with ui.Error() and ui.Die()

## Success Criteria

You must output the completion promise when ALL of the following are true:

### Phase 0 Complete:
- [ ] `test/features/` directory exists with README
- [ ] `test/harness.sh` script exists and is executable
- [ ] `test/harness.sh` can run muxtree commands (even if mxt doesn't exist yet)
- [ ] Go project initialized with `go.mod`
- [ ] Directory structure created with placeholder files
- [ ] `go build cmd/mxt/main.go` compiles successfully
- [ ] `tk list --status=closed | grep mux-rkrm` shows ticket closed
- [ ] `tk list --status=closed | grep mux-8f7k` shows ticket closed
- [ ] `tk list --status=closed | grep mux-605l` shows ticket closed
- [ ] `tk close mux-jz4d` succeeds

### Phase 1 Complete:
- [ ] All UI tests pass: `go test ./internal/ui/...`
- [ ] UI functions work correctly with color output
- [ ] Cobra CLI framework integrated
- [ ] `./mxt --help` shows help output
- [ ] Error handling framework implemented
- [ ] `tk list --status=closed | grep mux-llkf` shows ticket closed
- [ ] `tk list --status=closed | grep mux-yzos` shows ticket closed
- [ ] `tk list --status=closed | grep mux-gxpj` shows ticket closed
- [ ] `tk list --status=closed | grep mux-7lif` shows ticket closed
- [ ] `tk close mux-dsj1` succeeds

### Final Validation:
- [ ] `go test ./...` - all tests pass
- [ ] `go build cmd/mxt/main.go` - builds without errors
- [ ] `./mxt --help` - displays help (even if minimal)
- [ ] `tk show mux-jz4d` - shows status: closed
- [ ] `tk show mux-dsj1` - shows status: closed

## Completion Promise

When all success criteria are met, output:

```
<promise>PHASE 0 AND 1 COMPLETE - FOUNDATION READY</promise>
```

## Important Guidelines

### TDD Approach
- Write tests BEFORE implementation for all testable units
- Run tests frequently: `go test ./...`
- Ensure tests fail first, then make them pass

### Ticket Management
- Start tickets before working: `tk start <ticket-id>`
- Add notes during work: `tk add-note <ticket-id> "Progress update"`
- Close tickets when acceptance criteria met: `tk close <ticket-id>`
- Check progress: `tk list --status=closed | wc -l`

### Code Quality
- Use idiomatic Go
- Add comments for exported functions
- Handle errors properly
- Keep functions focused and testable

### Iterative Improvement
- If tests fail, examine output and fix issues
- If build fails, address compilation errors
- If tickets won't close, verify acceptance criteria met
- Use git commits to track progress

## Key Files to Create

```
test/features/README.md          - Feature spec documentation
test/harness.sh                  - Test harness script
cmd/mxt/main.go                  - Entry point
internal/ui/ui.go                - UI/output functions
internal/ui/ui_test.go           - UI tests
internal/errors/errors.go        - Error handling
go.mod                           - Go module definition
```

## Reference

- Full specification: `SPEC.md`
- Epic ticket: `mux-aofc`
- Phase 0 ticket: `mux-jz4d`
- Phase 1 ticket: `mux-dsj1`
- All subtasks: `tk show mux-jz4d`, `tk show mux-dsj1`

## Notes

- This is Phase 0 and 1 of 9 total phases
- Focus on foundation - no commands implemented yet
- Test harness will be used throughout project for feature parity validation
- UI functions will be used by all commands for consistent output
- Phase 2 (Configuration) depends on completion of Phase 1
