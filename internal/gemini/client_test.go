package gemini_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGeminiClient is a mock implementation of the ClientInterface.
type MockGeminiClient struct {
	mock.Mock
}

// GenerateAnswer mocks the GenerateAnswer method.
func (m *MockGeminiClient) GenerateAnswer(document, question string) (string, error) {
	args := m.Called(document, question)
	return args.String(0), args.Error(1)
}

// TestGenerateAnswer_Success tests the GenerateAnswer method with a successful response.
func TestGenerateAnswer_Success(t *testing.T) {
	// Create a mock Gemini client
	mockClient := new(MockGeminiClient)
	mockClient.On("GenerateAnswer", "test document", "test question").Return("test answer", nil)

	// Call the method
	answer, err := mockClient.GenerateAnswer("test document", "test question")

	// Assert the result
	assert.NoError(t, err)
	assert.Equal(t, "test answer", answer)
}

// TestGenerateAnswer_Error tests the GenerateAnswer method with an error response.
func TestGenerateAnswer_Error(t *testing.T) {
	// Create a mock Gemini client
	mockClient := new(MockGeminiClient)
	mockClient.On("GenerateAnswer", "test document", "test question").Return("", assert.AnError)

	// Call the method
	_, err := mockClient.GenerateAnswer("test document", "test question")

	// Assert the result
	assert.Error(t, err)
}