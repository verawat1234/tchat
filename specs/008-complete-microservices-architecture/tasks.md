# Tasks: Enterprise-Grade Tchat Backend Platform

**Input**: Design documents from `/specs/008-complete-microservices-architecture/`
**Prerequisites**: plan.md, research.md, data-model.md, contracts/, quickstart.md

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Extract: Go 1.22+ microservices, Gin HTTP, PostgreSQL/ScyllaDB/Redis/Kafka
   → Structure: Web application (backend microservices + frontend integration)
2. Load design documents:
   → data-model.md: 10 entities (User, Business, Message, Dialog, Transaction, Wallet, Product, Order, Notification, Session)
   → contracts/: auth-service.yaml, messaging-service.yaml
   → research.md: Technology stack decisions and regional compliance
   → quickstart.md: End-to-end test scenarios for SEA markets
3. Generate tasks by category:
   → Setup: Go microservices project, dependencies, database setup
   → Tests: Contract tests for APIs, integration tests for user scenarios
   → Core: Entity models, service implementations, API endpoints
   → Integration: Database connections, message queues, authentication
   → Polish: Performance testing, documentation, deployment
4. Apply task rules:
   → Different microservices = mark [P] for parallel
   → Shared database/config = sequential
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001-T080)
6. Generate dependency graph for microservices architecture
7. Create parallel execution examples for independent services
8. Validate: All contracts tested, entities modeled, endpoints implemented
9. Return: SUCCESS (80 tasks ready for enterprise platform execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different services/files, no dependencies)
- Include exact file paths for microservices architecture

## Path Conventions
**Backend microservices structure** (from plan.md):
```
backend/
├── services/
│   ├── auth-service/
│   ├── messaging-service/
│   ├── payment-service/
│   ├── commerce-service/
│   └── notification-service/
├── shared/
│   ├── models/
│   ├── middleware/
│   └── utils/
├── infrastructure/
├── scripts/
└── tests/
```

## Phase 3.1: Foundation Setup

### Project Structure & Dependencies
- [ ] T001 Create Go microservices project structure per plan.md
- [ ] T002 Initialize go.mod with Go 1.22+ and shared dependencies
- [ ] T003 [P] Configure air for hot reload development environment
- [ ] T004 [P] Setup golangci-lint configuration for code quality
- [ ] T005 [P] Create Docker development environment with docker-compose.dev.yml
- [ ] T006 [P] Setup PostgreSQL migration system using golang-migrate
- [ ] T007 [P] Setup ScyllaDB schema for message storage
- [ ] T008 [P] Configure Redis for session and cache management
- [ ] T009 [P] Setup Kafka for event streaming between services

### Environment Configuration
- [ ] T010 [P] Create environment configuration system in shared/config/
- [ ] T011 [P] Setup JWT authentication middleware in shared/middleware/
- [ ] T012 [P] Create database connection managers in shared/db/
- [ ] T013 [P] Setup structured logging with regional compliance in shared/logger/

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3

**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Auth Service Contract Tests
- [ ] T014 [P] Contract test POST /api/v1/auth/otp/send in tests/contract/auth/test_otp_send.go
- [ ] T015 [P] Contract test POST /api/v1/auth/otp/verify in tests/contract/auth/test_otp_verify.go
- [ ] T016 [P] Contract test POST /api/v1/auth/refresh in tests/contract/auth/test_refresh_token.go
- [ ] T017 [P] Contract test GET /api/v1/users/profile in tests/contract/auth/test_user_profile.go

### Messaging Service Contract Tests
- [ ] T018 [P] Contract test GET /api/v1/dialogs in tests/contract/messaging/test_dialogs_list.go
- [ ] T019 [P] Contract test POST /api/v1/dialogs in tests/contract/messaging/test_dialogs_create.go
- [ ] T020 [P] Contract test GET /api/v1/dialogs/{id}/messages in tests/contract/messaging/test_messages_get.go
- [ ] T021 [P] Contract test POST /api/v1/dialogs/{id}/messages in tests/contract/messaging/test_messages_send.go

### Integration Tests for User Scenarios
- [ ] T022 [P] Integration test: Thai user registration with OTP in tests/integration/test_thai_user_flow.go
- [ ] T023 [P] Integration test: Cross-country messaging (TH to SG) in tests/integration/test_cross_country_messaging.go
- [ ] T024 [P] Integration test: Multi-currency wallet operations in tests/integration/test_multicurrency_wallets.go
- [ ] T025 [P] Integration test: E-commerce shop creation flow in tests/integration/test_ecommerce_flow.go
- [ ] T026 [P] Integration test: Real-time WebSocket messaging in tests/integration/test_websocket_messaging.go

## Phase 3.3: Core Models & Data Layer (ONLY after tests are failing)

### Entity Models (from data-model.md)
- [ ] T027 [P] User entity model in shared/models/user.go
- [ ] T028 [P] Business entity model in shared/models/business.go
- [ ] T029 [P] Message entity model in shared/models/message.go
- [ ] T030 [P] Dialog entity model in shared/models/dialog.go
- [ ] T031 [P] Transaction entity model in shared/models/transaction.go
- [ ] T032 [P] Wallet entity model in shared/models/wallet.go
- [ ] T033 [P] Product entity model in shared/models/product.go
- [ ] T034 [P] Order entity model in shared/models/order.go
- [ ] T035 [P] Notification entity model in shared/models/notification.go
- [ ] T036 [P] Session entity model in shared/models/session.go

### Database Repositories
- [ ] T037 [P] User repository with PostgreSQL GORM in shared/repositories/user_repo.go
- [ ] T038 [P] Message repository with ScyllaDB in shared/repositories/message_repo.go
- [ ] T039 [P] Dialog repository with PostgreSQL in shared/repositories/dialog_repo.go
- [ ] T040 [P] Transaction repository with PostgreSQL in shared/repositories/transaction_repo.go
- [ ] T041 [P] Wallet repository with PostgreSQL in shared/repositories/wallet_repo.go

## Phase 3.4: Microservice Implementation

### Auth Service Implementation
- [ ] T042 OTP sending service with regional SMS providers in services/auth-service/internal/services/otp_service.go
- [ ] T043 OTP verification and JWT generation in services/auth-service/internal/services/auth_service.go
- [ ] T044 User profile management service in services/auth-service/internal/services/profile_service.go
- [ ] T045 POST /api/v1/auth/otp/send handler in services/auth-service/internal/handlers/otp_handler.go
- [ ] T046 POST /api/v1/auth/otp/verify handler in services/auth-service/internal/handlers/auth_handler.go
- [ ] T047 User profile endpoints in services/auth-service/internal/handlers/profile_handler.go
- [ ] T048 Auth service main server in services/auth-service/main.go

### Messaging Service Implementation
- [ ] T049 Real-time WebSocket connection manager in services/messaging-service/internal/services/websocket_service.go
- [ ] T050 Message encryption/decryption service in services/messaging-service/internal/services/crypto_service.go
- [ ] T051 Dialog management service in services/messaging-service/internal/services/dialog_service.go
- [ ] T052 Message delivery tracking service in services/messaging-service/internal/services/delivery_service.go
- [ ] T053 GET /api/v1/dialogs handler in services/messaging-service/internal/handlers/dialog_handler.go
- [ ] T054 Message CRUD handlers in services/messaging-service/internal/handlers/message_handler.go
- [ ] T055 WebSocket upgrade and message routing in services/messaging-service/internal/handlers/websocket_handler.go
- [ ] T056 Messaging service main server in services/messaging-service/main.go

### Payment Service Implementation
- [ ] T057 [P] Multi-currency wallet service in services/payment-service/internal/services/wallet_service.go
- [ ] T058 [P] Transaction processing with regional gateways in services/payment-service/internal/services/payment_service.go
- [ ] T059 [P] Currency conversion service in services/payment-service/internal/services/currency_service.go
- [ ] T060 [P] Payment service main server in services/payment-service/main.go

### Commerce Service Implementation
- [ ] T061 [P] Business verification service in services/commerce-service/internal/services/business_service.go
- [ ] T062 [P] Product catalog management in services/commerce-service/internal/services/product_service.go
- [ ] T063 [P] Order processing service in services/commerce-service/internal/services/order_service.go
- [ ] T064 [P] Commerce service main server in services/commerce-service/main.go

### Notification Service Implementation
- [ ] T065 [P] Multi-channel notification dispatch in services/notification-service/internal/services/notification_service.go
- [ ] T066 [P] Template management for regional languages in services/notification-service/internal/services/template_service.go
- [ ] T067 [P] Notification service main server in services/notification-service/main.go

## Phase 3.5: Integration & Infrastructure

### API Gateway & Load Balancing
- [ ] T068 API Gateway with service discovery in infrastructure/gateway/main.go
- [ ] T069 Rate limiting middleware for regional compliance in shared/middleware/rate_limiter.go
- [ ] T070 Health check endpoints for all services in shared/middleware/health.go

### Event-Driven Architecture
- [ ] T071 Kafka event producers for service communication in shared/events/producers.go
- [ ] T072 Kafka event consumers for cross-service notifications in shared/events/consumers.go
- [ ] T073 Event schema registry for message versioning in shared/events/schemas.go

### Security & Compliance
- [ ] T074 Regional data residency enforcement in shared/middleware/data_residency.go
- [ ] T075 Audit logging for compliance in shared/middleware/audit.go
- [ ] T076 Input validation for Southeast Asian formats in shared/validators/

## Phase 3.6: Polish & Deployment

### Performance & Testing
- [ ] T077 [P] Performance benchmarks achieving target metrics in tests/performance/
- [ ] T078 [P] End-to-end test execution per quickstart.md scenarios in tests/e2e/
- [ ] T079 [P] Documentation update with API examples in docs/
- [ ] T080 [P] Production deployment configuration with Docker and Kubernetes in infrastructure/deploy/

## Dependencies

### Critical Path Dependencies
- Foundation (T001-T013) → All subsequent phases
- Tests (T014-T026) → Implementation (T027-T067)
- Models (T027-T036) → Repositories (T037-T041) → Services (T042-T067)
- Auth Service (T042-T048) → All other services (user context required)
- Core services → Integration (T068-T076) → Polish (T077-T080)

### Service Dependencies
- T048 (Auth Service) blocks all user-dependent services
- T038 (Message Repository) blocks T049-T056 (Messaging Service)
- T037 (User Repository) blocks T042-T048 (Auth Service)
- T071-T073 (Event System) required before T068 (API Gateway)

## Parallel Execution Examples

### Phase 3.2: Contract Tests (All Parallel)
```bash
# Launch T014-T021 together (different contract files):
Task: "Contract test POST /api/v1/auth/otp/send in tests/contract/auth/test_otp_send.go"
Task: "Contract test POST /api/v1/auth/otp/verify in tests/contract/auth/test_otp_verify.go"
Task: "Contract test GET /api/v1/dialogs in tests/contract/messaging/test_dialogs_list.go"
Task: "Contract test POST /api/v1/dialogs/{id}/messages in tests/contract/messaging/test_messages_send.go"

# Launch T022-T026 together (different integration scenarios):
Task: "Integration test: Thai user registration with OTP in tests/integration/test_thai_user_flow.go"
Task: "Integration test: Cross-country messaging in tests/integration/test_cross_country_messaging.go"
Task: "Integration test: Multi-currency wallet operations in tests/integration/test_multicurrency_wallets.go"
```

### Phase 3.3: Entity Models (All Parallel)
```bash
# Launch T027-T036 together (different model files):
Task: "User entity model in shared/models/user.go"
Task: "Business entity model in shared/models/business.go"
Task: "Message entity model in shared/models/message.go"
Task: "Transaction entity model in shared/models/transaction.go"
```

### Phase 3.4: Independent Services (Parallel by Service)
```bash
# After T048 (Auth Service) completes, launch T057-T067 together:
Task: "Multi-currency wallet service in services/payment-service/internal/services/wallet_service.go"
Task: "Business verification service in services/commerce-service/internal/services/business_service.go"
Task: "Multi-channel notification dispatch in services/notification-service/internal/services/notification_service.go"
```

## Regional Compliance Requirements

### Southeast Asian Market Specific Tasks
- **Phone Validation**: Support TH, SG, ID, MY, PH, VN formats in auth service
- **Currency Support**: THB, SGD, IDR, MYR, PHP, VND, USD in payment service
- **Language Support**: en, th, id, ms, fil, vi in notification templates
- **Data Residency**: Regional database placement enforcement
- **Regulatory Compliance**: KYC/AML requirements per country

### Performance Targets (from plan.md)
- **Authentication**: 1000+ RPS sustained load
- **Messaging**: 10K concurrent WebSocket connections
- **API Latency**: <200ms P95 response time
- **Message Delivery**: <2s end-to-end latency
- **Uptime**: 99.9% availability target

## Validation Checklist
*GATE: Checked before task execution*

- [x] All auth service contracts have corresponding tests (T014-T017)
- [x] All messaging service contracts have corresponding tests (T018-T021)
- [x] All 10 entities have model tasks (T027-T036)
- [x] All critical user scenarios have integration tests (T022-T026)
- [x] All tests come before implementation (Phase 3.2 → 3.3)
- [x] Parallel tasks are truly independent (different files/services)
- [x] Each task specifies exact file path in microservices structure
- [x] No task modifies same file as another [P] task
- [x] Regional compliance requirements addressed in relevant tasks
- [x] Performance targets integrated into testing and benchmarking tasks

## Notes
- [P] tasks = different services/files, no dependencies
- Verify tests fail before implementing (TDD approach)
- Commit after each task completion
- Focus on Southeast Asian market requirements throughout
- Maintain 99.9% uptime and <200ms latency targets
- All services must support regional compliance and multi-currency operations