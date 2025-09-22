/**
 * ScreenComponent.kt
 * TchatApp
 *
 * Created by Claude on 22/09/2024.
 */

package com.tchat.models

import kotlinx.serialization.Serializable
import java.util.*

/**
 * Individual UI screens that correspond to web pages
 * Implements the ScreenComponent entity from data-model.md specification
 */
@Serializable
data class ScreenComponent(
    val id: String = UUID.randomUUID().toString(),
    val name: String,
    val route: String,
    val type: NavigationType,
    val platform: Platform = Platform.ANDROID,
    val webEquivalent: String,
    val requiredData: List<String> = emptyList(),
    val optionalData: List<String> = emptyList(),
    val accessLevel: AccessLevel = AccessLevel.PUBLIC,
    val cacheStrategy: CacheStrategy = CacheStrategy.SESSION,
    val offlineSupport: Boolean = false,

    // State management
    var currentState: ScreenState = ScreenState.LOADING,
    var lastStateChange: Long = System.currentTimeMillis(),
    var errorMessage: String? = null
) {

    // MARK: - Enums

    @Serializable
    enum class NavigationType(val value: String) {
        TAB("tab"),
        MODAL("modal"),
        PUSH("push")
    }

    @Serializable
    enum class Platform(val value: String) {
        IOS("ios"),
        ANDROID("android")
    }

    @Serializable
    enum class AccessLevel(val value: String) {
        PUBLIC("public"),
        AUTHENTICATED("authenticated"),
        PREMIUM("premium")
    }

    @Serializable
    enum class CacheStrategy(val value: String) {
        NONE("none"),
        SESSION("session"),
        PERSISTENT("persistent")
    }

    @Serializable
    enum class ScreenState(val value: String) {
        LOADING("loading"),
        READY("ready"),
        DISPLAYED("displayed"),
        NAVIGATING("navigating"),
        HIDDEN("hidden"),
        ERROR("error"),
        RETRY("retry")
    }

    // MARK: - Validation

    /**
     * Validates the screen component according to specification rules
     */
    @Throws(ScreenComponentException::class)
    fun validate() {
        if (id.isBlank()) {
            throw ScreenComponentException.InvalidId("ID cannot be empty")
        }

        if (name.isBlank()) {
            throw ScreenComponentException.InvalidName("Name cannot be empty")
        }

        if (route.isBlank()) {
            throw ScreenComponentException.InvalidRoute("Route cannot be empty")
        }

        if (!route.startsWith("/")) {
            throw ScreenComponentException.InvalidRoute("Route must start with '/'")
        }

        if (webEquivalent.isBlank()) {
            throw ScreenComponentException.InvalidWebEquivalent("Web equivalent route cannot be empty")
        }

        // Platform-specific route validation
        when (platform) {
            Platform.IOS -> {
                if (route.contains(":")) {
                    throw ScreenComponentException.InvalidRoute("iOS routes should not contain ':' parameters")
                }
            }
            Platform.ANDROID -> {
                // Android allows more flexible routing
            }
        }
    }

    // MARK: - State Transitions

    /**
     * Updates the screen state following valid transitions
     */
    @Throws(ScreenComponentException::class)
    fun updateState(newState: ScreenState, errorMessage: String? = null): ScreenComponent {
        if (!isValidTransition(currentState, newState)) {
            throw ScreenComponentException.InvalidStateTransition(
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
    private fun isValidTransition(from: ScreenState, to: ScreenState): Boolean {
        return when (from to to) {
            ScreenState.LOADING to ScreenState.READY,
            ScreenState.LOADING to ScreenState.ERROR -> true

            ScreenState.READY to ScreenState.DISPLAYED,
            ScreenState.READY to ScreenState.ERROR -> true

            ScreenState.DISPLAYED to ScreenState.NAVIGATING,
            ScreenState.DISPLAYED to ScreenState.HIDDEN,
            ScreenState.DISPLAYED to ScreenState.ERROR -> true

            ScreenState.NAVIGATING to ScreenState.HIDDEN,
            ScreenState.NAVIGATING to ScreenState.ERROR -> true

            ScreenState.HIDDEN to ScreenState.DISPLAYED,
            ScreenState.HIDDEN to ScreenState.LOADING -> true

            ScreenState.ERROR to ScreenState.RETRY -> true
            ScreenState.RETRY to ScreenState.LOADING -> true

            else -> false
        }
    }

    // MARK: - Data Dependencies

    /**
     * Checks if all required data dependencies are available
     */
    fun hasRequiredData(availableData: List<String>): Boolean {
        return requiredData.all { availableData.contains(it) }
    }

    /**
     * Gets the missing required data dependencies
     */
    fun getMissingData(availableData: List<String>): List<String> {
        return requiredData.filter { !availableData.contains(it) }
    }

    /**
     * Checks if screen can be displayed offline
     */
    fun canDisplayOffline(cachedData: List<String>): Boolean {
        return offlineSupport && hasRequiredData(cachedData)
    }

    companion object {

        /**
         * Creates a tab screen component
         */
        fun tab(
            name: String,
            route: String,
            webEquivalent: String,
            accessLevel: AccessLevel = AccessLevel.PUBLIC
        ): ScreenComponent {
            return ScreenComponent(
                name = name,
                route = route,
                type = NavigationType.TAB,
                webEquivalent = webEquivalent,
                accessLevel = accessLevel,
                cacheStrategy = CacheStrategy.PERSISTENT,
                offlineSupport = true
            )
        }

        /**
         * Creates a modal screen component
         */
        fun modal(
            name: String,
            route: String,
            webEquivalent: String,
            requiredData: List<String> = emptyList()
        ): ScreenComponent {
            return ScreenComponent(
                name = name,
                route = route,
                type = NavigationType.MODAL,
                webEquivalent = webEquivalent,
                requiredData = requiredData,
                cacheStrategy = CacheStrategy.SESSION
            )
        }

        /**
         * Creates a push navigation screen component
         */
        fun push(
            name: String,
            route: String,
            webEquivalent: String,
            requiredData: List<String> = emptyList()
        ): ScreenComponent {
            return ScreenComponent(
                name = name,
                route = route,
                type = NavigationType.PUSH,
                webEquivalent = webEquivalent,
                requiredData = requiredData,
                cacheStrategy = CacheStrategy.SESSION
            )
        }

        /**
         * Sample screen components for testing and development
         */
        val samples = listOf(
            tab("Chat", "/chat", "/chat", AccessLevel.AUTHENTICATED),
            tab("Store", "/store", "/store", AccessLevel.PUBLIC),
            tab("Social", "/social", "/social", AccessLevel.AUTHENTICATED),
            tab("Video", "/video", "/video", AccessLevel.PREMIUM),
            tab("More", "/more", "/more", AccessLevel.PUBLIC),
            push("Chat Room", "/chat/room", "/chat/room/:id", listOf("roomId", "userId")),
            modal("Settings", "/settings", "/settings")
        )
    }
}

// MARK: - Exception Types

sealed class ScreenComponentException(message: String) : Exception(message) {
    class InvalidId(message: String) : ScreenComponentException(message)
    class InvalidName(message: String) : ScreenComponentException(message)
    class InvalidRoute(message: String) : ScreenComponentException(message)
    class InvalidWebEquivalent(message: String) : ScreenComponentException(message)
    class InvalidStateTransition(message: String) : ScreenComponentException(message)
    class MissingRequiredData(data: List<String>) : ScreenComponentException(
        "Missing required data: ${data.joinToString(", ")}"
    )
}