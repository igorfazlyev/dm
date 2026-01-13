package orders

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

// GetMyOrders returns orders for the current user (patient or clinic)
func (h *Handler) GetMyOrders(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)
	userRole, _ := rbac.GetUserRole(c)

	var orders []database.Order
	query := database.DB.Preload("Offer.OfferRequest.PlanVersion.PlanItems").
		Preload("Patient").
		Preload("Clinic").
		Order("created_at DESC")

	if userRole == rbac.RolePatient {
		var patient database.Patient
		if err := database.DB.Where("user_id = ?", userID).First(&patient).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
			return
		}
		query = query.Where("patient_id = ?", patient.ID)
	} else if userRole == rbac.RoleClinicDoctor || userRole == rbac.RoleClinicManager {
		var clinic database.Clinic
		if err := database.DB.Where("user_id = ?", userID).First(&clinic).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "clinic not found"})
			return
		}
		query = query.Where("clinic_id = ?", clinic.ID)
	}

	if err := query.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder returns a single order by ID
func (h *Handler) GetOrder(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)
	userRole, _ := rbac.GetUserRole(c)
	orderIDParam := c.Param("id")

	orderID, err := uuid.Parse(orderIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	var order database.Order
	query := database.DB.Preload("Offer.OfferRequest.PlanVersion.PlanItems").
		Preload("Patient").
		Preload("Clinic").
		Preload("Slots")

	if err := query.Where("id = ?", orderID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Check access
	if userRole == rbac.RolePatient {
		var patient database.Patient
		if err := database.DB.Where("user_id = ?", userID).First(&patient).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if order.PatientID != patient.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
	} else if userRole == rbac.RoleClinicDoctor || userRole == rbac.RoleClinicManager {
		var clinic database.Clinic
		if err := database.DB.Where("user_id = ?", userID).First(&clinic).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if order.ClinicID != clinic.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
	}

	c.JSON(http.StatusOK, order)
}

type UpdateOrderStatusReq struct {
	Status string `json:"status" binding:"required"`
}

// UpdateOrderStatus updates order status (clinic only)
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	userID, _ := rbac.GetUserID(c)
	orderIDParam := c.Param("id")

	orderID, err := uuid.Parse(orderIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	var req UpdateOrderStatusReq
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

	// Get order
	var order database.Order
	if err := database.DB.Where("id = ? AND clinic_id = ?", orderID, clinic.ID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	order.Status = req.Status
	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order"})
		return
	}

	c.JSON(http.StatusOK, order)
}
