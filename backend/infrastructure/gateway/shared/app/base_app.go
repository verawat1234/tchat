package app

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

	"tchat.dev/shared/config"
	"tchat.dev/shared/database"
	"tchat.dev/shared/responses"
)

// ServiceInitializer defines the interface that each service must implement
type ServiceInitializer interface {
	// GetModels returns all models for database migration
	GetModels() []interface{}

	// InitializeRepositories sets up data access layer
	InitializeRepositories(db *gorm.DB) error

	// InitializeServices sets up business logic layer
	InitializeServices(db *gorm.DB) error

	// InitializeHandlers sets up HTTP handlers
	InitializeHandlers() error

	// RegisterRoutes registers all service-specific routes
	RegisterRoutes(router *gin.Engine) error

	// GetServiceInfo returns service name and version for health checks
	GetServiceInfo() (string, string)
}

// BaseApp provides common application functionality for all microservices
type BaseApp struct {
	config    *config.Config
	db        *gorm.DB
	router    *gin.Engine
	server    *http.Server
	validator *validator.Validate
	service   ServiceInitializer

	// Service info
	serviceName string
	version     string
}

// NewBaseApp creates a new base application instance
func NewBaseApp(cfg *config.Config, serviceImpl ServiceInitializer) *BaseApp {
	serviceName, version := serviceImpl.GetServiceInfo()

	return &BaseApp{
		config:      cfg,
		validator:   validator.New(),
		service:     serviceImpl,
		serviceName: serviceName,
		version:     version,
	}
}

// Initialize sets up the entire application stack
func (a *BaseApp) Initialize() error {
	log.Printf("Initializing %s v%s...", a.serviceName, a.version)

	// Initialize database
	if err := a.initDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Run migrations for service models
	if err := a.runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize service layers in correct order
	if err := a.service.InitializeRepositories(a.db); err != nil {
		return fmt.Errorf("failed to initialize repositories: %w", err)
	}

	if err := a.service.InitializeServices(a.db); err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	if err := a.service.InitializeHandlers(); err != nil {
		return fmt.Errorf("failed to initialize handlers: %w", err)
	}

	// Initialize router and routes
	if err := a.initRouter(); err != nil {
		return fmt.Errorf("failed to initialize router: %w", err)
	}

	// Initialize HTTP server
	a.initServer()

	log.Printf("%s initialized successfully on port %d", a.serviceName, a.config.Server.Port)
	return nil
}

// initDatabase initializes database connection with connection pooling
func (a *BaseApp) initDatabase() error {
	dsn := a.config.GetDatabaseURL()

	var gormLogger logger.Interface
	if a.config.Debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         gormLogger,
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

	a.db = db
	log.Printf("Database connection established for %s", a.serviceName)
	return nil
}

// runMigrations runs auto-migrations for service models
func (a *BaseApp) runMigrations() error {
	models := a.service.GetModels()
	if len(models) == 0 {
		log.Printf("No models to migrate for %s", a.serviceName)
		return nil
	}

	log.Printf("Running migrations for %s (%d models)...", a.serviceName, len(models))

	if err := a.db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to run auto-migrations: %w", err)
	}

	log.Printf("Migrations completed successfully for %s", a.serviceName)
	return nil
}

// initRouter initializes Gin router with standard middleware and routes
func (a *BaseApp) initRouter() error {
	// Set Gin mode based on debug setting
	if !a.config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add standard middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(a.corsMiddleware())
	router.Use(a.securityHeaders())

	// Add health check endpoints
	router.GET("/health", a.healthCheck)
	router.GET("/ready", a.readinessCheck)

	// Register service-specific routes
	if err := a.service.RegisterRoutes(router); err != nil {
		return fmt.Errorf("failed to register service routes: %w", err)
	}

	a.router = router
	log.Printf("Router initialized for %s", a.serviceName)
	return nil
}

// initServer initializes HTTP server with timeouts
func (a *BaseApp) initServer() {
	addr := fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port)

	a.server = &http.Server{
		Addr:         addr,
		Handler:      a.router,
		ReadTimeout:  a.config.Server.ReadTimeout,
		WriteTimeout: a.config.Server.WriteTimeout,
		IdleTimeout:  a.config.Server.IdleTimeout,
	}
}

// Run starts the HTTP server
func (a *BaseApp) Run() error {
	log.Printf("Starting %s on %s", a.serviceName, a.server.Addr)

	// Start server in goroutine
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start %s server: %v", a.serviceName, err)
		}
	}()

	log.Printf("%s is running on %s", a.serviceName, a.server.Addr)
	return nil
}

// Shutdown gracefully shuts down the application
func (a *BaseApp) Shutdown(ctx context.Context) error {
	log.Printf("Shutting down %s...", a.serviceName)

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

	log.Printf("%s shutdown completed", a.serviceName)
	return nil
}

// RunWithGracefulShutdown runs the application with graceful shutdown handling
func (a *BaseApp) RunWithGracefulShutdown() {
	// Setup graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start the application
	if err := a.Run(); err != nil {
		log.Fatalf("Failed to start %s: %v", a.serviceName, err)
	}

	// Wait for shutdown signal
	<-ctx.Done()
	log.Printf("Shutdown signal received for %s", a.serviceName)

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown %s gracefully: %v", a.serviceName, err)
		os.Exit(1)
	}

	log.Printf("%s stopped", a.serviceName)
}

// Health check endpoints

// healthCheck provides basic service health status
func (a *BaseApp) healthCheck(c *gin.Context) {
	responses.SendSuccessResponse(c, gin.H{
		"status":    "ok",
		"service":   a.serviceName,
		"version":   a.version,
		"timestamp": time.Now().UTC(),
	})
}

// readinessCheck provides detailed readiness status including database
func (a *BaseApp) readinessCheck(c *gin.Context) {
	// Check database connection
	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err != nil || sqlDB.Ping() != nil {
			responses.SendErrorResponse(c, http.StatusServiceUnavailable, "database_not_ready", "Database connection failed")
			return
		}
	}

	responses.SendSuccessResponse(c, gin.H{
		"status":   "ready",
		"service":  a.serviceName,
		"version":  a.version,
		"database": "connected",
	})
}

// Middleware

// corsMiddleware provides CORS support
func (a *BaseApp) corsMiddleware() gin.HandlerFunc {
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
func (a *BaseApp) securityHeaders() gin.HandlerFunc {
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

// GetDB returns the database instance (for testing or advanced usage)
func (a *BaseApp) GetDB() *gorm.DB {
	return a.db
}

// GetRouter returns the router instance (for testing or advanced usage)
func (a *BaseApp) GetRouter() *gin.Engine {
	return a.router
}

// GetConfig returns the configuration instance
func (a *BaseApp) GetConfig() *config.Config {
	return a.config
}