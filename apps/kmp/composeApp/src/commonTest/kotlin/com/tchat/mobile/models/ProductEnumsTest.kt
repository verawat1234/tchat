package com.tchat.mobile.models

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFalse
import kotlin.test.assertNull
import kotlin.test.assertTrue

class ProductEnumsTest {

    @Test
    fun testProductStatusFromValue() {
        assertEquals(ProductStatus.DRAFT, ProductStatus.fromValue("draft"))
        assertEquals(ProductStatus.ACTIVE, ProductStatus.fromValue("active"))
        assertEquals(ProductStatus.INACTIVE, ProductStatus.fromValue("inactive"))
        assertEquals(ProductStatus.ARCHIVED, ProductStatus.fromValue("archived"))
        assertEquals(ProductStatus.DELETED, ProductStatus.fromValue("deleted"))
        assertNull(ProductStatus.fromValue("invalid"))
    }

    @Test
    fun testProductStatusValues() {
        assertEquals("draft", ProductStatus.DRAFT.value)
        assertEquals("active", ProductStatus.ACTIVE.value)
        assertEquals("inactive", ProductStatus.INACTIVE.value)
        assertEquals("archived", ProductStatus.ARCHIVED.value)
        assertEquals("deleted", ProductStatus.DELETED.value)
    }

    @Test
    fun testProductStatusDisplayNames() {
        assertEquals("Draft", ProductStatus.DRAFT.displayName)
        assertEquals("Active", ProductStatus.ACTIVE.displayName)
        assertEquals("Inactive", ProductStatus.INACTIVE.displayName)
        assertEquals("Archived", ProductStatus.ARCHIVED.displayName)
        assertEquals("Deleted", ProductStatus.DELETED.displayName)
    }

    @Test
    fun testProductStatusIsAvailable() {
        assertTrue(ProductStatus.ACTIVE.isAvailable())
        assertFalse(ProductStatus.DRAFT.isAvailable())
        assertFalse(ProductStatus.INACTIVE.isAvailable())
        assertFalse(ProductStatus.ARCHIVED.isAvailable())
        assertFalse(ProductStatus.DELETED.isAvailable())
    }

    @Test
    fun testProductStatusStaticIsAvailable() {
        assertTrue(ProductStatus.isAvailable(ProductStatus.ACTIVE))
        assertFalse(ProductStatus.isAvailable(ProductStatus.DRAFT))
        assertFalse(ProductStatus.isAvailable(ProductStatus.INACTIVE))
        assertFalse(ProductStatus.isAvailable(ProductStatus.ARCHIVED))
        assertFalse(ProductStatus.isAvailable(ProductStatus.DELETED))
    }

    @Test
    fun testProductTypeFromValue() {
        assertEquals(ProductType.PHYSICAL, ProductType.fromValue("physical"))
        assertEquals(ProductType.DIGITAL, ProductType.fromValue("digital"))
        assertEquals(ProductType.SERVICE, ProductType.fromValue("service"))
        assertNull(ProductType.fromValue("invalid"))
    }

    @Test
    fun testProductTypeValues() {
        assertEquals("physical", ProductType.PHYSICAL.value)
        assertEquals("digital", ProductType.DIGITAL.value)
        assertEquals("service", ProductType.SERVICE.value)
    }

    @Test
    fun testProductTypeDisplayNames() {
        assertEquals("Physical Product", ProductType.PHYSICAL.displayName)
        assertEquals("Digital Product", ProductType.DIGITAL.displayName)
        assertEquals("Service", ProductType.SERVICE.displayName)
    }

    @Test
    fun testProductConditionFromValue() {
        assertEquals(ProductCondition.NEW, ProductCondition.fromValue("new"))
        assertEquals(ProductCondition.USED, ProductCondition.fromValue("used"))
        assertEquals(ProductCondition.REFURBISHED, ProductCondition.fromValue("refurbished"))
        assertNull(ProductCondition.fromValue("invalid"))
    }

    @Test
    fun testProductConditionValues() {
        assertEquals("new", ProductCondition.NEW.value)
        assertEquals("used", ProductCondition.USED.value)
        assertEquals("refurbished", ProductCondition.REFURBISHED.value)
    }

    @Test
    fun testProductConditionDisplayNames() {
        assertEquals("New", ProductCondition.NEW.displayName)
        assertEquals("Used", ProductCondition.USED.displayName)
        assertEquals("Refurbished", ProductCondition.REFURBISHED.displayName)
    }
}