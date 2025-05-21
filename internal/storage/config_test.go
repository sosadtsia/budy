package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigFunctions(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "budy-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: Failed to remove temp directory: %v", err)
		}
	}()

	// Test SaveConfig and LoadConfig
	t.Run("SaveAndLoadConfig", func(t *testing.T) {
		// Create test config
		testConfig := &Config{
			OpenAIAPIKey: "test-api-key",
			AIProvider:   ProviderOpenAI,
			OllamaURL:    "http://localhost:11434",
			OllamaModel:  "llama3",
		}

		// Save config
		err := SaveConfig(tempDir, testConfig)
		if err != nil {
			t.Fatalf("Failed to save config: %v", err)
		}

		// Verify config file exists
		configPath := filepath.Join(tempDir, "config.json")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Fatalf("Config file was not created at %s", configPath)
		}

		// Load config
		loadedConfig, err := LoadConfig(tempDir)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		// Verify data integrity
		if loadedConfig.OpenAIAPIKey != testConfig.OpenAIAPIKey {
			t.Errorf("Loaded config API key doesn't match. Got: %s, Expected: %s",
				loadedConfig.OpenAIAPIKey, testConfig.OpenAIAPIKey)
		}
		if loadedConfig.AIProvider != testConfig.AIProvider {
			t.Errorf("Loaded config AI provider doesn't match. Got: %s, Expected: %s",
				loadedConfig.AIProvider, testConfig.AIProvider)
		}
		if loadedConfig.OllamaURL != testConfig.OllamaURL {
			t.Errorf("Loaded config Ollama URL doesn't match. Got: %s, Expected: %s",
				loadedConfig.OllamaURL, testConfig.OllamaURL)
		}
		if loadedConfig.OllamaModel != testConfig.OllamaModel {
			t.Errorf("Loaded config Ollama model doesn't match. Got: %s, Expected: %s",
				loadedConfig.OllamaModel, testConfig.OllamaModel)
		}
	})

	// Test LoadConfig on non-existent config
	t.Run("LoadNonExistentConfig", func(t *testing.T) {
		// Create a new empty temp directory
		emptyDir, err := os.MkdirTemp("", "budy-empty-config-test")
		if err != nil {
			t.Fatalf("Failed to create empty temp directory: %v", err)
		}
		defer func() {
			if err := os.RemoveAll(emptyDir); err != nil {
				t.Logf("Warning: Failed to remove empty directory: %v", err)
			}
		}()

		// Load config from empty directory
		config, err := LoadConfig(emptyDir)
		if err != nil {
			t.Fatalf("Expected no error when loading non-existent config, but got: %v", err)
		}

		// Should get default empty config with default values
		if config.OpenAIAPIKey != "" {
			t.Errorf("Expected empty API key, but got: %s", config.OpenAIAPIKey)
		}
		if config.AIProvider != ProviderOllama {
			t.Errorf("Expected default AI provider %s, but got: %s", ProviderOllama, config.AIProvider)
		}
		if config.OllamaURL != "http://localhost:11434" {
			t.Errorf("Expected default Ollama URL http://localhost:11434, but got: %s", config.OllamaURL)
		}
		if config.OllamaModel != "llama3" {
			t.Errorf("Expected default Ollama model llama3, but got: %s", config.OllamaModel)
		}
	})

	// Test GetOpenAIKey from config
	t.Run("GetOpenAIKeyFromConfig", func(t *testing.T) {
		// Ensure OPENAI_API_KEY environment variable is not set
		originalEnvValue := os.Getenv("OPENAI_API_KEY")
		if err := os.Unsetenv("OPENAI_API_KEY"); err != nil {
			t.Fatalf("Failed to unset environment variable: %v", err)
		}
		defer func() {
			if err := os.Setenv("OPENAI_API_KEY", originalEnvValue); err != nil {
				t.Logf("Warning: Failed to restore environment variable: %v", err)
			}
		}()

		config := &Config{
			OpenAIAPIKey: "config-api-key",
		}

		key := GetOpenAIKey(config)
		if key != "config-api-key" {
			t.Errorf("Expected 'config-api-key', but got: %s", key)
		}
	})

	// Test GetOpenAIKey from environment
	t.Run("GetOpenAIKeyFromEnv", func(t *testing.T) {
		// Set environment variable
		originalEnvValue := os.Getenv("OPENAI_API_KEY")
		if err := os.Setenv("OPENAI_API_KEY", "env-api-key"); err != nil {
			t.Fatalf("Failed to set environment variable: %v", err)
		}
		defer func() {
			if err := os.Setenv("OPENAI_API_KEY", originalEnvValue); err != nil {
				t.Logf("Warning: Failed to restore environment variable: %v", err)
			}
		}()

		config := &Config{
			OpenAIAPIKey: "config-api-key",
		}

		key := GetOpenAIKey(config)
		if key != "env-api-key" {
			t.Errorf("Expected 'env-api-key', but got: %s", key)
		}
	})

	// Test SetOpenAIKey
	t.Run("SetOpenAIKey", func(t *testing.T) {
		config := &Config{}
		newKey := "new-api-key"

		err := SetOpenAIKey(tempDir, config, newKey)
		if err != nil {
			t.Fatalf("Failed to set API key: %v", err)
		}

		if config.OpenAIAPIKey != newKey {
			t.Errorf("Config API key not updated. Got: %s, Expected: %s", config.OpenAIAPIKey, newKey)
		}

		// Verify by loading the config
		loadedConfig, err := LoadConfig(tempDir)
		if err != nil {
			t.Fatalf("Failed to load config after setting key: %v", err)
		}

		if loadedConfig.OpenAIAPIKey != newKey {
			t.Errorf("Loaded config API key doesn't match. Got: %s, Expected: %s", loadedConfig.OpenAIAPIKey, newKey)
		}
	})

	// Test SetAIProvider
	t.Run("SetAIProvider", func(t *testing.T) {
		config := &Config{}

		// Test setting to OpenAI
		err := SetAIProvider(tempDir, config, ProviderOpenAI)
		if err != nil {
			t.Fatalf("Failed to set AI provider: %v", err)
		}

		if config.AIProvider != ProviderOpenAI {
			t.Errorf("Config AI provider not updated. Got: %s, Expected: %s", config.AIProvider, ProviderOpenAI)
		}

		// Test setting to Ollama
		err = SetAIProvider(tempDir, config, ProviderOllama)
		if err != nil {
			t.Fatalf("Failed to set AI provider: %v", err)
		}

		if config.AIProvider != ProviderOllama {
			t.Errorf("Config AI provider not updated. Got: %s, Expected: %s", config.AIProvider, ProviderOllama)
		}

		// Verify by loading the config
		loadedConfig, err := LoadConfig(tempDir)
		if err != nil {
			t.Fatalf("Failed to load config after setting provider: %v", err)
		}

		if loadedConfig.AIProvider != ProviderOllama {
			t.Errorf("Loaded config AI provider doesn't match. Got: %s, Expected: %s",
				loadedConfig.AIProvider, ProviderOllama)
		}
	})

	// Test SetOllamaSettings
	t.Run("SetOllamaSettings", func(t *testing.T) {
		config := &Config{
			OllamaURL:   "http://localhost:11434",
			OllamaModel: "llama3",
		}

		// Test setting URL only
		err := SetOllamaSettings(tempDir, config, "http://example.com:11434", "")
		if err != nil {
			t.Fatalf("Failed to set Ollama URL: %v", err)
		}

		if config.OllamaURL != "http://example.com:11434" {
			t.Errorf("Config Ollama URL not updated. Got: %s, Expected: %s",
				config.OllamaURL, "http://example.com:11434")
		}
		if config.OllamaModel != "llama3" {
			t.Errorf("Config Ollama model changed unexpectedly. Got: %s, Expected: %s",
				config.OllamaModel, "llama3")
		}

		// Test setting model only
		err = SetOllamaSettings(tempDir, config, "", "mistral")
		if err != nil {
			t.Fatalf("Failed to set Ollama model: %v", err)
		}

		if config.OllamaURL != "http://example.com:11434" {
			t.Errorf("Config Ollama URL changed unexpectedly. Got: %s, Expected: %s",
				config.OllamaURL, "http://example.com:11434")
		}
		if config.OllamaModel != "mistral" {
			t.Errorf("Config Ollama model not updated. Got: %s, Expected: %s",
				config.OllamaModel, "mistral")
		}

		// Test setting both URL and model
		err = SetOllamaSettings(tempDir, config, "http://localhost:11434", "llama3")
		if err != nil {
			t.Fatalf("Failed to set Ollama settings: %v", err)
		}

		if config.OllamaURL != "http://localhost:11434" {
			t.Errorf("Config Ollama URL not updated. Got: %s, Expected: %s",
				config.OllamaURL, "http://localhost:11434")
		}
		if config.OllamaModel != "llama3" {
			t.Errorf("Config Ollama model not updated. Got: %s, Expected: %s",
				config.OllamaModel, "llama3")
		}

		// Verify by loading the config
		loadedConfig, err := LoadConfig(tempDir)
		if err != nil {
			t.Fatalf("Failed to load config after setting Ollama settings: %v", err)
		}

		if loadedConfig.OllamaURL != "http://localhost:11434" {
			t.Errorf("Loaded config Ollama URL doesn't match. Got: %s, Expected: %s",
				loadedConfig.OllamaURL, "http://localhost:11434")
		}
		if loadedConfig.OllamaModel != "llama3" {
			t.Errorf("Loaded config Ollama model doesn't match. Got: %s, Expected: %s",
				loadedConfig.OllamaModel, "llama3")
		}
	})
}
