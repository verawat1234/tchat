package com.tchat.state

import android.content.Context
import android.content.SharedPreferences
import androidx.security.crypto.EncryptedSharedPreferences
import androidx.security.crypto.MasterKeys
import com.tchat.models.Workspace
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json

/**
 * Manages local data persistence for the Tchat app
 */
class PersistenceManager(private val context: Context? = null) {

    // MARK: - Properties
    private val sharedPreferences: SharedPreferences? by lazy {
        context?.getSharedPreferences("tchat_prefs", Context.MODE_PRIVATE)
    }

    private val encryptedPreferences: SharedPreferences? by lazy {
        context?.let { ctx ->
            val masterKeyAlias = MasterKeys.getOrCreate(MasterKeys.AES256_GCM_SPEC)
            EncryptedSharedPreferences.create(
                "tchat_secure_prefs",
                masterKeyAlias,
                ctx,
                EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
                EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
            )
        }
    }

    private val json = Json {
        ignoreUnknownKeys = true
        coerceInputValues = true
    }

    // MARK: - SharedPreferences Keys
    companion object {
        private const val KEY_THEME_PREFERENCES = "theme_preferences"
        private const val KEY_IS_AUTHENTICATED = "is_authenticated"
        private const val KEY_CURRENT_USER = "current_user"
        private const val KEY_CHAT_STATE = "chat_state"
        private const val KEY_STORE_STATE = "store_state"
        private const val KEY_SOCIAL_STATE = "social_state"
        private const val KEY_VIDEO_STATE = "video_state"
        private const val KEY_LAST_SYNC_TIMESTAMP = "last_sync_timestamp"
    }

    // MARK: - Generic Methods

    /**
     * Save object as JSON string
     */
    private inline fun <reified T> saveObject(obj: T, key: String, secure: Boolean = false) {
        val jsonString = json.encodeToString(obj)
        val prefs = if (secure) encryptedPreferences else sharedPreferences
        prefs?.edit()?.putString(key, jsonString)?.apply()
    }

    /**
     * Load object from JSON string
     */
    private inline fun <reified T> loadObject(key: String, secure: Boolean = false): T? {
        val prefs = if (secure) encryptedPreferences else sharedPreferences
        val jsonString = prefs?.getString(key, null) ?: return null
        return try {
            json.decodeFromString<T>(jsonString)
        } catch (e: Exception) {
            null
        }
    }

    /**
     * Save string value
     */
    fun saveString(key: String, value: String, secure: Boolean = false) {
        val prefs = if (secure) encryptedPreferences else sharedPreferences
        prefs?.edit()?.putString(key, value)?.apply()
    }

    /**
     * Load string value
     */
    fun loadString(key: String, secure: Boolean = false): String? {
        val prefs = if (secure) encryptedPreferences else sharedPreferences
        return prefs?.getString(key, null)
    }

    /**
     * Save boolean value
     */
    fun saveBoolean(key: String, value: Boolean) {
        sharedPreferences?.edit()?.putBoolean(key, value)?.apply()
    }

    /**
     * Load boolean value
     */
    fun loadBoolean(key: String, defaultValue: Boolean = false): Boolean {
        return sharedPreferences?.getBoolean(key, defaultValue) ?: defaultValue
    }

    /**
     * Save integer value
     */
    fun saveInt(key: String, value: Int) {
        sharedPreferences?.edit()?.putInt(key, value)?.apply()
    }

    /**
     * Load integer value
     */
    fun loadInt(key: String, defaultValue: Int = 0): Int {
        return sharedPreferences?.getInt(key, defaultValue) ?: defaultValue
    }

    /**
     * Save long value
     */
    fun saveLong(key: String, value: Long) {
        sharedPreferences?.edit()?.putLong(key, value)?.apply()
    }

    /**
     * Load long value
     */
    fun loadLong(key: String, defaultValue: Long = 0L): Long {
        return sharedPreferences?.getLong(key, defaultValue) ?: defaultValue
    }

    /**
     * Remove value for key
     */
    fun remove(key: String, secure: Boolean = false) {
        val prefs = if (secure) encryptedPreferences else sharedPreferences
        prefs?.edit()?.remove(key)?.apply()
    }

    /**
     * Clear all stored data
     */
    fun clearAll() {
        sharedPreferences?.edit()?.clear()?.apply()
        encryptedPreferences?.edit()?.clear()?.apply()
    }

    // MARK: - Specific State Methods

    /**
     * Save theme preferences
     */
    fun saveThemePreferences(preferences: ThemePreferences) {
        saveObject(preferences, KEY_THEME_PREFERENCES)
    }

    /**
     * Load theme preferences
     */
    fun loadThemePreferences(): ThemePreferences? {
        return loadObject<ThemePreferences>(KEY_THEME_PREFERENCES)
    }

    /**
     * Save user authentication state
     */
    fun saveAuthenticationState(isAuthenticated: Boolean) {
        saveBoolean(KEY_IS_AUTHENTICATED, isAuthenticated)
    }

    /**
     * Load user authentication state
     */
    fun loadAuthenticationState(): Boolean {
        return loadBoolean(KEY_IS_AUTHENTICATED)
    }

    /**
     * Save current user (encrypted)
     */
    fun saveCurrentUser(user: UserModel) {
        saveObject(user, KEY_CURRENT_USER, secure = true)
    }

    /**
     * Load current user (encrypted)
     */
    fun loadCurrentUser(): UserModel? {
        return loadObject<UserModel>(KEY_CURRENT_USER, secure = true)
    }

    /**
     * Save chat state
     */
    fun saveChatState(state: ChatState) {
        saveObject(state, KEY_CHAT_STATE)
    }

    /**
     * Load chat state
     */
    fun loadChatState(): ChatState? {
        return loadObject<ChatState>(KEY_CHAT_STATE)
    }

    /**
     * Save store state
     */
    fun saveStoreState(state: StoreState) {
        saveObject(state, KEY_STORE_STATE)
    }

    /**
     * Load store state
     */
    fun loadStoreState(): StoreState? {
        return loadObject<StoreState>(KEY_STORE_STATE)
    }

    /**
     * Save social state
     */
    fun saveSocialState(state: SocialState) {
        saveObject(state, KEY_SOCIAL_STATE)
    }

    /**
     * Load social state
     */
    fun loadSocialState(): SocialState? {
        return loadObject<SocialState>(KEY_SOCIAL_STATE)
    }

    /**
     * Save video state
     */
    fun saveVideoState(state: VideoState) {
        saveObject(state, KEY_VIDEO_STATE)
    }

    /**
     * Load video state
     */
    fun loadVideoState(): VideoState? {
        return loadObject<VideoState>(KEY_VIDEO_STATE)
    }

    /**
     * Save last sync timestamp
     */
    fun saveLastSyncTimestamp(timestamp: Long) {
        saveLong(KEY_LAST_SYNC_TIMESTAMP, timestamp)
    }

    /**
     * Load last sync timestamp
     */
    fun loadLastSyncTimestamp(): Long? {
        val timestamp = loadLong(KEY_LAST_SYNC_TIMESTAMP, -1L)
        return if (timestamp == -1L) null else timestamp
    }

    // MARK: - Secure Storage Methods

    /**
     * Save authentication token (encrypted)
     */
    fun saveAuthToken(token: String) {
        saveString("auth_token", token, secure = true)
    }

    /**
     * Load authentication token (encrypted)
     */
    fun loadAuthToken(): String? {
        return loadString("auth_token", secure = true)
    }

    /**
     * Remove authentication token
     */
    fun removeAuthToken() {
        remove("auth_token", secure = true)
    }

    /**
     * Save refresh token (encrypted)
     */
    fun saveRefreshToken(token: String) {
        saveString("refresh_token", token, secure = true)
    }

    /**
     * Load refresh token (encrypted)
     */
    fun loadRefreshToken(): String? {
        return loadString("refresh_token", secure = true)
    }

    /**
     * Remove refresh token
     */
    fun removeRefreshToken() {
        remove("refresh_token", secure = true)
    }

    /**
     * Save user credentials (encrypted)
     */
    fun saveUserCredentials(email: String, password: String) {
        saveString("user_email", email, secure = true)
        saveString("user_password", password, secure = true)
    }

    /**
     * Load user credentials (encrypted)
     */
    fun loadUserCredentials(): Pair<String?, String?> {
        val email = loadString("user_email", secure = true)
        val password = loadString("user_password", secure = true)
        return Pair(email, password)
    }

    /**
     * Remove user credentials
     */
    fun removeUserCredentials() {
        remove("user_email", secure = true)
        remove("user_password", secure = true)
    }

    // MARK: - Utility Methods

    /**
     * Check if key exists
     */
    fun contains(key: String, secure: Boolean = false): Boolean {
        val prefs = if (secure) encryptedPreferences else sharedPreferences
        return prefs?.contains(key) ?: false
    }

    /**
     * Get all keys
     */
    fun getAllKeys(secure: Boolean = false): Set<String> {
        val prefs = if (secure) encryptedPreferences else sharedPreferences
        return prefs?.all?.keys ?: emptySet()
    }

    /**
     * Get storage size (approximate)
     */
    fun getStorageSize(): Int {
        return (sharedPreferences?.all?.size ?: 0) + (encryptedPreferences?.all?.size ?: 0)
    }

    // MARK: - Model-specific save methods

    fun saveSession(session: UserSession) {
        saveObject(session, "current_session", secure = true)
    }

    fun saveWorkspace(workspace: Workspace) {
        saveObject(workspace, "current_workspace", secure = false)
    }

    fun loadSession(): UserSession? {
        return loadObject<UserSession>("current_session", secure = true)
    }

    fun loadWorkspace(): Workspace? {
        return loadObject<Workspace>("current_workspace", secure = false)
    }

    fun clearSession() {
        remove("current_session", secure = true)
        remove("current_workspace", secure = false)
    }
}