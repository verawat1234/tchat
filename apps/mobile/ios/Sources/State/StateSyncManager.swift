//
//  StateSyncManager.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import Foundation
import Network

/// Manages state synchronization between native app and web platform
public class StateSyncManager: ObservableObject {

    // MARK: - Properties
    private let baseURL = "https://api.tchat.app" // Replace with actual API endpoint
    private let networkMonitor = NWPathMonitor()
    private let queue = DispatchQueue(label: "StateSyncManager")

    @Published public var isConnected: Bool = false
    @Published public var lastSyncTimestamp: Date?

    // MARK: - Initialization
    public init() {
        startNetworkMonitoring()
    }

    deinit {
        networkMonitor.cancel()
    }

    // MARK: - Public Methods

    /// Sync app state with server
    public func syncState(_ appState: AppState) async throws {
        guard isConnected else {
            throw SyncError.networkUnavailable
        }

        let statePayload = createStatePayload(from: appState)

        do {
            // Upload current state
            try await uploadState(statePayload)

            // Download latest state
            let serverState = try await downloadState()
            await updateAppState(appState, with: serverState)

            await MainActor.run {
                lastSyncTimestamp = Date()
            }

        } catch {
            throw SyncError.syncFailed(error.localizedDescription)
        }
    }

    /// Download state from server
    public func downloadState() async throws -> ServerState {
        guard let url = URL(string: "\(baseURL)/api/state/sync") else {
            throw SyncError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        // Add authentication header if needed
        if let token = getAuthToken() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        let (data, response) = try await URLSession.shared.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse,
              httpResponse.statusCode == 200 else {
            throw SyncError.serverError
        }

        return try JSONDecoder().decode(ServerState.self, from: data)
    }

    /// Upload state to server
    public func uploadState(_ state: StatePayload) async throws {
        guard let url = URL(string: "\(baseURL)/api/state/sync") else {
            throw SyncError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        // Add authentication header if needed
        if let token = getAuthToken() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        let encoder = JSONEncoder()
        encoder.dateEncodingStrategy = .iso8601
        request.httpBody = try encoder.encode(state)

        let (_, response) = try await URLSession.shared.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse,
              httpResponse.statusCode == 200 else {
            throw SyncError.serverError
        }
    }

    // MARK: - Private Methods

    private func startNetworkMonitoring() {
        networkMonitor.pathUpdateHandler = { [weak self] path in
            DispatchQueue.main.async {
                self?.isConnected = path.status == .satisfied

                // Notify app state about connectivity change
                NotificationCenter.default.post(
                    name: .networkConnectivityChanged,
                    object: path.status == .satisfied
                )
            }
        }
        networkMonitor.start(queue: queue)
    }

    private func createStatePayload(from appState: AppState) -> StatePayload {
        return StatePayload(
            timestamp: Date(),
            platform: "ios",
            version: Bundle.main.infoDictionary?["CFBundleShortVersionString"] as? String ?? "1.0.0",
            userId: appState.currentUser?.id,
            themePreferences: appState.themePreferences,
            chatState: appState.chatState,
            storeState: appState.storeState,
            socialState: appState.socialState,
            videoState: appState.videoState
        )
    }

    @MainActor
    private func updateAppState(_ appState: AppState, with serverState: ServerState) {
        // Only update if server state is newer
        guard let serverTimestamp = serverState.timestamp,
              let lastSync = lastSyncTimestamp,
              serverTimestamp > lastSync else {
            return
        }

        // Update theme preferences if different
        if serverState.themePreferences != appState.themePreferences {
            appState.themePreferences = serverState.themePreferences
        }

        // Update chat state
        if serverState.chatState.unreadCount != appState.chatState.unreadCount {
            appState.chatState = serverState.chatState
        }

        // Update store state
        if serverState.storeState.cartItemCount != appState.storeState.cartItemCount {
            appState.storeState = serverState.storeState
        }

        // Update social state
        if serverState.socialState.friendRequestCount != appState.socialState.friendRequestCount {
            appState.socialState = serverState.socialState
        }

        // Update video state
        if serverState.videoState.watchHistoryCount != appState.videoState.watchHistoryCount {
            appState.videoState = serverState.videoState
        }
    }

    private func getAuthToken() -> String? {
        // Retrieve authentication token from keychain or secure storage
        return nil // Placeholder
    }
}

// MARK: - State Models

/// Payload sent to server
public struct StatePayload: Codable {
    let timestamp: Date
    let platform: String
    let version: String
    let userId: String?
    let themePreferences: ThemePreferences
    let chatState: ChatState
    let storeState: StoreState
    let socialState: SocialState
    let videoState: VideoState
}

/// State received from server
public struct ServerState: Codable {
    let timestamp: Date?
    let themePreferences: ThemePreferences
    let chatState: ChatState
    let storeState: StoreState
    let socialState: SocialState
    let videoState: VideoState
}

/// Sync errors
public enum SyncError: Error, LocalizedError {
    case networkUnavailable
    case invalidURL
    case serverError
    case syncFailed(String)
    case dataCorruption

    public var errorDescription: String? {
        switch self {
        case .networkUnavailable:
            return "Network is unavailable"
        case .invalidURL:
            return "Invalid server URL"
        case .serverError:
            return "Server error occurred"
        case .syncFailed(let message):
            return "Sync failed: \(message)"
        case .dataCorruption:
            return "Data corruption detected"
        }
    }
}

// MARK: - Extensions

extension ThemePreferences: Equatable {
    public static func == (lhs: ThemePreferences, rhs: ThemePreferences) -> Bool {
        return lhs.isDarkMode == rhs.isDarkMode &&
               lhs.accentColor == rhs.accentColor &&
               lhs.fontSize == rhs.fontSize
    }
}