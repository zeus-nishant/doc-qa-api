package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"doc-qa-api/internal/controllers"
)

// MockGeminiService is a mock implementation of the GeminiService interface
type MockGeminiService struct {
	mock.Mock
}

func (m *MockGeminiService) GenerateAnswer(document, question string) (string, error) {
	args := m.Called(document, question)
	return args.String(0), args.Error(1)
}

func TestQuestionController_AskQuestion(t *testing.T) {
	// Create a mock Gemini service
	mockService := new(MockGeminiService)
	mockService.On("GenerateAnswer", "test document", "test question").Return("test answer", nil)

	// Create the controller with the mock service
	controller := controllers.NewQuestionController(mockService)

	// Set up the Gin router
	router := gin.Default()
	router.POST("/ask", controller.AskQuestion)

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