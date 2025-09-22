
# Implementation Plan: Native Mobile UI Implementation

**Branch**: `003-creare-spec-for` | **Date**: 2025-09-21 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-creare-spec-for/spec.md`

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
Implement native iOS (Swift/SwiftUI) and Android (Kotlin/Compose) UI components that achieve pixel-perfect visual consistency with the existing comprehensive web design system (React 18.3.1 + Radix UI + TailwindCSS v4 + Framer Motion). The goal is seamless user experience across PWA and native platforms with 40+ component parity, 5-tab navigation architecture, and enhanced native performance while maintaining the established design language and interaction patterns.

## Technical Context
**Language/Version**: iOS: Swift 5.9+ with SwiftUI, Android: Kotlin 1.9+ with Jetpack Compose
**Primary Dependencies**: iOS: UIKit, SwiftUI, Combine; Android: Compose, Material3, Coroutines
**Storage**: CoreData (iOS), Room Database (Android), Shared session state with web
**Testing**: XCTest (iOS), Espresso/Compose Testing (Android), Cross-platform UI tests
**Target Platform**: iOS 15+, Android API 24+ (Android 7.0+)
**Project Type**: mobile - native iOS and Android applications
**Performance Goals**: <2s app launch, <100ms gesture response, 60fps animations, pixel-perfect rendering
**Constraints**: Platform design guidelines compliance, web design system consistency, offline capability
**Scale/Scope**: 40+ UI components, 5-tab architecture, video streaming, workspace management, commerce flows

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Initial Check**: PASS - Constitution template is generic and requires no specific constraints for native mobile UI implementation.

**Post-Design Check**: PASS - Design follows standard mobile development practices:
- Native technologies (Swift/SwiftUI for iOS, Kotlin/Compose for Android) aligned with platform conventions
- Clean architecture with separated concerns (design tokens, components, state management)
- Test-driven approach with contract tests and quickstart validation
- No unnecessary complexity or constitutional violations identified

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

**Structure Decision**: Option 3 (Mobile + API) - Existing project structure at `/apps/mobile/android/` and `/apps/mobile/ios/` with shared backend integration

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

**Output**: ✅ research.md complete - All research areas resolved with platform-specific decisions documented

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

**Output**: ✅ Complete
- data-model.md: Platform-specific entity definitions with Swift/Kotlin implementations
- contracts/: API contracts (design-tokens-api.yaml, component-registry-api.yaml, state-sync-api.yaml)
- quickstart.md: Comprehensive implementation and testing guide
- CLAUDE.md: Updated with iOS/Android native development context

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `.specify/templates/tasks-template.md` as base
- Generate native mobile implementation tasks from Phase 1 design artifacts
- Design Token Implementation: iOS DesignTokens.swift + Android DesignTokens.kt [P]
- Core Component Library: Platform-specific component implementations [P]
- Navigation Architecture: 5-tab navigation for both platforms
- State Synchronization: Cross-platform session and preference sync
- Contract Testing: API contract validation for each platform
- Integration Testing: Cross-platform visual consistency and performance tests

**Platform-Specific Ordering Strategy**:
- TDD order: Contract tests → Component tests → Implementation
- Foundation first: Design tokens → Base components → Navigation → Advanced features
- Parallel execution: iOS and Android implementation can proceed independently [P]
- Integration points: State sync and cross-platform testing require coordination

**Estimated Output**: 30-35 numbered, ordered tasks covering:
- 8-10 Design token and foundation tasks
- 15-20 Component implementation tasks (iOS + Android)
- 5-8 Navigation and state management tasks
- 5-7 Testing and validation tasks

**Platform Balance**: Equal task distribution between iOS and Android platforms

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
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented (none required)

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
