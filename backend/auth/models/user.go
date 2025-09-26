package models

import (
	"regexp"
	"strings"
)

// Country represents a country code for regional compliance
type Country string

// Southeast Asian country codes
const (
	CountryThailand   Country = "TH"
	CountrySingapore  Country = "SG"
	CountryIndonesia  Country = "ID"
	CountryMalaysia   Country = "MY"
	CountryPhilippines Country = "PH"
	CountryVietnam    Country = "VN"
)

// KYC tiers for auth service compatibility
const (
	KYCTier1 = 1
	KYCTier2 = 2
	KYCTier3 = 3
)

// ValidCountries returns list of supported countries
var ValidCountries = []Country{
	CountryThailand,
	CountrySingapore,
	CountryIndonesia,
	CountryMalaysia,
	CountryPhilippines,
	CountryVietnam,
}

// Phone number validation patterns for Southeast Asia
var phonePatterns = map[Country]*regexp.Regexp{
	CountryThailand:   regexp.MustCompile(`^\+66[0-9]{9}$`),
	CountrySingapore:  regexp.MustCompile(`^\+65[0-9]{8}$`),
	CountryIndonesia:  regexp.MustCompile(`^\+62[0-9]{9,12}$`),
	CountryMalaysia:   regexp.MustCompile(`^\+60[0-9]{9,10}$`),
	CountryPhilippines: regexp.MustCompile(`^\+63[0-9]{10}$`),
	CountryVietnam:    regexp.MustCompile(`^\+84[0-9]{9,10}$`),
}

// IsValidPhoneNumber validates phone number format for the given country
func IsValidPhoneNumber(phoneNumber string, country Country) bool {
	pattern, exists := phonePatterns[country]
	if !exists {
		return false
	}
	return pattern.MatchString(phoneNumber)
}

// NormalizePhoneNumber normalizes a phone number to E.164 format
func NormalizePhoneNumber(phoneNumber string, country Country) string {
	// Remove all non-digit characters
	cleaned := regexp.MustCompile(`[^\d]`).ReplaceAllString(phoneNumber, "")

	// Add country code if not present
	switch country {
	case CountryThailand:
		if !strings.HasPrefix(phoneNumber, "+66") && !strings.HasPrefix(cleaned, "66") {
			return "+66" + strings.TrimPrefix(cleaned, "0")
		}
	case CountrySingapore:
		if !strings.HasPrefix(phoneNumber, "+65") && !strings.HasPrefix(cleaned, "65") {
			return "+65" + cleaned
		}
	case CountryIndonesia:
		if !strings.HasPrefix(phoneNumber, "+62") && !strings.HasPrefix(cleaned, "62") {
			return "+62" + strings.TrimPrefix(cleaned, "0")
		}
	case CountryMalaysia:
		if !strings.HasPrefix(phoneNumber, "+60") && !strings.HasPrefix(cleaned, "60") {
			return "+60" + strings.TrimPrefix(cleaned, "0")
		}
	case CountryPhilippines:
		if !strings.HasPrefix(phoneNumber, "+63") && !strings.HasPrefix(cleaned, "63") {
			return "+63" + strings.TrimPrefix(cleaned, "0")
		}
	case CountryVietnam:
		if !strings.HasPrefix(phoneNumber, "+84") && !strings.HasPrefix(cleaned, "84") {
			return "+84" + strings.TrimPrefix(cleaned, "0")
		}
	}

	// If already in international format, return as-is (add + if missing)
	if !strings.HasPrefix(phoneNumber, "+") {
		return "+" + cleaned
	}
	return phoneNumber
}

// FormatPhoneNumberLocal formats phone number for local display
func FormatPhoneNumberLocal(phoneNumber string, country Country) string {
	switch country {
	case CountryThailand:
		if strings.HasPrefix(phoneNumber, "+66") {
			local := strings.TrimPrefix(phoneNumber, "+66")
			if len(local) == 9 {
				return "0" + local[:2] + " " + local[2:5] + " " + local[5:]
			}
		}
	case CountrySingapore:
		if strings.HasPrefix(phoneNumber, "+65") {
			local := strings.TrimPrefix(phoneNumber, "+65")
			if len(local) == 8 {
				return local[:4] + " " + local[4:]
			}
		}
	}
	return phoneNumber
}

// FormatPhoneNumberInternational formats phone number for international display
func FormatPhoneNumberInternational(phoneNumber string, country Country) string {
	normalized := NormalizePhoneNumber(phoneNumber, country)
	switch country {
	case CountryThailand:
		if strings.HasPrefix(normalized, "+66") {
			return "+66 " + normalized[3:5] + " " + normalized[5:8] + " " + normalized[8:]
		}
	case CountrySingapore:
		if strings.HasPrefix(normalized, "+65") {
			return "+65 " + normalized[3:7] + " " + normalized[7:]
		}
	}
	return normalized
}