/**
 * Shared Test Data for Mobile E2E Tests
 * Common test data and utilities used across iOS and Android tests
 */

package com.tchat.app.shared

object TestData {
    // Test Users
    object Users {
        val newUser = User(
            email = "new-user@test.com",
            password = "TestPassword123!",
            firstName = "New",
            lastName = "User"
        )

        val existingUser = User(
            email = "existing-user@test.com",
            password = "TestPassword123!",
            firstName = "Existing",
            lastName = "User"
        )

        val premiumUser = User(
            email = "premium-user@test.com",
            password = "TestPassword123!",
            firstName = "Premium",
            lastName = "User"
        )
    }

    // Test Products
    object Products {
        val electronics = Product(
            id = "electronics-smartphone-001",
            name = "Test Smartphone",
            category = "Electronics",
            price = 299.99,
            sku = "PHONE001"
        )

        val clothing = Product(
            id = "clothing-tshirt-001",
            name = "Test T-Shirt",
            category = "Clothing",
            price = 29.99,
            sku = "SHIRT001"
        )

        val books = Product(
            id = "books-novel-001",
            name = "Test Book",
            category = "Books",
            price = 19.99,
            sku = "BOOK001"
        )

        val outOfStock = Product(
            id = "test-out-of-stock",
            name = "Out of Stock Item",
            category = "Test",
            price = 49.99,
            sku = "STOCK000",
            stock = 0
        )
    }

    // Test Coupons
    object Coupons {
        val percentage = Coupon(
            code = "SAVE20",
            type = CouponType.PERCENTAGE,
            value = 20,
            minAmount = 100,
            maxDiscount = 50
        )

        val fixed = Coupon(
            code = "FIXED10",
            type = CouponType.FIXED,
            value = 10,
            minAmount = 50
        )

        val expired = Coupon(
            code = "EXPIRED",
            type = CouponType.PERCENTAGE,
            value = 15,
            expiresAt = "2020-01-01"
        )
    }

    // Test Addresses
    object Addresses {
        val defaultShipping = Address(
            street = "123 Test Street",
            city = "Test City",
            state = "CA",
            zipCode = "12345",
            country = "US"
        )

        val alternativeShipping = Address(
            street = "456 Alternative Ave",
            city = "Alt City",
            state = "NY",
            zipCode = "67890",
            country = "US"
        )
    }

    // Test Payment Methods
    object PaymentMethods {
        val validCreditCard = PaymentMethod(
            type = PaymentType.CREDIT_CARD,
            cardNumber = "4111111111111111",
            expiryMonth = "12",
            expiryYear = "2025",
            cvv = "123",
            name = "John Doe"
        )

        val invalidCreditCard = PaymentMethod(
            type = PaymentType.CREDIT_CARD,
            cardNumber = "1234567890123456",
            expiryMonth = "12",
            expiryYear = "2025",
            cvv = "123",
            name = "Test User"
        )

        val declinedCreditCard = PaymentMethod(
            type = PaymentType.CREDIT_CARD,
            cardNumber = "4000000000000002",
            expiryMonth = "12",
            expiryYear = "2025",
            cvv = "123",
            name = "Declined Card"
        )
    }

    // Deep Link URLs
    object DeepLinks {
        const val HOME = "tchat://home"
        const val STORE = "tchat://store"
        const val CART = "tchat://cart"
        const val ELECTRONICS_CATEGORY = "tchat://category/electronics"
        const val SMARTPHONE_PRODUCT = "tchat://product/electronics-smartphone-001"
        const val CHECKOUT = "tchat://checkout"
        const val PROFILE = "tchat://profile"
    }

    // Test Scenarios
    object Scenarios {
        val singleItemCart = CartScenario(
            products = listOf(
                CartItem(Products.electronics, quantity = 1)
            ),
            expectedTotal = 299.99
        )

        val multipleItemCart = CartScenario(
            products = listOf(
                CartItem(Products.electronics, quantity = 1),
                CartItem(Products.clothing, quantity = 2),
                CartItem(Products.books, quantity = 1)
            ),
            expectedTotal = 379.97 // 299.99 + (29.99 * 2) + 19.99
        )

        val cartWithCoupon = CartScenario(
            products = listOf(
                CartItem(Products.electronics, quantity = 1)
            ),
            coupon = Coupons.percentage,
            expectedSubtotal = 299.99,
            expectedDiscount = 50.0, // 20% of 299.99, capped at 50
            expectedTotal = 249.99
        )

        val freeShippingCart = CartScenario(
            products = listOf(
                CartItem(Products.electronics, quantity = 2)
            ),
            expectedTotal = 599.98,
            freeShipping = true
        )
    }

    // Performance Thresholds
    object PerformanceThresholds {
        const val APP_LAUNCH_TIME_MS = 2000
        const val CATEGORY_LOAD_TIME_MS = 1000
        const val PRODUCT_LOAD_TIME_MS = 800
        const val CART_UPDATE_TIME_MS = 500
        const val SEARCH_RESPONSE_TIME_MS = 1200
        const val CHECKOUT_PROCESS_TIME_MS = 5000
    }

    // Error Messages
    object ErrorMessages {
        const val INVALID_COUPON = "Invalid coupon code"
        const val EXPIRED_COUPON = "Coupon has expired"
        const val MINIMUM_ORDER_NOT_MET = "Minimum order amount not met"
        const val ITEM_OUT_OF_STOCK = "Item is out of stock"
        const val PAYMENT_DECLINED = "Payment was declined"
        const val NETWORK_ERROR = "Network connection error"
        const val OFFLINE_MODE = "You are currently offline"
        const val SESSION_EXPIRED = "Your session has expired"
        const val INVALID_EMAIL = "Please enter a valid email address"
        const val REQUIRED_FIELD = "This field is required"
        const val INVALID_CARD_NUMBER = "Invalid card number"
        const val CARD_EXPIRED = "Card has expired"
        const val INSUFFICIENT_FUNDS = "Insufficient funds"
    }

    // Success Messages
    object SuccessMessages {
        const val ADDED_TO_CART = "Added to Cart"
        const val COUPON_APPLIED = "Coupon applied successfully"
        const val ORDER_CONFIRMED = "Order Confirmed"
        const val ITEM_REMOVED = "Item removed from cart"
        const val CART_CLEARED = "Cart cleared"
        const val ADDRESS_SAVED = "Address saved successfully"
        const val PAYMENT_SUCCESSFUL = "Payment processed successfully"
        const val ACCOUNT_CREATED = "Account created successfully"
        const val PASSWORD_UPDATED = "Password updated successfully"
    }

    // API Endpoints (for mocking)
    object ApiEndpoints {
        const val BASE_URL = "https://api.tchat.com/v1"
        const val PRODUCTS = "$BASE_URL/commerce/products"
        const val CATEGORIES = "$BASE_URL/commerce/categories"
        const val CART = "$BASE_URL/commerce/cart"
        const val COUPONS = "$BASE_URL/commerce/coupons"
        const val ORDERS = "$BASE_URL/commerce/orders"
        const val SEARCH = "$BASE_URL/commerce/search"
        const val AUTH_LOGIN = "$BASE_URL/auth/login"
        const val AUTH_REGISTER = "$BASE_URL/auth/register"
        const val USER_PROFILE = "$BASE_URL/users/profile"
        const val PAYMENT_PROCESS = "$BASE_URL/payment/process"
    }

    // Test Tags for UI Testing
    object TestTags {
        // Navigation
        const val STORE_TAB = "StoreTab"
        const val CART_TAB = "CartTab"
        const val PROFILE_TAB = "ProfileTab"
        const val SEARCH_TAB = "SearchTab"
        const val MORE_TAB = "MoreTab"

        // Store/Category
        const val CATEGORY_GRID = "CategoryGrid"
        const val CATEGORY_CARD = "CategoryCard"
        const val PRODUCT_GRID = "ProductGrid"
        const val PRODUCT_CARD = "ProductCard"
        const val CATEGORY_TITLE = "CategoryTitle"

        // Product Detail
        const val PRODUCT_TITLE = "ProductTitle"
        const val PRODUCT_PRICE = "ProductPrice"
        const val PRODUCT_DESCRIPTION = "ProductDescription"
        const val PRODUCT_IMAGES = "ProductImages"
        const val QUANTITY_SELECTOR = "QuantitySelector"
        const val QUANTITY_INCREMENT = "QuantityIncrement"
        const val QUANTITY_DECREMENT = "QuantityDecrement"
        const val ADD_TO_CART_BUTTON = "AddToCartButton"
        const val BUY_NOW_BUTTON = "BuyNowButton"

        // Cart
        const val CART_CONTAINER = "CartContainer"
        const val CART_ITEM_LIST = "CartItemList"
        const val CART_ITEM = "CartItem"
        const val CART_ITEM_QUANTITY = "CartItemQuantity"
        const val CART_ITEM_QUANTITY_INCREMENT = "CartItemQuantityIncrement"
        const val CART_ITEM_QUANTITY_DECREMENT = "CartItemQuantityDecrement"
        const val REMOVE_ITEM_BUTTON = "RemoveItemButton"
        const val CLEAR_CART_BUTTON = "ClearCartButton"
        const val CART_SUBTOTAL = "CartSubtotal"
        const val CART_DISCOUNT = "CartDiscount"
        const val CART_TOTAL = "CartTotal"
        const val CHECKOUT_BUTTON = "CheckoutButton"
        const val EMPTY_CART_MESSAGE = "EmptyCartMessage"

        // Coupon
        const val COUPON_TEXT_FIELD = "CouponTextField"
        const val APPLY_COUPON_BUTTON = "ApplyCouponButton"
        const val REMOVE_COUPON_BUTTON = "RemoveCouponButton"
        const val COUPON_APPLIED_INDICATOR = "CouponAppliedIndicator"
        const val COUPON_ERROR_MESSAGE = "CouponErrorMessage"
        const val DISCOUNT_AMOUNT = "DiscountAmount"

        // Search
        const val PRODUCT_SEARCH_BAR = "ProductSearchBar"
        const val SEARCH_BUTTON = "SearchButton"
        const val SEARCH_RESULTS = "SearchResults"
        const val SEARCH_FILTER_BUTTON = "SearchFilterButton"

        // Checkout
        const val EMAIL_TEXT_FIELD = "EmailTextField"
        const val FIRST_NAME_TEXT_FIELD = "FirstNameTextField"
        const val LAST_NAME_TEXT_FIELD = "LastNameTextField"
        const val ADDRESS_TEXT_FIELD = "AddressTextField"
        const val CITY_TEXT_FIELD = "CityTextField"
        const val STATE_TEXT_FIELD = "StateTextField"
        const val ZIP_CODE_TEXT_FIELD = "ZipCodeTextField"
        const val PHONE_TEXT_FIELD = "PhoneTextField"
        const val CONTINUE_TO_PAYMENT_BUTTON = "ContinueToPaymentButton"

        // Payment
        const val CARD_NUMBER_TEXT_FIELD = "CardNumberTextField"
        const val EXPIRY_MONTH_TEXT_FIELD = "ExpiryMonthTextField"
        const val EXPIRY_YEAR_TEXT_FIELD = "ExpiryYearTextField"
        const val CVV_TEXT_FIELD = "CVVTextField"
        const val CARD_NAME_TEXT_FIELD = "CardNameTextField"
        const val PLACE_ORDER_BUTTON = "PlaceOrderButton"

        // Confirmation
        const val ORDER_CONFIRMATION = "OrderConfirmation"
        const val ORDER_NUMBER = "OrderNumber"
        const val ADD_TO_CART_CONFIRMATION = "AddToCartConfirmation"
        const val CONFIRMATION_OK_BUTTON = "ConfirmationOkButton"

        // Status Indicators
        const val LOADING_INDICATOR = "LoadingIndicator"
        const val OFFLINE_INDICATOR = "OfflineIndicator"
        const val ERROR_MESSAGE = "ErrorMessage"

        // Dialogs and Alerts
        const val CONFIRM_REMOVAL_BUTTON = "ConfirmRemovalButton"
        const val CONFIRM_CLEAR_CART_BUTTON = "ConfirmClearCartButton"
        const val OFFLINE_DIALOG = "OfflineDialog"

        // Navigation
        const val BACK_BUTTON = "BackButton"
        const val HOME_BUTTON = "HomeButton"
    }
}

// Data Classes
data class User(
    val email: String,
    val password: String,
    val firstName: String,
    val lastName: String,
    val phone: String? = null
)

data class Product(
    val id: String,
    val name: String,
    val category: String,
    val price: Double,
    val sku: String,
    val stock: Int = 100
)

data class Coupon(
    val code: String,
    val type: CouponType,
    val value: Int,
    val minAmount: Int? = null,
    val maxDiscount: Int? = null,
    val expiresAt: String? = null
)

data class Address(
    val street: String,
    val city: String,
    val state: String,
    val zipCode: String,
    val country: String
)

data class PaymentMethod(
    val type: PaymentType,
    val cardNumber: String,
    val expiryMonth: String,
    val expiryYear: String,
    val cvv: String,
    val name: String
)

data class CartItem(
    val product: Product,
    val quantity: Int
)

data class CartScenario(
    val products: List<CartItem>,
    val coupon: Coupon? = null,
    val expectedSubtotal: Double? = null,
    val expectedDiscount: Double? = null,
    val expectedTotal: Double,
    val freeShipping: Boolean = false
)

enum class CouponType {
    PERCENTAGE,
    FIXED
}

enum class PaymentType {
    CREDIT_CARD,
    PAYPAL,
    APPLE_PAY,
    GOOGLE_PAY
}