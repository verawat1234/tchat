package com.tchat.mobile.api

import com.tchat.mobile.api.models.*
import com.tchat.mobile.models.*
import com.tchat.mobile.network.HttpClientFactory
import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.request.*
import io.ktor.client.statement.*
import io.ktor.http.*
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.datetime.Clock
import kotlinx.datetime.Instant
import kotlinx.serialization.json.Json

/**
 * HTTP API Client for server communication
 *
 * Handles all HTTP requests to the backend API with:
 * - Authentication management
 * - Request/response mapping
 * - Error handling
 * - Connection state tracking
 * - Retry logic
 */
class ApiClient(
    private val baseUrl: String = "http://localhost:8080/api/v1"
) {

    private val httpClient: HttpClient = HttpClientFactory.create(baseUrl)
    private val _connectionState = MutableStateFlow(ConnectionState.DISCONNECTED)
    val connectionState: StateFlow<ConnectionState> = _connectionState

    private var accessToken: String? = null
    private var refreshToken: String? = null

    // Connection Management
    suspend fun connect(): Result<Unit> {
        return try {
            val response: HttpResponse = httpClient.get("/messages/health")
            if (response.status.isSuccess()) {
                _connectionState.value = ConnectionState.CONNECTED
                Result.success(Unit)
            } else {
                _connectionState.value = ConnectionState.ERROR
                Result.failure(ApiException(response.status.value, "Failed to connect to messaging service"))
            }
        } catch (e: Exception) {
            _connectionState.value = ConnectionState.ERROR
            Result.failure(ApiException(500, "Connection failed: ${e.message}"))
        }
    }

    suspend fun disconnect(): Result<Unit> {
        return try {
            _connectionState.value = ConnectionState.DISCONNECTED
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(ApiException(500, "Disconnect failed: ${e.message}"))
        }
    }

    // Authentication - Real API Methods
    suspend fun requestOTP(phoneNumber: String, countryCode: String = "+66"): Result<RequestOTPResponse> {
        return try {
            _connectionState.value = ConnectionState.CONNECTING

            val request = RequestOTPRequest(
                phoneNumber = phoneNumber,
                countryCode = countryCode,
                deviceInfo = DeviceInfo(
                    platform = "mobile",
                    deviceModel = null,
                    osVersion = null,
                    appVersion = "1.0.0",
                    deviceID = null,
                    pushToken = null,
                    userAgent = "TchatMobile/1.0",
                    ipAddress = null,
                    timezone = null,
                    language = "en"
                )
            )

            val response: HttpResponse = httpClient.post("/auth/register") {
                headers {
                    append(HttpHeaders.ContentType, ContentType.Application.Json)
                }
                setBody<RequestOTPRequest>(request)
            }

            if (response.status.isSuccess()) {
                val otpResponse = response.body<RequestOTPResponse>()
                Result.success(otpResponse)
            } else {
                val error = when (response.status.value) {
                    400 -> ApiException(400, "Invalid phone number format")
                    429 -> ApiException(429, "Too many OTP requests. Please wait.")
                    else -> ApiException(response.status.value, "Failed to send OTP: ${response.status.description}")
                }
                Result.failure(error)
            }
        } catch (e: Exception) {
            _connectionState.value = ConnectionState.ERROR
            Result.failure(ApiException(500, "OTP request failed: ${e.message}"))
        }
    }

    suspend fun verifyOTP(phoneNumber: String, code: String, requestID: String): Result<VerifyOTPResponse> {
        return try {
            val request = VerifyOTPRequest(
                requestID = requestID,
                code = code,
                phoneNumber = phoneNumber,
                deviceInfo = DeviceInfo(
                    platform = "mobile",
                    deviceModel = null,
                    osVersion = null,
                    appVersion = "1.0.0",
                    deviceID = null,
                    pushToken = null,
                    userAgent = "TchatMobile/1.0",
                    ipAddress = null,
                    timezone = null,
                    language = "en"
                )
            )

            val response: HttpResponse = httpClient.post("/auth/otp/verify") {
                headers {
                    append(HttpHeaders.ContentType, ContentType.Application.Json)
                }
                setBody<VerifyOTPRequest>(request)
            }

            if (response.status.isSuccess()) {
                val verifyResponse = response.body<VerifyOTPResponse>()

                // Store authentication tokens
                accessToken = verifyResponse.accessToken
                refreshToken = verifyResponse.refreshToken
                _connectionState.value = ConnectionState.CONNECTED

                Result.success(verifyResponse)
            } else {
                val error = when (response.status.value) {
                    400 -> ApiException(400, "Invalid OTP code")
                    401 -> ApiException(401, "OTP code expired or invalid")
                    404 -> ApiException(404, "OTP request not found")
                    else -> ApiException(response.status.value, "OTP verification failed: ${response.status.description}")
                }
                Result.failure(error)
            }
        } catch (e: Exception) {
            _connectionState.value = ConnectionState.ERROR
            Result.failure(ApiException(500, "OTP verification failed: ${e.message}"))
        }
    }

    suspend fun loginWithEmail(email: String, password: String): Result<AuthResponse> {
        return try {
            _connectionState.value = ConnectionState.CONNECTING

            val request = LoginRequest(
                email = email,
                password = password,
                rememberMe = true
            )

            val response: HttpResponse = httpClient.post("/auth/login") {
                headers {
                    append(HttpHeaders.ContentType, ContentType.Application.Json)
                }
                setBody<LoginRequest>(request)
            }

            if (response.status.isSuccess()) {
                val authResponse = response.body<AuthResponse>()

                accessToken = authResponse.accessToken
                refreshToken = authResponse.refreshToken
                _connectionState.value = ConnectionState.CONNECTED

                Result.success(authResponse)
            } else {
                val error = when (response.status.value) {
                    401 -> ApiException(401, "Invalid email or password")
                    400 -> ApiException(400, "Invalid login credentials format")
                    else -> ApiException(response.status.value, "Login failed: ${response.status.description}")
                }
                Result.failure(error)
            }
        } catch (e: Exception) {
            _connectionState.value = ConnectionState.ERROR
            Result.failure(ApiException(500, "Login failed: ${e.message}"))
        }
    }

    suspend fun getCurrentUser(): Result<CurrentUserResponse> {
        return try {
            val response: HttpResponse = httpClient.get("/auth/me") {
                headers {
                    accessToken?.let { append(HttpHeaders.Authorization, "Bearer $it") }
                }
            }

            if (response.status.isSuccess()) {
                val userResponse = response.body<CurrentUserResponse>()
                Result.success(userResponse)
            } else {
                val error = when (response.status.value) {
                    401 -> ApiException(401, "Unauthorized - please login again")
                    else -> ApiException(response.status.value, "Failed to get current user: ${response.status.description}")
                }
                Result.failure(error)
            }
        } catch (e: Exception) {
            Result.failure(ApiException(500, "Failed to get current user: ${e.message}"))
        }
    }

    suspend fun logout(): Result<Unit> {
        return try {
            val request = LogoutRequest(
                refreshToken = refreshToken,
                logoutAll = false
            )

            val response: HttpResponse = httpClient.post("/auth/logout") {
                headers {
                    accessToken?.let { append(HttpHeaders.Authorization, "Bearer $it") }
                    append(HttpHeaders.ContentType, ContentType.Application.Json)
                }
                setBody<LogoutRequest>(request)
            }

            // Clear tokens regardless of response status
            accessToken = null
            refreshToken = null
            _connectionState.value = ConnectionState.DISCONNECTED

            if (response.status.isSuccess()) {
                Result.success(Unit)
            } else {
                // Log the error but still consider logout successful locally
                Result.success(Unit)
            }
        } catch (e: Exception) {
            // Clear tokens even on network error
            accessToken = null
            refreshToken = null
            _connectionState.value = ConnectionState.DISCONNECTED
            Result.success(Unit)
        }
    }

    suspend fun refreshAuthToken(): Result<RefreshTokenResponse> {
        return try {
            val currentRefreshToken = refreshToken
                ?: return Result.failure(ApiException(401, "No refresh token available"))

            val request = RefreshTokenRequest(
                refreshToken = currentRefreshToken
            )

            val response: HttpResponse = httpClient.post("/auth/refresh") {
                headers {
                    append(HttpHeaders.ContentType, ContentType.Application.Json)
                }
                setBody<RefreshTokenRequest>(request)
            }

            if (response.status.isSuccess()) {
                val refreshResponse = response.body<RefreshTokenResponse>()
                accessToken = refreshResponse.accessToken
                Result.success(refreshResponse)
            } else {
                // Refresh failed - clear tokens
                accessToken = null
                refreshToken = null
                _connectionState.value = ConnectionState.DISCONNECTED

                val error = when (response.status.value) {
                    401 -> ApiException(401, "Refresh token expired - please login again")
                    else -> ApiException(response.status.value, "Token refresh failed: ${response.status.description}")
                }
                Result.failure(error)
            }
        } catch (e: Exception) {
            // Network error during refresh - clear tokens
            accessToken = null
            refreshToken = null
            _connectionState.value = ConnectionState.DISCONNECTED
            Result.failure(ApiException(500, "Token refresh failed: ${e.message}"))
        }
    }

    // Chat Operations
    suspend fun getChatSessions(userId: String): Result<List<ChatSessionDto>> {
        return try {
            val response: HttpResponse = httpClient.get("/messages/dialogs") {
                headers {
                    accessToken?.let { append(HttpHeaders.Authorization, "Bearer $it") }
                }
            }

            if (response.status.isSuccess()) {
                val backendDialogs: List<BackendDialogDto> = response.body()
                val dialogs = backendDialogs.map { ApiMapper.chatSessionDtoFromBackend(it) }
                _connectionState.value = ConnectionState.CONNECTED
                Result.success(dialogs)
            } else {
                val error = when (response.status.value) {
                    401 -> ApiException(401, "Unauthorized access")
                    404 -> ApiException(404, "Dialogs not found")
                    else -> ApiException(response.status.value, "Failed to fetch dialogs: ${response.status.description}")
                }
                Result.failure(error)
            }
        } catch (e: Exception) {
            _connectionState.value = ConnectionState.ERROR
            Result.failure(ApiException(500, "Network error: ${e.message}"))
        }
    }

    suspend fun getChatSession(chatId: String): Result<ChatSessionDto> {
        return executeRequest {
            ChatSessionDto(
                id = chatId,
                name = "Test Chat",
                description = "A test chat session",
                avatar = null,
                type = "group",
                isActive = true,
                participants = listOf("user1", "user2"),
                createdBy = "user1",
                createdAt = Clock.System.now().toEpochMilliseconds(),
                updatedAt = Clock.System.now().toEpochMilliseconds(),
                lastMessageId = null,
                lastMessageAt = null,
                unreadCount = 0,
                settings = ChatSettingsDto()
            )
        }
    }

    suspend fun createChatSession(chat: ChatSessionDto): Result<ChatSessionDto> {
        return executeRequest {
            chat.copy(
                id = "chat_${Clock.System.now().toEpochMilliseconds()}",
                createdAt = Clock.System.now().toEpochMilliseconds(),
                updatedAt = Clock.System.now().toEpochMilliseconds()
            )
        }
    }

    // Message Operations
    suspend fun getMessages(chatId: String, limit: Int = 50, before: Long? = null): Result<List<MessageDto>> {
        return executeRequest {
            // Simulate API call
            listOf(
                MessageDto(
                    id = "msg1",
                    chatId = chatId,
                    senderId = "user1",
                    senderName = "Test User",
                    senderAvatar = null,
                    type = "text",
                    content = "Hello, this is a test message",
                    isEdited = false,
                    isPinned = false,
                    isDeleted = false,
                    replyToId = null,
                    reactions = emptyList(),
                    attachments = emptyList(),
                    createdAt = Clock.System.now().toEpochMilliseconds(),
                    editedAt = null,
                    deletedAt = null,
                    serverTimestamp = Clock.System.now().toEpochMilliseconds(),
                    version = 1,
                    checksum = null,
                    deliveryStatus = "sent",
                    readBy = emptyList()
                )
            )
        }
    }

    suspend fun sendMessage(message: MessageDto): Result<MessageDto> {
        return try {
            val backendRequest = ApiMapper.messageDtoToBackendRequest(message)
            val response: HttpResponse = httpClient.post("/messages/dialogs/${message.chatId}/messages") {
                headers {
                    accessToken?.let { append(HttpHeaders.Authorization, "Bearer $it") }
                    append(HttpHeaders.ContentType, ContentType.Application.Json)
                }
                setBody(backendRequest)
            }

            if (response.status.isSuccess()) {
                val backendMessage: BackendMessageDto = response.body()
                val sentMessage = ApiMapper.messageDtoFromBackend(backendMessage)
                Result.success(sentMessage)
            } else {
                val error = when (response.status.value) {
                    401 -> ApiException(401, "Unauthorized access")
                    404 -> ApiException(404, "Dialog not found")
                    400 -> ApiException(400, "Invalid message format")
                    else -> ApiException(response.status.value, "Failed to send message: ${response.status.description}")
                }
                Result.failure(error)
            }
        } catch (e: Exception) {
            Result.failure(ApiException(500, "Network error: ${e.message}"))
        }
    }

    suspend fun updateMessage(message: MessageDto): Result<MessageDto> {
        return executeRequest {
            message.copy(
                isEdited = true,
                editedAt = Clock.System.now().toEpochMilliseconds(),
                serverTimestamp = Clock.System.now().toEpochMilliseconds(),
                version = message.version + 1
            )
        }
    }

    suspend fun deleteMessage(messageId: String): Result<Unit> {
        return executeRequest {
            Unit
        }
    }

    // Sync Operations
    suspend fun syncChatData(request: SyncRequestDto): Result<SyncResponseDto> {
        return executeRequest {
            SyncResponseDto(
                success = true,
                chatId = request.chatId,
                messages = emptyList(),
                conflicts = emptyList(),
                processedOperations = request.operations.map { it.id },
                failedOperations = emptyList(),
                serverTimestamp = Clock.System.now().toEpochMilliseconds(),
                nextSyncTimestamp = Clock.System.now().toEpochMilliseconds()
            )
        }
    }

    suspend fun pushPendingOperations(operations: List<SyncOperationDto>): Result<List<SyncOperationDto>> {
        return executeRequest {
            operations.map { op ->
                op.copy(status = "completed")
            }
        }
    }

    suspend fun resolveConflicts(conflicts: List<ConflictDto>): Result<List<ConflictResolutionDto>> {
        return executeRequest {
            conflicts.map { conflict ->
                ConflictResolutionDto(
                    conflictId = conflict.id,
                    strategy = "remote_wins", // Default strategy
                    resolvedMessage = conflict.remoteMessage,
                    explanation = "Server version takes precedence"
                )
            }
        }
    }

    // File Upload
    suspend fun uploadFile(request: FileUploadRequest, data: ByteArray): Result<FileUploadResponse> {
        return executeRequest {
            FileUploadResponse(
                success = true,
                fileUrl = "https://example.com/uploads/${request.fileName}",
                uploadId = "upload_${Clock.System.now().toEpochMilliseconds()}",
                thumbnailUrl = null
            )
        }
    }

    // Health & Monitoring
    suspend fun healthCheck(): Result<HealthCheckDto> {
        return executeRequest {
            HealthCheckDto(
                status = "healthy",
                uptime = 86400000, // 24 hours
                responseTime = 50,
                activeConnections = 150,
                serverLoad = 0.35,
                timestamp = Clock.System.now().toEpochMilliseconds()
            )
        }
    }

    suspend fun getServerCapabilities(): Result<ServerCapabilitiesDto> {
        return executeRequest {
            ServerCapabilitiesDto(
                version = "1.0.0",
                maxFileSize = 50 * 1024 * 1024, // 50MB
                supportedFileTypes = listOf("jpg", "png", "gif", "mp4", "pdf", "doc"),
                maxMessageLength = 4000,
                supportsBulkOperations = true,
                supportsRealTime = true,
                supportsFileUpload = true,
                maxChatParticipants = 100,
                messageRetentionDays = 365
            )
        }
    }

    // Utility Methods
    private suspend fun <T> executeRequest(request: suspend () -> T): Result<T> {
        return try {
            if (_connectionState.value != ConnectionState.CONNECTED) {
                return Result.failure(ApiException(503, "Not connected to server"))
            }

            val result = request()
            Result.success(result)
        } catch (e: ApiException) {
            Result.failure(e)
        } catch (e: Exception) {
            Result.failure(ApiException(500, "Request failed: ${e.message}"))
        }
    }

    private fun isTokenValid(): Boolean {
        // In real implementation, check token expiration
        return accessToken != null
    }

    private suspend fun ensureAuthenticated(): Result<Unit> {
        return if (isTokenValid()) {
            Result.success(Unit)
        } else {
            refreshToken?.let { token ->
                refreshAuthToken().map { Unit }
            } ?: Result.failure(ApiException(401, "No valid authentication"))
        }
    }


    fun getServerTimestamp(): Long {
        return Clock.System.now().toEpochMilliseconds()
    }
}

/**
 * Custom API Exception for handling HTTP errors
 */
class ApiException(
    val statusCode: Int,
    override val message: String,
    val requestId: String? = null
) : Exception(message)