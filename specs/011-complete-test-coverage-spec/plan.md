
# Implementation Plan: Complete Test Coverage Specification

**Branch**: `011-complete-test-coverage-spec` | **Date**: 2025-09-22 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/011-complete-test-coverage-spec/spec.md`

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
**Primary Requirement**: Complete enterprise-grade test coverage for Tchat Southeast Asian chat platform microservices. Address critical gaps: Content Service (zero coverage), Unit Tests (0% → 80%), Security Tests, Database Tests, Infrastructure Tests.

**Technical Approach**: 4-phase implementation using Go testing ecosystem with testify framework, establishing comprehensive test patterns across 7 microservices, implementing regional compliance validation, and achieving enterprise quality targets (80% unit, 70% integration, 95% critical path coverage).

## Technical Context
**Language/Version**: Go 1.22+ (microservices backend architecture)
**Primary Dependencies**: testify/suite, testify/mock, testify/assert, go-sqlmock, httptest, dockertest
**Storage**: PostgreSQL (primary), ScyllaDB (messages), Redis (cache/sessions)
**Testing**: Go testing package, testify framework, contract/integration/performance test suite
**Target Platform**: Linux containers (Docker), CI/CD pipelines, Southeast Asian deployment regions
**Project Type**: microservices - existing backend architecture with 7 services
**Performance Goals**: 80% unit coverage, 70% integration coverage, 95% critical path coverage, <200ms regional API response times
**Constraints**: Enterprise compliance, Southeast Asian regional requirements, zero production test failures
**Scale/Scope**: 7 microservices, 98 production files, 29 existing test files → comprehensive test ecosystem

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Test-First Principle**: ✅ PASS - Implementation follows TDD: Contract tests → Integration tests → Unit tests → Implementation
**Library-First Principle**: ✅ PASS - Test utilities and frameworks are self-contained, reusable libraries
**CLI Interface**: ✅ PASS - Test execution via standard Go CLI commands (go test, coverage tools)
**Integration Testing**: ✅ PASS - Comprehensive integration tests for service-to-service communication
**Observability**: ✅ PASS - Test results provide structured logging and coverage reporting
**Simplicity**: ✅ PASS - Standard Go testing patterns, no complex testing frameworks

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

**Structure Decision**: Existing microservices architecture - backend/ with 7 services (auth, messaging, payment, commerce, notification, content, gateway). Test structure follows service boundaries with backend/tests/ for cross-service testing.

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
- Content Service tests → highest priority contract test tasks [P]
- Unit test framework → foundation setup tasks [P]
- Each microservice → service-specific test suite tasks
- Integration workflows → cross-service test tasks
- Security/compliance → regional validation test tasks
- Performance benchmarks → regional performance test tasks

**Ordering Strategy**:
- **Phase 1 (Critical)**: Content Service complete test suite + Unit test framework
- **Phase 2 (High)**: Commerce/Notification/Gateway service tests
- **Phase 3 (Medium)**: Infrastructure + Performance tests
- **Phase 4 (Medium)**: Regional compliance + CI/CD integration
- TDD order: Contract tests → Integration tests → Unit tests → Implementation
- Mark [P] for parallel execution within service boundaries

**Test Task Categories**:
- Service Unit Tests: 7 services × 3-5 test files = ~25 tasks
- Integration Tests: 5 cross-service workflows = 5 tasks
- Security Tests: 6 regions × 3 compliance types = 10 tasks
- Performance Tests: 6 regions × benchmarks = 6 tasks
- Infrastructure Tests: Docker + DB + CI/CD = 5 tasks

**Estimated Output**: 50-60 numbered, ordered tasks in tasks.md (comprehensive test implementation)

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
- [x] Phase 0: Research complete (/plan command) ✅
- [x] Phase 1: Design complete (/plan command) ✅
- [x] Phase 2: Task planning complete (/plan command - describe approach only) ✅
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS ✅
- [x] Post-Design Constitution Check: PASS ✅
- [x] All NEEDS CLARIFICATION resolved ✅
- [x] Complexity deviations documented: N/A (no deviations) ✅

**Artifacts Generated**:
- [x] research.md: Technical stack research and patterns ✅
- [x] data-model.md: Test entities and relationships ✅
- [x] contracts/test-execution-api.md: API contract definitions ✅
- [x] contracts/coverage-reporting-schema.json: Coverage report schema ✅
- [x] quickstart.md: 5-minute validation workflow ✅
- [x] CLAUDE.md: Updated agent context ✅

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
