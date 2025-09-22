# Research: RTK Backend API Integration

## Overview
Research findings for implementing Redux Toolkit (RTK) with RTK Query for backend API integration in the Tchat application.

## Key Decisions

### 1. State Management Architecture
**Decision**: RTK Query with createApi for all API interactions
**Rationale**:
- Built-in caching with automatic cache invalidation
- Automatic re-fetching on focus/reconnect
- TypeScript code generation from OpenAPI schemas
- Optimistic updates with rollback support
- Reduced boilerplate vs traditional Redux

**Alternatives Considered**:
- Plain Redux + fetch: Too much boilerplate
- Redux + Axios + thunks: Manual cache management required
- TanStack Query: Good alternative but RTK Query integrates better with Redux

### 2. Authentication Method
**Decision**: JWT tokens with refresh token rotation
**Rationale**:
- Stateless authentication suitable for distributed systems
- Secure token storage in httpOnly cookies for refresh tokens
- Access tokens in memory (not localStorage)
- Automatic token refresh via RTK Query middleware

**Alternatives Considered**:
- Session cookies: Less suitable for mobile/multi-platform
- OAuth 2.0: Overkill for internal API
- Basic Auth: Insufficient security

### 3. Caching Strategy
**Decision**: Tag-based cache invalidation with 60-second default TTL
**Rationale**:
- Tags allow precise invalidation of related data
- 60-second TTL balances freshness with performance
- Manual refetch available for critical updates
- Optimistic updates for perceived performance

**Implementation**:
```typescript
// Tag types: 'User', 'Message', 'Chat', etc.
// Provides: ['User'], invalidates: ['User']
```

### 4. Error Handling Pattern
**Decision**: Centralized error handling with retry logic
**Rationale**:
- Consistent error messages across the app
- Exponential backoff for network errors (3 retries)
- User-friendly error translations
- Toast notifications for user feedback

**Error Categories**:
- Network errors: Auto-retry with backoff
- 401: Trigger token refresh
- 403: Show permission error
- 422: Show validation errors
- 500+: Show generic error with report option

### 5. Pagination Approach
**Decision**: Cursor-based pagination with infinite scroll
**Rationale**:
- Better performance for real-time data
- Handles data insertions/deletions gracefully
- Works well with RTK Query's merge function
- Supports bi-directional scrolling

**Implementation**:
```typescript
// Query params: { cursor?: string, limit: number }
// Response: { data: T[], nextCursor?: string, hasMore: boolean }
```

### 6. Optimistic Updates
**Decision**: Optimistic updates with automatic rollback
**Rationale**:
- Instant UI feedback for better UX
- Automatic rollback on server rejection
- Preserves UI state during rollback
- Works with RTK Query's onQueryStarted

**Supported Operations**:
- Create: Show immediately, remove on failure
- Update: Show changes, revert on failure
- Delete: Hide immediately, restore on failure

## API Structure

### Base Configuration
```typescript
baseUrl: process.env.VITE_API_URL || 'http://localhost:3001/api'
timeout: 30000 // 30 seconds
headers: {
  'Content-Type': 'application/json'
}
```

### Standard Endpoints Pattern
```
GET    /api/{resource}         - List with pagination
GET    /api/{resource}/{id}    - Get single item
POST   /api/{resource}         - Create new item
PUT    /api/{resource}/{id}    - Full update
PATCH  /api/{resource}/{id}    - Partial update
DELETE /api/{resource}/{id}    - Delete item
```

### Authentication Endpoints
```
POST   /api/auth/login         - Login with credentials
POST   /api/auth/logout        - Logout and invalidate tokens
POST   /api/auth/refresh       - Refresh access token
GET    /api/auth/me            - Get current user
```

## Performance Optimizations

### 1. Request Deduplication
- RTK Query automatically deduplicates identical concurrent requests
- Subscribers share the same cached result

### 2. Selective Refetching
- refetchOnFocus: true for active data
- refetchOnReconnect: true for critical data
- Manual refetch triggers for user-initiated updates

### 3. Bundle Splitting
- Lazy load API slices for feature-specific endpoints
- Core API (auth, user) loaded upfront
- Feature APIs loaded on-demand

### 4. Response Normalization
- Normalize nested data structures
- Use entity adapters for collections
- Reduce redundant data storage

## Testing Strategy

### 1. MSW (Mock Service Worker)
- Intercept API calls at network level
- Test real RTK Query behavior
- Support both browser and Node environments

### 2. Contract Tests
- Validate request/response schemas
- Ensure frontend/backend contract alignment
- Run as part of CI pipeline

### 3. Integration Tests
- Test complete user flows
- Verify cache invalidation logic
- Test error scenarios and recovery

## Security Considerations

### 1. Token Storage
- Access tokens: In-memory only (Redux store)
- Refresh tokens: httpOnly, secure, sameSite cookies
- No sensitive data in localStorage

### 2. CSRF Protection
- CSRF tokens for state-changing operations
- Double submit cookie pattern
- SameSite cookie attribute

### 3. Request Validation
- Client-side validation before API calls
- Schema validation with Zod
- Sanitize user inputs

## Migration Path

### Phase 1: Core Setup
1. Install RTK and dependencies
2. Configure store with RTK Query
3. Create base API service

### Phase 2: Authentication
1. Implement auth endpoints
2. Add token refresh middleware
3. Secure route protection

### Phase 3: Resource APIs
1. Migrate existing API calls to RTK Query
2. Implement caching strategies
3. Add optimistic updates

### Phase 4: Testing & Optimization
1. Add comprehensive tests
2. Optimize bundle size
3. Performance monitoring

## Dependencies

### Required Packages
```json
{
  "@reduxjs/toolkit": "^2.0.0",
  "react-redux": "^9.0.0",
  "msw": "^2.0.0" // for testing
}
```

### Optional Enhancements
```json
{
  "redux-persist": "^6.0.0", // for persistence
  "zod": "^3.22.0" // for validation
}
```

## Conclusion

RTK Query provides a robust solution for API integration with minimal boilerplate. The tag-based caching system, automatic re-fetching, and optimistic updates will significantly improve the user experience while maintaining code simplicity and type safety.