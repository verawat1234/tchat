package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"tchat.dev/shared/config"
	"tchat.dev/shared/database"
)

// TestHelper provides utilities for testing
type TestHelper struct {
	Infrastructure *config.Infrastructure
	Containers     *TestContainers
	Context        context.Context
}

// TestContainers manages test containers
type TestContainers struct {
	Postgres  testcontainers.Container
	Redis     testcontainers.Container
	Scylla    testcontainers.Container
	Kafka     testcontainers.Container
	Zookeeper testcontainers.Container
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	ctx := context.Background()

	// Skip integration tests if running unit tests only
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	helper := &TestHelper{
		Context:    ctx,
		Containers: &TestContainers{},
	}

	// Setup test containers
	helper.setupContainers(t)

	// Setup infrastructure
	helper.setupInfrastructure(t)

	return helper
}

// setupContainers starts test containers
func (h *TestHelper) setupContainers(t *testing.T) {
	var err error

	// Start PostgreSQL container
	h.Containers.Postgres, err = testcontainers.GenericContainer(h.Context, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "tchat_test",
				"POSTGRES_USER":     "test_user",
				"POSTGRES_PASSWORD": "test_pass",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	// Start Redis container
	h.Containers.Redis, err = testcontainers.GenericContainer(h.Context, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:7-alpine",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForListeningPort("6379/tcp").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	// Start ScyllaDB container
	h.Containers.Scylla, err = testcontainers.GenericContainer(h.Context, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "scylladb/scylla:5.2",
			ExposedPorts: []string{"9042/tcp"},
			Cmd:          []string{"--seeds", "127.0.0.1", "--smp", "1"},
			WaitingFor:   wait.ForListeningPort("9042/tcp").WithStartupTimeout(120 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	// Start Zookeeper container (needed for Kafka)
	h.Containers.Zookeeper, err = testcontainers.GenericContainer(h.Context, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "confluentinc/cp-zookeeper:7.4.0",
			ExposedPorts: []string{"2181/tcp"},
			Env: map[string]string{
				"ZOOKEEPER_CLIENT_PORT": "2181",
				"ZOOKEEPER_TICK_TIME":   "2000",
			},
			WaitingFor: wait.ForListeningPort("2181/tcp").WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	// Get Zookeeper endpoint
	zookeeperEndpoint, err := h.Containers.Zookeeper.Endpoint(h.Context, "")
	require.NoError(t, err)

	// Start Kafka container
	h.Containers.Kafka, err = testcontainers.GenericContainer(h.Context, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "confluentinc/cp-kafka:7.4.0",
			ExposedPorts: []string{"9092/tcp"},
			Env: map[string]string{
				"KAFKA_BROKER_ID":                 "1",
				"KAFKA_ZOOKEEPER_CONNECT":         zookeeperEndpoint,
				"KAFKA_ADVERTISED_LISTENERS":      "PLAINTEXT://localhost:9092",
				"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR": "1",
				"KAFKA_TRANSACTION_STATE_LOG_MIN_ISR":     "1",
				"KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR": "1",
			},
			WaitingFor: wait.ForListeningPort("9092/tcp").WithStartupTimeout(120 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	log.Println("All test containers started successfully")
}

// setupInfrastructure initializes test infrastructure
func (h *TestHelper) setupInfrastructure(t *testing.T) {
	// Get container endpoints
	postgresHost, err := h.Containers.Postgres.Host(h.Context)
	require.NoError(t, err)
	postgresPort, err := h.Containers.Postgres.MappedPort(h.Context, "5432")
	require.NoError(t, err)

	redisHost, err := h.Containers.Redis.Host(h.Context)
	require.NoError(t, err)
	redisPort, err := h.Containers.Redis.MappedPort(h.Context, "6379")
	require.NoError(t, err)

	scyllaHost, err := h.Containers.Scylla.Host(h.Context)
	require.NoError(t, err)
	scyllaPort, err := h.Containers.Scylla.MappedPort(h.Context, "9042")
	require.NoError(t, err)

	kafkaHost, err := h.Containers.Kafka.Host(h.Context)
	require.NoError(t, err)
	kafkaPort, err := h.Containers.Kafka.MappedPort(h.Context, "9092")
	require.NoError(t, err)

	// Create test configuration
	testConfig := &config.InfrastructureConfig{
		ServiceName: "test-service",
		Environment: "test",
		Database: &config.DatabaseConfig{
			Postgres: &database.PostgresConfig{
				Host:     postgresHost,
				Port:     postgresPort.Int(),
				User:     "test_user",
				Password: "test_pass",
				Database: "tchat_test",
				SSLMode:  "disable",
			},
			Scylla: &database.ScyllaConfig{
				Hosts:    []string{fmt.Sprintf("%s:%d", scyllaHost, scyllaPort.Int())},
				Keyspace: "tchat_test",
			},
		},
		Cache: &config.CacheConfig{
			Redis: &database.RedisConfig{
				Host: redisHost,
				Port: redisPort.Int(),
			},
		},
		Messaging: &config.MessagingConfig{
			Kafka: &messaging.KafkaConfig{
				Brokers:     []string{fmt.Sprintf("%s:%d", kafkaHost, kafkaPort.Int())},
				GroupID:     "tchat-test",
				TopicPrefix: "test",
			},
		},
	}

	// Initialize infrastructure
	h.Infrastructure = config.NewInfrastructure(testConfig)
	err = h.Infrastructure.Initialize(h.Context)
	require.NoError(t, err)

	// Run database migrations
	h.setupDatabases(t)

	log.Println("Test infrastructure initialized successfully")
}

// setupDatabases runs test database setup
func (h *TestHelper) setupDatabases(t *testing.T) {
	// Setup PostgreSQL test schema
	h.setupPostgreSQL(t)

	// Setup ScyllaDB test keyspace
	h.setupScyllaDB(t)
}

// setupPostgreSQL creates test tables
func (h *TestHelper) setupPostgreSQL(t *testing.T) {
	db := h.Infrastructure.PostgresDB

	// Create test tables
	testTables := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			phone_number VARCHAR(20) UNIQUE NOT NULL,
			country_code VARCHAR(5) NOT NULL DEFAULT '+66',
			username VARCHAR(50) UNIQUE,
			display_name VARCHAR(100),
			status VARCHAR(20) DEFAULT 'active',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,

		// Test dialogs table
		`CREATE TABLE IF NOT EXISTS dialogs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			type VARCHAR(20) NOT NULL,
			title VARCHAR(255),
			creator_id UUID NOT NULL,
			participant_count INT DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,

		// Test messages table
		`CREATE TABLE IF NOT EXISTS messages (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			dialog_id UUID NOT NULL,
			sender_id UUID NOT NULL,
			type VARCHAR(20) NOT NULL DEFAULT 'text',
			content TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,

		// Test wallets table
		`CREATE TABLE IF NOT EXISTS wallets (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL,
			currency VARCHAR(10) NOT NULL,
			balance DECIMAL(20,2) DEFAULT 0.00,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,

		// Test transactions table
		`CREATE TABLE IF NOT EXISTS transactions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			from_wallet_id UUID,
			to_wallet_id UUID,
			amount DECIMAL(20,2) NOT NULL,
			currency VARCHAR(10) NOT NULL,
			status VARCHAR(20) DEFAULT 'pending',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, table := range testTables {
		err := db.ExecInTransaction(h.Context, table)
		require.NoError(t, err)
	}

	log.Println("PostgreSQL test tables created")
}

// setupScyllaDB creates test keyspace and tables
func (h *TestHelper) setupScyllaDB(t *testing.T) {
	scylla := h.Infrastructure.ScyllaDB

	// Wait for ScyllaDB to be ready
	time.Sleep(10 * time.Second)

	// Create test tables
	testTables := []string{
		// Message timeline table
		`CREATE TABLE IF NOT EXISTS message_timeline (
			dialog_id UUID,
			bucket_time timestamp,
			message_id UUID,
			sender_id UUID,
			content text,
			created_at timestamp,
			PRIMARY KEY (dialog_id, bucket_time, message_id)
		) WITH CLUSTERING ORDER BY (bucket_time DESC, message_id DESC)`,

		// User message timeline
		`CREATE TABLE IF NOT EXISTS user_message_timeline (
			user_id UUID,
			bucket_time timestamp,
			message_id UUID,
			dialog_id UUID,
			content text,
			created_at timestamp,
			PRIMARY KEY (user_id, bucket_time, message_id)
		) WITH CLUSTERING ORDER BY (bucket_time DESC, message_id DESC)`,
	}

	for _, table := range testTables {
		err := scylla.ExecuteQuery(table)
		if err != nil {
			log.Printf("Warning: Failed to create ScyllaDB table: %v", err)
		}
	}

	log.Println("ScyllaDB test tables created")
}

// Cleanup cleans up test resources
func (h *TestHelper) Cleanup() {
	if h.Infrastructure != nil {
		h.Infrastructure.Shutdown(h.Context)
	}

	// Stop containers
	if h.Containers.Kafka != nil {
		h.Containers.Kafka.Terminate(h.Context)
	}
	if h.Containers.Zookeeper != nil {
		h.Containers.Zookeeper.Terminate(h.Context)
	}
	if h.Containers.Scylla != nil {
		h.Containers.Scylla.Terminate(h.Context)
	}
	if h.Containers.Redis != nil {
		h.Containers.Redis.Terminate(h.Context)
	}
	if h.Containers.Postgres != nil {
		h.Containers.Postgres.Terminate(h.Context)
	}

	log.Println("Test cleanup completed")
}

// CleanupTables cleans up test data
func (h *TestHelper) CleanupTables(t *testing.T) {
	db := h.Infrastructure.PostgresDB

	tables := []string{"transactions", "wallets", "messages", "dialogs", "users"}
	for _, table := range tables {
		err := db.ExecInTransaction(h.Context, fmt.Sprintf("DELETE FROM %s", table))
		require.NoError(t, err)
	}
}

// CreateTestUser creates a test user
func (h *TestHelper) CreateTestUser(t *testing.T, phone, countryCode string) string {
	db := h.Infrastructure.PostgresDB

	var userID string
	err := db.GetInTransaction(h.Context, &userID,
		`INSERT INTO users (phone_number, country_code) VALUES ($1, $2) RETURNING id`,
		phone, countryCode)
	require.NoError(t, err)

	return userID
}

// CreateTestDialog creates a test dialog
func (h *TestHelper) CreateTestDialog(t *testing.T, creatorID, title string) string {
	db := h.Infrastructure.PostgresDB

	var dialogID string
	err := db.GetInTransaction(h.Context, &dialogID,
		`INSERT INTO dialogs (creator_id, title, type) VALUES ($1, $2, $3) RETURNING id`,
		creatorID, title, "direct")
	require.NoError(t, err)

	return dialogID
}

// CreateTestWallet creates a test wallet
func (h *TestHelper) CreateTestWallet(t *testing.T, userID, currency string) string {
	db := h.Infrastructure.PostgresDB

	var walletID string
	err := db.GetInTransaction(h.Context, &walletID,
		`INSERT INTO wallets (user_id, currency) VALUES ($1, $2) RETURNING id`,
		userID, currency)
	require.NoError(t, err)

	return walletID
}

// WaitForKafkaMessage waits for a Kafka message on a topic
func (h *TestHelper) WaitForKafkaMessage(t *testing.T, topic string, timeout time.Duration) []byte {
	kafkaHost, err := h.Containers.Kafka.Host(h.Context)
	require.NoError(t, err)
	kafkaPort, err := h.Containers.Kafka.MappedPort(h.Context, "9092")
	require.NoError(t, err)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{fmt.Sprintf("%s:%d", kafkaHost, kafkaPort.Int())},
		Topic:   fmt.Sprintf("test.%s", topic),
		GroupID: "test-consumer",
	})
	defer reader.Close()

	ctx, cancel := context.WithTimeout(h.Context, timeout)
	defer cancel()

	message, err := reader.ReadMessage(ctx)
	require.NoError(t, err)

	return message.Value
}

// AssertEventuallyConsistent waits for eventual consistency
func (h *TestHelper) AssertEventuallyConsistent(t *testing.T, assertion func() bool, timeout time.Duration, interval time.Duration) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if assertion() {
			return // Success
		}
		time.Sleep(interval)
	}

	t.Fatalf("Assertion failed within timeout of %v", timeout)
}

// GenerateTestData generates test data for load testing
type TestDataGenerator struct {
	helper *TestHelper
}

// NewTestDataGenerator creates a new test data generator
func NewTestDataGenerator(helper *TestHelper) *TestDataGenerator {
	return &TestDataGenerator{helper: helper}
}

// GenerateUsers generates test users
func (g *TestDataGenerator) GenerateUsers(t *testing.T, count int) []string {
	var userIDs []string

	for i := 0; i < count; i++ {
		phone := fmt.Sprintf("+66%09d", 800000000+i)
		userID := g.helper.CreateTestUser(t, phone, "+66")
		userIDs = append(userIDs, userID)
	}

	return userIDs
}

// GenerateDialogs generates test dialogs
func (g *TestDataGenerator) GenerateDialogs(t *testing.T, userIDs []string, count int) []string {
	var dialogIDs []string

	for i := 0; i < count; i++ {
		creatorID := userIDs[i%len(userIDs)]
		title := fmt.Sprintf("Test Dialog %d", i+1)
		dialogID := g.helper.CreateTestDialog(t, creatorID, title)
		dialogIDs = append(dialogIDs, dialogID)
	}

	return dialogIDs
}

// MockExternalServices provides mocks for external services
type MockExternalServices struct {
	SMSService     *MockSMSService
	PaymentGateway *MockPaymentGateway
	StorageService *MockStorageService
}

// MockSMSService mocks SMS service
type MockSMSService struct {
	SentMessages []MockSMSMessage
}

// MockSMSMessage represents a mock SMS message
type MockSMSMessage struct {
	To      string
	Message string
	SentAt  time.Time
}

// SendSMS mocks sending SMS
func (m *MockSMSService) SendSMS(to, message string) error {
	m.SentMessages = append(m.SentMessages, MockSMSMessage{
		To:      to,
		Message: message,
		SentAt:  time.Now(),
	})
	return nil
}

// MockPaymentGateway mocks payment gateway
type MockPaymentGateway struct {
	Transactions []MockTransaction
}

// MockTransaction represents a mock transaction
type MockTransaction struct {
	ID       string
	Amount   float64
	Currency string
	Status   string
	SentAt   time.Time
}

// ProcessPayment mocks payment processing
func (m *MockPaymentGateway) ProcessPayment(amount float64, currency string) (*MockTransaction, error) {
	transaction := &MockTransaction{
		ID:       fmt.Sprintf("txn_%d", time.Now().UnixNano()),
		Amount:   amount,
		Currency: currency,
		Status:   "completed",
		SentAt:   time.Now(),
	}

	m.Transactions = append(m.Transactions, *transaction)
	return transaction, nil
}

// MockStorageService mocks storage service
type MockStorageService struct {
	StoredFiles []MockStoredFile
}

// MockStoredFile represents a mock stored file
type MockStoredFile struct {
	ID       string
	Name     string
	Size     int64
	URL      string
	SentAt   time.Time
}

// UploadFile mocks file upload
func (m *MockStorageService) UploadFile(name string, size int64) (*MockStoredFile, error) {
	file := &MockStoredFile{
		ID:     fmt.Sprintf("file_%d", time.Now().UnixNano()),
		Name:   name,
		Size:   size,
		URL:    fmt.Sprintf("https://cdn.example.com/%s", name),
		SentAt: time.Now(),
	}

	m.StoredFiles = append(m.StoredFiles, *file)
	return file, nil
}

// NewMockExternalServices creates mock external services
func NewMockExternalServices() *MockExternalServices {
	return &MockExternalServices{
		SMSService:     &MockSMSService{},
		PaymentGateway: &MockPaymentGateway{},
		StorageService: &MockStorageService{},
	}
}

// PerformanceTestHelper provides utilities for performance testing
type PerformanceTestHelper struct {
	StartTime    time.Time
	EndTime      time.Time
	RequestCount int64
	ErrorCount   int64
	ResponseTimes []time.Duration
}

// NewPerformanceTestHelper creates a new performance test helper
func NewPerformanceTestHelper() *PerformanceTestHelper {
	return &PerformanceTestHelper{
		StartTime:     time.Now(),
		ResponseTimes: make([]time.Duration, 0),
	}
}

// RecordRequest records a request
func (p *PerformanceTestHelper) RecordRequest(duration time.Duration, isError bool) {
	p.RequestCount++
	p.ResponseTimes = append(p.ResponseTimes, duration)

	if isError {
		p.ErrorCount++
	}
}

// GetStats returns performance statistics
func (p *PerformanceTestHelper) GetStats() PerformanceStats {
	p.EndTime = time.Now()
	totalDuration := p.EndTime.Sub(p.StartTime)

	var totalResponseTime time.Duration
	var minResponseTime time.Duration = time.Hour
	var maxResponseTime time.Duration

	for _, rt := range p.ResponseTimes {
		totalResponseTime += rt
		if rt < minResponseTime {
			minResponseTime = rt
		}
		if rt > maxResponseTime {
			maxResponseTime = rt
		}
	}

	avgResponseTime := time.Duration(0)
	if len(p.ResponseTimes) > 0 {
		avgResponseTime = totalResponseTime / time.Duration(len(p.ResponseTimes))
	}

	rps := float64(p.RequestCount) / totalDuration.Seconds()
	errorRate := float64(p.ErrorCount) / float64(p.RequestCount) * 100

	return PerformanceStats{
		TotalRequests:     p.RequestCount,
		ErrorCount:        p.ErrorCount,
		ErrorRate:         errorRate,
		RequestsPerSecond: rps,
		AvgResponseTime:   avgResponseTime,
		MinResponseTime:   minResponseTime,
		MaxResponseTime:   maxResponseTime,
		TotalDuration:     totalDuration,
	}
}

// PerformanceStats contains performance test results
type PerformanceStats struct {
	TotalRequests     int64         `json:"total_requests"`
	ErrorCount        int64         `json:"error_count"`
	ErrorRate         float64       `json:"error_rate"`
	RequestsPerSecond float64       `json:"requests_per_second"`
	AvgResponseTime   time.Duration `json:"avg_response_time"`
	MinResponseTime   time.Duration `json:"min_response_time"`
	MaxResponseTime   time.Duration `json:"max_response_time"`
	TotalDuration     time.Duration `json:"total_duration"`
}

// TestEnvironment manages test environment variables
func SetupTestEnvironment() {
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("DATABASE_URL", "postgres://test_user:test_pass@localhost/tchat_test?sslmode=disable")
	os.Setenv("REDIS_URL", "redis://localhost:6379/0")
	os.Setenv("KAFKA_BROKERS", "localhost:9092")
}

// CleanupTestEnvironment cleans up test environment
func CleanupTestEnvironment() {
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("REDIS_URL")
	os.Unsetenv("KAFKA_BROKERS")
}