package fixtures

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// TestDataManager handles database setup, seeding, and cleanup for integration tests
type TestDataManager struct {
	DB           *sql.DB
	dbName       string
	migrator     *migrate.Migrate
	fixturesPath string
}

// TestFixture represents a database fixture for testing
type TestFixture struct {
	Table string                 `json:"table"`
	Data  []map[string]interface{} `json:"data"`
}

// TestDataConfig holds configuration for test data setup
type TestDataConfig struct {
	DatabaseURL   string
	MigrationsDir string
	FixturesDir   string
	TestDBPrefix  string
}

// NewTestDataManager creates a new test data manager instance
func NewTestDataManager(config TestDataConfig) (*TestDataManager, error) {
	// Generate unique test database name
	dbName := fmt.Sprintf("%s_test_%d", config.TestDBPrefix, time.Now().UnixNano())

	// Create test database
	masterDB, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master database: %w", err)
	}
	defer masterDB.Close()

	_, err = masterDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		return nil, fmt.Errorf("failed to create test database: %w", err)
	}

	// Connect to test database
	testDBURL := config.DatabaseURL + "/" + dbName
	testDB, err := sql.Open("postgres", testDBURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Setup migrator
	driver, err := postgres.WithInstance(testDB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migration driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", config.MigrationsDir),
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return &TestDataManager{
		DB:           testDB,
		dbName:       dbName,
		migrator:     migrator,
		fixturesPath: config.FixturesDir,
	}, nil
}

// SetupDatabase runs migrations and sets up the database schema
func (tdm *TestDataManager) SetupDatabase() error {
	log.Printf("Setting up test database: %s", tdm.dbName)

	// Run migrations
	err := tdm.migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Printf("Database setup completed for: %s", tdm.dbName)
	return nil
}

// LoadFixtures loads test data from JSON fixtures
func (tdm *TestDataManager) LoadFixtures(fixtureNames ...string) error {
	ctx := context.Background()

	for _, fixtureName := range fixtureNames {
		err := tdm.loadFixture(ctx, fixtureName)
		if err != nil {
			return fmt.Errorf("failed to load fixture %s: %w", fixtureName, err)
		}
	}

	return nil
}

// loadFixture loads a single fixture file
func (tdm *TestDataManager) loadFixture(ctx context.Context, fixtureName string) error {
	fixturePath := filepath.Join(tdm.fixturesPath, fixtureName+".json")

	data, err := os.ReadFile(fixturePath)
	if err != nil {
		return fmt.Errorf("failed to read fixture file: %w", err)
	}

	var fixture TestFixture
	err = json.Unmarshal(data, &fixture)
	if err != nil {
		return fmt.Errorf("failed to unmarshal fixture: %w", err)
	}

	// Begin transaction for atomic fixture loading
	tx, err := tdm.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert fixture data
	for _, record := range fixture.Data {
		err = tdm.insertRecord(tx, fixture.Table, record)
		if err != nil {
			return fmt.Errorf("failed to insert record into %s: %w", fixture.Table, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit fixture transaction: %w", err)
	}

	log.Printf("Loaded fixture: %s (%d records)", fixtureName, len(fixture.Data))
	return nil
}

// insertRecord inserts a single record into the specified table
func (tdm *TestDataManager) insertRecord(tx *sql.Tx, table string, record map[string]interface{}) error {
	if len(record) == 0 {
		return nil
	}

	// Build INSERT query
	columns := make([]string, 0, len(record))
	placeholders := make([]string, 0, len(record))
	values := make([]interface{}, 0, len(record))

	i := 1
	for column, value := range record {
		columns = append(columns, column)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		values = append(values, value)
		i++
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		joinStrings(columns, ", "),
		joinStrings(placeholders, ", "),
	)

	_, err := tx.Exec(query, values...)
	return err
}

// CleanupTable truncates a specific table
func (tdm *TestDataManager) CleanupTable(tableName string) error {
	_, err := tdm.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName))
	if err != nil {
		return fmt.Errorf("failed to truncate table %s: %w", tableName, err)
	}
	return nil
}

// CleanupAllTables truncates all tables in the database
func (tdm *TestDataManager) CleanupAllTables() error {
	// Get all table names
	query := `
		SELECT tablename
		FROM pg_tables
		WHERE schemaname = 'public'
		AND tablename NOT LIKE 'schema_migrations%'
	`

	rows, err := tdm.DB.Query(query)
	if err != nil {
		return fmt.Errorf("failed to get table names: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, table)
	}

	// Truncate all tables
	if len(tables) > 0 {
		query = fmt.Sprintf("TRUNCATE TABLE %s CASCADE", joinStrings(tables, ", "))
		_, err = tdm.DB.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to truncate tables: %w", err)
		}
	}

	log.Printf("Cleaned up %d tables", len(tables))
	return nil
}

// ResetSequences resets all auto-increment sequences
func (tdm *TestDataManager) ResetSequences() error {
	query := `
		SELECT 'SELECT SETVAL(' || quote_literal(quote_ident(PGT.schemaname) || '.' || quote_ident(S.relname)) ||
		       ', COALESCE(MAX(' || quote_ident(C.attname) || '), 1)) FROM ' ||
		       quote_ident(PGT.schemaname) || '.' || quote_ident(T.relname) || ';'
		FROM pg_class AS S,
		     pg_depend AS D,
		     pg_class AS T,
		     pg_attribute AS C,
		     pg_tables AS PGT
		WHERE S.relkind = 'S'
		  AND S.oid = D.objid
		  AND D.refobjid = T.oid
		  AND D.refobjid = C.attrelid
		  AND D.refobjsubid = C.attnum
		  AND T.relname = PGT.tablename
		  AND PGT.schemaname = 'public'
	`

	rows, err := tdm.DB.Query(query)
	if err != nil {
		return fmt.Errorf("failed to get sequence reset queries: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var resetQuery string
		if err := rows.Scan(&resetQuery); err != nil {
			return fmt.Errorf("failed to scan reset query: %w", err)
		}

		_, err = tdm.DB.Exec(resetQuery)
		if err != nil {
			log.Printf("Warning: failed to reset sequence: %v", err)
		}
	}

	return nil
}

// TeardownDatabase drops the test database
func (tdm *TestDataManager) TeardownDatabase() error {
	log.Printf("Tearing down test database: %s", tdm.dbName)

	// Close database connection
	if tdm.DB != nil {
		tdm.DB.Close()
	}

	// Connect to master database to drop test database
	masterDBURL := os.Getenv("DATABASE_URL")
	if masterDBURL == "" {
		masterDBURL = "postgres://localhost/postgres?sslmode=disable"
	}

	masterDB, err := sql.Open("postgres", masterDBURL)
	if err != nil {
		return fmt.Errorf("failed to connect to master database: %w", err)
	}
	defer masterDB.Close()

	// Terminate existing connections to test database
	_, err = masterDB.Exec(fmt.Sprintf(`
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = '%s' AND pid <> pg_backend_pid()
	`, tdm.dbName))
	if err != nil {
		log.Printf("Warning: failed to terminate connections: %v", err)
	}

	// Drop test database
	_, err = masterDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", tdm.dbName))
	if err != nil {
		return fmt.Errorf("failed to drop test database: %w", err)
	}

	log.Printf("Test database dropped: %s", tdm.dbName)
	return nil
}

// GetTestTransaction returns a test transaction for isolated testing
func (tdm *TestDataManager) GetTestTransaction() (*sql.Tx, error) {
	return tdm.DB.Begin()
}

// CreateTestUser creates a test user and returns the user ID
func (tdm *TestDataManager) CreateTestUser(username, email string) (int64, error) {
	query := `
		INSERT INTO users (username, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, 'test_hash', NOW(), NOW())
		RETURNING id
	`

	var userID int64
	err := tdm.DB.QueryRow(query, username, email).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to create test user: %w", err)
	}

	return userID, nil
}

// CreateTestBusiness creates a test business and returns the business ID
func (tdm *TestDataManager) CreateTestBusiness(name, ownerID string) (int64, error) {
	query := `
		INSERT INTO businesses (name, owner_id, status, created_at, updated_at)
		VALUES ($1, $2, 'active', NOW(), NOW())
		RETURNING id
	`

	var businessID int64
	err := tdm.DB.QueryRow(query, name, ownerID).Scan(&businessID)
	if err != nil {
		return 0, fmt.Errorf("failed to create test business: %w", err)
	}

	return businessID, nil
}

// CreateTestCategory creates a test category and returns the category ID
func (tdm *TestDataManager) CreateTestCategory(name string, parentID *int64) (int64, error) {
	query := `
		INSERT INTO categories (name, parent_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id
	`

	var categoryID int64
	err := tdm.DB.QueryRow(query, name, parentID).Scan(&categoryID)
	if err != nil {
		return 0, fmt.Errorf("failed to create test category: %w", err)
	}

	return categoryID, nil
}

// CreateTestProduct creates a test product and returns the product ID
func (tdm *TestDataManager) CreateTestProduct(name string, businessID, categoryID int64, price float64) (int64, error) {
	query := `
		INSERT INTO products (name, business_id, category_id, price, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 'active', NOW(), NOW())
		RETURNING id
	`

	var productID int64
	err := tdm.DB.QueryRow(query, name, businessID, categoryID, price).Scan(&productID)
	if err != nil {
		return 0, fmt.Errorf("failed to create test product: %w", err)
	}

	return productID, nil
}

// Helper function to join strings
func joinStrings(strings []string, separator string) string {
	if len(strings) == 0 {
		return ""
	}

	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += separator + strings[i]
	}

	return result
}