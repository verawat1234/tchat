// Media Store Android Repository
// Generated for Media Store Tabs feature implementation

package com.tchat.app.data.repositories

import com.tchat.app.data.models.*
import com.tchat.app.data.network.MediaApiService
import com.tchat.app.data.local.MediaDao
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.catch
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class MediaRepository @Inject constructor(
    private val apiService: MediaApiService,
    private val localDao: MediaDao
) {

    // MARK: - Categories
    suspend fun getMediaCategories(forceRefresh: Boolean = false): Flow<Result<MediaCategoriesResponse>> = flow {
        try {
            // Check cache first if not forcing refresh
            if (!forceRefresh) {
                val cachedCategories = localDao.getAllCategories().first()
                if (cachedCategories.isNotEmpty()) {
                    emit(Result.success(MediaCategoriesResponse(cachedCategories, cachedCategories.size)))
                }
            }

            // Fetch from API
            val response = apiService.getMediaCategories()

            // Cache the results
            localDao.insertCategories(response.categories)

            emit(Result.success(response))
        } catch (e: Exception) {
            // Try to get cached data on error
            try {
                val cachedCategories = localDao.getAllCategories().first()
                if (cachedCategories.isNotEmpty()) {
                    emit(Result.success(MediaCategoriesResponse(cachedCategories, cachedCategories.size)))
                } else {
                    emit(Result.failure(e))
                }
            } catch (cacheError: Exception) {
                emit(Result.failure(e))
            }
        }
    }

    suspend fun getMediaCategory(categoryId: String, forceRefresh: Boolean = false): Flow<Result<MediaCategory>> = flow {
        try {
            // Check cache first if not forcing refresh
            if (!forceRefresh) {
                val cachedCategory = localDao.getCategoryById(categoryId).first()
                if (cachedCategory != null) {
                    emit(Result.success(cachedCategory))
                }
            }

            // Fetch from API
            val category = apiService.getMediaCategory(categoryId)

            // Cache the result
            localDao.insertCategories(listOf(category))

            emit(Result.success(category))
        } catch (e: Exception) {
            // Try to get cached data on error
            try {
                val cachedCategory = localDao.getCategoryById(categoryId).first()
                if (cachedCategory != null) {
                    emit(Result.success(cachedCategory))
                } else {
                    emit(Result.failure(e))
                }
            } catch (cacheError: Exception) {
                emit(Result.failure(e))
            }
        }
    }

    // MARK: - Subtabs
    suspend fun getMovieSubtabs(forceRefresh: Boolean = false): Flow<Result<MediaSubtabsResponse>> = flow {
        try {
            // Check cache first if not forcing refresh
            if (!forceRefresh) {
                val cachedSubtabs = localDao.getSubtabsByCategory("movies").first()
                if (cachedSubtabs.isNotEmpty()) {
                    emit(Result.success(MediaSubtabsResponse(cachedSubtabs, cachedSubtabs.firstOrNull()?.id ?: "")))
                }
            }

            // Fetch from API
            val response = apiService.getMovieSubtabs()

            // Cache the results
            localDao.insertSubtabs(response.subtabs)

            emit(Result.success(response))
        } catch (e: Exception) {
            // Try to get cached data on error
            try {
                val cachedSubtabs = localDao.getSubtabsByCategory("movies").first()
                if (cachedSubtabs.isNotEmpty()) {
                    emit(Result.success(MediaSubtabsResponse(cachedSubtabs, cachedSubtabs.firstOrNull()?.id ?: "")))
                } else {
                    emit(Result.failure(e))
                }
            } catch (cacheError: Exception) {
                emit(Result.failure(e))
            }
        }
    }

    // MARK: - Content
    suspend fun getContentByCategory(
        categoryId: String,
        page: Int = 1,
        limit: Int = 20,
        subtab: String? = null,
        forceRefresh: Boolean = false
    ): Flow<Result<MediaContentResponse>> = flow {
        try {
            // Check cache first if not forcing refresh and on first page
            if (!forceRefresh && page == 1) {
                val cachedContent = if (subtab != null) {
                    localDao.getContentByCategoryAndSubtab(categoryId, subtab).first()
                } else {
                    localDao.getContentByCategory(categoryId).first()
                }

                if (cachedContent.isNotEmpty()) {
                    emit(Result.success(MediaContentResponse(
                        items = cachedContent,
                        page = 1,
                        limit = cachedContent.size,
                        total = cachedContent.size,
                        hasMore = false
                    )))
                }
            }

            // Fetch from API
            val response = apiService.getContentByCategory(categoryId, page, limit, subtab)

            // Cache the results (only first page to avoid complex pagination caching)
            if (page == 1) {
                localDao.insertContent(response.items)
            }

            emit(Result.success(response))
        } catch (e: Exception) {
            // Try to get cached data on error (only for first page)
            if (page == 1) {
                try {
                    val cachedContent = if (subtab != null) {
                        localDao.getContentByCategoryAndSubtab(categoryId, subtab).first()
                    } else {
                        localDao.getContentByCategory(categoryId).first()
                    }

                    if (cachedContent.isNotEmpty()) {
                        emit(Result.success(MediaContentResponse(
                            items = cachedContent,
                            page = 1,
                            limit = cachedContent.size,
                            total = cachedContent.size,
                            hasMore = false
                        )))
                    } else {
                        emit(Result.failure(e))
                    }
                } catch (cacheError: Exception) {
                    emit(Result.failure(e))
                }
            } else {
                emit(Result.failure(e))
            }
        }
    }

    suspend fun getFeaturedContent(
        limit: Int = 10,
        categoryId: String? = null,
        forceRefresh: Boolean = false
    ): Flow<Result<MediaFeaturedResponse>> = flow {
        try {
            // Check cache first if not forcing refresh
            if (!forceRefresh) {
                val cachedFeatured = if (categoryId != null) {
                    localDao.getFeaturedContentByCategory(categoryId, limit).first()
                } else {
                    localDao.getFeaturedContent(limit).first()
                }

                if (cachedFeatured.isNotEmpty()) {
                    emit(Result.success(MediaFeaturedResponse(
                        items = cachedFeatured,
                        total = cachedFeatured.size,
                        hasMore = false
                    )))
                }
            }

            // Fetch from API
            val response = apiService.getFeaturedContent(limit, categoryId)

            // Cache the results
            localDao.insertContent(response.items)

            emit(Result.success(response))
        } catch (e: Exception) {
            // Try to get cached data on error
            try {
                val cachedFeatured = if (categoryId != null) {
                    localDao.getFeaturedContentByCategory(categoryId, limit).first()
                } else {
                    localDao.getFeaturedContent(limit).first()
                }

                if (cachedFeatured.isNotEmpty()) {
                    emit(Result.success(MediaFeaturedResponse(
                        items = cachedFeatured,
                        total = cachedFeatured.size,
                        hasMore = false
                    )))
                } else {
                    emit(Result.failure(e))
                }
            } catch (cacheError: Exception) {
                emit(Result.failure(e))
            }
        }
    }

    suspend fun searchMediaContent(
        query: String,
        categoryId: String? = null,
        page: Int = 1,
        limit: Int = 20
    ): Flow<Result<MediaSearchResponse>> = flow {
        try {
            // Search is always from API for fresh results
            val response = apiService.searchMediaContent(query, categoryId, page, limit)
            emit(Result.success(response))
        } catch (e: Exception) {
            emit(Result.failure(e))
        }
    }

    // MARK: - Store Integration
    suspend fun getMediaProducts(
        categoryId: String? = null,
        page: Int = 1,
        limit: Int = 20,
        forceRefresh: Boolean = false
    ): Flow<Result<MediaProductsResponse>> = flow {
        try {
            // Check cache first if not forcing refresh and on first page
            if (!forceRefresh && page == 1) {
                val cachedProducts = if (categoryId != null) {
                    localDao.getProductsByCategory(categoryId).first()
                } else {
                    localDao.getAllProducts().first()
                }

                if (cachedProducts.isNotEmpty()) {
                    emit(Result.success(MediaProductsResponse(
                        products = cachedProducts,
                        pagination = MediaProductsResponse.PaginationInfo(
                            page = 1,
                            limit = cachedProducts.size,
                            total = cachedProducts.size,
                            hasMore = false
                        )
                    )))
                }
            }

            // Fetch from API
            val response = apiService.getMediaProducts(categoryId, page, limit)

            // Cache the results (only first page)
            if (page == 1) {
                localDao.insertProducts(response.products)
            }

            emit(Result.success(response))
        } catch (e: Exception) {
            // Try to get cached data on error (only for first page)
            if (page == 1) {
                try {
                    val cachedProducts = if (categoryId != null) {
                        localDao.getProductsByCategory(categoryId).first()
                    } else {
                        localDao.getAllProducts().first()
                    }

                    if (cachedProducts.isNotEmpty()) {
                        emit(Result.success(MediaProductsResponse(
                            products = cachedProducts,
                            pagination = MediaProductsResponse.PaginationInfo(
                                page = 1,
                                limit = cachedProducts.size,
                                total = cachedProducts.size,
                                hasMore = false
                            )
                        )))
                    } else {
                        emit(Result.failure(e))
                    }
                } catch (cacheError: Exception) {
                    emit(Result.failure(e))
                }
            } else {
                emit(Result.failure(e))
            }
        }
    }

    // MARK: - Cart Operations
    suspend fun addMediaToCart(request: AddMediaToCartRequest): Flow<Result<AddMediaToCartResponse>> = flow {
        try {
            val response = apiService.addMediaToCart(request)

            // Update local cart cache
            localDao.insertCartItem(response.addedItem)

            emit(Result.success(response))
        } catch (e: Exception) {
            emit(Result.failure(e))
        }
    }

    suspend fun getUnifiedCart(forceRefresh: Boolean = false): Flow<Result<UnifiedCartResponse>> = flow {
        try {
            // Check cache first if not forcing refresh
            if (!forceRefresh) {
                val cachedCartItems = localDao.getAllCartItems().first()
                if (cachedCartItems.isNotEmpty()) {
                    val mediaItems = cachedCartItems.filter { it.mediaContentId != null }
                    val physicalItems = cachedCartItems.filter { it.mediaContentId == null }

                    emit(Result.success(UnifiedCartResponse(
                        cartId = cachedCartItems.firstOrNull()?.cartId ?: "",
                        physicalItems = physicalItems,
                        mediaItems = mediaItems,
                        totalPhysicalAmount = physicalItems.sumOf { it.totalPrice },
                        totalMediaAmount = mediaItems.sumOf { it.totalPrice },
                        totalAmount = cachedCartItems.sumOf { it.totalPrice },
                        currency = "USD",
                        itemsCount = cachedCartItems.sumOf { it.quantity }
                    )))
                }
            }

            // Fetch from API
            val response = apiService.getUnifiedCart()

            // Cache the results
            localDao.clearCartItems()
            localDao.insertCartItems(response.physicalItems + response.mediaItems)

            emit(Result.success(response))
        } catch (e: Exception) {
            // Try to get cached data on error
            try {
                val cachedCartItems = localDao.getAllCartItems().first()
                if (cachedCartItems.isNotEmpty()) {
                    val mediaItems = cachedCartItems.filter { it.mediaContentId != null }
                    val physicalItems = cachedCartItems.filter { it.mediaContentId == null }

                    emit(Result.success(UnifiedCartResponse(
                        cartId = cachedCartItems.firstOrNull()?.cartId ?: "",
                        physicalItems = physicalItems,
                        mediaItems = mediaItems,
                        totalPhysicalAmount = physicalItems.sumOf { it.totalPrice },
                        totalMediaAmount = mediaItems.sumOf { it.totalPrice },
                        totalAmount = cachedCartItems.sumOf { it.totalPrice },
                        currency = "USD",
                        itemsCount = cachedCartItems.sumOf { it.quantity }
                    )))
                } else {
                    emit(Result.failure(e))
                }
            } catch (cacheError: Exception) {
                emit(Result.failure(e))
            }
        }
    }

    suspend fun removeMediaFromCart(cartItemId: String): Flow<Result<Unit>> = flow {
        try {
            apiService.removeMediaFromCart(cartItemId)

            // Update local cache
            localDao.deleteCartItem(cartItemId)

            emit(Result.success(Unit))
        } catch (e: Exception) {
            emit(Result.failure(e))
        }
    }

    suspend fun updateMediaCartItem(cartItemId: String, quantity: Int): Flow<Result<MediaCartItem>> = flow {
        try {
            val response = apiService.updateMediaCartItem(cartItemId, quantity)

            // Update local cache
            localDao.insertCartItem(response)

            emit(Result.success(response))
        } catch (e: Exception) {
            emit(Result.failure(e))
        }
    }

    // MARK: - Checkout Operations
    suspend fun validateMediaCheckout(request: MediaCheckoutValidationRequest): Flow<Result<MediaCheckoutValidationResponse>> = flow {
        try {
            val response = apiService.validateMediaCheckout(request)
            emit(Result.success(response))
        } catch (e: Exception) {
            emit(Result.failure(e))
        }
    }

    suspend fun processMediaCheckout(request: ProcessMediaCheckoutRequest): Flow<Result<MediaOrder>> = flow {
        try {
            val response = apiService.processMediaCheckout(request)

            // Cache the order
            localDao.insertOrder(response)

            // Clear cart items (successful checkout)
            localDao.clearCartItems()

            emit(Result.success(response))
        } catch (e: Exception) {
            emit(Result.failure(e))
        }
    }

    // MARK: - Orders
    suspend fun getMediaOrders(
        page: Int = 1,
        limit: Int = 20,
        status: String? = null,
        forceRefresh: Boolean = false
    ): Flow<Result<MediaOrdersResponse>> = flow {
        try {
            // Check cache first if not forcing refresh and on first page
            if (!forceRefresh && page == 1) {
                val cachedOrders = if (status != null) {
                    localDao.getOrdersByStatus(status).first()
                } else {
                    localDao.getAllOrders().first()
                }

                if (cachedOrders.isNotEmpty()) {
                    emit(Result.success(MediaOrdersResponse(
                        orders = cachedOrders,
                        pagination = MediaOrdersResponse.PaginationInfo(
                            page = 1,
                            limit = cachedOrders.size,
                            total = cachedOrders.size,
                            hasMore = false
                        )
                    )))
                }
            }

            // Fetch from API
            val response = apiService.getMediaOrders(page, limit, status)

            // Cache the results (only first page)
            if (page == 1) {
                localDao.insertOrders(response.orders)
            }

            emit(Result.success(response))
        } catch (e: Exception) {
            // Try to get cached data on error (only for first page)
            if (page == 1) {
                try {
                    val cachedOrders = if (status != null) {
                        localDao.getOrdersByStatus(status).first()
                    } else {
                        localDao.getAllOrders().first()
                    }

                    if (cachedOrders.isNotEmpty()) {
                        emit(Result.success(MediaOrdersResponse(
                            orders = cachedOrders,
                            pagination = MediaOrdersResponse.PaginationInfo(
                                page = 1,
                                limit = cachedOrders.size,
                                total = cachedOrders.size,
                                hasMore = false
                            )
                        )))
                    } else {
                        emit(Result.failure(e))
                    }
                } catch (cacheError: Exception) {
                    emit(Result.failure(e))
                }
            } else {
                emit(Result.failure(e))
            }
        }
    }

    suspend fun getMediaOrder(orderId: String, forceRefresh: Boolean = false): Flow<Result<MediaOrder>> = flow {
        try {
            // Check cache first if not forcing refresh
            if (!forceRefresh) {
                val cachedOrder = localDao.getOrderById(orderId).first()
                if (cachedOrder != null) {
                    emit(Result.success(cachedOrder))
                }
            }

            // Fetch from API
            val order = apiService.getMediaOrder(orderId)

            // Cache the result
            localDao.insertOrder(order)

            emit(Result.success(order))
        } catch (e: Exception) {
            // Try to get cached data on error
            try {
                val cachedOrder = localDao.getOrderById(orderId).first()
                if (cachedOrder != null) {
                    emit(Result.success(cachedOrder))
                } else {
                    emit(Result.failure(e))
                }
            } catch (cacheError: Exception) {
                emit(Result.failure(e))
            }
        }
    }

    suspend fun downloadMediaContent(orderItemId: String): Flow<Result<MediaDownloadResponse>> = flow {
        try {
            val response = apiService.downloadMediaContent(orderItemId)
            emit(Result.success(response))
        } catch (e: Exception) {
            emit(Result.failure(e))
        }
    }

    // MARK: - Cache Management
    suspend fun clearAllCache() {
        localDao.clearAllTables()
    }

    suspend fun clearContentCache() {
        localDao.clearContent()
    }

    suspend fun clearCartCache() {
        localDao.clearCartItems()
    }
}

// MARK: - Additional Response Types
@Serializable
data class MediaProductsResponse(
    val products: List<MediaProduct>,
    val pagination: PaginationInfo
) {
    @Serializable
    data class PaginationInfo(
        val page: Int,
        val limit: Int,
        val total: Int,
        val hasMore: Boolean
    )
}

@Serializable
data class ProcessMediaCheckoutRequest(
    val cartId: String,
    val mediaItems: List<MediaCartItem>,
    val paymentMethod: String,
    val billingAddress: String? = null
)

@Serializable
data class MediaDownloadResponse(
    val downloadUrl: String,
    val expiresAt: String
)