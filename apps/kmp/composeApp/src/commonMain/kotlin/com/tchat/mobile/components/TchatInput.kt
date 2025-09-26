package com.tchat.mobile.components

import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.focus.FocusRequester
import androidx.compose.ui.graphics.vector.ImageVector

/**
 * TchatInput - Cross-platform input component with expect/actual pattern
 *
 * Features:
 * - Multiple input types (Text, Email, Password, Number, Search, Multiline)
 * - 3 validation states (None, Valid, Invalid) with animated borders
 * - Leading and trailing icons with interactive features
 * - Password visibility toggle
 * - 3 size variants (Small, Medium, Large)
 * - Advanced animations and focus management
 * - Full accessibility support
 * - Platform-native implementation (Material3 on Android, SwiftUI-style on iOS)
 */

enum class TchatInputType {
    Text,      // Standard text input
    Email,     // Email keyboard with validation patterns
    Password,  // Secure entry with visibility toggle
    Number,    // Numeric keyboard with input filtering
    Search,    // Search-optimized with leading search icon
    Multiline  // Multi-line text areas with configurable line limits
}

enum class TchatInputValidationState {
    None,    // Default neutral state
    Valid,   // Green success border with checkmark
    Invalid  // Red error border with error message
}

enum class TchatInputSize {
    Small,   // 14sp text, compact padding for dense layouts
    Medium,  // 16sp text, standard form field sizing
    Large    // 18sp text, prominent input fields
}

/**
 * Cross-platform input component using expect/actual pattern
 * Platform-specific implementations provide native keyboard handling and styling
 */
@Composable
expect fun TchatInput(
    value: String,
    onValueChange: (String) -> Unit,
    modifier: Modifier = Modifier,
    type: TchatInputType = TchatInputType.Text,
    validationState: TchatInputValidationState = TchatInputValidationState.None,
    size: TchatInputSize = TchatInputSize.Medium,
    placeholder: String = "",
    label: String? = null,
    supportingText: String? = null,
    errorMessage: String? = null,
    enabled: Boolean = true,
    readOnly: Boolean = false,
    leadingIcon: ImageVector? = null,
    trailingIcon: ImageVector? = null,
    onTrailingIconClick: (() -> Unit)? = null,
    maxLines: Int = if (type == TchatInputType.Multiline) Int.MAX_VALUE else 1,
    keyboardActions: KeyboardActions = KeyboardActions.Default,
    focusRequester: FocusRequester = remember { FocusRequester() },
    contentDescription: String? = null
)