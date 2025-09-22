package com.tchat.services

import com.tchat.models.*
import kotlinx.coroutines.flow.*
import kotlinx.coroutines.sync.Mutex
import kotlinx.coroutines.sync.withLock
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import java.util.UUID
import java.util.regex.Pattern

/**
 * Route registry service for managing application routes and route validation
 */
class RouteRegistry(
    private val persistenceService: RouteRegistryPersistence = RouteRegistryPersistence(),
    private val validationService: RouteValidationService = RouteValidationService()
) {

    private val _availableRoutes = MutableStateFlow<Set<NavigationRoute>>(emptySet())
    val availableRoutes: StateFlow<Set<NavigationRoute>> = _availableRoutes.asStateFlow()

    private val _routeHierarchy = MutableStateFlow<Map<String, List<String>>>(emptyMap())
    val routeHierarchy: StateFlow<Map<String, List<String>>> = _routeHierarchy.asStateFlow()

    private val _isLoading = MutableStateFlow(false)
    val isLoading: StateFlow<Boolean> = _isLoading.asStateFlow()

    private val routeCache = mutableMapOf<String, NavigationRoute>()
    private val parameterSchemas = mutableMapOf<String, Map<String, RouteParameter>>()
    private val registryMutex = Mutex()

    init {
        loadRoutes()
    }

    /**
     * Register a new route
     */
    suspend fun registerRoute(route: NavigationRoute) {
        registryMutex.withLock {
            // Validate route before registration
            validationService.validateRoute(route)

            // Check for conflicts
            val existingRoute = routeCache[route.id]
            if (existingRoute != null && existingRoute != route) {
                throw RouteRegistryError.RouteConflict(route.id)
            }

            // Register route
            routeCache[route.id] = route
            _availableRoutes.value = _availableRoutes.value + route

            // Update hierarchy
            updateRouteHierarchy(route)

            // Cache parameter schema
            parameterSchemas[route.id] = route.parameters

            // Persist changes
            persistenceService.saveRoutes(_availableRoutes.value.toList())
        }
    }

    /**
     * Unregister a route
     */
    suspend fun unregisterRoute(routeId: String) {
        registryMutex.withLock {
            val route = routeCache[routeId]
                ?: throw RouteRegistryError.RouteNotFound(routeId)

            // Remove from cache
            routeCache.remove(routeId)
            _availableRoutes.value = _availableRoutes.value - route

            // Update hierarchy
            removeFromRouteHierarchy(routeId)

            // Remove parameter schema
            parameterSchemas.remove(routeId)

            // Persist changes
            persistenceService.saveRoutes(_availableRoutes.value.toList())
        }
    }

    /**
     * Get route by ID
     */
    fun getRoute(routeId: String): NavigationRoute? {
        return routeCache[routeId]
    }

    /**
     * Check if route is valid
     */
    fun isValidRoute(routeId: String): Boolean {
        return routeCache.containsKey(routeId)
    }

    /**
     * Get routes matching pattern
     */
    fun getRoutes(pattern: String): List<NavigationRoute> {
        return _availableRoutes.value.filter { route ->
            route.path.contains(pattern) || route.title.contains(pattern)
        }
    }

    /**
     * Get child routes for a parent route
     */
    fun getChildRoutes(parentRouteId: String): List<NavigationRoute> {
        return _availableRoutes.value.filter { it.parentRouteId == parentRouteId }
    }

    /**
     * Get route hierarchy path
     */
    fun getRoutePath(routeId: String): List<NavigationRoute> {
        val path = mutableListOf<NavigationRoute>()
        var currentRouteId: String? = routeId

        while (currentRouteId != null) {
            val route = routeCache[currentRouteId] ?: break
            path.add(0, route)
            currentRouteId = route.parentRouteId
        }

        return path
    }

    /**
     * Validate route parameters
     */
    fun validateParameters(routeId: String, parameters: Map<String, Any>) {
        val route = routeCache[routeId]
            ?: throw RouteRegistryError.RouteNotFound(routeId)

        validationService.validateParameters(route, parameters)
    }

    /**
     * Get parameter schema for route
     */
    fun getParameterSchema(routeId: String): Map<String, RouteParameter>? {
        return parameterSchemas[routeId]
    }

    /**
     * Generate route URL with parameters
     */
    fun generateURL(routeId: String, parameters: Map<String, Any> = emptyMap()): String {
        val route = routeCache[routeId]
            ?: throw RouteRegistryError.RouteNotFound(routeId)

        return route.fullPath(parameters)
    }

    /**
     * Find routes by access level
     */
    fun getRoutes(withAccessLevel: AccessLevel): List<NavigationRoute> {
        return _availableRoutes.value.filter { it.accessLevel == withAccessLevel }
    }

    /**
     * Find routes available on current platform
     */
    fun getPlatformRoutes(): List<NavigationRoute> {
        return _availableRoutes.value.filter { it.isAvailableOnCurrentPlatform }
    }

    /**
     * Search routes by metadata tags
     */
    fun searchRoutes(withTags: List<String>): List<NavigationRoute> {
        return _availableRoutes.value.filter { route ->
            route.metadata.tags.any { it in withTags }
        }
    }

    // MARK: - Private Methods

    private fun loadRoutes() {
        CoroutineScope(Dispatchers.IO).launch {
            _isLoading.value = true

            try {
                // Load persisted routes
                val persistedRoutes = persistenceService.loadRoutes()

                // Load default routes if no persisted routes
                val routesToLoad = if (persistedRoutes.isEmpty()) {
                    NavigationRoute.defaultRoutes
                } else {
                    persistedRoutes
                }

                // Register all routes
                registryMutex.withLock {
                    val routes = mutableSetOf<NavigationRoute>()
                    for (route in routesToLoad) {
                        routeCache[route.id] = route
                        routes.add(route)
                        parameterSchemas[route.id] = route.parameters
                    }
                    _availableRoutes.value = routes
                }

                // Build hierarchy
                buildRouteHierarchy()

            } catch (e: Exception) {
                // Fallback to default routes
                registryMutex.withLock {
                    val routes = mutableSetOf<NavigationRoute>()
                    for (route in NavigationRoute.defaultRoutes) {
                        routeCache[route.id] = route
                        routes.add(route)
                        parameterSchemas[route.id] = route.parameters
                    }
                    _availableRoutes.value = routes
                }

                buildRouteHierarchy()
            } finally {
                _isLoading.value = false
            }
        }
    }

    private fun buildRouteHierarchy() {
        val hierarchy = mutableMapOf<String, MutableList<String>>()

        for (route in _availableRoutes.value) {
            val parentKey = route.parentRouteId ?: "root"
            hierarchy.getOrPut(parentKey) { mutableListOf() }.add(route.id)
        }

        _routeHierarchy.value = hierarchy.mapValues { it.value.toList() }
    }

    private fun updateRouteHierarchy(route: NavigationRoute) {
        val currentHierarchy = _routeHierarchy.value.toMutableMap()
        val parentKey = route.parentRouteId ?: "root"

        currentHierarchy.getOrPut(parentKey) { emptyList() }
        currentHierarchy[parentKey] = currentHierarchy[parentKey]!! + route.id

        _routeHierarchy.value = currentHierarchy
    }

    private fun removeFromRouteHierarchy(routeId: String) {
        val currentHierarchy = _routeHierarchy.value.toMutableMap()

        for ((parentId, children) in currentHierarchy) {
            if (children.contains(routeId)) {
                currentHierarchy[parentId] = children - routeId
                break
            }
        }

        _routeHierarchy.value = currentHierarchy
    }
}

// MARK: - Supporting Services

/**
 * Route registry persistence service
 */
class RouteRegistryPersistence {

    suspend fun saveRoutes(routes: List<NavigationRoute>) {
        // TODO: Implement persistence (SharedPreferences, Room, etc.)
    }

    suspend fun loadRoutes(): List<NavigationRoute> {
        // TODO: Implement loading
        return emptyList()
    }
}

/**
 * Route validation service
 */
class RouteValidationService {

    fun validateRoute(route: NavigationRoute) {
        // Validate route ID
        if (route.id.isEmpty()) {
            throw RouteRegistryError.InvalidRouteId()
        }

        // Validate path
        if (route.path.isEmpty()) {
            throw RouteRegistryError.InvalidPath()
        }

        // Validate component
        if (route.component.isEmpty()) {
            throw RouteRegistryError.InvalidComponent()
        }

        // Validate parameters
        for ((_, parameter) in route.parameters) {
            validateParameter(parameter)
        }
    }

    fun validateParameters(route: NavigationRoute, parameters: Map<String, Any>) {
        // Check required parameters
        for ((paramName, paramDef) in route.parameters) {
            if (paramDef.isRequired && !parameters.containsKey(paramName)) {
                throw RouteRegistryError.MissingRequiredParameter(paramName)
            }
        }

        // Validate parameter types and values
        for ((paramName, value) in parameters) {
            val paramDef = route.parameters[paramName] ?: continue // Unknown parameter, ignore

            validateParameterValue(value, paramDef)
        }
    }

    private fun validateParameter(parameter: RouteParameter) {
        if (parameter.name.isEmpty()) {
            throw RouteRegistryError.InvalidParameterName()
        }

        // Validate validation rules if present
        parameter.validation?.let { validateValidationRule(it) }
    }

    private fun validateParameterValue(value: Any, parameter: RouteParameter) {
        when (parameter.type) {
            ParameterType.STRING -> {
                if (value !is String) {
                    throw RouteRegistryError.InvalidParameterType(parameter.name, "String")
                }
            }
            ParameterType.INTEGER -> {
                if (value !is Int) {
                    throw RouteRegistryError.InvalidParameterType(parameter.name, "Int")
                }
            }
            ParameterType.UUID -> {
                if (value !is String) {
                    throw RouteRegistryError.InvalidParameterType(parameter.name, "UUID String")
                }
                try {
                    UUID.fromString(value)
                } catch (e: IllegalArgumentException) {
                    throw RouteRegistryError.InvalidParameterType(parameter.name, "Valid UUID")
                }
            }
            ParameterType.BOOLEAN -> {
                if (value !is Boolean) {
                    throw RouteRegistryError.InvalidParameterType(parameter.name, "Boolean")
                }
            }
            ParameterType.URL -> {
                if (value !is String) {
                    throw RouteRegistryError.InvalidParameterType(parameter.name, "URL String")
                }
                try {
                    java.net.URL(value)
                } catch (e: Exception) {
                    throw RouteRegistryError.InvalidParameterType(parameter.name, "Valid URL")
                }
            }
        }

        // Apply validation rules
        parameter.validation?.let { applyValidationRule(it, value, parameter.name) }
    }

    private fun validateValidationRule(rule: ValidationRule) {
        // Validate pattern if present
        rule.pattern?.let { pattern ->
            try {
                Pattern.compile(pattern)
            } catch (e: Exception) {
                throw RouteRegistryError.InvalidValidationPattern(pattern)
            }
        }

        // Validate length constraints
        val minLength = rule.minLength
        val maxLength = rule.maxLength
        if (minLength != null && maxLength != null && minLength > maxLength) {
            throw RouteRegistryError.InvalidLengthConstraints()
        }
    }

    private fun applyValidationRule(rule: ValidationRule, value: Any, parameterName: String) {
        if (value is String) {
            // Length validation
            rule.minLength?.let { minLength ->
                if (value.length < minLength) {
                    throw RouteRegistryError.ParameterTooShort(parameterName, minLength)
                }
            }

            rule.maxLength?.let { maxLength ->
                if (value.length > maxLength) {
                    throw RouteRegistryError.ParameterTooLong(parameterName, maxLength)
                }
            }

            // Pattern validation
            rule.pattern?.let { pattern ->
                if (!Pattern.matches(pattern, value)) {
                    throw RouteRegistryError.ParameterPatternMismatch(parameterName, pattern)
                }
            }

            // Allowed values validation
            rule.allowedValues?.let { allowedValues ->
                if (!allowedValues.contains(value)) {
                    throw RouteRegistryError.ParameterNotInAllowedValues(parameterName, allowedValues)
                }
            }
        }
    }
}

// MARK: - Error Types

sealed class RouteRegistryError : Exception() {
    class RouteNotFound(val id: String) : RouteRegistryError()
    class RouteConflict(val id: String) : RouteRegistryError()
    class InvalidRouteId : RouteRegistryError()
    class InvalidPath : RouteRegistryError()
    class InvalidComponent : RouteRegistryError()
    class InvalidParameterName : RouteRegistryError()
    class InvalidParameterType(val name: String, val expected: String) : RouteRegistryError()
    class MissingRequiredParameter(val name: String) : RouteRegistryError()
    class InvalidValidationPattern(val pattern: String) : RouteRegistryError()
    class InvalidLengthConstraints : RouteRegistryError()
    class ParameterTooShort(val name: String, val minLength: Int) : RouteRegistryError()
    class ParameterTooLong(val name: String, val maxLength: Int) : RouteRegistryError()
    class ParameterPatternMismatch(val name: String, val pattern: String) : RouteRegistryError()
    class ParameterNotInAllowedValues(val name: String, val values: List<String>) : RouteRegistryError()

    override val message: String?
        get() = when (this) {
            is RouteNotFound -> "Route not found: $id"
            is RouteConflict -> "Route conflict: $id"
            is InvalidRouteId -> "Invalid route ID"
            is InvalidPath -> "Invalid route path"
            is InvalidComponent -> "Invalid route component"
            is InvalidParameterName -> "Invalid parameter name"
            is InvalidParameterType -> "Invalid parameter type for $name, expected $expected"
            is MissingRequiredParameter -> "Missing required parameter: $name"
            is InvalidValidationPattern -> "Invalid validation pattern: $pattern"
            is InvalidLengthConstraints -> "Invalid length constraints"
            is ParameterTooShort -> "Parameter $name too short, minimum length: $minLength"
            is ParameterTooLong -> "Parameter $name too long, maximum length: $maxLength"
            is ParameterPatternMismatch -> "Parameter $name doesn't match pattern: $pattern"
            is ParameterNotInAllowedValues -> "Parameter $name not in allowed values: $values"
        }
}