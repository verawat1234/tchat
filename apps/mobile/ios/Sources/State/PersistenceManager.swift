//
//  PersistenceManager.swift
//  TchatApp
//
//  Created by Claude on 21/09/2024.
//

import Foundation

/// Manages local data persistence for the Tchat app
public class PersistenceManager {

    // MARK: - Properties
    private let userDefaults = UserDefaults.standard
    private let keychain = KeychainManager()

    // MARK: - UserDefaults Keys
    private enum Keys {
        static let themePreferences = "themePreferences"
        static let isAuthenticated = "isAuthenticated"
        static let currentUser = "currentUser"
        static let chatState = "chatState"
        static let storeState = "storeState"
        static let socialState = "socialState"
        static let videoState = "videoState"
        static let lastSyncTimestamp = "lastSyncTimestamp"
    }

    // MARK: - Public Methods

    /// Save codable object to UserDefaults
    public func save<T: Codable>(_ object: T, forKey key: String) {
        do {
            let data = try JSONEncoder().encode(object)
            userDefaults.set(data, forKey: key)
        } catch {
            print("Failed to save \(key): \(error)")
        }
    }

    /// Load codable object from UserDefaults
    public func load<T: Codable>(forKey key: String, type: T.Type = T.self) -> T? {
        guard let data = userDefaults.data(forKey: key) else { return nil }

        do {
            return try JSONDecoder().decode(T.self, from: data)
        } catch {
            print("Failed to load \(key): \(error)")
            return nil
        }
    }

    /// Save string to UserDefaults
    public func save(_ string: String, forKey key: String) {
        userDefaults.set(string, forKey: key)
    }

    /// Load string from UserDefaults
    public func loadString(forKey key: String) -> String? {
        return userDefaults.string(forKey: key)
    }

    /// Save boolean to UserDefaults
    public func save(_ bool: Bool, forKey key: String) {
        userDefaults.set(bool, forKey: key)
    }

    /// Load boolean from UserDefaults
    public func loadBool(forKey key: String) -> Bool {
        return userDefaults.bool(forKey: key)
    }

    /// Save integer to UserDefaults
    public func save(_ int: Int, forKey key: String) {
        userDefaults.set(int, forKey: key)
    }

    /// Load integer from UserDefaults
    public func loadInt(forKey key: String) -> Int {
        return userDefaults.integer(forKey: key)
    }

    /// Save date to UserDefaults
    public func save(_ date: Date, forKey key: String) {
        userDefaults.set(date, forKey: key)
    }

    /// Load date from UserDefaults
    public func loadDate(forKey key: String) -> Date? {
        return userDefaults.object(forKey: key) as? Date
    }

    /// Remove object for key
    public func remove(forKey key: String) {
        userDefaults.removeObject(forKey: key)
    }

    /// Clear all stored data
    public func clearAll() {
        let keys = [
            Keys.themePreferences,
            Keys.isAuthenticated,
            Keys.currentUser,
            Keys.chatState,
            Keys.storeState,
            Keys.socialState,
            Keys.videoState,
            Keys.lastSyncTimestamp
        ]

        keys.forEach { userDefaults.removeObject(forKey: $0) }
        keychain.clearAll()
    }

    // MARK: - Secure Storage (Keychain)

    /// Save sensitive data to keychain
    public func saveSecure(_ data: Data, forKey key: String) -> Bool {
        return keychain.save(data, forKey: key)
    }

    /// Load sensitive data from keychain
    public func loadSecure(forKey key: String) -> Data? {
        return keychain.load(forKey: key)
    }

    /// Save secure string to keychain
    public func saveSecureString(_ string: String, forKey key: String) -> Bool {
        guard let data = string.data(using: .utf8) else { return false }
        return keychain.save(data, forKey: key)
    }

    /// Load secure string from keychain
    public func loadSecureString(forKey key: String) -> String? {
        guard let data = keychain.load(forKey: key) else { return nil }
        return String(data: data, encoding: .utf8)
    }

    /// Remove secure data from keychain
    public func removeSecure(forKey key: String) -> Bool {
        return keychain.remove(forKey: key)
    }
}

// MARK: - Keychain Manager

/// Simple keychain manager for secure storage
public class KeychainManager {

    private let service = "com.tchat.app"

    /// Save data to keychain
    func save(_ data: Data, forKey key: String) -> Bool {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: key,
            kSecValueData as String: data
        ]

        // Delete existing item
        SecItemDelete(query as CFDictionary)

        // Add new item
        let status = SecItemAdd(query as CFDictionary, nil)
        return status == errSecSuccess
    }

    /// Load data from keychain
    func load(forKey key: String) -> Data? {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: key,
            kSecReturnData as String: true,
            kSecMatchLimit as String: kSecMatchLimitOne
        ]

        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)

        guard status == errSecSuccess else { return nil }
        return result as? Data
    }

    /// Remove data from keychain
    func remove(forKey key: String) -> Bool {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: key
        ]

        let status = SecItemDelete(query as CFDictionary)
        return status == errSecSuccess || status == errSecItemNotFound
    }

    /// Clear all keychain data for this service
    func clearAll() {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service
        ]

        SecItemDelete(query as CFDictionary)
    }
}

// MARK: - Convenience Extensions

extension PersistenceManager {

    /// Save theme preferences
    public func saveThemePreferences(_ preferences: ThemePreferences) {
        save(preferences, forKey: Keys.themePreferences)
    }

    /// Load theme preferences
    public func loadThemePreferences() -> ThemePreferences? {
        return load(forKey: Keys.themePreferences, type: ThemePreferences.self)
    }

    /// Save user authentication state
    public func saveAuthenticationState(_ isAuthenticated: Bool) {
        save(isAuthenticated, forKey: Keys.isAuthenticated)
    }

    /// Load user authentication state
    public func loadAuthenticationState() -> Bool {
        return loadBool(forKey: Keys.isAuthenticated)
    }

    /// Save current user
    public func saveCurrentUser(_ user: UserModel) {
        save(user, forKey: Keys.currentUser)
    }

    /// Load current user
    public func loadCurrentUser() -> UserModel? {
        return load(forKey: Keys.currentUser, type: UserModel.self)
    }

    /// Save chat state
    public func saveChatState(_ state: ChatState) {
        save(state, forKey: Keys.chatState)
    }

    /// Load chat state
    public func loadChatState() -> ChatState? {
        return load(forKey: Keys.chatState, type: ChatState.self)
    }

    /// Save store state
    public func saveStoreState(_ state: StoreState) {
        save(state, forKey: Keys.storeState)
    }

    /// Load store state
    public func loadStoreState() -> StoreState? {
        return load(forKey: Keys.storeState, type: StoreState.self)
    }

    /// Save social state
    public func saveSocialState(_ state: SocialState) {
        save(state, forKey: Keys.socialState)
    }

    /// Load social state
    public func loadSocialState() -> SocialState? {
        return load(forKey: Keys.socialState, type: SocialState.self)
    }

    /// Save video state
    public func saveVideoState(_ state: VideoState) {
        save(state, forKey: Keys.videoState)
    }

    /// Load video state
    public func loadVideoState() -> VideoState? {
        return load(forKey: Keys.videoState, type: VideoState.self)
    }

    /// Save last sync timestamp
    public func saveLastSyncTimestamp(_ timestamp: Date) {
        save(timestamp, forKey: Keys.lastSyncTimestamp)
    }

    /// Load last sync timestamp
    public func loadLastSyncTimestamp() -> Date? {
        return loadDate(forKey: Keys.lastSyncTimestamp)
    }
}