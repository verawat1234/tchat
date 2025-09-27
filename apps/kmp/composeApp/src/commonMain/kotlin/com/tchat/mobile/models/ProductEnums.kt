package com.tchat.mobile.models

/**
 * Product-related enums matching backend shared models
 * Based on backend/shared/models/product.go
 */

/**
 * Product status enum
 * Matches ProductStatus from backend
 */
enum class ProductStatus(val value: String, val displayName: String) {
    DRAFT("draft", "Draft"),
    ACTIVE("active", "Active"),
    INACTIVE("inactive", "Inactive"),
    ARCHIVED("archived", "Archived"),
    DELETED("deleted", "Deleted");

    companion object {
        fun fromValue(value: String): ProductStatus? {
            return values().find { it.value == value }
        }

        /**
         * Check if the product is available for purchase
         */
        fun isAvailable(status: ProductStatus): Boolean {
            return status == ACTIVE
        }
    }

    /**
     * Check if this status is available for purchase
     */
    fun isAvailable(): Boolean {
        return this == ACTIVE
    }
}

/**
 * Product type enum
 * Matches ProductType from backend
 */
enum class ProductType(val value: String, val displayName: String) {
    PHYSICAL("physical", "Physical Product"),
    DIGITAL("digital", "Digital Product"),
    SERVICE("service", "Service");

    companion object {
        fun fromValue(value: String): ProductType? {
            return values().find { it.value == value }
        }
    }
}

/**
 * Product condition enum
 * Matches ProductCondition from backend
 */
enum class ProductCondition(val value: String, val displayName: String) {
    NEW("new", "New"),
    USED("used", "Used"),
    REFURBISHED("refurbished", "Refurbished");

    companion object {
        fun fromValue(value: String): ProductCondition? {
            return values().find { it.value == value }
        }
    }
}