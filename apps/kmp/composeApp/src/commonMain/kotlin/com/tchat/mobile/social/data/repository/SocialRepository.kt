package com.tchat.mobile.social.data.repository

import com.tchat.mobile.database.TchatDatabase
import com.tchat.mobile.social.data.api.SocialApiClient
import com.tchat.mobile.social.domain.models.*
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.*
import kotlinx.coroutines.launch
import kotlinx.datetime.Clock
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json

/**
 * KMP Social Repository
 *
 * Offline-first repository implementing social features with:
 * - Cross-platform offline caching
 * - Incremental sync patterns
 * - Conflict resolution
 * - Southeast Asian regional optimization
 * - Mobile performance optimization
 */
class SocialRepository(
    private val apiClient: SocialApiClient,
    private val database: TchatDatabase,
    private val scope: CoroutineScope = CoroutineScope(Dispatchers.Default)
) {

    private val json = Json { ignoreUnknownKeys = true }
    private val _syncState = MutableStateFlow(SyncState.IDLE)
    val syncState: StateFlow<SyncState> = _syncState.asStateFlow()

    // Profile Management
    fun getProfileFlow(userId: String): Flow<SocialProfile?> = flow {
        // Emit cached profile first
        val cachedProfile = database.socialProfileQueries.getSocialProfile(userId).executeAsOneOrNull()
        cachedProfile?.let { emit(it.toDomainModel()) }

        // Fetch from server if online
        if (apiClient.isConnected()) {
            try {
                val result = apiClient.getSocialProfile(userId)
                result.getOrNull()?.let { profile ->
                    // Update cache
                    database.socialProfileQueries.insertSocialProfile(
                        id = profile.id,
                        username = profile.username,
                        display_name = profile.displayName,
                        bio = profile.bio,
                        avatar = profile.avatar,
                        interests = json.encodeToString(profile.interests),
                        social_links = profile.socialLinks?.let { json.encodeToString(it) },
                        social_preferences = profile.socialPreferences?.let { json.encodeToString(it) },
                        followers_count = profile.followersCount.toLong(),
                        following_count = profile.followingCount.toLong(),
                        posts_count = profile.postsCount.toLong(),
                        is_social_verified = if (profile.isSocialVerified) 1L else 0L,
                        country = profile.country,
                        region = profile.region,
                        social_created_at = profile.socialCreatedAt,
                        social_updated_at = profile.socialUpdatedAt,
                        last_sync_at = Clock.System.now().toString(),
                        sync_version = profile.syncVersion,
                        is_offline_edit = 0L,
                        sync_status = "synced"
                    )
                    emit(profile)
                }
            } catch (e: Exception) {
                // Silently continue with cached data
            }
        }
    }

    suspend fun updateProfile(userId: String, request: UpdateProfileRequest): Result<SocialProfile> {
        return try {
            val timestamp = Clock.System.now().toString()

            if (apiClient.isConnected()) {
                // Online update
                val result = apiClient.updateSocialProfile(userId, request)
                result.fold(
                    onSuccess = { profile ->
                        // Update cache with server response
                        database.socialProfileQueries.updateSocialProfile(
                            display_name = profile.displayName,
                            bio = profile.bio,
                            interests = json.encodeToString(profile.interests),
                            social_links = profile.socialLinks?.let { json.encodeToString(it) },
                            social_preferences = profile.socialPreferences?.let { json.encodeToString(it) },
                            social_updated_at = profile.socialUpdatedAt,
                            last_sync_at = timestamp,
                            is_offline_edit = 0L,
                            sync_status = "synced",
                            id = userId
                        )
                        Result.success(profile)
                    },
                    onFailure = { error ->
                        // Store offline operation
                        storeOfflineOperation("update", "profile", userId, request)
                        Result.failure(error)
                    }
                )
            } else {
                // Offline update
                database.socialProfileQueries.updateSocialProfile(
                    display_name = request.displayName,
                    bio = request.bio,
                    interests = request.interests?.let { json.encodeToString(it) } ?: "[]",
                    social_links = request.socialLinks?.let { json.encodeToString(it) },
                    social_preferences = null,
                    social_updated_at = timestamp,
                    last_sync_at = timestamp,
                    is_offline_edit = 1L,
                    sync_status = "pending",
                    id = userId
                )

                storeOfflineOperation("update", "profile", userId, request)

                // Return updated profile from cache
                val updatedProfile = database.socialProfileQueries.getSocialProfile(userId).executeAsOneOrNull()
                updatedProfile?.let {
                    Result.success(it.toDomainModel())
                } ?: Result.failure(Exception("Profile not found"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Feed Management
    fun getFeedFlow(
        userId: String,
        feedType: String = "home",
        region: String = "TH"
    ): Flow<SocialFeed> = flow {
        // Check cache first
        val cacheKey = "${userId}_${feedType}_${region}"
        val cachedFeed = database.socialProfileQueries.getFeedCache(
            user_id = userId,
            feed_type = feedType,
            region = region,
            expires_at = Clock.System.now().toString()
        ).executeAsOneOrNull()

        if (cachedFeed != null) {
            val postIds = json.decodeFromString<List<String>>(cachedFeed.post_ids)
            val posts = postIds.mapNotNull { postId ->
                database.socialProfileQueries.getSocialPost(postId).executeAsOneOrNull()?.toDomainModel()
            }
            emit(SocialFeed(
                posts = posts,
                hasMore = posts.size >= 20,
                nextOffset = posts.size,
                lastSyncAt = cachedFeed.last_updated_at,
                region = region,
                feedType = feedType
            ))
        }

        // Fetch from server if online
        if (apiClient.isConnected()) {
            try {
                val result = apiClient.getUserFeed(userId, limit = 20, offset = 0, feedType = feedType)
                result.getOrNull()?.let { feed ->
                    // Cache posts
                    feed.posts.forEach { post ->
                        cachePost(post)
                    }

                    // Cache feed
                    val postIds = feed.posts.map { it.id }
                    val expires = Clock.System.now().plus(kotlin.time.Duration.parse("1h")).toString()
                    database.socialProfileQueries.insertFeedCache(
                        id = cacheKey,
                        user_id = userId,
                        feed_type = feedType,
                        region = region,
                        post_ids = json.encodeToString(postIds),
                        last_updated_at = feed.lastSyncAt,
                        expires_at = expires
                    )

                    emit(feed)
                }
            } catch (e: Exception) {
                // Continue with cached data
            }
        }
    }

    suspend fun refreshFeed(
        userId: String,
        feedType: String = "home",
        region: String = "TH"
    ): Result<SocialFeed> {
        return if (apiClient.isConnected()) {
            apiClient.getUserFeed(userId, limit = 20, offset = 0, feedType = feedType)
        } else {
            Result.failure(Exception("No internet connection"))
        }
    }

    // Post Management
    suspend fun createPost(request: CreatePostRequest): Result<SocialPost> {
        return try {
            val timestamp = Clock.System.now().toString()
            val postId = "temp_${timestamp.hashCode()}"

            if (apiClient.isConnected()) {
                // Online creation
                val result = apiClient.createPost(request)
                result.fold(
                    onSuccess = { post ->
                        cachePost(post)
                        Result.success(post)
                    },
                    onFailure = { error ->
                        // Store offline operation
                        storeOfflineOperation("create", "post", postId, request)
                        Result.failure(error)
                    }
                )
            } else {
                // Offline creation - create temporary post
                val tempPost = SocialPost(
                    id = postId,
                    authorId = request.content, // This should come from user session
                    authorUsername = "temp_user", // This should come from user session
                    authorDisplayName = null,
                    authorAvatar = null,
                    content = request.content,
                    contentType = request.contentType,
                    mediaUrls = request.mediaUrls,
                    tags = request.tags,
                    mentions = request.mentions,
                    visibility = request.visibility,
                    createdAt = timestamp,
                    updatedAt = timestamp,
                    language = request.language,
                    region = request.region,
                    isOfflineEdit = true
                )

                cachePost(tempPost, isOffline = true)
                storeOfflineOperation("create", "post", postId, request)

                Result.success(tempPost)
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    fun getPostFlow(postId: String): Flow<SocialPost?> = flow {
        // Emit cached post first
        val cachedPost = database.socialProfileQueries.getSocialPost(postId).executeAsOneOrNull()
        cachedPost?.let { emit(it.toDomainModel()) }

        // Fetch from server if online
        if (apiClient.isConnected()) {
            try {
                val result = apiClient.getPost(postId)
                result.getOrNull()?.let { post ->
                    cachePost(post)
                    emit(post)
                }
            } catch (e: Exception) {
                // Continue with cached data
            }
        }
    }

    suspend fun likePost(postId: String, userId: String): Result<Unit> {
        return toggleInteraction(
            targetId = postId,
            targetType = "post",
            interactionType = "like",
            userId = userId
        )
    }

    suspend fun bookmarkPost(postId: String, userId: String): Result<Unit> {
        return toggleInteraction(
            targetId = postId,
            targetType = "post",
            interactionType = "bookmark",
            userId = userId
        )
    }

    // Follow Management
    suspend fun followUser(userId: String, followingId: String): Result<Unit> {
        return try {
            if (apiClient.isConnected()) {
                val result = apiClient.followUser(followingId)
                result.fold(
                    onSuccess = {
                        // Update local interaction
                        database.socialInteractionQueries.insertInteraction(
                            id = "${userId}_${followingId}_follow",
                            user_id = userId,
                            target_id = followingId,
                            target_type = "user",
                            interaction_type = "follow",
                            created_at = Clock.System.now().toString(),
                            updated_at = Clock.System.now().toString(),
                        )
                        Result.success(Unit)
                    },
                    onFailure = { error ->
                        // Store offline operation
                        storeOfflineOperation("create", "interaction",
                            "${userId}_${followingId}_follow",
                            mapOf("targetId" to followingId, "targetType" to "user", "interactionType" to "follow")
                        )
                        Result.failure(error)
                    }
                )
            } else {
                // Offline follow
                database.socialInteractionQueries.insertInteraction(
                    id = "${userId}_${followingId}_follow",
                    user_id = userId,
                    target_id = followingId,
                    target_type = "user",
                    interaction_type = "follow",
                    created_at = Clock.System.now().toString(),
                    updated_at = Clock.System.now().toString(),
                )

                storeOfflineOperation("create", "interaction",
                    "${userId}_${followingId}_follow",
                    mapOf("targetId" to followingId, "targetType" to "user", "interactionType" to "follow")
                )

                Result.success(Unit)
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun unfollowUser(userId: String, followingId: String): Result<Unit> {
        return try {
            if (apiClient.isConnected()) {
                val result = apiClient.unfollowUser(followingId)
                result.fold(
                    onSuccess = {
                        // Remove local interaction
                        database.socialInteractionQueries.removeInteraction(
                            user_id = userId,
                            target_id = followingId,
                            target_type = "user",
                            interaction_type = "follow"
                        )
                        Result.success(Unit)
                    },
                    onFailure = { error ->
                        // Store offline operation
                        storeOfflineOperation("delete", "interaction",
                            "${userId}_${followingId}_follow",
                            mapOf("targetId" to followingId, "targetType" to "user", "interactionType" to "follow")
                        )
                        Result.failure(error)
                    }
                )
            } else {
                // Offline unfollow
                database.socialInteractionQueries.removeInteraction(
                    user_id = userId,
                    target_id = followingId,
                    target_type = "user",
                    interaction_type = "follow"
                )

                storeOfflineOperation("delete", "interaction",
                    "${userId}_${followingId}_follow",
                    mapOf("targetId" to followingId, "targetType" to "user", "interactionType" to "follow")
                )

                Result.success(Unit)
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    fun getFollowingFlow(userId: String): Flow<List<SocialProfile>> = flow {
        val following = database.socialInteractionQueries.getFollowedUsers(
            user_id = userId
        ).executeAsList().map {
            SocialProfile(
                id = it.followed_user_id,
                username = it.username ?: "",
                displayName = it.display_name,
                bio = null,
                avatar = it.avatar,
                interests = emptyList(),
                socialLinks = null,
                socialPreferences = null,
                followersCount = 0,
                followingCount = 0,
                postsCount = 0,
                isSocialVerified = it.is_verified == 1L,
                country = "",
                region = "",
                socialCreatedAt = "",
                socialUpdatedAt = "",
                lastSyncAt = it.followed_at,
                syncVersion = "0"
            )
        }

        emit(following)
    }

    fun getFollowersFlow(userId: String): Flow<List<SocialProfile>> = flow {
        val followers = database.socialInteractionQueries.getFollowers(
            target_id = userId
        ).executeAsList().map {
            SocialProfile(
                id = it.follower_user_id,
                username = it.username ?: "",
                displayName = it.display_name,
                bio = null,
                avatar = it.avatar,
                interests = emptyList(),
                socialLinks = null,
                socialPreferences = null,
                followersCount = 0,
                followingCount = 0,
                postsCount = 0,
                isSocialVerified = it.is_verified == 1L,
                country = "",
                region = "",
                socialCreatedAt = "",
                socialUpdatedAt = "",
                lastSyncAt = it.followed_at,
                syncVersion = "0"
            )
        }

        emit(followers)
    }

    // Discovery
    suspend fun getDiscoveryFeed(userId: String, region: String = "TH"): Result<List<DiscoveryProfile>> {
        return if (apiClient.isConnected()) {
            apiClient.getDiscoveryFeed(userId, region)
        } else {
            // Return cached regional profiles
            val regionalProfiles = database.socialProfileQueries.getProfilesNeedingSync()
                .executeAsList()
                .filter { it.region == region && it.id != userId }
                .map {
                    DiscoveryProfile(
                        profile = it.toDomainModel(),
                        discoveryReason = "region",
                        score = 0.5
                    )
                }
                .take(10)

            Result.success(regionalProfiles)
        }
    }

    // Sync Management
    suspend fun performIncrementalSync(userId: String, lastSyncAt: String?): Result<SyncResponse> {
        if (!apiClient.isConnected()) {
            return Result.failure(Exception("No internet connection"))
        }

        _syncState.value = SyncState.SYNCING

        return try {
            // Get pending operations
            val pendingOps = database.socialProfileQueries.getPendingSyncOperations().executeAsList()

            if (pendingOps.isNotEmpty()) {
                // Apply client changes
                val operations = pendingOps.map { it.toDomainModel() }
                val result = apiClient.applyClientChanges(userId, operations)

                result.fold(
                    onSuccess = { syncResponse ->
                        // Update operation statuses
                        syncResponse.syncedOperations.forEach { opId ->
                            database.socialProfileQueries.updateSyncOperationStatus(
                                id = opId,
                                status = "synced",
                                error_message = null
                            )
                        }

                        syncResponse.failedOperations.forEach { opId ->
                            database.socialProfileQueries.updateSyncOperationStatus(
                                id = opId,
                                status = "failed",
                                error_message = "Sync failed"
                            )
                        }

                        _syncState.value = SyncState.SUCCESS
                        Result.success(syncResponse)
                    },
                    onFailure = { error ->
                        _syncState.value = SyncState.ERROR
                        Result.failure(error)
                    }
                )
            } else {
                // No pending operations, just fetch updates
                val profileChanges = apiClient.getProfileChanges(userId, lastSyncAt ?: "")
                _syncState.value = SyncState.SUCCESS
                profileChanges
            }
        } catch (e: Exception) {
            _syncState.value = SyncState.ERROR
            Result.failure(e)
        }
    }

    // Private helper methods
    private suspend fun toggleInteraction(
        targetId: String,
        targetType: String,
        interactionType: String,
        userId: String
    ): Result<Unit> {
        return try {
            val existingInteraction = database.socialInteractionQueries.getUserInteractionState(
                user_id = userId,
                target_id = targetId,
                target_type = targetType
            ).executeAsOneOrNull()

            if (existingInteraction != null) {
                // Remove interaction
                if (apiClient.isConnected()) {
                    val result = apiClient.removeInteraction(targetId, targetType, interactionType)
                    result.fold(
                        onSuccess = {
                            database.socialInteractionQueries.removeInteraction(
                                user_id = userId,
                                target_id = targetId,
                                target_type = targetType,
                                interaction_type = interactionType
                            )
                            Result.success(Unit)
                        },
                        onFailure = { error ->
                            storeOfflineOperation("delete", "interaction",
                                "${userId}_${targetId}_${interactionType}",
                                mapOf("targetId" to targetId, "targetType" to targetType, "interactionType" to interactionType)
                            )
                            Result.failure(error)
                        }
                    )
                } else {
                    // Offline removal
                    database.socialInteractionQueries.removeInteraction(
                        user_id = userId,
                        target_id = targetId,
                        target_type = targetType,
                        interaction_type = interactionType
                    )
                    storeOfflineOperation("delete", "interaction",
                        "${userId}_${targetId}_${interactionType}",
                        mapOf("targetId" to targetId, "targetType" to targetType, "interactionType" to interactionType)
                    )
                    Result.success(Unit)
                }
            } else {
                // Add interaction
                val interactionRequest = InteractionRequest(
                    targetId = targetId,
                    targetType = targetType,
                    interactionType = interactionType
                )

                if (apiClient.isConnected()) {
                    val result = apiClient.createInteraction(interactionRequest)
                    result.fold(
                        onSuccess = { interaction ->
                            database.socialInteractionQueries.insertInteraction(
                                id = interaction.id,
                                user_id = interaction.userId,
                                target_id = interaction.targetId,
                                target_type = interaction.targetType,
                                interaction_type = interaction.interactionType,
                                created_at = interaction.createdAt,
                                updated_at = interaction.updatedAt,
                            )
                            Result.success(Unit)
                        },
                        onFailure = { error ->
                            storeOfflineOperation("create", "interaction",
                                "${userId}_${targetId}_${interactionType}",
                                interactionRequest
                            )
                            Result.failure(error)
                        }
                    )
                } else {
                    // Offline addition
                    database.socialInteractionQueries.insertInteraction(
                        id = "${userId}_${targetId}_${interactionType}",
                        user_id = userId,
                        target_id = targetId,
                        target_type = targetType,
                        interaction_type = interactionType,
                        created_at = Clock.System.now().toString(),
                        updated_at = Clock.System.now().toString(),
                    )
                    storeOfflineOperation("create", "interaction",
                        "${userId}_${targetId}_${interactionType}",
                        interactionRequest
                    )
                    Result.success(Unit)
                }
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    private fun cachePost(post: SocialPost, isOffline: Boolean = false) {
        database.socialProfileQueries.insertSocialPost(
            id = post.id,
            author_id = post.authorId,
            author_username = post.authorUsername,
            author_display_name = post.authorDisplayName,
            author_avatar = post.authorAvatar,
            content = post.content,
            content_type = post.contentType,
            media_urls = json.encodeToString(post.mediaUrls),
            thumbnail_url = post.thumbnailUrl,
            tags = json.encodeToString(post.tags),
            mentions = json.encodeToString(post.mentions),
            visibility = post.visibility,
            likes_count = post.likesCount.toLong(),
            comments_count = post.commentsCount.toLong(),
            shares_count = post.sharesCount.toLong(),
            views_count = post.viewsCount.toLong(),
            is_liked_by_user = if (post.isLikedByUser) 1L else 0L,
            is_bookmarked_by_user = if (post.isBookmarkedByUser) 1L else 0L,
            created_at = post.createdAt,
            updated_at = post.updatedAt,
            language = post.language,
            region = post.region,
            is_regional_trending = if (post.isRegionalTrending) 1L else 0L,
            last_sync_at = post.lastSyncAt ?: Clock.System.now().toString(),
            is_offline_edit = if (isOffline) 1L else 0L,
            sync_status = if (isOffline) "pending" else "synced"
        )
    }

    private fun storeOfflineOperation(
        operation: String,
        resourceType: String,
        resourceId: String,
        data: Any
    ) {
        val timestamp = Clock.System.now().toString()
        database.socialProfileQueries.insertSyncOperation(
            id = "${operation}_${resourceType}_${resourceId}_${timestamp.hashCode()}",
            operation = operation,
            resource_type = resourceType,
            resource_id = resourceId,
            data_ = json.encodeToString(data),
            timestamp = timestamp,
            status = "pending",
            conflict_resolution = null,
            retry_count = 0L,
            error_message = null,
            created_at = timestamp
        )
    }

    private fun cleanupExpiredCache() {
        scope.launch {
            val now = Clock.System.now().toString()
            database.socialProfileQueries.deleteExpiredFeedCache(now)
        }
    }

    init {
        // Cleanup expired cache periodically
        cleanupExpiredCache()
    }
}

enum class SyncState {
    IDLE,
    SYNCING,
    SUCCESS,
    ERROR
}

// Extension functions to convert between database and domain models
private fun com.tchat.mobile.database.Social_profiles.toDomainModel(): SocialProfile {
    val json = Json { ignoreUnknownKeys = true }
    return SocialProfile(
        id = id,
        username = username,
        displayName = display_name,
        bio = bio,
        avatar = avatar,
        interests = interests?.let { json.decodeFromString<List<String>>(it) } ?: emptyList(),
        socialLinks = social_links?.let { json.decodeFromString(it) },
        socialPreferences = social_preferences?.let { json.decodeFromString(it) },
        followersCount = followers_count.toInt(),
        followingCount = following_count.toInt(),
        postsCount = posts_count.toInt(),
        isSocialVerified = is_social_verified == 1L,
        country = country,
        region = region,
        socialCreatedAt = social_created_at,
        socialUpdatedAt = social_updated_at,
        lastSyncAt = last_sync_at,
        syncVersion = sync_version
    )
}

private fun com.tchat.mobile.database.Social_posts.toDomainModel(): SocialPost {
    val json = Json { ignoreUnknownKeys = true }
    return SocialPost(
        id = id,
        authorId = author_id,
        authorUsername = author_username,
        authorDisplayName = author_display_name,
        authorAvatar = author_avatar,
        content = content,
        contentType = content_type,
        mediaUrls = json.decodeFromString<List<String>>(media_urls),
        thumbnailUrl = thumbnail_url,
        tags = json.decodeFromString<List<String>>(tags),
        mentions = json.decodeFromString<List<String>>(mentions),
        visibility = visibility,
        likesCount = likes_count.toInt(),
        commentsCount = comments_count.toInt(),
        sharesCount = shares_count.toInt(),
        viewsCount = views_count.toInt(),
        isLikedByUser = is_liked_by_user == 1L,
        isBookmarkedByUser = is_bookmarked_by_user == 1L,
        createdAt = created_at,
        updatedAt = updated_at,
        language = language,
        region = region,
        isRegionalTrending = is_regional_trending == 1L,
        lastSyncAt = last_sync_at,
        isOfflineEdit = is_offline_edit == 1L
    )
}

private fun com.tchat.mobile.database.Social_sync_operations.toDomainModel(): SyncOperation {
    return SyncOperation(
        id = id,
        operation = operation,
        resourceType = resource_type,
        resourceId = resource_id,
        data = data_?.let { Json.parseToJsonElement(it) },
        timestamp = timestamp,
        status = status,
        conflictResolution = conflict_resolution,
        retryCount = retry_count.toInt()
    )
}