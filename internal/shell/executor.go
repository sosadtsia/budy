package shell

import (
	"os"
	"os/exec"
	"strings"
)

// ShellExecutor implements the Executor interface
type ShellExecutor struct{}

// Ensure ShellExecutor implements the Executor interface
var _ Executor = (*ShellExecutor)(nil)

// NewExecutor creates a new shell executor
func NewExecutor() *ShellExecutor {
	return &ShellExecutor{}
}

// Execute runs a shell command and returns any error
func (e *ShellExecutor) Execute(command string) error {
	// Split the command into parts
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil
	}

	// Create the command
	cmd := exec.Command(parts[0], parts[1:]...)

	// Set up standard IO
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	return cmd.Run()
}
