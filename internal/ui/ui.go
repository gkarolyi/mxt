// Package ui provides terminal output formatting and color support.
package ui

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// Color constants for terminal output
var (
	Red    string
	Green  string
	Yellow string
	Blue   string
	Cyan   string
	Bold   string
	Dim    string
	Reset  string
)

// Symbol constants for terminal output
const (
	SymbolInfo     = "▸"
	SymbolSuccess  = "✓"
	SymbolWarning  = "⚠"
	SymbolError    = "✗"
	SymbolActive   = "●"
	SymbolInactive = "○"
)

func init() {
	InitColors()
}

// InitColors initializes color codes based on TTY detection
func InitColors() {
	if IsTTY() {
		Red = "\033[0;31m"
		Green = "\033[0;32m"
		Yellow = "\033[0;33m"
		Blue = "\033[0;34m"
		Cyan = "\033[0;36m"
		Bold = "\033[1m"
		Dim = "\033[2m"
		Reset = "\033[0m"
	} else {
		// Disable colors when not a TTY
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Cyan = ""
		Bold = ""
		Dim = ""
		Reset = ""
	}
}

// IsTTY returns true if stdout is a terminal
func IsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// Info prints an info message with blue arrow symbol
func Info(msg string) {
	fmt.Fprintf(os.Stdout, "%s%s%s %s\n", Blue, SymbolInfo, Reset, msg)
}

// Infof prints a formatted info message with blue arrow symbol
func Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	Info(msg)
}

// Success prints a success message with green checkmark symbol
func Success(msg string) {
	fmt.Fprintf(os.Stdout, "%s%s%s %s\n", Green, SymbolSuccess, Reset, msg)
}

// Successf prints a formatted success message with green checkmark symbol
func Successf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	Success(msg)
}

// Warn prints a warning message with yellow warning symbol
func Warn(msg string) {
	fmt.Fprintf(os.Stdout, "%s%s%s %s\n", Yellow, SymbolWarning, Reset, msg)
}

// Warnf prints a formatted warning message with yellow warning symbol
func Warnf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	Warn(msg)
}

// Error prints an error message with red X symbol to stderr
func Error(msg string) {
	fmt.Fprintf(os.Stderr, "%s%s%s %s\n", Red, SymbolError, Reset, msg)
}

// Errorf prints a formatted error message with red X symbol to stderr
func Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	Error(msg)
}

// Die prints an error message and exits with code 1
func Die(msg string) {
	Error(msg)
	os.Exit(1)
}

// Dief prints a formatted error message and exits with code 1
func Dief(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	Die(msg)
}
