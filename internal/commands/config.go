package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gkarolyi/mxt/internal/config"
	mxtErrors "github.com/gkarolyi/mxt/internal/errors"
	"github.com/gkarolyi/mxt/internal/ui"
)

const separator = "─────────────────────────────────"

// ConfigCommand implements the config command
func ConfigCommand() error {
	globalConfigPath := config.GetGlobalConfigPath()
	hasAny := false

	if _, err := os.Stat(globalConfigPath); err == nil {
		hasAny = true
		fmt.Printf("%sGlobal config:%s %s\n", ui.Bold, ui.Reset, globalConfigPath)
		fmt.Println(separator)

		content, err := os.ReadFile(globalConfigPath)
		if err != nil {
			return fmt.Errorf("failed to read global config: %w", err)
		}
		fmt.Print(string(content))
		fmt.Println()
	}

	projectExists := false
	gitRoot, err := config.FindGitRoot(".")
	if err == nil {
		projectConfigPath := filepath.Join(gitRoot, ".muxtree")
		if _, err := os.Stat(projectConfigPath); err == nil {
			hasAny = true
			projectExists = true
			fmt.Printf("%sProject config:%s %s %s(active)%s\n", ui.Bold, ui.Reset, projectConfigPath, ui.Green, ui.Reset)
			fmt.Println(separator)

			content, err := os.ReadFile(projectConfigPath)
			if err != nil {
				return fmt.Errorf("failed to read project config: %w", err)
			}
			fmt.Print(string(content))
			fmt.Println()
		}
	}

	if hasAny && !projectExists {
		fmt.Printf("%sNo project config. Use %s%smuxtree init --local%s%s to create one.%s\n", ui.Dim, ui.Reset, ui.Bold, ui.Reset, ui.Dim, ui.Reset)
	}

	if !hasAny {
		ui.Warn(fmt.Sprintf("No config found. Run %smuxtree init%s to create one.", ui.Bold, ui.Reset))
		return mxtErrors.ErrConfigNotFound{}
	}

	return nil
}
