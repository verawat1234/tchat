package com.tchat.mobile.components

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier

/**
 * TchatVideo - Video player component with platform-native implementations
 *
 * Features:
 * - Play/pause controls with platform-native styling
 * - Progress bar and time display with scrubbing support
 * - Fullscreen mode with orientation handling
 * - Platform-native video handling (AVPlayer on iOS, ExoPlayer on Android)
 * - Volume control and mute functionality
 * - Subtitle support with multiple languages
 * - Video quality selection and adaptive streaming
 * - Picture-in-Picture mode support where available
 * - Accessibility support with media session integration
 */

enum class TchatVideoState {
    Idle,       // Not loaded
    Loading,    // Loading video
    Ready,      // Ready to play
    Playing,    // Currently playing
    Paused,     // Paused by user
    Buffering,  // Buffering content
    Ended,      // Playback completed
    Error       // Error occurred
}

enum class TchatVideoQuality {
    Auto,       // Automatic quality selection
    Low,        // Low quality (240p-360p)
    Medium,     // Medium quality (480p-720p)
    High,       // High quality (1080p+)
    Original    // Original upload quality
}

enum class TchatVideoAspectRatio {
    Auto,       // Maintain original aspect ratio
    Square,     // 1:1 aspect ratio
    Portrait,   // 9:16 aspect ratio
    Landscape,  // 16:9 aspect ratio
    Cinema      // 21:9 aspect ratio
}

/**
 * Cross-platform video player component using expect/actual pattern
 * Platform-specific implementations provide native video playback
 */
@Composable
expect fun TchatVideo(
    source: VideoSource,
    modifier: Modifier = Modifier,
    aspectRatio: TchatVideoAspectRatio = TchatVideoAspectRatio.Auto,
    autoPlay: Boolean = false,
    loop: Boolean = false,
    muted: Boolean = false,
    showControls: Boolean = true,
    enableFullscreen: Boolean = true,
    enablePictureInPicture: Boolean = false,
    quality: TchatVideoQuality = TchatVideoQuality.Auto,
    startPosition: Long = 0L,
    onStateChange: ((TchatVideoState) -> Unit)? = null,
    onProgress: ((Long, Long) -> Unit)? = null, // current, duration
    onError: ((String) -> Unit)? = null,
    poster: ImageSource? = null,
    subtitles: List<VideoSubtitle> = emptyList()
)

/**
 * Image source configuration supporting multiple formats
 */
sealed class ImageSource {
    data class Url(val url: String, val headers: Map<String, String> = emptyMap()) : ImageSource()
    data class Local(val resourceName: String) : ImageSource()
    data class File(val path: String) : ImageSource()
}

/**
 * Video source configuration supporting multiple formats
 */
sealed class VideoSource {
    data class Url(val url: String, val headers: Map<String, String> = emptyMap()) : VideoSource()
    data class Local(val resourceName: String) : VideoSource()
    data class File(val path: String) : VideoSource()
    data class Streaming(
        val manifestUrl: String,
        val type: StreamingType = StreamingType.HLS
    ) : VideoSource()
}

/**
 * Streaming protocol types
 */
enum class StreamingType {
    HLS,        // HTTP Live Streaming
    DASH,       // Dynamic Adaptive Streaming
    RTMP,       // Real-Time Messaging Protocol
    WebRTC      // Web Real-Time Communication
}

/**
 * Video subtitle configuration
 */
data class VideoSubtitle(
    val url: String,
    val language: String,
    val label: String,
    val isDefault: Boolean = false
)

/**
 * Video player controls configuration
 */
data class VideoControlsConfig(
    val showPlayPause: Boolean = true,
    val showProgress: Boolean = true,
    val showTime: Boolean = true,
    val showVolume: Boolean = true,
    val showFullscreen: Boolean = true,
    val showQuality: Boolean = true,
    val showSubtitles: Boolean = true,
    val autoHideDelay: Long = 3000L, // ms
    val seekStep: Long = 10000L // 10 seconds
)

/**
 * Video playback statistics
 */
data class VideoPlaybackStats(
    val duration: Long,
    val currentPosition: Long,
    val bufferedPosition: Long,
    val playbackSpeed: Float,
    val videoWidth: Int,
    val videoHeight: Int,
    val bitrateKbps: Int,
    val droppedFrames: Int
)

/**
 * Global video player manager
 */
expect object TchatVideoManager {
    /**
     * Configure global video settings
     */
    fun configure(config: VideoPlayerConfig)

    /**
     * Pause all video players (useful for lifecycle management)
     */
    fun pauseAll()

    /**
     * Resume video players
     */
    fun resumeAll()

    /**
     * Clear video cache
     */
    fun clearCache()
}

/**
 * Global video player configuration
 */
data class VideoPlayerConfig(
    val cacheSize: Long = 100 * 1024 * 1024, // 100MB
    val enableCache: Boolean = true,
    val enableHardwareAcceleration: Boolean = true,
    val maxBitrate: Int = 0, // 0 = unlimited
    val preferredAudioLanguage: String? = null,
    val preferredSubtitleLanguage: String? = null
)