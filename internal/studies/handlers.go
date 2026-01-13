package studies

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/igorfazlyev/dm/internal/database"
	"github.com/igorfazlyev/dm/internal/rbac"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

type CreateStudyRequest struct {
	Modality  string `json:"modality"`
	StudyDate string `json:"study_date"` // YYYY-MM-DD
}

// CreateStudy creates a new study for the current patient
func (h *Handler) CreateStudy(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)

	var req CreateStudyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get patient
	var patient database.Patient
	if err := database.DB.Where("user_id = ?", userID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient profile not found"})
		return
	}

	// Parse study date if provided
	var studyDate *time.Time
	if req.StudyDate != "" {
		parsed, err := time.Parse("2006-01-02", req.StudyDate)
		if err == nil {
			studyDate = &parsed
		}
	}

	// Create study
	study := database.Study{
		PatientID: patient.ID,
		Status:    "created",
		Modality:  req.Modality,
		StudyDate: studyDate,
	}

	if err := database.DB.Create(&study).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create study"})
		return
	}

	// Create audit log
	auditLog := database.AuditLog{
		UserID:     &userID,
		Action:     "create_study",
		EntityType: "study",
		EntityID:   &study.ID,
		Details:    map[string]any{"modality": req.Modality},
		IPAddress:  c.ClientIP(),
	}
	database.DB.Create(&auditLog)

	c.JSON(http.StatusCreated, study)
}

// GetStudy returns a single study by ID
func (h *Handler) GetStudy(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)
	userRole, _ := rbac.GetUserRole(c)
	studyIDParam := c.Param("id")

	studyID, err := uuid.Parse(studyIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid study ID"})
		return
	}

	var study database.Study
	query := database.DB.Preload("Patient").Preload("PlanVersions.PlanItems")

	// Patients can only see their own studies
	if userRole == rbac.RolePatient {
		var patient database.Patient
		if err := database.DB.Where("user_id = ?", userID).First(&patient).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		query = query.Where("patient_id = ?", patient.ID)
	}

	if err := query.Where("id = ?", studyID).First(&study).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "study not found"})
		return
	}

	c.JSON(http.StatusOK, study)
}

// UpdateStudyStatus updates study status (for worker/admin use)
func (h *Handler) UpdateStudyStatus(c *gin.Context) {
	studyIDParam := c.Param("id")
	studyID, err := uuid.Parse(studyIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid study ID"})
		return
	}

	var req struct {
		Status       string         `json:"status" binding:"required"`
		ErrorMessage string         `json:"error_message"`
		ResultJSON   map[string]any `json:"result_json"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var study database.Study
	if err := database.DB.Where("id = ?", studyID).First(&study).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "study not found"})
		return
	}

	study.Status = req.Status
	if req.ErrorMessage != "" {
		study.ErrorMessage = req.ErrorMessage
	}
	if req.ResultJSON != nil {
		study.DiagnocatResultJSON = req.ResultJSON
	}
	if req.Status == "completed" {
		now := time.Now()
		study.CompletedAt = &now
	}

	if err := database.DB.Save(&study).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update study"})
		return
	}

	c.JSON(http.StatusOK, study)
}

// GetStudyPDF returns the PDF report for a study
func (h *Handler) GetStudyPDF(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)
	userRole, _ := rbac.GetUserRole(c)
	studyIDParam := c.Param("id")

	studyID, err := uuid.Parse(studyIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid study ID"})
		return
	}

	var study database.Study
	query := database.DB

	// Patients can only see their own studies
	if userRole == rbac.RolePatient {
		var patient database.Patient
		if err := database.DB.Where("user_id = ?", userID).First(&patient).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		query = query.Where("patient_id = ?", patient.ID)
	}

	if err := query.Where("id = ?", studyID).First(&study).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "study not found"})
		return
	}

	if len(study.DiagnocatReportPDF) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "PDF not available"})
		return
	}

	c.Data(http.StatusOK, "application/pdf", study.DiagnocatReportPDF)
}
