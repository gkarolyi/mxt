package config

import (
	"fmt"
	"strings"
)

// Shell metacharacters that could be used for injection
const dangerousChars = "`$;|&"

// isCommandKey returns true if the key is a command key that should allow metacharacters
func isCommandKey(key string) bool {
	return key == "pre_session_cmd" || key == "tmux_layout"
}

// containsMetacharacters checks if a value contains shell metacharacters
func containsMetacharacters(value string) bool {
	return strings.ContainsAny(value, dangerousChars)
}

// ValidateConfigValue validates a single config key-value pair for security issues.
// For non-command keys, it rejects values containing shell metacharacters.
// For command keys (pre_session_cmd, tmux_layout), it allows metacharacters.
func ValidateConfigValue(key, value string) error {
	// Command keys are allowed to have metacharacters
	if isCommandKey(key) {
		return nil
	}

	// Non-command keys must not have metacharacters
	if containsMetacharacters(value) {
		return fmt.Errorf("suspicious value for '%s': contains shell metacharacters", key)
	}

	return nil
}

// ValidateConfig validates all key-value pairs in a config map.
// Returns an error if any value fails security validation.
func ValidateConfig(config map[string]string) error {
	for key, value := range config {
		if err := ValidateConfigValue(key, value); err != nil {
			return err
		}
	}
	return nil
}
