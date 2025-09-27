package com.tchat.mobile.models

/**
 * Order-related enums matching backend shared models
 * Based on backend/shared/models/order.go
 */

/**
 * Order status enum
 * Matches OrderStatus from backend
 */
enum class OrderStatus(val value: String, val displayName: String) {
    PENDING("pending", "Pending"),
    CONFIRMED("confirmed", "Confirmed"),
    PROCESSING("processing", "Processing"),
    SHIPPED("shipped", "Shipped"),
    DELIVERED("delivered", "Delivered"),
    CANCELLED("cancelled", "Cancelled"),
    REFUNDED("refunded", "Refunded"),
    RETURNED("returned", "Returned");

    companion object {
        fun fromValue(value: String): OrderStatus? {
            return values().find { it.value == value }
        }
    }

    /**
     * Check if the order status is terminal (no further changes possible)
     */
    fun isTerminal(): Boolean {
        return this in listOf(DELIVERED, CANCELLED, REFUNDED, RETURNED)
    }

    /**
     * Check if the order can be cancelled
     */
    fun canCancel(): Boolean {
        return this == PENDING || this == CONFIRMED
    }
}

/**
 * Payment status enum
 * Matches PaymentStatus from backend
 */
enum class PaymentStatus(val value: String, val displayName: String) {
    PENDING("pending", "Pending"),
    AUTHORIZED("authorized", "Authorized"),
    PAID("paid", "Paid"),
    FAILED("failed", "Failed"),
    CANCELLED("cancelled", "Cancelled"),
    REFUNDED("refunded", "Refunded"),
    PARTIAL_REFUND("partial_refund", "Partially Refunded");

    companion object {
        fun fromValue(value: String): PaymentStatus? {
            return values().find { it.value == value }
        }
    }
}

/**
 * Fulfillment status enum
 * Matches FulfillmentStatus from backend
 */
enum class FulfillmentStatus(val value: String, val displayName: String) {
    UNFULFILLED("unfulfilled", "Unfulfilled"),
    PARTIALLY_FULFILLED("partially_fulfilled", "Partially Fulfilled"),
    FULFILLED("fulfilled", "Fulfilled"),
    SHIPPED("shipped", "Shipped"),
    DELIVERED("delivered", "Delivered"),
    RETURNED("returned", "Returned");

    companion object {
        fun fromValue(value: String): FulfillmentStatus? {
            return values().find { it.value == value }
        }
    }
}