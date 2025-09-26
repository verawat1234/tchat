package fixtures

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BaseFixture provides common fixture utilities
type BaseFixture struct {
	seed int64
}

// NewBaseFixture creates a new base fixture with optional seed
func NewBaseFixture(seed ...int64) *BaseFixture {
	var s int64
	if len(seed) > 0 {
		s = seed[0]
	} else {
		s = time.Now().UnixNano()
	}
	return &BaseFixture{seed: s}
}

// UUID generates a deterministic UUID based on input
func (b *BaseFixture) UUID(input string) uuid.UUID {
	// Create deterministic UUID from input string
	hash := fmt.Sprintf("%s-%d", input, b.seed)
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(hash))
}

// RandomUUID generates a random UUID
func (b *BaseFixture) RandomUUID() uuid.UUID {
	return uuid.New()
}

// Phone generates a valid Southeast Asian phone number
func (b *BaseFixture) Phone(country string) string {
	switch country {
	case "TH": // Thailand
		return "+66812345678"
	case "SG": // Singapore
		return "+6512345678"
	case "ID": // Indonesia
		return "+628123456789"
	case "MY": // Malaysia
		return "+60123456789"
	case "VN": // Vietnam
		return "+84123456789"
	case "PH": // Philippines
		return "+639123456789"
	default:
		return "+66812345678" // Default to Thailand
	}
}

// Email generates a test email address
func (b *BaseFixture) Email(username string, domain ...string) string {
	d := "tchat-test.com"
	if len(domain) > 0 {
		d = domain[0]
	}
	return fmt.Sprintf("%s@%s", username, d)
}

// Name generates test names appropriate for Southeast Asian context
func (b *BaseFixture) Name(country string, gender ...string) string {
	g := "male"
	if len(gender) > 0 {
		g = gender[0]
	}

	switch country {
	case "TH": // Thailand
		if g == "female" {
			return "สมหญิง เซอร์วิส" // Somying Service
		}
		return "สมชาย เทสต์" // Somchai Test
	case "SG", "MY": // Singapore/Malaysia (Chinese names)
		if g == "female" {
			return "Li Wei Ming"
		}
		return "Tan Wei Hao"
	case "ID": // Indonesia
		if g == "female" {
			return "Sari Dewi"
		}
		return "Budi Santoso"
	case "VN": // Vietnam
		if g == "female" {
			return "Nguyễn Thị Lan"
		}
		return "Trần Văn Nam"
	case "PH": // Philippines
		if g == "female" {
			return "Maria Santos"
		}
		return "Jose Dela Cruz"
	default:
		return "Test User"
	}
}

// FutureTime generates a future timestamp
func (b *BaseFixture) FutureTime(minutes int) time.Time {
	return time.Now().UTC().Add(time.Duration(minutes) * time.Minute)
}

// PastTime generates a past timestamp
func (b *BaseFixture) PastTime(minutes int) time.Time {
	return time.Now().UTC().Add(-time.Duration(minutes) * time.Minute)
}

// Token generates a secure test token
func (b *BaseFixture) Token(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to deterministic token
		return fmt.Sprintf("test-token-%d-%d", b.seed, length)
	}
	return hex.EncodeToString(bytes)
}

// Amount generates test monetary amounts (in cents)
func (b *BaseFixture) Amount(currency string) int64 {
	switch currency {
	case "THB": // Thai Baht
		return 10000 // 100.00 THB
	case "SGD": // Singapore Dollar
		return 1000 // 10.00 SGD
	case "IDR": // Indonesian Rupiah
		return 15000000 // 150,000.00 IDR
	case "MYR": // Malaysian Ringgit
		return 4200 // 42.00 MYR
	case "VND": // Vietnamese Dong
		return 23000000 // 230,000.00 VND
	case "PHP": // Philippine Peso
		return 55000 // 550.00 PHP
	case "USD": // US Dollar (fallback)
		return 1000 // 10.00 USD
	default:
		return 1000 // Default amount
	}
}

// Currency returns currency code for country
func (b *BaseFixture) Currency(country string) string {
	switch country {
	case "TH":
		return "THB"
	case "SG":
		return "SGD"
	case "ID":
		return "IDR"
	case "MY":
		return "MYR"
	case "VN":
		return "VND"
	case "PH":
		return "PHP"
	default:
		return "USD"
	}
}

// LoremText generates lorem ipsum text for content
func (b *BaseFixture) LoremText(words int) string {
	lorem := []string{
		"Lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
		"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
		"magna", "aliqua", "Ut", "enim", "ad", "minim", "veniam", "quis", "nostrud",
		"exercitation", "ullamco", "laboris", "nisi", "ut", "aliquip", "ex", "ea",
		"commodo", "consequat", "Duis", "aute", "irure", "in", "reprehenderit",
		"voluptate", "velit", "esse", "cillum", "fugiat", "nulla", "pariatur",
	}

	if words <= 0 {
		words = 10
	}

	result := make([]string, 0, words)
	for i := 0; i < words; i++ {
		result = append(result, lorem[i%len(lorem)])
	}

	return fmt.Sprintf("%s.", fmt.Sprintf("%s", result))
}

// SEAContent generates Southeast Asian appropriate content
func (b *BaseFixture) SEAContent(country, contentType string) string {
	switch country {
	case "TH":
		switch contentType {
		case "greeting":
			return "สวัสดีครับ/ค่ะ"
		case "product":
			return "สินค้าคุณภาพดี ราคาถูก"
		default:
			return "เนื้อหาทดสอบ"
		}
	case "VN":
		switch contentType {
		case "greeting":
			return "Xin chào"
		case "product":
			return "Sản phẩm chất lượng cao"
		default:
			return "Nội dung thử nghiệm"
		}
	case "ID":
		switch contentType {
		case "greeting":
			return "Selamat datang"
		case "product":
			return "Produk berkualitas tinggi"
		default:
			return "Konten uji coba"
		}
	default:
		return b.LoremText(5)
	}
}

// DeviceID generates realistic device identifiers
func (b *BaseFixture) DeviceID(platform string) string {
	switch platform {
	case "ios":
		return fmt.Sprintf("ios-device-%s", b.Token(8))
	case "android":
		return fmt.Sprintf("android-device-%s", b.Token(8))
	case "web":
		return fmt.Sprintf("web-session-%s", b.Token(8))
	default:
		return fmt.Sprintf("device-%s", b.Token(8))
	}
}

// UserAgent generates realistic user agent strings
func (b *BaseFixture) UserAgent(platform string) string {
	switch platform {
	case "ios":
		return "Tchat/1.0 (iPhone; iOS 17.0; Scale/3.00)"
	case "android":
		return "Tchat/1.0 (Android 14; Mobile; Samsung SM-G991B)"
	case "web":
		return "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Tchat/1.0"
	default:
		return "Tchat/1.0 Test Client"
	}
}

// IPAddress generates test IP addresses
func (b *BaseFixture) IPAddress(country string) string {
	switch country {
	case "TH": // Thailand IP range
		return "203.144.32.1"
	case "SG": // Singapore IP range
		return "103.233.0.1"
	case "ID": // Indonesia IP range
		return "36.92.0.1"
	case "MY": // Malaysia IP range
		return "103.16.0.1"
	case "VN": // Vietnam IP range
		return "14.224.0.1"
	case "PH": // Philippines IP range
		return "49.144.0.1"
	default:
		return "127.0.0.1" // Localhost for testing
	}
}

// Status generates valid status values for different models
func (b *BaseFixture) Status(modelType string) string {
	switch modelType {
	case "user":
		return "active"
	case "content":
		return "published"
	case "payment":
		return "completed"
	case "session":
		return "active"
	case "notification":
		return "sent"
	default:
		return "active"
	}
}

// KYCTier generates valid KYC tier values
func (b *BaseFixture) KYCTier() int {
	return 1 // Default to Tier 1 (basic verification)
}

// CountryCode returns standardized country codes
func (b *BaseFixture) CountryCode(country string) string {
	switch country {
	case "thailand":
		return "TH"
	case "singapore":
		return "SG"
	case "indonesia":
		return "ID"
	case "malaysia":
		return "MY"
	case "vietnam":
		return "VN"
	case "philippines":
		return "PH"
	default:
		return country
	}
}

// Locale returns appropriate locale for country
func (b *BaseFixture) Locale(country string) string {
	switch country {
	case "TH":
		return "th-TH"
	case "SG":
		return "en-SG"
	case "ID":
		return "id-ID"
	case "MY":
		return "ms-MY"
	case "VN":
		return "vi-VN"
	case "PH":
		return "en-PH"
	default:
		return "en-US"
	}
}

// RandomInt generates a random integer between min and max (inclusive)
func (b *BaseFixture) RandomInt(min, max int) int {
	if min > max {
		min, max = max, min
	}

	// Use seed for deterministic behavior if needed
	// For testing purposes, we can use a simple modulo approach
	range_ := max - min + 1
	return min + int(b.seed%int64(range_))
}