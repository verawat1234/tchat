package com.tchat.components

import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.animateOffsetAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.gestures.detectDragGestures
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.input.pointer.pointerInput
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.platform.LocalHapticFeedback
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing
import kotlin.math.abs
import kotlin.math.max
import kotlin.math.min

/**
 * Toggle switch component following Tchat design system
 */
@Composable
fun TchatSwitch(
    checked: Boolean,
    onCheckedChange: (Boolean) -> Unit,
    modifier: Modifier = Modifier,
    label: String? = null,
    description: String? = null,
    size: TchatSwitchSize = TchatSwitchSize.Medium,
    validationState: TchatValidationState = TchatValidationState.None,
    enabled: Boolean = true,
    showLabels: Boolean = false,
    onText: String = "ON",
    offText: String = "OFF"
) {
    val hapticFeedback = LocalHapticFeedback.current
    val density = LocalDensity.current

    var dragOffset by remember { mutableStateOf(0f) }
    var isDragging by remember { mutableStateOf(false) }

    // Animation values
    val trackColor by androidx.compose.animation.animateColorAsState(
        targetValue = when {
            !enabled -> Colors.surface.copy(alpha = 0.5f)
            validationState is TchatValidationState.Invalid -> {
                if (checked) Colors.error else Colors.error.copy(alpha = 0.3f)
            }
            validationState is TchatValidationState.Valid -> {
                if (checked) Colors.success else Colors.success.copy(alpha = 0.3f)
            }
            else -> if (checked) Colors.primary else Colors.border.copy(alpha = 0.3f)
        },
        animationSpec = tween(200),
        label = "track_color"
    )

    val thumbOffset by animateFloatAsState(
        targetValue = if (checked) {
            with(density) { (size.switchWidth - size.thumbSize - 6.dp).toPx() }
        } else {
            with(density) { 3.dp.toPx() }
        } + dragOffset,
        animationSpec = tween(200),
        label = "thumb_offset"
    )

    Column(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(Spacing.xs)
    ) {
        // Main switch row
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .clip(RoundedCornerShape(4.dp))
                .clickable(enabled = enabled && !isDragging) {
                    onCheckedChange(!checked)
                    hapticFeedback.performHapticFeedback(
                        androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                    )
                },
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(Spacing.sm)
        ) {
            // Switch control
            Box(
                modifier = Modifier
                    .size(width = size.switchWidth, height = size.switchHeight)
                    .background(
                        color = trackColor,
                        shape = RoundedCornerShape(size.switchHeight / 2)
                    )
                    .border(
                        width = when {
                            validationState !is TchatValidationState.None -> 2.dp
                            checked -> 0.dp
                            else -> 1.dp
                        },
                        color = when {
                            !enabled -> Colors.border.copy(alpha = 0.5f)
                            validationState is TchatValidationState.Invalid -> Colors.borderError
                            validationState is TchatValidationState.Valid -> Colors.success
                            checked -> Colors.primary
                            else -> Colors.border
                        },
                        shape = RoundedCornerShape(size.switchHeight / 2)
                    )
                    .pointerInput(enabled) {
                        if (!enabled) return@pointerInput

                        detectDragGestures(
                            onDragStart = {
                                isDragging = true
                            },
                            onDragEnd = {
                                isDragging = false
                                val threshold = size.switchWidth.toPx() / 2
                                val shouldToggle = abs(dragOffset) > 20f &&
                                        ((dragOffset > 0 && !checked) ||
                                         (dragOffset < 0 && checked))

                                dragOffset = 0f

                                if (shouldToggle) {
                                    onCheckedChange(!checked)
                                    hapticFeedback.performHapticFeedback(
                                        androidx.compose.ui.hapticfeedback.HapticFeedbackType.LongPress
                                    )
                                }
                            }
                        ) { _, dragAmount ->
                            val maxDrag = (size.switchWidth - size.thumbSize - 6.dp).toPx()
                            dragOffset = max(-3.dp.toPx(), min(maxDrag - 3.dp.toPx(), dragAmount.x))
                        }
                    }
            ) {
                // Labels inside track (if enabled)
                if (showLabels) {
                    Row(
                        modifier = Modifier
                            .fillMaxSize()
                            .padding(horizontal = 6.dp),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Text(
                            text = offText,
                            fontSize = size.labelFontSize,
                            fontWeight = FontWeight.Medium,
                            color = if (checked) Colors.textTertiary.copy(alpha = 0.5f) else Colors.textOnPrimary
                        )

                        Text(
                            text = onText,
                            fontSize = size.labelFontSize,
                            fontWeight = FontWeight.Medium,
                            color = if (checked) Colors.textOnPrimary else Colors.textTertiary.copy(alpha = 0.5f)
                        )
                    }
                }

                // Thumb
                Box(
                    modifier = Modifier
                        .size(size.thumbSize)
                        .offset(x = with(density) { thumbOffset.toDp() })
                        .shadow(2.dp, CircleShape)
                        .background(
                            color = if (enabled) Colors.background else Colors.surface,
                            shape = CircleShape
                        )
                        .align(Alignment.CenterStart)
                )
            }

            // Label and description
            if (label != null || description != null) {
                Column(
                    modifier = Modifier.weight(1f),
                    verticalArrangement = Arrangement.spacedBy(Spacing.xs)
                ) {
                    label?.let { labelText ->
                        Text(
                            text = labelText,
                            fontSize = size.fontSize,
                            color = if (enabled) Colors.textPrimary else Colors.textDisabled
                        )
                    }

                    description?.let { descriptionText ->
                        Text(
                            text = descriptionText,
                            fontSize = size.descriptionFontSize,
                            color = if (enabled) Colors.textSecondary else Colors.textDisabled
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
                modifier = Modifier.padding(start = size.switchWidth + Spacing.sm)
            )
        }
    }
}

/**
 * Switch group component
 */
@Composable
fun TchatSwitchGroup(
    title: String? = null,
    switches: List<TchatSwitchItem>,
    modifier: Modifier = Modifier,
    enabled: Boolean = true
) {
    Column(
        modifier = modifier,
        verticalArrangement = Arrangement.spacedBy(Spacing.md)
    ) {
        // Group title
        title?.let { titleText ->
            Text(
                text = titleText,
                fontSize = 18.sp,
                fontWeight = FontWeight.SemiBold,
                color = if (enabled) Colors.textPrimary else Colors.textDisabled
            )
        }

        // Switches
        Column(
            verticalArrangement = Arrangement.spacedBy(Spacing.sm)
        ) {
            switches.forEach { switchItem ->
                TchatSwitch(
                    checked = switchItem.checked,
                    onCheckedChange = switchItem.onCheckedChange,
                    label = switchItem.label,
                    description = switchItem.description,
                    validationState = switchItem.validationState,
                    enabled = enabled && switchItem.enabled
                )
            }
        }
    }
}

/**
 * Switch item data class for groups
 */
data class TchatSwitchItem(
    val id: String,
    val label: String,
    val description: String? = null,
    val checked: Boolean,
    val onCheckedChange: (Boolean) -> Unit,
    val validationState: TchatValidationState = TchatValidationState.None,
    val enabled: Boolean = true
)

/**
 * Switch size definitions
 */
enum class TchatSwitchSize(
    val switchWidth: androidx.compose.ui.unit.Dp,
    val switchHeight: androidx.compose.ui.unit.Dp,
    val thumbSize: androidx.compose.ui.unit.Dp,
    val fontSize: androidx.compose.ui.unit.TextUnit,
    val descriptionFontSize: androidx.compose.ui.unit.TextUnit,
    val labelFontSize: androidx.compose.ui.unit.TextUnit
) {
    Small(
        switchWidth = 40.dp,
        switchHeight = 24.dp,
        thumbSize = 18.dp,
        fontSize = 14.sp,
        descriptionFontSize = 12.sp,
        labelFontSize = 10.sp
    ),
    Medium(
        switchWidth = 50.dp,
        switchHeight = 30.dp,
        thumbSize = 24.dp,
        fontSize = 16.sp,
        descriptionFontSize = 14.sp,
        labelFontSize = 12.sp
    ),
    Large(
        switchWidth = 60.dp,
        switchHeight = 36.dp,
        thumbSize = 30.dp,
        fontSize = 18.sp,
        descriptionFontSize = 16.sp,
        labelFontSize = 14.sp
    )
}

// Preview
@Preview(showBackground = true)
@Composable
fun TchatSwitchPreview() {
    Column(
        modifier = Modifier.padding(Spacing.md),
        verticalArrangement = Arrangement.spacedBy(Spacing.lg)
    ) {
        // Basic switches
        TchatSwitch(
            checked = true,
            onCheckedChange = { },
            label = "Enable notifications"
        )

        TchatSwitch(
            checked = false,
            onCheckedChange = { },
            label = "Dark mode",
            description = "Use dark appearance throughout the app"
        )

        // Switch with labels
        TchatSwitch(
            checked = true,
            onCheckedChange = { },
            label = "Auto-save",
            description = "Automatically save changes",
            showLabels = true
        )

        // Validation states
        TchatSwitch(
            checked = true,
            onCheckedChange = { },
            label = "Valid setting",
            validationState = TchatValidationState.Valid
        )

        TchatSwitch(
            checked = false,
            onCheckedChange = { },
            label = "Required setting",
            validationState = TchatValidationState.Invalid("This setting is required")
        )

        // Disabled state
        TchatSwitch(
            checked = true,
            onCheckedChange = { },
            label = "Disabled switch",
            description = "This setting cannot be changed",
            enabled = false
        )

        // Different sizes
        Column(verticalArrangement = Arrangement.spacedBy(Spacing.sm)) {
            TchatSwitch(
                checked = true,
                onCheckedChange = { },
                label = "Small switch",
                size = TchatSwitchSize.Small
            )

            TchatSwitch(
                checked = true,
                onCheckedChange = { },
                label = "Medium switch",
                size = TchatSwitchSize.Medium
            )

            TchatSwitch(
                checked = true,
                onCheckedChange = { },
                label = "Large switch",
                size = TchatSwitchSize.Large
            )
        }

        Divider()

        // Switch group
        TchatSwitchGroup(
            title = "Notification Settings",
            switches = listOf(
                TchatSwitchItem(
                    id = "email",
                    label = "Email notifications",
                    description = "Receive notifications via email",
                    checked = true,
                    onCheckedChange = { }
                ),
                TchatSwitchItem(
                    id = "push",
                    label = "Push notifications",
                    description = "Receive push notifications on your device",
                    checked = false,
                    onCheckedChange = { }
                ),
                TchatSwitchItem(
                    id = "sms",
                    label = "SMS notifications",
                    description = "Receive notifications via SMS",
                    checked = true,
                    onCheckedChange = { },
                    validationState = TchatValidationState.Valid
                )
            )
        )
    }
}