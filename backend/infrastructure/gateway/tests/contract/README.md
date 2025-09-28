# Gateway Service Provider Verification Tests

This directory contains comprehensive Pact provider verification tests for the API Gateway service. The Gateway service acts as a reverse proxy that routes requests to backend microservices while providing authentication, service discovery, and load balancing.

## Overview

The Gateway provider verification tests validate that:

1. **Request Routing**: Gateway correctly routes requests to appropriate backend services
2. **Authentication Middleware**: Auth middleware properly validates JWT tokens
3. **Service Discovery**: Gateway can discover and route to healthy service instances
4. **Load Balancing**: Traffic is distributed across multiple service instances
5. **Error Handling**: Gateway returns appropriate errors when services are unavailable
6. **Header Forwarding**: Gateway properly forwards authentication and trace headers

## Provider States Tested

The following provider states are implemented and tested:

### Core Gateway States
- `gateway routes are configured` - Validates routing table setup
- `rate limiting is enabled` - Validates rate limiting middleware (if implemented)
- `authentication middleware is active` - Validates JWT auth middleware
- `service discovery is functional` - Validates service registry functionality

### Backend Service Health States
- `auth service is healthy` - Mock auth service returning valid responses
- `content service is healthy` - Mock content service returning valid responses

### Authentication States
- `user exists with valid credentials` - Mock valid login flow
- `user is authenticated with valid token` - Mock authenticated user context

### Content States
- `content items exist` - Mock content available for retrieval
- `user is authenticated and can create content` - Mock content creation permissions
- `content exists in mobile category` - Mock mobile-specific content

## Test Architecture

### Gateway Test Instance
The `GatewayTestInstance` provides:
- **Mock Backend Services**: HTTP servers simulating auth, content, messaging, etc.
- **Service Registry**: In-memory registry with mock service instances
- **Provider State Handlers**: Methods to configure test scenarios
- **Request Proxying**: Full reverse proxy functionality

### Mock Backend Services
Each backend service is mocked with:
- **Health Check Endpoint**: `/health` returning service status
- **Configurable Responses**: Mock responses for specific endpoints
- **Header Inspection**: Validation of forwarded headers
- **Error Simulation**: Ability to simulate service failures

## Setup Requirements

### 1. Install Pact CLI
```bash
# macOS
brew install pact-ruby-standalone

# Linux
curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash

# Windows
choco install pact
```

### 2. Install Native Dependencies
The Pact Go library requires native FFI libraries:

```bash
# macOS
brew install pact-foundation/pact-ruby/pact-ruby-standalone

# Ubuntu/Debian
wget -O - https://github.com/pact-foundation/pact-ruby-standalone/releases/download/v1.89.00/pact-1.89.00-linux-x86_64.tar.gz | tar xz
export PATH="$PATH:./pact/bin"

# Verify installation
pact-broker version
```

### 3. Start Pact Broker (Optional)
```bash
cd /Users/weerawat/Tchat/backend/infrastructure/pact-broker
docker-compose up -d
```

## Running Tests

### Basic Compilation Test
```bash
cd /Users/weerawat/Tchat/backend/infrastructure/gateway/tests/contract
GOWORK=off go test -c .
```

### Run Provider Verification
```bash
# With Pact Broker
PACT_BROKER_BASE_URL=http://localhost:9292 \
go test -v ./...

# With local contract files
PACT_FILES_DIR=/Users/weerawat/Tchat/specs/021-implement-pact-contract/contracts \
go test -v ./...
```

### Run Gateway Functionality Tests
```bash
# Test gateway routing without Pact verification
go test -v -run "TestGateway" ./...
```

## Test Scenarios

### 1. Authentication Flow Test
```bash
# Tests: POST /api/v1/auth/login → auth-service
# Validates: Request routing, response forwarding, token handling
go test -v -run "TestGatewayRouting.*auth" ./...
```

### 2. Content Service Test
```bash
# Tests: GET /api/v1/content → content-service
# Validates: Query parameter forwarding, pagination, response structure
go test -v -run "TestGatewayRouting.*content" ./...
```

### 3. Service Discovery Test
```bash
# Tests: Service registration, health checking, load balancing
# Validates: Multiple service instances, failover handling
go test -v -run "TestGatewayServiceRegistry" ./...
```

### 4. Health Check Test
```bash
# Tests: Gateway /health endpoint
# Validates: Gateway availability, service health aggregation
go test -v -run "TestGatewayHealthCheck" ./...
```

## Mock Service Configuration

### Authentication Service Mock
```go
authService.SetMockResponse("POST", "/login", MockResponse{
    StatusCode: 200,
    Headers: map[string]string{"Content-Type": "application/json"},
    Body: map[string]interface{}{
        "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
        "user": map[string]interface{}{
            "id": "123e4567-e89b-12d3-a456-426614174000",
            "phone_number": "+66812345678",
            "country_code": "TH",
        },
    },
})
```

### Content Service Mock
```go
contentService.SetMockResponse("GET", "/", MockResponse{
    StatusCode: 200,
    Body: map[string]interface{}{
        "items": []map[string]interface{}{{
            "id": "456e7890-e89b-12d3-a456-426614174000",
            "category": "general",
            "type": "text",
            "value": "Sample content",
        }},
        "pagination": map[string]interface{}{
            "current_page": 1,
            "total_items": 50,
        },
    },
})
```

## Contract Expectations

The Gateway provider verification validates against consumer contracts from:
- **Web Frontend**: `/api/v1/auth/*` requests via Gateway
- **Mobile App**: `/api/v1/content/*` requests via Gateway
- **Admin Dashboard**: `/registry/services` management requests

## Integration with CI/CD

### GitHub Actions Integration
```yaml
# .github/workflows/contract-tests.yml
- name: Run Gateway Provider Verification
  run: |
    cd backend/infrastructure/gateway/tests/contract
    PACT_BROKER_BASE_URL=${{ secrets.PACT_BROKER_URL }} \
    go test -v ./...
```

### Contract Publication
```bash
# Publish verification results to Pact Broker
go test -v \
  -env PACT_BROKER_PUBLISH_RESULTS=true \
  -env PACT_PROVIDER_VERSION=$BUILD_VERSION \
  ./...
```

## Troubleshooting

### Common Issues

**1. Native Library Not Found**
```
ld: library 'pact_ffi' not found
```
Solution: Install Pact Ruby Standalone or ensure native libraries are in PATH

**2. Pact Broker Connection Error**
```
Get "http://localhost:9292": connection refused
```
Solution: Start Pact Broker or use local contract files

**3. Service Mock Not Responding**
```
dial tcp: connection refused
```
Solution: Check mock service setup in test initialization

### Debug Mode
```bash
# Enable debug logging
PACT_LOG_LEVEL=DEBUG go test -v ./...

# Test specific provider state
go test -v -run "TestGatewayProviderVerification" ./...
```

## Performance Considerations

- **Mock Services**: Use httptest.Server for lightweight mocking
- **Parallel Tests**: Tests can run in parallel with isolated mock services
- **Resource Cleanup**: All test servers are properly closed after tests
- **Memory Usage**: Gateway test instance uses minimal memory footprint

## Security Testing

The provider verification includes security validation:
- **JWT Token Validation**: Ensures proper token forwarding and validation
- **Header Injection Prevention**: Validates header handling security
- **Service Isolation**: Ensures requests don't leak between services
- **Authentication Bypass Prevention**: Validates auth middleware enforcement

## Extension Points

### Adding New Provider States
1. Add state to `SetProviderState()` switch statement
2. Implement setup method (e.g., `setupNewState()`)
3. Configure appropriate mock service responses
4. Add corresponding consumer contract expectations

### Adding New Backend Services
1. Create mock service in `NewGatewayTestInstance()`
2. Register service in gateway registry
3. Add routing rules in `setupTestRoutes()`
4. Configure mock responses for provider states

### Custom Validation Rules
1. Extend `MockResponse` structure for complex validation
2. Add custom assertion logic in test methods
3. Implement request/response inspection middleware

This comprehensive test suite ensures the Gateway service correctly implements the API contract while maintaining reliability, security, and performance standards required for production deployment.