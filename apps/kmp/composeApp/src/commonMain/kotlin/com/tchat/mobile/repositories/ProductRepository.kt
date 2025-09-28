package com.tchat.mobile.repositories

import com.tchat.mobile.database.TchatDatabase
import com.tchat.mobile.models.*
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import kotlinx.datetime.Clock
import kotlinx.datetime.Instant

/**
 * Repository for product and store operations
 * Provides abstraction layer between UI and data layer
 */
class ProductRepository(
    private val database: TchatDatabase
) {
    private val productQueries = database.productQueries

    // Product operations
    suspend fun insertProduct(product: Product): Result<Unit> {
        return try {
            productQueries.insertProduct(
                id = product.id,
                name = product.name,
                description = product.description,
                short_description = product.shortDescription,
                sku = product.sku,
                category = product.category,
                brand = product.brand,
                price = product.price,
                original_price = product.originalPrice,
                currency = product.currency,
                thumbnail = product.thumbnail,
                availability = product.availability.name,
                stock = product.stock.toLong(),
                min_order_quantity = product.minOrderQuantity.toLong(),
                max_order_quantity = product.maxOrderQuantity.toLong(),
                weight = product.weight,
                rating = product.rating,
                review_count = product.reviewCount.toLong(),
                is_digital = if (product.isDigital) 1L else 0L,
                shipping_required = if (product.shippingRequired) 1L else 0L,
                taxable = if (product.taxable) 1L else 0L,
                status = product.status.name,
                seller_id = product.sellerId,
                store_name = product.storeName,
                created_at = product.createdAt.toEpochMilliseconds(),
                updated_at = product.updatedAt.toEpochMilliseconds(),
                delivery_time = product.deliveryTime,
                distance = product.distance,
                is_hot = if (product.isHot) 1L else 0L,
                orders = product.orders.toLong(),
                tags = product.tags.joinToString(","),
                images = product.images.joinToString(",")
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getProductById(id: String): Result<Product?> {
        return try {
            val productEntity = productQueries.getProductById(id).executeAsOneOrNull()
            val product = productEntity?.let { mapToProduct(it) }
            Result.success(product)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getAllProducts(): Result<List<Product>> {
        return try {
            val products = productQueries.getAllProducts().executeAsList().map { mapToProduct(it) }
            Result.success(products)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getProductsByCategory(category: String): Result<List<Product>> {
        return try {
            val products = productQueries.getProductsByCategory(category).executeAsList().map { mapToProduct(it) }
            Result.success(products)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getProductsByStore(sellerId: String): Result<List<Product>> {
        return try {
            val products = productQueries.getProductsByStore(sellerId).executeAsList().map { mapToProduct(it) }
            Result.success(products)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getHotProducts(): Result<List<Product>> {
        return try {
            val products = productQueries.getHotProducts().executeAsList().map { mapToProduct(it) }
            Result.success(products)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getTopRatedProducts(limit: Long = 10): Result<List<Product>> {
        return try {
            val products = productQueries.getTopRatedProducts(limit).executeAsList().map { mapToProduct(it) }
            Result.success(products)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun searchProducts(query: String): Result<List<Product>> {
        return try {
            val products = productQueries.searchProducts(query, query, query).executeAsList().map { mapToProduct(it) }
            Result.success(products)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Store operations
    suspend fun insertStore(store: Store): Result<Unit> {
        return try {
            productQueries.insertStore(
                id = store.id,
                name = store.name,
                description = store.description,
                avatar = store.avatar,
                cover_image = store.coverImage,
                rating = store.rating,
                delivery_time = store.deliveryTime,
                distance = store.distance,
                is_verified = if (store.isVerified) 1L else 0L,
                followers = store.followers.toLong(),
                total_products = store.totalProducts.toLong(),
                created_at = store.createdAt.toEpochMilliseconds(),
                updated_at = store.updatedAt.toEpochMilliseconds()
            )
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getStoreById(id: String): Result<Store?> {
        return try {
            val storeEntity = productQueries.getStoreById(id).executeAsOneOrNull()
            val store = storeEntity?.let { mapToStore(it) }
            Result.success(store)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getAllStores(): Result<List<Store>> {
        return try {
            val stores = productQueries.getAllStores().executeAsList().map { mapToStore(it) }
            Result.success(stores)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getVerifiedStores(): Result<List<Store>> {
        return try {
            val stores = productQueries.getVerifiedStores().executeAsList().map { mapToStore(it) }
            Result.success(stores)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun searchStores(query: String): Result<List<Store>> {
        return try {
            val stores = productQueries.searchStores(query, query).executeAsList().map { mapToStore(it) }
            Result.success(stores)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // UI compatibility methods
    suspend fun getProductItems(): Result<List<ProductItem>> {
        return try {
            val products = getAllProducts().getOrThrow()
            Result.success(products.map { it.toProductItem() })
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getShopItems(): Result<List<ShopItem>> {
        return try {
            val stores = getAllStores().getOrThrow()
            Result.success(stores.map { it.toShopItem() })
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getProductItemsByCategory(category: String): Result<List<ProductItem>> {
        return try {
            val products = getProductsByCategory(category).getOrThrow()
            Result.success(products.map { it.toProductItem() })
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun searchProductItems(query: String): Result<List<ProductItem>> {
        return try {
            val products = searchProducts(query).getOrThrow()
            Result.success(products.map { it.toProductItem() })
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun searchShopItems(query: String): Result<List<ShopItem>> {
        return try {
            val stores = searchStores(query).getOrThrow()
            Result.success(stores.map { it.toShopItem() })
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Seed data for initial population
    suspend fun seedInitialData(): Result<Unit> {
        return try {
            // Seed stores first
            val stores = getSeedStores()
            stores.forEach { store ->
                insertStore(store)
            }

            // Seed products
            val products = getSeedProducts()
            products.forEach { product ->
                insertProduct(product)
            }

            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Private helper methods
    private fun mapToProduct(entity: com.tchat.mobile.database.Product): Product {
        return Product(
            id = entity.id,
            name = entity.name,
            description = entity.description,
            shortDescription = entity.short_description,
            sku = entity.sku,
            category = entity.category,
            brand = entity.brand,
            price = entity.price,
            originalPrice = entity.original_price,
            currency = entity.currency,
            thumbnail = entity.thumbnail,
            availability = ProductAvailability.valueOf(entity.availability),
            stock = entity.stock.toInt(),
            minOrderQuantity = entity.min_order_quantity.toInt(),
            maxOrderQuantity = entity.max_order_quantity.toInt(),
            weight = entity.weight,
            rating = entity.rating,
            reviewCount = entity.review_count.toInt(),
            isDigital = entity.is_digital == 1L,
            shippingRequired = entity.shipping_required == 1L,
            taxable = entity.taxable == 1L,
            status = ProductStatus.valueOf(entity.status),
            sellerId = entity.seller_id,
            storeName = entity.store_name,
            createdAt = Instant.fromEpochMilliseconds(entity.created_at),
            updatedAt = Instant.fromEpochMilliseconds(entity.updated_at),
            deliveryTime = entity.delivery_time,
            distance = entity.distance,
            isHot = entity.is_hot == 1L,
            orders = entity.orders.toInt(),
            tags = if (entity.tags.isNotEmpty()) entity.tags.split(",") else emptyList(),
            images = if (entity.images.isNotEmpty()) entity.images.split(",") else emptyList()
        )
    }

    private fun mapToStore(entity: com.tchat.mobile.database.Store): Store {
        return Store(
            id = entity.id,
            name = entity.name,
            description = entity.description,
            avatar = entity.avatar,
            coverImage = entity.cover_image,
            rating = entity.rating,
            deliveryTime = entity.delivery_time,
            distance = entity.distance,
            isVerified = entity.is_verified == 1L,
            followers = entity.followers.toInt(),
            totalProducts = entity.total_products.toInt(),
            createdAt = Instant.fromEpochMilliseconds(entity.created_at),
            updatedAt = Instant.fromEpochMilliseconds(entity.updated_at)
        )
    }

    // Seed data functions
    private fun getSeedStores(): List<Store> {
        val now = Clock.System.now()
        return listOf(
            Store(
                id = "store_1",
                name = "Bangkok Street Food",
                description = "Authentic Thai street food and snacks",
                avatar = "BS",
                coverImage = "cover1",
                rating = 4.8,
                deliveryTime = "15-25 min",
                distance = "1.2 km",
                isVerified = true,
                followers = 2540,
                totalProducts = 45,
                createdAt = now,
                updatedAt = now
            ),
            Store(
                id = "store_2",
                name = "Tech Paradise",
                description = "Latest electronics and gadgets",
                avatar = "TP",
                coverImage = "cover2",
                rating = 4.6,
                deliveryTime = "30-45 min",
                distance = "3.1 km",
                isVerified = true,
                followers = 1890,
                totalProducts = 120,
                createdAt = now,
                updatedAt = now
            ),
            Store(
                id = "store_3",
                name = "Fashion House",
                description = "Trendy clothes and accessories",
                avatar = "FH",
                coverImage = "cover3",
                rating = 4.5,
                deliveryTime = "20-30 min",
                distance = "2.8 km",
                isVerified = false,
                followers = 980,
                totalProducts = 78,
                createdAt = now,
                updatedAt = now
            ),
            Store(
                id = "store_4",
                name = "Fresh Market",
                description = "Organic fruits and vegetables",
                avatar = "FM",
                coverImage = "cover4",
                rating = 4.7,
                deliveryTime = "10-20 min",
                distance = "0.8 km",
                isVerified = true,
                followers = 3200,
                totalProducts = 230,
                createdAt = now,
                updatedAt = now
            ),
            Store(
                id = "store_5",
                name = "Coffee Corner",
                description = "Premium coffee and desserts",
                avatar = "CC",
                coverImage = "cover5",
                rating = 4.9,
                deliveryTime = "5-15 min",
                distance = "0.5 km",
                isVerified = true,
                followers = 1560,
                totalProducts = 35,
                createdAt = now,
                updatedAt = now
            )
        )
    }

    private fun getSeedProducts(): List<Product> {
        val now = Clock.System.now()
        return listOf(
            Product(
                id = "product_1",
                name = "Pad Thai Goong",
                description = "Classic Thai stir-fried rice noodles with shrimp",
                sku = "PAD_THAI_001",
                category = "Food",
                price = 4500L, // 45.00 THB in cents
                originalPrice = 5500L, // 55.00 THB
                currency = "THB",
                thumbnail = "image1",
                rating = 4.8,
                reviewCount = 245,
                sellerId = "store_1",
                storeName = "Bangkok Street Food",
                deliveryTime = "15 min",
                distance = "1.2 km",
                isHot = true,
                orders = 245,
                createdAt = now,
                updatedAt = now
            ),
            Product(
                id = "product_2",
                name = "Wireless Headphones",
                description = "Premium noise-cancelling wireless headphones",
                sku = "WH_001",
                category = "Electronics",
                price = 249000L, // 2490.00 THB
                originalPrice = 299000L, // 2990.00 THB
                currency = "THB",
                thumbnail = "image2",
                rating = 4.6,
                reviewCount = 89,
                sellerId = "store_2",
                storeName = "Tech Paradise",
                deliveryTime = "30 min",
                distance = "3.1 km",
                isHot = false,
                orders = 89,
                createdAt = now,
                updatedAt = now
            ),
            Product(
                id = "product_3",
                name = "Cotton T-Shirt",
                description = "Comfortable cotton t-shirt in various colors",
                sku = "TSHIRT_001",
                category = "Fashion",
                price = 59000L, // 590.00 THB
                currency = "THB",
                thumbnail = "image3",
                rating = 4.3,
                reviewCount = 67,
                sellerId = "store_3",
                storeName = "Fashion House",
                deliveryTime = "25 min",
                distance = "2.8 km",
                isHot = false,
                orders = 67,
                createdAt = now,
                updatedAt = now
            ),
            Product(
                id = "product_4",
                name = "Fresh Mango",
                description = "Sweet and juicy organic mangoes from local farms",
                sku = "MANGO_001",
                category = "Food",
                price = 12000L, // 120.00 THB
                originalPrice = 15000L, // 150.00 THB
                currency = "THB",
                thumbnail = "image4",
                rating = 4.9,
                reviewCount = 156,
                sellerId = "store_4",
                storeName = "Fresh Market",
                deliveryTime = "15 min",
                distance = "0.8 km",
                isHot = true,
                orders = 156,
                createdAt = now,
                updatedAt = now
            ),
            Product(
                id = "product_5",
                name = "Iced Coffee",
                description = "Refreshing iced coffee made from premium beans",
                sku = "COFFEE_001",
                category = "Beverage",
                price = 8500L, // 85.00 THB
                currency = "THB",
                thumbnail = "image5",
                rating = 4.7,
                reviewCount = 98,
                sellerId = "store_5",
                storeName = "Coffee Corner",
                deliveryTime = "10 min",
                distance = "0.5 km",
                isHot = false,
                orders = 98,
                createdAt = now,
                updatedAt = now
            ),
            Product(
                id = "product_6",
                name = "Som Tam",
                description = "Spicy Thai green papaya salad",
                sku = "SOMTAM_001",
                category = "Food",
                price = 3500L, // 35.00 THB
                currency = "THB",
                thumbnail = "image6",
                rating = 4.5,
                reviewCount = 189,
                sellerId = "store_1",
                storeName = "Bangkok Street Food",
                deliveryTime = "15 min",
                distance = "1.2 km",
                isHot = true,
                orders = 189,
                createdAt = now,
                updatedAt = now
            ),
            Product(
                id = "product_7",
                name = "Smart Watch",
                description = "Feature-rich smartwatch with health tracking",
                sku = "WATCH_001",
                category = "Electronics",
                price = 899000L, // 8990.00 THB
                originalPrice = 1099000L, // 10990.00 THB
                currency = "THB",
                thumbnail = "image7",
                rating = 4.4,
                reviewCount = 45,
                sellerId = "store_2",
                storeName = "Tech Paradise",
                deliveryTime = "30 min",
                distance = "3.1 km",
                isHot = false,
                orders = 45,
                createdAt = now,
                updatedAt = now
            ),
            Product(
                id = "product_8",
                name = "Denim Jeans",
                description = "Classic denim jeans with modern fit",
                sku = "JEANS_001",
                category = "Fashion",
                price = 129000L, // 1290.00 THB
                originalPrice = 159000L, // 1590.00 THB
                currency = "THB",
                thumbnail = "image8",
                rating = 4.2,
                reviewCount = 78,
                sellerId = "store_3",
                storeName = "Fashion House",
                deliveryTime = "25 min",
                distance = "2.8 km",
                isHot = false,
                orders = 78,
                createdAt = now,
                updatedAt = now
            )
        )
    }
}