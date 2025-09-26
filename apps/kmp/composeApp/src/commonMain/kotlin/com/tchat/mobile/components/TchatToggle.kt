package com.tchat.mobile.components

import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector

/**
 * TchatToggle - Cross-platform toggle button component with expect/actual pattern
 *
 * Features:
 * - Pressed/Unpressed states with smooth transition animations
 * - Icon and text support with flexible content arrangement
 * - Group toggle functionality for multi-selection scenarios
 * - 3 size variants (Small, Medium, Large) for different use cases
 * - Advanced interaction states with haptic feedback
 * - Platform-native implementation (Material3 on Android, SwiftUI-style on iOS)
 * - Full accessibility support with semantic descriptions
 * - Design system integration with consistent theming
 */

enum class TchatToggleSize {
    Small,   // 32dp height - Compact toggle for dense layouts and toolbars
    Medium,  // 40dp height - Standard toggle size for most button use cases
    Large    // 48dp height - Prominent toggle for important actions
}

enum class TchatToggleVariant {
    Default,   // Standard toggle with background color change
    Outline,   // Outlined toggle with border style
    Text       // Text-only toggle for minimal styling
}

/**
 * Cross-platform toggle button component using expect/actual pattern
 * Platform-specific implementations provide native toggle behavior and styling
 */
@Composable
expect fun TchatToggle(
    pressed: Boolean,
    onPressedChange: ((Boolean) -> Unit)?,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
    size: TchatToggleSize = TchatToggleSize.Medium,
    variant: TchatToggleVariant = TchatToggleVariant.Default,
    text: String? = null,
    icon: ImageVector? = null,
    leadingIcon: ImageVector? = null,
    trailingIcon: ImageVector? = null,
    interactionSource: MutableInteractionSource = remember { MutableInteractionSource() },
    contentDescription: String? = null
)

/**
 * Toggle group for managing multiple related toggles
 */
@Composable
expect fun TchatToggleGroup(
    options: List<String>,
    selectedOptions: Set<String>,
    onSelectionChange: (Set<String>) -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
    size: TchatToggleSize = TchatToggleSize.Medium,
    variant: TchatToggleVariant = TchatToggleVariant.Default,
    allowMultipleSelection: Boolean = true,
    allowEmptySelection: Boolean = true,
    maxSelections: Int? = null,
    interactionSource: MutableInteractionSource = remember { MutableInteractionSource() }
)