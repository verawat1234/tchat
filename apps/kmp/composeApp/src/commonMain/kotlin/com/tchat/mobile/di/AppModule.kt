package com.tchat.mobile.di

import org.koin.core.module.Module
import org.koin.dsl.module

/**
 * Common application module for dependency injection
 * Contains shared dependencies across all platforms
 */
val appModule: Module = module {
    // Common services and dependencies will be added here
    // Examples:
    // single<ApiClient> { ApiClientImpl() }
    // single<DataRepository> { DataRepositoryImpl(get()) }
    // single<SettingsManager> { SettingsManagerImpl() }
}