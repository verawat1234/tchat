package com.tchat.app

import kotlinx.coroutines.test.TestDispatcher
import kotlinx.coroutines.test.StandardTestDispatcher
import kotlinx.coroutines.test.TestScope
import org.junit.Rule
import org.junit.rules.TestRule
import org.junit.runner.Description
import org.junit.runners.model.Statement

/**
 * Global test configuration and utilities for TchatApp
 */
object TestConfiguration {

    /** Test API base URL */
    const val API_BASE_URL = "https://api.tchat.test"

    /** Test timeout intervals */
    object Timeouts {
        const val SHORT = 2000L
        const val MEDIUM = 5000L
        const val LONG = 10000L
    }

    /** Test data utilities */
    object TestData {
        const val VALID_USER_ID = "test-user-123"
        const val VALID_WORKSPACE_ID = "test-workspace-456"
        const val VALID_SESSION_TOKEN = "test-session-token"
    }
}

/**
 * Base test class for all TchatApp tests
 */
abstract class TchatAppTestCase {

    @get:Rule
    val testCoroutineRule = TestCoroutineRule()

    protected val testDispatcher get() = testCoroutineRule.testDispatcher
    protected val testScope get() = testCoroutineRule.testScope
}

/**
 * Test rule for coroutine testing
 */
class TestCoroutineRule : TestRule {

    val testDispatcher: TestDispatcher = StandardTestDispatcher()
    val testScope = TestScope(testDispatcher)

    override fun apply(base: Statement, description: Description): Statement {
        return object : Statement() {
            override fun evaluate() {
                try {
                    base.evaluate()
                } finally {
                    // Modern test scope doesn't need manual cleanup
                }
            }
        }
    }
}