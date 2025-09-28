package com.tchat.mobile.models

import kotlinx.serialization.Serializable

/**
 * Commerce-specific message content models
 * Used for work chat system between users and shops
 */

@Serializable
data class ProductMessage(
    val productId: String,
    val name: String,
    val description: String,
    val price: Double,
    val currency: String = "THB",
    val imageUrl: String? = null,
    val category: String? = null,
    val shopId: String,
    val shopName: String,
    val isAvailable: Boolean = true,
    val variants: List<ProductVariant> = emptyList(),
    val specifications: Map<String, String> = emptyMap()
)

@Serializable
data class ProductVariant(
    val id: String,
    val name: String,
    val price: Double,
    val isAvailable: Boolean = true,
    val attributes: Map<String, String> = emptyMap() // e.g., "size" -> "M", "color" -> "Red"
)

@Serializable
data class InvoiceMessage(
    val invoiceId: String,
    val invoiceNumber: String,
    val items: List<InvoiceItem>,
    val subtotal: Double,
    val tax: Double,
    val total: Double,
    val currency: String = "THB",
    val dueDate: String,
    val status: InvoiceStatus,
    val shopId: String,
    val shopName: String,
    val customerNote: String? = null,
    val paymentTerms: String? = null
)

@Serializable
data class InvoiceItem(
    val productId: String,
    val name: String,
    val quantity: Int,
    val unitPrice: Double,
    val totalPrice: Double,
    val variantId: String? = null,
    val variantName: String? = null
)

@Serializable
enum class InvoiceStatus(val value: String, val displayName: String) {
    DRAFT("draft", "Draft"),
    SENT("sent", "Sent"),
    VIEWED("viewed", "Viewed"),
    PAID("paid", "Paid"),
    OVERDUE("overdue", "Overdue"),
    CANCELLED("cancelled", "Cancelled");

    companion object {
        fun fromValue(value: String): InvoiceStatus? {
            return values().find { it.value == value }
        }
    }
}

@Serializable
data class OrderMessage(
    val orderId: String,
    val orderNumber: String,
    val items: List<OrderItem>,
    val subtotal: Double,
    val shippingCost: Double,
    val tax: Double,
    val total: Double,
    val currency: String = "THB",
    val status: OrderStatus,
    val shopId: String,
    val shopName: String,
    val deliveryAddress: DeliveryAddress? = null,
    val estimatedDelivery: String? = null,
    val trackingNumber: String? = null,
    val customerNotes: String? = null
)

@Serializable
data class OrderItem(
    val productId: String,
    val name: String,
    val quantity: Int,
    val unitPrice: Double,
    val totalPrice: Double,
    val variantId: String? = null,
    val variantName: String? = null,
    val imageUrl: String? = null
)

@Serializable
data class DeliveryAddress(
    val fullName: String,
    val phoneNumber: String,
    val addressLine1: String,
    val addressLine2: String? = null,
    val city: String,
    val state: String,
    val postalCode: String,
    val country: String = "Thailand"
)

// Using the existing OrderStatus from OrderEnums.kt - adding color extension
fun OrderStatus.getStatusColor(): String {
    return when (this) {
        OrderStatus.PENDING -> "#F59E0B" // amber
        OrderStatus.CONFIRMED -> "#3B82F6" // blue
        OrderStatus.PROCESSING -> "#8B5CF6" // purple
        OrderStatus.SHIPPED -> "#06B6D4" // cyan
        OrderStatus.DELIVERED -> "#10B981" // green
        OrderStatus.CANCELLED -> "#EF4444" // red
        OrderStatus.REFUNDED -> "#6B7280" // gray
        OrderStatus.RETURNED -> "#6B7280" // gray
    }
}

fun OrderStatus.getActiveStatuses(): List<OrderStatus> {
    return listOf(OrderStatus.PENDING, OrderStatus.CONFIRMED, OrderStatus.PROCESSING, OrderStatus.SHIPPED)
}

fun OrderStatus.getCompletedStatuses(): List<OrderStatus> {
    return listOf(OrderStatus.DELIVERED, OrderStatus.CANCELLED, OrderStatus.REFUNDED, OrderStatus.RETURNED)
}

fun OrderStatus.isActive(): Boolean {
    return this in getActiveStatuses()
}

fun OrderStatus.isCompleted(): Boolean {
    return this in getCompletedStatuses()
}

@Serializable
data class OrderStatusUpdateMessage(
    val orderId: String,
    val orderNumber: String,
    val previousStatus: OrderStatus,
    val newStatus: OrderStatus,
    val updateMessage: String,
    val trackingNumber: String? = null,
    val estimatedDelivery: String? = null,
    val updateTimestamp: String,
    val shopId: String,
    val shopName: String
)

@Serializable
data class PaymentRequestMessage(
    val paymentId: String,
    val amount: Double,
    val currency: String = "THB",
    val description: String,
    val dueDate: String? = null,
    val paymentMethods: List<PaymentMethod>,
    val shopId: String,
    val shopName: String,
    val relatedOrderId: String? = null,
    val relatedInvoiceId: String? = null
)

@Serializable
data class PaymentMethod(
    val type: PaymentType,
    val displayName: String,
    val details: Map<String, String> = emptyMap() // e.g., bank account details
)

@Serializable
enum class PaymentType(val value: String, val displayName: String) {
    BANK_TRANSFER("bank_transfer", "Bank Transfer"),
    CREDIT_CARD("credit_card", "Credit Card"),
    DEBIT_CARD("debit_card", "Debit Card"),
    DIGITAL_WALLET("digital_wallet", "Digital Wallet"),
    CASH("cash", "Cash on Delivery"),
    QR_CODE("qr_code", "QR Code Payment");

    companion object {
        fun fromValue(value: String): PaymentType? {
            return values().find { it.value == value }
        }
    }
}

@Serializable
data class QuotationMessage(
    val quotationId: String,
    val quotationNumber: String,
    val items: List<QuotationItem>,
    val subtotal: Double,
    val tax: Double,
    val total: Double,
    val currency: String = "THB",
    val validUntil: String,
    val status: QuotationStatus,
    val shopId: String,
    val shopName: String,
    val notes: String? = null,
    val termsAndConditions: String? = null
)

@Serializable
data class QuotationItem(
    val productId: String,
    val name: String,
    val description: String? = null,
    val quantity: Int,
    val unitPrice: Double,
    val totalPrice: Double,
    val variantId: String? = null,
    val variantName: String? = null
)

@Serializable
enum class QuotationStatus(val value: String, val displayName: String) {
    DRAFT("draft", "Draft"),
    SENT("sent", "Sent"),
    VIEWED("viewed", "Viewed"),
    ACCEPTED("accepted", "Accepted"),
    REJECTED("rejected", "Rejected"),
    EXPIRED("expired", "Expired"),
    CONVERTED("converted", "Converted to Order");

    companion object {
        fun fromValue(value: String): QuotationStatus? {
            return values().find { it.value == value }
        }
    }
}

/**
 * Helper functions for commerce messages
 */

fun Message.getProductMessage(): ProductMessage? {
    return if (type == MessageType.PRODUCT) {
        try {
            kotlinx.serialization.json.Json.decodeFromString<ProductMessage>(getDisplayContent())
        } catch (e: Exception) {
            null
        }
    } else null
}

fun Message.getInvoiceMessage(): InvoiceMessage? {
    return if (type == MessageType.INVOICE) {
        try {
            kotlinx.serialization.json.Json.decodeFromString<InvoiceMessage>(getDisplayContent())
        } catch (e: Exception) {
            null
        }
    } else null
}

fun Message.getOrderMessage(): OrderMessage? {
    return if (type == MessageType.ORDER) {
        try {
            kotlinx.serialization.json.Json.decodeFromString<OrderMessage>(getDisplayContent())
        } catch (e: Exception) {
            null
        }
    } else null
}

fun Message.getOrderStatusUpdateMessage(): OrderStatusUpdateMessage? {
    return if (type == MessageType.ORDER_STATUS_UPDATE) {
        try {
            kotlinx.serialization.json.Json.decodeFromString<OrderStatusUpdateMessage>(getDisplayContent())
        } catch (e: Exception) {
            null
        }
    } else null
}

fun Message.getPaymentRequestMessage(): PaymentRequestMessage? {
    return if (type == MessageType.PAYMENT_REQUEST) {
        try {
            kotlinx.serialization.json.Json.decodeFromString<PaymentRequestMessage>(getDisplayContent())
        } catch (e: Exception) {
            null
        }
    } else null
}

fun Message.getQuotationMessage(): QuotationMessage? {
    return if (type == MessageType.QUOTATION) {
        try {
            kotlinx.serialization.json.Json.decodeFromString<QuotationMessage>(getDisplayContent())
        } catch (e: Exception) {
            null
        }
    } else null
}