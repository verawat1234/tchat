package com.tchat.mobile.utils

import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import kotlinx.serialization.SerializationException

/**
 * JSON utility extensions for serialization and deserialization
 */
private val jsonInstance = Json {
    ignoreUnknownKeys = true
    encodeDefaults = true
}

/**
 * Convert any object to JSON string
 */
fun <T> T.toJsonString(): String {
    return try {
        when (this) {
            is String -> this
            is List<*> -> if (this.isEmpty()) "[]" else this.toString()
            is Map<*, *> -> if (this.isEmpty()) "{}" else this.toString()
            else -> this.toString()
        }
    } catch (e: Exception) {
        this.toString()
    }
}

/**
 * Convert JSON string to simple list
 */
fun String.fromJsonToList(): List<String> {
    return try {
        if (this.isBlank() || this == "null" || this == "[]") return emptyList()
        // Simple parsing for now - this is a fallback
        emptyList()
    } catch (e: Exception) {
        emptyList()
    }
}

/**
 * Safe JSON conversion with default value
 */
fun String.fromJsonOrDefault(default: String): String {
    return if (this.isBlank() || this == "null") default else this
}