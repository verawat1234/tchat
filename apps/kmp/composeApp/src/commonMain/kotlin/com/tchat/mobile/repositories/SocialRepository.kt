package com.tchat.mobile.repositories

import com.tchat.mobile.models.*
// Explicit imports for social domain models
import com.tchat.mobile.models.Story
import com.tchat.mobile.models.SocialUserProfile
import com.tchat.mobile.models.Friend
import com.tchat.mobile.models.FriendshipStatus
import com.tchat.mobile.models.Event
import com.tchat.mobile.models.EventAttendanceStatus
import com.tchat.mobile.models.SocialInteraction
import com.tchat.mobile.models.InteractionType
import com.tchat.mobile.models.InteractionTargetType
import com.tchat.mobile.models.InteractionCounts
import com.tchat.mobile.models.SocialComment
import com.tchat.mobile.models.UserStats

/**
 * Repository interface for social data operations
 * Follows SQL sync architecture principles
 */
interface SocialRepository {

    // Stories operations
    suspend fun getStories(viewerId: String): Result<List<Story>>
    suspend fun getStoriesByAuthor(authorId: String, viewerId: String): Result<List<Story>>
    suspend fun getStoryById(storyId: String, viewerId: String): Result<Story?>
    suspend fun createStory(story: Story): Result<Story>
    suspend fun markStoryViewed(storyId: String, viewerId: String): Result<Unit>
    suspend fun deleteExpiredStories(): Result<Unit>

    // User profiles operations
    suspend fun getUserProfile(userId: String): Result<SocialUserProfile?>
    suspend fun getUserProfileByUsername(username: String): Result<SocialUserProfile?>
    suspend fun updateUserProfile(profile: SocialUserProfile): Result<SocialUserProfile>
    suspend fun updateUserOnlineStatus(userId: String, isOnline: Boolean): Result<Unit>

    // Friends operations
    suspend fun getFriends(userId: String, status: FriendshipStatus): Result<List<Friend>>
    suspend fun getPendingFriendRequests(userId: String): Result<List<Friend>>
    suspend fun getOnlineFriends(userId: String): Result<List<Friend>>
    suspend fun getFriendSuggestions(userId: String, limit: Int): Result<List<SocialUserProfile>>
    suspend fun sendFriendRequest(userId: String, targetUserId: String): Result<Unit>
    suspend fun acceptFriendRequest(userId: String, requesterId: String): Result<Unit>
    suspend fun rejectFriendRequest(userId: String, requesterId: String): Result<Unit>
    suspend fun removeFriend(userId: String, friendUserId: String): Result<Unit>
    suspend fun checkFriendship(userId: String, targetUserId: String): Result<FriendshipStatus?>

    // Events operations
    suspend fun getAllEvents(userId: String): Result<List<Event>>
    suspend fun getUpcomingEvents(userId: String): Result<List<Event>>
    suspend fun getEventsByCategory(category: String, userId: String): Result<List<Event>>
    suspend fun getEventById(eventId: String, userId: String): Result<Event?>
    suspend fun getUserEvents(userId: String): Result<List<Event>>
    suspend fun getEventAttendees(eventId: String): Result<List<SocialUserProfile>>
    suspend fun createEvent(event: Event): Result<Event>
    suspend fun updateEvent(event: Event): Result<Event>
    suspend fun deleteEvent(eventId: String, organizerId: String): Result<Unit>
    suspend fun rsvpToEvent(eventId: String, userId: String, status: EventAttendanceStatus): Result<Unit>
    suspend fun removeEventRsvp(eventId: String, userId: String): Result<Unit>
    suspend fun getEventCategories(): Result<List<Pair<String, Int>>>
    suspend fun searchEvents(query: String, userId: String): Result<List<Event>>

    // Social interactions operations
    suspend fun getInteractionsByUser(userId: String): Result<List<SocialInteraction>>
    suspend fun getInteractionsByTarget(targetId: String, targetType: InteractionTargetType): Result<List<SocialInteraction>>
    suspend fun getUserInteractionState(userId: String, targetId: String, targetType: InteractionTargetType): Result<Set<InteractionType>>
    suspend fun getInteractionCounts(targetId: String, targetType: InteractionTargetType): Result<InteractionCounts>
    suspend fun addInteraction(interaction: SocialInteraction): Result<Unit>
    suspend fun removeInteraction(userId: String, targetId: String, targetType: InteractionTargetType, interactionType: InteractionType): Result<Unit>
    suspend fun getFollowedUsers(userId: String): Result<List<SocialUserProfile>>
    suspend fun getFollowers(userId: String): Result<List<SocialUserProfile>>
    suspend fun getUserStats(userId: String): Result<UserStats>

    // Comments operations
    suspend fun getCommentsByTarget(targetId: String, targetType: InteractionTargetType, userId: String): Result<List<SocialComment>>
    suspend fun getCommentReplies(commentId: String, userId: String): Result<List<SocialComment>>
    suspend fun createComment(comment: SocialComment): Result<SocialComment>
    suspend fun updateComment(commentId: String, userId: String, content: String): Result<Unit>
    suspend fun deleteComment(commentId: String, userId: String): Result<Unit>

    // Analytics and recommendations
    suspend fun getPopularContent(since: Long, minInteractions: Int): Result<List<Triple<String, InteractionTargetType, Int>>>
    suspend fun getRecentActivity(userId: String, since: Long, limit: Int): Result<List<SocialInteraction>>
}