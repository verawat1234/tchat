# Tchat Backend Testing Standards

**Version**: 1.0
**Last Updated**: 2025-09-22
**Status**: Active Implementation

## Overview

This document establishes comprehensive testing standards for the Tchat Southeast Asian chat platform backend microservices. These standards ensure consistent, maintainable, and high-quality test coverage across all 7 microservices.

## Testing Framework Stack

### Core Testing Dependencies
```go
// Essential testing imports for all test files
import (
    "testing"
    "context"
    "time"

    // Testify framework for assertions and mocking
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
    "github.com/stretchr/testify/mock"

    // Database testing utilities
    "github.com/DATA-DOG/go-sqlmock"
    "net/http/httptest"

    // UUID and JSON utilities
    "github.com/google/uuid"
    "encoding/json"
)
```

### Testing Types & Hierarchy

1. **Unit Tests** (80% coverage target)
   - Individual function/method testing
   - Isolated business logic validation
   - Fast execution (<100ms per test)

2. **Integration Tests** (70% coverage target)
   - Service-to-service communication
   - Database integration workflows
   - End-to-end feature validation

3. **Contract Tests** (95% critical path coverage)
   - API endpoint contracts
   - Data format validation
   - Cross-service compatibility

4. **Performance Tests** (Regional benchmarks)
   - Response time validation (<200ms API, <100ms WebSocket)
   - Load testing (1000+ concurrent users)
   - Regional compliance (6 SEA regions)

## Directory Structure Standards

### Organized Test Hierarchy
```
backend/
├── auth/
│   ├── handlers/
│   │   ├── auth_handlers.go
│   │   └── auth_handlers_test.go          # Handler layer tests
│   ├── services/
│   │   ├── auth_service.go
│   │   └── auth_service_test.go           # Business logic tests
│   ├── repositories/
│   │   ├── user_repository.go
│   │   └── user_repository_test.go        # Data layer tests
│   └── mocks/                             # Auto-generated mocks
│       ├── mock_auth_service.go
│       └── mock_user_repository.go
├── messaging/                             # Similar structure for all services
├── payment/
├── commerce/
├── notification/
├── content/
├── gateway/
└── tests/
    ├── fixtures/                          # Shared test data
    ├── integration/                       # Cross-service tests
    ├── contract/                          # API contract tests
    ├── performance/                       # Load and benchmark tests
    └── compliance/                        # Regional compliance tests
```

### File Naming Conventions

**Test Files**: `*_test.go` suffix (Go standard)
```
auth_handlers_test.go      ✅ Correct
user_repository_test.go    ✅ Correct
test_auth.go              ❌ Incorrect
auth_test.go              ❌ Too generic
```

**Test Functions**: `Test[Type]_[Method]_[Scenario]_[ExpectedResult]`
```go
// ✅ Excellent naming - clear and descriptive
func TestUserService_CreateUser_ValidInput_ReturnsSuccess(t *testing.T)
func TestUserService_CreateUser_DuplicatePhone_ReturnsError(t *testing.T)
func TestUserService_CreateUser_InvalidCountryCode_ReturnsValidationError(t *testing.T)

// ❌ Poor naming - ambiguous or unclear
func TestCreateUser(t *testing.T)
func TestUser(t *testing.T)
func TestValidation(t *testing.T)
```

## Test Implementation Patterns

### 1. Unit Test Structure (AAA Pattern)

```go
func TestUserService_CreateUser_ValidInput_ReturnsSuccess(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    validUser := models.CreateUserRequest{
        Phone:    "+66891234567",
        Country:  "TH",
        Language: "th",
    }

    expectedUser := &models.User{
        ID:       uuid.New(),
        Phone:    validUser.Phone,
        Country:  validUser.Country,
        Language: validUser.Language,
        Status:   models.UserStatusActive,
    }

    mockRepo.On("Create", mock.MatchedBy(func(user *models.User) bool {
        return user.Phone == validUser.Phone
    })).Return(expectedUser, nil)

    // Act
    result, err := service.CreateUser(context.Background(), validUser)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, validUser.Phone, result.Phone)
    assert.Equal(t, validUser.Country, result.Country)
    assert.Equal(t, models.UserStatusActive, result.Status)

    mockRepo.AssertExpectations(t)
}
```

### 2. Repository Test Pattern (Database Layer)

```go
type UserRepositoryTestSuite struct {
    suite.Suite
    db     *sql.DB
    mock   sqlmock.Sqlmock
    repo   *PostgreSQLUserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
    db, mock, err := sqlmock.New()
    require.NoError(suite.T(), err)

    suite.db = db
    suite.mock = mock
    suite.repo = NewPostgreSQLUserRepository(db)
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
    suite.db.Close()
}

func (suite *UserRepositoryTestSuite) TestCreateUser_ValidUser_InsertsSuccessfully() {
    // Arrange
    user := &models.User{
        Phone:    "+66891234567",
        Country:  "TH",
        Language: "th",
    }

    suite.mock.ExpectExec(`INSERT INTO users`).
        WithArgs(sqlmock.AnyArg(), user.Phone, user.Country, user.Language, sqlmock.AnyArg()).
        WillReturnResult(sqlmock.NewResult(1, 1))

    // Act
    err := suite.repo.Create(context.Background(), user)

    // Assert
    assert.NoError(suite.T(), err)
    assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func TestUserRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(UserRepositoryTestSuite))
}
```

### 3. HTTP Handler Test Pattern

```go
func TestAuthHandlers_Login_ValidCredentials_ReturnsToken(t *testing.T) {
    // Arrange
    gin.SetMode(gin.TestMode)
    router := gin.New()

    mockService := new(MockAuthService)
    handlers := NewAuthHandlers(mockService)

    router.POST("/auth/login", handlers.Login)

    loginReq := models.LoginRequest{
        Phone:    "+66891234567",
        OTPCode:  "123456",
    }

    expectedResponse := &models.LoginResponse{
        AccessToken:  "jwt_access_token",
        RefreshToken: "jwt_refresh_token",
        User: models.User{
            ID:       uuid.New(),
            Phone:    loginReq.Phone,
            Status:   models.UserStatusActive,
        },
    }

    mockService.On("Login", mock.Anything, loginReq).Return(expectedResponse, nil)

    // Act
    jsonData, _ := json.Marshal(loginReq)
    req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "success", response["status"])
    assert.NotNil(t, response["data"])

    mockService.AssertExpectations(t)
}
```

### 4. Integration Test Pattern

```go
func TestAuthFlow_CompleteRegistration_Success(t *testing.T) {
    // Arrange - Setup test environment
    testDB := setupTestDatabase(t)
    defer cleanupTestDatabase(testDB)

    testServer := setupTestServer(testDB)
    defer testServer.Close()

    client := &http.Client{Timeout: 10 * time.Second}

    // Step 1: Request OTP
    otpReq := models.RequestOTPRequest{
        Phone:   "+66891234567",
        Country: "TH",
    }

    resp, err := sendJSONRequest(client, "POST", testServer.URL+"/auth/otp/request", otpReq)
    require.NoError(t, err)
    require.Equal(t, http.StatusOK, resp.StatusCode)

    // Step 2: Verify OTP (simulate)
    verifyReq := models.VerifyOTPRequest{
        Phone:   otpReq.Phone,
        OTPCode: "123456", // Test OTP
    }

    resp, err = sendJSONRequest(client, "POST", testServer.URL+"/auth/otp/verify", verifyReq)
    require.NoError(t, err)
    require.Equal(t, http.StatusOK, resp.StatusCode)

    // Assert complete flow worked
    var loginResp models.LoginResponse
    err = json.NewDecoder(resp.Body).Decode(&loginResp)
    require.NoError(t, err)
    assert.NotEmpty(t, loginResp.AccessToken)
    assert.Equal(t, otpReq.Phone, loginResp.User.Phone)
}
```

## Mock Implementation Standards

### Service Layer Mocks

```go
// MockAuthService provides mock implementation for testing
type MockAuthService struct {
    mock.Mock
}

func (m *MockAuthService) RequestOTP(ctx context.Context, req models.RequestOTPRequest) error {
    args := m.Called(ctx, req)
    return args.Error(0)
}

func (m *MockAuthService) VerifyOTP(ctx context.Context, req models.VerifyOTPRequest) (*models.LoginResponse, error) {
    args := m.Called(ctx, req)
    return args.Get(0).(*models.LoginResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, token string) (*models.TokenResponse, error) {
    args := m.Called(ctx, token)
    return args.Get(0).(*models.TokenResponse), args.Error(1)
}
```

### Repository Layer Mocks

```go
// MockUserRepository provides mock implementation for testing
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
    args := m.Called(ctx, phone)
    return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}
```

## Test Data Management

### Fixtures and Test Data

**File**: `backend/tests/fixtures/users.go`
```go
package fixtures

import (
    "time"
    "github.com/google/uuid"
    "tchat.dev/shared/models"
)

// ValidThailandUser creates a valid Thai user for testing
func ValidThailandUser() *models.User {
    return &models.User{
        ID:        uuid.New(),
        Phone:     "+66891234567",
        Country:   "TH",
        Language:  "th",
        Status:    models.UserStatusActive,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

// ValidSingaporeUser creates a valid Singaporean user for testing
func ValidSingaporeUser() *models.User {
    return &models.User{
        ID:        uuid.New(),
        Phone:     "+6591234567",
        Country:   "SG",
        Language:  "en",
        Status:    models.UserStatusActive,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

// CreateContentItems returns test content for different regions
func CreateContentItems() []models.ContentItem {
    return []models.ContentItem{
        {
            ID:       uuid.New(),
            Category: "navigation",
            Type:     models.ContentTypeText,
            Value:    models.ContentValue{"text": "หน้าหลัก", "lang": "th"},
            Status:   models.ContentStatusPublished,
        },
        {
            ID:       uuid.New(),
            Category: "navigation",
            Type:     models.ContentTypeText,
            Value:    models.ContentValue{"text": "Home", "lang": "en"},
            Status:   models.ContentStatusPublished,
        },
    }
}
```

### Database Test Setup

**File**: `backend/tests/fixtures/database.go`
```go
package fixtures

import (
    "database/sql"
    "testing"
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/require"
)

// SetupMockDB creates a mock database for testing
func SetupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    return db, mock
}

// ExpectUserCreation sets up common user creation expectations
func ExpectUserCreation(mock sqlmock.Sqlmock, user *models.User) {
    mock.ExpectExec(`INSERT INTO users`).
        WithArgs(
            sqlmock.AnyArg(), // ID
            user.Phone,
            user.Country,
            user.Language,
            user.Status,
            sqlmock.AnyArg(), // CreatedAt
            sqlmock.AnyArg(), // UpdatedAt
        ).
        WillReturnResult(sqlmock.NewResult(1, 1))
}
```

## Coverage Requirements & Quality Gates

### Coverage Targets by Service

| Service | Unit Tests | Integration Tests | Critical Path | Performance |
|---------|------------|-------------------|---------------|-------------|
| Auth Service | ≥80% | ≥70% | ≥95% | <150ms login |
| Messaging Service | ≥80% | ≥70% | ≥95% | <100ms delivery |
| Payment Service | ≥85% | ≥75% | ≥98% | <200ms transaction |
| Commerce Service | ≥80% | ≥70% | ≥95% | <250ms catalog |
| Notification Service | ≥80% | ≥70% | ≥90% | <500ms delivery |
| Content Service | ≥80% | ≥70% | ≥95% | <100ms retrieval |
| API Gateway | ≥75% | ≥80% | ≥95% | <50ms routing |

### Quality Gates (Must Pass)

```yaml
Pre-Commit Gates:
  - unit_tests_pass: required
  - code_coverage_threshold: ≥80%
  - linting_clean: required
  - security_scan_clean: required

Pre-Deploy Gates:
  - integration_tests_pass: required
  - performance_benchmarks_met: required
  - contract_tests_pass: required
  - regional_compliance_validated: required
```

## Error Handling & Edge Cases

### Error Testing Patterns

```go
func TestUserService_CreateUser_DatabaseError_ReturnsError(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    user := fixtures.ValidThailandUser()
    expectedError := errors.New("database connection failed")

    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).
        Return(expectedError)

    // Act
    result, err := service.CreateUser(context.Background(), models.CreateUserRequest{
        Phone:   user.Phone,
        Country: user.Country,
    })

    // Assert
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Contains(t, err.Error(), "database connection failed")

    mockRepo.AssertExpectations(t)
}
```

### Input Validation Testing

```go
func TestUserService_CreateUser_InvalidPhone_ReturnsValidationError(t *testing.T) {
    testCases := []struct {
        name        string
        phone       string
        expectedErr string
    }{
        {
            name:        "Empty phone number",
            phone:       "",
            expectedErr: "phone number is required",
        },
        {
            name:        "Invalid format",
            phone:       "123456789",
            expectedErr: "invalid phone number format",
        },
        {
            name:        "Invalid country code",
            phone:       "+1234567890",
            expectedErr: "unsupported country code",
        },
    }

    service := NewUserService(new(MockUserRepository))

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Act
            result, err := service.CreateUser(context.Background(), models.CreateUserRequest{
                Phone:   tc.phone,
                Country: "TH",
            })

            // Assert
            assert.Error(t, err)
            assert.Nil(t, result)
            assert.Contains(t, err.Error(), tc.expectedErr)
        })
    }
}
```

## Regional Compliance Testing

### Southeast Asian Regional Tests

```go
func TestUserService_CreateUser_RegionalCompliance_Success(t *testing.T) {
    testCases := []struct {
        name     string
        country  string
        phone    string
        language string
        dataLaws []string
    }{
        {
            name:     "Thailand PDPA compliance",
            country:  "TH",
            phone:    "+66891234567",
            language: "th",
            dataLaws: []string{"PDPA_TH", "GDPR"},
        },
        {
            name:     "Singapore PDPA compliance",
            country:  "SG",
            phone:    "+6591234567",
            language: "en",
            dataLaws: []string{"PDPA_SG", "GDPR"},
        },
        {
            name:     "Indonesia PDP compliance",
            country:  "ID",
            phone:    "+628123456789",
            language: "id",
            dataLaws: []string{"PDP_ID", "GDPR"},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Test data residency requirements
            assert.True(t, isDataResidencyCompliant(tc.country))

            // Test privacy law compliance
            for _, law := range tc.dataLaws {
                assert.True(t, isPrivacyCompliant(tc.country, law))
            }

            // Test localization support
            assert.True(t, isLanguageSupported(tc.language))
        })
    }
}
```

## Performance Testing Standards

### Response Time Benchmarks

```go
func BenchmarkAuthService_Login(b *testing.B) {
    service := setupAuthService()
    loginReq := models.LoginRequest{
        Phone:   "+66891234567",
        OTPCode: "123456",
    }

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        _, err := service.Login(context.Background(), loginReq)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func TestAuthService_Login_PerformanceRequirement(t *testing.T) {
    service := setupAuthService()
    loginReq := models.LoginRequest{
        Phone:   "+66891234567",
        OTPCode: "123456",
    }

    start := time.Now()
    _, err := service.Login(context.Background(), loginReq)
    duration := time.Since(start)

    assert.NoError(t, err)
    assert.Less(t, duration.Milliseconds(), int64(150), "Login should complete within 150ms")
}
```

## Test Execution & CI/CD Integration

### Running Tests

```bash
# Unit tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Integration tests
go test -v -tags=integration ./backend/tests/integration/...

# Performance benchmarks
go test -v -bench=. -benchtime=10s ./backend/tests/performance/...

# Coverage report
go tool cover -html=coverage.out -o coverage.html

# Coverage percentage
go tool cover -func=coverage.out | tail -1
```

### CI/CD Pipeline Configuration

**File**: `.github/workflows/tests.yml`
```yaml
name: Backend Tests
on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: 1.22

      - name: Run Unit Tests
        run: |
          cd backend
          go test -v -race -coverprofile=coverage.out ./...

      - name: Check Coverage
        run: |
          cd backend
          COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
          echo "Coverage: $COVERAGE%"
          if (( $(echo "$COVERAGE < 80" | bc -l) )); then
            echo "Coverage $COVERAGE% is below required 80%"
            exit 1
          fi

      - name: Upload Coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./backend/coverage.out

  integration-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test_password
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      redis:
        image: redis:7
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: 1.22

      - name: Run Integration Tests
        run: |
          cd backend
          go test -v -tags=integration ./tests/integration/...
```

## Common Anti-Patterns (Avoid These)

### ❌ Bad Testing Practices

```go
// DON'T: Test implementation details
func TestUserService_CreateUser_CallsRepositoryCreate(t *testing.T) {
    // This tests HOW, not WHAT
}

// DON'T: Giant test functions
func TestUserServiceEverything(t *testing.T) {
    // Tests 20+ scenarios in one function
}

// DON'T: Unclear test names
func TestCreateUser(t *testing.T) {
    // Which scenario? What's expected?
}

// DON'T: No assertions
func TestCreateUser_ValidInput(t *testing.T) {
    service.CreateUser(ctx, req)
    // Missing assertions!
}

// DON'T: Hardcoded values everywhere
func TestCreateUser_ValidInput(t *testing.T) {
    user := models.User{
        Phone: "+66891234567",
        Country: "TH",
        // ... repeated in every test
    }
}
```

### ✅ Good Testing Practices

```go
// DO: Test behavior, not implementation
func TestUserService_CreateUser_ValidInput_ReturnsActiveUser(t *testing.T) {
    // Tests WHAT the service does
}

// DO: Single responsibility per test
func TestUserService_CreateUser_DuplicatePhone_ReturnsError(t *testing.T) {
    // One scenario, clear expectation
}

// DO: Use test fixtures
func TestUserService_CreateUser_ValidInput_ReturnsActiveUser(t *testing.T) {
    user := fixtures.ValidThailandUser()
    // Reusable, maintainable test data
}

// DO: Clear assertions
func TestUserService_CreateUser_ValidInput_ReturnsActiveUser(t *testing.T) {
    result, err := service.CreateUser(ctx, req)

    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, models.UserStatusActive, result.Status)
}
```

## Best Practices Summary

### 🎯 Testing Principles

1. **Test Behavior, Not Implementation**: Focus on what the code does, not how it does it
2. **AAA Pattern**: Arrange, Act, Assert structure for all tests
3. **Single Responsibility**: One test case per scenario
4. **Descriptive Names**: Clear test names that describe scenario and expectation
5. **Fast & Reliable**: Tests should run quickly and consistently
6. **Independent**: Tests should not depend on each other
7. **Maintainable**: Easy to understand and modify as code evolves

### 📊 Quality Metrics

- **Unit Test Coverage**: ≥80% for all services
- **Integration Test Coverage**: ≥70% for critical workflows
- **Critical Path Coverage**: ≥95% for auth, payment, messaging
- **Test Execution Time**: <5 minutes for full suite
- **Test Reliability**: ≥98% pass rate in CI/CD

### 🌏 Regional Considerations

- **Data Residency**: Test data storage compliance for TH, SG, ID, MY, PH, VN
- **Privacy Laws**: Validate GDPR, PDPA compliance in all regions
- **Performance**: Regional response time targets (<200ms average)
- **Localization**: Multi-language support and cultural adaptations

### 🔄 Continuous Improvement

- **Regular Reviews**: Monthly testing standard reviews
- **Metric Tracking**: Weekly coverage and performance metrics
- **Pattern Evolution**: Quarterly pattern updates based on learnings
- **Team Training**: Ongoing education on testing best practices

---

This testing standards document serves as the foundation for implementing comprehensive test coverage across the Tchat platform, ensuring high-quality, maintainable code that meets enterprise standards and Southeast Asian regional requirements.