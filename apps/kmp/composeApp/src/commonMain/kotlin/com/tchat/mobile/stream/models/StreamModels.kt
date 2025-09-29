package com.tchat.mobile.stream.models

import kotlinx.serialization.Serializable

/**
 * Stream Store Tabs - Kotlin Multiplatform Models
 * Cross-platform data models for Stream content system
 */

@Serializable
data class StreamCategory(
    val id: String,
    val name: String,
    val displayOrder: Int,
    val iconName: String,
    val isActive: Boolean,
    val subtabs: List<StreamSubtab>? = null,
    val featuredContentEnabled: Boolean,
    val createdAt: String,
    val updatedAt: String
)

@Serializable
data class StreamSubtab(
    val id: String,
    val categoryId: String,
    val name: String,
    val displayOrder: Int,
    val filterCriteria: Map<String, String>, // Simplified from JSON for KMP compatibility
    val isActive: Boolean,
    val createdAt: String,
    val updatedAt: String
)

@Serializable
enum class StreamContentType {
    BOOK,
    PODCAST,
    CARTOON,
    SHORT_MOVIE,
    LONG_MOVIE,
    MUSIC,
    ART
}

@Serializable
enum class StreamAvailabilityStatus {
    AVAILABLE,
    COMING_SOON,
    UNAVAILABLE
}

@Serializable
data class StreamContentItem(
    val id: String,
    val categoryId: String,
    val title: String,
    val description: String,
    val thumbnailUrl: String,
    val contentType: StreamContentType,
    val duration: Int? = null, // in seconds, null for books
    val price: Double,
    val currency: String,
    val availabilityStatus: StreamAvailabilityStatus,
    val isFeatured: Boolean,
    val featuredOrder: Int? = null,
    val metadata: Map<String, String>, // Simplified from JSON for KMP compatibility
    val createdAt: String,
    val updatedAt: String
) {
    fun isBook(): Boolean = contentType == StreamContentType.BOOK

    fun isVideo(): Boolean = contentType in listOf(
        StreamContentType.SHORT_MOVIE,
        StreamContentType.LONG_MOVIE,
        StreamContentType.CARTOON
    )

    fun isAudio(): Boolean = contentType in listOf(
        StreamContentType.PODCAST,
        StreamContentType.MUSIC
    )

    fun isAvailable(): Boolean = availabilityStatus == StreamAvailabilityStatus.AVAILABLE

    fun canPurchase(): Boolean = isAvailable()

    fun getDurationString(): String {
        return duration?.let {
            val hours = it / 3600
            val minutes = (it % 3600) / 60
            val seconds = it % 60

            if (hours > 0) {
                String.format("%d:%02d:%02d", hours, minutes, seconds)
            } else {
                String.format("%d:%02d", minutes, seconds)
            }
        } ?: ""
    }
}

@Serializable
data class StreamProduct(
    val id: String,
    val name: String,
    val description: String,
    val price: Double,
    val currency: String,
    val productType: ProductType,
    val mediaContentId: String? = null,
    val mediaMetadata: MediaMetadata? = null,
    val category: String,
    val isActive: Boolean,
    val stockQuantity: Int? = null,
    val createdAt: String,
    val updatedAt: String
)

@Serializable
enum class ProductType {
    PHYSICAL,
    MEDIA
}

@Serializable
data class MediaMetadata(
    val contentType: StreamContentType,
    val duration: Int? = null,
    val format: String? = null,
    val license: String? = null
)

@Serializable
data class StreamCartItem(
    val id: String,
    val cartId: String,
    val productId: String,
    val mediaContentId: String? = null,
    val quantity: Int,
    val unitPrice: Double,
    val totalPrice: Double,
    val mediaLicense: MediaLicense? = null,
    val downloadFormat: DownloadFormat? = null,
    val createdAt: String,
    val updatedAt: String
)

@Serializable
enum class MediaLicense {
    PERSONAL,
    FAMILY
}

@Serializable
enum class DownloadFormat {
    PDF,
    EPUB,
    MP3,
    MP4,
    FLAC
}

@Serializable
data class ContentCollection(
    val id: String,
    val name: String,
    val categoryId: String,
    val collectionType: CollectionType,
    val displayOrder: Int,
    val isActive: Boolean,
    val itemIds: List<String>,
    val maxItems: Int,
    val createdAt: String,
    val updatedAt: String
)

@Serializable
enum class CollectionType {
    FEATURED,
    NEW_RELEASES,
    TRENDING,
    CURATED
}

// Navigation state management
@Serializable
data class TabNavigationState(
    val userId: String,
    val currentCategoryId: String,
    val currentSubtabId: String? = null,
    val lastVisitedAt: String,
    val sessionId: String
)

// API Response types - matching backend structure
@Serializable
data class StreamCategoriesResponse(
    val categories: List<StreamCategory>,
    val total: Int,
    val success: Boolean
)

@Serializable
data class StreamContentResponse(
    val content: List<StreamContentItem>,
    val page: Int,
    val limit: Int,
    val total: Int,
    val success: Boolean
)

@Serializable
data class StreamFeaturedResponse(
    val content: List<StreamContentItem>,
    val total: Int,
    val success: Boolean
)

@Serializable
data class StreamSubtabsResponse(
    val subtabs: List<StreamSubtab>,
    val total: Int,
    val success: Boolean
)

@Serializable
data class ContentCollectionsResponse(
    val collections: List<ContentCollection>,
    val total: Int,
    val success: Boolean
)

// Error types
@Serializable
data class StreamApiError(
    val error: String,
    val message: String,
    val details: Map<String, String>? = null
)

// Filter types
@Serializable
data class StreamFilters(
    val categoryId: String? = null,
    val contentType: StreamContentType? = null,
    val priceMin: Double? = null,
    val priceMax: Double? = null,
    val isFeatured: Boolean? = null,
    val availabilityStatus: StreamAvailabilityStatus? = null,
    val durationMin: Int? = null,
    val durationMax: Int? = null
)

@Serializable
data class StreamSortOptions(
    val field: SortField,
    val order: SortOrder
)

@Serializable
enum class SortField {
    TITLE,
    PRICE,
    CREATED_AT,
    FEATURED_ORDER
}

@Serializable
enum class SortOrder {
    ASC,
    DESC
}

// Request models for API communication
@Serializable
data class StreamContentRequest(
    val categoryId: String,
    val subtabId: String? = null,
    val page: Int = 1,
    val limit: Int = 20,
    val sortBy: String = "createdAt",
    val sortOrder: String = "desc"
)

@Serializable
data class AddToCartRequest(
    val productId: String,
    val quantity: Int = 1,
    val mediaLicense: String? = null,
    val downloadFormat: String? = null
)

@Serializable
data class PurchaseContentRequest(
    val contentId: String,
    val mediaLicense: String = "personal",
    val downloadFormat: String = "standard"
)

@Serializable
data class PurchaseResponse(
    val success: Boolean,
    val orderId: String?,
    val downloadUrls: List<String>? = null,
    val message: String
)