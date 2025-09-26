package com.tchat.mobile.components.chat

import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.models.MessageType

/**
 * Central message renderer that delegates to specific message type components
 */
@Composable
fun MessageRenderer(
    messageType: MessageType,
    content: String,
    isFromMe: Boolean,
    modifier: Modifier = Modifier
) {
    when (messageType) {
        MessageType.TEXT -> {
            Text(
                text = content,
                style = MaterialTheme.typography.bodyMedium,
                color = if (isFromMe) TchatColors.onPrimary else TchatColors.onSurface,
                modifier = modifier
            )
        }
        MessageType.IMAGE -> {
            ImageMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.VIDEO -> {
            VideoMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.AUDIO -> {
            AudioMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.FILE -> {
            FileMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.LOCATION -> {
            LocationMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.PAYMENT -> {
            PaymentMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.POLL -> {
            PollMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.FORM -> {
            FormMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.SYSTEM -> {
            SystemMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.STICKER -> {
            StickerMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.GIF -> {
            GifMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.CONTACT -> {
            ContactMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.EVENT -> {
            EventMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.EVENT_MESSAGE -> {
            EventMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.EMBED -> {
            EmbedMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.DELETED -> {
            DeletedMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.FILE_MESSAGE -> {
            FileMessage(
                content = content,
                modifier = modifier
            )
        }
        MessageType.LOCATION_MESSAGE -> {
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