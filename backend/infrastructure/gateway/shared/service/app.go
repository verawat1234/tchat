package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"tchat.dev/shared/config"
)

// App is the generic application framework
type App struct {
	// Core components
	config    *config.Config
	db        *gorm.DB
	router    *gin.Engine
	server    *http.Server
	validator *validator.Validate

	// Service metadata
	serviceInfo ServiceInfo

	// Component registry
	registry ServiceRegistry

	// Initializers
	dbInitializer      DatabaseInitializer
	repoInitializer    RepositoryInitializer
	serviceInitializer ServiceInitializer
	handlerInitializer HandlerInitializer
	routeRegistrar     RouteRegistrar
	healthChecker      HealthChecker
	middlewareProvider MiddlewareProvider

	// Managers
	dbManager     DatabaseManager
	serverManager ServerManager

	// Data storage for components
	repositories map[string]interface{}
	services     map[string]interface{}
	handlers     map[string]interface{}
}

// AppConfig holds configuration for app creation
type AppConfig struct {
	ServiceInfo        ServiceInfo
	DatabaseInitializer DatabaseInitializer
	RepositoryInitializer RepositoryInitializer
	ServiceInitializer ServiceInitializer
	HandlerInitializer HandlerInitializer
	RouteRegistrar     RouteRegistrar
	HealthChecker      HealthChecker
	MiddlewareProvider MiddlewareProvider
}

// NewApp creates a new application instance with the provided configuration
func NewApp(cfg *config.Config, appCfg AppConfig) *App {
	return &App{
		config:      cfg,
		validator:   validator.New(),
		serviceInfo: appCfg.ServiceInfo,
		registry:    NewDefaultServiceRegistry(),

		// Assign initializers
		dbInitializer:      appCfg.DatabaseInitializer,
		repoInitializer:    appCfg.RepositoryInitializer,
		serviceInitializer: appCfg.ServiceInitializer,
		handlerInitializer: appCfg.HandlerInitializer,
		routeRegistrar:     appCfg.RouteRegistrar,
		healthChecker:      appCfg.HealthChecker,
		middlewareProvider: appCfg.MiddlewareProvider,

		// Initialize managers
		dbManager:     NewDefaultDatabaseManager(),
		serverManager: NewDefaultServerManager(),

		// Initialize maps
		repositories: make(map[string]interface{}),
		services:     make(map[string]interface{}),
		handlers:     make(map[string]interface{}),
	}
}

// Initialize initializes all application components in the correct order
func (a *App) Initialize() error {
	log.Printf("Initializing %s service...", a.serviceInfo.Name)

	// Initialize database
	if err := a.initDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize repositories
	if err := a.initRepositories(); err != nil {
		return fmt.Errorf("failed to initialize repositories: %w", err)
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

	// Start service registry components
	if err := a.registry.StartAll(context.Background()); err != nil {
		return fmt.Errorf("failed to start service components: %w", err)
	}

	log.Printf("%s service initialized successfully on port %d", a.serviceInfo.Name, a.serviceInfo.Port)
	return nil
}

// Run starts the application server
func (a *App) Run() error {
	log.Printf("Starting %s service on %s", a.serviceInfo.Name, a.server.Addr)

	// Start server in a goroutine
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("%s service is running on %s", a.serviceInfo.Name, a.server.Addr)
	return nil
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	log.Printf("Shutting down %s service...", a.serviceInfo.Name)

	// Stop service registry components
	if err := a.registry.StopAll(ctx); err != nil {
		log.Printf("Error stopping service components: %v", err)
	}

	// Shutdown HTTP server
	if err := a.serverManager.Shutdown(ctx, a.server); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	// Close database connection
	if err := a.dbManager.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	log.Printf("%s service shutdown completed", a.serviceInfo.Name)
	return nil
}

// RunWithGracefulShutdown runs the application with graceful shutdown handling
func (a *App) RunWithGracefulShutdown() error {
	// Setup graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start the application
	if err := a.Run(); err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutdown signal received")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to shutdown gracefully: %w", err)
	}

	return nil
}

// initDatabase initializes the database connection and runs migrations
func (a *App) initDatabase() error {
	if a.dbInitializer == nil {
		return fmt.Errorf("database initializer not provided")
	}

	db, err := a.dbInitializer.InitializeDatabase(a.config)
	if err != nil {
		return err
	}

	if err := a.dbInitializer.RunMigrations(db); err != nil {
		return err
	}

	a.db = db
	log.Println("Database connection established and migrations completed")
	return nil
}

// initRepositories initializes all data access repositories
func (a *App) initRepositories() error {
	if a.repoInitializer == nil {
		return fmt.Errorf("repository initializer not provided")
	}

	if err := a.repoInitializer.InitializeRepositories(a.db); err != nil {
		return err
	}

	a.repositories = a.repoInitializer.GetRepositories()
	log.Println("Repositories initialized successfully")
	return nil
}

// initServices initializes all business logic services
func (a *App) initServices() error {
	if a.serviceInitializer == nil {
		return fmt.Errorf("service initializer not provided")
	}

	if err := a.serviceInitializer.InitializeServices(a.repositories, a.db); err != nil {
		return err
	}

	a.services = a.serviceInitializer.GetServices()
	log.Println("Services initialized successfully")
	return nil
}

// initHandlers initializes HTTP handlers
func (a *App) initHandlers() error {
	if a.handlerInitializer == nil {
		return fmt.Errorf("handler initializer not provided")
	}

	if err := a.handlerInitializer.InitializeHandlers(a.services); err != nil {
		return err
	}

	a.handlers = a.handlerInitializer.GetHandlers()
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

	// Configure middleware
	if a.middlewareProvider != nil {
		a.middlewareProvider.ConfigureMiddleware(router)
	} else {
		// Default middleware
		router.Use(gin.Logger())
		router.Use(gin.Recovery())
	}

	// Add health check endpoints
	if a.healthChecker != nil {
		router.GET("/health", a.healthChecker.HealthCheck)
		router.GET("/ready", a.healthChecker.ReadinessCheck)
	}

	// Register service-specific routes
	if a.routeRegistrar != nil {
		a.routeRegistrar.RegisterRoutes(router, a.handlers)
	}

	a.router = router
	log.Println("Router initialized successfully")
	return nil
}

// initServer initializes the HTTP server
func (a *App) initServer() {
	a.server = a.serverManager.CreateServer(a.config, a.router)
}

// Getter methods implementing ServiceApp interface
func (a *App) GetConfig() *config.Config {
	return a.config
}

func (a *App) GetServiceInfo() ServiceInfo {
	return a.serviceInfo
}

func (a *App) GetDB() *gorm.DB {
	return a.db
}

func (a *App) GetValidator() *validator.Validate {
	return a.validator
}

func (a *App) GetRouter() *gin.Engine {
	return a.router
}

func (a *App) GetServer() *http.Server {
	return a.server
}

// RegisterComponent registers a service component
func (a *App) RegisterComponent(name string, component ServiceComponent) error {
	return a.registry.Register(name, component)
}

// GetComponent retrieves a service component
func (a *App) GetComponent(name string) ServiceComponent {
	return a.registry.Get(name)
}

// GetRepository retrieves a repository by name
func (a *App) GetRepository(name string) interface{} {
	return a.repositories[name]
}

// GetService retrieves a service by name
func (a *App) GetService(name string) interface{} {
	return a.services[name]
}

// GetHandler retrieves a handler by name
func (a *App) GetHandler(name string) interface{} {
	return a.handlers[name]
}