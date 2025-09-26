package com.tchat.mobile.data.models

import kotlinx.datetime.Instant
import kotlinx.serialization.Serializable
import com.benasher44.uuid.uuid4

/**
 * Product represents an item available for purchase
 */
@Serializable
data class Product(
    val id: String = uuid4().toString(),
    val name: String = "",
    val description: String? = null,
    val shortDescription: String? = null,
    val sku: String = "",
    val category: String = "",
    val brand: String? = null,
    val price: Long = 0L, // in cents
    val originalPrice: Long? = null, // in cents, for discount display
    val currency: String = "THB",
    val images: List<String> = emptyList(),
    val thumbnail: String? = null,
    val availability: ProductAvailability = ProductAvailability.IN_STOCK,
    val stock: Int = 0,
    val minOrderQuantity: Int = 1,
    val maxOrderQuantity: Int = 100,
    val weight: Double? = null, // in grams
    val dimensions: ProductDimensions? = null,
    val tags: List<String> = emptyList(),
    val attributes: Map<String, String> = emptyMap(),
    val rating: Double = 0.0,
    val reviewCount: Int = 0,
    val isDigital: Boolean = false,
    val shippingRequired: Boolean = true,
    val taxable: Boolean = true,
    val status: ProductStatus = ProductStatus.ACTIVE,
    val sellerId: String = "",
    val storeName: String = "",
    val createdAt: Instant = kotlinx.datetime.Clock.System.now(),
    val updatedAt: Instant = kotlinx.datetime.Clock.System.now()
) {
    /**
     * Get formatted price with currency symbol
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
     * Convert to public product for API responses
     */
    fun toPublicProduct(): Map<String, Any?> = mapOf(
        "id" to id,
        "name" to name,
        "description" to description,
        "price" to getFormattedPrice(),
        "price_cents" to price,
        "currency" to currency,
        "images" to images,
        "thumbnail" to thumbnail,
        "availability" to availability.name.lowercase(),
        "stock" to stock,
        "rating" to rating,
        "review_count" to reviewCount,
        "has_discount" to hasDiscount(),
        "discount_percentage" to getDiscountPercentage(),
        "store_name" to storeName
    )
}

/**
 * Product Availability status
 */
@Serializable
enum class ProductAvailability {
    IN_STOCK,
    OUT_OF_STOCK,
    LIMITED_STOCK,
    DISCONTINUED,
    PRE_ORDER
}

/**
 * Product Status
 */
@Serializable
enum class ProductStatus {
    ACTIVE,
    INACTIVE,
    DRAFT,
    ARCHIVED
}

/**
 * Product Dimensions
 */
@Serializable
data class ProductDimensions(
    val length: Double, // in cm
    val width: Double,  // in cm
    val height: Double  // in cm
) {
    /**
     * Calculate volume in cubic centimeters
     */
    fun getVolume(): Double = length * width * height
}

/**
 * Order represents a purchase order
 */
@Serializable
data class Order(
    val id: String = uuid4().toString(),
    val orderNumber: String = "",
    val customerId: String = "",
    val customerEmail: String? = null,
    val customerPhone: String? = null,
    val items: List<OrderItem> = emptyList(),
    val subtotal: Long = 0L, // in cents
    val taxAmount: Long = 0L, // in cents
    val shippingAmount: Long = 0L, // in cents
    val discountAmount: Long = 0L, // in cents
    val totalAmount: Long = 0L, // in cents
    val currency: String = "THB",
    val status: OrderStatus = OrderStatus.PENDING,
    val paymentStatus: PaymentStatus = PaymentStatus.PENDING,
    val shippingAddress: Address? = null,
    val billingAddress: Address? = null,
    val shippingMethod: String? = null,
    val trackingNumber: String? = null,
    val notes: String? = null,
    val tags: List<String> = emptyList(),
    val metadata: Map<String, String> = emptyMap(),
    val createdAt: Instant = kotlinx.datetime.Clock.System.now(),
    val updatedAt: Instant = kotlinx.datetime.Clock.System.now(),
    val shippedAt: Instant? = null,
    val deliveredAt: Instant? = null,
    val cancelledAt: Instant? = null
) {
    /**
     * Get total item count
     */
    fun getTotalItemCount(): Int = items.sumOf { it.quantity }

    /**
     * Get formatted total amount with currency symbol
     */
    fun getFormattedTotalAmount(): String {
        val amount = totalAmount / 100.0
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
     * Check if order can be cancelled
     */
    fun canBeCancelled(): Boolean {
        return status in listOf(OrderStatus.PENDING, OrderStatus.CONFIRMED) &&
               paymentStatus != PaymentStatus.PAID
    }

    /**
     * Check if order can be shipped
     */
    fun canBeShipped(): Boolean {
        return status == OrderStatus.CONFIRMED &&
               paymentStatus == PaymentStatus.PAID &&
               shippedAt == null
    }

    /**
     * Cancel order
     */
    fun cancel(): Order = copy(
        status = OrderStatus.CANCELLED,
        cancelledAt = kotlinx.datetime.Clock.System.now(),
        updatedAt = kotlinx.datetime.Clock.System.now()
    )

    /**
     * Mark as shipped
     */
    fun ship(trackingNumber: String? = null): Order = copy(
        status = OrderStatus.SHIPPED,
        trackingNumber = trackingNumber,
        shippedAt = kotlinx.datetime.Clock.System.now(),
        updatedAt = kotlinx.datetime.Clock.System.now()
    )

    /**
     * Mark as delivered
     */
    fun deliver(): Order = copy(
        status = OrderStatus.DELIVERED,
        deliveredAt = kotlinx.datetime.Clock.System.now(),
        updatedAt = kotlinx.datetime.Clock.System.now()
    )

    /**
     * Convert to public order for API responses
     */
    fun toPublicOrder(): Map<String, Any?> = mapOf(
        "id" to id,
        "order_number" to orderNumber,
        "items" to items.map { it.toPublicOrderItem() },
        "total_amount" to getFormattedTotalAmount(),
        "total_amount_cents" to totalAmount,
        "currency" to currency,
        "status" to status.name.lowercase(),
        "payment_status" to paymentStatus.name.lowercase(),
        "item_count" to getTotalItemCount(),
        "shipping_address" to shippingAddress?.toPublicAddress(),
        "tracking_number" to trackingNumber,
        "created_at" to createdAt.toString(),
        "updated_at" to updatedAt.toString()
    )
}

/**
 * Order Status
 */
@Serializable
enum class OrderStatus {
    PENDING,
    CONFIRMED,
    PROCESSING,
    SHIPPED,
    DELIVERED,
    CANCELLED,
    REFUNDED,
    RETURNED
}

/**
 * Payment Status
 */
@Serializable
enum class PaymentStatus {
    PENDING,
    PAID,
    FAILED,
    REFUNDED,
    PARTIALLY_REFUNDED,
    CANCELLED
}

/**
 * Order Item represents an item in an order
 */
@Serializable
data class OrderItem(
    val id: String = uuid4().toString(),
    val productId: String = "",
    val productName: String = "",
    val productSku: String = "",
    val productImage: String? = null,
    val quantity: Int = 1,
    val unitPrice: Long = 0L, // in cents
    val totalPrice: Long = 0L, // in cents
    val currency: String = "THB",
    val attributes: Map<String, String> = emptyMap(), // color, size, etc.
    val notes: String? = null
) {
    /**
     * Get formatted unit price
     */
    fun getFormattedUnitPrice(): String {
        val amount = unitPrice / 100.0
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
     * Get formatted total price
     */
    fun getFormattedTotalPrice(): String {
        val amount = totalPrice / 100.0
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
     * Convert to public order item for API responses
     */
    fun toPublicOrderItem(): Map<String, Any?> = mapOf(
        "id" to id,
        "product_id" to productId,
        "product_name" to productName,
        "product_image" to productImage,
        "quantity" to quantity,
        "unit_price" to getFormattedUnitPrice(),
        "total_price" to getFormattedTotalPrice(),
        "attributes" to attributes
    )
}

/**
 * Address for shipping and billing
 */
@Serializable
data class Address(
    val id: String = uuid4().toString(),
    val firstName: String = "",
    val lastName: String = "",
    val company: String? = null,
    val address1: String = "",
    val address2: String? = null,
    val city: String = "",
    val state: String? = null,
    val postalCode: String = "",
    val country: String = "",
    val phone: String? = null,
    val isDefault: Boolean = false
) {
    /**
     * Get full name
     */
    fun getFullName(): String = "$firstName $lastName".trim()

    /**
     * Get formatted address
     */
    fun getFormattedAddress(): String {
        val parts = mutableListOf<String>()

        if (company?.isNotBlank() == true) parts.add(company!!)
        parts.add(getFullName())
        parts.add(address1)
        if (address2?.isNotBlank() == true) parts.add(address2!!)

        val cityStatePostal = listOfNotNull(
            city.takeIf { it.isNotBlank() },
            state?.takeIf { it.isNotBlank() },
            postalCode.takeIf { it.isNotBlank() }
        ).joinToString(" ")

        if (cityStatePostal.isNotBlank()) parts.add(cityStatePostal)
        parts.add(country)

        return parts.joinToString("\n")
    }

    /**
     * Validate address completeness
     */
    fun isComplete(): Boolean {
        return firstName.isNotBlank() &&
               lastName.isNotBlank() &&
               address1.isNotBlank() &&
               city.isNotBlank() &&
               postalCode.isNotBlank() &&
               country.isNotBlank()
    }

    /**
     * Convert to public address for API responses
     */
    fun toPublicAddress(): Map<String, Any?> = mapOf(
        "id" to id,
        "full_name" to getFullName(),
        "company" to company,
        "address1" to address1,
        "address2" to address2,
        "city" to city,
        "state" to state,
        "postal_code" to postalCode,
        "country" to country,
        "phone" to phone,
        "formatted_address" to getFormattedAddress()
    )
}

/**
 * Shopping Cart for managing items before checkout
 */
@Serializable
data class ShoppingCart(
    val id: String = uuid4().toString(),
    val userId: String = "",
    val items: List<CartItem> = emptyList(),
    val currency: String = "THB",
    val updatedAt: Instant = kotlinx.datetime.Clock.System.now()
) {
    /**
     * Get total item count
     */
    fun getTotalItemCount(): Int = items.sumOf { it.quantity }

    /**
     * Get total amount
     */
    fun getTotalAmount(): Long = items.sumOf { it.getTotalPrice() }

    /**
     * Get formatted total amount
     */
    fun getFormattedTotalAmount(): String {
        val amount = getTotalAmount() / 100.0
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
     * Add item to cart
     */
    fun addItem(product: Product, quantity: Int = 1): ShoppingCart {
        val existingItemIndex = items.indexOfFirst { it.productId == product.id }

        val updatedItems = if (existingItemIndex != -1) {
            // Update existing item quantity
            val existingItem = items[existingItemIndex]
            val newQuantity = existingItem.quantity + quantity
            items.toMutableList().apply {
                set(existingItemIndex, existingItem.copy(quantity = newQuantity))
            }
        } else {
            // Add new item
            items + CartItem(
                productId = product.id,
                productName = product.name,
                productSku = product.sku,
                productImage = product.getMainImageUrl(),
                quantity = quantity,
                unitPrice = product.price,
                currency = product.currency
            )
        }

        return copy(
            items = updatedItems,
            updatedAt = kotlinx.datetime.Clock.System.now()
        )
    }

    /**
     * Remove item from cart
     */
    fun removeItem(productId: String): ShoppingCart {
        return copy(
            items = items.filter { it.productId != productId },
            updatedAt = kotlinx.datetime.Clock.System.now()
        )
    }

    /**
     * Update item quantity
     */
    fun updateItemQuantity(productId: String, quantity: Int): ShoppingCart {
        val updatedItems = items.map { item ->
            if (item.productId == productId) {
                item.copy(quantity = quantity)
            } else {
                item
            }
        }

        return copy(
            items = updatedItems,
            updatedAt = kotlinx.datetime.Clock.System.now()
        )
    }

    /**
     * Clear cart
     */
    fun clear(): ShoppingCart = copy(
        items = emptyList(),
        updatedAt = kotlinx.datetime.Clock.System.now()
    )

    /**
     * Check if cart is empty
     */
    fun isEmpty(): Boolean = items.isEmpty()
}

/**
 * Cart Item represents an item in shopping cart
 */
@Serializable
data class CartItem(
    val productId: String = "",
    val productName: String = "",
    val productSku: String = "",
    val productImage: String? = null,
    val quantity: Int = 1,
    val unitPrice: Long = 0L, // in cents
    val currency: String = "THB",
    val attributes: Map<String, String> = emptyMap()
) {
    /**
     * Get total price for this item
     */
    fun getTotalPrice(): Long = unitPrice * quantity

    /**
     * Get formatted unit price
     */
    fun getFormattedUnitPrice(): String {
        val amount = unitPrice / 100.0
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
     * Get formatted total price
     */
    fun getFormattedTotalPrice(): String {
        val amount = getTotalPrice() / 100.0
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
}