package com.tchat.mobile.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
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
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.Dialog
import com.tchat.mobile.components.*
import com.tchat.mobile.components.chat.MessageRenderer
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import com.tchat.mobile.models.ChatMessageType

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ChatDetailScreen(
    chatId: String,
    chatName: String = "Alice Johnson",
    onBackClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    var messageText by remember { mutableStateOf("") }
    var messages by remember { mutableStateOf(getMessagesForChat(chatId)) }
    var showAttachmentMenu by remember { mutableStateOf(false) }

    // Chat action states
    var showVideoCallDialog by remember { mutableStateOf(false) }
    var showAudioCallDialog by remember { mutableStateOf(false) }
    var showMoreOptionsMenu by remember { mutableStateOf(false) }
    var showToast by remember { mutableStateOf("") }
    var showConfirmDialog by remember { mutableStateOf<ConfirmDialogState?>(null) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Row(
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        // User Avatar
                        Box(
                            modifier = Modifier
                                .size(40.dp)
                                .clip(CircleShape)
                                .background(TchatColors.primary),
                            contentAlignment = Alignment.Center
                        ) {
                            Text(
                                text = chatName.first().toString(),
                                color = TchatColors.onPrimary,
                                style = MaterialTheme.typography.titleMedium,
                                fontWeight = FontWeight.Bold
                            )
                        }
                        Spacer(modifier = Modifier.width(TchatSpacing.md))
                        Column {
                            Text(
                                text = chatName,
                                style = MaterialTheme.typography.titleMedium,
                                fontWeight = FontWeight.SemiBold,
                                color = TchatColors.onSurface
                            )
                            Text(
                                text = "Online",
                                style = MaterialTheme.typography.bodySmall,
                                color = TchatColors.primary
                            )
                        }
                    }
                },
                navigationIcon = {
                    IconButton(onClick = onBackClick) {
                        Icon(
                            Icons.Default.ArrowBack,
                            contentDescription = "Back",
                            tint = TchatColors.onSurface
                        )
                    }
                },
                actions = {
                    IconButton(onClick = {
                        showVideoCallDialog = true
                    }) {
                        Icon(
                            Icons.Default.VideoCall,
                            contentDescription = "Video Call",
                            tint = TchatColors.onSurface
                        )
                    }
                    IconButton(onClick = {
                        showAudioCallDialog = true
                    }) {
                        Icon(
                            Icons.Default.Call,
                            contentDescription = "Audio Call",
                            tint = TchatColors.onSurface
                        )
                    }
                    IconButton(onClick = {
                        showMoreOptionsMenu = true
                    }) {
                        Icon(
                            Icons.Default.MoreVert,
                            contentDescription = "More Options",
                            tint = TchatColors.onSurface
                        )
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = TchatColors.surface,
                    titleContentColor = TchatColors.onSurface
                )
            )
        },
        bottomBar = {
            ChatInputBar(
                message = messageText,
                onMessageChange = { messageText = it },
                onSendMessage = {
                    if (messageText.isNotBlank()) {
                        // Add new message to the list
                        val newMessage = MessageItem(
                            id = (messages.size + 1).toString(),
                            content = messageText,
                            timestamp = "Now",
                            isFromMe = true,
                            status = MessageStatus.SENT,
                            messageType = ChatMessageType.TEXT
                        )
                        messages = listOf(newMessage) + messages
                        messageText = ""
                    }
                },
                onAttachClick = {
                    showAttachmentMenu = true
                },
                showAttachmentMenu = showAttachmentMenu,
                onDismissAttachmentMenu = { showAttachmentMenu = false },
                onShowToast = { message -> showToast = message }
            )
        }
    ) { paddingValues ->
        LazyColumn(
            modifier = Modifier
                .fillMaxWidth()
                .padding(paddingValues)
                .background(TchatColors.background),
            contentPadding = PaddingValues(TchatSpacing.md),
            verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm),
            reverseLayout = true
        ) {
            items(messages) { message ->
                MessageBubble(message = message)
            }
        }
    }

    // Chat Action Dialogs
    if (showVideoCallDialog) {
        ChatActionDialog(
            title = "Start Video Call",
            message = "Would you like to start a video call with $chatName?",
            icon = Icons.Default.VideoCall,
            onConfirm = {
                showVideoCallDialog = false
                showToast = "Starting video call with $chatName..."
                // Here you would navigate to video call screen or start video call logic
            },
            onDismiss = { showVideoCallDialog = false }
        )
    }

    if (showAudioCallDialog) {
        ChatActionDialog(
            title = "Start Audio Call",
            message = "Would you like to start an audio call with $chatName?",
            icon = Icons.Default.Call,
            onConfirm = {
                showAudioCallDialog = false
                showToast = "Starting audio call with $chatName..."
                // Here you would navigate to audio call screen or start audio call logic
            },
            onDismiss = { showAudioCallDialog = false }
        )
    }

    if (showMoreOptionsMenu) {
        MoreOptionsDialog(
            chatName = chatName,
            onDismiss = { showMoreOptionsMenu = false },
            onClearChat = {
                showConfirmDialog = ConfirmDialogState(
                    title = "Clear Chat",
                    message = "Are you sure you want to clear all messages in this chat? This action cannot be undone.",
                    confirmText = "Clear",
                    onConfirm = {
                        messages = emptyList()
                        showToast = "Chat cleared successfully"
                        showConfirmDialog = null
                    }
                )
                showMoreOptionsMenu = false
            },
            onBlockUser = {
                showConfirmDialog = ConfirmDialogState(
                    title = "Block User",
                    message = "Are you sure you want to block $chatName? They won't be able to send you messages.",
                    confirmText = "Block",
                    isDestructive = true,
                    onConfirm = {
                        showToast = "$chatName has been blocked"
                        showConfirmDialog = null
                    }
                )
                showMoreOptionsMenu = false
            },
            onMuteNotifications = {
                showToast = "Notifications muted for $chatName"
                showMoreOptionsMenu = false
            },
            onExportChat = {
                showToast = "Exporting chat messages..."
                showMoreOptionsMenu = false
            }
        )
    }

    showConfirmDialog?.let { dialogState ->
        ConfirmationDialog(
            state = dialogState,
            onDismiss = { showConfirmDialog = null }
        )
    }

    // Toast notifications
    if (showToast.isNotEmpty()) {
        LaunchedEffect(showToast) {
            kotlinx.coroutines.delay(2000)
            showToast = ""
        }
        // Simple toast simulation with Snackbar-like behavior
        Box(
            modifier = Modifier.fillMaxSize(),
            contentAlignment = Alignment.BottomCenter
        ) {
            Card(
                modifier = Modifier.padding(16.dp),
                colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
                elevation = CardDefaults.cardElevation(defaultElevation = 8.dp)
            ) {
                Text(
                    text = showToast,
                    modifier = Modifier.padding(16.dp),
                    color = TchatColors.onSurface
                )
            }
        }
    }
}

@Composable
private fun ChatInputBar(
    message: String,
    onMessageChange: (String) -> Unit,
    onSendMessage: () -> Unit,
    onAttachClick: () -> Unit,
    showAttachmentMenu: Boolean,
    onDismissAttachmentMenu: () -> Unit,
    onShowToast: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 4.dp)
    ) {
        Column {
            // Attachment menu
            if (showAttachmentMenu) {
                AttachmentMenu(
                    onDismiss = onDismissAttachmentMenu,
                    onPhotoClick = { onShowToast("Opening camera/gallery...") },
                    onVideoClick = { onShowToast("Opening video recorder...") },
                    onFileClick = { onShowToast("Opening file picker...") },
                    onLocationClick = { onShowToast("Sharing current location...") }
                )
            }

            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(TchatSpacing.md),
                verticalAlignment = Alignment.Bottom
            ) {
                IconButton(
                    onClick = onAttachClick,
                    modifier = Modifier.size(40.dp)
                ) {
                    Icon(
                        Icons.Default.Add,
                        contentDescription = "Attach File",
                        tint = TchatColors.onSurfaceVariant
                    )
                }

                TchatInput(
                    value = message,
                    onValueChange = onMessageChange,
                    type = TchatInputType.Multiline,
                    placeholder = "Type a message...",
                    maxLines = 4,
                    modifier = Modifier.weight(1f)
                )

                Spacer(modifier = Modifier.width(TchatSpacing.sm))

                TchatButton(
                    onClick = onSendMessage,
                    text = "Send",
                    variant = TchatButtonVariant.Primary,
                    enabled = message.isNotBlank()
                )
            }
        }
    }
}

@Composable
private fun AttachmentMenu(
    onDismiss: () -> Unit,
    onPhotoClick: () -> Unit,
    onVideoClick: () -> Unit,
    onFileClick: () -> Unit,
    onLocationClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier
            .fillMaxWidth()
            .padding(TchatSpacing.md),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surfaceVariant)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md),
            horizontalArrangement = Arrangement.SpaceEvenly
        ) {
            AttachmentOption(
                icon = Icons.Default.Image,
                label = "Photo",
                onClick = {
                    onPhotoClick()
                    onDismiss()
                }
            )
            AttachmentOption(
                icon = Icons.Default.VideoCall,
                label = "Video",
                onClick = {
                    onVideoClick()
                    onDismiss()
                }
            )
            AttachmentOption(
                icon = Icons.Default.Description,
                label = "File",
                onClick = {
                    onFileClick()
                    onDismiss()
                }
            )
            AttachmentOption(
                icon = Icons.Default.LocationOn,
                label = "Location",
                onClick = {
                    onLocationClick()
                    onDismiss()
                }
            )
        }
    }
}

@Composable
private fun AttachmentOption(
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    label: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        modifier = modifier
            .clickable { onClick() }
            .padding(TchatSpacing.sm)
    ) {
        Icon(
            icon,
            contentDescription = label,
            tint = TchatColors.primary,
            modifier = Modifier.size(32.dp)
        )
        Spacer(modifier = Modifier.height(4.dp))
        Text(
            text = label,
            style = MaterialTheme.typography.labelSmall,
            color = TchatColors.onSurfaceVariant
        )
    }
}

@Composable
private fun MessageBubble(
    message: MessageItem,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier.fillMaxWidth(),
        horizontalAlignment = if (message.isFromMe) Alignment.End else Alignment.Start
    ) {
        Card(
            modifier = Modifier.widthIn(max = 280.dp),
            shape = RoundedCornerShape(
                topStart = 16.dp,
                topEnd = 16.dp,
                bottomStart = if (message.isFromMe) 16.dp else 4.dp,
                bottomEnd = if (message.isFromMe) 4.dp else 16.dp
            ),
            colors = CardDefaults.cardColors(
                containerColor = if (message.isFromMe) TchatColors.primary else TchatColors.surfaceDim
            )
        ) {
            Column(
                modifier = Modifier.padding(TchatSpacing.md)
            ) {
                // Render content using organized message components
                MessageRenderer(
                    messageType = message.messageType,
                    content = message.content,
                    isFromMe = message.isFromMe
                )

                Spacer(modifier = Modifier.height(4.dp))

                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        text = message.timestamp,
                        style = MaterialTheme.typography.labelSmall,
                        color = if (message.isFromMe)
                            TchatColors.onPrimary.copy(alpha = 0.7f)
                        else
                            TchatColors.onSurface.copy(alpha = 0.7f)
                    )

                    if (message.isFromMe) {
                        Spacer(modifier = Modifier.width(4.dp))
                        Icon(
                            when (message.status) {
                                MessageStatus.SENT -> Icons.Default.Done
                                MessageStatus.DELIVERED -> Icons.Default.Done
                                MessageStatus.READ -> Icons.Default.Done
                            },
                            contentDescription = message.status.name,
                            modifier = Modifier.size(16.dp),
                            tint = if (message.status == MessageStatus.READ)
                                TchatColors.primary
                            else
                                TchatColors.onPrimary.copy(alpha = 0.7f)
                        )
                    }
                }
            }
        }
    }
}

// Dialog Components

@Composable
private fun ChatActionDialog(
    title: String,
    message: String,
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    onConfirm: () -> Unit,
    onDismiss: () -> Unit
) {
    Dialog(onDismissRequest = onDismiss) {
        Card(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            shape = RoundedCornerShape(16.dp),
            colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
        ) {
            Column(
                modifier = Modifier.padding(24.dp),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Box(
                    modifier = Modifier
                        .size(64.dp)
                        .background(TchatColors.primary, CircleShape),
                    contentAlignment = Alignment.Center
                ) {
                    Icon(
                        imageVector = icon,
                        contentDescription = null,
                        modifier = Modifier.size(32.dp),
                        tint = TchatColors.onPrimary
                    )
                }

                Spacer(modifier = Modifier.height(16.dp))

                Text(
                    text = title,
                    style = MaterialTheme.typography.headlineSmall,
                    fontWeight = FontWeight.Bold,
                    color = TchatColors.onSurface,
                    textAlign = TextAlign.Center
                )

                Spacer(modifier = Modifier.height(8.dp))

                Text(
                    text = message,
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurfaceVariant,
                    textAlign = TextAlign.Center
                )

                Spacer(modifier = Modifier.height(24.dp))

                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    TchatButton(
                        onClick = onDismiss,
                        text = "Cancel",
                        variant = TchatButtonVariant.Secondary,
                        modifier = Modifier.weight(1f)
                    )
                    TchatButton(
                        onClick = onConfirm,
                        text = "Start Call",
                        variant = TchatButtonVariant.Primary,
                        modifier = Modifier.weight(1f)
                    )
                }
            }
        }
    }
}

@Composable
private fun MoreOptionsDialog(
    chatName: String,
    onDismiss: () -> Unit,
    onClearChat: () -> Unit,
    onBlockUser: () -> Unit,
    onMuteNotifications: () -> Unit,
    onExportChat: () -> Unit
) {
    Dialog(onDismissRequest = onDismiss) {
        Card(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            shape = RoundedCornerShape(16.dp),
            colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
        ) {
            Column(
                modifier = Modifier.padding(8.dp)
            ) {
                Text(
                    text = "More Options",
                    style = MaterialTheme.typography.titleLarge,
                    fontWeight = FontWeight.Bold,
                    color = TchatColors.onSurface,
                    modifier = Modifier.padding(16.dp)
                )

                MoreOptionItem(
                    icon = Icons.Default.VolumeOff,
                    text = "Mute Notifications",
                    onClick = onMuteNotifications
                )

                MoreOptionItem(
                    icon = Icons.Default.FileDownload,
                    text = "Export Chat",
                    onClick = onExportChat
                )

                MoreOptionItem(
                    icon = Icons.Default.Delete,
                    text = "Clear Chat",
                    onClick = onClearChat
                )

                MoreOptionItem(
                    icon = Icons.Default.Block,
                    text = "Block $chatName",
                    onClick = onBlockUser,
                    isDestructive = true
                )
            }
        }
    }
}

@Composable
private fun MoreOptionItem(
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    text: String,
    onClick: () -> Unit,
    isDestructive: Boolean = false
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable { onClick() }
            .padding(16.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Icon(
            imageVector = icon,
            contentDescription = null,
            modifier = Modifier.size(24.dp),
            tint = if (isDestructive) TchatColors.error else TchatColors.onSurface
        )
        Spacer(modifier = Modifier.width(16.dp))
        Text(
            text = text,
            style = MaterialTheme.typography.bodyLarge,
            color = if (isDestructive) TchatColors.error else TchatColors.onSurface
        )
    }
}

data class ConfirmDialogState(
    val title: String,
    val message: String,
    val confirmText: String = "Confirm",
    val cancelText: String = "Cancel",
    val isDestructive: Boolean = false,
    val onConfirm: () -> Unit
)

@Composable
private fun ConfirmationDialog(
    state: ConfirmDialogState,
    onDismiss: () -> Unit
) {
    Dialog(onDismissRequest = onDismiss) {
        Card(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            shape = RoundedCornerShape(16.dp),
            colors = CardDefaults.cardColors(containerColor = TchatColors.surface)
        ) {
            Column(
                modifier = Modifier.padding(24.dp)
            ) {
                Text(
                    text = state.title,
                    style = MaterialTheme.typography.headlineSmall,
                    fontWeight = FontWeight.Bold,
                    color = TchatColors.onSurface
                )

                Spacer(modifier = Modifier.height(16.dp))

                Text(
                    text = state.message,
                    style = MaterialTheme.typography.bodyMedium,
                    color = TchatColors.onSurfaceVariant
                )

                Spacer(modifier = Modifier.height(24.dp))

                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    TchatButton(
                        onClick = onDismiss,
                        text = state.cancelText,
                        variant = TchatButtonVariant.Secondary,
                        modifier = Modifier.weight(1f)
                    )
                    TchatButton(
                        onClick = {
                            state.onConfirm()
                            onDismiss()
                        },
                        text = state.confirmText,
                        variant = if (state.isDestructive) TchatButtonVariant.Destructive else TchatButtonVariant.Primary,
                        modifier = Modifier.weight(1f)
                    )
                }
            }
        }
    }
}

// Data models
private data class MessageItem(
    val id: String,
    val content: String,
    val timestamp: String,
    val isFromMe: Boolean,
    val status: MessageStatus = MessageStatus.SENT,
    val messageType: ChatMessageType = ChatMessageType.TEXT
)

private enum class MessageStatus {
    SENT, DELIVERED, READ
}

private fun getMessagesForChat(chatId: String): List<MessageItem> = when (chatId) {
    "family" -> listOf(
        MessageItem("f8", "See you at dinner tomorrow! 🍽️", "4:15 PM", false, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("f7", "Photo shared: Family vacation memories", "4:10 PM", true, MessageStatus.READ, ChatMessageType.IMAGE),
        MessageItem("f6", "Love you too mom! ❤️", "4:05 PM", true, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("f5", "Hope you're having a great day, love you!", "4:00 PM", false, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("f4", "Voice message: 0:32 duration", "3:55 PM", false, MessageStatus.READ, ChatMessageType.AUDIO),
        MessageItem("f3", "Thanks for the birthday wishes! 🎂", "3:50 PM", true, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("f2", "Happy birthday! Hope you have an amazing day! 🎉", "3:45 PM", false, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("f1", "Good morning family! ☀️", "9:00 AM", true, MessageStatus.READ, ChatMessageType.TEXT)
    )

    "customer" -> listOf(
        MessageItem("c6", "Perfect! Thank you for the quick response.", "2:30 PM", false, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("c5", "Your order has been shipped! Tracking: TC123456789", "2:25 PM", true, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("c4", "When will my order be ready?", "2:20 PM", false, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("c3", "Payment confirmed: $85.99 for Premium Package", "2:15 PM", true, MessageStatus.READ, ChatMessageType.PAYMENT),
        MessageItem("c2", "Thank you! I'll proceed with the payment.", "2:10 PM", false, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("c1", "Hello! Welcome to our service. How can I help you today?", "2:05 PM", true, MessageStatus.READ, ChatMessageType.TEXT)
    )

    "lead_score_95" -> listOf(
        MessageItem("l6", "Proposal sent: Enterprise Solution Package", "5:45 PM", true, MessageStatus.READ, ChatMessageType.FILE_MESSAGE),
        MessageItem("l5", "This looks exactly what we need. Can you send a proposal?", "5:40 PM", false, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("l4", "We handle 50,000+ transactions daily with 99.9% uptime", "5:35 PM", true, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("l3", "What's your current transaction volume capacity?", "5:30 PM", false, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("l2", "I'm interested in your enterprise solutions for our fintech company", "5:25 PM", false, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("l1", "Hello! I saw your demo and I'm very impressed.", "5:20 PM", false, MessageStatus.READ, ChatMessageType.TEXT)
    )

    "ai_bot_assistant" -> listOf(
        MessageItem("a5", "Is there anything else I can help you with today?", "6:15 PM", false, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("a4", "Perfect! I've scheduled that for you. You'll receive a confirmation email shortly.", "6:10 PM", false, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("a3", "Yes, please schedule it for 2 PM tomorrow", "6:05 PM", true, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("a2", "I found an available slot tomorrow at 2 PM. Would you like me to book it?", "6:00 PM", false, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("a1", "Hi! I need to schedule a product demo this week", "5:55 PM", true, MessageStatus.READ, ChatMessageType.TEXT)
    )

    "project_team" -> listOf(
        MessageItem("p8", "Great work everyone! Let's deploy to staging tomorrow 🚀", "7:30 PM", false, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("p7", "All tests are passing now! Ready for review.", "7:25 PM", true, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("p6", "Code review complete. Just a few minor suggestions.", "7:20 PM", false, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("p5", "Document attached: Technical_Specification_v3.pdf", "7:15 PM", true, MessageStatus.READ, ChatMessageType.FILE),
        MessageItem("p4", "Sprint review meeting: Tomorrow 10 AM in Conference Room B", "7:10 PM", false, MessageStatus.READ, ChatMessageType.EVENT_MESSAGE),
        MessageItem("p3", "I'll have the API integration ready by EOD", "7:05 PM", true, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("p2", "The new feature is 80% complete. Working on final testing.", "7:00 PM", false, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("p1", "Daily standup: Backend API integration progress", "6:55 PM", true, MessageStatus.READ, ChatMessageType.TEXT)
    )

    else -> listOf(
        // Default messages for unknown chats
        MessageItem("d3", "Thanks for reaching out! 😊", "3:00 PM", false, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("d2", "Hello! How can I help you today?", "2:58 PM", true, MessageStatus.READ, ChatMessageType.TEXT),
        MessageItem("d1", "Hi there!", "2:55 PM", false, MessageStatus.READ, ChatMessageType.TEXT)
    )
}