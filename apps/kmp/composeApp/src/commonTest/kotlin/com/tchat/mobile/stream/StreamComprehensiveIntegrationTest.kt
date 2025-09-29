package com.tchat.mobile.stream

import kotlin.test.*
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.delay
import kotlinx.serialization.json.Json

import com.tchat.mobile.stream.models.*

/**
 * Comprehensive integration tests for Stream Store Tabs functionality
 * Tests complete user flows across the KMP mobile application
 */
class StreamComprehensiveIntegrationTest {

    private val json = Json { ignoreUnknownKeys = true }

    // Test data matching backend implementation
    private val mockCategories = listOf(
        StreamCategory(
            id = "books",
            name = "Books",
            displayOrder = 1,
            iconName = "book-open",
            isActive = true,
            subtabs = null,
            featuredContentEnabled = true,
            createdAt = "2024-01-01T00:00:00Z",
            updatedAt = "2024-01-01T00:00:00Z"
        ),
        StreamCategory(
            id = "podcasts",
            name = "Podcasts",
            displayOrder = 2,
            iconName = "headphones",
            isActive = true,
            subtabs = null,
            featuredContentEnabled = true,
            createdAt = "2024-01-01T00:00:00Z",
            updatedAt = "2024-01-01T00:00:00Z"
        ),
        StreamCategory(
            id = "movies",
            name = "Movies",
            displayOrder = 4,
            iconName = "video",
            isActive = true,
            subtabs = listOf(
                StreamSubtab(
                    id = "movies_short",
                    categoryId = "movies",
                    name = "Short Films",
                    displayOrder = 1,
                    filterCriteria = mapOf("content_type" to "SHORT_MOVIE", "max_duration" to "1800"),
                    isActive = true,
                    createdAt = "2024-01-01T00:00:00Z",
                    updatedAt = "2024-01-01T00:00:00Z"
                ),
                StreamSubtab(
                    id = "movies_feature",
                    categoryId = "movies",
                    name = "Feature Films",
                    displayOrder = 2,
                    filterCriteria = mapOf("content_type" to "LONG_MOVIE", "min_duration" to "1800"),
                    isActive = true,
                    createdAt = "2024-01-01T00:00:00Z",
                    updatedAt = "2024-01-01T00:00:00Z"
                )
            ),
            featuredContentEnabled = true,
            createdAt = "2024-01-01T00:00:00Z",
            updatedAt = "2024-01-01T00:00:00Z"
        )
    )

    private val mockStreamContent = listOf(
        StreamContentItem(
            id = "book-1",
            categoryId = "books",
            title = "Test Book 1",
            description = "A fascinating test book",
            thumbnailUrl = "https://example.com/book1.jpg",
            contentType = StreamContentType.BOOK,
            duration = null,
            price = 9.99,
            currency = "USD",
            availabilityStatus = StreamAvailabilityStatus.AVAILABLE,
            isFeatured = true,
            featuredOrder = 1,
            metadata = mapOf("author" to "Test Author", "genre" to "fiction"),
            createdAt = "2024-01-01T00:00:00Z",
            updatedAt = "2024-01-01T00:00:00Z"
        ),
        StreamContentItem(
            id = "movie-short-1",
            categoryId = "movies",
            title = "Test Short Film",
            description = "A compelling short film",
            thumbnailUrl = "https://example.com/short1.jpg",
            contentType = StreamContentType.SHORT_MOVIE,
            duration = 1200, // 20 minutes
            price = 2.99,
            currency = "USD",
            availabilityStatus = StreamAvailabilityStatus.AVAILABLE,
            isFeatured = true,
            featuredOrder = 1,
            metadata = mapOf("director" to "Test Director", "year" to "2024"),
            createdAt = "2024-01-01T00:00:00Z",
            updatedAt = "2024-01-01T00:00:00Z"
        ),
        StreamContentItem(
            id = "movie-feature-1",
            categoryId = "movies",
            title = "Test Feature Film",
            description = "An epic feature film",
            thumbnailUrl = "https://example.com/feature1.jpg",
            contentType = StreamContentType.LONG_MOVIE,
            duration = 7200, // 2 hours
            price = 12.99,
            currency = "USD",
            availabilityStatus = StreamAvailabilityStatus.AVAILABLE,
            isFeatured = false,
            featuredOrder = null,
            metadata = mapOf("director" to "Test Director", "year" to "2024", "rating" to "PG-13"),
            createdAt = "2024-01-01T00:00:00Z",
            updatedAt = "2024-01-01T00:00:00Z"
        )
    )

    // Mock API client for testing
    class MockStreamApiClient {
        suspend fun getCategories(): StreamCategoriesResponse {
            delay(50) // Simulate network delay
            return StreamCategoriesResponse(
                categories = mockCategories,
                total = mockCategories.size,
                success = true
            )
        }

        suspend fun getCategoryDetail(categoryId: String): Result<StreamCategory> {
            delay(50)
            val category = mockCategories.find { it.id == categoryId }
            return if (category != null) {
                Result.success(category)
            } else {
                Result.failure(Exception("Category not found"))
            }
        }

        suspend fun getContent(request: StreamContentRequest): StreamContentResponse {
            delay(100) // Simulate network delay

            var filteredContent = mockStreamContent.filter { it.categoryId == request.categoryId }

            // Apply subtab filtering
            when (request.subtabId) {
                "movies_short" -> filteredContent = filteredContent.filter {
                    it.contentType == StreamContentType.SHORT_MOVIE && (it.duration ?: 0) <= 1800
                }
                "movies_feature" -> filteredContent = filteredContent.filter {
                    it.contentType == StreamContentType.LONG_MOVIE && (it.duration ?: 0) > 1800
                }
            }

            // Apply pagination
            val startIndex = (request.page - 1) * request.limit
            val endIndex = minOf(startIndex + request.limit, filteredContent.size)
            val paginatedContent = if (startIndex < filteredContent.size) {
                filteredContent.subList(startIndex, endIndex)
            } else {
                emptyList()
            }

            return StreamContentResponse(
                content = paginatedContent,
                page = request.page,
                limit = request.limit,
                total = filteredContent.size,
                success = true
            )
        }

        suspend fun getFeaturedContent(categoryId: String, limit: Int): StreamFeaturedResponse {
            delay(75)
            val featuredContent = mockStreamContent
                .filter { it.categoryId == categoryId && it.isFeatured }
                .sortedBy { it.featuredOrder ?: 999 }
                .take(limit)

            return StreamFeaturedResponse(
                content = featuredContent,
                total = featuredContent.size,
                success = true
            )
        }

        suspend fun purchaseContent(request: PurchaseContentRequest): PurchaseResponse {
            delay(200) // Simulate payment processing
            return PurchaseResponse(
                success = true,
                orderId = "test-order-123",
                downloadUrls = listOf("https://example.com/download/content"),
                message = "Purchase successful!"
            )
        }
    }

    private val mockApiClient = MockStreamApiClient()

    /**
     * FLOW TEST 1: Complete Navigation Flow
     * Tests user navigation between different stream categories
     */
    @Test
    fun testCompleteNavigationFlow() = runTest {
        // Step 1: Load categories
        val categoriesResponse = mockApiClient.getCategories()

        assertTrue(categoriesResponse.success)
        assertEquals(3, categoriesResponse.total)
        assertEquals("books", categoriesResponse.categories[0].id)
        assertEquals("podcasts", categoriesResponse.categories[1].id)
        assertEquals("movies", categoriesResponse.categories[2].id)

        // Step 2: Navigate to specific category (Movies)
        val moviesCategoryResult = mockApiClient.getCategoryDetail("movies")
        assertTrue(moviesCategoryResult.isSuccess)

        val moviesCategory = moviesCategoryResult.getOrNull()!!
        assertEquals("Movies", moviesCategory.name)
        assertEquals(2, moviesCategory.subtabs?.size)

        // Step 3: Load content for category
        val contentRequest = StreamContentRequest(
            categoryId = "movies",
            subtabId = null,
            page = 1,
            limit = 10
        )

        val contentResponse = mockApiClient.getContent(contentRequest)
        assertTrue(contentResponse.success)
        assertEquals(2, contentResponse.total) // 2 movies

        val movieContent = contentResponse.content
        assertTrue(movieContent.any { it.contentType == StreamContentType.SHORT_MOVIE })
        assertTrue(movieContent.any { it.contentType == StreamContentType.LONG_MOVIE })
    }

    /**
     * FLOW TEST 2: Featured Content Discovery Flow
     * Tests featured content carousel functionality
     */
    @Test
    fun testFeaturedContentFlow() = runTest {
        // Get featured content for books category
        val booksFeaturedResponse = mockApiClient.getFeaturedContent("books", 5)

        assertTrue(booksFeaturedResponse.success)
        assertEquals(1, booksFeaturedResponse.total) // 1 featured book

        val featuredBook = booksFeaturedResponse.content[0]
        assertTrue(featuredBook.isFeatured)
        assertEquals("Test Book 1", featuredBook.title)
        assertEquals("books", featuredBook.categoryId)

        // Get featured content for movies category
        val moviesFeaturedResponse = mockApiClient.getFeaturedContent("movies", 5)

        assertTrue(moviesFeaturedResponse.success)
        assertEquals(1, moviesFeaturedResponse.total) // 1 featured movie

        val featuredMovie = moviesFeaturedResponse.content[0]
        assertTrue(featuredMovie.isFeatured)
        assertEquals("Test Short Film", featuredMovie.title)
    }

    /**
     * FLOW TEST 3: Movies Subtab Navigation Flow
     * Tests subtab filtering functionality
     */
    @Test
    fun testSubtabFilteringFlow() = runTest {
        // Get short movies only
        val shortMoviesRequest = StreamContentRequest(
            categoryId = "movies",
            subtabId = "movies_short",
            page = 1,
            limit = 10
        )

        val shortMoviesResponse = mockApiClient.getContent(shortMoviesRequest)
        assertTrue(shortMoviesResponse.success)

        // Should only return short movies
        shortMoviesResponse.content.forEach { movie ->
            assertEquals(StreamContentType.SHORT_MOVIE, movie.contentType)
            assertTrue((movie.duration ?: 0) <= 1800) // ≤ 30 minutes
        }

        // Get feature movies only
        val featureMoviesRequest = StreamContentRequest(
            categoryId = "movies",
            subtabId = "movies_feature",
            page = 1,
            limit = 10
        )

        val featureMoviesResponse = mockApiClient.getContent(featureMoviesRequest)
        assertTrue(featureMoviesResponse.success)

        // Should only return feature movies
        featureMoviesResponse.content.forEach { movie ->
            assertEquals(StreamContentType.LONG_MOVIE, movie.contentType)
            assertTrue((movie.duration ?: 0) > 1800) // > 30 minutes
        }
    }

    /**
     * FLOW TEST 4: Content Purchase Flow
     * Tests the complete purchase workflow
     */
    @Test
    fun testContentPurchaseFlow() = runTest {
        val purchaseRequest = PurchaseContentRequest(
            contentId = "book-1",
            mediaLicense = "personal",
            downloadFormat = "standard"
        )

        val purchaseResponse = mockApiClient.purchaseContent(purchaseRequest)

        assertTrue(purchaseResponse.success)
        assertEquals("test-order-123", purchaseResponse.orderId)
        assertNotNull(purchaseResponse.downloadUrls)
        assertTrue(purchaseResponse.downloadUrls!!.isNotEmpty())
        assertEquals("Purchase successful!", purchaseResponse.message)
    }

    /**
     * FLOW TEST 5: Pagination Flow
     * Tests content pagination functionality
     */
    @Test
    fun testPaginationFlow() = runTest {
        // Test pagination with small limit
        val page1Request = StreamContentRequest(
            categoryId = "movies",
            page = 1,
            limit = 1
        )

        val page1Response = mockApiClient.getContent(page1Request)
        assertTrue(page1Response.success)
        assertEquals(1, page1Response.content.size)
        assertEquals(2, page1Response.total) // Total 2 movies
        assertEquals(1, page1Response.page)

        // Get page 2
        val page2Request = StreamContentRequest(
            categoryId = "movies",
            page = 2,
            limit = 1
        )

        val page2Response = mockApiClient.getContent(page2Request)
        assertTrue(page2Response.success)
        assertEquals(1, page2Response.content.size)
        assertEquals(2, page2Response.total)
        assertEquals(2, page2Response.page)

        // Pages should have different content
        assertNotEquals(page1Response.content[0].id, page2Response.content[0].id)
    }

    /**
     * FLOW TEST 6: Data Model Validation Flow
     * Tests cross-platform data consistency
     */
    @Test
    fun testDataModelConsistency() = runTest {
        val categoriesResponse = mockApiClient.getCategories()

        // Validate category structure
        categoriesResponse.categories.forEach { category ->
            assertNotNull(category.id)
            assertNotNull(category.name)
            assertTrue(category.displayOrder > 0)
            assertNotNull(category.iconName)
            assertNotNull(category.createdAt)
            assertNotNull(category.updatedAt)
        }

        // Test content structure
        val contentResponse = mockApiClient.getContent(
            StreamContentRequest("books", page = 1, limit = 10)
        )

        contentResponse.content.forEach { content ->
            assertNotNull(content.id)
            assertNotNull(content.title)
            assertNotNull(content.description)
            assertNotNull(content.thumbnailUrl)
            assertTrue(content.price >= 0.0)
            assertNotNull(content.currency)
            assertNotNull(content.availabilityStatus)
            assertNotNull(content.createdAt)
            assertNotNull(content.updatedAt)
        }
    }

    /**
     * FLOW TEST 7: Content Type Validation Flow
     * Tests content type specific business logic
     */
    @Test
    fun testContentTypeValidation() = runTest {
        val contentResponse = mockApiClient.getContent(
            StreamContentRequest("books", page = 1, limit = 10)
        )

        val bookContent = contentResponse.content.find { it.contentType == StreamContentType.BOOK }
        assertNotNull(bookContent)

        // Books should not have duration
        assertNull(bookContent.duration)
        assertTrue(bookContent.isBook())
        assertFalse(bookContent.isVideo())
        assertFalse(bookContent.isAudio())

        val movieResponse = mockApiClient.getContent(
            StreamContentRequest("movies", page = 1, limit = 10)
        )

        val movieContent = movieResponse.content.find { it.contentType == StreamContentType.SHORT_MOVIE }
        assertNotNull(movieContent)

        // Movies should have duration
        assertNotNull(movieContent.duration)
        assertTrue(movieContent.isVideo())
        assertFalse(movieContent.isBook())
        assertFalse(movieContent.isAudio())

        // Test duration formatting
        val durationString = movieContent.getDurationString()
        assertTrue(durationString.isNotEmpty())
        assertTrue(durationString.contains(":"))
    }

    /**
     * FLOW TEST 8: Error Handling Flow
     * Tests error conditions and edge cases
     */
    @Test
    fun testErrorHandling() = runTest {
        // Test invalid category ID
        val invalidCategoryResult = mockApiClient.getCategoryDetail("invalid")
        assertTrue(invalidCategoryResult.isFailure)

        // Test empty content request
        val emptyContentResponse = mockApiClient.getContent(
            StreamContentRequest("nonexistent", page = 1, limit = 10)
        )
        assertTrue(emptyContentResponse.success)
        assertEquals(0, emptyContentResponse.total)
        assertTrue(emptyContentResponse.content.isEmpty())

        // Test out-of-bounds pagination
        val outOfBoundsResponse = mockApiClient.getContent(
            StreamContentRequest("books", page = 999, limit = 10)
        )
        assertTrue(outOfBoundsResponse.success)
        assertTrue(outOfBoundsResponse.content.isEmpty())
    }

    /**
     * FLOW TEST 9: Performance Flow
     * Tests performance characteristics meet requirements
     */
    @Test
    fun testPerformanceRequirements() = runTest {
        val startTime = kotlinx.datetime.Clock.System.now()

        // Test category loading performance
        val categoriesResponse = mockApiClient.getCategories()
        assertTrue(categoriesResponse.success)

        val categoryLoadTime = kotlinx.datetime.Clock.System.now() - startTime

        // Should complete quickly (simulated network delay is 50ms)
        assertTrue(categoryLoadTime.inWholeMilliseconds < 1000) // <1s requirement

        // Test content loading performance
        val contentStartTime = kotlinx.datetime.Clock.System.now()

        val contentResponse = mockApiClient.getContent(
            StreamContentRequest("movies", page = 1, limit = 20)
        )
        assertTrue(contentResponse.success)

        val contentLoadTime = kotlinx.datetime.Clock.System.now() - contentStartTime

        // Should meet content loading requirements
        assertTrue(contentLoadTime.inWholeMilliseconds < 1000) // <1s requirement
    }

    /**
     * FLOW TEST 10: Cross-Platform State Flow
     * Tests navigation state management
     */
    @Test
    fun testNavigationStateManagement() = runTest {
        val userId = "test-user-123"
        val sessionId = "test-session-456"

        // Create navigation state
        val navigationState = TabNavigationState(
            userId = userId,
            currentCategoryId = "movies",
            currentSubtabId = "movies_short",
            lastVisitedAt = "2024-01-01T12:00:00Z",
            sessionId = sessionId
        )

        // Validate navigation state structure
        assertEquals(userId, navigationState.userId)
        assertEquals("movies", navigationState.currentCategoryId)
        assertEquals("movies_short", navigationState.currentSubtabId)
        assertEquals(sessionId, navigationState.sessionId)

        // Test state transitions
        val updatedState = navigationState.copy(
            currentCategoryId = "books",
            currentSubtabId = null,
            lastVisitedAt = "2024-01-01T12:05:00Z"
        )

        assertEquals("books", updatedState.currentCategoryId)
        assertNull(updatedState.currentSubtabId)
        assertEquals(userId, updatedState.userId) // Should remain unchanged
    }

    /**
     * FLOW TEST 11: Content Filtering Flow
     * Tests content filtering and availability
     */
    @Test
    fun testContentFiltering() = runTest {
        val contentResponse = mockApiClient.getContent(
            StreamContentRequest("movies", page = 1, limit = 10)
        )

        // All returned content should be available
        contentResponse.content.forEach { content ->
            assertEquals(StreamAvailabilityStatus.AVAILABLE, content.availabilityStatus)
            assertTrue(content.isAvailable())
            assertTrue(content.canPurchase())
        }

        // Test featured content filtering
        val allContent = contentResponse.content
        val featuredContent = allContent.filter { it.isFeatured }
        val regularContent = allContent.filter { !it.isFeatured }

        assertTrue(featuredContent.isNotEmpty())
        assertTrue(regularContent.isNotEmpty())

        // Featured content should have featured order
        featuredContent.forEach { content ->
            assertNotNull(content.featuredOrder)
        }
    }

    /**
     * FLOW TEST 12: Serialization Flow
     * Tests JSON serialization/deserialization
     */
    @Test
    fun testSerializationConsistency() {
        val originalCategory = mockCategories[0]

        // Serialize to JSON
        val categoryJson = json.encodeToString(StreamCategory.serializer(), originalCategory)
        assertNotNull(categoryJson)
        assertTrue(categoryJson.contains("\"id\":\"books\""))

        // Deserialize from JSON
        val deserializedCategory = json.decodeFromString(StreamCategory.serializer(), categoryJson)
        assertEquals(originalCategory.id, deserializedCategory.id)
        assertEquals(originalCategory.name, deserializedCategory.name)
        assertEquals(originalCategory.displayOrder, deserializedCategory.displayOrder)

        // Test content serialization
        val originalContent = mockStreamContent[0]
        val contentJson = json.encodeToString(StreamContentItem.serializer(), originalContent)
        val deserializedContent = json.decodeFromString(StreamContentItem.serializer(), contentJson)

        assertEquals(originalContent.id, deserializedContent.id)
        assertEquals(originalContent.title, deserializedContent.title)
        assertEquals(originalContent.price, deserializedContent.price)
        assertEquals(originalContent.contentType, deserializedContent.contentType)
    }
}