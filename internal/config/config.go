// Package config handles configuration file loading and parsing.
package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Config represents the mxt configuration
type Config struct {
	WorktreeDir   string
	Terminal      string
	CopyFiles     string
	PreSessionCmd string
	TmuxLayout    string
}

// Load loads the configuration from defaults, global config, and project config.
// It returns a Config struct with all values populated.
func Load() (*Config, error) {
	// Get current working directory for LoadConfig
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Load raw config map
	configMap, err := LoadConfig(workDir)
	if err != nil {
		return nil, err
	}

	// Convert map to struct
	cfg := &Config{
		WorktreeDir:   configMap["worktree_dir"],
		Terminal:      configMap["terminal"],
		CopyFiles:     configMap["copy_files"],
		PreSessionCmd: configMap["pre_session_cmd"],
		TmuxLayout:    configMap["tmux_layout"],
	}

	return cfg, nil
}

// ParseConfig parses a config file in TOML format.
// It supports standard TOML comments and validates keys.
// Returns a map of config key-value pairs.
func ParseConfig(r io.Reader) (map[string]string, error) {
	decoder := toml.NewDecoder(r)
	raw := map[string]any{}
	if err := decoder.Decode(&raw); err != nil {
		return nil, err
	}

	config := make(map[string]string)
	for key, value := range raw {
		switch key {
		case "worktree_dir", "terminal", "pre_session_cmd":
			parsed, err := parseStringValue(key, value)
			if err != nil {
				return nil, err
			}
			config[key] = parsed
		case "copy_files":
			parsed, err := parseStringOrArrayValue(key, value, ",")
			if err != nil {
				return nil, err
			}
			config[key] = parsed
		case "tmux_layout":
			parsed, err := parseStringOrArrayValue(key, value, " ")
			if err != nil {
				return nil, err
			}
			config[key] = normalizeTmuxLayout(parsed)
		default:
			return nil, fmt.Errorf("unknown config key %q", key)
		}
	}
	return config, nil
}

func parseStringValue(key string, value any) (string, error) {
	parsed, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("config key %q must be a string", key)
	}
	return parsed, nil
}

func parseStringOrArrayValue(key string, value any, joiner string) (string, error) {
	if parsed, ok := value.(string); ok {
		return parsed, nil
	}
	switch typed := value.(type) {
	case []string:
		return strings.Join(typed, joiner), nil
	case []any:
		result := make([]string, 0, len(typed))
		for _, item := range typed {
			str, ok := item.(string)
			if !ok {
				return "", fmt.Errorf("config key %q array values must be strings", key)
			}
			result = append(result, str)
		}
		return strings.Join(result, joiner), nil
	default:
		return "", fmt.Errorf("config key %q must be a string or array of strings", key)
	}
}

// normalizeTmuxLayout normalizes separators in tmux_layout values.
// It converts commas to semicolons and handles space-separated window definitions.
// Window definitions are in format: window_name:command|command
// Multiple windows can be separated by commas, semicolons, or newlines (spaces in multi-line)
func normalizeTmuxLayout(layout string) string {
	layout = strings.ReplaceAll(layout, "\n", " ")
	layout = strings.ReplaceAll(layout, "\r", " ")
	// First, replace commas with semicolons
	layout = strings.ReplaceAll(layout, ",", ";")

	// If already contains semicolons, we're mostly done
	if strings.Contains(layout, ";") {
		// Clean up double semicolons
		for strings.Contains(layout, ";;") {
			layout = strings.ReplaceAll(layout, ";;", ";")
		}
		// Trim leading/trailing semicolons
		layout = strings.Trim(layout, ";")
		return layout
	}

	// No semicolons yet - this means it's space-separated (from multi-line array)
	// We need to detect window boundaries by looking for window_name: pattern
	// A window starts with word(s) followed by colon

	// Strategy: look for pattern where we see "word:" that starts a window
	// We'll scan character by character to find word: patterns
	var result []string
	var currentWindow strings.Builder

	runes := []rune(layout)
	i := 0

	for i < len(runes) {
		// Skip leading whitespace
		for i < len(runes) && (runes[i] == ' ' || runes[i] == '\t') {
			i++
		}

		if i >= len(runes) {
			break
		}

		// Check if this looks like start of window (word followed by :)
		// Look ahead to find the next colon
		colonPos := -1
		for j := i; j < len(runes); j++ {
			if runes[j] == ':' {
				colonPos = j
				break
			}
			if runes[j] == ' ' || runes[j] == '\t' {
				// Space before colon, not a window name
				break
			}
		}

		// If we found a colon and we already have content, start new window
		if colonPos != -1 && currentWindow.Len() > 0 {
			result = append(result, strings.TrimSpace(currentWindow.String()))
			currentWindow.Reset()
		}

		// Read until we find the next potential window start or end
		if colonPos != -1 {
			// Read through this window definition
			// Find the end - either next window name pattern or end of string
			windowEnd := len(runes)

			// Look ahead for next window pattern (space followed by word:)
			for j := colonPos + 1; j < len(runes)-1; j++ {
				if runes[j] == ' ' {
					// Check if what follows looks like window_name:
					k := j + 1
					for k < len(runes) && runes[k] == ' ' {
						k++
					}
					if k < len(runes) {
						// Look for : after the word
						hasColon := false
						for m := k; m < len(runes) && runes[m] != ' '; m++ {
							if runes[m] == ':' {
								hasColon = true
								windowEnd = j
								break
							}
						}
						if hasColon {
							break
						}
					}
				}
			}

			currentWindow.WriteString(string(runes[i:windowEnd]))
			i = windowEnd
		} else {
			// No colon found, just read to end
			currentWindow.WriteString(string(runes[i:]))
			break
		}
	}

	// Don't forget the last window
	if currentWindow.Len() > 0 {
		result = append(result, strings.TrimSpace(currentWindow.String()))
	}

	// Join with semicolons
	if len(result) > 0 {
		layout = strings.Join(result, ";")
	}

	return layout
}
