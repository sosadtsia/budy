package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// OpenAIClient handles interactions with the OpenAI API
type OpenAIClient struct {
	apiKey string
}

// Message represents a message in the OpenAI chat
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIRequest represents a request to OpenAI API
type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// OpenAIResponse represents a response from OpenAI API
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey: apiKey,
	}
}

// Ask sends a question to the OpenAI API and displays the response
func (c *OpenAIClient) Ask(query string) error {
	if c.apiKey == "" {
		return fmt.Errorf("OpenAI API key not set (use export OPENAI_API_KEY=your_key)")
	}

	// Create request
	reqBody := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful terminal assistant for Unix/Linux/macOS systems. Provide concise answers for command line usage.",
			},
			{
				Role:    "user",
				Content: query,
			},
		},
	}

	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
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
		return fmt.Errorf("API error: %s", body)
	}

	// Parse response
	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return err
	}

	if len(openAIResp.Choices) > 0 {
		fmt.Println(openAIResp.Choices[0].Message.Content)
	}

	return nil
}
