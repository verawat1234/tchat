package com.tchat.mobile.video.platform

import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.interop.UIKitView
import com.tchat.mobile.video.models.PlaybackState
import com.tchat.mobile.video.models.VideoQuality
import kotlinx.cinterop.*
import platform.AVFoundation.*
import platform.AVKit.AVPlayerViewController
import platform.CoreGraphics.CGRectMake
import platform.Foundation.NSURL
import platform.Foundation.addObserver
import platform.Foundation.removeObserver
import platform.UIKit.UIView

/**
 * iOS Video Player Component
 *
 * AVPlayer-based video player implementation for iOS.
 * Provides native video playback with HLS support.
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
    val player = remember { createAVPlayer(videoUrl) }
    val playerViewController = remember { AVPlayerViewController() }

    DisposableEffect(videoUrl) {
        val url = NSURL.URLWithString(videoUrl)
        val playerItem = url?.let { AVPlayerItem(uRL = it) }

        playerItem?.let {
            player.replaceCurrentItemWithPlayerItem(it)

            if (autoPlay) {
                player.play()
            }
        }

        onDispose {
            player.pause()
            player.replaceCurrentItemWithPlayerItem(null)
        }
    }

    // Configure player settings
    LaunchedEffect(muted) {
        player.muted = muted
    }

    LaunchedEffect(loop) {
        // iOS AVPlayer doesn't have built-in loop, need to observe playback end
        // and restart manually
    }

    // Monitor playback state
    LaunchedEffect(player) {
        // Observe player status
        player.addObserver(
            observer = object : NSObject() {
                override fun observeValueForKeyPath(
                    keyPath: String?,
                    ofObject: Any?,
                    change: Map<Any?, *>?,
                    context: COpaquePointer?
                ) {
                    when (keyPath) {
                        "status" -> {
                            val status = player.status
                            when (status) {
                                AVPlayerStatusReadyToPlay -> {
                                    if (player.rate > 0f) {
                                        onPlaybackStateChange(PlaybackState.PLAYING)
                                    } else {
                                        onPlaybackStateChange(PlaybackState.PAUSED)
                                    }
                                }
                                AVPlayerStatusFailed -> {
                                    onError(player.error?.localizedDescription ?: "Playback failed")
                                    onPlaybackStateChange(PlaybackState.ERROR)
                                }
                                else -> {
                                    onPlaybackStateChange(PlaybackState.BUFFERING)
                                }
                            }
                        }
                        "rate" -> {
                            if (player.rate > 0f) {
                                onPlaybackStateChange(PlaybackState.PLAYING)
                            } else {
                                onPlaybackStateChange(PlaybackState.PAUSED)
                            }
                        }
                    }
                }
            },
            forKeyPath = "status",
            options = 0u,
            context = null
        )

        player.addObserver(
            observer = object : NSObject() {},
            forKeyPath = "rate",
            options = 0u,
            context = null
        )
    }

    // Update position periodically
    LaunchedEffect(player) {
        while (true) {
            if (player.rate > 0f) {
                val currentTime = player.currentTime()
                val position = CMTimeGetSeconds(currentTime)
                onPositionChange(position)

                player.currentItem?.duration?.let { duration ->
                    val durationSeconds = CMTimeGetSeconds(duration)
                    if (durationSeconds > 0) {
                        onDurationChange(durationSeconds)
                    }
                }
            }
            kotlinx.coroutines.delay(500) // Update every 500ms
        }
    }

    // Render AVPlayerViewController
    UIKitView(
        factory = {
            playerViewController.player = player
            playerViewController.showsPlaybackControls = false // Use custom controls
            playerViewController.view
        },
        modifier = modifier
    )
}

/**
 * Create AVPlayer instance with optimized configuration
 */
private fun createAVPlayer(videoUrl: String): AVPlayer {
    val url = NSURL.URLWithString(videoUrl)
    val playerItem = url?.let { AVPlayerItem(uRL = it) }

    return AVPlayer(playerItem = playerItem).apply {
        // Configure for adaptive streaming
        appliesMediaSelectionCriteriaAutomatically = true

        // Enable background playback if needed
        allowsExternalPlayback = true
    }
}

/**
 * Video Player Control Interface
 *
 * Provides control methods for the video player.
 */
class IOSVideoPlayerController(private val player: AVPlayer) {
    fun play() = player.play()

    fun pause() = player.pause()

    fun seekTo(positionSeconds: Double) {
        val time = CMTimeMakeWithSeconds(positionSeconds, preferredTimescale = 1)
        player.seekToTime(time)
    }

    fun setVolume(volume: Float) {
        player.volume = volume.coerceIn(0f, 1f)
    }

    fun setPlaybackSpeed(speed: Float) {
        player.rate = speed
    }

    fun getCurrentPosition(): Double {
        val currentTime = player.currentTime()
        return CMTimeGetSeconds(currentTime)
    }

    fun getDuration(): Double {
        return player.currentItem?.duration?.let { duration ->
            CMTimeGetSeconds(duration)
        } ?: 0.0
    }

    fun isPlaying(): Boolean = player.rate > 0f

    fun release() {
        player.pause()
        player.replaceCurrentItemWithPlayerItem(null)
    }
}

/**
 * Helper function to create CMTime from seconds
 */
private fun CMTimeMakeWithSeconds(seconds: Double, preferredTimescale: Int): platform.CoreMedia.CMTime {
    return platform.CoreMedia.CMTimeMakeWithSeconds(seconds, preferredTimescale)
}

/**
 * Helper function to get seconds from CMTime
 */
private fun CMTimeGetSeconds(time: platform.CoreMedia.CMTime): Double {
    return platform.CoreMedia.CMTimeGetSeconds(time)
}