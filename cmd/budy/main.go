package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

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

	// Initialize shell executor
	executor := shell.NewExecutor()

	// Initialize history manager
	history := shell.NewHistoryManager(store)

	// Initialize AI client based on configuration
	var aiClient ai.Client
	if config.AIProvider == storage.ProviderOpenAI {
		apiKey := storage.GetOpenAIKey(config)
		if apiKey == "" {
			fmt.Println("Warning: OpenAI API key not set, falling back to Ollama")
			aiClient = ai.NewOllamaClient(config.OllamaURL, config.OllamaModel)
			fmt.Printf("Using Ollama AI provider with model: %s\n", config.OllamaModel)
		} else {
			aiClient = ai.NewOpenAIClient(apiKey)
			fmt.Printf("Using OpenAI API\n")
		}
	} else {
		// Default to Ollama
		aiClient = ai.NewOllamaClient(config.OllamaURL, config.OllamaModel)
		fmt.Printf("Using Ollama AI provider with model: %s\n", config.OllamaModel)
	}

	// Initialize suggestion engine
	suggestionEngine := learning.NewSuggestionEngine(history)

	// Create a platform-specific terminal reader with history support
	terminal := shell.NewTerminalReader(history)

	fmt.Printf("%s v%s - Your AI Terminal Assistant\n", appName, appVersion)
	fmt.Println("Type commands normally or prefix with '?' to ask questions")
	fmt.Println("Type 'config set ai_provider <openai|ollama>' to switch between providers")
	fmt.Println("Type 'config set ollama_model <model_name>' to change the Ollama model")
	if config.AIProvider == storage.ProviderOpenAI {
		fmt.Println("Type 'config set openai_key <your_key>' to configure your OpenAI API key")
	}

	// Help text for history navigation
	fmt.Println("\nHistory navigation shortcuts:")
	fmt.Println("  !! - Repeat the most recent command")
	fmt.Println("  !1 - Execute the most recent command")
	fmt.Println("  !2 - Execute the second most recent command")
	fmt.Println("  !n - Execute the nth most recent command")

	fmt.Println("\nType 'exit' to quit")

	// Main interaction loop
	for {
		// Show suggestions
		suggestions := suggestionEngine.GetSuggestions()
		for _, suggestion := range suggestions {
			fmt.Println(suggestion)
		}

		// Get input with history support
		input, err := terminal.ReadLine("\n> ")
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			break
		}

		// Handle exit command
		if input == "exit" {
			break
		}

		// Empty input
		if input == "" {
			continue
		}

		// Process the input
		newClient := processInput(input, aiClient, executor, history, store.GetDataDir(), config)
		if newClient != nil {
			aiClient = newClient
		}
	}
}

// checkOllamaConnection tries to check if Ollama is running correctly
func checkOllamaConnection(url string) bool {
	client := &http.Client{
		Timeout: 2 * time.Second, // Short timeout to avoid hanging
	}
	resp, err := client.Get(url + "/api/tags")
	if err != nil {
		return false
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return resp.StatusCode == http.StatusOK
}

// printOllamaConnectionError prints a helpful error message for Ollama connection issues
func printOllamaConnectionError() {
	fmt.Println("")
	fmt.Println("==================================================")
	fmt.Println("ERROR: Could not connect to Ollama server.")
	fmt.Println("==================================================")
	fmt.Println("Make sure Ollama is installed and running:")
	fmt.Println("1. Install from https://ollama.ai")
	fmt.Println("2. Verify Ollama is running (it should start automatically after installation)")
	fmt.Println("3. Check that the Ollama server URL is correct with 'config set ollama_url <url>'")
	fmt.Println("")
	fmt.Println("You can also switch to OpenAI API if you have an API key:")
	fmt.Println("1. Set your API key with 'config set openai_key <your_key>'")
	fmt.Println("2. Switch to OpenAI with 'config set ai_provider openai'")
	fmt.Println("==================================================")
}

// processInput handles user input and dispatches to the appropriate handler
func processInput(
	input string,
	aiClient ai.Client,
	executor shell.Executor,
	history shell.HistoryManager,
	dataDir string,
	config *storage.Config,
) ai.Client {
	// Handle configuration commands
	if strings.HasPrefix(input, "config") {
		return processConfigCommand(input, aiClient, executor, history, dataDir, config)
	}

	// Handle question or command
	if strings.HasPrefix(input, "?") {
		query := strings.TrimSpace(input[1:])
		if err := aiClient.Ask(query); err != nil {
			// Check if this is an Ollama connection error
			if strings.Contains(err.Error(), "connection refused") &&
				(config.AIProvider == storage.ProviderOllama ||
					strings.Contains(err.Error(), "Ollama")) {
				printOllamaConnectionError()
			} else {
				fmt.Printf("Error: %v\n", err)
			}

			// If OpenAI fails, try to fall back to Ollama
			if config.AIProvider == storage.ProviderOpenAI {
				// Only try to fall back if Ollama appears to be running
				if checkOllamaConnection(config.OllamaURL) {
					fmt.Println("Trying fallback to Ollama...")
					ollamaClient := ai.NewOllamaClient(config.OllamaURL, config.OllamaModel)
					if err := ollamaClient.Ask(query); err != nil {
						fmt.Printf("Fallback also failed: %v\n", err)
					} else {
						return ollamaClient // If successful, switch to Ollama
					}
				}
			}
		}
		return nil
	} else {
		// Record and execute command
		if err := history.RecordCommand(input); err != nil {
			fmt.Printf("Warning: Failed to record command in history: %v\n", err)
		}
		if err := executor.Execute(input); err != nil {
			fmt.Printf("Error executing command: %v\n", err)
		}
		return nil
	}
}

// processConfigCommand handles configuration commands
func processConfigCommand(
	input string,
	aiClient ai.Client,
	executor shell.Executor,
	history shell.HistoryManager,
	dataDir string,
	config *storage.Config,
) ai.Client {
	parts := strings.Fields(input)
	if len(parts) < 3 {
		fmt.Println("Usage: config set <option> <value>")
		return nil
	}

	if parts[1] != "set" {
		fmt.Println("Unknown config command. Use 'config set <option> <value>'")
		return nil
	}

	switch parts[2] {
	case "openai_key":
		if len(parts) < 4 {
			fmt.Println("Usage: config set openai_key <your_api_key>")
			return nil
		}

		key := parts[3]
		if err := storage.SetOpenAIKey(dataDir, config, key); err != nil {
			fmt.Printf("Error setting API key: %v\n", err)
			return nil
		}

		fmt.Println("OpenAI API key set successfully")

		// If we're using OpenAI, update the client
		if config.AIProvider == storage.ProviderOpenAI {
			return ai.NewOpenAIClient(key)
		}

	case "ai_provider":
		if len(parts) < 4 {
			fmt.Println("Usage: config set ai_provider <openai|ollama>")
			return nil
		}

		provider := strings.ToLower(parts[3])
		if provider != storage.ProviderOpenAI && provider != storage.ProviderOllama {
			fmt.Printf("Invalid provider: %s. Use 'openai' or 'ollama'\n", provider)
			return nil
		}

		if provider == storage.ProviderOllama {
			// Check Ollama connection before switching
			if !checkOllamaConnection(config.OllamaURL) {
				printOllamaConnectionError()
				fmt.Println("Not switching to Ollama due to connection issues.")
				return nil
			}
		}

		if err := storage.SetAIProvider(dataDir, config, provider); err != nil {
			fmt.Printf("Error setting AI provider: %v\n", err)
			return nil
		}

		fmt.Printf("AI provider set to %s\n", provider)

		// Create and return new AI client based on provider
		if provider == storage.ProviderOpenAI {
			apiKey := storage.GetOpenAIKey(config)
			if apiKey == "" {
				fmt.Println("Warning: OpenAI API key not set. You need to set it with 'config set openai_key <your_key>'")
				fmt.Println("Staying with Ollama for now...")
				return ai.NewOllamaClient(config.OllamaURL, config.OllamaModel)
			}
			return ai.NewOpenAIClient(apiKey)
		} else {
			return ai.NewOllamaClient(config.OllamaURL, config.OllamaModel)
		}

	case "ollama_url":
		if len(parts) < 4 {
			fmt.Println("Usage: config set ollama_url <url>")
			return nil
		}

		url := parts[3]

		// Test the new URL before saving it
		fmt.Printf("Testing connection to Ollama at %s...\n", url)
		if !checkOllamaConnection(url) {
			fmt.Println("Warning: Could not connect to Ollama at this URL.")
			fmt.Println("Saving anyway, but you may need to correct it later.")
		}

		if err := storage.SetOllamaSettings(dataDir, config, url, ""); err != nil {
			fmt.Printf("Error setting Ollama URL: %v\n", err)
			return nil
		}

		fmt.Printf("Ollama URL set to %s\n", url)

		// If we're using Ollama, update the client
		if config.AIProvider == storage.ProviderOllama {
			return ai.NewOllamaClient(url, config.OllamaModel)
		}

	case "ollama_model":
		if len(parts) < 4 {
			fmt.Println("Usage: config set ollama_model <model_name>")
			return nil
		}

		model := parts[3]
		if err := storage.SetOllamaSettings(dataDir, config, "", model); err != nil {
			fmt.Printf("Error setting Ollama model: %v\n", err)
			return nil
		}

		fmt.Printf("Ollama model set to %s\n", model)

		// If we're using Ollama, update the client
		if config.AIProvider == storage.ProviderOllama {
			return ai.NewOllamaClient(config.OllamaURL, model)
		}

	default:
		fmt.Printf("Unknown config option: %s\n", parts[2])
	}

	return nil
}
