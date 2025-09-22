# Tchat Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-01-27

## Active Technologies
- Go 1.21+ (microservices backend architecture) + Go standard library, gRPC, protocol buffers, JWT authentication, message brokers (Kafka/RabbitMQ), WebSocket libraries (004-create-backend-spec)
- PostgreSQL (primary), ScyllaDB/Cassandra (timelines), Redis (caching), CDN (static assets) (004-create-backend-spec)
- TypeScript 5.3.0 + React 18.3.1 (web frontend testing) + Vitest, Testing Library (@testing-library/react), @testing-library/jest-dom, @storybook/test-runner, Playwright (E2E) (005-create-test-for)
- Test fixtures, snapshots, coverage reports (filesystem-based) (005-create-test-for)
- Swift 5.9+ (iOS), Kotlin 1.9+ (Android) + SwiftUI + Combine (iOS), Jetpack Compose + Coroutines (Android) (006-implement-native-mobile)
- Local state management, cross-platform sync with existing web backend (006-implement-native-mobile)
- Swift 5.9+ for iOS, Kotlin 1.9+ for Android, TypeScript 5.3.0 for shared web logic + SwiftUI (iOS), Jetpack Compose (Android), React 18.3.1 (web reference), Navigation frameworks (iOS NavigationStack, Android Navigation Component) (007-create-spec-of)
- Local state management (UserDefaults/SharedPreferences), cross-platform sync via existing API, CoreData (iOS) and Room (Android) for offline caching (007-create-spec-of)
- Go 1.22+ (microservices backend) + Gin HTTP framework, GORM ORM, JWT authentication, WebSocket libraries (008-complete-microservices-architecture)
- PostgreSQL (primary data), ScyllaDB (messages), Redis (cache/sessions), Kafka (events) (008-complete-microservices-architecture)
- TypeScript 5.3.0 / JavaScript ES2020+ + Redux Toolkit 2.0+, RTK Query, React 18.3.1 (009-create-rtk-follow)
- RTK store for client state, localStorage for persistence (009-create-rtk-follow)
- TypeScript 5.3.0 with React 18.3.1 + Redux Toolkit 2.0+, RTK Query, React-Redux 9.2.0 (010-now-i-think)
- RTK Query cache with backend persistence (existing RTK infrastructure) (010-now-i-think)
- **Dynamic Content Management System** + 12 RTK Query endpoints, localStorage fallback, performance optimization (<200ms), comprehensive E2E testing with Playwright (011-dynamic-content-system)
- Advanced caching with tag-based invalidation, optimistic updates, error recovery middleware, enterprise monitoring (011-dynamic-content-system)
- Go 1.22+ (microservices backend architecture) + estify/suite, testify/mock, testify/assert, go-sqlmock, httptest, dockertes (011-complete-test-coverage-spec)
- PostgreSQL (primary), ScyllaDB (messages), Redis (cache/sessions) (011-complete-test-coverage-spec)

### Web Platform
- TypeScript 5.3.0, React 18.3.1 + Vite 6.3.5, Radix UI components, TailwindCSS v4, Framer Motion 11.0.0
- Redux Toolkit 2.0+ with RTK Query for API state management, Redux Persist for offline support
- **Dynamic Content Management**: 12 RTK Query endpoints, localStorage fallback service, performance-optimized (<200ms load times)
- Authentication: JWT tokens with automatic refresh, secure token storage
- Caching: Advanced tag-based invalidation, optimistic updates with rollback, error recovery middleware
- Testing: Vitest, Testing Library, Storybook, MSW for API mocking, Playwright E2E testing, contract-driven TDD approach

### Mobile Platform (Native)
- **iOS**: Swift 5.9+ with SwiftUI, Combine, Alamofire 5.8+, Kingfisher 7.9+, Swift Package Manager
- **Android**: Kotlin 1.9.23 with Jetpack Compose, Material3, Coroutines, Gradle 8.4+
- **Architecture**: Design token-based system, cross-platform state synchronization, TDD approach
- **Navigation**: 5-tab architecture (Chat/Store/Social/Video/More) with platform-native patterns
- **Testing**: Contract tests, XCTest (iOS), JUnit + Espresso (Android)

## Project Structure
```
apps/
├── web/                    # React web application
├── mobile/
│   ├── ios/               # Native iOS Swift app
│   │   ├── Sources/       # Swift source code
│   │   │   ├── Components/    # UI components (TchatButton, TchatInput, TchatCard)
│   │   │   ├── DesignSystem/  # Design tokens (Colors, Typography, Spacing)
│   │   │   ├── Navigation/    # TabNavigationView
│   │   │   └── State/         # AppState, StateSyncManager, PersistenceManager
│   │   ├── Tests/         # iOS test suites
│   │   └── Package.swift  # Swift package configuration
│   └── android/           # Native Android Kotlin app
│       ├── app/src/main/java/com/tchat/
│       │   ├── components/     # Compose UI components
│       │   ├── designsystem/   # Design tokens (Colors, Typography, Spacing)
│       │   ├── navigation/     # Tab navigation composables
│       │   └── state/          # State management and sync
│       └── app/src/test/   # Android test suites
backend/
tests/
tools/
```

## Commands

### Web Development
```bash
# Web app development
npm run dev              # Start development server
npm run build           # Build for production
npm test                # Run tests
npm run storybook       # Start Storybook

# Component analysis
npm run analyze-components  # Validate component structure

# Content Management Testing
npm run test:content        # Run content management tests
npm run test:e2e:content    # Run content E2E performance tests
npm run test:fallback       # Test localStorage fallback system
```

### Mobile Development
```bash
# iOS development
cd apps/mobile/ios
swift build             # Build iOS app
swift test              # Run iOS tests
swiftlint              # Code linting
xcodebuild -scheme TchatApp  # Xcode build

# Android development
cd apps/mobile/android
./gradlew build         # Build Android app
./gradlew test          # Run Android tests
./gradlew ktlintCheck   # Code linting
./gradlew assembleDebug # Build debug APK
```

## Code Style

### Web (TypeScript/React)
- TypeScript 5.3.0, React 18.3.1: Follow standard conventions
- Radix UI patterns with TailwindCSS v4 styling
- Component-first architecture with design system consistency

### iOS (Swift/SwiftUI)
- SwiftLint configuration for code consistency
- Design token-based styling system
- Combine for reactive programming
- Platform-native navigation patterns

### Android (Kotlin/Compose)
- ktlint for code formatting
- Material3 design system integration
- Jetpack Compose UI patterns
- Coroutines for asynchronous operations

## Architecture Highlights

### Cross-Platform Design System
- **Design Tokens**: Translated from TailwindCSS v4 to native equivalents
- **Component Library**: Platform-native implementations with >95% visual consistency
- **5 Component Variants**: TchatButton (primary/secondary/ghost/destructive/outline), TchatInput (validation states), TchatCard (4 variants)

### Dynamic Content Management System
- **12 RTK Query Endpoints**: Complete CRUD operations, versioning, bulk updates, synchronization, category management
- **Advanced Fallback System**: localStorage persistence with automatic error recovery and offline support
- **Performance Optimization**: <200ms content load budget, Core Web Vitals monitoring, memory management
- **Enterprise Features**: Tag-based cache invalidation, optimistic updates, error recovery middleware, maintenance tasks
- **Type Safety**: Comprehensive TypeScript interfaces, request/response validation, content type definitions
- **Real-time Sync**: Incremental synchronization, conflict resolution, deleted content tracking

### State Management
- **Web-Native Sync**: Real-time state synchronization between web and mobile
- **Content State**: Centralized content management with Redux Toolkit and RTK Query
- **Persistence**: Secure storage (iOS Keychain, Android EncryptedSharedPreferences), localStorage fallback service
- **AppState Pattern**: Centralized state management with platform-specific optimizations

### Testing Strategy
- **TDD Approach**: Contract tests drive implementation, 50+ comprehensive test suites for content management
- **Cross-Platform Testing**: API contracts, design token validation, state sync testing
- **Performance Testing**: Playwright E2E testing with Core Web Vitals, memory usage, network efficiency validation
- **Content Testing**: 12 endpoint test suites, fallback service testing, error recovery validation
- **Platform-Specific**: XCTest (iOS), JUnit + Espresso (Android), Vitest (Web), Playwright (E2E)

## Recent Changes
- 011-complete-test-coverage-spec: Added Go 1.22+ (microservices backend architecture) + estify/suite, testify/mock, testify/assert, go-sqlmock, httptest, dockertes
- 011-dynamic-content-system: **COMPLETED** - Dynamic Content Management System Implementation
  - **Complete RTK Architecture**: 12 comprehensive endpoints (getContentItems, getContentItem, getContentByCategory, getContentCategories, getContentVersions, syncContent, createContentItem, updateContentItem, publishContent, archiveContent, bulkUpdateContent, revertContentVersion)
  - **Advanced Fallback System**: localStorage-based fallback service with automatic error recovery, offline support, and intelligent caching
  - **Performance Optimization**: <200ms content load budget achieved, Core Web Vitals monitoring, memory management (<100MB mobile, <500MB desktop)
  - **Enterprise Reliability**: Tag-based cache invalidation, optimistic updates with rollback, error recovery middleware, automated maintenance tasks
  - **Comprehensive Testing**: 50+ test suites covering all endpoints, E2E Playwright performance testing (882 lines), Core Web Vitals validation
  - **Type Safety**: Complete TypeScript interfaces, request/response validation, content type definitions for text, rich_text, config, and image types
  - **Real-time Features**: Incremental synchronization, conflict resolution, version control, bulk operations support
- 010-now-i-think: Added TypeScript 5.3.0 with React 18.3.1 + Redux Toolkit 2.0+, RTK Query, React-Redux 9.2.0
  - Redux store with RTK Query middleware and persistence
  - Complete API service layer with auth, users, messages, chats endpoints
  - JWT authentication with automatic token refresh middleware
  - Tag-based cache invalidation with optimistic updates
  - Error handling middleware with user-friendly notifications
  - Request retry logic with exponential backoff
  - Redux Persist for offline support with secure token storage
  - Prefetching service for performance optimization
  - Request deduplication and advanced caching strategies
  - Contract-driven TDD with MSW mocking (50 tests implemented)

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
