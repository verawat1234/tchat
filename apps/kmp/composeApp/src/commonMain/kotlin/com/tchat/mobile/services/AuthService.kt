package com.tchat.mobile.services

import com.tchat.mobile.api.ApiClient
import com.tchat.mobile.api.models.*
import com.tchat.mobile.models.*
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow

/**
 * Authentication service handling user authentication state and API calls
 */
class AuthService(
    private val apiClient: ApiClient,
    private val sessionManager: SessionManager
) {
    private val _authState = MutableStateFlow<AuthState>(AuthState.Loading)
    val authState: StateFlow<AuthState> = _authState.asStateFlow()

    private val _authStep = MutableStateFlow(AuthStep.INPUT)
    val authStep: StateFlow<AuthStep> = _authStep.asStateFlow()

    private val _authMethod = MutableStateFlow(AuthMethod.EMAIL)
    val authMethod: StateFlow<AuthMethod> = _authMethod.asStateFlow()

    private val _isLoading = MutableStateFlow(false)
    val isLoading: StateFlow<Boolean> = _isLoading.asStateFlow()

    // Store OTP request ID for verification
    private var otpRequestId: String? = null

    init {
        // Initially check if user is already authenticated
        checkAuthState()
    }

    /**
     * Initialize authentication service and restore session if available
     */
    suspend fun initialize() {
        _authState.value = AuthState.Loading

        try {
            val (accessToken, user) = sessionManager.initialize()

            if (accessToken != null && user != null) {
                // Validate token with server
                val result = apiClient.getCurrentUser()
                if (result.isSuccess) {
                    _authState.value = AuthState.Authenticated(user)
                } else {
                    // Token invalid, try to refresh
                    attemptTokenRefresh(user)
                }
            } else {
                _authState.value = AuthState.Unauthenticated
            }
        } catch (e: Exception) {
            _authState.value = AuthState.Unauthenticated
        }
    }

    /**
     * Attempt to refresh token and restore session
     */
    private suspend fun attemptTokenRefresh(user: User) {
        try {
            val result = apiClient.refreshAuthToken()
            if (result.isSuccess) {
                val refreshResponse = result.getOrThrow()
                sessionManager.updateAccessToken(refreshResponse.accessToken)
                _authState.value = AuthState.Authenticated(user)
            } else {
                // Refresh failed, clear session
                sessionManager.clearSession()
                _authState.value = AuthState.Unauthenticated
            }
        } catch (e: Exception) {
            sessionManager.clearSession()
            _authState.value = AuthState.Unauthenticated
        }
    }

    /**
     * Check current authentication state (e.g., from stored tokens)
     */
    private fun checkAuthState() {
        // Initial state is loading, actual initialization happens in initialize()
        _authState.value = AuthState.Loading
    }

    /**
     * Set authentication method
     */
    fun setAuthMethod(method: AuthMethod) {
        _authMethod.value = method
        _authStep.value = AuthStep.INPUT
    }

    /**
     * Login with email and password
     */
    suspend fun loginWithEmail(email: String, password: String): Result<User> {
        _isLoading.value = true

        return try {
            val result = apiClient.loginWithEmail(email, password)

            if (result.isSuccess) {
                val authResponse = result.getOrThrow()
                val user = User(
                    id = authResponse.user.id,
                    email = authResponse.user.email,
                    displayName = authResponse.user.displayName,
                    avatar = authResponse.user.avatar,
                    isActive = authResponse.user.isActive,
                    createdAt = authResponse.user.createdAt,
                    updatedAt = authResponse.user.updatedAt
                )

                // Save session
                sessionManager.saveSession(
                    authResponse.accessToken,
                    authResponse.refreshToken,
                    user
                )

                _authState.value = AuthState.Authenticated(user)
                _isLoading.value = false
                Result.success(user)
            } else {
                val error = result.exceptionOrNull()?.message ?: "Login failed"
                _authState.value = AuthState.Error(error)
                _isLoading.value = false
                Result.failure(result.exceptionOrNull() ?: Exception(error))
            }
        } catch (e: Exception) {
            _authState.value = AuthState.Error(e.message ?: "Login failed")
            _isLoading.value = false
            Result.failure(e)
        }
    }

    /**
     * Send OTP to phone number
     */
    suspend fun sendOtp(phone: String): Result<String> {
        _isLoading.value = true

        return try {
            // Extract country code and phone number
            val (countryCode, phoneNumber) = parsePhoneNumber(phone)

            val result = apiClient.requestOTP(phoneNumber, countryCode)

            if (result.isSuccess) {
                val otpResponse = result.getOrThrow()
                otpRequestId = otpResponse.requestID

                _authStep.value = AuthStep.VERIFY
                _isLoading.value = false
                Result.success(otpResponse.message)
            } else {
                val error = result.exceptionOrNull()?.message ?: "Failed to send OTP"
                _authState.value = AuthState.Error(error)
                _isLoading.value = false
                Result.failure(result.exceptionOrNull() ?: Exception(error))
            }
        } catch (e: Exception) {
            _authState.value = AuthState.Error(e.message ?: "Failed to send OTP")
            _isLoading.value = false
            Result.failure(e)
        }
    }

    /**
     * Parse phone number to extract country code and number
     * Returns ISO country code and phone number for backend API
     */
    private fun parsePhoneNumber(phone: String): Pair<String, String> {
        val cleanPhone = phone.replace(" ", "").replace("-", "").replace("(", "").replace(")", "")

        return when {
            cleanPhone.startsWith("+66") -> "TH" to cleanPhone.substring(3)  // Thailand
            cleanPhone.startsWith("+62") -> "ID" to cleanPhone.substring(3)  // Indonesia
            cleanPhone.startsWith("+63") -> "PH" to cleanPhone.substring(3)  // Philippines
            cleanPhone.startsWith("+84") -> "VN" to cleanPhone.substring(3)  // Vietnam
            cleanPhone.startsWith("+60") -> "MY" to cleanPhone.substring(3)  // Malaysia
            cleanPhone.startsWith("+65") -> "SG" to cleanPhone.substring(3)  // Singapore
            cleanPhone.startsWith("+1") -> "US" to cleanPhone.substring(2)   // US/Canada
            cleanPhone.startsWith("+") -> {
                // For other countries, extract the international dialing code
                val digits = cleanPhone.substring(1)
                val countryCodeLength = when {
                    digits.startsWith("1") -> 1
                    digits.startsWith("7") -> 1
                    digits.length >= 2 && digits.substring(0, 2).toIntOrNull() in 20..99 -> 2
                    digits.length >= 3 && digits.substring(0, 3).toIntOrNull() in 100..999 -> 3
                    else -> 2 // Default to 2-digit country code
                }
                // Return the international code for unknown countries
                "+${digits.substring(0, countryCodeLength)}" to digits.substring(countryCodeLength)
            }
            else -> "TH" to cleanPhone // Default to Thailand if no country code
        }
    }

    /**
     * Verify OTP code
     */
    suspend fun verifyOtp(phone: String, code: String): Result<User> {
        _isLoading.value = true

        return try {
            val requestId = otpRequestId
                ?: return Result.failure(Exception("No OTP request found. Please request OTP first."))

            val result = apiClient.verifyOTP(phone, code, requestId)

            if (result.isSuccess) {
                val verifyResponse = result.getOrThrow()
                val user = User(
                    id = verifyResponse.user.id,
                    email = verifyResponse.user.phoneNumber + "@sms.tchat.app", // Create email from phone
                    displayName = verifyResponse.user.displayName ?: "User ${verifyResponse.user.phoneNumber.takeLast(4)}",
                    avatar = verifyResponse.user.avatar,
                    isActive = verifyResponse.user.isActive,
                    createdAt = verifyResponse.user.createdAt,
                    updatedAt = verifyResponse.user.updatedAt
                )

                // Save session
                sessionManager.saveSession(
                    verifyResponse.accessToken,
                    verifyResponse.refreshToken,
                    user
                )

                _authState.value = AuthState.Authenticated(user)
                _isLoading.value = false

                // Clear the OTP request ID
                otpRequestId = null

                Result.success(user)
            } else {
                val error = result.exceptionOrNull()?.message ?: "Invalid OTP code"
                _authState.value = AuthState.Error(error)
                _isLoading.value = false
                Result.failure(result.exceptionOrNull() ?: Exception(error))
            }
        } catch (e: Exception) {
            _authState.value = AuthState.Error(e.message ?: "Verification failed")
            _isLoading.value = false
            Result.failure(e)
        }
    }

    /**
     * Go back to input step
     */
    fun goBackToInput() {
        _authStep.value = AuthStep.INPUT
    }

    /**
     * Logout user
     */
    suspend fun logout() {
        try {
            // Call API logout (best effort)
            apiClient.logout()
        } catch (e: Exception) {
            // Continue with local logout even if API call fails
        }

        // Clear session
        sessionManager.clearSession()

        _authState.value = AuthState.Unauthenticated
        _authStep.value = AuthStep.INPUT
        otpRequestId = null
    }

    /**
     * Get current authenticated user
     */
    fun getCurrentUser(): User? {
        return when (val state = _authState.value) {
            is AuthState.Authenticated -> state.user
            else -> null
        }
    }
}