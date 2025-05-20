package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	OpenAIAPIKey string `json:"openai_api_key"`
}

// LoadConfig loads application configuration from disk
func LoadConfig(dataDir string) (*Config, error) {
	configPath := filepath.Join(dataDir, "config.json")

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{}, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse config
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves application configuration to disk
func SaveConfig(dataDir string, config *Config) error {
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write config file
	configPath := filepath.Join(dataDir, "config.json")
	return os.WriteFile(configPath, data, 0644)
}

// GetOpenAIKey gets the OpenAI API key from either environment or config
func GetOpenAIKey(config *Config) string {
	// First check environment variable
	envKey := os.Getenv("OPENAI_API_KEY")
	if envKey != "" {
		return envKey
	}

	// Fall back to config
	return config.OpenAIAPIKey
}

// SetOpenAIKey sets the OpenAI API key in the config
func SetOpenAIKey(dataDir string, config *Config, key string) error {
	config.OpenAIAPIKey = key
	return SaveConfig(dataDir, config)
}
