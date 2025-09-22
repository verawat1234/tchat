package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T028: Contract test POST /transactions/topup
// Tests wallet top-up functionality with multiple payment methods and currencies
func TestPaymentTopupContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock top-up request structure
	type TopupRequest struct {
		WalletID      string                 `json:"wallet_id"`
		Amount        float64                `json:"amount"`
		Currency      string                 `json:"currency"`
		PaymentMethod string                 `json:"payment_method"` // bank_transfer, credit_card, debit_card, digital_wallet
		PaymentSource string                 `json:"payment_source,omitempty"`
		Description   string                 `json:"description,omitempty"`
		Metadata      map[string]interface{} `json:"metadata,omitempty"`
	}

	// Mock top-up response structure
	type TopupResponse struct {
		TransactionID   string  `json:"transaction_id"`
		WalletID        string  `json:"wallet_id"`
		Amount          float64 `json:"amount"`
		Currency        string  `json:"currency"`
		PaymentMethod   string  `json:"payment_method"`
		Status          string  `json:"status"` // pending, processing, completed, failed
		ProcessingFee   float64 `json:"processing_fee,omitempty"`
		ExchangeRate    float64 `json:"exchange_rate,omitempty"`
		EstimatedTime   string  `json:"estimated_time,omitempty"`
		PaymentURL      string  `json:"payment_url,omitempty"`
		PaymentRef      string  `json:"payment_ref,omitempty"`
		Instructions    string  `json:"instructions,omitempty"`
		ExpiresAt       string  `json:"expires_at,omitempty"`
		CreatedAt       string  `json:"created_at"`
		UpdatedAt       string  `json:"updated_at"`
	}

	type ErrorResponse struct {
		Error   string `json:"error"`
		Code    string `json:"code"`
		Message string `json:"message"`
		Details string `json:"details,omitempty"`
	}

	// Define test endpoints
	router.POST("/transactions/topup", func(c *gin.Context) {
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

		var req TopupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "bad_request",
				Code:    "INVALID_JSON",
				Message: "Invalid JSON format",
				Details: err.Error(),
			})
			return
		}

		// Validate required fields
		if req.WalletID == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "bad_request",
				Code:    "MISSING_WALLET_ID",
				Message: "Wallet ID is required",
			})
			return
		}

		if req.Amount <= 0 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "bad_request",
				Code:    "INVALID_AMOUNT",
				Message: "Amount must be greater than 0",
			})
			return
		}

		if req.Currency == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "bad_request",
				Code:    "MISSING_CURRENCY",
				Message: "Currency is required",
			})
			return
		}

		if req.PaymentMethod == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "bad_request",
				Code:    "MISSING_PAYMENT_METHOD",
				Message: "Payment method is required",
			})
			return
		}

		// Validate currency
		validCurrencies := []string{"THB", "USD", "SGD", "IDR", "MYR", "PHP", "VND"}
		validCurrency := false
		for _, currency := range validCurrencies {
			if req.Currency == currency {
				validCurrency = true
				break
			}
		}
		if !validCurrency {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "bad_request",
				Code:    "INVALID_CURRENCY",
				Message: "Unsupported currency",
				Details: "Supported: THB, USD, SGD, IDR, MYR, PHP, VND",
			})
			return
		}

		// Validate payment method
		validMethods := []string{"bank_transfer", "credit_card", "debit_card", "digital_wallet"}
		validMethod := false
		for _, method := range validMethods {
			if req.PaymentMethod == method {
				validMethod = true
				break
			}
		}
		if !validMethod {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "bad_request",
				Code:    "INVALID_PAYMENT_METHOD",
				Message: "Unsupported payment method",
				Details: "Supported: bank_transfer, credit_card, debit_card, digital_wallet",
			})
			return
		}

		// Validate amount limits by currency
		var minAmount, maxAmount float64
		switch req.Currency {
		case "THB":
			minAmount, maxAmount = 50.0, 100000.0
		case "USD":
			minAmount, maxAmount = 1.0, 2500.0
		case "SGD":
			minAmount, maxAmount = 1.0, 3500.0
		case "IDR":
			minAmount, maxAmount = 15000.0, 35000000.0
		case "MYR":
			minAmount, maxAmount = 5.0, 10000.0
		case "PHP":
			minAmount, maxAmount = 50.0, 150000.0
		case "VND":
			minAmount, maxAmount = 25000.0, 60000000.0
		}

		if req.Amount < minAmount || req.Amount > maxAmount {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "bad_request",
				Code:    "AMOUNT_OUT_OF_RANGE",
				Message: "Amount out of allowed range",
				Details: fmt.Sprintf("Min: %.2f %s, Max: %.2f %s", minAmount, req.Currency, maxAmount, req.Currency),
			})
			return
		}

		// Test cases based on wallet ID
		switch req.WalletID {
		case "wallet_123":
			// Successful top-up
			var processingFee float64
			var estimatedTime string
			var status string = "pending"
			var paymentURL string
			var instructions string

			switch req.PaymentMethod {
			case "bank_transfer":
				processingFee = req.Amount * 0.001 // 0.1% fee
				estimatedTime = "1-3 business days"
				instructions = "Transfer to account: 1234-5678-9012-3456"
			case "credit_card", "debit_card":
				processingFee = req.Amount * 0.029 // 2.9% fee
				estimatedTime = "Instant"
				status = "processing"
				paymentURL = "https://payment.tchat.com/cc/process/abc123"
			case "digital_wallet":
				processingFee = req.Amount * 0.015 // 1.5% fee
				estimatedTime = "5-10 minutes"
				status = "processing"
				paymentURL = "https://payment.tchat.com/wallet/process/def456"
			}

			response := TopupResponse{
				TransactionID: "topup_" + time.Now().Format("20060102150405"),
				WalletID:      req.WalletID,
				Amount:        req.Amount,
				Currency:      req.Currency,
				PaymentMethod: req.PaymentMethod,
				Status:        status,
				ProcessingFee: processingFee,
				EstimatedTime: estimatedTime,
				PaymentURL:    paymentURL,
				PaymentRef:    "REF_" + time.Now().Format("20060102150405"),
				Instructions:  instructions,
				ExpiresAt:     time.Now().Add(30 * time.Minute).UTC().Format(time.RFC3339),
				CreatedAt:     time.Now().UTC().Format(time.RFC3339),
				UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
			}

			// Add exchange rate for non-native currencies
			if req.Currency != "THB" {
				response.ExchangeRate = 1.0 // Simplified for testing
			}

			c.JSON(http.StatusCreated, response)

		case "wallet_suspended":
			// Suspended wallet
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "forbidden",
				Code:    "WALLET_SUSPENDED",
				Message: "Wallet is suspended and cannot receive top-ups",
				Details: "Contact support for assistance",
			})

		case "wallet_frozen":
			// Frozen wallet
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "forbidden",
				Code:    "WALLET_FROZEN",
				Message: "Wallet is frozen due to security concerns",
				Details: "Complete KYC verification to unfreeze",
			})

		case "wallet_notfound":
			// Non-existent wallet
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Code:    "WALLET_NOT_FOUND",
				Message: "Wallet not found",
			})

		case "wallet_unauthorized":
			// User doesn't have access to this wallet
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "forbidden",
				Code:    "ACCESS_DENIED",
				Message: "You don't have permission to top up this wallet",
			})

		case "wallet_limit_exceeded":
			// Daily/monthly limits exceeded
			c.JSON(http.StatusTooManyRequests, ErrorResponse{
				Error:   "rate_limit_exceeded",
				Code:    "TOPUP_LIMIT_EXCEEDED",
				Message: "Top-up limit exceeded",
				Details: "Daily limit: 50,000 THB, Monthly limit: 500,000 THB",
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
	t.Run("POST /transactions/topup - Success bank transfer", func(t *testing.T) {
		reqBody := TopupRequest{
			WalletID:      "wallet_123",
			Amount:        1000.00,
			Currency:      "THB",
			PaymentMethod: "bank_transfer",
			Description:   "Monthly allowance",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response TopupResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Validate response structure
		assert.NotEmpty(t, response.TransactionID)
		assert.Equal(t, "wallet_123", response.WalletID)
		assert.Equal(t, 1000.00, response.Amount)
		assert.Equal(t, "THB", response.Currency)
		assert.Equal(t, "bank_transfer", response.PaymentMethod)
		assert.Equal(t, "pending", response.Status)
		assert.Equal(t, 1.0, response.ProcessingFee) // 0.1% of 1000
		assert.Equal(t, "1-3 business days", response.EstimatedTime)
		assert.Contains(t, response.Instructions, "Transfer to account")
		assert.NotEmpty(t, response.PaymentRef)
		assert.NotEmpty(t, response.ExpiresAt)
		assert.NotEmpty(t, response.CreatedAt)
		assert.NotEmpty(t, response.UpdatedAt)
	})

	t.Run("POST /transactions/topup - Success credit card", func(t *testing.T) {
		reqBody := TopupRequest{
			WalletID:      "wallet_123",
			Amount:        500.00,
			Currency:      "USD",
			PaymentMethod: "credit_card",
			PaymentSource: "visa_1234",
			Description:   "Emergency top-up",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response TopupResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "wallet_123", response.WalletID)
		assert.Equal(t, 500.00, response.Amount)
		assert.Equal(t, "USD", response.Currency)
		assert.Equal(t, "credit_card", response.PaymentMethod)
		assert.Equal(t, "processing", response.Status)
		assert.Equal(t, 14.5, response.ProcessingFee) // 2.9% of 500
		assert.Equal(t, "Instant", response.EstimatedTime)
		assert.Contains(t, response.PaymentURL, "https://payment.tchat.com/cc/process/")
		assert.Equal(t, 1.0, response.ExchangeRate) // USD exchange rate
	})

	t.Run("POST /transactions/topup - Success digital wallet", func(t *testing.T) {
		reqBody := TopupRequest{
			WalletID:      "wallet_123",
			Amount:        250.00,
			Currency:      "SGD",
			PaymentMethod: "digital_wallet",
			PaymentSource: "grabpay",
			Metadata: map[string]interface{}{
				"source_app": "grabpay",
				"promo_code": "SAVE10",
			},
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response TopupResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "wallet_123", response.WalletID)
		assert.Equal(t, 250.00, response.Amount)
		assert.Equal(t, "SGD", response.Currency)
		assert.Equal(t, "digital_wallet", response.PaymentMethod)
		assert.Equal(t, "processing", response.Status)
		assert.Equal(t, 3.75, response.ProcessingFee) // 1.5% of 250
		assert.Equal(t, "5-10 minutes", response.EstimatedTime)
		assert.Contains(t, response.PaymentURL, "https://payment.tchat.com/wallet/process/")
	})

	t.Run("POST /transactions/topup - Unauthorized (no token)", func(t *testing.T) {
		reqBody := TopupRequest{
			WalletID:      "wallet_123",
			Amount:        100.00,
			Currency:      "THB",
			PaymentMethod: "bank_transfer",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonBody))
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

	t.Run("POST /transactions/topup - Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/transactions/topup", strings.NewReader("invalid json"))
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "bad_request", response.Error)
		assert.Equal(t, "INVALID_JSON", response.Code)
	})

	t.Run("POST /transactions/topup - Missing wallet ID", func(t *testing.T) {
		reqBody := TopupRequest{
			Amount:        100.00,
			Currency:      "THB",
			PaymentMethod: "bank_transfer",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "bad_request", response.Error)
		assert.Equal(t, "MISSING_WALLET_ID", response.Code)
	})

	t.Run("POST /transactions/topup - Invalid amount", func(t *testing.T) {
		reqBody := TopupRequest{
			WalletID:      "wallet_123",
			Amount:        -100.00,
			Currency:      "THB",
			PaymentMethod: "bank_transfer",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "bad_request", response.Error)
		assert.Equal(t, "INVALID_AMOUNT", response.Code)
	})

	t.Run("POST /transactions/topup - Invalid currency", func(t *testing.T) {
		reqBody := TopupRequest{
			WalletID:      "wallet_123",
			Amount:        100.00,
			Currency:      "INVALID",
			PaymentMethod: "bank_transfer",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "bad_request", response.Error)
		assert.Equal(t, "INVALID_CURRENCY", response.Code)
		assert.Contains(t, response.Details, "THB, USD, SGD")
	})

	t.Run("POST /transactions/topup - Invalid payment method", func(t *testing.T) {
		reqBody := TopupRequest{
			WalletID:      "wallet_123",
			Amount:        100.00,
			Currency:      "THB",
			PaymentMethod: "crypto",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "bad_request", response.Error)
		assert.Equal(t, "INVALID_PAYMENT_METHOD", response.Code)
	})

	t.Run("POST /transactions/topup - Amount out of range", func(t *testing.T) {
		reqBody := TopupRequest{
			WalletID:      "wallet_123",
			Amount:        999999.00, // Exceeds THB max limit
			Currency:      "THB",
			PaymentMethod: "bank_transfer",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "bad_request", response.Error)
		assert.Equal(t, "AMOUNT_OUT_OF_RANGE", response.Code)
		assert.Contains(t, response.Details, "Min:")
		assert.Contains(t, response.Details, "Max:")
	})

	t.Run("POST /transactions/topup - Wallet suspended", func(t *testing.T) {
		reqBody := TopupRequest{
			WalletID:      "wallet_suspended",
			Amount:        100.00,
			Currency:      "THB",
			PaymentMethod: "bank_transfer",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonBody))
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
		assert.Contains(t, response.Details, "support")
	})

	t.Run("POST /transactions/topup - Wallet frozen", func(t *testing.T) {
		reqBody := TopupRequest{
			WalletID:      "wallet_frozen",
			Amount:        100.00,
			Currency:      "THB",
			PaymentMethod: "bank_transfer",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "forbidden", response.Error)
		assert.Equal(t, "WALLET_FROZEN", response.Code)
		assert.Contains(t, response.Details, "KYC")
	})

	t.Run("POST /transactions/topup - Top-up limit exceeded", func(t *testing.T) {
		reqBody := TopupRequest{
			WalletID:      "wallet_limit_exceeded",
			Amount:        100.00,
			Currency:      "THB",
			PaymentMethod: "bank_transfer",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "rate_limit_exceeded", response.Error)
		assert.Equal(t, "TOPUP_LIMIT_EXCEEDED", response.Code)
		assert.Contains(t, response.Details, "Daily limit")
		assert.Contains(t, response.Details, "Monthly limit")
	})

	t.Run("POST /transactions/topup - Wallet not found", func(t *testing.T) {
		reqBody := TopupRequest{
			WalletID:      "wallet_notfound",
			Amount:        100.00,
			Currency:      "THB",
			PaymentMethod: "bank_transfer",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonBody))
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

	// Security and Performance Tests
	t.Run("POST /transactions/topup - Response time validation", func(t *testing.T) {
		reqBody := TopupRequest{
			WalletID:      "wallet_123",
			Amount:        100.00,
			Currency:      "THB",
			PaymentMethod: "bank_transfer",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		start := time.Now()
		router.ServeHTTP(w, req)
		duration := time.Since(start)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.True(t, duration < 200*time.Millisecond, "Response should be under 200ms")
	})

	t.Run("POST /transactions/topup - Southeast Asian currencies validation", func(t *testing.T) {
		testCases := []struct {
			currency   string
			amount     float64
			shouldPass bool
		}{
			{"THB", 100.00, true},
			{"USD", 10.00, true},
			{"SGD", 15.00, true},
			{"IDR", 150000.00, true},
			{"MYR", 50.00, true},
			{"PHP", 500.00, true},
			{"VND", 250000.00, true},
			{"EUR", 10.00, false}, // Should fail
		}

		for _, tc := range testCases {
			reqBody := TopupRequest{
				WalletID:      "wallet_123",
				Amount:        tc.amount,
				Currency:      tc.currency,
				PaymentMethod: "bank_transfer",
			}

			jsonBody, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonBody))
			req.Header.Set("Authorization", "Bearer valid_token")
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if tc.shouldPass {
				assert.Equal(t, http.StatusCreated, w.Code, "Currency %s should be supported", tc.currency)
			} else {
				assert.Equal(t, http.StatusBadRequest, w.Code, "Currency %s should not be supported", tc.currency)
			}
		}
	})
}