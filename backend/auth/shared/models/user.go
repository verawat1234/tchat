package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/shared/utils"
)

// UserStatus represents the status of a user account
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

// KYCTier represents the KYC verification level
type KYCTier int

const (
	KYCTierUnverified KYCTier = 0
	KYCTierBasic     KYCTier = 1
	KYCTierStandard  KYCTier = 2
	KYCTierPremium   KYCTier = 3
)

// Country represents supported countries
type Country string

const (
	CountryThailand   Country = "TH"
	CountrySingapore  Country = "SG"
	CountryIndonesia  Country = "ID"
	CountryMalaysia   Country = "MY"
	CountryPhilippines Country = "PH"
	CountryVietnam    Country = "VN"
)

// VerificationTier represents verification levels
type VerificationTier int

const (
	VerificationTierNone     VerificationTier = 0
	VerificationTierPhone    VerificationTier = 1
	VerificationTierEmail    VerificationTier = 2
	VerificationTierKYC      VerificationTier = 3
	VerificationTierFull     VerificationTier = 4
)

// UserPreferences represents user preferences - DISABLED due to missing database columns
// Current database schema doesn't have pref_* columns, only profile_* columns
// This struct is commented out until the database schema is updated to match
/*
type UserPreferences struct {
	NotificationsEmail    bool   `json:"notifications_email" gorm:"column:pref_notifications_email;default:true"`
	NotificationsPush     bool   `json:"notifications_push" gorm:"column:pref_notifications_push;default:true"`
	Language              string `json:"language" gorm:"column:pref_language;size:10"`
	Theme                 string `json:"theme" gorm:"column:pref_theme;size:20"`
	PrivacyLevel          string `json:"privacy_level" gorm:"column:pref_privacy_level;size:20"`
}
*/

// IsValid checks if the user status is valid
func (s UserStatus) IsValid() bool {
	switch s {
	case UserStatusActive, UserStatusSuspended, UserStatusDeleted:
		return true
	default:
		return false
	}
}

// UserProfile represents the user's profile information - DISABLED due to column naming conflicts
// Current database schema has profile_* columns, this struct expects different column names
// This struct is commented out until database schema is aligned
/*
type UserProfile struct {
	DisplayName string `json:"display_name" gorm:"column:display_name;size:100"`
	AvatarURL   string `json:"avatar_url,omitempty" gorm:"column:avatar;size:500"`
	Locale      string `json:"locale" gorm:"column:locale;size:10;default:'en'"`
	Timezone    string `json:"timezone" gorm:"column:timezone;size:50;default:'UTC'"`
}
*/

// User represents a platform user with Southeast Asian regional compliance
type User struct {
	ID          uuid.UUID   `json:"id" gorm:"column:id;type:uuid;primary_key;default:gen_random_uuid()"`
	Username    string      `json:"username,omitempty" gorm:"column:username;size:100;uniqueIndex"`
	Phone       string      `json:"phone,omitempty" gorm:"column:phone;size:20"`
	PhoneNumber string      `json:"phone_number,omitempty" gorm:"column:phone_number;size:20"`
	Email       string      `json:"email,omitempty" gorm:"column:email;size:255;uniqueIndex"`
	Name        string      `json:"name" gorm:"column:name;size:100;not null"`
	DisplayName string      `json:"display_name,omitempty" gorm:"column:display_name;size:100"`
	FirstName   string      `json:"first_name,omitempty" gorm:"column:first_name;size:100"`
	LastName    string      `json:"last_name,omitempty" gorm:"column:last_name;size:100"`
	Avatar      string      `json:"avatar,omitempty" gorm:"column:avatar;size:500"`
	Country     string      `json:"country" gorm:"column:country;size:5;not null;default:'TH'"`
	CountryCode string      `json:"country_code,omitempty" gorm:"column:country_code;size:5"`
	Locale      string      `json:"locale" gorm:"column:locale;size:10;not null;default:'th-TH'"`
	Language    string      `json:"language,omitempty" gorm:"column:language;size:10"`
	Timezone    string      `json:"timezone,omitempty" gorm:"column:timezone;size:50"`
	TimezoneAlias string    `json:"timezone_alias,omitempty" gorm:"column:timezone_alias;size:50"`
	Bio         string      `json:"bio,omitempty" gorm:"column:bio"`
	DateOfBirth *time.Time  `json:"date_of_birth,omitempty" gorm:"column:date_of_birth"`
	Gender      string      `json:"gender,omitempty" gorm:"column:gender;size:20"`

	// Status and verification
	Active       bool       `json:"is_active" gorm:"column:is_active;default:true"`
	KYCStatus    string     `json:"kyc_status" gorm:"column:kyc_status;size:20;default:'none'"`
	KYCTier      int        `json:"kyc_tier" gorm:"column:kyc_tier;default:0"`
	Status       string     `json:"status" gorm:"column:status;size:20;default:'offline'"`
	LastSeen     *time.Time `json:"last_seen,omitempty" gorm:"column:last_seen"`
	LastActiveAt *time.Time `json:"last_active_at,omitempty" gorm:"column:last_active_at"`

	// Verification flags
	Verified        bool `json:"is_verified" gorm:"column:is_verified;default:false"`
	EmailVerified   bool `json:"is_email_verified" gorm:"column:is_email_verified;default:false"`
	PhoneVerified   bool `json:"is_phone_verified" gorm:"column:is_phone_verified;default:false"`

	// User preferences - matching actual database schema
	PrefTheme             string `json:"pref_theme,omitempty" gorm:"column:pref_theme;size:20"`
	PrefLanguage          string `json:"pref_language,omitempty" gorm:"column:pref_language;size:10"`
	PrefNotificationsEmail bool  `json:"pref_notifications_email" gorm:"column:pref_notifications_email;default:true"`
	PrefNotificationsPush  bool  `json:"pref_notifications_push" gorm:"column:pref_notifications_push;default:true"`
	PrefPrivacyLevel      string `json:"pref_privacy_level,omitempty" gorm:"column:pref_privacy_level;size:20"`

	// Metadata - matching actual database schema
	Metadata json.RawMessage `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb;default:'{}'"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	// DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"` // DISABLED: column doesn't exist
}

// TableName returns the table name for the User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate sets up the user before creation
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	// Set default display name
	if u.DisplayName == "" {
		u.DisplayName = "User"
	}

	// Set default locale and timezone based on country
	if u.Locale == "" {
		u.Locale = GetDefaultLocaleForCountry(u.Country)
	}

	if u.Timezone == "" {
		u.Timezone = GetDefaultTimezoneForCountry(u.Country)
	}

	// Validate the user
	if err := u.Validate(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate validates the user before updating
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return u.Validate()
}

// Validate validates the user data for Southeast Asian compliance
func (u *User) Validate() error {
	// Validate country code
	if !IsValidSEACountry(u.Country) {
		return fmt.Errorf("invalid country code: %s", u.Country)
	}

	// Basic name validation
	if len(u.Name) == 0 || len(u.Name) > 100 {
		return fmt.Errorf("name must be between 1 and 100 characters")
	}

	// Validate display name length
	if u.DisplayName != "" && len(u.DisplayName) > 100 {
		return fmt.Errorf("display name cannot exceed 100 characters")
	}

	// Validate locale for the country
	if u.Locale != "" && !IsValidLocaleForCountry(u.Locale, u.Country) {
		return fmt.Errorf("invalid locale %s for country %s", u.Locale, u.Country)
	}

	// Validate timezone
	if u.Timezone != "" && !IsValidTimezone(u.Timezone) {
		return fmt.Errorf("invalid timezone: %s", u.Timezone)
	}

	return nil
}

// IsValidPhoneNumber validates phone number format for the user's country
func (u *User) IsValidPhoneNumber() bool {
	phoneToCheck := u.PhoneNumber
	if phoneToCheck == "" {
		phoneToCheck = u.Phone
	}
	countryCode := u.Country
	if u.CountryCode != "" {
		countryCode = u.CountryCode
	}
	return utils.IsValidPhoneNumber(phoneToCheck, countryCode)
}

// GetFullPhoneNumber returns the phone number in E.164 format
func (u *User) GetFullPhoneNumber() string {
	countryPrefixes := map[string]string{
		"TH": "+66",
		"SG": "+65",
		"ID": "+62",
		"MY": "+60",
		"PH": "+63",
		"VN": "+84",
	}

	phoneToFormat := u.PhoneNumber
	if phoneToFormat == "" {
		phoneToFormat = u.Phone
	}

	countryCode := u.Country
	if u.CountryCode != "" {
		countryCode = u.CountryCode
	}

	if prefix, exists := countryPrefixes[countryCode]; exists {
		// Remove leading zeros from local number
		localNumber := strings.TrimLeft(phoneToFormat, "0")
		return prefix + localNumber
	}

	return phoneToFormat
}

// IsActive checks if the user account is active
func (u *User) IsActive() bool {
	return u.Active && (u.Status == "active" || u.Status == "online")
}

// UpdateLastLogin updates the last login information
func (u *User) UpdateLastLogin(ipAddress string) {
	now := time.Now()
	u.LastActiveAt = &now
}

// GetAge returns the age of the user account
func (u *User) GetAge() time.Duration {
	return time.Since(u.CreatedAt)
}

// GetRegionalSettings returns regional settings for the user
func (u *User) GetRegionalSettings() map[string]interface{} {
	countryCode := u.Country
	if u.CountryCode != "" {
		countryCode = u.CountryCode
	}

	return map[string]interface{}{
		"country_code":  countryCode,
		"locale":        u.Locale,
		"timezone":      u.Timezone,
		"phone_format":  u.GetFullPhoneNumber(),
		"currency":      GetDefaultCurrencyForCountry(countryCode),
	}
}

// CanUpdateProfile checks if user can update their profile
func (u *User) CanUpdateProfile() bool {
	return u.Active
}

// CanSubmitKYC checks if user can submit KYC documents
func (u *User) CanSubmitKYC() bool {
	return u.Active && u.PhoneVerified
}

// CanUpgradeKYC checks if user can upgrade their KYC tier
func (u *User) CanUpgradeKYC() bool {
	return u.Active && u.PhoneVerified
}

// IsKYCVerified checks if user has completed KYC verification
func (u *User) IsKYCVerified() bool {
	return u.KYCTier >= 1 // KYCTierBasic
}

// GetMaxTransactionLimit returns transaction limit based on KYC tier
func (u *User) GetMaxTransactionLimit() float64 {
	switch u.KYCTier {
	case 0: // KYCTierUnverified
		return 100.0
	case 1: // KYCTierBasic
		return 1000.0
	case 2: // KYCTierStandard
		return 10000.0
	case 3: // KYCTierPremium
		return 100000.0
	default:
		return 0.0
	}
}

// UpdateKYCTier updates the user's KYC tier
func (u *User) UpdateKYCTier(newTier int) error {
	if !u.CanUpgradeKYC() {
		return fmt.Errorf("user cannot upgrade KYC tier")
	}

	if newTier < u.KYCTier {
		return fmt.Errorf("cannot downgrade KYC tier")
	}

	u.KYCTier = newTier
	return nil
}

// MarshalJSON customizes JSON serialization
func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		*Alias
		PhoneNumber string `json:"phone_number"`
	}{
		Alias:       (*Alias)(u),
		PhoneNumber: u.GetFullPhoneNumber(),
	})
}

// Southeast Asian regional helper functions

// IsValidSEACountry checks if the country code is a supported Southeast Asian country
func IsValidSEACountry(countryCode string) bool {
	supportedCountries := map[string]bool{
		"TH": true, // Thailand
		"SG": true, // Singapore
		"ID": true, // Indonesia
		"MY": true, // Malaysia
		"PH": true, // Philippines
		"VN": true, // Vietnam
	}
	return supportedCountries[countryCode]
}

// GetDataRegionForCountry returns the data region for a country
func GetDataRegionForCountry(countryCode string) string {
	regionMapping := map[string]string{
		"TH": "sea-central",  // Thailand
		"SG": "sea-central",  // Singapore
		"ID": "sea-central",  // Indonesia
		"MY": "sea-central",  // Malaysia
		"PH": "sea-east",     // Philippines
		"VN": "sea-north",    // Vietnam
	}

	if region, exists := regionMapping[countryCode]; exists {
		return region
	}
	return "sea-central" // Default region
}

// GetDefaultLocaleForCountry returns the default locale for a country
func GetDefaultLocaleForCountry(countryCode string) string {
	localeMapping := map[string]string{
		"TH": "th", // Thai
		"SG": "en", // English
		"ID": "id", // Indonesian
		"MY": "ms", // Malay
		"PH": "fil", // Filipino
		"VN": "vi", // Vietnamese
	}

	if locale, exists := localeMapping[countryCode]; exists {
		return locale
	}
	return "en" // Default to English
}

// GetDefaultTimezoneForCountry returns the default timezone for a country
func GetDefaultTimezoneForCountry(countryCode string) string {
	timezoneMapping := map[string]string{
		"TH": "Asia/Bangkok",
		"SG": "Asia/Singapore",
		"ID": "Asia/Jakarta",
		"MY": "Asia/Kuala_Lumpur",
		"PH": "Asia/Manila",
		"VN": "Asia/Ho_Chi_Minh",
	}

	if timezone, exists := timezoneMapping[countryCode]; exists {
		return timezone
	}
	return "UTC" // Default timezone
}

// GetDefaultCurrencyForCountry returns the default currency for a country
func GetDefaultCurrencyForCountry(countryCode string) string {
	currencyMapping := map[string]string{
		"TH": "THB", // Thai Baht
		"SG": "SGD", // Singapore Dollar
		"ID": "IDR", // Indonesian Rupiah
		"MY": "MYR", // Malaysian Ringgit
		"PH": "PHP", // Philippine Peso
		"VN": "VND", // Vietnamese Dong
	}

	if currency, exists := currencyMapping[countryCode]; exists {
		return currency
	}
	return "USD" // Default currency
}

// IsValidLocaleForCountry checks if a locale is valid for a country
func IsValidLocaleForCountry(locale, countryCode string) bool {
	validLocales := map[string][]string{
		"TH": {"th", "en"},
		"SG": {"en"},
		"ID": {"id", "en"},
		"MY": {"ms", "en"},
		"PH": {"fil", "en"},
		"VN": {"vi", "en"},
	}

	if locales, exists := validLocales[countryCode]; exists {
		for _, validLocale := range locales {
			if locale == validLocale {
				return true
			}
		}
	}
	return false
}

// IsValidTimezone checks if a timezone is valid
func IsValidTimezone(timezone string) bool {
	// Load timezone to validate
	_, err := time.LoadLocation(timezone)
	return err == nil
}