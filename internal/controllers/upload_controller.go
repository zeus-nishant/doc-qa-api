package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"doc-qa-api/internal/services"
)

type UploadController struct {
	pdfService *services.PDFService
}

func NewUploadController(pdfService *services.PDFService) *UploadController {
	return &UploadController{
		pdfService: pdfService,
	}
}

func (uc *UploadController) UploadHandler(c *gin.Context) {
	// Get the uploaded file
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read file"})
		return
	}
	defer file.Close()

	// Get the question from the form
	question := c.PostForm("question")
	if question == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Question is required"})
		return
	}

	// Process PDF in the service layer
	answer, err := uc.pdfService.ProcessPDF(file, question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process PDF"})
		return
	}

	// Return the answer
	c.JSON(http.StatusOK, gin.H{"answer": answer})
}
