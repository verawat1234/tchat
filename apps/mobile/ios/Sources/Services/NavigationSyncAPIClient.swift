//
//  NavigationSyncAPIClient.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
import Combine
import Security

/// Real implementation of Navigation Sync API Client
@MainActor
public class NavigationSyncAPIClient: ObservableObject {

    // MARK: - Properties

    private let httpClient: HTTPClient
    private let authenticationProvider: AuthenticationProvider
    private let baseURL: URL
    private var cancellables = Set<AnyCancellable>()

    // MARK: - Initialization

    public init(
        httpClient: HTTPClient = HTTPClient(),
        authenticationProvider: AuthenticationProvider = AuthenticationProvider(),
        baseURL: URL = URL(string: "https://api.tchat.app")!
    ) {
        self.httpClient = httpClient
        self.authenticationProvider = authenticationProvider
        self.baseURL = baseURL
    }

    // MARK: - Navigation Routes API

    /// GET /navigation/routes
    public func getNavigationRoutes(platform: String, userId: String) async throws -> NavigationRoutesResponse {
        let endpoint = baseURL.appendingPathComponent("navigation/routes")

        var urlComponents = URLComponents(url: endpoint, resolvingAgainstBaseURL: false)!
        urlComponents.queryItems = [
            URLQueryItem(name: "platform", value: platform),
            URLQueryItem(name: "userId", value: userId)
        ]

        guard let url = urlComponents.url else {
            throw NavigationSyncError(code: "INVALID_URL", message: "Failed to construct URL")
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Accept")

        // Add authentication
        if let token = await authenticationProvider.getAccessToken() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        } else {
            throw NavigationSyncError(code: "AUTHENTICATION_REQUIRED", message: "No valid authentication token available")
        }

        do {
            let (data, response) = try await httpClient.performRequest(request)

            guard let httpResponse = response as? HTTPURLResponse else {
                throw NavigationSyncError(code: "INVALID_RESPONSE", message: "Invalid response type")
            }

            switch httpResponse.statusCode {
            case 200:
                let routesResponse = try JSONDecoder().decode(NavigationRoutesResponse.self, from: data)
                return routesResponse
            case 400:
                throw NavigationSyncError(code: "INVALID_PLATFORM", message: "Invalid platform specified")
            case 401:
                // Try to refresh token and retry once
                if let refreshedToken = try? await authenticationProvider.refreshToken() {
                    request.setValue("Bearer \(refreshedToken)", forHTTPHeaderField: "Authorization")
                    let (retryData, retryResponse) = try await httpClient.performRequest(request)
                    if let retryHttpResponse = retryResponse as? HTTPURLResponse,
                       retryHttpResponse.statusCode == 200 {
                        let routesResponse = try JSONDecoder().decode(NavigationRoutesResponse.self, from: retryData)
                        return routesResponse
                    }
                }
                throw NavigationSyncError(code: "UNAUTHORIZED", message: "Authentication required")
            case 404:
                throw NavigationSyncError(code: "NOT_FOUND", message: "Routes not found")
            default:
                throw NavigationSyncError(code: "HTTP_ERROR", message: "HTTP \(httpResponse.statusCode)")
            }

        } catch let error as NavigationSyncError {
            throw error
        } catch {
            throw NavigationSyncError(code: "NETWORK_ERROR", message: error.localizedDescription)
        }
    }

    // MARK: - Navigation State Sync API

    /// POST /navigation/state/sync
    public func syncNavigationState(request: NavigationStateSyncRequest) async throws -> NavigationSyncResponse {
        let endpoint = baseURL.appendingPathComponent("navigation/state/sync")

        var urlRequest = URLRequest(url: endpoint)
        urlRequest.httpMethod = "POST"
        urlRequest.setValue("application/json", forHTTPHeaderField: "Content-Type")
        urlRequest.setValue("application/json", forHTTPHeaderField: "Accept")

        // Add authentication
        if let token = await authenticationProvider.getAccessToken() {
            urlRequest.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        } else {
            throw NavigationSyncError(code: "AUTHENTICATION_REQUIRED", message: "No valid authentication token available")
        }

        do {
            let requestData = try JSONEncoder().encode(request)
            urlRequest.httpBody = requestData

            let (data, response) = try await httpClient.performRequest(urlRequest)

            guard let httpResponse = response as? HTTPURLResponse else {
                throw NavigationSyncError(code: "INVALID_RESPONSE", message: "Invalid response type")
            }

            switch httpResponse.statusCode {
            case 200:
                let syncResponse = try JSONDecoder().decode(NavigationSyncResponse.self, from: data)
                return syncResponse
            case 409:
                let conflictData = try JSONDecoder().decode(NavigationSyncConflictResponse.self, from: data)
                throw NavigationSyncConflictError(
                    conflictType: conflictData.conflictType,
                    clientVersion: conflictData.clientVersion,
                    serverVersion: conflictData.serverVersion
                )
            case 400:
                throw NavigationSyncError(code: "INVALID_REQUEST", message: "Invalid sync request")
            case 401:
                // Try to refresh token and retry once
                if let refreshedToken = try? await authenticationProvider.refreshToken() {
                    urlRequest.setValue("Bearer \(refreshedToken)", forHTTPHeaderField: "Authorization")
                    let (retryData, retryResponse) = try await httpClient.performRequest(urlRequest)
                    if let retryHttpResponse = retryResponse as? HTTPURLResponse,
                       retryHttpResponse.statusCode == 200 {
                        let syncResponse = try JSONDecoder().decode(NavigationSyncResponse.self, from: retryData)
                        return syncResponse
                    }
                }
                throw NavigationSyncError(code: "UNAUTHORIZED", message: "Authentication required")
            default:
                throw NavigationSyncError(code: "HTTP_ERROR", message: "HTTP \(httpResponse.statusCode)")
            }

        } catch let error as NavigationSyncError {
            throw error
        } catch let error as NavigationSyncConflictError {
            throw error
        } catch {
            throw NavigationSyncError(code: "NETWORK_ERROR", message: error.localizedDescription)
        }
    }

    // MARK: - Deep Link Resolution API

    /// POST /navigation/deeplink/resolve
    public func resolveDeepLink(request: DeepLinkResolutionRequest) async throws -> DeepLinkResolution {
        let endpoint = baseURL.appendingPathComponent("navigation/deeplink/resolve")

        var urlRequest = URLRequest(url: endpoint)
        urlRequest.httpMethod = "POST"
        urlRequest.setValue("application/json", forHTTPHeaderField: "Content-Type")
        urlRequest.setValue("application/json", forHTTPHeaderField: "Accept")

        // Add authentication
        if let token = await authenticationProvider.getAccessToken() {
            urlRequest.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        } else {
            throw NavigationSyncError(code: "AUTHENTICATION_REQUIRED", message: "No valid authentication token available")
        }

        do {
            let requestData = try JSONEncoder().encode(request)
            urlRequest.httpBody = requestData

            let (data, response) = try await httpClient.performRequest(urlRequest)

            guard let httpResponse = response as? HTTPURLResponse else {
                throw NavigationSyncError(code: "INVALID_RESPONSE", message: "Invalid response type")
            }

            switch httpResponse.statusCode {
            case 200:
                let resolution = try JSONDecoder().decode(DeepLinkResolution.self, from: data)
                return resolution
            case 404:
                throw NavigationSyncError(code: "DEEP_LINK_NOT_FOUND", message: "Deep link not found")
            case 400:
                throw NavigationSyncError(code: "INVALID_REQUEST", message: "Invalid deep link request")
            case 401:
                // Try to refresh token and retry once
                if let refreshedToken = try? await authenticationProvider.refreshToken() {
                    urlRequest.setValue("Bearer \(refreshedToken)", forHTTPHeaderField: "Authorization")
                    let (retryData, retryResponse) = try await httpClient.performRequest(urlRequest)
                    if let retryHttpResponse = retryResponse as? HTTPURLResponse,
                       retryHttpResponse.statusCode == 200 {
                        let resolution = try JSONDecoder().decode(DeepLinkResolution.self, from: retryData)
                        return resolution
                    }
                }
                throw NavigationSyncError(code: "UNAUTHORIZED", message: "Authentication required")
            default:
                throw NavigationSyncError(code: "HTTP_ERROR", message: "HTTP \(httpResponse.statusCode)")
            }

        } catch let error as NavigationSyncError {
            throw error
        } catch {
            throw NavigationSyncError(code: "NETWORK_ERROR", message: error.localizedDescription)
        }
    }
}

// MARK: - Supporting Services

/// HTTP Client for API requests
public class HTTPClient {

    public init() {}

    public func performRequest(_ request: URLRequest) async throws -> (Data, URLResponse) {
        return try await URLSession.shared.data(for: request)
    }
}

/// Secure authentication provider with Keychain integration
public class AuthenticationProvider {

    // MARK: - Properties

    private let keychainService = "com.tchat.app.tokens"
    private let accessTokenKey = "access_token"
    private let refreshTokenKey = "refresh_token"
    private let tokenExpiryKey = "token_expiry"
    private let authBaseURL: URL

    private var isRefreshing = false
    private var refreshTask: Task<String, Error>?

    // MARK: - Initialization

    public init(authBaseURL: URL = URL(string: "https://auth.tchat.app")!) {
        self.authBaseURL = authBaseURL
    }

    // MARK: - Public Methods

    /// Get a valid access token, refreshing if necessary
    public func getAccessToken() async -> String? {
        // Check if we have a valid token
        if let token = getStoredAccessToken(),
           !isTokenExpired() {
            return token
        }

        // Try to refresh if token is expired or missing
        do {
            let newToken = try await refreshTokenIfNeeded()
            return newToken
        } catch {
            print("Failed to refresh token: \(error)")
            return nil
        }
    }

    /// Refresh the access token using the stored refresh token
    public func refreshToken() async throws -> String {
        // Prevent concurrent refresh attempts
        if isRefreshing, let existingTask = refreshTask {
            return try await existingTask.value
        }

        isRefreshing = true
        refreshTask = Task {
            defer {
                isRefreshing = false
                refreshTask = nil
            }

            guard let refreshToken = getStoredRefreshToken() else {
                throw AuthenticationError.noRefreshToken
            }

            let endpoint = authBaseURL.appendingPathComponent("auth/refresh")
            var request = URLRequest(url: endpoint)
            request.httpMethod = "POST"
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")

            let requestBody = RefreshTokenRequest(refreshToken: refreshToken)
            request.httpBody = try JSONEncoder().encode(requestBody)

            let (data, response) = try await URLSession.shared.data(for: request)

            guard let httpResponse = response as? HTTPURLResponse else {
                throw AuthenticationError.invalidResponse
            }

            switch httpResponse.statusCode {
            case 200:
                let tokenResponse = try JSONDecoder().decode(TokenResponse.self, from: data)
                try await storeTokens(
                    accessToken: tokenResponse.accessToken,
                    refreshToken: tokenResponse.refreshToken,
                    expiresIn: tokenResponse.expiresIn
                )
                return tokenResponse.accessToken

            case 401:
                // Refresh token is invalid, clear all tokens
                clearStoredTokens()
                throw AuthenticationError.refreshTokenExpired

            default:
                throw AuthenticationError.refreshFailed
            }
        }

        return try await refreshTask!.value
    }

    /// Store new authentication tokens securely
    public func storeTokens(accessToken: String, refreshToken: String, expiresIn: Int) async throws {
        let expiryDate = Date().addingTimeInterval(TimeInterval(expiresIn))

        try storeInKeychain(key: accessTokenKey, value: accessToken)
        try storeInKeychain(key: refreshTokenKey, value: refreshToken)
        try storeInKeychain(key: tokenExpiryKey, value: ISO8601DateFormatter().string(from: expiryDate))
    }

    /// Clear all stored tokens
    public func clearStoredTokens() {
        deleteFromKeychain(key: accessTokenKey)
        deleteFromKeychain(key: refreshTokenKey)
        deleteFromKeychain(key: tokenExpiryKey)
    }

    /// Check if user is authenticated (has valid tokens)
    public func isAuthenticated() -> Bool {
        return getStoredAccessToken() != nil && getStoredRefreshToken() != nil
    }

    // MARK: - Private Methods

    private func refreshTokenIfNeeded() async throws -> String {
        if !isTokenExpired(), let token = getStoredAccessToken() {
            return token
        }

        return try await refreshToken()
    }

    private func getStoredAccessToken() -> String? {
        return getFromKeychain(key: accessTokenKey)
    }

    private func getStoredRefreshToken() -> String? {
        return getFromKeychain(key: refreshTokenKey)
    }

    private func isTokenExpired() -> Bool {
        guard let expiryString = getFromKeychain(key: tokenExpiryKey),
              let expiryDate = ISO8601DateFormatter().date(from: expiryString) else {
            return true
        }

        // Consider token expired if it expires within the next 5 minutes
        return Date().addingTimeInterval(300) >= expiryDate
    }

    // MARK: - Keychain Operations

    private func storeInKeychain(key: String, value: String) throws {
        let data = value.data(using: .utf8)!

        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: keychainService,
            kSecAttrAccount as String: key,
            kSecValueData as String: data,
            kSecAttrAccessible as String: kSecAttrAccessibleWhenUnlockedThisDeviceOnly
        ]

        // Delete existing item first
        SecItemDelete(query as CFDictionary)

        let status = SecItemAdd(query as CFDictionary, nil)

        guard status == errSecSuccess else {
            throw AuthenticationError.keychainStoreFailed
        }
    }

    private func getFromKeychain(key: String) -> String? {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: keychainService,
            kSecAttrAccount as String: key,
            kSecReturnData as String: true,
            kSecMatchLimit as String: kSecMatchLimitOne
        ]

        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)

        guard status == errSecSuccess,
              let data = result as? Data,
              let string = String(data: data, encoding: .utf8) else {
            return nil
        }

        return string
    }

    private func deleteFromKeychain(key: String) {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: keychainService,
            kSecAttrAccount as String: key
        ]

        SecItemDelete(query as CFDictionary)
    }
}

// MARK: - Token Models

private struct RefreshTokenRequest: Codable {
    let refreshToken: String
}

private struct TokenResponse: Codable {
    let accessToken: String
    let refreshToken: String
    let expiresIn: Int
    let tokenType: String

    private enum CodingKeys: String, CodingKey {
        case accessToken = "access_token"
        case refreshToken = "refresh_token"
        case expiresIn = "expires_in"
        case tokenType = "token_type"
    }
}

// MARK: - Response Models

public struct NavigationSyncConflictResponse: Codable {
    public let conflictType: String
    public let clientVersion: Int
    public let serverVersion: Int
    public let conflictDetails: [String: Any]?

    private enum CodingKeys: String, CodingKey {
        case conflictType, clientVersion, serverVersion
    }

    public init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        conflictType = try container.decode(String.self, forKey: .conflictType)
        clientVersion = try container.decode(Int.self, forKey: .clientVersion)
        serverVersion = try container.decode(Int.self, forKey: .serverVersion)
        conflictDetails = nil // Simplified for demo
    }

    public func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(conflictType, forKey: .conflictType)
        try container.encode(clientVersion, forKey: .clientVersion)
        try container.encode(serverVersion, forKey: .serverVersion)
    }
}

// MARK: - Error Types

public enum AuthenticationError: Error {
    case refreshFailed
    case tokenExpired
    case invalidCredentials
    case noRefreshToken
    case refreshTokenExpired
    case keychainStoreFailed
    case invalidResponse

    var localizedDescription: String {
        switch self {
        case .refreshFailed:
            return "Failed to refresh authentication token"
        case .tokenExpired:
            return "Authentication token has expired"
        case .invalidCredentials:
            return "Invalid authentication credentials"
        case .noRefreshToken:
            return "No refresh token available"
        case .refreshTokenExpired:
            return "Refresh token has expired, please log in again"
        case .keychainStoreFailed:
            return "Failed to store token in Keychain"
        case .invalidResponse:
            return "Invalid response from authentication server"
        }
    }
}

// MARK: - API Client Factory

/// Factory for creating API clients with proper configuration
public class NavigationSyncAPIClientFactory {

    public static func create(environment: Environment = .development) -> NavigationSyncAPIClient {
        let baseURL: URL
        let authBaseURL: URL

        switch environment {
        case .development:
            baseURL = URL(string: "https://dev-api.tchat.app")!
            authBaseURL = URL(string: "https://dev-auth.tchat.app")!
        case .staging:
            baseURL = URL(string: "https://staging-api.tchat.app")!
            authBaseURL = URL(string: "https://staging-auth.tchat.app")!
        case .production:
            baseURL = URL(string: "https://api.tchat.app")!
            authBaseURL = URL(string: "https://auth.tchat.app")!
        }

        return NavigationSyncAPIClient(
            httpClient: HTTPClient(),
            authenticationProvider: AuthenticationProvider(authBaseURL: authBaseURL),
            baseURL: baseURL
        )
    }
}

public enum Environment {
    case development
    case staging
    case production
}