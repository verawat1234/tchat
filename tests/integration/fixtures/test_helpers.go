package fixtures

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// HTTPTestHelper provides utilities for HTTP testing
type HTTPTestHelper struct {
	Server  *httptest.Server
	Client  *http.Client
	BaseURL string
	Headers map[string]string
	t       *testing.T
}

// NewHTTPTestHelper creates a new HTTP test helper
func NewHTTPTestHelper(t *testing.T, handler http.Handler) *HTTPTestHelper {
	server := httptest.NewServer(handler)

	return &HTTPTestHelper{
		Server:  server,
		Client:  &http.Client{Timeout: 30 * time.Second},
		BaseURL: server.URL,
		Headers: make(map[string]string),
		t:       t,
	}
}

// SetHeader sets a header for all subsequent requests
func (h *HTTPTestHelper) SetHeader(key, value string) {
	h.Headers[key] = value
}

// SetAuthToken sets the Authorization header with Bearer token
func (h *HTTPTestHelper) SetAuthToken(token string) {
	h.SetHeader("Authorization", "Bearer "+token)
}

// GET performs a GET request and returns the response
func (h *HTTPTestHelper) GET(path string) *HTTPResponse {
	return h.Request("GET", path, nil)
}

// POST performs a POST request with JSON body
func (h *HTTPTestHelper) POST(path string, body interface{}) *HTTPResponse {
	return h.Request("POST", path, body)
}

// PUT performs a PUT request with JSON body
func (h *HTTPTestHelper) PUT(path string, body interface{}) *HTTPResponse {
	return h.Request("PUT", path, body)
}

// PATCH performs a PATCH request with JSON body
func (h *HTTPTestHelper) PATCH(path string, body interface{}) *HTTPResponse {
	return h.Request("PATCH", path, body)
}

// DELETE performs a DELETE request
func (h *HTTPTestHelper) DELETE(path string) *HTTPResponse {
	return h.Request("DELETE", path, nil)
}

// Request performs an HTTP request with the specified method, path, and body
func (h *HTTPTestHelper) Request(method, path string, body interface{}) *HTTPResponse {
	var reqBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(h.t, err, "Failed to marshal request body")
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, h.BaseURL+path, reqBody)
	require.NoError(h.t, err, "Failed to create HTTP request")

	// Set headers
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range h.Headers {
		req.Header.Set(key, value)
	}

	resp, err := h.Client.Do(req)
	require.NoError(h.t, err, "HTTP request failed")

	return &HTTPResponse{
		Response: resp,
		t:        h.t,
	}
}

// Close closes the test server
func (h *HTTPTestHelper) Close() {
	if h.Server != nil {
		h.Server.Close()
	}
}

// HTTPResponse wraps http.Response with testing utilities
type HTTPResponse struct {
	*http.Response
	t *testing.T
}

// AssertStatus asserts the response status code
func (r *HTTPResponse) AssertStatus(expectedStatus int) *HTTPResponse {
	assert.Equal(r.t, expectedStatus, r.StatusCode, "Unexpected status code")
	return r
}

// AssertOK asserts the response status is 200
func (r *HTTPResponse) AssertOK() *HTTPResponse {
	return r.AssertStatus(http.StatusOK)
}

// AssertCreated asserts the response status is 201
func (r *HTTPResponse) AssertCreated() *HTTPResponse {
	return r.AssertStatus(http.StatusCreated)
}

// AssertNoContent asserts the response status is 204
func (r *HTTPResponse) AssertNoContent() *HTTPResponse {
	return r.AssertStatus(http.StatusNoContent)
}

// AssertBadRequest asserts the response status is 400
func (r *HTTPResponse) AssertBadRequest() *HTTPResponse {
	return r.AssertStatus(http.StatusBadRequest)
}

// AssertUnauthorized asserts the response status is 401
func (r *HTTPResponse) AssertUnauthorized() *HTTPResponse {
	return r.AssertStatus(http.StatusUnauthorized)
}

// AssertForbidden asserts the response status is 403
func (r *HTTPResponse) AssertForbidden() *HTTPResponse {
	return r.AssertStatus(http.StatusForbidden)
}

// AssertNotFound asserts the response status is 404
func (r *HTTPResponse) AssertNotFound() *HTTPResponse {
	return r.AssertStatus(http.StatusNotFound)
}

// AssertInternalServerError asserts the response status is 500
func (r *HTTPResponse) AssertInternalServerError() *HTTPResponse {
	return r.AssertStatus(http.StatusInternalServerError)
}

// GetBodyBytes returns the response body as bytes
func (r *HTTPResponse) GetBodyBytes() []byte {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	require.NoError(r.t, err, "Failed to read response body")
	return body
}

// GetBodyString returns the response body as string
func (r *HTTPResponse) GetBodyString() string {
	return string(r.GetBodyBytes())
}

// UnmarshalJSON unmarshals the response body into the given interface
func (r *HTTPResponse) UnmarshalJSON(v interface{}) *HTTPResponse {
	body := r.GetBodyBytes()
	err := json.Unmarshal(body, v)
	require.NoError(r.t, err, "Failed to unmarshal response body")
	return r
}

// AssertJSONField asserts that a JSON field has the expected value
func (r *HTTPResponse) AssertJSONField(field string, expectedValue interface{}) *HTTPResponse {
	var responseData map[string]interface{}
	r.UnmarshalJSON(&responseData)

	actualValue, exists := responseData[field]
	require.True(r.t, exists, "Field %s not found in response", field)
	assert.Equal(r.t, expectedValue, actualValue, "Field %s has unexpected value", field)

	return r
}

// AssertJSONFieldExists asserts that a JSON field exists in the response
func (r *HTTPResponse) AssertJSONFieldExists(field string) *HTTPResponse {
	var responseData map[string]interface{}
	r.UnmarshalJSON(&responseData)

	_, exists := responseData[field]
	assert.True(r.t, exists, "Field %s should exist in response", field)

	return r
}

// AssertJSONFieldNotExists asserts that a JSON field does not exist in the response
func (r *HTTPResponse) AssertJSONFieldNotExists(field string) *HTTPResponse {
	var responseData map[string]interface{}
	r.UnmarshalJSON(&responseData)

	_, exists := responseData[field]
	assert.False(r.t, exists, "Field %s should not exist in response", field)

	return r
}

// AssertJSONArrayLength asserts the length of a JSON array field
func (r *HTTPResponse) AssertJSONArrayLength(field string, expectedLength int) *HTTPResponse {
	var responseData map[string]interface{}
	r.UnmarshalJSON(&responseData)

	value, exists := responseData[field]
	require.True(r.t, exists, "Field %s not found in response", field)

	array, ok := value.([]interface{})
	require.True(r.t, ok, "Field %s is not an array", field)

	assert.Equal(r.t, expectedLength, len(array), "Array %s has unexpected length", field)

	return r
}

// AssertHeader asserts that a response header has the expected value
func (r *HTTPResponse) AssertHeader(header, expectedValue string) *HTTPResponse {
	actualValue := r.Header.Get(header)
	assert.Equal(r.t, expectedValue, actualValue, "Header %s has unexpected value", header)
	return r
}

// AssertHeaderExists asserts that a response header exists
func (r *HTTPResponse) AssertHeaderExists(header string) *HTTPResponse {
	actualValue := r.Header.Get(header)
	assert.NotEmpty(r.t, actualValue, "Header %s should exist", header)
	return r
}

// AssertContentType asserts the response content type
func (r *HTTPResponse) AssertContentType(expectedContentType string) *HTTPResponse {
	return r.AssertHeader("Content-Type", expectedContentType)
}

// AssertJSONContentType asserts the response content type is JSON
func (r *HTTPResponse) AssertJSONContentType() *HTTPResponse {
	contentType := r.Header.Get("Content-Type")
	assert.Contains(r.t, contentType, "application/json", "Response should have JSON content type")
	return r
}

// AssertResponseTime asserts the response time is within the expected duration
func (r *HTTPResponse) AssertResponseTime(maxDuration time.Duration, startTime time.Time) *HTTPResponse {
	elapsed := time.Since(startTime)
	assert.LessOrEqual(r.t, elapsed, maxDuration, "Response time should be within %v, but was %v", maxDuration, elapsed)
	return r
}

// TestTimer helps measure response times in tests
type TestTimer struct {
	startTime time.Time
}

// NewTestTimer creates a new test timer
func NewTestTimer() *TestTimer {
	return &TestTimer{
		startTime: time.Now(),
	}
}

// Elapsed returns the elapsed time since the timer was created
func (tt *TestTimer) Elapsed() time.Duration {
	return time.Since(tt.startTime)
}

// AssertMaxDuration asserts that the elapsed time is within the maximum duration
func (tt *TestTimer) AssertMaxDuration(t *testing.T, maxDuration time.Duration) {
	elapsed := tt.Elapsed()
	assert.LessOrEqual(t, elapsed, maxDuration, "Operation should complete within %v, but took %v", maxDuration, elapsed)
}

// PerformanceAssertion provides performance testing utilities
type PerformanceAssertion struct {
	t *testing.T
}

// NewPerformanceAssertion creates a new performance assertion helper
func NewPerformanceAssertion(t *testing.T) *PerformanceAssertion {
	return &PerformanceAssertion{t: t}
}

// AssertAPIResponseTime asserts API response time is within acceptable limits
func (p *PerformanceAssertion) AssertAPIResponseTime(duration time.Duration) {
	// API response time should be < 200ms for normal operations
	maxAcceptable := 200 * time.Millisecond
	assert.LessOrEqual(p.t, duration, maxAcceptable, "API response time should be < 200ms, but was %v", duration)
}

// AssertDatabaseQueryTime asserts database query time is within acceptable limits
func (p *PerformanceAssertion) AssertDatabaseQueryTime(duration time.Duration) {
	// Database queries should be < 50ms
	maxAcceptable := 50 * time.Millisecond
	assert.LessOrEqual(p.t, duration, maxAcceptable, "Database query time should be < 50ms, but was %v", duration)
}

// AssertCacheResponseTime asserts cache response time is within acceptable limits
func (p *PerformanceAssertion) AssertCacheResponseTime(duration time.Duration) {
	// Cache operations should be < 10ms
	maxAcceptable := 10 * time.Millisecond
	assert.LessOrEqual(p.t, duration, maxAcceptable, "Cache response time should be < 10ms, but was %v", duration)
}

// ConcurrencyTestHelper provides utilities for concurrent testing
type ConcurrencyTestHelper struct {
	t *testing.T
}

// NewConcurrencyTestHelper creates a new concurrency test helper
func NewConcurrencyTestHelper(t *testing.T) *ConcurrencyTestHelper {
	return &ConcurrencyTestHelper{t: t}
}

// RunConcurrent runs multiple functions concurrently and waits for completion
func (c *ConcurrencyTestHelper) RunConcurrent(functions ...func()) {
	done := make(chan bool, len(functions))

	for _, fn := range functions {
		go func(f func()) {
			defer func() {
				if r := recover(); r != nil {
					c.t.Errorf("Panic in concurrent function: %v", r)
				}
				done <- true
			}()
			f()
		}(fn)
	}

	// Wait for all functions to complete
	for i := 0; i < len(functions); i++ {
		select {
		case <-done:
			// Function completed
		case <-time.After(30 * time.Second):
			c.t.Fatal("Concurrent function timed out after 30 seconds")
		}
	}
}

// RunConcurrentWithResults runs functions concurrently and collects results
func (c *ConcurrencyTestHelper) RunConcurrentWithResults(functions ...func() interface{}) []interface{} {
	results := make(chan interface{}, len(functions))

	for _, fn := range functions {
		go func(f func() interface{}) {
			defer func() {
				if r := recover(); r != nil {
					c.t.Errorf("Panic in concurrent function: %v", r)
					results <- nil
				}
			}()
			results <- f()
		}(fn)
	}

	// Collect all results
	var collectedResults []interface{}
	for i := 0; i < len(functions); i++ {
		select {
		case result := <-results:
			collectedResults = append(collectedResults, result)
		case <-time.After(30 * time.Second):
			c.t.Fatal("Concurrent function timed out after 30 seconds")
		}
	}

	return collectedResults
}

// ValidationHelper provides validation utilities for tests
type ValidationHelper struct {
	t *testing.T
}

// NewValidationHelper creates a new validation helper
func NewValidationHelper(t *testing.T) *ValidationHelper {
	return &ValidationHelper{t: t}
}

// ValidateEmail validates email format
func (v *ValidationHelper) ValidateEmail(email string) bool {
	// Simple email validation for testing
	return len(email) > 5 &&
		   len(email) < 255 &&
		   bytes.Contains([]byte(email), []byte("@")) &&
		   bytes.Contains([]byte(email), []byte("."))
}

// ValidateUUID validates UUID format
func (v *ValidationHelper) ValidateUUID(uuid string) bool {
	// Simple UUID validation for testing
	return len(uuid) == 36 &&
		   uuid[8] == '-' &&
		   uuid[13] == '-' &&
		   uuid[18] == '-' &&
		   uuid[23] == '-'
}

// ValidateURL validates URL format
func (v *ValidationHelper) ValidateURL(url string) bool {
	// Simple URL validation for testing
	return len(url) > 7 &&
		   (bytes.HasPrefix([]byte(url), []byte("http://")) ||
		    bytes.HasPrefix([]byte(url), []byte("https://")))
}

// ValidatePositiveNumber validates that a number is positive
func (v *ValidationHelper) ValidatePositiveNumber(num interface{}) bool {
	switch n := num.(type) {
	case int:
		return n > 0
	case int32:
		return n > 0
	case int64:
		return n > 0
	case float32:
		return n > 0
	case float64:
		return n > 0
	default:
		return false
	}
}

// AssertValidEmail asserts that an email is valid
func (v *ValidationHelper) AssertValidEmail(email string) {
	assert.True(v.t, v.ValidateEmail(email), "Invalid email format: %s", email)
}

// AssertValidUUID asserts that a UUID is valid
func (v *ValidationHelper) AssertValidUUID(uuid string) {
	assert.True(v.t, v.ValidateUUID(uuid), "Invalid UUID format: %s", uuid)
}

// AssertValidURL asserts that a URL is valid
func (v *ValidationHelper) AssertValidURL(url string) {
	assert.True(v.t, v.ValidateURL(url), "Invalid URL format: %s", url)
}

// AssertPositiveNumber asserts that a number is positive
func (v *ValidationHelper) AssertPositiveNumber(num interface{}) {
	assert.True(v.t, v.ValidatePositiveNumber(num), "Number should be positive: %v", num)
}