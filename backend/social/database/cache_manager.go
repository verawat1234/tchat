package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"tchat/social/models"
)

// CacheManager provides intelligent caching for social media data
type CacheManager struct {
	db          *gorm.DB
	cache       map[string]CacheEntry
	hitRate     float64
	totalHits   int64
	totalMisses int64
}

// CacheEntry represents a cached item with expiration
type CacheEntry struct {
	Data      interface{} `json:"data"`
	ExpiresAt time.Time   `json:"expires_at"`
	CreatedAt time.Time   `json:"created_at"`
	HitCount  int64       `json:"hit_count"`
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	DefaultTTL         time.Duration
	ProfileCacheTTL    time.Duration
	PostCacheTTL       time.Duration
	TrendingCacheTTL   time.Duration
	MaxCacheSize       int
	CleanupInterval    time.Duration
	PrefetchEnabled    bool
	CompressionEnabled bool
}

// NewCacheManager creates a new cache manager with optimized defaults
func NewCacheManager(db *gorm.DB) *CacheManager {
	cm := &CacheManager{
		db:    db,
		cache: make(map[string]CacheEntry),
	}

	// Start cleanup routine
	go cm.startCleanupRoutine()

	return cm
}

// GetDefaultCacheConfig returns optimized cache configuration for social media
func GetDefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		DefaultTTL:         time.Minute * 15,
		ProfileCacheTTL:    time.Hour * 2,     // User profiles change less frequently
		PostCacheTTL:       time.Minute * 30,  // Post data needs fresh metrics
		TrendingCacheTTL:   time.Minute * 5,   // Trending data needs frequent updates
		MaxCacheSize:       10000,             // Maximum cached items
		CleanupInterval:    time.Minute * 10,  // Cleanup expired items
		PrefetchEnabled:    true,              // Prefetch popular content
		CompressionEnabled: true,              // Compress large cached objects
	}
}

// CacheUserProfile caches user profile with intelligent TTL
func (cm *CacheManager) CacheUserProfile(userID string, profile *models.SocialProfile, ttl time.Duration) {
	key := fmt.Sprintf("profile:%s", userID)
	cm.setCache(key, profile, ttl)
}

// GetCachedUserProfile retrieves cached user profile
func (cm *CacheManager) GetCachedUserProfile(userID string) (*models.SocialProfile, bool) {
	key := fmt.Sprintf("profile:%s", userID)
	entry, exists := cm.getCache(key)
	if !exists {
		return nil, false
	}

	if profile, ok := entry.Data.(*models.SocialProfile); ok {
		return profile, true
	}

	// Type assertion failed, remove invalid entry
	cm.invalidateCache(key)
	return nil, false
}

// CachePost caches post data with metrics
func (cm *CacheManager) CachePost(postID string, post *models.Post, ttl time.Duration) {
	key := fmt.Sprintf("post:%s", postID)
	cm.setCache(key, post, ttl)
}

// GetCachedPost retrieves cached post
func (cm *CacheManager) GetCachedPost(postID string) (*models.Post, bool) {
	key := fmt.Sprintf("post:%s", postID)
	entry, exists := cm.getCache(key)
	if !exists {
		return nil, false
	}

	if post, ok := entry.Data.(*models.Post); ok {
		return post, true
	}

	cm.invalidateCache(key)
	return nil, false
}

// CacheTrendingPosts caches trending posts list
func (cm *CacheManager) CacheTrendingPosts(posts []*models.Post, ttl time.Duration) {
	key := "trending:posts"
	cm.setCache(key, posts, ttl)
}

// GetCachedTrendingPosts retrieves cached trending posts
func (cm *CacheManager) GetCachedTrendingPosts() ([]*models.Post, bool) {
	key := "trending:posts"
	entry, exists := cm.getCache(key)
	if !exists {
		return nil, false
	}

	if posts, ok := entry.Data.([]*models.Post); ok {
		return posts, true
	}

	cm.invalidateCache(key)
	return nil, false
}

// CacheUserFeed caches personalized user feed
func (cm *CacheManager) CacheUserFeed(userID string, posts []*models.Post, ttl time.Duration) {
	key := fmt.Sprintf("feed:%s", userID)
	cm.setCache(key, posts, ttl)
}

// GetCachedUserFeed retrieves cached user feed
func (cm *CacheManager) GetCachedUserFeed(userID string) ([]*models.Post, bool) {
	key := fmt.Sprintf("feed:%s", userID)
	entry, exists := cm.getCache(key)
	if !exists {
		return nil, false
	}

	if posts, ok := entry.Data.([]*models.Post); ok {
		return posts, true
	}

	cm.invalidateCache(key)
	return nil, false
}

// CacheEngagementMetrics caches engagement analytics
func (cm *CacheManager) CacheEngagementMetrics(key string, metrics interface{}, ttl time.Duration) {
	cacheKey := fmt.Sprintf("metrics:%s", key)
	cm.setCache(cacheKey, metrics, ttl)
}

// GetCachedEngagementMetrics retrieves cached engagement metrics
func (cm *CacheManager) GetCachedEngagementMetrics(key string) (interface{}, bool) {
	cacheKey := fmt.Sprintf("metrics:%s", key)
	entry, exists := cm.getCache(cacheKey)
	if !exists {
		return nil, false
	}

	return entry.Data, true
}

// InvalidateUserCache invalidates all cache entries for a user
func (cm *CacheManager) InvalidateUserCache(userID string) {
	patterns := []string{
		fmt.Sprintf("profile:%s", userID),
		fmt.Sprintf("feed:%s", userID),
		fmt.Sprintf("posts:author:%s", userID),
	}

	for _, pattern := range patterns {
		cm.invalidateCache(pattern)
	}

	// Also invalidate trending if user has trending posts
	cm.invalidatePattern("trending:")
}

// InvalidatePostCache invalidates post-related cache entries
func (cm *CacheManager) InvalidatePostCache(postID string) {
	patterns := []string{
		fmt.Sprintf("post:%s", postID),
		"trending:posts",
		"popular:posts",
	}

	for _, pattern := range patterns {
		cm.invalidateCache(pattern)
	}

	// Invalidate feeds that might contain this post
	cm.invalidatePattern("feed:")
}

// PrefetchPopularContent preloads frequently accessed content
func (cm *CacheManager) PrefetchPopularContent(ctx context.Context) error {
	// Prefetch trending posts
	var trendingPosts []*models.Post
	if err := cm.db.WithContext(ctx).
		Where("is_trending = true AND created_at > ?", time.Now().Add(-24*time.Hour)).
		Order("likes_count DESC, comments_count DESC").
		Limit(50).
		Find(&trendingPosts).Error; err != nil {
		return fmt.Errorf("failed to prefetch trending posts: %w", err)
	}

	cm.CacheTrendingPosts(trendingPosts, time.Minute*5)

	// Prefetch popular profiles
	var popularProfiles []*models.SocialProfile
	if err := cm.db.WithContext(ctx).
		Where("is_social_verified = true OR followers_count > 1000").
		Order("followers_count DESC").
		Limit(100).
		Find(&popularProfiles).Error; err != nil {
		return fmt.Errorf("failed to prefetch popular profiles: %w", err)
	}

	for _, profile := range popularProfiles {
		cm.CacheUserProfile(profile.ID.String(), profile, time.Hour*2)
	}

	return nil
}

// WarmupCache preloads essential data for better performance
func (cm *CacheManager) WarmupCache(ctx context.Context) error {
	// Prefetch popular content
	if err := cm.PrefetchPopularContent(ctx); err != nil {
		return fmt.Errorf("failed to prefetch popular content: %w", err)
	}

	// Preload community data
	var communities []*models.Community
	if err := cm.db.WithContext(ctx).
		Where("type = 'public' AND members_count > 100").
		Order("members_count DESC").
		Limit(50).
		Find(&communities).Error; err != nil {
		return fmt.Errorf("failed to prefetch communities: %w", err)
	}

	for _, community := range communities {
		key := fmt.Sprintf("community:%s", community.ID.String())
		cm.setCache(key, community, time.Hour)
	}

	return nil
}

// GetCacheStats returns cache performance statistics
func (cm *CacheManager) GetCacheStats() map[string]interface{} {
	totalRequests := cm.totalHits + cm.totalMisses
	hitRate := float64(0)
	if totalRequests > 0 {
		hitRate = float64(cm.totalHits) / float64(totalRequests) * 100
	}

	return map[string]interface{}{
		"total_entries":   len(cm.cache),
		"total_hits":      cm.totalHits,
		"total_misses":    cm.totalMisses,
		"hit_rate":        fmt.Sprintf("%.2f%%", hitRate),
		"total_requests":  totalRequests,
	}
}

// Private helper methods

// setCache stores an item in cache with TTL
func (cm *CacheManager) setCache(key string, data interface{}, ttl time.Duration) {
	entry := CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
		HitCount:  0,
	}

	cm.cache[key] = entry

	// Enforce cache size limit
	if len(cm.cache) > 10000 { // Max cache size
		cm.evictOldestEntries(1000) // Evict 1000 oldest entries
	}
}

// getCache retrieves an item from cache
func (cm *CacheManager) getCache(key string) (CacheEntry, bool) {
	entry, exists := cm.cache[key]
	if !exists {
		cm.totalMisses++
		return CacheEntry{}, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		delete(cm.cache, key)
		cm.totalMisses++
		return CacheEntry{}, false
	}

	// Update hit count
	entry.HitCount++
	cm.cache[key] = entry
	cm.totalHits++

	return entry, true
}

// invalidateCache removes a specific cache entry
func (cm *CacheManager) invalidateCache(key string) {
	delete(cm.cache, key)
}

// invalidatePattern removes all cache entries matching a pattern
func (cm *CacheManager) invalidatePattern(pattern string) {
	for key := range cm.cache {
		if len(key) >= len(pattern) && key[:len(pattern)] == pattern {
			delete(cm.cache, key)
		}
	}
}

// evictOldestEntries removes the oldest cache entries
func (cm *CacheManager) evictOldestEntries(count int) {
	type keyTime struct {
		key  string
		time time.Time
	}

	var entries []keyTime
	for key, entry := range cm.cache {
		entries = append(entries, keyTime{key: key, time: entry.CreatedAt})
	}

	// Sort by creation time (oldest first)
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].time.After(entries[j].time) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	// Remove oldest entries
	for i := 0; i < count && i < len(entries); i++ {
		delete(cm.cache, entries[i].key)
	}
}

// startCleanupRoutine runs periodic cache cleanup
func (cm *CacheManager) startCleanupRoutine() {
	ticker := time.NewTicker(time.Minute * 10) // Cleanup every 10 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cm.cleanupExpiredEntries()
		}
	}
}

// cleanupExpiredEntries removes expired cache entries
func (cm *CacheManager) cleanupExpiredEntries() {
	now := time.Now()
	expiredKeys := make([]string, 0)

	for key, entry := range cm.cache {
		if now.After(entry.ExpiresAt) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	for _, key := range expiredKeys {
		delete(cm.cache, key)
	}

	if len(expiredKeys) > 0 {
		fmt.Printf("Cleaned up %d expired cache entries\n", len(expiredKeys))
	}
}