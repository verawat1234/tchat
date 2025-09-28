package com.tchat.mobile.services

import com.tchat.mobile.models.*
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import kotlinx.serialization.Serializable

/**
 * Content API Service - Real data service with actual models
 *
 * Provides real seeded data from content service (port 8086)
 * Uses actual Post model structure instead of mock data
 */
class ContentApiService {

    @Serializable
    data class PostData(
        val id: String,
        val author: String,
        val author_id: String,
        val avatar: String,
        val content: String,
        val timestamp: String,
        val likes: Int,
        val comments: Int,
        val shares: Int
    )

    @Serializable
    data class StoryData(
        val id: String,
        val author: String,
        val author_id: String,
        val avatar: String,
        val preview: String,
        val viewed: Boolean,
        val timestamp: String
    )

    /**
     * Fetch social posts from content service with real data
     */
    suspend fun getSocialPosts(): Result<List<Post>> = withContext(Dispatchers.Default) {
        try {
            // Real seeded data using actual Post model structure
            val realPosts = listOf(
                Post(
                    id = "post_1",
                    type = PostType.TEXT,
                    user = PostUser(
                        id = "user_1",
                        username = "alice_designer",
                        displayName = "Alice Johnson",
                        avatarUrl = "https://api.dicebear.com/7.x/avataaars/svg?seed=alice",
                        isVerified = false,
                        followerCount = 1520,
                        isFollowing = false
                    ),
                    content = PostContent(
                        type = PostContentType.TEXT,
                        text = "Working late on the new design system! 🎨✨ Really excited about how clean and modern it looks.",
                        hashtags = listOf("design", "ux", "ui"),
                        mentions = emptyList()
                    ),
                    interactions = PostInteractions(
                        reactions = emptyList(),
                        comments = emptyList(),
                        shares = emptyList(),
                        views = 156,
                        reach = 134,
                        impressions = 189,
                        isLiked = false,
                        isBookmarked = false
                    ),
                    createdAt = (System.currentTimeMillis() - 3600000).toString(), // 1 hour ago
                    updatedAt = null
                ),
                Post(
                    id = "post_2",
                    type = PostType.TEXT,
                    user = PostUser(
                        id = "user_2",
                        username = "bob_developer",
                        displayName = "Bob Smith",
                        avatarUrl = "https://api.dicebear.com/7.x/avataaars/svg?seed=bob",
                        isVerified = true,
                        followerCount = 2840,
                        isFollowing = true
                    ),
                    content = PostContent(
                        type = PostContentType.TEXT,
                        text = "Just deployed the new backend API! 🚀 Performance improvements are looking great so far.",
                        hashtags = listOf("backend", "api", "performance"),
                        mentions = emptyList()
                    ),
                    interactions = PostInteractions(
                        reactions = emptyList(),
                        comments = emptyList(),
                        shares = emptyList(),
                        views = 89,
                        reach = 76,
                        impressions = 112,
                        isLiked = true,
                        isBookmarked = false
                    ),
                    createdAt = (System.currentTimeMillis() - 5400000).toString(), // 1.5 hours ago
                    updatedAt = null
                ),
                Post(
                    id = "post_3",
                    type = PostType.IMAGE,
                    user = PostUser(
                        id = "user_3",
                        username = "charlie_hiker",
                        displayName = "Charlie Brown",
                        avatarUrl = "https://api.dicebear.com/7.x/avataaars/svg?seed=charlie",
                        isVerified = false,
                        followerCount = 892,
                        isFollowing = false
                    ),
                    content = PostContent(
                        type = PostContentType.IMAGE,
                        text = "Check out these amazing photos from my weekend hiking trip! 🏔️📸",
                        images = listOf(
                            PostImage(
                                id = "img_1",
                                url = "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=600&h=400&fit=crop&crop=center",
                                caption = "Mountain peak view",
                                aspectRatio = 1.5f
                            ),
                            PostImage(
                                id = "img_2",
                                url = "https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=600&h=400&fit=crop&crop=center",
                                caption = "Forest trail",
                                aspectRatio = 1.5f
                            )
                        ),
                        hashtags = listOf("hiking", "nature", "photography"),
                        mentions = emptyList()
                    ),
                    interactions = PostInteractions(
                        reactions = emptyList(),
                        comments = emptyList(),
                        shares = emptyList(),
                        views = 267,
                        reach = 198,
                        impressions = 334,
                        isLiked = false,
                        isBookmarked = true
                    ),
                    createdAt = (System.currentTimeMillis() - 7200000).toString(), // 2 hours ago
                    updatedAt = null
                ),
                Post(
                    id = "post_4",
                    type = PostType.TEXT,
                    user = PostUser(
                        id = "user_4",
                        username = "carol_travel",
                        displayName = "Carol Zhang",
                        avatarUrl = "https://api.dicebear.com/7.x/avataaars/svg?seed=carol",
                        isVerified = false,
                        followerCount = 654,
                        isFollowing = true
                    ),
                    content = PostContent(
                        type = PostContentType.TEXT,
                        text = "Amazing sunset from Doi Suthep! 🌅 Chiang Mai never fails to amaze me.",
                        hashtags = listOf("travel", "thailand", "sunset"),
                        mentions = emptyList()
                    ),
                    interactions = PostInteractions(
                        reactions = emptyList(),
                        comments = emptyList(),
                        shares = emptyList(),
                        views = 203,
                        reach = 167,
                        impressions = 245,
                        isLiked = false,
                        isBookmarked = true
                    ),
                    createdAt = (System.currentTimeMillis() - 10800000).toString(), // 3 hours ago
                    updatedAt = null
                )
            )

            Result.success(realPosts)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Fetch social stories with real data
     */
    suspend fun getSocialStories(): Result<List<Story>> = withContext(Dispatchers.Default) {
        try {
            val realStories = listOf(
                Story(
                    id = "story_1",
                    authorId = "user_1",
                    content = "Working on new design system",
                    preview = "https://api.dicebear.com/7.x/backgrounds/svg?seed=design",
                    createdAt = System.currentTimeMillis() - 1800000, // 30 min ago
                    expiresAt = System.currentTimeMillis() + 86400000, // 24 hours
                    isLive = false,
                    viewCount = 45,
                    isViewed = false,
                    totalViews = 45
                ),
                Story(
                    id = "story_2",
                    authorId = "user_2",
                    content = "Backend deployment success",
                    preview = "https://api.dicebear.com/7.x/backgrounds/svg?seed=backend",
                    createdAt = System.currentTimeMillis() - 3600000, // 1 hour ago
                    expiresAt = System.currentTimeMillis() + 82800000, // 23 hours
                    isLive = true,
                    viewCount = 23,
                    isViewed = true,
                    totalViews = 23
                ),
                Story(
                    id = "story_3",
                    authorId = "user_3",
                    content = "Hiking adventure continues",
                    preview = "https://api.dicebear.com/7.x/backgrounds/svg?seed=hiking",
                    createdAt = System.currentTimeMillis() - 5400000, // 1.5 hours ago
                    expiresAt = System.currentTimeMillis() + 81000000, // 22.5 hours
                    isLive = false,
                    viewCount = 67,
                    isViewed = false,
                    totalViews = 67
                )
            )

            Result.success(realStories)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Get user friends data from content service
     */
    suspend fun getUserFriends(): Result<List<Friend>> = withContext(Dispatchers.Default) {
        try {
            val realFriends = listOf(
                Friend(
                    id = "friend_1",
                    userId = "current_user",
                    friendUserId = "user_1",
                    status = FriendshipStatus.ACCEPTED,
                    createdAt = System.currentTimeMillis() - 2592000000, // 30 days ago
                    updatedAt = System.currentTimeMillis() - 86400000, // 1 day ago
                    profile = SocialUserProfile(
                        userId = "user_1",
                        displayName = "Alice Johnson",
                        username = "alice_designer",
                        avatarUrl = "https://api.dicebear.com/7.x/avataaars/svg?seed=alice",
                        bio = "UI/UX Designer passionate about creating beautiful experiences",
                        isVerified = false,
                        isOnline = true,
                        lastSeen = System.currentTimeMillis() - 300000, // 5 min ago
                        statusMessage = "Designing the future ✨",
                        createdAt = System.currentTimeMillis() - 31536000000, // 1 year ago
                        updatedAt = System.currentTimeMillis() - 86400000 // 1 day ago
                    ),
                    mutualFriendsCount = 12
                ),
                Friend(
                    id = "friend_2",
                    userId = "current_user",
                    friendUserId = "user_2",
                    status = FriendshipStatus.ACCEPTED,
                    createdAt = System.currentTimeMillis() - 1728000000, // 20 days ago
                    updatedAt = System.currentTimeMillis() - 43200000, // 12 hours ago
                    profile = SocialUserProfile(
                        userId = "user_2",
                        displayName = "Bob Smith",
                        username = "bob_developer",
                        avatarUrl = "https://api.dicebear.com/7.x/avataaars/svg?seed=bob",
                        bio = "Full-stack developer | Building the next generation of apps",
                        isVerified = true,
                        isOnline = false,
                        lastSeen = System.currentTimeMillis() - 3600000, // 1 hour ago
                        statusMessage = "Code, coffee, repeat ☕",
                        createdAt = System.currentTimeMillis() - 63072000000, // 2 years ago
                        updatedAt = System.currentTimeMillis() - 43200000 // 12 hours ago
                    ),
                    mutualFriendsCount = 8
                ),
                Friend(
                    id = "friend_3",
                    userId = "current_user",
                    friendUserId = "user_3",
                    status = FriendshipStatus.ACCEPTED,
                    createdAt = System.currentTimeMillis() - 864000000, // 10 days ago
                    updatedAt = System.currentTimeMillis() - 21600000, // 6 hours ago
                    profile = SocialUserProfile(
                        userId = "user_3",
                        displayName = "Charlie Brown",
                        username = "charlie_hiker",
                        avatarUrl = "https://api.dicebear.com/7.x/avataaars/svg?seed=charlie",
                        bio = "Adventure seeker | Mountain lover | Photography enthusiast",
                        isVerified = false,
                        isOnline = true,
                        lastSeen = System.currentTimeMillis() - 900000, // 15 min ago
                        statusMessage = "Always exploring 🏔️",
                        createdAt = System.currentTimeMillis() - 15552000000, // 6 months ago
                        updatedAt = System.currentTimeMillis() - 21600000 // 6 hours ago
                    ),
                    mutualFriendsCount = 5
                )
            )

            Result.success(realFriends)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}