package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"doc-qa-api/internal/services"
)

type MockGeminiClient struct {
	mock.Mock
}

func (m *MockGeminiClient) GenerateAnswer(document, question string) (string, error) {
	args := m.Called(document, question)
	return args.String(0), args.Error(1)
}

func TestGeminiService_GenerateAnswer(t *testing.T) {
	mockClient := new(MockGeminiClient)
	mockClient.On("GenerateAnswer", "test document", "test question").Return("test answer", nil)

	service := services.NewGeminiService(mockClient)
	answer, err := service.GenerateAnswer("test document", "test question")

	assert.NoError(t, err)
	assert.Equal(t, "test answer", answer)
}