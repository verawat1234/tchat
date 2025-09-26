package com.tchat.mobile.components

import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier

/**
 * TchatButton - Cross-platform button component with expect/actual pattern
 *
 * Features:
 * - 5 visual variants (Primary, Secondary, Ghost, Destructive, Outline)
 * - 3 size variants (Small, Medium, Large)
 * - Loading states with animated progress indicators
 * - Advanced interaction states (press animations, focus, disabled)
 * - Full accessibility support with dynamic labels
 * - Platform-native implementation (Material3 on Android, SwiftUI-style on iOS)
 */

enum class TchatButtonVariant {
    Primary,      // Brand-colored call-to-action buttons
    Secondary,    // Subtle surface-based secondary actions
    Ghost,        // Transparent background with primary text
    Destructive,  // Error-colored for dangerous actions
    Outline       // Transparent with bordered outline
}

enum class TchatButtonSize {
    Small,   // 32dp height, 14sp text, compact touch targets
    Medium,  // 44dp height (iOS HIG compliance), 16sp text
    Large    // 48dp height, 18sp text, prominent actions
}

/**
 * Cross-platform button component using expect/actual pattern
 * Platform-specific implementations provide native look and feel
 */
@Composable
expect fun TchatButton(
    onClick: () -> Unit,
    text: String,
    modifier: Modifier = Modifier,
    variant: TchatButtonVariant = TchatButtonVariant.Primary,
    size: TchatButtonSize = TchatButtonSize.Medium,
    enabled: Boolean = true,
    loading: Boolean = false,
    leadingIcon: (@Composable () -> Unit)? = null,
    trailingIcon: (@Composable () -> Unit)? = null,
    interactionSource: MutableInteractionSource = remember { MutableInteractionSource() },
    contentDescription: String? = null
)