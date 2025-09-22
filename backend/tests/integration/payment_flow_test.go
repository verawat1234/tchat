package integration

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
	"github.com/stretchr/testify/suite"
)

// T031: Integration test payment processing flow
// Tests end-to-end payment workflow including:
// 1. Wallet creation → 2. KYC verification → 3. Top-up → 4. Transfer → 5. Balance verification → 6. Transaction history
type PaymentFlowTestSuite struct {
	suite.Suite
	router       *gin.Engine
	wallets      map[string]map[string]interface{}
	transactions map[string][]map[string]interface{}
	kycData      map[string]map[string]interface{}
}

func TestPaymentFlowSuite(t *testing.T) {
	suite.Run(t, new(PaymentFlowTestSuite))
}

func (suite *PaymentFlowTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.wallets = make(map[string]map[string]interface{})
	suite.transactions = make(map[string][]map[string]interface{})
	suite.kycData = make(map[string]map[string]interface{})

	// Setup payment service endpoints
	suite.setupPaymentEndpoints()
}

func (suite *PaymentFlowTestSuite) setupPaymentEndpoints() {
	// Mock authentication middleware
	authMiddleware := func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		if strings.HasPrefix(auth, "Bearer user_") {
			userID := strings.TrimPrefix(auth, "Bearer ")
			c.Set("user_id", userID)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			c.Abort()
			return
		}
		c.Next()
	}

	// KYC submission endpoint
	suite.router.POST("/users/kyc", authMiddleware, func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_json", "message": err.Error()})
			return
		}

		userID := c.GetString("user_id")

		// Validate required fields
		requiredFields := []string{"document_type", "document_number", "full_name", "date_of_birth", "country"}
		for _, field := range requiredFields {
			if _, exists := req[field]; !exists {
				c.JSON(http.StatusBadRequest, gin.H{"error": "missing_field", "field": field})
				return
			}
		}

		// Validate Southeast Asian countries
		validCountries := []string{"TH", "SG", "ID", "MY", "PH", "VN"}
		country := req["country"].(string)
		validCountry := false
		for _, vc := range validCountries {
			if country == vc {
				validCountry = true
				break
			}
		}

		if !validCountry {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_country", "message": "Country must be in Southeast Asia"})
			return
		}

		// Store KYC data
		kycID := fmt.Sprintf("kyc_%s_%d", userID, time.Now().Unix())
		kycRecord := map[string]interface{}{
			"id":              kycID,
			"user_id":         userID,
			"document_type":   req["document_type"],
			"document_number": req["document_number"],
			"full_name":       req["full_name"],
			"date_of_birth":   req["date_of_birth"],
			"country":         req["country"],
			"status":          "verified", // Auto-approve for testing
			"verified_at":     time.Now().UTC().Format(time.RFC3339),
			"created_at":      time.Now().UTC().Format(time.RFC3339),
		}

		suite.kycData[userID] = kycRecord

		c.JSON(http.StatusCreated, gin.H{
			"kyc_id":      kycID,
			"status":      "verified",
			"verified_at": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Create wallet endpoint
	suite.router.POST("/wallets", authMiddleware, func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_json", "message": err.Error()})
			return
		}

		userID := c.GetString("user_id")

		// Check KYC verification
		if _, kycExists := suite.kycData[userID]; !kycExists {
			c.JSON(http.StatusForbidden, gin.H{"error": "kyc_required", "message": "KYC verification required"})
			return
		}

		// Validate currency
		currency := req["currency"].(string)
		validCurrencies := []string{"THB", "USD", "SGD", "IDR", "MYR", "PHP", "VND"}
		validCurrency := false
		for _, vc := range validCurrencies {
			if currency == vc {
				validCurrency = true
				break
			}
		}

		if !validCurrency {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_currency"})
			return
		}

		walletID := fmt.Sprintf("wallet_%s_%s_%d", userID, currency, time.Now().Unix())

		wallet := map[string]interface{}{
			"id":                walletID,
			"user_id":           userID,
			"currency":          currency,
			"status":            "active",
			"available_balance": 0.0,
			"pending_balance":   0.0,
			"total_balance":     0.0,
			"created_at":        time.Now().UTC().Format(time.RFC3339),
			"updated_at":        time.Now().UTC().Format(time.RFC3339),
		}

		suite.wallets[walletID] = wallet
		suite.transactions[walletID] = make([]map[string]interface{}, 0)

		c.JSON(http.StatusCreated, wallet)
	})

	// Top-up endpoint
	suite.router.POST("/transactions/topup", authMiddleware, func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_json", "message": err.Error()})
			return
		}

		userID := c.GetString("user_id")
		walletID := req["wallet_id"].(string)
		amount := req["amount"].(float64)
		currency := req["currency"].(string)
		paymentMethod := req["payment_method"].(string)

		// Validate wallet ownership
		wallet, exists := suite.wallets[walletID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "wallet_not_found"})
			return
		}

		if wallet["user_id"] != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "wallet_access_denied"})
			return
		}

		if amount <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_amount"})
			return
		}

		// Calculate processing fee
		var processingFee float64
		switch paymentMethod {
		case "bank_transfer":
			processingFee = amount * 0.001 // 0.1%
		case "credit_card", "debit_card":
			processingFee = amount * 0.029 // 2.9%
		case "digital_wallet":
			processingFee = amount * 0.015 // 1.5%
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_payment_method"})
			return
		}

		// Create transaction
		txnID := fmt.Sprintf("topup_%d", time.Now().UnixNano())
		transaction := map[string]interface{}{
			"id":             txnID,
			"wallet_id":      walletID,
			"type":           "topup",
			"amount":         amount,
			"currency":       currency,
			"status":         "completed", // Auto-complete for testing
			"processing_fee": processingFee,
			"payment_method": paymentMethod,
			"description":    "Wallet top-up",
			"created_at":     time.Now().UTC().Format(time.RFC3339),
			"completed_at":   time.Now().UTC().Format(time.RFC3339),
		}

		// Update wallet balance
		wallet["available_balance"] = wallet["available_balance"].(float64) + amount
		wallet["total_balance"] = wallet["total_balance"].(float64) + amount
		wallet["updated_at"] = time.Now().UTC().Format(time.RFC3339)

		// Store transaction
		suite.transactions[walletID] = append(suite.transactions[walletID], transaction)

		c.JSON(http.StatusCreated, gin.H{
			"transaction_id": txnID,
			"status":         "completed",
			"amount":         amount,
			"processing_fee": processingFee,
			"new_balance":    wallet["total_balance"],
		})
	})

	// Transfer endpoint
	suite.router.POST("/transactions/transfer", authMiddleware, func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_json", "message": err.Error()})
			return
		}

		userID := c.GetString("user_id")
		fromWalletID := req["from_wallet"].(string)
		toWalletID := req["to_wallet"].(string)
		amount := req["amount"].(float64)

		// Validate source wallet
		fromWallet, exists := suite.wallets[fromWalletID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "source_wallet_not_found"})
			return
		}

		if fromWallet["user_id"] != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "source_wallet_access_denied"})
			return
		}

		// Validate destination wallet
		toWallet, exists := suite.wallets[toWalletID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "destination_wallet_not_found"})
			return
		}

		// Check currency match
		if fromWallet["currency"] != toWallet["currency"] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "currency_mismatch"})
			return
		}

		// Check sufficient balance
		if fromWallet["available_balance"].(float64) < amount {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient_balance"})
			return
		}

		if amount <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_amount"})
			return
		}

		// Calculate transfer fee
		transferFee := amount * 0.005 // 0.5% transfer fee
		totalDeduction := amount + transferFee

		if fromWallet["available_balance"].(float64) < totalDeduction {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient_balance_with_fee"})
			return
		}

		// Create transactions
		txnID := fmt.Sprintf("transfer_%d", time.Now().UnixNano())

		// Debit transaction
		debitTxn := map[string]interface{}{
			"id":           txnID + "_debit",
			"wallet_id":    fromWalletID,
			"type":         "transfer_out",
			"amount":       amount,
			"currency":     fromWallet["currency"],
			"status":       "completed",
			"to_wallet":    toWalletID,
			"transfer_fee": transferFee,
			"description":  "Transfer to " + toWalletID,
			"created_at":   time.Now().UTC().Format(time.RFC3339),
			"completed_at": time.Now().UTC().Format(time.RFC3339),
		}

		// Credit transaction
		creditTxn := map[string]interface{}{
			"id":           txnID + "_credit",
			"wallet_id":    toWalletID,
			"type":         "transfer_in",
			"amount":       amount,
			"currency":     toWallet["currency"],
			"status":       "completed",
			"from_wallet":  fromWalletID,
			"description":  "Transfer from " + fromWalletID,
			"created_at":   time.Now().UTC().Format(time.RFC3339),
			"completed_at": time.Now().UTC().Format(time.RFC3339),
		}

		// Update balances
		fromWallet["available_balance"] = fromWallet["available_balance"].(float64) - totalDeduction
		fromWallet["total_balance"] = fromWallet["total_balance"].(float64) - totalDeduction
		fromWallet["updated_at"] = time.Now().UTC().Format(time.RFC3339)

		toWallet["available_balance"] = toWallet["available_balance"].(float64) + amount
		toWallet["total_balance"] = toWallet["total_balance"].(float64) + amount
		toWallet["updated_at"] = time.Now().UTC().Format(time.RFC3339)

		// Store transactions
		suite.transactions[fromWalletID] = append(suite.transactions[fromWalletID], debitTxn)
		suite.transactions[toWalletID] = append(suite.transactions[toWalletID], creditTxn)

		c.JSON(http.StatusCreated, gin.H{
			"transaction_id":      txnID,
			"status":              "completed",
			"amount":              amount,
			"transfer_fee":        transferFee,
			"source_new_balance":  fromWallet["total_balance"],
			"dest_new_balance":    toWallet["total_balance"],
		})
	})

	// Get wallet balance endpoint
	suite.router.GET("/wallets/:id/balance", authMiddleware, func(c *gin.Context) {
		walletID := c.Param("id")
		userID := c.GetString("user_id")

		wallet, exists := suite.wallets[walletID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "wallet_not_found"})
			return
		}

		if wallet["user_id"] != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "wallet_access_denied"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"wallet_id":         walletID,
			"currency":          wallet["currency"],
			"available_balance": wallet["available_balance"],
			"pending_balance":   wallet["pending_balance"],
			"total_balance":     wallet["total_balance"],
			"updated_at":        wallet["updated_at"],
		})
	})

	// Get transaction history endpoint
	suite.router.GET("/wallets/:id/transactions", authMiddleware, func(c *gin.Context) {
		walletID := c.Param("id")
		userID := c.GetString("user_id")

		wallet, exists := suite.wallets[walletID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "wallet_not_found"})
			return
		}

		if wallet["user_id"] != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "wallet_access_denied"})
			return
		}

		transactions := suite.transactions[walletID]

		c.JSON(http.StatusOK, gin.H{
			"wallet_id":    walletID,
			"transactions": transactions,
			"total":        len(transactions),
		})
	})
}

func (suite *PaymentFlowTestSuite) TestCompletePaymentFlow() {
	userID := "user_123"

	// Step 1: KYC Verification
	suite.T().Log("Step 1: Submitting KYC verification")

	kycData := map[string]interface{}{
		"document_type":   "national_id",
		"document_number": "1234567890123",
		"full_name":       "Somchai Jaidee",
		"date_of_birth":   "1990-01-15",
		"country":         "TH",
	}

	jsonData, _ := json.Marshal(kycData)
	req := httptest.NewRequest("POST", "/users/kyc", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var kycResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &kycResponse)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), "verified", kycResponse["status"])
	assert.NotEmpty(suite.T(), kycResponse["kyc_id"])

	// Step 2: Create THB Wallet
	suite.T().Log("Step 2: Creating THB wallet")

	walletData := map[string]interface{}{
		"currency": "THB",
		"name":     "Primary THB Wallet",
	}

	jsonData, _ = json.Marshal(walletData)
	req = httptest.NewRequest("POST", "/wallets", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var walletResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &walletResponse)
	require.NoError(suite.T(), err)

	thbWalletID := walletResponse["id"].(string)
	assert.NotEmpty(suite.T(), thbWalletID)
	assert.Equal(suite.T(), "THB", walletResponse["currency"])
	assert.Equal(suite.T(), "active", walletResponse["status"])
	assert.Equal(suite.T(), 0.0, walletResponse["total_balance"])

	// Step 3: Create USD Wallet for Transfer Testing
	suite.T().Log("Step 3: Creating USD wallet")

	walletData["currency"] = "USD"
	jsonData, _ = json.Marshal(walletData)
	req = httptest.NewRequest("POST", "/wallets", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	json.Unmarshal(w.Body.Bytes(), &walletResponse)
	usdWalletID := walletResponse["id"].(string)

	// Step 4: Top-up THB Wallet
	suite.T().Log("Step 4: Topping up THB wallet")

	topupData := map[string]interface{}{
		"wallet_id":      thbWalletID,
		"amount":         5000.0,
		"currency":       "THB",
		"payment_method": "bank_transfer",
		"description":    "Initial funding",
	}

	jsonData, _ = json.Marshal(topupData)
	req = httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	start := time.Now()
	suite.router.ServeHTTP(w, req)
	topupDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	assert.True(suite.T(), topupDuration < 200*time.Millisecond, "Top-up should complete in <200ms")

	var topupResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &topupResponse)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), "completed", topupResponse["status"])
	assert.Equal(suite.T(), 5000.0, topupResponse["amount"])
	assert.Equal(suite.T(), 5.0, topupResponse["processing_fee"]) // 0.1% of 5000
	assert.Equal(suite.T(), 5000.0, topupResponse["new_balance"])

	// Step 5: Top-up USD Wallet
	suite.T().Log("Step 5: Topping up USD wallet")

	topupData = map[string]interface{}{
		"wallet_id":      usdWalletID,
		"amount":         100.0,
		"currency":       "USD",
		"payment_method": "credit_card",
		"description":    "USD funding",
	}

	jsonData, _ = json.Marshal(topupData)
	req = httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	json.Unmarshal(w.Body.Bytes(), &topupResponse)
	assert.Equal(suite.T(), 2.9, topupResponse["processing_fee"]) // 2.9% of 100

	// Step 6: Verify Balances
	suite.T().Log("Step 6: Verifying wallet balances")

	// Check THB wallet balance
	req = httptest.NewRequest("GET", fmt.Sprintf("/wallets/%s/balance", thbWalletID), nil)
	req.Header.Set("Authorization", "Bearer "+userID)

	w = httptest.NewRecorder()
	start = time.Now()
	suite.router.ServeHTTP(w, req)
	balanceDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.True(suite.T(), balanceDuration < 100*time.Millisecond, "Balance check should be <100ms")

	var balanceResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &balanceResponse)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), 5000.0, balanceResponse["total_balance"])
	assert.Equal(suite.T(), 5000.0, balanceResponse["available_balance"])
	assert.Equal(suite.T(), "THB", balanceResponse["currency"])

	// Step 7: Create Second User and Wallet for Transfer Testing
	suite.T().Log("Step 7: Setting up transfer recipient")

	secondUserID := "user_456"

	// KYC for second user
	kycData["full_name"] = "Jane Doe"
	kycData["document_number"] = "9876543210987"
	jsonData, _ = json.Marshal(kycData)
	req = httptest.NewRequest("POST", "/users/kyc", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+secondUserID)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	// Create THB wallet for second user
	walletData["currency"] = "THB"
	jsonData, _ = json.Marshal(walletData)
	req = httptest.NewRequest("POST", "/wallets", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+secondUserID)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &walletResponse)
	recipientWalletID := walletResponse["id"].(string)

	// Step 8: Transfer Money
	suite.T().Log("Step 8: Transferring money between wallets")

	transferData := map[string]interface{}{
		"from_wallet": thbWalletID,
		"to_wallet":   recipientWalletID,
		"amount":      1000.0,
		"description": "Payment to friend",
	}

	jsonData, _ = json.Marshal(transferData)
	req = httptest.NewRequest("POST", "/transactions/transfer", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	start = time.Now()
	suite.router.ServeHTTP(w, req)
	transferDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	assert.True(suite.T(), transferDuration < 300*time.Millisecond, "Transfer should complete in <300ms")

	var transferResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &transferResponse)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), "completed", transferResponse["status"])
	assert.Equal(suite.T(), 1000.0, transferResponse["amount"])
	assert.Equal(suite.T(), 5.0, transferResponse["transfer_fee"]) // 0.5% of 1000
	assert.Equal(suite.T(), 3995.0, transferResponse["source_new_balance"]) // 5000 - 1000 - 5
	assert.Equal(suite.T(), 1000.0, transferResponse["dest_new_balance"])

	// Step 9: Verify Final Balances
	suite.T().Log("Step 9: Verifying final balances")

	// Check sender balance
	req = httptest.NewRequest("GET", fmt.Sprintf("/wallets/%s/balance", thbWalletID), nil)
	req.Header.Set("Authorization", "Bearer "+userID)

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &balanceResponse)
	assert.Equal(suite.T(), 3995.0, balanceResponse["total_balance"])

	// Check recipient balance
	req = httptest.NewRequest("GET", fmt.Sprintf("/wallets/%s/balance", recipientWalletID), nil)
	req.Header.Set("Authorization", "Bearer "+secondUserID)

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &balanceResponse)
	assert.Equal(suite.T(), 1000.0, balanceResponse["total_balance"])

	// Step 10: Verify Transaction History
	suite.T().Log("Step 10: Verifying transaction history")

	req = httptest.NewRequest("GET", fmt.Sprintf("/wallets/%s/transactions", thbWalletID), nil)
	req.Header.Set("Authorization", "Bearer "+userID)

	w = httptest.NewRecorder()
	start = time.Now()
	suite.router.ServeHTTP(w, req)
	historyDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.True(suite.T(), historyDuration < 150*time.Millisecond, "Transaction history should load in <150ms")

	var historyResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &historyResponse)
	require.NoError(suite.T(), err)

	transactions := historyResponse["transactions"].([]interface{})
	assert.Equal(suite.T(), 2, len(transactions)) // Top-up + Transfer out
	assert.Equal(suite.T(), 2, int(historyResponse["total"].(float64)))

	// Verify transaction types
	txnTypes := make(map[string]bool)
	for _, txn := range transactions {
		transaction := txn.(map[string]interface{})
		txnTypes[transaction["type"].(string)] = true
	}
	assert.True(suite.T(), txnTypes["topup"])
	assert.True(suite.T(), txnTypes["transfer_out"])
}

func (suite *PaymentFlowTestSuite) TestPaymentErrorHandling() {
	userID := "user_123"

	// Test wallet creation without KYC
	suite.T().Log("Testing wallet creation without KYC")

	walletData := map[string]interface{}{
		"currency": "THB",
	}

	jsonData, _ := json.Marshal(walletData)
	req := httptest.NewRequest("POST", "/wallets", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)

	// Test transfer with insufficient balance
	suite.T().Log("Testing transfer with insufficient balance")

	// First complete KYC and create wallets
	kycData := map[string]interface{}{
		"document_type":   "national_id",
		"document_number": "1234567890123",
		"full_name":       "Test User",
		"date_of_birth":   "1990-01-15",
		"country":         "TH",
	}

	jsonData, _ = json.Marshal(kycData)
	req = httptest.NewRequest("POST", "/users/kyc", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Create wallets
	jsonData, _ = json.Marshal(walletData)
	req = httptest.NewRequest("POST", "/wallets", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var walletResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &walletResponse)
	walletID1 := walletResponse["id"].(string)

	req = httptest.NewRequest("POST", "/wallets", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &walletResponse)
	walletID2 := walletResponse["id"].(string)

	// Try transfer without sufficient balance
	transferData := map[string]interface{}{
		"from_wallet": walletID1,
		"to_wallet":   walletID2,
		"amount":      1000.0,
	}

	jsonData, _ = json.Marshal(transferData)
	req = httptest.NewRequest("POST", "/transactions/transfer", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var errorResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.Equal(suite.T(), "insufficient_balance", errorResponse["error"])
}

func (suite *PaymentFlowTestSuite) TestMultiCurrencySupport() {
	userID := "user_multicurrency"

	// Complete KYC
	kycData := map[string]interface{}{
		"document_type":   "national_id",
		"document_number": "MC123456789",
		"full_name":       "Multi Currency User",
		"date_of_birth":   "1985-06-20",
		"country":         "SG",
	}

	jsonData, _ := json.Marshal(kycData)
	req := httptest.NewRequest("POST", "/users/kyc", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Test all Southeast Asian currencies
	currencies := []string{"THB", "USD", "SGD", "IDR", "MYR", "PHP", "VND"}
	walletIDs := make(map[string]string)

	for _, currency := range currencies {
		suite.T().Logf("Testing %s currency support", currency)

		walletData := map[string]interface{}{
			"currency": currency,
		}

		jsonData, _ = json.Marshal(walletData)
		req = httptest.NewRequest("POST", "/wallets", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+userID)
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusCreated, w.Code, "Failed to create %s wallet", currency)

		var walletResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &walletResponse)
		walletIDs[currency] = walletResponse["id"].(string)

		// Test top-up for each currency
		var amount float64
		switch currency {
		case "THB":
			amount = 1000.0
		case "USD":
			amount = 50.0
		case "SGD":
			amount = 70.0
		case "IDR":
			amount = 750000.0
		case "MYR":
			amount = 200.0
		case "PHP":
			amount = 2500.0
		case "VND":
			amount = 1200000.0
		}

		topupData := map[string]interface{}{
			"wallet_id":      walletIDs[currency],
			"amount":         amount,
			"currency":       currency,
			"payment_method": "bank_transfer",
		}

		jsonData, _ = json.Marshal(topupData)
		req = httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+userID)
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusCreated, w.Code, "Failed to top-up %s wallet", currency)
	}
}

func (suite *PaymentFlowTestSuite) TestPaymentPerformance() {
	// Test high-frequency transaction processing
	suite.T().Log("Testing payment system performance")

	userID := "user_performance"

	// Setup user and wallet
	kycData := map[string]interface{}{
		"document_type":   "national_id",
		"document_number": "PERF123456789",
		"full_name":       "Performance Test User",
		"date_of_birth":   "1990-01-01",
		"country":         "TH",
	}

	jsonData, _ := json.Marshal(kycData)
	req := httptest.NewRequest("POST", "/users/kyc", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	walletData := map[string]interface{}{
		"currency": "THB",
	}

	jsonData, _ = json.Marshal(walletData)
	req = httptest.NewRequest("POST", "/wallets", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var walletResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &walletResponse)
	walletID := walletResponse["id"].(string)

	// Test multiple rapid top-ups
	topupCount := 10
	var totalDuration time.Duration

	for i := 0; i < topupCount; i++ {
		topupData := map[string]interface{}{
			"wallet_id":      walletID,
			"amount":         100.0,
			"currency":       "THB",
			"payment_method": "digital_wallet",
		}

		jsonData, _ = json.Marshal(topupData)
		req = httptest.NewRequest("POST", "/transactions/topup", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+userID)
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		start := time.Now()
		suite.router.ServeHTTP(w, req)
		duration := time.Since(start)
		totalDuration += duration

		assert.Equal(suite.T(), http.StatusCreated, w.Code)
		assert.True(suite.T(), duration < 200*time.Millisecond, "Each top-up should complete in <200ms")
	}

	avgDuration := totalDuration / time.Duration(topupCount)
	assert.True(suite.T(), avgDuration < 150*time.Millisecond, "Average top-up time should be <150ms")

	// Verify final balance calculation
	req = httptest.NewRequest("GET", fmt.Sprintf("/wallets/%s/balance", walletID), nil)
	req.Header.Set("Authorization", "Bearer "+userID)

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var balanceResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &balanceResponse)
	expectedBalance := float64(topupCount) * 100.0
	assert.Equal(suite.T(), expectedBalance, balanceResponse["total_balance"])

	// Test transaction history performance
	req = httptest.NewRequest("GET", fmt.Sprintf("/wallets/%s/transactions", walletID), nil)
	req.Header.Set("Authorization", "Bearer "+userID)

	w = httptest.NewRecorder()
	start := time.Now()
	suite.router.ServeHTTP(w, req)
	historyDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.True(suite.T(), historyDuration < 200*time.Millisecond, "Transaction history should load in <200ms with %d transactions", topupCount)

	var historyResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &historyResponse)
	assert.Equal(suite.T(), topupCount, int(historyResponse["total"].(float64)))
}