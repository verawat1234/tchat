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
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"tchat.dev/shared/config"
	"tchat.dev/shared/database"
	"tchat.dev/shared/events"
	"tchat.dev/shared/middleware"
	"tchat.dev/shared/responses"
	sharedModels "tchat.dev/shared/models"

	"tchat.dev/notification/handlers"
	"tchat.dev/notification/models"
	"tchat.dev/notification/providers"
	"tchat.dev/notification/repositories"
	"tchat.dev/notification/services"
	notificationConfig "tchat.dev/notification/config"
)

// App represents the main notification application
type App struct {
	config    *config.Config
	db        *gorm.DB
	router    *gin.Engine
	server    *http.Server
	validator *validator.Validate

	// Services
	notificationService services.NotificationService
	cacheService       services.CacheService
	eventService       services.EventService

	// Repositories
	notificationRepo repositories.NotificationRepository
	templateRepo     repositories.TemplateRepository

	// Handlers
	notificationHandler *handlers.NotificationHandler
}

// NewApp creates a new notification application instance
func NewApp(cfg *config.Config) *App {
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

	log.Printf("Notification service initialized successfully on port %d", a.config.Server.Port)
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

	// Run auto-migrations
	if err := a.runMigrations(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	a.db = db
	log.Println("Database connection established and migrations completed")
	return nil
}

// runMigrations runs database migrations for notification service models
func (a *App) runMigrations(db *gorm.DB) error {
	// Migration for notification service models
	return db.AutoMigrate(
		&sharedModels.Event{},
		&models.Notification{},
		&models.NotificationTemplate{},
		&models.NotificationSubscription{},
		&models.NotificationPreferences{},
	)
}

// initRepositories initializes repository instances
func (a *App) initRepositories() error {
	a.notificationRepo = repositories.NewNotificationRepository(a.db)
	a.templateRepo = repositories.NewTemplateRepository(a.db)

	log.Println("Repositories initialized successfully")
	return nil
}

// initServices initializes service instances
func (a *App) initServices() error {
	// Initialize cache service
	var err error
	redisURL := os.Getenv("REDIS_URL")
	if redisURL != "" {
		a.cacheService, err = services.NewRedisCacheService(redisURL)
		if err != nil {
			log.Printf("Failed to connect to Redis, falling back to in-memory cache: %v", err)
			a.cacheService = services.NewInMemoryCacheService()
		} else {
			log.Println("Redis cache service initialized")
		}
	} else {
		a.cacheService = services.NewInMemoryCacheService()
		log.Println("In-memory cache service initialized")
	}

	// Initialize event service
	a.eventService = services.NewEventService(true) // Enable async processing
	log.Println("Event service initialized")

	// Initialize notification service
	notificationConfig := services.DefaultNotificationConfig()
	a.notificationService = services.NewNotificationService(
		a.notificationRepo,
		a.templateRepo,
		a.cacheService,
		a.eventService,
		notificationConfig,
	)

	// Initialize and register providers
	if err := a.initProviders(); err != nil {
		return fmt.Errorf("failed to initialize providers: %w", err)
	}

	log.Println("Services initialized successfully")
	return nil
}

// initProviders initializes notification providers
func (a *App) initProviders() error {
	// Initialize Email Provider
	if smtpHost := os.Getenv("SMTP_HOST"); smtpHost != "" {
		emailConfig := providers.EmailConfig{
			SMTPHost:     smtpHost,
			SMTPPort:     getEnvOrDefault("SMTP_PORT", "587"),
			SMTPUsername: os.Getenv("SMTP_USERNAME"),
			SMTPPassword: os.Getenv("SMTP_PASSWORD"),
			FromEmail:    os.Getenv("FROM_EMAIL"),
			FromName:     getEnvOrDefault("FROM_NAME", "Tchat Notifications"),
		}
		emailProvider := providers.NewEmailProvider(emailConfig)
		if err := emailProvider.ValidateConfig(); err != nil {
			log.Printf("Email provider configuration invalid: %v", err)
		} else {
			// Register email provider with service
			log.Println("Email provider initialized")
		}
	}

	// Initialize SMS Provider (Twilio)
	if twilioSID := os.Getenv("TWILIO_ACCOUNT_SID"); twilioSID != "" {
		smsConfig := providers.SMSConfig{
			TwilioAccountSID: twilioSID,
			TwilioAuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
			TwilioFromNumber: os.Getenv("TWILIO_FROM_NUMBER"),
		}
		smsProvider := providers.NewSMSProvider(smsConfig)
		if err := smsProvider.ValidateConfig(); err != nil {
			log.Printf("SMS provider configuration invalid: %v", err)
		} else {
			log.Println("SMS provider initialized")
		}
	}

	// Initialize Push Provider
	if fcmKey := os.Getenv("FCM_SERVER_KEY"); fcmKey != "" {
		pushConfig := providers.PushConfig{
			FCMServerKey: fcmKey,
			APNSKeyID:    os.Getenv("APNS_KEY_ID"),
			APNSTeamID:   os.Getenv("APNS_TEAM_ID"),
			APNSKeyPath:  os.Getenv("APNS_KEY_PATH"),
			BundleID:     os.Getenv("BUNDLE_ID"),
		}
		pushProvider := providers.NewPushProvider(pushConfig)
		if err := pushProvider.ValidateConfig(); err != nil {
			log.Printf("Push provider configuration invalid: %v", err)
		} else {
			log.Println("Push provider initialized")
		}
	}

	// Initialize In-App Provider
	wsManager := providers.NewSimpleWebSocketManager()
	inAppProvider := providers.NewInAppProvider(wsManager, nil) // Repository will be added later
	if err := inAppProvider.ValidateConfig(); err != nil {
		log.Printf("In-app provider configuration invalid: %v", err)
	} else {
		log.Println("In-app provider initialized")
	}

	return nil
}

// initHandlers initializes HTTP handlers
func (a *App) initHandlers() error {
	// Create proper zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Create logger adapter for event bus
	loggerAdapter := &zapLoggerAdapter{logger: logger}

	// Create event bus with proper configuration
	eventBusConfig := &events.EventBusConfig{
		MaxConcurrency:  10,
		BufferSize:      100,
		MaxRetries:      3,
		RetryDelay:      time.Second,
		HandlerTimeout:  30 * time.Second,
		EnableMetrics:   true,
	}
	eventBus := events.NewEventBus(eventBusConfig, loggerAdapter)

	// Create notification config from existing config
	notifConfig := &notificationConfig.NotificationConfig{
		Config: a.config,
	}
	// Set defaults for notification-specific settings
	notifConfig.Notification.QueueSize = 1000
	notifConfig.Notification.WorkerCount = 5
	notifConfig.Notification.RetryAttempts = 3
	notifConfig.Notification.RetryDelay = 5
	notifConfig.Notification.RateLimitPerUser = 10
	notifConfig.Notification.RateLimitGlobal = 1000

	// Initialize notification handler
	a.notificationHandler = handlers.NewNotificationHandler(
		a.notificationService,
		logger,
		eventBus,
		notifConfig,
	)

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
	router.Use(middleware.CORS())
	router.Use(middleware.SecurityHeaders())

	// Health check endpoints
	router.GET("/health", a.healthCheck)
    router.GET("/v1/healthcheck", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":  "ok",
            "version": "1.0.0",
        })
    })
	router.GET("/ready", a.readinessCheck)

	// Register notification routes
	a.notificationHandler.RegisterRoutes(router)

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
	log.Printf("Starting notification service on %s", a.server.Addr)

	// Start server in a goroutine
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Notification service is running on %s", a.server.Addr)
	return nil
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down notification service...")

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

	// Close cache service if it supports closing
	if redisCacheService, ok := a.cacheService.(*services.RedisCacheService); ok {
		redisCacheService.Close()
	}

	log.Println("Notification service shutdown completed")
	return nil
}

// Health check endpoint
func (a *App) healthCheck(c *gin.Context) {
	responses.SendSuccessResponse(c, gin.H{
		"status":    "ok",
		"service":   "notification-service",
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
			responses.SendErrorResponse(c, http.StatusServiceUnavailable, "Database not ready", "Database connection failed")
			return
		}
	}

	responses.SendSuccessResponse(c, gin.H{
		"status":   "ready",
		"service":  "notification-service",
		"database": "connected",
		"cache":    "available",
		"events":   "available",
	})
}

// Helper function to get environment variable with default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Simplified implementations for dependencies
// In production, these would be proper implementations

type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, fields ...interface{}) {
	log.Printf("[INFO] %s %v", msg, fields)
}

func (l *SimpleLogger) Error(msg string, fields ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, fields)
}

func (l *SimpleLogger) Debug(msg string, fields ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, fields)
}

func (l *SimpleLogger) Warn(msg string, fields ...interface{}) {
	log.Printf("[WARN] %s %v", msg, fields)
}

type SimpleEventBus struct{}

func (e *SimpleEventBus) Publish(event string, data interface{}) {
	log.Printf("Event published: %s with data: %+v", event, data)
}

type SimpleNotificationConfig struct{}

// zapLoggerAdapter adapts zap.Logger to events.Logger interface
type zapLoggerAdapter struct {
	logger *zap.Logger
}

func (z *zapLoggerAdapter) Info(msg string, fields ...interface{}) {
	z.logger.Info(msg, zap.Any("fields", fields))
}

func (z *zapLoggerAdapter) Error(msg string, err error, fields ...interface{}) {
	z.logger.Error(msg, zap.Error(err), zap.Any("fields", fields))
}

func (z *zapLoggerAdapter) Debug(msg string, fields ...interface{}) {
	z.logger.Debug(msg, zap.Any("fields", fields))
}

func (z *zapLoggerAdapter) Warn(msg string, fields ...interface{}) {
	z.logger.Warn(msg, zap.Any("fields", fields))
}

func main() {
	// Load configuration with notification service specific port (8089)
	cfg, err := config.LoadWithServicePort("notification", 8089)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
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

	log.Println("Notification service stopped")
}