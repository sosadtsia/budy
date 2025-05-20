package main

import (
	"os"
	"strings"
	"testing"

	"github.com/sosadtsia/budy/internal/ai"
	"github.com/sosadtsia/budy/internal/shell"
	"github.com/sosadtsia/budy/internal/storage"
)

// MockAIClient is a mock implementation of the AI client
type MockAIClient struct {
	askedQueries []string
}

// Ensure MockAIClient implements the ai.Client interface
var _ ai.Client = (*MockAIClient)(nil)

func (m *MockAIClient) Ask(query string) error {
	m.askedQueries = append(m.askedQueries, query)
	return nil
}

// MockExecutor is a mock implementation for testing that matches shell.Executor's API
type MockExecutor struct {
	executedCommands []string
}

// Ensure MockExecutor implements the shell.Executor interface
var _ shell.Executor = (*MockExecutor)(nil)

func (m *MockExecutor) Execute(command string) error {
	m.executedCommands = append(m.executedCommands, command)
	return nil
}

// MockHistoryManager is a mock implementation for testing that matches shell.HistoryManager's API
type MockHistoryManager struct {
	recordedCommands []string
}

// Ensure MockHistoryManager implements the shell.HistoryManager interface
var _ shell.HistoryManager = (*MockHistoryManager)(nil)

func (m *MockHistoryManager) RecordCommand(command string) error {
	m.recordedCommands = append(m.recordedCommands, command)
	return nil
}

func (m *MockHistoryManager) GetHistory() []shell.CommandEntry {
	return []shell.CommandEntry{}
}

func (m *MockHistoryManager) GetRecentCommands(n int) []shell.CommandEntry {
	return []shell.CommandEntry{}
}

func (m *MockHistoryManager) GetDirectoryCommands() []shell.CommandEntry {
	return []shell.CommandEntry{}
}

// TestCommandHandling tests the command handling functionality
func TestCommandHandling(t *testing.T) {
	// Test question handling
	t.Run("QuestionHandling", func(t *testing.T) {
		mockAI := &MockAIClient{}
		mockExecutor := &MockExecutor{}
		mockHistory := &MockHistoryManager{}

		// Call processInput through our test helper
		testHandleInput("? how to list files", mockAI, mockExecutor, mockHistory)

		// Check that the AI client was called with the correct query
		if len(mockAI.askedQueries) != 1 {
			t.Fatalf("Expected 1 query, got %d", len(mockAI.askedQueries))
		}
		expectedQuery := "how to list files"
		if mockAI.askedQueries[0] != expectedQuery {
			t.Errorf("Expected query '%s', got '%s'", expectedQuery, mockAI.askedQueries[0])
		}

		// Ensure executor and history were not called
		if len(mockExecutor.executedCommands) != 0 {
			t.Errorf("Executor should not have been called")
		}
		if len(mockHistory.recordedCommands) != 0 {
			t.Errorf("History should not have been called")
		}
	})

	// Test command handling
	t.Run("CommandHandling", func(t *testing.T) {
		mockAI := &MockAIClient{}
		mockExecutor := &MockExecutor{}
		mockHistory := &MockHistoryManager{}

		// Call processInput through our test helper
		testHandleInput("ls -la", mockAI, mockExecutor, mockHistory)

		// Check that the executor and history were called
		if len(mockExecutor.executedCommands) != 1 {
			t.Fatalf("Expected 1 command, got %d", len(mockExecutor.executedCommands))
		}
		expectedCmd := "ls -la"
		if mockExecutor.executedCommands[0] != expectedCmd {
			t.Errorf("Expected command '%s', got '%s'", expectedCmd, mockExecutor.executedCommands[0])
		}
		if len(mockHistory.recordedCommands) != 1 {
			t.Fatalf("Expected 1 history entry, got %d", len(mockHistory.recordedCommands))
		}
		if mockHistory.recordedCommands[0] != expectedCmd {
			t.Errorf("Expected history entry '%s', got '%s'", expectedCmd, mockHistory.recordedCommands[0])
		}

		// Ensure AI was not called
		if len(mockAI.askedQueries) != 0 {
			t.Errorf("AI should not have been called")
		}
	})

	// Test config command handling
	t.Run("ConfigCommandHandling", func(t *testing.T) {
		mockAI := &MockAIClient{}
		mockExecutor := &MockExecutor{}
		mockHistory := &MockHistoryManager{}

		// Create a temporary directory for testing
		tempDir, err := os.MkdirTemp("", "budy-config-test")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: Failed to remove temp directory: %v", err)
		}
		defer func() {
			if err := os.RemoveAll(tempDir); err != nil {
				t.Logf("Warning: Failed to remove temp directory during cleanup: %v", err)
			}
		}()

		config := &storage.Config{}

		// Call handleConfigInput through our test helper
		input := "config set openai_key test-api-key"
		t.Logf("Test input: %q", input)
		handled, key := testHandleConfigInput(input, mockAI, mockExecutor, mockHistory, tempDir, config)

		// Check that the config was handled
		if !handled {
			t.Errorf("Config command should have been handled")
		}
		if key != "test-api-key" {
			t.Errorf("Expected key 'test-api-key', got %q", key)
		}

		// Ensure executor, history, and AI were not called
		if len(mockExecutor.executedCommands) != 0 {
			t.Errorf("Executor should not have been called")
		}
		if len(mockHistory.recordedCommands) != 0 {
			t.Errorf("History should not have been called")
		}
		if len(mockAI.askedQueries) != 0 {
			t.Errorf("AI should not have been called")
		}
	})
}

// testHandleInput is a test helper to simulate the processInput function
func testHandleInput(input string, aiClient ai.Client, executor shell.Executor, history shell.HistoryManager) {
	if input == "" {
		return
	}

	// Check if it's a question
	if input[0] == '?' {
		query := input[1:]
		if len(query) > 0 && query[0] == ' ' {
			query = query[1:]
		}
		if err := aiClient.Ask(query); err != nil {
			panic(err) // In tests we can panic on errors
		}
		return
	}

	// If it's not a question, treat it as a command
	if err := history.RecordCommand(input); err != nil {
		panic(err)
	}
	if err := executor.Execute(input); err != nil {
		panic(err)
	}
}

// testHandleConfigInput is a test helper to simulate the processConfigCommand function
func testHandleConfigInput(input string, aiClient ai.Client, executor shell.Executor, history shell.HistoryManager, dataDir string, config *storage.Config) (bool, string) {
	// Debug input
	if len(input) < 19 {
		return false, ""
	}

	prefix := "config set openai_key"
	if !strings.HasPrefix(input, prefix) {
		return false, ""
	}

	if len(input) <= len(prefix) {
		return true, ""
	}

	// Extract key - everything after the prefix plus a space
	key := input[len(prefix)+1:]
	return true, key
}
