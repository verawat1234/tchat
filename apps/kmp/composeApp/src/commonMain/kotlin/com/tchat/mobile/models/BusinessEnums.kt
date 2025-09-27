package com.tchat.mobile.models

/**
 * Business-related enums matching backend shared models
 * Based on backend/shared/models/business.go
 */

/**
 * Business verification status enum
 * Matches BusinessVerificationStatus from backend
 */
enum class BusinessVerificationStatus(val value: String) {
    PENDING("pending"),
    VERIFIED("verified"),
    REJECTED("rejected"),
    SUSPENDED("suspended");

    companion object {
        fun fromValue(value: String): BusinessVerificationStatus? {
            return values().find { it.value == value }
        }
    }
}

/**
 * Business category enum
 * Matches valid business categories from backend IsValidCategory()
 */
enum class BusinessCategory(val value: String, val displayName: String) {
    ELECTRONICS("electronics", "Electronics"),
    FASHION("fashion", "Fashion"),
    FOOD("food", "Food & Beverages"),
    HEALTH("health", "Health & Wellness"),
    BEAUTY("beauty", "Beauty & Personal Care"),
    HOME("home", "Home & Garden"),
    SPORTS("sports", "Sports & Recreation"),
    AUTOMOTIVE("automotive", "Automotive"),
    BOOKS("books", "Books & Media"),
    TOYS("toys", "Toys & Games"),
    SERVICES("services", "Services"),
    DIGITAL("digital", "Digital Products"),
    AGRICULTURE("agriculture", "Agriculture"),
    CRAFTS("crafts", "Arts & Crafts"),
    JEWELRY("jewelry", "Jewelry & Accessories"),
    TRAVEL("travel", "Travel & Tourism"),
    EDUCATION("education", "Education"),
    FINANCE("finance", "Finance"),
    REAL_ESTATE("real_estate", "Real Estate"),
    ENTERTAINMENT("entertainment", "Entertainment");

    companion object {
        fun fromValue(value: String): BusinessCategory? {
            return values().find { it.value == value }
        }

        fun getAllCategories(): List<BusinessCategory> {
            return values().toList()
        }
    }
}