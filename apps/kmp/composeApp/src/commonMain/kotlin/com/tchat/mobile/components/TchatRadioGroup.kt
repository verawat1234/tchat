package com.tchat.mobile.components

import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

/**
 * TchatRadioGroup - Cross-platform radio button group component with expect/actual pattern
 *
 * Features:
 * - Options list with single selection
 * - 3 validation states (None, Valid, Invalid) with animated borders
 * - Horizontal/Vertical orientation support
 * - Custom option rendering with labels
 * - Full accessibility support with semantic descriptions
 * - Platform-native implementation (Material3 on Android, SwiftUI-style on iOS)
 * - Advanced interaction states and animations
 */

data class RadioOption(
    val value: String,
    val label: String,
    val enabled: Boolean = true,
    val description: String? = null
)

enum class TchatRadioGroupOrientation {
    Horizontal,  // Options arranged in a row
    Vertical     // Options arranged in a column
}

enum class TchatRadioGroupValidationState {
    None,    // Default neutral state
    Valid,   // Green success border with checkmark
    Invalid  // Red error border with error message
}

enum class TchatRadioGroupSize {
    Small,   // 16dp radio buttons, 12sp text, compact spacing
    Medium,  // 20dp radio buttons, 14sp text, standard spacing
    Large    // 24dp radio buttons, 16sp text, spacious layout
}

/**
 * Cross-platform radio group component using expect/actual pattern
 * Platform-specific implementations provide native styling and interactions
 */
@Composable
expect fun TchatRadioGroup(
    options: List<RadioOption>,
    selectedValue: String?,
    onSelectionChange: (String) -> Unit,
    modifier: Modifier = Modifier,
    orientation: TchatRadioGroupOrientation = TchatRadioGroupOrientation.Vertical,
    validationState: TchatRadioGroupValidationState = TchatRadioGroupValidationState.None,
    size: TchatRadioGroupSize = TchatRadioGroupSize.Medium,
    enabled: Boolean = true,
    label: String? = null,
    supportingText: String? = null,
    errorMessage: String? = null,
    spacing: Int = when (size) {
        TchatRadioGroupSize.Small -> 8
        TchatRadioGroupSize.Medium -> 12
        TchatRadioGroupSize.Large -> 16
    },
    interactionSource: MutableInteractionSource = remember { MutableInteractionSource() },
    contentDescription: String? = null
)