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
