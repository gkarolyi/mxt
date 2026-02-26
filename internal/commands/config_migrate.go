package commands

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gkarolyi/mxt/internal/config"
	"github.com/gkarolyi/mxt/internal/ui"
)

// ConfigMigrateCommand migrates legacy key=value config files to TOML.
func ConfigMigrateCommand() error {
	migrated := false

	migratedGlobal, err := migrateConfigFile(config.GetGlobalConfigPath())
	if err != nil {
		return err
	}
	migrated = migrated || migratedGlobal

	gitRoot, err := config.FindGitRoot(".")
	if err == nil {
		projectPath := filepath.Join(gitRoot, ".mxt")
		migratedProject, err := migrateConfigFile(projectPath)
		if err != nil {
			return err
		}
		migrated = migrated || migratedProject
	} else {
		ui.Info("No git repository detected; skipping project config migration.")
	}

	if !migrated {
		ui.Warn("No legacy config files found to migrate.")
	}

	return nil
}

func migrateConfigFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to stat config file %s: %w", path, err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	if _, err := config.ParseConfig(bytes.NewReader(content)); err == nil {
		ui.Info(fmt.Sprintf("Config already in TOML format: %s", path))
		return false, nil
	}

	legacyConfig, err := config.ParseLegacyConfig(bytes.NewReader(content))
	if err != nil {
		return false, fmt.Errorf("failed to parse legacy config %s: %w", path, err)
	}

	tomlContent, err := config.EncodeConfig(legacyConfig)
	if err != nil {
		return false, fmt.Errorf("failed to encode TOML config %s: %w", path, err)
	}

	if err := os.WriteFile(path, []byte(tomlContent), info.Mode().Perm()); err != nil {
		return false, fmt.Errorf("failed to write TOML config %s: %w", path, err)
	}

	ui.Success(fmt.Sprintf("Migrated config to TOML: %s", path))
	return true, nil
}
