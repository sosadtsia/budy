package shell

import (
	"testing"
)

// MockStorage is a simple in-memory storage implementation for testing
type MockStorage struct {
	data map[string]interface{}
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		data: make(map[string]interface{}),
	}
}

func (m *MockStorage) Save(key string, value interface{}) error {
	m.data[key] = value
	return nil
}

func (m *MockStorage) Load(key string, value interface{}) error {
	if data, ok := m.data[key]; ok {
		// For testing, we'll just assign the data directly if it's a slice of CommandEntry
		if entries, ok := data.([]CommandEntry); ok {
			if v, ok := value.(*[]CommandEntry); ok {
				*v = entries
				return nil
			}
		}
	}
	// If key doesn't exist or types don't match, initialize with empty slice
	if v, ok := value.(*[]CommandEntry); ok {
		*v = []CommandEntry{}
	}
	return nil
}

func (m *MockStorage) GetDataDir() string {
	return "/mock/data/dir"
}

// TestHistoryManager tests the history manager
func TestHistoryManager(t *testing.T) {
	storage := NewMockStorage()
	history := NewHistoryManager(storage)

	// Test initial state
	if len(history.GetHistory()) != 0 {
		t.Errorf("Expected empty history, got %d entries", len(history.GetHistory()))
	}

	// Test recording a command
	cmd1 := "test command 1"
	if err := history.RecordCommand(cmd1); err != nil {
		t.Errorf("Error recording command: %v", err)
	}

	// Test history has the command
	entries := history.GetHistory()
	if len(entries) != 1 {
		t.Errorf("Expected 1 history entry, got %d", len(entries))
	}
	if entries[0].Command != cmd1 {
		t.Errorf("Expected command '%s', got '%s'", cmd1, entries[0].Command)
	}

	// Add more commands
	cmd2 := "test command 2"
	cmd3 := "test command 3"
	if err := history.RecordCommand(cmd2); err != nil {
		t.Errorf("Error recording command 2: %v", err)
	}
	if err := history.RecordCommand(cmd3); err != nil {
		t.Errorf("Error recording command 3: %v", err)
	}

	// Test GetRecentCommands
	recent := history.GetRecentCommands(2)
	if len(recent) != 2 {
		t.Errorf("Expected 2 recent entries, got %d", len(recent))
	}
	if recent[0].Command != cmd2 || recent[1].Command != cmd3 {
		t.Errorf("Recent commands in wrong order: %v", recent)
	}

	// Test retrieving all history
	all := history.GetHistory()
	if len(all) != 3 {
		t.Errorf("Expected 3 history entries, got %d", len(all))
	}
	if all[0].Command != cmd1 || all[1].Command != cmd2 || all[2].Command != cmd3 {
		t.Errorf("History commands in wrong order: %v", all)
	}
}
