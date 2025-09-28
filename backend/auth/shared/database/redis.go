package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host               string        `mapstructure:"host" validate:"required"`
	Port               int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	Password           string        `mapstructure:"password"`
	Database           int           `mapstructure:"database" validate:"min=0,max=15"`
	MaxRetries         int           `mapstructure:"max_retries"`
	MinRetryBackoff    time.Duration `mapstructure:"min_retry_backoff"`
	MaxRetryBackoff    time.Duration `mapstructure:"max_retry_backoff"`
	DialTimeout        time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout        time.Duration `mapstructure:"read_timeout"`
	WriteTimeout       time.Duration `mapstructure:"write_timeout"`
	PoolSize           int           `mapstructure:"pool_size"`
	MinIdleConns       int           `mapstructure:"min_idle_conns"`
	MaxIdleConns       int           `mapstructure:"max_idle_conns"`
	ConnMaxIdleTime    time.Duration `mapstructure:"conn_max_idle_time"`
	ConnMaxLifetime    time.Duration `mapstructure:"conn_max_lifetime"`
	EnableTLS          bool          `mapstructure:"enable_tls"`
	InsecureSkipVerify bool          `mapstructure:"insecure_skip_verify"`
}

// RedisDB wraps redis.Client with additional functionality
type RedisDB struct {
	Client *redis.Client
	Config *RedisConfig
	ctx    context.Context
}

// DefaultRedisConfig returns default Redis configuration
func DefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:               "localhost",
		Port:               6379,
		Database:           0,
		MaxRetries:         3,
		MinRetryBackoff:    8 * time.Millisecond,
		MaxRetryBackoff:    512 * time.Millisecond,
		DialTimeout:        5 * time.Second,
		ReadTimeout:        3 * time.Second,
		WriteTimeout:       3 * time.Second,
		PoolSize:           10,
		MinIdleConns:       5,
		MaxIdleConns:       10,
		ConnMaxIdleTime:    30 * time.Minute,
		ConnMaxLifetime:    time.Hour,
		EnableTLS:          false,
		InsecureSkipVerify: false,
	}
}

// NewRedisDB creates a new Redis database connection
func NewRedisDB(config *RedisConfig) (*RedisDB, error) {
	if config == nil {
		config = DefaultRedisConfig()
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	options := &redis.Options{
		Addr:               addr,
		Password:           config.Password,
		DB:                 config.Database,
		MaxRetries:         config.MaxRetries,
		MinRetryBackoff:    config.MinRetryBackoff,
		MaxRetryBackoff:    config.MaxRetryBackoff,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		ConnMaxIdleTime:    config.ConnMaxIdleTime,
		ConnMaxLifetime:    config.ConnMaxLifetime,
	}

	if config.EnableTLS {
		options.TLSConfig = &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
		}
	}

	client := redis.NewClient(options)

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Printf("Successfully connected to Redis: %s", addr)

	return &RedisDB{
		Client: client,
		Config: config,
		ctx:    ctx,
	}, nil
}

// Close closes the Redis connection
func (r *RedisDB) Close() error {
	if r.Client != nil {
		log.Println("Closing Redis connection")
		return r.Client.Close()
	}
	return nil
}

// HealthCheck performs a health check on the Redis connection
func (r *RedisDB) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}

	return nil
}

// GetStats returns Redis client statistics
func (r *RedisDB) GetStats() *redis.PoolStats {
	return r.Client.PoolStats()
}

// Basic Operations

// Set sets a key-value pair with optional expiration
func (r *RedisDB) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func (r *RedisDB) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// Del deletes one or more keys
func (r *RedisDB) Del(ctx context.Context, keys ...string) error {
	return r.Client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (r *RedisDB) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.Client.Exists(ctx, keys...).Result()
}

// Expire sets expiration for a key
func (r *RedisDB) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.Client.Expire(ctx, key, expiration).Err()
}

// TTL returns the time to live for a key
func (r *RedisDB) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.Client.TTL(ctx, key).Result()
}

// Hash Operations

// HSet sets a field in a hash
func (r *RedisDB) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.Client.HSet(ctx, key, values...).Err()
}

// HGet gets a field from a hash
func (r *RedisDB) HGet(ctx context.Context, key, field string) (string, error) {
	return r.Client.HGet(ctx, key, field).Result()
}

// HGetAll gets all fields from a hash
func (r *RedisDB) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, key).Result()
}

// HDel deletes fields from a hash
func (r *RedisDB) HDel(ctx context.Context, key string, fields ...string) error {
	return r.Client.HDel(ctx, key, fields...).Err()
}

// HExists checks if a field exists in a hash
func (r *RedisDB) HExists(ctx context.Context, key, field string) (bool, error) {
	return r.Client.HExists(ctx, key, field).Result()
}

// List Operations

// LPush pushes elements to the head of a list
func (r *RedisDB) LPush(ctx context.Context, key string, values ...interface{}) error {
	return r.Client.LPush(ctx, key, values...).Err()
}

// RPush pushes elements to the tail of a list
func (r *RedisDB) RPush(ctx context.Context, key string, values ...interface{}) error {
	return r.Client.RPush(ctx, key, values...).Err()
}

// LPop pops an element from the head of a list
func (r *RedisDB) LPop(ctx context.Context, key string) (string, error) {
	return r.Client.LPop(ctx, key).Result()
}

// RPop pops an element from the tail of a list
func (r *RedisDB) RPop(ctx context.Context, key string) (string, error) {
	return r.Client.RPop(ctx, key).Result()
}

// LRange gets a range of elements from a list
func (r *RedisDB) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.Client.LRange(ctx, key, start, stop).Result()
}

// LLen gets the length of a list
func (r *RedisDB) LLen(ctx context.Context, key string) (int64, error) {
	return r.Client.LLen(ctx, key).Result()
}

// Set Operations

// SAdd adds members to a set
func (r *RedisDB) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.Client.SAdd(ctx, key, members...).Err()
}

// SRem removes members from a set
func (r *RedisDB) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.Client.SRem(ctx, key, members...).Err()
}

// SMembers gets all members of a set
func (r *RedisDB) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.Client.SMembers(ctx, key).Result()
}

// SIsMember checks if a member exists in a set
func (r *RedisDB) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return r.Client.SIsMember(ctx, key, member).Result()
}

// SCard gets the cardinality of a set
func (r *RedisDB) SCard(ctx context.Context, key string) (int64, error) {
	return r.Client.SCard(ctx, key).Result()
}

// Sorted Set Operations

// ZAdd adds members to a sorted set
func (r *RedisDB) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return r.Client.ZAdd(ctx, key, members...).Err()
}

// ZRem removes members from a sorted set
func (r *RedisDB) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return r.Client.ZRem(ctx, key, members...).Err()
}

// ZRange gets a range of members from a sorted set
func (r *RedisDB) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.Client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores gets a range of members with scores from a sorted set
func (r *RedisDB) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return r.Client.ZRangeWithScores(ctx, key, start, stop).Result()
}

// ZScore gets the score of a member in a sorted set
func (r *RedisDB) ZScore(ctx context.Context, key, member string) (float64, error) {
	return r.Client.ZScore(ctx, key, member).Result()
}

// ZCard gets the cardinality of a sorted set
func (r *RedisDB) ZCard(ctx context.Context, key string) (int64, error) {
	return r.Client.ZCard(ctx, key).Result()
}

// Pub/Sub Operations

// Publish publishes a message to a channel
func (r *RedisDB) Publish(ctx context.Context, channel string, message interface{}) error {
	return r.Client.Publish(ctx, channel, message).Err()
}

// Subscribe subscribes to channels
func (r *RedisDB) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return r.Client.Subscribe(ctx, channels...)
}

// PSubscribe subscribes to channel patterns
func (r *RedisDB) PSubscribe(ctx context.Context, patterns ...string) *redis.PubSub {
	return r.Client.PSubscribe(ctx, patterns...)
}

// Transaction Operations

// TxPipeline creates a transaction pipeline
func (r *RedisDB) TxPipeline() redis.Pipeliner {
	return r.Client.TxPipeline()
}

// Watch watches keys for changes
func (r *RedisDB) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	return r.Client.Watch(ctx, fn, keys...)
}

// Pipeline Operations

// Pipeline creates a pipeline
func (r *RedisDB) Pipeline() redis.Pipeliner {
	return r.Client.Pipeline()
}

// Script Operations

// Eval evaluates a Lua script
func (r *RedisDB) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	return r.Client.Eval(ctx, script, keys, args...)
}

// EvalSha evaluates a Lua script by SHA1
func (r *RedisDB) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	return r.Client.EvalSha(ctx, sha1, keys, args...)
}

// ScriptExists checks if scripts exist
func (r *RedisDB) ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd {
	return r.Client.ScriptExists(ctx, hashes...)
}

// ScriptLoad loads a script
func (r *RedisDB) ScriptLoad(ctx context.Context, script string) *redis.StringCmd {
	return r.Client.ScriptLoad(ctx, script)
}

// Atomic Operations

// Incr increments a key
func (r *RedisDB) Incr(ctx context.Context, key string) (int64, error) {
	return r.Client.Incr(ctx, key).Result()
}

// IncrBy increments a key by a value
func (r *RedisDB) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.Client.IncrBy(ctx, key, value).Result()
}

// Decr decrements a key
func (r *RedisDB) Decr(ctx context.Context, key string) (int64, error) {
	return r.Client.Decr(ctx, key).Result()
}

// DecrBy decrements a key by a value
func (r *RedisDB) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.Client.DecrBy(ctx, key, value).Result()
}

// SetNX sets a key only if it doesn't exist
func (r *RedisDB) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.Client.SetNX(ctx, key, value, expiration).Result()
}

// GetSet sets a key and returns the old value
func (r *RedisDB) GetSet(ctx context.Context, key string, value interface{}) (string, error) {
	return r.Client.GetSet(ctx, key, value).Result()
}

// Geo Operations

// GeoAdd adds geospatial items
func (r *RedisDB) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) error {
	return r.Client.GeoAdd(ctx, key, geoLocation...).Err()
}

// GeoRadius searches for members within a radius
func (r *RedisDB) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error) {
	return r.Client.GeoRadius(ctx, key, longitude, latitude, query).Result()
}

// FlushDB flushes the current database
func (r *RedisDB) FlushDB(ctx context.Context) error {
	return r.Client.FlushDB(ctx).Err()
}

// FlushAll flushes all databases
func (r *RedisDB) FlushAll(ctx context.Context) error {
	return r.Client.FlushAll(ctx).Err()
}