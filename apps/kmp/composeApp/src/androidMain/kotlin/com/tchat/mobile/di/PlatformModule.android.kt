package com.tchat.mobile.di

import android.content.Context
import app.cash.sqldelight.db.SqlDriver
import app.cash.sqldelight.driver.android.AndroidSqliteDriver
import com.tchat.mobile.database.TchatDatabase
import com.tchat.mobile.services.UserSeedingService
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import org.koin.android.ext.koin.androidContext
import org.koin.core.module.Module
import org.koin.dsl.module

/**
 * Android-specific dependency injection module
 * Provides Android platform implementations
 */
actual val platformModule: Module = module {
    // SQLDelight Android driver
    single<SqlDriver> {
        AndroidSqliteDriver(
            schema = TchatDatabase.Schema,
            context = androidContext(),
            name = "tchat.db"
        )
    }

    // Database instance
    single<TchatDatabase> {
        TchatDatabase(get())
    }

    // Application scope for background operations
    single<CoroutineScope> {
        CoroutineScope(Dispatchers.Main + SupervisorJob())
    }

    // User seeding service
    single<UserSeedingService> {
        UserSeedingService(
            database = get(),
            socialRepository = get(),
            chatRepository = get()
        )
    }
}