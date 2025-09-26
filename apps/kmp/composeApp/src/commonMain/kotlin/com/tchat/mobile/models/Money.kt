package com.tchat.mobile.models

import kotlinx.serialization.Serializable
import kotlin.math.round

/**
 * Money model for handling currency amounts across the application
 * Supports multiple currencies and arithmetic operations
 */
@Serializable
data class Money(
    val amount: Double,
    val currency: String = "USD"
) {
    fun add(other: Money): Money {
        require(currency == other.currency) { "Cannot add different currencies: $currency and ${other.currency}" }
        return Money(amount + other.amount, currency)
    }

    fun subtract(other: Money): Money {
        require(currency == other.currency) { "Cannot subtract different currencies: $currency and ${other.currency}" }
        return Money(amount - other.amount, currency)
    }

    fun multiply(multiplier: Int): Money = Money(amount * multiplier, currency)
    fun multiply(multiplier: Double): Money = Money(amount * multiplier, currency)

    companion object {
        fun zero(currency: String = "USD") = Money(0.0, currency)
    }
}

/**
 * Utility function to format money for display
 */
fun formatMoney(amount: Double, currency: String): String {
    // Simple cross-platform formatting - round to 2 decimal places
    val rounded = round(amount * 100) / 100
    val formattedAmount = rounded.toString()

    return when (currency) {
        "USD" -> "$$formattedAmount"
        "EUR" -> "€$formattedAmount"
        "GBP" -> "£$formattedAmount"
        "THB" -> "฿$formattedAmount"
        "SGD" -> "S$$formattedAmount"
        else -> "$currency $formattedAmount"
    }
}