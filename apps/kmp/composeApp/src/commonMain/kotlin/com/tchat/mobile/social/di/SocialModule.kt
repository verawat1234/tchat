package com.tchat.mobile.social.di

import com.tchat.mobile.social.data.api.SocialApiClient
import com.tchat.mobile.social.data.repository.SocialRepository
import com.tchat.mobile.social.presentation.SocialViewModel
import com.tchat.mobile.social.services.RegionalContentService
import org.koin.core.module.dsl.singleOf
import org.koin.dsl.module

/**
 * KMP Social Dependency Injection Module
 *
 * Provides social feature dependencies with:
 * - Cross-platform compatibility
 * - Singleton instances for optimal performance
 * - Proper lifecycle management
 */
val socialModule = module {

    // API Client
    single<SocialApiClient> {
        SocialApiClient(
            baseUrl = "http://localhost:8080/api/v1/social"
        )
    }

    // Repository
    single<SocialRepository> {
        SocialRepository(
            apiClient = get(),
            database = get()
        )
    }

    // Regional Content Service
    single<RegionalContentService> { RegionalContentService() }

    // ViewModel - Use factory instead of viewModelOf for KMP compatibility
    factory { SocialViewModel(get(), get()) }
}