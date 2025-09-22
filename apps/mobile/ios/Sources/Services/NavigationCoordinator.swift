//
//  NavigationCoordinator.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
import SwiftUI
import Combine

/// Central navigation coordinator for managing app-wide navigation state and synchronization
@MainActor
public class NavigationCoordinator: ObservableObject {

    // MARK: - Published Properties

    @Published public var currentNavigationState: NavigationState?
    @Published public var isNavigating: Bool = false
    @Published public var lastError: NavigationError?

    // MARK: - Navigation Properties

    @Published public var navigationPath = NavigationPath()
    @Published public var selectedTab: TabNavigationView.Tab = .chat

    // MARK: - Private Properties

    private var cancellables = Set<AnyCancellable>()
    private let syncService: StateSyncService
    private let performanceMonitor: PerformanceMonitor
    private let routeValidator: RouteValidator
    private let deepLinkHandler: DeepLinkHandlerManager
    private let stateManager: NavigationStateManager

    // MARK: - Initialization

    public init(
        syncService: StateSyncService = StateSyncService.shared,
        performanceMonitor: PerformanceMonitor = PerformanceMonitor.shared,
        routeValidator: RouteValidator = RouteValidator(),
        deepLinkHandler: DeepLinkHandlerManager = DeepLinkHandlerManager(),
        stateManager: NavigationStateManager = NavigationStateManager()
    ) {
        self.syncService = syncService
        self.performanceMonitor = performanceMonitor
        self.routeValidator = routeValidator
        self.deepLinkHandler = deepLinkHandler
        self.stateManager = stateManager

        setupObservers()
        loadInitialState()
    }

    // MARK: - Public Methods

    /// Navigate to a specific route with parameters
    public func navigate(
        to routeId: String,
        parameters: [String: Any] = [:],
        transition: NavigationTransition = .push
    ) async throws {
        guard !isNavigating else {
            throw NavigationError.navigationInProgress
        }

        isNavigating = true
        defer { isNavigating = false }

        do {
            // Validate route
            try await routeValidator.validate(routeId: routeId, parameters: parameters)

            // Update local state
            try await updateNavigationState { state in
                switch transition {
                case .push:
                    state.push(routeId: routeId, parameters: parameters)
                case .replace:
                    state.replace(routeId: routeId, parameters: parameters)
                case .reset:
                    state.reset(toRoute: routeId, parameters: parameters)
                default:
                    state.push(routeId: routeId, parameters: parameters)
                }
            }

            // Sync with remote if needed
            try await syncNavigationState()

            // Update navigation path
            navigationPath.append(routeId)

        } catch {
            lastError = error as? NavigationError ?? NavigationError.unknown(error)
            throw error
        }
    }

    /// Go back in navigation stack
    public func goBack() async throws {
        guard let state = currentNavigationState, state.canGoBack else {
            throw NavigationError.cannotGoBack
        }

        isNavigating = true
        defer { isNavigating = false }

        do {
            try await updateNavigationState { state in
                state.pop()
            }

            try await syncNavigationState()
            navigationPath.removeLast()

        } catch {
            lastError = error as? NavigationError ?? NavigationError.unknown(error)
            throw error
        }
    }

    /// Pop to root of navigation stack
    public func popToRoot() async throws {
        guard let state = currentNavigationState, state.depth > 1 else {
            return
        }

        isNavigating = true
        defer { isNavigating = false }

        do {
            try await updateNavigationState { state in
                state.popToRoot()
            }

            try await syncNavigationState()
            navigationPath.removeLast(navigationPath.count)

        } catch {
            lastError = error as? NavigationError ?? NavigationError.unknown(error)
            throw error
        }
    }

    /// Handle deep link navigation
    public func handleDeepLink(_ url: URL) async throws {
        guard deepLinkHandler.handle(url: url) else {
            throw NavigationError.invalidDeepLink(url)
        }

        guard let routeId = deepLinkHandler.processDeepLink(url) else {
            throw NavigationError.deepLinkProcessingFailed(url)
        }

        // Extract parameters from URL
        let parameters = extractParameters(from: url)

        try await navigate(to: routeId, parameters: parameters, transition: .deepLink)
    }

    /// Reset navigation to specific state
    public func resetNavigation(
        to routeId: String,
        parameters: [String: Any] = [:]
    ) async throws {
        isNavigating = true
        defer { isNavigating = false }

        do {
            let newState = NavigationState.defaultState(
                userId: getCurrentUserId(),
                sessionId: getCurrentSessionId()
            )

            var mutableState = newState
            mutableState.reset(toRoute: routeId, parameters: parameters)

            currentNavigationState = mutableState
            try await syncNavigationState()
            navigationPath.removeLast(navigationPath.count)

        } catch {
            lastError = error as? NavigationError ?? NavigationError.unknown(error)
            throw error
        }
    }

    /// Restore navigation state from sync
    public func restoreNavigationState(_ state: NavigationState) async throws {
        currentNavigationState = state

        // Restore navigation path
        if let currentRoute = state.currentRoute {
            navigationPath.append(currentRoute)
        }
    }

    // MARK: - Private Methods

    private func setupObservers() {
        // Note: NavigationPath doesn't expose internal changes in native SwiftUI
        // Path changes are managed through the navigate methods above

        // Observe deep link events
        deepLinkHandler.$pendingDeepLink
            .compactMap { $0 }
            .sink { [weak self] url in
                Task { @MainActor in
                    try? await self?.handleDeepLink(url)
                    self?.deepLinkHandler.pendingDeepLink = nil
                }
            }
            .store(in: &cancellables)
    }

    private func loadInitialState() {
        Task {
            do {
                // Try to load saved navigation state
                if let savedState = try await stateManager.loadNavigationState() {
                    currentNavigationState = savedState
                } else {
                    // Create default state
                    currentNavigationState = NavigationState.defaultState(
                        userId: getCurrentUserId(),
                        sessionId: getCurrentSessionId()
                    )
                }
            } catch {
                // Fallback to default state
                currentNavigationState = NavigationState.defaultState(
                    userId: getCurrentUserId(),
                    sessionId: getCurrentSessionId()
                )
            }
        }
    }

    private func handleRouteChange(_ route: String) async {
        // Update current navigation state to reflect route change
        guard var state = currentNavigationState else { return }

        if state.currentRoute != route {
            state.currentRoute = route
            currentNavigationState = state

            // Save state
            try? await stateManager.saveNavigationState(state)
        }
    }

    private func updateNavigationState(_ update: (inout NavigationState) -> Void) async throws {
        guard var state = currentNavigationState else {
            throw NavigationError.noNavigationState
        }

        update(&state)
        currentNavigationState = state

        // Save updated state
        try await stateManager.saveNavigationState(state)
    }

    private func syncNavigationState() async throws {
        guard let state = currentNavigationState else { return }

        do {
            let syncRequest = state.createSyncRequest()
            let response = try await syncService.syncNavigationState(request: syncRequest)

            if response.success {
                var updatedState = state
                updatedState = updatedState.applySyncResponse(response)
                currentNavigationState = updatedState
            }
        } catch {
            // Log sync error but don't fail navigation
            print("Navigation sync failed: \(error)")
        }
    }

    private func extractParameters(from url: URL) -> [String: Any] {
        guard let components = URLComponents(url: url, resolvingAgainstBaseURL: false),
              let queryItems = components.queryItems else {
            return [:]
        }

        var parameters: [String: Any] = [:]
        for item in queryItems {
            parameters[item.name] = item.value
        }

        return parameters
    }

    private func getCurrentUserId() -> String {
        // TODO: Get from authentication service
        return "current_user_id"
    }

    private func getCurrentSessionId() -> String {
        // TODO: Get from session service
        return UUID().uuidString
    }
}

// MARK: - Supporting Services

/// Navigation sync service for remote synchronization
public class NavigationSyncService {

    public init() {}

    public func syncNavigationState(request: NavigationStateSyncRequest) async throws -> NavigationStateSyncResponse {
        // TODO: Implement actual sync with backend
        return NavigationStateSyncResponse(
            success: true,
            syncVersion: request.syncVersion + 1,
            conflictsResolved: [],
            timestamp: Date()
        )
    }
}

/// Navigation state manager for local persistence
public class NavigationStateManager {

    public init() {}

    public func saveNavigationState(_ state: NavigationState) async throws {
        // TODO: Implement local persistence (UserDefaults, Core Data, etc.)
    }

    public func loadNavigationState() async throws -> NavigationState? {
        // TODO: Implement local loading
        return nil
    }
}

/// Route validator for navigation validation
public class RouteValidator {

    private let routeRegistry: RouteRegistryManager

    public init(routeRegistry: RouteRegistryManager = RouteRegistryManager()) {
        self.routeRegistry = routeRegistry
    }

    public func validate(routeId: String, parameters: [String: Any]) async throws {
        guard routeRegistry.isValidRoute(routeId) else {
            throw NavigationError.invalidRoute(routeId)
        }

        // TODO: Validate parameters against route schema
    }
}

// MARK: - Error Types

public enum NavigationError: Error, LocalizedError {
    case navigationInProgress
    case cannotGoBack
    case invalidRoute(String)
    case invalidDeepLink(URL)
    case deepLinkProcessingFailed(URL)
    case noNavigationState
    case syncFailed(Error)
    case unknown(Error)

    public var errorDescription: String? {
        switch self {
        case .navigationInProgress:
            return "Navigation already in progress"
        case .cannotGoBack:
            return "Cannot go back from current navigation state"
        case .invalidRoute(let route):
            return "Invalid route: \(route)"
        case .invalidDeepLink(let url):
            return "Invalid deep link: \(url)"
        case .deepLinkProcessingFailed(let url):
            return "Failed to process deep link: \(url)"
        case .noNavigationState:
            return "No navigation state available"
        case .syncFailed(let error):
            return "Navigation sync failed: \(error.localizedDescription)"
        case .unknown(let error):
            return "Unknown navigation error: \(error.localizedDescription)"
        }
    }
}