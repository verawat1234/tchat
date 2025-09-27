package com.tchat.mobile.components.chat

import androidx.compose.animation.*
import androidx.compose.foundation.ExperimentalFoundationApi
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyRow
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
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.tchat.mobile.components.*
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import com.tchat.mobile.models.*

/**
 * Comprehensive message input system with attachment support, voice recording,
 * quick reactions, and advanced chat features
 */
@OptIn(ExperimentalFoundationApi::class)
@Composable
fun MessageInputSystem(
    message: String,
    onMessageChange: (String) -> Unit,
    onSendMessage: (MessageType, String, List<MessageAttachment>) -> Unit,
    replyToMessage: Message? = null,
    onCancelReply: () -> Unit = {},
    isTyping: Boolean = false,
    onTypingStart: () -> Unit = {},
    onTypingStop: () -> Unit = {},
    onVoiceRecordStart: () -> Unit = {},
    onVoiceRecordStop: () -> Unit = {},
    onVoiceRecordCancel: () -> Unit = {},
    modifier: Modifier = Modifier
) {
    var showAttachmentMenu by remember { mutableStateOf(false) }
    var showEmojiPicker by remember { mutableStateOf(false) }
    var isRecordingVoice by remember { mutableStateOf(false) }
    var recordingDuration by remember { mutableStateOf(0) }
    var selectedAttachments by remember { mutableStateOf<List<MessageAttachment>>(emptyList()) }

    Card(
        modifier = modifier
            .fillMaxWidth()
            .testTag("enhanced-message-input-container"),
        colors = CardDefaults.cardColors(containerColor = TchatColors.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 4.dp)
    ) {
        Column {
            // Reply indicator
            AnimatedVisibility(
                visible = replyToMessage != null,
                enter = slideInVertically() + fadeIn(),
                exit = slideOutVertically() + fadeOut()
            ) {
                replyToMessage?.let { reply ->
                    ReplyToIndicator(
                        message = reply,
                        onCancel = onCancelReply,
                        modifier = Modifier.testTag("reply-indicator")
                    )
                }
            }

            // Attachment previews
            AnimatedVisibility(
                visible = selectedAttachments.isNotEmpty(),
                enter = slideInVertically() + fadeIn(),
                exit = slideOutVertically() + fadeOut()
            ) {
                AttachmentPreviews(
                    attachments = selectedAttachments,
                    onRemoveAttachment = { attachment ->
                        selectedAttachments = selectedAttachments.filter { it.id != attachment.id }
                    },
                    modifier = Modifier.testTag("attachment-previews")
                )
            }

            // Enhanced attachment menu
            AnimatedVisibility(
                visible = showAttachmentMenu,
                enter = slideInVertically() + fadeIn(),
                exit = slideOutVertically() + fadeOut()
            ) {
                EnhancedAttachmentMenu(
                    onDismiss = { showAttachmentMenu = false },
                    onAttachmentSelected = { attachment ->
                        selectedAttachments = selectedAttachments + attachment
                        showAttachmentMenu = false
                    },
                    modifier = Modifier.testTag("enhanced-attachment-menu")
                )
            }

            // Emoji picker
            AnimatedVisibility(
                visible = showEmojiPicker,
                enter = slideInVertically() + fadeIn(),
                exit = slideOutVertically() + fadeOut()
            ) {
                EmojiPicker(
                    onEmojiSelected = { emoji ->
                        onMessageChange(message + emoji)
                        showEmojiPicker = false
                    },
                    onDismiss = { showEmojiPicker = false },
                    modifier = Modifier.testTag("emoji-picker")
                )
            }

            // Main input row
            if (isRecordingVoice) {
                VoiceRecordingInput(
                    duration = recordingDuration,
                    onStopRecording = {
                        isRecordingVoice = false
                        onVoiceRecordStop()
                    },
                    onCancelRecording = {
                        isRecordingVoice = false
                        onVoiceRecordCancel()
                    },
                    modifier = Modifier.testTag("voice-recording-input")
                )
            } else {
                MainInputRow(
                    message = message,
                    onMessageChange = { newMessage ->
                        onMessageChange(newMessage)
                        if (newMessage.isNotEmpty() && !isTyping) {
                            onTypingStart()
                        } else if (newMessage.isEmpty() && isTyping) {
                            onTypingStop()
                        }
                    },
                    onSendMessage = {
                        if (message.isNotBlank() || selectedAttachments.isNotEmpty()) {
                            val messageType = when {
                                selectedAttachments.any { it.type == AttachmentType.IMAGE } -> MessageType.IMAGE
                                selectedAttachments.any { it.type == AttachmentType.VIDEO } -> MessageType.VIDEO
                                selectedAttachments.any { it.type == AttachmentType.AUDIO } -> MessageType.AUDIO
                                selectedAttachments.any { it.type == AttachmentType.FILE } -> MessageType.FILE
                                selectedAttachments.any { it.type == AttachmentType.LOCATION } -> MessageType.LOCATION
                                else -> MessageType.TEXT
                            }
                            onSendMessage(messageType, message, selectedAttachments)
                            selectedAttachments = emptyList()
                            onTypingStop()
                        }
                    },
                    onAttachClick = { showAttachmentMenu = !showAttachmentMenu },
                    onEmojiClick = { showEmojiPicker = !showEmojiPicker },
                    onVoiceRecordStart = {
                        isRecordingVoice = true
                        onVoiceRecordStart()
                    },
                    hasAttachments = selectedAttachments.isNotEmpty(),
                    modifier = Modifier.testTag("main-input-row")
                )
            }
        }
    }
}

@Composable
private fun ReplyToIndicator(
    message: Message,
    onCancel: () -> Unit,
    modifier: Modifier = Modifier
) {
    Surface(
        modifier = modifier
            .fillMaxWidth()
            .padding(TchatSpacing.md),
        color = MaterialTheme.colorScheme.primary.copy(alpha = 0.1f),
        shape = RoundedCornerShape(8.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(TchatSpacing.md),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Icon(
                Icons.Default.Reply,
                contentDescription = "Replying to",
                modifier = Modifier.size(16.dp),
                tint = TchatColors.primary
            )

            Spacer(modifier = Modifier.width(TchatSpacing.sm))

            Column(
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = "Replying to ${message.senderName}",
                    style = MaterialTheme.typography.labelMedium,
                    color = TchatColors.primary,
                    fontWeight = FontWeight.Medium
                )
                Text(
                    text = message.getDisplayContent(),
                    style = MaterialTheme.typography.bodySmall,
                    color = TchatColors.onSurfaceVariant,
                    maxLines = 1
                )
            }

            IconButton(
                onClick = onCancel,
                modifier = Modifier.size(24.dp)
            ) {
                Icon(
                    Icons.Default.Close,
                    contentDescription = "Cancel reply",
                    modifier = Modifier.size(16.dp),
                    tint = TchatColors.onSurfaceVariant
                )
            }
        }
    }
}

@Composable
private fun AttachmentPreviews(
    attachments: List<MessageAttachment>,
    onRemoveAttachment: (MessageAttachment) -> Unit,
    modifier: Modifier = Modifier
) {
    LazyRow(
        modifier = modifier
            .fillMaxWidth()
            .padding(horizontal = TchatSpacing.md),
        horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
    ) {
        items(attachments) { attachment ->
            AttachmentPreview(
                attachment = attachment,
                onRemove = { onRemoveAttachment(attachment) }
            )
        }
    }
}

@Composable
private fun AttachmentPreview(
    attachment: MessageAttachment,
    onRemove: () -> Unit,
    modifier: Modifier = Modifier
) {
    Box(
        modifier = modifier
            .size(80.dp)
            .clip(RoundedCornerShape(8.dp))
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
                Icon(
                    when (attachment.type) {
                        AttachmentType.IMAGE -> Icons.Default.Image
                        AttachmentType.VIDEO -> Icons.Default.VideoFile
                        AttachmentType.AUDIO -> Icons.Default.AudioFile
                        AttachmentType.FILE -> Icons.Default.Description
                        AttachmentType.LOCATION -> Icons.Default.LocationOn
                    },
                    contentDescription = attachment.type.name,
                    modifier = Modifier.size(24.dp),
                    tint = TchatColors.primary
                )
                attachment.filename?.let { filename ->
                    Text(
                        text = filename.take(10) + if (filename.length > 10) "..." else "",
                        style = MaterialTheme.typography.labelSmall,
                        color = TchatColors.onSurfaceVariant,
                        textAlign = TextAlign.Center
                    )
                }
            }
        }

        // Remove button
        Surface(
            modifier = Modifier
                .align(Alignment.TopEnd)
                .offset(4.dp, (-4).dp)
                .size(20.dp)
                .clip(CircleShape)
                .clickable { onRemove() },
            color = TchatColors.error
        ) {
            Icon(
                Icons.Default.Close,
                contentDescription = "Remove",
                modifier = Modifier
                    .fillMaxSize()
                    .padding(2.dp),
                tint = MaterialTheme.colorScheme.onError
            )
        }
    }
}

@Composable
private fun EnhancedAttachmentMenu(
    onDismiss: () -> Unit,
    onAttachmentSelected: (MessageAttachment) -> Unit,
    modifier: Modifier = Modifier
) {
    Surface(
        modifier = modifier
            .fillMaxWidth()
            .padding(TchatSpacing.md),
        color = TchatColors.surfaceVariant,
        shape = RoundedCornerShape(16.dp)
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            Text(
                text = "Attachments",
                style = MaterialTheme.typography.titleMedium,
                color = TchatColors.onSurfaceVariant,
                fontWeight = FontWeight.Medium,
                modifier = Modifier.padding(bottom = TchatSpacing.md)
            )

            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceEvenly
            ) {
                EnhancedAttachmentOption(
                    icon = Icons.Default.PhotoCamera,
                    label = "Camera",
                    onClick = {
                        // Mock camera attachment
                        onAttachmentSelected(
                            MessageAttachment(
                                id = "camera_${System.currentTimeMillis()}",
                                type = AttachmentType.IMAGE,
                                url = "camera://capture",
                                filename = "photo.jpg"
                            )
                        )
                    }
                )

                EnhancedAttachmentOption(
                    icon = Icons.Default.PhotoLibrary,
                    label = "Gallery",
                    onClick = {
                        // Mock gallery attachment
                        onAttachmentSelected(
                            MessageAttachment(
                                id = "gallery_${System.currentTimeMillis()}",
                                type = AttachmentType.IMAGE,
                                url = "gallery://image",
                                filename = "image.jpg"
                            )
                        )
                    }
                )

                EnhancedAttachmentOption(
                    icon = Icons.Default.Videocam,
                    label = "Video",
                    onClick = {
                        onAttachmentSelected(
                            MessageAttachment(
                                id = "video_${System.currentTimeMillis()}",
                                type = AttachmentType.VIDEO,
                                url = "video://record",
                                filename = "video.mp4"
                            )
                        )
                    }
                )

                EnhancedAttachmentOption(
                    icon = Icons.Default.Description,
                    label = "File",
                    onClick = {
                        onAttachmentSelected(
                            MessageAttachment(
                                id = "file_${System.currentTimeMillis()}",
                                type = AttachmentType.FILE,
                                url = "file://document",
                                filename = "document.pdf"
                            )
                        )
                    }
                )
            }

            Spacer(modifier = Modifier.height(TchatSpacing.md))

            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceEvenly
            ) {
                EnhancedAttachmentOption(
                    icon = Icons.Default.LocationOn,
                    label = "Location",
                    onClick = {
                        onAttachmentSelected(
                            MessageAttachment(
                                id = "location_${System.currentTimeMillis()}",
                                type = AttachmentType.LOCATION,
                                url = "location://current",
                                filename = "Current Location"
                            )
                        )
                    }
                )

                EnhancedAttachmentOption(
                    icon = Icons.Default.ContactPage,
                    label = "Contact",
                    onClick = {
                        // Mock contact sharing
                        onDismiss()
                    }
                )

                EnhancedAttachmentOption(
                    icon = Icons.Default.Poll,
                    label = "Poll",
                    onClick = {
                        // Mock poll creation
                        onDismiss()
                    }
                )

                EnhancedAttachmentOption(
                    icon = Icons.Default.EventNote,
                    label = "Event",
                    onClick = {
                        // Mock event creation
                        onDismiss()
                    }
                )
            }
        }
    }
}

@Composable
private fun EnhancedAttachmentOption(
    icon: ImageVector,
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
        Surface(
            modifier = Modifier.size(48.dp),
            shape = CircleShape,
            color = TchatColors.primary.copy(alpha = 0.1f)
        ) {
            Icon(
                icon,
                contentDescription = label,
                modifier = Modifier
                    .fillMaxSize()
                    .padding(12.dp),
                tint = TchatColors.primary
            )
        }
        Spacer(modifier = Modifier.height(4.dp))
        Text(
            text = label,
            style = MaterialTheme.typography.labelSmall,
            color = TchatColors.onSurfaceVariant
        )
    }
}

@Composable
private fun EmojiPicker(
    onEmojiSelected: (String) -> Unit,
    onDismiss: () -> Unit,
    modifier: Modifier = Modifier
) {
    val emojis = listOf(
        "😀", "😃", "😄", "😁", "😆", "😅", "😂", "🤣",
        "😊", "😇", "🙂", "🙃", "😉", "😌", "😍", "🥰",
        "😘", "😗", "😙", "😚", "😋", "😛", "😝", "😜",
        "🤪", "🤨", "🧐", "🤓", "😎", "🤩", "🥳", "😏",
        "❤️", "🧡", "💛", "💚", "💙", "💜", "🖤", "🤍",
        "👍", "👎", "👌", "✌️", "🤞", "🤟", "🤘", "🤙",
        "🎉", "🎊", "🎈", "🎁", "🎀", "🎂", "🍰", "🧁"
    )

    Surface(
        modifier = modifier
            .fillMaxWidth()
            .padding(TchatSpacing.md),
        color = TchatColors.surfaceVariant,
        shape = RoundedCornerShape(16.dp)
    ) {
        Column(
            modifier = Modifier.padding(TchatSpacing.md)
        ) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(
                    text = "Emojis",
                    style = MaterialTheme.typography.titleMedium,
                    color = TchatColors.onSurfaceVariant,
                    fontWeight = FontWeight.Medium
                )
                IconButton(
                    onClick = onDismiss,
                    modifier = Modifier.size(24.dp)
                ) {
                    Icon(
                        Icons.Default.Close,
                        contentDescription = "Close",
                        tint = TchatColors.onSurfaceVariant
                    )
                }
            }

            // Emoji grid
            val chunkedEmojis = emojis.chunked(8)
            chunkedEmojis.forEach { row ->
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceEvenly
                ) {
                    row.forEach { emoji ->
                        Surface(
                            modifier = Modifier
                                .size(40.dp)
                                .clip(CircleShape)
                                .clickable { onEmojiSelected(emoji) },
                            color = TchatColors.surface
                        ) {
                            Box(
                                contentAlignment = Alignment.Center,
                                modifier = Modifier.fillMaxSize()
                            ) {
                                Text(
                                    text = emoji,
                                    style = MaterialTheme.typography.titleMedium
                                )
                            }
                        }
                    }
                }
                Spacer(modifier = Modifier.height(4.dp))
            }
        }
    }
}

@Composable
private fun MainInputRow(
    message: String,
    onMessageChange: (String) -> Unit,
    onSendMessage: () -> Unit,
    onAttachClick: () -> Unit,
    onEmojiClick: () -> Unit,
    onVoiceRecordStart: () -> Unit,
    hasAttachments: Boolean,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier
            .fillMaxWidth()
            .padding(TchatSpacing.md),
        verticalAlignment = Alignment.Bottom
    ) {
        // Attachment button
        IconButton(
            onClick = onAttachClick,
            modifier = Modifier
                .size(40.dp)
                .testTag("enhanced-attachment-button")
        ) {
            Icon(
                Icons.Default.AttachFile,
                contentDescription = "Attachments",
                tint = if (hasAttachments) TchatColors.primary else TchatColors.onSurfaceVariant
            )
        }

        // Emoji button
        IconButton(
            onClick = onEmojiClick,
            modifier = Modifier
                .size(40.dp)
                .testTag("emoji-button")
        ) {
            Icon(
                Icons.Default.EmojiEmotions,
                contentDescription = "Emojis",
                tint = TchatColors.onSurfaceVariant
            )
        }

        // Message input field
        TchatInput(
            value = message,
            onValueChange = onMessageChange,
            type = TchatInputType.Multiline,
            placeholder = "Type a message...",
            maxLines = 4,
            modifier = Modifier
                .weight(1f)
                .testTag("enhanced-message-field")
        )

        Spacer(modifier = Modifier.width(TchatSpacing.sm))

        // Send or Voice button
        if (message.isNotBlank() || hasAttachments) {
            TchatButton(
                onClick = onSendMessage,
                text = "Send",
                variant = TchatButtonVariant.Primary,
                modifier = Modifier.testTag("enhanced-send-button")
            )
        } else {
            IconButton(
                onClick = onVoiceRecordStart,
                modifier = Modifier
                    .size(48.dp)
                    .background(TchatColors.primary, CircleShape)
                    .testTag("voice-record-button")
            ) {
                Icon(
                    Icons.Default.Mic,
                    contentDescription = "Voice message",
                    tint = TchatColors.onPrimary
                )
            }
        }
    }
}

@Composable
private fun VoiceRecordingInput(
    duration: Int,
    onStopRecording: () -> Unit,
    onCancelRecording: () -> Unit,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier
            .fillMaxWidth()
            .padding(TchatSpacing.md)
            .background(MaterialTheme.colorScheme.error.copy(alpha = 0.1f), RoundedCornerShape(24.dp))
            .padding(TchatSpacing.md),
        verticalAlignment = Alignment.CenterVertically
    ) {
        // Cancel button
        IconButton(
            onClick = onCancelRecording,
            modifier = Modifier.testTag("voice-cancel-button")
        ) {
            Icon(
                Icons.Default.Close,
                contentDescription = "Cancel recording",
                tint = TchatColors.error
            )
        }

        // Recording indicator
        Row(
            modifier = Modifier.weight(1f),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.Center
        ) {
            // Animated recording dot
            Surface(
                modifier = Modifier.size(12.dp),
                shape = CircleShape,
                color = TchatColors.error
            ) {}

            Spacer(modifier = Modifier.width(8.dp))

            Text(
                text = "Recording... ${duration}s",
                style = MaterialTheme.typography.bodyMedium,
                color = TchatColors.onSurface
            )
        }

        // Stop/Send button
        IconButton(
            onClick = onStopRecording,
            modifier = Modifier.testTag("voice-stop-button")
        ) {
            Icon(
                Icons.Default.Send,
                contentDescription = "Send recording",
                tint = TchatColors.primary
            )
        }
    }
}