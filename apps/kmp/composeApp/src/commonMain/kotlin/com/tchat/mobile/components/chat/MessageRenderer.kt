package com.tchat.mobile.components.chat

import androidx.compose.foundation.ExperimentalFoundationApi
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.combinedClickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
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
import androidx.compose.ui.platform.testTag
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import com.tchat.mobile.models.*

/**
 * Central message renderer with comprehensive chat actions and interactions
 */
@OptIn(ExperimentalFoundationApi::class)
@Composable
fun MessageRenderer(
    message: Message,
    isFromMe: Boolean,
    onReactionClick: (String) -> Unit = {},
    onLongPress: () -> Unit = {},
    onReplyClick: () -> Unit = {},
    showReactions: Boolean = true,
    showStatus: Boolean = true,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .combinedClickable(
                onClick = {},
                onLongClick = onLongPress
            )
    ) {
        // Reply indicator if this is a reply
        message.replyToId?.let {
            ReplyIndicator(
                replyToContent = "Replying to previous message", // In real implementation, fetch the original message
                modifier = Modifier.padding(bottom = TchatSpacing.xs)
            )
        }

        // Main message content
        MessageContent(
            message = message,
            isFromMe = isFromMe,
            modifier = Modifier.testTag("message-content-${message.id}")
        )

        // Message reactions
        if (showReactions && message.hasReactions()) {
            MessageReactions(
                reactions = message.reactions,
                onReactionClick = onReactionClick,
                modifier = Modifier
                    .padding(top = TchatSpacing.xs)
                    .testTag("message-reactions-${message.id}")
            )
        }

        // Message status and timestamp
        if (showStatus) {
            MessageStatusRow(
                message = message,
                isFromMe = isFromMe,
                modifier = Modifier
                    .padding(top = 4.dp)
                    .testTag("message-status-${message.id}")
            )
        }
    }
}

@Composable
private fun MessageContent(
    message: Message,
    isFromMe: Boolean,
    modifier: Modifier = Modifier
) {
    when (message.type) {
        MessageType.TEXT -> {
            Column {
                if (message.isEdited) {
                    Text(
                        text = "edited",
                        style = MaterialTheme.typography.labelSmall,
                        color = (if (isFromMe) TchatColors.onPrimary else TchatColors.onSurface).copy(alpha = 0.6f),
                        modifier = Modifier.padding(bottom = 2.dp)
                    )
                }
                Text(
                    text = message.content,
                    style = MaterialTheme.typography.bodyMedium,
                    color = if (isFromMe) TchatColors.onPrimary else TchatColors.onSurface,
                    modifier = modifier.testTag("message-text-content")
                )
            }
        }
        MessageType.IMAGE -> {
            ImageMessage(
                message = message,
                modifier = modifier.testTag("message-image-content")
            )
        }
        MessageType.VIDEO -> {
            VideoMessage(
                message = message,
                modifier = modifier.testTag("message-video-content")
            )
        }
        MessageType.AUDIO -> {
            AudioMessage(
                message = message,
                modifier = modifier.testTag("message-audio-content")
            )
        }
        MessageType.FILE -> {
            FileMessage(
                message = message,
                modifier = modifier.testTag("message-file-content")
            )
        }
        MessageType.LOCATION -> {
            LocationMessage(
                message = message,
                modifier = modifier.testTag("message-location-content")
            )
        }
        MessageType.STICKER -> {
            StickerMessage(
                message = message,
                modifier = modifier.testTag("message-sticker-content")
            )
        }
        MessageType.SYSTEM -> {
            SystemMessage(
                message = message,
                modifier = modifier.testTag("message-system-content")
            )
        }
    }
}

@Composable
private fun ReplyIndicator(
    replyToContent: String,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier
            .fillMaxWidth()
            .padding(start = 8.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Icon(
            imageVector = Icons.Default.Reply,
            contentDescription = "Reply",
            modifier = Modifier.size(16.dp),
            tint = TchatColors.onSurfaceVariant
        )
        Spacer(modifier = Modifier.width(4.dp))
        Text(
            text = replyToContent,
            style = MaterialTheme.typography.labelMedium,
            color = TchatColors.onSurfaceVariant,
            maxLines = 1,
            overflow = TextOverflow.Ellipsis
        )
    }
}

@Composable
private fun MessageReactions(
    reactions: List<MessageReaction>,
    onReactionClick: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    if (reactions.isEmpty()) return

    // Group reactions by emoji
    val groupedReactions = reactions.groupBy { it.emoji }

    LazyRow(
        modifier = modifier,
        horizontalArrangement = Arrangement.spacedBy(4.dp)
    ) {
        items(groupedReactions.toList()) { (emoji, reactionList) ->
            ReactionChip(
                emoji = emoji,
                count = reactionList.size,
                onClick = { onReactionClick(emoji) }
            )
        }
    }
}

@Composable
private fun ReactionChip(
    emoji: String,
    count: Int,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Surface(
        modifier = modifier
            .clip(RoundedCornerShape(12.dp)),
        onClick = onClick,
        color = TchatColors.surfaceVariant,
        contentColor = TchatColors.onSurfaceVariant
    ) {
        Row(
            modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                text = emoji,
                style = MaterialTheme.typography.labelMedium
            )
            if (count > 1) {
                Spacer(modifier = Modifier.width(4.dp))
                Text(
                    text = count.toString(),
                    style = MaterialTheme.typography.labelSmall,
                    fontWeight = FontWeight.Medium
                )
            }
        }
    }
}

@Composable
private fun MessageStatusRow(
    message: Message,
    isFromMe: Boolean,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
    ) {
        Row(
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                text = formatTimestamp(message.createdAt),
                style = MaterialTheme.typography.labelSmall,
                color = if (isFromMe)
                    TchatColors.onPrimary.copy(alpha = 0.7f)
                else
                    TchatColors.onSurface.copy(alpha = 0.7f)
            )

            if (isFromMe) {
                Spacer(modifier = Modifier.width(4.dp))
                DeliveryStatusIcon(
                    status = message.deliveryStatus,
                    modifier = Modifier.size(16.dp)
                )
            }

            if (message.isPinned) {
                Spacer(modifier = Modifier.width(4.dp))
                Icon(
                    imageVector = Icons.Default.PushPin,
                    contentDescription = "Pinned",
                    modifier = Modifier.size(14.dp),
                    tint = TchatColors.primary
                )
            }
        }

        // Read status information (like web chat interfaces)
        if (message.deliveryStatus == MessageDeliveryStatus.READ && message.readBy.isNotEmpty()) {
            Row(
                modifier = Modifier.padding(top = 2.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Icon(
                    imageVector = Icons.Default.Visibility,
                    contentDescription = "Read by",
                    modifier = Modifier.size(12.dp),
                    tint = TchatColors.primary.copy(alpha = 0.6f)
                )
                Spacer(modifier = Modifier.width(2.dp))
                Text(
                    text = if (message.readBy.size == 1) {
                        "Read"
                    } else {
                        "Read by ${message.readBy.size} people"
                    },
                    style = MaterialTheme.typography.labelSmall,
                    color = if (isFromMe)
                        TchatColors.onPrimary.copy(alpha = 0.6f)
                    else
                        TchatColors.onSurface.copy(alpha = 0.6f)
                )
            }
        }

        // Edited status
        if (message.isEdited && message.editedAt != null) {
            Text(
                text = "Edited ${formatTimestamp(message.editedAt!!)}",
                style = MaterialTheme.typography.labelSmall,
                color = if (isFromMe)
                    TchatColors.onPrimary.copy(alpha = 0.5f)
                else
                    TchatColors.onSurface.copy(alpha = 0.5f),
                modifier = Modifier.padding(top = 1.dp)
            )
        }
    }
}

@Composable
private fun DeliveryStatusIcon(
    status: MessageDeliveryStatus,
    modifier: Modifier = Modifier
) {
    val (icon, tint) = when (status) {
        MessageDeliveryStatus.SENDING -> Icons.Default.Schedule to TchatColors.onPrimary.copy(alpha = 0.5f)
        MessageDeliveryStatus.SENT -> Icons.Default.Done to TchatColors.onPrimary.copy(alpha = 0.7f)
        MessageDeliveryStatus.DELIVERED -> Icons.Default.DoneAll to TchatColors.onPrimary.copy(alpha = 0.7f)
        MessageDeliveryStatus.READ -> Icons.Default.DoneAll to TchatColors.primary
        MessageDeliveryStatus.FAILED -> Icons.Default.Error to TchatColors.error
    }

    Icon(
        imageVector = icon,
        contentDescription = status.name,
        modifier = modifier,
        tint = tint
    )
}

// Message Action Menu
@Composable
fun MessageActionMenu(
    message: Message,
    isFromMe: Boolean,
    onDismiss: () -> Unit,
    onReply: () -> Unit,
    onEdit: () -> Unit,
    onDelete: () -> Unit,
    onPin: () -> Unit,
    onCopy: () -> Unit,
    onForward: () -> Unit,
    onAddReaction: () -> Unit,
    modifier: Modifier = Modifier
) {
    DropdownMenu(
        expanded = true,
        onDismissRequest = onDismiss,
        modifier = modifier
    ) {
        // Reply action
        if (message.canBeRepliedTo()) {
            DropdownMenuItem(
                text = { Text("Reply") },
                onClick = {
                    onReply()
                    onDismiss()
                },
                leadingIcon = {
                    Icon(Icons.Default.Reply, contentDescription = "Reply")
                }
            )
        }

        // Edit action (only for own text messages)
        if (message.canBeEdited("current_user")) {
            DropdownMenuItem(
                text = { Text("Edit") },
                onClick = {
                    onEdit()
                    onDismiss()
                },
                leadingIcon = {
                    Icon(Icons.Default.Edit, contentDescription = "Edit")
                }
            )
        }

        // Copy action
        if (message.type == MessageType.TEXT) {
            DropdownMenuItem(
                text = { Text("Copy") },
                onClick = {
                    onCopy()
                    onDismiss()
                },
                leadingIcon = {
                    Icon(Icons.Default.ContentCopy, contentDescription = "Copy")
                }
            )
        }

        // Forward action
        DropdownMenuItem(
            text = { Text("Forward") },
            onClick = {
                onForward()
                onDismiss()
            },
            leadingIcon = {
                Icon(Icons.Default.Forward, contentDescription = "Forward")
            }
        )

        // Pin/Unpin action
        DropdownMenuItem(
            text = { Text(if (message.isPinned) "Unpin" else "Pin") },
            onClick = {
                onPin()
                onDismiss()
            },
            leadingIcon = {
                Icon(
                    if (message.isPinned) Icons.Default.PushPin else Icons.Default.PushPin,
                    contentDescription = if (message.isPinned) "Unpin" else "Pin"
                )
            }
        )

        // Add reaction
        DropdownMenuItem(
            text = { Text("Add Reaction") },
            onClick = {
                onAddReaction()
                onDismiss()
            },
            leadingIcon = {
                Icon(Icons.Default.EmojiEmotions, contentDescription = "Add Reaction")
            }
        )

        // Delete action (only for own messages)
        if (message.canBeDeleted("current_user")) {
            Divider()
            DropdownMenuItem(
                text = { Text("Delete") },
                onClick = {
                    onDelete()
                    onDismiss()
                },
                leadingIcon = {
                    Icon(
                        Icons.Default.Delete,
                        contentDescription = "Delete",
                        tint = TchatColors.error
                    )
                },
                colors = MenuDefaults.itemColors(
                    textColor = TchatColors.error
                )
            )
        }
    }
}

// Quick Reaction Picker
@Composable
fun QuickReactionPicker(
    onReactionSelect: (String) -> Unit,
    onDismiss: () -> Unit,
    modifier: Modifier = Modifier
) {
    val quickReactions = listOf("❤️", "👍", "👎", "😂", "😮", "😢", "😡", "🎉")

    Surface(
        modifier = modifier,
        shape = RoundedCornerShape(20.dp),
        color = TchatColors.surface,
        shadowElevation = 8.dp
    ) {
        LazyRow(
            modifier = Modifier.padding(8.dp),
            horizontalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            items(quickReactions) { emoji ->
                Surface(
                    modifier = Modifier
                        .size(40.dp)
                        .clip(CircleShape),
                    onClick = {
                        onReactionSelect(emoji)
                        onDismiss()
                    },
                    color = TchatColors.surfaceVariant
                ) {
                    Box(
                        contentAlignment = Alignment.Center
                    ) {
                        Text(
                            text = emoji,
                            style = MaterialTheme.typography.titleMedium
                        )
                    }
                }
            }
        }
    }
}

// Enhanced Message Context Menu with web-like features
@Composable
fun EnhancedMessageContextMenu(
    message: Message,
    onReactionSelect: (String) -> Unit,
    onCopyMessage: () -> Unit,
    onForwardMessage: () -> Unit,
    onPinMessage: () -> Unit,
    onDeleteMessage: () -> Unit,
    onReplyToMessage: () -> Unit,
    onEditMessage: () -> Unit,
    onDismiss: () -> Unit,
    modifier: Modifier = Modifier
) {
    val quickReactions = listOf("❤️", "👍", "👎", "😂", "😮", "😢", "😡", "🎉")

    Surface(
        modifier = modifier,
        shape = RoundedCornerShape(16.dp),
        color = TchatColors.surface,
        shadowElevation = 12.dp
    ) {
        Column(
            modifier = Modifier
                .padding(TchatSpacing.md)
                .widthIn(min = 280.dp, max = 320.dp)
        ) {
            // Quick reactions section
            Text(
                text = "React",
                style = MaterialTheme.typography.labelMedium,
                color = TchatColors.onSurfaceVariant,
                modifier = Modifier.padding(bottom = TchatSpacing.xs)
            )

            LazyRow(
                horizontalArrangement = Arrangement.spacedBy(TchatSpacing.xs),
                modifier = Modifier.padding(bottom = TchatSpacing.md)
            ) {
                items(quickReactions) { emoji ->
                    Surface(
                        modifier = Modifier
                            .size(36.dp)
                            .clip(CircleShape),
                        onClick = {
                            onReactionSelect(emoji)
                            onDismiss()
                        },
                        color = TchatColors.surfaceVariant
                    ) {
                        Box(contentAlignment = Alignment.Center) {
                            Text(
                                text = emoji,
                                style = MaterialTheme.typography.titleSmall
                            )
                        }
                    }
                }
            }

            HorizontalDivider(
                modifier = Modifier.padding(vertical = TchatSpacing.xs),
                color = TchatColors.outline.copy(alpha = 0.5f)
            )

            // Action items section
            ContextMenuItem(
                icon = Icons.Default.Reply,
                label = "Reply",
                onClick = {
                    onReplyToMessage()
                    onDismiss()
                }
            )

            ContextMenuItem(
                icon = Icons.Default.ContentCopy,
                label = "Copy",
                onClick = {
                    onCopyMessage()
                    onDismiss()
                }
            )

            ContextMenuItem(
                icon = Icons.Default.Forward,
                label = "Forward",
                onClick = {
                    onForwardMessage()
                    onDismiss()
                }
            )

            ContextMenuItem(
                icon = Icons.Default.PushPin,
                label = if (message.isPinned) "Unpin" else "Pin",
                onClick = {
                    onPinMessage()
                    onDismiss()
                }
            )

            if (message.senderId == "current_user") {
                ContextMenuItem(
                    icon = Icons.Default.Edit,
                    label = "Edit",
                    onClick = {
                        onEditMessage()
                        onDismiss()
                    }
                )
            }

            HorizontalDivider(
                modifier = Modifier.padding(vertical = TchatSpacing.xs),
                color = TchatColors.outline.copy(alpha = 0.5f)
            )

            // Destructive actions
            ContextMenuItem(
                icon = Icons.Default.Delete,
                label = "Delete",
                onClick = {
                    onDeleteMessage()
                    onDismiss()
                },
                isDestructive = true
            )
        }
    }
}

@Composable
private fun ContextMenuItem(
    icon: ImageVector,
    label: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    isDestructive: Boolean = false
) {
    Surface(
        modifier = modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(8.dp)),
        onClick = onClick,
        color = TchatColors.surface
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Icon(
                imageVector = icon,
                contentDescription = label,
                tint = if (isDestructive) TchatColors.error else TchatColors.onSurface,
                modifier = Modifier.size(20.dp)
            )

            Spacer(modifier = Modifier.width(TchatSpacing.md))

            Text(
                text = label,
                style = MaterialTheme.typography.bodyMedium,
                color = if (isDestructive) TchatColors.error else TchatColors.onSurface
            )
        }
    }
}

// Chat Member List with Online Status - Web-like feature
@Composable
fun ChatMemberList(
    chatSession: ChatSession,
    currentUserId: String,
    onMemberClick: (ChatParticipant) -> Unit = {},
    onMemberAction: (ChatParticipant, MemberAction) -> Unit = { _, _ -> },
    modifier: Modifier = Modifier
) {
    val onlineMembers = chatSession.participants.filter { it.status == ParticipantStatus.ONLINE }
    val awayMembers = chatSession.participants.filter { it.status == ParticipantStatus.AWAY || it.status == ParticipantStatus.BUSY }
    val offlineMembers = chatSession.participants.filter { it.status == ParticipantStatus.OFFLINE || it.status == ParticipantStatus.INVISIBLE }

    LazyColumn(
        modifier = modifier
            .fillMaxWidth()
            .testTag("chat-member-list"),
        contentPadding = PaddingValues(TchatSpacing.md),
        verticalArrangement = Arrangement.spacedBy(TchatSpacing.xs)
    ) {
        // Header with member count and online status
        item {
            MemberListHeader(
                totalMembers = chatSession.participants.size,
                onlineCount = chatSession.getOnlineCount(),
                chatType = chatSession.type
            )
        }

        // Online members section
        if (onlineMembers.isNotEmpty()) {
            item {
                MemberSectionHeader(
                    title = "Online",
                    count = onlineMembers.size,
                    statusColor = TchatColors.primary
                )
            }
            items(onlineMembers) { member ->
                MemberListItem(
                    member = member,
                    isCurrentUser = member.id == currentUserId,
                    onClick = { onMemberClick(member) },
                    onAction = { action -> onMemberAction(member, action) }
                )
            }
        }

        // Away/Busy members section
        if (awayMembers.isNotEmpty()) {
            item {
                MemberSectionHeader(
                    title = "Away",
                    count = awayMembers.size,
                    statusColor = TchatColors.warning
                )
            }
            items(awayMembers) { member ->
                MemberListItem(
                    member = member,
                    isCurrentUser = member.id == currentUserId,
                    onClick = { onMemberClick(member) },
                    onAction = { action -> onMemberAction(member, action) }
                )
            }
        }

        // Offline members section
        if (offlineMembers.isNotEmpty()) {
            item {
                MemberSectionHeader(
                    title = "Offline",
                    count = offlineMembers.size,
                    statusColor = TchatColors.onSurfaceVariant.copy(alpha = 0.6f)
                )
            }
            items(offlineMembers) { member ->
                MemberListItem(
                    member = member,
                    isCurrentUser = member.id == currentUserId,
                    onClick = { onMemberClick(member) },
                    onAction = { action -> onMemberAction(member, action) }
                )
            }
        }
    }
}

@Composable
private fun MemberListHeader(
    totalMembers: Int,
    onlineCount: Int,
    chatType: ChatType,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier
            .fillMaxWidth()
            .testTag("member-list-header"),
        colors = CardDefaults.cardColors(
            containerColor = TchatColors.primaryLight.copy(alpha = 0.3f)
        )
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            Text(
                text = when (chatType) {
                    ChatType.GROUP -> "Group Members"
                    ChatType.CHANNEL -> "Channel Participants"
                    ChatType.DIRECT -> "Conversation"
                    ChatType.SUPPORT -> "Support Chat"
                    ChatType.SYSTEM -> "System Chat"
                },
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.SemiBold,
                color = TchatColors.onSurface
            )

            Spacer(modifier = Modifier.height(TchatSpacing.xs))

            Row(
                verticalAlignment = Alignment.CenterVertically
            ) {
                // Online indicator
                Box(
                    modifier = Modifier
                        .size(8.dp)
                        .background(TchatColors.primary, CircleShape)
                )

                Spacer(modifier = Modifier.width(TchatSpacing.xs))

                Text(
                    text = "$onlineCount of $totalMembers online",
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurfaceVariant
                )
            }
        }
    }
}

@Composable
private fun MemberSectionHeader(
    title: String,
    count: Int,
    statusColor: androidx.compose.ui.graphics.Color,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier
            .fillMaxWidth()
            .padding(vertical = TchatSpacing.sm),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Box(
            modifier = Modifier
                .size(6.dp)
                .background(statusColor, CircleShape)
        )

        Spacer(modifier = Modifier.width(TchatSpacing.sm))

        Text(
            text = "$title ($count)",
            style = MaterialTheme.typography.labelLarge,
            fontWeight = FontWeight.Medium,
            color = TchatColors.onSurfaceVariant
        )
    }
}

@Composable
private fun MemberListItem(
    member: ChatParticipant,
    isCurrentUser: Boolean,
    onClick: () -> Unit,
    onAction: (MemberAction) -> Unit,
    modifier: Modifier = Modifier
) {
    var showActionMenu by remember { mutableStateOf(false) }

    Surface(
        modifier = modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(8.dp))
            .combinedClickable(
                onClick = onClick,
                onLongClick = { showActionMenu = true }
            )
            .testTag("member-item-${member.id}"),
        color = if (isCurrentUser)
            TchatColors.primaryLight.copy(alpha = 0.2f)
        else
            TchatColors.surface
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Avatar with status indicator
            Box {
                // Avatar background
                Box(
                    modifier = Modifier
                        .size(40.dp)
                        .clip(CircleShape)
                        .background(TchatColors.primary),
                    contentAlignment = Alignment.Center
                ) {
                    Text(
                        text = member.name.firstOrNull()?.toString() ?: "?",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold,
                        color = TchatColors.onPrimary
                    )
                }

                // Status indicator
                Box(
                    modifier = Modifier
                        .align(Alignment.BottomEnd)
                        .size(12.dp)
                        .background(
                            when (member.status) {
                                ParticipantStatus.ONLINE -> TchatColors.primary
                                ParticipantStatus.AWAY, ParticipantStatus.BUSY -> TchatColors.warning
                                ParticipantStatus.TYPING -> TchatColors.success
                                else -> TchatColors.onSurfaceVariant.copy(alpha = 0.6f)
                            },
                            CircleShape
                        )
                        .border(2.dp, TchatColors.surface, CircleShape)
                )
            }

            Spacer(modifier = Modifier.width(TchatSpacing.md))

            // Member info
            Column(
                modifier = Modifier.weight(1f)
            ) {
                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        text = if (isCurrentUser) "${member.name} (You)" else member.name,
                        style = MaterialTheme.typography.bodyLarge,
                        fontWeight = if (isCurrentUser) FontWeight.SemiBold else FontWeight.Normal,
                        color = TchatColors.onSurface
                    )

                    // Role badge
                    if (member.role != ChatRole.MEMBER) {
                        Spacer(modifier = Modifier.width(TchatSpacing.xs))

                        Surface(
                            shape = RoundedCornerShape(4.dp),
                            color = when (member.role) {
                                ChatRole.OWNER -> TchatColors.errorLight
                                ChatRole.ADMIN -> TchatColors.primaryLight
                                ChatRole.MODERATOR -> TchatColors.surfaceVariant
                                else -> TchatColors.surfaceVariant
                            }
                        ) {
                            Text(
                                text = member.role.name.lowercase().replaceFirstChar { it.uppercase() },
                                style = MaterialTheme.typography.labelSmall,
                                color = when (member.role) {
                                    ChatRole.OWNER -> TchatColors.onPrimary
                                    ChatRole.ADMIN -> TchatColors.onPrimary
                                    ChatRole.MODERATOR -> TchatColors.onSurface
                                    else -> TchatColors.onSurfaceVariant
                                },
                                modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp)
                            )
                        }
                    }

                    // Bot indicator
                    if (member.isBot) {
                        Spacer(modifier = Modifier.width(TchatSpacing.xs))

                        Icon(
                            Icons.Default.SmartToy,
                            contentDescription = "Bot",
                            modifier = Modifier.size(16.dp),
                            tint = TchatColors.onSurfaceVariant
                        )
                    }
                }

                // Status text
                Text(
                    text = member.getDisplayStatus(),
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant
                )

                // Custom title if available
                member.customTitle?.let { title ->
                    Text(
                        text = title,
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.primary,
                        fontStyle = FontStyle.Italic
                    )
                }
            }

            // Action menu trigger
            IconButton(
                onClick = { showActionMenu = true },
                modifier = Modifier.size(32.dp)
            ) {
                Icon(
                    Icons.Default.MoreVert,
                    contentDescription = "Member Actions",
                    modifier = Modifier.size(20.dp),
                    tint = TchatColors.onSurfaceVariant
                )
            }
        }
    }

    // Member action dropdown menu
    if (showActionMenu) {
        MemberActionMenu(
            member = member,
            isCurrentUser = isCurrentUser,
            onAction = { action ->
                onAction(action)
                showActionMenu = false
            },
            onDismiss = { showActionMenu = false }
        )
    }
}

@Composable
private fun MemberActionMenu(
    member: ChatParticipant,
    isCurrentUser: Boolean,
    onAction: (MemberAction) -> Unit,
    onDismiss: () -> Unit,
    modifier: Modifier = Modifier
) {
    DropdownMenu(
        expanded = true,
        onDismissRequest = onDismiss,
        modifier = modifier
    ) {
        if (!isCurrentUser) {
            DropdownMenuItem(
                text = { Text("View Profile") },
                onClick = { onAction(MemberAction.VIEW_PROFILE) },
                leadingIcon = {
                    Icon(Icons.Default.Person, contentDescription = null)
                }
            )

            DropdownMenuItem(
                text = { Text("Send Message") },
                onClick = { onAction(MemberAction.SEND_MESSAGE) },
                leadingIcon = {
                    Icon(Icons.Default.Message, contentDescription = null)
                }
            )

            if (member.canModerate()) {
                HorizontalDivider()

                DropdownMenuItem(
                    text = { Text("Manage Role") },
                    onClick = { onAction(MemberAction.MANAGE_ROLE) },
                    leadingIcon = {
                        Icon(Icons.Default.AdminPanelSettings, contentDescription = null)
                    }
                )

                DropdownMenuItem(
                    text = { Text("Remove Member", color = TchatColors.error) },
                    onClick = { onAction(MemberAction.REMOVE_MEMBER) },
                    leadingIcon = {
                        Icon(Icons.Default.PersonRemove, contentDescription = null, tint = TchatColors.error)
                    }
                )
            }
        }
    }
}

enum class MemberAction {
    VIEW_PROFILE,
    SEND_MESSAGE,
    MANAGE_ROLE,
    REMOVE_MEMBER
}

// Comprehensive Message Pinning System - Web-like pinned messages management

/**
 * Pinned Messages Bar - Displays all pinned messages in a collapsible section
 * Similar to Slack/Discord pinned messages functionality
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun PinnedMessagesBar(
    pinnedMessages: List<Message>,
    onPinnedMessageClick: (Message) -> Unit,
    onUnpinMessage: (Message) -> Unit,
    onClosePinnedBar: () -> Unit,
    modifier: Modifier = Modifier
) {
    var isExpanded by remember { mutableStateOf(false) }

    if (pinnedMessages.isEmpty()) return

    Card(
        modifier = modifier
            .fillMaxWidth()
            .testTag("pinned-messages-bar"),
        colors = CardDefaults.cardColors(
            containerColor = TchatColors.primaryLight.copy(alpha = 0.3f)
        ),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column {
            // Header section - always visible
            Surface(
                modifier = Modifier
                    .fillMaxWidth()
                    .combinedClickable(
                        onClick = { isExpanded = !isExpanded }
                    ),
                color = TchatColors.primaryLight.copy(alpha = 0.1f)
            ) {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(TchatSpacing.md),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        Icons.Default.PushPin,
                        contentDescription = "Pinned Messages",
                        tint = TchatColors.primary,
                        modifier = Modifier.size(20.dp)
                    )

                    Spacer(modifier = Modifier.width(TchatSpacing.sm))

                    Text(
                        text = "${pinnedMessages.size} pinned message${if (pinnedMessages.size != 1) "s" else ""}",
                        style = MaterialTheme.typography.bodyMedium,
                        fontWeight = FontWeight.Medium,
                        color = TchatColors.primary,
                        modifier = Modifier.weight(1f)
                    )

                    // Expand/Collapse indicator
                    Icon(
                        if (isExpanded) Icons.Default.ExpandLess else Icons.Default.ExpandMore,
                        contentDescription = if (isExpanded) "Collapse" else "Expand",
                        tint = TchatColors.primary,
                        modifier = Modifier.size(24.dp)
                    )

                    Spacer(modifier = Modifier.width(TchatSpacing.xs))

                    // Close button
                    IconButton(
                        onClick = onClosePinnedBar,
                        modifier = Modifier.size(32.dp)
                    ) {
                        Icon(
                            Icons.Default.Close,
                            contentDescription = "Close pinned messages",
                            tint = TchatColors.primary,
                            modifier = Modifier.size(18.dp)
                        )
                    }
                }
            }

            // Expandable content section
            if (isExpanded) {
                LazyColumn(
                    modifier = Modifier
                        .fillMaxWidth()
                        .heightIn(max = 300.dp),
                    contentPadding = PaddingValues(TchatSpacing.md),
                    verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                ) {
                    items(pinnedMessages, key = { it.id }) { message ->
                        PinnedMessageCard(
                            message = message,
                            onClick = { onPinnedMessageClick(message) },
                            onUnpin = { onUnpinMessage(message) }
                        )
                    }
                }
            } else {
                // Compact view - show most recent pinned message
                pinnedMessages.firstOrNull()?.let { recentMessage ->
                    PinnedMessageCard(
                        message = recentMessage,
                        onClick = { onPinnedMessageClick(recentMessage) },
                        onUnpin = { onUnpinMessage(recentMessage) },
                        isCompact = true,
                        modifier = Modifier.padding(TchatSpacing.md)
                    )
                }
            }
        }
    }
}

/**
 * Individual Pinned Message Card
 */
@Composable
fun PinnedMessageCard(
    message: Message,
    onClick: () -> Unit,
    onUnpin: () -> Unit,
    modifier: Modifier = Modifier,
    isCompact: Boolean = false
) {
    Surface(
        modifier = modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(12.dp))
            .combinedClickable(onClick = onClick)
            .testTag("pinned-message-${message.id}"),
        color = TchatColors.surface,
        shadowElevation = 1.dp
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md),
            verticalAlignment = Alignment.Top
        ) {
            // Pin indicator
            Icon(
                Icons.Default.PushPin,
                contentDescription = "Pinned",
                tint = TchatColors.primary,
                modifier = Modifier
                    .size(16.dp)
                    .padding(top = 2.dp)
            )

            Spacer(modifier = Modifier.width(TchatSpacing.sm))

            // Message content
            Column(
                modifier = Modifier.weight(1f)
            ) {
                // Sender name and timestamp
                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        text = message.senderName,
                        style = MaterialTheme.typography.labelMedium,
                        fontWeight = FontWeight.SemiBold,
                        color = TchatColors.primary
                    )

                    Spacer(modifier = Modifier.width(TchatSpacing.xs))

                    Text(
                        text = formatTimestamp(message.createdAt),
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }

                Spacer(modifier = Modifier.height(TchatSpacing.xs))

                // Message preview
                Text(
                    text = when (message.type) {
                        MessageType.TEXT -> message.content
                        MessageType.IMAGE -> "📷 Image"
                        MessageType.VIDEO -> "🎥 Video"
                        MessageType.AUDIO -> "🎵 Audio"
                        MessageType.FILE -> "📄 ${message.attachments.firstOrNull()?.filename ?: "File"}"
                        MessageType.LOCATION -> "📍 Location"
                        MessageType.STICKER -> "😀 Sticker"
                        else -> message.content
                    },
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurface,
                    maxLines = if (isCompact) 1 else 3,
                    overflow = TextOverflow.Ellipsis
                )

                // Show reactions if any
                if (message.hasReactions() && !isCompact) {
                    Spacer(modifier = Modifier.height(TchatSpacing.xs))

                    LazyRow(
                        horizontalArrangement = Arrangement.spacedBy(4.dp)
                    ) {
                        items(message.reactions.groupBy { it.emoji }.toList().take(3)) { (emoji, reactionList) ->
                            Surface(
                                shape = RoundedCornerShape(8.dp),
                                color = TchatColors.surfaceVariant
                            ) {
                                Text(
                                    text = "$emoji ${reactionList.size}",
                                    style = MaterialTheme.typography.labelSmall,
                                    modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp)
                                )
                            }
                        }

                        val groupedReactions = message.reactions.groupBy { it.emoji }
                        if (groupedReactions.size > 3) {
                            item {
                                Surface(
                                    shape = RoundedCornerShape(8.dp),
                                    color = TchatColors.surfaceVariant
                                ) {
                                    Text(
                                        text = "+${groupedReactions.size - 3}",
                                        style = MaterialTheme.typography.labelSmall,
                                        modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp)
                                    )
                                }
                            }
                        }
                    }
                }
            }

            // Unpin action
            IconButton(
                onClick = onUnpin,
                modifier = Modifier.size(32.dp)
            ) {
                Icon(
                    Icons.Default.Close,
                    contentDescription = "Unpin message",
                    tint = TchatColors.onSurfaceVariant,
                    modifier = Modifier.size(16.dp)
                )
            }
        }
    }
}

/**
 * Enhanced Message Pin Indicator - Shows pin status prominently
 */
@Composable
fun EnhancedPinIndicator(
    isPinned: Boolean,
    onPinToggle: () -> Unit,
    modifier: Modifier = Modifier,
    showLabel: Boolean = false
) {
    if (isPinned) {
        Surface(
            modifier = modifier
                .clip(RoundedCornerShape(12.dp))
                .combinedClickable(onClick = onPinToggle),
            color = TchatColors.primaryLight.copy(alpha = 0.8f)
        ) {
            Row(
                modifier = Modifier.padding(
                    horizontal = if (showLabel) TchatSpacing.sm else TchatSpacing.xs,
                    vertical = TchatSpacing.xs
                ),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Icon(
                    Icons.Default.PushPin,
                    contentDescription = "Pinned message",
                    tint = TchatColors.primary,
                    modifier = Modifier.size(14.dp)
                )

                if (showLabel) {
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(
                        text = "Pinned",
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.primary,
                        fontWeight = FontWeight.Medium
                    )
                }
            }
        }
    }
}

/**
 * Pin Management Dialog - For managing all pinned messages
 */
@Composable
fun PinManagementDialog(
    pinnedMessages: List<Message>,
    onUnpinMessage: (Message) -> Unit,
    onUnpinAll: () -> Unit,
    onDismiss: () -> Unit,
    modifier: Modifier = Modifier
) {
    AlertDialog(
        onDismissRequest = onDismiss,
        modifier = modifier.testTag("pin-management-dialog"),
        title = {
            Row(
                verticalAlignment = Alignment.CenterVertically
            ) {
                Icon(
                    Icons.Default.PushPin,
                    contentDescription = null,
                    tint = TchatColors.primary,
                    modifier = Modifier.size(24.dp)
                )
                Spacer(modifier = Modifier.width(TchatSpacing.sm))
                Text(
                    text = "Pinned Messages",
                    style = MaterialTheme.typography.headlineSmall
                )
            }
        },
        text = {
            if (pinnedMessages.isEmpty()) {
                Text(
                    text = "No messages are currently pinned.",
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurfaceVariant
                )
            } else {
                LazyColumn(
                    modifier = Modifier.heightIn(max = 400.dp),
                    verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                ) {
                    items(pinnedMessages, key = { it.id }) { message ->
                        PinnedMessageListItem(
                            message = message,
                            onUnpin = { onUnpinMessage(message) }
                        )
                    }
                }
            }
        },
        confirmButton = {
            if (pinnedMessages.isNotEmpty()) {
                TextButton(
                    onClick = {
                        onUnpinAll()
                        onDismiss()
                    }
                ) {
                    Text(
                        "Unpin All",
                        color = TchatColors.error
                    )
                }
            }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) {
                Text("Close")
            }
        }
    )
}

@Composable
private fun PinnedMessageListItem(
    message: Message,
    onUnpin: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(
            containerColor = TchatColors.surfaceVariant.copy(alpha = 0.5f)
        )
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Column(
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = message.senderName,
                    style = MaterialTheme.typography.labelMedium,
                    fontWeight = FontWeight.SemiBold,
                    color = TchatColors.primary
                )

                Spacer(modifier = Modifier.height(2.dp))

                Text(
                    text = message.content,
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurface,
                    maxLines = 2,
                    overflow = TextOverflow.Ellipsis
                )

                Spacer(modifier = Modifier.height(2.dp))

                Text(
                    text = formatTimestamp(message.createdAt),
                    style = MaterialTheme.typography.labelSmall,
                    color = TchatColors.onSurfaceVariant
                )
            }

            IconButton(
                onClick = onUnpin,
                modifier = Modifier.size(32.dp)
            ) {
                Icon(
                    Icons.Default.Close,
                    contentDescription = "Unpin",
                    tint = TchatColors.error,
                    modifier = Modifier.size(18.dp)
                )
            }
        }
    }
}

// Advanced Message Forwarding System - Web-like multi-chat forwarding modal

/**
 * Advanced Forwarding Modal - Comprehensive message forwarding with chat selection
 * Similar to WhatsApp/Telegram/Slack forwarding functionality
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun AdvancedForwardingModal(
    message: Message,
    availableChats: List<ChatSession>,
    onForwardToChats: (List<String>, String?) -> Unit, // chatIds, optional comment
    onDismiss: () -> Unit,
    modifier: Modifier = Modifier
) {
    var selectedChatIds by remember { mutableStateOf(setOf<String>()) }
    var searchQuery by remember { mutableStateOf("") }
    var forwardingComment by remember { mutableStateOf("") }
    var currentFilter by remember { mutableStateOf(ChatFilter.ALL) }
    var showCommentField by remember { mutableStateOf(false) }

    // Filter chats based on search and filter type
    val filteredChats = availableChats.filter { chat ->
        val matchesSearch = if (searchQuery.isBlank()) true else {
            chat.getDisplayName("current_user").lowercase().contains(searchQuery.lowercase()) ||
            chat.participants.any { it.name.lowercase().contains(searchQuery.lowercase()) }
        }

        val matchesFilter = when (currentFilter) {
            ChatFilter.ALL -> true
            ChatFilter.DIRECT -> chat.type == ChatType.DIRECT
            ChatFilter.GROUPS -> chat.type == ChatType.GROUP
            ChatFilter.CHANNELS -> chat.type == ChatType.CHANNEL
        }

        matchesSearch && matchesFilter
    }

    AlertDialog(
        onDismissRequest = onDismiss,
        modifier = modifier
            .fillMaxWidth()
            .heightIn(max = 600.dp)
            .testTag("advanced-forwarding-modal"),
        title = {
            Column {
                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        Icons.Default.Forward,
                        contentDescription = null,
                        tint = TchatColors.primary,
                        modifier = Modifier.size(24.dp)
                    )
                    Spacer(modifier = Modifier.width(TchatSpacing.sm))
                    Text(
                        text = "Forward Message",
                        style = MaterialTheme.typography.headlineSmall
                    )
                }

                // Selection count
                if (selectedChatIds.isNotEmpty()) {
                    Spacer(modifier = Modifier.height(TchatSpacing.xs))
                    Text(
                        text = "${selectedChatIds.size} chat${if (selectedChatIds.size != 1) "s" else ""} selected",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.primary
                    )
                }
            }
        },
        text = {
            Column {
                // Message preview
                ForwardingMessagePreview(
                    message = message,
                    modifier = Modifier.padding(bottom = TchatSpacing.md)
                )

                // Search field
                OutlinedTextField(
                    value = searchQuery,
                    onValueChange = { searchQuery = it },
                    label = { Text("Search chats") },
                    leadingIcon = {
                        Icon(Icons.Default.Search, contentDescription = "Search")
                    },
                    trailingIcon = {
                        if (searchQuery.isNotEmpty()) {
                            IconButton(onClick = { searchQuery = "" }) {
                                Icon(Icons.Default.Clear, contentDescription = "Clear")
                            }
                        }
                    },
                    modifier = Modifier
                        .fillMaxWidth()
                        .testTag("forward-search-field"),
                    singleLine = true
                )

                Spacer(modifier = Modifier.height(TchatSpacing.sm))

                // Filter chips
                ForwardingFilterChips(
                    currentFilter = currentFilter,
                    onFilterChange = { currentFilter = it },
                    chatCounts = mapOf(
                        ChatFilter.ALL to availableChats.size,
                        ChatFilter.DIRECT to availableChats.count { it.type == ChatType.DIRECT },
                        ChatFilter.GROUPS to availableChats.count { it.type == ChatType.GROUP },
                        ChatFilter.CHANNELS to availableChats.count { it.type == ChatType.CHANNEL }
                    )
                )

                Spacer(modifier = Modifier.height(TchatSpacing.sm))

                // Chat selection list
                LazyColumn(
                    modifier = Modifier
                        .fillMaxWidth()
                        .heightIn(max = 280.dp),
                    verticalArrangement = Arrangement.spacedBy(TchatSpacing.xs)
                ) {
                    items(filteredChats, key = { it.id }) { chat ->
                        ForwardingChatItem(
                            chat = chat,
                            isSelected = selectedChatIds.contains(chat.id),
                            onSelectionChange = { isSelected ->
                                selectedChatIds = if (isSelected) {
                                    selectedChatIds + chat.id
                                } else {
                                    selectedChatIds - chat.id
                                }
                            }
                        )
                    }

                    if (filteredChats.isEmpty()) {
                        item {
                            Box(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .padding(TchatSpacing.lg),
                                contentAlignment = Alignment.Center
                            ) {
                                Column(
                                    horizontalAlignment = Alignment.CenterHorizontally
                                ) {
                                    Icon(
                                        Icons.Default.SearchOff,
                                        contentDescription = null,
                                        modifier = Modifier.size(48.dp),
                                        tint = TchatColors.onSurfaceVariant.copy(alpha = 0.6f)
                                    )
                                    Spacer(modifier = Modifier.height(TchatSpacing.sm))
                                    Text(
                                        text = if (searchQuery.isNotEmpty()) "No chats found" else "No chats available",
                                        style = MaterialTheme.typography.bodyMedium,
                                        color = TchatColors.onSurfaceVariant
                                    )
                                }
                            }
                        }
                    }
                }

                // Comment field toggle
                if (selectedChatIds.isNotEmpty()) {
                    Spacer(modifier = Modifier.height(TchatSpacing.sm))

                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Text(
                            text = "Add comment",
                            style = MaterialTheme.typography.bodyMedium,
                            color = TchatColors.onSurface
                        )

                        Switch(
                            checked = showCommentField,
                            onCheckedChange = {
                                showCommentField = it
                                if (!it) forwardingComment = ""
                            }
                        )
                    }

                    // Comment input field
                    if (showCommentField) {
                        Spacer(modifier = Modifier.height(TchatSpacing.sm))

                        OutlinedTextField(
                            value = forwardingComment,
                            onValueChange = { forwardingComment = it },
                            label = { Text("Add a comment...") },
                            modifier = Modifier
                                .fillMaxWidth()
                                .testTag("forward-comment-field"),
                            maxLines = 3,
                            placeholder = { Text("Optional comment to send with the forwarded message") }
                        )
                    }
                }
            }
        },
        confirmButton = {
            Button(
                onClick = {
                    onForwardToChats(
                        selectedChatIds.toList(),
                        forwardingComment.takeIf { it.isNotBlank() }
                    )
                    onDismiss()
                },
                enabled = selectedChatIds.isNotEmpty(),
                modifier = Modifier.testTag("forward-confirm-button")
            ) {
                Icon(
                    Icons.Default.Send,
                    contentDescription = null,
                    modifier = Modifier.size(18.dp)
                )
                Spacer(modifier = Modifier.width(TchatSpacing.xs))
                Text("Forward")
            }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) {
                Text("Cancel")
            }
        }
    )
}

/**
 * Message preview for forwarding modal
 */
@Composable
private fun ForwardingMessagePreview(
    message: Message,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(
            containerColor = TchatColors.surfaceVariant.copy(alpha = 0.5f)
        )
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            Row(
                verticalAlignment = Alignment.CenterVertically
            ) {
                Icon(
                    Icons.Default.Forward,
                    contentDescription = "Forwarding",
                    tint = TchatColors.primary,
                    modifier = Modifier.size(16.dp)
                )
                Spacer(modifier = Modifier.width(TchatSpacing.xs))
                Text(
                    text = "Forwarding from ${message.senderName}",
                    style = MaterialTheme.typography.labelMedium,
                    color = TchatColors.primary,
                    fontWeight = FontWeight.Medium
                )
            }

            Spacer(modifier = Modifier.height(TchatSpacing.xs))

            // Message content preview
            Text(
                text = when (message.type) {
                    MessageType.TEXT -> message.content
                    MessageType.IMAGE -> "📷 Image"
                    MessageType.VIDEO -> "🎥 Video"
                    MessageType.AUDIO -> "🎵 Audio"
                    MessageType.FILE -> "📄 ${message.attachments.firstOrNull()?.filename ?: "File"}"
                    MessageType.LOCATION -> "📍 Location"
                    MessageType.STICKER -> "😀 Sticker"
                    else -> message.content
                },
                style = MaterialTheme.typography.bodyMedium,
                color = TchatColors.onSurface,
                maxLines = 2,
                overflow = TextOverflow.Ellipsis
            )

            // Show attachment info if available
            if (message.attachments.isNotEmpty()) {
                Spacer(modifier = Modifier.height(TchatSpacing.xs))
                val attachment = message.attachments.first()
                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        when (message.type) {
                            MessageType.IMAGE -> Icons.Default.Image
                            MessageType.VIDEO -> Icons.Default.VideoFile
                            MessageType.AUDIO -> Icons.Default.AudioFile
                            else -> Icons.Default.AttachFile
                        },
                        contentDescription = null,
                        modifier = Modifier.size(14.dp),
                        tint = TchatColors.onSurfaceVariant
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(
                        text = attachment.filename ?: "Attachment",
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            }
        }
    }
}

/**
 * Filter chips for chat types
 */
@Composable
private fun ForwardingFilterChips(
    currentFilter: ChatFilter,
    onFilterChange: (ChatFilter) -> Unit,
    chatCounts: Map<ChatFilter, Int>,
    modifier: Modifier = Modifier
) {
    LazyRow(
        modifier = modifier,
        horizontalArrangement = Arrangement.spacedBy(TchatSpacing.xs)
    ) {
        items(ChatFilter.entries) { filter ->
            FilterChip(
                selected = currentFilter == filter,
                onClick = { onFilterChange(filter) },
                label = {
                    Text(
                        text = "${filter.displayName} (${chatCounts[filter] ?: 0})",
                        style = MaterialTheme.typography.labelMedium
                    )
                },
                leadingIcon = {
                    Icon(
                        when (filter) {
                            ChatFilter.ALL -> Icons.Default.Chat
                            ChatFilter.DIRECT -> Icons.Default.Person
                            ChatFilter.GROUPS -> Icons.Default.Group
                            ChatFilter.CHANNELS -> Icons.Default.Campaign
                        },
                        contentDescription = null,
                        modifier = Modifier.size(16.dp)
                    )
                }
            )
        }
    }
}

/**
 * Individual chat item in the forwarding list
 */
@Composable
private fun ForwardingChatItem(
    chat: ChatSession,
    isSelected: Boolean,
    onSelectionChange: (Boolean) -> Unit,
    modifier: Modifier = Modifier
) {
    Surface(
        modifier = modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(8.dp))
            .combinedClickable(
                onClick = { onSelectionChange(!isSelected) }
            )
            .testTag("forward-chat-item-${chat.id}"),
        color = if (isSelected)
            TchatColors.primaryLight.copy(alpha = 0.3f)
        else
            TchatColors.surface
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Selection checkbox
            Checkbox(
                checked = isSelected,
                onCheckedChange = onSelectionChange,
                colors = CheckboxDefaults.colors(
                    checkedColor = TchatColors.primary,
                    uncheckedColor = TchatColors.onSurfaceVariant
                )
            )

            Spacer(modifier = Modifier.width(TchatSpacing.sm))

            // Chat avatar
            Box(
                modifier = Modifier
                    .size(40.dp)
                    .clip(CircleShape)
                    .background(
                        when (chat.type) {
                            ChatType.DIRECT -> TchatColors.primary
                            ChatType.GROUP -> TchatColors.warning
                            ChatType.CHANNEL -> TchatColors.success
                            else -> TchatColors.surfaceVariant
                        }
                    ),
                contentAlignment = Alignment.Center
            ) {
                if (chat.type == ChatType.DIRECT) {
                    // For direct chats, show first letter of other participant
                    val otherParticipant = chat.getOtherParticipant("current_user")
                    Text(
                        text = otherParticipant?.name?.firstOrNull()?.toString() ?: "?",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold,
                        color = TchatColors.onPrimary
                    )
                } else {
                    // For groups/channels, show appropriate icon
                    Icon(
                        when (chat.type) {
                            ChatType.GROUP -> Icons.Default.Group
                            ChatType.CHANNEL -> Icons.Default.Campaign
                            else -> Icons.Default.Chat
                        },
                        contentDescription = null,
                        tint = TchatColors.onPrimary,
                        modifier = Modifier.size(20.dp)
                    )
                }
            }

            Spacer(modifier = Modifier.width(TchatSpacing.md))

            // Chat info
            Column(
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = chat.getDisplayName("current_user"),
                    style = MaterialTheme.typography.bodyLarge,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.onSurface,
                    maxLines = 1,
                    overflow = TextOverflow.Ellipsis
                )

                Spacer(modifier = Modifier.height(2.dp))

                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    // Chat type indicator
                    Text(
                        text = when (chat.type) {
                            ChatType.DIRECT -> "Direct message"
                            ChatType.GROUP -> "${chat.participants.size} members"
                            ChatType.CHANNEL -> "Channel • ${chat.participants.size} subscribers"
                            ChatType.SUPPORT -> "Support chat"
                            ChatType.SYSTEM -> "System chat"
                        },
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )

                    // Online status for direct chats
                    if (chat.type == ChatType.DIRECT) {
                        val otherParticipant = chat.getOtherParticipant("current_user")
                        if (otherParticipant?.status == ParticipantStatus.ONLINE) {
                            Spacer(modifier = Modifier.width(TchatSpacing.xs))
                            Box(
                                modifier = Modifier
                                    .size(6.dp)
                                    .background(TchatColors.primary, CircleShape)
                            )
                        }
                    }
                }
            }

            // Additional chat info
            Column(
                horizontalAlignment = Alignment.End
            ) {
                // Last message time
                chat.lastMessage?.timestamp?.let { timestamp ->
                    Text(
                        text = formatTimestamp(timestamp),
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }

                // Unread count or muted indicator
                if (chat.hasUnreadMessages() && !chat.isMuted) {
                    Spacer(modifier = Modifier.height(2.dp))
                    Surface(
                        shape = CircleShape,
                        color = TchatColors.primary
                    ) {
                        Text(
                            text = if (chat.unreadCount > 99) "99+" else chat.unreadCount.toString(),
                            style = MaterialTheme.typography.labelSmall,
                            color = TchatColors.onPrimary,
                            fontWeight = FontWeight.Bold,
                            modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp)
                        )
                    }
                } else if (chat.isMuted) {
                    Spacer(modifier = Modifier.height(2.dp))
                    Icon(
                        Icons.Default.VolumeOff,
                        contentDescription = "Muted",
                        modifier = Modifier.size(14.dp),
                        tint = TchatColors.onSurfaceVariant
                    )
                }
            }
        }
    }
}

/**
 * Chat filter enum for forwarding modal
 */
enum class ChatFilter(val displayName: String) {
    ALL("All"),
    DIRECT("Direct"),
    GROUPS("Groups"),
    CHANNELS("Channels")
}

/**
 * Forwarding Progress Dialog - Shows progress when forwarding to multiple chats
 */
@Composable
fun ForwardingProgressDialog(
    isVisible: Boolean,
    currentProgress: Int,
    totalChats: Int,
    currentChatName: String,
    onDismiss: () -> Unit,
    modifier: Modifier = Modifier
) {
    if (isVisible) {
        AlertDialog(
            onDismissRequest = { /* Prevent dismissal during forwarding */ },
            modifier = modifier.testTag("forwarding-progress-dialog"),
            title = {
                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        Icons.Default.Forward,
                        contentDescription = null,
                        tint = TchatColors.primary,
                        modifier = Modifier.size(24.dp)
                    )
                    Spacer(modifier = Modifier.width(TchatSpacing.sm))
                    Text(
                        text = "Forwarding Message",
                        style = MaterialTheme.typography.headlineSmall
                    )
                }
            },
            text = {
                Column {
                    Text(
                        text = "Sending to $currentChatName...",
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.onSurface
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.md))

                    LinearProgressIndicator(
                        progress = { currentProgress.toFloat() / totalChats.toFloat() },
                        modifier = Modifier.fillMaxWidth(),
                        color = TchatColors.primary
                    )

                    Spacer(modifier = Modifier.height(TchatSpacing.sm))

                    Text(
                        text = "$currentProgress of $totalChats chats",
                        style = MaterialTheme.typography.bodySmall,
                        color = TchatColors.onSurfaceVariant
                    )
                }
            },
            confirmButton = {
                if (currentProgress >= totalChats) {
                    Button(onClick = onDismiss) {
                        Text("Done")
                    }
                }
            }
        )
    }
}

// Helper function to format timestamp
private fun formatTimestamp(timestamp: String): String {
    try {
        // Parse ISO timestamp or simple format
        val messageTime = if (timestamp.contains("T")) {
            // ISO format like "2024-01-20T10:30:00Z"
            timestamp.substringBefore("T").split("-").let { parts ->
                if (parts.size >= 3) {
                    val hour = timestamp.substringAfter("T").substringBefore(":").toIntOrNull() ?: 0
                    val minute = timestamp.substringAfter("T").substringAfter(":").substringBefore(":").toIntOrNull() ?: 0
                    Triple(parts[2].toIntOrNull() ?: 1, hour, minute)
                } else Triple(1, 0, 0)
            }
        } else {
            // Simple format, treat as current time
            Triple(1, 0, 0)
        }

        // For demo purposes, return relative time based on patterns
        return when {
            timestamp.contains("just_now") || timestamp.isEmpty() -> "Just now"
            timestamp.contains("1m") -> "1m ago"
            timestamp.contains("5m") -> "5m ago"
            timestamp.contains("30m") -> "30m ago"
            timestamp.contains("1h") -> "1h ago"
            timestamp.contains("today") -> "Today ${messageTime.second}:${messageTime.third.toString().padStart(2, '0')}"
            timestamp.contains("yesterday") -> "Yesterday"
            else -> {
                // Format time for display (simplified)
                val hour = messageTime.second
                val minute = messageTime.third
                val amPm = if (hour >= 12) "PM" else "AM"
                val displayHour = if (hour == 0) 12 else if (hour > 12) hour - 12 else hour
                "${displayHour}:${minute.toString().padStart(2, '0')} $amPm"
            }
        }
    } catch (e: Exception) {
        // Fallback for any parsing errors
        return "Now"
    }
}

// Message Type Components

@Composable
fun ImageMessage(
    message: Message,
    modifier: Modifier = Modifier
) {
    Column(modifier = modifier) {
        if (message.attachments.isNotEmpty()) {
            val attachment = message.attachments.first()
            Card(
                modifier = Modifier
                    .size(200.dp, 150.dp)
                    .clip(RoundedCornerShape(8.dp)),
                colors = CardDefaults.cardColors(containerColor = TchatColors.surfaceVariant)
            ) {
                Box(
                    contentAlignment = Alignment.Center,
                    modifier = Modifier.fillMaxSize()
                ) {
                    Column(
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        Icon(
                            Icons.Default.Image,
                            contentDescription = "Image",
                            modifier = Modifier.size(48.dp),
                            tint = TchatColors.primary
                        )
                        Spacer(modifier = Modifier.height(8.dp))
                        Text(
                            text = attachment.filename ?: "Image",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }
            }
            attachment.caption?.let { caption ->
                Spacer(modifier = Modifier.height(8.dp))
                Text(
                    text = caption,
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurface
                )
            }
        } else {
            Text(
                text = "📷 Image: ${message.content}",
                style = MaterialTheme.typography.bodyMedium,
                color = TchatColors.onSurface
            )
        }
    }
}

@Composable
fun VideoMessage(
    message: Message,
    modifier: Modifier = Modifier
) {
    Column(modifier = modifier) {
        if (message.attachments.isNotEmpty()) {
            val attachment = message.attachments.first()
            Card(
                modifier = Modifier
                    .size(250.dp, 150.dp)
                    .clip(RoundedCornerShape(8.dp)),
                colors = CardDefaults.cardColors(containerColor = TchatColors.surfaceVariant)
            ) {
                Box(
                    contentAlignment = Alignment.Center,
                    modifier = Modifier.fillMaxSize()
                ) {
                    Column(
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        Icon(
                            Icons.Default.PlayArrow,
                            contentDescription = "Video",
                            modifier = Modifier.size(64.dp),
                            tint = TchatColors.primary
                        )
                        Spacer(modifier = Modifier.height(8.dp))
                        Text(
                            text = attachment.filename ?: "Video",
                            style = MaterialTheme.typography.bodySmall,
                            color = TchatColors.onSurfaceVariant
                        )
                        attachment.duration?.let { duration ->
                            Text(
                                text = "${duration / 1000}s",
                                style = MaterialTheme.typography.labelSmall,
                                color = TchatColors.onSurfaceVariant
                            )
                        }
                    }
                }
            }
            attachment.caption?.let { caption ->
                Spacer(modifier = Modifier.height(8.dp))
                Text(
                    text = caption,
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurface
                )
            }
        } else {
            Text(
                text = "🎥 Video: ${message.content}",
                style = MaterialTheme.typography.bodyMedium,
                color = TchatColors.onSurface
            )
        }
    }
}

@Composable
fun AudioMessage(
    message: Message,
    modifier: Modifier = Modifier
) {
    var isPlaying by remember { mutableStateOf(false) }

    Card(
        modifier = modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(16.dp)),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surfaceVariant)
    ) {
        Row(
            modifier = Modifier
                .padding(TchatSpacing.md)
                .fillMaxWidth(),
            verticalAlignment = Alignment.CenterVertically
        ) {
            IconButton(
                onClick = { isPlaying = !isPlaying }
            ) {
                Icon(
                    if (isPlaying) Icons.Default.Pause else Icons.Default.PlayArrow,
                    contentDescription = if (isPlaying) "Pause" else "Play",
                    tint = TchatColors.primary
                )
            }

            Spacer(modifier = Modifier.width(TchatSpacing.sm))

            Column(
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = "Voice Message",
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurfaceVariant
                )
                if (message.attachments.isNotEmpty()) {
                    message.attachments.first().duration?.let { duration ->
                        Text(
                            text = "${duration / 1000}s",
                            style = MaterialTheme.typography.labelSmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }
            }

            Icon(
                Icons.Default.GraphicEq,
                contentDescription = "Audio waveform",
                tint = TchatColors.primary.copy(alpha = 0.6f)
            )
        }
    }
}

@Composable
fun FileMessage(
    message: Message,
    modifier: Modifier = Modifier
) {
    if (message.attachments.isNotEmpty()) {
        val attachment = message.attachments.first()
        Card(
            modifier = modifier
                .fillMaxWidth()
                .clip(RoundedCornerShape(8.dp)),
            colors = CardDefaults.cardColors(containerColor = TchatColors.surfaceVariant)
        ) {
            Row(
                modifier = Modifier
                    .padding(TchatSpacing.md)
                    .fillMaxWidth(),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Icon(
                    Icons.Default.Description,
                    contentDescription = "File",
                    modifier = Modifier.size(40.dp),
                    tint = TchatColors.primary
                )

                Spacer(modifier = Modifier.width(TchatSpacing.md))

                Column(
                    modifier = Modifier.weight(1f)
                ) {
                    Text(
                        text = attachment.filename ?: "File",
                        style = MaterialTheme.typography.bodyMedium,
                        color = TchatColors.onSurface,
                        maxLines = 1,
                        overflow = TextOverflow.Ellipsis
                    )
                    attachment.fileSize?.let { size ->
                        Text(
                            text = formatFileSize(size),
                            style = MaterialTheme.typography.labelSmall,
                            color = TchatColors.onSurfaceVariant
                        )
                    }
                }

                IconButton(onClick = { /* Download file */ }) {
                    Icon(
                        Icons.Default.Download,
                        contentDescription = "Download",
                        tint = TchatColors.primary
                    )
                }
            }
        }
    } else {
        Text(
            text = "📄 File: ${message.content}",
            style = MaterialTheme.typography.bodyMedium,
            color = TchatColors.onSurface
        )
    }
}

@Composable
fun LocationMessage(
    message: Message,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(8.dp)),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surfaceVariant)
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            Row(
                verticalAlignment = Alignment.CenterVertically
            ) {
                Icon(
                    Icons.Default.LocationOn,
                    contentDescription = "Location",
                    modifier = Modifier.size(24.dp),
                    tint = TchatColors.primary
                )
                Spacer(modifier = Modifier.width(TchatSpacing.sm))
                Text(
                    text = "Location Shared",
                    style = MaterialTheme.typography.bodyMedium,
                    fontWeight = FontWeight.Medium,
                    color = TchatColors.onSurface
                )
            }

            Spacer(modifier = Modifier.height(TchatSpacing.sm))

            // Mock map preview
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(120.dp)
                    .clip(RoundedCornerShape(8.dp)),
                contentAlignment = Alignment.Center
            ) {
                Surface(
                    color = MaterialTheme.colorScheme.primary,
                    modifier = Modifier.fillMaxSize()
                ) {
                    Column(
                        modifier = Modifier.fillMaxSize(),
                        horizontalAlignment = Alignment.CenterHorizontally,
                        verticalArrangement = Arrangement.Center
                    ) {
                        Icon(
                            Icons.Default.Map,
                            contentDescription = "Map",
                            modifier = Modifier.size(32.dp),
                            tint = MaterialTheme.colorScheme.onPrimary
                        )
                        Text(
                            text = message.content,
                            style = MaterialTheme.typography.bodySmall,
                            color = MaterialTheme.colorScheme.onPrimary,
                            textAlign = androidx.compose.ui.text.style.TextAlign.Center
                        )
                    }
                }
            }
        }
    }
}

@Composable
fun StickerMessage(
    message: Message,
    modifier: Modifier = Modifier
) {
    Box(
        modifier = modifier
            .size(120.dp)
            .clip(RoundedCornerShape(16.dp)),
        contentAlignment = Alignment.Center
    ) {
        Surface(
            color = TchatColors.surfaceVariant,
            modifier = Modifier.fillMaxSize()
        ) {
            Column(
                modifier = Modifier.fillMaxSize(),
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.Center
            ) {
                Text(
                    text = "\uD83D\uDE00", // Placeholder sticker emoji
                    style = MaterialTheme.typography.displayMedium
                )
                Text(
                    text = "Sticker",
                    style = MaterialTheme.typography.labelSmall,
                    color = TchatColors.onSurfaceVariant
                )
            }
        }
    }
}

@Composable
fun SystemMessage(
    message: Message,
    modifier: Modifier = Modifier
) {
    Box(
        modifier = modifier.fillMaxWidth(),
        contentAlignment = Alignment.Center
    ) {
        Surface(
            color = TchatColors.surfaceVariant.copy(alpha = 0.7f),
            shape = RoundedCornerShape(12.dp)
        ) {
            Text(
                text = message.content,
                style = MaterialTheme.typography.bodySmall,
                color = TchatColors.onSurfaceVariant,
                modifier = Modifier.padding(horizontal = 12.dp, vertical = 6.dp),
                textAlign = androidx.compose.ui.text.style.TextAlign.Center
            )
        }
    }
}

// Helper function
private fun formatFileSize(bytes: Long): String {
    val kb = bytes / 1024.0
    val mb = kb / 1024.0
    val gb = mb / 1024.0

    return when {
        gb >= 1 -> "%.1f GB".format(gb)
        mb >= 1 -> "%.1f MB".format(mb)
        kb >= 1 -> "%.1f KB".format(kb)
        else -> "$bytes B"
    }
}