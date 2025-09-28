package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"tchat.dev/shared/config"
	"tchat.dev/shared/responses"
)

// ServicePattern defines a reusable microservice pattern
type ServicePattern struct {
	config       *config.Config
	db           *gorm.DB
	router       *gin.Engine
	server       *http.Server
	validator    *validator.Validate
	modelManager *ModelManager
}

// ServiceInitializer interface for service-specific initialization
type ServiceInitializer interface {
	GetModels() []interface{}
	InitializeRepositories(db *gorm.DB) error
	InitializeServices(db *gorm.DB) error
	InitializeHandlers() error
	RegisterRoutes(router *gin.Engine) error
	GetServiceInfo() (name, version string)
}

// NewServicePattern creates a new service pattern instance
func NewServicePattern(cfg *config.Config) *ServicePattern {
	return &ServicePattern{
		config:    cfg,
		validator: validator.New(),
	}
}

// Initialize initializes the service with the given initializer
func (s *ServicePattern) Initialize(initializer ServiceInitializer) error {
	// Initialize database
	if err := s.initDatabase(initializer); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize service-specific components
	if err := initializer.InitializeRepositories(s.db); err != nil {
		return fmt.Errorf("failed to initialize repositories: %w", err)
	}

	if err := initializer.InitializeServices(s.db); err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	if err := initializer.InitializeHandlers(); err != nil {
		return fmt.Errorf("failed to initialize handlers: %w", err)
	}

	// Initialize router
	if err := s.initRouter(initializer); err != nil {
		return fmt.Errorf("failed to initialize router: %w", err)
	}

	// Initialize server
	s.initServer()

	name, version := initializer.GetServiceInfo()
	log.Printf("%s (v%s) initialized successfully on port %d", name, version, s.config.Server.Port)
	return nil
}

// initDatabase initializes database with models from initializer
func (s *ServicePattern) initDatabase(initializer ServiceInitializer) error {
	// Create database connection using existing database URL
	dsn := s.config.GetDatabaseURL()

	gormConfig := &gorm.Config{}
	if s.config.Debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(s.config.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(s.config.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(s.config.Database.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(s.config.Database.ConnMaxIdleTime)

	// Run migrations with models from initializer
	models := initializer.GetModels()
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	s.db = db
	s.modelManager = NewModelManager(db)
	log.Println("Database connection established and migrations completed")
	return nil
}

// initRouter initializes router with routes from initializer
func (s *ServicePattern) initRouter(initializer ServiceInitializer) error {
	if !s.config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add common middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())
	router.Use(s.securityHeaders())

	// Add common health endpoints
	name, version := initializer.GetServiceInfo()
	router.GET("/health", s.healthCheck(name, version))
	router.GET("/ready", s.readinessCheck(name))

	// Register service-specific routes
	if err := initializer.RegisterRoutes(router); err != nil {
		return fmt.Errorf("failed to register routes: %w", err)
	}

	s.router = router
	log.Println("Router initialized successfully")
	return nil
}

// initServer initializes HTTP server
func (s *ServicePattern) initServer() {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)

	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
		IdleTimeout:  s.config.Server.IdleTimeout,
	}
}

// RunWithGracefulShutdown runs the service with graceful shutdown
func (s *ServicePattern) RunWithGracefulShutdown() error {
	log.Printf("Starting service on %s", s.server.Addr)

	// Start server in a goroutine
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Setup graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Printf("Service is running on %s", s.server.Addr)

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutdown signal received")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.shutdown(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown gracefully: %v", err)
		return err
	}

	log.Println("Service stopped gracefully")
	return nil
}

// shutdown gracefully shuts down the service
func (s *ServicePattern) shutdown(ctx context.Context) error {
	log.Println("Shutting down service...")

	// Shutdown HTTP server
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	// Close database connection
	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	return nil
}

// healthCheck returns a health check handler
func (s *ServicePattern) healthCheck(serviceName, version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		responses.SendSuccessResponse(c, gin.H{
			"status":    "ok",
			"service":   serviceName,
			"version":   version,
			"timestamp": time.Now().UTC(),
		})
	}
}

// readinessCheck returns a readiness check handler
func (s *ServicePattern) readinessCheck(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check database connection
		if s.db != nil {
			sqlDB, err := s.db.DB()
			if err != nil || sqlDB.Ping() != nil {
				responses.SendErrorResponse(c, http.StatusServiceUnavailable, "database_not_ready", "Database connection failed")
				return
			}
		}

		responses.SendSuccessResponse(c, gin.H{
			"status":   "ready",
			"service":  serviceName,
			"database": "connected",
		})
	}
}

// corsMiddleware provides CORS support
func (s *ServicePattern) corsMiddleware() gin.HandlerFunc {
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
func (s *ServicePattern) securityHeaders() gin.HandlerFunc {
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