//
//  NavigationEnvironment.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import SwiftUI
import Combine

/// Core navigation environment for managing app-wide navigation state
public class NavigationEnvironment: ObservableObject {

    /// Main navigation path for the app
    @Published public var navigationPath = NavigationPathManager()

    /// Deep link handler
    @Published public var deepLinkHandler = DeepLinkHandlerManager()

    /// Route registry for managing available routes
    @Published public var routeRegistry = RouteRegistryManager()

    public init() {}
}

// MARK: - NavigationPathManager

public class NavigationPathManager: ObservableObject {
    @Published public var path = NavigationPath()
    @Published public var currentRoute: String?

    public init() {}

    public func navigate(to route: String, parameters: [String: Any] = [:]) {
        currentRoute = route
        // Implementation will be expanded with route-specific navigation
    }

    public func goBack() {
        if !path.isEmpty {
            path.removeLast()
        }
    }

    public func popToRoot() {
        path = NavigationPath()
        currentRoute = nil
    }
}

// MARK: - Deep Link Handler

// MARK: - DeepLinkHandlerManager

public class DeepLinkHandlerManager: ObservableObject {
    @Published public var pendingDeepLink: URL?

    public init() {}

    public func handle(url: URL) -> Bool {
        // Parse URL and determine if it's a valid deep link
        guard url.scheme == "tchat" else { return false }

        pendingDeepLink = url
        return true
    }

    public func processDeepLink(_ url: URL) -> String? {
        // Convert URL to route string
        let pathComponents = url.pathComponents.filter { $0 != "/" }
        return pathComponents.joined(separator: "/")
    }
}

// MARK: - Route Registry

// MARK: - RouteRegistryManager

public class RouteRegistryManager: ObservableObject {
    @Published public var availableRoutes: Set<String> = []

    public init() {
        setupDefaultRoutes()
    }

    private func setupDefaultRoutes() {
        availableRoutes = [
            "chat",
            "store",
            "social",
            "video",
            "more",
            "chat/user",
            "store/products",
            "social/feed",
            "video/call",
            "more/settings"
        ]
    }

    public func isValidRoute(_ route: String) -> Bool {
        return availableRoutes.contains(route)
    }

    public func registerRoute(_ route: String) {
        availableRoutes.insert(route)
    }
}