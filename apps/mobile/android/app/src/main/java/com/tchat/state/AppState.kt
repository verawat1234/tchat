package com.tchat.state

import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.getValue
import androidx.compose.runtime.setValue
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.launch
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.delay
import java.util.Date

/**
 * Global application state management following Tchat design system
 */
class AppState : ViewModel() {

    // MARK: - State Properties

    private val _currentUser = MutableStateFlow<UserModel?>(null)
    val currentUser: StateFlow<UserModel?> = _currentUser.asStateFlow()

    private val _isAuthenticated = MutableStateFlow(false)
    val isAuthenticated: StateFlow<Boolean> = _isAuthenticated.asStateFlow()

    private val _themePreferences = MutableStateFlow(ThemePreferences())
    val themePreferences: StateFlow<ThemePreferences> = _themePreferences.asStateFlow()

    private val _chatState = MutableStateFlow(ChatState())
    val chatState: StateFlow<ChatState> = _chatState.asStateFlow()

    private val _storeState = MutableStateFlow(StoreState())
    val storeState: StateFlow<StoreState> = _storeState.asStateFlow()

    private val _socialState = MutableStateFlow(SocialState())
    val socialState: StateFlow<SocialState> = _socialState.asStateFlow()

    private val _videoState = MutableStateFlow(VideoState())
    val videoState: StateFlow<VideoState> = _videoState.asStateFlow()

    private val _isOnline = MutableStateFlow(true)
    val isOnline: StateFlow<Boolean> = _isOnline.asStateFlow()

    private val _syncStatus = MutableStateFlow(SyncStatus.Idle)
    val syncStatus: StateFlow<SyncStatus> = _syncStatus.asStateFlow()

    private val _lastError = MutableStateFlow<AppError?>(null)
    val lastError: StateFlow<AppError?> = _lastError.asStateFlow()

    // MARK: - Dependencies
    private val syncManager = StateSyncManager()
    private val persistence = PersistenceManager()

    init {
        loadPersistedState()
        startSyncTimer()
    }

    // MARK: - Public Methods

    /**
     * Update user information
     */
    fun updateUser(user: UserModel) {
        _currentUser.value = user
        _isAuthenticated.value = true
        syncWithServer()
    }

    /**
     * Sign out user
     */
    fun signOut() {
        _currentUser.value = null
        _isAuthenticated.value = false
        clearAllState()
        syncWithServer()
    }

    /**
     * Update theme preferences
     */
    fun updateThemePreferences(preferences: ThemePreferences) {
        _themePreferences.value = preferences
        persistence.saveThemePreferences(preferences)
        syncWithServer()
    }

    /**
     * Update chat state
     */
    fun updateChatState(state: ChatState) {
        _chatState.value = state
        persistence.saveChatState(state)
        syncWithServer()
    }

    /**
     * Update store state
     */
    fun updateStoreState(state: StoreState) {
        _storeState.value = state
        persistence.saveStoreState(state)
        syncWithServer()
    }

    /**
     * Update social state
     */
    fun updateSocialState(state: SocialState) {
        _socialState.value = state
        persistence.saveSocialState(state)
        syncWithServer()
    }

    /**
     * Update video state
     */
    fun updateVideoState(state: VideoState) {
        _videoState.value = state
        persistence.saveVideoState(state)
        syncWithServer()
    }

    /**
     * Update network connectivity
     */
    fun updateNetworkConnectivity(isConnected: Boolean) {
        _isOnline.value = isConnected
        if (isConnected) {
            syncWithServer()
        }
    }

    /**
     * Force sync with server
     */
    fun forceSyncWithServer() {
        _syncStatus.value = SyncStatus.Syncing
        syncWithServer()
    }

    /**
     * Handle deep link
     */
    fun handleDeepLink(url: String) {
        // Deep link routing logic
        println("Handling deep link: $url")
    }

    /**
     * Clear error state
     */
    fun clearError() {
        _lastError.value = null
    }

    // MARK: - Private Methods

    private fun loadPersistedState() {
        // Load theme preferences
        persistence.loadThemePreferences()?.let { preferences ->
            _themePreferences.value = preferences
        }

        // Load authentication state
        _isAuthenticated.value = persistence.loadAuthenticationState()

        // Load user data
        persistence.loadCurrentUser()?.let { user ->
            _currentUser.value = user
        }

        // Load other states
        persistence.loadChatState()?.let { state ->
            _chatState.value = state
        }

        persistence.loadStoreState()?.let { state ->
            _storeState.value = state
        }

        persistence.loadSocialState()?.let { state ->
            _socialState.value = state
        }

        persistence.loadVideoState()?.let { state ->
            _videoState.value = state
        }
    }

    private fun saveState() {
        persistence.saveThemePreferences(_themePreferences.value)
        persistence.saveAuthenticationState(_isAuthenticated.value)

        _currentUser.value?.let { user ->
            persistence.saveCurrentUser(user)
        }

        persistence.saveChatState(_chatState.value)
        persistence.saveStoreState(_storeState.value)
        persistence.saveSocialState(_socialState.value)
        persistence.saveVideoState(_videoState.value)
    }

    private fun clearAllState() {
        _chatState.value = ChatState()
        _storeState.value = StoreState()
        _socialState.value = SocialState()
        _videoState.value = VideoState()
        persistence.clearAll()
    }

    private fun startSyncTimer() {
        viewModelScope.launch {
            while (true) {
                delay(30000) // 30 seconds
                if (_isAuthenticated.value && _isOnline.value) {
                    syncWithServer()
                }
            }
        }
    }

    private fun syncWithServer() {
        if (!_isOnline.value || !_isAuthenticated.value) return

        _syncStatus.value = SyncStatus.Syncing

        viewModelScope.launch {
            try {
                syncManager.syncState(this@AppState)
                _syncStatus.value = SyncStatus.Success
                saveState()
            } catch (e: Exception) {
                _syncStatus.value = SyncStatus.Failed
                _lastError.value = AppError.SyncFailed(e.message ?: "Unknown error")
            }
        }
    }

    override fun onCleared() {
        super.onCleared()
        saveState()
    }
}

// MARK: - State Models

/**
 * User model
 */
data class UserModel(
    val id: String,
    val username: String,
    val email: String,
    val displayName: String,
    val avatarUrl: String? = null,
    val preferences: UserPreferences = UserPreferences()
)

/**
 * User preferences
 */
data class UserPreferences(
    val language: String = "en",
    val notificationsEnabled: Boolean = true,
    val soundEnabled: Boolean = true
)

/**
 * Theme preferences
 */
data class ThemePreferences(
    val isDarkMode: Boolean = false,
    val accentColor: String = "#3B82F6",
    val fontSize: FontSize = FontSize.Medium
)

enum class FontSize(val value: String) {
    Small("small"),
    Medium("medium"),
    Large("large"),
    ExtraLarge("extraLarge")
}

/**
 * Chat state
 */
data class ChatState(
    val unreadCount: Int = 0,
    val lastMessageTimestamp: Date? = null,
    val activeConversations: List<String> = emptyList()
)

/**
 * Store state
 */
data class StoreState(
    val cartItemCount: Int = 0,
    val wishlistCount: Int = 0,
    val lastPurchaseDate: Date? = null
)

/**
 * Social state
 */
data class SocialState(
    val friendRequestCount: Int = 0,
    val notificationCount: Int = 0,
    val lastActivityDate: Date? = null
)

/**
 * Video state
 */
data class VideoState(
    val watchHistoryCount: Int = 0,
    val subscriptionCount: Int = 0,
    val lastWatchedDate: Date? = null
)

/**
 * Sync status
 */
enum class SyncStatus {
    Idle,
    Syncing,
    Success,
    Failed
}

/**
 * Application errors
 */
sealed class AppError(val message: String) {
    class SyncFailed(message: String) : AppError("Sync failed: $message")
    class AuthenticationFailed(message: String) : AppError("Authentication failed: $message")
    class NetworkError(message: String) : AppError("Network error: $message")
    class DataCorruption(message: String) : AppError("Data corruption: $message")
}