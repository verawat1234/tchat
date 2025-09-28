package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tchat/backend/audit/handlers"
	"github.com/tchat/backend/audit/models"
	"github.com/tchat/backend/audit/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Initialize database
	db, err := initDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Auto-migrate database tables
	if err := autoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize services
	placeholderService := services.NewPlaceholderService(db)
	serviceCompletionService := services.NewServiceCompletionService(db)
	validationService := services.NewValidationService(db)

	// Initialize handlers
	placeholdersHandler := handlers.NewPlaceholdersHandler(placeholderService)
	serviceCompletionHandler := handlers.NewServiceCompletionHandler(serviceCompletionService)
	validationHandler := handlers.NewValidationHandler(validationService)

	// Setup routes
	router := setupRoutes(placeholdersHandler, serviceCompletionHandler, validationHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8093" // Default audit service port
	}

	log.Printf("Starting Audit Service on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initDatabase() (*gorm.DB, error) {
	// Get database URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://user:password@localhost:5432/tchat_audit?sslmode=disable"
	}

	// Configure GORM logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Info,
			Colorful: true,
		},
	)

	// Connect to database
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	// Auto-migrate all audit models
	return db.AutoMigrate(
		&models.PlaceholderItem{},
		&models.ServiceCompletion{},
		&models.CompletionAudit{},
	)
}

func setupRoutes(
	placeholdersHandler *handlers.PlaceholdersHandler,
	serviceCompletionHandler *handlers.ServiceCompletionHandler,
	validationHandler *handlers.ValidationHandler,
) *gin.Engine {
	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "audit",
			"version": "1.0.0",
		})
	})

	// API versioning
	v1 := router.Group("/api/v1")

	// Placeholder management routes
	placeholdersGroup := v1.Group("/audit/placeholders")
	{
		placeholdersGroup.GET("", placeholdersHandler.GetPlaceholders)
		placeholdersGroup.POST("", placeholdersHandler.CreatePlaceholder)
		placeholdersGroup.PATCH("/:id", placeholdersHandler.UpdatePlaceholder)
		placeholdersGroup.GET("/stats", placeholdersHandler.GetPlaceholderStats)
		placeholdersGroup.POST("/bulk-update", placeholdersHandler.BulkUpdatePlaceholders)
		placeholdersGroup.GET("/region/:region", placeholdersHandler.GetPlaceholdersByRegion)
		placeholdersGroup.POST("/:id/archive", placeholdersHandler.ArchivePlaceholder)
	}

	// Service completion routes
	servicesGroup := v1.Group("/audit/services")
	{
		servicesGroup.GET("/:serviceId/completion", serviceCompletionHandler.GetServiceCompletion)
		servicesGroup.GET("/completion", serviceCompletionHandler.GetAllServiceCompletions)
		servicesGroup.PUT("/:serviceId/completion", serviceCompletionHandler.UpdateServiceCompletion)
		servicesGroup.GET("/health", serviceCompletionHandler.GetServiceHealth)
		servicesGroup.GET("/metrics", serviceCompletionHandler.GetServiceMetrics)
		servicesGroup.POST("/:serviceId/refresh", serviceCompletionHandler.TriggerServiceRefresh)
		servicesGroup.GET("/regional", serviceCompletionHandler.GetRegionalOptimization)
		servicesGroup.GET("/dependencies", serviceCompletionHandler.GetDependencyGraph)
		servicesGroup.GET("/trends", serviceCompletionHandler.GetCompletionTrends)
	}

	// Validation routes
	validationGroup := v1.Group("/audit/validation")
	{
		validationGroup.POST("", validationHandler.RunValidation)
		validationGroup.GET("/:auditId/status", validationHandler.GetValidationStatus)
		validationGroup.GET("/:auditId/results", validationHandler.GetValidationResults)
		validationGroup.POST("/:auditId/cancel", validationHandler.CancelValidation)
		validationGroup.GET("/:auditId/report", validationHandler.GetValidationReport)
		validationGroup.GET("/metrics", validationHandler.GetValidationMetrics)
		validationGroup.POST("/quick", validationHandler.RunQuickValidation)
	}

	// Additional validation routes
	v1.GET("/audit/validations", validationHandler.ListValidations)

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}