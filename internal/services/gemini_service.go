package services

import (
	"doc-qa-api/internal/gemini"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

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
	log.WithFields(logrus.Fields{
		"document_length": len(document),
		"question":        question,
	}).Info("Processing GenerateAnswer request")

	answer, err := gs.client.GenerateAnswer(document, question)
	if err != nil {
		log.WithError(err).Error("Error generating answer from Gemini")
		return "", err
	}

	log.WithFields(logrus.Fields{
		"answer_length": len(answer),
	}).Info("Successfully generated answer")

	return answer, nil
}