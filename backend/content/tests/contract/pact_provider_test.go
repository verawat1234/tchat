package contract

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// Import content service packages - adjust paths as needed
	"tchat.dev/content/handlers"
	contentModels "tchat.dev/content/models"
	"tchat.dev/content/services"
	"tchat.dev/content/utils"
)

// ProviderTestState holds the test application state
type ProviderTestState struct {
	db              *gorm.DB
	router          *gin.Engine
	server          *http.Server
	contentService  *services.ContentService
	contentHandlers *handlers.ContentHandlers
	testData        *TestContentData
}

// TestContentData holds test data for provider states
type TestContentData struct {
	ContentItems     []contentModels.ContentItem
	Categories       []contentModels.ContentCategory
	AuthenticatedUser uuid.UUID
}

// setupContentService initializes the content service for testing
func setupContentService(t *testing.T) *ProviderTestState {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup test database connection
	dbURL := getTestDatabaseURL()
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent, // Silent for tests
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: gormLogger,
	})
	require.NoError(t, err, "Failed to connect to test database")

	// Auto-migrate models
	err = db.AutoMigrate(
		&contentModels.ContentItem{},
		&contentModels.ContentCategory{},
		&contentModels.ContentVersion{},
	)
	require.NoError(t, err, "Failed to run migrations")

	// Initialize repositories
	contentRepo := services.NewPostgreSQLContentRepository(db)
	categoryRepo := services.NewPostgreSQLCategoryRepository(db)
	versionRepo := services.NewPostgreSQLVersionRepository(db)

	// Initialize services
	contentService := services.NewContentService(
		contentRepo,
		categoryRepo,
		versionRepo,
		db,
	)

	// Initialize handlers
	contentHandlers := handlers.NewContentHandlers(contentService)

	// Setup router
	router := gin.New()
	router.Use(gin.Recovery())

	// Add authentication middleware mock for testing
	router.Use(mockAuthMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		utils.SuccessResponse(c, gin.H{"status": "ok"})
	})

	// API routes
	v1 := router.Group("/api/v1")
	handlers.RegisterContentRoutes(v1, contentHandlers)

	return &ProviderTestState{
		db:              db,
		router:          router,
		contentService:  contentService,
		contentHandlers: contentHandlers,
		testData:        &TestContentData{},
	}
}

// mockAuthMiddleware provides mock authentication for testing
func mockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "Bearer mobile-token" {
			// Set authenticated user context
			c.Set("user_id", uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"))
			c.Set("authenticated", true)
		}
		c.Next()
	}
}

// getTestDatabaseURL returns the database URL for testing
func getTestDatabaseURL() string {
	// Use environment variable or default test database
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/tchat_content_test?sslmode=disable"
	}
	return dbURL
}

// startTestServer starts the HTTP server on a random available port
func (pts *ProviderTestState) startTestServer() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}

	port := listener.Addr().(*net.TCPAddr).Port

	pts.server = &http.Server{
		Handler: pts.router,
	}

	go func() {
		if err := pts.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("Test server error: %v", err)
		}
	}()

	return port, nil
}

// stopTestServer stops the HTTP server
func (pts *ProviderTestState) stopTestServer() error {
	if pts.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return pts.server.Shutdown(ctx)
	}
	return nil
}

// cleanupDatabase removes test data from the database
func (pts *ProviderTestState) cleanupDatabase() error {
	// Clean up in reverse order of dependencies
	tables := []string{
		"content_versions",
		"content_items",
		"content_categories",
	}

	for _, table := range tables {
		if err := pts.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	return nil
}

// setupProviderState sets up test data for specific provider states
func (pts *ProviderTestState) setupProviderState(state string) error {
	// Clean database before setting up state
	if err := pts.cleanupDatabase(); err != nil {
		return fmt.Errorf("failed to cleanup database: %w", err)
	}

	switch state {
	case "content items exist":
		return pts.setupContentItemsExistState()
	case "user is authenticated and can create content":
		return pts.setupAuthenticatedUserState()
	case "content exists in mobile category":
		return pts.setupMobileCategoryContentState()
	default:
		return fmt.Errorf("unknown provider state: %s", state)
	}
}

// setupContentItemsExistState creates test content items for pagination testing
func (pts *ProviderTestState) setupContentItemsExistState() error {

	// Create general category
	generalCategory := contentModels.ContentCategory{
		ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Name:        "general",
		Description: stringPtr("General content category"),
		IsActive:    true,
		SortOrder:   1,
	}
	if err := pts.db.Create(&generalCategory).Error; err != nil {
		return fmt.Errorf("failed to create general category: %w", err)
	}

	// Create test content items
	testItems := []contentModels.ContentItem{}
	for i := 1; i <= 15; i++ {
		itemID := uuid.New()
		item := contentModels.ContentItem{
			ID:       itemID,
			Category: "general",
			Type:     contentModels.ContentTypeText,
			Value: contentModels.ContentValue{
				"content": fmt.Sprintf("Sample content value %d", i),
			},
			Metadata: contentModels.ContentMetadata{
				"title":       fmt.Sprintf("Sample Title %d", i),
				"description": fmt.Sprintf("Sample description %d", i),
			},
			Status: contentModels.ContentStatusPublished,
			Tags:   []string{"mobile", "general"},
		}

		if err := pts.db.Create(&item).Error; err != nil {
			return fmt.Errorf("failed to create content item %d: %w", i, err)
		}

		testItems = append(testItems, item)
	}

	pts.testData.ContentItems = testItems
	log.Printf("Created %d test content items", len(testItems))
	return nil
}

// setupAuthenticatedUserState sets up authenticated user context
func (pts *ProviderTestState) setupAuthenticatedUserState() error {
	// Set authenticated user ID
	pts.testData.AuthenticatedUser = uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	// Create user-generated category
	userCategory := contentModels.ContentCategory{
		ID:          uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		Name:        "user-generated",
		Description: stringPtr("User generated content category"),
		IsActive:    true,
		SortOrder:   2,
	}
	if err := pts.db.Create(&userCategory).Error; err != nil {
		return fmt.Errorf("failed to create user-generated category: %w", err)
	}

	log.Printf("Set up authenticated user state with user ID: %s", pts.testData.AuthenticatedUser)
	return nil
}

// setupMobileCategoryContentState creates mobile-optimized content
func (pts *ProviderTestState) setupMobileCategoryContentState() error {

	// Create mobile category
	mobileCategory := contentModels.ContentCategory{
		ID:          uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		Name:        "mobile",
		Description: stringPtr("Mobile optimized content category"),
		IsActive:    true,
		SortOrder:   3,
	}
	if err := pts.db.Create(&mobileCategory).Error; err != nil {
		return fmt.Errorf("failed to create mobile category: %w", err)
	}

	// Create mobile-optimized content with image type
	mobileContentID := uuid.New()
	mobileContent := contentModels.ContentItem{
		ID:       mobileContentID,
		Category: "mobile",
		Type:     "image", // Note: This might need to be added to ContentType enum
		Value: contentModels.ContentValue{
			"url": "optimized-mobile-image.jpg",
		},
		Metadata: contentModels.ContentMetadata{
			"title":    "Mobile Optimized Image",
			"alt_text": "Mobile friendly image",
			"dimensions": map[string]interface{}{
				"width":  375,
				"height": 667,
			},
			"mobile_specific": map[string]interface{}{
				"retina_url":       "optimized-mobile-image@2x.jpg",
				"webp_url":         "optimized-mobile-image.webp",
				"loading_priority": "high",
			},
		},
		Status: contentModels.ContentStatusPublished,
		Tags:   []string{"mobile", "optimized", "image"},
	}

	if err := pts.db.Create(&mobileContent).Error; err != nil {
		return fmt.Errorf("failed to create mobile content: %w", err)
	}

	pts.testData.ContentItems = []contentModels.ContentItem{mobileContent}
	log.Printf("Created mobile-optimized content with ID: %s", mobileContentID)
	return nil
}

// TestContentServiceProviderVerification runs the Pact provider verification
func TestContentServiceProviderVerification(t *testing.T) {
	// Setup content service
	providerState := setupContentService(t)
	defer func() {
		if err := providerState.cleanupDatabase(); err != nil {
			t.Logf("Warning: Failed to cleanup database: %v", err)
		}
		if sqlDB, err := providerState.db.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Start test server
	port, err := providerState.startTestServer()
	require.NoError(t, err, "Failed to start test server")
	defer func() {
		if err := providerState.stopTestServer(); err != nil {
			t.Logf("Warning: Failed to stop test server: %v", err)
		}
	}()

	// Wait for server to be ready
	time.Sleep(100 * time.Millisecond)

	// Get the contract file path
	contractPath := getContractFilePath()
	t.Logf("Using contract file: %s", contractPath)

	// Configure Pact provider verification
	verifier := provider.NewVerifier()

	// Run verification
	err = verifier.VerifyProvider(t, provider.VerifyRequest{
		ProviderBaseURL: fmt.Sprintf("http://localhost:%d", port),
		Provider:        "content-service",

		// Pact file source
		PactFiles: []string{contractPath},

		// Provider state handlers
		StateHandlers: models.StateHandlers{
			"content items exist": func(setup bool, state models.ProviderState) (models.ProviderStateResponse, error) {
				if setup {
					err := providerState.setupProviderState("content items exist")
					return nil, err
				} else {
					err := providerState.cleanupDatabase()
					return nil, err
				}
			},
			"user is authenticated and can create content": func(setup bool, state models.ProviderState) (models.ProviderStateResponse, error) {
				if setup {
					err := providerState.setupProviderState("user is authenticated and can create content")
					return nil, err
				} else {
					err := providerState.cleanupDatabase()
					return nil, err
				}
			},
			"content exists in mobile category": func(setup bool, state models.ProviderState) (models.ProviderStateResponse, error) {
				if setup {
					err := providerState.setupProviderState("content exists in mobile category")
					return nil, err
				} else {
					err := providerState.cleanupDatabase()
					return nil, err
				}
			},
		},

		// Request filters for authentication
		RequestFilter: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Ensure the Authorization header is properly set for test requests
				if r.Header.Get("Authorization") == "Bearer mobile-token" {
					// Header is already correct, pass through
				}
				next.ServeHTTP(w, r)
			})
		},

		// Verification options
		PublishVerificationResults: false, // Set to true when integrating with Pact Broker
		ProviderVersion:           "1.0.0",
		BrokerURL:                 "", // Set when using Pact Broker
	})

	assert.NoError(t, err, "Provider verification should pass")
}

// getContractFilePath returns the path to the mobile consumer contract
func getContractFilePath() string {
	// Look for contract file in specs directory
	contractPath := filepath.Join("..", "..", "..", "specs", "021-implement-pact-contract", "contracts", "pact-consumer-mobile.json")

	// Check if file exists
	if _, err := os.Stat(contractPath); os.IsNotExist(err) {
		// Fallback to relative path
		contractPath = "../../../specs/021-implement-pact-contract/contracts/pact-consumer-mobile.json"
	}

	return contractPath
}


// TestProviderStateSetup tests the provider state setup functions independently
func TestProviderStateSetup(t *testing.T) {
	providerState := setupContentService(t)
	defer func() {
		if err := providerState.cleanupDatabase(); err != nil {
			t.Logf("Warning: Failed to cleanup database: %v", err)
		}
		if sqlDB, err := providerState.db.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	tests := []struct {
		name  string
		state string
	}{
		{"ContentItemsExist", "content items exist"},
		{"AuthenticatedUser", "user is authenticated and can create content"},
		{"MobileCategoryContent", "content exists in mobile category"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := providerState.setupProviderState(tt.state)
			assert.NoError(t, err, "Provider state setup should succeed")

			// Verify data was created
			switch tt.state {
			case "content items exist":
				var count int64
				err := providerState.db.Model(&contentModels.ContentItem{}).Count(&count).Error
				assert.NoError(t, err)
				assert.Greater(t, count, int64(0), "Content items should be created")
			case "user is authenticated and can create content":
				var categoryCount int64
				err := providerState.db.Model(&contentModels.ContentCategory{}).Where("name = ?", "user-generated").Count(&categoryCount).Error
				assert.NoError(t, err)
				assert.Equal(t, int64(1), categoryCount, "User-generated category should be created")
			case "content exists in mobile category":
				var mobileCount int64
				err := providerState.db.Model(&contentModels.ContentItem{}).Where("category = ?", "mobile").Count(&mobileCount).Error
				assert.NoError(t, err)
				assert.Greater(t, mobileCount, int64(0), "Mobile content should be created")
			}

			// Cleanup after each test
			err = providerState.cleanupDatabase()
			assert.NoError(t, err, "Cleanup should succeed")
		})
	}
}

// TestContentServiceHealthCheck tests the health check endpoint
func TestContentServiceHealthCheck(t *testing.T) {
	providerState := setupContentService(t)
	defer func() {
		if sqlDB, err := providerState.db.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	port, err := providerState.startTestServer()
	require.NoError(t, err)
	defer providerState.stopTestServer()

	// Wait for server to be ready
	time.Sleep(100 * time.Millisecond)

	// Test health check
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", port))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}