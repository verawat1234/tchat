# Tasks: Common UI Component Breakdown

**Input**: Design documents from `/specs/001-agent-frontend-specialist/`
**Prerequisites**: plan.md (✓), research.md (✓), data-model.md (✓), contracts/ (✓), quickstart.md (✓)

## Execution Flow (main)
```
1. Load plan.md from feature directory
   ✅ Implementation plan loaded - React/TypeScript component library
   ✅ Extract: TypeScript 5.3.0, React 18.3.1, Vite 6.3.5, TailwindCSS v4
2. Load optional design documents:
   ✅ data-model.md: 8 component entities → model tasks
   ✅ contracts/: 10 interface files → contract test tasks
   ✅ research.md: Technology decisions → setup tasks
3. Generate tasks by category:
   ✅ Setup: project structure, dependencies, testing tools
   ✅ Tests: contract tests, component tests, integration tests
   ✅ Core: component implementations, stories, documentation
   ✅ Integration: theme system, existing components
   ✅ Polish: accessibility, performance, final documentation
4. Apply task rules:
   ✅ Different files = mark [P] for parallel
   ✅ Same file = sequential (no [P])
   ✅ Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
   ✅ 42 tasks generated across 5 phases
6. Generate dependency graph
   ✅ Dependencies mapped with TDD ordering
7. Create parallel execution examples
   ✅ Parallel groups identified and documented
8. Validate task completeness:
   ✅ All contracts have tests
   ✅ All components have implementations
   ✅ All stories have Storybook coverage
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Web app**: `apps/web/src/components/common/`, `apps/web/tests/components/common/`
- Based on plan.md structure: React frontend with component library in `/common/`

## Phase 3.1: Setup
- [ ] **T001** Create component library structure in `apps/web/src/components/common/`
- [ ] **T002** Install testing dependencies (React Testing Library, Vitest, Storybook)
- [ ] **T003** [P] Configure Vitest config in `apps/web/vitest.config.ts`
- [ ] **T004** [P] Configure Storybook in `apps/web/.storybook/main.ts`
- [ ] **T005** [P] Create test setup file in `apps/web/src/test-setup.ts`
- [ ] **T006** [P] Create utility functions in `apps/web/src/utils/cn.ts` for className merging

## Phase 3.2: Contract Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Base Component Tests
- [ ] **T007** [P] Contract test for BaseComponent interface in `apps/web/tests/components/common/base.test.tsx`
- [ ] **T008** [P] Contract test for ThemeTokens in `apps/web/tests/components/common/theme.test.tsx`

### Component Contract Tests
- [ ] **T009** [P] Contract test for Badge props in `apps/web/tests/components/common/Badge.test.tsx`
- [ ] **T010** [P] Contract test for Card props in `apps/web/tests/components/common/Card.test.tsx`
- [ ] **T011** [P] Contract test for Pagination props in `apps/web/tests/components/common/Pagination.test.tsx`
- [ ] **T012** [P] Contract test for Tabs props in `apps/web/tests/components/common/Tabs.test.tsx`
- [ ] **T013** [P] Contract test for Layout props in `apps/web/tests/components/common/Layout.test.tsx`
- [ ] **T014** [P] Contract test for Header props in `apps/web/tests/components/common/Header.test.tsx`
- [ ] **T015** [P] Contract test for Sidebar props in `apps/web/tests/components/common/Sidebar.test.tsx`
- [ ] **T016** [P] Contract test for ChatMessage props in `apps/web/tests/components/common/ChatMessage.test.tsx`

### Integration Tests
- [ ] **T017** [P] Integration test for theme system in `apps/web/tests/integration/theme-integration.test.tsx`
- [ ] **T018** [P] Integration test for component composition in `apps/web/tests/integration/component-composition.test.tsx`
- [ ] **T019** [P] Integration test for accessibility compliance in `apps/web/tests/integration/accessibility.test.tsx`

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Foundation Components
- [ ] **T020** [P] Badge component in `apps/web/src/components/common/badge/Badge.tsx`
- [ ] **T021** [P] Badge export index in `apps/web/src/components/common/badge/index.ts`
- [ ] **T022** [P] Badge Storybook story in `apps/web/src/components/common/badge/Badge.stories.tsx`

### Layout Components
- [ ] **T023** [P] Flex component in `apps/web/src/components/common/layout/Flex.tsx`
- [ ] **T024** [P] Grid component in `apps/web/src/components/common/layout/Grid.tsx`
- [ ] **T025** [P] Stack component in `apps/web/src/components/common/layout/Stack.tsx`
- [ ] **T026** [P] Container component in `apps/web/src/components/common/layout/Container.tsx`
- [ ] **T027** Layout components export index in `apps/web/src/components/common/layout/index.ts`
- [ ] **T028** [P] Layout Storybook stories in `apps/web/src/components/common/layout/Layout.stories.tsx`

### Card Component
- [ ] **T029** [P] Card root component in `apps/web/src/components/common/card/Card.tsx`
- [ ] **T030** [P] CardHeader component in `apps/web/src/components/common/card/CardHeader.tsx`
- [ ] **T031** [P] CardContent component in `apps/web/src/components/common/card/CardContent.tsx`
- [ ] **T032** [P] CardFooter component in `apps/web/src/components/common/card/CardFooter.tsx`
- [ ] **T033** Card components export index in `apps/web/src/components/common/card/index.ts`
- [ ] **T034** [P] Card Storybook stories in `apps/web/src/components/common/card/Card.stories.tsx`

### Navigation Components
- [ ] **T035** [P] Pagination component in `apps/web/src/components/common/pagination/Pagination.tsx`
- [ ] **T036** [P] Pagination Storybook story in `apps/web/src/components/common/pagination/Pagination.stories.tsx`
- [ ] **T037** [P] Tabs root component in `apps/web/src/components/common/tabs/Tabs.tsx`
- [ ] **T038** [P] TabsList component in `apps/web/src/components/common/tabs/TabsList.tsx`
- [ ] **T039** [P] TabsTrigger component in `apps/web/src/components/common/tabs/TabsTrigger.tsx`
- [ ] **T040** [P] TabsContent component in `apps/web/src/components/common/tabs/TabsContent.tsx`
- [ ] **T041** Tabs components export index in `apps/web/src/components/common/tabs/index.ts`
- [ ] **T042** [P] Tabs Storybook stories in `apps/web/src/components/common/tabs/Tabs.stories.tsx`

## Phase 3.4: Complex Components

### Header Component
- [ ] **T043** [P] Header component in `apps/web/src/components/common/header/Header.tsx`
- [ ] **T044** [P] Breadcrumb component in `apps/web/src/components/common/header/Breadcrumb.tsx`
- [ ] **T045** Header components export index in `apps/web/src/components/common/header/index.ts`
- [ ] **T046** [P] Header Storybook stories in `apps/web/src/components/common/header/Header.stories.tsx`

### Sidebar Component
- [ ] **T047** [P] Sidebar root component in `apps/web/src/components/common/sidebar/Sidebar.tsx`
- [ ] **T048** [P] SidebarContent component in `apps/web/src/components/common/sidebar/SidebarContent.tsx`
- [ ] **T049** [P] SidebarNav component in `apps/web/src/components/common/sidebar/SidebarNav.tsx`
- [ ] **T050** [P] SidebarToggle component in `apps/web/src/components/common/sidebar/SidebarToggle.tsx`
- [ ] **T051** Sidebar components export index in `apps/web/src/components/common/sidebar/index.ts`
- [ ] **T052** [P] Sidebar Storybook stories in `apps/web/src/components/common/sidebar/Sidebar.stories.tsx`

### ChatMessage Component
- [ ] **T053** [P] ChatMessage root component in `apps/web/src/components/common/chat-message/ChatMessage.tsx`
- [ ] **T054** [P] MessageBubble component in `apps/web/src/components/common/chat-message/MessageBubble.tsx`
- [ ] **T055** [P] TypingIndicator component in `apps/web/src/components/common/chat-message/TypingIndicator.tsx`
- [ ] **T056** ChatMessage components export index in `apps/web/src/components/common/chat-message/index.ts`
- [ ] **T057** [P] ChatMessage Storybook stories in `apps/web/src/components/common/chat-message/ChatMessage.stories.tsx`

## Phase 3.5: Integration & Polish

### Library Integration
- [ ] **T058** Main export index in `apps/web/src/components/common/index.ts`
- [ ] **T059** Update main components index in `apps/web/src/components/index.ts`
- [ ] **T060** [P] Type definitions export in `apps/web/src/types/components.ts`

### Documentation & Performance
- [ ] **T061** [P] Component usage documentation in `apps/web/docs/components.md`
- [ ] **T062** [P] Performance testing script in `apps/web/scripts/test-performance.js`
- [ ] **T063** [P] Bundle size analysis in `apps/web/scripts/analyze-bundle.js`
- [ ] **T064** [P] Accessibility audit script in `apps/web/scripts/audit-a11y.js`

### Final Validation
- [ ] **T065** Run all tests and ensure they pass
- [ ] **T066** Build Storybook and verify all stories
- [ ] **T067** Run performance benchmarks (<100ms render, <50KB bundle)
- [ ] **T068** Accessibility compliance verification (WCAG 2.1 AA)
- [ ] **T069** Integration with existing codebase validation
- [ ] **T070** Final documentation review and quickstart validation

## Dependencies

### Sequential Dependencies
- **Setup** (T001-T006) → **Contract Tests** (T007-T019) → **Implementation** (T020-T057) → **Integration** (T058-T070)
- **T027** (Layout index) blocks **T028** (Layout stories)
- **T033** (Card index) blocks **T034** (Card stories)
- **T041** (Tabs index) blocks **T042** (Tabs stories)
- **T045** (Header index) blocks **T046** (Header stories)
- **T051** (Sidebar index) blocks **T052** (Sidebar stories)
- **T056** (ChatMessage index) blocks **T057** (ChatMessage stories)
- **T058-T060** (Integration) must complete before **T065-T070** (Validation)

### Parallel Opportunities
All [P] marked tasks can run concurrently within their phase, enabling significant time savings through parallel execution.

## Parallel Execution Examples

### Contract Tests (Phase 3.2)
```bash
# Launch T009-T016 together (component contract tests):
Task: "Contract test for Badge props in apps/web/tests/components/common/Badge.test.tsx"
Task: "Contract test for Card props in apps/web/tests/components/common/Card.test.tsx"
Task: "Contract test for Pagination props in apps/web/tests/components/common/Pagination.test.tsx"
Task: "Contract test for Tabs props in apps/web/tests/components/common/Tabs.test.tsx"
# ... (all 8 component contract tests)
```

### Foundation Components (Phase 3.3)
```bash
# Launch T020-T022 together (Badge implementation):
Task: "Badge component in apps/web/src/components/common/badge/Badge.tsx"
Task: "Badge export index in apps/web/src/components/common/badge/index.ts"
Task: "Badge Storybook story in apps/web/src/components/common/badge/Badge.stories.tsx"
```

### Layout Components (Phase 3.3)
```bash
# Launch T023-T026 together (Layout components):
Task: "Flex component in apps/web/src/components/common/layout/Flex.tsx"
Task: "Grid component in apps/web/src/components/common/layout/Grid.tsx"
Task: "Stack component in apps/web/src/components/common/layout/Stack.tsx"
Task: "Container component in apps/web/src/components/common/layout/Container.tsx"
```

## Notes
- **[P] tasks** = different files, no dependencies
- **Verify tests fail** before implementing components
- **Commit after each task** for clean history
- **Component order**: Simple → Complex (Badge → Layout → Card → Tabs → Pagination → Header → Sidebar → ChatMessage)
- **Storybook integration** throughout development for visual validation
- **Performance monitoring** from day one with <100ms render time target

## Task Generation Rules
*Applied during main() execution*

1. **From Contracts**:
   ✅ 10 contract files → 10 contract test tasks [P]
   ✅ Each component → implementation task

2. **From Data Model**:
   ✅ 8 component entities → 8 implementation groupings
   ✅ Compound components → multiple component tasks

3. **From User Stories**:
   ✅ Developer usage scenarios → integration tests [P]
   ✅ Quickstart scenarios → validation tasks

4. **Ordering**:
   ✅ Setup → Tests → Foundation → Navigation → Complex → Integration
   ✅ TDD ordering maintained throughout

## Validation Checklist
*GATE: Checked by main() before returning*

- [x] All contracts have corresponding tests (T007-T016)
- [x] All component entities have implementation tasks (T020-T057)
- [x] All tests come before implementation (T007-T019 before T020+)
- [x] Parallel tasks truly independent (different files/components)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] TDD ordering enforced (tests must fail before implementation)
- [x] Performance and accessibility requirements included
- [x] Integration with existing codebase planned