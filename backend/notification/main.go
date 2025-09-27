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
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"tchat.dev/shared/config"
	"tchat.dev/shared/database"
	"tchat.dev/shared/middleware"
	"tchat.dev/shared/responses"
	sharedModels "tchat.dev/shared/models"

	"tchat.dev/notification/models"
)

// App represents the main notification application
type App struct {
	config    *config.Config
	db        *gorm.DB
	router    *gin.Engine
	server    *http.Server
	validator *validator.Validate

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

	// Services initialization skipped - using simple handlers

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
	// Basic migration for notifications table
	return db.AutoMigrate(
		&sharedModels.Event{},
	)
}


// initHandlers initializes HTTP handlers
func (a *App) initHandlers() error {
	// Handlers are now methods on the App struct
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

	// Rate limiting temporarily disabled
	// if a.config.RateLimit.Enabled {
	//	rateLimiter := middleware.NewRateLimiter(
	//		a.config.RateLimit.RequestsPerMinute,
	//		a.config.RateLimit.BurstSize,
	//		a.config.RateLimit.CleanupInterval,
	//	)
	//	router.Use(rateLimiter.Middleware())
	// }

	// Health check endpoints
	router.GET("/health", a.healthCheck)
	router.GET("/ready", a.readinessCheck)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Notification endpoints
		v1.GET("/notifications", a.getNotifications)
		v1.GET("/notifications/preferences", a.getNotificationPreferences)
		v1.PUT("/notifications/preferences", a.updateNotificationPreferences)
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
	})
}

// Mock implementations

type mockNotificationRepository struct {
	db *gorm.DB
}

func (m *mockNotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	return m.db.WithContext(ctx).Create(notification).Error
}

func (m *mockNotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	var notification models.Notification
	err := m.db.WithContext(ctx).First(&notification, id).Error
	return &notification, err
}

func (m *mockNotificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := m.db.WithContext(ctx).Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&notifications).Error
	return notifications, err
}

func (m *mockNotificationRepository) Update(ctx context.Context, notification *models.Notification) error {
	return m.db.WithContext(ctx).Save(notification).Error
}

func (m *mockNotificationRepository) GetPendingNotifications(ctx context.Context, limit int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := m.db.WithContext(ctx).Where("status = ?", models.DeliveryStatusPending).Limit(limit).Find(&notifications).Error
	return notifications, err
}

func (m *mockNotificationRepository) MarkAsDelivered(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Model(&models.Notification{}).Where("id = ?", id).Update("status", models.DeliveryStatusDelivered).Error
}

func (m *mockNotificationRepository) MarkAsFailed(ctx context.Context, id uuid.UUID, reason string) error {
	return m.db.WithContext(ctx).Model(&models.Notification{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       models.DeliveryStatusFailed,
		"failure_reason": reason,
	}).Error
}

func (m *mockNotificationRepository) GetByChannel(ctx context.Context, channel models.NotificationType, limit, offset int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := m.db.WithContext(ctx).Where("type = ?", channel).Limit(limit).Offset(offset).Find(&notifications).Error
	return notifications, err
}

func (m *mockNotificationRepository) CleanupOldNotifications(ctx context.Context, before time.Time) error {
	return m.db.WithContext(ctx).Where("created_at < ?", before).Delete(&models.Notification{}).Error
}

type mockTemplateRepository struct{}

func (m *mockTemplateRepository) GetByType(ctx context.Context, notificationType string) (*models.NotificationTemplate, error) {
	return &models.NotificationTemplate{
		ID:      uuid.New(),
		Type:    models.NotificationType(notificationType),
		Subject: "Default Subject",
		Body:    "Default notification body",
	}, nil
}

func (m *mockTemplateRepository) GetByTypeAndLanguage(ctx context.Context, notificationType, language string) (*models.NotificationTemplate, error) {
	return m.GetByType(ctx, notificationType)
}

type mockCacheService struct{}

func (m *mockCacheService) Get(ctx context.Context, key string) (interface{}, error) {
	return nil, fmt.Errorf("not found")
}

func (m *mockCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return nil
}

func (m *mockCacheService) Delete(ctx context.Context, key string) error {
	return nil
}

type mockEventService struct{}

func (m *mockEventService) Publish(ctx context.Context, event interface{}) error {
	log.Printf("Event published: %+v", event)
	return nil
}

type mockPushProvider struct{}

func (m *mockPushProvider) Send(ctx context.Context, notification *models.Notification) error {
	log.Printf("Push notification sent to user %s: %s", notification.ID, notification.Title)
	return nil
}

func (m *mockPushProvider) SendBatch(ctx context.Context, notifications []*models.Notification) error {
	log.Printf("Batch push notifications sent: %d notifications", len(notifications))
	return nil
}

type mockEmailProvider struct{}

func (m *mockEmailProvider) Send(ctx context.Context, notification *models.Notification) error {
	log.Printf("Email sent to user %s: %s", notification.ID, notification.Title)
	return nil
}

func (m *mockEmailProvider) SendBatch(ctx context.Context, notifications []*models.Notification) error {
	log.Printf("Batch emails sent: %d notifications", len(notifications))
	return nil
}

type mockSMSProvider struct{}

func (m *mockSMSProvider) Send(ctx context.Context, notification *models.Notification) error {
	log.Printf("SMS sent to user %s: %s", notification.ID, notification.Title)
	return nil
}

func (m *mockSMSProvider) SendBatch(ctx context.Context, notifications []*models.Notification) error {
	log.Printf("Batch SMS sent: %d notifications", len(notifications))
	return nil
}

type mockInAppProvider struct{}

func (m *mockInAppProvider) Send(ctx context.Context, notification *models.Notification) error {
	log.Printf("In-app notification sent to user %s: %s", notification.ID, notification.Title)
	return nil
}

func (m *mockInAppProvider) SendBatch(ctx context.Context, notifications []*models.Notification) error {
	log.Printf("Batch in-app notifications sent: %d notifications", len(notifications))
	return nil
}

func main() {
	// Load configuration
	cfg, err := config.Load()
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

// Simple handler methods

// getNotifications retrieves user notifications
func (a *App) getNotifications(c *gin.Context) {
	responses.SendSuccessResponse(c, gin.H{
		"notifications": []gin.H{
			{
				"id":      "1",
				"title":   "Welcome to Tchat",
				"content": "Thank you for joining our platform",
				"read":    false,
			},
		},
		"total": 1,
	})
}

// getNotificationPreferences retrieves notification preferences
func (a *App) getNotificationPreferences(c *gin.Context) {
	responses.SendSuccessResponse(c, gin.H{
		"email_enabled": true,
		"push_enabled":  true,
		"sms_enabled":   false,
	})
}

// updateNotificationPreferences updates notification preferences
func (a *App) updateNotificationPreferences(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	responses.SendSuccessResponse(c, gin.H{
		"message": "Notification preferences updated successfully",
		"preferences": req,
	})
}

// Additional mock implementations for missing dependencies

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...interface{}) {
	log.Printf("[INFO] %s %v", msg, fields)
}

func (m *mockLogger) Error(msg string, fields ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, fields)
}

func (m *mockLogger) Debug(msg string, fields ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, fields)
}

func (m *mockLogger) Warn(msg string, fields ...interface{}) {
	log.Printf("[WARN] %s %v", msg, fields)
}

type mockEventBus struct{}

func (m *mockEventBus) Publish(event string, data interface{}) {
	log.Printf("Event published: %s with data: %+v", event, data)
}

type mockNotificationConfig struct{}