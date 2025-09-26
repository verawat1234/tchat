package com.tchat.mobile.components

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier

/**
 * TchatCard - Cross-platform card component with expect/actual pattern
 *
 * Features:
 * - 4 visual variants (Elevated, Outlined, Filled, Glass)
 * - 3 size variants (Compact, Standard, Expanded)
 * - Interactive support with press animations
 * - Header and footer component slots
 * - Flexible content composition
 * - Platform-native implementation (Material3 on Android, SwiftUI-style on iOS)
 */

enum class TchatCardVariant {
    Elevated,    // 4dp shadow elevation with white background
    Outlined,    // 1dp border without elevation for subtle containers
    Filled,      // Surface color background for grouped content sections
    Glass        // Semi-transparent glassmorphism effect (80% opacity)
}

enum class TchatCardSize {
    Compact,   // 8dp padding for dense information display
    Standard,  // 16dp padding for typical card content
    Expanded   // 24dp padding for spacious layouts with breathing room
}

/**
 * Cross-platform card component using expect/actual pattern
 * Platform-specific implementations provide native elevation and interaction handling
 */
@Composable
expect fun TchatCard(
    modifier: Modifier = Modifier,
    variant: TchatCardVariant = TchatCardVariant.Elevated,
    size: TchatCardSize = TchatCardSize.Standard,
    onClick: (() -> Unit)? = null,
    enabled: Boolean = true,
    content: @Composable () -> Unit
)