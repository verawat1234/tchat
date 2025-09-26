package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TransactionType represents the type of financial transaction
type TransactionType string

const (
	TransactionTypeDeposit     TransactionType = "deposit"
	TransactionTypeWithdrawal  TransactionType = "withdrawal"
	TransactionTypeTransferIn  TransactionType = "transfer_in"
	TransactionTypeTransferOut TransactionType = "transfer_out"
	TransactionTypePayment     TransactionType = "payment"
	TransactionTypeRefund      TransactionType = "refund"
	TransactionTypeConversion  TransactionType = "conversion"
	TransactionTypeFee         TransactionType = "fee"
)

// IsValid checks if the transaction type is valid
func (tt TransactionType) IsValid() bool {
	switch tt {
	case TransactionTypeDeposit, TransactionTypeWithdrawal, TransactionTypeTransferIn,
		 TransactionTypeTransferOut, TransactionTypePayment, TransactionTypeRefund,
		 TransactionTypeConversion, TransactionTypeFee:
		return true
	default:
		return false
	}
}

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending    TransactionStatus = "pending"
	TransactionStatusProcessing TransactionStatus = "processing"
	TransactionStatusCompleted  TransactionStatus = "completed"
	TransactionStatusFailed     TransactionStatus = "failed"
	TransactionStatusCancelled  TransactionStatus = "cancelled"
	TransactionStatusExpired    TransactionStatus = "expired"
)

// IsValid checks if the transaction status is valid
func (ts TransactionStatus) IsValid() bool {
	switch ts {
	case TransactionStatusPending, TransactionStatusProcessing, TransactionStatusCompleted,
		 TransactionStatusFailed, TransactionStatusCancelled, TransactionStatusExpired:
		return true
	default:
		return false
	}
}

// IsTerminal checks if the transaction status is terminal (final)
func (ts TransactionStatus) IsTerminal() bool {
	switch ts {
	case TransactionStatusCompleted, TransactionStatusFailed, TransactionStatusCancelled, TransactionStatusExpired:
		return true
	default:
		return false
	}
}

// PaymentGateway represents the payment gateway used
type PaymentGateway string

const (
	PaymentGatewayStripe       PaymentGateway = "stripe"
	PaymentGatewayOmise        PaymentGateway = "omise"        // Thailand
	PaymentGatewayMidtrans     PaymentGateway = "midtrans"     // Indonesia
	PaymentGatewayRazorpay     PaymentGateway = "razorpay"     // Regional
	PaymentGateway2C2P         PaymentGateway = "2c2p"         // Southeast Asia
	PaymentGatewayPayPal       PaymentGateway = "paypal"
	PaymentGatewayBankTransfer PaymentGateway = "bank_transfer"
	PaymentGatewayEWallet      PaymentGateway = "e_wallet"
	PaymentGatewayInternal     PaymentGateway = "internal"     // Internal wallet transfers
)

// TransactionFees represents the fee structure for a transaction
type TransactionFees struct {
	ProcessingFee decimal.Decimal `json:"processing_fee" gorm:"column:processing_fee;type:decimal(20,8)"`
	ConversionFee decimal.Decimal `json:"conversion_fee" gorm:"column:conversion_fee;type:decimal(20,8)"`
	NetworkFee    decimal.Decimal `json:"network_fee" gorm:"column:network_fee;type:decimal(20,8)"`
	PlatformFee   decimal.Decimal `json:"platform_fee" gorm:"column:platform_fee;type:decimal(20,8)"`
	TotalFees     decimal.Decimal `json:"total_fees" gorm:"column:total_fees;type:decimal(20,8)"`
}

// AuditTrailEntry represents an entry in the transaction audit trail
type AuditTrailEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	PreviousStatus string `json:"previous_status"`
	NewStatus   string    `json:"new_status"`
	Reason      string    `json:"reason,omitempty"`
	OperatorID  uuid.UUID `json:"operator_id,omitempty"`
	GatewayData map[string]interface{} `json:"gateway_data,omitempty"`
}

// Transaction represents a financial transaction in the system
type Transaction struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	WalletID uuid.UUID `json:"wallet_id" gorm:"type:uuid;not null;index"`

	// Transaction details
	Type        TransactionType   `json:"type" gorm:"column:type;type:varchar(20);not null"`
	Status      TransactionStatus `json:"status" gorm:"column:status;type:varchar(20);not null;default:'pending'"`
	Amount      decimal.Decimal   `json:"amount" gorm:"column:amount;type:decimal(20,8);not null"`
	Currency    string            `json:"currency" gorm:"column:currency;size:3;not null"`
	Description string            `json:"description" gorm:"column:description;size:500"`

	// Payment gateway integration
	Gateway                PaymentGateway `json:"gateway" gorm:"column:gateway;type:varchar(20);not null"`
	GatewayTransactionID   string         `json:"gateway_transaction_id,omitempty" gorm:"column:gateway_transaction_id;size:255"`
	GatewayReference       string         `json:"gateway_reference,omitempty" gorm:"column:gateway_reference;size:255"`
	GatewayStatus          string         `json:"gateway_status,omitempty" gorm:"column:gateway_status;size:50"`
	GatewayResponseCode    string         `json:"gateway_response_code,omitempty" gorm:"column:gateway_response_code;size:20"`
	GatewayResponseMessage string         `json:"gateway_response_message,omitempty" gorm:"column:gateway_response_message;size:500"`

	// Currency conversion (if applicable)
	OriginalAmount   *decimal.Decimal `json:"original_amount,omitempty" gorm:"column:original_amount;type:decimal(20,8)"`
	OriginalCurrency string           `json:"original_currency,omitempty" gorm:"column:original_currency;size:3"`
	ExchangeRate     *decimal.Decimal `json:"exchange_rate,omitempty" gorm:"column:exchange_rate;type:decimal(20,8)"`
	RateProvider     string           `json:"rate_provider,omitempty" gorm:"column:rate_provider;size:50"`
	RateTimestamp    *time.Time       `json:"rate_timestamp,omitempty" gorm:"column:rate_timestamp"`

	// Fees structure
	Fees TransactionFees `json:"fees" gorm:"embedded;embeddedPrefix:fee_"`

	// Related transactions
	ParentTransactionID *uuid.UUID `json:"parent_transaction_id,omitempty" gorm:"column:parent_transaction_id;type:uuid;index"`
	RelatedTransactionID *uuid.UUID `json:"related_transaction_id,omitempty" gorm:"column:related_transaction_id;type:uuid"`
	BatchID             *uuid.UUID `json:"batch_id,omitempty" gorm:"column:batch_id;type:uuid;index"`

	// Regional compliance
	DataRegion         string `json:"data_region" gorm:"column:data_region;size:20"`
	ComplianceData     map[string]interface{} `json:"compliance_data,omitempty" gorm:"column:compliance_data;type:jsonb"`
	RequiresApproval   bool   `json:"requires_approval" gorm:"column:requires_approval;default:false"`
	ApprovedBy         *uuid.UUID `json:"approved_by,omitempty" gorm:"column:approved_by;type:uuid"`
	ApprovedAt         *time.Time `json:"approved_at,omitempty" gorm:"column:approved_at"`

	// Risk management
	RiskScore          float64 `json:"risk_score" gorm:"column:risk_score;default:0"`
	RiskFlags          []string `json:"risk_flags,omitempty" gorm:"column:risk_flags;type:jsonb"`
	FraudCheck         bool     `json:"fraud_check" gorm:"column:fraud_check;default:false"`
	FraudCheckResult   string   `json:"fraud_check_result,omitempty" gorm:"column:fraud_check_result;size:20"`
	FraudCheckProvider string   `json:"fraud_check_provider,omitempty" gorm:"column:fraud_check_provider;size:50"`

	// Metadata and additional data
	Metadata map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
	Tags     []string                `json:"tags,omitempty" gorm:"column:tags;type:jsonb"`

	// Audit trail
	AuditTrail []AuditTrailEntry `json:"audit_trail" gorm:"column:audit_trail;type:jsonb"`

	// Timestamps
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"column:updated_at;not null"`
	ProcessedAt *time.Time `json:"processed_at,omitempty" gorm:"column:processed_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" gorm:"column:completed_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" gorm:"column:expires_at"`

	// Soft delete
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Relationships
	User             *User         `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Wallet           *Wallet       `json:"wallet,omitempty" gorm:"foreignKey:WalletID;references:ID"`
	ParentTransaction *Transaction `json:"parent_transaction,omitempty" gorm:"foreignKey:ParentTransactionID;references:ID"`
	ChildTransactions []Transaction `json:"child_transactions,omitempty" gorm:"foreignKey:ParentTransactionID;references:ID"`
}

// TableName returns the table name for the Transaction model
func (Transaction) TableName() string {
	return "transactions"
}

// BeforeCreate sets up the transaction before creation
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}

	// Initialize audit trail
	if len(t.AuditTrail) == 0 {
		t.AuditTrail = []AuditTrailEntry{
			{
				Timestamp:   time.Now(),
				PreviousStatus: "",
				NewStatus:   string(t.Status),
				Reason:      "Transaction created",
			},
		}
	}

	// Set expiry for pending transactions (24 hours)
	if t.Status == TransactionStatusPending && t.ExpiresAt == nil {
		expiry := time.Now().Add(24 * time.Hour)
		t.ExpiresAt = &expiry
	}

	// Validate the transaction
	if err := t.Validate(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate validates the transaction before updating
func (t *Transaction) BeforeUpdate(tx *gorm.DB) error {
	return t.Validate()
}

// Validate validates the transaction data
func (t *Transaction) Validate() error {
	// Validate UUIDs
	if t.ID == uuid.Nil {
		return fmt.Errorf("transaction ID cannot be nil")
	}
	if t.UserID == uuid.Nil {
		return fmt.Errorf("user ID cannot be nil")
	}
	if t.WalletID == uuid.Nil {
		return fmt.Errorf("wallet ID cannot be nil")
	}

	// Validate type and status
	if !t.Type.IsValid() {
		return fmt.Errorf("invalid transaction type: %s", t.Type)
	}
	if !t.Status.IsValid() {
		return fmt.Errorf("invalid transaction status: %s", t.Status)
	}

	// Validate amount
	if t.Amount.IsZero() || t.Amount.IsNegative() {
		return fmt.Errorf("transaction amount must be positive")
	}

	// Validate currency
	if !IsValidCurrency(t.Currency) {
		return fmt.Errorf("invalid currency: %s", t.Currency)
	}

	// Validate currency-specific precision
	if err := t.ValidateCurrencyPrecision(); err != nil {
		return err
	}

	// Validate gateway
	if t.Gateway == "" {
		return fmt.Errorf("payment gateway is required")
	}

	// Validate exchange rate data consistency
	if err := t.ValidateExchangeRateData(); err != nil {
		return err
	}

	return nil
}

// ValidateCurrencyPrecision validates amount precision for different currencies
func (t *Transaction) ValidateCurrencyPrecision() error {
	// Currencies with no decimal places
	noDecimalCurrencies := map[string]bool{
		"IDR": true, // Indonesian Rupiah
		"VND": true, // Vietnamese Dong
		"JPY": true, // Japanese Yen
		"KRW": true, // Korean Won
	}

	if noDecimalCurrencies[t.Currency] {
		if !t.Amount.Equal(t.Amount.Truncate(0)) {
			return fmt.Errorf("currency %s does not support decimal places", t.Currency)
		}
	} else {
		// Most currencies support 2 decimal places
		if t.Amount.Exponent() < -2 {
			return fmt.Errorf("currency %s supports maximum 2 decimal places", t.Currency)
		}
	}

	return nil
}

// ValidateExchangeRateData validates exchange rate data consistency
func (t *Transaction) ValidateExchangeRateData() error {
	// If exchange rate is provided, original amount and currency must also be provided
	if t.ExchangeRate != nil {
		if t.OriginalAmount == nil {
			return fmt.Errorf("original amount is required when exchange rate is provided")
		}
		if t.OriginalCurrency == "" {
			return fmt.Errorf("original currency is required when exchange rate is provided")
		}
		if t.OriginalCurrency == t.Currency {
			return fmt.Errorf("original currency cannot be the same as target currency")
		}

		// Validate calculated amount
		calculatedAmount := t.OriginalAmount.Mul(*t.ExchangeRate)
		if !calculatedAmount.Sub(t.Amount).Abs().LessThan(decimal.NewFromFloat(0.01)) {
			return fmt.Errorf("amount does not match original amount * exchange rate")
		}
	}

	return nil
}

// UpdateStatus updates the transaction status with audit trail
func (t *Transaction) UpdateStatus(newStatus TransactionStatus, reason string, operatorID *uuid.UUID) error {
	if !newStatus.IsValid() {
		return fmt.Errorf("invalid transaction status: %s", newStatus)
	}

	// Check if status transition is valid
	if !t.IsValidStatusTransition(newStatus) {
		return fmt.Errorf("invalid status transition from %s to %s", t.Status, newStatus)
	}

	// Add audit trail entry
	auditEntry := AuditTrailEntry{
		Timestamp:      time.Now(),
		PreviousStatus: string(t.Status),
		NewStatus:      string(newStatus),
		Reason:         reason,
	}
	if operatorID != nil {
		auditEntry.OperatorID = *operatorID
	}

	t.AuditTrail = append(t.AuditTrail, auditEntry)
	t.Status = newStatus
	t.UpdatedAt = time.Now()

	// Set timestamps for specific status changes
	switch newStatus {
	case TransactionStatusProcessing:
		if t.ProcessedAt == nil {
			now := time.Now()
			t.ProcessedAt = &now
		}
	case TransactionStatusCompleted:
		if t.CompletedAt == nil {
			now := time.Now()
			t.CompletedAt = &now
		}
	}

	return nil
}

// IsValidStatusTransition checks if a status transition is valid
func (t *Transaction) IsValidStatusTransition(newStatus TransactionStatus) bool {
	// Terminal statuses cannot be changed
	if t.Status.IsTerminal() {
		return false
	}

	validTransitions := map[TransactionStatus][]TransactionStatus{
		TransactionStatusPending: {
			TransactionStatusProcessing,
			TransactionStatusCancelled,
			TransactionStatusExpired,
			TransactionStatusFailed,
		},
		TransactionStatusProcessing: {
			TransactionStatusCompleted,
			TransactionStatusFailed,
		},
	}

	allowedTransitions, exists := validTransitions[t.Status]
	if !exists {
		return false
	}

	for _, allowed := range allowedTransitions {
		if newStatus == allowed {
			return true
		}
	}

	return false
}

// CalculateFees calculates transaction fees based on amount and gateway
func (t *Transaction) CalculateFees() {
	// Fee calculation based on gateway and region
	gatewayFeeRates := map[PaymentGateway]map[string]decimal.Decimal{
		PaymentGatewayStripe: {
			"processing_rate": decimal.NewFromFloat(0.029), // 2.9%
			"fixed_fee":       decimal.NewFromFloat(0.30),  // $0.30
		},
		PaymentGatewayOmise: {
			"processing_rate": decimal.NewFromFloat(0.0275), // 2.75%
			"fixed_fee":       decimal.NewFromFloat(0.25),   // Varies by currency
		},
		PaymentGateway2C2P: {
			"processing_rate": decimal.NewFromFloat(0.030), // 3.0%
			"fixed_fee":       decimal.NewFromFloat(0.20),
		},
		PaymentGatewayInternal: {
			"processing_rate": decimal.NewFromFloat(0.001), // 0.1% for internal transfers
			"fixed_fee":       decimal.Zero,
		},
	}

	if rates, exists := gatewayFeeRates[t.Gateway]; exists {
		processingRate := rates["processing_rate"]
		fixedFee := rates["fixed_fee"]

		t.Fees.ProcessingFee = t.Amount.Mul(processingRate).Add(fixedFee)
		t.Fees.TotalFees = t.Fees.ProcessingFee.Add(t.Fees.ConversionFee).Add(t.Fees.NetworkFee).Add(t.Fees.PlatformFee)
	}
}

// SetConversionData sets currency conversion data
func (t *Transaction) SetConversionData(originalAmount decimal.Decimal, originalCurrency string, exchangeRate decimal.Decimal, rateProvider string) {
	t.OriginalAmount = &originalAmount
	t.OriginalCurrency = originalCurrency
	t.ExchangeRate = &exchangeRate
	t.RateProvider = rateProvider
	now := time.Now()
	t.RateTimestamp = &now

	// Calculate conversion fee (0.5% for cross-currency transactions)
	conversionRate := decimal.NewFromFloat(0.005)
	t.Fees.ConversionFee = originalAmount.Mul(conversionRate)
}

// IsExpired checks if the transaction has expired
func (t *Transaction) IsExpired() bool {
	return t.ExpiresAt != nil && t.ExpiresAt.Before(time.Now())
}

// IsCompleted checks if the transaction is completed
func (t *Transaction) IsCompleted() bool {
	return t.Status == TransactionStatusCompleted
}

// IsPending checks if the transaction is pending
func (t *Transaction) IsPending() bool {
	return t.Status == TransactionStatusPending
}

// CanCancel checks if the transaction can be cancelled
func (t *Transaction) CanCancel() bool {
	return t.Status == TransactionStatusPending && !t.IsExpired()
}

// GetNetAmount returns the amount after fees
func (t *Transaction) GetNetAmount() decimal.Decimal {
	return t.Amount.Sub(t.Fees.TotalFees)
}

// GetRegionalSettings returns regional compliance settings
func (t *Transaction) GetRegionalSettings() map[string]interface{} {
	return map[string]interface{}{
		"data_region":      t.DataRegion,
		"compliance_data":  t.ComplianceData,
		"requires_approval": t.RequiresApproval,
		"risk_score":       t.RiskScore,
		"fraud_check":      t.FraudCheck,
	}
}

// MarshalJSON customizes JSON serialization
func (t *Transaction) MarshalJSON() ([]byte, error) {
	type Alias Transaction
	return json.Marshal(&struct {
		*Alias
		IsCompleted      bool            `json:"is_completed"`
		IsPending        bool            `json:"is_pending"`
		IsExpired        bool            `json:"is_expired"`
		CanCancel        bool            `json:"can_cancel"`
		NetAmount        decimal.Decimal `json:"net_amount"`
		RegionalSettings map[string]interface{} `json:"regional_settings"`
	}{
		Alias:           (*Alias)(t),
		IsCompleted:     t.IsCompleted(),
		IsPending:       t.IsPending(),
		IsExpired:       t.IsExpired(),
		CanCancel:       t.CanCancel(),
		NetAmount:       t.GetNetAmount(),
		RegionalSettings: t.GetRegionalSettings(),
	})
}

// Helper functions for currency validation

// IsValidCurrency checks if a currency code is valid
func IsValidCurrency(currency string) bool {
	validCurrencies := map[string]bool{
		"THB": true, // Thai Baht
		"SGD": true, // Singapore Dollar
		"IDR": true, // Indonesian Rupiah
		"MYR": true, // Malaysian Ringgit
		"PHP": true, // Philippine Peso
		"VND": true, // Vietnamese Dong
		"USD": true, // US Dollar
		"EUR": true, // Euro
	}
	return validCurrencies[currency]
}

// GetCurrencyPrecision returns the decimal precision for a currency
func GetCurrencyPrecision(currency string) int {
	noDecimalCurrencies := map[string]bool{
		"IDR": true,
		"VND": true,
		"JPY": true,
		"KRW": true,
	}

	if noDecimalCurrencies[currency] {
		return 0
	}
	return 2 // Most currencies use 2 decimal places
}

// FormatCurrencyAmount formats amount according to currency rules
func FormatCurrencyAmount(amount decimal.Decimal, currency string) string {
	precision := GetCurrencyPrecision(currency)
	return amount.StringFixed(int32(precision))
}

// GenerateSearchKeywords generates search keywords for the transaction
func (t *Transaction) GenerateSearchKeywords() []string {
	keywords := []string{
		string(t.Type),
		string(t.Status),
		t.Currency,
		string(t.Gateway),
	}

	// Add description keywords
	if t.Description != "" {
		descWords := strings.Fields(strings.ToLower(t.Description))
		keywords = append(keywords, descWords...)
	}

	// Add gateway reference if available
	if t.GatewayReference != "" {
		keywords = append(keywords, t.GatewayReference)
	}

	// Add tags
	keywords = append(keywords, t.Tags...)

	// Remove duplicates and empty strings
	seen := make(map[string]bool)
	var unique []string
	for _, keyword := range keywords {
		if keyword != "" && !seen[keyword] {
			seen[keyword] = true
			unique = append(unique, keyword)
		}
	}

	return unique
}

// GetGatewaysForCountry returns available payment gateways for a country
func GetGatewaysForCountry(countryCode string) []PaymentGateway {
	gatewaysByCountry := map[string][]PaymentGateway{
		"TH": {PaymentGatewayOmise, PaymentGatewayStripe, PaymentGateway2C2P, PaymentGatewayEWallet, PaymentGatewayBankTransfer},
		"SG": {PaymentGatewayStripe, PaymentGateway2C2P, PaymentGatewayPayPal, PaymentGatewayBankTransfer},
		"ID": {PaymentGatewayMidtrans, PaymentGateway2C2P, PaymentGatewayEWallet, PaymentGatewayBankTransfer},
		"MY": {PaymentGatewayStripe, PaymentGateway2C2P, PaymentGatewayRazorpay, PaymentGatewayEWallet},
		"PH": {PaymentGateway2C2P, PaymentGatewayStripe, PaymentGatewayEWallet, PaymentGatewayBankTransfer},
		"VN": {PaymentGateway2C2P, PaymentGatewayEWallet, PaymentGatewayBankTransfer},
	}

	if gateways, exists := gatewaysByCountry[countryCode]; exists {
		// Always include internal transfers and add global gateways
		result := []PaymentGateway{PaymentGatewayInternal}
		result = append(result, gateways...)
		return result
	}

	// Default gateways for unsupported countries
	return []PaymentGateway{PaymentGatewayInternal, PaymentGatewayStripe, PaymentGatewayPayPal}
}

// GetTransactionLimits returns transaction limits based on user verification level and country
func GetTransactionLimits(userVerified bool, countryCode string, currency string) map[string]decimal.Decimal {
	limits := make(map[string]decimal.Decimal)

	// Base limits (in USD equivalent)
	if userVerified {
		limits["daily_limit"] = decimal.NewFromInt(10000)    // $10,000
		limits["monthly_limit"] = decimal.NewFromInt(100000) // $100,000
		limits["single_transaction"] = decimal.NewFromInt(5000) // $5,000
	} else {
		limits["daily_limit"] = decimal.NewFromInt(1000)     // $1,000
		limits["monthly_limit"] = decimal.NewFromInt(5000)   // $5,000
		limits["single_transaction"] = decimal.NewFromInt(500) // $500
	}

	// Country-specific adjustments (simplified conversion rates)
	countryMultipliers := map[string]decimal.Decimal{
		"TH": decimal.NewFromInt(35),    // THB
		"SG": decimal.NewFromFloat(1.35), // SGD
		"ID": decimal.NewFromInt(15000), // IDR
		"MY": decimal.NewFromFloat(4.5),  // MYR
		"PH": decimal.NewFromInt(55),    // PHP
		"VN": decimal.NewFromInt(24000), // VND
	}

	if multiplier, exists := countryMultipliers[countryCode]; exists && currency != "USD" {
		for key, value := range limits {
			limits[key] = value.Mul(multiplier)
		}
	}

	return limits
}

// GetMinimumTransactionAmount returns minimum transaction amount for a currency
func GetMinimumTransactionAmount(currency string) decimal.Decimal {
	minimums := map[string]decimal.Decimal{
		"THB": decimal.NewFromInt(1),      // 1 THB
		"SGD": decimal.NewFromFloat(0.01), // 1 cent
		"IDR": decimal.NewFromInt(1000),   // 1000 IDR
		"MYR": decimal.NewFromFloat(0.01), // 1 sen
		"PHP": decimal.NewFromInt(1),      // 1 PHP
		"VND": decimal.NewFromInt(1000),   // 1000 VND
		"USD": decimal.NewFromFloat(0.01), // 1 cent
		"EUR": decimal.NewFromFloat(0.01), // 1 cent
	}

	if min, exists := minimums[currency]; exists {
		return min
	}
	return decimal.NewFromFloat(0.01) // Default minimum
}