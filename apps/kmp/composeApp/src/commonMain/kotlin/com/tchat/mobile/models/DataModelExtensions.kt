package com.tchat.mobile.models

import kotlinx.serialization.Serializable

/**
 * KMP Data Model Extensions and Utilities
 *
 * Additional abstractions for cross-platform data handling
 */

/**
 * Platform-agnostic result wrapper for KMP
 * Alternative to platform-specific Result types
 */
@Serializable
sealed class DataResult<out T> {
    @Serializable
    data class Success<T>(val data: T) : DataResult<T>()

    @Serializable
    data class Error(val message: String, val code: String? = null) : DataResult<Nothing>()

    @Serializable
    object Loading : DataResult<Nothing>()
}

/**
 * Platform-agnostic pagination wrapper
 */
@Serializable
data class PaginatedData<T>(
    val items: List<T>,
    val currentPage: Int,
    val totalPages: Int,
    val totalItems: Int,
    val hasNextPage: Boolean,
    val hasPreviousPage: Boolean
)

/**
 * Cross-platform cache metadata
 */
@Serializable
data class CacheMetadata(
    val timestamp: Long,
    val expiryTime: Long,
    val version: String,
    val source: CacheSource
)

enum class CacheSource {
    NETWORK,
    DISK,
    MEMORY,
    FALLBACK
}

/**
 * Platform-agnostic file metadata
 */
@Serializable
data class FileMetadata(
    val name: String,
    val size: Long,
    val mimeType: String,
    val extension: String,
    val lastModified: Long? = null
)

/**
 * Cross-platform error types
 */
enum class AppErrorType {
    NETWORK,
    VALIDATION,
    AUTHENTICATION,
    AUTHORIZATION,
    NOT_FOUND,
    SERVER_ERROR,
    UNKNOWN
}

/**
 * KMP-friendly error class
 */
@Serializable
data class AppError(
    val type: AppErrorType,
    val message: String,
    val code: String? = null,
    val details: Map<String, String> = emptyMap()
)

/**
 * Extension functions for better KMP usability
 */

// Post extensions for cross-platform usage
fun Post.toPlatformShareableData(): Map<String, Any> {
    return mapOf(
        "title" to (content.text ?: "Check out this post"),
        "url" to "https://tchat.app/posts/$id",
        "imageUrl" to (content.images.firstOrNull()?.url ?: ""),
        "type" to type.name
    )
}

// Data validation for cross-platform consistency
fun Post.validateForPlatform(): List<String> {
    val errors = mutableListOf<String>()

    // Cross-platform validation rules
    if (id.isEmpty()) errors.add("Post ID is required")
    if (user.id.isEmpty()) errors.add("User ID is required")

    // Platform-specific content validation
    when (type) {
        PostType.IMAGE -> {
            if (content.images.isEmpty()) {
                errors.add("Image posts must contain at least one image")
            }
        }
        PostType.VIDEO -> {
            if (content.videos.isEmpty()) {
                errors.add("Video posts must contain video content")
            }
        }
        PostType.TEXT -> {
            if (content.text.isNullOrBlank()) {
                errors.add("Text posts must contain text content")
            }
        }
        else -> {
            // Other types validated by existing PostValidationRules
        }
    }

    return errors
}

// Cross-platform serialization helpers
fun Post.toJsonString(): String {
    return kotlinx.serialization.json.Json.encodeToString(Post.serializer(), this)
}

fun String.toPost(): Post? {
    return try {
        kotlinx.serialization.json.Json.decodeFromString(Post.serializer(), this)
    } catch (e: Exception) {
        null
    }
}

// Platform-agnostic data transformation
fun List<Post>.filterByPlatformCapabilities(platform: String): List<Post> {
    return when (platform.lowercase()) {
        "ios" -> {
            // iOS-specific filtering
            this.filter { post ->
                // Example: iOS might have different video format support
                if (post.type == PostType.VIDEO) {
                    post.content.videos.all { it.url.endsWith(".mp4") || it.url.endsWith(".mov") }
                } else {
                    true
                }
            }
        }
        "android" -> {
            // Android-specific filtering
            this.filter { post ->
                // Example: Android might support more formats
                true // More permissive
            }
        }
        else -> this
    }
}

/**
 * KMP-friendly constants
 */
object PlatformConstants {
    const val MAX_POST_TEXT_LENGTH = 280
    const val MAX_IMAGE_SIZE_MB = 10
    const val MAX_VIDEO_SIZE_MB = 100
    const val MAX_IMAGES_PER_POST = 10
    const val MAX_VIDEOS_PER_POST = 1

    val SUPPORTED_IMAGE_FORMATS = listOf("jpg", "jpeg", "png", "gif", "webp")
    val SUPPORTED_VIDEO_FORMATS = listOf("mp4", "mov", "avi", "mkv")
}

/**
 * Cross-platform utility functions
 */
object KmpDataUtils {

    fun generateId(): String {
        // Simple UUID-like generation for cross-platform compatibility
        return "${kotlinx.datetime.Clock.System.now().toEpochMilliseconds()}-${(1000..9999).random()}"
    }

    fun getCurrentTimestamp(): String {
        return kotlinx.datetime.Clock.System.now().toString()
    }

    fun isValidUrl(url: String): Boolean {
        return url.startsWith("http://") || url.startsWith("https://")
    }

    fun sanitizeText(text: String): String {
        return text.trim()
            .replace(Regex("\\s+"), " ") // Normalize whitespace
            .take(PlatformConstants.MAX_POST_TEXT_LENGTH)
    }
}