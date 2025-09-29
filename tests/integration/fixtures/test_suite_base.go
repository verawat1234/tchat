package fixtures

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// BaseTestSuite provides common test setup and teardown functionality
type BaseTestSuite struct {
	suite.Suite
	DataManager    *TestDataManager
	TestStartTime  time.Time
	TestName       string
	SkipTeardown   bool
	PreserveData   bool
	FixturesLoaded []string
}

// TestConfig holds configuration for test execution
type TestConfig struct {
	DatabaseURL     string
	MigrationsDir   string
	FixturesDir     string
	TestDBPrefix    string
	SkipTeardown    bool
	PreserveData    bool
	LoadFixtures    []string
	MaxTestDuration time.Duration
}

// GetDefaultTestConfig returns default test configuration
func GetDefaultTestConfig() TestConfig {
	return TestConfig{
		DatabaseURL:     getEnvOrDefault("TEST_DATABASE_URL", "postgres://localhost/tchat_test?sslmode=disable"),
		MigrationsDir:   getEnvOrDefault("MIGRATIONS_DIR", "../../../backend/migrations"),
		FixturesDir:     getEnvOrDefault("FIXTURES_DIR", "./fixtures"),
		TestDBPrefix:    getEnvOrDefault("TEST_DB_PREFIX", "tchat"),
		SkipTeardown:    getEnvOrDefault("SKIP_TEARDOWN", "false") == "true",
		PreserveData:    getEnvOrDefault("PRESERVE_DATA", "false") == "true",
		LoadFixtures:    []string{"users_fixtures", "commerce_fixtures", "businesses_fixtures", "products_fixtures"},
		MaxTestDuration: 5 * time.Minute,
	}
}

// SetupSuite runs once before all tests in the suite
func (s *BaseTestSuite) SetupSuite() {
	s.TestStartTime = time.Now()
	s.TestName = s.T().Name()

	log.Printf("=== Starting Test Suite: %s ===", s.TestName)

	config := GetDefaultTestConfig()

	// Create test data manager
	dataManager, err := NewTestDataManager(TestDataConfig{
		DatabaseURL:   config.DatabaseURL,
		MigrationsDir: config.MigrationsDir,
		FixturesDir:   config.FixturesDir,
		TestDBPrefix:  config.TestDBPrefix,
	})
	s.Require().NoError(err, "Failed to create test data manager")

	s.DataManager = dataManager
	s.SkipTeardown = config.SkipTeardown
	s.PreserveData = config.PreserveData

	// Setup database schema
	err = s.DataManager.SetupDatabase()
	s.Require().NoError(err, "Failed to setup database")

	// Load initial fixtures
	if len(config.LoadFixtures) > 0 {
		err = s.DataManager.LoadFixtures(config.LoadFixtures...)
		s.Require().NoError(err, "Failed to load fixtures")
		s.FixturesLoaded = config.LoadFixtures
	}

	log.Printf("Test suite setup completed in %v", time.Since(s.TestStartTime))
}

// TearDownSuite runs once after all tests in the suite
func (s *BaseTestSuite) TearDownSuite() {
	totalDuration := time.Since(s.TestStartTime)
	log.Printf("=== Finishing Test Suite: %s (Duration: %v) ===", s.TestName, totalDuration)

	if s.DataManager != nil && !s.SkipTeardown {
		err := s.DataManager.TeardownDatabase()
		if err != nil {
			log.Printf("Warning: Failed to teardown database: %v", err)
		}
	}

	log.Printf("Test suite teardown completed")
}

// SetupTest runs before each individual test
func (s *BaseTestSuite) SetupTest() {
	testName := s.T().Name()
	log.Printf("--- Starting Test: %s ---", testName)

	if !s.PreserveData {
		// Clean up data from previous test but preserve schema
		err := s.DataManager.CleanupAllTables()
		s.Require().NoError(err, "Failed to cleanup tables")

		err = s.DataManager.ResetSequences()
		s.Require().NoError(err, "Failed to reset sequences")

		// Reload fixtures for fresh test data
		if len(s.FixturesLoaded) > 0 {
			err = s.DataManager.LoadFixtures(s.FixturesLoaded...)
			s.Require().NoError(err, "Failed to reload fixtures")
		}
	}
}

// TearDownTest runs after each individual test
func (s *BaseTestSuite) TearDownTest() {
	testName := s.T().Name()

	// Log test result
	if s.T().Failed() {
		log.Printf("--- Test FAILED: %s ---", testName)
		s.logDatabaseState()
	} else {
		log.Printf("--- Test PASSED: %s ---", testName)
	}
}

// LoadAdditionalFixtures loads additional fixtures during test execution
func (s *BaseTestSuite) LoadAdditionalFixtures(fixtureNames ...string) {
	err := s.DataManager.LoadFixtures(fixtureNames...)
	s.Require().NoError(err, "Failed to load additional fixtures")
}

// CleanupTables cleans specific tables during test execution
func (s *BaseTestSuite) CleanupTables(tableNames ...string) {
	for _, tableName := range tableNames {
		err := s.DataManager.CleanupTable(tableName)
		s.Require().NoError(err, "Failed to cleanup table: %s", tableName)
	}
}

// GetTestTransaction returns a database transaction for isolated testing
func (s *BaseTestSuite) GetTestTransaction() context.Context {
	tx, err := s.DataManager.GetTestTransaction()
	s.Require().NoError(err, "Failed to get test transaction")

	// Create context with transaction
	ctx := context.Background()
	return context.WithValue(ctx, "tx", tx)
}

// CreateTestEntities provides helper methods to create test entities
func (s *BaseTestSuite) CreateTestEntities() *TestEntityFactory {
	return &TestEntityFactory{
		DataManager: s.DataManager,
		Suite:       s,
	}
}

// AssertDatabaseState validates database state during tests
func (s *BaseTestSuite) AssertDatabaseState() *DatabaseAssertion {
	return &DatabaseAssertion{
		DataManager: s.DataManager,
		Suite:       s,
	}
}

// logDatabaseState logs current database state for debugging
func (s *BaseTestSuite) logDatabaseState() {
	log.Printf("=== Database State Debug ===")

	// Log table counts
	tables := []string{"users", "businesses", "categories", "products", "carts", "cart_items"}
	for _, table := range tables {
		count := s.getTableCount(table)
		log.Printf("Table %s: %d rows", table, count)
	}
}

// getTableCount returns the number of rows in a table
func (s *BaseTestSuite) getTableCount(tableName string) int {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := s.DataManager.DB.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Error counting rows in %s: %v", tableName, err)
		return -1
	}
	return count
}

// TestEntityFactory provides helper methods to create test entities
type TestEntityFactory struct {
	DataManager *TestDataManager
	Suite       *BaseTestSuite
}

// CreateUser creates a test user
func (f *TestEntityFactory) CreateUser(username, email string) int64 {
	userID, err := f.DataManager.CreateTestUser(username, email)
	f.Suite.Require().NoError(err, "Failed to create test user")
	return userID
}

// CreateBusiness creates a test business
func (f *TestEntityFactory) CreateBusiness(name string, ownerID int64) int64 {
	businessID, err := f.DataManager.CreateTestBusiness(name, fmt.Sprintf("%d", ownerID))
	f.Suite.Require().NoError(err, "Failed to create test business")
	return businessID
}

// CreateCategory creates a test category
func (f *TestEntityFactory) CreateCategory(name string, parentID *int64) int64 {
	categoryID, err := f.DataManager.CreateTestCategory(name, parentID)
	f.Suite.Require().NoError(err, "Failed to create test category")
	return categoryID
}

// CreateProduct creates a test product
func (f *TestEntityFactory) CreateProduct(name string, businessID, categoryID int64, price float64) int64 {
	productID, err := f.DataManager.CreateTestProduct(name, businessID, categoryID, price)
	f.Suite.Require().NoError(err, "Failed to create test product")
	return productID
}

// DatabaseAssertion provides database assertion helpers
type DatabaseAssertion struct {
	DataManager *TestDataManager
	Suite       *BaseTestSuite
}

// AssertTableCount verifies the number of rows in a table
func (a *DatabaseAssertion) AssertTableCount(tableName string, expectedCount int) {
	var actualCount int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := a.DataManager.DB.QueryRow(query).Scan(&actualCount)
	a.Suite.Require().NoError(err, "Failed to count rows in table %s", tableName)
	a.Suite.Assert().Equal(expectedCount, actualCount, "Table %s should have %d rows, but has %d", tableName, expectedCount, actualCount)
}

// AssertRecordExists verifies a record exists with the given conditions
func (a *DatabaseAssertion) AssertRecordExists(tableName string, conditions map[string]interface{}) {
	whereClause, values := buildWhereClause(conditions)
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", tableName, whereClause)

	var count int
	err := a.DataManager.DB.QueryRow(query, values...).Scan(&count)
	a.Suite.Require().NoError(err, "Failed to check record existence in table %s", tableName)
	a.Suite.Assert().Greater(count, 0, "Record should exist in table %s with conditions %v", tableName, conditions)
}

// AssertRecordNotExists verifies a record does not exist with the given conditions
func (a *DatabaseAssertion) AssertRecordNotExists(tableName string, conditions map[string]interface{}) {
	whereClause, values := buildWhereClause(conditions)
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", tableName, whereClause)

	var count int
	err := a.DataManager.DB.QueryRow(query, values...).Scan(&count)
	a.Suite.Require().NoError(err, "Failed to check record non-existence in table %s", tableName)
	a.Suite.Assert().Equal(0, count, "Record should not exist in table %s with conditions %v", tableName, conditions)
}

// buildWhereClause builds a WHERE clause from conditions map
func buildWhereClause(conditions map[string]interface{}) (string, []interface{}) {
	if len(conditions) == 0 {
		return "1=1", []interface{}{}
	}

	var clauses []string
	var values []interface{}
	i := 1

	for column, value := range conditions {
		clauses = append(clauses, fmt.Sprintf("%s = $%d", column, i))
		values = append(values, value)
		i++
	}

	return joinStrings(clauses, " AND "), values
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}