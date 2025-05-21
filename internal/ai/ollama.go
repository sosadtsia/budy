package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Ensure OllamaClient implements the Client interface
var _ Client = (*OllamaClient)(nil)

// OllamaClient handles interactions with the Ollama API
type OllamaClient struct {
	serverURL string
	model     string
}

// OllamaRequest represents a request to Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	System string `json:"system,omitempty"`
}

// OllamaResponse represents a response from Ollama API
type OllamaResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at,omitempty"`
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(serverURL string, model string) *OllamaClient {
	// Use default values if not provided
	if serverURL == "" {
		serverURL = "http://localhost:11434"
	}
	if model == "" {
		model = "llama3" // Default model
	}

	return &OllamaClient{
		serverURL: serverURL,
		model:     model,
	}
}

// Ask sends a question to the Ollama API and displays the response
func (c *OllamaClient) Ask(query string) error {
	// Create request
	reqBody := OllamaRequest{
		Model:  c.model,
		Prompt: query,
		Stream: false,
		System: "You are a helpful terminal assistant for Unix/Linux/macOS systems. Provide concise answers for command line usage.",
	}

	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", c.serverURL+"/api/generate", bytes.NewBuffer(reqData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama server: %v", err)
	}

	// Use a closure to properly handle the error from Body.Close()
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error closing response body: %v\n", err)
		}
	}()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, body)
	}

	// Parse response
	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return err
	}

	// Print the response
	fmt.Println(ollamaResp.Response)

	return nil
}
