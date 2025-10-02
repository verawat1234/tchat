package com.tchat.mobile.models

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFalse
import kotlin.test.assertNull
import kotlin.test.assertTrue

class CommunicationEnumsTest {

    @Test
    fun testMessageTypeFromValue() {
        assertEquals(MessageType.TEXT, MessageType.fromValue("text"))
        assertEquals(MessageType.IMAGE, MessageType.fromValue("image"))
        assertEquals(MessageType.VIDEO, MessageType.fromValue("video"))
        assertEquals(MessageType.AUDIO, MessageType.fromValue("audio"))
        assertEquals(MessageType.FILE, MessageType.fromValue("file"))
        assertEquals(MessageType.LOCATION, MessageType.fromValue("location"))
        assertEquals(MessageType.STICKER, MessageType.fromValue("sticker"))
        assertEquals(MessageType.SYSTEM, MessageType.fromValue("system"))
        assertNull(MessageType.fromValue("invalid"))
    }

    @Test
    fun testMessageTypeValues() {
        assertEquals("text", MessageType.TEXT.value)
        assertEquals("image", MessageType.IMAGE.value)
        assertEquals("video", MessageType.VIDEO.value)
        assertEquals("audio", MessageType.AUDIO.value)
        assertEquals("file", MessageType.FILE.value)
        assertEquals("location", MessageType.LOCATION.value)
        assertEquals("sticker", MessageType.STICKER.value)
        assertEquals("system", MessageType.SYSTEM.value)
    }

    @Test
    fun testMessageTypeDisplayNames() {
        assertEquals("Text Message", MessageType.TEXT.displayName)
        assertEquals("Image", MessageType.IMAGE.displayName)
        assertEquals("Video", MessageType.VIDEO.displayName)
        assertEquals("Audio", MessageType.AUDIO.displayName)
        assertEquals("File", MessageType.FILE.displayName)
        assertEquals("Location", MessageType.LOCATION.displayName)
        assertEquals("Sticker", MessageType.STICKER.displayName)
        assertEquals("System Message", MessageType.SYSTEM.displayName)
    }

    @Test
    fun testMessageTypeIsMediaType() {
        assertFalse(MessageType.TEXT.isMediaType())
        assertTrue(MessageType.IMAGE.isMediaType())
        assertTrue(MessageType.VIDEO.isMediaType())
        assertTrue(MessageType.AUDIO.isMediaType())
        assertTrue(MessageType.FILE.isMediaType())
        assertFalse(MessageType.LOCATION.isMediaType())
        assertFalse(MessageType.STICKER.isMediaType())
        assertFalse(MessageType.SYSTEM.isMediaType())
    }

    @Test
    fun testMessageTypeGetMediaTypes() {
        val mediaTypes = MessageType.getMediaTypes()
        assertEquals(4, mediaTypes.size)
        assertTrue(mediaTypes.contains(MessageType.IMAGE))
        assertTrue(mediaTypes.contains(MessageType.VIDEO))
        assertTrue(mediaTypes.contains(MessageType.AUDIO))
        assertTrue(mediaTypes.contains(MessageType.FILE))
    }

    @Test
    fun testMessageTypeGetInteractiveTypes() {
        val interactiveTypes = MessageType.getInteractiveTypes()
        assertEquals(7, interactiveTypes.size)
        assertTrue(interactiveTypes.contains(MessageType.TEXT))
        assertTrue(interactiveTypes.contains(MessageType.IMAGE))
        assertTrue(interactiveTypes.contains(MessageType.VIDEO))
        assertTrue(interactiveTypes.contains(MessageType.AUDIO))
        assertTrue(interactiveTypes.contains(MessageType.FILE))
        assertTrue(interactiveTypes.contains(MessageType.LOCATION))
        assertTrue(interactiveTypes.contains(MessageType.STICKER))
        assertFalse(interactiveTypes.contains(MessageType.SYSTEM))
    }

    @Test
    fun testEventTypeFromValue() {
        assertEquals(EventType.USER_REGISTERED, EventType.fromValue("user.registered"))
        assertEquals(EventType.MESSAGE_SENT, EventType.fromValue("message.sent"))
        assertEquals(EventType.ORDER_CREATED, EventType.fromValue("order.created"))
        assertEquals(EventType.PAYMENT_COMPLETED, EventType.fromValue("payment.completed"))
        assertEquals(EventType.SECURITY_LOGIN_SUCCESS, EventType.fromValue("security.login_success"))
        assertEquals(EventType.NOTIFICATION_SENT, EventType.fromValue("notification.sent"))
        assertNull(EventType.fromValue("invalid.event"))
    }

    @Test
    fun testEventTypeValues() {
        assertEquals("user.registered", EventType.USER_REGISTERED.value)
        assertEquals("message.sent", EventType.MESSAGE_SENT.value)
        assertEquals("order.created", EventType.ORDER_CREATED.value)
        assertEquals("payment.completed", EventType.PAYMENT_COMPLETED.value)
        assertEquals("security.login_success", EventType.SECURITY_LOGIN_SUCCESS.value)
        assertEquals("notification.sent", EventType.NOTIFICATION_SENT.value)
    }

    @Test
    fun testEventTypeCategories() {
        assertEquals("user", EventType.USER_REGISTERED.category)
        assertEquals("message", EventType.MESSAGE_SENT.category)
        assertEquals("order", EventType.ORDER_CREATED.category)
        assertEquals("payment", EventType.PAYMENT_COMPLETED.category)
        assertEquals("security", EventType.SECURITY_LOGIN_SUCCESS.category)
        assertEquals("notification", EventType.NOTIFICATION_SENT.category)
    }

    @Test
    fun testEventTypeDisplayNames() {
        assertEquals("User Registered", EventType.USER_REGISTERED.displayName)
        assertEquals("Message Sent", EventType.MESSAGE_SENT.displayName)
        assertEquals("Order Created", EventType.ORDER_CREATED.displayName)
        assertEquals("Payment Completed", EventType.PAYMENT_COMPLETED.displayName)
        assertEquals("Login Success", EventType.SECURITY_LOGIN_SUCCESS.displayName)
        assertEquals("Notification Sent", EventType.NOTIFICATION_SENT.displayName)
    }

    @Test
    fun testEventTypeGetByCategory() {
        val userEvents = EventType.getByCategory("user")
        assertTrue(userEvents.contains(EventType.USER_REGISTERED))
        assertTrue(userEvents.contains(EventType.USER_PROFILE_UPDATED))
        assertTrue(userEvents.contains(EventType.USER_KYC_VERIFIED))
        assertTrue(userEvents.contains(EventType.USER_SESSION_CREATED))
        assertTrue(userEvents.contains(EventType.USER_SESSION_EXPIRED))
        assertTrue(userEvents.contains(EventType.USER_PRESENCE_CHANGED))

        val messageEvents = EventType.getByCategory("message")
        assertTrue(messageEvents.contains(EventType.MESSAGE_SENT))
        assertTrue(messageEvents.contains(EventType.MESSAGE_DELIVERED))
        assertTrue(messageEvents.contains(EventType.MESSAGE_READ))

        val paymentEvents = EventType.getByCategory("payment")
        assertTrue(paymentEvents.contains(EventType.PAYMENT_INITIATED))
        assertTrue(paymentEvents.contains(EventType.PAYMENT_COMPLETED))
        assertTrue(paymentEvents.contains(EventType.PAYMENT_FAILED))
        assertTrue(paymentEvents.contains(EventType.WALLET_BALANCE_CHANGED))
        assertTrue(paymentEvents.contains(EventType.TRANSACTION_CREATED))
    }

    @Test
    fun testEventTypeGetCategories() {
        val categories = EventType.getCategories()
        assertTrue(categories.contains("user"))
        assertTrue(categories.contains("message"))
        assertTrue(categories.contains("dialog"))
        assertTrue(categories.contains("payment"))
        assertTrue(categories.contains("order"))
        assertTrue(categories.contains("product"))
        assertTrue(categories.contains("shop"))
        assertTrue(categories.contains("system"))
        assertTrue(categories.contains("security"))
        assertTrue(categories.contains("notification"))
        assertEquals(categories, categories.sorted()) // Should be sorted
    }

    @Test
    fun testNotificationTypeFromValue() {
        assertEquals(NotificationType.SYSTEM, NotificationType.fromValue("system"))
        assertEquals(NotificationType.MESSAGE, NotificationType.fromValue("message"))
        assertEquals(NotificationType.ORDER, NotificationType.fromValue("order"))
        assertEquals(NotificationType.PAYMENT, NotificationType.fromValue("payment"))
        assertEquals(NotificationType.PROMOTION, NotificationType.fromValue("promotion"))
        assertEquals(NotificationType.FRIEND, NotificationType.fromValue("friend"))
        assertEquals(NotificationType.REVIEW, NotificationType.fromValue("review"))
        assertEquals(NotificationType.SECURITY, NotificationType.fromValue("security"))
        assertEquals(NotificationType.MARKETING, NotificationType.fromValue("marketing"))
        assertEquals(NotificationType.REMINDER, NotificationType.fromValue("reminder"))
        assertNull(NotificationType.fromValue("invalid"))
    }

    @Test
    fun testNotificationTypeGetCriticalTypes() {
        val criticalTypes = NotificationType.getCriticalTypes()
        assertEquals(3, criticalTypes.size)
        assertTrue(criticalTypes.contains(NotificationType.SECURITY))
        assertTrue(criticalTypes.contains(NotificationType.PAYMENT))
        assertTrue(criticalTypes.contains(NotificationType.ORDER))
    }

    @Test
    fun testNotificationTypeGetMarketingTypes() {
        val marketingTypes = NotificationType.getMarketingTypes()
        assertEquals(2, marketingTypes.size)
        assertTrue(marketingTypes.contains(NotificationType.PROMOTION))
        assertTrue(marketingTypes.contains(NotificationType.MARKETING))
    }

    @Test
    fun testNotificationStatusFromValue() {
        assertEquals(NotificationStatus.PENDING, NotificationStatus.fromValue("pending"))
        assertEquals(NotificationStatus.SENT, NotificationStatus.fromValue("sent"))
        assertEquals(NotificationStatus.DELIVERED, NotificationStatus.fromValue("delivered"))
        assertEquals(NotificationStatus.READ, NotificationStatus.fromValue("read"))
        assertEquals(NotificationStatus.FAILED, NotificationStatus.fromValue("failed"))
        assertEquals(NotificationStatus.CANCELLED, NotificationStatus.fromValue("cancelled"))
        assertNull(NotificationStatus.fromValue("invalid"))
    }

    @Test
    fun testNotificationStatusIsTerminal() {
        assertFalse(NotificationStatus.PENDING.isTerminal())
        assertFalse(NotificationStatus.SENT.isTerminal())
        assertFalse(NotificationStatus.DELIVERED.isTerminal())
        assertTrue(NotificationStatus.READ.isTerminal())
        assertTrue(NotificationStatus.FAILED.isTerminal())
        assertTrue(NotificationStatus.CANCELLED.isTerminal())
    }

    @Test
    fun testNotificationPriorityFromValue() {
        assertEquals(NotificationPriority.LOW, NotificationPriority.fromValue("low"))
        assertEquals(NotificationPriority.NORMAL, NotificationPriority.fromValue("normal"))
        assertEquals(NotificationPriority.HIGH, NotificationPriority.fromValue("high"))
        assertEquals(NotificationPriority.CRITICAL, NotificationPriority.fromValue("critical"))
        assertNull(NotificationPriority.fromValue("invalid"))
    }

    @Test
    fun testNotificationPriorityFromLevel() {
        assertEquals(NotificationPriority.LOW, NotificationPriority.fromLevel(1))
        assertEquals(NotificationPriority.NORMAL, NotificationPriority.fromLevel(2))
        assertEquals(NotificationPriority.HIGH, NotificationPriority.fromLevel(3))
        assertEquals(NotificationPriority.CRITICAL, NotificationPriority.fromLevel(4))
        assertNull(NotificationPriority.fromLevel(99))
    }

    @Test
    fun testNotificationPriorityLevels() {
        assertEquals(1, NotificationPriority.LOW.level)
        assertEquals(2, NotificationPriority.NORMAL.level)
        assertEquals(3, NotificationPriority.HIGH.level)
        assertEquals(4, NotificationPriority.CRITICAL.level)
    }

    @Test
    fun testDeliveryStatusFromValue() {
        assertEquals(DeliveryStatus.PENDING, DeliveryStatus.fromValue("pending"))
        assertEquals(DeliveryStatus.SENT, DeliveryStatus.fromValue("sent"))
        assertEquals(DeliveryStatus.DELIVERED, DeliveryStatus.fromValue("delivered"))
        assertEquals(DeliveryStatus.READ, DeliveryStatus.fromValue("read"))
        assertEquals(DeliveryStatus.FAILED, DeliveryStatus.fromValue("failed"))
        assertNull(DeliveryStatus.fromValue("invalid"))
    }

    @Test
    fun testDeliveryStatusIsDelivered() {
        assertFalse(DeliveryStatus.PENDING.isDelivered())
        assertFalse(DeliveryStatus.SENT.isDelivered())
        assertTrue(DeliveryStatus.DELIVERED.isDelivered())
        assertTrue(DeliveryStatus.READ.isDelivered())
        assertFalse(DeliveryStatus.FAILED.isDelivered())
    }
}

// Extension functions for MessageType companion object functionality
object MessageTypeExtensions {
    fun getMediaTypes(): List<MessageType> = listOf(
        MessageType.IMAGE,
        MessageType.VIDEO,
        MessageType.AUDIO,
        MessageType.FILE
    )

    fun getInteractiveTypes(): List<MessageType> = listOf(
        MessageType.TEXT,
        MessageType.IMAGE,
        MessageType.VIDEO,
        MessageType.AUDIO,
        MessageType.LOCATION,
        MessageType.CONTACT,
        MessageType.POLL
    )
}

// Make these available as static-like functions on MessageType
fun MessageType.Companion.getMediaTypes() = MessageTypeExtensions.getMediaTypes()
fun MessageType.Companion.getInteractiveTypes() = MessageTypeExtensions.getInteractiveTypes()