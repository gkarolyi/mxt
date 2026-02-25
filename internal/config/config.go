// Package config handles configuration file loading and parsing.
package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Config represents the muxtree configuration
type Config struct {
	WorktreeDir    string
	Terminal       string
	CopyFiles      string
	PreSessionCmd  string
	TmuxLayout     string
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

// ParseConfig parses a config file in key=value format.
// It handles comments (lines starting with #), empty lines, and whitespace trimming.
// It also supports multi-line arrays with key=[...] syntax.
// Returns a map of config key-value pairs.
func ParseConfig(r io.Reader) (map[string]string, error) {
	config := make(map[string]string)
	scanner := bufio.NewScanner(r)
	lineNum := 0

	var multiLineKey string
	var multiLineValues []string
	inMultiLine := false

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Trim leading and trailing whitespace
		line = strings.TrimSpace(line)

		// Skip empty lines when not in multi-line mode
		if line == "" && !inMultiLine {
			continue
		}

		// Skip comments when not in multi-line mode
		if strings.HasPrefix(line, "#") && !inMultiLine {
			continue
		}

		// Handle multi-line array accumulation
		if inMultiLine {
			// Skip comments within multi-line arrays
			if strings.HasPrefix(line, "#") {
				continue
			}

			// Check for closing bracket
			if line == "]" || strings.HasSuffix(line, "]") {
				// Extract value before ]
				if line != "]" {
					value := strings.TrimSpace(strings.TrimSuffix(line, "]"))
					if value != "" {
						multiLineValues = append(multiLineValues, value)
					}
				}

				// Join accumulated values
				finalValue := strings.Join(multiLineValues, " ")

				// Apply tmux_layout normalization if needed
				if multiLineKey == "tmux_layout" {
					finalValue = normalizeTmuxLayout(finalValue)
				}

				config[multiLineKey] = finalValue
				inMultiLine = false
				multiLineKey = ""
				multiLineValues = nil
				continue
			}

			// Accumulate non-empty lines
			if line != "" {
				multiLineValues = append(multiLineValues, line)
			}
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid format (expected key=value): %s", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Empty key is invalid
		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", lineNum)
		}

		// Check for multi-line array start
		if strings.HasPrefix(value, "[") {
			// Check if it's a single-line array [value1 value2]
			if strings.HasSuffix(value, "]") {
				// Single-line array
				arrayContent := strings.TrimSpace(value[1 : len(value)-1])
				if key == "tmux_layout" {
					arrayContent = normalizeTmuxLayout(arrayContent)
				}
				config[key] = arrayContent
				continue
			}

			// Multi-line array
			multiLineKey = key
			inMultiLine = true
			// Check if there's content after the opening bracket on the same line
			content := strings.TrimSpace(strings.TrimPrefix(value, "["))
			if content != "" {
				multiLineValues = append(multiLineValues, content)
			}
			continue
		}

		// Apply tmux_layout normalization for single-line values
		if key == "tmux_layout" {
			value = normalizeTmuxLayout(value)
		}

		config[key] = value
	}

	// Check for unclosed multi-line array
	if inMultiLine {
		return nil, fmt.Errorf("unclosed multi-line array for key %q", multiLineKey)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	return config, nil
}

// normalizeTmuxLayout normalizes separators in tmux_layout values.
// It converts commas to semicolons and handles space-separated window definitions.
// Window definitions are in format: window_name:command|command
// Multiple windows can be separated by commas, semicolons, or newlines (spaces in multi-line)
func normalizeTmuxLayout(layout string) string {
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
