package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the authentication system
type User struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	Phone      *string    `json:"phone,omitempty" db:"phone"`
	Email      *string    `json:"email,omitempty" db:"email"`
	Name       string     `json:"name" db:"name"`
	Avatar     *string    `json:"avatar,omitempty" db:"avatar"`
	Country    Country    `json:"country" db:"country"`
	Locale     string     `json:"locale" db:"locale"`
	KYCTier    KYCTier    `json:"kyc_tier" db:"kyc_tier"`
	Status     UserStatus `json:"status" db:"status"`
	LastSeen   *time.Time `json:"last_seen,omitempty" db:"last_seen"`
	IsVerified bool       `json:"is_verified" db:"is_verified"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

// Country represents supported Southeast Asian countries
type Country string

const (
	CountryThailand   Country = "TH"
	CountryIndonesia  Country = "ID"
	CountryMalaysia   Country = "MY"
	CountryVietnam    Country = "VN"
	CountrySingapore  Country = "SG"
	CountryPhilippines Country = "PH"
)

// ValidCountries returns all supported countries
func ValidCountries() []Country {
	return []Country{
		CountryThailand,
		CountryIndonesia,
		CountryMalaysia,
		CountryVietnam,
		CountrySingapore,
		CountryPhilippines,
	}
}

// IsValid validates if the country code is supported
func (c Country) IsValid() bool {
	for _, valid := range ValidCountries() {
		if c == valid {
			return true
		}
	}
	return false
}

// String returns the string representation of Country
func (c Country) String() string {
	return string(c)
}

// Value implements the driver.Valuer interface for database storage
func (c Country) Value() (driver.Value, error) {
	return string(c), nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (c *Country) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if str, ok := value.(string); ok {
		*c = Country(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into Country", value)
}

// KYCTier represents Know Your Customer verification levels
type KYCTier int

const (
	KYCTier1 KYCTier = 1 // Basic verification (phone/email)
	KYCTier2 KYCTier = 2 // Identity verification (ID document)
	KYCTier3 KYCTier = 3 // Enhanced verification (proof of address)
)

// IsValid validates if the KYC tier is supported
func (k KYCTier) IsValid() bool {
	return k >= KYCTier1 && k <= KYCTier3
}

// String returns the string representation of KYCTier
func (k KYCTier) String() string {
	switch k {
	case KYCTier1:
		return "Basic"
	case KYCTier2:
		return "Identity"
	case KYCTier3:
		return "Enhanced"
	default:
		return "Unknown"
	}
}

// UserStatus represents the current status of a user
type UserStatus string

const (
	UserStatusOnline UserStatus = "online"
	UserStatusOffline UserStatus = "offline"
	UserStatusAway   UserStatus = "away"
	UserStatusBusy   UserStatus = "busy"
)

// IsValid validates if the user status is supported
func (s UserStatus) IsValid() bool {
	validStatuses := []UserStatus{
		UserStatusOnline,
		UserStatusOffline,
		UserStatusAway,
		UserStatusBusy,
	}
	for _, valid := range validStatuses {
		if s == valid {
			return true
		}
	}
	return false
}

// Validate performs comprehensive validation on the User model
func (u *User) Validate() error {
	var errs []string

	// Phone OR Email required (not both null)
	if u.Phone == nil && u.Email == nil {
		errs = append(errs, "either phone or email is required")
	}

	// Validate phone format if provided
	if u.Phone != nil && *u.Phone != "" {
		if err := u.validatePhoneFormat(*u.Phone); err != nil {
			errs = append(errs, fmt.Sprintf("invalid phone format: %v", err))
		}
	}

	// Validate email format if provided
	if u.Email != nil && *u.Email != "" {
		if err := u.validateEmailFormat(*u.Email); err != nil {
			errs = append(errs, fmt.Sprintf("invalid email format: %v", err))
		}
	}

	// Name validation (2-100 characters)
	if len(strings.TrimSpace(u.Name)) < 2 {
		errs = append(errs, "name must be at least 2 characters")
	}
	if len(u.Name) > 100 {
		errs = append(errs, "name must not exceed 100 characters")
	}

	// Country validation
	if !u.Country.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid country code: %s", u.Country))
	}

	// Locale validation
	if err := u.validateLocale(); err != nil {
		errs = append(errs, fmt.Sprintf("invalid locale: %v", err))
	}

	// KYC tier validation
	if !u.KYCTier.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid KYC tier: %d", u.KYCTier))
	}

	// User status validation
	if !u.Status.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid user status: %s", u.Status))
	}

	// Avatar URL validation if provided
	if u.Avatar != nil && *u.Avatar != "" {
		if err := u.validateAvatarURL(*u.Avatar); err != nil {
			errs = append(errs, fmt.Sprintf("invalid avatar URL: %v", err))
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

// validatePhoneFormat validates phone number format based on country
func (u *User) validatePhoneFormat(phone string) error {
	// Country-specific phone validation patterns
	patterns := map[Country]string{
		CountryThailand:   `^\+66[0-9]{8,9}$`,
		CountryIndonesia:  `^\+62[0-9]{8,12}$`,
		CountryMalaysia:   `^\+60[0-9]{8,10}$`,
		CountryVietnam:    `^\+84[0-9]{8,10}$`,
		CountrySingapore:  `^\+65[0-9]{8}$`,
		CountryPhilippines: `^\+63[0-9]{9,10}$`,
	}

	pattern, exists := patterns[u.Country]
	if !exists {
		return fmt.Errorf("no phone validation pattern for country %s", u.Country)
	}

	matched, err := regexp.MatchString(pattern, phone)
	if err != nil {
		return fmt.Errorf("regex error: %v", err)
	}

	if !matched {
		return fmt.Errorf("phone number does not match pattern for country %s", u.Country)
	}

	return nil
}

// validateEmailFormat validates email format using RFC 5322 regex
func (u *User) validateEmailFormat(email string) error {
	// Basic email validation pattern
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return fmt.Errorf("regex error: %v", err)
	}

	if !matched {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// validateLocale validates locale format (language-country)
func (u *User) validateLocale() error {
	// Expected locale formats for supported countries
	validLocales := map[Country][]string{
		CountryThailand:   {"th-TH", "en-TH"},
		CountryIndonesia:  {"id-ID", "en-ID"},
		CountryMalaysia:   {"ms-MY", "en-MY", "zh-MY"},
		CountryVietnam:    {"vi-VN", "en-VN"},
		CountrySingapore:  {"en-SG", "ms-SG", "zh-SG", "ta-SG"},
		CountryPhilippines: {"fil-PH", "en-PH"},
	}

	locales, exists := validLocales[u.Country]
	if !exists {
		return fmt.Errorf("no valid locales defined for country %s", u.Country)
	}

	for _, validLocale := range locales {
		if u.Locale == validLocale {
			return nil
		}
	}

	return fmt.Errorf("locale %s is not valid for country %s", u.Locale, u.Country)
}

// validateAvatarURL validates avatar URL format
func (u *User) validateAvatarURL(url string) error {
	// Basic URL validation pattern
	pattern := `^https?://[^\s/$.?#].[^\s]*$`
	matched, err := regexp.MatchString(pattern, url)
	if err != nil {
		return fmt.Errorf("regex error: %v", err)
	}

	if !matched {
		return fmt.Errorf("invalid URL format")
	}

	return nil
}

// BeforeCreate sets up the user before database creation
func (u *User) BeforeCreate() error {
	// Generate UUID if not set
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now

	// Set default values
	if u.KYCTier == 0 {
		u.KYCTier = KYCTier1
	}
	if u.Status == "" {
		u.Status = UserStatusOffline
	}

	// Validate before creation
	return u.Validate()
}

// BeforeUpdate sets up the user before database update
func (u *User) BeforeUpdate() error {
	// Update timestamp
	u.UpdatedAt = time.Now().UTC()

	// Validate before update
	return u.Validate()
}

// GetDisplayName returns the user's display name
func (u *User) GetDisplayName() string {
	if u.Name != "" {
		return u.Name
	}
	if u.Phone != nil {
		return *u.Phone
	}
	if u.Email != nil {
		return *u.Email
	}
	return "Unknown User"
}

// IsOnline checks if the user is currently online
func (u *User) IsOnline() bool {
	return u.Status == UserStatusOnline
}

// CanSendPayments checks if user can send payments based on KYC tier
func (u *User) CanSendPayments() bool {
	return u.KYCTier >= KYCTier2 && u.IsVerified
}

// GetMaxDailyLimit returns maximum daily transaction limit based on KYC tier
func (u *User) GetMaxDailyLimit() int64 {
	switch u.KYCTier {
	case KYCTier1:
		return 100000 // 1,000 THB equivalent in cents
	case KYCTier2:
		return 5000000 // 50,000 THB equivalent in cents
	case KYCTier3:
		return 50000000 // 500,000 THB equivalent in cents
	default:
		return 0
	}
}

// UpdateLastSeen updates the user's last seen timestamp
func (u *User) UpdateLastSeen() {
	now := time.Now().UTC()
	u.LastSeen = &now
	u.UpdatedAt = now
}

// ToPublicUser returns a sanitized version for public API responses
func (u *User) ToPublicUser() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"name":       u.Name,
		"avatar":     u.Avatar,
		"country":    u.Country,
		"status":     u.Status,
		"is_verified": u.IsVerified,
		"created_at": u.CreatedAt,
	}
}

// ToResponse returns a response structure for API responses
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:         u.ID,
		Name:       u.Name,
		Avatar:     u.Avatar,
		Country:    string(u.Country),
		Locale:     u.Locale,
		KYCTier:    int(u.KYCTier),
		Status:     string(u.Status),
		IsVerified: u.IsVerified,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

// UserResponse represents user data for API responses
type UserResponse struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Avatar     *string    `json:"avatar,omitempty"`
	Country    string     `json:"country"`
	Locale     string     `json:"locale"`
	KYCTier    int        `json:"kyc_tier"`
	Status     string     `json:"status"`
	IsVerified bool       `json:"is_verified"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// Global phone number utility functions

// IsValidPhoneNumber validates phone number format for a given country
func IsValidPhoneNumber(phone, countryCode string) bool {
	// Country-specific phone validation patterns
	patterns := map[string]string{
		"TH": `^\+66[0-9]{8,9}$`,
		"ID": `^\+62[0-9]{8,12}$`,
		"MY": `^\+60[0-9]{8,10}$`,
		"VN": `^\+84[0-9]{8,10}$`,
		"SG": `^\+65[0-9]{8}$`,
		"PH": `^\+63[0-9]{9,10}$`,
	}

	pattern, exists := patterns[countryCode]
	if !exists {
		return false
	}

	matched, err := regexp.MatchString(pattern, phone)
	return err == nil && matched
}

// NormalizePhoneNumber normalizes phone number to E.164 format
func NormalizePhoneNumber(phone string) string {
	// Remove all non-digit characters except +
	phone = strings.TrimSpace(phone)

	// Basic normalization - in production, use a library like libphonenumber
	if !strings.HasPrefix(phone, "+") {
		return phone // Return as-is if no country code
	}

	return phone
}

// FormatPhoneNumberLocal formats phone number for local display
func FormatPhoneNumberLocal(phone, countryCode string) string {
	// Remove country code for local formatting
	// This is a simplified implementation
	switch countryCode {
	case "TH":
		if strings.HasPrefix(phone, "+66") {
			local := strings.TrimPrefix(phone, "+66")
			if len(local) >= 8 {
				return fmt.Sprintf("0%s-%s-%s", local[:2], local[2:5], local[5:])
			}
		}
	case "SG":
		if strings.HasPrefix(phone, "+65") {
			local := strings.TrimPrefix(phone, "+65")
			if len(local) == 8 {
				return fmt.Sprintf("%s %s", local[:4], local[4:])
			}
		}
	case "ID":
		if strings.HasPrefix(phone, "+62") {
			local := strings.TrimPrefix(phone, "+62")
			return fmt.Sprintf("0%s", local)
		}
	}

	// Fallback to international format
	return phone
}

// FormatPhoneNumberInternational formats phone number for international display
func FormatPhoneNumberInternational(phone string) string {
	// Ensure it's in E.164 format
	if !strings.HasPrefix(phone, "+") {
		return phone // Return as-is if not international format
	}

	// Add spacing for readability
	if len(phone) > 3 {
		countryCode := phone[1:4] // Assume 3-digit country code
		number := phone[4:]

		// Add spacing every 3-4 digits
		if len(number) > 6 {
			return fmt.Sprintf("+%s %s %s", countryCode, number[:3], number[3:])
		} else if len(number) > 3 {
			return fmt.Sprintf("+%s %s", countryCode, number)
		}
	}

	return phone
}

// Additional UserStatus constants for compatibility with services
const (
	UserStatusPending UserStatus = "pending"
	UserStatusActive  UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusInactive UserStatus = "inactive"
)

// Additional methods needed by UserService

// CanUpdateProfile checks if user can update their profile
func (u *User) CanUpdateProfile() bool {
	return u.Status == UserStatusActive || u.Status == UserStatusOnline
}

// CanTransitionToStatus checks if user can transition to given status
func (u *User) CanTransitionToStatus(newStatus UserStatus) bool {
	switch u.Status {
	case UserStatusPending:
		return newStatus == UserStatusActive || newStatus == UserStatusInactive
	case UserStatusActive:
		return newStatus == UserStatusSuspended || newStatus == UserStatusInactive || newStatus == UserStatusOnline || newStatus == UserStatusOffline
	case UserStatusSuspended:
		return newStatus == UserStatusActive || newStatus == UserStatusInactive
	case UserStatusInactive:
		return newStatus == UserStatusActive
	default:
		return true // Allow any transition for other statuses
	}
}

// CanUpgradeToTier checks if user can upgrade to given KYC tier
func (u *User) CanUpgradeToTier(newTier KYCTier) bool {
	return newTier > u.KYCTier && newTier.IsValid()
}

// CanBeDeleted checks if user can be deleted
func (u *User) CanBeDeleted() bool {
	return u.Status == UserStatusInactive || u.Status == UserStatusSuspended
}

// CanSubmitKYC checks if user can submit KYC verification
func (u *User) CanSubmitKYC() bool {
	return u.Status == UserStatusActive && u.IsVerified
}

// Global validation helper functions

// IsValidCountry checks if country code is valid
func IsValidCountry(countryCode string) bool {
	country := Country(countryCode)
	return country.IsValid()
}

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	return err == nil && matched
}

// IsValidUsername validates username format
func IsValidUsername(username string) bool {
	// Username should be 3-30 characters, alphanumeric + underscore
	pattern := `^[a-zA-Z0-9_]{3,30}$`
	matched, err := regexp.MatchString(pattern, username)
	return err == nil && matched
}