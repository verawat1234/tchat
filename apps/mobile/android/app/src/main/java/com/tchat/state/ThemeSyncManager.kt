package com.tchat.state

import android.content.Context
import kotlinx.coroutines.*
import kotlinx.coroutines.flow.*
import kotlinx.serialization.Serializable
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import kotlinx.serialization.Contextual
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody
import java.util.*
import java.util.concurrent.TimeUnit

/**
 * Manages theme synchronization between web and mobile platforms
 */
class ThemeSyncManager(
    private val context: Context,
    private val stateSyncManager: StateSyncManager,
    private val persistenceManager: PersistenceManager
) {

    // MARK: - Properties
    private val baseUrl = "https://api.tchat.app"
    private val client = OkHttpClient.Builder()
        .connectTimeout(30, TimeUnit.SECONDS)
        .readTimeout(30, TimeUnit.SECONDS)
        .writeTimeout(30, TimeUnit.SECONDS)
        .build()

    private val json = Json {
        ignoreUnknownKeys = true
        coerceInputValues = true
    }

    private val _isThemeSyncing = MutableStateFlow(false)
    val isThemeSyncing: StateFlow<Boolean> = _isThemeSyncing.asStateFlow()

    private val _lastThemeSync = MutableStateFlow<Date?>(null)
    val lastThemeSync: StateFlow<Date?> = _lastThemeSync.asStateFlow()

    private var autoSyncJob: Job? = null
    private val scope = CoroutineScope(Dispatchers.Main + SupervisorJob())

    // MARK: - Public Methods

    /**
     * Sync theme changes to server immediately
     */
    suspend fun syncThemeToServer(themePreferences: ThemePreferences) {
        _isThemeSyncing.value = true

        try {
            val themePayload = ThemeSyncPayload(
                timestamp = Date(),
                platform = "android",
                preferences = themePreferences,
                syncReason = ThemeSyncReason.USER_CHANGED
            )

            uploadThemeChanges(themePayload)
            _lastThemeSync.value = Date()

        } finally {
            _isThemeSyncing.value = false
        }
    }

    /**
     * Poll for theme changes from other platforms
     */
    suspend fun pollThemeChanges(): ThemePreferences? {
        val serverTheme = downloadLatestTheme()

        // Check if server theme is newer than our last sync
        val serverTimestamp = serverTheme.timestamp
        val lastSync = _lastThemeSync.value

        return if (serverTimestamp != null && lastSync != null && serverTimestamp > lastSync) {
            serverTheme.preferences
        } else {
            null
        }
    }

    /**
     * Apply theme changes with cross-platform compatibility
     */
    suspend fun applyThemeChanges(preferences: ThemePreferences, fromPlatform: String) {
        withContext(Dispatchers.Main) {
            // Store theme locally
            persistenceManager.saveThemePreferences(preferences)

            // Notify app components of theme change
            // This would typically be done through a shared state or event bus
            notifyThemeChanged(preferences, fromPlatform)

            _lastThemeSync.value = Date()
        }
    }

    /**
     * Start automatic theme synchronization
     */
    fun startAutoSync(intervalSeconds: Long = 30L) {
        autoSyncJob?.cancel()
        autoSyncJob = scope.launch {
            while (isActive) {
                try {
                    checkForThemeUpdates()
                    delay(intervalSeconds * 1000)
                } catch (e: Exception) {
                    // Log error but continue syncing
                    delay(intervalSeconds * 1000)
                }
            }
        }
    }

    /**
     * Stop automatic theme synchronization
     */
    fun stopAutoSync() {
        autoSyncJob?.cancel()
        autoSyncJob = null
    }

    /**
     * Cleanup resources
     */
    fun cleanup() {
        stopAutoSync()
        scope.cancel()
    }

    // MARK: - Private Methods

    private suspend fun checkForThemeUpdates() {
        if (_isThemeSyncing.value) return

        val newTheme = pollThemeChanges()
        if (newTheme != null) {
            applyThemeChanges(newTheme, "server")
        }
    }

    private suspend fun uploadThemeChanges(payload: ThemeSyncPayload) = withContext(Dispatchers.IO) {
        val url = "$baseUrl/api/theme/sync"
        val jsonString = json.encodeToString(payload)
        val mediaType = "application/json; charset=utf-8".toMediaType()
        val requestBody = jsonString.toRequestBody(mediaType)

        val requestBuilder = Request.Builder()
            .url(url)
            .post(requestBody)
            .addHeader("Content-Type", "application/json")

        // Add authentication header if needed
        getAuthToken()?.let { token ->
            requestBuilder.addHeader("Authorization", "Bearer $token")
        }

        val request = requestBuilder.build()
        val response = client.newCall(request).execute()

        if (!response.isSuccessful) {
            throw ThemeSyncError.UploadFailed
        }
    }

    private suspend fun downloadLatestTheme(): ThemeServerResponse = withContext(Dispatchers.IO) {
        val url = "$baseUrl/api/theme/latest"

        val requestBuilder = Request.Builder()
            .url(url)
            .get()
            .addHeader("Content-Type", "application/json")

        // Add authentication header if needed
        getAuthToken()?.let { token ->
            requestBuilder.addHeader("Authorization", "Bearer $token")
        }

        val request = requestBuilder.build()
        val response = client.newCall(request).execute()

        if (!response.isSuccessful) {
            throw ThemeSyncError.DownloadFailed
        }

        val responseBody = response.body?.string() ?: throw ThemeSyncError.DownloadFailed
        json.decodeFromString<ThemeServerResponse>(responseBody)
    }

    private fun notifyThemeChanged(preferences: ThemePreferences, source: String) {
        // In a real app, this would use an event bus or shared state management
        // For now, we'll use a simple callback approach
        // You could implement this with Android's LocalBroadcastManager or a similar mechanism
    }

    private fun getAuthToken(): String? {
        // Retrieve authentication token from secure storage
        return null // Placeholder
    }
}

// MARK: - Theme Sync Models

@Serializable
data class ThemeSyncPayload(
    val timestamp: @Contextual Date,
    val platform: String,
    val preferences: @Contextual ThemePreferences,
    val syncReason: ThemeSyncReason
)

@Serializable
data class ThemeServerResponse(
    val timestamp: @Contextual Date? = null,
    val preferences: @Contextual ThemePreferences,
    val lastModifiedBy: String? = null
)

enum class ThemeSyncReason {
    USER_CHANGED,
    SYSTEM_CHANGED,
    CROSS_PLATFORM_SYNC,
    STARTUP
}

/**
 * Theme sync errors
 */
sealed class ThemeSyncError(message: String) : Exception(message) {
    object InvalidUrl : ThemeSyncError("Invalid theme sync URL")
    object UploadFailed : ThemeSyncError("Failed to upload theme changes")
    object DownloadFailed : ThemeSyncError("Failed to download theme updates")
    object NetworkUnavailable : ThemeSyncError("Network unavailable for theme sync")
}