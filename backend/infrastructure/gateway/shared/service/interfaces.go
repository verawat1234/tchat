package service

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"tchat.dev/shared/config"
)

// ServiceApp defines the main application interface
type ServiceApp interface {
	// Lifecycle methods
	Initialize() error
	Run() error
	Shutdown(ctx context.Context) error

	// Configuration
	GetConfig() *config.Config
	GetServiceInfo() ServiceInfo

	// Database access
	GetDB() *gorm.DB
	GetValidator() *validator.Validate

	// HTTP components
	GetRouter() *gin.Engine
	GetServer() *http.Server
}

// ServiceInfo holds metadata about the service
type ServiceInfo struct {
	Name        string
	Version     string
	Description string
	Port        int
}

// DatabaseInitializer handles database setup and migrations
type DatabaseInitializer interface {
	InitializeDatabase(cfg *config.Config) (*gorm.DB, error)
	RunMigrations(db *gorm.DB) error
	GetModels() []interface{}
}

// RepositoryInitializer handles repository setup
type RepositoryInitializer interface {
	InitializeRepositories(db *gorm.DB) error
	GetRepositories() map[string]interface{}
}

// ServiceInitializer handles business service setup
type ServiceInitializer interface {
	InitializeServices(repos map[string]interface{}, db *gorm.DB) error
	GetServices() map[string]interface{}
}

// HandlerInitializer handles HTTP handler setup
type HandlerInitializer interface {
	InitializeHandlers(services map[string]interface{}) error
	GetHandlers() map[string]interface{}
}

// RouteRegistrar handles route registration
type RouteRegistrar interface {
	RegisterRoutes(router *gin.Engine, handlers map[string]interface{})
}

// HealthChecker provides health and readiness checks
type HealthChecker interface {
	HealthCheck(c *gin.Context)
	ReadinessCheck(c *gin.Context)
	GetHealthData() map[string]interface{}
}

// MiddlewareProvider provides middleware configuration
type MiddlewareProvider interface {
	GetMiddlewares() []gin.HandlerFunc
	ConfigureMiddleware(router *gin.Engine)
}

// ServiceComponent represents a pluggable service component
type ServiceComponent interface {
	Initialize(ctx context.Context, cfg *config.Config, db *gorm.DB) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Name() string
	IsHealthy() bool
}

// ServiceRegistry manages service components
type ServiceRegistry interface {
	Register(name string, component ServiceComponent) error
	Unregister(name string) error
	Get(name string) ServiceComponent
	List() map[string]ServiceComponent
	StartAll(ctx context.Context) error
	StopAll(ctx context.Context) error
}

// ConfigurationProvider handles configuration management
type ConfigurationProvider interface {
	LoadConfiguration() (*config.Config, error)
	ValidateConfiguration(cfg *config.Config) error
	GetServicePort() int
	GetServiceName() string
}

// DatabaseManager handles database operations
type DatabaseManager interface {
	Connect(cfg *config.Config) (*gorm.DB, *sql.DB, error)
	ConfigureConnection(db *sql.DB, cfg *config.Config) error
	Migrate(db *gorm.DB, models []interface{}) error
	Close() error
	Ping() error
	GetConnection() *gorm.DB
}

// ServerManager handles HTTP server lifecycle
type ServerManager interface {
	CreateServer(cfg *config.Config, router *gin.Engine) *http.Server
	Start(server *http.Server) error
	Shutdown(ctx context.Context, server *http.Server) error
	GetServerInfo() ServerInfo
}

// ServerInfo holds server metadata
type ServerInfo struct {
	Address     string
	Port        int
	TLSEnabled  bool
	StartedAt   string
	Version     string
	Environment string
}

// LoggerProvider handles logging configuration
type LoggerProvider interface {
	ConfigureLogger(cfg *config.Config) error
	GetLogger() interface{}
}

// MetricsProvider handles metrics collection
type MetricsProvider interface {
	InitializeMetrics(cfg *config.Config) error
	RecordMetric(name string, value float64, labels map[string]string)
	GetMetrics() map[string]interface{}
}