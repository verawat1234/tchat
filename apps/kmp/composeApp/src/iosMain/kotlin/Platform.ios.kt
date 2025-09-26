package com.tchat.mobile

import org.koin.core.module.Module
import org.koin.dsl.module
import com.tchat.mobile.di.appModule
import com.tchat.mobile.di.platformModule
import platform.Foundation.NSBundle
import platform.Foundation.NSString
import platform.Foundation.NSURL
import platform.Foundation.stringWithFormat
import platform.UIKit.UIDevice

/**
 * iOS platform-specific implementations and dependencies
 */
class IOSPlatform: Platform {
    override val name: String = UIDevice.currentDevice.systemName() + " " + UIDevice.currentDevice.systemVersion
}

fun getPlatform(): Platform = IOSPlatform()


/**
 * Initialize Koin for iOS
 */
fun initKoin() {
    org.koin.core.context.startKoin {
        modules(
            appModule,
            platformModule
        )
    }
}

/**
 * Platform interface
 */
interface Platform {
    val name: String
}