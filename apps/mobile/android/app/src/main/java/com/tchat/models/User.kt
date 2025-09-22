package com.tchat.models

import kotlinx.serialization.Serializable

/**
 * Simple User model for session management
 */
@Serializable
data class User(
    val id: String,
    val name: String,
    val email: String? = null
)

/**
 * Workspace model for session management
 */
@Serializable
data class Workspace(
    val id: String,
    val name: String,
    val description: String? = null,
    val iconUrl: String? = null,
    val memberCount: Int = 0,
    val isPersonal: Boolean = false
)