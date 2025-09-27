//go:build contracts
// +build contracts

package contract

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"tchat.dev/shared/config"
	"tchat.dev/shared/database"

	// Import all contract test suites
	authContracts "tchat.dev/auth/contracts"
	commerceContracts "tchat.dev/commerce/contracts"
	contentContracts "tchat.dev/content/contracts"
	notificationContracts "tchat.dev/notification/contracts"
)

// ContractIntegrationTestSuite runs all provider verification tests
type ContractIntegrationTestSuite struct {
	suite.Suite
	db     *sql.DB
	config *config.Config
}

// SetupSuite initializes the integration test environment
func (suite *ContractIntegrationTestSuite) SetupSuite() {
	// Load test configuration
	suite.config = config.LoadTestConfig()
	if suite.config == nil {
		suite.config = config.LoadConfig()
	}

	// Initialize test database
	var err error
	suite.db, err = database.NewConnection(suite.config.Database)
	suite.Require().NoError(err, "Failed to connect to test database")

	// Setup test data and infrastructure
	suite.setupTestInfrastructure()

	// Wait for services to be ready
	suite.waitForServicesReady()
}

// TearDownSuite cleans up the integration test environment
func (suite *ContractIntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.cleanupTestData()
		suite.db.Close()
	}
}

// TestAllProviderContracts runs all provider verification tests
func (suite *ContractIntegrationTestSuite) TestAllProviderContracts() {
	suite.T().Run("AuthServiceContracts", func(t *testing.T) {
		authSuite := new(authContracts.AuthServiceProviderTestSuite)
		suite.Run(authSuite, t)
	})

	suite.T().Run("ContentServiceContracts", func(t *testing.T) {
		contentSuite := new(contentContracts.ContentServiceProviderTestSuite)
		suite.Run(contentSuite, t)
	})

	suite.T().Run("CommerceServiceContracts", func(t *testing.T) {
		commerceSuite := new(commerceContracts.CommerceServiceProviderTestSuite)
		suite.Run(commerceSuite, t)
	})

	suite.T().Run("NotificationServiceContracts", func(t *testing.T) {
		notificationSuite := new(notificationContracts.NotificationServiceProviderTestSuite)
		suite.Run(notificationSuite, t)
	})
}

// TestAuthServiceContracts specifically tests auth service contracts
func (suite *ContractIntegrationTestSuite) TestAuthServiceContracts() {
	authSuite := new(authContracts.AuthServiceProviderTestSuite)
	suite.Run(authSuite, suite.T())
}

// TestContentServiceContracts specifically tests content service contracts
func (suite *ContractIntegrationTestSuite) TestContentServiceContracts() {
	contentSuite := new(contentContracts.ContentServiceProviderTestSuite)
	suite.Run(contentSuite, suite.T())
}

// TestCommerceServiceContracts specifically tests commerce service contracts
func (suite *ContractIntegrationTestSuite) TestCommerceServiceContracts() {
	commerceSuite := new(commerceContracts.CommerceServiceProviderTestSuite)
	suite.Run(commerceSuite, suite.T())
}

// TestNotificationServiceContracts specifically tests notification service contracts
func (suite *ContractIntegrationTestSuite) TestNotificationServiceContracts() {
	notificationSuite := new(notificationContracts.NotificationServiceProviderTestSuite)
	suite.Run(notificationSuite, suite.T())
}

// setupTestInfrastructure prepares the test environment
func (suite *ContractIntegrationTestSuite) setupTestInfrastructure() {
	ctx := context.Background()

	// Create test database schema if needed
	suite.createTestSchema(ctx)

	// Setup test data that's common across all services
	suite.setupCommonTestData(ctx)

	// Setup service-specific test infrastructure
	suite.setupServiceInfrastructure(ctx)
}

// createTestSchema creates database schema for testing
func (suite *ContractIntegrationTestSuite) createTestSchema(ctx context.Context) {
	// Create basic test tables that might be needed across services
	schemas := []string{
		`CREATE TABLE IF NOT EXISTS test_users (
			id UUID PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE,
			phone_number VARCHAR(20) UNIQUE,
			country VARCHAR(2),
			kyc_tier VARCHAR(20),
			is_verified BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS test_sessions (
			id UUID PRIMARY KEY,
			user_id UUID REFERENCES test_users(id),
			device_id VARCHAR(255),
			status VARCHAR(20),
			expires_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS test_content (
			id VARCHAR(255) PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			type VARCHAR(50),
			category VARCHAR(100),
			status VARCHAR(20),
			version INTEGER DEFAULT 1,
			author_id UUID,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW(),
			published_at TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_test_content_status ON test_content(status);`,
		`CREATE INDEX IF NOT EXISTS idx_test_content_category ON test_content(category);`,
	}

	for _, schema := range schemas {
		_, err := suite.db.ExecContext(ctx, schema)
		if err != nil {
			suite.T().Logf("Warning: Failed to create test schema: %v", err)
		}
	}
}

// setupCommonTestData creates test data used across multiple services
func (suite *ContractIntegrationTestSuite) setupCommonTestData(ctx context.Context) {
	// Create test users
	testUsers := []map[string]interface{}{
		{
			"id":           "123e4567-e89b-12d3-a456-426614174000",
			"name":         "John Doe",
			"email":        "john.doe@example.com",
			"phone_number": "+66812345678",
			"country":      "TH",
			"kyc_tier":     "basic",
			"is_verified":  true,
		},
		{
			"id":           "456e7890-e89b-12d3-a456-426614174001",
			"name":         "Jane Smith",
			"email":        "jane.smith@example.com",
			"phone_number": "+66987654321",
			"country":      "TH",
			"kyc_tier":     "premium",
			"is_verified":  true,
		},
	}

	for _, user := range testUsers {
		_, err := suite.db.ExecContext(ctx,
			`INSERT INTO test_users (id, name, email, phone_number, country, kyc_tier, is_verified)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)
			 ON CONFLICT (id) DO NOTHING`,
			user["id"], user["name"], user["email"], user["phone_number"],
			user["country"], user["kyc_tier"], user["is_verified"])
		if err != nil {
			suite.T().Logf("Warning: Failed to create test user: %v", err)
		}
	}

	// Create test content
	testContent := []map[string]interface{}{
		{
			"id":           "content_123",
			"title":        "Sample Content",
			"type":         "text",
			"category":     "general",
			"status":       "published",
			"version":      1,
			"author_id":    "123e4567-e89b-12d3-a456-426614174000",
			"published_at": time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			"id":           "content_456",
			"title":        "Featured Content",
			"type":         "text",
			"category":     "featured",
			"status":       "published",
			"version":      1,
			"author_id":    "456e7890-e89b-12d3-a456-426614174001",
			"published_at": time.Date(2025, 1, 1, 1, 0, 0, 0, time.UTC),
		},
	}

	for _, content := range testContent {
		_, err := suite.db.ExecContext(ctx,
			`INSERT INTO test_content (id, title, type, category, status, version, author_id, published_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			 ON CONFLICT (id) DO NOTHING`,
			content["id"], content["title"], content["type"], content["category"],
			content["status"], content["version"], content["author_id"], content["published_at"])
		if err != nil {
			suite.T().Logf("Warning: Failed to create test content: %v", err)
		}
	}
}

// setupServiceInfrastructure sets up infrastructure specific to each service
func (suite *ContractIntegrationTestSuite) setupServiceInfrastructure(ctx context.Context) {
	// This could include setting up service-specific test data,
	// mock external dependencies, or service configurations
	suite.T().Log("Setting up service-specific test infrastructure...")

	// Setup auth service infrastructure
	suite.setupAuthInfrastructure(ctx)

	// Setup content service infrastructure
	suite.setupContentInfrastructure(ctx)

	// Setup commerce service infrastructure
	suite.setupCommerceInfrastructure(ctx)

	// Setup notification service infrastructure
	suite.setupNotificationInfrastructure(ctx)
}

// setupAuthInfrastructure sets up auth-specific test infrastructure
func (suite *ContractIntegrationTestSuite) setupAuthInfrastructure(ctx context.Context) {
	// Create test sessions for auth testing
	_, err := suite.db.ExecContext(ctx,
		`INSERT INTO test_sessions (id, user_id, device_id, status, expires_at)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (id) DO NOTHING`,
		"456e7890-e89b-12d3-a456-426614174001",
		"123e4567-e89b-12d3-a456-426614174000",
		"device_web_123",
		"active",
		time.Now().Add(7*24*time.Hour))
	if err != nil {
		suite.T().Logf("Warning: Failed to create test session: %v", err)
	}
}

// setupContentInfrastructure sets up content-specific test infrastructure
func (suite *ContractIntegrationTestSuite) setupContentInfrastructure(ctx context.Context) {
	// Content-specific setup could include creating categories, templates, etc.
	suite.T().Log("Content service infrastructure ready")
}

// setupCommerceInfrastructure sets up commerce-specific test infrastructure
func (suite *ContractIntegrationTestSuite) setupCommerceInfrastructure(ctx context.Context) {
	// Commerce-specific setup could include creating test shops, products, etc.
	suite.T().Log("Commerce service infrastructure ready")
}

// setupNotificationInfrastructure sets up notification-specific test infrastructure
func (suite *ContractIntegrationTestSuite) setupNotificationInfrastructure(ctx context.Context) {
	// Notification-specific setup could include creating templates, channels, etc.
	suite.T().Log("Notification service infrastructure ready")
}

// waitForServicesReady waits for all services to be ready for testing
func (suite *ContractIntegrationTestSuite) waitForServicesReady() {
	// In a real environment, this might check service health endpoints
	// For now, we'll just add a brief delay to ensure everything is initialized
	time.Sleep(2 * time.Second)
	suite.T().Log("All services ready for contract testing")
}

// cleanupTestData removes test data after all tests complete
func (suite *ContractIntegrationTestSuite) cleanupTestData() {
	ctx := context.Background()

	// Clean up in reverse order of creation
	cleanupQueries := []string{
		"DELETE FROM test_sessions WHERE device_id LIKE 'test_%' OR device_id = 'device_web_123'",
		"DELETE FROM test_content WHERE id LIKE 'content_%'",
		"DELETE FROM test_users WHERE phone_number IN ('+66812345678', '+66987654321')",
	}

	for _, query := range cleanupQueries {
		_, err := suite.db.ExecContext(ctx, query)
		if err != nil {
			suite.T().Logf("Warning: Cleanup query failed: %v", err)
		}
	}

	suite.T().Log("Test data cleanup completed")
}

// Helper method to run individual test suites
func (suite *ContractIntegrationTestSuite) Run(testSuite suite.TestingSuite, t *testing.T) {
	// This is a custom runner that can add integration-specific setup
	suite.T().Logf("Running contract test suite: %T", testSuite)
	suite.Suite.Run(t, testSuite)
}

// TestContractIntegrationSuite runs the integration contract test suite
func TestContractIntegrationSuite(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("Skipping contract integration tests in short mode")
	}

	// Check if required environment is available
	if os.Getenv("SKIP_CONTRACT_TESTS") == "true" {
		t.Skip("Contract tests disabled by SKIP_CONTRACT_TESTS environment variable")
	}

	// Run the suite
	suite.Run(t, new(ContractIntegrationTestSuite))
}

// Individual service test functions for more granular testing

// TestAuthServiceProviderContracts tests just the auth service contracts
func TestAuthServiceProviderContracts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping auth contract tests in short mode")
	}

	suite.Run(t, new(authContracts.AuthServiceProviderTestSuite))
}

// TestContentServiceProviderContracts tests just the content service contracts
func TestContentServiceProviderContracts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping content contract tests in short mode")
	}

	suite.Run(t, new(contentContracts.ContentServiceProviderTestSuite))
}

// TestCommerceServiceProviderContracts tests just the commerce service contracts
func TestCommerceServiceProviderContracts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping commerce contract tests in short mode")
	}

	suite.Run(t, new(commerceContracts.CommerceServiceProviderTestSuite))
}

// TestNotificationServiceProviderContracts tests just the notification service contracts
func TestNotificationServiceProviderContracts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping notification contract tests in short mode")
	}

	suite.Run(t, new(notificationContracts.NotificationServiceProviderTestSuite))
}