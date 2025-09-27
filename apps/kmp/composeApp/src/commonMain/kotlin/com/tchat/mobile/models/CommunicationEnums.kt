package com.tchat.mobile.models

/**
 * Communication-related enums (Message, Event, Notification) matching backend shared models
 * Based on backend/shared/models/message.go, event.go, and notification.go
 */

/**
 * Message type enum
 * Matches MessageType from backend
 */
enum class MessageType(val value: String, val displayName: String) {
    TEXT("text", "Text Message"),
    IMAGE("image", "Image"),
    VIDEO("video", "Video"),
    AUDIO("audio", "Audio"),
    FILE("file", "File"),
    LOCATION("location", "Location"),
    STICKER("sticker", "Sticker"),
    SYSTEM("system", "System Message");

    companion object {
        fun fromValue(value: String): MessageType? {
            return values().find { it.value == value }
        }

        fun getMediaTypes(): List<MessageType> {
            return listOf(IMAGE, VIDEO, AUDIO, FILE)
        }

        fun getInteractiveTypes(): List<MessageType> {
            return listOf(TEXT, IMAGE, VIDEO, AUDIO, FILE, LOCATION, STICKER)
        }
    }

    fun isMediaType(): Boolean {
        return this in listOf(IMAGE, VIDEO, AUDIO, FILE)
    }
}

/**
 * Event type enum - comprehensive system events
 * Matches EventType from backend
 */
enum class EventType(val value: String, val category: String, val displayName: String) {
    // User Events
    USER_REGISTERED("user.registered", "user", "User Registered"),
    USER_PROFILE_UPDATED("user.profile_updated", "user", "Profile Updated"),
    USER_KYC_VERIFIED("user.kyc_verified", "user", "KYC Verified"),
    USER_SESSION_CREATED("user.session_created", "user", "Session Created"),
    USER_SESSION_EXPIRED("user.session_expired", "user", "Session Expired"),
    USER_PRESENCE_CHANGED("user.presence_changed", "user", "Presence Changed"),

    // Message Events
    MESSAGE_SENT("message.sent", "message", "Message Sent"),
    MESSAGE_DELIVERED("message.delivered", "message", "Message Delivered"),
    MESSAGE_READ("message.read", "message", "Message Read"),

    // Dialog Events
    DIALOG_CREATED("dialog.created", "dialog", "Dialog Created"),
    DIALOG_PARTICIPANT_ADDED("dialog.participant_added", "dialog", "Participant Added"),

    // Payment Events
    PAYMENT_INITIATED("payment.initiated", "payment", "Payment Initiated"),
    PAYMENT_COMPLETED("payment.completed", "payment", "Payment Completed"),
    PAYMENT_FAILED("payment.failed", "payment", "Payment Failed"),
    WALLET_BALANCE_CHANGED("wallet.balance_changed", "payment", "Wallet Balance Changed"),
    TRANSACTION_CREATED("transaction.created", "payment", "Transaction Created"),

    // Order Events
    ORDER_CREATED("order.created", "order", "Order Created"),
    ORDER_UPDATED("order.updated", "order", "Order Updated"),
    ORDER_FULFILLED("order.fulfilled", "order", "Order Fulfilled"),
    ORDER_CANCELLED("order.cancelled", "order", "Order Cancelled"),

    // Product Events
    PRODUCT_CREATED("product.created", "product", "Product Created"),
    PRODUCT_UPDATED("product.updated", "product", "Product Updated"),

    // Shop Events
    SHOP_CREATED("shop.created", "shop", "Shop Created"),
    SHOP_STATUS_CHANGED("shop.status_changed", "shop", "Shop Status Changed"),

    // System Events
    SYSTEM_STARTUP("system.startup", "system", "System Startup"),
    SYSTEM_SHUTDOWN("system.shutdown", "system", "System Shutdown"),
    SYSTEM_HEALTH_CHECK("system.health_check", "system", "Health Check"),
    SYSTEM_BACKUP_CREATED("system.backup_created", "system", "Backup Created"),
    SYSTEM_MIGRATION_STARTED("system.migration_started", "system", "Migration Started"),
    SYSTEM_MIGRATION_COMPLETED("system.migration_completed", "system", "Migration Completed"),

    // Security Events
    SECURITY_LOGIN_ATTEMPT("security.login_attempt", "security", "Login Attempt"),
    SECURITY_LOGIN_SUCCESS("security.login_success", "security", "Login Success"),
    SECURITY_LOGIN_FAILED("security.login_failed", "security", "Login Failed"),
    SECURITY_PASSWORD_CHANGED("security.password_changed", "security", "Password Changed"),
    SECURITY_SUSPICIOUS_ACTIVITY("security.suspicious_activity", "security", "Suspicious Activity"),
    SECURITY_ACCOUNT_LOCKED("security.account_locked", "security", "Account Locked"),

    // Notification Events
    NOTIFICATION_SENT("notification.sent", "notification", "Notification Sent"),
    NOTIFICATION_DELIVERED("notification.delivered", "notification", "Notification Delivered"),
    NOTIFICATION_FAILED("notification.failed", "notification", "Notification Failed");

    companion object {
        fun fromValue(value: String): EventType? {
            return values().find { it.value == value }
        }

        fun getByCategory(category: String): List<EventType> {
            return values().filter { it.category == category }
        }

        fun getCategories(): List<String> {
            return values().map { it.category }.distinct().sorted()
        }
    }
}

/**
 * Notification type enum
 * Matches NotificationType from backend
 */
enum class NotificationType(val value: String, val displayName: String) {
    SYSTEM("system", "System Notification"),
    MESSAGE("message", "New Message"),
    ORDER("order", "Order Update"),
    PAYMENT("payment", "Payment Update"),
    PROMOTION("promotion", "Promotion"),
    FRIEND("friend", "Friend Activity"),
    REVIEW("review", "New Review"),
    SECURITY("security", "Security Alert"),
    MARKETING("marketing", "Marketing"),
    REMINDER("reminder", "Reminder");

    companion object {
        fun fromValue(value: String): NotificationType? {
            return values().find { it.value == value }
        }

        fun getCriticalTypes(): List<NotificationType> {
            return listOf(SECURITY, PAYMENT, ORDER)
        }

        fun getMarketingTypes(): List<NotificationType> {
            return listOf(PROMOTION, MARKETING)
        }
    }
}

/**
 * Notification status enum
 * Matches NotificationStatus from backend
 */
enum class NotificationStatus(val value: String, val displayName: String) {
    PENDING("pending", "Pending"),
    SENT("sent", "Sent"),
    DELIVERED("delivered", "Delivered"),
    READ("read", "Read"),
    FAILED("failed", "Failed"),
    CANCELLED("cancelled", "Cancelled");

    companion object {
        fun fromValue(value: String): NotificationStatus? {
            return values().find { it.value == value }
        }
    }

    fun isTerminal(): Boolean {
        return this in listOf(READ, FAILED, CANCELLED)
    }
}

/**
 * Notification priority enum
 * Matches NotificationPriority from backend
 */
enum class NotificationPriority(val value: String, val displayName: String, val level: Int) {
    LOW("low", "Low", 1),
    NORMAL("normal", "Normal", 2),
    HIGH("high", "High", 3),
    CRITICAL("critical", "Critical", 4);

    companion object {
        fun fromValue(value: String): NotificationPriority? {
            return values().find { it.value == value }
        }

        fun fromLevel(level: Int): NotificationPriority? {
            return values().find { it.level == level }
        }
    }
}

/**
 * Delivery status enum for messages
 * Matches DeliveryStatus from backend message model
 */
enum class DeliveryStatus(val value: String, val displayName: String) {
    PENDING("pending", "Pending"),
    SENT("sent", "Sent"),
    DELIVERED("delivered", "Delivered"),
    READ("read", "Read"),
    FAILED("failed", "Failed");

    companion object {
        fun fromValue(value: String): DeliveryStatus? {
            return values().find { it.value == value }
        }
    }

    fun isDelivered(): Boolean {
        return this in listOf(DELIVERED, READ)
    }
}