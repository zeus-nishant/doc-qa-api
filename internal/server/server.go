package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"doc-qa-api/internal/gemini"
	"doc-qa-api/internal/handlers"
)

// Router is the Gin router used by the server.
var Router *gin.Engine

// Start starts the HTTP server.
func Start() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get the Gemini API key from environment variables
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable missing")
	}

	// Initialize Gemini client
	geminiClient, err := gemini.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}

	// Create a new Gin router
	Router = gin.Default()

	// Define routes
	Router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Document Q&A API!",
		})
	})

	// Add the /ask endpoint
	Router.POST("/ask", handlers.AskHandler(geminiClient))

	// Start the server
	log.Println("Server running on :8080")
	Router.Run(":8080")
}