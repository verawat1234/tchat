package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
)

// ValidationResult holds validation results
type ValidationResult struct {
	IsValid bool
	Errors  []string
}

// Validator provides common validation functions
type Validator struct {
	errors []string
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		errors: make([]string, 0),
	}
}

// AddError adds a validation error
func (v *Validator) AddError(field, message string) {
	v.errors = append(v.errors, fmt.Sprintf("%s: %s", field, message))
}

// AddErrorf adds a formatted validation error
func (v *Validator) AddErrorf(field, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	v.AddError(field, message)
}

// HasErrors returns true if there are validation errors
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// GetErrors returns all validation errors
func (v *Validator) GetErrors() []string {
	return v.errors
}

// GetError returns a single error combining all validation errors
func (v *Validator) GetError() error {
	if !v.HasErrors() {
		return nil
	}
	return errors.New(strings.Join(v.errors, "; "))
}

// Reset clears all validation errors
func (v *Validator) Reset() {
	v.errors = make([]string, 0)
}

// String validation functions

// Required validates that a string is not empty
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, "is required")
	}
	return v
}

// MinLength validates minimum string length
func (v *Validator) MinLength(field, value string, min int) *Validator {
	if utf8.RuneCountInString(value) < min {
		v.AddErrorf(field, "must be at least %d characters long", min)
	}
	return v
}

// MaxLength validates maximum string length
func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if utf8.RuneCountInString(value) > max {
		v.AddErrorf(field, "must not exceed %d characters", max)
	}
	return v
}

// Length validates exact string length
func (v *Validator) Length(field, value string, length int) *Validator {
	if utf8.RuneCountInString(value) != length {
		v.AddErrorf(field, "must be exactly %d characters long", length)
	}
	return v
}

// Pattern validates string against regex pattern
func (v *Validator) Pattern(field, value, pattern string) *Validator {
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		v.AddErrorf(field, "pattern validation failed: %v", err)
		return v
	}
	if !matched {
		v.AddError(field, "does not match required format")
	}
	return v
}

// OneOf validates that value is one of the allowed values
func (v *Validator) OneOf(field, value string, allowed []string) *Validator {
	for _, allowed := range allowed {
		if value == allowed {
			return v
		}
	}
	v.AddErrorf(field, "must be one of: %s", strings.Join(allowed, ", "))
	return v
}

// Numeric validation functions

// Min validates minimum numeric value
func (v *Validator) Min(field string, value, min int64) *Validator {
	if value < min {
		v.AddErrorf(field, "must be at least %d", min)
	}
	return v
}

// Max validates maximum numeric value
func (v *Validator) Max(field string, value, max int64) *Validator {
	if value > max {
		v.AddErrorf(field, "must not exceed %d", max)
	}
	return v
}

// Range validates numeric value within range
func (v *Validator) Range(field string, value, min, max int64) *Validator {
	if value < min || value > max {
		v.AddErrorf(field, "must be between %d and %d", min, max)
	}
	return v
}

// Positive validates positive numeric value
func (v *Validator) Positive(field string, value int64) *Validator {
	if value <= 0 {
		v.AddError(field, "must be positive")
	}
	return v
}

// NonNegative validates non-negative numeric value
func (v *Validator) NonNegative(field string, value int64) *Validator {
	if value < 0 {
		v.AddError(field, "must be non-negative")
	}
	return v
}

// Format validation functions

// Email validates email format
func (v *Validator) Email(field, value string) *Validator {
	if value == "" {
		return v // Allow empty if not required
	}

	_, err := mail.ParseAddress(value)
	if err != nil {
		v.AddError(field, "must be a valid email address")
	}
	return v
}

// URL validates URL format
func (v *Validator) URL(field, value string) *Validator {
	if value == "" {
		return v // Allow empty if not required
	}

	_, err := url.ParseRequestURI(value)
	if err != nil {
		v.AddError(field, "must be a valid URL")
	}
	return v
}

// UUID validates UUID format
func (v *Validator) UUID(field, value string) *Validator {
	if value == "" {
		return v // Allow empty if not required
	}

	_, err := uuid.Parse(value)
	if err != nil {
		v.AddError(field, "must be a valid UUID")
	}
	return v
}

// IP validates IP address format
func (v *Validator) IP(field, value string) *Validator {
	if value == "" {
		return v // Allow empty if not required
	}

	ip := net.ParseIP(value)
	if ip == nil {
		v.AddError(field, "must be a valid IP address")
	}
	return v
}

// Phone validates phone number format
func (v *Validator) Phone(field, value string) *Validator {
	if value == "" {
		return v // Allow empty if not required
	}

	// Basic international phone number validation
	phonePattern := `^\+[1-9]\d{1,14}$`
	matched, err := regexp.MatchString(phonePattern, value)
	if err != nil {
		v.AddErrorf(field, "phone validation failed: %v", err)
		return v
	}
	if !matched {
		v.AddError(field, "must be a valid international phone number (+country code followed by number)")
	}
	return v
}

// Time validation functions

// After validates that time is after the specified time
func (v *Validator) After(field string, value, after time.Time) *Validator {
	if !value.After(after) {
		v.AddErrorf(field, "must be after %s", after.Format(time.RFC3339))
	}
	return v
}

// Before validates that time is before the specified time
func (v *Validator) Before(field string, value, before time.Time) *Validator {
	if !value.Before(before) {
		v.AddErrorf(field, "must be before %s", before.Format(time.RFC3339))
	}
	return v
}

// DateRange validates that date is within range
func (v *Validator) DateRange(field string, value, start, end time.Time) *Validator {
	if value.Before(start) || value.After(end) {
		v.AddErrorf(field, "must be between %s and %s", start.Format(time.RFC3339), end.Format(time.RFC3339))
	}
	return v
}

// Collection validation functions

// SliceNotEmpty validates that slice is not empty
func (v *Validator) SliceNotEmpty(field string, slice []interface{}) *Validator {
	if len(slice) == 0 {
		v.AddError(field, "must not be empty")
	}
	return v
}

// SliceMinLength validates minimum slice length
func (v *Validator) SliceMinLength(field string, slice []interface{}, min int) *Validator {
	if len(slice) < min {
		v.AddErrorf(field, "must contain at least %d items", min)
	}
	return v
}

// SliceMaxLength validates maximum slice length
func (v *Validator) SliceMaxLength(field string, slice []interface{}, max int) *Validator {
	if len(slice) > max {
		v.AddErrorf(field, "must not contain more than %d items", max)
	}
	return v
}

// SliceLength validates exact slice length
func (v *Validator) SliceLength(field string, slice []interface{}, length int) *Validator {
	if len(slice) != length {
		v.AddErrorf(field, "must contain exactly %d items", length)
	}
	return v
}

// MapNotEmpty validates that map is not empty
func (v *Validator) MapNotEmpty(field string, m map[string]interface{}) *Validator {
	if len(m) == 0 {
		v.AddError(field, "must not be empty")
	}
	return v
}

// Advanced validation functions

// AlphaNumeric validates alphanumeric characters only
func (v *Validator) AlphaNumeric(field, value string) *Validator {
	if value == "" {
		return v
	}

	for _, r := range value {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			v.AddError(field, "must contain only letters and numbers")
			break
		}
	}
	return v
}

// Alpha validates alphabetic characters only
func (v *Validator) Alpha(field, value string) *Validator {
	if value == "" {
		return v
	}

	for _, r := range value {
		if !unicode.IsLetter(r) {
			v.AddError(field, "must contain only letters")
			break
		}
	}
	return v
}

// Numeric validates numeric characters only
func (v *Validator) Numeric(field, value string) *Validator {
	if value == "" {
		return v
	}

	for _, r := range value {
		if !unicode.IsDigit(r) {
			v.AddError(field, "must contain only numbers")
			break
		}
	}
	return v
}

// NoWhitespace validates no whitespace characters
func (v *Validator) NoWhitespace(field, value string) *Validator {
	if value == "" {
		return v
	}

	for _, r := range value {
		if unicode.IsSpace(r) {
			v.AddError(field, "must not contain whitespace")
			break
		}
	}
	return v
}

// PrintableASCII validates printable ASCII characters only
func (v *Validator) PrintableASCII(field, value string) *Validator {
	if value == "" {
		return v
	}

	for _, r := range value {
		if r < 32 || r > 126 {
			v.AddError(field, "must contain only printable ASCII characters")
			break
		}
	}
	return v
}

// Password validates password strength
func (v *Validator) Password(field, value string, minLength int, requireSpecial bool) *Validator {
	if value == "" {
		return v
	}

	// Check minimum length
	if len(value) < minLength {
		v.AddErrorf(field, "must be at least %d characters long", minLength)
	}

	// Check for required character types
	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, r := range value {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if !hasUpper {
		v.AddError(field, "must contain at least one uppercase letter")
	}
	if !hasLower {
		v.AddError(field, "must contain at least one lowercase letter")
	}
	if !hasDigit {
		v.AddError(field, "must contain at least one digit")
	}
	if requireSpecial && !hasSpecial {
		v.AddError(field, "must contain at least one special character")
	}

	return v
}

// CreditCard validates credit card number using Luhn algorithm
func (v *Validator) CreditCard(field, value string) *Validator {
	if value == "" {
		return v
	}

	// Remove spaces and dashes
	cleaned := strings.ReplaceAll(strings.ReplaceAll(value, " ", ""), "-", "")

	// Check if all characters are digits
	for _, r := range cleaned {
		if !unicode.IsDigit(r) {
			v.AddError(field, "must contain only digits")
			return v
		}
	}

	// Check length (most cards are 13-19 digits)
	if len(cleaned) < 13 || len(cleaned) > 19 {
		v.AddError(field, "must be between 13 and 19 digits")
		return v
	}

	// Validate using Luhn algorithm
	if !v.luhnCheck(cleaned) {
		v.AddError(field, "is not a valid credit card number")
	}

	return v
}

// luhnCheck implements the Luhn algorithm for credit card validation
func (v *Validator) luhnCheck(number string) bool {
	var sum int
	alternate := false

	// Process digits from right to left
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}

// JSON validates JSON format
func (v *Validator) JSON(field, value string) *Validator {
	if value == "" {
		return v
	}

	if !IsValidJSON(value) {
		v.AddError(field, "must be valid JSON")
	}
	return v
}

// Custom validation function
type CustomValidationFunc func(value interface{}) error

// Custom validates using a custom function
func (v *Validator) Custom(field string, value interface{}, fn CustomValidationFunc) *Validator {
	if err := fn(value); err != nil {
		v.AddError(field, err.Error())
	}
	return v
}

// Utility functions

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// IsValidURL validates URL format
func IsValidURL(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}

// IsValidUUID validates UUID format
func IsValidUUID(uuidStr string) bool {
	_, err := uuid.Parse(uuidStr)
	return err == nil
}

// IsValidIP validates IP address format
func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsValidPhone validates international phone number format
func IsValidPhone(phone string) bool {
	phonePattern := `^\+[1-9]\d{1,14}$`
	matched, err := regexp.MatchString(phonePattern, phone)
	return err == nil && matched
}

// IsValidJSON validates JSON format
func IsValidJSON(jsonStr string) bool {
	var js interface{}
	return json.Unmarshal([]byte(jsonStr), &js) == nil
}

// IsAlphaNumeric checks if string contains only alphanumeric characters
func IsAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsAlpha checks if string contains only alphabetic characters
func IsAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// IsNumeric checks if string contains only numeric characters
func IsNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsPrintableASCII checks if string contains only printable ASCII characters
func IsPrintableASCII(s string) bool {
	for _, r := range s {
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}

// SanitizeString removes or replaces dangerous characters
func SanitizeString(s string) string {
	// Remove null bytes
	s = strings.ReplaceAll(s, "\x00", "")

	// Remove control characters except tab, newline, and carriage return
	var result strings.Builder
	for _, r := range s {
		if r >= 32 || r == '\t' || r == '\n' || r == '\r' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// TruncateString truncates string to maximum length
func TruncateString(s string, maxLength int) string {
	if utf8.RuneCountInString(s) <= maxLength {
		return s
	}

	runes := []rune(s)
	if len(runes) > maxLength {
		return string(runes[:maxLength])
	}

	return s
}

// NormalizeEmail normalizes email address (lowercase, trim)
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// NormalizePhone normalizes phone number (remove non-digits except +)
func NormalizePhone(phone string) string {
	var result strings.Builder
	for i, r := range phone {
		if unicode.IsDigit(r) || (i == 0 && r == '+') {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Validation helpers for common scenarios

// ValidateUserInput validates common user input fields
func ValidateUserInput(name, email, phone string) *ValidationResult {
	v := NewValidator()

	v.Required("name", name).MinLength("name", name, 2).MaxLength("name", name, 100)

	if email != "" {
		v.Email("email", email)
	}

	if phone != "" {
		v.Phone("phone", phone)
	}

	return &ValidationResult{
		IsValid: !v.HasErrors(),
		Errors:  v.GetErrors(),
	}
}

// ValidatePassword validates password strength
func ValidatePassword(password string) *ValidationResult {
	v := NewValidator()
	v.Required("password", password).Password("password", password, 8, true)

	return &ValidationResult{
		IsValid: !v.HasErrors(),
		Errors:  v.GetErrors(),
	}
}

// ValidateProductInput validates product input
func ValidateProductInput(name, description string, price int64, inventory int) *ValidationResult {
	v := NewValidator()

	v.Required("name", name).MinLength("name", name, 3).MaxLength("name", name, 200)

	if description != "" {
		v.MaxLength("description", description, 2000)
	}

	v.Positive("price", price)
	v.NonNegative("inventory", int64(inventory))

	return &ValidationResult{
		IsValid: !v.HasErrors(),
		Errors:  v.GetErrors(),
	}
}

// Southeast Asian specific validation functions

// IsValidPhoneNumber validates phone number format for Southeast Asian countries
func IsValidPhoneNumber(phoneNumber, countryCode string) bool {
	// Remove any non-digit characters for validation
	phoneDigits := regexp.MustCompile(`\D`).ReplaceAllString(phoneNumber, "")

	// Country-specific phone number validation
	switch countryCode {
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

// IsValidSEACountryCode validates if a country code is supported in Southeast Asia
func IsValidSEACountryCode(countryCode string) bool {
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

// IsValidEmailAddress validates email format (alternate name for consistency)
func IsValidEmailAddress(email string) bool {
	return IsValidEmail(email)
}

// IsValidUsername validates username format
func IsValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 30 {
		return false
	}

	// Username should start with letter or number, can contain letters, numbers, underscores, and hyphens
	for i, r := range username {
		if i == 0 {
			// First character must be letter or number
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
				return false
			}
		} else {
			// Other characters can be letters, numbers, underscores, or hyphens
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-') {
				return false
			}
		}
	}

	return true
}