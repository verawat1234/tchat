//
//  PlatformAdapterAPIClient.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
import Combine

/// Real implementation of Platform Adapter API Client
@MainActor
public class PlatformAdapterAPIClient: ObservableObject {

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

    // MARK: - Platform Capabilities API

    /// GET /platform/capabilities
    public func getPlatformCapabilities(platform: String, version: String) async throws -> PlatformCapabilitiesAPIResponse {
        let endpoint = baseURL.appendingPathComponent("platform/capabilities")

        var urlComponents = URLComponents(url: endpoint, resolvingAgainstBaseURL: false)!
        urlComponents.queryItems = [
            URLQueryItem(name: "platform", value: platform),
            URLQueryItem(name: "version", value: version)
        ]

        guard let url = urlComponents.url else {
            throw PlatformAdapterAPIError(code: "INVALID_URL", message: "Failed to construct URL")
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
                throw PlatformAdapterAPIError(code: "INVALID_RESPONSE", message: "Invalid response type")
            }

            switch httpResponse.statusCode {
            case 200:
                let capabilitiesResponse = try JSONDecoder().decode(PlatformCapabilitiesAPIResponse.self, from: data)
                return capabilitiesResponse
            case 400:
                throw PlatformAdapterAPIError(code: "INVALID_PLATFORM", message: "Invalid platform or version")
            case 401:
                throw PlatformAdapterAPIError(code: "UNAUTHORIZED", message: "Authentication required")
            case 404:
                throw PlatformAdapterAPIError(code: "PLATFORM_NOT_SUPPORTED", message: "Platform not supported")
            default:
                throw PlatformAdapterAPIError(code: "HTTP_ERROR", message: "HTTP \(httpResponse.statusCode)")
            }

        } catch let error as PlatformAdapterAPIError {
            throw error
        } catch {
            throw PlatformAdapterAPIError(code: "NETWORK_ERROR", message: error.localizedDescription)
        }
    }

    // MARK: - Gesture Handling API

    /// POST /platform/gestures/handle
    public func handleGesture(request: GestureHandlingAPIRequest) async throws -> GestureHandlingAPIResponse {
        let endpoint = baseURL.appendingPathComponent("platform/gestures/handle")

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
                throw PlatformAdapterAPIError(code: "INVALID_RESPONSE", message: "Invalid response type")
            }

            switch httpResponse.statusCode {
            case 200:
                let gestureResponse = try JSONDecoder().decode(GestureHandlingAPIResponse.self, from: data)
                return gestureResponse
            case 400:
                throw PlatformAdapterAPIError(code: "INVALID_GESTURE", message: "Invalid gesture request")
            case 401:
                throw PlatformAdapterAPIError(code: "UNAUTHORIZED", message: "Authentication required")
            case 422:
                throw PlatformAdapterAPIError(code: "GESTURE_NOT_SUPPORTED", message: "Gesture not supported on platform")
            default:
                throw PlatformAdapterAPIError(code: "HTTP_ERROR", message: "HTTP \(httpResponse.statusCode)")
            }

        } catch let error as PlatformAdapterAPIError {
            throw error
        } catch {
            throw PlatformAdapterAPIError(code: "NETWORK_ERROR", message: error.localizedDescription)
        }
    }

    // MARK: - Animation Execution API

    /// POST /platform/animations/execute
    public func executeAnimation(request: AnimationExecutionAPIRequest) async throws -> AnimationExecutionAPIResponse {
        let endpoint = baseURL.appendingPathComponent("platform/animations/execute")

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
                throw PlatformAdapterAPIError(code: "INVALID_RESPONSE", message: "Invalid response type")
            }

            switch httpResponse.statusCode {
            case 200:
                let animationResponse = try JSONDecoder().decode(AnimationExecutionAPIResponse.self, from: data)
                return animationResponse
            case 400:
                throw PlatformAdapterAPIError(code: "INVALID_ANIMATION", message: "Invalid animation request")
            case 401:
                throw PlatformAdapterAPIError(code: "UNAUTHORIZED", message: "Authentication required")
            case 422:
                throw PlatformAdapterAPIError(code: "ANIMATION_NOT_SUPPORTED", message: "Animation not supported on platform")
            default:
                throw PlatformAdapterAPIError(code: "HTTP_ERROR", message: "HTTP \(httpResponse.statusCode)")
            }

        } catch let error as PlatformAdapterAPIError {
            throw error
        } catch {
            throw PlatformAdapterAPIError(code: "NETWORK_ERROR", message: error.localizedDescription)
        }
    }

    // MARK: - UI Conventions API

    /// GET /platform/ui-conventions
    public func getUIConventions(platform: String) async throws -> UIConventionsAPIResponse {
        let endpoint = baseURL.appendingPathComponent("platform/ui-conventions")

        var urlComponents = URLComponents(url: endpoint, resolvingAgainstBaseURL: false)!
        urlComponents.queryItems = [
            URLQueryItem(name: "platform", value: platform)
        ]

        guard let url = urlComponents.url else {
            throw PlatformAdapterAPIError(code: "INVALID_URL", message: "Failed to construct URL")
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
                throw PlatformAdapterAPIError(code: "INVALID_RESPONSE", message: "Invalid response type")
            }

            switch httpResponse.statusCode {
            case 200:
                let conventionsResponse = try JSONDecoder().decode(UIConventionsAPIResponse.self, from: data)
                return conventionsResponse
            case 400:
                throw PlatformAdapterAPIError(code: "INVALID_PLATFORM", message: "Invalid platform specified")
            case 401:
                throw PlatformAdapterAPIError(code: "UNAUTHORIZED", message: "Authentication required")
            case 404:
                throw PlatformAdapterAPIError(code: "CONVENTIONS_NOT_FOUND", message: "UI conventions not found for platform")
            default:
                throw PlatformAdapterAPIError(code: "HTTP_ERROR", message: "HTTP \(httpResponse.statusCode)")
            }

        } catch let error as PlatformAdapterAPIError {
            throw error
        } catch {
            throw PlatformAdapterAPIError(code: "NETWORK_ERROR", message: error.localizedDescription)
        }
    }

    // MARK: - Device Features API

    /// POST /platform/device/features
    public func reportDeviceFeatures(request: DeviceFeaturesReportRequest) async throws -> DeviceFeaturesReportResponse {
        let endpoint = baseURL.appendingPathComponent("platform/device/features")

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
                throw PlatformAdapterAPIError(code: "INVALID_RESPONSE", message: "Invalid response type")
            }

            switch httpResponse.statusCode {
            case 200:
                let featuresResponse = try JSONDecoder().decode(DeviceFeaturesReportResponse.self, from: data)
                return featuresResponse
            case 400:
                throw PlatformAdapterAPIError(code: "INVALID_DEVICE_DATA", message: "Invalid device features data")
            case 401:
                throw PlatformAdapterAPIError(code: "UNAUTHORIZED", message: "Authentication required")
            default:
                throw PlatformAdapterAPIError(code: "HTTP_ERROR", message: "HTTP \(httpResponse.statusCode)")
            }

        } catch let error as PlatformAdapterAPIError {
            throw error
        } catch {
            throw PlatformAdapterAPIError(code: "NETWORK_ERROR", message: error.localizedDescription)
        }
    }

    // MARK: - Performance Metrics API

    /// POST /platform/performance/metrics
    public func reportPerformanceMetrics(request: PerformanceMetricsRequest) async throws -> PerformanceMetricsResponse {
        let endpoint = baseURL.appendingPathComponent("platform/performance/metrics")

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
                throw PlatformAdapterAPIError(code: "INVALID_RESPONSE", message: "Invalid response type")
            }

            switch httpResponse.statusCode {
            case 200, 202: // 200 = processed, 202 = accepted for async processing
                let metricsResponse = try JSONDecoder().decode(PerformanceMetricsResponse.self, from: data)
                return metricsResponse
            case 400:
                throw PlatformAdapterAPIError(code: "INVALID_METRICS", message: "Invalid performance metrics data")
            case 401:
                throw PlatformAdapterAPIError(code: "UNAUTHORIZED", message: "Authentication required")
            case 413:
                throw PlatformAdapterAPIError(code: "METRICS_TOO_LARGE", message: "Performance metrics payload too large")
            default:
                throw PlatformAdapterAPIError(code: "HTTP_ERROR", message: "HTTP \(httpResponse.statusCode)")
            }

        } catch let error as PlatformAdapterAPIError {
            throw error
        } catch {
            throw PlatformAdapterAPIError(code: "NETWORK_ERROR", message: error.localizedDescription)
        }
    }
}

// MARK: - Request Models

public struct GestureHandlingAPIRequest: Codable {
    public let gestureType: String
    public let direction: String
    public let velocity: Double
    public let position: GesturePosition
    public let platform: String
    public let componentId: String
    public let metadata: [String: AnyCodable]
    public let timestamp: Date
    public let sessionId: String

    public init(
        gestureType: String,
        direction: String,
        velocity: Double,
        position: GesturePosition,
        platform: String,
        componentId: String,
        metadata: [String: AnyCodable],
        timestamp: Date,
        sessionId: String
    ) {
        self.gestureType = gestureType
        self.direction = direction
        self.velocity = velocity
        self.position = position
        self.platform = platform
        self.componentId = componentId
        self.metadata = metadata
        self.timestamp = timestamp
        self.sessionId = sessionId
    }
}

public struct AnimationExecutionAPIRequest: Codable {
    public let animationType: String
    public let duration: Double
    public let easing: String
    public let properties: [String: AnyCodable]
    public let platform: String
    public let componentId: String
    public let timestamp: Date
    public let sessionId: String

    public init(
        animationType: String,
        duration: Double,
        easing: String,
        properties: [String: AnyCodable],
        platform: String,
        componentId: String,
        timestamp: Date,
        sessionId: String
    ) {
        self.animationType = animationType
        self.duration = duration
        self.easing = easing
        self.properties = properties
        self.platform = platform
        self.componentId = componentId
        self.timestamp = timestamp
        self.sessionId = sessionId
    }
}

public struct DeviceFeaturesReportRequest: Codable {
    public let platform: String
    public let platformVersion: String
    public let deviceModel: String
    public let features: [String: Bool]
    public let capabilities: [String: AnyCodable]
    public let limitations: [String]
    public let userId: String
    public let sessionId: String
    public let timestamp: Date

    public init(
        platform: String,
        platformVersion: String,
        deviceModel: String,
        features: [String: Bool],
        capabilities: [String: AnyCodable],
        limitations: [String],
        userId: String,
        sessionId: String,
        timestamp: Date
    ) {
        self.platform = platform
        self.platformVersion = platformVersion
        self.deviceModel = deviceModel
        self.features = features
        self.capabilities = capabilities
        self.limitations = limitations
        self.userId = userId
        self.sessionId = sessionId
        self.timestamp = timestamp
    }
}

public struct PerformanceMetricsRequest: Codable {
    public let platform: String
    public let metrics: PerformanceMetrics
    public let userId: String
    public let sessionId: String
    public let timestamp: Date

    public init(platform: String, metrics: PerformanceMetrics, userId: String, sessionId: String, timestamp: Date) {
        self.platform = platform
        self.metrics = metrics
        self.userId = userId
        self.sessionId = sessionId
        self.timestamp = timestamp
    }
}

public struct PerformanceMetrics: Codable {
    public let appLaunchTime: Double
    public let navigationTime: Double
    public let renderTime: Double
    public let memoryUsage: Double
    public let cpuUsage: Double
    public let networkLatency: Double
    public let frameRate: Double
    public let customMetrics: [String: Double]

    public init(
        appLaunchTime: Double,
        navigationTime: Double,
        renderTime: Double,
        memoryUsage: Double,
        cpuUsage: Double,
        networkLatency: Double,
        frameRate: Double,
        customMetrics: [String: Double]
    ) {
        self.appLaunchTime = appLaunchTime
        self.navigationTime = navigationTime
        self.renderTime = renderTime
        self.memoryUsage = memoryUsage
        self.cpuUsage = cpuUsage
        self.networkLatency = networkLatency
        self.frameRate = frameRate
        self.customMetrics = customMetrics
    }
}

// MARK: - Response Models

public struct PlatformCapabilitiesAPIResponse: Codable {
    public let platform: String
    public let version: String
    public let capabilities: [PlatformCapability]
    public let limitations: [String]
    public let supportedGestures: [String]
    public let supportedAnimations: [String]
    public let performanceTargets: PerformanceTargets

    public init(
        platform: String,
        version: String,
        capabilities: [PlatformCapability],
        limitations: [String],
        supportedGestures: [String],
        supportedAnimations: [String],
        performanceTargets: PerformanceTargets
    ) {
        self.platform = platform
        self.version = version
        self.capabilities = capabilities
        self.limitations = limitations
        self.supportedGestures = supportedGestures
        self.supportedAnimations = supportedAnimations
        self.performanceTargets = performanceTargets
    }
}

public struct PerformanceTargets: Codable {
    public let maxAppLaunchTime: Double
    public let maxNavigationTime: Double
    public let minFrameRate: Double
    public let maxMemoryUsage: Double

    public init(maxAppLaunchTime: Double, maxNavigationTime: Double, minFrameRate: Double, maxMemoryUsage: Double) {
        self.maxAppLaunchTime = maxAppLaunchTime
        self.maxNavigationTime = maxNavigationTime
        self.minFrameRate = minFrameRate
        self.maxMemoryUsage = maxMemoryUsage
    }
}

public struct GestureHandlingAPIResponse: Codable {
    public let handled: Bool
    public let action: String?
    public let gestureType: String
    public let timestamp: Date
    public let preventDefaultBehavior: Bool
    public let platformOptimizations: [String: AnyCodable]

    public init(
        handled: Bool,
        action: String?,
        gestureType: String,
        timestamp: Date,
        preventDefaultBehavior: Bool,
        platformOptimizations: [String: AnyCodable]
    ) {
        self.handled = handled
        self.action = action
        self.gestureType = gestureType
        self.timestamp = timestamp
        self.preventDefaultBehavior = preventDefaultBehavior
        self.platformOptimizations = platformOptimizations
    }
}

public struct AnimationExecutionAPIResponse: Codable {
    public let started: Bool
    public let animationType: String
    public let duration: Double
    public let animationId: String
    public let timestamp: Date
    public let platformOptimizations: [String: AnyCodable]

    public init(
        started: Bool,
        animationType: String,
        duration: Double,
        animationId: String,
        timestamp: Date,
        platformOptimizations: [String: AnyCodable]
    ) {
        self.started = started
        self.animationType = animationType
        self.duration = duration
        self.animationId = animationId
        self.timestamp = timestamp
        self.platformOptimizations = platformOptimizations
    }
}

public struct UIConventionsAPIResponse: Codable {
    public let platform: String
    public let designSystem: DesignSystem
    public let navigationPatterns: [String: String]
    public let gestureConventions: [String: AnyCodable]
    public let animationSpecs: [String: AnyCodable]
    public let accessibilityGuidelines: [String: String]

    public init(
        platform: String,
        designSystem: DesignSystem,
        navigationPatterns: [String: String],
        gestureConventions: [String: AnyCodable],
        animationSpecs: [String: AnyCodable],
        accessibilityGuidelines: [String: String]
    ) {
        self.platform = platform
        self.designSystem = designSystem
        self.navigationPatterns = navigationPatterns
        self.gestureConventions = gestureConventions
        self.animationSpecs = animationSpecs
        self.accessibilityGuidelines = accessibilityGuidelines
    }
}

public struct DeviceFeaturesReportResponse: Codable {
    public let acknowledged: Bool
    public let optimizations: [String: AnyCodable]
    public let recommendations: [String]
    public let timestamp: Date

    public init(acknowledged: Bool, optimizations: [String: AnyCodable], recommendations: [String], timestamp: Date) {
        self.acknowledged = acknowledged
        self.optimizations = optimizations
        self.recommendations = recommendations
        self.timestamp = timestamp
    }
}

public struct PerformanceMetricsResponse: Codable {
    public let received: Bool
    public let processed: Bool
    public let insights: [String]
    public let optimizationSuggestions: [String]
    public let timestamp: Date

    public init(received: Bool, processed: Bool, insights: [String], optimizationSuggestions: [String], timestamp: Date) {
        self.received = received
        self.processed = processed
        self.insights = insights
        self.optimizationSuggestions = optimizationSuggestions
        self.timestamp = timestamp
    }
}

// MARK: - Error Types

public struct PlatformAdapterAPIError: Error, LocalizedError {
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

// MARK: - API Client Factory

/// Factory for creating Platform Adapter API clients with proper configuration
public class PlatformAdapterAPIClientFactory {

    public static func create(environment: Environment = .development) -> PlatformAdapterAPIClient {
        let baseURL: URL

        switch environment {
        case .development:
            baseURL = URL(string: "https://dev-api.tchat.app")!
        case .staging:
            baseURL = URL(string: "https://staging-api.tchat.app")!
        case .production:
            baseURL = URL(string: "https://api.tchat.app")!
        }

        return PlatformAdapterAPIClient(
            httpClient: HTTPClient(),
            authenticationProvider: AuthenticationProvider(),
            baseURL: baseURL
        )
    }
}