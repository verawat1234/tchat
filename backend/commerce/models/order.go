package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Order represents a customer order in the commerce system
type Order struct {
	ID                uuid.UUID       `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ShopID            uuid.UUID       `json:"shop_id" gorm:"type:varchar(36);not null;index"`
	CustomerID        uuid.UUID       `json:"customer_id" gorm:"type:varchar(36);not null;index"`
	OrderNumber       string          `json:"order_number" gorm:"type:varchar(50);uniqueIndex;not null"`
	Status            OrderStatus     `json:"status" gorm:"type:varchar(30);not null;index"`
	PaymentStatus     PaymentStatus   `json:"payment_status" gorm:"type:varchar(30);not null;index"`
	FulfillmentStatus FulfillmentStatus `json:"fulfillment_status" gorm:"type:varchar(30);not null;index"`
	Currency          Currency        `json:"currency" gorm:"type:varchar(3);not null"`
	TotalAmount       int64           `json:"total_amount" gorm:"not null"`      // Total in cents
	SubtotalAmount    int64           `json:"subtotal_amount" gorm:"not null"`   // Subtotal before tax/shipping
	TaxAmount         int64           `json:"tax_amount" gorm:"default:0"`       // Tax amount in cents
	ShippingAmount    int64           `json:"shipping_amount" gorm:"default:0"`  // Shipping cost in cents
	DiscountAmount    int64           `json:"discount_amount" gorm:"default:0"`  // Discount amount in cents
	RefundedAmount    int64           `json:"refunded_amount" gorm:"default:0"`  // Refunded amount in cents
	Items             OrderItemSlice  `json:"items" gorm:"type:json"`
	ShippingAddress   Address         `json:"shipping_address" gorm:"type:json"`
	BillingAddress    Address         `json:"billing_address" gorm:"type:json"`
	ShippingMethod    ShippingMethod  `json:"shipping_method" gorm:"type:json"`
	PaymentMethod     PaymentInfo     `json:"payment_method" gorm:"type:json"`
	DiscountCodes     DiscountSlice   `json:"discount_codes" gorm:"type:json"`
	Notes             *string         `json:"notes,omitempty" gorm:"type:text"`
	CustomerNotes     *string         `json:"customer_notes,omitempty" gorm:"type:text"`
	Tags              StringSlice     `json:"tags" gorm:"type:json"`
	Source            OrderSource     `json:"source" gorm:"type:varchar(30);not null"`
	SourceDetails     SourceDetails   `json:"source_details" gorm:"type:json"`
	Metadata          OrderMetadata   `json:"metadata" gorm:"type:json"`
	TrackingInfo      TrackingInfo    `json:"tracking_info" gorm:"type:json"`
	TimestampInfo     TimestampInfo   `json:"timestamps" gorm:"type:json"`
	ExternalID        *string         `json:"external_id,omitempty" gorm:"type:varchar(255);index"`
	ChannelID         *uuid.UUID      `json:"channel_id,omitempty" gorm:"type:varchar(36);index"`
	CreatedAt         time.Time       `json:"created_at" gorm:"not null"`
	UpdatedAt         time.Time       `json:"updated_at" gorm:"not null"`
}

// OrderStatus represents the overall status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
	OrderStatusFailed     OrderStatus = "failed"
	OrderStatusExpired    OrderStatus = "expired"
	OrderStatusOnHold     OrderStatus = "on_hold"
)

// PaymentStatus represents the payment status of an order
type PaymentStatus string

const (
	PaymentStatusPending     PaymentStatus = "pending"
	PaymentStatusPaid        PaymentStatus = "paid"
	PaymentStatusPartiallyPaid PaymentStatus = "partially_paid"
	PaymentStatusRefunded    PaymentStatus = "refunded"
	PaymentStatusPartiallyRefunded PaymentStatus = "partially_refunded"
	PaymentStatusFailed      PaymentStatus = "failed"
	PaymentStatusCancelled   PaymentStatus = "cancelled"
	PaymentStatusAuthorized  PaymentStatus = "authorized"
	PaymentStatusVoided      PaymentStatus = "voided"
)

// FulfillmentStatus represents the fulfillment status of an order
type FulfillmentStatus string

const (
	FulfillmentStatusUnfulfilled     FulfillmentStatus = "unfulfilled"
	FulfillmentStatusPartiallyFulfilled FulfillmentStatus = "partially_fulfilled"
	FulfillmentStatusFulfilled       FulfillmentStatus = "fulfilled"
	FulfillmentStatusShipped         FulfillmentStatus = "shipped"
	FulfillmentStatusDelivered       FulfillmentStatus = "delivered"
	FulfillmentStatusReturned        FulfillmentStatus = "returned"
	FulfillmentStatusCancelled       FulfillmentStatus = "cancelled"
)

// OrderSource represents where the order originated from
type OrderSource string

const (
	OrderSourceWeb       OrderSource = "web"
	OrderSourceMobile    OrderSource = "mobile"
	OrderSourceAPI       OrderSource = "api"
	OrderSourceChat      OrderSource = "chat"
	OrderSourceSocial    OrderSource = "social"
	OrderSourceMarketplace OrderSource = "marketplace"
	OrderSourcePOS       OrderSource = "pos"
	OrderSourcePhone     OrderSource = "phone"
	OrderSourceEmail     OrderSource = "email"
)

// OrderItemSlice represents a slice of order items
type OrderItemSlice []OrderItem

// DiscountSlice represents a slice of discount codes
type DiscountSlice []DiscountCode

// OrderItem represents an item within an order
type OrderItem struct {
	ID              string     `json:"id"`
	ProductID       uuid.UUID  `json:"product_id"`
	VariantID       *string    `json:"variant_id,omitempty"`
	SKU             string     `json:"sku"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	Image           *ProductImage `json:"image,omitempty"`
	Quantity        int        `json:"quantity"`
	UnitPrice       int64      `json:"unit_price"`      // Price per unit in cents
	TotalPrice      int64      `json:"total_price"`     // Total price for this item
	OriginalPrice   *int64     `json:"original_price,omitempty"` // Original price before discounts
	Weight          *float64   `json:"weight,omitempty"`
	Dimensions      *Dimensions `json:"dimensions,omitempty"`
	IsDigital       bool       `json:"is_digital"`
	RequiresShipping bool      `json:"requires_shipping"`
	TaxRate         float64    `json:"tax_rate"`        // Tax rate as percentage
	TaxAmount       int64      `json:"tax_amount"`      // Tax amount in cents
	Discounts       []ItemDiscount `json:"discounts,omitempty"`
	CustomAttributes map[string]interface{} `json:"custom_attributes,omitempty"`
	FulfillmentStatus FulfillmentStatus `json:"fulfillment_status"`
	GiftMessage     *string    `json:"gift_message,omitempty"`
}

// ItemDiscount represents a discount applied to an order item
type ItemDiscount struct {
	ID          string  `json:"id"`
	Code        string  `json:"code"`
	Type        string  `json:"type"`        // percentage, fixed_amount
	Value       float64 `json:"value"`       // Percentage or amount
	Amount      int64   `json:"amount"`      // Actual discount amount in cents
	Description string  `json:"description"`
}

// Address represents shipping or billing address
type Address struct {
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Company      *string `json:"company,omitempty"`
	AddressLine1 string  `json:"address_line1"`
	AddressLine2 *string `json:"address_line2,omitempty"`
	City         string  `json:"city"`
	State        string  `json:"state"`
	PostalCode   string  `json:"postal_code"`
	Country      string  `json:"country"`     // ISO country code
	Phone        *string `json:"phone,omitempty"`
	Email        *string `json:"email,omitempty"`
	IsDefault    bool    `json:"is_default"`
	AddressType  string  `json:"address_type"` // home, office, other
	Instructions *string `json:"instructions,omitempty"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
}

// ShippingMethod represents the shipping method selected
type ShippingMethod struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Description     *string        `json:"description,omitempty"`
	Carrier         string         `json:"carrier"`         // DHL, FedEx, etc.
	ServiceType     string         `json:"service_type"`    // standard, express, overnight
	Cost            int64          `json:"cost"`            // Cost in cents
	EstimatedDays   *int           `json:"estimated_days,omitempty"`
	TrackingNumber  *string        `json:"tracking_number,omitempty"`
	TrackingURL     *string        `json:"tracking_url,omitempty"`
	Insurance       *Insurance     `json:"insurance,omitempty"`
	Signature       bool           `json:"signature"`       // Requires signature
	PackagingType   string         `json:"packaging_type"`  // box, envelope, tube
	Weight          *float64       `json:"weight,omitempty"`
	Dimensions      *Dimensions    `json:"dimensions,omitempty"`
	PickupLocation  *Address       `json:"pickup_location,omitempty"`
	DropoffLocation *Address       `json:"dropoff_location,omitempty"`
}

// Insurance represents shipping insurance
type Insurance struct {
	Provider     string `json:"provider"`
	Cost         int64  `json:"cost"`         // Insurance cost in cents
	CoverageAmount int64 `json:"coverage_amount"` // Coverage amount in cents
	PolicyNumber *string `json:"policy_number,omitempty"`
}

// PaymentInfo represents payment method information
type PaymentInfo struct {
	PaymentMethodID  *uuid.UUID `json:"payment_method_id,omitempty"`
	Provider         string     `json:"provider"`         // stripe, paypal, etc.
	Type             string     `json:"type"`             // card, bank_transfer, wallet
	LastFourDigits   *string    `json:"last_four_digits,omitempty"`
	Brand            *string    `json:"brand,omitempty"`  // visa, mastercard
	ExpiryMonth      *int       `json:"expiry_month,omitempty"`
	ExpiryYear       *int       `json:"expiry_year,omitempty"`
	TransactionID    *string    `json:"transaction_id,omitempty"`
	ExternalID       *string    `json:"external_id,omitempty"`
	AuthorizationCode *string   `json:"authorization_code,omitempty"`
	ProcessorResponse *string   `json:"processor_response,omitempty"`
	RiskScore        *float64   `json:"risk_score,omitempty"`
	CVVResult        *string    `json:"cvv_result,omitempty"`
	AVSResult        *string    `json:"avs_result,omitempty"`
}

// DiscountCode represents an applied discount code
type DiscountCode struct {
	ID          string  `json:"id"`
	Code        string  `json:"code"`
	Type        string  `json:"type"`        // percentage, fixed_amount, free_shipping
	Value       float64 `json:"value"`       // Percentage or amount
	Amount      int64   `json:"amount"`      // Actual discount amount in cents
	Description string  `json:"description"`
	MinimumAmount *int64 `json:"minimum_amount,omitempty"`
	MaximumAmount *int64 `json:"maximum_amount,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	UsageCount  int     `json:"usage_count"`
	MaxUsage    *int    `json:"max_usage,omitempty"`
}

// SourceDetails represents details about the order source
type SourceDetails struct {
	UserAgent    *string `json:"user_agent,omitempty"`
	IPAddress    *string `json:"ip_address,omitempty"`
	ReferrerURL  *string `json:"referrer_url,omitempty"`
	LandingPage  *string `json:"landing_page,omitempty"`
	UTMSource    *string `json:"utm_source,omitempty"`
	UTMCampaign  *string `json:"utm_campaign,omitempty"`
	UTMMedium    *string `json:"utm_medium,omitempty"`
	UTMTerm      *string `json:"utm_term,omitempty"`
	UTMContent   *string `json:"utm_content,omitempty"`
	SessionID    *string `json:"session_id,omitempty"`
	DeviceType   *string `json:"device_type,omitempty"`
	BrowserName  *string `json:"browser_name,omitempty"`
	OSName       *string `json:"os_name,omitempty"`
	Country      *string `json:"country,omitempty"`
	Region       *string `json:"region,omitempty"`
	City         *string `json:"city,omitempty"`
}

// OrderMetadata represents additional order metadata
type OrderMetadata struct {
	InternalNotes     []InternalNote `json:"internal_notes,omitempty"`
	CustomerHistory   CustomerHistory `json:"customer_history"`
	RiskAssessment    RiskAssessment `json:"risk_assessment"`
	MarketingData     MarketingData  `json:"marketing_data"`
	CustomFields      map[string]interface{} `json:"custom_fields,omitempty"`
	IntegrationData   map[string]interface{} `json:"integration_data,omitempty"`
	AnalyticsData     AnalyticsData  `json:"analytics_data"`
	ComplianceData    ComplianceData `json:"compliance_data"`
}

// InternalNote represents internal notes about the order
type InternalNote struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	IsPublic  bool      `json:"is_public"`
}

// CustomerHistory represents customer's order history summary
type CustomerHistory struct {
	TotalOrders       int     `json:"total_orders"`
	TotalSpent        int64   `json:"total_spent"`
	AverageOrderValue int64   `json:"average_order_value"`
	FirstOrderDate    *time.Time `json:"first_order_date,omitempty"`
	LastOrderDate     *time.Time `json:"last_order_date,omitempty"`
	CustomerSegment   string  `json:"customer_segment"` // new, returning, vip
	LoyaltyTier       string  `json:"loyalty_tier,omitempty"`
	PreferredPayment  string  `json:"preferred_payment,omitempty"`
}

// RiskAssessment represents fraud risk assessment
type RiskAssessment struct {
	RiskScore        float64   `json:"risk_score"`        // 0-100
	RiskLevel        string    `json:"risk_level"`        // low, medium, high
	FraudIndicators  []string  `json:"fraud_indicators,omitempty"`
	ReviewRequired   bool      `json:"review_required"`
	AutoApproved     bool      `json:"auto_approved"`
	AssessedAt       time.Time `json:"assessed_at"`
	AssessedBy       string    `json:"assessed_by"`       // system, manual
	CVVCheck         string    `json:"cvv_check,omitempty"`
	AVSCheck         string    `json:"avs_check,omitempty"`
	IPGeolocation    string    `json:"ip_geolocation,omitempty"`
	BillingAddressMatch bool   `json:"billing_address_match"`
	ShippingAddressMatch bool  `json:"shipping_address_match"`
}

// MarketingData represents marketing and attribution data
type MarketingData struct {
	Channel          string    `json:"channel,omitempty"`
	Campaign         string    `json:"campaign,omitempty"`
	Source           string    `json:"source,omitempty"`
	Medium           string    `json:"medium,omitempty"`
	Content          string    `json:"content,omitempty"`
	Term             string    `json:"term,omitempty"`
	CouponCodes      []string  `json:"coupon_codes,omitempty"`
	AffiliateID      *string   `json:"affiliate_id,omitempty"`
	CommissionRate   *float64  `json:"commission_rate,omitempty"`
	CustomerLifetimeValue *int64 `json:"customer_lifetime_value,omitempty"`
	AcquisitionCost  *int64    `json:"acquisition_cost,omitempty"`
}

// AnalyticsData represents analytics and tracking data
type AnalyticsData struct {
	SessionDuration   *time.Duration `json:"session_duration,omitempty"`
	PageViews         int            `json:"page_views"`
	ProductViews      int            `json:"product_views"`
	CartAbandonment   bool           `json:"cart_abandonment"`
	ConversionPath    []string       `json:"conversion_path,omitempty"`
	TimeToConvert     *time.Duration `json:"time_to_convert,omitempty"`
	DeviceFingerprint *string        `json:"device_fingerprint,omitempty"`
	ABTestVariant     *string        `json:"ab_test_variant,omitempty"`
	PersonalizationID *string        `json:"personalization_id,omitempty"`
}

// ComplianceData represents compliance and regulatory data
type ComplianceData struct {
	TaxExempt         bool      `json:"tax_exempt"`
	TaxExemptionID    *string   `json:"tax_exemption_id,omitempty"`
	VATNumber         *string   `json:"vat_number,omitempty"`
	BusinessTaxID     *string   `json:"business_tax_id,omitempty"`
	RequiredDocuments []string  `json:"required_documents,omitempty"`
	ComplianceChecks  []ComplianceCheck `json:"compliance_checks,omitempty"`
	DataProcessingConsent bool   `json:"data_processing_consent"`
	MarketingConsent  bool      `json:"marketing_consent"`
	ConsentTimestamp  *time.Time `json:"consent_timestamp,omitempty"`
	GDPRApplicable    bool      `json:"gdpr_applicable"`
	CCPAApplicable    bool      `json:"ccpa_applicable"`
}

// ComplianceCheck represents a compliance check
type ComplianceCheck struct {
	Type      string    `json:"type"`      // kyc, sanctions, export_control
	Status    string    `json:"status"`    // passed, failed, pending
	CheckedAt time.Time `json:"checked_at"`
	CheckedBy string    `json:"checked_by"`
	Reference *string   `json:"reference,omitempty"`
	Notes     *string   `json:"notes,omitempty"`
}

// TrackingInfo represents shipment tracking information
type TrackingInfo struct {
	TrackingNumber    string          `json:"tracking_number,omitempty"`
	Carrier           string          `json:"carrier,omitempty"`
	TrackingURL       string          `json:"tracking_url,omitempty"`
	Status            string          `json:"status,omitempty"`
	EstimatedDelivery *time.Time      `json:"estimated_delivery,omitempty"`
	TrackingEvents    []TrackingEvent `json:"tracking_events,omitempty"`
	SignedBy          *string         `json:"signed_by,omitempty"`
	DeliveryNotes     *string         `json:"delivery_notes,omitempty"`
}

// TrackingEvent represents a tracking event
type TrackingEvent struct {
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Location    string    `json:"location,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	Details     *string   `json:"details,omitempty"`
}

// TimestampInfo represents important timestamps
type TimestampInfo struct {
	ConfirmedAt      *time.Time `json:"confirmed_at,omitempty"`
	ProcessingAt     *time.Time `json:"processing_at,omitempty"`
	ShippedAt        *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt      *time.Time `json:"delivered_at,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	CancelledAt      *time.Time `json:"cancelled_at,omitempty"`
	RefundedAt       *time.Time `json:"refunded_at,omitempty"`
	FirstViewedAt    *time.Time `json:"first_viewed_at,omitempty"`
	LastModifiedAt   *time.Time `json:"last_modified_at,omitempty"`
	PaymentDueAt     *time.Time `json:"payment_due_at,omitempty"`
	FulfillmentDueAt *time.Time `json:"fulfillment_due_at,omitempty"`
}

// Value implementations for JSON fields

func (ois OrderItemSlice) Value() (driver.Value, error) {
	return json.Marshal(ois)
}

func (ois *OrderItemSlice) Scan(value interface{}) error {
	if value == nil {
		*ois = nil
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into OrderItemSlice", value)
	}

	return json.Unmarshal(jsonData, ois)
}

func (ds DiscountSlice) Value() (driver.Value, error) {
	return json.Marshal(ds)
}

func (ds *DiscountSlice) Scan(value interface{}) error {
	if value == nil {
		*ds = nil
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into DiscountSlice", value)
	}

	return json.Unmarshal(jsonData, ds)
}

func (a Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Address) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into Address", value)
	}

	return json.Unmarshal(jsonData, a)
}

func (sm ShippingMethod) Value() (driver.Value, error) {
	return json.Marshal(sm)
}

func (sm *ShippingMethod) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into ShippingMethod", value)
	}

	return json.Unmarshal(jsonData, sm)
}

func (pi PaymentInfo) Value() (driver.Value, error) {
	return json.Marshal(pi)
}

func (pi *PaymentInfo) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into PaymentInfo", value)
	}

	return json.Unmarshal(jsonData, pi)
}

func (sd SourceDetails) Value() (driver.Value, error) {
	return json.Marshal(sd)
}

func (sd *SourceDetails) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into SourceDetails", value)
	}

	return json.Unmarshal(jsonData, sd)
}

func (om OrderMetadata) Value() (driver.Value, error) {
	return json.Marshal(om)
}

func (om *OrderMetadata) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into OrderMetadata", value)
	}

	return json.Unmarshal(jsonData, om)
}

func (ti TrackingInfo) Value() (driver.Value, error) {
	return json.Marshal(ti)
}

func (ti *TrackingInfo) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into TrackingInfo", value)
	}

	return json.Unmarshal(jsonData, ti)
}

func (tsi TimestampInfo) Value() (driver.Value, error) {
	return json.Marshal(tsi)
}

func (tsi *TimestampInfo) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into TimestampInfo", value)
	}

	return json.Unmarshal(jsonData, tsi)
}

// Validation helper functions

func ValidOrderStatuses() []OrderStatus {
	return []OrderStatus{
		OrderStatusPending, OrderStatusConfirmed, OrderStatusProcessing,
		OrderStatusShipped, OrderStatusDelivered, OrderStatusCompleted,
		OrderStatusCancelled, OrderStatusRefunded, OrderStatusFailed,
		OrderStatusExpired, OrderStatusOnHold,
	}
}

func ValidPaymentStatuses() []PaymentStatus {
	return []PaymentStatus{
		PaymentStatusPending, PaymentStatusPaid, PaymentStatusPartiallyPaid,
		PaymentStatusRefunded, PaymentStatusPartiallyRefunded, PaymentStatusFailed,
		PaymentStatusCancelled, PaymentStatusAuthorized, PaymentStatusVoided,
	}
}

func ValidFulfillmentStatuses() []FulfillmentStatus {
	return []FulfillmentStatus{
		FulfillmentStatusUnfulfilled, FulfillmentStatusPartiallyFulfilled,
		FulfillmentStatusFulfilled, FulfillmentStatusShipped,
		FulfillmentStatusDelivered, FulfillmentStatusReturned,
		FulfillmentStatusCancelled,
	}
}

func ValidOrderSources() []OrderSource {
	return []OrderSource{
		OrderSourceWeb, OrderSourceMobile, OrderSourceAPI, OrderSourceChat,
		OrderSourceSocial, OrderSourceMarketplace, OrderSourcePOS,
		OrderSourcePhone, OrderSourceEmail,
	}
}

// Validation methods

func (os OrderStatus) IsValid() bool {
	for _, valid := range ValidOrderStatuses() {
		if os == valid {
			return true
		}
	}
	return false
}

func (ps PaymentStatus) IsValid() bool {
	for _, valid := range ValidPaymentStatuses() {
		if ps == valid {
			return true
		}
	}
	return false
}

func (fs FulfillmentStatus) IsValid() bool {
	for _, valid := range ValidFulfillmentStatuses() {
		if fs == valid {
			return true
		}
	}
	return false
}

func (osrc OrderSource) IsValid() bool {
	for _, valid := range ValidOrderSources() {
		if osrc == valid {
			return true
		}
	}
	return false
}

// Business logic methods

func (os OrderStatus) CanBeCancelled() bool {
	cancellableStatuses := []OrderStatus{
		OrderStatusPending, OrderStatusConfirmed, OrderStatusOnHold,
	}
	for _, status := range cancellableStatuses {
		if os == status {
			return true
		}
	}
	return false
}

func (os OrderStatus) CanBeRefunded() bool {
	refundableStatuses := []OrderStatus{
		OrderStatusCompleted, OrderStatusDelivered, OrderStatusShipped,
	}
	for _, status := range refundableStatuses {
		if os == status {
			return true
		}
	}
	return false
}

func (os OrderStatus) IsFinalState() bool {
	finalStates := []OrderStatus{
		OrderStatusCompleted, OrderStatusCancelled, OrderStatusRefunded,
		OrderStatusFailed, OrderStatusExpired,
	}
	for _, state := range finalStates {
		if os == state {
			return true
		}
	}
	return false
}

func (ps PaymentStatus) IsFullyPaid() bool {
	return ps == PaymentStatusPaid
}

func (ps PaymentStatus) RequiresPayment() bool {
	return ps == PaymentStatusPending || ps == PaymentStatusPartiallyPaid
}

func (fs FulfillmentStatus) IsCompleted() bool {
	return fs == FulfillmentStatusFulfilled || fs == FulfillmentStatusDelivered
}

func (fs FulfillmentStatus) RequiresFulfillment() bool {
	return fs == FulfillmentStatusUnfulfilled || fs == FulfillmentStatusPartiallyFulfilled
}

// Main Order validation

func (o *Order) Validate() error {
	var errs []string

	// Shop ID validation
	if o.ShopID == uuid.Nil {
		errs = append(errs, "shop_id is required")
	}

	// Customer ID validation
	if o.CustomerID == uuid.Nil {
		errs = append(errs, "customer_id is required")
	}

	// Order number validation
	if strings.TrimSpace(o.OrderNumber) == "" {
		errs = append(errs, "order_number is required")
	}

	// Status validations
	if !o.Status.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid order status: %s", o.Status))
	}

	if !o.PaymentStatus.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid payment status: %s", o.PaymentStatus))
	}

	if !o.FulfillmentStatus.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid fulfillment status: %s", o.FulfillmentStatus))
	}

	// Currency validation
	if !o.Currency.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid currency: %s", o.Currency))
	}

	// Amount validations
	if o.TotalAmount < 0 {
		errs = append(errs, "total_amount cannot be negative")
	}

	if o.SubtotalAmount < 0 {
		errs = append(errs, "subtotal_amount cannot be negative")
	}

	if o.TaxAmount < 0 {
		errs = append(errs, "tax_amount cannot be negative")
	}

	if o.ShippingAmount < 0 {
		errs = append(errs, "shipping_amount cannot be negative")
	}

	if o.DiscountAmount < 0 {
		errs = append(errs, "discount_amount cannot be negative")
	}

	if o.RefundedAmount < 0 {
		errs = append(errs, "refunded_amount cannot be negative")
	}

	if o.RefundedAmount > o.TotalAmount {
		errs = append(errs, "refunded_amount cannot exceed total_amount")
	}

	// Items validation
	if len(o.Items) == 0 {
		errs = append(errs, "order must have at least one item")
	}

	// Validate total amount calculation
	expectedTotal := o.SubtotalAmount + o.TaxAmount + o.ShippingAmount - o.DiscountAmount
	if o.TotalAmount != expectedTotal {
		errs = append(errs, fmt.Sprintf("total_amount mismatch: expected %d, got %d", expectedTotal, o.TotalAmount))
	}

	// Source validation
	if !o.Source.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid order source: %s", o.Source))
	}

	// Address validations
	if err := o.validateAddresses(); err != nil {
		errs = append(errs, err.Error())
	}

	// Items validation
	if err := o.validateItems(); err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

func (o *Order) validateAddresses() error {
	// Shipping address validation
	if err := o.validateAddress(o.ShippingAddress, "shipping"); err != nil {
		return err
	}

	// Billing address validation
	if err := o.validateAddress(o.BillingAddress, "billing"); err != nil {
		return err
	}

	return nil
}

func (o *Order) validateAddress(addr Address, addressType string) error {
	if strings.TrimSpace(addr.FirstName) == "" {
		return fmt.Errorf("%s address first_name is required", addressType)
	}

	if strings.TrimSpace(addr.LastName) == "" {
		return fmt.Errorf("%s address last_name is required", addressType)
	}

	if strings.TrimSpace(addr.AddressLine1) == "" {
		return fmt.Errorf("%s address address_line1 is required", addressType)
	}

	if strings.TrimSpace(addr.City) == "" {
		return fmt.Errorf("%s address city is required", addressType)
	}

	if strings.TrimSpace(addr.Country) == "" {
		return fmt.Errorf("%s address country is required", addressType)
	}

	if len(addr.Country) != 2 {
		return fmt.Errorf("%s address country must be a 2-letter ISO code", addressType)
	}

	return nil
}

func (o *Order) validateItems() error {
	itemIDs := make(map[string]bool)
	calculatedSubtotal := int64(0)

	for i, item := range o.Items {
		// Check for duplicate item IDs
		if itemIDs[item.ID] {
			return fmt.Errorf("duplicate item ID: %s", item.ID)
		}
		itemIDs[item.ID] = true

		// Validate quantity
		if item.Quantity <= 0 {
			return fmt.Errorf("item %d: quantity must be positive", i)
		}

		// Validate prices
		if item.UnitPrice < 0 {
			return fmt.Errorf("item %d: unit_price cannot be negative", i)
		}

		if item.TotalPrice < 0 {
			return fmt.Errorf("item %d: total_price cannot be negative", i)
		}

		// Validate total price calculation
		expectedTotal := item.UnitPrice * int64(item.Quantity)
		if item.TotalPrice != expectedTotal {
			return fmt.Errorf("item %d: total_price mismatch: expected %d, got %d", i, expectedTotal, item.TotalPrice)
		}

		// Validate tax amount
		if item.TaxAmount < 0 {
			return fmt.Errorf("item %d: tax_amount cannot be negative", i)
		}

		// Validate fulfillment status
		if !item.FulfillmentStatus.IsValid() {
			return fmt.Errorf("item %d: invalid fulfillment status: %s", i, item.FulfillmentStatus)
		}

		calculatedSubtotal += item.TotalPrice
	}

	// Validate subtotal calculation
	if o.SubtotalAmount != calculatedSubtotal {
		return fmt.Errorf("subtotal_amount mismatch: expected %d, got %d", calculatedSubtotal, o.SubtotalAmount)
	}

	return nil
}

// Order lifecycle methods

func (o *Order) BeforeCreate() error {
	// Generate UUID if not set
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now().UTC()
	o.CreatedAt = now
	o.UpdatedAt = now

	// Generate order number if not provided
	if o.OrderNumber == "" {
		o.OrderNumber = o.generateOrderNumber()
	}

	// Set default values
	if o.Status == "" {
		o.Status = OrderStatusPending
	}

	if o.PaymentStatus == "" {
		o.PaymentStatus = PaymentStatusPending
	}

	if o.FulfillmentStatus == "" {
		o.FulfillmentStatus = FulfillmentStatusUnfulfilled
	}

	if o.Source == "" {
		o.Source = OrderSourceWeb
	}

	// Initialize metadata if empty
	if o.Metadata == (OrderMetadata{}) {
		o.Metadata = OrderMetadata{
			CustomerHistory: CustomerHistory{
				CustomerSegment: "new",
			},
			RiskAssessment: RiskAssessment{
				RiskScore:    50.0,
				RiskLevel:    "medium",
				AssessedAt:   now,
				AssessedBy:   "system",
				AutoApproved: false,
			},
			AnalyticsData: AnalyticsData{
				PageViews:    1,
				ProductViews: len(o.Items),
			},
			ComplianceData: ComplianceData{
				DataProcessingConsent: true,
			},
		}
	}

	// Initialize timestamps
	if o.TimestampInfo == (TimestampInfo{}) {
		o.TimestampInfo = TimestampInfo{
			FirstViewedAt:   &now,
			LastModifiedAt: &now,
		}
	}

	// Set item defaults
	for i := range o.Items {
		if o.Items[i].ID == "" {
			o.Items[i].ID = uuid.New().String()
		}
		if o.Items[i].FulfillmentStatus == "" {
			o.Items[i].FulfillmentStatus = FulfillmentStatusUnfulfilled
		}
	}

	// Validate before creation
	return o.Validate()
}

func (o *Order) BeforeUpdate() error {
	// Update timestamp
	o.UpdatedAt = time.Now().UTC()
	if o.TimestampInfo.LastModifiedAt != nil {
		*o.TimestampInfo.LastModifiedAt = o.UpdatedAt
	}

	// Validate before update
	return o.Validate()
}

func (o *Order) generateOrderNumber() string {
	// Generate order number: ORD-YYYYMMDD-HHMMSS-XXXX
	now := time.Now().UTC()
	timestamp := now.Format("20060102-150405")
	shortID := o.ID.String()[:8]
	return fmt.Sprintf("ORD-%s-%s", timestamp, strings.ToUpper(shortID))
}

// State transition methods

func (o *Order) Confirm() error {
	if o.Status != OrderStatusPending {
		return fmt.Errorf("can only confirm pending orders, current status: %s", o.Status)
	}

	now := time.Now().UTC()
	o.Status = OrderStatusConfirmed
	o.TimestampInfo.ConfirmedAt = &now
	o.UpdatedAt = now

	return o.Validate()
}

func (o *Order) StartProcessing() error {
	allowedStatuses := []OrderStatus{OrderStatusConfirmed, OrderStatusOnHold}
	allowed := false
	for _, status := range allowedStatuses {
		if o.Status == status {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("cannot start processing from status: %s", o.Status)
	}

	now := time.Now().UTC()
	o.Status = OrderStatusProcessing
	o.TimestampInfo.ProcessingAt = &now
	o.UpdatedAt = now

	return o.Validate()
}

func (o *Order) Ship(trackingNumber, carrier string) error {
	if o.Status != OrderStatusProcessing {
		return fmt.Errorf("can only ship processing orders, current status: %s", o.Status)
	}

	now := time.Now().UTC()
	o.Status = OrderStatusShipped
	o.FulfillmentStatus = FulfillmentStatusShipped
	o.TimestampInfo.ShippedAt = &now

	// Update tracking info
	o.TrackingInfo.TrackingNumber = trackingNumber
	o.TrackingInfo.Carrier = carrier
	o.TrackingInfo.Status = "shipped"

	// Update shipping method if it exists
	if o.ShippingMethod.TrackingNumber == nil {
		o.ShippingMethod.TrackingNumber = &trackingNumber
	}

	o.UpdatedAt = now

	return o.Validate()
}

func (o *Order) Deliver() error {
	if o.Status != OrderStatusShipped {
		return fmt.Errorf("can only deliver shipped orders, current status: %s", o.Status)
	}

	now := time.Now().UTC()
	o.Status = OrderStatusDelivered
	o.FulfillmentStatus = FulfillmentStatusDelivered
	o.TimestampInfo.DeliveredAt = &now

	// Update tracking info
	o.TrackingInfo.Status = "delivered"

	o.UpdatedAt = now

	return o.Validate()
}

func (o *Order) Complete() error {
	if o.Status != OrderStatusDelivered {
		return fmt.Errorf("can only complete delivered orders, current status: %s", o.Status)
	}

	now := time.Now().UTC()
	o.Status = OrderStatusCompleted
	o.TimestampInfo.CompletedAt = &now
	o.UpdatedAt = now

	return o.Validate()
}

func (o *Order) Cancel(reason string) error {
	if !o.Status.CanBeCancelled() {
		return fmt.Errorf("cannot cancel order with status: %s", o.Status)
	}

	now := time.Now().UTC()
	o.Status = OrderStatusCancelled
	o.FulfillmentStatus = FulfillmentStatusCancelled
	o.TimestampInfo.CancelledAt = &now

	// Add internal note
	if reason != "" {
		note := InternalNote{
			ID:        uuid.New().String(),
			Author:    "system",
			Note:      fmt.Sprintf("Order cancelled: %s", reason),
			CreatedAt: now,
			IsPublic:  false,
		}
		o.Metadata.InternalNotes = append(o.Metadata.InternalNotes, note)
	}

	o.UpdatedAt = now

	return o.Validate()
}

func (o *Order) Refund(amount int64, reason string) error {
	if !o.Status.CanBeRefunded() {
		return fmt.Errorf("cannot refund order with status: %s", o.Status)
	}

	if amount <= 0 {
		return errors.New("refund amount must be positive")
	}

	if o.RefundedAmount+amount > o.TotalAmount {
		return errors.New("total refund amount cannot exceed order total")
	}

	now := time.Now().UTC()
	o.RefundedAmount += amount

	// Update status if fully refunded
	if o.RefundedAmount == o.TotalAmount {
		o.Status = OrderStatusRefunded
		o.PaymentStatus = PaymentStatusRefunded
		o.TimestampInfo.RefundedAt = &now
	} else {
		o.PaymentStatus = PaymentStatusPartiallyRefunded
	}

	// Add internal note
	if reason != "" {
		note := InternalNote{
			ID:        uuid.New().String(),
			Author:    "system",
			Note:      fmt.Sprintf("Refund processed: %s (%s)", o.Currency.FormatAmount(amount), reason),
			CreatedAt: now,
			IsPublic:  false,
		}
		o.Metadata.InternalNotes = append(o.Metadata.InternalNotes, note)
	}

	o.UpdatedAt = now

	return o.Validate()
}

func (o *Order) PutOnHold(reason string) error {
	if o.Status.IsFinalState() {
		return fmt.Errorf("cannot put order on hold with status: %s", o.Status)
	}

	now := time.Now().UTC()
	o.Status = OrderStatusOnHold

	// Add internal note
	if reason != "" {
		note := InternalNote{
			ID:        uuid.New().String(),
			Author:    "system",
			Note:      fmt.Sprintf("Order on hold: %s", reason),
			CreatedAt: now,
			IsPublic:  false,
		}
		o.Metadata.InternalNotes = append(o.Metadata.InternalNotes, note)
	}

	o.UpdatedAt = now

	return o.Validate()
}

// Payment status updates

func (o *Order) MarkAsPaid(transactionID string) error {
	if o.PaymentStatus == PaymentStatusPaid {
		return errors.New("order is already paid")
	}

	now := time.Now().UTC()
	o.PaymentStatus = PaymentStatusPaid

	if o.PaymentMethod.TransactionID == nil {
		o.PaymentMethod.TransactionID = &transactionID
	}

	o.UpdatedAt = now

	return o.Validate()
}

// Utility methods

func (o *Order) GetFormattedTotal() string {
	return o.Currency.FormatAmount(o.TotalAmount)
}

func (o *Order) GetFormattedSubtotal() string {
	return o.Currency.FormatAmount(o.SubtotalAmount)
}

func (o *Order) GetFormattedTax() string {
	return o.Currency.FormatAmount(o.TaxAmount)
}

func (o *Order) GetFormattedShipping() string {
	return o.Currency.FormatAmount(o.ShippingAmount)
}

func (o *Order) GetFormattedDiscount() string {
	return o.Currency.FormatAmount(o.DiscountAmount)
}

func (o *Order) GetFormattedRefunded() string {
	return o.Currency.FormatAmount(o.RefundedAmount)
}

func (o *Order) HasPhysicalItems() bool {
	for _, item := range o.Items {
		if item.RequiresShipping {
			return true
		}
	}
	return false
}

func (o *Order) HasDigitalItems() bool {
	for _, item := range o.Items {
		if item.IsDigital {
			return true
		}
	}
	return false
}

func (o *Order) GetTotalItems() int {
	total := 0
	for _, item := range o.Items {
		total += item.Quantity
	}
	return total
}

func (o *Order) IsFullyRefunded() bool {
	return o.RefundedAmount == o.TotalAmount
}

func (o *Order) IsPartiallyRefunded() bool {
	return o.RefundedAmount > 0 && o.RefundedAmount < o.TotalAmount
}

func (o *Order) GetDaysSinceCreated() int {
	return int(time.Since(o.CreatedAt).Hours() / 24)
}

// Public API response methods

func (o *Order) ToPublicOrder() map[string]interface{} {
	return map[string]interface{}{
		"id":              o.ID,
		"order_number":    o.OrderNumber,
		"status":          o.Status,
		"payment_status":  o.PaymentStatus,
		"fulfillment_status": o.FulfillmentStatus,
		"currency":        o.Currency,
		"total_amount":    o.TotalAmount,
		"formatted_total": o.GetFormattedTotal(),
		"subtotal_amount": o.SubtotalAmount,
		"formatted_subtotal": o.GetFormattedSubtotal(),
		"tax_amount":      o.TaxAmount,
		"formatted_tax":   o.GetFormattedTax(),
		"shipping_amount": o.ShippingAmount,
		"formatted_shipping": o.GetFormattedShipping(),
		"discount_amount": o.DiscountAmount,
		"formatted_discount": o.GetFormattedDiscount(),
		"items":           o.Items,
		"shipping_address": o.ShippingAddress,
		"tracking_info":   o.TrackingInfo,
		"created_at":      o.CreatedAt,
		"updated_at":      o.UpdatedAt,
	}
}

func (o *Order) ToAdminOrder() map[string]interface{} {
	response := o.ToPublicOrder()

	// Add admin-only fields
	response["shop_id"] = o.ShopID
	response["customer_id"] = o.CustomerID
	response["billing_address"] = o.BillingAddress
	response["payment_method"] = o.PaymentMethod
	response["discount_codes"] = o.DiscountCodes
	response["notes"] = o.Notes
	response["tags"] = o.Tags
	response["source"] = o.Source
	response["source_details"] = o.SourceDetails
	response["metadata"] = o.Metadata
	response["timestamps"] = o.TimestampInfo
	response["refunded_amount"] = o.RefundedAmount
	response["formatted_refunded"] = o.GetFormattedRefunded()
	response["external_id"] = o.ExternalID
	response["channel_id"] = o.ChannelID

	return response
}

// Request structures

type OrderCreateRequest struct {
	ShopID          uuid.UUID      `json:"shop_id" validate:"required"`
	CustomerID      uuid.UUID      `json:"customer_id" validate:"required"`
	Items           []OrderItem    `json:"items" validate:"required,min=1"`
	ShippingAddress Address        `json:"shipping_address" validate:"required"`
	BillingAddress  *Address       `json:"billing_address,omitempty"`
	ShippingMethod  *ShippingMethod `json:"shipping_method,omitempty"`
	PaymentMethod   *PaymentInfo   `json:"payment_method,omitempty"`
	DiscountCodes   []DiscountCode `json:"discount_codes,omitempty"`
	Notes           *string        `json:"notes,omitempty"`
	CustomerNotes   *string        `json:"customer_notes,omitempty"`
	Tags            []string       `json:"tags,omitempty"`
	Source          OrderSource    `json:"source"`
	Currency        Currency       `json:"currency" validate:"required"`
	ChannelID       *uuid.UUID     `json:"channel_id,omitempty"`
}

func (req *OrderCreateRequest) ToOrder() *Order {
	// Calculate amounts (simplified - real implementation would be more complex)
	subtotal := int64(0)
	for _, item := range req.Items {
		subtotal += item.TotalPrice
	}

	shippingAmount := int64(0)
	if req.ShippingMethod != nil {
		shippingAmount = req.ShippingMethod.Cost
	}

	discountAmount := int64(0)
	for _, discount := range req.DiscountCodes {
		discountAmount += discount.Amount
	}

	total := subtotal + shippingAmount - discountAmount

	billingAddr := req.ShippingAddress
	if req.BillingAddress != nil {
		billingAddr = *req.BillingAddress
	}

	return &Order{
		ShopID:            req.ShopID,
		CustomerID:        req.CustomerID,
		Currency:          req.Currency,
		TotalAmount:       total,
		SubtotalAmount:    subtotal,
		ShippingAmount:    shippingAmount,
		DiscountAmount:    discountAmount,
		Items:             req.Items,
		ShippingAddress:   req.ShippingAddress,
		BillingAddress:    billingAddr,
		Notes:             req.Notes,
		CustomerNotes:     req.CustomerNotes,
		Tags:              req.Tags,
		Source:            req.Source,
		ChannelID:         req.ChannelID,
	}
}

// Order manager

type OrderManager struct {
	// Add dependencies
}

func NewOrderManager() *OrderManager {
	return &OrderManager{}
}

func (om *OrderManager) CreateOrder(req *OrderCreateRequest) (*Order, error) {
	order := req.ToOrder()

	if err := order.BeforeCreate(); err != nil {
		return nil, fmt.Errorf("order creation failed: %v", err)
	}

	return order, nil
}