package com.tchat.components

import androidx.compose.animation.animateContentSize
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.border
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.BasicTextField
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.focus.onFocusChanged
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.input.TextFieldValue
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing

/**
 * Multi-line text input component following Tchat design system
 */
@Composable
fun TchatTextarea(
    value: String,
    onValueChange: (String) -> Unit,
    modifier: Modifier = Modifier,
    placeholder: String = "Enter text...",
    size: TchatTextareaSize = TchatTextareaSize.Medium,
    validationState: TchatValidationState = TchatValidationState.None,
    enabled: Boolean = true,
    resizeBehavior: TchatTextareaResizeBehavior = TchatTextareaResizeBehavior.AutoResize(
        minLines = 3,
        maxLines = 8
    ),
    characterLimit: Int? = null,
    showCharacterCount: Boolean = false,
    leadingIcon: ImageVector? = null
) {
    var isFocused by remember { mutableStateOf(false) }

    // Character limit handling
    val remainingCharacters = characterLimit?.let { limit ->
        maxOf(0, limit - value.length)
    }

    val isCharacterLimitExceeded = characterLimit?.let { limit ->
        value.length > limit
    } ?: false

    // Dynamic height calculation
    val minHeight = when (resizeBehavior) {
        is TchatTextareaResizeBehavior.Fixed -> resizeBehavior.height
        is TchatTextareaResizeBehavior.AutoResize -> (resizeBehavior.minLines * size.lineHeight.value).dp + (size.verticalPadding * 2)
        TchatTextareaResizeBehavior.Expandable -> size.lineHeight.value.dp + (size.verticalPadding * 2)
    }

    val maxHeight = when (resizeBehavior) {
        is TchatTextareaResizeBehavior.Fixed -> resizeBehavior.height
        is TchatTextareaResizeBehavior.AutoResize -> (resizeBehavior.maxLines * size.lineHeight.value).dp + (size.verticalPadding * 2)
        TchatTextareaResizeBehavior.Expandable -> Dp.Unspecified
    }

    // Border animation
    val borderColor by animateFloatAsState(
        targetValue = when {
            !enabled -> 0.5f
            isCharacterLimitExceeded -> 1f
            validationState is TchatValidationState.Invalid -> 1f
            validationState is TchatValidationState.Valid -> 1f
            isFocused -> 1f
            else -> 1f
        },
        label = "border_color"
    )

    Column(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(Spacing.xs)
    ) {
        // Main textarea field
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .border(
                    BorderStroke(
                        width = when {
                            isCharacterLimitExceeded -> 2.dp
                            validationState !is TchatValidationState.None -> 2.dp
                            isFocused -> 2.dp
                            else -> 1.dp
                        },
                        color = when {
                            !enabled -> Colors.border.copy(alpha = 0.5f)
                            isCharacterLimitExceeded -> Colors.borderError
                            validationState is TchatValidationState.Invalid -> Colors.borderError
                            validationState is TchatValidationState.Valid -> Colors.success
                            isFocused -> Colors.borderFocus
                            else -> Colors.border
                        }
                    ),
                    RoundedCornerShape(Spacing.sm)
                )
                .animateContentSize(),
            verticalAlignment = Alignment.Top
        ) {
            // Leading icon
            leadingIcon?.let { icon ->
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    tint = when {
                        !enabled -> Colors.textDisabled
                        isFocused -> Colors.primary
                        else -> Colors.textSecondary
                    },
                    modifier = Modifier
                        .size(16.dp)
                        .padding(
                            start = size.horizontalPadding,
                            top = size.verticalPadding
                        )
                )
            }

            // Text field
            BasicTextField(
                value = value,
                onValueChange = { newValue ->
                    // Handle character limit
                    if (characterLimit != null && newValue.length > characterLimit) {
                        onValueChange(newValue.take(characterLimit))
                    } else {
                        onValueChange(newValue)
                    }
                },
                modifier = Modifier
                    .weight(1f)
                    .heightIn(
                        min = minHeight,
                        max = if (maxHeight == Dp.Unspecified) Dp.Infinity else maxHeight
                    )
                    .padding(
                        horizontal = size.horizontalPadding,
                        vertical = size.verticalPadding
                    )
                    .onFocusChanged { focusState ->
                        isFocused = focusState.isFocused
                    },
                enabled = enabled,
                textStyle = TextStyle(
                    fontSize = size.fontSize,
                    color = if (enabled) Colors.textPrimary else Colors.textDisabled,
                    lineHeight = size.lineHeight.value.sp
                ),
                decorationBox = { innerTextField ->
                    Box(
                        modifier = Modifier.fillMaxWidth()
                    ) {
                        // Placeholder
                        if (value.isEmpty()) {
                            Text(
                                text = placeholder,
                                fontSize = size.fontSize,
                                color = Colors.textTertiary,
                                modifier = Modifier.alpha(if (enabled) 1f else 0.5f)
                            )
                        }
                        innerTextField()
                    }
                }
            )
        }

        // Bottom section (character count)
        if (showCharacterCount || characterLimit != null) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.End
            ) {
                when {
                    remainingCharacters != null -> {
                        Text(
                            text = "$remainingCharacters remaining",
                            fontSize = 12.sp,
                            color = if (isCharacterLimitExceeded) Colors.error else Colors.textTertiary
                        )
                    }
                    showCharacterCount -> {
                        Text(
                            text = "${value.length} characters",
                            fontSize = 12.sp,
                            color = Colors.textTertiary
                        )
                    }
                }
            }
        }

        // Validation message
        if (validationState is TchatValidationState.Invalid) {
            Text(
                text = validationState.message,
                fontSize = 12.sp,
                color = Colors.error,
                modifier = Modifier.padding(horizontal = Spacing.xs)
            )
        }
    }
}

/**
 * Textarea resize behavior definitions
 */
sealed class TchatTextareaResizeBehavior {
    data class Fixed(val height: Dp) : TchatTextareaResizeBehavior()
    data class AutoResize(val minLines: Int, val maxLines: Int) : TchatTextareaResizeBehavior()
    object Expandable : TchatTextareaResizeBehavior()
}

/**
 * Textarea size definitions
 */
enum class TchatTextareaSize(
    val fontSize: androidx.compose.ui.unit.TextUnit,
    val lineHeight: androidx.compose.ui.unit.Dp,
    val horizontalPadding: androidx.compose.ui.unit.Dp,
    val verticalPadding: androidx.compose.ui.unit.Dp
) {
    Small(
        fontSize = 14.sp,
        lineHeight = 18.dp,
        horizontalPadding = Spacing.sm,
        verticalPadding = Spacing.xs
    ),
    Medium(
        fontSize = 16.sp,
        lineHeight = 22.dp,
        horizontalPadding = Spacing.md,
        verticalPadding = Spacing.sm
    ),
    Large(
        fontSize = 18.sp,
        lineHeight = 26.dp,
        horizontalPadding = Spacing.lg,
        verticalPadding = Spacing.md
    )
}

/**
 * Convenience composables for common use cases
 */
@Composable
fun TchatTextareaFixed(
    value: String,
    onValueChange: (String) -> Unit,
    height: Dp,
    modifier: Modifier = Modifier,
    placeholder: String = "Enter text...",
    size: TchatTextareaSize = TchatTextareaSize.Medium,
    validationState: TchatValidationState = TchatValidationState.None,
    enabled: Boolean = true,
    characterLimit: Int? = null,
    showCharacterCount: Boolean = false,
    leadingIcon: ImageVector? = null
) {
    TchatTextarea(
        value = value,
        onValueChange = onValueChange,
        modifier = modifier,
        placeholder = placeholder,
        size = size,
        validationState = validationState,
        enabled = enabled,
        resizeBehavior = TchatTextareaResizeBehavior.Fixed(height),
        characterLimit = characterLimit,
        showCharacterCount = showCharacterCount,
        leadingIcon = leadingIcon
    )
}

@Composable
fun TchatTextareaExpandable(
    value: String,
    onValueChange: (String) -> Unit,
    modifier: Modifier = Modifier,
    placeholder: String = "Enter text...",
    size: TchatTextareaSize = TchatTextareaSize.Medium,
    validationState: TchatValidationState = TchatValidationState.None,
    enabled: Boolean = true,
    characterLimit: Int? = null,
    showCharacterCount: Boolean = false,
    leadingIcon: ImageVector? = null
) {
    TchatTextarea(
        value = value,
        onValueChange = onValueChange,
        modifier = modifier,
        placeholder = placeholder,
        size = size,
        validationState = validationState,
        enabled = enabled,
        resizeBehavior = TchatTextareaResizeBehavior.Expandable,
        characterLimit = characterLimit,
        showCharacterCount = showCharacterCount,
        leadingIcon = leadingIcon
    )
}

// Preview
@Preview(showBackground = true)
@Composable
fun TchatTextareaPreview() {
    Column(
        modifier = Modifier.padding(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.md)
    ) {
        // Auto-resize textarea
        TchatTextarea(
            value = "This is a sample text that demonstrates the auto-resize functionality of the textarea component.",
            onValueChange = { },
            placeholder = "Enter your message...",
            resizeBehavior = TchatTextareaResizeBehavior.AutoResize(minLines = 3, maxLines = 6),
            leadingIcon = Icons.Default.Edit
        )

        // Fixed height textarea
        TchatTextareaFixed(
            value = "Fixed height textarea",
            onValueChange = { },
            height = 100.dp,
            placeholder = "Fixed height (100dp)",
            characterLimit = 500,
            showCharacterCount = true
        )

        // Expandable textarea
        TchatTextareaExpandable(
            value = "Expandable textarea that grows with content",
            onValueChange = { },
            placeholder = "Type to see expansion...",
            characterLimit = 280,
            showCharacterCount = true
        )

        // Validation states
        TchatTextarea(
            value = "Valid content",
            onValueChange = { },
            placeholder = "Valid textarea",
            validationState = TchatValidationState.Valid,
            resizeBehavior = TchatTextareaResizeBehavior.AutoResize(minLines = 2, maxLines = 4),
            leadingIcon = Icons.Default.CheckCircle
        )

        TchatTextarea(
            value = "",
            onValueChange = { },
            placeholder = "Required field",
            validationState = TchatValidationState.Invalid("This field is required"),
            resizeBehavior = TchatTextareaResizeBehavior.AutoResize(minLines = 2, maxLines = 4),
            leadingIcon = Icons.Default.Warning
        )

        // Disabled state
        TchatTextarea(
            value = "Disabled textarea content",
            onValueChange = { },
            placeholder = "Disabled",
            enabled = false,
            resizeBehavior = TchatTextareaResizeBehavior.AutoResize(minLines = 3, maxLines = 5)
        )

        // Different sizes
        Column(verticalArrangement = Arrangement.spacedBy(Spacing.sm)) {
            TchatTextarea(
                value = "Small size",
                onValueChange = { },
                placeholder = "Small textarea",
                size = TchatTextareaSize.Small,
                resizeBehavior = TchatTextareaResizeBehavior.AutoResize(minLines = 2, maxLines = 3)
            )

            TchatTextarea(
                value = "Large size",
                onValueChange = { },
                placeholder = "Large textarea",
                size = TchatTextareaSize.Large,
                resizeBehavior = TchatTextareaResizeBehavior.AutoResize(minLines = 2, maxLines = 3)
            )
        }

        // Character limit example
        TchatTextarea(
            value = "This demonstrates character limit functionality",
            onValueChange = { },
            placeholder = "Tweet-like input (280 chars)",
            characterLimit = 280,
            showCharacterCount = true,
            resizeBehavior = TchatTextareaResizeBehavior.AutoResize(minLines = 3, maxLines = 6),
            leadingIcon = Icons.Default.AlternateEmail
        )
    }
}