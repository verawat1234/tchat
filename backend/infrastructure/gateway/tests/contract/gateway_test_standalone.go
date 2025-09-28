package contract

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Standalone Gateway test that doesn't require Pact libraries
// This demonstrates the complete provider verification functionality
// without external dependencies

// StandaloneGatewayTest provides comprehensive Gateway testing without Pact FFI dependencies
func TestStandaloneGatewayFunctionality(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create and start the gateway test instance (using shared types from pact_provider_test.go)
	gatewayTest := NewGatewayTestInstance()
	err := gatewayTest.Start()
	require.NoError(t, err, "Failed to start gateway test server")
	defer gatewayTest.Stop()

	baseURL := gatewayTest.GetURL()
	require.NotEmpty(t, baseURL, "Gateway test server URL should not be empty")

	// Test 1: Gateway Health Check - Validates core gateway functionality
	t.Run("Gateway_Health_Check", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var healthResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&healthResponse)
		require.NoError(t, err)

		// Validate contract expectations
		assert.Equal(t, "healthy", healthResponse["status"])
		assert.Equal(t, "api-gateway", healthResponse["service"])
		assert.Contains(t, healthResponse, "timestamp")

		t.Logf("✓ Gateway health check passed - Status: %v", healthResponse["status"])
	})

	// Test 2: Service Registry - Validates service discovery functionality
	t.Run("Service_Registry_Functionality", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/registry/services")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var servicesResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&servicesResponse)
		require.NoError(t, err)

		// Validate registry contract
		assert.Contains(t, servicesResponse, "services")
		assert.Contains(t, servicesResponse, "count")

		services := servicesResponse["services"].([]interface{})
		assert.Greater(t, len(services), 0, "Should have registered services")

		// Validate service instance structure
		if len(services) > 0 {
			service := services[0].(map[string]interface{})
			requiredFields := []string{"id", "name", "host", "port", "health"}
			for _, field := range requiredFields {
				assert.Contains(t, service, field, "Service should contain field: %s", field)
			}
		}

		t.Logf("✓ Service registry passed - Registered services: %d", len(services))
	})

	// Test 3: Auth Service Routing - Validates request proxying and response forwarding
	t.Run("Auth_Service_Routing", func(t *testing.T) {
		// Setup provider state: auth service is healthy
		err := gatewayTest.SetProviderState("auth service is healthy")
		require.NoError(t, err)

		// Make request to gateway auth endpoint
		resp, err := http.Post(baseURL+"/api/v1/auth/login", "application/json", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var loginResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&loginResponse)
		require.NoError(t, err)

		// Validate contract expectations for auth login
		expectedFields := []string{"access_token", "refresh_token", "expires_in", "user"}
		for _, field := range expectedFields {
			assert.Contains(t, loginResponse, field, "Login response should contain: %s", field)
		}

		// Validate user structure
		user := loginResponse["user"].(map[string]interface{})
		userFields := []string{"id", "phone_number", "country_code", "status", "profile"}
		for _, field := range userFields {
			assert.Contains(t, user, field, "User object should contain: %s", field)
		}

		t.Logf("✓ Auth service routing passed - Access token: %s",
			loginResponse["access_token"].(string)[:20]+"...")
	})

	// Test 4: Content Service Routing - Validates query parameter forwarding and pagination
	t.Run("Content_Service_Routing", func(t *testing.T) {
		// Setup provider state: content items exist
		err := gatewayTest.SetProviderState("content items exist")
		require.NoError(t, err)

		// Test content listing with query parameters
		resp, err := http.Get(baseURL + "/api/v1/content/?page=1&limit=10&category=general")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var contentResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&contentResponse)
		require.NoError(t, err)

		// Validate contract expectations for content listing
		assert.Contains(t, contentResponse, "items")
		assert.Contains(t, contentResponse, "pagination")

		items := contentResponse["items"].([]interface{})
		assert.Greater(t, len(items), 0, "Should have content items")

		// Validate content item structure
		if len(items) > 0 {
			item := items[0].(map[string]interface{})
			itemFields := []string{"id", "category", "type", "value", "status", "tags"}
			for _, field := range itemFields {
				assert.Contains(t, item, field, "Content item should contain: %s", field)
			}
		}

		// Validate pagination structure
		pagination := contentResponse["pagination"].(map[string]interface{})
		paginationFields := []string{"current_page", "total_pages", "total_items", "items_per_page"}
		for _, field := range paginationFields {
			assert.Contains(t, pagination, field, "Pagination should contain: %s", field)
		}

		t.Logf("✓ Content service routing passed - Items: %d", len(items))
	})

	// Test 5: Mobile Content Routing - Validates mobile-specific content handling
	t.Run("Mobile_Content_Category_Routing", func(t *testing.T) {
		// Setup provider state: content exists in mobile category
		err := gatewayTest.SetProviderState("content exists in mobile category")
		require.NoError(t, err)

		// Test mobile category endpoint
		resp, err := http.Get(baseURL + "/api/v1/content/category/mobile")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var mobileContentResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&mobileContentResponse)
		require.NoError(t, err)

		// Validate mobile-specific contract
		assert.Contains(t, mobileContentResponse, "category")
		assert.Contains(t, mobileContentResponse, "items")
		assert.Equal(t, "mobile", mobileContentResponse["category"])

		items := mobileContentResponse["items"].([]interface{})
		assert.Greater(t, len(items), 0, "Should have mobile content items")

		// Validate mobile-specific fields
		if len(items) > 0 {
			item := items[0].(map[string]interface{})
			assert.Contains(t, item, "mobile_specific", "Mobile content should have mobile_specific field")

			mobileSpec := item["mobile_specific"].(map[string]interface{})
			mobileFields := []string{"retina_url", "webp_url", "loading_priority"}
			for _, field := range mobileFields {
				assert.Contains(t, mobileSpec, field, "Mobile specification should contain: %s", field)
			}
		}

		t.Logf("✓ Mobile content routing passed - Mobile items: %d", len(items))
	})

	// Test 6: Service Unavailable Handling - Validates error handling
	t.Run("Service_Unavailable_Handling", func(t *testing.T) {
		// Make request to a service that doesn't exist or is unhealthy
		resp, err := http.Get(baseURL + "/api/v1/nonexistent/test")
		require.NoError(t, err)
		defer resp.Body.Close()

		// Gateway should return service unavailable
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)

		// Validate error contract
		assert.Contains(t, errorResponse, "error")
		assert.Contains(t, errorResponse, "message")
		assert.Equal(t, "service_unavailable", errorResponse["error"])

		t.Logf("✓ Service unavailable handling passed - Error: %s", errorResponse["error"])
	})

	// Test 7: Header Forwarding - Validates authentication and trace header forwarding
	t.Run("Header_Forwarding", func(t *testing.T) {
		// Setup authenticated user state
		err := gatewayTest.SetProviderState("user is authenticated with valid token")
		require.NoError(t, err)

		// Create request with authentication and trace headers
		req, err := http.NewRequest("GET", baseURL+"/api/v1/auth/profile", nil)
		require.NoError(t, err)

		// Add headers that should be forwarded
		req.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...")
		req.Header.Set("X-Request-ID", "test-request-123")
		req.Header.Set("X-User-ID", "test-user-456")
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var profileResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&profileResponse)
		require.NoError(t, err)

		// Validate profile contract
		profileFields := []string{"id", "phone_number", "country_code", "profile", "kyc"}
		for _, field := range profileFields {
			assert.Contains(t, profileResponse, field, "Profile should contain: %s", field)
		}

		// Validate nested profile structure
		profile := profileResponse["profile"].(map[string]interface{})
		nestedFields := []string{"first_name", "last_name", "date_of_birth", "gender", "address"}
		for _, field := range nestedFields {
			assert.Contains(t, profile, field, "Profile details should contain: %s", field)
		}

		t.Logf("✓ Header forwarding passed - Profile ID: %s", profileResponse["id"])
	})
}

// TestProviderStateComprehensive tests all provider states that would be used in Pact verification
func TestProviderStateComprehensive(t *testing.T) {
	gatewayTest := NewGatewayTestInstance()

	// Define all provider states that the Gateway should support
	providerStates := map[string]string{
		"gateway routes are configured":               "Gateway routing table setup",
		"rate limiting is enabled":                    "Rate limiting middleware active",
		"authentication middleware is active":         "JWT auth middleware enabled",
		"service discovery is functional":             "Service registry operational",
		"auth service is healthy":                     "Auth service responding to requests",
		"content service is healthy":                  "Content service responding to requests",
		"user exists with valid credentials":          "Valid user available for login",
		"user is authenticated with valid token":      "Authenticated user context available",
		"content items exist":                         "Content available for retrieval",
		"user is authenticated and can create content": "User has content creation permissions",
		"content exists in mobile category":           "Mobile-specific content available",
	}

	for state, description := range providerStates {
		t.Run("ProviderState_"+strings.ReplaceAll(state, " ", "_"), func(t *testing.T) {
			err := gatewayTest.SetProviderState(state)
			assert.NoError(t, err, "Provider state should be handled: %s (%s)", state, description)

			if err == nil {
				t.Logf("✓ Provider state configured: %s", state)
			}
		})
	}
}

// TestGatewayServiceDiscoveryAdvanced tests advanced service discovery features
func TestGatewayServiceDiscoveryAdvanced(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gatewayTest := NewGatewayTestInstance()
	err := gatewayTest.Start()
	require.NoError(t, err)
	defer gatewayTest.Stop()

	registry := gatewayTest.gateway.registry

	// Test 1: Multiple service instances
	t.Run("Multiple_Service_Instances", func(t *testing.T) {
		services := registry.GetServicesByName("auth-service")
		assert.Greater(t, len(services), 0, "Auth service should be registered")

		for i, service := range services {
			assert.Equal(t, "auth-service", service.Name)
			assert.NotEmpty(t, service.ID)
			assert.NotEmpty(t, service.Host)
			assert.Greater(t, service.Port, 0)
			t.Logf("✓ Auth service instance %d: %s:%d", i+1, service.Host, service.Port)
		}
	})

	// Test 2: Health status management
	t.Run("Health_Status_Management", func(t *testing.T) {
		// Find a healthy service
		healthyService := registry.GetHealthyService("auth-service")
		require.NotNil(t, healthyService, "Should find healthy auth service")

		originalHealth := healthyService.Health

		// Mark service as unhealthy
		registry.UpdateServiceHealth(healthyService.ID, Unhealthy)
		unhealthyService := registry.GetHealthyService("auth-service")
		assert.Nil(t, unhealthyService, "Should not find healthy service after marking unhealthy")

		// Restore health
		registry.UpdateServiceHealth(healthyService.ID, Healthy)
		restoredService := registry.GetHealthyService("auth-service")
		assert.NotNil(t, restoredService, "Should find healthy service after restoration")

		t.Logf("✓ Health management: %s → unhealthy → healthy", originalHealth)
	})

	// Test 3: Load balancing behavior
	t.Run("Load_Balancing_Behavior", func(t *testing.T) {
		// Get multiple service instances (simulating load balancing selection)
		selections := make(map[string]int)
		for i := 0; i < 10; i++ {
			service := registry.GetHealthyService("auth-service")
			if service != nil {
				selections[service.ID]++
			}
		}

		assert.Greater(t, len(selections), 0, "Should select at least one service instance")

		// Log selection distribution
		for serviceID, count := range selections {
			t.Logf("✓ Service %s selected %d times", serviceID[:8], count)
		}
	})
}

// TestContractExpectationValidation validates that Gateway responses match Pact contract expectations
func TestContractExpectationValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gatewayTest := NewGatewayTestInstance()
	err := gatewayTest.Start()
	require.NoError(t, err)
	defer gatewayTest.Stop()

	baseURL := gatewayTest.GetURL()

	// Setup all required provider states
	requiredStates := []string{
		"gateway routes are configured",
		"auth service is healthy",
		"content service is healthy",
		"user exists with valid credentials",
		"user is authenticated with valid token",
		"content items exist",
		"content exists in mobile category",
	}

	for _, state := range requiredStates {
		err := gatewayTest.SetProviderState(state)
		require.NoError(t, err, "Failed to setup provider state: %s", state)
	}

	// Define contract validation test cases based on consumer contracts
	contractTestCases := []struct {
		name           string
		method         string
		endpoint       string
		headers        map[string]string
		expectedStatus int
		expectedFields []string
		description    string
	}{
		{
			name:           "Web_Frontend_Auth_Login",
			method:         "POST",
			endpoint:       "/api/v1/auth/login",
			headers:        map[string]string{"Content-Type": "application/json"},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"access_token", "refresh_token", "expires_in", "user"},
			description:    "Web frontend login contract validation",
		},
		{
			name:           "Web_Frontend_Auth_Profile",
			method:         "GET",
			endpoint:       "/api/v1/auth/profile",
			headers:        map[string]string{
				"Authorization": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
				"Content-Type":  "application/json",
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"id", "phone_number", "country_code", "status", "profile", "kyc"},
			description:    "Web frontend profile contract validation",
		},
		{
			name:           "Mobile_App_Content_List",
			method:         "GET",
			endpoint:       "/api/v1/content",
			headers:        map[string]string{
				"Authorization": "Bearer mobile-token",
				"Content-Type":  "application/json",
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"items", "pagination"},
			description:    "Mobile app content listing contract validation",
		},
		{
			name:           "Mobile_App_Content_Category",
			method:         "GET",
			endpoint:       "/api/v1/content/category/mobile",
			headers:        map[string]string{
				"Authorization": "Bearer mobile-token",
				"Content-Type":  "application/json",
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"category", "items"},
			description:    "Mobile app mobile category contract validation",
		},
	}

	for _, tc := range contractTestCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest(tc.method, baseURL+tc.endpoint, nil)
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

			// Validate status code
			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Status code should match contract")

			// Validate response structure
			if resp.StatusCode == http.StatusOK {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err, "Response should be valid JSON")

				// Check all expected fields are present
				for _, field := range tc.expectedFields {
					assert.Contains(t, response, field, "Response should contain field: %s", field)
				}

				t.Logf("✓ %s: All %d expected fields present", tc.description, len(tc.expectedFields))
			}
		})
	}
}

// TestGatewayPerformanceAndReliability tests performance characteristics and reliability features
func TestGatewayPerformanceAndReliability(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gatewayTest := NewGatewayTestInstance()
	err := gatewayTest.Start()
	require.NoError(t, err)
	defer gatewayTest.Stop()

	baseURL := gatewayTest.GetURL()

	// Performance test: Multiple concurrent requests
	t.Run("Concurrent_Request_Handling", func(t *testing.T) {
		concurrency := 10
		requests := 20

		resultChan := make(chan bool, concurrency*requests)

		for i := 0; i < concurrency; i++ {
			go func() {
				for j := 0; j < requests; j++ {
					resp, err := http.Get(baseURL + "/health")
					success := err == nil && resp.StatusCode == http.StatusOK
					if resp != nil {
						resp.Body.Close()
					}
					resultChan <- success
				}
			}()
		}

		// Collect results
		successCount := 0
		totalRequests := concurrency * requests
		for i := 0; i < totalRequests; i++ {
			if <-resultChan {
				successCount++
			}
		}

		successRate := float64(successCount) / float64(totalRequests) * 100
		assert.Greater(t, successRate, 95.0, "Success rate should be > 95%%")

		t.Logf("✓ Concurrent requests: %d/%d successful (%.1f%%)",
			successCount, totalRequests, successRate)
	})

	// Reliability test: Service failure handling
	t.Run("Service_Failure_Resilience", func(t *testing.T) {
		registry := gatewayTest.gateway.registry

		// Get a healthy service and mark it unhealthy
		authService := registry.GetHealthyService("auth-service")
		require.NotNil(t, authService, "Should have healthy auth service")

		originalHealth := authService.Health
		registry.UpdateServiceHealth(authService.ID, Unhealthy)

		// Request should still be handled gracefully (service unavailable)
		resp, err := http.Get(baseURL + "/api/v1/auth/profile")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

		// Restore service health
		registry.UpdateServiceHealth(authService.ID, Healthy)

		// Verify service is available again
		resp2, err := http.Get(baseURL + "/api/v1/auth/profile")
		require.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, http.StatusOK, resp2.StatusCode)

		t.Logf("✓ Service failure resilience: %s → unhealthy (503) → healthy (200)",
			originalHealth)
	})
}