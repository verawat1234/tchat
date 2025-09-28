package fixtures

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"tchat.dev/shared/models"
)

// CommerceFixtures provides fixtures for commerce-related models
type CommerceFixtures struct {
	*BaseFixture
}

// NewCommerceFixtures creates a new commerce fixtures instance
func NewCommerceFixtures(seed ...int64) *CommerceFixtures {
	return &CommerceFixtures{
		BaseFixture: NewBaseFixture(seed...),
	}
}

// Southeast Asian Business Data
var SEABusinessCategories = []string{
	"electronics", "fashion", "food", "health", "beauty", "home",
	"sports", "automotive", "books", "toys", "crafts", "jewelry",
}

var SEABusinessNames = map[string][]string{
	"TH": {"Bangkok Electronics", "Siam Fashion", "Thai Spice Co", "Wellness Thailand", "Thai Beauty Hub"},
	"SG": {"Singapore Tech", "Marina Fashion", "Lion City Foods", "SG Health Plus", "Beauty Singapore"},
	"ID": {"Jakarta Digital", "Indonesian Style", "Nusantara Foods", "Indo Health", "Bali Beauty Co"},
	"MY": {"KL Electronics", "Malaysia Fashion", "Mamak Kitchen", "Health Malaysia", "Beauty KL"},
	"VN": {"Saigon Tech", "Vietnam Fashion", "Pho Kitchen", "Vietnam Health", "Hanoi Beauty"},
	"PH": {"Manila Electronics", "Filipino Fashion", "Adobo House", "PH Health Care", "Manila Beauty"},
}

// BasicBusiness creates a basic business for testing
func (cf *CommerceFixtures) BasicBusiness(country string) *models.Business {
	businessNames := SEABusinessNames[country]
	if len(businessNames) == 0 {
		businessNames = SEABusinessNames["TH"] // fallback
	}

	name := businessNames[cf.RandomInt(0, len(businessNames))]
	category := SEABusinessCategories[cf.RandomInt(0, len(SEABusinessCategories))]
	ownerID := cf.UUID(fmt.Sprintf("owner-%s-%s", country, name))

	return &models.Business{
		ID:          cf.UUID(fmt.Sprintf("business-%s-%s", country, name)),
		OwnerID:     ownerID,
		Name:        name,
		Description: fmt.Sprintf("A premier %s business serving %s market", category, country),
		Category:    category,

		VerificationStatus: models.BusinessVerificationPending,

		ContactInfo: models.BusinessContactInfo{
			Phone:   cf.Phone(country),
			Email:   fmt.Sprintf("contact@%s.com", name),
			Website: fmt.Sprintf("https://%s.com", name),
		},

		Address: models.BusinessAddress{
			Street:     fmt.Sprintf("%d Main Street", cf.RandomInt(1, 999)),
			City:       cf.getCity(country),
			State:      cf.getState(country),
			PostalCode: cf.getPostalCode(country),
			Country:    country,
		},

		BusinessSettings: models.BusinessSettings{
			SupportedCurrencies: []string{cf.Currency(country)},
			SupportedLanguages:  []string{"en", cf.getLanguage(country)},
			ShippingCountries:   []string{country},
			PaymentMethods:      []string{"card", "bank_transfer", "digital_wallet"},
		},

		IsActive:  true,
		CreatedAt: cf.PastTime(cf.RandomInt(1, 365*24*60)), // Random time in past year
		UpdatedAt: time.Now(),
	}
}

// ElectronicsBusiness creates an electronics business
func (cf *CommerceFixtures) ElectronicsBusiness(country string) *models.Business {
	business := cf.BasicBusiness(country)
	business.Category = "electronics"
	business.Name = fmt.Sprintf("Elite Electronics %s", country)
	business.Description = "Premium electronics retailer with latest technology"
	return business
}

// FashionBusiness creates a fashion business
func (cf *CommerceFixtures) FashionBusiness(country string) *models.Business {
	business := cf.BasicBusiness(country)
	business.Category = "fashion"
	business.Name = fmt.Sprintf("Fashion Forward %s", country)
	business.Description = "Trendy fashion retailer with latest styles"
	return business
}

// BasicProduct creates a basic product for testing
func (cf *CommerceFixtures) BasicProduct(businessID uuid.UUID, country string) *models.Product {
	category := SEABusinessCategories[cf.RandomInt(0, len(SEABusinessCategories))]
	productNames := []string{"Premium Product", "Quality Item", "Best Seller"}
	name := productNames[cf.RandomInt(0, len(productNames))]

	return &models.Product{
		ID:          cf.UUID(fmt.Sprintf("product-%s-%s", country, name)),
		BusinessID:  businessID,
		Name:        name,
		Description: fmt.Sprintf("High-quality %s product from %s", category, country),
		Category:    category,
		Brand:       fmt.Sprintf("%s Brand", country),
		CreatedAt:   cf.PastTime(cf.RandomInt(1, 180*24*60)), // Random time in past 6 months
		UpdatedAt:   time.Now(),
	}
}

// ElectronicsProduct creates an electronics product
func (cf *CommerceFixtures) ElectronicsProduct(businessID uuid.UUID, country string) *models.Product {
	product := cf.BasicProduct(businessID, country)
	product.Category = "electronics"
	product.Name = fmt.Sprintf("Smart Electronics Device %s", country)
	product.Description = "Advanced electronics with latest technology features"
	return product
}

// FashionProduct creates a fashion product
func (cf *CommerceFixtures) FashionProduct(businessID uuid.UUID, country string) *models.Product {
	product := cf.BasicProduct(businessID, country)
	product.Category = "fashion"
	product.Name = fmt.Sprintf("Stylish Fashion Item %s", country)
	product.Description = "Trendy fashion item with premium quality materials"
	return product
}

// FoodProduct creates a food product
func (cf *CommerceFixtures) FoodProduct(businessID uuid.UUID, country string) *models.Product {
	product := cf.BasicProduct(businessID, country)
	product.Category = "food"
	product.Name = fmt.Sprintf("Gourmet Food %s", country)
	product.Description = "Authentic local cuisine ingredients and specialties"
	return product
}

// Helper methods for location data
func (cf *CommerceFixtures) getCity(country string) string {
	cities := map[string][]string{
		"TH": {"Bangkok", "Chiang Mai", "Phuket", "Pattaya"},
		"SG": {"Singapore", "Jurong", "Tampines", "Woodlands"},
		"ID": {"Jakarta", "Surabaya", "Bandung", "Medan"},
		"MY": {"Kuala Lumpur", "George Town", "Johor Bahru", "Ipoh"},
		"VN": {"Ho Chi Minh City", "Hanoi", "Da Nang", "Can Tho"},
		"PH": {"Manila", "Cebu City", "Davao", "Quezon City"},
	}

	if cityList, ok := cities[country]; ok {
		return cityList[cf.RandomInt(0, len(cityList))]
	}
	return "Metro City"
}

func (cf *CommerceFixtures) getState(country string) string {
	states := map[string][]string{
		"TH": {"Bangkok", "Chiang Mai", "Phuket", "Khon Kaen"},
		"SG": {"Central", "East", "North", "West"},
		"ID": {"Jakarta", "West Java", "East Java", "Bali"},
		"MY": {"Kuala Lumpur", "Selangor", "Penang", "Johor"},
		"VN": {"Ho Chi Minh", "Hanoi", "Da Nang", "An Giang"},
		"PH": {"Metro Manila", "Cebu", "Davao", "Iloilo"},
	}

	if stateList, ok := states[country]; ok {
		return stateList[cf.RandomInt(0, len(stateList))]
	}
	return "Central Region"
}

func (cf *CommerceFixtures) getPostalCode(country string) string {
	switch country {
	case "TH":
		return fmt.Sprintf("%05d", cf.RandomInt(10000, 99999))
	case "SG":
		return fmt.Sprintf("%06d", cf.RandomInt(100000, 999999))
	case "ID":
		return fmt.Sprintf("%05d", cf.RandomInt(10000, 99999))
	case "MY":
		return fmt.Sprintf("%05d", cf.RandomInt(10000, 99999))
	case "VN":
		return fmt.Sprintf("%06d", cf.RandomInt(100000, 999999))
	case "PH":
		return fmt.Sprintf("%04d", cf.RandomInt(1000, 9999))
	default:
		return "12345"
	}
}

func (cf *CommerceFixtures) getLanguage(country string) string {
	languages := map[string]string{
		"TH": "th",
		"SG": "zh",
		"ID": "id",
		"MY": "ms",
		"VN": "vi",
		"PH": "tl",
	}

	if lang, ok := languages[country]; ok {
		return lang
	}
	return "en"
}