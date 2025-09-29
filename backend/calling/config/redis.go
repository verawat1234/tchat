package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient holds the Redis client instance
type RedisClient struct {
	Client *redis.Client
	ctx    context.Context
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	MaxRetries int
	PoolSize   int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewRedisConfig creates a new Redis configuration from environment variables
func NewRedisConfig() *RedisConfig {
	config := &RedisConfig{
		Host:         getEnv("REDIS_HOST", "localhost"),
		Port:         getEnvAsInt("REDIS_PORT", 6379),
		Password:     getEnv("REDIS_PASSWORD", ""),
		DB:           getEnvAsInt("REDIS_DB", 0),
		MaxRetries:   getEnvAsInt("REDIS_MAX_RETRIES", 3),
		PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
		MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
		DialTimeout:  getEnvAsDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
		ReadTimeout:  getEnvAsDuration("REDIS_READ_TIMEOUT", 3*time.Second),
		WriteTimeout: getEnvAsDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
	}

	return config
}

// NewRedisClient creates a new Redis client connection
func NewRedisClient() (*RedisClient, error) {
	config := NewRedisConfig()

	// Create Redis client options
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		MaxRetries:   config.MaxRetries,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	ctx := context.Background()

	// Test the connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Printf("Successfully connected to Redis at %s:%d", config.Host, config.Port)

	return &RedisClient{
		Client: rdb,
		ctx:    ctx,
	}, nil
}

// UserPresenceService provides Redis-based user presence operations
type UserPresenceService struct {
	*BaseRedisService
	metrics *MetricsCollector
}

// NewUserPresenceService creates a new user presence service
func NewUserPresenceService(redisClient *RedisClient) *UserPresenceService {
	return &UserPresenceService{
		BaseRedisService: NewBaseRedisService(redisClient),
		metrics:         NewMetricsCollector(),
	}
}

// SetUserOnline sets a user as online with TTL
func (ups *UserPresenceService) SetUserOnline(userID string, callID string) error {
	start := time.Now()
	key := PresenceKeyPattern.GenerateKey(userID)

	presenceData := map[string]interface{}{
		"status":    "online",
		"in_call":   callID != "",
		"call_id":   callID,
		"last_seen": time.Now().Unix(),
	}

	err := ups.SetHashWithTTL(key, presenceData, 5*time.Minute)

	// Record metrics
	ups.metrics.RecordOperation("presence_set_online", time.Since(start), err != nil)
	return err
}

// SetUserOffline removes user presence
func (ups *UserPresenceService) SetUserOffline(userID string) error {
	start := time.Now()
	key := PresenceKeyPattern.GenerateKey(userID)

	err := ups.DeleteKey(key)
	ups.metrics.RecordOperation("presence_set_offline", time.Since(start), err != nil)
	return err
}

// GetUserPresence retrieves user presence status
func (ups *UserPresenceService) GetUserPresence(userID string) (map[string]string, error) {
	start := time.Now()
	key := PresenceKeyPattern.GenerateKey(userID)

	presence, err := ups.GetHash(key)
	ups.metrics.RecordOperation("presence_get", time.Since(start), err != nil)

	if err != nil {
		return nil, err
	}

	// If no presence data, user is offline
	if len(presence) == 0 {
		return map[string]string{
			"status":    "offline",
			"in_call":   "false",
			"call_id":   "",
			"last_seen": "0",
		}, nil
	}

	return presence, nil
}

// GetMetrics returns service metrics
func (ups *UserPresenceService) GetMetrics() *ServiceMetrics {
	return ups.metrics.GetMetrics("presence")
}

// CallStateService provides Redis-based call state operations
type CallStateService struct {
	*BaseRedisService
	metrics *MetricsCollector
}

// NewCallStateService creates a new call state service
func NewCallStateService(redisClient *RedisClient) *CallStateService {
	return &CallStateService{
		BaseRedisService: NewBaseRedisService(redisClient),
		metrics:         NewMetricsCollector(),
	}
}

// SetCallState stores call session state in Redis
func (css *CallStateService) SetCallState(callID string, state map[string]interface{}) error {
	start := time.Now()
	key := CallStateKeyPattern.GenerateKey(callID)

	err := css.SetHashWithTTL(key, state, 2*time.Hour)
	css.metrics.RecordOperation("call_state_set", time.Since(start), err != nil)
	return err
}

// GetCallState retrieves call session state from Redis
func (css *CallStateService) GetCallState(callID string) (map[string]string, error) {
	start := time.Now()
	key := CallStateKeyPattern.GenerateKey(callID)

	result, err := css.GetHash(key)
	css.metrics.RecordOperation("call_state_get", time.Since(start), err != nil)
	return result, err
}

// RemoveCallState removes call session state from Redis
func (css *CallStateService) RemoveCallState(callID string) error {
	start := time.Now()
	key := CallStateKeyPattern.GenerateKey(callID)

	err := css.DeleteKey(key)
	css.metrics.RecordOperation("call_state_remove", time.Since(start), err != nil)
	return err
}

// GetMetrics returns service metrics
func (css *CallStateService) GetMetrics() *ServiceMetrics {
	return css.metrics.GetMetrics("call_state")
}

// SignalingService provides Redis-based WebSocket signaling operations
type SignalingService struct {
	*BaseRedisService
	metrics *MetricsCollector
}

// NewSignalingService creates a new signaling service
func NewSignalingService(redisClient *RedisClient) *SignalingService {
	return &SignalingService{
		BaseRedisService: NewBaseRedisService(redisClient),
		metrics:         NewMetricsCollector(),
	}
}

// StoreSignalingMessage stores WebRTC signaling message temporarily
func (ss *SignalingService) StoreSignalingMessage(fromUserID string, toUserID string, message string) error {
	start := time.Now()
	key := SignalingKeyPattern.GenerateKey(toUserID)

	messageData := map[string]interface{}{
		"from":      fromUserID,
		"timestamp": time.Now().Unix(),
		"message":   message,
	}

	err := ss.ListPushWithTTLAndTrim(key, messageData, 10*time.Minute, 100)
	ss.metrics.RecordOperation("signaling_store", time.Since(start), err != nil)
	return err
}

// GetSignalingMessages retrieves signaling messages for a user
func (ss *SignalingService) GetSignalingMessages(userID string) ([]string, error) {
	start := time.Now()
	key := SignalingKeyPattern.GenerateKey(userID)

	messages, err := ss.GetList(key)
	ss.metrics.RecordOperation("signaling_get", time.Since(start), err != nil)
	return messages, err
}

// GetMetrics returns service metrics
func (ss *SignalingService) GetMetrics() *ServiceMetrics {
	return ss.metrics.GetMetrics("signaling")
}

// Close closes the Redis connection
func (rc *RedisClient) Close() error {
	return rc.Client.Close()
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}