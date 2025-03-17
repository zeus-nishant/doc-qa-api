package routes

import (
	"github.com/gin-gonic/gin"
	"doc-qa-api/internal/controllers"
)

func SetupRoutes(router *gin.Engine, questionController *controllers.QuestionController) {
	api := router.Group("/api")
	{
		api.POST("/ask", questionController.AskQuestion)
	}
}