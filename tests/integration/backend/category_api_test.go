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

// CategoryAPITestSuite provides comprehensive integration testing for category API endpoints
type CategoryAPITestSuite struct {
	suite.Suite
	baseURL         string
	httpClient      *http.Client
	gatewayPort     int
	testUserID      string
	createdCategories []string
}

// Category represents category structure
type Category struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Slug           string                 `json:"slug"`
	Description    *string                `json:"description,omitempty"`
	ParentID       *string                `json:"parentId,omitempty"`
	Level          int                    `json:"level"`
	Path           string                 `json:"path"`
	Children       []Category             `json:"children,omitempty"`
	Image          *CategoryImage         `json:"image,omitempty"`
	Icon           *string                `json:"icon,omitempty"`
	Color          *string                `json:"color,omitempty"`
	Status         string                 `json:"status"`
	Visibility     string                 `json:"visibility"`
	SortOrder      int                    `json:"sortOrder"`
	ProductCount   int                    `json:"productCount"`
	Attributes     map[string]interface{} `json:"attributes,omitempty"`
	SEO            *CategorySEO           `json:"seo,omitempty"`
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt"`
}

// CategoryImage represents category image
type CategoryImage struct {
	URL     string `json:"url"`
	AltText string `json:"altText"`
	Width   *int   `json:"width,omitempty"`
	Height  *int   `json:"height,omitempty"`
}

// CategorySEO represents category SEO data
type CategorySEO struct {
	Title         *string `json:"title,omitempty"`
	Description   *string `json:"description,omitempty"`
	Keywords      *string `json:"keywords,omitempty"`
	MetaTitle     *string `json:"metaTitle,omitempty"`
	MetaDesc      *string `json:"metaDescription,omitempty"`
	CanonicalURL  *string `json:"canonicalUrl,omitempty"`
}

// CreateCategoryRequest represents request to create category
type CreateCategoryRequest struct {
	Name         string                 `json:"name"`
	Slug         *string                `json:"slug,omitempty"`
	Description  *string                `json:"description,omitempty"`
	ParentID     *string                `json:"parentId,omitempty"`
	Image        *CategoryImage         `json:"image,omitempty"`
	Icon         *string                `json:"icon,omitempty"`
	Color        *string                `json:"color,omitempty"`
	Status       *string                `json:"status,omitempty"`
	Visibility   *string                `json:"visibility,omitempty"`
	SortOrder    *int                   `json:"sortOrder,omitempty"`
	Attributes   map[string]interface{} `json:"attributes,omitempty"`
	SEO          *CategorySEO           `json:"seo,omitempty"`
}

// UpdateCategoryRequest represents request to update category
type UpdateCategoryRequest struct {
	Name         *string                `json:"name,omitempty"`
	Slug         *string                `json:"slug,omitempty"`
	Description  *string                `json:"description,omitempty"`
	ParentID     *string                `json:"parentId,omitempty"`
	Image        *CategoryImage         `json:"image,omitempty"`
	Icon         *string                `json:"icon,omitempty"`
	Color        *string                `json:"color,omitempty"`
	Status       *string                `json:"status,omitempty"`
	Visibility   *string                `json:"visibility,omitempty"`
	SortOrder    *int                   `json:"sortOrder,omitempty"`
	Attributes   map[string]interface{} `json:"attributes,omitempty"`
	SEO          *CategorySEO           `json:"seo,omitempty"`
}

// CategoryResponse represents API response for category operations
type CategoryResponse struct {
	Success    bool       `json:"success"`
	Message    string     `json:"message"`
	Category   *Category  `json:"category,omitempty"`
	Categories []Category `json:"categories,omitempty"`
	Total      int        `json:"total,omitempty"`
	Error      *string    `json:"error,omitempty"`
	Timestamp  string     `json:"timestamp"`
}

// CategoryAnalytics represents category analytics data
type CategoryAnalytics struct {
	CategoryID      string            `json:"categoryId"`
	ProductCount    int               `json:"productCount"`
	ViewCount       int64             `json:"viewCount"`
	ConversionRate  float64           `json:"conversionRate"`
	Revenue         float64           `json:"revenue"`
	TopProducts     []ProductSummary  `json:"topProducts"`
	TrendData       []TrendPoint      `json:"trendData"`
	Demographics    Demographics      `json:"demographics"`
	Period          string            `json:"period"`
	LastUpdated     time.Time         `json:"lastUpdated"`
}

// ProductSummary represents basic product information
type ProductSummary struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Revenue  float64 `json:"revenue"`
	Units    int     `json:"units"`
}

// TrendPoint represents a data point in trend analysis
type TrendPoint struct {
	Date   time.Time `json:"date"`
	Value  float64   `json:"value"`
	Metric string    `json:"metric"`
}

// Demographics represents demographic data
type Demographics struct {
	AgeGroups     map[string]int `json:"ageGroups"`
	GenderSplit   map[string]int `json:"genderSplit"`
	TopLocations  []LocationData `json:"topLocations"`
}

// LocationData represents location-based data
type LocationData struct {
	Country string `json:"country"`
	Count   int    `json:"count"`
}

// SetupSuite initializes the test suite
func (suite *CategoryAPITestSuite) SetupSuite() {
	suite.baseURL = "http://localhost"
	suite.gatewayPort = 8080
	suite.httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
	suite.testUserID = uuid.New().String()
	suite.createdCategories = []string{}

	// Wait for services
	suite.waitForServices()
}

// TearDownSuite cleans up after test suite
func (suite *CategoryAPITestSuite) TearDownSuite() {
	suite.cleanupTestData()
}

// TestCategoryHierarchy tests category hierarchy management
func (suite *CategoryAPITestSuite) TestCategoryHierarchy() {
	// 1. Create root category
	electronicsReq := CreateCategoryRequest{
		Name:        "Electronics",
		Description: stringPtr("Electronic devices and accessories"),
		Image: &CategoryImage{
			URL:     "https://cdn.example.com/categories/electronics.jpg",
			AltText: "Electronics Category",
			Width:   intPtr(800),
			Height:  intPtr(600),
		},
		Icon:       stringPtr("📱"),
		Color:      stringPtr("#2563eb"),
		Status:     stringPtr("active"),
		Visibility: stringPtr("public"),
		SortOrder:  intPtr(1),
		Attributes: map[string]interface{}{
			"featured":     true,
			"commission":   15.0,
			"tax_category": "electronics",
		},
		SEO: &CategorySEO{
			Title:       stringPtr("Electronics - Latest Gadgets & Devices"),
			Description: stringPtr("Shop the latest electronics including smartphones, laptops, and accessories"),
			Keywords:    stringPtr("electronics, gadgets, smartphones, laptops"),
			MetaTitle:   stringPtr("Buy Electronics Online"),
			MetaDesc:    stringPtr("Discover the latest electronics with free shipping"),
		},
	}

	electronics := suite.createCategory(electronicsReq)
	electronicsID := electronics.ID
	suite.createdCategories = append(suite.createdCategories, electronicsID)

	// Verify root category
	assert.Equal(suite.T(), "Electronics", electronics.Name)
	assert.Equal(suite.T(), "electronics", electronics.Slug)
	assert.Equal(suite.T(), 0, electronics.Level)
	assert.Nil(suite.T(), electronics.ParentID)
	assert.Equal(suite.T(), "/electronics", electronics.Path)

	// 2. Create subcategory
	smartphonesReq := CreateCategoryRequest{
		Name:        "Smartphones",
		Description: stringPtr("Mobile phones and accessories"),
		ParentID:    &electronicsID,
		Icon:        stringPtr("📱"),
		Color:       stringPtr("#059669"),
		Status:      stringPtr("active"),
		Visibility:  stringPtr("public"),
		SortOrder:   intPtr(1),
	}

	smartphones := suite.createCategory(smartphonesReq)
	smartphonesID := smartphones.ID
	suite.createdCategories = append(suite.createdCategories, smartphonesID)

	// Verify subcategory
	assert.Equal(suite.T(), "Smartphones", smartphones.Name)
	assert.Equal(suite.T(), "smartphones", smartphones.Slug)
	assert.Equal(suite.T(), 1, smartphones.Level)
	assert.Equal(suite.T(), electronicsID, *smartphones.ParentID)
	assert.Equal(suite.T(), "/electronics/smartphones", smartphones.Path)

	// 3. Create sub-subcategory
	iphoneReq := CreateCategoryRequest{
		Name:        "iPhone",
		Description: stringPtr("Apple iPhone models"),
		ParentID:    &smartphonesID,
		Icon:        stringPtr("🍎"),
		Color:       stringPtr("#000000"),
		Status:      stringPtr("active"),
		Visibility:  stringPtr("public"),
		SortOrder:   intPtr(1),
	}

	iphone := suite.createCategory(iphoneReq)
	iphoneID := iphone.ID
	suite.createdCategories = append(suite.createdCategories, iphoneID)

	// Verify sub-subcategory
	assert.Equal(suite.T(), "iPhone", iphone.Name)
	assert.Equal(suite.T(), 2, iphone.Level)
	assert.Equal(suite.T(), smartphonesID, *iphone.ParentID)
	assert.Equal(suite.T(), "/electronics/smartphones/iphone", iphone.Path)

	// 4. Test hierarchy retrieval
	hierarchy := suite.getCategoryHierarchy()
	assert.GreaterOrEqual(suite.T(), len(hierarchy.Categories), 1)

	// Find electronics category with children
	var electronicsWithChildren *Category
	for _, cat := range hierarchy.Categories {
		if cat.ID == electronicsID {
			electronicsWithChildren = &cat
			break
		}
	}
	assert.NotNil(suite.T(), electronicsWithChildren)
	assert.GreaterOrEqual(suite.T(), len(electronicsWithChildren.Children), 1)

	// 5. Test moving category (change parent)
	updateReq := UpdateCategoryRequest{
		ParentID: nil, // Move to root
	}
	movedCategory := suite.updateCategory(smartphonesID, updateReq)
	assert.Nil(suite.T(), movedCategory.ParentID)
	assert.Equal(suite.T(), 0, movedCategory.Level)
	assert.Equal(suite.T(), "/smartphones", movedCategory.Path)
}

// TestCategoryCRUD tests basic category CRUD operations
func (suite *CategoryAPITestSuite) TestCategoryCRUD() {
	// 1. Create category
	createReq := CreateCategoryRequest{
		Name:        "Fashion & Clothing",
		Description: stringPtr("Clothing, shoes, and fashion accessories"),
		Image: &CategoryImage{
			URL:     "https://cdn.example.com/categories/fashion.jpg",
			AltText: "Fashion Category",
		},
		Icon:       stringPtr("👗"),
		Color:      stringPtr("#ec4899"),
		Status:     stringPtr("active"),
		Visibility: stringPtr("public"),
		SortOrder:  intPtr(2),
		Attributes: map[string]interface{}{
			"seasonal":     true,
			"target_age":   "18-65",
			"style_guide":  "modern",
		},
	}

	category := suite.createCategory(createReq)
	categoryID := category.ID
	suite.createdCategories = append(suite.createdCategories, categoryID)

	// Verify creation
	assert.Equal(suite.T(), createReq.Name, category.Name)
	assert.Equal(suite.T(), "fashion-clothing", category.Slug)
	assert.Equal(suite.T(), "active", category.Status)
	assert.Equal(suite.T(), "public", category.Visibility)
	assert.Equal(suite.T(), 2, category.SortOrder)

	// 2. Get category by ID
	retrievedCategory := suite.getCategory(categoryID)
	assert.Equal(suite.T(), category.ID, retrievedCategory.ID)
	assert.Equal(suite.T(), category.Name, retrievedCategory.Name)
	assert.Equal(suite.T(), category.Slug, retrievedCategory.Slug)

	// 3. Update category
	updateReq := UpdateCategoryRequest{
		Name:        stringPtr("Fashion & Apparel"),
		Description: stringPtr("Updated description for fashion category"),
		Color:       stringPtr("#db2777"),
		SortOrder:   intPtr(3),
		Attributes: map[string]interface{}{
			"seasonal":      true,
			"target_age":    "16-70",
			"style_guide":   "contemporary",
			"new_attribute": "value",
		},
	}

	updatedCategory := suite.updateCategory(categoryID, updateReq)
	assert.Equal(suite.T(), "Fashion & Apparel", updatedCategory.Name)
	assert.Equal(suite.T(), "fashion-apparel", updatedCategory.Slug)
	assert.Equal(suite.T(), "#db2777", *updatedCategory.Color)
	assert.Equal(suite.T(), 3, updatedCategory.SortOrder)

	// 4. List categories
	categories := suite.listCategories()
	assert.GreaterOrEqual(suite.T(), len(categories.Categories), 1)

	// 5. Delete category
	suite.deleteCategory(categoryID)
	suite.expectCategoryNotFound(categoryID)
}

// TestCategoryValidation tests category validation rules
func (suite *CategoryAPITestSuite) TestCategoryValidation() {
	testCases := []struct {
		name        string
		createReq   CreateCategoryRequest
		expectError bool
		statusCode  int
	}{
		{
			name: "Missing required name",
			createReq: CreateCategoryRequest{
				Description: stringPtr("Category without name"),
			},
			expectError: true,
			statusCode:  http.StatusBadRequest,
		},
		{
			name: "Duplicate slug",
			createReq: CreateCategoryRequest{
				Name: "Duplicate Category",
				Slug: stringPtr("duplicate-slug"),
			},
			expectError: true,
			statusCode:  http.StatusConflict,
		},
		{
			name: "Invalid parent ID",
			createReq: CreateCategoryRequest{
				Name:     "Invalid Parent Category",
				ParentID: stringPtr("invalid-parent-id"),
			},
			expectError: true,
			statusCode:  http.StatusBadRequest,
		},
		{
			name: "Circular reference",
			createReq: CreateCategoryRequest{
				Name:     "Circular Category",
				ParentID: stringPtr("self-reference"),
			},
			expectError: true,
			statusCode:  http.StatusBadRequest,
		},
		{
			name: "Invalid status",
			createReq: CreateCategoryRequest{
				Name:   "Invalid Status Category",
				Status: stringPtr("invalid_status"),
			},
			expectError: true,
			statusCode:  http.StatusBadRequest,
		},
		{
			name: "Invalid visibility",
			createReq: CreateCategoryRequest{
				Name:       "Invalid Visibility Category",
				Visibility: stringPtr("invalid_visibility"),
			},
			expectError: true,
			statusCode:  http.StatusBadRequest,
		},
	}

	// First create a category with specific slug for duplicate test
	firstCategory := CreateCategoryRequest{
		Name: "First Category",
		Slug: stringPtr("duplicate-slug"),
	}
	category := suite.createCategory(firstCategory)
	suite.createdCategories = append(suite.createdCategories, category.ID)

	// Run validation tests
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			if tc.expectError {
				suite.expectCategoryCreationError(tc.statusCode, tc.createReq)
			} else {
				cat := suite.createCategory(tc.createReq)
				suite.createdCategories = append(suite.createdCategories, cat.ID)
				assert.NotEmpty(t, cat.ID)
			}
		})
	}
}

// TestCategoryAnalytics tests category analytics functionality
func (suite *CategoryAPITestSuite) TestCategoryAnalytics() {
	// Create test category
	createReq := CreateCategoryRequest{
		Name:        "Analytics Test Category",
		Description: stringPtr("Category for analytics testing"),
		Status:      stringPtr("active"),
		Visibility:  stringPtr("public"),
	}

	category := suite.createCategory(createReq)
	categoryID := category.ID
	suite.createdCategories = append(suite.createdCategories, categoryID)

	// Test analytics endpoints
	analytics := suite.getCategoryAnalytics(categoryID, "last_30_days")
	assert.Equal(suite.T(), categoryID, analytics.CategoryID)
	assert.GreaterOrEqual(suite.T(), analytics.ProductCount, 0)
	assert.GreaterOrEqual(suite.T(), analytics.ViewCount, int64(0))
	assert.GreaterOrEqual(suite.T(), analytics.ConversionRate, 0.0)
	assert.Equal(suite.T(), "last_30_days", analytics.Period)

	// Test different time periods
	periods := []string{"last_7_days", "last_30_days", "last_90_days", "last_year"}
	for _, period := range periods {
		suite.T().Run(fmt.Sprintf("period_%s", period), func(t *testing.T) {
			analytics := suite.getCategoryAnalytics(categoryID, period)
			assert.Equal(t, categoryID, analytics.CategoryID)
			assert.Equal(t, period, analytics.Period)
		})
	}
}

// TestCategorySearch tests category search functionality
func (suite *CategoryAPITestSuite) TestCategorySearch() {
	// Create test categories
	categories := []CreateCategoryRequest{
		{
			Name:        "Home & Garden",
			Description: stringPtr("Home improvement and garden supplies"),
			Icon:        stringPtr("🏠"),
			Status:      stringPtr("active"),
		},
		{
			Name:        "Garden Tools",
			Description: stringPtr("Tools for gardening"),
			Icon:        stringPtr("🛠️"),
			Status:      stringPtr("active"),
		},
		{
			Name:        "Home Decor",
			Description: stringPtr("Decorative items for home"),
			Icon:        stringPtr("🖼️"),
			Status:      stringPtr("draft"),
		},
	}

	createdIDs := []string{}
	for _, req := range categories {
		cat := suite.createCategory(req)
		createdIDs = append(createdIDs, cat.ID)
		suite.createdCategories = append(suite.createdCategories, cat.ID)
	}

	// Test search scenarios
	testCases := []struct {
		name          string
		query         string
		status        string
		expectedCount int
	}{
		{
			name:          "Search by name",
			query:         "Garden",
			expectedCount: 2,
		},
		{
			name:          "Search by status active",
			status:        "active",
			expectedCount: 2,
		},
		{
			name:          "Search by status draft",
			status:        "draft",
			expectedCount: 1,
		},
		{
			name:          "Search by query and status",
			query:         "Home",
			status:        "active",
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			results := suite.searchCategories(tc.query, tc.status)
			assert.GreaterOrEqual(t, len(results.Categories), tc.expectedCount)
		})
	}
}

// TestCategoryConcurrency tests concurrent category operations
func (suite *CategoryAPITestSuite) TestCategoryConcurrency() {
	// Create a category for concurrent updates
	createReq := CreateCategoryRequest{
		Name:        "Concurrency Test Category",
		Description: stringPtr("Category for concurrency testing"),
		Status:      stringPtr("active"),
		SortOrder:   intPtr(10),
	}

	category := suite.createCategory(createReq)
	categoryID := category.ID
	suite.createdCategories = append(suite.createdCategories, categoryID)

	// Perform concurrent updates
	concurrency := 5
	resultChan := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(index int) {
			updateReq := UpdateCategoryRequest{
				SortOrder: intPtr(10 + index),
				Attributes: map[string]interface{}{
					"update_index": index,
					"timestamp":    time.Now().Unix(),
				},
			}

			err := suite.updateCategoryWithError(categoryID, updateReq)
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

	// Final category should be consistent
	finalCategory := suite.getCategory(categoryID)
	assert.GreaterOrEqual(suite.T(), finalCategory.SortOrder, 10)
	assert.LessOrEqual(suite.T(), finalCategory.SortOrder, 14)
}

// Helper methods

func (suite *CategoryAPITestSuite) createCategory(req CreateCategoryRequest) *Category {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/categories", suite.baseURL, suite.gatewayPort)

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

	var categoryResp CategoryResponse
	err = json.NewDecoder(resp.Body).Decode(&categoryResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), categoryResp.Success)
	require.NotNil(suite.T(), categoryResp.Category)

	return categoryResp.Category
}

func (suite *CategoryAPITestSuite) getCategory(categoryID string) *Category {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/categories/%s", suite.baseURL, suite.gatewayPort, categoryID)

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(suite.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var categoryResp CategoryResponse
	err = json.NewDecoder(resp.Body).Decode(&categoryResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), categoryResp.Success)
	require.NotNil(suite.T(), categoryResp.Category)

	return categoryResp.Category
}

func (suite *CategoryAPITestSuite) updateCategory(categoryID string, req UpdateCategoryRequest) *Category {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/categories/%s", suite.baseURL, suite.gatewayPort, categoryID)

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

	var categoryResp CategoryResponse
	err = json.NewDecoder(resp.Body).Decode(&categoryResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), categoryResp.Success)
	require.NotNil(suite.T(), categoryResp.Category)

	return categoryResp.Category
}

func (suite *CategoryAPITestSuite) updateCategoryWithError(categoryID string, req UpdateCategoryRequest) error {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/categories/%s", suite.baseURL, suite.gatewayPort, categoryID)

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

func (suite *CategoryAPITestSuite) deleteCategory(categoryID string) {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/categories/%s", suite.baseURL, suite.gatewayPort, categoryID)

	req, err := http.NewRequest("DELETE", url, nil)
	require.NoError(suite.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *CategoryAPITestSuite) listCategories() *CategoryResponse {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/categories", suite.baseURL, suite.gatewayPort)

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(suite.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var categoryResp CategoryResponse
	err = json.NewDecoder(resp.Body).Decode(&categoryResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), categoryResp.Success)

	return &categoryResp
}

func (suite *CategoryAPITestSuite) getCategoryHierarchy() *CategoryResponse {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/categories/hierarchy", suite.baseURL, suite.gatewayPort)

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(suite.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var categoryResp CategoryResponse
	err = json.NewDecoder(resp.Body).Decode(&categoryResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), categoryResp.Success)

	return &categoryResp
}

func (suite *CategoryAPITestSuite) searchCategories(query, status string) *CategoryResponse {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/categories/search", suite.baseURL, suite.gatewayPort)

	params := make(map[string]string)
	if query != "" {
		params["query"] = query
	}
	if status != "" {
		params["status"] = status
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

	var categoryResp CategoryResponse
	err = json.NewDecoder(resp.Body).Decode(&categoryResp)
	require.NoError(suite.T(), err)
	require.True(suite.T(), categoryResp.Success)

	return &categoryResp
}

func (suite *CategoryAPITestSuite) getCategoryAnalytics(categoryID, period string) *CategoryAnalytics {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/categories/%s/analytics?period=%s",
		suite.baseURL, suite.gatewayPort, categoryID, period)

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(suite.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var analytics CategoryAnalytics
	err = json.NewDecoder(resp.Body).Decode(&analytics)
	require.NoError(suite.T(), err)

	return &analytics
}

func (suite *CategoryAPITestSuite) expectCategoryNotFound(categoryID string) {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/categories/%s", suite.baseURL, suite.gatewayPort, categoryID)

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(suite.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer test-token-%s", suite.testUserID))

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *CategoryAPITestSuite) expectCategoryCreationError(expectedStatus int, req CreateCategoryRequest) {
	url := fmt.Sprintf("%s:%d/api/v1/commerce/categories", suite.baseURL, suite.gatewayPort)

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

func (suite *CategoryAPITestSuite) waitForServices() {
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

func (suite *CategoryAPITestSuite) cleanupTestData() {
	// Clean up created categories in reverse order (children first)
	for i := len(suite.createdCategories) - 1; i >= 0; i-- {
		suite.deleteCategory(suite.createdCategories[i])
	}
}

// Utility functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

// TestCategoryAPI runs the category API test suite
func TestCategoryAPI(t *testing.T) {
	suite.Run(t, new(CategoryAPITestSuite))
}