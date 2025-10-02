package com.tchat.mobile.video.platform

import android.content.Context
import android.view.ViewGroup
import android.widget.FrameLayout
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.viewinterop.AndroidView
import androidx.media3.common.MediaItem
import androidx.media3.common.Player
import androidx.media3.exoplayer.ExoPlayer
import androidx.media3.ui.PlayerView
import com.tchat.mobile.video.models.PlaybackState
import com.tchat.mobile.video.models.VideoQuality

/**
 * Android Video Player Component
 *
 * ExoPlayer-based video player implementation for Android.
 * Provides native video playback with HLS/DASH support.
 */
@Composable
actual fun PlatformVideoPlayer(
    videoUrl: String,
    autoPlay: Boolean,
    muted: Boolean,
    loop: Boolean,
    quality: VideoQuality,
    onPlaybackStateChange: (PlaybackState) -> Unit,
    onPositionChange: (Double) -> Unit,
    onDurationChange: (Double) -> Unit,
    onError: (String) -> Unit,
    modifier: Modifier
) {
    val context = LocalContext.current
    val exoPlayer = remember { createExoPlayer(context) }

    DisposableEffect(videoUrl) {
        val mediaItem = MediaItem.fromUri(videoUrl)
        exoPlayer.setMediaItem(mediaItem)
        exoPlayer.prepare()

        if (autoPlay) {
            exoPlayer.play()
        }

        onDispose {
            exoPlayer.release()
        }
    }

    LaunchedEffect(exoPlayer) {
        // Listen to playback state changes
        val listener = object : Player.Listener {
            override fun onPlaybackStateChanged(playbackState: Int) {
                val state = when (playbackState) {
                    Player.STATE_IDLE -> PlaybackState.PAUSED
                    Player.STATE_BUFFERING -> PlaybackState.BUFFERING
                    Player.STATE_READY -> {
                        if (exoPlayer.isPlaying) PlaybackState.PLAYING
                        else PlaybackState.PAUSED
                    }
                    Player.STATE_ENDED -> PlaybackState.ENDED
                    else -> PlaybackState.PAUSED
                }
                onPlaybackStateChange(state)
            }

            override fun onPlayerError(error: androidx.media3.common.PlaybackException) {
                onError(error.message ?: "Unknown playback error")
                onPlaybackStateChange(PlaybackState.ERROR)
            }
        }

        exoPlayer.addListener(listener)
    }

    // Update position periodically
    LaunchedEffect(exoPlayer) {
        while (true) {
            if (exoPlayer.isPlaying) {
                val position = exoPlayer.currentPosition / 1000.0 // Convert to seconds
                onPositionChange(position)

                val duration = exoPlayer.duration / 1000.0
                if (duration > 0) {
                    onDurationChange(duration)
                }
            }
            kotlinx.coroutines.delay(500) // Update every 500ms
        }
    }

    // Configure ExoPlayer settings
    LaunchedEffect(muted) {
        exoPlayer.volume = if (muted) 0f else 1f
    }

    LaunchedEffect(loop) {
        exoPlayer.repeatMode = if (loop) Player.REPEAT_MODE_ONE else Player.REPEAT_MODE_OFF
    }

    // Render PlayerView
    AndroidView(
        factory = { ctx ->
            PlayerView(ctx).apply {
                player = exoPlayer
                layoutParams = FrameLayout.LayoutParams(
                    ViewGroup.LayoutParams.MATCH_PARENT,
                    ViewGroup.LayoutParams.MATCH_PARENT
                )
                useController = false // Use custom controls from Compose
            }
        },
        modifier = modifier
    )
}

/**
 * Create ExoPlayer instance with optimized configuration
 */
private fun createExoPlayer(context: Context): ExoPlayer {
    return ExoPlayer.Builder(context)
        .build()
        .apply {
            // Configure for adaptive streaming
            videoScalingMode = androidx.media3.common.C.VIDEO_SCALING_MODE_SCALE_TO_FIT_WITH_CROPPING

            // Enable audio focus handling
            setAudioAttributes(
                androidx.media3.common.AudioAttributes.Builder()
                    .setContentType(androidx.media3.common.C.AUDIO_CONTENT_TYPE_MOVIE)
                    .setUsage(androidx.media3.common.C.USAGE_MEDIA)
                    .build(),
                true
            )
        }
}

/**
 * Video Player Control Interface
 *
 * Provides control methods for the video player.
 */
class AndroidVideoPlayerController(private val player: ExoPlayer) {
    fun play() = player.play()
    fun pause() = player.pause()
    fun seekTo(positionMs: Long) = player.seekTo(positionMs)
    fun setVolume(volume: Float) {
        player.volume = volume.coerceIn(0f, 1f)
    }
    fun setPlaybackSpeed(speed: Float) {
        player.setPlaybackSpeed(speed)
    }
    fun getCurrentPosition(): Long = player.currentPosition
    fun getDuration(): Long = player.duration
    fun isPlaying(): Boolean = player.isPlaying
    fun release() = player.release()
}