package handlers

import (
	"net/http"

	"doc-qa-api/internal/gemini"

	"github.com/gin-gonic/gin"
)

type QuestionRequest struct {
	Document string `json:"document" binding:"required"`
	Question string `json:"question" binding:"required"`
}

func AskHandler(client gemini.ClientInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req QuestionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		answer, err := client.GenerateAnswer(req.Document, req.Question)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "AI service failure"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"answer": answer})
	}
}
