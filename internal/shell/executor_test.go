package shell

import (
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
