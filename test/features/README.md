# Feature Specification Tests

This directory contains feature spec files for mxt commands. Each spec captures the exact expected terminal output for feature parity validation with muxtree.

## Purpose

Feature specs serve as:
1. **Documentation**: Exact behavior of each command
2. **Test cases**: Expected outputs for validation
3. **Regression tests**: Ensure changes don't break existing behavior

## File Format

Each feature spec file follows this format:

```markdown
# Feature Spec: <command-name>

## Test Case: <scenario-name>

**Input:**
```
$ mxt <command> <args>
```

**Expected Output:**
```
<exact terminal output including colors/formatting>
```

**Exit Code:** 0

**Side Effects:**
- Worktree created at path X
- Tmux session "name" exists
- Files copied: .env, config.yml

## Test Case: <another-scenario>
...
```

## Required Test Cases per Command

For each command, document:
- **Happy path**: Successful execution with typical inputs
- **Error cases**: Invalid input, missing dependencies, etc.
- **Edge cases**: Special characters, empty values, boundary conditions

## Creating Feature Specs

To create a feature spec:

1. Run `muxtree <command>` with specific inputs
2. Capture stdout, stderr, exit code
3. Note side effects (files created, sessions launched, etc.)
4. Document in a new spec file
5. Later, compare `mxt <command>` output to validate feature parity

## Testing with Harness

Use the test harness to compare muxtree vs mxt:

```bash
./test/harness.sh <command> <args>
```

The harness will:
- Run both muxtree and mxt with identical inputs
- Compare outputs and exit codes
- Report any differences
- Validate side effects

## Completion Criteria

A command is fully implemented when:
1. `mxt <command>` produces **exactly the same output** as `muxtree <command>`
2. Exit codes match
3. Side effects match (files created, git operations, tmux sessions, etc.)

## Example: new Command

```markdown
# Feature Spec: new

## Test Case: Create simple worktree

**Input:**
```
$ mxt new feature-auth
```

**Expected Output:**
```
▸ Creating worktree at /Users/user/worktrees/myapp/feature-auth
✓ Worktree created (branch feature-auth from main)
▸ Creating tmux session...
✓ Created session myapp_feature-auth (windows: dev, agent)
✓ Ready! Worktree: /Users/user/worktrees/myapp/feature-auth
```

**Exit Code:** 0

**Side Effects:**
- Git worktree created at /Users/user/worktrees/myapp/feature-auth
- Branch "feature-auth" created from "main"
- Tmux session "myapp_feature-auth" running
- Session has 2 windows: "dev" and "agent"
```

## Feature Spec Files

As commands are implemented, create spec files:

- `init_spec.md` - Config initialization
- `config_spec.md` - Config display
- `new_spec.md` - Worktree creation
- `list_spec.md` - Worktree listing
- `delete_spec.md` - Worktree deletion
- `sessions_spec.md` - Session management
- `help_spec.md` - Help display
- `version_spec.md` - Version display
