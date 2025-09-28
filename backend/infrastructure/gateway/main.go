package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"tchat.dev/shared/config"
	"tchat.dev/shared/logger"
	"tchat.dev/shared/middleware"
)

// ServiceRegistry manages registered microservices
type ServiceRegistry struct {
	services map[string]*ServiceInstance
	mu       sync.RWMutex
}

// ServiceInstance represents a registered microservice
type ServiceInstance struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Host     string    `json:"host"`
	Port     int       `json:"port"`
	Health   string    `json:"health"`
	Version  string    `json:"version"`
	Tags     []string  `json:"tags"`
	LastSeen time.Time `json:"last_seen"`
}

// HealthStatus represents service health
type HealthStatus string

const (
	Healthy   HealthStatus = "healthy"
	Unhealthy HealthStatus = "unhealthy"
	Unknown   HealthStatus = "unknown"
)

// Gateway represents the API Gateway
type Gateway struct {
	config          *config.Config
	logger          *logger.TchatLogger
	registry        *ServiceRegistry
	router          *gin.Engine
	server          *http.Server
	healthCheckers  map[string]*HealthChecker
	loadBalancers   map[string]*LoadBalancer
	requestTimeout  time.Duration
	maxRetries      int
}

// NewGateway creates a new API Gateway instance
func NewGateway(cfg *config.Config) *Gateway {
	// Initialize logger
	logConfig := &logger.LoggerConfig{
		Level:           logger.InfoLevel,
		Format:          "json",
		ServiceName:     "api-gateway",
		ServiceVersion:  "1.0.0",
		Environment:     cfg.Environment,
		Region:          "sea-central",
		OutputPath:      "stdout",
		ComplianceConfig: logger.DefaultLoggerConfig().ComplianceConfig,
	}
	log := logger.NewTchatLogger(logConfig)

	// Initialize service registry
	registry := &ServiceRegistry{
		services: make(map[string]*ServiceInstance),
	}

	// Configure Gin
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(corsMiddleware(cfg))
	router.Use(requestIDMiddleware())
	router.Use(loggingMiddleware(log))

	return &Gateway{
		config:         cfg,
		logger:         log,
		registry:       registry,
		router:         router,
		healthCheckers: make(map[string]*HealthChecker),
		loadBalancers:  make(map[string]*LoadBalancer),
		requestTimeout: 30 * time.Second,
		maxRetries:     3,
	}
}

// Start starts the API Gateway server
func (g *Gateway) Start() error {
	// Register routes
	g.setupRoutes()

	// Register default services
	g.registerDefaultServices()

	// Start health checkers
	g.startHealthCheckers()

	// Configure server
	g.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", g.config.Server.Port),
		Handler:      g.router,
		ReadTimeout:  g.config.Server.ReadTimeout,
		WriteTimeout: g.config.Server.WriteTimeout,
		IdleTimeout:  g.config.Server.IdleTimeout,
	}

	g.logger.WithFields(logrus.Fields{
		"port": g.config.Server.Port,
		"env":  g.config.Environment,
	}).Info("Starting API Gateway")

	return g.server.ListenAndServe()
}

// Stop gracefully stops the API Gateway
func (g *Gateway) Stop(ctx context.Context) error {
	g.logger.Info("Shutting down API Gateway")

	// Stop health checkers
	g.stopHealthCheckers()

	// Shutdown server
	return g.server.Shutdown(ctx)
}

// setupRoutes configures all gateway routes
func (g *Gateway) setupRoutes() {
	// Health check endpoint
	g.router.GET("/health", g.healthHandler)
	g.router.GET("/ready", g.readinessHandler)

	// Service registry endpoints
	registry := g.router.Group("/registry")
	{
		registry.GET("/services", g.listServicesHandler)
		registry.POST("/services", g.registerServiceHandler)
		registry.DELETE("/services/:id", g.deregisterServiceHandler)
	}

	// API versioning
	v1 := g.router.Group("/api/v1")
	{
		// Auth service routes
		auth := v1.Group("/auth")
		{
			auth.Any("/*path", g.proxyHandler("auth-service"))
		}

		// Messaging service routes
		messaging := v1.Group("/messages", g.authMiddleware())
		{
			messaging.Any("/*path", g.proxyHandler("messaging-service"))
		}

		// Payment service routes
		payment := v1.Group("/payments", g.authMiddleware(), g.kycMiddleware(1))
		{
			payment.Any("/*path", g.proxyHandler("payment-service"))
		}

		// Commerce service routes
		commerce := v1.Group("/commerce", g.authMiddleware())
		{
			commerce.Any("/*path", g.proxyHandler("commerce-service"))
		}

		// Notification service routes
		notifications := v1.Group("/notifications", g.authMiddleware())
		{
			notifications.Any("/*path", g.proxyHandler("notification-service"))
		}

		// Content service routes
		content := v1.Group("/content")
		{
			content.Any("/*path", g.proxyHandler("content-service"))
		}

		// Video service routes
		videos := v1.Group("/videos")
		{
			videos.Any("/*path", g.proxyHandler("video-service"))
		}

		// Channels service routes (also handled by video-service)
		v1.GET("/channels", g.proxyHandler("video-service"))
		channels := v1.Group("/channels")
		{
			channels.Any("/*path", g.proxyHandler("video-service"))
		}

		// Social service routes
		social := v1.Group("/social", g.authMiddleware())
		{
			social.Any("/*path", g.proxyHandler("social-service"))
		}
	}

	// WebSocket proxy for real-time messaging
	g.router.GET("/ws", g.websocketProxyHandler("messaging-service"))

	// Admin endpoints
	admin := g.router.Group("/admin", g.adminAuthMiddleware())
	{
		admin.GET("/metrics", g.metricsHandler)
		admin.GET("/services/health", g.servicesHealthHandler)
		admin.POST("/services/:name/restart", g.restartServiceHandler)
	}
}

// registerDefaultServices registers the core microservices
func (g *Gateway) registerDefaultServices() {
	services := []ServiceInstance{
		{
			ID:      uuid.New().String(),
			Name:    "auth-service",
			Host:    "localhost",
			Port:    8081,
			Health:  string(Unknown),
			Version: "1.0.0",
			Tags:    []string{"auth", "core"},
		},
		{
			ID:      uuid.New().String(),
			Name:    "messaging-service",
			Host:    "localhost",
			Port:    8082,
			Health:  string(Unknown),
			Version: "1.0.0",
			Tags:    []string{"messaging", "realtime"},
		},
		{
			ID:      uuid.New().String(),
			Name:    "commerce-service",
			Host:    "localhost",
			Port:    8083,
			Health:  string(Unknown),
			Version: "1.0.0",
			Tags:    []string{"commerce", "business"},
		},
		{
			ID:      uuid.New().String(),
			Name:    "payment-service",
			Host:    "localhost",
			Port:    8084,
			Health:  string(Unknown),
			Version: "1.0.0",
			Tags:    []string{"payment", "financial"},
		},
		{
			ID:      uuid.New().String(),
			Name:    "notification-service",
			Host:    "localhost",
			Port:    8085,
			Health:  string(Unknown),
			Version: "1.0.0",
			Tags:    []string{"notification", "communication"},
		},
		{
			ID:      uuid.New().String(),
			Name:    "content-service",
			Host:    "localhost",
			Port:    8086,
			Health:  string(Unknown),
			Version: "1.0.0",
			Tags:    []string{"content", "cms"},
		},
		{
			ID:      uuid.New().String(),
			Name:    "video-service",
			Host:    "localhost",
			Port:    8091,
			Health:  string(Unknown),
			Version: "1.0.0",
			Tags:    []string{"video", "media"},
		},
		{
			ID:      uuid.New().String(),
			Name:    "social-service",
			Host:    "localhost",
			Port:    8092,
			Health:  string(Unknown),
			Version: "1.0.0",
			Tags:    []string{"social", "community"},
		},
	}

	for _, service := range services {
		g.registry.RegisterService(&service)
		g.logger.WithFields(logrus.Fields{
			"service": service.Name,
			"host":    service.Host,
			"port":    service.Port,
		}).Info("Registered service")
	}
}

// proxyHandler creates a reverse proxy handler for a service
func (g *Gateway) proxyHandler(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get service instance
		service := g.registry.GetHealthyService(serviceName)
		if service == nil {
			g.logger.WithFields(logrus.Fields{
				"service":    serviceName,
				"request_id": c.GetString("request_id"),
			}).Error("Service not available")
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "service_unavailable",
				"message": fmt.Sprintf("Service %s is not available", serviceName),
				"code":    "SERVICE_UNAVAILABLE",
			})
			return
		}

		// Create proxy target
		target, err := url.Parse(fmt.Sprintf("http://%s:%d", service.Host, service.Port))
		if err != nil {
			g.logger.WithFields(logrus.Fields{
				"service": serviceName,
				"error":   err.Error(),
			}).Error("Failed to parse service URL")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Failed to proxy request",
				"code":    "PROXY_ERROR",
			})
			return
		}

		// Create reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(target)

		// Configure proxy
		proxy.Director = func(req *http.Request) {
			req.URL.Host = target.Host
			req.URL.Scheme = target.Scheme
			req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
			req.Header.Set("X-Forwarded-Proto", "http")
			req.Header.Set("X-Gateway-Service", serviceName)

			// Add request ID
			if requestID := c.GetString("request_id"); requestID != "" {
				req.Header.Set("X-Request-ID", requestID)
			}

			// Forward user context
			if userID := c.GetString("user_id"); userID != "" {
				req.Header.Set("X-User-ID", userID)
			}
			if countryCode := c.GetString("country_code"); countryCode != "" {
				req.Header.Set("X-Country-Code", countryCode)
			}
		}

		proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
			g.logger.WithFields(logrus.Fields{
				"service": serviceName,
				"error":   err.Error(),
				"path":    req.URL.Path,
			}).Error("Proxy error")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error":"bad_gateway","message":"Service temporarily unavailable","code":"BAD_GATEWAY"}`))
		}

		// Log request
		start := time.Now()
		defer func() {
			duration := time.Since(start)
			g.logger.PerformanceLog(fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path), duration, map[string]interface{}{
				"service":      serviceName,
				"service_host": service.Host,
				"service_port": service.Port,
				"status_code":  c.Writer.Status(),
				"method":       c.Request.Method,
				"path":         c.Request.URL.Path,
			})
		}()

		// Forward the request
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// websocketProxyHandler handles WebSocket connections
func (g *Gateway) websocketProxyHandler(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get service instance
		service := g.registry.GetHealthyService(serviceName)
		if service == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "service_unavailable",
				"message": fmt.Sprintf("WebSocket service %s is not available", serviceName),
			})
			return
		}

		// Proxy WebSocket connection
		target := fmt.Sprintf("ws://%s:%d/ws", service.Host, service.Port)

		g.logger.WithFields(logrus.Fields{
			"service":    serviceName,
			"target":     target,
			"request_id": c.GetString("request_id"),
			"event":      "connection_proxied",
		}).Info("WebSocket connection proxied")

		// In a real implementation, this would handle WebSocket proxying
		// For now, return the target information
		c.JSON(http.StatusOK, gin.H{
			"websocket_url": target,
			"service":       serviceName,
		})
	}
}

// Middleware functions

func (g *Gateway) authMiddleware() gin.HandlerFunc {
	authMiddleware := middleware.NewAuthMiddleware(g.config)
	return authMiddleware.RequireAuth()
}

func (g *Gateway) kycMiddleware(minLevel int) gin.HandlerFunc {
	authMiddleware := middleware.NewAuthMiddleware(g.config)
	return authMiddleware.RequireKYC(minLevel)
}

func (g *Gateway) adminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple admin authentication - in production use proper admin auth
		token := c.GetHeader("X-Admin-Token")
		if token != "admin-secret-token" { // This should come from config
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Admin access required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Handler functions

func (g *Gateway) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "api-gateway",
		"version":   g.config.Version,
	})
}

func (g *Gateway) readinessHandler(c *gin.Context) {
	healthyServices := 0
	totalServices := 0

	g.registry.mu.RLock()
	for _, service := range g.registry.services {
		totalServices++
		if service.Health == string(Healthy) {
			healthyServices++
		}
	}
	g.registry.mu.RUnlock()

	ready := healthyServices > 0 && float64(healthyServices)/float64(totalServices) >= 0.5

	status := http.StatusOK
	if !ready {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{
		"ready":            ready,
		"healthy_services": healthyServices,
		"total_services":   totalServices,
		"timestamp":        time.Now().UTC(),
	})
}

func (g *Gateway) listServicesHandler(c *gin.Context) {
	services := g.registry.GetAllServices()
	c.JSON(http.StatusOK, gin.H{
		"services": services,
		"count":    len(services),
	})
}

func (g *Gateway) registerServiceHandler(c *gin.Context) {
	var service ServiceInstance
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	service.ID = uuid.New().String()
	service.LastSeen = time.Now()
	service.Health = string(Unknown)

	g.registry.RegisterService(&service)

	g.logger.WithFields(logrus.Fields{
		"service_id":   service.ID,
		"service_name": service.Name,
		"host":         service.Host,
		"port":         service.Port,
	}).Info("Service registered")

	c.JSON(http.StatusCreated, service)
}

func (g *Gateway) deregisterServiceHandler(c *gin.Context) {
	serviceID := c.Param("id")
	if g.registry.DeregisterService(serviceID) {
		g.logger.WithFields(logrus.Fields{
			"service_id": serviceID,
		}).Info("Service deregistered")
		c.JSON(http.StatusOK, gin.H{"message": "Service deregistered"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "not_found",
			"message": "Service not found",
		})
	}
}

func (g *Gateway) metricsHandler(c *gin.Context) {
	// Return basic metrics - in production integrate with Prometheus
	c.JSON(http.StatusOK, gin.H{
		"requests_total":   0, // Would track actual metrics
		"request_duration": 0,
		"services_healthy": g.registry.GetHealthyServiceCount(),
		"services_total":   g.registry.GetServiceCount(),
	})
}

func (g *Gateway) servicesHealthHandler(c *gin.Context) {
	services := g.registry.GetAllServices()
	healthStatus := make(map[string]interface{})

	for _, service := range services {
		healthStatus[service.Name] = gin.H{
			"status":    service.Health,
			"last_seen": service.LastSeen,
			"endpoint":  fmt.Sprintf("%s:%d", service.Host, service.Port),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"services": healthStatus,
		"summary": gin.H{
			"total":   len(services),
			"healthy": g.registry.GetHealthyServiceCount(),
		},
	})
}

func (g *Gateway) restartServiceHandler(c *gin.Context) {
	serviceName := c.Param("name")

	// This would trigger service restart - implementation depends on orchestration
	g.logger.WithFields(logrus.Fields{
		"service": serviceName,
		"admin":   c.GetHeader("X-Admin-User"),
	}).Info("Service restart requested")

	c.JSON(http.StatusAccepted, gin.H{
		"message": fmt.Sprintf("Restart request submitted for %s", serviceName),
		"service": serviceName,
	})
}

// Utility middleware

func corsMiddleware(cfg *config.Config) gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = cfg.Server.CORS.AllowOrigins
	config.AllowMethods = cfg.Server.CORS.AllowMethods
	config.AllowHeaders = cfg.Server.CORS.AllowHeaders
	config.ExposeHeaders = cfg.Server.CORS.ExposeHeaders
	config.AllowCredentials = cfg.Server.CORS.AllowCredentials
	config.MaxAge = time.Duration(cfg.Server.CORS.MaxAge) * time.Second

	return cors.New(config)
}

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func loggingMiddleware(logger *logger.TchatLogger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.PerformanceLog(
			fmt.Sprintf("%s %s", param.Method, param.Path),
			param.Latency,
			map[string]interface{}{
				"client_ip":   param.ClientIP,
				"user_agent":  param.Request.UserAgent(),
				"status_code": param.StatusCode,
				"method":      param.Method,
				"path":        param.Path,
			},
		)
		return ""
	})
}

// Health checking

func (g *Gateway) startHealthCheckers() {
	g.registry.mu.RLock()
	defer g.registry.mu.RUnlock()

	for _, service := range g.registry.services {
		checker := NewHealthChecker(service, g.logger)
		g.healthCheckers[service.ID] = checker
		go checker.Start(g.registry)
	}
}

func (g *Gateway) stopHealthCheckers() {
	for _, checker := range g.healthCheckers {
		checker.Stop()
	}
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create gateway
	gateway := NewGateway(cfg)

	// Setup graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start gateway in goroutine
	go func() {
		if err := gateway.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start gateway: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := gateway.Stop(shutdownCtx); err != nil {
		log.Printf("Gateway shutdown error: %v", err)
	}

	log.Println("Gateway stopped")
}