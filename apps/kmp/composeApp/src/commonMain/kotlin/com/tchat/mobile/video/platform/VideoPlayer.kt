// apps/kmp/composeApp/src/commonMain/kotlin/com/tchat/mobile/video/platform/VideoPlayer.kt
// Common expect declaration for platform-specific video player implementations

package com.tchat.mobile.video.platform

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import com.tchat.mobile.video.models.PlaybackState
import com.tchat.mobile.video.models.VideoQuality

/**
 * Platform Video Player - Expect Declaration
 *
 * This is the common interface for platform-specific video player implementations.
 * Android uses ExoPlayer, iOS uses AVPlayer.
 *
 * @param videoUrl URL of the video to play
 * @param autoPlay Whether to start playing automatically
 * @param muted Whether to mute audio
 * @param loop Whether to loop playback
 * @param quality Video quality setting
 * @param onPlaybackStateChange Callback for playback state changes
 * @param onPositionChange Callback for position updates (in seconds)
 * @param onDurationChange Callback for duration updates (in seconds)
 * @param onError Callback for error events
 * @param modifier Compose modifier for styling
 */
@Composable
expect fun PlatformVideoPlayer(
    videoUrl: String,
    autoPlay: Boolean = false,
    muted: Boolean = false,
    loop: Boolean = false,
    quality: VideoQuality = VideoQuality.AUTO,
    onPlaybackStateChange: (PlaybackState) -> Unit = {},
    onPositionChange: (Double) -> Unit = {},
    onDurationChange: (Double) -> Unit = {},
    onError: (String) -> Unit = {},
    modifier: Modifier = Modifier
)