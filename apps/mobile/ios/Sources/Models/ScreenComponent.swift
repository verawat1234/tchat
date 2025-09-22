//
//  ScreenComponent.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
import SwiftUI

/// Individual UI screens that correspond to web pages
/// Implements the ScreenComponent entity from data-model.md specification
public struct ScreenComponent: Codable, Identifiable, Equatable, Hashable {

    // MARK: - Types

    public enum NavigationType: String, Codable, CaseIterable {
        case tab = "tab"
        case modal = "modal"
        case push = "push"
    }

    public enum Platform: String, Codable, CaseIterable {
        case ios = "ios"
        case android = "android"
    }

    public enum AccessLevel: String, Codable, CaseIterable {
        case `public` = "public"
        case authenticated = "authenticated"
        case premium = "premium"
    }

    public enum CacheStrategy: String, Codable, CaseIterable {
        case none = "none"
        case session = "session"
        case persistent = "persistent"
    }

    public enum ScreenState: String, Codable, CaseIterable {
        case loading = "loading"
        case ready = "ready"
        case displayed = "displayed"
        case navigating = "navigating"
        case hidden = "hidden"
        case error = "error"
        case retry = "retry"
    }

    // MARK: - Properties

    public let id: String
    public let name: String
    public let route: String
    public let type: NavigationType
    public let platform: Platform
    public let webEquivalent: String
    public let requiredData: [String]
    public let optionalData: [String]
    public let accessLevel: AccessLevel
    public let cacheStrategy: CacheStrategy
    public let offlineSupport: Bool

    // State management
    public var currentState: ScreenState
    public var lastStateChange: Date
    public var errorMessage: String?

    // MARK: - Initialization

    public init(
        id: String = UUID().uuidString,
        name: String,
        route: String,
        type: NavigationType,
        platform: Platform = .ios,
        webEquivalent: String,
        requiredData: [String] = [],
        optionalData: [String] = [],
        accessLevel: AccessLevel = .public,
        cacheStrategy: CacheStrategy = .session,
        offlineSupport: Bool = false,
        currentState: ScreenState = .loading,
        lastStateChange: Date = Date(),
        errorMessage: String? = nil
    ) {
        self.id = id
        self.name = name
        self.route = route
        self.type = type
        self.platform = platform
        self.webEquivalent = webEquivalent
        self.requiredData = requiredData
        self.optionalData = optionalData
        self.accessLevel = accessLevel
        self.cacheStrategy = cacheStrategy
        self.offlineSupport = offlineSupport
        self.currentState = currentState
        self.lastStateChange = lastStateChange
        self.errorMessage = errorMessage
    }

    // MARK: - Validation

    /// Validates the screen component according to specification rules
    public func validate() throws {
        guard !id.isEmpty else {
            throw ScreenComponentError.invalidId("ID cannot be empty")
        }

        guard !name.isEmpty else {
            throw ScreenComponentError.invalidName("Name cannot be empty")
        }

        guard !route.isEmpty else {
            throw ScreenComponentError.invalidRoute("Route cannot be empty")
        }

        guard route.hasPrefix("/") else {
            throw ScreenComponentError.invalidRoute("Route must start with '/'")
        }

        guard !webEquivalent.isEmpty else {
            throw ScreenComponentError.invalidWebEquivalent("Web equivalent route cannot be empty")
        }

        // Platform-specific route validation
        switch platform {
        case .ios:
            guard route.contains(":") == false else {
                throw ScreenComponentError.invalidRoute("iOS routes should not contain ':' parameters")
            }
        case .android:
            break // Android allows more flexible routing
        }
    }

    // MARK: - State Transitions

    /// Updates the screen state following valid transitions
    public mutating func updateState(to newState: ScreenState, errorMessage: String? = nil) throws {
        guard isValidTransition(from: currentState, to: newState) else {
            throw ScreenComponentError.invalidStateTransition("Cannot transition from \(currentState) to \(newState)")
        }

        self.currentState = newState
        self.lastStateChange = Date()
        self.errorMessage = errorMessage
    }

    /// Validates state transitions according to specification
    private func isValidTransition(from: ScreenState, to: ScreenState) -> Bool {
        switch (from, to) {
        case (.loading, .ready), (.loading, .error):
            return true
        case (.ready, .displayed), (.ready, .error):
            return true
        case (.displayed, .navigating), (.displayed, .hidden), (.displayed, .error):
            return true
        case (.navigating, .hidden), (.navigating, .error):
            return true
        case (.hidden, .displayed), (.hidden, .loading):
            return true
        case (.error, .retry):
            return true
        case (.retry, .loading):
            return true
        default:
            return false
        }
    }

    // MARK: - Data Dependencies

    /// Checks if all required data dependencies are available
    public func hasRequiredData(_ availableData: [String]) -> Bool {
        return requiredData.allSatisfy { availableData.contains($0) }
    }

    /// Gets the missing required data dependencies
    public func getMissingData(_ availableData: [String]) -> [String] {
        return requiredData.filter { !availableData.contains($0) }
    }

    /// Checks if screen can be displayed offline
    public func canDisplayOffline(_ cachedData: [String]) -> Bool {
        guard offlineSupport else { return false }
        return hasRequiredData(cachedData)
    }
}

// MARK: - Error Types

public enum ScreenComponentError: LocalizedError {
    case invalidId(String)
    case invalidName(String)
    case invalidRoute(String)
    case invalidWebEquivalent(String)
    case invalidStateTransition(String)
    case missingRequiredData([String])

    public var errorDescription: String? {
        switch self {
        case .invalidId(let message),
             .invalidName(let message),
             .invalidRoute(let message),
             .invalidWebEquivalent(let message),
             .invalidStateTransition(let message):
            return message
        case .missingRequiredData(let data):
            return "Missing required data: \(data.joined(separator: ", "))"
        }
    }
}

// MARK: - Convenience Extensions

extension ScreenComponent {

    /// Creates a tab screen component
    public static func tab(
        name: String,
        route: String,
        webEquivalent: String,
        accessLevel: AccessLevel = .public
    ) -> ScreenComponent {
        return ScreenComponent(
            name: name,
            route: route,
            type: .tab,
            webEquivalent: webEquivalent,
            accessLevel: accessLevel,
            cacheStrategy: .persistent,
            offlineSupport: true
        )
    }

    /// Creates a modal screen component
    public static func modal(
        name: String,
        route: String,
        webEquivalent: String,
        requiredData: [String] = []
    ) -> ScreenComponent {
        return ScreenComponent(
            name: name,
            route: route,
            type: .modal,
            webEquivalent: webEquivalent,
            requiredData: requiredData,
            cacheStrategy: .session
        )
    }

    /// Creates a push navigation screen component
    public static func push(
        name: String,
        route: String,
        webEquivalent: String,
        requiredData: [String] = []
    ) -> ScreenComponent {
        return ScreenComponent(
            name: name,
            route: route,
            type: .push,
            webEquivalent: webEquivalent,
            requiredData: requiredData,
            cacheStrategy: .session
        )
    }
}

// MARK: - Sample Data

extension ScreenComponent {

    /// Sample screen components for testing and development
    public static let samples: [ScreenComponent] = [
        .tab(name: "Chat", route: "/chat", webEquivalent: "/chat", accessLevel: .authenticated),
        .tab(name: "Store", route: "/store", webEquivalent: "/store", accessLevel: .public),
        .tab(name: "Social", route: "/social", webEquivalent: "/social", accessLevel: .authenticated),
        .tab(name: "Video", route: "/video", webEquivalent: "/video", accessLevel: .premium),
        .tab(name: "More", route: "/more", webEquivalent: "/more", accessLevel: .public),
        .push(name: "Chat Room", route: "/chat/room", webEquivalent: "/chat/room/:id", requiredData: ["roomId", "userId"]),
        .modal(name: "Settings", route: "/settings", webEquivalent: "/settings")
    ]
}