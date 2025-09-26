package com.tchat.mobile.components

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier

/**
 * TchatToast - Toast notification component with platform-native implementations
 *
 * Features:
 * - Auto-dismiss with configurable timeout (1-10 seconds)
 * - Position variants (Top, Bottom, Center) for flexible placement
 * - Queue management for multiple toasts with smart ordering
 * - Platform-native animation styles (slide, fade, bounce)
 * - Semantic variants with appropriate icons and colors
 * - Manual dismiss support with swipe gestures
 * - Accessibility announcements for screen readers
 * - Rich content support with icons and actions
 */

enum class TchatToastVariant {
    Info,       // Informational messages
    Success,    // Success confirmations
    Warning,    // Warning notifications
    Error       // Error alerts
}

enum class TchatToastPosition {
    Top,        // Top of screen with safe area respect
    Bottom,     // Bottom of screen with safe area respect
    Center      // Center of screen for important messages
}

/**
 * Cross-platform toast component using expect/actual pattern
 * Platform-specific implementations provide native animations and positioning
 */
@Composable
expect fun TchatToast(
    message: String,
    variant: TchatToastVariant = TchatToastVariant.Info,
    position: TchatToastPosition = TchatToastPosition.Bottom,
    modifier: Modifier = Modifier,
    duration: Long = 3000L, // Duration in milliseconds (1000-10000)
    dismissible: Boolean = true,
    onDismiss: (() -> Unit)? = null,
    action: ToastAction? = null,
    icon: (@Composable () -> Unit)? = null
)

/**
 * Toast action configuration for interactive toasts
 */
data class ToastAction(
    val label: String,
    val onClick: () -> Unit
)

/**
 * Toast queue manager for handling multiple toasts
 */
expect class TchatToastManager {
    /**
     * Show a toast with queue management
     * Automatically handles positioning and timing
     */
    fun showToast(
        message: String,
        variant: TchatToastVariant = TchatToastVariant.Info,
        position: TchatToastPosition = TchatToastPosition.Bottom,
        duration: Long = 3000L,
        action: ToastAction? = null
    )

    /**
     * Clear all pending toasts
     */
    fun clearAll()

    /**
     * Clear toasts of specific variant
     */
    fun clearVariant(variant: TchatToastVariant)
}