package com.tchat.mobile.data.models

import kotlinx.datetime.Instant
import kotlinx.serialization.Serializable
import com.benasher44.uuid.uuid4

/**
 * Dialog represents a conversation between users
 */
@Serializable
data class Dialog(
    val id: String = uuid4().toString(),
    val type: DialogType = DialogType.PRIVATE,
    val title: String? = null,
    val description: String? = null,
    val avatar: String? = null,
    val participants: List<String> = emptyList(), // User IDs
    val adminIds: List<String> = emptyList(),
    val createdBy: String = "",
    val lastMessage: Message? = null,
    val lastMessageId: String? = null,
    val lastMessageAt: Instant? = null,
    val unreadCount: Int = 0,
    val isArchived: Boolean = false,
    val isMuted: Boolean = false,
    val isPinned: Boolean = false,
    val isDeleted: Boolean = false,
    val settings: DialogSettings = DialogSettings(),
    val metadata: Map<String, String> = emptyMap(),
    val createdAt: Instant = kotlinx.datetime.Clock.System.now(),
    val updatedAt: Instant = kotlinx.datetime.Clock.System.now(),
    val deletedAt: Instant? = null
) {
    /**
     * Check if user is a participant
     */
    fun hasParticipant(userId: String): Boolean = participants.contains(userId)

    /**
     * Check if user is an admin
     */
    fun isAdmin(userId: String): Boolean = adminIds.contains(userId)

    /**
     * Check if user can send messages
     */
    fun canSendMessage(userId: String): Boolean {
        return hasParticipant(userId) && !isDeleted
    }

    /**
     * Check if user can manage dialog (add/remove participants, change settings)
     */
    fun canManageDialog(userId: String): Boolean {
        return when (type) {
            DialogType.PRIVATE -> false // Private chats can't be managed
            DialogType.GROUP -> isAdmin(userId) || createdBy == userId
            DialogType.CHANNEL -> isAdmin(userId) || createdBy == userId
            DialogType.BROADCAST -> createdBy == userId
        }
    }

    /**
     * Add participant to dialog
     */
    fun addParticipant(userId: String, addedBy: String): Dialog? {
        if (!canManageDialog(addedBy) || hasParticipant(userId)) {
            return null
        }

        return copy(
            participants = participants + userId,
            updatedAt = kotlinx.datetime.Clock.System.now()
        )
    }

    /**
     * Remove participant from dialog
     */
    fun removeParticipant(userId: String, removedBy: String): Dialog? {
        if (!canManageDialog(removedBy) || !hasParticipant(userId)) {
            return null
        }

        return copy(
            participants = participants - userId,
            adminIds = adminIds - userId, // Remove from admins if present
            updatedAt = kotlinx.datetime.Clock.System.now()
        )
    }

    /**
     * Promote user to admin
     */
    fun promoteToAdmin(userId: String, promotedBy: String): Dialog? {
        if (!canManageDialog(promotedBy) || !hasParticipant(userId) || isAdmin(userId)) {
            return null
        }

        return copy(
            adminIds = adminIds + userId,
            updatedAt = kotlinx.datetime.Clock.System.now()
        )
    }

    /**
     * Demote admin to regular participant
     */
    fun demoteAdmin(userId: String, demotedBy: String): Dialog? {
        if (!canManageDialog(demotedBy) || !isAdmin(userId) || userId == createdBy) {
            return null
        }

        return copy(
            adminIds = adminIds - userId,
            updatedAt = kotlinx.datetime.Clock.System.now()
        )
    }

    /**
     * Update last message
     */
    fun updateLastMessage(message: Message): Dialog {
        return copy(
            lastMessage = message,
            lastMessageId = message.id,
            lastMessageAt = message.createdAt,
            updatedAt = kotlinx.datetime.Clock.System.now()
        )
    }

    /**
     * Mark messages as read for user
     */
    fun markAsRead(userId: String): Dialog {
        return if (hasParticipant(userId)) {
            copy(unreadCount = 0)
        } else {
            this
        }
    }

    /**
     * Increment unread count
     */
    fun incrementUnreadCount(): Dialog = copy(unreadCount = unreadCount + 1)

    /**
     * Archive dialog
     */
    fun archive(): Dialog = copy(
        isArchived = true,
        updatedAt = kotlinx.datetime.Clock.System.now()
    )

    /**
     * Unarchive dialog
     */
    fun unarchive(): Dialog = copy(
        isArchived = false,
        updatedAt = kotlinx.datetime.Clock.System.now()
    )

    /**
     * Mute dialog
     */
    fun mute(): Dialog = copy(
        isMuted = true,
        updatedAt = kotlinx.datetime.Clock.System.now()
    )

    /**
     * Unmute dialog
     */
    fun unmute(): Dialog = copy(
        isMuted = false,
        updatedAt = kotlinx.datetime.Clock.System.now()
    )

    /**
     * Pin dialog
     */
    fun pin(): Dialog = copy(
        isPinned = true,
        updatedAt = kotlinx.datetime.Clock.System.now()
    )

    /**
     * Unpin dialog
     */
    fun unpin(): Dialog = copy(
        isPinned = false,
        updatedAt = kotlinx.datetime.Clock.System.now()
    )

    /**
     * Soft delete dialog
     */
    fun softDelete(): Dialog = copy(
        isDeleted = true,
        deletedAt = kotlinx.datetime.Clock.System.now(),
        updatedAt = kotlinx.datetime.Clock.System.now()
    )

    /**
     * Get display name for dialog
     */
    fun getDisplayName(currentUserId: String, users: List<User> = emptyList()): String {
        return when (type) {
            DialogType.PRIVATE -> {
                // For private chats, show the other participant's name
                val otherParticipant = participants.find { it != currentUserId }
                val otherUser = users.find { it.id == otherParticipant }
                otherUser?.getDisplayName() ?: "Unknown User"
            }
            DialogType.GROUP, DialogType.CHANNEL, DialogType.BROADCAST -> {
                title ?: "Unnamed ${type.name.lowercase().replaceFirstChar { it.uppercase() }}"
            }
        }
    }

    /**
     * Get participant count
     */
    fun getParticipantCount(): Int = participants.size

    /**
     * Convert to public dialog for API responses
     */
    fun toPublicDialog(): Map<String, Any?> = mapOf(
        "id" to id,
        "type" to type.name.lowercase(),
        "title" to title,
        "description" to description,
        "avatar" to avatar,
        "participant_count" to participants.size,
        "last_message_at" to lastMessageAt?.toString(),
        "unread_count" to unreadCount,
        "is_archived" to isArchived,
        "is_muted" to isMuted,
        "is_pinned" to isPinned,
        "created_at" to createdAt.toString(),
        "updated_at" to updatedAt.toString()
    )
}

/**
 * Dialog Type represents the type of conversation
 */
@Serializable
enum class DialogType {
    PRIVATE,    // One-on-one conversation
    GROUP,      // Group conversation with multiple participants
    CHANNEL,    // Channel where only admins can post
    BROADCAST;  // Broadcast channel for announcements

    companion object {
        fun fromString(type: String): DialogType? {
            return values().find { it.name.lowercase() == type.lowercase() }
        }
    }
}

/**
 * Dialog Settings for customization and permissions
 */
@Serializable
data class DialogSettings(
    val allowInvites: Boolean = true,
    val allowMembershipRequests: Boolean = true,
    val requireApprovalForMessages: Boolean = false,
    val allowMediaSharing: Boolean = true,
    val allowFileSharing: Boolean = true,
    val allowVoiceMessages: Boolean = true,
    val allowPayments: Boolean = true,
    val allowLocationSharing: Boolean = true,
    val autoDeleteMessages: Boolean = false,
    val autoDeleteDuration: Long = 0, // in milliseconds
    val slowModeDelay: Int = 0, // in seconds, 0 = disabled
    val maxMembers: Int = 200000, // Maximum number of participants
    val welcomeMessage: String? = null,
    val pinnedMessageId: String? = null,
    val customPermissions: Map<String, Boolean> = emptyMap()
) {
    /**
     * Check if a specific permission is enabled
     */
    fun hasPermission(permission: String): Boolean {
        return customPermissions[permission] ?: true
    }
}

/**
 * Dialog Participant represents a participant in a dialog
 */
@Serializable
data class DialogParticipant(
    val userId: String,
    val role: ParticipantRole = ParticipantRole.MEMBER,
    val joinedAt: Instant,
    val invitedBy: String? = null,
    val lastReadMessageId: String? = null,
    val lastReadAt: Instant? = null,
    val isNotificationEnabled: Boolean = true,
    val customTitle: String? = null,
    val permissions: Map<String, Boolean> = emptyMap()
) {
    /**
     * Check if participant has specific permission
     */
    fun hasPermission(permission: String): Boolean {
        return permissions[permission] ?: role.getDefaultPermissions()[permission] ?: false
    }
}

/**
 * Participant Role in a dialog
 */
@Serializable
enum class ParticipantRole {
    OWNER,      // Dialog creator with all permissions
    ADMIN,      // Administrator with management permissions
    MODERATOR,  // Moderator with limited management permissions
    MEMBER,     // Regular participant
    RESTRICTED; // Restricted participant with limited permissions

    /**
     * Get default permissions for this role
     */
    fun getDefaultPermissions(): Map<String, Boolean> {
        return when (this) {
            OWNER -> mapOf(
                "send_messages" to true,
                "send_media" to true,
                "add_members" to true,
                "remove_members" to true,
                "change_info" to true,
                "pin_messages" to true,
                "delete_messages" to true,
                "ban_users" to true,
                "promote_members" to true
            )
            ADMIN -> mapOf(
                "send_messages" to true,
                "send_media" to true,
                "add_members" to true,
                "remove_members" to true,
                "pin_messages" to true,
                "delete_messages" to true,
                "ban_users" to true,
                "promote_members" to false
            )
            MODERATOR -> mapOf(
                "send_messages" to true,
                "send_media" to true,
                "add_members" to false,
                "remove_members" to true,
                "pin_messages" to true,
                "delete_messages" to true,
                "ban_users" to true,
                "promote_members" to false
            )
            MEMBER -> mapOf(
                "send_messages" to true,
                "send_media" to true,
                "add_members" to false,
                "remove_members" to false,
                "pin_messages" to false,
                "delete_messages" to false,
                "ban_users" to false,
                "promote_members" to false
            )
            RESTRICTED -> mapOf(
                "send_messages" to false,
                "send_media" to false,
                "add_members" to false,
                "remove_members" to false,
                "pin_messages" to false,
                "delete_messages" to false,
                "ban_users" to false,
                "promote_members" to false
            )
        }
    }
}

/**
 * Dialog validation utility
 */
object DialogValidator {
    /**
     * Validate dialog based on its type and settings
     */
    fun validateDialog(dialog: Dialog): List<String> {
        val errors = mutableListOf<String>()

        // Basic validation
        if (dialog.createdBy.isBlank()) errors.add("Created by user ID is required")

        // Type-specific validation
        when (dialog.type) {
            DialogType.PRIVATE -> {
                if (dialog.participants.size != 2) {
                    errors.add("Private dialog must have exactly 2 participants")
                }
                if (dialog.title != null) {
                    errors.add("Private dialog cannot have a title")
                }
            }
            DialogType.GROUP -> {
                if (dialog.participants.size < 3) {
                    errors.add("Group dialog must have at least 3 participants")
                }
                if (dialog.participants.size > dialog.settings.maxMembers) {
                    errors.add("Group exceeds maximum member limit")
                }
            }
            DialogType.CHANNEL -> {
                if (dialog.title.isNullOrBlank()) {
                    errors.add("Channel must have a title")
                }
            }
            DialogType.BROADCAST -> {
                if (dialog.title.isNullOrBlank()) {
                    errors.add("Broadcast channel must have a title")
                }
                if (dialog.adminIds.isNotEmpty() && !dialog.adminIds.contains(dialog.createdBy)) {
                    errors.add("Broadcast channel creator must be an admin")
                }
            }
        }

        // Participants validation
        if (!dialog.participants.contains(dialog.createdBy)) {
            errors.add("Dialog creator must be a participant")
        }

        // Admin validation
        dialog.adminIds.forEach { adminId ->
            if (!dialog.participants.contains(adminId)) {
                errors.add("Admin $adminId must be a participant")
            }
        }

        return errors
    }

    /**
     * Validate dialog settings
     */
    fun validateDialogSettings(settings: DialogSettings): List<String> {
        val errors = mutableListOf<String>()

        if (settings.maxMembers <= 0) {
            errors.add("Maximum members must be positive")
        }

        if (settings.slowModeDelay < 0) {
            errors.add("Slow mode delay cannot be negative")
        }

        if (settings.autoDeleteMessages && settings.autoDeleteDuration <= 0) {
            errors.add("Auto delete duration must be positive when auto delete is enabled")
        }

        return errors
    }
}