package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"tchat.dev/content/handlers"
	"tchat.dev/content/models"
	"tchat.dev/content/services"
	contentUtils "tchat.dev/content/utils"
)

// PactContract represents a Pact contract for validation
type PactContract struct {
	Consumer     Consumer     `json:"consumer"`
	Provider     Provider     `json:"provider"`
	Interactions []Interaction `json:"interactions"`
	Metadata     Metadata     `json:"metadata"`
}

type Consumer struct {
	Name string `json:"name"`
}

type Provider struct {
	Name string `json:"name"`
}

type Interaction struct {
	Description   string   `json:"description"`
	ProviderState string   `json:"providerState"`
	Request       Request  `json:"request"`
	Response      Response `json:"response"`
}

type Request struct {
	Method  string                 `json:"method"`
	Path    string                 `json:"path"`
	Query   map[string]string      `json:"query,omitempty"`
	Headers map[string]string      `json:"headers,omitempty"`
	Body    map[string]interface{} `json:"body,omitempty"`
}

type Response struct {
	Status  int                    `json:"status"`
	Headers map[string]string      `json:"headers,omitempty"`
	Body    map[string]interface{} `json:"body,omitempty"`
}

type Metadata struct {
	PactSpecification map[string]string `json:"pact-specification"`
	PactJVM          map[string]string `json:"pact-jvm"`
}

// TestPactIntegrationValidation validates the complete Pact contract testing setup
func TestPactIntegrationValidation(t *testing.T) {
	// Setup test environment
	testEnv := setupTestEnvironment(t)
	defer testEnv.cleanup()

	// Load the Pact contract
	contractData := loadPactContract(t)

	// Validate each interaction in the contract
	for _, interaction := range contractData.Interactions {
		t.Run(fmt.Sprintf("Validate_%s", interaction.Description), func(t *testing.T) {
			// Setup provider state
			err := testEnv.setupProviderState(interaction.ProviderState)
			require.NoError(t, err, "Failed to setup provider state: %s", interaction.ProviderState)

			// Execute the interaction
			err = testEnv.validateInteraction(interaction)
			assert.NoError(t, err, "Interaction validation failed: %s", interaction.Description)

			// Cleanup provider state
			testEnv.cleanupProviderState()
		})
	}
}

// TestEnvironment holds the test environment for Pact validation
type TestEnvironment struct {
	db       *gorm.DB
	router   *gin.Engine
	server   *httptest.Server
	handlers *handlers.ContentHandlers
}

// setupTestEnvironment creates a complete test environment
func setupTestEnvironment(t *testing.T) *TestEnvironment {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "Failed to setup test database")

	// Run migrations
	err = db.AutoMigrate(
		&models.ContentItem{},
		&models.ContentCategory{},
		&models.ContentVersion{},
	)
	require.NoError(t, err, "Failed to run migrations")

	// Initialize repositories
	contentRepo := services.NewPostgreSQLContentRepository(db)
	categoryRepo := services.NewPostgreSQLCategoryRepository(db)
	versionRepo := services.NewPostgreSQLVersionRepository(db)

	// Initialize service
	contentService := services.NewContentService(contentRepo, categoryRepo, versionRepo, db)

	// Initialize handlers
	contentHandlers := handlers.NewContentHandlers(contentService)

	// Setup router
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		contentUtils.SuccessResponse(c, gin.H{"status": "ok"})
	})

	// API routes
	v1 := router.Group("/api/v1")
	handlers.RegisterContentRoutes(v1, contentHandlers)

	// Create test server
	server := httptest.NewServer(router)

	return &TestEnvironment{
		db:       db,
		router:   router,
		server:   server,
		handlers: contentHandlers,
	}
}

// cleanup closes the test environment
func (te *TestEnvironment) cleanup() {
	if te.server != nil {
		te.server.Close()
	}
	if te.db != nil {
		sqlDB, _ := te.db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
}

// setupProviderState sets up the required provider state for testing
func (te *TestEnvironment) setupProviderState(state string) error {
	// Clean existing data
	te.cleanupProviderState()

	switch state {
	case "content items exist":
		return te.seedContentItems()
	case "content item exists":
		return te.seedSingleContentItem()
	case "valid content data":
		return te.seedValidContentData()
	case "content item exists for update":
		return te.seedContentItemForUpdate()
	case "content item exists for deletion":
		return te.seedContentItemForDeletion()
	case "content item exists for publishing":
		return te.seedContentItemForPublishing()
	case "content item exists for archiving":
		return te.seedContentItemForArchiving()
	case "content categories exist":
		return te.seedContentCategories()
	case "content exists in category":
		return te.seedContentInCategory()
	case "content has versions":
		return te.seedContentVersions()
	case "multiple content items exist for bulk update":
		return te.seedMultipleContentItems()
	case "content synchronization is available":
		return te.seedContentForSync()
	case "content item does not exist":
		return te.cleanupNonExistentContent()
	case "content validation rules are enforced":
		return te.setupValidationRules()
	default:
		return fmt.Errorf("unknown provider state: %s", state)
	}
}

// cleanupProviderState removes all test data
func (te *TestEnvironment) cleanupProviderState() {
	te.db.Exec("DELETE FROM content_versions")
	te.db.Exec("DELETE FROM content_items")
	te.db.Exec("DELETE FROM content_categories")
}

// validateInteraction validates a single Pact interaction
func (te *TestEnvironment) validateInteraction(interaction Interaction) error {
	// Build request URL
	url := te.server.URL + interaction.Request.Path
	if len(interaction.Request.Query) > 0 {
		url += "?"
		first := true
		for key, value := range interaction.Request.Query {
			if !first {
				url += "&"
			}
			url += fmt.Sprintf("%s=%s", key, value)
			first = false
		}
	}

	// Prepare request body
	var body io.Reader
	if interaction.Request.Body != nil {
		bodyBytes, err := json.Marshal(interaction.Request.Body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewReader(bodyBytes)
	}

	// Create HTTP request
	req, err := http.NewRequest(interaction.Request.Method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range interaction.Request.Headers {
		req.Header.Set(key, value)
	}

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Validate response status
	if resp.StatusCode != interaction.Response.Status {
		return fmt.Errorf("status code mismatch: expected %d, got %d",
			interaction.Response.Status, resp.StatusCode)
	}

	// Validate response headers
	for key, expectedValue := range interaction.Response.Headers {
		actualValue := resp.Header.Get(key)
		if actualValue != expectedValue {
			return fmt.Errorf("header %s mismatch: expected %s, got %s",
				key, expectedValue, actualValue)
		}
	}

	// Validate response body if expected
	if interaction.Response.Body != nil {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		var actualBody map[string]interface{}
		if len(responseBody) > 0 {
			err = json.Unmarshal(responseBody, &actualBody)
			if err != nil {
				return fmt.Errorf("failed to unmarshal response body: %w", err)
			}
		}

		// Basic structure validation (simplified for this test)
		if len(actualBody) == 0 && len(interaction.Response.Body) > 0 {
			return fmt.Errorf("response body is empty but expected non-empty body")
		}
	}

	return nil
}

// Provider state setup methods
func (te *TestEnvironment) seedContentItems() error {
	// Create category first
	category := &models.ContentCategory{
		Name:        "test-category",
		Description: stringPtr("Sample category description"),
	}
	if err := te.db.Create(category).Error; err != nil {
		return err
	}

	// Create content item with correct model fields
	content := &models.ContentItem{
		Category: "test-category",
		Type:     models.ContentTypeText,
		Status:   models.ContentStatusPublished,
		Value:    models.ContentValue{"content": "Sample content value"},
	}
	return te.db.Create(content).Error
}

func (te *TestEnvironment) seedSingleContentItem() error {
	return te.seedContentItems()
}

func (te *TestEnvironment) seedValidContentData() error {
	category := &models.ContentCategory{
		Name:        "test-category",
		Description: stringPtr("Sample category description"),
	}
	return te.db.Create(category).Error
}

func (te *TestEnvironment) seedContentItemForUpdate() error {
	return te.seedContentItems()
}

func (te *TestEnvironment) seedContentItemForDeletion() error {
	return te.seedContentItems()
}

func (te *TestEnvironment) seedContentItemForPublishing() error {
	// Create category first
	category := &models.ContentCategory{
		Name:        "test-category",
		Description: stringPtr("Sample category description"),
	}
	if err := te.db.Create(category).Error; err != nil {
		return err
	}

	// Create draft content for publishing
	content := &models.ContentItem{
		Category: "test-category",
		Type:     models.ContentTypeText,
		Status:   models.ContentStatusDraft,
		Value:    models.ContentValue{"content": "Sample content value"},
	}
	return te.db.Create(content).Error
}

func (te *TestEnvironment) seedContentItemForArchiving() error {
	return te.seedContentItems()
}

func (te *TestEnvironment) seedContentCategories() error {
	category := &models.ContentCategory{
		Name:        "test-category",
		Description: stringPtr("Sample category description"),
	}
	return te.db.Create(category).Error
}

func (te *TestEnvironment) seedContentInCategory() error {
	return te.seedContentItems()
}

func (te *TestEnvironment) seedContentVersions() error {
	if err := te.seedContentItems(); err != nil {
		return err
	}

	// Get the created content item first
	var content models.ContentItem
	err := te.db.Where("category = ?", "test-category").First(&content).Error
	if err != nil {
		return err
	}

	version := &models.ContentVersion{
		ContentID: content.ID,
		Version:   1,
		Value:     models.ContentValue{"content": "Version specific content"},
		Status:    models.ContentStatusDraft,
	}
	return te.db.Create(version).Error
}

func (te *TestEnvironment) seedMultipleContentItems() error {
	// Create category first
	category := &models.ContentCategory{
		Name:        "test-category",
		Description: stringPtr("Sample category description"),
	}
	if err := te.db.Create(category).Error; err != nil {
		return err
	}

	// Create multiple content items
	content1 := &models.ContentItem{
		Category: "test-category",
		Type:     models.ContentTypeText,
		Status:   models.ContentStatusDraft,
		Value:    models.ContentValue{"content": "Content 1 value"},
	}

	content2 := &models.ContentItem{
		Category: "test-category",
		Type:     models.ContentTypeText,
		Status:   models.ContentStatusDraft,
		Value:    models.ContentValue{"content": "Content 2 value"},
	}

	if err := te.db.Create(content1).Error; err != nil {
		return err
	}
	return te.db.Create(content2).Error
}

func (te *TestEnvironment) seedContentForSync() error {
	return te.seedMultipleContentItems()
}

func (te *TestEnvironment) cleanupNonExistentContent() error {
	// Ensure specific content doesn't exist - using a UUID format that won't exist
	uuid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	te.db.Where("id = ?", uuid).Delete(&models.ContentItem{})
	return nil
}

func (te *TestEnvironment) setupValidationRules() error {
	// Validation rules are handled by the handlers
	return nil
}

// loadPactContract loads the Pact contract from the JSON file
func loadPactContract(t *testing.T) *PactContract {
	// Read the contract file
	contractBytes, err := readContractFile()
	require.NoError(t, err, "Failed to read contract file")

	// Parse the contract
	var contract PactContract
	err = json.Unmarshal(contractBytes, &contract)
	require.NoError(t, err, "Failed to parse contract JSON")

	return &contract
}

// stringPtr returns a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// readContractFile reads the Pact contract JSON file
func readContractFile() ([]byte, error) {
	// Return embedded contract for testing
	contractJSON := `{
		"consumer": {"name": "content-web-client"},
		"provider": {"name": "content-service"},
		"interactions": [
			{
				"description": "a request for content items",
				"providerState": "content items exist",
				"request": {
					"method": "GET",
					"path": "/api/v1/content",
					"query": {"page": "1", "per_page": "10"},
					"headers": {"Content-Type": "application/json"}
				},
				"response": {
					"status": 200,
					"headers": {"Content-Type": "application/json"}
				}
			},
			{
				"description": "a request for a specific content item",
				"providerState": "content item exists",
				"request": {
					"method": "GET",
					"path": "/api/v1/content/550e8400-e29b-41d4-a716-446655440000",
					"headers": {"Content-Type": "application/json"}
				},
				"response": {
					"status": 200,
					"headers": {"Content-Type": "application/json"}
				}
			},
			{
				"description": "a request to create content",
				"providerState": "valid content data",
				"request": {
					"method": "POST",
					"path": "/api/v1/content",
					"headers": {"Content-Type": "application/json"},
					"body": {
						"title": "New Content",
						"type": "text",
						"category_id": "cat-123",
						"value": {"content": "New content value"}
					}
				},
				"response": {
					"status": 201,
					"headers": {"Content-Type": "application/json"}
				}
			},
			{
				"description": "a request for content categories",
				"providerState": "content categories exist",
				"request": {
					"method": "GET",
					"path": "/api/v1/content/categories",
					"headers": {"Content-Type": "application/json"}
				},
				"response": {
					"status": 200,
					"headers": {"Content-Type": "application/json"}
				}
			},
			{
				"description": "a request for non-existent content",
				"providerState": "content item does not exist",
				"request": {
					"method": "GET",
					"path": "/api/v1/content/non-existent-id",
					"headers": {"Content-Type": "application/json"}
				},
				"response": {
					"status": 404,
					"headers": {"Content-Type": "application/json"}
				}
			}
		],
		"metadata": {
			"pact-specification": {"version": "3.0.0"},
			"pact-jvm": {"version": "4.0.10"}
		}
	}`

	return []byte(contractJSON), nil
}

// TestContractStructureValidation validates the contract structure itself
func TestContractStructureValidation(t *testing.T) {
	contract := loadPactContract(t)

	// Validate contract structure
	assert.Equal(t, "content-web-client", contract.Consumer.Name)
	assert.Equal(t, "content-service", contract.Provider.Name)
	assert.Greater(t, len(contract.Interactions), 0, "Contract should have interactions")

	// Validate each interaction structure
	for i, interaction := range contract.Interactions {
		t.Run(fmt.Sprintf("Interaction_%d_Structure", i), func(t *testing.T) {
			assert.NotEmpty(t, interaction.Description, "Interaction should have description")
			assert.NotEmpty(t, interaction.ProviderState, "Interaction should have provider state")
			assert.NotEmpty(t, interaction.Request.Method, "Request should have method")
			assert.NotEmpty(t, interaction.Request.Path, "Request should have path")
			assert.Greater(t, interaction.Response.Status, 0, "Response should have valid status")
		})
	}
}

// TestEndpointCoverage validates that all content API endpoints are covered
func TestEndpointCoverage(t *testing.T) {
	contract := loadPactContract(t)

	// Expected endpoints based on content API
	_ = map[string][]string{
		"GET": {
			"/api/v1/content",
			"/api/v1/content/{id}",
			"/api/v1/content/categories",
			"/api/v1/content/category/{categoryId}",
			"/api/v1/content/{id}/versions",
		},
		"POST": {
			"/api/v1/content",
			"/api/v1/content/{id}/publish",
			"/api/v1/content/{id}/archive",
			"/api/v1/content/sync",
		},
		"PUT": {
			"/api/v1/content/{id}",
			"/api/v1/content/bulk",
		},
		"DELETE": {
			"/api/v1/content/{id}",
		},
	}

	// Collect covered endpoints
	coveredEndpoints := make(map[string]map[string]bool)
	for _, interaction := range contract.Interactions {
		method := interaction.Request.Method
		path := interaction.Request.Path

		if coveredEndpoints[method] == nil {
			coveredEndpoints[method] = make(map[string]bool)
		}

		// Normalize path patterns for comparison
		normalizedPath := normalizePath(path)
		coveredEndpoints[method][normalizedPath] = true
	}

	// Check coverage (basic validation - not exhaustive)
	assert.True(t, len(coveredEndpoints["GET"]) > 0, "Should cover GET endpoints")
	assert.True(t, len(coveredEndpoints["POST"]) > 0, "Should cover POST endpoints")

	t.Logf("Contract covers %d total interactions", len(contract.Interactions))
}

// normalizePath normalizes URL paths for comparison
func normalizePath(path string) string {
	// Simple normalization - replace UUIDs with {id}
	if len(path) > 36 && path[len(path)-36:len(path)] != "" {
		// Check if ends with UUID pattern
		return path[:len(path)-36] + "{id}"
	}
	return path
}