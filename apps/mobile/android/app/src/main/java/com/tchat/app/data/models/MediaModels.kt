// Media Store Android data classes
// Generated for Media Store Tabs feature implementation

package com.tchat.app.data.models

import kotlinx.serialization.Serializable
import java.time.Instant
import java.time.LocalDateTime
import java.time.ZoneOffset
import java.time.format.DateTimeFormatter

// MARK: - Media Category
@Serializable
data class MediaCategory(
    val id: String,
    val name: String,
    val displayOrder: Int,
    val iconName: String,
    val isActive: Boolean,
    val featuredContentEnabled: Boolean,
    val createdAt: String,
    val updatedAt: String
) {
    val createdAtDate: LocalDateTime
        get() = LocalDateTime.parse(createdAt, DateTimeFormatter.ISO_DATE_TIME)

    val updatedAtDate: LocalDateTime
        get() = LocalDateTime.parse(updatedAt, DateTimeFormatter.ISO_DATE_TIME)
}

// MARK: - Media Subtab
@Serializable
data class MediaSubtab(
    val id: String,
    val categoryId: String,
    val name: String,
    val displayOrder: Int,
    val filterCriteria: Map<String, kotlinx.serialization.json.JsonElement>,
    val isActive: Boolean,
    val createdAt: String,
    val updatedAt: String
) {
    val createdAtDate: LocalDateTime
        get() = LocalDateTime.parse(createdAt, DateTimeFormatter.ISO_DATE_TIME)

    val updatedAtDate: LocalDateTime
        get() = LocalDateTime.parse(updatedAt, DateTimeFormatter.ISO_DATE_TIME)
}

// MARK: - Content Type
@Serializable
enum class MediaContentType(val value: String) {
    BOOK("book"),
    PODCAST("podcast"),
    VIDEO("video"),
    CARTOON("cartoon");

    companion object {
        fun fromString(value: String): MediaContentType? {
            return values().find { it.value == value }
        }
    }
}

// MARK: - Availability Status
@Serializable
enum class MediaAvailabilityStatus(val value: String) {
    AVAILABLE("available"),
    COMING_SOON("coming_soon"),
    UNAVAILABLE("unavailable");

    companion object {
        fun fromString(value: String): MediaAvailabilityStatus? {
            return values().find { it.value == value }
        }
    }
}

// MARK: - Media Content Item
@Serializable
data class MediaContentItem(
    val id: String,
    val categoryId: String,
    val title: String,
    val description: String,
    val thumbnailUrl: String,
    val contentUrl: String? = null,
    val contentType: MediaContentType,
    val duration: Int? = null,
    val price: Double,
    val currency: String,
    val availabilityStatus: MediaAvailabilityStatus,
    val isFeatured: Boolean,
    val featuredOrder: Int? = null,
    val metadata: Map<String, kotlinx.serialization.json.JsonElement>,
    val createdAt: String,
    val updatedAt: String
) {
    val createdAtDate: LocalDateTime
        get() = LocalDateTime.parse(createdAt, DateTimeFormatter.ISO_DATE_TIME)

    val updatedAtDate: LocalDateTime
        get() = LocalDateTime.parse(updatedAt, DateTimeFormatter.ISO_DATE_TIME)

    val durationString: String?
        get() = duration?.let {
            val minutes = it / 60
            val seconds = it % 60
            String.format("%d:%02d", minutes, seconds)
        }
}

// MARK: - Product Type
@Serializable
enum class MediaProductType(val value: String) {
    PHYSICAL("physical"),
    MEDIA("media");

    companion object {
        fun fromString(value: String): MediaProductType? {
            return values().find { it.value == value }
        }
    }
}

// MARK: - Media License
@Serializable
enum class MediaLicense(val value: String) {
    PERSONAL("personal"),
    COMMERCIAL("commercial"),
    EDUCATIONAL("educational");

    companion object {
        fun fromString(value: String): MediaLicense? {
            return values().find { it.value == value }
        }
    }
}

// MARK: - Download Format
@Serializable
enum class MediaDownloadFormat(val value: String) {
    PDF("PDF"),
    EPUB("EPUB"),
    MP3("MP3"),
    MP4("MP4"),
    FLAC("FLAC");

    companion object {
        fun fromString(value: String): MediaDownloadFormat? {
            return values().find { it.value == value }
        }
    }
}

// MARK: - Media Product
@Serializable
data class MediaProduct(
    val id: String,
    val name: String,
    val description: String,
    val price: Double,
    val currency: String,
    val productType: MediaProductType,
    val mediaContentId: String? = null,
    val mediaMetadata: MediaMetadata? = null,
    val category: String,
    val isActive: Boolean,
    val stockQuantity: Int? = null,
    val createdAt: String,
    val updatedAt: String
) {
    @Serializable
    data class MediaMetadata(
        val contentType: MediaContentType,
        val duration: Int? = null,
        val format: String? = null,
        val license: String? = null
    )

    val createdAtDate: LocalDateTime
        get() = LocalDateTime.parse(createdAt, DateTimeFormatter.ISO_DATE_TIME)

    val updatedAtDate: LocalDateTime
        get() = LocalDateTime.parse(updatedAt, DateTimeFormatter.ISO_DATE_TIME)
}

// MARK: - Media Cart Item
@Serializable
data class MediaCartItem(
    val id: String,
    val cartId: String,
    val productId: String,
    val mediaContentId: String? = null,
    val quantity: Int,
    val unitPrice: Double,
    val totalPrice: Double,
    val mediaLicense: MediaLicense? = null,
    val downloadFormat: MediaDownloadFormat? = null,
    val createdAt: String,
    val updatedAt: String
) {
    val createdAtDate: LocalDateTime
        get() = LocalDateTime.parse(createdAt, DateTimeFormatter.ISO_DATE_TIME)

    val updatedAtDate: LocalDateTime
        get() = LocalDateTime.parse(updatedAt, DateTimeFormatter.ISO_DATE_TIME)
}

// MARK: - Order Status
@Serializable
enum class MediaOrderStatus(val value: String) {
    PENDING("pending"),
    PROCESSING("processing"),
    COMPLETED("completed"),
    CANCELLED("cancelled");

    companion object {
        fun fromString(value: String): MediaOrderStatus? {
            return values().find { it.value == value }
        }
    }
}

// MARK: - Delivery Status
@Serializable
enum class MediaDeliveryStatus(val value: String) {
    PENDING("pending"),
    DELIVERED("delivered"),
    FAILED("failed");

    companion object {
        fun fromString(value: String): MediaDeliveryStatus? {
            return values().find { it.value == value }
        }
    }
}

// MARK: - Media Order
@Serializable
data class MediaOrder(
    val id: String,
    val userId: String,
    val status: MediaOrderStatus,
    val totalPhysicalAmount: Double,
    val totalMediaAmount: Double,
    val totalAmount: Double,
    val currency: String,
    val mediaDeliveryStatus: MediaDeliveryStatus,
    val shippingAddress: String? = null,
    val items: List<MediaOrderItem>,
    val createdAt: String,
    val updatedAt: String
) {
    val createdAtDate: LocalDateTime
        get() = LocalDateTime.parse(createdAt, DateTimeFormatter.ISO_DATE_TIME)

    val updatedAtDate: LocalDateTime
        get() = LocalDateTime.parse(updatedAt, DateTimeFormatter.ISO_DATE_TIME)
}

// MARK: - Media Order Item
@Serializable
data class MediaOrderItem(
    val id: String,
    val orderId: String,
    val productId: String,
    val mediaContentId: String? = null,
    val quantity: Int,
    val unitPrice: Double,
    val totalPrice: Double,
    val mediaLicense: MediaLicense? = null,
    val downloadFormat: MediaDownloadFormat? = null,
    val deliveryStatus: MediaDeliveryStatus? = null,
    val downloadUrl: String? = null,
    val createdAt: String,
    val updatedAt: String
) {
    val createdAtDate: LocalDateTime
        get() = LocalDateTime.parse(createdAt, DateTimeFormatter.ISO_DATE_TIME)

    val updatedAtDate: LocalDateTime
        get() = LocalDateTime.parse(updatedAt, DateTimeFormatter.ISO_DATE_TIME)
}

// MARK: - API Response Types
@Serializable
data class MediaCategoriesResponse(
    val categories: List<MediaCategory>,
    val total: Int
)

@Serializable
data class MediaContentResponse(
    val items: List<MediaContentItem>,
    val page: Int,
    val limit: Int,
    val total: Int,
    val hasMore: Boolean
)

@Serializable
data class MediaFeaturedResponse(
    val items: List<MediaContentItem>,
    val total: Int,
    val hasMore: Boolean
)

@Serializable
data class MediaSubtabsResponse(
    val subtabs: List<MediaSubtab>,
    val defaultSubtab: String
)

@Serializable
data class MediaSearchResponse(
    val items: List<MediaContentItem>,
    val query: String,
    val total: Int,
    val page: Int
)

// MARK: - Store Integration Types
@Serializable
data class AddMediaToCartRequest(
    val mediaContentId: String,
    val quantity: Int,
    val mediaLicense: MediaLicense,
    val downloadFormat: MediaDownloadFormat
)

@Serializable
data class AddMediaToCartResponse(
    val cartId: String,
    val itemsCount: Int,
    val totalAmount: Double,
    val currency: String,
    val addedItem: MediaCartItem
)

@Serializable
data class UnifiedCartResponse(
    val cartId: String,
    val physicalItems: List<MediaCartItem>,
    val mediaItems: List<MediaCartItem>,
    val totalPhysicalAmount: Double,
    val totalMediaAmount: Double,
    val totalAmount: Double,
    val currency: String,
    val itemsCount: Int
)

@Serializable
data class MediaCheckoutValidationRequest(
    val cartId: String,
    val mediaItems: List<MediaCartItem>
)

@Serializable
data class MediaCheckoutValidationResponse(
    val isValid: Boolean,
    val validItems: List<MediaCartItem>,
    val invalidItems: List<MediaCartItem>,
    val totalMediaAmount: Double,
    val estimatedDeliveryTime: String
)

@Serializable
data class MediaOrdersResponse(
    val orders: List<MediaOrder>,
    val pagination: PaginationInfo
) {
    @Serializable
    data class PaginationInfo(
        val page: Int,
        val limit: Int,
        val total: Int,
        val hasMore: Boolean
    )
}

// MARK: - Error Types
@Serializable
data class MediaApiError(
    val error: String,
    val message: String,
    val details: Map<String, kotlinx.serialization.json.JsonElement>? = null
)