package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sosadtsia/budy/internal/ai"
	"github.com/sosadtsia/budy/internal/learning"
	"github.com/sosadtsia/budy/internal/shell"
	"github.com/sosadtsia/budy/internal/storage"
)

// Version information
const (
	appName    = "budy"
	appVersion = "0.0.1"
)

func main() {
	// Get OpenAI API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")

	// Initialize storage
	store, err := storage.NewFileStorage()
	if err != nil {
		fmt.Printf("Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	// Initialize shell executor
	executor := shell.NewExecutor()

	// Initialize history manager
	history := shell.NewHistoryManager(store)

	// Initialize AI client
	aiClient := ai.NewOpenAIClient(apiKey)

	// Initialize suggestion engine
	suggestionEngine := learning.NewSuggestionEngine(history)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("%s v%s - Your AI Terminal Assistant\n", appName, appVersion)
	fmt.Println("Type commands normally or prefix with '?' to ask questions")
	fmt.Println("Type 'exit' to quit")

	// Main interaction loop
	for {
		// Show suggestions
		suggestions := suggestionEngine.GetSuggestions()
		for _, suggestion := range suggestions {
			fmt.Println(suggestion)
		}

		// Display prompt
		fmt.Print("\n> ")

		// Get input
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		// Handle exit command
		if input == "exit" {
			break
		}

		// Empty input
		if input == "" {
			continue
		}

		// Handle question or command
		if strings.HasPrefix(input, "?") {
			query := strings.TrimSpace(input[1:])
			if err := aiClient.Ask(query); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		} else {
			// Record and execute command
			if err := history.RecordCommand(input); err != nil {
				fmt.Printf("Warning: Failed to record command in history: %v\n", err)
			}
			if err := executor.Execute(input); err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}
		}
	}
}
