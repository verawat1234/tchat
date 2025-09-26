package com.tchat.mobile.components

import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier

/**
 * TchatCheckbox - Cross-platform checkbox component with expect/actual pattern
 *
 * Features:
 * - Material3 (Android) vs SwiftUI-inspired (iOS) styling
 * - Multiple size variants for different use cases
 * - Comprehensive interaction states and animations
 * - Full accessibility support with semantic descriptions
 * - Design system integration with TchatColors
 */

enum class TchatCheckboxSize {
    Small,      // 16dp - Compact checkbox for dense layouts
    Medium,     // 20dp - Standard checkbox size for most use cases
    Large       // 24dp - Prominent checkbox for important selections
}

/**
 * Cross-platform checkbox component using expect/actual pattern
 * Platform-specific implementations provide native styling and interactions
 */
@Composable
expect fun TchatCheckbox(
    checked: Boolean,
    onCheckedChange: ((Boolean) -> Unit)?,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
    size: TchatCheckboxSize = TchatCheckboxSize.Medium,
    interactionSource: MutableInteractionSource = remember { MutableInteractionSource() },
    label: String? = null
)