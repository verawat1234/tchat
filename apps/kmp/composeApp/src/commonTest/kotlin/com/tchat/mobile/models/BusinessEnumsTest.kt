package com.tchat.mobile.models

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNull
import kotlin.test.assertTrue

class BusinessEnumsTest {

    @Test
    fun testBusinessVerificationStatusFromValue() {
        assertEquals(BusinessVerificationStatus.PENDING, BusinessVerificationStatus.fromValue("pending"))
        assertEquals(BusinessVerificationStatus.VERIFIED, BusinessVerificationStatus.fromValue("verified"))
        assertEquals(BusinessVerificationStatus.REJECTED, BusinessVerificationStatus.fromValue("rejected"))
        assertEquals(BusinessVerificationStatus.SUSPENDED, BusinessVerificationStatus.fromValue("suspended"))
        assertNull(BusinessVerificationStatus.fromValue("invalid"))
    }

    @Test
    fun testBusinessVerificationStatusValues() {
        assertEquals("pending", BusinessVerificationStatus.PENDING.value)
        assertEquals("verified", BusinessVerificationStatus.VERIFIED.value)
        assertEquals("rejected", BusinessVerificationStatus.REJECTED.value)
        assertEquals("suspended", BusinessVerificationStatus.SUSPENDED.value)
    }

    @Test
    fun testBusinessCategoryFromValue() {
        assertEquals(BusinessCategory.ELECTRONICS, BusinessCategory.fromValue("electronics"))
        assertEquals(BusinessCategory.FASHION, BusinessCategory.fromValue("fashion"))
        assertEquals(BusinessCategory.FOOD, BusinessCategory.fromValue("food"))
        assertEquals(BusinessCategory.HEALTH, BusinessCategory.fromValue("health"))
        assertEquals(BusinessCategory.BEAUTY, BusinessCategory.fromValue("beauty"))
        assertEquals(BusinessCategory.HOME, BusinessCategory.fromValue("home"))
        assertEquals(BusinessCategory.SPORTS, BusinessCategory.fromValue("sports"))
        assertEquals(BusinessCategory.AUTOMOTIVE, BusinessCategory.fromValue("automotive"))
        assertEquals(BusinessCategory.BOOKS, BusinessCategory.fromValue("books"))
        assertEquals(BusinessCategory.TOYS, BusinessCategory.fromValue("toys"))
        assertEquals(BusinessCategory.SERVICES, BusinessCategory.fromValue("services"))
        assertEquals(BusinessCategory.DIGITAL, BusinessCategory.fromValue("digital"))
        assertEquals(BusinessCategory.AGRICULTURE, BusinessCategory.fromValue("agriculture"))
        assertEquals(BusinessCategory.CRAFTS, BusinessCategory.fromValue("crafts"))
        assertEquals(BusinessCategory.JEWELRY, BusinessCategory.fromValue("jewelry"))
        assertEquals(BusinessCategory.TRAVEL, BusinessCategory.fromValue("travel"))
        assertEquals(BusinessCategory.EDUCATION, BusinessCategory.fromValue("education"))
        assertEquals(BusinessCategory.FINANCE, BusinessCategory.fromValue("finance"))
        assertEquals(BusinessCategory.REAL_ESTATE, BusinessCategory.fromValue("real_estate"))
        assertEquals(BusinessCategory.ENTERTAINMENT, BusinessCategory.fromValue("entertainment"))
        assertNull(BusinessCategory.fromValue("invalid"))
    }

    @Test
    fun testBusinessCategoryDisplayNames() {
        assertEquals("Electronics", BusinessCategory.ELECTRONICS.displayName)
        assertEquals("Fashion", BusinessCategory.FASHION.displayName)
        assertEquals("Food & Beverages", BusinessCategory.FOOD.displayName)
        assertEquals("Health & Wellness", BusinessCategory.HEALTH.displayName)
        assertEquals("Beauty & Personal Care", BusinessCategory.BEAUTY.displayName)
        assertEquals("Home & Garden", BusinessCategory.HOME.displayName)
        assertEquals("Sports & Recreation", BusinessCategory.SPORTS.displayName)
        assertEquals("Automotive", BusinessCategory.AUTOMOTIVE.displayName)
        assertEquals("Books & Media", BusinessCategory.BOOKS.displayName)
        assertEquals("Toys & Games", BusinessCategory.TOYS.displayName)
        assertEquals("Services", BusinessCategory.SERVICES.displayName)
        assertEquals("Digital Products", BusinessCategory.DIGITAL.displayName)
        assertEquals("Agriculture", BusinessCategory.AGRICULTURE.displayName)
        assertEquals("Arts & Crafts", BusinessCategory.CRAFTS.displayName)
        assertEquals("Jewelry & Accessories", BusinessCategory.JEWELRY.displayName)
        assertEquals("Travel & Tourism", BusinessCategory.TRAVEL.displayName)
        assertEquals("Education", BusinessCategory.EDUCATION.displayName)
        assertEquals("Finance", BusinessCategory.FINANCE.displayName)
        assertEquals("Real Estate", BusinessCategory.REAL_ESTATE.displayName)
        assertEquals("Entertainment", BusinessCategory.ENTERTAINMENT.displayName)
    }

    @Test
    fun testBusinessCategoryGetAllCategories() {
        val allCategories = BusinessCategory.getAllCategories()
        assertEquals(20, allCategories.size)
        assertTrue(allCategories.contains(BusinessCategory.ELECTRONICS))
        assertTrue(allCategories.contains(BusinessCategory.ENTERTAINMENT))
    }

    @Test
    fun testBusinessCategoryValues() {
        assertEquals("electronics", BusinessCategory.ELECTRONICS.value)
        assertEquals("fashion", BusinessCategory.FASHION.value)
        assertEquals("real_estate", BusinessCategory.REAL_ESTATE.value)
    }
}