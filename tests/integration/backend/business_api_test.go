package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// BusinessAPITestSuite provides comprehensive integration testing for business API endpoints
type BusinessAPITestSuite struct {
	suite.Suite
	baseURL          string
	httpClient       *http.Client
	gatewayPort      int
	testUserID       string
	createdBusinesses []string
}

// Business represents business structure
type Business struct {
	ID             string                 `json:"id"`
	OwnerID        string                 `json:"ownerId"`
	Name           string                 `json:"name"`
	LegalName      *string                `json:"legalName,omitempty"`
	Description    *string                `json:"description,omitempty"`
	BusinessType   string                 `json:"businessType"`
	Industry       string                 `json:"industry"`
	TaxID          *string                `json:"taxId,omitempty"`
	RegistrationNo *string                `json:"registrationNumber,omitempty"`
	Email          string                 `json:"email"`
	Phone          *string                `json:"phone,omitempty"`
	Website        *string                `json:"website,omitempty"`
	Address        BusinessAddress        `json:"address"`
	BillingAddress *BusinessAddress       `json:"billingAddress,omitempty"`
	ContactPerson  *ContactPerson         `json:"contactPerson,omitempty"`
	BankAccount    *BankAccount           `json:"bankAccount,omitempty"`
	Verification   BusinessVerification   `json:"verification"`
	Settings       BusinessSettings       `json:"settings"`
	Status         string                 `json:"status"`
	Metrics        BusinessMetrics        `json:"metrics"`
	Documents      []BusinessDocument     `json:"documents"`
	Licenses       []BusinessLicense      `json:"licenses"`
	Attributes     map[string]interface{} `json:"attributes,omitempty"`
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt"`
	VerifiedAt     *time.Time             `json:"verifiedAt,omitempty"`
}

// BusinessAddress represents business address
type BusinessAddress struct {
	Street     string  `json:"street"`
	Street2    *string `json:"street2,omitempty"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	PostalCode string  `json:"postalCode"`
	Country    string  `json:"country"`
	Latitude   *float64 `json:"latitude,omitempty"`
	Longitude  *float64 `json:"longitude,omitempty"`
}

// ContactPerson represents business contact person
type ContactPerson struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Title     *string `json:"title,omitempty"`
	Email     string  `json:"email"`
	Phone     *string `json:"phone,omitempty"`
}

// BankAccount represents business bank account
type BankAccount struct {
	BankName      string `json:"bankName"`
	AccountNumber string `json:"accountNumber"`
	RoutingNumber string `json:"routingNumber"`
	AccountType   string `json:"accountType"`
	Currency      string `json:"currency"`
	IsVerified    bool   `json:"isVerified"`
}

// BusinessVerification represents verification status
type BusinessVerification struct {
	Status            string     `json:"status"`
	Level             string     `json:"level"`
	SubmittedAt       *time.Time `json:"submittedAt,omitempty"`
	ReviewedAt        *time.Time `json:"reviewedAt,omitempty"`
	VerifiedAt        *time.Time `json:"verifiedAt,omitempty"`
	RejectedAt        *time.Time `json:"rejectedAt,omitempty"`
	RejectionReason   *string    `json:"rejectionReason,omitempty"`
	RequiredDocuments []string   `json:"requiredDocuments"`
	SubmittedDocs     []string   `json:"submittedDocuments"`
	ReviewNotes       *string    `json:"reviewNotes,omitempty"`
}

// BusinessSettings represents business settings
type BusinessSettings struct {
	Timezone          string                 `json:"timezone"`
	Currency          string                 `json:"currency"`
	Language          string                 `json:"language"`
	TaxSettings       TaxSettings            `json:"taxSettings"`
	PaymentSettings   PaymentSettings        `json:"paymentSettings"`
	ShippingSettings  ShippingSettings       `json:"shippingSettings"`
	NotificationPrefs NotificationPrefs      `json:"notificationPrefs"`
	CustomFields      map[string]interface{} `json:"customFields,omitempty"`
}

// TaxSettings represents tax configuration
type TaxSettings struct {
	TaxEnabled       bool                    `json:"taxEnabled"`
	DefaultTaxRate   float64                 `json:"defaultTaxRate"`
	TaxIncluded      bool                    `json:"taxIncluded"`
	TaxExemptNumbers []string                `json:"taxExemptNumbers"`
	TaxZones         []TaxZone               `json:"taxZones"`
	TaxCategories    map[string]float64      `json:"taxCategories"`
}

// TaxZone represents tax zone configuration
type TaxZone struct {
	Name     string   `json:"name"`
	Countries []string `json:"countries"`
	States   []string `json:"states,omitempty"`
	Rate     float64  `json:"rate"`
}

// PaymentSettings represents payment configuration
type PaymentSettings struct {
	AcceptedMethods   []string               `json:"acceptedMethods"`
	DefaultCurrency   string                 `json:"defaultCurrency"`
	PaymentProviders  []PaymentProvider      `json:"paymentProviders"`
	AutoCapture       bool                   `json:"autoCapture"`
	PaymentTerms      *string                `json:"paymentTerms,omitempty"`
	MinOrderAmount    *float64               `json:"minOrderAmount,omitempty"`
	MaxOrderAmount    *float64               `json:"maxOrderAmount,omitempty"`
}

// PaymentProvider represents payment provider configuration
type PaymentProvider struct {
	Name       string                 `json:"name"`
	Enabled    bool                   `json:"enabled"`
	Config     map[string]interface{} `json:"config"`
	TestMode   bool                   `json:"testMode"`
}

// ShippingSettings represents shipping configuration
type ShippingSettings struct {
	DefaultProvider  string           `json:"defaultProvider"`
	ShippingZones    []ShippingZone   `json:"shippingZones"`
	FreeShipping     *FreeShipping    `json:"freeShipping,omitempty"`
	HandlingTime     *string          `json:"handlingTime,omitempty"`
	PackagingOptions []PackageOption  `json:"packagingOptions"`
}

// ShippingZone represents shipping zone configuration
type ShippingZone struct {
	Name      string           `json:"name"`
	Countries []string         `json:"countries"`
	Methods   []ShippingMethod `json:"methods"`
}

// ShippingMethod represents shipping method
type ShippingMethod struct {
	Name        string   `json:"name"`
	Cost        float64  `json:"cost"`
	MinDelivery *int     `json:"minDeliveryDays,omitempty"`
	MaxDelivery *int     `json:"maxDeliveryDays,omitempty"`
	Enabled     bool     `json:"enabled"`
}

// FreeShipping represents free shipping configuration
type FreeShipping struct {
	Enabled      bool     `json:"enabled"`
	MinAmount    float64  `json:"minAmount"`
	Countries    []string `json:"countries,omitempty"`
}

// PackageOption represents packaging option
type PackageOption struct {
	Name   string  `json:"name"`
	Cost   float64 `json:"cost"`
	Weight float64 `json:"maxWeight"`
}

// NotificationPrefs represents notification preferences
type NotificationPrefs struct {
	OrderNotifications    bool `json:"orderNotifications"`
	PaymentNotifications  bool `json:"paymentNotifications"`
	ShippingNotifications bool `json:"shippingNotifications"`
	ReviewNotifications   bool `json:"reviewNotifications"`
	MarketingEmails       bool `json:"marketingEmails"`
	EmailDigest           bool `json:"emailDigest"`
	SMSNotifications      bool `json:"smsNotifications"`
}

// BusinessMetrics represents business metrics
type BusinessMetrics struct {
	TotalRevenue    float64   `json:"totalRevenue"`
	MonthlyRevenue  float64   `json:"monthlyRevenue"`
	TotalOrders     int       `json:"totalOrders"`
	MonthlyOrders   int       `json:"monthlyOrders"`
	TotalProducts   int       `json:"totalProducts"`
	ActiveProducts  int       `json:"activeProducts"`
	CustomerCount   int       `json:"customerCount"`
	AverageRating   float64   `json:"averageRating"`
	ReviewCount     int       `json:"reviewCount"`
	ConversionRate  float64   `json:"conversionRate"`
	LastOrderDate   *time.Time `json:"lastOrderDate,omitempty"`
	LastUpdated     time.Time  `json:"lastUpdated"`
}

// BusinessDocument represents business document
type BusinessDocument struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	URL          string    `json:"url"`
	MimeType     string    `json:"mimeType"`
	Size         int64     `json:"size"`
	Status       string    `json:"status"`
	UploadedAt   time.Time `json:"uploadedAt"`
	VerifiedAt   *time.Time `json:"verifiedAt,omitempty"`
	RejectedAt   *time.Time `json:"rejectedAt,omitempty"`
	RejectionReason *string `json:"rejectionReason,omitempty"`
}

// BusinessLicense represents business license
type BusinessLicense struct {
	ID           string     `json:"id"`
	Type         string     `json:"type"`
	Number       string     `json:"number"`
	IssuedBy     string     `json:"issuedBy"`
	IssuedDate   time.Time  `json:"issuedDate"`
	ExpiryDate   *time.Time `json:"expiryDate,omitempty"`
	Status       string     `json:"status"`
	DocumentURL  *string    `json:"documentUrl,omitempty"`
}

// CreateBusinessRequest represents request to create business
type CreateBusinessRequest struct {
	Name           string                 `json:"name"`
	LegalName      *string                `json:"legalName,omitempty"`
	Description    *string                `json:"description,omitempty"`
	BusinessType   string                 `json:"businessType"`
	Industry       string                 `json:"industry"`
	TaxID          *string                `json:"taxId,omitempty"`
	RegistrationNo *string                `json:"registrationNumber,omitempty"`
	Email          string                 `json:"email"`
	Phone          *string                `json:"phone,omitempty"`
	Website        *string                `json:"website,omitempty"`
	Address        BusinessAddress        `json:"address"`
	BillingAddress *BusinessAddress       `json:"billingAddress,omitempty"`
	ContactPerson  *ContactPerson         `json:"contactPerson,omitempty"`
	BankAccount    *BankAccount           `json:"bankAccount,omitempty"`
	Settings       *BusinessSettings      `json:"settings,omitempty"`
	Attributes     map[string]interface{} `json:"attributes,omitempty"`
}

// UpdateBusinessRequest represents request to update business
type UpdateBusinessRequest struct {
	Name           *string                `json:"name,omitempty"`
	LegalName      *string                `json:"legalName,omitempty"`
	Description    *string                `json:"description,omitempty"`
	BusinessType   *string                `json:"businessType,omitempty"`
	Industry       *string                `json:"industry,omitempty"`
	TaxID          *string                `json:"taxId,omitempty"`
	RegistrationNo *string                `json:"registrationNumber,omitempty"`
	Email          *string                `json:"email,omitempty"`
	Phone          *string                `json:"phone,omitempty"`
	Website        *string                `json:"website,omitempty"`
	Address        *BusinessAddress       `json:"address,omitempty"`
	BillingAddress *BusinessAddress       `json:"billingAddress,omitempty"`
	ContactPerson  *ContactPerson         `json:"contactPerson,omitempty"`
	BankAccount    *BankAccount           `json:"bankAccount,omitempty"`
	Settings       *BusinessSettings      `json:"settings,omitempty"`
	Attributes     map[string]interface{} `json:"attributes,omitempty"`
}

// BusinessResponse represents API response for business operations
type BusinessResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Business  *Business   `json:"business,omitempty"`
	Businesses []Business `json:"businesses,omitempty"`
	Total     int         `json:"total,omitempty"`
	Error     *string     `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// SetupSuite initializes the test suite
func (suite *BusinessAPITestSuite) SetupSuite() {
	suite.baseURL = "http://localhost"
	suite.gatewayPort = 8080
	suite.httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
	suite.testUserID = uuid.New().String()
	suite.createdBusinesses = []string{}

	// Wait for services
	suite.waitForServices()
}

// TearDownSuite cleans up after test suite
func (suite *BusinessAPITestSuite) TearDownSuite() {
	suite.cleanupTestData()
}

// TestBusinessRegistration tests complete business registration flow
func (suite *BusinessAPITestSuite) TestBusinessRegistration() {
	// 1. Create business registration
	createReq := CreateBusinessRequest{
		Name:         "TechCorp Solutions Ltd.",
		LegalName:    stringPtr("TechCorp Solutions Limited"),
		Description:  stringPtr("Technology solutions and consulting services"),
		BusinessType: "corporation",
		Industry:     "technology",
		TaxID:        stringPtr("123456789"),
		RegistrationNo: stringPtr("REG-123456"),
		Email:        "contact@techcorp.com",
		Phone:        stringPtr("+1-555-0123"),
		Website:      stringPtr("https://techcorp.com"),
		Address: BusinessAddress{
			Street:     "123 Tech Street",
			Street2:    stringPtr("Suite 456"),
			City:       "San Francisco",
			State:      "CA",
			PostalCode: "94105",
			Country:    "US",
			Latitude:   float64Ptr(37.7749),
			Longitude:  float64Ptr(-122.4194),
		},
		BillingAddress: &BusinessAddress{
			Street:     "789 Billing Ave",
			City:       "San Francisco",
			State:      "CA",
			PostalCode: "94107",
			Country:    "US",
		},
		ContactPerson: &ContactPerson{
			FirstName: "John",
			LastName:  "Doe",
			Title:     stringPtr("CEO"),
			Email:     "john.doe@techcorp.com",
			Phone:     stringPtr("+1-555-0124"),
		},
		BankAccount: &BankAccount{
			BankName:      "First National Bank",
			AccountNumber: "****1234",
			RoutingNumber: "123456789",
			AccountType:   "business_checking",
			Currency:      "USD",
			IsVerified:    false,
		},
		Settings: &BusinessSettings{
			Timezone: "America/Los_Angeles",
			Currency: "USD",
			Language: "en",
			TaxSettings: TaxSettings{
				TaxEnabled:     true,
				DefaultTaxRate: 8.75,
				TaxIncluded:    false,
				TaxZones: []TaxZone{
					{
						Name:      "California",
						Countries: []string{"US"},
						States:    []string{"CA"},
						Rate:      8.75,
					},
				},
			},
			PaymentSettings: PaymentSettings{
				AcceptedMethods: []string{"credit_card", "paypal", "bank_transfer"},
				DefaultCurrency: "USD",
				AutoCapture:     true,
				MinOrderAmount:  float64Ptr(10.00),
				MaxOrderAmount:  float64Ptr(10000.00),
			},
			ShippingSettings: ShippingSettings{
				DefaultProvider: "ups",
				ShippingZones: []ShippingZone{
					{
						Name:      "Domestic",
						Countries: []string{"US"},
						Methods: []ShippingMethod{
							{
								Name:        "Standard",
								Cost:        9.99,
								MinDelivery: intPtr(3),
								MaxDelivery: intPtr(7),
								Enabled:     true,
							},
							{
								Name:        "Express",
								Cost:        19.99,
								MinDelivery: intPtr(1),
								MaxDelivery: intPtr(3),
								Enabled:     true,
							},
						},
					},
				},
				FreeShipping: &FreeShipping{
					Enabled:   true,
					MinAmount: 50.00,
					Countries: []string{"US"},
				},
			},
			NotificationPrefs: NotificationPrefs{
				OrderNotifications:    true,
				PaymentNotifications:  true,
				ShippingNotifications: true,
				ReviewNotifications:   true,
				MarketingEmails:       false,
				EmailDigest:           true,
				SMSNotifications:      false,
			},
		},
		Attributes: map[string]interface{}{
			"founded_year":    2020,
			"employee_count":  "10-50",
			"annual_revenue":  "1M-5M",
			"certifications":  []string{"ISO9001", "SOC2"},
		},
	}

	business := suite.createBusiness(createReq)
	businessID := business.ID
	suite.createdBusinesses = append(suite.createdBusinesses, businessID)

	// Verify creation
	assert.Equal(suite.T(), createReq.Name, business.Name)
	assert.Equal(suite.T(), createReq.Email, business.Email)
	assert.Equal(suite.T(), createReq.BusinessType, business.BusinessType)
	assert.Equal(suite.T(), createReq.Industry, business.Industry)
	assert.Equal(suite.T(), "pending", business.Status)
	assert.Equal(suite.T(), "unverified", business.Verification.Status)
	assert.NotNil(suite.T(), business.ContactPerson)
	assert.NotNil(suite.T(), business.BankAccount)
	assert.False(suite.T(), business.BankAccount.IsVerified)

	// 2. Get business details
	retrievedBusiness := suite.getBusiness(businessID)
	assert.Equal(suite.T(), business.ID, retrievedBusiness.ID)
	assert.Equal(suite.T(), business.Name, retrievedBusiness.Name)

	// 3. Update business information
	updateReq := UpdateBusinessRequest{
		Description: stringPtr("Updated: Advanced technology solutions and consulting services"),
		Phone:       stringPtr("+1-555-0125"),
		Website:     stringPtr("https://www.techcorp.com"),
		Settings: &BusinessSettings{
			Timezone: "America/Los_Angeles",
			Currency: "USD",
			Language: "en",
			TaxSettings: TaxSettings{
				TaxEnabled:     true,
				DefaultTaxRate: 9.25, // Updated tax rate
				TaxIncluded:    false,
			},
			PaymentSettings: PaymentSettings{
				AcceptedMethods: []string{"credit_card", "paypal", "bank_transfer", "crypto"},
				DefaultCurrency: "USD",
				AutoCapture:     true,
			},
			NotificationPrefs: NotificationPrefs{
				OrderNotifications:    true,
				PaymentNotifications:  true,
				ShippingNotifications: true,
				ReviewNotifications:   true,
				MarketingEmails:       true, // Updated preference
				EmailDigest:           true,
				SMSNotifications:      true, // Updated preference
			},
		},
		Attributes: map[string]interface{}{
			"founded_year":    2020,
			"employee_count":  "50-100", // Updated
			"annual_revenue":  "5M-10M", // Updated
			"certifications":  []string{"ISO9001", "SOC2", "PCI-DSS"}, // Added certification
			"headquarters":    "San Francisco, CA",
		},
	}

	updatedBusiness := suite.updateBusiness(businessID, updateReq)
	assert.Contains(suite.T(), *updatedBusiness.Description, "Updated:")
	assert.Equal(suite.T(), "+1-555-0125", *updatedBusiness.Phone)
	assert.Equal(suite.T(), 9.25, updatedBusiness.Settings.TaxSettings.DefaultTaxRate)
	assert.True(suite.T(), updatedBusiness.Settings.NotificationPrefs.MarketingEmails)
	assert.True(suite.T(), updatedBusiness.Settings.NotificationPrefs.SMSNotifications)

	// 4. Test business verification flow
	suite.testBusinessVerification(businessID)

	// 5. Test business metrics
	metrics := suite.getBusinessMetrics(businessID)
	assert.Equal(suite.T(), businessID, metrics.Business.ID)
	assert.GreaterOrEqual(suite.T(), metrics.TotalRevenue, 0.0)
	assert.GreaterOrEqual(suite.T(), metrics.TotalOrders, 0)
}

// TestBusinessVerification tests business verification process
func (suite *BusinessAPITestSuite) TestBusinessVerification() {
	// Create business for verification testing
	createReq := CreateBusinessRequest{
		Name:         "Verification Test Business",
		BusinessType: "llc",
		Industry:     "retail",
		Email:        "verify@testbusiness.com",
		Address: BusinessAddress{
			Street:     "456 Verify Street",
			City:       "Los Angeles",
			State:      "CA",
			PostalCode: "90210",
			Country:    "US",
		},
	}

	business := suite.createBusiness(createReq)
	businessID := business.ID
	suite.createdBusinesses = append(suite.createdBusinesses, businessID)

	// Test verification workflow
	suite.testBusinessVerification(businessID)
}

// TestBusinessValidation tests business validation rules
func (suite *BusinessAPITestSuite) TestBusinessValidation() {
	testCases := []struct {
		name        string
		createReq   CreateBusinessRequest
		expectError bool
		statusCode  int
	}{
		{
			name: "Missing required name",
			createReq: CreateBusinessRequest{
				BusinessType: "llc",
				Industry:     "technology",
				Email:        "test@example.com",
				Address: BusinessAddress{
					Street:     "123 Test St",
					City:       "Test City",
					State:      "CA",
					PostalCode: "12345",
					Country:    "US",
				},
			},
			expectError: true,
			statusCode:  http.StatusBadRequest,
		},
		{
			name: "Invalid business type",
			createReq: CreateBusinessRequest{
				Name:         "Invalid Type Business",
				BusinessType: "invalid_type",
				Industry:     "technology",
				Email:        "test@example.com",
				Address: BusinessAddress{
					Street:     "123 Test St",
					City:       "Test City",
					State:      "CA",
					PostalCode: "12345",
					Country:    "US",
				},
			},
			expectError: true,
			statusCode:  http.StatusBadRequest,
		},
		{
			name: "Invalid email format",
			createReq: CreateBusinessRequest{
				Name:         "Invalid Email Business",
				BusinessType: "llc",
				Industry:     "technology",
				Email:        "invalid-email-format",
				Address: BusinessAddress{
					Street:     "123 Test St",
					City:       "Test City",
					State:      "CA",
					PostalCode: "12345",
					Country:    "US",
				},
			},
			expectError: true,
			statusCode:  http.StatusBadRequest,
		},
		{
			name: "Duplicate tax ID",
			createReq: CreateBusinessRequest{
				Name:         "Duplicate Tax ID Business",
				BusinessType: "corporation",
				Industry:     "technology",
				Email:        "duplicate@example.com",
				TaxID:        stringPtr("DUPLICATE-TAX-123"),
				Address: BusinessAddress{
					Street:     "123 Test St",
					City:       "Test City",
					State:      "CA",
					PostalCode: "12345",
					Country:    "US",
				},
			},
			expectError: true,
			statusCode:  http.StatusConflict,
		},
		{
			name: "Missing required address",
			createReq: CreateBusinessRequest{
				Name:         "Missing Address Business",
				BusinessType: "sole_proprietorship",
				Industry:     "consulting",
				Email:        "missing@example.com",
			},
			expectError: true,
			statusCode:  http.StatusBadRequest,
		},
	}

	// First create a business with specific tax ID for duplicate test
	firstBusiness := CreateBusinessRequest{
		Name:         "First Business",
		BusinessType: "llc",
		Industry:     "technology",
		Email:        "first@example.com",
		TaxID:        stringPtr("DUPLICATE-TAX-123"),
		Address: BusinessAddress{
			Street:     "123 First St",
			City:       "Test City",
			State:      "CA",
			PostalCode: "12345",
			Country:    "US",
		},
	}
	business := suite.createBusiness(firstBusiness)
	suite.createdBusinesses = append(suite.createdBusinesses, business.ID)

	// Run validation tests
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			if tc.expectError {
				suite.expectBusinessCreationError(tc.statusCode, tc.createReq)
			} else {
				bus := suite.createBusiness(tc.createReq)
				suite.createdBusinesses = append(suite.createdBusinesses, bus.ID)
				assert.NotEmpty(t, bus.ID)
			}
		})
	}
}

// TestBusinessSearch tests business search functionality
func (suite *BusinessAPITestSuite) TestBusinessSearch() {
	// Create test businesses
	businesses := []CreateBusinessRequest{
		{
			Name:         "Tech Innovations Inc",
			BusinessType: "corporation",
			Industry:     "technology",
			Email:        "contact@techinnovations.com",
			Address: BusinessAddress{
				Street: "123 Innovation St", City: "San Francisco", State: "CA", PostalCode: "94105", Country: "US",
			},
		},
		{
			Name:         "Green Solutions LLC",
			BusinessType: "llc",
			Industry:     "environmental",
			Email:        "info@greensolutions.com",
			Address: BusinessAddress{
				Street: "456 Green Ave", City: "Portland", State: "OR", PostalCode: "97201", Country: "US",
			},
		},
		{
			Name:         "Creative Tech Studio",
			BusinessType: "llc",
			Industry:     "technology",
			Email:        "hello@creativetech.com",
			Address: BusinessAddress{
				Street: "789 Creative Blvd", City: "Austin", State: "TX", PostalCode: "73301", Country: "US",
			},
		},
	}

	createdIDs := []string{}
	for _, req := range businesses {
		bus := suite.createBusiness(req)
		createdIDs = append(createdIDs, bus.ID)
		suite.createdBusinesses = append(suite.createdBusinesses, bus.ID)
	}

	// Test search scenarios
	testCases := []struct {
		name          string
		query         string
		industry      string
		businessType  string
		expectedCount int
	}{
		{
			name:          "Search by name",
			query:         "Tech",
			expectedCount: 2,
		},
		{
			name:          "Search by industry",
			industry:      "technology",
			expectedCount: 2,
		},
		{
			name:          "Search by business type",
			businessType:  "llc",
			expectedCount: 2,
		},
		{
			name:          "Search by query and industry",
			query:         "Tech",
			industry:      "technology",
			expectedCount: 2,
		},
		{
			name:          "Search by specific name",
			query:         "Green Solutions",
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			results := suite.searchBusinesses(tc.query, tc.industry, tc.businessType)
			assert.GreaterOrEqual(t, len(results.Businesses), tc.expectedCount)
		})
	}
}

// TestBusinessConcurrency tests concurrent business operations
func (suite *BusinessAPITestSuite) TestBusinessConcurrency() {
	// Create a business for concurrent updates
	createReq := CreateBusinessRequest{
		Name:         "Concurrency Test Business",
		BusinessType: "llc",
		Industry:     "technology",
		Email:        "concurrency@test.com",
		Address: BusinessAddress{
			Street:     "123 Concurrency St",
			City:       "Test City",
			State:      "CA",
			PostalCode: "12345",
			Country:    "US",
		},
	}

	business := suite.createBusiness(createReq)
	businessID := business.ID
	suite.createdBusinesses = append(suite.createdBusinesses, businessID)

	// Perform concurrent updates
	concurrency := 5
	resultChan := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(index int) {
			updateReq := UpdateBusinessRequest{
				Description: stringPtr(fmt.Sprintf("Updated description %d", index)),
				Attributes: map[string]interface{}{
					"update_index": index,
					"timestamp":    time.Now().Unix(),
				},
			}

			err := suite.updateBusinessWithError(businessID, updateReq)
			resultChan <- err == nil
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < concurrency; i++ {
		if <-resultChan {
			successCount++
		}
	}

	// At least some should succeed
	assert.GreaterOrEqual(suite.T(), successCount, 1)

	// Final business should be consistent
	finalBusiness := suite.getBusiness(businessID)
	assert.NotNil(suite.T(), finalBusiness.Description)
	assert.Contains(suite.T(), *finalBusiness.Description, "Updated description")
}

// Helper methods

func (suite *BusinessAPITestSuite) createBusiness(req CreateBusinessRequest) *Business {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/businesses", suite.baseURL, suite.gatewayPort)

	body, err := json.Marshal(req)
	require.NoError(suite.T(), err)

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	require.NoError(suite.T(), err)
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(httpReq)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var businessResp BusinessResponse
	err = json.NewDecoder(resp.Body).Decode(&businessResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), businessResp.Success)
	require.NotNil(suite.T(), businessResp.Business)

	return businessResp.Business
}

func (suite *BusinessAPITestSuite) getBusiness(businessID string) *Business {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/businesses/%s", suite.baseURL, suite.gatewayPort, businessID)

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(suite.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var businessResp BusinessResponse
	err = json.NewDecoder(resp.Body).Decode(&businessResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), businessResp.Success)
	require.NotNil(suite.T(), businessResp.Business)

	return businessResp.Business
}

func (suite *BusinessAPITestSuite) updateBusiness(businessID string, req UpdateBusinessRequest) *Business {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/businesses/%s", suite.baseURL, suite.gatewayPort, businessID)

	body, err := json.Marshal(req)
	require.NoError(suite.T(), err)

	httpReq, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	require.NoError(suite.T(), err)
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(httpReq)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var businessResp BusinessResponse
	err = json.NewDecoder(resp.Body).Decode(&businessResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), businessResp.Success)
	require.NotNil(suite.T(), businessResp.Business)

	return businessResp.Business
}

func (suite *BusinessAPITestSuite) updateBusinessWithError(businessID string, req UpdateBusinessRequest) error {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/businesses/%s", suite.baseURL, suite.gatewayPort, businessID)

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	return nil
}

func (suite *BusinessAPITestSuite) deleteBusiness(businessID string) {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/businesses/%s", suite.baseURL, suite.gatewayPort, businessID)

	req, err := http.NewRequest("DELETE", url, nil)
	require.NoError(suite.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *BusinessAPITestSuite) searchBusinesses(query, industry, businessType string) *BusinessResponse {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/businesses/search", suite.baseURL, suite.gatewayPort)

	params := make(map[string]string)
	if query != "" {
		params["query"] = query
	}
	if industry != "" {
		params["industry"] = industry
	}
	if businessType != "" {
		params["businessType"] = businessType
	}

	if len(params) > 0 {
		queryParams := ""
		for key, value := range params {
			if queryParams != "" {
				queryParams += "&"
			}
			queryParams += fmt.Sprintf("%s=%s", key, value)
		}
		url += "?" + queryParams
	}

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(suite.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var businessResp BusinessResponse
	err = json.NewDecoder(resp.Body).Decode(&businessResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), businessResp.Success)

	return &businessResp
}

func (suite *BusinessAPITestSuite) getBusinessMetrics(businessID string) *Business {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/businesses/%s/metrics", suite.baseURL, suite.gatewayPort, businessID)

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(suite.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var businessResp BusinessResponse
	err = json.NewDecoder(resp.Body).Decode(&businessResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), businessResp.Success)
	require.NotNil(suite.T(), businessResp.Business)

	return businessResp.Business
}

func (suite *BusinessAPITestSuite) testBusinessVerification(businessID string) {
	// Submit for verification
	url := fmt.Sprintf("%s:%d/api/v1/commerce/businesses/%s/verification/submit", suite.baseURL, suite.gatewayPort, businessID)

	req, err := http.NewRequest("POST", url, nil)
	require.NoError(suite.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Check verification status
	business := suite.getBusiness(businessID)
	assert.Equal(suite.T(), "submitted", business.Verification.Status)
	assert.NotNil(suite.T(), business.Verification.SubmittedAt)
}

func (suite *BusinessAPITestSuite) expectBusinessCreationError(expectedStatus int, req CreateBusinessRequest) {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/businesses", suite.baseURL, suite.gatewayPort)

	body, err := json.Marshal(req)
	require.NoError(suite.T(), err)

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	require.NoError(suite.T(), err)
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(httpReq)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), expectedStatus, resp.StatusCode)
}

func (suite *BusinessAPITestSuite) waitForServices() {
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		url := fmt.Sprintf("%s:%d/health", suite.baseURL, suite.gatewayPort)
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(1 * time.Second)
	}
	suite.T().Fatal("Services not ready after waiting")
}

func (suite *BusinessAPITestSuite) cleanupTestData() {
	// Clean up created businesses
	for _, businessID := range suite.createdBusinesses {
		suite.deleteBusiness(businessID)
	}
}

// Utility functions
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}

// TestBusinessAPI runs the business API test suite
func TestBusinessAPI(t *testing.T) {
	suite.Run(t, new(BusinessAPITestSuite))
}