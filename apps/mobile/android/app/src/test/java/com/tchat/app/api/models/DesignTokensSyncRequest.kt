package com.tchat.app.api.models

/**
 * Design tokens sync request for testing
 */
data class DesignTokensSyncRequest(
    val platform: String,
    val version: String,
    val tokens: Map<String, Any>
)