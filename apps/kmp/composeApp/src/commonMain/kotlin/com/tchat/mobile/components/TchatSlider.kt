package com.tchat.mobile.components

import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector

/**
 * TchatSlider - Cross-platform range slider component with expect/actual pattern
 *
 * Features:
 * - Single/Range selection with smooth thumb interactions
 * - Step increments with snap-to-step behavior
 * - Value labels and tooltips with customizable formatting
 * - Platform-native styling (Material3 on Android, SwiftUI-style on iOS)
 * - 3 size variants (Small, Medium, Large) for different use cases
 * - Advanced accessibility support with semantic descriptions
 * - Custom value formatting and display options
 * - Loading states and disabled interactions
 */

enum class TchatSliderMode {
    Single,    // Single value slider with one thumb
    Range      // Range slider with two thumbs (start and end values)
}

enum class TchatSliderSize {
    Small,   // Compact slider for dense layouts and settings
    Medium,  // Standard slider size for most form use cases
    Large    // Prominent slider for important value selection
}

/**
 * Cross-platform slider component using expect/actual pattern
 * Platform-specific implementations provide native slider behavior and styling
 */
@Composable
expect fun TchatSlider(
    value: Float,
    onValueChange: (Float) -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
    valueRange: ClosedFloatingPointRange<Float> = 0f..1f,
    steps: Int = 0,
    size: TchatSliderSize = TchatSliderSize.Medium,
    showValueLabel: Boolean = false,
    valueFormatter: (Float) -> String = { it.toString() },
    label: String? = null,
    leadingIcon: ImageVector? = null,
    trailingIcon: ImageVector? = null,
    interactionSource: MutableInteractionSource = remember { MutableInteractionSource() },
    contentDescription: String? = null
)

/**
 * Range slider for selecting a range of values
 */
@Composable
expect fun TchatRangeSlider(
    value: ClosedFloatingPointRange<Float>,
    onValueChange: (ClosedFloatingPointRange<Float>) -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
    valueRange: ClosedFloatingPointRange<Float> = 0f..1f,
    steps: Int = 0,
    size: TchatSliderSize = TchatSliderSize.Medium,
    showValueLabels: Boolean = false,
    valueFormatter: (Float) -> String = { it.toString() },
    label: String? = null,
    leadingIcon: ImageVector? = null,
    trailingIcon: ImageVector? = null,
    startInteractionSource: MutableInteractionSource = remember { MutableInteractionSource() },
    endInteractionSource: MutableInteractionSource = remember { MutableInteractionSource() },
    contentDescription: String? = null
)