package com.tchat.mobile.components.chat

import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.models.ChatMessageType

/**
 * Central message renderer that delegates to specific message type components
 */
@Composable
fun MessageRenderer(
    messageType: ChatMessageType,
    content: String,
    isFromMe: Boolean,
    modifier: Modifier = Modifier
) {
    when (messageType) {
        ChatMessageType.TEXT -> {
            Text(
                text = content,
                style = MaterialTheme.typography.bodyMedium,
                color = if (isFromMe) TchatColors.onPrimary else TchatColors.onSurface,
                modifier = modifier
            )
        }
        ChatMessageType.IMAGE -> {
            ImageMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.VIDEO -> {
            VideoMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.AUDIO -> {
            AudioMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.FILE -> {
            FileMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.LOCATION -> {
            LocationMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.PAYMENT -> {
            PaymentMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.POLL -> {
            PollMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.FORM -> {
            FormMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.SYSTEM -> {
            SystemMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.STICKER -> {
            StickerMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.GIF -> {
            GifMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.CONTACT -> {
            ContactMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.EVENT -> {
            EventMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.EVENT_MESSAGE -> {
            EventMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.EMBED -> {
            EmbedMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.DELETED -> {
            DeletedMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.FILE_MESSAGE -> {
            FileMessage(
                content = content,
                modifier = modifier
            )
        }
        ChatMessageType.LOCATION_MESSAGE -> {
            LocationMessage(
                content = content,
                modifier = modifier
            )
        }
        // Handle any remaining message types with default text rendering
        else -> {
            Text(
                text = content,
                style = MaterialTheme.typography.bodyMedium,
                color = if (isFromMe) TchatColors.onPrimary else TchatColors.onSurface,
                modifier = modifier
            )
        }
    }
}