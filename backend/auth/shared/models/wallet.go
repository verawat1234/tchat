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

// WalletType represents the type of wallet
type WalletType string

const (
	WalletTypePersonal  WalletType = "personal"  // Individual user wallet
	WalletTypeBusiness  WalletType = "business"  // Business wallet
	WalletTypeEscrow    WalletType = "escrow"    // Escrow wallet for transactions
	WalletTypeSavings   WalletType = "savings"   // Savings wallet
	WalletTypeRewards   WalletType = "rewards"   // Rewards/loyalty wallet
)

// IsValid checks if the wallet type is valid
func (wt WalletType) IsValid() bool {
	switch wt {
	case WalletTypePersonal, WalletTypeBusiness, WalletTypeEscrow, WalletTypeSavings, WalletTypeRewards:
		return true
	default:
		return false
	}
}

// WalletStatus represents the status of a wallet
type WalletStatus string

const (
	WalletStatusActive    WalletStatus = "active"
	WalletStatusFrozen    WalletStatus = "frozen"
	WalletStatusSuspended WalletStatus = "suspended"
	WalletStatusClosed    WalletStatus = "closed"
)

// IsValid checks if the wallet status is valid
func (ws WalletStatus) IsValid() bool {
	switch ws {
	case WalletStatusActive, WalletStatusFrozen, WalletStatusSuspended, WalletStatusClosed:
		return true
	default:
		return false
	}
}

// IsOperational checks if the wallet can process transactions
func (ws WalletStatus) IsOperational() bool {
	return ws == WalletStatusActive
}

// WalletBalance represents a currency balance in the wallet
type WalletBalance struct {
	Currency        string          `json:"currency" gorm:"column:currency;size:3;not null"`
	AvailableAmount decimal.Decimal `json:"available_amount" gorm:"column:available_amount;type:decimal(20,8);not null;default:0"`
	PendingAmount   decimal.Decimal `json:"pending_amount" gorm:"column:pending_amount;type:decimal(20,8);not null;default:0"`
	ReservedAmount  decimal.Decimal `json:"reserved_amount" gorm:"column:reserved_amount;type:decimal(20,8);not null;default:0"`
	TotalAmount     decimal.Decimal `json:"total_amount" gorm:"column:total_amount;type:decimal(20,8);not null;default:0"`
	LastUpdated     time.Time       `json:"last_updated" gorm:"column:last_updated;not null"`
}

// WalletSettings represents wallet configuration and preferences
type WalletSettings struct {
	DefaultCurrency      string   `json:"default_currency" gorm:"column:default_currency;size:3;not null;default:'USD'"`
	AllowedCurrencies    []string `json:"allowed_currencies" gorm:"column:allowed_currencies;type:jsonb"`
	AutoConvert          bool     `json:"auto_convert" gorm:"column:auto_convert;default:false"`
	PreferredExchange    string   `json:"preferred_exchange,omitempty" gorm:"column:preferred_exchange;size:50"`
	LowBalanceThreshold  decimal.Decimal `json:"low_balance_threshold" gorm:"column:low_balance_threshold;type:decimal(20,8);default:0"`
	HighBalanceThreshold decimal.Decimal `json:"high_balance_threshold" gorm:"column:high_balance_threshold;type:decimal(20,8);default:0"`
	NotificationsEnabled bool     `json:"notifications_enabled" gorm:"column:notifications_enabled;default:true"`
	TransactionLimits    map[string]interface{} `json:"transaction_limits,omitempty" gorm:"column:transaction_limits;type:jsonb"`
}

// WalletSecurity represents security settings and verification
type WalletSecurity struct {
	RequiresPIN         bool       `json:"requires_pin" gorm:"column:requires_pin;default:false"`
	Require2FA          bool       `json:"require_2fa" gorm:"column:require_2fa;default:false"`
	WithdrawalPassword  bool       `json:"withdrawal_password" gorm:"column:withdrawal_password;default:false"`
	LastSecurityCheck   *time.Time `json:"last_security_check,omitempty" gorm:"column:last_security_check"`
	SecurityLevel       string     `json:"security_level" gorm:"column:security_level;size:20;default:'basic'"`
	WhitelistedAddresses []string  `json:"whitelisted_addresses,omitempty" gorm:"column:whitelisted_addresses;type:jsonb"`
	DailyWithdrawalLimit decimal.Decimal `json:"daily_withdrawal_limit" gorm:"column:daily_withdrawal_limit;type:decimal(20,8)"`
	MonthlyWithdrawalLimit decimal.Decimal `json:"monthly_withdrawal_limit" gorm:"column:monthly_withdrawal_limit;type:decimal(20,8)"`
}

// Wallet represents a multi-currency digital wallet
type Wallet struct {
	ID       uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID   *uuid.UUID `json:"user_id,omitempty" gorm:"type:uuid;index"`
	BusinessID *uuid.UUID `json:"business_id,omitempty" gorm:"type:uuid;index"`

	// Wallet details
	Name        string       `json:"name" gorm:"column:name;size:100;not null"`
	Type        WalletType   `json:"type" gorm:"column:type;type:varchar(20);not null"`
	Status      WalletStatus `json:"status" gorm:"column:status;type:varchar(20);not null;default:'active'"`
	Description string       `json:"description,omitempty" gorm:"column:description;size:500"`

	// Multi-currency balances (stored as JSONB for flexibility)
	Balances []WalletBalance `json:"balances" gorm:"column:balances;type:jsonb"`

	// Configuration and security
	Settings WalletSettings `json:"settings" gorm:"embedded;embeddedPrefix:settings_"`
	Security WalletSecurity `json:"security" gorm:"embedded;embeddedPrefix:security_"`

	// Regional compliance
	DataRegion     string `json:"data_region" gorm:"column:data_region;size:20"`
	ComplianceData map[string]interface{} `json:"compliance_data,omitempty" gorm:"column:compliance_data;type:jsonb"`
	KYCLevel       string `json:"kyc_level" gorm:"column:kyc_level;size:20;default:'basic'"`

	// Verification and limits
	IsVerified       bool      `json:"is_verified" gorm:"column:is_verified;default:false"`
	VerifiedAt       *time.Time `json:"verified_at,omitempty" gorm:"column:verified_at"`
	VerificationLevel string   `json:"verification_level" gorm:"column:verification_level;size:20;default:'unverified'"`

	// Activity tracking
	LastActivityAt   *time.Time `json:"last_activity_at,omitempty" gorm:"column:last_activity_at"`
	TransactionCount int64      `json:"transaction_count" gorm:"column:transaction_count;default:0"`
	TotalVolume      decimal.Decimal `json:"total_volume" gorm:"column:total_volume;type:decimal(20,8);default:0"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
	Tags     []string               `json:"tags,omitempty" gorm:"column:tags;type:jsonb"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Relationships
	User         *User          `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Business     *Business      `json:"business,omitempty" gorm:"foreignKey:BusinessID;references:ID"`
	Transactions []Transaction  `json:"transactions,omitempty" gorm:"foreignKey:WalletID;references:ID"`
}

// TableName returns the table name for the Wallet model
func (Wallet) TableName() string {
	return "wallets"
}

// BeforeCreate sets up the wallet before creation
func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}

	// Ensure either UserID or BusinessID is set, but not both
	if w.UserID == nil && w.BusinessID == nil {
		return fmt.Errorf("wallet must belong to either a user or business")
	}
	if w.UserID != nil && w.BusinessID != nil {
		return fmt.Errorf("wallet cannot belong to both user and business")
	}

	// Set default data region based on user/business
	if w.DataRegion == "" {
		if w.UserID != nil {
			// Get user's country for data region
			var user User
			if err := tx.First(&user, w.UserID).Error; err == nil {
				w.DataRegion = GetDataRegionForCountry(user.CountryCode)
			}
		} else if w.BusinessID != nil {
			// Get business's country for data region
			var business Business
			if err := tx.First(&business, w.BusinessID).Error; err == nil {
				w.DataRegion = GetDataRegionForCountry(business.Address.Country)
			}
		}

		if w.DataRegion == "" {
			w.DataRegion = "sea-central" // Default region
		}
	}

	// Initialize default settings
	if len(w.Settings.AllowedCurrencies) == 0 {
		w.Settings.AllowedCurrencies = []string{"USD", "THB", "SGD", "IDR", "MYR", "PHP", "VND"}
	}

	// Initialize empty balances if not set
	if len(w.Balances) == 0 {
		w.Balances = []WalletBalance{}
	}

	// Set default security limits
	if w.Security.DailyWithdrawalLimit.IsZero() {
		w.Security.DailyWithdrawalLimit = decimal.NewFromInt(1000) // $1000 USD equivalent
	}
	if w.Security.MonthlyWithdrawalLimit.IsZero() {
		w.Security.MonthlyWithdrawalLimit = decimal.NewFromInt(10000) // $10000 USD equivalent
	}

	// Validate the wallet
	if err := w.Validate(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate validates the wallet before updating
func (w *Wallet) BeforeUpdate(tx *gorm.DB) error {
	return w.Validate()
}

// Validate validates the wallet data
func (w *Wallet) Validate() error {
	// Validate UUIDs
	if w.ID == uuid.Nil {
		return fmt.Errorf("wallet ID cannot be nil")
	}

	// Validate ownership
	if w.UserID == nil && w.BusinessID == nil {
		return fmt.Errorf("wallet must belong to either a user or business")
	}
	if w.UserID != nil && w.BusinessID != nil {
		return fmt.Errorf("wallet cannot belong to both user and business")
	}

	// Validate type and status
	if !w.Type.IsValid() {
		return fmt.Errorf("invalid wallet type: %s", w.Type)
	}
	if !w.Status.IsValid() {
		return fmt.Errorf("invalid wallet status: %s", w.Status)
	}

	// Validate name
	if len(w.Name) == 0 || len(w.Name) > 100 {
		return fmt.Errorf("wallet name must be between 1 and 100 characters")
	}

	// Validate default currency
	if !IsValidCurrency(w.Settings.DefaultCurrency) {
		return fmt.Errorf("invalid default currency: %s", w.Settings.DefaultCurrency)
	}

	// Validate allowed currencies
	for _, currency := range w.Settings.AllowedCurrencies {
		if !IsValidCurrency(currency) {
			return fmt.Errorf("invalid allowed currency: %s", currency)
		}
	}

	// Validate balances
	for i, balance := range w.Balances {
		if err := w.validateBalance(balance, i); err != nil {
			return err
		}
	}

	return nil
}

// validateBalance validates a single wallet balance
func (w *Wallet) validateBalance(balance WalletBalance, index int) error {
	if !IsValidCurrency(balance.Currency) {
		return fmt.Errorf("invalid currency in balance %d: %s", index, balance.Currency)
	}

	if balance.AvailableAmount.IsNegative() {
		return fmt.Errorf("available amount cannot be negative for %s", balance.Currency)
	}

	if balance.PendingAmount.IsNegative() {
		return fmt.Errorf("pending amount cannot be negative for %s", balance.Currency)
	}

	if balance.ReservedAmount.IsNegative() {
		return fmt.Errorf("reserved amount cannot be negative for %s", balance.Currency)
	}

	// Verify total amount calculation
	calculatedTotal := balance.AvailableAmount.Add(balance.PendingAmount).Add(balance.ReservedAmount)
	if !calculatedTotal.Equal(balance.TotalAmount) {
		return fmt.Errorf("total amount mismatch for %s: expected %s, got %s",
			balance.Currency, calculatedTotal.String(), balance.TotalAmount.String())
	}

	return nil
}

// GetBalance returns the balance for a specific currency
func (w *Wallet) GetBalance(currency string) (*WalletBalance, error) {
	for i := range w.Balances {
		if w.Balances[i].Currency == currency {
			return &w.Balances[i], nil
		}
	}
	return nil, fmt.Errorf("currency %s not found in wallet", currency)
}

// GetAvailableBalance returns the available balance for a currency
func (w *Wallet) GetAvailableBalance(currency string) decimal.Decimal {
	balance, err := w.GetBalance(currency)
	if err != nil {
		return decimal.Zero
	}
	return balance.AvailableAmount
}

// GetTotalBalance returns the total balance for a currency
func (w *Wallet) GetTotalBalance(currency string) decimal.Decimal {
	balance, err := w.GetBalance(currency)
	if err != nil {
		return decimal.Zero
	}
	return balance.TotalAmount
}

// HasSufficientBalance checks if wallet has sufficient available balance
func (w *Wallet) HasSufficientBalance(currency string, amount decimal.Decimal) bool {
	available := w.GetAvailableBalance(currency)
	return available.GreaterThanOrEqual(amount)
}

// AddBalance adds or updates a currency balance
func (w *Wallet) AddBalance(currency string, amount decimal.Decimal, balanceType string) error {
	if !IsValidCurrency(currency) {
		return fmt.Errorf("invalid currency: %s", currency)
	}

	if amount.IsNegative() {
		return fmt.Errorf("amount cannot be negative")
	}

	// Find existing balance or create new one
	var balance *WalletBalance
	found := false
	for i := range w.Balances {
		if w.Balances[i].Currency == currency {
			balance = &w.Balances[i]
			found = true
			break
		}
	}

	if !found {
		// Create new balance entry
		newBalance := WalletBalance{
			Currency:        currency,
			AvailableAmount: decimal.Zero,
			PendingAmount:   decimal.Zero,
			ReservedAmount:  decimal.Zero,
			TotalAmount:     decimal.Zero,
			LastUpdated:     time.Now(),
		}
		w.Balances = append(w.Balances, newBalance)
		balance = &w.Balances[len(w.Balances)-1]
	}

	// Update the appropriate balance type
	switch balanceType {
	case "available":
		balance.AvailableAmount = balance.AvailableAmount.Add(amount)
	case "pending":
		balance.PendingAmount = balance.PendingAmount.Add(amount)
	case "reserved":
		balance.ReservedAmount = balance.ReservedAmount.Add(amount)
	default:
		return fmt.Errorf("invalid balance type: %s", balanceType)
	}

	// Recalculate total
	balance.TotalAmount = balance.AvailableAmount.Add(balance.PendingAmount).Add(balance.ReservedAmount)
	balance.LastUpdated = time.Now()

	return nil
}

// DeductBalance deducts from a specific balance type
func (w *Wallet) DeductBalance(currency string, amount decimal.Decimal, balanceType string) error {
	if amount.IsNegative() {
		return fmt.Errorf("amount cannot be negative")
	}

	balance, err := w.GetBalance(currency)
	if err != nil {
		return err
	}

	// Check if sufficient balance exists
	var currentAmount decimal.Decimal
	switch balanceType {
	case "available":
		currentAmount = balance.AvailableAmount
	case "pending":
		currentAmount = balance.PendingAmount
	case "reserved":
		currentAmount = balance.ReservedAmount
	default:
		return fmt.Errorf("invalid balance type: %s", balanceType)
	}

	if currentAmount.LessThan(amount) {
		return fmt.Errorf("insufficient %s balance for %s: need %s, have %s",
			balanceType, currency, amount.String(), currentAmount.String())
	}

	// Deduct from the appropriate balance type
	switch balanceType {
	case "available":
		balance.AvailableAmount = balance.AvailableAmount.Sub(amount)
	case "pending":
		balance.PendingAmount = balance.PendingAmount.Sub(amount)
	case "reserved":
		balance.ReservedAmount = balance.ReservedAmount.Sub(amount)
	}

	// Recalculate total
	balance.TotalAmount = balance.AvailableAmount.Add(balance.PendingAmount).Add(balance.ReservedAmount)
	balance.LastUpdated = time.Now()

	return nil
}

// TransferBalance transfers amount between balance types within the same currency
func (w *Wallet) TransferBalance(currency string, amount decimal.Decimal, fromType, toType string) error {
	if err := w.DeductBalance(currency, amount, fromType); err != nil {
		return err
	}

	if err := w.AddBalance(currency, amount, toType); err != nil {
		// Rollback the deduction
		_ = w.AddBalance(currency, amount, fromType)
		return err
	}

	return nil
}

// ReserveAmount reserves an amount for pending transactions
func (w *Wallet) ReserveAmount(currency string, amount decimal.Decimal) error {
	return w.TransferBalance(currency, amount, "available", "reserved")
}

// ReleaseReserve releases reserved amount back to available
func (w *Wallet) ReleaseReserve(currency string, amount decimal.Decimal) error {
	return w.TransferBalance(currency, amount, "reserved", "available")
}

// ConfirmPending confirms pending amount to available
func (w *Wallet) ConfirmPending(currency string, amount decimal.Decimal) error {
	return w.TransferBalance(currency, amount, "pending", "available")
}

// IsOperational checks if the wallet can process transactions
func (w *Wallet) IsOperational() bool {
	return w.Status.IsOperational()
}

// CanWithdraw checks if amount can be withdrawn considering limits
func (w *Wallet) CanWithdraw(currency string, amount decimal.Decimal) bool {
	if !w.IsOperational() {
		return false
	}

	if !w.HasSufficientBalance(currency, amount) {
		return false
	}

	// Check withdrawal limits (simplified check against daily limit)
	if amount.GreaterThan(w.Security.DailyWithdrawalLimit) {
		return false
	}

	return true
}

// GetSupportedCurrencies returns list of supported currencies
func (w *Wallet) GetSupportedCurrencies() []string {
	return w.Settings.AllowedCurrencies
}

// IsCurrencySupported checks if a currency is supported
func (w *Wallet) IsCurrencySupported(currency string) bool {
	for _, supported := range w.Settings.AllowedCurrencies {
		if supported == currency {
			return true
		}
	}
	return false
}

// UpdateActivity updates wallet activity tracking
func (w *Wallet) UpdateActivity() {
	now := time.Now()
	w.LastActivityAt = &now
	w.TransactionCount++
	w.UpdatedAt = now
}

// GetWalletSummary returns a summary of wallet information
func (w *Wallet) GetWalletSummary() map[string]interface{} {
	summary := map[string]interface{}{
		"id":              w.ID,
		"name":            w.Name,
		"type":            w.Type,
		"status":          w.Status,
		"is_operational":  w.IsOperational(),
		"is_verified":     w.IsVerified,
		"verification_level": w.VerificationLevel,
		"kyc_level":       w.KYCLevel,
		"data_region":     w.DataRegion,
		"default_currency": w.Settings.DefaultCurrency,
		"supported_currencies": w.GetSupportedCurrencies(),
		"transaction_count": w.TransactionCount,
		"total_volume":    w.TotalVolume,
		"created_at":      w.CreatedAt,
		"last_activity":   w.LastActivityAt,
	}

	// Add balance summary
	balanceSummary := make(map[string]interface{})
	for _, balance := range w.Balances {
		balanceSummary[balance.Currency] = map[string]interface{}{
			"available": balance.AvailableAmount,
			"pending":   balance.PendingAmount,
			"reserved":  balance.ReservedAmount,
			"total":     balance.TotalAmount,
		}
	}
	summary["balances"] = balanceSummary

	return summary
}

// GenerateSearchKeywords generates search keywords for the wallet
func (w *Wallet) GenerateSearchKeywords() []string {
	keywords := []string{
		w.Name,
		string(w.Type),
		string(w.Status),
		w.Settings.DefaultCurrency,
		w.VerificationLevel,
		w.KYCLevel,
	}

	// Add supported currencies
	keywords = append(keywords, w.Settings.AllowedCurrencies...)

	// Add tags
	keywords = append(keywords, w.Tags...)

	// Remove duplicates and empty strings
	seen := make(map[string]bool)
	var unique []string
	for _, keyword := range keywords {
		if keyword != "" && !seen[strings.ToLower(keyword)] {
			seen[strings.ToLower(keyword)] = true
			unique = append(unique, strings.ToLower(keyword))
		}
	}

	return unique
}

// MarshalJSON customizes JSON serialization
func (w *Wallet) MarshalJSON() ([]byte, error) {
	type Alias Wallet
	return json.Marshal(&struct {
		*Alias
		IsOperational     bool                   `json:"is_operational"`
		WalletSummary     map[string]interface{} `json:"wallet_summary"`
		SearchKeywords    []string               `json:"search_keywords,omitempty"`
	}{
		Alias:          (*Alias)(w),
		IsOperational:  w.IsOperational(),
		WalletSummary:  w.GetWalletSummary(),
		SearchKeywords: w.GenerateSearchKeywords(),
	})
}

// Helper functions for wallet management

// CreatePersonalWallet creates a personal wallet for a user
func CreatePersonalWallet(userID uuid.UUID, name string, defaultCurrency string) *Wallet {
	return &Wallet{
		UserID: &userID,
		Name:   name,
		Type:   WalletTypePersonal,
		Status: WalletStatusActive,
		Settings: WalletSettings{
			DefaultCurrency:   defaultCurrency,
			AllowedCurrencies: []string{defaultCurrency, "USD"},
		},
		Security: WalletSecurity{
			SecurityLevel:          "basic",
			DailyWithdrawalLimit:   decimal.NewFromInt(1000),
			MonthlyWithdrawalLimit: decimal.NewFromInt(10000),
		},
		Balances: []WalletBalance{},
	}
}

// CreateBusinessWallet creates a business wallet
func CreateBusinessWallet(businessID uuid.UUID, name string, defaultCurrency string, allowedCurrencies []string) *Wallet {
	if len(allowedCurrencies) == 0 {
		allowedCurrencies = []string{defaultCurrency, "USD"}
	}

	return &Wallet{
		BusinessID: &businessID,
		Name:       name,
		Type:       WalletTypeBusiness,
		Status:     WalletStatusActive,
		Settings: WalletSettings{
			DefaultCurrency:   defaultCurrency,
			AllowedCurrencies: allowedCurrencies,
		},
		Security: WalletSecurity{
			SecurityLevel:          "enhanced",
			DailyWithdrawalLimit:   decimal.NewFromInt(50000),   // Higher limits for business
			MonthlyWithdrawalLimit: decimal.NewFromInt(500000),
		},
		Balances: []WalletBalance{},
	}
}