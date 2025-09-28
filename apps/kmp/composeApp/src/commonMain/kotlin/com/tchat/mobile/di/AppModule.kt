package com.tchat.mobile.di

import com.tchat.mobile.repositories.ChatRepository
import com.tchat.mobile.repositories.SQLDelightChatRepository
import com.tchat.mobile.repositories.HybridChatRepository
import com.tchat.mobile.repositories.ProductRepository
import com.tchat.mobile.repositories.EventRepository
import com.tchat.mobile.repositories.datasource.ApiRemoteDataSource
import com.tchat.mobile.services.SocialContentService
import com.tchat.mobile.services.SocialSyncService
import com.tchat.mobile.services.ContentApiService
import com.tchat.mobile.services.AuthService
import com.tchat.mobile.api.ApiClient
import com.tchat.mobile.social.di.socialModule
import org.koin.core.module.Module
import org.koin.dsl.module

/**
 * Common application module for dependency injection
 * Contains shared dependencies across all platforms
 */
val appModule: Module = module {
    // API and data sources
    single<ApiClient> { ApiClient() }
    single<ApiRemoteDataSource> { ApiRemoteDataSource(get()) }

    // Database and repositories - NOW WITH REAL API INTEGRATION!
    single<ChatRepository> { HybridChatRepository(get(), get()) }
    single<ProductRepository> { ProductRepository(get()) }
    single<EventRepository> { EventRepository(get()) }

    // Services
    single<AuthService> { AuthService(get(), get()) }
    single<ContentApiService> { ContentApiService() }
    single<SocialContentService> {
        SocialContentService(
            socialRepository = get(),
            currentUserId = "current_user" // Real user ID for testing
        )
    }
    single<SocialSyncService> { SocialSyncService(get()) }

    // Include social module
    includes(socialModule)
}