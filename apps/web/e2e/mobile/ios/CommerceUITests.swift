/**
 * iOS Commerce E2E Tests
 * Comprehensive testing of commerce workflows on iOS using XCUITest
 */

import XCTest

class CommerceUITests: XCTestCase {
    var app: XCUIApplication!

    override func setUpWithError() throws {
        continueAfterFailure = false
        app = XCUIApplication()

        // Launch arguments for testing
        app.launchArguments = [
            "-UITesting",
            "-MockAPI",
            "-ResetUserDefaults"
        ]

        app.launch()

        // Wait for app to fully load
        let homeTabButton = app.tabBars.buttons["Store"]
        XCTAssertTrue(homeTabButton.waitForExistence(timeout: 10))
    }

    override func tearDownWithError() throws {
        // Clear cart after each test
        clearCart()
        app.terminate()
    }

    // MARK: - Helper Methods

    private func clearCart() {
        // Navigate to cart and clear it
        let cartTabButton = app.tabBars.buttons["Cart"]
        if cartTabButton.exists {
            cartTabButton.tap()

            let clearCartButton = app.buttons["ClearCartButton"]
            if clearCartButton.exists {
                clearCartButton.tap()

                let confirmButton = app.alerts.buttons["Clear Cart"]
                if confirmButton.exists {
                    confirmButton.tap()
                }
            }
        }
    }

    private func navigateToCategory(_ categoryName: String) {
        let storeTabButton = app.tabBars.buttons["Store"]
        storeTabButton.tap()

        let categoryCell = app.collectionViews.cells.containing(.staticText, identifier: categoryName).firstMatch
        XCTAssertTrue(categoryCell.waitForExistence(timeout: 5))
        categoryCell.tap()
    }

    private func addProductToCart(productIndex: Int = 0, quantity: Int = 1) {
        let productCells = app.collectionViews.cells.matching(identifier: "ProductCard")
        XCTAssertGreaterThan(productCells.count, productIndex)

        let productCell = productCells.element(boundBy: productIndex)
        productCell.tap()

        // Set quantity if needed
        if quantity > 1 {
            let quantitySelector = app.steppers["QuantitySelector"]
            for _ in 1..<quantity {
                quantitySelector.buttons["Increment"].tap()
            }
        }

        let addToCartButton = app.buttons["AddToCartButton"]
        XCTAssertTrue(addToCartButton.waitForExistence(timeout: 5))
        addToCartButton.tap()

        // Wait for add to cart confirmation
        let confirmationAlert = app.alerts["Added to Cart"]
        XCTAssertTrue(confirmationAlert.waitForExistence(timeout: 5))
        confirmationAlert.buttons["OK"].tap()

        // Navigate back to category
        let backButton = app.navigationBars.buttons.element(boundBy: 0)
        backButton.tap()
    }

    // MARK: - Test Cases

    func testCategoryBrowsing() throws {
        // Test browsing electronics category
        navigateToCategory("Electronics")

        // Verify category loaded
        let categoryTitle = app.navigationBars.staticTexts["Electronics"]
        XCTAssertTrue(categoryTitle.exists)

        // Verify products are displayed
        let productCells = app.collectionViews.cells.matching(identifier: "ProductCard")
        XCTAssertGreaterThan(productCells.count, 0)

        // Test product interaction
        let firstProduct = productCells.element(boundBy: 0)
        firstProduct.tap()

        // Verify product detail page
        let productTitle = app.staticTexts["ProductTitle"]
        XCTAssertTrue(productTitle.waitForExistence(timeout: 5))

        let productPrice = app.staticTexts["ProductPrice"]
        XCTAssertTrue(productPrice.exists)

        let addToCartButton = app.buttons["AddToCartButton"]
        XCTAssertTrue(addToCartButton.exists)
    }

    func testAddToCart() throws {
        // Navigate to electronics category
        navigateToCategory("Electronics")

        // Add product to cart
        addProductToCart(productIndex: 0, quantity: 1)

        // Navigate to cart
        let cartTabButton = app.tabBars.buttons["Cart"]
        cartTabButton.tap()

        // Verify item in cart
        let cartItems = app.tables.cells.matching(identifier: "CartItem")
        XCTAssertEqual(cartItems.count, 1)

        let firstCartItem = cartItems.element(boundBy: 0)
        XCTAssertTrue(firstCartItem.exists)

        // Verify cart total
        let cartTotal = app.staticTexts["CartTotal"]
        XCTAssertTrue(cartTotal.exists)
        XCTAssertTrue(cartTotal.label.contains("$"))
    }

    func testCartQuantityUpdate() throws {
        // Add product to cart
        navigateToCategory("Electronics")
        addProductToCart(productIndex: 0, quantity: 1)

        // Navigate to cart
        let cartTabButton = app.tabBars.buttons["Cart"]
        cartTabButton.tap()

        // Update quantity
        let quantityStepper = app.steppers["CartItemQuantityStepper"]
        XCTAssertTrue(quantityStepper.waitForExistence(timeout: 5))

        let incrementButton = quantityStepper.buttons["Increment"]
        incrementButton.tap()
        incrementButton.tap() // Quantity should now be 3

        // Verify quantity updated
        let quantityLabel = app.staticTexts["CartItemQuantity"]
        XCTAssertTrue(quantityLabel.waitForExistence(timeout: 2))
        XCTAssertEqual(quantityLabel.label, "3")

        // Verify total updated
        let cartTotal = app.staticTexts["CartTotal"]
        let updatedTotal = cartTotal.label
        XCTAssertTrue(updatedTotal.contains("$"))
    }

    func testRemoveFromCart() throws {
        // Add multiple products to cart
        navigateToCategory("Electronics")
        addProductToCart(productIndex: 0, quantity: 1)
        addProductToCart(productIndex: 1, quantity: 1)

        // Navigate to cart
        let cartTabButton = app.tabBars.buttons["Cart"]
        cartTabButton.tap()

        // Verify two items in cart
        let cartItems = app.tables.cells.matching(identifier: "CartItem")
        XCTAssertEqual(cartItems.count, 2)

        // Remove first item
        let firstCartItem = cartItems.element(boundBy: 0)
        firstCartItem.swipeLeft()

        let deleteButton = app.buttons["Delete"]
        XCTAssertTrue(deleteButton.waitForExistence(timeout: 3))
        deleteButton.tap()

        // Verify item removed
        XCTAssertEqual(cartItems.count, 1)
    }

    func testCouponApplication() throws {
        // Add product to cart
        navigateToCategory("Electronics")
        addProductToCart(productIndex: 0, quantity: 1)

        // Navigate to cart
        let cartTabButton = app.tabBars.buttons["Cart"]
        cartTabButton.tap()

        // Apply coupon
        let couponTextField = app.textFields["CouponTextField"]
        XCTAssertTrue(couponTextField.waitForExistence(timeout: 5))
        couponTextField.tap()
        couponTextField.typeText("SAVE20")

        let applyCouponButton = app.buttons["ApplyCouponButton"]
        applyCouponButton.tap()

        // Verify coupon applied
        let couponSuccess = app.staticTexts["CouponApplied"]
        XCTAssertTrue(couponSuccess.waitForExistence(timeout: 5))

        let discountAmount = app.staticTexts["DiscountAmount"]
        XCTAssertTrue(discountAmount.exists)
        XCTAssertTrue(discountAmount.label.contains("-$"))
    }

    func testInvalidCoupon() throws {
        // Add product to cart
        navigateToCategory("Electronics")
        addProductToCart(productIndex: 0, quantity: 1)

        // Navigate to cart
        let cartTabButton = app.tabBars.buttons["Cart"]
        cartTabButton.tap()

        // Apply invalid coupon
        let couponTextField = app.textFields["CouponTextField"]
        couponTextField.tap()
        couponTextField.typeText("INVALID123")

        let applyCouponButton = app.buttons["ApplyCouponButton"]
        applyCouponButton.tap()

        // Verify error message
        let errorAlert = app.alerts["Invalid Coupon"]
        XCTAssertTrue(errorAlert.waitForExistence(timeout: 5))
        errorAlert.buttons["OK"].tap()
    }

    func testCheckoutFlow() throws {
        // Add product to cart
        navigateToCategory("Electronics")
        addProductToCart(productIndex: 0, quantity: 1)

        // Navigate to cart and proceed to checkout
        let cartTabButton = app.tabBars.buttons["Cart"]
        cartTabButton.tap()

        let checkoutButton = app.buttons["CheckoutButton"]
        XCTAssertTrue(checkoutButton.waitForExistence(timeout: 5))
        checkoutButton.tap()

        // Fill shipping information
        let emailTextField = app.textFields["EmailTextField"]
        XCTAssertTrue(emailTextField.waitForExistence(timeout: 5))
        emailTextField.tap()
        emailTextField.typeText("test@example.com")

        let firstNameTextField = app.textFields["FirstNameTextField"]
        firstNameTextField.tap()
        firstNameTextField.typeText("John")

        let lastNameTextField = app.textFields["LastNameTextField"]
        lastNameTextField.tap()
        lastNameTextField.typeText("Doe")

        let addressTextField = app.textFields["AddressTextField"]
        addressTextField.tap()
        addressTextField.typeText("123 Test Street")

        let cityTextField = app.textFields["CityTextField"]
        cityTextField.tap()
        cityTextField.typeText("Test City")

        let zipCodeTextField = app.textFields["ZipCodeTextField"]
        zipCodeTextField.tap()
        zipCodeTextField.typeText("12345")

        // Continue to payment
        let continueButton = app.buttons["ContinueToPaymentButton"]
        continueButton.tap()

        // Fill payment information
        let cardNumberTextField = app.textFields["CardNumberTextField"]
        XCTAssertTrue(cardNumberTextField.waitForExistence(timeout: 5))
        cardNumberTextField.tap()
        cardNumberTextField.typeText("4111111111111111")

        let expiryTextField = app.textFields["ExpiryTextField"]
        expiryTextField.tap()
        expiryTextField.typeText("1225")

        let cvvTextField = app.textFields["CVVTextField"]
        cvvTextField.tap()
        cvvTextField.typeText("123")

        let cardNameTextField = app.textFields["CardNameTextField"]
        cardNameTextField.tap()
        cardNameTextField.typeText("John Doe")

        // Place order
        let placeOrderButton = app.buttons["PlaceOrderButton"]
        placeOrderButton.tap()

        // Verify order confirmation
        let orderConfirmation = app.staticTexts["OrderConfirmation"]
        XCTAssertTrue(orderConfirmation.waitForExistence(timeout: 10))

        let orderNumber = app.staticTexts["OrderNumber"]
        XCTAssertTrue(orderNumber.exists)
    }

    func testProductSearch() throws {
        // Navigate to store
        let storeTabButton = app.tabBars.buttons["Store"]
        storeTabButton.tap()

        // Tap search bar
        let searchBar = app.searchFields["ProductSearchBar"]
        XCTAssertTrue(searchBar.waitForExistence(timeout: 5))
        searchBar.tap()
        searchBar.typeText("smartphone")

        // Tap search button on keyboard
        app.keyboards.buttons["Search"].tap()

        // Verify search results
        let searchResults = app.collectionViews.cells.matching(identifier: "ProductCard")
        XCTAssertGreaterThan(searchResults.count, 0)

        // Verify search results contain search term
        let firstProduct = searchResults.element(boundBy: 0)
        firstProduct.tap()

        let productTitle = app.staticTexts["ProductTitle"]
        XCTAssertTrue(productTitle.waitForExistence(timeout: 5))

        let titleText = productTitle.label.lowercased()
        XCTAssertTrue(titleText.contains("smartphone") || titleText.contains("phone"))
    }

    func testCartPersistence() throws {
        // Add product to cart
        navigateToCategory("Electronics")
        addProductToCart(productIndex: 0, quantity: 2)

        // Navigate to cart and verify items
        let cartTabButton = app.tabBars.buttons["Cart"]
        cartTabButton.tap()

        let cartItems = app.tables.cells.matching(identifier: "CartItem")
        XCTAssertEqual(cartItems.count, 1)

        let quantityLabel = app.staticTexts["CartItemQuantity"]
        XCTAssertEqual(quantityLabel.label, "2")

        // Terminate and relaunch app
        app.terminate()
        app.launch()

        // Wait for app to load
        let homeTabButton = app.tabBars.buttons["Store"]
        XCTAssertTrue(homeTabButton.waitForExistence(timeout: 10))

        // Navigate to cart and verify persistence
        cartTabButton.tap()

        let persistedCartItems = app.tables.cells.matching(identifier: "CartItem")
        XCTAssertEqual(persistedCartItems.count, 1)

        let persistedQuantityLabel = app.staticTexts["CartItemQuantity"]
        XCTAssertEqual(persistedQuantityLabel.label, "2")
    }

    func testOfflineMode() throws {
        // Add product to cart while online
        navigateToCategory("Electronics")
        addProductToCart(productIndex: 0, quantity: 1)

        // Simulate offline mode
        app.launchArguments.append("-OfflineMode")
        app.terminate()
        app.launch()

        // Wait for app to load
        let homeTabButton = app.tabBars.buttons["Store"]
        XCTAssertTrue(homeTabButton.waitForExistence(timeout: 10))

        // Verify cart still accessible
        let cartTabButton = app.tabBars.buttons["Cart"]
        cartTabButton.tap()

        let cartItems = app.tables.cells.matching(identifier: "CartItem")
        XCTAssertEqual(cartItems.count, 1)

        // Verify offline indicator
        let offlineIndicator = app.staticTexts["OfflineIndicator"]
        XCTAssertTrue(offlineIndicator.waitForExistence(timeout: 5))

        // Try to checkout (should show offline message)
        let checkoutButton = app.buttons["CheckoutButton"]
        checkoutButton.tap()

        let offlineAlert = app.alerts["Offline Mode"]
        XCTAssertTrue(offlineAlert.waitForExistence(timeout: 5))
        offlineAlert.buttons["OK"].tap()
    }

    func testAccessibility() throws {
        // Navigate to store
        let storeTabButton = app.tabBars.buttons["Store"]
        storeTabButton.tap()

        // Test VoiceOver accessibility
        navigateToCategory("Electronics")

        let productCells = app.collectionViews.cells.matching(identifier: "ProductCard")
        let firstProduct = productCells.element(boundBy: 0)

        // Verify accessibility labels exist
        XCTAssertNotNil(firstProduct.label)
        XCTAssertNotEqual(firstProduct.label, "")

        // Test cart accessibility
        addProductToCart(productIndex: 0, quantity: 1)

        let cartTabButton = app.tabBars.buttons["Cart"]
        cartTabButton.tap()

        let cartItems = app.tables.cells.matching(identifier: "CartItem")
        let firstCartItem = cartItems.element(boundBy: 0)

        XCTAssertNotNil(firstCartItem.label)
        XCTAssertNotEqual(firstCartItem.label, "")

        // Test button accessibility
        let checkoutButton = app.buttons["CheckoutButton"]
        XCTAssertNotNil(checkoutButton.label)
        XCTAssertTrue(checkoutButton.isHittable)
    }

    func testDeepLinking() throws {
        // Test deep link to specific product
        let productDeepLink = "tchat://product/electronics-smartphone-001"

        // Simulate opening deep link (in real test, this would be done externally)
        app.terminate()
        app.launchArguments.append("-DeepLinkURL")
        app.launchArguments.append(productDeepLink)
        app.launch()

        // Verify product page opened directly
        let productTitle = app.staticTexts["ProductTitle"]
        XCTAssertTrue(productTitle.waitForExistence(timeout: 10))

        let addToCartButton = app.buttons["AddToCartButton"]
        XCTAssertTrue(addToCartButton.exists)

        // Test deep link to cart
        let cartDeepLink = "tchat://cart"
        app.terminate()
        app.launchArguments = ["-DeepLinkURL", cartDeepLink]
        app.launch()

        // Verify cart opened directly
        let cartTitle = app.navigationBars.staticTexts["Cart"]
        XCTAssertTrue(cartTitle.waitForExistence(timeout: 10))
    }
}

// MARK: - Performance Tests

class CommercePerformanceTests: XCTestCase {
    var app: XCUIApplication!

    override func setUpWithError() throws {
        continueAfterFailure = false
        app = XCUIApplication()
        app.launchArguments = ["-UITesting", "-MockAPI"]
        app.launch()
    }

    func testAppLaunchPerformance() throws {
        measure(metrics: [XCTApplicationLaunchMetric()]) {
            app.terminate()
            app.launch()

            let storeTabButton = app.tabBars.buttons["Store"]
            XCTAssertTrue(storeTabButton.waitForExistence(timeout: 10))
        }
    }

    func testCategoryLoadingPerformance() throws {
        let storeTabButton = app.tabBars.buttons["Store"]
        storeTabButton.tap()

        measure(metrics: [XCTClockMetric()]) {
            let categoryCell = app.collectionViews.cells.containing(.staticText, identifier: "Electronics").firstMatch
            categoryCell.tap()

            let productCells = app.collectionViews.cells.matching(identifier: "ProductCard")
            XCTAssertTrue(productCells.firstMatch.waitForExistence(timeout: 5))

            // Navigate back
            let backButton = app.navigationBars.buttons.element(boundBy: 0)
            backButton.tap()
        }
    }

    func testScrollingPerformance() throws {
        // Navigate to category with many products
        let storeTabButton = app.tabBars.buttons["Store"]
        storeTabButton.tap()

        let categoryCell = app.collectionViews.cells.containing(.staticText, identifier: "Electronics").firstMatch
        categoryCell.tap()

        let collectionView = app.collectionViews.firstMatch

        measure(metrics: [XCTClockMetric()]) {
            // Scroll through products
            for _ in 0..<10 {
                collectionView.swipeUp()
            }

            for _ in 0..<10 {
                collectionView.swipeDown()
            }
        }
    }
}