package shell

// Executor defines the interface for shell command execution
type Executor interface {
	Execute(command string) error
}

// HistoryManager defines the interface for command history management
type HistoryManager interface {
	RecordCommand(command string) error
	GetHistory() []CommandEntry
	GetRecentCommands(n int) []CommandEntry
	GetDirectoryCommands() []CommandEntry
}
