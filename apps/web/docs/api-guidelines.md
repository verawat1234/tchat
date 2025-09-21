# API Design Guidelines for Telegram SEA Edition

This document outlines the comprehensive API design guidelines for the Telegram SEA Edition application, based on the schema analysis and business requirements.

## Overview

The API is designed to support a real-time, multi-tenant social commerce platform optimized for Southeast Asian markets. The API follows RESTful principles with real-time WebSocket capabilities and GraphQL support for complex queries.

### Key Features
- **RESTful Design**: Standard HTTP methods and status codes
- **Real-time Support**: WebSocket events for live features
- **Multi-language**: Content localization and i18n support
- **Multi-currency**: Regional payment and pricing support
- **High Performance**: Caching, pagination, and optimization
- **Security**: Authentication, authorization, and rate limiting

---

## API Architecture

### Base URL Structure
```
Production: https://api.telegram-sea.com/v1
Staging: https://api-staging.telegram-sea.com/v1
Development: http://localhost:3000/api/v1
```

### Versioning Strategy
- **URL Versioning**: `/v1/`, `/v2/` for major versions
- **Header Versioning**: `Accept: application/vnd.telegram-sea.v1+json` for minor versions
- **Backward Compatibility**: Maintain at least 2 versions simultaneously
- **Deprecation Policy**: 6-month notice for version deprecation

### Content Types
- **Request**: `application/json` (primary), `multipart/form-data` (file uploads)
- **Response**: `application/json` (primary), `text/event-stream` (SSE)
- **Real-time**: WebSocket connections for live updates

---

## Authentication & Authorization

### Authentication Methods

#### 1. JWT Bearer Tokens
```http
POST /auth/login
Content-Type: application/json

{
  "phone": "+66812345678",
  "otp": "123456"
}

Response:
{
  "success": true,
  "data": {
    "accessToken": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9...",
    "refreshToken": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9...",
    "expiresIn": 3600,
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "John Doe",
      "avatar": "https://cdn.telegram-sea.com/avatars/123.jpg"
    }
  }
}
```

#### 2. API Keys (for service-to-service)
```http
GET /users/profile
Authorization: Bearer sk_live_abc123...
X-API-Version: v1
```

#### 3. OAuth2 (for third-party integrations)
```http
GET /oauth/authorize?client_id=abc123&response_type=code&scope=read_profile&redirect_uri=...
```

### Authorization Levels
1. **Public**: No authentication required
2. **User**: Requires valid user token
3. **Business**: Requires business account
4. **Admin**: Requires admin privileges
5. **System**: Requires service account

### Security Headers
```http
Authorization: Bearer <token>
X-API-Key: <api_key>
X-Request-ID: <uuid>
X-Client-Version: 1.2.3
X-Device-ID: <device_uuid>
X-User-Agent: TelegramSEA/1.2.3 (iOS; 14.5)
```

---

## Request/Response Format

### Standard Response Structure
```typescript
interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: Record<string, any>;
    field?: string;
    suggestion?: string;
  };
  meta?: {
    pagination?: PaginationMeta;
    timing?: { requestTime: number; totalTime: number };
    cache?: { hit: boolean; ttl: number };
  };
  timestamp: string;
  requestId: string;
}
```

### Success Response Examples
```json
// Single Resource
{
  "success": true,
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "John Doe",
    "email": "john@example.com"
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "requestId": "req_abc123"
}

// Collection with Pagination
{
  "success": true,
  "data": {
    "items": [...],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "totalPages": 8,
      "hasNextPage": true,
      "hasPreviousPage": false
    }
  },
  "meta": {
    "timing": {
      "requestTime": 45,
      "totalTime": 120
    }
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "requestId": "req_abc123"
}
```

### Error Response Examples
```json
// Validation Error
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "field": "email",
    "details": {
      "email": ["Must be a valid email address"]
    },
    "suggestion": "Please provide a valid email address"
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "requestId": "req_abc123"
}

// Business Logic Error
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_BALANCE",
    "message": "Insufficient wallet balance for transaction",
    "details": {
      "required": 1000,
      "available": 500,
      "currency": "THB"
    },
    "suggestion": "Please top up your wallet or reduce the transaction amount"
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "requestId": "req_abc123"
}
```

---

## Domain API Endpoints

## 1. User Management

### User Profile
```http
# Get current user profile
GET /users/me

# Update user profile
PATCH /users/me
{
  "name": "Updated Name",
  "bio": "Updated bio",
  "preferences": {
    "language": "th-TH",
    "currency": "THB"
  }
}

# Upload avatar
POST /users/me/avatar
Content-Type: multipart/form-data
```

### Authentication
```http
# Send OTP
POST /auth/send-otp
{
  "phone": "+66812345678",
  "type": "login"
}

# Verify OTP and login
POST /auth/verify-otp
{
  "phone": "+66812345678",
  "otp": "123456"
}

# Refresh token
POST /auth/refresh
{
  "refreshToken": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9..."
}

# Logout
POST /auth/logout
```

### KYC Management
```http
# Get KYC status
GET /users/me/kyc

# Submit KYC document
POST /users/me/kyc/documents
Content-Type: multipart/form-data

# Get KYC limits
GET /users/me/kyc/limits
```

### Friend Management
```http
# Get friends list
GET /users/me/friends?status=accepted&page=1&limit=20

# Send friend request
POST /users/{userId}/friends

# Accept friend request
PATCH /users/me/friends/{friendshipId}
{
  "status": "accepted"
}

# Remove friend
DELETE /users/me/friends/{friendshipId}
```

---

## 2. Messaging & Chat

### Dialogs
```http
# Get dialogs list
GET /dialogs?type=user,group&page=1&limit=20

# Create dialog
POST /dialogs
{
  "type": "group",
  "name": "Project Team",
  "participants": ["user1", "user2"]
}

# Update dialog
PATCH /dialogs/{dialogId}
{
  "name": "Updated Name",
  "avatar": "https://cdn.example.com/avatar.jpg"
}

# Get dialog details
GET /dialogs/{dialogId}

# Leave dialog
DELETE /dialogs/{dialogId}/members/me
```

### Messages
```http
# Get messages
GET /dialogs/{dialogId}/messages?page=1&limit=50&before=messageId

# Send message
POST /dialogs/{dialogId}/messages
{
  "type": "text",
  "content": {
    "text": "Hello world!"
  },
  "replyToId": "optional-message-id"
}

# Send media message
POST /dialogs/{dialogId}/messages
Content-Type: multipart/form-data

# Edit message
PATCH /dialogs/{dialogId}/messages/{messageId}
{
  "content": {
    "text": "Updated message"
  }
}

# Delete message
DELETE /dialogs/{dialogId}/messages/{messageId}

# React to message
POST /dialogs/{dialogId}/messages/{messageId}/reactions
{
  "emoji": "👍"
}
```

### Message Search
```http
# Search messages
GET /messages/search?q=hello&dialogId=123&type=text&limit=20

# Advanced search
POST /messages/search
{
  "query": "hello world",
  "filters": {
    "dialogId": "123",
    "senderId": "456",
    "dateRange": {
      "start": "2024-01-01T00:00:00Z",
      "end": "2024-01-31T23:59:59Z"
    },
    "hasAttachments": true
  },
  "sorting": {
    "sortBy": "createdAt",
    "sortOrder": "desc"
  }
}
```

---

## 3. E-Commerce

### Shop Management
```http
# Get user's shops
GET /shops/me

# Create shop
POST /shops
{
  "name": "My Amazing Shop",
  "description": "Best products in town",
  "category": "electronics",
  "contact": {
    "email": "shop@example.com",
    "phone": "+66812345678"
  }
}

# Update shop
PATCH /shops/{shopId}
{
  "name": "Updated Shop Name",
  "settings": {
    "currency": "THB",
    "timezone": "Asia/Bangkok"
  }
}

# Get shop analytics
GET /shops/{shopId}/analytics?period=30d
```

### Product Management
```http
# Get products
GET /products?shopId=123&status=active&category=electronics&page=1&limit=20

# Search products
GET /products/search?q=iPhone&price_min=1000&price_max=50000&currency=THB

# Get product details
GET /products/{productId}

# Create product
POST /products
{
  "title": "iPhone 15 Pro",
  "description": "Latest iPhone model",
  "price": 35000,
  "currency": "THB",
  "inventory": {
    "quantity": 10,
    "trackQuantity": true
  },
  "images": [
    {
      "url": "https://cdn.example.com/iphone1.jpg",
      "isPrimary": true
    }
  ]
}

# Update product
PATCH /products/{productId}
{
  "price": 33000,
  "inventory": {
    "quantity": 8
  }
}

# Upload product images
POST /products/{productId}/images
Content-Type: multipart/form-data
```

### Shopping Cart
```http
# Get cart
GET /cart

# Add to cart
POST /cart/items
{
  "productId": "123",
  "variantId": "456",
  "quantity": 2
}

# Update cart item
PATCH /cart/items/{itemId}
{
  "quantity": 3
}

# Remove from cart
DELETE /cart/items/{itemId}

# Apply coupon
POST /cart/coupons
{
  "code": "SAVE20"
}

# Calculate shipping
POST /cart/shipping/calculate
{
  "address": {
    "country": "TH",
    "province": "Bangkok",
    "city": "Bangkok",
    "postalCode": "10110"
  }
}
```

### Orders
```http
# Create order
POST /orders
{
  "shippingAddress": {
    "firstName": "John",
    "lastName": "Doe",
    "address1": "123 Main St",
    "city": "Bangkok",
    "province": "Bangkok",
    "country": "TH",
    "postalCode": "10110",
    "phone": "+66812345678"
  },
  "paymentMethod": {
    "type": "qr_payment",
    "provider": "promptpay"
  }
}

# Get orders
GET /orders?status=processing&page=1&limit=20

# Get order details
GET /orders/{orderId}

# Update order status (seller)
PATCH /orders/{orderId}/status
{
  "status": "shipped",
  "trackingNumber": "TH123456789"
}

# Cancel order
POST /orders/{orderId}/cancel
{
  "reason": "Changed mind"
}
```

---

## 4. Financial Services

### Wallet Management
```http
# Get wallets
GET /wallets

# Get wallet balance
GET /wallets/{walletId}/balance

# Get transaction history
GET /wallets/{walletId}/transactions?type=send,receive&page=1&limit=50

# Create transaction
POST /wallets/{walletId}/transactions
{
  "type": "send",
  "amount": 1000,
  "currency": "THB",
  "description": "Payment for lunch",
  "counterpart": {
    "phone": "+66812345678"
  }
}
```

### Payment Methods
```http
# Get payment methods
GET /payment-methods

# Add payment method
POST /payment-methods
{
  "type": "bank_account",
  "details": {
    "bankName": "Kasikornbank",
    "accountNumber": "1234567890",
    "accountName": "John Doe"
  }
}

# Set default payment method
PATCH /payment-methods/{methodId}/default
```

### QR Payments
```http
# Generate QR payment
POST /qr-payments
{
  "amount": 500,
  "currency": "THB",
  "description": "Coffee payment"
}

# Process QR payment
POST /qr-payments/{qrId}/process
{
  "paymentMethodId": "123"
}
```

---

## 5. Events & Ticketing

### Events
```http
# Get events
GET /events?category=music&location=Bangkok&date_from=2024-01-01&page=1&limit=20

# Search events
GET /events/search?q=concert&latitude=13.7563&longitude=100.5018&radius=10

# Get event details
GET /events/{eventId}

# Create event (organizer)
POST /events
{
  "title": "Bangkok Music Festival 2024",
  "description": "Amazing music festival",
  "category": "music",
  "venue": {
    "name": "Impact Arena",
    "address": "Muang Thong Thani",
    "capacity": 10000
  },
  "schedule": {
    "startDate": "2024-06-15T19:00:00+07:00",
    "endDate": "2024-06-15T23:00:00+07:00"
  }
}

# Express interest
POST /events/{eventId}/interest

# Get event attendees
GET /events/{eventId}/attendees?status=going
```

### Ticketing
```http
# Get ticket types
GET /events/{eventId}/tickets/types

# Purchase tickets
POST /events/{eventId}/tickets/purchase
{
  "tickets": [
    {
      "typeId": "vip",
      "quantity": 2
    }
  ],
  "paymentMethodId": "123"
}

# Get user tickets
GET /tickets?status=active&page=1&limit=20

# Transfer ticket
POST /tickets/{ticketId}/transfer
{
  "recipientId": "user123"
}

# Check in ticket
POST /tickets/{ticketId}/checkin
{
  "qrCode": "scanned_qr_data"
}
```

---

## Real-time Features (WebSocket)

### Connection Setup
```javascript
const ws = new WebSocket('wss://api.telegram-sea.com/ws');

// Authentication
ws.send(JSON.stringify({
  type: 'authenticate',
  token: 'bearer_token_here'
}));

// Subscribe to events
ws.send(JSON.stringify({
  type: 'subscribe',
  channels: ['dialog:123', 'user:notifications']
}));
```

### Event Types

#### Message Events
```json
{
  "type": "message.new",
  "channel": "dialog:123",
  "data": {
    "id": "msg123",
    "dialogId": "123",
    "senderId": "user456",
    "type": "text",
    "content": {
      "text": "Hello world!"
    },
    "createdAt": "2024-01-15T10:30:00Z"
  }
}

{
  "type": "message.read",
  "channel": "dialog:123",
  "data": {
    "messageId": "msg123",
    "userId": "user789",
    "readAt": "2024-01-15T10:31:00Z"
  }
}
```

#### Typing Indicators
```json
{
  "type": "typing.start",
  "channel": "dialog:123",
  "data": {
    "userId": "user456",
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

#### Payment Events
```json
{
  "type": "payment.received",
  "channel": "user:notifications",
  "data": {
    "transactionId": "tx123",
    "amount": 1000,
    "currency": "THB",
    "from": {
      "name": "John Doe",
      "avatar": "https://cdn.example.com/avatar.jpg"
    }
  }
}
```

---

## Pagination

### Cursor-based Pagination (Recommended)
```http
GET /dialogs/{dialogId}/messages?limit=50&before=msg123

Response:
{
  "success": true,
  "data": {
    "items": [...],
    "pagination": {
      "limit": 50,
      "hasNextPage": true,
      "hasPreviousPage": false,
      "nextCursor": "msg456",
      "previousCursor": null
    }
  }
}
```

### Offset-based Pagination
```http
GET /products?page=2&limit=20

Response:
{
  "success": true,
  "data": {
    "items": [...],
    "pagination": {
      "page": 2,
      "limit": 20,
      "total": 1000,
      "totalPages": 50,
      "hasNextPage": true,
      "hasPreviousPage": true
    }
  }
}
```

---

## Rate Limiting

### Rate Limit Headers
```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1390441200
X-RateLimit-Window: 3600
```

### Rate Limit Rules
| Endpoint Category | Authenticated | Rate Limit |
|------------------|---------------|------------|
| Authentication | No | 5 requests/minute |
| User Profile | Yes | 100 requests/hour |
| Messaging | Yes | 1000 requests/hour |
| File Upload | Yes | 50 requests/hour |
| Payment | Yes | 20 requests/hour |
| Public API | No | 100 requests/hour |

### Rate Limit Error Response
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMITED",
    "message": "Too many requests",
    "details": {
      "limit": 1000,
      "window": 3600,
      "retryAfter": 3600
    }
  }
}
```

---

## Localization & Internationalization

### Accept-Language Header
```http
GET /products
Accept-Language: th-TH, en-US;q=0.9
```

### Localized Responses
```json
{
  "success": true,
  "data": {
    "title": {
      "th-TH": "สินค้าที่ดีที่สุด",
      "en-US": "Best Product"
    },
    "description": {
      "th-TH": "คำอธิบายสินค้า",
      "en-US": "Product description"
    },
    "price": {
      "amount": 1000,
      "currency": "THB",
      "formatted": "฿1,000"
    }
  }
}
```

### Currency Conversion
```http
GET /products/{productId}?currency=SGD

Response:
{
  "price": {
    "amount": 38.50,
    "currency": "SGD",
    "originalAmount": 1000,
    "originalCurrency": "THB",
    "exchangeRate": 0.0385,
    "formatted": "S$38.50"
  }
}
```

---

## Error Handling

### Standard Error Codes
| Code | HTTP Status | Description |
|------|-------------|-------------|
| VALIDATION_ERROR | 400 | Invalid input data |
| UNAUTHORIZED | 401 | Authentication required |
| FORBIDDEN | 403 | Insufficient permissions |
| NOT_FOUND | 404 | Resource not found |
| METHOD_NOT_ALLOWED | 405 | HTTP method not allowed |
| CONFLICT | 409 | Resource conflict |
| UNPROCESSABLE_ENTITY | 422 | Business logic error |
| RATE_LIMITED | 429 | Too many requests |
| INTERNAL_ERROR | 500 | Server error |
| SERVICE_UNAVAILABLE | 503 | Service temporarily unavailable |

### Domain-specific Error Codes
| Code | Domain | Description |
|------|--------|-------------|
| INSUFFICIENT_BALANCE | Payment | Not enough wallet balance |
| KYC_REQUIRED | User | KYC verification needed |
| PRODUCT_OUT_OF_STOCK | Commerce | Product unavailable |
| EVENT_SOLD_OUT | Events | No tickets available |
| MESSAGE_TOO_LONG | Chat | Message exceeds limit |
| FILE_TOO_LARGE | Upload | File size exceeds limit |

---

## Security Best Practices

### Input Validation
- **Sanitize all inputs**: Prevent XSS and injection attacks
- **Validate data types**: Ensure proper type conversion
- **Check ranges**: Validate numeric ranges and string lengths
- **Whitelist approach**: Only allow expected values

### Authentication Security
- **JWT expiration**: Short-lived access tokens (1 hour)
- **Refresh token rotation**: Issue new refresh tokens
- **Device binding**: Tie tokens to specific devices
- **Suspicious activity detection**: Monitor for unusual patterns

### API Security
- **HTTPS only**: All API calls must use TLS
- **CORS configuration**: Restrict cross-origin requests
- **Request signing**: Sign sensitive API calls
- **IP whitelisting**: For admin and system APIs

### Data Privacy
- **PII encryption**: Encrypt personally identifiable information
- **Data minimization**: Only return necessary fields
- **Consent tracking**: Track user consent for data usage
- **Right to erasure**: Support data deletion requests

---

## Performance Optimization

### Caching Strategy
```http
# Cache headers
Cache-Control: public, max-age=3600
ETag: "abc123"
Last-Modified: Wed, 15 Jan 2024 10:30:00 GMT

# Conditional requests
If-None-Match: "abc123"
If-Modified-Since: Wed, 15 Jan 2024 10:30:00 GMT
```

### Response Compression
```http
Accept-Encoding: gzip, deflate, br
Content-Encoding: gzip
```

### Field Selection
```http
GET /users/me?fields=id,name,avatar
GET /products?fields=id,title,price,currency,images.url
```

### Bulk Operations
```http
POST /messages/bulk
{
  "messages": [
    {
      "dialogId": "123",
      "type": "text",
      "content": {"text": "Message 1"}
    },
    {
      "dialogId": "456",
      "type": "text",
      "content": {"text": "Message 2"}
    }
  ]
}
```

---

## Monitoring & Analytics

### Request Tracking
```http
X-Request-ID: req_abc123
X-Trace-ID: trace_xyz789
X-User-ID: user_123
X-Session-ID: session_456
```

### Performance Metrics
- **Response time**: Track P95 and P99 latencies
- **Error rates**: Monitor 4xx and 5xx responses
- **Throughput**: Requests per second by endpoint
- **Database performance**: Query execution times

### Business Metrics
- **User engagement**: API usage patterns
- **Feature adoption**: New feature usage rates
- **Revenue metrics**: Payment processing success rates
- **Geographic distribution**: Usage by country/region

---

## API Documentation

### OpenAPI Specification
```yaml
openapi: 3.0.3
info:
  title: Telegram SEA Edition API
  version: 1.0.0
  description: Social commerce platform API for Southeast Asia
  contact:
    name: API Support
    email: api-support@telegram-sea.com
    url: https://docs.telegram-sea.com

servers:
  - url: https://api.telegram-sea.com/v1
    description: Production server
  - url: https://api-staging.telegram-sea.com/v1
    description: Staging server

paths:
  /users/me:
    get:
      summary: Get current user profile
      tags: [Users]
      security:
        - bearerAuth: []
      responses:
        '200':
          description: User profile
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
```

### SDK Generation
- **TypeScript**: For web and Node.js applications
- **Swift**: For iOS native applications
- **Kotlin**: For Android native applications
- **PHP**: For server-side integrations
- **Python**: For data analysis and automation

---

This API design provides a comprehensive foundation for the Telegram SEA Edition application, supporting all major features while maintaining high performance, security, and developer experience.