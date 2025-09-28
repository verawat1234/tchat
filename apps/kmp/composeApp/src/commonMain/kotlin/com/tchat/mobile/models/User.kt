package com.tchat.mobile.models

import kotlinx.serialization.Serializable

/**
 * Domain model for authenticated user
 */
@Serializable
data class User(
    val id: String,
    val email: String,
    val displayName: String,
    val avatar: String? = null,
    val isActive: Boolean = true,
    val createdAt: Long,
    val updatedAt: Long
)

/**
 * Authentication state
 */
sealed class AuthState {
    object Loading : AuthState()
    object Unauthenticated : AuthState()
    data class Authenticated(val user: User) : AuthState()
    data class Error(val message: String) : AuthState()
}

/**
 * Authentication step for OTP flow
 */
enum class AuthStep {
    INPUT,    // Enter email/phone
    VERIFY    // Enter OTP code
}

/**
 * Authentication method
 */
enum class AuthMethod {
    EMAIL,    // Email + password
    PHONE     // Phone + OTP
}


/**
 * OTP request model
 */
@Serializable
data class OtpRequest(
    val phone: String,
    val countryCode: String = "+66"
)

/**
 * OTP verification model
 */
@Serializable
data class OtpVerification(
    val phone: String,
    val code: String,
    val countryCode: String = "+66"
)

/**
 * Authentication tokens
 */
@Serializable
data class AuthTokens(
    val accessToken: String,
    val refreshToken: String,
    val expiresAt: Long
)