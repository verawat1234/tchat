# Quickstart Guide: Backend Services Architecture

**Date**: 2025-09-22
**Feature**: Backend Services Architecture for Telegram SEA Edition

## Overview

This quickstart guide provides step-by-step instructions to validate the backend services architecture through contract testing and integration scenarios. All tests should initially fail (no implementation yet) according to TDD principles.

## Prerequisites

### Development Environment
```bash
# Go 1.21+ installation
go version  # Should show 1.21 or higher

# Required tools
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker for test containers
docker --version
docker-compose --version
```

### Infrastructure Dependencies
```bash
# Start local development stack
docker-compose -f docker-compose.dev.yml up -d

# Services started:
# - PostgreSQL (auth, commerce data)
# - ScyllaDB (message timelines)
# - Redis (caching, sessions)
# - Kafka (event streaming)
# - Redis (real-time messaging)
```

## Service Setup

### 1. Project Structure Creation
```bash
# Create backend service structure
mkdir -p backend/{auth,messaging,payment,commerce,notification}/
mkdir -p backend/shared/{events,models,middleware}/
mkdir -p backend/tests/{contract,integration,unit}/

# Initialize Go modules for each service
cd backend/auth && go mod init auth-service
cd ../messaging && go mod init messaging-service
cd ../payment && go mod init payment-service
cd ../commerce && go mod init commerce-service
cd ../notification && go mod init notification-service
cd ../shared && go mod init shared
```

### 2. Generate API Contracts
```bash
# Generate Go code from OpenAPI specs
cd backend/auth
oapi-codegen -package auth ../../specs/004-create-backend-spec/contracts/auth-service.yaml > api.go

cd ../messaging
oapi-codegen -package messaging ../../specs/004-create-backend-spec/contracts/messaging-service.yaml > api.go

cd ../payment
oapi-codegen -package payment ../../specs/004-create-backend-spec/contracts/payment-service.yaml > api.go
```

## Contract Test Scenarios

### 3. Auth Service Contract Tests
```bash
# Create auth service contract tests
cd backend/tests/contract
```

Create `auth_contract_test.go`:
```go
package contract_test

import (
    "testing"
    "net/http"
    "bytes"
    "encoding/json"
    "github.com/stretchr/testify/assert"
)

func TestAuthService_SendOTP_Contract(t *testing.T) {
    // Test contract for /auth/otp/send endpoint

    // Valid request payload
    payload := map[string]interface{}{
        "identifier": "+66812345678",
        "type": "phone",
        "country": "TH",
    }

    body, _ := json.Marshal(payload)
    resp, err := http.Post("http://localhost:8001/auth/otp/send",
        "application/json", bytes.NewBuffer(body))

    // Should fail initially (no implementation)
    assert.Error(t, err) // Connection refused expected

    // Expected response structure when implemented:
    // Status: 200
    // Body: {"success": true, "session_id": "uuid", "expires_at": "timestamp"}
}

func TestAuthService_VerifyOTP_Contract(t *testing.T) {
    // Test contract for /auth/otp/verify endpoint

    payload := map[string]interface{}{
        "session_id": "550e8400-e29b-41d4-a716-446655440000",
        "code": "123456",
    }

    body, _ := json.Marshal(payload)
    resp, err := http.Post("http://localhost:8001/auth/otp/verify",
        "application/json", bytes.NewBuffer(body))

    // Should fail initially
    assert.Error(t, err)

    // Expected response when implemented:
    // Status: 200
    // Body: {"access_token": "jwt", "refresh_token": "token", "user": {...}}
}

func TestAuthService_GetUserProfile_Contract(t *testing.T) {
    // Test contract for /users/profile endpoint

    req, _ := http.NewRequest("GET", "http://localhost:8001/users/profile", nil)
    req.Header.Set("Authorization", "Bearer valid-jwt-token")

    client := &http.Client{}
    resp, err := client.Do(req)

    // Should fail initially
    assert.Error(t, err)

    // Expected response when implemented:
    // Status: 200
    // Body: User object matching schema
}
```

### 4. Messaging Service Contract Tests
```bash
# Create messaging service contract tests
```

Create `messaging_contract_test.go`:
```go
package contract_test

import (
    "testing"
    "net/http"
    "github.com/stretchr/testify/assert"
    "github.com/gorilla/websocket"
)

func TestMessagingService_GetDialogs_Contract(t *testing.T) {
    // Test contract for /dialogs endpoint

    req, _ := http.NewRequest("GET", "http://localhost:8002/dialogs", nil)
    req.Header.Set("Authorization", "Bearer valid-jwt-token")

    client := &http.Client{}
    resp, err := client.Do(req)

    // Should fail initially
    assert.Error(t, err)

    // Expected response when implemented:
    // Status: 200
    // Body: {"dialogs": [...], "total": 0, "has_more": false}
}

func TestMessagingService_WebSocket_Contract(t *testing.T) {
    // Test WebSocket connection contract

    url := "ws://localhost:8002/websocket?token=valid-jwt-token"
    conn, resp, err := websocket.DefaultDialer.Dial(url, nil)

    // Should fail initially
    assert.Error(t, err)

    // Expected when implemented:
    // - WebSocket connection established
    // - Real-time message delivery
    // - Presence updates
}

func TestMessagingService_SendMessage_Contract(t *testing.T) {
    // Test contract for POST /dialogs/{id}/messages

    dialogID := "550e8400-e29b-41d4-a716-446655440000"
    payload := map[string]interface{}{
        "type": "text",
        "content": map[string]string{"text": "Hello World"},
    }

    body, _ := json.Marshal(payload)
    url := fmt.Sprintf("http://localhost:8002/dialogs/%s/messages", dialogID)

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
    req.Header.Set("Authorization", "Bearer valid-jwt-token")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)

    // Should fail initially
    assert.Error(t, err)

    // Expected response when implemented:
    // Status: 201
    // Body: Message object with ID, timestamps, etc.
}
```

### 5. Payment Service Contract Tests
Create `payment_contract_test.go`:
```go
package contract_test

import (
    "testing"
    "net/http"
    "github.com/stretchr/testify/assert"
)

func TestPaymentService_GetWallets_Contract(t *testing.T) {
    // Test contract for /wallets endpoint

    req, _ := http.NewRequest("GET", "http://localhost:8003/wallets", nil)
    req.Header.Set("Authorization", "Bearer valid-jwt-token")

    client := &http.Client{}
    resp, err := client.Do(req)

    // Should fail initially
    assert.Error(t, err)

    // Expected response when implemented:
    // Status: 200
    // Body: Array of wallet objects
}

func TestPaymentService_SendMoney_Contract(t *testing.T) {
    // Test contract for /transactions/send endpoint

    payload := map[string]interface{}{
        "wallet_id": "550e8400-e29b-41d4-a716-446655440000",
        "recipient_id": "660e8400-e29b-41d4-a716-446655440000",
        "amount": 10000, // 100.00 in cents
        "currency": "THB",
        "description": "Test payment",
    }

    body, _ := json.Marshal(payload)
    resp, err := http.Post("http://localhost:8003/transactions/send",
        "application/json", bytes.NewBuffer(body))

    // Should fail initially
    assert.Error(t, err)

    // Expected response when implemented:
    // Status: 201
    // Body: Transaction object with status, fees, etc.
}
```

## Integration Test Scenarios

### 6. User Registration & Authentication Flow
Create `integration/auth_flow_test.go`:
```go
package integration_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCompleteAuthFlow_Integration(t *testing.T) {
    // Test complete user authentication flow

    t.Run("1. Send OTP to new user", func(t *testing.T) {
        // POST /auth/otp/send with phone number
        // Should create user if not exists
        // Should send SMS via external provider
        // Should return session_id

        t.Skip("Implementation required")
    })

    t.Run("2. Verify OTP and get tokens", func(t *testing.T) {
        // POST /auth/otp/verify with session_id and code
        // Should validate OTP
        // Should return JWT tokens
        // Should create user session

        t.Skip("Implementation required")
    })

    t.Run("3. Access protected endpoint", func(t *testing.T) {
        // GET /users/profile with JWT token
        // Should validate token
        // Should return user profile

        t.Skip("Implementation required")
    })

    t.Run("4. Refresh expired token", func(t *testing.T) {
        // POST /auth/token/refresh with refresh token
        // Should validate refresh token
        // Should return new access token

        t.Skip("Implementation required")
    })
}
```

### 7. Real-time Messaging Flow
Create `integration/messaging_flow_test.go`:
```go
package integration_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestMessagingFlow_Integration(t *testing.T) {
    // Test complete messaging workflow

    t.Run("1. Create dialog between users", func(t *testing.T) {
        // POST /dialogs with participants
        // Should create dialog in database
        // Should notify participants via WebSocket

        t.Skip("Implementation required")
    })

    t.Run("2. Establish WebSocket connections", func(t *testing.T) {
        // Connect both users via WebSocket
        // Should authenticate via JWT
        // Should establish presence

        t.Skip("Implementation required")
    })

    t.Run("3. Send and receive messages", func(t *testing.T) {
        // POST message to dialog
        // Should store in ScyllaDB
        // Should deliver via WebSocket <200ms
        // Should update unread counts

        t.Skip("Implementation required")
    })

    t.Run("4. Mark messages as read", func(t *testing.T) {
        // POST /messages/{id}/read
        // Should update read status
        // Should notify sender via WebSocket

        t.Skip("Implementation required")
    })
}
```

### 8. Payment Processing Flow
Create `integration/payment_flow_test.go`:
```go
package integration_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestPaymentFlow_Integration(t *testing.T) {
    // Test complete payment workflow

    t.Run("1. Create user wallets", func(t *testing.T) {
        // POST /wallets for each currency
        // Should create wallet with zero balance
        // Should set appropriate limits based on KYC tier

        t.Skip("Implementation required")
    })

    t.Run("2. Top up wallet via PromptPay", func(t *testing.T) {
        // POST /transactions/topup with PromptPay method
        // Should generate QR code
        // Should integrate with PromptPay API
        // Should update balance on confirmation

        t.Skip("Implementation required")
    })

    t.Run("3. Send money between users", func(t *testing.T) {
        // POST /transactions/send
        // Should validate sufficient balance
        // Should use distributed transaction (saga)
        // Should debit sender, credit receiver
        // Should maintain audit trail

        t.Skip("Implementation required")
    })

    t.Run("4. Handle payment failures", func(t *testing.T) {
        // Test insufficient funds
        // Test daily limit exceeded
        // Test external service failures
        // Should implement proper rollback

        t.Skip("Implementation required")
    })
}
```

## Performance Validation

### 9. Load Testing Scenarios
Create `performance/load_test.go`:
```go
package performance_test

import (
    "testing"
    "time"
    "sync"
)

func TestMessageDeliveryLatency_Performance(t *testing.T) {
    // Test message delivery under load

    t.Run("Validate <200ms delivery requirement", func(t *testing.T) {
        // Send 1000 messages concurrently
        // Measure delivery latency via WebSocket
        // Assert 95th percentile < 200ms

        t.Skip("Implementation required")
    })

    t.Run("Test 100K concurrent WebSocket connections", func(t *testing.T) {
        // Establish 100K WebSocket connections
        // Send messages through random connections
        // Measure memory usage and response times

        t.Skip("Implementation required")
    })
}

func TestPaymentThroughput_Performance(t *testing.T) {
    // Test payment processing throughput

    t.Run("Process 1000 transactions per minute", func(t *testing.T) {
        // Send payment requests at 1000 TPS
        // Validate all transactions processed
        // Ensure no double-spending

        t.Skip("Implementation required")
    })
}
```

## Running Tests

### Execute Contract Tests
```bash
# Run all contract tests (should fail initially)
cd backend/tests
go test ./contract/... -v

# Expected output:
# FAIL: All tests should fail with "connection refused" or similar
# This confirms no implementation exists yet (TDD approach)
```

### Execute Integration Tests
```bash
# Run integration tests (should skip initially)
go test ./integration/... -v

# Expected output:
# SKIP: All tests should skip with "Implementation required"
# Tests will be enabled as implementation progresses
```

### Execute Performance Tests
```bash
# Run performance tests (when implementation ready)
go test ./performance/... -v -timeout=10m

# Expected output:
# Performance baselines to be established during implementation
```

## Validation Checklist

### Contract Compliance
- [ ] All OpenAPI contracts generate valid Go code
- [ ] Request/response schemas match data models
- [ ] Authentication patterns consistent across services
- [ ] Error responses follow RFC 7807 standard

### Integration Readiness
- [ ] Database schemas align with data models
- [ ] Event sourcing patterns defined for cross-service communication
- [ ] External service integration points identified
- [ ] Monitoring and observability hooks planned

### Performance Targets
- [ ] Message delivery latency baseline <200ms
- [ ] API response time baseline <100ms (95th percentile)
- [ ] Payment processing throughput 1000+ TPS
- [ ] WebSocket connection capacity 100K+ concurrent

### Security & Compliance
- [ ] JWT authentication flow validated
- [ ] PII data encryption strategies defined
- [ ] PDPA compliance workflows planned
- [ ] PCI DSS requirements mapped to implementation

## Next Steps

1. **Phase 2**: Generate detailed implementation tasks from this quickstart
2. **Phase 3**: Implement services following TDD approach
3. **Phase 4**: Enable and validate integration tests
4. **Phase 5**: Performance testing and optimization

All tests should initially fail, confirming proper TDD approach before implementation begins.