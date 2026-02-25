# Feature Spec: Terminal Integration

## Overview

This spec documents the behavior of terminal integration for all four supported terminal types. The terminal integration is triggered when creating or opening tmux sessions without the `--bg` flag.

---

## Test Case: Terminal.app (Default)

**Setup:**
- macOS system with Terminal.app installed
- Configuration: `terminal=terminal` (or not set)
- Existing tmux session: `test-repo_feature-auth`

**Input:**
```
$ mxt sessions open feature-auth
```

**Expected Behavior:**
- Terminal.app activates (comes to foreground)
- New Terminal.app window opens
- New window automatically runs: `tmux attach -t test-repo_feature-auth`
- User sees tmux interface in the new window

**Implementation Method:**
- Uses AppleScript via `osascript`
- Script:
  ```applescript
  tell application "Terminal"
      activate
      do script "tmux attach -t test-repo_feature-auth"
  end tell
  ```

**Error Handling:**
- If Terminal.app not found or script fails, warning is displayed
- Command continues (doesn't fail)
- User sees suggestion: "Run: tmux attach -t test-repo_feature-auth"

---

## Test Case: iTerm2

**Setup:**
- macOS system with iTerm2 installed
- Configuration: `terminal=iterm2`
- Existing tmux session: `test-repo_feature-api`

**Input:**
```
$ mxt sessions open feature-api
```

**Expected Behavior:**
- iTerm2 activates (comes to foreground)
- New iTerm2 window created with default profile
- New window automatically runs: `tmux attach -t test-repo_feature-api`
- User sees tmux interface in the new window

**Implementation Method:**
- Uses AppleScript via `osascript`
- Script:
  ```applescript
  tell application "iTerm"
      activate
      create window with default profile
      tell current session of current window
          write text "tmux attach -t test-repo_feature-api"
      end tell
  end tell
  ```

**Error Handling:**
- If iTerm2 not found or script fails, warning is displayed
- Command continues (doesn't fail)
- User sees suggestion: "Run: tmux attach -t test-repo_feature-api"

---

## Test Case: Ghostty

**Setup:**
- macOS system with Ghostty.app installed
- Configuration: `terminal=ghostty`
- Existing tmux session: `test-repo_fix-bug`

**Input:**
```
$ mxt sessions open fix-bug
```

**Expected Behavior:**
- Ghostty activates (comes to foreground)
- New Ghostty window/tab opens
- Automatically attaches to tmux session
- User sees tmux interface in Ghostty

**Implementation Method:**
- Uses `open` command with args:
  ```bash
  open -a Ghostty --args -e tmux attach -t test-repo_fix-bug
  ```

**Error Handling:**
- If Ghostty.app not found:
  ```
  ⚠ Failed to open Ghostty. Ensure Ghostty.app is installed.
  ⚠ Falling back to current terminal...
  ▸ Run: tmux attach -t test-repo_fix-bug
  ```
- Command continues (doesn't fail)
- User can manually run the attach command

---

## Test Case: Current Terminal

**Setup:**
- Any terminal application (Terminal.app, iTerm2, Ghostty, Alacritty, etc.)
- Configuration: `terminal=current`
- Existing tmux session: `test-repo_feature-new`

**Input:**
```
$ mxt sessions open feature-new
```

**Expected Output:**
```
▸ Attaching to session in current terminal: test-repo_feature-new
(terminal switches to tmux attached mode)
```

**Expected Behavior:**
- NO new window or tab is opened
- Current terminal attaches to tmux session
- Process is replaced with `tmux attach`
- User stays in the same terminal window/tab

**Implementation Method:**
- Direct execution:
  ```bash
  tmux attach -t test-repo_feature-new
  ```
- stdin, stdout, stderr are passed through

**Error Handling:**
- If attach fails:
  ```
  ⚠ Could not attach automatically. Run: tmux attach -t test-repo_feature-new
  ```

**Use Cases:**
- When already working in a terminal and don't want new windows
- Works with any terminal application
- Useful in SSH sessions or remote terminals
- Good for terminal multiplexer users

---

## Test Case: Session Name Escaping

**Setup:**
- Repository: test-repo
- Branch with special characters: `feature/auth-v2.0`
- Sanitized session name: `test-repo_feature-auth-v2.0`
- Terminal: Terminal.app or iTerm2

**Expected Behavior:**
- Session name is properly escaped for AppleScript
- Backslashes are escaped: `\` → `\\`
- Double quotes are escaped: `"` → `\"`
- Session opens without errors

**Test Cases:**
- Session name with dash: Works correctly
- Session name with underscore: Works correctly
- Session name with dot: Works correctly
- Session name with escaped characters: Properly handled

---

## Test Case: Background Mode (--bg)

**Setup:**
- Any terminal configuration
- Existing tmux session: `test-repo_feature-auth`

**Input:**
```
$ mxt sessions open feature-auth --bg
```

**Expected Behavior:**
- NO terminal window is opened
- Tmux session is created in detached mode
- User remains in their current terminal
- Session is ready and can be attached later

**Implementation:**
- Terminal integration is completely skipped when `--bg` flag is present
- Applies to all terminal types (terminal, iterm2, ghostty, current)

---

## Test Case: Invalid Terminal Type

**Setup:**
- Configuration: `terminal=invalid`

**Input:**
```
$ mxt sessions open feature-auth
```

**Expected Output:**
```
⚠ Failed to open terminal: unknown terminal type: invalid (use terminal, iterm2, ghostty, or current)
▸ Run: tmux attach -t test-repo_feature-auth
```

**Expected Behavior:**
- Warning is displayed
- Command continues and completes successfully
- Session is created
- User can manually attach

---

## Test Case: Terminal Not Installed

**Setup:**
- Configuration: `terminal=iterm2`
- iTerm2 is NOT installed on the system

**Input:**
```
$ mxt sessions open feature-auth
```

**Expected Output:**
```
⚠ Failed to open terminal: failed to open iTerm2: <error>
▸ Run: tmux attach -t test-repo_feature-auth
```

**Expected Behavior:**
- Warning is displayed
- Command continues and completes successfully
- Session is created
- User can manually attach using any terminal

---

## Implementation Notes

### AppleScript Safety
- Session names are sanitized before being embedded in AppleScript
- Special characters are properly escaped to prevent injection
- Scripts use double quotes around commands for safety

### Error Philosophy
- Terminal integration failures are NON-FATAL
- Sessions are always created successfully
- Users get clear instructions for manual attachment
- This ensures the tool is robust even when terminal integration fails

### Supported Terminals
1. **terminal** (Terminal.app) - Default, universally available on macOS
2. **iterm2** - Popular third-party terminal
3. **ghostty** - Modern GPU-accelerated terminal
4. **current** - Works with any terminal (no new window)

### Configuration
- Terminal type is set in config: `terminal=<type>`
- Can be set globally (`~/.muxtree/config`)
- Cannot be overridden per-project (intentional design choice)
- Can be overridden with environment variable (future enhancement)

### Design Principles
1. **Graceful degradation**: Failures don't block session creation
2. **Universal compatibility**: "current" works everywhere
3. **Clear feedback**: Users always know what to do next
4. **No assumptions**: Don't assume terminal availability
