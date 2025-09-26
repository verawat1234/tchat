package com.tchat.mobile.data.models

import kotlinx.datetime.Instant
import kotlinx.serialization.Serializable
import com.benasher44.uuid.Uuid
import com.benasher44.uuid.uuid4

/**
 * User represents a user in the authentication system with Southeast Asian localization support
 */
@Serializable
data class User(
    val id: String = uuid4().toString(),
    val username: String? = null,
    val phone: String? = null,
    val phoneNumber: String? = null, // Alias for compatibility
    val email: String? = null,
    val name: String = "",
    val firstName: String? = null,
    val lastName: String? = null,
    val avatar: String? = null,
    val country: Country = Country.THAILAND,
    val countryCode: String? = null, // Alias for compatibility
    val locale: String = "th-TH",
    val language: String? = null,
    val timeZone: String? = null,
    val timezone: String? = null, // Alias for compatibility
    val bio: String? = null,
    val dateOfBirth: Instant? = null,
    val gender: String? = null,
    val isActive: Boolean = true,
    val kycStatus: KYCStatus = KYCStatus.NONE,
    val kycTier: VerificationTier = VerificationTier.NONE,
    val status: UserStatus = UserStatus.OFFLINE,
    val lastSeen: Instant? = null,
    val lastActiveAt: Instant? = null,
    val isVerified: Boolean = false,
    val isEmailVerified: Boolean = false,
    val isPhoneVerified: Boolean = false,
    val preferences: UserPreferences? = null,
    val metadata: Map<String, String> = emptyMap(),
    val createdAt: Instant = kotlinx.datetime.Clock.System.now(),
    val updatedAt: Instant = kotlinx.datetime.Clock.System.now()
) {
    /**
     * Get display name for the user
     */
    fun getDisplayName(): String {
        return when {
            name.isNotBlank() -> name
            phone != null -> phone!!
            email != null -> email!!
            else -> "Unknown User"
        }
    }

    /**
     * Check if user is currently online
     */
    fun isOnline(): Boolean = status == UserStatus.ONLINE

    /**
     * Check if user can send payments based on KYC tier
     */
    fun canSendPayments(): Boolean = kycTier >= VerificationTier.IDENTITY && isVerified

    /**
     * Get maximum daily limit based on KYC tier (in cents)
     */
    fun getMaxDailyLimit(): Long {
        return when (kycTier) {
            VerificationTier.BASIC -> 100000L // 1,000 THB equivalent
            VerificationTier.IDENTITY -> 5000000L // 50,000 THB equivalent
            VerificationTier.ENHANCED -> 50000000L // 500,000 THB equivalent
            else -> 0L
        }
    }

    /**
     * Check if user can update their profile
     */
    fun canUpdateProfile(): Boolean = status == UserStatus.ACTIVE || status == UserStatus.ONLINE

    /**
     * Convert to public user data for API responses
     */
    fun toPublicUser(): Map<String, Any> = mapOf(
        "id" to id,
        "name" to name,
        "avatar" to (avatar ?: ""),
        "country" to country.code,
        "status" to status.name.lowercase(),
        "is_verified" to isVerified,
        "created_at" to createdAt.toString()
    )
}

/**
 * Country represents supported Southeast Asian countries
 */
@Serializable
enum class Country(val code: String, val displayName: String) {
    THAILAND("TH", "Thailand"),
    INDONESIA("ID", "Indonesia"),
    MALAYSIA("MY", "Malaysia"),
    VIETNAM("VN", "Vietnam"),
    SINGAPORE("SG", "Singapore"),
    PHILIPPINES("PH", "Philippines");

    companion object {
        fun fromCode(code: String): Country? = values().find { it.code == code }

        fun validCountries(): List<Country> = values().toList()
    }

    /**
     * Get phone validation pattern for this country
     */
    fun getPhonePattern(): String {
        return when (this) {
            THAILAND -> "^\\+66[0-9]{8,9}$"
            INDONESIA -> "^\\+62[0-9]{8,12}$"
            MALAYSIA -> "^\\+60[0-9]{8,10}$"
            VIETNAM -> "^\\+84[0-9]{8,10}$"
            SINGAPORE -> "^\\+65[0-9]{8}$"
            PHILIPPINES -> "^\\+63[0-9]{9,10}$"
        }
    }

    /**
     * Get valid locales for this country
     */
    fun getValidLocales(): List<String> {
        return when (this) {
            THAILAND -> listOf("th-TH", "en-TH")
            INDONESIA -> listOf("id-ID", "en-ID")
            MALAYSIA -> listOf("ms-MY", "en-MY", "zh-MY")
            VIETNAM -> listOf("vi-VN", "en-VN")
            SINGAPORE -> listOf("en-SG", "ms-SG", "zh-SG", "ta-SG")
            PHILIPPINES -> listOf("fil-PH", "en-PH")
        }
    }
}

/**
 * KYC Status represents Know Your Customer verification status
 */
@Serializable
enum class KYCStatus {
    NONE,
    PENDING,
    APPROVED,
    REJECTED,
    EXPIRED
}

/**
 * Verification Tier represents KYC verification levels
 */
@Serializable
enum class VerificationTier(val level: Int) {
    NONE(0),
    BASIC(1), // Basic verification (phone/email)
    IDENTITY(2), // Identity verification (ID document)
    ENHANCED(3); // Enhanced verification (proof of address)

    companion object {
        fun fromLevel(level: Int): VerificationTier = values().find { it.level == level } ?: NONE
    }

    override fun toString(): String {
        return when (this) {
            BASIC -> "Basic"
            IDENTITY -> "Identity"
            ENHANCED -> "Enhanced"
            else -> "Unknown"
        }
    }
}

/**
 * User Status represents the current status of a user
 */
@Serializable
enum class UserStatus {
    ONLINE,
    OFFLINE,
    AWAY,
    BUSY,
    PENDING,
    ACTIVE,
    SUSPENDED,
    INACTIVE
}

/**
 * User Preferences for profile customization
 */
@Serializable
data class UserPreferences(
    val theme: String? = null,
    val language: String? = null,
    val notificationsEmail: Boolean = true,
    val notificationsPush: Boolean = true,
    val privacyLevel: String? = null
)

/**
 * User Response for API responses
 */
@Serializable
data class UserResponse(
    val id: String,
    val name: String,
    val avatar: String? = null,
    val country: String,
    val locale: String,
    val kycTier: Int,
    val status: String,
    val isVerified: Boolean,
    val createdAt: String,
    val updatedAt: String
)

/**
 * Validation utility functions
 */
object UserValidator {
    /**
     * Validate phone number format for a given country
     */
    fun isValidPhoneNumber(phone: String, country: Country): Boolean {
        return phone.matches(Regex(country.getPhonePattern()))
    }

    /**
     * Validate email format
     */
    fun isValidEmail(email: String): Boolean {
        val emailRegex = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
        return email.matches(Regex(emailRegex))
    }

    /**
     * Validate username format
     */
    fun isValidUsername(username: String): Boolean {
        val usernameRegex = "^[a-zA-Z0-9_]{3,30}$"
        return username.matches(Regex(usernameRegex))
    }

    /**
     * Validate locale for country
     */
    fun isValidLocale(locale: String, country: Country): Boolean {
        return country.getValidLocales().contains(locale)
    }

    /**
     * Format phone number for local display
     */
    fun formatPhoneNumberLocal(phone: String, country: Country): String {
        return when (country) {
            Country.THAILAND -> {
                if (phone.startsWith("+66")) {
                    val local = phone.removePrefix("+66")
                    if (local.length >= 8) {
                        "0${local.substring(0, 2)}-${local.substring(2, 5)}-${local.substring(5)}"
                    } else phone
                } else phone
            }
            Country.SINGAPORE -> {
                if (phone.startsWith("+65")) {
                    val local = phone.removePrefix("+65")
                    if (local.length == 8) {
                        "${local.substring(0, 4)} ${local.substring(4)}"
                    } else phone
                } else phone
            }
            Country.INDONESIA -> {
                if (phone.startsWith("+62")) {
                    val local = phone.removePrefix("+62")
                    "0$local"
                } else phone
            }
            else -> phone // Return international format for other countries
        }
    }
}