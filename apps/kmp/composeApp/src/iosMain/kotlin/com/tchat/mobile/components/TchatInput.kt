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
 * iOS implementation of TchatInput using Compose Multiplatform
 * Provides iOS-native styling patterns with SwiftUI-inspired design language
 *
 * Features iOS-specific behavior:
 * - iOS-style keyboard handling
 * - Native text field appearance following iOS HIG
 * - iOS-style focus states and animations
 * - Enhanced accessibility for VoiceOver
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

    // iOS HIG-compliant size configuration
    val (height, textStyle, paddingHorizontal, paddingVertical) = when (size) {
        TchatInputSize.Small -> IOSInputSizeConfig(
            36.dp, // iOS minimum touch target
            TchatTypography.typography.bodyMedium,
            14.dp,
            8.dp
        )
        TchatInputSize.Medium -> IOSInputSizeConfig(
            TchatSpacing.inputMinHeight,
            TchatTypography.typography.bodyLarge,
            TchatSpacing.inputPaddingHorizontal,
            TchatSpacing.inputPaddingVertical
        )
        TchatInputSize.Large -> IOSInputSizeConfig(
            54.dp, // Larger for iOS prominence
            TchatTypography.typography.bodyLarge.copy(fontSize = 18.sp),
            22.dp,
            16.dp
        )
    }

    // iOS-style border colors with softer appearance
    val borderColor = when (validationState) {
        TchatInputValidationState.Valid -> TchatColors.success.copy(alpha = 0.8f)
        TchatInputValidationState.Invalid -> TchatColors.error.copy(alpha = 0.8f)
        TchatInputValidationState.None -> if (isFocused) TchatColors.focus.copy(alpha = 0.9f) else TchatColors.outline.copy(alpha = 0.6f)
    }

    // iOS-style smooth animations (longer duration)
    val animatedBorderColor by animateColorAsState(
        targetValue = borderColor,
        animationSpec = tween(250), // iOS spring timing
        label = "ios_border_color"
    )

    val borderWidth by animateFloatAsState(
        targetValue = if (isFocused || validationState != TchatInputValidationState.None) 2.5f else 1.5f, // Thicker iOS borders
        animationSpec = tween(250),
        label = "ios_border_width"
    )

    // iOS-optimized keyboard options
    val keyboardOptions = when (type) {
        TchatInputType.Email -> KeyboardOptions(
            keyboardType = KeyboardType.Email,
            imeAction = ImeAction.Next,
            capitalization = KeyboardCapitalization.None
        )
        TchatInputType.Password -> KeyboardOptions(
            keyboardType = KeyboardType.Password,
            imeAction = ImeAction.Done,
            capitalization = KeyboardCapitalization.None
        )
        TchatInputType.Number -> KeyboardOptions(
            keyboardType = KeyboardType.NumberPassword, // iOS-style number pad
            imeAction = ImeAction.Done
        )
        TchatInputType.Search -> KeyboardOptions(
            keyboardType = KeyboardType.Text,
            imeAction = ImeAction.Search,
            capitalization = KeyboardCapitalization.None
        )
        else -> KeyboardOptions(
            keyboardType = KeyboardType.Text,
            imeAction = if (type == TchatInputType.Multiline) ImeAction.Default else ImeAction.Next,
            capitalization = if (type == TchatInputType.Multiline) KeyboardCapitalization.Sentences else KeyboardCapitalization.None
        )
    }

    // Visual transformation for password
    val visualTransformation = if (type == TchatInputType.Password && !passwordVisible) {
        PasswordVisualTransformation()
    } else {
        VisualTransformation.None
    }

    // iOS-style icon selection (SF Symbol equivalents)
    val effectiveLeadingIcon = leadingIcon ?: when (type) {
        TchatInputType.Email -> Icons.Default.Email
        TchatInputType.Password -> Icons.Default.Lock
        TchatInputType.Search -> Icons.Default.Search
        else -> null
    }

    // iOS-style trailing icons
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
        // iOS-style label with enhanced styling
        label?.let { labelText ->
            Text(
                text = labelText,
                style = TchatTypography.typography.bodyMedium.copy(
                    fontSize = 15.sp // iOS-standard label size
                ),
                color = if (enabled) TchatColors.onSurface.copy(alpha = 0.8f) else TchatColors.disabled.copy(alpha = 0.6f),
                modifier = Modifier.padding(bottom = TchatSpacing.xs * 1.25f)
            )
        }

        // iOS-style input field
        OutlinedTextField(
            value = value,
            onValueChange = { newValue ->
                // iOS-style numeric filtering
                if (type == TchatInputType.Number) {
                    if (newValue.all { it.isDigit() || it == '.' || it == '-' || it == '+' }) {
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
                    shape = RoundedCornerShape(TchatSpacing.inputBorderRadius * 1.3f) // More rounded for iOS
                ),
            enabled = enabled,
            readOnly = readOnly,
            textStyle = textStyle.copy(
                letterSpacing = 0.2.sp // iOS-style letter spacing
            ),
            placeholder = if (placeholder.isNotEmpty()) {
                {
                    Text(
                        text = placeholder,
                        color = TchatColors.onSurfaceVariant.copy(alpha = 0.7f) // Softer iOS placeholder
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
                                TchatInputValidationState.Valid -> TchatColors.success.copy(alpha = 0.8f)
                                TchatInputValidationState.Invalid -> TchatColors.error.copy(alpha = 0.8f)
                                else -> TchatColors.onSurfaceVariant.copy(alpha = 0.7f)
                            }
                        } else TchatColors.disabled.copy(alpha = 0.5f),
                        modifier = Modifier.size(TchatSpacing.iconSize * 1.1f) // Slightly larger for iOS
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
                                TchatInputValidationState.Valid -> TchatColors.success.copy(alpha = 0.8f)
                                TchatInputValidationState.Invalid -> TchatColors.error.copy(alpha = 0.8f)
                                else -> if (enabled) TchatColors.onSurfaceVariant.copy(alpha = 0.7f) else TchatColors.disabled.copy(alpha = 0.5f)
                            },
                            modifier = Modifier.size(TchatSpacing.iconSize * 1.1f)
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
                unfocusedTextColor = TchatColors.onSurface.copy(alpha = 0.9f),
                disabledTextColor = TchatColors.disabled.copy(alpha = 0.6f),
                // iOS-style background
                focusedContainerColor = TchatColors.surface.copy(alpha = 0.5f),
                unfocusedContainerColor = TchatColors.surface.copy(alpha = 0.3f)
            ),
            shape = RoundedCornerShape(TchatSpacing.inputBorderRadius * 1.3f)
        )

        // iOS-style supporting text or error message
        val displayText = if (validationState == TchatInputValidationState.Invalid && errorMessage != null) {
            errorMessage
        } else {
            supportingText
        }

        displayText?.let { text ->
            Text(
                text = text,
                style = TchatTypography.typography.bodySmall.copy(
                    fontSize = 13.sp // iOS-standard supporting text size
                ),
                color = when (validationState) {
                    TchatInputValidationState.Invalid -> TchatColors.error.copy(alpha = 0.8f)
                    TchatInputValidationState.Valid -> TchatColors.success.copy(alpha = 0.8f)
                    else -> TchatColors.onSurfaceVariant.copy(alpha = 0.7f)
                },
                modifier = Modifier.padding(
                    start = TchatSpacing.md,
                    top = TchatSpacing.xs * 1.25f
                )
            )
        }
    }
}

// Helper data class for iOS configuration
private data class IOSInputSizeConfig(
    val height: androidx.compose.ui.unit.Dp,
    val textStyle: androidx.compose.ui.text.TextStyle,
    val horizontalPadding: androidx.compose.ui.unit.Dp,
    val verticalPadding: androidx.compose.ui.unit.Dp
)