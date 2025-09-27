// Journey 3: E-Commerce & Payment API Integration Tests
// Tests all API endpoints involved in e-commerce and payment functionality

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Note: AuthenticatedUser is now defined in types.go

type Journey03EcommerceAPISuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	ctx        context.Context
	seller     *AuthenticatedUser // Arif (Indonesia)
	buyer      *AuthenticatedUser // Maria (Philippines)
}

// Note: CreateProductRequest is now defined in types.go

type ProductImage struct {
	URL    string `json:"url"`
	IsMain bool   `json:"isMain"`
	Alt    string `json:"alt,omitempty"`
}

type ShippingOption struct {
	Type          string   `json:"type"` // "standard", "express", "overnight"
	Price         int64    `json:"price"`
	EstimatedDays int      `json:"estimatedDays"`
	Regions       []string `json:"regions"`
}

type ProductInventory struct {
	Quantity  int              `json:"quantity"`
	LowStock  int              `json:"lowStock"`
	Variants  []ProductVariant `json:"variants,omitempty"`
}

type ProductVariant struct {
	Name     string `json:"name"`
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
	Price    int64  `json:"price,omitempty"` // Optional price override
}

type ProductResponse struct {
	ID              string                 `json:"id"`
	SellerID        string                 `json:"sellerId"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Price           int64                  `json:"price"`
	Currency        string                 `json:"currency"`
	Category        string                 `json:"category"`
	Tags            []string               `json:"tags,omitempty"`
	Images          []ProductImage         `json:"images,omitempty"`
	Specifications  map[string]interface{} `json:"specifications,omitempty"`
	ShippingOptions []ShippingOption       `json:"shippingOptions,omitempty"`
	Inventory       ProductInventory       `json:"inventory"`
	Status          string                 `json:"status"`
	CreatedAt       string                 `json:"createdAt"`
	UpdatedAt       string                 `json:"updatedAt"`
}

type ProductSearchResponse struct {
	Products     []ProductSummary `json:"products"`
	TotalCount   int              `json:"totalCount"`
	Page         int              `json:"page"`
	PageSize     int              `json:"pageSize"`
	Filters      SearchFilters    `json:"filters"`
}

type ProductSummary struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       int64   `json:"price"`
	Currency    string  `json:"currency"`
	PricePHP    float64 `json:"pricePHP,omitempty"`
	MainImage   string  `json:"mainImage"`
	SellerID    string  `json:"sellerId"`
	SellerName  string  `json:"sellerName"`
	Rating      float64 `json:"rating"`
	ReviewCount int     `json:"reviewCount"`
}

type SearchFilters struct {
	Region       string   `json:"region"`
	Currency     string   `json:"currency"`
	MinPrice     int64    `json:"minPrice,omitempty"`
	MaxPrice     int64    `json:"maxPrice,omitempty"`
	Categories   []string `json:"categories,omitempty"`
	Tags         []string `json:"tags,omitempty"`
}

type AddToCartRequest struct {
	ProductID  string                 `json:"productId"`
	Quantity   int                    `json:"quantity"`
	VariantSKU string                 `json:"variantSku,omitempty"`
	Options    map[string]interface{} `json:"options,omitempty"`
}

type CartResponse struct {
	ID         string      `json:"id"`
	UserID     string      `json:"userId"`
	Items      []CartItem  `json:"items"`
	TotalPrice int64       `json:"totalPrice"`
	Currency   string      `json:"currency"`
	UpdatedAt  string      `json:"updatedAt"`
}

type CartItem struct {
	ID         string                 `json:"id"`
	ProductID  string                 `json:"productId"`
	Name       string                 `json:"name"`
	Price      int64                  `json:"price"`
	Quantity   int                    `json:"quantity"`
	VariantSKU string                 `json:"variantSku,omitempty"`
	Image      string                 `json:"image"`
	Options    map[string]interface{} `json:"options,omitempty"`
}

type CheckoutRequest struct {
	OrderID        string         `json:"orderId,omitempty"` // For existing orders
	PaymentMethod  PaymentMethod  `json:"paymentMethod"`
	BillingAddress BillingAddress `json:"billingAddress"`
	ShippingAddress ShippingAddress `json:"shippingAddress,omitempty"`
	PromoCode      string         `json:"promoCode,omitempty"`
}

type PaymentMethod struct {
	Type     string `json:"type"`     // "card", "bank_transfer", "ewallet", "crypto"
	Provider string `json:"provider"` // "stripe", "promptpay", "grabpay", "gopay", etc.
	Currency string `json:"currency"`
}

type BillingAddress struct {
	FullName   string `json:"fullName"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

type ShippingAddress struct {
	FullName   string `json:"fullName"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
	Phone      string `json:"phone,omitempty"`
}

type CheckoutResponse struct {
	PaymentIntentID  string  `json:"paymentIntentId"`
	OrderID         string  `json:"orderId"`
	Currency        string  `json:"currency"`
	TotalAmountPHP  float64 `json:"totalAmountPHP"`
	TotalAmountIDR  int64   `json:"totalAmountIDR"`
	ExchangeRate    string  `json:"exchangeRate"`
	PaymentURL      string  `json:"paymentUrl,omitempty"`
	ExpiresAt       string  `json:"expiresAt"`
}

type ProcessPaymentRequest struct {
	PaymentIntentID   string                `json:"paymentIntentId"`
	OrderID          string                `json:"orderId"`
	PaymentMethod    PaymentMethodDetails  `json:"paymentMethod"`
	SavePaymentMethod bool                 `json:"savePaymentMethod,omitempty"`
}

type PaymentMethodDetails struct {
	Type         string `json:"type"`
	CardNumber   string `json:"cardNumber,omitempty"`
	ExpiryMonth  int    `json:"expiryMonth,omitempty"`
	ExpiryYear   int    `json:"expiryYear,omitempty"`
	CVV          string `json:"cvv,omitempty"`
	HolderName   string `json:"holderName,omitempty"`
	BankCode     string `json:"bankCode,omitempty"`
	AccountNumber string `json:"accountNumber,omitempty"`
	WalletID     string `json:"walletId,omitempty"`
}

type PaymentResponse struct {
	ID            string `json:"id"`
	Status        string `json:"status"` // "pending", "processing", "completed", "failed"
	TransactionID string `json:"transactionId"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	ProcessedAt   string `json:"processedAt,omitempty"`
}

type OrderResponse struct {
	ID            string          `json:"id"`
	BuyerID       string          `json:"buyerId"`
	SellerID      string          `json:"sellerId"`
	Items         []OrderItem     `json:"items"`
	TotalAmount   int64           `json:"totalAmount"`
	Currency      string          `json:"currency"`
	Status        string          `json:"status"` // "pending", "confirmed", "processing", "shipped", "delivered", "cancelled"
	PaymentStatus string          `json:"paymentStatus"` // "pending", "paid", "refunded"
	Shipping      ShippingDetails `json:"shipping"`
	CreatedAt     string          `json:"createdAt"`
	UpdatedAt     string          `json:"updatedAt"`
}

// Note: OrderItem is now defined in types.go

type ShippingDetails struct {
	Method        string `json:"method"`
	TrackingNumber string `json:"trackingNumber,omitempty"`
	Carrier       string `json:"carrier,omitempty"`
	Address       ShippingAddress `json:"address"`
	ShippedAt     string `json:"shippedAt,omitempty"`
	EstimatedDelivery string `json:"estimatedDelivery,omitempty"`
}

type CurrencyConversionResponse struct {
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
	Amount       int64   `json:"amount"`
	ConvertedAmount float64 `json:"convertedAmount"`
	ExchangeRate string  `json:"exchangeRate"`
	Timestamp    string  `json:"timestamp"`
}

func (suite *Journey03EcommerceAPISuite) SetupSuite() {
	suite.baseURL = "http://localhost:8081" // Auth Service Direct
	suite.httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
	suite.ctx = context.Background()

	// Create authenticated test users
	suite.seller = suite.createBusinessUser("arif@test.com", "ID", "id")
	suite.buyer = suite.createAuthenticatedUser("maria@test.com", "PH", "tl")
}

// Test 3.1: Product Management API
func (suite *Journey03EcommerceAPISuite) TestProductManagementAPI() {
	// Step 1: POST /api/v1/products - Create product
	productReq := CreateProductRequest{
		Name:        "Traditional Batik Scarf",
		Description: "Handmade batik scarf from Central Java, Indonesia",
		Price:       85000000, // 850,000 IDR in cents (IDR doesn't have cents, but using consistent format)
		Currency:    "IDR",
		Category:    "traditional_crafts",
		Tags:        []string{"batik", "scarf", "handmade", "java", "traditional", "indonesia"},
		Images: []ProductImage{
			{URL: "https://example.com/batik-scarf-1.jpg", IsMain: true, Alt: "Blue batik scarf front view"},
			{URL: "https://example.com/batik-scarf-2.jpg", IsMain: false, Alt: "Blue batik scarf detail"},
		},
		Specifications: map[string]interface{}{
			"material":    "100% silk",
			"dimensions":  "180cm x 60cm",
			"care":        "dry clean only",
			"origin":      "Central Java, Indonesia",
			"weight":      "120g",
		},
		ShippingOptions: []ShippingOption{
			{
				Type:          "standard",
				Price:         5000000, // 50,000 IDR
				EstimatedDays: 7,
				Regions:       []string{"ID", "PH", "SG", "TH", "MY", "VN"},
			},
			{
				Type:          "express",
				Price:         12000000, // 120,000 IDR
				EstimatedDays: 3,
				Regions:       []string{"ID", "PH", "SG"},
			},
		},
		Inventory: ProductInventory{
			Quantity: 15,
			LowStock: 3,
			Variants: []ProductVariant{
				{Name: "Blue Pattern", SKU: "BATIK-BLUE-001", Quantity: 5},
				{Name: "Red Pattern", SKU: "BATIK-RED-001", Quantity: 5},
				{Name: "Green Pattern", SKU: "BATIK-GREEN-001", Quantity: 5},
			},
		},
		Status: "active",
	}

	sellerHeaders := map[string]string{
		"Authorization":   "Bearer " + suite.seller.AccessToken,
		"Accept-Language": "id",
	}

	productResp, statusCode := suite.makeAPICall("POST", "/api/v1/products", productReq, sellerHeaders)
	assert.Equal(suite.T(), 201, statusCode, "Product creation should succeed")

	var product ProductResponse
	err := json.Unmarshal(productResp, &product)
	require.NoError(suite.T(), err, "Should parse product response")

	productID := product.ID
	assert.NotEmpty(suite.T(), productID, "Should return product ID")
	assert.Equal(suite.T(), "Traditional Batik Scarf", product.Name, "Name should match")
	assert.Equal(suite.T(), suite.seller.UserID, product.SellerID, "Seller ID should match")
	assert.Equal(suite.T(), "IDR", product.Currency, "Currency should be IDR")
	assert.Equal(suite.T(), int64(85000000), product.Price, "Price should match")

	// Step 2: GET /api/v1/products/{id} - Retrieve product
	getProductResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/products/%s", productID), nil, nil)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve product")

	var retrievedProduct ProductResponse
	err = json.Unmarshal(getProductResp, &retrievedProduct)
	require.NoError(suite.T(), err, "Should parse retrieved product")
	assert.Equal(suite.T(), productID, retrievedProduct.ID, "Product ID should match")

	// Step 3: PUT /api/v1/products/{id} - Update product
	updateReq := map[string]interface{}{
		"description": "Handmade batik scarf from Central Java, Indonesia - Updated with new patterns!",
		"tags":        []string{"batik", "scarf", "handmade", "java", "traditional", "indonesia", "updated"},
	}

	_, statusCode = suite.makeAPICall("PUT",
		fmt.Sprintf("/api/v1/products/%s", productID), updateReq, sellerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Product update should succeed")

	// Step 4: GET /api/v1/products/seller/{sellerId} - List seller products
	listProductsResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/products/seller/%s", suite.seller.UserID), nil, sellerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should list seller products")

	var sellerProducts []ProductResponse
	err = json.Unmarshal(listProductsResp, &sellerProducts)
	require.NoError(suite.T(), err, "Should parse seller products")
	assert.Greater(suite.T(), len(sellerProducts), 0, "Should have products")

	// Find our product
	found := false
	for _, p := range sellerProducts {
		if p.ID == productID {
			found = true
			assert.Contains(suite.T(), p.Description, "Updated with new patterns!", "Should have updated description")
			break
		}
	}
	assert.True(suite.T(), found, "Should find created product in seller's list")
}

// Test 3.2: Product Search and Discovery API
func (suite *Journey03EcommerceAPISuite) TestProductSearchAPI() {
	// Create test product first
	productID := suite.createTestProduct()

	// Step 1: GET /api/v1/products/search - Basic search
	searchResp, statusCode := suite.makeAPICall("GET",
		"/api/v1/products/search?q=batik&limit=10", nil, nil)
	assert.Equal(suite.T(), 200, statusCode, "Basic search should succeed")

	var searchResults ProductSearchResponse
	err := json.Unmarshal(searchResp, &searchResults)
	require.NoError(suite.T(), err, "Should parse search results")
	assert.Greater(suite.T(), searchResults.TotalCount, 0, "Should have search results")

	// Step 2: GET /api/v1/products/search - Category search
	categorySearchResp, statusCode := suite.makeAPICall("GET",
		"/api/v1/products/search?category=traditional_crafts&region=ID", nil, nil)
	assert.Equal(suite.T(), 200, statusCode, "Category search should succeed")

	var categoryResults ProductSearchResponse
	err = json.Unmarshal(categorySearchResp, &categoryResults)
	require.NoError(suite.T(), err, "Should parse category results")
	assert.Greater(suite.T(), categoryResults.TotalCount, 0, "Should have category results")

	// Step 3: GET /api/v1/products/search - Cross-region search with currency conversion
	buyerHeaders := map[string]string{
		"Authorization":        "Bearer " + suite.buyer.AccessToken,
		"Accept-Language":      "tl", // Filipino
		"X-Preferred-Currency": "PHP",
	}

	crossRegionResp, statusCode := suite.makeAPICall("GET",
		"/api/v1/products/search?q=batik&region=ID&currency=PHP", nil, buyerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Cross-region search should succeed")

	var crossRegionResults ProductSearchResponse
	err = json.Unmarshal(crossRegionResp, &crossRegionResults)
	require.NoError(suite.T(), err, "Should parse cross-region results")

	// Find our product and verify currency conversion
	var foundProduct *ProductSummary
	for _, product := range crossRegionResults.Products {
		if product.ID == productID {
			foundProduct = &product
			break
		}
	}
	require.NotNil(suite.T(), foundProduct, "Should find product in cross-region search")
	assert.Greater(suite.T(), foundProduct.PricePHP, float64(0), "Should have converted price in PHP")

	// Step 4: GET /api/v1/products/categories - List categories
	categoriesResp, statusCode := suite.makeAPICall("GET", "/api/v1/products/categories", nil, nil)
	assert.Equal(suite.T(), 200, statusCode, "Should list categories")

	var categories []map[string]interface{}
	err = json.Unmarshal(categoriesResp, &categories)
	require.NoError(suite.T(), err, "Should parse categories")
	assert.Greater(suite.T(), len(categories), 0, "Should have categories")

	// Find traditional_crafts category
	found := false
	for _, category := range categories {
		if category["name"] == "traditional_crafts" {
			found = true
			assert.Greater(suite.T(), category["productCount"], float64(0), "Should have product count")
			break
		}
	}
	assert.True(suite.T(), found, "Should find traditional_crafts category")
}

// Test 3.3: Shopping Cart Management API
func (suite *Journey03EcommerceAPISuite) TestShoppingCartAPI() {
	productID := suite.createTestProduct()

	buyerHeaders := map[string]string{
		"Authorization": "Bearer " + suite.buyer.AccessToken,
	}

	// Step 1: POST /api/v1/cart/items - Add item to cart
	addToCartReq := AddToCartRequest{
		ProductID:  productID,
		Quantity:   2,
		VariantSKU: "BATIK-BLUE-001",
		Options: map[string]interface{}{
			"gift_wrap":     true,
			"special_note":  "Please wrap carefully for international shipping",
		},
	}

	_, statusCode := suite.makeAPICall("POST", "/api/v1/cart/items", addToCartReq, buyerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Adding to cart should succeed")

	// Step 2: GET /api/v1/cart - Retrieve cart
	cartResp, statusCode := suite.makeAPICall("GET", "/api/v1/cart", nil, buyerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve cart")

	var cart CartResponse
	err := json.Unmarshal(cartResp, &cart)
	require.NoError(suite.T(), err, "Should parse cart response")

	assert.NotEmpty(suite.T(), cart.ID, "Should have cart ID")
	assert.Equal(suite.T(), suite.buyer.UserID, cart.UserID, "Cart user ID should match")
	assert.Len(suite.T(), cart.Items, 1, "Should have 1 item in cart")
	assert.Equal(suite.T(), productID, cart.Items[0].ProductID, "Product ID should match")
	assert.Equal(suite.T(), 2, cart.Items[0].Quantity, "Quantity should match")
	assert.Equal(suite.T(), "BATIK-BLUE-001", cart.Items[0].VariantSKU, "Variant SKU should match")

	cartItemID := cart.Items[0].ID

	// Step 3: PUT /api/v1/cart/items/{id} - Update cart item
	updateCartReq := map[string]interface{}{
		"quantity": 3,
		"options": map[string]interface{}{
			"gift_wrap":     true,
			"special_note":  "Please wrap very carefully - this is a special gift",
		},
	}

	_, statusCode = suite.makeAPICall("PUT",
		fmt.Sprintf("/api/v1/cart/items/%s", cartItemID), updateCartReq, buyerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Cart item update should succeed")

	// Step 4: GET /api/v1/cart - Verify update
	updatedCartResp, statusCode := suite.makeAPICall("GET", "/api/v1/cart", nil, buyerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve updated cart")

	var updatedCart CartResponse
	err = json.Unmarshal(updatedCartResp, &updatedCart)
	require.NoError(suite.T(), err, "Should parse updated cart")
	assert.Equal(suite.T(), 3, updatedCart.Items[0].Quantity, "Quantity should be updated")

	// Step 5: DELETE /api/v1/cart/items/{id} - Remove item from cart
	_, statusCode = suite.makeAPICall("DELETE",
		fmt.Sprintf("/api/v1/cart/items/%s", cartItemID), nil, buyerHeaders)
	assert.Equal(suite.T(), 204, statusCode, "Item removal should succeed")

	// Step 6: GET /api/v1/cart - Verify item removed
	finalCartResp, statusCode := suite.makeAPICall("GET", "/api/v1/cart", nil, buyerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve final cart")

	var finalCart CartResponse
	err = json.Unmarshal(finalCartResp, &finalCart)
	require.NoError(suite.T(), err, "Should parse final cart")
	assert.Len(suite.T(), finalCart.Items, 0, "Cart should be empty after removal")
}

// Test 3.4: Checkout and Payment Processing API
func (suite *Journey03EcommerceAPISuite) TestCheckoutPaymentAPI() {
	productID := suite.createTestProduct()

	// Add item to cart first
	suite.addItemToCart(productID, 1)

	buyerHeaders := map[string]string{
		"Authorization": "Bearer " + suite.buyer.AccessToken,
	}

	// Step 1: POST /api/v1/checkout - Create checkout session
	checkoutReq := CheckoutRequest{
		PaymentMethod: PaymentMethod{
			Type:     "card",
			Provider: "stripe",
			Currency: "PHP",
		},
		BillingAddress: BillingAddress{
			FullName:   "Maria Santos",
			Street:     "123 Roxas Boulevard",
			City:       "Manila",
			State:      "NCR",
			PostalCode: "1000",
			Country:    "PH",
		},
		ShippingAddress: ShippingAddress{
			FullName:   "Maria Santos",
			Street:     "123 Roxas Boulevard",
			City:       "Manila",
			State:      "NCR",
			PostalCode: "1000",
			Country:    "PH",
			Phone:      "+639171234567",
		},
	}

	checkoutResp, statusCode := suite.makeAPICall("POST", "/api/v1/checkout", checkoutReq, buyerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Checkout should succeed")

	var checkout CheckoutResponse
	err := json.Unmarshal(checkoutResp, &checkout)
	require.NoError(suite.T(), err, "Should parse checkout response")

	assert.NotEmpty(suite.T(), checkout.PaymentIntentID, "Should return payment intent ID")
	assert.NotEmpty(suite.T(), checkout.OrderID, "Should return order ID")
	assert.Equal(suite.T(), "PHP", checkout.Currency, "Currency should be PHP")
	assert.Greater(suite.T(), checkout.TotalAmountPHP, float64(0), "Should have total amount in PHP")
	assert.Greater(suite.T(), checkout.TotalAmountIDR, int64(0), "Should have total amount in IDR")
	assert.NotEmpty(suite.T(), checkout.ExchangeRate, "Should have exchange rate")

	// Step 2: POST /api/v1/payments/process - Process payment
	paymentReq := ProcessPaymentRequest{
		PaymentIntentID: checkout.PaymentIntentID,
		OrderID:        checkout.OrderID,
		PaymentMethod: PaymentMethodDetails{
			Type:        "card",
			CardNumber:  "4242424242424242", // Test card
			ExpiryMonth: 12,
			ExpiryYear:  2025,
			CVV:         "123",
			HolderName:  "Maria Santos",
		},
		SavePaymentMethod: true,
	}

	paymentResp, statusCode := suite.makeAPICall("POST", "/api/v1/payments/process", paymentReq, buyerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Payment processing should succeed")

	var payment PaymentResponse
	err = json.Unmarshal(paymentResp, &payment)
	require.NoError(suite.T(), err, "Should parse payment response")

	assert.NotEmpty(suite.T(), payment.ID, "Should return payment ID")
	assert.Equal(suite.T(), "completed", payment.Status, "Payment status should be completed")
	assert.NotEmpty(suite.T(), payment.TransactionID, "Should return transaction ID")
	assert.Greater(suite.T(), payment.Amount, int64(0), "Should have payment amount")

	// Step 3: GET /api/v1/orders/{id} - Verify order status
	orderResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/orders/%s", checkout.OrderID), nil, buyerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve order")

	var order OrderResponse
	err = json.Unmarshal(orderResp, &order)
	require.NoError(suite.T(), err, "Should parse order response")

	assert.Equal(suite.T(), checkout.OrderID, order.ID, "Order ID should match")
	assert.Equal(suite.T(), suite.buyer.UserID, order.BuyerID, "Buyer ID should match")
	assert.Equal(suite.T(), suite.seller.UserID, order.SellerID, "Seller ID should match")
	assert.Equal(suite.T(), "confirmed", order.Status, "Order status should be confirmed")
	assert.Equal(suite.T(), "paid", order.PaymentStatus, "Payment status should be paid")
	assert.Len(suite.T(), order.Items, 1, "Should have 1 item")
}

// Test 3.5: Multi-Currency and Cross-Border Features API
func (suite *Journey03EcommerceAPISuite) TestMultiCurrencyAPI() {
	// Step 1: GET /api/v1/commerce/currencies - List supported currencies
	currenciesResp, statusCode := suite.makeAPICall("GET", "/api/v1/commerce/currencies", nil, nil)
	assert.Equal(suite.T(), 200, statusCode, "Should list currencies")

	var currencies []map[string]interface{}
	err := json.Unmarshal(currenciesResp, &currencies)
	require.NoError(suite.T(), err, "Should parse currencies")

	expectedCurrencies := []string{"THB", "SGD", "IDR", "MYR", "PHP", "VND"}
	for _, expectedCurrency := range expectedCurrencies {
		found := false
		for _, currency := range currencies {
			if currency["code"] == expectedCurrency {
				found = true
				assert.Greater(suite.T(), currency["exchangeRate"], float64(0), "Should have exchange rate")
				break
			}
		}
		assert.True(suite.T(), found, fmt.Sprintf("Should support currency: %s", expectedCurrency))
	}

	// Step 2: POST /api/v1/commerce/currency/convert - Currency conversion
	conversionReq := map[string]interface{}{
		"fromCurrency": "IDR",
		"toCurrency":   "PHP",
		"amount":       850000, // 850,000 IDR
	}

	conversionResp, statusCode := suite.makeAPICall("POST", "/api/v1/commerce/currency/convert", conversionReq, nil)
	assert.Equal(suite.T(), 200, statusCode, "Currency conversion should succeed")

	var conversion CurrencyConversionResponse
	err = json.Unmarshal(conversionResp, &conversion)
	require.NoError(suite.T(), err, "Should parse conversion response")

	assert.Equal(suite.T(), "IDR", conversion.FromCurrency, "From currency should match")
	assert.Equal(suite.T(), "PHP", conversion.ToCurrency, "To currency should match")
	assert.Equal(suite.T(), int64(850000), conversion.Amount, "Amount should match")
	assert.Greater(suite.T(), conversion.ConvertedAmount, float64(0), "Should have converted amount")
	assert.NotEmpty(suite.T(), conversion.ExchangeRate, "Should have exchange rate")

	// Step 3: Regional payment methods testing
	regions := []struct {
		country        string
		currency       string
		paymentMethods []string
	}{
		{"TH", "THB", []string{"promptpay", "card", "bank_transfer"}},
		{"ID", "IDR", []string{"gopay", "ovo", "dana", "card", "bank_transfer"}},
		{"PH", "PHP", []string{"gcash", "paymaya", "card", "bank_transfer"}},
		{"SG", "SGD", []string{"paynow", "grabpay", "card", "bank_transfer"}},
		{"MY", "MYR", []string{"tng", "grabpay", "card", "bank_transfer"}},
		{"VN", "VND", []string{"momo", "zalopay", "card", "bank_transfer"}},
	}

	for _, region := range regions {
		suite.T().Run(fmt.Sprintf("PaymentMethods_%s", region.country), func(t *testing.T) {
			paymentMethodsResp, statusCode := suite.makeAPICall("GET",
				fmt.Sprintf("/api/v1/payments/methods?country=%s", region.country), nil, nil)
			assert.Equal(t, 200, statusCode, fmt.Sprintf("Should get payment methods for %s", region.country))

			var methods []map[string]interface{}
			err := json.Unmarshal(paymentMethodsResp, &methods)
			require.NoError(t, err, "Should parse payment methods")

			assert.Greater(t, len(methods), 0, "Should have payment methods")

			// Verify expected methods are available
			for _, expectedMethod := range region.paymentMethods {
				found := false
				for _, method := range methods {
					if method["type"] == expectedMethod {
						found = true
						assert.True(t, method["enabled"].(bool), fmt.Sprintf("%s should be enabled", expectedMethod))
						break
					}
				}
				assert.True(t, found, fmt.Sprintf("Should support payment method: %s", expectedMethod))
			}
		})
	}
}

// Test 3.6: Order Management and Fulfillment API
func (suite *Journey03EcommerceAPISuite) TestOrderManagementAPI() {
	// Create and pay for order first
	orderID := suite.createPaidOrder()

	sellerHeaders := map[string]string{
		"Authorization": "Bearer " + suite.seller.AccessToken,
	}

	buyerHeaders := map[string]string{
		"Authorization": "Bearer " + suite.buyer.AccessToken,
	}

	// Step 1: GET /api/v1/orders/seller/{sellerId} - List seller orders
	sellerOrdersResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/orders/seller/%s", suite.seller.UserID), nil, sellerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should list seller orders")

	var sellerOrders []OrderResponse
	err := json.Unmarshal(sellerOrdersResp, &sellerOrders)
	require.NoError(suite.T(), err, "Should parse seller orders")
	assert.Greater(suite.T(), len(sellerOrders), 0, "Seller should have orders")

	// Step 2: PUT /api/v1/orders/{id}/status - Update order status to processing
	updateStatusReq := map[string]interface{}{
		"status":    "processing",
		"notes":     "Order is being prepared for shipment",
		"updatedBy": suite.seller.UserID,
	}

	_, statusCode = suite.makeAPICall("PUT",
		fmt.Sprintf("/api/v1/orders/%s/status", orderID), updateStatusReq, sellerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Order status update should succeed")

	// Step 3: PUT /api/v1/orders/{id}/shipping - Update shipping information
	shippingUpdateReq := map[string]interface{}{
		"status":           "shipped",
		"trackingNumber":   "JNE123456789",
		"carrier":          "JNE Express",
		"shippedAt":        time.Now().Format(time.RFC3339),
		"estimatedDelivery": time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
		"notes":            "Package has been dispatched from Jakarta distribution center",
	}

	_, statusCode = suite.makeAPICall("PUT",
		fmt.Sprintf("/api/v1/orders/%s/shipping", orderID), shippingUpdateReq, sellerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Shipping update should succeed")

	// Step 4: GET /api/v1/orders/{id} - Verify shipping update
	updatedOrderResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/orders/%s", orderID), nil, buyerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve updated order")

	var updatedOrder OrderResponse
	err = json.Unmarshal(updatedOrderResp, &updatedOrder)
	require.NoError(suite.T(), err, "Should parse updated order")

	assert.Equal(suite.T(), "shipped", updatedOrder.Status, "Order status should be shipped")
	assert.Equal(suite.T(), "JNE123456789", updatedOrder.Shipping.TrackingNumber, "Tracking number should match")
	assert.Equal(suite.T(), "JNE Express", updatedOrder.Shipping.Carrier, "Carrier should match")

	// Step 5: GET /api/v1/orders/{id}/tracking - Track order
	trackingResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/orders/%s/tracking", orderID), nil, buyerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should get tracking information")

	var trackingInfo map[string]interface{}
	err = json.Unmarshal(trackingResp, &trackingInfo)
	require.NoError(suite.T(), err, "Should parse tracking info")

	assert.Equal(suite.T(), "JNE123456789", trackingInfo["trackingNumber"], "Tracking number should match")
	assert.Equal(suite.T(), "shipped", trackingInfo["status"], "Tracking status should be shipped")

	// Step 6: GET /api/v1/orders/buyer/{buyerId} - List buyer orders
	buyerOrdersResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/orders/buyer/%s", suite.buyer.UserID), nil, buyerHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should list buyer orders")

	var buyerOrders []OrderResponse
	err = json.Unmarshal(buyerOrdersResp, &buyerOrders)
	require.NoError(suite.T(), err, "Should parse buyer orders")
	assert.Greater(suite.T(), len(buyerOrders), 0, "Buyer should have orders")

	// Find our order
	found := false
	for _, order := range buyerOrders {
		if order.ID == orderID {
			found = true
			assert.Equal(suite.T(), "shipped", order.Status, "Order should be shipped")
			break
		}
	}
	assert.True(suite.T(), found, "Should find order in buyer's order list")
}

// Helper methods
func (suite *Journey03EcommerceAPISuite) makeAPICall(method, endpoint string, body interface{}, headers map[string]string) ([]byte, int) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(suite.ctx, method, suite.baseURL+endpoint, reqBody)
	require.NoError(suite.T(), err)

	req.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(suite.T(), err)

	return respBody, resp.StatusCode
}

func (suite *Journey03EcommerceAPISuite) createAuthenticatedUser(email, country, language string) *AuthenticatedUser {
	regReq := map[string]interface{}{
		"email":     email,
		"password":  "SecurePass123!",
		"firstName": "Test",
		"lastName":  "User",
		"country":   country,
		"language":  language,
	}

	regResp, statusCode := suite.makeAPICall("POST", "/api/v1/auth/register", regReq, nil)
	require.Equal(suite.T(), 201, statusCode)

	var regResult map[string]interface{}
	err := json.Unmarshal(regResp, &regResult)
	require.NoError(suite.T(), err)

	verifyReq := map[string]string{
		"userId": regResult["userId"].(string),
		"code":   regResult["verifyCode"].(string),
	}

	verifyResp, statusCode := suite.makeAPICall("POST", "/api/v1/auth/verify", verifyReq, nil)
	require.Equal(suite.T(), 200, statusCode)

	var verifyResult map[string]interface{}
	err = json.Unmarshal(verifyResp, &verifyResult)
	require.NoError(suite.T(), err)

	return &AuthenticatedUser{
		UserID:       regResult["userId"].(string),
		Email:        email,
		AccessToken:  verifyResult["accessToken"].(string),
		RefreshToken: verifyResult["refreshToken"].(string),
		Country:      country,
		Language:     language,
	}
}

func (suite *Journey03EcommerceAPISuite) createBusinessUser(email, country, language string) *AuthenticatedUser {
	user := suite.createAuthenticatedUser(email, country, language)

	// Upgrade to business account
	businessReq := map[string]interface{}{
		"accountType":   "business",
		"businessName":  "Traditional Crafts Indonesia",
		"businessType":  "retail",
		"taxId":         "ID123456789",
	}

	headers := map[string]string{
		"Authorization": "Bearer " + user.AccessToken,
	}

	_, statusCode := suite.makeAPICall("PUT", fmt.Sprintf("/api/v1/users/%s/business", user.UserID), businessReq, headers)
	require.Equal(suite.T(), 200, statusCode, "Business upgrade should succeed")

	return user
}

func (suite *Journey03EcommerceAPISuite) createTestProduct() string {
	productReq := CreateProductRequest{
		Name:        "Test Batik Scarf",
		Description: "Test product for integration testing",
		Price:       85000000, // 850,000 IDR
		Currency:    "IDR",
		Category:    "traditional_crafts",
		Tags:        []string{"batik", "test"},
		Inventory: ProductInventory{
			Quantity: 10,
			LowStock: 2,
		},
		Status: "active",
	}

	headers := map[string]string{
		"Authorization": "Bearer " + suite.seller.AccessToken,
	}

	productResp, statusCode := suite.makeAPICall("POST", "/api/v1/products", productReq, headers)
	require.Equal(suite.T(), 201, statusCode)

	var product ProductResponse
	err := json.Unmarshal(productResp, &product)
	require.NoError(suite.T(), err)

	return product.ID
}

func (suite *Journey03EcommerceAPISuite) addItemToCart(productID string, quantity int) {
	addToCartReq := AddToCartRequest{
		ProductID: productID,
		Quantity:  quantity,
	}

	headers := map[string]string{
		"Authorization": "Bearer " + suite.buyer.AccessToken,
	}

	_, statusCode := suite.makeAPICall("POST", "/api/v1/cart/items", addToCartReq, headers)
	require.Equal(suite.T(), 200, statusCode)
}

func (suite *Journey03EcommerceAPISuite) createPaidOrder() string {
	productID := suite.createTestProduct()
	suite.addItemToCart(productID, 1)

	checkoutReq := CheckoutRequest{
		PaymentMethod: PaymentMethod{
			Type:     "card",
			Provider: "stripe",
			Currency: "PHP",
		},
		BillingAddress: BillingAddress{
			FullName:   "Test Buyer",
			Street:     "123 Test Street",
			City:       "Test City",
			PostalCode: "1000",
			Country:    "PH",
		},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + suite.buyer.AccessToken,
	}

	checkoutResp, statusCode := suite.makeAPICall("POST", "/api/v1/checkout", checkoutReq, headers)
	require.Equal(suite.T(), 200, statusCode)

	var checkout CheckoutResponse
	err := json.Unmarshal(checkoutResp, &checkout)
	require.NoError(suite.T(), err)

	paymentReq := ProcessPaymentRequest{
		PaymentIntentID: checkout.PaymentIntentID,
		OrderID:        checkout.OrderID,
		PaymentMethod: PaymentMethodDetails{
			Type:       "card",
			CardNumber: "4242424242424242",
			ExpiryMonth: 12,
			ExpiryYear: 2025,
			CVV:        "123",
			HolderName: "Test Buyer",
		},
	}

	_, statusCode = suite.makeAPICall("POST", "/api/v1/payments/process", paymentReq, headers)
	require.Equal(suite.T(), 200, statusCode)

	return checkout.OrderID
}

func TestJourney03EcommerceAPISuite(t *testing.T) {
	suite.Run(t, new(Journey03EcommerceAPISuite))
}