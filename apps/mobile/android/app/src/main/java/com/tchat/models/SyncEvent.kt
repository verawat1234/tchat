/**
 * SyncEvent.kt
 * TchatApp
 *
 * Created by Claude on 22/09/2024.
 */

package com.tchat.models

import kotlinx.serialization.Serializable
import kotlinx.serialization.json.JsonElement
import java.util.*
import kotlin.math.pow

/**
 * Real-time data updates between web and mobile platforms
 * Implements the SyncEvent entity from data-model.md specification
 */
@Serializable
data class SyncEvent(
    val id: String = UUID.randomUUID().toString(),
    val type: EventType,
    val source: Platform,
    val target: Target = Target.ALL,
    val payload: Map<String, JsonElement> = emptyMap(),
    val userId: String,
    val sessionId: String,
    val timestamp: Long = System.currentTimeMillis(),
    val version: Int = 1,
    val requiresAck: Boolean = false,

    // Mutable state
    var retryCount: Int = 0,
    var status: EventStatus = EventStatus.PENDING,
    var lastRetryAt: Long? = null,
    var acknowledgedAt: Long? = null,
    var errorMessage: String? = null
) {

    // MARK: - Enums

    @Serializable
    enum class EventType(val value: String) {
        STATE_UPDATE("state_update"),
        NAVIGATION("navigation"),
        DATA_CHANGE("data_change")
    }

    @Serializable
    enum class Platform(val value: String) {
        IOS("ios"),
        ANDROID("android"),
        WEB("web")
    }

    @Serializable
    enum class Target(val value: String) {
        ALL("all"),
        IOS("ios"),
        ANDROID("android"),
        WEB("web")
    }

    @Serializable
    enum class EventStatus(val value: String) {
        PENDING("pending"),
        SENT("sent"),
        ACKNOWLEDGED("acknowledged"),
        FAILED("failed"),
        EXPIRED("expired"),
        DISCARDED("discarded")
    }

    // MARK: - Validation

    /**
     * Validates the sync event according to specification rules
     */
    @Throws(SyncEventException::class)
    fun validate() {
        if (id.isBlank()) {
            throw SyncEventException.InvalidId("ID cannot be empty")
        }

        if (userId.isBlank()) {
            throw SyncEventException.InvalidUserId("User ID cannot be empty")
        }

        if (sessionId.isBlank()) {
            throw SyncEventException.InvalidSessionId("Session ID cannot be empty")
        }

        if (timestamp > System.currentTimeMillis()) {
            throw SyncEventException.InvalidTimestamp("Timestamp cannot be in the future")
        }

        if (version <= 0) {
            throw SyncEventException.InvalidVersion("Version must be greater than 0")
        }

        if (retryCount < 0) {
            throw SyncEventException.InvalidRetryCount("Retry count cannot be negative")
        }

        // Validate that event is not expired (older than 1 hour)
        val oneHourAgo = System.currentTimeMillis() - 3600000
        if (timestamp < oneHourAgo) {
            throw SyncEventException.EventExpired("Event is older than 1 hour")
        }
    }

    // MARK: - State Transitions

    /**
     * Updates the event status following valid transitions
     */
    @Throws(SyncEventException::class)
    fun updateStatus(newStatus: EventStatus, errorMessage: String? = null): SyncEvent {
        if (!isValidStatusTransition(status, newStatus)) {
            throw SyncEventException.InvalidStatusTransition(
                "Cannot transition from $status to $newStatus"
            )
        }

        return copy(
            status = newStatus,
            errorMessage = errorMessage,
            acknowledgedAt = if (newStatus == EventStatus.ACKNOWLEDGED) System.currentTimeMillis() else acknowledgedAt,
            lastRetryAt = if (newStatus == EventStatus.FAILED) System.currentTimeMillis() else lastRetryAt
        )
    }

    /**
     * Validates status transitions according to specification
     */
    private fun isValidStatusTransition(from: EventStatus, to: EventStatus): Boolean {
        return when (from to to) {
            EventStatus.PENDING to EventStatus.SENT,
            EventStatus.PENDING to EventStatus.FAILED -> true

            EventStatus.SENT to EventStatus.ACKNOWLEDGED,
            EventStatus.SENT to EventStatus.FAILED -> true

            EventStatus.FAILED to EventStatus.PENDING,
            EventStatus.FAILED to EventStatus.DISCARDED -> true

            EventStatus.EXPIRED to EventStatus.DISCARDED -> true

            else -> false
        }
    }

    // MARK: - Retry Logic

    /**
     * Increments retry count and updates retry timestamp
     */
    @Throws(SyncEventException::class)
    fun incrementRetry(): SyncEvent {
        if (status != EventStatus.FAILED) {
            throw SyncEventException.InvalidRetryState("Cannot retry event with status: $status")
        }

        if (retryCount >= MAX_RETRIES) {
            throw SyncEventException.MaxRetriesExceeded("Maximum retry count exceeded")
        }

        return copy(
            retryCount = retryCount + 1,
            lastRetryAt = System.currentTimeMillis(),
            status = EventStatus.PENDING
        )
    }

    /**
     * Checks if the event should be retried based on retry policy
     */
    val shouldRetry: Boolean
        get() {
            if (status != EventStatus.FAILED) return false
            if (retryCount >= MAX_RETRIES) return false

            // Exponential backoff: wait 2^retryCount seconds
            lastRetryAt?.let { lastRetry ->
                val backoffTime = (2.0.pow(retryCount.toDouble()) * 1000).toLong() // 2^retryCount seconds in ms
                return System.currentTimeMillis() - lastRetry >= backoffTime
            }

            return true
        }

    /**
     * Checks if the event has expired
     */
    val isExpired: Boolean
        get() {
            val expiryTime = timestamp + EVENT_TTL
            return System.currentTimeMillis() > expiryTime
        }

    // MARK: - Acknowledgment

    /**
     * Marks the event as acknowledged
     */
    @Throws(SyncEventException::class)
    fun acknowledge(): SyncEvent {
        if (!requiresAck) {
            throw SyncEventException.AcknowledgmentNotRequired("Event does not require acknowledgment")
        }

        if (status != EventStatus.SENT) {
            throw SyncEventException.InvalidAckState("Cannot acknowledge event with status: $status")
        }

        return updateStatus(EventStatus.ACKNOWLEDGED)
    }

    // MARK: - Payload Helpers

    /**
     * Gets a string value from the payload
     */
    fun getStringValue(key: String): String? {
        return payload[key]?.let {
            if (it.toString().startsWith("\"") && it.toString().endsWith("\"")) {
                it.toString().removeSurrounding("\"")
            } else {
                it.toString()
            }
        }
    }

    /**
     * Gets an integer value from the payload
     */
    fun getIntValue(key: String): Int? {
        return payload[key]?.toString()?.toIntOrNull()
    }

    /**
     * Gets a boolean value from the payload
     */
    fun getBooleanValue(key: String): Boolean? {
        return payload[key]?.toString()?.toBooleanStrictOrNull()
    }

    companion object {

        private const val MAX_RETRIES = 3
        private const val EVENT_TTL = 3600000L // 1 hour in milliseconds

        /**
         * Creates a state update event
         */
        fun stateUpdate(
            userId: String,
            sessionId: String,
            payload: Map<String, JsonElement>,
            source: Platform = Platform.ANDROID,
            target: Target = Target.ALL
        ): SyncEvent {
            return SyncEvent(
                type = EventType.STATE_UPDATE,
                source = source,
                target = target,
                payload = payload,
                userId = userId,
                sessionId = sessionId,
                requiresAck = true
            )
        }

        /**
         * Creates a navigation event
         */
        fun navigation(
            userId: String,
            sessionId: String,
            fromRoute: String,
            toRoute: String,
            source: Platform = Platform.ANDROID
        ): SyncEvent {
            return SyncEvent(
                type = EventType.NAVIGATION,
                source = source,
                target = Target.ALL,
                payload = mapOf(
                    "fromRoute" to kotlinx.serialization.json.JsonPrimitive(fromRoute),
                    "toRoute" to kotlinx.serialization.json.JsonPrimitive(toRoute),
                    "trigger" to kotlinx.serialization.json.JsonPrimitive("user")
                ),
                userId = userId,
                sessionId = sessionId,
                requiresAck = false
            )
        }

        /**
         * Creates a data change event
         */
        fun dataChange(
            userId: String,
            sessionId: String,
            entity: String,
            action: String,
            entityId: String,
            source: Platform = Platform.ANDROID
        ): SyncEvent {
            return SyncEvent(
                type = EventType.DATA_CHANGE,
                source = source,
                target = Target.ALL,
                payload = mapOf(
                    "entity" to kotlinx.serialization.json.JsonPrimitive(entity),
                    "action" to kotlinx.serialization.json.JsonPrimitive(action),
                    "entityId" to kotlinx.serialization.json.JsonPrimitive(entityId)
                ),
                userId = userId,
                sessionId = sessionId,
                requiresAck = true
            )
        }
    }
}

// MARK: - Exception Types

sealed class SyncEventException(message: String) : Exception(message) {
    class InvalidId(message: String) : SyncEventException(message)
    class InvalidUserId(message: String) : SyncEventException(message)
    class InvalidSessionId(message: String) : SyncEventException(message)
    class InvalidTimestamp(message: String) : SyncEventException(message)
    class InvalidVersion(message: String) : SyncEventException(message)
    class InvalidRetryCount(message: String) : SyncEventException(message)
    class EventExpired(message: String) : SyncEventException(message)
    class InvalidStatusTransition(message: String) : SyncEventException(message)
    class InvalidRetryState(message: String) : SyncEventException(message)
    class MaxRetriesExceeded(message: String) : SyncEventException(message)
    class AcknowledgmentNotRequired(message: String) : SyncEventException(message)
    class InvalidAckState(message: String) : SyncEventException(message)
}