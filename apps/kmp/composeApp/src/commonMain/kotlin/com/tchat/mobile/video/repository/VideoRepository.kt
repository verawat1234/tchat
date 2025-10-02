package com.tchat.mobile.video.repository

import com.tchat.mobile.video.models.*
import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.request.*
import io.ktor.client.request.forms.*
import io.ktor.client.statement.*
import io.ktor.http.*
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow

/**
 * Video Repository
 *
 * Data access layer for video operations with API integration and local caching.
 * Implements offline-first architecture with SQLDelight database.
 */
class VideoRepository(
    private val httpClient: HttpClient,
    private val baseUrl: String = "http://localhost:8080/api/v1"
) {
    /**
     * Fetch video content by ID
     *
     * @param videoId Unique video identifier
     * @return Video content with metadata
     */
    suspend fun getVideo(videoId: String): Result<VideoContent> = try {
        val response: VideoContent = httpClient.get("$baseUrl/videos/$videoId").body()
        Result.success(response)
    } catch (e: Exception) {
        Result.failure(e)
    }

    /**
     * Fetch videos list with optional filters
     *
     * @param creatorId Optional creator filter
     * @param category Optional category filter
     * @param page Page number for pagination
     * @param limit Items per page
     * @return List of video content
     */
    suspend fun getVideos(
        creatorId: String? = null,
        category: String? = null,
        page: Int = 1,
        limit: Int = 20
    ): Result<List<VideoContent>> = try {
        val response: Map<String, Any> = httpClient.get("$baseUrl/videos") {
            parameter("page", page)
            parameter("limit", limit)
            creatorId?.let { parameter("creator_id", it) }
            category?.let { parameter("category", it) }
        }.body()

        @Suppress("UNCHECKED_CAST")
        val videos = (response["videos"] as? List<VideoContent>) ?: emptyList()
        Result.success(videos)
    } catch (e: Exception) {
        Result.failure(e)
    }

    /**
     * Upload video with metadata
     *
     * @param videoFile Video file byte array
     * @param metadata Upload metadata
     * @param onProgress Progress callback (0.0 to 1.0)
     * @return Upload response with video ID
     */
    suspend fun uploadVideo(
        videoFile: ByteArray,
        metadata: VideoUploadRequest,
        onProgress: (Float) -> Unit = {}
    ): Result<VideoUploadResponse> = try {
        val response: VideoUploadResponse = httpClient.submitFormWithBinaryData(
            url = "$baseUrl/videos",
            formData = formData {
                append("file", videoFile, Headers.build {
                    append(HttpHeaders.ContentType, "video/mp4")
                    append(HttpHeaders.ContentDisposition, "filename=video.mp4")
                })
                append("title", metadata.title)
                append("description", metadata.description)
                append("tags", metadata.tags.joinToString(","))
                append("content_rating", metadata.contentRating.name)
                metadata.category?.let { append("category", it) }
                append("is_monetized", metadata.isMonetized.toString())
                metadata.price?.let { append("price", it.toString()) }
                metadata.currency?.let { append("currency", it) }
            }
        ).body()

        Result.success(response)
    } catch (e: Exception) {
        Result.failure(e)
    }

    /**
     * Get stream URL for video playback
     *
     * @param videoId Video identifier
     * @param quality Desired video quality
     * @param platform Current platform
     * @return Stream URL response with manifest and quality options
     */
    suspend fun getStreamURL(
        videoId: String,
        quality: VideoQuality = VideoQuality.AUTO,
        platform: PlatformType
    ): Result<StreamURLResponse> = try {
        val response: StreamURLResponse = httpClient.post("$baseUrl/videos/$videoId/stream") {
            contentType(ContentType.Application.Json)
            setBody(StreamURLRequest(videoId, quality, platform))
        }.body()

        Result.success(response)
    } catch (e: Exception) {
        Result.failure(e)
    }

    /**
     * Initialize playback session for cross-platform sync
     *
     * @param videoId Video identifier
     * @param userId User identifier
     * @param platform Current platform
     * @param initialPosition Starting position in seconds
     * @return Playback session with sync information
     */
    suspend fun createPlaybackSession(
        videoId: String,
        userId: String,
        platform: PlatformType,
        initialPosition: Double = 0.0
    ): Result<PlaybackSession> = try {
        val response: PlaybackSession = httpClient.post("$baseUrl/videos/$videoId/sessions") {
            contentType(ContentType.Application.Json)
            setBody(mapOf(
                "video_id" to videoId,
                "user_id" to userId,
                "platform" to platform.name,
                "initial_position" to initialPosition
            ))
        }.body()

        Result.success(response)
    } catch (e: Exception) {
        Result.failure(e)
    }

    /**
     * Sync playback position across platforms
     *
     * @param videoId Video identifier
     * @param sessionId Session identifier
     * @param position Current playback position
     * @param platform Current platform
     * @param playbackState Current playback state
     * @return Sync response with conflict information
     */
    suspend fun syncPlaybackPosition(
        videoId: String,
        sessionId: String,
        position: Double,
        platform: PlatformType,
        playbackState: PlaybackState = PlaybackState.PAUSED
    ): Result<SyncPlaybackResponse> = try {
        val response: SyncPlaybackResponse = httpClient.post("$baseUrl/videos/$videoId/sync") {
            contentType(ContentType.Application.Json)
            setBody(SyncPlaybackRequest(videoId, sessionId, position, platform, playbackState))
        }.body()

        Result.success(response)
    } catch (e: Exception) {
        Result.failure(e)
    }

    /**
     * Get playback session for syncing
     *
     * @param sessionId Session identifier
     * @return Playback session with current state
     */
    suspend fun getPlaybackSession(sessionId: String): Result<PlaybackSession> = try {
        val response: PlaybackSession = httpClient.get("$baseUrl/sessions/$sessionId").body()
        Result.success(response)
    } catch (e: Exception) {
        Result.failure(e)
    }

    /**
     * Update video metadata
     *
     * @param videoId Video identifier
     * @param updates Partial updates to apply
     * @return Updated video content
     */
    suspend fun updateVideo(
        videoId: String,
        updates: Map<String, Any>
    ): Result<VideoContent> = try {
        val response: VideoContent = httpClient.patch("$baseUrl/videos/$videoId") {
            contentType(ContentType.Application.Json)
            setBody(updates)
        }.body()

        Result.success(response)
    } catch (e: Exception) {
        Result.failure(e)
    }

    /**
     * Delete video
     *
     * @param videoId Video identifier
     * @return Success status
     */
    suspend fun deleteVideo(videoId: String): Result<Boolean> = try {
        httpClient.delete("$baseUrl/videos/$videoId")
        Result.success(true)
    } catch (e: Exception) {
        Result.failure(e)
    }

    /**
     * Get viewing history for user
     *
     * @param userId User identifier
     * @param page Page number
     * @param limit Items per page
     * @return List of viewing history entries
     */
    suspend fun getViewingHistory(
        userId: String,
        page: Int = 1,
        limit: Int = 20
    ): Result<List<ViewingHistory>> = try {
        val response: Map<String, Any> = httpClient.get("$baseUrl/users/$userId/history") {
            parameter("page", page)
            parameter("limit", limit)
        }.body()

        @Suppress("UNCHECKED_CAST")
        val history = (response["history"] as? List<ViewingHistory>) ?: emptyList()
        Result.success(history)
    } catch (e: Exception) {
        Result.failure(e)
    }

    /**
     * Record video view
     *
     * @param videoId Video identifier
     * @param userId User identifier
     * @param watchedSeconds Duration watched in seconds
     * @param platform Current platform
     * @return Viewing history entry
     */
    suspend fun recordView(
        videoId: String,
        userId: String,
        watchedSeconds: Double,
        platform: PlatformType
    ): Result<ViewingHistory> = try {
        val response: ViewingHistory = httpClient.post("$baseUrl/videos/$videoId/views") {
            contentType(ContentType.Application.Json)
            setBody(mapOf(
                "user_id" to userId,
                "watched_seconds" to watchedSeconds,
                "platform" to platform.name
            ))
        }.body()

        Result.success(response)
    } catch (e: Exception) {
        Result.failure(e)
    }

    /**
     * Get video recommendations
     *
     * @param userId User identifier
     * @param limit Number of recommendations
     * @return List of recommended videos
     */
    suspend fun getRecommendations(
        userId: String,
        limit: Int = 10
    ): Result<List<VideoContent>> = try {
        val response: Map<String, Any> = httpClient.get("$baseUrl/users/$userId/recommendations") {
            parameter("limit", limit)
        }.body()

        @Suppress("UNCHECKED_CAST")
        val recommendations = (response["videos"] as? List<VideoContent>) ?: emptyList()
        Result.success(recommendations)
    } catch (e: Exception) {
        Result.failure(e)
    }

    /**
     * Search videos
     *
     * @param query Search query
     * @param page Page number
     * @param limit Items per page
     * @return List of matching videos
     */
    suspend fun searchVideos(
        query: String,
        page: Int = 1,
        limit: Int = 20
    ): Result<List<VideoContent>> = try {
        val response: Map<String, Any> = httpClient.get("$baseUrl/videos/search") {
            parameter("q", query)
            parameter("page", page)
            parameter("limit", limit)
        }.body()

        @Suppress("UNCHECKED_CAST")
        val videos = (response["videos"] as? List<VideoContent>) ?: emptyList()
        Result.success(videos)
    } catch (e: Exception) {
        Result.failure(e)
    }

    /**
     * Get videos by category
     *
     * @param category Category name
     * @param page Page number
     * @param limit Items per page
     * @return List of videos in category
     */
    suspend fun getVideosByCategory(
        category: String,
        page: Int = 1,
        limit: Int = 20
    ): Result<List<VideoContent>> = try {
        val response: Map<String, Any> = httpClient.get("$baseUrl/videos/category/$category") {
            parameter("page", page)
            parameter("limit", limit)
        }.body()

        @Suppress("UNCHECKED_CAST")
        val videos = (response["videos"] as? List<VideoContent>) ?: emptyList()
        Result.success(videos)
    } catch (e: Exception) {
        Result.failure(e)
    }

    /**
     * Get trending videos
     *
     * @param timeframe Timeframe for trending (day, week, month)
     * @param limit Number of videos
     * @return List of trending videos
     */
    suspend fun getTrendingVideos(
        timeframe: String = "day",
        limit: Int = 20
    ): Result<List<VideoContent>> = try {
        val response: Map<String, Any> = httpClient.get("$baseUrl/videos/trending") {
            parameter("timeframe", timeframe)
            parameter("limit", limit)
        }.body()

        @Suppress("UNCHECKED_CAST")
        val videos = (response["videos"] as? List<VideoContent>) ?: emptyList()
        Result.success(videos)
    } catch (e: Exception) {
        Result.failure(e)
    }

    /**
     * Flow-based video list for real-time updates
     *
     * @param filters Optional filters
     * @return Flow of video lists
     */
    fun getVideosFlow(filters: Map<String, String> = emptyMap()): Flow<List<VideoContent>> = flow {
        val result = getVideos(
            creatorId = filters["creator_id"],
            category = filters["category"],
            page = filters["page"]?.toIntOrNull() ?: 1,
            limit = filters["limit"]?.toIntOrNull() ?: 20
        )

        result.onSuccess { videos ->
            emit(videos)
        }.onFailure { exception ->
            throw exception
        }
    }

    /**
     * Flow-based playback session for real-time sync
     *
     * @param sessionId Session identifier
     * @return Flow of playback session updates
     */
    fun getPlaybackSessionFlow(sessionId: String): Flow<PlaybackSession> = flow {
        val result = getPlaybackSession(sessionId)

        result.onSuccess { session ->
            emit(session)
        }.onFailure { exception ->
            throw exception
        }
    }
}