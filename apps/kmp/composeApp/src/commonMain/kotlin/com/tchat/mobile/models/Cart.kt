package com.tchat.mobile.models

import kotlinx.serialization.Serializable
import kotlinx.datetime.*

/**
 * T029: CartItem and ShoppingCart models
 *
 * Shopping cart management with item variants, pricing calculations,
 * and cross-platform synchronization. Supports offline cart persistence.
 */
@Serializable
data class ShoppingCart(
    val id: String,
    val userId: String,
    val sessionId: String? = null,
    val items: List<CartItem>,
    val summary: CartSummary,
    val appliedCoupons: List<AppliedCoupon> = emptyList(),
    val shippingAddress: ShippingAddress? = null,
    val billingAddress: BillingAddress? = null,
    val deliveryMethod: DeliveryMethod? = null,
    val paymentMethod: String? = null,
    val notes: String? = null,
    val currency: String = "USD",
    val expiresAt: String? = null, // For guest carts
    val lastSyncAt: String? = null,
    val createdAt: String,
    val updatedAt: String,
    val metadata: Map<String, String> = emptyMap()
)

@Serializable
data class CartItem(
    val id: String,
    val productId: String,
    val productName: String,
    val productImage: String? = null,
    val productSlug: String? = null,
    val sellerId: String,
    val sellerName: String,
    val variantId: String? = null,
    val variantName: String? = null,
    val quantity: Int,
    val unitPrice: Money,
    val originalUnitPrice: Money? = null, // For discounted items
    val totalPrice: Money,
    val selectedAttributes: Map<String, String> = emptyMap(), // color: "red", size: "L"
    val customization: Map<String, String> = emptyMap(), // Custom text, engraving
    val gift: GiftOptions? = null,
    val availability: CartItemAvailability,
    val shipping: CartItemShipping? = null,
    val addedAt: String,
    val updatedAt: String,
    val isSelected: Boolean = true, // For partial checkout
    val isSaveForLater: Boolean = false,
    val maxQuantity: Int? = null,
    val metadata: Map<String, String> = emptyMap()
)

@Serializable
data class CartSummary(
    val itemCount: Int, // Total number of items
    val totalItems: Int, // Sum of all quantities
    val subtotal: Money, // Sum of all item totals
    val discounts: List<DiscountBreakdown> = emptyList(),
    val totalDiscount: Money = Money(0.0, subtotal.currency),
    val shippingCost: Money = Money(0.0, subtotal.currency),
    val taxAmount: Money = Money(0.0, subtotal.currency),
    val fees: List<FeeBreakdown> = emptyList(),
    val totalFees: Money = Money(0.0, subtotal.currency),
    val total: Money,
    val estimatedTotal: Money? = null, // Including estimated tax/shipping
    val savings: Money? = null, // Total amount saved from original prices
    val freeShippingThreshold: Money? = null,
    val freeShippingProgress: Double = 0.0, // 0.0 to 1.0
    val currency: String = subtotal.currency
)

@Serializable
data class CartItemAvailability(
    val inStock: Boolean,
    val stockQuantity: Int? = null,
    val availableQuantity: Int? = null,
    val backorderAllowed: Boolean = false,
    val estimatedRestockDate: String? = null,
    val priceChanged: Boolean = false,
    val previousPrice: Money? = null,
    val discontinued: Boolean = false,
    val status: CartItemStatus = if (inStock) CartItemStatus.AVAILABLE else CartItemStatus.OUT_OF_STOCK
)

enum class CartItemStatus {
    AVAILABLE,
    OUT_OF_STOCK,
    LIMITED_STOCK,
    BACKORDER,
    PRICE_CHANGED,
    DISCONTINUED,
    UNAVAILABLE
}

@Serializable
data class CartItemShipping(
    val method: String? = null,
    val cost: Money? = null,
    val estimatedDays: Int? = null,
    val isFree: Boolean = false,
    val restrictions: List<String> = emptyList()
)

@Serializable
data class GiftOptions(
    val isGift: Boolean = false,
    val giftMessage: String? = null,
    val giftWrap: String? = null,
    val giftWrapCost: Money? = null,
    val recipientName: String? = null,
    val recipientEmail: String? = null,
    val deliveryDate: String? = null
)

@Serializable
data class AppliedCoupon(
    val id: String,
    val code: String,
    val title: String,
    val description: String? = null,
    val type: CouponType,
    val value: Money,
    val discountAmount: Money,
    val appliedTo: List<String> = emptyList(), // Item IDs
    val minOrderAmount: Money? = null,
    val maxDiscountAmount: Money? = null,
    val expiresAt: String? = null,
    val usageLimit: Int? = null,
    val usedCount: Int = 0,
    val metadata: Map<String, String> = emptyMap()
)

enum class CouponType {
    PERCENTAGE,
    FIXED_AMOUNT,
    FREE_SHIPPING,
    BUY_X_GET_Y,
    BULK_DISCOUNT
}

@Serializable
data class DiscountBreakdown(
    val id: String,
    val name: String,
    val type: DiscountType,
    val amount: Money,
    val appliedTo: List<String> = emptyList(), // Item IDs
    val source: String? = null // Coupon code, loyalty program, etc.
)

enum class DiscountType {
    ITEM_DISCOUNT,
    CART_DISCOUNT,
    COUPON_DISCOUNT,
    LOYALTY_DISCOUNT,
    BULK_DISCOUNT,
    SEASONAL_DISCOUNT,
    MEMBERSHIP_DISCOUNT
}

@Serializable
data class FeeBreakdown(
    val id: String,
    val name: String,
    val type: FeeType,
    val amount: Money,
    val description: String? = null,
    val optional: Boolean = false
)

enum class FeeType {
    PROCESSING_FEE,
    SERVICE_FEE,
    CONVENIENCE_FEE,
    HANDLING_FEE,
    INSURANCE_FEE,
    RUSH_PROCESSING_FEE
}

@Serializable
data class ShippingAddress(
    val id: String? = null,
    val firstName: String,
    val lastName: String,
    val company: String? = null,
    val addressLine1: String,
    val addressLine2: String? = null,
    val city: String,
    val state: String,
    val zipCode: String,
    val country: String,
    val phone: String? = null,
    val isDefault: Boolean = false,
    val instructions: String? = null,
    val addressType: String = "home", // "home", "office", "other"
    val validated: Boolean = false,
    val coordinates: Coordinates? = null
)

@Serializable
data class BillingAddress(
    val id: String? = null,
    val firstName: String,
    val lastName: String,
    val company: String? = null,
    val addressLine1: String,
    val addressLine2: String? = null,
    val city: String,
    val state: String,
    val zipCode: String,
    val country: String,
    val phone: String? = null,
    val email: String? = null,
    val isDefault: Boolean = false,
    val sameAsShipping: Boolean = false
)

@Serializable
data class Coordinates(
    val latitude: Double,
    val longitude: Double
)

@Serializable
data class DeliveryMethod(
    val id: String,
    val name: String,
    val type: DeliveryType,
    val cost: Money,
    val estimatedDays: Int,
    val description: String? = null,
    val carrier: String? = null,
    val trackingAvailable: Boolean = false,
    val insuranceIncluded: Boolean = false,
    val signatureRequired: Boolean = false,
    val cutoffTime: String? = null // "2:00 PM" for same-day delivery
)

enum class DeliveryType {
    STANDARD,
    EXPEDITED,
    OVERNIGHT,
    TWO_DAY,
    SAME_DAY,
    PICKUP,
    LOCAL_DELIVERY,
    DIGITAL_DELIVERY
}

/**
 * Shopping cart utilities and extensions
 */
fun ShoppingCart.isEmpty(): Boolean = items.isEmpty()

fun ShoppingCart.isNotEmpty(): Boolean = items.isNotEmpty()

fun ShoppingCart.getActiveItems(): List<CartItem> = items.filter { it.isSelected && !it.isSaveForLater }

fun ShoppingCart.getSavedItems(): List<CartItem> = items.filter { it.isSaveForLater }

fun ShoppingCart.getTotalWeight(): Double? {
    // Would need product weight data to calculate
    return null
}

fun ShoppingCart.hasUnavailableItems(): Boolean =
    items.any { !it.availability.inStock && !it.availability.backorderAllowed }

fun ShoppingCart.hasGiftItems(): Boolean = items.any { it.gift?.isGift == true }

fun ShoppingCart.needsShipping(): Boolean = items.any { /* check if physical product */ true }

fun ShoppingCart.requiresAge(): Boolean = items.any { /* check if age-restricted */ false }

fun ShoppingCart.getSellerGroups(): Map<String, List<CartItem>> = items.groupBy { it.sellerId }

fun ShoppingCart.canApplyCoupon(coupon: AppliedCoupon): Boolean {
    // Check if coupon is valid and applicable
    val minOrderMet = coupon.minOrderAmount?.let { summary.subtotal.amount >= it.amount } ?: true
    val notExpired = coupon.expiresAt?.let {
        Clock.System.now() < Instant.parse(it)
    } ?: true

    return minOrderMet && notExpired
}

/**
 * Cart item utilities
 */
fun CartItem.isAvailable(): Boolean = availability.inStock || availability.backorderAllowed

fun CartItem.hasDiscount(): Boolean = originalUnitPrice != null && originalUnitPrice.amount > unitPrice.amount

fun CartItem.getDiscountAmount(): Money? {
    return if (hasDiscount() && originalUnitPrice != null) {
        originalUnitPrice.subtract(unitPrice).multiply(quantity.toDouble())
    } else null
}

fun CartItem.getDiscountPercentage(): Int? {
    return if (hasDiscount() && originalUnitPrice != null) {
        val discount = originalUnitPrice.amount - unitPrice.amount
        ((discount.toDouble() / originalUnitPrice.amount) * 100).toInt()
    } else null
}

fun CartItem.canIncreaseQuantity(): Boolean {
    val currentMax = maxQuantity ?: availability.stockQuantity ?: Int.MAX_VALUE
    return quantity < currentMax
}

fun CartItem.canDecreaseQuantity(): Boolean = quantity > 1

fun CartItem.getMaxSelectableQuantity(): Int {
    return minOf(
        maxQuantity ?: Int.MAX_VALUE,
        availability.availableQuantity ?: availability.stockQuantity ?: Int.MAX_VALUE,
        99 // Reasonable UI limit
    )
}

fun CartItem.needsUpdate(): Boolean {
    return availability.priceChanged || availability.discontinued || !availability.inStock
}

fun CartItem.getDisplayName(): String {
    return if (variantName != null) {
        "$productName - $variantName"
    } else {
        productName
    }
}

fun CartItem.getAttributesDisplay(): String? {
    return if (selectedAttributes.isNotEmpty()) {
        selectedAttributes.values.joinToString(", ")
    } else null
}

/**
 * Cart operations
 */
fun ShoppingCart.addItem(item: CartItem): ShoppingCart {
    val existingIndex = items.indexOfFirst {
        it.productId == item.productId &&
        it.variantId == item.variantId &&
        it.selectedAttributes == item.selectedAttributes &&
        it.customization == item.customization &&
        !it.isSaveForLater
    }

    val updatedItems = if (existingIndex >= 0) {
        val existingItem = items[existingIndex]
        val newQuantity = existingItem.quantity + item.quantity
        val updatedItem = existingItem.copy(
            quantity = newQuantity,
            totalPrice = existingItem.unitPrice.multiply(newQuantity.toDouble()),
            updatedAt = Clock.System.now().toString()
        )
        items.toMutableList().apply { set(existingIndex, updatedItem) }
    } else {
        items + item.copy(
            totalPrice = item.unitPrice.multiply(item.quantity.toDouble()),
            updatedAt = Clock.System.now().toString()
        )
    }

    return copy(
        items = updatedItems,
        summary = recalculateSummary(updatedItems),
        updatedAt = Clock.System.now().toString()
    )
}

fun ShoppingCart.removeItem(itemId: String): ShoppingCart {
    val updatedItems = items.filter { it.id != itemId }
    return copy(
        items = updatedItems,
        summary = recalculateSummary(updatedItems),
        updatedAt = Clock.System.now().toString()
    )
}

fun ShoppingCart.updateItemQuantity(itemId: String, quantity: Int): ShoppingCart {
    val updatedItems = items.map { item ->
        if (item.id == itemId) {
            item.copy(
                quantity = quantity,
                totalPrice = item.unitPrice.multiply(quantity.toDouble()),
                updatedAt = Clock.System.now().toString()
            )
        } else item
    }

    return copy(
        items = updatedItems,
        summary = recalculateSummary(updatedItems),
        updatedAt = Clock.System.now().toString()
    )
}

fun ShoppingCart.saveItemForLater(itemId: String): ShoppingCart {
    val updatedItems = items.map { item ->
        if (item.id == itemId) {
            item.copy(
                isSaveForLater = true,
                isSelected = false,
                updatedAt = Clock.System.now().toString()
            )
        } else item
    }

    return copy(
        items = updatedItems,
        summary = recalculateSummary(updatedItems.filter { it.isSelected && !it.isSaveForLater }),
        updatedAt = Clock.System.now().toString()
    )
}

fun ShoppingCart.moveItemToCart(itemId: String): ShoppingCart {
    val updatedItems = items.map { item ->
        if (item.id == itemId) {
            item.copy(
                isSaveForLater = false,
                isSelected = true,
                updatedAt = Clock.System.now().toString()
            )
        } else item
    }

    return copy(
        items = updatedItems,
        summary = recalculateSummary(updatedItems.filter { it.isSelected && !it.isSaveForLater }),
        updatedAt = Clock.System.now().toString()
    )
}

fun ShoppingCart.selectItems(itemIds: List<String>): ShoppingCart {
    val updatedItems = items.map { item ->
        item.copy(isSelected = item.id in itemIds)
    }

    return copy(
        items = updatedItems,
        summary = recalculateSummary(updatedItems.filter { it.isSelected && !it.isSaveForLater }),
        updatedAt = Clock.System.now().toString()
    )
}

fun ShoppingCart.applyCoupon(coupon: AppliedCoupon): ShoppingCart {
    if (!canApplyCoupon(coupon)) return this

    val updatedCoupons = appliedCoupons + coupon
    val activeItems = getActiveItems()

    return copy(
        appliedCoupons = updatedCoupons,
        summary = recalculateSummary(activeItems, updatedCoupons),
        updatedAt = Clock.System.now().toString()
    )
}

fun ShoppingCart.removeCoupon(couponId: String): ShoppingCart {
    val updatedCoupons = appliedCoupons.filter { it.id != couponId }
    val activeItems = getActiveItems()

    return copy(
        appliedCoupons = updatedCoupons,
        summary = recalculateSummary(activeItems, updatedCoupons),
        updatedAt = Clock.System.now().toString()
    )
}

/**
 * Private helper functions for cart calculations
 */
private fun recalculateSummary(
    items: List<CartItem>,
    coupons: List<AppliedCoupon> = emptyList()
): CartSummary {
    val activeItems = items.filter { it.isSelected && !it.isSaveForLater && it.isAvailable() }

    val itemCount = activeItems.size
    val totalItems = activeItems.sumOf { it.quantity }
    val subtotal = activeItems.fold(Money(0.0, "USD")) { acc, item ->
        acc.add(item.totalPrice)
    }

    // Calculate discounts (simplified)
    val totalDiscount = coupons.fold(Money(0.0, subtotal.currency)) { acc, coupon ->
        acc.add(coupon.discountAmount)
    }

    // Calculate total
    val total = subtotal.subtract(totalDiscount)

    return CartSummary(
        itemCount = itemCount,
        totalItems = totalItems,
        subtotal = subtotal,
        totalDiscount = totalDiscount,
        total = total,
        currency = subtotal.currency
    )
}