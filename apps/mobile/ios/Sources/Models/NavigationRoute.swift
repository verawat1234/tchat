//
//  NavigationRoute.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation

/// Core navigation route entity with cross-platform compatibility
public struct NavigationRoute: Codable, Identifiable, Equatable, Hashable {

    // MARK: - Properties

    public let id: String
    public let path: String
    public let title: String
    public let component: String
    public let parameters: [String: RouteParameter]
    public let isDeepLinkable: Bool
    public let platformRestrictions: [String]
    public let parentRouteId: String?
    public let accessLevel: AccessLevel
    public let metadata: RouteMetadata

    // MARK: - Initialization

    public init(
        id: String,
        path: String,
        title: String,
        component: String,
        parameters: [String: RouteParameter] = [:],
        isDeepLinkable: Bool = true,
        platformRestrictions: [String] = [],
        parentRouteId: String? = nil,
        accessLevel: AccessLevel = .public,
        metadata: RouteMetadata = RouteMetadata()
    ) {
        self.id = id
        self.path = path
        self.title = title
        self.component = component
        self.parameters = parameters
        self.isDeepLinkable = isDeepLinkable
        self.platformRestrictions = platformRestrictions
        self.parentRouteId = parentRouteId
        self.accessLevel = accessLevel
        self.metadata = metadata
    }

    // MARK: - Computed Properties

    /// Check if route is available on current platform
    public var isAvailableOnCurrentPlatform: Bool {
        return platformRestrictions.isEmpty || platformRestrictions.contains("ios")
    }

    /// Generate full route path with parameters
    public func fullPath(with parameters: [String: Any] = [:]) -> String {
        var fullPath = path

        for (key, value) in parameters {
            fullPath = fullPath.replacingOccurrences(of: ":\(key)", with: "\(value)")
        }

        return fullPath
    }

    /// Check if route requires authentication
    public var requiresAuthentication: Bool {
        return accessLevel == .private || accessLevel == .admin
    }

    /// Get route hierarchy depth
    public var depth: Int {
        return path.components(separatedBy: "/").filter { !$0.isEmpty }.count
    }
}

// MARK: - Supporting Types

/// Route parameter definition
public struct RouteParameter: Codable, Equatable, Hashable {
    public let name: String
    public let type: ParameterType
    public let isRequired: Bool
    public let defaultValue: String?
    public let validation: ValidationRule?

    public init(
        name: String,
        type: ParameterType,
        isRequired: Bool = true,
        defaultValue: String? = nil,
        validation: ValidationRule? = nil
    ) {
        self.name = name
        self.type = type
        self.isRequired = isRequired
        self.defaultValue = defaultValue
        self.validation = validation
    }
}

/// Parameter data types
public enum ParameterType: String, Codable, CaseIterable {
    case string
    case integer
    case uuid
    case boolean
    case url
}

/// Access level for route protection
public enum AccessLevel: String, Codable, CaseIterable {
    case `public`
    case `private`
    case admin
    case beta
}

/// Validation rules for parameters
public struct ValidationRule: Codable, Equatable, Hashable {
    public let pattern: String?
    public let minLength: Int?
    public let maxLength: Int?
    public let allowedValues: [String]?

    public init(
        pattern: String? = nil,
        minLength: Int? = nil,
        maxLength: Int? = nil,
        allowedValues: [String]? = nil
    ) {
        self.pattern = pattern
        self.minLength = minLength
        self.maxLength = maxLength
        self.allowedValues = allowedValues
    }
}

/// Route metadata for additional information
public struct RouteMetadata: Codable, Equatable, Hashable {
    public let createdAt: Date
    public let updatedAt: Date
    public let version: String
    public let description: String?
    public let tags: [String]
    public let analyticsEnabled: Bool
    public let cacheStrategy: CacheStrategy

    public init(
        createdAt: Date = Date(),
        updatedAt: Date = Date(),
        version: String = "1.0.0",
        description: String? = nil,
        tags: [String] = [],
        analyticsEnabled: Bool = true,
        cacheStrategy: CacheStrategy = .default
    ) {
        self.createdAt = createdAt
        self.updatedAt = updatedAt
        self.version = version
        self.description = description
        self.tags = tags
        self.analyticsEnabled = analyticsEnabled
        self.cacheStrategy = cacheStrategy
    }
}

/// Caching strategy for route data
public enum CacheStrategy: String, Codable, CaseIterable {
    case `default`
    case noCache
    case aggressive
    case conditional
}

// MARK: - Default Routes

extension NavigationRoute {

    /// Default application routes
    public static let defaultRoutes: [NavigationRoute] = [
        NavigationRoute(
            id: "chat",
            path: "/chat",
            title: "Chat",
            component: "ChatView",
            metadata: RouteMetadata(description: "Main chat interface")
        ),
        NavigationRoute(
            id: "chat-user",
            path: "/chat/user/:userId",
            title: "User Chat",
            component: "UserChatView",
            parameters: [
                "userId": RouteParameter(
                    name: "userId",
                    type: .uuid,
                    validation: ValidationRule(pattern: "[0-9a-fA-F-]{36}")
                )
            ],
            parentRouteId: "chat",
            metadata: RouteMetadata(description: "Direct user chat")
        ),
        NavigationRoute(
            id: "store",
            path: "/store",
            title: "Store",
            component: "StoreView",
            metadata: RouteMetadata(description: "Shopping interface")
        ),
        NavigationRoute(
            id: "store-products",
            path: "/store/products",
            title: "Products",
            component: "ProductsView",
            parentRouteId: "store",
            metadata: RouteMetadata(description: "Product listing")
        ),
        NavigationRoute(
            id: "social",
            path: "/social",
            title: "Social",
            component: "SocialView",
            metadata: RouteMetadata(description: "Social feed")
        ),
        NavigationRoute(
            id: "social-feed",
            path: "/social/feed",
            title: "Feed",
            component: "FeedView",
            parentRouteId: "social",
            metadata: RouteMetadata(description: "Activity feed")
        ),
        NavigationRoute(
            id: "video",
            path: "/video",
            title: "Video",
            component: "VideoView",
            metadata: RouteMetadata(description: "Video calls")
        ),
        NavigationRoute(
            id: "video-call",
            path: "/video/call/:callId",
            title: "Video Call",
            component: "VideoCallView",
            parameters: [
                "callId": RouteParameter(
                    name: "callId",
                    type: .string,
                    validation: ValidationRule(minLength: 8, maxLength: 32)
                )
            ],
            parentRouteId: "video",
            accessLevel: .private,
            metadata: RouteMetadata(description: "Active video call")
        ),
        NavigationRoute(
            id: "more",
            path: "/more",
            title: "More",
            component: "MoreView",
            metadata: RouteMetadata(description: "Additional options")
        ),
        NavigationRoute(
            id: "more-settings",
            path: "/more/settings",
            title: "Settings",
            component: "SettingsView",
            parentRouteId: "more",
            accessLevel: .private,
            metadata: RouteMetadata(description: "App settings")
        )
    ]
}