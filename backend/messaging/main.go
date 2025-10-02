package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"tchat.dev/messaging/external"
	"tchat.dev/messaging/handlers"
	"tchat.dev/messaging/models"
	"tchat.dev/messaging/repositories"
	"tchat.dev/messaging/services"
	"tchat.dev/shared/config"
	sharedModels "tchat.dev/shared/models"
)

// App represents the main messaging application
type App struct {
	config    *config.Config
	db        *gorm.DB
	router    *mux.Router
	server    *http.Server
	validator *validator.Validate

	// Services
	messageService    *services.MessageService
	dialogService     *services.DialogService
	presenceService   *services.PresenceService
	messagingService  services.MessagingService

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
		NamingStrategy: schema.NamingStrategy{SingularTable: false},
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

// initServices initializes all business logic services using real implementations
func (a *App) initServices() error {
	// Real repository implementations
	messageRepo := repositories.NewMessageRepository(a.db)
	dialogRepo := repositories.NewDialogRepository(a.db)
	presenceRepo := repositories.NewPresenceRepository(a.db)

	// External service implementations
	eventPublisher := external.NewEventPublisher()
	wsManager := external.NewWebSocketManager()
	locationService := external.NewLocationService()
	notificationService := external.NewNotificationService()
	contentModerator := external.NewContentModerator()
	mediaProcessor := external.NewMediaProcessor()

	// Message delivery service depends on notification and websocket services
	deliveryService := external.NewMessageDeliveryService(notificationService, wsManager)

	// Initialize core business services
	a.dialogService = services.NewDialogService(dialogRepo, eventPublisher, notificationService, a.db)
	a.messageService = services.NewMessageService(messageRepo, a.dialogService, deliveryService, contentModerator, mediaProcessor, eventPublisher, a.db)
	a.presenceService = services.NewPresenceService(presenceRepo, wsManager, locationService, eventPublisher, a.db)

	// Initialize messaging service
	a.messagingService = services.NewMessagingService(a.dialogService, a.messageService, nil, nil, nil)

	log.Println("Services initialized successfully with real implementations")
	return nil
}

// initHandlers initializes HTTP handlers
func (a *App) initHandlers() error {
	// Initialize handlers
	a.messagingHandlers = handlers.NewMessagingHandler(a.messagingService)

	log.Println("Handlers initialized successfully")
	return nil
}

// initRouter initializes the HTTP router with all routes
func (a *App) initRouter() error {
	router := mux.NewRouter()

	// Add test authentication middleware for development
	router.Use(a.testAuthMiddleware)

	// Health check endpoints
	router.HandleFunc("/health", a.healthCheck).Methods("GET")
	router.HandleFunc("/ready", a.readinessCheck).Methods("GET")

	// API routes
	v1 := router.PathPrefix("/api/v1").Subrouter()

	// Health endpoint for messaging
	v1.HandleFunc("/messaging/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "messaging service healthy"}`))
	}).Methods("GET")

	// Register messaging routes
	a.messagingHandlers.RegisterRoutes(v1)

	// Start messaging handler
	a.messagingHandlers.Start()

	a.router = router
	log.Println("Router initialized successfully with messaging routes")
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

// testAuthMiddleware provides test authentication for development
func (a *App) testAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health endpoints
		if r.URL.Path == "/health" || r.URL.Path == "/ready" {
			next.ServeHTTP(w, r)
			return
		}

		// Create a test user for development/testing
		testUser := &sharedModels.User{
			ID:          uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), // Fixed test UUID
			PhoneNumber: "+1234567890",
			CountryCode: "US",
			Active:      true,
		}

		// Add user to request context
		ctx := context.WithValue(r.Context(), "user", testUser)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Health check endpoint
func (a *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":    "ok",
		"service":   "messaging-service",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC(),
	}

	json.NewEncoder(w).Encode(response)
}

// Readiness check endpoint
func (a *App) readinessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check database connection
	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err != nil || sqlDB.Ping() != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			response := map[string]interface{}{
				"status":  "error",
				"message": "Database not ready",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":   "ready",
		"service":  "messaging-service",
		"database": "connected",
	}
	json.NewEncoder(w).Encode(response)
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