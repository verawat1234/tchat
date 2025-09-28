package com.tchat.mobile.social.data.api

import com.tchat.mobile.network.HttpClientFactory
import com.tchat.mobile.social.domain.models.*
import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.request.*
import io.ktor.client.statement.*
import io.ktor.http.*
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.datetime.Clock

/**
 * KMP Social API Client
 *
 * Mobile-optimized API client for social features with:
 * - Southeast Asian regional support
 * - Incremental sync patterns
 * - Offline-first architecture
 * - Conflict resolution
 * - Mobile performance optimization
 */
class SocialApiClient(
    private val baseUrl: String = "http://localhost:8080/api/v1/social"
) {

    private val httpClient: HttpClient = HttpClientFactory.create(baseUrl)
    private val _connectionState = MutableStateFlow(SocialConnectionState.DISCONNECTED)
    val connectionState: StateFlow<SocialConnectionState> = _connectionState

    private var accessToken: String? = null

    // Connection Management
    suspend fun connect(): Result<Unit> {
        return try {
            val response: HttpResponse = httpClient.get("/mobile/health")
            if (response.status.isSuccess()) {
                _connectionState.value = SocialConnectionState.CONNECTED
                Result.success(Unit)
            } else {
                _connectionState.value = SocialConnectionState.ERROR
                Result.failure(SocialApiException(response.status.value, "Failed to connect to social service"))
            }
        } catch (e: Exception) {
            _connectionState.value = SocialConnectionState.ERROR
            Result.failure(SocialApiException(500, "Connection failed: ${e.message}"))
        }
    }

    fun setAuthToken(token: String) {
        accessToken = token
    }

    // Initial Data Loading
    suspend fun getInitialUserData(userId: String): Result<InitialDataResponse> {
        return executeAuthenticatedRequest {
            httpClient.get("/mobile/init/$userId")
        }
    }

    // Profile Management
    suspend fun getSocialProfile(userId: String): Result<SocialProfile> {
        return executeAuthenticatedRequest {
            httpClient.get("/mobile/profile/$userId")
        }
    }

    suspend fun updateSocialProfile(userId: String, request: UpdateProfileRequest): Result<SocialProfile> {
        return executeAuthenticatedRequest {
            httpClient.put("/mobile/profile/$userId") {
                contentType(ContentType.Application.Json)
                setBody(request)
            }
        }
    }

    // Feed Management
    suspend fun getUserFeed(
        userId: String,
        limit: Int = 20,
        offset: Int = 0,
        feedType: String = "home"
    ): Result<SocialFeed> {
        return executeAuthenticatedRequest {
            httpClient.get("/mobile/feed/$userId") {
                parameter("limit", limit)
                parameter("offset", offset)
                parameter("type", feedType)
            }
        }
    }

    suspend fun getDiscoveryFeed(
        userId: String,
        region: String = "TH",
        limit: Int = 10
    ): Result<List<DiscoveryProfile>> {
        return executeAuthenticatedRequest {
            httpClient.get("/mobile/discover/$userId") {
                parameter("region", region)
                parameter("limit", limit)
            }
        }
    }

    suspend fun getTrendingContent(
        region: String = "TH",
        limit: Int = 10,
        timeWindow: String = "today"
    ): Result<RegionalTrending> {
        return executeAuthenticatedRequest {
            httpClient.get("/mobile/trending") {
                parameter("region", region)
                parameter("limit", limit)
                parameter("window", timeWindow)
            }
        }
    }

    // Post Management
    suspend fun createPost(request: CreatePostRequest): Result<SocialPost> {
        return executeAuthenticatedRequest {
            httpClient.post("/posts") {
                contentType(ContentType.Application.Json)
                setBody(request)
            }
        }
    }

    suspend fun getPost(postId: String): Result<SocialPost> {
        return executeAuthenticatedRequest {
            httpClient.get("/posts/$postId")
        }
    }

    suspend fun updatePost(postId: String, request: CreatePostRequest): Result<SocialPost> {
        return executeAuthenticatedRequest {
            httpClient.put("/posts/$postId") {
                contentType(ContentType.Application.Json)
                setBody(request)
            }
        }
    }

    suspend fun deletePost(postId: String): Result<Unit> {
        return executeAuthenticatedRequest<String> {
            httpClient.delete("/posts/$postId")
        }.map { Unit }
    }

    suspend fun getPostComments(
        postId: String,
        limit: Int = 20,
        offset: Int = 0
    ): Result<List<SocialComment>> {
        return executeAuthenticatedRequest {
            httpClient.get("/posts/$postId/comments") {
                parameter("limit", limit)
                parameter("offset", offset)
            }
        }
    }

    // Comments Management
    suspend fun createComment(request: CreateCommentRequest): Result<SocialComment> {
        return executeAuthenticatedRequest {
            httpClient.post("/comments") {
                contentType(ContentType.Application.Json)
                setBody(request)
            }
        }
    }

    suspend fun updateComment(commentId: String, content: String): Result<SocialComment> {
        return executeAuthenticatedRequest {
            httpClient.put("/comments/$commentId") {
                contentType(ContentType.Application.Json)
                setBody(mapOf("content" to content))
            }
        }
    }

    suspend fun deleteComment(commentId: String): Result<Unit> {
        return executeAuthenticatedRequest<String> {
            httpClient.delete("/comments/$commentId")
        }.map { Unit }
    }

    // Interactions (likes, follows, bookmarks)
    suspend fun createInteraction(request: InteractionRequest): Result<SocialInteraction> {
        return executeAuthenticatedRequest {
            httpClient.post("/interactions") {
                contentType(ContentType.Application.Json)
                setBody(request)
            }
        }
    }

    suspend fun removeInteraction(
        targetId: String,
        targetType: String,
        interactionType: String
    ): Result<Unit> {
        return executeAuthenticatedRequest<String> {
            httpClient.delete("/interactions") {
                parameter("targetId", targetId)
                parameter("targetType", targetType)
                parameter("interactionType", interactionType)
            }
        }.map { Unit }
    }

    suspend fun getInteractionState(
        targetId: String,
        targetType: String
    ): Result<List<SocialInteraction>> {
        return executeAuthenticatedRequest {
            httpClient.get("/interactions/state") {
                parameter("targetId", targetId)
                parameter("targetType", targetType)
            }
        }
    }

    // Follow Management
    suspend fun followUser(followingId: String): Result<Unit> {
        return executeAuthenticatedRequest<String> {
            httpClient.post("/mobile/follow") {
                contentType(ContentType.Application.Json)
                setBody(mapOf(
                    "followingId" to followingId,
                    "source" to "mobile_app"
                ))
            }
        }.map { Unit }
    }

    suspend fun unfollowUser(followingId: String): Result<Unit> {
        return executeAuthenticatedRequest<String> {
            httpClient.delete("/mobile/follow") {
                contentType(ContentType.Application.Json)
                setBody(mapOf("followingId" to followingId))
            }
        }.map { Unit }
    }

    suspend fun getFollowers(
        userId: String,
        limit: Int = 20,
        offset: Int = 0
    ): Result<List<SocialProfile>> {
        return executeAuthenticatedRequest {
            httpClient.get("/users/$userId/followers") {
                parameter("limit", limit)
                parameter("offset", offset)
            }
        }
    }

    suspend fun getFollowing(
        userId: String,
        limit: Int = 20,
        offset: Int = 0
    ): Result<List<SocialProfile>> {
        return executeAuthenticatedRequest {
            httpClient.get("/users/$userId/following") {
                parameter("limit", limit)
                parameter("offset", offset)
            }
        }
    }

    // Search
    suspend fun searchSocial(request: SocialSearchRequest): Result<SocialSearchResponse> {
        return executeAuthenticatedRequest {
            httpClient.get("/search") {
                parameter("query", request.query)
                parameter("type", request.type)
                parameter("region", request.region)
                parameter("limit", request.limit)
                parameter("offset", request.offset)
            }
        }
    }

    // Mobile Sync Operations
    suspend fun getProfileChanges(
        userId: String,
        since: String
    ): Result<SyncResponse> {
        return executeAuthenticatedRequest {
            httpClient.get("/mobile/sync/profile/$userId") {
                parameter("since", since)
            }
        }
    }

    suspend fun getPostChanges(
        userId: String,
        since: String
    ): Result<SyncResponse> {
        return executeAuthenticatedRequest {
            httpClient.get("/mobile/sync/posts/$userId") {
                parameter("since", since)
            }
        }
    }

    suspend fun getFollowChanges(
        userId: String,
        since: String
    ): Result<SyncResponse> {
        return executeAuthenticatedRequest {
            httpClient.get("/mobile/sync/follows/$userId") {
                parameter("since", since)
            }
        }
    }

    suspend fun applyClientChanges(
        userId: String,
        operations: List<SyncOperation>
    ): Result<SyncResponse> {
        return executeAuthenticatedRequest {
            httpClient.post("/mobile/apply/$userId") {
                contentType(ContentType.Application.Json)
                setBody(mapOf(
                    "userId" to userId,
                    "operations" to operations,
                    "timestamp" to Clock.System.now().toString()
                ))
            }
        }
    }

    suspend fun resolveConflicts(
        userId: String,
        conflicts: List<SyncConflict>,
        resolutions: Map<String, String>
    ): Result<SyncResponse> {
        return executeAuthenticatedRequest {
            httpClient.post("/mobile/resolve/$userId") {
                contentType(ContentType.Application.Json)
                setBody(mapOf(
                    "conflicts" to conflicts,
                    "resolutions" to resolutions,
                    "timestamp" to Clock.System.now().toString()
                ))
            }
        }
    }

    // Utility Methods
    private suspend inline fun <reified T> executeAuthenticatedRequest(
        crossinline request: suspend () -> HttpResponse
    ): Result<T> {
        return try {
            if (_connectionState.value != SocialConnectionState.CONNECTED) {
                return Result.failure(SocialApiException(503, "Not connected to social service"))
            }

            val response = request()

            if (response.status.isSuccess()) {
                val body = response.body<T>()
                Result.success(body)
            } else {
                val error = when (response.status.value) {
                    401 -> SocialApiException(401, "Unauthorized - please login again")
                    403 -> SocialApiException(403, "Access forbidden")
                    404 -> SocialApiException(404, "Resource not found")
                    429 -> SocialApiException(429, "Rate limit exceeded")
                    else -> SocialApiException(response.status.value, "Request failed: ${response.status.description}")
                }
                Result.failure(error)
            }
        } catch (e: SocialApiException) {
            Result.failure(e)
        } catch (e: Exception) {
            _connectionState.value = SocialConnectionState.ERROR
            Result.failure(SocialApiException(500, "Network error: ${e.message}"))
        }
    }

    private fun HttpRequestBuilder.authenticatedHeaders() {
        accessToken?.let { token ->
            headers {
                append(HttpHeaders.Authorization, "Bearer $token")
            }
        }
    }

    fun getServerTimestamp(): String {
        return Clock.System.now().toString()
    }

    fun isConnected(): Boolean {
        return _connectionState.value == SocialConnectionState.CONNECTED
    }
}

enum class SocialConnectionState {
    DISCONNECTED,
    CONNECTING,
    CONNECTED,
    ERROR
}

/**
 * Custom Social API Exception for handling HTTP errors
 */
class SocialApiException(
    val statusCode: Int,
    override val message: String,
    val requestId: String? = null
) : Exception(message)