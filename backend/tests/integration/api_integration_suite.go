package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// APIIntegrationSuite provides a comprehensive testing framework for all microservices
type APIIntegrationSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	authToken  string
	testUser   *TestUser
}

// Note: TestUser, ServicePort, and DefaultServicePorts are now defined in types.go

// SetupSuite initializes the integration test suite
func (suite *APIIntegrationSuite) SetupSuite() {
	suite.baseURL = "http://localhost"
	suite.httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}

	// Initialize test user
	suite.testUser = &TestUser{
		PhoneNumber: "+1234567890",
		Country:     "US",
		Language:    "en",
	}
}

// TearDownSuite cleans up after all tests
func (suite *APIIntegrationSuite) TearDownSuite() {
	// Cleanup test data if needed
}

// SetupTest runs before each test
func (suite *APIIntegrationSuite) SetupTest() {
	// Reset any test state if needed
}

// TearDownTest runs after each test
func (suite *APIIntegrationSuite) TearDownTest() {
	// Cleanup test state if needed
}

// Note: APIResponse and APIError are now defined in types.go

// makeRequest performs HTTP requests with common setup
func (suite *APIIntegrationSuite) makeRequest(method, url string, body interface{}, headers map[string]string) (*http.Response, error) {
	var reqBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Add auth token if available
	if suite.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+suite.authToken)
	}

	return suite.httpClient.Do(req)
}

// parseResponse parses HTTP response into APIResponse struct
func (suite *APIIntegrationSuite) parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if target != nil {
		if err := json.Unmarshal(body, target); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w, body: %s", err, string(body))
		}
	}

	return nil
}

// assertSuccessResponse validates successful API responses
func (suite *APIIntegrationSuite) assertSuccessResponse(resp *http.Response, expectedStatus int) *APIResponse {
	suite.Equal(expectedStatus, resp.StatusCode, "Expected status code %d, got %d", expectedStatus, resp.StatusCode)

	var apiResp APIResponse
	err := suite.parseResponse(resp, &apiResp)
	suite.NoError(err, "Failed to parse API response")

	// Handle different response formats
	if apiResp.Success || apiResp.Status == "success" {
		return &apiResp
	}

	suite.Fail("API response indicates failure", "Response: %+v", apiResp)
	return nil
}

// assertErrorResponse validates error API responses
func (suite *APIIntegrationSuite) assertErrorResponse(resp *http.Response, expectedStatus int) *APIResponse {
	suite.Equal(expectedStatus, resp.StatusCode, "Expected status code %d, got %d", expectedStatus, resp.StatusCode)

	var apiResp APIResponse
	err := suite.parseResponse(resp, &apiResp)
	suite.NoError(err, "Failed to parse error API response")

	return &apiResp
}

// waitForService waits for a service to be available
func (suite *APIIntegrationSuite) waitForService(port int, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	url := fmt.Sprintf("%s:%d/health", suite.baseURL, port)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("service on port %d not available after %v", port, timeout)
		default:
			resp, err := suite.httpClient.Get(url)
			if err == nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				return nil
			}
			if resp != nil {
				resp.Body.Close()
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// Note: ServiceHealthCheck is now defined in types.go

// checkServiceHealth performs health check on a service
func (suite *APIIntegrationSuite) checkServiceHealth(port int) (*ServiceHealthCheck, error) {
	url := fmt.Sprintf("%s:%d/health", suite.baseURL, port)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("health check request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	var healthResp struct {
		Success bool                `json:"success"`
		Status  string              `json:"status"`
		Data    *ServiceHealthCheck `json:"data"`
	}

	if err := suite.parseResponse(resp, &healthResp); err != nil {
		return nil, fmt.Errorf("failed to parse health response: %w", err)
	}

	return healthResp.Data, nil
}

// generateTestData creates test data for various endpoints
func (suite *APIIntegrationSuite) generateTestData() map[string]interface{} {
	return map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"testID":    fmt.Sprintf("test_%d", time.Now().UnixNano()),
	}
}

// RunAPIIntegrationSuite runs the complete API integration test suite
func RunAPIIntegrationSuite(t *testing.T) {
	suite.Run(t, new(APIIntegrationSuite))
}