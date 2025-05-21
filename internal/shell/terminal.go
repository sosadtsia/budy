package shell

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// TerminalReader defines the interface for terminal input with history support
type TerminalReader interface {
	ReadLine(prompt string) (string, error)
}

// NewTerminalReader creates a new platform-specific terminal reader
func NewTerminalReader(history HistoryManager) TerminalReader {
	switch runtime.GOOS {
	case "darwin":
		// Use macOS specific implementation on Darwin
		return NewMacOSTerminalReader(history)
	default:
		// Fall back to simple implementation on other platforms
		return NewSimpleTerminalReader(history)
	}
}

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

// SimpleTerminalReader provides a basic readline-like interface with history
type SimpleTerminalReader struct {
	history HistoryManager
	scanner *bufio.Scanner
}

// NewSimpleTerminalReader creates a new terminal reader
func NewSimpleTerminalReader(history HistoryManager) *SimpleTerminalReader {
	return &SimpleTerminalReader{
		history: history,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// ReadLine reads a line of input with basic history management
func (t *SimpleTerminalReader) ReadLine(prompt string) (string, error) {
	// Show command history suggestions
	recentCommands := t.history.GetRecentCommands(5)
	if len(recentCommands) > 0 {
		fmt.Println("\nRecent commands (use !n to recall):")
		for i, cmd := range recentCommands {
			fmt.Printf("  !%d: %s\n", i+1, cmd.Command)
		}
	}

	// Display prompt
	fmt.Print(prompt)

	// Get input
	if !t.scanner.Scan() {
		return "", fmt.Errorf("error reading input")
	}

	input := strings.TrimSpace(t.scanner.Text())

	// Process history expansion
	if strings.HasPrefix(input, "!") {
		return t.expandHistory(input, recentCommands)
	}

	return input, nil
}

// expandHistory handles history expansion (!n commands)
func (t *SimpleTerminalReader) expandHistory(input string, recentCommands []CommandEntry) (string, error) {
	// Handle !n for command recall
	if len(input) > 1 && input[1] >= '1' && input[1] <= '9' {
		n := int(input[1] - '0')
		if n <= len(recentCommands) {
			// Adjust for 1-based indexing
			cmd := recentCommands[n-1].Command
			fmt.Printf("Executing: %s\n", cmd)
			return cmd, nil
		}
	}

	// Handle !! for most recent command
	if input == "!!" && len(recentCommands) > 0 {
		cmd := recentCommands[0].Command
		fmt.Printf("Executing: %s\n", cmd)
		return cmd, nil
	}

	return input, nil
}
