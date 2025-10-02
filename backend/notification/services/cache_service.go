package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCacheService implements CacheService using Redis
type RedisCacheService struct {
	client *redis.Client
}

// NewRedisCacheService creates a new Redis cache service
func NewRedisCacheService(redisURL string) (CacheService, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCacheService{client: client}, nil
}

// Get retrieves a value from cache
func (r *RedisCacheService) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrCacheKeyNotFound
		}
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		// If JSON unmarshal fails, return as string
		return val, nil
	}

	return result, nil
}

// Set stores a value in cache with expiration
func (r *RedisCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
	}

	if err := r.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return nil
}

// Delete removes a key from cache
func (r *RedisCacheService) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// Exists checks if a key exists in cache
func (r *RedisCacheService) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}
	return count > 0, nil
}

// GetMulti retrieves multiple values from cache
func (r *RedisCacheService) GetMulti(ctx context.Context, keys []string) (map[string]interface{}, error) {
	if len(keys) == 0 {
		return make(map[string]interface{}), nil
	}

	values, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get multiple keys: %w", err)
	}

	result := make(map[string]interface{})
	for i, key := range keys {
		if values[i] != nil {
			var value interface{}
			if strVal, ok := values[i].(string); ok {
				if err := json.Unmarshal([]byte(strVal), &value); err != nil {
					// If JSON unmarshal fails, use as string
					value = strVal
				}
			} else {
				value = values[i]
			}
			result[key] = value
		}
	}

	return result, nil
}

// SetMulti stores multiple values in cache with expiration
func (r *RedisCacheService) SetMulti(ctx context.Context, items map[string]interface{}, expiration time.Duration) error {
	pipe := r.client.Pipeline()

	for key, value := range items {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}
		pipe.Set(ctx, key, data, expiration)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to set multiple keys: %w", err)
	}

	return nil
}

// FlushAll removes all keys from cache
func (r *RedisCacheService) FlushAll(ctx context.Context) error {
	if err := r.client.FlushAll(ctx).Err(); err != nil {
		return fmt.Errorf("failed to flush cache: %w", err)
	}
	return nil
}

// Increment increments a numeric value in cache
func (r *RedisCacheService) Increment(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key %s: %w", key, err)
	}
	return val, nil
}

// SetWithExpiry stores a value in cache with expiration (alias for Set)
func (r *RedisCacheService) SetWithExpiry(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Set(ctx, key, value, expiration)
}

// Incr increments a counter in cache
func (r *RedisCacheService) Incr(ctx context.Context, key string) error {
	if err := r.client.Incr(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to increment key %s: %w", key, err)
	}
	return nil
}

// IncrementWithExpiry increments a numeric value and sets expiration
func (r *RedisCacheService) IncrementWithExpiry(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	pipe := r.client.Pipeline()

	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, expiration)

	if _, err := pipe.Exec(ctx); err != nil {
		return 0, fmt.Errorf("failed to increment key %s with expiry: %w", key, err)
	}

	return incrCmd.Val(), nil
}

// Close closes the Redis connection
func (r *RedisCacheService) Close() error {
	return r.client.Close()
}

// InMemoryCacheService provides a simple in-memory cache for development/testing
type InMemoryCacheService struct {
	data map[string]cacheItem
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

// NewInMemoryCacheService creates a new in-memory cache service
func NewInMemoryCacheService() CacheService {
	return &InMemoryCacheService{
		data: make(map[string]cacheItem),
	}
}

// Get retrieves a value from in-memory cache
func (m *InMemoryCacheService) Get(ctx context.Context, key string) (interface{}, error) {
	item, exists := m.data[key]
	if !exists {
		return nil, ErrCacheKeyNotFound
	}

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		delete(m.data, key)
		return nil, ErrCacheKeyNotFound
	}

	return item.value, nil
}

// Set stores a value in in-memory cache
func (m *InMemoryCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var expiresAt time.Time
	if expiration > 0 {
		expiresAt = time.Now().Add(expiration)
	}

	m.data[key] = cacheItem{
		value:     value,
		expiresAt: expiresAt,
	}

	return nil
}

// Delete removes a key from in-memory cache
func (m *InMemoryCacheService) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

// Exists checks if a key exists in in-memory cache
func (m *InMemoryCacheService) Exists(ctx context.Context, key string) (bool, error) {
	item, exists := m.data[key]
	if !exists {
		return false, nil
	}

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		delete(m.data, key)
		return false, nil
	}

	return true, nil
}

// GetMulti retrieves multiple values from in-memory cache
func (m *InMemoryCacheService) GetMulti(ctx context.Context, keys []string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, key := range keys {
		if value, err := m.Get(ctx, key); err == nil {
			result[key] = value
		}
	}
	return result, nil
}

// SetMulti stores multiple values in in-memory cache
func (m *InMemoryCacheService) SetMulti(ctx context.Context, items map[string]interface{}, expiration time.Duration) error {
	for key, value := range items {
		if err := m.Set(ctx, key, value, expiration); err != nil {
			return err
		}
	}
	return nil
}

// FlushAll removes all keys from in-memory cache
func (m *InMemoryCacheService) FlushAll(ctx context.Context) error {
	m.data = make(map[string]cacheItem)
	return nil
}

// SetWithExpiry stores a value in cache with expiration (alias for Set)
func (m *InMemoryCacheService) SetWithExpiry(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return m.Set(ctx, key, value, expiration)
}

// Incr increments a counter in cache
func (m *InMemoryCacheService) Incr(ctx context.Context, key string) error {
	item, exists := m.data[key]
	if !exists {
		m.data[key] = cacheItem{value: int64(1), expiresAt: time.Time{}}
		return nil
	}

	if val, ok := item.value.(int64); ok {
		item.value = val + 1
		m.data[key] = item
	} else {
		m.data[key] = cacheItem{value: int64(1), expiresAt: item.expiresAt}
	}
	return nil
}

// Cache-related errors
var (
	ErrCacheKeyNotFound = fmt.Errorf("cache key not found")
)