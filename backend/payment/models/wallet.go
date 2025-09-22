package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Wallet represents a user's payment wallet for a specific currency
type Wallet struct {
	ID               uuid.UUID     `json:"id" db:"id"`
	UserID           uuid.UUID     `json:"user_id" db:"user_id"`
	Balance          int64         `json:"balance" db:"balance"`                   // Amount in cents
	Currency         Currency      `json:"currency" db:"currency"`
	FrozenBalance    int64         `json:"frozen_balance" db:"frozen_balance"`     // Frozen amount in cents
	DailyLimit       int64         `json:"daily_limit" db:"daily_limit"`           // Daily limit in cents
	MonthlyLimit     int64         `json:"monthly_limit" db:"monthly_limit"`       // Monthly limit in cents
	UsedThisDay      int64         `json:"used_this_day" db:"used_this_day"`       // Used today in cents
	UsedThisMonth    int64         `json:"used_this_month" db:"used_this_month"`   // Used this month in cents
	LastResetDay     time.Time     `json:"last_reset_day" db:"last_reset_day"`     // Last daily reset
	LastResetMonth   time.Time     `json:"last_reset_month" db:"last_reset_month"` // Last monthly reset
	Status           WalletStatus  `json:"status" db:"status"`
	IsPrimary        bool          `json:"is_primary" db:"is_primary"`
	CreatedAt        time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at" db:"updated_at"`
}

// Currency represents supported currencies in Southeast Asia
type Currency string

const (
	CurrencyTHB Currency = "THB" // Thai Baht
	CurrencySGD Currency = "SGD" // Singapore Dollar
	CurrencyIDR Currency = "IDR" // Indonesian Rupiah
	CurrencyMYR Currency = "MYR" // Malaysian Ringgit
	CurrencyPHP Currency = "PHP" // Philippine Peso
	CurrencyVND Currency = "VND" // Vietnamese Dong
	CurrencyUSD Currency = "USD" // US Dollar (reference currency)
)

// WalletStatus represents the current status of a wallet
type WalletStatus string

const (
	WalletStatusActive    WalletStatus = "active"
	WalletStatusSuspended WalletStatus = "suspended"
	WalletStatusClosed    WalletStatus = "closed"
	WalletStatusFrozen    WalletStatus = "frozen"
)

// Default limits by currency (in cents)
var DefaultDailyLimits = map[Currency]int64{
	CurrencyTHB: 10000000, // 100,000 THB
	CurrencySGD: 500000,   // 5,000 SGD
	CurrencyIDR: 1500000000, // 15,000,000 IDR
	CurrencyMYR: 2000000,  // 20,000 MYR
	CurrencyPHP: 25000000, // 250,000 PHP
	CurrencyVND: 23000000000, // 230,000,000 VND
	CurrencyUSD: 300000,   // 3,000 USD
}

var DefaultMonthlyLimits = map[Currency]int64{
	CurrencyTHB: 300000000, // 3,000,000 THB
	CurrencySGD: 15000000,  // 150,000 SGD
	CurrencyIDR: 45000000000, // 450,000,000 IDR
	CurrencyMYR: 60000000,  // 600,000 MYR
	CurrencyPHP: 750000000, // 7,500,000 PHP
	CurrencyVND: 690000000000, // 6,900,000,000 VND
	CurrencyUSD: 9000000,   // 90,000 USD
}

// ValidCurrencies returns all supported currencies
func ValidCurrencies() []Currency {
	return []Currency{
		CurrencyTHB,
		CurrencySGD,
		CurrencyIDR,
		CurrencyMYR,
		CurrencyPHP,
		CurrencyVND,
		CurrencyUSD,
	}
}

// IsValid validates if the currency is supported
func (c Currency) IsValid() bool {
	for _, valid := range ValidCurrencies() {
		if c == valid {
			return true
		}
	}
	return false
}

// String returns the string representation of Currency
func (c Currency) String() string {
	return string(c)
}

// GetSymbol returns the currency symbol
func (c Currency) GetSymbol() string {
	symbols := map[Currency]string{
		CurrencyTHB: "฿",
		CurrencySGD: "S$",
		CurrencyIDR: "Rp",
		CurrencyMYR: "RM",
		CurrencyPHP: "₱",
		CurrencyVND: "₫",
		CurrencyUSD: "$",
	}
	return symbols[c]
}

// GetDecimalPlaces returns the number of decimal places for the currency
func (c Currency) GetDecimalPlaces() int {
	// Most currencies use 2 decimal places, but some Asian currencies don't use decimals
	switch c {
	case CurrencyIDR, CurrencyVND:
		return 0 // No decimal places
	default:
		return 2 // Standard decimal places
	}
}

// FormatAmount formats an amount in cents to a human-readable string
func (c Currency) FormatAmount(amountInCents int64) string {
	decimalPlaces := c.GetDecimalPlaces()
	symbol := c.GetSymbol()

	if decimalPlaces == 0 {
		return fmt.Sprintf("%s %d", symbol, amountInCents)
	}

	divisor := int64(100) // Convert cents to main currency unit
	whole := amountInCents / divisor
	fraction := amountInCents % divisor

	return fmt.Sprintf("%s %d.%02d", symbol, whole, fraction)
}

// Value implements the driver.Valuer interface for database storage
func (c Currency) Value() (driver.Value, error) {
	return string(c), nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (c *Currency) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if str, ok := value.(string); ok {
		*c = Currency(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into Currency", value)
}

// IsValid validates if the wallet status is supported
func (ws WalletStatus) IsValid() bool {
	validStatuses := []WalletStatus{
		WalletStatusActive,
		WalletStatusSuspended,
		WalletStatusClosed,
		WalletStatusFrozen,
	}
	for _, valid := range validStatuses {
		if ws == valid {
			return true
		}
	}
	return false
}

// String returns the string representation of WalletStatus
func (ws WalletStatus) String() string {
	return string(ws)
}

// CanTransact checks if wallet can perform transactions
func (ws WalletStatus) CanTransact() bool {
	return ws == WalletStatusActive
}

// CanReceive checks if wallet can receive funds
func (ws WalletStatus) CanReceive() bool {
	return ws == WalletStatusActive || ws == WalletStatusSuspended
}

// Validate performs comprehensive validation on the Wallet model
func (w *Wallet) Validate() error {
	var errs []string

	// User ID validation
	if w.UserID == uuid.Nil {
		errs = append(errs, "user_id is required")
	}

	// Currency validation
	if !w.Currency.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid currency: %s", w.Currency))
	}

	// Balance validation
	if w.Balance < 0 {
		errs = append(errs, "balance cannot be negative")
	}

	// Frozen balance validation
	if w.FrozenBalance < 0 {
		errs = append(errs, "frozen_balance cannot be negative")
	}

	if w.FrozenBalance > w.Balance {
		errs = append(errs, "frozen_balance cannot exceed total balance")
	}

	// Limit validation
	if w.DailyLimit < 0 {
		errs = append(errs, "daily_limit cannot be negative")
	}

	if w.MonthlyLimit < 0 {
		errs = append(errs, "monthly_limit cannot be negative")
	}

	if w.DailyLimit > w.MonthlyLimit && w.MonthlyLimit > 0 {
		errs = append(errs, "daily_limit cannot exceed monthly_limit")
	}

	// Usage validation
	if w.UsedThisDay < 0 {
		errs = append(errs, "used_this_day cannot be negative")
	}

	if w.UsedThisMonth < 0 {
		errs = append(errs, "used_this_month cannot be negative")
	}

	// Status validation
	if !w.Status.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid wallet status: %s", w.Status))
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

// BeforeCreate sets up the wallet before database creation
func (w *Wallet) BeforeCreate() error {
	// Generate UUID if not set
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now().UTC()
	w.CreatedAt = now
	w.UpdatedAt = now

	// Set default values
	if w.Status == "" {
		w.Status = WalletStatusActive
	}

	// Set default limits based on currency
	if w.DailyLimit == 0 {
		if defaultLimit, exists := DefaultDailyLimits[w.Currency]; exists {
			w.DailyLimit = defaultLimit
		}
	}

	if w.MonthlyLimit == 0 {
		if defaultLimit, exists := DefaultMonthlyLimits[w.Currency]; exists {
			w.MonthlyLimit = defaultLimit
		}
	}

	// Initialize reset timestamps
	if w.LastResetDay.IsZero() {
		w.LastResetDay = now
	}

	if w.LastResetMonth.IsZero() {
		w.LastResetMonth = now
	}

	// Validate before creation
	return w.Validate()
}

// BeforeUpdate sets up the wallet before database update
func (w *Wallet) BeforeUpdate() error {
	// Update timestamp
	w.UpdatedAt = time.Now().UTC()

	// Check if daily/monthly usage needs reset
	w.CheckAndResetUsage()

	// Validate before update
	return w.Validate()
}

// CheckAndResetUsage resets daily and monthly usage if needed
func (w *Wallet) CheckAndResetUsage() {
	now := time.Now().UTC()

	// Reset daily usage if it's a new day
	if w.LastResetDay.Day() != now.Day() || w.LastResetDay.Month() != now.Month() || w.LastResetDay.Year() != now.Year() {
		w.UsedThisDay = 0
		w.LastResetDay = now
	}

	// Reset monthly usage if it's a new month
	if w.LastResetMonth.Month() != now.Month() || w.LastResetMonth.Year() != now.Year() {
		w.UsedThisMonth = 0
		w.LastResetMonth = now
	}
}

// GetAvailableBalance returns the available balance (total - frozen)
func (w *Wallet) GetAvailableBalance() int64 {
	return w.Balance - w.FrozenBalance
}

// CanSend checks if wallet can send a specific amount
func (w *Wallet) CanSend(amount int64) error {
	if !w.Status.CanTransact() {
		return fmt.Errorf("wallet is %s and cannot send funds", w.Status)
	}

	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	// Check available balance
	availableBalance := w.GetAvailableBalance()
	if amount > availableBalance {
		return fmt.Errorf("insufficient funds: available %d, requested %d", availableBalance, amount)
	}

	// Check daily limit
	w.CheckAndResetUsage()
	if w.DailyLimit > 0 && (w.UsedThisDay+amount) > w.DailyLimit {
		return fmt.Errorf("daily limit exceeded: used %d, limit %d, requested %d",
			w.UsedThisDay, w.DailyLimit, amount)
	}

	// Check monthly limit
	if w.MonthlyLimit > 0 && (w.UsedThisMonth+amount) > w.MonthlyLimit {
		return fmt.Errorf("monthly limit exceeded: used %d, limit %d, requested %d",
			w.UsedThisMonth, w.MonthlyLimit, amount)
	}

	return nil
}

// CanReceive checks if wallet can receive funds
func (w *Wallet) CanReceive(amount int64) error {
	if !w.Status.CanReceive() {
		return fmt.Errorf("wallet is %s and cannot receive funds", w.Status)
	}

	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	// Check for overflow
	if w.Balance > 0 && amount > (9223372036854775807-w.Balance) {
		return errors.New("transaction would cause balance overflow")
	}

	return nil
}

// DebitAmount debits an amount from the wallet with limit tracking
func (w *Wallet) DebitAmount(amount int64) error {
	if err := w.CanSend(amount); err != nil {
		return err
	}

	w.Balance -= amount
	w.UsedThisDay += amount
	w.UsedThisMonth += amount
	w.UpdatedAt = time.Now().UTC()

	return w.Validate()
}

// CreditAmount credits an amount to the wallet
func (w *Wallet) CreditAmount(amount int64) error {
	if err := w.CanReceive(amount); err != nil {
		return err
	}

	w.Balance += amount
	w.UpdatedAt = time.Now().UTC()

	return w.Validate()
}

// FreezeAmount freezes a specific amount in the wallet
func (w *Wallet) FreezeAmount(amount int64) error {
	if amount <= 0 {
		return errors.New("freeze amount must be positive")
	}

	if amount > w.GetAvailableBalance() {
		return errors.New("insufficient available balance to freeze")
	}

	w.FrozenBalance += amount
	w.UpdatedAt = time.Now().UTC()

	return w.Validate()
}

// UnfreezeAmount unfreezes a specific amount in the wallet
func (w *Wallet) UnfreezeAmount(amount int64) error {
	if amount <= 0 {
		return errors.New("unfreeze amount must be positive")
	}

	if amount > w.FrozenBalance {
		return errors.New("cannot unfreeze more than frozen balance")
	}

	w.FrozenBalance -= amount
	w.UpdatedAt = time.Now().UTC()

	return w.Validate()
}

// SetPrimary marks this wallet as the primary wallet
func (w *Wallet) SetPrimary() error {
	w.IsPrimary = true
	w.UpdatedAt = time.Now().UTC()
	return nil
}

// Suspend suspends the wallet
func (w *Wallet) Suspend() error {
	if w.Status == WalletStatusClosed {
		return errors.New("cannot suspend a closed wallet")
	}

	w.Status = WalletStatusSuspended
	w.UpdatedAt = time.Now().UTC()

	return w.Validate()
}

// Activate activates a suspended wallet
func (w *Wallet) Activate() error {
	if w.Status == WalletStatusClosed {
		return errors.New("cannot activate a closed wallet")
	}

	w.Status = WalletStatusActive
	w.UpdatedAt = time.Now().UTC()

	return w.Validate()
}

// Freeze freezes the wallet
func (w *Wallet) Freeze() error {
	if w.Status == WalletStatusClosed {
		return errors.New("cannot freeze a closed wallet")
	}

	w.Status = WalletStatusFrozen
	w.UpdatedAt = time.Now().UTC()

	return w.Validate()
}

// Close closes the wallet permanently
func (w *Wallet) Close() error {
	if w.Balance > 0 {
		return errors.New("cannot close wallet with positive balance")
	}

	if w.FrozenBalance > 0 {
		return errors.New("cannot close wallet with frozen balance")
	}

	w.Status = WalletStatusClosed
	w.UpdatedAt = time.Now().UTC()

	return w.Validate()
}

// UpdateLimits updates the daily and monthly limits
func (w *Wallet) UpdateLimits(dailyLimit, monthlyLimit int64) error {
	if dailyLimit < 0 || monthlyLimit < 0 {
		return errors.New("limits cannot be negative")
	}

	if dailyLimit > monthlyLimit && monthlyLimit > 0 {
		return errors.New("daily limit cannot exceed monthly limit")
	}

	w.DailyLimit = dailyLimit
	w.MonthlyLimit = monthlyLimit
	w.UpdatedAt = time.Now().UTC()

	return w.Validate()
}

// GetRemainingDailyLimit returns the remaining daily limit
func (w *Wallet) GetRemainingDailyLimit() int64 {
	w.CheckAndResetUsage()
	if w.DailyLimit == 0 {
		return 0 // No limit
	}
	remaining := w.DailyLimit - w.UsedThisDay
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetRemainingMonthlyLimit returns the remaining monthly limit
func (w *Wallet) GetRemainingMonthlyLimit() int64 {
	w.CheckAndResetUsage()
	if w.MonthlyLimit == 0 {
		return 0 // No limit
	}
	remaining := w.MonthlyLimit - w.UsedThisMonth
	if remaining < 0 {
		return 0
	}
	return remaining
}

// ToPublicWallet returns a sanitized version for public API responses
func (w *Wallet) ToPublicWallet() map[string]interface{} {
	return map[string]interface{}{
		"id":                w.ID,
		"balance":           w.Balance,
		"currency":          w.Currency,
		"currency_symbol":   w.Currency.GetSymbol(),
		"available_balance": w.GetAvailableBalance(),
		"frozen_balance":    w.FrozenBalance,
		"daily_limit":       w.DailyLimit,
		"monthly_limit":     w.MonthlyLimit,
		"remaining_daily":   w.GetRemainingDailyLimit(),
		"remaining_monthly": w.GetRemainingMonthlyLimit(),
		"status":            w.Status,
		"is_primary":        w.IsPrimary,
		"created_at":        w.CreatedAt,
		"updated_at":        w.UpdatedAt,
	}
}

// ToBalanceInfo returns basic balance information
func (w *Wallet) ToBalanceInfo() map[string]interface{} {
	return map[string]interface{}{
		"currency":          w.Currency,
		"balance":           w.Balance,
		"available_balance": w.GetAvailableBalance(),
		"formatted_balance": w.Currency.FormatAmount(w.Balance),
	}
}

// WalletCreateRequest represents a request to create a new wallet
type WalletCreateRequest struct {
	Currency  Currency `json:"currency" validate:"required"`
	IsPrimary bool     `json:"is_primary"`
}

// ToWallet converts a create request to a Wallet model
func (req *WalletCreateRequest) ToWallet(userID uuid.UUID) *Wallet {
	return &Wallet{
		UserID:    userID,
		Currency:  req.Currency,
		IsPrimary: req.IsPrimary,
	}
}

// WalletManager provides wallet management utilities
type WalletManager struct {
	// Add dependencies like database, cache, etc.
}

// NewWalletManager creates a new wallet manager
func NewWalletManager() *WalletManager {
	return &WalletManager{}
}

// CreateWallet creates a new wallet with proper validation
func (wm *WalletManager) CreateWallet(req *WalletCreateRequest, userID uuid.UUID) (*Wallet, error) {
	wallet := req.ToWallet(userID)

	if err := wallet.BeforeCreate(); err != nil {
		return nil, fmt.Errorf("wallet creation failed: %v", err)
	}

	return wallet, nil
}