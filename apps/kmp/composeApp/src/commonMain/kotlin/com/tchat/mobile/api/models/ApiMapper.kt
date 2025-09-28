package com.tchat.mobile.api.models

import kotlinx.serialization.Serializable
import kotlinx.datetime.Clock

/**
 * Backend API Models that match the Go messaging service format
 * These are the actual DTOs expected by the backend
 */

@Serializable
data class BackendDialogDto(
    val id: String,
    val type: String, // "user", "group", "channel", "business"
    val name: String? = null,
    val title: String? = null, // alias for name
    val avatar: String? = null,
    val description: String? = null,
    val participants: List<String>,
    val participant_count: Int,
    val admin_ids: List<String> = emptyList(),
    val last_message_id: String? = null,
    val unread_count: Int = 0,
    val is_pinned: Boolean = false,
    val is_archived: Boolean = false,
    val is_muted: Boolean = false,
    val settings: BackendDialogSettingsDto = BackendDialogSettingsDto(),
    val metadata: Map<String, String>? = null,
    val created_at: String,
    val updated_at: String
)

@Serializable
data class BackendDialogSettingsDto(
    val notifications_enabled: Boolean = true,
    val sound_enabled: Boolean = true,
    val vibration_enabled: Boolean = true,
    val message_preview: Boolean = true
)

@Serializable
data class BackendMessageDto(
    val id: String,
    val dialog_id: String,
    val sender_id: String,
    val type: String, // "text", "voice", "file", "image", "video", "payment", "location"
    val content: BackendMessageContentDto,
    val media_url: String? = null,
    val thumbnail_url: String? = null,
    val status: String, // "sent", "delivered", "read", "failed"
    val metadata: Map<String, String>? = null,
    val reply_to_id: String? = null,
    val reply_to: BackendMessageReplyDto? = null,
    val is_edited: Boolean = false,
    val is_pinned: Boolean = false,
    val is_deleted: Boolean = false,
    val mentions: List<String> = emptyList(),
    val reactions: List<BackendReactionDto> = emptyList(),
    val sent_at: String,
    val created_at: String,
    val updated_at: String,
    val edited_at: String? = null,
    val deleted_at: String? = null
)

@Serializable
data class BackendMessageContentDto(
    val text: String? = null,
    val media_url: String? = null,
    val file_name: String? = null,
    val file_size: Long? = null,
    val mime_type: String? = null,
    val thumbnail_url: String? = null,
    val duration: Long? = null,
    val location: BackendLocationDto? = null,
    val payment: BackendPaymentDto? = null
)

@Serializable
data class BackendMessageReplyDto(
    val id: String,
    val sender_id: String,
    val content: String,
    val type: String
)

@Serializable
data class BackendReactionDto(
    val user_id: String,
    val emoji: String,
    val created_at: String
)

@Serializable
data class BackendLocationDto(
    val latitude: Double,
    val longitude: Double,
    val address: String? = null,
    val place_name: String? = null
)

@Serializable
data class BackendPaymentDto(
    val amount: Double,
    val currency: String,
    val description: String? = null,
    val transaction_id: String? = null
)

/**
 * Request DTOs for backend API
 */
@Serializable
data class SendMessageRequest(
    val type: String,
    val content: BackendMessageContentDto,
    val reply_to_id: String? = null,
    val mentions: List<String> = emptyList()
)

@Serializable
data class CreateDialogRequest(
    val type: String,
    val name: String? = null,
    val description: String? = null,
    val participants: List<String>,
    val settings: BackendDialogSettingsDto = BackendDialogSettingsDto()
)

/**
 * Mapper functions to convert between KMP DTOs and Backend DTOs
 */
object ApiMapper {

    fun chatSessionDtoFromBackend(backend: BackendDialogDto): ChatSessionDto {
        return ChatSessionDto(
            id = backend.id,
            name = backend.name ?: backend.title ?: "Chat",
            description = backend.description,
            avatar = backend.avatar,
            type = when (backend.type) {
                "user" -> "direct"
                "group" -> "group"
                "channel" -> "channel"
                "business" -> "support"
                else -> "group"
            },
            isActive = !backend.is_archived,
            participants = backend.participants,
            createdBy = backend.participants.firstOrNull() ?: "",
            createdAt = parseTimestamp(backend.created_at),
            updatedAt = parseTimestamp(backend.updated_at),
            lastMessageId = backend.last_message_id,
            lastMessageAt = null, // would need last message timestamp
            unreadCount = backend.unread_count,
            settings = ChatSettingsDto(
                isMuted = backend.is_muted,
                muteUntil = null,
                notificationsEnabled = backend.settings.notifications_enabled,
                messageRetention = 30,
                allowInvites = true,
                isPrivate = backend.type == "user"
            )
        )
    }

    fun messageDtoFromBackend(backend: BackendMessageDto): MessageDto {
        return MessageDto(
            id = backend.id,
            chatId = backend.dialog_id,
            senderId = backend.sender_id,
            senderName = "User", // would need user lookup
            senderAvatar = null,
            type = backend.type,
            content = backend.content.text ?: "",
            isEdited = backend.is_edited,
            isPinned = backend.is_pinned,
            isDeleted = backend.is_deleted,
            replyToId = backend.reply_to_id,
            reactions = backend.reactions.map { reaction ->
                ReactionDto(
                    messageId = backend.id,
                    userId = reaction.user_id,
                    emoji = reaction.emoji,
                    timestamp = parseTimestamp(reaction.created_at)
                )
            },
            attachments = emptyList(), // would need to map media
            createdAt = parseTimestamp(backend.created_at),
            editedAt = backend.edited_at?.let { parseTimestamp(it) },
            deletedAt = backend.deleted_at?.let { parseTimestamp(it) },
            serverTimestamp = parseTimestamp(backend.sent_at),
            version = 1,
            checksum = null,
            deliveryStatus = when (backend.status) {
                "sent" -> "sent"
                "delivered" -> "delivered"
                "read" -> "read"
                "failed" -> "failed"
                else -> "sent"
            },
            readBy = emptyList()
        )
    }

    fun messageDtoToBackendRequest(message: MessageDto): SendMessageRequest {
        return SendMessageRequest(
            type = message.type,
            content = BackendMessageContentDto(
                text = if (message.type == "text") message.content else null,
                media_url = message.attachments.firstOrNull()?.url,
                file_name = message.attachments.firstOrNull()?.filename,
                file_size = message.attachments.firstOrNull()?.fileSize,
                mime_type = message.attachments.firstOrNull()?.mimeType,
                thumbnail_url = message.attachments.firstOrNull()?.thumbnail,
                duration = message.attachments.firstOrNull()?.duration
            ),
            reply_to_id = message.replyToId,
            mentions = emptyList() // would need to extract mentions from content
        )
    }

    private fun parseTimestamp(timestamp: String): Long {
        return try {
            // Backend sends ISO 8601 timestamps, convert to epoch millis
            kotlinx.datetime.Instant.parse(timestamp).toEpochMilliseconds()
        } catch (e: Exception) {
            // Fallback to current time if parsing fails
            Clock.System.now().toEpochMilliseconds()
        }
    }

    fun formatTimestamp(epochMillis: Long): String {
        return kotlinx.datetime.Instant.fromEpochMilliseconds(epochMillis).toString()
    }
}