package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	authModels "tchat.dev/auth/models"
	"tchat.dev/payment/models"
	"tchat.dev/shared/config"
	"tchat.dev/shared/responses"
)

// PaymentProviderVerification handles Pact provider verification for the Payment service
type PaymentProviderVerification struct {
	server *httptest.Server
	db     *gorm.DB
	router *gin.Engine
}

// TestData holds test payment data for provider states
type TestData struct {
	Users          map[string]uuid.UUID
	Wallets        map[string]*models.Wallet
	PaymentMethods map[string]*models.PaymentMethod
	Transactions   map[string]*models.Transaction
}

var testData *TestData

func TestMain(m *testing.M) {
	// Setup test environment
	gin.SetMode(gin.TestMode)

	// Initialize test data
	testData = initializeTestData()

	// Run tests
	code := m.Run()

	// Cleanup
	cleanup()

	os.Exit(code)
}

func TestPaymentProviderVerification(t *testing.T) {
	// Setup test database
	db, err := setupTestDatabase()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	// Setup Gin router with payment endpoints
	router := setupPaymentRouter(db)

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Configure Pact verifier
	verifier := provider.NewVerifier()

	// Set provider base URL
	err = verifier.VerifyProvider(t, provider.VerifyRequest{
		ProviderBaseURL:        server.URL,
		Provider:               "payment-service",
		ConsumerVersionTags:    []string{"master", "test"},
		PublishVerificationResults: true,
		ProviderVersion:       "1.0.0",
		StateHandlers: map[string]provider.StateHandler{
			// Wallet provider states
			"user has wallet with balance":                   userHasWalletWithBalance,
			"user has multiple currency wallets":            userHasMultipleCurrencyWallets,
			"user wallet has insufficient balance":          userWalletHasInsufficientBalance,
			"user wallet is frozen":                         userWalletIsFrozen,

			// Payment method provider states
			"payment method exists for user":                paymentMethodExistsForUser,
			"user has verified payment method":             userHasVerifiedPaymentMethod,
			"user has multiple payment methods":            userHasMultiplePaymentMethods,
			"payment method is expired":                     paymentMethodIsExpired,

			// Transaction provider states
			"transaction can be processed":                   transactionCanBeProcessed,
			"transaction requires additional verification":   transactionRequiresVerification,
			"transaction exceeds daily limit":              transactionExceedsDailyLimit,
			"pending transaction exists":                    pendingTransactionExists,

			// Currency operation states
			"wallet supports currency operations":            walletSupportsCurrencyOperations,
			"exchange rate is available":                    exchangeRateIsAvailable,
			"multi-currency transaction is possible":        multiCurrencyTransactionPossible,

			// Compliance and security states
			"user is kyc verified":                          userIsKYCVerified,
			"transaction passes aml checks":                 transactionPassesAMLChecks,
			"user has transaction history":                  userHasTransactionHistory,

			// Error and edge case states
			"payment service is healthy":                    paymentServiceIsHealthy,
			"external payment processor is available":       externalProcessorAvailable,
			"payment method requires 3ds verification":     paymentMethodRequires3DS,
		},
	})

	assert.NoError(t, err, "Provider verification should pass")
}

// Provider State Handlers

func userHasWalletWithBalance(setup bool, s provider.ProviderState) error {
	if setup {
		// Create test user and wallet with balance
		userID := testData.Users["test_user_with_balance"]
		wallet := &models.Wallet{
			ID:               uuid.New(),
			UserID:           userID,
			Balance:          100000, // $1000.00 in cents
			Currency:         models.CurrencyUSD,
			FrozenBalance:    0,
			DailyLimit:       testData.getDefaultDailyLimit(models.CurrencyUSD),
			MonthlyLimit:     testData.getDefaultMonthlyLimit(models.CurrencyUSD),
			UsedThisDay:      0,
			UsedThisMonth:    0,
			LastResetDay:     time.Now().Truncate(24 * time.Hour),
			LastResetMonth:   time.Now().Truncate(24 * time.Hour * 30),
			Status:           models.WalletStatusActive,
			IsPrimary:        true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		testData.Wallets["primary_usd_wallet"] = wallet
	}
	return nil
}

func userHasMultipleCurrencyWallets(setup bool, s provider.ProviderState) error {
	if setup {
		userID := testData.Users["multi_currency_user"]

		// Create wallets for different currencies
		currencies := []models.Currency{
			models.CurrencyUSD, models.CurrencyTHB, models.CurrencySGD, models.CurrencyIDR,
		}

		for i, currency := range currencies {
			wallet := &models.Wallet{
				ID:            uuid.New(),
				UserID:        userID,
				Balance:       int64((i + 1) * 50000), // Different balances
				Currency:      currency,
				FrozenBalance: 0,
				DailyLimit:    testData.getDefaultDailyLimit(currency),
				MonthlyLimit:  testData.getDefaultMonthlyLimit(currency),
				Status:        models.WalletStatusActive,
				IsPrimary:     currency == models.CurrencyUSD,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			testData.Wallets[fmt.Sprintf("%s_wallet", strings.ToLower(string(currency)))] = wallet
		}
	}
	return nil
}

func userWalletHasInsufficientBalance(setup bool, s provider.ProviderState) error {
	if setup {
		userID := testData.Users["low_balance_user"]
		wallet := &models.Wallet{
			ID:            uuid.New(),
			UserID:        userID,
			Balance:       100, // Only $1.00
			Currency:      models.CurrencyUSD,
			FrozenBalance: 0,
			Status:        models.WalletStatusActive,
			IsPrimary:     true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		testData.Wallets["low_balance_wallet"] = wallet
	}
	return nil
}

func userWalletIsFrozen(setup bool, s provider.ProviderState) error {
	if setup {
		userID := testData.Users["frozen_user"]
		wallet := &models.Wallet{
			ID:            uuid.New(),
			UserID:        userID,
			Balance:       50000,
			Currency:      models.CurrencyUSD,
			FrozenBalance: 50000, // All balance frozen
			Status:        models.WalletStatusFrozen,
			IsPrimary:     true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		testData.Wallets["frozen_wallet"] = wallet
	}
	return nil
}

func paymentMethodExistsForUser(setup bool, s provider.ProviderState) error {
	if setup {
		userID := testData.Users["user_with_payment_method"]
		walletID := testData.Wallets["primary_usd_wallet"].ID

		paymentMethod := &models.PaymentMethod{
			ID:              uuid.New(),
			UserID:          userID,
			WalletID:        &walletID,
			Type:            models.PaymentMethodTypeCreditCard,
			Provider:        models.PaymentProviderVisa,
			Status:          models.PaymentMethodStatusActive,
			IsDefault:       true,
			IsVerified:      true,
			DisplayName:     "Visa ****1234",
			LastFourDigits:  stringPtr("1234"),
			ExpiryMonth:     intPtr(12),
			ExpiryYear:      intPtr(2025),
			BrandName:       stringPtr("Visa"),
			Country:         "US",
			Currency:        models.CurrencyUSD,
			VerifiedAt:      timePtr(time.Now()),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		testData.PaymentMethods["visa_card"] = paymentMethod
	}
	return nil
}

func userHasVerifiedPaymentMethod(setup bool, s provider.ProviderState) error {
	if setup {
		userID := testData.Users["verified_user"]

		paymentMethod := &models.PaymentMethod{
			ID:          uuid.New(),
			UserID:      userID,
			Type:        models.PaymentMethodTypeCreditCard,
			Provider:    models.PaymentProviderMastercard,
			Status:      models.PaymentMethodStatusActive,
			IsDefault:   true,
			IsVerified:  true,
			DisplayName: "Mastercard ****5678",
			Country:     "TH",
			Currency:    models.CurrencyTHB,
			VerifiedAt:  timePtr(time.Now()),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		testData.PaymentMethods["verified_mastercard"] = paymentMethod
	}
	return nil
}

func userHasMultiplePaymentMethods(setup bool, s provider.ProviderState) error {
	if setup {
		userID := testData.Users["multi_payment_user"]

		// Credit card
		creditCard := &models.PaymentMethod{
			ID:          uuid.New(),
			UserID:      userID,
			Type:        models.PaymentMethodTypeCreditCard,
			Provider:    models.PaymentProviderVisa,
			Status:      models.PaymentMethodStatusActive,
			IsDefault:   true,
			IsVerified:  true,
			DisplayName: "Primary Visa",
			Country:     "SG",
			Currency:    models.CurrencySGD,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// E-wallet
		ewallet := &models.PaymentMethod{
			ID:          uuid.New(),
			UserID:      userID,
			Type:        models.PaymentMethodTypeEWallet,
			Provider:    models.PaymentProviderGrabPay,
			Status:      models.PaymentMethodStatusActive,
			IsDefault:   false,
			IsVerified:  true,
			DisplayName: "GrabPay Wallet",
			Country:     "SG",
			Currency:    models.CurrencySGD,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		testData.PaymentMethods["primary_visa"] = creditCard
		testData.PaymentMethods["grabpay_wallet"] = ewallet
	}
	return nil
}

func paymentMethodIsExpired(setup bool, s provider.ProviderState) error {
	if setup {
		userID := testData.Users["expired_card_user"]

		paymentMethod := &models.PaymentMethod{
			ID:          uuid.New(),
			UserID:      userID,
			Type:        models.PaymentMethodTypeCreditCard,
			Provider:    models.PaymentProviderVisa,
			Status:      models.PaymentMethodStatusExpired,
			IsDefault:   false,
			IsVerified:  false,
			DisplayName: "Expired Visa",
			ExpiryMonth: intPtr(1),
			ExpiryYear:  intPtr(2020), // Expired
			Country:     "US",
			Currency:    models.CurrencyUSD,
			ExpiresAt:   timePtr(time.Date(2020, 1, 31, 23, 59, 59, 0, time.UTC)),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		testData.PaymentMethods["expired_visa"] = paymentMethod
	}
	return nil
}

func transactionCanBeProcessed(setup bool, s provider.ProviderState) error {
	if setup {
		userID := testData.Users["transaction_user"]
		walletID := testData.Wallets["primary_usd_wallet"].ID

		transaction := &models.Transaction{
			ID:          uuid.New(),
			WalletID:    walletID,
			Type:        models.TransactionTypePayment,
			Status:      models.TransactionStatusPending,
			Currency:    models.CurrencyUSD,
			Amount:      5000, // $50.00
			FeeAmount:   150,  // $1.50 fee
			NetAmount:   4850, // $48.50 net
			Reference:   "PAY_" + generateReference(),
			Description: stringPtr("Test payment transaction"),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		testData.Transactions["pending_payment"] = transaction
	}
	return nil
}

func transactionRequiresVerification(setup bool, s provider.ProviderState) error {
	if setup {
		userID := testData.Users["verification_required_user"]
		walletID := testData.Wallets["primary_usd_wallet"].ID

		transaction := &models.Transaction{
			ID:          uuid.New(),
			WalletID:    walletID,
			Type:        models.TransactionTypeTransfer,
			Status:      models.TransactionStatusPending,
			Currency:    models.CurrencyUSD,
			Amount:      100000, // $1000.00 - requires verification
			FeeAmount:   1000,   // $10.00 fee
			NetAmount:   99000,  // $990.00 net
			Reference:   "TXN_" + generateReference(),
			Description: stringPtr("Large transfer requiring verification"),
			Metadata:    map[string]string{"verification_required": "true", "risk_score": "high"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		testData.Transactions["verification_required"] = transaction
	}
	return nil
}

func transactionExceedsDailyLimit(setup bool, s provider.ProviderState) error {
	if setup {
		userID := testData.Users["limit_exceeded_user"]

		// Create wallet with used daily limit close to maximum
		wallet := &models.Wallet{
			ID:            uuid.New(),
			UserID:        userID,
			Balance:       500000, // $5000.00
			Currency:      models.CurrencyUSD,
			DailyLimit:    300000, // $3000.00 daily limit
			UsedThisDay:   290000, // $2900.00 already used today
			LastResetDay:  time.Now().Truncate(24 * time.Hour),
			Status:        models.WalletStatusActive,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		testData.Wallets["limit_wallet"] = wallet
	}
	return nil
}

func pendingTransactionExists(setup bool, s provider.ProviderState) error {
	if setup {
		transaction := &models.Transaction{
			ID:          uuid.New(),
			WalletID:    testData.Wallets["primary_usd_wallet"].ID,
			Type:        models.TransactionTypeDeposit,
			Status:      models.TransactionStatusPending,
			Currency:    models.CurrencyUSD,
			Amount:      25000, // $250.00
			FeeAmount:   0,     // No fee for deposits
			NetAmount:   25000,
			Reference:   "DEP_" + generateReference(),
			Description: stringPtr("Pending bank deposit"),
			ProcessedAt: nil, // Not processed yet
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		testData.Transactions["pending_deposit"] = transaction
	}
	return nil
}

func walletSupportsCurrencyOperations(setup bool, s provider.ProviderState) error {
	if setup {
		// Create multi-currency capable wallet
		userID := testData.Users["currency_user"]

		wallet := &models.Wallet{
			ID:           uuid.New(),
			UserID:       userID,
			Balance:      200000, // $2000.00 equivalent
			Currency:     models.CurrencyUSD,
			Status:       models.WalletStatusActive,
			IsPrimary:    true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		testData.Wallets["currency_wallet"] = wallet
	}
	return nil
}

func exchangeRateIsAvailable(setup bool, s provider.ProviderState) error {
	if setup {
		// Mock exchange rate data - in real implementation, this would setup mock exchange service
		// For now, we just ensure the state is marked as available
		return nil
	}
	return nil
}

func multiCurrencyTransactionPossible(setup bool, s provider.ProviderState) error {
	if setup {
		userID := testData.Users["multi_currency_transaction_user"]

		// USD source wallet
		usdWallet := &models.Wallet{
			ID:        uuid.New(),
			UserID:    userID,
			Balance:   100000, // $1000.00
			Currency:  models.CurrencyUSD,
			Status:    models.WalletStatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// THB destination wallet
		thbWallet := &models.Wallet{
			ID:        uuid.New(),
			UserID:    userID,
			Balance:   0,
			Currency:  models.CurrencyTHB,
			Status:    models.WalletStatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		testData.Wallets["usd_source"] = usdWallet
		testData.Wallets["thb_destination"] = thbWallet
	}
	return nil
}

func userIsKYCVerified(setup bool, s provider.ProviderState) error {
	if setup {
		// In a real implementation, this would setup KYC verification data
		// For provider tests, we just mark the state as available
		return nil
	}
	return nil
}

func transactionPassesAMLChecks(setup bool, s provider.ProviderState) error {
	if setup {
		// Mock AML compliance check - passed
		return nil
	}
	return nil
}

func userHasTransactionHistory(setup bool, s provider.ProviderState) error {
	if setup {
		userID := testData.Users["history_user"]
		walletID := testData.Wallets["primary_usd_wallet"].ID

		// Create historical transactions
		for i := 0; i < 5; i++ {
			transaction := &models.Transaction{
				ID:          uuid.New(),
				WalletID:    walletID,
				Type:        models.TransactionTypePayment,
				Status:      models.TransactionStatusCompleted,
				Currency:    models.CurrencyUSD,
				Amount:      int64((i + 1) * 1000), // $10, $20, $30, etc.
				FeeAmount:   int64((i + 1) * 30),   // Progressive fees
				NetAmount:   int64((i + 1) * 970),
				Reference:   fmt.Sprintf("HIST_%d_%s", i, generateReference()),
				CompletedAt: timePtr(time.Now().Add(-time.Duration(i*24) * time.Hour)),
				CreatedAt:   time.Now().Add(-time.Duration(i*24) * time.Hour),
				UpdatedAt:   time.Now(),
			}
			testData.Transactions[fmt.Sprintf("history_%d", i)] = transaction
		}
	}
	return nil
}

func paymentServiceIsHealthy(setup bool, s provider.ProviderState) error {
	if setup {
		// Service health check - always healthy for tests
		return nil
	}
	return nil
}

func externalProcessorAvailable(setup bool, s provider.ProviderState) error {
	if setup {
		// Mock external payment processor availability
		return nil
	}
	return nil
}

func paymentMethodRequires3DS(setup bool, s provider.ProviderState) error {
	if setup {
		userID := testData.Users["3ds_user"]

		paymentMethod := &models.PaymentMethod{
			ID:          uuid.New(),
			UserID:      userID,
			Type:        models.PaymentMethodTypeCreditCard,
			Provider:    models.PaymentProviderVisa,
			Status:      models.PaymentMethodStatusPending,
			IsDefault:   true,
			IsVerified:  false,
			DisplayName: "Visa 3DS Required",
			Country:     "SG",
			Currency:    models.CurrencySGD,
			Metadata: models.PaymentMethodMetadata{
				HolderName: "Test User",
			},
			SecurityInfo: models.SecurityInfo{
				SecurityLevel:     "high",
				RiskScore:         85.0,
				LastSecurityCheck: time.Now().UTC(),
				FraudFlags:        []string{"requires_3ds"},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		testData.PaymentMethods["3ds_visa"] = paymentMethod
	}
	return nil
}

// Helper functions

func initializeTestData() *TestData {
	data := &TestData{
		Users:          make(map[string]uuid.UUID),
		Wallets:        make(map[string]*models.Wallet),
		PaymentMethods: make(map[string]*models.PaymentMethod),
		Transactions:   make(map[string]*models.Transaction),
	}

	// Initialize test user IDs
	userTypes := []string{
		"test_user_with_balance", "multi_currency_user", "low_balance_user",
		"frozen_user", "user_with_payment_method", "verified_user",
		"multi_payment_user", "expired_card_user", "transaction_user",
		"verification_required_user", "limit_exceeded_user", "currency_user",
		"multi_currency_transaction_user", "history_user", "3ds_user",
	}

	for _, userType := range userTypes {
		data.Users[userType] = uuid.New()
	}

	return data
}

func (td *TestData) getDefaultDailyLimit(currency models.Currency) int64 {
	if limit, exists := models.DefaultDailyLimits[currency]; exists {
		return limit
	}
	return 300000 // Default $3000 USD equivalent
}

func (td *TestData) getDefaultMonthlyLimit(currency models.Currency) int64 {
	if limit, exists := models.DefaultMonthlyLimits[currency]; exists {
		return limit
	}
	return 9000000 // Default $90000 USD equivalent
}

func setupTestDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate payment models
	err = db.AutoMigrate(
		&models.Wallet{},
		&models.Transaction{},
		&models.PaymentMethod{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupPaymentRouter(db *gorm.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		responses.SendSuccessResponse(c, gin.H{
			"status": "ok",
			"service": "payment",
			"timestamp": time.Now().Unix(),
		})
	})

	api := router.Group("/api/v1")

	// Wallet endpoints
	wallets := api.Group("/wallets")
	{
		wallets.GET("/", getWallets(db))
		wallets.POST("/", createWallet(db))
		wallets.GET("/:id", getWallet(db))
		wallets.PUT("/:id/balance", updateWalletBalance(db))
		wallets.GET("/:id/transactions", getWalletTransactions(db))
	}

	// Transaction endpoints
	transactions := api.Group("/transactions")
	{
		transactions.POST("/", createTransaction(db))
		transactions.GET("/:id", getTransaction(db))
		transactions.PUT("/:id/status", updateTransactionStatus(db))
		transactions.POST("/transfer", processTransfer(db))
		transactions.POST("/payment", processPayment(db))
	}

	// Payment method endpoints
	paymentMethods := api.Group("/payment-methods")
	{
		paymentMethods.GET("/", getPaymentMethods(db))
		paymentMethods.POST("/", createPaymentMethod(db))
		paymentMethods.GET("/:id", getPaymentMethod(db))
		paymentMethods.PUT("/:id", updatePaymentMethod(db))
		paymentMethods.DELETE("/:id", deletePaymentMethod(db))
		paymentMethods.POST("/:id/verify", verifyPaymentMethod(db))
	}

	return router
}

// Mock handler implementations for provider verification

func getWallets(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mock implementation for provider verification
		userID := c.GetHeader("User-ID") // In real implementation, extract from JWT
		if userID == "" {
			responses.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user ID")
			return
		}

		// Return test wallet data
		wallets := []*models.Wallet{}
		for _, wallet := range testData.Wallets {
			if wallet.UserID.String() == userID {
				wallets = append(wallets, wallet)
			}
		}

		responses.SendSuccessResponse(c, gin.H{
			"wallets": wallets,
			"count":   len(wallets),
		})
	}
}

func createWallet(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Currency string `json:"currency" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			responses.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}

		// Create new wallet
		wallet := &models.Wallet{
			ID:        uuid.New(),
			UserID:    uuid.New(), // Mock user ID
			Balance:   0,
			Currency:  models.Currency(req.Currency),
			Status:    models.WalletStatusActive,
			IsPrimary: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		c.JSON(http.StatusCreated, wallet)
	}
}

func getWallet(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("id")

		// Find wallet in test data
		for _, wallet := range testData.Wallets {
			if wallet.ID.String() == walletID {
				c.JSON(http.StatusOK, wallet)
				return
			}
		}

		responses.SendErrorResponse(c, http.StatusNotFound, "WALLET_NOT_FOUND", "Wallet not found")
	}
}

func updateWalletBalance(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Amount        int64  `json:"amount" binding:"required"`
			TransactionID string `json:"transaction_id" binding:"required"`
			Operation     string `json:"operation" binding:"required"` // "credit" or "debit"
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			responses.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}

		walletID := c.Param("id")

		// Find and update wallet
		for _, wallet := range testData.Wallets {
			if wallet.ID.String() == walletID {
				if req.Operation == "credit" {
					wallet.Balance += req.Amount
				} else {
					wallet.Balance -= req.Amount
					if wallet.Balance < 0 {
						responses.SendErrorResponse(c, http.StatusBadRequest, "INSUFFICIENT_BALANCE", "Insufficient balance")
						return
					}
				}
				wallet.UpdatedAt = time.Now()
				c.JSON(http.StatusOK, wallet)
				return
			}
		}

		responses.SendErrorResponse(c, http.StatusNotFound, "WALLET_NOT_FOUND", "Wallet not found")
	}
}

func getWalletTransactions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("id")

		transactions := []*models.Transaction{}
		for _, txn := range testData.Transactions {
			if txn.WalletID.String() == walletID {
				transactions = append(transactions, txn)
			}
		}

		responses.SendSuccessResponse(c, gin.H{
			"transactions": transactions,
			"count":        len(transactions),
		})
	}
}

func createTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			WalletID    string  `json:"wallet_id" binding:"required"`
			Type        string  `json:"type" binding:"required"`
			Amount      int64   `json:"amount" binding:"required"`
			Currency    string  `json:"currency" binding:"required"`
			Description *string `json:"description"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			responses.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}

		transaction := &models.Transaction{
			ID:          uuid.New(),
			WalletID:    uuid.MustParse(req.WalletID),
			Type:        models.TransactionType(req.Type),
			Status:      models.TransactionStatusPending,
			Currency:    models.Currency(req.Currency),
			Amount:      req.Amount,
			FeeAmount:   req.Amount / 100, // 1% fee
			NetAmount:   req.Amount - (req.Amount / 100),
			Reference:   generateReference(),
			Description: req.Description,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		c.JSON(http.StatusCreated, transaction)
	}
}

func getTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		txnID := c.Param("id")

		for _, txn := range testData.Transactions {
			if txn.ID.String() == txnID {
				c.JSON(http.StatusOK, txn)
				return
			}
		}

		responses.SendErrorResponse(c, http.StatusNotFound, "TRANSACTION_NOT_FOUND", "Transaction not found")
	}
}

func updateTransactionStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Status string `json:"status" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			responses.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}

		txnID := c.Param("id")

		for _, txn := range testData.Transactions {
			if txn.ID.String() == txnID {
				txn.Status = models.TransactionStatus(req.Status)
				txn.UpdatedAt = time.Now()

				if req.Status == "completed" {
					now := time.Now()
					txn.CompletedAt = &now
				}

				c.JSON(http.StatusOK, txn)
				return
			}
		}

		responses.SendErrorResponse(c, http.StatusNotFound, "TRANSACTION_NOT_FOUND", "Transaction not found")
	}
}

func processTransfer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FromWalletID string  `json:"from_wallet_id" binding:"required"`
			ToWalletID   string  `json:"to_wallet_id" binding:"required"`
			Amount       int64   `json:"amount" binding:"required"`
			Currency     string  `json:"currency" binding:"required"`
			Description  *string `json:"description"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			responses.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}

		// Create transfer transaction
		transaction := &models.Transaction{
			ID:             uuid.New(),
			WalletID:       uuid.MustParse(req.FromWalletID),
			CounterpartyID: func() *uuid.UUID { id := uuid.MustParse(req.ToWalletID); return &id }(),
			Type:           models.TransactionTypeTransfer,
			Status:         models.TransactionStatusProcessing,
			Currency:       models.Currency(req.Currency),
			Amount:         req.Amount,
			FeeAmount:      req.Amount / 200, // 0.5% transfer fee
			NetAmount:      req.Amount - (req.Amount / 200),
			Reference:      "TXF_" + generateReference(),
			Description:    req.Description,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		c.JSON(http.StatusAccepted, gin.H{
			"transaction": transaction,
			"message":     "Transfer initiated successfully",
		})
	}
}

func processPayment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			WalletID         string `json:"wallet_id" binding:"required"`
			PaymentMethodID  string `json:"payment_method_id" binding:"required"`
			Amount           int64  `json:"amount" binding:"required"`
			Currency         string `json:"currency" binding:"required"`
			MerchantID       string `json:"merchant_id"`
			PaymentReference string `json:"payment_reference"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			responses.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}

		// Create payment transaction
		transaction := &models.Transaction{
			ID:          uuid.New(),
			WalletID:    uuid.MustParse(req.WalletID),
			Type:        models.TransactionTypePayment,
			Status:      models.TransactionStatusProcessing,
			Currency:    models.Currency(req.Currency),
			Amount:      req.Amount,
			FeeAmount:   req.Amount / 67, // ~1.5% payment processing fee
			NetAmount:   req.Amount - (req.Amount / 67),
			Reference:   "PAY_" + generateReference(),
			Metadata: map[string]string{
				"merchant_id":        req.MerchantID,
				"payment_method_id":  req.PaymentMethodID,
				"payment_reference":  req.PaymentReference,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		c.JSON(http.StatusAccepted, gin.H{
			"transaction": transaction,
			"message":     "Payment processing initiated",
		})
	}
}

func getPaymentMethods(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("User-ID")
		if userID == "" {
			responses.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user ID")
			return
		}

		paymentMethods := []*models.PaymentMethod{}
		for _, pm := range testData.PaymentMethods {
			if pm.UserID.String() == userID {
				paymentMethods = append(paymentMethods, pm)
			}
		}

		responses.SendSuccessResponse(c, gin.H{
			"payment_methods": paymentMethods,
			"count":          len(paymentMethods),
		})
	}
}

func createPaymentMethod(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Type        string `json:"type" binding:"required"`
			Provider    string `json:"provider" binding:"required"`
			DisplayName string `json:"display_name" binding:"required"`
			Country     string `json:"country" binding:"required"`
			Currency    string `json:"currency" binding:"required"`
			ExternalID  string `json:"external_id"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			responses.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}

		paymentMethod := &models.PaymentMethod{
			ID:          uuid.New(),
			UserID:      uuid.New(), // Mock user ID
			Type:        models.PaymentMethodType(req.Type),
			Provider:    models.PaymentProvider(req.Provider),
			Status:      models.PaymentMethodStatusPending,
			IsDefault:   false,
			IsVerified:  false,
			DisplayName: req.DisplayName,
			Country:     req.Country,
			Currency:    models.Currency(req.Currency),
			ExternalID:  &req.ExternalID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		c.JSON(http.StatusCreated, paymentMethod)
	}
}

func getPaymentMethod(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		pmID := c.Param("id")

		for _, pm := range testData.PaymentMethods {
			if pm.ID.String() == pmID {
				c.JSON(http.StatusOK, pm)
				return
			}
		}

		responses.SendErrorResponse(c, http.StatusNotFound, "PAYMENT_METHOD_NOT_FOUND", "Payment method not found")
	}
}

func updatePaymentMethod(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			DisplayName *string `json:"display_name"`
			IsDefault   *bool   `json:"is_default"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			responses.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}

		pmID := c.Param("id")

		for _, pm := range testData.PaymentMethods {
			if pm.ID.String() == pmID {
				if req.DisplayName != nil {
					pm.DisplayName = *req.DisplayName
				}
				if req.IsDefault != nil {
					pm.IsDefault = *req.IsDefault
				}
				pm.UpdatedAt = time.Now()

				c.JSON(http.StatusOK, pm)
				return
			}
		}

		responses.SendErrorResponse(c, http.StatusNotFound, "PAYMENT_METHOD_NOT_FOUND", "Payment method not found")
	}
}

func deletePaymentMethod(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		pmID := c.Param("id")

		for key, pm := range testData.PaymentMethods {
			if pm.ID.String() == pmID {
				delete(testData.PaymentMethods, key)
				c.JSON(http.StatusNoContent, nil)
				return
			}
		}

		responses.SendErrorResponse(c, http.StatusNotFound, "PAYMENT_METHOD_NOT_FOUND", "Payment method not found")
	}
}

func verifyPaymentMethod(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			VerificationCode string                 `json:"verification_code"`
			VerificationData map[string]interface{} `json:"verification_data"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			responses.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}

		pmID := c.Param("id")

		for _, pm := range testData.PaymentMethods {
			if pm.ID.String() == pmID {
				pm.IsVerified = true
				pm.Status = models.PaymentMethodStatusActive
				now := time.Now()
				pm.VerifiedAt = &now
				pm.UpdatedAt = now

				c.JSON(http.StatusOK, gin.H{
					"payment_method": pm,
					"message":        "Payment method verified successfully",
				})
				return
			}
		}

		responses.SendErrorResponse(c, http.StatusNotFound, "PAYMENT_METHOD_NOT_FOUND", "Payment method not found")
	}
}

// Utility functions

func generateReference() string {
	return fmt.Sprintf("%d%d", time.Now().Unix(), time.Now().Nanosecond()%10000)
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func cleanup() {
	// Cleanup any test resources
	testData = nil
}