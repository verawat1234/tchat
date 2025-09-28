package com.tchat.mobile.models

import kotlinx.serialization.json.Json
import kotlinx.serialization.encodeToString
import kotlinx.serialization.decodeFromString
import kotlinx.serialization.SerializationException
import com.tchat.mobile.database.SyncMetadata

/**
 * JSON serialization extensions for sync operations
 *
 * Provides safe serialization/deserialization utilities for:
 * - Chat messages and metadata
 * - Sync operation data
 * - Error handling with fallback mechanisms
 * - Type-safe conversions
 */

internal val json = Json {
    ignoreUnknownKeys = true
    encodeDefaults = true
    isLenient = true
}

/**
 * Safely serialize an object to JSON string
 */
inline fun <reified T> T.toJsonString(): String {
    return try {
        Json {
            ignoreUnknownKeys = true
            encodeDefaults = true
            isLenient = true
        }.encodeToString(this)
    } catch (e: SerializationException) {
        "{\"error\": \"serialization_failed\", \"type\": \"${T::class.simpleName}\"}"
    }
}

/**
 * Safely deserialize JSON string to object
 */
inline fun <reified T> String.fromJsonString(): T? {
    return try {
        Json {
            ignoreUnknownKeys = true
            encodeDefaults = true
            isLenient = true
        }.decodeFromString<T>(this)
    } catch (e: SerializationException) {
        null
    }
}

/**
 * Convert string list to JSON array string
 */
fun List<String>.toJsonArray(): String {
    return try {
        Json {
            ignoreUnknownKeys = true
            encodeDefaults = true
            isLenient = true
        }.encodeToString(this)
    } catch (e: SerializationException) {
        "[]"
    }
}

/**
 * Convert JSON array string to string list
 */
fun String.fromJsonArray(): List<String> {
    return try {
        Json {
            ignoreUnknownKeys = true
            encodeDefaults = true
            isLenient = true
        }.decodeFromString<List<String>>(this)
    } catch (e: SerializationException) {
        emptyList()
    }
}

/**
 * Convert map to JSON object string
 */
fun Map<String, String>.toJsonObject(): String {
    return try {
        Json {
            ignoreUnknownKeys = true
            encodeDefaults = true
            isLenient = true
        }.encodeToString(this)
    } catch (e: SerializationException) {
        "{}"
    }
}

/**
 * Convert JSON object string to map
 */
fun String.fromJsonObject(): Map<String, String> {
    return try {
        Json {
            ignoreUnknownKeys = true
            encodeDefaults = true
            isLenient = true
        }.decodeFromString<Map<String, String>>(this)
    } catch (e: SerializationException) {
        emptyMap()
    }
}

/**
 * Safe JSON serialization for sync metadata
 */
fun SyncMetadata.toJson(): String = this.toJsonString()

/**
 * Safe JSON deserialization for sync metadata
 */
fun String.toSyncMetadata(): SyncMetadata? = this.fromJsonString<SyncMetadata>()

/**
 * Safe JSON serialization for sync operations
 */
fun SyncOperation.toJson(): String = this.toJsonString()

/**
 * Safe JSON deserialization for sync operations
 */
fun String.toSyncOperation(): SyncOperation? = this.fromJsonString<SyncOperation>()

/**
 * Safe JSON serialization for message conflicts
 */
fun MessageConflict.toJson(): String = this.toJsonString()

/**
 * Safe JSON deserialization for message conflicts
 */
fun String.toMessageConflict(): MessageConflict? = this.fromJsonString<MessageConflict>()

/**
 * Utility for handling nullable JSON serialization
 */
inline fun <reified T> T?.toJsonStringOrNull(): String? {
    return this?.toJsonString()
}

/**
 * Utility for handling nullable JSON deserialization
 */
inline fun <reified T> String?.fromJsonStringOrNull(): T? {
    return this?.fromJsonString<T>()
}

/**
 * Convert any serializable object to pretty JSON string for debugging
 */
inline fun <reified T> T.toPrettyJsonString(): String {
    return try {
        Json {
            prettyPrint = true
            ignoreUnknownKeys = true
            encodeDefaults = true
        }.encodeToString(this)
    } catch (e: SerializationException) {
        "{\n  \"error\": \"serialization_failed\",\n  \"type\": \"${T::class.simpleName}\"\n}"
    }
}

/**
 * Validate JSON string format
 */
fun String.isValidJson(): Boolean {
    return try {
        Json {
            ignoreUnknownKeys = true
            encodeDefaults = true
            isLenient = true
        }.parseToJsonElement(this)
        true
    } catch (e: SerializationException) {
        false
    }
}

/**
 * Get safe JSON size for storage estimation
 */
fun String.jsonByteSize(): Int {
    return try {
        this.toByteArray(Charsets.UTF_8).size
    } catch (e: Exception) {
        0
    }
}