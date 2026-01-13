package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igorfazlyev/dm/internal/config"
	"github.com/igorfazlyev/dm/internal/database"
	jwtpkg "github.com/igorfazlyev/dm/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Role      string `json:"role" binding:"required,oneof=patient clinic_doctor clinic_manager"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var existingUser database.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Create user
	user := database.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		IsActive:     true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// Create patient profile if role is patient
	if req.Role == "patient" {
		patient := database.Patient{
			UserID:    user.ID,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}
		if err := database.DB.Create(&patient).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create patient profile"})
			return
		}
	}

	// Generate tokens
	tokens, err := jwtpkg.GenerateTokenPair(
		user.ID,
		user.Email,
		user.Role,
		h.cfg.JWT.Secret,
		h.cfg.JWT.AccessTokenTTL,
		h.cfg.JWT.RefreshTokenTTL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":   user,
		"tokens": tokens,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user database.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "account is inactive"})
		return
	}

	// Generate tokens
	tokens, err := jwtpkg.GenerateTokenPair(
		user.ID,
		user.Email,
		user.Role,
		h.cfg.JWT.Secret,
		h.cfg.JWT.AccessTokenTTL,
		h.cfg.JWT.RefreshTokenTTL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":   user,
		"tokens": tokens,
	})
}

func (h *Handler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user database.User
	if err := database.DB.Preload("Patient").Preload("Clinic").Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
