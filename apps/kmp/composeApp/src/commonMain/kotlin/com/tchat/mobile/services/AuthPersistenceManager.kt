package com.tchat.mobile.services

import com.tchat.mobile.models.User
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow

/**
 * Manages persistence of authentication data across app sessions
 *
 * This is a common interface for storing authentication tokens and user data.
 * Platform-specific implementations will handle secure storage:
 * - Android: EncryptedSharedPreferences
 * - iOS: Keychain Services
 */
expect class AuthPersistenceManager {
    suspend fun saveAuthTokens(accessToken: String, refreshToken: String)
    suspend fun getAccessToken(): String?
    suspend fun getRefreshToken(): String?
    suspend fun clearAuthTokens()

    suspend fun saveUser(user: User)
    suspend fun getUser(): User?
    suspend fun clearUser()

    suspend fun isUserAuthenticated(): Boolean
}

/**
 * Authentication data holder for persistence
 */
data class AuthData(
    val accessToken: String,
    val refreshToken: String,
    val user: User,
    val expiresAt: Long
)

/**
 * Session management service for handling authentication persistence
 */
class SessionManager(
    private val persistenceManager: AuthPersistenceManager
) {
    private val _isInitialized = MutableStateFlow(false)
    val isInitialized: StateFlow<Boolean> = _isInitialized.asStateFlow()

    /**
     * Initialize session manager and restore authentication state if available
     */
    suspend fun initialize(): Pair<String?, User?> {
        return try {
            val accessToken = persistenceManager.getAccessToken()
            val user = persistenceManager.getUser()

            _isInitialized.value = true

            if (accessToken != null && user != null) {
                accessToken to user
            } else {
                null to null
            }
        } catch (e: Exception) {
            _isInitialized.value = true
            null to null
        }
    }

    /**
     * Save authentication session
     */
    suspend fun saveSession(
        accessToken: String,
        refreshToken: String,
        user: User
    ): Result<Unit> {
        return try {
            persistenceManager.saveAuthTokens(accessToken, refreshToken)
            persistenceManager.saveUser(user)
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Clear authentication session
     */
    suspend fun clearSession(): Result<Unit> {
        return try {
            persistenceManager.clearAuthTokens()
            persistenceManager.clearUser()
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Get stored tokens
     */
    suspend fun getTokens(): Pair<String?, String?> {
        return try {
            val accessToken = persistenceManager.getAccessToken()
            val refreshToken = persistenceManager.getRefreshToken()
            accessToken to refreshToken
        } catch (e: Exception) {
            null to null
        }
    }

    /**
     * Update access token (for token refresh)
     */
    suspend fun updateAccessToken(newAccessToken: String): Result<Unit> {
        return try {
            val refreshToken = persistenceManager.getRefreshToken()
            if (refreshToken != null) {
                persistenceManager.saveAuthTokens(newAccessToken, refreshToken)
                Result.success(Unit)
            } else {
                Result.failure(Exception("No refresh token found"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Check if user is authenticated
     */
    suspend fun isAuthenticated(): Boolean {
        return try {
            persistenceManager.isUserAuthenticated()
        } catch (e: Exception) {
            false
        }
    }
}