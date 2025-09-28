package com.tchat.mobile.models

import kotlinx.serialization.Serializable

/**
 * Message model for chat functionality
 *
 * Comprehensive message structure supporting rich content,
 * attachments, reactions, and real-time features.
 */
@Serializable
data class Message(
    val id: String,
    val chatId: String,
    val senderId: String,
    val senderName: String,
    val senderAvatar: String? = null,
    val type: MessageType,
    val content: String,
    val isEdited: Boolean = false,
    val isPinned: Boolean = false,
    val isDeleted: Boolean = false,
    val replyToId: String? = null,
    val reactions: List<MessageReaction> = emptyList(),
    val attachments: List<MessageAttachment> = emptyList(),
    val createdAt: String,
    val editedAt: String? = null,
    val deletedAt: String? = null,
    val deliveryStatus: DeliveryStatus = DeliveryStatus.SENT,
    val readBy: List<String> = emptyList() // User IDs who have read this message
)



@Serializable
data class MessageAttachment(
    val id: String,
    val type: AttachmentType,
    val url: String,
    val thumbnail: String? = null,
    val filename: String? = null,
    val fileSize: Long? = null,
    val mimeType: String? = null,
    val width: Int? = null,
    val height: Int? = null,
    val duration: Int? = null, // in milliseconds for audio/video
    val caption: String? = null,
    val metadata: Map<String, String> = emptyMap()
)

@Serializable
enum class AttachmentType {
    IMAGE,
    VIDEO,
    AUDIO,
    FILE,
    LOCATION
}

@Serializable
data class MessageReaction(
    val emoji: String,
    val userId: String,
    val userName: String,
    val timestamp: String
)

@Serializable
data class MessageReply(
    val messageId: String,
    val senderId: String,
    val senderName: String,
    val content: String,
    val type: MessageType,
    val timestamp: String
)

/**
 * Message utilities and extensions
 */
fun Message.isFromCurrentUser(currentUserId: String): Boolean = senderId == currentUserId

fun Message.hasAttachments(): Boolean = attachments.isNotEmpty()

fun Message.hasReactions(): Boolean = reactions.isNotEmpty()

fun Message.getReactionCount(): Int = reactions.size

fun Message.hasUserReacted(userId: String): Boolean = reactions.any { it.userId == userId }

fun Message.getUserReaction(userId: String): MessageReaction? = reactions.find { it.userId == userId }

fun Message.getDisplayContent(): String {
    return when {
        isDeleted -> "This message was deleted"
        type == MessageType.IMAGE && attachments.isNotEmpty() ->
            attachments.first().caption ?: "📷 Image"
        type == MessageType.VIDEO && attachments.isNotEmpty() ->
            attachments.first().caption ?: "🎥 Video"
        type == MessageType.AUDIO -> "🎵 Audio"
        type == MessageType.FILE && attachments.isNotEmpty() ->
            "📄 ${attachments.first().filename ?: "File"}"
        type == MessageType.LOCATION -> "📍 Location"
        type == MessageType.STICKER -> "Sticker"
        type == MessageType.SYSTEM -> content
        else -> content
    }
}

fun Message.isMediaMessage(): Boolean {
    return type in listOf(MessageType.IMAGE, MessageType.VIDEO, MessageType.AUDIO, MessageType.FILE)
}

fun Message.canBeEdited(currentUserId: String): Boolean {
    return senderId == currentUserId && !isDeleted && type == MessageType.TEXT
}

fun Message.canBeDeleted(currentUserId: String): Boolean {
    return senderId == currentUserId && !isDeleted
}

fun Message.canBeRepliedTo(): Boolean {
    return !isDeleted && type != MessageType.SYSTEM
}

/**
 * Message grouping utilities for chat UI
 */
fun List<Message>.groupByDate(): Map<String, List<Message>> {
    return groupBy { message ->
        // Extract date from timestamp (assuming ISO format)
        message.createdAt.split("T").firstOrNull() ?: message.createdAt
    }
}

fun List<Message>.groupConsecutiveMessages(): List<List<Message>> {
    if (isEmpty()) return emptyList()

    val groups = mutableListOf<MutableList<Message>>()
    var currentGroup = mutableListOf(first())

    for (i in 1 until size) {
        val current = this[i]
        val previous = this[i - 1]

        // Group if same sender and within 5 minutes
        if (current.senderId == previous.senderId &&
            shouldGroupMessages(previous, current)) {
            currentGroup.add(current)
        } else {
            groups.add(currentGroup)
            currentGroup = mutableListOf(current)
        }
    }

    groups.add(currentGroup)
    return groups
}

private fun shouldGroupMessages(previous: Message, current: Message): Boolean {
    // Simple time-based grouping (in a real app, you'd parse timestamps properly)
    return true // For now, always group consecutive messages from same sender
}

/**
 * Message search utilities
 */
fun List<Message>.searchByContent(query: String): List<Message> {
    if (query.isBlank()) return this

    val lowerQuery = query.lowercase()
    return filter { message ->
        message.getDisplayContent().lowercase().contains(lowerQuery) ||
        message.senderName.lowercase().contains(lowerQuery)
    }
}

fun List<Message>.filterByType(type: MessageType): List<Message> {
    return filter { it.type == type }
}

fun List<Message>.filterMediaMessages(): List<Message> {
    return filter { it.isMediaMessage() }
}

fun List<Message>.filterUnread(lastReadMessageId: String?): List<Message> {
    if (lastReadMessageId == null) return this

    val lastReadIndex = indexOfFirst { it.id == lastReadMessageId }
    return if (lastReadIndex >= 0) {
        drop(lastReadIndex + 1)
    } else {
        this
    }
}