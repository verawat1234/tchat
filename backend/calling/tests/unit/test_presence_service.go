package unit

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"tchat.dev/calling/config"
)

// PresenceServiceTestSuite defines the test suite for PresenceService
type PresenceServiceTestSuite struct {
	suite.Suite
	presenceService *config.UserPresenceService
	redisClient     redis.Cmdable
	redisMock       redismock.ClientMock
}

// SetupTest sets up test dependencies
func (suite *PresenceServiceTestSuite) SetupTest() {
	// Create Redis mock
	db, mock := redismock.NewClientMock()
	suite.redisClient = db
	suite.redisMock = mock

	// Create Redis client wrapper
	redisWrapper := &config.RedisClient{
		Client: db.(*redis.Client),
	}

	// Create presence service
	suite.presenceService = config.NewUserPresenceService(redisWrapper)
}

// TearDownTest cleans up after each test
func (suite *PresenceServiceTestSuite) TearDownTest() {
	assert.NoError(suite.T(), suite.redisMock.ExpectationsWereMet())
}

// TestSetUserOnline tests setting user online status
func (suite *PresenceServiceTestSuite) TestSetUserOnline() {
	userID := uuid.New().String()
	callID := ""

	// Mock Redis operations
	key := "presence:user:" + userID
	presenceData := map[string]interface{}{
		"status":    "online",
		"in_call":   false,
		"call_id":   callID,
		"last_seen": time.Now().Unix(),
	}

	suite.redisMock.ExpectHMSet(key, presenceData).SetVal(true)
	suite.redisMock.ExpectExpire(key, 5*time.Minute).SetVal(true)

	// Execute test
	err := suite.presenceService.SetUserOnline(userID, callID)

	// Assertions
	assert.NoError(suite.T(), err)
}

// TestSetUserOnlineWithCall tests setting user online with active call
func (suite *PresenceServiceTestSuite) TestSetUserOnlineWithCall() {
	userID := uuid.New().String()
	callID := uuid.New().String()

	// Mock Redis operations
	key := "presence:user:" + userID
	presenceData := map[string]interface{}{
		"status":    "online",
		"in_call":   true,
		"call_id":   callID,
		"last_seen": time.Now().Unix(),
	}

	suite.redisMock.ExpectHMSet(key, presenceData).SetVal(true)
	suite.redisMock.ExpectExpire(key, 5*time.Minute).SetVal(true)

	// Execute test
	err := suite.presenceService.SetUserOnline(userID, callID)

	// Assertions
	assert.NoError(suite.T(), err)
}

// TestSetUserOffline tests setting user offline status
func (suite *PresenceServiceTestSuite) TestSetUserOffline() {
	userID := uuid.New().String()
	key := "presence:user:" + userID

	// Mock Redis operations
	suite.redisMock.ExpectDel(key).SetVal(1)

	// Execute test
	err := suite.presenceService.SetUserOffline(userID)

	// Assertions
	assert.NoError(suite.T(), err)
}

// TestGetUserPresence tests retrieving user presence
func (suite *PresenceServiceTestSuite) TestGetUserPresence() {
	userID := uuid.New().String()
	key := "presence:user:" + userID

	// Expected presence data
	expectedPresence := map[string]string{
		"status":    "online",
		"in_call":   "false",
		"call_id":   "",
		"last_seen": "1640995200",
	}

	// Mock Redis operations
	suite.redisMock.ExpectHGetAll(key).SetVal(expectedPresence)

	// Execute test
	presence, err := suite.presenceService.GetUserPresence(userID)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedPresence, presence)
}

// TestGetUserPresenceOffline tests retrieving offline user presence
func (suite *PresenceServiceTestSuite) TestGetUserPresenceOffline() {
	userID := uuid.New().String()
	key := "presence:user:" + userID

	// Mock Redis operations - empty result indicates offline user
	suite.redisMock.ExpectHGetAll(key).SetVal(map[string]string{})

	// Execute test
	presence, err := suite.presenceService.GetUserPresence(userID)

	// Expected offline presence
	expectedOfflinePresence := map[string]string{
		"status":    "offline",
		"in_call":   "false",
		"call_id":   "",
		"last_seen": "0",
	}

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedOfflinePresence, presence)
}

// TestPresenceExpiration tests presence data expiration
func (suite *PresenceServiceTestSuite) TestPresenceExpiration() {
	userID := uuid.New().String()
	key := "presence:user:" + userID

	// Test data
	callID := ""
	presenceData := map[string]interface{}{
		"status":    "online",
		"in_call":   false,
		"call_id":   callID,
		"last_seen": time.Now().Unix(),
	}

	// Mock Redis operations with TTL verification
	suite.redisMock.ExpectHMSet(key, presenceData).SetVal(true)
	suite.redisMock.ExpectExpire(key, 5*time.Minute).SetVal(true)

	// Execute test
	err := suite.presenceService.SetUserOnline(userID, callID)

	// Assertions
	assert.NoError(suite.T(), err)
}

// TestConcurrentPresenceUpdates tests concurrent presence updates
func (suite *PresenceServiceTestSuite) TestConcurrentPresenceUpdates() {
	userID := uuid.New().String()
	numGoroutines := 10

	// Mock Redis operations for concurrent updates
	for i := 0; i < numGoroutines; i++ {
		key := "presence:user:" + userID
		presenceData := map[string]interface{}{
			"status":    "online",
			"in_call":   false,
			"call_id":   "",
			"last_seen": time.Now().Unix(),
		}

		suite.redisMock.ExpectHMSet(key, presenceData).SetVal(true)
		suite.redisMock.ExpectExpire(key, 5*time.Minute).SetVal(true)
	}

	// Execute concurrent updates
	done := make(chan bool, numGoroutines)
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			err := suite.presenceService.SetUserOnline(userID, "")
			if err != nil {
				errors <- err
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Check for errors
	select {
	case err := <-errors:
		assert.NoError(suite.T(), err)
	default:
		// No errors, which is expected
	}
}

// TestPresenceStateTransitions tests valid presence state transitions
func (suite *PresenceServiceTestSuite) TestPresenceStateTransitions() {
	userID := uuid.New().String()

	// Test sequence: offline -> online -> in_call -> online -> offline
	transitions := []struct {
		description string
		callID      string
		expectedInCall bool
	}{
		{"Set user online", "", false},
		{"Set user in call", uuid.New().String(), true},
		{"Set user online without call", "", false},
	}

	key := "presence:user:" + userID

	for _, transition := range transitions {
		suite.T().Run(transition.description, func(t *testing.T) {
			// Mock Redis operations
			presenceData := map[string]interface{}{
				"status":    "online",
				"in_call":   transition.expectedInCall,
				"call_id":   transition.callID,
				"last_seen": time.Now().Unix(),
			}

			suite.redisMock.ExpectHMSet(key, presenceData).SetVal(true)
			suite.redisMock.ExpectExpire(key, 5*time.Minute).SetVal(true)

			// Execute test
			err := suite.presenceService.SetUserOnline(userID, transition.callID)
			assert.NoError(t, err)
		})
	}

	// Final offline transition
	suite.redisMock.ExpectDel(key).SetVal(1)
	err := suite.presenceService.SetUserOffline(userID)
	assert.NoError(suite.T(), err)
}

// TestRedisConnectionFailure tests handling of Redis connection failures
func (suite *PresenceServiceTestSuite) TestRedisConnectionFailure() {
	userID := uuid.New().String()
	key := "presence:user:" + userID

	// Mock Redis connection failure
	suite.redisMock.ExpectHMSet(key, map[string]interface{}{
		"status":    "online",
		"in_call":   false,
		"call_id":   "",
		"last_seen": time.Now().Unix(),
	}).SetErr(redis.ErrClosed)

	// Execute test
	err := suite.presenceService.SetUserOnline(userID, "")

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), redis.ErrClosed, err)
}

// TestPresenceDataValidation tests presence data validation
func (suite *PresenceServiceTestSuite) TestPresenceDataValidation() {
	tests := []struct {
		name     string
		userID   string
		callID   string
		expectError bool
	}{
		{
			name:     "Valid user ID",
			userID:   uuid.New().String(),
			callID:   "",
			expectError: false,
		},
		{
			name:     "Valid user ID with call",
			userID:   uuid.New().String(),
			callID:   uuid.New().String(),
			expectError: false,
		},
		{
			name:     "Empty user ID",
			userID:   "",
			callID:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			// Validate user ID
			if tt.userID == "" {
				// Skip Redis mock for invalid input
				err := suite.presenceService.SetUserOnline(tt.userID, tt.callID)
				if tt.expectError {
					assert.Error(t, err)
				}
				return
			}

			// Mock Redis operations for valid input
			key := "presence:user:" + tt.userID
			presenceData := map[string]interface{}{
				"status":    "online",
				"in_call":   tt.callID != "",
				"call_id":   tt.callID,
				"last_seen": time.Now().Unix(),
			}

			suite.redisMock.ExpectHMSet(key, presenceData).SetVal(true)
			suite.redisMock.ExpectExpire(key, 5*time.Minute).SetVal(true)

			err := suite.presenceService.SetUserOnline(tt.userID, tt.callID)
			assert.NoError(t, err)
		})
	}
}

// TestPresenceServiceTestSuite runs the test suite
func TestPresenceServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PresenceServiceTestSuite))
}

// Benchmark tests

func BenchmarkSetUserOnline(b *testing.B) {
	// Setup mock Redis
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisWrapper := &config.RedisClient{
		Client: db.(*redis.Client),
	}
	presenceService := config.NewUserPresenceService(redisWrapper)

	userID := uuid.New().String()
	key := "presence:user:" + userID

	// Setup expectations for benchmark iterations
	for i := 0; i < b.N; i++ {
		presenceData := map[string]interface{}{
			"status":    "online",
			"in_call":   false,
			"call_id":   "",
			"last_seen": time.Now().Unix(),
		}
		mock.ExpectHMSet(key, presenceData).SetVal(true)
		mock.ExpectExpire(key, 5*time.Minute).SetVal(true)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		presenceService.SetUserOnline(userID, "")
	}
}

func BenchmarkGetUserPresence(b *testing.B) {
	// Setup mock Redis
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisWrapper := &config.RedisClient{
		Client: db.(*redis.Client),
	}
	presenceService := config.NewUserPresenceService(redisWrapper)

	userID := uuid.New().String()
	key := "presence:user:" + userID

	expectedPresence := map[string]string{
		"status":    "online",
		"in_call":   "false",
		"call_id":   "",
		"last_seen": "1640995200",
	}

	// Setup expectations for benchmark iterations
	for i := 0; i < b.N; i++ {
		mock.ExpectHGetAll(key).SetVal(expectedPresence)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		presenceService.GetUserPresence(userID)
	}
}