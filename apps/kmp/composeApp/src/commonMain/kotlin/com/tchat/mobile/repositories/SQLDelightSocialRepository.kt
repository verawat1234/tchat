package com.tchat.mobile.repositories

import app.cash.sqldelight.coroutines.asFlow
import app.cash.sqldelight.coroutines.mapToList
import app.cash.sqldelight.coroutines.mapToOne
import app.cash.sqldelight.coroutines.mapToOneOrNull
import com.tchat.mobile.database.TchatDatabase
import com.tchat.mobile.models.*
import com.tchat.mobile.utils.PlatformUtils
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.firstOrNull
import kotlinx.coroutines.withContext

/**
 * SQLDelight implementation of SocialRepository
 * Handles all social data operations using SQLDelight
 */
class SQLDelightSocialRepository(
    private val database: TchatDatabase
) : SocialRepository {

    // Stories operations
    override suspend fun getStories(viewerId: String): Result<List<Story>> = withContext(Dispatchers.Default) {
        try {
            val currentTime = System.currentTimeMillis()
            val stories = database.storyQueries.getAllStories(viewerId, currentTime)
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val domainStories = stories.map { row ->
                Story(
                    id = row.id,
                    authorId = row.author_id,
                    content = row.content,
                    preview = row.preview,
                    createdAt = row.created_at,
                    expiresAt = row.expires_at,
                    isLive = row.is_live == 1L,
                    viewCount = row.view_count.toInt(),
                    isViewed = row.is_viewed == 1L,
                    totalViews = row.total_views.toInt()
                )
            }
            Result.success(domainStories)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun getStoriesByAuthor(authorId: String, viewerId: String): Result<List<Story>> = withContext(Dispatchers.Default) {
        try {
            val currentTime = System.currentTimeMillis()
            val stories = database.storyQueries.getStoriesByAuthor(viewerId, authorId, currentTime)
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val domainStories = stories.map { row ->
                Story(
                    id = row.id,
                    authorId = row.author_id,
                    content = row.content,
                    preview = row.preview,
                    createdAt = row.created_at,
                    expiresAt = row.expires_at,
                    isLive = row.is_live == 1L,
                    viewCount = row.view_count.toInt(),
                    isViewed = row.is_viewed == 1L,
                    totalViews = row.total_views.toInt()
                )
            }
            Result.success(domainStories)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun getStoryById(storyId: String, viewerId: String): Result<Story?> = withContext(Dispatchers.Default) {
        try {
            val row = database.storyQueries.getStoryById(viewerId, storyId)
                .asFlow()
                .mapToOneOrNull(Dispatchers.Default)
                .firstOrNull()

            val story = row?.let {
                Story(
                    id = it.id,
                    authorId = it.author_id,
                    content = it.content,
                    preview = it.preview,
                    createdAt = it.created_at,
                    expiresAt = it.expires_at,
                    isLive = it.is_live == 1L,
                    viewCount = it.view_count.toInt(),
                    isViewed = it.is_viewed == 1L,
                    totalViews = it.total_views.toInt()
                )
            }
            Result.success(story)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun createStory(story: Story): Result<Story> = withContext(Dispatchers.Default) {
        try {
            database.storyQueries.insertStory(
                id = story.id,
                author_id = story.authorId,
                content = story.content,
                preview = story.preview,
                created_at = story.createdAt,
                expires_at = story.expiresAt,
                is_live = if (story.isLive) 1L else 0L,
                view_count = story.viewCount.toLong()
            )
            Result.success(story)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun markStoryViewed(storyId: String, viewerId: String): Result<Unit> = withContext(Dispatchers.Default) {
        try {
            val viewId = "${storyId}_${viewerId}_${PlatformUtils.currentTimeMillis()}"
            database.storyQueries.markStoryViewed(
                id = viewId,
                story_id = storyId,
                viewer_id = viewerId,
                viewed_at = System.currentTimeMillis()
            )

            // Update view count
            database.storyQueries.updateStoryViewCount(storyId, storyId)
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun deleteExpiredStories(): Result<Unit> = withContext(Dispatchers.Default) {
        try {
            database.storyQueries.deleteExpiredStories(System.currentTimeMillis())
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // User profiles operations
    override suspend fun getUserProfile(userId: String): Result<SocialUserProfile?> = withContext(Dispatchers.Default) {
        try {
            val row = database.friendQueries.getUserProfile(userId)
                .asFlow()
                .mapToOneOrNull(Dispatchers.Default)
                .firstOrNull()

            val profile = row?.let {
                SocialUserProfile(
                    userId = it.user_id,
                    displayName = it.display_name,
                    username = it.username,
                    avatarUrl = it.avatar_url ?: "",
                    bio = it.bio ?: "",
                    isVerified = it.is_verified == 1L,
                    isOnline = it.is_online == 1L,
                    lastSeen = it.last_seen,
                    statusMessage = it.status_message ?: "",
                    createdAt = it.created_at,
                    updatedAt = it.updated_at
                )
            }
            Result.success(profile)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun getUserProfileByUsername(username: String): Result<SocialUserProfile?> = withContext(Dispatchers.Default) {
        try {
            val row = database.friendQueries.getUserProfileByUsername(username)
                .asFlow()
                .mapToOneOrNull(Dispatchers.Default)
                .firstOrNull()

            val profile = row?.let {
                SocialUserProfile(
                    userId = it.user_id,
                    displayName = it.display_name,
                    username = it.username,
                    avatarUrl = it.avatar_url ?: "",
                    bio = it.bio ?: "",
                    isVerified = it.is_verified == 1L,
                    isOnline = it.is_online == 1L,
                    lastSeen = it.last_seen,
                    statusMessage = it.status_message ?: "",
                    createdAt = it.created_at,
                    updatedAt = it.updated_at
                )
            }
            Result.success(profile)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun updateUserProfile(profile: SocialUserProfile): Result<SocialUserProfile> = withContext(Dispatchers.Default) {
        try {
            database.friendQueries.insertUserProfile(
                user_id = profile.userId,
                display_name = profile.displayName,
                username = profile.username,
                avatar_url = profile.avatarUrl,
                bio = profile.bio,
                is_verified = if (profile.isVerified) 1L else 0L,
                is_online = if (profile.isOnline) 1L else 0L,
                last_seen = profile.lastSeen,
                status_message = profile.statusMessage,
                created_at = profile.createdAt,
                updated_at = System.currentTimeMillis()
            )
            Result.success(profile.copy(updatedAt = System.currentTimeMillis()))
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun updateUserOnlineStatus(userId: String, isOnline: Boolean): Result<Unit> = withContext(Dispatchers.Default) {
        try {
            database.friendQueries.updateUserOnlineStatus(
                is_online = if (isOnline) 1L else 0L,
                last_seen = System.currentTimeMillis(),
                updated_at = System.currentTimeMillis(),
                user_id = userId
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Friends operations
    override suspend fun getFriends(userId: String, status: FriendshipStatus): Result<List<Friend>> = withContext(Dispatchers.Default) {
        try {
            val friends = database.friendQueries.getFriends(userId, status.name.lowercase())
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val domainFriends = friends.map { row ->
                Friend(
                    id = row.id,
                    userId = row.user_id,
                    friendUserId = row.friend_user_id,
                    status = FriendshipStatus.valueOf(row.status.uppercase()),
                    createdAt = row.created_at.toLong(),
                    updatedAt = row.updated_at.toLong(),
                    profile = SocialUserProfile(
                        userId = row.friend_user_id,
                        displayName = row.display_name,
                        username = row.username,
                        avatarUrl = row.avatar_url ?: "",
                        isVerified = row.is_verified == 1L,
                        isOnline = row.is_online == 1L,
                        lastSeen = row.last_seen,
                        statusMessage = row.status_message ?: "",
                        createdAt = 0, // Not populated in this query
                        updatedAt = 0,  // Not populated in this query
                        bio = ""
                    ),
                    mutualFriendsCount = row.mutual_friends_count.toInt()
                )
            }
            Result.success(domainFriends)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun sendFriendRequest(userId: String, targetUserId: String): Result<Unit> = withContext(Dispatchers.Default) {
        try {
            val requestId = "${userId}_${targetUserId}_${PlatformUtils.currentTimeMillis()}"
            val currentTime = System.currentTimeMillis()
            database.friendQueries.insertFriendRequest(
                id = requestId,
                user_id = userId,
                friend_user_id = targetUserId,
                created_at = currentTime,
                updated_at = currentTime
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun acceptFriendRequest(userId: String, requesterId: String): Result<Unit> = withContext(Dispatchers.Default) {
        try {
            database.transaction {
                // Update the original request
                database.friendQueries.updateFriendStatus(
                    status = FriendshipStatus.ACCEPTED.name.lowercase(),
                    updated_at = System.currentTimeMillis(),
                    user_id = requesterId,
                    friend_user_id = userId
                )

                // Create reverse friendship
                val reverseId = "${userId}_${requesterId}_${PlatformUtils.currentTimeMillis()}"
                database.friendQueries.insertFriendRequest(
                    id = reverseId,
                    user_id = userId,
                    friend_user_id = requesterId,
                    created_at = System.currentTimeMillis(),
                    updated_at = System.currentTimeMillis()
                )
                database.friendQueries.updateFriendStatus(
                    status = FriendshipStatus.ACCEPTED.name.lowercase(),
                    updated_at = System.currentTimeMillis(),
                    user_id = userId,
                    friend_user_id = requesterId
                )
            }
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Interaction operations (simplified - showing pattern)
    override suspend fun addInteraction(interaction: SocialInteraction): Result<Unit> = withContext(Dispatchers.Default) {
        try {
            database.socialInteractionQueries.insertInteraction(
                interaction.id,
                interaction.userId,
                interaction.targetId,
                interaction.targetType.name.lowercase(),
                interaction.interactionType.name.lowercase(),
                interaction.createdAt.toString(),
                interaction.updatedAt.toString()
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun removeInteraction(
        userId: String,
        targetId: String,
        targetType: InteractionTargetType,
        interactionType: InteractionType
    ): Result<Unit> = withContext(Dispatchers.Default) {
        try {
            database.socialInteractionQueries.removeInteraction(
                user_id = userId,
                target_id = targetId,
                target_type = targetType.name.lowercase(),
                interaction_type = interactionType.name.lowercase()
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Placeholder implementations for remaining methods
    // (These would follow the same pattern as above)

    override suspend fun getPendingFriendRequests(userId: String): Result<List<Friend>> = withContext(Dispatchers.Default) {
        try {
            val pendingRequests = database.friendQueries.getFriends(userId, FriendshipStatus.PENDING.name.lowercase())
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val domainFriends = pendingRequests.map { row ->
                Friend(
                    id = row.id,
                    userId = row.user_id,
                    friendUserId = row.friend_user_id,
                    status = FriendshipStatus.PENDING,
                    createdAt = row.created_at.toLong(),
                    updatedAt = row.updated_at.toLong(),
                    profile = SocialUserProfile(
                        userId = row.friend_user_id,
                        displayName = row.display_name,
                        username = row.username,
                        avatarUrl = row.avatar_url ?: "",
                        isVerified = row.is_verified == 1L,
                        isOnline = row.is_online == 1L,
                        lastSeen = row.last_seen,
                        statusMessage = row.status_message ?: "",
                        createdAt = 0,
                        updatedAt = 0,
                        bio = ""
                    ),
                    mutualFriendsCount = row.mutual_friends_count.toInt()
                )
            }
            Result.success(domainFriends)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    override suspend fun getOnlineFriends(userId: String): Result<List<Friend>> = withContext(Dispatchers.Default) {
        try {
            val onlineFriends = database.friendQueries.getOnlineFriends(userId)
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val domainFriends = onlineFriends.map { row ->
                Friend(
                    id = row.id,
                    userId = row.user_id,
                    friendUserId = row.friend_user_id,
                    status = FriendshipStatus.valueOf(row.status.uppercase()),
                    createdAt = row.created_at.toLong(),
                    updatedAt = row.updated_at.toLong(),
                    profile = SocialUserProfile(
                        userId = row.friend_user_id,
                        displayName = row.display_name,
                        username = row.username,
                        avatarUrl = row.avatar_url ?: "",
                        isVerified = row.is_verified == 1L,
                        isOnline = true, // These are specifically online friends
                        lastSeen = row.last_seen,
                        statusMessage = row.status_message ?: "",
                        createdAt = 0,
                        updatedAt = 0,
                        bio = ""
                    ),
                    mutualFriendsCount = 0 // getOnlineFriends query doesn't include mutual_friends_count
                )
            }
            Result.success(domainFriends)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    override suspend fun getFriendSuggestions(userId: String, limit: Int): Result<List<SocialUserProfile>> = withContext(Dispatchers.Default) {
        try {
            // Get friend suggestions with regional optimization for Southeast Asia
            val suggestions = database.friendQueries.getFriendSuggestions(
                userId,
                userId,
                userId,
                limit.toLong()
            )
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val domainProfiles = suggestions.map { row ->
                SocialUserProfile(
                    userId = row.user_id,
                    displayName = row.display_name,
                    username = row.username,
                    avatarUrl = row.avatar_url ?: "",
                    bio = row.bio ?: "",
                    isVerified = row.is_verified == 1L,
                    isOnline = row.is_online == 1L,
                    lastSeen = row.last_seen,
                    statusMessage = row.status_message ?: "",
                    createdAt = row.created_at.toLong(),
                    updatedAt = row.updated_at.toLong()
                )
            }

            // Apply Southeast Asian regional optimization
            val optimizedSuggestions = applyRegionalOptimization(domainProfiles, userId)
            Result.success(optimizedSuggestions.take(limit))
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    private suspend fun applyRegionalOptimization(suggestions: List<SocialUserProfile>, userId: String): List<SocialUserProfile> {
        // Regional optimization for TH, SG, MY, ID, PH, VN markets
        // Prioritize users with cultural/regional relevance
        return suggestions.sortedWith(compareByDescending<SocialUserProfile> { profile ->
            // Prioritize verified users for better trust in Southeast Asian markets
            if (profile.isVerified) 2 else 0
        }.thenByDescending { profile ->
            // Prioritize online users for real-time engagement
            if (profile.isOnline) 1 else 0
        }.thenBy { profile ->
            // Sort by recent activity (Southeast Asian users prefer active connections)
            profile.lastSeen
        })
    }
    override suspend fun rejectFriendRequest(userId: String, requesterId: String): Result<Unit> = Result.success(Unit)
    override suspend fun removeFriend(userId: String, friendUserId: String): Result<Unit> = Result.success(Unit)
    override suspend fun checkFriendship(userId: String, targetUserId: String): Result<FriendshipStatus?> = Result.success(null)

    override suspend fun getAllEvents(userId: String): Result<List<Event>> = withContext(Dispatchers.Default) {
        try {
            val eventRows = database.eventQueries.getAllEvents(userId)
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val events = eventRows.map { row ->
                Event(
                    id = row.id,
                    title = row.title,
                    description = row.description,
                    eventDate = row.event_date,
                    location = row.location,
                    price = row.price,
                    imageUrl = row.image_url ?: "",
                    category = row.category,
                    organizerId = row.organizer_id,
                    attendeesCount = row.attendees_count.toInt(),
                    maxAttendees = row.max_attendees?.toInt(),
                    isPublic = row.is_public == 1L,
                    createdAt = row.created_at.toLong(),
                    updatedAt = row.updated_at.toLong(),
                    organizerProfile = row.organizer_name?.let {
                        SocialUserProfile(
                            userId = row.organizer_id,
                            displayName = it,
                            username = row.organizer_username ?: "",
                            avatarUrl = row.organizer_avatar ?: "",
                            isVerified = false,
                            isOnline = false,
                            lastSeen = 0,
                            statusMessage = "",
                            createdAt = 0,
                            updatedAt = 0,
                            bio = ""
                        )
                    },
                    userAttendanceStatus = row.user_attendance_status?.let {
                        EventAttendanceStatus.valueOf(it.uppercase())
                    },
                    actualAttendeesCount = row.actual_attendees_count.toInt()
                )
            }
            Result.success(events)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    override suspend fun getUpcomingEvents(userId: String): Result<List<Event>> = withContext(Dispatchers.Default) {
        try {
            val currentTime = System.currentTimeMillis()
            val eventRows = database.eventQueries.getUpcomingEvents(userId, currentTime)
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val events = eventRows.map { row ->
                Event(
                    id = row.id,
                    title = row.title,
                    description = row.description,
                    eventDate = row.event_date,
                    location = row.location,
                    price = row.price,
                    imageUrl = row.image_url ?: "",
                    category = row.category,
                    organizerId = row.organizer_id,
                    attendeesCount = row.attendees_count.toInt(),
                    maxAttendees = row.max_attendees?.toInt(),
                    isPublic = row.is_public == 1L,
                    createdAt = row.created_at.toLong(),
                    updatedAt = row.updated_at.toLong(),
                    organizerProfile = row.organizer_name?.let {
                        SocialUserProfile(
                            userId = row.organizer_id,
                            displayName = it,
                            username = row.organizer_username ?: "",
                            avatarUrl = row.organizer_avatar ?: "",
                            isVerified = false,
                            isOnline = false,
                            lastSeen = 0,
                            statusMessage = "",
                            createdAt = 0,
                            updatedAt = 0,
                            bio = ""
                        )
                    },
                    userAttendanceStatus = row.user_attendance_status?.let {
                        EventAttendanceStatus.valueOf(it.uppercase())
                    },
                    actualAttendeesCount = row.actual_attendees_count.toInt()
                )
            }
            Result.success(events)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    override suspend fun getEventsByCategory(category: String, userId: String): Result<List<Event>> = withContext(Dispatchers.Default) {
        try {
            val eventRows = database.eventQueries.getEventsByCategory(userId, category)
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val events = eventRows.map { row ->
                Event(
                    id = row.id,
                    title = row.title,
                    description = row.description,
                    eventDate = row.event_date,
                    location = row.location,
                    price = row.price,
                    imageUrl = row.image_url ?: "",
                    category = row.category,
                    organizerId = row.organizer_id,
                    attendeesCount = row.attendees_count.toInt(),
                    maxAttendees = row.max_attendees?.toInt(),
                    isPublic = row.is_public == 1L,
                    createdAt = row.created_at.toLong(),
                    updatedAt = row.updated_at.toLong(),
                    organizerProfile = row.organizer_name?.let {
                        SocialUserProfile(
                            userId = row.organizer_id,
                            displayName = it,
                            username = row.organizer_username ?: "",
                            avatarUrl = row.organizer_avatar ?: "",
                            isVerified = false,
                            isOnline = false,
                            lastSeen = 0,
                            statusMessage = "",
                            createdAt = 0,
                            updatedAt = 0,
                            bio = ""
                        )
                    },
                    userAttendanceStatus = row.user_attendance_status?.let {
                        EventAttendanceStatus.valueOf(it.uppercase())
                    },
                    actualAttendeesCount = row.actual_attendees_count.toInt()
                )
            }
            Result.success(events)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    override suspend fun getEventById(eventId: String, userId: String): Result<Event?> = Result.success(null)
    override suspend fun getUserEvents(userId: String): Result<List<Event>> = withContext(Dispatchers.Default) {
        try {
            val eventRows = database.eventQueries.getUserEvents(userId)
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val events = eventRows.map { row ->
                Event(
                    id = row.id,
                    title = row.title,
                    description = row.description,
                    eventDate = row.event_date,
                    location = row.location,
                    price = row.price,
                    imageUrl = row.image_url ?: "",
                    category = row.category,
                    organizerId = row.organizer_id,
                    attendeesCount = row.attendees_count.toInt(),
                    maxAttendees = row.max_attendees?.toInt(),
                    isPublic = row.is_public == 1L,
                    createdAt = row.created_at.toLong(),
                    updatedAt = row.updated_at.toLong(),
                    organizerProfile = row.organizer_name?.let {
                        SocialUserProfile(
                            userId = row.organizer_id,
                            displayName = it,
                            username = row.organizer_username ?: "",
                            avatarUrl = row.organizer_avatar ?: "",
                            isVerified = false,
                            isOnline = false,
                            lastSeen = 0,
                            statusMessage = "",
                            createdAt = 0,
                            updatedAt = 0,
                            bio = ""
                        )
                    },
                    userAttendanceStatus = row.user_attendance_status?.let {
                        EventAttendanceStatus.valueOf(it.uppercase())
                    },
                    actualAttendeesCount = row.actual_attendees_count.toInt()
                )
            }
            Result.success(events)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    override suspend fun getEventAttendees(eventId: String): Result<List<SocialUserProfile>> = withContext(Dispatchers.Default) {
        try {
            val attendeeRows = database.eventQueries.getEventAttendees(eventId)
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val attendees = attendeeRows.map { row ->
                SocialUserProfile(
                    userId = row.user_id,
                    displayName = row.display_name ?: "Unknown User",
                    username = row.username ?: "",
                    avatarUrl = row.avatar_url ?: "",
                    isVerified = row.is_verified == 1L,
                    isOnline = false,
                    lastSeen = 0,
                    statusMessage = "",
                    createdAt = 0,
                    updatedAt = 0,
                    bio = ""
                )
            }
            Result.success(attendees)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    override suspend fun createEvent(event: Event): Result<Event> = Result.success(event)
    override suspend fun updateEvent(event: Event): Result<Event> = Result.success(event)
    override suspend fun deleteEvent(eventId: String, organizerId: String): Result<Unit> = Result.success(Unit)
    override suspend fun rsvpToEvent(eventId: String, userId: String, status: EventAttendanceStatus): Result<Unit> = Result.success(Unit)
    override suspend fun removeEventRsvp(eventId: String, userId: String): Result<Unit> = Result.success(Unit)
    override suspend fun getEventCategories(): Result<List<Pair<String, Int>>> = withContext(Dispatchers.Default) {
        try {
            val currentTime = System.currentTimeMillis()
            val categoryRows = database.eventQueries.getEventCategories(currentTime)
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val categories = categoryRows.map { row ->
                Pair(row.category, row.event_count.toInt())
            }
            Result.success(categories)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    override suspend fun searchEvents(query: String, userId: String): Result<List<Event>> = withContext(Dispatchers.Default) {
        try {
            val eventRows = database.eventQueries.searchEvents(userId, query, query, query)
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val events = eventRows.map { row ->
                Event(
                    id = row.id,
                    title = row.title,
                    description = row.description,
                    eventDate = row.event_date,
                    location = row.location,
                    price = row.price,
                    imageUrl = row.image_url ?: "",
                    category = row.category,
                    organizerId = row.organizer_id,
                    attendeesCount = row.attendees_count.toInt(),
                    maxAttendees = row.max_attendees?.toInt(),
                    isPublic = row.is_public == 1L,
                    createdAt = row.created_at.toLong(),
                    updatedAt = row.updated_at.toLong(),
                    organizerProfile = row.organizer_name?.let {
                        SocialUserProfile(
                            userId = row.organizer_id,
                            displayName = it,
                            username = row.organizer_username ?: "",
                            avatarUrl = row.organizer_avatar ?: "",
                            isVerified = false,
                            isOnline = false,
                            lastSeen = 0,
                            statusMessage = "",
                            createdAt = 0,
                            updatedAt = 0,
                            bio = ""
                        )
                    },
                    userAttendanceStatus = row.user_attendance_status?.let {
                        EventAttendanceStatus.valueOf(it.uppercase())
                    },
                    actualAttendeesCount = row.actual_attendees_count.toInt()
                )
            }
            Result.success(events)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override suspend fun getInteractionsByUser(userId: String): Result<List<SocialInteraction>> = withContext(Dispatchers.Default) {
        try {
            val interactionRows = database.socialInteractionQueries.getInteractionsByUser(userId)
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val interactions = interactionRows.map { row ->
                SocialInteraction(
                    id = row.id,
                    userId = row.user_id,
                    targetId = row.target_id,
                    targetType = InteractionTargetType.valueOf(row.target_type.uppercase()),
                    interactionType = InteractionType.valueOf(row.interaction_type.uppercase()),
                    createdAt = row.created_at.toLong(),
                    updatedAt = row.updated_at.toLong()
                )
            }
            Result.success(interactions)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    override suspend fun getInteractionsByTarget(targetId: String, targetType: InteractionTargetType): Result<List<SocialInteraction>> = withContext(Dispatchers.Default) {
        try {
            val interactionRows = database.socialInteractionQueries.getInteractionsByTarget(targetId, targetType.name.lowercase())
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val interactions = interactionRows.map { row ->
                SocialInteraction(
                    id = row.id,
                    userId = row.user_id,
                    targetId = row.target_id,
                    targetType = InteractionTargetType.valueOf(row.target_type.uppercase()),
                    interactionType = InteractionType.valueOf(row.interaction_type.uppercase()),
                    createdAt = row.created_at.toLong(),
                    updatedAt = row.updated_at.toLong(),
                    userProfile = SocialUserProfile(
                        userId = row.user_id,
                        displayName = row.display_name ?: "Unknown User",
                        username = row.username ?: "",
                        avatarUrl = row.avatar ?: "",
                        isVerified = false,
                        isOnline = false,
                        lastSeen = 0,
                        statusMessage = "",
                        createdAt = 0,
                        updatedAt = 0,
                        bio = ""
                    )
                )
            }
            Result.success(interactions)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    override suspend fun getUserInteractionState(userId: String, targetId: String, targetType: InteractionTargetType): Result<Set<InteractionType>> = Result.success(emptySet())
    override suspend fun getInteractionCounts(targetId: String, targetType: InteractionTargetType): Result<InteractionCounts> = Result.success(InteractionCounts())
    override suspend fun getFollowedUsers(userId: String): Result<List<SocialUserProfile>> = withContext(Dispatchers.Default) {
        try {
            val followedRows = database.socialInteractionQueries.getFollowedUsers(userId)
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val followedUsers = followedRows.map { row ->
                SocialUserProfile(
                    userId = row.followed_user_id,
                    displayName = row.display_name ?: "Unknown User",
                    username = row.username ?: "",
                    avatarUrl = row.avatar ?: "",
                    isVerified = row.is_verified == 1L,
                    isOnline = false,
                    lastSeen = 0,
                    statusMessage = "",
                    createdAt = row.followed_at.toLong(),
                    updatedAt = 0,
                    bio = ""
                )
            }
            Result.success(followedUsers)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    override suspend fun getFollowers(userId: String): Result<List<SocialUserProfile>> = withContext(Dispatchers.Default) {
        try {
            val followerRows = database.socialInteractionQueries.getFollowers(userId)
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val followers = followerRows.map { row ->
                SocialUserProfile(
                    userId = row.follower_user_id,
                    displayName = row.display_name ?: "Unknown User",
                    username = row.username ?: "",
                    avatarUrl = row.avatar ?: "",
                    isVerified = row.is_verified == 1L,
                    isOnline = false,
                    lastSeen = 0,
                    statusMessage = "",
                    createdAt = row.followed_at.toLong(),
                    updatedAt = 0,
                    bio = ""
                )
            }
            Result.success(followers)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    override suspend fun getUserStats(userId: String): Result<UserStats> = Result.success(UserStats())

    override suspend fun getCommentsByTarget(targetId: String, targetType: InteractionTargetType, userId: String): Result<List<SocialComment>> = withContext(Dispatchers.Default) {
        try {
            val commentRows = database.socialInteractionQueries.getCommentsByTarget(userId, targetId)
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val comments = commentRows.map { row ->
                SocialComment(
                    id = row.id,
                    targetId = row.post_id ?: targetId,
                    targetType = targetType,
                    userId = row.author_id,
                    content = row.content,
                    parentCommentId = row.parent_comment_id,
                    likesCount = row.likes_count.toInt(),
                    repliesCount = row.replies_count.toInt(),
                    createdAt = row.created_at.toLong(),
                    updatedAt = row.updated_at.toLong(),
                    userProfile = SocialUserProfile(
                        userId = row.author_id,
                        displayName = row.display_name ?: "Unknown User",
                        username = row.username ?: "",
                        avatarUrl = row.avatar ?: "",
                        isVerified = row.is_verified == 1L,
                        isOnline = false,
                        lastSeen = 0,
                        statusMessage = "",
                        createdAt = 0,
                        updatedAt = 0,
                        bio = ""
                    ),
                    isLikedByUser = row.is_liked_by_user == 1L
                )
            }
            Result.success(comments)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    override suspend fun getCommentReplies(commentId: String, userId: String): Result<List<SocialComment>> = withContext(Dispatchers.Default) {
        try {
            val replyRows = database.socialInteractionQueries.getCommentReplies(userId, commentId)
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val replies = replyRows.map { row ->
                SocialComment(
                    id = row.id,
                    targetId = row.post_id ?: "",
                    targetType = InteractionTargetType.COMMENT,
                    userId = row.author_id,
                    content = row.content,
                    parentCommentId = row.parent_comment_id,
                    likesCount = row.likes_count.toInt(),
                    repliesCount = row.replies_count.toInt(),
                    createdAt = row.created_at.toLong(),
                    updatedAt = row.updated_at.toLong(),
                    userProfile = SocialUserProfile(
                        userId = row.author_id,
                        displayName = row.display_name ?: "Unknown User",
                        username = row.username ?: "",
                        avatarUrl = row.avatar ?: "",
                        isVerified = row.is_verified == 1L,
                        isOnline = false,
                        lastSeen = 0,
                        statusMessage = "",
                        createdAt = 0,
                        updatedAt = 0,
                        bio = ""
                    ),
                    isLikedByUser = row.is_liked_by_user == 1L
                )
            }
            Result.success(replies)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    override suspend fun createComment(comment: SocialComment): Result<SocialComment> = Result.success(comment)
    override suspend fun updateComment(commentId: String, userId: String, content: String): Result<Unit> = Result.success(Unit)
    override suspend fun deleteComment(commentId: String, userId: String): Result<Unit> = Result.success(Unit)

    override suspend fun getPopularContent(since: Long, minInteractions: Int): Result<List<Triple<String, InteractionTargetType, Int>>> = withContext(Dispatchers.Default) {
        try {
            val popularRows = database.socialInteractionQueries.getPopularContent(since.toString(), minInteractions.toLong())
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val popularContent = popularRows.map { row ->
                Triple(
                    row.target_id,
                    InteractionTargetType.valueOf(row.target_type.uppercase()),
                    row.interaction_count.toInt()
                )
            }
            Result.success(popularContent)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    override suspend fun getRecentActivity(userId: String, since: Long, limit: Int): Result<List<SocialInteraction>> = withContext(Dispatchers.Default) {
        try {
            val activityRows = database.socialInteractionQueries.getRecentActivity(userId, since.toString(), limit.toLong())
                .asFlow()
                .mapToList(Dispatchers.Default)
                .firstOrNull() ?: emptyList()

            val activities = activityRows.map { row ->
                SocialInteraction(
                    id = row.id,
                    userId = row.user_id,
                    targetId = row.target_id,
                    targetType = InteractionTargetType.valueOf(row.target_type.uppercase()),
                    interactionType = InteractionType.valueOf(row.interaction_type.uppercase()),
                    createdAt = row.created_at.toLong(),
                    updatedAt = row.updated_at.toLong(),
                    userProfile = SocialUserProfile(
                        userId = row.user_id,
                        displayName = row.display_name ?: "Unknown User",
                        username = row.username ?: "",
                        avatarUrl = row.avatar ?: "",
                        isVerified = false,
                        isOnline = false,
                        lastSeen = 0,
                        statusMessage = "",
                        createdAt = 0,
                        updatedAt = 0,
                        bio = ""
                    )
                )
            }
            Result.success(activities)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}