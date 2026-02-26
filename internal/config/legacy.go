package config

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ParseLegacyConfig parses the legacy key=value config format.
// This is only used for migrating existing configs to TOML.
func ParseLegacyConfig(r io.Reader) (map[string]string, error) {
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
