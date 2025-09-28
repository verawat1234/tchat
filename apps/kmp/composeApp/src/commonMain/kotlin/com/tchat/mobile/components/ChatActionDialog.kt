package com.tchat.mobile.components

import androidx.compose.animation.*
import androidx.compose.animation.core.*
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.semantics.*
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.Dialog
import androidx.compose.ui.window.DialogProperties
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.models.ChatParticipant
import com.tchat.mobile.models.ParticipantStatus

/**
 * ChatActionDialog - Confirmation dialog for chat actions
 *
 * Features:
 * - Video and voice call confirmations
 * - Participant information display
 * - Smooth animations and transitions
 * - Accessibility support
 * - Loading states for call initiation
 */

enum class ChatActionType(
    val title: String,
    val icon: ImageVector,
    val description: String
) {
    VIDEO_CALL(
        title = "Start Video Call",
        icon = Icons.Default.VideoCall,
        description = "Start a video call with"
    ),
    VOICE_CALL(
        title = "Start Voice Call",
        icon = Icons.Default.Call,
        description = "Start a voice call with"
    )
}

@Composable
fun ChatActionDialog(
    isVisible: Boolean,
    actionType: ChatActionType,
    participants: List<ChatParticipant>,
    onConfirm: () -> Unit,
    onDismiss: () -> Unit,
    isLoading: Boolean = false,
    modifier: Modifier = Modifier
) {
    AnimatedVisibility(
        visible = isVisible,
        enter = fadeIn(animationSpec = tween(200)) +
                scaleIn(animationSpec = tween(200, easing = FastOutSlowInEasing)),
        exit = fadeOut(animationSpec = tween(150)) +
               scaleOut(animationSpec = tween(150, easing = FastOutSlowInEasing))
    ) {
        Dialog(
            onDismissRequest = onDismiss,
            properties = DialogProperties(
                dismissOnBackPress = !isLoading,
                dismissOnClickOutside = !isLoading
            )
        ) {
            Surface(
                modifier = modifier,
                shape = RoundedCornerShape(20.dp),
                color = TchatColors.surface,
                shadowElevation = 8.dp
            ) {
                Column(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(24.dp),
                    horizontalAlignment = Alignment.CenterHorizontally
                ) {
                    // Action icon with animated background
                    ActionIcon(
                        actionType = actionType,
                        isLoading = isLoading
                    )

                    Spacer(modifier = Modifier.height(24.dp))

                    // Title
                    Text(
                        text = actionType.title,
                        style = MaterialTheme.typography.headlineSmall,
                        fontWeight = FontWeight.SemiBold,
                        color = TchatColors.onSurface,
                        textAlign = TextAlign.Center
                    )

                    Spacer(modifier = Modifier.height(8.dp))

                    // Description with participants
                    ParticipantsSection(
                        actionType = actionType,
                        participants = participants
                    )

                    Spacer(modifier = Modifier.height(32.dp))

                    // Action buttons
                    ActionButtons(
                        onConfirm = onConfirm,
                        onDismiss = onDismiss,
                        isLoading = isLoading,
                        actionType = actionType
                    )
                }
            }
        }
    }
}

@Composable
private fun ActionIcon(
    actionType: ChatActionType,
    isLoading: Boolean
) {
    val backgroundColor = when (actionType) {
        ChatActionType.VIDEO_CALL -> TchatColors.primary
        ChatActionType.VOICE_CALL -> TchatColors.success
    }

    val infiniteTransition = rememberInfiniteTransition()
    val pulseScale by infiniteTransition.animateFloat(
        initialValue = 1f,
        targetValue = if (isLoading) 1.1f else 1f,
        animationSpec = infiniteRepeatable(
            animation = tween(1000, easing = FastOutSlowInEasing),
            repeatMode = RepeatMode.Reverse
        )
    )

    Box(
        modifier = Modifier
            .size(80.dp)
            .clip(CircleShape)
            .background(backgroundColor.copy(alpha = 0.1f)),
        contentAlignment = Alignment.Center
    ) {
        Box(
            modifier = Modifier
                .size(60.dp)
                .clip(CircleShape)
                .background(backgroundColor)
                .then(
                    if (isLoading) {
                        Modifier.animateContentSize()
                    } else {
                        Modifier
                    }
                ),
            contentAlignment = Alignment.Center
        ) {
            if (isLoading) {
                CircularProgressIndicator(
                    modifier = Modifier.size(24.dp),
                    strokeWidth = 2.dp,
                    color = TchatColors.onPrimary
                )
            } else {
                Icon(
                    imageVector = actionType.icon,
                    contentDescription = actionType.title,
                    modifier = Modifier.size(28.dp),
                    tint = TchatColors.onPrimary
                )
            }
        }
    }
}

@Composable
private fun ParticipantsSection(
    actionType: ChatActionType,
    participants: List<ChatParticipant>
) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            text = actionType.description,
            style = MaterialTheme.typography.bodyLarge,
            color = TchatColors.onSurfaceVariant,
            textAlign = TextAlign.Center
        )

        Spacer(modifier = Modifier.height(12.dp))

        // Show participants
        if (participants.isNotEmpty()) {
            when {
                participants.size == 1 -> {
                    // Single participant
                    SingleParticipant(participants.first())
                }
                participants.size <= 3 -> {
                    // Multiple participants (small group)
                    MultipleParticipants(participants)
                }
                else -> {
                    // Large group
                    LargeGroupParticipants(participants)
                }
            }
        }
    }
}

@Composable
private fun SingleParticipant(participant: ChatParticipant) {
    Row(
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        // Avatar
        TchatAvatar(
            name = participant.name,
            imageUrl = participant.avatar,
            size = TchatAvatarSize.MD,
            modifier = Modifier.size(40.dp)
        )

        // Name and status
        Column {
            Text(
                text = participant.name,
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.Medium,
                color = TchatColors.onSurface
            )

            if (participant.status != ParticipantStatus.OFFLINE) {
                Text(
                    text = participant.status.name.lowercase().replaceFirstChar { it.uppercase() },
                    style = MaterialTheme.typography.bodySmall,
                    color = when (participant.status) {
                        ParticipantStatus.ONLINE -> TchatColors.success
                        ParticipantStatus.BUSY -> TchatColors.warning
                        ParticipantStatus.AWAY -> TchatColors.warning
                        else -> TchatColors.onSurfaceVariant
                    }
                )
            }
        }
    }
}

@Composable
private fun MultipleParticipants(participants: List<ChatParticipant>) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        // Participant names
        Text(
            text = participants.joinToString(", ") { it.name },
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.Medium,
            color = TchatColors.onSurface,
            textAlign = TextAlign.Center
        )

        // Participant count
        Text(
            text = "${participants.size} participants",
            style = MaterialTheme.typography.bodySmall,
            color = TchatColors.onSurfaceVariant
        )
    }
}

@Composable
private fun LargeGroupParticipants(participants: List<ChatParticipant>) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        // First few participants
        val displayParticipants = participants.take(3)
        Text(
            text = displayParticipants.joinToString(", ") { it.name } +
                  if (participants.size > 3) " and ${participants.size - 3} others" else "",
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.Medium,
            color = TchatColors.onSurface,
            textAlign = TextAlign.Center
        )

        // Total count
        Text(
            text = "${participants.size} participants",
            style = MaterialTheme.typography.bodySmall,
            color = TchatColors.onSurfaceVariant
        )
    }
}

@Composable
private fun ActionButtons(
    onConfirm: () -> Unit,
    onDismiss: () -> Unit,
    isLoading: Boolean,
    actionType: ChatActionType
) {
    Row(
        modifier = Modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        // Cancel button
        OutlinedButton(
            onClick = onDismiss,
            modifier = Modifier
                .weight(1f)
                .height(48.dp),
            enabled = !isLoading,
            colors = ButtonDefaults.outlinedButtonColors(
                contentColor = TchatColors.onSurface,
                disabledContentColor = TchatColors.disabled
            ),
            border = androidx.compose.foundation.BorderStroke(
                width = 1.dp,
                color = if (isLoading) TchatColors.disabled else TchatColors.outline
            )
        ) {
            Text(
                text = "Cancel",
                style = MaterialTheme.typography.labelLarge,
                fontWeight = FontWeight.Medium
            )
        }

        // Confirm button
        Button(
            onClick = onConfirm,
            modifier = Modifier
                .weight(1f)
                .height(48.dp)
                .semantics {
                    contentDescription = "Start ${actionType.title.lowercase()}"
                },
            enabled = !isLoading,
            colors = ButtonDefaults.buttonColors(
                containerColor = when (actionType) {
                    ChatActionType.VIDEO_CALL -> TchatColors.primary
                    ChatActionType.VOICE_CALL -> TchatColors.success
                },
                contentColor = TchatColors.onPrimary,
                disabledContainerColor = TchatColors.disabled,
                disabledContentColor = TchatColors.disabled
            )
        ) {
            if (isLoading) {
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    CircularProgressIndicator(
                        modifier = Modifier.size(16.dp),
                        strokeWidth = 2.dp,
                        color = TchatColors.onPrimary
                    )
                    Text(
                        text = "Starting...",
                        style = MaterialTheme.typography.labelLarge,
                        fontWeight = FontWeight.Medium
                    )
                }
            } else {
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    Icon(
                        imageVector = actionType.icon,
                        contentDescription = null,
                        modifier = Modifier.size(18.dp)
                    )
                    Text(
                        text = "Start Call",
                        style = MaterialTheme.typography.labelLarge,
                        fontWeight = FontWeight.Medium
                    )
                }
            }
        }
    }
}