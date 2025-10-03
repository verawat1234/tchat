package com.tchat.mobile.config

/**
 * Network Configuration
 *
 * Centralized configuration for API endpoints with environment-based URL switching.
 *
 * Environment Detection:
 * - DEBUG builds → Local development server (localhost:8080)
 * - RELEASE builds → Railway production server (HTTPS)
 *
 * Usage:
 * ```kotlin
 * val baseUrl = NetworkConfig.getBaseUrl()
 * val socialUrl = NetworkConfig.getSocialApiUrl()
 * ```
 */
object NetworkConfig {

    /**
     * Environment Types
     */
    enum class Environment {
        LOCAL,      // Local development (localhost)
        RAILWAY     // Railway production (HTTPS)
    }

    /**
     * API Endpoints Configuration
     */
    private object Endpoints {
        // Local Development URLs (for emulator/simulator)
        object Local {
            const val BASE_URL = "http://10.0.2.2:8080"  // Android emulator special IP
            const val IOS_BASE_URL = "http://localhost:8080"  // iOS simulator
            const val API_VERSION = "/api/v1"
        }

        // Railway Production URLs
        object Railway {
            const val GATEWAY_URL = "https://gateway-service-production-d78d.up.railway.app"
            const val API_VERSION = "/api/v1"
        }
    }

    /**
     * Current environment
     * Can be overridden for testing
     */
    var currentEnvironment: Environment = detectEnvironment()
        private set

    /**
     * Detect environment based on build configuration
     */
    private fun detectEnvironment(): Environment {
        return if (isDebugBuild()) {
            Environment.LOCAL
        } else {
            Environment.RAILWAY
        }
    }

    /**
     * Check if this is a debug build
     * Platform-specific implementation required
     */
    expect fun isDebugBuild(): Boolean

    /**
     * Get platform-specific localhost URL
     * Platform-specific implementation required
     */
    expect fun getLocalHostUrl(): String

    /**
     * Override environment (for testing)
     */
    fun setEnvironment(environment: Environment) {
        currentEnvironment = environment
    }

    /**
     * Get base API URL based on current environment
     */
    fun getBaseUrl(): String {
        return when (currentEnvironment) {
            Environment.LOCAL -> "${getLocalHostUrl()}${Endpoints.Local.API_VERSION}"
            Environment.RAILWAY -> "${Endpoints.Railway.GATEWAY_URL}${Endpoints.Railway.API_VERSION}"
        }
    }

    /**
     * Get authentication API URL
     */
    fun getAuthApiUrl(): String {
        return "${getBaseUrl()}/auth"
    }

    /**
     * Get messaging API URL
     */
    fun getMessagingApiUrl(): String {
        return "${getBaseUrl()}/messages"
    }

    /**
     * Get video API URL
     */
    fun getVideoApiUrl(): String {
        return "${getBaseUrl()}/videos"
    }

    /**
     * Get social API URL
     */
    fun getSocialApiUrl(): String {
        return "${getBaseUrl()}/social"
    }

    /**
     * Get commerce API URL
     */
    fun getCommerceApiUrl(): String {
        return "${getBaseUrl()}/commerce"
    }

    /**
     * Get content API URL
     */
    fun getContentApiUrl(): String {
        return "${getBaseUrl()}/content"
    }

    /**
     * Get payment API URL
     */
    fun getPaymentApiUrl(): String {
        return "${getBaseUrl()}/payment"
    }

    /**
     * Get notification API URL
     */
    fun getNotificationApiUrl(): String {
        return "${getBaseUrl()}/notifications"
    }

    /**
     * Check if using Railway environment
     */
    fun isRailwayEnvironment(): Boolean {
        return currentEnvironment == Environment.RAILWAY
    }

    /**
     * Check if using local environment
     */
    fun isLocalEnvironment(): Boolean {
        return currentEnvironment == Environment.LOCAL
    }

    /**
     * Get WebSocket URL for real-time features
     */
    fun getWebSocketUrl(): String {
        return when (currentEnvironment) {
            Environment.LOCAL -> "ws://${getLocalHostUrl().removePrefix("http://")}/ws"
            Environment.RAILWAY -> "wss://${Endpoints.Railway.GATEWAY_URL.removePrefix("https://")}/ws"
        }
    }

    /**
     * Get current environment name (for logging)
     */
    fun getEnvironmentName(): String {
        return currentEnvironment.name
    }

    /**
     * Configuration info for debugging
     */
    fun getConfigInfo(): String {
        return """
            Network Configuration:
            - Environment: ${currentEnvironment.name}
            - Base URL: ${getBaseUrl()}
            - WebSocket: ${getWebSocketUrl()}
            - Is Debug: ${isDebugBuild()}
        """.trimIndent()
    }
}
