# Notification Service Provider Verification Tests

## Overview

This directory contains comprehensive Pact provider verification tests for the Notification service, implementing T016 from the contract testing specification. These tests verify that the notification service properly implements the contracts expected by consumer applications (web, iOS, Android).

## Test Coverage

The provider verification tests cover:

- **Push Notifications**: Send notifications to mobile and web clients
- **Notification Preferences**: User preference management for different notification types
- **Template Management**: Notification template creation and localization
- **Delivery APIs**: Notification delivery, status tracking, and read receipts
- **User Management**: Per-user notification history and unread counts
- **Device Management**: Device token registration for push notifications

## Provider States

The following provider states are implemented to support consumer contract verification:

### 1. `user has notification preferences`
- **Purpose**: Sets up user notification preferences for testing preference APIs
- **Test Data**: Creates mock preferences with push/email/SMS settings, category preferences, and quiet hours
- **Parameters**: `user_id` (optional, defaults to test user)

### 2. `notifications exist for user`
- **Purpose**: Creates test notifications for a specific user
- **Test Data**: Generates sample notifications with different priorities and statuses
- **Parameters**: `user_id` (optional, defaults to test user)

### 3. `user can receive push notifications`
- **Purpose**: Sets up device tokens and push preferences for push notification testing
- **Test Data**: Mock iOS and Android device tokens, enabled push preferences
- **Parameters**: `user_id` (optional, defaults to test user)

### 4. `notification templates are available`
- **Purpose**: Creates notification templates for template-based notification testing
- **Test Data**: Welcome, payment success, and security alert templates with localization
- **Parameters**: None

## API Endpoints Tested

### Core Notification Endpoints
- `POST /api/v1/notifications/` - Send single notification
- `POST /api/v1/notifications/bulk` - Send bulk notifications
- `PUT /api/v1/notifications/:id/read` - Mark notification as read

### Extended Endpoints (added for comprehensive testing)
- `GET /api/v1/notifications/preferences/:user_id` - Get user notification preferences
- `PUT /api/v1/notifications/preferences/:user_id` - Update user notification preferences
- `GET /api/v1/notifications/history/:user_id` - Get notification history with pagination
- `POST /api/v1/notifications/device-tokens` - Register device token for push notifications
- `GET /api/v1/notifications/unread-count/:user_id` - Get unread notification count

### Template Management Endpoints
- `POST /api/v1/admin/notifications/templates` - Create notification template
- `GET /api/v1/admin/notifications/templates/:type` - Get template by type

## Test Structure

### Mock Implementations
The test file includes comprehensive mock implementations:

- **MockNotificationRepository**: In-memory notification storage for testing
- **MockTemplateRepository**: Template management for testing
- **MockCacheService**: Cache service for preferences and device tokens
- **MockEventService**: Event publishing for notifications

### Test Fixtures
Pre-defined test data includes:
- Test user IDs for different scenarios
- Sample notification templates with English and Thai localization
- Mock notification data with various statuses and priorities

## Running the Tests

### Prerequisites
```bash
# Install Pact CLI (if not already installed)
go install github.com/pact-foundation/pact-go/v2@latest

# Ensure Pact dependencies are available
cd /Users/weerawat/Tchat/backend/tests/contract
go mod tidy
```

### Run Provider Verification
```bash
# Run all provider verification tests
cd /Users/weerawat/Tchat/backend/notification/tests/contract
go test -v

# Run specific test
go test -v -run TestNotificationServiceProviderVerification

# Run with coverage
go test -v -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Run Individual Provider State Tests
```bash
# Test provider state setup
go test -v -run TestProviderStateSetup

# Test contract scenarios
go test -v -run TestNotificationServiceContracts
```

## Configuration

### Environment Variables
The tests can be configured with the following environment variables:

```bash
# Pact Broker configuration
export PACT_BROKER_BASE_URL=http://localhost:9292
export PACT_BROKER_USERNAME=pact_user
export PACT_BROKER_PASSWORD=pact_password

# Provider configuration
export PROVIDER_NAME=notification-service
export PROVIDER_VERSION=1.0.0
export PROVIDER_BASE_URL=http://localhost:8080

# Test configuration
export TEST_DATABASE_URL=postgres://test:test@localhost:5432/notification_test
export TEST_CACHE_URL=redis://localhost:6379/1
```

### Pact Broker Integration
To run verification against Pact files from the broker:

```bash
# Verify against latest consumer pacts
go test -v -run TestNotificationServiceProviderVerification \
  -pact-broker-url=http://localhost:9292 \
  -provider-name=notification-service

# Verify against specific consumer version
go test -v -run TestNotificationServiceProviderVerification \
  -consumer-version-selector='{"tag": "main", "latest": true}'
```

## Integration with CI/CD

### GitHub Actions Integration
Add to your `.github/workflows/contract-tests.yml`:

```yaml
notification-provider-verification:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: 1.22

    - name: Install Pact CLI
      run: |
        curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash
        sudo ln -s /tmp/pact/bin/* /usr/local/bin/

    - name: Run Provider Verification
      run: |
        cd backend/notification/tests/contract
        go test -v -run TestNotificationServiceProviderVerification
      env:
        PACT_BROKER_BASE_URL: ${{ secrets.PACT_BROKER_URL }}
        PACT_BROKER_TOKEN: ${{ secrets.PACT_BROKER_TOKEN }}
```

## Localization Testing

The tests include comprehensive localization support:

### Supported Languages
- **English (en)**: Default language
- **Thai (th)**: Primary Southeast Asian market
- **Indonesian (id)**: Mobile-first market
- **Malay (ms)**: Malaysia and Brunei
- **Filipino (tl)**: Philippines market
- **Vietnamese (vi)**: Growing mobile market

### Localized Templates
Templates are tested with localized content for:
- Welcome notifications (onboarding)
- Payment confirmations (commerce)
- Security alerts (authentication)
- System notifications (maintenance, updates)

## Performance Testing

### Response Time Requirements
All notification API endpoints must respond within:
- **Sync operations**: < 200ms p95
- **Async operations**: < 500ms for queuing
- **Template rendering**: < 100ms per template

### Load Testing Integration
The provider tests can be extended for load testing:

```bash
# Run with concurrent requests
go test -v -run TestNotificationServiceProviderVerification \
  -test.parallel 10 \
  -test.count 100
```

## Security Testing

### Authentication Testing
Provider verification includes authentication scenarios:
- Valid JWT tokens
- Expired tokens
- Invalid tokens
- Missing authentication headers

### Data Privacy Testing
Tests verify that sensitive data is not exposed:
- No password/token leakage in responses
- User isolation (users can only see their own data)
- Proper error messages without data exposure

## Troubleshooting

### Common Issues

1. **Pact CLI Not Found**
   ```bash
   # Install Pact CLI
   curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash
   ```

2. **Test Database Connection Issues**
   ```bash
   # Set up test database
   createdb notification_test
   export TEST_DATABASE_URL=postgres://localhost/notification_test
   ```

3. **Mock Service Startup Issues**
   ```bash
   # Check port availability
   lsof -i :8080
   # Kill processes if needed
   kill $(lsof -ti :8080)
   ```

4. **Provider State Setup Failures**
   - Verify mock repositories are properly initialized
   - Check test data creation in setupTestData()
   - Ensure proper UUID parsing in state handlers

### Debug Mode
Run tests with debug output:

```bash
go test -v -run TestNotificationServiceProviderVerification \
  -test.v -test.debug
```

### Logging
Enable detailed logging during tests:

```bash
export PACT_LOG_LEVEL=DEBUG
export GIN_MODE=debug
go test -v
```

## Contributing

### Adding New Provider States
1. Add state handler to `stateHandlers()` map
2. Implement setup logic for the provider state
3. Add test case to `TestProviderStateSetup`
4. Update this documentation

### Adding New API Endpoints
1. Add endpoint to `setupRouter()` function
2. Implement mock response logic
3. Add contract test to `TestNotificationServiceContracts`
4. Update API endpoint documentation

### Localization Support
1. Add new language to supported languages list
2. Update template creation with new localized versions
3. Test with consumer contracts requiring the new language
4. Update documentation

## References

- [Pact Foundation Documentation](https://docs.pact.io/)
- [Pact Go v2 Library](https://github.com/pact-foundation/pact-go)
- [TChat Notification Service API](../../handlers/notification_handler.go)
- [Contract Testing Specification](../../../../specs/021-implement-pact-contract/)

## Support

For issues related to notification service provider verification:
1. Check the troubleshooting section above
2. Review existing GitHub issues
3. Create new issue with test output and configuration details