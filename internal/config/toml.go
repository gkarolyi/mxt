package config

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
)

var knownConfigKeys = map[string]struct{}{
	"worktree_dir":    {},
	"terminal":        {},
	"copy_files":      {},
	"pre_session_cmd": {},
	"tmux_layout":     {},
}

func validateConfigKeys(config map[string]string) error {
	for key := range config {
		if _, ok := knownConfigKeys[key]; !ok {
			return fmt.Errorf("unknown config key %q", key)
		}
	}
	return nil
}

// EncodeConfig renders the provided config map as TOML.
func EncodeConfig(config map[string]string) (string, error) {
	if err := validateConfigKeys(config); err != nil {
		return "", err
	}
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(tomlConfigFromMap(config)); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func tomlConfigFromMap(config map[string]string) tomlConfig {
	var output tomlConfig
	if value, ok := config["worktree_dir"]; ok {
		valueCopy := value
		output.WorktreeDir = &valueCopy
	}
	if value, ok := config["terminal"]; ok {
		valueCopy := value
		output.Terminal = &valueCopy
	}
	if value, ok := config["copy_files"]; ok {
		valueCopy := value
		output.CopyFiles = &valueCopy
	}
	if value, ok := config["pre_session_cmd"]; ok {
		valueCopy := value
		output.PreSessionCmd = &valueCopy
	}
	if value, ok := config["tmux_layout"]; ok {
		valueCopy := value
		output.TmuxLayout = &valueCopy
	}
	return output
}
