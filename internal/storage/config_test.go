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
			t.Errorf("Loaded config doesn't match original. Got: %+v, Expected: %+v", loadedConfig, testConfig)
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

		// Should get default empty config
		if config.OpenAIAPIKey != "" {
			t.Errorf("Expected empty API key, but got: %s", config.OpenAIAPIKey)
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
}
