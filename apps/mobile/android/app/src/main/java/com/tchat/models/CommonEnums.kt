package com.tchat.models

import kotlinx.serialization.Serializable

/**
 * Common enums used across multiple model files
 */

/**
 * Access level for various resources and operations
 */
@Serializable
enum class AccessLevel(val value: String) {
    PUBLIC("public"),
    PRIVATE("private"),
    PROTECTED("protected"),
    INTERNAL("internal"),
    ADMIN("admin"),
    DEVELOPER("developer")
}

/**
 * Conflict resolution strategy for state management
 */
@Serializable
enum class ConflictResolutionStrategy {
    NEWEST_WINS,
    HIGHEST_VERSION_WINS,
    CLIENT_WINS,
    SERVER_WINS,
    MERGE,
    MANUAL,
    PROMPT
}