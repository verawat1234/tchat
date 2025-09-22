package com.tchat.state

import android.content.Context
import android.net.ConnectivityManager
import android.net.Network
import android.net.NetworkCapabilities
import android.net.NetworkRequest
import androidx.core.content.ContextCompat.getSystemService
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import kotlinx.serialization.Serializable
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import kotlinx.serialization.Contextual
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody
import java.util.Date
import java.util.concurrent.TimeUnit

/**
 * Manages state synchronization between native app and web platform
 */
class StateSyncManager(private val context: Context? = null) {

    // MARK: - Properties
    private val baseUrl = "https://api.tchat.app" // Replace with actual API endpoint
    private val client = OkHttpClient.Builder()
        .connectTimeout(30, TimeUnit.SECONDS)
        .readTimeout(30, TimeUnit.SECONDS)
        .writeTimeout(30, TimeUnit.SECONDS)
        .build()

    private val json = Json {
        ignoreUnknownKeys = true
        coerceInputValues = true
    }

    private var isConnected: Boolean = false
    var lastSyncTimestamp: Date? = null

    init {
        context?.let { startNetworkMonitoring(it) }
    }

    // MARK: - Public Methods

    /**
     * Sync app state with server
     */
    suspend fun syncState(appState: AppState) {
        if (!isConnected) {
            throw SyncError.NetworkUnavailable
        }

        withContext(Dispatchers.IO) {
            try {
                val statePayload = createStatePayload(appState)

                // Upload current state
                uploadState(statePayload)

                // Download latest state
                val serverState = downloadState()
                updateAppState(appState, serverState)

                lastSyncTimestamp = Date()

            } catch (e: Exception) {
                throw SyncError.SyncFailed(e.message ?: "Unknown error")
            }
        }
    }

    /**
     * Download state from server
     */
    private suspend fun downloadState(): ServerState = withContext(Dispatchers.IO) {
        val url = "$baseUrl/api/state/sync"

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
            throw SyncError.ServerError
        }

        val responseBody = response.body?.string() ?: throw SyncError.ServerError
        json.decodeFromString<ServerState>(responseBody)
    }

    /**
     * Upload state to server
     */
    private suspend fun uploadState(state: StatePayload) = withContext(Dispatchers.IO) {
        val url = "$baseUrl/api/state/sync"
        val jsonString = json.encodeToString(state)
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
            throw SyncError.ServerError
        }
    }

    // MARK: - Private Methods

    private fun startNetworkMonitoring(context: Context) {
        val connectivityManager = getSystemService(context, ConnectivityManager::class.java)
        val networkRequest = NetworkRequest.Builder()
            .addCapability(NetworkCapabilities.NET_CAPABILITY_INTERNET)
            .build()

        val networkCallback = object : ConnectivityManager.NetworkCallback() {
            override fun onAvailable(network: Network) {
                isConnected = true
            }

            override fun onLost(network: Network) {
                isConnected = false
            }
        }

        connectivityManager?.registerNetworkCallback(networkRequest, networkCallback)

        // Check initial connectivity
        val activeNetwork = connectivityManager?.activeNetwork
        val networkCapabilities = connectivityManager?.getNetworkCapabilities(activeNetwork)
        isConnected = networkCapabilities?.hasCapability(NetworkCapabilities.NET_CAPABILITY_INTERNET) == true
    }

    private suspend fun createStatePayload(appState: AppState): StatePayload {
        return StatePayload(
            timestamp = Date(),
            platform = "android",
            version = "1.0.0", // Get from BuildConfig
            userId = appState.currentUser.value?.id,
            themePreferences = appState.themePreferences.value,
            chatState = appState.chatState.value,
            storeState = appState.storeState.value,
            socialState = appState.socialState.value,
            videoState = appState.videoState.value
        )
    }

    private suspend fun updateAppState(appState: AppState, serverState: ServerState) {
        withContext(Dispatchers.Main) {
            // Only update if server state is newer
            val serverTimestamp = serverState.timestamp
            val lastSync = lastSyncTimestamp

            if (serverTimestamp != null && lastSync != null && serverTimestamp <= lastSync) {
                return@withContext
            }

            // Update theme preferences if different
            if (serverState.themePreferences != appState.themePreferences.value) {
                appState.updateThemePreferences(serverState.themePreferences)
            }

            // Update chat state
            if (serverState.chatState.unreadCount != appState.chatState.value.unreadCount) {
                appState.updateChatState(serverState.chatState)
            }

            // Update store state
            if (serverState.storeState.cartItemCount != appState.storeState.value.cartItemCount) {
                appState.updateStoreState(serverState.storeState)
            }

            // Update social state
            if (serverState.socialState.friendRequestCount != appState.socialState.value.friendRequestCount) {
                appState.updateSocialState(serverState.socialState)
            }

            // Update video state
            if (serverState.videoState.watchHistoryCount != appState.videoState.value.watchHistoryCount) {
                appState.updateVideoState(serverState.videoState)
            }
        }
    }

    private fun getAuthToken(): String? {
        // Retrieve authentication token from secure storage
        return null // Placeholder
    }
}

// MARK: - State Models

/**
 * Payload sent to server
 */
@Serializable
data class StatePayload(
    val timestamp: @Contextual Date,
    val platform: String,
    val version: String,
    val userId: String? = null,
    val themePreferences: @Contextual ThemePreferences,
    val chatState: @Contextual ChatState,
    val storeState: @Contextual StoreState,
    val socialState: @Contextual SocialState,
    val videoState: @Contextual VideoState
)

/**
 * State received from server
 */
@Serializable
data class ServerState(
    val timestamp: @Contextual Date? = null,
    val themePreferences: @Contextual ThemePreferences,
    val chatState: @Contextual ChatState,
    val storeState: @Contextual StoreState,
    val socialState: @Contextual SocialState,
    val videoState: @Contextual VideoState
)

/**
 * Sync errors
 */
sealed class SyncError(message: String) : Exception(message) {
    object NetworkUnavailable : SyncError("Network is unavailable")
    object InvalidUrl : SyncError("Invalid server URL")
    object ServerError : SyncError("Server error occurred")
    class SyncFailed(message: String) : SyncError("Sync failed: $message")
    object DataCorruption : SyncError("Data corruption detected")
}