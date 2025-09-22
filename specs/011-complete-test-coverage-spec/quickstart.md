# Quickstart: Test Coverage Implementation

## Prerequisites

1. **Go Development Environment**
   ```bash
   go version  # Should be 1.22+
   ```

2. **Testing Dependencies**
   ```bash
   go get github.com/stretchr/testify/suite
   go get github.com/stretchr/testify/mock
   go get github.com/stretchr/testify/assert
   go get github.com/DATA-DOG/go-sqlmock
   ```

3. **Infrastructure Services Running**
   ```bash
   make dev-backend  # Start PostgreSQL, Redis, Kafka
   ```

## Quick Start (5 minutes)

### 1. Content Service Unit Tests (Phase 1 - Critical)

Create the missing Content Service test suite:

```bash
# Navigate to Content Service
cd backend/content

# Create test file for handlers
cat > handlers/content_handlers_test.go << 'EOF'
package handlers

import (
    "testing"
    "github.com/stretchr/testify/suite"
    "github.com/stretchr/testify/mock"
)

type ContentHandlersTestSuite struct {
    suite.Suite
    mockService *MockContentService
    handlers    *ContentHandlers
}

func (suite *ContentHandlersTestSuite) SetupTest() {
    suite.mockService = new(MockContentService)
    suite.handlers = NewContentHandlers(suite.mockService)
}

func TestContentHandlersTestSuite(t *testing.T) {
    suite.Run(t, new(ContentHandlersTestSuite))
}

// Test placeholder - will fail until implementation
func (suite *ContentHandlersTestSuite) TestCreateContent_ValidRequest_ReturnsSuccess() {
    suite.T().Skip("TODO: Implement Content Service tests")
}
EOF

# Run the test (should fail - TDD approach)
go test ./handlers/ -v
```

**Expected Output**: Test should **FAIL** with skip message - this is correct for TDD approach.

### 2. Verify Current Test Structure

```bash
# From project root
cd backend

# Check existing test structure
find . -name "*_test.go" | head -10

# Run existing contract tests
go test ./tests/contract/auth_otp_test.go -v

# Check coverage for existing services
go test -cover ./auth/... ./messaging/...
```

**Expected Output**: Should show existing 29 test files and current coverage metrics.

### 3. Test Infrastructure Validation

```bash
# Check Docker services
make ps

# Validate service health
make health

# Run integration tests
go test ./tests/integration/auth_flow_test.go -v
```

**Expected Output**: All infrastructure services healthy, integration tests pass.

## Development Workflow (15 minutes)

### 1. Unit Test Development Pattern

Follow this TDD pattern for each service:

```bash
# 1. Create failing test first
cat > services/content_service_test.go << 'EOF'
package services

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockContentRepository struct {
    mock.Mock
}

func (m *MockContentRepository) Create(content *ContentItem) error {
    args := m.Called(content)
    return args.Error(0)
}

func TestContentService_CreateContent_ValidInput_ReturnsSuccess(t *testing.T) {
    // Arrange
    mockRepo := new(MockContentRepository)
    service := NewContentService(mockRepo)
    content := &ContentItem{Title: "Test Content"}

    mockRepo.On("Create", content).Return(nil)

    // Act
    err := service.CreateContent(content)

    // Assert
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
EOF

# 2. Run test (should fail)
go test ./services/ -v

# 3. Implement minimal code to make test pass
# (Implementation step)

# 4. Refactor and add more tests
```

### 2. Coverage Measurement

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Check coverage percentage
go tool cover -func=coverage.out | tail -1
```

### 3. Integration Test Setup

```bash
# Create integration test
mkdir -p tests/integration
cat > tests/integration/content_flow_test.go << 'EOF'
package integration

import (
    "testing"
    "github.com/stretchr/testify/suite"
)

type ContentIntegrationTestSuite struct {
    suite.Suite
    // Database, HTTP client setup
}

func (suite *ContentIntegrationTestSuite) SetupSuite() {
    // Setup test database, start service
}

func (suite *ContentIntegrationTestSuite) TestContentWorkflow_CreatePublishRetrieve_Success() {
    suite.T().Skip("TODO: Implement integration test")
}

func TestContentIntegrationTestSuite(t *testing.T) {
    suite.Run(t, new(ContentIntegrationTestSuite))
}
EOF
```

## Validation Workflow (10 minutes)

### 1. Coverage Validation

```bash
# Check current coverage
make test-backend

# Validate coverage targets
go test -cover ./content/... | grep -E "coverage: [0-9]+\.[0-9]+%"

# Should show progression toward 80% unit test target
```

### 2. Quality Gates Check

```bash
# Run linting
go vet ./...
golangci-lint run

# Security scan
gosec ./...

# Performance benchmarks
go test -bench=. ./tests/performance/
```

### 3. Regional Compliance Test

```bash
# Test data residency compliance
go test ./tests/integration/ -run "TestRegional.*" -v

# Validate Southeast Asian specific features
go test ./tests/contract/ -run "TestPayment.*Thailand|Singapore|Indonesia" -v
```

## Success Criteria Validation

### Phase 1 Success (Content Service)
- [ ] Content Service has >10 unit tests
- [ ] All tests follow testify/suite pattern
- [ ] Tests are currently failing (TDD)
- [ ] Coverage measurement working
- [ ] Integration test structure created

### Quick Validation Commands
```bash
# Count Content Service tests
find backend/content -name "*_test.go" -exec grep -l "func Test" {} \; | wc -l

# Check test execution
go test ./backend/content/... -v | grep -E "PASS|FAIL|SKIP"

# Verify coverage measurement
go test -cover ./backend/content/... | grep "coverage:"
```

## Troubleshooting

### Common Issues

1. **Import Path Errors**
   ```bash
   # Fix module imports
   go mod tidy
   go mod verify
   ```

2. **Database Connection Failures**
   ```bash
   # Restart infrastructure
   make infra-down
   make infra-up
   ```

3. **Test Discovery Issues**
   ```bash
   # Verify test file naming
   ls -la *_test.go
   # Files must end with _test.go
   ```

4. **Mock Generation Failures**
   ```bash
   # Install mockery for auto-generation
   go install github.com/vektra/mockery/v2@latest

   # Generate mocks
   mockery --all --output mocks/
   ```

## Next Steps

After completing this quickstart:

1. **Proceed to Phase 2**: Commerce and Notification Service tests
2. **Establish CI/CD**: Add automated test execution
3. **Security Testing**: Implement auth and input validation tests
4. **Performance Baseline**: Create regional benchmark tests

## Production Readiness Checklist

- [ ] All services have ≥80% unit test coverage
- [ ] Integration tests cover critical workflows
- [ ] Security tests validate auth and compliance
- [ ] Performance tests meet regional targets
- [ ] CI/CD pipeline enforces quality gates
- [ ] Coverage reports generated automatically

**Total Time Investment**: ~30 minutes for quickstart, 3-4 sprints for full implementation.

This quickstart provides immediate validation of the testing approach and establishes the foundation for comprehensive test coverage across the Tchat platform.