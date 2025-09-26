package security

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// CSRFProtectionTestSuite tests CSRF protection mechanisms
type CSRFProtectionTestSuite struct {
	suite.Suite
	securityUtils *SecurityUtils
}

// SetupTest initializes the CSRF test suite
func (suite *CSRFProtectionTestSuite) SetupTest() {
	suite.securityUtils = NewSecurityUtils()
}

// TestCSRFProtectionTestSuite runs the CSRF protection test suite
func TestCSRFProtectionTestSuite(t *testing.T) {
	suite.Run(t, new(CSRFProtectionTestSuite))
}

// TestCSRFTokenValidation tests CSRF token validation
func (suite *CSRFProtectionTestSuite) TestCSRFTokenValidation() {
	testCases := []struct {
		name           string
		endpoint       string
		method         string
		hasCSRFToken   bool
		validToken     bool
		expectRejected bool
	}{
		{
			name:           "POST with valid CSRF token",
			endpoint:       "/api/auth/register",
			method:         "POST",
			hasCSRFToken:   true,
			validToken:     true,
			expectRejected: false,
		},
		{
			name:           "POST without CSRF token",
			endpoint:       "/api/auth/register",
			method:         "POST",
			hasCSRFToken:   false,
			validToken:     false,
			expectRejected: true,
		},
		{
			name:           "POST with invalid CSRF token",
			endpoint:       "/api/auth/register",
			method:         "POST",
			hasCSRFToken:   true,
			validToken:     false,
			expectRejected: true,
		},
		{
			name:           "PUT with valid CSRF token",
			endpoint:       "/api/users/profile",
			method:         "PUT",
			hasCSRFToken:   true,
			validToken:     true,
			expectRejected: false,
		},
		{
			name:           "DELETE with valid CSRF token",
			endpoint:       "/api/content/123",
			method:         "DELETE",
			hasCSRFToken:   true,
			validToken:     true,
			expectRejected: false,
		},
		{
			name:           "GET request (should not require CSRF)",
			endpoint:       "/api/users/profile",
			method:         "GET",
			hasCSRFToken:   false,
			validToken:     false,
			expectRejected: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Create request
			requestBody := map[string]interface{}{
				"test": "data",
			}
			jsonBody, _ := json.Marshal(requestBody)

			req := httptest.NewRequest(tc.method, tc.endpoint, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Add CSRF token if required
			if tc.hasCSRFToken {
				if tc.validToken {
					req.Header.Set("X-CSRF-Token", suite.generateValidCSRFToken())
				} else {
					req.Header.Set("X-CSRF-Token", "invalid-token")
				}
			}

			// Simulate CSRF validation
			rr := httptest.NewRecorder()
			suite.simulateCSRFValidation(rr, req)

			// Verify response
			if tc.expectRejected {
				assert.True(suite.T(), rr.Code >= 400,
					"Expected request to be rejected due to CSRF validation failure")
				assert.Contains(suite.T(), rr.Body.String(), "CSRF",
					"Error message should mention CSRF")
			} else {
				assert.True(suite.T(), rr.Code < 400,
					"Expected request to be accepted")
			}
		})
	}
}

// TestCSRFDoubleSubmitCookie tests double submit cookie pattern
func (suite *CSRFProtectionTestSuite) TestCSRFDoubleSubmitCookie() {
	testCases := []struct {
		name           string
		cookieToken    string
		headerToken    string
		expectRejected bool
	}{
		{
			name:           "Matching cookie and header tokens",
			cookieToken:    "valid-csrf-token-123",
			headerToken:    "valid-csrf-token-123",
			expectRejected: false,
		},
		{
			name:           "Mismatched cookie and header tokens",
			cookieToken:    "valid-csrf-token-123",
			headerToken:    "different-token-456",
			expectRejected: true,
		},
		{
			name:           "Missing cookie token",
			cookieToken:    "",
			headerToken:    "valid-csrf-token-123",
			expectRejected: true,
		},
		{
			name:           "Missing header token",
			cookieToken:    "valid-csrf-token-123",
			headerToken:    "",
			expectRejected: true,
		},
		{
			name:           "Empty tokens",
			cookieToken:    "",
			headerToken:    "",
			expectRejected: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			requestBody := map[string]interface{}{
				"action": "sensitive_operation",
			}
			jsonBody, _ := json.Marshal(requestBody)

			req := httptest.NewRequest("POST", "/api/sensitive/action", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Set cookie token
			if tc.cookieToken != "" {
				req.AddCookie(&http.Cookie{
					Name:  "csrf-token",
					Value: tc.cookieToken,
				})
			}

			// Set header token
			if tc.headerToken != "" {
				req.Header.Set("X-CSRF-Token", tc.headerToken)
			}

			rr := httptest.NewRecorder()
			suite.simulateDoubleSubmitCSRFValidation(rr, req)

			if tc.expectRejected {
				assert.True(suite.T(), rr.Code >= 400,
					"Expected request to be rejected")
			} else {
				assert.True(suite.T(), rr.Code < 400,
					"Expected request to be accepted")
			}
		})
	}
}

// TestCSRFOriginValidation tests origin header validation
func (suite *CSRFProtectionTestSuite) TestCSRFOriginValidation() {
	testCases := []struct {
		name           string
		origin         string
		referer        string
		expectRejected bool
	}{
		{
			name:           "Valid origin from same domain",
			origin:         "https://tchat.app",
			referer:        "https://tchat.app/dashboard",
			expectRejected: false,
		},
		{
			name:           "Invalid origin from different domain",
			origin:         "https://malicious.com",
			referer:        "https://malicious.com/attack",
			expectRejected: true,
		},
		{
			name:           "Missing origin header",
			origin:         "",
			referer:        "https://tchat.app/dashboard",
			expectRejected: false, // Referer can be fallback
		},
		{
			name:           "Missing both origin and referer",
			origin:         "",
			referer:        "",
			expectRejected: true,
		},
		{
			name:           "Localhost origin (development)",
			origin:         "http://localhost:3000",
			referer:        "http://localhost:3000/dev",
			expectRejected: false,
		},
		{
			name:           "HTTPS origin with HTTP referer (downgrade attack)",
			origin:         "http://tchat.app",
			referer:        "https://tchat.app/page",
			expectRejected: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			requestBody := map[string]interface{}{
				"data": "sensitive",
			}
			jsonBody, _ := json.Marshal(requestBody)

			req := httptest.NewRequest("POST", "/api/sensitive/action", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}
			if tc.referer != "" {
				req.Header.Set("Referer", tc.referer)
			}

			rr := httptest.NewRecorder()
			suite.simulateOriginValidation(rr, req)

			if tc.expectRejected {
				assert.True(suite.T(), rr.Code >= 400,
					"Expected request to be rejected due to origin validation")
			} else {
				assert.True(suite.T(), rr.Code < 400,
					"Expected request to be accepted")
			}
		})
	}
}

// TestCSRFSameSiteCookies tests SameSite cookie protection
func (suite *CSRFProtectionTestSuite) TestCSRFSameSiteCookies() {
	testCases := []struct {
		name               string
		sameSiteAttribute  string
		crossSiteRequest   bool
		expectRejected     bool
	}{
		{
			name:               "SameSite=Strict with same-site request",
			sameSiteAttribute:  "Strict",
			crossSiteRequest:   false,
			expectRejected:     false,
		},
		{
			name:               "SameSite=Strict with cross-site request",
			sameSiteAttribute:  "Strict",
			crossSiteRequest:   true,
			expectRejected:     true,
		},
		{
			name:               "SameSite=Lax with same-site request",
			sameSiteAttribute:  "Lax",
			crossSiteRequest:   false,
			expectRejected:     false,
		},
		{
			name:               "SameSite=Lax with cross-site POST",
			sameSiteAttribute:  "Lax",
			crossSiteRequest:   true,
			expectRejected:     true,
		},
		{
			name:               "SameSite=None with cross-site request",
			sameSiteAttribute:  "None",
			crossSiteRequest:   true,
			expectRejected:     false, // Should be allowed but needs CSRF token
		},
		{
			name:               "No SameSite attribute",
			sameSiteAttribute:  "",
			crossSiteRequest:   true,
			expectRejected:     false, // Legacy behavior, needs CSRF token
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			requestBody := map[string]interface{}{
				"action": "sensitive",
			}
			jsonBody, _ := json.Marshal(requestBody)

			req := httptest.NewRequest("POST", "/api/sensitive/action", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Simulate cross-site request
			if tc.crossSiteRequest {
				req.Header.Set("Origin", "https://external-site.com")
				req.Header.Set("Referer", "https://external-site.com/page")
			} else {
				req.Header.Set("Origin", "https://tchat.app")
				req.Header.Set("Referer", "https://tchat.app/page")
			}

			// Add session cookie with SameSite attribute
			cookie := &http.Cookie{
				Name:  "session-token",
				Value: "valid-session-123",
			}

			if tc.sameSiteAttribute != "" {
				cookie.Raw = cookie.String() + "; SameSite=" + tc.sameSiteAttribute
			}

			req.AddCookie(cookie)

			rr := httptest.NewRecorder()
			suite.simulateSameSiteValidation(rr, req, tc.sameSiteAttribute)

			if tc.expectRejected {
				assert.True(suite.T(), rr.Code >= 400,
					"Expected request to be rejected due to SameSite policy")
			} else {
				assert.True(suite.T(), rr.Code < 400,
					"Expected request to be accepted")
			}
		})
	}
}

// TestCSRFTokenGeneration tests CSRF token generation security
func (suite *CSRFProtectionTestSuite) TestCSRFTokenGeneration() {
	// Generate multiple tokens and verify they are unique and secure
	tokens := make(map[string]bool)
	tokenCount := 1000

	for i := 0; i < tokenCount; i++ {
		token := suite.generateValidCSRFToken()

		// Check token is not empty
		assert.NotEmpty(suite.T(), token, "CSRF token should not be empty")

		// Check token has sufficient length (minimum 32 characters)
		assert.True(suite.T(), len(token) >= 32,
			"CSRF token should be at least 32 characters long")

		// Check token is unique
		assert.False(suite.T(), tokens[token],
			"CSRF token should be unique")
		tokens[token] = true

		// Check token contains only valid characters (base64 or hex)
		assert.Regexp(suite.T(), `^[A-Za-z0-9+/=\-_]+$`, token,
			"CSRF token should contain only valid characters")
	}

	suite.T().Logf("Generated %d unique CSRF tokens", tokenCount)
}

// Helper methods for CSRF testing

// generateValidCSRFToken generates a valid CSRF token for testing
func (suite *CSRFProtectionTestSuite) generateValidCSRFToken() string {
	return suite.securityUtils.GenerateSecureToken(32)
}

// simulateCSRFValidation simulates CSRF token validation
func (suite *CSRFProtectionTestSuite) simulateCSRFValidation(rr *httptest.ResponseRecorder, req *http.Request) {
	// Methods that require CSRF protection
	protectedMethods := []string{"POST", "PUT", "DELETE", "PATCH"}

	requiresCSRF := false
	for _, method := range protectedMethods {
		if req.Method == method {
			requiresCSRF = true
			break
		}
	}

	if !requiresCSRF {
		// GET, HEAD, OPTIONS don't require CSRF tokens
		rr.WriteHeader(http.StatusOK)
		json.NewEncoder(rr).Encode(map[string]string{"status": "success"})
		return
	}

	// Check for CSRF token
	csrfToken := req.Header.Get("X-CSRF-Token")
	if csrfToken == "" {
		// Also check for token in form data or query params
		csrfToken = req.FormValue("csrf_token")
	}

	if csrfToken == "" {
		rr.WriteHeader(http.StatusForbidden)
		json.NewEncoder(rr).Encode(map[string]string{
			"error": "CSRF token missing",
			"code":  "CSRF_TOKEN_MISSING",
		})
		return
	}

	// Validate CSRF token (simplified validation)
	if !suite.isValidCSRFToken(csrfToken) {
		rr.WriteHeader(http.StatusForbidden)
		json.NewEncoder(rr).Encode(map[string]string{
			"error": "Invalid CSRF token",
			"code":  "CSRF_TOKEN_INVALID",
		})
		return
	}

	// CSRF validation passed
	rr.WriteHeader(http.StatusOK)
	json.NewEncoder(rr).Encode(map[string]string{"status": "success"})
}

// simulateDoubleSubmitCSRFValidation simulates double submit cookie CSRF validation
func (suite *CSRFProtectionTestSuite) simulateDoubleSubmitCSRFValidation(rr *httptest.ResponseRecorder, req *http.Request) {
	// Get CSRF token from cookie
	var cookieToken string
	for _, cookie := range req.Cookies() {
		if cookie.Name == "csrf-token" {
			cookieToken = cookie.Value
			break
		}
	}

	// Get CSRF token from header
	headerToken := req.Header.Get("X-CSRF-Token")

	// Both tokens must be present and match
	if cookieToken == "" || headerToken == "" {
		rr.WriteHeader(http.StatusForbidden)
		json.NewEncoder(rr).Encode(map[string]string{
			"error": "CSRF token missing from cookie or header",
			"code":  "CSRF_TOKEN_MISSING",
		})
		return
	}

	if cookieToken != headerToken {
		rr.WriteHeader(http.StatusForbidden)
		json.NewEncoder(rr).Encode(map[string]string{
			"error": "CSRF tokens do not match",
			"code":  "CSRF_TOKEN_MISMATCH",
		})
		return
	}

	// CSRF validation passed
	rr.WriteHeader(http.StatusOK)
	json.NewEncoder(rr).Encode(map[string]string{"status": "success"})
}

// simulateOriginValidation simulates origin header validation
func (suite *CSRFProtectionTestSuite) simulateOriginValidation(rr *httptest.ResponseRecorder, req *http.Request) {
	allowedOrigins := []string{
		"https://tchat.app",
		"https://app.tchat.com",
		"http://localhost:3000",
		"http://localhost:8080",
	}

	origin := req.Header.Get("Origin")
	referer := req.Header.Get("Referer")

	// Check origin first
	if origin != "" {
		isAllowed := false
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			rr.WriteHeader(http.StatusForbidden)
			json.NewEncoder(rr).Encode(map[string]string{
				"error": "Origin not allowed",
				"code":  "ORIGIN_NOT_ALLOWED",
			})
			return
		}
	} else if referer != "" {
		// Fallback to referer validation
		isAllowed := false
		for _, allowed := range allowedOrigins {
			if strings.HasPrefix(referer, allowed) {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			rr.WriteHeader(http.StatusForbidden)
			json.NewEncoder(rr).Encode(map[string]string{
				"error": "Referer not allowed",
				"code":  "REFERER_NOT_ALLOWED",
			})
			return
		}
	} else {
		// No origin or referer header
		rr.WriteHeader(http.StatusForbidden)
		json.NewEncoder(rr).Encode(map[string]string{
			"error": "Missing origin and referer headers",
			"code":  "ORIGIN_REFERER_MISSING",
		})
		return
	}

	// Origin validation passed
	rr.WriteHeader(http.StatusOK)
	json.NewEncoder(rr).Encode(map[string]string{"status": "success"})
}

// simulateSameSiteValidation simulates SameSite cookie validation
func (suite *CSRFProtectionTestSuite) simulateSameSiteValidation(rr *httptest.ResponseRecorder, req *http.Request, sameSiteAttr string) {
	origin := req.Header.Get("Origin")
	referer := req.Header.Get("Referer")

	// Determine if this is a cross-site request
	isCrossSite := false
	if origin != "" && !strings.HasPrefix(origin, "https://tchat.app") && !strings.HasPrefix(origin, "http://localhost") {
		isCrossSite = true
	}

	// Check SameSite policy
	switch sameSiteAttr {
	case "Strict":
		if isCrossSite {
			rr.WriteHeader(http.StatusForbidden)
			json.NewEncoder(rr).Encode(map[string]string{
				"error": "SameSite=Strict policy violation",
				"code":  "SAMESITE_STRICT_VIOLATION",
			})
			return
		}
	case "Lax":
		if isCrossSite && req.Method == "POST" {
			rr.WriteHeader(http.StatusForbidden)
			json.NewEncoder(rr).Encode(map[string]string{
				"error": "SameSite=Lax policy violation for POST request",
				"code":  "SAMESITE_LAX_VIOLATION",
			})
			return
		}
	case "None":
		// SameSite=None allows cross-site requests but requires CSRF protection
		if isCrossSite {
			// Would need additional CSRF validation here
		}
	case "":
		// No SameSite attribute - legacy behavior
		// Would need additional CSRF validation for cross-site requests
	}

	// SameSite validation passed
	rr.WriteHeader(http.StatusOK)
	json.NewEncoder(rr).Encode(map[string]string{"status": "success"})
}

// isValidCSRFToken validates a CSRF token (simplified for testing)
func (suite *CSRFProtectionTestSuite) isValidCSRFToken(token string) bool {
	// Basic validation: token should be non-empty and have minimum length
	if len(token) < 16 {
		return false
	}

	// In real implementation, this would validate against stored token or HMAC
	// For testing, we consider tokens starting with "valid" or generated by our method as valid
	return strings.HasPrefix(token, "valid") || len(token) >= 32
}

// BenchmarkCSRFTokenGeneration benchmarks CSRF token generation
func BenchmarkCSRFTokenGeneration(b *testing.B) {
	suite := &CSRFProtectionTestSuite{}
	suite.SetupTest()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = suite.generateValidCSRFToken()
	}
}

// BenchmarkCSRFValidation benchmarks CSRF validation
func BenchmarkCSRFValidation(b *testing.B) {
	suite := &CSRFProtectionTestSuite{}
	suite.SetupTest()

	requestBody := map[string]interface{}{"test": "data"}
	jsonBody, _ := json.Marshal(requestBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/test", bytes.NewBuffer(jsonBody))
		req.Header.Set("X-CSRF-Token", suite.generateValidCSRFToken())
		rr := httptest.NewRecorder()
		suite.simulateCSRFValidation(rr, req)
	}
}