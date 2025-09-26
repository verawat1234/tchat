package com.tchat.mobile.utils

/**
 * Cross-platform utilities for KMP
 */
expect object PlatformUtils {
    fun currentTimeMillis(): Long
    fun formatDecimal(value: Double, decimalPlaces: Int): String
}