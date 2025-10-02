package com.tchat.mobile.streaming

import kotlinx.coroutines.flow.StateFlow

/**
 * Cross-platform WebRTC client for live streaming
 *
 * Provides platform-agnostic API for WebRTC operations:
 * - Android: Uses org.webrtc:google-webrtc library
 * - iOS: Uses native WebRTC.framework via Swift interop
 *
 * Quality Layers:
 * - Low: 360p @ 500 Kbps
 * - Mid: 720p @ 1500 Kbps
 * - High: 1080p @ 3000 Kbps
 */
expect class WebRTCClient() {

    /**
     * Connection state of the WebRTC peer connection
     */
    val connectionState: StateFlow<ConnectionState>

    /**
     * Current quality layer being received
     */
    val currentQuality: StateFlow<QualityLayer>

    /**
     * Real-time WebRTC statistics
     */
    val stats: StateFlow<WebRTCStats>

    /**
     * Connect to a live stream
     *
     * @param streamId The stream ID to connect to
     * @param sdpOffer The SDP offer from the broadcaster
     * @param onTrackReceived Callback when media track is received
     * @param onConnectionStateChange Callback for connection state changes
     * @param onStatsUpdate Callback for statistics updates (every 2s)
     */
    suspend fun connect(
        streamId: String,
        sdpOffer: String,
        onTrackReceived: (MediaTrack) -> Unit,
        onConnectionStateChange: (ConnectionState) -> Unit,
        onStatsUpdate: (WebRTCStats) -> Unit
    ): Result<String> // Returns SDP answer

    /**
     * Disconnect from the current stream
     */
    suspend fun disconnect()

    /**
     * Switch to a different quality layer
     *
     * @param quality Target quality layer
     */
    suspend fun switchQuality(quality: QualityLayer): Result<Unit>

    /**
     * Get current WebRTC statistics
     *
     * @return Current statistics snapshot
     */
    suspend fun getStats(): WebRTCStats

    /**
     * Send signaling message (ICE candidate, etc.)
     *
     * @param message Signaling message to send
     */
    suspend fun sendSignalingMessage(message: String): Result<Unit>
}

/**
 * WebRTC connection states
 */
enum class ConnectionState {
    NEW,           // Initial state
    CONNECTING,    // Attempting connection
    CONNECTED,     // Successfully connected
    RECONNECTING,  // Connection lost, attempting reconnection
    FAILED,        // Connection failed permanently
    CLOSED         // Connection closed by user
}

/**
 * Quality layers for adaptive bitrate streaming
 */
enum class QualityLayer(
    val resolution: String,
    val bitrate: Int,
    val label: String
) {
    LOW("360p", 500_000, "Low (360p)"),      // 500 Kbps
    MID("720p", 1_500_000, "Mid (720p)"),    // 1.5 Mbps
    HIGH("1080p", 3_000_000, "High (1080p)") // 3 Mbps
}

/**
 * Media track information
 */
data class MediaTrack(
    val id: String,
    val kind: String, // "video" or "audio"
    val enabled: Boolean
)

/**
 * WebRTC statistics for monitoring and quality adaptation
 */
data class WebRTCStats(
    val timestamp: Long = 0L, // Platform-specific timestamp set by actual implementation

    // Bandwidth metrics
    val currentBandwidth: Int = 0,        // Current bandwidth in Kbps
    val estimatedBandwidth: Int = 0,      // Estimated available bandwidth in Kbps

    // Packet metrics
    val packetsReceived: Long = 0,
    val packetsLost: Long = 0,
    val packetLossRate: Float = 0f,       // Percentage (0-100)

    // Frame metrics
    val framesReceived: Long = 0,
    val framesDropped: Long = 0,
    val frameRate: Int = 0,                // Current FPS

    // Quality metrics
    val currentQuality: QualityLayer = QualityLayer.MID,
    val jitter: Long = 0,                  // Jitter in milliseconds
    val roundTripTime: Long = 0            // RTT in milliseconds
) {
    /**
     * Calculate quality score (0-100)
     * Based on packet loss, frame rate, and jitter
     */
    fun calculateQualityScore(): Int {
        val lossScore = (100 - (packetLossRate * 2)).coerceIn(0f, 100f)
        val fpsScore = (frameRate / 30f * 100).coerceIn(0f, 100f)
        val jitterScore = ((200 - jitter) / 2f).coerceIn(0f, 100f)

        return ((lossScore + fpsScore + jitterScore) / 3).toInt()
    }

    /**
     * Suggest optimal quality layer based on current bandwidth
     */
    fun suggestQualityLayer(): QualityLayer {
        return when {
            estimatedBandwidth >= 2_500 -> QualityLayer.HIGH  // >2.5 Mbps
            estimatedBandwidth >= 1_200 -> QualityLayer.MID   // >1.2 Mbps
            else -> QualityLayer.LOW                           // <1.2 Mbps
        }
    }
}