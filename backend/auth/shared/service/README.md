# Tchat Microservice Framework

A reusable, generic service framework for building microservices in the Tchat application. This framework provides a clean architecture with dependency injection, interface segregation, and standardized patterns for configuration management, error handling, logging, and testing support.

## Features

- **Clean Architecture**: Separation of concerns with clear interfaces
- **Dependency Injection**: Interface-based dependency management
- **Configuration Management**: Centralized configuration with validation
- **Database Management**: GORM-based database handling with connection pooling
- **Health Checks**: Standardized health and readiness endpoints
- **Middleware Factory**: Reusable middleware components (CORS, security, rate limiting)
- **Service Registry**: Dynamic service component registration and lifecycle management
- **Graceful Shutdown**: Proper resource cleanup and shutdown handling
- **Builder Pattern**: Fluent API for service construction
- **Testing Support**: Interfaces designed for easy mocking and testing

## Architecture

### Core Interfaces

```go
// Main application interface
type ServiceApp interface {
    Initialize() error
    Run() error
    Shutdown(ctx context.Context) error
    GetConfig() *config.Config
    GetServiceInfo() ServiceInfo
}

// Component lifecycle management
type ServiceComponent interface {
    Initialize(ctx context.Context, cfg *config.Config, db *gorm.DB) error
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Name() string
    IsHealthy() bool
}
```

### Framework Components

1. **App**: Main application orchestrator
2. **ServiceRegistry**: Manages pluggable service components
3. **DatabaseManager**: Handles database connections and migrations
4. **ServerManager**: HTTP server lifecycle management
5. **HealthChecker**: Health and readiness monitoring
6. **MiddlewareProvider**: Configurable middleware stack

## Quick Start

### Basic Service Creation

```go
package main

import (
    "tchat.dev/shared/config"
    "tchat.dev/shared/service"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Create service using builder pattern
    err = service.NewServiceBuilder("my-service", cfg).
        WithModels(&MyModel{}).
        WithDefaultPort(8080).
        WithServiceInfo(service.ServiceInfo{
            Name:        "my-service",
            Version:     "1.0.0",
            Description: "My microservice",
        }).
        WithRepositoryInitializer(NewMyRepositoryInitializer()).
        WithServiceInitializer(NewMyServiceInitializer()).
        WithHandlerInitializer(NewMyHandlerInitializer()).
        WithRouteRegistrar(NewMyRouteRegistrar()).
        WithConfigurableMiddleware().
        WithHealthChecks(true).
        BuildAndRun()

    if err != nil {
        log.Fatalf("Service failed: %v", err)
    }
}
```

### Implementing Service Components

#### 1. Repository Initializer

```go
type MyRepositoryInitializer struct {
    repositories map[string]interface{}
}

func NewMyRepositoryInitializer() service.RepositoryInitializer {
    return &MyRepositoryInitializer{
        repositories: make(map[string]interface{}),
    }
}

func (r *MyRepositoryInitializer) InitializeRepositories(db *gorm.DB) error {
    myRepo := NewMyRepository(db)
    r.repositories["my"] = myRepo
    return nil
}

func (r *MyRepositoryInitializer) GetRepositories() map[string]interface{} {
    return r.repositories
}
```

#### 2. Service Initializer

```go
type MyServiceInitializer struct {
    services map[string]interface{}
}

func NewMyServiceInitializer() service.ServiceInitializer {
    return &MyServiceInitializer{
        services: make(map[string]interface{}),
    }
}

func (s *MyServiceInitializer) InitializeServices(repos map[string]interface{}, db *gorm.DB) error {
    myRepo := repos["my"].(MyRepository)
    myService := NewMyService(myRepo)
    s.services["my"] = myService
    return nil
}
```

#### 3. Handler Initializer

```go
type MyHandlerInitializer struct {
    handlers map[string]interface{}
}

func (h *MyHandlerInitializer) InitializeHandlers(services map[string]interface{}) error {
    myService := services["my"].(*MyService)
    myHandler := NewMyHandler(myService)
    h.handlers["my"] = myHandler
    return nil
}
```

#### 4. Route Registrar

```go
type MyRouteRegistrar struct{}

func (r *MyRouteRegistrar) RegisterRoutes(router *gin.Engine, handlers map[string]interface{}) {
    myHandler := handlers["my"].(*MyHandler)

    v1 := router.Group("/api/v1")
    {
        my := v1.Group("/my")
        {
            my.GET("", myHandler.GetAll)
            my.POST("", myHandler.Create)
            my.GET("/:id", myHandler.GetByID)
        }
    }
}
```

## Advanced Features

### Custom Service Components

```go
// Create a background service component
backgroundTask := service.NewBackgroundService("data-sync", func(ctx context.Context) error {
    // Your background task logic here
    return nil
})

app.RegisterComponent("data-sync", backgroundTask)
```

### Custom Health Checks

```go
healthChecker := service.NewDefaultHealthChecker("my-service", "1.0.0", db, cfg, registry)
healthChecker.AddCustomCheck("external-api", func() error {
    // Check external API connectivity
    return checkExternalAPI()
})
```

### Custom Middleware

```go
middlewareProvider := service.NewDefaultMiddlewareProvider(true, true)
middlewareProvider.AddMiddleware(func(c *gin.Context) {
    // Your custom middleware logic
    c.Next()
})
```

## Configuration

The framework uses the shared configuration system from `tchat.dev/shared/config`. Services can override specific settings:

```go
cfg.Server.Port = 8091  // Override port for this service
```

## Health Endpoints

The framework automatically provides:

- `GET /health` - Basic health check
- `GET /ready` - Readiness check (includes database, components, custom checks)

Health check response example:
```json
{
    "status": "ready",
    "service": "my-service",
    "timestamp": "2024-01-01T00:00:00Z",
    "checks": {
        "database": {"status": "healthy"},
        "components": {
            "data-sync": {"status": "healthy"}
        },
        "custom": {
            "external-api": {"status": "healthy"}
        }
    }
}
```

## Testing

The framework is designed for easy testing with interface-based dependency injection:

```go
func TestMyService(t *testing.T) {
    // Mock repository
    mockRepo := &MockMyRepository{}

    // Create service with mock
    service := NewMyService(mockRepo)

    // Test service logic
    // ...
}
```

## Migration from Existing Services

To migrate an existing service to use the framework:

1. Extract your models, repositories, services, and handlers
2. Implement the framework interfaces (RepositoryInitializer, ServiceInitializer, etc.)
3. Replace your main.go with the builder pattern
4. Update imports to use the framework components

See `backend/video/main_refactored.go` for a complete example of migrating the video service.

## Framework Components Reference

### ServiceBuilder

- `WithModels()` - Set database models for migration
- `WithDefaultPort()` - Set default service port
- `WithServiceInfo()` - Set service metadata
- `WithRepositoryInitializer()` - Set repository initialization logic
- `WithServiceInitializer()` - Set business service initialization
- `WithHandlerInitializer()` - Set HTTP handler initialization
- `WithRouteRegistrar()` - Set route registration logic
- `WithConfigurableMiddleware()` - Use config-based middleware
- `WithDefaultMiddleware()` - Use default middleware stack
- `WithHealthChecks()` - Enable/disable health checking
- `WithMetrics()` - Enable/disable metrics collection
- `WithComponent()` - Add custom service components

### Middleware Options

- CORS support with configurable origins
- Security headers (XSS, CSRF, etc.)
- Rate limiting (when enabled in config)
- Request ID tracking
- Monitoring and metrics collection
- Custom middleware injection

### Database Features

- Connection pooling with configurable limits
- Automatic migrations
- Health checking
- Graceful connection closure
- GORM integration

## Examples

See `backend/shared/service/examples/example_service.go` for complete working examples of:

- Basic service creation
- Custom components
- Different builder patterns
- Repository/Service/Handler implementations

## Best Practices

1. **Interface Segregation**: Keep interfaces focused and small
2. **Dependency Injection**: Use interfaces, not concrete types
3. **Error Handling**: Return errors, don't panic
4. **Configuration**: Use environment variables through the config system
5. **Testing**: Write unit tests with mocked dependencies
6. **Logging**: Use structured logging
7. **Health Checks**: Implement meaningful health checks for your service
8. **Graceful Shutdown**: Ensure proper resource cleanup