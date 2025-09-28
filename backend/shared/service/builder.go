package service

import (
	"fmt"

	"tchat.dev/shared/config"
)

// AppBuilder provides a fluent interface for building applications
type AppBuilder struct {
	config             *config.Config
	serviceInfo        ServiceInfo
	dbInitializer      DatabaseInitializer
	repoInitializer    RepositoryInitializer
	serviceInitializer ServiceInitializer
	handlerInitializer HandlerInitializer
	routeRegistrar     RouteRegistrar
	healthChecker      HealthChecker
	middlewareProvider MiddlewareProvider
	components         map[string]ServiceComponent
}

// NewAppBuilder creates a new application builder
func NewAppBuilder(cfg *config.Config) *AppBuilder {
	return &AppBuilder{
		config:     cfg,
		components: make(map[string]ServiceComponent),
	}
}

// WithServiceInfo sets the service information
func (b *AppBuilder) WithServiceInfo(info ServiceInfo) *AppBuilder {
	b.serviceInfo = info
	return b
}

// WithDatabaseInitializer sets the database initializer
func (b *AppBuilder) WithDatabaseInitializer(initializer DatabaseInitializer) *AppBuilder {
	b.dbInitializer = initializer
	return b
}

// WithRepositoryInitializer sets the repository initializer
func (b *AppBuilder) WithRepositoryInitializer(initializer RepositoryInitializer) *AppBuilder {
	b.repoInitializer = initializer
	return b
}

// WithServiceInitializer sets the service initializer
func (b *AppBuilder) WithServiceInitializer(initializer ServiceInitializer) *AppBuilder {
	b.serviceInitializer = initializer
	return b
}

// WithHandlerInitializer sets the handler initializer
func (b *AppBuilder) WithHandlerInitializer(initializer HandlerInitializer) *AppBuilder {
	b.handlerInitializer = initializer
	return b
}

// WithRouteRegistrar sets the route registrar
func (b *AppBuilder) WithRouteRegistrar(registrar RouteRegistrar) *AppBuilder {
	b.routeRegistrar = registrar
	return b
}

// WithHealthChecker sets the health checker
func (b *AppBuilder) WithHealthChecker(checker HealthChecker) *AppBuilder {
	b.healthChecker = checker
	return b
}

// WithMiddlewareProvider sets the middleware provider
func (b *AppBuilder) WithMiddlewareProvider(provider MiddlewareProvider) *AppBuilder {
	b.middlewareProvider = provider
	return b
}

// WithComponent adds a service component
func (b *AppBuilder) WithComponent(name string, component ServiceComponent) *AppBuilder {
	b.components[name] = component
	return b
}

// WithDefaultHealthChecker creates and sets a default health checker
func (b *AppBuilder) WithDefaultHealthChecker() *AppBuilder {
	// Will be created after the app is built since it needs db and registry
	return b
}

// WithDefaultMiddleware creates and sets default middleware
func (b *AppBuilder) WithDefaultMiddleware() *AppBuilder {
	b.middlewareProvider = NewDefaultMiddlewareProvider(true, true)
	return b
}

// WithConfigurableMiddleware creates configurable middleware from config
func (b *AppBuilder) WithConfigurableMiddleware() *AppBuilder {
	b.middlewareProvider = NewConfigurableMiddlewareProvider(b.config)
	return b
}

// Build creates the application instance
func (b *AppBuilder) Build() (*App, error) {
	// Validate required components
	if err := b.validate(); err != nil {
		return nil, err
	}

	// Create app configuration
	appCfg := AppConfig{
		ServiceInfo:           b.serviceInfo,
		DatabaseInitializer:   b.dbInitializer,
		RepositoryInitializer: b.repoInitializer,
		ServiceInitializer:    b.serviceInitializer,
		HandlerInitializer:    b.handlerInitializer,
		RouteRegistrar:        b.routeRegistrar,
		HealthChecker:         b.healthChecker,
		MiddlewareProvider:    b.middlewareProvider,
	}

	// Create the application
	app := NewApp(b.config, appCfg)

	// Register components
	for name, component := range b.components {
		if err := app.RegisterComponent(name, component); err != nil {
			return nil, fmt.Errorf("failed to register component %s: %w", name, err)
		}
	}

	// Create default health checker if none provided
	if b.healthChecker == nil {
		registry := app.registry
		healthChecker := NewDefaultHealthChecker(
			b.serviceInfo.Name,
			b.serviceInfo.Version,
			nil, // Will be set after initialization
			b.config,
			registry,
		)
		app.healthChecker = healthChecker
	}

	return app, nil
}

// BuildAndRun builds the application and runs it with graceful shutdown
func (b *AppBuilder) BuildAndRun() error {
	app, err := b.Build()
	if err != nil {
		return fmt.Errorf("failed to build application: %w", err)
	}

	if err := app.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}

	return app.RunWithGracefulShutdown()
}

// validate validates the builder configuration
func (b *AppBuilder) validate() error {
	if b.config == nil {
		return fmt.Errorf("configuration is required")
	}

	if b.serviceInfo.Name == "" {
		return fmt.Errorf("service name is required")
	}

	if b.serviceInfo.Version == "" {
		b.serviceInfo.Version = "1.0.0" // Default version
	}

	if b.serviceInfo.Port == 0 {
		b.serviceInfo.Port = b.config.Server.Port
	}

	// Ensure we have at least basic components
	if b.dbInitializer == nil {
		return fmt.Errorf("database initializer is required")
	}

	if b.routeRegistrar == nil {
		return fmt.Errorf("route registrar is required")
	}

	// Set defaults for optional components
	if b.middlewareProvider == nil {
		b.middlewareProvider = NewDefaultMiddlewareProvider(true, true)
	}

	return nil
}

// ServiceBuilder provides a specialized builder for microservices
type ServiceBuilder struct {
	*AppBuilder
	models             []interface{}
	defaultPort        int
	enableHealthChecks bool
	enableMetrics      bool
}

// NewServiceBuilder creates a new service builder
func NewServiceBuilder(serviceName string, cfg *config.Config) *ServiceBuilder {
	return &ServiceBuilder{
		AppBuilder:         NewAppBuilder(cfg),
		defaultPort:        8080,
		enableHealthChecks: true,
		enableMetrics:      false,
	}
}

// WithModels sets the database models for the service
func (b *ServiceBuilder) WithModels(models ...interface{}) *ServiceBuilder {
	b.models = append(b.models, models...)
	return b
}

// WithDefaultPort sets the default port for the service
func (b *ServiceBuilder) WithDefaultPort(port int) *ServiceBuilder {
	b.defaultPort = port
	return b
}

// WithHealthChecks enables or disables health checks
func (b *ServiceBuilder) WithHealthChecks(enabled bool) *ServiceBuilder {
	b.enableHealthChecks = enabled
	return b
}

// WithMetrics enables or disables metrics collection
func (b *ServiceBuilder) WithMetrics(enabled bool) *ServiceBuilder {
	b.enableMetrics = enabled
	return b
}

// Build creates the service application
func (b *ServiceBuilder) Build() (*App, error) {
	// Set up service info if not already set
	if b.serviceInfo.Name == "" {
		return nil, fmt.Errorf("service name is required")
	}

	// Set default port if not configured
	if b.config.Server.Port == 0 {
		b.config.Server.Port = b.defaultPort
	}

	// Set up database initializer with models
	if b.dbInitializer == nil && len(b.models) > 0 {
		b.dbInitializer = NewDefaultDatabaseInitializer(b.models)
	}

	// Set up default middleware if not provided
	if b.middlewareProvider == nil {
		b.middlewareProvider = NewConfigurableMiddlewareProvider(b.config)
	}

	// Add database health component if health checks are enabled
	if b.enableHealthChecks {
		dbManager := NewDefaultDatabaseManager()
		healthComponent := NewDatabaseHealthComponent(dbManager)
		b.WithComponent("database-health", healthComponent)
	}

	return b.AppBuilder.Build()
}