package fixtures

import (
	"tchat-backend/auth/models"
	"time"

	"github.com/google/uuid"
	"tchat-backend/payment/models"
)

// PaymentFixtures provides test data for Payment models
type PaymentFixtures struct {
	*BaseFixture
}

// NewPaymentFixtures creates a new payment fixtures instance
func NewPaymentFixtures(seed ...int64) *PaymentFixtures {
	return &PaymentFixtures{
		BaseFixture: NewBaseFixture(seed...),
	}
}

// WalletFixtures provides test data for Wallet models
type WalletFixtures struct {
	*BaseFixture
}

// NewWalletFixtures creates a new wallet fixtures instance
func NewWalletFixtures(seed ...int64) *WalletFixtures {
	return &WalletFixtures{
		BaseFixture: NewBaseFixture(seed...),
	}
}

// BasicWallet creates a basic wallet for testing
func (w *WalletFixtures) BasicWallet(userID uuid.UUID, country string) *models.Wallet {
	currency := w.Currency(country)
	balance := w.Amount(currency)

	return &models.Wallet{
		ID:             w.UUID("wallet-" + userID.String() + "-" + currency),
		UserID:         userID,
		Balance:        balance,
		Currency:       models.Currency(currency),
		FrozenBalance:  0,
		DailyLimit:     balance * 10,   // 10x balance as daily limit
		MonthlyLimit:   balance * 100,  // 100x balance as monthly limit
		UsedThisDay:    0,
		UsedThisMonth:  0,
		LastResetDay:   w.PastTime(1440), // Reset yesterday
		LastResetMonth: time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC),
		Status:         models.WalletStatusActive,
		IsPrimary:      true,
		CreatedAt:      w.PastTime(2880), // Created 2 days ago
		UpdatedAt:      w.PastTime(60),   // Updated 1 hour ago
	}
}

// FrozenWallet creates a wallet with frozen funds for testing
func (w *WalletFixtures) FrozenWallet(userID uuid.UUID, country string) *models.Wallet {
	wallet := w.BasicWallet(userID, country)
	wallet.ID = w.UUID("frozen-wallet-" + userID.String())
	wallet.FrozenBalance = wallet.Balance / 2 // Freeze half the balance
	wallet.Status = models.WalletStatusFrozen
	return wallet
}

// LimitedWallet creates a wallet with usage near limits for testing
func (w *WalletFixtures) LimitedWallet(userID uuid.UUID, country string) *models.Wallet {
	wallet := w.BasicWallet(userID, country)
	wallet.ID = w.UUID("limited-wallet-" + userID.String())
	wallet.UsedThisDay = wallet.DailyLimit * 8 / 10    // 80% of daily limit used
	wallet.UsedThisMonth = wallet.MonthlyLimit * 6 / 10 // 60% of monthly limit used
	return wallet
}

// MultiCurrencyWallets creates wallets for multiple currencies
func (w *WalletFixtures) MultiCurrencyWallets(userID uuid.UUID) []*models.Wallet {
	countries := []string{"TH", "SG", "ID", "MY", "VN", "PH"}
	wallets := make([]*models.Wallet, 0, len(countries))

	for i, country := range countries {
		wallet := w.BasicWallet(userID, country)
		wallet.ID = w.UUID("multi-wallet-" + country + "-" + userID.String())
		wallet.IsPrimary = (i == 0) // First wallet is primary
		wallets = append(wallets, wallet)
	}

	return wallets
}

// TransactionFixtures provides test data for Transaction models
type TransactionFixtures struct {
	*BaseFixture
}

// NewTransactionFixtures creates a new transaction fixtures instance
func NewTransactionFixtures(seed ...int64) *TransactionFixtures {
	return &TransactionFixtures{
		BaseFixture: NewBaseFixture(seed...),
	}
}

// BasicTransaction creates a basic transaction for testing
func (t *TransactionFixtures) BasicTransaction(fromUserID, toUserID uuid.UUID, currency string) *models.Transaction {
	amount := t.Amount(currency)

	return &models.Transaction{
		ID:              t.UUID("transaction-" + fromUserID.String() + "-" + toUserID.String()),
		WalletID:        fromUserID, // Assuming fromUserID is the wallet ID
		CounterpartyID:  &toUserID,
		Amount:          amount,
		Currency:        models.Currency(currency),
		Type:            models.TransactionTypeTransfer,
		Status:          models.TransactionStatusCompleted,
		Description:     &[]string{"Test transfer transaction"}[0],
		Reference:       t.Token(16),
		ExternalID:      nil,
		FeeAmount:       amount / 100, // 1% fee
		NetAmount:       amount - (amount / 100),
		Metadata: map[string]string{
			"source":      "test_fixture",
			"category":    "peer_to_peer",
			"method":      "wallet",
		},
		ProcessedAt:     &[]time.Time{t.PastTime(30)}[0], // Processed 30 minutes ago
		CreatedAt:       t.PastTime(60),  // Created 1 hour ago
		UpdatedAt:       t.PastTime(30),  // Updated 30 minutes ago
	}
}

// PendingTransaction creates a pending transaction for testing
func (t *TransactionFixtures) PendingTransaction(fromUserID, toUserID uuid.UUID, currency string) *models.Transaction {
	transaction := t.BasicTransaction(fromUserID, toUserID, currency)
	transaction.ID = t.UUID("pending-transaction-" + fromUserID.String())
	transaction.Status = models.TransactionStatusPending
	transaction.ProcessedAt = nil
	transaction.UpdatedAt = transaction.CreatedAt
	return transaction
}

// FailedTransaction creates a failed transaction for testing
func (t *TransactionFixtures) FailedTransaction(fromUserID, toUserID uuid.UUID, currency string) *models.Transaction {
	transaction := t.BasicTransaction(fromUserID, toUserID, currency)
	transaction.ID = t.UUID("failed-transaction-" + fromUserID.String())
	transaction.Status = models.TransactionStatusFailed
	description := "Test failed transaction - insufficient funds"
	transaction.Description = &description
	transaction.Metadata["failure_reason"] = "insufficient_funds"
	transaction.Metadata["error_code"] = "E001"
	return transaction
}

// TopUpTransaction creates a top-up transaction for testing
func (t *TransactionFixtures) TopUpTransaction(userID uuid.UUID, currency string) *models.Transaction {
	amount := t.Amount(currency)

	return &models.Transaction{
		ID:          t.UUID("topup-transaction-" + userID.String()),
		WalletID:    userID, // User's wallet receiving the top-up
		CounterpartyID: nil, // No counterparty for top-ups
		Amount:      amount,
		Currency:    models.Currency(currency),
		Type:        models.TransactionTypeDeposit,
		Status:      models.TransactionStatusCompleted,
		Description: &[]string{"Test wallet top-up"}[0],
		Reference:   t.Token(16),
		ExternalID:  &[]string{"ext-topup-" + t.Token(8)}[0],
		FeeAmount:   0, // No fee for top-ups
		NetAmount:   amount,
		Metadata: map[string]string{
			"source":        "test_fixture",
			"payment_method": "bank_transfer",
			"gateway":       "test_gateway",
		},
		ProcessedAt: &[]time.Time{t.PastTime(15)}[0], // Processed 15 minutes ago
		CreatedAt:   t.PastTime(30), // Created 30 minutes ago
		UpdatedAt:   t.PastTime(15), // Updated 15 minutes ago
	}
}

// WithdrawalTransaction creates a withdrawal transaction for testing
func (t *TransactionFixtures) WithdrawalTransaction(userID uuid.UUID, currency string) *models.Transaction {
	amount := t.Amount(currency)

	return &models.Transaction{
		ID:          t.UUID("withdrawal-transaction-" + userID.String()),
		WalletID:    userID, // User's wallet for withdrawal
		CounterpartyID: nil, // No counterparty for withdrawals
		Amount:      amount,
		Currency:    models.Currency(currency),
		Type:        models.TransactionTypeWithdrawal,
		Status:      models.TransactionStatusCompleted,
		Description: &[]string{"Test wallet withdrawal"}[0],
		Reference:   t.Token(16),
		ExternalID:  &[]string{"ext-withdrawal-" + t.Token(8)}[0],
		FeeAmount:   amount / 50, // 2% fee for withdrawals
		NetAmount:   amount - (amount / 50),
		Metadata: map[string]string{
			"source":        "test_fixture",
			"payment_method": "bank_transfer",
			"gateway":       "test_gateway",
			"bank_account":  "****1234",
		},
		ProcessedAt: &[]time.Time{t.PastTime(45)}[0], // Processed 45 minutes ago
		CreatedAt:   t.PastTime(60), // Created 1 hour ago
		UpdatedAt:   t.PastTime(45), // Updated 45 minutes ago
	}
}

// PaymentMethodFixtures provides test data for PaymentMethod models
type PaymentMethodFixtures struct {
	*BaseFixture
}

// NewPaymentMethodFixtures creates a new payment method fixtures instance
func NewPaymentMethodFixtures(seed ...int64) *PaymentMethodFixtures {
	return &PaymentMethodFixtures{
		BaseFixture: NewBaseFixture(seed...),
	}
}

// BankAccountPaymentMethod creates a bank account payment method for testing
func (p *PaymentMethodFixtures) BankAccountPaymentMethod(userID uuid.UUID, country string) *models.PaymentMethod {
	return &models.PaymentMethod{
		ID:       p.UUID("bank-payment-" + userID.String()),
		UserID:   userID,
		Type:     models.PaymentMethodTypeBankAccount,
		Provider: models.PaymentProviderDBS, // Use DBS as test bank provider
		Status:   models.PaymentMethodStatusActive,
		DisplayName:    "Test Bank Account",
		Country:        country,
		Currency:       models.Currency(p.Currency(country)),
		Metadata: models.PaymentMethodMetadata{
			HolderName:  p.Name(country),
			BankName:    "Test Bank " + country,
			BankCode:    "TESTBANK" + country,
			AccountType: "checking",
		},
		SecurityInfo: models.SecurityInfo{
			CVVVerified:       true,
			AddressVerified:   true,
			PhoneVerified:     true,
			EmailVerified:     true,
			RiskScore:         10, // Low risk
			FraudFlags:        []string{},
			LastSecurityCheck: p.PastTime(60), // Last check 1 hour ago
			SecurityLevel:     "high",
		},
		IsDefault:    true,
		IsVerified:   true,
		ExpiresAt:    nil, // Bank accounts don't expire
		CreatedAt:    p.PastTime(2880), // Created 2 days ago
		UpdatedAt:    p.PastTime(60),   // Updated 1 hour ago
	}
}

// CreditCardPaymentMethod creates a credit card payment method for testing
func (p *PaymentMethodFixtures) CreditCardPaymentMethod(userID uuid.UUID) *models.PaymentMethod {
	expiresAt := p.FutureTime(525600) // Expires in 1 year

	return &models.PaymentMethod{
		ID:       p.UUID("card-payment-" + userID.String()),
		UserID:   userID,
		Type:     models.PaymentMethodTypeCreditCard,
		Provider: models.PaymentProviderVisa,
		Status:   models.PaymentMethodStatusActive,
		DisplayName:    "Visa ****1234",
		LastFourDigits: &[]string{"1234"}[0],
		ExpiryMonth:    &[]int{12}[0],
		ExpiryYear:     &[]int{2025}[0],
		BrandName:      &[]string{"Visa"}[0],
		Country:        "TH",
		Currency:       models.CurrencyTHB,
		Metadata: models.PaymentMethodMetadata{
			HolderName: "Test Cardholder",
		},
		SecurityInfo: models.SecurityInfo{
			CVVVerified:       true,
			AddressVerified:   true,
			PhoneVerified:     true,
			EmailVerified:     true,
			RiskScore:         25, // Medium-low risk
			FraudFlags:        []string{},
			LastSecurityCheck: p.PastTime(60), // Last check 1 hour ago
			SecurityLevel:     "medium",
		},
		IsDefault:    false,
		IsVerified:   true,
		ExpiresAt:    &expiresAt,
		CreatedAt:    p.PastTime(1440), // Created yesterday
		UpdatedAt:    p.PastTime(60),   // Updated 1 hour ago
	}
}

// EWalletPaymentMethod creates an e-wallet payment method for testing
func (p *PaymentMethodFixtures) EWalletPaymentMethod(userID uuid.UUID, country string) *models.PaymentMethod {
	provider := "test_ewallet"
	switch country {
	case "TH":
		provider = "truemoney"
	case "SG":
		provider = "grabpay"
	case "ID":
		provider = "gopay"
	case "MY":
		provider = "tng"
	case "VN":
		provider = "momo"
	case "PH":
		provider = "gcash"
	}

	return &models.PaymentMethod{
		ID:       p.UUID("ewallet-payment-" + userID.String()),
		UserID:   userID,
		Type:     models.PaymentMethodTypeEWallet,
		Provider: models.PaymentProvider(provider),
		Status:   models.PaymentMethodStatusActive,
		DisplayName:    provider + " Wallet",
		Country:        country,
		Currency:       models.Currency(p.Currency(country)),
		Metadata: models.PaymentMethodMetadata{
			HolderName: p.Name(country),
		},
		SecurityInfo: models.SecurityInfo{
			CVVVerified:       true,
			AddressVerified:   true,
			PhoneVerified:     true,
			EmailVerified:     true,
			RiskScore:         15, // Low-medium risk
			FraudFlags:        []string{},
			LastSecurityCheck: p.PastTime(60), // Last check 1 hour ago
			SecurityLevel:     "medium",
		},
		IsDefault:    false,
		IsVerified:   true,
		ExpiresAt:    nil, // E-wallets don't expire
		CreatedAt:    p.PastTime(720), // Created 12 hours ago
		UpdatedAt:    p.PastTime(60),  // Updated 1 hour ago
	}
}

// TestPaymentData creates a comprehensive set of payment test data
func (p *PaymentFixtures) TestPaymentData(userID uuid.UUID, country string) map[string]interface{} {
	walletFixtures := NewWalletFixtures()
	transactionFixtures := NewTransactionFixtures()
	paymentMethodFixtures := NewPaymentMethodFixtures()

	// Create wallets
	wallets := []*models.Wallet{
		walletFixtures.BasicWallet(userID, country),
		walletFixtures.FrozenWallet(userID, country),
		walletFixtures.LimitedWallet(userID, country),
	}

	// Add multi-currency wallets
	multiCurrencyWallets := walletFixtures.MultiCurrencyWallets(userID)
	wallets = append(wallets, multiCurrencyWallets...)

	// Create transactions
	toUserID := p.UUID("recipient-user")
	currency := p.Currency(country)

	transactions := []*models.Transaction{
		transactionFixtures.BasicTransaction(userID, toUserID, currency),
		transactionFixtures.PendingTransaction(userID, toUserID, currency),
		transactionFixtures.FailedTransaction(userID, toUserID, currency),
		transactionFixtures.TopUpTransaction(userID, currency),
		transactionFixtures.WithdrawalTransaction(userID, currency),
	}

	// Create payment methods
	paymentMethods := []*models.PaymentMethod{
		paymentMethodFixtures.BankAccountPaymentMethod(userID, country),
		paymentMethodFixtures.CreditCardPaymentMethod(userID),
		paymentMethodFixtures.EWalletPaymentMethod(userID, country),
	}

	return map[string]interface{}{
		"wallets":         wallets,
		"transactions":    transactions,
		"payment_methods": paymentMethods,
	}
}

// AllPaymentFixtures creates a complete set of payment-related test data
func AllPaymentFixtures(seed ...int64) (*PaymentFixtures, *WalletFixtures, *TransactionFixtures, *PaymentMethodFixtures) {
	return NewPaymentFixtures(seed...), NewWalletFixtures(seed...), NewTransactionFixtures(seed...), NewPaymentMethodFixtures(seed...)
}