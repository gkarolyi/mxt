package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gkarolyi/mxt/internal/config"
	"github.com/gkarolyi/mxt/internal/git"
	"github.com/gkarolyi/mxt/internal/ui"
)

// WorktreeInfo represents information about a single worktree.
type WorktreeInfo struct {
	BranchName    string
	Path          string
	Insertions    int
	Deletions     int
	SessionName   string
	SessionActive bool
}

// ListCommand lists all managed worktrees for the current repository.
func ListCommand() error {
	// Step 1: Check if inside git repository
	if !git.IsInsideWorkTree() {
		return fmt.Errorf("Not inside a git repository. Run muxtree from within your repo.")
	}

	// Step 2: Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Step 3: Get repository name
	repoName, err := git.GetRepoName()
	if err != nil {
		return fmt.Errorf("failed to get repository name: %w", err)
	}

	// Step 4: Print header
	fmt.Printf("%sWorktrees for %s\n", ui.Bold, ui.CyanText(repoName))
	fmt.Println("════════════════════════════════════════════════════════════════")

	managedDir := filepath.Join(cfg.WorktreeDir, repoName)
	if _, err := os.Stat(managedDir); os.IsNotExist(err) {
		ui.Info(fmt.Sprintf("No worktrees found. Use %s to create one.", ui.BoldText("muxtree new <branch>")))
		return nil
	}

	// Step 5: Get managed worktrees
	worktrees, err := getManagedWorktrees(cfg.WorktreeDir, repoName)
	if err != nil {
		return fmt.Errorf("failed to list worktrees: %w", err)
	}

	// Step 6: Handle no worktrees case
	if len(worktrees) == 0 {
		ui.Info(fmt.Sprintf("No managed worktrees found. Use %s to create one.", ui.BoldText("muxtree new <branch>")))
		fmt.Println()
		return nil
	}

	// Step 7: Display each worktree
	for _, wt := range worktrees {
		fmt.Println() // Blank line before each worktree
		displayWorktree(wt)
	}
	fmt.Println()

	return nil
}

// getManagedWorktrees returns a list of worktrees managed by muxtree (in $WORKTREE_DIR/<repo>/).
func getManagedWorktrees(worktreeDir, repoName string) ([]WorktreeInfo, error) {
	// Run git worktree list --porcelain
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git worktree list failed: %w", err)
	}

	// Parse output
	lines := strings.Split(stdout.String(), "\n")
	var worktrees []WorktreeInfo
	managedBase := filepath.Join(worktreeDir, repoName)

	var currentPath string
	var currentBranch string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "worktree ") {
			// Save previous worktree if it was managed
			if currentPath != "" && currentBranch != "" && strings.HasPrefix(currentPath, managedBase) {
				wt := createWorktreeInfo(currentPath, currentBranch, repoName)
				worktrees = append(worktrees, wt)
			}

			// Start new worktree
			currentPath = strings.TrimPrefix(line, "worktree ")
			currentBranch = ""
		} else if strings.HasPrefix(line, "branch ") {
			// Extract branch name from refs/heads/<branch>
			refPath := strings.TrimPrefix(line, "branch ")
			currentBranch = strings.TrimPrefix(refPath, "refs/heads/")
		} else if line == "" {
			// End of worktree entry - save it if managed
			if currentPath != "" && currentBranch != "" && strings.HasPrefix(currentPath, managedBase) {
				wt := createWorktreeInfo(currentPath, currentBranch, repoName)
				worktrees = append(worktrees, wt)
			}
			currentPath = ""
			currentBranch = ""
		}
	}

	// Handle last worktree if file doesn't end with blank line
	if currentPath != "" && currentBranch != "" && strings.HasPrefix(currentPath, managedBase) {
		wt := createWorktreeInfo(currentPath, currentBranch, repoName)
		worktrees = append(worktrees, wt)
	}

	return worktrees, nil
}

// createWorktreeInfo creates a WorktreeInfo from path and branch, calculating stats and session status.
func createWorktreeInfo(path, branch, repoName string) WorktreeInfo {
	wt := WorktreeInfo{
		BranchName: branch,
		Path:       path,
	}

	// Calculate change statistics
	insertions, deletions := calculateChangeStats(path)
	wt.Insertions = insertions
	wt.Deletions = deletions

	// Check session status
	sessionName := git.GenerateSessionName(repoName, branch)
	wt.SessionName = sessionName
	wt.SessionActive = isSessionActive(sessionName)

	return wt
}

// calculateChangeStats calculates total insertions and deletions for a worktree.
// Returns (insertions, deletions).
func calculateChangeStats(worktreePath string) (int, int) {
	totalInsertions := 0
	totalDeletions := 0

	// Get unstaged changes
	unstagedIns, unstagedDel := getGitDiffStats(worktreePath, false)
	totalInsertions += unstagedIns
	totalDeletions += unstagedDel

	// Get staged changes
	stagedIns, stagedDel := getGitDiffStats(worktreePath, true)
	totalInsertions += stagedIns
	totalDeletions += stagedDel

	return totalInsertions, totalDeletions
}

// getGitDiffStats runs git diff --stat and parses insertions/deletions.
// If staged is true, uses --cached flag.
// Returns (insertions, deletions).
func getGitDiffStats(worktreePath string, staged bool) (int, int) {
	args := []string{"-C", worktreePath, "diff", "--stat", "HEAD"}
	if staged {
		args = []string{"-C", worktreePath, "diff", "--cached", "--stat", "HEAD"}
	}

	cmd := exec.Command("git", args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		// If command fails, return 0 (don't fail the whole operation)
		return 0, 0
	}

	// Parse last line for stats
	// Format: "2 files changed, 5 insertions(+), 3 deletions(-)"
	output := stdout.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return 0, 0
	}

	lastLine := lines[len(lines)-1]
	insertions := extractNumber(lastLine, `(\d+) insertion`)
	deletions := extractNumber(lastLine, `(\d+) deletion`)

	return insertions, deletions
}

// extractNumber extracts a number from a string using a regex pattern.
// Returns 0 if not found.
func extractNumber(text, pattern string) int {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		num, err := strconv.Atoi(matches[1])
		if err == nil {
			return num
		}
	}
	return 0
}

// isSessionActive checks if a tmux session exists.
// Returns true if session is active, false otherwise.
func isSessionActive(sessionName string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", sessionName)
	err := cmd.Run()
	return err == nil // Exit code 0 means session exists
}

// displayWorktree prints information about a single worktree.
func displayWorktree(wt WorktreeInfo) {
	// Line 1: Branch name + change stats
	branchText := ui.BoldText(ui.CyanText(wt.BranchName))
	insertionsText := ui.GreenText(fmt.Sprintf("+%d", wt.Insertions))
	deletionsText := ui.RedText(fmt.Sprintf("-%d", wt.Deletions))
	fmt.Printf("  %s  %s %s\n", branchText, insertionsText, deletionsText)

	// Line 2: Worktree path
	fmt.Printf("  %s\n", ui.DimText(wt.Path))

	// Line 3: Session status
	var statusSymbol string
	if wt.SessionActive {
		statusSymbol = ui.GreenText("●")
	} else {
		statusSymbol = ui.DimText("○")
	}
	fmt.Printf("  Session: %s %s\n", statusSymbol, wt.SessionName)
}
