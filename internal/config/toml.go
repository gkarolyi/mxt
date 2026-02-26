package config

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
)

var knownConfigKeys = map[string]struct{}{
	"worktree_dir":    {},
	"terminal":        {},
	"sandbox_tool":    {},
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
	data, err := toml.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
