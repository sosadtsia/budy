package ai

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOllamaClient(t *testing.T) {
	// Setup a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", contentType)
		}

		// Check endpoint
		if r.URL.Path != "/api/generate" {
			t.Errorf("Expected /api/generate endpoint, got %s", r.URL.Path)
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"model":"llama3","response":"To list files in a directory, use the 'ls' command.","done":true}`))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create a client with the mock server URL
	client := NewOllamaClient(server.URL, "llama3")

	// Test Ask method
	err := client.Ask("How do I list files?")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestOllamaClientDefaults(t *testing.T) {
	// Test default values
	client := NewOllamaClient("", "")

	// Default server URL
	if client.serverURL != "http://localhost:11434" {
		t.Errorf("Expected default server URL to be http://localhost:11434, got %s", client.serverURL)
	}

	// Default model
	if client.model != "llama3" {
		t.Errorf("Expected default model to be llama3, got %s", client.model)
	}
}
