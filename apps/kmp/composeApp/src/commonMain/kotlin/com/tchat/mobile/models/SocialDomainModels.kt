package com.tchat.mobile.models

import kotlinx.serialization.Serializable

/**
 * Domain models for social features following SQL sync architecture
 * These models represent the domain entities and work with SQLDelight
 */

// Enums for type safety
enum class InteractionType {
    LIKE, BOOKMARK, FOLLOW, ATTEND, VIEW, SHARE
}

enum class InteractionTargetType {
    POST, STORY, EVENT, USER, COMMENT
}

enum class FriendshipStatus {
    PENDING, ACCEPTED, BLOCKED, REJECTED
}

enum class EventAttendanceStatus {
    ATTENDING, MAYBE, NOT_ATTENDING
}

// Note: Using SocialPostType and PostSource from SocialTypes.kt to avoid conflicts

// Core domain models
@Serializable
data class Story(
    val id: String,
    val authorId: String,
    val content: String,
    val preview: String = "",
    val createdAt: Long,
    val expiresAt: Long,
    val isLive: Boolean = false,
    val viewCount: Int = 0,
    // UI state (not persisted)
    val isViewed: Boolean = false,
    val totalViews: Int = 0,
    val author: SocialUserProfile? = null
) {
    val isExpired: Boolean get() = System.currentTimeMillis() > expiresAt
    val remainingTime: Long get() = maxOf(0, expiresAt - System.currentTimeMillis())
}

@Serializable
data class SocialUserProfile(
    val userId: String,
    val displayName: String,
    val username: String,
    val avatarUrl: String = "",
    val bio: String = "",
    val isVerified: Boolean = false,
    val isOnline: Boolean = false,
    val lastSeen: Long? = null,
    val statusMessage: String = "",
    val createdAt: Long,
    val updatedAt: Long
) {
    val isActive: Boolean get() = isOnline || (lastSeen?.let { System.currentTimeMillis() - it < 300000 } ?: false) // 5 minutes
}

@Serializable
data class Friend(
    val id: String,
    val userId: String,
    val friendUserId: String,
    val status: FriendshipStatus,
    val createdAt: Long,
    val updatedAt: Long,
    // Populated from joins
    val profile: SocialUserProfile? = null,
    val mutualFriendsCount: Int = 0
) {
    val isAccepted: Boolean get() = status == FriendshipStatus.ACCEPTED
}

@Serializable
data class Event(
    val id: String,
    val title: String,
    val description: String,
    val eventDate: Long,
    val location: String,
    val price: String = "Free",
    val imageUrl: String = "",
    val category: String,
    val organizerId: String,
    val attendeesCount: Int = 0,
    val maxAttendees: Int? = null,
    val isPublic: Boolean = true,
    val createdAt: Long,
    val updatedAt: Long,
    // UI state (not persisted)
    val organizerProfile: SocialUserProfile? = null,
    val userAttendanceStatus: EventAttendanceStatus? = null,
    val actualAttendeesCount: Int = 0
) {
    val isUpcoming: Boolean get() = eventDate > System.currentTimeMillis()
    val isFull: Boolean get() = maxAttendees?.let { attendeesCount >= it } ?: false
    val daysUntilEvent: Long get() = (eventDate - System.currentTimeMillis()) / (24 * 60 * 60 * 1000)
}

@Serializable
data class SocialInteraction(
    val id: String,
    val userId: String,
    val targetId: String,
    val targetType: InteractionTargetType,
    val interactionType: InteractionType,
    val createdAt: Long,
    val updatedAt: Long,
    // Populated from joins
    val userProfile: SocialUserProfile? = null
)

@Serializable
data class SocialComment(
    val id: String,
    val targetId: String,
    val targetType: InteractionTargetType,
    val userId: String,
    val content: String,
    val parentCommentId: String? = null,
    val likesCount: Int = 0,
    val repliesCount: Int = 0,
    val createdAt: Long,
    val updatedAt: Long,
    // UI state (not persisted)
    val userProfile: SocialUserProfile? = null,
    val isLikedByUser: Boolean = false,
    val replies: List<SocialComment> = emptyList()
) {
    val isReply: Boolean get() = parentCommentId != null
}

// Data transfer objects for UI
@Serializable
data class InteractionCounts(
    val likes: Int = 0,
    val bookmarks: Int = 0,
    val shares: Int = 0,
    val views: Int = 0,
    val comments: Int = 0
)

@Serializable
data class UserInteractionState(
    val isLiked: Boolean = false,
    val isBookmarked: Boolean = false,
    val isFollowing: Boolean = false,
    val isAttending: Boolean = false
)

@Serializable
data class UserStats(
    val followersCount: Int = 0,
    val followingCount: Int = 0,
    val totalLikesReceived: Int = 0,
    val postsCount: Int = 0,
    val storiesCount: Int = 0
)

// Note: Legacy UI models (StoryItem, FriendItem, etc.) are defined in SocialTypes.kt
// Conversion logic is handled in the service layer to avoid conflicts

// Domain model to UI model conversion logic moved to service layer

// Utility functions for formatting
private fun formatTimestamp(timestamp: Long): String {
    val now = System.currentTimeMillis()
    val diff = now - timestamp
    return when {
        diff < 60000 -> "just now"
        diff < 3600000 -> "${diff / 60000}m ago"
        diff < 86400000 -> "${diff / 3600000}h ago"
        diff < 604800000 -> "${diff / 86400000}d ago"
        else -> "${diff / 604800000}w ago"
    }
}

private fun formatTimeRemaining(milliseconds: Long): String {
    val hours = milliseconds / (60 * 60 * 1000)
    return when {
        hours <= 0 -> "Expired"
        hours < 24 -> "${hours}h remaining"
        else -> "${hours / 24}d remaining"
    }
}

private fun formatEventDate(timestamp: Long): String {
    // This should use proper date formatting based on platform
    return formatTimestamp(timestamp)
}