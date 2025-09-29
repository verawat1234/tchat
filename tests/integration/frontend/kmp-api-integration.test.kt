/**
 * Comprehensive KMP (Kotlin Multiplatform) Integration Tests
 * Tests KMP API integration with commerce backend endpoints
 */

package com.tchat.tests.integration

import com.tchat.shared.api.CommerceApi
import com.tchat.shared.api.CartApi
import com.tchat.shared.api.ProductApi
import com.tchat.shared.api.CategoryApi
import com.tchat.shared.api.BusinessApi
import com.tchat.shared.models.*
import com.tchat.shared.network.NetworkClient
import com.tchat.shared.cache.CacheManager
import com.tchat.shared.database.DatabaseManager

import kotlinx.coroutines.test.*
import kotlinx.coroutines.flow.*
import kotlinx.coroutines.delay
import kotlin.test.*
import io.ktor.client.plugins.*
import io.ktor.http.*
import io.mockk.*
import kotlinx.datetime.Clock
import kotlinx.datetime.Instant

@OptIn(ExperimentalCoroutinesApi::class)
class KMPApiIntegrationTest {

    private lateinit var testDispatcher: TestDispatcher
    private lateinit var testScope: TestScope

    private lateinit var mockNetworkClient: NetworkClient
    private lateinit var mockCacheManager: CacheManager
    private lateinit var mockDatabaseManager: DatabaseManager

    private lateinit var commerceApi: CommerceApi
    private lateinit var cartApi: CartApi
    private lateinit var productApi: ProductApi
    private lateinit var categoryApi: CategoryApi
    private lateinit var businessApi: BusinessApi

    // Test data
    private val mockProducts = listOf(
        Product(
            id = "1",
            name = "iPhone 15 Pro",
            price = 999.99,
            currency = "USD",
            category = "smartphones",
            shopId = "shop-1",
            status = ProductStatus.Active,
            inventory = ProductInventory(
                quantity = 50,
                stockStatus = StockStatus.InStock
            ),
            images = listOf(
                ProductImage(
                    url = "https://example.com/iphone15-main.jpg",
                    altText = "iPhone 15 Pro",
                    isMain = true
                )
            ),
            variants = listOf(
                ProductVariant(
                    id = "variant-1",
                    name = "128GB Space Black",
                    options = mapOf("storage" to "128GB", "color" to "Space Black")
                )
            ),
            createdAt = Clock.System.now(),
            updatedAt = Clock.System.now()
        ),
        Product(
            id = "2",
            name = "MacBook Air M2",
            price = 1199.99,
            currency = "USD",
            category = "laptops",
            shopId = "shop-1",
            status = ProductStatus.Active,
            inventory = ProductInventory(
                quantity = 25,
                stockStatus = StockStatus.InStock
            ),
            images = listOf(
                ProductImage(
                    url = "https://example.com/macbook-air.jpg",
                    altText = "MacBook Air M2",
                    isMain = true
                )
            ),
            variants = emptyList(),
            createdAt = Clock.System.now(),
            updatedAt = Clock.System.now()
        )
    )

    private val mockCart = Cart(
        id = "cart-1",
        userId = "user-1",
        items = listOf(
            CartItem(
                id = "item-1",
                productId = "1",
                quantity = 1,
                unitPrice = 999.99,
                totalPrice = 999.99,
                name = "iPhone 15 Pro",
                createdAt = Clock.System.now(),
                updatedAt = Clock.System.now()
            )
        ),
        total = 999.99,
        currency = "USD",
        status = CartStatus.Active,
        createdAt = Clock.System.now(),
        updatedAt = Clock.System.now()
    )

    private val mockCategories = listOf(
        Category(
            id = "cat-1",
            name = "Electronics",
            slug = "electronics",
            level = 0,
            path = "/electronics",
            productCount = 2,
            children = listOf(
                Category(
                    id = "cat-2",
                    name = "Smartphones",
                    slug = "smartphones",
                    parentId = "cat-1",
                    level = 1,
                    path = "/electronics/smartphones",
                    productCount = 1,
                    createdAt = Clock.System.now(),
                    updatedAt = Clock.System.now()
                ),
                Category(
                    id = "cat-3",
                    name = "Laptops",
                    slug = "laptops",
                    parentId = "cat-1",
                    level = 1,
                    path = "/electronics/laptops",
                    productCount = 1,
                    createdAt = Clock.System.now(),
                    updatedAt = Clock.System.now()
                )
            ),
            createdAt = Clock.System.now(),
            updatedAt = Clock.System.now()
        )
    )

    private val mockBusiness = Business(
        id = "business-1",
        ownerId = "user-1",
        name = "TechCorp Solutions",
        businessType = BusinessType.Corporation,
        industry = "technology",
        email = "contact@techcorp.com",
        status = BusinessStatus.Active,
        verification = BusinessVerification(
            status = VerificationStatus.Verified,
            level = VerificationLevel.Full
        ),
        address = BusinessAddress(
            street = "123 Tech Street",
            city = "San Francisco",
            state = "CA",
            postalCode = "94105",
            country = "US"
        ),
        createdAt = Clock.System.now(),
        updatedAt = Clock.System.now()
    )

    @BeforeTest
    fun setup() {
        testDispatcher = StandardTestDispatcher()
        testScope = TestScope(testDispatcher)

        mockNetworkClient = mockk()
        mockCacheManager = mockk()
        mockDatabaseManager = mockk()

        // Setup default mock responses
        setupMockResponses()

        // Initialize APIs
        commerceApi = CommerceApi(mockNetworkClient, mockCacheManager)
        cartApi = CartApi(mockNetworkClient, mockCacheManager, mockDatabaseManager)
        productApi = ProductApi(mockNetworkClient, mockCacheManager, mockDatabaseManager)
        categoryApi = CategoryApi(mockNetworkClient, mockCacheManager, mockDatabaseManager)
        businessApi = BusinessApi(mockNetworkClient, mockCacheManager)
    }

    @AfterTest
    fun tearDown() {
        clearAllMocks()
    }

    private fun setupMockResponses() {
        // Cache manager mocks
        every { mockCacheManager.get<Any>(any()) } returns null
        every { mockCacheManager.put(any(), any<Any>(), any()) } just Runs
        every { mockCacheManager.invalidate(any()) } just Runs
        every { mockCacheManager.clear() } just Runs

        // Database manager mocks
        every { mockDatabaseManager.getProducts() } returns flowOf(mockProducts)
        every { mockDatabaseManager.getCart() } returns flowOf(mockCart)
        every { mockDatabaseManager.getCategories() } returns flowOf(mockCategories)
        every { mockDatabaseManager.insertProducts(any()) } just Runs
        every { mockDatabaseManager.insertCart(any()) } just Runs
        every { mockDatabaseManager.insertCategories(any()) } just Runs
    }

    @Test
    fun `should fetch products list successfully`() = testScope.runTest {
        // Given
        coEvery {
            mockNetworkClient.get<ProductListResponse>(any())
        } returns ApiResponse.Success(
            ProductListResponse(
                success = true,
                products = mockProducts,
                total = mockProducts.size,
                page = 1,
                limit = 20
            )
        )

        // When
        val result = productApi.getProducts(page = 1, limit = 20)

        // Then
        assertTrue(result.isSuccess)
        val products = result.getOrNull()
        assertNotNull(products)
        assertEquals(mockProducts.size, products.size)
        assertEquals("iPhone 15 Pro", products[0].name)
        assertEquals("MacBook Air M2", products[1].name)

        // Verify network call
        coVerify {
            mockNetworkClient.get<ProductListResponse>("/api/v1/commerce/products?page=1&limit=20")
        }
    }

    @Test
    fun `should fetch product by ID successfully`() = testScope.runTest {
        // Given
        val productId = "1"
        coEvery {
            mockNetworkClient.get<ProductResponse>("/api/v1/commerce/products/$productId")
        } returns ApiResponse.Success(
            ProductResponse(
                success = true,
                product = mockProducts[0]
            )
        )

        // When
        val result = productApi.getProduct(productId)

        // Then
        assertTrue(result.isSuccess)
        val product = result.getOrNull()
        assertNotNull(product)
        assertEquals(productId, product.id)
        assertEquals("iPhone 15 Pro", product.name)
        assertEquals(999.99, product.price)
    }

    @Test
    fun `should handle product not found error`() = testScope.runTest {
        // Given
        val productId = "999"
        coEvery {
            mockNetworkClient.get<ProductResponse>("/api/v1/commerce/products/$productId")
        } returns ApiResponse.Error(
            error = NetworkError.NotFound("Product not found"),
            statusCode = 404
        )

        // When
        val result = productApi.getProduct(productId)

        // Then
        assertTrue(result.isFailure)
        val exception = result.exceptionOrNull()
        assertNotNull(exception)
        assertTrue(exception is NetworkError.NotFound)
    }

    @Test
    fun `should search products with query`() = testScope.runTest {
        // Given
        val query = "iPhone"
        val filteredProducts = mockProducts.filter { it.name.contains(query, ignoreCase = true) }

        coEvery {
            mockNetworkClient.get<ProductListResponse>(any())
        } returns ApiResponse.Success(
            ProductListResponse(
                success = true,
                products = filteredProducts,
                total = filteredProducts.size,
                page = 1,
                limit = 20
            )
        )

        // When
        val result = productApi.searchProducts(query = query)

        // Then
        assertTrue(result.isSuccess)
        val products = result.getOrNull()
        assertNotNull(products)
        assertEquals(1, products.size)
        assertTrue(products[0].name.contains("iPhone"))

        // Verify correct endpoint was called
        coVerify {
            mockNetworkClient.get<ProductListResponse>(match { url ->
                url.contains("query=iPhone")
            })
        }
    }

    @Test
    fun `should filter products by category`() = testScope.runTest {
        // Given
        val category = "smartphones"
        val filteredProducts = mockProducts.filter { it.category == category }

        coEvery {
            mockNetworkClient.get<ProductListResponse>(any())
        } returns ApiResponse.Success(
            ProductListResponse(
                success = true,
                products = filteredProducts,
                total = filteredProducts.size,
                page = 1,
                limit = 20
            )
        )

        // When
        val result = productApi.getProductsByCategory(category)

        // Then
        assertTrue(result.isSuccess)
        val products = result.getOrNull()
        assertNotNull(products)
        assertEquals(1, products.size)
        assertEquals(category, products[0].category)
    }

    @Test
    fun `should create product successfully`() = testScope.runTest {
        // Given
        val shopId = "shop-1"
        val newProduct = CreateProductRequest(
            name = "Test Product",
            price = 199.99,
            currency = "USD",
            category = "test",
            inventory = ProductInventory(
                quantity = 10,
                stockStatus = StockStatus.InStock
            )
        )

        val createdProduct = mockProducts[0].copy(
            id = "product-${Clock.System.now().toEpochMilliseconds()}",
            name = newProduct.name,
            price = newProduct.price,
            category = newProduct.category
        )

        coEvery {
            mockNetworkClient.post<ProductResponse>(any(), any())
        } returns ApiResponse.Success(
            ProductResponse(
                success = true,
                product = createdProduct
            )
        )

        // When
        val result = productApi.createProduct(shopId, newProduct)

        // Then
        assertTrue(result.isSuccess)
        val product = result.getOrNull()
        assertNotNull(product)
        assertEquals(newProduct.name, product.name)
        assertEquals(newProduct.price, product.price)
        assertEquals(shopId, product.shopId)

        // Verify network call
        coVerify {
            mockNetworkClient.post<ProductResponse>(
                "/api/v1/commerce/shops/$shopId/products",
                newProduct
            )
        }
    }

    @Test
    fun `should update product successfully`() = testScope.runTest {
        // Given
        val productId = "1"
        val updates = UpdateProductRequest(
            name = "Updated iPhone 15 Pro",
            price = 1099.99
        )

        val updatedProduct = mockProducts[0].copy(
            name = updates.name!!,
            price = updates.price!!,
            updatedAt = Clock.System.now()
        )

        coEvery {
            mockNetworkClient.put<ProductResponse>(any(), any())
        } returns ApiResponse.Success(
            ProductResponse(
                success = true,
                product = updatedProduct
            )
        )

        // When
        val result = productApi.updateProduct(productId, updates)

        // Then
        assertTrue(result.isSuccess)
        val product = result.getOrNull()
        assertNotNull(product)
        assertEquals(updates.name, product.name)
        assertEquals(updates.price, product.price)

        // Verify cache invalidation
        verify { mockCacheManager.invalidate("products_list") }
        verify { mockCacheManager.invalidate("product_$productId") }
    }

    @Test
    fun `should delete product successfully`() = testScope.runTest {
        // Given
        val productId = "1"

        coEvery {
            mockNetworkClient.delete<ApiResponse<Unit>>(any())
        } returns ApiResponse.Success(Unit)

        // When
        val result = productApi.deleteProduct(productId)

        // Then
        assertTrue(result.isSuccess)

        // Verify network call
        coVerify {
            mockNetworkClient.delete<ApiResponse<Unit>>("/api/v1/commerce/products/$productId")
        }

        // Verify cache invalidation
        verify { mockCacheManager.invalidate("products_list") }
        verify { mockCacheManager.invalidate("product_$productId") }
    }

    @Test
    fun `should fetch cart successfully`() = testScope.runTest {
        // Given
        coEvery {
            mockNetworkClient.get<CartResponse>(any())
        } returns ApiResponse.Success(
            CartResponse(
                success = true,
                cart = mockCart
            )
        )

        // When
        val result = cartApi.getCart()

        // Then
        assertTrue(result.isSuccess)
        val cart = result.getOrNull()
        assertNotNull(cart)
        assertEquals("cart-1", cart.id)
        assertEquals(1, cart.items.size)
        assertEquals(999.99, cart.total)

        // Verify database sync
        verify { mockDatabaseManager.insertCart(cart) }
    }

    @Test
    fun `should add item to cart successfully`() = testScope.runTest {
        // Given
        val addToCartRequest = AddToCartRequest(
            productId = "2",
            quantity = 1,
            unitPrice = 1199.99
        )

        val updatedCart = mockCart.copy(
            items = mockCart.items + CartItem(
                id = "item-2",
                productId = addToCartRequest.productId,
                quantity = addToCartRequest.quantity,
                unitPrice = addToCartRequest.unitPrice,
                totalPrice = addToCartRequest.quantity * addToCartRequest.unitPrice,
                name = "MacBook Air M2",
                createdAt = Clock.System.now(),
                updatedAt = Clock.System.now()
            ),
            total = mockCart.total + (addToCartRequest.quantity * addToCartRequest.unitPrice),
            updatedAt = Clock.System.now()
        )

        coEvery {
            mockNetworkClient.post<CartResponse>(any(), any())
        } returns ApiResponse.Success(
            CartResponse(
                success = true,
                cart = updatedCart
            )
        )

        // When
        val result = cartApi.addToCart(addToCartRequest)

        // Then
        assertTrue(result.isSuccess)
        val cart = result.getOrNull()
        assertNotNull(cart)
        assertEquals(2, cart.items.size)
        assertTrue(cart.total > mockCart.total)

        // Verify optimistic update was applied
        verify { mockDatabaseManager.insertCart(any()) }
    }

    @Test
    fun `should update cart item successfully`() = testScope.runTest {
        // Given
        val itemId = "item-1"
        val updateRequest = UpdateCartItemRequest(
            quantity = 2
        )

        val updatedCart = mockCart.copy(
            items = mockCart.items.map { item ->
                if (item.id == itemId) {
                    item.copy(
                        quantity = updateRequest.quantity,
                        totalPrice = updateRequest.quantity * item.unitPrice,
                        updatedAt = Clock.System.now()
                    )
                } else item
            },
            total = updateRequest.quantity * mockCart.items[0].unitPrice,
            updatedAt = Clock.System.now()
        )

        coEvery {
            mockNetworkClient.put<CartResponse>(any(), any())
        } returns ApiResponse.Success(
            CartResponse(
                success = true,
                cart = updatedCart
            )
        )

        // When
        val result = cartApi.updateCartItem(itemId, updateRequest)

        // Then
        assertTrue(result.isSuccess)
        val cart = result.getOrNull()
        assertNotNull(cart)
        assertEquals(2, cart.items[0].quantity)
        assertEquals(1999.98, cart.total, 0.01)
    }

    @Test
    fun `should remove item from cart successfully`() = testScope.runTest {
        // Given
        val itemId = "item-1"

        val updatedCart = mockCart.copy(
            items = emptyList(),
            total = 0.0,
            updatedAt = Clock.System.now()
        )

        coEvery {
            mockNetworkClient.delete<CartResponse>(any())
        } returns ApiResponse.Success(
            CartResponse(
                success = true,
                cart = updatedCart
            )
        )

        // When
        val result = cartApi.removeCartItem(itemId)

        // Then
        assertTrue(result.isSuccess)
        val cart = result.getOrNull()
        assertNotNull(cart)
        assertEquals(0, cart.items.size)
        assertEquals(0.0, cart.total)
    }

    @Test
    fun `should clear cart successfully`() = testScope.runTest {
        // Given
        val clearedCart = mockCart.copy(
            items = emptyList(),
            total = 0.0,
            updatedAt = Clock.System.now()
        )

        coEvery {
            mockNetworkClient.delete<CartResponse>(any())
        } returns ApiResponse.Success(
            CartResponse(
                success = true,
                cart = clearedCart
            )
        )

        // When
        val result = cartApi.clearCart()

        // Then
        assertTrue(result.isSuccess)
        val cart = result.getOrNull()
        assertNotNull(cart)
        assertEquals(0, cart.items.size)
        assertEquals(0.0, cart.total)
    }

    @Test
    fun `should fetch categories successfully`() = testScope.runTest {
        // Given
        coEvery {
            mockNetworkClient.get<CategoryListResponse>(any())
        } returns ApiResponse.Success(
            CategoryListResponse(
                success = true,
                categories = mockCategories,
                total = mockCategories.size
            )
        )

        // When
        val result = categoryApi.getCategories()

        // Then
        assertTrue(result.isSuccess)
        val categories = result.getOrNull()
        assertNotNull(categories)
        assertEquals(1, categories.size)
        assertEquals("Electronics", categories[0].name)
        assertEquals(2, categories[0].children?.size)
    }

    @Test
    fun `should fetch category hierarchy successfully`() = testScope.runTest {
        // Given
        coEvery {
            mockNetworkClient.get<CategoryListResponse>(any())
        } returns ApiResponse.Success(
            CategoryListResponse(
                success = true,
                categories = mockCategories
            )
        )

        // When
        val result = categoryApi.getCategoryHierarchy()

        // Then
        assertTrue(result.isSuccess)
        val categories = result.getOrNull()
        assertNotNull(categories)
        assertEquals(1, categories.size)

        val rootCategory = categories[0]
        assertEquals(0, rootCategory.level)
        assertEquals(2, rootCategory.children?.size)

        val smartphone = rootCategory.children?.find { it.name == "Smartphones" }
        assertNotNull(smartphone)
        assertEquals(1, smartphone.level)
        assertEquals("cat-1", smartphone.parentId)
    }

    @Test
    fun `should fetch category by ID successfully`() = testScope.runTest {
        // Given
        val categoryId = "cat-1"
        coEvery {
            mockNetworkClient.get<CategoryResponse>("/api/v1/commerce/categories/$categoryId")
        } returns ApiResponse.Success(
            CategoryResponse(
                success = true,
                category = mockCategories[0]
            )
        )

        // When
        val result = categoryApi.getCategory(categoryId)

        // Then
        assertTrue(result.isSuccess)
        val category = result.getOrNull()
        assertNotNull(category)
        assertEquals(categoryId, category.id)
        assertEquals("Electronics", category.name)
    }

    @Test
    fun `should create category successfully`() = testScope.runTest {
        // Given
        val createRequest = CreateCategoryRequest(
            name = "New Category",
            description = "Test category"
        )

        val createdCategory = Category(
            id = "category-${Clock.System.now().toEpochMilliseconds()}",
            name = createRequest.name,
            slug = "new-category",
            description = createRequest.description,
            level = 0,
            path = "/new-category",
            productCount = 0,
            createdAt = Clock.System.now(),
            updatedAt = Clock.System.now()
        )

        coEvery {
            mockNetworkClient.post<CategoryResponse>(any(), any())
        } returns ApiResponse.Success(
            CategoryResponse(
                success = true,
                category = createdCategory
            )
        )

        // When
        val result = categoryApi.createCategory(createRequest)

        // Then
        assertTrue(result.isSuccess)
        val category = result.getOrNull()
        assertNotNull(category)
        assertEquals(createRequest.name, category.name)
        assertEquals(0, category.level)
    }

    @Test
    fun `should fetch businesses successfully`() = testScope.runTest {
        // Given
        coEvery {
            mockNetworkClient.get<BusinessListResponse>(any())
        } returns ApiResponse.Success(
            BusinessListResponse(
                success = true,
                businesses = listOf(mockBusiness),
                total = 1
            )
        )

        // When
        val result = businessApi.getBusinesses()

        // Then
        assertTrue(result.isSuccess)
        val businesses = result.getOrNull()
        assertNotNull(businesses)
        assertEquals(1, businesses.size)
        assertEquals("TechCorp Solutions", businesses[0].name)
    }

    @Test
    fun `should fetch business by ID successfully`() = testScope.runTest {
        // Given
        val businessId = "business-1"
        coEvery {
            mockNetworkClient.get<BusinessResponse>("/api/v1/commerce/businesses/$businessId")
        } returns ApiResponse.Success(
            BusinessResponse(
                success = true,
                business = mockBusiness
            )
        )

        // When
        val result = businessApi.getBusiness(businessId)

        // Then
        assertTrue(result.isSuccess)
        val business = result.getOrNull()
        assertNotNull(business)
        assertEquals(businessId, business.id)
        assertEquals("TechCorp Solutions", business.name)
        assertEquals(VerificationStatus.Verified, business.verification.status)
    }

    @Test
    fun `should create business successfully`() = testScope.runTest {
        // Given
        val createRequest = CreateBusinessRequest(
            name = "New Business",
            businessType = BusinessType.LLC,
            industry = "retail",
            email = "contact@newbusiness.com",
            address = BusinessAddress(
                street = "123 New St",
                city = "New City",
                state = "NY",
                postalCode = "12345",
                country = "US"
            )
        )

        val createdBusiness = mockBusiness.copy(
            id = "business-${Clock.System.now().toEpochMilliseconds()}",
            name = createRequest.name,
            businessType = createRequest.businessType,
            industry = createRequest.industry,
            email = createRequest.email,
            status = BusinessStatus.Pending,
            verification = BusinessVerification(
                status = VerificationStatus.Unverified,
                level = VerificationLevel.None
            )
        )

        coEvery {
            mockNetworkClient.post<BusinessResponse>(any(), any())
        } returns ApiResponse.Success(
            BusinessResponse(
                success = true,
                business = createdBusiness
            )
        )

        // When
        val result = businessApi.createBusiness(createRequest)

        // Then
        assertTrue(result.isSuccess)
        val business = result.getOrNull()
        assertNotNull(business)
        assertEquals(createRequest.name, business.name)
        assertEquals(BusinessStatus.Pending, business.status)
    }

    @Test
    fun `should handle network errors gracefully`() = testScope.runTest {
        // Given
        coEvery {
            mockNetworkClient.get<ProductListResponse>(any())
        } returns ApiResponse.Error(
            error = NetworkError.ServerError("Internal server error"),
            statusCode = 500
        )

        // When
        val result = productApi.getProducts()

        // Then
        assertTrue(result.isFailure)
        val exception = result.exceptionOrNull()
        assertNotNull(exception)
        assertTrue(exception is NetworkError.ServerError)
    }

    @Test
    fun `should handle network timeout gracefully`() = testScope.runTest {
        // Given
        coEvery {
            mockNetworkClient.get<ProductListResponse>(any())
        } coAnswers {
            delay(5000)
            ApiResponse.Error(
                error = NetworkError.Timeout("Request timeout"),
                statusCode = 408
            )
        }

        // When
        val result = productApi.getProducts()

        // Then
        assertTrue(result.isFailure)
        val exception = result.exceptionOrNull()
        assertNotNull(exception)
        assertTrue(exception is NetworkError.Timeout)
    }

    @Test
    fun `should retry failed requests according to configuration`() = testScope.runTest {
        // Given
        var callCount = 0
        coEvery {
            mockNetworkClient.get<ProductListResponse>(any())
        } coAnswers {
            callCount++
            if (callCount < 3) {
                ApiResponse.Error(
                    error = NetworkError.ServerError("Server error"),
                    statusCode = 500
                )
            } else {
                ApiResponse.Success(
                    ProductListResponse(
                        success = true,
                        products = mockProducts,
                        total = mockProducts.size
                    )
                )
            }
        }

        // When
        val result = productApi.getProducts()

        // Then
        assertTrue(result.isSuccess)
        assertEquals(3, callCount) // Should have retried twice
    }

    @Test
    fun `should use cached data when available`() = testScope.runTest {
        // Given
        val cachedProducts = mockProducts
        every {
            mockCacheManager.get<List<Product>>("products_list")
        } returns cachedProducts

        // When
        val result = productApi.getProducts()

        // Then
        assertTrue(result.isSuccess)
        val products = result.getOrNull()
        assertEquals(cachedProducts, products)

        // Verify network was not called
        coVerify(exactly = 0) {
            mockNetworkClient.get<ProductListResponse>(any())
        }
    }

    @Test
    fun `should sync data with local database`() = testScope.runTest {
        // Given
        coEvery {
            mockNetworkClient.get<ProductListResponse>(any())
        } returns ApiResponse.Success(
            ProductListResponse(
                success = true,
                products = mockProducts,
                total = mockProducts.size
            )
        )

        // When
        val result = productApi.getProducts()

        // Then
        assertTrue(result.isSuccess)

        // Verify data was synced to database
        verify { mockDatabaseManager.insertProducts(mockProducts) }
    }

    @Test
    fun `should handle offline mode gracefully`() = testScope.runTest {
        // Given - network is offline
        coEvery {
            mockNetworkClient.get<ProductListResponse>(any())
        } returns ApiResponse.Error(
            error = NetworkError.NoConnection("No internet connection"),
            statusCode = 0
        )

        // But we have cached data in database
        every { mockDatabaseManager.getProducts() } returns flowOf(mockProducts)

        // When
        val result = productApi.getProducts()

        // Then
        assertTrue(result.isSuccess)
        val products = result.getOrNull()
        assertEquals(mockProducts, products)
    }

    @Test
    fun `should invalidate cache when mutations occur`() = testScope.runTest {
        // Given
        val newProduct = CreateProductRequest(
            name = "Cache Test Product",
            price = 99.99,
            currency = "USD",
            category = "test"
        )

        coEvery {
            mockNetworkClient.post<ProductResponse>(any(), any())
        } returns ApiResponse.Success(
            ProductResponse(
                success = true,
                product = mockProducts[0].copy(name = newProduct.name)
            )
        )

        // When
        val result = productApi.createProduct("shop-1", newProduct)

        // Then
        assertTrue(result.isSuccess)

        // Verify cache was invalidated
        verify { mockCacheManager.invalidate("products_list") }
        verify { mockCacheManager.invalidate("shop_shop-1_products") }
    }

    @Test
    fun `should handle optimistic updates for cart operations`() = testScope.runTest {
        // Given
        val addToCartRequest = AddToCartRequest(
            productId = "2",
            quantity = 1,
            unitPrice = 199.99
        )

        // Network call succeeds
        coEvery {
            mockNetworkClient.post<CartResponse>(any(), any())
        } returns ApiResponse.Success(
            CartResponse(
                success = true,
                cart = mockCart.copy(
                    items = mockCart.items + CartItem(
                        id = "item-2",
                        productId = addToCartRequest.productId,
                        quantity = addToCartRequest.quantity,
                        unitPrice = addToCartRequest.unitPrice,
                        totalPrice = addToCartRequest.quantity * addToCartRequest.unitPrice,
                        name = "Test Product",
                        createdAt = Clock.System.now(),
                        updatedAt = Clock.System.now()
                    )
                )
            )
        )

        // When
        val result = cartApi.addToCart(addToCartRequest)

        // Then
        assertTrue(result.isSuccess)

        // Verify optimistic update was applied to database immediately
        verify { mockDatabaseManager.insertCart(any()) }

        // And cache was invalidated
        verify { mockCacheManager.invalidate("cart") }
    }

    @Test
    fun `should handle real-time data synchronization`() = testScope.runTest {
        // Given
        val cartFlow = MutableSharedFlow<Cart>()
        every { mockDatabaseManager.getCart() } returns cartFlow.asSharedFlow()

        // When
        val flow = cartApi.observeCart()
        val emissions = mutableListOf<Cart>()

        val job = launch {
            flow.collect { cart ->
                emissions.add(cart)
            }
        }

        // Emit updates
        cartFlow.emit(mockCart)
        cartFlow.emit(mockCart.copy(total = 1500.0))
        cartFlow.emit(mockCart.copy(total = 2000.0))

        // Wait for emissions
        testScheduler.advanceUntilIdle()

        // Then
        assertEquals(3, emissions.size)
        assertEquals(999.99, emissions[0].total, 0.01)
        assertEquals(1500.0, emissions[1].total, 0.01)
        assertEquals(2000.0, emissions[2].total, 0.01)

        job.cancel()
    }
}

// Extension functions for test utilities
private fun TestScope.launch(block: suspend CoroutineScope.() -> Unit) =
    kotlinx.coroutines.launch(testDispatcher, block = block)

private val TestScope.testScheduler: TestCoroutineScheduler
    get() = testDispatcher.scheduler