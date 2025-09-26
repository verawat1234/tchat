package com.tchat.mobile.utils

import java.text.DecimalFormat

actual object PlatformUtils {
    actual fun currentTimeMillis(): Long = System.currentTimeMillis()

    actual fun formatDecimal(value: Double, decimalPlaces: Int): String {
        val pattern = "#." + "#".repeat(decimalPlaces)
        val formatter = DecimalFormat(pattern)
        return formatter.format(value)
    }
}