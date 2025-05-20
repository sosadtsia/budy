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
	// Initialize storage
	store, err := storage.NewFileStorage()
	if err != nil {
		fmt.Printf("Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	// Load configuration
	config, err := storage.LoadConfig(store.GetDataDir())
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Get OpenAI API key from configuration or environment
	apiKey := storage.GetOpenAIKey(config)

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
	fmt.Println("Type 'config set openai_key <your_key>' to configure your OpenAI API key")
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

		// Process the input
		processInput(input, aiClient, executor, history, store.GetDataDir(), config)
	}
}

// processInput handles user input and dispatches to the appropriate handler
func processInput(
	input string,
	aiClient ai.Client,
	executor shell.Executor,
	history shell.HistoryManager,
	dataDir string,
	config *storage.Config,
) bool {
	// Handle configuration commands
	if strings.HasPrefix(input, "config") {
		return processConfigCommand(input, aiClient, executor, history, dataDir, config, storage.SetOpenAIKey)
	}

	// Handle question or command
	if strings.HasPrefix(input, "?") {
		query := strings.TrimSpace(input[1:])
		if err := aiClient.Ask(query); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		return true
	} else {
		// Record and execute command
		if err := history.RecordCommand(input); err != nil {
			fmt.Printf("Warning: Failed to record command in history: %v\n", err)
		}
		if err := executor.Execute(input); err != nil {
			fmt.Printf("Error executing command: %v\n", err)
		}
		return true
	}
}

// Helper type for SetOpenAIKey function
type SetOpenAIKeyFunc func(dataDir string, config *storage.Config, key string) error

// processConfigCommand handles configuration commands
func processConfigCommand(
	input string,
	aiClient ai.Client,
	executor shell.Executor,
	history shell.HistoryManager,
	dataDir string,
	config *storage.Config,
	setOpenAIKeyFunc SetOpenAIKeyFunc,
) bool {
	parts := strings.Fields(input)
	if len(parts) >= 3 && parts[1] == "set" && parts[2] == "openai_key" {
		if len(parts) >= 4 {
			key := parts[3]
			if err := setOpenAIKeyFunc(dataDir, config, key); err != nil {
				fmt.Printf("Error setting API key: %v\n", err)
			} else {
				fmt.Println("OpenAI API key set successfully")
				// In a real app, we would update the AI client here with the new key
				return true
			}
		} else {
			fmt.Println("Usage: config set openai_key <your_api_key>")
		}
		return true
	}
	return false
}
