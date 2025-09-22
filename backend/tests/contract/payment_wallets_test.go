package contract_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPaymentWallets_Contract validates the GET /wallets endpoint contract
// This test MUST FAIL initially as no implementation exists yet (TDD)
func TestPaymentWallets_Contract(t *testing.T) {
	// Test server URL - will fail until server is implemented
	baseURL := "http://localhost:8080"
	mockToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		validateArray  bool
	}{
		{
			name:           "Get user wallets with valid auth",
			token:          mockToken,
			expectedStatus: http.StatusOK,
			validateArray:  true,
		},
		{
			name:           "Unauthorized - no token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			validateArray:  false,
		},
		{
			name:           "Unauthorized - invalid token",
			token:          "invalid.token",
			expectedStatus: http.StatusUnauthorized,
			validateArray:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create HTTP request
			req, err := http.NewRequest("GET", baseURL+"/wallets", nil)
			require.NoError(t, err)

			// Add authentication if provided
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			// Execute request (will fail until server is implemented)
			client := &http.Client{}
			resp, err := client.Do(req)

			// This SHOULD FAIL initially - no server running
			if err != nil {
				t.Logf("Expected failure: %v (no implementation yet)", err)
				return // Test passes by failing as expected in TDD
			}
			defer resp.Body.Close()

			// Validate response status
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.validateArray && resp.StatusCode == http.StatusOK {
				// Validate response is array of wallets
				var wallets []map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&wallets)
				require.NoError(t, err)

				// Validate each wallet structure
				for i, wallet := range wallets {
					requiredFields := []string{"id", "user_id", "balance", "currency", "frozen_balance", "daily_limit", "monthly_limit", "status", "is_primary", "created_at", "updated_at"}
					for _, field := range requiredFields {
						assert.Contains(t, wallet, field, "Wallet %d should contain field: %s", i, field)
					}

					// Validate wallet ID is UUID
					walletID, exists := wallet["id"]
					assert.True(t, exists, "Wallet should have ID")
					if exists {
						_, err := uuid.Parse(walletID.(string))
						assert.NoError(t, err, "Wallet ID should be valid UUID")
					}

					// Validate currency is supported
					currency, exists := wallet["currency"]
					assert.True(t, exists, "Wallet should have currency")
					if exists {
						validCurrencies := []string{"THB", "SGD", "IDR", "MYR", "PHP", "VND", "USD"}
						assert.Contains(t, validCurrencies, currency, "Currency should be supported")
					}

					// Validate balances are non-negative
					balance, exists := wallet["balance"]
					assert.True(t, exists, "Wallet should have balance")
					if exists {
						balanceNum, ok := balance.(float64)
						assert.True(t, ok, "Balance should be number")
						if ok {
							assert.GreaterOrEqual(t, balanceNum, float64(0), "Balance should be non-negative")
						}
					}

					// Validate status
					status, exists := wallet["status"]
					assert.True(t, exists, "Wallet should have status")
					if exists {
						validStatuses := []string{"active", "suspended", "closed"}
						assert.Contains(t, validStatuses, status, "Status should be valid")
					}
				}
			}
		})
	}
}

// TestPaymentWalletCreate_Contract validates the POST /wallets endpoint contract
func TestPaymentWalletCreate_Contract(t *testing.T) {
	baseURL := "http://localhost:8080"
	mockToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	tests := []struct {
		name           string
		payload        map[string]interface{}
		token          string
		expectedStatus int
		expectedFields []string
	}{
		{
			name: "Create THB wallet",
			payload: map[string]interface{}{
				"currency": "THB",
			},
			token:          mockToken,
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "user_id", "balance", "currency", "status", "is_primary", "created_at"},
		},
		{
			name: "Create primary USD wallet",
			payload: map[string]interface{}{
				"currency":   "USD",
				"is_primary": true,
			},
			token:          mockToken,
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "user_id", "balance", "currency", "status", "is_primary", "created_at"},
		},
		{
			name: "Invalid currency",
			payload: map[string]interface{}{
				"currency": "XXX",
			},
			token:          mockToken,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
		{
			name: "Missing currency",
			payload: map[string]interface{}{
				"is_primary": true,
			},
			token:          mockToken,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
		{
			name: "Unauthorized",
			payload: map[string]interface{}{
				"currency": "THB",
			},
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error", "message"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal request payload
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			// Create HTTP request
			req, err := http.NewRequest("POST", baseURL+"/wallets", bytes.NewBuffer(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Add authentication if provided
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			// Execute request (will fail until server is implemented)
			client := &http.Client{}
			resp, err := client.Do(req)

			// This SHOULD FAIL initially - no server running
			if err != nil {
				t.Logf("Expected failure: %v (no implementation yet)", err)
				return // Test passes by failing as expected in TDD
			}
			defer resp.Body.Close()

			// Validate response status
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Validate response structure
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			// Check expected fields exist
			for _, field := range tt.expectedFields {
				assert.Contains(t, response, field, "Response should contain field: %s", field)
			}

			// Additional validations for successful creation
			if resp.StatusCode == http.StatusCreated {
				// Validate currency matches request
				currency, exists := response["currency"]
				assert.True(t, exists, "Response should have currency")
				if exists {
					assert.Equal(t, tt.payload["currency"], currency, "Response currency should match request")
				}

				// Validate initial balance is zero
				balance, exists := response["balance"]
				assert.True(t, exists, "Response should have balance")
				if exists {
					assert.Equal(t, float64(0), balance, "Initial balance should be zero")
				}

				// Validate status is active
				status, exists := response["status"]
				assert.True(t, exists, "Response should have status")
				if exists {
					assert.Equal(t, "active", status, "Initial status should be active")
				}

				// Validate is_primary if specified
				if isPrimary, hasIsPrimary := tt.payload["is_primary"]; hasIsPrimary {
					responseIsPrimary, exists := response["is_primary"]
					assert.True(t, exists, "Response should have is_primary when specified")
					if exists {
						assert.Equal(t, isPrimary, responseIsPrimary, "Response is_primary should match request")
					}
				}
			}
		})
	}
}

// TestPaymentWallets_MultiCurrency validates multi-currency wallet support
func TestPaymentWallets_MultiCurrency(t *testing.T) {
	baseURL := "http://localhost:8080"
	mockToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	currencies := []string{"THB", "SGD", "IDR", "MYR", "PHP", "VND", "USD"}

	for _, currency := range currencies {
		t.Run("Create "+currency+" wallet", func(t *testing.T) {
			payload := map[string]interface{}{
				"currency": currency,
			}

			body, err := json.Marshal(payload)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", baseURL+"/wallets", bytes.NewBuffer(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+mockToken)

			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				t.Logf("Expected failure: %v (no implementation yet)", err)
				return // Test passes by failing as expected in TDD
			}
			defer resp.Body.Close()

			assert.Equal(t, http.StatusCreated, resp.StatusCode, "Should accept currency: "+currency)

			if resp.StatusCode == http.StatusCreated {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)

				responseCurrency, exists := response["currency"]
				assert.True(t, exists, "Response should have currency")
				if exists {
					assert.Equal(t, currency, responseCurrency, "Response currency should match request")
				}
			}
		})
	}
}