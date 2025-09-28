package com.tchat.mobile.di

import com.tchat.mobile.repositories.ChatRepository
import com.tchat.mobile.repositories.SQLDelightChatRepository
import com.tchat.mobile.repositories.SocialRepository
import com.tchat.mobile.repositories.SQLDelightSocialRepository
import com.tchat.mobile.repositories.ProductRepository
import com.tchat.mobile.repositories.EventRepository
import com.tchat.mobile.services.SocialContentService
import com.tchat.mobile.services.SocialSyncService
import com.tchat.mobile.services.ContentApiService
import org.koin.core.module.Module
import org.koin.dsl.module

/**
 * Common application module for dependency injection
 * Contains shared dependencies across all platforms
 */
val appModule: Module = module {
    // Database and repositories
    single<ChatRepository> { SQLDelightChatRepository(get()) }
    single<SocialRepository> { SQLDelightSocialRepository(get()) }
    single<ProductRepository> { ProductRepository(get()) }
    single<EventRepository> { EventRepository(get()) }

    // Services
    single<ContentApiService> { ContentApiService() }
    single<SocialContentService> {
        SocialContentService(
            socialRepository = get(),
            currentUserId = "current_user_id" // TODO: Get from auth service
        )
    }
    single<SocialSyncService> { SocialSyncService(get()) }
}