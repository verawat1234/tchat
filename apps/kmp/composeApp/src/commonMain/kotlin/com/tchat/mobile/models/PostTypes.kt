package com.tchat.mobile.models

import kotlinx.serialization.Serializable

/**
 * Generic Post System - Social Media Platform Architecture
 *
 * This system treats everything as posts (reviews, social posts, videos, etc.)
 * Similar to how Instagram, TikTok, and other platforms work
 */

enum class PostType {
    REVIEW,         // Product/shop reviews
    SOCIAL,         // Social media posts
    VIDEO,          // TikTok-style videos
    STORY,          // Story-style content
    ANNOUNCEMENT,   // Official announcements
    AD              // Sponsored content
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
data class PostInteractions(
    val likes: Int = 0,
    val dislikes: Int = 0,
    val comments: Int = 0,
    val shares: Int = 0,
    val views: Int = 0,
    val saves: Int = 0,
    val isLiked: Boolean = false,
    val isDisliked: Boolean = false,
    val isBookmarked: Boolean = false,
    val isFollowing: Boolean = false
)

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
    val sponsorName: String? = null
)

@Serializable
data class Post(
    val id: String,
    val type: PostType,
    val user: PostUser,
    val content: PostContent,
    val interactions: PostInteractions,
    val metadata: PostMetadata? = null,
    val createdAt: String,
    val updatedAt: String? = null,
    val isEdited: Boolean = false,
    val visibility: PostVisibility = PostVisibility.PUBLIC
)

enum class PostVisibility {
    PUBLIC,
    FRIENDS,
    PRIVATE,
    UNLISTED
}

@Serializable
data class PostComment(
    val id: String,
    val postId: String,
    val user: PostUser,
    val text: String,
    val parentCommentId: String? = null, // For reply threads
    val likes: Int = 0,
    val isLiked: Boolean = false,
    val createdAt: String,
    val isEdited: Boolean = false
)

@Serializable
data class PostHashtag(
    val tag: String,
    val count: Int,
    val isFollowing: Boolean = false,
    val category: String? = null
)

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
            likes = this.likes,
            dislikes = this.dislikes,
            comments = this.comments,
            shares = this.shares,
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
        dislikes = this.interactions.dislikes,
        comments = this.interactions.comments,
        shares = this.interactions.shares,
        isLiked = this.interactions.isLiked,
        isBookmarked = this.interactions.isBookmarked,
        createdAt = this.createdAt,
        updatedAt = this.updatedAt
    )
}