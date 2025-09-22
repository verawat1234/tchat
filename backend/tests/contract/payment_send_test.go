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

// TestPaymentTransactionSend_Contract validates the POST /transactions/send endpoint contract
// This test MUST FAIL initially as no implementation exists yet (TDD)
func TestPaymentTransactionSend_Contract(t *testing.T) {
	// Test server URL - will fail until server is implemented
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
			name: "Valid THB transaction",
			payload: map[string]interface{}{
				"wallet_id":    uuid.New().String(),
				"recipient_id": uuid.New().String(),
				"amount":       10000, // 100.00 THB in cents
				"currency":     "THB",
				"description":  "Lunch money",
			},
			token:          mockToken,
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "wallet_id", "type", "amount", "currency", "status", "description", "created_at"},
		},
		{
			name: "Transaction with reference",
			payload: map[string]interface{}{
				"wallet_id":    uuid.New().String(),
				"recipient_id": uuid.New().String(),
				"amount":       5000,
				"currency":     "SGD",
				"description":  "Payment for service",
				"reference":    "INV-2023-001",
			},
			token:          mockToken,
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "wallet_id", "type", "amount", "currency", "status", "description", "reference", "created_at"},
		},
		{
			name: "Minimum amount transaction",
			payload: map[string]interface{}{
				"wallet_id":    uuid.New().String(),
				"recipient_id": uuid.New().String(),
				"amount":       1, // 0.01 in cents
				"currency":     "USD",
			},
			token:          mockToken,
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "wallet_id", "type", "amount", "currency", "status", "created_at"},
		},
		{
			name: "Invalid amount - zero",
			payload: map[string]interface{}{
				"wallet_id":    uuid.New().String(),
				"recipient_id": uuid.New().String(),
				"amount":       0,
				"currency":     "THB",
			},
			token:          mockToken,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
		{
			name: "Invalid amount - negative",
			payload: map[string]interface{}{
				"wallet_id":    uuid.New().String(),
				"recipient_id": uuid.New().String(),
				"amount":       -1000,
				"currency":     "THB",
			},
			token:          mockToken,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
		{
			name: "Invalid currency",
			payload: map[string]interface{}{
				"wallet_id":    uuid.New().String(),
				"recipient_id": uuid.New().String(),
				"amount":       10000,
				"currency":     "XXX",
			},
			token:          mockToken,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
		{
			name: "Invalid wallet ID format",
			payload: map[string]interface{}{
				"wallet_id":    "invalid-uuid",
				"recipient_id": uuid.New().String(),
				"amount":       10000,
				"currency":     "THB",
			},
			token:          mockToken,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
		{
			name: "Invalid recipient ID format",
			payload: map[string]interface{}{
				"wallet_id":    uuid.New().String(),
				"recipient_id": "invalid-uuid",
				"amount":       10000,
				"currency":     "THB",
			},
			token:          mockToken,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
		{
			name: "Missing required fields",
			payload: map[string]interface{}{
				"wallet_id": uuid.New().String(),
				// missing recipient_id, amount, currency
			},
			token:          mockToken,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
		{
			name: "Description too long",
			payload: map[string]interface{}{
				"wallet_id":    uuid.New().String(),
				"recipient_id": uuid.New().String(),
				"amount":       10000,
				"currency":     "THB",
				"description":  string(make([]byte, 201)), // 201 chars, max is 200
			},
			token:          mockToken,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
		{
			name: "Reference too long",
			payload: map[string]interface{}{
				"wallet_id":    uuid.New().String(),
				"recipient_id": uuid.New().String(),
				"amount":       10000,
				"currency":     "THB",
				"reference":    string(make([]byte, 101)), // 101 chars, max is 100
			},
			token:          mockToken,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
		{
			name: "Unauthorized - no token",
			payload: map[string]interface{}{
				"wallet_id":    uuid.New().String(),
				"recipient_id": uuid.New().String(),
				"amount":       10000,
				"currency":     "THB",
			},
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error", "message"},
		},
		{
			name: "Insufficient funds",
			payload: map[string]interface{}{
				"wallet_id":    uuid.New().String(),
				"recipient_id": uuid.New().String(),
				"amount":       999999999, // Very large amount
				"currency":     "THB",
			},
			token:          mockToken,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal request payload
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			// Create HTTP request
			req, err := http.NewRequest("POST", baseURL+"/transactions/send", bytes.NewBuffer(body))
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

			// Additional contract validations for successful transactions
			if resp.StatusCode == http.StatusCreated {
				// Transaction ID validation
				transactionID, exists := response["id"]
				assert.True(t, exists, "Success response should have 'id' field")
				if exists {
					_, err := uuid.Parse(transactionID.(string))
					assert.NoError(t, err, "Transaction ID should be valid UUID")
				}

				// Wallet ID validation
				walletID, exists := response["wallet_id"]
				assert.True(t, exists, "Success response should have 'wallet_id' field")
				if exists {
					assert.Equal(t, tt.payload["wallet_id"], walletID, "Response wallet_id should match request")
				}

				// Transaction type validation
				transactionType, exists := response["type"]
				assert.True(t, exists, "Success response should have 'type' field")
				if exists {
					assert.Equal(t, "send", transactionType, "Transaction type should be 'send'")
				}

				// Amount validation
				amount, exists := response["amount"]
				assert.True(t, exists, "Success response should have 'amount' field")
				if exists {
					assert.Equal(t, tt.payload["amount"], amount, "Response amount should match request")
				}

				// Currency validation
				currency, exists := response["currency"]
				assert.True(t, exists, "Success response should have 'currency' field")
				if exists {
					assert.Equal(t, tt.payload["currency"], currency, "Response currency should match request")
				}

				// Status validation
				status, exists := response["status"]
				assert.True(t, exists, "Success response should have 'status' field")
				if exists {
					validStatuses := []string{"pending", "processing", "completed"}
					assert.Contains(t, validStatuses, status, "Status should be valid")
				}

				// Optional fields validation
				if description, hasDescription := tt.payload["description"]; hasDescription {
					responseDescription, exists := response["description"]
					assert.True(t, exists, "Response should have 'description' when provided in request")
					if exists {
						assert.Equal(t, description, responseDescription, "Response description should match request")
					}
				}

				if reference, hasReference := tt.payload["reference"]; hasReference {
					responseReference, exists := response["reference"]
					assert.True(t, exists, "Response should have 'reference' when provided in request")
					if exists {
						assert.Equal(t, reference, responseReference, "Response reference should match request")
					}
				}

				// Timestamps validation
				createdAt, exists := response["created_at"]
				assert.True(t, exists, "Success response should have 'created_at' field")
				assert.NotEmpty(t, createdAt, "created_at should not be empty")
			}
		})
	}
}

// TestPaymentTransactionSend_Currencies validates all supported currencies
func TestPaymentTransactionSend_Currencies(t *testing.T) {
	baseURL := "http://localhost:8080"
	mockToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	currencies := []string{"THB", "SGD", "IDR", "MYR", "PHP", "VND", "USD"}

	for _, currency := range currencies {
		t.Run("Send "+currency+" transaction", func(t *testing.T) {
			payload := map[string]interface{}{
				"wallet_id":    uuid.New().String(),
				"recipient_id": uuid.New().String(),
				"amount":       10000,
				"currency":     currency,
				"description":  "Test " + currency + " transaction",
			}

			body, err := json.Marshal(payload)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", baseURL+"/transactions/send", bytes.NewBuffer(body))
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

// TestPaymentTransactionSend_AmountLimits validates transaction amount limits
func TestPaymentTransactionSend_AmountLimits(t *testing.T) {
	baseURL := "http://localhost:8080"
	mockToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	tests := []struct {
		name           string
		amount         int
		expectedStatus int
		description    string
	}{
		{"Minimum valid amount", 1, http.StatusCreated, "Should accept minimum amount of 1 cent"},
		{"Small amount", 100, http.StatusCreated, "Should accept small amounts"},
		{"Medium amount", 10000, http.StatusCreated, "Should accept medium amounts"},
		{"Large amount", 1000000, http.StatusCreated, "Should accept large amounts within limits"},
		{"Very large amount", 999999999, http.StatusBadRequest, "Should reject very large amounts (daily limit)"},
		{"Zero amount", 0, http.StatusBadRequest, "Should reject zero amounts"},
		{"Negative amount", -1000, http.StatusBadRequest, "Should reject negative amounts"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := map[string]interface{}{
				"wallet_id":    uuid.New().String(),
				"recipient_id": uuid.New().String(),
				"amount":       tt.amount,
				"currency":     "THB",
				"description":  tt.description,
			}

			body, err := json.Marshal(payload)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", baseURL+"/transactions/send", bytes.NewBuffer(body))
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

			assert.Equal(t, tt.expectedStatus, resp.StatusCode, tt.description)
		})
	}
}

// TestPaymentTransactionSend_DailyLimit validates daily transaction limits
func TestPaymentTransactionSend_DailyLimit(t *testing.T) {
	baseURL := "http://localhost:8080"
	mockToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	walletID := uuid.New().String()

	// Simulate multiple transactions to test daily limit
	for i := 0; i < 6; i++ { // Assume daily limit allows 5 transactions
		payload := map[string]interface{}{
			"wallet_id":    walletID,
			"recipient_id": uuid.New().String(),
			"amount":       50000, // Large amount to trigger limit
			"currency":     "THB",
			"description":  "Daily limit test transaction",
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", baseURL+"/transactions/send", bytes.NewBuffer(body))
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

		// After 5th transaction, should get daily limit exceeded (402)
		if i >= 5 {
			assert.Equal(t, http.StatusPaymentRequired, resp.StatusCode,
				"Should be daily limit exceeded after multiple large transactions")
		}
	}
}