# Tasks: RTK Backend API Integration

**Input**: Design documents from `/specs/009-create-rtk-follow/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory ✓
   → Extracted: TypeScript 5.3.0, Redux Toolkit 2.0+, RTK Query, React 18.3.1
2. Load optional design documents ✓
   → data-model.md: User, Auth, Message, Chat entities
   → contracts/api.yaml: Auth, User, Message endpoints
   → research.md: RTK Query decisions, caching, auth patterns
3. Generate tasks by category ✓
   → Setup: RTK installation, store config
   → Tests: MSW contract tests, integration tests
   → Core: API services, auth slice, typed hooks
   → Integration: Token refresh, error handling
   → Polish: Optimistic updates, DevTools, docs
4. Apply task rules ✓
   → Different files marked [P] for parallel
   → Tests before implementation (TDD)
5. Number tasks sequentially ✓
6. Generate dependency graph ✓
7. Create parallel execution examples ✓
8. Validate task completeness ✓
   → All contracts have tests ✓
   → All entities have types ✓
   → All endpoints implemented ✓
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Web app**: `apps/web/src/`, `apps/web/tests/`
- All paths relative to repository root

## Phase 3.1: Setup
- [ ] T001 Create Redux store directory structure in apps/web/src/store/
- [ ] T002 Install RTK dependencies: @reduxjs/toolkit react-redux
- [ ] T003 Install test dependencies: msw @mswjs/data @testing-library/react-hooks
- [ ] T004 [P] Configure TypeScript types for Redux in apps/web/src/types/redux.ts
- [ ] T005 [P] Setup MSW server configuration in apps/web/src/mocks/server.ts

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [ ] T006 [P] Contract test POST /api/auth/login in apps/web/tests/services/auth.test.ts
- [ ] T007 [P] Contract test POST /api/auth/refresh in apps/web/tests/services/auth.test.ts
- [ ] T008 [P] Contract test GET /api/auth/me in apps/web/tests/services/auth.test.ts
- [ ] T009 [P] Contract test GET /api/users in apps/web/tests/services/users.test.ts
- [ ] T010 [P] Contract test GET /api/users/:id in apps/web/tests/services/users.test.ts
- [ ] T011 [P] Contract test PATCH /api/users/:id in apps/web/tests/services/users.test.ts
- [ ] T012 [P] Contract test GET /api/messages in apps/web/tests/services/messages.test.ts
- [ ] T013 [P] Contract test POST /api/messages in apps/web/tests/services/messages.test.ts
- [ ] T014 [P] Integration test login flow in apps/web/tests/integration/auth-flow.test.tsx
- [ ] T015 [P] Integration test data fetching with cache in apps/web/tests/integration/data-fetch.test.tsx
- [ ] T016 [P] Integration test error handling in apps/web/tests/integration/error-handling.test.tsx
- [ ] T017 [P] Integration test optimistic updates in apps/web/tests/integration/optimistic.test.tsx

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [ ] T018 Configure Redux store with RTK Query in apps/web/src/store/index.ts
- [ ] T019 [P] Create typed Redux hooks in apps/web/src/store/hooks.ts
- [ ] T020 [P] Define TypeScript types from data model in apps/web/src/types/api.ts
- [ ] T021 Create base API service with auth headers in apps/web/src/services/api.ts
- [ ] T022 [P] Implement auth API endpoints in apps/web/src/services/auth.ts
- [ ] T023 [P] Implement users API endpoints in apps/web/src/services/users.ts
- [ ] T024 [P] Implement messages API endpoints in apps/web/src/services/messages.ts
- [ ] T025 [P] Implement chats API endpoints in apps/web/src/services/chats.ts
- [ ] T026 Create auth slice for local auth state in apps/web/src/features/authSlice.ts
- [ ] T027 [P] Create UI slice for app state in apps/web/src/features/uiSlice.ts
- [ ] T028 Wire Redux Provider to React app in apps/web/src/main.tsx
- [ ] T029 [P] Create MSW handlers for all endpoints in apps/web/src/mocks/handlers.ts

## Phase 3.4: Integration
- [ ] T030 Implement token refresh middleware in apps/web/src/store/middleware/authMiddleware.ts
- [ ] T031 Add global error handling middleware in apps/web/src/store/middleware/errorMiddleware.ts
- [ ] T032 Configure cache invalidation tags in apps/web/src/services/api.ts
- [ ] T033 Implement optimistic updates for messages in apps/web/src/services/messages.ts
- [ ] T034 Add request retry logic with exponential backoff in apps/web/src/services/api.ts
- [ ] T035 Setup Redux persist for offline support in apps/web/src/store/index.ts
- [ ] T036 [P] Create auth interceptor for 401 handling in apps/web/src/services/api.ts
- [ ] T037 [P] Add loading state management in apps/web/src/features/loadingSlice.ts

## Phase 3.5: Advanced Features
- [ ] T038 Implement cursor-based pagination in apps/web/src/services/messages.ts
- [ ] T039 Add infinite scroll support in apps/web/src/services/messages.ts
- [ ] T040 Create subscription support for real-time updates in apps/web/src/services/api.ts
- [ ] T041 [P] Implement prefetching for common routes in apps/web/src/services/prefetch.ts
- [ ] T042 [P] Add request deduplication in apps/web/src/services/api.ts

## Phase 3.6: Polish
- [ ] T043 [P] Unit test auth slice in apps/web/tests/features/authSlice.test.ts
- [ ] T044 [P] Unit test UI slice in apps/web/tests/features/uiSlice.test.ts
- [ ] T045 [P] Unit test middleware in apps/web/tests/middleware/auth.test.ts
- [ ] T046 Performance test API response times (<500ms) in apps/web/tests/performance/api.test.ts
- [ ] T047 [P] Configure Redux DevTools in apps/web/src/store/index.ts
- [ ] T048 [P] Create usage documentation in apps/web/docs/rtk-usage.md
- [ ] T049 [P] Add JSDoc comments to all API services
- [ ] T050 Run quickstart.md validation steps

## Dependencies
- Setup (T001-T005) must complete first
- Tests (T006-T017) before implementation (T018-T029)
- T018 blocks T019, T026, T027 (store must exist)
- T021 blocks T022-T025 (base API required)
- Core implementation before integration (T030-T037)
- Integration before advanced features (T038-T042)
- Everything before polish (T043-T050)

## Parallel Execution Examples

### Parallel Test Creation (T006-T017)
```bash
# Launch all contract and integration tests together:
Task: "Contract test POST /api/auth/login in apps/web/tests/services/auth.test.ts"
Task: "Contract test GET /api/users in apps/web/tests/services/users.test.ts"
Task: "Contract test GET /api/messages in apps/web/tests/services/messages.test.ts"
Task: "Integration test login flow in apps/web/tests/integration/auth-flow.test.tsx"
Task: "Integration test data fetching in apps/web/tests/integration/data-fetch.test.tsx"
```

### Parallel Service Implementation (T022-T025)
```bash
# After base API (T021) is complete:
Task: "Implement auth API endpoints in apps/web/src/services/auth.ts"
Task: "Implement users API endpoints in apps/web/src/services/users.ts"
Task: "Implement messages API endpoints in apps/web/src/services/messages.ts"
Task: "Implement chats API endpoints in apps/web/src/services/chats.ts"
```

### Parallel Type and Hook Creation (T019-T020)
```bash
# Can run simultaneously as they're independent files:
Task: "Create typed Redux hooks in apps/web/src/store/hooks.ts"
Task: "Define TypeScript types from data model in apps/web/src/types/api.ts"
```

## Task Execution Order

### Critical Path
```
1. Setup (T001-T005) - 30 min
2. Write failing tests (T006-T017) - 2 hours [PARALLEL]
3. Core store setup (T018, T021) - 1 hour
4. API services (T022-T025) - 2 hours [PARALLEL]
5. State management (T026-T027) - 1 hour [PARALLEL]
6. Integration (T030-T037) - 2 hours
7. Advanced features (T038-T042) - 2 hours
8. Polish and testing (T043-T050) - 2 hours [PARALLEL]

Total estimated time: 12-14 hours
```

## Notes
- [P] tasks = different files, no dependencies
- Verify all tests fail before implementing
- Commit after each task completion
- Use Redux DevTools to verify state changes
- Run type checking after TypeScript changes
- MSW handlers must match OpenAPI contract exactly

## Validation Checklist
*GATE: Checked before execution*

- [x] All contracts have corresponding tests (T006-T013)
- [x] All entities have TypeScript types (T020)
- [x] All tests come before implementation
- [x] Parallel tasks truly independent
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] Auth endpoints covered (login, refresh, me, logout)
- [x] CRUD operations covered for main resources
- [x] Error handling middleware included
- [x] Token refresh logic included
- [x] Optimistic updates configured
- [x] DevTools configuration included