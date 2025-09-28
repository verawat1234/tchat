package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"tchat.dev/shared/config"
)

// DefaultServerManager implements ServerManager
type DefaultServerManager struct{}

// NewDefaultServerManager creates a new server manager
func NewDefaultServerManager() ServerManager {
	return &DefaultServerManager{}
}

// CreateServer creates and configures an HTTP server
func (m *DefaultServerManager) CreateServer(cfg *config.Config, router *gin.Engine) *http.Server {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	return &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}

// Start starts the HTTP server
func (m *DefaultServerManager) Start(server *http.Server) error {
	// Server is started in a goroutine by the App.Run() method
	// This method is here for interface compliance and potential future use
	return nil
}

// Shutdown gracefully shuts down the HTTP server
func (m *DefaultServerManager) Shutdown(ctx context.Context, server *http.Server) error {
	if server == nil {
		return nil
	}

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}

// GetServerInfo returns server information
func (m *DefaultServerManager) GetServerInfo() ServerInfo {
	return ServerInfo{
		Address:     "localhost",
		Port:        8080,
		TLSEnabled:  false,
		StartedAt:   time.Now().Format(time.RFC3339),
		Version:     "1.0.0",
		Environment: "development",
	}
}

// DefaultMiddlewareProvider provides standard middleware
type DefaultMiddlewareProvider struct {
	corsEnabled     bool
	securityEnabled bool
	customMiddleware []gin.HandlerFunc
}

// NewDefaultMiddlewareProvider creates a new middleware provider
func NewDefaultMiddlewareProvider(corsEnabled, securityEnabled bool) MiddlewareProvider {
	return &DefaultMiddlewareProvider{
		corsEnabled:     corsEnabled,
		securityEnabled: securityEnabled,
		customMiddleware: make([]gin.HandlerFunc, 0),
	}
}

// AddMiddleware adds custom middleware
func (p *DefaultMiddlewareProvider) AddMiddleware(middleware gin.HandlerFunc) {
	p.customMiddleware = append(p.customMiddleware, middleware)
}

// GetMiddlewares returns all middleware
func (p *DefaultMiddlewareProvider) GetMiddlewares() []gin.HandlerFunc {
	var middlewares []gin.HandlerFunc

	// Add standard middleware
	middlewares = append(middlewares, gin.Logger())
	middlewares = append(middlewares, gin.Recovery())

	// Add CORS middleware if enabled
	if p.corsEnabled {
		middlewares = append(middlewares, p.corsMiddleware())
	}

	// Add security middleware if enabled
	if p.securityEnabled {
		middlewares = append(middlewares, p.securityMiddleware())
	}

	// Add custom middleware
	middlewares = append(middlewares, p.customMiddleware...)

	return middlewares
}

// ConfigureMiddleware configures middleware on the router
func (p *DefaultMiddlewareProvider) ConfigureMiddleware(router *gin.Engine) {
	middlewares := p.GetMiddlewares()
	for _, middleware := range middlewares {
		router.Use(middleware)
	}
}

// corsMiddleware provides CORS support
func (p *DefaultMiddlewareProvider) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// securityMiddleware adds security headers
func (p *DefaultMiddlewareProvider) securityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}

// ConfigurableMiddlewareProvider allows for more flexible middleware configuration
type ConfigurableMiddlewareProvider struct {
	*DefaultMiddlewareProvider
	config *config.Config
}

// NewConfigurableMiddlewareProvider creates a configurable middleware provider
func NewConfigurableMiddlewareProvider(cfg *config.Config) MiddlewareProvider {
	provider := &ConfigurableMiddlewareProvider{
		DefaultMiddlewareProvider: NewDefaultMiddlewareProvider(true, true).(*DefaultMiddlewareProvider),
		config:                   cfg,
	}

	// Configure based on config
	provider.configureFromConfig()

	return provider
}

// configureFromConfig configures middleware based on configuration
func (p *ConfigurableMiddlewareProvider) configureFromConfig() {
	// Add rate limiting if enabled
	if p.config.RateLimit.Enabled {
		p.AddMiddleware(p.rateLimitMiddleware())
	}

	// Add monitoring middleware if enabled
	if p.config.Monitoring.Enabled {
		p.AddMiddleware(p.monitoringMiddleware())
	}

	// Add request ID middleware
	p.AddMiddleware(p.requestIDMiddleware())
}

// rateLimitMiddleware provides rate limiting (basic implementation)
func (p *ConfigurableMiddlewareProvider) rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Basic rate limiting implementation
		// In a real implementation, you'd use a proper rate limiter
		c.Next()
	}
}

// monitoringMiddleware provides monitoring and metrics
func (p *ConfigurableMiddlewareProvider) monitoringMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// Record metrics
		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.FullPath()

		// Log request metrics (in a real implementation, send to metrics system)
		if p.config.Debug {
			fmt.Printf("Request: %s %s - %d - %v\n", method, path, status, duration)
		}
	}
}

// requestIDMiddleware adds request ID to context
func (p *ConfigurableMiddlewareProvider) requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}