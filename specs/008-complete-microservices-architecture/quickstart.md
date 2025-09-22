# Quickstart: Tchat Backend Platform

**Purpose**: Validate core user scenarios through end-to-end testing
**Prerequisites**: Docker, Docker Compose, Go 1.22+, curl or API testing tool

## Quick Setup

### 1. Environment Setup (5 minutes)
```bash
# Clone repository
git clone <repo-url>
cd tchat

# Start infrastructure services
docker-compose -f docker-compose.dev.yml up -d

# Wait for services to be ready (check health)
./scripts/wait-for-services.sh

# Run database migrations
make migrate-up

# Start backend services
make dev
```

### 2. Verify Service Health (1 minute)
```bash
# Check all services are running
curl http://localhost:8081/health  # Auth Service
curl http://localhost:8082/health  # Messaging Service
curl http://localhost:8083/health  # Payment Service
curl http://localhost:8084/health  # Notification Service
curl http://localhost:8080/health  # API Gateway

# Expected response for each:
# {"status":"healthy","timestamp":"2025-09-22T10:00:00Z"}
```

## Core User Scenarios

### Scenario 1: User Registration & Authentication (Thailand)
**Goal**: Verify Southeast Asian user can register and authenticate

```bash
# 1. Send OTP to Thai phone number
curl -X POST http://localhost:8080/api/v1/auth/otp/send \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "812345678",
    "country_code": "TH",
    "locale": "th"
  }'

# Expected: {"message":"OTP sent successfully","expires_in":300,"retry_after":60}

# 2. Verify OTP (use test code 123456 in dev environment)
curl -X POST http://localhost:8080/api/v1/auth/otp/verify \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "812345678",
    "country_code": "TH",
    "otp_code": "123456",
    "device_info": {
      "device_id": "test_device_thai",
      "platform": "web",
      "app_version": "1.0.0"
    }
  }'

# Expected: Authentication response with access_token and user profile
# Save the access_token for subsequent requests
export TOKEN="<access_token_from_response>"
```

### Scenario 2: Real-time Messaging
**Goal**: Verify users can send and receive messages

```bash
# 1. Create a direct conversation (requires 2 authenticated users)
# Register second user first (Singapore)
curl -X POST http://localhost:8080/api/v1/auth/otp/send \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "81234567",
    "country_code": "SG",
    "locale": "en"
  }'

curl -X POST http://localhost:8080/api/v1/auth/otp/verify \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "81234567",
    "country_code": "SG",
    "otp_code": "123456"
  }'

export TOKEN2="<second_user_token>"

# 2. Get user dialogs (should show auto-created direct dialog)
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/messaging/dialogs

# Expected: List of dialogs, note dialog_id for messaging

# 3. Send a message
export DIALOG_ID="<dialog_id_from_response>"
curl -X POST http://localhost:8080/api/v1/messaging/dialogs/$DIALOG_ID/messages \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "สวัสดีครับ! (Hello in Thai)",
    "message_type": "text"
  }'

# Expected: Message object with delivery confirmation

# 4. Retrieve messages as second user
curl -H "Authorization: Bearer $TOKEN2" \
  http://localhost:8080/api/v1/messaging/dialogs/$DIALOG_ID/messages

# Expected: List containing the sent message
```

### Scenario 3: Multi-currency Payment Processing
**Goal**: Verify payment processing in local Southeast Asian currencies

```bash
# 1. Get user wallets
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/payments/wallets

# Expected: List of wallets for different currencies

# 2. Deposit funds to THB wallet (Thailand user)
export WALLET_ID="<thb_wallet_id>"
curl -X POST http://localhost:8080/api/v1/payments/wallets/$WALLET_ID/deposit \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": "100.00",
    "payment_method": "test_card",
    "description": "Test deposit"
  }'

# Expected: Transaction object with pending status
# Note: In test environment, transactions auto-complete

# 3. Check transaction history
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/payments/transactions

# Expected: List showing the deposit transaction
```

### Scenario 4: E-commerce Shop Management
**Goal**: Verify business can create shop and list products

```bash
# 1. Create a business/shop
curl -X POST http://localhost:8080/api/v1/commerce/shops \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bangkok Electronics Store",
    "description": "Premium electronics in Bangkok",
    "category": "electronics",
    "address": {
      "street": "123 Sukhumvit Road",
      "city": "Bangkok",
      "country": "TH",
      "postal_code": "10110"
    },
    "contact": {
      "phone": "+66812345678",
      "email": "shop@bangkok-electronics.com"
    }
  }'

# Expected: Shop object with verification status

# 2. Add a product to the shop
export SHOP_ID="<shop_id_from_response>"
curl -X POST http://localhost:8080/api/v1/commerce/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shop_id": "'$SHOP_ID'",
    "name": "iPhone 15 Pro Max",
    "description": "Latest iPhone with advanced features",
    "category": "smartphones",
    "price": "45000.00",
    "currency": "THB",
    "stock_quantity": 10
  }'

# Expected: Product object with Thai Baht pricing

# 3. Search for products
curl "http://localhost:8080/api/v1/commerce/products?search=iPhone&currency=THB"

# Expected: Product listings with THB prices
```

### Scenario 5: Multi-channel Notifications
**Goal**: Verify notification system works across channels

```bash
# 1. Send a test notification
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "in_app",
    "channel": "test_notifications",
    "subject": "Welcome to Tchat!",
    "content": "ยินดีต้อนรับสู่ Tchat! (Welcome to Tchat in Thai)",
    "priority": "medium"
  }'

# Expected: Notification object with delivery status

# 2. Get user notifications
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/notifications

# Expected: List including the sent notification
```

## WebSocket Real-time Testing

### Test Real-time Messaging
```bash
# Install wscat for WebSocket testing
npm install -g wscat

# Connect to messaging WebSocket
wscat -c ws://localhost:8080/ws/messaging \
  -H "Authorization: Bearer $TOKEN"

# Once connected, send join dialog message:
{"type":"join_dialog","data":{"dialog_id":"<dialog_id>"}}

# Send a real-time message:
{"type":"send_message","data":{"dialog_id":"<dialog_id>","content":"Real-time message!","message_type":"text"}}

# Expected: Immediate message delivery to other connected users
```

## Performance Validation

### Load Testing (Optional)
```bash
# Test authentication load
go run tests/performance/auth_load_test.go -users=100 -duration=30s

# Expected: >400 RPS, <200ms P95 latency

# Test messaging load
go run tests/performance/messaging_load_test.go -connections=500 -duration=30s

# Expected: All connections established, <50ms message delivery
```

## Integration Testing

### Contract Test Validation
```bash
# Run contract tests to verify API compliance
go test ./tests/contract/...

# Expected: All contract tests pass, API responses match OpenAPI specs

# Run integration tests
go test ./tests/integration/...

# Expected: End-to-end workflows complete successfully
```

## Regional Compliance Verification

### Test Country-Specific Features
```bash
# Test Indonesian user with IDR currency (no decimals)
curl -X POST http://localhost:8080/api/v1/auth/otp/send \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "812345678",
    "country_code": "ID",
    "locale": "id"
  }'

# Verify OTP sent in Indonesian language
# Test IDR wallet operations (amounts should be whole numbers)

# Test Vietnamese user
curl -X POST http://localhost:8080/api/v1/auth/otp/send \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "901234567",
    "country_code": "VN",
    "locale": "vi"
  }'

# Verify VND currency handling (no decimals)
```

## Success Criteria

### ✅ Core Functionality
- [ ] Users from all 6 SEA countries can register and authenticate
- [ ] Real-time messaging works with <2s delivery time
- [ ] Multi-currency payments process correctly (7 currencies)
- [ ] E-commerce shops can be created and products listed
- [ ] Notifications deliver across multiple channels
- [ ] WebSocket connections handle concurrent users

### ✅ Performance Targets
- [ ] Authentication: >400 RPS sustained
- [ ] Messaging: >500 concurrent WebSocket connections
- [ ] API latency: <200ms P95 response time
- [ ] Database queries: <100ms P95 query time
- [ ] Message delivery: <2s end-to-end latency

### ✅ Regional Compliance
- [ ] Local language support for OTP and notifications
- [ ] Currency handling matches regional requirements (decimals)
- [ ] Phone number validation for each country
- [ ] Payment method integration per region
- [ ] Data residency options available

### ✅ Enterprise Requirements
- [ ] 99.9% uptime during testing period
- [ ] Zero-downtime service updates
- [ ] Comprehensive monitoring and alerting
- [ ] Backup and recovery procedures tested
- [ ] Security audit compliance

## Troubleshooting

### Common Issues
1. **Services not starting**: Check Docker containers and port conflicts
2. **Database connection fails**: Verify PostgreSQL is running and migrations applied
3. **OTP not sending**: Check Twilio configuration in development environment
4. **WebSocket connection fails**: Verify JWT token format and network access
5. **Currency errors**: Ensure decimal handling matches currency requirements

### Debug Commands
```bash
# Check service logs
docker-compose logs -f auth-service
docker-compose logs -f messaging-service

# Verify database connections
go run tools/db-test/main.go

# Test external service connectivity
curl http://localhost:8081/debug/health
```

### Performance Issues
```bash
# Monitor resource usage
docker stats

# Check database performance
go run tools/db-bench/main.go

# Profile WebSocket connections
go run tools/ws-bench/main.go -connections=100
```

This quickstart validates that the Tchat backend platform successfully serves Southeast Asian users with region-specific features, multi-currency support, and enterprise-grade performance and reliability.