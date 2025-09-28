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

	// Services
	messagingService services.MessagingService

	// Handlers
	messagingHandlers *handlers.MessagingHandler
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

	// Mock notification service
	notificationService := &mockNotificationService{}

	// Mock delivery services
	deliveryService := &mockDeliveryService{}
	contentModerator := &mockContentModerator{}
	mediaProcessor := &mockMediaProcessor{}

	// Initialize services
	a.dialogService = services.NewDialogService(dialogRepo, eventPublisher, notificationService, a.db)
	a.messageService = services.NewMessageService(messageRepo, a.dialogService, deliveryService, contentModerator, mediaProcessor, eventPublisher, a.db)
	a.presenceService = services.NewPresenceService(presenceRepo, wsManager, locationService, eventPublisher, a.db)

	// Initialize messaging service
	a.messagingService = services.NewMessagingService(a.dialogService, a.messageService, nil, nil, nil)

	log.Println("Services initialized successfully")
	return nil
}

// initHandlers initializes HTTP handlers
func (a *App) initHandlers() error {
	// Initialize handlers - pass the struct, not interface
	a.messagingHandlers = handlers.NewMessagingHandler(a.messagingService)

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

	// Rate limiting (commented out for now)
	// if a.config.RateLimit.Enabled {
	// 	rateLimiter := middleware.NewRateLimiter(
	// 		a.config.RateLimit.RequestsPerMinute,
	// 		a.config.RateLimit.BurstSize,
	// 		a.config.RateLimit.CleanupInterval,
	// 	)
	// 	router.Use(rateLimiter.Middleware())
	// }

	// Health check endpoints
	router.GET("/health", a.healthCheck)
	router.GET("/ready", a.readinessCheck)

	// API routes - for now, manually register routes
	v1 := router.Group("/api/v1")
	{
		// Health endpoint for messaging
		v1.GET("/messaging/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "messaging service healthy"})
		})

		// Note: Individual route registration would be implemented here
		// For now, we'll keep it simple to get the service running
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
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "messaging-service",
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
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "error",
				"message": "Database not ready",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ready",
		"service":  "messaging-service",
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

func (m *mockMessageRepository) GetByDialogID(ctx context.Context, dialogID uuid.UUID, filters services.MessageFilters, pagination services.Pagination) ([]*models.Message, int64, error) {
	var messages []*models.Message
	var total int64
	query := m.db.WithContext(ctx).Where("dialog_id = ?", dialogID)

	// Count total
	query.Model(&models.Message{}).Count(&total)

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.PageSize
	err := query.Limit(pagination.PageSize).Offset(offset).Find(&messages).Error
	return messages, total, err
}

func (m *mockMessageRepository) Update(ctx context.Context, message *models.Message) error {
	return m.db.WithContext(ctx).Save(message).Error
}

func (m *mockMessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Delete(&models.Message{}, id).Error
}

func (m *mockMessageRepository) MarkAsRead(ctx context.Context, messageID, userID uuid.UUID) error {
	// Implementation for marking message as read
	log.Printf("Marking message %s as read for user %s", messageID, userID)
	return nil
}

func (m *mockMessageRepository) GetUnreadCount(ctx context.Context, dialogID, userID uuid.UUID) (int, error) {
	return 0, nil
}

func (m *mockMessageRepository) MarkAsDelivered(ctx context.Context, messageID, userID uuid.UUID) error {
	log.Printf("Marking message %s as delivered for user %s", messageID, userID)
	return nil
}

func (m *mockMessageRepository) SearchMessages(ctx context.Context, dialogID uuid.UUID, query string, limit int) ([]*models.Message, error) {
	var messages []*models.Message
	err := m.db.WithContext(ctx).Where("dialog_id = ? AND content ILIKE ?", dialogID, "%"+query+"%").Limit(limit).Find(&messages).Error
	return messages, err
}

func (m *mockMessageRepository) GetMessageStats(ctx context.Context, dialogID uuid.UUID) (*services.MessageStats, error) {
	return &services.MessageStats{
		TotalMessages:  100,
		TextMessages:   80,
		MediaMessages:  15,
		SystemMessages: 5,
		AverageLength:  50.5,
		MessagesPerDay: 10,
		ActiveSenders:  5,
	}, nil
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

func (m *mockDialogRepository) GetByUserID(ctx context.Context, userID uuid.UUID, filters services.DialogFilters, pagination services.Pagination) ([]*models.Dialog, int64, error) {
	var dialogs []*models.Dialog
	var total int64
	query := m.db.WithContext(ctx).Where("? = ANY(participants)", userID)

	// Count total
	query.Model(&models.Dialog{}).Count(&total)

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.PageSize
	err := query.Limit(pagination.PageSize).Offset(offset).Find(&dialogs).Error
	return dialogs, total, err
}

func (m *mockDialogRepository) GetParticipants(ctx context.Context, dialogID uuid.UUID) ([]*models.DialogParticipant, error) {
	var participants []*models.DialogParticipant
	err := m.db.WithContext(ctx).Where("dialog_id = ?", dialogID).Find(&participants).Error
	return participants, err
}

func (m *mockDialogRepository) AddParticipant(ctx context.Context, participant *models.DialogParticipant) error {
	return m.db.WithContext(ctx).Create(participant).Error
}

func (m *mockDialogRepository) RemoveParticipant(ctx context.Context, dialogID, userID uuid.UUID) error {
	return m.db.WithContext(ctx).Where("dialog_id = ? AND user_id = ?", dialogID, userID).Delete(&models.DialogParticipant{}).Error
}

func (m *mockDialogRepository) UpdateParticipant(ctx context.Context, participant *models.DialogParticipant) error {
	return m.db.WithContext(ctx).Save(participant).Error
}

func (m *mockDialogRepository) GetAdmins(ctx context.Context, dialogID uuid.UUID) ([]*models.DialogParticipant, error) {
	var admins []*models.DialogParticipant
	err := m.db.WithContext(ctx).Where("dialog_id = ? AND role IN ?", dialogID, []string{"admin", "owner"}).Find(&admins).Error
	return admins, err
}

func (m *mockDialogRepository) SearchDialogs(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.Dialog, error) {
	var dialogs []*models.Dialog
	err := m.db.WithContext(ctx).Where("? = ANY(participants) AND name ILIKE ?", userID, "%"+query+"%").Limit(limit).Find(&dialogs).Error
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

// Additional mock services needed for service constructors

type mockNotificationService struct{}

func (m *mockNotificationService) SendNotification(ctx context.Context, userID uuid.UUID, notificationType string, data map[string]interface{}) error {
	log.Printf("Notification sent to user %s: %s - %+v", userID, notificationType, data)
	return nil
}

type mockDeliveryService struct{}

func (m *mockDeliveryService) DeliverMessage(ctx context.Context, message *models.Message, recipientIDs []uuid.UUID) error {
	log.Printf("Message %s delivered to %d recipients", message.ID, len(recipientIDs))
	return nil
}

func (m *mockDeliveryService) SendPushNotification(ctx context.Context, userID uuid.UUID, message *models.Message) error {
	log.Printf("Push notification sent to user %s for message %s", userID, message.ID)
	return nil
}

type mockContentModerator struct{}

func (m *mockContentModerator) ModerateContent(ctx context.Context, content string, contentType models.MessageType) (*services.ModerationResult, error) {
	return &services.ModerationResult{
		IsApproved:      true,
		Violations:      []string{},
		Confidence:      0.95,
		FilteredContent: content,
	}, nil
}

func (m *mockContentModerator) DetectSpam(ctx context.Context, senderID uuid.UUID, content string) (*services.SpamDetectionResult, error) {
	return &services.SpamDetectionResult{
		IsSpam:     false,
		Confidence: 0.1,
		Reasons:    []string{},
	}, nil
}

type mockMediaProcessor struct{}

func (m *mockMediaProcessor) ProcessImageUpload(ctx context.Context, imageData []byte, metadata map[string]interface{}) (*services.ProcessedMedia, error) {
	return &services.ProcessedMedia{
		URL:          "https://example.com/image.jpg",
		ThumbnailURL: "https://example.com/thumb.jpg",
		Size:         int64(len(imageData)),
		Width:        800,
		Height:       600,
		Format:       "jpeg",
		Metadata:     metadata,
	}, nil
}

func (m *mockMediaProcessor) ProcessVideoUpload(ctx context.Context, videoData []byte, metadata map[string]interface{}) (*services.ProcessedMedia, error) {
	duration := 30.0
	return &services.ProcessedMedia{
		URL:          "https://example.com/video.mp4",
		ThumbnailURL: "https://example.com/video_thumb.jpg",
		Size:         int64(len(videoData)),
		Width:        1920,
		Height:       1080,
		Duration:     &duration,
		Format:       "mp4",
		Metadata:     metadata,
	}, nil
}

func (m *mockMediaProcessor) ProcessAudioUpload(ctx context.Context, audioData []byte, metadata map[string]interface{}) (*services.ProcessedMedia, error) {
	duration := 60.0
	return &services.ProcessedMedia{
		URL:      "https://example.com/audio.mp3",
		Size:     int64(len(audioData)),
		Duration: &duration,
		Format:   "mp3",
		Metadata: metadata,
	}, nil
}

func (m *mockMediaProcessor) GenerateThumbnail(ctx context.Context, mediaURL string, mediaType string) (string, error) {
	return "https://example.com/thumbnail.jpg", nil
}

func main() {
	// Load configuration with messaging service specific port (8082)
	cfg, err := config.LoadWithServicePort("messaging", 8082)
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