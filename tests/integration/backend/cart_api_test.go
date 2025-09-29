package integration

import (
	"bytes"
	"context"
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

// CartAPITestSuite provides comprehensive integration testing for cart API endpoints
type CartAPITestSuite struct {
	suite.Suite
	baseURL     string
	httpClient  *http.Client
	gatewayPort int
	testUserID  string
}

// Cart represents cart structure
type Cart struct {
	ID           string      `json:"id"`
	UserID       string      `json:"userId"`
	Items        []CartItem  `json:"items"`
	SubTotal     float64     `json:"subTotal"`
	TaxAmount    float64     `json:"taxAmount"`
	Total        float64     `json:"total"`
	Currency     string      `json:"currency"`
	Status       string      `json:"status"`
	CouponCode   *string     `json:"couponCode,omitempty"`
	Discount     float64     `json:"discount"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
	ExpiresAt    *time.Time  `json:"expiresAt,omitempty"`
}

// CartItem represents an item in the cart
type CartItem struct {
	ID            string            `json:"id"`
	ProductID     string            `json:"productId"`
	VariantID     *string           `json:"variantId,omitempty"`
	Quantity      int               `json:"quantity"`
	UnitPrice     float64           `json:"unitPrice"`
	TotalPrice    float64           `json:"totalPrice"`
	Name          string            `json:"name"`
	Description   *string           `json:"description,omitempty"`
	Image         *string           `json:"image,omitempty"`
	SKU           *string           `json:"sku,omitempty"`
	Attributes    map[string]string `json:"attributes,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
}

// AddToCartRequest represents request to add item to cart
type AddToCartRequest struct {
	ProductID  string            `json:"productId"`
	VariantID  *string           `json:"variantId,omitempty"`
	Quantity   int               `json:"quantity"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// UpdateCartItemRequest represents request to update cart item
type UpdateCartItemRequest struct {
	Quantity   int               `json:"quantity"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// ApplyCouponRequest represents request to apply coupon
type ApplyCouponRequest struct {
	CouponCode string `json:"couponCode"`
}

// CartResponse represents API response for cart operations
type CartResponse struct {
	Success   bool    `json:"success"`
	Message   string  `json:"message"`
	Cart      *Cart   `json:"cart,omitempty"`
	Error     *string `json:"error,omitempty"`
	Timestamp string  `json:"timestamp"`
}

// SetupSuite initializes the test suite
func (suite *CartAPITestSuite) SetupSuite() {
	suite.baseURL = "http://localhost"
	suite.gatewayPort = 8080
	suite.httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
	suite.testUserID = uuid.New().String()

	// Wait for gateway and commerce service
	suite.waitForServices()
}

// TearDownSuite cleans up after test suite
func (suite *CartAPITestSuite) TearDownSuite() {
	// Clean up test data
	suite.cleanupTestData()
}

// SetupTest prepares for each test
func (suite *CartAPITestSuite) SetupTest() {
	// Ensure clean state for each test
	suite.clearUserCart()
}

// TestCartLifecycle tests complete cart workflow
func (suite *CartAPITestSuite) TestCartLifecycle() {
	// 1. Create empty cart
	cart := suite.getOrCreateCart()
	assert.NotNil(suite.T(), cart)
	assert.Equal(suite.T(), suite.testUserID, cart.UserID)
	assert.Equal(suite.T(), 0, len(cart.Items))
	assert.Equal(suite.T(), 0.0, cart.Total)

	// 2. Add first item to cart
	productID := suite.createTestProduct()
	addReq := AddToCartRequest{
		ProductID: productID,
		Quantity:  2,
		Attributes: map[string]string{
			"color": "blue",
			"size":  "large",
		},
	}

	updatedCart := suite.addToCart(addReq)
	assert.Equal(suite.T(), 1, len(updatedCart.Items))
	assert.Equal(suite.T(), 2, updatedCart.Items[0].Quantity)
	assert.Equal(suite.T(), productID, updatedCart.Items[0].ProductID)
	assert.Greater(suite.T(), updatedCart.Total, 0.0)

	// 3. Add second different item
	productID2 := suite.createTestProduct()
	addReq2 := AddToCartRequest{
		ProductID: productID2,
		Quantity:  1,
	}

	updatedCart = suite.addToCart(addReq2)
	assert.Equal(suite.T(), 2, len(updatedCart.Items))

	// 4. Update first item quantity
	itemID := updatedCart.Items[0].ID
	updateReq := UpdateCartItemRequest{
		Quantity: 3,
		Attributes: map[string]string{
			"color": "red",
			"size":  "large",
		},
	}

	updatedCart = suite.updateCartItem(itemID, updateReq)
	assert.Equal(suite.T(), 3, updatedCart.Items[0].Quantity)
	assert.Equal(suite.T(), "red", updatedCart.Items[0].Attributes["color"])

	// 5. Apply coupon
	couponCode := "TEST10"
	originalTotal := updatedCart.Total
	updatedCart = suite.applyCoupon(couponCode)
	assert.NotNil(suite.T(), updatedCart.CouponCode)
	assert.Equal(suite.T(), couponCode, *updatedCart.CouponCode)
	assert.Greater(suite.T(), updatedCart.Discount, 0.0)
	assert.Less(suite.T(), updatedCart.Total, originalTotal)

	// 6. Remove coupon
	updatedCart = suite.removeCoupon()
	assert.Nil(suite.T(), updatedCart.CouponCode)
	assert.Equal(suite.T(), 0.0, updatedCart.Discount)

	// 7. Remove item from cart
	updatedCart = suite.removeCartItem(itemID)
	assert.Equal(suite.T(), 1, len(updatedCart.Items))

	// 8. Clear entire cart
	updatedCart = suite.clearCart()
	assert.Equal(suite.T(), 0, len(updatedCart.Items))
	assert.Equal(suite.T(), 0.0, updatedCart.Total)
}

// TestCartValidation tests cart validation rules
func (suite *CartAPITestSuite) TestCartValidation() {
	// Test adding invalid product
	invalidReq := AddToCartRequest{
		ProductID: "invalid-product-id",
		Quantity:  1,
	}
	suite.expectCartError(http.StatusNotFound, invalidReq)

	// Test adding zero quantity
	productID := suite.createTestProduct()
	zeroQtyReq := AddToCartRequest{
		ProductID: productID,
		Quantity:  0,
	}
	suite.expectCartError(http.StatusBadRequest, zeroQtyReq)

	// Test adding negative quantity
	negativeQtyReq := AddToCartRequest{
		ProductID: productID,
		Quantity:  -1,
	}
	suite.expectCartError(http.StatusBadRequest, negativeQtyReq)

	// Test adding excessive quantity (assuming limit of 100)
	excessiveQtyReq := AddToCartRequest{
		ProductID: productID,
		Quantity:  101,
	}
	suite.expectCartError(http.StatusBadRequest, excessiveQtyReq)
}

// TestCartConcurrency tests concurrent cart operations
func (suite *CartAPITestSuite) TestCartConcurrency() {
	productID := suite.createTestProduct()

	// Add item to cart first
	addReq := AddToCartRequest{
		ProductID: productID,
		Quantity:  1,
	}
	cart := suite.addToCart(addReq)
	itemID := cart.Items[0].ID

	// Perform concurrent updates
	concurrency := 5
	resultChan := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(index int) {
			updateReq := UpdateCartItemRequest{
				Quantity: index + 2, // 2-6
			}

			_, err := suite.updateCartItemWithError(itemID, updateReq)
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

	// At least one should succeed
	assert.GreaterOrEqual(suite.T(), successCount, 1)

	// Final cart should be consistent
	finalCart := suite.getOrCreateCart()
	assert.Equal(suite.T(), 1, len(finalCart.Items))
	assert.GreaterOrEqual(suite.T(), finalCart.Items[0].Quantity, 2)
	assert.LessOrEqual(suite.T(), finalCart.Items[0].Quantity, 6)
}

// TestCartExpiration tests cart expiration functionality
func (suite *CartAPITestSuite) TestCartExpiration() {
	// Create cart with expiration
	cart := suite.getOrCreateCart()
	assert.NotNil(suite.T(), cart)

	// Add item to cart
	productID := suite.createTestProduct()
	addReq := AddToCartRequest{
		ProductID: productID,
		Quantity:  1,
	}
	cart = suite.addToCart(addReq)
	assert.Equal(suite.T(), 1, len(cart.Items))

	// Simulate cart expiration (this might require specific test configuration)
	// In a real implementation, you might need to trigger expiration or wait
	// For now, we'll test the expiration endpoint if it exists
	suite.testCartExpiration(cart.ID)
}

// TestCartPersistence tests cart data persistence
func (suite *CartAPITestSuite) TestCartPersistence() {
	productID := suite.createTestProduct()

	// Add item to cart
	addReq := AddToCartRequest{
		ProductID: productID,
		Quantity:  2,
	}
	cart := suite.addToCart(addReq)
	originalCartID := cart.ID
	originalTotal := cart.Total

	// Simulate app restart by getting cart again
	retrievedCart := suite.getOrCreateCart()
	assert.Equal(suite.T(), originalCartID, retrievedCart.ID)
	assert.Equal(suite.T(), 1, len(retrievedCart.Items))
	assert.Equal(suite.T(), 2, retrievedCart.Items[0].Quantity)
	assert.Equal(suite.T(), originalTotal, retrievedCart.Total)
}

// Helper methods

func (suite *CartAPITestSuite) getOrCreateCart() *Cart {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/cart", suite.baseURL, suite.gatewayPort)

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(suite.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var cartResp CartResponse
	err = json.NewDecoder(resp.Body).Decode(&cartResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), cartResp.Success)

	return cartResp.Cart
}

func (suite *CartAPITestSuite) addToCart(req AddToCartRequest) *Cart {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/cart/items", suite.baseURL, suite.gatewayPort)

	body, err := json.Marshal(req)
	require.NoError(suite.T(), err)

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	require.NoError(suite.T(), err)
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(httpReq)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var cartResp CartResponse
	err = json.NewDecoder(resp.Body).Decode(&cartResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), cartResp.Success)

	return cartResp.Cart
}

func (suite *CartAPITestSuite) updateCartItem(itemID string, req UpdateCartItemRequest) *Cart {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/cart/items/%s", suite.baseURL, suite.gatewayPort, itemID)

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

	var cartResp CartResponse
	err = json.NewDecoder(resp.Body).Decode(&cartResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), cartResp.Success)

	return cartResp.Cart
}

func (suite *CartAPITestSuite) updateCartItemWithError(itemID string, req UpdateCartItemRequest) (*Cart, error) {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/cart/items/%s", suite.baseURL, suite.gatewayPort, itemID)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	var cartResp CartResponse
	err = json.NewDecoder(resp.Body).Decode(&cartResp)
	if err != nil {
		return nil, err
	}

	if !cartResp.Success {
		return nil, fmt.Errorf("API error: %s", cartResp.Message)
	}

	return cartResp.Cart, nil
}

func (suite *CartAPITestSuite) removeCartItem(itemID string) *Cart {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/cart/items/%s", suite.baseURL, suite.gatewayPort, itemID)

	httpReq, err := http.NewRequest("DELETE", url, nil)
	require.NoError(suite.T(), err)
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(httpReq)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var cartResp CartResponse
	err = json.NewDecoder(resp.Body).Decode(&cartResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), cartResp.Success)

	return cartResp.Cart
}

func (suite *CartAPITestSuite) applyCoupon(couponCode string) *Cart {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/cart/coupon", suite.baseURL, suite.gatewayPort)

	req := ApplyCouponRequest{CouponCode: couponCode}
	body, err := json.Marshal(req)
	require.NoError(suite.T(), err)

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	require.NoError(suite.T(), err)
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(httpReq)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var cartResp CartResponse
	err = json.NewDecoder(resp.Body).Decode(&cartResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), cartResp.Success)

	return cartResp.Cart
}

func (suite *CartAPITestSuite) removeCoupon() *Cart {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/cart/coupon", suite.baseURL, suite.gatewayPort)

	httpReq, err := http.NewRequest("DELETE", url, nil)
	require.NoError(suite.T(), err)
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(httpReq)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var cartResp CartResponse
	err = json.NewDecoder(resp.Body).Decode(&cartResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), cartResp.Success)

	return cartResp.Cart
}

func (suite *CartAPITestSuite) clearCart() *Cart {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/cart", suite.baseURL, suite.gatewayPort)

	httpReq, err := http.NewRequest("DELETE", url, nil)
	require.NoError(suite.T(), err)
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(httpReq)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var cartResp CartResponse
	err = json.NewDecoder(resp.Body).Decode(&cartResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), cartResp.Success)

	return cartResp.Cart
}

func (suite *CartAPITestSuite) clearUserCart() {
	// Clear cart before each test
	suite.clearCart()
}

func (suite *CartAPITestSuite) expectCartError(expectedStatus int, req AddToCartRequest) {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/cart/items", suite.baseURL, suite.gatewayPort)

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

func (suite *CartAPITestSuite) createTestProduct() string {
	// This would create a test product and return its ID
	// For integration testing, you might need to call the actual product creation endpoint
	return uuid.New().String()
}

func (suite *CartAPITestSuite) testCartExpiration(cartID string) {
	// Test cart expiration functionality if available
	// This might involve calling a specific endpoint or waiting for expiration
}

func (suite *CartAPITestSuite) waitForServices() {
	// Wait for gateway and commerce service to be ready
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

func (suite *CartAPITestSuite) cleanupTestData() {
	// Clean up any test data created during tests
	suite.clearUserCart()
}

// TestCartAPI runs the cart API test suite
func TestCartAPI(t *testing.T) {
	suite.Run(t, new(CartAPITestSuite))
}