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

	"github.com/gocql/gocql"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"tchat.dev/messaging/external"
	"tchat.dev/messaging/handlers"
	"tchat.dev/messaging/repositories"
	"tchat.dev/messaging/services"
	"tchat.dev/shared/config"
	sharedModels "tchat.dev/shared/models"
)

// App represents the main messaging application
type App struct {
	config    *config.Config
	session   *gocql.Session
	router    *mux.Router
	server    *http.Server
	validator *validator.Validate

	// Services
	messageService    *services.MessageService
	dialogService     *services.DialogService
	presenceService   *services.PresenceService
	messagingService  services.MessagingService
	wsManager         services.WebSocketManager

	// Handlers
	messagingHandlers *handlers.MessagingHandler
	wsHandler         *handlers.WebSocketHandler
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

// initDatabase initializes the ScyllaDB connection and creates keyspace/tables
func (a *App) initDatabase() error {
	// Get ScyllaDB connection details from environment
	scyllaHost := os.Getenv("SCYLLA_HOST")
	if scyllaHost == "" {
		scyllaHost = "localhost" // Default to localhost for local development
	}

	scyllaPort := os.Getenv("SCYLLA_PORT")
	if scyllaPort == "" {
		scyllaPort = "9042"
	}

	scyllaUsername := os.Getenv("SCYLLA_USERNAME")
	if scyllaUsername == "" {
		scyllaUsername = "cassandra"
	}

	scyllaPassword := os.Getenv("SCYLLA_PASSWORD")
	if scyllaPassword == "" {
		scyllaPassword = "cassandra"
	}

	// Create ScyllaDB cluster configuration
	cluster := gocql.NewCluster(scyllaHost)
	cluster.Port = 9042
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: scyllaUsername,
		Password: scyllaPassword,
	}
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.DisableInitialHostLookup = true // Disable host discovery for Railway environment

	// Create keyspace if it doesn't exist
	keyspace := "messaging"
	cluster.Keyspace = "system"

	systemSession, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("failed to connect to ScyllaDB: %w", err)
	}

	// Create keyspace with SimpleStrategy and replication factor 1
	if err := systemSession.Query(fmt.Sprintf(`
		CREATE KEYSPACE IF NOT EXISTS %s
		WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}
	`, keyspace)).Exec(); err != nil {
		systemSession.Close()
		return fmt.Errorf("failed to create keyspace: %w", err)
	}
	systemSession.Close()

	// Connect to the messaging keyspace
	cluster.Keyspace = keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("failed to connect to messaging keyspace: %w", err)
	}

	// Create tables
	if err := a.createTables(session); err != nil {
		session.Close()
		return fmt.Errorf("failed to create tables: %w", err)
	}

	a.session = session
	log.Println("ScyllaDB connection established and tables created")
	return nil
}

// createTables creates ScyllaDB tables for messaging service
func (a *App) createTables(session *gocql.Session) error {
	// Create dialogs table
	if err := session.Query(`
		CREATE TABLE IF NOT EXISTS dialogs (
			id UUID PRIMARY KEY,
			type TEXT,
			participant_ids SET<UUID>,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			last_message_id UUID,
			last_message_text TEXT,
			last_message_at TIMESTAMP
		)
	`).Exec(); err != nil {
		return fmt.Errorf("failed to create dialogs table: %w", err)
	}

	// Create messages table with time-series optimized structure
	if err := session.Query(`
		CREATE TABLE IF NOT EXISTS messages (
			dialog_id UUID,
			id UUID,
			sender_id UUID,
			text TEXT,
			media_url TEXT,
			message_type TEXT,
			status TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			deleted_at TIMESTAMP,
			PRIMARY KEY (dialog_id, created_at, id)
		) WITH CLUSTERING ORDER BY (created_at DESC, id DESC)
	`).Exec(); err != nil {
		return fmt.Errorf("failed to create messages table: %w", err)
	}

	// Create presence table
	if err := session.Query(`
		CREATE TABLE IF NOT EXISTS presence (
			user_id UUID PRIMARY KEY,
			status TEXT,
			last_seen TIMESTAMP,
			location TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)
	`).Exec(); err != nil {
		return fmt.Errorf("failed to create presence table: %w", err)
	}

	// Create user_dialogs table for efficient user-to-dialogs queries
	if err := session.Query(`
		CREATE TABLE IF NOT EXISTS user_dialogs (
			user_id UUID,
			dialog_id UUID,
			type TEXT,
			last_message_at TIMESTAMP,
			is_archived BOOLEAN,
			is_muted BOOLEAN,
			unread_count INT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			PRIMARY KEY (user_id, last_message_at, dialog_id)
		) WITH CLUSTERING ORDER BY (last_message_at DESC, dialog_id DESC)
	`).Exec(); err != nil {
		return fmt.Errorf("failed to create user_dialogs table: %w", err)
	}

	log.Println("ScyllaDB tables created successfully")
	return nil
}

// initServices initializes all business logic services using real implementations
func (a *App) initServices() error {
	// Real repository implementations with ScyllaDB session
	messageRepo := repositories.NewScyllaMessageRepository(a.session)
	dialogRepo := repositories.NewScyllaDialogRepository(a.session)
	presenceRepo := repositories.NewScyllaPresenceRepository(a.session)

	// External service implementations
	eventPublisher := external.NewEventPublisher()
	a.wsManager = external.NewWebSocketManager()
	locationService := external.NewLocationService()
	notificationService := external.NewNotificationService()
	contentModerator := external.NewContentModerator()
	mediaProcessor := external.NewMediaProcessor()

	// Message delivery service depends on notification and websocket services
	deliveryService := external.NewMessageDeliveryService(notificationService, a.wsManager)

	// Initialize core business services with ScyllaDB session
	// Note: passing nil for db parameter as we're using ScyllaDB repositories
	a.dialogService = services.NewDialogService(dialogRepo, eventPublisher, notificationService, nil)
	a.messageService = services.NewMessageService(messageRepo, a.dialogService, deliveryService, contentModerator, mediaProcessor, eventPublisher, nil)
	a.presenceService = services.NewPresenceService(presenceRepo, a.wsManager, locationService, eventPublisher, nil)

	// Initialize messaging service
	a.messagingService = services.NewMessagingService(a.dialogService, a.messageService, nil, nil, nil)

	log.Println("Services initialized successfully with real implementations")
	return nil
}

// initHandlers initializes HTTP handlers
func (a *App) initHandlers() error {
	// Initialize handlers
	a.messagingHandlers = handlers.NewMessagingHandler(a.messagingService)
	a.wsHandler = handlers.NewWebSocketHandler(a.wsManager)

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
	v1 := router.PathPrefix("/v1").Subrouter()

	// Gateway-compatible health check endpoint
	v1.HandleFunc("/healthcheck", a.healthCheck).Methods("GET")

	// API v1 routes
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// Health endpoint for messaging
	apiV1.HandleFunc("/messaging/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "messaging service healthy"}`))
	}).Methods("GET")

	// Register messaging routes
	a.messagingHandlers.RegisterRoutes(apiV1)

	// Register WebSocket routes
	a.wsHandler.RegisterRoutes(apiV1)

	// Start messaging handler
	a.messagingHandlers.Start()

	a.router = router
	log.Println("Router initialized successfully with messaging and WebSocket routes")
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

	// Close ScyllaDB session
	if a.session != nil {
		a.session.Close()
		log.Println("ScyllaDB session closed")
	}

	log.Println("Messaging service shutdown completed")
	return nil
}

// testAuthMiddleware provides test authentication for development
func (a *App) testAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health endpoints
		if r.URL.Path == "/health" || r.URL.Path == "/ready" || r.URL.Path == "/v1/healthcheck" {
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

	// Check ScyllaDB connection
	if a.session != nil {
		if err := a.session.Query("SELECT now() FROM system.local").Exec(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			response := map[string]interface{}{
				"status":  "error",
				"message": "ScyllaDB not ready",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":   "ready",
		"service":  "messaging-service",
		"database": "scylladb-connected",
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
