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

	"tchat.dev/commerce/handlers"
	"tchat.dev/commerce/models"
	"tchat.dev/commerce/services"
	"tchat.dev/shared/config"
	"tchat.dev/shared/database"
	"tchat.dev/shared/middleware"
	"tchat.dev/shared/responses"
	sharedModels "tchat.dev/shared/models"
)

// App represents the main commerce application
type App struct {
	config    *config.Config
	db        *gorm.DB
	router    *gin.Engine
	server    *http.Server
	validator *validator.Validate

	// Services
	commerceService *services.CommerceService

	// Handlers
	commerceHandlers *handlers.CommerceHandler
}

// NewApp creates a new commerce application instance
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

	log.Printf("Commerce service initialized successfully on port %d", a.config.Server.Port)
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

// runMigrations runs database migrations for commerce service models
func (a *App) runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Shop{},
		&models.Product{},
		&models.Order{},
		&sharedModels.Event{},
	)
}

// initServices initializes all business logic services
func (a *App) initServices() error {
	// Mock repositories and services
	shopRepo := &mockShopRepository{db: a.db}
	productRepo := &mockProductRepository{db: a.db}
	orderRepo := &mockOrderRepository{db: a.db}
	cacheService := &mockCacheService{}
	eventService := &mockEventService{}
	paymentService := &mockPaymentService{}

	// Initialize commerce service
	a.commerceService = services.NewCommerceService(
		shopRepo,
		productRepo,
		orderRepo,
		cacheService,
		eventService,
		paymentService,
		services.DefaultCommerceConfig(),
	)

	log.Println("Services initialized successfully")
	return nil
}

// initHandlers initializes HTTP handlers
func (a *App) initHandlers() error {
	// Initialize handlers
	a.commerceHandlers = handlers.NewCommerceHandler(a.commerceService)

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
		// Commerce routes
		handlers.RegisterCommerceRoutes(v1, a.commerceService)
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
	log.Printf("Starting commerce service on %s", a.server.Addr)

	// Start server in a goroutine
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Commerce service is running on %s", a.server.Addr)
	return nil
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down commerce service...")

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

	log.Println("Commerce service shutdown completed")
	return nil
}

// Health check endpoint
func (a *App) healthCheck(c *gin.Context) {
	responses.SuccessResponse(c, gin.H{
		"status":    "ok",
		"service":   "commerce",
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
		"service":  "commerce",
		"database": "connected",
	})
}

// Mock implementations

type mockShopRepository struct {
	db *gorm.DB
}

func (m *mockShopRepository) Create(ctx context.Context, shop *models.Shop) error {
	return m.db.WithContext(ctx).Create(shop).Error
}

func (m *mockShopRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Shop, error) {
	var shop models.Shop
	err := m.db.WithContext(ctx).First(&shop, id).Error
	return &shop, err
}

func (m *mockShopRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Shop, error) {
	var shops []*models.Shop
	err := m.db.WithContext(ctx).Where("owner_id = ?", userID).Find(&shops).Error
	return shops, err
}

func (m *mockShopRepository) Update(ctx context.Context, shop *models.Shop) error {
	return m.db.WithContext(ctx).Save(shop).Error
}

func (m *mockShopRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Delete(&models.Shop{}, id).Error
}

func (m *mockShopRepository) Search(ctx context.Context, query string, limit, offset int) ([]*models.Shop, error) {
	var shops []*models.Shop
	err := m.db.WithContext(ctx).Where("name ILIKE ?", "%"+query+"%").Limit(limit).Offset(offset).Find(&shops).Error
	return shops, err
}

type mockProductRepository struct {
	db *gorm.DB
}

func (m *mockProductRepository) Create(ctx context.Context, product *models.Product) error {
	return m.db.WithContext(ctx).Create(product).Error
}

func (m *mockProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := m.db.WithContext(ctx).First(&product, id).Error
	return &product, err
}

func (m *mockProductRepository) GetByShopID(ctx context.Context, shopID uuid.UUID, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	err := m.db.WithContext(ctx).Where("shop_id = ?", shopID).Limit(limit).Offset(offset).Find(&products).Error
	return products, err
}

func (m *mockProductRepository) Update(ctx context.Context, product *models.Product) error {
	return m.db.WithContext(ctx).Save(product).Error
}

func (m *mockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

func (m *mockProductRepository) Search(ctx context.Context, query string, filters map[string]interface{}, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	err := m.db.WithContext(ctx).Where("name ILIKE ?", "%"+query+"%").Limit(limit).Offset(offset).Find(&products).Error
	return products, err
}

func (m *mockProductRepository) GetByCategory(ctx context.Context, category string, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	err := m.db.WithContext(ctx).Where("category = ?", category).Limit(limit).Offset(offset).Find(&products).Error
	return products, err
}

func (m *mockProductRepository) UpdateInventory(ctx context.Context, productID uuid.UUID, quantity int) error {
	return m.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", productID).Update("inventory_count", quantity).Error
}

type mockOrderRepository struct {
	db *gorm.DB
}

func (m *mockOrderRepository) Create(ctx context.Context, order *models.Order) error {
	return m.db.WithContext(ctx).Create(order).Error
}

func (m *mockOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := m.db.WithContext(ctx).First(&order, id).Error
	return &order, err
}

func (m *mockOrderRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	err := m.db.WithContext(ctx).Where("buyer_id = ?", userID).Limit(limit).Offset(offset).Find(&orders).Error
	return orders, err
}

func (m *mockOrderRepository) GetByShopID(ctx context.Context, shopID uuid.UUID, limit, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	err := m.db.WithContext(ctx).Where("shop_id = ?", shopID).Limit(limit).Offset(offset).Find(&orders).Error
	return orders, err
}

func (m *mockOrderRepository) Update(ctx context.Context, order *models.Order) error {
	return m.db.WithContext(ctx).Save(order).Error
}

func (m *mockOrderRepository) GetByStatus(ctx context.Context, status models.OrderStatus, limit, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	err := m.db.WithContext(ctx).Where("status = ?", status).Limit(limit).Offset(offset).Find(&orders).Error
	return orders, err
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

type mockPaymentService struct{}

func (m *mockPaymentService) ProcessPayment(ctx context.Context, orderID uuid.UUID, amount float64, paymentMethodID uuid.UUID) error {
	log.Printf("Payment processed for order %s: $%.2f", orderID, amount)
	return nil
}

func (m *mockPaymentService) RefundPayment(ctx context.Context, orderID uuid.UUID, amount float64) error {
	log.Printf("Refund processed for order %s: $%.2f", orderID, amount)
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

	log.Println("Commerce service stopped")
}