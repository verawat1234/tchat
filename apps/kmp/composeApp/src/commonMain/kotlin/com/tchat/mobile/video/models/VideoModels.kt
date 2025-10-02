package com.tchat.mobile.video.models

import kotlinx.serialization.Serializable

/**
 * Video Content Model
 *
 * Represents video metadata and content information across platforms.
 * Matches backend video_content schema from backend/video/models/video_content.go
 */
@Serializable
data class VideoContent(
    val id: String,
    val creatorId: String,
    val title: String,
    val description: String,
    val thumbnailUrl: String? = null,
    val videoUrl: String,
    val durationSeconds: Int,
    val fileSize: Long,
    val mimeType: String,
    val resolution: String,
    val bitrate: Int,
    val framerate: Double,
    val codec: String,
    val uploadStatus: VideoStatus = VideoStatus.PROCESSING,
    val availabilityStatus: VideoStatus = VideoStatus.AVAILABLE,
    val contentRating: ContentRating = ContentRating.GENERAL,
    val tags: List<String> = emptyList(),
    val category: String? = null,
    val viewCount: Long = 0,
    val likeCount: Long = 0,
    val commentCount: Long = 0,
    val shareCount: Long = 0,
    val isMonetized: Boolean = false,
    val price: Double? = null,
    val currency: String? = null,
    val socialMetrics: SocialMetrics = SocialMetrics(),
    val createdAt: String,
    val updatedAt: String,
    val publishedAt: String? = null,
    val deletedAt: String? = null
)

/**
 * Video Status Enum
 *
 * Represents the current state of video processing and availability.
 */
@Serializable
enum class VideoStatus {
    PROCESSING,
    AVAILABLE,
    UNAVAILABLE,
    ARCHIVED,
    DELETED
}

/**
 * Content Rating Enum
 *
 * Age-appropriate content rating system.
 */
@Serializable
enum class ContentRating {
    GENERAL,
    TEEN,
    MATURE,
    ADULT
}

/**
 * Social Metrics Model
 *
 * Engagement metrics for video content.
 */
@Serializable
data class SocialMetrics(
    val viewCount: Long = 0,
    val likeCount: Long = 0,
    val dislikeCount: Long = 0,
    val commentCount: Long = 0,
    val shareCount: Long = 0,
    val favoriteCount: Long = 0
)

/**
 * Playback Session Model
 *
 * Represents an active video playback session for cross-platform sync.
 * Matches backend playback_session schema from backend/video/models/playback_session.go
 */
@Serializable
data class PlaybackSession(
    val id: String,
    val videoId: String,
    val userId: String,
    val platform: PlatformType,
    val deviceId: String,
    val currentPosition: Double,
    val totalDuration: Double,
    val playbackSpeed: Double = 1.0,
    val quality: VideoQuality = VideoQuality.AUTO,
    val volume: Double = 1.0,
    val state: PlaybackState = PlaybackState.PAUSED,
    val bufferHealth: BufferHealth = BufferHealth(),
    val isFullscreen: Boolean = false,
    val isPictureInPicture: Boolean = false,
    val lastSyncTime: String,
    val createdAt: String,
    val updatedAt: String
)

/**
 * Platform Type Enum
 *
 * Identifies the platform where video is being played.
 */
@Serializable
enum class PlatformType {
    WEB,
    IOS,
    ANDROID,
    MOBILE_WEB
}

/**
 * Playback State Enum
 *
 * Current playback state of the video player.
 */
@Serializable
enum class PlaybackState {
    PLAYING,
    PAUSED,
    BUFFERING,
    ENDED,
    ERROR
}

/**
 * Video Quality Enum
 *
 * Available video quality options.
 */
@Serializable
enum class VideoQuality {
    AUTO,
    LOW,      // 360p
    MEDIUM,   // 720p
    HIGH,     // 1080p
    ULTRA     // 4K
}

/**
 * Buffer Health Model
 *
 * Tracks video buffer status for adaptive streaming.
 */
@Serializable
data class BufferHealth(
    val bufferedSeconds: Double = 0.0,
    val bufferPercentage: Double = 0.0,
    val isStalled: Boolean = false,
    val lastStallTime: String? = null
)

/**
 * Viewing History Model
 *
 * Tracks user's video viewing history for recommendations.
 * Matches backend viewing_history schema from backend/video/models/viewing_history.go
 */
@Serializable
data class ViewingHistory(
    val id: String,
    val userId: String,
    val videoId: String,
    val watchedSeconds: Double,
    val completionPercentage: Double,
    val lastWatchedPosition: Double,
    val watchCount: Int = 1,
    val platform: PlatformType,
    val deviceId: String,
    val isCompleted: Boolean = false,
    val createdAt: String,
    val updatedAt: String,
    val lastWatchedAt: String
)

/**
 * Synchronization State Model
 *
 * Manages cross-platform sync state and conflict resolution.
 * Matches backend sync_state schema from backend/video/models/sync_state.go
 */
@Serializable
data class SyncState(
    val id: String,
    val userId: String,
    val videoId: String,
    val sessionId: String,
    val syncedPlatforms: List<PlatformType>,
    val lastSyncTime: String,
    val conflictDetected: Boolean = false,
    val conflictResolution: ConflictResolutionStrategy? = null,
    val syncVersion: Int = 1,
    val pendingChanges: List<String> = emptyList(),
    val createdAt: String,
    val updatedAt: String
)

/**
 * Conflict Resolution Strategy Enum
 *
 * Strategy for resolving sync conflicts between platforms.
 */
@Serializable
enum class ConflictResolutionStrategy {
    LATEST_WINS,
    AUTHORITY_PLATFORM,
    AVERAGE_POSITION,
    MANUAL_RESOLUTION
}

/**
 * Video Upload Request Model
 *
 * Request model for video upload operations.
 */
@Serializable
data class VideoUploadRequest(
    val title: String,
    val description: String,
    val tags: List<String> = emptyList(),
    val contentRating: ContentRating = ContentRating.GENERAL,
    val category: String? = null,
    val isMonetized: Boolean = false,
    val price: Double? = null,
    val currency: String? = null
)

/**
 * Video Upload Response Model
 *
 * Response model containing upload results.
 */
@Serializable
data class VideoUploadResponse(
    val videoId: String,
    val status: String,
    val message: String,
    val uploadUrl: String? = null,
    val thumbnailUploadUrl: String? = null
)

/**
 * Stream URL Request Model
 *
 * Request model for obtaining video stream URLs.
 */
@Serializable
data class StreamURLRequest(
    val videoId: String,
    val quality: VideoQuality = VideoQuality.AUTO,
    val platform: PlatformType
)

/**
 * Stream URL Response Model
 *
 * Response containing streaming URLs and available qualities.
 */
@Serializable
data class StreamURLResponse(
    val videoId: String,
    val streamUrl: String,
    val manifestUrl: String,
    val availableQualities: List<QualityOption>,
    val expiresAt: String,
    val cdn: String
)

/**
 * Quality Option Model
 *
 * Available quality option with metadata.
 */
@Serializable
data class QualityOption(
    val quality: VideoQuality,
    val resolution: String,
    val bitrate: Int,
    val url: String
)

/**
 * Sync Playback Request Model
 *
 * Request model for cross-platform playback synchronization.
 */
@Serializable
data class SyncPlaybackRequest(
    val videoId: String,
    val sessionId: String,
    val position: Double,
    val platform: PlatformType,
    val playbackState: PlaybackState
)

/**
 * Sync Playback Response Model
 *
 * Response containing sync status and conflict information.
 */
@Serializable
data class SyncPlaybackResponse(
    val videoId: String,
    val sessionId: String,
    val syncedPlatforms: List<PlatformType>,
    val latency: Long,
    val conflictDetected: Boolean = false,
    val conflictResolution: ConflictResolutionStrategy? = null,
    val timestamp: String
)

/**
 * Video Error Model
 *
 * Error information for failed video operations.
 */
@Serializable
data class VideoError(
    val code: String,
    val message: String,
    val details: Map<String, String>? = null
)