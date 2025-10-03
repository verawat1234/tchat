package com.tchat.mobile.config

import platform.Foundation.NSBundle

/**
 * iOS-specific Network Configuration
 */

/**
 * Check if this is a debug build on iOS
 * Uses DEBUG preprocessor flag from Xcode build configuration
 */
actual fun isDebugBuild(): Boolean {
    // Check if the app is running in debug mode
    // In iOS, we check for the DEBUG flag in the bundle
    return try {
        val bundle = NSBundle.mainBundle
        val bundlePath = bundle.bundlePath
        // Debug builds typically have "Debug" in their path
        bundlePath.contains("Debug", ignoreCase = true)
    } catch (e: Exception) {
        // Default to production if we can't determine
        false
    }
}

/**
 * Get iOS simulator localhost URL
 * iOS simulator can access localhost directly
 */
actual fun getLocalHostUrl(): String {
    return "http://localhost:8080"
}
