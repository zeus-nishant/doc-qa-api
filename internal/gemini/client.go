package gemini

import (
	"context"
	"fmt"
	"log"

	generativeai "github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// ClientInterface defines the methods that the Gemini client must implement.
type ClientInterface interface {
	GenerateAnswer(document, question string) (string, error)
}

// Client implements the Gemini client.
type Client struct {
	client *generativeai.Client
}

// NewClient creates a new Gemini client.
func NewClient(apiKey string) (*Client, error) {
	ctx := context.Background()
	client, err := generativeai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("Failed to create Gemini client: %v", err)
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}
	return &Client{client: client}, nil
}

// GenerateAnswer generates an answer using the Gemini API.
func (c *Client) GenerateAnswer(document, question string) (string, error) {
	ctx := context.Background()

	// Create a prompt
	prompt := fmt.Sprintf("Analyze this document and answer the question.\nDocument: %s\nQuestion: %s\nAnswer:", document, question)
	log.Printf("Generated prompt: %s", prompt)

	// Initialize the model
	model := c.client.GenerativeModel("gemini-1.5-flash")

	// Generate content
	resp, err := model.GenerateContent(ctx, generativeai.Text(prompt))
	if err != nil {
		log.Printf("API call failed: %v", err)
		return "", fmt.Errorf("API call failed: %w", err)
	}

	// Log the full response for debugging
	log.Printf("API response: %+v", resp)

	// Extract the answer
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		log.Println("No answer generated")
		return "", fmt.Errorf("no answer generated")
	}

	answer := resp.Candidates[0].Content.Parts[0].(generativeai.Text)
	log.Printf("Generated answer: %s", answer)
	return string(answer), nil
}