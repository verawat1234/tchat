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
                val currentReactions = post.interactions.reactions.toMutableList()
                val userId = "current_user" // In real app, get from auth service

                // Toggle like reaction
                val existingLike = currentReactions.find { it.userId == userId && it.type == ReactionType.LIKE }
                if (existingLike != null) {
                    currentReactions.remove(existingLike)
                } else {
                    currentReactions.add(PostReaction(
                        type = ReactionType.LIKE,
                        userId = userId,
                        timestamp = getCurrentTimestamp()
                    ))
                }

                val newInteractions = post.interactions.copy(
                    reactions = currentReactions,
                    isLiked = currentReactions.any { it.userId == userId && it.type == ReactionType.LIKE }
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
                            post.interactions.saves.filter { it != "current_user" }
                        } else {
                            post.interactions.saves + "current_user"
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
                        shares = post.interactions.shares + PostShare(
                            id = "share_${getCurrentTimestamp()}",
                            userId = "current_user",
                            userName = "You",
                            timestamp = getCurrentTimestamp(),
                            shareType = ShareType.DIRECT_SHARE
                        )
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

        val postComments = comments
            .sortedByDescending { it.timestamp }

        return Result.success(postComments)
    }

    override suspend fun addComment(postId: String, text: String): Result<PostComment> {
        delay(400)

        val newComment = PostComment(
            id = "comment_${PlatformUtils.currentTimeMillis()}",
            userId = "current_user",
            userName = "You",
            userAvatar = null,
            content = text,
            timestamp = getCurrentTimestamp()
        )

        _comments.value = _comments.value + newComment

        // Update post comment count
        _posts.value = _posts.value.map { post ->
            if (post.id == postId) {
                post.copy(
                    interactions = post.interactions.copy(
                        comments = post.interactions.comments + newComment
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
                val currentReactions = comment.reactions.toMutableList()
                val userId = "current_user" // In real app, get from auth service

                // Toggle like reaction
                val existingLike = currentReactions.find { it.userId == userId && it.type == ReactionType.LIKE }
                if (existingLike != null) {
                    currentReactions.remove(existingLike)
                } else {
                    currentReactions.add(PostReaction(
                        type = ReactionType.LIKE,
                        userId = userId,
                        timestamp = getCurrentTimestamp()
                    ))
                }

                comment.copy(
                    reactions = currentReactions
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
        return listOf(
            // ========== CORE CONTENT TYPES (8) ==========

            // 1. TEXT - Simple status updates
            Post(
                id = "text_1",
                type = PostType.TEXT,
                user = PostUser("user1", "mindful_jane", "Jane Wilson", "https://images.unsplash.com/photo-1494790108755-2616b332c2c2?w=100", false, 1240, false),
                content = PostContent(
                    type = PostContentType.TEXT,
                    text = "Sometimes the best moments are the quiet ones ☕ Starting my morning with gratitude and intention. What's bringing you joy today?",
                    hashtags = listOf("#mindfulness", "#morningritual", "#gratitude", "#selfcare", "#positivevibes")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "2 hours ago"
            ),

            // 2. IMAGE - Single/multiple photos
            Post(
                id = "image_1",
                type = PostType.IMAGE,
                user = PostUser("user2", "urban_photographer", "Alex Chen 📸", "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100", true, 15600, true),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "Golden hour in the city never gets old 🌅 The way light dances between buildings creates pure magic",
                    images = listOf(
                        PostImage("img1", "https://images.unsplash.com/photo-1514565131-fce0801e5785?w=800", "Cityscape at golden hour", 1.5f),
                        PostImage("img2", "https://images.unsplash.com/photo-1449824913935-59a10b8d2000?w=800", "Urban architecture", 1.0f)
                    ),
                    hashtags = listOf("#goldenhour", "#cityscape", "#photography", "#urban", "#architecture", "#lightplay")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "4 hours ago"
            ),

            // 3. VIDEO - Video posts
            Post(
                id = "video_1",
                type = PostType.VIDEO,
                user = PostUser("user3", "chef_maria", "Chef Maria 👩‍🍳", "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=100", true, 89400, true),
                content = PostContent(
                    type = PostContentType.VIDEO,
                    text = "30-second pasta technique that will change your life! Who's trying this tonight? 🍝✨",
                    videos = listOf(
                        PostVideo("vid1", "https://sample-videos.com/zip/10/mp4/SampleVideo_1280x720_1mb.mp4", "https://images.unsplash.com/photo-1516684669134-de6f7c473a2a?w=400", "0:30", "Perfect pasta technique")
                    ),
                    hashtags = listOf("#cooking", "#pasta", "#technique", "#chef", "#viral", "#foodhack")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "1 day ago"
            ),

            // 4. AUDIO - Voice notes, music
            Post(
                id = "audio_1",
                type = PostType.AUDIO,
                user = PostUser("user4", "indie_musician", "Sam Rivers 🎵", "https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=100", true, 34500, false),
                content = PostContent(
                    type = PostContentType.TEXT,
                    text = "Late night studio session vibes 🎹 Here's a snippet of something I'm working on. Should I finish this track?",
                    hashtags = listOf("#latenight", "#studio", "#indie", "#musician", "#wip", "#ambient")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "8 hours ago"
            ),

            // 5. LINK_SHARE - Shared articles/websites
            Post(
                id = "link_1",
                type = PostType.LINK_SHARE,
                user = PostUser("user5", "tech_sarah", "Sarah Tech 💻", "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=100", false, 8900, false),
                content = PostContent(
                    type = PostContentType.TEXT,
                    text = "This article about sustainable tech really opened my eyes 🌱 The future of development is green! What are your thoughts on eco-friendly coding practices?",
                    hashtags = listOf("#sustainabletech", "#greencomputing", "#ecofriendly", "#techforgood", "#sustainability")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "6 hours ago"
            ),

            // 6. POST_MESSAGE - Message posted to timeline
            Post(
                id = "message_1",
                type = PostType.POST_MESSAGE,
                user = PostUser("user6", "birthday_friend", "Mike Johnson", "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100", false, 450, true),
                content = PostContent(
                    type = PostContentType.TEXT,
                    text = "@jane_wilson Happy birthday! 🎉 Hope your special day is as amazing as you are! Let's celebrate soon 🥳",
                    mentions = listOf("jane_wilson"),
                    hashtags = listOf("#birthday", "#celebration", "#friend")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.FRIENDS,
                createdAt = "12 hours ago"
            ),

            // 7. REVIEW - Reviews of places, products, services
            Post(
                id = "review_1",
                type = PostType.REVIEW,
                user = PostUser("user7", "foodie_explorer", "David Kim 🍜", "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100", true, 23400, false),
                content = PostContent(
                    type = PostContentType.MIXED,
                    text = "Just had the most incredible ramen experience! The broth was rich and complex, noodles perfectly al dente. Service was outstanding too ⭐⭐⭐⭐⭐",
                    images = listOf(
                        PostImage("review_img1", "https://images.unsplash.com/photo-1569718212165-3a8278d5f624?w=800", "Amazing ramen bowl", 1.0f)
                    ),
                    hashtags = listOf("#ramen", "#foodreview", "#authentic", "#delicious", "#mustvisit"),
                    location = "Ichiran Ramen, Tokyo"
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                metadata = PostMetadata(
                    targetType = "restaurant",
                    targetId = "rest_ichiran_1",
                    targetName = "Ichiran Ramen",
                    rating = 4.8f
                ),
                createdAt = "1 day ago"
            ),

            // 8. ALBUM - Photo collections
            Post(
                id = "album_1",
                type = PostType.ALBUM,
                user = PostUser("user8", "travel_emma", "Emma Travel ✈️", "https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=100", true, 67800, true),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "Santorini memories 🇬🇷 Every sunset painted the sky in different colors. Already planning my return trip!",
                    images = listOf(
                        PostImage("album1", "https://images.unsplash.com/photo-1613395877344-13d4a8e0d49e?w=800", "Santorini sunset", 1.3f),
                        PostImage("album2", "https://images.unsplash.com/photo-1570077188670-e3a8d69ac5ff?w=800", "White buildings", 0.8f),
                        PostImage("album3", "https://images.unsplash.com/photo-1613395877344-13d4a8e0d49e?w=800", "Aegean Sea view", 1.5f)
                    ),
                    hashtags = listOf("#santorini", "#greece", "#travel", "#sunset", "#wanderlust", "#mediterranean"),
                    location = "Santorini, Greece"
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "3 days ago"
            ),

            // ========== RICH MEDIA TYPES (6) ==========

            // 9. STORY - Ephemeral 24h content
            Post(
                id = "story_1",
                type = PostType.STORY,
                user = PostUser("user9", "daily_vibes", "Chris Daily", "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=100", false, 2340, true),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "Coffee shop vibes ☕",
                    images = listOf(
                        PostImage("story1", "https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=800", "Coffee shop", 0.56f)
                    ),
                    hashtags = listOf("#coffeetime", "#vibes", "#story")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                metadata = PostMetadata(expiresAt = "22:30:00"),
                createdAt = "6 hours ago"
            ),

            // 10. REEL - Short-form vertical video
            Post(
                id = "reel_1",
                type = PostType.REEL,
                user = PostUser("user10", "dance_moves", "Maya Dance 💃", "https://images.unsplash.com/photo-1494790108755-2616b332c2c2?w=100", true, 234000, false),
                content = PostContent(
                    type = PostContentType.VIDEO,
                    text = "New choreography to this trending sound! 🕺 Who's learning this with me?",
                    videos = listOf(
                        PostVideo("reel1", "https://sample-videos.com/zip/10/mp4/SampleVideo_1280x720_1mb.mp4", "https://images.unsplash.com/photo-1547036967-23d11aacaee0?w=400", "0:15", "Dance choreography")
                    ),
                    hashtags = listOf("#dance", "#choreography", "#trending", "#viral", "#reels", "#dancechallenge")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "2 days ago"
            ),

            // 11. LIVE_STREAM - Live video broadcasts
            Post(
                id = "live_1",
                type = PostType.LIVE_STREAM,
                user = PostUser("user11", "gaming_pro", "Pro Gamer 🎮", "https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=100", true, 145000, true),
                content = PostContent(
                    type = PostContentType.LIVE,
                    text = "🔴 LIVE: Climbing to Legendary rank! Drop your tips in the chat!",
                    hashtags = listOf("#live", "#gaming", "#legendary", "#stream", "#pro", "#tips")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "30 minutes ago"
            ),

            // 12. PLAYLIST - Music/video collections
            Post(
                id = "playlist_1",
                type = PostType.PLAYLIST,
                user = PostUser("user12", "music_curator", "Beats Curator 🎧", "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=100", true, 56700, false),
                content = PostContent(
                    type = PostContentType.TEXT,
                    text = "Sunday Chill Vibes playlist is here 🎵 Perfect for lazy afternoons and coffee moments. 20 handpicked tracks ☕",
                    hashtags = listOf("#playlist", "#chillvibes", "#sunday", "#coffee", "#relax", "#curated")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "1 day ago"
            ),

            // 13. MOOD_BOARD - Visual inspiration
            Post(
                id = "mood_1",
                type = PostType.MOOD_BOARD,
                user = PostUser("user13", "design_inspiration", "Luna Design ✨", "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=100", true, 34500, true),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "Spring 2024 mood board 🌸 Soft pastels, natural textures, and minimalist vibes. What's inspiring you this season?",
                    images = listOf(
                        PostImage("mood1", "https://images.unsplash.com/photo-1615800001374-d4a8e1b4b6c2?w=400", "Pastel colors", 1.0f),
                        PostImage("mood2", "https://images.unsplash.com/photo-1586880244386-c398dc1b6f9a?w=400", "Natural textures", 1.0f),
                        PostImage("mood3", "https://images.unsplash.com/photo-1586880244386-c398dc1b6f9a?w=400", "Minimalist design", 1.0f)
                    ),
                    hashtags = listOf("#moodboard", "#spring2024", "#pastels", "#design", "#inspiration", "#minimal")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "2 days ago"
            ),

            // 14. TUTORIAL - How-to content
            Post(
                id = "tutorial_1",
                type = PostType.TUTORIAL,
                user = PostUser("user14", "diy_expert", "DIY Expert 🔧", "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100", true, 78900, false),
                content = PostContent(
                    type = PostContentType.VIDEO,
                    text = "How to create the perfect gallery wall! 🖼️ Save this for your next room makeover. Step-by-step guide in comments!",
                    videos = listOf(
                        PostVideo("tutorial1", "https://sample-videos.com/zip/10/mp4/SampleVideo_1280x720_1mb.mp4", "https://images.unsplash.com/photo-1586880244386-c398dc1b6f9a?w=400", "2:15", "Gallery wall tutorial")
                    ),
                    hashtags = listOf("#tutorial", "#diy", "#gallerywall", "#homedecor", "#stepbystep", "#interior")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "4 days ago"
            ),

            // ========== INTERACTIVE CONTENT (6) ==========

            // 15. POLL - Voting posts
            Post(
                id = "poll_1",
                type = PostType.POLL,
                user = PostUser("user15", "decision_maker", "Sam Polls", "https://images.unsplash.com/photo-1494790108755-2616b332c2c2?w=100", false, 3450, false),
                content = PostContent(
                    type = PostContentType.POLL,
                    text = "Help me decide my next vacation destination! 🌍 Where would you go?",
                    poll = PostPoll(
                        question = "Next vacation destination?",
                        options = listOf("Bali, Indonesia 🏝️", "Iceland 🏔️", "Japan 🗾", "Peru 🏛️"),
                        votes = mapOf(0 to 234, 1 to 156, 2 to 345, 3 to 123),
                        expiresAt = "2024-12-31T23:59:59Z"
                    ),
                    hashtags = listOf("#poll", "#travel", "#vacation", "#wanderlust", "#help")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "8 hours ago"
            ),

            // 16. QUIZ - Trivia, personality tests
            Post(
                id = "quiz_1",
                type = PostType.QUIZ,
                user = PostUser("user16", "trivia_master", "Quiz Master 🧠", "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100", true, 23400, true),
                content = PostContent(
                    type = PostContentType.TEXT,
                    text = "🧠 BRAIN TEASER: I have cities, but no houses. I have mountains, but no trees. I have water, but no fish. What am I? First correct answer gets a shoutout! 🏆",
                    hashtags = listOf("#quiz", "#brainteaser", "#riddle", "#challenge", "#thinking", "#puzzle")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "5 hours ago"
            ),

            // 17. SURVEY - Feedback collection
            Post(
                id = "survey_1",
                type = PostType.SURVEY,
                user = PostUser("user17", "feedback_collector", "Research Pro 📊", "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=100", false, 5670, false),
                content = PostContent(
                    type = PostContentType.TEXT,
                    text = "📊 Quick survey for my thesis research! What motivates you most in your daily work? Your responses help shape future workplace policies 💼",
                    hashtags = listOf("#survey", "#research", "#workplace", "#motivation", "#thesis", "#feedback")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "1 day ago"
            ),

            // 18. Q_AND_A - Ask me anything
            Post(
                id = "qna_1",
                type = PostType.Q_AND_A,
                user = PostUser("user18", "startup_founder", "Tech Founder 🚀", "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100", true, 89400, false),
                content = PostContent(
                    type = PostContentType.TEXT,
                    text = "🚀 AMA: Just closed our Series A funding! Ask me anything about building a startup, raising capital, or the entrepreneurial journey. I'll answer top questions in a live session tomorrow 💡",
                    hashtags = listOf("#ama", "#startup", "#seriesA", "#entrepreneur", "#funding", "#askme")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "12 hours ago"
            ),

            // 19. CHALLENGE - Viral challenges/trends
            Post(
                id = "challenge_1",
                type = PostType.CHALLENGE,
                user = PostUser("user19", "fitness_coach", "Fit Coach 💪", "https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=100", true, 156000, true),
                content = PostContent(
                    type = PostContentType.VIDEO,
                    text = "💪 #30DayPlankChallenge starts Monday! Who's joining me? Let's build that core strength together! Tag 3 friends to join 🔥",
                    videos = listOf(
                        PostVideo("challenge1", "https://sample-videos.com/zip/10/mp4/SampleVideo_1280x720_1mb.mp4", "https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=400", "0:45", "Plank challenge demo")
                    ),
                    hashtags = listOf("#30DayPlankChallenge", "#fitness", "#challenge", "#core", "#strength", "#together")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "6 hours ago"
            ),

            // 20. PETITION - Social causes
            Post(
                id = "petition_1",
                type = PostType.PETITION,
                user = PostUser("user20", "eco_warrior", "Green Activist 🌱", "https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=100", true, 45600, false),
                content = PostContent(
                    type = PostContentType.TEXT,
                    text = "🌱 SIGN THE PETITION: Save our local park from development! This green space is crucial for our community's wellbeing and biodiversity. Every signature counts! Link in bio 📝",
                    hashtags = listOf("#petition", "#saveourpark", "#environment", "#community", "#greenspace", "#activism")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "2 days ago"
            ),

            // ========== SOCIAL & LOCATION (8) ==========

            // 21. CHECK_IN - Location-based posts
            Post(
                id = "checkin_1",
                type = PostType.CHECK_IN,
                user = PostUser("user21", "local_explorer", "City Explorer 📍", "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=100", false, 8900, true),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "Finally tried this hidden gem everyone's been talking about! ✨ The vibe is incredible and the view is unmatched 📍",
                    images = listOf(
                        PostImage("checkin1", "https://images.unsplash.com/photo-1559847844-5315695dadae?w=800", "Rooftop view", 1.2f)
                    ),
                    hashtags = listOf("#checkin", "#hiddengem", "#rooftopbar", "#view", "#vibes", "#local"),
                    location = "Sky Bar, Downtown"
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                metadata = PostMetadata(location = PostLocation("Sky Bar", "Downtown Plaza", 40.7589, -73.9851, "New York", "USA", LocationCategory.ENTERTAINMENT)),
                createdAt = "4 hours ago"
            ),

            // 22. TRAVEL_LOG - Trip updates/itinerary
            Post(
                id = "travel_1",
                type = PostType.TRAVEL_LOG,
                user = PostUser("user22", "backpacker_soul", "Solo Traveler 🎒", "https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=100", true, 67800, false),
                content = PostContent(
                    type = PostContentType.MIXED,
                    text = "Day 5 in Southeast Asia! 🌏 Today's adventure: sunrise at Angkor Wat → local cooking class → night market exploration. This trip is teaching me so much about different cultures and myself 🙏",
                    images = listOf(
                        PostImage("travel1", "https://images.unsplash.com/photo-1578662996442-48f60103fc96?w=800", "Angkor Wat sunrise", 1.3f),
                        PostImage("travel2", "https://images.unsplash.com/photo-1555217851-6141535bd771?w=800", "Cooking class", 1.0f)
                    ),
                    hashtags = listOf("#travellog", "#southeastasia", "#angkorwat", "#solotravel", "#backpacking", "#culture"),
                    location = "Siem Reap, Cambodia"
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "1 day ago"
            ),

            // 23. LIFE_EVENT - Major life moments
            Post(
                id = "life_1",
                type = PostType.LIFE_EVENT,
                user = PostUser("user23", "new_graduate", "Jessica Smith", "https://images.unsplash.com/photo-1494790108755-2616b332c2c2?w=100", false, 890, false),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "WE DID IT! 🎓 After 4 years of late nights, countless coffee cups, and amazing friendships, I'm officially a Computer Science graduate! Thank you to everyone who supported me on this journey 💙",
                    images = listOf(
                        PostImage("grad1", "https://images.unsplash.com/photo-1523050854058-8df90110c9d1?w=800", "Graduation ceremony", 1.0f)
                    ),
                    hashtags = listOf("#graduation", "#computerscience", "#milestone", "#grateful", "#newbeginnings", "#classof2024")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.FRIENDS,
                createdAt = "3 days ago"
            ),

            // 24. MILESTONE - Personal achievements
            Post(
                id = "milestone_1",
                type = PostType.MILESTONE,
                user = PostUser("user24", "marathon_runner", "Run Fast 🏃‍♀️", "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=100", false, 5670, true),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "26.2 miles DONE! ✅ My first marathon is in the books! 6 months of training, 3:45:22 finish time, and the most incredible feeling of accomplishment. Next goal: Boston qualifier! 🏃‍♀️💪",
                    images = listOf(
                        PostImage("marathon1", "https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=800", "Marathon finish line", 1.0f)
                    ),
                    hashtags = listOf("#marathon", "#firstmarathon", "#milestone", "#running", "#achievement", "#bostonqualifier")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "2 days ago"
            ),

            // 25. MEMORY - Throwback/flashback posts
            Post(
                id = "memory_1",
                type = PostType.MEMORY,
                user = PostUser("user25", "nostalgic_soul", "Memory Lane", "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100", false, 2340, false),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "One year ago today 📅 This photo captures one of the best days of my life! Sometimes I need to remind myself how far I've come 🌟 #ThrowbackThursday",
                    images = listOf(
                        PostImage("memory1", "https://images.unsplash.com/photo-1469474968028-56623f02e42e?w=800", "Mountain summit", 1.2f)
                    ),
                    hashtags = listOf("#memory", "#oneyearago", "#throwbackthursday", "#grateful", "#progress", "#reflection")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "1 day ago"
            ),

            // 26. ANNIVERSARY - Yearly memories
            Post(
                id = "anniversary_1",
                type = PostType.ANNIVERSARY,
                user = PostUser("user26", "business_owner", "Café Owner ☕", "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=100", true, 12300, false),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "5 YEARS! ☕ Today marks 5 incredible years since we opened our little café! From serving 10 customers a day to being the neighborhood's favorite spot. Thank you to our amazing community! 🙏❤️",
                    images = listOf(
                        PostImage("anniversary1", "https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=800", "Café interior", 1.0f)
                    ),
                    hashtags = listOf("#anniversary", "#5years", "#café", "#community", "#grateful", "#milestone")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "6 hours ago"
            ),

            // 27. RECOMMENDATION - Place/product suggestions
            Post(
                id = "recommendation_1",
                type = PostType.RECOMMENDATION,
                user = PostUser("user27", "book_lover", "Page Turner 📚", "https://images.unsplash.com/photo-1494790108755-2616b332c2c2?w=100", true, 15600, true),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "📚 BOOK RECOMMENDATION: Just finished 'The Seven Husbands of Evelyn Hugo' and WOW! If you love stories about ambition, love, and Hollywood secrets, this is for you! Have you read it? No spoilers please! 🤐✨",
                    images = listOf(
                        PostImage("book1", "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=800", "Book cover", 0.75f)
                    ),
                    hashtags = listOf("#bookrecommendation", "#evelynhugo", "#reading", "#bookstagram", "#mustread", "#pageturner")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "8 hours ago"
            ),

            // 28. GROUP_ACTIVITY - Group-specific content
            Post(
                id = "group_1",
                type = PostType.GROUP_ACTIVITY,
                user = PostUser("user28", "hiking_group", "Mountain Hikers 🏔️", "https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=100", false, 3450, false),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "🏔️ HIKING GROUP UPDATE: This Saturday's trail has been changed to Eagle Peak due to weather conditions. Meet at 7 AM sharp at the usual parking spot. Bring extra layers! See everyone there! 🥾",
                    images = listOf(
                        PostImage("hike1", "https://images.unsplash.com/photo-1469474968028-56623f02e42e?w=800", "Eagle Peak trail", 1.3f)
                    ),
                    hashtags = listOf("#hikinggroup", "#eaglepeak", "#saturday", "#weather", "#layers", "#adventure")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.CUSTOM,
                createdAt = "12 hours ago"
            ),

            // ========== COMMERCIAL & BUSINESS (6) ==========

            // 29. PRODUCT_SHOWCASE - Selling items
            Post(
                id = "product_1",
                type = PostType.PRODUCT_SHOWCASE,
                user = PostUser("user29", "handmade_crafts", "Artisan Crafts 🎨", "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=100", true, 23400, false),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "✨ NEW ARRIVAL: Hand-woven macramé wall hangings! Each piece is unique and made with love using sustainable materials. Perfect for adding boho vibes to any space 🌿 DM for pricing!",
                    images = listOf(
                        PostImage("product1", "https://images.unsplash.com/photo-1586880244386-c398dc1b6f9a?w=800", "Macramé wall hanging", 0.8f),
                        PostImage("product2", "https://images.unsplash.com/photo-1615800001374-d4a8e1b4b6c2?w=800", "Room setup", 1.2f)
                    ),
                    hashtags = listOf("#handmade", "#macrame", "#wallart", "#boho", "#sustainable", "#artisan")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                metadata = PostMetadata(
                    targetType = "product",
                    targetId = "macrame_001",
                    targetName = "Macramé Wall Hanging",
                    price = "$85-$120",
                    isPromoted = true
                ),
                createdAt = "2 hours ago"
            ),

            // 30. SERVICE_PROMOTION - Business services
            Post(
                id = "service_1",
                type = PostType.SERVICE_PROMOTION,
                user = PostUser("user30", "yoga_instructor", "Zen Yoga 🧘‍♀️", "https://images.unsplash.com/photo-1494790108755-2616b332c2c2?w=100", true, 12300, true),
                content = PostContent(
                    type = PostContentType.VIDEO,
                    text = "🧘‍♀️ NEW ONLINE CLASSES STARTING MONDAY! Join me for 'Morning Flow' sessions designed to energize your day. First week FREE for new students! Link in bio to book your spot 🌅",
                    videos = listOf(
                        PostVideo("service1", "https://sample-videos.com/zip/10/mp4/SampleVideo_1280x720_1mb.mp4", "https://images.unsplash.com/photo-1506629905645-b178f3d67b5b?w=400", "1:30", "Morning yoga flow preview")
                    ),
                    hashtags = listOf("#yogaclasses", "#morningflow", "#onlineyoga", "#wellness", "#meditation", "#firstweekfree")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                metadata = PostMetadata(
                    targetType = "service",
                    targetId = "yoga_classes_001",
                    targetName = "Morning Flow Yoga Classes",
                    isPromoted = true
                ),
                createdAt = "1 day ago"
            ),

            // 31. EVENT_PROMOTION - Events/meetups
            Post(
                id = "event_1",
                type = PostType.EVENT_PROMOTION,
                user = PostUser("user31", "event_organizer", "City Events 🎉", "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100", false, 8900, false),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "🎵 SUMMER MUSIC FESTIVAL 2024! July 15-17 at Riverside Park. 3 days of incredible music, local food trucks, and good vibes! Early bird tickets available until June 1st 🎸✨",
                    images = listOf(
                        PostImage("event1", "https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=800", "Music festival poster", 1.0f),
                        PostImage("event2", "https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=800", "Previous year crowd", 1.5f)
                    ),
                    hashtags = listOf("#summermusicfest", "#festival", "#livemusic", "#riverside", "#earlybird", "#july2024"),
                    location = "Riverside Park"
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                metadata = PostMetadata(
                    targetType = "event",
                    targetId = "summer_fest_2024",
                    targetName = "Summer Music Festival 2024",
                    isPromoted = true
                ),
                createdAt = "3 days ago"
            ),

            // 32. JOB_POSTING - Hiring/career opportunities
            Post(
                id = "job_1",
                type = PostType.JOB_POSTING,
                user = PostUser("user32", "tech_startup", "InnovateTech 💼", "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100", true, 45600, false),
                content = PostContent(
                    type = PostContentType.TEXT,
                    text = "🚀 WE'RE HIRING: Senior Flutter Developer to join our growing team! Remote-first company, competitive salary, great benefits, and the chance to work on cutting-edge mobile apps. Apply through the link in our bio! #NowHiring",
                    hashtags = listOf("#hiring", "#flutter", "#remotework", "#developer", "#seniordeveloper", "#techcareer")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                metadata = PostMetadata(
                    targetType = "job",
                    targetId = "flutter_dev_001",
                    targetName = "Senior Flutter Developer",
                    isPromoted = true
                ),
                createdAt = "2 days ago"
            ),

            // 33. FUNDRAISER - Charity/personal causes
            Post(
                id = "fundraiser_1",
                type = PostType.FUNDRAISER,
                user = PostUser("user33", "charity_organizer", "Hearts for Hope ❤️", "https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=100", true, 34500, false),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "❤️ HELP US REACH OUR GOAL! We're 70% of the way to providing clean water access to 3 remote villages. Every donation, no matter the size, makes a real difference. Link in bio to contribute 🚰💙",
                    images = listOf(
                        PostImage("fundraiser1", "https://images.unsplash.com/photo-1578662996442-48f60103fc96?w=800", "Village water project", 1.0f)
                    ),
                    hashtags = listOf("#fundraiser", "#cleanwater", "#charity", "#help", "#community", "#donate")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                metadata = PostMetadata(
                    targetType = "fundraiser",
                    targetId = "clean_water_001",
                    targetName = "Clean Water Access Project",
                    isPromoted = false
                ),
                createdAt = "4 hours ago"
            ),

            // 34. COLLABORATION - Creative projects
            Post(
                id = "collab_1",
                type = PostType.COLLABORATION,
                user = PostUser("user34", "indie_filmmaker", "Indie Films 🎬", "https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=100", false, 8900, true),
                content = PostContent(
                    type = PostContentType.VIDEO,
                    text = "🎬 CALLING ALL CREATIVES! Looking for a talented composer to create an original soundtrack for my upcoming short film about urban life. This is a passion project with potential for festival submissions. DM me your samples! 🎵",
                    videos = listOf(
                        PostVideo("collab1", "https://sample-videos.com/zip/10/mp4/SampleVideo_1280x720_1mb.mp4", "https://images.unsplash.com/photo-1514565131-fce0801e5785?w=400", "1:45", "Film teaser")
                    ),
                    hashtags = listOf("#collaboration", "#filmmaker", "#composer", "#soundtrack", "#shortfilm", "#creative")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "6 hours ago"
            ),

            // ========== SPECIALIZED CONTENT (8) ==========

            // 35. RECIPE - Cooking/food content
            Post(
                id = "recipe_1",
                type = PostType.RECIPE,
                user = PostUser("user35", "home_chef", "Kitchen Magic 👩‍🍳", "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=100", true, 56700, true),
                content = PostContent(
                    type = PostContentType.MIXED,
                    text = "🍰 VIRAL CHOCOLATE LAVA CAKE that's ready in 15 minutes! Perfect for when you need dessert RIGHT NOW. Recipe in comments - you probably have all ingredients already! Who's making this tonight? 🤤",
                    images = listOf(
                        PostImage("recipe1", "https://images.unsplash.com/photo-1578985545062-69928b1d9587?w=800", "Chocolate lava cake", 1.0f),
                        PostImage("recipe2", "https://images.unsplash.com/photo-1571115764595-644a1f56a55c?w=800", "Ingredients layout", 1.2f)
                    ),
                    hashtags = listOf("#recipe", "#chocolatelavacake", "#15minutes", "#dessert", "#homemade", "#viral")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "3 hours ago"
            ),

            // 36. WORKOUT - Fitness routines
            Post(
                id = "workout_1",
                type = PostType.WORKOUT,
                user = PostUser("user36", "fitness_trainer", "Strong & Fit 💪", "https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=100", true, 78900, false),
                content = PostContent(
                    type = PostContentType.VIDEO,
                    text = "💪 10-MINUTE MORNING ENERGY BOOST! No equipment needed - just your body and determination. Perfect for busy mornings when you want to feel energized all day. Save this for tomorrow! ⚡",
                    videos = listOf(
                        PostVideo("workout1", "https://sample-videos.com/zip/10/mp4/SampleVideo_1280x720_1mb.mp4", "https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=400", "10:15", "Morning workout routine")
                    ),
                    hashtags = listOf("#workout", "#10minutes", "#morningworkout", "#noequipment", "#energy", "#fitness")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "1 day ago"
            ),

            // 37. BOOK_REVIEW - Reading updates
            Post(
                id = "book_review_1",
                type = PostType.BOOK_REVIEW,
                user = PostUser("user37", "avid_reader", "Book Worm 📖", "https://images.unsplash.com/photo-1494790108755-2616b332c2c2?w=100", false, 12300, true),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "📖 JUST FINISHED: 'Project Hail Mary' by Andy Weir and my mind is BLOWN! 🤯 If you love science, humor, and unexpected friendships, this is your next read. No spoilers, but I ugly-cried at the end. Rating: ⭐⭐⭐⭐⭐",
                    images = listOf(
                        PostImage("book1", "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=800", "Project Hail Mary book", 0.75f)
                    ),
                    hashtags = listOf("#bookreview", "#projecthailmary", "#andyweir", "#scifi", "#mustread", "#5stars")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                metadata = PostMetadata(
                    targetType = "book",
                    targetId = "project_hail_mary",
                    targetName = "Project Hail Mary",
                    rating = 5.0f
                ),
                createdAt = "2 days ago"
            ),

            // 38. MOOD_UPDATE - Emotional status
            Post(
                id = "mood_1",
                type = PostType.MOOD_UPDATE,
                user = PostUser("user38", "emotional_soul", "Feeling Human", "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=100", false, 2340, false),
                content = PostContent(
                    type = PostContentType.TEXT,
                    text = "Feeling grateful today 💛 Sometimes life throws curveballs, but I'm learning to appreciate the small moments that make everything worthwhile. Hope everyone is taking care of their mental health 🌱",
                    hashtags = listOf("#moodupdate", "#grateful", "#mentalhealth", "#selfcare", "#smallmoments", "#wellbeing")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                metadata = PostMetadata(
                    mood = MoodType.GRATEFUL,
                    feeling = FeelingType.GOOD
                ),
                createdAt = "5 hours ago"
            ),

            // 39. ACHIEVEMENT - Gaming/app achievements
            Post(
                id = "achievement_1",
                type = PostType.ACHIEVEMENT,
                user = PostUser("user39", "pro_gamer", "Gaming Pro 🏆", "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100", true, 89400, true),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "🏆 FINALLY DID IT! After 200+ hours, I've reached Grandmaster in Competitive mode! The grind was real but so worth it. Thank you to everyone who believed in me! Next stop: Tournament play 🎮⚡",
                    images = listOf(
                        PostImage("achievement1", "https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=800", "Grandmaster rank screen", 1.0f)
                    ),
                    hashtags = listOf("#achievement", "#grandmaster", "#competitive", "#gaming", "#200hours", "#tournament")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                metadata = PostMetadata(
                    targetType = "game_achievement",
                    targetId = "grandmaster_001",
                    targetName = "Grandmaster Rank"
                ),
                createdAt = "4 hours ago"
            ),

            // 40. QUOTE - Inspirational quotes
            Post(
                id = "quote_1",
                type = PostType.QUOTE,
                user = PostUser("user40", "daily_inspiration", "Daily Wisdom ✨", "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=100", false, 23400, false),
                content = PostContent(
                    type = PostContentType.IMAGE,
                    text = "✨ Monday Motivation ✨\n\n\"The only way to do great work is to love what you do.\" - Steve Jobs\n\nReminder: Your passion is your power. What you love will lead you to where you need to be 💫",
                    images = listOf(
                        PostImage("quote1", "https://images.unsplash.com/photo-1586880244386-c398dc1b6f9a?w=800", "Motivational quote design", 1.0f)
                    ),
                    hashtags = listOf("#quote", "#mondayMotivation", "#stevejobs", "#passion", "#inspiration", "#wisdom")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                createdAt = "8 hours ago"
            ),

            // 41. MUSIC - Music sharing/streaming
            Post(
                id = "music_1",
                type = PostType.MUSIC,
                user = PostUser("user41", "melody_maker", "Tune Sharer 🎵", "https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=100", true, 45600, true),
                content = PostContent(
                    type = PostContentType.TEXT,
                    text = "🎵 CURRENTLY OBSESSED with this track! 'Midnight City' by M83 - the synth work is absolutely incredible. Perfect for night drives or creative sessions. What's your current song obsession? Drop it below! 🌃",
                    hashtags = listOf("#music", "#midnightcity", "#m83", "#synth", "#obsessed", "#nightdrive")
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                metadata = PostMetadata(
                    targetType = "song",
                    targetId = "midnight_city_m83",
                    targetName = "Midnight City - M83"
                ),
                createdAt = "6 hours ago"
            ),

            // 42. VENUE - Venue information/reviews
            Post(
                id = "venue_1",
                type = PostType.VENUE,
                user = PostUser("user42", "venue_scout", "Event Finder 📍", "https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=100", false, 12300, false),
                content = PostContent(
                    type = PostContentType.MIXED,
                    text = "📍 VENUE SPOTLIGHT: The Rooftop Garden! Perfect for intimate events, amazing city views, and the staff is incredibly professional. Capacity: 80 people. Great for weddings, corporate events, or celebrations! 🌿✨",
                    images = listOf(
                        PostImage("venue1", "https://images.unsplash.com/photo-1559847844-5315695dadae?w=800", "Rooftop garden setup", 1.3f),
                        PostImage("venue2", "https://images.unsplash.com/photo-1514565131-fce0801e5785?w=800", "City view from venue", 1.5f)
                    ),
                    hashtags = listOf("#venue", "#rooftopgarden", "#events", "#wedding", "#corporate", "#cityview"),
                    location = "The Rooftop Garden, Downtown"
                ),
                interactions = PostInteractions(),
                privacy = PostPrivacy.PUBLIC,
                metadata = PostMetadata(
                    targetType = "venue",
                    targetId = "rooftop_garden_001",
                    targetName = "The Rooftop Garden",
                    rating = 4.7f,
                    location = PostLocation("The Rooftop Garden", "Downtown Plaza", 40.7589, -73.9851, "New York", "USA", LocationCategory.ENTERTAINMENT)
                ),
                createdAt = "1 day ago"
            )
        )
    }

    private fun generateMockComments(): List<PostComment> {
        // TODO: Implement proper mock comments with correct constructor parameters
        return emptyList()
        /*
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
        */
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

    // ========== HELPER FUNCTIONS FOR POST TYPE EXPLORATION ==========

    /**
     * Get posts filtered by category
     */
    suspend fun getPostsByCategory(category: String): Result<List<Post>> {
        delay(200)

        val categoryPosts = when(category.lowercase()) {
            "core" -> posts.filter { it.type in listOf(
                PostType.TEXT, PostType.IMAGE, PostType.VIDEO, PostType.AUDIO,
                PostType.LINK_SHARE, PostType.POST_MESSAGE, PostType.REVIEW, PostType.ALBUM
            )}
            "rich_media" -> posts.filter { it.type in listOf(
                PostType.STORY, PostType.REEL, PostType.LIVE_STREAM,
                PostType.PLAYLIST, PostType.MOOD_BOARD, PostType.TUTORIAL
            )}
            "interactive" -> posts.filter { it.type in listOf(
                PostType.POLL, PostType.QUIZ, PostType.SURVEY,
                PostType.Q_AND_A, PostType.CHALLENGE, PostType.PETITION
            )}
            "social" -> posts.filter { it.type in listOf(
                PostType.CHECK_IN, PostType.TRAVEL_LOG, PostType.LIFE_EVENT,
                PostType.MILESTONE, PostType.MEMORY, PostType.ANNIVERSARY,
                PostType.RECOMMENDATION, PostType.GROUP_ACTIVITY
            )}
            "commercial" -> posts.filter { it.type in listOf(
                PostType.PRODUCT_SHOWCASE, PostType.SERVICE_PROMOTION, PostType.EVENT_PROMOTION,
                PostType.JOB_POSTING, PostType.FUNDRAISER, PostType.COLLABORATION
            )}
            "specialized" -> posts.filter { it.type in listOf(
                PostType.RECIPE, PostType.WORKOUT, PostType.BOOK_REVIEW,
                PostType.MOOD_UPDATE, PostType.ACHIEVEMENT, PostType.QUOTE,
                PostType.MUSIC, PostType.VENUE
            )}
            else -> emptyList()
        }

        return Result.success(categoryPosts.sortedByDescending { it.createdAt })
    }

    /**
     * Get all post types with sample counts
     */
    suspend fun getAllPostTypesWithCounts(): Result<Map<PostType, Int>> {
        delay(100)

        val typeCounts = posts.groupBy { it.type }.mapValues { it.value.size }
        return Result.success(typeCounts)
    }

    /**
     * Get posts with specific engagement features (reactions, comments, etc.)
     */
    suspend fun getPostsWithHighEngagement(minLikes: Int = 1000): Result<List<Post>> {
        delay(150)

        val highEngagementPosts = posts.filter { post ->
            post.interactions.reactions.count { it.type == ReactionType.LIKE } >= minLikes ||
            post.interactions.comments.size >= 50 ||
            post.interactions.shares.size >= 100
        }.sortedByDescending { it.interactions.reactions.count { r -> r.type == ReactionType.LIKE } }

        return Result.success(highEngagementPosts)
    }

    /**
     * Get posts by privacy level
     */
    suspend fun getPostsByPrivacy(privacy: PostPrivacy): Result<List<Post>> {
        delay(100)

        val privacyPosts = posts.filter { it.privacy == privacy }
            .sortedByDescending { it.createdAt }

        return Result.success(privacyPosts)
    }

    /**
     * Get posts with media content (images/videos)
     */
    suspend fun getMediaPosts(): Result<List<Post>> {
        delay(150)

        val mediaPosts = posts.filter { post ->
            post.content.images.isNotEmpty() || post.content.videos.isNotEmpty()
        }.sortedByDescending { it.createdAt }

        return Result.success(mediaPosts)
    }

    /**
     * Get posts from verified users
     */
    suspend fun getVerifiedUserPosts(): Result<List<Post>> {
        delay(100)

        val verifiedPosts = posts.filter { it.user.isVerified }
            .sortedByDescending { it.interactions.likes }

        return Result.success(verifiedPosts)
    }

    /**
     * Get sample posts for demo/testing purposes
     * Returns one post of each type for comprehensive testing
     */
    suspend fun getAllPostTypesSample(): Result<Map<PostType, Post>> {
        delay(200)

        val samplePosts = posts.groupBy { it.type }
            .mapValues { it.value.first() }

        return Result.success(samplePosts)
    }
}