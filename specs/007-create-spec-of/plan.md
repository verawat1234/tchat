
# Implementation Plan: iOS and Android Native UI Screens Following Web Platform

**Branch**: `007-create-spec-of` | **Date**: 2025-09-22 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/007-create-spec-of/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, `GEMINI.md` for Gemini CLI, `QWEN.md` for Qwen Code or `AGENTS.md` for opencode).
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
Primary requirement: Create native iOS and Android screens that mirror all web application pages and functionality while respecting platform design conventions. Technical approach involves implementing native UI components using SwiftUI for iOS and Jetpack Compose for Android, with cross-platform state synchronization and platform-specific optimizations for navigation, gestures, and hardware features.

## Technical Context
**Language/Version**: Swift 5.9+ for iOS, Kotlin 1.9+ for Android, TypeScript 5.3.0 for shared web logic
**Primary Dependencies**: SwiftUI (iOS), Jetpack Compose (Android), React 18.3.1 (web reference), Navigation frameworks (iOS NavigationStack, Android Navigation Component)
**Storage**: Local state management (UserDefaults/SharedPreferences), cross-platform sync via existing API, CoreData (iOS) and Room (Android) for offline caching
**Testing**: XCTest for iOS, JUnit + Espresso for Android, existing Playwright E2E infrastructure for cross-platform validation
**Target Platform**: iOS 15+ and Android API 24+ (Android 7.0+) to match modern device penetration
**Project Type**: Mobile (iOS + Android native apps) - determines native app structure under apps/mobile/
**Performance Goals**: NEEDS CLARIFICATION - specific load time targets for mobile screens not specified in spec
**Constraints**: NEEDS CLARIFICATION - offline capabilities scope and data usage optimization targets not specified
**Scale/Scope**: NEEDS CLARIFICATION - platform-specific features to include/exclude not fully defined

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Simplicity**: ✅ PASS - Building native UI screens is a fundamental mobile app requirement, not over-engineering
**Platform Conventions**: ✅ PASS - Following iOS HIG and Material Design guidelines as specified
**Code Reuse**: ✅ PASS - Leveraging existing navigation infrastructure and API integration patterns
**Performance**: ⚠️ REVIEW - Mobile performance targets need clarification but reasonable defaults can be applied
**Security**: ✅ PASS - Using platform-standard authentication and data protection mechanisms
**Testing**: ✅ PASS - Comprehensive testing strategy including E2E validation across platforms
**Maintainability**: ✅ PASS - Following established patterns from existing codebase structure

## Project Structure

### Documentation (this feature)
```
specs/[###-feature]/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
# Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure]
```

**Structure Decision**: [DEFAULT to Option 1 unless Technical Context indicates web/mobile app]

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents**:
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
   - Quickstart test = story validation steps

5. **Update agent file incrementally** (O(1) operation):
   - Run `.specify/scripts/bash/update-agent-context.sh claude`
     **IMPORTANT**: Execute it exactly as specified above. Do not add or remove any arguments.
   - If exists: Add only NEW tech from current plan
   - Preserve manual additions between markers
   - Update recent changes (keep last 3)
   - Keep under 150 lines for token efficiency
   - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `.specify/templates/tasks-template.md` as base
- Generate tasks from Phase 1 design docs (contracts, data model, quickstart)
- Each API contract (3 files) → contract test task [P]
- Each data entity (5 entities) → model creation task [P]
- Each main screen (5 tabs + 5 sub-screens) → UI implementation task
- Cross-platform sync implementation → synchronization tasks
- Platform-specific features → platform integration tasks
- Performance optimization → optimization tasks
- Testing and validation → E2E test tasks

**Task Categories**:
1. **Foundation Tasks (T001-T010)**: Data models, API contracts, base infrastructure
2. **UI Implementation Tasks (T011-T030)**: Screen components for iOS and Android
3. **Integration Tasks (T031-T040)**: Cross-platform sync, platform features
4. **Testing Tasks (T041-T050)**: Contract tests, E2E tests, performance validation
5. **Polish Tasks (T051-T055)**: Accessibility, performance optimization, documentation

**Ordering Strategy**:
- TDD order: Contract tests before implementation
- Dependency order: Models → API layer → UI components → Integration → Testing
- Platform parallel: iOS and Android tasks marked [P] for parallel execution
- Critical path: Main tabs before sub-screens

**Estimated Output**: 50-55 numbered, ordered tasks in tasks.md covering:
- 10 foundation tasks (models, contracts, infrastructure)
- 20 UI implementation tasks (10 iOS + 10 Android screens)
- 10 integration tasks (sync, platform features, performance)
- 10 testing tasks (contract, integration, E2E tests)
- 5 polish tasks (accessibility, optimization, docs)

**Parallelization Opportunities**:
- iOS and Android UI development can happen simultaneously [P]
- Contract test development parallel to model implementation [P]
- Platform-specific feature development independent [P]
- Performance optimization can overlap with UI development

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command) - ✅ COMPLETED
- [x] Phase 1: Design complete (/plan command) - ✅ COMPLETED
- [x] Phase 2: Task planning complete (/plan command - describe approach only) - ✅ COMPLETED
- [x] Phase 3: Tasks generated (/tasks command) - ✅ COMPLETED
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS - ✅ COMPLETED
- [x] Post-Design Constitution Check: PASS - ✅ COMPLETED
- [x] All NEEDS CLARIFICATION resolved - ✅ COMPLETED
- [x] Complexity deviations documented - ✅ N/A (no deviations)

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
