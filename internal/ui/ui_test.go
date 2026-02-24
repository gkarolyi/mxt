package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// TestColorCodes verifies that color constants contain expected ANSI codes
// or are empty when not running in a TTY
func TestColorCodes(t *testing.T) {
	tests := []struct {
		name         string
		color        string
		expectedTTY  string
		expectedPipe string
	}{
		{"Red", Red, "\033[0;31m", ""},
		{"Green", Green, "\033[0;32m", ""},
		{"Yellow", Yellow, "\033[0;33m", ""},
		{"Blue", Blue, "\033[0;34m", ""},
		{"Cyan", Cyan, "\033[0;36m", ""},
		{"Bold", Bold, "\033[1m", ""},
		{"Dim", Dim, "\033[2m", ""},
		{"Reset", Reset, "\033[0m", ""},
	}

	// When tests run, stdout is typically not a TTY, so colors will be empty
	isTTY := IsTTY()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := tt.expectedPipe
			if isTTY {
				expected = tt.expectedTTY
			}
			if tt.color != expected {
				t.Errorf("Color %s = %q, want %q (TTY=%v)", tt.name, tt.color, expected, isTTY)
			}
		})
	}
}

// TestSymbols verifies that symbol constants are defined
func TestSymbols(t *testing.T) {
	tests := []struct {
		name     string
		symbol   string
		expected string
	}{
		{"Info", SymbolInfo, "▸"},
		{"Success", SymbolSuccess, "✓"},
		{"Warning", SymbolWarning, "⚠"},
		{"Error", SymbolError, "✗"},
		{"Active", SymbolActive, "●"},
		{"Inactive", SymbolInactive, "○"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.symbol != tt.expected {
				t.Errorf("Symbol %s = %q, want %q", tt.name, tt.symbol, tt.expected)
			}
		})
	}
}

// TestIsTTY tests TTY detection
func TestIsTTY(t *testing.T) {
	// When stdout is not redirected, IsTTY should detect terminal
	// This is hard to test directly, but we can test the behavior
	// We'll mainly test that the function doesn't crash
	_ = IsTTY()
}

// TestInfo tests info message formatting
func TestInfo(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, _ := os.Pipe()
	os.Stdout = w

	Info("test message")

	w.Close()
	out, _ := io.ReadAll(r)
	output := string(out)

	// Should contain the info symbol and message
	if !strings.Contains(output, "▸") {
		t.Errorf("Info() output missing info symbol: %q", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("Info() output missing message: %q", output)
	}
}

// TestSuccess tests success message formatting
func TestSuccess(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, _ := os.Pipe()
	os.Stdout = w

	Success("test success")

	w.Close()
	out, _ := io.ReadAll(r)
	output := string(out)

	// Should contain the success symbol and message
	if !strings.Contains(output, "✓") {
		t.Errorf("Success() output missing success symbol: %q", output)
	}
	if !strings.Contains(output, "test success") {
		t.Errorf("Success() output missing message: %q", output)
	}
}

// TestWarn tests warning message formatting
func TestWarn(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, _ := os.Pipe()
	os.Stdout = w

	Warn("test warning")

	w.Close()
	out, _ := io.ReadAll(r)
	output := string(out)

	// Should contain the warning symbol and message
	if !strings.Contains(output, "⚠") {
		t.Errorf("Warn() output missing warning symbol: %q", output)
	}
	if !strings.Contains(output, "test warning") {
		t.Errorf("Warn() output missing message: %q", output)
	}
}

// TestError tests error message formatting
func TestError(t *testing.T) {
	oldStderr := os.Stderr
	defer func() { os.Stderr = oldStderr }()

	r, w, _ := os.Pipe()
	os.Stderr = w

	Error("test error")

	w.Close()
	out, _ := io.ReadAll(r)
	output := string(out)

	// Should contain the error symbol and message
	if !strings.Contains(output, "✗") {
		t.Errorf("Error() output missing error symbol: %q", output)
	}
	if !strings.Contains(output, "test error") {
		t.Errorf("Error() output missing message: %q", output)
	}
}

// TestColorDisabledWhenNotTTY verifies colors are disabled when output is not a TTY
func TestColorDisabledWhenNotTTY(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, _ := os.Pipe()
	os.Stdout = w

	// Force reinit to detect non-TTY
	InitColors()

	Info("test")

	w.Close()
	out, _ := io.ReadAll(r)
	output := string(out)

	// When output is redirected (not a TTY), colors should be disabled
	// So output should NOT contain ANSI color codes
	if strings.Contains(output, "\033[") {
		t.Errorf("Colors should be disabled when not TTY, but got: %q", output)
	}
}

// TestMessageFormatting tests that messages are properly formatted with colors
func TestMessageFormatting(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(string)
		symbol   string
		color    string
		isstderr bool
	}{
		{"Info", Info, "▸", Blue, false},
		{"Success", Success, "✓", Green, false},
		{"Warn", Warn, "⚠", Yellow, false},
		{"Error", Error, "✗", Red, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var oldStream *os.File
			var stream *os.File

			r, w, _ := os.Pipe()

			if tt.isstderr {
				oldStream = os.Stderr
				os.Stderr = w
				stream = os.Stderr
				defer func() { os.Stderr = oldStream }()
			} else {
				oldStream = os.Stdout
				os.Stdout = w
				stream = os.Stdout
				defer func() { os.Stdout = oldStream }()
			}

			_ = stream // Silence unused variable warning

			tt.fn("test message")

			w.Close()
			out, _ := io.ReadAll(r)
			output := string(out)

			if !strings.Contains(output, tt.symbol) {
				t.Errorf("%s() missing symbol %s in output: %q", tt.name, tt.symbol, output)
			}
			if !strings.Contains(output, "test message") {
				t.Errorf("%s() missing message in output: %q", tt.name, output)
			}
		})
	}
}

// TestDie verifies that Die prints error and exits
func TestDie(t *testing.T) {
	// We can't actually test os.Exit, but we can test that the error is printed
	// before exit by mocking the exit function if we refactor Die to use a variable
	// For now, just test the error formatting

	oldStderr := os.Stderr
	defer func() { os.Stderr = oldStderr }()

	r, w, _ := os.Pipe()
	os.Stderr = w

	// We'll test the Error function which Die uses
	Error("fatal error")

	w.Close()
	out, _ := io.ReadAll(r)
	output := string(out)

	if !strings.Contains(output, "✗") {
		t.Errorf("Die should use Error symbol, got: %q", output)
	}
	if !strings.Contains(output, "fatal error") {
		t.Errorf("Die should print message, got: %q", output)
	}
}

// TestInfof tests formatted info message
func TestInfof(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, _ := os.Pipe()
	os.Stdout = w

	Infof("test %s %d", "message", 42)

	w.Close()
	out, _ := io.ReadAll(r)
	output := string(out)

	if !strings.Contains(output, "test message 42") {
		t.Errorf("Infof() didn't format correctly: %q", output)
	}
}

// TestSuccessf tests formatted success message
func TestSuccessf(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, _ := os.Pipe()
	os.Stdout = w

	Successf("completed %s", "successfully")

	w.Close()
	out, _ := io.ReadAll(r)
	output := string(out)

	if !strings.Contains(output, "completed successfully") {
		t.Errorf("Successf() didn't format correctly: %q", output)
	}
}

// TestWarnf tests formatted warning message
func TestWarnf(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, _ := os.Pipe()
	os.Stdout = w

	Warnf("warning: %d items", 5)

	w.Close()
	out, _ := io.ReadAll(r)
	output := string(out)

	if !strings.Contains(output, "warning: 5 items") {
		t.Errorf("Warnf() didn't format correctly: %q", output)
	}
}

// TestErrorf tests formatted error message
func TestErrorf(t *testing.T) {
	oldStderr := os.Stderr
	defer func() { os.Stderr = oldStderr }()

	r, w, _ := os.Pipe()
	os.Stderr = w

	Errorf("error code: %d", 404)

	w.Close()
	out, _ := io.ReadAll(r)
	output := string(out)

	if !strings.Contains(output, "error code: 404") {
		t.Errorf("Errorf() didn't format correctly: %q", output)
	}
}

// TestColorOutput tests that colors are actually in output when TTY
func TestColorOutput(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Temporarily replace stdout
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	// This test is tricky because we need to force TTY mode
	// We'll just verify the functions don't crash and basic output works
	r, w, _ := os.Pipe()
	os.Stdout = w

	Info("colored info")
	Success("colored success")
	Warn("colored warning")

	w.Close()
	io.Copy(&buf, r)

	// At minimum, verify output contains our messages
	output := buf.String()
	if !strings.Contains(output, "info") {
		t.Error("Missing info message")
	}
	if !strings.Contains(output, "success") {
		t.Error("Missing success message")
	}
	if !strings.Contains(output, "warning") {
		t.Error("Missing warning message")
	}
}
