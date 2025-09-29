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

// ProductAPITestSuite provides comprehensive integration testing for product API endpoints
type ProductAPITestSuite struct {
	suite.Suite
	baseURL      string
	httpClient   *http.Client
	gatewayPort  int
	testShopID   string
	testUserID   string
	createdProducts []string
}

// Product represents product structure
type Product struct {
	ID             string                 `json:"id"`
	ShopID         string                 `json:"shopId"`
	Name           string                 `json:"name"`
	Description    *string                `json:"description,omitempty"`
	ShortDesc      *string                `json:"shortDescription,omitempty"`
	Category       string                 `json:"category"`
	SubCategory    *string                `json:"subCategory,omitempty"`
	Brand          *string                `json:"brand,omitempty"`
	SKU            string                 `json:"sku"`
	Price          float64                `json:"price"`
	ComparePrice   *float64               `json:"comparePrice,omitempty"`
	Currency       string                 `json:"currency"`
	Status         string                 `json:"status"`
	Visibility     string                 `json:"visibility"`
	Inventory      ProductInventory       `json:"inventory"`
	Images         []ProductImage         `json:"images"`
	Variants       []ProductVariant       `json:"variants"`
	Attributes     map[string]interface{} `json:"attributes"`
	Tags           []string               `json:"tags"`
	Weight         *float64               `json:"weight,omitempty"`
	Dimensions     *ProductDimensions     `json:"dimensions,omitempty"`
	SEO            *ProductSEO            `json:"seo,omitempty"`
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt"`
	PublishedAt    *time.Time             `json:"publishedAt,omitempty"`
}

// ProductInventory represents inventory information
type ProductInventory struct {
	TrackQuantity   bool   `json:"trackQuantity"`
	Quantity        int    `json:"quantity"`
	ReservedQty     int    `json:"reservedQuantity"`
	AvailableQty    int    `json:"availableQuantity"`
	StockStatus     string `json:"stockStatus"`
	LowStockLevel   *int   `json:"lowStockLevel,omitempty"`
	BackorderAllowed bool  `json:"backorderAllowed"`
}

// ProductImage represents product image
type ProductImage struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	AltText  string `json:"altText"`
	Position int    `json:"position"`
	IsMain   bool   `json:"isMain"`
}

// ProductVariant represents product variant
type ProductVariant struct {
	ID          string                 `json:"id"`
	SKU         string                 `json:"sku"`
	Name        string                 `json:"name"`
	Price       *float64               `json:"price,omitempty"`
	ComparePrice *float64              `json:"comparePrice,omitempty"`
	Options     map[string]string      `json:"options"`
	Inventory   ProductInventory       `json:"inventory"`
	Images      []ProductImage         `json:"images"`
	Attributes  map[string]interface{} `json:"attributes"`
	Status      string                 `json:"status"`
}

// ProductDimensions represents product dimensions
type ProductDimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"`
}

// ProductSEO represents SEO data
type ProductSEO struct {
	Title         *string `json:"title,omitempty"`
	Description   *string `json:"description,omitempty"`
	Keywords      *string `json:"keywords,omitempty"`
	MetaTitle     *string `json:"metaTitle,omitempty"`
	MetaDesc      *string `json:"metaDescription,omitempty"`
}

// CreateProductRequest represents request to create product
type CreateProductRequest struct {
	Name           string                 `json:"name"`
	Description    *string                `json:"description,omitempty"`
	ShortDesc      *string                `json:"shortDescription,omitempty"`
	Category       string                 `json:"category"`
	SubCategory    *string                `json:"subCategory,omitempty"`
	Brand          *string                `json:"brand,omitempty"`
	SKU            string                 `json:"sku"`
	Price          float64                `json:"price"`
	ComparePrice   *float64               `json:"comparePrice,omitempty"`
	Currency       string                 `json:"currency"`
	Status         *string                `json:"status,omitempty"`
	Visibility     *string                `json:"visibility,omitempty"`
	Inventory      ProductInventory       `json:"inventory"`
	Images         []ProductImage         `json:"images"`
	Variants       []ProductVariant       `json:"variants"`
	Attributes     map[string]interface{} `json:"attributes"`
	Tags           []string               `json:"tags"`
	Weight         *float64               `json:"weight,omitempty"`
	Dimensions     *ProductDimensions     `json:"dimensions,omitempty"`
	SEO            *ProductSEO            `json:"seo,omitempty"`
}

// UpdateProductRequest represents request to update product
type UpdateProductRequest struct {
	Name           *string                `json:"name,omitempty"`
	Description    *string                `json:"description,omitempty"`
	ShortDesc      *string                `json:"shortDescription,omitempty"`
	Category       *string                `json:"category,omitempty"`
	SubCategory    *string                `json:"subCategory,omitempty"`
	Brand          *string                `json:"brand,omitempty"`
	Price          *float64               `json:"price,omitempty"`
	ComparePrice   *float64               `json:"comparePrice,omitempty"`
	Status         *string                `json:"status,omitempty"`
	Visibility     *string                `json:"visibility,omitempty"`
	Inventory      *ProductInventory      `json:"inventory,omitempty"`
	Images         []ProductImage         `json:"images,omitempty"`
	Variants       []ProductVariant       `json:"variants,omitempty"`
	Attributes     map[string]interface{} `json:"attributes,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	Weight         *float64               `json:"weight,omitempty"`
	Dimensions     *ProductDimensions     `json:"dimensions,omitempty"`
	SEO            *ProductSEO            `json:"seo,omitempty"`
}

// ProductSearchRequest represents search parameters
type ProductSearchRequest struct {
	Query        string   `json:"query,omitempty"`
	Category     string   `json:"category,omitempty"`
	SubCategory  string   `json:"subCategory,omitempty"`
	Brand        string   `json:"brand,omitempty"`
	MinPrice     *float64 `json:"minPrice,omitempty"`
	MaxPrice     *float64 `json:"maxPrice,omitempty"`
	InStock      *bool    `json:"inStock,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Status       string   `json:"status,omitempty"`
	Visibility   string   `json:"visibility,omitempty"`
	SortBy       string   `json:"sortBy,omitempty"`
	SortOrder    string   `json:"sortOrder,omitempty"`
	Page         int      `json:"page,omitempty"`
	Limit        int      `json:"limit,omitempty"`
}

// ProductResponse represents API response for product operations
type ProductResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Product   *Product  `json:"product,omitempty"`
	Products  []Product `json:"products,omitempty"`
	Total     int       `json:"total,omitempty"`
	Page      int       `json:"page,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	Error     *string   `json:"error,omitempty"`
	Timestamp string    `json:"timestamp"`
}

// SetupSuite initializes the test suite
func (suite *ProductAPITestSuite) SetupSuite() {
	suite.baseURL = "http://localhost"
	suite.gatewayPort = 8080
	suite.httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
	suite.testUserID = uuid.New().String()
	suite.testShopID = uuid.New().String()
	suite.createdProducts = []string{}

	// Wait for services and create test shop
	suite.waitForServices()
	suite.createTestShop()
}

// TearDownSuite cleans up after test suite
func (suite *ProductAPITestSuite) TearDownSuite() {
	suite.cleanupTestData()
}

// TestProductCRUD tests complete product CRUD operations
func (suite *ProductAPITestSuite) TestProductCRUD() {
	// 1. Create product
	createReq := CreateProductRequest{
		Name:        "Test Smartphone Pro",
		Description: stringPtr("Advanced smartphone with premium features"),
		ShortDesc:   stringPtr("Premium smartphone"),
		Category:    "electronics",
		SubCategory: stringPtr("smartphones"),
		Brand:       stringPtr("TestBrand"),
		SKU:         fmt.Sprintf("TSP-%d", time.Now().Unix()),
		Price:       999.99,
		ComparePrice: floatPtr(1199.99),
		Currency:    "USD",
		Status:      stringPtr("active"),
		Visibility:  stringPtr("public"),
		Inventory: ProductInventory{
			TrackQuantity:    true,
			Quantity:         100,
			StockStatus:      "in_stock",
			LowStockLevel:    intPtr(10),
			BackorderAllowed: false,
		},
		Images: []ProductImage{
			{
				ID:       uuid.New().String(),
				URL:      "https://cdn.example.com/smartphone-main.jpg",
				AltText:  "Test Smartphone Pro - Main View",
				Position: 1,
				IsMain:   true,
			},
			{
				ID:       uuid.New().String(),
				URL:      "https://cdn.example.com/smartphone-back.jpg",
				AltText:  "Test Smartphone Pro - Back View",
				Position: 2,
				IsMain:   false,
			},
		},
		Variants: []ProductVariant{
			{
				ID:   uuid.New().String(),
				SKU:  "TSP-128-BLACK",
				Name: "128GB Black",
				Options: map[string]string{
					"storage": "128GB",
					"color":   "black",
				},
				Inventory: ProductInventory{
					TrackQuantity: true,
					Quantity:      50,
					StockStatus:   "in_stock",
				},
				Status: "active",
			},
			{
				ID:   uuid.New().String(),
				SKU:  "TSP-256-WHITE",
				Name: "256GB White",
				Price: floatPtr(1099.99),
				Options: map[string]string{
					"storage": "256GB",
					"color":   "white",
				},
				Inventory: ProductInventory{
					TrackQuantity: true,
					Quantity:      30,
					StockStatus:   "in_stock",
				},
				Status: "active",
			},
		},
		Attributes: map[string]interface{}{
			"processor":     "A15 Bionic",
			"screen_size":   "6.1 inches",
			"battery_life":  "24 hours",
			"water_resistant": true,
			"warranty_years": 2,
		},
		Tags: []string{"smartphone", "premium", "5g", "wireless_charging"},
		Weight: floatPtr(174.0),
		Dimensions: &ProductDimensions{
			Length: 146.7,
			Width:  71.5,
			Height: 7.65,
			Unit:   "mm",
		},
		SEO: &ProductSEO{
			Title:       stringPtr("Test Smartphone Pro - Premium 5G Phone"),
			Description: stringPtr("Experience the latest technology with Test Smartphone Pro"),
			Keywords:    stringPtr("smartphone, 5g, premium, wireless charging"),
			MetaTitle:   stringPtr("Buy Test Smartphone Pro Online"),
			MetaDesc:    stringPtr("Shop the latest Test Smartphone Pro with free shipping"),
		},
	}

	product := suite.createProduct(createReq)
	productID := product.ID
	suite.createdProducts = append(suite.createdProducts, productID)

	// Verify creation
	assert.Equal(suite.T(), createReq.Name, product.Name)
	assert.Equal(suite.T(), createReq.SKU, product.SKU)
	assert.Equal(suite.T(), createReq.Price, product.Price)
	assert.Equal(suite.T(), suite.testShopID, product.ShopID)
	assert.Equal(suite.T(), "active", product.Status)
	assert.Equal(suite.T(), 2, len(product.Images))
	assert.Equal(suite.T(), 2, len(product.Variants))
	assert.Equal(suite.T(), 4, len(product.Tags))

	// 2. Get product by ID
	retrievedProduct := suite.getProduct(productID)
	assert.Equal(suite.T(), product.ID, retrievedProduct.ID)
	assert.Equal(suite.T(), product.Name, retrievedProduct.Name)
	assert.Equal(suite.T(), product.SKU, retrievedProduct.SKU)

	// 3. Update product
	updateReq := UpdateProductRequest{
		Name:        stringPtr("Test Smartphone Pro Max"),
		Price:       floatPtr(1099.99),
		ComparePrice: floatPtr(1299.99),
		Attributes: map[string]interface{}{
			"processor":     "A16 Bionic",
			"screen_size":   "6.7 inches",
			"battery_life":  "28 hours",
			"water_resistant": true,
			"warranty_years": 3,
		},
		Tags: []string{"smartphone", "premium", "5g", "wireless_charging", "pro_max"},
	}

	updatedProduct := suite.updateProduct(productID, updateReq)
	assert.Equal(suite.T(), "Test Smartphone Pro Max", updatedProduct.Name)
	assert.Equal(suite.T(), 1099.99, updatedProduct.Price)
	assert.Equal(suite.T(), 5, len(updatedProduct.Tags))
	assert.Equal(suite.T(), "A16 Bionic", updatedProduct.Attributes["processor"])

	// 4. Delete product
	suite.deleteProduct(productID)
	suite.expectProductNotFound(productID)
}

// TestProductSearch tests product search and filtering
func (suite *ProductAPITestSuite) TestProductSearch() {
	// Create test products with different attributes
	products := []CreateProductRequest{
		{
			Name:     "iPhone 15 Pro",
			Category: "electronics",
			SubCategory: stringPtr("smartphones"),
			Brand:    stringPtr("Apple"),
			SKU:      fmt.Sprintf("IP15P-%d", time.Now().Unix()),
			Price:    999.99,
			Currency: "USD",
			Tags:     []string{"smartphone", "premium", "apple"},
			Inventory: ProductInventory{
				TrackQuantity: true,
				Quantity:      50,
				StockStatus:   "in_stock",
			},
		},
		{
			Name:     "Samsung Galaxy S24",
			Category: "electronics",
			SubCategory: stringPtr("smartphones"),
			Brand:    stringPtr("Samsung"),
			SKU:      fmt.Sprintf("SGS24-%d", time.Now().Unix()),
			Price:    799.99,
			Currency: "USD",
			Tags:     []string{"smartphone", "android", "samsung"},
			Inventory: ProductInventory{
				TrackQuantity: true,
				Quantity:      0,
				StockStatus:   "out_of_stock",
			},
		},
		{
			Name:     "MacBook Air M2",
			Category: "electronics",
			SubCategory: stringPtr("laptops"),
			Brand:    stringPtr("Apple"),
			SKU:      fmt.Sprintf("MBA2-%d", time.Now().Unix()),
			Price:    1199.99,
			Currency: "USD",
			Tags:     []string{"laptop", "premium", "apple"},
			Inventory: ProductInventory{
				TrackQuantity: true,
				Quantity:      25,
				StockStatus:   "in_stock",
			},
		},
	}

	// Create products
	createdProductIDs := []string{}
	for _, req := range products {
		product := suite.createProduct(req)
		createdProductIDs = append(createdProductIDs, product.ID)
		suite.createdProducts = append(suite.createdProducts, product.ID)
	}

	// Test various search scenarios
	testCases := []struct {
		name           string
		searchReq      ProductSearchRequest
		expectedCount  int
		expectedBrands []string
	}{
		{
			name: "Search by query",
			searchReq: ProductSearchRequest{
				Query: "iPhone",
			},
			expectedCount:  1,
			expectedBrands: []string{"Apple"},
		},
		{
			name: "Search by category",
			searchReq: ProductSearchRequest{
				Category: "electronics",
			},
			expectedCount: 3,
		},
		{
			name: "Search by subcategory",
			searchReq: ProductSearchRequest{
				SubCategory: "smartphones",
			},
			expectedCount: 2,
		},
		{
			name: "Search by brand",
			searchReq: ProductSearchRequest{
				Brand: "Apple",
			},
			expectedCount:  2,
			expectedBrands: []string{"Apple", "Apple"},
		},
		{
			name: "Search by price range",
			searchReq: ProductSearchRequest{
				MinPrice: floatPtr(800.0),
				MaxPrice: floatPtr(1000.0),
			},
			expectedCount: 2,
		},
		{
			name: "Search in stock only",
			searchReq: ProductSearchRequest{
				InStock: boolPtr(true),
			},
			expectedCount: 2,
		},
		{
			name: "Search by tags",
			searchReq: ProductSearchRequest{
				Tags: []string{"premium"},
			},
			expectedCount: 2,
		},
		{
			name: "Complex search",
			searchReq: ProductSearchRequest{
				Category: "electronics",
				Brand:    "Apple",
				InStock:  boolPtr(true),
			},
			expectedCount: 2,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			results := suite.searchProducts(tc.searchReq)
			assert.Equal(t, tc.expectedCount, len(results.Products))

			if tc.expectedBrands != nil {
				brands := []string{}
				for _, product := range results.Products {
					if product.Brand != nil {
						brands = append(brands, *product.Brand)
					}
				}
				assert.ElementsMatch(t, tc.expectedBrands, brands)
			}
		})
	}
}

// TestProductInventoryManagement tests inventory operations
func (suite *ProductAPITestSuite) TestProductInventoryManagement() {
	// Create product with inventory tracking
	createReq := CreateProductRequest{
		Name:     "Inventory Test Product",
		Category: "test",
		SKU:      fmt.Sprintf("ITP-%d", time.Now().Unix()),
		Price:    99.99,
		Currency: "USD",
		Inventory: ProductInventory{
			TrackQuantity:    true,
			Quantity:         100,
			StockStatus:      "in_stock",
			LowStockLevel:    intPtr(20),
			BackorderAllowed: false,
		},
	}

	product := suite.createProduct(createReq)
	productID := product.ID
	suite.createdProducts = append(suite.createdProducts, productID)

	// Test inventory updates
	suite.updateInventory(productID, 150)
	updatedProduct := suite.getProduct(productID)
	assert.Equal(suite.T(), 150, updatedProduct.Inventory.Quantity)

	// Test low stock
	suite.updateInventory(productID, 15)
	updatedProduct = suite.getProduct(productID)
	assert.Equal(suite.T(), 15, updatedProduct.Inventory.Quantity)
	// Note: Stock status might change based on business logic

	// Test out of stock
	suite.updateInventory(productID, 0)
	updatedProduct = suite.getProduct(productID)
	assert.Equal(suite.T(), 0, updatedProduct.Inventory.Quantity)
	assert.Equal(suite.T(), "out_of_stock", updatedProduct.Inventory.StockStatus)
}

// TestProductValidation tests product validation rules
func (suite *ProductAPITestSuite) TestProductValidation() {
	testCases := []struct {
		name        string
		createReq   CreateProductRequest
		expectError bool
		statusCode  int
	}{
		{
			name: "Missing required name",
			createReq: CreateProductRequest{
				Category: "test",
				SKU:      "TEST-001",
				Price:    99.99,
				Currency: "USD",
			},
			expectError: true,
			statusCode:  http.StatusBadRequest,
		},
		{
			name: "Missing required category",
			createReq: CreateProductRequest{
				Name:     "Test Product",
				SKU:      "TEST-002",
				Price:    99.99,
				Currency: "USD",
			},
			expectError: true,
			statusCode:  http.StatusBadRequest,
		},
		{
			name: "Duplicate SKU",
			createReq: CreateProductRequest{
				Name:     "Duplicate SKU Product",
				Category: "test",
				SKU:      "DUPLICATE-SKU-001",
				Price:    99.99,
				Currency: "USD",
			},
			expectError: true,
			statusCode:  http.StatusConflict,
		},
		{
			name: "Invalid price",
			createReq: CreateProductRequest{
				Name:     "Invalid Price Product",
				Category: "test",
				SKU:      "INVALID-PRICE-001",
				Price:    -99.99,
				Currency: "USD",
			},
			expectError: true,
			statusCode:  http.StatusBadRequest,
		},
		{
			name: "Invalid currency",
			createReq: CreateProductRequest{
				Name:     "Invalid Currency Product",
				Category: "test",
				SKU:      "INVALID-CURR-001",
				Price:    99.99,
				Currency: "INVALID",
			},
			expectError: true,
			statusCode:  http.StatusBadRequest,
		},
	}

	// First create a product with a specific SKU for duplicate test
	firstProduct := CreateProductRequest{
		Name:     "First Product",
		Category: "test",
		SKU:      "DUPLICATE-SKU-001",
		Price:    99.99,
		Currency: "USD",
	}
	product := suite.createProduct(firstProduct)
	suite.createdProducts = append(suite.createdProducts, product.ID)

	// Run validation tests
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			if tc.expectError {
				suite.expectProductCreationError(tc.statusCode, tc.createReq)
			} else {
				product := suite.createProduct(tc.createReq)
				suite.createdProducts = append(suite.createdProducts, product.ID)
				assert.NotEmpty(t, product.ID)
			}
		})
	}
}

// TestProductConcurrency tests concurrent product operations
func (suite *ProductAPITestSuite) TestProductConcurrency() {
	// Create a product for concurrent updates
	createReq := CreateProductRequest{
		Name:     "Concurrency Test Product",
		Category: "test",
		SKU:      fmt.Sprintf("CTP-%d", time.Now().Unix()),
		Price:    99.99,
		Currency: "USD",
		Inventory: ProductInventory{
			TrackQuantity: true,
			Quantity:      100,
			StockStatus:   "in_stock",
		},
	}

	product := suite.createProduct(createReq)
	productID := product.ID
	suite.createdProducts = append(suite.createdProducts, productID)

	// Perform concurrent inventory updates
	concurrency := 5
	resultChan := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(index int) {
			newQuantity := 50 + index*10 // 50, 60, 70, 80, 90
			err := suite.updateInventoryWithError(productID, newQuantity)
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

	// Final product should be consistent
	finalProduct := suite.getProduct(productID)
	assert.GreaterOrEqual(suite.T(), finalProduct.Inventory.Quantity, 50)
	assert.LessOrEqual(suite.T(), finalProduct.Inventory.Quantity, 90)
}

// Helper methods

func (suite *ProductAPITestSuite) createProduct(req CreateProductRequest) *Product {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/shops/%s/products", suite.baseURL, suite.gatewayPort, suite.testShopID)

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

	var productResp ProductResponse
	err = json.NewDecoder(resp.Body).Decode(&productResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), productResp.Success)
	require.NotNil(suite.T(), productResp.Product)

	return productResp.Product
}

func (suite *ProductAPITestSuite) getProduct(productID string) *Product {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/products/%s", suite.baseURL, suite.gatewayPort, productID)

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(suite.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var productResp ProductResponse
	err = json.NewDecoder(resp.Body).Decode(&productResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), productResp.Success)
	require.NotNil(suite.T(), productResp.Product)

	return productResp.Product
}

func (suite *ProductAPITestSuite) updateProduct(productID string, req UpdateProductRequest) *Product {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/products/%s", suite.baseURL, suite.gatewayPort, productID)

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

	var productResp ProductResponse
	err = json.NewDecoder(resp.Body).Decode(&productResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), productResp.Success)
	require.NotNil(suite.T(), productResp.Product)

	return productResp.Product
}

func (suite *ProductAPITestSuite) deleteProduct(productID string) {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/products/%s", suite.baseURL, suite.gatewayPort, productID)

	req, err := http.NewRequest("DELETE", url, nil)
	require.NoError(suite.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *ProductAPITestSuite) searchProducts(req ProductSearchRequest) *ProductResponse {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/products/search", suite.baseURL, suite.gatewayPort)

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

	var productResp ProductResponse
	err = json.NewDecoder(resp.Body).Decode(&productResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), productResp.Success)

	return &productResp
}

func (suite *ProductAPITestSuite) updateInventory(productID string, quantity int) {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/products/%s/inventory", suite.baseURL, suite.gatewayPort, productID)

	reqBody := map[string]interface{}{
		"quantity": quantity,
	}
	body, err := json.Marshal(reqBody)
	require.NoError(suite.T(), err)

	httpReq, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	require.NoError(suite.T(), err)
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(httpReq)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *ProductAPITestSuite) updateInventoryWithError(productID string, quantity int) error {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/products/%s/inventory", suite.baseURL, suite.gatewayPort, productID)

	reqBody := map[string]interface{}{
		"quantity": quantity,
	}
	body, err := json.Marshal(reqBody)
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

func (suite *ProductAPITestSuite) expectProductNotFound(productID string) {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/products/%s", suite.baseURL, suite.gatewayPort, productID)

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(suite.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *ProductAPITestSuite) expectProductCreationError(expectedStatus int, req CreateProductRequest) {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/shops/%s/products", suite.baseURL, suite.gatewayPort, suite.testShopID)

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

func (suite *ProductAPITestSuite) createTestShop() {
	// This would create a test shop - simplified for integration testing
	suite.testShopID = uuid.New().String()
}

func (suite *ProductAPITestSuite) waitForServices() {
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

func (suite *ProductAPITestSuite) cleanupTestData() {
	// Clean up created products
	for _, productID := range suite.createdProducts {
		suite.deleteProduct(productID)
	}
}

// Utility functions
func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

// TestProductAPI runs the product API test suite
func TestProductAPI(t *testing.T) {
	suite.Run(t, new(ProductAPITestSuite))
}