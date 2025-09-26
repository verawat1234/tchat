package com.tchat.mobile.data.network

import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.request.*
import io.ktor.http.*
import com.tchat.mobile.data.models.User
import com.tchat.mobile.data.models.Country
import kotlinx.serialization.Serializable
import io.github.aakira.napier.Napier

/**
 * Authentication API client for user authentication and profile management
 */
class AuthApiClient(
    private val httpClient: HttpClient,
    private val baseUrl: String
) {

    /**
     * Login with email/phone and password
     */
    suspend fun login(request: LoginRequest): ApiResult<LoginResponse> {
        return try {
            val response: ApiResponse<LoginResponse> = httpClient.post("$baseUrl/auth/login") {
                contentType(ContentType.Application.Json)
                setBody(request)
            }.body()

            if (response.success && response.data != null) {
                ApiResult.Success(response.data!!)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Login failed"))
            }
        } catch (e: Exception) {
            Napier.e("Login error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Register new user with phone or email
     */
    suspend fun register(request: RegisterRequest): ApiResult<RegisterResponse> {
        return try {
            val response: ApiResponse<RegisterResponse> = httpClient.post("$baseUrl/auth/register") {
                contentType(ContentType.Application.Json)
                setBody(request)
            }.body()

            if (response.success && response.data != null) {
                ApiResult.Success(response.data!!)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Registration failed"))
            }
        } catch (e: Exception) {
            Napier.e("Registration error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Verify OTP for phone/email verification
     */
    suspend fun verifyOtp(request: OtpVerificationRequest): ApiResult<OtpVerificationResponse> {
        return try {
            val response: ApiResponse<OtpVerificationResponse> = httpClient.post("$baseUrl/auth/verify-otp") {
                contentType(ContentType.Application.Json)
                setBody(request)
            }.body()

            if (response.success && response.data != null) {
                ApiResult.Success(response.data!!)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "OTP verification failed"))
            }
        } catch (e: Exception) {
            Napier.e("OTP verification error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Refresh access token using refresh token
     */
    suspend fun refreshToken(refreshToken: String): ApiResult<TokenResponse> {
        return try {
            val response: ApiResponse<TokenResponse> = httpClient.post("$baseUrl/auth/refresh") {
                contentType(ContentType.Application.Json)
                setBody(RefreshTokenRequest(refreshToken))
            }.body()

            if (response.success && response.data != null) {
                ApiResult.Success(response.data!!)
            } else {
                ApiResult.Error(ApiError.UnauthorizedError("Token refresh failed"))
            }
        } catch (e: Exception) {
            Napier.e("Token refresh error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Logout user and invalidate tokens
     */
    suspend fun logout(token: String): ApiResult<Unit> {
        return try {
            val response: ApiResponse<Unit> = httpClient.post("$baseUrl/auth/logout") {
                header("Authorization", "Bearer $token")
                contentType(ContentType.Application.Json)
            }.body()

            if (response.success) {
                ApiResult.Success(Unit)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Logout failed"))
            }
        } catch (e: Exception) {
            Napier.e("Logout error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Get current user profile
     */
    suspend fun getProfile(token: String): ApiResult<User> {
        return try {
            val response: ApiResponse<User> = httpClient.get("$baseUrl/auth/profile") {
                header("Authorization", "Bearer $token")
            }.body()

            if (response.success && response.data != null) {
                ApiResult.Success(response.data!!)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to get profile"))
            }
        } catch (e: Exception) {
            Napier.e("Get profile error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Update user profile
     */
    suspend fun updateProfile(token: String, request: UpdateProfileRequest): ApiResult<User> {
        return try {
            val response: ApiResponse<User> = httpClient.put("$baseUrl/auth/profile") {
                header("Authorization", "Bearer $token")
                contentType(ContentType.Application.Json)
                setBody(request)
            }.body()

            if (response.success && response.data != null) {
                ApiResult.Success(response.data!!)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to update profile"))
            }
        } catch (e: Exception) {
            Napier.e("Update profile error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Change user password
     */
    suspend fun changePassword(token: String, request: ChangePasswordRequest): ApiResult<Unit> {
        return try {
            val response: ApiResponse<Unit> = httpClient.put("$baseUrl/auth/change-password") {
                header("Authorization", "Bearer $token")
                contentType(ContentType.Application.Json)
                setBody(request)
            }.body()

            if (response.success) {
                ApiResult.Success(Unit)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to change password"))
            }
        } catch (e: Exception) {
            Napier.e("Change password error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Request password reset via email/phone
     */
    suspend fun requestPasswordReset(request: PasswordResetRequest): ApiResult<Unit> {
        return try {
            val response: ApiResponse<Unit> = httpClient.post("$baseUrl/auth/forgot-password") {
                contentType(ContentType.Application.Json)
                setBody(request)
            }.body()

            if (response.success) {
                ApiResult.Success(Unit)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to request password reset"))
            }
        } catch (e: Exception) {
            Napier.e("Password reset request error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Reset password with OTP
     */
    suspend fun resetPassword(request: ResetPasswordRequest): ApiResult<Unit> {
        return try {
            val response: ApiResponse<Unit> = httpClient.post("$baseUrl/auth/reset-password") {
                contentType(ContentType.Application.Json)
                setBody(request)
            }.body()

            if (response.success) {
                ApiResult.Success(Unit)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to reset password"))
            }
        } catch (e: Exception) {
            Napier.e("Password reset error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Submit KYC verification documents
     */
    suspend fun submitKyc(token: String, request: KycSubmissionRequest): ApiResult<KycResponse> {
        return try {
            val response: ApiResponse<KycResponse> = httpClient.post("$baseUrl/auth/kyc/submit") {
                header("Authorization", "Bearer $token")
                contentType(ContentType.Application.Json)
                setBody(request)
            }.body()

            if (response.success && response.data != null) {
                ApiResult.Success(response.data!!)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "KYC submission failed"))
            }
        } catch (e: Exception) {
            Napier.e("KYC submission error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Get KYC status
     */
    suspend fun getKycStatus(token: String): ApiResult<KycResponse> {
        return try {
            val response: ApiResponse<KycResponse> = httpClient.get("$baseUrl/auth/kyc/status") {
                header("Authorization", "Bearer $token")
            }.body()

            if (response.success && response.data != null) {
                ApiResult.Success(response.data!!)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to get KYC status"))
            }
        } catch (e: Exception) {
            Napier.e("Get KYC status error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Delete user account
     */
    suspend fun deleteAccount(token: String, request: DeleteAccountRequest): ApiResult<Unit> {
        return try {
            val response: ApiResponse<Unit> = httpClient.delete("$baseUrl/auth/account") {
                header("Authorization", "Bearer $token")
                contentType(ContentType.Application.Json)
                setBody(request)
            }.body()

            if (response.success) {
                ApiResult.Success(Unit)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to delete account"))
            }
        } catch (e: Exception) {
            Napier.e("Delete account error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Handle common exceptions and convert to ApiError
     */
    private fun handleException(e: Exception): ApiError {
        return when (e) {
            is io.ktor.client.plugins.ClientRequestException -> {
                when (e.response.status.value) {
                    400 -> ApiError.ClientError(400, "Bad request")
                    401 -> ApiError.UnauthorizedError("Unauthorized")
                    403 -> ApiError.ForbiddenError("Forbidden")
                    404 -> ApiError.NotFoundError("Not found")
                    else -> ApiError.ClientError(e.response.status.value, e.message)
                }
            }
            is io.ktor.client.plugins.ServerResponseException -> {
                ApiError.ServerError(e.response.status.value, "Server error: ${e.message}")
            }
            is kotlinx.serialization.SerializationException -> {
                ApiError.SerializationError("Serialization error: ${e.message}")
            }
            else -> {
                ApiError.NetworkError("Network error: ${e.message}")
            }
        }
    }
}

// Request/Response models for Authentication API

@Serializable
data class LoginRequest(
    val email: String? = null,
    val phone: String? = null,
    val password: String,
    val country: String? = null,
    val deviceId: String? = null,
    val deviceInfo: String? = null
)

@Serializable
data class LoginResponse(
    val user: User,
    val tokens: TokenResponse,
    val sessionId: String,
    val expiresAt: String,
    val isFirstLogin: Boolean = false
)

@Serializable
data class RegisterRequest(
    val name: String,
    val email: String? = null,
    val phone: String? = null,
    val password: String,
    val country: String,
    val locale: String = "en",
    val referralCode: String? = null,
    val deviceId: String? = null,
    val acceptTerms: Boolean = true
)

@Serializable
data class RegisterResponse(
    val user: User,
    val tokens: TokenResponse? = null,
    val needsVerification: Boolean = true,
    val verificationMethod: String = "otp", // otp, email
    val sessionId: String? = null
)

@Serializable
data class TokenResponse(
    val accessToken: String,
    val refreshToken: String,
    val tokenType: String = "Bearer",
    val expiresIn: Int, // seconds
    val refreshExpiresIn: Int, // seconds
    val scope: String = "all"
)

@Serializable
data class OtpVerificationRequest(
    val phone: String? = null,
    val email: String? = null,
    val otp: String,
    val type: String = "registration" // registration, password_reset, phone_verification
)

@Serializable
data class OtpVerificationResponse(
    val verified: Boolean,
    val user: User? = null,
    val tokens: TokenResponse? = null,
    val message: String? = null
)

@Serializable
data class RefreshTokenRequest(
    val refreshToken: String
)

@Serializable
data class UpdateProfileRequest(
    val name: String? = null,
    val displayName: String? = null,
    val firstName: String? = null,
    val lastName: String? = null,
    val bio: String? = null,
    val avatar: String? = null,
    val dateOfBirth: String? = null,
    val gender: String? = null,
    val language: String? = null,
    val timezone: String? = null,
    val preferences: Map<String, String> = emptyMap()
)

@Serializable
data class ChangePasswordRequest(
    val currentPassword: String,
    val newPassword: String,
    val confirmPassword: String
)

@Serializable
data class PasswordResetRequest(
    val email: String? = null,
    val phone: String? = null,
    val country: String? = null
)

@Serializable
data class ResetPasswordRequest(
    val email: String? = null,
    val phone: String? = null,
    val otp: String,
    val newPassword: String,
    val confirmPassword: String
)

@Serializable
data class KycSubmissionRequest(
    val tier: Int, // 1, 2, 3
    val documents: List<KycDocument>,
    val personalInfo: KycPersonalInfo? = null,
    val addressInfo: KycAddressInfo? = null
)

@Serializable
data class KycDocument(
    val type: String, // id_card, passport, driver_license, proof_of_address
    val frontImageUrl: String,
    val backImageUrl: String? = null,
    val documentNumber: String? = null,
    val expiryDate: String? = null,
    val issueDate: String? = null,
    val issuer: String? = null
)

@Serializable
data class KycPersonalInfo(
    val firstName: String,
    val lastName: String,
    val middleName: String? = null,
    val dateOfBirth: String,
    val placeOfBirth: String? = null,
    val nationality: String,
    val occupation: String? = null,
    val sourceOfIncome: String? = null
)

@Serializable
data class KycAddressInfo(
    val address1: String,
    val address2: String? = null,
    val city: String,
    val state: String? = null,
    val postalCode: String,
    val country: String
)

@Serializable
data class KycResponse(
    val tier: Int,
    val status: String, // pending, approved, rejected, expired
    val message: String? = null,
    val documents: List<KycDocument> = emptyList(),
    val submittedAt: String? = null,
    val reviewedAt: String? = null,
    val expiresAt: String? = null
)

@Serializable
data class DeleteAccountRequest(
    val password: String,
    val reason: String? = null,
    val feedback: String? = null
)