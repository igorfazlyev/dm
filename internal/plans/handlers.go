package plans

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/igorfazlyev/dm/internal/database"
	"github.com/igorfazlyev/dm/internal/rbac"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

type CreatePlanRequest struct {
	StudyID uuid.UUID      `json:"study_id" binding:"required"`
	Source  string         `json:"source"` // diagnocat, manual, modified
	Items   []PlanItemData `json:"items" binding:"required,dive"`
}

type PlanItemData struct {
	ToothNumber   *int           `json:"tooth_number"`
	Specialty     string         `json:"specialty" binding:"required"`
	ProcedureCode string         `json:"procedure_code"`
	ProcedureName string         `json:"procedure_name" binding:"required"`
	Diagnosis     string         `json:"diagnosis"`
	Quantity      int            `json:"quantity"`
	Notes         string         `json:"notes"`
	Metadata      map[string]any `json:"metadata"`
}

// CreatePlan creates a new treatment plan version
func (h *Handler) CreatePlan(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)
	userRole, _ := rbac.GetUserRole(c)

	var req CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify study exists and user has access
	var study database.Study
	query := database.DB

	if userRole == rbac.RolePatient {
		var patient database.Patient
		if err := database.DB.Where("user_id = ?", userID).First(&patient).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		query = query.Where("patient_id = ?", patient.ID)
	}

	if err := query.Where("id = ?", req.StudyID).First(&study).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "study not found"})
		return
	}

	// Get next version number
	var count int64
	database.DB.Model(&database.PlanVersion{}).Where("study_id = ?", req.StudyID).Count(&count)

	// Create plan version
	planVersion := database.PlanVersion{
		StudyID: req.StudyID,
		Version: int(count) + 1,
		Source:  req.Source,
	}

	if err := database.DB.Create(&planVersion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create plan version"})
		return
	}

	// Create plan items
	for _, itemData := range req.Items {
		item := database.PlanItem{
			PlanVersionID: planVersion.ID,
			ToothNumber:   itemData.ToothNumber,
			Specialty:     itemData.Specialty,
			ProcedureCode: itemData.ProcedureCode,
			ProcedureName: itemData.ProcedureName,
			Diagnosis:     itemData.Diagnosis,
			Quantity:      itemData.Quantity,
			Notes:         itemData.Notes,
			Metadata:      itemData.Metadata,
		}

		if item.Quantity == 0 {
			item.Quantity = 1
		}

		if err := database.DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create plan item"})
			return
		}
	}

	// Reload with items
	database.DB.Preload("PlanItems").First(&planVersion, planVersion.ID)

	c.JSON(http.StatusCreated, planVersion)
}

// GetPlan returns a treatment plan by ID
func (h *Handler) GetPlan(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)
	userRole, _ := rbac.GetUserRole(c)
	planIDParam := c.Param("id")

	planID, err := uuid.Parse(planIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	var planVersion database.PlanVersion
	query := database.DB.Preload("PlanItems").Preload("Study.Patient")

	if err := query.Where("id = ?", planID).First(&planVersion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}

	// Check access for patients
	if userRole == rbac.RolePatient {
		var patient database.Patient
		if err := database.DB.Where("user_id = ?", userID).First(&patient).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		if planVersion.Study.PatientID != patient.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
	}

	c.JSON(http.StatusOK, planVersion)
}

// GetPlansByStudy returns all plan versions for a study
func (h *Handler) GetPlansByStudy(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)
	userRole, _ := rbac.GetUserRole(c)
	studyIDParam := c.Param("study_id")

	studyID, err := uuid.Parse(studyIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid study ID"})
		return
	}

	// Verify study access
	var study database.Study
	query := database.DB

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

	// Get plan versions
	var planVersions []database.PlanVersion
	if err := database.DB.Preload("PlanItems").
		Where("study_id = ?", studyID).
		Order("version DESC").
		Find(&planVersions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch plans"})
		return
	}

	c.JSON(http.StatusOK, planVersions)
}

// GetEstimate calculates price estimates by specialty
func (h *Handler) GetEstimate(c *gin.Context) {
	planIDParam := c.Param("id")
	planID, err := uuid.Parse(planIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	var planVersion database.PlanVersion
	if err := database.DB.Preload("PlanItems").Where("id = ?", planID).First(&planVersion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}

	// Group by specialty
	estimates := map[string]map[string]float64{
		"therapy":      {"min": 0, "max": 0, "count": 0},
		"orthopedics":  {"min": 0, "max": 0, "count": 0},
		"surgery":      {"min": 0, "max": 0, "count": 0},
		"hygiene":      {"min": 0, "max": 0, "count": 0},
		"periodontics": {"min": 0, "max": 0, "count": 0},
	}

	for _, item := range planVersion.PlanItems {
		if specialty, exists := estimates[item.Specialty]; exists {
			specialty["count"]++
			// This is placeholder - real implementation would look up prices
			// For now, just count items
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"plan_id":   planID,
		"estimates": estimates,
		"note":      "Estimates require clinic pricelist integration",
	})
}
