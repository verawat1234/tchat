package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BusinessVerificationStatus represents the verification status of a business
type BusinessVerificationStatus string

const (
	BusinessVerificationPending  BusinessVerificationStatus = "pending"
	BusinessVerificationVerified BusinessVerificationStatus = "verified"
	BusinessVerificationRejected BusinessVerificationStatus = "rejected"
	BusinessVerificationSuspended BusinessVerificationStatus = "suspended"
)

// IsValid checks if the verification status is valid
func (s BusinessVerificationStatus) IsValid() bool {
	switch s {
	case BusinessVerificationPending, BusinessVerificationVerified, BusinessVerificationRejected, BusinessVerificationSuspended:
		return true
	default:
		return false
	}
}

// BusinessContactInfo represents contact information for a business
type BusinessContactInfo struct {
	Phone   string `json:"phone" gorm:"column:contact_phone;size:20"`
	Email   string `json:"email" gorm:"column:contact_email;size:255"`
	Website string `json:"website,omitempty" gorm:"column:contact_website;size:255"`
}

// BusinessAddress represents the physical address of a business
type BusinessAddress struct {
	Street     string `json:"street" gorm:"column:address_street;size:255;not null"`
	City       string `json:"city" gorm:"column:address_city;size:100;not null"`
	State      string `json:"state,omitempty" gorm:"column:address_state;size:100"`
	PostalCode string `json:"postal_code" gorm:"column:address_postal_code;size:20"`
	Country    string `json:"country" gorm:"column:address_country;size:2;not null"`
}

// BusinessSettings represents business operational settings
type BusinessSettings struct {
	SupportedCurrencies []string `json:"supported_currencies" gorm:"column:settings_currencies;type:jsonb"`
	SupportedLanguages  []string `json:"supported_languages" gorm:"column:settings_languages;type:jsonb"`
	ShippingCountries   []string `json:"shipping_countries" gorm:"column:settings_shipping;type:jsonb"`
	TaxSettings         map[string]interface{} `json:"tax_settings" gorm:"column:settings_tax;type:jsonb"`
	BusinessHours       map[string]interface{} `json:"business_hours,omitempty" gorm:"column:settings_hours;type:jsonb"`
	PaymentMethods      []string `json:"payment_methods,omitempty" gorm:"column:settings_payments;type:jsonb"`
}

// BusinessComplianceData represents regulatory compliance information
type BusinessComplianceData struct {
	TaxID              string                 `json:"tax_id,omitempty" gorm:"column:compliance_tax_id;size:50"`
	BusinessLicense    string                 `json:"business_license,omitempty" gorm:"column:compliance_license;size:100"`
	VATNumber          string                 `json:"vat_number,omitempty" gorm:"column:compliance_vat;size:50"`
	RegulationCategory string                 `json:"regulation_category,omitempty" gorm:"column:compliance_category;size:100"`
	ComplianceData     map[string]interface{} `json:"compliance_data,omitempty" gorm:"column:compliance_data;type:jsonb"`
}

// Business represents a commercial entity on the platform
type Business struct {
	ID          uuid.UUID                  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OwnerID     uuid.UUID                  `json:"owner_id" gorm:"type:uuid;not null;index"`
	Name        string                     `json:"name" gorm:"column:name;size:100;not null"`
	Description string                     `json:"description" gorm:"column:description;size:1000"`
	Category    string                     `json:"category" gorm:"column:category;size:50;not null"`

	// Verification
	VerificationStatus BusinessVerificationStatus `json:"verification_status" gorm:"column:verification_status;type:varchar(20);not null;default:'pending'"`
	VerifiedAt         *time.Time                 `json:"verified_at,omitempty" gorm:"column:verified_at"`
	VerificationNotes  string                     `json:"verification_notes,omitempty" gorm:"column:verification_notes;size:500"`

	// Contact and location
	ContactInfo BusinessContactInfo `json:"contact_info" gorm:"embedded;embeddedPrefix:contact_"`
	Address     BusinessAddress     `json:"address" gorm:"embedded;embeddedPrefix:address_"`

	// Business settings and compliance
	BusinessSettings   BusinessSettings       `json:"business_settings" gorm:"embedded;embeddedPrefix:settings_"`
	ComplianceData     BusinessComplianceData `json:"compliance_data" gorm:"embedded;embeddedPrefix:compliance_"`

	// Platform metrics
	TotalProducts    int     `json:"total_products" gorm:"column:total_products;default:0"`
	TotalOrders      int     `json:"total_orders" gorm:"column:total_orders;default:0"`
	AverageRating    float64 `json:"average_rating" gorm:"column:average_rating;default:0"`
	TotalReviews     int     `json:"total_reviews" gorm:"column:total_reviews;default:0"`
	IsActive         bool    `json:"is_active" gorm:"column:is_active;default:true"`

	// Regional compliance
	DataRegion           string `json:"data_region,omitempty" gorm:"column:data_region;size:20"`
	RequiresKYB          bool   `json:"requires_kyb" gorm:"column:requires_kyb;default:false"`
	KYBStatus            string `json:"kyb_status,omitempty" gorm:"column:kyb_status;size:20"`
	KYBCompletedAt       *time.Time `json:"kyb_completed_at,omitempty" gorm:"column:kyb_completed_at"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Relationships
	Owner *User `json:"owner,omitempty" gorm:"foreignKey:OwnerID;references:ID"`
}

// TableName returns the table name for the Business model
func (Business) TableName() string {
	return "businesses"
}

// BeforeCreate sets up the business before creation
func (b *Business) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}

	// Set regional defaults
	if b.DataRegion == "" {
		b.DataRegion = GetDataRegionForCountry(b.Address.Country)
	}

	// Set default currencies and languages based on country
	if len(b.BusinessSettings.SupportedCurrencies) == 0 {
		b.BusinessSettings.SupportedCurrencies = []string{GetDefaultCurrencyForCountry(b.Address.Country)}
	}

	if len(b.BusinessSettings.SupportedLanguages) == 0 {
		b.BusinessSettings.SupportedLanguages = []string{GetDefaultLocaleForCountry(b.Address.Country), "en"}
	}

	// Set shipping countries to include business country by default
	if len(b.BusinessSettings.ShippingCountries) == 0 {
		b.BusinessSettings.ShippingCountries = []string{b.Address.Country}
	}

	// Check if KYB is required for this country and category
	b.RequiresKYB = b.IsKYBRequired()
	if b.RequiresKYB {
		b.KYBStatus = "pending"
	}

	// Validate the business
	if err := b.Validate(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate validates the business before updating
func (b *Business) BeforeUpdate(tx *gorm.DB) error {
	return b.Validate()
}

// Validate validates the business data for Southeast Asian compliance
func (b *Business) Validate() error {
	// Validate country code
	if !IsValidSEACountry(b.Address.Country) {
		return fmt.Errorf("invalid country code: %s", b.Address.Country)
	}

	// Validate verification status
	if !b.VerificationStatus.IsValid() {
		return fmt.Errorf("invalid verification status: %s", b.VerificationStatus)
	}

	// Validate business name
	if len(b.Name) == 0 || len(b.Name) > 100 {
		return fmt.Errorf("business name must be between 1 and 100 characters")
	}

	// Validate category
	if !b.IsValidCategory() {
		return fmt.Errorf("invalid business category: %s", b.Category)
	}

	// Validate contact info
	if err := b.ValidateContactInfo(); err != nil {
		return fmt.Errorf("contact info validation failed: %w", err)
	}

	// Validate address
	if err := b.ValidateAddress(); err != nil {
		return fmt.Errorf("address validation failed: %w", err)
	}

	// Validate currencies are supported for the country
	if err := b.ValidateCurrencies(); err != nil {
		return fmt.Errorf("currency validation failed: %w", err)
	}

	return nil
}

// IsValidCategory checks if the business category is valid
func (b *Business) IsValidCategory() bool {
	validCategories := map[string]bool{
		"electronics":   true,
		"fashion":       true,
		"food":          true,
		"health":        true,
		"beauty":        true,
		"home":          true,
		"sports":        true,
		"automotive":    true,
		"books":         true,
		"toys":          true,
		"services":      true,
		"digital":       true,
		"agriculture":   true,
		"crafts":        true,
		"jewelry":       true,
		"travel":        true,
		"education":     true,
		"finance":       true,
		"real_estate":   true,
		"entertainment": true,
	}
	return validCategories[b.Category]
}

// ValidateContactInfo validates the contact information
func (b *Business) ValidateContactInfo() error {
	// Validate email format
	if b.ContactInfo.Email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(b.ContactInfo.Email) {
			return fmt.Errorf("invalid email format: %s", b.ContactInfo.Email)
		}
	}

	// Validate phone number format (basic validation)
	if b.ContactInfo.Phone != "" {
		phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{6,14}$`)
		if !phoneRegex.MatchString(b.ContactInfo.Phone) {
			return fmt.Errorf("invalid phone number format: %s", b.ContactInfo.Phone)
		}
	}

	// Validate website URL format
	if b.ContactInfo.Website != "" {
		urlRegex := regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
		if !urlRegex.MatchString(b.ContactInfo.Website) {
			return fmt.Errorf("invalid website URL format: %s", b.ContactInfo.Website)
		}
	}

	return nil
}

// ValidateAddress validates the business address
func (b *Business) ValidateAddress() error {
	// Validate required fields
	if b.Address.Street == "" {
		return fmt.Errorf("street address is required")
	}
	if b.Address.City == "" {
		return fmt.Errorf("city is required")
	}
	if b.Address.Country == "" {
		return fmt.Errorf("country is required")
	}

	// Validate postal code format for specific countries
	if err := b.ValidatePostalCode(); err != nil {
		return err
	}

	return nil
}

// ValidatePostalCode validates postal code format for Southeast Asian countries
func (b *Business) ValidatePostalCode() error {
	if b.Address.PostalCode == "" {
		return nil // Postal code is optional for some countries
	}

	switch b.Address.Country {
	case "TH": // Thailand: 5 digits
		if !regexp.MustCompile(`^\d{5}$`).MatchString(b.Address.PostalCode) {
			return fmt.Errorf("invalid postal code for Thailand: %s", b.Address.PostalCode)
		}
	case "SG": // Singapore: 6 digits
		if !regexp.MustCompile(`^\d{6}$`).MatchString(b.Address.PostalCode) {
			return fmt.Errorf("invalid postal code for Singapore: %s", b.Address.PostalCode)
		}
	case "MY": // Malaysia: 5 digits
		if !regexp.MustCompile(`^\d{5}$`).MatchString(b.Address.PostalCode) {
			return fmt.Errorf("invalid postal code for Malaysia: %s", b.Address.PostalCode)
		}
	case "PH": // Philippines: 4 digits
		if !regexp.MustCompile(`^\d{4}$`).MatchString(b.Address.PostalCode) {
			return fmt.Errorf("invalid postal code for Philippines: %s", b.Address.PostalCode)
		}
	case "ID": // Indonesia: 5 digits
		if !regexp.MustCompile(`^\d{5}$`).MatchString(b.Address.PostalCode) {
			return fmt.Errorf("invalid postal code for Indonesia: %s", b.Address.PostalCode)
		}
	case "VN": // Vietnam: 6 digits
		if !regexp.MustCompile(`^\d{6}$`).MatchString(b.Address.PostalCode) {
			return fmt.Errorf("invalid postal code for Vietnam: %s", b.Address.PostalCode)
		}
	}

	return nil
}

// ValidateCurrencies validates that currencies are supported for the business country
func (b *Business) ValidateCurrencies() error {
	supportedCurrencies := GetSupportedCurrenciesForCountry(b.Address.Country)

	for _, currency := range b.BusinessSettings.SupportedCurrencies {
		found := false
		for _, supported := range supportedCurrencies {
			if currency == supported {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("currency %s not supported for country %s", currency, b.Address.Country)
		}
	}

	return nil
}

// IsKYBRequired checks if Know Your Business verification is required
func (b *Business) IsKYBRequired() bool {
	// KYB requirements based on country and category
	kybRequiredCategories := map[string][]string{
		"TH": {"finance", "real_estate", "jewelry", "automotive"},
		"SG": {"finance", "real_estate", "jewelry"},
		"ID": {"finance", "real_estate", "automotive"},
		"MY": {"finance", "real_estate", "jewelry"},
		"PH": {"finance", "real_estate"},
		"VN": {"finance", "real_estate", "automotive"},
	}

	if categories, exists := kybRequiredCategories[b.Address.Country]; exists {
		for _, category := range categories {
			if b.Category == category {
				return true
			}
		}
	}

	return false
}

// IsVerified checks if the business is verified
func (b *Business) IsVerified() bool {
	return b.VerificationStatus == BusinessVerificationVerified
}

// IsOperational checks if the business can operate (verified and active)
func (b *Business) IsOperational() bool {
	return b.IsVerified() && b.IsActive && (!b.RequiresKYB || b.KYBStatus == "completed")
}

// GetRegionalSettings returns regional settings for the business
func (b *Business) GetRegionalSettings() map[string]interface{} {
	return map[string]interface{}{
		"country":              b.Address.Country,
		"data_region":          b.DataRegion,
		"supported_currencies": b.BusinessSettings.SupportedCurrencies,
		"supported_languages":  b.BusinessSettings.SupportedLanguages,
		"shipping_countries":   b.BusinessSettings.ShippingCountries,
		"requires_kyb":         b.RequiresKYB,
		"kyb_status":           b.KYBStatus,
	}
}

// UpdateMetrics updates business metrics
func (b *Business) UpdateMetrics(products, orders, reviews int, rating float64) {
	b.TotalProducts = products
	b.TotalOrders = orders
	b.TotalReviews = reviews
	b.AverageRating = rating
}

// Activate activates the business
func (b *Business) Activate() error {
	if !b.IsVerified() {
		return fmt.Errorf("business must be verified before activation")
	}
	if b.RequiresKYB && b.KYBStatus != "completed" {
		return fmt.Errorf("KYB verification must be completed before activation")
	}
	b.IsActive = true
	return nil
}

// Deactivate deactivates the business
func (b *Business) Deactivate() {
	b.IsActive = false
}

// Helper functions for currency support

// GetSupportedCurrenciesForCountry returns supported currencies for a country
func GetSupportedCurrenciesForCountry(countryCode string) []string {
	currencyMapping := map[string][]string{
		"TH": {"THB", "USD"},
		"SG": {"SGD", "USD"},
		"ID": {"IDR", "USD"},
		"MY": {"MYR", "USD"},
		"PH": {"PHP", "USD"},
		"VN": {"VND", "USD"},
	}

	if currencies, exists := currencyMapping[countryCode]; exists {
		return currencies
	}
	return []string{"USD"} // Default
}

// MarshalJSON customizes JSON serialization
func (b *Business) MarshalJSON() ([]byte, error) {
	type Alias Business
	return json.Marshal(&struct {
		*Alias
		IsOperational   bool                   `json:"is_operational"`
		RegionalSettings map[string]interface{} `json:"regional_settings"`
	}{
		Alias:           (*Alias)(b),
		IsOperational:   b.IsOperational(),
		RegionalSettings: b.GetRegionalSettings(),
	})
}