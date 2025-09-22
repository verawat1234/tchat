//
//  NavigationState.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation

/// Core navigation state entity for cross-platform state synchronization
public struct NavigationState: Codable, Equatable, Hashable {

    // MARK: - Properties

    public let id: String
    public let userId: String
    public let sessionId: String
    public let platform: String
    public var navigationStack: [NavigationStackEntry]
    public var currentRoute: String?
    public var previousRoute: String?
    public let timestamp: Date
    public var syncVersion: Int
    public let metadata: NavigationStateMetadata

    // MARK: - Initialization

    public init(
        id: String = UUID().uuidString,
        userId: String,
        sessionId: String,
        platform: String = "ios",
        navigationStack: [NavigationStackEntry] = [],
        currentRoute: String? = nil,
        previousRoute: String? = nil,
        timestamp: Date = Date(),
        syncVersion: Int = 1,
        metadata: NavigationStateMetadata = NavigationStateMetadata()
    ) {
        self.id = id
        self.userId = userId
        self.sessionId = sessionId
        self.platform = platform
        self.navigationStack = navigationStack
        self.currentRoute = currentRoute
        self.previousRoute = previousRoute
        self.timestamp = timestamp
        self.syncVersion = syncVersion
        self.metadata = metadata
    }

    // MARK: - Computed Properties

    /// Get current navigation depth
    public var depth: Int {
        return navigationStack.count
    }

    /// Check if can go back
    public var canGoBack: Bool {
        return navigationStack.count > 1
    }

    /// Get root route
    public var rootRoute: String? {
        return navigationStack.first?.routeId
    }

    /// Get current route parameters
    public var currentParameters: [String: Any] {
        return navigationStack.last?.parameters ?? [:]
    }

    /// Check if state is synchronized
    public var isSynchronized: Bool {
        return metadata.lastSyncTimestamp != nil &&
               Date().timeIntervalSince(metadata.lastSyncTimestamp!) < 30.0 // 30 seconds
    }

    // MARK: - Navigation Operations

    /// Push new route to navigation stack
    public mutating func push(
        routeId: String,
        parameters: [String: Any] = [:],
        timestamp: Date = Date()
    ) {
        let entry = NavigationStackEntry(
            routeId: routeId,
            parameters: parameters,
            timestamp: timestamp,
            transition: .push
        )

        previousRoute = currentRoute
        currentRoute = routeId
        navigationStack.append(entry)
        incrementVersion()
    }

    /// Pop current route from navigation stack
    public mutating func pop() {
        guard navigationStack.count > 1 else { return }

        let poppedEntry = navigationStack.removeLast()
        previousRoute = currentRoute
        currentRoute = navigationStack.last?.routeId

        // Add reverse transition entry for sync
        var reverseEntry = poppedEntry
        reverseEntry.transition = .pop
        reverseEntry.timestamp = Date()

        incrementVersion()
    }

    /// Replace current route
    public mutating func replace(
        routeId: String,
        parameters: [String: Any] = [:],
        timestamp: Date = Date()
    ) {
        guard !navigationStack.isEmpty else {
            push(routeId: routeId, parameters: parameters, timestamp: timestamp)
            return
        }

        let entry = NavigationStackEntry(
            routeId: routeId,
            parameters: parameters,
            timestamp: timestamp,
            transition: .replace
        )

        previousRoute = currentRoute
        currentRoute = routeId
        navigationStack[navigationStack.count - 1] = entry
        incrementVersion()
    }

    /// Reset navigation to root
    public mutating func popToRoot() {
        guard navigationStack.count > 1 else { return }

        let rootEntry = navigationStack.first!
        previousRoute = currentRoute
        currentRoute = rootEntry.routeId
        navigationStack = [rootEntry]
        incrementVersion()
    }

    /// Reset entire navigation state
    public mutating func reset(
        toRoute routeId: String,
        parameters: [String: Any] = [:]
    ) {
        let entry = NavigationStackEntry(
            routeId: routeId,
            parameters: parameters,
            timestamp: Date(),
            transition: .reset
        )

        previousRoute = currentRoute
        currentRoute = routeId
        navigationStack = [entry]
        incrementVersion()
    }

    // MARK: - Sync Operations

    /// Mark state as synchronized
    public mutating func markSynchronized() {
        metadata = NavigationStateMetadata(
            lastSyncTimestamp: Date(),
            syncAttempts: metadata.syncAttempts,
            conflictResolutionStrategy: metadata.conflictResolutionStrategy,
            customData: metadata.customData
        )
    }

    /// Increment sync version
    private mutating func incrementVersion() {
        syncVersion += 1
    }

    /// Create sync request payload
    public func createSyncRequest() -> NavigationStateSyncRequest {
        return NavigationStateSyncRequest(
            userId: userId,
            sessionId: sessionId,
            platform: platform,
            navigationStack: navigationStack,
            timestamp: Date(),
            syncVersion: syncVersion
        )
    }

    /// Apply sync response
    public mutating func applySyncResponse(_ response: NavigationStateSyncResponse) {
        if response.success {
            self.syncVersion = response.syncVersion
            markSynchronized()
        }
    }
}

// MARK: - Supporting Types

/// Navigation stack entry with transition information
public struct NavigationStackEntry: Codable, Equatable, Hashable, Identifiable {
    public let id: String
    public let routeId: String
    public let parameters: [String: Any]
    public let timestamp: Date
    public var transition: NavigationTransition

    public init(
        id: String = UUID().uuidString,
        routeId: String,
        parameters: [String: Any] = [:],
        timestamp: Date = Date(),
        transition: NavigationTransition = .push
    ) {
        self.id = id
        self.routeId = routeId
        self.parameters = parameters
        self.timestamp = timestamp
        self.transition = transition
    }

    // MARK: - Codable Conformance

    private enum CodingKeys: String, CodingKey {
        case id, routeId, timestamp, transition
        // parameters requires custom handling
    }

    public init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        id = try container.decode(String.self, forKey: .id)
        routeId = try container.decode(String.self, forKey: .routeId)
        timestamp = try container.decode(Date.self, forKey: .timestamp)
        transition = try container.decode(NavigationTransition.self, forKey: .transition)
        parameters = [:] // Simplified for demo
    }

    public func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(id, forKey: .id)
        try container.encode(routeId, forKey: .routeId)
        try container.encode(timestamp, forKey: .timestamp)
        try container.encode(transition, forKey: .transition)
        // parameters encoding simplified for demo
    }
}

/// Navigation transition types
public enum NavigationTransition: String, Codable, CaseIterable {
    case push
    case pop
    case replace
    case reset
    case deepLink
}

/// Navigation state metadata
public struct NavigationStateMetadata: Codable, Equatable, Hashable {
    public let lastSyncTimestamp: Date?
    public let syncAttempts: Int
    public let conflictResolutionStrategy: ConflictResolutionStrategy
    public let customData: [String: String]

    public init(
        lastSyncTimestamp: Date? = nil,
        syncAttempts: Int = 0,
        conflictResolutionStrategy: ConflictResolutionStrategy = .clientWins,
        customData: [String: String] = [:]
    ) {
        self.lastSyncTimestamp = lastSyncTimestamp
        self.syncAttempts = syncAttempts
        self.conflictResolutionStrategy = conflictResolutionStrategy
        self.customData = customData
    }
}

/// Conflict resolution strategies for sync conflicts
public enum ConflictResolutionStrategy: String, Codable, CaseIterable {
    case clientWins
    case serverWins
    case merge
    case prompt
}

/// Navigation state sync request
public struct NavigationStateSyncRequest: Codable {
    public let userId: String
    public let sessionId: String
    public let platform: String
    public let navigationStack: [NavigationStackEntry]
    public let timestamp: Date
    public let syncVersion: Int

    public init(
        userId: String,
        sessionId: String,
        platform: String,
        navigationStack: [NavigationStackEntry],
        timestamp: Date,
        syncVersion: Int
    ) {
        self.userId = userId
        self.sessionId = sessionId
        self.platform = platform
        self.navigationStack = navigationStack
        self.timestamp = timestamp
        self.syncVersion = syncVersion
    }
}

/// Navigation state sync response
public struct NavigationStateSyncResponse: Codable {
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

// MARK: - Default Navigation States

extension NavigationState {

    /// Create default navigation state for main app flow
    public static func defaultState(
        userId: String,
        sessionId: String
    ) -> NavigationState {
        let rootEntry = NavigationStackEntry(
            routeId: "chat",
            parameters: [:],
            timestamp: Date(),
            transition: .reset
        )

        return NavigationState(
            userId: userId,
            sessionId: sessionId,
            platform: "ios",
            navigationStack: [rootEntry],
            currentRoute: "chat",
            previousRoute: nil,
            timestamp: Date(),
            syncVersion: 1,
            metadata: NavigationStateMetadata()
        )
    }

    /// Create navigation state from deep link
    public static func fromDeepLink(
        userId: String,
        sessionId: String,
        deepLinkUrl: String
    ) -> NavigationState? {
        // Parse deep link and create appropriate navigation state
        // This is a simplified implementation
        guard let url = URL(string: deepLinkUrl),
              url.scheme == "tchat" else {
            return nil
        }

        let pathComponents = url.pathComponents.filter { $0 != "/" }
        guard !pathComponents.isEmpty else {
            return defaultState(userId: userId, sessionId: sessionId)
        }

        let routeId = pathComponents.joined(separator: "/")
        let entry = NavigationStackEntry(
            routeId: routeId,
            parameters: [:], // Could extract query parameters
            timestamp: Date(),
            transition: .deepLink
        )

        return NavigationState(
            userId: userId,
            sessionId: sessionId,
            platform: "ios",
            navigationStack: [entry],
            currentRoute: routeId,
            previousRoute: nil,
            timestamp: Date(),
            syncVersion: 1,
            metadata: NavigationStateMetadata()
        )
    }
}