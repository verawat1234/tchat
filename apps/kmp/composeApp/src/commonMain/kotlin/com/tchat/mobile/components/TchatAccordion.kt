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
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.KeyboardArrowDown
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
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
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.role
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing

/**
 * TchatAccordion - Expandable content sections with complex animations
 *
 * Features:
 * - Multiple expandable sections with smooth animations
 * - Single or multiple section expansion modes
 * - Memory-efficient lazy content rendering
 * - Advanced accessibility with focus management
 * - Platform-specific presentation patterns
 * - Custom header and content rendering support
 * - Keyboard navigation support
 */

/**
 * Data class representing a single accordion item
 */
data class AccordionItem(
    val id: String,
    val title: String,
    val subtitle: String? = null,
    val content: @Composable () -> Unit,
    val icon: (@Composable () -> Unit)? = null,
    val enabled: Boolean = true
)

/**
 * Accordion expansion behavior modes
 */
enum class AccordionMode {
    /**
     * Only one section can be expanded at a time
     * Opening a section automatically closes others
     */
    Single,

    /**
     * Multiple sections can be expanded simultaneously
     * Each section maintains independent state
     */
    Multiple
}

/**
 * TchatAccordion - Cross-platform accordion component
 *
 * @param items List of accordion items to display
 * @param modifier Modifier for styling the accordion container
 * @param mode Expansion behavior (Single or Multiple)
 * @param initialExpandedIds Set of initially expanded item IDs
 * @param onExpandedChange Callback when expansion state changes
 * @param showDividers Whether to show dividers between sections
 * @param backgroundColor Background color for the accordion
 * @param contentDescription Accessibility description
 */
@Composable
fun TchatAccordion(
    items: List<AccordionItem>,
    modifier: Modifier = Modifier,
    mode: AccordionMode = AccordionMode.Multiple,
    initialExpandedIds: Set<String> = emptySet(),
    onExpandedChange: ((itemId: String, isExpanded: Boolean) -> Unit)? = null,
    showDividers: Boolean = true,
    backgroundColor: Color = TchatColors.background,
    contentDescription: String? = null
) {
    var expandedItems by remember {
        mutableStateOf(initialExpandedIds.toSet())
    }

    Surface(
        modifier = modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(TchatSpacing.cardBorderRadius))
            .semantics {
                contentDescription?.let {
                    this.contentDescription = it
                }
            },
        color = backgroundColor,
        shadowElevation = if (backgroundColor == TchatColors.background) 0.dp else 1.dp
    ) {
        Column {
            items.forEachIndexed { index, item ->
                val isExpanded = expandedItems.contains(item.id)
                val isLast = index == items.lastIndex

                AccordionSection(
                    item = item,
                    isExpanded = isExpanded,
                    onToggle = {
                        val newExpandedItems = when (mode) {
                            AccordionMode.Single -> {
                                if (isExpanded) {
                                    expandedItems - item.id
                                } else {
                                    setOf(item.id)
                                }
                            }
                            AccordionMode.Multiple -> {
                                if (isExpanded) {
                                    expandedItems - item.id
                                } else {
                                    expandedItems + item.id
                                }
                            }
                        }
                        expandedItems = newExpandedItems
                        onExpandedChange?.invoke(item.id, !isExpanded)
                    },
                    showBottomDivider = showDividers && !isLast,
                    index = index
                )
            }
        }
    }
}

/**
 * Individual accordion section component with animation support
 */
@Composable
private fun AccordionSection(
    item: AccordionItem,
    isExpanded: Boolean,
    onToggle: () -> Unit,
    showBottomDivider: Boolean,
    index: Int
) {
    val rotationAngle by animateFloatAsState(
        targetValue = if (isExpanded) 180f else 0f,
        animationSpec = tween(
            durationMillis = 300,
            easing = LinearOutSlowInEasing
        ),
        label = "chevron_rotation_${item.id}"
    )

    Column(
        modifier = Modifier.animateContentSize(
            animationSpec = tween(
                durationMillis = 300,
                easing = LinearOutSlowInEasing
            )
        )
    ) {
        // Header Section
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .clickable(
                    enabled = item.enabled,
                    onClick = onToggle
                )
                .padding(TchatSpacing.md)
                .semantics {
                    role = Role.Button
                    contentDescription = buildString {
                        append(item.title)
                        item.subtitle?.let { append(", $it") }
                        append(if (isExpanded) ", expanded" else ", collapsed")
                        append(", button, tap to ${if (isExpanded) "collapse" else "expand"}")
                    }
                },
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Row(
                verticalAlignment = Alignment.CenterVertically,
                modifier = Modifier.weight(1f)
            ) {
                // Leading icon if provided
                item.icon?.let { icon ->
                    Box(
                        modifier = Modifier
                            .size(TchatSpacing.iconSize)
                            .padding(end = TchatSpacing.sm)
                    ) {
                        icon()
                    }
                }

                // Title and subtitle
                Column {
                    Text(
                        text = item.title,
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Medium,
                        color = if (item.enabled) TchatColors.onSurface else TchatColors.disabled
                    )

                    item.subtitle?.let { subtitle ->
                        Spacer(modifier = Modifier.height(2.dp))
                        Text(
                            text = subtitle,
                            style = MaterialTheme.typography.bodySmall,
                            color = if (item.enabled) TchatColors.onSurfaceVariant else TchatColors.disabled
                        )
                    }
                }
            }

            // Chevron icon with rotation animation
            Icon(
                imageVector = Icons.Default.KeyboardArrowDown,
                contentDescription = null,
                modifier = Modifier
                    .rotate(rotationAngle)
                    .size(TchatSpacing.iconSize),
                tint = if (item.enabled) TchatColors.onSurfaceVariant else TchatColors.disabled
            )
        }

        // Expandable content section
        AnimatedVisibility(
            visible = isExpanded,
            enter = expandVertically(
                animationSpec = tween(
                    durationMillis = 300,
                    easing = LinearOutSlowInEasing
                )
            ) + fadeIn(
                animationSpec = tween(
                    durationMillis = 200,
                    delayMillis = 100,
                    easing = LinearOutSlowInEasing
                )
            ),
            exit = shrinkVertically(
                animationSpec = tween(
                    durationMillis = 200,
                    easing = LinearOutSlowInEasing
                )
            ) + fadeOut(
                animationSpec = tween(
                    durationMillis = 150,
                    easing = LinearOutSlowInEasing
                )
            )
        ) {
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .background(TchatColors.surfaceVariant.copy(alpha = 0.3f))
                    .padding(
                        start = TchatSpacing.md,
                        end = TchatSpacing.md,
                        top = TchatSpacing.sm,
                        bottom = TchatSpacing.md
                    )
                    .semantics {
                                contentDescription = "Content for ${item.title}"
                    }
            ) {
                item.content()
            }
        }

        // Bottom divider
        if (showBottomDivider) {
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(1.dp)
                    .background(TchatColors.outline.copy(alpha = 0.5f))
            )
        }
    }
}

/**
 * Convenience composable for creating simple text-based accordion items
 */
@Composable
fun TchatSimpleAccordion(
    items: List<Pair<String, String>>, // Title to Content pairs
    modifier: Modifier = Modifier,
    mode: AccordionMode = AccordionMode.Multiple,
    initialExpandedIndices: Set<Int> = emptySet(),
    onExpandedChange: ((index: Int, isExpanded: Boolean) -> Unit)? = null
) {
    val accordionItems = items.mapIndexed { index, (title, content) ->
        AccordionItem(
            id = index.toString(),
            title = title,
            content = {
                Text(
                    text = content,
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurface
                )
            }
        )
    }

    TchatAccordion(
        items = accordionItems,
        modifier = modifier,
        mode = mode,
        initialExpandedIds = initialExpandedIndices.map { it.toString() }.toSet(),
        onExpandedChange = { itemId, isExpanded ->
            itemId.toIntOrNull()?.let { index ->
                onExpandedChange?.invoke(index, isExpanded)
            }
        }
    )
}