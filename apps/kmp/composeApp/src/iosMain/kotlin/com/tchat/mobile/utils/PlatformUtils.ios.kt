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
        return "%.${decimalPlaces}f".format(rounded)
    }
}