package server

import (
	"log"
	"net/http"
	"os"

	"doc-qa-api/internal/controllers"
	"doc-qa-api/internal/gemini"
	"doc-qa-api/internal/routes"
	"doc-qa-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	pdfService := services.NewPDFService(geminiClient)

	// Initialize controllers
	questionController := controllers.NewQuestionController(geminiService)
	uploadController := controllers.NewUploadController(pdfService)

	// Create a new Gin router
	Router = gin.Default()

	// Apply CORS middleware
	Router.Use(corsMiddleware())

	// Set up routes
	routes.SetupRoutes(Router, questionController, uploadController)

	// Start the server
	log.Println("Server running on :8080")
	Router.Run(":8080")
}

// corsMiddleware handles CORS settings properly
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}