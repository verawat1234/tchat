package com.tchat.navigation

import androidx.compose.runtime.Composable
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.getValue
import androidx.compose.runtime.setValue
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.navigation.NavController
import androidx.navigation.NavHostController
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import java.net.URL

/**
 * Core navigation environment for managing app-wide navigation state
 */
class NavigationEnvironment(
    private val navigationPathManager: NavigationPathManager = NavigationPathManager(),
    private val deepLinkHandlerManager: DeepLinkHandlerManager = DeepLinkHandlerManager(),
    private val routeRegistryManager: RouteRegistryManager = RouteRegistryManager()
) {
    val navigationPath = navigationPathManager
    val deepLinkHandler = deepLinkHandlerManager
    val routeRegistry = routeRegistryManager
}

/**
 * Manages navigation path and current route state
 */
class NavigationPathManager {
    private val _currentRoute = MutableStateFlow<String?>(null)
    val currentRoute: StateFlow<String?> = _currentRoute.asStateFlow()

    private val _navigationStack = MutableStateFlow<List<String>>(emptyList())
    val navigationStack: StateFlow<List<String>> = _navigationStack.asStateFlow()

    private var navController: NavHostController? = null

    fun navigate(to: String, parameters: Map<String, Any> = emptyMap()) {
        _currentRoute.value = to
        _navigationStack.value = _navigationStack.value + to

        // Navigate using NavController if available
        navController?.navigate(to) {
            // Add navigation options as needed
            launchSingleTop = true
        }
    }

    fun goBack() {
        val currentStack = _navigationStack.value
        if (currentStack.isNotEmpty()) {
            _navigationStack.value = currentStack.dropLast(1)
            _currentRoute.value = currentStack.getOrNull(currentStack.size - 2)
            navController?.popBackStack()
        }
    }

    fun popToRoot() {
        _navigationStack.value = emptyList()
        _currentRoute.value = null
        navController?.popBackStack(
            route = "main", // Assuming "main" is the root route
            inclusive = false
        )
    }

    fun setNavController(controller: NavHostController) {
        navController = controller
    }
}

/**
 * Handles deep link processing and URL parsing
 */
class DeepLinkHandlerManager {
    private val _pendingDeepLink = MutableStateFlow<String?>(null)
    val pendingDeepLink: StateFlow<String?> = _pendingDeepLink.asStateFlow()

    fun handle(url: String): Boolean {
        // Parse URL and determine if it's a valid deep link
        return try {
            val uri = android.net.Uri.parse(url)
            if (uri.scheme == "tchat") {
                _pendingDeepLink.value = url
                true
            } else {
                false
            }
        } catch (e: Exception) {
            false
        }
    }

    fun processDeepLink(url: String): String? {
        return try {
            val uri = android.net.Uri.parse(url)
            val pathSegments = uri.pathSegments
            pathSegments?.joinToString("/")
        } catch (e: Exception) {
            null
        }
    }

    fun clearPendingDeepLink() {
        _pendingDeepLink.value = null
    }
}

/**
 * Manages available routes and route validation
 */
class RouteRegistryManager {
    private val _availableRoutes = MutableStateFlow<Set<String>>(emptySet())
    val availableRoutes: StateFlow<Set<String>> = _availableRoutes.asStateFlow()

    init {
        setupDefaultRoutes()
    }

    private fun setupDefaultRoutes() {
        _availableRoutes.value = setOf(
            "chat",
            "store",
            "social",
            "video",
            "more",
            "chat/user",
            "store/products",
            "social/feed",
            "video/call",
            "more/settings"
        )
    }

    fun isValidRoute(route: String): Boolean {
        return _availableRoutes.value.contains(route)
    }

    fun registerRoute(route: String) {
        _availableRoutes.value = _availableRoutes.value + route
    }

    fun getAvailableRoutes(): Set<String> {
        return _availableRoutes.value
    }
}

/**
 * Navigation ViewModel for Compose integration
 */
class NavigationViewModel(
    private val navigationEnvironment: NavigationEnvironment = NavigationEnvironment()
) : ViewModel() {

    val currentRoute = navigationEnvironment.navigationPath.currentRoute
    val navigationStack = navigationEnvironment.navigationPath.navigationStack
    val availableRoutes = navigationEnvironment.routeRegistry.availableRoutes

    fun navigate(route: String, parameters: Map<String, Any> = emptyMap()) {
        viewModelScope.launch {
            navigationEnvironment.navigationPath.navigate(route, parameters)
        }
    }

    fun goBack() {
        viewModelScope.launch {
            navigationEnvironment.navigationPath.goBack()
        }
    }

    fun popToRoot() {
        viewModelScope.launch {
            navigationEnvironment.navigationPath.popToRoot()
        }
    }

    fun handleDeepLink(url: String): Boolean {
        return navigationEnvironment.deepLinkHandler.handle(url)
    }
}