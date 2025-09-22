// User types
export interface User {
  id: string;
  email: string;
  username: string;
  displayName: string;
  avatar?: string;
  role: UserRole;
  createdAt: string;
  updatedAt: string;
  lastLoginAt?: string;
  preferences?: UserPreferences;
}

export enum UserRole {
  ADMIN = 'admin',
  USER = 'user',
  GUEST = 'guest'
}

export interface UserPreferences {
  theme: 'light' | 'dark' | 'system';
  language: string;
  notifications: NotificationSettings;
}

export interface NotificationSettings {
  email: boolean;
  push: boolean;
  inApp: boolean;
}

// Authentication types
export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
  tokenType: 'Bearer';
}

export interface LoginRequest {
  email: string;
  password: string;
  rememberMe?: boolean;
}

export interface LoginResponse {
  user: User;
  tokens: AuthTokens;
}

export interface RefreshTokenRequest {
  refreshToken: string;
}

export interface RefreshTokenResponse {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

// Message types
export interface Message {
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
  readBy?: ReadReceipt[];
}

export enum MessageType {
  TEXT = 'text',
  IMAGE = 'image',
  FILE = 'file',
  SYSTEM = 'system'
}

export interface Attachment {
  id: string;
  url: string;
  type: string;
  size: number;
  name: string;
}

export interface ReadReceipt {
  userId: string;
  readAt: string;
}

export interface CreateMessageRequest {
  chatId: string;
  content: string;
  type: MessageType;
  attachments?: Attachment[];
}

// Chat types
export interface Chat {
  id: string;
  name: string;
  type: ChatType;
  participants: string[];
  lastMessage?: Message;
  unreadCount: number;
  createdAt: string;
  updatedAt: string;
}

export enum ChatType {
  DIRECT = 'direct',
  GROUP = 'group',
  CHANNEL = 'channel'
}

// API Response types
export interface ApiResponse<T> {
  success: true;
  data: T;
  meta?: ResponseMeta;
}

export interface ApiError {
  success: false;
  error: {
    code: string;
    message: string;
    details?: Record<string, any>;
    timestamp: string;
  };
}

export interface ResponseMeta {
  pagination?: PaginationMeta;
  timestamp: string;
  version: string;
}

export interface PaginationMeta {
  total: number;
  page: number;
  limit: number;
  totalPages: number;
  hasNext: boolean;
  hasPrev: boolean;
}

export interface PaginatedResponse<T> {
  items: T[];
  pagination: {
    cursor?: string;
    nextCursor?: string;
    prevCursor?: string;
    hasMore: boolean;
    total?: number;
  };
}

// Request state
export interface RequestState {
  isLoading: boolean;
  isSuccess: boolean;
  isError: boolean;
  error?: ApiError;
}

// Update user request
export interface UpdateUserRequest {
  displayName?: string;
  avatar?: string;
  preferences?: Partial<UserPreferences>;
}

// Error codes
export enum ErrorCode {
  UNAUTHORIZED = 'UNAUTHORIZED',
  TOKEN_EXPIRED = 'TOKEN_EXPIRED',
  INVALID_CREDENTIALS = 'INVALID_CREDENTIALS',
  FORBIDDEN = 'FORBIDDEN',
  INSUFFICIENT_PERMISSIONS = 'INSUFFICIENT_PERMISSIONS',
  VALIDATION_ERROR = 'VALIDATION_ERROR',
  INVALID_INPUT = 'INVALID_INPUT',
  NOT_FOUND = 'NOT_FOUND',
  CONFLICT = 'CONFLICT',
  RATE_LIMITED = 'RATE_LIMITED',
  INTERNAL_ERROR = 'INTERNAL_ERROR',
  SERVICE_UNAVAILABLE = 'SERVICE_UNAVAILABLE',
  NETWORK_ERROR = 'NETWORK_ERROR',
  TIMEOUT = 'TIMEOUT'
}

// Type guards
export function isApiError(error: unknown): error is ApiError {
  return (
    typeof error === 'object' &&
    error !== null &&
    'success' in error &&
    (error as any).success === false
  );
}

export function isAuthTokens(data: unknown): data is AuthTokens {
  return (
    typeof data === 'object' &&
    data !== null &&
    'accessToken' in data &&
    'refreshToken' in data
  );
}