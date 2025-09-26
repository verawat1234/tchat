package com.tchat.mobile.di

import org.koin.core.module.Module

/**
 * Platform-specific dependency injection module
 * Uses expect/actual pattern for cross-platform DI
 */
expect val platformModule: Module