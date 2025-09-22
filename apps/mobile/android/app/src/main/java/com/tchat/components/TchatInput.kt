package com.tchat.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.border
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.focus.FocusRequester
import androidx.compose.ui.focus.focusRequester
import androidx.compose.ui.focus.onFocusChanged
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing

/**
 * Input field component following Tchat design system
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TchatInput(
    value: String,
    onValueChange: (String) -> Unit,
    placeholder: String,
    modifier: Modifier = Modifier,
    type: TchatInputType = TchatInputType.Text,
    size: TchatInputSize = TchatInputSize.Medium,
    validationState: TchatValidationState = TchatValidationState.None,
    enabled: Boolean = true,
    leadingIcon: ImageVector? = null,
    trailingIcon: ImageVector? = null,
    onTrailingIconClick: (() -> Unit)? = null,
    keyboardOptions: KeyboardOptions = KeyboardOptions.Default,
    keyboardActions: KeyboardActions = KeyboardActions.Default,
    singleLine: Boolean = true,
    maxLines: Int = if (singleLine) 1 else Int.MAX_VALUE,
    minLines: Int = 1
) {
    var isFocused by remember { mutableStateOf(false) }
    var passwordVisible by remember { mutableStateOf(false) }
    val focusRequester = remember { FocusRequester() }

    // Animation values
    val borderColor by animateColorAsState(
        targetValue = when {
            !enabled -> Colors.border.copy(alpha = 0.5f)
            validationState is TchatValidationState.Invalid -> Colors.borderError
            validationState is TchatValidationState.Valid -> Colors.success
            isFocused -> Colors.borderFocus
            else -> Colors.border
        },
        label = "border_color"
    )

    val borderWidth by animateFloatAsState(
        targetValue = when {
            validationState !is TchatValidationState.None -> 2f
            isFocused -> 2f
            else -> 1f
        },
        label = "border_width"
    )

    Column(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(Spacing.xs)
    ) {
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .border(
                    BorderStroke(borderWidth.dp, borderColor),
                    RoundedCornerShape(Spacing.sm)
                )
        ) {
            OutlinedTextField(
                value = value,
                onValueChange = onValueChange,
                placeholder = {
                    Text(
                        text = placeholder,
                        fontSize = size.fontSize,
                        color = Colors.textTertiary
                    )
                },
                modifier = Modifier
                    .fillMaxWidth()
                    .focusRequester(focusRequester)
                    .onFocusChanged { isFocused = it.isFocused },
                enabled = enabled,
                singleLine = singleLine,
                maxLines = maxLines,
                minLines = minLines,
                textStyle = LocalTextStyle.current.copy(
                    fontSize = size.fontSize,
                    color = if (enabled) Colors.textPrimary else Colors.textDisabled
                ),
                keyboardOptions = keyboardOptions.copy(
                    keyboardType = type.keyboardType,
                    imeAction = keyboardOptions.imeAction
                ),
                keyboardActions = keyboardActions,
                visualTransformation = when (type) {
                    is TchatInputType.Password -> {
                        if (passwordVisible) VisualTransformation.None
                        else PasswordVisualTransformation()
                    }
                    else -> VisualTransformation.None
                },
                leadingIcon = if (leadingIcon != null) {
                    {
                        Icon(
                            imageVector = leadingIcon,
                            contentDescription = null,
                            tint = if (enabled) {
                                if (isFocused) Colors.primary else Colors.textSecondary
                            } else Colors.textDisabled,
                            modifier = Modifier.size(16.dp)
                        )
                    }
                } else null,
                trailingIcon = if (type is TchatInputType.Password || trailingIcon != null) {
                    {
                        when (type) {
                            is TchatInputType.Password -> {
                                IconButton(
                                    onClick = { passwordVisible = !passwordVisible }
                                ) {
                                    Icon(
                                        imageVector = if (passwordVisible) Icons.Default.VisibilityOff
                                        else Icons.Default.Visibility,
                                        contentDescription = if (passwordVisible) "Hide password"
                                        else "Show password",
                                        tint = if (enabled) Colors.textSecondary else Colors.textDisabled,
                                        modifier = Modifier.size(16.dp)
                                    )
                                }
                            }
                            else -> {
                                if (trailingIcon != null) {
                                    IconButton(
                                        onClick = { onTrailingIconClick?.invoke() }
                                    ) {
                                        Icon(
                                            imageVector = trailingIcon,
                                            contentDescription = null,
                                            tint = if (enabled) Colors.textSecondary else Colors.textDisabled,
                                            modifier = Modifier.size(16.dp)
                                        )
                                    }
                                }
                            }
                        }
                    }
                } else null,
                shape = RoundedCornerShape(Spacing.sm),
                colors = OutlinedTextFieldDefaults.colors(
                    focusedBorderColor = Color.Transparent,
                    unfocusedBorderColor = Color.Transparent,
                    disabledBorderColor = Color.Transparent,
                    errorBorderColor = Color.Transparent,
                    focusedContainerColor = if (enabled) Colors.background else Colors.surface.copy(alpha = 0.5f),
                    unfocusedContainerColor = if (enabled) Colors.background else Colors.surface.copy(alpha = 0.5f),
                    disabledContainerColor = Colors.surface.copy(alpha = 0.5f)
                )
            )
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
 * Input type definitions
 */
sealed class TchatInputType(val keyboardType: KeyboardType) {
    object Text : TchatInputType(KeyboardType.Text)
    object Email : TchatInputType(KeyboardType.Email)
    object Password : TchatInputType(KeyboardType.Password)
    object Number : TchatInputType(KeyboardType.Number)
    object Search : TchatInputType(KeyboardType.Text)
    object Multiline : TchatInputType(KeyboardType.Text)
}

/**
 * Input size definitions
 */
enum class TchatInputSize(
    val fontSize: androidx.compose.ui.unit.TextUnit,
    val horizontalPadding: androidx.compose.ui.unit.Dp,
    val verticalPadding: androidx.compose.ui.unit.Dp
) {
    Small(
        fontSize = 14.sp,
        horizontalPadding = Spacing.sm,
        verticalPadding = Spacing.xs
    ),
    Medium(
        fontSize = 16.sp,
        horizontalPadding = Spacing.md,
        verticalPadding = Spacing.sm
    ),
    Large(
        fontSize = 18.sp,
        horizontalPadding = Spacing.lg,
        verticalPadding = Spacing.md
    )
}

/**
 * Validation state definitions
 */
sealed class TchatValidationState {
    object None : TchatValidationState()
    object Valid : TchatValidationState()
    data class Invalid(val message: String) : TchatValidationState()
}

// Preview
@Preview(showBackground = true)
@Composable
fun TchatInputPreview() {
    Column(
        modifier = Modifier.padding(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.md)
    ) {
        TchatInput(
            value = "",
            onValueChange = { },
            placeholder = "Enter your email",
            type = TchatInputType.Email,
            leadingIcon = Icons.Default.Email
        )

        TchatInput(
            value = "",
            onValueChange = { },
            placeholder = "Enter password",
            type = TchatInputType.Password,
            leadingIcon = Icons.Default.Lock
        )

        TchatInput(
            value = "Search...",
            onValueChange = { },
            placeholder = "Search",
            type = TchatInputType.Search,
            leadingIcon = Icons.Default.Search,
            trailingIcon = Icons.Default.Clear,
            onTrailingIconClick = { }
        )

        TchatInput(
            value = "Valid input",
            onValueChange = { },
            placeholder = "Valid input",
            validationState = TchatValidationState.Valid,
            leadingIcon = Icons.Default.CheckCircle
        )

        TchatInput(
            value = "Invalid input",
            onValueChange = { },
            placeholder = "Invalid input",
            validationState = TchatValidationState.Invalid("This field is required"),
            leadingIcon = Icons.Default.Warning
        )

        TchatInput(
            value = "",
            onValueChange = { },
            placeholder = "Disabled input",
            enabled = false
        )

        TchatInput(
            value = "",
            onValueChange = { },
            placeholder = "Type your message here...",
            type = TchatInputType.Multiline,
            singleLine = false,
            maxLines = 6,
            minLines = 3
        )

        Row(
            horizontalArrangement = Arrangement.spacedBy(Spacing.sm)
        ) {
            TchatInput(
                value = "Small",
                onValueChange = { },
                placeholder = "Small",
                size = TchatInputSize.Small,
                modifier = Modifier.weight(1f)
            )
            TchatInput(
                value = "Medium",
                onValueChange = { },
                placeholder = "Medium",
                size = TchatInputSize.Medium,
                modifier = Modifier.weight(1f)
            )
            TchatInput(
                value = "Large",
                onValueChange = { },
                placeholder = "Large",
                size = TchatInputSize.Large,
                modifier = Modifier.weight(1f)
            )
        }
    }
}