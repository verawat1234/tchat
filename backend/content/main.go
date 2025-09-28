package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"tchat.dev/content/handlers"
	"tchat.dev/content/models"
	"tchat.dev/content/services"
	"tchat.dev/content/utils"
)

// App represents the main application
type App struct {
	config    *utils.Config
	db        *gorm.DB
	router    *gin.Engine
	server    *http.Server
	validator *validator.Validate

	// Repositories
	contentRepo  services.ContentRepository
	categoryRepo services.CategoryRepository
	versionRepo  services.VersionRepository

	// Services
	contentService *services.ContentService

	// Handlers
	contentHandlers *handlers.ContentHandlers
}

// NewApp creates a new application instance
func NewApp(cfg *utils.Config) *App {
	return &App{
		config:    cfg,
		validator: validator.New(),
	}
}

// Initialize initializes all application components
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

	log.Printf("Content service initialized successfully on port %d", a.config.Server.Port)
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

// runMigrations runs database migrations for content service models
func (a *App) runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.ContentItem{},
		&models.ContentCategory{},
		&models.ContentVersion{},
	)
}

// initRepositories initializes all data access repositories
func (a *App) initRepositories() error {
	a.contentRepo = services.NewPostgreSQLContentRepository(a.db)
	a.categoryRepo = services.NewPostgreSQLCategoryRepository(a.db)
	a.versionRepo = services.NewPostgreSQLVersionRepository(a.db)

	log.Println("Repositories initialized successfully")
	return nil
}

// initServices initializes all business logic services
func (a *App) initServices() error {
	a.contentService = services.NewContentService(
		a.contentRepo,
		a.categoryRepo,
		a.versionRepo,
		a.db,
	)

	log.Println("Services initialized successfully")
	return nil
}

// initHandlers initializes HTTP handlers
func (a *App) initHandlers() error {
	a.contentHandlers = handlers.NewContentHandlers(a.contentService)

	log.Println("Handlers initialized successfully")
	return nil
}

// initRouter initializes the HTTP router with all routes
func (a *App) initRouter() error {
	// Set Gin mode
	if !a.config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(securityHeaders())

	// Rate limiting (simplified for now)
	// TODO: Implement proper rate limiting

	// Health check endpoints
	router.GET("/health", a.healthCheck)
	router.GET("/ready", a.readinessCheck)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Content routes - complete registration
		content := v1.Group("/content")
		{
			// Content CRUD operations
			content.GET("", a.contentHandlers.GetContentItems)
			content.POST("", a.contentHandlers.CreateContent)
			content.GET("/:id", a.contentHandlers.GetContent)
			content.PUT("/:id", a.contentHandlers.UpdateContent)
			content.DELETE("/:id", a.contentHandlers.DeleteContent)

			// Content operations
			content.POST("/:id/publish", a.contentHandlers.PublishContent)
			content.POST("/:id/archive", a.contentHandlers.ArchiveContent)
			content.PUT("/bulk", a.contentHandlers.BulkUpdateContent)
			content.POST("/sync", a.contentHandlers.SyncContent)

			// Category operations
			content.GET("/categories", a.contentHandlers.GetContentCategories)
			content.GET("/category/:category", a.contentHandlers.GetContentByCategory)

			// Version operations
			content.GET("/:id/versions", a.contentHandlers.GetContentVersions)
			content.POST("/:id/versions/:version/revert", a.contentHandlers.RevertContentVersion)

			// Health check
			content.GET("/health", a.contentHealth)
		}
	}

	// Swagger documentation (if enabled)
	if a.config.Debug {
		// Add swagger routes here if needed
	}

	a.router = router
	log.Println("Router initialized successfully")
	return nil
}

// initServer initializes the HTTP server
func (a *App) initServer() {
	addr := fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port)

	a.server = &http.Server{
		Addr:         addr,
		Handler:      a.router,
		ReadTimeout:  a.config.Server.ReadTimeout,
		WriteTimeout: a.config.Server.WriteTimeout,
		IdleTimeout:  a.config.Server.IdleTimeout,
	}
}

// Run starts the application server
func (a *App) Run() error {
	log.Printf("Starting content service on %s", a.server.Addr)

	// Start server in a goroutine
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Content service is running on %s", a.server.Addr)
	return nil
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down content service...")

	// Shutdown HTTP server
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	// Close database connection
	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	log.Println("Content service shutdown completed")
	return nil
}

// Health check endpoint
func (a *App) healthCheck(c *gin.Context) {
	utils.SuccessResponse(c, gin.H{
		"status":    "ok",
		"service":   "content-service",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC(),
	})
}

// Readiness check endpoint
func (a *App) readinessCheck(c *gin.Context) {
	// Check database connection
	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err != nil || sqlDB.Ping() != nil {
			utils.ErrorResponse(c, http.StatusServiceUnavailable, "Database not ready", "Database connection failed")
			return
		}
	}

	utils.SuccessResponse(c, gin.H{
		"status":   "ready",
		"service":  "content-service",
		"database": "connected",
	})
}

// contentHealth provides a simple content service health check endpoint for API routes
func (a *App) contentHealth(c *gin.Context) {
	utils.SuccessResponse(c, gin.H{
		"status":    "ok",
		"service":   "content-service",
		"api":       "available",
		"timestamp": time.Now().UTC(),
	})
}

func main() {
	// Load configuration
	cfg := utils.LoadConfig()

	// Override port for content service if not set
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8086 // Content service port (matches Docker configuration)
	}

	// Create and initialize application
	app := NewApp(cfg)
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Setup graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start the application
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutdown signal received")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown gracefully: %v", err)
		os.Exit(1)
	}

	log.Println("Content service stopped")
}

// corsMiddleware provides CORS support
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// securityHeaders adds security headers
func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}