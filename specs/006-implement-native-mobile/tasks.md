# Tasks: Native Mobile UI & Routing Parity

**Input**: Design documents from `/specs/006-implement-native-mobile/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Tech stack: Swift 5.9+ (iOS), Kotlin 1.9+ (Android), SwiftUI + Combine, Jetpack Compose + Coroutines
   → Project structure: Mobile apps at /apps/mobile/ios/ and /apps/mobile/android/
2. Load design documents:
   → data-model.md: 6 entities (NavigationRoute, UIComponentState, NavigationStack, DeepLinkHandler, PlatformAdapter, LayoutContainer)
   → contracts/: 3 API contracts (navigation-sync, ui-component-sync, platform-adapter)
   → research.md: Platform-specific implementation decisions
   → quickstart.md: 10 testing scenarios and validation criteria
3. Generate tasks by category:
   → Setup: iOS/Android project extension, dependencies, linting
   → Tests: Contract tests, integration tests, visual validation tests
   → Core: Navigation system, UI state management, platform adapters
   → Integration: API integration, cross-platform sync, deep linking
   → Polish: Performance testing, accessibility validation, documentation
4. Apply task rules:
   → iOS/Android parallel implementation [P]
   → Platform-specific files can run in parallel
   → Shared integration points are sequential
5. Generate 40 numbered tasks (T001-T040)
6. Dependencies: Setup → Tests → Implementation → Integration → Polish
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files/platforms, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **iOS**: `/apps/mobile/ios/Sources/` for implementation, `/apps/mobile/ios/Tests/` for tests
- **Android**: `/apps/mobile/android/app/src/main/java/` for implementation, `/apps/mobile/android/app/src/test/java/` for tests
- **Shared**: `/specs/006-implement-native-mobile/contracts/` for API contracts

## Phase 4.1: Setup & Infrastructure Extension
- [ ] T001 [P] Extend iOS project dependencies for navigation and UI sync in /apps/mobile/ios/Package.swift
- [ ] T002 [P] Extend Android project dependencies for navigation and UI sync in /apps/mobile/android/app/build.gradle
- [ ] T003 [P] Setup iOS navigation framework integration in /apps/mobile/ios/Sources/Navigation/
- [ ] T004 [P] Setup Android navigation framework integration in /apps/mobile/android/app/src/main/java/com/tchat/navigation/

## Phase 4.2: Contract Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 4.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [ ] T005 [P] Contract test navigation-sync API in /apps/mobile/ios/Tests/ContractTests/NavigationSyncAPITests.swift
- [ ] T006 [P] Contract test navigation-sync API in /apps/mobile/android/app/src/test/java/ContractTests/NavigationSyncAPITest.kt
- [ ] T007 [P] Contract test ui-component-sync API in /apps/mobile/ios/Tests/ContractTests/UIComponentSyncAPITests.swift
- [ ] T008 [P] Contract test ui-component-sync API in /apps/mobile/android/app/src/test/java/ContractTests/UIComponentSyncAPITest.kt
- [ ] T009 [P] Contract test platform-adapter API in /apps/mobile/ios/Tests/ContractTests/PlatformAdapterAPITests.swift
- [ ] T010 [P] Contract test platform-adapter API in /apps/mobile/android/app/src/test/java/ContractTests/PlatformAdapterAPITest.kt

## Phase 4.3: Core Entity Models (ONLY after tests are failing)
- [ ] T011 [P] NavigationRoute model in /apps/mobile/ios/Sources/Models/NavigationRoute.swift
- [ ] T012 [P] NavigationRoute model in /apps/mobile/android/app/src/main/java/com/tchat/models/NavigationRoute.kt
- [ ] T013 [P] UIComponentState model in /apps/mobile/ios/Sources/Models/UIComponentState.swift
- [ ] T014 [P] UIComponentState model in /apps/mobile/android/app/src/main/java/com/tchat/models/UIComponentState.kt
- [ ] T015 [P] NavigationStack model in /apps/mobile/ios/Sources/Models/NavigationStack.swift
- [ ] T016 [P] NavigationStack model in /apps/mobile/android/app/src/main/java/com/tchat/models/NavigationStack.kt
- [ ] T017 [P] DeepLinkHandler model in /apps/mobile/ios/Sources/Models/DeepLinkHandler.swift
- [ ] T018 [P] DeepLinkHandler model in /apps/mobile/android/app/src/main/java/com/tchat/models/DeepLinkHandler.kt
- [ ] T019 [P] PlatformAdapter model in /apps/mobile/ios/Sources/Models/PlatformAdapter.swift
- [ ] T020 [P] PlatformAdapter model in /apps/mobile/android/app/src/main/java/com/tchat/models/PlatformAdapter.kt
- [ ] T021 [P] LayoutContainer model in /apps/mobile/ios/Sources/Models/LayoutContainer.swift
- [ ] T022 [P] LayoutContainer model in /apps/mobile/android/app/src/main/java/com/tchat/models/LayoutContainer.kt

## Phase 4.4: Navigation Architecture Implementation
- [ ] T023 [P] NavigationCoordinator service in /apps/mobile/ios/Sources/Services/NavigationCoordinator.swift
- [ ] T024 [P] NavigationCoordinator service in /apps/mobile/android/app/src/main/java/com/tchat/services/NavigationCoordinator.kt
- [ ] T025 [P] RouteRegistry service in /apps/mobile/ios/Sources/Services/RouteRegistry.swift
- [ ] T026 [P] RouteRegistry service in /apps/mobile/android/app/src/main/java/com/tchat/services/RouteRegistry.kt
- [ ] T027 [P] DeepLinkProcessor service in /apps/mobile/ios/Sources/Services/DeepLinkProcessor.swift
- [ ] T028 [P] DeepLinkProcessor service in /apps/mobile/android/app/src/main/java/com/tchat/services/DeepLinkProcessor.kt

## Phase 4.5: UI State Management System
- [ ] T029 [P] UIStateManager service in /apps/mobile/ios/Sources/Services/UIStateManager.swift
- [ ] T030 [P] UIStateManager service in /apps/mobile/android/app/src/main/java/com/tchat/services/UIStateManager.kt
- [ ] T031 [P] ComponentSyncService in /apps/mobile/ios/Sources/Services/ComponentSyncService.swift
- [ ] T032 [P] ComponentSyncService in /apps/mobile/android/app/src/main/java/com/tchat/services/ComponentSyncService.kt

## Phase 4.6: Platform-Specific Adapters
- [ ] T033 [P] iOS PlatformAdapterImpl in /apps/mobile/ios/Sources/Adapters/PlatformAdapterImpl.swift
- [ ] T034 [P] Android PlatformAdapterImpl in /apps/mobile/android/app/src/main/java/com/tchat/adapters/PlatformAdapterImpl.kt
- [ ] T035 [P] iOS GestureHandler in /apps/mobile/ios/Sources/Adapters/GestureHandler.swift
- [ ] T036 [P] Android GestureHandler in /apps/mobile/android/app/src/main/java/com/tchat/adapters/GestureHandler.kt

## Phase 4.7: Integration & API Connectivity
- [ ] T037 Navigation sync API integration in both platforms (requires coordination)
- [ ] T038 UI component sync API integration in both platforms (requires coordination)
- [ ] T039 Platform adapter API integration in both platforms (requires coordination)

## Phase 4.8: Validation & Testing
- [ ] T040 [P] Cross-platform navigation consistency validation tests and performance benchmarking

## Dependencies
- Setup (T001-T004) before everything else
- Contract tests (T005-T010) before implementation (T011-T040)
- Models (T011-T022) before services (T023-T036)
- Services before integration (T037-T039)
- Implementation before validation (T040)

## Parallel Execution Examples

### Phase 4.1 - Setup (All Parallel)
```bash
# Launch T001-T004 together:
Task: "Extend iOS project dependencies for navigation and UI sync in /apps/mobile/ios/Package.swift"
Task: "Extend Android project dependencies for navigation and UI sync in /apps/mobile/android/app/build.gradle"
Task: "Setup iOS navigation framework integration in /apps/mobile/ios/Sources/Navigation/"
Task: "Setup Android navigation framework integration in /apps/mobile/android/app/src/main/java/com/tchat/navigation/"
```

### Phase 4.2 - Contract Tests (All Parallel)
```bash
# Launch T005-T010 together:
Task: "Contract test navigation-sync API in /apps/mobile/ios/Tests/ContractTests/NavigationSyncAPITests.swift"
Task: "Contract test navigation-sync API in /apps/mobile/android/app/src/test/java/ContractTests/NavigationSyncAPITest.kt"
Task: "Contract test ui-component-sync API in /apps/mobile/ios/Tests/ContractTests/UIComponentSyncAPITests.swift"
Task: "Contract test ui-component-sync API in /apps/mobile/android/app/src/test/java/ContractTests/UIComponentSyncAPITest.kt"
Task: "Contract test platform-adapter API in /apps/mobile/ios/Tests/ContractTests/PlatformAdapterAPITests.swift"
Task: "Contract test platform-adapter API in /apps/mobile/android/app/src/test/java/ContractTests/PlatformAdapterAPITest.kt"
```

### Phase 4.3 - Models (Platform Pairs in Parallel)
```bash
# Launch T011-T012 together (NavigationRoute):
Task: "NavigationRoute model in /apps/mobile/ios/Sources/Models/NavigationRoute.swift"
Task: "NavigationRoute model in /apps/mobile/android/app/src/main/java/com/tchat/models/NavigationRoute.kt"

# Launch T013-T014 together (UIComponentState):
Task: "UIComponentState model in /apps/mobile/ios/Sources/Models/UIComponentState.swift"
Task: "UIComponentState model in /apps/mobile/android/app/src/main/java/com/tchat/models/UIComponentState.kt"
```

### Phase 4.4 - Services (Platform Pairs in Parallel)
```bash
# Launch T023-T024 together (NavigationCoordinator):
Task: "NavigationCoordinator service in /apps/mobile/ios/Sources/Services/NavigationCoordinator.swift"
Task: "NavigationCoordinator service in /apps/mobile/android/app/src/main/java/com/tchat/services/NavigationCoordinator.kt"

# Launch T025-T026 together (RouteRegistry):
Task: "RouteRegistry service in /apps/mobile/ios/Sources/Services/RouteRegistry.swift"
Task: "RouteRegistry service in /apps/mobile/android/app/src/main/java/com/tchat/services/RouteRegistry.kt"
```

## Task Validation Criteria

### Navigation Implementation Requirements
- URL-like routing structure compatible with web application
- Platform-specific navigation patterns (iOS Navigation Stack, Android Navigation Component)
- Deep linking support with parameter parsing and validation
- Navigation state persistence across app lifecycle events
- Cross-platform navigation consistency within mobile platform conventions

### UI State Management Requirements
- Component state synchronization across platforms
- Design system parity with web application (within 1% color tolerance)
- Responsive layout adaptation for different device sizes and orientations
- Accessibility compliance (VoiceOver/TalkBack support, WCAG AA standards)
- Performance optimization (60fps animations, <100ms gesture response)

### Integration Requirements
- API contract compliance for all three service endpoints
- Cross-platform state synchronization within 5 seconds
- Offline capability with graceful degradation
- Error handling and retry mechanisms for network operations
- Security compliance for authentication and data transmission

### Performance Requirements
- App launch time <2 seconds with new navigation system
- Navigation response time <100ms
- 60fps maintained during transitions and animations
- Memory usage within platform limits with new state management
- Network request optimization with caching and batching

### Platform-Specific Requirements
- iOS: SwiftUI declarative UI patterns, Combine reactive programming
- Android: Jetpack Compose declarative UI, Coroutines + Flow reactive streams
- Both: Platform accessibility APIs, native gesture recognition, hardware-accelerated animations

## Notes
- [P] tasks = different platforms/files, no dependencies
- Building upon existing mobile foundation (T001-T035 from previous implementation)
- Verify all contract tests fail before implementing
- Commit after each task completion
- Run platform-specific linting before each commit
- Test on both iOS simulator and Android emulator
- Validate cross-platform consistency with automated testing tools
- Ensure accessibility compliance with platform tools (Accessibility Inspector, Accessibility Scanner)

## Task Generation Rules Applied

1. **From Contracts**: 3 API contracts × 2 platforms = 6 contract test tasks [P]
2. **From Data Model**: 6 entities × 2 platforms = 12 model tasks [P]
3. **From Navigation Requirements**: Navigation services across 2 platforms = 6 service tasks [P]
4. **From UI State Requirements**: State management services across 2 platforms = 4 service tasks [P]
5. **From Platform Requirements**: Platform adapters across 2 platforms = 4 adapter tasks [P]
6. **From Integration**: 3 API integrations requiring coordination = 3 sequential tasks
7. **From Validation**: Cross-platform testing requirements = 1 validation task [P]

## Validation Checklist
*GATE: Checked before task execution*

- [x] All contracts have corresponding tests (T005-T010)
- [x] All entities have model tasks (T011-T022)
- [x] All tests come before implementation (T005-T010 before T011-T040)
- [x] Parallel tasks truly independent (different platforms/files)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] Platform-specific implementation maintains parity
- [x] Cross-platform integration points identified and sequenced

**Total Tasks**: 40 (T001-T040)
**Parallel Tasks**: 34 (85% parallelizable)
**Sequential Dependencies**: 6 coordination points