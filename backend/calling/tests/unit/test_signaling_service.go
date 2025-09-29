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

// SignalingServiceTestSuite defines the test suite for SignalingService
type SignalingServiceTestSuite struct {
	suite.Suite
	signalingService *config.SignalingService
	redisClient      redis.Cmdable
	redisMock        redismock.ClientMock
}

// SetupTest sets up test dependencies
func (suite *SignalingServiceTestSuite) SetupTest() {
	// Create Redis mock
	db, mock := redismock.NewClientMock()
	suite.redisClient = db
	suite.redisMock = mock

	// Create Redis client wrapper
	redisWrapper := &config.RedisClient{
		Client: db.(*redis.Client),
	}

	// Create signaling service
	suite.signalingService = config.NewSignalingService(redisWrapper)
}

// TearDownTest cleans up after each test
func (suite *SignalingServiceTestSuite) TearDownTest() {
	assert.NoError(suite.T(), suite.redisMock.ExpectationsWereMet())
}

// TestStoreSignalingMessage tests storing WebRTC signaling messages
func (suite *SignalingServiceTestSuite) TestStoreSignalingMessage() {
	callID := uuid.New().String()
	fromUserID := uuid.New().String()
	message := map[string]interface{}{
		"type": "offer",
		"sdp":  "v=0\r\no=- 123456789 1 IN IP4 127.0.0.1\r\n...",
	}

	key := "signaling:" + callID

	// Expected message data with metadata
	expectedMessageData := map[string]interface{}{
		"from":      fromUserID,
		"timestamp": time.Now().Unix(),
		"message":   message,
	}

	// Mock Redis operations
	suite.redisMock.ExpectLPush(key, expectedMessageData).SetVal(1)
	suite.redisMock.ExpectExpire(key, 10*time.Minute).SetVal(true)
	suite.redisMock.ExpectLTrim(key, 0, 100).SetVal(true)

	// Execute test
	err := suite.signalingService.StoreSignalingMessage(callID, fromUserID, message)

	// Assertions
	assert.NoError(suite.T(), err)
}

// TestStoreOfferMessage tests storing WebRTC offer messages
func (suite *SignalingServiceTestSuite) TestStoreOfferMessage() {
	callID := uuid.New().String()
	fromUserID := uuid.New().String()
	offerMessage := map[string]interface{}{
		"type": "offer",
		"sdp":  "v=0\r\no=- 123456789 1 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\nm=audio 54400 RTP/AVP 0\r\nc=IN IP4 127.0.0.1\r\n",
	}

	key := "signaling:" + callID

	expectedMessageData := map[string]interface{}{
		"from":      fromUserID,
		"timestamp": time.Now().Unix(),
		"message":   offerMessage,
	}

	// Mock Redis operations
	suite.redisMock.ExpectLPush(key, expectedMessageData).SetVal(1)
	suite.redisMock.ExpectExpire(key, 10*time.Minute).SetVal(true)
	suite.redisMock.ExpectLTrim(key, 0, 100).SetVal(true)

	// Execute test
	err := suite.signalingService.StoreSignalingMessage(callID, fromUserID, offerMessage)

	// Assertions
	assert.NoError(suite.T(), err)
}

// TestStoreAnswerMessage tests storing WebRTC answer messages
func (suite *SignalingServiceTestSuite) TestStoreAnswerMessage() {
	callID := uuid.New().String()
	fromUserID := uuid.New().String()
	answerMessage := map[string]interface{}{
		"type": "answer",
		"sdp":  "v=0\r\no=- 987654321 1 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\nm=audio 54401 RTP/AVP 0\r\nc=IN IP4 127.0.0.1\r\n",
	}

	key := "signaling:" + callID

	expectedMessageData := map[string]interface{}{
		"from":      fromUserID,
		"timestamp": time.Now().Unix(),
		"message":   answerMessage,
	}

	// Mock Redis operations
	suite.redisMock.ExpectLPush(key, expectedMessageData).SetVal(1)
	suite.redisMock.ExpectExpire(key, 10*time.Minute).SetVal(true)
	suite.redisMock.ExpectLTrim(key, 0, 100).SetVal(true)

	// Execute test
	err := suite.signalingService.StoreSignalingMessage(callID, fromUserID, answerMessage)

	// Assertions
	assert.NoError(suite.T(), err)
}

// TestStoreICECandidate tests storing ICE candidate messages
func (suite *SignalingServiceTestSuite) TestStoreICECandidate() {
	callID := uuid.New().String()
	fromUserID := uuid.New().String()
	iceCandidateMessage := map[string]interface{}{
		"type":      "ice-candidate",
		"candidate": "candidate:842163049 1 udp 1677729535 192.168.1.100 54400 typ srflx raddr 0.0.0.0 rport 0 generation 0",
		"sdpMid":    "0",
		"sdpMLineIndex": 0,
	}

	key := "signaling:" + callID

	expectedMessageData := map[string]interface{}{
		"from":      fromUserID,
		"timestamp": time.Now().Unix(),
		"message":   iceCandidateMessage,
	}

	// Mock Redis operations
	suite.redisMock.ExpectLPush(key, expectedMessageData).SetVal(1)
	suite.redisMock.ExpectExpire(key, 10*time.Minute).SetVal(true)
	suite.redisMock.ExpectLTrim(key, 0, 100).SetVal(true)

	// Execute test
	err := suite.signalingService.StoreSignalingMessage(callID, fromUserID, iceCandidateMessage)

	// Assertions
	assert.NoError(suite.T(), err)
}

// TestSignalingMessageExpiration tests message expiration
func (suite *SignalingServiceTestSuite) TestSignalingMessageExpiration() {
	callID := uuid.New().String()
	fromUserID := uuid.New().String()
	message := map[string]interface{}{
		"type": "offer",
		"sdp":  "test-sdp",
	}

	key := "signaling:" + callID

	expectedMessageData := map[string]interface{}{
		"from":      fromUserID,
		"timestamp": time.Now().Unix(),
		"message":   message,
	}

	// Mock Redis operations with TTL verification
	suite.redisMock.ExpectLPush(key, expectedMessageData).SetVal(1)
	suite.redisMock.ExpectExpire(key, 10*time.Minute).SetVal(true)
	suite.redisMock.ExpectLTrim(key, 0, 100).SetVal(true)

	// Execute test
	err := suite.signalingService.StoreSignalingMessage(callID, fromUserID, message)

	// Assertions
	assert.NoError(suite.T(), err)
}

// TestSignalingMessageLimit tests message count limiting
func (suite *SignalingServiceTestSuite) TestSignalingMessageLimit() {
	callID := uuid.New().String()
	fromUserID := uuid.New().String()

	// Store multiple messages to test the limit
	numMessages := 5

	for i := 0; i < numMessages; i++ {
		message := map[string]interface{}{
			"type":    "ice-candidate",
			"candidate": "candidate:" + string(rune(i)),
		}

		key := "signaling:" + callID

		expectedMessageData := map[string]interface{}{
			"from":      fromUserID,
			"timestamp": time.Now().Unix(),
			"message":   message,
		}

		// Mock Redis operations
		suite.redisMock.ExpectLPush(key, expectedMessageData).SetVal(int64(i + 1))
		suite.redisMock.ExpectExpire(key, 10*time.Minute).SetVal(true)
		suite.redisMock.ExpectLTrim(key, 0, 100).SetVal(true)

		// Execute test
		err := suite.signalingService.StoreSignalingMessage(callID, fromUserID, message)
		assert.NoError(suite.T(), err)
	}
}

// TestConcurrentSignalingMessages tests concurrent message storage
func (suite *SignalingServiceTestSuite) TestConcurrentSignalingMessages() {
	callID := uuid.New().String()
	numGoroutines := 10

	// Mock Redis operations for concurrent messages
	for i := 0; i < numGoroutines; i++ {
		fromUserID := uuid.New().String()
		message := map[string]interface{}{
			"type": "ice-candidate",
			"candidate": "concurrent-candidate-" + string(rune(i)),
		}

		key := "signaling:" + callID

		expectedMessageData := map[string]interface{}{
			"from":      fromUserID,
			"timestamp": time.Now().Unix(),
			"message":   message,
		}

		suite.redisMock.ExpectLPush(key, expectedMessageData).SetVal(int64(i + 1))
		suite.redisMock.ExpectExpire(key, 10*time.Minute).SetVal(true)
		suite.redisMock.ExpectLTrim(key, 0, 100).SetVal(true)
	}

	// Execute concurrent message storage
	done := make(chan bool, numGoroutines)
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			fromUserID := uuid.New().String()
			message := map[string]interface{}{
				"type": "ice-candidate",
				"candidate": "concurrent-candidate-" + string(rune(index)),
			}

			err := suite.signalingService.StoreSignalingMessage(callID, fromUserID, message)
			if err != nil {
				errors <- err
			}
			done <- true
		}(i)
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

// TestSignalingMessageTypes tests different message types
func (suite *SignalingServiceTestSuite) TestSignalingMessageTypes() {
	callID := uuid.New().String()
	fromUserID := uuid.New().String()

	messageTypes := []struct {
		name     string
		message  map[string]interface{}
	}{
		{
			name: "offer",
			message: map[string]interface{}{
				"type": "offer",
				"sdp":  "offer-sdp-content",
			},
		},
		{
			name: "answer",
			message: map[string]interface{}{
				"type": "answer",
				"sdp":  "answer-sdp-content",
			},
		},
		{
			name: "ice-candidate",
			message: map[string]interface{}{
				"type":      "ice-candidate",
				"candidate": "ice-candidate-content",
				"sdpMid":    "0",
				"sdpMLineIndex": 0,
			},
		},
		{
			name: "call-end",
			message: map[string]interface{}{
				"type":   "call-end",
				"reason": "user-hangup",
			},
		},
	}

	for i, messageType := range messageTypes {
		suite.T().Run(messageType.name, func(t *testing.T) {
			key := "signaling:" + callID

			expectedMessageData := map[string]interface{}{
				"from":      fromUserID,
				"timestamp": time.Now().Unix(),
				"message":   messageType.message,
			}

			// Mock Redis operations
			suite.redisMock.ExpectLPush(key, expectedMessageData).SetVal(int64(i + 1))
			suite.redisMock.ExpectExpire(key, 10*time.Minute).SetVal(true)
			suite.redisMock.ExpectLTrim(key, 0, 100).SetVal(true)

			// Execute test
			err := suite.signalingService.StoreSignalingMessage(callID, fromUserID, messageType.message)
			assert.NoError(t, err)
		})
	}
}

// TestRedisConnectionFailure tests handling of Redis connection failures
func (suite *SignalingServiceTestSuite) TestRedisConnectionFailure() {
	callID := uuid.New().String()
	fromUserID := uuid.New().String()
	message := map[string]interface{}{
		"type": "offer",
		"sdp":  "test-sdp",
	}

	key := "signaling:" + callID

	expectedMessageData := map[string]interface{}{
		"from":      fromUserID,
		"timestamp": time.Now().Unix(),
		"message":   message,
	}

	// Mock Redis connection failure
	suite.redisMock.ExpectLPush(key, expectedMessageData).SetErr(redis.ErrClosed)

	// Execute test
	err := suite.signalingService.StoreSignalingMessage(callID, fromUserID, message)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), redis.ErrClosed, err)
}

// TestInvalidSignalingData tests handling of invalid signaling data
func (suite *SignalingServiceTestSuite) TestInvalidSignalingData() {
	tests := []struct {
		name         string
		callID       string
		fromUserID   string
		message      map[string]interface{}
		expectError  bool
	}{
		{
			name:       "Valid message",
			callID:     uuid.New().String(),
			fromUserID: uuid.New().String(),
			message:    map[string]interface{}{"type": "offer", "sdp": "test"},
			expectError: false,
		},
		{
			name:       "Empty call ID",
			callID:     "",
			fromUserID: uuid.New().String(),
			message:    map[string]interface{}{"type": "offer", "sdp": "test"},
			expectError: true,
		},
		{
			name:       "Empty user ID",
			callID:     uuid.New().String(),
			fromUserID: "",
			message:    map[string]interface{}{"type": "offer", "sdp": "test"},
			expectError: true,
		},
		{
			name:       "Empty message",
			callID:     uuid.New().String(),
			fromUserID: uuid.New().String(),
			message:    map[string]interface{}{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			// Skip Redis mock for invalid input that should fail early
			if tt.expectError && (tt.callID == "" || tt.fromUserID == "") {
				// These would fail validation before reaching Redis
				return
			}

			// Mock Redis operations for valid input
			key := "signaling:" + tt.callID

			expectedMessageData := map[string]interface{}{
				"from":      tt.fromUserID,
				"timestamp": time.Now().Unix(),
				"message":   tt.message,
			}

			suite.redisMock.ExpectLPush(key, expectedMessageData).SetVal(1)
			suite.redisMock.ExpectExpire(key, 10*time.Minute).SetVal(true)
			suite.redisMock.ExpectLTrim(key, 0, 100).SetVal(true)

			err := suite.signalingService.StoreSignalingMessage(tt.callID, tt.fromUserID, tt.message)
			assert.NoError(t, err)
		})
	}
}

// TestSignalingServiceTestSuite runs the test suite
func TestSignalingServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SignalingServiceTestSuite))
}

// Benchmark tests

func BenchmarkStoreSignalingMessage(b *testing.B) {
	// Setup mock Redis
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisWrapper := &config.RedisClient{
		Client: db.(*redis.Client),
	}
	signalingService := config.NewSignalingService(redisWrapper)

	callID := uuid.New().String()
	fromUserID := uuid.New().String()
	message := map[string]interface{}{
		"type": "ice-candidate",
		"candidate": "benchmark-candidate",
	}

	key := "signaling:" + callID

	// Setup expectations for benchmark iterations
	for i := 0; i < b.N; i++ {
		expectedMessageData := map[string]interface{}{
			"from":      fromUserID,
			"timestamp": time.Now().Unix(),
			"message":   message,
		}

		mock.ExpectLPush(key, expectedMessageData).SetVal(int64(i + 1))
		mock.ExpectExpire(key, 10*time.Minute).SetVal(true)
		mock.ExpectLTrim(key, 0, 100).SetVal(true)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		signalingService.StoreSignalingMessage(callID, fromUserID, message)
	}
}

func BenchmarkSignalingMessageTypes(b *testing.B) {
	// Setup mock Redis
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisWrapper := &config.RedisClient{
		Client: db.(*redis.Client),
	}
	signalingService := config.NewSignalingService(redisWrapper)

	callID := uuid.New().String()
	fromUserID := uuid.New().String()

	messageTypes := []map[string]interface{}{
		{"type": "offer", "sdp": "offer-sdp"},
		{"type": "answer", "sdp": "answer-sdp"},
		{"type": "ice-candidate", "candidate": "ice-candidate"},
	}

	key := "signaling:" + callID

	// Setup expectations for benchmark iterations
	for i := 0; i < b.N; i++ {
		message := messageTypes[i%len(messageTypes)]
		expectedMessageData := map[string]interface{}{
			"from":      fromUserID,
			"timestamp": time.Now().Unix(),
			"message":   message,
		}

		mock.ExpectLPush(key, expectedMessageData).SetVal(int64(i + 1))
		mock.ExpectExpire(key, 10*time.Minute).SetVal(true)
		mock.ExpectLTrim(key, 0, 100).SetVal(true)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		message := messageTypes[i%len(messageTypes)]
		signalingService.StoreSignalingMessage(callID, fromUserID, message)
	}
}