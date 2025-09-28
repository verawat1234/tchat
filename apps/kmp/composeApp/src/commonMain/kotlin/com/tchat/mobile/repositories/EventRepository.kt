package com.tchat.mobile.repositories

import com.tchat.mobile.database.TchatDatabase
import com.tchat.mobile.models.*
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.IO
import kotlinx.coroutines.withContext

/**
 * Repository for event-related data operations
 * Handles events, categories, and event posts with SQLDelight integration
 */
class EventRepository(private val database: TchatDatabase) {

    /**
     * Seed event categories with predefined data
     */
    suspend fun seedEventCategories(): Result<Unit> = withContext(Dispatchers.IO) {
        try {
            val currentTime = System.currentTimeMillis()
            database.eventQueries.seedEventCategories(
                currentTime, currentTime, // 1st category
                currentTime, currentTime, // 2nd category
                currentTime, currentTime, // 3rd category
                currentTime, currentTime, // 4th category
                currentTime, currentTime, // 5th category
                currentTime, currentTime, // 6th category
                currentTime, currentTime, // 7th category
                currentTime, currentTime  // 8th category
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Seed event posts with predefined data
     */
    suspend fun seedEventPosts(): Result<Unit> = withContext(Dispatchers.IO) {
        try {
            val currentTime = System.currentTimeMillis()
            database.eventQueries.seedEventPosts(
                currentTime, currentTime, // 1st post
                currentTime, currentTime, // 2nd post
                currentTime, currentTime, // 3rd post
                currentTime, currentTime, // 4th post
                currentTime, currentTime, // 5th post
                currentTime, currentTime, // 6th post
                currentTime, currentTime, // 7th post
                currentTime, currentTime  // 8th post
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Get all event categories
     */
    suspend fun getAllEventCategories(): Result<List<EventCategory>> = withContext(Dispatchers.IO) {
        try {
            val categories = database.eventQueries.getAllEventCategories().executeAsList().map { row ->
                EventCategory(
                    id = row.id,
                    name = row.name,
                    description = row.description ?: "",
                    icon = row.icon ?: "",
                    color = row.color ?: "#3B82F6",
                    eventCount = row.event_count.toInt(),
                    isFeatured = row.is_featured == 1L,
                    sortOrder = row.sort_order.toInt(),
                    createdAt = row.created_at,
                    updatedAt = row.updated_at
                )
            }
            Result.success(categories)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Get featured event categories only
     */
    suspend fun getFeaturedEventCategories(): Result<List<EventCategory>> = withContext(Dispatchers.IO) {
        try {
            val categories = database.eventQueries.getFeaturedEventCategories().executeAsList().map { row ->
                EventCategory(
                    id = row.id,
                    name = row.name,
                    description = row.description ?: "",
                    icon = row.icon ?: "",
                    color = row.color ?: "#3B82F6",
                    eventCount = row.event_count.toInt(),
                    isFeatured = row.is_featured == 1L,
                    sortOrder = row.sort_order.toInt(),
                    createdAt = row.created_at,
                    updatedAt = row.updated_at
                )
            }
            Result.success(categories)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Get event categories with live event counts
     */
    suspend fun getEventCategoriesWithCounts(): Result<List<EventCategoryWithCount>> = withContext(Dispatchers.IO) {
        try {
            val currentTime = System.currentTimeMillis()
            val categories = database.eventQueries.getEventCategoriesWithCounts(currentTime).executeAsList().map { row ->
                EventCategoryWithCount(
                    category = EventCategory(
                        id = row.id,
                        name = row.name,
                        description = row.description ?: "",
                        icon = row.icon ?: "",
                        color = row.color ?: "#3B82F6",
                        eventCount = row.event_count.toInt(),
                        isFeatured = row.is_featured == 1L,
                        sortOrder = row.sort_order.toInt(),
                        createdAt = row.created_at,
                        updatedAt = row.updated_at
                    ),
                    liveEventCount = row.live_event_count?.toInt() ?: 0
                )
            }
            Result.success(categories)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Get all event posts across all events
     */
    suspend fun getAllEventPosts(): Result<List<EventPostWithDetails>> = withContext(Dispatchers.IO) {
        try {
            val posts = database.eventQueries.getAllEventPosts().executeAsList().map { row ->
                EventPostWithDetails(
                    post = EventPost(
                        id = row.id,
                        eventId = row.event_id,
                        authorId = row.author_id,
                        content = row.content,
                        imageUrl = row.image_url ?: "",
                        postType = EventPostType.valueOf(row.post_type.uppercase()),
                        likesCount = row.likes_count.toInt(),
                        commentsCount = row.comments_count.toInt(),
                        sharesCount = row.shares_count.toInt(),
                        isPinned = row.is_pinned == 1L,
                        createdAt = row.created_at,
                        updatedAt = row.updated_at
                    ),
                    authorName = row.author_name ?: "",
                    authorUsername = row.author_username ?: "",
                    authorAvatar = row.author_avatar ?: "",
                    authorVerified = row.author_verified == 1L,
                    eventTitle = row.event_title ?: "",
                    eventCategory = row.event_category ?: ""
                )
            }
            Result.success(posts)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Get event posts for a specific event
     */
    suspend fun getEventPosts(eventId: String): Result<List<EventPostWithDetails>> = withContext(Dispatchers.IO) {
        try {
            val posts = database.eventQueries.getEventPosts(eventId).executeAsList().map { row ->
                EventPostWithDetails(
                    post = EventPost(
                        id = row.id,
                        eventId = row.event_id,
                        authorId = row.author_id,
                        content = row.content,
                        imageUrl = row.image_url ?: "",
                        postType = EventPostType.valueOf(row.post_type.uppercase()),
                        likesCount = row.likes_count.toInt(),
                        commentsCount = row.comments_count.toInt(),
                        sharesCount = row.shares_count.toInt(),
                        isPinned = row.is_pinned == 1L,
                        createdAt = row.created_at,
                        updatedAt = row.updated_at
                    ),
                    authorName = row.author_name ?: "",
                    authorUsername = row.author_username ?: "",
                    authorAvatar = row.author_avatar ?: "",
                    authorVerified = row.author_verified == 1L,
                    eventTitle = "",
                    eventCategory = ""
                )
            }
            Result.success(posts)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Get pinned event posts across all events
     */
    suspend fun getPinnedEventPosts(): Result<List<EventPostWithDetails>> = withContext(Dispatchers.IO) {
        try {
            val posts = database.eventQueries.getPinnedEventPosts().executeAsList().map { row ->
                EventPostWithDetails(
                    post = EventPost(
                        id = row.id,
                        eventId = row.event_id,
                        authorId = row.author_id,
                        content = row.content,
                        imageUrl = row.image_url ?: "",
                        postType = EventPostType.valueOf(row.post_type.uppercase()),
                        likesCount = row.likes_count.toInt(),
                        commentsCount = row.comments_count.toInt(),
                        sharesCount = row.shares_count.toInt(),
                        isPinned = row.is_pinned == 1L,
                        createdAt = row.created_at,
                        updatedAt = row.updated_at
                    ),
                    authorName = row.author_name ?: "",
                    authorUsername = row.author_username ?: "",
                    authorAvatar = row.author_avatar ?: "",
                    authorVerified = row.author_verified == 1L,
                    eventTitle = row.event_title ?: "",
                    eventCategory = ""
                )
            }
            Result.success(posts)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Convert event categories to UI models for Social screen
     */
    suspend fun getEventCategoriesForUI(): List<Pair<String, Int>> {
        return try {
            val categories = getFeaturedEventCategories().getOrNull() ?: emptyList()
            categories.map { it.name to it.eventCount }
        } catch (e: Exception) {
            // Fallback to hardcoded categories
            listOf(
                "Music" to 15,
                "Food" to 23,
                "Technology" to 8,
                "Arts & Culture" to 12
            )
        }
    }

    /**
     * Convert event posts to UI models for Social screen
     */
    suspend fun getEventPostsForUI(): List<PostItem> {
        return try {
            val posts = getAllEventPosts().getOrNull() ?: emptyList()
            posts.map { postWithDetails ->
                PostItem(
                    id = postWithDetails.post.id,
                    author = UserItem(
                        id = postWithDetails.post.authorId,
                        name = postWithDetails.authorName,
                        username = postWithDetails.authorUsername,
                        avatar = postWithDetails.authorAvatar,
                        isVerified = postWithDetails.authorVerified
                    ),
                    content = postWithDetails.post.content,
                    timestamp = formatTimestamp(postWithDetails.post.createdAt),
                    likes = postWithDetails.post.likesCount,
                    comments = postWithDetails.post.commentsCount,
                    shares = postWithDetails.post.sharesCount,
                    imageUrl = if (postWithDetails.post.imageUrl.isNotEmpty()) postWithDetails.post.imageUrl else null,
                    tags = listOf(postWithDetails.eventCategory),
                    type = when (postWithDetails.post.postType) {
                        EventPostType.IMAGE -> SocialPostType.IMAGE
                        EventPostType.VIDEO -> SocialPostType.LIVE
                        else -> SocialPostType.TEXT
                    },
                    source = PostSource.FOLLOWING,
                    isLiked = false
                )
            }
        } catch (e: Exception) {
            emptyList()
        }
    }

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

    /**
     * Initialize repository with seed data
     */
    suspend fun initializeWithSeedData(): Result<Unit> {
        return try {
            seedEventCategories().getOrThrow()
            seedEventPosts().getOrThrow()
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}

// Data models for event categories and posts
data class EventCategory(
    val id: String,
    val name: String,
    val description: String,
    val icon: String,
    val color: String,
    val eventCount: Int,
    val isFeatured: Boolean,
    val sortOrder: Int,
    val createdAt: Long,
    val updatedAt: Long
)

data class EventCategoryWithCount(
    val category: EventCategory,
    val liveEventCount: Int
)

data class EventPost(
    val id: String,
    val eventId: String,
    val authorId: String,
    val content: String,
    val imageUrl: String,
    val postType: EventPostType,
    val likesCount: Int,
    val commentsCount: Int,
    val sharesCount: Int,
    val isPinned: Boolean,
    val createdAt: Long,
    val updatedAt: Long
)

data class EventPostWithDetails(
    val post: EventPost,
    val authorName: String,
    val authorUsername: String,
    val authorAvatar: String,
    val authorVerified: Boolean,
    val eventTitle: String,
    val eventCategory: String
)

enum class EventPostType {
    TEXT, IMAGE, VIDEO, LIVE
}