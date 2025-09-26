package com.tchat.mobile.data.network

import io.ktor.client.*
import io.ktor.client.plugins.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.plugins.logging.*
import io.ktor.client.plugins.websocket.*
import io.ktor.client.request.*
import io.ktor.http.*
import io.ktor.serialization.kotlinx.json.*
import kotlinx.serialization.json.Json
import io.github.aakira.napier.Napier

/**
 * HTTP Client Configuration for Tchat API services
 * Provides configured Ktor client instances for different environments
 */
object HttpClientConfig {

    /**
     * JSON configuration for serialization
     */
    val json = Json {
        ignoreUnknownKeys = true
        isLenient = true
        encodeDefaults = true
        prettyPrint = false
        coerceInputValues = true
    }

    /**
     * Create HTTP client with common configuration
     */
    fun createHttpClient(config: ClientConfig = ClientConfig()): HttpClient = HttpClient {
        // Content negotiation for JSON
        install(ContentNegotiation) {
            json(json)
        }

        // WebSocket support for real-time messaging
        install(WebSockets) {
            pingInterval = kotlin.time.Duration.parse("20s")
        }

        // Request/Response logging
        install(Logging) {
            logger = object : Logger {
                override fun log(message: String) {
                    Napier.v(message, tag = "HTTP")
                }
            }
            level = if (config.enableLogging) LogLevel.INFO else LogLevel.NONE
            filter { request ->
                // Don't log sensitive endpoints
                !request.url.encodedPath.contains("auth") ||
                config.logSensitiveData
            }
        }

        // Default request configuration
        install(DefaultRequest) {
            // Set default headers
            header("Accept", "application/json")
            header("Content-Type", "application/json")

            // Add API version header
            header("API-Version", "v1")

            // Add user agent
            header("User-Agent", "Tchat-Mobile/${config.appVersion}")

            // Set base URL if provided
            config.baseUrl?.let { url(it) }
        }

        // Request timeout configuration
        install(HttpTimeout) {
            requestTimeoutMillis = config.requestTimeout
            connectTimeoutMillis = config.connectTimeout
            socketTimeoutMillis = config.socketTimeout
        }

        // Retry configuration
        install(HttpRequestRetry) {
            retryOnServerErrors(maxRetries = config.maxRetries)
            retryOnException(maxRetries = config.maxRetries, retryOnTimeout = true)
            exponentialDelay()
        }

        // Custom engine configuration if needed
        engine {
            // Platform-specific engine configuration
        }
    }

    /**
     * Create authenticated HTTP client with token management
     */
    fun createAuthenticatedClient(
        authToken: String?,
        config: ClientConfig = ClientConfig()
    ): HttpClient = createHttpClient(config).config {
        install(DefaultRequest) {
            // Add authorization header if token is available
            authToken?.let { token ->
                header("Authorization", "Bearer $token")
            }
        }
    }
}

/**
 * Configuration for HTTP client behavior
 */
data class ClientConfig(
    val baseUrl: String? = null,
    val enableLogging: Boolean = true,
    val logSensitiveData: Boolean = false,
    val requestTimeout: Long = 30_000L, // 30 seconds
    val connectTimeout: Long = 10_000L, // 10 seconds
    val socketTimeout: Long = 30_000L,  // 30 seconds
    val maxRetries: Int = 3,
    val appVersion: String = "1.0.0"
)

/**
 * API endpoints configuration for different environments
 */
object ApiConfig {

    /**
     * Development environment endpoints
     */
    object Development {
        const val BASE_URL = "http://localhost"
        const val AUTH_SERVICE = "$BASE_URL:8080"
        const val CONTENT_SERVICE = "$BASE_URL:8081"
        const val COMMERCE_SERVICE = "$BASE_URL:8082"
        const val MESSAGING_SERVICE = "$BASE_URL:8083"
        const val PAYMENT_SERVICE = "$BASE_URL:8084"
        const val NOTIFICATION_SERVICE = "$BASE_URL:8085"
        const val GATEWAY_SERVICE = "$BASE_URL:8086"

        // WebSocket endpoints
        const val MESSAGING_WS = "ws://localhost:8083/ws"
        const val NOTIFICATION_WS = "ws://localhost:8085/ws"
    }

    /**
     * Production environment endpoints
     */
    object Production {
        const val BASE_URL = "https://api.tchat.app"
        const val AUTH_SERVICE = "$BASE_URL/auth"
        const val CONTENT_SERVICE = "$BASE_URL/content"
        const val COMMERCE_SERVICE = "$BASE_URL/commerce"
        const val MESSAGING_SERVICE = "$BASE_URL/messaging"
        const val PAYMENT_SERVICE = "$BASE_URL/payment"
        const val NOTIFICATION_SERVICE = "$BASE_URL/notification"
        const val GATEWAY_SERVICE = BASE_URL

        // WebSocket endpoints
        const val MESSAGING_WS = "wss://api.tchat.app/messaging/ws"
        const val NOTIFICATION_WS = "wss://api.tchat.app/notification/ws"
    }

    /**
     * Get endpoints based on environment
     */
    fun getEndpoints(isDevelopment: Boolean = true): ApiEndpoints {
        return if (isDevelopment) {
            ApiEndpoints(
                auth = Development.AUTH_SERVICE,
                content = Development.CONTENT_SERVICE,
                commerce = Development.COMMERCE_SERVICE,
                messaging = Development.MESSAGING_SERVICE,
                payment = Development.PAYMENT_SERVICE,
                notification = Development.NOTIFICATION_SERVICE,
                gateway = Development.GATEWAY_SERVICE,
                messagingWs = Development.MESSAGING_WS,
                notificationWs = Development.NOTIFICATION_WS
            )
        } else {
            ApiEndpoints(
                auth = Production.AUTH_SERVICE,
                content = Production.CONTENT_SERVICE,
                commerce = Production.COMMERCE_SERVICE,
                messaging = Production.MESSAGING_SERVICE,
                payment = Production.PAYMENT_SERVICE,
                notification = Production.NOTIFICATION_SERVICE,
                gateway = Production.GATEWAY_SERVICE,
                messagingWs = Production.MESSAGING_WS,
                notificationWs = Production.NOTIFICATION_WS
            )
        }
    }
}

/**
 * Data class containing all API endpoints
 */
data class ApiEndpoints(
    val auth: String,
    val content: String,
    val commerce: String,
    val messaging: String,
    val payment: String,
    val notification: String,
    val gateway: String,
    val messagingWs: String,
    val notificationWs: String
)

/**
 * Common API response wrapper
 */
@kotlinx.serialization.Serializable
data class ApiResponse<T>(
    val success: Boolean = false,
    val data: T? = null,
    val message: String? = null,
    val error: String? = null,
    val errors: List<String> = emptyList(),
    val timestamp: String = kotlinx.datetime.Clock.System.now().toString(),
    val requestId: String? = null
)

/**
 * Pagination response wrapper
 */
@kotlinx.serialization.Serializable
data class PaginatedResponse<T>(
    val success: Boolean = false,
    val data: List<T> = emptyList(),
    val pagination: PaginationInfo = PaginationInfo(),
    val message: String? = null,
    val error: String? = null
)

/**
 * Pagination information
 */
@kotlinx.serialization.Serializable
data class PaginationInfo(
    val page: Int = 1,
    val limit: Int = 20,
    val total: Int = 0,
    val totalPages: Int = 0,
    val hasNext: Boolean = false,
    val hasPrevious: Boolean = false
)

/**
 * API Error types for proper error handling
 */
sealed class ApiError : Exception() {
    data class NetworkError(override val message: String) : ApiError()
    data class ServerError(val code: Int, override val message: String) : ApiError()
    data class ClientError(val code: Int, override val message: String) : ApiError()
    data class UnauthorizedError(override val message: String) : ApiError()
    data class ForbiddenError(override val message: String) : ApiError()
    data class NotFoundError(override val message: String) : ApiError()
    data class ValidationError(val errors: List<String>) : ApiError() {
        override val message: String = errors.joinToString(", ")
    }
    data class SerializationError(override val message: String) : ApiError()
    data class UnknownError(override val message: String) : ApiError()
}

/**
 * Result wrapper for API calls
 */
sealed class ApiResult<out T> {
    data class Success<T>(val data: T) : ApiResult<T>()
    data class Error(val error: ApiError) : ApiResult<Nothing>()
    data object Loading : ApiResult<Nothing>()

    fun isSuccess(): Boolean = this is Success
    fun isError(): Boolean = this is Error
    fun isLoading(): Boolean = this is Loading

    fun getOrNull(): T? = when (this) {
        is Success -> data
        else -> null
    }

    fun getErrorOrNull(): ApiError? = when (this) {
        is Error -> error
        else -> null
    }
}

/**
 * Extension function to map ApiResult
 */
inline fun <T, R> ApiResult<T>.map(transform: (T) -> R): ApiResult<R> {
    return when (this) {
        is ApiResult.Success -> ApiResult.Success(transform(data))
        is ApiResult.Error -> this
        is ApiResult.Loading -> this
    }
}

/**
 * Extension function to handle ApiResult
 */
inline fun <T> ApiResult<T>.onSuccess(action: (T) -> Unit): ApiResult<T> {
    if (this is ApiResult.Success) {
        action(data)
    }
    return this
}

inline fun <T> ApiResult<T>.onError(action: (ApiError) -> Unit): ApiResult<T> {
    if (this is ApiResult.Error) {
        action(error)
    }
    return this
}

inline fun <T> ApiResult<T>.onLoading(action: () -> Unit): ApiResult<T> {
    if (this is ApiResult.Loading) {
        action()
    }
    return this
}