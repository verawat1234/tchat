package com.tchat.mobile.data.models

import kotlinx.datetime.Instant
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.JsonElement
import kotlinx.serialization.json.JsonObject
import com.benasher44.uuid.uuid4

/**
 * Message represents a message in a dialog/conversation with rich content support
 */
@Serializable
data class Message(
    val id: String = uuid4().toString(),
    val dialogId: String = "",
    val senderId: String = "",
    val type: MessageType = MessageType.TEXT,
    val content: MessageContent = MessageContent(),
    val replyToId: String? = null,
    val isEdited: Boolean = false,
    val isPinned: Boolean = false,
    val isDeleted: Boolean = false,
    val mentions: List<String> = emptyList(),
    val reactions: Map<String, List<String>> = emptyMap(), // emoji -> user IDs
    val createdAt: Instant = kotlinx.datetime.Clock.System.now(),
    val editedAt: Instant? = null,
    val deletedAt: Instant? = null
) {
    /**
     * Add reaction from a user
     */
    fun addReaction(emoji: String, userId: String): Message {
        val currentReactions = reactions.toMutableMap()
        val userList = currentReactions[emoji]?.toMutableList() ?: mutableListOf()

        if (!userList.contains(userId)) {
            userList.add(userId)
            currentReactions[emoji] = userList
        }

        return copy(reactions = currentReactions)
    }

    /**
     * Remove reaction from a user
     */
    fun removeReaction(emoji: String, userId: String): Message {
        val currentReactions = reactions.toMutableMap()
        val userList = currentReactions[emoji]?.toMutableList()

        userList?.let { list ->
            list.remove(userId)
            if (list.isEmpty()) {
                currentReactions.remove(emoji)
            } else {
                currentReactions[emoji] = list
            }
        }

        return copy(reactions = currentReactions)
    }

    /**
     * Get total reaction count
     */
    fun getReactionCount(): Int = reactions.values.sumOf { it.size }

    /**
     * Check if user has reacted to this message
     */
    fun hasUserReacted(userId: String): Boolean {
        return reactions.values.any { it.contains(userId) }
    }

    /**
     * Check if user is mentioned in this message
     */
    fun isMentioned(userId: String): Boolean = mentions.contains(userId)

    /**
     * Soft delete the message
     */
    fun softDelete(): Message {
        return copy(
            isDeleted = true,
            deletedAt = kotlinx.datetime.Clock.System.now()
        )
    }

    /**
     * Convert to public message for API responses
     */
    fun toPublicMessage(forUserId: String? = null): Map<String, Any?> {
        val response = mutableMapOf<String, Any?>(
            "id" to id,
            "dialog_id" to dialogId,
            "sender_id" to senderId,
            "type" to type.name.lowercase(),
            "is_edited" to isEdited,
            "is_pinned" to isPinned,
            "created_at" to createdAt.toString()
        )

        // Include content if not deleted or if user is sender
        if (!isDeleted || senderId == forUserId) {
            response["content"] = content.data
            response["mentions"] = mentions
            response["reactions"] = reactions

            replyToId?.let { response["reply_to_id"] = it }

            if (isEdited) {
                editedAt?.let { response["edited_at"] = it.toString() }
            }
        } else {
            // Show deleted message placeholder
            response["content"] = mapOf("text" to "This message was deleted")
            response["is_deleted"] = true
        }

        return response
    }
}

/**
 * Message Type represents the type of message content
 */
@Serializable
enum class MessageType {
    TEXT,
    VOICE,
    FILE,
    IMAGE,
    VIDEO,
    PAYMENT,
    LOCATION,
    STICKER,
    SYSTEM;

    /**
     * Check if this message type requires content
     */
    fun requiresContent(): Boolean = this != SYSTEM

    /**
     * Check if this is a media message type
     */
    fun isMedia(): Boolean = this in listOf(VOICE, FILE, IMAGE, VIDEO)

    companion object {
        fun fromString(type: String): MessageType? {
            return values().find { it.name.lowercase() == type.lowercase() }
        }
    }
}

/**
 * Message Status represents the status of a message
 */
@Serializable
enum class MessageStatus {
    SENT,
    DELIVERED,
    READ,
    FAILED
}

/**
 * Message Content wrapper for different content types
 */
@Serializable
data class MessageContent(
    val data: Map<String, JsonElement> = emptyMap()
) {
    /**
     * Get text content
     */
    fun getTextContent(): TextContent? {
        return try {
            TextContent(
                text = data["text"]?.toString()?.removeSurrounding("\"") ?: "",
                entities = emptyList() // TODO: Parse entities from data
            )
        } catch (e: Exception) {
            null
        }
    }

    /**
     * Get voice content
     */
    fun getVoiceContent(): VoiceContent? {
        return try {
            VoiceContent(
                url = data["url"]?.toString()?.removeSurrounding("\"") ?: "",
                duration = data["duration"]?.toString()?.toIntOrNull() ?: 0,
                fileSize = data["file_size"]?.toString()?.toLongOrNull() ?: 0L,
                mimeType = data["mime_type"]?.toString()?.removeSurrounding("\"") ?: ""
            )
        } catch (e: Exception) {
            null
        }
    }

    /**
     * Get image content
     */
    fun getImageContent(): ImageContent? {
        return try {
            ImageContent(
                url = data["url"]?.toString()?.removeSurrounding("\"") ?: "",
                thumbnail = data["thumbnail"]?.toString()?.removeSurrounding("\""),
                width = data["width"]?.toString()?.toIntOrNull() ?: 0,
                height = data["height"]?.toString()?.toIntOrNull() ?: 0,
                fileSize = data["file_size"]?.toString()?.toLongOrNull() ?: 0L,
                caption = data["caption"]?.toString()?.removeSurrounding("\"")
            )
        } catch (e: Exception) {
            null
        }
    }

    /**
     * Get payment content
     */
    fun getPaymentContent(): PaymentContent? {
        return try {
            PaymentContent(
                amount = data["amount"]?.toString()?.toLongOrNull() ?: 0L,
                currency = data["currency"]?.toString()?.removeSurrounding("\"") ?: "THB",
                description = data["description"]?.toString()?.removeSurrounding("\"") ?: "",
                reference = data["reference"]?.toString()?.removeSurrounding("\""),
                status = data["status"]?.toString()?.removeSurrounding("\"") ?: "pending"
            )
        } catch (e: Exception) {
            null
        }
    }
}

/**
 * Text message content structure
 */
@Serializable
data class TextContent(
    val text: String,
    val entities: List<MessageEntity> = emptyList(),
    val metadata: Map<String, String> = emptyMap()
)

/**
 * Voice message content structure
 */
@Serializable
data class VoiceContent(
    val url: String,
    val duration: Int, // in milliseconds
    val waveform: List<Int> = emptyList(),
    val fileSize: Long,
    val mimeType: String
)

/**
 * File message content structure
 */
@Serializable
data class FileContent(
    val url: String,
    val filename: String,
    val fileSize: Long,
    val mimeType: String,
    val caption: String? = null
)

/**
 * Image message content structure
 */
@Serializable
data class ImageContent(
    val url: String,
    val thumbnail: String? = null,
    val width: Int,
    val height: Int,
    val fileSize: Long,
    val caption: String? = null
)

/**
 * Video message content structure
 */
@Serializable
data class VideoContent(
    val url: String,
    val thumbnail: String? = null,
    val duration: Int, // in milliseconds
    val width: Int,
    val height: Int,
    val fileSize: Long,
    val caption: String? = null
)

/**
 * Payment message content structure
 */
@Serializable
data class PaymentContent(
    val amount: Long, // in cents
    val currency: String,
    val description: String,
    val reference: String? = null,
    val status: String // pending, completed, failed, cancelled
) {
    companion object {
        val VALID_CURRENCIES = listOf("THB", "SGD", "IDR", "MYR", "PHP", "VND", "USD")
        val VALID_STATUSES = listOf("pending", "completed", "failed", "cancelled")
    }

    /**
     * Validate payment content
     */
    fun isValid(): Boolean {
        return amount > 0 &&
               VALID_CURRENCIES.contains(currency) &&
               VALID_STATUSES.contains(status) &&
               description.isNotBlank()
    }
}

/**
 * Location message content structure
 */
@Serializable
data class LocationContent(
    val latitude: Double,
    val longitude: Double,
    val address: String? = null,
    val venue: String? = null
) {
    /**
     * Validate location coordinates
     */
    fun isValid(): Boolean {
        return latitude in -90.0..90.0 && longitude in -180.0..180.0
    }
}

/**
 * Sticker message content structure
 */
@Serializable
data class StickerContent(
    val stickerId: String,
    val packId: String,
    val url: String,
    val width: Int,
    val height: Int
)

/**
 * System message content structure
 */
@Serializable
data class SystemContent(
    val type: String, // user_joined, user_left, name_changed, etc.
    val message: String, // Human-readable message
    val data: Map<String, String> = emptyMap() // Additional system data
)

/**
 * Message Entity represents entities within text messages (mentions, links, etc.)
 */
@Serializable
data class MessageEntity(
    val type: String, // mention, hashtag, url, email, phone, bold, italic, code
    val offset: Int, // Start position in text
    val length: Int, // Length of entity
    val url: String? = null, // For url entities
    val userId: String? = null // For mention entities
)

/**
 * Message Reply represents a reply to a message
 */
@Serializable
data class MessageReply(
    val messageId: String,
    val senderId: String,
    val content: String,
    val type: MessageType
)

/**
 * Message Forward represents a forwarded message
 */
@Serializable
data class MessageForward(
    val fromDialogId: String,
    val originalId: String,
    val senderId: String,
    val forwardedBy: String,
    val timestamp: Instant
)

/**
 * Message validation utility
 */
object MessageValidator {
    /**
     * Validate message based on its type
     */
    fun validateMessage(message: Message): List<String> {
        val errors = mutableListOf<String>()

        // Basic validation
        if (message.dialogId.isBlank()) errors.add("Dialog ID is required")
        if (message.senderId.isBlank()) errors.add("Sender ID is required")

        // Type-specific validation
        when (message.type) {
            MessageType.TEXT -> validateTextContent(message, errors)
            MessageType.VOICE -> validateVoiceContent(message, errors)
            MessageType.IMAGE -> validateImageContent(message, errors)
            MessageType.VIDEO -> validateVideoContent(message, errors)
            MessageType.PAYMENT -> validatePaymentContent(message, errors)
            MessageType.LOCATION -> validateLocationContent(message, errors)
            else -> { /* Other types have minimal validation */ }
        }

        return errors
    }

    private fun validateTextContent(message: Message, errors: MutableList<String>) {
        val textContent = message.content.getTextContent()
        if (textContent == null || textContent.text.isBlank()) {
            errors.add("Text content cannot be empty")
        } else if (textContent.text.length > 4096) {
            errors.add("Text content cannot exceed 4096 characters")
        }
    }

    private fun validateVoiceContent(message: Message, errors: MutableList<String>) {
        val voiceContent = message.content.getVoiceContent()
        if (voiceContent == null) {
            errors.add("Voice content is invalid")
            return
        }

        if (voiceContent.url.isBlank()) errors.add("Voice URL is required")
        if (voiceContent.duration <= 0 || voiceContent.duration > 300000) {
            errors.add("Voice duration must be between 1ms and 5 minutes")
        }
        if (voiceContent.fileSize <= 0 || voiceContent.fileSize > 50 * 1024 * 1024) {
            errors.add("Voice file size must be between 1 byte and 50MB")
        }
    }

    private fun validateImageContent(message: Message, errors: MutableList<String>) {
        val imageContent = message.content.getImageContent()
        if (imageContent == null) {
            errors.add("Image content is invalid")
            return
        }

        if (imageContent.url.isBlank()) errors.add("Image URL is required")
        if (imageContent.width <= 0 || imageContent.width > 8192) {
            errors.add("Image width must be between 1 and 8192 pixels")
        }
        if (imageContent.height <= 0 || imageContent.height > 8192) {
            errors.add("Image height must be between 1 and 8192 pixels")
        }
    }

    private fun validateVideoContent(message: Message, errors: MutableList<String>) {
        val videoContent = message.content.getVoiceContent() // Using voice content structure for now
        if (videoContent == null) {
            errors.add("Video content is invalid")
            return
        }

        if (videoContent.duration <= 0 || videoContent.duration > 1800000) {
            errors.add("Video duration must be between 1ms and 30 minutes")
        }
        if (videoContent.fileSize > 500 * 1024 * 1024) {
            errors.add("Video file size must be less than 500MB")
        }
    }

    private fun validatePaymentContent(message: Message, errors: MutableList<String>) {
        val paymentContent = message.content.getPaymentContent()
        if (paymentContent == null) {
            errors.add("Payment content is invalid")
            return
        }

        if (!paymentContent.isValid()) {
            errors.add("Payment content is invalid")
        }
    }

    private fun validateLocationContent(message: Message, errors: MutableList<String>) {
        // Location validation would go here
    }
}