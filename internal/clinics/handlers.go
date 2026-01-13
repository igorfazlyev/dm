package clinics

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

type CreateClinicRequest struct {
	Name            string `json:"name" binding:"required"`
	LegalName       string `json:"legal_name"`
	LicenseNumber   string `json:"license_number" binding:"required"`
	YearEstablished int    `json:"year_established"`
	City            string `json:"city" binding:"required"`
	District        string `json:"district"`
	Address         string `json:"address" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	Email           string `json:"email"`
	Website         string `json:"website"`
	PriceSegment    string `json:"price_segment"`
}

type UpdateClinicRequest struct {
	Name         string `json:"name"`
	District     string `json:"district"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Website      string `json:"website"`
	PriceSegment string `json:"price_segment"`
}

// CreateClinic creates a new clinic profile
func (h *Handler) CreateClinic(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)

	var req CreateClinicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if clinic already exists for this user
	var existing database.Clinic
	if err := database.DB.Where("user_id = ?", userID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "clinic profile already exists"})
		return
	}

	clinic := database.Clinic{
		UserID:          userID,
		Name:            req.Name,
		LegalName:       req.LegalName,
		LicenseNumber:   req.LicenseNumber,
		YearEstablished: req.YearEstablished,
		City:            req.City,
		District:        req.District,
		Address:         req.Address,
		Phone:           req.Phone,
		Email:           req.Email,
		Website:         req.Website,
		PriceSegment:    req.PriceSegment,
		IsActive:        false, // Requires admin approval
	}

	if err := database.DB.Create(&clinic).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create clinic"})
		return
	}

	c.JSON(http.StatusCreated, clinic)
}

// GetMyClinic returns the current user's clinic
func (h *Handler) GetMyClinic(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)

	var clinic database.Clinic
	if err := database.DB.Where("user_id = ?", userID).First(&clinic).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "clinic not found"})
		return
	}

	c.JSON(http.StatusOK, clinic)
}

// UpdateMyClinic updates the current user's clinic
func (h *Handler) UpdateMyClinic(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)

	var req UpdateClinicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var clinic database.Clinic
	if err := database.DB.Where("user_id = ?", userID).First(&clinic).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "clinic not found"})
		return
	}

	// Update fields
	if req.Name != "" {
		clinic.Name = req.Name
	}
	if req.District != "" {
		clinic.District = req.District
	}
	if req.Address != "" {
		clinic.Address = req.Address
	}
	if req.Phone != "" {
		clinic.Phone = req.Phone
	}
	if req.Email != "" {
		clinic.Email = req.Email
	}
	if req.Website != "" {
		clinic.Website = req.Website
	}
	if req.PriceSegment != "" {
		clinic.PriceSegment = req.PriceSegment
	}

	if err := database.DB.Save(&clinic).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update clinic"})
		return
	}

	c.JSON(http.StatusOK, clinic)
}

// ListClinics returns all active clinics (with filters)
func (h *Handler) ListClinics(c *gin.Context) {
	city := c.Query("city")
	district := c.Query("district")
	priceSegment := c.Query("price_segment")

	query := database.DB.Where("is_active = ?", true)

	if city != "" {
		query = query.Where("city = ?", city)
	}
	if district != "" {
		query = query.Where("district = ?", district)
	}
	if priceSegment != "" {
		query = query.Where("price_segment = ?", priceSegment)
	}

	var clinics []database.Clinic
	if err := query.Find(&clinics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch clinics"})
		return
	}

	c.JSON(http.StatusOK, clinics)
}

// GetClinic returns a single clinic by ID
func (h *Handler) GetClinic(c *gin.Context) {
	clinicIDParam := c.Param("id")
	clinicID, err := uuid.Parse(clinicIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}

	var clinic database.Clinic
	if err := database.DB.Where("id = ?", clinicID).First(&clinic).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "clinic not found"})
		return
	}

	c.JSON(http.StatusOK, clinic)
}

// Pricelist handlers

type AddPriceItemRequest struct {
	Specialty     string  `json:"specialty" binding:"required"`
	ProcedureCode string  `json:"procedure_code"`
	ProcedureName string  `json:"procedure_name" binding:"required"`
	PriceFrom     float64 `json:"price_from" binding:"required"`
	PriceTo       float64 `json:"price_to"`
}

// AddPriceItem adds an item to clinic's pricelist
func (h *Handler) AddPriceItem(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)

	var req AddPriceItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get clinic
	var clinic database.Clinic
	if err := database.DB.Where("user_id = ?", userID).First(&clinic).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "clinic not found"})
		return
	}

	priceItem := database.PriceListItem{
		ClinicID:      clinic.ID,
		Specialty:     req.Specialty,
		ProcedureCode: req.ProcedureCode,
		ProcedureName: req.ProcedureName,
		PriceFrom:     req.PriceFrom,
		PriceTo:       req.PriceTo,
		IsActive:      true,
	}

	if err := database.DB.Create(&priceItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add price item"})
		return
	}

	c.JSON(http.StatusCreated, priceItem)
}

// GetMyPricelist returns the current clinic's pricelist
func (h *Handler) GetMyPricelist(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)

	var clinic database.Clinic
	if err := database.DB.Where("user_id = ?", userID).First(&clinic).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "clinic not found"})
		return
	}

	var priceList []database.PriceListItem
	if err := database.DB.Where("clinic_id = ? AND is_active = ?", clinic.ID, true).
		Order("specialty, procedure_name").
		Find(&priceList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pricelist"})
		return
	}

	c.JSON(http.StatusOK, priceList)
}

// DeletePriceItem soft deletes a price item
func (h *Handler) DeletePriceItem(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)
	itemIDParam := c.Param("id")

	itemID, err := uuid.Parse(itemIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	// Get clinic
	var clinic database.Clinic
	if err := database.DB.Where("user_id = ?", userID).First(&clinic).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "clinic not found"})
		return
	}

	// Verify ownership and deactivate
	result := database.DB.Model(&database.PriceListItem{}).
		Where("id = ? AND clinic_id = ?", itemID, clinic.ID).
		Update("is_active", false)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete item"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "price item deleted"})
}
