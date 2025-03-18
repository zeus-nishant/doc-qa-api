package services

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"doc-qa-api/internal/gemini"
)

var logger = logrus.New()

// PDFService processes PDFs and interacts with Gemini API
type PDFService struct {
	geminiClient gemini.ClientInterface
}

// NewPDFService creates a new instance of PDFService
func NewPDFService(geminiClient gemini.ClientInterface) *PDFService {
	return &PDFService{geminiClient: geminiClient}
}

func (s *PDFService) ProcessPDF(file io.Reader, question string) (string, error) {
	// Save the uploaded file to a temporary location
	tempFile, err := os.CreateTemp("", "uploaded-*.pdf")
	if err != nil {
		logger.Error("Failed to create temp file: ", err)
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, file); err != nil {
		logger.Error("Failed to save uploaded file: ", err)
		return "", fmt.Errorf("failed to save uploaded file: %v", err)
	}

	tempFile.Close()

	// Try to extract text directly using pdftotext
	extractedText, err := extractPDFText(tempFile.Name())
	if err != nil {
		logger.Error("Failed to extract text directly: ", err)
	}

	// If no text was extracted or extraction failed, try OCR
	if extractedText == "" || err != nil {
		ocrText, ocrErr := performOCR(tempFile.Name())
		if ocrErr != nil {
			logger.Error("OCR processing failed: ", ocrErr)
			return "", fmt.Errorf("both text extraction and OCR failed: %v, OCR error: %v", err, ocrErr)
		}
		extractedText = ocrText
	}

	if extractedText == "" {
		return "Unable to extract any text from the provided document. The document might be corrupted or contain only images without recognizable text.", nil
	}

	// Send extracted text to Gemini API
	answer, err := s.geminiClient.GenerateAnswer(extractedText, question)
	if err != nil {
		logger.Error("Failed to generate answer from Gemini: ", err)
		return "", fmt.Errorf("failed to generate answer from Gemini: %v", err)
	}

	return answer, nil
}

// Helper function to extract text from PDF using pdftotext
func extractPDFText(filename string) (string, error) {
	cmd := exec.Command("pdftotext", filename, "-")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// Helper function to perform OCR on PDF using Tesseract
func performOCR(pdfPath string) (string, error) {
	tempDir, err := os.MkdirTemp("", "ocr-images-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cmd := exec.Command("pdftoppm", "-png", pdfPath, filepath.Join(tempDir, "page"))
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to convert PDF to images: %v", err)
	}

	imageFiles, err := filepath.Glob(filepath.Join(tempDir, "page-*.png"))
	if err != nil {
		return "", fmt.Errorf("failed to list image files: %v", err)
	}

	var allText strings.Builder

	for _, imgFile := range imageFiles {
		outputBase := filepath.Join(tempDir, fmt.Sprintf("out_%s", filepath.Base(imgFile)))
		cmd := exec.Command("tesseract", imgFile, outputBase)
		if err := cmd.Run(); err != nil {
			logger.Error(fmt.Sprintf("OCR failed for %s: %v", imgFile, err))
			continue
		}

		txtFile := outputBase + ".txt"
		textBytes, err := os.ReadFile(txtFile)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to read OCR result for %s: %v", imgFile, err))
			continue
		}

		allText.WriteString(string(textBytes))
		allText.WriteString("\n")
	}

	return allText.String(), nil
}