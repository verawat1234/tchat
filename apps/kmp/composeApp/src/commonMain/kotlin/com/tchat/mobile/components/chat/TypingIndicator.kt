package com.tchat.mobile.components.chat

import androidx.compose.animation.*
import androidx.compose.animation.core.*
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.scale
import androidx.compose.ui.platform.testTag
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import com.tchat.mobile.models.TypingIndicator

/**
 * Typing indicator components for real-time chat UX
 */
@Composable
fun TypingIndicatorRow(
    typingUsers: List<TypingIndicator>,
    modifier: Modifier = Modifier
) {
    AnimatedVisibility(
        visible = typingUsers.isNotEmpty(),
        enter = slideInVertically(
            initialOffsetY = { it },
            animationSpec = tween(300)
        ) + fadeIn(animationSpec = tween(300)),
        exit = slideOutVertically(
            targetOffsetY = { it },
            animationSpec = tween(300)
        ) + fadeOut(animationSpec = tween(300)),
        modifier = modifier
    ) {
        TypingIndicatorContent(
            typingUsers = typingUsers,
            modifier = Modifier.testTag("typing-indicator-content")
        )
    }
}

@Composable
private fun TypingIndicatorContent(
    typingUsers: List<TypingIndicator>,
    modifier: Modifier = Modifier
) {
    Surface(
        modifier = modifier
            .padding(horizontal = TchatSpacing.md, vertical = TchatSpacing.xs)
            .clip(RoundedCornerShape(16.dp)),
        color = TchatColors.surfaceVariant.copy(alpha = 0.8f)
    ) {
        Row(
            modifier = Modifier.padding(
                horizontal = TchatSpacing.md,
                vertical = TchatSpacing.sm
            ),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Animated typing dots
            AnimatedTypingDots(
                modifier = Modifier.testTag("typing-dots")
            )

            Spacer(modifier = Modifier.width(TchatSpacing.sm))

            // Typing users text
            Text(
                text = getTypingText(typingUsers),
                style = MaterialTheme.typography.bodySmall,
                color = TchatColors.onSurfaceVariant,
                fontStyle = FontStyle.Italic,
                modifier = Modifier.testTag("typing-text")
            )
        }
    }
}

@Composable
private fun AnimatedTypingDots(
    modifier: Modifier = Modifier
) {
    val infiniteTransition = rememberInfiniteTransition(label = "typing_animation")

    Row(
        modifier = modifier,
        horizontalArrangement = Arrangement.spacedBy(2.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        repeat(3) { index ->
            val animationDelay = index * 200
            val scale by infiniteTransition.animateFloat(
                initialValue = 0.5f,
                targetValue = 1.2f,
                animationSpec = infiniteRepeatable(
                    animation = tween(
                        durationMillis = 600,
                        delayMillis = animationDelay,
                        easing = FastOutSlowInEasing
                    ),
                    repeatMode = RepeatMode.Reverse
                ),
                label = "dot_scale_$index"
            )

            Box(
                modifier = Modifier
                    .size(6.dp)
                    .scale(scale)
                    .background(TchatColors.primary, CircleShape)
            )
        }
    }
}

@Composable
fun HeaderTypingIndicator(
    typingUsers: List<TypingIndicator>,
    modifier: Modifier = Modifier
) {
    AnimatedVisibility(
        visible = typingUsers.isNotEmpty(),
        enter = slideInVertically { -it } + fadeIn(),
        exit = slideOutVertically { -it } + fadeOut(),
        modifier = modifier
    ) {
        Row(
            verticalAlignment = Alignment.CenterVertically,
            modifier = Modifier.testTag("header-typing-indicator")
        ) {
            AnimatedTypingDots(
                modifier = Modifier.size(16.dp, 6.dp)
            )

            Spacer(modifier = Modifier.width(4.dp))

            Text(
                text = getShortTypingText(typingUsers),
                style = MaterialTheme.typography.labelSmall,
                color = TchatColors.primary,
                fontStyle = FontStyle.Italic
            )
        }
    }
}

@Composable
fun InlineTypingBubble(
    typingUsers: List<TypingIndicator>,
    modifier: Modifier = Modifier
) {
    AnimatedVisibility(
        visible = typingUsers.isNotEmpty(),
        enter = scaleIn(
            animationSpec = spring(
                dampingRatio = Spring.DampingRatioMediumBouncy,
                stiffness = Spring.StiffnessMedium
            )
        ) + fadeIn(),
        exit = scaleOut(
            animationSpec = tween(200)
        ) + fadeOut(),
        modifier = modifier
    ) {
        Surface(
            modifier = Modifier
                .widthIn(min = 60.dp, max = 120.dp)
                .testTag("inline-typing-bubble"),
            shape = RoundedCornerShape(
                topStart = 16.dp,
                topEnd = 16.dp,
                bottomStart = 4.dp,
                bottomEnd = 16.dp
            ),
            color = TchatColors.surfaceDim
        ) {
            Box(
                modifier = Modifier.padding(TchatSpacing.md),
                contentAlignment = Alignment.Center
            ) {
                AnimatedTypingDots()
            }
        }
    }
}

// Enhanced real-time status bar component with web-like features
@Composable
fun ChatStatusBar(
    isOnline: Boolean,
    lastSeen: String?,
    typingUsers: List<TypingIndicator>,
    modifier: Modifier = Modifier,
    participantCount: Int = 1,
    isActiveNow: Boolean = false,
    customStatus: String? = null
) {
    Column(
        modifier = modifier,
        horizontalAlignment = Alignment.Start
    ) {
        // Main status with enhanced information
        when {
            typingUsers.isNotEmpty() -> {
                HeaderTypingIndicator(typingUsers = typingUsers)
            }
            isActiveNow -> {
                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Box(
                        modifier = Modifier
                            .size(8.dp)
                            .background(TchatColors.primary, CircleShape)
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(
                        text = "Active now",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.primary,
                        fontWeight = FontWeight.Medium
                    )
                }
            }
            isOnline -> {
                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Box(
                        modifier = Modifier
                            .size(8.dp)
                            .background(TchatColors.primary, CircleShape)
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(
                        text = if (participantCount > 1) "Online ($participantCount members)" else "Online",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.primary
                    )
                }
            }
            customStatus != null -> {
                Text(
                    text = customStatus,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant,
                    fontStyle = FontStyle.Italic
                )
            }
            lastSeen != null -> {
                Text(
                    text = formatLastSeenTime(lastSeen),
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )
            }
            else -> {
                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Box(
                        modifier = Modifier
                            .size(8.dp)
                            .background(TchatColors.onSurfaceVariant.copy(alpha = 0.5f), CircleShape)
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(
                        text = if (participantCount > 1) "Offline ($participantCount members)" else "Offline",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            }
        }
    }
}

// Helper function to format last seen time like web interfaces
private fun formatLastSeenTime(lastSeen: String): String {
    return when {
        lastSeen.contains("minute") -> "Last seen ${lastSeen.replace("_", " ")}"
        lastSeen.contains("hour") -> "Last seen ${lastSeen.replace("_", " ")}"
        lastSeen.contains("today") -> "Last seen today"
        lastSeen.contains("yesterday") -> "Last seen yesterday"
        lastSeen.contains("week") -> "Last seen this week"
        lastSeen.contains("online") -> "Last seen recently"
        else -> "Last seen $lastSeen"
    }
}

// Message read indicators
@Composable
fun MessageReadIndicators(
    readBy: List<String>,
    totalParticipants: Int,
    modifier: Modifier = Modifier
) {
    if (readBy.isEmpty()) return

    Row(
        modifier = modifier,
        verticalAlignment = Alignment.CenterVertically
    ) {
        // Show read count if more than one person
        if (readBy.size > 1) {
            Surface(
                modifier = Modifier.size(16.dp),
                shape = CircleShape,
                color = TchatColors.primary
            ) {
                Box(
                    contentAlignment = Alignment.Center
                ) {
                    Text(
                        text = readBy.size.toString(),
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.onPrimary
                    )
                }
            }
        }

        // Show read status text
        Text(
            text = when {
                readBy.size == totalParticipants - 1 -> "Read by all"
                readBy.size > 1 -> "Read by ${readBy.size}"
                readBy.size == 1 -> "Read"
                else -> ""
            },
            style = MaterialTheme.typography.labelSmall,
            color = TchatColors.primary,
            modifier = Modifier.padding(start = if (readBy.size > 1) 4.dp else 0.dp)
        )
    }
}

// Helper functions
private fun getTypingText(typingUsers: List<TypingIndicator>): String {
    return when (typingUsers.size) {
        0 -> ""
        1 -> "${typingUsers.first().userName} is typing..."
        2 -> "${typingUsers[0].userName} and ${typingUsers[1].userName} are typing..."
        else -> "${typingUsers[0].userName} and ${typingUsers.size - 1} others are typing..."
    }
}

private fun getShortTypingText(typingUsers: List<TypingIndicator>): String {
    return when (typingUsers.size) {
        0 -> ""
        1 -> "typing..."
        else -> "${typingUsers.size} typing..."
    }
}