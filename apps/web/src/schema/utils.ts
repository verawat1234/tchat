/**
 * Utility Types and Common Helpers
 *
 * Shared types, validators, and helper functions for the schema system
 * Includes API response wrappers, pagination, sorting, and validation utilities
 */

import { UUID, Timestamp, Currency, Locale, CountryCode } from './schema';

// =============================================================================
// API RESPONSE WRAPPERS
// =============================================================================

export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: ApiError;
  meta?: ResponseMeta;
  timestamp: Timestamp;
  requestId?: string;
}

export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, any>;
  field?: string;
  suggestion?: string;
  documentation_url?: string;
}

export interface ResponseMeta {
  pagination?: PaginationMeta;
  sorting?: SortingMeta;
  filtering?: FilteringMeta;
  timing?: TimingMeta;
  cache?: CacheMeta;
}

// =============================================================================
// PAGINATION
// =============================================================================

export interface PaginationParams {
  page?: number;
  limit?: number;
  offset?: number;
  cursor?: string;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
  hasNextPage: boolean;
  hasPreviousPage: boolean;
  nextCursor?: string;
  previousCursor?: string;
}

export interface PaginatedResponse<T> {
  items: T[];
  pagination: PaginationMeta;
}

// =============================================================================
// SORTING & FILTERING
// =============================================================================

export interface SortingParams {
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
  sortFields?: SortField[];
}

export interface SortField {
  field: string;
  direction: 'asc' | 'desc';
  nullsFirst?: boolean;
}

export interface SortingMeta {
  sortBy: string;
  sortOrder: 'asc' | 'desc';
  availableFields: string[];
}

export interface FilteringParams {
  filters?: Record<string, any>;
  search?: string;
  dateRange?: DateRange;
  priceRange?: PriceRange;
  status?: string[];
  category?: string[];
  tags?: string[];
}

export interface FilteringMeta {
  appliedFilters: Record<string, any>;
  availableFilters: FilterOption[];
  totalMatchingItems: number;
}

export interface FilterOption {
  field: string;
  type: 'string' | 'number' | 'date' | 'boolean' | 'select' | 'multiselect';
  label: string;
  options?: { value: any; label: string; count?: number }[];
  range?: { min: any; max: any };
}

export interface DateRange {
  start: Timestamp;
  end: Timestamp;
  timezone?: string;
}

export interface PriceRange {
  min: number;
  max: number;
  currency: Currency;
}

// =============================================================================
// TIMING & PERFORMANCE
// =============================================================================

export interface TimingMeta {
  requestTime: number; // milliseconds
  dbQueryTime?: number;
  cacheTime?: number;
  totalTime: number;
  queries?: number;
}

export interface CacheMeta {
  hit: boolean;
  key?: string;
  ttl?: number; // seconds
  age?: number; // seconds
  region?: string;
}

// =============================================================================
// VALIDATION TYPES
// =============================================================================

export interface ValidationResult {
  isValid: boolean;
  errors: ValidationError[];
  warnings: ValidationWarning[];
}

export interface ValidationError {
  field: string;
  code: string;
  message: string;
  value?: any;
  constraint?: string;
}

export interface ValidationWarning {
  field: string;
  code: string;
  message: string;
  suggestion?: string;
}

// =============================================================================
// SEARCH TYPES
// =============================================================================

export interface SearchParams {
  query: string;
  type?: string[];
  filters?: FilteringParams;
  sorting?: SortingParams;
  pagination?: PaginationParams;
  highlight?: boolean;
  fuzzy?: boolean;
  boost?: SearchBoost[];
}

export interface SearchBoost {
  field: string;
  factor: number;
}

export interface SearchResult<T = any> {
  item: T;
  score: number;
  highlights?: SearchHighlight[];
  explanation?: string;
}

export interface SearchHighlight {
  field: string;
  value: string;
  matched: string[];
}

export interface SearchResponse<T = any> {
  results: SearchResult<T>[];
  total: number;
  maxScore: number;
  took: number; // milliseconds
  suggestions?: SearchSuggestion[];
  facets?: SearchFacet[];
}

export interface SearchSuggestion {
  text: string;
  score: number;
  category?: string;
}

export interface SearchFacet {
  field: string;
  values: { value: string; count: number }[];
}

// =============================================================================
// GEOLOCATION TYPES
// =============================================================================

export interface Coordinates {
  latitude: number;
  longitude: number;
  altitude?: number;
  accuracy?: number;
  heading?: number;
  speed?: number;
}

export interface BoundingBox {
  northEast: Coordinates;
  southWest: Coordinates;
}

export interface DistanceFilter {
  center: Coordinates;
  radius: number;
  unit: 'km' | 'miles' | 'meters';
}

// =============================================================================
// FILE & MEDIA TYPES
// =============================================================================

export interface FileUpload {
  file: File | Blob;
  fileName?: string;
  mimeType?: string;
  size?: number;
  metadata?: FileMetadata;
}

export interface FileMetadata {
  width?: number;
  height?: number;
  duration?: number;
  bitrate?: number;
  fps?: number;
  format?: string;
  checksum?: string;
  uploadedBy?: UUID;
  uploadedAt?: Timestamp;
}

export interface ImageResolution {
  width: number;
  height: number;
  quality?: number;
  format?: 'jpeg' | 'png' | 'webp' | 'avif';
}

export interface VideoTranscoding {
  resolution: string;
  bitrate: number;
  fps: number;
  codec: string;
  format: string;
}

// =============================================================================
// LOCALIZATION TYPES
// =============================================================================

export interface LocalizedString {
  [locale: string]: string;
}

export interface LocalizedContent {
  [locale: string]: {
    title?: string;
    description?: string;
    content?: string;
    keywords?: string[];
    metadata?: Record<string, any>;
  };
}

export interface CurrencyAmount {
  amount: number;
  currency: Currency;
  formatted?: string;
  exchangeRate?: number;
  baseAmount?: number;
  baseCurrency?: Currency;
}

export interface Translation {
  key: string;
  locale: Locale;
  value: string;
  pluralForms?: Record<string, string>;
  context?: string;
  namespace?: string;
}

// =============================================================================
// AUDIT & TRACKING TYPES
// =============================================================================

export interface AuditLog {
  id: UUID;
  entityType: string;
  entityId: UUID;
  action: AuditAction;
  userId?: UUID;
  changes: AuditChange[];
  metadata: AuditMetadata;
  createdAt: Timestamp;
}

export type AuditAction = 'create' | 'update' | 'delete' | 'view' | 'export' | 'login' | 'logout';

export interface AuditChange {
  field: string;
  oldValue?: any;
  newValue?: any;
  type: 'added' | 'modified' | 'removed';
}

export interface AuditMetadata {
  ipAddress?: string;
  userAgent?: string;
  source: 'web' | 'mobile' | 'api' | 'system';
  sessionId?: string;
  requestId?: string;
  correlationId?: string;
}

export interface TrackingEvent {
  id: UUID;
  userId?: UUID;
  sessionId: string;
  event: string;
  category: string;
  action: string;
  label?: string;
  value?: number;
  properties: Record<string, any>;
  context: EventContext;
  createdAt: Timestamp;
}

export interface EventContext {
  page?: string;
  referrer?: string;
  userAgent?: string;
  device?: DeviceInfo;
  location?: GeoLocation;
  campaign?: CampaignInfo;
}

export interface DeviceInfo {
  type: 'desktop' | 'mobile' | 'tablet';
  os: string;
  osVersion?: string;
  browser?: string;
  browserVersion?: string;
  screen?: { width: number; height: number };
  language?: string;
  timezone?: string;
}

export interface GeoLocation {
  country?: CountryCode;
  region?: string;
  city?: string;
  coordinates?: Coordinates;
  timezone?: string;
  isp?: string;
}

export interface CampaignInfo {
  source?: string;
  medium?: string;
  campaign?: string;
  term?: string;
  content?: string;
}

// =============================================================================
// WEBHOOK TYPES
// =============================================================================

export interface WebhookEvent<T = any> {
  id: UUID;
  type: string;
  data: T;
  metadata: WebhookMetadata;
  createdAt: Timestamp;
}

export interface WebhookMetadata {
  version: string;
  source: string;
  retryCount?: number;
  signature?: string;
  correlationId?: string;
}

export interface WebhookDelivery {
  id: UUID;
  webhookId: UUID;
  eventId: UUID;
  url: string;
  httpStatus?: number;
  responseBody?: string;
  responseHeaders?: Record<string, string>;
  deliveredAt?: Timestamp;
  failedAt?: Timestamp;
  retryCount: number;
  nextRetryAt?: Timestamp;
  error?: string;
}

// =============================================================================
// RATE LIMITING TYPES
// =============================================================================

export interface RateLimit {
  limit: number;
  remaining: number;
  reset: Timestamp;
  window: number; // seconds
  type: 'user' | 'ip' | 'api_key' | 'endpoint';
}

export interface RateLimitConfig {
  requests: number;
  window: number; // seconds
  burstLimit?: number;
  skipFailedRequests?: boolean;
  skipSuccessfulRequests?: boolean;
  keyGenerator?: (req: any) => string;
}

// =============================================================================
// FEATURE FLAGS TYPES
// =============================================================================

export interface FeatureFlag {
  key: string;
  name: string;
  description?: string;
  isEnabled: boolean;
  rolloutPercentage?: number;
  targetUsers?: UUID[];
  targetCountries?: CountryCode[];
  targetPlatforms?: string[];
  startDate?: Timestamp;
  endDate?: Timestamp;
  metadata?: Record<string, any>;
}

export interface FeatureFlagEvaluation {
  key: string;
  isEnabled: boolean;
  variant?: string;
  reason: string;
  metadata?: Record<string, any>;
}

// =============================================================================
// UTILITY FUNCTIONS
// =============================================================================

/**
 * Type guard functions
 */
export const isUUID = (value: any): value is UUID => {
  return typeof value === 'string' && /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(value);
};

export const isTimestamp = (value: any): value is Timestamp => {
  return typeof value === 'string' && !isNaN(Date.parse(value));
};

export const isCurrency = (value: any): value is Currency => {
  return ['THB', 'SGD', 'IDR', 'MYR', 'PHP', 'VND', 'USD'].includes(value);
};

export const isLocale = (value: any): value is Locale => {
  return ['th-TH', 'id-ID', 'ms-MY', 'vi-VN', 'en-US'].includes(value);
};

export const isCountryCode = (value: any): value is CountryCode => {
  return ['TH', 'ID', 'MY', 'VN', 'SG', 'PH'].includes(value);
};

/**
 * Default pagination parameters
 */
export const DEFAULT_PAGINATION: PaginationParams = {
  page: 1,
  limit: 20,
  offset: 0
};

/**
 * Default sorting parameters
 */
export const DEFAULT_SORTING: SortingParams = {
  sortBy: 'createdAt',
  sortOrder: 'desc'
};

/**
 * Common error codes
 */
export const ERROR_CODES = {
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  NOT_FOUND: 'NOT_FOUND',
  UNAUTHORIZED: 'UNAUTHORIZED',
  FORBIDDEN: 'FORBIDDEN',
  RATE_LIMITED: 'RATE_LIMITED',
  INTERNAL_ERROR: 'INTERNAL_ERROR',
  SERVICE_UNAVAILABLE: 'SERVICE_UNAVAILABLE',
  INVALID_INPUT: 'INVALID_INPUT',
  DUPLICATE_ENTRY: 'DUPLICATE_ENTRY',
  PAYMENT_FAILED: 'PAYMENT_FAILED',
  INSUFFICIENT_BALANCE: 'INSUFFICIENT_BALANCE',
  KYC_REQUIRED: 'KYC_REQUIRED',
  FEATURE_DISABLED: 'FEATURE_DISABLED'
} as const;

/**
 * HTTP status codes
 */
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  ACCEPTED: 202,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  METHOD_NOT_ALLOWED: 405,
  CONFLICT: 409,
  UNPROCESSABLE_ENTITY: 422,
  TOO_MANY_REQUESTS: 429,
  INTERNAL_SERVER_ERROR: 500,
  BAD_GATEWAY: 502,
  SERVICE_UNAVAILABLE: 503,
  GATEWAY_TIMEOUT: 504
} as const;

/**
 * Common regular expressions
 */
export const REGEX_PATTERNS = {
  EMAIL: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  PHONE: /^\+?[1-9]\d{1,14}$/,
  UUID: /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i,
  PASSWORD: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/,
  URL: /^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$/,
  SLUG: /^[a-z0-9]+(?:-[a-z0-9]+)*$/,
  HEX_COLOR: /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/,
  BASE64: /^[A-Za-z0-9+/]*={0,2}$/
} as const;

/**
 * File size limits (in bytes)
 */
export const FILE_SIZE_LIMITS = {
  AVATAR: 5 * 1024 * 1024, // 5MB
  IMAGE: 10 * 1024 * 1024, // 10MB
  VIDEO: 100 * 1024 * 1024, // 100MB
  AUDIO: 50 * 1024 * 1024, // 50MB
  DOCUMENT: 25 * 1024 * 1024, // 25MB
  VOICE_MESSAGE: 10 * 1024 * 1024 // 10MB
} as const;

/**
 * Content type mappings
 */
export const CONTENT_TYPES = {
  'image/jpeg': 'jpg',
  'image/png': 'png',
  'image/gif': 'gif',
  'image/webp': 'webp',
  'image/svg+xml': 'svg',
  'video/mp4': 'mp4',
  'video/webm': 'webm',
  'video/quicktime': 'mov',
  'audio/mpeg': 'mp3',
  'audio/wav': 'wav',
  'audio/ogg': 'ogg',
  'application/pdf': 'pdf',
  'application/msword': 'doc',
  'application/vnd.openxmlformats-officedocument.wordprocessingml.document': 'docx',
  'application/vnd.ms-excel': 'xls',
  'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet': 'xlsx',
  'text/plain': 'txt',
  'text/csv': 'csv'
} as const;