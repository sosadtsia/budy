package shell

import (
	"os"
	"os/exec"
	"strings"
)

// Executor handles shell command execution
type Executor struct{}

// NewExecutor creates a new shell executor
func NewExecutor() *Executor {
	return &Executor{}
}

// Execute runs a shell command and returns any error
func (e *Executor) Execute(command string) error {
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
