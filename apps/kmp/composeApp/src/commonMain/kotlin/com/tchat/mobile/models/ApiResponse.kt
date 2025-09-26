package com.tchat.mobile.models

import kotlinx.serialization.Serializable
import kotlinx.datetime.*

/**
 * T034: ApiResponse model for standardized API responses
 *
 * Comprehensive API response handling with pagination, error management,
 * metadata, caching, and cross-platform consistency. Provides unified
 * response structure for all API endpoints with type safety.
 */
@Serializable
data class ApiResponse<T>(
    val success: Boolean,
    val data: T? = null,
    val error: ApiError? = null,
    val message: String? = null,
    val code: String? = null,         // Response code (e.g., "USER_NOT_FOUND")
    val timestamp: String,            // ISO 8601 timestamp
    val requestId: String? = null,    // Trace ID for debugging
    val pagination: Pagination? = null,
    val metadata: ResponseMetadata = ResponseMetadata(),
    val links: Map<String, String> = emptyMap(), // HATEOAS links
    val warnings: List<ApiWarning> = emptyList(),
    val debug: DebugInfo? = null      // Debug info (development only)
)

@Serializable
data class ApiError(
    val type: ErrorType,
    val code: String,                 // Error code (e.g., "VALIDATION_FAILED")
    val message: String,              // Human-readable message
    val details: String? = null,      // Technical details
    val field: String? = null,        // Field that caused the error
    val validation: ValidationErrors? = null,
    val retryable: Boolean = false,   // Whether client should retry
    val retryAfter: Long? = null,     // Seconds to wait before retry
    val documentation: String? = null, // Link to error documentation
    val context: Map<String, String> = emptyMap(),
    val traceId: String? = null,
    val correlationId: String? = null,
    val innerError: ApiError? = null,  // Nested error for chaining
    val timestamp: String,
    val severity: ErrorSeverity = ErrorSeverity.ERROR
)

enum class ErrorType {
    // Client Errors (4xx)
    BAD_REQUEST,          // 400
    UNAUTHORIZED,         // 401
    FORBIDDEN,           // 403
    NOT_FOUND,           // 404
    METHOD_NOT_ALLOWED,  // 405
    NOT_ACCEPTABLE,      // 406
    REQUEST_TIMEOUT,     // 408
    CONFLICT,            // 409
    GONE,                // 410
    PRECONDITION_FAILED, // 412
    PAYLOAD_TOO_LARGE,   // 413
    UNSUPPORTED_MEDIA,   // 415
    UNPROCESSABLE_ENTITY, // 422
    TOO_MANY_REQUESTS,   // 429

    // Server Errors (5xx)
    INTERNAL_SERVER_ERROR, // 500
    NOT_IMPLEMENTED,     // 501
    BAD_GATEWAY,         // 502
    SERVICE_UNAVAILABLE, // 503
    GATEWAY_TIMEOUT,     // 504

    // Application Errors
    VALIDATION_ERROR,
    BUSINESS_RULE_VIOLATION,
    RESOURCE_EXHAUSTED,
    DEPENDENCY_FAILURE,

    // Network/Client Errors
    NETWORK_ERROR,
    CONNECTION_ERROR,
    TIMEOUT_ERROR,
    DNS_ERROR,
    SSL_ERROR,

    // Data Errors
    PARSE_ERROR,
    SERIALIZATION_ERROR,
    ENCODING_ERROR,

    // Authentication/Authorization
    TOKEN_EXPIRED,
    TOKEN_INVALID,
    INSUFFICIENT_PERMISSIONS,
    ACCOUNT_LOCKED,
    ACCOUNT_SUSPENDED,

    // Rate Limiting
    RATE_LIMIT_EXCEEDED,
    QUOTA_EXCEEDED,
    CONCURRENT_LIMIT_EXCEEDED,

    // Feature/Service
    FEATURE_DISABLED,
    MAINTENANCE_MODE,
    DEPRECATED_ENDPOINT,

    // Unknown/Fallback
    UNKNOWN_ERROR
}

enum class ErrorSeverity {
    DEBUG, INFO, WARNING, ERROR, CRITICAL, FATAL
}

@Serializable
data class ValidationErrors(
    val fields: Map<String, List<FieldError>>,
    val global: List<String> = emptyList(),
    val count: Int = fields.values.sumOf { it.size } + global.size
)

@Serializable
data class FieldError(
    val code: String,                 // Error code (e.g., "REQUIRED", "EMAIL_INVALID")
    val message: String,              // User-friendly message
    val value: String? = null,        // The invalid value
    val constraint: String? = null,   // The constraint that was violated
    val path: String? = null          // JSON path to the field
)

@Serializable
data class ApiWarning(
    val code: String,
    val message: String,
    val field: String? = null,
    val severity: WarningSeverity = WarningSeverity.MEDIUM,
    val dismissible: Boolean = true,
    val action: String? = null,       // Suggested action
    val learnMoreUrl: String? = null
)

enum class WarningSeverity {
    LOW, MEDIUM, HIGH
}

@Serializable
data class Pagination(
    val page: Int,                    // Current page (1-based)
    val pageSize: Int,                // Items per page
    val totalPages: Int,              // Total number of pages
    val totalItems: Long,             // Total number of items
    val hasNext: Boolean,             // Whether there's a next page
    val hasPrev: Boolean,             // Whether there's a previous page
    val itemsOnPage: Int,             // Items on current page
    val startIndex: Long,             // Start index of items on page (0-based)
    val endIndex: Long,               // End index of items on page (0-based)
    val sortBy: String? = null,       // Current sort field
    val sortDirection: SortDirection = SortDirection.ASC,
    val filters: Map<String, String> = emptyMap(),
    val cursors: PaginationCursors? = null, // Cursor-based pagination
    val links: PaginationLinks? = null,     // HATEOAS navigation links
    val performance: PaginationPerformance? = null
)

enum class SortDirection {
    ASC, DESC
}

@Serializable
data class PaginationCursors(
    val before: String? = null,       // Cursor for previous page
    val after: String? = null,        // Cursor for next page
    val current: String? = null       // Current cursor
)

@Serializable
data class PaginationLinks(
    val first: String? = null,        // URL to first page
    val previous: String? = null,     // URL to previous page
    val next: String? = null,         // URL to next page
    val last: String? = null          // URL to last page
)

@Serializable
data class PaginationPerformance(
    val queryTime: Long,              // Query execution time (ms)
    val totalTime: Long,              // Total request time (ms)
    val cached: Boolean = false,      // Whether result was cached
    val estimatedTotal: Boolean = false // Whether total is estimated
)

@Serializable
data class ResponseMetadata(
    val version: String = "1.0",      // API version
    val environment: String? = null,  // Environment (prod, staging, dev)
    val region: String? = null,       // Server region
    val processingTime: Long? = null, // Processing time in milliseconds
    val cached: Boolean = false,      // Whether response was cached
    val cacheAge: Long? = null,       // Cache age in seconds
    val cacheTtl: Long? = null,       // Cache TTL in seconds
    val etag: String? = null,         // ETag for caching
    val lastModified: String? = null, // Last modified timestamp
    val rateLimits: RateLimitInfo? = null,
    val features: List<String> = emptyList(), // Enabled features
    val experiments: List<String> = emptyList(), // Active experiments
    val clientInfo: ClientInfo? = null,
    val serverInfo: ServerInfo? = null,
    val security: SecurityInfo? = null,
    val compliance: ComplianceInfo? = null,
    val custom: Map<String, String> = emptyMap()
)

@Serializable
data class RateLimitInfo(
    val limit: Int,                   // Request limit
    val remaining: Int,               // Remaining requests
    val reset: Long,                  // Reset time (Unix timestamp)
    val resetTime: String,            // Reset time (ISO 8601)
    val window: Long,                 // Window duration in seconds
    val policy: String? = null        // Rate limiting policy
)

@Serializable
data class ClientInfo(
    val userAgent: String? = null,
    val ipAddress: String? = null,
    val platform: String? = null,    // "android", "ios", "web", "desktop"
    val version: String? = null,      // Client app version
    val buildNumber: String? = null,
    val locale: String? = null,
    val timezone: String? = null,
    val sessionId: String? = null
)

@Serializable
data class ServerInfo(
    val hostname: String? = null,
    val version: String? = null,      // API server version
    val commit: String? = null,       // Git commit hash
    val buildTime: String? = null,    // Build timestamp
    val uptime: Long? = null,         // Server uptime in seconds
    val loadAverage: Double? = null   // System load average
)

@Serializable
data class SecurityInfo(
    val encrypted: Boolean = true,    // Whether response is encrypted
    val signed: Boolean = false,      // Whether response is signed
    val checksum: String? = null,     // Response checksum
    val certificate: String? = null,  // Certificate info
    val permissions: List<String> = emptyList() // User permissions
)

@Serializable
data class ComplianceInfo(
    val gdpr: Boolean = false,        // GDPR compliance
    val ccpa: Boolean = false,        // CCPA compliance
    val coppa: Boolean = false,       // COPPA compliance
    val pci: Boolean = false,         // PCI compliance
    val hipaa: Boolean = false,       // HIPAA compliance
    val dataRetention: Long? = null,  // Data retention period (days)
    val consentId: String? = null     // Consent tracking ID
)

@Serializable
data class DebugInfo(
    val sql: List<SqlQuery> = emptyList(),    // SQL queries executed
    val cache: CacheDebug? = null,            // Cache debug info
    val performance: PerformanceDebug? = null, // Performance metrics
    val stack: List<String> = emptyList(),    // Stack trace (errors only)
    val logs: List<LogEntry> = emptyList(),   // Debug logs
    val memory: MemoryInfo? = null,           // Memory usage
    val database: DatabaseInfo? = null        // Database metrics
)

@Serializable
data class SqlQuery(
    val query: String,
    val parameters: List<String> = emptyList(),
    val executionTime: Long,          // Milliseconds
    val rowsAffected: Int = 0,
    val rowsReturned: Int = 0,
    val cached: Boolean = false
)

@Serializable
data class CacheDebug(
    val hits: Int = 0,
    val misses: Int = 0,
    val keys: List<String> = emptyList(),
    val expiration: Map<String, String> = emptyMap()
)

@Serializable
data class PerformanceDebug(
    val totalTime: Long,              // Total request time (ms)
    val dbTime: Long = 0,             // Database time (ms)
    val cacheTime: Long = 0,          // Cache time (ms)
    val networkTime: Long = 0,        // Network time (ms)
    val renderTime: Long = 0,         // Rendering time (ms)
    val bottlenecks: List<String> = emptyList()
)

@Serializable
data class LogEntry(
    val level: String,                // DEBUG, INFO, WARN, ERROR
    val message: String,
    val timestamp: String,
    val context: Map<String, String> = emptyMap()
)

@Serializable
data class MemoryInfo(
    val used: Long,                   // Used memory (bytes)
    val available: Long,              // Available memory (bytes)
    val peak: Long = 0,               // Peak memory usage (bytes)
    val gcCollections: Int = 0        // GC collection count
)

@Serializable
data class DatabaseInfo(
    val connections: Int = 0,         // Active connections
    val queries: Int = 0,             // Total queries
    val slowQueries: Int = 0,         // Slow queries
    val locks: Int = 0,               // Database locks
    val version: String? = null       // Database version
)

/**
 * Specialized response types
 */
@Serializable
data class ListResponse<T>(
    val items: List<T>,
    val pagination: Pagination,
    val filters: Map<String, String> = emptyMap(),
    val sorting: SortInfo? = null,
    val aggregations: Map<String, AggregationResult> = emptyMap()
)

@Serializable
data class SortInfo(
    val field: String,
    val direction: SortDirection,
    val available: List<String> = emptyList() // Available sort fields
)

@Serializable
data class AggregationResult(
    val type: AggregationType,
    val field: String,
    val value: Double,
    val count: Long? = null
)

enum class AggregationType {
    COUNT, SUM, AVERAGE, MIN, MAX, DISTINCT_COUNT
}

@Serializable
data class BulkResponse<T>(
    val successful: List<BulkItem<T>>,
    val failed: List<BulkError>,
    val summary: BulkSummary
)

@Serializable
data class BulkItem<T>(
    val index: Int,
    val data: T,
    val id: String? = null,
    val status: BulkStatus = BulkStatus.SUCCESS
)

@Serializable
data class BulkError(
    val index: Int,
    val error: ApiError,
    val originalData: String? = null // JSON representation
)

enum class BulkStatus {
    SUCCESS, PARTIAL_SUCCESS, FAILED, SKIPPED
}

@Serializable
data class BulkSummary(
    val total: Int,
    val successful: Int,
    val failed: Int,
    val skipped: Int = 0,
    val processingTime: Long,         // Milliseconds
    val validationErrors: Int = 0,
    val businessErrors: Int = 0
)

@Serializable
data class HealthResponse(
    val status: HealthStatus,
    val timestamp: String,
    val version: String,
    val uptime: Long,                 // Seconds
    val checks: Map<String, HealthCheck> = emptyMap(),
    val metrics: Map<String, Double> = emptyMap()
)

enum class HealthStatus {
    HEALTHY, DEGRADED, UNHEALTHY, MAINTENANCE
}

@Serializable
data class HealthCheck(
    val status: HealthStatus,
    val responseTime: Long? = null,   // Milliseconds
    val message: String? = null,
    val lastCheck: String,
    val metadata: Map<String, String> = emptyMap()
)

/**
 * Response state management
 */
@Serializable
data class ResponseState<T>(
    val data: T? = null,
    val loading: Boolean = false,
    val error: ApiError? = null,
    val refreshing: Boolean = false,
    val lastFetched: String? = null,
    val cacheValid: Boolean = true,
    val retryCount: Int = 0,
    val maxRetries: Int = 3,
    val retryDelay: Long = 1000,      // Milliseconds
    val offline: Boolean = false
)

enum class ResponseStatus {
    IDLE, LOADING, SUCCESS, ERROR, REFRESHING, OFFLINE
}

/**
 * ApiResponse utilities and extensions
 */
fun <T> ApiResponse<T>.isSuccess(): Boolean = success && error == null

fun <T> ApiResponse<T>.isError(): Boolean = !success || error != null

fun <T> ApiResponse<T>.hasData(): Boolean = data != null

fun <T> ApiResponse<T>.getDataOrNull(): T? = if (isSuccess()) data else null

fun <T> ApiResponse<T>.getErrorOrNull(): ApiError? = if (isError()) error else null

fun <T> ApiResponse<T>.requireData(): T =
    data ?: throw IllegalStateException("Response data is null")

fun <T> ApiResponse<T>.getMessageOrDefault(default: String): String =
    message ?: error?.message ?: default

fun <T> ApiResponse<T>.isPaginated(): Boolean = pagination != null

fun <T> ApiResponse<T>.hasNextPage(): Boolean = pagination?.hasNext == true

fun <T> ApiResponse<T>.hasPreviousPage(): Boolean = pagination?.hasPrev == true

fun <T> ApiResponse<T>.isRetryable(): Boolean = error?.retryable == true

fun <T> ApiResponse<T>.shouldRetry(): Boolean =
    isError() && isRetryable() && error?.retryAfter?.let { it <= 300 } != false

fun <T> ApiResponse<T>.getCacheAge(): Long? = metadata.cacheAge

fun <T> ApiResponse<T>.isCached(): Boolean = metadata.cached

fun <T> ApiResponse<T>.isRateLimited(): Boolean = error?.type == ErrorType.TOO_MANY_REQUESTS

fun <T> ApiResponse<T>.getRateLimitReset(): Long? = error?.retryAfter

fun <T> ApiResponse<T>.hasWarnings(): Boolean = warnings.isNotEmpty()

fun <T> ApiResponse<T>.getHighPriorityWarnings(): List<ApiWarning> =
    warnings.filter { it.severity == WarningSeverity.HIGH }

fun <T> ApiResponse<T>.isFromProduction(): Boolean =
    metadata.environment == "production" || metadata.environment == "prod"

fun <T> ApiResponse<T>.getProcessingTime(): Long =
    metadata.processingTime ?: debug?.performance?.totalTime ?: 0L

/**
 * Error handling utilities
 */
fun ApiError.isClientError(): Boolean = type.name.startsWith("4") ||
    type in listOf(ErrorType.BAD_REQUEST, ErrorType.UNAUTHORIZED, ErrorType.FORBIDDEN, ErrorType.NOT_FOUND)

fun ApiError.isServerError(): Boolean = type in listOf(
    ErrorType.INTERNAL_SERVER_ERROR, ErrorType.BAD_GATEWAY,
    ErrorType.SERVICE_UNAVAILABLE, ErrorType.GATEWAY_TIMEOUT
)

fun ApiError.isNetworkError(): Boolean = type in listOf(
    ErrorType.NETWORK_ERROR, ErrorType.CONNECTION_ERROR,
    ErrorType.TIMEOUT_ERROR, ErrorType.DNS_ERROR, ErrorType.SSL_ERROR
)

fun ApiError.isAuthenticationError(): Boolean = type in listOf(
    ErrorType.UNAUTHORIZED, ErrorType.TOKEN_EXPIRED, ErrorType.TOKEN_INVALID,
    ErrorType.ACCOUNT_LOCKED, ErrorType.ACCOUNT_SUSPENDED
)

fun ApiError.isValidationError(): Boolean = type == ErrorType.VALIDATION_ERROR ||
    validation != null

fun ApiError.getUserFriendlyMessage(): String = when (type) {
    ErrorType.NETWORK_ERROR -> "Please check your internet connection and try again"
    ErrorType.UNAUTHORIZED -> "Please log in to continue"
    ErrorType.FORBIDDEN -> "You don't have permission to access this resource"
    ErrorType.NOT_FOUND -> "The requested resource was not found"
    ErrorType.TOO_MANY_REQUESTS -> "Too many requests. Please try again later"
    ErrorType.SERVICE_UNAVAILABLE -> "Service is temporarily unavailable"
    ErrorType.INTERNAL_SERVER_ERROR -> "Something went wrong. Please try again"
    else -> message
}

fun ApiError.getValidationMessage(field: String): String? =
    validation?.fields?.get(field)?.firstOrNull()?.message

/**
 * Pagination utilities
 */
fun Pagination.getNextPage(): Int? = if (hasNext) page + 1 else null

fun Pagination.getPreviousPage(): Int? = if (hasPrev) page - 1 else null

fun Pagination.getPageRange(): IntRange = 1..totalPages

fun Pagination.getCurrentRange(): LongRange = startIndex..endIndex

fun Pagination.getProgress(): Double =
    if (totalPages > 0) page.toDouble() / totalPages else 0.0

fun Pagination.isEmpty(): Boolean = totalItems == 0L

/**
 * Factory functions for common responses
 */
fun <T> successResponse(
    data: T,
    message: String? = null,
    pagination: Pagination? = null
): ApiResponse<T> = ApiResponse(
    success = true,
    data = data,
    message = message,
    pagination = pagination,
    timestamp = Clock.System.now().toString()
)

fun <T> errorResponse(
    error: ApiError,
    message: String? = null
): ApiResponse<T> = ApiResponse(
    success = false,
    error = error,
    message = message ?: error.message,
    timestamp = Clock.System.now().toString()
)

fun <T> loadingResponse(): ResponseState<T> = ResponseState(
    loading = true,
    lastFetched = Clock.System.now().toString()
)

fun <T> cachedResponse(data: T, cacheAge: Long): ApiResponse<T> = ApiResponse(
    success = true,
    data = data,
    timestamp = Clock.System.now().toString(),
    metadata = ResponseMetadata(
        cached = true,
        cacheAge = cacheAge
    )
)

/**
 * Transformation utilities
 */
fun <T, R> ApiResponse<T>.map(transform: (T) -> R): ApiResponse<R> =
    if (isSuccess() && data != null) {
        ApiResponse(
            success = success,
            data = transform(data!!),
            error = error,
            message = message,
            code = code,
            timestamp = timestamp,
            requestId = requestId,
            pagination = pagination,
            metadata = metadata,
            links = links,
            warnings = warnings,
            debug = debug
        )
    } else {
        ApiResponse(
            success = success,
            data = null,
            error = error,
            message = message,
            code = code,
            timestamp = timestamp,
            requestId = requestId,
            pagination = pagination,
            metadata = metadata,
            links = links,
            warnings = warnings,
            debug = debug
        )
    }

fun <T> ApiResponse<T>.recover(fallback: T): ApiResponse<T> =
    if (isError()) copy(success = true, data = fallback, error = null)
    else this

fun <T> ApiResponse<T>.onSuccess(action: (T) -> Unit): ApiResponse<T> = apply {
    if (isSuccess() && data != null) action(data)
}

fun <T> ApiResponse<T>.onError(action: (ApiError) -> Unit): ApiResponse<T> = apply {
    if (isError() && error != null) action(error)
}

fun <T> ApiResponse<T>.fold(
    onSuccess: (T) -> Unit,
    onError: (ApiError) -> Unit
): ApiResponse<T> = apply {
    when {
        isSuccess() && data != null -> onSuccess(data)
        isError() && error != null -> onError(error)
    }
}