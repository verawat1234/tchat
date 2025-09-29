/**
 * Android Commerce E2E Tests
 * Comprehensive testing of commerce workflows on Android using Espresso
 */

package com.tchat.app.e2e

import androidx.compose.ui.test.*
import androidx.compose.ui.test.junit4.createComposeRule
import androidx.test.espresso.Espresso.*
import androidx.test.espresso.action.ViewActions.*
import androidx.test.espresso.assertion.ViewAssertions.*
import androidx.test.espresso.matcher.ViewMatchers.*
import androidx.test.ext.junit.rules.ActivityScenarioRule
import androidx.test.ext.junit.runners.AndroidJUnit4
import androidx.test.platform.app.InstrumentationRegistry
import androidx.test.uiautomator.UiDevice
import com.tchat.app.MainActivity
import com.tchat.app.R
import org.hamcrest.Matchers.*
import org.junit.Before
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith

@RunWith(AndroidJUnit4::class)
class CommerceE2ETest {

    @get:Rule
    val composeTestRule = createComposeRule()

    @get:Rule
    val activityRule = ActivityScenarioRule(MainActivity::class.java)

    private lateinit var device: UiDevice

    @Before
    fun setUp() {
        device = UiDevice.getInstance(InstrumentationRegistry.getInstrumentation())

        // Set up test environment
        activityRule.scenario.onActivity { activity ->
            activity.intent.putExtra("test_mode", true)
            activity.intent.putExtra("mock_api", true)
            activity.intent.putExtra("reset_data", true)
        }

        // Wait for app to load
        composeTestRule.waitForIdle()
        Thread.sleep(2000) // Additional wait for network/initialization
    }

    @Test
    fun testCategoryBrowsing() {
        // Navigate to Store tab
        composeTestRule.onNodeWithTag("StoreTab").performClick()

        // Wait for categories to load
        composeTestRule.waitForIdle()

        // Tap on Electronics category
        composeTestRule.onNodeWithTag("CategoryCard_Electronics").performClick()

        // Verify category page loaded
        composeTestRule.onNodeWithTag("CategoryTitle").assertTextContains("Electronics")

        // Verify products are displayed
        composeTestRule.onNodeWithTag("ProductGrid").assertExists()
        composeTestRule.onAllNodesWithTag("ProductCard").assertCountGreaterThan(0)

        // Test product interaction
        composeTestRule.onAllNodesWithTag("ProductCard")[0].performClick()

        // Verify product detail page
        composeTestRule.onNodeWithTag("ProductTitle").assertExists()
        composeTestRule.onNodeWithTag("ProductPrice").assertExists()
        composeTestRule.onNodeWithTag("AddToCartButton").assertExists()
    }

    @Test
    fun testAddToCart() {
        // Navigate to Electronics category
        navigateToCategory("Electronics")

        // Add first product to cart
        composeTestRule.onAllNodesWithTag("ProductCard")[0].performClick()

        // Set quantity to 2
        composeTestRule.onNodeWithTag("QuantitySelector").performClick()
        composeTestRule.onNodeWithTag("QuantityIncrement").performClick()

        // Add to cart
        composeTestRule.onNodeWithTag("AddToCartButton").performClick()

        // Verify add to cart confirmation
        composeTestRule.onNodeWithTag("AddToCartConfirmation").assertExists()
        composeTestRule.onNodeWithText("Added to Cart").assertExists()

        // Close confirmation
        composeTestRule.onNodeWithTag("ConfirmationOkButton").performClick()

        // Navigate to cart
        composeTestRule.onNodeWithTag("CartTab").performClick()

        // Verify item in cart
        composeTestRule.onNodeWithTag("CartItemList").assertExists()
        composeTestRule.onAllNodesWithTag("CartItem").assertCountEquals(1)

        // Verify quantity
        composeTestRule.onNodeWithTag("CartItemQuantity").assertTextEquals("2")

        // Verify cart total exists
        composeTestRule.onNodeWithTag("CartTotal").assertExists()
    }

    @Test
    fun testCartQuantityUpdate() {
        // Add product to cart first
        addProductToCart(productIndex = 0, quantity = 1)

        // Navigate to cart
        composeTestRule.onNodeWithTag("CartTab").performClick()

        // Increase quantity
        composeTestRule.onNodeWithTag("CartItemQuantityIncrement").performClick()
        composeTestRule.onNodeWithTag("CartItemQuantityIncrement").performClick()

        // Verify quantity updated to 3
        composeTestRule.waitForIdle()
        composeTestRule.onNodeWithTag("CartItemQuantity").assertTextEquals("3")

        // Verify total updated
        composeTestRule.onNodeWithTag("CartTotal").assertExists()

        // Decrease quantity
        composeTestRule.onNodeWithTag("CartItemQuantityDecrement").performClick()

        // Verify quantity decreased to 2
        composeTestRule.waitForIdle()
        composeTestRule.onNodeWithTag("CartItemQuantity").assertTextEquals("2")
    }

    @Test
    fun testRemoveFromCart() {
        // Add multiple products to cart
        addProductToCart(productIndex = 0, quantity = 1)
        addProductToCart(productIndex = 1, quantity = 1)

        // Navigate to cart
        composeTestRule.onNodeWithTag("CartTab").performClick()

        // Verify two items in cart
        composeTestRule.onAllNodesWithTag("CartItem").assertCountEquals(2)

        // Remove first item
        composeTestRule.onAllNodesWithTag("RemoveItemButton")[0].performClick()

        // Confirm removal
        composeTestRule.onNodeWithTag("ConfirmRemovalButton").performClick()

        // Verify item removed
        composeTestRule.waitForIdle()
        composeTestRule.onAllNodesWithTag("CartItem").assertCountEquals(1)
    }

    @Test
    fun testClearCart() {
        // Add multiple products to cart
        addProductToCart(productIndex = 0, quantity = 1)
        addProductToCart(productIndex = 1, quantity = 1)

        // Navigate to cart
        composeTestRule.onNodeWithTag("CartTab").performClick()

        // Clear cart
        composeTestRule.onNodeWithTag("ClearCartButton").performClick()

        // Confirm clear cart
        composeTestRule.onNodeWithTag("ConfirmClearCartButton").performClick()

        // Verify cart is empty
        composeTestRule.waitForIdle()
        composeTestRule.onNodeWithTag("EmptyCartMessage").assertExists()
        composeTestRule.onNodeWithText("Your cart is empty").assertExists()
    }

    @Test
    fun testCouponApplication() {
        // Add product to cart
        addProductToCart(productIndex = 0, quantity = 1)

        // Navigate to cart
        composeTestRule.onNodeWithTag("CartTab").performClick()

        // Apply valid coupon
        composeTestRule.onNodeWithTag("CouponTextField").performTextInput("SAVE20")
        composeTestRule.onNodeWithTag("ApplyCouponButton").performClick()

        // Verify coupon applied
        composeTestRule.waitForIdle()
        composeTestRule.onNodeWithTag("CouponAppliedIndicator").assertExists()
        composeTestRule.onNodeWithTag("DiscountAmount").assertExists()

        // Verify discount in total
        composeTestRule.onNodeWithTag("CartSubtotal").assertExists()
        composeTestRule.onNodeWithTag("CartDiscount").assertExists()
        composeTestRule.onNodeWithTag("CartTotal").assertExists()
    }

    @Test
    fun testInvalidCoupon() {
        // Add product to cart
        addProductToCart(productIndex = 0, quantity = 1)

        // Navigate to cart
        composeTestRule.onNodeWithTag("CartTab").performClick()

        // Apply invalid coupon
        composeTestRule.onNodeWithTag("CouponTextField").performTextInput("INVALID123")
        composeTestRule.onNodeWithTag("ApplyCouponButton").performClick()

        // Verify error message
        composeTestRule.waitForIdle()
        composeTestRule.onNodeWithTag("CouponErrorMessage").assertExists()
        composeTestRule.onNodeWithText("Invalid coupon code").assertExists()
    }

    @Test
    fun testProductSearch() {
        // Navigate to Store tab
        composeTestRule.onNodeWithTag("StoreTab").performClick()

        // Tap search bar
        composeTestRule.onNodeWithTag("ProductSearchBar").performClick()

        // Type search query
        composeTestRule.onNodeWithTag("ProductSearchBar").performTextInput("smartphone")

        // Tap search button
        composeTestRule.onNodeWithTag("SearchButton").performClick()

        // Verify search results
        composeTestRule.waitForIdle()
        composeTestRule.onNodeWithTag("SearchResults").assertExists()
        composeTestRule.onAllNodesWithTag("ProductCard").assertCountGreaterThan(0)

        // Verify search results relevance
        composeTestRule.onAllNodesWithTag("ProductCard")[0].performClick()
        composeTestRule.onNodeWithTag("ProductTitle").assertExists()

        // The product title should contain search term or related terms
        val productTitle = composeTestRule.onNodeWithTag("ProductTitle")
        // In a real test, we would assert the title contains "smartphone" or "phone"
    }

    @Test
    fun testCheckoutFlow() {
        // Add product to cart
        addProductToCart(productIndex = 0, quantity = 1)

        // Navigate to cart and proceed to checkout
        composeTestRule.onNodeWithTag("CartTab").performClick()
        composeTestRule.onNodeWithTag("CheckoutButton").performClick()

        // Fill shipping information
        composeTestRule.onNodeWithTag("EmailTextField").performTextInput("test@example.com")
        composeTestRule.onNodeWithTag("FirstNameTextField").performTextInput("John")
        composeTestRule.onNodeWithTag("LastNameTextField").performTextInput("Doe")
        composeTestRule.onNodeWithTag("AddressTextField").performTextInput("123 Test Street")
        composeTestRule.onNodeWithTag("CityTextField").performTextInput("Test City")
        composeTestRule.onNodeWithTag("StateTextField").performTextInput("CA")
        composeTestRule.onNodeWithTag("ZipCodeTextField").performTextInput("12345")

        // Continue to payment
        composeTestRule.onNodeWithTag("ContinueToPaymentButton").performClick()

        // Fill payment information
        composeTestRule.waitForIdle()
        composeTestRule.onNodeWithTag("CardNumberTextField").performTextInput("4111111111111111")
        composeTestRule.onNodeWithTag("ExpiryMonthTextField").performTextInput("12")
        composeTestRule.onNodeWithTag("ExpiryYearTextField").performTextInput("2025")
        composeTestRule.onNodeWithTag("CVVTextField").performTextInput("123")
        composeTestRule.onNodeWithTag("CardNameTextField").performTextInput("John Doe")

        // Place order
        composeTestRule.onNodeWithTag("PlaceOrderButton").performClick()

        // Verify order confirmation
        composeTestRule.waitForIdle()
        composeTestRule.onNodeWithTag("OrderConfirmation").assertExists()
        composeTestRule.onNodeWithTag("OrderNumber").assertExists()
        composeTestRule.onNodeWithText("Order Confirmed").assertExists()
    }

    @Test
    fun testCartPersistence() {
        // Add product to cart
        addProductToCart(productIndex = 0, quantity = 2)

        // Verify cart has item
        composeTestRule.onNodeWithTag("CartTab").performClick()
        composeTestRule.onAllNodesWithTag("CartItem").assertCountEquals(1)
        composeTestRule.onNodeWithTag("CartItemQuantity").assertTextEquals("2")

        // Simulate app restart by recreating activity
        activityRule.scenario.recreate()
        composeTestRule.waitForIdle()
        Thread.sleep(2000)

        // Navigate to cart and verify persistence
        composeTestRule.onNodeWithTag("CartTab").performClick()
        composeTestRule.onAllNodesWithTag("CartItem").assertCountEquals(1)
        composeTestRule.onNodeWithTag("CartItemQuantity").assertTextEquals("2")
    }

    @Test
    fun testOfflineMode() {
        // Add product to cart while online
        addProductToCart(productIndex = 0, quantity = 1)

        // Simulate offline mode
        activityRule.scenario.onActivity { activity ->
            activity.intent.putExtra("offline_mode", true)
        }

        activityRule.scenario.recreate()
        composeTestRule.waitForIdle()

        // Verify offline indicator
        composeTestRule.onNodeWithTag("OfflineIndicator").assertExists()

        // Verify cart still accessible
        composeTestRule.onNodeWithTag("CartTab").performClick()
        composeTestRule.onAllNodesWithTag("CartItem").assertCountEquals(1)

        // Try to checkout (should show offline message)
        composeTestRule.onNodeWithTag("CheckoutButton").performClick()
        composeTestRule.onNodeWithTag("OfflineDialog").assertExists()
        composeTestRule.onNodeWithText("You are currently offline").assertExists()
    }

    @Test
    fun testAccessibility() {
        // Test semantic properties of main components
        composeTestRule.onNodeWithTag("StoreTab").assertHasClickAction()

        // Navigate to store
        composeTestRule.onNodeWithTag("StoreTab").performClick()

        // Test category accessibility
        navigateToCategory("Electronics")

        // Test product card accessibility
        val productCard = composeTestRule.onAllNodesWithTag("ProductCard")[0]
        productCard.assertHasClickAction()

        productCard.performClick()

        // Test product detail accessibility
        composeTestRule.onNodeWithTag("ProductTitle").assertExists()
        composeTestRule.onNodeWithTag("ProductPrice").assertExists()
        composeTestRule.onNodeWithTag("AddToCartButton").assertHasClickAction()

        // Test cart accessibility
        addToCartAndNavigateToCart()

        composeTestRule.onNodeWithTag("CartItemQuantityIncrement").assertHasClickAction()
        composeTestRule.onNodeWithTag("CartItemQuantityDecrement").assertHasClickAction()
        composeTestRule.onNodeWithTag("RemoveItemButton").assertHasClickAction()
    }

    @Test
    fun testDeepLinking() {
        // Test deep link to specific product
        val productIntent = android.content.Intent().apply {
            action = android.content.Intent.ACTION_VIEW
            data = android.net.Uri.parse("tchat://product/electronics-smartphone-001")
        }

        activityRule.scenario.onActivity { activity ->
            activity.onNewIntent(productIntent)
        }

        // Verify product page opened directly
        composeTestRule.waitForIdle()
        composeTestRule.onNodeWithTag("ProductTitle").assertExists()
        composeTestRule.onNodeWithTag("AddToCartButton").assertExists()

        // Test deep link to cart
        val cartIntent = android.content.Intent().apply {
            action = android.content.Intent.ACTION_VIEW
            data = android.net.Uri.parse("tchat://cart")
        }

        activityRule.scenario.onActivity { activity ->
            activity.onNewIntent(cartIntent)
        }

        // Verify cart opened directly
        composeTestRule.waitForIdle()
        composeTestRule.onNodeWithTag("CartContainer").assertExists()
    }

    @Test
    fun testResponsiveLayout() {
        // Test different screen orientations
        device.setOrientationLandscape()
        composeTestRule.waitForIdle()

        // Verify layout adapts to landscape
        composeTestRule.onNodeWithTag("StoreTab").performClick()
        navigateToCategory("Electronics")

        // In landscape, should show more products per row
        composeTestRule.onNodeWithTag("ProductGrid").assertExists()

        // Rotate back to portrait
        device.setOrientationPortrait()
        composeTestRule.waitForIdle()

        // Verify layout adapts back to portrait
        composeTestRule.onNodeWithTag("ProductGrid").assertExists()
    }

    // MARK: - Helper Methods

    private fun navigateToCategory(categoryName: String) {
        composeTestRule.onNodeWithTag("StoreTab").performClick()
        composeTestRule.waitForIdle()
        composeTestRule.onNodeWithTag("CategoryCard_$categoryName").performClick()
        composeTestRule.waitForIdle()
    }

    private fun addProductToCart(productIndex: Int = 0, quantity: Int = 1) {
        navigateToCategory("Electronics")

        // Tap product
        composeTestRule.onAllNodesWithTag("ProductCard")[productIndex].performClick()

        // Set quantity if greater than 1
        if (quantity > 1) {
            repeat(quantity - 1) {
                composeTestRule.onNodeWithTag("QuantityIncrement").performClick()
            }
        }

        // Add to cart
        composeTestRule.onNodeWithTag("AddToCartButton").performClick()

        // Close confirmation
        composeTestRule.onNodeWithTag("ConfirmationOkButton").performClick()

        // Navigate back
        composeTestRule.onNodeWithTag("BackButton").performClick()
    }

    private fun addToCartAndNavigateToCart() {
        addProductToCart(productIndex = 0, quantity = 1)
        composeTestRule.onNodeWithTag("CartTab").performClick()
    }

    private fun ComposeContentTestRule.onAllNodesWithTag(tag: String) =
        onAllNodes(hasTestTag(tag))

    private fun SemanticsNodeInteractionCollection.assertCountGreaterThan(count: Int) {
        assert(fetchSemanticsNodes().size > count) {
            "Expected more than $count nodes, but found ${fetchSemanticsNodes().size}"
        }
    }

    private fun SemanticsNodeInteractionCollection.assertCountEquals(count: Int) {
        assert(fetchSemanticsNodes().size == count) {
            "Expected $count nodes, but found ${fetchSemanticsNodes().size}"
        }
    }
}

// MARK: - Performance Tests

@RunWith(AndroidJUnit4::class)
class CommercePerformanceTest {

    @get:Rule
    val composeTestRule = createComposeRule()

    @get:Rule
    val activityRule = ActivityScenarioRule(MainActivity::class.java)

    @Test
    fun testAppStartupPerformance() {
        val startTime = System.currentTimeMillis()

        // Wait for app to fully load
        composeTestRule.onNodeWithTag("StoreTab").assertExists()
        composeTestRule.waitForIdle()

        val loadTime = System.currentTimeMillis() - startTime

        // Assert app loads within acceptable time (2 seconds)
        assert(loadTime < 2000) {
            "App took too long to load: ${loadTime}ms"
        }
    }

    @Test
    fun testCategoryLoadingPerformance() {
        composeTestRule.onNodeWithTag("StoreTab").performClick()

        val startTime = System.currentTimeMillis()

        // Navigate to category
        composeTestRule.onNodeWithTag("CategoryCard_Electronics").performClick()

        // Wait for products to load
        composeTestRule.onAllNodesWithTag("ProductCard")[0].assertExists()

        val loadTime = System.currentTimeMillis() - startTime

        // Assert category loads within acceptable time (1 second)
        assert(loadTime < 1000) {
            "Category took too long to load: ${loadTime}ms"
        }
    }

    @Test
    fun testScrollPerformance() {
        // Navigate to category with many products
        composeTestRule.onNodeWithTag("StoreTab").performClick()
        composeTestRule.onNodeWithTag("CategoryCard_Electronics").performClick()

        val productGrid = composeTestRule.onNodeWithTag("ProductGrid")

        val startTime = System.currentTimeMillis()

        // Perform multiple scroll gestures
        repeat(10) {
            productGrid.performScrollToIndex(it * 5)
            composeTestRule.waitForIdle()
        }

        val scrollTime = System.currentTimeMillis() - startTime

        // Assert scrolling performance is acceptable
        assert(scrollTime < 3000) {
            "Scrolling took too long: ${scrollTime}ms"
        }
    }

    @Test
    fun testCartUpdatePerformance() {
        // Add product to cart
        composeTestRule.onNodeWithTag("StoreTab").performClick()
        composeTestRule.onNodeWithTag("CategoryCard_Electronics").performClick()
        composeTestRule.onAllNodesWithTag("ProductCard")[0].performClick()
        composeTestRule.onNodeWithTag("AddToCartButton").performClick()
        composeTestRule.onNodeWithTag("ConfirmationOkButton").performClick()

        // Navigate to cart
        composeTestRule.onNodeWithTag("CartTab").performClick()

        val startTime = System.currentTimeMillis()

        // Update quantity multiple times
        repeat(5) {
            composeTestRule.onNodeWithTag("CartItemQuantityIncrement").performClick()
            composeTestRule.waitForIdle()
        }

        val updateTime = System.currentTimeMillis() - startTime

        // Assert cart updates are fast
        assert(updateTime < 1000) {
            "Cart updates took too long: ${updateTime}ms"
        }
    }

    private fun ComposeContentTestRule.onAllNodesWithTag(tag: String) =
        onAllNodes(hasTestTag(tag))
}