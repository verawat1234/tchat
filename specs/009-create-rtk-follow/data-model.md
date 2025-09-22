# Data Model: RTK Backend API Integration

## Core Entities

### User
Represents an authenticated user in the system.

```typescript
interface User {
  id: string;
  email: string;
  username: string;
  displayName: string;
  avatar?: string;
  role: UserRole;
  createdAt: string; // ISO 8601
  updatedAt: string; // ISO 8601
  lastLoginAt?: string; // ISO 8601
  preferences?: UserPreferences;
}

enum UserRole {
  ADMIN = 'admin',
  USER = 'user',
  GUEST = 'guest'
}

interface UserPreferences {
  theme: 'light' | 'dark' | 'system';
  language: string;
  notifications: NotificationSettings;
}
```

### Authentication
Authentication tokens and session data.

```typescript
interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  expiresIn: number; // seconds
  tokenType: 'Bearer';
}

interface LoginRequest {
  email: string;
  password: string;
  rememberMe?: boolean;
}

interface LoginResponse {
  user: User;
  tokens: AuthTokens;
}

interface RefreshTokenRequest {
  refreshToken: string;
}

interface RefreshTokenResponse {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}
```

### API Response Wrappers
Standard response formats for consistency.

```typescript
// Success response
interface ApiResponse<T> {
  success: true;
  data: T;
  meta?: ResponseMeta;
}

// Error response
interface ApiError {
  success: false;
  error: {
    code: string;
    message: string;
    details?: Record<string, any>;
    timestamp: string;
  };
}

// Pagination meta
interface ResponseMeta {
  pagination?: PaginationMeta;
  timestamp: string;
  version: string;
}

interface PaginationMeta {
  total: number;
  page: number;
  limit: number;
  totalPages: number;
  hasNext: boolean;
  hasPrev: boolean;
}
```

### Paginated Response
For list endpoints with pagination.

```typescript
interface PaginatedResponse<T> {
  items: T[];
  pagination: {
    cursor?: string;
    nextCursor?: string;
    prevCursor?: string;
    hasMore: boolean;
    total?: number;
  };
}
```

### Request State
Track loading states for UI feedback.

```typescript
interface RequestState {
  isLoading: boolean;
  isSuccess: boolean;
  isError: boolean;
  error?: ApiError;
}

interface CachedData<T> {
  data: T;
  timestamp: number;
  isStale: boolean;
}
```

## Resource Entities (Examples)

### Message
Example resource entity for chat messages.

```typescript
interface Message {
  id: string;
  chatId: string;
  userId: string;
  content: string;
  type: MessageType;
  attachments?: Attachment[];
  createdAt: string;
  updatedAt: string;
  editedAt?: string;
  deletedAt?: string;
  readBy: ReadReceipt[];
}

enum MessageType {
  TEXT = 'text',
  IMAGE = 'image',
  FILE = 'file',
  SYSTEM = 'system'
}

interface Attachment {
  id: string;
  url: string;
  type: string;
  size: number;
  name: string;
}

interface ReadReceipt {
  userId: string;
  readAt: string;
}
```

### Chat
Example resource entity for chat rooms.

```typescript
interface Chat {
  id: string;
  name: string;
  type: ChatType;
  participants: string[]; // User IDs
  lastMessage?: Message;
  unreadCount: number;
  createdAt: string;
  updatedAt: string;
}

enum ChatType {
  DIRECT = 'direct',
  GROUP = 'group',
  CHANNEL = 'channel'
}
```

## State Shape

### Redux Store Structure
```typescript
interface RootState {
  // RTK Query managed
  api: {
    queries: Record<string, CachedData<any>>;
    mutations: Record<string, RequestState>;
    provided: Record<string, string[]>;
    subscriptions: Record<string, number>;
  };

  // Feature slices
  auth: AuthState;
  ui: UIState;
  // ... other feature slices
}

interface AuthState {
  isAuthenticated: boolean;
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  expiresAt: number | null;
}

interface UIState {
  theme: 'light' | 'dark' | 'system';
  sidebarOpen: boolean;
  activeModal: string | null;
  notifications: Notification[];
}
```

## Validation Rules

### User Validation
```typescript
const UserSchema = z.object({
  email: z.string().email(),
  username: z.string().min(3).max(30).regex(/^[a-zA-Z0-9_-]+$/),
  displayName: z.string().min(1).max(100),
  password: z.string().min(8).max(100) // for registration
});
```

### Message Validation
```typescript
const MessageSchema = z.object({
  content: z.string().min(1).max(5000),
  type: z.enum(['text', 'image', 'file', 'system']),
  attachments: z.array(AttachmentSchema).optional()
});
```

## State Transitions

### Authentication Flow
```
INITIAL → LOGGING_IN → AUTHENTICATED
        ↓            ↓
     FAILED ←────── LOGGING_OUT
        ↑               ↓
        └────────── LOGGED_OUT
```

### Request Lifecycle
```
IDLE → PENDING → FULFILLED
     ↓        ↓
     └─→ REJECTED
```

### Cache States
```
EMPTY → LOADING → CACHED → STALE → INVALIDATED
                     ↓        ↓         ↓
                     └────────┴─────> REFETCHING
```

## Error Codes

### Standard Error Codes
```typescript
enum ErrorCode {
  // Authentication
  UNAUTHORIZED = 'UNAUTHORIZED',
  TOKEN_EXPIRED = 'TOKEN_EXPIRED',
  INVALID_CREDENTIALS = 'INVALID_CREDENTIALS',

  // Authorization
  FORBIDDEN = 'FORBIDDEN',
  INSUFFICIENT_PERMISSIONS = 'INSUFFICIENT_PERMISSIONS',

  // Validation
  VALIDATION_ERROR = 'VALIDATION_ERROR',
  INVALID_INPUT = 'INVALID_INPUT',

  // Resources
  NOT_FOUND = 'NOT_FOUND',
  CONFLICT = 'CONFLICT',

  // Rate Limiting
  RATE_LIMITED = 'RATE_LIMITED',

  // Server
  INTERNAL_ERROR = 'INTERNAL_ERROR',
  SERVICE_UNAVAILABLE = 'SERVICE_UNAVAILABLE',

  // Network
  NETWORK_ERROR = 'NETWORK_ERROR',
  TIMEOUT = 'TIMEOUT'
}
```

## Type Guards

### Utility Type Guards
```typescript
function isApiError(error: unknown): error is ApiError {
  return (
    typeof error === 'object' &&
    error !== null &&
    'success' in error &&
    error.success === false
  );
}

function isAuthTokens(data: unknown): data is AuthTokens {
  return (
    typeof data === 'object' &&
    data !== null &&
    'accessToken' in data &&
    'refreshToken' in data
  );
}
```

## Notes

1. All timestamps use ISO 8601 format for consistency
2. IDs are strings to support various backend ID strategies (UUID, MongoDB ObjectId, etc.)
3. Enums use string values for better debugging and API compatibility
4. Optional fields use TypeScript's optional property syntax (?)
5. Arrays default to empty arrays, not null/undefined
6. All entities include createdAt and updatedAt for audit trails