package com.tchat.mobile.components

import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector

/**
 * TchatSelect - Cross-platform dropdown selection component with expect/actual pattern
 *
 * Features:
 * - Single/Multi-selection modes with advanced state management
 * - Search/filter capabilities with real-time filtering
 * - Custom option templates with icons and descriptions
 * - Platform-native dropdown behavior (Material3 ExposedDropdownMenu on Android, SwiftUI-style on iOS)
 * - 3 validation states (None, Valid, Invalid) with animated feedback
 * - 3 size variants (Small, Medium, Large) for different use cases
 * - Advanced keyboard navigation and accessibility support
 * - Loading states with skeleton animations
 */

data class SelectOption(
    val value: String,
    val label: String,
    val description: String? = null,
    val icon: ImageVector? = null,
    val enabled: Boolean = true,
    val group: String? = null
)

enum class TchatSelectMode {
    Single,     // Single selection dropdown
    Multiple    // Multi-selection with chips/tags
}

enum class TchatSelectValidationState {
    None,    // Default neutral state
    Valid,   // Green success border with checkmark
    Invalid  // Red error border with error message
}

enum class TchatSelectSize {
    Small,   // 14sp text, 40dp height, compact touch targets
    Medium,  // 16sp text, 48dp height, standard dropdown size
    Large    // 18sp text, 56dp height, prominent selection field
}

/**
 * Cross-platform select component using expect/actual pattern
 * Platform-specific implementations provide native dropdown behavior and styling
 */
@Composable
expect fun TchatSelect(
    options: List<SelectOption>,
    selectedValues: List<String>,
    onSelectionChange: (List<String>) -> Unit,
    modifier: Modifier = Modifier,
    mode: TchatSelectMode = TchatSelectMode.Single,
    validationState: TchatSelectValidationState = TchatSelectValidationState.None,
    size: TchatSelectSize = TchatSelectSize.Medium,
    enabled: Boolean = true,
    placeholder: String = "Select an option",
    label: String? = null,
    supportingText: String? = null,
    errorMessage: String? = null,
    searchEnabled: Boolean = false,
    searchPlaceholder: String = "Search options...",
    maxSelections: Int? = null,
    isLoading: Boolean = false,
    leadingIcon: ImageVector? = null,
    trailingIcon: ImageVector? = null,
    onTrailingIconClick: (() -> Unit)? = null,
    interactionSource: MutableInteractionSource = remember { MutableInteractionSource() },
    contentDescription: String? = null
)

/**
 * Single-selection convenience function
 */
@Composable
expect fun TchatSingleSelect(
    options: List<SelectOption>,
    selectedValue: String?,
    onSelectionChange: (String?) -> Unit,
    modifier: Modifier = Modifier,
    validationState: TchatSelectValidationState = TchatSelectValidationState.None,
    size: TchatSelectSize = TchatSelectSize.Medium,
    enabled: Boolean = true,
    placeholder: String = "Select an option",
    label: String? = null,
    supportingText: String? = null,
    errorMessage: String? = null,
    searchEnabled: Boolean = false,
    searchPlaceholder: String = "Search options...",
    isLoading: Boolean = false,
    leadingIcon: ImageVector? = null,
    trailingIcon: ImageVector? = null,
    onTrailingIconClick: (() -> Unit)? = null,
    interactionSource: MutableInteractionSource = remember { MutableInteractionSource() },
    contentDescription: String? = null
)