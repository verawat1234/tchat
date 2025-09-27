package integration

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// PaymentIntegrationSuite tests the Payment service endpoints
type PaymentIntegrationSuite struct {
	APIIntegrationSuite
	ports ServicePort
}

// Wallet represents a user's wallet
type Wallet struct {
	ID        string            `json:"id"`
	UserID    string            `json:"userId"`
	Currency  string            `json:"currency"`
	Balance   float64           `json:"balance"`
	Status    string            `json:"status"`
	Type      string            `json:"type"`
	Settings  map[string]string `json:"settings"`
	CreatedAt string            `json:"createdAt"`
	UpdatedAt string            `json:"updatedAt"`
}

// Transaction represents a financial transaction
type Transaction struct {
	ID              string            `json:"id"`
	WalletID        string            `json:"walletId"`
	Type            string            `json:"type"`
	Amount          float64           `json:"amount"`
	Currency        string            `json:"currency"`
	Status          string            `json:"status"`
	Description     string            `json:"description"`
	Reference       string            `json:"reference"`
	Metadata        map[string]string `json:"metadata"`
	FromWalletID    *string           `json:"fromWalletId,omitempty"`
	ToWalletID      *string           `json:"toWalletId,omitempty"`
	ExternalTxnID   *string           `json:"externalTxnId,omitempty"`
	ProcessedAt     *string           `json:"processedAt,omitempty"`
	CompletedAt     *string           `json:"completedAt,omitempty"`
	CreatedAt       string            `json:"createdAt"`
	UpdatedAt       string            `json:"updatedAt"`
}

// PaymentMethod represents a payment method
type PaymentMethod struct {
	ID           string            `json:"id"`
	UserID       string            `json:"userId"`
	Type         string            `json:"type"`
	Provider     string            `json:"provider"`
	DisplayName  string            `json:"displayName"`
	LastFour     *string           `json:"lastFour,omitempty"`
	ExpiryMonth  *int              `json:"expiryMonth,omitempty"`
	ExpiryYear   *int              `json:"expiryYear,omitempty"`
	IsDefault    bool              `json:"isDefault"`
	IsActive     bool              `json:"isActive"`
	Metadata     map[string]string `json:"metadata"`
	CreatedAt    string            `json:"createdAt"`
	UpdatedAt    string            `json:"updatedAt"`
}

// CreateWalletRequest represents wallet creation request
type CreateWalletRequest struct {
	UserID   string            `json:"userId"`
	Currency string            `json:"currency"`
	Type     string            `json:"type"`
	Settings map[string]string `json:"settings"`
}

// CreateTransactionRequest represents transaction creation request
type CreateTransactionRequest struct {
	WalletID      string            `json:"walletId"`
	Type          string            `json:"type"`
	Amount        float64           `json:"amount"`
	Currency      string            `json:"currency"`
	Description   string            `json:"description"`
	Reference     string            `json:"reference"`
	Metadata      map[string]string `json:"metadata"`
	ToWalletID    *string           `json:"toWalletId,omitempty"`
	ExternalTxnID *string           `json:"externalTxnId,omitempty"`
}

// AddPaymentMethodRequest represents payment method addition request
type AddPaymentMethodRequest struct {
	Type        string            `json:"type"`
	Provider    string            `json:"provider"`
	DisplayName string            `json:"displayName"`
	Token       string            `json:"token"`
	IsDefault   bool              `json:"isDefault"`
	Metadata    map[string]string `json:"metadata"`
}

// TransferRequest represents money transfer request
type TransferRequest struct {
	FromWalletID string            `json:"fromWalletId"`
	ToWalletID   string            `json:"toWalletId"`
	Amount       float64           `json:"amount"`
	Currency     string            `json:"currency"`
	Description  string            `json:"description"`
	Reference    string            `json:"reference"`
	Metadata     map[string]string `json:"metadata"`
}

// PaymentResponse represents payment API response
type PaymentResponse struct {
	Success       bool            `json:"success"`
	Status        string          `json:"status"`
	Message       string          `json:"message"`
	Wallet        *Wallet         `json:"wallet,omitempty"`
	Wallets       []Wallet        `json:"wallets,omitempty"`
	Transaction   *Transaction    `json:"transaction,omitempty"`
	Transactions  []Transaction   `json:"transactions,omitempty"`
	PaymentMethod *PaymentMethod  `json:"paymentMethod,omitempty"`
	PaymentMethods []PaymentMethod `json:"paymentMethods,omitempty"`
	Total         int             `json:"total,omitempty"`
	Page          int             `json:"page,omitempty"`
	Limit         int             `json:"limit,omitempty"`
	Timestamp     string          `json:"timestamp"`
}

// SetupSuite initializes the payment integration test suite
func (suite *PaymentIntegrationSuite) SetupSuite() {
	suite.APIIntegrationSuite.SetupSuite()
	suite.ports = DefaultServicePorts()

	// Wait for payment service to be available
	err := suite.waitForService(suite.ports.Payment, 30*time.Second)
	if err != nil {
		suite.T().Fatalf("Payment service not available: %v", err)
	}
}

// TestPaymentServiceHealth verifies payment service health endpoint
func (suite *PaymentIntegrationSuite) TestPaymentServiceHealth() {
	healthCheck, err := suite.checkServiceHealth(suite.ports.Payment)
	require.NoError(suite.T(), err, "Health check should succeed")

	assert.Equal(suite.T(), "healthy", healthCheck.Status)
	assert.Equal(suite.T(), "payment-service", healthCheck.Service)
	assert.NotEmpty(suite.T(), healthCheck.Timestamp)
}

// TestCreateWallet tests wallet creation
func (suite *PaymentIntegrationSuite) TestCreateWallet() {
	url := fmt.Sprintf("%s:%d/api/v1/wallets", suite.baseURL, suite.ports.Payment)

	createReq := CreateWalletRequest{
		UserID:   "test-user-wallet",
		Currency: "USD",
		Type:     "personal",
		Settings: map[string]string{
			"daily_limit":     "1000.00",
			"monthly_limit":   "10000.00",
			"auto_top_up":     "false",
		},
	}

	resp, err := suite.makeRequest("POST", url, createReq, nil)
	require.NoError(suite.T(), err, "Create wallet request should succeed")
	defer resp.Body.Close()

	// Should return 201 for successful creation
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var paymentResp PaymentResponse
	err = suite.parseResponse(resp, &paymentResp)
	require.NoError(suite.T(), err, "Should parse wallet creation response")

	assert.True(suite.T(), paymentResp.Success)
	assert.Equal(suite.T(), "success", paymentResp.Status)
	assert.NotNil(suite.T(), paymentResp.Wallet)
	assert.NotEmpty(suite.T(), paymentResp.Wallet.ID)
	assert.Equal(suite.T(), createReq.UserID, paymentResp.Wallet.UserID)
	assert.Equal(suite.T(), createReq.Currency, paymentResp.Wallet.Currency)
	assert.Equal(suite.T(), createReq.Type, paymentResp.Wallet.Type)
	assert.Equal(suite.T(), 0.0, paymentResp.Wallet.Balance) // Initial balance
	assert.Equal(suite.T(), "active", paymentResp.Wallet.Status) // Default status
}

// TestGetWallet tests retrieving a specific wallet
func (suite *PaymentIntegrationSuite) TestGetWallet() {
	// First create a wallet
	createURL := fmt.Sprintf("%s:%d/api/v1/wallets", suite.baseURL, suite.ports.Payment)
	createReq := CreateWalletRequest{
		UserID:   "test-user-get-wallet",
		Currency: "EUR",
		Type:     "business",
		Settings: map[string]string{"daily_limit": "5000.00"},
	}

	createResp, err := suite.makeRequest("POST", createURL, createReq, nil)
	require.NoError(suite.T(), err, "Create wallet for get test should succeed")
	defer createResp.Body.Close()

	var createPaymentResp PaymentResponse
	err = suite.parseResponse(createResp, &createPaymentResp)
	require.NoError(suite.T(), err, "Should parse create response")
	require.NotNil(suite.T(), createPaymentResp.Wallet)

	walletID := createPaymentResp.Wallet.ID

	// Now get the wallet
	getURL := fmt.Sprintf("%s:%d/api/v1/wallets/%s", suite.baseURL, suite.ports.Payment, walletID)
	resp, err := suite.makeRequest("GET", getURL, nil, nil)
	require.NoError(suite.T(), err, "Get wallet request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var paymentResp PaymentResponse
	err = suite.parseResponse(resp, &paymentResp)
	require.NoError(suite.T(), err, "Should parse get wallet response")

	assert.True(suite.T(), paymentResp.Success)
	assert.NotNil(suite.T(), paymentResp.Wallet)
	assert.Equal(suite.T(), walletID, paymentResp.Wallet.ID)
	assert.Equal(suite.T(), createReq.UserID, paymentResp.Wallet.UserID)
}

// TestListWallets tests listing wallets
func (suite *PaymentIntegrationSuite) TestListWallets() {
	url := fmt.Sprintf("%s:%d/api/v1/wallets", suite.baseURL, suite.ports.Payment)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "List wallets request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var paymentResp PaymentResponse
	err = suite.parseResponse(resp, &paymentResp)
	require.NoError(suite.T(), err, "Should parse list wallets response")

	assert.True(suite.T(), paymentResp.Success)
	assert.NotNil(suite.T(), paymentResp.Wallets)
	assert.GreaterOrEqual(suite.T(), paymentResp.Total, 0)
}

// TestCreateTransaction tests transaction creation
func (suite *PaymentIntegrationSuite) TestCreateTransaction() {
	// First create a wallet
	createWalletURL := fmt.Sprintf("%s:%d/api/v1/wallets", suite.baseURL, suite.ports.Payment)
	createWalletReq := CreateWalletRequest{
		UserID:   "test-user-transaction",
		Currency: "USD",
		Type:     "personal",
	}

	createWalletResp, err := suite.makeRequest("POST", createWalletURL, createWalletReq, nil)
	require.NoError(suite.T(), err, "Create wallet for transaction test should succeed")
	defer createWalletResp.Body.Close()

	var createWalletPaymentResp PaymentResponse
	err = suite.parseResponse(createWalletResp, &createWalletPaymentResp)
	require.NoError(suite.T(), err, "Should parse create wallet response")
	require.NotNil(suite.T(), createWalletPaymentResp.Wallet)

	walletID := createWalletPaymentResp.Wallet.ID

	// Now create a transaction
	createTxnURL := fmt.Sprintf("%s:%d/api/v1/transactions", suite.baseURL, suite.ports.Payment)
	createTxnReq := CreateTransactionRequest{
		WalletID:    walletID,
		Type:        "deposit",
		Amount:      100.50,
		Currency:    "USD",
		Description: "Test deposit transaction",
		Reference:   "TEST-DEP-001",
		Metadata: map[string]string{
			"source":     "bank_transfer",
			"bank_ref":   "BT123456789",
			"test_mode":  "true",
		},
	}

	resp, err := suite.makeRequest("POST", createTxnURL, createTxnReq, nil)
	require.NoError(suite.T(), err, "Create transaction request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var paymentResp PaymentResponse
	err = suite.parseResponse(resp, &paymentResp)
	require.NoError(suite.T(), err, "Should parse transaction creation response")

	assert.True(suite.T(), paymentResp.Success)
	assert.NotNil(suite.T(), paymentResp.Transaction)
	assert.NotEmpty(suite.T(), paymentResp.Transaction.ID)
	assert.Equal(suite.T(), walletID, paymentResp.Transaction.WalletID)
	assert.Equal(suite.T(), createTxnReq.Type, paymentResp.Transaction.Type)
	assert.Equal(suite.T(), createTxnReq.Amount, paymentResp.Transaction.Amount)
	assert.Equal(suite.T(), "pending", paymentResp.Transaction.Status) // Default status
}

// TestGetTransaction tests retrieving a specific transaction
func (suite *PaymentIntegrationSuite) TestGetTransaction() {
	// First create a wallet and transaction
	createWalletURL := fmt.Sprintf("%s:%d/api/v1/wallets", suite.baseURL, suite.ports.Payment)
	createWalletReq := CreateWalletRequest{
		UserID:   "test-user-get-txn",
		Currency: "USD",
		Type:     "personal",
	}

	createWalletResp, err := suite.makeRequest("POST", createWalletURL, createWalletReq, nil)
	require.NoError(suite.T(), err, "Create wallet should succeed")
	defer createWalletResp.Body.Close()

	var createWalletPaymentResp PaymentResponse
	err = suite.parseResponse(createWalletResp, &createWalletPaymentResp)
	require.NoError(suite.T(), err, "Should parse create wallet response")
	require.NotNil(suite.T(), createWalletPaymentResp.Wallet)

	walletID := createWalletPaymentResp.Wallet.ID

	// Create transaction
	createTxnURL := fmt.Sprintf("%s:%d/api/v1/transactions", suite.baseURL, suite.ports.Payment)
	createTxnReq := CreateTransactionRequest{
		WalletID:    walletID,
		Type:        "withdrawal",
		Amount:      50.25,
		Currency:    "USD",
		Description: "Test withdrawal transaction",
		Reference:   "TEST-WTH-001",
	}

	createTxnResp, err := suite.makeRequest("POST", createTxnURL, createTxnReq, nil)
	require.NoError(suite.T(), err, "Create transaction should succeed")
	defer createTxnResp.Body.Close()

	var createTxnPaymentResp PaymentResponse
	err = suite.parseResponse(createTxnResp, &createTxnPaymentResp)
	require.NoError(suite.T(), err, "Should parse create transaction response")
	require.NotNil(suite.T(), createTxnPaymentResp.Transaction)

	transactionID := createTxnPaymentResp.Transaction.ID

	// Now get the transaction
	getURL := fmt.Sprintf("%s:%d/api/v1/transactions/%s", suite.baseURL, suite.ports.Payment, transactionID)
	resp, err := suite.makeRequest("GET", getURL, nil, nil)
	require.NoError(suite.T(), err, "Get transaction request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var paymentResp PaymentResponse
	err = suite.parseResponse(resp, &paymentResp)
	require.NoError(suite.T(), err, "Should parse get transaction response")

	assert.True(suite.T(), paymentResp.Success)
	assert.NotNil(suite.T(), paymentResp.Transaction)
	assert.Equal(suite.T(), transactionID, paymentResp.Transaction.ID)
	assert.Equal(suite.T(), createTxnReq.Amount, paymentResp.Transaction.Amount)
}

// TestListTransactions tests listing transactions
func (suite *PaymentIntegrationSuite) TestListTransactions() {
	url := fmt.Sprintf("%s:%d/api/v1/transactions", suite.baseURL, suite.ports.Payment)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "List transactions request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var paymentResp PaymentResponse
	err = suite.parseResponse(resp, &paymentResp)
	require.NoError(suite.T(), err, "Should parse list transactions response")

	assert.True(suite.T(), paymentResp.Success)
	assert.NotNil(suite.T(), paymentResp.Transactions)
	assert.GreaterOrEqual(suite.T(), paymentResp.Total, 0)
}

// TestAddPaymentMethod tests adding a payment method
func (suite *PaymentIntegrationSuite) TestAddPaymentMethod() {
	userID := "test-user-payment-method"
	url := fmt.Sprintf("%s:%d/api/v1/users/%s/payment-methods", suite.baseURL, suite.ports.Payment, userID)

	addReq := AddPaymentMethodRequest{
		Type:        "card",
		Provider:    "stripe",
		DisplayName: "Visa **** 4242",
		Token:       "tok_test_4242424242424242",
		IsDefault:   true,
		Metadata: map[string]string{
			"last_four":    "4242",
			"expiry_month": "12",
			"expiry_year":  "2025",
			"brand":        "visa",
		},
	}

	resp, err := suite.makeRequest("POST", url, addReq, nil)
	require.NoError(suite.T(), err, "Add payment method request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var paymentResp PaymentResponse
	err = suite.parseResponse(resp, &paymentResp)
	require.NoError(suite.T(), err, "Should parse add payment method response")

	assert.True(suite.T(), paymentResp.Success)
	assert.NotNil(suite.T(), paymentResp.PaymentMethod)
	assert.NotEmpty(suite.T(), paymentResp.PaymentMethod.ID)
	assert.Equal(suite.T(), userID, paymentResp.PaymentMethod.UserID)
	assert.Equal(suite.T(), addReq.Type, paymentResp.PaymentMethod.Type)
	assert.Equal(suite.T(), addReq.Provider, paymentResp.PaymentMethod.Provider)
	assert.Equal(suite.T(), addReq.IsDefault, paymentResp.PaymentMethod.IsDefault)
	assert.True(suite.T(), paymentResp.PaymentMethod.IsActive) // Default active
}

// TestListPaymentMethods tests listing payment methods for a user
func (suite *PaymentIntegrationSuite) TestListPaymentMethods() {
	userID := "test-user-list-payment-methods"
	url := fmt.Sprintf("%s:%d/api/v1/users/%s/payment-methods", suite.baseURL, suite.ports.Payment, userID)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "List payment methods request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var paymentResp PaymentResponse
	err = suite.parseResponse(resp, &paymentResp)
	require.NoError(suite.T(), err, "Should parse list payment methods response")

	assert.True(suite.T(), paymentResp.Success)
	assert.NotNil(suite.T(), paymentResp.PaymentMethods)
	assert.GreaterOrEqual(suite.T(), paymentResp.Total, 0)
}

// TestWalletTransfer tests money transfer between wallets
func (suite *PaymentIntegrationSuite) TestWalletTransfer() {
	// First create two wallets
	createWalletURL := fmt.Sprintf("%s:%d/api/v1/wallets", suite.baseURL, suite.ports.Payment)

	// Source wallet
	createSourceWalletReq := CreateWalletRequest{
		UserID:   "test-user-source",
		Currency: "USD",
		Type:     "personal",
	}

	sourceResp, err := suite.makeRequest("POST", createWalletURL, createSourceWalletReq, nil)
	require.NoError(suite.T(), err, "Create source wallet should succeed")
	defer sourceResp.Body.Close()

	var sourcePaymentResp PaymentResponse
	err = suite.parseResponse(sourceResp, &sourcePaymentResp)
	require.NoError(suite.T(), err, "Should parse source wallet response")
	require.NotNil(suite.T(), sourcePaymentResp.Wallet)

	sourceWalletID := sourcePaymentResp.Wallet.ID

	// Destination wallet
	createDestWalletReq := CreateWalletRequest{
		UserID:   "test-user-dest",
		Currency: "USD",
		Type:     "personal",
	}

	destResp, err := suite.makeRequest("POST", createWalletURL, createDestWalletReq, nil)
	require.NoError(suite.T(), err, "Create destination wallet should succeed")
	defer destResp.Body.Close()

	var destPaymentResp PaymentResponse
	err = suite.parseResponse(destResp, &destPaymentResp)
	require.NoError(suite.T(), err, "Should parse destination wallet response")
	require.NotNil(suite.T(), destPaymentResp.Wallet)

	destWalletID := destPaymentResp.Wallet.ID

	// Perform transfer
	transferURL := fmt.Sprintf("%s:%d/api/v1/transfers", suite.baseURL, suite.ports.Payment)
	transferReq := TransferRequest{
		FromWalletID: sourceWalletID,
		ToWalletID:   destWalletID,
		Amount:       75.00,
		Currency:     "USD",
		Description:  "Test wallet transfer",
		Reference:    "TEST-TRANSFER-001",
		Metadata: map[string]string{
			"type":        "peer_to_peer",
			"category":    "friends_family",
			"test_mode":   "true",
		},
	}

	resp, err := suite.makeRequest("POST", transferURL, transferReq, nil)
	require.NoError(suite.T(), err, "Transfer request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var paymentResp PaymentResponse
	err = suite.parseResponse(resp, &paymentResp)
	require.NoError(suite.T(), err, "Should parse transfer response")

	assert.True(suite.T(), paymentResp.Success)
	assert.NotNil(suite.T(), paymentResp.Transaction)
	assert.Equal(suite.T(), "transfer", paymentResp.Transaction.Type)
	assert.Equal(suite.T(), transferReq.Amount, paymentResp.Transaction.Amount)
	assert.Equal(suite.T(), &sourceWalletID, paymentResp.Transaction.FromWalletID)
	assert.Equal(suite.T(), &destWalletID, paymentResp.Transaction.ToWalletID)
}

// TestGetNonExistentWallet tests retrieving non-existent wallet
func (suite *PaymentIntegrationSuite) TestGetNonExistentWallet() {
	url := fmt.Sprintf("%s:%d/api/v1/wallets/non-existent-id", suite.baseURL, suite.ports.Payment)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "Get non-existent wallet should complete")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

// TestCreateWalletInvalidData tests wallet creation with invalid data
func (suite *PaymentIntegrationSuite) TestCreateWalletInvalidData() {
	url := fmt.Sprintf("%s:%d/api/v1/wallets", suite.baseURL, suite.ports.Payment)

	// Test with missing required fields
	invalidReq := CreateWalletRequest{
		// Missing userID and currency
		Type: "personal",
	}

	resp, err := suite.makeRequest("POST", url, invalidReq, nil)
	require.NoError(suite.T(), err, "Invalid wallet creation should complete")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

// TestTransactionsByWallet tests filtering transactions by wallet
func (suite *PaymentIntegrationSuite) TestTransactionsByWallet() {
	// First create a wallet and transaction
	createWalletURL := fmt.Sprintf("%s:%d/api/v1/wallets", suite.baseURL, suite.ports.Payment)
	createWalletReq := CreateWalletRequest{
		UserID:   "test-user-filter-txn",
		Currency: "USD",
		Type:     "personal",
	}

	createWalletResp, err := suite.makeRequest("POST", createWalletURL, createWalletReq, nil)
	require.NoError(suite.T(), err, "Create wallet should succeed")
	defer createWalletResp.Body.Close()

	var createWalletPaymentResp PaymentResponse
	err = suite.parseResponse(createWalletResp, &createWalletPaymentResp)
	require.NoError(suite.T(), err, "Should parse create wallet response")
	require.NotNil(suite.T(), createWalletPaymentResp.Wallet)

	walletID := createWalletPaymentResp.Wallet.ID

	// Create a transaction for this wallet
	createTxnURL := fmt.Sprintf("%s:%d/api/v1/transactions", suite.baseURL, suite.ports.Payment)
	createTxnReq := CreateTransactionRequest{
		WalletID:    walletID,
		Type:        "deposit",
		Amount:      200.00,
		Currency:    "USD",
		Description: "Test filter transaction",
		Reference:   "TEST-FILTER-001",
	}

	createTxnResp, err := suite.makeRequest("POST", createTxnURL, createTxnReq, nil)
	require.NoError(suite.T(), err, "Create transaction should succeed")
	createTxnResp.Body.Close()

	// Now filter transactions by wallet
	url := fmt.Sprintf("%s:%d/api/v1/transactions?walletId=%s", suite.baseURL, suite.ports.Payment, walletID)
	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "Filter transactions request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var paymentResp PaymentResponse
	err = suite.parseResponse(resp, &paymentResp)
	require.NoError(suite.T(), err, "Should parse filter transactions response")

	assert.True(suite.T(), paymentResp.Success)

	// All returned transactions should belong to the specified wallet
	for _, transaction := range paymentResp.Transactions {
		assert.Equal(suite.T(), walletID, transaction.WalletID)
	}
}

// TestInvalidHTTPMethods tests endpoints with invalid HTTP methods
func (suite *PaymentIntegrationSuite) TestInvalidHTTPMethods() {
	baseURL := fmt.Sprintf("%s:%d/api/v1/wallets", suite.baseURL, suite.ports.Payment)

	testCases := []struct {
		url    string
		method string
	}{
		{baseURL, "PATCH"},           // List endpoint with invalid method
		{baseURL + "/123", "POST"},   // Get endpoint with invalid method
		{baseURL + "/123", "PATCH"},  // Update endpoint with invalid method
	}

	for _, tc := range testCases {
		resp, err := suite.makeRequest(tc.method, tc.url, nil, nil)
		require.NoError(suite.T(), err, "Invalid method request should complete")
		defer resp.Body.Close()

		assert.Equal(suite.T(), http.StatusMethodNotAllowed, resp.StatusCode,
			"URL: %s, Method: %s", tc.url, tc.method)
	}
}

// TestPaymentServiceConcurrency tests concurrent requests to payment service
func (suite *PaymentIntegrationSuite) TestPaymentServiceConcurrency() {
	url := fmt.Sprintf("%s:%d/api/v1/wallets", suite.baseURL, suite.ports.Payment)

	// Create 5 concurrent wallet creation requests
	concurrency := 5
	results := make(chan int, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(index int) {
			createReq := CreateWalletRequest{
				UserID:   fmt.Sprintf("concurrent-user-%d", index),
				Currency: "USD",
				Type:     "concurrency-test",
			}

			resp, err := suite.makeRequest("POST", url, createReq, nil)
			if err != nil {
				results <- 0
				return
			}
			defer resp.Body.Close()

			results <- resp.StatusCode
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < concurrency; i++ {
		statusCode := <-results
		if statusCode == http.StatusCreated {
			successCount++
		}
	}

	// At least 80% of concurrent requests should succeed
	assert.GreaterOrEqual(suite.T(), successCount, 4, "Concurrent requests should mostly succeed")
}

// RunPaymentIntegrationTests runs the payment integration test suite
func RunPaymentIntegrationTests(t *testing.T) {
	suite.Run(t, new(PaymentIntegrationSuite))
}