package com.tchat.mobile.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.selection.toggleable
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.scale
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.semantics.Role
import androidx.compose.foundation.BorderStroke
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatTypography

/**
 * Android implementation of TchatToggle using Material3 styling
 * Provides native Material Design toggle with comprehensive theming and animations
 */
@Composable
actual fun TchatToggle(
    pressed: Boolean,
    onPressedChange: ((Boolean) -> Unit)?,
    modifier: Modifier,
    enabled: Boolean,
    size: TchatToggleSize,
    variant: TchatToggleVariant,
    text: String?,
    icon: ImageVector?,
    leadingIcon: ImageVector?,
    trailingIcon: ImageVector?,
    interactionSource: MutableInteractionSource,
    contentDescription: String?
) {
    val buttonHeight = when (size) {
        TchatToggleSize.Small -> 32.dp
        TchatToggleSize.Medium -> 40.dp
        TchatToggleSize.Large -> 48.dp
    }

    val textStyle = when (size) {
        TchatToggleSize.Small -> TchatTypography.typography.bodySmall
        TchatToggleSize.Medium -> TchatTypography.typography.bodyMedium
        TchatToggleSize.Large -> TchatTypography.typography.bodyLarge
    }

    val iconSize = when (size) {
        TchatToggleSize.Small -> 16.dp
        TchatToggleSize.Medium -> 20.dp
        TchatToggleSize.Large -> 24.dp
    }

    val horizontalPadding = when (size) {
        TchatToggleSize.Small -> 12.dp
        TchatToggleSize.Medium -> 16.dp
        TchatToggleSize.Large -> 20.dp
    }

    // Animation for background color
    val backgroundColor by animateColorAsState(
        targetValue = when (variant) {
            TchatToggleVariant.Default -> if (pressed) {
                TchatColors.primary
            } else {
                TchatColors.surface
            }
            TchatToggleVariant.Outline -> if (pressed) {
                TchatColors.primary.copy(alpha = 0.12f)
            } else {
                Color.Transparent
            }
            TchatToggleVariant.Text -> Color.Transparent
        },
        animationSpec = tween(200),
        label = "BackgroundColor"
    )

    // Animation for text/icon color
    val contentColor by animateColorAsState(
        targetValue = when (variant) {
            TchatToggleVariant.Default -> if (pressed) {
                TchatColors.onPrimary
            } else {
                TchatColors.onSurface
            }
            TchatToggleVariant.Outline -> if (pressed) {
                TchatColors.primary
            } else {
                TchatColors.onSurface
            }
            TchatToggleVariant.Text -> if (pressed) {
                TchatColors.primary
            } else {
                TchatColors.onSurface
            }
        },
        animationSpec = tween(200),
        label = "ContentColor"
    )

    // Animation for border
    val borderColor by animateColorAsState(
        targetValue = when (variant) {
            TchatToggleVariant.Outline -> if (pressed) {
                TchatColors.primary
            } else {
                TchatColors.outline
            }
            else -> Color.Transparent
        },
        animationSpec = tween(200),
        label = "BorderColor"
    )

    // Scale animation for press feedback
    val scale by animateFloatAsState(
        targetValue = if (pressed) 0.95f else 1f,
        animationSpec = tween(100),
        label = "Scale"
    )

    Surface(
        modifier = modifier
            .height(buttonHeight)
            .scale(scale)
            .clip(RoundedCornerShape(8.dp))
            .toggleable(
                value = pressed,
                enabled = enabled,
                role = Role.Button,
                interactionSource = interactionSource,
                indication = null,
                onValueChange = onPressedChange ?: {}
            ),
        color = backgroundColor,
        shape = RoundedCornerShape(8.dp),
        border = if (variant == TchatToggleVariant.Outline) {
            BorderStroke(1.dp, borderColor)
        } else null
    ) {
        Row(
            modifier = Modifier
                .fillMaxHeight()
                .padding(horizontal = horizontalPadding),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            // Leading icon
            leadingIcon?.let { leadingIconVector ->
                Icon(
                    imageVector = leadingIconVector,
                    contentDescription = null,
                    modifier = Modifier.size(iconSize),
                    tint = if (enabled) contentColor else contentColor.copy(alpha = 0.38f)
                )
            }

            // Main icon (if no text)
            if (text == null) {
                icon?.let { iconVector ->
                    Icon(
                        imageVector = iconVector,
                        contentDescription = contentDescription,
                        modifier = Modifier.size(iconSize),
                        tint = if (enabled) contentColor else contentColor.copy(alpha = 0.38f)
                    )
                }
            }

            // Text
            text?.let { textContent ->
                Text(
                    text = textContent,
                    style = textStyle.copy(
                        fontWeight = if (pressed) FontWeight.Medium else FontWeight.Normal
                    ),
                    color = if (enabled) contentColor else contentColor.copy(alpha = 0.38f)
                )
            }

            // Trailing icon
            trailingIcon?.let { trailingIconVector ->
                Icon(
                    imageVector = trailingIconVector,
                    contentDescription = null,
                    modifier = Modifier.size(iconSize),
                    tint = if (enabled) contentColor else contentColor.copy(alpha = 0.38f)
                )
            }
        }
    }
}

/**
 * Android implementation of TchatToggleGroup using Material3 styling
 */
@Composable
actual fun TchatToggleGroup(
    options: List<String>,
    selectedOptions: Set<String>,
    onSelectionChange: (Set<String>) -> Unit,
    modifier: Modifier,
    enabled: Boolean,
    size: TchatToggleSize,
    variant: TchatToggleVariant,
    allowMultipleSelection: Boolean,
    allowEmptySelection: Boolean,
    maxSelections: Int?,
    interactionSource: MutableInteractionSource
) {
    Row(
        modifier = modifier,
        horizontalArrangement = Arrangement.spacedBy(2.dp) // Material3 segmented button spacing
    ) {
        options.forEachIndexed { index, option ->
            val isSelected = selectedOptions.contains(option)
            val isFirstItem = index == 0
            val isLastItem = index == options.size - 1

            val shape = when {
                options.size == 1 -> RoundedCornerShape(8.dp)
                isFirstItem -> RoundedCornerShape(topStart = 8.dp, bottomStart = 8.dp)
                isLastItem -> RoundedCornerShape(topEnd = 8.dp, bottomEnd = 8.dp)
                else -> RoundedCornerShape(0.dp)
            }

            TchatToggle(
                pressed = isSelected,
                onPressedChange = { _ ->
                    handleToggleGroupSelection(
                        option = option,
                        currentSelection = selectedOptions,
                        allowMultipleSelection = allowMultipleSelection,
                        allowEmptySelection = allowEmptySelection,
                        maxSelections = maxSelections,
                        onSelectionChange = onSelectionChange
                    )
                },
                modifier = Modifier.weight(1f),
                enabled = enabled,
                size = size,
                variant = variant,
                text = option,
                interactionSource = remember { MutableInteractionSource() }
            )
        }
    }
}

private fun handleToggleGroupSelection(
    option: String,
    currentSelection: Set<String>,
    allowMultipleSelection: Boolean,
    allowEmptySelection: Boolean,
    maxSelections: Int?,
    onSelectionChange: (Set<String>) -> Unit
) {
    val newSelection = if (currentSelection.contains(option)) {
        // Deselecting
        if (allowEmptySelection || currentSelection.size > 1) {
            currentSelection - option
        } else {
            currentSelection // Don't allow empty selection if not allowed
        }
    } else {
        // Selecting
        if (allowMultipleSelection) {
            if (maxSelections == null || currentSelection.size < maxSelections) {
                currentSelection + option
            } else {
                currentSelection // Max selections reached
            }
        } else {
            setOf(option) // Single selection mode
        }
    }

    onSelectionChange(newSelection)
}