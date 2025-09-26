package com.tchat.mobile.components

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier

/**
 * TchatSpinner - Cross-platform loading indicator component with expect/actual pattern
 *
 * Features:
 * - 4 semantic variants (Default, Success, Warning, Error) for contextual feedback
 * - 3 size variants (Small, Medium, Large) for different UI contexts
 * - Platform-native animations (Material ripple vs iOS spring animations)
 * - Indeterminate and determinate progress modes
 * - Accessibility support with loading state announcements
 * - Custom color theming based on design system
 * - Smooth rotation animations with platform-optimized timing
 */

enum class TchatSpinnerVariant {
    Default,   // Primary brand color for general loading states
    Success,   // Green for success operations and confirmations
    Warning,   // Amber for warning states and caution operations
    Error      // Red for error states and failed operations
}

enum class TchatSpinnerSize {
    Small,   // 16dp size, 2dp stroke - For inline loading in buttons/fields
    Medium,  // 20dp size, 2.5dp stroke - Standard loading for content areas
    Large    // 24dp size, 3dp stroke - Prominent loading for full screen
}

/**
 * Cross-platform spinner component using expect/actual pattern
 * Platform-specific implementations provide native animation timing and styling
 */
@Composable
expect fun TchatSpinner(
    modifier: Modifier = Modifier,
    variant: TchatSpinnerVariant = TchatSpinnerVariant.Default,
    size: TchatSpinnerSize = TchatSpinnerSize.Medium,
    progress: Float? = null, // 0.0-1.0 for determinate progress, null for indeterminate
    strokeWidth: Float? = null, // Override default stroke width
    contentDescription: String? = null
)