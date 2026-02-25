package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gkarolyi/mxt/internal/config"
)

const separator = "─────────────────────────────────"

// ConfigCommand implements the config command
func ConfigCommand() error {
	// Check for global config
	globalConfigPath := config.GetGlobalConfigPath()

	if _, err := os.Stat(globalConfigPath); err == nil {
		fmt.Printf("Global config: %s\n", globalConfigPath)
		fmt.Println(separator)

		content, err := os.ReadFile(globalConfigPath)
		if err != nil {
			return fmt.Errorf("failed to read global config: %w", err)
		}
		fmt.Print(string(content))
		fmt.Println()
	} else {
		fmt.Println("No global config. Use muxtree init to create one.")
		fmt.Println()
	}

	// Check for project config
	gitRoot, err := config.FindGitRoot(".")
	projectExists := false

	if err == nil {
		// We're in a git repo, check for project config
		projectConfigPath := filepath.Join(gitRoot, ".muxtree")

		if _, err := os.Stat(projectConfigPath); err == nil {
			projectExists = true
			fmt.Printf("Project config: %s (active)\n", projectConfigPath)
			fmt.Println(separator)

			content, err := os.ReadFile(projectConfigPath)
			if err != nil {
				return fmt.Errorf("failed to read project config: %w", err)
			}
			fmt.Print(string(content))
			fmt.Println()
		}
	}

	// If no project config found, show message
	if !projectExists {
		fmt.Println("No project config. Use muxtree init --local to create one.")
	}

	return nil
}
