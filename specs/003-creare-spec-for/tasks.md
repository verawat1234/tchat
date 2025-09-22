# Tasks: Native Mobile UI Implementation

**Input**: Design documents from `/specs/003-creare-spec-for/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Tech stack: iOS (Swift/SwiftUI), Android (Kotlin/Compose)
   → Project structure: Mobile apps at /apps/mobile/ios/ and /apps/mobile/android/
2. Load design documents:
   → data-model.md: 5 entities (DesignTokens, Components, Gestures, StateSync, Integrations)
   → contracts/: 3 API contracts (design-tokens, component-registry, state-sync)
   → research.md: Platform-specific implementation decisions
   → quickstart.md: Testing scenarios and validation criteria
3. Generate tasks by category:
   → Setup: iOS/Android project configuration
   → Tests: Contract tests, component tests, integration tests
   → Core: Design tokens, components, navigation, state sync
   → Integration: Platform integrations, performance optimization
   → Polish: Accessibility, testing, documentation
4. Apply task rules:
   → iOS/Android parallel implementation [P]
   → Platform-specific files can run in parallel
   → Shared integration points are sequential
5. Generate 35 numbered tasks (T001-T035)
6. Dependencies: Setup → Tests → Implementation → Integration → Polish
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files/platforms, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **iOS**: `/apps/mobile/ios/Sources/` for implementation, `/apps/mobile/ios/Tests/` for tests
- **Android**: `/apps/mobile/android/app/src/main/java/` for implementation, `/apps/mobile/android/app/src/test/java/` for tests
- **Shared**: `/specs/003-creare-spec-for/contracts/` for API contracts

## Phase 3.1: Setup & Infrastructure
- [x] T001 [P] Configure iOS project dependencies (SwiftUI, Combine) in /apps/mobile/ios/Package.swift
- [x] T002 [P] Configure Android project dependencies (Compose, Material3, Coroutines) in /apps/mobile/android/app/build.gradle
- [x] T003 [P] Setup iOS linting and formatting (SwiftLint) in /apps/mobile/ios/.swiftlint.yml
- [x] T004 [P] Setup Android linting and formatting (ktlint) in /apps/mobile/android/app/build.gradle
- [x] T005 [P] Configure iOS testing framework (XCTest) in /apps/mobile/ios/Tests/
- [x] T006 [P] Configure Android testing framework (Espresso, Compose Testing) in /apps/mobile/android/app/src/test/

## Phase 3.2: Contract Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [x] T007 [P] Contract test design-tokens API in /apps/mobile/ios/Tests/ContractTests/DesignTokensAPITests.swift
- [x] T008 [P] Contract test design-tokens API in /apps/mobile/android/app/src/test/java/ContractTests/DesignTokensAPITest.kt
- [x] T009 [P] Contract test component-registry API in /apps/mobile/ios/Tests/ContractTests/ComponentRegistryAPITests.swift
- [x] T010 [P] Contract test component-registry API in /apps/mobile/android/app/src/test/java/ContractTests/ComponentRegistryAPITest.kt
- [x] T011 [P] Contract test state-sync API in /apps/mobile/ios/Tests/ContractTests/StateSyncAPITests.swift
- [x] T012 [P] Contract test state-sync API in /apps/mobile/android/app/src/test/java/ContractTests/StateSyncAPITest.kt

## Phase 3.3: Design Token Implementation (ONLY after tests are failing)
- [x] T013 [P] iOS DesignTokens struct in /apps/mobile/ios/Sources/DesignSystem/DesignTokens.swift
- [x] T014 [P] Android DesignTokens object in /apps/mobile/android/app/src/main/java/com/tchat/designsystem/DesignTokens.kt
- [x] T015 [P] iOS Typography tokens in /apps/mobile/ios/Sources/DesignSystem/Typography.swift
- [x] T016 [P] Android Typography tokens in /apps/mobile/android/app/src/main/java/com/tchat/designsystem/Typography.kt
- [x] T017 [P] iOS Colors system in /apps/mobile/ios/Sources/DesignSystem/Colors.swift
- [x] T018 [P] Android Colors system in /apps/mobile/android/app/src/main/java/com/tchat/designsystem/Colors.kt
- [x] T019 [P] iOS Spacing constants in /apps/mobile/ios/Sources/DesignSystem/Spacing.swift
- [x] T020 [P] Android Spacing constants in /apps/mobile/android/app/src/main/java/com/tchat/designsystem/Spacing.kt

## Phase 3.4: Core Component Library
- [x] T021 [P] iOS TchatButton component in /apps/mobile/ios/Sources/Components/TchatButton.swift
- [x] T022 [P] Android TchatButton component in /apps/mobile/android/app/src/main/java/com/tchat/components/TchatButton.kt
- [x] T023 [P] iOS TchatInput component in /apps/mobile/ios/Sources/Components/TchatInput.swift
- [x] T024 [P] Android TchatInput component in /apps/mobile/android/app/src/main/java/com/tchat/components/TchatInput.kt
- [x] T025 [P] iOS TchatCard component in /apps/mobile/ios/Sources/Components/TchatCard.swift
- [x] T026 [P] Android TchatCard component in /apps/mobile/android/app/src/main/java/com/tchat/components/TchatCard.kt

## Phase 3.5: Navigation Architecture
- [x] T027 [P] iOS MainTabView (5-tab navigation) in /apps/mobile/ios/Sources/Navigation/MainTabView.swift
- [x] T028 [P] Android MainScreen (5-tab navigation) in /apps/mobile/android/app/src/main/java/com/tchat/navigation/MainScreen.kt
- [x] T029 [P] iOS tab view implementations (Chat/Store/Social/Video/More) in /apps/mobile/ios/Sources/Screens/
- [x] T030 [P] Android screen implementations (Chat/Store/Social/Video/More) in /apps/mobile/android/app/src/main/java/com/tchat/screens/

## Phase 3.6: State Synchronization & Integration
- [x] T031 State sync service integration in both platforms (requires coordination)
- [x] T032 Cross-platform session management implementation
- [x] T033 Theme synchronization between web and mobile platforms

## Phase 3.7: Polish & Validation
- [x] T034 [P] Visual consistency validation tests in /apps/mobile/ios/Tests/VisualTests/ and /apps/mobile/android/app/src/test/java/VisualTests/
- [x] T035 [P] Performance benchmarking and accessibility compliance validation

## Dependencies
- Setup (T001-T006) before everything else
- Contract tests (T007-T012) before implementation (T013-T035)
- Design tokens (T013-T020) before components (T021-T026)
- Components before navigation (T027-T030)
- Core implementation before integration (T031-T033)
- Implementation before polish (T034-T035)

## Parallel Execution Examples

### Phase 3.1 - Setup (All Parallel)
```bash
# Launch T001-T006 together:
Task: "Configure iOS project dependencies (SwiftUI, Combine) in /apps/mobile/ios/Package.swift"
Task: "Configure Android project dependencies (Compose, Material3, Coroutines) in /apps/mobile/android/app/build.gradle"
Task: "Setup iOS linting and formatting (SwiftLint) in /apps/mobile/ios/.swiftlint.yml"
Task: "Setup Android linting and formatting (ktlint) in /apps/mobile/android/app/build.gradle"
Task: "Configure iOS testing framework (XCTest) in /apps/mobile/ios/Tests/"
Task: "Configure Android testing framework (Espresso, Compose Testing) in /apps/mobile/android/app/src/test/"
```

### Phase 3.2 - Contract Tests (All Parallel)
```bash
# Launch T007-T012 together:
Task: "Contract test design-tokens API in /apps/mobile/ios/Tests/ContractTests/DesignTokensAPITests.swift"
Task: "Contract test design-tokens API in /apps/mobile/android/app/src/test/java/ContractTests/DesignTokensAPITest.kt"
Task: "Contract test component-registry API in /apps/mobile/ios/Tests/ContractTests/ComponentRegistryAPITests.swift"
Task: "Contract test component-registry API in /apps/mobile/android/app/src/test/java/ContractTests/ComponentRegistryAPITest.kt"
Task: "Contract test state-sync API in /apps/mobile/ios/Tests/ContractTests/StateSyncAPITests.swift"
Task: "Contract test state-sync API in /apps/mobile/android/app/src/test/java/ContractTests/StateSyncAPITest.kt"
```

### Phase 3.3 - Design Tokens (Platform Pairs in Parallel)
```bash
# Launch T013-T014 together (DesignTokens):
Task: "iOS DesignTokens struct in /apps/mobile/ios/Sources/DesignSystem/DesignTokens.swift"
Task: "Android DesignTokens object in /apps/mobile/android/app/src/main/java/com/tchat/designsystem/DesignTokens.kt"

# Launch T015-T016 together (Typography):
Task: "iOS Typography tokens in /apps/mobile/ios/Sources/DesignSystem/Typography.swift"
Task: "Android Typography tokens in /apps/mobile/android/app/src/main/java/com/tchat/designsystem/Typography.kt"
```

### Phase 3.4 - Components (Platform Pairs in Parallel)
```bash
# Launch T021-T022 together (Button):
Task: "iOS TchatButton component in /apps/mobile/ios/Sources/Components/TchatButton.swift"
Task: "Android TchatButton component in /apps/mobile/android/app/src/main/java/com/tchat/components/TchatButton.kt"

# Launch T023-T024 together (Input):
Task: "iOS TchatInput component in /apps/mobile/ios/Sources/Components/TchatInput.swift"
Task: "Android TchatInput component in /apps/mobile/android/app/src/main/java/com/tchat/components/TchatInput.kt"
```

## Task Validation Criteria

### Visual Consistency Requirements
- Color values must match web design system within 1% tolerance
- Typography scales must be pixel-perfect with web implementation
- Spacing must follow 4px grid system exactly
- Animation timing must match web within 50ms tolerance

### Performance Requirements
- App launch time <2 seconds
- Gesture response time <100ms
- 60fps maintained during animations
- Memory usage within platform limits

### Accessibility Requirements
- Touch targets ≥44pt (iOS) / ≥48dp (Android)
- Screen reader compatibility (VoiceOver/TalkBack)
- Dynamic type support
- High contrast mode support

### Cross-Platform Sync Requirements
- Theme changes sync within 5 seconds
- Workspace switching maintains context
- Session state preserved across platforms
- Notification preferences sync correctly

## Notes
- [P] tasks = different platforms/files, no dependencies
- Verify all contract tests fail before implementing
- Commit after each task completion
- Run platform-specific linting before each commit
- Test on both iOS simulator and Android emulator
- Validate accessibility compliance with platform tools

## Task Generation Rules Applied

1. **From Contracts**: 3 API contracts × 2 platforms = 6 contract test tasks [P]
2. **From Data Model**: 5 entities across 2 platforms = 10 platform-specific tasks [P]
3. **From Quickstart**: Testing scenarios → validation tasks
4. **Ordering**: Setup → Tests → Design Tokens → Components → Navigation → Integration → Polish
5. **Platform Parallelization**: iOS and Android implementations can run simultaneously

## Validation Checklist
*GATE: Checked before task execution*

- [x] All contracts have corresponding tests (T007-T012)
- [x] All entities have model tasks (T013-T030)
- [x] All tests come before implementation (T007-T012 before T013-T035)
- [x] Parallel tasks truly independent (different platforms/files)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] Platform-specific implementation maintains parity
- [x] Cross-platform integration points identified and sequenced