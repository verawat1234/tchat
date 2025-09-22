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

	"tchat.dev/auth/handlers"
	"tchat.dev/auth/models"
	"tchat.dev/auth/services"
	"tchat.dev/shared/config"
	"tchat.dev/shared/database"
	"tchat.dev/shared/middleware"
	"tchat.dev/shared/responses"
	sharedModels "tchat.dev/shared/models"
)

// App represents the main application
type App struct {
	config     *config.Config
	db         *gorm.DB
	router     *gin.Engine
	server     *http.Server
	validator  *validator.Validate

	// Services
	authService    *services.AuthService
	userService    *services.UserService
	sessionService *services.SessionService
	jwtService     *services.JWTService
	kycService     *services.KYCService

	// Handlers
	authHandlers    *handlers.AuthHandlers
	profileHandlers *handlers.ProfileHandler
}

// NewApp creates a new application instance
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

	log.Printf("Auth service initialized successfully on port %d", a.config.Server.Port)
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

// runMigrations runs database migrations for auth service models
func (a *App) runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.KYC{},
		&sharedModels.Event{},
	)
}

// initServices initializes all business logic services
func (a *App) initServices() error {
	// Mock event publisher for now
	eventPublisher := &mockEventPublisher{}

	// Mock repositories for now - these would be implemented based on your repository pattern
	userRepo := &mockUserRepository{db: a.db}
	sessionRepo := &mockSessionRepository{db: a.db}
	kycRepo := &mockKYCRepository{db: a.db}

	// Initialize services
	a.userService = services.NewUserService(userRepo, eventPublisher, a.db)
	a.sessionService = services.NewSessionService(sessionRepo, eventPublisher, a.db)
	a.jwtService = services.NewJWTService(a.config)
	a.kycService = services.NewKYCService(kycRepo, a.userService, nil, nil, nil, eventPublisher, a.db)
	a.authService = services.NewAuthService(userRepo, a.sessionService, a.jwtService, eventPublisher, a.db)

	log.Println("Services initialized successfully")
	return nil
}

// initHandlers initializes HTTP handlers
func (a *App) initHandlers() error {
	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(a.jwtService, a.sessionService, a.userService)

	// Initialize handlers
	a.authHandlers = handlers.NewAuthHandlers(
		a.authService,
		a.userService,
		a.sessionService,
		a.jwtService,
	)

	a.profileHandlers = handlers.NewProfileHandler(a.userService, a.kycService)

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

	// Health check endpoint
	router.GET("/health", a.healthCheck)
	router.GET("/ready", a.readinessCheck)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		handlers.RegisterAuthRoutes(v1, a.authHandlers, middleware.NewAuthMiddleware(a.jwtService, a.sessionService, a.userService).RequireAuth())

		// Profile routes
		handlers.RegisterProfileRoutes(v1, a.userService, a.kycService)
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
	log.Printf("Starting auth service on %s", a.server.Addr)

	// Start server in a goroutine
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Auth service is running on %s", a.server.Addr)
	return nil
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down auth service...")

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

	log.Println("Auth service shutdown completed")
	return nil
}

// Health check endpoint
func (a *App) healthCheck(c *gin.Context) {
	responses.SuccessResponse(c, gin.H{
		"status":    "ok",
		"service":   "auth",
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
		"service":  "auth",
		"database": "connected",
	})
}

// Mock implementations - these would be replaced with actual repository implementations
type mockEventPublisher struct{}

func (m *mockEventPublisher) Publish(ctx context.Context, event *sharedModels.Event) error {
	log.Printf("Event published: %s - %s", event.Type, event.Subject)
	return nil
}

type mockUserRepository struct {
	db *gorm.DB
}

func (m *mockUserRepository) Create(ctx context.Context, user *models.User) error {
	return m.db.WithContext(ctx).Create(user).Error
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := m.db.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func (m *mockUserRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	var user models.User
	err := m.db.WithContext(ctx).Where("phone_number = ?", phoneNumber).First(&user).Error
	return &user, err
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := m.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (m *mockUserRepository) Update(ctx context.Context, user *models.User) error {
	return m.db.WithContext(ctx).Save(user).Error
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (m *mockUserRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.User, error) {
	var users []*models.User
	err := m.db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error
	return users, err
}

func (m *mockUserRepository) GetUsers(ctx context.Context, filters services.UserFilters, pagination services.Pagination) ([]*models.User, error) {
	var users []*models.User
	err := m.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (m *mockUserRepository) GetUserStats(ctx context.Context, userID uuid.UUID) (*services.UserStats, error) {
	return &services.UserStats{
		TotalDialogs:     25,
		TotalMessages:    1540,
		TotalContacts:    89,
		WalletBalance:    1250.75,
		TransactionCount: 47,
		JoinedDaysAgo:    156,
		LastActiveHours:  2,
	}, nil
}

type mockSessionRepository struct {
	db *gorm.DB
}

func (m *mockSessionRepository) Create(ctx context.Context, session *models.Session) error {
	return m.db.WithContext(ctx).Create(session).Error
}

func (m *mockSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := m.db.WithContext(ctx).First(&session, id).Error
	return &session, err
}

func (m *mockSessionRepository) GetByAccessToken(ctx context.Context, accessToken string) (*models.Session, error) {
	var session models.Session
	err := m.db.WithContext(ctx).Where("access_token = ?", accessToken).First(&session).Error
	return &session, err
}

func (m *mockSessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	var session models.Session
	err := m.db.WithContext(ctx).Where("refresh_token = ?", refreshToken).First(&session).Error
	return &session, err
}

func (m *mockSessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Session, error) {
	var sessions []*models.Session
	err := m.db.WithContext(ctx).Where("user_id = ?", userID).Find(&sessions).Error
	return sessions, err
}

func (m *mockSessionRepository) Update(ctx context.Context, session *models.Session) error {
	return m.db.WithContext(ctx).Save(session).Error
}

func (m *mockSessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Delete(&models.Session{}, id).Error
}

func (m *mockSessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return m.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.Session{}).Error
}

func (m *mockSessionRepository) CleanupExpiredSessions(ctx context.Context) error {
	return m.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&models.Session{}).Error
}

type mockKYCRepository struct {
	db *gorm.DB
}

func (m *mockKYCRepository) Create(ctx context.Context, kyc *models.KYC) error {
	return m.db.WithContext(ctx).Create(kyc).Error
}

func (m *mockKYCRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.KYC, error) {
	var kyc models.KYC
	err := m.db.WithContext(ctx).First(&kyc, id).Error
	return &kyc, err
}

func (m *mockKYCRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.KYC, error) {
	var kyc models.KYC
	err := m.db.WithContext(ctx).Where("user_id = ?", userID).First(&kyc).Error
	return &kyc, err
}

func (m *mockKYCRepository) Update(ctx context.Context, kyc *models.KYC) error {
	return m.db.WithContext(ctx).Save(kyc).Error
}

func (m *mockKYCRepository) GetByStatus(ctx context.Context, status models.KYCStatus) ([]*models.KYC, error) {
	var kycs []*models.KYC
	err := m.db.WithContext(ctx).Where("status = ?", status).Find(&kycs).Error
	return kycs, err
}

func (m *mockKYCRepository) GetPendingReviews(ctx context.Context) ([]*models.KYC, error) {
	var kycs []*models.KYC
	err := m.db.WithContext(ctx).Where("status = ?", models.KYCStatusPending).Find(&kycs).Error
	return kycs, err
}

func (m *mockKYCRepository) GetStatistics(ctx context.Context) (*services.KYCStatistics, error) {
	return &services.KYCStatistics{
		TotalSubmissions:      150,
		PendingReviews:        25,
		ApprovedToday:         12,
		RejectedToday:         3,
		AverageProcessingTime: 24 * time.Hour,
		ApprovalRate:          0.85,
	}, nil
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

	log.Println("Auth service stopped")
}