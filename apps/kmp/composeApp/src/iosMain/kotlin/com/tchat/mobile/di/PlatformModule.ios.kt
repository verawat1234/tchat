package com.tchat.mobile.di

import org.koin.core.module.Module
import org.koin.dsl.module
import com.tchat.mobile.Platform
import com.tchat.mobile.IOSPlatform

/**
 * iOS-specific dependency injection module
 * Provides iOS platform implementations
 */
actual val platformModule: Module = module {
    // iOS-specific services and dependencies

    // Platform info
    single<Platform> { IOSPlatform() }

    // iOS-specific services can be added here
    // Examples:
    // single<BiometricAuthService> { IOSBiometricAuthService() }
    // single<PushNotificationService> { IOSPushNotificationService() }
    // single<KeychainService> { IOSKeychainService() }
    // single<DeepLinkService> { IOSDeepLinkService() }
}