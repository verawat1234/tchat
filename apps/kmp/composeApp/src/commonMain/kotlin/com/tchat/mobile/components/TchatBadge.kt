package com.tchat.mobile.components

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.TextStyle

/**
 * TchatBadge - Cross-platform badge component with expect/actual pattern
 *
 * Features:
 * - 5 semantic variants (Default, Success, Warning, Error, Info)
 * - 3 size variants (Small, Medium, Large)
 * - Numeric and text content support
 * - Platform-specific styling (rounded corners, shadows, borders)
 * - Accessibility support with semantic descriptions
 * - Auto-hide for zero values in numeric badges
 */

enum class TchatBadgeVariant {
    Default,      // Primary brand color for general notifications
    Success,      // Green for success states and positive notifications
    Warning,      // Amber/Orange for warnings and caution states
    Error,        // Red for errors and critical alerts
    Info          // Blue/Gray for informational content
}

enum class TchatBadgeSize {
    Small,   // 16dp height, 12sp text, compact badges
    Medium,  // 20dp height, 14sp text, standard size
    Large    // 24dp height, 16sp text, prominent badges
}

/**
 * Cross-platform badge component using expect/actual pattern
 * Platform-specific implementations provide native styling and animations
 */
@Composable
expect fun TchatBadge(
    text: String,
    modifier: Modifier = Modifier,
    variant: TchatBadgeVariant = TchatBadgeVariant.Default,
    size: TchatBadgeSize = TchatBadgeSize.Medium,
    count: Int? = null, // For numeric badges - auto-hides when zero
    maxCount: Int = 99,  // Maximum count to display (99+)
    showZero: Boolean = false, // Whether to show badge when count is 0
    contentDescription: String? = null
)