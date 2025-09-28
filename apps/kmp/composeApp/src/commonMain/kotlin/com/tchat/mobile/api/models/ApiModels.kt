package com.tchat.mobile.api.models

import kotlinx.datetime.Instant
import kotlinx.serialization.Serializable
import kotlinx.serialization.SerialName
import kotlinx.serialization.Contextual

/**
 * API Data Transfer Objects (DTOs) for server communication
 *
 * These models represent the JSON structure expected by the API
 * and are separate from domain models to allow independent evolution
 */

// Authentication DTOs
@Serializable
data class AuthRequest(
    val token: String,
    val deviceId: String? = null,
    val platform: String = "mobile"
)

@Serializable
data class AuthResponse(
    val success: Boolean,
    val user: UserDto,
    val accessToken: String,
    val refreshToken: String,
    val expiresAt: Long
)

@Serializable
data class RefreshTokenRequest(
    val refreshToken: String
)

@Serializable
data class RefreshTokenResponse(
    val accessToken: String,
    val expiresAt: Long
)

// User DTOs
@Serializable
data class UserDto(
    val id: String,
    val email: String,
    val displayName: String,
    val avatar: String? = null,
    val isActive: Boolean = true,
    val createdAt: Long,
    val updatedAt: Long
)

@Serializable
data class UserProfileDto(
    val userId: String,
    val displayName: String,
    val username: String,
    val avatar: String? = null,
    val bio: String? = null,
    val isVerified: Boolean = false,
    val isOnline: Boolean = false,
    val lastSeen: Long? = null,
    val statusMessage: String? = null
)

// Chat DTOs
@Serializable
data class ChatSessionDto(
    val id: String,
    val name: String,
    val description: String? = null,
    val avatar: String? = null,
    val type: String, // "direct", "group", "channel"
    val isActive: Boolean = true,
    val participants: List<String>,
    val createdBy: String,
    val createdAt: Long,
    val updatedAt: Long,
    val lastMessageId: String? = null,
    val lastMessageAt: Long? = null,
    val unreadCount: Int = 0,
    val settings: ChatSettingsDto? = null
)

@Serializable
data class ChatSettingsDto(
    val isMuted: Boolean = false,
    val muteUntil: Long? = null,
    val notificationsEnabled: Boolean = true,
    val messageRetention: Int = 30, // days
    val allowInvites: Boolean = true,
    val isPrivate: Boolean = false
)

// Message DTOs
@Serializable
data class MessageDto(
    val id: String,
    val chatId: String,
    val senderId: String,
    val senderName: String,
    val senderAvatar: String? = null,
    val type: String, // "text", "image", "video", "audio", "file", "location"
    val content: String,
    val isEdited: Boolean = false,
    val isPinned: Boolean = false,
    val isDeleted: Boolean = false,
    val replyToId: String? = null,
    val reactions: List<ReactionDto> = emptyList(),
    val attachments: List<AttachmentDto> = emptyList(),
    val createdAt: Long,
    val editedAt: Long? = null,
    val deletedAt: Long? = null,
    val serverTimestamp: Long,
    val version: Int = 1,
    val checksum: String? = null,
    val deliveryStatus: String = "sending", // "sending", "sent", "delivered", "read", "failed"
    val readBy: List<String> = emptyList()
)

@Serializable
data class ReactionDto(
    val messageId: String,
    val userId: String,
    val emoji: String,
    val timestamp: Long
)

@Serializable
data class AttachmentDto(
    val id: String,
    val messageId: String,
    val type: String, // "image", "video", "audio", "file", "location"
    val url: String,
    val thumbnail: String? = null,
    val filename: String? = null,
    val fileSize: Long? = null,
    val mimeType: String? = null,
    val width: Int? = null,
    val height: Int? = null,
    val duration: Long? = null, // milliseconds
    val caption: String? = null,
    val metadata: Map<String, String> = emptyMap()
)

// Sync DTOs
@Serializable
data class SyncOperationDto(
    val id: String,
    val type: String, // "send_message", "edit_message", "delete_message", etc.
    val chatId: String,
    val data: String, // JSON serialized operation data
    val timestamp: Long,
    val retryCount: Int = 0,
    val maxRetries: Int = 3,
    val status: String = "pending" // "pending", "processing", "completed", "failed"
)

@Serializable
data class SyncRequestDto(
    val chatId: String,
    val lastSyncTimestamp: Long? = null,
    val operations: List<SyncOperationDto> = emptyList()
)

@Serializable
data class SyncResponseDto(
    val success: Boolean,
    val chatId: String,
    val messages: List<MessageDto> = emptyList(),
    val conflicts: List<ConflictDto> = emptyList(),
    val processedOperations: List<String> = emptyList(),
    val failedOperations: List<SyncOperationDto> = emptyList(),
    val serverTimestamp: Long,
    val nextSyncTimestamp: Long? = null
)

@Serializable
data class ConflictDto(
    val id: String,
    val messageId: String,
    val chatId: String,
    val type: String, // "edit_conflict", "status_conflict", "delete_conflict"
    val severity: String, // "critical", "high", "medium", "low"
    val localMessage: MessageDto,
    val remoteMessage: MessageDto,
    val detectedAt: Long,
    val autoResolvable: Boolean = false,
    val suggestedResolution: String? = null
)

@Serializable
data class ConflictResolutionDto(
    val conflictId: String,
    val strategy: String, // "local_wins", "remote_wins", "merge", "user_choice_required"
    val resolvedMessage: MessageDto? = null,
    val explanation: String? = null
)

// Real-time DTOs
@Serializable
data class PresenceUpdateDto(
    val userId: String,
    val isOnline: Boolean,
    val lastSeen: Long? = null,
    val activity: String? = null // "typing", "recording_audio", "recording_video", "idle"
)

@Serializable
data class TypingIndicatorDto(
    val chatId: String,
    val userId: String,
    val isTyping: Boolean,
    val timestamp: Long
)

// File Upload DTOs
@Serializable
data class FileUploadRequest(
    val fileName: String,
    val mimeType: String,
    val fileSize: Long,
    val checksum: String? = null
)

@Serializable
data class FileUploadResponse(
    val success: Boolean,
    val fileUrl: String,
    val uploadId: String,
    val thumbnailUrl: String? = null
)

// Error DTOs
@Serializable
data class ApiErrorDto(
    val code: String,
    val message: String,
    val details: String? = null,
    val timestamp: Long,
    val requestId: String? = null
)

@Serializable
data class ApiResponseDto<T>(
    val success: Boolean,
    val data: T? = null,
    val error: ApiErrorDto? = null,
    val timestamp: Long,
    val requestId: String? = null
)

// Health & Monitoring DTOs
@Serializable
data class HealthCheckDto(
    val status: String, // "healthy", "degraded", "unhealthy"
    val uptime: Long,
    val responseTime: Long,
    val activeConnections: Int,
    val serverLoad: Double,
    val timestamp: Long
)

@Serializable
data class ServerCapabilitiesDto(
    val version: String,
    val maxFileSize: Long,
    val supportedFileTypes: List<String>,
    val maxMessageLength: Int,
    val supportsBulkOperations: Boolean,
    val supportsRealTime: Boolean,
    val supportsFileUpload: Boolean,
    val maxChatParticipants: Int,
    val messageRetentionDays: Int
)

// Analytics DTOs
@Serializable
data class MessageAnalyticsDto(
    val chatId: String,
    val totalMessages: Int,
    val averageResponseTime: Long,
    val activeParticipants: Int,
    val messageFrequency: Map<String, Int>,
    val peakHours: List<Int>,
    val reportPeriod: String,
    val generatedAt: Long
)

// OTP Authentication DTOs
@Serializable
data class RequestOTPRequest(
    @SerialName("phone_number") val phoneNumber: String,
    @SerialName("country_code") val countryCode: String,
    @SerialName("device_info") val deviceInfo: DeviceInfo? = null
)

@Serializable
data class RequestOTPResponse(
    val success: Boolean,
    val message: String,
    val requestID: String,
    val expiresIn: Int
)

@Serializable
data class VerifyOTPRequest(
    val requestID: String,
    val code: String,
    val phoneNumber: String? = null,
    val deviceInfo: DeviceInfo? = null
)

@Serializable
data class VerifyOTPResponse(
    val success: Boolean,
    val message: String,
    val accessToken: String,
    val refreshToken: String,
    val tokenType: String,
    val expiresIn: Int,
    val user: UserInfo,
    val session: SessionInfo
)

@Serializable
data class DeviceInfo(
    val platform: String,
    val deviceModel: String? = null,
    val osVersion: String? = null,
    val appVersion: String? = null,
    val deviceID: String? = null,
    val pushToken: String? = null,
    val userAgent: String? = null,
    val ipAddress: String? = null,
    val timezone: String? = null,
    val language: String? = null
)

@Serializable
data class UserInfo(
    val id: String,
    val phoneNumber: String,
    val countryCode: String,
    val displayName: String? = null,
    val avatar: String? = null,
    val kycStatus: String,
    val kycTier: String,
    val isActive: Boolean,
    val createdAt: Long,
    val updatedAt: Long
)

@Serializable
data class SessionInfo(
    val id: String,
    val deviceInfo: String,
    val ipAddress: String,
    val createdAt: Long,
    val expiresAt: Long
)

@Serializable
data class LoginRequest(
    val email: String,
    val password: String,
    val rememberMe: Boolean = true
)

@Serializable
data class LogoutRequest(
    val refreshToken: String? = null,
    val logoutAll: Boolean = false
)

@Serializable
data class CurrentUserResponse(
    val user: UserDto,
    val authenticated: Boolean,
    val session: SessionInfo? = null
)