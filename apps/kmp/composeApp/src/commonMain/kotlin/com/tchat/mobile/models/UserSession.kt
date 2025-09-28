package com.tchat.mobile.models

import com.tchat.mobile.api.models.DeviceInfo
import kotlinx.serialization.Serializable
import kotlinx.datetime.*

/**
 * T024: UserSession data class
 *
 * Core user session management model with authentication and profile data.
 * Aligned with contract tests and cross-platform consistency requirements.
 */
@Serializable
data class UserSession(
    val accessToken: String,
    val refreshToken: String,
    val expiresIn: Long, // Token expiration time in seconds
    val expiresAt: String, // ISO 8601 timestamp when token expires
    val user: UserProfile,
    val sessionId: String? = null,
    val deviceInfo: DeviceInfo? = null,
    val loginTimestamp: String, // ISO 8601 timestamp of login
    val lastActivityTimestamp: String? = null,
    val preferences: UserPreferences? = null
)

@Serializable
data class UserProfile(
    val id: String,
    val email: String,
    val firstName: String? = null,
    val lastName: String? = null,
    val fullName: String? = null,
    val username: String? = null,
    val avatar: String? = null,
    val verified: Boolean = false,
    val phoneNumber: String? = null,
    val dateOfBirth: String? = null, // ISO 8601 date
    val bio: String? = null,
    val location: String? = null,
    val website: String? = null,
    val privacySettings: PrivacySettings = PrivacySettings(),
    val notificationSettings: NotificationSettings = NotificationSettings(),
    val createdAt: String,
    val updatedAt: String
)


@Serializable
data class UserPreferences(
    val language: String = "en", // ISO 639-1 language code
    val timezone: String = "UTC", // IANA timezone identifier
    val theme: String = "system", // "light" | "dark" | "system"
    val dateFormat: String = "YYYY-MM-DD",
    val timeFormat: String = "24h", // "12h" | "24h"
    val currency: String = "USD", // ISO 4217 currency code
    val measurementUnit: String = "metric", // "metric" | "imperial"
    val firstDayOfWeek: Int = 1, // 1 = Monday, 0 = Sunday
    val accessibility: AccessibilitySettings = AccessibilitySettings()
)

@Serializable
data class PrivacySettings(
    val profileVisibility: String = "public", // "public" | "friends" | "private"
    val showOnlineStatus: Boolean = true,
    val allowMessagesFrom: String = "everyone", // "everyone" | "friends" | "none"
    val allowFriendRequests: Boolean = true,
    val showEmail: Boolean = false,
    val showPhoneNumber: Boolean = false,
    val dataProcessingConsent: Boolean = false,
    val analyticsConsent: Boolean = false,
    val marketingConsent: Boolean = false
)

@Serializable
data class NotificationSettings(
    val enablePushNotifications: Boolean = true,
    val enableEmailNotifications: Boolean = true,
    val enableSMSNotifications: Boolean = false,
    val messageNotifications: Boolean = true,
    val commentNotifications: Boolean = true,
    val likeNotifications: Boolean = true,
    val friendRequestNotifications: Boolean = true,
    val systemNotifications: Boolean = true,
    val marketingNotifications: Boolean = false,
    val quietHours: QuietHours = QuietHours(),
    val soundEnabled: Boolean = true,
    val vibrationEnabled: Boolean = true,
    val badgeCountEnabled: Boolean = true
)

@Serializable
data class QuietHours(
    val enabled: Boolean = false,
    val startTime: String = "22:00", // HH:MM format
    val endTime: String = "08:00", // HH:MM format
    val daysOfWeek: List<Int> = listOf(1, 2, 3, 4, 5, 6, 7) // 1=Monday, 7=Sunday
)

@Serializable
data class AccessibilitySettings(
    val fontSize: String = "medium", // "small" | "medium" | "large" | "extra_large"
    val highContrast: Boolean = false,
    val reduceMotion: Boolean = false,
    val screenReaderEnabled: Boolean = false,
    val voiceOverEnabled: Boolean = false,
    val magnificationEnabled: Boolean = false,
    val colorBlindnessSupport: String = "none" // "none" | "protanopia" | "deuteranopia" | "tritanopia"
)

/**
 * Session state management enums
 */
enum class SessionState {
    AUTHENTICATED,
    UNAUTHENTICATED,
    REFRESHING_TOKEN,
    EXPIRED,
    INVALID,
    LOGOUT_PENDING
}

enum class LoginMethod {
    EMAIL_PASSWORD,
    GOOGLE,
    APPLE,
    FACEBOOK,
    PHONE_OTP,
    MAGIC_LINK,
    BIOMETRIC
}

/**
 * Extension functions for UserSession management
 */
fun UserSession.isExpired(): Boolean {
    val currentTime = Clock.System.now().epochSeconds
    val expirationTime = try {
        Instant.parse(expiresAt).epochSeconds
    } catch (e: Exception) {
        return true // If we can't parse expiration, consider expired for safety
    }
    return currentTime >= expirationTime
}

fun UserSession.needsRefresh(): Boolean {
    val currentTime = Clock.System.now().epochSeconds
    val expirationTime = try {
        Instant.parse(expiresAt).epochSeconds
    } catch (e: Exception) {
        return true
    }
    // Refresh if less than 5 minutes remaining
    return (expirationTime - currentTime) < 300
}

fun UserSession.timeUntilExpiry(): Long {
    val currentTime = Clock.System.now().epochSeconds
    val expirationTime = try {
        Instant.parse(expiresAt).epochSeconds
    } catch (e: Exception) {
        return 0L
    }
    return maxOf(0L, expirationTime - currentTime)
}

fun UserProfile.getDisplayName(): String {
    return when {
        fullName?.isNotBlank() == true -> fullName
        firstName?.isNotBlank() == true && lastName?.isNotBlank() == true -> "$firstName $lastName"
        firstName?.isNotBlank() == true -> firstName
        username?.isNotBlank() == true -> username
        else -> email.substringBefore("@")
    }
}

fun UserProfile.getInitials(): String {
    val displayName = getDisplayName()
    return if (displayName.contains(" ")) {
        displayName.split(" ").take(2).mapNotNull { it.firstOrNull()?.toString() }.joinToString("")
    } else {
        displayName.take(2)
    }.uppercase()
}