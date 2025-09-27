package com.tchat.mobile.models

import kotlinx.serialization.Serializable

/**
 * T026: ChatSession model
 *
 * Chat session management with participant handling, message previews,
 * and real-time state synchronization. Aligned with messaging contract tests.
 */
@Serializable
data class ChatSession(
    val id: String,
    val name: String? = null, // null for direct messages, set for groups
    val type: ChatType,
    val participants: List<ChatParticipant>,
    val lastMessage: MessagePreview? = null,
    val unreadCount: Int = 0,
    val isPinned: Boolean = false,
    val isMuted: Boolean = false,
    val isArchived: Boolean = false,
    val isBlocked: Boolean = false,
    val metadata: ChatMetadata,
    val permissions: ChatPermissions = ChatPermissions(),
    val createdAt: String,
    val updatedAt: String,
    val lastActivityAt: String? = null
)

enum class ChatType {
    DIRECT, // 1-on-1 conversation
    GROUP, // Multi-participant group chat
    CHANNEL, // Broadcast channel with admins and subscribers
    SYSTEM, // System/bot conversations
    SUPPORT // Customer support chat
}

@Serializable
data class ChatParticipant(
    val id: String,
    val name: String,
    val avatar: String? = null,
    val role: ChatRole = ChatRole.MEMBER,
    val status: ParticipantStatus = ParticipantStatus.OFFLINE,
    val lastSeen: String? = null,
    val joinedAt: String,
    val permissions: ParticipantPermissions = ParticipantPermissions(),
    val customTitle: String? = null,
    val isBot: Boolean = false,
    val metadata: Map<String, String> = emptyMap()
)

enum class ChatRole {
    OWNER, // Can delete chat, manage all participants
    ADMIN, // Can manage participants, edit chat settings
    MODERATOR, // Can moderate messages, kick members
    MEMBER, // Regular participant
    GUEST, // Limited participant (in channels)
    BOT // Automated participant
}

enum class ParticipantStatus {
    ONLINE,
    OFFLINE,
    AWAY,
    BUSY,
    INVISIBLE,
    TYPING
}

@Serializable
data class MessagePreview(
    val id: String,
    val content: String, // Truncated/preview content
    val senderId: String,
    val senderName: String,
    val timestamp: String,
    val type: ChatMessageType = ChatMessageType.TEXT,
    val isEdited: Boolean = false,
    val reactions: List<String> = emptyList(), // Emoji reactions
    val attachmentCount: Int = 0,
    val replyToId: String? = null
)

enum class ChatMessageType {
    TEXT,
    IMAGE,
    VIDEO,
    AUDIO,
    FILE,
    LOCATION,
    CONTACT,
    STICKER,
    GIF,
    POLL,
    EVENT,
    SYSTEM,
    DELETED,
    // Extended message types (Phase D: T036-T041)
    EMBED,           // T036: Rich embeds with media previews and link cards
    EVENT_MESSAGE,   // T037: Advanced calendar events with RSVP and reminders
    FORM,            // T038: Interactive forms, surveys, and data collection
    LOCATION_MESSAGE,// T039: Enhanced location with maps, places, check-ins
    PAYMENT,         // T040: Transactions, invoices, and receipts
    FILE_MESSAGE     // T041: Advanced document sharing with version control
}

@Serializable
data class ChatMetadata(
    val description: String? = null,
    val avatar: String? = null,
    val backgroundColor: String? = null,
    val theme: String? = null,
    val language: String = "en",
    val timezone: String = "UTC",
    val tags: List<String> = emptyList(),
    val category: String? = null,
    val isPublic: Boolean = false,
    val inviteLink: String? = null,
    val maxParticipants: Int? = null,
    val autoDeleteMessages: Boolean = false,
    val autoDeleteDays: Int? = null,
    val encryptionEnabled: Boolean = false,
    val backupEnabled: Boolean = true,
    val customFields: Map<String, String> = emptyMap()
)

@Serializable
data class ChatPermissions(
    val canAddMembers: Boolean = true,
    val canRemoveMembers: Boolean = false,
    val canEditInfo: Boolean = false,
    val canPinMessages: Boolean = false,
    val canDeleteMessages: Boolean = false,
    val canInviteByLink: Boolean = true,
    val canSendMessages: Boolean = true,
    val canSendMedia: Boolean = true,
    val canSendPolls: Boolean = true,
    val canSendFiles: Boolean = true,
    val requireApproval: Boolean = false
)

@Serializable
data class ParticipantPermissions(
    val canSendMessages: Boolean = true,
    val canSendMedia: Boolean = true,
    val canAddMembers: Boolean = false,
    val canPinMessages: Boolean = false,
    val canDeleteOwnMessages: Boolean = true,
    val canDeleteOtherMessages: Boolean = false,
    val canEditChatInfo: Boolean = false,
    val isMuted: Boolean = false,
    val mutedUntil: String? = null
)

/**
 * Chat session state for UI management
 */
@Serializable
data class ChatSessionState(
    val sessionId: String,
    val isLoading: Boolean = false,
    val hasMoreMessages: Boolean = true,
    val isTyping: List<String> = emptyList(), // User IDs currently typing
    val isDraftSaving: Boolean = false,
    val draftMessage: String = "",
    val replyToMessage: String? = null,
    val editingMessage: String? = null,
    val selectedMessages: List<String> = emptyList(),
    val searchQuery: String? = null,
    val scrollToMessage: String? = null,
    val networkStatus: String = "connected",
    val lastSync: String? = null,
    val pendingMessages: List<String> = emptyList(),
    val failedMessages: List<String> = emptyList()
)

/**
 * Typing indicator data
 */
@Serializable
data class TypingIndicator(
    val userId: String,
    val userName: String,
    val startedAt: String,
    val expiresAt: String
)

/**
 * Chat session utilities and extensions
 */
fun ChatSession.isDirectMessage(): Boolean = type == ChatType.DIRECT

fun ChatSession.isGroup(): Boolean = type == ChatType.GROUP || type == ChatType.CHANNEL

fun ChatSession.getOtherParticipant(currentUserId: String): ChatParticipant? {
    return if (isDirectMessage()) {
        participants.find { it.id != currentUserId }
    } else null
}

fun ChatSession.getDisplayName(currentUserId: String): String {
    return when {
        name?.isNotBlank() == true -> name
        isDirectMessage() -> {
            getOtherParticipant(currentUserId)?.name ?: "Unknown User"
        }
        participants.size <= 3 -> {
            participants.filter { it.id != currentUserId }
                .joinToString(", ") { it.name }
        }
        else -> {
            val others = participants.filter { it.id != currentUserId }
            "${others.take(2).joinToString(", ") { it.name }} and ${others.size - 2} others"
        }
    }
}

fun ChatSession.getDisplayAvatar(): String? {
    return metadata.avatar ?: if (isDirectMessage()) {
        participants.firstOrNull()?.avatar
    } else null
}

fun ChatSession.canUserPerformAction(userId: String, action: ChatAction): Boolean {
    val participant = participants.find { it.id == userId } ?: return false
    val role = participant.role

    return when (action) {
        ChatAction.SEND_MESSAGE -> permissions.canSendMessages && participant.permissions.canSendMessages
        ChatAction.SEND_MEDIA -> permissions.canSendMedia && participant.permissions.canSendMedia
        ChatAction.ADD_MEMBERS -> {
            permissions.canAddMembers &&
            (role == ChatRole.OWNER || role == ChatRole.ADMIN || participant.permissions.canAddMembers)
        }
        ChatAction.REMOVE_MEMBERS -> {
            permissions.canRemoveMembers &&
            (role == ChatRole.OWNER || role == ChatRole.ADMIN)
        }
        ChatAction.EDIT_INFO -> {
            permissions.canEditInfo &&
            (role == ChatRole.OWNER || role == ChatRole.ADMIN || participant.permissions.canEditChatInfo)
        }
        ChatAction.PIN_MESSAGES -> {
            permissions.canPinMessages &&
            (role == ChatRole.OWNER || role == ChatRole.ADMIN || role == ChatRole.MODERATOR)
        }
        ChatAction.DELETE_MESSAGES -> {
            permissions.canDeleteMessages ||
            (role == ChatRole.OWNER || role == ChatRole.ADMIN || role == ChatRole.MODERATOR)
        }
    }
}

enum class ChatAction {
    SEND_MESSAGE,
    SEND_MEDIA,
    ADD_MEMBERS,
    REMOVE_MEMBERS,
    EDIT_INFO,
    PIN_MESSAGES,
    DELETE_MESSAGES
}

fun ChatSession.getActiveParticipants(): List<ChatParticipant> {
    return participants.filter {
        it.status == ParticipantStatus.ONLINE || it.status == ParticipantStatus.AWAY
    }
}

fun ChatSession.getOnlineCount(): Int {
    return participants.count { it.status == ParticipantStatus.ONLINE }
}

fun ChatSession.hasUnreadMessages(): Boolean = unreadCount > 0

fun ChatSession.isActive(): Boolean {
    return !isArchived && !isBlocked && lastActivityAt != null
}

fun ChatSession.needsAttention(): Boolean {
    return hasUnreadMessages() && !isMuted && !isArchived
}

/**
 * Chat session filtering and sorting
 */
fun List<ChatSession>.filterActive(): List<ChatSession> {
    return filter { !it.isArchived && !it.isBlocked }
}

fun List<ChatSession>.filterUnread(): List<ChatSession> {
    return filter { it.hasUnreadMessages() }
}

fun List<ChatSession>.sortByActivity(): List<ChatSession> {
    return sortedWith(compareByDescending<ChatSession> { it.isPinned }
        .thenByDescending { it.lastMessage?.timestamp ?: it.updatedAt }
        .thenBy { it.getDisplayName("") })
}

fun List<ChatSession>.sortByUnread(): List<ChatSession> {
    return sortedWith(compareByDescending<ChatSession> { it.unreadCount }
        .thenByDescending { it.lastMessage?.timestamp ?: it.updatedAt })
}

/**
 * Search functionality
 */
fun List<ChatSession>.searchByName(query: String, currentUserId: String): List<ChatSession> {
    if (query.isBlank()) return this

    val lowerQuery = query.lowercase()
    return filter { session ->
        session.getDisplayName(currentUserId).lowercase().contains(lowerQuery) ||
        session.participants.any { it.name.lowercase().contains(lowerQuery) } ||
        session.lastMessage?.content?.lowercase()?.contains(lowerQuery) == true
    }
}

/**
 * Participant management utilities
 */
fun ChatParticipant.isOnline(): Boolean {
    return status == ParticipantStatus.ONLINE
}

fun ChatParticipant.canModerate(): Boolean {
    return role == ChatRole.OWNER || role == ChatRole.ADMIN || role == ChatRole.MODERATOR
}

fun ChatParticipant.getDisplayStatus(): String {
    return when (status) {
        ParticipantStatus.ONLINE -> "Online"
        ParticipantStatus.OFFLINE -> lastSeen?.let { "Last seen $it" } ?: "Offline"
        ParticipantStatus.AWAY -> "Away"
        ParticipantStatus.BUSY -> "Busy"
        ParticipantStatus.INVISIBLE -> "Offline" // Show as offline for privacy
        ParticipantStatus.TYPING -> "Typing..."
    }
}