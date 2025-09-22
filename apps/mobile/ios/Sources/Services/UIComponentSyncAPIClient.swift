//
//  UIComponentSyncAPIClient.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
import Combine

/// Real implementation of UI Component Sync API Client
@MainActor
public class UIComponentSyncAPIClient: ObservableObject {

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

    // MARK: - Component State Sync API

    /// POST /components/state/sync
    public func syncComponentState(request: ComponentStateSyncRequest) async throws -> ComponentStateSyncResponse {
        let endpoint = baseURL.appendingPathComponent("components/state/sync")

        var urlRequest = URLRequest(url: endpoint)
        urlRequest.httpMethod = "POST"
        urlRequest.setValue("application/json", forHTTPHeaderField: "Content-Type")
        urlRequest.setValue("application/json", forHTTPHeaderField: "Accept")

        // Add authentication
        if let token = await authenticationProvider.getAccessToken() {
            urlRequest.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        do {
            let requestData = try JSONEncoder().encode(request)
            urlRequest.httpBody = requestData

            let (data, response) = try await httpClient.performRequest(urlRequest)

            guard let httpResponse = response as? HTTPURLResponse else {
                throw UIComponentSyncError(code: "INVALID_RESPONSE", message: "Invalid response type")
            }

            switch httpResponse.statusCode {
            case 200:
                let syncResponse = try JSONDecoder().decode(ComponentStateSyncResponse.self, from: data)
                return syncResponse
            case 409:
                let conflictData = try JSONDecoder().decode(ComponentSyncConflictResponse.self, from: data)
                throw UIComponentSyncConflictError(
                    conflictType: conflictData.conflictType,
                    componentId: conflictData.componentId,
                    clientVersion: conflictData.clientVersion,
                    serverVersion: conflictData.serverVersion
                )
            case 400:
                throw UIComponentSyncError(code: "INVALID_REQUEST", message: "Invalid sync request")
            case 401:
                throw UIComponentSyncError(code: "UNAUTHORIZED", message: "Authentication required")
            case 422:
                throw UIComponentSyncError(code: "VALIDATION_ERROR", message: "Component state validation failed")
            default:
                throw UIComponentSyncError(code: "HTTP_ERROR", message: "HTTP \(httpResponse.statusCode)")
            }

        } catch let error as UIComponentSyncError {
            throw error
        } catch let error as UIComponentSyncConflictError {
            throw error
        } catch {
            throw UIComponentSyncError(code: "NETWORK_ERROR", message: error.localizedDescription)
        }
    }

    // MARK: - Component State Retrieval API

    /// GET /components/state/{instanceId}
    public func getComponentState(instanceId: String, userId: String) async throws -> ComponentStateResponse {
        let endpoint = baseURL.appendingPathComponent("components/state/\(instanceId)")

        var urlComponents = URLComponents(url: endpoint, resolvingAgainstBaseURL: false)!
        urlComponents.queryItems = [
            URLQueryItem(name: "userId", value: userId)
        ]

        guard let url = urlComponents.url else {
            throw UIComponentSyncError(code: "INVALID_URL", message: "Failed to construct URL")
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Accept")

        // Add authentication
        if let token = await authenticationProvider.getAccessToken() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        do {
            let (data, response) = try await httpClient.performRequest(request)

            guard let httpResponse = response as? HTTPURLResponse else {
                throw UIComponentSyncError(code: "INVALID_RESPONSE", message: "Invalid response type")
            }

            switch httpResponse.statusCode {
            case 200:
                let stateResponse = try JSONDecoder().decode(ComponentStateResponse.self, from: data)
                return stateResponse
            case 404:
                throw UIComponentSyncError(code: "COMPONENT_NOT_FOUND", message: "Component state not found")
            case 401:
                throw UIComponentSyncError(code: "UNAUTHORIZED", message: "Authentication required")
            case 403:
                throw UIComponentSyncError(code: "ACCESS_DENIED", message: "Access denied to component state")
            default:
                throw UIComponentSyncError(code: "HTTP_ERROR", message: "HTTP \(httpResponse.statusCode)")
            }

        } catch let error as UIComponentSyncError {
            throw error
        } catch {
            throw UIComponentSyncError(code: "NETWORK_ERROR", message: error.localizedDescription)
        }
    }

    // MARK: - Bulk Component State Sync API

    /// POST /components/state/bulk-sync
    public func bulkSyncComponentStates(request: BulkComponentStateSyncRequest) async throws -> BulkComponentStateSyncResponse {
        let endpoint = baseURL.appendingPathComponent("components/state/bulk-sync")

        var urlRequest = URLRequest(url: endpoint)
        urlRequest.httpMethod = "POST"
        urlRequest.setValue("application/json", forHTTPHeaderField: "Content-Type")
        urlRequest.setValue("application/json", forHTTPHeaderField: "Accept")

        // Add authentication
        if let token = await authenticationProvider.getAccessToken() {
            urlRequest.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        do {
            let requestData = try JSONEncoder().encode(request)
            urlRequest.httpBody = requestData

            let (data, response) = try await httpClient.performRequest(urlRequest)

            guard let httpResponse = response as? HTTPURLResponse else {
                throw UIComponentSyncError(code: "INVALID_RESPONSE", message: "Invalid response type")
            }

            switch httpResponse.statusCode {
            case 200:
                let bulkResponse = try JSONDecoder().decode(BulkComponentStateSyncResponse.self, from: data)
                return bulkResponse
            case 207:
                // Multi-status response with partial success
                let bulkResponse = try JSONDecoder().decode(BulkComponentStateSyncResponse.self, from: data)
                return bulkResponse
            case 400:
                throw UIComponentSyncError(code: "INVALID_REQUEST", message: "Invalid bulk sync request")
            case 401:
                throw UIComponentSyncError(code: "UNAUTHORIZED", message: "Authentication required")
            case 413:
                throw UIComponentSyncError(code: "PAYLOAD_TOO_LARGE", message: "Bulk request too large")
            default:
                throw UIComponentSyncError(code: "HTTP_ERROR", message: "HTTP \(httpResponse.statusCode)")
            }

        } catch let error as UIComponentSyncError {
            throw error
        } catch {
            throw UIComponentSyncError(code: "NETWORK_ERROR", message: error.localizedDescription)
        }
    }

    // MARK: - Component Schema Validation API

    /// POST /components/schema/validate
    public func validateComponentSchema(request: ComponentSchemaValidationRequest) async throws -> ComponentSchemaValidationResponse {
        let endpoint = baseURL.appendingPathComponent("components/schema/validate")

        var urlRequest = URLRequest(url: endpoint)
        urlRequest.httpMethod = "POST"
        urlRequest.setValue("application/json", forHTTPHeaderField: "Content-Type")
        urlRequest.setValue("application/json", forHTTPHeaderField: "Accept")

        // Add authentication
        if let token = await authenticationProvider.getAccessToken() {
            urlRequest.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        do {
            let requestData = try JSONEncoder().encode(request)
            urlRequest.httpBody = requestData

            let (data, response) = try await httpClient.performRequest(urlRequest)

            guard let httpResponse = response as? HTTPURLResponse else {
                throw UIComponentSyncError(code: "INVALID_RESPONSE", message: "Invalid response type")
            }

            switch httpResponse.statusCode {
            case 200:
                let validationResponse = try JSONDecoder().decode(ComponentSchemaValidationResponse.self, from: data)
                return validationResponse
            case 400:
                throw UIComponentSyncError(code: "INVALID_SCHEMA", message: "Invalid component schema")
            case 401:
                throw UIComponentSyncError(code: "UNAUTHORIZED", message: "Authentication required")
            case 422:
                let validationResponse = try JSONDecoder().decode(ComponentSchemaValidationResponse.self, from: data)
                return validationResponse // Return with validation errors
            default:
                throw UIComponentSyncError(code: "HTTP_ERROR", message: "HTTP \(httpResponse.statusCode)")
            }

        } catch let error as UIComponentSyncError {
            throw error
        } catch {
            throw UIComponentSyncError(code: "NETWORK_ERROR", message: error.localizedDescription)
        }
    }
}

// MARK: - Response Models

public struct ComponentSyncConflictResponse: Codable {
    public let conflictType: String
    public let componentId: String
    public let clientVersion: Int
    public let serverVersion: Int
    public let conflictDetails: [String: AnyCodable]?

    private enum CodingKeys: String, CodingKey {
        case conflictType, componentId, clientVersion, serverVersion
    }

    public init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        conflictType = try container.decode(String.self, forKey: .conflictType)
        componentId = try container.decode(String.self, forKey: .componentId)
        clientVersion = try container.decode(Int.self, forKey: .clientVersion)
        serverVersion = try container.decode(Int.self, forKey: .serverVersion)
        conflictDetails = nil // Simplified for demo
    }

    public func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(conflictType, forKey: .conflictType)
        try container.encode(componentId, forKey: .componentId)
        try container.encode(clientVersion, forKey: .clientVersion)
        try container.encode(serverVersion, forKey: .serverVersion)
    }
}

public struct ComponentStateResponse: Codable {
    public let componentState: ComponentState
    public let version: Int
    public let lastModified: Date
    public let permissions: [String]

    public init(componentState: ComponentState, version: Int, lastModified: Date, permissions: [String]) {
        self.componentState = componentState
        self.version = version
        self.lastModified = lastModified
        self.permissions = permissions
    }
}

public struct BulkComponentStateSyncRequest: Codable {
    public let userId: String
    public let sessionId: String
    public let platform: String
    public let componentStates: [ComponentState]
    public let timestamp: Date
    public let syncVersion: Int
    public let options: BulkSyncOptions?

    public init(
        userId: String,
        sessionId: String,
        platform: String,
        componentStates: [ComponentState],
        timestamp: Date,
        syncVersion: Int,
        options: BulkSyncOptions? = nil
    ) {
        self.userId = userId
        self.sessionId = sessionId
        self.platform = platform
        self.componentStates = componentStates
        self.timestamp = timestamp
        self.syncVersion = syncVersion
        self.options = options
    }
}

public struct BulkSyncOptions: Codable {
    public let maxConcurrency: Int?
    public let timeoutSeconds: Int?
    public let stopOnFirstError: Bool?
    public let includeMetrics: Bool?

    public init(maxConcurrency: Int? = nil, timeoutSeconds: Int? = nil, stopOnFirstError: Bool? = nil, includeMetrics: Bool? = nil) {
        self.maxConcurrency = maxConcurrency
        self.timeoutSeconds = timeoutSeconds
        self.stopOnFirstError = stopOnFirstError
        self.includeMetrics = includeMetrics
    }
}

public struct BulkComponentStateSyncResponse: Codable {
    public let success: Bool
    public let totalRequested: Int
    public let totalSucceeded: Int
    public let totalFailed: Int
    public let results: [ComponentSyncResult]
    public let syncVersion: Int
    public let timestamp: Date
    public let metrics: BulkSyncMetrics?

    public init(
        success: Bool,
        totalRequested: Int,
        totalSucceeded: Int,
        totalFailed: Int,
        results: [ComponentSyncResult],
        syncVersion: Int,
        timestamp: Date,
        metrics: BulkSyncMetrics? = nil
    ) {
        self.success = success
        self.totalRequested = totalRequested
        self.totalSucceeded = totalSucceeded
        self.totalFailed = totalFailed
        self.results = results
        self.syncVersion = syncVersion
        self.timestamp = timestamp
        self.metrics = metrics
    }
}

public struct ComponentSyncResult: Codable {
    public let instanceId: String
    public let success: Bool
    public let error: String?
    public let version: Int?

    public init(instanceId: String, success: Bool, error: String? = nil, version: Int? = nil) {
        self.instanceId = instanceId
        self.success = success
        self.error = error
        self.version = version
    }
}

public struct BulkSyncMetrics: Codable {
    public let totalDurationMs: Int
    public let avgComponentSyncMs: Double
    public let networkLatencyMs: Int
    public let processingTimeMs: Int

    public init(totalDurationMs: Int, avgComponentSyncMs: Double, networkLatencyMs: Int, processingTimeMs: Int) {
        self.totalDurationMs = totalDurationMs
        self.avgComponentSyncMs = avgComponentSyncMs
        self.networkLatencyMs = networkLatencyMs
        self.processingTimeMs = processingTimeMs
    }
}

public struct ComponentSchemaValidationRequest: Codable {
    public let componentId: String
    public let schema: [String: AnyCodable]
    public let platform: String
    public let version: String

    public init(componentId: String, schema: [String: AnyCodable], platform: String, version: String) {
        self.componentId = componentId
        self.schema = schema
        self.platform = platform
        self.version = version
    }
}

public struct ComponentSchemaValidationResponse: Codable {
    public let valid: Bool
    public let errors: [ValidationError]
    public let warnings: [ValidationWarning]
    public let suggestions: [String]

    public init(valid: Bool, errors: [ValidationError], warnings: [ValidationWarning], suggestions: [String]) {
        self.valid = valid
        self.errors = errors
        self.warnings = warnings
        self.suggestions = suggestions
    }
}

public struct ValidationError: Codable {
    public let field: String
    public let message: String
    public let code: String

    public init(field: String, message: String, code: String) {
        self.field = field
        self.message = message
        self.code = code
    }
}

public struct ValidationWarning: Codable {
    public let field: String
    public let message: String
    public let severity: String

    public init(field: String, message: String, severity: String) {
        self.field = field
        self.message = message
        self.severity = severity
    }
}

// MARK: - Error Types

public struct UIComponentSyncError: Error, LocalizedError {
    public let code: String
    public let message: String

    public init(code: String, message: String) {
        self.code = code
        self.message = message
    }

    public var errorDescription: String? {
        return "\(code): \(message)"
    }
}

public struct UIComponentSyncConflictError: Error, LocalizedError {
    public let conflictType: String
    public let componentId: String
    public let clientVersion: Int
    public let serverVersion: Int

    public init(conflictType: String, componentId: String, clientVersion: Int, serverVersion: Int) {
        self.conflictType = conflictType
        self.componentId = componentId
        self.clientVersion = clientVersion
        self.serverVersion = serverVersion
    }

    public var errorDescription: String? {
        return "Component sync conflict (\(conflictType)) for \(componentId): client v\(clientVersion) vs server v\(serverVersion)"
    }
}

// Helper for encoding/decoding arbitrary JSON values (internal to avoid collision)
internal struct UIComponentAnyCodable: Codable {
    public let value: Any

    public init<T>(_ value: T?) {
        self.value = value ?? ()
    }
}

extension UIComponentAnyCodable: ExpressibleByNilLiteral {
    public init(nilLiteral _: ()) {
        self.value = ()
    }
}

extension AnyCodable {
    public init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()

        if container.decodeNil() {
            self.init(nil)
        } else if let bool = try? container.decode(Bool.self) {
            self.init(bool)
        } else if let int = try? container.decode(Int.self) {
            self.init(int)
        } else if let double = try? container.decode(Double.self) {
            self.init(double)
        } else if let string = try? container.decode(String.self) {
            self.init(string)
        } else if let array = try? container.decode([AnyCodable].self) {
            self.init(array.map { $0.value })
        } else if let dictionary = try? container.decode([String: AnyCodable].self) {
            self.init(dictionary.mapValues { $0.value })
        } else {
            throw DecodingError.dataCorruptedError(in: container, debugDescription: "AnyCodable value cannot be decoded")
        }
    }

    public func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()

        switch value {
        case is Void:
            try container.encodeNil()
        case let bool as Bool:
            try container.encode(bool)
        case let int as Int:
            try container.encode(int)
        case let double as Double:
            try container.encode(double)
        case let string as String:
            try container.encode(string)
        case let array as [Any]:
            try container.encode(array.map { AnyCodable($0) })
        case let dictionary as [String: Any]:
            try container.encode(dictionary.mapValues { AnyCodable($0) })
        default:
            let context = EncodingError.Context(codingPath: container.codingPath, debugDescription: "AnyCodable value cannot be encoded")
            throw EncodingError.invalidValue(value, context)
        }
    }
}

// MARK: - API Client Factory

/// Factory for creating UI Component Sync API clients with proper configuration
public class UIComponentSyncAPIClientFactory {

    public static func create(environment: Environment = .development) -> UIComponentSyncAPIClient {
        let baseURL: URL

        switch environment {
        case .development:
            baseURL = URL(string: "https://dev-api.tchat.app")!
        case .staging:
            baseURL = URL(string: "https://staging-api.tchat.app")!
        case .production:
            baseURL = URL(string: "https://api.tchat.app")!
        }

        return UIComponentSyncAPIClient(
            httpClient: HTTPClient(),
            authenticationProvider: AuthenticationProvider(),
            baseURL: baseURL
        )
    }
}