package com.tchat.mobile.services

import com.tchat.mobile.models.*
import com.tchat.mobile.repositories.SocialRepository
import com.tchat.mobile.utils.PlatformUtils
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow

/**
 * Service for handling social content operations
 * Contains business logic for stories, friends, events, and comments
 */
class SocialContentService(
    private val socialRepository: SocialRepository,
    private val currentUserId: String
) {

    // Local content cache for performance
    private val _stories = MutableStateFlow<List<Story>>(emptyList())
    val stories: StateFlow<List<Story>> = _stories.asStateFlow()

    private val _friends = MutableStateFlow<List<Friend>>(emptyList())
    val friends: StateFlow<List<Friend>> = _friends.asStateFlow()

    private val _events = MutableStateFlow<List<Event>>(emptyList())
    val events: StateFlow<List<Event>> = _events.asStateFlow()

    private val _userProfile = MutableStateFlow<SocialUserProfile?>(null)
    val userProfile: StateFlow<SocialUserProfile?> = _userProfile.asStateFlow()

    // Stories operations
    suspend fun getPersonalizedStories(): Result<List<Story>> {
        return try {
            val result = socialRepository.getStories(currentUserId)

            result.onSuccess { storiesList ->
                _stories.value = storiesList
            }

            result
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getStoriesByAuthor(authorId: String): Result<List<Story>> {
        return socialRepository.getStoriesByAuthor(authorId, currentUserId)
    }

    suspend fun createStory(content: String, preview: String = "", isLive: Boolean = false): Result<Story> {
        return try {
            val story = Story(
                id = "story_${currentUserId}_${PlatformUtils.currentTimeMillis()}",
                authorId = currentUserId,
                content = content,
                preview = preview.ifEmpty { content.take(50) },
                createdAt = System.currentTimeMillis(),
                expiresAt = System.currentTimeMillis() + (24 * 60 * 60 * 1000), // 24 hours
                isLive = isLive,
                viewCount = 0
            )

            val result = socialRepository.createStory(story)

            result.onSuccess {
                // Add to local cache
                _stories.value = listOf(story) + _stories.value
            }

            result
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun viewStory(storyId: String): Result<Unit> {
        return try {
            socialRepository.markStoryViewed(storyId, currentUserId)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Friends operations
    suspend fun getFriendsWithStatus(status: FriendshipStatus = FriendshipStatus.ACCEPTED): Result<List<Friend>> {
        return try {
            val result = socialRepository.getFriends(currentUserId, status)

            result.onSuccess { friendsList ->
                if (status == FriendshipStatus.ACCEPTED) {
                    _friends.value = friendsList
                }
            }

            result
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getOnlineFriends(): Result<List<Friend>> {
        return socialRepository.getOnlineFriends(currentUserId)
    }

    suspend fun getFriendSuggestions(limit: Int = 10): Result<List<SocialUserProfile>> {
        return socialRepository.getFriendSuggestions(currentUserId, limit)
    }

    suspend fun sendFriendRequest(targetUserId: String): Result<Unit> {
        return socialRepository.sendFriendRequest(currentUserId, targetUserId)
    }

    suspend fun acceptFriendRequest(requesterId: String): Result<Unit> {
        return try {
            val result = socialRepository.acceptFriendRequest(currentUserId, requesterId)

            result.onSuccess {
                // Refresh friends list
                getFriendsWithStatus()
            }

            result
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun rejectFriendRequest(requesterId: String): Result<Unit> {
        return socialRepository.rejectFriendRequest(currentUserId, requesterId)
    }

    suspend fun removeFriend(friendUserId: String): Result<Unit> {
        return try {
            val result = socialRepository.removeFriend(currentUserId, friendUserId)

            result.onSuccess {
                // Remove from local cache
                _friends.value = _friends.value.filter { it.friendUserId != friendUserId }
            }

            result
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun checkFriendshipStatus(targetUserId: String): Result<FriendshipStatus?> {
        return socialRepository.checkFriendship(currentUserId, targetUserId)
    }

    // Events operations
    suspend fun getUpcomingEvents(): Result<List<Event>> {
        return try {
            val result = socialRepository.getUpcomingEvents(currentUserId)

            result.onSuccess { eventsList ->
                _events.value = eventsList
            }

            result
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getAllEvents(): Result<List<Event>> {
        return socialRepository.getAllEvents(currentUserId)
    }

    suspend fun getEventsByCategory(category: String): Result<List<Event>> {
        return socialRepository.getEventsByCategory(category, currentUserId)
    }

    suspend fun getUserEvents(): Result<List<Event>> {
        return socialRepository.getUserEvents(currentUserId)
    }

    suspend fun searchEvents(query: String): Result<List<Event>> {
        return socialRepository.searchEvents(query, currentUserId)
    }

    suspend fun createEvent(
        title: String,
        description: String,
        eventDate: Long,
        location: String,
        price: String = "Free",
        category: String,
        imageUrl: String = "",
        maxAttendees: Int? = null,
        isPublic: Boolean = true
    ): Result<Event> {
        return try {
            val event = Event(
                id = "event_${currentUserId}_${PlatformUtils.currentTimeMillis()}",
                title = title,
                description = description,
                eventDate = eventDate,
                location = location,
                price = price,
                imageUrl = imageUrl,
                category = category,
                organizerId = currentUserId,
                attendeesCount = 0,
                maxAttendees = maxAttendees,
                isPublic = isPublic,
                createdAt = System.currentTimeMillis(),
                updatedAt = System.currentTimeMillis()
            )

            val result = socialRepository.createEvent(event)

            result.onSuccess {
                // Add to local cache
                _events.value = listOf(event) + _events.value
            }

            result
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Comments operations
    suspend fun getComments(targetId: String, targetType: InteractionTargetType): Result<List<SocialComment>> {
        return socialRepository.getCommentsByTarget(targetId, targetType, currentUserId)
    }

    suspend fun getCommentReplies(commentId: String): Result<List<SocialComment>> {
        return socialRepository.getCommentReplies(commentId, currentUserId)
    }

    suspend fun createComment(
        targetId: String,
        targetType: InteractionTargetType,
        content: String,
        parentCommentId: String? = null
    ): Result<SocialComment> {
        return try {
            val comment = SocialComment(
                id = "comment_${currentUserId}_${PlatformUtils.currentTimeMillis()}",
                targetId = targetId,
                targetType = targetType,
                userId = currentUserId,
                content = content,
                parentCommentId = parentCommentId,
                likesCount = 0,
                repliesCount = 0,
                createdAt = System.currentTimeMillis(),
                updatedAt = System.currentTimeMillis()
            )

            socialRepository.createComment(comment)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun updateComment(commentId: String, content: String): Result<Unit> {
        return socialRepository.updateComment(commentId, currentUserId, content)
    }

    suspend fun deleteComment(commentId: String): Result<Unit> {
        return socialRepository.deleteComment(commentId, currentUserId)
    }

    // User profile operations
    suspend fun getCurrentUserProfile(): Result<SocialUserProfile?> {
        return try {
            val result = socialRepository.getUserProfile(currentUserId)

            result.onSuccess { profile ->
                _userProfile.value = profile
            }

            result
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getUserProfile(userId: String): Result<SocialUserProfile?> {
        return socialRepository.getUserProfile(userId)
    }

    suspend fun getUserProfileByUsername(username: String): Result<SocialUserProfile?> {
        return socialRepository.getUserProfileByUsername(username)
    }

    suspend fun updateUserProfile(
        displayName: String,
        bio: String = "",
        avatarUrl: String = "",
        statusMessage: String = ""
    ): Result<SocialUserProfile> {
        return try {
            val currentProfile = _userProfile.value
                ?: return Result.failure(Exception("No current user profile found"))

            val updatedProfile = currentProfile.copy(
                displayName = displayName,
                bio = bio,
                avatarUrl = avatarUrl,
                statusMessage = statusMessage,
                updatedAt = System.currentTimeMillis()
            )

            val result = socialRepository.updateUserProfile(updatedProfile)

            result.onSuccess { profile ->
                _userProfile.value = profile
            }

            result
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun updateOnlineStatus(isOnline: Boolean): Result<Unit> {
        return try {
            val result = socialRepository.updateUserOnlineStatus(currentUserId, isOnline)

            result.onSuccess {
                _userProfile.value = _userProfile.value?.copy(
                    isOnline = isOnline,
                    lastSeen = if (!isOnline) System.currentTimeMillis() else null,
                    updatedAt = System.currentTimeMillis()
                )
            }

            result
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Analytics and insights
    suspend fun getUserStats(): Result<UserStats> {
        return socialRepository.getUserStats(currentUserId)
    }

    suspend fun getPopularContent(since: Long, minInteractions: Int = 10): Result<List<Triple<String, InteractionTargetType, Int>>> {
        return socialRepository.getPopularContent(since, minInteractions)
    }

    suspend fun getRecentActivity(since: Long, limit: Int = 50): Result<List<SocialInteraction>> {
        return socialRepository.getRecentActivity(currentUserId, since, limit)
    }

    // Cleanup operations
    suspend fun cleanupExpiredStories(): Result<Unit> {
        return try {
            val result = socialRepository.deleteExpiredStories()

            result.onSuccess {
                // Remove expired stories from local cache
                val currentTime = System.currentTimeMillis()
                _stories.value = _stories.value.filter { !it.isExpired }
            }

            result
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Batch operations for performance
    suspend fun refreshAllContent(): Result<Unit> {
        return try {
            // Run all refresh operations in parallel
            val storiesResult = getPersonalizedStories()
            val friendsResult = getFriendsWithStatus()
            val eventsResult = getUpcomingEvents()
            val profileResult = getCurrentUserProfile()

            // Check if any failed
            listOf(storiesResult, friendsResult, eventsResult, profileResult).forEach { result ->
                if (result.isFailure) {
                    return result.map { }
                }
            }

            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Additional social content methods
    suspend fun getSocialStories(): Result<List<Story>> {
        // Alias for getPersonalizedStories for compatibility
        return getPersonalizedStories()
    }

    // Legacy compatibility methods (for gradual migration)
    suspend fun getLegacyStories(): List<StoryItem> {
        return try {
            getPersonalizedStories().getOrNull()?.map { story ->
                // Convert domain model to legacy UI model
                val userProfile = getUserProfile(story.authorId).getOrNull()
                val authorProfile = userProfile ?: createFallbackProfile(story.authorId)
                StoryItem(
                    id = story.id,
                    author = UserItem(
                        id = authorProfile.userId,
                        name = authorProfile.displayName,
                        username = authorProfile.username,
                        avatar = authorProfile.avatarUrl,
                        isVerified = authorProfile.isVerified,
                        isOnline = authorProfile.isOnline,
                        lastSeen = authorProfile.lastSeen?.toString() ?: "",
                        mutualFriends = 0,
                        status = ""
                    ),
                    preview = story.preview,
                    content = story.content,
                    timestamp = story.createdAt.toString(),
                    isViewed = story.isViewed,
                    isLive = story.isLive,
                    expiresAt = story.expiresAt.toString()
                )
            } ?: emptyList()
        } catch (e: Exception) {
            emptyList()
        }
    }

    suspend fun getLegacyFriends(): List<FriendItem> {
        return try {
            getFriendsWithStatus().getOrNull()?.map { friend ->
                FriendItem(
                    id = friend.id,
                    name = friend.profile?.displayName ?: "Unknown",
                    username = friend.profile?.username ?: "",
                    avatar = friend.profile?.avatarUrl ?: "",
                    isOnline = friend.profile?.isOnline ?: false,
                    isFollowing = friend.status == FriendshipStatus.ACCEPTED,
                    mutualFriends = friend.mutualFriendsCount,
                    status = friend.status.name.lowercase()
                )
            } ?: emptyList()
        } catch (e: Exception) {
            emptyList()
        }
    }

    suspend fun getLegacyEvents(): List<EventItem> {
        return try {
            getUpcomingEvents().getOrNull()?.map { event ->
                EventItem(
                    id = event.id,
                    title = event.title,
                    description = event.description,
                    date = event.eventDate.toString(),
                    location = event.location,
                    price = event.price,
                    imageUrl = event.imageUrl,
                    attendeesCount = event.attendeesCount,
                    category = event.category,
                    isAttending = event.userAttendanceStatus == EventAttendanceStatus.ATTENDING
                )
            } ?: emptyList()
        } catch (e: Exception) {
            emptyList()
        }
    }

    private fun createFallbackProfile(userId: String): SocialUserProfile {
        return SocialUserProfile(
            userId = userId,
            displayName = "User $userId",
            username = "user_$userId",
            avatarUrl = "",
            bio = "",
            isVerified = false,
            isOnline = false,
            lastSeen = null,
            statusMessage = "",
            createdAt = System.currentTimeMillis(),
            updatedAt = System.currentTimeMillis()
        )
    }
}