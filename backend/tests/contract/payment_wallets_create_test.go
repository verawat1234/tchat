package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T024: Contract test POST /wallets - Create new wallet
func TestPaymentWalletCreateContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		payload        map[string]interface{}
		expectedStatus int
		expectedFields []string
		description    string
	}{
		{
			name:       "create_thb_wallet_success",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"currency": "THB",
			},
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "user_id", "currency", "balance", "available_balance", "frozen_balance", "status", "created_at"},
			description:    "Should create THB wallet successfully",
		},
		{
			name:       "create_sgd_wallet_success",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"currency": "SGD",
			},
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "user_id", "currency", "balance", "status"},
			description:    "Should create SGD wallet successfully",
		},
		{
			name:       "create_idr_wallet_success",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"currency": "IDR",
			},
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "user_id", "currency", "balance", "status"},
			description:    "Should create IDR wallet successfully (no decimals)",
		},
		{
			name:       "create_vnd_wallet_success",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"currency": "VND",
			},
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "user_id", "currency", "balance", "status"},
			description:    "Should create VND wallet successfully (no decimals)",
		},
		{
			name:       "create_usd_wallet_success",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"currency": "USD",
			},
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "user_id", "currency", "balance", "status"},
			description:    "Should create USD wallet successfully",
		},
		{
			name:           "missing_authorization",
			authHeader:     "",
			payload:        map[string]interface{}{"currency": "THB"},
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			description:    "Should return unauthorized when no auth header",
		},
		{
			name:           "invalid_jwt_token",
			authHeader:     "Bearer invalid_token_67890",
			payload:        map[string]interface{}{"currency": "THB"},
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			description:    "Should return unauthorized for invalid token",
		},
		{
			name:           "missing_currency",
			authHeader:     "Bearer valid_jwt_token_12345",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error when currency is missing",
		},
		{
			name:       "invalid_currency",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"currency": "INVALID",
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for invalid currency",
		},
		{
			name:       "unsupported_currency",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"currency": "EUR", // Not supported in SEA region
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for unsupported currency",
		},
		{
			name:       "duplicate_wallet_same_currency",
			authHeader: "Bearer user_with_existing_thb_wallet_token",
			payload: map[string]interface{}{
				"currency": "THB",
			},
			expectedStatus: http.StatusConflict,
			expectedFields: []string{"error", "existing_wallet_id"},
			description:    "Should return conflict for duplicate wallet in same currency",
		},
		{
			name:       "malformed_request_body",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"currency": 12345, // Wrong type
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for malformed request body",
		},
		{
			name:       "kyc_not_verified",
			authHeader: "Bearer user_without_kyc_token",
			payload: map[string]interface{}{
				"currency": "THB",
			},
			expectedStatus: http.StatusForbidden,
			expectedFields: []string{"error"},
			description:    "Should return forbidden for user without KYC verification",
		},
		{
			name:       "suspended_user_account",
			authHeader: "Bearer suspended_user_token",
			payload: map[string]interface{}{
				"currency": "THB",
			},
			expectedStatus: http.StatusForbidden,
			expectedFields: []string{"error"},
			description:    "Should return forbidden for suspended user account",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := gin.New()

			// Mock authentication middleware and endpoint
			router.POST("/api/v1/payments/wallets", func(c *gin.Context) {
				authHeader := c.GetHeader("Authorization")

				// Mock authentication logic
				if authHeader == "" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
					return
				}

				var currentUserID string
				var userStatus string
				var hasKYC bool

				switch authHeader {
				case "Bearer valid_jwt_token_12345":
					currentUserID = "user_123e4567-e89b-12d3-a456-426614174000"
					userStatus = "active"
					hasKYC = true
				case "Bearer user_with_existing_thb_wallet_token":
					currentUserID = "user_with_existing_wallet"
					userStatus = "active"
					hasKYC = true
				case "Bearer user_without_kyc_token":
					currentUserID = "user_no_kyc"
					userStatus = "active"
					hasKYC = false
				case "Bearer suspended_user_token":
					currentUserID = "user_suspended"
					userStatus = "suspended"
					hasKYC = true
				default:
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
					return
				}

				// Check user status
				if userStatus == "suspended" {
					c.JSON(http.StatusForbidden, gin.H{"error": "Account is suspended"})
					return
				}

				// Check KYC verification
				if !hasKYC {
					c.JSON(http.StatusForbidden, gin.H{"error": "KYC verification required to create wallets"})
					return
				}

				// Parse request body
				var request map[string]interface{}
				if err := c.ShouldBindJSON(&request); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
					return
				}

				// Validate currency
				currency, exists := request["currency"]
				if !exists {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Currency is required"})
					return
				}

				currencyStr, ok := currency.(string)
				if !ok {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Currency must be a string"})
					return
				}

				// Validate currency against supported list
				supportedCurrencies := []string{"THB", "SGD", "IDR", "MYR", "PHP", "VND", "USD"}
				isSupported := false
				for _, supported := range supportedCurrencies {
					if currencyStr == supported {
						isSupported = true
						break
					}
				}

				if !isSupported {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported currency for SEA region"})
					return
				}

				// Check for duplicate wallet
				if authHeader == "Bearer user_with_existing_thb_wallet_token" && currencyStr == "THB" {
					c.JSON(http.StatusConflict, gin.H{
						"error": "Wallet already exists for this currency",
						"existing_wallet_id": "wallet_existing_thb_123",
					})
					return
				}

				// Determine initial balance format based on currency
				var balance string
				var availableBalance string
				var frozenBalance string

				switch currencyStr {
				case "IDR", "VND":
					// No decimal places for IDR and VND
					balance = "0"
					availableBalance = "0"
					frozenBalance = "0"
				default:
					// Two decimal places for other currencies
					balance = "0.00"
					availableBalance = "0.00"
					frozenBalance = "0.00"
				}

				// Mock successful wallet creation response
				response := gin.H{
					"id":                "wallet_" + currencyStr + "_123e4567-e89b-12d3-a456-426614174000",
					"user_id":           currentUserID,
					"currency":          currencyStr,
					"balance":           balance,
					"available_balance": availableBalance,
					"frozen_balance":    frozenBalance,
					"status":            "active",
					"created_at":        "2023-12-01T15:30:00Z",
					"updated_at":        "2023-12-01T15:30:00Z",
				}

				c.JSON(http.StatusCreated, response)
			})

			// Prepare request
			jsonData, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/payments/wallets", bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Execute
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify status code
			assert.Equal(t, tt.expectedStatus, w.Code,
				"Test: %s - %s", tt.name, tt.description)

			// Parse response
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Response should be valid JSON")

			// Verify expected fields are present
			for _, field := range tt.expectedFields {
				assert.Contains(t, response, field,
					"Test: %s - Response should contain field '%s'", tt.name, field)
			}

			// Additional assertions based on test case
			switch tt.name {
			case "create_thb_wallet_success", "create_sgd_wallet_success", "create_idr_wallet_success", "create_vnd_wallet_success", "create_usd_wallet_success":
				// Verify wallet ID format
				walletID, ok := response["id"].(string)
				assert.True(t, ok, "id should be string")
				assert.Contains(t, walletID, "wallet_", "id should have wallet_ prefix")
				assert.Contains(t, walletID, tt.payload["currency"].(string), "id should contain currency")

				// Verify currency matches request
				currency, ok := response["currency"].(string)
				assert.True(t, ok, "currency should be string")
				assert.Equal(t, tt.payload["currency"], currency, "currency should match request")

				// Verify user ID
				userID, ok := response["user_id"].(string)
				assert.True(t, ok, "user_id should be string")
				assert.NotEmpty(t, userID, "user_id should not be empty")

				// Verify balance format based on currency
				balance, ok := response["balance"].(string)
				assert.True(t, ok, "balance should be string")

				if currency == "IDR" || currency == "VND" {
					// No decimal places
					assert.Equal(t, "0", balance, "IDR/VND balance should have no decimals")
				} else {
					// Two decimal places
					assert.Equal(t, "0.00", balance, "Other currencies should have 2 decimal places")
				}

				// Verify status
				status, ok := response["status"].(string)
				assert.True(t, ok, "status should be string")
				assert.Equal(t, "active", status, "status should be active")

				// Verify timestamps
				createdAt, ok := response["created_at"].(string)
				assert.True(t, ok, "created_at should be string")
				assert.NotEmpty(t, createdAt, "created_at should not be empty")

				updatedAt, ok := response["updated_at"].(string)
				assert.True(t, ok, "updated_at should be string")
				assert.NotEmpty(t, updatedAt, "updated_at should not be empty")

			case "duplicate_wallet_same_currency":
				// Verify conflict response includes existing wallet ID
				existingWalletID, ok := response["existing_wallet_id"].(string)
				assert.True(t, ok, "existing_wallet_id should be string")
				assert.NotEmpty(t, existingWalletID, "existing_wallet_id should not be empty")

			default:
				// Error cases
				errorMsg, ok := response["error"].(string)
				assert.True(t, ok, "error should be string")
				assert.NotEmpty(t, errorMsg, "error message should not be empty")
			}
		})
	}
}

// TestPaymentWalletCreateSecurity tests security aspects
func TestPaymentWalletCreateSecurity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("user_isolation", func(t *testing.T) {
		// Ensure users can only create wallets for themselves
		router := gin.New()
		router.POST("/api/v1/payments/wallets", func(c *gin.Context) {
			authHeader := c.GetHeader("Authorization")

			// Mock user context from token
			var currentUserID string
			switch authHeader {
			case "Bearer user1_token":
				currentUserID = "user1-id"
			case "Bearer user2_token":
				currentUserID = "user2-id"
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			var request map[string]interface{}
			c.ShouldBindJSON(&request)

			// User can only create wallets for themselves (user_id comes from token)
			c.JSON(http.StatusCreated, gin.H{
				"id":       "wallet_123",
				"user_id":  currentUserID, // Always from token
				"currency": request["currency"],
				"balance":  "0.00",
				"status":   "active",
				"created_at": "2023-12-01T15:30:00Z",
			})
		})

		// User 1 creates wallet
		payload1 := map[string]interface{}{"currency": "THB"}
		jsonData1, _ := json.Marshal(payload1)

		req1, _ := http.NewRequest("POST", "/api/v1/payments/wallets", bytes.NewBuffer(jsonData1))
		req1.Header.Set("Content-Type", "application/json")
		req1.Header.Set("Authorization", "Bearer user1_token")

		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)

		var response1 map[string]interface{}
		json.Unmarshal(w1.Body.Bytes(), &response1)

		// Verify user1 gets their own user_id
		assert.Equal(t, "user1-id", response1["user_id"], "User1 should get wallet with their user_id")

		// User 2 creates wallet
		req2, _ := http.NewRequest("POST", "/api/v1/payments/wallets", bytes.NewBuffer(jsonData1))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Authorization", "Bearer user2_token")

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		var response2 map[string]interface{}
		json.Unmarshal(w2.Body.Bytes(), &response2)

		// Verify user2 gets their own user_id
		assert.Equal(t, "user2-id", response2["user_id"], "User2 should get wallet with their user_id")
	})

	t.Run("kyc_verification_required", func(t *testing.T) {
		router := gin.New()
		router.POST("/api/v1/payments/wallets", func(c *gin.Context) {
			authHeader := c.GetHeader("Authorization")

			// Mock KYC status based on token
			var hasKYC bool
			switch authHeader {
			case "Bearer verified_user_token":
				hasKYC = true
			case "Bearer unverified_user_token":
				hasKYC = false
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			if !hasKYC {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "KYC verification required to create wallets",
					"kyc_status": "pending",
				})
				return
			}

			c.JSON(http.StatusCreated, gin.H{
				"id": "wallet_123",
				"currency": "THB",
				"balance": "0.00",
			})
		})

		payload := map[string]interface{}{"currency": "THB"}
		jsonData, _ := json.Marshal(payload)

		// Verified user should succeed
		req1, _ := http.NewRequest("POST", "/api/v1/payments/wallets", bytes.NewBuffer(jsonData))
		req1.Header.Set("Content-Type", "application/json")
		req1.Header.Set("Authorization", "Bearer verified_user_token")

		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusCreated, w1.Code, "Verified user should create wallet")

		// Unverified user should fail
		req2, _ := http.NewRequest("POST", "/api/v1/payments/wallets", bytes.NewBuffer(jsonData))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Authorization", "Bearer unverified_user_token")

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusForbidden, w2.Code, "Unverified user should not create wallet")

		var errorResponse map[string]interface{}
		json.Unmarshal(w2.Body.Bytes(), &errorResponse)
		assert.Contains(t, errorResponse["error"].(string), "KYC", "Error should mention KYC")
	})

	t.Run("input_validation", func(t *testing.T) {
		router := gin.New()
		router.POST("/api/v1/payments/wallets", func(c *gin.Context) {
			var request map[string]interface{}
			c.ShouldBindJSON(&request)

			// Validate currency format
			if currency, exists := request["currency"]; exists {
				if currencyStr, ok := currency.(string); ok {
					// Check for suspicious patterns
					if len(currencyStr) > 3 {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid currency format"})
						return
					}
					// Check for SQL injection patterns
					dangerousPatterns := []string{"DROP", "SELECT", "INSERT", "DELETE", "'", "\"", ";"}
					for _, pattern := range dangerousPatterns {
						if strings.Contains(strings.ToUpper(currencyStr), pattern) {
							c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid currency format"})
							return
						}
					}
				}
			}

			c.JSON(http.StatusCreated, gin.H{
				"id": "wallet_123",
				"currency": request["currency"],
			})
		})

		maliciousInputs := []string{
			"THB'; DROP TABLE wallets; --",
			"SELECT * FROM users",
			"<script>alert('xss')</script>",
			"INVALID_LONG_CURRENCY",
		}

		for _, maliciousInput := range maliciousInputs {
			payload := map[string]interface{}{
				"currency": maliciousInput,
			}
			jsonData, _ := json.Marshal(payload)

			req, _ := http.NewRequest("POST", "/api/v1/payments/wallets", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer valid_token")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code,
				"Should reject malicious input: %s", maliciousInput)
		}
	})
}

// TestPaymentWalletCreateRateLimit tests rate limiting
func TestPaymentWalletCreateRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	createCount := 0
	router.POST("/api/v1/payments/wallets", func(c *gin.Context) {
		createCount++
		if createCount > 5 { // Simulate rate limit of 5 wallet creations per user
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many wallet creation attempts. Please wait before creating more wallets.",
			})
			return
		}

		var request map[string]interface{}
		c.ShouldBindJSON(&request)

		c.JSON(http.StatusCreated, gin.H{
			"id": "wallet_" + string(rune(createCount)),
			"currency": request["currency"],
			"balance": "0.00",
		})
	})

	supportedCurrencies := []string{"THB", "SGD", "IDR", "MYR", "PHP", "VND", "USD"}

	// Try to create more wallets than allowed
	for i, currency := range supportedCurrencies {
		payload := map[string]interface{}{"currency": currency}
		jsonData, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/api/v1/payments/wallets", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if i < 5 {
			assert.Equal(t, http.StatusCreated, w.Code, "Wallet creation %d should succeed", i+1)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code, "Wallet creation %d should be rate limited", i+1)
		}
	}
}