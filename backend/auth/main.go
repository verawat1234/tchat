package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
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

// generateRandomOTP generates a random 6-digit OTP for development mode
func generateRandomOTP() string {
	source := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(source)
	otp := rnd.Intn(900000) + 100000 // Generates number between 100000-999999
	return fmt.Sprintf("%06d", otp)
}

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
	authHandlers    *handlers.AuthHandler
	// profileHandlers *handlers.ProfileHandler // Temporarily disabled
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
	// Pre-migration: Fix existing NULL values in phone_number and country_code
	if err := a.fixNullPhoneAndCountryCode(db); err != nil {
		log.Printf("Warning: Failed to fix NULL phone/country values: %v", err)
	}

	// Fix NULL profile values before migration attempts
	if err := a.fixNullProfileDisplayNames(db); err != nil {
		log.Printf("Warning: Failed to fix NULL profile_display_name values: %v", err)
	}

	// Handle user table migration manually due to NOT NULL constraints
	if err := a.migrateUserProfileColumns(db); err != nil {
		return fmt.Errorf("failed to migrate user profile columns: %w", err)
	}

	// Run auto migration for all models including User
	if err := db.AutoMigrate(
		&sharedModels.User{},
		&models.Session{},
		&models.KYC{},
		&sharedModels.Event{},
		&services.OTP{},
	); err != nil {
		return err
	}

	log.Printf("All database migrations completed successfully")

	return nil
}

// fixNullProfileDisplayNames fixes existing NULL profile_display_name values
func (a *App) fixNullProfileDisplayNames(db *gorm.DB) error {
	// Check if users table exists
	if !db.Migrator().HasTable("users") {
		return nil
	}

	// Check if profile_display_name column exists
	if !db.Migrator().HasColumn(&sharedModels.User{}, "profile_display_name") {
		log.Println("profile_display_name column doesn't exist yet, skipping NULL fixes")
		return nil
	}

	// Migrate data from display_name to profile_display_name if display_name column exists
	var displayNameColumnCount int64
	if err := db.Raw("SELECT count(*) FROM information_schema.columns WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'users' AND column_name = 'display_name'").Scan(&displayNameColumnCount).Error; err != nil {
		log.Printf("Warning: Could not check for display_name column: %v", err)
	} else if displayNameColumnCount > 0 {
		// Migrate display_name to profile_display_name
		result := db.Exec("UPDATE users SET profile_display_name = COALESCE(display_name, 'User') WHERE profile_display_name IS NULL OR profile_display_name = ''")
		if result.Error != nil {
			log.Printf("Warning: Failed to migrate display_name to profile_display_name: %v", result.Error)
		} else {
			log.Printf("Migrated %d users from display_name to profile_display_name", result.RowsAffected)
		}
	}

	// Update NULL profile_display_name values with a default name
	result := db.Exec("UPDATE users SET profile_display_name = 'User' WHERE profile_display_name IS NULL OR profile_display_name = ''")
	if result.Error != nil {
		return result.Error
	}

	log.Printf("Fixed %d users with NULL profile_display_name", result.RowsAffected)
	return nil
}

// fixNullPhoneAndCountryCode fixes existing NULL phone_number and country_code values
func (a *App) fixNullPhoneAndCountryCode(db *gorm.DB) error {
	// Check if users table exists
	if !db.Migrator().HasTable("users") {
		return nil
	}

	// Check if phone_number column exists before updating
	if db.Migrator().HasColumn(&sharedModels.User{}, "phone_number") {
		// Update NULL phone_number values with a default phone number
		result := db.Exec("UPDATE users SET phone_number = '+66000000000' WHERE phone_number IS NULL OR phone_number = ''")
		if result.Error != nil {
			return result.Error
		}
		log.Printf("Fixed %d users with NULL phone_number", result.RowsAffected)
	}

	// Check if country_code column exists before updating
	if db.Migrator().HasColumn(&sharedModels.User{}, "country_code") {
		// Update NULL country_code values with a default country code
		result := db.Exec("UPDATE users SET country_code = 'TH' WHERE country_code IS NULL OR country_code = ''")
		if result.Error != nil {
			return result.Error
		}
		log.Printf("Fixed %d users with NULL country_code", result.RowsAffected)
	}

	return nil
}

// migrateProfileColumns migrates data from old profile columns to new embedded profile columns
func (a *App) migrateProfileColumns(db *gorm.DB) error {
	// Check if users table exists
	if !db.Migrator().HasTable("users") {
		return nil
	}

	// Check if new profile columns exist after GORM migration
	if !db.Migrator().HasColumn(&sharedModels.User{}, "profile_display_name") {
		return nil // GORM migration hasn't run yet
	}

	// Migrate locale data (existing locale -> profile_locale)
	var localeColumnCount int64
	if err := db.Raw("SELECT count(*) FROM information_schema.columns WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'users' AND column_name = 'locale'").Scan(&localeColumnCount).Error; err == nil && localeColumnCount > 0 {
		result := db.Exec("UPDATE users SET profile_locale = COALESCE(locale, 'en') WHERE profile_locale IS NULL OR profile_locale = ''")
		if result.Error == nil {
			log.Printf("Migrated %d users from locale to profile_locale", result.RowsAffected)
		}
	}

	// Migrate timezone data (existing timezone -> profile_timezone)
	var timezoneColumnCount int64
	if err := db.Raw("SELECT count(*) FROM information_schema.columns WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'users' AND column_name = 'timezone'").Scan(&timezoneColumnCount).Error; err == nil && timezoneColumnCount > 0 {
		result := db.Exec("UPDATE users SET profile_timezone = COALESCE(timezone, 'UTC') WHERE profile_timezone IS NULL OR profile_timezone = ''")
		if result.Error == nil {
			log.Printf("Migrated %d users from timezone to profile_timezone", result.RowsAffected)
		}
	}

	// Migrate avatar data (existing avatar -> profile_avatar_url)
	var avatarColumnCount int64
	if err := db.Raw("SELECT count(*) FROM information_schema.columns WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'users' AND column_name = 'avatar'").Scan(&avatarColumnCount).Error; err == nil && avatarColumnCount > 0 {
		result := db.Exec("UPDATE users SET profile_avatar_url = avatar WHERE avatar IS NOT NULL AND avatar != '' AND (profile_avatar_url IS NULL OR profile_avatar_url = '')")
		if result.Error == nil {
			log.Printf("Migrated %d users from avatar to profile_avatar_url", result.RowsAffected)
		}
	}

	return nil
}

// migrateUserProfileColumns creates the necessary profile columns with proper data migration
func (a *App) migrateUserProfileColumns(db *gorm.DB) error {
	// Check if users table exists
	if !db.Migrator().HasTable("users") {
		log.Println("Users table doesn't exist, will be created by AutoMigrate")
		return nil
	}

	log.Println("Starting user profile column migration...")

	// Manually add profile columns as nullable first to avoid constraint violations
	profileColumns := []struct {
		name string
		dataType string
		oldColumn string
		defaultValue string
	}{
		{"profile_display_name", "varchar(100)", "display_name", "User"},
		{"profile_avatar_url", "varchar(500)", "avatar", ""},
		{"profile_locale", "varchar(5) DEFAULT 'en'", "locale", "en"},
		{"profile_timezone", "varchar(50) DEFAULT 'UTC'", "timezone", "UTC"},
	}

	for _, col := range profileColumns {
		// Check if column exists
		var columnCount int64
		db.Raw("SELECT count(*) FROM information_schema.columns WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'users' AND column_name = ?", col.name).Scan(&columnCount)

		if columnCount == 0 {
			// Add column as nullable first
			if err := db.Exec(fmt.Sprintf("ALTER TABLE users ADD COLUMN %s %s", col.name, col.dataType)).Error; err != nil {
				log.Printf("Warning: Failed to add column %s: %v", col.name, err)
				continue
			}
			log.Printf("Added column %s", col.name)

			// Migrate data from old column if it exists
			var oldColumnCount int64
			db.Raw("SELECT count(*) FROM information_schema.columns WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'users' AND column_name = ?", col.oldColumn).Scan(&oldColumnCount)

			if oldColumnCount > 0 {
				query := fmt.Sprintf("UPDATE users SET %s = COALESCE(%s, '%s') WHERE %s IS NULL", col.name, col.oldColumn, col.defaultValue, col.name)
				if result := db.Exec(query); result.Error == nil {
					log.Printf("Migrated %d users from %s to %s", result.RowsAffected, col.oldColumn, col.name)
				}
			} else {
				// Set default values for existing records
				query := fmt.Sprintf("UPDATE users SET %s = '%s' WHERE %s IS NULL", col.name, col.defaultValue, col.name)
				if result := db.Exec(query); result.Error == nil {
					log.Printf("Set default values for %d users in %s", result.RowsAffected, col.name)
				}
			}
		}
	}

	log.Println("User profile column migration completed")
	return nil
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
	a.jwtService = services.NewJWTService(a.config)
	a.sessionService = services.NewSessionService(sessionRepo, eventPublisher, a.jwtService, a.db)
	a.kycService = services.NewKYCService(kycRepo, a.userService, nil, nil, nil, eventPublisher, a.db)
	// Create mock dependencies for AuthService
	mockOTPRepo := &mockOTPRepository{}
	mockSMSProvider := &mockSMSProvider{}
	mockRateLimiter := &mockRateLimiter{}
	mockSecurityLogger := &mockSecurityLogger{}
	a.authService = services.NewAuthService(a.userService, a.sessionService, mockOTPRepo, mockSMSProvider, mockRateLimiter, mockSecurityLogger, eventPublisher, a.db)

	log.Println("Services initialized successfully")
	return nil
}

// initHandlers initializes HTTP handlers
func (a *App) initHandlers() error {
	// Initialize handlers
	a.authHandlers = handlers.NewAuthHandler(
		a.authService,
		a.sessionService,
		a.userService,
	)

	// a.profileHandlers = handlers.NewProfileHandler(a.userService, a.kycService) // Temporarily disabled

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

	// Rate limiting - temporarily disabled
	// if a.config.RateLimit.Enabled {
	// 	rateLimiter := middleware.NewRateLimiter(
	// 		a.config.RateLimit.RequestsPerMinute,
	// 		a.config.RateLimit.BurstSize,
	// 		a.config.RateLimit.CleanupInterval,
	// 	)
	// 	router.Use(rateLimiter.Middleware())
	// }

	// Health check endpoint
	router.GET("/health", a.healthCheck)
	router.GET("/ready", a.readinessCheck)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		handlers.RegisterAuthRoutes(v1, a.authHandlers, middleware.NewAuthMiddleware(a.config).RequireAuth())

		// Profile routes - temporarily disabled
		// handlers.RegisterProfileRoutes(v1, a.userService, a.kycService)
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
	responses.SendSuccessResponse(c, gin.H{
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
			responses.SendErrorResponse(c, http.StatusServiceUnavailable, "database_not_ready", "Database connection failed")
			return
		}
	}

	responses.SendSuccessResponse(c, gin.H{
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

func (m *mockUserRepository) Create(ctx context.Context, user *sharedModels.User) error {
	// User is already a shared GORM model, use directly
	return m.db.WithContext(ctx).Create(user).Error
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*sharedModels.User, error) {
	var user sharedModels.User
	err := m.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}

	// Return shared model directly
	return &user, nil
}

func (m *mockUserRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*sharedModels.User, error) {
	var user sharedModels.User
	err := m.db.WithContext(ctx).Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		return nil, err
	}

	// Return shared model directly
	return &user, nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*sharedModels.User, error) {
	var user sharedModels.User
	err := m.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	// Return shared model directly
	return &user, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *sharedModels.User) error {
	// User is already a shared GORM model, use directly
	return m.db.WithContext(ctx).Save(user).Error
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Delete(&sharedModels.User{}, id).Error
}

func (m *mockUserRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*sharedModels.User, error) {
	var users []*sharedModels.User
	err := m.db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error
	return users, err
}

func (m *mockUserRepository) GetUsers(ctx context.Context, filters services.UserFilters, pagination services.Pagination) ([]*sharedModels.User, error) {
	var users []*sharedModels.User
	err := m.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (m *mockUserRepository) SearchByUsername(ctx context.Context, username string, limit int) ([]*sharedModels.User, error) {
	var users []*sharedModels.User
	err := m.db.WithContext(ctx).Where("username ILIKE ?", "%"+username+"%").Limit(limit).Find(&users).Error
	return users, err
}

func (m *mockUserRepository) GetUserStats(ctx context.Context, userID uuid.UUID) (*services.UserStats, error) {
	return &services.UserStats{
		// Add actual fields based on services.UserStats struct
	}, nil
}

func (m *mockUserRepository) List(ctx context.Context, filters services.UserFilters, pagination services.Pagination) ([]*sharedModels.User, int64, error) {
	// Simple mock implementation
	var users []*sharedModels.User
	err := m.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	count := int64(len(users))
	return users, count, nil
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
	// Parse JWT to extract session ID (access tokens are not stored in DB)
	// For now, return error since we should validate tokens through JWT parsing
	// This method should not be used for access token validation
	return nil, errors.New("access tokens are not stored in database - use JWT validation instead")
}

func (m *mockSessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	var session models.Session
	err := m.db.WithContext(ctx).Where("refresh_token_hash = ?", refreshToken).First(&session).Error
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

// Mock implementations for missing dependencies
type mockOTPRepository struct{}

func (m *mockOTPRepository) Store(ctx context.Context, phoneNumber string, otp string, expiresAt time.Time) error {
	log.Printf("Mock OTP stored for %s", phoneNumber)
	return nil
}

func (m *mockOTPRepository) Verify(ctx context.Context, phoneNumber string, otp string) (bool, error) {
	// Simple mock: accept "123456" as valid OTP
	return otp == "123456", nil
}

func (m *mockOTPRepository) Delete(ctx context.Context, phoneNumber string) error {
	log.Printf("Mock OTP deleted for %s", phoneNumber)
	return nil
}

func (m *mockOTPRepository) GetAttemptCount(ctx context.Context, phoneNumber string, timeWindow time.Duration) (int, error) {
	return 1, nil
}

func (m *mockOTPRepository) IncrementAttempts(ctx context.Context, phoneNumber string) error {
	return nil
}

func (m *mockOTPRepository) Create(ctx context.Context, otp *services.OTP) error {
	// Generate real OTP for dev mode
	if otp.Code == "" || otp.Code == "123456" {
		otp.Code = generateRandomOTP()
	}
	log.Printf("DEV MODE: Mock OTP created for %s with code %s", otp.PhoneNumber, otp.Code)
	return nil
}

func (m *mockOTPRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*services.OTP, error) {
	// Mock implementation with real OTP generation for dev mode
	realOTP := generateRandomOTP()
	log.Printf("DEV MODE: Generated real OTP %s for %s", realOTP, phoneNumber)

	otp := &services.OTP{
		ID:          uuid.New(),
		PhoneNumber: phoneNumber,
		Code:        realOTP,
		ExpiresAt:   time.Now().Add(5 * time.Minute),
		Status:      services.OTPStatusPending,
		Type:        services.OTPTypeLogin,
		CreatedAt:   time.Now(),
	}
	return otp, nil
}

func (m *mockOTPRepository) GetByID(ctx context.Context, id uuid.UUID) (*services.OTP, error) {
	// Mock implementation with real OTP generation for dev mode
	realOTP := generateRandomOTP()
	log.Printf("DEV MODE: Generated real OTP %s for ID %s", realOTP, id.String())

	otp := &services.OTP{
		ID:          id,
		PhoneNumber: "+66812345678",
		Code:        realOTP,
		ExpiresAt:   time.Now().Add(5 * time.Minute),
		Status:      services.OTPStatusPending,
		Type:        services.OTPTypeLogin,
		CreatedAt:   time.Now(),
	}
	return otp, nil
}

func (m *mockOTPRepository) Update(ctx context.Context, otp *services.OTP) error {
	log.Printf("Mock OTP updated for %s", otp.PhoneNumber)
	return nil
}


func (m *mockOTPRepository) DeleteExpired(ctx context.Context) error {
	log.Printf("Mock delete expired OTPs")
	return nil
}


type mockSMSProvider struct{}

func (m *mockSMSProvider) SendSMS(ctx context.Context, phoneNumber string, message string) error {
	log.Printf("Mock SMS sent to %s: %s", phoneNumber, message)
	return nil
}

func (m *mockSMSProvider) GetRemainingCredits(ctx context.Context) (int, error) {
	return 1000, nil
}

func (m *mockSMSProvider) SendOTP(ctx context.Context, phoneNumber, otp, template string) error {
	message := fmt.Sprintf("Your OTP code is: %s", otp)
	return m.SendSMS(ctx, phoneNumber, message)
}

type mockRateLimiter struct{}

func (m *mockRateLimiter) Allow(ctx context.Context, key string, limit int, duration time.Duration) (bool, error) {
	return true, nil
}

func (m *mockRateLimiter) Reset(ctx context.Context, key string) error {
	return nil
}

func (m *mockRateLimiter) GetRemainingAttempts(ctx context.Context, key string, limit int, window time.Duration) (int, error) {
	return limit - 1, nil
}


type mockSecurityLogger struct{}

func (m *mockSecurityLogger) LogSecurityEvent(ctx context.Context, event string, details map[string]interface{}) {
	log.Printf("Mock security event: %s %+v", event, details)
}

func (m *mockSecurityLogger) LogLoginAttempt(ctx context.Context, phoneNumber, userAgent, ipAddress string, success bool, reason string) {
	log.Printf("Mock login attempt: %s from %s (%s) - Success: %t, Reason: %s", phoneNumber, ipAddress, userAgent, success, reason)
}

func (m *mockSecurityLogger) LogOTPGeneration(ctx context.Context, phoneNumber, ipAddress string) {
	log.Printf("Mock OTP generation: %s from %s", phoneNumber, ipAddress)
}

func (m *mockSecurityLogger) LogSuspiciousActivity(ctx context.Context, userID uuid.UUID, activity, reason string) {
	log.Printf("Mock suspicious activity: User %s - %s (%s)", userID, activity, reason)
}

