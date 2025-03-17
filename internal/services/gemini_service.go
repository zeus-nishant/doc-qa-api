package services

import (
	"doc-qa-api/internal/gemini"
)

// GeminiService defines the interface for the Gemini service
type GeminiService interface {
	GenerateAnswer(document, question string) (string, error)
}

// geminiServiceImpl is the concrete implementation of GeminiService
type geminiServiceImpl struct {
	client gemini.ClientInterface
}

// NewGeminiService creates a new instance of the Gemini service
func NewGeminiService(client gemini.ClientInterface) GeminiService {
	return &geminiServiceImpl{client: client}
}

// GenerateAnswer implements the GeminiService interface
func (gs *geminiServiceImpl) GenerateAnswer(document, question string) (string, error) {
	return gs.client.GenerateAnswer(document, question)
}