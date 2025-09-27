package com.tchat.mobile.components

import android.content.Context
import android.view.ViewGroup
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import androidx.compose.ui.viewinterop.AndroidView
import androidx.media3.common.MediaItem
import androidx.media3.common.Player
import androidx.media3.exoplayer.ExoPlayer
import androidx.media3.ui.PlayerView
import com.tchat.mobile.designsystem.TchatColors

/**
 * Android implementation of TchatVideo using ExoPlayer
 */
@Composable
actual fun TchatVideo(
    source: VideoSource,
    modifier: Modifier,
    aspectRatio: TchatVideoAspectRatio,
    autoPlay: Boolean,
    loop: Boolean,
    muted: Boolean,
    showControls: Boolean,
    enableFullscreen: Boolean,
    enablePictureInPicture: Boolean,
    quality: TchatVideoQuality,
    startPosition: Long,
    onStateChange: ((TchatVideoState) -> Unit)?,
    onProgress: ((Long, Long) -> Unit)?,
    onError: ((String) -> Unit)?,
    poster: ImageSource?,
    subtitles: List<VideoSubtitle>
) {
    val context = LocalContext.current

    // Get video URL from source
    val videoUrl = when (source) {
        is VideoSource.Url -> source.url
        is VideoSource.Local -> source.resourceName  // For local assets
        is VideoSource.File -> source.path  // For file assets
        is VideoSource.Streaming -> source.manifestUrl  // For streaming
    }

    if (videoUrl.isNullOrEmpty()) {
        // Fallback UI for empty video URL
        Box(
            modifier = modifier
                .fillMaxWidth()
                .aspectRatio(
                    when (aspectRatio) {
                        TchatVideoAspectRatio.Auto -> 16f / 9f  // Default to landscape
                        TchatVideoAspectRatio.Portrait -> 9f / 16f
                        TchatVideoAspectRatio.Landscape -> 16f / 9f
                        TchatVideoAspectRatio.Square -> 1f
                        TchatVideoAspectRatio.Cinema -> 21f / 9f
                    }
                )
                .clip(RoundedCornerShape(8.dp)),
            contentAlignment = Alignment.Center
        ) {
            Card(
                modifier = Modifier.fillMaxSize(),
                colors = CardDefaults.cardColors(containerColor = TchatColors.surfaceVariant)
            ) {
                Column(
                    modifier = Modifier.fillMaxSize(),
                    horizontalAlignment = Alignment.CenterHorizontally,
                    verticalArrangement = Arrangement.Center
                ) {
                    Text("No video source", style = MaterialTheme.typography.bodyMedium)
                }
            }
        }
        return
    }

    // Create and remember ExoPlayer
    val exoPlayer = remember(videoUrl) {
        ExoPlayer.Builder(context).build().apply {
            val mediaItem = MediaItem.fromUri(videoUrl)
            setMediaItem(mediaItem)
            prepare()

            // Configure player settings
            volume = if (muted) 0f else 1f
            repeatMode = if (loop) Player.REPEAT_MODE_ONE else Player.REPEAT_MODE_OFF
            playWhenReady = autoPlay

            if (startPosition > 0) {
                seekTo(startPosition)
            }

            // Add listener for state changes
            addListener(object : Player.Listener {
                override fun onPlaybackStateChanged(playbackState: Int) {
                    val state = when (playbackState) {
                        Player.STATE_IDLE -> TchatVideoState.Idle
                        Player.STATE_BUFFERING -> TchatVideoState.Loading
                        Player.STATE_READY -> if (isPlaying) TchatVideoState.Playing else TchatVideoState.Paused
                        Player.STATE_ENDED -> TchatVideoState.Ended
                        else -> TchatVideoState.Idle
                    }
                    onStateChange?.invoke(state)
                }

                override fun onPlayerError(error: androidx.media3.common.PlaybackException) {
                    onError?.invoke(error.message ?: "Video playback error")
                }
            })
        }
    }

    // Progress tracking
    LaunchedEffect(exoPlayer, onProgress) {
        if (onProgress != null) {
            while (true) {
                val currentPosition = exoPlayer.currentPosition
                val duration = exoPlayer.duration
                if (duration > 0) {
                    onProgress(currentPosition, duration)
                }
                kotlinx.coroutines.delay(100) // Update every 100ms
            }
        }
    }

    // Dispose player when leaving composition
    DisposableEffect(exoPlayer) {
        onDispose {
            exoPlayer.release()
        }
    }

    // Create PlayerView
    AndroidView(
        factory = { ctx ->
            PlayerView(ctx).apply {
                player = exoPlayer
                layoutParams = ViewGroup.LayoutParams(
                    ViewGroup.LayoutParams.MATCH_PARENT,
                    ViewGroup.LayoutParams.MATCH_PARENT
                )
                useController = showControls
                resizeMode = androidx.media3.ui.AspectRatioFrameLayout.RESIZE_MODE_ZOOM
            }
        },
        modifier = modifier
            .fillMaxWidth()
            .aspectRatio(
                when (aspectRatio) {
                    TchatVideoAspectRatio.Auto -> 16f / 9f  // Default to landscape
                    TchatVideoAspectRatio.Portrait -> 9f / 16f
                    TchatVideoAspectRatio.Landscape -> 16f / 9f
                    TchatVideoAspectRatio.Square -> 1f
                    TchatVideoAspectRatio.Cinema -> 21f / 9f
                }
            )
            .clip(RoundedCornerShape(8.dp))
    )
}

/**
 * Android implementation of TchatVideoManager
 */
actual object TchatVideoManager {
    actual fun configure(config: VideoPlayerConfig) {
        // Implementation for configuring video settings
    }

    actual fun pauseAll() {
        // Implementation for pausing all videos
    }

    actual fun resumeAll() {
        // Implementation for resuming all videos
    }

    actual fun clearCache() {
        // Implementation for clearing video cache
    }
}