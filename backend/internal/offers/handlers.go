package offers

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

type CreateOfferRequestReq struct {
	PlanVersionID     uuid.UUID `json:"plan_version_id" binding:"required"`
	SelectedItemIDs   []string  `json:"selected_item_ids" binding:"required"`
	PreferredCity     string    `json:"preferred_city"`
	PreferredDistrict string    `json:"preferred_district"`
	PriceSegment      string    `json:"price_segment"`
}

// CreateOfferRequest creates a request for offers from clinics
func (h *Handler) CreateOfferRequest(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)

	var req CreateOfferRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get patient
	var patient database.Patient
	if err := database.DB.Where("user_id = ?", userID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	// Verify plan exists and belongs to patient
	var planVersion database.PlanVersion
	if err := database.DB.Preload("Study").Where("id = ?", req.PlanVersionID).First(&planVersion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}

	if planVersion.Study.PatientID != patient.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// Create offer request
	offerRequest := database.OfferRequest{
		PatientID:         patient.ID,
		PlanVersionID:     req.PlanVersionID,
		SelectedItemIDs:   req.SelectedItemIDs, // Direct assignment
		PreferredCity:     req.PreferredCity,
		PreferredDistrict: req.PreferredDistrict,
		PriceSegment:      req.PriceSegment,
		Status:            "open",
	}

	if err := database.DB.Create(&offerRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create offer request"})
		return
	}

	c.JSON(http.StatusCreated, offerRequest)
}

// GetMyOfferRequests returns all offer requests for the current patient
func (h *Handler) GetMyOfferRequests(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)

	var patient database.Patient
	if err := database.DB.Where("user_id = ?", userID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	var offerRequests []database.OfferRequest
	if err := database.DB.Preload("PlanVersion.PlanItems").Preload("Offers.Clinic").
		Where("patient_id = ?", patient.ID).
		Order("created_at DESC").
		Find(&offerRequests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch offer requests"})
		return
	}

	c.JSON(http.StatusOK, offerRequests)
}

// GetOfferRequest returns a single offer request with all offers
func (h *Handler) GetOfferRequest(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)
	userRole, _ := rbac.GetUserRole(c)
	requestIDParam := c.Param("id")

	requestID, err := uuid.Parse(requestIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request ID"})
		return
	}

	var offerRequest database.OfferRequest
	query := database.DB.Preload("PlanVersion.PlanItems").Preload("Offers.Clinic")

	if err := query.Where("id = ?", requestID).First(&offerRequest).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "offer request not found"})
		return
	}

	// Check access
	if userRole == rbac.RolePatient {
		var patient database.Patient
		if err := database.DB.Where("user_id = ?", userID).First(&patient).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		if offerRequest.PatientID != patient.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
	}

	c.JSON(http.StatusOK, offerRequest)
}

// Clinic creates an offer

type CreateOfferReq struct {
	OfferRequestID   uuid.UUID `json:"offer_request_id" binding:"required"`
	TotalPrice       float64   `json:"total_price" binding:"required"`
	DiscountPercent  float64   `json:"discount_percent"`
	HasInstallment   bool      `json:"has_installment"`
	InstallmentTerms string    `json:"installment_terms"`
	SpecialOffer     string    `json:"special_offer"`
	EstimatedDays    int       `json:"estimated_days"`
}

// CreateOffer allows a clinic to submit an offer
func (h *Handler) CreateOffer(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)

	var req CreateOfferReq
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

	if !clinic.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "clinic is not active"})
		return
	}

	// Verify offer request exists
	var offerRequest database.OfferRequest
	if err := database.DB.Where("id = ?", req.OfferRequestID).First(&offerRequest).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "offer request not found"})
		return
	}

	if offerRequest.Status != "open" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "offer request is closed"})
		return
	}

	// Create offer
	offer := database.Offer{
		OfferRequestID:   req.OfferRequestID,
		ClinicID:         clinic.ID,
		TotalPrice:       req.TotalPrice,
		DiscountPercent:  req.DiscountPercent,
		HasInstallment:   req.HasInstallment,
		InstallmentTerms: req.InstallmentTerms,
		SpecialOffer:     req.SpecialOffer,
		EstimatedDays:    req.EstimatedDays,
		Status:           "pending",
	}

	if err := database.DB.Create(&offer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create offer"})
		return
	}

	c.JSON(http.StatusCreated, offer)
}

// GetMyOffers returns all offers created by the current clinic
func (h *Handler) GetMyOffers(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)

	var clinic database.Clinic
	if err := database.DB.Where("user_id = ?", userID).First(&clinic).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "clinic not found"})
		return
	}

	var offers []database.Offer
	if err := database.DB.Preload("OfferRequest.PlanVersion.PlanItems").
		Where("clinic_id = ?", clinic.ID).
		Order("created_at DESC").
		Find(&offers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch offers"})
		return
	}

	c.JSON(http.StatusOK, offers)
}

// AcceptOffer allows a patient to accept an offer
func (h *Handler) AcceptOffer(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)
	offerIDParam := c.Param("id")

	offerID, err := uuid.Parse(offerIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offer ID"})
		return
	}

	var patient database.Patient
	if err := database.DB.Where("user_id = ?", userID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	var offer database.Offer
	if err := database.DB.Preload("OfferRequest").Where("id = ?", offerID).First(&offer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "offer not found"})
		return
	}

	// Verify patient owns the offer request
	if offer.OfferRequest.PatientID != patient.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if offer.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "offer already processed"})
		return
	}

	// Start transaction
	tx := database.DB.Begin()

	// Update offer status
	offer.Status = "accepted"
	if err := tx.Save(&offer).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to accept offer"})
		return
	}

	// Create order
	order := database.Order{
		OfferID:   offer.ID,
		PatientID: patient.ID,
		ClinicID:  offer.ClinicID,
		Status:    "new",
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "offer accepted",
		"order":   order,
	})
}
