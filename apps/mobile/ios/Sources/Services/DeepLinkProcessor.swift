//
//  DeepLinkProcessor.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
import Combine

/// Deep link processor for handling deep link resolution and routing
@MainActor
public class DeepLinkProcessor: ObservableObject {

    // MARK: - Published Properties

    @Published public var pendingDeepLinks: [DeepLink] = []
    @Published public var processingStatus: ProcessingStatus = .idle
    @Published public var lastProcessedLink: DeepLink?

    // MARK: - Private Properties

    private let routeRegistry: RouteRegistryManager
    private let resolutionService: DeepLinkResolutionService
    private let cacheService: DeepLinkCacheService
    private let analyticsService: DeepLinkAnalyticsService
    private var cancellables = Set<AnyCancellable>()

    // MARK: - Initialization

    public init(
        routeRegistry: RouteRegistryManager = RouteRegistryManager(),
        resolutionService: DeepLinkResolutionService = DeepLinkResolutionService(),
        cacheService: DeepLinkCacheService = DeepLinkCacheService(),
        analyticsService: DeepLinkAnalyticsService = DeepLinkAnalyticsService()
    ) {
        self.routeRegistry = routeRegistry
        self.resolutionService = resolutionService
        self.cacheService = cacheService
        self.analyticsService = analyticsService

        setupProcessing()
    }

    // MARK: - Public Methods

    /// Process a deep link URL
    public func processDeepLink(_ urlString: String, userId: String) async throws -> DeepLinkResolution {
        guard let deepLink = DeepLink.fromURL(urlString) else {
            throw DeepLinkProcessorError.invalidURL(urlString)
        }

        return try await processDeepLink(deepLink, userId: userId)
    }

    /// Process a deep link object
    public func processDeepLink(_ deepLink: DeepLink, userId: String) async throws -> DeepLinkResolution {
        processingStatus = .processing

        do {
            // Check cache first
            if let cachedResolution = await cacheService.getCachedResolution(for: deepLink) {
                processingStatus = .completed
                lastProcessedLink = deepLink
                await analyticsService.trackDeepLink(deepLink, resolution: cachedResolution, fromCache: true)
                return cachedResolution
            }

            // Validate deep link
            try validateDeepLink(deepLink)

            // Create resolution request
            let request = deepLink.createResolutionRequest(userId: userId)

            // Resolve deep link
            let resolution = try await resolutionService.resolveDeepLink(request: request)

            // Cache resolution
            await cacheService.cacheResolution(resolution, for: deepLink)

            // Track analytics
            await analyticsService.trackDeepLink(deepLink, resolution: resolution, fromCache: false)

            processingStatus = .completed
            lastProcessedLink = deepLink

            return resolution

        } catch {
            processingStatus = .failed(error)
            throw error
        }
    }

    /// Process multiple deep links in batch
    public func processDeepLinks(_ deepLinks: [DeepLink], userId: String) async -> [Result<DeepLinkResolution, Error>] {
        var results: [Result<DeepLinkResolution, Error>] = []

        for deepLink in deepLinks {
            do {
                let resolution = try await processDeepLink(deepLink, userId: userId)
                results.append(.success(resolution))
            } catch {
                results.append(.failure(error))
            }
        }

        return results
    }

    /// Queue deep link for processing
    public func queueDeepLink(_ deepLink: DeepLink) {
        pendingDeepLinks.append(deepLink)
    }

    /// Process queued deep links
    public func processQueuedDeepLinks(userId: String) async {
        guard !pendingDeepLinks.isEmpty else { return }

        let linksToProcess = pendingDeepLinks
        pendingDeepLinks.removeAll()

        for deepLink in linksToProcess {
            try? await processDeepLink(deepLink, userId: userId)
        }
    }

    /// Check if URL can be handled
    public func canHandle(_ urlString: String) -> Bool {
        guard let deepLink = DeepLink.fromURL(urlString) else { return false }
        return deepLink.isValid
    }

    /// Extract route information from URL
    public func extractRoute(from urlString: String) -> (routeId: String, parameters: [String: Any])? {
        guard let deepLink = DeepLink.fromURL(urlString), deepLink.isValid else {
            return nil
        }

        // Find matching route pattern
        let availableRoutes = routeRegistry.availableRoutes

        for routeId in availableRoutes {
            if deepLink.matchesRouteId(routeId) {
                let parameters = deepLink.extractParameters()
                return (routeId, parameters)
            }
        }

        return nil
    }

    /// Get deep link statistics
    public func getStatistics() async -> DeepLinkStatistics {
        return await analyticsService.getStatistics()
    }

    /// Clear cache
    public func clearCache() async {
        await cacheService.clearCache()
    }

    // MARK: - Private Methods

    private func setupProcessing() {
        // Auto-process pending deep links when queue is not empty
        $pendingDeepLinks
            .filter { !$0.isEmpty }
            .debounce(for: .milliseconds(500), scheduler: RunLoop.main)
            .sink { [weak self] _ in
                Task { @MainActor in
                    await self?.processQueuedDeepLinks(userId: "current_user") // TODO: Get actual user ID
                }
            }
            .store(in: &cancellables)
    }

    private func validateDeepLink(_ deepLink: DeepLink) throws {
        // Check if expired
        if deepLink.isExpired {
            throw DeepLinkProcessorError.expiredDeepLink(deepLink.url)
        }

        // Check platform support
        if !deepLink.platform.isEmpty && deepLink.platform != "ios" {
            throw DeepLinkProcessorError.unsupportedPlatform(deepLink.platform)
        }

        // Check authentication requirements
        if deepLink.requiresAuthentication {
            // TODO: Check if user is authenticated
        }
    }
}

// MARK: - Supporting Services

/// Deep link resolution service
public class DeepLinkResolutionService {

    public init() {}

    public func resolveDeepLink(request: DeepLinkResolutionRequest) async throws -> DeepLinkResolution {
        // Parse URL components
        guard let url = URL(string: request.url),
              let components = URLComponents(url: url, resolvingAgainstBaseURL: false) else {
            throw DeepLinkProcessorError.invalidURL(request.url)
        }

        let path = components.path
        let queryParams = components.queryItems?.reduce(into: [String: Any]()) { result, item in
            result[item.name] = item.value
        } ?? [:]

        // Extract route ID from path
        let routeId = path.trimmingCharacters(in: CharacterSet(charactersIn: "/"))
            .replacingOccurrences(of: "/", with: "/")

        // TODO: Validate route with backend API
        // For now, return a basic resolution

        return DeepLinkResolution(
            routeId: routeId,
            parameters: queryParams,
            isValid: !routeId.isEmpty,
            fallbackAction: routeId.isEmpty ? "goto_home" : nil,
            requiresAuth: false
        )
    }
}

/// Deep link cache service
public class DeepLinkCacheService {

    private var cache: [String: DeepLinkResolution] = [:]
    private let cacheExpiry: TimeInterval = 300 // 5 minutes

    public init() {}

    public func getCachedResolution(for deepLink: DeepLink) async -> DeepLinkResolution? {
        return cache[deepLink.url]
    }

    public func cacheResolution(_ resolution: DeepLinkResolution, for deepLink: DeepLink) async {
        cache[deepLink.url] = resolution

        // Auto-expire cache entries
        Task {
            try await Task.sleep(nanoseconds: UInt64(cacheExpiry * 1_000_000_000))
            cache.removeValue(forKey: deepLink.url)
        }
    }

    public func clearCache() async {
        cache.removeAll()
    }
}

/// Deep link analytics service
public class DeepLinkAnalyticsService {

    private var statistics = DeepLinkStatistics()

    public init() {}

    public func trackDeepLink(
        _ deepLink: DeepLink,
        resolution: DeepLinkResolution,
        fromCache: Bool
    ) async {
        statistics.totalProcessed += 1

        if resolution.isValid {
            statistics.successfulResolutions += 1
        } else {
            statistics.failedResolutions += 1
        }

        if fromCache {
            statistics.cacheHits += 1
        } else {
            statistics.cacheMisses += 1
        }

        // Track by source
        statistics.sourceBreakdown[deepLink.metadata.source.rawValue, default: 0] += 1
    }

    public func getStatistics() async -> DeepLinkStatistics {
        return statistics
    }
}

// MARK: - Supporting Types

public enum ProcessingStatus: Equatable {
    case idle
    case processing
    case completed
    case failed(Error)

    public static func == (lhs: ProcessingStatus, rhs: ProcessingStatus) -> Bool {
        switch (lhs, rhs) {
        case (.idle, .idle), (.processing, .processing), (.completed, .completed):
            return true
        case (.failed, .failed):
            return true
        default:
            return false
        }
    }
}

public struct DeepLinkStatistics: Codable {
    public var totalProcessed: Int = 0
    public var successfulResolutions: Int = 0
    public var failedResolutions: Int = 0
    public var cacheHits: Int = 0
    public var cacheMisses: Int = 0
    public var sourceBreakdown: [String: Int] = [:]

    public var successRate: Double {
        guard totalProcessed > 0 else { return 0 }
        return Double(successfulResolutions) / Double(totalProcessed)
    }

    public var cacheHitRate: Double {
        let totalRequests = cacheHits + cacheMisses
        guard totalRequests > 0 else { return 0 }
        return Double(cacheHits) / Double(totalRequests)
    }
}

// MARK: - Error Types

public enum DeepLinkProcessorError: Error, LocalizedError {
    case invalidURL(String)
    case expiredDeepLink(String)
    case unsupportedPlatform(String)
    case resolutionFailed(String)
    case authenticationRequired
    case rateLimitExceeded

    public var errorDescription: String? {
        switch self {
        case .invalidURL(let url):
            return "Invalid deep link URL: \(url)"
        case .expiredDeepLink(let url):
            return "Expired deep link: \(url)"
        case .unsupportedPlatform(let platform):
            return "Unsupported platform: \(platform)"
        case .resolutionFailed(let reason):
            return "Deep link resolution failed: \(reason)"
        case .authenticationRequired:
            return "Authentication required for this deep link"
        case .rateLimitExceeded:
            return "Deep link processing rate limit exceeded"
        }
    }
}