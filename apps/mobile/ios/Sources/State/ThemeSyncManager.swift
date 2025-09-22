//
//  ThemeSyncManager.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
import Combine

/// Manages theme synchronization between web and mobile platforms
public class ThemeSyncManager: ObservableObject {

    // MARK: - Properties
    private let stateSyncManager: StateSyncManager
    private let persistenceManager: PersistenceManager
    private var cancellables = Set<AnyCancellable>()

    @Published public var isThemeSyncing: Bool = false
    @Published public var lastThemeSync: Date?

    // MARK: - Initialization
    public init(stateSyncManager: StateSyncManager, persistenceManager: PersistenceManager) {
        self.stateSyncManager = stateSyncManager
        self.persistenceManager = persistenceManager
        setupThemeChangeListener()
    }

    // MARK: - Public Methods

    /// Sync theme changes to server immediately
    public func syncThemeToServer(_ themePreferences: ThemePreferences) async throws {
        isThemeSyncing = true
        defer { isThemeSyncing = false }

        let themePayload = ThemeSyncPayload(
            timestamp: Date(),
            platform: "ios",
            preferences: themePreferences,
            syncReason: .userChanged
        )

        try await uploadThemeChanges(themePayload)

        await MainActor.run {
            lastThemeSync = Date()
        }
    }

    /// Poll for theme changes from other platforms
    public func pollThemeChanges() async throws -> ThemePreferences? {
        let serverTheme = try await downloadLatestTheme()

        // Check if server theme is newer than our last sync
        guard let serverTimestamp = serverTheme.timestamp,
              let lastSync = lastThemeSync,
              serverTimestamp > lastSync else {
            return nil
        }

        return serverTheme.preferences
    }

    /// Apply theme changes with cross-platform compatibility
    public func applyThemeChanges(_ preferences: ThemePreferences, fromPlatform platform: String) async {
        await MainActor.run {
            // Store theme locally
            persistenceManager.saveThemePreferences(preferences)

            // Notify app components of theme change
            NotificationCenter.default.post(
                name: .themeDidChange,
                object: preferences,
                userInfo: ["source": platform]
            )

            lastThemeSync = Date()
        }
    }

    /// Start automatic theme synchronization
    public func startAutoSync(interval: TimeInterval = 30.0) {
        Timer.publish(every: interval, on: .main, in: .common)
            .autoconnect()
            .sink { [weak self] _ in
                Task {
                    try await self?.checkForThemeUpdates()
                }
            }
            .store(in: &cancellables)
    }

    /// Stop automatic theme synchronization
    public func stopAutoSync() {
        cancellables.removeAll()
    }

    // MARK: - Private Methods

    private func setupThemeChangeListener() {
        NotificationCenter.default.publisher(for: .themePreferencesChanged)
            .compactMap { $0.object as? ThemePreferences }
            .sink { [weak self] preferences in
                Task {
                    try await self?.syncThemeToServer(preferences)
                }
            }
            .store(in: &cancellables)
    }

    private func checkForThemeUpdates() async throws {
        guard !isThemeSyncing else { return }

        if let newTheme = try await pollThemeChanges() {
            await applyThemeChanges(newTheme, fromPlatform: "server")
        }
    }

    private func uploadThemeChanges(_ payload: ThemeSyncPayload) async throws {
        guard let url = URL(string: "https://api.tchat.app/api/theme/sync") else {
            throw ThemeSyncError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        if let token = getAuthToken() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        let encoder = JSONEncoder()
        encoder.dateEncodingStrategy = .iso8601
        request.httpBody = try encoder.encode(payload)

        let (_, response) = try await URLSession.shared.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse,
              httpResponse.statusCode == 200 else {
            throw ThemeSyncError.uploadFailed
        }
    }

    private func downloadLatestTheme() async throws -> ThemeServerResponse {
        guard let url = URL(string: "https://api.tchat.app/api/theme/latest") else {
            throw ThemeSyncError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        if let token = getAuthToken() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        let (data, response) = try await URLSession.shared.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse,
              httpResponse.statusCode == 200 else {
            throw ThemeSyncError.downloadFailed
        }

        return try JSONDecoder().decode(ThemeServerResponse.self, from: data)
    }

    private func getAuthToken() -> String? {
        // Retrieve authentication token from keychain or secure storage
        return nil // Placeholder
    }
}

// MARK: - Theme Sync Models

public struct ThemeSyncPayload: Codable {
    let timestamp: Date
    let platform: String
    let preferences: ThemePreferences
    let syncReason: ThemeSyncReason
}

public struct ThemeServerResponse: Codable {
    let timestamp: Date?
    let preferences: ThemePreferences
    let lastModifiedBy: String?
}

public enum ThemeSyncReason: String, Codable {
    case userChanged
    case systemChanged
    case crossPlatformSync
    case startup
}

public enum ThemeSyncError: Error, LocalizedError {
    case invalidURL
    case uploadFailed
    case downloadFailed
    case networkUnavailable

    public var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Invalid theme sync URL"
        case .uploadFailed:
            return "Failed to upload theme changes"
        case .downloadFailed:
            return "Failed to download theme updates"
        case .networkUnavailable:
            return "Network unavailable for theme sync"
        }
    }
}

// MARK: - Notification Extensions

extension Notification.Name {
    static let themeDidChange = Notification.Name("themeDidChange")
    static let themePreferencesChanged = Notification.Name("themePreferencesChanged")
    static let networkConnectivityChanged = Notification.Name("networkConnectivityChanged")
}