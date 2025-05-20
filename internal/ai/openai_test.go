package ai

import (
	"testing"
)

func TestNewOpenAIClient(t *testing.T) {
	client := NewOpenAIClient("test-api-key")

	if client.apiKey != "test-api-key" {
		t.Errorf("Expected API key 'test-api-key', got '%s'", client.apiKey)
	}
}

func TestAskNoAPIKey(t *testing.T) {
	client := NewOpenAIClient("")

	err := client.Ask("test question")

	if err == nil {
		t.Error("Expected error when API key is empty, got nil")
	}
}
