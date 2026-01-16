package studies

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/igorfazlyev/dm/internal/config"
	"github.com/igorfazlyev/dm/internal/database"
	"github.com/igorfazlyev/dm/internal/services"
)

type Handler struct {
	cfg              *config.Config
	diagnocatService *services.DiagnocatService
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg:              cfg,
		diagnocatService: services.NewDiagnocatService(),
	}
}

type CreateStudyRequest struct {
	Modality  string `json:"modality" binding:"required"`
	StudyDate string `json:"study_date"`
}

func (h *Handler) CreateStudy(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateStudyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get patient profile
	var patient database.Patient
	if err := database.DB.Where("user_id = ?", userID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient profile not found"})
		return
	}

	study := database.Study{
		PatientID: patient.ID,
		Modality:  req.Modality,
		Status:    "created",
	}

	if req.StudyDate != "" {
		study.StudyDate = &req.StudyDate
	}

	if err := database.DB.Create(&study).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create study"})
		return
	}

	c.JSON(http.StatusCreated, study)
}

func (h *Handler) InitiateDICOMUpload(c *gin.Context) {
	userID := c.GetString("user_id")
	studyID := c.Param("id")

	// Verify study belongs to user
	var study database.Study
	if err := database.DB.Preload("Patient").Where("id = ?", studyID).First(&study).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "study not found"})
		return
	}

	if study.Patient.UserID.String() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Update study status to ready for upload
	study.Status = "uploading"
	if err := database.DB.Save(&study).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update study"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "ready for upload",
		"study_id":        studyID,
		"upload_endpoint": fmt.Sprintf("/api/v1/studies/%s/upload", studyID),
	})
}

func (h *Handler) UploadDICOMFile(c *gin.Context) {
	userID := c.GetString("user_id")
	studyID := c.Param("id")

	// Verify study
	var study database.Study
	if err := database.DB.Preload("Patient").Where("id = ?", studyID).First(&study).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "study not found"})
		return
	}

	if study.Patient.UserID.String() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Get uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}
	defer file.Close()

	// Save file temporarily
	tempDir := os.TempDir()
	tempFilePath := filepath.Join(tempDir, fmt.Sprintf("study_%s_%s", studyID, header.Filename))

	outFile, err := os.Create(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	defer func() {
		outFile.Close()
		os.Remove(tempFilePath) // Clean up temp file
	}()

	if _, err := io.Copy(outFile, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	outFile.Close() // Close before uploading

	// Get or create Diagnocat patient ID
	diagnocatPatientID := study.Patient.DiagnocatPatientID
	if diagnocatPatientID == nil || *diagnocatPatientID == "" {
		newPatientID := fmt.Sprintf("patient-%s", study.Patient.ID.String())
		diagnocatPatientID = &newPatientID

		if err := database.DB.Model(&study.Patient).Update("diagnocat_patient_id", newPatientID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update patient"})
			return
		}
	}

	// Update study status
	study.Status = "processing"
	if err := database.DB.Save(&study).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update study"})
		return
	}

	// Upload to Diagnocat asynchronously
	go func() {
		analysisResp, err := h.diagnocatService.UploadStudy(*diagnocatPatientID, tempFilePath)
		if err != nil {
			fmt.Printf("Failed to upload to Diagnocat: %v\n", err)
			database.DB.Model(&study).Update("status", "failed")
			return
		}

		// Store analysis IDs
		diagnocatStudyUID := analysisResp.UID
		if diagnocatStudyUID == "" {
			diagnocatStudyUID = analysisResp.IDV3
		}

		updates := map[string]interface{}{
			"diagnocat_study_uid": diagnocatStudyUID,
			"status":              "processing",
		}
		database.DB.Model(&study).Updates(updates)

		fmt.Printf("✅ Upload complete. Analysis ID: %s\n", diagnocatStudyUID)
	}()

	c.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded successfully, processing started",
		"filename": header.Filename,
		"status":   "processing",
	})
}

func (h *Handler) GetStudy(c *gin.Context) {
	studyID := c.Param("id")

	var study database.Study
	if err := database.DB.Preload("Patient").Where("id = ?", studyID).First(&study).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "study not found"})
		return
	}

	c.JSON(http.StatusOK, study)
}

func (h *Handler) CheckStudyStatus(c *gin.Context) {
	studyID := c.Param("id")
	userID := c.GetString("user_id")

	var study database.Study
	if err := database.DB.Preload("Patient").Where("id = ?", studyID).First(&study).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "study not found"})
		return
	}

	if study.Patient.UserID.String() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// If study is processing and has Diagnocat study UID, check status
	if study.Status == "processing" && study.DiagnocatStudyUID != nil {
		reportStatus, err := h.diagnocatService.GetAnalysisStatus(*study.DiagnocatStudyUID)
		if err == nil {
			if reportStatus.Complete || reportStatus.Status == "complete" {
				study.Status = "completed"
				study.DiagnocatReportURL = &reportStatus.PDFUrl
				database.DB.Save(&study)
			} else if reportStatus.Status == "error" {
				study.Status = "failed"
				database.DB.Save(&study)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":                   study.ID,
		"status":               study.Status,
		"modality":             study.Modality,
		"created_at":           study.CreatedAt,
		"diagnocat_study_uid":  study.DiagnocatStudyUID,
		"diagnocat_report_url": study.DiagnocatReportURL,
	})
}

func (h *Handler) GetStudyPDF(c *gin.Context) {
	studyID := c.Param("id")
	userID := c.GetString("user_id")

	var study database.Study
	if err := database.DB.Preload("Patient").Where("id = ?", studyID).First(&study).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "study not found"})
		return
	}

	if study.Patient.UserID.String() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if study.DiagnocatStudyUID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no diagnocat report available"})
		return
	}

	// Download PDF to temp file
	tempDir := os.TempDir()
	pdfPath := filepath.Join(tempDir, fmt.Sprintf("report_%s.pdf", studyID))

	if err := h.diagnocatService.DownloadReportPDF(*study.DiagnocatStudyUID, pdfPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to download PDF: %v", err)})
		return
	}
	defer os.Remove(pdfPath)

	// Serve the PDF
	c.FileAttachment(pdfPath, fmt.Sprintf("report_%s.pdf", studyID))
}
