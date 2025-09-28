package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"tchat.dev/shared/config"
	"tchat.dev/shared/database"
	"tchat.dev/shared/middleware"
	"tchat.dev/shared/responses"
	sharedModels "tchat.dev/shared/models"
	"tchat.dev/tests/fixtures"
)

// App represents the main commerce application
type App struct {
	config    *config.Config
	db        *gorm.DB
	router    *gin.Engine
	server    *http.Server
	validator *validator.Validate
}

// NewApp creates a new commerce application instance
func NewApp(cfg *config.Config) *App {
	return &App{
		config:    cfg,
		validator: validator.New(),
	}
}

// Initialize sets up the application
func (a *App) Initialize() error {
	// Initialize database
	if err := a.initDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize router
	if err := a.initRouter(); err != nil {
		return fmt.Errorf("failed to initialize router: %w", err)
	}

	// Initialize server
	a.initServer()

	log.Printf("Commerce service initialized successfully on port %d", a.config.Server.Port)
	return nil
}

// initDatabase initializes the database connection and runs migrations
func (a *App) initDatabase() error {
	// Create database connection
	dsn := a.config.GetDatabaseURL()

	var gormLogger logger.Interface
	if a.config.Debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NamingStrategy: database.NamingStrategy{},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(a.config.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(a.config.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(a.config.Database.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(a.config.Database.ConnMaxIdleTime)

	// Run auto-migrations
	if err := a.runMigrations(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	a.db = db
	log.Println("Database connection established and migrations completed")
	return nil
}

// runMigrations runs database migrations for commerce service models
func (a *App) runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&sharedModels.Business{},
		&sharedModels.Product{},
		&sharedModels.Order{},
		&sharedModels.Event{},
	)
}

// initRouter initializes the HTTP router
func (a *App) initRouter() error {
	gin.SetMode(gin.ReleaseMode)
	if a.config.Debug {
		gin.SetMode(gin.DebugMode)
	}

	a.router = gin.New()

	// Add middleware
	a.router.Use(gin.Logger())
	a.router.Use(gin.Recovery())
	a.router.Use(middleware.CORS())
	a.router.Use(middleware.SecurityHeaders())

	// Health check endpoints
	a.router.GET("/health", a.healthCheck)
	a.router.GET("/ready", a.readinessCheck)

	// Mock commerce endpoints
	v1 := a.router.Group("/api/v1")
	{
		v1.GET("/shops", a.getShops)
		v1.POST("/shops", a.createShop)
		v1.GET("/products", a.getProducts)
		v1.POST("/products", a.createProduct)
	}

	// Seed data endpoints (for development/testing only)
	if a.config.Debug {
		seed := a.router.Group("/seed")
		{
			seed.POST("/commerce", a.seedCommerceData)
			seed.POST("/commerce/businesses", a.seedBusinesses)
			seed.POST("/commerce/products", a.seedProducts)
			seed.DELETE("/commerce", a.clearCommerceData)
		}
	}

	return nil
}

// initServer initializes the HTTP server
func (a *App) initServer() {
	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.Server.Port),
		Handler: a.router,
	}
}

// healthCheck handles health check requests
func (a *App) healthCheck(c *gin.Context) {
	responses.SendSuccessResponse(c, map[string]string{
		"status":    "healthy",
		"service":   "commerce-service",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// readinessCheck handles readiness probe requests
func (a *App) readinessCheck(c *gin.Context) {
	// Check database connection
	sqlDB, err := a.db.DB()
	if err != nil {
		responses.SendErrorResponse(c, http.StatusServiceUnavailable, "database connection not available", "DATABASE_ERROR")
		return
	}

	if err := sqlDB.Ping(); err != nil {
		responses.SendErrorResponse(c, http.StatusServiceUnavailable, "database ping failed", "DATABASE_PING_ERROR")
		return
	}

	responses.SendSuccessResponse(c, map[string]string{
		"status":    "ready",
		"service":   "commerce-service",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// Mock commerce endpoints
func (a *App) getShops(c *gin.Context) {
	responses.SendSuccessResponse(c, []map[string]interface{}{
		{
			"id":     uuid.New().String(),
			"name":   "Sample Shop",
			"status": "active",
		},
	})
}

func (a *App) createShop(c *gin.Context) {
	responses.SendSuccessResponse(c, map[string]interface{}{
		"id":     uuid.New().String(),
		"name":   "New Shop",
		"status": "created",
	})
}

// getProducts handles getting products
func (a *App) getProducts(c *gin.Context) {
	var products []sharedModels.Product

	result := a.db.Find(&products)
	if result.Error != nil {
		responses.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch products", "DATABASE_ERROR")
		return
	}

	responses.SendSuccessResponse(c, products)
}

// createProduct handles creating a new product
func (a *App) createProduct(c *gin.Context) {
	var product sharedModels.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "Invalid product data", "VALIDATION_ERROR")
		return
	}

	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}

	result := a.db.Create(&product)
	if result.Error != nil {
		responses.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create product", "DATABASE_ERROR")
		return
	}

	responses.SendSuccessResponse(c, product)
}

// seedCommerceData handles seeding all commerce data
func (a *App) seedCommerceData(c *gin.Context) {
	country := c.DefaultQuery("country", "TH")

	// Create fixtures
	masterFixtures := fixtures.NewMasterFixtures()
	commerceFixtures := masterFixtures.Commerce

	// Generate businesses and products for the specified country
	businesses := make([]*sharedModels.Business, 0)
	products := make([]*sharedModels.Product, 0)

	// Create multiple businesses for variety
	businessTypes := []string{"electronics", "fashion", "food"}
	for _, businessType := range businessTypes {
		var business *sharedModels.Business
		switch businessType {
		case "electronics":
			business = commerceFixtures.ElectronicsBusiness(country)
		case "fashion":
			business = commerceFixtures.FashionBusiness(country)
		default:
			business = commerceFixtures.BasicBusiness(country)
		}

		businesses = append(businesses, business)

		// Create products for each business
		electronics := commerceFixtures.ElectronicsProduct(business.ID, country)
		fashion := commerceFixtures.FashionProduct(business.ID, country)
		food := commerceFixtures.FoodProduct(business.ID, country)

		products = append(products, electronics, fashion, food)
	}

	// Save to database
	for _, business := range businesses {
		if err := a.db.Create(business).Error; err != nil {
			log.Printf("Failed to create business %s: %v", business.Name, err)
		}
	}

	for _, product := range products {
		if err := a.db.Create(product).Error; err != nil {
			log.Printf("Failed to create product %s: %v", product.Name, err)
		}
	}

	responses.SendSuccessResponse(c, map[string]interface{}{
		"message":        "Commerce data seeded successfully",
		"country":        country,
		"businesses_count": len(businesses),
		"products_count":   len(products),
		"businesses":     businesses,
		"products":       products,
	})
}

// seedBusinesses handles seeding business data only
func (a *App) seedBusinesses(c *gin.Context) {
	country := c.DefaultQuery("country", "TH")
	count := 5 // Default count

	if countStr := c.Query("count"); countStr != "" {
		if parsedCount, err := strconv.Atoi(countStr); err == nil && parsedCount > 0 {
			count = parsedCount
		}
	}

	commerceFixtures := fixtures.NewCommerceFixtures()
	businesses := make([]*sharedModels.Business, 0, count)

	for i := 0; i < count; i++ {
		var business *sharedModels.Business
		switch i % 3 {
		case 0:
			business = commerceFixtures.ElectronicsBusiness(country)
		case 1:
			business = commerceFixtures.FashionBusiness(country)
		default:
			business = commerceFixtures.BasicBusiness(country)
		}

		if err := a.db.Create(business).Error; err != nil {
			log.Printf("Failed to create business %s: %v", business.Name, err)
			continue
		}

		businesses = append(businesses, business)
	}

	responses.SendSuccessResponse(c, map[string]interface{}{
		"message":    "Businesses seeded successfully",
		"country":    country,
		"count":      len(businesses),
		"businesses": businesses,
	})
}

// seedProducts handles seeding product data only
func (a *App) seedProducts(c *gin.Context) {
	country := c.DefaultQuery("country", "TH")
	businessIDStr := c.Query("business_id")

	var businessID uuid.UUID
	if businessIDStr != "" {
		var err error
		businessID, err = uuid.Parse(businessIDStr)
		if err != nil {
			responses.SendErrorResponse(c, http.StatusBadRequest, "Invalid business_id", "VALIDATION_ERROR")
			return
		}
	} else {
		// Get a random business or create one
		var business sharedModels.Business
		result := a.db.First(&business)
		if result.Error != nil {
			// Create a new business
			commerceFixtures := fixtures.NewCommerceFixtures()
			newBusiness := commerceFixtures.BasicBusiness(country)
			if err := a.db.Create(newBusiness).Error; err != nil {
				responses.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create business for products", "DATABASE_ERROR")
				return
			}
			businessID = newBusiness.ID
		} else {
			businessID = business.ID
		}
	}

	commerceFixtures := fixtures.NewCommerceFixtures()
	products := make([]*sharedModels.Product, 0)

	// Create different types of products
	electronics := commerceFixtures.ElectronicsProduct(businessID, country)
	fashion := commerceFixtures.FashionProduct(businessID, country)
	food := commerceFixtures.FoodProduct(businessID, country)
	basic := commerceFixtures.BasicProduct(businessID, country)

	allProducts := []*sharedModels.Product{electronics, fashion, food, basic}

	for _, product := range allProducts {
		if err := a.db.Create(product).Error; err != nil {
			log.Printf("Failed to create product %s: %v", product.Name, err)
			continue
		}
		products = append(products, product)
	}

	responses.SendSuccessResponse(c, map[string]interface{}{
		"message":     "Products seeded successfully",
		"country":     country,
		"business_id": businessID.String(),
		"count":       len(products),
		"products":    products,
	})
}

// clearCommerceData handles clearing all commerce data
func (a *App) clearCommerceData(c *gin.Context) {
	confirm := c.Query("confirm")
	if confirm != "true" {
		responses.SendErrorResponse(c, http.StatusBadRequest, "Add ?confirm=true to confirm deletion", "CONFIRMATION_REQUIRED")
		return
	}

	// Delete in order (products first due to foreign key constraints)
	productResult := a.db.Unscoped().Delete(&sharedModels.Product{}, "1 = 1")
	businessResult := a.db.Unscoped().Delete(&sharedModels.Business{}, "1 = 1")

	responses.SendSuccessResponse(c, map[string]interface{}{
		"message":           "Commerce data cleared successfully",
		"products_deleted":  productResult.RowsAffected,
		"businesses_deleted": businessResult.RowsAffected,
	})
}

// Run starts the HTTP server
func (a *App) Run() error {
	// Graceful shutdown
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Commerce service is running on port %d", a.config.Server.Port)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	// Close database connection
	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	log.Println("Commerce service stopped")
	return nil
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Override with environment variables if present
	if dbHost := os.Getenv("DATABASE_HOST"); dbHost != "" {
		cfg.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DATABASE_PORT"); dbPort != "" {
		if port, err := strconv.Atoi(dbPort); err == nil {
			cfg.Database.Port = port
		}
	}
	if dbPassword := os.Getenv("DATABASE_PASSWORD"); dbPassword != "" {
		cfg.Database.Password = dbPassword
	}
	if dbName := os.Getenv("COMMERCE_DATABASE_NAME"); dbName != "" {
		cfg.Database.Database = dbName
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		cfg.JWT.Secret = jwtSecret
	}
	// Check for service-specific port first, then general SERVER_PORT
	if serverPort := os.Getenv("COMMERCE_SERVICE_PORT"); serverPort != "" {
		if port, err := strconv.Atoi(serverPort); err == nil {
			cfg.Server.Port = port
		}
	} else if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		if port, err := strconv.Atoi(serverPort); err == nil {
			cfg.Server.Port = port
		}
	}

	// Default to 8083 for commerce service if not set
	if cfg.Server.Port == 0 || cfg.Server.Port == 8080 {
		cfg.Server.Port = 8083
	}

	// Create and initialize app
	app := NewApp(cfg)
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start the server
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}