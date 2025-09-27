package com.tchat.mobile.models

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFalse
import kotlin.test.assertNull
import kotlin.test.assertTrue

class OrderEnumsTest {

    @Test
    fun testOrderStatusFromValue() {
        assertEquals(OrderStatus.PENDING, OrderStatus.fromValue("pending"))
        assertEquals(OrderStatus.CONFIRMED, OrderStatus.fromValue("confirmed"))
        assertEquals(OrderStatus.PROCESSING, OrderStatus.fromValue("processing"))
        assertEquals(OrderStatus.SHIPPED, OrderStatus.fromValue("shipped"))
        assertEquals(OrderStatus.DELIVERED, OrderStatus.fromValue("delivered"))
        assertEquals(OrderStatus.CANCELLED, OrderStatus.fromValue("cancelled"))
        assertEquals(OrderStatus.REFUNDED, OrderStatus.fromValue("refunded"))
        assertEquals(OrderStatus.RETURNED, OrderStatus.fromValue("returned"))
        assertNull(OrderStatus.fromValue("invalid"))
    }

    @Test
    fun testOrderStatusValues() {
        assertEquals("pending", OrderStatus.PENDING.value)
        assertEquals("confirmed", OrderStatus.CONFIRMED.value)
        assertEquals("processing", OrderStatus.PROCESSING.value)
        assertEquals("shipped", OrderStatus.SHIPPED.value)
        assertEquals("delivered", OrderStatus.DELIVERED.value)
        assertEquals("cancelled", OrderStatus.CANCELLED.value)
        assertEquals("refunded", OrderStatus.REFUNDED.value)
        assertEquals("returned", OrderStatus.RETURNED.value)
    }

    @Test
    fun testOrderStatusDisplayNames() {
        assertEquals("Pending", OrderStatus.PENDING.displayName)
        assertEquals("Confirmed", OrderStatus.CONFIRMED.displayName)
        assertEquals("Processing", OrderStatus.PROCESSING.displayName)
        assertEquals("Shipped", OrderStatus.SHIPPED.displayName)
        assertEquals("Delivered", OrderStatus.DELIVERED.displayName)
        assertEquals("Cancelled", OrderStatus.CANCELLED.displayName)
        assertEquals("Refunded", OrderStatus.REFUNDED.displayName)
        assertEquals("Returned", OrderStatus.RETURNED.displayName)
    }

    @Test
    fun testOrderStatusIsTerminal() {
        assertFalse(OrderStatus.PENDING.isTerminal())
        assertFalse(OrderStatus.CONFIRMED.isTerminal())
        assertFalse(OrderStatus.PROCESSING.isTerminal())
        assertFalse(OrderStatus.SHIPPED.isTerminal())
        assertTrue(OrderStatus.DELIVERED.isTerminal())
        assertTrue(OrderStatus.CANCELLED.isTerminal())
        assertTrue(OrderStatus.REFUNDED.isTerminal())
        assertTrue(OrderStatus.RETURNED.isTerminal())
    }

    @Test
    fun testOrderStatusCanCancel() {
        assertTrue(OrderStatus.PENDING.canCancel())
        assertTrue(OrderStatus.CONFIRMED.canCancel())
        assertFalse(OrderStatus.PROCESSING.canCancel())
        assertFalse(OrderStatus.SHIPPED.canCancel())
        assertFalse(OrderStatus.DELIVERED.canCancel())
        assertFalse(OrderStatus.CANCELLED.canCancel())
        assertFalse(OrderStatus.REFUNDED.canCancel())
        assertFalse(OrderStatus.RETURNED.canCancel())
    }

    @Test
    fun testPaymentStatusFromValue() {
        assertEquals(PaymentStatus.PENDING, PaymentStatus.fromValue("pending"))
        assertEquals(PaymentStatus.AUTHORIZED, PaymentStatus.fromValue("authorized"))
        assertEquals(PaymentStatus.PAID, PaymentStatus.fromValue("paid"))
        assertEquals(PaymentStatus.FAILED, PaymentStatus.fromValue("failed"))
        assertEquals(PaymentStatus.CANCELLED, PaymentStatus.fromValue("cancelled"))
        assertEquals(PaymentStatus.REFUNDED, PaymentStatus.fromValue("refunded"))
        assertEquals(PaymentStatus.PARTIAL_REFUND, PaymentStatus.fromValue("partial_refund"))
        assertNull(PaymentStatus.fromValue("invalid"))
    }

    @Test
    fun testPaymentStatusValues() {
        assertEquals("pending", PaymentStatus.PENDING.value)
        assertEquals("authorized", PaymentStatus.AUTHORIZED.value)
        assertEquals("paid", PaymentStatus.PAID.value)
        assertEquals("failed", PaymentStatus.FAILED.value)
        assertEquals("cancelled", PaymentStatus.CANCELLED.value)
        assertEquals("refunded", PaymentStatus.REFUNDED.value)
        assertEquals("partial_refund", PaymentStatus.PARTIAL_REFUND.value)
    }

    @Test
    fun testPaymentStatusDisplayNames() {
        assertEquals("Pending", PaymentStatus.PENDING.displayName)
        assertEquals("Authorized", PaymentStatus.AUTHORIZED.displayName)
        assertEquals("Paid", PaymentStatus.PAID.displayName)
        assertEquals("Failed", PaymentStatus.FAILED.displayName)
        assertEquals("Cancelled", PaymentStatus.CANCELLED.displayName)
        assertEquals("Refunded", PaymentStatus.REFUNDED.displayName)
        assertEquals("Partially Refunded", PaymentStatus.PARTIAL_REFUND.displayName)
    }

    @Test
    fun testFulfillmentStatusFromValue() {
        assertEquals(FulfillmentStatus.UNFULFILLED, FulfillmentStatus.fromValue("unfulfilled"))
        assertEquals(FulfillmentStatus.PARTIALLY_FULFILLED, FulfillmentStatus.fromValue("partially_fulfilled"))
        assertEquals(FulfillmentStatus.FULFILLED, FulfillmentStatus.fromValue("fulfilled"))
        assertEquals(FulfillmentStatus.SHIPPED, FulfillmentStatus.fromValue("shipped"))
        assertEquals(FulfillmentStatus.DELIVERED, FulfillmentStatus.fromValue("delivered"))
        assertEquals(FulfillmentStatus.RETURNED, FulfillmentStatus.fromValue("returned"))
        assertNull(FulfillmentStatus.fromValue("invalid"))
    }

    @Test
    fun testFulfillmentStatusValues() {
        assertEquals("unfulfilled", FulfillmentStatus.UNFULFILLED.value)
        assertEquals("partially_fulfilled", FulfillmentStatus.PARTIALLY_FULFILLED.value)
        assertEquals("fulfilled", FulfillmentStatus.FULFILLED.value)
        assertEquals("shipped", FulfillmentStatus.SHIPPED.value)
        assertEquals("delivered", FulfillmentStatus.DELIVERED.value)
        assertEquals("returned", FulfillmentStatus.RETURNED.value)
    }

    @Test
    fun testFulfillmentStatusDisplayNames() {
        assertEquals("Unfulfilled", FulfillmentStatus.UNFULFILLED.displayName)
        assertEquals("Partially Fulfilled", FulfillmentStatus.PARTIALLY_FULFILLED.displayName)
        assertEquals("Fulfilled", FulfillmentStatus.FULFILLED.displayName)
        assertEquals("Shipped", FulfillmentStatus.SHIPPED.displayName)
        assertEquals("Delivered", FulfillmentStatus.DELIVERED.displayName)
        assertEquals("Returned", FulfillmentStatus.RETURNED.displayName)
    }
}