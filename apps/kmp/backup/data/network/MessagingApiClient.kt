package com.tchat.mobile.data.network

import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.request.*
import io.ktor.http.*
import com.tchat.mobile.data.models.Message
import com.tchat.mobile.data.models.Dialog
import com.tchat.mobile.data.models.MessageType
import com.tchat.mobile.data.models.DialogType
import kotlinx.serialization.Serializable
import io.github.aakira.napier.Napier

/**
 * Messaging API client for chat and messaging functionality
 */
class MessagingApiClient(
    private val httpClient: HttpClient,
    private val baseUrl: String
) {

    // Dialog Management

    /**
     * Get user's dialogs with pagination
     */
    suspend fun getDialogs(
        token: String,
        page: Int = 1,
        limit: Int = 20,
        type: String? = null,
        archived: Boolean = false
    ): ApiResult<PaginatedResponse<Dialog>> {
        return try {
            val response: PaginatedResponse<Dialog> = httpClient.get("$baseUrl/dialogs") {
                header("Authorization", "Bearer $token")
                parameter("page", page)
                parameter("limit", limit)
                type?.let { parameter("type", it) }
                parameter("archived", archived)
            }.body()

            if (response.success) {
                ApiResult.Success(response)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to get dialogs"))
            }
        } catch (e: Exception) {
            Napier.e("Get dialogs error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Get specific dialog by ID
     */
    suspend fun getDialog(token: String, dialogId: String): ApiResult<Dialog> {
        return try {
            val response: ApiResponse<Dialog> = httpClient.get("$baseUrl/dialogs/$dialogId") {
                header("Authorization", "Bearer $token")
            }.body()

            if (response.success && response.data != null) {
                ApiResult.Success(response.data!!)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to get dialog"))
            }
        } catch (e: Exception) {
            Napier.e("Get dialog error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Create new dialog (group, channel, etc.)
     */
    suspend fun createDialog(token: String, request: CreateDialogRequest): ApiResult<Dialog> {
        return try {
            val response: ApiResponse<Dialog> = httpClient.post("$baseUrl/dialogs") {
                header("Authorization", "Bearer $token")
                contentType(ContentType.Application.Json)
                setBody(request)
            }.body()

            if (response.success && response.data != null) {
                ApiResult.Success(response.data!!)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to create dialog"))
            }
        } catch (e: Exception) {
            Napier.e("Create dialog error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Update dialog settings
     */
    suspend fun updateDialog(
        token: String,
        dialogId: String,
        request: UpdateDialogRequest
    ): ApiResult<Dialog> {
        return try {
            val response: ApiResponse<Dialog> = httpClient.put("$baseUrl/dialogs/$dialogId") {
                header("Authorization", "Bearer $token")
                contentType(ContentType.Application.Json)
                setBody(request)
            }.body()

            if (response.success && response.data != null) {
                ApiResult.Success(response.data!!)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to update dialog"))
            }
        } catch (e: Exception) {
            Napier.e("Update dialog error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Delete/leave dialog
     */
    suspend fun deleteDialog(token: String, dialogId: String): ApiResult<Unit> {
        return try {
            val response: ApiResponse<Unit> = httpClient.delete("$baseUrl/dialogs/$dialogId") {
                header("Authorization", "Bearer $token")
            }.body()

            if (response.success) {
                ApiResult.Success(Unit)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to delete dialog"))
            }
        } catch (e: Exception) {
            Napier.e("Delete dialog error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    // Dialog Participants Management

    /**
     * Add participants to dialog
     */
    suspend fun addParticipants(
        token: String,
        dialogId: String,
        request: ManageParticipantsRequest
    ): ApiResult<Unit> {
        return try {
            val response: ApiResponse<Unit> = httpClient.post("$baseUrl/dialogs/$dialogId/participants") {
                header("Authorization", "Bearer $token")
                contentType(ContentType.Application.Json)
                setBody(request)
            }.body()

            if (response.success) {
                ApiResult.Success(Unit)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to add participants"))
            }
        } catch (e: Exception) {
            Napier.e("Add participants error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Remove participants from dialog
     */
    suspend fun removeParticipants(
        token: String,
        dialogId: String,
        request: ManageParticipantsRequest
    ): ApiResult<Unit> {
        return try {
            val response: ApiResponse<Unit> = httpClient.delete("$baseUrl/dialogs/$dialogId/participants") {
                header("Authorization", "Bearer $token")
                contentType(ContentType.Application.Json)
                setBody(request)
            }.body()

            if (response.success) {
                ApiResult.Success(Unit)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to remove participants"))
            }
        } catch (e: Exception) {
            Napier.e("Remove participants error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    // Message Management

    /**
     * Get messages from a dialog with pagination
     */
    suspend fun getMessages(
        token: String,
        dialogId: String,
        page: Int = 1,
        limit: Int = 50,
        beforeMessageId: String? = null,
        afterMessageId: String? = null
    ): ApiResult<PaginatedResponse<Message>> {
        return try {
            val response: PaginatedResponse<Message> = httpClient.get("$baseUrl/dialogs/$dialogId/messages") {
                header("Authorization", "Bearer $token")
                parameter("page", page)
                parameter("limit", limit)
                beforeMessageId?.let { parameter("before", it) }
                afterMessageId?.let { parameter("after", it) }
            }.body()

            if (response.success) {
                ApiResult.Success(response)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to get messages"))
            }
        } catch (e: Exception) {
            Napier.e("Get messages error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Send a new message
     */
    suspend fun sendMessage(
        token: String,
        dialogId: String,
        request: SendMessageRequest
    ): ApiResult<Message> {
        return try {
            val response: ApiResponse<Message> = httpClient.post("$baseUrl/dialogs/$dialogId/messages") {
                header("Authorization", "Bearer $token")
                contentType(ContentType.Application.Json)
                setBody(request)
            }.body()

            if (response.success && response.data != null) {
                ApiResult.Success(response.data!!)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to send message"))
            }
        } catch (e: Exception) {
            Napier.e("Send message error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Edit an existing message
     */
    suspend fun editMessage(
        token: String,
        dialogId: String,
        messageId: String,
        request: EditMessageRequest
    ): ApiResult<Message> {
        return try {
            val response: ApiResponse<Message> = httpClient.put("$baseUrl/dialogs/$dialogId/messages/$messageId") {
                header("Authorization", "Bearer $token")
                contentType(ContentType.Application.Json)
                setBody(request)
            }.body()

            if (response.success && response.data != null) {
                ApiResult.Success(response.data!!)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to edit message"))
            }
        } catch (e: Exception) {
            Napier.e("Edit message error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Delete a message
     */
    suspend fun deleteMessage(
        token: String,
        dialogId: String,
        messageId: String,
        forEveryone: Boolean = false
    ): ApiResult<Unit> {
        return try {
            val response: ApiResponse<Unit> = httpClient.delete("$baseUrl/dialogs/$dialogId/messages/$messageId") {
                header("Authorization", "Bearer $token")
                parameter("for_everyone", forEveryone)
            }.body()

            if (response.success) {
                ApiResult.Success(Unit)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to delete message"))
            }
        } catch (e: Exception) {
            Napier.e("Delete message error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Add reaction to message
     */
    suspend fun addReaction(
        token: String,
        dialogId: String,
        messageId: String,
        emoji: String
    ): ApiResult<Unit> {
        return try {
            val response: ApiResponse<Unit> = httpClient.post("$baseUrl/dialogs/$dialogId/messages/$messageId/reactions") {
                header("Authorization", "Bearer $token")
                contentType(ContentType.Application.Json)
                setBody(mapOf("emoji" to emoji))
            }.body()

            if (response.success) {
                ApiResult.Success(Unit)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to add reaction"))
            }
        } catch (e: Exception) {
            Napier.e("Add reaction error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Remove reaction from message
     */
    suspend fun removeReaction(
        token: String,
        dialogId: String,
        messageId: String,
        emoji: String
    ): ApiResult<Unit> {
        return try {
            val response: ApiResponse<Unit> = httpClient.delete("$baseUrl/dialogs/$dialogId/messages/$messageId/reactions/$emoji") {
                header("Authorization", "Bearer $token")
            }.body()

            if (response.success) {
                ApiResult.Success(Unit)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to remove reaction"))
            }
        } catch (e: Exception) {
            Napier.e("Remove reaction error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Mark messages as read
     */
    suspend fun markMessagesAsRead(
        token: String,
        dialogId: String,
        lastMessageId: String
    ): ApiResult<Unit> {
        return try {
            val response: ApiResponse<Unit> = httpClient.post("$baseUrl/dialogs/$dialogId/read") {
                header("Authorization", "Bearer $token")
                contentType(ContentType.Application.Json)
                setBody(mapOf("last_message_id" to lastMessageId))
            }.body()

            if (response.success) {
                ApiResult.Success(Unit)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to mark as read"))
            }
        } catch (e: Exception) {
            Napier.e("Mark as read error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Search messages in dialog
     */
    suspend fun searchMessages(
        token: String,
        dialogId: String,
        query: String,
        page: Int = 1,
        limit: Int = 20
    ): ApiResult<PaginatedResponse<Message>> {
        return try {
            val response: PaginatedResponse<Message> = httpClient.get("$baseUrl/dialogs/$dialogId/messages/search") {
                header("Authorization", "Bearer $token")
                parameter("q", query)
                parameter("page", page)
                parameter("limit", limit)
            }.body()

            if (response.success) {
                ApiResult.Success(response)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to search messages"))
            }
        } catch (e: Exception) {
            Napier.e("Search messages error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Upload file for messaging
     */
    suspend fun uploadFile(
        token: String,
        fileData: ByteArray,
        fileName: String,
        mimeType: String
    ): ApiResult<FileUploadResponse> {
        return try {
            val response: ApiResponse<FileUploadResponse> = httpClient.post("$baseUrl/upload") {
                header("Authorization", "Bearer $token")
                contentType(ContentType.MultiPart.FormData)
                // TODO: Implement multipart form data upload
                // This would require additional Ktor multipart dependencies
            }.body()

            if (response.success && response.data != null) {
                ApiResult.Success(response.data!!)
            } else {
                ApiResult.Error(ApiError.ClientError(400, response.error ?: "Failed to upload file"))
            }
        } catch (e: Exception) {
            Napier.e("Upload file error: ${e.message}", e)
            ApiResult.Error(handleException(e))
        }
    }

    /**
     * Handle common exceptions and convert to ApiError
     */
    private fun handleException(e: Exception): ApiError {
        return when (e) {
            is io.ktor.client.plugins.ClientRequestException -> {
                when (e.response.status.value) {
                    400 -> ApiError.ClientError(400, "Bad request")
                    401 -> ApiError.UnauthorizedError("Unauthorized")
                    403 -> ApiError.ForbiddenError("Forbidden")
                    404 -> ApiError.NotFoundError("Not found")
                    else -> ApiError.ClientError(e.response.status.value, e.message)
                }
            }
            is io.ktor.client.plugins.ServerResponseException -> {
                ApiError.ServerError(e.response.status.value, "Server error: ${e.message}")
            }
            is kotlinx.serialization.SerializationException -> {
                ApiError.SerializationError("Serialization error: ${e.message}")
            }
            else -> {
                ApiError.NetworkError("Network error: ${e.message}")
            }
        }
    }
}

// Request/Response models for Messaging API

@Serializable
data class CreateDialogRequest(
    val type: String, // private, group, channel, broadcast
    val title: String? = null,
    val description: String? = null,
    val participants: List<String> = emptyList(), // User IDs
    val settings: DialogSettingsRequest? = null
)

@Serializable
data class UpdateDialogRequest(
    val title: String? = null,
    val description: String? = null,
    val avatar: String? = null,
    val settings: DialogSettingsRequest? = null
)

@Serializable
data class DialogSettingsRequest(
    val allowInvites: Boolean? = null,
    val allowMembershipRequests: Boolean? = null,
    val requireApprovalForMessages: Boolean? = null,
    val allowMediaSharing: Boolean? = null,
    val allowFileSharing: Boolean? = null,
    val allowVoiceMessages: Boolean? = null,
    val allowPayments: Boolean? = null,
    val allowLocationSharing: Boolean? = null,
    val autoDeleteMessages: Boolean? = null,
    val autoDeleteDuration: Long? = null,
    val slowModeDelay: Int? = null,
    val welcomeMessage: String? = null
)

@Serializable
data class ManageParticipantsRequest(
    val userIds: List<String>
)

@Serializable
data class SendMessageRequest(
    val type: String, // text, voice, file, image, video, payment, location, sticker
    val content: Map<String, String> = emptyMap(),
    val replyToId: String? = null,
    val mentions: List<String> = emptyList()
)

@Serializable
data class EditMessageRequest(
    val content: Map<String, String> = emptyMap()
)

@Serializable
data class FileUploadResponse(
    val fileId: String,
    val url: String,
    val filename: String,
    val mimeType: String,
    val size: Long,
    val thumbnail: String? = null
)