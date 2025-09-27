package com.tchat.mobile.database

import app.cash.sqldelight.db.SqlDriver

/**
 * Database driver factory for SQLDelight
 * Platform-specific implementations provide the actual driver
 */
expect class DatabaseDriverFactory {
    fun createDriver(): SqlDriver
}