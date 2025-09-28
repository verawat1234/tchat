package com.tchat.mobile.social.domain.models

import kotlinx.datetime.Instant
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.JsonElement

/**
 * KMP Social Domain Models
 *
 * Cross-platform data models for social features with:
 * - Southeast Asian regional optimization
 * - Mobile-first design patterns
 * - Offline-first sync capabilities
 * - Type-safe serialization
 */

@Serializable
data class SocialProfile(
    val id: String,
    val username: String,
    val displayName: String?,
    val bio: String?,
    val avatar: String?,
    val interests: List<String> = emptyList(),
    val socialLinks: Map<String, JsonElement>? = null,
    val socialPreferences: Map<String, JsonElement>? = null,
    val followersCount: Int = 0,
    val followingCount: Int = 0,
    val postsCount: Int = 0,
    val isSocialVerified: Boolean = false,
    val country: String = "TH",
    val region: String = "TH",
    val socialCreatedAt: String,
    val socialUpdatedAt: String,
    // Mobile optimization
    val lastSyncAt: String? = null,
    val syncVersion: String = "1.0"
)

@Serializable
data class SocialPost(
    val id: String,
    val authorId: String,
    val authorUsername: String,
    val authorDisplayName: String?,
    val authorAvatar: String?,
    val content: String,
    val contentType: String = "text", // text, image, video, link
    val mediaUrls: List<String> = emptyList(),
    val thumbnailUrl: String? = null,
    val tags: List<String> = emptyList(),
    val mentions: List<String> = emptyList(),
    val visibility: String = "public", // public, followers, private
    val likesCount: Int = 0,
    val commentsCount: Int = 0,
    val sharesCount: Int = 0,
    val viewsCount: Int = 0,
    val isLikedByUser: Boolean = false,
    val isBookmarkedByUser: Boolean = false,
    val createdAt: String,
    val updatedAt: String,
    // Southeast Asian features
    val language: String = "en",
    val region: String = "TH",
    val isRegionalTrending: Boolean = false,
    // Mobile optimization
    val lastSyncAt: String? = null,
    val isOfflineEdit: Boolean = false
)

@Serializable
data class SocialComment(
    val id: String,
    val postId: String,
    val authorId: String,
    val authorUsername: String,
    val authorDisplayName: String?,
    val authorAvatar: String?,
    val content: String,
    val parentCommentId: String? = null,
    val likesCount: Int = 0,
    val repliesCount: Int = 0,
    val isLikedByUser: Boolean = false,
    val createdAt: String,
    val updatedAt: String,
    // Mobile optimization
    val lastSyncAt: String? = null,
    val isOfflineEdit: Boolean = false
)

@Serializable
data class SocialInteraction(
    val id: String,
    val userId: String,
    val targetId: String,
    val targetType: String, // post, comment, user, story
    val interactionType: String, // like, bookmark, follow, share, view
    val createdAt: String,
    val updatedAt: String,
    // Mobile optimization
    val lastSyncAt: String? = null,
    val isOfflineAction: Boolean = false
)

@Serializable
data class SocialFeed(
    val posts: List<SocialPost>,
    val hasMore: Boolean,
    val nextOffset: Int? = null,
    val lastSyncAt: String,
    val region: String,
    val feedType: String // home, discover, trending, following
)

@Serializable
data class DiscoveryProfile(
    val profile: SocialProfile,
    val mutualFriends: Int = 0,
    val commonInterests: List<String> = emptyList(),
    val discoveryReason: String, // region, interests, mutual_friends, trending
    val score: Double = 0.0
)

@Serializable
data class RegionalTrending(
    val posts: List<SocialPost>,
    val hashtags: List<TrendingHashtag>,
    val region: String,
    val timeWindow: String, // today, week, month
    val lastUpdatedAt: String
)

@Serializable
data class TrendingHashtag(
    val tag: String,
    val count: Int,
    val region: String,
    val growth: Double,
    val posts: List<SocialPost> = emptyList()
)

// Sync and conflict resolution models
@Serializable
data class SyncOperation(
    val id: String,
    val operation: String, // create, update, delete
    val resourceType: String, // post, comment, interaction, profile
    val resourceId: String,
    val data: JsonElement? = null,
    val timestamp: String,
    val status: String = "pending", // pending, synced, failed, conflict
    val conflictResolution: String? = null,
    val retryCount: Int = 0
)

@Serializable
data class SyncConflict(
    val id: String,
    val operationId: String,
    val resourceType: String,
    val resourceId: String,
    val clientData: JsonElement,
    val serverData: JsonElement,
    val conflictType: String, // version, deleted, field_conflict
    val suggestedResolution: String, // client_wins, server_wins, merge
    val createdAt: String
)

@Serializable
data class SyncResponse(
    val success: Boolean,
    val syncedOperations: List<String>,
    val failedOperations: List<String>,
    val conflicts: List<SyncConflict>,
    val serverTimestamp: String,
    val nextSyncTimestamp: String
)

// Mobile-specific request/response models
@Serializable
data class MobileSyncRequest(
    val userId: String,
    val lastSyncAt: String?,
    val operations: List<SyncOperation> = emptyList(),
    val region: String = "TH",
    val deviceInfo: MobileDeviceInfo
)

@Serializable
data class MobileDeviceInfo(
    val platform: String, // ios, android
    val version: String,
    val language: String = "en",
    val timezone: String = "Asia/Bangkok",
    val region: String = "TH"
)

@Serializable
data class InitialDataResponse(
    val profile: SocialProfile,
    val recentFeed: SocialFeed,
    val discoveryProfiles: List<DiscoveryProfile>,
    val regionalTrending: RegionalTrending,
    val unreadNotifications: Int,
    val lastSyncAt: String
)

// Social action request models
@Serializable
data class CreatePostRequest(
    val content: String,
    val contentType: String = "text",
    val mediaUrls: List<String> = emptyList(),
    val tags: List<String> = emptyList(),
    val mentions: List<String> = emptyList(),
    val visibility: String = "public",
    val language: String = "en",
    val region: String = "TH"
)

@Serializable
data class CreateCommentRequest(
    val postId: String,
    val content: String,
    val parentCommentId: String? = null
)

@Serializable
data class InteractionRequest(
    val targetId: String,
    val targetType: String,
    val interactionType: String
)

@Serializable
data class UpdateProfileRequest(
    val displayName: String? = null,
    val bio: String? = null,
    val interests: List<String>? = null,
    val socialLinks: Map<String, JsonElement>? = null
)

// Social search and discovery models
@Serializable
data class SocialSearchRequest(
    val query: String,
    val type: String = "all", // all, posts, users, hashtags
    val region: String = "TH",
    val limit: Int = 20,
    val offset: Int = 0
)

@Serializable
data class SocialSearchResponse(
    val posts: List<SocialPost> = emptyList(),
    val users: List<SocialProfile> = emptyList(),
    val hashtags: List<TrendingHashtag> = emptyList(),
    val hasMore: Boolean,
    val total: Int
)

// Southeast Asian localization models
@Serializable
data class LocalizedContent(
    val originalText: String,
    val translations: Map<String, String> = emptyMap(),
    val detectedLanguage: String,
    val region: String
)

@Serializable
data class RegionalSettings(
    val region: String,
    val language: String,
    val timezone: String,
    val dateFormat: String,
    val culturalPreferences: Map<String, JsonElement> = emptyMap()
)

// Error handling models
@Serializable
data class SocialError(
    val code: String,
    val message: String,
    val details: Map<String, JsonElement>? = null,
    val timestamp: String
)