package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserStatus represents the status of a user account
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

// IsValid checks if the user status is valid
func (s UserStatus) IsValid() bool {
	switch s {
	case UserStatusActive, UserStatusSuspended, UserStatusDeleted:
		return true
	default:
		return false
	}
}

// UserProfile represents the user's profile information
type UserProfile struct {
	DisplayName string `json:"display_name" gorm:"column:display_name;size:100;not null"`
	AvatarURL   string `json:"avatar_url,omitempty" gorm:"column:avatar_url;size:500"`
	Locale      string `json:"locale" gorm:"column:locale;size:5;not null;default:'en'"`
	Timezone    string `json:"timezone" gorm:"column:timezone;size:50;not null;default:'UTC'"`
}

// User represents a platform user with Southeast Asian regional compliance
type User struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PhoneNumber string      `json:"phone_number" gorm:"column:phone_number;size:20;not null;uniqueIndex"`
	CountryCode string      `json:"country_code" gorm:"column:country_code;size:2;not null"`
	Status      UserStatus  `json:"status" gorm:"column:status;type:varchar(20);not null;default:'active'"`
	Profile     UserProfile `json:"profile" gorm:"embedded;embeddedPrefix:profile_"`

	// Regional compliance fields
	DataRegion       string `json:"data_region,omitempty" gorm:"column:data_region;size:20"`
	ConsentDate      *time.Time `json:"consent_date,omitempty" gorm:"column:consent_date"`
	ConsentVersion   string `json:"consent_version,omitempty" gorm:"column:consent_version;size:10"`
	LegalBasis       string `json:"legal_basis,omitempty" gorm:"column:legal_basis;size:50"`
	ProcessingPurpose string `json:"processing_purpose,omitempty" gorm:"column:processing_purpose;size:100"`

	// Verification status
	PhoneVerified    bool       `json:"phone_verified" gorm:"column:phone_verified;default:false"`
	PhoneVerifiedAt  *time.Time `json:"phone_verified_at,omitempty" gorm:"column:phone_verified_at"`
	EmailVerified    bool       `json:"email_verified" gorm:"column:email_verified;default:false"`
	EmailVerifiedAt  *time.Time `json:"email_verified_at,omitempty" gorm:"column:email_verified_at"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Audit fields
	LastLoginAt    *time.Time `json:"last_login_at,omitempty" gorm:"column:last_login_at"`
	LastLoginIP    string     `json:"last_login_ip,omitempty" gorm:"column:last_login_ip;size:45"`
	FailedAttempts int        `json:"failed_attempts" gorm:"column:failed_attempts;default:0"`
	LockedUntil    *time.Time `json:"locked_until,omitempty" gorm:"column:locked_until"`
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

	// Set default profile values
	if u.Profile.DisplayName == "" {
		u.Profile.DisplayName = "User"
	}

	// Set regional defaults based on country
	if u.DataRegion == "" {
		u.DataRegion = GetDataRegionForCountry(u.CountryCode)
	}

	// Set default locale and timezone based on country
	if u.Profile.Locale == "" {
		u.Profile.Locale = GetDefaultLocaleForCountry(u.CountryCode)
	}

	if u.Profile.Timezone == "" {
		u.Profile.Timezone = GetDefaultTimezoneForCountry(u.CountryCode)
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
	if !IsValidSEACountry(u.CountryCode) {
		return fmt.Errorf("invalid country code: %s", u.CountryCode)
	}

	// Validate phone number format for the country
	if !u.IsValidPhoneNumber() {
		return fmt.Errorf("invalid phone number format for country %s: %s", u.CountryCode, u.PhoneNumber)
	}

	// Validate status
	if !u.Status.IsValid() {
		return fmt.Errorf("invalid user status: %s", u.Status)
	}

	// Validate locale for the country
	if !IsValidLocaleForCountry(u.Profile.Locale, u.CountryCode) {
		return fmt.Errorf("invalid locale %s for country %s", u.Profile.Locale, u.CountryCode)
	}

	// Validate timezone
	if !IsValidTimezone(u.Profile.Timezone) {
		return fmt.Errorf("invalid timezone: %s", u.Profile.Timezone)
	}

	// Validate display name length
	if len(u.Profile.DisplayName) == 0 || len(u.Profile.DisplayName) > 100 {
		return fmt.Errorf("display name must be between 1 and 100 characters")
	}

	return nil
}

// IsValidPhoneNumber validates phone number format for the user's country
func (u *User) IsValidPhoneNumber() bool {
	// Remove any non-digit characters for validation
	phoneDigits := regexp.MustCompile(`\D`).ReplaceAllString(u.PhoneNumber, "")

	// Country-specific phone number validation
	switch u.CountryCode {
	case "TH": // Thailand: 8-9 digits starting with 6-9
		return regexp.MustCompile(`^[6-9]\d{7,8}$`).MatchString(phoneDigits)
	case "SG": // Singapore: 8 digits starting with 6, 8, or 9
		return regexp.MustCompile(`^[689]\d{7}$`).MatchString(phoneDigits)
	case "ID": // Indonesia: 8-12 digits starting with 8
		return regexp.MustCompile(`^8\d{7,11}$`).MatchString(phoneDigits)
	case "MY": // Malaysia: 9-10 digits starting with 1
		return regexp.MustCompile(`^1\d{8,9}$`).MatchString(phoneDigits)
	case "PH": // Philippines: 10 digits starting with 9
		return regexp.MustCompile(`^9\d{9}$`).MatchString(phoneDigits)
	case "VN": // Vietnam: 9-10 digits starting with 9
		return regexp.MustCompile(`^9\d{8,9}$`).MatchString(phoneDigits)
	default:
		return false
	}
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

	if prefix, exists := countryPrefixes[u.CountryCode]; exists {
		// Remove leading zeros from local number
		localNumber := strings.TrimLeft(u.PhoneNumber, "0")
		return prefix + localNumber
	}

	return u.PhoneNumber
}

// IsActive checks if the user account is active
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive && (u.LockedUntil == nil || u.LockedUntil.Before(time.Now()))
}

// IsLocked checks if the user account is temporarily locked
func (u *User) IsLocked() bool {
	return u.LockedUntil != nil && u.LockedUntil.After(time.Now())
}

// LockAccount temporarily locks the user account
func (u *User) LockAccount(duration time.Duration) {
	lockUntil := time.Now().Add(duration)
	u.LockedUntil = &lockUntil
}

// UnlockAccount unlocks the user account and resets failed attempts
func (u *User) UnlockAccount() {
	u.LockedUntil = nil
	u.FailedAttempts = 0
}

// IncrementFailedAttempts increments the failed login attempts
func (u *User) IncrementFailedAttempts() {
	u.FailedAttempts++

	// Auto-lock after 5 failed attempts for 30 minutes
	if u.FailedAttempts >= 5 {
		u.LockAccount(30 * time.Minute)
	}
}

// UpdateLastLogin updates the last login information
func (u *User) UpdateLastLogin(ipAddress string) {
	now := time.Now()
	u.LastLoginAt = &now
	u.LastLoginIP = ipAddress
	u.FailedAttempts = 0 // Reset failed attempts on successful login
	u.LockedUntil = nil  // Unlock account on successful login
}

// GetAge returns the age of the user account
func (u *User) GetAge() time.Duration {
	return time.Since(u.CreatedAt)
}

// GetRegionalSettings returns regional settings for the user
func (u *User) GetRegionalSettings() map[string]interface{} {
	return map[string]interface{}{
		"country_code":  u.CountryCode,
		"locale":        u.Profile.Locale,
		"timezone":      u.Profile.Timezone,
		"data_region":   u.DataRegion,
		"phone_format":  u.GetFullPhoneNumber(),
		"currency":      GetDefaultCurrencyForCountry(u.CountryCode),
	}
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