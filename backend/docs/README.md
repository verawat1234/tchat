# Tchat Backend API Documentation

## Overview

Tchat (Telegram SEA Edition) is a comprehensive messaging and commerce platform designed specifically for the Southeast Asian market. This backend provides a robust microservices architecture supporting real-time messaging, multi-currency payments, e-commerce functionality, and multi-channel notifications.

## Architecture

### Microservices Structure

```
backend/
├── auth/              # Authentication Service
├── messaging/         # Real-time Messaging Service
├── payment/           # Multi-currency Payment Service
├── commerce/          # E-commerce Platform Service
├── notification/      # Multi-channel Notification Service
└── shared/            # Shared Utilities and Configuration
```

### Technology Stack

- **Language**: Go 1.22+
- **Framework**: Gin HTTP framework
- **Database**: PostgreSQL with GORM
- **Real-time**: WebSocket for messaging
- **Authentication**: JWT with refresh tokens
- **Validation**: Custom validation with regional support
- **Logging**: Zap structured logging
- **Events**: Custom event bus for microservice communication

## Regional Focus

### Supported Countries
- **TH** - Thailand
- **SG** - Singapore
- **ID** - Indonesia
- **MY** - Malaysia
- **PH** - Philippines
- **VN** - Vietnam

### Supported Currencies
- **THB** - Thai Baht
- **SGD** - Singapore Dollar
- **IDR** - Indonesian Rupiah (no decimals)
- **MYR** - Malaysian Ringgit
- **PHP** - Philippine Peso
- **VND** - Vietnamese Dong (no decimals)
- **USD** - US Dollar (international transactions)

### Supported Languages
- **en** - English (default)
- **th** - Thai
- **id** - Bahasa Indonesia
- **ms** - Bahasa Malaysia
- **fil** - Filipino
- **vi** - Vietnamese

## Quick Start

### Prerequisites

- Go 1.22 or higher
- PostgreSQL 13+
- Redis (for caching and sessions)

### Environment Setup

Create a `.env` file in the root directory:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=tchat_dev
DB_USER=tchat_user
DB_PASSWORD=your_password
DB_SSL_MODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_ACCESS_EXPIRY=1h
JWT_REFRESH_EXPIRY=720h

# OTP
OTP_PROVIDER=twilio
TWILIO_ACCOUNT_SID=your_twilio_sid
TWILIO_AUTH_TOKEN=your_twilio_token
TWILIO_PHONE_NUMBER=+1234567890

# Payment Providers
STRIPE_SECRET_KEY=sk_test_your_stripe_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret
PAYPAL_CLIENT_ID=your_paypal_client_id
PAYPAL_CLIENT_SECRET=your_paypal_secret

# Notification Providers
SENDGRID_API_KEY=your_sendgrid_key
FCM_SERVER_KEY=your_fcm_key

# Server
PORT=8080
GIN_MODE=debug
```

### Installation and Setup

1. **Clone and install dependencies**:
```bash
cd backend
go mod download
```

2. **Database setup**:
```bash
# Create database
createdb tchat_dev

# Run migrations (implement as needed)
go run migrations/migrate.go
```

3. **Start the services**:
```bash
# Start all services
go run cmd/auth/main.go &
go run cmd/messaging/main.go &
go run cmd/payment/main.go &
go run cmd/commerce/main.go &
go run cmd/notification/main.go &
```

## API Documentation

### OpenAPI Specification

The complete API documentation is available in OpenAPI 3.0 format:
- **File**: `docs/api/openapi.yaml`
- **Interactive Documentation**: Use Swagger UI or similar tools

### Base URLs

- **Production**: `https://api.tchat.sea/v1`
- **Staging**: `https://staging-api.tchat.sea/v1`
- **Local**: `http://localhost:8080/api/v1`

### Authentication

All endpoints (except public ones) require JWT authentication:

```http
Authorization: Bearer <your_jwt_token>
```

### Rate Limiting

- **Authentication endpoints**: 10 requests/minute per IP
- **Messaging endpoints**: 100 requests/minute per user
- **Payment endpoints**: 20 requests/minute per user
- **Commerce endpoints**: 50 requests/minute per user
- **Notification endpoints**: 100 requests/minute per user

## Service Details

### 1. Authentication Service

**Endpoints**: `/api/v1/auth/*`

Features:
- OTP-based phone authentication
- JWT access/refresh token management
- Multi-device session support
- Regional phone number validation

**Key Endpoints**:
- `POST /auth/otp/send` - Send OTP to phone
- `POST /auth/otp/verify` - Verify OTP and login
- `POST /auth/refresh` - Refresh access token
- `POST /auth/logout` - Logout user

### 2. Messaging Service

**Endpoints**: `/api/v1/messaging/*`
**WebSocket**: `/ws/messaging`

Features:
- Real-time messaging with WebSocket
- Group and channel support
- Message reactions and replies
- Read receipts and typing indicators
- File and media sharing
- Message search and filtering

**Key Endpoints**:
- `GET /messaging/dialogs` - Get user dialogs
- `POST /messaging/dialogs` - Create new dialog
- `GET /messaging/dialogs/{id}/messages` - Get dialog messages
- `POST /messaging/dialogs/{id}/messages` - Send message

### 3. Payment Service

**Endpoints**: `/api/v1/payments/*`

Features:
- Multi-currency wallet system
- Deposit, withdrawal, transfer operations
- Payment processor integration (Stripe, PayPal)
- Transaction history and reporting
- Webhook handling for payment updates
- Regional currency support

**Key Endpoints**:
- `GET /payments/wallets` - Get user wallets
- `POST /payments/wallets/{id}/deposit` - Deposit funds
- `POST /payments/wallets/{id}/withdraw` - Withdraw funds
- `POST /payments/wallets/{id}/transfer` - Transfer funds
- `GET /payments/transactions` - Get transaction history

### 4. Commerce Service

**Endpoints**: `/api/v1/commerce/*`

Features:
- Shop and product management
- Advanced product search and filtering
- Order processing and management
- Inventory tracking
- Review and rating system
- Multi-currency pricing

**Key Endpoints**:
- `GET /commerce/shops` - Search shops
- `POST /commerce/shops` - Create shop
- `GET /commerce/products` - Search products
- `POST /commerce/products` - Create product
- `POST /commerce/orders` - Create order

### 5. Notification Service

**Endpoints**: `/api/v1/notifications/*`

Features:
- Multi-channel notifications (email, SMS, push, in-app)
- Template management with localization
- Bulk and broadcast notifications
- User preferences and subscriptions
- Delivery tracking and analytics
- Webhook integration with providers

**Key Endpoints**:
- `GET /notifications` - Get user notifications
- `POST /notifications/send` - Send notification
- `GET /notifications/templates` - Get templates
- `PUT /notifications/preferences` - Update preferences

## Development Guidelines

### Code Structure

Each service follows a clean architecture pattern:

```
service/
├── handlers/          # HTTP handlers and routing
├── services/          # Business logic
├── models/           # Data models and validation
├── repositories/     # Data access layer
└── tests/           # Service-specific tests
```

### Shared Components

```
shared/
├── config/          # Configuration management
├── middleware/      # HTTP middleware (auth, CORS, logging)
├── utils/           # Utility functions and validation
├── events/          # Event bus for microservice communication
└── database/        # Database connection and models
```

### Testing Strategy

- **Unit Tests**: Service logic testing
- **Integration Tests**: API endpoint testing
- **Contract Tests**: Microservice interface testing
- **E2E Tests**: Full workflow testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific service tests
go test ./auth/...
```

### Error Handling

All APIs return consistent error responses:

```json
{
  "error": "Human-readable error message",
  "details": "Additional error details (optional)",
  "code": "MACHINE_READABLE_ERROR_CODE"
}
```

Common error codes:
- `VALIDATION_ERROR` - Input validation failed
- `UNAUTHORIZED` - Authentication required
- `FORBIDDEN` - Access denied
- `NOT_FOUND` - Resource not found
- `RATE_LIMIT_EXCEEDED` - Rate limit exceeded
- `INTERNAL_ERROR` - Internal server error

### Logging

All services use structured logging with Zap:

```go
logger.Info("User authenticated",
    zap.String("user_id", userID),
    zap.String("country", "TH"),
    zap.Duration("duration", time.Since(start)),
)
```

Log levels:
- `DEBUG` - Development debugging
- `INFO` - General information
- `WARN` - Warning conditions
- `ERROR` - Error conditions
- `FATAL` - Fatal errors (service shutdown)

## Deployment

### Docker Support

Each service includes a Dockerfile for containerized deployment:

```bash
# Build service image
docker build -t tchat-auth:latest ./auth

# Run with docker-compose
docker-compose up -d
```

### Kubernetes

Kubernetes manifests are provided for production deployment:

```bash
# Apply all services
kubectl apply -f k8s/
```

### Environment Configuration

Different configurations for each environment:
- `config/development.yaml`
- `config/staging.yaml`
- `config/production.yaml`

## Security Considerations

### Authentication
- JWT tokens with short expiry (1 hour)
- Refresh token rotation
- Device-based session management
- Rate limiting on auth endpoints

### Data Protection
- Input validation and sanitization
- SQL injection prevention with GORM
- XSS protection with content validation
- CORS configuration for web clients

### Payment Security
- PCI DSS compliance for card data
- Webhook signature verification
- Transaction encryption
- Audit trails for all financial operations

### API Security
- Rate limiting per user/IP
- Request size limits
- HTTPS enforcement in production
- API key management for webhooks

## Monitoring and Observability

### Metrics
- Request/response metrics
- Database query performance
- WebSocket connection counts
- Payment transaction metrics
- Notification delivery rates

### Health Checks
- `/health` endpoint for service health
- Database connectivity checks
- External service dependency checks
- Resource utilization monitoring

### Logging
- Structured JSON logging
- Request/response logging
- Error tracking and alerting
- Performance monitoring

## Contributing

### Development Workflow

1. **Create feature branch**: `git checkout -b feature/new-feature`
2. **Implement changes**: Follow coding standards
3. **Write tests**: Maintain test coverage
4. **Run tests**: `go test ./...`
5. **Submit PR**: Include description and test results

### Coding Standards

- Go formatting with `gofmt`
- Linting with `golangci-lint`
- Documentation with godoc
- Error handling best practices
- Consistent naming conventions

### Pull Request Requirements

- [ ] Tests pass
- [ ] Code coverage maintained
- [ ] Documentation updated
- [ ] API changes documented
- [ ] Breaking changes noted

## API Examples

### Authentication Flow

```bash
# 1. Send OTP
curl -X POST http://localhost:8080/api/v1/auth/otp/send \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "812345678", "country_code": "TH"}'

# 2. Verify OTP
curl -X POST http://localhost:8080/api/v1/auth/otp/verify \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "812345678",
    "country_code": "TH",
    "otp_code": "123456"
  }'
```

### Send Message

```bash
curl -X POST http://localhost:8080/api/v1/messaging/dialogs/{dialog_id}/messages \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hello!",
    "message_type": "text"
  }'
```

### Create Product

```bash
curl -X POST http://localhost:8080/api/v1/commerce/products \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro",
    "price": "45000.00",
    "currency": "THB",
    "category": "smartphones"
  }'
```

## Support

For API support and questions:
- **Documentation**: This README and OpenAPI spec
- **Issues**: GitHub Issues
- **Email**: api-support@tchat.sea
- **Community**: Developer Discord server

## License

This project is licensed under the MIT License - see the LICENSE file for details.