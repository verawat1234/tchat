package com.tchat.mobile.components

import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector

/**
 * TchatSwitch - Cross-platform toggle switch component with expect/actual pattern
 *
 * Features:
 * - On/Off states with smooth spring animations
 * - Loading state support with spinner indicator
 * - 3 size variants (Small, Medium, Large) for different use cases
 * - Custom labels and descriptions with rich typography
 * - Advanced interaction states and haptic feedback
 * - Platform-native implementation (Material3 on Android, SwiftUI-style on iOS)
 * - Full accessibility support with semantic descriptions
 * - Design system integration with consistent theming
 */

enum class TchatSwitchSize {
    Small,   // 24x14dp - Compact switch for dense layouts and settings lists
    Medium,  // 32x20dp - Standard switch size for most form use cases
    Large    // 40x24dp - Prominent switch for important toggles and primary actions
}

/**
 * Cross-platform switch component using expect/actual pattern
 * Platform-specific implementations provide native toggle behavior and styling
 */
@Composable
expect fun TchatSwitch(
    checked: Boolean,
    onCheckedChange: ((Boolean) -> Unit)?,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
    size: TchatSwitchSize = TchatSwitchSize.Medium,
    isLoading: Boolean = false,
    label: String? = null,
    description: String? = null,
    leadingIcon: ImageVector? = null,
    interactionSource: MutableInteractionSource = remember { MutableInteractionSource() },
    contentDescription: String? = null
)