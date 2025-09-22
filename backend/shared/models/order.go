package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
	OrderStatusReturned   OrderStatus = "returned"
)

// IsValid checks if the order status is valid
func (os OrderStatus) IsValid() bool {
	switch os {
	case OrderStatusPending, OrderStatusConfirmed, OrderStatusProcessing,
		 OrderStatusShipped, OrderStatusDelivered, OrderStatusCancelled,
		 OrderStatusRefunded, OrderStatusReturned:
		return true
	default:
		return false
	}
}

// IsTerminal checks if the order status is terminal
func (os OrderStatus) IsTerminal() bool {
	switch os {
	case OrderStatusDelivered, OrderStatusCancelled, OrderStatusRefunded, OrderStatusReturned:
		return true
	default:
		return false
	}
}

// CanCancel checks if the order can be cancelled
func (os OrderStatus) CanCancel() bool {
	return os == OrderStatusPending || os == OrderStatusConfirmed
}

// PaymentStatus represents the payment status of an order
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusAuthorized PaymentStatus = "authorized"
	PaymentStatusPaid       PaymentStatus = "paid"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
	PaymentStatusRefunded   PaymentStatus = "refunded"
	PaymentStatusPartialRefund PaymentStatus = "partial_refund"
)

// IsValid checks if the payment status is valid
func (ps PaymentStatus) IsValid() bool {
	switch ps {
	case PaymentStatusPending, PaymentStatusAuthorized, PaymentStatusPaid,
		 PaymentStatusFailed, PaymentStatusCancelled, PaymentStatusRefunded,
		 PaymentStatusPartialRefund:
		return true
	default:
		return false
	}
}

// FulfillmentStatus represents the fulfillment status of an order
type FulfillmentStatus string

const (
	FulfillmentStatusUnfulfilled    FulfillmentStatus = "unfulfilled"
	FulfillmentStatusPartiallyFulfilled FulfillmentStatus = "partially_fulfilled"
	FulfillmentStatusFulfilled      FulfillmentStatus = "fulfilled"
	FulfillmentStatusShipped        FulfillmentStatus = "shipped"
	FulfillmentStatusDelivered      FulfillmentStatus = "delivered"
	FulfillmentStatusReturned       FulfillmentStatus = "returned"
)

// IsValid checks if the fulfillment status is valid
func (fs FulfillmentStatus) IsValid() bool {
	switch fs {
	case FulfillmentStatusUnfulfilled, FulfillmentStatusPartiallyFulfilled,
		 FulfillmentStatusFulfilled, FulfillmentStatusShipped,
		 FulfillmentStatusDelivered, FulfillmentStatusReturned:
		return true
	default:
		return false
	}
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID              uuid.UUID `json:"id" gorm:"column:id;type:uuid;default:gen_random_uuid()"`
	ProductID       uuid.UUID `json:"product_id" gorm:"column:product_id;type:uuid;not null"`
	ProductVariantID *uuid.UUID `json:"product_variant_id,omitempty" gorm:"column:product_variant_id;type:uuid"`
	ProductName     string    `json:"product_name" gorm:"column:product_name;size:200;not null"`
	ProductSKU      string    `json:"product_sku,omitempty" gorm:"column:product_sku;size:100"`
	Quantity        int       `json:"quantity" gorm:"column:quantity;not null"`
	UnitPrice       decimal.Decimal `json:"unit_price" gorm:"column:unit_price;type:decimal(20,8);not null"`
	TotalPrice      decimal.Decimal `json:"total_price" gorm:"column:total_price;type:decimal(20,8);not null"`
	Currency        string    `json:"currency" gorm:"column:currency;size:3;not null"`
	TaxAmount       decimal.Decimal `json:"tax_amount" gorm:"column:tax_amount;type:decimal(20,8);default:0"`
	DiscountAmount  decimal.Decimal `json:"discount_amount" gorm:"column:discount_amount;type:decimal(20,8);default:0"`
	Weight          *decimal.Decimal `json:"weight,omitempty" gorm:"column:weight;type:decimal(10,3)"`
	Metadata        map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
}

// ShippingAddress represents a shipping address
type ShippingAddress struct {
	FirstName   string `json:"first_name" gorm:"column:first_name;size:100;not null"`
	LastName    string `json:"last_name" gorm:"column:last_name;size:100;not null"`
	Company     string `json:"company,omitempty" gorm:"column:company;size:100"`
	AddressLine1 string `json:"address_line1" gorm:"column:address_line1;size:200;not null"`
	AddressLine2 string `json:"address_line2,omitempty" gorm:"column:address_line2;size:200"`
	City        string `json:"city" gorm:"column:city;size:100;not null"`
	State       string `json:"state,omitempty" gorm:"column:state;size:100"`
	PostalCode  string `json:"postal_code" gorm:"column:postal_code;size:20"`
	Country     string `json:"country" gorm:"column:country;size:2;not null"`
	Phone       string `json:"phone,omitempty" gorm:"column:phone;size:20"`
	Email       string `json:"email,omitempty" gorm:"column:email;size:255"`
}

// BillingAddress represents a billing address
type BillingAddress struct {
	FirstName   string `json:"first_name" gorm:"column:first_name;size:100;not null"`
	LastName    string `json:"last_name" gorm:"column:last_name;size:100;not null"`
	Company     string `json:"company,omitempty" gorm:"column:company;size:100"`
	AddressLine1 string `json:"address_line1" gorm:"column:address_line1;size:200;not null"`
	AddressLine2 string `json:"address_line2,omitempty" gorm:"column:address_line2;size:200"`
	City        string `json:"city" gorm:"column:city;size:100;not null"`
	State       string `json:"state,omitempty" gorm:"column:state;size:100"`
	PostalCode  string `json:"postal_code" gorm:"column:postal_code;size:20"`
	Country     string `json:"country" gorm:"column:country;size:2;not null"`
	Phone       string `json:"phone,omitempty" gorm:"column:phone;size:20"`
	Email       string `json:"email,omitempty" gorm:"column:email;size:255"`
	TaxID       string `json:"tax_id,omitempty" gorm:"column:tax_id;size:50"`
}

// ShippingMethod represents shipping method information
type ShippingMethod struct {
	Name         string          `json:"name" gorm:"column:name;size:100;not null"`
	Provider     string          `json:"provider" gorm:"column:provider;size:50"`
	ServiceType  string          `json:"service_type" gorm:"column:service_type;size:50"`
	Cost         decimal.Decimal `json:"cost" gorm:"column:cost;type:decimal(20,8);not null"`
	Currency     string          `json:"currency" gorm:"column:currency;size:3;not null"`
	EstimatedDays int            `json:"estimated_days" gorm:"column:estimated_days"`
	TrackingNumber string        `json:"tracking_number,omitempty" gorm:"column:tracking_number;size:100"`
	TrackingURL  string          `json:"tracking_url,omitempty" gorm:"column:tracking_url;size:500"`
}

// OrderTotals represents the order totals breakdown
type OrderTotals struct {
	SubtotalAmount  decimal.Decimal `json:"subtotal_amount" gorm:"column:subtotal_amount;type:decimal(20,8);not null"`
	TaxAmount       decimal.Decimal `json:"tax_amount" gorm:"column:tax_amount;type:decimal(20,8);default:0"`
	ShippingAmount  decimal.Decimal `json:"shipping_amount" gorm:"column:shipping_amount;type:decimal(20,8);default:0"`
	DiscountAmount  decimal.Decimal `json:"discount_amount" gorm:"column:discount_amount;type:decimal(20,8);default:0"`
	TotalAmount     decimal.Decimal `json:"total_amount" gorm:"column:total_amount;type:decimal(20,8);not null"`
	Currency        string          `json:"currency" gorm:"column:currency;size:3;not null"`
}

// OrderDiscount represents discount information
type OrderDiscount struct {
	Code        string          `json:"code,omitempty" gorm:"column:code;size:50"`
	Type        string          `json:"type" gorm:"column:type;size:20;not null"`
	Amount      decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(20,8);not null"`
	Currency    string          `json:"currency" gorm:"column:currency;size:3;not null"`
	Description string          `json:"description,omitempty" gorm:"column:description;size:200"`
}

// OrderTimestamps represents important order timestamps
type OrderTimestamps struct {
	ConfirmedAt  *time.Time `json:"confirmed_at,omitempty" gorm:"column:confirmed_at"`
	ProcessedAt  *time.Time `json:"processed_at,omitempty" gorm:"column:processed_at"`
	ShippedAt    *time.Time `json:"shipped_at,omitempty" gorm:"column:shipped_at"`
	DeliveredAt  *time.Time `json:"delivered_at,omitempty" gorm:"column:delivered_at"`
	CancelledAt  *time.Time `json:"cancelled_at,omitempty" gorm:"column:cancelled_at"`
	RefundedAt   *time.Time `json:"refunded_at,omitempty" gorm:"column:refunded_at"`
}

// Order represents an order in the system
type Order struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderNumber string   `json:"order_number" gorm:"column:order_number;size:50;not null;uniqueIndex"`
	BusinessID uuid.UUID `json:"business_id" gorm:"type:uuid;not null;index"`
	CustomerID uuid.UUID `json:"customer_id" gorm:"type:uuid;not null;index"`

	// Order status
	Status            OrderStatus       `json:"status" gorm:"column:status;type:varchar(20);not null;default:'pending'"`
	PaymentStatus     PaymentStatus     `json:"payment_status" gorm:"column:payment_status;type:varchar(20);not null;default:'pending'"`
	FulfillmentStatus FulfillmentStatus `json:"fulfillment_status" gorm:"column:fulfillment_status;type:varchar(20);not null;default:'unfulfilled'"`

	// Order items
	Items []OrderItem `json:"items" gorm:"column:items;type:jsonb"`
	ItemCount int      `json:"item_count" gorm:"column:item_count;default:0"`

	// Addresses
	ShippingAddress ShippingAddress `json:"shipping_address" gorm:"embedded;embeddedPrefix:shipping_"`
	BillingAddress  BillingAddress  `json:"billing_address" gorm:"embedded;embeddedPrefix:billing_"`

	// Shipping information
	ShippingMethod ShippingMethod `json:"shipping_method" gorm:"embedded;embeddedPrefix:shipping_method_"`
	RequiresShipping bool         `json:"requires_shipping" gorm:"column:requires_shipping;default:true"`

	// Financial information
	Totals    OrderTotals     `json:"totals" gorm:"embedded;embeddedPrefix:totals_"`
	Discounts []OrderDiscount `json:"discounts,omitempty" gorm:"column:discounts;type:jsonb"`

	// Payment information
	PaymentGateway        string `json:"payment_gateway" gorm:"column:payment_gateway;size:50"`
	PaymentTransactionID  string `json:"payment_transaction_id,omitempty" gorm:"column:payment_transaction_id;size:255"`
	PaymentReference      string `json:"payment_reference,omitempty" gorm:"column:payment_reference;size:255"`

	// Customer information
	CustomerEmail string `json:"customer_email" gorm:"column:customer_email;size:255;not null"`
	CustomerPhone string `json:"customer_phone,omitempty" gorm:"column:customer_phone;size:20"`
	CustomerNotes string `json:"customer_notes,omitempty" gorm:"column:customer_notes;size:1000"`

	// Regional compliance
	DataRegion     string `json:"data_region" gorm:"column:data_region;size:20"`
	TaxRegion      string `json:"tax_region" gorm:"column:tax_region;size:20"`
	ComplianceData map[string]interface{} `json:"compliance_data,omitempty" gorm:"column:compliance_data;type:jsonb"`

	// Order processing
	Source         string `json:"source" gorm:"column:source;size:50;default:'web'"`
	SourceDetails  map[string]interface{} `json:"source_details,omitempty" gorm:"column:source_details;type:jsonb"`
	ProcessingNotes string `json:"processing_notes,omitempty" gorm:"column:processing_notes;size:1000"`

	// Important timestamps
	Timestamps OrderTimestamps `json:"timestamps" gorm:"embedded;embeddedPrefix:ts_"`

	// Metadata and tags
	Metadata map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
	Tags     []string               `json:"tags,omitempty" gorm:"column:tags;type:jsonb"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Relationships
	Business *Business `json:"business,omitempty" gorm:"foreignKey:BusinessID;references:ID"`
	Customer *User     `json:"customer,omitempty" gorm:"foreignKey:CustomerID;references:ID"`
}

// TableName returns the table name for the Order model
func (Order) TableName() string {
	return "orders"
}

// BeforeCreate sets up the order before creation
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}

	// Generate order number if not provided
	if o.OrderNumber == "" {
		o.OrderNumber = o.generateOrderNumber()
	}

	// Set data region based on business or shipping address
	if o.DataRegion == "" {
		if o.ShippingAddress.Country != "" {
			o.DataRegion = GetDataRegionForCountry(o.ShippingAddress.Country)
		} else {
			var business Business
			if err := tx.First(&business, o.BusinessID).Error; err == nil {
				o.DataRegion = GetDataRegionForCountry(business.Address.Country)
			} else {
				o.DataRegion = "sea-central" // Default region
			}
		}
	}

	// Set tax region
	if o.TaxRegion == "" {
		o.TaxRegion = o.ShippingAddress.Country
	}

	// Calculate totals
	o.calculateTotals()

	// Set item count
	o.ItemCount = len(o.Items)

	// Validate the order
	if err := o.Validate(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate validates the order before updating
func (o *Order) BeforeUpdate(tx *gorm.DB) error {
	// Recalculate totals if items changed
	o.calculateTotals()
	o.ItemCount = len(o.Items)

	return o.Validate()
}

// Validate validates the order data
func (o *Order) Validate() error {
	// Validate UUIDs
	if o.ID == uuid.Nil {
		return fmt.Errorf("order ID cannot be nil")
	}
	if o.BusinessID == uuid.Nil {
		return fmt.Errorf("business ID cannot be nil")
	}
	if o.CustomerID == uuid.Nil {
		return fmt.Errorf("customer ID cannot be nil")
	}

	// Validate statuses
	if !o.Status.IsValid() {
		return fmt.Errorf("invalid order status: %s", o.Status)
	}
	if !o.PaymentStatus.IsValid() {
		return fmt.Errorf("invalid payment status: %s", o.PaymentStatus)
	}
	if !o.FulfillmentStatus.IsValid() {
		return fmt.Errorf("invalid fulfillment status: %s", o.FulfillmentStatus)
	}

	// Validate order number
	if o.OrderNumber == "" {
		return fmt.Errorf("order number is required")
	}

	// Validate customer email
	if o.CustomerEmail == "" {
		return fmt.Errorf("customer email is required")
	}

	// Validate items
	if len(o.Items) == 0 {
		return fmt.Errorf("order must have at least one item")
	}

	if err := o.validateItems(); err != nil {
		return err
	}

	// Validate addresses
	if err := o.validateShippingAddress(); err != nil {
		return err
	}

	if err := o.validateBillingAddress(); err != nil {
		return err
	}

	// Validate totals
	if err := o.validateTotals(); err != nil {
		return err
	}

	return nil
}

// validateItems validates order items
func (o *Order) validateItems() error {
	for i, item := range o.Items {
		if item.ProductID == uuid.Nil {
			return fmt.Errorf("product ID is required for item %d", i)
		}

		if item.ProductName == "" {
			return fmt.Errorf("product name is required for item %d", i)
		}

		if item.Quantity <= 0 {
			return fmt.Errorf("quantity must be positive for item %d", i)
		}

		if item.UnitPrice.IsNegative() || item.UnitPrice.IsZero() {
			return fmt.Errorf("unit price must be positive for item %d", i)
		}

		if !IsValidCurrency(item.Currency) {
			return fmt.Errorf("invalid currency for item %d: %s", i, item.Currency)
		}

		// Validate calculated total price
		expectedTotal := item.UnitPrice.Mul(decimal.NewFromInt(int64(item.Quantity)))
		expectedTotal = expectedTotal.Sub(item.DiscountAmount).Add(item.TaxAmount)
		if !expectedTotal.Equal(item.TotalPrice) {
			return fmt.Errorf("item %d total price mismatch: expected %s, got %s",
				i, expectedTotal.String(), item.TotalPrice.String())
		}
	}

	return nil
}

// validateShippingAddress validates the shipping address
func (o *Order) validateShippingAddress() error {
	if o.RequiresShipping {
		if o.ShippingAddress.FirstName == "" {
			return fmt.Errorf("shipping first name is required")
		}
		if o.ShippingAddress.LastName == "" {
			return fmt.Errorf("shipping last name is required")
		}
		if o.ShippingAddress.AddressLine1 == "" {
			return fmt.Errorf("shipping address line 1 is required")
		}
		if o.ShippingAddress.City == "" {
			return fmt.Errorf("shipping city is required")
		}
		if o.ShippingAddress.Country == "" {
			return fmt.Errorf("shipping country is required")
		}
		if !IsValidSEACountry(o.ShippingAddress.Country) {
			return fmt.Errorf("invalid shipping country: %s", o.ShippingAddress.Country)
		}
	}

	return nil
}

// validateBillingAddress validates the billing address
func (o *Order) validateBillingAddress() error {
	if o.BillingAddress.FirstName == "" {
		return fmt.Errorf("billing first name is required")
	}
	if o.BillingAddress.LastName == "" {
		return fmt.Errorf("billing last name is required")
	}
	if o.BillingAddress.AddressLine1 == "" {
		return fmt.Errorf("billing address line 1 is required")
	}
	if o.BillingAddress.City == "" {
		return fmt.Errorf("billing city is required")
	}
	if o.BillingAddress.Country == "" {
		return fmt.Errorf("billing country is required")
	}
	if !IsValidSEACountry(o.BillingAddress.Country) {
		return fmt.Errorf("invalid billing country: %s", o.BillingAddress.Country)
	}

	return nil
}

// validateTotals validates the order totals
func (o *Order) validateTotals() error {
	if !IsValidCurrency(o.Totals.Currency) {
		return fmt.Errorf("invalid currency: %s", o.Totals.Currency)
	}

	if o.Totals.SubtotalAmount.IsNegative() {
		return fmt.Errorf("subtotal amount cannot be negative")
	}

	if o.Totals.TotalAmount.IsNegative() {
		return fmt.Errorf("total amount cannot be negative")
	}

	// Validate calculated total
	expectedTotal := o.Totals.SubtotalAmount.
		Add(o.Totals.TaxAmount).
		Add(o.Totals.ShippingAmount).
		Sub(o.Totals.DiscountAmount)

	if !expectedTotal.Equal(o.Totals.TotalAmount) {
		return fmt.Errorf("total amount mismatch: expected %s, got %s",
			expectedTotal.String(), o.Totals.TotalAmount.String())
	}

	return nil
}

// generateOrderNumber generates a unique order number
func (o *Order) generateOrderNumber() string {
	timestamp := time.Now().Format("20060102")
	return fmt.Sprintf("ORD-%s-%s", timestamp, o.ID.String()[:8])
}

// calculateTotals calculates order totals from items
func (o *Order) calculateTotals() {
	subtotal := decimal.Zero
	totalTax := decimal.Zero
	totalDiscount := decimal.Zero

	// Get currency from first item or default to USD
	currency := "USD"
	if len(o.Items) > 0 {
		currency = o.Items[0].Currency
	}

	// Calculate item totals
	for _, item := range o.Items {
		itemSubtotal := item.UnitPrice.Mul(decimal.NewFromInt(int64(item.Quantity)))
		subtotal = subtotal.Add(itemSubtotal)
		totalTax = totalTax.Add(item.TaxAmount)
		totalDiscount = totalDiscount.Add(item.DiscountAmount)
	}

	// Add order-level discounts
	for _, discount := range o.Discounts {
		totalDiscount = totalDiscount.Add(discount.Amount)
	}

	// Set totals
	o.Totals.SubtotalAmount = subtotal
	o.Totals.TaxAmount = totalTax
	o.Totals.DiscountAmount = totalDiscount
	o.Totals.Currency = currency

	// Calculate final total
	o.Totals.TotalAmount = subtotal.
		Add(o.Totals.TaxAmount).
		Add(o.Totals.ShippingAmount).
		Sub(o.Totals.DiscountAmount)
}

// UpdateStatus updates the order status with appropriate timestamps
func (o *Order) UpdateStatus(newStatus OrderStatus) error {
	if !newStatus.IsValid() {
		return fmt.Errorf("invalid order status: %s", newStatus)
	}

	oldStatus := o.Status
	o.Status = newStatus
	now := time.Now()

	// Set appropriate timestamps
	switch newStatus {
	case OrderStatusConfirmed:
		if o.Timestamps.ConfirmedAt == nil {
			o.Timestamps.ConfirmedAt = &now
		}
	case OrderStatusProcessing:
		if o.Timestamps.ProcessedAt == nil {
			o.Timestamps.ProcessedAt = &now
		}
	case OrderStatusShipped:
		if o.Timestamps.ShippedAt == nil {
			o.Timestamps.ShippedAt = &now
		}
		o.FulfillmentStatus = FulfillmentStatusShipped
	case OrderStatusDelivered:
		if o.Timestamps.DeliveredAt == nil {
			o.Timestamps.DeliveredAt = &now
		}
		o.FulfillmentStatus = FulfillmentStatusDelivered
	case OrderStatusCancelled:
		if o.Timestamps.CancelledAt == nil {
			o.Timestamps.CancelledAt = &now
		}
	case OrderStatusRefunded:
		if o.Timestamps.RefundedAt == nil {
			o.Timestamps.RefundedAt = &now
		}
		o.PaymentStatus = PaymentStatusRefunded
	}

	o.UpdatedAt = now
	return nil
}

// UpdatePaymentStatus updates the payment status
func (o *Order) UpdatePaymentStatus(newStatus PaymentStatus) error {
	if !newStatus.IsValid() {
		return fmt.Errorf("invalid payment status: %s", newStatus)
	}

	o.PaymentStatus = newStatus
	o.UpdatedAt = time.Now()

	// Auto-update order status based on payment
	if newStatus == PaymentStatusPaid && o.Status == OrderStatusPending {
		return o.UpdateStatus(OrderStatusConfirmed)
	}

	return nil
}

// UpdateFulfillmentStatus updates the fulfillment status
func (o *Order) UpdateFulfillmentStatus(newStatus FulfillmentStatus) error {
	if !newStatus.IsValid() {
		return fmt.Errorf("invalid fulfillment status: %s", newStatus)
	}

	o.FulfillmentStatus = newStatus
	o.UpdatedAt = time.Now()
	return nil
}

// AddItem adds an item to the order
func (o *Order) AddItem(item OrderItem) error {
	// Generate ID for item if not set
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}

	// Calculate item total price
	itemTotal := item.UnitPrice.Mul(decimal.NewFromInt(int64(item.Quantity)))
	item.TotalPrice = itemTotal.Sub(item.DiscountAmount).Add(item.TaxAmount)

	o.Items = append(o.Items, item)
	o.ItemCount = len(o.Items)
	o.calculateTotals()
	o.UpdatedAt = time.Now()

	return nil
}

// RemoveItem removes an item from the order
func (o *Order) RemoveItem(itemID uuid.UUID) error {
	for i, item := range o.Items {
		if item.ID == itemID {
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
			o.ItemCount = len(o.Items)
			o.calculateTotals()
			o.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("item with ID %s not found", itemID)
}

// GetItem returns an item by ID
func (o *Order) GetItem(itemID uuid.UUID) (*OrderItem, error) {
	for i, item := range o.Items {
		if item.ID == itemID {
			return &o.Items[i], nil
		}
	}
	return nil, fmt.Errorf("item with ID %s not found", itemID)
}

// ApplyDiscount applies a discount to the order
func (o *Order) ApplyDiscount(discount OrderDiscount) {
	o.Discounts = append(o.Discounts, discount)
	o.calculateTotals()
	o.UpdatedAt = time.Now()
}

// SetShippingMethod sets the shipping method and cost
func (o *Order) SetShippingMethod(method ShippingMethod) {
	o.ShippingMethod = method
	o.Totals.ShippingAmount = method.Cost
	o.calculateTotals()
	o.UpdatedAt = time.Now()
}

// CanCancel checks if the order can be cancelled
func (o *Order) CanCancel() bool {
	return o.Status.CanCancel()
}

// CanRefund checks if the order can be refunded
func (o *Order) CanRefund() bool {
	return o.PaymentStatus == PaymentStatusPaid &&
		   (o.Status == OrderStatusDelivered || o.Status == OrderStatusShipped)
}

// IsCompleted checks if the order is completed
func (o *Order) IsCompleted() bool {
	return o.Status == OrderStatusDelivered
}

// IsCancelled checks if the order is cancelled
func (o *Order) IsCancelled() bool {
	return o.Status == OrderStatusCancelled
}

// GetTotalWeight returns the total weight of all items
func (o *Order) GetTotalWeight() decimal.Decimal {
	totalWeight := decimal.Zero
	for _, item := range o.Items {
		if item.Weight != nil {
			itemWeight := item.Weight.Mul(decimal.NewFromInt(int64(item.Quantity)))
			totalWeight = totalWeight.Add(itemWeight)
		}
	}
	return totalWeight
}

// GenerateSearchKeywords generates search keywords for the order
func (o *Order) GenerateSearchKeywords() []string {
	keywords := []string{
		o.OrderNumber,
		o.CustomerEmail,
		o.CustomerPhone,
		string(o.Status),
		string(o.PaymentStatus),
		string(o.FulfillmentStatus),
		o.Source,
		o.Totals.Currency,
	}

	// Add item names and SKUs
	for _, item := range o.Items {
		keywords = append(keywords, item.ProductName)
		if item.ProductSKU != "" {
			keywords = append(keywords, item.ProductSKU)
		}
	}

	// Add shipping and billing info
	keywords = append(keywords,
		o.ShippingAddress.FirstName,
		o.ShippingAddress.LastName,
		o.ShippingAddress.Company,
		o.BillingAddress.FirstName,
		o.BillingAddress.LastName,
		o.BillingAddress.Company,
	)

	// Add tags
	keywords = append(keywords, o.Tags...)

	// Remove duplicates and empty strings
	seen := make(map[string]bool)
	var unique []string
	for _, keyword := range keywords {
		cleaned := strings.ToLower(strings.TrimSpace(keyword))
		if cleaned != "" && !seen[cleaned] {
			seen[cleaned] = true
			unique = append(unique, cleaned)
		}
	}

	return unique
}

// GetOrderSummary returns a summary of order information
func (o *Order) GetOrderSummary() map[string]interface{} {
	return map[string]interface{}{
		"id":                 o.ID,
		"order_number":       o.OrderNumber,
		"status":             o.Status,
		"payment_status":     o.PaymentStatus,
		"fulfillment_status": o.FulfillmentStatus,
		"item_count":         o.ItemCount,
		"total_amount":       o.Totals.TotalAmount,
		"currency":           o.Totals.Currency,
		"customer_email":     o.CustomerEmail,
		"requires_shipping":  o.RequiresShipping,
		"can_cancel":         o.CanCancel(),
		"can_refund":         o.CanRefund(),
		"is_completed":       o.IsCompleted(),
		"is_cancelled":       o.IsCancelled(),
		"source":             o.Source,
		"created_at":         o.CreatedAt,
		"confirmed_at":       o.Timestamps.ConfirmedAt,
		"shipped_at":         o.Timestamps.ShippedAt,
		"delivered_at":       o.Timestamps.DeliveredAt,
	}
}

// MarshalJSON customizes JSON serialization
func (o *Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	return json.Marshal(&struct {
		*Alias
		CanCancel     bool                   `json:"can_cancel"`
		CanRefund     bool                   `json:"can_refund"`
		IsCompleted   bool                   `json:"is_completed"`
		IsCancelled   bool                   `json:"is_cancelled"`
		TotalWeight   decimal.Decimal        `json:"total_weight"`
		OrderSummary  map[string]interface{} `json:"order_summary"`
		SearchKeywords []string              `json:"search_keywords,omitempty"`
	}{
		Alias:         (*Alias)(o),
		CanCancel:     o.CanCancel(),
		CanRefund:     o.CanRefund(),
		IsCompleted:   o.IsCompleted(),
		IsCancelled:   o.IsCancelled(),
		TotalWeight:   o.GetTotalWeight(),
		OrderSummary:  o.GetOrderSummary(),
		SearchKeywords: o.GenerateSearchKeywords(),
	})
}