package com.tchat.mobile.components

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier

/**
 * TchatSkeleton - Cross-platform skeleton loading component with expect/actual pattern
 *
 * Features:
 * - Platform-native shimmer animations (Material ripple vs iOS pulse)
 * - Multiple shape variants for different content types
 * - Size customization for flexible loading states
 * - Accessibility announcements for loading states
 * - Design system integration with TchatColors
 */

enum class TchatSkeletonShape {
    Rectangle,  // Standard rectangular placeholder for text/content blocks
    Circle,     // Circular placeholder for avatars and profile images
    Rounded,    // Rounded rectangle for buttons and cards
    Line        // Thin line placeholder for single lines of text
}

enum class TchatSkeletonSize {
    Small,      // Compact skeleton for dense layouts
    Medium,     // Standard skeleton size for typical content
    Large       // Prominent skeleton for hero content
}

/**
 * Cross-platform skeleton loading component using expect/actual pattern
 * Platform-specific implementations provide native loading animations and accessibility
 */
@Composable
expect fun TchatSkeleton(
    modifier: Modifier = Modifier,
    shape: TchatSkeletonShape = TchatSkeletonShape.Rectangle,
    size: TchatSkeletonSize = TchatSkeletonSize.Medium,
    animated: Boolean = true
)