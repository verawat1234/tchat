//
//  SyncEvent.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation

/// Real-time data updates between web and mobile platforms
/// Implements the SyncEvent entity from data-model.md specification
public struct SyncEvent: Codable, Identifiable, Equatable, Hashable {

    // MARK: - Types

    public enum EventType: String, Codable, CaseIterable {
        case stateUpdate = "state_update"
        case navigation = "navigation"
        case dataChange = "data_change"
    }

    public enum Platform: String, Codable, CaseIterable {
        case ios = "ios"
        case android = "android"
        case web = "web"
    }

    public enum Target: String, Codable, CaseIterable {
        case all = "all"
        case ios = "ios"
        case android = "android"
        case web = "web"
    }

    public enum EventStatus: String, Codable, CaseIterable {
        case pending = "pending"
        case sent = "sent"
        case acknowledged = "acknowledged"
        case failed = "failed"
        case expired = "expired"
        case discarded = "discarded"
    }

    // MARK: - Properties

    public let id: String
    public let type: EventType
    public let source: Platform
    public let target: Target
    public let payload: [String: AnyCodable]
    public let userId: String
    public let sessionId: String
    public let timestamp: Date
    public let version: Int
    public let requiresAck: Bool

    // Mutable state
    public var retryCount: Int
    public var status: EventStatus
    public var lastRetryAt: Date?
    public var acknowledgedAt: Date?
    public var errorMessage: String?

    // MARK: - Initialization

    public init(
        id: String = UUID().uuidString,
        type: EventType,
        source: Platform,
        target: Target = .all,
        payload: [String: AnyCodable] = [:],
        userId: String,
        sessionId: String,
        timestamp: Date = Date(),
        version: Int = 1,
        requiresAck: Bool = false,
        retryCount: Int = 0,
        status: EventStatus = .pending,
        lastRetryAt: Date? = nil,
        acknowledgedAt: Date? = nil,
        errorMessage: String? = nil
    ) {
        self.id = id
        self.type = type
        self.source = source
        self.target = target
        self.payload = payload
        self.userId = userId
        self.sessionId = sessionId
        self.timestamp = timestamp
        self.version = version
        self.requiresAck = requiresAck
        self.retryCount = retryCount
        self.status = status
        self.lastRetryAt = lastRetryAt
        self.acknowledgedAt = acknowledgedAt
        self.errorMessage = errorMessage
    }

    // MARK: - Validation

    /// Validates the sync event according to specification rules
    public func validate() throws {
        guard !id.isEmpty else {
            throw SyncEventError.invalidId("ID cannot be empty")
        }

        guard !userId.isEmpty else {
            throw SyncEventError.invalidUserId("User ID cannot be empty")
        }

        guard !sessionId.isEmpty else {
            throw SyncEventError.invalidSessionId("Session ID cannot be empty")
        }

        guard timestamp <= Date() else {
            throw SyncEventError.invalidTimestamp("Timestamp cannot be in the future")
        }

        guard version > 0 else {
            throw SyncEventError.invalidVersion("Version must be greater than 0")
        }

        guard retryCount >= 0 else {
            throw SyncEventError.invalidRetryCount("Retry count cannot be negative")
        }

        // Validate that event is not expired (older than 1 hour)
        let oneHourAgo = Date().addingTimeInterval(-3600)
        guard timestamp > oneHourAgo else {
            throw SyncEventError.eventExpired("Event is older than 1 hour")
        }
    }

    // MARK: - State Transitions

    /// Updates the event status following valid transitions
    public mutating func updateStatus(to newStatus: EventStatus, errorMessage: String? = nil) throws {
        guard isValidStatusTransition(from: status, to: newStatus) else {
            throw SyncEventError.invalidStatusTransition("Cannot transition from \(status) to \(newStatus)")
        }

        self.status = newStatus
        self.errorMessage = errorMessage

        switch newStatus {
        case .acknowledged:
            self.acknowledgedAt = Date()
        case .failed:
            self.lastRetryAt = Date()
        default:
            break
        }
    }

    /// Validates status transitions according to specification
    private func isValidStatusTransition(from: EventStatus, to: EventStatus) -> Bool {
        switch (from, to) {
        case (.pending, .sent), (.pending, .failed):
            return true
        case (.sent, .acknowledged), (.sent, .failed):
            return true
        case (.failed, .pending), (.failed, .discarded):
            return true
        case (.expired, .discarded):
            return true
        default:
            return false
        }
    }

    // MARK: - Retry Logic

    /// Increments retry count and updates retry timestamp
    public mutating func incrementRetry() throws {
        guard status == .failed else {
            throw SyncEventError.invalidRetryState("Cannot retry event with status: \(status)")
        }

        guard retryCount < Constants.maxRetries else {
            throw SyncEventError.maxRetriesExceeded("Maximum retry count exceeded")
        }

        self.retryCount += 1
        self.lastRetryAt = Date()
        self.status = .pending
    }

    /// Checks if the event should be retried based on retry policy
    public var shouldRetry: Bool {
        guard status == .failed else { return false }
        guard retryCount < Constants.maxRetries else { return false }

        // Exponential backoff: wait 2^retryCount seconds
        if let lastRetry = lastRetryAt {
            let backoffTime = TimeInterval(1 << retryCount) // 2^retryCount
            return Date().timeIntervalSince(lastRetry) >= backoffTime
        }

        return true
    }

    /// Checks if the event has expired
    public var isExpired: Bool {
        let expiryTime = timestamp.addingTimeInterval(Constants.eventTTL)
        return Date() > expiryTime
    }

    // MARK: - Acknowledgment

    /// Marks the event as acknowledged
    public mutating func acknowledge() throws {
        guard requiresAck else {
            throw SyncEventError.acknowledgmentNotRequired("Event does not require acknowledgment")
        }

        guard status == .sent else {
            throw SyncEventError.invalidAckState("Cannot acknowledge event with status: \(status)")
        }

        try updateStatus(to: .acknowledged)
    }

    // MARK: - Payload Helpers

    /// Gets a typed value from the payload
    public func getPayloadValue<T: Codable>(_ key: String, as type: T.Type) -> T? {
        guard let anyCodable = payload[key] else { return nil }
        return anyCodable.value as? T
    }

    /// Sets a typed value in the payload (returns new instance)
    public func withPayloadValue<T: Codable>(_ key: String, value: T) -> SyncEvent {
        var newPayload = payload
        newPayload[key] = AnyCodable(value)

        return SyncEvent(
            id: id, type: type, source: source, target: target,
            payload: newPayload, userId: userId, sessionId: sessionId,
            timestamp: timestamp, version: version, requiresAck: requiresAck,
            retryCount: retryCount, status: status, lastRetryAt: lastRetryAt,
            acknowledgedAt: acknowledgedAt, errorMessage: errorMessage
        )
    }
}

// MARK: - Constants

private extension SyncEvent {
    enum Constants {
        static let maxRetries = 3
        static let eventTTL: TimeInterval = 3600 // 1 hour
    }
}

// MARK: - Error Types

public enum SyncEventError: LocalizedError {
    case invalidId(String)
    case invalidUserId(String)
    case invalidSessionId(String)
    case invalidTimestamp(String)
    case invalidVersion(String)
    case invalidRetryCount(String)
    case eventExpired(String)
    case invalidStatusTransition(String)
    case invalidRetryState(String)
    case maxRetriesExceeded(String)
    case acknowledgmentNotRequired(String)
    case invalidAckState(String)

    public var errorDescription: String? {
        switch self {
        case .invalidId(let message),
             .invalidUserId(let message),
             .invalidSessionId(let message),
             .invalidTimestamp(let message),
             .invalidVersion(let message),
             .invalidRetryCount(let message),
             .eventExpired(let message),
             .invalidStatusTransition(let message),
             .invalidRetryState(let message),
             .maxRetriesExceeded(let message),
             .acknowledgmentNotRequired(let message),
             .invalidAckState(let message):
            return message
        }
    }
}

// MARK: - Factory Methods

extension SyncEvent {

    /// Creates a state update event
    public static func stateUpdate(
        userId: String,
        sessionId: String,
        payload: [String: AnyCodable],
        source: Platform = .ios,
        target: Target = .all
    ) -> SyncEvent {
        return SyncEvent(
            type: .stateUpdate,
            source: source,
            target: target,
            payload: payload,
            userId: userId,
            sessionId: sessionId,
            requiresAck: true
        )
    }

    /// Creates a navigation event
    public static func navigation(
        userId: String,
        sessionId: String,
        fromRoute: String,
        toRoute: String,
        source: Platform = .ios
    ) -> SyncEvent {
        let payload: [String: AnyCodable] = [
            "fromRoute": AnyCodable(fromRoute),
            "toRoute": AnyCodable(toRoute),
            "trigger": AnyCodable("user")
        ]

        return SyncEvent(
            type: .navigation,
            source: source,
            target: .all,
            payload: payload,
            userId: userId,
            sessionId: sessionId,
            requiresAck: false
        )
    }

    /// Creates a data change event
    public static func dataChange(
        userId: String,
        sessionId: String,
        entity: String,
        action: String,
        entityId: String,
        source: Platform = .ios
    ) -> SyncEvent {
        let payload: [String: AnyCodable] = [
            "entity": AnyCodable(entity),
            "action": AnyCodable(action),
            "entityId": AnyCodable(entityId)
        ]

        return SyncEvent(
            type: .dataChange,
            source: source,
            target: .all,
            payload: payload,
            userId: userId,
            sessionId: sessionId,
            requiresAck: true
        )
    }
}

// MARK: - AnyCodable Helper

public struct AnyCodable: Codable, Hashable {
    public let value: Any

    public init<T>(_ value: T?) {
        self.value = value ?? ()
    }

    public init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()

        if container.decodeNil() {
            self.init(())
        } else if let bool = try? container.decode(Bool.self) {
            self.init(bool)
        } else if let int = try? container.decode(Int.self) {
            self.init(int)
        } else if let double = try? container.decode(Double.self) {
            self.init(double)
        } else if let string = try? container.decode(String.self) {
            self.init(string)
        } else if let array = try? container.decode([AnyCodable].self) {
            self.init(array.map { $0.value })
        } else if let dictionary = try? container.decode([String: AnyCodable].self) {
            self.init(dictionary.mapValues { $0.value })
        } else {
            throw DecodingError.dataCorruptedError(in: container, debugDescription: "AnyCodable value cannot be decoded")
        }
    }

    public func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()

        switch value {
        case is Void:
            try container.encodeNil()
        case let bool as Bool:
            try container.encode(bool)
        case let int as Int:
            try container.encode(int)
        case let double as Double:
            try container.encode(double)
        case let string as String:
            try container.encode(string)
        case let array as [Any?]:
            try container.encode(array.map { AnyCodable($0) })
        case let dictionary as [String: Any?]:
            try container.encode(dictionary.mapValues { AnyCodable($0) })
        default:
            let context = EncodingError.Context(codingPath: container.codingPath, debugDescription: "AnyCodable value cannot be encoded")
            throw EncodingError.invalidValue(value, context)
        }
    }

    public func hash(into hasher: inout Hasher) {
        switch value {
        case let bool as Bool:
            hasher.combine(bool)
        case let int as Int:
            hasher.combine(int)
        case let double as Double:
            hasher.combine(double)
        case let string as String:
            hasher.combine(string)
        default:
            hasher.combine(0)
        }
    }

    public static func == (lhs: AnyCodable, rhs: AnyCodable) -> Bool {
        switch (lhs.value, rhs.value) {
        case is (Void, Void):
            return true
        case (let lhs as Bool, let rhs as Bool):
            return lhs == rhs
        case (let lhs as Int, let rhs as Int):
            return lhs == rhs
        case (let lhs as Double, let rhs as Double):
            return lhs == rhs
        case (let lhs as String, let rhs as String):
            return lhs == rhs
        default:
            return false
        }
    }
}