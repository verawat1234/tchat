package com.tchat.mobile.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.border
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.focus.FocusRequester
import androidx.compose.ui.focus.focusRequester
import androidx.compose.ui.focus.onFocusChanged
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.text.input.*
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import com.tchat.mobile.designsystem.TchatTypography

/**
 * Android implementation of TchatInput using Jetpack Compose and Material3
 * Provides Material Design 3 styling with advanced input field features
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
actual fun TchatInput(
    value: String,
    onValueChange: (String) -> Unit,
    modifier: Modifier,
    type: TchatInputType,
    validationState: TchatInputValidationState,
    size: TchatInputSize,
    placeholder: String,
    label: String?,
    supportingText: String?,
    errorMessage: String?,
    enabled: Boolean,
    readOnly: Boolean,
    leadingIcon: ImageVector?,
    trailingIcon: ImageVector?,
    onTrailingIconClick: (() -> Unit)?,
    maxLines: Int,
    keyboardActions: KeyboardActions,
    focusRequester: FocusRequester,
    contentDescription: String?
) {
    var isFocused by remember { mutableStateOf(false) }
    var passwordVisible by remember { mutableStateOf(false) }

    // Size configuration
    val (height, textStyle, paddingHorizontal, paddingVertical) = when (size) {
        TchatInputSize.Small -> AndroidInputSizeConfig(
            32.dp,
            TchatTypography.typography.bodyMedium,
            12.dp,
            6.dp
        )
        TchatInputSize.Medium -> AndroidInputSizeConfig(
            TchatSpacing.inputMinHeight,
            TchatTypography.typography.bodyLarge,
            TchatSpacing.inputPaddingHorizontal,
            TchatSpacing.inputPaddingVertical
        )
        TchatInputSize.Large -> AndroidInputSizeConfig(
            52.dp,
            TchatTypography.typography.bodyLarge.copy(fontSize = 18.sp),
            20.dp,
            14.dp
        )
    }

    // Border color based on validation state and focus
    val borderColor = when (validationState) {
        TchatInputValidationState.Valid -> TchatColors.success
        TchatInputValidationState.Invalid -> TchatColors.error
        TchatInputValidationState.None -> if (isFocused) TchatColors.focus else TchatColors.outline
    }

    // Animated border color and width
    val animatedBorderColor by animateColorAsState(
        targetValue = borderColor,
        animationSpec = tween(150),
        label = "border_color"
    )

    val borderWidth by animateFloatAsState(
        targetValue = if (isFocused || validationState != TchatInputValidationState.None) 2f else 1f,
        animationSpec = tween(150),
        label = "border_width"
    )

    // Keyboard options based on input type
    val keyboardOptions = when (type) {
        TchatInputType.Email -> KeyboardOptions(
            keyboardType = KeyboardType.Email,
            imeAction = ImeAction.Next
        )
        TchatInputType.Password -> KeyboardOptions(
            keyboardType = KeyboardType.Password,
            imeAction = ImeAction.Done
        )
        TchatInputType.Number -> KeyboardOptions(
            keyboardType = KeyboardType.Number,
            imeAction = ImeAction.Done
        )
        TchatInputType.Search -> KeyboardOptions(
            keyboardType = KeyboardType.Text,
            imeAction = ImeAction.Search
        )
        else -> KeyboardOptions(
            keyboardType = KeyboardType.Text,
            imeAction = if (type == TchatInputType.Multiline) ImeAction.Default else ImeAction.Next
        )
    }

    // Visual transformation for password
    val visualTransformation = if (type == TchatInputType.Password && !passwordVisible) {
        PasswordVisualTransformation()
    } else {
        VisualTransformation.None
    }

    // Leading icon based on type
    val effectiveLeadingIcon = leadingIcon ?: when (type) {
        TchatInputType.Email -> Icons.Default.Email
        TchatInputType.Password -> Icons.Default.Lock
        TchatInputType.Search -> Icons.Default.Search
        else -> null
    }

    // Trailing icon with password visibility toggle
    val effectiveTrailingIcon = when {
        type == TchatInputType.Password -> if (passwordVisible) Icons.Default.Close else Icons.Default.Info
        validationState == TchatInputValidationState.Valid -> Icons.Default.CheckCircle
        validationState == TchatInputValidationState.Invalid -> Icons.Default.Warning
        else -> trailingIcon
    }

    val effectiveTrailingIconClick: (() -> Unit)? = when {
        type == TchatInputType.Password -> ({ passwordVisible = !passwordVisible })
        else -> onTrailingIconClick
    }

    Column(
        modifier = modifier.semantics {
            this.contentDescription = contentDescription ?: when (type) {
                TchatInputType.Email -> "Email input field"
                TchatInputType.Password -> "Password input field"
                TchatInputType.Search -> "Search input field"
                TchatInputType.Number -> "Number input field"
                TchatInputType.Multiline -> "Multiline text input"
                else -> "Text input field"
            }
        }
    ) {
        // Label
        label?.let { labelText ->
            Text(
                text = labelText,
                style = TchatTypography.typography.bodyMedium,
                color = if (enabled) TchatColors.onSurface else TchatColors.disabled,
                modifier = Modifier.padding(bottom = TchatSpacing.xs)
            )
        }

        // Input field
        OutlinedTextField(
            value = value,
            onValueChange = { newValue ->
                // Filter numeric input
                if (type == TchatInputType.Number) {
                    if (newValue.all { it.isDigit() || it == '.' || it == '-' }) {
                        onValueChange(newValue)
                    }
                } else {
                    onValueChange(newValue)
                }
            },
            modifier = Modifier
                .fillMaxWidth()
                .let { baseModifier ->
                    if (type != TchatInputType.Multiline) {
                        baseModifier.heightIn(min = height)
                    } else baseModifier
                }
                .focusRequester(focusRequester)
                .onFocusChanged { focusState ->
                    isFocused = focusState.isFocused
                }
                .border(
                    width = borderWidth.dp,
                    color = animatedBorderColor,
                    shape = RoundedCornerShape(TchatSpacing.inputBorderRadius)
                ),
            enabled = enabled,
            readOnly = readOnly,
            textStyle = textStyle,
            placeholder = if (placeholder.isNotEmpty()) {
                {
                    Text(
                        text = placeholder,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            } else null,
            leadingIcon = effectiveLeadingIcon?.let { icon ->
                {
                    Icon(
                        imageVector = icon,
                        contentDescription = null,
                        tint = if (enabled) {
                            when (validationState) {
                                TchatInputValidationState.Valid -> TchatColors.success
                                TchatInputValidationState.Invalid -> TchatColors.error
                                else -> TchatColors.onSurfaceVariant
                            }
                        } else TchatColors.disabled,
                        modifier = Modifier.size(TchatSpacing.iconSize)
                    )
                }
            },
            trailingIcon = effectiveTrailingIcon?.let { icon ->
                {
                    IconButton(
                        onClick = { effectiveTrailingIconClick?.invoke() },
                        enabled = enabled
                    ) {
                        Icon(
                            imageVector = icon,
                            contentDescription = when {
                                type == TchatInputType.Password && passwordVisible -> "Hide password"
                                type == TchatInputType.Password && !passwordVisible -> "Show password"
                                validationState == TchatInputValidationState.Valid -> "Valid input"
                                validationState == TchatInputValidationState.Invalid -> "Invalid input"
                                else -> null
                            },
                            tint = when (validationState) {
                                TchatInputValidationState.Valid -> TchatColors.success
                                TchatInputValidationState.Invalid -> TchatColors.error
                                else -> if (enabled) TchatColors.onSurfaceVariant else TchatColors.disabled
                            },
                            modifier = Modifier.size(TchatSpacing.iconSize)
                        )
                    }
                }
            },
            visualTransformation = visualTransformation,
            keyboardOptions = keyboardOptions,
            keyboardActions = keyboardActions,
            singleLine = type != TchatInputType.Multiline,
            maxLines = maxLines,
            colors = OutlinedTextFieldDefaults.colors(
                focusedBorderColor = Color.Transparent, // We handle border ourselves
                unfocusedBorderColor = Color.Transparent,
                errorBorderColor = Color.Transparent,
                disabledBorderColor = Color.Transparent,
                focusedTextColor = TchatColors.onSurface,
                unfocusedTextColor = TchatColors.onSurface,
                disabledTextColor = TchatColors.disabled
            ),
            shape = RoundedCornerShape(TchatSpacing.inputBorderRadius)
        )

        // Supporting text or error message
        val displayText = if (validationState == TchatInputValidationState.Invalid && errorMessage != null) {
            errorMessage
        } else {
            supportingText
        }

        displayText?.let { text ->
            Text(
                text = text,
                style = TchatTypography.typography.bodySmall,
                color = when (validationState) {
                    TchatInputValidationState.Invalid -> TchatColors.error
                    TchatInputValidationState.Valid -> TchatColors.success
                    else -> TchatColors.onSurfaceVariant
                },
                modifier = Modifier.padding(
                    start = TchatSpacing.md,
                    top = TchatSpacing.xs
                )
            )
        }
    }
}

// Helper data class for Android configuration
private data class AndroidInputSizeConfig(
    val height: androidx.compose.ui.unit.Dp,
    val textStyle: androidx.compose.ui.text.TextStyle,
    val horizontalPadding: androidx.compose.ui.unit.Dp,
    val verticalPadding: androidx.compose.ui.unit.Dp
)