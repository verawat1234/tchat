package com.tchat.state

import android.content.Context
import android.provider.Settings
import com.tchat.models.User
import com.tchat.models.Workspace
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
 * Manages cross-platform session state and workspace switching
 */
class SessionManager(
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

    private val _currentSession = MutableStateFlow<UserSession?>(null)
    val currentSession: StateFlow<UserSession?> = _currentSession.asStateFlow()

    private val _activeWorkspace = MutableStateFlow<Workspace?>(null)
    val activeWorkspace: StateFlow<Workspace?> = _activeWorkspace.asStateFlow()

    private val _isSessionSyncing = MutableStateFlow(false)
    val isSessionSyncing: StateFlow<Boolean> = _isSessionSyncing.asStateFlow()

    private val _lastSessionSync = MutableStateFlow<Date?>(null)
    val lastSessionSync: StateFlow<Date?> = _lastSessionSync.asStateFlow()

    // Session timeout configuration
    private val sessionTimeout = 30 * 60 * 1000L // 30 minutes in milliseconds
    private var sessionMonitoringJob: Job? = null
    private var lastSyncTime = Date(0)

    private val scope = CoroutineScope(Dispatchers.Main + SupervisorJob())

    init {
        loadPersistedSession()
    }

    // MARK: - Session Management

    /**
     * Start a new session with user authentication
     */
    suspend fun startSession(user: User, workspace: Workspace) {
        val session = UserSession(
            id = UUID.randomUUID().toString(),
            userId = user.id,
            workspaceId = workspace.id,
            platform = "android",
            deviceId = getDeviceId(),
            startTime = Date(),
            lastActivity = Date(),
            isActive = true
        )

        _currentSession.value = session
        _activeWorkspace.value = workspace

        // Sync session to server
        syncSessionToServer(session)

        // Start session monitoring
        startSessionMonitoring()

        // Persist session locally
        persistenceManager.saveSession(session)
        persistenceManager.saveWorkspace(workspace)

        // Notify session start
        notifySessionStart(session)
    }

    /**
     * Switch to a different workspace
     */
    suspend fun switchWorkspace(workspace: Workspace) {
        val session = _currentSession.value
            ?: throw SessionError.NoActiveSession

        // Update session with new workspace
        val updatedSession = session.copy(
            workspaceId = workspace.id,
            lastActivity = Date()
        )

        _currentSession.value = updatedSession
        _activeWorkspace.value = workspace

        // Sync workspace change to server
        syncWorkspaceSwitch(updatedSession, workspace)

        // Persist changes
        persistenceManager.saveSession(updatedSession)
        persistenceManager.saveWorkspace(workspace)

        // Notify workspace change
        notifyWorkspaceChange(workspace)
    }

    /**
     * Update session activity
     */
    suspend fun updateActivity() {
        val session = _currentSession.value ?: return

        val updatedSession = session.copy(
            lastActivity = Date()
        )

        _currentSession.value = updatedSession
        persistenceManager.saveSession(updatedSession)

        // Sync activity update to server (throttled)
        throttledSessionSync(updatedSession)
    }

    /**
     * End current session
     */
    suspend fun endSession() {
        val session = _currentSession.value
            ?: throw SessionError.NoActiveSession

        val endedSession = session.copy(
            lastActivity = Date(),
            isActive = false
        )

        // Sync session end to server
        syncSessionToServer(endedSession)

        // Clear local session
        _currentSession.value = null
        _activeWorkspace.value = null
        stopSessionMonitoring()

        // Clear persisted session
        persistenceManager.clearSession()

        // Notify session end
        notifySessionEnd(endedSession)
    }

    /**
     * Get active sessions across all platforms
     */
    suspend fun getActiveSessions(): List<CrossPlatformSession> = withContext(Dispatchers.IO) {
        val url = "$baseUrl/api/sessions/active"

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
            throw SessionError.ServerError
        }

        val responseBody = response.body?.string() ?: throw SessionError.ServerError
        json.decodeFromString<List<CrossPlatformSession>>(responseBody)
    }

    /**
     * Cleanup resources
     */
    fun cleanup() {
        stopSessionMonitoring()
        scope.cancel()
    }

    // MARK: - Private Methods

    private fun loadPersistedSession() {
        val session = persistenceManager.loadSession()
        val workspace = persistenceManager.loadWorkspace()

        if (session != null && workspace != null) {
            // Check if session is still valid (not expired)
            val timeSinceLastActivity = Date().time - session.lastActivity.time
            if (timeSinceLastActivity < sessionTimeout) {
                _currentSession.value = session
                _activeWorkspace.value = workspace
                startSessionMonitoring()
            } else {
                // Session expired, clear it
                persistenceManager.clearSession()
            }
        }
    }

    private fun startSessionMonitoring() {
        sessionMonitoringJob?.cancel()
        sessionMonitoringJob = scope.launch {
            while (isActive) {
                checkSessionExpiry()
                delay(60_000) // Check every minute
            }
        }
    }

    private fun stopSessionMonitoring() {
        sessionMonitoringJob?.cancel()
        sessionMonitoringJob = null
    }

    private suspend fun checkSessionExpiry() {
        val session = _currentSession.value ?: return

        val timeSinceLastActivity = Date().time - session.lastActivity.time
        if (timeSinceLastActivity >= sessionTimeout) {
            try {
                endSession()
            } catch (e: Exception) {
                // Handle session end error
            }
        }
    }

    private suspend fun syncSessionToServer(session: UserSession) = withContext(Dispatchers.IO) {
        _isSessionSyncing.value = true

        try {
            val url = "$baseUrl/api/sessions/sync"
            val jsonString = json.encodeToString(session)
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
                throw SessionError.SyncFailed
            }

            _lastSessionSync.value = Date()

        } finally {
            _isSessionSyncing.value = false
        }
    }

    private suspend fun syncWorkspaceSwitch(session: UserSession, workspace: Workspace) = withContext(Dispatchers.IO) {
        val payload = WorkspaceSwitchPayload(
            sessionId = session.id,
            userId = session.userId,
            newWorkspaceId = workspace.id,
            timestamp = Date(),
            platform = "android"
        )

        val url = "$baseUrl/api/sessions/workspace-switch"
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
            throw SessionError.SyncFailed
        }
    }

    private suspend fun throttledSessionSync(session: UserSession) {
        val now = Date()
        val timeSinceLastSync = now.time - lastSyncTime.time

        // Only sync if more than 30 seconds have passed
        if (timeSinceLastSync >= 30_000) {
            syncSessionToServer(session)
            lastSyncTime = now
        }
    }

    private fun getDeviceId(): String {
        return Settings.Secure.getString(
            context.contentResolver,
            Settings.Secure.ANDROID_ID
        ) ?: UUID.randomUUID().toString()
    }

    private fun getAuthToken(): String? {
        // Retrieve authentication token from secure storage
        return null // Placeholder
    }

    private fun notifySessionStart(session: UserSession) {
        // In a real app, this would use an event bus or broadcast
        // For now, this is a placeholder for session start notifications
    }

    private fun notifySessionEnd(session: UserSession) {
        // In a real app, this would use an event bus or broadcast
        // For now, this is a placeholder for session end notifications
    }

    private fun notifyWorkspaceChange(workspace: Workspace) {
        // In a real app, this would use an event bus or broadcast
        // For now, this is a placeholder for workspace change notifications
    }
}

// MARK: - Session Models

@Serializable
data class UserSession(
    val id: String,
    val userId: String,
    val workspaceId: String,
    val platform: String,
    val deviceId: String,
    val startTime: @Contextual Date,
    val lastActivity: @Contextual Date,
    val isActive: Boolean
)


@Serializable
data class CrossPlatformSession(
    val sessionId: String,
    val platform: String,
    val deviceType: String,
    val lastActivity: @Contextual Date,
    val isCurrentSession: Boolean
)

@Serializable
data class WorkspaceSwitchPayload(
    val sessionId: String,
    val userId: String,
    val newWorkspaceId: String,
    val timestamp: @Contextual Date,
    val platform: String
)

/**
 * Session errors
 */
sealed class SessionError(message: String) : Exception(message) {
    object NoActiveSession : SessionError("No active session found")
    object InvalidUrl : SessionError("Invalid session management URL")
    object ServerError : SessionError("Session server error")
    object SyncFailed : SessionError("Failed to sync session")
    object SessionExpired : SessionError("Session has expired")
}