package utils

import (
	"fmt"
	"os"
	"strings"
)

// PrintColorized prints text with ANSI color
func PrintColorized(text string, colorCode string) {
	fmt.Printf("%s%s\033[0m", colorCode, text)
}

// PrintSuccess prints text in green
func PrintSuccess(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("\033[32m%s\033[0m\n", message)
}

// PrintError prints text in red
func PrintError(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("\033[31m%s\033[0m\n", message)
}

// PrintWarning prints text in yellow
func PrintWarning(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("\033[33m%s\033[0m\n", message)
}

// PrintInfo prints text in blue
func PrintInfo(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("\033[34m%s\033[0m\n", message)
}

// Confirm shows a confirmation prompt and returns the user's choice
func Confirm(prompt string) bool {
	fmt.Printf("%s (y/n): ", prompt)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// IsTerminal checks if the current session is interactive
func IsTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
