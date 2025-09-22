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

	"tchat.dev/notification/handlers"
	"tchat.dev/notification/models"
	"tchat.dev/notification/services"
	"tchat.dev/shared/config"
	"tchat.dev/shared/database"
	"tchat.dev/shared/middleware"
	"tchat.dev/shared/responses"
	sharedModels "tchat.dev/shared/models"
)

// App represents the main notification application
type App struct {
	config    *config.Config
	db        *gorm.DB
	router    *gin.Engine
	server    *http.Server
	validator *validator.Validate

	// Services
	notificationService *services.NotificationService

	// Handlers
	notificationHandlers *handlers.NotificationHandler
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
	return db.AutoMigrate(
		&models.Notification{},
		&sharedModels.Event{},
	)
}

// initServices initializes all business logic services
func (a *App) initServices() error {
	// Mock repositories and services
	notificationRepo := &mockNotificationRepository{db: a.db}
	templateRepo := &mockTemplateRepository{}
	cacheService := &mockCacheService{}
	eventService := &mockEventService{}

	// Mock channel providers
	providers := map[models.NotificationChannel]services.ChannelProvider{
		models.ChannelPush:  &mockPushProvider{},
		models.ChannelEmail: &mockEmailProvider{},
		models.ChannelSMS:   &mockSMSProvider{},
		models.ChannelInApp: &mockInAppProvider{},
	}

	// Initialize notification service
	a.notificationService = services.NewNotificationService(
		notificationRepo,
		templateRepo,
		cacheService,
		eventService,
		providers,
		services.DefaultNotificationConfig(),
	)

	log.Println("Services initialized successfully")
	return nil
}

// initHandlers initializes HTTP handlers
func (a *App) initHandlers() error {
	// Initialize handlers
	a.notificationHandlers = handlers.NewNotificationHandler(a.notificationService)

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

	// Rate limiting
	if a.config.RateLimit.Enabled {
		rateLimiter := middleware.NewRateLimiter(
			a.config.RateLimit.RequestsPerMinute,
			a.config.RateLimit.BurstSize,
			a.config.RateLimit.CleanupInterval,
		)
		router.Use(rateLimiter.Middleware())
	}

	// Health check endpoints
	router.GET("/health", a.healthCheck)
	router.GET("/ready", a.readinessCheck)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Notification routes
		handlers.RegisterNotificationRoutes(v1, a.notificationService)
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
	responses.SuccessResponse(c, gin.H{
		"status":    "ok",
		"service":   "notification",
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
			responses.ErrorResponse(c, http.StatusServiceUnavailable, "Database not ready", "Database connection failed")
			return
		}
	}

	responses.SuccessResponse(c, gin.H{
		"status":   "ready",
		"service":  "notification",
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
	err := m.db.WithContext(ctx).Where("status = ?", models.StatusPending).Limit(limit).Find(&notifications).Error
	return notifications, err
}

func (m *mockNotificationRepository) MarkAsDelivered(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Model(&models.Notification{}).Where("id = ?", id).Update("status", models.StatusDelivered).Error
}

func (m *mockNotificationRepository) MarkAsFailed(ctx context.Context, id uuid.UUID, reason string) error {
	return m.db.WithContext(ctx).Model(&models.Notification{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       models.StatusFailed,
		"failure_reason": reason,
	}).Error
}

func (m *mockNotificationRepository) GetByChannel(ctx context.Context, channel models.NotificationChannel, limit, offset int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := m.db.WithContext(ctx).Where("channel = ?", channel).Limit(limit).Offset(offset).Find(&notifications).Error
	return notifications, err
}

func (m *mockNotificationRepository) CleanupOldNotifications(ctx context.Context, before time.Time) error {
	return m.db.WithContext(ctx).Where("created_at < ?", before).Delete(&models.Notification{}).Error
}

type mockTemplateRepository struct{}

func (m *mockTemplateRepository) GetByType(ctx context.Context, notificationType string) (*services.NotificationTemplate, error) {
	return &services.NotificationTemplate{
		ID:      uuid.New(),
		Type:    notificationType,
		Subject: "Default Subject",
		Body:    "Default notification body",
	}, nil
}

func (m *mockTemplateRepository) GetByTypeAndLanguage(ctx context.Context, notificationType, language string) (*services.NotificationTemplate, error) {
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
	log.Printf("Push notification sent to user %s: %s", notification.UserID, notification.Title)
	return nil
}

func (m *mockPushProvider) SendBatch(ctx context.Context, notifications []*models.Notification) error {
	log.Printf("Batch push notifications sent: %d notifications", len(notifications))
	return nil
}

type mockEmailProvider struct{}

func (m *mockEmailProvider) Send(ctx context.Context, notification *models.Notification) error {
	log.Printf("Email sent to user %s: %s", notification.UserID, notification.Title)
	return nil
}

func (m *mockEmailProvider) SendBatch(ctx context.Context, notifications []*models.Notification) error {
	log.Printf("Batch emails sent: %d notifications", len(notifications))
	return nil
}

type mockSMSProvider struct{}

func (m *mockSMSProvider) Send(ctx context.Context, notification *models.Notification) error {
	log.Printf("SMS sent to user %s: %s", notification.UserID, notification.Title)
	return nil
}

func (m *mockSMSProvider) SendBatch(ctx context.Context, notifications []*models.Notification) error {
	log.Printf("Batch SMS sent: %d notifications", len(notifications))
	return nil
}

type mockInAppProvider struct{}

func (m *mockInAppProvider) Send(ctx context.Context, notification *models.Notification) error {
	log.Printf("In-app notification sent to user %s: %s", notification.UserID, notification.Title)
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