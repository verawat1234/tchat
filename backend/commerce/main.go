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
	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		if port, err := strconv.Atoi(serverPort); err == nil {
			cfg.Server.Port = port
		}
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