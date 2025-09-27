package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// PaymentMethod represents a payment method linked to a user's wallet
type PaymentMethod struct {
	ID              uuid.UUID             `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID          uuid.UUID             `json:"user_id" gorm:"type:varchar(36);not null;index"`
	WalletID        *uuid.UUID            `json:"wallet_id,omitempty" gorm:"type:varchar(36);index"`
	Type            PaymentMethodType     `json:"type" gorm:"type:varchar(30);not null"`
	Provider        PaymentProvider       `json:"provider" gorm:"type:varchar(50);not null"`
	Status          PaymentMethodStatus   `json:"status" gorm:"type:varchar(20);default:'active'"`
	IsDefault       bool                  `json:"is_default" gorm:"default:false"`
	IsVerified      bool                  `json:"is_verified" gorm:"default:false"`
	DisplayName     string                `json:"display_name" gorm:"type:varchar(100)"`
	LastFourDigits  *string               `json:"last_four_digits,omitempty" gorm:"type:varchar(4)"`
	ExpiryMonth     *int                  `json:"expiry_month,omitempty" gorm:"type:int"`
	ExpiryYear      *int                  `json:"expiry_year,omitempty" gorm:"type:int"`
	BrandName       *string               `json:"brand_name,omitempty" gorm:"type:varchar(50)"`
	Country         string                `json:"country" gorm:"type:varchar(2);not null"` // ISO country code
	Currency        Currency              `json:"currency" gorm:"type:varchar(3);not null"`
	Metadata        PaymentMethodMetadata `json:"metadata" gorm:"type:json"`
	ProviderData    ProviderData          `json:"provider_data" gorm:"type:json"`
	SecurityInfo    SecurityInfo          `json:"security_info" gorm:"type:json"`
	UsageStats      UsageStats            `json:"usage_stats" gorm:"type:json"`
	ExternalID      *string               `json:"external_id,omitempty" gorm:"type:varchar(255);index"`
	LastUsedAt      *time.Time            `json:"last_used_at,omitempty" gorm:"index"`
	VerifiedAt      *time.Time            `json:"verified_at,omitempty"`
	ExpiresAt       *time.Time            `json:"expires_at,omitempty"`
	CreatedAt       time.Time             `json:"created_at" gorm:"not null"`
	UpdatedAt       time.Time             `json:"updated_at" gorm:"not null"`
}

// PaymentMethodType represents the type of payment method
type PaymentMethodType string

const (
	PaymentMethodTypeCreditCard      PaymentMethodType = "credit_card"
	PaymentMethodTypeDebitCard       PaymentMethodType = "debit_card"
	PaymentMethodTypeBankAccount     PaymentMethodType = "bank_account"
	PaymentMethodTypeEWallet         PaymentMethodType = "e_wallet"
	PaymentMethodTypeMobileWallet    PaymentMethodType = "mobile_wallet"
	PaymentMethodTypeDigitalPayment  PaymentMethodType = "digital_payment"
	PaymentMethodTypeCrypto          PaymentMethodType = "crypto"
	PaymentMethodTypePrepaidCard     PaymentMethodType = "prepaid_card"
	PaymentMethodTypeGiftCard        PaymentMethodType = "gift_card"
	PaymentMethodTypeBuyNowPayLater  PaymentMethodType = "buy_now_pay_later"
)

// PaymentProvider represents the payment service provider
type PaymentProvider string

const (
	// Southeast Asian providers
	PaymentProviderPromptPay     PaymentProvider = "promptpay"      // Thailand
	PaymentProviderTrueMoney     PaymentProvider = "truemoney"      // Thailand
	PaymentProviderShopeePay     PaymentProvider = "shopeepay"      // Regional
	PaymentProviderGrabPay       PaymentProvider = "grabpay"        // Regional
	PaymentProviderGoPay         PaymentProvider = "gopay"          // Indonesia
	PaymentProviderOVO           PaymentProvider = "ovo"            // Indonesia
	PaymentProviderDANA          PaymentProvider = "dana"           // Indonesia
	PaymentProviderTouchNGo      PaymentProvider = "touchngo"       // Malaysia
	PaymentProviderBoost         PaymentProvider = "boost"          // Malaysia
	PaymentProviderGCash         PaymentProvider = "gcash"          // Philippines
	PaymentProviderPayMaya       PaymentProvider = "paymaya"        // Philippines
	PaymentProviderZaloPay       PaymentProvider = "zalopay"        // Vietnam
	PaymentProviderMomoPay       PaymentProvider = "momopay"        // Vietnam
	PaymentProviderDBS           PaymentProvider = "dbs"            // Singapore
	PaymentProviderOCBC          PaymentProvider = "ocbc"           // Singapore
	PaymentProviderUOB           PaymentProvider = "uob"            // Singapore

	// International providers
	PaymentProviderVisa          PaymentProvider = "visa"
	PaymentProviderMastercard    PaymentProvider = "mastercard"
	PaymentProviderAmex          PaymentProvider = "amex"
	PaymentProviderPayPal        PaymentProvider = "paypal"
	PaymentProviderStripe        PaymentProvider = "stripe"
	PaymentProviderAdyen         PaymentProvider = "adyen"

	// Crypto providers
	PaymentProviderBinance       PaymentProvider = "binance"
	PaymentProviderCoinbase      PaymentProvider = "coinbase"
	PaymentProviderMetaMask      PaymentProvider = "metamask"
)

// PaymentMethodStatus represents the current status of a payment method
type PaymentMethodStatus string

const (
	PaymentMethodStatusActive    PaymentMethodStatus = "active"
	PaymentMethodStatusInactive  PaymentMethodStatus = "inactive"
	PaymentMethodStatusExpired   PaymentMethodStatus = "expired"
	PaymentMethodStatusBlocked   PaymentMethodStatus = "blocked"
	PaymentMethodStatusPending   PaymentMethodStatus = "pending"
	PaymentMethodStatusFailed    PaymentMethodStatus = "failed"
)

// PaymentMethodMetadata represents additional metadata for payment methods
type PaymentMethodMetadata struct {
	BillingAddress   *BillingAddress `json:"billing_address,omitempty"`
	HolderName       string          `json:"holder_name,omitempty"`
	BankName         string          `json:"bank_name,omitempty"`
	BankCode         string          `json:"bank_code,omitempty"`
	AccountType      string          `json:"account_type,omitempty"`     // savings, checking, business
	AccountNumber    string          `json:"account_number,omitempty"`   // Masked/encrypted
	RoutingNumber    string          `json:"routing_number,omitempty"`
	SwiftCode        string          `json:"swift_code,omitempty"`
	PhoneNumber      string          `json:"phone_number,omitempty"`     // For mobile wallets
	EmailAddress     string          `json:"email_address,omitempty"`    // For digital payments
	WalletAddress    string          `json:"wallet_address,omitempty"`   // For crypto
	NetworkType      string          `json:"network_type,omitempty"`     // For crypto (ethereum, bitcoin, etc.)
	TokenSymbol      string          `json:"token_symbol,omitempty"`     // For crypto tokens
	InstallmentPlans []string        `json:"installment_plans,omitempty"` // For BNPL
	CreditLimit      *int64          `json:"credit_limit,omitempty"`     // For credit products
	AvailableCredit  *int64          `json:"available_credit,omitempty"` // For credit products
}

// BillingAddress represents billing address information
type BillingAddress struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Company     string `json:"company,omitempty"`
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	Country     string `json:"country,omitempty"` // ISO country code
}

// ProviderData represents provider-specific data (encrypted/tokenized)
type ProviderData struct {
	TokenID           string            `json:"token_id,omitempty"`
	ProviderAccountID string            `json:"provider_account_id,omitempty"`
	ProviderUserID    string            `json:"provider_user_id,omitempty"`
	TokenType         string            `json:"token_type,omitempty"`
	Fingerprint       string            `json:"fingerprint,omitempty"`
	CustomData        map[string]string `json:"custom_data,omitempty"`
}

// SecurityInfo represents security-related information
type SecurityInfo struct {
	CVVVerified       bool      `json:"cvv_verified"`
	AddressVerified   bool      `json:"address_verified"`
	PhoneVerified     bool      `json:"phone_verified"`
	EmailVerified     bool      `json:"email_verified"`
	RiskScore         float64   `json:"risk_score,omitempty"`       // 0-100
	FraudFlags        []string  `json:"fraud_flags,omitempty"`
	LastSecurityCheck time.Time `json:"last_security_check,omitempty"`
	SecurityLevel     string    `json:"security_level,omitempty"`   // low, medium, high
}

// UsageStats represents usage statistics for the payment method
type UsageStats struct {
	TotalTransactions      int       `json:"total_transactions"`
	TotalAmountProcessed   int64     `json:"total_amount_processed"`   // In cents
	SuccessfulTransactions int       `json:"successful_transactions"`
	FailedTransactions     int       `json:"failed_transactions"`
	LastTransactionDate    time.Time `json:"last_transaction_date,omitempty"`
	AverageTransactionSize int64     `json:"average_transaction_size"` // In cents
	MonthlyUsageCount      int       `json:"monthly_usage_count"`
	PreferredCurrency      Currency  `json:"preferred_currency,omitempty"`
}

// Value implements the driver.Valuer interface for PaymentMethodMetadata
func (pmm PaymentMethodMetadata) Value() (driver.Value, error) {
	return json.Marshal(pmm)
}

// Scan implements the sql.Scanner interface for PaymentMethodMetadata
func (pmm *PaymentMethodMetadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into PaymentMethodMetadata", value)
	}

	return json.Unmarshal(jsonData, pmm)
}

// Value implements the driver.Valuer interface for ProviderData
func (pd ProviderData) Value() (driver.Value, error) {
	return json.Marshal(pd)
}

// Scan implements the sql.Scanner interface for ProviderData
func (pd *ProviderData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into ProviderData", value)
	}

	return json.Unmarshal(jsonData, pd)
}

// Value implements the driver.Valuer interface for SecurityInfo
func (si SecurityInfo) Value() (driver.Value, error) {
	return json.Marshal(si)
}

// Scan implements the sql.Scanner interface for SecurityInfo
func (si *SecurityInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into SecurityInfo", value)
	}

	return json.Unmarshal(jsonData, si)
}

// Value implements the driver.Valuer interface for UsageStats
func (us UsageStats) Value() (driver.Value, error) {
	return json.Marshal(us)
}

// Scan implements the sql.Scanner interface for UsageStats
func (us *UsageStats) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into UsageStats", value)
	}

	return json.Unmarshal(jsonData, us)
}

// ValidPaymentMethodTypes returns all supported payment method types
func ValidPaymentMethodTypes() []PaymentMethodType {
	return []PaymentMethodType{
		PaymentMethodTypeCreditCard,
		PaymentMethodTypeDebitCard,
		PaymentMethodTypeBankAccount,
		PaymentMethodTypeEWallet,
		PaymentMethodTypeMobileWallet,
		PaymentMethodTypeDigitalPayment,
		PaymentMethodTypeCrypto,
		PaymentMethodTypePrepaidCard,
		PaymentMethodTypeGiftCard,
		PaymentMethodTypeBuyNowPayLater,
	}
}

// ValidPaymentProviders returns all supported payment providers
func ValidPaymentProviders() []PaymentProvider {
	return []PaymentProvider{
		PaymentProviderPromptPay, PaymentProviderTrueMoney, PaymentProviderShopeePay,
		PaymentProviderGrabPay, PaymentProviderGoPay, PaymentProviderOVO, PaymentProviderDANA,
		PaymentProviderTouchNGo, PaymentProviderBoost, PaymentProviderGCash, PaymentProviderPayMaya,
		PaymentProviderZaloPay, PaymentProviderMomoPay, PaymentProviderDBS, PaymentProviderOCBC,
		PaymentProviderUOB, PaymentProviderVisa, PaymentProviderMastercard, PaymentProviderAmex,
		PaymentProviderPayPal, PaymentProviderStripe, PaymentProviderAdyen, PaymentProviderBinance,
		PaymentProviderCoinbase, PaymentProviderMetaMask,
	}
}

// ValidPaymentMethodStatuses returns all supported payment method statuses
func ValidPaymentMethodStatuses() []PaymentMethodStatus {
	return []PaymentMethodStatus{
		PaymentMethodStatusActive,
		PaymentMethodStatusInactive,
		PaymentMethodStatusExpired,
		PaymentMethodStatusBlocked,
		PaymentMethodStatusPending,
		PaymentMethodStatusFailed,
	}
}

// IsValid validates if the payment method type is supported
func (pmt PaymentMethodType) IsValid() bool {
	for _, valid := range ValidPaymentMethodTypes() {
		if pmt == valid {
			return true
		}
	}
	return false
}

// IsValid validates if the payment provider is supported
func (pp PaymentProvider) IsValid() bool {
	for _, valid := range ValidPaymentProviders() {
		if pp == valid {
			return true
		}
	}
	return false
}

// IsValid validates if the payment method status is supported
func (pms PaymentMethodStatus) IsValid() bool {
	for _, valid := range ValidPaymentMethodStatuses() {
		if pms == valid {
			return true
		}
	}
	return false
}

// String returns the string representation of PaymentMethodType
func (pmt PaymentMethodType) String() string {
	return string(pmt)
}

// String returns the string representation of PaymentProvider
func (pp PaymentProvider) String() string {
	return string(pp)
}

// String returns the string representation of PaymentMethodStatus
func (pms PaymentMethodStatus) String() string {
	return string(pms)
}

// RequiresCardData checks if this payment method type requires card information
func (pmt PaymentMethodType) RequiresCardData() bool {
	cardTypes := []PaymentMethodType{
		PaymentMethodTypeCreditCard,
		PaymentMethodTypeDebitCard,
		PaymentMethodTypePrepaidCard,
		PaymentMethodTypeGiftCard,
	}
	for _, cardType := range cardTypes {
		if pmt == cardType {
			return true
		}
	}
	return false
}

// RequiresBankData checks if this payment method type requires bank information
func (pmt PaymentMethodType) RequiresBankData() bool {
	return pmt == PaymentMethodTypeBankAccount
}

// RequiresDigitalData checks if this payment method type requires digital account information
func (pmt PaymentMethodType) RequiresDigitalData() bool {
	digitalTypes := []PaymentMethodType{
		PaymentMethodTypeEWallet,
		PaymentMethodTypeMobileWallet,
		PaymentMethodTypeDigitalPayment,
		PaymentMethodTypeCrypto,
	}
	for _, digitalType := range digitalTypes {
		if pmt == digitalType {
			return true
		}
	}
	return false
}

// GetRegion returns the primary region for this payment provider
func (pp PaymentProvider) GetRegion() string {
	seaProviders := map[PaymentProvider]string{
		PaymentProviderPromptPay: "TH",
		PaymentProviderTrueMoney: "TH",
		PaymentProviderGoPay:     "ID",
		PaymentProviderOVO:       "ID",
		PaymentProviderDANA:      "ID",
		PaymentProviderTouchNGo:  "MY",
		PaymentProviderBoost:     "MY",
		PaymentProviderGCash:     "PH",
		PaymentProviderPayMaya:   "PH",
		PaymentProviderZaloPay:   "VN",
		PaymentProviderMomoPay:   "VN",
		PaymentProviderDBS:       "SG",
		PaymentProviderOCBC:      "SG",
		PaymentProviderUOB:       "SG",
	}

	if region, exists := seaProviders[pp]; exists {
		return region
	}
	return "GLOBAL" // For international providers
}

// CanTransact checks if this payment method can process transactions
func (pms PaymentMethodStatus) CanTransact() bool {
	return pms == PaymentMethodStatusActive
}

// CanReceive checks if this payment method can receive funds
func (pms PaymentMethodStatus) CanReceive() bool {
	return pms == PaymentMethodStatusActive
}

// IsExpirable checks if this payment method type can expire
func (pmt PaymentMethodType) IsExpirable() bool {
	expirableTypes := []PaymentMethodType{
		PaymentMethodTypeCreditCard,
		PaymentMethodTypeDebitCard,
		PaymentMethodTypePrepaidCard,
		PaymentMethodTypeGiftCard,
	}
	for _, expirableType := range expirableTypes {
		if pmt == expirableType {
			return true
		}
	}
	return false
}

// Validate performs comprehensive validation on the PaymentMethod model
func (pm *PaymentMethod) Validate() error {
	var errs []string

	// User ID validation
	if pm.UserID == uuid.Nil {
		errs = append(errs, "user_id is required")
	}

	// Type validation
	if !pm.Type.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid payment method type: %s", pm.Type))
	}

	// Provider validation
	if !pm.Provider.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid payment provider: %s", pm.Provider))
	}

	// Status validation
	if !pm.Status.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid payment method status: %s", pm.Status))
	}

	// Currency validation
	if !pm.Currency.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid currency: %s", pm.Currency))
	}

	// Country validation
	if len(pm.Country) != 2 {
		errs = append(errs, "country must be a 2-letter ISO code")
	}

	// Display name validation
	if strings.TrimSpace(pm.DisplayName) == "" {
		errs = append(errs, "display_name is required")
	}
	if len(pm.DisplayName) > 100 {
		errs = append(errs, "display_name cannot exceed 100 characters")
	}

	// Card-specific validation
	if pm.Type.RequiresCardData() {
		if err := pm.validateCardData(); err != nil {
			errs = append(errs, err.Error())
		}
	}

	// Bank-specific validation
	if pm.Type.RequiresBankData() {
		if err := pm.validateBankData(); err != nil {
			errs = append(errs, err.Error())
		}
	}

	// Digital payment validation
	if pm.Type.RequiresDigitalData() {
		if err := pm.validateDigitalData(); err != nil {
			errs = append(errs, err.Error())
		}
	}

	// Expiry validation for expirable types
	if pm.Type.IsExpirable() {
		if err := pm.validateExpiry(); err != nil {
			errs = append(errs, err.Error())
		}
	}

	// Security info validation
	if err := pm.validateSecurityInfo(); err != nil {
		errs = append(errs, err.Error())
	}

	// Provider compatibility validation
	if err := pm.validateProviderCompatibility(); err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

// validateCardData validates card-specific data
func (pm *PaymentMethod) validateCardData() error {
	// Last four digits validation
	if pm.LastFourDigits == nil || len(*pm.LastFourDigits) != 4 {
		return errors.New("last_four_digits is required and must be 4 digits for card types")
	}

	// Validate digits only
	digitRegex := regexp.MustCompile(`^\d{4}$`)
	if !digitRegex.MatchString(*pm.LastFourDigits) {
		return errors.New("last_four_digits must contain only digits")
	}

	return nil
}

// validateBankData validates bank account specific data
func (pm *PaymentMethod) validateBankData() error {
	metadata := pm.Metadata

	// Bank name is required for bank accounts
	if strings.TrimSpace(metadata.BankName) == "" {
		return errors.New("bank_name is required for bank account types")
	}

	// Account type validation
	validAccountTypes := []string{"savings", "checking", "business", "current"}
	if metadata.AccountType != "" {
		valid := false
		for _, validType := range validAccountTypes {
			if metadata.AccountType == validType {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid account_type: %s", metadata.AccountType)
		}
	}

	return nil
}

// validateDigitalData validates digital payment method data
func (pm *PaymentMethod) validateDigitalData() error {
	metadata := pm.Metadata

	// For mobile wallets, phone number is often required
	if pm.Type == PaymentMethodTypeMobileWallet {
		if strings.TrimSpace(metadata.PhoneNumber) == "" {
			return errors.New("phone_number is required for mobile wallet types")
		}
		// Basic phone validation
		phoneRegex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
		if !phoneRegex.MatchString(metadata.PhoneNumber) {
			return errors.New("phone_number must be in international format (+1234567890)")
		}
	}

	// For crypto, wallet address is required
	if pm.Type == PaymentMethodTypeCrypto {
		if strings.TrimSpace(metadata.WalletAddress) == "" {
			return errors.New("wallet_address is required for crypto types")
		}
		if strings.TrimSpace(metadata.NetworkType) == "" {
			return errors.New("network_type is required for crypto types")
		}
	}

	return nil
}

// validateExpiry validates expiry date for expirable payment methods
func (pm *PaymentMethod) validateExpiry() error {
	if pm.ExpiryMonth == nil || pm.ExpiryYear == nil {
		return errors.New("expiry_month and expiry_year are required for expirable payment methods")
	}

	// Validate month range
	if *pm.ExpiryMonth < 1 || *pm.ExpiryMonth > 12 {
		return errors.New("expiry_month must be between 1 and 12")
	}

	// Validate year range (current year to 20 years in future)
	currentYear := time.Now().Year()
	if *pm.ExpiryYear < currentYear || *pm.ExpiryYear > currentYear+20 {
		return fmt.Errorf("expiry_year must be between %d and %d", currentYear, currentYear+20)
	}

	// Check if already expired
	now := time.Now()
	if *pm.ExpiryYear < now.Year() || (*pm.ExpiryYear == now.Year() && *pm.ExpiryMonth < int(now.Month())) {
		return errors.New("payment method has already expired")
	}

	return nil
}

// validateSecurityInfo validates security information
func (pm *PaymentMethod) validateSecurityInfo() error {
	security := pm.SecurityInfo

	// Risk score validation
	if security.RiskScore < 0 || security.RiskScore > 100 {
		return errors.New("risk_score must be between 0 and 100")
	}

	// Security level validation
	if security.SecurityLevel != "" {
		validLevels := []string{"low", "medium", "high"}
		valid := false
		for _, level := range validLevels {
			if security.SecurityLevel == level {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid security_level: %s", security.SecurityLevel)
		}
	}

	return nil
}

// validateProviderCompatibility validates provider and type compatibility
func (pm *PaymentMethod) validateProviderCompatibility() error {
	// Provider-Type compatibility matrix
	compatibilityMatrix := map[PaymentProvider][]PaymentMethodType{
		PaymentProviderVisa:        {PaymentMethodTypeCreditCard, PaymentMethodTypeDebitCard},
		PaymentProviderMastercard:  {PaymentMethodTypeCreditCard, PaymentMethodTypeDebitCard},
		PaymentProviderAmex:        {PaymentMethodTypeCreditCard},
		PaymentProviderPromptPay:   {PaymentMethodTypeMobileWallet, PaymentMethodTypeDigitalPayment},
		PaymentProviderGrabPay:     {PaymentMethodTypeEWallet, PaymentMethodTypeMobileWallet},
		PaymentProviderGoPay:       {PaymentMethodTypeEWallet, PaymentMethodTypeMobileWallet},
		PaymentProviderPayPal:      {PaymentMethodTypeDigitalPayment, PaymentMethodTypeEWallet},
		PaymentProviderBinance:     {PaymentMethodTypeCrypto},
		PaymentProviderCoinbase:    {PaymentMethodTypeCrypto},
		PaymentProviderMetaMask:    {PaymentMethodTypeCrypto},
	}

	if compatibleTypes, exists := compatibilityMatrix[pm.Provider]; exists {
		for _, compatibleType := range compatibleTypes {
			if pm.Type == compatibleType {
				return nil // Found compatible combination
			}
		}
		return fmt.Errorf("payment provider %s is not compatible with type %s", pm.Provider, pm.Type)
	}

	// If not in matrix, assume compatibility (for flexibility)
	return nil
}

// BeforeCreate sets up the payment method before database creation
func (pm *PaymentMethod) BeforeCreate() error {
	// Generate UUID if not set
	if pm.ID == uuid.Nil {
		pm.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now().UTC()
	pm.CreatedAt = now
	pm.UpdatedAt = now

	// Set default values
	if pm.Status == "" {
		pm.Status = PaymentMethodStatusPending
	}

	// Initialize security info
	if pm.SecurityInfo.SecurityLevel == "" {
		pm.SecurityInfo = SecurityInfo{
			RiskScore:     50.0, // Medium risk by default
			SecurityLevel: "medium",
		}
	}

	// Initialize usage stats
	if pm.UsageStats == (UsageStats{}) {
		pm.UsageStats = UsageStats{
			PreferredCurrency: pm.Currency,
		}
	}

	// Set display name if empty
	if strings.TrimSpace(pm.DisplayName) == "" {
		pm.DisplayName = pm.generateDisplayName()
	}

	// Set expiry for expirable types
	if pm.Type.IsExpirable() && pm.ExpiresAt == nil {
		if pm.ExpiryMonth != nil && pm.ExpiryYear != nil {
			// Set to last day of expiry month
			expiryDate := time.Date(*pm.ExpiryYear, time.Month(*pm.ExpiryMonth+1), 0, 23, 59, 59, 0, time.UTC)
			pm.ExpiresAt = &expiryDate
		}
	}

	// Validate before creation
	return pm.Validate()
}

// BeforeUpdate sets up the payment method before database update
func (pm *PaymentMethod) BeforeUpdate() error {
	// Update timestamp
	pm.UpdatedAt = time.Now().UTC()

	// Check expiry and update status
	if pm.Type.IsExpirable() && pm.ExpiresAt != nil {
		if time.Now().UTC().After(*pm.ExpiresAt) && pm.Status != PaymentMethodStatusExpired {
			pm.Status = PaymentMethodStatusExpired
		}
	}

	// Update verification timestamp if newly verified
	if pm.IsVerified && pm.VerifiedAt == nil {
		now := time.Now().UTC()
		pm.VerifiedAt = &now
	}

	// Validate before update
	return pm.Validate()
}

// generateDisplayName generates a display name based on type and provider
func (pm *PaymentMethod) generateDisplayName() string {
	switch pm.Type {
	case PaymentMethodTypeCreditCard, PaymentMethodTypeDebitCard:
		if pm.LastFourDigits != nil {
			return fmt.Sprintf("%s •••• %s", pm.Provider, *pm.LastFourDigits)
		}
		return string(pm.Provider)
	case PaymentMethodTypeBankAccount:
		if pm.Metadata.BankName != "" {
			return fmt.Sprintf("%s Bank", pm.Metadata.BankName)
		}
		return "Bank Account"
	case PaymentMethodTypeEWallet, PaymentMethodTypeMobileWallet:
		return string(pm.Provider)
	case PaymentMethodTypeCrypto:
		if pm.Metadata.TokenSymbol != "" {
			return fmt.Sprintf("%s Wallet", pm.Metadata.TokenSymbol)
		}
		return "Crypto Wallet"
	default:
		return string(pm.Type)
	}
}

// IsExpired checks if the payment method has expired
func (pm *PaymentMethod) IsExpired() bool {
	if !pm.Type.IsExpirable() {
		return false
	}
	return pm.ExpiresAt != nil && time.Now().UTC().After(*pm.ExpiresAt)
}

// CanProcessPayment checks if this payment method can process payments
func (pm *PaymentMethod) CanProcessPayment() bool {
	return pm.Status.CanTransact() && pm.IsVerified && !pm.IsExpired()
}

// MarkAsUsed updates usage statistics when payment method is used
func (pm *PaymentMethod) MarkAsUsed(amount int64, currency Currency, success bool) error {
	now := time.Now().UTC()
	pm.LastUsedAt = &now
	pm.UpdatedAt = now

	// Update usage stats
	pm.UsageStats.TotalTransactions++
	if success {
		pm.UsageStats.SuccessfulTransactions++
		pm.UsageStats.TotalAmountProcessed += amount
		pm.UsageStats.LastTransactionDate = now

		// Update average transaction size
		if pm.UsageStats.SuccessfulTransactions > 0 {
			pm.UsageStats.AverageTransactionSize = pm.UsageStats.TotalAmountProcessed / int64(pm.UsageStats.SuccessfulTransactions)
		}
	} else {
		pm.UsageStats.FailedTransactions++
	}

	// Update monthly usage if same month
	currentMonth := now.Month()
	if pm.UsageStats.LastTransactionDate.Month() == currentMonth {
		pm.UsageStats.MonthlyUsageCount++
	} else {
		pm.UsageStats.MonthlyUsageCount = 1
	}

	return pm.Validate()
}

// SetAsDefault marks this payment method as the default for the user
func (pm *PaymentMethod) SetAsDefault() error {
	pm.IsDefault = true
	pm.UpdatedAt = time.Now().UTC()
	return nil
}

// Verify marks the payment method as verified
func (pm *PaymentMethod) Verify() error {
	if pm.Status != PaymentMethodStatusPending {
		return fmt.Errorf("can only verify pending payment methods, current status: %s", pm.Status)
	}

	now := time.Now().UTC()
	pm.IsVerified = true
	pm.VerifiedAt = &now
	pm.Status = PaymentMethodStatusActive
	pm.UpdatedAt = now

	return pm.Validate()
}

// Block blocks the payment method
func (pm *PaymentMethod) Block(reason string) error {
	pm.Status = PaymentMethodStatusBlocked
	pm.UpdatedAt = time.Now().UTC()

	// Add to fraud flags
	if reason != "" {
		pm.SecurityInfo.FraudFlags = append(pm.SecurityInfo.FraudFlags, reason)
	}

	return pm.Validate()
}

// Activate activates a blocked or inactive payment method
func (pm *PaymentMethod) Activate() error {
	if pm.Status == PaymentMethodStatusExpired {
		return errors.New("cannot activate expired payment method")
	}
	if pm.Status == PaymentMethodStatusFailed {
		return errors.New("cannot activate failed payment method")
	}

	pm.Status = PaymentMethodStatusActive
	pm.UpdatedAt = time.Now().UTC()

	return pm.Validate()
}

// Deactivate deactivates the payment method
func (pm *PaymentMethod) Deactivate() error {
	pm.Status = PaymentMethodStatusInactive
	pm.IsDefault = false // Remove default status
	pm.UpdatedAt = time.Now().UTC()

	return pm.Validate()
}

// UpdateSecurityInfo updates security-related information
func (pm *PaymentMethod) UpdateSecurityInfo(cvvVerified, addressVerified bool, riskScore float64) error {
	pm.SecurityInfo.CVVVerified = cvvVerified
	pm.SecurityInfo.AddressVerified = addressVerified
	pm.SecurityInfo.RiskScore = riskScore
	pm.SecurityInfo.LastSecurityCheck = time.Now().UTC()

	// Determine security level based on risk score
	if riskScore <= 30 {
		pm.SecurityInfo.SecurityLevel = "low"
	} else if riskScore <= 70 {
		pm.SecurityInfo.SecurityLevel = "medium"
	} else {
		pm.SecurityInfo.SecurityLevel = "high"
	}

	pm.UpdatedAt = time.Now().UTC()
	return pm.Validate()
}

// GetSuccessRate returns the success rate percentage
func (pm *PaymentMethod) GetSuccessRate() float64 {
	if pm.UsageStats.TotalTransactions == 0 {
		return 0.0
	}
	return float64(pm.UsageStats.SuccessfulTransactions) / float64(pm.UsageStats.TotalTransactions) * 100
}

// GetFormattedTotalProcessed returns formatted total amount processed
func (pm *PaymentMethod) GetFormattedTotalProcessed() string {
	return pm.Currency.FormatAmount(pm.UsageStats.TotalAmountProcessed)
}

// GetFormattedAverageTransaction returns formatted average transaction size
func (pm *PaymentMethod) GetFormattedAverageTransaction() string {
	return pm.Currency.FormatAmount(pm.UsageStats.AverageTransactionSize)
}

// ToPublicPaymentMethod returns a sanitized version for public API responses
func (pm *PaymentMethod) ToPublicPaymentMethod() map[string]interface{} {
	result := map[string]interface{}{
		"id":              pm.ID,
		"type":            pm.Type,
		"provider":        pm.Provider,
		"status":          pm.Status,
		"is_default":      pm.IsDefault,
		"is_verified":     pm.IsVerified,
		"display_name":    pm.DisplayName,
		"last_four_digits": pm.LastFourDigits,
		"brand_name":      pm.BrandName,
		"country":         pm.Country,
		"currency":        pm.Currency,
		"last_used_at":    pm.LastUsedAt,
		"expires_at":      pm.ExpiresAt,
		"created_at":      pm.CreatedAt,
	}

	// Add expiry information for expirable types
	if pm.Type.IsExpirable() {
		result["expiry_month"] = pm.ExpiryMonth
		result["expiry_year"] = pm.ExpiryYear
		result["is_expired"] = pm.IsExpired()
	}

	// Add usage statistics
	result["usage_stats"] = map[string]interface{}{
		"total_transactions":    pm.UsageStats.TotalTransactions,
		"success_rate":          pm.GetSuccessRate(),
		"monthly_usage_count":   pm.UsageStats.MonthlyUsageCount,
		"formatted_total_processed": pm.GetFormattedTotalProcessed(),
	}

	return result
}

// ToSecurePaymentMethod returns a version with security information for admin purposes
func (pm *PaymentMethod) ToSecurePaymentMethod() map[string]interface{} {
	result := pm.ToPublicPaymentMethod()

	// Add security information
	result["security_info"] = pm.SecurityInfo
	result["provider_data"] = pm.ProviderData
	result["external_id"] = pm.ExternalID

	return result
}

// PaymentMethodCreateRequest represents a request to create a new payment method
type PaymentMethodCreateRequest struct {
	Type         PaymentMethodType     `json:"type" validate:"required"`
	Provider     PaymentProvider       `json:"provider" validate:"required"`
	Currency     Currency              `json:"currency" validate:"required"`
	Country      string                `json:"country" validate:"required,len=2"`
	DisplayName  string                `json:"display_name,omitempty"`
	IsDefault    bool                  `json:"is_default"`
	Metadata     PaymentMethodMetadata `json:"metadata,omitempty"`
	ProviderData ProviderData          `json:"provider_data,omitempty"`
	ExternalID   *string               `json:"external_id,omitempty"`
}

// ToPaymentMethod converts a create request to a PaymentMethod model
func (req *PaymentMethodCreateRequest) ToPaymentMethod(userID uuid.UUID, walletID *uuid.UUID) *PaymentMethod {
	return &PaymentMethod{
		UserID:       userID,
		WalletID:     walletID,
		Type:         req.Type,
		Provider:     req.Provider,
		Currency:     req.Currency,
		Country:      req.Country,
		DisplayName:  req.DisplayName,
		IsDefault:    req.IsDefault,
		Metadata:     req.Metadata,
		ProviderData: req.ProviderData,
		ExternalID:   req.ExternalID,
	}
}

// PaymentMethodManager provides payment method management utilities
type PaymentMethodManager struct {
	// Add dependencies like database, payment processors, etc.
}

// NewPaymentMethodManager creates a new payment method manager
func NewPaymentMethodManager() *PaymentMethodManager {
	return &PaymentMethodManager{}
}

// CreatePaymentMethod creates a new payment method with proper validation
func (pmm *PaymentMethodManager) CreatePaymentMethod(req *PaymentMethodCreateRequest, userID uuid.UUID, walletID *uuid.UUID) (*PaymentMethod, error) {
	paymentMethod := req.ToPaymentMethod(userID, walletID)

	if err := paymentMethod.BeforeCreate(); err != nil {
		return nil, fmt.Errorf("payment method creation failed: %v", err)
	}

	return paymentMethod, nil
}