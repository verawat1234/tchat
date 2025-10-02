package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"tchat.dev/shared/database"
)

// CacheManager handles caching operations using Redis
type CacheManager struct {
	redis  *database.RedisDB
	prefix string
}

// NewCacheManager creates a new cache manager
func NewCacheManager(redis *database.RedisDB, prefix string) *CacheManager {
	return &CacheManager{
		redis:  redis,
		prefix: prefix,
	}
}

// Cache keys
const (
	// User caching
	UserProfileKey    = "user:profile:%s"
	UserSessionKey    = "user:session:%s"
	UserSettingsKey   = "user:settings:%s"
	UserOnlineKey     = "user:online:%s"
	UserTokenKey      = "user:token:%s"

	// Dialog caching
	DialogKey           = "dialog:%s"
	DialogParticipants  = "dialog:participants:%s"
	DialogMessages      = "dialog:messages:%s"
	DialogUnreadCount   = "dialog:unread:%s:%s" // dialog_id:user_id
	DialogLastMessage   = "dialog:last_message:%s"

	// Message caching
	MessageKey        = "message:%s"
	MessageReactions  = "message:reactions:%s"
	MessageReadBy     = "message:read_by:%s"

	// Real-time messaging
	UserConnectionsKey = "connections:user:%s"
	DialogConnectionsKey = "connections:dialog:%s"
	TypingIndicatorsKey = "typing:%s" // dialog_id
	PresenceKey = "presence:%s" // user_id

	// Rate limiting
	RateLimitKey = "rate_limit:%s:%s" // type:identifier

	// OTP caching
	OTPKey = "otp:%s:%s" // phone:purpose

	// Payment caching
	WalletBalanceKey = "wallet:balance:%s" // wallet_id
	TransactionKey = "transaction:%s"
	PaymentSessionKey = "payment:session:%s"

	// Search caching
	SearchResultsKey = "search:%s:%s" // type:query_hash

	// Analytics caching
	UserStatsKey = "stats:user:%s:%s" // user_id:date
	DialogStatsKey = "stats:dialog:%s:%s" // dialog_id:date
)

// TTL constants
const (
	ShortTTL   = 5 * time.Minute
	MediumTTL  = 30 * time.Minute
	LongTTL    = 2 * time.Hour
	DayTTL     = 24 * time.Hour
	WeekTTL    = 7 * 24 * time.Hour
)

// formatKey formats a cache key with prefix
func (c *CacheManager) formatKey(key string) string {
	if c.prefix != "" {
		return fmt.Sprintf("%s:%s", c.prefix, key)
	}
	return key
}

// Generic caching operations

// Set stores a value in cache with TTL
func (c *CacheManager) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.redis.Set(ctx, c.formatKey(key), data, ttl)
}

// Get retrieves a value from cache
func (c *CacheManager) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.redis.Get(ctx, c.formatKey(key))
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal cached value: %w", err)
	}

	return nil
}

// Delete removes a key from cache
func (c *CacheManager) Delete(ctx context.Context, keys ...string) error {
	formattedKeys := make([]string, len(keys))
	for i, key := range keys {
		formattedKeys[i] = c.formatKey(key)
	}
	return c.redis.Del(ctx, formattedKeys...)
}

// Exists checks if a key exists in cache
func (c *CacheManager) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.redis.Exists(ctx, c.formatKey(key))
	return count > 0, err
}

// User caching operations

// SetUserProfile caches user profile data
func (c *CacheManager) SetUserProfile(ctx context.Context, userID string, profile interface{}) error {
	key := fmt.Sprintf(UserProfileKey, userID)
	return c.Set(ctx, key, profile, LongTTL)
}

// GetUserProfile retrieves cached user profile
func (c *CacheManager) GetUserProfile(ctx context.Context, userID string, dest interface{}) error {
	key := fmt.Sprintf(UserProfileKey, userID)
	return c.Get(ctx, key, dest)
}

// SetUserSession caches user session data
func (c *CacheManager) SetUserSession(ctx context.Context, sessionID string, session interface{}) error {
	key := fmt.Sprintf(UserSessionKey, sessionID)
	return c.Set(ctx, key, session, DayTTL)
}

// GetUserSession retrieves cached user session
func (c *CacheManager) GetUserSession(ctx context.Context, sessionID string, dest interface{}) error {
	key := fmt.Sprintf(UserSessionKey, sessionID)
	return c.Get(ctx, key, dest)
}

// SetUserOnline marks user as online
func (c *CacheManager) SetUserOnline(ctx context.Context, userID string) error {
	key := fmt.Sprintf(UserOnlineKey, userID)
	return c.redis.Set(ctx, c.formatKey(key), time.Now().Unix(), ShortTTL)
}

// IsUserOnline checks if user is online
func (c *CacheManager) IsUserOnline(ctx context.Context, userID string) (bool, error) {
	key := fmt.Sprintf(UserOnlineKey, userID)
	return c.Exists(ctx, key)
}

// Dialog caching operations

// SetDialog caches dialog data
func (c *CacheManager) SetDialog(ctx context.Context, dialogID string, dialog interface{}) error {
	key := fmt.Sprintf(DialogKey, dialogID)
	return c.Set(ctx, key, dialog, MediumTTL)
}

// GetDialog retrieves cached dialog
func (c *CacheManager) GetDialog(ctx context.Context, dialogID string, dest interface{}) error {
	key := fmt.Sprintf(DialogKey, dialogID)
	return c.Get(ctx, key, dest)
}

// SetDialogParticipants caches dialog participants
func (c *CacheManager) SetDialogParticipants(ctx context.Context, dialogID string, participants interface{}) error {
	key := fmt.Sprintf(DialogParticipants, dialogID)
	return c.Set(ctx, key, participants, MediumTTL)
}

// GetDialogParticipants retrieves cached dialog participants
func (c *CacheManager) GetDialogParticipants(ctx context.Context, dialogID string, dest interface{}) error {
	key := fmt.Sprintf(DialogParticipants, dialogID)
	return c.Get(ctx, key, dest)
}

// SetUnreadCount sets unread message count for user in dialog
func (c *CacheManager) SetUnreadCount(ctx context.Context, dialogID, userID string, count int) error {
	key := fmt.Sprintf(DialogUnreadCount, dialogID, userID)
	return c.redis.Set(ctx, c.formatKey(key), count, DayTTL)
}

// GetUnreadCount gets unread message count for user in dialog
func (c *CacheManager) GetUnreadCount(ctx context.Context, dialogID, userID string) (int, error) {
	key := fmt.Sprintf(DialogUnreadCount, dialogID, userID)
	data, err := c.redis.Get(ctx, c.formatKey(key))
	if err != nil {
		return 0, err
	}

	var count int
	if err := json.Unmarshal([]byte(data), &count); err != nil {
		return 0, err
	}

	return count, nil
}

// IncrementUnreadCount increments unread count for user in dialog
func (c *CacheManager) IncrementUnreadCount(ctx context.Context, dialogID, userID string) error {
	key := fmt.Sprintf(DialogUnreadCount, dialogID, userID)
	_, err := c.redis.Incr(ctx, c.formatKey(key))
	if err == nil {
		// Set TTL if key was created
		c.redis.Expire(ctx, c.formatKey(key), DayTTL)
	}
	return err
}

// Real-time messaging operations

// AddUserConnection adds a WebSocket connection for a user
func (c *CacheManager) AddUserConnection(ctx context.Context, userID, connectionID string) error {
	key := fmt.Sprintf(UserConnectionsKey, userID)
	return c.redis.SAdd(ctx, c.formatKey(key), connectionID)
}

// RemoveUserConnection removes a WebSocket connection for a user
func (c *CacheManager) RemoveUserConnection(ctx context.Context, userID, connectionID string) error {
	key := fmt.Sprintf(UserConnectionsKey, userID)
	return c.redis.SRem(ctx, c.formatKey(key), connectionID)
}

// GetUserConnections gets all connections for a user
func (c *CacheManager) GetUserConnections(ctx context.Context, userID string) ([]string, error) {
	key := fmt.Sprintf(UserConnectionsKey, userID)
	return c.redis.SMembers(ctx, c.formatKey(key))
}

// AddDialogConnection adds a connection to a dialog room
func (c *CacheManager) AddDialogConnection(ctx context.Context, dialogID, connectionID string) error {
	key := fmt.Sprintf(DialogConnectionsKey, dialogID)
	return c.redis.SAdd(ctx, c.formatKey(key), connectionID)
}

// RemoveDialogConnection removes a connection from a dialog room
func (c *CacheManager) RemoveDialogConnection(ctx context.Context, dialogID, connectionID string) error {
	key := fmt.Sprintf(DialogConnectionsKey, dialogID)
	return c.redis.SRem(ctx, c.formatKey(key), connectionID)
}

// GetDialogConnections gets all connections in a dialog room
func (c *CacheManager) GetDialogConnections(ctx context.Context, dialogID string) ([]string, error) {
	key := fmt.Sprintf(DialogConnectionsKey, dialogID)
	return c.redis.SMembers(ctx, c.formatKey(key))
}

// SetTypingIndicator sets typing indicator for user in dialog
func (c *CacheManager) SetTypingIndicator(ctx context.Context, dialogID, userID string, isTyping bool) error {
	key := fmt.Sprintf(TypingIndicatorsKey, dialogID)

	if isTyping {
		// Add user to typing set with TTL
		if err := c.redis.SAdd(ctx, c.formatKey(key), userID); err != nil {
			return err
		}
		return c.redis.Expire(ctx, c.formatKey(key), 30*time.Second)
	} else {
		// Remove user from typing set
		return c.redis.SRem(ctx, c.formatKey(key), userID)
	}
}

// GetTypingUsers gets users currently typing in dialog
func (c *CacheManager) GetTypingUsers(ctx context.Context, dialogID string) ([]string, error) {
	key := fmt.Sprintf(TypingIndicatorsKey, dialogID)
	return c.redis.SMembers(ctx, c.formatKey(key))
}

// Rate limiting operations

// CheckRateLimit checks if an action is rate limited
func (c *CacheManager) CheckRateLimit(ctx context.Context, limitType, identifier string, maxAttempts int, window time.Duration) (bool, error) {
	key := fmt.Sprintf(RateLimitKey, limitType, identifier)

	current, err := c.redis.Incr(ctx, c.formatKey(key))
	if err != nil {
		return false, err
	}

	if current == 1 {
		// First request, set TTL
		if err := c.redis.Expire(ctx, c.formatKey(key), window); err != nil {
			log.Printf("Failed to set TTL for rate limit key %s: %v", key, err)
		}
	}

	return current > int64(maxAttempts), nil
}

// OTP caching operations

// SetOTP caches OTP code
func (c *CacheManager) SetOTP(ctx context.Context, phone, purpose, otpHash string, ttl time.Duration) error {
	key := fmt.Sprintf(OTPKey, phone, purpose)
	return c.redis.Set(ctx, c.formatKey(key), otpHash, ttl)
}

// GetOTP retrieves cached OTP
func (c *CacheManager) GetOTP(ctx context.Context, phone, purpose string) (string, error) {
	key := fmt.Sprintf(OTPKey, phone, purpose)
	return c.redis.Get(ctx, c.formatKey(key))
}

// DeleteOTP removes OTP from cache
func (c *CacheManager) DeleteOTP(ctx context.Context, phone, purpose string) error {
	key := fmt.Sprintf(OTPKey, phone, purpose)
	return c.Delete(ctx, key)
}

// Payment caching operations

// SetWalletBalance caches wallet balance
func (c *CacheManager) SetWalletBalance(ctx context.Context, walletID string, balance interface{}) error {
	key := fmt.Sprintf(WalletBalanceKey, walletID)
	return c.Set(ctx, key, balance, MediumTTL)
}

// GetWalletBalance retrieves cached wallet balance
func (c *CacheManager) GetWalletBalance(ctx context.Context, walletID string, dest interface{}) error {
	key := fmt.Sprintf(WalletBalanceKey, walletID)
	return c.Get(ctx, key, dest)
}

// InvalidateWalletBalance removes wallet balance from cache
func (c *CacheManager) InvalidateWalletBalance(ctx context.Context, walletID string) error {
	key := fmt.Sprintf(WalletBalanceKey, walletID)
	return c.Delete(ctx, key)
}

// Message caching operations

// SetMessage caches message data
func (c *CacheManager) SetMessage(ctx context.Context, messageID string, message interface{}) error {
	key := fmt.Sprintf(MessageKey, messageID)
	return c.Set(ctx, key, message, LongTTL)
}

// GetMessage retrieves cached message
func (c *CacheManager) GetMessage(ctx context.Context, messageID string, dest interface{}) error {
	key := fmt.Sprintf(MessageKey, messageID)
	return c.Get(ctx, key, dest)
}

// Search caching operations

// SetSearchResults caches search results
func (c *CacheManager) SetSearchResults(ctx context.Context, searchType, queryHash string, results interface{}) error {
	key := fmt.Sprintf(SearchResultsKey, searchType, queryHash)
	return c.Set(ctx, key, results, MediumTTL)
}

// GetSearchResults retrieves cached search results
func (c *CacheManager) GetSearchResults(ctx context.Context, searchType, queryHash string, dest interface{}) error {
	key := fmt.Sprintf(SearchResultsKey, searchType, queryHash)
	return c.Get(ctx, key, dest)
}

// Batch operations

// InvalidateUserCache invalidates all cache entries for a user
func (c *CacheManager) InvalidateUserCache(ctx context.Context, userID string) error {
	keys := []string{
		fmt.Sprintf(UserProfileKey, userID),
		fmt.Sprintf(UserSettingsKey, userID),
		fmt.Sprintf(UserOnlineKey, userID),
		fmt.Sprintf(UserConnectionsKey, userID),
		fmt.Sprintf(PresenceKey, userID),
	}

	return c.Delete(ctx, keys...)
}

// InvalidateDialogCache invalidates all cache entries for a dialog
func (c *CacheManager) InvalidateDialogCache(ctx context.Context, dialogID string) error {
	keys := []string{
		fmt.Sprintf(DialogKey, dialogID),
		fmt.Sprintf(DialogParticipants, dialogID),
		fmt.Sprintf(DialogMessages, dialogID),
		fmt.Sprintf(DialogLastMessage, dialogID),
		fmt.Sprintf(DialogConnectionsKey, dialogID),
		fmt.Sprintf(TypingIndicatorsKey, dialogID),
	}

	return c.Delete(ctx, keys...)
}

// Health check
func (c *CacheManager) HealthCheck(ctx context.Context) error {
	return c.redis.HealthCheck(ctx)
}