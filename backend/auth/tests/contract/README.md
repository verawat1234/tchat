# Auth Service Pact Provider Verification Tests

This directory contains Pact provider verification tests for the Auth service, validating that the service implements the contracts expected by web frontend consumers.

## Overview

The provider verification tests ensure that the Auth service correctly implements the API contracts defined by the web frontend. The tests use Pact Go v2 framework to verify that the actual service responses match the expected contract specifications.

## Architecture

### Test Components

1. **Provider Verification Test** (`pact_provider_test.go`)
   - Main test file that runs Pact provider verification
   - Sets up test server with auth handlers
   - Configures provider state handlers
   - Adapts current OTP-based implementation to contract expectations

2. **Mock Services** (`mock_services.go`)
   - Mock implementations of UserService, AuthService, SessionService
   - Test data management and state setup
   - JWT token validation for testing

3. **Service Interfaces** (`service_interfaces.go`)
   - Interface definitions for contract testing
   - Request/Response structures
   - Type compatibility layer

4. **Utilities** (`utils.go`)
   - Helper functions for test execution
   - String utilities and type conversions

### Contract Adaptation

The tests handle the mismatch between the current OTP-based authentication system and the password-based contracts:

- **Login Endpoint**: Adapts password requests to OTP verification
- **Profile Endpoints**: Uses JWT middleware for authentication
- **Response Mapping**: Transforms service responses to match contract expectations

## Contract Specifications

The tests verify three main interaction scenarios:

### 1. Login with Valid Credentials
- **Provider State**: "user exists with valid credentials"
- **Request**: POST `/api/v1/auth/login`
- **Validates**: Token generation, user response format, session creation

### 2. Get User Profile
- **Provider State**: "user is authenticated with valid token"
- **Request**: GET `/api/v1/auth/profile`
- **Validates**: Profile data format, authentication middleware

### 3. Update User Profile
- **Provider State**: "user is authenticated and can update profile"
- **Request**: PUT `/api/v1/auth/profile`
- **Validates**: Profile updates, response format consistency

## Running the Tests

### Prerequisites

1. **Pact Consumer Contracts**: Ensure consumer contracts are generated and available
2. **Go Dependencies**: All required dependencies installed via `go mod tidy`
3. **Contract Files**: Consumer contract JSON files in the expected location

### Execution Methods

#### Method 1: From Backend Root (Recommended)
```bash
# From /Users/weerawat/Tchat/backend
cd /Users/weerawat/Tchat/backend
go test ./auth/tests/contract -v
```

#### Method 2: Direct Execution
```bash
# From contract test directory
cd /Users/weerawat/Tchat/backend/auth/tests/contract
go test -v .
```

### Test Configuration

The tests can be configured through environment variables:

```bash
# Set contract file path
export PACT_CONTRACT_PATH="/path/to/pact-consumer-web.json"

# Set test server port (optional, auto-detects available port)
export TEST_SERVER_PORT=8080

# Enable verbose logging
export PACT_LOG_LEVEL=DEBUG
```

### Expected Output

Successful test execution should show:
```
=== RUN   TestAuthServiceProvider
--- PASS: TestAuthServiceProvider (2.34s)
PASS
ok  	tchat.dev/auth/tests/contract	2.567s
```

## Provider State Handlers

### User Exists with Valid Credentials
- Creates test user with ID: `123e4567-e89b-12d3-a456-426614174000`
- Phone: `0123456789`, Country: Thailand
- Sets up mock OTP verification to succeed with code `123456`
- Configures password-to-OTP adaptation

### User is Authenticated with Valid Token
- Extends "user exists" setup
- Generates valid JWT token using test configuration
- Stores token in mock JWT service for middleware validation
- Creates mock session for authenticated requests

### User Can Update Profile
- Extends authentication setup
- Ensures user has active status for profile updates
- Configures mock user service for profile modification

## Test Data Management

### Mock User Service
- Stores test users in memory during test execution
- Provides user lookup by ID and phone number
- Supports profile updates and user creation

### Mock Auth Service
- Handles OTP verification with configurable success/failure
- Maps password-based requests to OTP verification
- Manages session creation and validation

### Mock Session Service
- Creates and manages test sessions
- Provides session lookup and termination
- Integrates with JWT token generation

### Mock JWT Service
- Validates test tokens against stored user data
- Provides claims extraction for middleware
- Handles token-to-user mapping

## Error Handling

### Common Issues and Solutions

1. **Contract File Not Found**
   ```
   panic: Contract file not found at any expected location
   ```
   **Solution**: Ensure contract file exists at expected path or set `PACT_CONTRACT_PATH`

2. **Port Already in Use**
   ```
   Failed to start test server: listen tcp :8080: bind: address already in use
   ```
   **Solution**: Test automatically finds available port, or set `TEST_SERVER_PORT`

3. **Module Resolution Conflicts**
   ```
   conflicting replacements for tchat.dev/auth
   ```
   **Solution**: Run tests from backend root directory using workspace

4. **Provider State Setup Failure**
   ```
   Provider state setup failed
   ```
   **Solution**: Check mock service initialization and test data setup

### Debugging Tips

1. **Enable Verbose Logging**
   ```bash
   go test -v . -args -test.v
   ```

2. **Check Test Data Setup**
   - Verify mock services are properly initialized
   - Ensure test user creation succeeds
   - Validate JWT token generation

3. **Verify Contract File**
   - Confirm contract file is valid JSON
   - Check interaction definitions match test expectations
   - Validate request/response formats

## Integration with CI/CD

### Test Automation
```bash
# Add to CI pipeline
- name: Run Provider Verification Tests
  run: |
    cd backend
    go test ./auth/tests/contract -v -timeout=10m
```

### Contract Validation
- Tests should be run after consumer contract generation
- Verify against latest consumer contracts
- Fail build if provider verification fails

### Performance Monitoring
- Track test execution time
- Monitor test stability and flakiness
- Alert on provider verification failures

## Security Considerations

### Test Data Security
- Mock services use non-production secrets
- Test users created with minimal data
- JWT tokens use test-specific signing keys

### Isolation
- Tests run in isolated environment
- No external dependencies or real database
- Mock services prevent data leakage

## Maintenance

### Contract Updates
When consumer contracts change:
1. Update provider state handlers if needed
2. Modify response mappings for new fields
3. Add new provider states for new scenarios
4. Update mock services for new functionality

### Service Evolution
When auth service changes:
1. Update mock service implementations
2. Modify request adaptation logic
3. Update test data setup procedures
4. Verify compatibility with existing contracts

### Performance Optimization
- Monitor test execution time
- Optimize mock service operations
- Reduce test setup/teardown overhead
- Parallelize independent test scenarios

## Related Documentation

- [Pact Go Documentation](https://docs.pact.io/implementation_guides/go)
- [Contract Testing Best Practices](https://docs.pact.io/getting_started/how_pact_works)
- [Auth Service API Documentation](../../README.md)
- [Backend Architecture Guide](../../../README.md)