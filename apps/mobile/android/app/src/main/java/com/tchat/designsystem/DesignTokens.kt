package com.tchat.designsystem

import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.unit.Dp

/**
 * Central design token system for TchatApp
 * Provides platform-specific implementation of TailwindCSS v4 design tokens
 */
object DesignTokens {

    /** Current design system version */
    const val VERSION = "1.0.0"

    /** Typography design tokens */
    val typography = Typography

    /** Color design tokens */
    val colors = Colors

    /** Spacing design tokens */
    val spacing = Spacing

    /** Animation design tokens */
    val animations = Animations

    /** Border radius design tokens */
    val borderRadius = BorderRadius

    /** Shadow design tokens */
    val shadows = Shadows
}

/**
 * Animation design tokens
 */
object Animations {
    val fast = 0.2
    val normal = 0.3
    val slow = 0.5
}

/**
 * Border radius design tokens
 */
object BorderRadius {
    val small = 4.dp
    val medium = 8.dp
    val large = 12.dp
    val full = 9999.dp
}

/**
 * Shadow design tokens
 */
object Shadows {
    val small = androidx.compose.ui.graphics.Shadow(
        color = Color.Black.copy(alpha = 0.1f),
        offset = androidx.compose.ui.geometry.Offset(0f, 1f),
        blurRadius = 2f
    )
    val medium = androidx.compose.ui.graphics.Shadow(
        color = Color.Black.copy(alpha = 0.15f),
        offset = androidx.compose.ui.geometry.Offset(0f, 4f),
        blurRadius = 8f
    )
    val large = androidx.compose.ui.graphics.Shadow(
        color = Color.Black.copy(alpha = 0.2f),
        offset = androidx.compose.ui.geometry.Offset(0f, 8f),
        blurRadius = 16f
    )
}

// Design Token Extensions

/**
 * Get design tokens for current theme
 */
fun DesignTokens.tokensFor(theme: AppTheme): DesignTokens {
    // In a full implementation, this would return theme-specific tokens
    return DesignTokens
}

/**
 * Validate design token consistency
 */
fun DesignTokens.validateTokens(): Boolean {
    // Validate that all required tokens are present and valid
    return VERSION.isNotEmpty() &&
            colors.primary != Color.Unspecified &&
            spacing.md > 0.dp
}

/**
 * App theme enumeration
 */
enum class AppTheme(val value: String) {
    LIGHT("light"),
    DARK("dark"),
    AUTO("auto");

    val displayName: String
        get() = when (this) {
            LIGHT -> "Light"
            DARK -> "Dark"
            AUTO -> "Auto"
        }

    companion object {
        fun fromString(value: String): AppTheme? {
            return values().find { it.value == value }
        }
    }
}

// Color Extensions

/**
 * Create color from hex string
 */
fun Color.Companion.fromHex(hex: String): Color {
    val cleanHex = hex.replace("#", "")
    return when (cleanHex.length) {
        6 -> {
            val colorInt = cleanHex.toLong(16)
            Color(
                red = ((colorInt shr 16) and 0xFF) / 255f,
                green = ((colorInt shr 8) and 0xFF) / 255f,
                blue = (colorInt and 0xFF) / 255f,
                alpha = 1f
            )
        }
        8 -> {
            val colorInt = cleanHex.toLong(16)
            Color(
                red = ((colorInt shr 16) and 0xFF) / 255f,
                green = ((colorInt shr 8) and 0xFF) / 255f,
                blue = (colorInt and 0xFF) / 255f,
                alpha = ((colorInt shr 24) and 0xFF) / 255f
            )
        }
        else -> Color.Black
    }
}