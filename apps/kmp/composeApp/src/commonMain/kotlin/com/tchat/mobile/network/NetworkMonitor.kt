package com.tchat.mobile.network

import com.tchat.mobile.models.NetworkState
import com.tchat.mobile.models.ConnectionState
import kotlinx.coroutines.flow.*
import kotlinx.coroutines.delay
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.launch
import kotlin.time.Duration.Companion.seconds

/**
 * Network connectivity monitoring and state management
 *
 * Provides:
 * - Real-time network connectivity status
 * - Connection quality assessment
 * - Automatic reconnection logic
 * - Network state persistence
 */
class NetworkMonitor(
    private val scope: CoroutineScope
) {

    private val _networkState = MutableStateFlow(NetworkState.CONNECTED)
    val networkState: StateFlow<NetworkState> = _networkState.asStateFlow()

    private val _connectionQuality = MutableStateFlow(ConnectionQuality.GOOD)
    val connectionQuality: StateFlow<ConnectionQuality> = _connectionQuality.asStateFlow()

    private val _isOnline = MutableStateFlow(true)
    val isOnline: StateFlow<Boolean> = _isOnline.asStateFlow()

    private val _latency = MutableStateFlow(50L) // milliseconds
    val latency: StateFlow<Long> = _latency.asStateFlow()

    private var isMonitoring = false

    init {
        startMonitoring()
    }

    /**
     * Start network monitoring
     */
    fun startMonitoring() {
        if (isMonitoring) return
        isMonitoring = true

        scope.launch {
            while (isMonitoring) {
                checkNetworkConnectivity()
                delay(5.seconds) // Check every 5 seconds
            }
        }
    }

    /**
     * Stop network monitoring
     */
    fun stopMonitoring() {
        isMonitoring = false
    }

    /**
     * Force a network connectivity check
     */
    suspend fun checkConnectivity(): NetworkState {
        checkNetworkConnectivity()
        return _networkState.value
    }

    /**
     * Get current network information
     */
    fun getNetworkInfo(): NetworkInfo {
        return NetworkInfo(
            state = _networkState.value,
            quality = _connectionQuality.value,
            latency = _latency.value,
            isOnline = _isOnline.value
        )
    }

    private suspend fun checkNetworkConnectivity() {
        try {
            // Simulate network check - replace with actual platform-specific implementation
            val isConnected = simulateConnectivityCheck()
            val currentLatency = simulateLatencyCheck()

            _latency.value = currentLatency
            _isOnline.value = isConnected

            when {
                !isConnected -> {
                    _networkState.value = NetworkState.DISCONNECTED
                    _connectionQuality.value = ConnectionQuality.POOR
                }
                currentLatency > 1000 -> {
                    _networkState.value = NetworkState.LIMITED_CONNECTIVITY
                    _connectionQuality.value = ConnectionQuality.POOR
                }
                currentLatency > 500 -> {
                    _networkState.value = NetworkState.CONNECTED
                    _connectionQuality.value = ConnectionQuality.FAIR
                }
                else -> {
                    _networkState.value = NetworkState.CONNECTED
                    _connectionQuality.value = ConnectionQuality.GOOD
                }
            }
        } catch (e: Exception) {
            _networkState.value = NetworkState.DISCONNECTED
            _connectionQuality.value = ConnectionQuality.POOR
            _isOnline.value = false
        }
    }

    private suspend fun simulateConnectivityCheck(): Boolean {
        // Simulate network check - replace with actual implementation
        // For Android: use ConnectivityManager
        // For iOS: use NWPathMonitor
        delay(100) // Simulate network request
        return true // Assume connected for simulation
    }

    private suspend fun simulateLatencyCheck(): Long {
        // Simulate ping test - replace with actual implementation
        delay(50) // Simulate ping
        return (30..200).random().toLong() // Random latency between 30-200ms
    }
}

/**
 * Connection quality levels
 */
enum class ConnectionQuality {
    EXCELLENT, // < 50ms latency
    GOOD,      // 50-200ms latency
    FAIR,      // 200-500ms latency
    POOR       // > 500ms latency or unstable
}

/**
 * Network information data class
 */
data class NetworkInfo(
    val state: NetworkState,
    val quality: ConnectionQuality,
    val latency: Long,
    val isOnline: Boolean
)

/**
 * Network state manager for coordinating between local and remote data sources
 */
class NetworkStateManager(
    private val networkMonitor: NetworkMonitor,
    private val scope: CoroutineScope
) {

    private val _shouldSyncOfflineChanges = MutableStateFlow(false)
    val shouldSyncOfflineChanges: StateFlow<Boolean> = _shouldSyncOfflineChanges.asStateFlow()

    private val _isInOfflineMode = MutableStateFlow(false)
    val isInOfflineMode: StateFlow<Boolean> = _isInOfflineMode.asStateFlow()

    private var wasOffline = false

    init {
        observeNetworkChanges()
    }

    private fun observeNetworkChanges() {
        scope.launch {
            networkMonitor.networkState.collect { networkState ->
                val isCurrentlyOffline = networkState == NetworkState.DISCONNECTED

                _isInOfflineMode.value = isCurrentlyOffline

                // Detect transition from offline to online
                if (wasOffline && !isCurrentlyOffline) {
                    _shouldSyncOfflineChanges.value = true
                }

                wasOffline = isCurrentlyOffline
            }
        }
    }

    /**
     * Mark offline changes as synced
     */
    fun markOfflineChangesSynced() {
        _shouldSyncOfflineChanges.value = false
    }

    /**
     * Check if we should use cache-first strategy
     */
    fun shouldUseCacheFirst(): Boolean {
        val networkInfo = networkMonitor.getNetworkInfo()
        return !networkInfo.isOnline || networkInfo.quality == ConnectionQuality.POOR
    }

    /**
     * Check if we should attempt network requests
     */
    fun shouldAttemptNetworkRequest(): Boolean {
        val networkInfo = networkMonitor.getNetworkInfo()
        return networkInfo.isOnline && networkInfo.quality != ConnectionQuality.POOR
    }

    /**
     * Get current network information
     */
    fun getNetworkInfo(): NetworkInfo {
        return networkMonitor.getNetworkInfo()
    }

    /**
     * Get recommended sync strategy based on network conditions
     */
    fun getRecommendedSyncStrategy(): SyncStrategy {
        val networkInfo = networkMonitor.getNetworkInfo()

        return when {
            !networkInfo.isOnline -> SyncStrategy.CACHE_ONLY
            networkInfo.quality == ConnectionQuality.POOR -> SyncStrategy.CACHE_FIRST
            networkInfo.quality in listOf(ConnectionQuality.FAIR, ConnectionQuality.GOOD) -> SyncStrategy.NETWORK_FIRST
            else -> SyncStrategy.NETWORK_ONLY
        }
    }
}

/**
 * Sync strategies based on network conditions
 */
enum class SyncStrategy {
    CACHE_ONLY,     // No network requests, use cached data only
    CACHE_FIRST,    // Try cache first, fallback to network
    NETWORK_FIRST,  // Try network first, fallback to cache
    NETWORK_ONLY    // Network requests only, no cache fallback
}

/**
 * Connection retry manager with exponential backoff
 */
class ConnectionRetryManager(
    private val scope: CoroutineScope
) {
    private var currentRetryCount = 0
    private val maxRetries = 5
    private val baseDelayMs = 1000L

    /**
     * Execute operation with retry logic
     */
    suspend fun <T> executeWithRetry(
        operation: suspend () -> Result<T>
    ): Result<T> {
        repeat(maxRetries) { attempt ->
            val result = operation()

            if (result.isSuccess) {
                currentRetryCount = 0 // Reset on success
                return result
            }

            if (attempt < maxRetries - 1) {
                val delayMs = baseDelayMs * (1L shl attempt) // Exponential backoff
                delay(delayMs)
            }
        }

        currentRetryCount++
        return operation() // Final attempt
    }

    /**
     * Reset retry counter
     */
    fun reset() {
        currentRetryCount = 0
    }

    /**
     * Check if we should continue retrying
     */
    fun shouldRetry(): Boolean {
        return currentRetryCount < maxRetries
    }
}

/**
 * Platform-specific network utilities
 * Implement these in androidMain and iosMain source sets
 */
expect class PlatformNetworkMonitor() {
    fun isNetworkAvailable(): Boolean
    suspend fun pingServer(host: String): Long
    fun getNetworkType(): String
    fun getSignalStrength(): Int
}