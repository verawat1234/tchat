package com.tchat.mobile.repositories

import com.tchat.mobile.models.*
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.map
import com.tchat.mobile.utils.PlatformUtils

/**
 * Post Repository - Mock API Implementation
 *
 * Provides social media platform functionality with repository pattern
 * Ready for real API integration
 */

interface PostRepository {
    // Core Post Operations
    suspend fun getPosts(type: PostType? = null, userId: String? = null): Result<List<Post>>
    suspend fun getPost(postId: String): Result<Post>
    suspend fun createPost(post: Post): Result<Post>
    suspend fun updatePost(postId: String, post: Post): Result<Post>
    suspend fun deletePost(postId: String): Result<Boolean>

    // Interactions
    suspend fun likePost(postId: String): Result<PostInteractions>
    suspend fun unlikePost(postId: String): Result<PostInteractions>
    suspend fun bookmarkPost(postId: String): Result<Boolean>
    suspend fun sharePost(postId: String, platform: String): Result<Boolean>

    // Comments
    suspend fun getComments(postId: String): Result<List<PostComment>>
    suspend fun addComment(postId: String, text: String): Result<PostComment>
    suspend fun likeComment(commentId: String): Result<PostComment>
    suspend fun deleteComment(commentId: String): Result<Boolean>

    // Hashtags & Search
    suspend fun getHashtagPosts(hashtag: String): Result<List<Post>>
    suspend fun getTrendingHashtags(): Result<List<PostHashtag>>
    suspend fun searchPosts(query: String): Result<List<Post>>

    // User Operations
    suspend fun getUserPosts(userId: String): Result<List<Post>>
    suspend fun followUser(userId: String): Result<Boolean>
    suspend fun unfollowUser(userId: String): Result<Boolean>

    // Real-time updates
    fun observePost(postId: String): Flow<Post?>
    fun observePostInteractions(postId: String): Flow<PostInteractions>
}

class MockPostRepository : PostRepository {

    private val _posts = MutableStateFlow(generateMockPosts())
    private val _comments = MutableStateFlow(generateMockComments())
    private val _hashtags = MutableStateFlow(generateMockHashtags())

    private val posts: List<Post> get() = _posts.value
    private val comments: List<PostComment> get() = _comments.value
    private val hashtags: List<PostHashtag> get() = _hashtags.value

    override suspend fun getPosts(type: PostType?, userId: String?): Result<List<Post>> {
        delay(300) // Simulate network delay

        var filteredPosts = posts

        if (type != null) {
            filteredPosts = filteredPosts.filter { it.type == type }
        }

        if (userId != null) {
            filteredPosts = filteredPosts.filter { it.user.id == userId }
        }

        return Result.success(filteredPosts.sortedByDescending { it.createdAt })
    }

    override suspend fun getPost(postId: String): Result<Post> {
        delay(200)

        val post = posts.find { it.id == postId }
        return if (post != null) {
            Result.success(post)
        } else {
            Result.failure(Exception("Post not found"))
        }
    }

    override suspend fun createPost(post: Post): Result<Post> {
        delay(500)

        val newPost = post.copy(
            id = "post_${PlatformUtils.currentTimeMillis()}",
            createdAt = getCurrentTimestamp()
        )

        _posts.value = _posts.value + newPost
        return Result.success(newPost)
    }

    override suspend fun updatePost(postId: String, post: Post): Result<Post> {
        delay(400)

        val updatedPosts = _posts.value.map { currentPost ->
            if (currentPost.id == postId) {
                post.copy(
                    id = postId,
                    updatedAt = getCurrentTimestamp(),
                    isEdited = true
                )
            } else {
                currentPost
            }
        }

        _posts.value = updatedPosts
        val updatedPost = updatedPosts.find { it.id == postId }

        return if (updatedPost != null) {
            Result.success(updatedPost)
        } else {
            Result.failure(Exception("Post not found"))
        }
    }

    override suspend fun deletePost(postId: String): Result<Boolean> {
        delay(300)

        val originalSize = _posts.value.size
        _posts.value = _posts.value.filter { it.id != postId }

        return Result.success(_posts.value.size < originalSize)
    }

    override suspend fun likePost(postId: String): Result<PostInteractions> {
        delay(150)

        val updatedPosts = _posts.value.map { post ->
            if (post.id == postId) {
                val newInteractions = post.interactions.copy(
                    likes = if (post.interactions.isLiked) {
                        post.interactions.likes - 1
                    } else {
                        post.interactions.likes + 1
                    },
                    isLiked = !post.interactions.isLiked,
                    isDisliked = false // Remove dislike if exists
                )
                post.copy(interactions = newInteractions)
            } else {
                post
            }
        }

        _posts.value = updatedPosts
        val updatedPost = updatedPosts.find { it.id == postId }

        return if (updatedPost != null) {
            Result.success(updatedPost.interactions)
        } else {
            Result.failure(Exception("Post not found"))
        }
    }

    override suspend fun unlikePost(postId: String): Result<PostInteractions> {
        return likePost(postId) // Same logic for toggle
    }

    override suspend fun bookmarkPost(postId: String): Result<Boolean> {
        delay(200)

        val updatedPosts = _posts.value.map { post ->
            if (post.id == postId) {
                post.copy(
                    interactions = post.interactions.copy(
                        isBookmarked = !post.interactions.isBookmarked,
                        saves = if (post.interactions.isBookmarked) {
                            post.interactions.saves - 1
                        } else {
                            post.interactions.saves + 1
                        }
                    )
                )
            } else {
                post
            }
        }

        _posts.value = updatedPosts
        val updatedPost = updatedPosts.find { it.id == postId }

        return Result.success(updatedPost?.interactions?.isBookmarked ?: false)
    }

    override suspend fun sharePost(postId: String, platform: String): Result<Boolean> {
        delay(250)

        val updatedPosts = _posts.value.map { post ->
            if (post.id == postId) {
                post.copy(
                    interactions = post.interactions.copy(
                        shares = post.interactions.shares + 1
                    )
                )
            } else {
                post
            }
        }

        _posts.value = updatedPosts

        // Mock platform-specific sharing
        println("Shared post $postId to $platform")

        return Result.success(true)
    }

    override suspend fun getComments(postId: String): Result<List<PostComment>> {
        delay(250)

        val postComments = comments.filter { it.postId == postId }
            .sortedByDescending { it.createdAt }

        return Result.success(postComments)
    }

    override suspend fun addComment(postId: String, text: String): Result<PostComment> {
        delay(400)

        val newComment = PostComment(
            id = "comment_${PlatformUtils.currentTimeMillis()}",
            postId = postId,
            user = PostUser(
                id = "current_user",
                username = "current_user",
                displayName = "You",
                avatarUrl = null,
                isVerified = false
            ),
            text = text,
            createdAt = getCurrentTimestamp()
        )

        _comments.value = _comments.value + newComment

        // Update post comment count
        _posts.value = _posts.value.map { post ->
            if (post.id == postId) {
                post.copy(
                    interactions = post.interactions.copy(
                        comments = post.interactions.comments + 1
                    )
                )
            } else {
                post
            }
        }

        return Result.success(newComment)
    }

    override suspend fun likeComment(commentId: String): Result<PostComment> {
        delay(150)

        val updatedComments = _comments.value.map { comment ->
            if (comment.id == commentId) {
                comment.copy(
                    likes = if (comment.isLiked) comment.likes - 1 else comment.likes + 1,
                    isLiked = !comment.isLiked
                )
            } else {
                comment
            }
        }

        _comments.value = updatedComments
        val updatedComment = updatedComments.find { it.id == commentId }

        return if (updatedComment != null) {
            Result.success(updatedComment)
        } else {
            Result.failure(Exception("Comment not found"))
        }
    }

    override suspend fun deleteComment(commentId: String): Result<Boolean> {
        delay(200)

        val originalSize = _comments.value.size
        _comments.value = _comments.value.filter { it.id != commentId }

        return Result.success(_comments.value.size < originalSize)
    }

    override suspend fun getHashtagPosts(hashtag: String): Result<List<Post>> {
        delay(300)

        val hashtagPosts = posts.filter { post ->
            post.content.hashtags.any { tag ->
                tag.lowercase().contains(hashtag.lowercase().removePrefix("#"))
            }
        }.sortedByDescending { it.createdAt }

        return Result.success(hashtagPosts)
    }

    override suspend fun getTrendingHashtags(): Result<List<PostHashtag>> {
        delay(200)
        return Result.success(hashtags.sortedByDescending { it.count })
    }

    override suspend fun searchPosts(query: String): Result<List<Post>> {
        delay(400)

        val searchResults = posts.filter { post ->
            post.content.text?.contains(query, ignoreCase = true) == true ||
            post.content.hashtags.any { it.contains(query, ignoreCase = true) } ||
            post.user.username.contains(query, ignoreCase = true) ||
            post.metadata?.targetName?.contains(query, ignoreCase = true) == true
        }.sortedByDescending { it.createdAt }

        return Result.success(searchResults)
    }

    override suspend fun getUserPosts(userId: String): Result<List<Post>> {
        return getPosts(userId = userId)
    }

    override suspend fun followUser(userId: String): Result<Boolean> {
        delay(300)

        _posts.value = _posts.value.map { post ->
            if (post.user.id == userId) {
                post.copy(
                    user = post.user.copy(
                        isFollowing = true,
                        followerCount = post.user.followerCount + 1
                    )
                )
            } else {
                post
            }
        }

        return Result.success(true)
    }

    override suspend fun unfollowUser(userId: String): Result<Boolean> {
        delay(300)

        _posts.value = _posts.value.map { post ->
            if (post.user.id == userId) {
                post.copy(
                    user = post.user.copy(
                        isFollowing = false,
                        followerCount = maxOf(0, post.user.followerCount - 1)
                    )
                )
            } else {
                post
            }
        }

        return Result.success(true)
    }

    override fun observePost(postId: String): Flow<Post?> {
        return _posts.asStateFlow().map { posts ->
            posts.find { it.id == postId }
        }
    }

    override fun observePostInteractions(postId: String): Flow<PostInteractions> {
        return _posts.asStateFlow().map { posts ->
            posts.find { it.id == postId }?.interactions ?: PostInteractions()
        }
    }

    // Helper functions
    private fun getCurrentTimestamp(): String {
        // In real app, this would be proper timestamp formatting
        return "${PlatformUtils.currentTimeMillis()}"
    }

    private fun generateMockPosts(): List<Post> {
        // Convert existing reviews to posts for backwards compatibility
        val reviewPosts = listOf(
            // Mock social posts mixed with reviews
            Post(
                id = "social_1",
                type = PostType.SOCIAL,
                user = PostUser(
                    id = "social_user_1",
                    username = "food_explorer",
                    displayName = "Bangkok Food Explorer 🍜",
                    avatarUrl = "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=100",
                    isVerified = true,
                    followerCount = 12500,
                    isFollowing = false
                ),
                content = PostContent(
                    type = PostContentType.MIXED,
                    text = "Just discovered this amazing hidden gem in Bangkok! The Pad Thai here is absolutely incredible. Who else loves finding those secret foodie spots? 🥢✨",
                    images = listOf(
                        PostImage("img_social_1", "https://images.unsplash.com/photo-1559847844-5315695dadae?w=800", "Amazing Pad Thai", 1.2f),
                        PostImage("img_social_2", "https://images.unsplash.com/photo-1565299624946-b28f40a0ca4b?w=800", "Restaurant interior", 0.8f)
                    ),
                    hashtags = listOf("#bangkok", "#foodie", "#padthai", "#hiddengem", "#authentic", "#thai"),
                    mentions = listOf("@bangkokstreetfood"),
                    location = "Bangkok, Thailand"
                ),
                interactions = PostInteractions(
                    likes = 2347,
                    comments = 89,
                    shares = 56,
                    views = 15670,
                    saves = 234,
                    isLiked = false,
                    isBookmarked = true
                ),
                createdAt = "2 hours ago",
                visibility = PostVisibility.PUBLIC
            ),

            Post(
                id = "video_1",
                type = PostType.VIDEO,
                user = PostUser(
                    id = "video_user_1",
                    username = "chef_mike",
                    displayName = "Chef Mike 👨‍🍳",
                    avatarUrl = "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100",
                    isVerified = true,
                    followerCount = 89400,
                    isFollowing = true
                ),
                content = PostContent(
                    type = PostContentType.VIDEO,
                    text = "Secret technique for perfect fried rice! Save this for later 🔥 #cookinghacks",
                    videos = listOf(
                        PostVideo(
                            id = "vid_cooking_1",
                            url = "https://sample-videos.com/zip/10/mp4/SampleVideo_1280x720_1mb.mp4",
                            thumbnailUrl = "https://images.unsplash.com/photo-1516684669134-de6f7c473a2a?w=400",
                            duration = "0:47",
                            caption = "Perfect fried rice technique",
                            isAutoPlay = true
                        )
                    ),
                    hashtags = listOf("#cooking", "#friedrice", "#recipe", "#cheftips", "#viral", "#foodtok"),
                    location = "Professional Kitchen"
                ),
                interactions = PostInteractions(
                    likes = 45670,
                    comments = 1234,
                    shares = 2890,
                    views = 234500,
                    saves = 8970,
                    isLiked = true,
                    isBookmarked = true
                ),
                createdAt = "5 hours ago",
                visibility = PostVisibility.PUBLIC
            )
        )

        return reviewPosts
    }

    private fun generateMockComments(): List<PostComment> {
        return listOf(
            PostComment(
                id = "comment_1",
                postId = "social_1",
                user = PostUser("user_1", "food_lover", "Food Lover", null, false),
                text = "OMG I need to visit this place! Thanks for sharing 🤤",
                likes = 23,
                isLiked = false,
                createdAt = "1 hour ago"
            ),
            PostComment(
                id = "comment_2",
                postId = "social_1",
                user = PostUser("user_2", "bangkok_local", "Bangkok Local", null, true),
                text = "I've been there! The Tom Yum is also incredible. Hidden gem for sure!",
                likes = 45,
                isLiked = true,
                createdAt = "45 minutes ago"
            ),
            PostComment(
                id = "comment_3",
                postId = "video_1",
                user = PostUser("user_3", "home_chef", "Home Chef", null, false),
                text = "This technique changed my cooking game! Thanks Chef Mike! 👨‍🍳",
                likes = 156,
                isLiked = false,
                createdAt = "3 hours ago"
            ),
            PostComment(
                id = "comment_4",
                postId = "video_1",
                user = PostUser("user_4", "cooking_student", "Cooking Student", null, false),
                text = "Can you do a video on perfect scrambled eggs next? 🥚",
                likes = 89,
                isLiked = true,
                createdAt = "2 hours ago"
            )
        )
    }

    private fun generateMockHashtags(): List<PostHashtag> {
        return listOf(
            PostHashtag("#bangkok", 15670, false, "location"),
            PostHashtag("#foodie", 89450, true, "lifestyle"),
            PostHashtag("#cooking", 67890, true, "hobby"),
            PostHashtag("#viral", 234500, false, "trending"),
            PostHashtag("#recipe", 45670, true, "cooking"),
            PostHashtag("#thai", 23450, false, "cuisine"),
            PostHashtag("#authentic", 12340, false, "quality"),
            PostHashtag("#padthai", 8790, true, "dish"),
            PostHashtag("#cheftips", 34560, false, "education"),
            PostHashtag("#hiddengem", 5670, false, "discovery")
        )
    }
}