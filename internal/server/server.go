package server

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"doc-qa-api/internal/gemini"
	"doc-qa-api/internal/routes"
	"doc-qa-api/internal/services"
	"doc-qa-api/internal/controllers"
)

var Router *gin.Engine

func Start() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize Gemini client
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable missing")
	}
	geminiClient, err := gemini.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}

	// Initialize services
	geminiService := services.NewGeminiService(geminiClient)

	// Initialize controllers
	questionController := controllers.NewQuestionController(geminiService)

	// Create a new Gin router
	Router = gin.Default()

	// Set up routes
	routes.SetupRoutes(Router, questionController)

	// Start the server
	log.Println("Server running on :8080")
	Router.Run(":8080")
}