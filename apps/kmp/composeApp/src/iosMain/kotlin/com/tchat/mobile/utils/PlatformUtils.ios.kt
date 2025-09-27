package com.tchat.mobile.utils

import kotlin.math.pow
import kotlin.math.round

actual object PlatformUtils {
    actual fun currentTimeMillis(): Long {
        return kotlinx.datetime.Clock.System.now().toEpochMilliseconds()
    }

    actual fun formatDecimal(value: Double, decimalPlaces: Int): String {
        val multiplier = 10.0.pow(decimalPlaces)
        val rounded = round(value * multiplier) / multiplier
        // Simple formatting for iOS without using String.format
        return when (decimalPlaces) {
            0 -> rounded.toInt().toString()
            1 -> "${(rounded * 10).toInt() / 10.0}"
            2 -> "${(rounded * 100).toInt() / 100.0}"
            else -> rounded.toString()
        }
    }
}