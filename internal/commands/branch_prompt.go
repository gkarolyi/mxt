package commands

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// PromptBranchName prompts for a branch name using the provided reader/writer.
func PromptBranchName(reader io.Reader, writer io.Writer) (string, error) {
	fmt.Fprint(writer, "Branch name: ")
	inputReader := bufio.NewReader(reader)
	input, err := inputReader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read branch name: %w", err)
	}
	branch := strings.TrimSpace(input)
	if branch == "" {
		return "", fmt.Errorf("Branch name is required.")
	}
	return branch, nil
}
