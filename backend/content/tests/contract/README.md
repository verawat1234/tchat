# Content Service Provider Verification Tests

This directory contains Pact provider verification tests for the Content service, verifying contracts against mobile consumer applications.

## Overview

The provider verification tests ensure that the Content service meets the contract expectations from mobile apps (iOS/Android). These tests verify:

1. **Content pagination endpoints** - GET `/api/v1/content` with query parameters
2. **Content creation endpoints** - POST `/api/v1/content` for user-generated content
3. **Category-based content retrieval** - GET `/api/v1/content/category/{category}` for mobile-optimized content

## Test Architecture

### Provider States

The tests implement three provider states:

#### 1. "content items exist"
- **Purpose**: Tests pagination functionality
- **Setup**: Creates 15 test content items in the "general" category
- **Verifies**: Pagination response structure, content format, and filtering

#### 2. "user is authenticated and can create content"
- **Purpose**: Tests content creation with authentication
- **Setup**: Sets up authenticated user context and "user-generated" category
- **Verifies**: Content creation, authentication handling, and response format

#### 3. "content exists in mobile category"
- **Purpose**: Tests mobile-optimized content retrieval
- **Setup**: Creates mobile-specific content with image metadata
- **Verifies**: Mobile-optimized response format with retina URLs and dimensions

### Test Data Management

- **Setup**: Each provider state creates specific test data in PostgreSQL
- **Isolation**: Tests run with database cleanup between states
- **Cleanup**: Automatic teardown removes all test data after execution

## Prerequisites

### 1. Database Setup

Create a test database for provider verification:

```bash
# Create test database
createdb tchat_content_test

# Or use Docker
docker run --name tchat-test-db -e POSTGRES_PASSWORD=password -e POSTGRES_DB=tchat_content_test -p 5433:5432 -d postgres:15
```

### 2. Environment Variables

Set the test database URL:

```bash
export TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/tchat_content_test?sslmode=disable"

# Or for Docker setup
export TEST_DATABASE_URL="postgres://postgres:password@localhost:5433/tchat_content_test?sslmode=disable"
```

### 3. Dependencies

Install Go dependencies:

```bash
cd backend/content/tests/contract
go mod download
```

## Running the Tests

### Full Provider Verification

Run all provider verification tests:

```bash
cd backend/content/tests/contract
go test -v -run TestContentServiceProviderVerification
```

### Individual Provider State Tests

Test provider state setup in isolation:

```bash
# Test all provider states
go test -v -run TestProviderStateSetup

# Test specific provider state
go test -v -run TestProviderStateSetup/ContentItemsExist
go test -v -run TestProviderStateSetup/AuthenticatedUser
go test -v -run TestProviderStateSetup/MobileCategoryContent
```

### Health Check Verification

Verify the test server setup:

```bash
go test -v -run TestContentServiceHealthCheck
```

### Verbose Output

For detailed verification output:

```bash
go test -v -run TestContentServiceProviderVerification -args -verbose
```

## Contract File Location

The tests reference the mobile consumer contract at:
```
specs/021-implement-pact-contract/contracts/pact-consumer-mobile.json
```

Ensure this file exists with the expected mobile app contract definitions.

## Expected Test Output

Successful verification output should show:

```
=== RUN   TestContentServiceProviderVerification
--- PASS: TestContentServiceProviderVerification (2.34s)
    pact_provider_test.go:XXX: Using contract file: ../../../specs/021-implement-pact-contract/contracts/pact-consumer-mobile.json
    pact_provider_test.go:XXX: Created 15 test content items
    pact_provider_test.go:XXX: Set up authenticated user state with user ID: 550e8400-e29b-41d4-a716-446655440000
    pact_provider_test.go:XXX: Created mobile-optimized content with ID: [uuid]
PASS
```

## Verification Details

### Authentication Handling
- Mock authentication middleware recognizes `Bearer mobile-token`
- Sets authenticated user ID: `550e8400-e29b-41d4-a716-446655440000`
- Enables content creation permissions

### Content Type Handling
The tests handle multiple content types:
- **Text content**: Standard text with metadata
- **Image content**: Mobile-optimized images with dimensions and retina URLs
- **User-generated content**: Content created by authenticated users

### Response Format Verification
Tests verify:
- **Status codes**: 200 (GET), 201 (POST)
- **Content-Type headers**: `application/json`
- **Response structure**: Items, pagination, metadata
- **Field presence**: Required fields per contract
- **Data types**: UUID format validation, array structures

## Troubleshooting

### Database Connection Issues
```bash
# Check database connectivity
psql -h localhost -p 5432 -U postgres -d tchat_content_test -c "SELECT 1;"

# For Docker setup
psql -h localhost -p 5433 -U postgres -d tchat_content_test -c "SELECT 1;"
```

### Contract File Not Found
```bash
# Verify contract file exists
ls -la ../../../specs/021-implement-pact-contract/contracts/pact-consumer-mobile.json

# Check current working directory
pwd
```

### Import Path Issues
If Go cannot find content service packages, verify the module path in `go.mod`:
```bash
# Check module path
head -n 5 ../../../../go.mod
```

### Port Conflicts
Tests use dynamic port allocation. If conflicts occur:
```bash
# Check for processes using test ports
netstat -tulpn | grep :808*
```

## Integration with CI/CD

### GitHub Actions Example
```yaml
- name: Run Content Service Provider Verification
  run: |
    cd backend/content/tests/contract
    go test -v -run TestContentServiceProviderVerification
  env:
    TEST_DATABASE_URL: postgres://postgres:password@localhost:5432/tchat_content_test?sslmode=disable
```

### Pact Broker Integration
To publish verification results to a Pact Broker, update the verification config:

```go
PublishVerificationResults: true,
BrokerURL: "https://your-pact-broker.com",
BrokerToken: os.Getenv("PACT_BROKER_TOKEN"),
```

## Performance Considerations

- **Test Duration**: ~2-5 seconds for full verification
- **Database Operations**: Bulk inserts for test data creation
- **Memory Usage**: Minimal - test data is cleaned up after each state
- **Network**: Local HTTP server for verification

## Security Notes

- **Test Database**: Uses isolated test database, never production
- **Authentication**: Mock authentication for testing only
- **Network**: Local server binding for security
- **Data Cleanup**: Automatic cleanup prevents data leakage

## Contract Evolution

When mobile consumer contracts change:
1. Update the contract file at `specs/021-implement-pact-contract/contracts/pact-consumer-mobile.json`
2. Run verification tests to identify breaking changes
3. Update provider state handlers if needed
4. Ensure backward compatibility for existing mobile apps

## Related Documentation

- [Pact Go v2 Documentation](https://docs.pact.io/implementation_guides/go/readme)
- [Content Service API Specification](../../README.md)
- [Contract Testing Strategy](../../../../specs/021-implement-pact-contract/README.md)