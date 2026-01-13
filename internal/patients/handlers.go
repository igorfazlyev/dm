package patients

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igorfazlyev/dm/internal/database"
	"github.com/igorfazlyev/dm/internal/rbac"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

type CreatePatientRequest struct {
	FirstName             string `json:"first_name" binding:"required"`
	LastName              string `json:"last_name" binding:"required"`
	DateOfBirth           string `json:"date_of_birth"` // YYYY-MM-DD
	Phone                 string `json:"phone"`
	PreferredCity         string `json:"preferred_city"`
	PreferredDistrict     string `json:"preferred_district"`
	PreferredPriceSegment string `json:"preferred_price_segment"`
}

type UpdatePatientRequest struct {
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	Phone                 string `json:"phone"`
	PreferredCity         string `json:"preferred_city"`
	PreferredDistrict     string `json:"preferred_district"`
	PreferredPriceSegment string `json:"preferred_price_segment"`
}

// GetMyProfile returns the current patient's profile
func (h *Handler) GetMyProfile(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)

	var patient database.Patient
	if err := database.DB.Where("user_id = ?", userID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient profile not found"})
		return
	}

	c.JSON(http.StatusOK, patient)
}

// UpdateMyProfile updates the current patient's profile
func (h *Handler) UpdateMyProfile(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)

	var req UpdatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var patient database.Patient
	if err := database.DB.Where("user_id = ?", userID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient profile not found"})
		return
	}

	// Update fields
	if req.FirstName != "" {
		patient.FirstName = req.FirstName
	}
	if req.LastName != "" {
		patient.LastName = req.LastName
	}
	if req.Phone != "" {
		patient.Phone = req.Phone
	}
	if req.PreferredCity != "" {
		patient.PreferredCity = req.PreferredCity
	}
	if req.PreferredDistrict != "" {
		patient.PreferredDistrict = req.PreferredDistrict
	}
	if req.PreferredPriceSegment != "" {
		patient.PreferredPriceSegment = req.PreferredPriceSegment
	}

	if err := database.DB.Save(&patient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, patient)
}

// GetMyStudies returns all studies for the current patient
func (h *Handler) GetMyStudies(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)

	var patient database.Patient
	if err := database.DB.Where("user_id = ?", userID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	var studies []database.Study
	if err := database.DB.Where("patient_id = ?", patient.ID).
		Order("created_at DESC").
		Find(&studies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch studies"})
		return
	}

	c.JSON(http.StatusOK, studies)
}
