package shell

import (
	"os"
	"time"

	"github.com/sosadtsia/budy/internal/storage"
)

// CommandEntry represents a single command in history
type CommandEntry struct {
	Command   string    `json:"command"`
	Timestamp time.Time `json:"timestamp"`
	Directory string    `json:"directory"`
}

// HistoryManager handles command history
type HistoryManager struct {
	storage storage.Storage
	history []CommandEntry
}

// NewHistoryManager creates a new history manager
func NewHistoryManager(storage storage.Storage) *HistoryManager {
	// Load existing history from storage
	var history []CommandEntry
	if err := storage.Load("history", &history); err != nil {
		// If error, start with empty history
		history = []CommandEntry{}
	}

	return &HistoryManager{
		storage: storage,
		history: history,
	}
}

// RecordCommand adds a command to history
func (h *HistoryManager) RecordCommand(command string) error {
	// Get current directory
	dir, err := os.Getwd()
	if err != nil {
		dir = "" // Use empty string if we can't get the directory
	}

	// Create history entry
	entry := CommandEntry{
		Command:   command,
		Timestamp: time.Now(),
		Directory: dir,
	}

	// Add to history
	h.history = append(h.history, entry)

	// Save to storage
	return h.storage.Save("history", h.history)
}

// GetHistory returns the command history
func (h *HistoryManager) GetHistory() []CommandEntry {
	return h.history
}

// GetRecentCommands returns the n most recent commands
func (h *HistoryManager) GetRecentCommands(n int) []CommandEntry {
	if len(h.history) <= n {
		return h.history
	}
	return h.history[len(h.history)-n:]
}

// GetDirectoryCommands returns commands executed in the current directory
func (h *HistoryManager) GetDirectoryCommands() []CommandEntry {
	dir, err := os.Getwd()
	if err != nil {
		return []CommandEntry{}
	}

	var dirCommands []CommandEntry
	for _, entry := range h.history {
		if entry.Directory == dir {
			dirCommands = append(dirCommands, entry)
		}
	}
	return dirCommands
}
