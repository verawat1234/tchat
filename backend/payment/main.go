package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"tchat.dev/payment/handlers"
	"tchat.dev/payment/internal/config"
	"tchat.dev/payment/internal/database"
	"tchat.dev/payment/repositories"
	"tchat.dev/payment/services"
	"tchat.dev/shared/middleware"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate tables
	if err := autoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	walletRepo := repositories.NewWalletRepository(db)
	txRepo := repositories.NewTransactionRepository(db)
	eventPublisher := &services.NoOpEventPublisher{}

	// Initialize services
	paymentService := services.NewPaymentService(walletRepo, txRepo, eventPublisher, db)

	// Initialize handlers
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "payment-service",
			"version": "1.0.0",
		})
	})

	// Register API routes
	api := router.Group("/api/v1")
	paymentHandler.RegisterRoutes(api)

	// Start server
	port := cfg.Server.Port
	if port == "" {
		port = "8084"  // Payment service runs on 8084
	}

	log.Printf("Payment service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func autoMigrate(db *gorm.DB) error {
	// Auto-migrate wallet and transaction tables if they don't exist
	// The shared models should handle their own migration
	return nil
}