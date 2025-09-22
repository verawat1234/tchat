//
//  SessionManager.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
import Combine

/// Manages cross-platform session state and workspace switching
public class SessionManager: ObservableObject {

    // MARK: - Properties
    private let stateSyncManager: StateSyncManager
    private let persistenceManager: PersistenceManager
    private var cancellables = Set<AnyCancellable>()

    @Published public var currentSession: UserSession?
    @Published public var activeWorkspace: Workspace?
    @Published public var isSessionSyncing: Bool = false
    @Published public var lastSessionSync: Date?

    // Session timeout configuration
    private let sessionTimeout: TimeInterval = 30 * 60 // 30 minutes
    private var sessionTimer: Timer?

    // MARK: - Initialization
    public init(stateSyncManager: StateSyncManager, persistenceManager: PersistenceManager) {
        self.stateSyncManager = stateSyncManager
        self.persistenceManager = persistenceManager

        loadPersistedSession()
        setupSessionMonitoring()
    }

    deinit {
        sessionTimer?.invalidate()
    }

    // MARK: - Session Management

    /// Start a new session with user authentication
    public func startSession(for user: User, workspace: Workspace) async throws {
        let session = UserSession(
            id: UUID().uuidString,
            userId: user.id,
            workspaceId: workspace.id,
            platform: "ios",
            deviceId: getDeviceId(),
            startTime: Date(),
            lastActivity: Date(),
            isActive: true
        )

        currentSession = session
        activeWorkspace = workspace

        // Sync session to server
        try await syncSessionToServer(session)

        // Start session monitoring
        startSessionMonitoring()

        // Persist session locally
        persistenceManager.saveSession(session)
        persistenceManager.saveWorkspace(workspace)

        await MainActor.run {
            NotificationCenter.default.post(
                name: .sessionDidStart,
                object: session
            )
        }
    }

    /// Switch to a different workspace
    public func switchWorkspace(_ workspace: Workspace) async throws {
        guard let session = currentSession else {
            throw SessionError.noActiveSession
        }

        // Update session with new workspace
        let updatedSession = UserSession(
            id: session.id,
            userId: session.userId,
            workspaceId: workspace.id,
            platform: session.platform,
            deviceId: session.deviceId,
            startTime: session.startTime,
            lastActivity: Date(),
            isActive: true
        )

        currentSession = updatedSession
        activeWorkspace = workspace

        // Sync workspace change to server
        try await syncWorkspaceSwitch(updatedSession, workspace)

        // Persist changes
        persistenceManager.saveSession(updatedSession)
        persistenceManager.saveWorkspace(workspace)

        await MainActor.run {
            NotificationCenter.default.post(
                name: .workspaceDidChange,
                object: workspace,
                userInfo: ["previousWorkspace": activeWorkspace as Any]
            )
        }
    }

    /// Update session activity
    public func updateActivity() async {
        guard let session = currentSession else { return }

        let updatedSession = UserSession(
            id: session.id,
            userId: session.userId,
            workspaceId: session.workspaceId,
            platform: session.platform,
            deviceId: session.deviceId,
            startTime: session.startTime,
            lastActivity: Date(),
            isActive: true
        )

        currentSession = updatedSession
        persistenceManager.saveSession(updatedSession)

        // Sync activity update to server (throttled)
        Task {
            try await throttledSessionSync(updatedSession)
        }
    }

    /// End current session
    public func endSession() async throws {
        guard let session = currentSession else {
            throw SessionError.noActiveSession
        }

        let endedSession = UserSession(
            id: session.id,
            userId: session.userId,
            workspaceId: session.workspaceId,
            platform: session.platform,
            deviceId: session.deviceId,
            startTime: session.startTime,
            lastActivity: Date(),
            isActive: false
        )

        // Sync session end to server
        try await syncSessionToServer(endedSession)

        // Clear local session
        currentSession = nil
        activeWorkspace = nil
        stopSessionMonitoring()

        // Clear persisted session
        persistenceManager.clearSession()

        await MainActor.run {
            NotificationCenter.default.post(
                name: .sessionDidEnd,
                object: endedSession
            )
        }
    }

    /// Get active sessions across all platforms
    public func getActiveSessions() async throws -> [CrossPlatformSession] {
        guard let url = URL(string: "https://api.tchat.app/api/sessions/active") else {
            throw SessionError.invalidURL
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
            throw SessionError.serverError
        }

        return try JSONDecoder().decode([CrossPlatformSession].self, from: data)
    }

    // MARK: - Private Methods

    private func loadPersistedSession() {
        if let session = persistenceManager.loadSession(),
           let workspace = persistenceManager.loadWorkspace() {

            // Check if session is still valid (not expired)
            let timeSinceLastActivity = Date().timeIntervalSince(session.lastActivity)
            if timeSinceLastActivity < sessionTimeout {
                currentSession = session
                activeWorkspace = workspace
                startSessionMonitoring()
            } else {
                // Session expired, clear it
                persistenceManager.clearSession()
            }
        }
    }

    private func setupSessionMonitoring() {
        // Monitor app lifecycle events
        NotificationCenter.default.publisher(for: UIApplication.willEnterForegroundNotification)
            .sink { [weak self] _ in
                Task {
                    await self?.updateActivity()
                }
            }
            .store(in: &cancellables)

        NotificationCenter.default.publisher(for: UIApplication.didEnterBackgroundNotification)
            .sink { [weak self] _ in
                Task {
                    await self?.updateActivity()
                }
            }
            .store(in: &cancellables)
    }

    private func startSessionMonitoring() {
        sessionTimer?.invalidate()
        sessionTimer = Timer.scheduledTimer(withTimeInterval: 60.0, repeats: true) { [weak self] _ in
            Task {
                await self?.checkSessionExpiry()
            }
        }
    }

    private func stopSessionMonitoring() {
        sessionTimer?.invalidate()
        sessionTimer = nil
    }

    private func checkSessionExpiry() async {
        guard let session = currentSession else { return }

        let timeSinceLastActivity = Date().timeIntervalSince(session.lastActivity)
        if timeSinceLastActivity >= sessionTimeout {
            try? await endSession()
        }
    }

    private func syncSessionToServer(_ session: UserSession) async throws {
        guard let url = URL(string: "https://api.tchat.app/api/sessions/sync") else {
            throw SessionError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        if let token = getAuthToken() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        let encoder = JSONEncoder()
        encoder.dateEncodingStrategy = .iso8601
        request.httpBody = try encoder.encode(session)

        let (_, response) = try await URLSession.shared.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse,
              httpResponse.statusCode == 200 else {
            throw SessionError.syncFailed
        }

        lastSessionSync = Date()
    }

    private func syncWorkspaceSwitch(_ session: UserSession, _ workspace: Workspace) async throws {
        let payload = WorkspaceSwitchPayload(
            sessionId: session.id,
            userId: session.userId,
            newWorkspaceId: workspace.id,
            timestamp: Date(),
            platform: "ios"
        )

        guard let url = URL(string: "https://api.tchat.app/api/sessions/workspace-switch") else {
            throw SessionError.invalidURL
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
            throw SessionError.syncFailed
        }
    }

    private var lastSyncTime: Date = Date.distantPast
    private func throttledSessionSync(_ session: UserSession) async throws {
        let now = Date()
        let timeSinceLastSync = now.timeIntervalSince(lastSyncTime)

        // Only sync if more than 30 seconds have passed
        if timeSinceLastSync >= 30 {
            try await syncSessionToServer(session)
            lastSyncTime = now
        }
    }

    private func getDeviceId() -> String {
        return UIDevice.current.identifierForVendor?.uuidString ?? UUID().uuidString
    }

    private func getAuthToken() -> String? {
        // Retrieve authentication token from keychain or secure storage
        return nil // Placeholder
    }
}

// MARK: - Session Models

public struct UserSession: Codable, Identifiable {
    public let id: String
    public let userId: String
    public let workspaceId: String
    public let platform: String
    public let deviceId: String
    public let startTime: Date
    public let lastActivity: Date
    public let isActive: Bool
}

public struct Workspace: Codable, Identifiable {
    public let id: String
    public let name: String
    public let description: String?
    public let iconUrl: String?
    public let memberCount: Int
    public let isPersonal: Bool
}

public struct CrossPlatformSession: Codable {
    public let sessionId: String
    public let platform: String
    public let deviceType: String
    public let lastActivity: Date
    public let isCurrentSession: Bool
}

public struct WorkspaceSwitchPayload: Codable {
    public let sessionId: String
    public let userId: String
    public let newWorkspaceId: String
    public let timestamp: Date
    public let platform: String
}

public enum SessionError: Error, LocalizedError {
    case noActiveSession
    case invalidURL
    case serverError
    case syncFailed
    case sessionExpired

    public var errorDescription: String? {
        switch self {
        case .noActiveSession:
            return "No active session found"
        case .invalidURL:
            return "Invalid session management URL"
        case .serverError:
            return "Session server error"
        case .syncFailed:
            return "Failed to sync session"
        case .sessionExpired:
            return "Session has expired"
        }
    }
}

// MARK: - Session Notification Extensions

extension Notification.Name {
    static let sessionDidStart = Notification.Name("sessionDidStart")
    static let sessionDidEnd = Notification.Name("sessionDidEnd")
    static let workspaceDidChange = Notification.Name("workspaceDidChange")
}