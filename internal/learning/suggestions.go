package learning

import (
	"fmt"
	"time"

	"github.com/sosadtsia/budy/internal/shell"
)

// SuggestionEngine generates command suggestions based on history
type SuggestionEngine struct {
	history *shell.HistoryManager
}

// NewSuggestionEngine creates a new suggestion engine
func NewSuggestionEngine(history *shell.HistoryManager) *SuggestionEngine {
	return &SuggestionEngine{
		history: history,
	}
}

// GetSuggestions returns command suggestions based on patterns
func (s *SuggestionEngine) GetSuggestions() []string {
	suggestions := []string{}

	// Add time-based suggestions
	timeSuggestions := s.getTimeSuggestions()
	suggestions = append(suggestions, timeSuggestions...)

	// Add directory-based suggestions
	dirSuggestions := s.getDirectorySuggestions()
	suggestions = append(suggestions, dirSuggestions...)

	// Limit to 3 suggestions to avoid overwhelming the user
	if len(suggestions) > 3 {
		suggestions = suggestions[len(suggestions)-3:]
	}

	return suggestions
}

// getTimeSuggestions generates suggestions based on time patterns
func (s *SuggestionEngine) getTimeSuggestions() []string {
	suggestions := []string{}

	// Get all history
	history := s.history.GetHistory()
	if len(history) == 0 {
		return suggestions
	}

	// Get current hour
	currentHour := time.Now().Hour()

	// Count command occurrences by hour
	hourlyCommands := make(map[string]int)

	// Look for commands frequently used at this hour
	for _, entry := range history {
		if entry.Timestamp.Hour() == currentHour {
			hourlyCommands[entry.Command]++
		}
	}

	// Find common commands for this hour
	for cmd, count := range hourlyCommands {
		if count >= 2 {
			suggestions = append(suggestions, fmt.Sprintf("Suggestion: %s (used %d times at this hour)", cmd, count))
		}
	}

	return suggestions
}

// getDirectorySuggestions generates suggestions based on current directory
func (s *SuggestionEngine) getDirectorySuggestions() []string {
	suggestions := []string{}

	// Get commands executed in current directory
	dirCommands := s.history.GetDirectoryCommands()
	if len(dirCommands) == 0 {
		return suggestions
	}

	// Count command occurrences
	commandCounts := make(map[string]int)
	for _, entry := range dirCommands {
		commandCounts[entry.Command]++
	}

	// Find common commands for this directory
	for cmd, count := range commandCounts {
		if count >= 2 {
			suggestions = append(suggestions, fmt.Sprintf("Suggestion: %s (used %d times in this directory)", cmd, count))
		}
	}

	return suggestions
}
