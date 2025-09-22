# Tasks: Dynamic Content Management System

**Input**: Design documents from `/Users/weerawat/Tchat/specs/010-now-i-think/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Tech stack: TypeScript 5.3.0, React 18.3.1, RTK 2.0+, RTK Query
   → Structure: Web app (frontend RTK integration)
2. Load optional design documents:
   → data-model.md: ContentItem, ContentCategory, ContentVersion, ContentValue entities
   → contracts/: content-api.ts with 12 endpoints, content-api.test.ts contract tests
   → research.md: RTK Query patterns, caching strategy, fallback system
3. Generate tasks by category:
   → Setup: Content types, RTK slice, hooks
   → Tests: Contract tests for all endpoints, integration tests
   → Core: Content slice, API endpoints, hooks, selectors
   → Integration: Component updates, fallback system
   → Polish: E2E tests, performance optimization, validation
4. Apply task rules:
   → Different files = mark [P] for parallel
   → RTK slice before API endpoints, API before UI hooks
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001-T055)
6. Generate dependency graph
7. Create parallel execution examples
8. SUCCESS (55 tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Web app**: `apps/web/src/` for frontend, existing RTK infrastructure
- Paths integrate with existing 009-create-rtk-follow structure

## Phase 3.1: Setup & Foundation
- [ ] T001 Create content types and interfaces in `apps/web/src/types/content.ts`
- [ ] T002 [P] Create content slice initial structure in `apps/web/src/features/contentSlice.ts`
- [ ] T003 [P] Create content selectors in `apps/web/src/features/contentSelectors.ts`

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Contract Tests [P] - All can run in parallel
- [ ] T004 [P] Contract test getContentItems in `apps/web/src/services/__tests__/content-api-get-items.test.ts`
- [ ] T005 [P] Contract test getContentItem in `apps/web/src/services/__tests__/content-api-get-item.test.ts`
- [ ] T006 [P] Contract test getContentByCategory in `apps/web/src/services/__tests__/content-api-get-category.test.ts`
- [ ] T007 [P] Contract test getContentCategories in `apps/web/src/services/__tests__/content-api-get-categories.test.ts`
- [ ] T008 [P] Contract test getContentVersions in `apps/web/src/services/__tests__/content-api-get-versions.test.ts`
- [ ] T009 [P] Contract test syncContent in `apps/web/src/services/__tests__/content-api-sync.test.ts`
- [ ] T010 [P] Contract test createContentItem in `apps/web/src/services/__tests__/content-api-create.test.ts`
- [ ] T011 [P] Contract test updateContentItem in `apps/web/src/services/__tests__/content-api-update.test.ts`
- [ ] T012 [P] Contract test publishContent in `apps/web/src/services/__tests__/content-api-publish.test.ts`
- [ ] T013 [P] Contract test archiveContent in `apps/web/src/services/__tests__/content-api-archive.test.ts`
- [ ] T014 [P] Contract test bulkUpdateContent in `apps/web/src/services/__tests__/content-api-bulk.test.ts`
- [ ] T015 [P] Contract test revertContentVersion in `apps/web/src/services/__tests__/content-api-revert.test.ts`

### Integration Tests [P] - All can run in parallel
- [ ] T016 [P] Integration test dynamic content loading in `apps/web/src/__tests__/integration/content-loading.test.tsx`
- [ ] T017 [P] Integration test real-time content updates in `apps/web/src/__tests__/integration/content-updates.test.tsx`
- [ ] T018 [P] Integration test fallback content display in `apps/web/src/__tests__/integration/content-fallback.test.tsx`
- [ ] T019 [P] Integration test content loading performance in `apps/web/src/__tests__/integration/content-performance.test.tsx`

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Content Types and Models
- [ ] T020 [P] Implement ContentItem type in `apps/web/src/types/content.ts`
- [ ] T021 [P] Implement ContentCategory type in `apps/web/src/types/content.ts`
- [ ] T022 [P] Implement ContentValue types in `apps/web/src/types/content.ts`
- [ ] T023 [P] Implement ContentVersion type in `apps/web/src/types/content.ts`
- [ ] T024 [P] Implement ContentMetadata type in `apps/web/src/types/content.ts`

### Content Slice Implementation
- [ ] T025 Implement content slice initial state in `apps/web/src/features/contentSlice.ts`
- [ ] T026 Add content slice reducers in `apps/web/src/features/contentSlice.ts`
- [ ] T027 Add content slice extraReducers for API integration in `apps/web/src/features/contentSlice.ts`
- [ ] T028 Implement content selectors in `apps/web/src/features/contentSelectors.ts`

### RTK Query API Endpoints
- [ ] T029 Inject getContentItems endpoint into API in `apps/web/src/services/contentApi.ts`
- [ ] T030 Inject getContentItem endpoint into API in `apps/web/src/services/contentApi.ts`
- [ ] T031 Inject getContentByCategory endpoint into API in `apps/web/src/services/contentApi.ts`
- [ ] T032 Inject getContentCategories endpoint into API in `apps/web/src/services/contentApi.ts`
- [ ] T033 Inject getContentVersions endpoint into API in `apps/web/src/services/contentApi.ts`
- [ ] T034 Inject syncContent endpoint into API in `apps/web/src/services/contentApi.ts`
- [ ] T035 Inject createContentItem mutation into API in `apps/web/src/services/contentApi.ts`
- [ ] T036 Inject updateContentItem mutation into API in `apps/web/src/services/contentApi.ts`
- [ ] T037 Inject publishContent mutation into API in `apps/web/src/services/contentApi.ts`
- [ ] T038 Inject archiveContent mutation into API in `apps/web/src/services/contentApi.ts`
- [ ] T039 Inject bulkUpdateContent mutation into API in `apps/web/src/services/contentApi.ts`
- [ ] T040 Inject revertContentVersion mutation into API in `apps/web/src/services/contentApi.ts`

### Content Hooks and Utilities
- [ ] T041 [P] Create useContent hook in `apps/web/src/hooks/useContent.ts`
- [ ] T042 [P] Create useContentByCategory hook in `apps/web/src/hooks/useContentByCategory.ts`
- [ ] T043 [P] Create useContentFallback hook in `apps/web/src/hooks/useContentFallback.ts`
- [ ] T044 [P] Create content validation utilities in `apps/web/src/utils/contentValidation.ts`
- [ ] T045 [P] Create content transformation utilities in `apps/web/src/utils/contentTransform.ts`

## Phase 3.4: Integration & UI Updates

### Component Content Integration
- [ ] T046 Update App.tsx to include content slice in store in `apps/web/src/App.tsx`
- [ ] T047 Replace hardcoded navigation text in `apps/web/src/components/Navigation.tsx`
- [ ] T048 Replace hardcoded header content in `apps/web/src/components/Header.tsx`
- [ ] T049 Replace hardcoded form labels in sign-in components in `apps/web/src/components/Auth/`
- [ ] T050 Replace hardcoded error messages in error components in `apps/web/src/components/Error/`

### Fallback System Implementation
- [ ] T051 Implement localStorage fallback in `apps/web/src/services/contentFallback.ts`
- [ ] T052 Add error boundary for content failures in `apps/web/src/components/ContentErrorBoundary.tsx`
- [ ] T053 Implement content loading indicators in `apps/web/src/components/ContentLoader.tsx`

## Phase 3.5: Polish & Validation

### Performance & Optimization
- [ ] T054 [P] Add content caching optimization in `apps/web/src/features/contentSlice.ts`
- [ ] T055 [P] Implement content prefetching strategy in `apps/web/src/hooks/useContentPrefetch.ts`

### E2E Testing with Playwright
- [ ] T056 [P] E2E test dynamic content loading flow in `apps/web/tests/e2e/content-loading.spec.ts`
- [ ] T057 [P] E2E test content update propagation in `apps/web/tests/e2e/content-updates.spec.ts`
- [ ] T058 [P] E2E test fallback content behavior in `apps/web/tests/e2e/content-fallback.spec.ts`
- [ ] T059 [P] E2E test content performance requirements in `apps/web/tests/e2e/content-performance.spec.ts`

### Final Validation
- [ ] T060 Run quickstart validation scenarios in `specs/010-now-i-think/quickstart.md`
- [ ] T061 Performance validation: ensure <200ms content load times
- [ ] T062 Verify no hardcoded content remains in production build
- [ ] T063 Test fallback system under network failure conditions

## Dependencies

### Critical Path Dependencies
- **Types** (T020-T024) → **Slice** (T025-T028) → **API** (T029-T040) → **Hooks** (T041-T045) → **Components** (T046-T050)
- **Setup** (T001-T003) before all other phases
- **All tests** (T004-T019) before implementation (T020-T063)
- **Content slice** (T025-T028) before API endpoints (T029-T040)
- **API endpoints** (T029-T040) before hooks (T041-T045)
- **Hooks** (T041-T045) before component integration (T046-T050)

### Blocking Relationships
- T025 blocks T026, T027
- T028 blocks T041, T042, T043
- T029-T040 (API endpoints) block T041-T045 (hooks)
- T041-T045 (hooks) block T046-T050 (component updates)
- T046 blocks T047-T050 (component updates must start with store integration)

## Parallel Execution Examples

### Phase 1: Foundation Setup [P]
```bash
# T001-T003 can run in parallel
Task: "Create content types and interfaces in apps/web/src/types/content.ts"
Task: "Create content slice initial structure in apps/web/src/features/contentSlice.ts"
Task: "Create content selectors in apps/web/src/features/contentSelectors.ts"
```

### Phase 2: Contract Tests [P]
```bash
# T004-T015 can all run in parallel (different test files)
Task: "Contract test getContentItems in apps/web/src/services/__tests__/content-api-get-items.test.ts"
Task: "Contract test getContentItem in apps/web/src/services/__tests__/content-api-get-item.test.ts"
Task: "Contract test getContentByCategory in apps/web/src/services/__tests__/content-api-get-category.test.ts"
Task: "Contract test createContentItem in apps/web/src/services/__tests__/content-api-create.test.ts"
# ... (all 12 contract tests)
```

### Phase 3: Integration Tests [P]
```bash
# T016-T019 can all run in parallel (different integration test files)
Task: "Integration test dynamic content loading in apps/web/src/__tests__/integration/content-loading.test.tsx"
Task: "Integration test real-time content updates in apps/web/src/__tests__/integration/content-updates.test.tsx"
Task: "Integration test fallback content display in apps/web/src/__tests__/integration/content-fallback.test.tsx"
Task: "Integration test content loading performance in apps/web/src/__tests__/integration/content-performance.test.tsx"
```

### Phase 4: Type Definitions [P]
```bash
# T020-T024 can run in parallel within same file (non-conflicting types)
Task: "Implement ContentItem type in apps/web/src/types/content.ts"
Task: "Implement ContentCategory type in apps/web/src/types/content.ts"
Task: "Implement ContentValue types in apps/web/src/types/content.ts"
Task: "Implement ContentVersion type in apps/web/src/types/content.ts"
Task: "Implement ContentMetadata type in apps/web/src/types/content.ts"
```

### Phase 5: Utility Hooks [P]
```bash
# T041-T045 can run in parallel (different files)
Task: "Create useContent hook in apps/web/src/hooks/useContent.ts"
Task: "Create useContentByCategory hook in apps/web/src/hooks/useContentByCategory.ts"
Task: "Create useContentFallback hook in apps/web/src/hooks/useContentFallback.ts"
Task: "Create content validation utilities in apps/web/src/utils/contentValidation.ts"
Task: "Create content transformation utilities in apps/web/src/utils/contentTransform.ts"
```

### Phase 6: E2E Tests [P]
```bash
# T056-T059 can run in parallel (different E2E test files)
Task: "E2E test dynamic content loading flow in apps/web/tests/e2e/content-loading.spec.ts"
Task: "E2E test content update propagation in apps/web/tests/e2e/content-updates.spec.ts"
Task: "E2E test fallback content behavior in apps/web/tests/e2e/content-fallback.spec.ts"
Task: "E2E test content performance requirements in apps/web/tests/e2e/content-performance.spec.ts"
```

## Notes
- [P] tasks target different files or non-conflicting code sections
- Verify all contract tests fail before implementing API endpoints
- Content slice must be implemented before API integration
- Component updates should be done incrementally to maintain app stability
- Fallback system is critical for production reliability

## Task Generation Rules Applied

1. **From Contracts** (content-api.ts):
   - 12 endpoints → 12 contract test tasks (T004-T015)
   - 12 endpoints → 12 implementation tasks (T029-T040)

2. **From Data Model**:
   - 5 main entities → 5 type definition tasks (T020-T024)
   - Content state → slice implementation tasks (T025-T028)

3. **From User Stories** (quickstart.md):
   - 4 test scenarios → 4 integration test tasks (T016-T019)
   - Performance requirements → performance validation tasks (T061)

4. **From Research Decisions**:
   - Fallback system → fallback implementation tasks (T051-T053)
   - Caching strategy → optimization tasks (T054-T055)

## Validation Checklist ✅

- [x] All contracts have corresponding tests (T004-T015 cover all 12 endpoints)
- [x] All entities have model tasks (T020-T024 cover all 5 entities)
- [x] All tests come before implementation (T004-T019 before T020+)
- [x] Parallel tasks truly independent (different files or non-conflicting sections)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task (except type definitions which are additive)
- [x] TDD order maintained: Tests → Types → Slice → API → Hooks → Components
- [x] Integration with existing RTK infrastructure (009-create-rtk-follow)

## Success Criteria

After completing all 63 tasks:
- ✅ All hardcoded content replaced with dynamic RTK-managed data
- ✅ Content loading performance under 200ms
- ✅ Robust fallback system for offline/error scenarios
- ✅ Real-time content updates without page refresh
- ✅ Comprehensive test coverage for all content functionality
- ✅ Seamless integration with existing RTK infrastructure