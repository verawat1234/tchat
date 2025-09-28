package com.tchat.mobile.services

import android.content.Context
import android.content.SharedPreferences
import androidx.security.crypto.EncryptedSharedPreferences
import androidx.security.crypto.MasterKey
import com.tchat.mobile.models.User
import kotlinx.serialization.encodeToString
import kotlinx.serialization.decodeFromString
import kotlinx.serialization.json.Json

/**
 * Android implementation of AuthPersistenceManager using EncryptedSharedPreferences
 * for secure storage of authentication tokens and user data
 */
actual class AuthPersistenceManager(private val context: Context) {

    constructor() : this(throw IllegalStateException("Context required for Android implementation"))

    companion object {
        private const val PREFS_NAME = "tchat_auth_prefs"
        private const val KEY_ACCESS_TOKEN = "access_token"
        private const val KEY_REFRESH_TOKEN = "refresh_token"
        private const val KEY_USER_DATA = "user_data"
        private const val KEY_IS_AUTHENTICATED = "is_authenticated"
    }

    private val masterKey = MasterKey.Builder(context)
        .setKeyScheme(MasterKey.KeyScheme.AES256_GCM)
        .build()

    private val encryptedPrefs: SharedPreferences = EncryptedSharedPreferences.create(
        context,
        PREFS_NAME,
        masterKey,
        EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
        EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
    )

    private val json = Json {
        ignoreUnknownKeys = true
        encodeDefaults = true
    }

    actual suspend fun saveAuthTokens(accessToken: String, refreshToken: String) {
        try {
            encryptedPrefs.edit()
                .putString(KEY_ACCESS_TOKEN, accessToken)
                .putString(KEY_REFRESH_TOKEN, refreshToken)
                .putBoolean(KEY_IS_AUTHENTICATED, true)
                .apply()
        } catch (e: Exception) {
            throw Exception("Failed to save auth tokens: ${e.message}")
        }
    }

    actual suspend fun getAccessToken(): String? {
        return try {
            encryptedPrefs.getString(KEY_ACCESS_TOKEN, null)
        } catch (e: Exception) {
            null
        }
    }

    actual suspend fun getRefreshToken(): String? {
        return try {
            encryptedPrefs.getString(KEY_REFRESH_TOKEN, null)
        } catch (e: Exception) {
            null
        }
    }

    actual suspend fun clearAuthTokens() {
        try {
            encryptedPrefs.edit()
                .remove(KEY_ACCESS_TOKEN)
                .remove(KEY_REFRESH_TOKEN)
                .putBoolean(KEY_IS_AUTHENTICATED, false)
                .apply()
        } catch (e: Exception) {
            throw Exception("Failed to clear auth tokens: ${e.message}")
        }
    }

    actual suspend fun saveUser(user: User) {
        try {
            val userJson = json.encodeToString(user)
            encryptedPrefs.edit()
                .putString(KEY_USER_DATA, userJson)
                .apply()
        } catch (e: Exception) {
            throw Exception("Failed to save user data: ${e.message}")
        }
    }

    actual suspend fun getUser(): User? {
        return try {
            val userJson = encryptedPrefs.getString(KEY_USER_DATA, null)
            userJson?.let { json.decodeFromString<User>(it) }
        } catch (e: Exception) {
            null
        }
    }

    actual suspend fun clearUser() {
        try {
            encryptedPrefs.edit()
                .remove(KEY_USER_DATA)
                .apply()
        } catch (e: Exception) {
            throw Exception("Failed to clear user data: ${e.message}")
        }
    }

    actual suspend fun isUserAuthenticated(): Boolean {
        return try {
            encryptedPrefs.getBoolean(KEY_IS_AUTHENTICATED, false) &&
                    getAccessToken() != null &&
                    getUser() != null
        } catch (e: Exception) {
            false
        }
    }
}