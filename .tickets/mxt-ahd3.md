---
id: mxt-ahd3
status: closed
deps: []
links: []
created: 2026-02-26T03:10:45Z
type: feature
priority: 2
assignee: gkarolyi
---
# Prompt for name when mxt new has no args

If mxt new is invoked without a name argument, prompt the user for a branch/worktree name.

## Behavior
- When args are empty and stdin is a TTY, prompt: "Branch name: ".
- Trim whitespace; if empty, print a clear error ("Branch name is required.") and exit non-zero.
- If stdin is not a TTY, keep existing usage error (no prompt).

## Scope
- Update cmd/mxt/main.go to accept optional branch name and call the prompt.
- Keep NewCommand signature unchanged; pass the prompted branch.
- Update help/usage strings and README examples to show branch name as optional with prompt.

## Acceptance Criteria
- `mxt new` with no args prompts for a name and proceeds using that value.
- Empty input is rejected with a clear error message.
- Non-interactive usage still requires explicit branch arg.

## Testing
- Add a testable prompt helper (inject reader) and cover empty input + whitespace trimming.
## Notes

**2026-02-26T03:54:58Z**

Added interactive branch prompt when no args, updated usage/docs, and added tests.
