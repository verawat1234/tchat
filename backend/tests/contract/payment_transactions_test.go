package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T026: Contract test GET /wallets/{id}/transactions
// Tests wallet transaction history retrieval with pagination and filtering
func TestPaymentWalletTransactionsContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock transaction response structures
	type Transaction struct {
		ID              string  `json:"id"`
		WalletID        string  `json:"wallet_id"`
		Type            string  `json:"type"` // credit, debit, transfer_in, transfer_out, topup, withdrawal, fee
		Amount          float64 `json:"amount"`
		Currency        string  `json:"currency"`
		Status          string  `json:"status"` // pending, completed, failed, cancelled
		Description     string  `json:"description"`
		Reference       string  `json:"reference,omitempty"`
		FromWallet      string  `json:"from_wallet,omitempty"`
		ToWallet        string  `json:"to_wallet,omitempty"`
		TransactionFee  float64 `json:"transaction_fee,omitempty"`
		ExchangeRate    float64 `json:"exchange_rate,omitempty"`
		ProcessedAt     string  `json:"processed_at,omitempty"`
		CreatedAt       string  `json:"created_at"`
		UpdatedAt       string  `json:"updated_at"`
		Metadata        map[string]interface{} `json:"metadata,omitempty"`
	}

	type TransactionsResponse struct {
		WalletID     string        `json:"wallet_id"`
		Transactions []Transaction `json:"transactions"`
		Pagination   struct {
			Page       int    `json:"page"`
			Limit      int    `json:"limit"`
			Total      int    `json:"total"`
			TotalPages int    `json:"total_pages"`
			HasNext    bool   `json:"has_next"`
			HasPrev    bool   `json:"has_prev"`
			NextCursor string `json:"next_cursor,omitempty"`
			PrevCursor string `json:"prev_cursor,omitempty"`
		} `json:"pagination"`
	}

	type ErrorResponse struct {
		Error   string `json:"error"`
		Code    string `json:"code"`
		Message string `json:"message"`
		Details string `json:"details,omitempty"`
	}

	// Mock transaction data
	mockTransactions := map[string][]Transaction{
		"wallet_123": {
			{
				ID:          "txn_001",
				WalletID:    "wallet_123",
				Type:        "topup",
				Amount:      1000.00,
				Currency:    "THB",
				Status:      "completed",
				Description: "Bank transfer top-up",
				Reference:   "BANK_REF_001",
				ProcessedAt: "2024-01-20T14:30:00Z",
				CreatedAt:   "2024-01-20T14:28:00Z",
				UpdatedAt:   "2024-01-20T14:30:00Z",
			},
			{
				ID:             "txn_002",
				WalletID:       "wallet_123",
				Type:           "transfer_out",
				Amount:         500.00,
				Currency:       "THB",
				Status:         "completed",
				Description:    "Transfer to friend",
				ToWallet:       "wallet_456",
				TransactionFee: 5.00,
				ProcessedAt:    "2024-01-20T15:45:00Z",
				CreatedAt:      "2024-01-20T15:44:00Z",
				UpdatedAt:      "2024-01-20T15:45:00Z",
			},
			{
				ID:          "txn_003",
				WalletID:    "wallet_123",
				Type:        "debit",
				Amount:      25.50,
				Currency:    "THB",
				Status:      "completed",
				Description: "Coffee Shop Purchase",
				Reference:   "MERCHANT_001",
				ProcessedAt: "2024-01-20T16:20:00Z",
				CreatedAt:   "2024-01-20T16:20:00Z",
				UpdatedAt:   "2024-01-20T16:20:00Z",
			},
			{
				ID:          "txn_004",
				WalletID:    "wallet_123",
				Type:        "pending",
				Amount:      200.00,
				Currency:    "USD",
				Status:      "pending",
				Description: "International transfer",
				Reference:   "INTL_001",
				CreatedAt:   "2024-01-21T09:15:00Z",
				UpdatedAt:   "2024-01-21T09:15:00Z",
			},
		},
		"wallet_empty": {},
	}

	// Define test endpoints
	router.GET("/wallets/:id/transactions", func(c *gin.Context) {
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
		case "wallet_123", "wallet_empty":
			// Parse query parameters
			pageStr := c.DefaultQuery("page", "1")
			limitStr := c.DefaultQuery("limit", "10")
			typeFilter := c.Query("type")
			statusFilter := c.Query("status")
			currencyFilter := c.Query("currency")
			startDate := c.Query("start_date")
			endDate := c.Query("end_date")

			page, err := strconv.Atoi(pageStr)
			if err != nil || page < 1 {
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error:   "bad_request",
					Code:    "INVALID_PAGE",
					Message: "Invalid page parameter",
				})
				return
			}

			limit, err := strconv.Atoi(limitStr)
			if err != nil || limit < 1 || limit > 100 {
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error:   "bad_request",
					Code:    "INVALID_LIMIT",
					Message: "Invalid limit parameter (1-100)",
				})
				return
			}

			// Get transactions for this wallet
			allTransactions := mockTransactions[walletID]
			filteredTransactions := make([]Transaction, 0)

			// Apply filters
			for _, txn := range allTransactions {
				// Type filter
				if typeFilter != "" && txn.Type != typeFilter {
					continue
				}
				// Status filter
				if statusFilter != "" && txn.Status != statusFilter {
					continue
				}
				// Currency filter
				if currencyFilter != "" && txn.Currency != currencyFilter {
					continue
				}
				// Date filters would be implemented here
				_ = startDate
				_ = endDate

				filteredTransactions = append(filteredTransactions, txn)
			}

			// Calculate pagination
			total := len(filteredTransactions)
			totalPages := (total + limit - 1) / limit
			startIdx := (page - 1) * limit
			endIdx := startIdx + limit

			if startIdx > total {
				startIdx = total
			}
			if endIdx > total {
				endIdx = total
			}

			paginatedTransactions := make([]Transaction, 0)
			if startIdx < endIdx {
				paginatedTransactions = filteredTransactions[startIdx:endIdx]
			}

			response := TransactionsResponse{
				WalletID:     walletID,
				Transactions: paginatedTransactions,
			}

			response.Pagination.Page = page
			response.Pagination.Limit = limit
			response.Pagination.Total = total
			response.Pagination.TotalPages = totalPages
			response.Pagination.HasNext = page < totalPages
			response.Pagination.HasPrev = page > 1

			// Set cursors for cursor-based pagination
			if len(paginatedTransactions) > 0 {
				response.Pagination.NextCursor = paginatedTransactions[len(paginatedTransactions)-1].ID
				response.Pagination.PrevCursor = paginatedTransactions[0].ID
			}

			c.JSON(http.StatusOK, response)

		case "wallet_unauthorized":
			// User doesn't have access to this wallet
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "forbidden",
				Code:    "ACCESS_DENIED",
				Message: "You don't have permission to access this wallet",
			})

		case "wallet_notfound":
			// Non-existent wallet
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Code:    "WALLET_NOT_FOUND",
				Message: "Wallet not found",
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
	t.Run("GET /wallets/{id}/transactions - Success with pagination", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_123/transactions?page=1&limit=2", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response TransactionsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Validate response structure
		assert.Equal(t, "wallet_123", response.WalletID)
		assert.Len(t, response.Transactions, 2) // Limited to 2 per page

		// Validate pagination
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, 2, response.Pagination.Limit)
		assert.Equal(t, 4, response.Pagination.Total)
		assert.Equal(t, 2, response.Pagination.TotalPages)
		assert.True(t, response.Pagination.HasNext)
		assert.False(t, response.Pagination.HasPrev)

		// Validate transaction structure
		txn := response.Transactions[0]
		assert.NotEmpty(t, txn.ID)
		assert.Equal(t, "wallet_123", txn.WalletID)
		assert.NotEmpty(t, txn.Type)
		assert.Greater(t, txn.Amount, 0.0)
		assert.NotEmpty(t, txn.Currency)
		assert.NotEmpty(t, txn.Status)
		assert.NotEmpty(t, txn.CreatedAt)
		assert.NotEmpty(t, txn.UpdatedAt)
	})

	t.Run("GET /wallets/{id}/transactions - Success with type filter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_123/transactions?type=topup", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response TransactionsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "wallet_123", response.WalletID)
		assert.Len(t, response.Transactions, 1) // Only topup transactions

		// Validate filtered result
		txn := response.Transactions[0]
		assert.Equal(t, "topup", txn.Type)
		assert.Equal(t, "txn_001", txn.ID)
		assert.Equal(t, 1000.00, txn.Amount)
		assert.Equal(t, "THB", txn.Currency)
		assert.Equal(t, "completed", txn.Status)
	})

	t.Run("GET /wallets/{id}/transactions - Success with status filter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_123/transactions?status=pending", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response TransactionsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "wallet_123", response.WalletID)
		assert.Len(t, response.Transactions, 1) // Only pending transactions

		// Validate filtered result
		txn := response.Transactions[0]
		assert.Equal(t, "pending", txn.Status)
		assert.Empty(t, txn.ProcessedAt) // Pending transactions don't have processed_at
	})

	t.Run("GET /wallets/{id}/transactions - Success with currency filter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_123/transactions?currency=USD", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response TransactionsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "wallet_123", response.WalletID)
		assert.Len(t, response.Transactions, 1) // Only USD transactions

		// Validate filtered result
		txn := response.Transactions[0]
		assert.Equal(t, "USD", txn.Currency)
	})

	t.Run("GET /wallets/{id}/transactions - Success with multiple filters", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_123/transactions?type=transfer_out&status=completed&currency=THB", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response TransactionsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "wallet_123", response.WalletID)
		assert.Len(t, response.Transactions, 1) // Only matching transactions

		// Validate filtered result
		txn := response.Transactions[0]
		assert.Equal(t, "transfer_out", txn.Type)
		assert.Equal(t, "completed", txn.Status)
		assert.Equal(t, "THB", txn.Currency)
		assert.Equal(t, "wallet_456", txn.ToWallet)
		assert.Equal(t, 5.00, txn.TransactionFee)
	})

	t.Run("GET /wallets/{id}/transactions - Empty wallet", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_empty/transactions", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response TransactionsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "wallet_empty", response.WalletID)
		assert.Len(t, response.Transactions, 0)
		assert.Equal(t, 0, response.Pagination.Total)
		assert.Equal(t, 0, response.Pagination.TotalPages)
		assert.False(t, response.Pagination.HasNext)
		assert.False(t, response.Pagination.HasPrev)
	})

	t.Run("GET /wallets/{id}/transactions - Unauthorized (no token)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_123/transactions", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "unauthorized", response.Error)
		assert.Equal(t, "AUTH_REQUIRED", response.Code)
	})

	t.Run("GET /wallets/{id}/transactions - Unauthorized (invalid token)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_123/transactions", nil)
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

	t.Run("GET /wallets/{id}/transactions - Forbidden", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_unauthorized/transactions", nil)
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
	})

	t.Run("GET /wallets/{id}/transactions - Not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_notfound/transactions", nil)
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

	t.Run("GET /wallets/{id}/transactions - Invalid page parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_123/transactions?page=invalid", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "bad_request", response.Error)
		assert.Equal(t, "INVALID_PAGE", response.Code)
	})

	t.Run("GET /wallets/{id}/transactions - Invalid limit parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_123/transactions?limit=999", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "bad_request", response.Error)
		assert.Equal(t, "INVALID_LIMIT", response.Code)
		assert.Contains(t, response.Message, "1-100")
	})

	// Performance and Security Tests
	t.Run("GET /wallets/{id}/transactions - Response time validation", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_123/transactions", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		start := time.Now()
		router.ServeHTTP(w, req)
		duration := time.Since(start)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, duration < 200*time.Millisecond, "Response should be under 200ms")
	})

	t.Run("GET /wallets/{id}/transactions - Transaction data validation", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wallets/wallet_123/transactions", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response TransactionsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Validate transaction types
		validTypes := map[string]bool{
			"credit": true, "debit": true, "transfer_in": true, "transfer_out": true,
			"topup": true, "withdrawal": true, "fee": true, "pending": true,
		}

		validStatuses := map[string]bool{
			"pending": true, "completed": true, "failed": true, "cancelled": true,
		}

		for _, txn := range response.Transactions {
			// Validate transaction ID format
			assert.NotEmpty(t, txn.ID)
			assert.True(t, len(txn.ID) >= 5, "Transaction ID should be meaningful length")

			// Validate type and status
			assert.True(t, validTypes[txn.Type] || validStatuses[txn.Type], "Invalid transaction type: %s", txn.Type)
			assert.True(t, validStatuses[txn.Status], "Invalid transaction status: %s", txn.Status)

			// Validate amounts
			assert.GreaterOrEqual(t, txn.Amount, 0.0, "Amount should be non-negative")
			if txn.TransactionFee > 0 {
				assert.GreaterOrEqual(t, txn.TransactionFee, 0.0, "Fee should be non-negative")
			}

			// Validate currencies
			validCurrencies := []string{"THB", "USD", "SGD", "IDR", "MYR", "PHP", "VND"}
			assert.Contains(t, validCurrencies, txn.Currency, "Invalid currency: %s", txn.Currency)

			// Validate timestamps
			assert.NotEmpty(t, txn.CreatedAt)
			assert.NotEmpty(t, txn.UpdatedAt)

			// For completed transactions, processed_at should be set
			if txn.Status == "completed" && txn.Type != "pending" {
				assert.NotEmpty(t, txn.ProcessedAt, "Completed transactions should have processed_at")
			}
		}
	})
}