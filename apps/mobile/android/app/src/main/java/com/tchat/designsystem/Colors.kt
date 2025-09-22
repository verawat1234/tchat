package com.tchat.designsystem

import androidx.compose.ui.graphics.Color

/**
 * Color design tokens matching TailwindCSS v4 color palette
 */
object Colors {

    // Brand Colors

    /** Primary brand color - Blue 500
     *  Maps to TailwindCSS: blue-500 */
    val primary = Color(0xFF3B82F6)

    /** Secondary brand color - Gray 600
     *  Maps to TailwindCSS: gray-600 */
    val secondary = Color(0xFF4B5563)

    /** Accent color - Indigo 500
     *  Maps to TailwindCSS: indigo-500 */
    val accent = Color(0xFF6366F1)

    // Semantic Colors

    /** Success color - Green 500
     *  Maps to TailwindCSS: green-500 */
    val success = Color(0xFF10B981)

    /** Warning color - Amber 500
     *  Maps to TailwindCSS: amber-500 */
    val warning = Color(0xFFF59E0B)

    /** Error color - Red 500
     *  Maps to TailwindCSS: red-500 */
    val error = Color(0xFFEF4444)

    /** Info color - Blue 400
     *  Maps to TailwindCSS: blue-400 */
    val info = Color(0xFF60A5FA)

    // Surface Colors

    /** Background color - White
     *  Maps to TailwindCSS: white */
    val background = Color(0xFFFFFFFF)

    /** Surface color - Gray 50
     *  Maps to TailwindCSS: gray-50 */
    val surface = Color(0xFFF9FAFB)

    /** Card background - White with subtle shadow
     *  Maps to TailwindCSS: white */
    val cardBackground = Color(0xFFFFFFFF)

    /** Modal overlay - Black with opacity
     *  Maps to TailwindCSS: black/50 */
    val overlay = Color(0x80000000)

    // Text Colors

    /** Primary text color - Gray 900
     *  Maps to TailwindCSS: gray-900 */
    val textPrimary = Color(0xFF111827)

    /** Secondary text color - Gray 600
     *  Maps to TailwindCSS: gray-600 */
    val textSecondary = Color(0xFF4B5563)

    /** Tertiary text color - Gray 400
     *  Maps to TailwindCSS: gray-400 */
    val textTertiary = Color(0xFF9CA3AF)

    /** Disabled text color - Gray 300
     *  Maps to TailwindCSS: gray-300 */
    val textDisabled = Color(0xFFD1D5DB)

    /** Text on primary color - White */
    val textOnPrimary = Color.White

    /** Text on dark backgrounds - White */
    val textOnDark = Color.White

    // Border Colors

    /** Default border color - Gray 200
     *  Maps to TailwindCSS: gray-200 */
    val border = Color(0xFFE5E7EB)

    /** Focus border color - Blue 500
     *  Maps to TailwindCSS: blue-500 */
    val borderFocus = Color(0xFF3B82F6)

    /** Error border color - Red 300
     *  Maps to TailwindCSS: red-300 */
    val borderError = Color(0xFFFCA5A5)

    /** Divider color - Gray 200
     *  Maps to TailwindCSS: gray-200 */
    val divider = Color(0xFFE5E7EB)

    // Interactive States

    /** Hover state colors */
    object Hover {
        val primary = Color(0xFF2563EB)     // blue-600
        val secondary = Color(0xFF374151)   // gray-700
        val surface = Color(0xFFF3F4F6)     // gray-100
    }

    /** Pressed state colors */
    object Pressed {
        val primary = Color(0xFF1D4ED8)     // blue-700
        val secondary = Color(0xFF1F2937)   // gray-800
        val surface = Color(0xFFE5E7EB)     // gray-200
    }

    /** Disabled state colors */
    object Disabled {
        val background = Color(0xFFF3F4F6)  // gray-100
        val text = Color(0xFF9CA3AF)        // gray-400
        val border = Color(0xFFD1D5DB)      // gray-300
    }

    // Tab/Navigation Colors

    /** Tab bar background - Surface with transparency */
    val tabBarBackground = Color(0xF2F9FAFB)

    /** Navigation bar background - Background */
    val navigationBackground = background

    /** Selected tab color - Primary */
    val tabSelected = primary

    /** Unselected tab color - Gray 400 */
    val tabUnselected = Color(0xFF9CA3AF)

    // Shadow Colors

    /** Light shadow color */
    val shadowLight = Color(0x1A000000)

    /** Medium shadow color */
    val shadowMedium = Color(0x26000000)

    /** Heavy shadow color */
    val shadowHeavy = Color(0x40000000)
}

// Dark Mode Colors

/** Dark mode color variants */
object DarkColors {

    // Brand colors remain the same
    val primary = Color(0xFF3B82F6)
    val secondary = Color(0xFF6B7280)
    val accent = Color(0xFF6366F1)

    // Semantic colors adjusted for dark mode
    val success = Color(0xFF10B981)
    val warning = Color(0xFFF59E0B)
    val error = Color(0xFFEF4444)
    val info = Color(0xFF60A5FA)

    // Dark surfaces
    val background = Color(0xFF111827)     // gray-900
    val surface = Color(0xFF1F2937)        // gray-800
    val cardBackground = Color(0xFF374151) // gray-700
    val overlay = Color(0xB3000000)

    // Dark text colors
    val textPrimary = Color(0xFFF9FAFB)    // gray-50
    val textSecondary = Color(0xFFD1D5DB)  // gray-300
    val textTertiary = Color(0xFF9CA3AF)   // gray-400
    val textDisabled = Color(0xFF6B7280)   // gray-500

    // Dark borders
    val border = Color(0xFF374151)         // gray-700
    val borderFocus = Color(0xFF3B82F6)    // blue-500
    val borderError = Color(0xFFEF4444)    // red-500
    val divider = Color(0xFF374151)        // gray-700

    // Dark navigation
    val tabBarBackground = Color(0xF21F2937)
    val navigationBackground = Color(0xFF111827)

    // Interactive states for dark mode
    object Hover {
        val primary = Color(0xFF2563EB)     // blue-600
        val secondary = Color(0xFF4B5563)   // gray-600
        val surface = Color(0xFF374151)     // gray-700
    }

    object Pressed {
        val primary = Color(0xFF1D4ED8)     // blue-700
        val secondary = Color(0xFF374151)   // gray-700
        val surface = Color(0xFF4B5563)     // gray-600
    }

    object Disabled {
        val background = Color(0xFF374151)  // gray-700
        val text = Color(0xFF6B7280)        // gray-500
        val border = Color(0xFF4B5563)      // gray-600
    }
}

// Color Utilities

/**
 * Get colors for current theme
 */
fun Colors.colorsFor(isDark: Boolean): Any {
    return if (isDark) DarkColors else Colors
}

/**
 * Validate color accessibility contrast
 */
fun validateContrast(foreground: Color, background: Color): Boolean {
    // Implementation would check WCAG contrast requirements
    // This is a placeholder - full implementation would calculate luminance
    return true
}

/**
 * Create color with alpha
 */
fun Color.withAlpha(alpha: Float): Color {
    return this.copy(alpha = alpha)
}

/**
 * Blend two colors
 */
fun Color.blendWith(other: Color, ratio: Float): Color {
    val r1 = this.red
    val g1 = this.green
    val b1 = this.blue
    val a1 = this.alpha

    val r2 = other.red
    val g2 = other.green
    val b2 = other.blue
    val a2 = other.alpha

    return Color(
        red = r1 + (r2 - r1) * ratio,
        green = g1 + (g2 - g1) * ratio,
        blue = b1 + (b2 - b1) * ratio,
        alpha = a1 + (a2 - a1) * ratio
    )
}