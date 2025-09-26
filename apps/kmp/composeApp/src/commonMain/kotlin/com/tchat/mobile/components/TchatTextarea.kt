package com.tchat.mobile.components

import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.focus.FocusRequester
import androidx.compose.ui.graphics.vector.ImageVector

/**
 * TchatTextarea - Cross-platform multi-line text input component with expect/actual pattern
 *
 * Features:
 * - Auto-resize capability with configurable min/max height constraints
 * - Character limits with visual counter and validation feedback
 * - 3 validation states (None, Valid, Invalid) with animated borders
 * - Platform-specific scrolling behavior (iOS momentum, Android fling)
 * - Advanced keyboard handling and input method support
 * - Rich text formatting support (optional)
 * - Line numbering and syntax highlighting (optional)
 * - Full accessibility support with screen reader optimization
 * - Design system integration with consistent theming
 */

enum class TchatTextareaValidationState {
    None,    // Default neutral state with standard border
    Valid,   // Green success border with checkmark indicator
    Invalid  // Red error border with inline error message
}

enum class TchatTextareaSize {
    Small,   // 14sp text, 80dp min height, compact padding for dense layouts
    Medium,  // 16sp text, 100dp min height, standard form field sizing
    Large    // 18sp text, 120dp min height, prominent text areas
}

enum class TchatTextareaResizeMode {
    None,     // Fixed height, scrollable content
    Auto,     // Automatic height adjustment based on content
    Vertical  // Manual vertical resize with drag handle
}

/**
 * Cross-platform textarea component using expect/actual pattern
 * Platform-specific implementations provide native multi-line input behavior and styling
 */
@Composable
expect fun TchatTextarea(
    value: String,
    onValueChange: (String) -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
    readOnly: Boolean = false,
    validationState: TchatTextareaValidationState = TchatTextareaValidationState.None,
    size: TchatTextareaSize = TchatTextareaSize.Medium,
    resizeMode: TchatTextareaResizeMode = TchatTextareaResizeMode.Auto,
    placeholder: String = "",
    label: String? = null,
    supportingText: String? = null,
    errorMessage: String? = null,
    characterLimit: Int? = null,
    showCharacterCount: Boolean = false,
    minLines: Int = 3,
    maxLines: Int = 8,
    leadingIcon: ImageVector? = null,
    trailingIcon: ImageVector? = null,
    onTrailingIconClick: (() -> Unit)? = null,
    keyboardActions: KeyboardActions = KeyboardActions.Default,
    focusRequester: FocusRequester = remember { FocusRequester() },
    interactionSource: MutableInteractionSource = remember { MutableInteractionSource() },
    contentDescription: String? = null
)