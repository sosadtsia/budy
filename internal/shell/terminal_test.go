package shell

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// MockHistoryManager implements HistoryManager for testing
type MockHistoryManager struct {
	commands []CommandEntry
}

func NewMockHistoryManager(commands []string) *MockHistoryManager {
	entries := []CommandEntry{}
	for _, cmd := range commands {
		entries = append(entries, CommandEntry{Command: cmd})
	}
	return &MockHistoryManager{commands: entries}
}

func (m *MockHistoryManager) RecordCommand(command string) error {
	m.commands = append(m.commands, CommandEntry{Command: command})
	return nil
}

func (m *MockHistoryManager) GetHistory() []CommandEntry {
	return m.commands
}

func (m *MockHistoryManager) GetRecentCommands(n int) []CommandEntry {
	if len(m.commands) <= n {
		return m.commands
	}
	return m.commands[len(m.commands)-n:]
}

func (m *MockHistoryManager) GetDirectoryCommands() []CommandEntry {
	return m.commands
}

// TestNewTerminalReader tests that the appropriate reader is created based on platform
func TestNewTerminalReader(t *testing.T) {
	history := NewMockHistoryManager([]string{})
	reader := NewTerminalReader(history)

	// We can't test platform-specific behavior easily, so just check it's not nil
	if reader == nil {
		t.Error("NewTerminalReader returned nil")
	}
}

// TestSimpleTerminalReader tests the SimpleTerminalReader
func TestSimpleTerminalReader(t *testing.T) {
	// Setup test history
	history := NewMockHistoryManager([]string{"ls -la", "cd /tmp", "echo test"})

	// Create reader
	reader := NewSimpleTerminalReader(history)

	// Test that history is retrieved correctly
	if len(history.GetRecentCommands(5)) != 3 {
		t.Errorf("Expected 3 commands in history, got %d", len(history.GetRecentCommands(5)))
	}

	// Test expandHistory with !1
	result, err := reader.expandHistory("!1", history.GetRecentCommands(5))
	if err != nil {
		t.Errorf("Error expanding history: %v", err)
	}
	if result != "ls -la" {
		t.Errorf("Expected 'ls -la', got '%s'", result)
	}

	// Test expandHistory with !!
	result, err = reader.expandHistory("!!", history.GetRecentCommands(5))
	if err != nil {
		t.Errorf("Error expanding history: %v", err)
	}
	if result != "echo test" {
		t.Errorf("Expected 'echo test', got '%s'", result)
	}
}

// TestMacOSTerminalReader tests the MacOSTerminalReader
func TestMacOSTerminalReader(t *testing.T) {
	// Setup test history
	history := NewMockHistoryManager([]string{"ls -la", "cd /tmp", "echo test"})

	// Create reader
	reader := NewMacOSTerminalReader(history)

	// Test expandHistoryCommand with !1
	result := reader.expandHistoryCommand("!1")
	if result != "echo test" {
		t.Errorf("Expected 'echo test', got '%s'", result)
	}

	// Test expandHistoryCommand with !!
	result = reader.expandHistoryCommand("!!")
	if result != "echo test" {
		t.Errorf("Expected 'echo test', got '%s'", result)
	}

	// Test expandHistoryCommand with !-1
	result = reader.expandHistoryCommand("!-1")
	if result != "echo test" {
		t.Errorf("Expected 'echo test', got '%s'", result)
	}
}

// TestSimpleTerminalReaderReadLine tests the ReadLine method with simulated input
func TestSimpleTerminalReaderReadLine(t *testing.T) {
	// Save original stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Create a pipe to simulate input
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdin = r

	// Write test input to the pipe
	go func() {
		_, err := w.Write([]byte("test command\n"))
		if err != nil {
			t.Errorf("Error writing to pipe: %v", err)
			return
		}

		err = w.Close()
		if err != nil {
			t.Errorf("Error closing pipe writer: %v", err)
		}
	}()

	// Create reader with empty history
	history := NewMockHistoryManager([]string{})
	reader := &SimpleTerminalReader{
		history: history,
		scanner: bufio.NewScanner(os.Stdin),
	}

	// Capture stdout
	oldStdout := os.Stdout
	r2, w2, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}
	os.Stdout = w2

	// Test ReadLine
	input, err := reader.ReadLine("> ")
	if err != nil {
		t.Errorf("Error reading line: %v", err)
	}

	// Restore stdout
	os.Stdout = oldStdout
	err = w2.Close()
	if err != nil {
		t.Errorf("Error closing stdout pipe writer: %v", err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r2)
	if err != nil {
		t.Errorf("Error copying from stdout pipe: %v", err)
	}

	if input != "test command" {
		t.Errorf("Expected 'test command', got '%s'", input)
	}

	// Verify prompt was printed
	if !strings.Contains(buf.String(), "> ") {
		t.Errorf("Prompt not found in output: %s", buf.String())
	}
}
