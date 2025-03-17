package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"doc-qa-api/internal/services"
)

type QuestionController struct {
	geminiService services.GeminiService
}

func NewQuestionController(geminiService services.GeminiService) *QuestionController {
	return &QuestionController{geminiService: geminiService}
}

func (qc *QuestionController) AskQuestion(c *gin.Context) {
	var requestBody struct {
		Document string `json:"document" binding:"required"`
		Question string `json:"question" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	answer, err := qc.geminiService.GenerateAnswer(requestBody.Document, requestBody.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI service failure"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer": answer})
}