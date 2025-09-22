//
//  RouteRegistry.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
import Combine

/// Route registry service for managing application routes and route validation
@MainActor
public class RouteRegistry: ObservableObject {

    // MARK: - Published Properties

    @Published public var availableRoutes: Set<NavigationRoute> = []
    @Published public var routeHierarchy: [String: [String]] = [:]
    @Published public var isLoading: Bool = false

    // MARK: - Private Properties

    private var routeCache: [String: NavigationRoute] = [:]
    private var parameterSchemas: [String: [String: RouteParameter]] = [:]
    private let persistenceService: RouteRegistryPersistence
    private let validationService: RouteValidationService

    // MARK: - Initialization

    public init(
        persistenceService: RouteRegistryPersistence = RouteRegistryPersistence(),
        validationService: RouteValidationService = RouteValidationService()
    ) {
        self.persistenceService = persistenceService
        self.validationService = validationService

        loadRoutes()
    }

    // MARK: - Public Methods

    /// Register a new route
    public func registerRoute(_ route: NavigationRoute) async throws {
        // Validate route before registration
        try await validationService.validateRoute(route)

        // Check for conflicts
        if let existingRoute = routeCache[route.id], existingRoute != route {
            throw RouteRegistryError.routeConflict(route.id)
        }

        // Register route
        routeCache[route.id] = route
        availableRoutes.insert(route)

        // Update hierarchy
        updateRouteHierarchy(for: route)

        // Cache parameter schema
        parameterSchemas[route.id] = route.parameters

        // Persist changes
        try await persistenceService.saveRoutes(Array(availableRoutes))
    }

    /// Unregister a route
    public func unregisterRoute(_ routeId: String) async throws {
        guard let route = routeCache[routeId] else {
            throw RouteRegistryError.routeNotFound(routeId)
        }

        // Remove from cache
        routeCache.removeValue(forKey: routeId)
        availableRoutes.remove(route)

        // Update hierarchy
        removeFromRouteHierarchy(routeId: routeId)

        // Remove parameter schema
        parameterSchemas.removeValue(forKey: routeId)

        // Persist changes
        try await persistenceService.saveRoutes(Array(availableRoutes))
    }

    /// Get route by ID
    public func getRoute(_ routeId: String) -> NavigationRoute? {
        return routeCache[routeId]
    }

    /// Check if route is valid
    public func isValidRoute(_ routeId: String) -> Bool {
        return routeCache[routeId] != nil
    }

    /// Get routes by category
    public func getRoutes(matching pattern: String) -> [NavigationRoute] {
        return availableRoutes.filter { route in
            route.path.contains(pattern) || route.title.contains(pattern)
        }
    }

    /// Get child routes for a parent route
    public func getChildRoutes(for parentRouteId: String) -> [NavigationRoute] {
        return availableRoutes.filter { $0.parentRouteId == parentRouteId }
    }

    /// Get route hierarchy path
    public func getRoutePath(for routeId: String) -> [NavigationRoute] {
        var path: [NavigationRoute] = []
        var currentRouteId: String? = routeId

        while let routeId = currentRouteId, let route = routeCache[routeId] {
            path.insert(route, at: 0)
            currentRouteId = route.parentRouteId
        }

        return path
    }

    /// Validate route parameters
    public func validateParameters(
        for routeId: String,
        parameters: [String: Any]
    ) throws {
        guard let route = routeCache[routeId] else {
            throw RouteRegistryError.routeNotFound(routeId)
        }

        try validationService.validateParameters(route: route, parameters: parameters)
    }

    /// Get parameter schema for route
    public func getParameterSchema(for routeId: String) -> [String: RouteParameter]? {
        return parameterSchemas[routeId]
    }

    /// Generate route URL with parameters
    public func generateURL(
        for routeId: String,
        parameters: [String: Any] = [:]
    ) throws -> String {
        guard let route = routeCache[routeId] else {
            throw RouteRegistryError.routeNotFound(routeId)
        }

        return route.fullPath(with: parameters)
    }

    /// Find routes by access level
    public func getRoutes(withAccessLevel accessLevel: AccessLevel) -> [NavigationRoute] {
        return availableRoutes.filter { $0.accessLevel == accessLevel }
    }

    /// Find routes available on current platform
    public func getPlatformRoutes() -> [NavigationRoute] {
        return availableRoutes.filter { $0.isAvailableOnCurrentPlatform }
    }

    /// Search routes by metadata tags
    public func searchRoutes(withTags tags: [String]) -> [NavigationRoute] {
        return availableRoutes.filter { route in
            !Set(route.metadata.tags).isDisjoint(with: Set(tags))
        }
    }

    // MARK: - Private Methods

    private func loadRoutes() {
        Task {
            isLoading = true
            defer { isLoading = false }

            do {
                // Load persisted routes
                let persistedRoutes = try await persistenceService.loadRoutes()

                // Load default routes if no persisted routes
                let routesToLoad = persistedRoutes.isEmpty ? NavigationRoute.defaultRoutes : persistedRoutes

                // Register all routes
                for route in routesToLoad {
                    routeCache[route.id] = route
                    availableRoutes.insert(route)
                    parameterSchemas[route.id] = route.parameters
                }

                // Build hierarchy
                buildRouteHierarchy()

            } catch {
                // Fallback to default routes
                for route in NavigationRoute.defaultRoutes {
                    routeCache[route.id] = route
                    availableRoutes.insert(route)
                    parameterSchemas[route.id] = route.parameters
                }

                buildRouteHierarchy()
            }
        }
    }

    private func buildRouteHierarchy() {
        routeHierarchy.removeAll()

        for route in availableRoutes {
            if let parentId = route.parentRouteId {
                routeHierarchy[parentId, default: []].append(route.id)
            } else {
                routeHierarchy["root", default: []].append(route.id)
            }
        }
    }

    private func updateRouteHierarchy(for route: NavigationRoute) {
        if let parentId = route.parentRouteId {
            routeHierarchy[parentId, default: []].append(route.id)
        } else {
            routeHierarchy["root", default: []].append(route.id)
        }
    }

    private func removeFromRouteHierarchy(routeId: String) {
        for (parentId, children) in routeHierarchy {
            if let index = children.firstIndex(of: routeId) {
                routeHierarchy[parentId]?.remove(at: index)
                break
            }
        }
    }
}

// MARK: - Supporting Services

/// Route registry persistence service
public class RouteRegistryPersistence {

    public init() {}

    public func saveRoutes(_ routes: [NavigationRoute]) async throws {
        // TODO: Implement persistence (UserDefaults, Core Data, etc.)
    }

    public func loadRoutes() async throws -> [NavigationRoute] {
        // TODO: Implement loading
        return []
    }
}

/// Route validation service
public class RouteValidationService {

    public init() {}

    public func validateRoute(_ route: NavigationRoute) async throws {
        // Validate route ID
        guard !route.id.isEmpty else {
            throw RouteRegistryError.invalidRouteId
        }

        // Validate path
        guard !route.path.isEmpty else {
            throw RouteRegistryError.invalidPath
        }

        // Validate component
        guard !route.component.isEmpty else {
            throw RouteRegistryError.invalidComponent
        }

        // Validate parameters
        for (_, parameter) in route.parameters {
            try validateParameter(parameter)
        }
    }

    public func validateParameters(
        route: NavigationRoute,
        parameters: [String: Any]
    ) throws {
        // Check required parameters
        for (paramName, paramDef) in route.parameters {
            if paramDef.isRequired && parameters[paramName] == nil {
                throw RouteRegistryError.missingRequiredParameter(paramName)
            }
        }

        // Validate parameter types and values
        for (paramName, value) in parameters {
            guard let paramDef = route.parameters[paramName] else {
                continue // Unknown parameter, ignore
            }

            try validateParameterValue(value, against: paramDef)
        }
    }

    private func validateParameter(_ parameter: RouteParameter) throws {
        guard !parameter.name.isEmpty else {
            throw RouteRegistryError.invalidParameterName
        }

        // Validate validation rules if present
        if let validation = parameter.validation {
            try validateValidationRule(validation)
        }
    }

    private func validateParameterValue(_ value: Any, against parameter: RouteParameter) throws {
        switch parameter.type {
        case .string:
            guard value is String else {
                throw RouteRegistryError.invalidParameterType(parameter.name, expected: "String")
            }
        case .integer:
            guard value is Int else {
                throw RouteRegistryError.invalidParameterType(parameter.name, expected: "Int")
            }
        case .uuid:
            guard let stringValue = value as? String,
                  UUID(uuidString: stringValue) != nil else {
                throw RouteRegistryError.invalidParameterType(parameter.name, expected: "UUID")
            }
        case .boolean:
            guard value is Bool else {
                throw RouteRegistryError.invalidParameterType(parameter.name, expected: "Bool")
            }
        case .url:
            guard let stringValue = value as? String,
                  URL(string: stringValue) != nil else {
                throw RouteRegistryError.invalidParameterType(parameter.name, expected: "URL")
            }
        }

        // Apply validation rules
        if let validation = parameter.validation {
            try applyValidationRule(validation, to: value, parameterName: parameter.name)
        }
    }

    private func validateValidationRule(_ rule: ValidationRule) throws {
        // Validate pattern if present
        if let pattern = rule.pattern {
            do {
                _ = try NSRegularExpression(pattern: pattern, options: [])
            } catch {
                throw RouteRegistryError.invalidValidationPattern(pattern)
            }
        }

        // Validate length constraints
        if let minLength = rule.minLength, let maxLength = rule.maxLength {
            guard minLength <= maxLength else {
                throw RouteRegistryError.invalidLengthConstraints
            }
        }
    }

    private func applyValidationRule(
        _ rule: ValidationRule,
        to value: Any,
        parameterName: String
    ) throws {
        if let stringValue = value as? String {
            // Length validation
            if let minLength = rule.minLength, stringValue.count < minLength {
                throw RouteRegistryError.parameterTooShort(parameterName, minLength)
            }

            if let maxLength = rule.maxLength, stringValue.count > maxLength {
                throw RouteRegistryError.parameterTooLong(parameterName, maxLength)
            }

            // Pattern validation
            if let pattern = rule.pattern {
                let regex = try NSRegularExpression(pattern: pattern, options: [])
                let range = NSRange(location: 0, length: stringValue.count)
                if regex.firstMatch(in: stringValue, options: [], range: range) == nil {
                    throw RouteRegistryError.parameterPatternMismatch(parameterName, pattern)
                }
            }

            // Allowed values validation
            if let allowedValues = rule.allowedValues, !allowedValues.contains(stringValue) {
                throw RouteRegistryError.parameterNotInAllowedValues(parameterName, allowedValues)
            }
        }
    }
}

// MARK: - Error Types

public enum RouteRegistryError: Error, LocalizedError {
    case routeNotFound(String)
    case routeConflict(String)
    case invalidRouteId
    case invalidPath
    case invalidComponent
    case invalidParameterName
    case invalidParameterType(String, expected: String)
    case missingRequiredParameter(String)
    case invalidValidationPattern(String)
    case invalidLengthConstraints
    case parameterTooShort(String, Int)
    case parameterTooLong(String, Int)
    case parameterPatternMismatch(String, String)
    case parameterNotInAllowedValues(String, [String])

    public var errorDescription: String? {
        switch self {
        case .routeNotFound(let id):
            return "Route not found: \(id)"
        case .routeConflict(let id):
            return "Route conflict: \(id)"
        case .invalidRouteId:
            return "Invalid route ID"
        case .invalidPath:
            return "Invalid route path"
        case .invalidComponent:
            return "Invalid route component"
        case .invalidParameterName:
            return "Invalid parameter name"
        case .invalidParameterType(let name, let expected):
            return "Invalid parameter type for \(name), expected \(expected)"
        case .missingRequiredParameter(let name):
            return "Missing required parameter: \(name)"
        case .invalidValidationPattern(let pattern):
            return "Invalid validation pattern: \(pattern)"
        case .invalidLengthConstraints:
            return "Invalid length constraints"
        case .parameterTooShort(let name, let minLength):
            return "Parameter \(name) too short, minimum length: \(minLength)"
        case .parameterTooLong(let name, let maxLength):
            return "Parameter \(name) too long, maximum length: \(maxLength)"
        case .parameterPatternMismatch(let name, let pattern):
            return "Parameter \(name) doesn't match pattern: \(pattern)"
        case .parameterNotInAllowedValues(let name, let values):
            return "Parameter \(name) not in allowed values: \(values)"
        }
    }
}