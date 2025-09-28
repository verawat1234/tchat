package com.tchat.mobile.models

import kotlinx.serialization.Serializable

/**
 * Generic Post System - Social Media Platform Architecture
 *
 * This system treats everything as posts (reviews, social posts, videos, etc.)
 * Similar to how Instagram, TikTok, and other platforms work
 */

/**
 * Unified Post Type System - 42 Post Types
 * Aligned with web platform PostData.ts for consistent cross-platform experience
 */
enum class PostType {
    // Core Content Types (8)
    TEXT,                    // Simple status updates
    IMAGE,                   // Single/multiple photos
    VIDEO,                   // Video posts
    AUDIO,                   // Voice notes, music
    LINK_SHARE,              // Shared articles/websites
    POST_MESSAGE,            // Message posted to someone's timeline/wall
    REVIEW,                  // Reviews of places, products, services
    ALBUM,                   // Photo collections

    // Rich Media Types (6)
    STORY,                   // Ephemeral 24h content
    REEL,                    // Short-form vertical video
    LIVE_STREAM,             // Live video broadcasts
    PLAYLIST,                // Music/video collections
    MOOD_BOARD,              // Visual inspiration
    TUTORIAL,                // How-to content

    // Interactive Content (6)
    POLL,                    // Voting posts
    QUIZ,                    // Trivia, personality tests
    SURVEY,                  // Feedback collection
    Q_AND_A,                 // Ask me anything
    CHALLENGE,               // Viral challenges/trends
    PETITION,                // Social causes

    // Social & Location (8)
    CHECK_IN,                // Location-based posts
    TRAVEL_LOG,              // Trip updates/itinerary
    LIFE_EVENT,              // Major life moments
    MILESTONE,               // Personal achievements
    MEMORY,                  // Throwback/flashback posts
    ANNIVERSARY,             // Yearly memories
    RECOMMENDATION,          // Place/product suggestions
    GROUP_ACTIVITY,          // Group-specific content

    // Commercial & Business (6)
    PRODUCT_SHOWCASE,        // Selling items
    SERVICE_PROMOTION,       // Business services
    EVENT_PROMOTION,         // Events/meetups
    JOB_POSTING,            // Hiring/career opportunities
    FUNDRAISER,             // Charity/personal causes
    COLLABORATION,          // Creative projects

    // Specialized Content (8)
    RECIPE,                 // Cooking/food content
    WORKOUT,                // Fitness routines
    BOOK_REVIEW,            // Reading updates
    MOOD_UPDATE,            // Emotional status
    ACHIEVEMENT,            // Gaming/app achievements
    QUOTE,                  // Inspirational quotes
    MUSIC,                  // Music sharing/streaming
    VENUE                   // Venue information/reviews
}

enum class PostContentType {
    TEXT,
    IMAGE,
    VIDEO,
    MIXED,
    LIVE,
    POLL
}

@Serializable
data class PostContent(
    val type: PostContentType,
    val text: String? = null,
    val images: List<PostImage> = emptyList(),
    val videos: List<PostVideo> = emptyList(),
    val hashtags: List<String> = emptyList(),
    val mentions: List<String> = emptyList(),
    val location: String? = null,
    val poll: PostPoll? = null
)

@Serializable
data class PostImage(
    val id: String,
    val url: String,
    val caption: String? = null,
    val aspectRatio: Float = 1f,
    val filters: List<String> = emptyList()
)

@Serializable
data class PostVideo(
    val id: String,
    val url: String,
    val thumbnailUrl: String? = null,
    val duration: String,
    val caption: String? = null,
    val isAutoPlay: Boolean = true
)

@Serializable
data class PostPoll(
    val question: String,
    val options: List<String>,
    val votes: Map<Int, Int> = emptyMap(), // option index -> vote count
    val expiresAt: String? = null
)


@Serializable
data class PostReaction(
    val type: ReactionType,
    val userId: String,
    val timestamp: String,
    val userName: String? = null
)

enum class ShareType {
    DIRECT_SHARE,    // Simple reshare
    QUOTE_SHARE,     // Share with comment
    STORY_SHARE,     // Share to story
    MESSAGE_SHARE,   // Share via DM
    EXTERNAL_SHARE   // Share outside platform
}

@Serializable
data class PostShare(
    val id: String,
    val userId: String,
    val userName: String,
    val timestamp: String,
    val shareType: ShareType,
    val addedComment: String? = null,
    val sharedToGroups: List<String> = emptyList()
)

@Serializable
data class PostComment(
    val id: String,
    val userId: String,
    val userName: String,
    val userAvatar: String? = null,
    val content: String,
    val timestamp: String,
    val replies: List<PostComment> = emptyList(),
    val reactions: List<PostReaction> = emptyList(),
    val isEdited: Boolean = false,
    val isDeleted: Boolean = false,
    val mentionedUsers: List<String> = emptyList()
)

@Serializable
data class PostInteractions(
    val reactions: List<PostReaction> = emptyList(),
    val comments: List<PostComment> = emptyList(),
    val shares: List<PostShare> = emptyList(),
    val saves: List<String> = emptyList(),        // User IDs who saved
    val views: Int = 0,
    val reach: Int = 0,                           // Unique users reached
    val impressions: Int = 0,                     // Total views
    val clickThroughs: Int = 0,                   // Link/action clicks
    val engagementRate: Float = 0f,               // Calculated engagement %

    // User-specific interaction states
    val isLiked: Boolean = false,
    val isBookmarked: Boolean = false
) {
    // Computed properties for legacy compatibility
    val likes: Int get() = reactions.count { it.type == ReactionType.LIKE }
    val commentsCount: Int get() = comments.size
    val sharesCount: Int get() = shares.size
    val savesCount: Int get() = saves.size
}

@Serializable
data class PostUser(
    val id: String,
    val username: String,
    val displayName: String? = null,
    val avatarUrl: String? = null,
    val isVerified: Boolean = false,
    val followerCount: Int = 0,
    val isFollowing: Boolean = false
)

@Serializable
data class PostMetadata(
    val targetType: String? = null,      // "product", "shop", "user", etc.
    val targetId: String? = null,        // ID of the target
    val targetName: String? = null,      // Name of the target
    val rating: Float? = null,           // For reviews
    val price: String? = null,           // For product posts
    val category: String? = null,        // Content category
    val tags: List<String> = emptyList(),
    val isSponsored: Boolean = false,
    val sponsorName: String? = null,
    val editedAt: String? = null,
    val originalPostId: String? = null,              // If this is a shared post
    val isPromoted: Boolean = false,                 // Sponsored/boosted content
    val mentionedUsers: List<String> = emptyList(),  // @mentions
    val location: PostLocation? = null,
    val mood: MoodType? = null,
    val feeling: FeelingType? = null,
    val activity: ActivityType? = null,
    val contentWarning: ContentWarningType? = null,
    val language: String? = null,
    val isArchived: Boolean = false,
    val archivedAt: String? = null,
    val expiresAt: String? = null,                   // For stories/temporary content
    val allowComments: Boolean = true,
    val allowShares: Boolean = true,
    val allowSaves: Boolean = true,
    val isPinned: Boolean = false,
    val isSticky: Boolean = false                    // Stays at top of profile/group
)

/**
 * Enhanced Post Data Class - Unified Social Platform Architecture
 * Aligned with web platform PostData interface for cross-platform consistency
 */
@Serializable
data class Post(
    val id: String,
    val type: PostType,
    val user: PostUser,
    val content: PostContent,
    val interactions: PostInteractions,
    val privacy: PostPrivacy = PostPrivacy.PUBLIC,
    val audience: PostAudience? = null,
    val metadata: PostMetadata? = null,
    val createdAt: String,
    val updatedAt: String? = null,
    val isEdited: Boolean = false,

    // Legacy field for backward compatibility
    val visibility: PostVisibility = PostVisibility.PUBLIC
)

/**
 * Enhanced Privacy Controls - 8 Privacy Levels
 * Aligned with web platform PostData.ts
 */
enum class PostPrivacy {
    PUBLIC,              // Visible to everyone
    FRIENDS,             // Friends only
    CLOSE_FRIENDS,       // Close friends list
    FOLLOWERS,           // Followers only
    MUTUAL_FRIENDS,      // Mutual connections
    CUSTOM,              // Custom audience
    UNLISTED,            // Hidden from feeds but shareable
    PRIVATE              // Only author can see
}

@Serializable
data class PostAudience(
    val includedUsers: List<String> = emptyList(),    // Specific users who can see
    val excludedUsers: List<String> = emptyList(),    // Users who cannot see
    val includedGroups: List<String> = emptyList(),   // Specific groups/circles
    val excludedGroups: List<String> = emptyList(),   // Excluded groups
    val locationRestriction: LocationRestriction? = null,
    val ageRange: AgeRange? = null,
    val interests: List<String> = emptyList()         // Interest-based targeting
)

@Serializable
data class LocationRestriction(
    val countries: List<String> = emptyList(),
    val regions: List<String> = emptyList(),
    val cities: List<String> = emptyList(),
    val radius: GeofenceRadius? = null
)

@Serializable
data class GeofenceRadius(
    val latitude: Double,
    val longitude: Double,
    val distance: Int                                 // Distance in kilometers
)

@Serializable
data class AgeRange(
    val min: Int? = null,
    val max: Int? = null
)

// Duplicate PostMetadata removed - using the enhanced first definition

enum class MoodType {
    HAPPY, EXCITED, GRATEFUL, LOVED, BLESSED, RELAXED, CONTENT, MOTIVATED,
    PROUD, ACCOMPLISHED, TIRED, STRESSED, SAD, ANXIOUS, CONFUSED,
    FRUSTRATED, ANGRY, LONELY, NOSTALGIC, CONTEMPLATIVE
}

enum class FeelingType {
    AMAZING, FANTASTIC, GOOD, OKAY, MEH, NOT_GREAT, TERRIBLE
}

enum class ActivityType {
    EATING, DRINKING, TRAVELING, EXERCISING, WORKING, STUDYING, READING,
    WATCHING, LISTENING, PLAYING, COOKING, SHOPPING, CELEBRATING,
    RELAXING, SLEEPING
}

enum class ContentWarningType {
    NONE, SENSITIVE_CONTENT, GRAPHIC_VIOLENCE, ADULT_CONTENT,
    DISTURBING_CONTENT, SPOILER, FLASHING_LIGHTS, POLITICAL_CONTENT
}

@Serializable
data class PostLocation(
    val name: String,
    val address: String? = null,
    val latitude: Double? = null,
    val longitude: Double? = null,
    val city: String? = null,
    val country: String? = null,
    val category: LocationCategory? = null
)

enum class LocationCategory {
    RESTAURANT, HOTEL, ATTRACTION, SHOPPING, ENTERTAINMENT, OUTDOORS,
    TRANSPORTATION, HOME, WORK, EDUCATION, HEALTHCARE, GOVERNMENT,
    RELIGIOUS, SPORTS, OTHER
}

// Legacy enum for backward compatibility
enum class PostVisibility {
    PUBLIC,
    FRIENDS,
    PRIVATE,
    UNLISTED
}

// Duplicate PostComment removed - using the enhanced first definition

@Serializable
data class PostHashtag(
    val tag: String,
    val count: Int,
    val isFollowing: Boolean = false,
    val category: String? = null
)

/**
 * Post Type Guards and Validation Utilities
 * Similar to web platform's TypeScript type guards
 */
object PostTypeValidator {

    fun isImagePost(post: Post): Boolean {
        return post.content.type == PostContentType.IMAGE && post.content.images.isNotEmpty()
    }

    fun isVideoPost(post: Post): Boolean {
        return (post.content.type == PostContentType.VIDEO && post.content.videos.isNotEmpty()) ||
               (post.type == PostType.VIDEO || post.type == PostType.REEL || post.type == PostType.LIVE_STREAM)
    }

    fun isPollPost(post: Post): Boolean {
        return post.content.poll != null || post.type == PostType.POLL
    }

    fun isStoryPost(post: Post): Boolean {
        return post.type == PostType.STORY && post.metadata?.expiresAt != null
    }

    fun isReviewPost(post: Post): Boolean {
        return post.type == PostType.REVIEW && post.metadata?.rating != null
    }

    fun isLocationPost(post: Post): Boolean {
        return post.type == PostType.CHECK_IN ||
               post.type == PostType.TRAVEL_LOG ||
               post.metadata?.location != null
    }

    fun isInteractivePost(post: Post): Boolean {
        return post.type in listOf(
            PostType.POLL, PostType.QUIZ, PostType.SURVEY,
            PostType.Q_AND_A, PostType.CHALLENGE, PostType.PETITION
        )
    }

    fun isCommercialPost(post: Post): Boolean {
        return post.type in listOf(
            PostType.PRODUCT_SHOWCASE, PostType.SERVICE_PROMOTION,
            PostType.EVENT_PROMOTION, PostType.JOB_POSTING, PostType.FUNDRAISER
        ) || post.metadata?.isPromoted == true
    }

    fun isEphemeralPost(post: Post): Boolean {
        return post.type == PostType.STORY || post.metadata?.expiresAt != null
    }

    fun getPostDominantContent(post: Post): PostContentType {
        return when {
            post.content.images.isNotEmpty() -> PostContentType.IMAGE
            post.content.videos.isNotEmpty() -> PostContentType.VIDEO
            post.content.poll != null -> PostContentType.POLL
            post.type == PostType.LIVE_STREAM -> PostContentType.LIVE
            else -> PostContentType.TEXT
        }
    }
}

/**
 * Post Validation Rules
 */
object PostValidationRules {

    fun validatePost(post: Post): List<String> {
        val errors = mutableListOf<String>()

        // Required field validation
        if (post.id.isEmpty()) {
            errors.add("Post ID is required")
        }

        if (post.user.id.isEmpty()) {
            errors.add("User ID is required")
        }

        // Content validation based on type
        when (post.type) {
            PostType.IMAGE, PostType.ALBUM -> {
                if (post.content.images.isEmpty()) {
                    errors.add("Image posts must contain at least one image")
                }
            }
            PostType.VIDEO, PostType.REEL -> {
                if (post.content.videos.isEmpty()) {
                    errors.add("Video posts must contain video content")
                }
            }
            PostType.POLL -> {
                if (post.content.poll == null || post.content.poll.options.size < 2) {
                    errors.add("Poll posts must have at least 2 options")
                }
            }
            PostType.TEXT -> {
                if (post.content.text.isNullOrBlank()) {
                    errors.add("Text posts must contain text content")
                }
            }
            PostType.REVIEW -> {
                if (post.metadata?.rating == null) {
                    errors.add("Review posts must include a rating")
                }
            }
            PostType.CHECK_IN -> {
                if (post.metadata?.location == null) {
                    errors.add("Check-in posts must include location")
                }
            }
            else -> {} // Other types have flexible content requirements
        }

        // Privacy validation
        if (post.privacy == PostPrivacy.CUSTOM && post.audience == null) {
            errors.add("Custom privacy requires audience specification")
        }

        // Engagement validation
        if (post.interactions.engagementRate < 0 || post.interactions.engagementRate > 1) {
            errors.add("Engagement rate must be between 0 and 1")
        }

        return errors
    }

    fun isValidPostType(type: PostType, content: PostContent): Boolean {
        return when (type) {
            PostType.IMAGE, PostType.ALBUM -> content.images.isNotEmpty()
            PostType.VIDEO, PostType.REEL, PostType.LIVE_STREAM -> content.videos.isNotEmpty()
            PostType.POLL -> content.poll != null
            PostType.TEXT -> !content.text.isNullOrBlank()
            else -> true // Other types are flexible
        }
    }
}

/**
 * Post Engagement Calculator
 */
object PostEngagementCalculator {

    fun calculateEngagementRate(post: Post): Float {
        val totalEngagements = post.interactions.reactions.size +
                             post.interactions.comments.size +
                             post.interactions.shares.size
        val views = post.interactions.views

        return if (views > 0) {
            (totalEngagements.toFloat() / views.toFloat()).coerceIn(0f, 1f)
        } else {
            0f
        }
    }

    fun getTopReaction(post: Post): ReactionType? {
        return post.interactions.reactions
            .groupBy { it.type }
            .maxByOrNull { it.value.size }
            ?.key
    }

    fun getEngagementSummary(post: Post): Map<String, Int> {
        return mapOf(
            "reactions" to post.interactions.reactions.size,
            "comments" to post.interactions.comments.size,
            "shares" to post.interactions.shares.size,
            "saves" to post.interactions.saves.size,
            "views" to post.interactions.views
        )
    }
}

/**
 * Extension functions to convert between Review and Post
 */
fun Review.toPost(): Post {
    return Post(
        id = this.id,
        type = PostType.REVIEW,
        user = PostUser(
            id = this.userId,
            username = this.userName,
            displayName = this.userName,
            avatarUrl = this.userAvatar,
            isVerified = this.isVerifiedPurchase
        ),
        content = PostContent(
            type = when (this.content.type) {
                ReviewType.TEXT -> PostContentType.TEXT
                ReviewType.IMAGE -> PostContentType.IMAGE
                ReviewType.VIDEO -> PostContentType.VIDEO
                ReviewType.MIXED, ReviewType.DETAILED -> PostContentType.MIXED
                ReviewType.QUICK -> PostContentType.TEXT
            },
            text = this.content.text,
            images = this.content.images.map { img ->
                PostImage(
                    id = img.id,
                    url = img.url,
                    caption = img.caption,
                    aspectRatio = img.aspectRatio
                )
            },
            videos = this.content.videos.map { vid ->
                PostVideo(
                    id = vid.id,
                    url = vid.url,
                    thumbnailUrl = vid.thumbnailUrl,
                    duration = vid.duration,
                    caption = vid.caption
                )
            },
            hashtags = this.content.hashtags,
            mentions = this.content.mentions
        ),
        interactions = PostInteractions(
            views = this.likes, // Use likes as views for compatibility
            isLiked = this.isLiked,
            isBookmarked = this.isBookmarked
        ),
        metadata = PostMetadata(
            targetType = this.targetType.name.lowercase(),
            targetId = this.targetId,
            targetName = this.targetName,
            rating = this.rating.value
        ),
        createdAt = this.createdAt,
        updatedAt = this.updatedAt
    )
}

fun Post.toReview(): Review? {
    if (this.type != PostType.REVIEW) return null

    return Review(
        id = this.id,
        userId = this.user.id,
        userName = this.user.username,
        userAvatar = this.user.avatarUrl,
        targetType = ReviewTargetType.valueOf(
            this.metadata?.targetType?.uppercase() ?: "PRODUCT"
        ),
        targetId = this.metadata?.targetId ?: "",
        targetName = this.metadata?.targetName ?: "",
        rating = ReviewRating(
            type = ReviewRatingType.STARS_5,
            value = this.metadata?.rating ?: 0.8f,
            displayValue = "${((this.metadata?.rating ?: 0.8f) * 5).toInt()}/5"
        ),
        content = ReviewContent(
            type = when (this.content.type) {
                PostContentType.TEXT -> ReviewType.TEXT
                PostContentType.IMAGE -> ReviewType.IMAGE
                PostContentType.VIDEO -> ReviewType.VIDEO
                PostContentType.MIXED -> ReviewType.MIXED
                else -> ReviewType.TEXT
            },
            text = this.content.text,
            images = this.content.images.map { img ->
                ReviewImage(
                    id = img.id,
                    url = img.url,
                    caption = img.caption,
                    aspectRatio = img.aspectRatio
                )
            },
            videos = this.content.videos.map { vid ->
                ReviewVideo(
                    id = vid.id,
                    url = vid.url,
                    thumbnailUrl = vid.thumbnailUrl,
                    duration = vid.duration,
                    caption = vid.caption
                )
            },
            hashtags = this.content.hashtags,
            mentions = this.content.mentions
        ),
        isVerifiedPurchase = this.user.isVerified,
        likes = this.interactions.likes,
        dislikes = 0, // No dislikes in new system
        comments = this.interactions.commentsCount,
        shares = this.interactions.sharesCount,
        isLiked = this.interactions.isLiked,
        isBookmarked = this.interactions.isBookmarked,
        createdAt = this.createdAt,
        updatedAt = this.updatedAt
    )
}