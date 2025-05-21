package learning

import (
	"testing"
	"time"

	"github.com/sosadtsia/budy/internal/shell"
)

// MockHistoryManager is a simple in-memory history manager for testing
type MockHistoryManager struct {
	commands       []shell.CommandEntry
	dirCommands    []shell.CommandEntry
	recentCommands []shell.CommandEntry
}

func NewMockHistoryManager() *MockHistoryManager {
	return &MockHistoryManager{
		commands:       []shell.CommandEntry{},
		dirCommands:    []shell.CommandEntry{},
		recentCommands: []shell.CommandEntry{},
	}
}

func (m *MockHistoryManager) RecordCommand(command string) error {
	entry := shell.CommandEntry{
		Command:   command,
		Timestamp: time.Now(),
		Directory: "/test/dir",
	}
	m.commands = append(m.commands, entry)
	return nil
}

func (m *MockHistoryManager) GetHistory() []shell.CommandEntry {
	return m.commands
}

func (m *MockHistoryManager) GetRecentCommands(n int) []shell.CommandEntry {
	return m.recentCommands
}

func (m *MockHistoryManager) GetDirectoryCommands() []shell.CommandEntry {
	return m.dirCommands
}

// TestSuggestionEngine tests the suggestion engine
func TestSuggestionEngine(t *testing.T) {
	// Create mock history manager
	history := NewMockHistoryManager()

	// Prepare test data with current hour
	now := time.Now()

	// Add commands from current hour
	history.commands = append(history.commands,
		shell.CommandEntry{
			Command:   "hourly-cmd",
			Timestamp: now,
			Directory: "/some/dir",
		},
		shell.CommandEntry{
			Command:   "hourly-cmd",
			Timestamp: now,
			Directory: "/some/dir",
		},
	)

	// Add directory-specific commands
	history.dirCommands = append(history.dirCommands,
		shell.CommandEntry{
			Command:   "dir-cmd",
			Timestamp: now,
			Directory: "/current/dir",
		},
		shell.CommandEntry{
			Command:   "dir-cmd",
			Timestamp: now,
			Directory: "/current/dir",
		},
		shell.CommandEntry{
			Command:   "dir-cmd-2",
			Timestamp: now,
			Directory: "/current/dir",
		},
	)

	// Create suggestion engine
	engine := NewSuggestionEngine(history)

	// Test GetSuggestions
	suggestions := engine.GetSuggestions()

	// Should have suggestions based on time and directory
	if len(suggestions) == 0 {
		t.Errorf("Expected suggestions, got none")
	}

	// Check for time-based suggestion
	foundHourly := false
	for _, s := range suggestions {
		if s == "Suggestion: hourly-cmd (used 2 times at this hour)" {
			foundHourly = true
			break
		}
	}

	if !foundHourly {
		t.Errorf("Expected to find hourly command suggestion, but didn't find it in: %v", suggestions)
	}

	// Check for directory-based suggestion
	foundDir := false
	for _, s := range suggestions {
		if s == "Suggestion: dir-cmd (used 2 times in this directory)" {
			foundDir = true
			break
		}
	}

	if !foundDir {
		t.Errorf("Expected to find directory command suggestion, but didn't find it in: %v", suggestions)
	}

	// Test empty history
	emptyHistory := NewMockHistoryManager()
	emptyEngine := NewSuggestionEngine(emptyHistory)
	emptySuggestions := emptyEngine.GetSuggestions()

	if len(emptySuggestions) != 0 {
		t.Errorf("Expected no suggestions for empty history, got %d", len(emptySuggestions))
	}
}
