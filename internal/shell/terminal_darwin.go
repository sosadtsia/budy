//go:build darwin
// +build darwin

package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// MacOSTerminalReader provides a readline-like interface for MacOS
type MacOSTerminalReader struct {
	history             HistoryManager
	reader              *bufio.Reader
	currentHistoryIndex int
}

// NewMacOSTerminalReader creates a new terminal reader for MacOS
func NewMacOSTerminalReader(history HistoryManager) *MacOSTerminalReader {
	return &MacOSTerminalReader{
		history:             history,
		reader:              bufio.NewReader(os.Stdin),
		currentHistoryIndex: -1,
	}
}

// ReadLine reads a line of input
func (t *MacOSTerminalReader) ReadLine(prompt string) (string, error) {
	// Display prompt
	fmt.Print(prompt)

	// Read input line
	input, err := t.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Trim newline characters
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")

	// Handle history command shortcuts
	if strings.HasPrefix(input, "!") {
		historyCmd := t.expandHistoryCommand(input)
		if historyCmd != input {
			fmt.Printf("Executing: %s\n", historyCmd)
			return historyCmd, nil
		}
	}

	return input, nil
}

// expandHistoryCommand handles history expansion (!n commands)
func (t *MacOSTerminalReader) expandHistoryCommand(input string) string {
	history := t.history.GetHistory()
	if len(history) == 0 {
		return input
	}

	// Handle !! for most recent command
	if input == "!!" && len(history) > 0 {
		return history[len(history)-1].Command
	}

	// Handle !n for command recall
	if len(input) > 1 && input[1] >= '1' && input[1] <= '9' {
		n := int(input[1] - '0')
		if n <= len(history) {
			// Adjust for 1-based indexing
			return history[len(history)-n].Command
		}
	}

	// Handle !-n for nth previous command
	if len(input) > 2 && input[1] == '-' && input[2] >= '1' && input[2] <= '9' {
		n := int(input[2] - '0')
		if n <= len(history) {
			return history[len(history)-n].Command
		}
	}

	return input
}
