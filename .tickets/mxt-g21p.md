---
id: mxt-g21p
status: open
deps: []
links: []
created: 2026-02-26T03:10:41Z
type: feature
priority: 2
assignee: gkarolyi
---
# Interactive session selection for sessions commands

When sessions open/close/etc are run without a branch argument, open an interactive selector to choose the target.

## Context
- sessions subcommands currently require a branch argument and error otherwise.

## Behavior
- If branch is missing and stdin/stdout are TTY, prompt with fzf (or equivalent).
- If fzf is not installed, return an actionable error: "install fzf or pass a branch name".
- If the user cancels selection (empty output or fzf exit code 130), print an info message and return nil.

## Selection sources
- open / relaunch: list managed worktrees for the current repo (branch names from getManagedWorktrees).
- close / attach: list only worktrees whose tmux session is active (SessionActive == true).
- Items should include status text, but selection must resolve to the branch name (e.g., "branch\t(active)").

## Implementation Notes
- Reuse getManagedWorktrees from list.go; move to shared helper if needed.
- Build selection after config.Load + git.GetRepoName validation (same as existing command flow).
- Keep action aliases (launch/start, kill/stop, restart) working with the selector path.

## Acceptance Criteria
- Running `mxt sessions <action>` with no branch opens a selector and runs the action for the chosen branch.
- open/relaunch can target any managed worktree; close/attach only show active sessions.
- Non-TTY usage still errors with usage text.
- Help text and usage mention interactive selection.

## Testing
- Add unit tests for list filtering (open vs close/attach) and for "no candidates" messaging.
- If selector logic is abstracted, add tests for cancel handling and non-tty behavior.