package com.tchat.mobile.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material.icons.filled.Error
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.focus.FocusRequester
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardCapitalization
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography

/**
 * Android implementation of TchatTextarea using Material3 OutlinedTextField
 * Provides native Material Design multi-line text field with comprehensive theming
 */
@Composable
actual fun TchatTextarea(
    value: String,
    onValueChange: (String) -> Unit,
    modifier: Modifier,
    enabled: Boolean,
    readOnly: Boolean,
    validationState: TchatTextareaValidationState,
    size: TchatTextareaSize,
    resizeMode: TchatTextareaResizeMode,
    placeholder: String,
    label: String?,
    supportingText: String?,
    errorMessage: String?,
    characterLimit: Int?,
    showCharacterCount: Boolean,
    minLines: Int,
    maxLines: Int,
    leadingIcon: ImageVector?,
    trailingIcon: ImageVector?,
    onTrailingIconClick: (() -> Unit)?,
    keyboardActions: KeyboardActions,
    focusRequester: FocusRequester,
    interactionSource: MutableInteractionSource,
    contentDescription: String?
) {
    val minHeight = when (size) {
        TchatTextareaSize.Small -> 80.dp
        TchatTextareaSize.Medium -> 100.dp
        TchatTextareaSize.Large -> 120.dp
    }

    val textStyle = when (size) {
        TchatTextareaSize.Small -> TchatTypography.typography.bodySmall
        TchatTextareaSize.Medium -> TchatTypography.typography.bodyMedium
        TchatTextareaSize.Large -> TchatTypography.typography.bodyLarge
    }

    val iconSize = when (size) {
        TchatTextareaSize.Small -> 16.dp
        TchatTextareaSize.Medium -> 20.dp
        TchatTextareaSize.Large -> 24.dp
    }

    // Character count validation
    val isOverLimit = characterLimit != null && value.length > characterLimit
    val effectiveValidationState = when {
        isOverLimit -> TchatTextareaValidationState.Invalid
        else -> validationState
    }

    val borderColor by animateColorAsState(
        targetValue = when (effectiveValidationState) {
            TchatTextareaValidationState.Valid -> TchatColors.success
            TchatTextareaValidationState.Invalid -> TchatColors.error
            TchatTextareaValidationState.None -> TchatColors.outline
        },
        animationSpec = tween(300),
        label = "BorderColor"
    )

    val borderAlpha by animateFloatAsState(
        targetValue = when (effectiveValidationState) {
            TchatTextareaValidationState.None -> 1f
            else -> 1f
        },
        animationSpec = tween(300),
        label = "BorderAlpha"
    )

    val textFieldColors = OutlinedTextFieldDefaults.colors(
        focusedBorderColor = borderColor,
        unfocusedBorderColor = borderColor.copy(alpha = borderAlpha),
        disabledBorderColor = TchatColors.outline.copy(alpha = 0.12f),
        focusedLabelColor = borderColor,
        unfocusedLabelColor = TchatColors.onSurface.copy(alpha = 0.6f),
        disabledLabelColor = TchatColors.onSurface.copy(alpha = 0.38f),
        cursorColor = TchatColors.primary,
        focusedTextColor = TchatColors.onSurface,
        unfocusedTextColor = TchatColors.onSurface,
        disabledTextColor = TchatColors.onSurface.copy(alpha = 0.38f)
    )

    Column(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(4.dp)
    ) {
        // Main textarea field
        OutlinedTextField(
            value = value,
            onValueChange = { newValue ->
                // Apply character limit if set
                val filteredValue = if (characterLimit != null) {
                    newValue.take(characterLimit)
                } else {
                    newValue
                }
                onValueChange(filteredValue)
            },
            modifier = Modifier
                .fillMaxWidth()
                .then(
                    if (resizeMode == TchatTextareaResizeMode.None) {
                        Modifier.heightIn(min = minHeight)
                    } else {
                        Modifier
                    }
                ),
            enabled = enabled,
            readOnly = readOnly,
            textStyle = textStyle,
            label = label?.let { labelText ->
                {
                    Text(
                        text = labelText,
                        style = TchatTypography.typography.bodySmall
                    )
                }
            },
            placeholder = if (placeholder.isNotEmpty()) {
                {
                    Text(
                        text = placeholder,
                        style = textStyle,
                        color = TchatColors.onSurface.copy(alpha = 0.6f)
                    )
                }
            } else null,
            leadingIcon = leadingIcon?.let { icon ->
                {
                    Icon(
                        imageVector = icon,
                        contentDescription = null,
                        modifier = Modifier.size(iconSize),
                        tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f)
                    )
                }
            },
            trailingIcon = {
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(4.dp)
                ) {
                    // Validation icon
                    when (effectiveValidationState) {
                        TchatTextareaValidationState.Valid -> {
                            Icon(
                                imageVector = Icons.Default.CheckCircle,
                                contentDescription = "Valid",
                                modifier = Modifier.size(iconSize),
                                tint = TchatColors.success
                            )
                        }
                        TchatTextareaValidationState.Invalid -> {
                            Icon(
                                imageVector = Icons.Default.Error,
                                contentDescription = "Invalid",
                                modifier = Modifier.size(iconSize),
                                tint = TchatColors.error
                            )
                        }
                        TchatTextareaValidationState.None -> {}
                    }

                    // Custom trailing icon
                    trailingIcon?.let { icon ->
                        IconButton(
                            onClick = { onTrailingIconClick?.invoke() },
                            enabled = enabled
                        ) {
                            Icon(
                                imageVector = icon,
                                contentDescription = null,
                                modifier = Modifier.size(iconSize),
                                tint = if (enabled) TchatColors.onSurface else TchatColors.onSurface.copy(alpha = 0.38f)
                            )
                        }
                    }
                }
            },
            supportingText = {
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    // Supporting text or error message
                    Column(
                        modifier = Modifier.weight(1f)
                    ) {
                        when {
                            effectiveValidationState == TchatTextareaValidationState.Invalid && errorMessage != null -> {
                                Text(
                                    text = errorMessage,
                                    style = TchatTypography.typography.bodySmall,
                                    color = TchatColors.error
                                )
                            }
                            isOverLimit && characterLimit != null -> {
                                Text(
                                    text = "Character limit exceeded",
                                    style = TchatTypography.typography.bodySmall,
                                    color = TchatColors.error
                                )
                            }
                            supportingText != null -> {
                                Text(
                                    text = supportingText,
                                    style = TchatTypography.typography.bodySmall,
                                    color = TchatColors.onSurface.copy(alpha = 0.6f)
                                )
                            }
                        }
                    }

                    // Character count
                    if (showCharacterCount || characterLimit != null) {
                        val characterCountText = if (characterLimit != null) {
                            "${value.length}/$characterLimit"
                        } else {
                            "${value.length}"
                        }

                        Text(
                            text = characterCountText,
                            style = TchatTypography.typography.bodySmall,
                            color = when {
                                isOverLimit -> TchatColors.error
                                characterLimit != null && value.length > characterLimit * 0.8f -> TchatColors.warning
                                else -> TchatColors.onSurface.copy(alpha = 0.6f)
                            },
                            textAlign = TextAlign.End
                        )
                    }
                }
            },
            minLines = minLines,
            maxLines = if (resizeMode == TchatTextareaResizeMode.Auto) Int.MAX_VALUE else maxLines,
            keyboardOptions = KeyboardOptions(
                keyboardType = KeyboardType.Text,
                imeAction = ImeAction.Default,
                capitalization = KeyboardCapitalization.Sentences
            ),
            keyboardActions = keyboardActions,
            colors = textFieldColors,
            interactionSource = interactionSource
        )
    }
}