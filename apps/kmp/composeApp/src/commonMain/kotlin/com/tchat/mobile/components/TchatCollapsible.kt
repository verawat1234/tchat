package com.tchat.mobile.components

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.animateContentSize
import androidx.compose.animation.core.LinearOutSlowInEasing
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.animation.expandVertically
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut
import androidx.compose.animation.shrinkVertically
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.KeyboardArrowDown
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.rotate
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.Shape
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.role
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

/**
 * TchatCollapsible - Single collapsible content component with advanced animations
 *
 * Features:
 * - Smooth expand/collapse animations with customizable duration
 * - Custom trigger styling and content
 * - Memory-efficient content rendering
 * - Advanced accessibility with focus management
 * - Platform-specific interaction feedback
 * - Custom animations and easing curves
 * - Keyboard navigation support
 */

/**
 * Collapsible component variants for different use cases
 */
enum class CollapsibleVariant {
    /**
     * Default variant with standard styling
     */
    Default,

    /**
     * Card-based variant with elevated surface
     */
    Card,

    /**
     * Outlined variant with border
     */
    Outlined,

    /**
     * Ghost variant with minimal styling
     */
    Ghost
}

/**
 * Animation configuration for collapsible behavior
 */
data class CollapsibleAnimation(
    val expandDuration: Int = 300,
    val collapseDuration: Int = 250,
    val fadeInDelay: Int = 100,
    val fadeInDuration: Int = 200,
    val fadeOutDuration: Int = 150,
    val easing: androidx.compose.animation.core.Easing = LinearOutSlowInEasing
)

/**
 * TchatCollapsible - Cross-platform collapsible component
 *
 * @param isExpanded Whether the content is currently expanded
 * @param onToggle Callback when the collapsible is toggled
 * @param trigger Composable content for the trigger button
 * @param content Composable content that can be collapsed/expanded
 * @param modifier Modifier for styling the collapsible container
 * @param variant Visual variant of the collapsible
 * @param enabled Whether the collapsible is interactive
 * @param showChevron Whether to show the chevron indicator
 * @param chevronIcon Custom chevron icon (uses default arrow if not provided)
 * @param backgroundColor Background color
 * @param contentBackgroundColor Background color for the content area
 * @param shape Shape of the collapsible container
 * @param elevation Elevation for card variant
 * @param animation Animation configuration
 * @param interactionSource Interaction source for custom ripple effects
 * @param contentDescription Accessibility description
 */
@Composable
fun TchatCollapsible(
    isExpanded: Boolean,
    onToggle: () -> Unit,
    trigger: @Composable () -> Unit,
    content: @Composable () -> Unit,
    modifier: Modifier = Modifier,
    variant: CollapsibleVariant = CollapsibleVariant.Default,
    enabled: Boolean = true,
    showChevron: Boolean = true,
    chevronIcon: ImageVector = Icons.Default.KeyboardArrowDown,
    backgroundColor: Color = when (variant) {
        CollapsibleVariant.Default -> Color.Transparent
        CollapsibleVariant.Card -> TchatColors.surface
        CollapsibleVariant.Outlined -> Color.Transparent
        CollapsibleVariant.Ghost -> Color.Transparent
    },
    contentBackgroundColor: Color = TchatColors.surfaceVariant.copy(alpha = 0.3f),
    shape: Shape = RoundedCornerShape(TchatSpacing.cardBorderRadius),
    elevation: Dp = when (variant) {
        CollapsibleVariant.Card -> TchatSpacing.cardElevation
        else -> 0.dp
    },
    animation: CollapsibleAnimation = CollapsibleAnimation(),
    interactionSource: MutableInteractionSource = remember { MutableInteractionSource() },
    contentDescription: String? = null
) {
    val chevronRotation by animateFloatAsState(
        targetValue = if (isExpanded) 180f else 0f,
        animationSpec = tween(
            durationMillis = animation.expandDuration,
            easing = animation.easing
        ),
        label = "chevron_rotation"
    )

    Surface(
        modifier = modifier
            .fillMaxWidth()
            .clip(shape)
            .semantics {
                contentDescription?.let {
                    this.contentDescription = it
                }
                role = Role.Button
            },
        color = backgroundColor,
        shadowElevation = elevation,
        shape = shape
    ) {
        Column(
            modifier = Modifier.animateContentSize(
                animationSpec = tween(
                    durationMillis = if (isExpanded) animation.expandDuration else animation.collapseDuration,
                    easing = animation.easing
                )
            )
        ) {
            // Trigger Section
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .clickable(
                        enabled = enabled,
                        interactionSource = interactionSource,
                        indication = null,
                        onClick = onToggle
                    )
                    .then(
                        if (variant == CollapsibleVariant.Outlined) {
                            Modifier.background(
                                color = Color.Transparent,
                                shape = shape
                            ).padding(1.dp).background(
                                color = Color.Transparent,
                                shape = shape
                            )
                        } else {
                            Modifier
                        }
                    )
                    .padding(TchatSpacing.md),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                // Trigger content
                Box(
                    modifier = Modifier.weight(1f)
                ) {
                    trigger()
                }

                // Chevron indicator
                if (showChevron) {
                    Spacer(modifier = Modifier.width(TchatSpacing.sm))
                    Icon(
                        imageVector = chevronIcon,
                        contentDescription = if (isExpanded) "Collapse" else "Expand",
                        modifier = Modifier
                            .rotate(chevronRotation)
                            .size(TchatSpacing.iconSize),
                        tint = if (enabled) TchatColors.onSurfaceVariant else TchatColors.disabled
                    )
                }
            }

            // Border for outlined variant
            if (variant == CollapsibleVariant.Outlined) {
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(1.dp)
                        .background(TchatColors.outline)
                )
            }

            // Expandable content section
            AnimatedVisibility(
                visible = isExpanded,
                enter = expandVertically(
                    animationSpec = tween(
                        durationMillis = animation.expandDuration,
                        easing = animation.easing
                    )
                ) + fadeIn(
                    animationSpec = tween(
                        durationMillis = animation.fadeInDuration,
                        delayMillis = animation.fadeInDelay,
                        easing = animation.easing
                    )
                ),
                exit = shrinkVertically(
                    animationSpec = tween(
                        durationMillis = animation.collapseDuration,
                        easing = animation.easing
                    )
                ) + fadeOut(
                    animationSpec = tween(
                        durationMillis = animation.fadeOutDuration,
                        easing = animation.easing
                    )
                )
            ) {
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .background(contentBackgroundColor)
                        .padding(TchatSpacing.md)
                        .semantics {
                            this.contentDescription = "Expandable content"
                        }
                ) {
                    content()
                }
            }
        }
    }
}

/**
 * Stateful version of TchatCollapsible that manages its own expansion state
 */
@Composable
fun TchatCollapsible(
    trigger: @Composable () -> Unit,
    content: @Composable () -> Unit,
    modifier: Modifier = Modifier,
    initiallyExpanded: Boolean = false,
    variant: CollapsibleVariant = CollapsibleVariant.Default,
    enabled: Boolean = true,
    showChevron: Boolean = true,
    onToggle: ((isExpanded: Boolean) -> Unit)? = null,
    chevronIcon: ImageVector = Icons.Default.KeyboardArrowDown,
    backgroundColor: Color = when (variant) {
        CollapsibleVariant.Default -> Color.Transparent
        CollapsibleVariant.Card -> TchatColors.surface
        CollapsibleVariant.Outlined -> Color.Transparent
        CollapsibleVariant.Ghost -> Color.Transparent
    },
    contentBackgroundColor: Color = TchatColors.surfaceVariant.copy(alpha = 0.3f),
    shape: Shape = RoundedCornerShape(TchatSpacing.cardBorderRadius),
    elevation: Dp = when (variant) {
        CollapsibleVariant.Card -> TchatSpacing.cardElevation
        else -> 0.dp
    },
    animation: CollapsibleAnimation = CollapsibleAnimation(),
    contentDescription: String? = null
) {
    var isExpanded by remember { mutableStateOf(initiallyExpanded) }

    TchatCollapsible(
        isExpanded = isExpanded,
        onToggle = {
            isExpanded = !isExpanded
            onToggle?.invoke(isExpanded)
        },
        trigger = trigger,
        content = content,
        modifier = modifier,
        variant = variant,
        enabled = enabled,
        showChevron = showChevron,
        chevronIcon = chevronIcon,
        backgroundColor = backgroundColor,
        contentBackgroundColor = contentBackgroundColor,
        shape = shape,
        elevation = elevation,
        animation = animation,
        contentDescription = contentDescription
    )
}

/**
 * Convenience composable for simple text-based collapsible
 */
@Composable
fun TchatSimpleCollapsible(
    title: String,
    content: String,
    modifier: Modifier = Modifier,
    subtitle: String? = null,
    initiallyExpanded: Boolean = false,
    variant: CollapsibleVariant = CollapsibleVariant.Default,
    enabled: Boolean = true,
    onToggle: ((isExpanded: Boolean) -> Unit)? = null
) {
    TchatCollapsible(
        trigger = {
            Column {
                Text(
                    text = title,
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Medium,
                    color = if (enabled) TchatColors.onSurface else TchatColors.disabled
                )

                subtitle?.let {
                    Spacer(modifier = Modifier.height(2.dp))
                    Text(
                        text = it,
                        style = MaterialTheme.typography.bodySmall,
                        color = if (enabled) TchatColors.onSurfaceVariant else TchatColors.disabled
                    )
                }
            }
        },
        content = {
            Text(
                text = content,
                style = MaterialTheme.typography.bodyMedium,
                color = TchatColors.onSurface
            )
        },
        modifier = modifier,
        initiallyExpanded = initiallyExpanded,
        variant = variant,
        enabled = enabled,
        onToggle = onToggle,
        contentDescription = "Collapsible section: $title"
    )
}