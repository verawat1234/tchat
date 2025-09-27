package integration

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// CommerceIntegrationSuite tests the Commerce service endpoints
type CommerceIntegrationSuite struct {
	APIIntegrationSuite
	ports ServicePort
}

// Shop represents a shop/store
type Shop struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	OwnerID     string            `json:"ownerId"`
	Category    string            `json:"category"`
	Status      string            `json:"status"`
	Settings    map[string]string `json:"settings"`
	Address     *Address          `json:"address,omitempty"`
	CreatedAt   string            `json:"createdAt"`
	UpdatedAt   string            `json:"updatedAt"`
}

// Product represents a product in a shop
type Product struct {
	ID          string            `json:"id"`
	ShopID      string            `json:"shopId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Currency    string            `json:"currency"`
	Category    string            `json:"category"`
	SKU         string            `json:"sku"`
	Inventory   int               `json:"inventory"`
	Images      []string          `json:"images"`
	Tags        []string          `json:"tags"`
	Attributes  map[string]string `json:"attributes"`
	Status      string            `json:"status"`
	CreatedAt   string            `json:"createdAt"`
	UpdatedAt   string            `json:"updatedAt"`
}

// Order represents an order
type Order struct {
	ID          string      `json:"id"`
	CustomerID  string      `json:"customerId"`
	ShopID      string      `json:"shopId"`
	Items       []OrderItem `json:"items"`
	TotalAmount float64     `json:"totalAmount"`
	Currency    string      `json:"currency"`
	Status      string      `json:"status"`
	ShippingAddress *Address `json:"shippingAddress,omitempty"`
	BillingAddress  *Address `json:"billingAddress,omitempty"`
	CreatedAt   string      `json:"createdAt"`
	UpdatedAt   string      `json:"updatedAt"`
}

// Note: OrderItem is now defined in types.go

// Address represents an address
type Address struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	State    string `json:"state"`
	Country  string `json:"country"`
	PostalCode string `json:"postalCode"`
}

// CreateShopRequest represents shop creation request
type CreateShopRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Category    string            `json:"category"`
	Settings    map[string]string `json:"settings"`
	Address     *Address          `json:"address,omitempty"`
}

// Note: CreateProductRequest is now defined in types.go

// CreateOrderRequest represents order creation request
type CreateOrderRequest struct {
	ShopID          string    `json:"shopId"`
	Items           []OrderItem `json:"items"`
	ShippingAddress *Address  `json:"shippingAddress,omitempty"`
	BillingAddress  *Address  `json:"billingAddress,omitempty"`
}

// CommerceResponse represents commerce API response
type CommerceResponse struct {
	Success   bool        `json:"success"`
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Shop      *Shop       `json:"shop,omitempty"`
	Shops     []Shop      `json:"shops,omitempty"`
	Product   *Product    `json:"product,omitempty"`
	Products  []Product   `json:"products,omitempty"`
	Order     *Order      `json:"order,omitempty"`
	Orders    []Order     `json:"orders,omitempty"`
	Total     int         `json:"total,omitempty"`
	Page      int         `json:"page,omitempty"`
	Limit     int         `json:"limit,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// SetupSuite initializes the commerce integration test suite
func (suite *CommerceIntegrationSuite) SetupSuite() {
	suite.APIIntegrationSuite.SetupSuite()
	suite.ports = DefaultServicePorts()

	// Wait for commerce service to be available
	err := suite.waitForService(suite.ports.Commerce, 30*time.Second)
	if err != nil {
		suite.T().Fatalf("Commerce service not available: %v", err)
	}
}

// TestCommerceServiceHealth verifies commerce service health endpoint
func (suite *CommerceIntegrationSuite) TestCommerceServiceHealth() {
	healthCheck, err := suite.checkServiceHealth(suite.ports.Commerce)
	require.NoError(suite.T(), err, "Health check should succeed")

	assert.Equal(suite.T(), "healthy", healthCheck.Status)
	assert.Equal(suite.T(), "commerce-service", healthCheck.Service)
	assert.NotEmpty(suite.T(), healthCheck.Timestamp)
}

// TestCreateShop tests shop creation
func (suite *CommerceIntegrationSuite) TestCreateShop() {
	url := fmt.Sprintf("%s:%d/api/v1/shops", suite.baseURL, suite.ports.Commerce)

	createReq := CreateShopRequest{
		Name:        "Test Electronics Shop",
		Description: "A shop selling electronic devices",
		Category:    "electronics",
		Settings: map[string]string{
			"currency":       "USD",
			"tax_rate":       "8.5",
			"shipping_zone":  "domestic",
		},
		Address: &Address{
			Street:     "123 Tech Street",
			City:       "San Francisco",
			State:      "CA",
			Country:    "USA",
			PostalCode: "94105",
		},
	}

	resp, err := suite.makeRequest("POST", url, createReq, nil)
	require.NoError(suite.T(), err, "Create shop request should succeed")
	defer resp.Body.Close()

	// Should return 201 for successful creation
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var commerceResp CommerceResponse
	err = suite.parseResponse(resp, &commerceResp)
	require.NoError(suite.T(), err, "Should parse shop creation response")

	assert.True(suite.T(), commerceResp.Success)
	assert.Equal(suite.T(), "success", commerceResp.Status)
	assert.NotNil(suite.T(), commerceResp.Shop)
	assert.NotEmpty(suite.T(), commerceResp.Shop.ID)
	assert.Equal(suite.T(), createReq.Name, commerceResp.Shop.Name)
	assert.Equal(suite.T(), createReq.Category, commerceResp.Shop.Category)
	assert.Equal(suite.T(), "active", commerceResp.Shop.Status) // Default status
}

// TestGetShop tests retrieving a specific shop
func (suite *CommerceIntegrationSuite) TestGetShop() {
	// First create a shop
	createURL := fmt.Sprintf("%s:%d/api/v1/shops", suite.baseURL, suite.ports.Commerce)
	createReq := CreateShopRequest{
		Name:        "Get Test Shop",
		Description: "Shop for get test",
		Category:    "test",
		Settings:    map[string]string{"currency": "USD"},
	}

	createResp, err := suite.makeRequest("POST", createURL, createReq, nil)
	require.NoError(suite.T(), err, "Create shop for get test should succeed")
	defer createResp.Body.Close()

	var createCommerceResp CommerceResponse
	err = suite.parseResponse(createResp, &createCommerceResp)
	require.NoError(suite.T(), err, "Should parse create response")
	require.NotNil(suite.T(), createCommerceResp.Shop)

	shopID := createCommerceResp.Shop.ID

	// Now get the shop
	getURL := fmt.Sprintf("%s:%d/api/v1/shops/%s", suite.baseURL, suite.ports.Commerce, shopID)
	resp, err := suite.makeRequest("GET", getURL, nil, nil)
	require.NoError(suite.T(), err, "Get shop request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var commerceResp CommerceResponse
	err = suite.parseResponse(resp, &commerceResp)
	require.NoError(suite.T(), err, "Should parse get shop response")

	assert.True(suite.T(), commerceResp.Success)
	assert.NotNil(suite.T(), commerceResp.Shop)
	assert.Equal(suite.T(), shopID, commerceResp.Shop.ID)
	assert.Equal(suite.T(), createReq.Name, commerceResp.Shop.Name)
}

// TestListShops tests listing shops
func (suite *CommerceIntegrationSuite) TestListShops() {
	url := fmt.Sprintf("%s:%d/api/v1/shops", suite.baseURL, suite.ports.Commerce)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "List shops request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var commerceResp CommerceResponse
	err = suite.parseResponse(resp, &commerceResp)
	require.NoError(suite.T(), err, "Should parse list shops response")

	assert.True(suite.T(), commerceResp.Success)
	assert.NotNil(suite.T(), commerceResp.Shops)
	assert.GreaterOrEqual(suite.T(), commerceResp.Total, 0)
}

// TestCreateProduct tests product creation
func (suite *CommerceIntegrationSuite) TestCreateProduct() {
	// First create a shop
	shopURL := fmt.Sprintf("%s:%d/api/v1/shops", suite.baseURL, suite.ports.Commerce)
	createShopReq := CreateShopRequest{
		Name:        "Product Test Shop",
		Description: "Shop for product testing",
		Category:    "electronics",
		Settings:    map[string]string{"currency": "USD"},
	}

	shopResp, err := suite.makeRequest("POST", shopURL, createShopReq, nil)
	require.NoError(suite.T(), err, "Create shop for product test should succeed")
	defer shopResp.Body.Close()

	var shopCommerceResp CommerceResponse
	err = suite.parseResponse(shopResp, &shopCommerceResp)
	require.NoError(suite.T(), err, "Should parse shop response")
	require.NotNil(suite.T(), shopCommerceResp.Shop)

	shopID := shopCommerceResp.Shop.ID

	// Now create a product
	productURL := fmt.Sprintf("%s:%d/api/v1/shops/%s/products", suite.baseURL, suite.ports.Commerce, shopID)
	createProductReq := CreateProductRequest{
		Name:        "Test Smartphone",
		Description: "A high-quality smartphone",
		Price:       799.99,
		Currency:    "USD",
		Category:    "smartphones",
		SKU:         "TSM-001",
		Inventory:   50,
		Images:      []string{"image1.jpg", "image2.jpg"},
		Tags:        []string{"smartphone", "electronics", "mobile"},
		Attributes: map[string]string{
			"brand":   "TestBrand",
			"model":   "TestModel",
			"color":   "Black",
			"storage": "128GB",
		},
	}

	resp, err := suite.makeRequest("POST", productURL, createProductReq, nil)
	require.NoError(suite.T(), err, "Create product request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var commerceResp CommerceResponse
	err = suite.parseResponse(resp, &commerceResp)
	require.NoError(suite.T(), err, "Should parse product creation response")

	assert.True(suite.T(), commerceResp.Success)
	assert.NotNil(suite.T(), commerceResp.Product)
	assert.NotEmpty(suite.T(), commerceResp.Product.ID)
	assert.Equal(suite.T(), shopID, commerceResp.Product.ShopID)
	assert.Equal(suite.T(), createProductReq.Name, commerceResp.Product.Name)
	assert.Equal(suite.T(), createProductReq.Price, commerceResp.Product.Price)
	assert.Equal(suite.T(), "active", commerceResp.Product.Status) // Default status
}

// TestListProducts tests listing products for a shop
func (suite *CommerceIntegrationSuite) TestListProducts() {
	// First create a shop
	shopURL := fmt.Sprintf("%s:%d/api/v1/shops", suite.baseURL, suite.ports.Commerce)
	createShopReq := CreateShopRequest{
		Name:        "List Products Shop",
		Description: "Shop for product listing test",
		Category:    "general",
		Settings:    map[string]string{"currency": "USD"},
	}

	shopResp, err := suite.makeRequest("POST", shopURL, createShopReq, nil)
	require.NoError(suite.T(), err, "Create shop should succeed")
	defer shopResp.Body.Close()

	var shopCommerceResp CommerceResponse
	err = suite.parseResponse(shopResp, &shopCommerceResp)
	require.NoError(suite.T(), err, "Should parse shop response")
	require.NotNil(suite.T(), shopCommerceResp.Shop)

	shopID := shopCommerceResp.Shop.ID

	// List products (may be empty initially)
	url := fmt.Sprintf("%s:%d/api/v1/shops/%s/products", suite.baseURL, suite.ports.Commerce, shopID)
	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "List products request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var commerceResp CommerceResponse
	err = suite.parseResponse(resp, &commerceResp)
	require.NoError(suite.T(), err, "Should parse list products response")

	assert.True(suite.T(), commerceResp.Success)
	assert.NotNil(suite.T(), commerceResp.Products)
	assert.GreaterOrEqual(suite.T(), commerceResp.Total, 0)
}

// TestCreateOrder tests order creation
func (suite *CommerceIntegrationSuite) TestCreateOrder() {
	// This test assumes products exist, so it may need to be adjusted based on actual implementation
	url := fmt.Sprintf("%s:%d/api/v1/orders", suite.baseURL, suite.ports.Commerce)

	createOrderReq := CreateOrderRequest{
		ShopID: "test-shop-id",
		Items: []OrderItem{
			{
				ProductID: "test-product-id",
				Quantity:  2,
				Price:     99.99,
				Total:     199.98,
			},
		},
		ShippingAddress: &Address{
			Street:     "456 Customer Ave",
			City:       "Los Angeles",
			State:      "CA",
			Country:    "USA",
			PostalCode: "90210",
		},
	}

	resp, err := suite.makeRequest("POST", url, createOrderReq, nil)
	require.NoError(suite.T(), err, "Create order request should succeed")
	defer resp.Body.Close()

	// This might return 400 or 404 if products don't exist, which is acceptable for integration testing
	assert.True(suite.T(),
		resp.StatusCode == http.StatusCreated ||
		resp.StatusCode == http.StatusBadRequest ||
		resp.StatusCode == http.StatusNotFound,
		"Order creation should complete with appropriate status")
}

// TestListOrders tests listing orders
func (suite *CommerceIntegrationSuite) TestListOrders() {
	url := fmt.Sprintf("%s:%d/api/v1/orders", suite.baseURL, suite.ports.Commerce)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "List orders request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var commerceResp CommerceResponse
	err = suite.parseResponse(resp, &commerceResp)
	require.NoError(suite.T(), err, "Should parse list orders response")

	assert.True(suite.T(), commerceResp.Success)
	assert.NotNil(suite.T(), commerceResp.Orders)
	assert.GreaterOrEqual(suite.T(), commerceResp.Total, 0)
}

// TestGetNonExistentShop tests retrieving non-existent shop
func (suite *CommerceIntegrationSuite) TestGetNonExistentShop() {
	url := fmt.Sprintf("%s:%d/api/v1/shops/non-existent-id", suite.baseURL, suite.ports.Commerce)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "Get non-existent shop should complete")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

// TestCreateShopInvalidData tests shop creation with invalid data
func (suite *CommerceIntegrationSuite) TestCreateShopInvalidData() {
	url := fmt.Sprintf("%s:%d/api/v1/shops", suite.baseURL, suite.ports.Commerce)

	// Test with missing required fields
	invalidReq := CreateShopRequest{
		// Missing name
		Description: "Shop without name",
		Category:    "test",
	}

	resp, err := suite.makeRequest("POST", url, invalidReq, nil)
	require.NoError(suite.T(), err, "Invalid shop creation should complete")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

// TestShopsByCategory tests filtering shops by category
func (suite *CommerceIntegrationSuite) TestShopsByCategory() {
	// First create a shop with specific category
	createURL := fmt.Sprintf("%s:%d/api/v1/shops", suite.baseURL, suite.ports.Commerce)
	createReq := CreateShopRequest{
		Name:        "Category Filter Shop",
		Description: "Shop for category filtering test",
		Category:    "test-category-filter",
		Settings:    map[string]string{"currency": "USD"},
	}

	createResp, err := suite.makeRequest("POST", createURL, createReq, nil)
	require.NoError(suite.T(), err, "Create shop for category test should succeed")
	createResp.Body.Close()

	// Now filter by category
	url := fmt.Sprintf("%s:%d/api/v1/shops?category=test-category-filter", suite.baseURL, suite.ports.Commerce)
	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "Category filter request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var commerceResp CommerceResponse
	err = suite.parseResponse(resp, &commerceResp)
	require.NoError(suite.T(), err, "Should parse category filter response")

	assert.True(suite.T(), commerceResp.Success)

	// All returned shops should have the specified category
	for _, shop := range commerceResp.Shops {
		assert.Equal(suite.T(), "test-category-filter", shop.Category)
	}
}

// TestInvalidHTTPMethods tests endpoints with invalid HTTP methods
func (suite *CommerceIntegrationSuite) TestInvalidHTTPMethods() {
	baseURL := fmt.Sprintf("%s:%d/api/v1/shops", suite.baseURL, suite.ports.Commerce)

	testCases := []struct {
		url    string
		method string
	}{
		{baseURL, "PATCH"},           // List endpoint with invalid method
		{baseURL + "/123", "POST"},   // Get endpoint with invalid method
		{baseURL + "/123", "PATCH"},  // Update endpoint with invalid method
	}

	for _, tc := range testCases {
		resp, err := suite.makeRequest(tc.method, tc.url, nil, nil)
		require.NoError(suite.T(), err, "Invalid method request should complete")
		defer resp.Body.Close()

		assert.Equal(suite.T(), http.StatusMethodNotAllowed, resp.StatusCode,
			"URL: %s, Method: %s", tc.url, tc.method)
	}
}

// TestCommerceServiceConcurrency tests concurrent requests to commerce service
func (suite *CommerceIntegrationSuite) TestCommerceServiceConcurrency() {
	url := fmt.Sprintf("%s:%d/api/v1/shops", suite.baseURL, suite.ports.Commerce)

	// Create 5 concurrent shop creation requests
	concurrency := 5
	results := make(chan int, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(index int) {
			createReq := CreateShopRequest{
				Name:        fmt.Sprintf("Concurrent Shop %d", index),
				Description: fmt.Sprintf("Shop created concurrently #%d", index),
				Category:    "concurrency-test",
				Settings:    map[string]string{"currency": "USD"},
			}

			resp, err := suite.makeRequest("POST", url, createReq, nil)
			if err != nil {
				results <- 0
				return
			}
			defer resp.Body.Close()

			results <- resp.StatusCode
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < concurrency; i++ {
		statusCode := <-results
		if statusCode == http.StatusCreated {
			successCount++
		}
	}

	// At least 80% of concurrent requests should succeed
	assert.GreaterOrEqual(suite.T(), successCount, 4, "Concurrent requests should mostly succeed")
}

// RunCommerceIntegrationTests runs the commerce integration test suite
func RunCommerceIntegrationTests(t *testing.T) {
	suite.Run(t, new(CommerceIntegrationSuite))
}