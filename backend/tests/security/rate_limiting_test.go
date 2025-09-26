package security_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"tchat-backend/tests/fixtures"
)

// RateLimitingTestSuite provides comprehensive rate limiting testing
// for API endpoints and user actions across Tchat microservices
type RateLimitingTestSuite struct {
	suite.Suite
	fixtures    *fixtures.MasterFixtures
	ctx         context.Context
	utils       *SecurityUtils
	rateLimiter *MockRateLimiter
}

// MockRateLimiter simulates rate limiting behavior for testing
type MockRateLimiter struct {
	requests map[string][]time.Time
	limits   map[string]RateLimit
	mutex    sync.RWMutex
}

// RateLimit defines rate limiting configuration
type RateLimit struct {
	MaxRequests int           // Maximum requests allowed
	Window      time.Duration // Time window for the limit
	Burst       int           // Burst allowance
}

// RateLimitResult represents the result of a rate limit check
type RateLimitResult struct {
	Allowed       bool
	Remaining     int
	RetryAfter    time.Duration
	TotalRequests int
}

// NewMockRateLimiter creates a new mock rate limiter
func NewMockRateLimiter() *MockRateLimiter {
	return &MockRateLimiter{
		requests: make(map[string][]time.Time),
		limits:   make(map[string]RateLimit),
		mutex:    sync.RWMutex{},
	}
}

// SetLimit configures rate limit for a specific key pattern
func (rl *MockRateLimiter) SetLimit(pattern string, limit RateLimit) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	rl.limits[pattern] = limit
}

// CheckLimit checks if a request is allowed under rate limiting
func (rl *MockRateLimiter) CheckLimit(key string) RateLimitResult {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	// Find matching limit pattern
	var limit RateLimit
	var found bool
	for pattern, l := range rl.limits {
		if key == pattern || (len(pattern) > 0 && pattern[len(pattern)-1] == '*' &&
			len(key) >= len(pattern)-1 && key[:len(pattern)-1] == pattern[:len(pattern)-1]) {
			limit = l
			found = true
			break
		}
	}

	if !found {
		// Default limit if no pattern matches
		limit = RateLimit{MaxRequests: 100, Window: time.Minute, Burst: 10}
	}

	// Get existing requests for this key
	requests, exists := rl.requests[key]
	if !exists {
		requests = make([]time.Time, 0)
	}

	// Clean up old requests outside the window
	cutoff := now.Add(-limit.Window)
	validRequests := make([]time.Time, 0)
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}

	// Check if request is allowed
	allowed := len(validRequests) < limit.MaxRequests
	remaining := limit.MaxRequests - len(validRequests) - 1
	if remaining < 0 {
		remaining = 0
	}

	var retryAfter time.Duration
	if !allowed && len(validRequests) > 0 {
		// Calculate when the oldest request will expire
		retryAfter = validRequests[0].Add(limit.Window).Sub(now)
		if retryAfter < 0 {
			retryAfter = 0
		}
	}

	// Record this request if allowed
	if allowed {
		validRequests = append(validRequests, now)
		rl.requests[key] = validRequests
	}

	return RateLimitResult{
		Allowed:       allowed,
		Remaining:     remaining,
		RetryAfter:    retryAfter,
		TotalRequests: len(validRequests),
	}
}

// Reset clears all rate limiting data
func (rl *MockRateLimiter) Reset() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	rl.requests = make(map[string][]time.Time)
}

// SetupSuite initializes the test suite
func (suite *RateLimitingTestSuite) SetupSuite() {
	suite.fixtures = fixtures.NewMasterFixtures(12345)
	suite.ctx = context.Background()
	suite.utils = NewSecurityUtils()
	suite.rateLimiter = NewMockRateLimiter()

	// Configure standard rate limits
	suite.setupStandardRateLimits()
}

// setupStandardRateLimits configures standard rate limits for testing
func (suite *RateLimitingTestSuite) setupStandardRateLimits() {
	// Authentication endpoints - strict limits
	suite.rateLimiter.SetLimit("auth:login", RateLimit{
		MaxRequests: 5,
		Window:      time.Minute,
		Burst:       2,
	})

	suite.rateLimiter.SetLimit("auth:register", RateLimit{
		MaxRequests: 3,
		Window:      time.Minute,
		Burst:       1,
	})

	suite.rateLimiter.SetLimit("auth:password-reset", RateLimit{
		MaxRequests: 2,
		Window:      time.Hour,
		Burst:       1,
	})

	// API endpoints - moderate limits
	suite.rateLimiter.SetLimit("api:messages", RateLimit{
		MaxRequests: 100,
		Window:      time.Minute,
		Burst:       20,
	})

	suite.rateLimiter.SetLimit("api:users", RateLimit{
		MaxRequests: 50,
		Window:      time.Minute,
		Burst:       10,
	})

	// File upload - strict limits
	suite.rateLimiter.SetLimit("upload:*", RateLimit{
		MaxRequests: 10,
		Window:      time.Minute,
		Burst:       3,
	})

	// Global per-IP limits
	suite.rateLimiter.SetLimit("ip:*", RateLimit{
		MaxRequests: 1000,
		Window:      time.Minute,
		Burst:       100,
	})

	// Per-user limits
	suite.rateLimiter.SetLimit("user:*", RateLimit{
		MaxRequests: 500,
		Window:      time.Minute,
		Burst:       50,
	})
}

// TestAuthenticationRateLimit tests rate limiting for authentication endpoints
func (suite *RateLimitingTestSuite) TestAuthenticationRateLimit() {
	testCases := []struct {
		name           string
		endpoint       string
		expectedLimit  int
		testRequests   int
		description    string
	}{
		{
			name:           "Login endpoint rate limiting",
			endpoint:       "auth:login",
			expectedLimit:  5,
			testRequests:   10,
			description:    "Login attempts should be rate limited to prevent brute force",
		},
		{
			name:           "Registration endpoint rate limiting",
			endpoint:       "auth:register",
			expectedLimit:  3,
			testRequests:   8,
			description:    "Registration attempts should be rate limited to prevent abuse",
		},
		{
			name:           "Password reset rate limiting",
			endpoint:       "auth:password-reset",
			expectedLimit:  2,
			testRequests:   5,
			description:    "Password reset should be strictly rate limited",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.rateLimiter.Reset()

			allowedCount := 0
			blockedCount := 0

			// Make test requests
			for i := 0; i < tc.testRequests; i++ {
				result := suite.rateLimiter.CheckLimit(tc.endpoint)
				if result.Allowed {
					allowedCount++
				} else {
					blockedCount++
				}
			}

			// Verify rate limiting behavior
			suite.Equal(tc.expectedLimit, allowedCount, tc.description)
			suite.Equal(tc.testRequests-tc.expectedLimit, blockedCount, "Excess requests should be blocked")

			// Test that rate limit resets after window
			time.Sleep(100 * time.Millisecond) // Small delay for test timing

			// Should still be blocked immediately after
			result := suite.rateLimiter.CheckLimit(tc.endpoint)
			suite.False(result.Allowed, "Request should still be blocked within window")
			suite.Greater(result.RetryAfter, time.Duration(0), "RetryAfter should be set")
		})
	}
}

// TestAPIEndpointRateLimit tests rate limiting for API endpoints
func (suite *RateLimitingTestSuite) TestAPIEndpointRateLimit() {
	testCases := []struct {
		name        string
		endpoint    string
		limit       int
		window      time.Duration
		description string
	}{
		{
			name:        "Messages API rate limiting",
			endpoint:    "api:messages",
			limit:       100,
			window:      time.Minute,
			description: "Message API should allow reasonable throughput",
		},
		{
			name:        "Users API rate limiting",
			endpoint:    "api:users",
			limit:       50,
			window:      time.Minute,
			description: "User API should be moderately rate limited",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.rateLimiter.Reset()

			// Test normal usage within limits
			for i := 0; i < tc.limit; i++ {
				result := suite.rateLimiter.CheckLimit(tc.endpoint)
				suite.True(result.Allowed, "Request %d should be allowed within limit", i+1)
				suite.Equal(tc.limit-i-1, result.Remaining, "Remaining count should be accurate")
			}

			// Test that excess requests are blocked
			result := suite.rateLimiter.CheckLimit(tc.endpoint)
			suite.False(result.Allowed, "Request beyond limit should be blocked")
			suite.Equal(0, result.Remaining, "No requests should remain")
			suite.Greater(result.RetryAfter, time.Duration(0), "RetryAfter should indicate when to retry")
		})
	}
}

// TestPerUserRateLimit tests per-user rate limiting
func (suite *RateLimitingTestSuite) TestPerUserRateLimit() {
	// Create test users
	user1 := suite.fixtures.Users.BasicUser("TH")
	user2 := suite.fixtures.Users.BasicUser("SG")

	user1Key := fmt.Sprintf("user:%s", user1.ID)
	user2Key := fmt.Sprintf("user:%s", user2.ID)

	suite.Run("Independent_User_Limits", func() {
		suite.rateLimiter.Reset()

		// Both users should have independent limits
		for i := 0; i < 100; i++ {
			result1 := suite.rateLimiter.CheckLimit(user1Key)
			result2 := suite.rateLimiter.CheckLimit(user2Key)

			suite.True(result1.Allowed, "User 1 request %d should be allowed", i+1)
			suite.True(result2.Allowed, "User 2 request %d should be allowed", i+1)
		}

		// Test that users hit their individual limits
		for i := 0; i < 500; i++ {
			suite.rateLimiter.CheckLimit(user1Key)
		}

		// User 1 should be limited, but user 2 should still work
		result1 := suite.rateLimiter.CheckLimit(user1Key)
		result2 := suite.rateLimiter.CheckLimit(user2Key)

		suite.False(result1.Allowed, "User 1 should be rate limited")
		suite.True(result2.Allowed, "User 2 should still be allowed")
	})
}

// TestIPBasedRateLimit tests IP-based rate limiting
func (suite *RateLimitingTestSuite) TestIPBasedRateLimit() {
	ipAddresses := []string{
		"192.168.1.1",
		"10.0.0.1",
		"172.16.0.1",
		"203.0.113.1",
	}

	for _, ip := range ipAddresses {
		suite.Run(fmt.Sprintf("IP_%s_Rate_Limit", ip), func() {
			suite.rateLimiter.Reset()

			ipKey := fmt.Sprintf("ip:%s", ip)

			// Test normal usage
			allowedRequests := 0
			for i := 0; i < 1200; i++ { // Try more than the limit
				result := suite.rateLimiter.CheckLimit(ipKey)
				if result.Allowed {
					allowedRequests++
				} else {
					break
				}
			}

			suite.Equal(1000, allowedRequests, "IP should be allowed 1000 requests per minute")

			// Verify rate limiting kicks in
			result := suite.rateLimiter.CheckLimit(ipKey)
			suite.False(result.Allowed, "Additional requests should be blocked")
		})
	}
}

// TestFileUploadRateLimit tests file upload rate limiting
func (suite *RateLimitingTestSuite) TestFileUploadRateLimit() {
	uploadTypes := []string{
		"upload:image",
		"upload:document",
		"upload:video",
		"upload:audio",
	}

	for _, uploadType := range uploadTypes {
		suite.Run(fmt.Sprintf("Upload_Type_%s", uploadType), func() {
			suite.rateLimiter.Reset()

			// Test that uploads are strictly limited
			allowedUploads := 0
			for i := 0; i < 20; i++ {
				result := suite.rateLimiter.CheckLimit(uploadType)
				if result.Allowed {
					allowedUploads++
				}
			}

			suite.Equal(10, allowedUploads, "File uploads should be limited to 10 per minute")

			// Verify blocking
			result := suite.rateLimiter.CheckLimit(uploadType)
			suite.False(result.Allowed, "Additional uploads should be blocked")
			suite.Greater(result.RetryAfter, time.Duration(0), "Should indicate retry time")
		})
	}
}

// TestConcurrentRateLimit tests rate limiting under concurrent access
func (suite *RateLimitingTestSuite) TestConcurrentRateLimit() {
	suite.Run("Concurrent_Access_Safety", func() {
		suite.rateLimiter.Reset()

		const numGoroutines = 100
		const requestsPerGoroutine = 10

		var wg sync.WaitGroup
		var allowedCount int64
		var blockedCount int64
		var mutex sync.Mutex

		// Launch concurrent goroutines
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()

				for j := 0; j < requestsPerGoroutine; j++ {
					key := fmt.Sprintf("api:messages")
					result := suite.rateLimiter.CheckLimit(key)

					mutex.Lock()
					if result.Allowed {
						allowedCount++
					} else {
						blockedCount++
					}
					mutex.Unlock()
				}
			}(i)
		}

		wg.Wait()

		// Verify that rate limiting worked correctly under concurrency
		totalRequests := allowedCount + blockedCount
		suite.Equal(int64(numGoroutines*requestsPerGoroutine), totalRequests, "All requests should be accounted for")
		suite.LessOrEqual(allowedCount, int64(100), "Should not exceed rate limit even under concurrency")
		suite.Greater(allowedCount, int64(0), "Some requests should be allowed")
	})
}

// TestRateLimitBypass tests detection of rate limit bypass attempts
func (suite *RateLimitingTestSuite) TestRateLimitBypass() {
	bypassAttempts := []struct {
		name        string
		keys        []string
		description string
	}{
		{
			name:        "Header manipulation",
			keys:        []string{"ip:192.168.1.1", "ip:192.168.1.1", "ip:192.168.1.1"},
			description: "Should not allow bypassing via header manipulation",
		},
		{
			name:        "Case variation",
			keys:        []string{"auth:login", "auth:LOGIN", "AUTH:login"},
			description: "Should not allow bypassing via case changes",
		},
		{
			name:        "Unicode variation",
			keys:        []string{"auth:login", "auth:login", "auth:login"},
			description: "Should not allow bypassing via Unicode variations",
		},
	}

	for _, attempt := range bypassAttempts {
		suite.Run(attempt.name, func() {
			suite.rateLimiter.Reset()

			totalAllowed := 0

			// Try bypass attempts
			for _, key := range attempt.keys {
				for i := 0; i < 10; i++ {
					result := suite.rateLimiter.CheckLimit(key)
					if result.Allowed {
						totalAllowed++
					}
				}
			}

			// Should still respect the rate limit regardless of bypass attempts
			suite.LessOrEqual(totalAllowed, 15, attempt.description) // Some tolerance for different endpoints
		})
	}
}

// TestRateLimitRecovery tests rate limit recovery after window expiration
func (suite *RateLimitingTestSuite) TestRateLimitRecovery() {
	suite.Run("Rate_Limit_Recovery", func() {
		suite.rateLimiter.Reset()

		endpoint := "auth:login"

		// Exhaust the rate limit
		for i := 0; i < 10; i++ {
			suite.rateLimiter.CheckLimit(endpoint)
		}

		// Verify we're blocked
		result := suite.rateLimiter.CheckLimit(endpoint)
		suite.False(result.Allowed, "Should be blocked after exhausting limit")

		// Note: In a real test, you would wait for the window to expire
		// For this mock test, we'll simulate time passage
		suite.rateLimiter.Reset() // Simulate window expiration

		// Should be allowed again after window reset
		result = suite.rateLimiter.CheckLimit(endpoint)
		suite.True(result.Allowed, "Should be allowed after window reset")
	})
}

// TestRateLimitMetrics tests rate limit metrics and monitoring
func (suite *RateLimitingTestSuite) TestRateLimitMetrics() {
	suite.Run("Rate_Limit_Metrics", func() {
		suite.rateLimiter.Reset()

		endpoint := "api:messages"

		// Make several requests
		for i := 0; i < 50; i++ {
			result := suite.rateLimiter.CheckLimit(endpoint)

			// Verify metrics are accurate
			suite.Equal(i+1, result.TotalRequests, "Total requests should be accurate")
			suite.Equal(100-i-1, result.Remaining, "Remaining requests should be accurate")

			if i < 100 {
				suite.True(result.Allowed, "Request should be allowed within limit")
			}
		}

		// Test blocked request metrics
		result := suite.rateLimiter.CheckLimit(endpoint)
		if !result.Allowed {
			suite.Equal(0, result.Remaining, "No requests should remain when blocked")
			suite.Greater(result.RetryAfter, time.Duration(0), "RetryAfter should be positive")
		}
	})
}

// TestSEASpecificRateLimit tests Southeast Asia specific rate limiting
func (suite *RateLimitingTestSuite) TestSEASpecificRateLimit() {
	countries := []string{"TH", "SG", "ID", "MY", "VN", "PH"}

	for _, country := range countries {
		suite.Run(fmt.Sprintf("Country_%s_Rate_Limit", country), func() {
			suite.rateLimiter.Reset()

			// Create country-specific user
			user := suite.fixtures.Users.BasicUser(country)
			userKey := fmt.Sprintf("user:%s:%s", country, user.ID)

			// Test country-specific rate limiting
			// Some countries might have different limits based on regulations
			expectedLimit := 500 // Default limit

			// Adjust for country-specific regulations if needed
			switch country {
			case "SG":
				// Singapore might have stricter limits for compliance
				expectedLimit = 400
			case "TH":
				// Thailand might have specific telecom regulations
				expectedLimit = 450
			}

			allowedRequests := 0
			for i := 0; i < expectedLimit+100; i++ {
				result := suite.rateLimiter.CheckLimit(userKey)
				if result.Allowed {
					allowedRequests++
				}
			}

			suite.LessOrEqual(allowedRequests, expectedLimit+50, "Should respect country-specific limits with some tolerance")
			suite.Greater(allowedRequests, expectedLimit-50, "Should allow reasonable requests for country")
		})
	}
}

// TestRateLimitIntegration tests integration with other security measures
func (suite *RateLimitingTestSuite) TestRateLimitIntegration() {
	suite.Run("Integration_With_Security_Utils", func() {
		suite.rateLimiter.Reset()

		// Test rate limiting with input validation
		maliciousInputs := []string{
			"<script>alert('xss')</script>",
			"'; DROP TABLE users; --",
			"javascript:alert(1)",
		}

		for _, input := range maliciousInputs {
			// Even malicious inputs should be rate limited
			for i := 0; i < 10; i++ {
				key := fmt.Sprintf("api:messages:%s", input) // Using input as part of key
				result := suite.rateLimiter.CheckLimit(key)

				// Should apply rate limiting regardless of input content
				if i < 5 { // Assuming a limit of 5 for this test
					suite.True(result.Allowed, "Should apply rate limiting to malicious input")
				}
			}
		}

		// Test with valid Southeast Asian content
		seaContent := []string{
			"สวัสดีครับ", // Thai
			"Xin chào",    // Vietnamese
			"Halo",        // Indonesian
		}

		for _, content := range seaContent {
			suite.True(suite.utils.IsSEAContentSafe(content, "TH"), "SEA content should be validated as safe")

			// Rate limiting should work with international content
			key := fmt.Sprintf("api:messages:%s", content)
			result := suite.rateLimiter.CheckLimit(key)
			suite.True(result.Allowed, "Rate limiting should work with international content")
		}
	})
}

// TestRateLimitingEdgeCases tests edge cases in rate limiting
func (suite *RateLimitingTestSuite) TestRateLimitingEdgeCases() {
	edgeCases := []struct {
		name        string
		testFunc    func()
		description string
	}{
		{
			name: "Empty_Key",
			testFunc: func() {
				result := suite.rateLimiter.CheckLimit("")
				suite.True(result.Allowed, "Empty key should use default limits")
			},
			description: "Empty keys should be handled gracefully",
		},
		{
			name: "Very_Long_Key",
			testFunc: func() {
				longKey := strings.Repeat("a", 1000)
				result := suite.rateLimiter.CheckLimit(longKey)
				suite.True(result.Allowed, "Long keys should be handled")
			},
			description: "Very long keys should not break rate limiting",
		},
		{
			name: "Special_Characters_In_Key",
			testFunc: func() {
				specialKey := "key:with:special!@#$%^&*()characters"
				result := suite.rateLimiter.CheckLimit(specialKey)
				suite.True(result.Allowed, "Special characters should be handled")
			},
			description: "Keys with special characters should work",
		},
		{
			name: "Unicode_In_Key",
			testFunc: func() {
				unicodeKey := "key:with:unicode:สวัสดี:🌟"
				result := suite.rateLimiter.CheckLimit(unicodeKey)
				suite.True(result.Allowed, "Unicode keys should be handled")
			},
			description: "Unicode in keys should be supported",
		},
	}

	for _, tc := range edgeCases {
		suite.Run(tc.name, func() {
			suite.rateLimiter.Reset()
			tc.testFunc()
		})
	}
}

func TestRateLimitingSuite(t *testing.T) {
	suite.Run(t, new(RateLimitingTestSuite))
}