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

	"tchat.dev/messaging/handlers"
	"tchat.dev/messaging/models"
	"tchat.dev/messaging/services"
	"tchat.dev/shared/config"
	"tchat.dev/shared/database"
	"tchat.dev/shared/middleware"
	"tchat.dev/shared/responses"
	sharedModels "tchat.dev/shared/models"
)

// App represents the main messaging application
type App struct {
	config    *config.Config
	db        *gorm.DB
	router    *gin.Engine
	server    *http.Server
	validator *validator.Validate

	// Services
	messageService  *services.MessageService
	dialogService   *services.DialogService
	presenceService *services.PresenceService

	// Handlers
	messagingHandlers *handlers.MessagingHandler
	websocketHandler  *handlers.WebSocketHandler
}

// NewApp creates a new messaging application instance
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

	log.Printf("Messaging service initialized successfully on port %d", a.config.Server.Port)
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

// runMigrations runs database migrations for messaging service models
func (a *App) runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Dialog{},
		&models.Message{},
		&models.Presence{},
		&sharedModels.Event{},
	)
}

// initServices initializes all business logic services
func (a *App) initServices() error {
	// Mock event publisher for now
	eventPublisher := &mockEventPublisher{}

	// Mock repositories for now - these would be implemented based on your repository pattern
	messageRepo := &mockMessageRepository{db: a.db}
	dialogRepo := &mockDialogRepository{db: a.db}
	presenceRepo := &mockPresenceRepository{db: a.db}

	// Mock WebSocket manager for now
	wsManager := &mockWebSocketManager{}

	// Mock location service for now
	locationService := &mockLocationService{}

	// Initialize services
	a.messageService = services.NewMessageService(messageRepo, dialogRepo, eventPublisher, a.db)
	a.dialogService = services.NewDialogService(dialogRepo, eventPublisher, a.db)
	a.presenceService = services.NewPresenceService(presenceRepo, wsManager, locationService, eventPublisher, a.db)

	log.Println("Services initialized successfully")
	return nil
}

// initHandlers initializes HTTP handlers
func (a *App) initHandlers() error {
	// Initialize handlers
	a.messagingHandlers = handlers.NewMessagingHandler(
		a.messageService,
		a.dialogService,
	)

	a.websocketHandler = handlers.NewWebSocketHandler(
		a.messageService,
		a.presenceService,
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

	// WebSocket endpoint
	router.GET("/websocket", a.websocketHandler.HandleWebSocket)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Messaging routes
		handlers.RegisterMessagingRoutes(v1, a.messageService, a.dialogService)

		// WebSocket routes
		handlers.RegisterWebSocketRoutes(v1, a.messageService, a.presenceService)
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
	log.Printf("Starting messaging service on %s", a.server.Addr)

	// Start server in a goroutine
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Messaging service is running on %s", a.server.Addr)
	return nil
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down messaging service...")

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

	log.Println("Messaging service shutdown completed")
	return nil
}

// Health check endpoint
func (a *App) healthCheck(c *gin.Context) {
	responses.SuccessResponse(c, gin.H{
		"status":    "ok",
		"service":   "messaging",
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
		"service":  "messaging",
		"database": "connected",
	})
}

// Mock implementations for repositories and services

type mockEventPublisher struct{}

func (m *mockEventPublisher) Publish(ctx context.Context, event *sharedModels.Event) error {
	log.Printf("Event published: %s - %s", event.Type, event.Subject)
	return nil
}

type mockMessageRepository struct {
	db *gorm.DB
}

func (m *mockMessageRepository) Create(ctx context.Context, message *models.Message) error {
	return m.db.WithContext(ctx).Create(message).Error
}

func (m *mockMessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	var message models.Message
	err := m.db.WithContext(ctx).First(&message, id).Error
	return &message, err
}

func (m *mockMessageRepository) GetByDialogID(ctx context.Context, dialogID uuid.UUID, limit, offset int) ([]*models.Message, error) {
	var messages []*models.Message
	err := m.db.WithContext(ctx).Where("dialog_id = ?", dialogID).Limit(limit).Offset(offset).Find(&messages).Error
	return messages, err
}

func (m *mockMessageRepository) Update(ctx context.Context, message *models.Message) error {
	return m.db.WithContext(ctx).Save(message).Error
}

func (m *mockMessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Delete(&models.Message{}, id).Error
}

func (m *mockMessageRepository) MarkAsRead(ctx context.Context, messageIDs []uuid.UUID, userID uuid.UUID) error {
	// Implementation for marking messages as read
	return nil
}

func (m *mockMessageRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID, dialogID uuid.UUID) (int64, error) {
	return 0, nil
}

type mockDialogRepository struct {
	db *gorm.DB
}

func (m *mockDialogRepository) Create(ctx context.Context, dialog *models.Dialog) error {
	return m.db.WithContext(ctx).Create(dialog).Error
}

func (m *mockDialogRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Dialog, error) {
	var dialog models.Dialog
	err := m.db.WithContext(ctx).First(&dialog, id).Error
	return &dialog, err
}

func (m *mockDialogRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Dialog, error) {
	var dialogs []*models.Dialog
	err := m.db.WithContext(ctx).Where("? = ANY(participant_ids)", userID).Limit(limit).Offset(offset).Find(&dialogs).Error
	return dialogs, err
}

func (m *mockDialogRepository) Update(ctx context.Context, dialog *models.Dialog) error {
	return m.db.WithContext(ctx).Save(dialog).Error
}

func (m *mockDialogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Delete(&models.Dialog{}, id).Error
}

type mockPresenceRepository struct {
	db *gorm.DB
}

func (m *mockPresenceRepository) Create(ctx context.Context, presence *models.Presence) error {
	return m.db.WithContext(ctx).Create(presence).Error
}

func (m *mockPresenceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Presence, error) {
	var presence models.Presence
	err := m.db.WithContext(ctx).Where("user_id = ?", userID).First(&presence).Error
	return &presence, err
}

func (m *mockPresenceRepository) Update(ctx context.Context, presence *models.Presence) error {
	return m.db.WithContext(ctx).Save(presence).Error
}

func (m *mockPresenceRepository) GetByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.Presence, error) {
	var presences []*models.Presence
	err := m.db.WithContext(ctx).Where("user_id IN ?", userIDs).Find(&presences).Error
	return presences, err
}

func (m *mockPresenceRepository) GetOnlineUsers(ctx context.Context, limit int) ([]*models.Presence, error) {
	var presences []*models.Presence
	err := m.db.WithContext(ctx).Where("is_online = ?", true).Limit(limit).Find(&presences).Error
	return presences, err
}

func (m *mockPresenceRepository) CleanupStalePresence(ctx context.Context, staleThreshold time.Duration) error {
	threshold := time.Now().Add(-staleThreshold)
	return m.db.WithContext(ctx).Model(&models.Presence{}).
		Where("last_updated < ? AND is_online = ?", threshold, true).
		Update("is_online", false).Error
}

func (m *mockPresenceRepository) GetPresenceStats(ctx context.Context) (*services.PresenceStats, error) {
	return &services.PresenceStats{
		TotalUsers:      1000,
		OnlineUsers:     250,
		AwayUsers:       50,
		BusyUsers:       25,
		OfflineUsers:    675,
		AverageUptime:   4 * time.Hour,
		PeakOnlineTime:  time.Now().Add(-2 * time.Hour),
		PeakOnlineCount: 300,
	}, nil
}

type mockWebSocketManager struct{}

func (m *mockWebSocketManager) BroadcastToUser(ctx context.Context, userID uuid.UUID, message interface{}) error {
	log.Printf("Broadcasting to user %s: %+v", userID, message)
	return nil
}

func (m *mockWebSocketManager) BroadcastToUsers(ctx context.Context, userIDs []uuid.UUID, message interface{}) error {
	log.Printf("Broadcasting to %d users: %+v", len(userIDs), message)
	return nil
}

func (m *mockWebSocketManager) GetConnectedUsers(ctx context.Context) []uuid.UUID {
	// Return some mock connected users
	return []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
}

func (m *mockWebSocketManager) IsUserConnected(ctx context.Context, userID uuid.UUID) bool {
	// Mock implementation - would check actual connections
	return true
}

type mockLocationService struct{}

func (m *mockLocationService) UpdateUserLocation(ctx context.Context, userID uuid.UUID, location models.Location) error {
	log.Printf("Updated location for user %s: %+v", userID, location)
	return nil
}

func (m *mockLocationService) GetNearbyUsers(ctx context.Context, userID uuid.UUID, radius float64) ([]uuid.UUID, error) {
	// Return some mock nearby users
	return []uuid.UUID{uuid.New(), uuid.New()}, nil
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

	log.Println("Messaging service stopped")
}