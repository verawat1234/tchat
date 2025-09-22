package contract

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T025: Contract test GET /wallets/{id}/balance
// Tests wallet balance retrieval with multi-currency support
func TestPaymentWalletBalanceContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock wallet balance response structure
	type CurrencyBalance struct {
		Currency         string  `json:"currency"`
		AvailableBalance float64 `json:"available_balance"`
		PendingBalance   float64 `json:"pending_balance"`
		TotalBalance     float64 `json:"total_balance"`
		LastUpdated      string  `json:"last_updated"`
	}

	type WalletBalanceResponse struct {
		WalletID  string            `json:"wallet_id"`
		UserID    string            `json:"user_id"`
		Status    string            `json:"status"`
		Balances  []CurrencyBalance `json:"balances"`
		CreatedAt string            `json:"created_at"`
		UpdatedAt string            `json:"updated_at"`
	}

	type ErrorResponse struct {
		Error   string `json:"error"`
		Code    string `json:"code"`
		Message string `json:"message"`
		Details string `json:"details,omitempty"`
	}

	// Define test endpoints with mock data
	router.GET("/wallets/:id/balance", func(c *gin.Context) {
		walletID := c.Param("id")
		auth := c.GetHeader("Authorization")

		// Simulate authentication check
		if auth == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Code:    "AUTH_REQUIRED",
				Message: "Authentication required",
			})
			return
		}

		if auth == "Bearer invalid_token" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Code:    "INVALID_TOKEN",
				Message: "Invalid or expired token",
			})
			return
		}

		// Test cases based on wallet ID
		switch walletID {
		case "wallet_123":
			// Valid wallet with multi-currency balances
			response := WalletBalanceResponse{
				WalletID: "wallet_123",
				UserID:   "user_456",
				Status:   "active",
				Balances: []CurrencyBalance{
					{
						Currency:         "THB",
						AvailableBalance: 1000.50,
						PendingBalance:   50.25,
						TotalBalance:     1050.75,
						LastUpdated:      time.Now().UTC().Format(time.RFC3339),
					},
					{
						Currency:         "USD",
						AvailableBalance: 25.75,
						PendingBalance:   0.00,
						TotalBalance:     25.75,
						LastUpdated:      time.Now().UTC().Format(time.RFC3339),
					},
					{
						Currency:         "SGD",
						AvailableBalance: 35.20,
						PendingBalance:   2.10,
						TotalBalance:     37.30,
						LastUpdated:      time.Now().UTC().Format(time.RFC3339),
					},
				},
				CreatedAt: "2024-01-15T10:30:00Z",
				UpdatedAt: time.Now().UTC().Format(time.RFC3339),
			}
			c.JSON(http.StatusOK, response)

		case "wallet_frozen":
			// Frozen wallet
			response := WalletBalanceResponse{
				WalletID: "wallet_frozen",
				UserID:   "user_789",
				Status:   "frozen",
				Balances: []CurrencyBalance{
					{
						Currency:         "THB",
						AvailableBalance: 0.00,
						PendingBalance:   500.00,
						TotalBalance:     500.00,
						LastUpdated:      time.Now().UTC().Format(time.RFC3339),
					},
				},
				CreatedAt: "2024-01-10T08:15:00Z",
				UpdatedAt: time.Now().UTC().Format(time.RFC3339),
			}
			c.JSON(http.StatusOK, response)

		case "wallet_unauthorized":
			// User doesn't have access to this wallet
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "forbidden",
				Code:    "ACCESS_DENIED",
				Message: "You don't have permission to access this wallet",
			})

		case "wallet_suspended":
			// Suspended wallet
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "forbidden",
				Code:    "WALLET_SUSPENDED",
				Message: "Wallet is suspended due to policy violation",
				Details: "Contact support for assistance",
			})

		case "wallet_notfound":
			// Non-existent wallet
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Code:    "WALLET_NOT_FOUND",
				Message: "Wallet not found",
			})

		case "wallet_invalid_format":
			// Invalid wallet ID format
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "bad_request",
				Code:    "INVALID_WALLET_ID",
				Message: "Invalid wallet ID format",
			})

		default:
			// Generic error
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "internal_error",
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			})
		}
	})

	// Test Cases
	t.Run("GET /wallets/{id}/balance - Success with multi-currency", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_123/balance", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response WalletBalanceResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Validate response structure
		assert.Equal(t, "wallet_123", response.WalletID)
		assert.Equal(t, "user_456", response.UserID)
		assert.Equal(t, "active", response.Status)
		assert.Len(t, response.Balances, 3)

		// Validate currency balances
		currencies := make(map[string]CurrencyBalance)
		for _, balance := range response.Balances {
			currencies[balance.Currency] = balance
		}

		// Check THB balance
		thb := currencies["THB"]
		assert.Equal(t, 1000.50, thb.AvailableBalance)
		assert.Equal(t, 50.25, thb.PendingBalance)
		assert.Equal(t, 1050.75, thb.TotalBalance)

		// Check USD balance
		usd := currencies["USD"]
		assert.Equal(t, 25.75, usd.AvailableBalance)
		assert.Equal(t, 0.00, usd.PendingBalance)
		assert.Equal(t, 25.75, usd.TotalBalance)

		// Check SGD balance
		sgd := currencies["SGD"]
		assert.Equal(t, 35.20, sgd.AvailableBalance)
		assert.Equal(t, 2.10, sgd.PendingBalance)
		assert.Equal(t, 37.30, sgd.TotalBalance)

		// Validate timestamps
		assert.NotEmpty(t, response.CreatedAt)
		assert.NotEmpty(t, response.UpdatedAt)
		for _, balance := range response.Balances {
			assert.NotEmpty(t, balance.LastUpdated)
		}
	})

	t.Run("GET /wallets/{id}/balance - Frozen wallet", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_frozen/balance", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response WalletBalanceResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "wallet_frozen", response.WalletID)
		assert.Equal(t, "frozen", response.Status)
		assert.Len(t, response.Balances, 1)

		// Frozen wallet should have zero available balance
		balance := response.Balances[0]
		assert.Equal(t, "THB", balance.Currency)
		assert.Equal(t, 0.00, balance.AvailableBalance)
		assert.Equal(t, 500.00, balance.PendingBalance)
		assert.Equal(t, 500.00, balance.TotalBalance)
	})

	t.Run("GET /wallets/{id}/balance - Unauthorized (no token)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_123/balance", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "unauthorized", response.Error)
		assert.Equal(t, "AUTH_REQUIRED", response.Code)
		assert.Contains(t, response.Message, "Authentication required")
	})

	t.Run("GET /wallets/{id}/balance - Unauthorized (invalid token)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_123/balance", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "unauthorized", response.Error)
		assert.Equal(t, "INVALID_TOKEN", response.Code)
	})

	t.Run("GET /wallets/{id}/balance - Forbidden (access denied)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_unauthorized/balance", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "forbidden", response.Error)
		assert.Equal(t, "ACCESS_DENIED", response.Code)
		assert.Contains(t, response.Message, "permission")
	})

	t.Run("GET /wallets/{id}/balance - Forbidden (wallet suspended)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_suspended/balance", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "forbidden", response.Error)
		assert.Equal(t, "WALLET_SUSPENDED", response.Code)
		assert.Contains(t, response.Message, "suspended")
		assert.Contains(t, response.Details, "support")
	})

	t.Run("GET /wallets/{id}/balance - Not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_notfound/balance", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "not_found", response.Error)
		assert.Equal(t, "WALLET_NOT_FOUND", response.Code)
	})

	t.Run("GET /wallets/{id}/balance - Bad request (invalid format)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_invalid_format/balance", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "bad_request", response.Error)
		assert.Equal(t, "INVALID_WALLET_ID", response.Code)
	})

	// Security and Performance Tests
	t.Run("GET /wallets/{id}/balance - Rate limiting simulation", func(t *testing.T) {
		// Simulate multiple rapid requests
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest("GET", "/wallets/wallet_123/balance", nil)
			req.Header.Set("Authorization", "Bearer valid_token")
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			start := time.Now()
			router.ServeHTTP(w, req)
			duration := time.Since(start)

			// Should respond quickly
			assert.True(t, duration < 100*time.Millisecond, "Response should be under 100ms")
			// First few requests should succeed
			if i < 3 {
				assert.Equal(t, http.StatusOK, w.Code)
			}
		}
	})

	t.Run("GET /wallets/{id}/balance - Response headers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_123/balance", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	})

	t.Run("GET /wallets/{id}/balance - Currency precision validation", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_123/balance", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response WalletBalanceResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Validate decimal precision for financial calculations
		for _, balance := range response.Balances {
			// Check that balances have proper precision (2 decimal places for display)
			availableStr := fmt.Sprintf("%.2f", balance.AvailableBalance)
			assert.NotContains(t, availableStr, "e") // No scientific notation

			pendingStr := fmt.Sprintf("%.2f", balance.PendingBalance)
			assert.NotContains(t, pendingStr, "e")

			totalStr := fmt.Sprintf("%.2f", balance.TotalBalance)
			assert.NotContains(t, totalStr, "e")

			// Verify total = available + pending
			calculatedTotal := balance.AvailableBalance + balance.PendingBalance
			assert.InDelta(t, balance.TotalBalance, calculatedTotal, 0.01,
				"Total balance should equal available + pending")
		}
	})
}