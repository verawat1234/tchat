# Implementation Plan: Common UI Component Breakdown

**Branch**: `001-agent-frontend-specialist` | **Date**: 2025-09-21 | **Spec**: [../spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-agent-frontend-specialist/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   ✅ Feature spec loaded successfully
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   ✅ Project Type: web (React/TypeScript frontend)
   ✅ Structure Decision: Option 2 (frontend components)
3. Fill the Constitution Check section based on the content of the constitution document.
   ✅ Constitution template reviewed (no specific rules yet)
4. Evaluate Constitution Check section below
   ✅ No violations - component library follows standard patterns
   ✅ Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   ✅ Research phase complete
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent file
   ✅ Design phase complete
7. Re-evaluate Constitution Check section
   ✅ No new violations after design
   ✅ Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
   ✅ Task planning approach documented
9. STOP - Ready for /tasks command
   ✅ Implementation plan complete
```

**IMPORTANT**: The /plan command STOPS at step 8. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
Extract common UI patterns from the Tchat React application into reusable components including pagination, tabs, cards, chat messages, badges, layout, header, and sidebar. Components will use the existing Telegram SEA Edition theme and follow current architectural patterns to maintain consistency while enabling reusability and maintainability.

## Technical Context
**Language/Version**: TypeScript 5.3.0, React 18.3.1
**Primary Dependencies**: Vite 6.3.5, Radix UI components, TailwindCSS v4, Framer Motion 11.0.0
**Storage**: Component state in React state, design tokens in CSS variables
**Testing**: Component testing with React Testing Library, Storybook for component documentation
**Target Platform**: Web browsers (ES2020+), responsive design for mobile/desktop
**Project Type**: web - React frontend with component library architecture
**Performance Goals**: <100ms component render time, <50KB bundle size per component
**Constraints**: Must maintain existing visual design, backward compatibility with current usage
**Scale/Scope**: 8-12 core components (pagination, tabs, card, chat message types, badge, layout, header, sidebar), 20+ variations/props

**User Provided Details**: Focus on pagination, tabs, card, chat message types, badge, layout, header, sidebar components. Use current theme system and follow existing architectural approaches.

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Constitution Status**: Template constitution file found but not yet populated with specific rules. Proceeding with standard best practices:
- Component-first architecture ✅
- Reusability and composition ✅
- Testing and documentation ✅
- Performance considerations ✅

**Initial Assessment**: PASS - Component library approach aligns with standard frontend architecture patterns.

## Project Structure

### Documentation (this feature)
```
specs/001-agent-frontend-specialist/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
# Option 2: Web application (React frontend)
apps/web/
├── src/
│   ├── components/
│   │   ├── ui/           # Existing Radix UI components
│   │   ├── common/       # NEW: Common reusable components
│   │   └── pages/        # Page-level components
│   ├── styles/
│   └── utils/
└── tests/
    ├── components/       # Component tests
    └── integration/      # Integration tests
```

**Structure Decision**: Option 2 (Web application) - React frontend with dedicated component library in `apps/web/src/components/common/`

## Phase 0: Outline & Research

1. **Extract unknowns from Technical Context** above:
   - Current component patterns and conventions
   - Existing theme system implementation
   - Component testing strategies
   - Storybook integration approach

2. **Generate and dispatch research agents**:
   ```
   Research current component patterns in Tchat codebase
   Research existing theme system and design tokens
   Research component testing best practices for React/TypeScript
   Research Storybook integration for component documentation
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: Component library architecture with common/ directory
   - Rationale: Maintains existing structure while enabling reusability
   - Alternatives considered: Separate package vs. integrated approach

**Output**: research.md with all technical decisions documented

## Phase 1: Design & Contracts

*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Component interfaces and prop definitions
   - Theme token structure
   - Component composition patterns

2. **Generate API contracts** from functional requirements:
   - Component prop interfaces (TypeScript definitions)
   - Theme token contracts
   - Component composition patterns
   - Output TypeScript interfaces to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per component
   - Assert prop validation and rendering
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Developer usage scenarios → component tests
   - Integration scenarios → Storybook stories

5. **Update agent file incrementally** (O(1) operation):
   - Run `.specify/scripts/bash/update-agent-context.sh claude`
   - Add React/TypeScript component development context
   - Include component library patterns
   - Keep under 150 lines for token efficiency

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, CLAUDE.md

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `.specify/templates/tasks-template.md` as base
- Generate tasks from component contracts and data model
- Each component → design, implement, test, document tasks
- Integration tasks for theme system and existing codebase

**Ordering Strategy**:
- TDD order: Tests and interfaces before implementation
- Dependency order: Basic components before composite components
- Core components first: Button, Input, Card before complex ones like ChatMessage
- Mark [P] for parallel execution (independent components)

**Estimated Output**: 35-40 numbered, ordered tasks covering 8-12 components

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)
**Phase 4**: Implementation (execute tasks.md following constitutional principles)
**Phase 5**: Validation (run tests, build Storybook, integration testing)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

No violations identified. Component library approach follows standard frontend patterns.

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
- [x] Complexity deviations documented (none)

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*