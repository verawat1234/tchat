//
//  DeepLink.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation

/// Deep link entity for handling and resolving deep link navigation
public struct DeepLink: Codable, Identifiable, Equatable, Hashable {

    // MARK: - Properties

    public let id: String
    public let url: String
    public let scheme: String
    public let host: String?
    public let path: String
    public let queryParameters: [String: String]
    public let fragment: String?
    public let platform: String
    public let createdAt: Date
    public let metadata: DeepLinkMetadata

    // MARK: - Initialization

    public init(
        id: String = UUID().uuidString,
        url: String,
        scheme: String,
        host: String? = nil,
        path: String,
        queryParameters: [String: String] = [:],
        fragment: String? = nil,
        platform: String = "ios",
        createdAt: Date = Date(),
        metadata: DeepLinkMetadata = DeepLinkMetadata()
    ) {
        self.id = id
        self.url = url
        self.scheme = scheme
        self.host = host
        self.path = path
        self.queryParameters = queryParameters
        self.fragment = fragment
        self.platform = platform
        self.createdAt = createdAt
        self.metadata = metadata
    }

    // MARK: - Computed Properties

    /// Check if deep link is valid
    public var isValid: Bool {
        return scheme == "tchat" && !path.isEmpty
    }

    /// Get route ID from path
    public var routeId: String {
        return path.trimmingCharacters(in: CharacterSet(charactersIn: "/"))
            .replacingOccurrences(of: "/", with: "/")
    }

    /// Get path components
    public var pathComponents: [String] {
        return path.components(separatedBy: "/").filter { !$0.isEmpty }
    }

    /// Check if requires authentication
    public var requiresAuthentication: Bool {
        return metadata.accessLevel == .private || metadata.accessLevel == .admin
    }

    /// Get expiration status
    public var isExpired: Bool {
        guard let expiresAt = metadata.expiresAt else { return false }
        return Date() > expiresAt
    }

    /// Generate full URL
    public var fullURL: URL? {
        var components = URLComponents()
        components.scheme = scheme
        components.host = host
        components.path = path

        if !queryParameters.isEmpty {
            components.queryItems = queryParameters.map { URLQueryItem(name: $0.key, value: $0.value) }
        }

        components.fragment = fragment

        return components.url
    }

    // MARK: - Static Factory Methods

    /// Create deep link from URL string
    public static func fromURL(_ urlString: String) -> DeepLink? {
        guard let url = URL(string: urlString),
              let components = URLComponents(url: url, resolvingAgainstBaseURL: false),
              let scheme = components.scheme else {
            return nil
        }

        let queryParams = components.queryItems?.reduce(into: [String: String]()) { result, item in
            result[item.name] = item.value
        } ?? [:]

        return DeepLink(
            url: urlString,
            scheme: scheme,
            host: components.host,
            path: components.path ?? "",
            queryParameters: queryParams,
            fragment: components.fragment,
            platform: "ios"
        )
    }

    /// Create deep link for route with parameters
    public static func forRoute(
        _ routeId: String,
        parameters: [String: String] = [:],
        platform: String = "ios"
    ) -> DeepLink {
        let path = "/\(routeId)"
        let url = "tchat://\(path)"

        return DeepLink(
            url: url,
            scheme: "tchat",
            host: nil,
            path: path,
            queryParameters: parameters,
            fragment: nil,
            platform: platform
        )
    }

    // MARK: - Resolution

    /// Create resolution request for this deep link
    public func createResolutionRequest(userId: String) -> DeepLinkResolutionRequest {
        return DeepLinkResolutionRequest(
            url: url,
            platform: platform,
            userId: userId
        )
    }

    /// Check if deep link matches a route pattern
    public func matchesRoute(_ routePattern: String) -> Bool {
        let pattern = routePattern.replacingOccurrences(of: ":", with: "([^/]+)")
        let regex = try? NSRegularExpression(pattern: "^\(pattern)$", options: [])
        let range = NSRange(location: 0, length: routeId.count)
        return regex?.firstMatch(in: routeId, options: [], range: range) != nil
    }

    /// Extract parameters from route pattern
    public func extractParameters(from routePattern: String) -> [String: String] {
        let pathComponents = self.pathComponents
        let patternComponents = routePattern.components(separatedBy: "/").filter { !$0.isEmpty }

        var parameters: [String: String] = [:]

        for (index, component) in patternComponents.enumerated() {
            if component.hasPrefix(":") && index < pathComponents.count {
                let paramName = String(component.dropFirst())
                parameters[paramName] = pathComponents[index]
            }
        }

        // Add query parameters
        parameters.merge(queryParameters) { (_, new) in new }

        return parameters
    }
}

// MARK: - Supporting Types

/// Deep link metadata
public struct DeepLinkMetadata: Codable, Equatable, Hashable {
    public let source: DeepLinkSource
    public let accessLevel: DeepLinkAccessLevel
    public let expiresAt: Date?
    public let campaign: String?
    public let referrer: String?
    public let customData: [String: String]

    public init(
        source: DeepLinkSource = .app,
        accessLevel: DeepLinkAccessLevel = .public,
        expiresAt: Date? = nil,
        campaign: String? = nil,
        referrer: String? = nil,
        customData: [String: String] = [:]
    ) {
        self.source = source
        self.accessLevel = accessLevel
        self.expiresAt = expiresAt
        self.campaign = campaign
        self.referrer = referrer
        self.customData = customData
    }
}

/// Deep link source
public enum DeepLinkSource: String, Codable, CaseIterable {
    case app
    case web
    case push
    case sms
    case email
    case social
    case qr
    case clipboard
}

/// Access level for deep links
public enum DeepLinkAccessLevel: String, Codable, CaseIterable {
    case `public`
    case `private`
    case admin
    case beta
}

/// Deep link resolution request
public struct DeepLinkResolutionRequest: Codable {
    public let url: String
    public let platform: String
    public let userId: String

    public init(url: String, platform: String, userId: String) {
        self.url = url
        self.platform = platform
        self.userId = userId
    }
}

/// Deep link resolution response
public struct DeepLinkResolution: Codable {
    public let routeId: String
    public let parameters: [String: Any]
    public let isValid: Bool
    public let fallbackAction: String?
    public let requiresAuth: Bool

    public init(
        routeId: String,
        parameters: [String: Any],
        isValid: Bool,
        fallbackAction: String? = nil,
        requiresAuth: Bool = false
    ) {
        self.routeId = routeId
        self.parameters = parameters
        self.isValid = isValid
        self.fallbackAction = fallbackAction
        self.requiresAuth = requiresAuth
    }

    // MARK: - Codable Conformance

    private enum CodingKeys: String, CodingKey {
        case routeId, isValid, fallbackAction, requiresAuth
        // parameters requires custom handling
    }

    public init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        routeId = try container.decode(String.self, forKey: .routeId)
        isValid = try container.decode(Bool.self, forKey: .isValid)
        fallbackAction = try container.decodeIfPresent(String.self, forKey: .fallbackAction)
        requiresAuth = try container.decode(Bool.self, forKey: .requiresAuth)
        parameters = [:] // Simplified for demo
    }

    public func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(routeId, forKey: .routeId)
        try container.encode(isValid, forKey: .isValid)
        try container.encodeIfPresent(fallbackAction, forKey: .fallbackAction)
        try container.encode(requiresAuth, forKey: .requiresAuth)
        // parameters encoding simplified for demo
    }
}

// MARK: - Default Deep Links

extension DeepLink {

    /// Common deep link patterns for the application
    public static let commonPatterns: [String: String] = [
        "chat": "/chat",
        "chat-user": "/chat/user/:userId",
        "store": "/store",
        "store-product": "/store/product/:productId",
        "social": "/social",
        "social-post": "/social/post/:postId",
        "video": "/video",
        "video-call": "/video/call/:callId",
        "more": "/more",
        "settings": "/more/settings"
    ]

    /// Create deep link for chat with user
    public static func chatWithUser(_ userId: String) -> DeepLink {
        return forRoute("chat/user/\(userId)", platform: "ios")
    }

    /// Create deep link for store product
    public static func storeProduct(_ productId: String) -> DeepLink {
        return forRoute("store/product/\(productId)", platform: "ios")
    }

    /// Create deep link for video call
    public static func videoCall(_ callId: String) -> DeepLink {
        return forRoute("video/call/\(callId)", platform: "ios")
    }

    /// Create deep link for social post
    public static func socialPost(_ postId: String) -> DeepLink {
        return forRoute("social/post/\(postId)", platform: "ios")
    }
}