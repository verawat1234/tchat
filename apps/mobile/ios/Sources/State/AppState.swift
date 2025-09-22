//
//  AppState.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import SwiftUI
import Combine

/// Global application state management following Tchat design system
@MainActor
public class AppState: ObservableObject {

    // MARK: - Singleton
    public static let shared = AppState()

    // MARK: - Published Properties

    /// Current user information
    @Published public var currentUser: UserModel?

    /// Authentication state
    @Published public var isAuthenticated: Bool = false

    /// Theme preferences
    @Published public var themePreferences: ThemePreferences = ThemePreferences()

    /// Chat state
    @Published public var chatState: ChatState = ChatState()

    /// Store state
    @Published public var storeState: StoreState = StoreState()

    /// Social state
    @Published public var socialState: SocialState = SocialState()

    /// Video state
    @Published public var videoState: VideoState = VideoState()

    /// Network connectivity
    @Published public var isOnline: Bool = true

    /// Sync status
    @Published public var syncStatus: SyncStatus = .idle

    /// Error state
    @Published public var lastError: AppError?

    // MARK: - Private Properties
    private var cancellables = Set<AnyCancellable>()
    private let syncManager: StateSyncManager
    private let persistence: PersistenceManager

    // MARK: - Initialization
    private init() {
        self.syncManager = StateSyncManager()
        self.persistence = PersistenceManager()

        setupObservers()
        loadPersistedState()
        startSyncTimer()
    }

    // MARK: - Public Methods

    /// Update user information
    public func updateUser(_ user: UserModel) {
        currentUser = user
        isAuthenticated = true
        syncWithServer()
    }

    /// Sign out user
    public func signOut() {
        currentUser = nil
        isAuthenticated = false
        clearAllState()
        syncWithServer()
    }

    /// Update theme preferences
    public func updateThemePreferences(_ preferences: ThemePreferences) {
        themePreferences = preferences
        persistence.save(preferences, forKey: "themePreferences")
        syncWithServer()
    }

    /// Force sync with server
    public func forceSyncWithServer() {
        syncStatus = .syncing
        syncWithServer()
    }

    /// Handle deep link
    public func handleDeepLink(_ url: URL) {
        // Deep link routing logic
        print("Handling deep link: \(url)")
    }

    // MARK: - Private Methods

    private func setupObservers() {
        // Monitor authentication state changes
        $isAuthenticated
            .dropFirst()
            .sink { [weak self] isAuth in
                self?.handleAuthenticationChange(isAuth)
            }
            .store(in: &cancellables)

        // Monitor network connectivity
        NotificationCenter.default.publisher(for: .networkConnectivityChanged)
            .sink { [weak self] notification in
                if let isConnected = notification.object as? Bool {
                    self?.isOnline = isConnected
                    if isConnected {
                        self?.syncWithServer()
                    }
                }
            }
            .store(in: &cancellables)

        // Monitor app lifecycle
        NotificationCenter.default.publisher(for: UIApplication.willResignActiveNotification)
            .sink { [weak self] _ in
                self?.saveState()
            }
            .store(in: &cancellables)

        NotificationCenter.default.publisher(for: UIApplication.didBecomeActiveNotification)
            .sink { [weak self] _ in
                self?.syncWithServer()
            }
            .store(in: &cancellables)
    }

    private func loadPersistedState() {
        // Load theme preferences
        if let savedPreferences: ThemePreferences = persistence.load(forKey: "themePreferences") {
            themePreferences = savedPreferences
        }

        // Load authentication state
        isAuthenticated = persistence.loadBool(forKey: "isAuthenticated")

        // Load user data
        if let savedUser: UserModel = persistence.load(forKey: "currentUser") {
            currentUser = savedUser
        }
    }

    private func saveState() {
        persistence.save(themePreferences, forKey: "themePreferences")
        persistence.save(isAuthenticated, forKey: "isAuthenticated")

        if let user = currentUser {
            persistence.save(user, forKey: "currentUser")
        }
    }

    private func handleAuthenticationChange(_ isAuthenticated: Bool) {
        if isAuthenticated {
            startSyncTimer()
        } else {
            stopSyncTimer()
            clearAllState()
        }
    }

    private func clearAllState() {
        chatState = ChatState()
        storeState = StoreState()
        socialState = SocialState()
        videoState = VideoState()
        persistence.clearAll()
    }

    private var syncTimer: Timer?

    private func startSyncTimer() {
        syncTimer?.invalidate()
        syncTimer = Timer.scheduledTimer(withTimeInterval: 30.0, repeats: true) { [weak self] _ in
            self?.syncWithServer()
        }
    }

    private func stopSyncTimer() {
        syncTimer?.invalidate()
        syncTimer = nil
    }

    private func syncWithServer() {
        guard isOnline && isAuthenticated else { return }

        syncStatus = .syncing

        Task {
            do {
                try await syncManager.syncState(self)
                await MainActor.run {
                    syncStatus = .success
                }
            } catch {
                await MainActor.run {
                    syncStatus = .failed
                    lastError = AppError.syncFailed(error.localizedDescription)
                }
            }
        }
    }
}

// MARK: - State Models

/// User model
public struct UserModel: Codable, Identifiable {
    public let id: String
    public let username: String
    public let email: String
    public let displayName: String
    public let avatarURL: String?
    public let preferences: UserPreferences

    public init(id: String, username: String, email: String, displayName: String, avatarURL: String? = nil, preferences: UserPreferences = UserPreferences()) {
        self.id = id
        self.username = username
        self.email = email
        self.displayName = displayName
        self.avatarURL = avatarURL
        self.preferences = preferences
    }
}

/// User preferences
public struct UserPreferences: Codable {
    public let language: String
    public let notificationsEnabled: Bool
    public let soundEnabled: Bool

    public init(language: String = "en", notificationsEnabled: Bool = true, soundEnabled: Bool = true) {
        self.language = language
        self.notificationsEnabled = notificationsEnabled
        self.soundEnabled = soundEnabled
    }
}

/// Theme preferences
public struct ThemePreferences: Codable {
    public let isDarkMode: Bool
    public let accentColor: String
    public let fontSize: FontSize

    public init(isDarkMode: Bool = false, accentColor: String = "#3B82F6", fontSize: FontSize = .medium) {
        self.isDarkMode = isDarkMode
        self.accentColor = accentColor
        self.fontSize = fontSize
    }
}

public enum FontSize: String, Codable, CaseIterable {
    case small = "small"
    case medium = "medium"
    case large = "large"
    case extraLarge = "extraLarge"
}

/// Chat state
public struct ChatState: Codable {
    public let unreadCount: Int
    public let lastMessageTimestamp: Date?
    public let activeConversations: [String]

    public init(unreadCount: Int = 0, lastMessageTimestamp: Date? = nil, activeConversations: [String] = []) {
        self.unreadCount = unreadCount
        self.lastMessageTimestamp = lastMessageTimestamp
        self.activeConversations = activeConversations
    }
}

/// Store state
public struct StoreState: Codable {
    public let cartItemCount: Int
    public let wishlistCount: Int
    public let lastPurchaseDate: Date?

    public init(cartItemCount: Int = 0, wishlistCount: Int = 0, lastPurchaseDate: Date? = nil) {
        self.cartItemCount = cartItemCount
        self.wishlistCount = wishlistCount
        self.lastPurchaseDate = lastPurchaseDate
    }
}

/// Social state
public struct SocialState: Codable {
    public let friendRequestCount: Int
    public let notificationCount: Int
    public let lastActivityDate: Date?

    public init(friendRequestCount: Int = 0, notificationCount: Int = 0, lastActivityDate: Date? = nil) {
        self.friendRequestCount = friendRequestCount
        self.notificationCount = notificationCount
        self.lastActivityDate = lastActivityDate
    }
}

/// Video state
public struct VideoState: Codable {
    public let watchHistoryCount: Int
    public let subscriptionCount: Int
    public let lastWatchedDate: Date?

    public init(watchHistoryCount: Int = 0, subscriptionCount: Int = 0, lastWatchedDate: Date? = nil) {
        self.watchHistoryCount = watchHistoryCount
        self.subscriptionCount = subscriptionCount
        self.lastWatchedDate = lastWatchedDate
    }
}

/// Sync status
public enum SyncStatus {
    case idle
    case syncing
    case success
    case failed
}

/// Application errors
public enum AppError: Error, LocalizedError {
    case syncFailed(String)
    case authenticationFailed(String)
    case networkError(String)
    case dataCorruption(String)

    public var errorDescription: String? {
        switch self {
        case .syncFailed(let message):
            return "Sync failed: \(message)"
        case .authenticationFailed(let message):
            return "Authentication failed: \(message)"
        case .networkError(let message):
            return "Network error: \(message)"
        case .dataCorruption(let message):
            return "Data corruption: \(message)"
        }
    }
}

// MARK: - Notifications
extension Notification.Name {
    static let networkConnectivityChanged = Notification.Name("networkConnectivityChanged")
}