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

// Shop represents a merchant shop/store in the commerce system
type Shop struct {
	ID              uuid.UUID       `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OwnerID         uuid.UUID       `json:"owner_id" gorm:"type:varchar(36);not null;index"`
	Name            string          `json:"name" gorm:"type:varchar(255);not null"`
	Description     *string         `json:"description,omitempty" gorm:"type:text"`
	Slug            string          `json:"slug" gorm:"type:varchar(100);uniqueIndex;not null"`
	Category        ShopCategory    `json:"category" gorm:"type:varchar(50);not null"`
	Email           string          `json:"email" gorm:"type:varchar(255);not null"`
	Phone           *string         `json:"phone,omitempty" gorm:"type:varchar(20)"`
	Website         *string         `json:"website,omitempty" gorm:"type:varchar(255)"`
	LogoURL         *string         `json:"logo_url,omitempty" gorm:"type:varchar(512)"`
	BannerURL       *string         `json:"banner_url,omitempty" gorm:"type:varchar(512)"`
	Status          ShopStatus      `json:"status" gorm:"type:varchar(20);default:'pending'"`
	Visibility      ShopVisibility  `json:"visibility" gorm:"type:varchar(20);default:'public'"`
	Country         string          `json:"country" gorm:"type:varchar(2);not null"` // ISO country code
	Currency        Currency        `json:"currency" gorm:"type:varchar(3);not null"`
	Language        string          `json:"language" gorm:"type:varchar(5);not null"` // ISO language code
	Timezone        string          `json:"timezone" gorm:"type:varchar(50);not null"` // IANA timezone
	Address         ShopAddress     `json:"address" gorm:"type:json"`
	BusinessInfo    BusinessInfo    `json:"business_info" gorm:"type:json"`
	Settings        ShopSettings    `json:"settings" gorm:"type:json"`
	Features        ShopFeatures    `json:"features" gorm:"type:json"`
	Subscription    Subscription    `json:"subscription" gorm:"type:json"`
	Analytics       ShopAnalytics   `json:"analytics" gorm:"type:json"`
	Compliance      ComplianceInfo  `json:"compliance" gorm:"type:json"`
	SocialMedia     SocialMedia     `json:"social_media" gorm:"type:json"`
	Metadata        ShopMetadata    `json:"metadata" gorm:"type:json"`
	Tags            StringSlice     `json:"tags" gorm:"type:json"`
	IsVerified      bool            `json:"is_verified" gorm:"default:false"`
	IsFeatured      bool            `json:"is_featured" gorm:"default:false"`
	VerifiedAt      *time.Time      `json:"verified_at,omitempty"`
	FeaturedAt      *time.Time      `json:"featured_at,omitempty"`
	LastActivityAt  *time.Time      `json:"last_activity_at,omitempty"`
	CreatedAt       time.Time       `json:"created_at" gorm:"not null"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"not null"`
}

// ShopCategory represents shop business categories
type ShopCategory string

const (
	CategoryElectronicsShop     ShopCategory = "electronics"
	CategoryFashionShop         ShopCategory = "fashion"
	CategoryHomeGardenShop      ShopCategory = "home_garden"
	CategoryBeautyShop          ShopCategory = "beauty"
	CategorySportsShop          ShopCategory = "sports"
	CategoryBooksShop           ShopCategory = "books"
	CategoryFoodBeverageShop    ShopCategory = "food_beverage"
	CategoryHealthShop          ShopCategory = "health"
	CategoryToysShop            ShopCategory = "toys"
	CategoryAutomotiveShop      ShopCategory = "automotive"
	CategoryServicesShop        ShopCategory = "services"
	CategoryDigitalShop         ShopCategory = "digital"
	CategoryHandmadeShop        ShopCategory = "handmade"
	CategoryLocalShop           ShopCategory = "local"
	CategoryHalalShop           ShopCategory = "halal"
	CategoryOrganicShop         ShopCategory = "organic"
	CategoryWholesaleShop       ShopCategory = "wholesale"
	CategoryRetailShop          ShopCategory = "retail"
	CategoryMarketplaceShop     ShopCategory = "marketplace"
	CategoryDropshippingShop    ShopCategory = "dropshipping"
)

// ShopStatus represents the current status of a shop
type ShopStatus string

const (
	ShopStatusPending    ShopStatus = "pending"
	ShopStatusActive     ShopStatus = "active"
	ShopStatusInactive   ShopStatus = "inactive"
	ShopStatusSuspended  ShopStatus = "suspended"
	ShopStatusClosed     ShopStatus = "closed"
	ShopStatusUnderReview ShopStatus = "under_review"
	ShopStatusRejected   ShopStatus = "rejected"
)

// ShopVisibility represents shop visibility settings
type ShopVisibility string

const (
	ShopVisibilityPublic    ShopVisibility = "public"
	ShopVisibilityPrivate   ShopVisibility = "private"
	ShopVisibilityHidden    ShopVisibility = "hidden"
	ShopVisibilityComingSoon ShopVisibility = "coming_soon"
)

// ShopAddress represents the shop's physical address
type ShopAddress struct {
	BusinessName string  `json:"business_name"`
	AddressLine1 string  `json:"address_line1"`
	AddressLine2 *string `json:"address_line2,omitempty"`
	City         string  `json:"city"`
	State        string  `json:"state"`
	PostalCode   string  `json:"postal_code"`
	Country      string  `json:"country"` // ISO country code
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
	IsHeadquarters bool   `json:"is_headquarters"`
	IsWarehouse    bool   `json:"is_warehouse"`
	IsRetailStore  bool   `json:"is_retail_store"`
}

// BusinessInfo represents business information and registration
type BusinessInfo struct {
	BusinessType       string    `json:"business_type"`       // sole_proprietorship, partnership, corporation, etc.
	RegistrationNumber *string   `json:"registration_number,omitempty"`
	TaxID              *string   `json:"tax_id,omitempty"`
	VATNumber          *string   `json:"vat_number,omitempty"`
	LicenseNumber      *string   `json:"license_number,omitempty"`
	Industry           string    `json:"industry,omitempty"`
	YearEstablished    *int      `json:"year_established,omitempty"`
	EmployeeCount      *int      `json:"employee_count,omitempty"`
	AnnualRevenue      *string   `json:"annual_revenue,omitempty"`
	CertificationFiles []string  `json:"certification_files,omitempty"`
	InsuranceInfo      *Insurance `json:"insurance_info,omitempty"`
	BankingInfo        BankingInfo `json:"banking_info,omitempty"`
}

// BankingInfo represents banking and financial information
type BankingInfo struct {
	BankName        string  `json:"bank_name"`
	AccountNumber   string  `json:"account_number"`   // Masked/encrypted
	AccountType     string  `json:"account_type"`     // checking, savings, business
	RoutingNumber   *string `json:"routing_number,omitempty"`
	SwiftCode       *string `json:"swift_code,omitempty"`
	AccountHolder   string  `json:"account_holder"`
	IsVerified      bool    `json:"is_verified"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty"`
}

// ShopSettings represents shop configuration settings
type ShopSettings struct {
	General         GeneralSettings   `json:"general"`
	Checkout        CheckoutSettings  `json:"checkout"`
	Shipping        ShippingSettings  `json:"shipping"`
	Tax             TaxSettings       `json:"tax"`
	Inventory       InventorySettings `json:"inventory"`
	Notifications   NotificationSettings `json:"notifications"`
	SEO             SEOSettings       `json:"seo"`
	Security        SecuritySettings  `json:"security"`
	Integrations    IntegrationSettings `json:"integrations"`
}

// GeneralSettings represents general shop settings
type GeneralSettings struct {
	AutoApproveOrders  bool     `json:"auto_approve_orders"`
	AllowGuestCheckout bool     `json:"allow_guest_checkout"`
	RequireEmailVerification bool `json:"require_email_verification"`
	DefaultProductVisibility string `json:"default_product_visibility"`
	EnableReviews      bool     `json:"enable_reviews"`
	EnableWishlist     bool     `json:"enable_wishlist"`
	EnableCompareProducts bool  `json:"enable_compare_products"`
	EnableSocialSharing bool    `json:"enable_social_sharing"`
	EnableLiveChat     bool     `json:"enable_live_chat"`
	SupportedLanguages []string `json:"supported_languages"`
	DefaultLanguage    string   `json:"default_language"`
	EnableMultiCurrency bool    `json:"enable_multi_currency"`
	SupportedCurrencies []string `json:"supported_currencies"`
}

// CheckoutSettings represents checkout configuration
type CheckoutSettings struct {
	EnableTermsAndConditions bool     `json:"enable_terms_and_conditions"`
	TermsURL                 *string  `json:"terms_url,omitempty"`
	EnablePrivacyPolicy      bool     `json:"enable_privacy_policy"`
	PrivacyPolicyURL         *string  `json:"privacy_policy_url,omitempty"`
	EnableNewsletterSignup   bool     `json:"enable_newsletter_signup"`
	EnableOrderNotes         bool     `json:"enable_order_notes"`
	EnableCoupons            bool     `json:"enable_coupons"`
	EnableGiftCards          bool     `json:"enable_gift_cards"`
	MinimumOrderAmount       *int64   `json:"minimum_order_amount,omitempty"`
	MaximumOrderAmount       *int64   `json:"maximum_order_amount,omitempty"`
	EnableAddressValidation  bool     `json:"enable_address_validation"`
	RequiredCheckoutFields   []string `json:"required_checkout_fields"`
}

// ShippingSettings represents shipping configuration
type ShippingSettings struct {
	EnableShipping        bool           `json:"enable_shipping"`
	FreeShippingThreshold *int64         `json:"free_shipping_threshold,omitempty"`
	ShippingZones         []ShippingZone `json:"shipping_zones"`
	DefaultProcessingTime int            `json:"default_processing_time"` // days
	MaxProcessingTime     int            `json:"max_processing_time"`     // days
	EnableShippingCalculator bool        `json:"enable_shipping_calculator"`
	EnableSignatureRequired  bool        `json:"enable_signature_required"`
	EnableInsurance         bool         `json:"enable_insurance"`
	EnableTrackingNumbers   bool         `json:"enable_tracking_numbers"`
	RestrictedCountries     []string     `json:"restricted_countries"`
	EnableLocalDelivery     bool         `json:"enable_local_delivery"`
	LocalDeliveryRadius     *float64     `json:"local_delivery_radius,omitempty"` // km
}

// ShippingZone represents a shipping zone
type ShippingZone struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Countries   []string `json:"countries"`
	Methods     []ShippingZoneMethod `json:"methods"`
	IsDefault   bool     `json:"is_default"`
}

// ShippingZoneMethod represents a shipping method in a zone
type ShippingZoneMethod struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`        // flat_rate, free, calculated
	Cost        *int64   `json:"cost,omitempty"`
	MinCost     *int64   `json:"min_cost,omitempty"`
	MaxCost     *int64   `json:"max_cost,omitempty"`
	Enabled     bool     `json:"enabled"`
	Carrier     *string  `json:"carrier,omitempty"`
	ServiceType *string  `json:"service_type,omitempty"`
}

// TaxSettings represents tax configuration
type TaxSettings struct {
	EnableTax           bool        `json:"enable_tax"`
	TaxIncludedInPrices bool        `json:"tax_included_in_prices"`
	DefaultTaxRate      float64     `json:"default_tax_rate"`
	TaxZones            []TaxZone   `json:"tax_zones"`
	TaxExemptCustomers  []string    `json:"tax_exempt_customers"`
	EnableVAT           bool        `json:"enable_vat"`
	VATRegistrationNumber *string   `json:"vat_registration_number,omitempty"`
}

// TaxZone represents a tax zone
type TaxZone struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Countries []string   `json:"countries"`
	States    []string   `json:"states,omitempty"`
	Rate      float64    `json:"rate"`
	Type      string     `json:"type"` // percentage, fixed
}

// InventorySettings represents inventory management settings
type InventorySettings struct {
	TrackInventory           bool `json:"track_inventory"`
	AllowBackorders          bool `json:"allow_backorders"`
	ShowStockQuantity        bool `json:"show_stock_quantity"`
	LowStockThreshold        int  `json:"low_stock_threshold"`
	OutOfStockThreshold      int  `json:"out_of_stock_threshold"`
	EnableInventoryReports   bool `json:"enable_inventory_reports"`
	EnableAutomaticReordering bool `json:"enable_automatic_reordering"`
	ReorderPoint             int  `json:"reorder_point"`
	ReorderQuantity          int  `json:"reorder_quantity"`
}

// NotificationSettings represents notification preferences
type NotificationSettings struct {
	OrderNotifications      OrderNotifications      `json:"order_notifications"`
	InventoryNotifications  InventoryNotifications  `json:"inventory_notifications"`
	CustomerNotifications   CustomerNotifications   `json:"customer_notifications"`
	MarketingNotifications  MarketingNotifications  `json:"marketing_notifications"`
	SecurityNotifications   SecurityNotifications   `json:"security_notifications"`
}

// OrderNotifications represents order-related notification settings
type OrderNotifications struct {
	NewOrder         bool `json:"new_order"`
	OrderCancelled   bool `json:"order_cancelled"`
	OrderRefunded    bool `json:"order_refunded"`
	OrderShipped     bool `json:"order_shipped"`
	OrderDelivered   bool `json:"order_delivered"`
	PaymentReceived  bool `json:"payment_received"`
	PaymentFailed    bool `json:"payment_failed"`
}

// InventoryNotifications represents inventory notification settings
type InventoryNotifications struct {
	LowStock       bool `json:"low_stock"`
	OutOfStock     bool `json:"out_of_stock"`
	RestockAlert   bool `json:"restock_alert"`
	NewProduct     bool `json:"new_product"`
	ProductUpdated bool `json:"product_updated"`
}

// CustomerNotifications represents customer notification settings
type CustomerNotifications struct {
	NewCustomer      bool `json:"new_customer"`
	CustomerUpdated  bool `json:"customer_updated"`
	CustomerReview   bool `json:"customer_review"`
	CustomerMessage  bool `json:"customer_message"`
	CustomerReturn   bool `json:"customer_return"`
}

// MarketingNotifications represents marketing notification settings
type MarketingNotifications struct {
	CampaignLaunched  bool `json:"campaign_launched"`
	PromotionStarted  bool `json:"promotion_started"`
	PromotionEnded    bool `json:"promotion_ended"`
	NewsletterSignup  bool `json:"newsletter_signup"`
	AbandonedCart     bool `json:"abandoned_cart"`
}

// SecurityNotifications represents security notification settings
type SecurityNotifications struct {
	SuspiciousActivity bool `json:"suspicious_activity"`
	LoginAttempts      bool `json:"login_attempts"`
	PasswordChanged    bool `json:"password_changed"`
	DataExport         bool `json:"data_export"`
	APIKeyGenerated    bool `json:"api_key_generated"`
}

// SEOSettings represents SEO configuration
type SEOSettings struct {
	MetaTitle          string   `json:"meta_title,omitempty"`
	MetaDescription    string   `json:"meta_description,omitempty"`
	MetaKeywords       []string `json:"meta_keywords,omitempty"`
	OpenGraphTitle     string   `json:"og_title,omitempty"`
	OpenGraphDescription string `json:"og_description,omitempty"`
	OpenGraphImage     string   `json:"og_image,omitempty"`
	TwitterCard        string   `json:"twitter_card,omitempty"`
	GoogleAnalyticsID  *string  `json:"google_analytics_id,omitempty"`
	GoogleTagManagerID *string  `json:"google_tag_manager_id,omitempty"`
	FacebookPixelID    *string  `json:"facebook_pixel_id,omitempty"`
	EnableSitemap      bool     `json:"enable_sitemap"`
	EnableRobotsTxt    bool     `json:"enable_robots_txt"`
	EnableStructuredData bool   `json:"enable_structured_data"`
}

// SecuritySettings represents security configuration
type SecuritySettings struct {
	EnableTwoFactor       bool     `json:"enable_two_factor"`
	RequireStrongPasswords bool    `json:"require_strong_passwords"`
	EnableLoginNotifications bool  `json:"enable_login_notifications"`
	AllowedIPAddresses    []string `json:"allowed_ip_addresses,omitempty"`
	BlockedIPAddresses    []string `json:"blocked_ip_addresses,omitempty"`
	EnableSSL             bool     `json:"enable_ssl"`
	EnableCSRFProtection  bool     `json:"enable_csrf_protection"`
	EnableRateLimiting    bool     `json:"enable_rate_limiting"`
	MaxLoginAttempts      int      `json:"max_login_attempts"`
	SessionTimeoutMinutes int      `json:"session_timeout_minutes"`
}

// IntegrationSettings represents third-party integrations
type IntegrationSettings struct {
	PaymentGateways   []PaymentGateway   `json:"payment_gateways"`
	ShippingProviders []ShippingProvider `json:"shipping_providers"`
	EmailProviders    []EmailProvider    `json:"email_providers"`
	AnalyticsProviders []AnalyticsProvider `json:"analytics_providers"`
	MarketingTools    []MarketingTool    `json:"marketing_tools"`
	ERP               *ERPIntegration    `json:"erp,omitempty"`
	CRM               *CRMIntegration    `json:"crm,omitempty"`
	Accounting        *AccountingIntegration `json:"accounting,omitempty"`
}

// PaymentGateway represents a payment gateway integration
type PaymentGateway struct {
	Provider    string            `json:"provider"`
	Enabled     bool              `json:"enabled"`
	IsDefault   bool              `json:"is_default"`
	Config      map[string]string `json:"config"`
	SupportedMethods []string     `json:"supported_methods"`
}

// ShippingProvider represents a shipping provider integration
type ShippingProvider struct {
	Provider    string            `json:"provider"`
	Enabled     bool              `json:"enabled"`
	IsDefault   bool              `json:"is_default"`
	Config      map[string]string `json:"config"`
	Services    []string          `json:"services"`
}

// EmailProvider represents an email service integration
type EmailProvider struct {
	Provider string            `json:"provider"`
	Enabled  bool              `json:"enabled"`
	Config   map[string]string `json:"config"`
}

// AnalyticsProvider represents an analytics service integration
type AnalyticsProvider struct {
	Provider string            `json:"provider"`
	Enabled  bool              `json:"enabled"`
	Config   map[string]string `json:"config"`
}

// MarketingTool represents a marketing tool integration
type MarketingTool struct {
	Tool     string            `json:"tool"`
	Enabled  bool              `json:"enabled"`
	Config   map[string]string `json:"config"`
}

// ERPIntegration represents ERP system integration
type ERPIntegration struct {
	System   string            `json:"system"`
	Enabled  bool              `json:"enabled"`
	Config   map[string]string `json:"config"`
	SyncFrequency string       `json:"sync_frequency"`
}

// CRMIntegration represents CRM system integration
type CRMIntegration struct {
	System   string            `json:"system"`
	Enabled  bool              `json:"enabled"`
	Config   map[string]string `json:"config"`
	SyncCustomers bool         `json:"sync_customers"`
	SyncOrders    bool         `json:"sync_orders"`
}

// AccountingIntegration represents accounting system integration
type AccountingIntegration struct {
	System   string            `json:"system"`
	Enabled  bool              `json:"enabled"`
	Config   map[string]string `json:"config"`
	SyncInvoices bool          `json:"sync_invoices"`
	SyncPayments bool          `json:"sync_payments"`
}

// ShopFeatures represents enabled features for the shop
type ShopFeatures struct {
	MultiChannel      bool `json:"multi_channel"`
	Dropshipping      bool `json:"dropshipping"`
	Wholesale         bool `json:"wholesale"`
	Subscriptions     bool `json:"subscriptions"`
	DigitalProducts   bool `json:"digital_products"`
	GiftCards         bool `json:"gift_cards"`
	LoyaltyProgram    bool `json:"loyalty_program"`
	AffiliateProgram  bool `json:"affiliate_program"`
	AdvancedAnalytics bool `json:"advanced_analytics"`
	APIAccess         bool `json:"api_access"`
	CustomBranding    bool `json:"custom_branding"`
	MultiLanguage     bool `json:"multi_language"`
	MultiCurrency     bool `json:"multi_currency"`
	MobileApp         bool `json:"mobile_app"`
	POSIntegration    bool `json:"pos_integration"`
}

// Subscription represents the shop's subscription plan
type Subscription struct {
	PlanID          string     `json:"plan_id"`
	PlanName        string     `json:"plan_name"`
	PlanType        string     `json:"plan_type"`      // free, basic, premium, enterprise
	BillingCycle    string     `json:"billing_cycle"`  // monthly, yearly
	Price           int64      `json:"price"`          // in cents
	Currency        Currency   `json:"currency"`
	Status          string     `json:"status"`         // active, cancelled, expired, suspended
	StartDate       time.Time  `json:"start_date"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	NextBillingDate *time.Time `json:"next_billing_date,omitempty"`
	Features        []string   `json:"features"`
	Limits          PlanLimits `json:"limits"`
	IsTrialPeriod   bool       `json:"is_trial_period"`
	TrialEndDate    *time.Time `json:"trial_end_date,omitempty"`
	AutoRenew       bool       `json:"auto_renew"`
}

// PlanLimits represents subscription plan limits
type PlanLimits struct {
	MaxProducts      int `json:"max_products"`
	MaxOrders        int `json:"max_orders"`
	MaxStorage       int `json:"max_storage"`       // in MB
	MaxBandwidth     int `json:"max_bandwidth"`     // in GB
	MaxStaffAccounts int `json:"max_staff_accounts"`
	MaxAPIRequests   int `json:"max_api_requests"`  // per month
}

// ShopAnalytics represents analytics data
type ShopAnalytics struct {
	TotalProducts    int            `json:"total_products"`
	ActiveProducts   int            `json:"active_products"`
	TotalOrders      int            `json:"total_orders"`
	TotalRevenue     int64          `json:"total_revenue"`     // in cents
	MonthlyRevenue   int64          `json:"monthly_revenue"`   // in cents
	TotalCustomers   int            `json:"total_customers"`
	ActiveCustomers  int            `json:"active_customers"`
	ConversionRate   float64        `json:"conversion_rate"`
	AverageOrderValue int64         `json:"average_order_value"` // in cents
	TopProducts      []ProductStats `json:"top_products"`
	TopCategories    []CategoryStats `json:"top_categories"`
	MonthlyStats     MonthlyStats   `json:"monthly_stats"`
	TrafficSources   TrafficSources `json:"traffic_sources"`
	LastUpdated      time.Time      `json:"last_updated"`
}

// ProductStats represents product statistics
type ProductStats struct {
	ProductID   uuid.UUID `json:"product_id"`
	Name        string    `json:"name"`
	TotalSales  int       `json:"total_sales"`
	Revenue     int64     `json:"revenue"`
	ViewCount   int       `json:"view_count"`
}

// CategoryStats represents category statistics
type CategoryStats struct {
	Category   string `json:"category"`
	TotalSales int    `json:"total_sales"`
	Revenue    int64  `json:"revenue"`
	ProductCount int  `json:"product_count"`
}

// MonthlyStats represents monthly statistics
type MonthlyStats struct {
	Orders     int   `json:"orders"`
	Revenue    int64 `json:"revenue"`
	Customers  int   `json:"customers"`
	Products   int   `json:"products"`
	Views      int   `json:"views"`
	Month      int   `json:"month"`
	Year       int   `json:"year"`
}

// TrafficSources represents traffic source analytics
type TrafficSources struct {
	Direct     int `json:"direct"`
	Search     int `json:"search"`
	Social     int `json:"social"`
	Email      int `json:"email"`
	Referral   int `json:"referral"`
	Paid       int `json:"paid"`
	Other      int `json:"other"`
}

// ComplianceInfo represents compliance and regulatory information
type ComplianceInfo struct {
	GDPR         GDPRCompliance `json:"gdpr"`
	CCPA         CCPACompliance `json:"ccpa"`
	SOX          bool           `json:"sox"`
	PCI          PCICompliance  `json:"pci"`
	DataRetention DataRetention `json:"data_retention"`
	PrivacyPolicy PrivacyPolicy `json:"privacy_policy"`
	TermsOfService TermsOfService `json:"terms_of_service"`
	CookiePolicy  CookiePolicy   `json:"cookie_policy"`
	AccessibilityCompliance AccessibilityCompliance `json:"accessibility"`
}

// GDPRCompliance represents GDPR compliance status
type GDPRCompliance struct {
	Applicable        bool       `json:"applicable"`
	ConsentManagement bool       `json:"consent_management"`
	DataPortability   bool       `json:"data_portability"`
	RightToErasure    bool       `json:"right_to_erasure"`
	PrivacyByDesign   bool       `json:"privacy_by_design"`
	DPOAppointed      bool       `json:"dpo_appointed"`
	LastAssessment    *time.Time `json:"last_assessment,omitempty"`
}

// CCPACompliance represents CCPA compliance status
type CCPACompliance struct {
	Applicable         bool       `json:"applicable"`
	DoNotSellLink      bool       `json:"do_not_sell_link"`
	PrivacyRightsNotice bool      `json:"privacy_rights_notice"`
	ConsumerRequests   bool       `json:"consumer_requests"`
	LastAssessment     *time.Time `json:"last_assessment,omitempty"`
}

// PCICompliance represents PCI DSS compliance
type PCICompliance struct {
	Level             string     `json:"level"`
	CertificationDate *time.Time `json:"certification_date,omitempty"`
	ExpiryDate        *time.Time `json:"expiry_date,omitempty"`
	AOCOnFile         bool       `json:"aoc_on_file"`
	QSAValidated      bool       `json:"qsa_validated"`
}

// DataRetention represents data retention policies
type DataRetention struct {
	CustomerData   int `json:"customer_data"`   // months
	OrderData      int `json:"order_data"`      // months
	PaymentData    int `json:"payment_data"`    // months
	AnalyticsData  int `json:"analytics_data"`  // months
	LogData        int `json:"log_data"`        // months
	AutoDeletion   bool `json:"auto_deletion"`
}

// PrivacyPolicy represents privacy policy information
type PrivacyPolicy struct {
	URL           string     `json:"url"`
	LastUpdated   time.Time  `json:"last_updated"`
	Version       string     `json:"version"`
	Language      string     `json:"language"`
	Approved      bool       `json:"approved"`
	ApprovedBy    *string    `json:"approved_by,omitempty"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
}

// TermsOfService represents terms of service information
type TermsOfService struct {
	URL           string     `json:"url"`
	LastUpdated   time.Time  `json:"last_updated"`
	Version       string     `json:"version"`
	Language      string     `json:"language"`
	Approved      bool       `json:"approved"`
	ApprovedBy    *string    `json:"approved_by,omitempty"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
}

// CookiePolicy represents cookie policy information
type CookiePolicy struct {
	URL           string     `json:"url"`
	LastUpdated   time.Time  `json:"last_updated"`
	Version       string     `json:"version"`
	ConsentBanner bool       `json:"consent_banner"`
	Categories    []string   `json:"categories"`
}

// AccessibilityCompliance represents accessibility compliance
type AccessibilityCompliance struct {
	WCAGLevel       string     `json:"wcag_level"`       // A, AA, AAA
	LastAudit       *time.Time `json:"last_audit,omitempty"`
	AuditReport     *string    `json:"audit_report,omitempty"`
	Compliant       bool       `json:"compliant"`
	RemediationPlan *string    `json:"remediation_plan,omitempty"`
}

// SocialMedia represents social media links
type SocialMedia struct {
	Facebook  *string `json:"facebook,omitempty"`
	Instagram *string `json:"instagram,omitempty"`
	Twitter   *string `json:"twitter,omitempty"`
	YouTube   *string `json:"youtube,omitempty"`
	LinkedIn  *string `json:"linkedin,omitempty"`
	TikTok    *string `json:"tiktok,omitempty"`
	Pinterest *string `json:"pinterest,omitempty"`
	WhatsApp  *string `json:"whatsapp,omitempty"`
	Line      *string `json:"line,omitempty"`
	Telegram  *string `json:"telegram,omitempty"`
	WeChat    *string `json:"wechat,omitempty"`
}

// ShopMetadata represents additional shop metadata
type ShopMetadata struct {
	FoundedYear       *int               `json:"founded_year,omitempty"`
	StaffCount        *int               `json:"staff_count,omitempty"`
	MonthlyVisitors   *int               `json:"monthly_visitors,omitempty"`
	CustomerRating    *float64           `json:"customer_rating,omitempty"`
	ReviewCount       *int               `json:"review_count,omitempty"`
	Awards            []Award            `json:"awards,omitempty"`
	Certifications    []Certification    `json:"certifications,omitempty"`
	Partnerships      []Partnership      `json:"partnerships,omitempty"`
	MediaMentions     []MediaMention     `json:"media_mentions,omitempty"`
	Keywords          []string           `json:"keywords,omitempty"`
	CompetitorAnalysis CompetitorAnalysis `json:"competitor_analysis,omitempty"`
	MarketPosition    MarketPosition     `json:"market_position,omitempty"`
	GrowthMetrics     GrowthMetrics      `json:"growth_metrics,omitempty"`
}

// Award represents an award or recognition
type Award struct {
	Name        string    `json:"name"`
	Organization string   `json:"organization"`
	Year        int       `json:"year"`
	Description *string   `json:"description,omitempty"`
	URL         *string   `json:"url,omitempty"`
}

// Certification represents a business certification
type Certification struct {
	Name          string     `json:"name"`
	Issuer        string     `json:"issuer"`
	IssueDate     time.Time  `json:"issue_date"`
	ExpiryDate    *time.Time `json:"expiry_date,omitempty"`
	CertificateURL *string   `json:"certificate_url,omitempty"`
	Verified      bool       `json:"verified"`
}

// Partnership represents a business partnership
type Partnership struct {
	PartnerName   string  `json:"partner_name"`
	PartnerType   string  `json:"partner_type"`   // supplier, distributor, affiliate, etc.
	StartDate     time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	Description   *string `json:"description,omitempty"`
	URL           *string `json:"url,omitempty"`
	IsActive      bool    `json:"is_active"`
}

// MediaMention represents media coverage
type MediaMention struct {
	PublicationName string    `json:"publication_name"`
	Title           string    `json:"title"`
	PublishDate     time.Time `json:"publish_date"`
	URL             *string   `json:"url,omitempty"`
	Sentiment       string    `json:"sentiment"` // positive, negative, neutral
	ReachEstimate   *int      `json:"reach_estimate,omitempty"`
}

// CompetitorAnalysis represents competitive analysis data
type CompetitorAnalysis struct {
	MainCompetitors   []string `json:"main_competitors,omitempty"`
	MarketShare       *float64 `json:"market_share,omitempty"`
	PricePosition     string   `json:"price_position"`     // premium, mid-range, budget
	QualityPosition   string   `json:"quality_position"`   // high, medium, low
	UniqueValueProp   *string  `json:"unique_value_prop,omitempty"`
	CompetitiveAdvantages []string `json:"competitive_advantages,omitempty"`
	LastAnalysisDate  *time.Time `json:"last_analysis_date,omitempty"`
}

// MarketPosition represents market positioning
type MarketPosition struct {
	TargetMarket      string   `json:"target_market"`
	MarketSegments    []string `json:"market_segments"`
	GeographicReach   string   `json:"geographic_reach"` // local, regional, national, international
	CustomerSegments  []string `json:"customer_segments"`
	BrandPositioning  string   `json:"brand_positioning"`
	ValueProposition  string   `json:"value_proposition"`
}

// GrowthMetrics represents growth and performance metrics
type GrowthMetrics struct {
	MonthlyGrowthRate  *float64   `json:"monthly_growth_rate,omitempty"`
	YearlyGrowthRate   *float64   `json:"yearly_growth_rate,omitempty"`
	CustomerRetentionRate *float64 `json:"customer_retention_rate,omitempty"`
	CustomerAcquisitionCost *int64 `json:"customer_acquisition_cost,omitempty"`
	CustomerLifetimeValue *int64   `json:"customer_lifetime_value,omitempty"`
	ChurnRate          *float64   `json:"churn_rate,omitempty"`
	RevenuePerCustomer *int64     `json:"revenue_per_customer,omitempty"`
	LastCalculated     *time.Time `json:"last_calculated,omitempty"`
}

// Value implementations for JSON fields

func (sa ShopAddress) Value() (driver.Value, error) {
	return json.Marshal(sa)
}

func (sa *ShopAddress) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into ShopAddress", value)
	}

	return json.Unmarshal(jsonData, sa)
}

func (bi BusinessInfo) Value() (driver.Value, error) {
	return json.Marshal(bi)
}

func (bi *BusinessInfo) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into BusinessInfo", value)
	}

	return json.Unmarshal(jsonData, bi)
}

func (ss ShopSettings) Value() (driver.Value, error) {
	return json.Marshal(ss)
}

func (ss *ShopSettings) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into ShopSettings", value)
	}

	return json.Unmarshal(jsonData, ss)
}

func (sf ShopFeatures) Value() (driver.Value, error) {
	return json.Marshal(sf)
}

func (sf *ShopFeatures) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into ShopFeatures", value)
	}

	return json.Unmarshal(jsonData, sf)
}

func (s Subscription) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *Subscription) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into Subscription", value)
	}

	return json.Unmarshal(jsonData, s)
}

func (sa ShopAnalytics) Value() (driver.Value, error) {
	return json.Marshal(sa)
}

func (sa *ShopAnalytics) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into ShopAnalytics", value)
	}

	return json.Unmarshal(jsonData, sa)
}

func (ci ComplianceInfo) Value() (driver.Value, error) {
	return json.Marshal(ci)
}

func (ci *ComplianceInfo) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into ComplianceInfo", value)
	}

	return json.Unmarshal(jsonData, ci)
}

func (sm SocialMedia) Value() (driver.Value, error) {
	return json.Marshal(sm)
}

func (sm *SocialMedia) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into SocialMedia", value)
	}

	return json.Unmarshal(jsonData, sm)
}

func (sm ShopMetadata) Value() (driver.Value, error) {
	return json.Marshal(sm)
}

func (sm *ShopMetadata) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into ShopMetadata", value)
	}

	return json.Unmarshal(jsonData, sm)
}

// Validation helper functions

func ValidShopCategories() []ShopCategory {
	return []ShopCategory{
		CategoryElectronicsShop, CategoryFashionShop, CategoryHomeGardenShop,
		CategoryBeautyShop, CategorySportsShop, CategoryBooksShop,
		CategoryFoodBeverageShop, CategoryHealthShop, CategoryToysShop,
		CategoryAutomotiveShop, CategoryServicesShop, CategoryDigitalShop,
		CategoryHandmadeShop, CategoryLocalShop, CategoryHalalShop,
		CategoryOrganicShop, CategoryWholesaleShop, CategoryRetailShop,
		CategoryMarketplaceShop, CategoryDropshippingShop,
	}
}

func ValidShopStatuses() []ShopStatus {
	return []ShopStatus{
		ShopStatusPending, ShopStatusActive, ShopStatusInactive,
		ShopStatusSuspended, ShopStatusClosed, ShopStatusUnderReview,
		ShopStatusRejected,
	}
}

func ValidShopVisibilities() []ShopVisibility {
	return []ShopVisibility{
		ShopVisibilityPublic, ShopVisibilityPrivate,
		ShopVisibilityHidden, ShopVisibilityComingSoon,
	}
}

// Validation methods

func (sc ShopCategory) IsValid() bool {
	for _, valid := range ValidShopCategories() {
		if sc == valid {
			return true
		}
	}
	return false
}

func (ss ShopStatus) IsValid() bool {
	for _, valid := range ValidShopStatuses() {
		if ss == valid {
			return true
		}
	}
	return false
}

func (sv ShopVisibility) IsValid() bool {
	for _, valid := range ValidShopVisibilities() {
		if sv == valid {
			return true
		}
	}
	return false
}

// Business logic methods

func (ss ShopStatus) CanAcceptOrders() bool {
	return ss == ShopStatusActive
}

func (ss ShopStatus) IsOperational() bool {
	return ss == ShopStatusActive || ss == ShopStatusInactive
}

func (sv ShopVisibility) IsPubliclyVisible() bool {
	return sv == ShopVisibilityPublic || sv == ShopVisibilityComingSoon
}

// Main Shop validation

func (s *Shop) Validate() error {
	var errs []string

	// Owner ID validation
	if s.OwnerID == uuid.Nil {
		errs = append(errs, "owner_id is required")
	}

	// Name validation
	if strings.TrimSpace(s.Name) == "" {
		errs = append(errs, "name is required")
	}
	if len(s.Name) > 255 {
		errs = append(errs, "name cannot exceed 255 characters")
	}

	// Slug validation
	if strings.TrimSpace(s.Slug) == "" {
		errs = append(errs, "slug is required")
	}
	if !s.isValidSlug(s.Slug) {
		errs = append(errs, "slug must contain only lowercase letters, numbers, and hyphens")
	}

	// Category validation
	if !s.Category.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid category: %s", s.Category))
	}

	// Email validation
	if strings.TrimSpace(s.Email) == "" {
		errs = append(errs, "email is required")
	}
	if !s.isValidEmail(s.Email) {
		errs = append(errs, "invalid email format")
	}

	// Status validation
	if !s.Status.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid status: %s", s.Status))
	}

	// Visibility validation
	if !s.Visibility.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid visibility: %s", s.Visibility))
	}

	// Country validation
	if len(s.Country) != 2 {
		errs = append(errs, "country must be a 2-letter ISO code")
	}

	// Currency validation
	if !s.Currency.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid currency: %s", s.Currency))
	}

	// Language validation
	if len(s.Language) < 2 || len(s.Language) > 5 {
		errs = append(errs, "language must be a valid ISO language code")
	}

	// Timezone validation
	if strings.TrimSpace(s.Timezone) == "" {
		errs = append(errs, "timezone is required")
	}

	// Phone validation if provided
	if s.Phone != nil && *s.Phone != "" {
		if !s.isValidPhone(*s.Phone) {
			errs = append(errs, "invalid phone number format")
		}
	}

	// Website validation if provided
	if s.Website != nil && *s.Website != "" {
		if !s.isValidURL(*s.Website) {
			errs = append(errs, "invalid website URL format")
		}
	}

	// Address validation
	if err := s.validateAddress(); err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

func (s *Shop) isValidSlug(slug string) bool {
	// Only lowercase letters, numbers, and hyphens
	slugRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
	return slugRegex.MatchString(slug) && len(slug) <= 100
}

func (s *Shop) isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func (s *Shop) isValidPhone(phone string) bool {
	// International format validation
	phoneRegex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return phoneRegex.MatchString(phone)
}

func (s *Shop) isValidURL(url string) bool {
	urlRegex := regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}/?.*$`)
	return urlRegex.MatchString(url)
}

func (s *Shop) validateAddress() error {
	addr := s.Address

	if strings.TrimSpace(addr.BusinessName) == "" {
		return errors.New("address business_name is required")
	}

	if strings.TrimSpace(addr.AddressLine1) == "" {
		return errors.New("address address_line1 is required")
	}

	if strings.TrimSpace(addr.City) == "" {
		return errors.New("address city is required")
	}

	if strings.TrimSpace(addr.Country) == "" {
		return errors.New("address country is required")
	}

	if len(addr.Country) != 2 {
		return errors.New("address country must be a 2-letter ISO code")
	}

	return nil
}

// Shop lifecycle methods

func (s *Shop) BeforeCreate() error {
	// Generate UUID if not set
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now().UTC()
	s.CreatedAt = now
	s.UpdatedAt = now
	s.LastActivityAt = &now

	// Set default values
	if s.Status == "" {
		s.Status = ShopStatusPending
	}

	if s.Visibility == "" {
		s.Visibility = ShopVisibilityPublic
	}

	// Generate slug if not provided
	if strings.TrimSpace(s.Slug) == "" {
		s.Slug = s.generateSlug()
	}

	// Initialize default settings
	s.initializeDefaultSettings()

	// Initialize default features
	s.initializeDefaultFeatures()

	// Initialize default analytics
	s.initializeDefaultAnalytics()

	// Validate before creation
	return s.Validate()
}

func (s *Shop) BeforeUpdate() error {
	// Update timestamp
	s.UpdatedAt = time.Now().UTC()

	// Update last activity
	if s.Status == ShopStatusActive {
		now := time.Now().UTC()
		s.LastActivityAt = &now
	}

	// Validate before update
	return s.Validate()
}

func (s *Shop) generateSlug() string {
	// Simple slug generation from shop name
	slug := strings.ToLower(s.Name)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters (simplified)
	slug = strings.ReplaceAll(slug, "'", "")
	slug = strings.ReplaceAll(slug, "\"", "")

	// Limit length
	if len(slug) > 100 {
		slug = slug[:100]
	}

	return slug
}

func (s *Shop) initializeDefaultSettings() {
	s.Settings = ShopSettings{
		General: GeneralSettings{
			AutoApproveOrders:       true,
			AllowGuestCheckout:      true,
			RequireEmailVerification: true,
			DefaultProductVisibility: "public",
			EnableReviews:           true,
			EnableWishlist:          true,
			EnableCompareProducts:   false,
			EnableSocialSharing:     true,
			EnableLiveChat:          false,
			SupportedLanguages:      []string{s.Language},
			DefaultLanguage:         s.Language,
			EnableMultiCurrency:     false,
			SupportedCurrencies:     []string{string(s.Currency)},
		},
		Checkout: CheckoutSettings{
			EnableTermsAndConditions: true,
			EnablePrivacyPolicy:      true,
			EnableNewsletterSignup:   false,
			EnableOrderNotes:         true,
			EnableCoupons:           true,
			EnableGiftCards:         false,
			EnableAddressValidation: true,
			RequiredCheckoutFields:  []string{"email", "name", "address", "phone"},
		},
		Shipping: ShippingSettings{
			EnableShipping:          true,
			DefaultProcessingTime:   2,
			MaxProcessingTime:       7,
			EnableShippingCalculator: true,
			EnableSignatureRequired: false,
			EnableInsurance:         false,
			EnableTrackingNumbers:   true,
			EnableLocalDelivery:     false,
		},
		Tax: TaxSettings{
			EnableTax:           true,
			TaxIncludedInPrices: false,
			DefaultTaxRate:      0.0,
			EnableVAT:           false,
		},
		Inventory: InventorySettings{
			TrackInventory:            true,
			AllowBackorders:           false,
			ShowStockQuantity:         true,
			LowStockThreshold:         10,
			OutOfStockThreshold:       0,
			EnableInventoryReports:    true,
			EnableAutomaticReordering: false,
			ReorderPoint:              20,
			ReorderQuantity:           50,
		},
		SEO: SEOSettings{
			EnableSitemap:        true,
			EnableRobotsTxt:      true,
			EnableStructuredData: true,
		},
		Security: SecuritySettings{
			EnableTwoFactor:         false,
			RequireStrongPasswords: true,
			EnableLoginNotifications: true,
			EnableSSL:              true,
			EnableCSRFProtection:   true,
			EnableRateLimiting:     true,
			MaxLoginAttempts:       5,
			SessionTimeoutMinutes:  1440, // 24 hours
		},
	}
}

func (s *Shop) initializeDefaultFeatures() {
	s.Features = ShopFeatures{
		MultiChannel:      false,
		Dropshipping:      false,
		Wholesale:         false,
		Subscriptions:     false,
		DigitalProducts:   true,
		GiftCards:         false,
		LoyaltyProgram:    false,
		AffiliateProgram:  false,
		AdvancedAnalytics: false,
		APIAccess:         false,
		CustomBranding:    false,
		MultiLanguage:     false,
		MultiCurrency:     false,
		MobileApp:         false,
		POSIntegration:    false,
	}
}

func (s *Shop) initializeDefaultAnalytics() {
	s.Analytics = ShopAnalytics{
		TotalProducts:    0,
		ActiveProducts:   0,
		TotalOrders:      0,
		TotalRevenue:     0,
		MonthlyRevenue:   0,
		TotalCustomers:   0,
		ActiveCustomers:  0,
		ConversionRate:   0.0,
		AverageOrderValue: 0,
		TopProducts:      []ProductStats{},
		TopCategories:    []CategoryStats{},
		MonthlyStats:     MonthlyStats{},
		TrafficSources:   TrafficSources{},
		LastUpdated:      time.Now().UTC(),
	}
}

// State transition methods

func (s *Shop) Activate() error {
	if s.Status == ShopStatusActive {
		return errors.New("shop is already active")
	}

	if s.Status == ShopStatusClosed {
		return errors.New("cannot activate a closed shop")
	}

	s.Status = ShopStatusActive
	s.UpdatedAt = time.Now().UTC()

	return s.Validate()
}

func (s *Shop) Deactivate() error {
	if s.Status == ShopStatusInactive {
		return errors.New("shop is already inactive")
	}

	s.Status = ShopStatusInactive
	s.UpdatedAt = time.Now().UTC()

	return s.Validate()
}

func (s *Shop) Suspend(reason string) error {
	if s.Status == ShopStatusSuspended {
		return errors.New("shop is already suspended")
	}

	s.Status = ShopStatusSuspended
	s.UpdatedAt = time.Now().UTC()

	// Add to metadata (simplified)
	if s.Metadata.Keywords == nil {
		s.Metadata.Keywords = []string{}
	}
	s.Metadata.Keywords = append(s.Metadata.Keywords, fmt.Sprintf("suspended:%s", reason))

	return s.Validate()
}

func (s *Shop) Close() error {
	if s.Status == ShopStatusClosed {
		return errors.New("shop is already closed")
	}

	s.Status = ShopStatusClosed
	s.Visibility = ShopVisibilityHidden
	s.UpdatedAt = time.Now().UTC()

	return s.Validate()
}

func (s *Shop) Verify() error {
	if s.IsVerified {
		return errors.New("shop is already verified")
	}

	now := time.Now().UTC()
	s.IsVerified = true
	s.VerifiedAt = &now
	s.UpdatedAt = now

	return s.Validate()
}

func (s *Shop) Feature() error {
	if s.IsFeatured {
		return errors.New("shop is already featured")
	}

	now := time.Now().UTC()
	s.IsFeatured = true
	s.FeaturedAt = &now
	s.UpdatedAt = now

	return s.Validate()
}

func (s *Shop) Unfeature() error {
	if !s.IsFeatured {
		return errors.New("shop is not featured")
	}

	s.IsFeatured = false
	s.FeaturedAt = nil
	s.UpdatedAt = time.Now().UTC()

	return s.Validate()
}

// Business logic methods

func (s *Shop) CanAcceptOrders() bool {
	return s.Status.CanAcceptOrders() && s.Visibility.IsPubliclyVisible()
}

func (s *Shop) IsOperational() bool {
	return s.Status.IsOperational()
}

func (s *Shop) GetShopURL() string {
	return fmt.Sprintf("https://tchat.shop/%s", s.Slug)
}

func (s *Shop) GetDaysSinceCreated() int {
	return int(time.Since(s.CreatedAt).Hours() / 24)
}

func (s *Shop) UpdateLastActivity() {
	now := time.Now().UTC()
	s.LastActivityAt = &now
	s.UpdatedAt = now
}

// Analytics methods

func (s *Shop) UpdateAnalytics(analytics ShopAnalytics) {
	analytics.LastUpdated = time.Now().UTC()
	s.Analytics = analytics
	s.UpdatedAt = time.Now().UTC()
}

func (s *Shop) GetFormattedRevenue() string {
	return s.Currency.FormatAmount(s.Analytics.TotalRevenue)
}

func (s *Shop) GetFormattedMonthlyRevenue() string {
	return s.Currency.FormatAmount(s.Analytics.MonthlyRevenue)
}

func (s *Shop) GetFormattedAverageOrderValue() string {
	return s.Currency.FormatAmount(s.Analytics.AverageOrderValue)
}

// Public API response methods

func (s *Shop) ToPublicShop() map[string]interface{} {
	return map[string]interface{}{
		"id":            s.ID,
		"name":          s.Name,
		"description":   s.Description,
		"slug":          s.Slug,
		"category":      s.Category,
		"email":         s.Email,
		"phone":         s.Phone,
		"website":       s.Website,
		"logo_url":      s.LogoURL,
		"banner_url":    s.BannerURL,
		"status":        s.Status,
		"visibility":    s.Visibility,
		"country":       s.Country,
		"currency":      s.Currency,
		"language":      s.Language,
		"address":       s.Address,
		"social_media":  s.SocialMedia,
		"is_verified":   s.IsVerified,
		"is_featured":   s.IsFeatured,
		"verified_at":   s.VerifiedAt,
		"featured_at":   s.FeaturedAt,
		"shop_url":      s.GetShopURL(),
		"created_at":    s.CreatedAt,
		"updated_at":    s.UpdatedAt,
	}
}

func (s *Shop) ToOwnerShop() map[string]interface{} {
	response := s.ToPublicShop()

	// Add owner-only fields
	response["owner_id"] = s.OwnerID
	response["business_info"] = s.BusinessInfo
	response["settings"] = s.Settings
	response["features"] = s.Features
	response["subscription"] = s.Subscription
	response["analytics"] = s.Analytics
	response["compliance"] = s.Compliance
	response["metadata"] = s.Metadata
	response["tags"] = s.Tags
	response["last_activity_at"] = s.LastActivityAt

	// Add formatted analytics
	response["formatted_revenue"] = s.GetFormattedRevenue()
	response["formatted_monthly_revenue"] = s.GetFormattedMonthlyRevenue()
	response["formatted_average_order_value"] = s.GetFormattedAverageOrderValue()

	return response
}

// Request structures

type ShopCreateRequest struct {
	OwnerID     uuid.UUID    `json:"owner_id" validate:"required"`
	Name        string       `json:"name" validate:"required,max=255"`
	Description *string      `json:"description,omitempty"`
	Category    ShopCategory `json:"category" validate:"required"`
	Email       string       `json:"email" validate:"required,email"`
	Phone       *string      `json:"phone,omitempty"`
	Website     *string      `json:"website,omitempty"`
	Country     string       `json:"country" validate:"required,len=2"`
	Currency    Currency     `json:"currency" validate:"required"`
	Language    string       `json:"language" validate:"required"`
	Timezone    string       `json:"timezone" validate:"required"`
	Address     ShopAddress  `json:"address" validate:"required"`
}

func (req *ShopCreateRequest) ToShop() *Shop {
	return &Shop{
		OwnerID:     req.OwnerID,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Email:       req.Email,
		Phone:       req.Phone,
		Website:     req.Website,
		Country:     req.Country,
		Currency:    req.Currency,
		Language:    req.Language,
		Timezone:    req.Timezone,
		Address:     req.Address,
	}
}

// Shop manager

type ShopManager struct {
	// Add dependencies
}

func NewShopManager() *ShopManager {
	return &ShopManager{}
}

func (sm *ShopManager) CreateShop(req *ShopCreateRequest) (*Shop, error) {
	shop := req.ToShop()

	if err := shop.BeforeCreate(); err != nil {
		return nil, fmt.Errorf("shop creation failed: %v", err)
	}

	return shop, nil
}