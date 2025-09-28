package contract

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGatewayFunctionalityWithoutPact tests the Gateway functionality without requiring Pact FFI libraries
// This test validates the core Gateway behavior and can run in environments where Pact is not fully set up
func TestGatewayFunctionalityWithoutPact(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create and start the gateway test instance
	gatewayTest := NewGatewayTestInstance()
	err := gatewayTest.Start()
	require.NoError(t, err, "Failed to start gateway test server")
	defer gatewayTest.Stop()

	baseURL := gatewayTest.GetURL()
	require.NotEmpty(t, baseURL, "Gateway test server URL should not be empty")

	// Test 1: Gateway Health Check
	t.Run("Gateway Health Check", func(t *testing.T) {
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
	})

	// Test 2: Service Registry
	t.Run("Service Registry", func(t *testing.T) {
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

		// Validate service structure
		if len(services) > 0 {
			service := services[0].(map[string]interface{})
			assert.Contains(t, service, "id")
			assert.Contains(t, service, "name")
			assert.Contains(t, service, "host")
			assert.Contains(t, service, "port")
			assert.Contains(t, service, "health")
		}
	})

	// Test 3: Auth Service Routing
	t.Run("Auth Service Routing", func(t *testing.T) {
		// Setup auth service mock response
		authService := gatewayTest.mockServices["auth-service"]
		authService.SetMockResponse("POST", "/login", MockResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: map[string]interface{}{
				"access_token":  "test-jwt-token",
				"refresh_token": "test-refresh-token",
				"expires_in":    900,
				"user": map[string]interface{}{
					"id":           "test-user-id",
					"phone_number": "0123456789",
					"country_code": "+66",
					"status":       "active",
				},
			},
		})

		// Make request to gateway auth endpoint
		resp, err := http.Post(baseURL+"/api/v1/auth/login", "application/json", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var loginResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&loginResponse)
		require.NoError(t, err)

		// Validate response structure matches contract expectations
		assert.Contains(t, loginResponse, "access_token")
		assert.Contains(t, loginResponse, "user")

		user := loginResponse["user"].(map[string]interface{})
		assert.Contains(t, user, "id")
		assert.Contains(t, user, "phone_number")
		assert.Contains(t, user, "country_code")
	})

	// Test 4: Content Service Routing
	t.Run("Content Service Routing", func(t *testing.T) {
		// Setup content service mock response
		contentService := gatewayTest.mockServices["content-service"]
		contentService.SetMockResponse("GET", "/", MockResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"id":       "test-content-id",
						"category": "general",
						"type":     "text",
						"value":    "Test content value",
						"status":   "published",
						"tags":     []string{"test"},
					},
				},
				"pagination": map[string]interface{}{
					"current_page":    1,
					"total_pages":     1,
					"total_items":     1,
					"items_per_page":  10,
				},
			},
		})

		// Make request to gateway content endpoint
		resp, err := http.Get(baseURL + "/api/v1/content/")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var contentResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&contentResponse)
		require.NoError(t, err)

		// Validate response structure matches contract expectations
		assert.Contains(t, contentResponse, "items")
		assert.Contains(t, contentResponse, "pagination")

		items := contentResponse["items"].([]interface{})
		assert.Len(t, items, 1)

		item := items[0].(map[string]interface{})
		assert.Contains(t, item, "id")
		assert.Contains(t, item, "category")
		assert.Contains(t, item, "type")
		assert.Contains(t, item, "value")
	})

	// Test 5: Service Unavailable Handling
	t.Run("Service Unavailable Handling", func(t *testing.T) {
		// Make request to a service that doesn't exist
		resp, err := http.Get(baseURL + "/api/v1/nonexistent/test")
		require.NoError(t, err)
		defer resp.Body.Close()

		// Gateway should return service unavailable
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)

		assert.Contains(t, errorResponse, "error")
		assert.Contains(t, errorResponse, "message")
		assert.Equal(t, "service_unavailable", errorResponse["error"])
	})

	// Test 6: Provider State Configuration
	t.Run("Provider State Configuration", func(t *testing.T) {
		// Test various provider states
		states := []string{
			"gateway routes are configured",
			"service discovery is functional",
			"auth service is healthy",
			"content service is healthy",
		}

		for _, state := range states {
			t.Run("State: "+state, func(t *testing.T) {
				err := gatewayTest.SetProviderState(state)
				assert.NoError(t, err, "Provider state setup should succeed for: %s", state)
			})
		}
	})

	// Test 7: Header Forwarding
	t.Run("Header Forwarding", func(t *testing.T) {
		// Setup mock to capture forwarded headers
		authService := gatewayTest.mockServices["auth-service"]
		authService.SetMockResponse("GET", "/profile", MockResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: map[string]interface{}{
				"id":       "test-user",
				"profile":  map[string]interface{}{},
				"forwarded_headers": "captured", // Mock would capture actual headers
			},
		})

		// Create request with auth headers
		req, err := http.NewRequest("GET", baseURL+"/api/v1/auth/profile", nil)
		require.NoError(t, err)

		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("X-Request-ID", "test-request-123")

		// Execute request
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify that gateway processed the request (exact header validation would require mock inspection)
		var profileResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&profileResponse)
		require.NoError(t, err)

		assert.Contains(t, profileResponse, "id")
	})
}

// TestProviderStateValidation tests that all expected provider states are properly handled
func TestProviderStateValidation(t *testing.T) {
	gatewayTest := NewGatewayTestInstance()

	providerStates := map[string]bool{
		"gateway routes are configured":               true,
		"rate limiting is enabled":                    true,
		"authentication middleware is active":         true,
		"service discovery is functional":             true,
		"auth service is healthy":                     true,
		"content service is healthy":                  true,
		"user exists with valid credentials":          true,
		"user is authenticated with valid token":      true,
		"content items exist":                         true,
		"user is authenticated and can create content": true,
		"content exists in mobile category":           true,
		"unknown state":                               true, // Should handle gracefully
	}

	for state, shouldSucceed := range providerStates {
		t.Run("ProviderState_"+state, func(t *testing.T) {
			err := gatewayTest.SetProviderState(state)
			if shouldSucceed {
				assert.NoError(t, err, "Provider state should be handled: %s", state)
			} else {
				assert.Error(t, err, "Provider state should fail: %s", state)
			}
		})
	}
}

// TestGatewayServiceDiscovery tests the service discovery and load balancing functionality
func TestGatewayServiceDiscovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gatewayTest := NewGatewayTestInstance()
	err := gatewayTest.Start()
	require.NoError(t, err)
	defer gatewayTest.Stop()

	// Test service registry functionality
	registry := gatewayTest.gateway.registry

	// Verify services are registered
	services := registry.GetServicesByName("auth-service")
	assert.Greater(t, len(services), 0, "Auth service should be registered")

	// Test healthy service selection
	healthyService := registry.GetHealthyService("auth-service")
	assert.NotNil(t, healthyService, "Should find healthy auth service")
	assert.Equal(t, "auth-service", healthyService.Name)

	// Test service health updates
	if healthyService != nil {
		registry.UpdateServiceHealth(healthyService.ID, Unhealthy)
		unhealthyService := registry.GetHealthyService("auth-service")
		assert.Nil(t, unhealthyService, "Should not find healthy service after marking unhealthy")

		// Restore health
		registry.UpdateServiceHealth(healthyService.ID, Healthy)
		restoredService := registry.GetHealthyService("auth-service")
		assert.NotNil(t, restoredService, "Should find healthy service after restoration")
	}
}

// TestContractExpectationCompliance tests that our gateway responses match expected contract structures
func TestContractExpectationCompliance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gatewayTest := NewGatewayTestInstance()
	err := gatewayTest.Start()
	require.NoError(t, err)
	defer gatewayTest.Stop()

	baseURL := gatewayTest.GetURL()

	// Setup all provider states for comprehensive testing
	states := []string{
		"gateway routes are configured",
		"auth service is healthy",
		"user exists with valid credentials",
		"content service is healthy",
		"content items exist",
	}

	for _, state := range states {
		err := gatewayTest.SetProviderState(state)
		require.NoError(t, err, "Failed to setup provider state: %s", state)
	}

	// Test contract compliance for different endpoints
	testCases := []struct {
		name           string
		method         string
		endpoint       string
		expectedStatus int
		expectedFields []string
	}{
		{
			name:           "Auth Login Contract",
			method:         "POST",
			endpoint:       "/api/v1/auth/login",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"access_token", "user"},
		},
		{
			name:           "Auth Profile Contract",
			method:         "GET",
			endpoint:       "/api/v1/auth/profile",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"id", "phone_number", "country_code", "profile"},
		},
		{
			name:           "Content List Contract",
			method:         "GET",
			endpoint:       "/api/v1/content/",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"items", "pagination"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, baseURL+tc.endpoint, nil)
			require.NoError(t, err)

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if resp.StatusCode == http.StatusOK {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)

				for _, field := range tc.expectedFields {
					assert.Contains(t, response, field, "Response should contain field: %s", field)
				}
			}
		})
	}
}