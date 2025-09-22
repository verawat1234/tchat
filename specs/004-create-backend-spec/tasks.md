# Tasks: Backend Services Architecture

**Input**: Design documents from `/Users/weerawat/Tchat/specs/004-create-backend-spec/`
**Prerequisites**: plan.md (✓), research.md (✓), data-model.md (✓), contracts/ (✓), quickstart.md (✓)

## Execution Flow (main)
```
1. Load plan.md from feature directory
   ✓ Extract: Go 1.21+, gRPC microservices, PostgreSQL/ScyllaDB/Redis
2. Load optional design documents:
   ✓ data-model.md: 12+ entities across 5 services
   ✓ contracts/: 3 OpenAPI specs (auth, messaging, payment)
   ✓ research.md: Technology decisions and patterns
   ✓ quickstart.md: TDD test scenarios
3. Generate tasks by category:
   ✓ Setup: microservices structure, Go modules, Docker
   ✓ Tests: contract tests (3), integration tests (8)
   ✓ Core: models (12), services (5), API endpoints (20+)
   ✓ Integration: databases, message brokers, external services
   ✓ Polish: performance tests, documentation, monitoring
4. Apply task rules:
   ✓ Different services/files = [P] for parallel execution
   ✓ TDD: All tests before implementation
   ✓ Service dependencies respected
5. Number tasks sequentially (T001-T075)
6. Generate dependency graph and parallel execution examples
7. Validate task completeness
8. Return: SUCCESS (75 tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files/services, no dependencies)
- Exact file paths included for each task

## Path Conventions
Based on plan.md structure decision: **Web application (Option 2)**
- **Backend**: `backend/` directory with microservices
- **Services**: Each service in separate directory
- **Shared**: Common libraries and utilities

## Phase 3.1: Project Setup
- [ ] **T001** Create backend microservices directory structure per implementation plan
- [ ] **T002** [P] Initialize auth-service Go module in `backend/auth/`
- [ ] **T003** [P] Initialize messaging-service Go module in `backend/messaging/`
- [ ] **T004** [P] Initialize payment-service Go module in `backend/payment/`
- [ ] **T005** [P] Initialize commerce-service Go module in `backend/commerce/`
- [ ] **T006** [P] Initialize notification-service Go module in `backend/notification/`
- [ ] **T007** [P] Initialize shared utilities Go module in `backend/shared/`
- [ ] **T008** Configure Docker Compose for development stack in `docker-compose.dev.yml`
- [ ] **T009** [P] Setup golangci-lint configuration in `backend/.golangci.yml`
- [ ] **T010** [P] Setup Git hooks and CI/CD pipeline in `.github/workflows/backend.yml`

## Phase 3.2: Contract Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Auth Service Contract Tests
- [ ] **T011** [P] Contract test POST /auth/otp/send in `backend/tests/contract/auth_otp_send_test.go`
- [ ] **T012** [P] Contract test POST /auth/otp/verify in `backend/tests/contract/auth_otp_verify_test.go`
- [ ] **T013** [P] Contract test POST /auth/token/refresh in `backend/tests/contract/auth_token_refresh_test.go`
- [ ] **T014** [P] Contract test GET /users/profile in `backend/tests/contract/auth_profile_test.go`
- [ ] **T015** [P] Contract test PUT /users/profile in `backend/tests/contract/auth_profile_update_test.go`
- [ ] **T016** [P] Contract test POST /users/kyc in `backend/tests/contract/auth_kyc_test.go`

### Messaging Service Contract Tests
- [ ] **T017** [P] Contract test GET /dialogs in `backend/tests/contract/messaging_dialogs_test.go`
- [ ] **T018** [P] Contract test POST /dialogs in `backend/tests/contract/messaging_dialogs_create_test.go`
- [ ] **T019** [P] Contract test GET /dialogs/{id}/messages in `backend/tests/contract/messaging_messages_get_test.go`
- [ ] **T020** [P] Contract test POST /dialogs/{id}/messages in `backend/tests/contract/messaging_messages_send_test.go`
- [ ] **T021** [P] Contract test POST /messages/{id}/read in `backend/tests/contract/messaging_read_test.go`
- [ ] **T022** [P] Contract test WebSocket /websocket in `backend/tests/contract/messaging_websocket_test.go`

### Payment Service Contract Tests
- [ ] **T023** [P] Contract test GET /wallets in `backend/tests/contract/payment_wallets_test.go`
- [ ] **T024** [P] Contract test POST /wallets in `backend/tests/contract/payment_wallets_create_test.go`
- [ ] **T025** [P] Contract test GET /wallets/{id}/balance in `backend/tests/contract/payment_balance_test.go`
- [ ] **T026** [P] Contract test GET /wallets/{id}/transactions in `backend/tests/contract/payment_transactions_test.go`
- [ ] **T027** [P] Contract test POST /transactions/send in `backend/tests/contract/payment_send_test.go`
- [ ] **T028** [P] Contract test POST /transactions/topup in `backend/tests/contract/payment_topup_test.go`

### Integration Test Scenarios
- [ ] **T029** [P] Integration test complete auth flow in `backend/tests/integration/auth_flow_test.go`
- [ ] **T030** [P] Integration test real-time messaging flow in `backend/tests/integration/messaging_flow_test.go`
- [ ] **T031** [P] Integration test payment processing flow in `backend/tests/integration/payment_flow_test.go`
- [ ] **T032** [P] Integration test cross-service user registration in `backend/tests/integration/user_registration_test.go`
- [ ] **T033** [P] Integration test message payment workflow in `backend/tests/integration/message_payment_test.go`
- [ ] **T034** [P] Integration test KYC verification workflow in `backend/tests/integration/kyc_flow_test.go`
- [ ] **T035** [P] Integration test multi-region data sync in `backend/tests/integration/region_sync_test.go`
- [ ] **T036** [P] Integration test external service failures in `backend/tests/integration/service_failures_test.go`

## Phase 3.3: Data Models (ONLY after tests are failing)

### Auth Service Models
- [ ] **T037** [P] User model with validation in `backend/auth/models/user.go`
- [ ] **T038** [P] Session model with state transitions in `backend/auth/models/session.go`
- [ ] **T039** [P] KYC model with document handling in `backend/auth/models/kyc.go`

### Messaging Service Models
- [ ] **T040** [P] Dialog model with participants in `backend/messaging/models/dialog.go`
- [ ] **T041** [P] Message model with content types in `backend/messaging/models/message.go`
- [ ] **T042** [P] Presence model for online status in `backend/messaging/models/presence.go`

### Payment Service Models
- [ ] **T043** [P] Wallet model with multi-currency in `backend/payment/models/wallet.go`
- [ ] **T044** [P] Transaction model with state machine in `backend/payment/models/transaction.go`
- [ ] **T045** [P] PaymentMethod model for provider integration in `backend/payment/models/payment_method.go`

### Commerce Service Models
- [ ] **T046** [P] Product model with variants in `backend/commerce/models/product.go`
- [ ] **T047** [P] Order model with fulfillment in `backend/commerce/models/order.go`
- [ ] **T048** [P] Shop model with verification in `backend/commerce/models/shop.go`

### Notification Service Models
- [ ] **T049** [P] Notification model with multi-channel in `backend/notification/models/notification.go`

### Shared Models
- [ ] **T050** [P] Event sourcing models in `backend/shared/models/event.go`
- [ ] **T051** [P] Saga execution models in `backend/shared/models/saga.go`

## Phase 3.4: Service Layer Implementation

### Auth Service
- [ ] **T052** UserService CRUD operations in `backend/auth/services/user_service.go`
- [ ] **T053** AuthService with OTP generation in `backend/auth/services/auth_service.go`
- [ ] **T054** SessionService with JWT management in `backend/auth/services/session_service.go`
- [ ] **T055** KYCService with verification workflow in `backend/auth/services/kyc_service.go`

### Messaging Service
- [ ] **T056** DialogService with participant management in `backend/messaging/services/dialog_service.go`
- [ ] **T057** MessageService with real-time delivery in `backend/messaging/services/message_service.go`
- [ ] **T058** PresenceService with WebSocket management in `backend/messaging/services/presence_service.go`

### Payment Service
- [ ] **T059** WalletService with balance management in `backend/payment/services/wallet_service.go`
- [ ] **T060** TransactionService with distributed processing in `backend/payment/services/transaction_service.go`
- [ ] **T061** PaymentGatewayService with provider integration in `backend/payment/services/payment_gateway_service.go`

### Notification Service
- [ ] **T062** NotificationService with channel routing in `backend/notification/services/notification_service.go`

## Phase 3.5: API Endpoints Implementation

### Auth Service Endpoints
- [ ] **T063** POST /auth/otp/send endpoint in `backend/auth/handlers/otp.go`
- [ ] **T064** POST /auth/otp/verify endpoint in `backend/auth/handlers/auth.go`
- [ ] **T065** GET /users/profile endpoint in `backend/auth/handlers/profile.go`

### Messaging Service Endpoints
- [ ] **T066** GET /dialogs endpoint in `backend/messaging/handlers/dialogs.go`
- [ ] **T067** POST /dialogs/{id}/messages endpoint in `backend/messaging/handlers/messages.go`
- [ ] **T068** WebSocket /websocket endpoint in `backend/messaging/handlers/websocket.go`

### Payment Service Endpoints
- [ ] **T069** GET /wallets endpoint in `backend/payment/handlers/wallets.go`
- [ ] **T070** POST /transactions/send endpoint in `backend/payment/handlers/transactions.go`

## Phase 3.6: Infrastructure Integration
- [ ] **T071** PostgreSQL database connections and migrations
- [ ] **T072** ScyllaDB setup for message timelines
- [ ] **T073** Redis integration for caching and real-time messaging
- [ ] **T074** Kafka event streaming for cross-service communication
- [ ] **T075** External service integration (SMS, payment gateways, CDN)

## Dependencies

### Setup Dependencies
- T001 must complete before all other tasks
- T002-T007 can run in parallel (different directories)
- T008-T010 depend on T001 completion

### Test Dependencies (TDD Critical Path)
- **ALL contract tests (T011-T028) MUST complete before ANY implementation**
- **ALL integration tests (T029-T036) MUST complete before service implementation**
- Tests are independent and can run in parallel [P]

### Implementation Dependencies
- Data models (T037-T051) can run in parallel within services
- Service layer (T052-T062) depends on corresponding models
- API endpoints (T063-T070) depend on corresponding services
- Infrastructure (T071-T075) can run parallel with implementation

### Service Dependencies
- Auth Service models/services must complete before other services (user management)
- Messaging Service depends on Auth Service (user authentication)
- Payment Service depends on Auth Service (user verification)
- Notification Service depends on all other services (event handling)

## Parallel Execution Examples

### Phase 3.1 - Project Setup (6 parallel tasks)
```bash
# Launch T002-T007 together:
Task: "Initialize auth-service Go module in backend/auth/"
Task: "Initialize messaging-service Go module in backend/messaging/"
Task: "Initialize payment-service Go module in backend/payment/"
Task: "Initialize commerce-service Go module in backend/commerce/"
Task: "Initialize notification-service Go module in backend/notification/"
Task: "Initialize shared utilities Go module in backend/shared/"
```

### Phase 3.2 - Contract Tests (18 parallel tasks)
```bash
# Launch T011-T028 together (Auth, Messaging, Payment contract tests):
Task: "Contract test POST /auth/otp/send in backend/tests/contract/auth_otp_send_test.go"
Task: "Contract test POST /auth/otp/verify in backend/tests/contract/auth_otp_verify_test.go"
Task: "Contract test GET /dialogs in backend/tests/contract/messaging_dialogs_test.go"
Task: "Contract test POST /transactions/send in backend/tests/contract/payment_send_test.go"
# ... (all 18 contract tests can run in parallel)

# Launch T029-T036 together (Integration tests):
Task: "Integration test complete auth flow in backend/tests/integration/auth_flow_test.go"
Task: "Integration test real-time messaging flow in backend/tests/integration/messaging_flow_test.go"
Task: "Integration test payment processing flow in backend/tests/integration/payment_flow_test.go"
# ... (all 8 integration tests can run in parallel)
```

### Phase 3.3 - Data Models (15 parallel tasks)
```bash
# Launch T037-T051 together (all models across services):
Task: "User model with validation in backend/auth/models/user.go"
Task: "Dialog model with participants in backend/messaging/models/dialog.go"
Task: "Wallet model with multi-currency in backend/payment/models/wallet.go"
Task: "Product model with variants in backend/commerce/models/product.go"
Task: "Notification model with multi-channel in backend/notification/models/notification.go"
# ... (all models can run in parallel as they're in different files)
```

## Performance Requirements Validation
- Message delivery latency: <200ms (validated in T030)
- API response times: <100ms 95th percentile (validated in contract tests)
- Payment throughput: 1000+ TPS (validated in T031)
- Concurrent users: 100K+ WebSocket connections (validated in T022)

## Security & Compliance Validation
- JWT authentication: 15-minute expiry (T013, T014)
- Data encryption: AES-256 for sensitive data (T037, T043)
- PDPA compliance: Data export/deletion (T032, T034)
- PCI DSS: Payment tokenization (T061, T070)

## Notes
- **[P] tasks** = different files/services, no dependencies - can run in parallel
- **TDD Critical**: All tests (T011-T036) MUST fail before implementation starts
- **Service Independence**: Each microservice can be developed in parallel after models
- **External Dependencies**: Infrastructure tasks (T071-T075) require external services
- **Commit Strategy**: Commit after each task completion for incremental progress
- **Rollback Safety**: Each task is atomic and can be rolled back independently

## Task Generation Rules Applied
1. **From Contracts**: 3 OpenAPI files → 18 contract test tasks [P]
2. **From Data Model**: 12+ entities → 15 model creation tasks [P]
3. **From Quickstart**: 8 integration scenarios → 8 integration tests [P]
4. **Service Architecture**: 5 microservices → parallel development streams
5. **TDD Ordering**: All tests before any implementation (critical requirement)

## Validation Checklist
- [x] All contracts have corresponding tests (T011-T028)
- [x] All entities have model tasks (T037-T051)
- [x] All tests come before implementation (T011-T036 before T037+)
- [x] Parallel tasks truly independent (different files/services)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] Performance and compliance requirements addressed
- [x] Microservices architecture properly decomposed
- [x] External service integrations planned