package com.tchat.mobile.models

import kotlinx.serialization.Serializable

/**
 * Review Type System - Similar to MessageType
 *
 * Defines different types of reviews that users can create:
 * - TEXT: Text-only reviews with ratings and hashtags
 * - IMAGE: Photo reviews with multiple images
 * - VIDEO: Video reviews with TikTok-style media
 * - MIXED: Combined media (text + images + videos)
 */
enum class ReviewType {
    TEXT,       // Text-only review with rating and hashtags
    IMAGE,      // Photo review with image gallery
    VIDEO,      // Video review with video content
    MIXED,      // Combined media review (text + images + videos)
    QUICK,      // Quick rating with minimal text
    DETAILED    // Comprehensive review with all elements
}

/**
 * Review Rating Type - Different rating systems
 */
enum class ReviewRatingType {
    STARS_5,    // 1-5 star rating
    STARS_10,   // 1-10 star rating
    THUMBS,     // Thumbs up/down
    EMOJI,      // Emoji-based rating (😍, 😊, 😐, 😞, 😠)
    PERCENTAGE  // Percentage-based rating (0-100%)
}

/**
 * Review Content Data Classes
 */
@Serializable
data class ReviewContent(
    val type: ReviewType,
    val text: String? = null,
    val images: List<ReviewImage> = emptyList(),
    val videos: List<ReviewVideo> = emptyList(),
    val hashtags: List<String> = emptyList(),
    val mentions: List<String> = emptyList()
)

@Serializable
data class ReviewImage(
    val id: String,
    val url: String,
    val caption: String? = null,
    val aspectRatio: Float = 1f
)

@Serializable
data class ReviewVideo(
    val id: String,
    val url: String,
    val thumbnailUrl: String? = null,
    val duration: String,
    val caption: String? = null
)

@Serializable
data class ReviewRating(
    val type: ReviewRatingType,
    val value: Float, // Normalized to 0-1 scale
    val displayValue: String, // Human-readable display (e.g., "4/5", "👍", "85%")
    val categories: Map<String, Float> = emptyMap() // Optional category ratings
)

@Serializable
data class Review(
    val id: String,
    val userId: String,
    val userName: String,
    val userAvatar: String? = null,
    val targetType: ReviewTargetType,
    val targetId: String,
    val targetName: String,
    val rating: ReviewRating,
    val content: ReviewContent,
    val isVerifiedPurchase: Boolean = false,
    val isHelpful: Int = 0,
    val isReported: Boolean = false,
    val likes: Int = 0,
    val dislikes: Int = 0,
    val comments: Int = 0,
    val shares: Int = 0,
    val isLiked: Boolean = false,
    val isBookmarked: Boolean = false,
    val createdAt: String,
    val updatedAt: String? = null
)

enum class ReviewTargetType {
    PRODUCT,
    SHOP,
    SERVICE,
    USER,
    EXPERIENCE
}

/**
 * Review Statistics for Analytics
 */
@Serializable
data class ReviewStats(
    val totalReviews: Int,
    val averageRating: Float,
    val ratingDistribution: Map<Int, Int>, // Star level -> count
    val typeDistribution: Map<ReviewType, Int>,
    val verifiedPurchaseCount: Int,
    val helpfulCount: Int,
    val recentReviewsCount: Int // Last 30 days
)

/**
 * Review Filter Options
 */
data class ReviewFilter(
    val type: ReviewType? = null,
    val minRating: Float? = null,
    val maxRating: Float? = null,
    val verifiedOnly: Boolean = false,
    val withImages: Boolean = false,
    val withVideos: Boolean = false,
    val timeRange: ReviewTimeRange = ReviewTimeRange.ALL,
    val sortBy: ReviewSortBy = ReviewSortBy.NEWEST
)

enum class ReviewTimeRange {
    LAST_WEEK,
    LAST_MONTH,
    LAST_3_MONTHS,
    LAST_YEAR,
    ALL
}

enum class ReviewSortBy {
    NEWEST,
    OLDEST,
    HIGHEST_RATING,
    LOWEST_RATING,
    MOST_HELPFUL,
    MOST_LIKES
}