package com.tchat.mobile.models

/**
 * User-related enums matching backend shared models
 * Based on backend/shared/models/user.go
 */

/**
 * User status enum
 * Matches UserStatus from backend
 */
enum class UserStatus(val value: String) {
    ACTIVE("active"),
    SUSPENDED("suspended"),
    DELETED("deleted");

    companion object {
        fun fromValue(value: String): UserStatus? {
            return values().find { it.value == value }
        }
    }
}

/**
 * KYC (Know Your Customer) tier enum
 * Matches KYCTier from backend
 */
enum class KYCTier(val value: Int, val displayName: String) {
    UNVERIFIED(0, "Unverified"),
    BASIC(1, "Basic"),
    STANDARD(2, "Standard"),
    PREMIUM(3, "Premium");

    companion object {
        fun fromValue(value: Int): KYCTier? {
            return values().find { it.value == value }
        }
    }
}

/**
 * Country codes for Southeast Asian countries
 * Matches Country from backend
 */
enum class Country(val code: String, val displayName: String) {
    THAILAND("TH", "Thailand"),
    SINGAPORE("SG", "Singapore"),
    INDONESIA("ID", "Indonesia"),
    MALAYSIA("MY", "Malaysia"),
    PHILIPPINES("PH", "Philippines"),
    VIETNAM("VN", "Vietnam");

    companion object {
        fun fromCode(code: String): Country? {
            return values().find { it.code == code }
        }

        fun getAllCountries(): List<Country> {
            return values().toList()
        }
    }
}

/**
 * Verification tier enum
 * Matches VerificationTier from backend
 */
enum class VerificationTier(val value: Int, val displayName: String) {
    NONE(0, "No Verification"),
    PHONE(1, "Phone Verified"),
    EMAIL(2, "Email Verified"),
    KYC(3, "KYC Verified"),
    FULL(4, "Full Verification");

    companion object {
        fun fromValue(value: Int): VerificationTier? {
            return values().find { it.value == value }
        }
    }
}