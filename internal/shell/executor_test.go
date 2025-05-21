package shell

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestNewExecutor(t *testing.T) {
	executor := NewExecutor()

	if executor == nil {
		t.Error("NewExecutor returned nil")
	}
}

func TestExecuteEmptyCommand(t *testing.T) {
	executor := NewExecutor()

	err := executor.Execute("")

	if err != nil {
		t.Errorf("Expected no error for empty command, got %v", err)
	}
}

func TestExecuteEchoCommand(t *testing.T) {
	// This test is skipped in automated testing environments
	// since it would attempt to actually execute the command
	if testing.Short() {
		t.Skip("Skipping test that executes commands")
	}

	executor := NewExecutor()

	err := executor.Execute("echo test")

	if err != nil {
		t.Errorf("Expected no error for echo command, got %v", err)
	}
}

// TestShellExecutor_Execute tests the Execute method
func TestShellExecutor_Execute(t *testing.T) {
	// Save original stdout and create a pipe
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Restore stdout when the test ends
	defer func() {
		os.Stdout = oldStdout
	}()

	// Create a test executor
	executor := NewExecutor()

	// Test empty command
	if err := executor.Execute(""); err != nil {
		t.Errorf("Expected no error for empty command, got %v", err)
	}

	// Test a simple echo command
	err := executor.Execute("echo test")

	// Close the writer to capture the output
	if err := w.Close(); err != nil {
		t.Errorf("Error closing pipe writer: %v", err)
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Errorf("Error copying from pipe: %v", err)
	}

	// Check the result
	if err != nil {
		t.Errorf("Error executing echo command: %v", err)
	}

	// Check that stdout contains the expected output
	if output := buf.String(); !contains(output, "test") {
		t.Errorf("Expected output to contain 'test', got: %q", output)
	}
}

// contains checks if a string contains a substring, trimming whitespace
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}

// TestHelperProcess isn't a real test - it's used as a helper process
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Get the command line args after "--"
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		os.Exit(1)
	}

	// Handle the commands we want to mock
	switch args[0] {
	case "echo":
		fmt.Println(strings.Join(args[1:], " "))
	default:
		os.Exit(1)
	}
	os.Exit(0)
}
