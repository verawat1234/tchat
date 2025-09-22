//
//  ComponentState.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
import SwiftUI

/// Component state entity for UI component state management and synchronization
public struct ComponentState: Codable, Identifiable {

    // MARK: - Properties

    public let id: String
    public let componentId: String
    public let instanceId: String
    public let userId: String
    public let sessionId: String
    public var state: [String: Any]
    public let timestamp: Date
    public var version: Int
    public let platform: String
    public var metadata: ComponentStateMetadata

    // MARK: - Initialization

    public init(
        id: String = UUID().uuidString,
        componentId: String,
        instanceId: String,
        userId: String,
        sessionId: String,
        state: [String: Any] = [:],
        timestamp: Date = Date(),
        version: Int = 1,
        platform: String = "ios",
        metadata: ComponentStateMetadata = ComponentStateMetadata()
    ) {
        self.id = id
        self.componentId = componentId
        self.instanceId = instanceId
        self.userId = userId
        self.sessionId = sessionId
        self.state = state
        self.timestamp = timestamp
        self.version = version
        self.platform = platform
        self.metadata = metadata
    }

    // MARK: - Computed Properties

    /// Check if state is synchronized
    public var isSynchronized: Bool {
        return metadata.lastSyncTimestamp != nil &&
               Date().timeIntervalSince(metadata.lastSyncTimestamp!) < 10.0 // 10 seconds
    }

    /// Get state keys
    public var stateKeys: [String] {
        return Array(state.keys)
    }

    /// Check if state is empty
    public var isEmpty: Bool {
        return state.isEmpty
    }

    /// Get state size (number of properties)
    public var stateSize: Int {
        return state.count
    }

    /// Calculate state hash for conflict detection
    public var stateHash: String {
        let sortedKeys = state.keys.sorted()
        let stateString = sortedKeys.map { "\($0):\(state[$0] ?? "")" }.joined(separator: "|")
        return String(stateString.hashValue)
    }

    // MARK: - State Operations

    /// Update specific state property
    public mutating func updateProperty(_ key: String, value: Any) {
        state[key] = value
        incrementVersion()
    }

    /// Update multiple state properties
    public mutating func updateProperties(_ updates: [String: Any]) {
        for (key, value) in updates {
            state[key] = value
        }
        incrementVersion()
    }

    /// Remove state property
    public mutating func removeProperty(_ key: String) {
        state.removeValue(forKey: key)
        incrementVersion()
    }

    /// Clear all state
    public mutating func clearState() {
        state.removeAll()
        incrementVersion()
    }

    /// Merge state with another state
    public mutating func mergeState(_ otherState: [String: Any], strategy: MergeStrategy = .overwrite) {
        switch strategy {
        case .overwrite:
            for (key, value) in otherState {
                state[key] = value
            }
        case .keepExisting:
            for (key, value) in otherState {
                if state[key] == nil {
                    state[key] = value
                }
            }
        case .merge:
            for (key, value) in otherState {
                state[key] = value
            }
        }
        incrementVersion()
    }

    /// Get state property with type safety
    public func getProperty<T>(_ key: String, as type: T.Type) -> T? {
        return state[key] as? T
    }

    /// Get state property with default value
    public func getProperty<T>(_ key: String, default defaultValue: T) -> T {
        return state[key] as? T ?? defaultValue
    }

    // MARK: - Sync Operations

    /// Create state sync request
    public func createSyncRequest() -> ComponentStateSyncRequest {
        return ComponentStateSyncRequest(
            userId: userId,
            sessionId: sessionId,
            platform: platform,
            componentStates: [self],
            timestamp: Date(),
            syncVersion: version
        )
    }

    /// Apply sync response
    public mutating func applySyncResponse(_ response: ComponentStateSyncResponse) {
        if response.success {
            self.version = response.syncVersion
            markSynchronized()
        }
    }

    /// Mark state as synchronized
    public mutating func markSynchronized() {
        metadata = ComponentStateMetadata(
            lastSyncTimestamp: Date(),
            syncAttempts: metadata.syncAttempts,
            conflictCount: metadata.conflictCount,
            customData: metadata.customData
        )
    }

    /// Increment version
    private mutating func incrementVersion() {
        version += 1
    }

    /// Detect conflicts with another state
    public func hasConflictWith(_ otherState: ComponentState) -> Bool {
        return componentId == otherState.componentId &&
               instanceId == otherState.instanceId &&
               version != otherState.version &&
               stateHash != otherState.stateHash
    }

    /// Resolve conflict with another state
    public mutating func resolveConflictWith(
        _ otherState: ComponentState,
        strategy: ComponentConflictResolutionStrategy = .newestWins
    ) {
        switch strategy {
        case .newestWins:
            if otherState.timestamp > timestamp {
                state = otherState.state
                version = otherState.version
            }
        case .highestVersionWins:
            if otherState.version > version {
                state = otherState.state
                version = otherState.version
            }
        case .merge:
            mergeState(otherState.state, strategy: .merge)
        case .manual:
            // Manual resolution required - emit conflict event
            break
        }
    }
}

// MARK: - Supporting Types

/// Merge strategy for state updates
public enum MergeStrategy: String, Codable, CaseIterable {
    case overwrite
    case keepExisting
    case merge
}

/// Conflict resolution strategy for component state
public enum ComponentConflictResolutionStrategy: String, Codable, CaseIterable {
    case newestWins
    case highestVersionWins
    case merge
    case manual
}

/// Component state metadata
public struct ComponentStateMetadata: Codable, Equatable, Hashable {
    public let lastSyncTimestamp: Date?
    public let syncAttempts: Int
    public let conflictCount: Int
    public let customData: [String: String]

    public init(
        lastSyncTimestamp: Date? = nil,
        syncAttempts: Int = 0,
        conflictCount: Int = 0,
        customData: [String: String] = [:]
    ) {
        self.lastSyncTimestamp = lastSyncTimestamp
        self.syncAttempts = syncAttempts
        self.conflictCount = conflictCount
        self.customData = customData
    }
}

/// Component state sync request
public struct ComponentStateSyncRequest: Codable {
    public let userId: String
    public let sessionId: String
    public let platform: String
    public let componentStates: [ComponentState]
    public let timestamp: Date
    public let syncVersion: Int

    public init(
        userId: String,
        sessionId: String,
        platform: String,
        componentStates: [ComponentState],
        timestamp: Date,
        syncVersion: Int
    ) {
        self.userId = userId
        self.sessionId = sessionId
        self.platform = platform
        self.componentStates = componentStates
        self.timestamp = timestamp
        self.syncVersion = syncVersion
    }
}

/// Component state sync response
public struct ComponentStateSyncResponse: Codable {
    public let success: Bool
    public let syncVersion: Int
    public let conflictsResolved: [String]
    public let timestamp: Date

    public init(
        success: Bool,
        syncVersion: Int,
        conflictsResolved: [String] = [],
        timestamp: Date = Date()
    ) {
        self.success = success
        self.syncVersion = syncVersion
        self.conflictsResolved = conflictsResolved
        self.timestamp = timestamp
    }
}

// MARK: - Codable Conformance

extension ComponentState {
    private enum CodingKeys: String, CodingKey {
        case id, componentId, instanceId, userId, sessionId
        case timestamp, version, platform, metadata
        // state requires custom encoding/decoding
    }

    public init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        id = try container.decode(String.self, forKey: .id)
        componentId = try container.decode(String.self, forKey: .componentId)
        instanceId = try container.decode(String.self, forKey: .instanceId)
        userId = try container.decode(String.self, forKey: .userId)
        sessionId = try container.decode(String.self, forKey: .sessionId)
        timestamp = try container.decode(Date.self, forKey: .timestamp)
        version = try container.decode(Int.self, forKey: .version)
        platform = try container.decode(String.self, forKey: .platform)
        metadata = try container.decode(ComponentStateMetadata.self, forKey: .metadata)
        state = [:] // Simplified for demo
    }

    public func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(id, forKey: .id)
        try container.encode(componentId, forKey: .componentId)
        try container.encode(instanceId, forKey: .instanceId)
        try container.encode(userId, forKey: .userId)
        try container.encode(sessionId, forKey: .sessionId)
        try container.encode(timestamp, forKey: .timestamp)
        try container.encode(version, forKey: .version)
        try container.encode(platform, forKey: .platform)
        try container.encode(metadata, forKey: .metadata)
        // state encoding simplified for demo
    }
}

// MARK: - Default Component States

extension ComponentState {

    /// Create default component state for chat message
    public static func chatMessageState(
        instanceId: String,
        userId: String,
        sessionId: String,
        isRead: Bool = false,
        isSelected: Bool = false
    ) -> ComponentState {
        return ComponentState(
            componentId: "chat-message",
            instanceId: instanceId,
            userId: userId,
            sessionId: sessionId,
            state: [
                "isRead": isRead,
                "isSelected": isSelected,
                "timestamp": Date()
            ],
            platform: "ios"
        )
    }

    /// Create default component state for user avatar
    public static func userAvatarState(
        instanceId: String,
        userId: String,
        sessionId: String,
        isOnline: Bool = false,
        lastSeen: Date? = nil
    ) -> ComponentState {
        var state: [String: Any] = [
            "isOnline": isOnline
        ]

        if let lastSeen = lastSeen {
            state["lastSeen"] = lastSeen
        }

        return ComponentState(
            componentId: "user-avatar",
            instanceId: instanceId,
            userId: userId,
            sessionId: sessionId,
            state: state,
            platform: "ios"
        )
    }

    /// Create default component state for navigation tab
    public static func navigationTabState(
        instanceId: String,
        userId: String,
        sessionId: String,
        isActive: Bool = false,
        hasNotification: Bool = false
    ) -> ComponentState {
        return ComponentState(
            componentId: "navigation-tab",
            instanceId: instanceId,
            userId: userId,
            sessionId: sessionId,
            state: [
                "isActive": isActive,
                "hasNotification": hasNotification
            ],
            platform: "ios"
        )
    }
}