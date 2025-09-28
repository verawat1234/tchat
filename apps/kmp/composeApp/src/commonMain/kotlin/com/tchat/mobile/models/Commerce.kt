package com.tchat.mobile.models

import kotlinx.datetime.Instant
import kotlinx.serialization.Serializable
import kotlinx.datetime.Clock

/**
 * Product represents an item available for purchase
 * Enhanced with Southeast Asian e-commerce features
 */
@Serializable
data class Product(
    val id: String,
    val name: String,
    val description: String? = null,
    val shortDescription: String? = null,
    val sku: String,
    val category: String,
    val brand: String? = null,
    val price: Long, // in cents
    val originalPrice: Long? = null, // in cents, for discount display
    val currency: String = "THB",
    val images: List<String> = emptyList(),
    val thumbnail: String? = null,
    val availability: ProductAvailability = ProductAvailability.IN_STOCK,
    val stock: Int = 0,
    val minOrderQuantity: Int = 1,
    val maxOrderQuantity: Int = 100,
    val weight: Double? = null, // in grams
    val rating: Double = 0.0,
    val reviewCount: Int = 0,
    val isDigital: Boolean = false,
    val shippingRequired: Boolean = true,
    val taxable: Boolean = true,
    val status: ProductStatus = ProductStatus.ACTIVE,
    val sellerId: String,
    val storeName: String,
    val createdAt: Instant = Clock.System.now(),
    val updatedAt: Instant = Clock.System.now(),
    // UI-specific fields for compatibility with existing screens
    val deliveryTime: String = "30 min",
    val distance: String = "2.5 km",
    val isHot: Boolean = false,
    val orders: Int = 0,
    val tags: List<String> = emptyList()
) {
    /**
     * Get formatted price with Southeast Asian currency symbols
     */
    fun getFormattedPrice(): String {
        val amount = price / 100.0
        return when (currency) {
            "THB" -> "฿${String.format("%.2f", amount)}"
            "SGD" -> "S$${String.format("%.2f", amount)}"
            "IDR" -> "Rp${String.format("%.0f", amount)}"
            "MYR" -> "RM${String.format("%.2f", amount)}"
            "PHP" -> "₱${String.format("%.2f", amount)}"
            "VND" -> "${String.format("%.0f", amount)}₫"
            else -> "$${String.format("%.2f", amount)}"
        }
    }

    /**
     * Check if product has discount
     */
    fun hasDiscount(): Boolean = originalPrice != null && originalPrice!! > price

    /**
     * Get discount percentage
     */
    fun getDiscountPercentage(): Int {
        return if (hasDiscount()) {
            ((originalPrice!! - price) * 100 / originalPrice!!).toInt()
        } else 0
    }

    /**
     * Check if product is available for purchase
     */
    fun isAvailableForPurchase(quantity: Int = 1): Boolean {
        return status == ProductStatus.ACTIVE &&
               availability == ProductAvailability.IN_STOCK &&
               stock >= quantity &&
               quantity >= minOrderQuantity &&
               quantity <= maxOrderQuantity
    }

    /**
     * Get main image URL
     */
    fun getMainImageUrl(): String? = thumbnail ?: images.firstOrNull()

    /**
     * Convert to ProductItem for UI compatibility
     */
    fun toProductItem(): ProductItem {
        return ProductItem(
            id = id,
            name = name,
            price = price / 100.0, // Convert cents to dollars
            originalPrice = originalPrice?.let { it / 100.0 },
            rating = rating,
            category = category,
            merchant = storeName,
            image = getMainImageUrl() ?: "",
            deliveryTime = deliveryTime,
            distance = distance,
            isHot = isHot,
            discount = if (hasDiscount()) getDiscountPercentage() else null,
            orders = orders
        )
    }
}

/**
 * Store represents a merchant/shop
 */
@Serializable
data class Store(
    val id: String,
    val name: String,
    val description: String,
    val avatar: String,
    val coverImage: String,
    val rating: Double = 0.0,
    val deliveryTime: String = "30 min",
    val distance: String = "2.5 km",
    val isVerified: Boolean = false,
    val followers: Int = 0,
    val totalProducts: Int = 0,
    val createdAt: Instant = Clock.System.now(),
    val updatedAt: Instant = Clock.System.now()
) {
    /**
     * Convert to ShopItem for UI compatibility
     */
    fun toShopItem(): ShopItem {
        return ShopItem(
            id = id,
            name = name,
            description = description,
            avatar = avatar,
            coverImage = coverImage,
            rating = rating,
            deliveryTime = deliveryTime,
            distance = distance,
            isVerified = isVerified,
            followers = followers,
            totalProducts = totalProducts
        )
    }
}

/**
 * Product Availability status
 */
@Serializable
enum class ProductAvailability(val displayName: String) {
    IN_STOCK("In Stock"),
    OUT_OF_STOCK("Out of Stock"),
    LIMITED_STOCK("Limited Stock"),
    DISCONTINUED("Discontinued"),
    PRE_ORDER("Pre-order")
}


// Legacy UI models for compatibility (from StoreScreen.kt)
data class ShopItem(
    val id: String,
    val name: String,
    val description: String,
    val avatar: String,
    val coverImage: String,
    val rating: Double,
    val deliveryTime: String,
    val distance: String,
    val isVerified: Boolean = false,
    val followers: Int = 0,
    val totalProducts: Int = 0
)

data class ProductItem(
    val id: String,
    val name: String,
    val price: Double,
    val originalPrice: Double? = null,
    val rating: Double,
    val category: String,
    val merchant: String,
    val image: String,
    val deliveryTime: String = "30 min",
    val distance: String = "2.5 km",
    val isHot: Boolean = false,
    val discount: Int? = null,
    val orders: Int = 0
)

data class LiveStreamItem(
    val id: String,
    val title: String,
    val merchant: String,
    val viewers: Int,
    val thumbnail: String,
    val products: List<ProductItem>,
    val isLive: Boolean = true,
    val duration: String = "00:45:32"
)