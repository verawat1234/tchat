package config

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// BaseRedisService provides common Redis operations for all services
type BaseRedisService struct {
	redis *RedisClient
	ctx   context.Context
}

// NewBaseRedisService creates a new base Redis service
func NewBaseRedisService(redisClient *RedisClient) *BaseRedisService {
	return &BaseRedisService{
		redis: redisClient,
		ctx:   context.Background(),
	}
}

// SetWithTTL sets a key-value pair with TTL
func (brs *BaseRedisService) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	return brs.redis.Client.Set(brs.ctx, key, value, ttl).Err()
}

// SetHashWithTTL sets a hash with TTL using pipeline for atomicity
func (brs *BaseRedisService) SetHashWithTTL(key string, data map[string]interface{}, ttl time.Duration) error {
	pipe := brs.redis.Client.Pipeline()
	pipe.HMSet(brs.ctx, key, data)
	pipe.Expire(brs.ctx, key, ttl)

	_, err := pipe.Exec(brs.ctx)
	return err
}

// GetHash retrieves all fields from a hash
func (brs *BaseRedisService) GetHash(key string) (map[string]string, error) {
	return brs.redis.Client.HGetAll(brs.ctx, key).Result()
}

// DeleteKey deletes a key
func (brs *BaseRedisService) DeleteKey(key string) error {
	return brs.redis.Client.Del(brs.ctx, key).Err()
}

// ListPushWithTTLAndTrim pushes to list, sets TTL, and trims in a pipeline
func (brs *BaseRedisService) ListPushWithTTLAndTrim(key string, value interface{}, ttl time.Duration, maxLength int64) error {
	pipe := brs.redis.Client.Pipeline()
	pipe.LPush(brs.ctx, key, value)
	pipe.Expire(brs.ctx, key, ttl)
	pipe.LTrim(brs.ctx, key, 0, maxLength-1)

	_, err := pipe.Exec(brs.ctx)
	return err
}

// GetList retrieves all items from a list
func (brs *BaseRedisService) GetList(key string) ([]string, error) {
	return brs.redis.Client.LRange(brs.ctx, key, 0, -1).Result()
}

// KeyExists checks if a key exists
func (brs *BaseRedisService) KeyExists(key string) (bool, error) {
	result := brs.redis.Client.Exists(brs.ctx, key)
	count, err := result.Result()
	return count > 0, err
}

// SetTTL sets TTL for an existing key
func (brs *BaseRedisService) SetTTL(key string, ttl time.Duration) error {
	return brs.redis.Client.Expire(brs.ctx, key, ttl).Err()
}

// GetTTL gets remaining TTL for a key
func (brs *BaseRedisService) GetTTL(key string) (time.Duration, error) {
	return brs.redis.Client.TTL(brs.ctx, key).Result()
}

// BatchOperation executes multiple Redis operations in a pipeline
func (brs *BaseRedisService) BatchOperation(operations func(pipe redis.Pipeliner)) error {
	pipe := brs.redis.Client.Pipeline()
	operations(pipe)
	_, err := pipe.Exec(brs.ctx)
	return err
}

// KeyPattern generates consistent key patterns for services
type KeyPattern struct {
	Prefix string
	Type   string
}

// GenerateKey creates a consistent key with pattern
func (kp KeyPattern) GenerateKey(id string) string {
	return fmt.Sprintf("%s:%s:%s", kp.Prefix, kp.Type, id)
}

// Common key patterns used across the calling service
var (
	PresenceKeyPattern  = KeyPattern{Prefix: "presence", Type: "user"}
	CallStateKeyPattern = KeyPattern{Prefix: "call", Type: "state"}
	SignalingKeyPattern = KeyPattern{Prefix: "signaling", Type: "messages"}
	CallRoomKeyPattern  = KeyPattern{Prefix: "call", Type: "room"}
)

// RedisServiceInterface defines common operations for Redis-based services
type RedisServiceInterface interface {
	GetKey(id string) string
	SetTTL() time.Duration
	Cleanup(id string) error
}

// ServiceMetrics holds common metrics for Redis services
type ServiceMetrics struct {
	OperationsCount int64
	ErrorsCount     int64
	AvgLatency      time.Duration
	LastOperation   time.Time
}

// MetricsCollector provides basic metrics collection for Redis services
type MetricsCollector struct {
	metrics map[string]*ServiceMetrics
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: make(map[string]*ServiceMetrics),
	}
}

// RecordOperation records an operation for metrics
func (mc *MetricsCollector) RecordOperation(serviceName string, duration time.Duration, isError bool) {
	if mc.metrics[serviceName] == nil {
		mc.metrics[serviceName] = &ServiceMetrics{}
	}

	metrics := mc.metrics[serviceName]
	metrics.OperationsCount++
	if isError {
		metrics.ErrorsCount++
	}

	// Simple moving average for latency
	if metrics.OperationsCount == 1 {
		metrics.AvgLatency = duration
	} else {
		metrics.AvgLatency = (metrics.AvgLatency + duration) / 2
	}

	metrics.LastOperation = time.Now()
}

// GetMetrics returns metrics for a service
func (mc *MetricsCollector) GetMetrics(serviceName string) *ServiceMetrics {
	return mc.metrics[serviceName]
}