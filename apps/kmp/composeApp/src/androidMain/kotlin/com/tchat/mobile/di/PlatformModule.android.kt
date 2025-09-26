package com.tchat.mobile.di

import org.koin.core.module.Module
import org.koin.dsl.module

/**
 * Android-specific dependency injection module
 * Provides Android platform implementations
 */
actual val platformModule: Module = module {
    // Android-specific services and dependencies

    // Android-specific services can be added here
    // Examples:
    // single<BiometricAuthService> { AndroidBiometricAuthService(get()) }
    // single<PushNotificationService> { AndroidPushNotificationService(get()) }
    // single<SharedPreferencesService> { AndroidSharedPreferencesService(get()) }
    // single<DeepLinkService> { AndroidDeepLinkService(get()) }
}