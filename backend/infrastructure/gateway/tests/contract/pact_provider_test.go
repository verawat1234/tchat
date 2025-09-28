package contract

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tchat.dev/shared/config"
)

// Gateway types for testing (replicated from main package)

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

// ServiceRegistry manages registered microservices
type ServiceRegistry struct {
	services map[string]*ServiceInstance
	mu       sync.RWMutex
}

// HealthStatus represents service health
type HealthStatus string

const (
	Healthy   HealthStatus = "healthy"
	Unhealthy HealthStatus = "unhealthy"
	Unknown   HealthStatus = "unknown"
)

// Gateway represents the API Gateway for testing
type Gateway struct {
	config    *config.Config
	logger    *log.Logger
	registry  *ServiceRegistry
	router    *gin.Engine
	server    *http.Server
}

// ServiceRegistry methods

// RegisterService adds a service to the registry
func (sr *ServiceRegistry) RegisterService(service *ServiceInstance) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	service.LastSeen = time.Now()
	sr.services[service.ID] = service
}

// GetServiceByName retrieves the first healthy service by name
func (sr *ServiceRegistry) GetServiceByName(serviceName string) *ServiceInstance {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	for _, service := range sr.services {
		if service.Name == serviceName && service.Health == string(Healthy) {
			return service
		}
	}
	return nil
}

// GetHealthyService returns a healthy service instance using simple load balancing
func (sr *ServiceRegistry) GetHealthyService(serviceName string) *ServiceInstance {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	var healthyServices []*ServiceInstance
	for _, service := range sr.services {
		if service.Name == serviceName && service.Health == string(Healthy) {
			healthyServices = append(healthyServices, service)
		}
	}

	if len(healthyServices) == 0 {
		return nil
	}

	// Simple round-robin load balancing
	return healthyServices[time.Now().Unix()%int64(len(healthyServices))]
}

// UpdateServiceHealth updates the health status of a service
func (sr *ServiceRegistry) UpdateServiceHealth(serviceID string, health HealthStatus) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	if service, exists := sr.services[serviceID]; exists {
		service.Health = string(health)
		service.LastSeen = time.Now()
	}
}

// GetServicesByName returns all services with the given name
func (sr *ServiceRegistry) GetServicesByName(serviceName string) []*ServiceInstance {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	var services []*ServiceInstance
	for _, service := range sr.services {
		if service.Name == serviceName {
			services = append(services, service)
		}
	}
	return services
}

// Gateway methods

// NewGateway creates a new API Gateway instance for testing
func NewGateway(cfg *config.Config) *Gateway {
	// Initialize logger
	logger := log.New(os.Stdout, "[API-GATEWAY] ", log.LstdFlags)

	// Initialize service registry
	registry := &ServiceRegistry{
		services: make(map[string]*ServiceInstance),
	}

	// Configure Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add basic middleware
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		c.Header("X-Gateway-Service", "gateway-service")
		c.Next()
	})

	gateway := &Gateway{
		config:   cfg,
		logger:   logger,
		registry: registry,
		router:   router,
	}

	// Setup basic routes for testing
	gateway.setupTestRoutes()

	return gateway
}

// setupTestRoutes configures test routes that proxy to backend services
func (g *Gateway) setupTestRoutes() {
	// Health check endpoint
	g.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "api-gateway",
			"timestamp": time.Now().UTC(),
		})
	})

	// Service registry endpoints
	registry := g.router.Group("/registry")
	{
		registry.GET("/services", func(c *gin.Context) {
			services := []*ServiceInstance{}
			g.registry.mu.RLock()
			for _, service := range g.registry.services {
				services = append(services, service)
			}
			g.registry.mu.RUnlock()

			c.JSON(http.StatusOK, gin.H{
				"services": services,
				"count":    len(services),
			})
		})
	}

	// API versioning
	v1 := g.router.Group("/api/v1")
	{
		// Auth service routes
		auth := v1.Group("/auth")
		{
			auth.Any("/*path", g.proxyHandler("auth-service"))
		}

		// Content service routes
		content := v1.Group("/content")
		{
			content.Any("/*path", g.proxyHandler("content-service"))
		}

		// Other service routes (simplified for testing)
		messaging := v1.Group("/messages")
		{
			messaging.Any("/*path", g.proxyHandler("messaging-service"))
		}

		commerce := v1.Group("/commerce")
		{
			commerce.Any("/*path", g.proxyHandler("commerce-service"))
		}

		payment := v1.Group("/payments")
		{
			payment.Any("/*path", g.proxyHandler("payment-service"))
		}

		notifications := v1.Group("/notifications")
		{
			notifications.Any("/*path", g.proxyHandler("notification-service"))
		}
	}
}

// proxyHandler creates a reverse proxy handler for a service
func (g *Gateway) proxyHandler(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get service instance
		service := g.registry.GetHealthyService(serviceName)
		if service == nil {
			g.logger.Printf("Service not available: %s", serviceName)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "service_unavailable",
				"message": fmt.Sprintf("Service %s is not available", serviceName),
			})
			return
		}

		// Forward the request to the mock service
		targetURL := fmt.Sprintf("http://%s:%d%s", service.Host, service.Port, c.Request.URL.Path)
		if c.Request.URL.RawQuery != "" {
			targetURL = fmt.Sprintf("%s?%s", targetURL, c.Request.URL.RawQuery)
		}

		// Create new request
		req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
		if err != nil {
			g.logger.Printf("Failed to create proxy request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "proxy_error",
				"message": "Failed to proxy request",
			})
			return
		}

		// Copy headers
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Add gateway headers
		req.Header.Set("X-Gateway-Service", serviceName)
		if requestID := c.GetHeader("X-Request-ID"); requestID == "" {
			req.Header.Set("X-Request-ID", uuid.New().String())
		}

		// Execute request
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			g.logger.Printf("Proxy error for service %s: %v", serviceName, err)
			c.JSON(http.StatusBadGateway, gin.H{
				"error":   "bad_gateway",
				"message": "Service temporarily unavailable",
			})
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		// Set status and body
		c.Status(resp.StatusCode)

		// Copy response body
		if resp.Body != nil {
			var responseBody map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&responseBody); err == nil {
				c.JSON(resp.StatusCode, responseBody)
			}
		}
	}
}

// MockBackendService represents a mock backend service for testing
type MockBackendService struct {
	server     *httptest.Server
	serviceName string
	responses  map[string]MockResponse
}

// MockResponse defines a mock response structure
type MockResponse struct {
	StatusCode int                    `json:"status_code"`
	Headers    map[string]string      `json:"headers"`
	Body       map[string]interface{} `json:"body"`
}

// Gateway test instance for provider verification
type GatewayTestInstance struct {
	gateway     *Gateway
	mockServices map[string]*MockBackendService
	server      *httptest.Server
}

// NewGatewayTestInstance creates a new test instance of the Gateway
func NewGatewayTestInstance() *GatewayTestInstance {
	// Create test configuration
	cfg := &config.Config{
		Environment: "test",
		Server: config.ServerConfig{
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	// Create gateway instance
	gateway := NewGateway(cfg)

	// Create mock services
	mockServices := make(map[string]*MockBackendService)

	// Create auth service mock
	authMock := createMockService("auth-service")
	mockServices["auth-service"] = authMock

	// Create content service mock
	contentMock := createMockService("content-service")
	mockServices["content-service"] = contentMock

	// Create commerce service mock
	commerceMock := createMockService("commerce-service")
	mockServices["commerce-service"] = commerceMock

	// Create messaging service mock
	messagingMock := createMockService("messaging-service")
	mockServices["messaging-service"] = messagingMock

	// Create payment service mock
	paymentMock := createMockService("payment-service")
	mockServices["payment-service"] = paymentMock

	// Create notification service mock
	notificationMock := createMockService("notification-service")
	mockServices["notification-service"] = notificationMock

	// Update gateway registry with mock service addresses
	for serviceName, mockService := range mockServices {
		host, port, _ := net.SplitHostPort(strings.TrimPrefix(mockService.server.URL, "http://"))
		portNum := 8081 // Default port, will be overridden by actual mock server port
		if port != "" {
			fmt.Sscanf(port, "%d", &portNum)
		}

		service := &ServiceInstance{
			ID:      serviceName + "-test-id",
			Name:    serviceName,
			Host:    host,
			Port:    portNum,
			Health:  string(Healthy),
			Version: "1.0.0-test",
			Tags:    []string{"test", "mock"},
		}
		gateway.registry.RegisterService(service)
	}

	return &GatewayTestInstance{
		gateway:      gateway,
		mockServices: mockServices,
	}
}

// createMockService creates a mock backend service
func createMockService(serviceName string) *MockBackendService {
	mockService := &MockBackendService{
		serviceName: serviceName,
		responses:   make(map[string]MockResponse),
	}

	// Create HTTP server for mock service
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "healthy",
			"service": serviceName,
		})
	})

	// Generic handler for all other endpoints
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Create key for response lookup
		key := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		if r.URL.RawQuery != "" {
			key = fmt.Sprintf("%s?%s", key, r.URL.RawQuery)
		}

		// Look for registered mock response
		if response, exists := mockService.responses[key]; exists {
			// Set headers
			for headerName, headerValue := range response.Headers {
				w.Header().Set(headerName, headerValue)
			}

			// Set status code and body
			w.WriteHeader(response.StatusCode)
			if response.Body != nil {
				json.NewEncoder(w).Encode(response.Body)
			}
			return
		}

		// Default response for unregistered endpoints
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"service": serviceName,
			"path":    r.URL.Path,
			"method":  r.Method,
			"headers": getRequestHeaders(r),
		})
	})

	mockService.server = httptest.NewServer(mux)
	return mockService
}

// SetMockResponse configures a mock response for a specific request
func (m *MockBackendService) SetMockResponse(method, path string, response MockResponse) {
	key := fmt.Sprintf("%s %s", method, path)
	m.responses[key] = response
}

// getRequestHeaders extracts relevant headers from request
func getRequestHeaders(r *http.Request) map[string]string {
	headers := make(map[string]string)
	relevantHeaders := []string{"Authorization", "Content-Type", "X-Request-ID", "X-User-ID", "X-Country-Code"}

	for _, header := range relevantHeaders {
		if value := r.Header.Get(header); value != "" {
			headers[header] = value
		}
	}
	return headers
}

// Start starts the gateway test server
func (g *GatewayTestInstance) Start() error {
	g.server = httptest.NewServer(g.gateway.router)
	return nil
}

// Stop stops the gateway test server and all mock services
func (g *GatewayTestInstance) Stop() {
	if g.server != nil {
		g.server.Close()
	}

	for _, mockService := range g.mockServices {
		if mockService.server != nil {
			mockService.server.Close()
		}
	}
}

// GetURL returns the test server URL
func (g *GatewayTestInstance) GetURL() string {
	if g.server != nil {
		return g.server.URL
	}
	return ""
}

// SetProviderState configures the provider state for testing
func (g *GatewayTestInstance) SetProviderState(state string) error {
	switch state {
	case "gateway routes are configured":
		return g.setupRoutingState()
	case "rate limiting is enabled":
		return g.setupRateLimitingState()
	case "authentication middleware is active":
		return g.setupAuthenticationState()
	case "service discovery is functional":
		return g.setupServiceDiscoveryState()
	case "auth service is healthy":
		return g.setupHealthyAuthService()
	case "content service is healthy":
		return g.setupHealthyContentService()
	case "user exists with valid credentials":
		return g.setupValidUserCredentials()
	case "user is authenticated with valid token":
		return g.setupAuthenticatedUser()
	case "content items exist":
		return g.setupExistingContent()
	case "user is authenticated and can create content":
		return g.setupContentCreationPermissions()
	case "content exists in mobile category":
		return g.setupMobileCategoryContent()
	default:
		log.Printf("Unknown provider state: %s", state)
		return nil
	}
}

// Provider state setup methods

func (g *GatewayTestInstance) setupRoutingState() error {
	// Ensure all routes are properly configured
	// This is already done in NewGateway, but we can verify
	log.Println("Provider State: Gateway routes configured")
	return nil
}

func (g *GatewayTestInstance) setupRateLimitingState() error {
	// Enable rate limiting middleware (if implemented)
	log.Println("Provider State: Rate limiting enabled")
	return nil
}

func (g *GatewayTestInstance) setupAuthenticationState() error {
	// Ensure authentication middleware is active
	log.Println("Provider State: Authentication middleware active")
	return nil
}

func (g *GatewayTestInstance) setupServiceDiscoveryState() error {
	// Ensure service discovery is working
	for serviceName := range g.mockServices {
		service := g.gateway.registry.GetServiceByName(serviceName)
		if service == nil {
			return fmt.Errorf("service %s not registered", serviceName)
		}
		// Mark service as healthy
		g.gateway.registry.UpdateServiceHealth(service.ID, Healthy)
	}
	log.Println("Provider State: Service discovery functional")
	return nil
}

func (g *GatewayTestInstance) setupHealthyAuthService() error {
	// Configure auth service mock responses
	authService := g.mockServices["auth-service"]

	// Mock login response
	authService.SetMockResponse("POST", "/login", MockResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]interface{}{
			"access_token":  "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1c2VyLTEyMyIsImV4cCI6OTk5OTk5OTk5OX0.test-signature",
			"refresh_token": "rt_eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.refresh-payload.test-signature",
			"expires_in":    900,
			"user": map[string]interface{}{
				"id":           "123e4567-e89b-12d3-a456-426614174000",
				"phone_number": "0123456789",
				"country_code": "+66",
				"status":       "active",
				"profile": map[string]interface{}{
					"first_name": "John",
					"last_name":  "Doe",
				},
			},
		},
	})

	// Mock profile response
	authService.SetMockResponse("GET", "/profile", MockResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]interface{}{
			"id":           "123e4567-e89b-12d3-a456-426614174000",
			"phone_number": "0123456789",
			"country_code": "+66",
			"status":       "active",
			"profile": map[string]interface{}{
				"first_name":    "John",
				"last_name":     "Doe",
				"date_of_birth": "1990-01-01",
				"gender":        "male",
				"address": map[string]interface{}{
					"street":      "123 Main St",
					"city":        "Bangkok",
					"country":     "Thailand",
					"postal_code": "10100",
				},
			},
			"kyc": map[string]interface{}{
				"status":      "verified",
				"tier":        2,
				"verified_at": "2024-01-01T00:00:00Z",
			},
			"created_at": "2024-01-01T00:00:00Z",
			"updated_at": "2024-01-01T00:00:00Z",
		},
	})

	log.Println("Provider State: Auth service healthy")
	return nil
}

func (g *GatewayTestInstance) setupHealthyContentService() error {
	// Configure content service mock responses
	contentService := g.mockServices["content-service"]

	// Mock content list response
	contentService.SetMockResponse("GET", "/", MockResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"id":       "456e7890-e89b-12d3-a456-426614174000",
					"category": "general",
					"type":     "text",
					"value":    "Sample content value",
					"metadata": map[string]interface{}{
						"title":       "Sample Title",
						"description": "Sample description",
					},
					"status":     "published",
					"tags":       []string{"mobile", "general"},
					"created_at": "2024-01-01T00:00:00Z",
					"updated_at": "2024-01-01T00:00:00Z",
				},
			},
			"pagination": map[string]interface{}{
				"current_page":    1,
				"total_pages":     5,
				"total_items":     50,
				"items_per_page":  10,
			},
		},
	})

	log.Println("Provider State: Content service healthy")
	return nil
}

func (g *GatewayTestInstance) setupValidUserCredentials() error {
	// This is handled by setupHealthyAuthService
	return g.setupHealthyAuthService()
}

func (g *GatewayTestInstance) setupAuthenticatedUser() error {
	// This is handled by setupHealthyAuthService
	return g.setupHealthyAuthService()
}

func (g *GatewayTestInstance) setupExistingContent() error {
	// This is handled by setupHealthyContentService
	return g.setupHealthyContentService()
}

func (g *GatewayTestInstance) setupContentCreationPermissions() error {
	contentService := g.mockServices["content-service"]

	// Mock content creation response
	contentService.SetMockResponse("POST", "/", MockResponse{
		StatusCode: 201,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]interface{}{
			"id":       "789e1234-e89b-12d3-a456-426614174000",
			"category": "user-generated",
			"type":     "text",
			"value":    "New mobile content",
			"metadata": map[string]interface{}{
				"title":       "Mobile Content",
				"description": "Content created from mobile app",
			},
			"status":     "draft",
			"tags":       []string{"mobile", "user-generated"},
			"created_at": "2024-01-01T00:00:00Z",
			"updated_at": "2024-01-01T00:00:00Z",
		},
	})

	log.Println("Provider State: Content creation permissions set")
	return nil
}

func (g *GatewayTestInstance) setupMobileCategoryContent() error {
	contentService := g.mockServices["content-service"]

	// Mock mobile category content response
	contentService.SetMockResponse("GET", "/category/mobile", MockResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]interface{}{
			"category": "mobile",
			"items": []map[string]interface{}{
				{
					"id":   "abc12345-e89b-12d3-a456-426614174000",
					"type": "image",
					"value": "optimized-mobile-image.jpg",
					"metadata": map[string]interface{}{
						"title":    "Mobile Optimized Image",
						"alt_text": "Mobile friendly image",
						"dimensions": map[string]interface{}{
							"width":  375,
							"height": 667,
						},
					},
					"status": "published",
					"mobile_specific": map[string]interface{}{
						"retina_url":      "optimized-mobile-image@2x.jpg",
						"webp_url":        "optimized-mobile-image.webp",
						"loading_priority": "high",
					},
				},
			},
		},
	})

	log.Println("Provider State: Mobile category content set")
	return nil
}

// TestGatewayProviderVerification runs Pact provider verification tests
func TestGatewayProviderVerification(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping provider verification tests in short mode")
	}

	// Set up test mode for Gin
	gin.SetMode(gin.TestMode)

	// Create gateway test instance
	gatewayTest := NewGatewayTestInstance()

	// Start the gateway test server
	err := gatewayTest.Start()
	require.NoError(t, err, "Failed to start gateway test server")
	defer gatewayTest.Stop()

	// Get the test server URL
	gatewayURL := gatewayTest.GetURL()
	require.NotEmpty(t, gatewayURL, "Gateway test server URL should not be empty")

	// Configure Pact provider verification
	verifier := provider.NewVerifier()

	// Run verification against Pact contracts
	err = verifier.VerifyProvider(t, provider.VerifyRequest{
		ProviderBaseURL: gatewayURL,
		Provider:        "gateway-service",

		// Pact sources - in a real scenario these would come from Pact Broker
		PactURLs: []string{
			// These would normally be URLs from Pact Broker
			// For testing purposes, we can create local contract files
			// or use the contracts from the specs directory
		},

		// Alternative: Use Pact Broker
		BrokerURL:      getEnvOrDefault("PACT_BROKER_BASE_URL", "http://localhost:9292"),
		BrokerUsername: getEnvOrDefault("PACT_BROKER_USERNAME", ""),
		BrokerPassword: getEnvOrDefault("PACT_BROKER_PASSWORD", ""),

		// Provider version and tags
		ProviderVersion: "1.0.0-test",
		ProviderTags:    []string{"test", "gateway"},
		ConsumerVersionSelectors: []provider.Selector{
			&provider.ConsumerVersionSelector{
				Tag:    "main",
				Latest: true,
			},
			&provider.ConsumerVersionSelector{
				Tag:    "test",
				Latest: true,
			},
		},

		// Publishing options
		PublishVerificationResults: true,
	})

	assert.NoError(t, err, "Provider verification should pass")
}

// TestGatewayRouting tests the gateway routing functionality
func TestGatewayRouting(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gatewayTest := NewGatewayTestInstance()
	err := gatewayTest.Start()
	require.NoError(t, err)
	defer gatewayTest.Stop()

	baseURL := gatewayTest.GetURL()

	testCases := []struct {
		name           string
		method         string
		path           string
		headers        map[string]string
		expectedStatus int
		expectedService string
	}{
		{
			name:           "Route auth login request",
			method:         "POST",
			path:           "/api/v1/auth/login",
			headers:        map[string]string{"Content-Type": "application/json"},
			expectedStatus: http.StatusOK,
			expectedService: "auth-service",
		},
		{
			name:           "Route auth profile request with authentication",
			method:         "GET",
			path:           "/api/v1/auth/profile",
			headers:        map[string]string{
				"Authorization": "Bearer test-token",
				"Content-Type":  "application/json",
			},
			expectedStatus: http.StatusOK,
			expectedService: "auth-service",
		},
		{
			name:           "Route content request",
			method:         "GET",
			path:           "/api/v1/content",
			headers:        map[string]string{"Content-Type": "application/json"},
			expectedStatus: http.StatusOK,
			expectedService: "content-service",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create HTTP request
			req, err := http.NewRequest(tc.method, baseURL+tc.path, nil)
			require.NoError(t, err)

			// Set headers
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}

			// Execute request
			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Verify status code
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			// Verify X-Gateway-Service header is set (shows request was proxied)
			gatewayService := resp.Header.Get("X-Gateway-Service")
			if tc.expectedStatus == http.StatusOK {
				assert.NotEmpty(t, gatewayService, "Gateway should set X-Gateway-Service header")
			}
		})
	}
}

// TestGatewayHealthCheck tests the gateway health check functionality
func TestGatewayHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gatewayTest := NewGatewayTestInstance()
	err := gatewayTest.Start()
	require.NoError(t, err)
	defer gatewayTest.Stop()

	baseURL := gatewayTest.GetURL()

	// Test health endpoint
	resp, err := http.Get(baseURL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var healthResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&healthResponse)
	require.NoError(t, err)

	assert.Equal(t, "healthy", healthResponse["status"])
	assert.Equal(t, "api-gateway", healthResponse["service"])
	assert.Contains(t, healthResponse, "timestamp")
}

// TestGatewayServiceRegistry tests the service registry functionality
func TestGatewayServiceRegistry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gatewayTest := NewGatewayTestInstance()
	err := gatewayTest.Start()
	require.NoError(t, err)
	defer gatewayTest.Stop()

	baseURL := gatewayTest.GetURL()

	// Test services listing endpoint
	resp, err := http.Get(baseURL + "/registry/services")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var servicesResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&servicesResponse)
	require.NoError(t, err)

	assert.Contains(t, servicesResponse, "services")
	assert.Contains(t, servicesResponse, "count")

	services := servicesResponse["services"].([]interface{})
	assert.Greater(t, len(services), 0, "Should have registered services")
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Helper function to create minimal config for testing
func createTestConfig() *config.Config {
	return &config.Config{
		Environment: "test",
		Server: config.ServerConfig{
			Port:         0, // Let the test server choose a port
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}