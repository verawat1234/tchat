package com.tchat.mobile

import android.app.Application
import com.tchat.mobile.di.appModule
import com.tchat.mobile.di.platformModule
import com.tchat.mobile.services.UserSeedingService
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.launch
import org.koin.android.ext.android.inject
import org.koin.android.ext.koin.androidContext
import org.koin.core.context.startKoin

/**
 * Application class for Tchat KMP app
 * Initializes dependency injection and seeds initial data
 */
class TchatApplication : Application() {

    override fun onCreate() {
        super.onCreate()

        // Initialize Koin dependency injection
        startKoin {
            androidContext(this@TchatApplication)
            modules(
                appModule,
                platformModule
            )
        }

        // Seed initial data for app startup
        seedInitialData()
    }

    private fun seedInitialData() {
        try {
            // Get services from Koin after initialization
            val userSeedingService: UserSeedingService by inject()
            val scope: CoroutineScope by inject()

            scope.launch {
                try {
                    // Check if data is already seeded
                    if (!userSeedingService.isDataSeeded()) {
                        println("🌱 Seeding initial data for first app launch...")

                        val result = userSeedingService.seedInitialData()
                        if (result.isSuccess) {
                            println("✅ App ready with test data!")
                        } else {
                            println("❌ Failed to seed data: ${result.exceptionOrNull()?.message}")
                        }
                    } else {
                        println("✅ App data already exists, ready to go!")
                    }
                } catch (e: Exception) {
                    println("❌ Error during data seeding: ${e.message}")
                    e.printStackTrace()
                }
            }
        } catch (e: Exception) {
            println("❌ Failed to initialize seeding service: ${e.message}")
            e.printStackTrace()
        }
    }
}