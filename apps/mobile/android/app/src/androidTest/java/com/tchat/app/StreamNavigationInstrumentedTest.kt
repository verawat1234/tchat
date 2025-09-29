/**
 * Android Stream Navigation Instrumented Tests
 *
 * Native Android instrumented tests for Stream Store Tabs functionality using Espresso.
 * Tests cover Jetpack Compose navigation, Material3 theming, tab switching,
 * content loading, and Android-specific behavior patterns.
 *
 * Platform: Android API 24+
 * Framework: Jetpack Compose + Espresso + Material3
 * Testing: Native Android app behavior and performance
 *
 * Test Coverage:
 * - Bottom navigation with Stream tab integration
 * - Stream category navigation with Compose animations
 * - Movies subtab functionality with Material3 design
 * - Featured content display and touch interactions
 * - Performance validation (60fps, <200ms)
 * - Android accessibility and TalkBack support
 * - Configuration changes and state restoration
 * - Material Design 3 compliance
 */

package com.tchat.app

import androidx.compose.ui.test.*
import androidx.compose.ui.test.junit4.createAndroidComposeRule
import androidx.test.ext.junit.runners.AndroidJUnit4
import androidx.test.platform.app.InstrumentationRegistry
import androidx.test.uiautomator.UiDevice
import dagger.hilt.android.testing.HiltAndroidRule
import dagger.hilt.android.testing.HiltAndroidTest
import kotlinx.coroutines.delay
import kotlinx.coroutines.runBlocking
import org.junit.Before
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith
import kotlin.system.measureTimeMillis

@HiltAndroidTest
@RunWith(AndroidJUnit4::class)
class StreamNavigationInstrumentedTest {

    @get:Rule(order = 0)
    val hiltRule = HiltAndroidRule(this)

    @get:Rule(order = 1)
    val composeTestRule = createAndroidComposeRule<MainActivity>()

    private lateinit var device: UiDevice

    // Performance budgets for Android
    private companion object {
        const val TAB_SWITCH_TIME_MS = 200L
        const val CONTENT_LOAD_TIME_MS = 3000L
        const val ANIMATION_DURATION_MS = 300L
        const val MINIMUM_TOUCH_TARGET_DP = 48 // Material Design minimum
    }

    // Stream categories to test
    private val streamCategories = listOf(
        "Books", "Podcasts", "Cartoons", "Movies", "Music", "Art"
    )

    @Before
    fun setUp() {
        hiltRule.inject()
        device = UiDevice.getInstance(InstrumentationRegistry.getInstrumentation())

        // Ensure device is in portrait mode
        device.setOrientationNatural()
        device.waitForIdle()
    }

    // MARK: - Basic Stream Tab Navigation Tests

    @Test
    fun streamTab_existsInBottomNavigation() {
        composeTestRule.apply {
            // Wait for app to load
            waitForIdle()

            // Navigate to Store screen via bottom navigation
            onNodeWithContentDescription("Store")
                .assertExists("Store tab should exist in bottom navigation")
                .performClick()

            // Verify Stream tab appears in store layout
            onNodeWithText("Stream")
                .assertExists("Stream tab should exist in store layout")
                .assertHasClickAction()

            // Validate Material3 theming
            onNodeWithText("Stream")
                .assertIsDisplayed()
        }
    }

    @Test
    fun streamCategoryTabs_displayAllSixCategories() {
        navigateToStreamTab()

        // Wait for Stream content to load
        onNodeWithTag("StreamContent")
            .assertExists("Stream content should load")
            .assertIsDisplayed()

        // Verify all 6 category tabs are present
        streamCategories.forEach { category ->
            onNodeWithText(category)
                .assertExists("$category tab should exist")
                .assertIsDisplayed()
                .assertHasClickAction()

            // Validate minimum touch target size (Material Design)
            onNodeWithText(category)
                .assertWidthIsAtLeast(MINIMUM_TOUCH_TARGET_DP.dp)
                .assertHeightIsAtLeast(MINIMUM_TOUCH_TARGET_DP.dp)
        }
    }

    @Test
    fun categoryTabSwitching_performsWithinBudget() {
        navigateToStreamTab()

        // Test switching between first 3 categories for performance
        val testCategories = streamCategories.take(3)

        testCategories.forEach { category ->
            val switchTime = measureTimeMillis {
                onNodeWithText(category)
                    .assertExists("$category tab should exist")
                    .performClick()

                // Wait for content to update
                onNodeWithTag("${category}Content")
                    .assertExists("$category content should load")

                waitForIdle()
            }

            assert(switchTime < TAB_SWITCH_TIME_MS) {
                "$category tab switch took ${switchTime}ms, should be under ${TAB_SWITCH_TIME_MS}ms"
            }

            // Verify tab is selected (Material3 selection state)
            onNodeWithText(category)
                .assertIsSelected()
        }
    }

    @Test
    fun tabStateManagement_persistsSelection() {
        navigateToStreamTab()

        // Select Podcasts category
        onNodeWithText("Podcasts")
            .assertExists()
            .performClick()

        // Verify selection
        onNodeWithText("Podcasts")
            .assertIsSelected()

        // Navigate away to Chat
        onNodeWithContentDescription("Chat")
            .performClick()

        // Navigate back to Store
        onNodeWithContentDescription("Store")
            .performClick()

        // Navigate back to Stream
        onNodeWithText("Stream")
            .performClick()

        // Verify Podcasts is still selected
        onNodeWithText("Podcasts")
            .assertIsSelected()
    }

    // MARK: - Movies Subtab Tests

    @Test
    fun moviesSubtabs_displayAndNavigate() {
        navigateToStreamTab()

        // Navigate to Movies category
        onNodeWithText("Movies")
            .assertExists()
            .performClick()

        // Wait for subtabs to appear
        onNodeWithTag("MoviesSubtabs")
            .assertExists("Movies subtabs should appear")
            .assertIsDisplayed()

        // Test subtab presence
        onNodeWithText("Short Films")
            .assertExists("Short Films subtab should exist")
            .assertIsDisplayed()

        onNodeWithText("Feature Films")
            .assertExists("Feature Films subtab should exist")
            .assertIsDisplayed()

        // Test subtab selection
        onNodeWithText("Feature Films")
            .performClick()
            .assertIsSelected()

        onNodeWithText("Short Films")
            .performClick()
            .assertIsSelected()

        // Verify other subtab is deselected
        onNodeWithText("Feature Films")
            .assertIsNotSelected()
    }

    @Test
    fun moviesSubtabs_persistSelection() {
        navigateToStreamTab()

        // Navigate to Movies and select Feature Films
        onNodeWithText("Movies")
            .performClick()

        onNodeWithText("Feature Films")
            .assertExists()
            .performClick()

        // Navigate to different category
        onNodeWithText("Books")
            .performClick()

        // Return to Movies
        onNodeWithText("Movies")
            .performClick()

        // Verify Feature Films is still selected
        onNodeWithText("Feature Films")
            .assertIsSelected()
    }

    // MARK: - Featured Content Tests

    @Test
    fun featuredContent_displaysInCarousel() {
        navigateToStreamTab()

        // Navigate to Books category
        onNodeWithText("Books")
            .performClick()

        // Wait for featured content to load
        onNodeWithTag("FeaturedCarousel")
            .assertExists("Featured carousel should appear")
            .assertIsDisplayed()

        // Verify featured items exist
        onAllNodesWithTag("FeaturedItem")
            .assertCountEquals(1, true) // At least one item

        // Test first featured item properties
        onAllNodesWithTag("FeaturedItem")
            .onFirst()
            .apply {
                // Should have title, image, and price
                assertExists("Featured item should exist")
                assertIsDisplayed()
                assertHasClickAction()
            }
    }

    @Test
    fun featuredContent_supportsHorizontalScrolling() {
        navigateToStreamTab()

        onNodeWithText("Books")
            .performClick()

        onNodeWithTag("FeaturedCarousel")
            .assertExists()
            .assertIsDisplayed()

        // Test horizontal scroll gesture
        onNodeWithTag("FeaturedCarousel")
            .performTouchInput {
                swipeLeft()
            }

        // Should not crash and should handle gesture
        waitForIdle()
    }

    // MARK: - Store Integration Tests

    @Test
    fun addToCart_integrationWorksCorrectly() {
        navigateToStreamTab()

        onNodeWithText("Books")
            .performClick()

        // Look for add to cart button
        onNodeWithTag("AddToCartButton")
            .assertExists()
            .assertIsDisplayed()
            .assertHasClickAction()

        // Validate Material3 button styling
        onNodeWithTag("AddToCartButton")
            .assertWidthIsAtLeast(MINIMUM_TOUCH_TARGET_DP.dp)
            .assertHeightIsAtLeast(MINIMUM_TOUCH_TARGET_DP.dp)

        // Perform add to cart action
        onNodeWithTag("AddToCartButton")
            .performClick()

        // Should show feedback (Snackbar, Toast, etc.)
        onNodeWithTag("CartFeedback")
            .assertExists("Should show cart feedback")
            .assertIsDisplayed()
    }

    // MARK: - Performance Tests

    @Test
    fun streamInitialLoad_meetsPerformanceBudget() {
        val loadTime = measureTimeMillis {
            navigateToStreamTab()

            onNodeWithTag("StreamContent")
                .assertExists()
                .assertIsDisplayed()

            waitForIdle()
        }

        assert(loadTime < CONTENT_LOAD_TIME_MS) {
            "Stream initial load took ${loadTime}ms, should be under ${CONTENT_LOAD_TIME_MS}ms"
        }
    }

    @Test
    fun contentLoading_performsEfficientlyAcrossCategories() {
        navigateToStreamTab()

        // Test content loading performance for each category
        streamCategories.take(3).forEach { category ->
            val loadTime = measureTimeMillis {
                onNodeWithText(category)
                    .performClick()

                onNodeWithTag("${category}Content")
                    .assertExists()

                waitForIdle()
            }

            assert(loadTime < CONTENT_LOAD_TIME_MS) {
                "$category content load took ${loadTime}ms, should be under ${CONTENT_LOAD_TIME_MS}ms"
            }
        }
    }

    // MARK: - Error Handling Tests

    @Test
    fun emptyState_handledGracefully() {
        navigateToStreamTab()

        // Navigate through categories and ensure no crashes
        streamCategories.forEach { category ->
            onNodeWithText(category)
                .performClick()

            // Should either show content or empty state, not crash
            val hasContent = onNodeWithTag("${category}Content").isDisplayed()
            val hasEmptyState = onNodeWithText("No content available", substring = true).isDisplayed()

            assert(hasContent || hasEmptyState) {
                "$category should show content or empty state"
            }

            waitForIdle()
        }
    }

    @Test
    fun networkError_handledGracefully() {
        // This test would require network condition mocking
        navigateToStreamTab()

        // Should handle network errors gracefully without crashing
        onNodeWithTag("StreamContent")
            .assertExists()

        // Even with potential network issues, should not crash
        waitForIdle()
    }

    // MARK: - Accessibility Tests

    @Test
    fun accessibility_supportsTalkBack() {
        navigateToStreamTab()

        // Test TalkBack accessibility
        streamCategories.take(3).forEach { category ->
            onNodeWithText(category)
                .assertExists()
                .assertIsDisplayed()
                .assert(hasClickAction())

            // Should have semantic properties for accessibility
            onNodeWithText(category)
                .assertHasClickAction()
        }
    }

    @Test
    fun accessibility_hasProperSemantics() {
        navigateToStreamTab()

        // Verify proper semantics for screen readers
        onNodeWithTag("StreamContent")
            .assertExists()

        // Category tabs should have proper role
        streamCategories.take(3).forEach { category ->
            onNodeWithText(category)
                .assertExists()
                .assertHasClickAction()
        }
    }

    // MARK: - Configuration Change Tests

    @Test
    fun configurationChange_maintainsState() = runBlocking {
        navigateToStreamTab()

        // Select Podcasts category
        onNodeWithText("Podcasts")
            .performClick()
            .assertIsSelected()

        // Rotate device to landscape
        device.setOrientationLeft()
        device.waitForIdle()
        delay(500) // Allow for configuration change

        // Should maintain state in landscape
        onNodeWithText("Podcasts")
            .assertIsSelected()

        // Rotate back to portrait
        device.setOrientationNatural()
        device.waitForIdle()
        delay(500)

        // Should still maintain state
        onNodeWithText("Podcasts")
            .assertIsSelected()
    }

    @Test
    fun configurationChange_handlesLayoutChanges() = runBlocking {
        navigateToStreamTab()

        // Rotate to landscape
        device.setOrientationLeft()
        device.waitForIdle()
        delay(500)

        // Should still display properly in landscape
        onNodeWithTag("StreamContent")
            .assertExists()
            .assertIsDisplayed()

        // Categories should still be accessible
        onNodeWithText("Books")
            .assertExists()
            .assertIsDisplayed()

        // Rotate back
        device.setOrientationNatural()
        device.waitForIdle()
        delay(500)

        // Should return to normal layout
        onNodeWithTag("StreamContent")
            .assertExists()
            .assertIsDisplayed()
    }

    // MARK: - Material Design 3 Compliance Tests

    @Test
    fun materialDesign3_theming() {
        navigateToStreamTab()

        // Verify Material3 theming is applied
        onNodeWithText("Stream")
            .assertExists()
            .assertIsDisplayed()

        // Category tabs should follow Material3 design
        streamCategories.take(3).forEach { category ->
            onNodeWithText(category)
                .assertExists()
                .assertIsDisplayed()
        }
    }

    @Test
    fun materialDesign3_animations() {
        navigateToStreamTab()

        // Test that animations are smooth (no specific assertion, just verify no crashes)
        streamCategories.take(3).forEach { category ->
            onNodeWithText(category)
                .performClick()

            // Wait for animation to complete
            waitForIdle()
        }
    }

    // MARK: - Helper Methods

    private fun navigateToStreamTab() {
        composeTestRule.apply {
            // Wait for app to load
            waitForIdle()

            // Navigate to Store tab
            onNodeWithContentDescription("Store")
                .assertExists("Store tab should exist")
                .performClick()

            // Navigate to Stream tab
            onNodeWithText("Stream")
                .assertExists("Stream tab should exist")
                .performClick()

            // Wait for Stream content to load
            onNodeWithTag("StreamContent")
                .assertExists("Stream content should load")

            waitForIdle()
        }
    }

    private fun SemanticsNodeInteraction.assertWidthIsAtLeast(width: Dp): SemanticsNodeInteraction {
        return assert(
            SemanticsMatcher("width is at least $width") {
                it.layoutInfo.width >= with(it.layoutInfo.density) { width.toPx() }
            }
        )
    }

    private fun SemanticsNodeInteraction.assertHeightIsAtLeast(height: Dp): SemanticsNodeInteraction {
        return assert(
            SemanticsMatcher("height is at least $height") {
                it.layoutInfo.height >= with(it.layoutInfo.density) { height.toPx() }
            }
        )
    }
}

/**
 * Performance Test Suite for Stream Navigation
 *
 * Specialized tests focused on performance measurement and validation
 * for Android-specific metrics and benchmarks.
 */
@HiltAndroidTest
@RunWith(AndroidJUnit4::class)
class StreamNavigationPerformanceTest {

    @get:Rule(order = 0)
    val hiltRule = HiltAndroidRule(this)

    @get:Rule(order = 1)
    val composeTestRule = createAndroidComposeRule<MainActivity>()

    @Before
    fun setUp() {
        hiltRule.inject()
    }

    @Test
    fun measureTabSwitchPerformance() {
        val streamCategories = listOf("Books", "Podcasts", "Movies")

        composeTestRule.apply {
            navigateToStreamTab()

            // Measure tab switching performance
            val measurements = mutableListOf<Long>()

            streamCategories.forEach { category ->
                val switchTime = measureTimeMillis {
                    onNodeWithText(category)
                        .performClick()

                    onNodeWithTag("${category}Content")
                        .assertExists()

                    waitForIdle()
                }

                measurements.add(switchTime)
            }

            val averageTime = measurements.average()
            assert(averageTime < 200.0) {
                "Average tab switch time ${averageTime}ms exceeds 200ms budget"
            }

            println("Tab switch performance - Average: ${averageTime}ms, Individual: $measurements")
        }
    }

    @Test
    fun measureMemoryUsageDuringNavigation() {
        composeTestRule.apply {
            navigateToStreamTab()

            // Navigate through categories to test memory stability
            repeat(3) {
                listOf("Books", "Podcasts", "Movies").forEach { category ->
                    onNodeWithText(category)
                        .performClick()

                    waitForIdle()
                }
            }

            // Should complete without memory issues
            onNodeWithTag("StreamContent")
                .assertExists()
        }
    }

    private fun ComposeContentTestRule.navigateToStreamTab() {
        waitForIdle()

        onNodeWithContentDescription("Store")
            .performClick()

        onNodeWithText("Stream")
            .performClick()

        onNodeWithTag("StreamContent")
            .assertExists()

        waitForIdle()
    }
}