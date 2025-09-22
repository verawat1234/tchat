/**
 * PlatformIntegration.kt
 * TchatApp
 *
 * Created by Claude on 22/09/2024.
 */

package com.tchat.models

import kotlinx.serialization.Serializable
import java.util.*

/**
 * Native mobile features that enhance web functionality
 * Implements the PlatformIntegration entity from data-model.md specification
 */
@Serializable
data class PlatformIntegration(
    val id: String = UUID.randomUUID().toString(),
    val name: String,
    val platform: Platform = Platform.ANDROID,
    val capability: String,
    val isAvailable: Boolean = false,
    val permissions: List<Permission> = emptyList(),
    val configuration: Map<String, String> = emptyMap(),
    val fallbackBehavior: FallbackBehavior = FallbackBehavior.WEB_EQUIVALENT,

    // State management
    var currentState: IntegrationState = IntegrationState.REQUESTED,
    var lastStateChange: Long = System.currentTimeMillis(),
    var errorMessage: String? = null
) {

    // MARK: - Enums

    @Serializable
    enum class Platform(val value: String) {
        IOS("ios"),
        ANDROID("android")
    }

    @Serializable
    enum class FallbackBehavior(val value: String) {
        DISABLE("disable"),
        WEB_EQUIVALENT("web_equivalent"),
        ALTERNATE("alternate")
    }

    @Serializable
    enum class IntegrationState(val value: String) {
        REQUESTED("requested"),
        CHECKING("checking"),
        AVAILABLE("available"),
        UNAVAILABLE("unavailable"),
        CONFIGURED("configured"),
        ACTIVE("active"),
        DENIED("denied"),
        FALLBACK("fallback"),
        DISABLED("disabled")
    }

    // MARK: - Validation

    /**
     * Validates the platform integration according to specification rules
     */
    @Throws(PlatformIntegrationException::class)
    fun validate() {
        if (id.isBlank()) {
            throw PlatformIntegrationException.InvalidId("ID cannot be empty")
        }

        if (name.isBlank()) {
            throw PlatformIntegrationException.InvalidName("Name cannot be empty")
        }

        if (capability.isBlank()) {
            throw PlatformIntegrationException.InvalidCapability("Capability cannot be empty")
        }

        if (!isValidCapability(capability)) {
            throw PlatformIntegrationException.InvalidCapability(
                "Capability '$capability' is not valid for platform '$platform'"
            )
        }

        // Validate permissions for capability
        permissions.forEach { permission ->
            permission.validate()
            if (!isValidPermissionForCapability(permission.name, capability)) {
                throw PlatformIntegrationException.InvalidPermission(
                    "Permission '${permission.name}' is not valid for capability '$capability'"
                )
            }
        }
    }

    /**
     * Validates that capability is supported on the platform
     */
    private fun isValidCapability(capability: String): Boolean {
        return when (platform) {
            Platform.ANDROID -> AndroidCapabilities.values().any { it.value == capability }
            Platform.IOS -> true // iOS allows more flexible capabilities
        }
    }

    /**
     * Validates that permission is required for the capability
     */
    private fun isValidPermissionForCapability(permission: String, capability: String): Boolean {
        return AndroidCapabilities.values()
            .find { it.value == capability }
            ?.requiredPermissions
            ?.contains(permission) ?: true
    }

    // MARK: - State Transitions

    /**
     * Updates the integration state following valid transitions
     */
    @Throws(PlatformIntegrationException::class)
    fun updateState(newState: IntegrationState, errorMessage: String? = null): PlatformIntegration {
        if (!isValidStateTransition(currentState, newState)) {
            throw PlatformIntegrationException.InvalidStateTransition(
                "Cannot transition from $currentState to $newState"
            )
        }

        return copy(
            currentState = newState,
            lastStateChange = System.currentTimeMillis(),
            errorMessage = errorMessage
        )
    }

    /**
     * Validates state transitions according to specification
     */
    private fun isValidStateTransition(from: IntegrationState, to: IntegrationState): Boolean {
        return when (from to to) {
            IntegrationState.REQUESTED to IntegrationState.CHECKING -> true

            IntegrationState.CHECKING to IntegrationState.AVAILABLE,
            IntegrationState.CHECKING to IntegrationState.UNAVAILABLE -> true

            IntegrationState.AVAILABLE to IntegrationState.CONFIGURED,
            IntegrationState.AVAILABLE to IntegrationState.DENIED -> true

            IntegrationState.CONFIGURED to IntegrationState.ACTIVE,
            IntegrationState.CONFIGURED to IntegrationState.FALLBACK -> true

            IntegrationState.DENIED to IntegrationState.FALLBACK,
            IntegrationState.DENIED to IntegrationState.DISABLED -> true

            IntegrationState.UNAVAILABLE to IntegrationState.FALLBACK,
            IntegrationState.UNAVAILABLE to IntegrationState.DISABLED -> true

            else -> false
        }
    }

    // MARK: - Permission Management

    /**
     * Checks if all required permissions are granted
     */
    val hasRequiredPermissions: Boolean
        get() {
            val requiredPermissions = permissions.filter { it.required }
            return requiredPermissions.all { it.status == Permission.PermissionStatus.GRANTED }
        }

    /**
     * Gets missing required permissions
     */
    val missingRequiredPermissions: List<Permission>
        get() = permissions.filter { it.required && it.status != Permission.PermissionStatus.GRANTED }

    /**
     * Checks if integration can function with current permissions
     */
    val canFunction: Boolean
        get() = isAvailable && hasRequiredPermissions

    companion object {

        /**
         * Checks runtime availability of the integration
         */
        fun checkAvailability(capability: String): Boolean {
            return AndroidCapabilities.values()
                .find { it.value == capability }
                ?.isAvailable ?: false
        }

        /**
         * Creates a camera integration
         */
        fun camera(
            name: String = "Camera Access",
            configuration: Map<String, String> = emptyMap()
        ): PlatformIntegration {
            return PlatformIntegration(
                name = name,
                capability = AndroidCapabilities.CAMERA.value,
                isAvailable = AndroidCapabilities.CAMERA.isAvailable,
                permissions = listOf(
                    Permission(
                        name = "android.permission.CAMERA",
                        required = true,
                        requestReason = "This app needs camera access to take photos and videos."
                    )
                ),
                configuration = configuration,
                fallbackBehavior = FallbackBehavior.WEB_EQUIVALENT
            )
        }

        /**
         * Creates a notifications integration
         */
        fun notifications(
            name: String = "Push Notifications",
            configuration: Map<String, String> = emptyMap()
        ): PlatformIntegration {
            return PlatformIntegration(
                name = name,
                capability = AndroidCapabilities.NOTIFICATIONS.value,
                isAvailable = AndroidCapabilities.NOTIFICATIONS.isAvailable,
                permissions = listOf(
                    Permission(
                        name = "android.permission.POST_NOTIFICATIONS",
                        required = true,
                        requestReason = "This app needs notification access to send you important updates."
                    )
                ),
                configuration = configuration,
                fallbackBehavior = FallbackBehavior.DISABLE
            )
        }

        /**
         * Creates a biometrics integration
         */
        fun biometrics(
            name: String = "Biometric Authentication",
            configuration: Map<String, String> = emptyMap()
        ): PlatformIntegration {
            return PlatformIntegration(
                name = name,
                capability = AndroidCapabilities.BIOMETRICS.value,
                isAvailable = AndroidCapabilities.BIOMETRICS.isAvailable,
                permissions = listOf(
                    Permission(
                        name = "android.permission.USE_BIOMETRIC",
                        required = true,
                        requestReason = "This app needs biometric access for secure authentication."
                    )
                ),
                configuration = configuration,
                fallbackBehavior = FallbackBehavior.ALTERNATE
            )
        }

        /**
         * Creates a location integration
         */
        fun location(
            name: String = "Location Services",
            configuration: Map<String, String> = emptyMap()
        ): PlatformIntegration {
            return PlatformIntegration(
                name = name,
                capability = AndroidCapabilities.LOCATION.value,
                isAvailable = AndroidCapabilities.LOCATION.isAvailable,
                permissions = listOf(
                    Permission(
                        name = "android.permission.ACCESS_FINE_LOCATION",
                        required = true,
                        requestReason = "This app needs location access to provide location-based features."
                    ),
                    Permission(
                        name = "android.permission.ACCESS_COARSE_LOCATION",
                        required = false,
                        requestReason = "This app needs approximate location access."
                    )
                ),
                configuration = configuration,
                fallbackBehavior = FallbackBehavior.WEB_EQUIVALENT
            )
        }

        /**
         * Creates a microphone integration
         */
        fun microphone(
            name: String = "Microphone Access",
            configuration: Map<String, String> = emptyMap()
        ): PlatformIntegration {
            return PlatformIntegration(
                name = name,
                capability = AndroidCapabilities.MICROPHONE.value,
                isAvailable = AndroidCapabilities.MICROPHONE.isAvailable,
                permissions = listOf(
                    Permission(
                        name = "android.permission.RECORD_AUDIO",
                        required = true,
                        requestReason = "This app needs microphone access to record audio and voice messages."
                    )
                ),
                configuration = configuration,
                fallbackBehavior = FallbackBehavior.WEB_EQUIVALENT
            )
        }
    }
}

// MARK: - Permission Model

@Serializable
data class Permission(
    val name: String,
    val status: PermissionStatus = PermissionStatus.NOT_DETERMINED,
    val required: Boolean = true,
    val requestReason: String
) {

    @Serializable
    enum class PermissionStatus(val value: String) {
        GRANTED("granted"),
        DENIED("denied"),
        NOT_DETERMINED("not_determined")
    }

    /**
     * Validates the permission
     */
    @Throws(PlatformIntegrationException::class)
    fun validate() {
        if (name.isBlank()) {
            throw PlatformIntegrationException.InvalidPermission("Permission name cannot be empty")
        }

        if (requestReason.isBlank()) {
            throw PlatformIntegrationException.InvalidPermission("Request reason cannot be empty")
        }
    }
}

// MARK: - Android Capabilities

enum class AndroidCapabilities(
    val value: String,
    val isAvailable: Boolean,
    val requiredPermissions: List<String>
) {
    CAMERA(
        "camera",
        true,
        listOf("android.permission.CAMERA")
    ),
    MICROPHONE(
        "microphone",
        true,
        listOf("android.permission.RECORD_AUDIO")
    ),
    NOTIFICATIONS(
        "notifications",
        true,
        listOf("android.permission.POST_NOTIFICATIONS")
    ),
    LOCATION(
        "location",
        true,
        listOf("android.permission.ACCESS_FINE_LOCATION", "android.permission.ACCESS_COARSE_LOCATION")
    ),
    BIOMETRICS(
        "biometrics",
        true,
        listOf("android.permission.USE_BIOMETRIC")
    ),
    VIBRATION(
        "vibration",
        true,
        listOf("android.permission.VIBRATE")
    ),
    SHARE_INTENT(
        "share_intent",
        true,
        emptyList()
    ),
    DEEP_LINKING(
        "deep_linking",
        true,
        emptyList()
    ),
    BACKGROUND_PROCESSING(
        "background_processing",
        true,
        listOf("android.permission.WAKE_LOCK")
    ),
    PUSH_NOTIFICATIONS(
        "push_notifications",
        true,
        listOf("android.permission.POST_NOTIFICATIONS")
    ),
    STORAGE(
        "storage",
        true,
        listOf("android.permission.READ_EXTERNAL_STORAGE", "android.permission.WRITE_EXTERNAL_STORAGE")
    ),
    CONTACTS(
        "contacts",
        true,
        listOf("android.permission.READ_CONTACTS")
    )
}

// MARK: - Exception Types

sealed class PlatformIntegrationException(message: String) : Exception(message) {
    class InvalidId(message: String) : PlatformIntegrationException(message)
    class InvalidName(message: String) : PlatformIntegrationException(message)
    class InvalidCapability(message: String) : PlatformIntegrationException(message)
    class InvalidPermission(message: String) : PlatformIntegrationException(message)
    class InvalidStateTransition(message: String) : PlatformIntegrationException(message)
    class PermissionDenied(message: String) : PlatformIntegrationException(message)
    class UnavailableCapability(message: String) : PlatformIntegrationException(message)
}