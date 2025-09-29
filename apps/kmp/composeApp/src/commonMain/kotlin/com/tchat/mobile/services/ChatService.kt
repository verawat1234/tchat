package com.tchat.mobile.services

import com.tchat.mobile.models.*
import com.tchat.mobile.repositories.ChatRepository
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.first

/**
 * ChatService - Business logic layer for chat operations
 *
 * Coordinates between repository layer and UI, handling business rules,
 * validation, and complex operations that involve multiple data sources.
 */
class ChatService(
    private val chatRepository: ChatRepository
) {

    // Message operations
    suspend fun searchMessages(chatId: String, query: String): Result<List<Message>> {
        return if (query.trim().isEmpty()) {
            Result.success(emptyList())
        } else {
            chatRepository.searchMessages(chatId, query.trim())
        }
    }

    suspend fun sendMessage(
        chatId: String,
        content: String,
        replyToId: String? = null
    ): Result<Message> {
        return if (content.trim().isEmpty()) {
            Result.failure(IllegalArgumentException("Message content cannot be empty"))
        } else {
            chatRepository.sendTextMessage(chatId, content.trim(), replyToId)
        }
    }

    suspend fun deleteMessage(messageId: String): Result<Boolean> {
        return chatRepository.deleteMessage(messageId)
    }

    // Chat session operations
    suspend fun muteChat(chatId: String, muted: Boolean): Result<Boolean> {
        return if (muted) {
            chatRepository.muteChat(chatId)
        } else {
            chatRepository.unmuteChat(chatId)
        }
    }

    suspend fun pinChat(chatId: String, pinned: Boolean): Result<Boolean> {
        // Use updateChatSession to modify pin status
        return try {
            val session = chatRepository.getChatSession(chatId).getOrThrow()
            val updatedSession = session.copy(isPinned = pinned)
            chatRepository.updateChatSession(chatId, updatedSession)
                .map { true }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun archiveChat(chatId: String, archived: Boolean): Result<Boolean> {
        return if (archived) {
            chatRepository.archiveChat(chatId)
        } else {
            chatRepository.unarchiveChat(chatId)
        }
    }

    suspend fun blockChat(chatId: String, blocked: Boolean): Result<Boolean> {
        // Use updateChatSession to modify block status
        return try {
            val session = chatRepository.getChatSession(chatId).getOrThrow()
            val updatedSession = session.copy(isBlocked = blocked)
            chatRepository.updateChatSession(chatId, updatedSession)
                .map { true }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun deleteChat(chatId: String): Result<Boolean> {
        return chatRepository.deleteChatSession(chatId)
    }

    suspend fun clearChatHistory(chatId: String): Result<Boolean> {
        // Business logic: confirm this is a destructive operation
        return try {
            val messages = chatRepository.getMessages(chatId).getOrNull() ?: emptyList()
            if (messages.isNotEmpty()) {
                // Delete all messages in the chat
                messages.forEach { message ->
                    chatRepository.deleteMessage(message.id)
                }
            }
            Result.success(true)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun exportChat(chatId: String): Result<String> {
        return try {
            val messages = chatRepository.getMessages(chatId).getOrNull() ?: emptyList()
            val chatSession = chatRepository.getChatSession(chatId).getOrNull()

            // Format chat export
            val exportContent = buildString {
                appendLine("Chat Export: ${chatSession?.name ?: "Unknown Chat"}")
                appendLine("Exported on: ${kotlinx.datetime.Clock.System.now()}")
                appendLine("Total messages: ${messages.size}")
                appendLine()

                messages.forEach { message ->
                    appendLine("${message.senderName} (${message.createdAt}): ${message.getDisplayContent()}")
                }
            }

            Result.success(exportContent)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Call operations (temporarily disabled)
    // suspend fun startVideoCall(chatId: String, participants: List<ChatParticipant>): Result<Boolean> {
    //     return Result.failure(IllegalStateException("Call service not available"))
    // }

    // suspend fun startVoiceCall(chatId: String, participants: List<ChatParticipant>): Result<Boolean> {
    //     return Result.failure(IllegalStateException("Call service not available"))
    // }

    // Chat info operations
    suspend fun getChatInfo(chatId: String): Result<ChatInfo> {
        return try {
            val session = chatRepository.getChatSession(chatId).getOrThrow()
            val participants = session.participants
            val messageCount = chatRepository.getMessages(chatId).getOrNull()?.size ?: 0

            val chatInfo = ChatInfo(
                session = session,
                participants = participants,
                messageCount = messageCount
            )

            Result.success(chatInfo)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun reportChat(chatId: String, reason: String): Result<Boolean> {
        return try {
            // Business logic for reporting
            val session = chatRepository.getChatSession(chatId).getOrThrow()

            // Log the report (in real app, this would go to moderation system)
            println("Chat reported: ${session.name}, Reason: $reason")

            Result.success(true)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Observables
    fun observeMessages(chatId: String): Flow<List<Message>> {
        return chatRepository.observeMessages(chatId)
    }

    fun observeTypingIndicators(chatId: String): Flow<List<TypingIndicator>> {
        return chatRepository.observeTypingIndicators(chatId)
    }
}

/**
 * CallService - Handles video and voice call operations (temporarily disabled)
 */
// interface CallService {
//     suspend fun startVideoCall(chatId: String, participants: List<ChatParticipant>): Result<Boolean>
//     suspend fun startVoiceCall(chatId: String, participants: List<ChatParticipant>): Result<Boolean>
// }

/**
 * ChatInfo - Comprehensive chat information
 */
data class ChatInfo(
    val session: ChatSession,
    val participants: List<ChatParticipant>,
    val messageCount: Int
)