# Tchat Backend Development Guide

This guide covers development setup, coding standards, testing procedures, and contribution guidelines for the Tchat backend services.

## Table of Contents

1. [Development Setup](#development-setup)
2. [Project Structure](#project-structure)
3. [Coding Standards](#coding-standards)
4. [Testing Guidelines](#testing-guidelines)
5. [Database Development](#database-development)
6. [API Development](#api-development)
7. [Contributing](#contributing)
8. [Debugging](#debugging)

## Development Setup

### Prerequisites

- Go 1.22 or later
- Docker and Docker Compose
- Make
- Git
- IDE with Go support (VS Code, GoLand, etc.)

### Local Environment Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd tchat/backend
   ```

2. **Start development infrastructure**
   ```bash
   # Start PostgreSQL, Redis, ScyllaDB, and Kafka
   docker-compose -f docker-compose.dev.yml up -d

   # Wait for services to be ready
   make wait-for-services
   ```

3. **Install dependencies**
   ```bash
   go mod download
   go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

4. **Setup environment variables**
   ```bash
   cp .env.example .env.local
   # Edit .env.local with your local configuration
   ```

5. **Run database migrations**
   ```bash
   make migrate-up
   ```

6. **Start services**
   ```bash
   # Start all services
   make dev

   # Or start individual services
   make dev-auth
   make dev-messaging
   make dev-payment
   make dev-notification
   ```

### Environment Configuration

```bash
# .env.local
ENVIRONMENT=development
LOG_LEVEL=debug

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=tchat_dev
POSTGRES_PASSWORD=dev_password
POSTGRES_SSL_MODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# ScyllaDB
SCYLLA_HOSTS=localhost:9042
SCYLLA_CONSISTENCY=one

# Kafka
KAFKA_BROKERS=localhost:9092

# External Services (use test credentials)
TWILIO_ACCOUNT_SID=test_sid
TWILIO_AUTH_TOKEN=test_token
STRIPE_SECRET_KEY=sk_test_...
OMISE_SECRET_KEY=skey_test_...
```

### Development Tools

```bash
# Install development tools
make install-tools

# This installs:
# - golangci-lint (code linting)
# - migrate (database migrations)
# - air (live reload)
# - mockgen (mock generation)
# - protoc-gen-go (Protocol Buffers)
```

## Project Structure

```
backend/
├── cmd/                    # Service entry points
│   ├── auth/              # Auth service main
│   ├── messaging/         # Messaging service main
│   ├── payment/           # Payment service main
│   ├── notification/      # Notification service main
│   └── gateway/           # API Gateway main
├── internal/              # Private application code
│   ├── auth/              # Auth service implementation
│   ├── messaging/         # Messaging service implementation
│   ├── payment/           # Payment service implementation
│   └── notification/      # Notification service implementation
├── shared/                # Shared packages
│   ├── config/            # Configuration management
│   ├── database/          # Database clients
│   ├── cache/             # Caching utilities
│   ├── messaging/         # Event streaming
│   ├── external/          # External service clients
│   ├── middleware/        # HTTP middleware
│   ├── validation/        # Input validation
│   └── utils/             # Utility functions
├── migrations/            # Database migrations
│   ├── postgres/          # PostgreSQL migrations
│   └── scylla/            # ScyllaDB schema
├── tests/                 # Test utilities and data
│   ├── testutil/          # Test helpers
│   ├── contract/          # Contract tests
│   ├── integration/       # Integration tests
│   └── performance/       # Load tests
├── docs/                  # Documentation
├── scripts/               # Development scripts
├── deployments/           # Deployment configurations
└── tools/                 # Development tools
```

### Service Architecture

Each service follows this structure:

```
internal/auth/
├── handlers/              # HTTP handlers
│   ├── otp.go            # OTP endpoint handlers
│   ├── auth.go           # Authentication handlers
│   └── middleware.go     # Service-specific middleware
├── services/              # Business logic
│   ├── auth_service.go   # Auth business logic
│   ├── otp_service.go    # OTP business logic
│   └── user_service.go   # User management logic
├── repositories/          # Data access layer
│   ├── user_repo.go      # User repository
│   └── session_repo.go   # Session repository
├── models/                # Domain models
│   ├── user.go           # User model
│   ├── session.go        # Session model
│   └── otp.go            # OTP model
└── server.go              # HTTP server setup
```

## Coding Standards

### Go Style Guide

We follow the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md) with these additions:

#### Package Organization

```go
// Package declaration with clear purpose
// Package auth provides authentication and authorization services
// for the Tchat messaging platform.
package auth

import (
    // Standard library imports first
    "context"
    "fmt"
    "time"

    // Third-party imports
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    // Local imports last
    "github.com/tchat/backend/shared/config"
    "github.com/tchat/backend/shared/database"
)
```

#### Error Handling

```go
// Define package-specific errors
var (
    ErrUserNotFound    = errors.New("user not found")
    ErrInvalidOTP      = errors.New("invalid OTP code")
    ErrOTPExpired      = errors.New("OTP code has expired")
)

// Use error wrapping for context
func (s *AuthService) VerifyOTP(ctx context.Context, phone, code string) error {
    user, err := s.userRepo.GetByPhone(ctx, phone)
    if err != nil {
        return fmt.Errorf("failed to get user by phone %s: %w", phone, err)
    }

    if user == nil {
        return ErrUserNotFound
    }

    // Business logic...
    return nil
}
```

#### Struct Design

```go
// Use clear, descriptive names
type AuthService struct {
    userRepo    repositories.UserRepository
    otpService  OTPService
    jwtManager  *jwt.Manager
    logger      *slog.Logger
}

// Constructor with interface dependencies
func NewAuthService(
    userRepo repositories.UserRepository,
    otpService OTPService,
    jwtManager *jwt.Manager,
    logger *slog.Logger,
) *AuthService {
    return &AuthService{
        userRepo:   userRepo,
        otpService: otpService,
        jwtManager: jwtManager,
        logger:     logger,
    }
}
```

#### Interface Design

```go
// Keep interfaces small and focused
type UserRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
    GetByPhone(ctx context.Context, phone string) (*models.User, error)
    Create(ctx context.Context, user *models.User) error
    Update(ctx context.Context, user *models.User) error
}

type OTPService interface {
    Send(ctx context.Context, phone, countryCode string) error
    Verify(ctx context.Context, phone, code string) error
}
```

### Code Quality Tools

#### Linting Configuration

```yaml
# .golangci.yml
run:
  timeout: 5m
  modules-download-mode: readonly

linters:
  enable:
    - gofmt
    - goimports
    - govet
    - ineffassign
    - misspell
    - revive
    - staticcheck
    - unused
    - errcheck
    - gosimple
    - typecheck

linters-settings:
  revive:
    rules:
      - name: exported
        arguments: [true]
      - name: package-comments
      - name: var-naming
```

#### Pre-commit Hooks

```bash
#!/bin/sh
# .git/hooks/pre-commit

# Run linting
golangci-lint run

# Run tests
go test ./...

# Check for security issues
gosec ./...

# Format code
gofmt -w .
goimports -w .
```

### Documentation Standards

#### Code Comments

```go
// Package-level documentation
// Package auth provides authentication services including OTP verification,
// JWT token management, and user session handling for the Tchat platform.
package auth

// AuthService handles user authentication and authorization.
// It provides methods for OTP verification, JWT token generation,
// and user session management.
type AuthService struct {
    // Dependencies are injected through the constructor
    userRepo   repositories.UserRepository
    otpService OTPService
}

// VerifyOTP verifies an OTP code for the given phone number.
// It checks if the code is valid and not expired, then marks it as used.
//
// Returns ErrInvalidOTP if the code is invalid or ErrOTPExpired if expired.
func (s *AuthService) VerifyOTP(ctx context.Context, phone, code string) error {
    // Implementation...
}
```

#### API Documentation

```go
// SendOTP sends an OTP code to the specified phone number
// @Summary Send OTP code
// @Description Send a one-time password to the phone number for authentication
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.SendOTPRequest true "OTP request"
// @Success 200 {object} models.SendOTPResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 429 {object} models.ErrorResponse
// @Router /auth/otp/send [post]
func (h *AuthHandler) SendOTP(c *gin.Context) {
    // Implementation...
}
```

## Testing Guidelines

### Test Organization

```
tests/
├── unit/                  # Unit tests (next to source files)
├── integration/           # Integration tests
├── contract/             # API contract tests
├── performance/          # Load and performance tests
└── testutil/             # Test utilities and helpers
```

### Unit Testing

```go
// auth_service_test.go
package auth_test

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"

    "github.com/tchat/backend/internal/auth"
    "github.com/tchat/backend/internal/auth/mocks"
    "github.com/tchat/backend/internal/auth/models"
)

func TestAuthService_VerifyOTP(t *testing.T) {
    tests := []struct {
        name          string
        phone         string
        code          string
        setupMocks    func(*mocks.UserRepository, *mocks.OTPService)
        expectedError error
    }{
        {
            name:  "valid OTP",
            phone: "+1234567890",
            code:  "123456",
            setupMocks: func(userRepo *mocks.UserRepository, otpSvc *mocks.OTPService) {
                userRepo.On("GetByPhone", mock.Anything, "+1234567890").
                    Return(&models.User{ID: uuid.New(), Phone: "+1234567890"}, nil)
                otpSvc.On("Verify", mock.Anything, "+1234567890", "123456").
                    Return(nil)
            },
            expectedError: nil,
        },
        {
            name:  "invalid OTP",
            phone: "+1234567890",
            code:  "invalid",
            setupMocks: func(userRepo *mocks.UserRepository, otpSvc *mocks.OTPService) {
                userRepo.On("GetByPhone", mock.Anything, "+1234567890").
                    Return(&models.User{ID: uuid.New(), Phone: "+1234567890"}, nil)
                otpSvc.On("Verify", mock.Anything, "+1234567890", "invalid").
                    Return(auth.ErrInvalidOTP)
            },
            expectedError: auth.ErrInvalidOTP,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mocks
            userRepo := mocks.NewUserRepository(t)
            otpService := mocks.NewOTPService(t)

            tt.setupMocks(userRepo, otpService)

            // Create service
            authService := auth.NewAuthService(userRepo, otpService, nil, nil)

            // Execute
            err := authService.VerifyOTP(context.Background(), tt.phone, tt.code)

            // Assert
            if tt.expectedError != nil {
                assert.ErrorIs(t, err, tt.expectedError)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Integration Testing

```go
// auth_integration_test.go
package integration_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/tchat/backend/tests/testutil"
)

func TestAuthIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    // Setup test environment
    helper := testutil.NewTestHelper(t)
    defer helper.Cleanup()

    server := helper.StartAuthService()
    defer server.Close()

    t.Run("OTP flow", func(t *testing.T) {
        // Send OTP
        sendReq := map[string]interface{}{
            "phone_number":  "1234567890",
            "country_code":  "US",
            "locale":        "en",
        }

        reqBody, _ := json.Marshal(sendReq)
        resp, err := http.Post(
            server.URL+"/auth/otp/send",
            "application/json",
            bytes.NewBuffer(reqBody),
        )
        require.NoError(t, err)
        defer resp.Body.Close()

        assert.Equal(t, http.StatusOK, resp.StatusCode)

        // Verify OTP (using test code)
        verifyReq := map[string]interface{}{
            "phone_number":  "1234567890",
            "country_code":  "US",
            "otp_code":      "123456", // Test OTP code
        }

        reqBody, _ = json.Marshal(verifyReq)
        resp, err = http.Post(
            server.URL+"/auth/otp/verify",
            "application/json",
            bytes.NewBuffer(reqBody),
        )
        require.NoError(t, err)
        defer resp.Body.Close()

        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var verifyResp map[string]interface{}
        err = json.NewDecoder(resp.Body).Decode(&verifyResp)
        require.NoError(t, err)

        assert.Contains(t, verifyResp, "access_token")
        assert.Contains(t, verifyResp, "refresh_token")
    })
}
```

### Test Data Management

```go
// testutil/fixtures.go
package testutil

import (
    "github.com/google/uuid"
    "github.com/tchat/backend/internal/auth/models"
)

// UserFixtures provides test user data
type UserFixtures struct{}

func (f *UserFixtures) ValidUser() *models.User {
    return &models.User{
        ID:          uuid.New(),
        Phone:       "+1234567890",
        CountryCode: "US",
        Status:      models.UserStatusActive,
        Profile: models.UserProfile{
            DisplayName: "Test User",
            Locale:      "en",
        },
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

func (f *UserFixtures) CreateUsers(count int) []*models.User {
    users := make([]*models.User, count)
    for i := 0; i < count; i++ {
        users[i] = &models.User{
            ID:          uuid.New(),
            Phone:       fmt.Sprintf("+123456%04d", i),
            CountryCode: "US",
            Status:      models.UserStatusActive,
            CreatedAt:   time.Now(),
            UpdatedAt:   time.Now(),
        }
    }
    return users
}
```

### Performance Testing

```go
// performance/auth_load_test.go
package performance_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "sync"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestAuthServiceLoad(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping load test in short mode")
    }

    const (
        concurrentUsers = 100
        requestsPerUser = 10
        targetRPS       = 500
    )

    var (
        totalRequests   = concurrentUsers * requestsPerUser
        successCount    int64
        errorCount      int64
        responseTimes   []time.Duration
        mutex           sync.Mutex
    )

    start := time.Now()

    var wg sync.WaitGroup
    for i := 0; i < concurrentUsers; i++ {
        wg.Add(1)
        go func(userID int) {
            defer wg.Done()

            for j := 0; j < requestsPerUser; j++ {
                reqStart := time.Now()

                // Send OTP request
                reqBody, _ := json.Marshal(map[string]interface{}{
                    "phone_number": fmt.Sprintf("123456%04d", userID),
                    "country_code": "US",
                    "locale":       "en",
                })

                resp, err := http.Post(
                    "http://localhost:8081/auth/otp/send",
                    "application/json",
                    bytes.NewBuffer(reqBody),
                )

                duration := time.Since(reqStart)

                mutex.Lock()
                if err != nil || resp.StatusCode != http.StatusOK {
                    errorCount++
                } else {
                    successCount++
                }
                responseTimes = append(responseTimes, duration)
                mutex.Unlock()

                if resp != nil {
                    resp.Body.Close()
                }
            }
        }(i)
    }

    wg.Wait()
    totalDuration := time.Since(start)

    // Calculate metrics
    actualRPS := float64(totalRequests) / totalDuration.Seconds()
    successRate := float64(successCount) / float64(totalRequests)

    // Sort response times for percentiles
    sort.Slice(responseTimes, func(i, j int) bool {
        return responseTimes[i] < responseTimes[j]
    })

    p50 := responseTimes[len(responseTimes)/2]
    p95 := responseTimes[len(responseTimes)*95/100]
    p99 := responseTimes[len(responseTimes)*99/100]

    t.Logf("Load test results:")
    t.Logf("  Total requests: %d", totalRequests)
    t.Logf("  Success rate: %.2f%%", successRate*100)
    t.Logf("  Actual RPS: %.2f", actualRPS)
    t.Logf("  P50 latency: %v", p50)
    t.Logf("  P95 latency: %v", p95)
    t.Logf("  P99 latency: %v", p99)

    // Performance assertions
    assert.Greater(t, actualRPS, float64(targetRPS*0.8), "RPS too low")
    assert.Greater(t, successRate, 0.95, "Success rate too low")
    assert.Less(t, p95, 500*time.Millisecond, "P95 latency too high")
}
```

## Database Development

### Migration Guidelines

#### PostgreSQL Migrations

```sql
-- migrations/postgres/001_create_users_table.up.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone_number VARCHAR(20) NOT NULL UNIQUE,
    country_code CHAR(2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    profile JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Add indexes
CREATE INDEX idx_users_phone_number ON users(phone_number);
CREATE INDEX idx_users_country_code ON users(country_code);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Add triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

```sql
-- migrations/postgres/001_create_users_table.down.sql
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS users;
```

#### ScyllaDB Schema

```sql
-- migrations/scylla/001_messaging_keyspace.cql
CREATE KEYSPACE IF NOT EXISTS tchat_messaging
WITH replication = {
    'class': 'SimpleStrategy',
    'replication_factor': 3
};

USE tchat_messaging;

CREATE TABLE messages (
    dialog_id UUID,
    message_id UUID,
    sender_id UUID,
    content TEXT,
    message_type TEXT,
    metadata MAP<TEXT, TEXT>,
    reply_to_id UUID,
    reactions LIST<FROZEN<reaction>>,
    created_at TIMESTAMP,
    PRIMARY KEY (dialog_id, created_at, message_id)
) WITH CLUSTERING ORDER BY (created_at DESC);

-- User-defined type for reactions
CREATE TYPE reaction (
    emoji TEXT,
    user_ids SET<UUID>
);

-- Create secondary indexes
CREATE INDEX ON messages (sender_id);
CREATE INDEX ON messages (message_type);
```

### Repository Pattern

```go
// repositories/user_repository.go
package repositories

import (
    "context"
    "database/sql"
    "encoding/json"

    "github.com/google/uuid"
    "github.com/lib/pq"

    "github.com/tchat/backend/internal/auth/models"
)

type PostgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
    query := `
        SELECT id, phone_number, country_code, status, profile, created_at, updated_at
        FROM users
        WHERE id = $1
    `

    var user models.User
    var profileJSON []byte

    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Phone,
        &user.CountryCode,
        &user.Status,
        &profileJSON,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }

    if err := json.Unmarshal(profileJSON, &user.Profile); err != nil {
        return nil, err
    }

    return &user, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
    query := `
        INSERT INTO users (id, phone_number, country_code, status, profile)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING created_at, updated_at
    `

    profileJSON, err := json.Marshal(user.Profile)
    if err != nil {
        return err
    }

    return r.db.QueryRowContext(
        ctx, query,
        user.ID, user.Phone, user.CountryCode, user.Status, profileJSON,
    ).Scan(&user.CreatedAt, &user.UpdatedAt)
}
```

## API Development

### Request/Response Models

```go
// models/requests.go
package models

import (
    "github.com/go-playground/validator/v10"
)

type SendOTPRequest struct {
    PhoneNumber string `json:"phone_number" validate:"required,e164"`
    CountryCode string `json:"country_code" validate:"required,iso3166_1_alpha2"`
    Locale      string `json:"locale" validate:"omitempty,locale"`
}

type VerifyOTPRequest struct {
    PhoneNumber string      `json:"phone_number" validate:"required,e164"`
    CountryCode string      `json:"country_code" validate:"required,iso3166_1_alpha2"`
    OTPCode     string      `json:"otp_code" validate:"required,len=6,numeric"`
    DeviceInfo  *DeviceInfo `json:"device_info,omitempty"`
}

type DeviceInfo struct {
    DeviceID   string `json:"device_id" validate:"required,max=100"`
    Platform   string `json:"platform" validate:"required,oneof=web mobile_ios mobile_android"`
    AppVersion string `json:"app_version" validate:"omitempty,semver"`
}

// Custom validator for phone numbers
func ValidateE164(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    // E.164 format validation logic
    return regexp.MustCompile(`^\+[1-9]\d{1,14}$`).MatchString(phone)
}
```

### Handler Implementation

```go
// handlers/auth_handler.go
package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"

    "github.com/tchat/backend/internal/auth/models"
    "github.com/tchat/backend/internal/auth/services"
    "github.com/tchat/backend/shared/middleware"
)

type AuthHandler struct {
    authService services.AuthService
    validator   *validator.Validate
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
    v := validator.New()
    v.RegisterValidation("e164", models.ValidateE164)

    return &AuthHandler{
        authService: authService,
        validator:   v,
    }
}

func (h *AuthHandler) SendOTP(c *gin.Context) {
    var req models.SendOTPRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    if err := h.validator.Struct(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Validation failed",
            "details": err.Error(),
        })
        return
    }

    ctx := c.Request.Context()

    if err := h.authService.SendOTP(ctx, req.PhoneNumber, req.CountryCode); err != nil {
        switch err {
        case services.ErrRateLimitExceeded:
            c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":    "OTP sent successfully",
        "expires_in": 300, // 5 minutes
    })
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
    auth := router.Group("/auth")
    {
        auth.POST("/otp/send", h.SendOTP)
        auth.POST("/otp/verify", h.VerifyOTP)
        auth.POST("/refresh", h.RefreshToken)
        auth.POST("/logout", middleware.RequireAuth(), h.Logout)
    }
}
```

### Middleware Development

```go
// middleware/auth_middleware.go
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"

    "github.com/tchat/backend/shared/auth"
)

func RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
            c.Abort()
            return
        }

        claims, err := auth.ValidateJWT(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Add user info to context
        c.Set("user_id", claims.UserID)
        c.Set("phone", claims.Phone)
        c.Next()
    }
}

func RateLimit(requestsPerMinute int) gin.HandlerFunc {
    // Implementation using Redis or in-memory store
    return func(c *gin.Context) {
        // Rate limiting logic
        c.Next()
    }
}
```

## Contributing

### Git Workflow

1. **Feature Development**
   ```bash
   # Create feature branch
   git checkout -b feature/user-authentication

   # Make changes and commit
   git add .
   git commit -m "feat: implement OTP verification"

   # Push and create PR
   git push origin feature/user-authentication
   ```

2. **Commit Message Format**
   ```
   type(scope): description

   feat: new feature
   fix: bug fix
   docs: documentation only changes
   style: formatting, missing semi colons, etc
   refactor: code change that neither fixes a bug nor adds a feature
   test: adding missing tests
   chore: maintain
   ```

### Code Review Guidelines

#### Checklist for Reviewers

- [ ] Code follows style guide and conventions
- [ ] All tests pass and coverage is maintained
- [ ] Documentation is updated
- [ ] No security vulnerabilities introduced
- [ ] Performance impact is acceptable
- [ ] Database migrations are reversible
- [ ] Error handling is appropriate
- [ ] Logging is adequate for debugging

#### Pull Request Template

```markdown
## Summary
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added for new functionality
```

## Debugging

### Local Debugging

#### VS Code Configuration

```json
// .vscode/launch.json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug Auth Service",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/auth",
            "env": {
                "ENVIRONMENT": "development",
                "LOG_LEVEL": "debug",
                "DB_HOST": "localhost",
                "DB_PORT": "5432"
            },
            "args": []
        }
    ]
}
```

#### Debug Tools

```bash
# Memory profiling
go tool pprof http://localhost:8081/debug/pprof/heap

# CPU profiling
go tool pprof http://localhost:8081/debug/pprof/profile?seconds=30

# Goroutine debugging
go tool pprof http://localhost:8081/debug/pprof/goroutine

# Race condition detection
go run -race cmd/auth/main.go
```

### Production Debugging

#### Log Analysis

```bash
# Search for errors in specific service
kubectl logs -n tchat -l app=auth-service | grep ERROR

# Follow logs in real-time
kubectl logs -n tchat -f deployment/auth-service

# Aggregate logs from all replicas
kubectl logs -n tchat -l app=auth-service --prefix=true
```

#### Performance Debugging

```bash
# Check resource usage
kubectl top pods -n tchat

# Check service endpoints
kubectl get endpoints -n tchat

# Check service health
curl -f http://auth-service:8081/health

# Check metrics
curl http://auth-service:8081/metrics | grep http_requests_total
```

### Common Issues

#### Database Connection Problems

```bash
# Test database connectivity
go run tools/db-test/main.go

# Check connection pool
curl http://localhost:8081/debug/vars | jq '.db_stats'

# PostgreSQL connection check
psql -h localhost -U tchat_user -d tchat_auth -c "SELECT 1;"
```

#### Memory Leaks

```bash
# Monitor memory usage over time
watch kubectl top pod auth-service-xxx -n tchat

# Generate heap dump
curl http://localhost:8081/debug/pprof/heap > heap.prof
go tool pprof heap.prof
```

#### Performance Issues

```bash
# Load testing
go run tests/load/auth_test.go -users=100 -duration=5m

# Trace analysis
curl http://localhost:8081/debug/pprof/trace?seconds=30 > trace.out
go tool trace trace.out
```

## Development Commands

### Make Targets

```makefile
# Development
.PHONY: dev dev-auth dev-messaging dev-payment dev-notification
dev: ## Start all services
dev-auth: ## Start auth service only
dev-messaging: ## Start messaging service only
dev-payment: ## Start payment service only
dev-notification: ## Start notification service only

# Testing
.PHONY: test test-unit test-integration test-performance
test: ## Run all tests
test-unit: ## Run unit tests only
test-integration: ## Run integration tests only
test-performance: ## Run performance tests

# Database
.PHONY: migrate-up migrate-down migrate-force
migrate-up: ## Run all migrations
migrate-down: ## Rollback last migration
migrate-force: ## Force migration version

# Quality
.PHONY: lint fmt check
lint: ## Run linter
fmt: ## Format code
check: ## Run all quality checks

# Build
.PHONY: build build-auth build-messaging build-payment build-notification
build: ## Build all services
build-auth: ## Build auth service
build-messaging: ## Build messaging service
build-payment: ## Build payment service
build-notification: ## Build notification service
```

### Useful Scripts

```bash
# scripts/dev-setup.sh
#!/bin/bash
echo "Setting up development environment..."

# Start infrastructure
docker-compose -f docker-compose.dev.yml up -d

# Wait for services
./scripts/wait-for-services.sh

# Run migrations
make migrate-up

# Install tools
make install-tools

echo "Development environment ready!"
```

This comprehensive development guide provides all the necessary information for developers to contribute effectively to the Tchat backend services.