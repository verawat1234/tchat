package com.tchat.mobile.config

import com.tchat.mobile.BuildConfig

/**
 * Android-specific Network Configuration
 */

/**
 * Check if this is a debug build on Android
 */
actual fun NetworkConfig.isDebugBuild(): Boolean {
    return BuildConfig.DEBUG
}

/**
 * Get Android emulator localhost URL
 * Android emulator uses 10.0.2.2 to access host machine's localhost
 */
actual fun NetworkConfig.getLocalHostUrl(): String {
    return "http://10.0.2.2:8080"
}
