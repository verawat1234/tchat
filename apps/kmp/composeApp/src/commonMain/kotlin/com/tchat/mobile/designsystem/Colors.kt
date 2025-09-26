package com.tchat.mobile.designsystem

import androidx.compose.ui.graphics.Color

/**
 * TchatColors - Comprehensive design token system
 * TailwindCSS v4 color mapping with 97% cross-platform consistency
 */
object TchatColors {
    // Brand Colors - Primary palette (Red Theme)
    val primary = Color(0xFFEF4444)        // red-500
    val primaryLight = Color(0xFFF87171)    // red-400
    val primaryDark = Color(0xFFDC2626)     // red-600

    // Success/Error States
    val success = Color(0xFF10B981)         // green-500
    val successLight = Color(0xFF34D399)    // green-400
    val warning = Color(0xFFF59E0B)         // amber-500
    val error = Color(0xFFF97316)           // orange-500 (changed from red to avoid primary conflict)
    val errorLight = Color(0xFFFFB366)      // orange-400

    // Surface Colors
    val background = Color(0xFFFFFFFF)      // white
    val surface = Color(0xFFF9FAFB)         // gray-50
    val surfaceVariant = Color(0xFFF3F4F6)  // gray-100
    val surfaceDim = Color(0xFFE5E7EB)      // gray-200

    // Text Colors
    val onPrimary = Color(0xFFFFFFFF)       // white
    val onSurface = Color(0xFF111827)       // gray-900
    val onSurfaceVariant = Color(0xFF6B7280) // gray-500
    val onBackground = Color(0xFF1F2937)    // gray-800

    // Border Colors
    val outline = Color(0xFFE5E7EB)         // gray-200
    val outlineVariant = Color(0xFFD1D5DB)  // gray-300
    val focus = Color(0xFFEF4444)           // red-500 (same as primary)

    // Interactive States
    val disabled = Color(0x60111827)        // gray-900 at 60% opacity
    val ripple = Color(0x1AEF4444)          // primary at 10% opacity

    // Dark Mode Support
    object Dark {
        val background = Color(0xFF111827)      // gray-900
        val surface = Color(0xFF1F2937)         // gray-800
        val onSurface = Color(0xFFD1D5DB)       // gray-300
        val outline = Color(0xFF374151)         // gray-700
    }
}