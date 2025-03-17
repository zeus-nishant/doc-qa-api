package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"doc-qa-api/internal/server"
)

func TestStart(t *testing.T) {
	// Load environment variables from .env file
	err := godotenv.Load("../../.env") // Path to the .env file from the test directory
	if err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	// Start the server in a goroutine
	go server.Start()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	// Create a test request
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	resp := httptest.NewRecorder()

	// Serve the request using the Gin router
	server.Router.ServeHTTP(resp, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, resp.Code)

	// Parse the JSON response
	var response map[string]string
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert the response message
	assert.Equal(t, "Welcome to the Document Q&A API!", response["message"])
}