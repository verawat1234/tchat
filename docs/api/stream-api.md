# Stream API Documentation

**Version:** 1.0.0
**Base URL:** `/api/v1/stream`
**Service:** Commerce Service
**Gateway:** Port 8080 → Commerce Service

## Overview

The Stream API provides comprehensive access to the Stream Store Tabs functionality, supporting 6 content categories (Books, Podcasts, Cartoons, Movies, Music, Art) with advanced browsing, searching, purchasing, and user state management capabilities.

### Key Features

- **Content Discovery**: Browse categories, featured content, and subtab filtering
- **Advanced Search**: Full-text search with category filtering and pagination
- **Purchase Integration**: Unified shopping cart for both physical and digital products
- **User State Management**: Navigation state persistence and content view tracking
- **Performance Optimization**: <1s content load, <200ms API response time
- **Cross-Platform Support**: Web (React/RTK Query) and Mobile (KMP/SQLDelight) clients

## Authentication

Authentication is required for user-specific endpoints (purchasing, navigation state, preferences).

**Method:** Bearer Token (JWT)
**Header:** `Authorization: Bearer <token>`
**Test Header:** `X-User-ID: <user-id>` (development only)

## Content Categories

The Stream API supports 6 predefined content categories:

| Category | ID | Content Types | Icon |
|----------|----|--------------|----- |
| Books | `books` | E-books, Audiobooks | 📚 |
| Podcasts | `podcasts` | Audio shows, Episodes | 🎙️ |
| Cartoons | `cartoons` | Animated series, Shorts | 🎨 |
| Movies | `movies` | Short films, Feature films | 🎬 |
| Music | `music` | Songs, Albums, Playlists | 🎵 |
| Art | `art` | Digital art, NFTs, Galleries | 🖼️ |

## Endpoints

### Categories

#### GET /categories

Retrieves all active stream categories with display order.

**Parameters:** None

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/stream/categories" \
  -H "Content-Type: application/json"
```

**Example Response:**
```json
{
  "categories": [
    {
      "id": "books",
      "name": "Books",
      "displayOrder": 1,
      "iconName": "book",
      "isActive": true,
      "featuredContentEnabled": true,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    },
    {
      "id": "podcasts",
      "name": "Podcasts",
      "displayOrder": 2,
      "iconName": "podcast",
      "isActive": true,
      "featuredContentEnabled": true,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 6,
  "success": true
}
```

**Response Fields:**
- `categories` (array): List of category objects
- `total` (integer): Total number of categories
- `success` (boolean): Operation success status

#### GET /categories/:id

Retrieves detailed information about a specific category including subtabs and statistics.

**Parameters:**
- `id` (path, required): Category ID (e.g., "books", "podcasts")

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/stream/categories/books" \
  -H "Content-Type: application/json"
```

**Example Response:**
```json
{
  "category": {
    "id": "books",
    "name": "Books",
    "displayOrder": 1,
    "iconName": "book",
    "isActive": true,
    "featuredContentEnabled": true,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "subtabs": [
    {
      "id": "bestsellers",
      "categoryId": "books",
      "name": "Bestsellers",
      "displayOrder": 1,
      "filterCriteria": {
        "featured": true,
        "minRating": 4.5
      },
      "isActive": true,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "stats": {
    "totalItems": 150,
    "featuredItems": 12,
    "avgRating": 4.3
  },
  "success": true
}
```

**Error Responses:**
- `404 Not Found`: Category does not exist
- `400 Bad Request`: Invalid category ID format

### Content Browsing

#### GET /content

Retrieves paginated content items for a specific category with optional subtab filtering.

**Parameters:**
- `categoryId` (query, required): Category ID to filter content
- `subtabId` (query, optional): Subtab ID for additional filtering
- `page` (query, optional): Page number (default: 1, min: 1)
- `limit` (query, optional): Items per page (default: 20, min: 1, max: 100)

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/stream/content?categoryId=books&page=1&limit=10" \
  -H "Content-Type: application/json"
```

**Example Response:**
```json
{
  "items": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "categoryId": "books",
      "title": "The Great Gatsby",
      "description": "A classic American novel about the Jazz Age and the American Dream.",
      "thumbnailUrl": "https://example.com/images/great-gatsby.jpg",
      "contentType": "book",
      "price": 12.99,
      "currency": "USD",
      "availabilityStatus": "available",
      "isFeatured": true,
      "featuredOrder": 1,
      "metadata": {
        "author": "F. Scott Fitzgerald",
        "publishYear": 1925,
        "pages": 180,
        "isbn": "978-0-7432-7356-5"
      },
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 150,
  "hasMore": true,
  "success": true
}
```

**Content Types:**
- `book`: E-books and audiobooks
- `podcast`: Podcast episodes and series
- `cartoon`: Animated content
- `short_movie`: Short films (<30 minutes)
- `long_movie`: Feature films (≥30 minutes)
- `music`: Songs and albums
- `art`: Digital artwork and NFTs

**Availability Status:**
- `available`: Ready for purchase and consumption
- `coming_soon`: Available for pre-order or wishlist
- `unavailable`: Temporarily or permanently unavailable

#### GET /featured

Retrieves featured content for a specific category with priority ordering.

**Parameters:**
- `categoryId` (query, required): Category ID to filter featured content
- `limit` (query, optional): Number of featured items (default: 10, min: 1, max: 50)

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/stream/featured?categoryId=books&limit=5" \
  -H "Content-Type: application/json"
```

**Example Response:**
```json
{
  "items": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "categoryId": "books",
      "title": "Featured Novel",
      "description": "Editor's choice for this month.",
      "thumbnailUrl": "https://example.com/images/featured-novel.jpg",
      "contentType": "book",
      "price": 15.99,
      "currency": "USD",
      "availabilityStatus": "available",
      "isFeatured": true,
      "featuredOrder": 1,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 5,
  "hasMore": false,
  "success": true
}
```

#### GET /content/:id

Retrieves detailed information about a specific content item and tracks view for authenticated users.

**Parameters:**
- `id` (path, required): Content item UUID

**Headers (Optional):**
- `X-User-ID`: User ID for view tracking
- `X-Session-ID`: Session ID for analytics

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/stream/content/550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123"
```

**Example Response:**
```json
{
  "content": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "categoryId": "books",
    "title": "The Great Gatsby",
    "description": "A detailed description of this classic American novel...",
    "thumbnailUrl": "https://example.com/images/great-gatsby.jpg",
    "contentType": "book",
    "price": 12.99,
    "currency": "USD",
    "availabilityStatus": "available",
    "isFeatured": true,
    "featuredOrder": 1,
    "metadata": {
      "author": "F. Scott Fitzgerald",
      "publishYear": 1925,
      "pages": 180,
      "isbn": "978-0-7432-7356-5",
      "genres": ["Fiction", "Classic Literature"],
      "downloadFormats": ["PDF", "EPUB"],
      "preview": "https://example.com/previews/great-gatsby-sample.pdf"
    },
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "success": true
}
```

**Error Responses:**
- `404 Not Found`: Content item does not exist
- `400 Bad Request`: Invalid content ID format

### Search

#### GET /search

Performs full-text search across content titles, descriptions, and metadata with optional category filtering.

**Parameters:**
- `q` (query, required): Search query string
- `categoryId` (query, optional): Category ID to limit search scope
- `page` (query, optional): Page number (default: 1, min: 1)
- `limit` (query, optional): Items per page (default: 20, min: 1, max: 100)

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/stream/search?q=gatsby&categoryId=books&page=1&limit=10" \
  -H "Content-Type: application/json"
```

**Example Response:**
```json
{
  "items": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "categoryId": "books",
      "title": "The Great Gatsby",
      "description": "A classic American novel...",
      "thumbnailUrl": "https://example.com/images/great-gatsby.jpg",
      "contentType": "book",
      "price": 12.99,
      "currency": "USD",
      "availabilityStatus": "available",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "hasMore": false,
  "success": true
}
```

**Search Features:**
- Full-text search across title, description, and metadata
- Stemming and fuzzy matching support
- Category-specific filtering
- Relevance-based ranking
- Real-time search suggestions (future enhancement)

### Purchasing

#### POST /content/purchase

Processes a content purchase request with unified cart integration.

**Authentication:** Required
**Headers:**
- `Authorization: Bearer <jwt-token>`
- `Content-Type: application/json`

**Request Body:**
```json
{
  "mediaContentId": "550e8400-e29b-41d4-a716-446655440000",
  "quantity": 1,
  "mediaLicense": "personal",
  "downloadFormat": "PDF",
  "cartId": "cart_abc123"
}
```

**Request Fields:**
- `mediaContentId` (string, required): UUID of content to purchase
- `quantity` (integer, required): Number of items (min: 1)
- `mediaLicense` (string, required): License type - "personal", "family", "commercial"
- `downloadFormat` (string, optional): Preferred format - "PDF", "EPUB", "MP3", "MP4", "FLAC"
- `cartId` (string, optional): Shopping cart ID for unified cart integration

**Example Request:**
```bash
curl -X POST "http://localhost:8080/api/v1/stream/content/purchase" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "mediaContentId": "550e8400-e29b-41d4-a716-446655440000",
    "quantity": 1,
    "mediaLicense": "personal",
    "downloadFormat": "PDF"
  }'
```

**Example Response:**
```json
{
  "orderId": "order_789xyz",
  "totalAmount": 12.99,
  "currency": "USD",
  "success": true,
  "message": "Purchase completed successfully"
}
```

**Error Responses:**
- `401 Unauthorized`: Authentication required
- `400 Bad Request`: Invalid request format or content unavailable
- `409 Conflict`: Content already owned by user

### User State Management

#### GET /navigation

Retrieves the current user's navigation state including active category and subtab.

**Authentication:** Required
**Headers:** `Authorization: Bearer <jwt-token>`

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/stream/navigation" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Example Response:**
```json
{
  "state": {
    "userId": "user123",
    "sessionId": "session_abc123",
    "currentCategoryId": "books",
    "currentSubtabId": "bestsellers",
    "lastUpdated": "2024-01-01T12:00:00Z"
  },
  "success": true
}
```

#### PUT /navigation

Updates the user's navigation state when switching categories or subtabs.

**Authentication:** Required
**Headers:**
- `Authorization: Bearer <jwt-token>`
- `Content-Type: application/json`

**Request Body:**
```json
{
  "categoryId": "podcasts",
  "subtabId": "new-releases"
}
```

**Request Fields:**
- `categoryId` (string, required): Target category ID
- `subtabId` (string, optional): Target subtab ID (null for category main view)

**Example Request:**
```bash
curl -X PUT "http://localhost:8080/api/v1/stream/navigation" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "categoryId": "podcasts",
    "subtabId": "new-releases"
  }'
```

**Example Response:**
```json
{
  "state": {
    "userId": "user123",
    "sessionId": "session_abc123",
    "currentCategoryId": "podcasts",
    "currentSubtabId": "new-releases",
    "lastUpdated": "2024-01-01T12:30:00Z"
  },
  "success": true
}
```

#### PUT /content/:id/progress

Updates the user's content viewing progress for media with duration (videos, audio).

**Authentication:** Required
**Parameters:**
- `id` (path, required): Content item UUID

**Headers:**
- `Authorization: Bearer <jwt-token>`
- `Content-Type: application/json`

**Request Body:**
```json
{
  "progress": 0.75
}
```

**Request Fields:**
- `progress` (number, required): Progress as decimal (0.0 - 1.0)

**Example Request:**
```bash
curl -X PUT "http://localhost:8080/api/v1/stream/content/550e8400-e29b-41d4-a716-446655440000/progress" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "progress": 0.75
  }'
```

**Example Response:**
```json
{
  "success": true,
  "message": "Progress updated successfully"
}
```

#### GET /preferences

Retrieves the user's content preferences and settings.

**Authentication:** Required
**Headers:** `Authorization: Bearer <jwt-token>`

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/stream/preferences" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Example Response:**
```json
{
  "preferences": {
    "autoplay": true,
    "defaultQuality": "high",
    "downloadFormat": "PDF",
    "notifications": {
      "newReleases": true,
      "personalizedRecommendations": true,
      "promotions": false
    },
    "parentalControls": {
      "enabled": false,
      "maxRating": "PG-13"
    }
  },
  "success": true
}
```

#### PUT /preferences

Updates the user's content preferences and settings.

**Authentication:** Required
**Headers:**
- `Authorization: Bearer <jwt-token>`
- `Content-Type: application/json`

**Request Body:** Complete preferences object (same structure as GET response)

**Example Request:**
```bash
curl -X PUT "http://localhost:8080/api/v1/stream/preferences" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "autoplay": false,
    "defaultQuality": "medium",
    "downloadFormat": "EPUB",
    "notifications": {
      "newReleases": true,
      "personalizedRecommendations": true,
      "promotions": true
    }
  }'
```

**Example Response:**
```json
{
  "success": true,
  "message": "Preferences updated successfully"
}
```

## Health Checks

#### GET /health

General service health check.

**Example Request:**
```bash
curl -X GET "http://localhost:8080/health"
```

**Example Response:**
```json
{
  "status": "healthy",
  "service": "stream",
  "version": "1.0.0"
}
```

#### GET /health/stream

Detailed stream service health check with endpoint information.

**Example Request:**
```bash
curl -X GET "http://localhost:8080/health/stream"
```

**Example Response:**
```json
{
  "status": "healthy",
  "service": "stream",
  "component": "stream-content-management",
  "version": "1.0.0",
  "endpoints": {
    "categories": "/api/v1/stream/categories",
    "content": "/api/v1/stream/content",
    "featured": "/api/v1/stream/featured",
    "search": "/api/v1/stream/search",
    "purchase": "/api/v1/stream/content/purchase",
    "navigation": "/api/v1/stream/navigation",
    "preferences": "/api/v1/stream/preferences"
  }
}
```

## Error Handling

All endpoints return consistent error responses:

**Error Response Format:**
```json
{
  "success": false,
  "error": "Error description"
}
```

**HTTP Status Codes:**
- `200 OK`: Request successful
- `400 Bad Request`: Invalid request format or parameters
- `401 Unauthorized`: Authentication required or invalid
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict (e.g., already owned)
- `500 Internal Server Error`: Server error

## Performance Specifications

**Response Time Targets:**
- Content listing: <200ms
- Content search: <500ms
- Content detail: <100ms
- Purchase processing: <1s

**Pagination Limits:**
- Default page size: 20 items
- Maximum page size: 100 items
- Maximum search results: 1000 items

**Caching:**
- Categories: 1 hour TTL
- Featured content: 15 minutes TTL
- Content listings: 5 minutes TTL
- User state: No caching (real-time)

## Integration Examples

### Web Client (React + RTK Query)

```typescript
// hooks/streamApi.ts
import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'

export const streamApi = createApi({
  reducerPath: 'streamApi',
  baseQuery: fetchBaseQuery({
    baseUrl: '/api/v1/stream',
    prepareHeaders: (headers, { getState }) => {
      const token = getState().auth.token
      if (token) {
        headers.set('authorization', `Bearer ${token}`)
      }
      return headers
    },
  }),
  tagTypes: ['Categories', 'Content', 'Navigation'],
  endpoints: (builder) => ({
    getCategories: builder.query<CategoriesResponse, void>({
      query: () => '/categories',
      providesTags: ['Categories'],
    }),
    getContent: builder.query<ContentResponse, ContentParams>({
      query: ({ categoryId, page = 1, limit = 20, subtabId }) => ({
        url: '/content',
        params: { categoryId, page, limit, subtabId },
      }),
      providesTags: ['Content'],
    }),
    purchaseContent: builder.mutation<PurchaseResponse, PurchaseRequest>({
      query: (body) => ({
        url: '/content/purchase',
        method: 'POST',
        body,
      }),
    }),
  }),
})

export const {
  useGetCategoriesQuery,
  useGetContentQuery,
  usePurchaseContentMutation,
} = streamApi
```

### Mobile Client (KMP + SQLDelight)

```kotlin
// StreamApiClient.kt
class StreamApiClient(
    private val httpClient: HttpClient,
    private val database: StreamDatabase
) {
    suspend fun getCategories(): List<StreamCategory> {
        return try {
            val response = httpClient.get("/api/v1/stream/categories")
            val categories = response.body<CategoriesResponse>().categories

            // Cache in SQLDelight database
            database.streamCategoryQueries.transaction {
                categories.forEach { category ->
                    database.streamCategoryQueries.insertOrReplace(category)
                }
            }

            categories
        } catch (e: Exception) {
            // Fallback to cached data
            database.streamCategoryQueries.selectAll().executeAsList()
        }
    }

    suspend fun getContent(
        categoryId: String,
        page: Int = 1,
        limit: Int = 20,
        subtabId: String? = null
    ): ContentResponse {
        return httpClient.get("/api/v1/stream/content") {
            parameter("categoryId", categoryId)
            parameter("page", page)
            parameter("limit", limit)
            subtabId?.let { parameter("subtabId", it) }
        }.body()
    }
}
```

## Version History

- **v1.0.0** (2024-01-01): Initial release with full Stream Store Tabs functionality
- Categories, content browsing, search, purchase, and user state management
- Cross-platform support (Web + Mobile KMP)
- Performance optimization and caching
- Unified shopping cart integration

## Support

For API support and integration assistance:
- **Technical Documentation**: `/docs/api/stream-integration-guide.md`
- **OpenAPI Specification**: Available via `/api/v1/stream/openapi.json`
- **Performance Monitoring**: Real-time metrics at `/health/stream`