package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/igorfazlyev/dm/internal/auth"
	"github.com/igorfazlyev/dm/internal/clinics"
	"github.com/igorfazlyev/dm/internal/config"
	"github.com/igorfazlyev/dm/internal/database"
	"github.com/igorfazlyev/dm/internal/offers"
	"github.com/igorfazlyev/dm/internal/orders"
	"github.com/igorfazlyev/dm/internal/patients"
	"github.com/igorfazlyev/dm/internal/plans"
	"github.com/igorfazlyev/dm/internal/rbac"
	"github.com/igorfazlyev/dm/internal/studies"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migrations completed successfully")

	// Initialize Gin
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		//AllowOrigins:     cfg.Server.AllowedOrigins,
		//AllowOrigins: []string{
		//	"http://localhost:3000",
		//	"https://*.app.github.dev", // This allows all GitHub Codespaces
		//},
		AllowAllOrigins:  true, 
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize handlers
	authHandler := auth.NewHandler(cfg)
	patientsHandler := patients.NewHandler()
	studiesHandler := studies.NewHandler()
	plansHandler := plans.NewHandler()
	clinicsHandler := clinics.NewHandler()
	offersHandler := offers.NewHandler()
	ordersHandler := orders.NewHandler()

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		public.GET("/clinics", clinicsHandler.ListClinics)
		public.GET("/clinics/:id", clinicsHandler.GetClinic)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(rbac.AuthMiddleware(cfg))
	{
		// Auth
		protected.GET("/auth/me", authHandler.Me)

		// Patient routes
		patientRoutes := protected.Group("/patient")
		patientRoutes.Use(rbac.RequireRole(rbac.RolePatient))
		{
			patientRoutes.GET("/profile", patientsHandler.GetMyProfile)
			patientRoutes.PUT("/profile", patientsHandler.UpdateMyProfile)
			patientRoutes.GET("/studies", patientsHandler.GetMyStudies)
			patientRoutes.GET("/offer-requests", offersHandler.GetMyOfferRequests)
			patientRoutes.POST("/offer-requests", offersHandler.CreateOfferRequest)
			patientRoutes.POST("/offers/:id/accept", offersHandler.AcceptOffer)
			patientRoutes.GET("/orders", ordersHandler.GetMyOrders)
		}

		// Study routes (patients and clinics)
		studyRoutes := protected.Group("/studies")
		{
			studyRoutes.POST("", rbac.RequireRole(rbac.RolePatient), studiesHandler.CreateStudy)
			studyRoutes.GET("/:id", studiesHandler.GetStudy)
			studyRoutes.PATCH("/:id/status", rbac.RequireRole(rbac.RoleAdmin), studiesHandler.UpdateStudyStatus)
			studyRoutes.GET("/:id/pdf", studiesHandler.GetStudyPDF)
		}

		// Treatment plan routes
		planRoutes := protected.Group("/plans")
		{
			planRoutes.POST("", plansHandler.CreatePlan)
			planRoutes.GET("/:id", plansHandler.GetPlan)
			planRoutes.GET("/:id/estimate", plansHandler.GetEstimate)
			planRoutes.GET("/study/:study_id", plansHandler.GetPlansByStudy)
		}

		// Clinic routes
		clinicRoutes := protected.Group("/clinic")
		clinicRoutes.Use(rbac.RequireRole(rbac.RoleClinicDoctor, rbac.RoleClinicManager))
		{
			clinicRoutes.POST("/profile", clinicsHandler.CreateClinic)
			clinicRoutes.GET("/profile", clinicsHandler.GetMyClinic)
			clinicRoutes.PUT("/profile", clinicsHandler.UpdateMyClinic)

			// Pricelist
			clinicRoutes.GET("/pricelist", clinicsHandler.GetMyPricelist)
			clinicRoutes.POST("/pricelist", clinicsHandler.AddPriceItem)
			clinicRoutes.DELETE("/pricelist/:id", clinicsHandler.DeletePriceItem)

			// Offers
			clinicRoutes.GET("/offers", offersHandler.GetMyOffers)
			clinicRoutes.POST("/offers", offersHandler.CreateOffer)

			// Orders
			clinicRoutes.GET("/orders", ordersHandler.GetMyOrders)
			clinicRoutes.PATCH("/orders/:id/status", ordersHandler.UpdateOrderStatus)
		}

		// Offer request routes (accessible by multiple roles)
		protected.GET("/offer-requests/:id", offersHandler.GetOfferRequest)

		// Order routes (accessible by patient and clinic)
		protected.GET("/orders/:id", ordersHandler.GetOrder)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "dental-marketplace-api",
		})
	})

	// Start server
	addr := ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
