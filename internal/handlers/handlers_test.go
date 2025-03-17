package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"doc-qa-api/internal/handlers"
)

// MockGeminiClient is a mock implementation of the Gemini client.
type MockGeminiClient struct {
	mock.Mock
}

// GenerateAnswer mocks the GenerateAnswer method of the Gemini client.
func (m *MockGeminiClient) GenerateAnswer(document, question string) (string, error) {
	args := m.Called(document, question)
	return args.String(0), args.Error(1)
}

// TestAskHandler_Success tests the /ask endpoint with a valid request.
func TestAskHandler_Success(t *testing.T) {
	// Create a mock Gemini client
	mockClient := new(MockGeminiClient)
	mockClient.On("GenerateAnswer", "test document", "test question").Return("test answer", nil)

	// Create a Gin router
	router := gin.Default()
	router.POST("/ask", handlers.AskHandler(mockClient))

	// Create a test request
	reqBody := `{"document":"test document","question":"test question"}`
	req, _ := http.NewRequest("POST", "/ask", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"answer":"test answer"}`, w.Body.String())
}

// TestAskHandler_InvalidRequest tests the /ask endpoint with an invalid request.
func TestAskHandler_InvalidRequest(t *testing.T) {
	// Create a mock Gemini client
	mockClient := new(MockGeminiClient)

	// Create a Gin router
	router := gin.Default()
	router.POST("/ask", handlers.AskHandler(mockClient))

	// Create a test request with invalid JSON
	reqBody := `{"document":"test document"`
	req, _ := http.NewRequest("POST", "/ask", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request format")
}