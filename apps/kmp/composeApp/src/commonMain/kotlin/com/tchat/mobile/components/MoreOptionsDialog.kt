package com.tchat.mobile.components

import androidx.compose.animation.*
import androidx.compose.animation.core.*
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.semantics.*
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.Dialog
import androidx.compose.ui.window.DialogProperties
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.models.ChatSession

/**
 * MoreOptionsDialog - Chat options menu dialog
 *
 * Features:
 * - Comprehensive chat management options
 * - Dynamic options based on chat state
 * - Smooth animations and transitions
 * - Accessibility support
 * - Destructive action confirmations
 */

data class ChatOption(
    val id: String,
    val title: String,
    val subtitle: String? = null,
    val icon: ImageVector,
    val enabled: Boolean = true,
    val destructive: Boolean = false,
    val action: () -> Unit
)

@Composable
fun MoreOptionsDialog(
    isVisible: Boolean,
    chatSession: ChatSession?,
    onDismiss: () -> Unit,
    onPinChat: () -> Unit,
    onMuteChat: () -> Unit,
    onArchiveChat: () -> Unit,
    onBlockChat: () -> Unit,
    onDeleteChat: () -> Unit,
    onClearHistory: () -> Unit,
    onExportChat: () -> Unit,
    onChatInfo: () -> Unit,
    onReportChat: () -> Unit,
    modifier: Modifier = Modifier
) {
    var showDeleteConfirmation by remember { mutableStateOf(false) }
    var showClearConfirmation by remember { mutableStateOf(false) }

    // Generate options based on chat state
    val options = remember(chatSession) {
        generateChatOptions(
            chatSession = chatSession,
            onPinChat = onPinChat,
            onMuteChat = onMuteChat,
            onArchiveChat = onArchiveChat,
            onBlockChat = onBlockChat,
            onDeleteChat = { showDeleteConfirmation = true },
            onClearHistory = { showClearConfirmation = true },
            onExportChat = onExportChat,
            onChatInfo = onChatInfo,
            onReportChat = onReportChat
        )
    }

    AnimatedVisibility(
        visible = isVisible,
        enter = fadeIn(animationSpec = tween(200)) +
                slideInVertically(animationSpec = tween(200)) { it / 2 },
        exit = fadeOut(animationSpec = tween(150)) +
               slideOutVertically(animationSpec = tween(150)) { it / 2 }
    ) {
        Dialog(
            onDismissRequest = onDismiss,
            properties = DialogProperties(
                dismissOnBackPress = true,
                dismissOnClickOutside = true
            )
        ) {
            Surface(
                modifier = modifier.fillMaxWidth(),
                shape = RoundedCornerShape(16.dp),
                color = TchatColors.surface,
                shadowElevation = 8.dp
            ) {
                Column(
                    modifier = Modifier.fillMaxWidth()
                ) {
                    // Header
                    OptionsHeader(
                        chatSession = chatSession,
                        onDismiss = onDismiss
                    )

                    // Options list
                    LazyColumn(
                        modifier = Modifier.fillMaxWidth(),
                        contentPadding = PaddingValues(bottom = 16.dp)
                    ) {
                        items(options) { option ->
                            OptionItem(
                                option = option,
                                onClick = {
                                    option.action()
                                    if (!option.destructive) {
                                        onDismiss()
                                    }
                                }
                            )
                        }
                    }
                }
            }
        }
    }

    // Confirmation dialogs
    DeleteConfirmationDialog(
        isVisible = showDeleteConfirmation,
        onConfirm = {
            onDeleteChat()
            showDeleteConfirmation = false
            onDismiss()
        },
        onDismiss = { showDeleteConfirmation = false }
    )

    ClearHistoryConfirmationDialog(
        isVisible = showClearConfirmation,
        onConfirm = {
            onClearHistory()
            showClearConfirmation = false
            onDismiss()
        },
        onDismiss = { showClearConfirmation = false }
    )
}

@Composable
private fun OptionsHeader(
    chatSession: ChatSession?,
    onDismiss: () -> Unit
) {
    Surface(
        modifier = Modifier.fillMaxWidth(),
        color = TchatColors.surfaceVariant,
        shape = RoundedCornerShape(topStart = 16.dp, topEnd = 16.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Column(
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = "Chat Options",
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.SemiBold,
                    color = TchatColors.onSurfaceVariant
                )

                chatSession?.name?.let { name ->
                    Text(
                        text = name,
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.onSurfaceVariant,
                        maxLines = 1,
                        overflow = TextOverflow.Ellipsis
                    )
                }
            }

            IconButton(
                onClick = onDismiss,
                modifier = Modifier.semantics {
                    contentDescription = "Close options"
                }
            ) {
                Icon(
                    Icons.Default.Close,
                    contentDescription = "Close",
                    tint = TchatColors.onSurfaceVariant
                )
            }
        }
    }
}

@Composable
private fun OptionItem(
    option: ChatOption,
    onClick: () -> Unit
) {
    val animatedBackgroundColor by animateColorAsState(
        targetValue = Color.Transparent,
        animationSpec = tween(200)
    )

    Surface(
        modifier = Modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(0.dp))
            .clickable(enabled = option.enabled) { onClick() }
            .semantics {
                role = Role.Button
                contentDescription = option.title
                if (!option.enabled) {
                    disabled()
                }
            },
        color = animatedBackgroundColor
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 20.dp, vertical = 16.dp),
            horizontalArrangement = Arrangement.spacedBy(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Icon
            Icon(
                imageVector = option.icon,
                contentDescription = null,
                modifier = Modifier.size(24.dp),
                tint = when {
                    !option.enabled -> TchatColors.disabled
                    option.destructive -> TchatColors.error
                    else -> TchatColors.onSurface
                }
            )

            // Text content
            Column(
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = option.title,
                    style = MaterialTheme.typography.bodyLarge,
                    fontWeight = FontWeight.Medium,
                    color = when {
                        !option.enabled -> TchatColors.disabled
                        option.destructive -> TchatColors.error
                        else -> TchatColors.onSurface
                    }
                )

                option.subtitle?.let { subtitle ->
                    Text(
                        text = subtitle,
                        style = MaterialTheme.typography.bodySmall,
                        color = when {
                            !option.enabled -> TchatColors.disabled
                            option.destructive -> TchatColors.error.copy(alpha = 0.7f)
                            else -> TchatColors.onSurfaceVariant
                        }
                    )
                }
            }

            // Trailing indicator for destructive actions
            if (option.destructive) {
                Icon(
                    Icons.Default.Warning,
                    contentDescription = "Destructive action",
                    modifier = Modifier.size(16.dp),
                    tint = TchatColors.error.copy(alpha = 0.6f)
                )
            }
        }
    }
}

@Composable
private fun DeleteConfirmationDialog(
    isVisible: Boolean,
    onConfirm: () -> Unit,
    onDismiss: () -> Unit
) {
    TchatDialog(
        isVisible = isVisible,
        onDismissRequest = onDismiss,
        variant = DialogVariant.Confirmation,
        title = "Delete Chat",
        message = "Are you sure you want to delete this chat? This action cannot be undone.",
        primaryButton = DialogButton(
            text = "Delete",
            onClick = onConfirm,
            variant = TchatButtonVariant.Destructive
        ),
        dismissButton = DialogButton(
            text = "Cancel",
            onClick = onDismiss,
            variant = TchatButtonVariant.Secondary
        )
    )
}

@Composable
private fun ClearHistoryConfirmationDialog(
    isVisible: Boolean,
    onConfirm: () -> Unit,
    onDismiss: () -> Unit
) {
    TchatDialog(
        isVisible = isVisible,
        onDismissRequest = onDismiss,
        variant = DialogVariant.Confirmation,
        title = "Clear Chat History",
        message = "Are you sure you want to clear all messages in this chat? This action cannot be undone.",
        primaryButton = DialogButton(
            text = "Clear",
            onClick = onConfirm,
            variant = TchatButtonVariant.Destructive
        ),
        dismissButton = DialogButton(
            text = "Cancel",
            onClick = onDismiss,
            variant = TchatButtonVariant.Secondary
        )
    )
}

// Helper function to generate chat options based on state
private fun generateChatOptions(
    chatSession: ChatSession?,
    onPinChat: () -> Unit,
    onMuteChat: () -> Unit,
    onArchiveChat: () -> Unit,
    onBlockChat: () -> Unit,
    onDeleteChat: () -> Unit,
    onClearHistory: () -> Unit,
    onExportChat: () -> Unit,
    onChatInfo: () -> Unit,
    onReportChat: () -> Unit
): List<ChatOption> {
    val isPinned = chatSession?.isPinned ?: false
    val isMuted = chatSession?.isMuted ?: false
    val isArchived = chatSession?.isArchived ?: false
    val isBlocked = chatSession?.isBlocked ?: false

    return listOf(
        // Chat info
        ChatOption(
            id = "info",
            title = "Chat Info",
            subtitle = "View participants and settings",
            icon = Icons.Default.Info,
            action = onChatInfo
        ),

        // Pin/Unpin
        ChatOption(
            id = "pin",
            title = if (isPinned) "Unpin Chat" else "Pin Chat",
            subtitle = if (isPinned) "Remove from pinned chats" else "Keep chat at the top",
            icon = if (isPinned) Icons.Default.PushPin else Icons.Default.PushPin,
            action = onPinChat
        ),

        // Mute/Unmute
        ChatOption(
            id = "mute",
            title = if (isMuted) "Unmute Chat" else "Mute Chat",
            subtitle = if (isMuted) "Enable notifications" else "Disable notifications",
            icon = if (isMuted) Icons.Default.VolumeUp else Icons.Default.VolumeOff,
            action = onMuteChat
        ),

        // Archive/Unarchive
        ChatOption(
            id = "archive",
            title = if (isArchived) "Unarchive Chat" else "Archive Chat",
            subtitle = if (isArchived) "Move back to main chats" else "Move to archived chats",
            icon = if (isArchived) Icons.Default.Unarchive else Icons.Default.Archive,
            action = onArchiveChat
        ),

        // Export chat
        ChatOption(
            id = "export",
            title = "Export Chat",
            subtitle = "Save chat history to file",
            icon = Icons.Default.Download,
            action = onExportChat
        ),

        // Clear history
        ChatOption(
            id = "clear",
            title = "Clear History",
            subtitle = "Delete all messages",
            icon = Icons.Default.DeleteSweep,
            destructive = true,
            action = onClearHistory
        ),

        // Block/Unblock
        ChatOption(
            id = "block",
            title = if (isBlocked) "Unblock Chat" else "Block Chat",
            subtitle = if (isBlocked) "Allow messages from this chat" else "Block messages from this chat",
            icon = if (isBlocked) Icons.Default.PersonAdd else Icons.Default.Block,
            destructive = !isBlocked,
            action = onBlockChat
        ),

        // Report
        ChatOption(
            id = "report",
            title = "Report Chat",
            subtitle = "Report inappropriate content",
            icon = Icons.Default.Report,
            destructive = true,
            action = onReportChat
        ),

        // Delete chat
        ChatOption(
            id = "delete",
            title = "Delete Chat",
            subtitle = "Permanently delete this chat",
            icon = Icons.Default.Delete,
            destructive = true,
            action = onDeleteChat
        )
    )
}