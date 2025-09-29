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

	"tchat.dev/commerce/handlers"
	commerceMiddleware "tchat.dev/commerce/middleware"
	"tchat.dev/commerce/models"
	"tchat.dev/commerce/repository"
	"tchat.dev/commerce/services"
)

// App represents the main commerce application
type App struct {
	config    *config.Config
	db        *gorm.DB
	router    *gin.Engine
	server    *http.Server
	validator *validator.Validate

	// Performance middleware
	performanceMiddleware *commerceMiddleware.PerformanceMiddleware

	// Repositories
	businessRepo         repository.BusinessRepository
	productRepo          repository.ProductRepository
	reviewRepo           repository.ReviewRepository
	wishlistRepo         repository.WishlistRepository
	productFollowRepo    repository.ProductFollowRepository
	wishlistShareRepo    repository.WishlistShareRepository
	cartRepo             repository.CartRepository
	cartAbandonmentRepo  repository.CartAbandonmentRepository
	categoryRepo         repository.CategoryRepository
	productCategoryRepo  repository.ProductCategoryRepository
	categoryViewRepo     repository.CategoryViewRepository
	streamRepo           repository.StreamRepository

	// Services
	businessService services.BusinessService
	productService  services.ProductService
	reviewService   services.ReviewService
	wishlistService services.WishlistService
	cartService     services.CartService
	categoryService services.CategoryService
	// Stream services
	streamCategoryService *services.StreamCategoryService
	streamContentService  *services.StreamContentService
	streamSessionService  *services.StreamSessionService
	streamPurchaseService *services.StreamPurchaseService

	// Handlers
	businessHandler *handlers.BusinessHandler
	productHandler  *handlers.ProductHandler
	reviewHandler   *handlers.ReviewHandler
	wishlistHandler *handlers.WishlistHandler
	cartHandler     *handlers.CartHandler
	categoryHandler *handlers.CategoryHandler
	streamHandler   *handlers.StreamHandler
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

	// Initialize repositories
	if err := a.initRepositories(); err != nil {
		return fmt.Errorf("failed to initialize repositories: %w", err)
	}

	// Initialize services
	if err := a.initServices(); err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	// Initialize handlers
	if err := a.initHandlers(); err != nil {
		return fmt.Errorf("failed to initialize handlers: %w", err)
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
		// Shared models
		&sharedModels.Business{},
		&sharedModels.Product{},
		&sharedModels.Order{},
		&sharedModels.Event{},
		// Commerce-specific models
		&models.Review{},
		&models.Wishlist{},
		&models.ProductFollow{},
		&models.WishlistShare{},
		&models.Cart{},
		&models.CartAbandonmentTracking{},
		&models.Category{},
		&models.ProductCategory{},
		&models.CategoryView{},
		// Stream models
		&models.StreamCategory{},
		&models.StreamSubtab{},
		&models.StreamContentItem{},
		&models.TabNavigationState{},
		&models.StreamUserSession{},
		&models.StreamContentView{},
		&models.StreamUserPreference{},
	)
}

// initRepositories initializes all data access repositories
func (a *App) initRepositories() error {
	a.businessRepo = repository.NewBusinessRepository(a.db)
	a.productRepo = repository.NewProductRepository(a.db)
	a.reviewRepo = repository.NewReviewRepository(a.db)
	a.wishlistRepo = repository.NewWishlistRepository(a.db)
	a.productFollowRepo = repository.NewProductFollowRepository(a.db)
	a.wishlistShareRepo = repository.NewWishlistShareRepository(a.db)
	a.cartRepo = repository.NewCartRepository(a.db)
	a.cartAbandonmentRepo = repository.NewCartAbandonmentRepository(a.db)
	a.categoryRepo = repository.NewCategoryRepository(a.db)
	a.productCategoryRepo = repository.NewProductCategoryRepository(a.db)
	a.categoryViewRepo = repository.NewCategoryViewRepository(a.db)
	a.streamRepo = repository.NewStreamRepository(a.db)

	log.Println("Repositories initialized successfully")
	return nil
}

// initServices initializes all business logic services
func (a *App) initServices() error {
	a.businessService = services.NewBusinessService(a.businessRepo, a.productRepo)
	a.productService = services.NewProductService(a.productRepo, a.businessRepo)
	a.reviewService = services.NewReviewService(a.reviewRepo, a.db)
	a.wishlistService = services.NewWishlistService(a.wishlistRepo, a.productFollowRepo, a.wishlistShareRepo, a.db)
	a.cartService = services.NewCartService(a.cartRepo, a.cartAbandonmentRepo, a.db)
	a.categoryService = services.NewCategoryService(a.categoryRepo, a.productCategoryRepo, a.categoryViewRepo, a.db)
	// Initialize stream services
	a.streamCategoryService = services.NewStreamCategoryService(a.streamRepo)
	a.streamContentService = services.NewStreamContentService(a.streamRepo)
	a.streamSessionService = services.NewStreamSessionService(a.streamRepo)
	a.streamPurchaseService = services.NewStreamPurchaseService(a.streamRepo)

	log.Println("Services initialized successfully")
	return nil
}

// initHandlers initializes HTTP handlers
func (a *App) initHandlers() error {
	a.businessHandler = handlers.NewBusinessHandler(a.businessService)
	a.productHandler = handlers.NewProductHandler(a.productService)
	a.reviewHandler = handlers.NewReviewHandler(a.reviewService)
	a.wishlistHandler = handlers.NewWishlistHandler(a.wishlistService)
	a.cartHandler = handlers.NewCartHandler(a.cartService)
	a.categoryHandler = handlers.NewCategoryHandler(a.categoryService)
	a.streamHandler = handlers.NewStreamHandler(
		a.streamCategoryService,
		a.streamContentService,
		a.streamSessionService,
		a.streamPurchaseService,
	)

	log.Println("Handlers initialized successfully")
	return nil
}

// initRouter initializes the HTTP router
func (a *App) initRouter() error {
	gin.SetMode(gin.ReleaseMode)
	if a.config.Debug {
		gin.SetMode(gin.DebugMode)
	}

	a.router = gin.New()

	// Initialize performance middleware
	performanceConfig := commerceMiddleware.DefaultPerformanceConfig()
	a.performanceMiddleware = commerceMiddleware.NewPerformanceMiddleware(performanceConfig)

	// Add middleware
	a.router.Use(gin.Logger())
	a.router.Use(gin.Recovery())
	a.router.Use(middleware.CORS())
	a.router.Use(middleware.SecurityHeaders())
	a.router.Use(a.performanceMiddleware.Handler()) // Add performance middleware

	// Health check endpoints
	a.router.GET("/health", a.healthCheck)
	a.router.GET("/ready", a.readinessCheck)

	// Performance monitoring endpoints
	perf := a.router.Group("/performance")
	{
		perf.GET("/metrics", a.getPerformanceMetrics)
		perf.GET("/cache/stats", a.getCacheStats)
		perf.POST("/cache/clear", a.clearCache)
		perf.POST("/metrics/reset", a.resetMetrics)
	}

	// Commerce API endpoints
	v1 := a.router.Group("/api/v1")
	{
		// Business (Shop) routes
		shops := v1.Group("/shops")
		{
			shops.GET("", a.businessHandler.GetBusinesses)
			shops.POST("", a.businessHandler.CreateBusiness)
			shops.GET("/:id", a.businessHandler.GetBusiness)
			shops.PUT("/:id", a.businessHandler.UpdateBusiness)
			shops.DELETE("/:id", a.businessHandler.DeleteBusiness)
			shops.GET("/:id/products", a.businessHandler.GetBusinessProducts)
		}

		// Product routes
		products := v1.Group("/products")
		{
			products.GET("", a.productHandler.GetProducts)
			products.POST("", a.productHandler.CreateProduct)
			products.GET("/:id", a.productHandler.GetProduct)
			products.PUT("/:id", a.productHandler.UpdateProduct)
			products.DELETE("/:id", a.productHandler.DeleteProduct)
			products.GET("/:productId/reviews", a.reviewHandler.GetProductReviews)
			products.POST("/:productId/follow", a.wishlistHandler.FollowProduct)
			products.DELETE("/:productId/follow", a.wishlistHandler.UnfollowProduct)
		}

		// Review routes
		reviews := v1.Group("/reviews")
		{
			reviews.GET("", a.reviewHandler.ListReviews)
			reviews.POST("", a.reviewHandler.CreateReview)
			reviews.GET("/:id", a.reviewHandler.GetReview)
			reviews.PUT("/:id", a.reviewHandler.UpdateReview)
			reviews.DELETE("/:id", a.reviewHandler.DeleteReview)
			reviews.POST("/:reviewId/helpful", a.reviewHandler.MarkReviewHelpful)
			reviews.POST("/:reviewId/report", a.reviewHandler.ReportReview)
			reviews.POST("/:reviewId/moderate", a.reviewHandler.ModerateReview)
			reviews.GET("/average-rating", a.reviewHandler.GetAverageRating)
		}

		// Business review routes
		businesses := v1.Group("/businesses")
		{
			businesses.GET("/:businessId/reviews", a.reviewHandler.GetBusinessReviews)
			businesses.GET("/:businessId/categories", a.categoryHandler.GetBusinessCategories)
		}

		// Wishlist routes
		wishlists := v1.Group("/wishlists")
		{
			wishlists.GET("", a.wishlistHandler.ListUserWishlists)
			wishlists.POST("", a.wishlistHandler.CreateWishlist)
			wishlists.GET("/default", a.wishlistHandler.GetDefaultWishlist)
			wishlists.GET("/shared", a.wishlistHandler.GetSharedWishlists)
			wishlists.GET("/shared/:token", a.wishlistHandler.GetWishlistByShareToken)
			wishlists.GET("/:id", a.wishlistHandler.GetWishlist)
			wishlists.PUT("/:id", a.wishlistHandler.UpdateWishlist)
			wishlists.DELETE("/:id", a.wishlistHandler.DeleteWishlist)
			wishlists.POST("/:id/items", a.wishlistHandler.AddToWishlist)
			wishlists.DELETE("/:id/items/:productId", a.wishlistHandler.RemoveFromWishlist)
			wishlists.POST("/:id/share", a.wishlistHandler.ShareWishlist)
		}

		// Product following routes
		following := v1.Group("/products/following")
		{
			following.GET("", a.wishlistHandler.ListFollowedProducts)
		}

		// Cart routes
		carts := v1.Group("/carts")
		{
			carts.GET("", a.cartHandler.GetCart)
			carts.POST("/items", a.cartHandler.AddToCart)
			carts.PUT("/items/:itemId", a.cartHandler.UpdateCartItem)
			carts.DELETE("/items/:itemId", a.cartHandler.RemoveFromCart)
			carts.POST("/clear", a.cartHandler.ClearCart)
			carts.POST("/merge", a.cartHandler.MergeCart)
			carts.GET("/abandoned", a.cartHandler.GetAbandonedCarts)
			carts.POST("/abandonment", a.cartHandler.CreateAbandonmentTracking)
			carts.GET("/abandonment/analytics", a.cartHandler.GetAbandonmentAnalytics)
			carts.POST("/:cartId/coupons", a.cartHandler.ApplyCoupon)
			carts.DELETE("/:cartId/coupons", a.cartHandler.RemoveCoupon)
			carts.GET("/:cartId/validate", a.cartHandler.ValidateCart)
		}

		// Category routes
		categories := v1.Group("/categories")
		{
			categories.GET("", a.categoryHandler.ListCategories)
			categories.POST("", a.categoryHandler.CreateCategory)
			categories.GET("/global", a.categoryHandler.GetGlobalCategories)
			categories.GET("/featured", a.categoryHandler.GetFeaturedCategories)
			categories.GET("/root", a.categoryHandler.GetRootCategories)
			categories.GET("/path/:path", a.categoryHandler.GetCategoryByPath)
			categories.GET("/:id", a.categoryHandler.GetCategory)
			categories.PUT("/:id", a.categoryHandler.UpdateCategory)
			categories.DELETE("/:id", a.categoryHandler.DeleteCategory)
			categories.GET("/:id/children", a.categoryHandler.GetCategoryChildren)
			categories.GET("/:categoryId/products", a.categoryHandler.GetCategoryProducts)
			categories.POST("/:categoryId/products", a.categoryHandler.AddProductToCategory)
			categories.DELETE("/:categoryId/products/:productId", a.categoryHandler.RemoveProductFromCategory)
			categories.POST("/:categoryId/views", a.categoryHandler.TrackCategoryView)
			categories.GET("/:categoryId/analytics", a.categoryHandler.GetCategoryAnalytics)
		}

		// Stream routes
		stream := v1.Group("/stream")
		{
			// Category routes
			stream.GET("/categories", a.streamHandler.GetStreamCategories)
			stream.GET("/categories/:id", a.streamHandler.GetStreamCategoryDetail)

			// Content routes
			stream.GET("/content", a.streamHandler.GetStreamContent)
			stream.GET("/content/:id", a.streamHandler.GetStreamContentDetail)
			stream.GET("/featured", a.streamHandler.GetStreamFeatured)
			stream.GET("/search", a.streamHandler.SearchStreamContent)

			// Purchase routes
			stream.POST("/content/purchase", a.streamHandler.PostStreamContentPurchase)

			// Navigation routes
			stream.GET("/navigation", a.streamHandler.GetUserNavigationState)
			stream.PUT("/navigation", a.streamHandler.UpdateUserNavigationState)

			// Progress tracking routes
			stream.PUT("/content/:id/progress", a.streamHandler.UpdateContentViewProgress)

			// User preferences routes
			stream.GET("/preferences", a.streamHandler.GetUserPreferences)
			stream.PUT("/preferences", a.streamHandler.UpdateUserPreferences)
		}
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

// Commerce API endpoints are now handled by proper handlers in the handlers package

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

// Performance monitoring handlers
func (a *App) getPerformanceMetrics(c *gin.Context) {
	if a.performanceMiddleware == nil {
		responses.SendErrorResponse(c, http.StatusServiceUnavailable, "Performance monitoring not available", "PERF_NOT_AVAILABLE")
		return
	}

	metrics := a.performanceMiddleware.GetMetrics()
	responses.SendSuccessResponse(c, map[string]interface{}{
		"metrics": metrics,
		"service": "commerce-service",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (a *App) getCacheStats(c *gin.Context) {
	if a.performanceMiddleware == nil {
		responses.SendErrorResponse(c, http.StatusServiceUnavailable, "Performance monitoring not available", "PERF_NOT_AVAILABLE")
		return
	}

	stats := a.performanceMiddleware.GetCacheStats()
	responses.SendSuccessResponse(c, map[string]interface{}{
		"cache_stats": stats,
		"service": "commerce-service",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (a *App) clearCache(c *gin.Context) {
	if a.performanceMiddleware == nil {
		responses.SendErrorResponse(c, http.StatusServiceUnavailable, "Performance monitoring not available", "PERF_NOT_AVAILABLE")
		return
	}

	a.performanceMiddleware.ClearCache()
	responses.SendSuccessResponse(c, map[string]interface{}{
		"message": "Cache cleared successfully",
		"service": "commerce-service",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (a *App) resetMetrics(c *gin.Context) {
	if a.performanceMiddleware == nil {
		responses.SendErrorResponse(c, http.StatusServiceUnavailable, "Performance monitoring not available", "PERF_NOT_AVAILABLE")
		return
	}

	a.performanceMiddleware.ResetMetrics()
	responses.SendSuccessResponse(c, map[string]interface{}{
		"message": "Performance metrics reset successfully",
		"service": "commerce-service",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
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