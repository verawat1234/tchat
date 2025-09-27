package com.tchat.mobile.di

import com.tchat.mobile.services.NavigationService
import com.tchat.mobile.services.SharingService
import com.tchat.mobile.services.MockNavigationService
import com.tchat.mobile.services.MockSharingService
import org.koin.dsl.module

/**
 * Common dependency injection module
 * Provides shared services and repositories for KMP
 */
val commonModule = module {

    // Platform-specific services using Mock implementations for now
    single<SharingService> { MockSharingService() }
    single<NavigationService> { MockNavigationService() }

    // Add other common services here
    // single<ApiService> { ApiServiceImpl(get()) }
    // single<DatabaseService> { DatabaseServiceImpl() }
    // single<PostRepository> { PostRepositoryImpl(get()) }
    // single<UserRepository> { UserRepositoryImpl(get()) }
}

/**
 * Platform modules should be created in androidMain and iosMain
 * to provide platform-specific dependencies
 */