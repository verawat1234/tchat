# Tasks: iOS and Android Native UI Screens Following Web Platform

**Input**: Design documents from `/specs/007-create-spec-of/`
**Prerequisites**: plan.md (✅), research.md (✅), data-model.md (✅), contracts/ (✅)

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → ✅ Found: iOS/Android native UI implementation
   → Tech stack: Swift 5.9+/SwiftUI, Kotlin 1.9+/Jetpack Compose, TypeScript 5.3.0
2. Load design documents:
   → data-model.md: 5 entities (ScreenComponent, NavigationState, ComponentState, SyncEvent, PlatformIntegration)
   → contracts/: 3 API contracts (navigation-sync, ui-component-sync, platform-adapter)
   → research.md: Performance targets, offline scope, platform features
3. Generate tasks by category:
   → Setup: Mobile app structure, dependencies, platform configs
   → Tests: Contract tests, E2E tests, platform validation
   → Core: Data models, screen components, navigation
   → Integration: Cross-platform sync, platform features
   → Polish: Accessibility, performance, documentation
4. Apply task rules:
   → iOS and Android parallel development [P]
   → Contract tests before implementation (TDD)
   → Main tabs before sub-screens
5. Number tasks sequentially (T001-T055)
6. Generate dependency graph and parallel execution
7. SUCCESS: 55 tasks ready for execution across 5 categories
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files/platforms, no dependencies)
- Exact file paths included in descriptions

## Path Conventions
**Mobile project structure** (as defined in plan.md):
- **iOS**: `apps/mobile/ios/src/`, `apps/mobile/ios/Tests/`
- **Android**: `apps/mobile/android/app/src/`, `apps/mobile/android/app/src/test/`
- **Shared**: `apps/mobile/shared/`, existing navigation infrastructure

## Phase 3.1: Foundation & Setup (T001-T010)

- [ ] T001 Create mobile app directory structure per plan.md specifications
- [ ] T002 [P] Initialize iOS project with SwiftUI and iOS NavigationStack dependencies
- [ ] T003 [P] Initialize Android project with Jetpack Compose and Navigation Component dependencies
- [ ] T004 [P] Configure iOS project settings (iOS 15+ target, Swift 5.9+, Xcode build)
- [ ] T005 [P] Configure Android project settings (API 24+ target, Kotlin 1.9+, Gradle build)
- [ ] T006 [P] Set up iOS testing framework (XCTest, test targets, CI integration)
- [ ] T007 [P] Set up Android testing framework (JUnit, Espresso, test configurations)
- [ ] T008 [P] Configure iOS linting and formatting tools (SwiftLint, SwiftFormat)
- [ ] T009 [P] Configure Android linting and formatting tools (ktlint, detekt)
- [ ] T010 [P] Set up cross-platform E2E testing infrastructure with Playwright

## Phase 3.2: Contract Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

- [ ] T011 [P] Contract test Navigation Sync API in apps/mobile/android/app/src/test/java/com/tchat/NavigationSyncAPITest.kt
- [ ] T012 [P] Contract test UI Component Sync API in apps/mobile/android/app/src/test/java/com/tchat/UIComponentSyncAPITest.kt
- [ ] T013 [P] Contract test Platform Adapter API in apps/mobile/android/app/src/test/java/com/tchat/PlatformAdapterAPITest.kt
- [ ] T014 [P] Contract test Navigation Sync API in apps/mobile/ios/Tests/NavigationSyncAPITests.swift
- [ ] T015 [P] Contract test UI Component Sync API in apps/mobile/ios/Tests/UIComponentSyncAPITests.swift
- [ ] T016 [P] Contract test Platform Adapter API in apps/mobile/ios/Tests/PlatformAdapterAPITests.swift
- [ ] T017 [P] Integration test cross-platform navigation sync in apps/mobile/android/app/src/test/java/com/tchat/NavigationConsistencyValidationTest.kt
- [ ] T018 [P] Integration test cross-platform navigation sync in apps/mobile/ios/Tests/NavigationConsistencyValidationTests.swift
- [ ] T019 [P] E2E test main tab navigation flow in tests/e2e/mobile-navigation.spec.ts
- [ ] T020 [P] E2E test cross-platform state synchronization in tests/e2e/cross-platform-sync.spec.ts

## Phase 3.3: Data Models (ONLY after tests are failing) (T021-T030)

- [ ] T021 [P] ScreenComponent model in apps/mobile/ios/src/models/ScreenComponent.swift
- [ ] T022 [P] NavigationState model in apps/mobile/ios/src/models/NavigationState.swift
- [ ] T023 [P] ComponentState model in apps/mobile/ios/src/models/ComponentState.swift
- [ ] T024 [P] SyncEvent model in apps/mobile/ios/src/models/SyncEvent.swift
- [ ] T025 [P] PlatformIntegration model in apps/mobile/ios/src/models/PlatformIntegration.swift
- [ ] T026 [P] ScreenComponent model in apps/mobile/android/app/src/main/java/com/tchat/models/ScreenComponent.kt
- [ ] T027 [P] NavigationState model in apps/mobile/android/app/src/main/java/com/tchat/models/NavigationState.kt
- [ ] T028 [P] ComponentState model in apps/mobile/android/app/src/main/java/com/tchat/models/ComponentState.kt
- [ ] T029 [P] SyncEvent model in apps/mobile/android/app/src/main/java/com/tchat/models/SyncEvent.kt
- [ ] T030 [P] PlatformIntegration model in apps/mobile/android/app/src/main/java/com/tchat/models/PlatformIntegration.kt

## Phase 3.4: API Service Layer (T031-T036)

- [ ] T031 [P] NavigationSyncAPIClient in apps/mobile/ios/src/services/NavigationSyncAPIClient.swift
- [ ] T032 [P] UIComponentSyncAPIClient in apps/mobile/ios/src/services/UIComponentSyncAPIClient.swift
- [ ] T033 [P] PlatformAdapterAPIClient in apps/mobile/ios/src/services/PlatformAdapterAPIClient.swift
- [ ] T034 [P] NavigationSyncAPIClient in apps/mobile/android/app/src/main/java/com/tchat/services/NavigationSyncAPIClient.kt
- [ ] T035 [P] UIComponentSyncAPIClient in apps/mobile/android/app/src/main/java/com/tchat/services/UIComponentSyncAPIClient.kt
- [ ] T036 [P] PlatformAdapterAPIClient in apps/mobile/android/app/src/main/java/com/tchat/services/PlatformAdapterAPIClient.kt

## Phase 3.5: Core Screen Components - Main Tabs (T037-T046)

- [ ] T037 [P] Chat Tab screen component in apps/mobile/ios/src/screens/ChatTabView.swift
- [ ] T038 [P] Store Tab screen component in apps/mobile/ios/src/screens/StoreTabView.swift
- [ ] T039 [P] Social Tab screen component in apps/mobile/ios/src/screens/SocialTabView.swift
- [ ] T040 [P] Video Tab screen component in apps/mobile/ios/src/screens/VideoTabView.swift
- [ ] T041 [P] More Tab screen component in apps/mobile/ios/src/screens/MoreTabView.swift
- [ ] T042 [P] Chat Tab screen component in apps/mobile/android/app/src/main/java/com/tchat/screens/ChatTabScreen.kt
- [ ] T043 [P] Store Tab screen component in apps/mobile/android/app/src/main/java/com/tchat/screens/StoreTabScreen.kt
- [ ] T044 [P] Social Tab screen component in apps/mobile/android/app/src/main/java/com/tchat/screens/SocialTabScreen.kt
- [ ] T045 [P] Video Tab screen component in apps/mobile/android/app/src/main/java/com/tchat/screens/VideoTabScreen.kt
- [ ] T046 [P] More Tab screen component in apps/mobile/android/app/src/main/java/com/tchat/screens/MoreTabScreen.kt

## Phase 3.6: Navigation Integration (T047-T048)

- [ ] T047 Main tab navigation controller with iOS NavigationStack in apps/mobile/ios/src/navigation/MainTabNavigation.swift
- [ ] T048 Main tab navigation controller with Android Navigation Component in apps/mobile/android/app/src/main/java/com/tchat/navigation/MainTabNavigation.kt

## Phase 3.7: Cross-Platform Synchronization (T049-T050)

- [ ] T049 [P] Cross-platform state synchronization service in apps/mobile/ios/src/services/StateSyncService.swift
- [ ] T050 [P] Cross-platform state synchronization service in apps/mobile/android/app/src/main/java/com/tchat/services/StateSyncService.kt

## Phase 3.8: Platform-Specific Features (T051-T052)

- [ ] T051 [P] iOS platform adapter with haptic feedback, Face ID, share sheet in apps/mobile/ios/src/platform/IOSPlatformAdapter.swift
- [ ] T052 [P] Android platform adapter with Material You, biometrics, sharing intents in apps/mobile/android/app/src/main/java/com/tchat/platform/AndroidPlatformAdapter.kt

## Phase 3.9: Performance & Polish (T053-T055)

- [ ] T053 [P] Performance monitoring and optimization implementation for iOS and Android
- [ ] T054 [P] Accessibility compliance implementation (VoiceOver, TalkBack, Dynamic Type) for both platforms
- [ ] T055 Execute quickstart.md validation scenarios and verify all 6 test scenarios pass

## Dependencies

**Critical Path**:
- Setup (T001-T010) → Contract Tests (T011-T020) → Models (T021-T030) → API Services (T031-T036) → Screen Components (T037-T046) → Navigation (T047-T048) → Sync (T049-T050) → Platform Features (T051-T052) → Polish (T053-T055)

**Blocking Dependencies**:
- T011-T020 (contract tests) MUST complete before T021-T055 (implementation)
- T021-T030 (models) block T031-T036 (API services)
- T031-T036 (API services) block T037-T046 (screen components)
- T037-T046 (screen components) block T047-T048 (navigation)
- T047-T048 (navigation) blocks T049-T050 (sync services)
- T049-T050 (sync services) blocks T051-T052 (platform features)

**Platform Independence**:
- iOS tasks (T002, T004, T006, T008, T014-T016, T018, T021-T025, T031-T033, T037-T041, T047, T049, T051) can run parallel to Android tasks
- Android tasks (T003, T005, T007, T009, T011-T013, T017, T026-T030, T034-T036, T042-T046, T048, T050, T052) can run parallel to iOS tasks

## Parallel Execution Examples

### Foundation Setup (T002-T009)
```bash
# Launch iOS and Android setup together:
Task: "Initialize iOS project with SwiftUI and iOS NavigationStack dependencies"
Task: "Initialize Android project with Jetpack Compose and Navigation Component dependencies"
Task: "Configure iOS project settings (iOS 15+ target, Swift 5.9+, Xcode build)"
Task: "Configure Android project settings (API 24+ target, Kotlin 1.9+, Gradle build)"
```

### Contract Tests (T011-T020)
```bash
# Launch all contract tests in parallel:
Task: "Contract test Navigation Sync API in apps/mobile/android/app/src/test/java/com/tchat/NavigationSyncAPITest.kt"
Task: "Contract test UI Component Sync API in apps/mobile/android/app/src/test/java/com/tchat/UIComponentSyncAPITest.kt"
Task: "Contract test Platform Adapter API in apps/mobile/android/app/src/test/java/com/tchat/PlatformAdapterAPITest.kt"
Task: "Contract test Navigation Sync API in apps/mobile/ios/Tests/NavigationSyncAPITests.swift"
Task: "Contract test UI Component Sync API in apps/mobile/ios/Tests/UIComponentSyncAPITests.swift"
Task: "Contract test Platform Adapter API in apps/mobile/ios/Tests/PlatformAdapterAPITests.swift"
```

### Screen Components (T037-T046)
```bash
# Launch main tab implementations in parallel:
Task: "Chat Tab screen component in apps/mobile/ios/src/screens/ChatTabView.swift"
Task: "Store Tab screen component in apps/mobile/ios/src/screens/StoreTabView.swift"
Task: "Chat Tab screen component in apps/mobile/android/app/src/main/java/com/tchat/screens/ChatTabScreen.kt"
Task: "Store Tab screen component in apps/mobile/android/app/src/main/java/com/tchat/screens/StoreTabScreen.kt"
```

## Notes
- [P] tasks = different files/platforms, no dependencies
- Verify contract tests fail before implementing features (TDD approach)
- iOS and Android development can proceed in parallel throughout
- Commit after each completed task
- Performance targets: 3s launch, 300ms navigation, 60 FPS, <300MB memory
- Accessibility: VoiceOver, TalkBack, Dynamic Type compliance required

## Validation Checklist
*GATE: Checked before task execution*

- [x] All contracts (3) have corresponding tests (T011-T016)
- [x] All entities (5) have model tasks for both platforms (T021-T030)
- [x] All tests (T011-T020) come before implementation (T021-T055)
- [x] Parallel tasks truly independent (iOS vs Android, different files)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] Critical path clearly defined with dependencies
- [x] 55 tasks cover complete implementation scope
- [x] Quickstart validation included (T055)