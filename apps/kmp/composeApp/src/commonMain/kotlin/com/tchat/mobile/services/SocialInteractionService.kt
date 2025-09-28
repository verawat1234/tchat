package com.tchat.mobile.services

import com.tchat.mobile.models.*
import com.tchat.mobile.repositories.SocialRepository
import com.tchat.mobile.utils.PlatformUtils
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow

/**
 * Service for handling social interactions
 * Contains all business logic for likes, follows, bookmarks, etc.
 */
class SocialInteractionService(
    private val socialRepository: SocialRepository,
    private val currentUserId: String // Should come from authentication service
) {

    // Local state for real-time UI updates
    private val _interactionStates = MutableStateFlow<Map<String, UserInteractionState>>(emptyMap())
    val interactionStates: StateFlow<Map<String, UserInteractionState>> = _interactionStates.asStateFlow()

    private val _interactionCounts = MutableStateFlow<Map<String, InteractionCounts>>(emptyMap())
    val interactionCounts: StateFlow<Map<String, InteractionCounts>> = _interactionCounts.asStateFlow()

    // Like operations
    suspend fun likePost(postId: String): Result<Boolean> {
        return try {
            val interaction = SocialInteraction(
                id = "${currentUserId}_${postId}_like_${PlatformUtils.currentTimeMillis()}",
                userId = currentUserId,
                targetId = postId,
                targetType = InteractionTargetType.POST,
                interactionType = InteractionType.LIKE,
                createdAt = System.currentTimeMillis(),
                updatedAt = System.currentTimeMillis()
            )

            val result = socialRepository.addInteraction(interaction)

            if (result.isSuccess) {
                // Update local state for immediate UI feedback
                updateLocalInteractionState(postId, InteractionType.LIKE, true)
                refreshInteractionCounts(postId, InteractionTargetType.POST)
            }

            result.map { true }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun unlikePost(postId: String): Result<Boolean> {
        return try {
            val result = socialRepository.removeInteraction(
                userId = currentUserId,
                targetId = postId,
                targetType = InteractionTargetType.POST,
                interactionType = InteractionType.LIKE
            )

            if (result.isSuccess) {
                updateLocalInteractionState(postId, InteractionType.LIKE, false)
                refreshInteractionCounts(postId, InteractionTargetType.POST)
            }

            result.map { false }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun toggleLike(postId: String): Result<Boolean> {
        return try {
            val currentState = getCurrentInteractionState(postId)
            if (currentState.isLiked) {
                unlikePost(postId)
            } else {
                likePost(postId)
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Bookmark operations
    suspend fun bookmarkPost(postId: String): Result<Boolean> {
        return try {
            val interaction = SocialInteraction(
                id = "${currentUserId}_${postId}_bookmark_${PlatformUtils.currentTimeMillis()}",
                userId = currentUserId,
                targetId = postId,
                targetType = InteractionTargetType.POST,
                interactionType = InteractionType.BOOKMARK,
                createdAt = System.currentTimeMillis(),
                updatedAt = System.currentTimeMillis()
            )

            val result = socialRepository.addInteraction(interaction)

            if (result.isSuccess) {
                updateLocalInteractionState(postId, InteractionType.BOOKMARK, true)
            }

            result.map { true }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun removeBookmark(postId: String): Result<Boolean> {
        return try {
            val result = socialRepository.removeInteraction(
                userId = currentUserId,
                targetId = postId,
                targetType = InteractionTargetType.POST,
                interactionType = InteractionType.BOOKMARK
            )

            if (result.isSuccess) {
                updateLocalInteractionState(postId, InteractionType.BOOKMARK, false)
            }

            result.map { false }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun toggleBookmark(postId: String): Result<Boolean> {
        return try {
            val currentState = getCurrentInteractionState(postId)
            if (currentState.isBookmarked) {
                removeBookmark(postId)
            } else {
                bookmarkPost(postId)
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Follow operations
    suspend fun followUser(userId: String): Result<Boolean> {
        return try {
            val interaction = SocialInteraction(
                id = "${currentUserId}_${userId}_follow_${PlatformUtils.currentTimeMillis()}",
                userId = currentUserId,
                targetId = userId,
                targetType = InteractionTargetType.USER,
                interactionType = InteractionType.FOLLOW,
                createdAt = System.currentTimeMillis(),
                updatedAt = System.currentTimeMillis()
            )

            val result = socialRepository.addInteraction(interaction)

            if (result.isSuccess) {
                updateLocalInteractionState(userId, InteractionType.FOLLOW, true)
            }

            result.map { true }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun unfollowUser(userId: String): Result<Boolean> {
        return try {
            val result = socialRepository.removeInteraction(
                userId = currentUserId,
                targetId = userId,
                targetType = InteractionTargetType.USER,
                interactionType = InteractionType.FOLLOW
            )

            if (result.isSuccess) {
                updateLocalInteractionState(userId, InteractionType.FOLLOW, false)
            }

            result.map { false }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun toggleFollow(userId: String): Result<Boolean> {
        return try {
            val currentState = getCurrentInteractionState(userId)
            if (currentState.isFollowing) {
                unfollowUser(userId)
            } else {
                followUser(userId)
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Event attendance operations
    suspend fun attendEvent(eventId: String): Result<Boolean> {
        return try {
            val interaction = SocialInteraction(
                id = "${currentUserId}_${eventId}_attend_${PlatformUtils.currentTimeMillis()}",
                userId = currentUserId,
                targetId = eventId,
                targetType = InteractionTargetType.EVENT,
                interactionType = InteractionType.ATTEND,
                createdAt = System.currentTimeMillis(),
                updatedAt = System.currentTimeMillis()
            )

            val result = socialRepository.addInteraction(interaction)

            if (result.isSuccess) {
                updateLocalInteractionState(eventId, InteractionType.ATTEND, true)
                refreshInteractionCounts(eventId, InteractionTargetType.EVENT)
            }

            result.map { true }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun leaveEvent(eventId: String): Result<Boolean> {
        return try {
            val result = socialRepository.removeInteraction(
                userId = currentUserId,
                targetId = eventId,
                targetType = InteractionTargetType.EVENT,
                interactionType = InteractionType.ATTEND
            )

            if (result.isSuccess) {
                updateLocalInteractionState(eventId, InteractionType.ATTEND, false)
                refreshInteractionCounts(eventId, InteractionTargetType.EVENT)
            }

            result.map { false }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun toggleEventAttendance(eventId: String): Result<Boolean> {
        return try {
            val currentState = getCurrentInteractionState(eventId)
            if (currentState.isAttending) {
                leaveEvent(eventId)
            } else {
                attendEvent(eventId)
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // State management
    suspend fun loadInteractionState(targetId: String, targetType: InteractionTargetType): Result<UserInteractionState> {
        return try {
            val interactions = socialRepository.getUserInteractionState(currentUserId, targetId, targetType)

            interactions.map { interactionTypes ->
                val state = UserInteractionState(
                    isLiked = InteractionType.LIKE in interactionTypes,
                    isBookmarked = InteractionType.BOOKMARK in interactionTypes,
                    isFollowing = InteractionType.FOLLOW in interactionTypes,
                    isAttending = InteractionType.ATTEND in interactionTypes
                )

                // Update local state
                _interactionStates.value = _interactionStates.value + (targetId to state)

                state
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun loadInteractionCounts(targetId: String, targetType: InteractionTargetType): Result<InteractionCounts> {
        return try {
            val counts = socialRepository.getInteractionCounts(targetId, targetType)

            counts.onSuccess { interactionCounts ->
                _interactionCounts.value = _interactionCounts.value + (targetId to interactionCounts)
            }

            counts
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Bulk operations for performance
    suspend fun loadMultipleInteractionStates(targets: List<Pair<String, InteractionTargetType>>): Result<Map<String, UserInteractionState>> {
        return try {
            val states = mutableMapOf<String, UserInteractionState>()

            targets.forEach { (targetId, targetType) ->
                loadInteractionState(targetId, targetType).onSuccess { state ->
                    states[targetId] = state
                }
            }

            Result.success(states)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Helper methods
    private fun getCurrentInteractionState(targetId: String): UserInteractionState {
        return _interactionStates.value[targetId] ?: UserInteractionState()
    }

    private fun updateLocalInteractionState(targetId: String, interactionType: InteractionType, isActive: Boolean) {
        val currentState = getCurrentInteractionState(targetId)
        val newState = when (interactionType) {
            InteractionType.LIKE -> currentState.copy(isLiked = isActive)
            InteractionType.BOOKMARK -> currentState.copy(isBookmarked = isActive)
            InteractionType.FOLLOW -> currentState.copy(isFollowing = isActive)
            InteractionType.ATTEND -> currentState.copy(isAttending = isActive)
            else -> currentState
        }

        _interactionStates.value = _interactionStates.value + (targetId to newState)
    }

    private suspend fun refreshInteractionCounts(targetId: String, targetType: InteractionTargetType) {
        loadInteractionCounts(targetId, targetType)
    }

    // Public getters for current state
    fun getInteractionState(targetId: String): UserInteractionState {
        return getCurrentInteractionState(targetId)
    }

    fun getInteractionCounts(targetId: String): InteractionCounts {
        return _interactionCounts.value[targetId] ?: InteractionCounts()
    }
}