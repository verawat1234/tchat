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

	"tchat.dev/payment/handlers"
	"tchat.dev/payment/models"
	"tchat.dev/payment/services"
	"tchat.dev/shared/config"
	"tchat.dev/shared/database"
	"tchat.dev/shared/middleware"
	"tchat.dev/shared/responses"
	sharedModels "tchat.dev/shared/models"
)

// App represents the main payment application
type App struct {
	config    *config.Config
	db        *gorm.DB
	router    *gin.Engine
	server    *http.Server
	validator *validator.Validate

	// Services
	paymentService        *services.PaymentService
	walletService         *services.WalletService
	transactionService    *services.TransactionService
	paymentGatewayService *services.PaymentGatewayService

	// Handlers
	paymentHandlers *handlers.PaymentHandler
	walletHandlers  *handlers.WalletHandler
}

// NewApp creates a new payment application instance
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

	log.Printf("Payment service initialized successfully on port %d", a.config.Server.Port)
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

// runMigrations runs database migrations for payment service models
func (a *App) runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Wallet{},
		&models.Transaction{},
		&models.PaymentMethod{},
		&sharedModels.Event{},
	)
}

// initServices initializes all business logic services
func (a *App) initServices() error {
	// Mock event publisher for now
	eventPublisher := &mockEventPublisher{}

	// Mock repositories for now - these would be implemented based on your repository pattern
	walletRepo := &mockWalletRepository{db: a.db}
	transactionRepo := &mockTransactionRepository{db: a.db}
	paymentMethodRepo := &mockPaymentMethodRepository{db: a.db}

	// Mock external services
	complianceChecker := &mockComplianceChecker{}
	exchangeRateService := &mockExchangeRateService{}
	paymentProcessor := &mockPaymentProcessor{}
	cacheService := &mockCacheService{}

	// Initialize services
	a.walletService = services.NewWalletService(
		walletRepo,
		transactionRepo,
		complianceChecker,
		exchangeRateService,
		eventPublisher,
		a.db,
	)

	a.transactionService = services.NewTransactionService(
		transactionRepo,
		walletRepo,
		complianceChecker,
		eventPublisher,
		a.db,
	)

	a.paymentGatewayService = services.NewPaymentGatewayService(
		paymentProcessor,
		transactionRepo,
		eventPublisher,
	)

	a.paymentService = services.NewPaymentService(
		walletRepo,
		transactionRepo,
		cacheService,
		eventPublisher,
		paymentProcessor,
		a.db,
	)

	log.Println("Services initialized successfully")
	return nil
}

// initHandlers initializes HTTP handlers
func (a *App) initHandlers() error {
	// Initialize handlers
	a.paymentHandlers = handlers.NewPaymentHandler(
		a.paymentService,
		a.walletService,
		a.transactionService,
	)

	a.walletHandlers = handlers.NewWalletHandler(
		a.walletService,
		a.transactionService,
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

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Payment routes
		handlers.RegisterPaymentRoutes(v1, a.paymentService, a.walletService, a.transactionService)

		// Wallet routes
		handlers.RegisterWalletRoutes(v1, a.walletService, a.transactionService)

		// Transaction routes
		handlers.RegisterTransactionRoutes(v1, a.transactionService)
	}

	// Webhook endpoints for payment gateways
	webhooks := router.Group("/webhooks")
	{
		handlers.RegisterWebhookRoutes(webhooks, a.paymentGatewayService)
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
	log.Printf("Starting payment service on %s", a.server.Addr)

	// Start server in a goroutine
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Payment service is running on %s", a.server.Addr)
	return nil
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down payment service...")

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

	log.Println("Payment service shutdown completed")
	return nil
}

// Health check endpoint
func (a *App) healthCheck(c *gin.Context) {
	responses.SuccessResponse(c, gin.H{
		"status":    "ok",
		"service":   "payment-service",
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
		"service":  "payment-service",
		"database": "connected",
	})
}

// Mock implementations for repositories and services

type mockEventPublisher struct{}

func (m *mockEventPublisher) Publish(ctx context.Context, event *sharedModels.Event) error {
	log.Printf("Event published: %s - %s", event.Type, event.Subject)
	return nil
}

type mockWalletRepository struct {
	db *gorm.DB
}

func (m *mockWalletRepository) Create(ctx context.Context, wallet *models.Wallet) error {
	return m.db.WithContext(ctx).Create(wallet).Error
}

func (m *mockWalletRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	err := m.db.WithContext(ctx).First(&wallet, id).Error
	return &wallet, err
}

func (m *mockWalletRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Wallet, error) {
	var wallets []*models.Wallet
	err := m.db.WithContext(ctx).Where("user_id = ?", userID).Find(&wallets).Error
	return wallets, err
}

func (m *mockWalletRepository) GetByUserIDAndCurrency(ctx context.Context, userID uuid.UUID, currency models.Currency) (*models.Wallet, error) {
	var wallet models.Wallet
	err := m.db.WithContext(ctx).Where("user_id = ? AND currency = ?", userID, currency).First(&wallet).Error
	return &wallet, err
}

func (m *mockWalletRepository) Update(ctx context.Context, wallet *models.Wallet) error {
	return m.db.WithContext(ctx).Save(wallet).Error
}

func (m *mockWalletRepository) UpdateBalance(ctx context.Context, walletID uuid.UUID, amount float64, transactionID uuid.UUID) error {
	return m.db.WithContext(ctx).Model(&models.Wallet{}).Where("id = ?", walletID).Update("balance", amount).Error
}

func (m *mockWalletRepository) FreezeAmount(ctx context.Context, walletID uuid.UUID, amount float64, reason string) error {
	return m.db.WithContext(ctx).Model(&models.Wallet{}).Where("id = ?", walletID).Updates(map[string]interface{}{
		"frozen_amount": gorm.Expr("frozen_amount + ?", amount),
		"balance":       gorm.Expr("balance - ?", amount),
	}).Error
}

func (m *mockWalletRepository) UnfreezeAmount(ctx context.Context, walletID uuid.UUID, amount float64) error {
	return m.db.WithContext(ctx).Model(&models.Wallet{}).Where("id = ?", walletID).Updates(map[string]interface{}{
		"frozen_amount": gorm.Expr("frozen_amount - ?", amount),
		"balance":       gorm.Expr("balance + ?", amount),
	}).Error
}

func (m *mockWalletRepository) GetWalletHistory(ctx context.Context, walletID uuid.UUID, limit int) ([]*services.WalletBalanceHistory, error) {
	// Mock implementation
	return []*services.WalletBalanceHistory{}, nil
}

func (m *mockWalletRepository) GetWalletStats(ctx context.Context, userID uuid.UUID) (*services.WalletStats, error) {
	return &services.WalletStats{
		TotalBalance:       1000.50,
		TotalTransactions:  45,
		ThisMonthSpending:  250.75,
		ThisMonthIncome:    300.00,
		LargestTransaction: 150.00,
		AverageTransaction: 22.23,
	}, nil
}

type mockTransactionRepository struct {
	db *gorm.DB
}

func (m *mockTransactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	return m.db.WithContext(ctx).Create(transaction).Error
}

func (m *mockTransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	err := m.db.WithContext(ctx).First(&transaction, id).Error
	return &transaction, err
}

func (m *mockTransactionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	err := m.db.WithContext(ctx).Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&transactions).Error
	return transactions, err
}

func (m *mockTransactionRepository) Update(ctx context.Context, transaction *models.Transaction) error {
	return m.db.WithContext(ctx).Save(transaction).Error
}

func (m *mockTransactionRepository) GetByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	err := m.db.WithContext(ctx).Where("from_wallet_id = ? OR to_wallet_id = ?", walletID, walletID).Limit(limit).Offset(offset).Find(&transactions).Error
	return transactions, err
}

type mockPaymentMethodRepository struct {
	db *gorm.DB
}

func (m *mockPaymentMethodRepository) Create(ctx context.Context, method *models.PaymentMethod) error {
	return m.db.WithContext(ctx).Create(method).Error
}

func (m *mockPaymentMethodRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.PaymentMethod, error) {
	var method models.PaymentMethod
	err := m.db.WithContext(ctx).First(&method, id).Error
	return &method, err
}

func (m *mockPaymentMethodRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.PaymentMethod, error) {
	var methods []*models.PaymentMethod
	err := m.db.WithContext(ctx).Where("user_id = ?", userID).Find(&methods).Error
	return methods, err
}

func (m *mockPaymentMethodRepository) Update(ctx context.Context, method *models.PaymentMethod) error {
	return m.db.WithContext(ctx).Save(method).Error
}

func (m *mockPaymentMethodRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Delete(&models.PaymentMethod{}, id).Error
}

type mockComplianceChecker struct{}

func (m *mockComplianceChecker) CheckTransactionLimits(ctx context.Context, userID uuid.UUID, amount float64, currency models.Currency, transactionType string) error {
	if amount > 10000 {
		return fmt.Errorf("transaction amount exceeds daily limit")
	}
	return nil
}

func (m *mockComplianceChecker) CheckAMLCompliance(ctx context.Context, userID uuid.UUID, amount float64, currency models.Currency) error {
	if amount > 50000 {
		return fmt.Errorf("transaction requires AML verification")
	}
	return nil
}

func (m *mockComplianceChecker) ReportSuspiciousActivity(ctx context.Context, userID uuid.UUID, activity string, metadata map[string]interface{}) error {
	log.Printf("Suspicious activity reported for user %s: %s", userID, activity)
	return nil
}

type mockExchangeRateService struct{}

func (m *mockExchangeRateService) GetExchangeRate(ctx context.Context, fromCurrency, toCurrency models.Currency) (float64, error) {
	// Mock exchange rates
	return 1.0, nil
}

func (m *mockExchangeRateService) ConvertAmount(ctx context.Context, amount float64, fromCurrency, toCurrency models.Currency) (float64, error) {
	rate, err := m.GetExchangeRate(ctx, fromCurrency, toCurrency)
	if err != nil {
		return 0, err
	}
	return amount * rate, nil
}

type mockPaymentProcessor struct{}

func (m *mockPaymentProcessor) ProcessPayment(ctx context.Context, req *services.PaymentRequest) (*services.PaymentResponse, error) {
	return &services.PaymentResponse{
		TransactionID: uuid.New(),
		Status:        "completed",
		Message:       "Payment processed successfully",
	}, nil
}

func (m *mockPaymentProcessor) RefundPayment(ctx context.Context, transactionID uuid.UUID, amount float64) (*services.RefundResponse, error) {
	return &services.RefundResponse{
		RefundID: uuid.New(),
		Status:   "completed",
		Message:  "Refund processed successfully",
	}, nil
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

	log.Println("Payment service stopped")
}