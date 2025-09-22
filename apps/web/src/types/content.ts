/**
 * Content Management System Types
 *
 * This file defines the TypeScript types and interfaces for the dynamic content
 * management system, supporting flexible content storage, localization, and
 * version control with RTK Query integration.
 */

// Re-export enhanced ContentCategory types for backwards compatibility
export {
  type PermissionAction,
  type CategoryPermissions,
  type CategoryHierarchy,
  type CategoryAccessContext,
  type EnhancedContentCategory,
  isEnhancedContentCategory,
  isValidCategoryPermissions,
  isValidPermissionOverrides,
  isValidCategoryHierarchy,
  isValidCategoryMetadata,
  createContentCategory,
  createCategoryPermissions,
  computeCategoryHierarchy,
  getCategoryParents,
  getCategoryChildren,
  updateCategoryHierarchies,
  hasPermission,
  getEffectivePermissions,
  canAccessCategoryTree,
  validateCategoryHierarchy,
  validateCategoryMove,
} from './content-category';

/**
 * Content type enumeration defining the supported content data types
 */
export enum ContentType {
  /** Plain text content */
  TEXT = 'text',
  /** HTML or Markdown formatted content */
  RICH_TEXT = 'rich_text',
  /** Image URL content with metadata */
  IMAGE_URL = 'image_url',
  /** Configuration values (boolean, number, string, object) */
  CONFIG = 'config',
  /** Localized content with multiple language support */
  TRANSLATION = 'translation'
}

/**
 * Content lifecycle status enumeration
 */
export enum ContentStatus {
  /** Work in progress, not visible to end users */
  DRAFT = 'draft',
  /** Live content visible to end users */
  PUBLISHED = 'published',
  /** Deprecated content no longer in use */
  ARCHIVED = 'archived'
}

/**
 * Permissions interface for content category access control
 */
export interface CategoryPermissions {
  /** User roles that can view content in this category */
  read: string[];
  /** User roles that can edit content in this category */
  write: string[];
  /** User roles that can publish content in this category */
  publish: string[];
}

/**
 * Content category for organizing and grouping related content items
 */
export interface ContentCategory {
  /** Unique identifier (e.g., "navigation", "errors", "help") */
  id: string;
  /** Human-readable name for display */
  name: string;
  /** Purpose and usage description */
  description: string;
  /** Parent category ID for hierarchical organization */
  parentId?: string;
  /** Access control permissions */
  permissions: CategoryPermissions;
}

/**
 * Plain text content with optional length constraints
 */
export interface TextContent {
  type: 'text';
  /** The text content */
  value: string;
  /** Maximum allowed character length */
  maxLength?: number;
}

/**
 * Rich text content supporting HTML or Markdown
 */
export interface RichTextContent {
  type: 'rich_text';
  /** The formatted content (HTML or Markdown) */
  value: string;
  /** Content format specification */
  format: 'html' | 'markdown';
  /** Allowed HTML tags for security (whitelist) */
  allowedTags?: string[];
}

/**
 * Image content with URL and metadata
 */
export interface ImageContent {
  type: 'image_url';
  /** Image URL */
  url: string;
  /** Alternative text for accessibility */
  alt: string;
  /** Image width in pixels */
  width?: number;
  /** Image height in pixels */
  height?: number;
  /** Image format (e.g., 'png', 'jpg', 'webp') */
  format?: string;
}

/**
 * Configuration content supporting various data types
 */
export interface ConfigContent {
  type: 'config';
  /** Configuration value (boolean, number, string, or object) */
  value: boolean | number | string | object;
  /** JSON schema for validation (optional) */
  schema?: Record<string, any>;
}

/**
 * Translation content with multi-language support
 */
export interface TranslationContent {
  type: 'translation';
  /** Mapping of locale codes to translated text */
  values: Record<string, string>;
  /** Default locale for fallback */
  defaultLocale: string;
}

/**
 * Union type for all supported content value types
 */
export type ContentValue =
  | TextContent
  | RichTextContent
  | ImageContent
  | ConfigContent
  | TranslationContent;

/**
 * User action type enumeration for audit tracking
 */
export enum AuditAction {
  CREATE = 'create',
  UPDATE = 'update',
  PUBLISH = 'publish',
  ARCHIVE = 'archive',
  RESTORE = 'restore',
  DELETE = 'delete',
  REVERT = 'revert',
  BULK_UPDATE = 'bulk_update',
  BULK_PUBLISH = 'bulk_publish',
  BULK_ARCHIVE = 'bulk_archive',
  STATUS_CHANGE = 'status_change',
  TAG_ADDED = 'tag_added',
  TAG_REMOVED = 'tag_removed',
  METADATA_UPDATE = 'metadata_update'
}

/**
 * Metadata change severity levels
 */
export enum ChangeSeverity {
  MINOR = 'minor',
  MAJOR = 'major',
  CRITICAL = 'critical'
}

/**
 * Content annotation interface for editorial notes and reviews
 */
export interface ContentAnnotation {
  /** Unique annotation ID */
  id: string;
  /** Annotation type */
  type: 'note' | 'review' | 'approval' | 'warning' | 'todo';
  /** Annotation content */
  content: string;
  /** User who created the annotation */
  createdBy: string;
  /** Creation timestamp */
  createdAt: string;
  /** Whether annotation is resolved */
  resolved: boolean;
  /** Resolution timestamp */
  resolvedAt?: string;
  /** User who resolved the annotation */
  resolvedBy?: string;
  /** Associated tags for categorization */
  tags?: string[];
  /** Priority level */
  priority?: 'low' | 'medium' | 'high' | 'urgent';
}

/**
 * Audit trail entry for tracking all content changes
 */
export interface AuditEntry {
  /** Unique audit entry ID */
  id: string;
  /** Content item ID this audit entry relates to */
  contentId: string;
  /** Action performed */
  action: AuditAction;
  /** User who performed the action */
  performedBy: string;
  /** Timestamp of the action */
  performedAt: string;
  /** Version after the action */
  version: number;
  /** Previous version (if applicable) */
  previousVersion?: number;
  /** Change description */
  description: string;
  /** Severity of the change */
  severity: ChangeSeverity;
  /** IP address of the user */
  ipAddress?: string;
  /** User agent string */
  userAgent?: string;
  /** Additional context data */
  context?: Record<string, any>;
  /** Changes made (field-level diff) */
  changes?: Record<string, {
    before: any;
    after: any;
  }>;
}

/**
 * Metadata validation rules interface
 */
export interface MetadataValidationRules {
  /** Required tags pattern */
  requiredTags?: string[];
  /** Maximum number of tags allowed */
  maxTags?: number;
  /** Allowed tag patterns (regex) */
  allowedTagPatterns?: string[];
  /** Required annotations for certain actions */
  requiredAnnotations?: Record<AuditAction, string[]>;
  /** Maximum annotation length */
  maxAnnotationLength?: number;
  /** Custom validation functions */
  customValidators?: Array<(metadata: ContentMetadata) => ValidationResult>;
}

/**
 * Validation result interface
 */
export interface ValidationResult {
  /** Whether validation passed */
  valid: boolean;
  /** Validation errors */
  errors: Array<{
    field: string;
    message: string;
    code: string;
  }>;
  /** Validation warnings */
  warnings: Array<{
    field: string;
    message: string;
    code: string;
  }>;
}

/**
 * Comprehensive metadata for tracking content changes and management information
 */
export interface ContentMetadata {
  /** ISO timestamp of creation */
  createdAt: string;
  /** User ID who created the content */
  createdBy: string;
  /** ISO timestamp of last update */
  updatedAt: string;
  /** User ID who last updated the content */
  updatedBy: string;
  /** Incremental version number */
  version: number;
  /** Searchable tags for content organization */
  tags: string[];
  /** Editorial notes for content management */
  notes?: string;
  /** Content annotations for editorial workflow */
  annotations: ContentAnnotation[];
  /** Complete audit trail */
  auditTrail: AuditEntry[];
  /** Content workflow state */
  workflow?: {
    /** Current workflow stage */
    stage: 'draft' | 'review' | 'approved' | 'published' | 'archived';
    /** Assigned reviewer */
    assignedTo?: string;
    /** Review deadline */
    reviewDeadline?: string;
    /** Approval status */
    approvalStatus?: 'pending' | 'approved' | 'rejected';
    /** Approval timestamp */
    approvedAt?: string;
    /** User who approved */
    approvedBy?: string;
  };
  /** Content governance metadata */
  governance?: {
    /** Content owner */
    owner: string;
    /** Data classification */
    classification: 'public' | 'internal' | 'confidential' | 'restricted';
    /** Retention policy */
    retentionPolicy?: {
      /** Retention period in days */
      retentionDays: number;
      /** Auto-archive date */
      autoArchiveDate?: string;
      /** Auto-delete date */
      autoDeleteDate?: string;
    };
    /** Compliance tags */
    complianceTags?: string[];
    /** Legal hold status */
    legalHold?: boolean;
  };
  /** Quality metrics */
  quality?: {
    /** Content quality score (0-100) */
    score?: number;
    /** Last quality assessment */
    lastAssessment?: string;
    /** Quality issues */
    issues?: Array<{
      type: string;
      description: string;
      severity: ChangeSeverity;
    }>;
  };
  /** Usage analytics metadata */
  analytics?: {
    /** View count */
    viewCount?: number;
    /** Last accessed timestamp */
    lastAccessed?: string;
    /** Access frequency */
    accessFrequency?: number;
    /** Performance metrics */
    performance?: {
      /** Average load time in ms */
      avgLoadTime?: number;
      /** Error rate percentage */
      errorRate?: number;
    };
  };
  /** Content relationships */
  relationships?: {
    /** Parent content ID */
    parentId?: string;
    /** Child content IDs */
    childIds?: string[];
    /** Related content IDs */
    relatedIds?: string[];
    /** Dependencies */
    dependencies?: string[];
  };
  /** Localization metadata */
  localization?: {
    /** Source language */
    sourceLanguage: string;
    /** Available translations */
    translations: Record<string, {
      /** Translation status */
      status: 'pending' | 'in_progress' | 'completed' | 'outdated';
      /** Translator ID */
      translatedBy?: string;
      /** Translation timestamp */
      translatedAt?: string;
      /** Translation quality score */
      qualityScore?: number;
    }>;
    /** Translation memory references */
    tmReferences?: string[];
  };
  /** Custom metadata fields */
  custom?: Record<string, any>;
}

/**
 * Core content entity representing any piece of information displayed to users
 */
export interface ContentItem {
  /** Unique identifier (e.g., "header.title", "error.network") */
  id: string;
  /** Semantic key for content item (e.g., "welcome.title", "chat.subtitle") */
  key: string;
  /** Category ID reference */
  categoryId: string;
  /** Semantic grouping category */
  category: ContentCategory;
  /** Data type and rendering information */
  type: ContentType;
  /** The actual content data */
  value: ContentValue;
  /** Management and tracking information */
  metadata: ContentMetadata;
  /** Current lifecycle state */
  status: ContentStatus;
  /** Optional tags for content organization */
  tags?: string[];
  /** Optional editorial notes */
  notes?: string;
  /** Optional change log for version history */
  changeLog?: string;
}

/**
 * Change types for version diff calculation
 */
export enum VersionChangeType {
  /** Content was added */
  ADDED = 'added',
  /** Content was modified */
  MODIFIED = 'modified',
  /** Content was deleted */
  DELETED = 'deleted',
  /** Content was moved/renamed */
  MOVED = 'moved',
  /** Metadata or status change */
  METADATA = 'metadata'
}

/**
 * Version comparison result enumeration
 */
export enum VersionComparison {
  /** First version is older than second */
  OLDER = -1,
  /** Versions are equal */
  EQUAL = 0,
  /** First version is newer than second */
  NEWER = 1
}

/**
 * Detailed change information for version diffs
 */
export interface VersionChange {
  /** Type of change */
  type: VersionChangeType;
  /** Field that was changed (for MODIFIED type) */
  field?: string;
  /** Previous value (for MODIFIED/DELETED) */
  oldValue?: any;
  /** New value (for MODIFIED/ADDED) */
  newValue?: any;
  /** Path to the changed property (dot notation) */
  path?: string;
  /** Human-readable description of the change */
  description: string;
}

/**
 * Diff result between two content versions
 */
export interface VersionDiff {
  /** Version that changes are compared from */
  fromVersion: number;
  /** Version that changes are compared to */
  toVersion: number;
  /** Array of changes between versions */
  changes: VersionChange[];
  /** Whether the diff contains breaking changes */
  hasBreakingChanges: boolean;
  /** Summary statistics about the changes */
  summary: {
    /** Total number of changes */
    totalChanges: number;
    /** Number of additions */
    additions: number;
    /** Number of modifications */
    modifications: number;
    /** Number of deletions */
    deletions: number;
  };
  /** ISO timestamp when diff was calculated */
  calculatedAt: string;
}

/**
 * Enhanced version history tracking for content items
 */
export interface ContentVersion {
  /** Unique version identifier */
  id: string;
  /** Parent content item ID */
  contentId: string;
  /** Version number (semantic versioning supported) */
  version: number;
  /** Semantic version string (e.g., "1.2.3") */
  semanticVersion?: string;
  /** Content value at this version */
  value: ContentValue;
  /** Version-specific metadata */
  metadata: ContentMetadata;
  /** Description of changes made in this version */
  changeLog: string;
  /** Type of change that created this version */
  changeType: VersionChangeType;
  /** Previous version number (null for initial version) */
  previousVersion?: number;
  /** Branch name for version control workflows */
  branch?: string;
  /** Parent version for merge operations */
  parentVersion?: number;
  /** Merge information if this version is a merge */
  mergeInfo?: {
    /** Source branch that was merged */
    sourceBranch: string;
    /** Target branch that received the merge */
    targetBranch: string;
    /** User who performed the merge */
    mergedBy: string;
    /** Timestamp of merge operation */
    mergedAt: string;
    /** Merge commit message */
    mergeMessage?: string;
  };
  /** Approval information for published versions */
  approval?: {
    /** User who approved the version */
    approvedBy: string;
    /** Timestamp of approval */
    approvedAt: string;
    /** Approval notes or comments */
    approvalNotes?: string;
    /** Required approval level (e.g., 'content-editor', 'admin') */
    approvalLevel: string;
  };
  /** Checksum for version integrity validation */
  checksum?: string;
  /** File size in bytes for large content tracking */
  size?: number;
  /** Version expiration date for time-limited content */
  expiresAt?: string;
  /** Rollback information if this version was created by revert */
  rollbackInfo?: {
    /** Version that this was rolled back from */
    rolledBackFrom: number;
    /** User who performed the rollback */
    rolledBackBy: string;
    /** Timestamp of rollback operation */
    rolledBackAt: string;
    /** Reason for rollback */
    rollbackReason?: string;
  };
}

/**
 * Content slice state for Redux store
 */
export interface ContentState {
  /** UI state for language selection */
  selectedLanguage: string;
  /** User content preferences */
  contentPreferences: {
    /** Whether to show draft content */
    showDrafts: boolean;
    /** Use compact view for content lists */
    compactView: boolean;
  };
  /** ISO timestamp of last successful sync */
  lastSyncTime: string;
  /** Current synchronization status */
  syncStatus: 'idle' | 'syncing' | 'error';
  /** Whether using local fallback content */
  fallbackMode: boolean;
  /** Local fallback content cache */
  fallbackContent: Record<string, ContentValue>;
}

/**
 * Request interface for creating new content items
 */
export interface CreateContentItemRequest {
  /** Content item ID */
  id: string;
  /** Category ID */
  categoryId: string;
  /** Content type */
  type: ContentType;
  /** Content value */
  value: ContentValue;
  /** Optional tags */
  tags?: string[];
  /** Optional editorial notes */
  notes?: string;
}

/**
 * Request interface for updating existing content items
 */
export interface UpdateContentItemRequest {
  /** Updated content value */
  value?: ContentValue;
  /** Updated status */
  status?: ContentStatus;
  /** Updated tags */
  tags?: string[];
  /** Editorial notes for the update */
  notes?: string;
  /** Description of changes for version history */
  changeLog?: string;
  /** Expected version for optimistic locking */
  expectedVersion?: number;
}

/**
 * Request interface for bulk content operations
 */
export interface BulkContentOperationRequest {
  /** Content item IDs to operate on */
  contentIds: string[];
  /** Operation to perform */
  operation: 'publish' | 'archive' | 'delete';
  /** Optional change log for version history */
  changeLog?: string;
}

/**
 * Single item update request for bulk operations
 */
export interface BulkUpdateItemRequest {
  /** Content item ID to update */
  id: string;
  /** Update request data */
  update: UpdateContentItemRequest;
}

/**
 * Request interface for bulk content updates
 */
export interface BulkUpdateContentRequest {
  /** Array of content items to update */
  items: BulkUpdateItemRequest[];
  /** Whether to process updates atomically (all or none) */
  atomic?: boolean;
  /** Whether to continue processing on individual item failures */
  continueOnError?: boolean;
  /** Global change log for all updates */
  changeLog?: string;
}

/**
 * Result for a single item in bulk update operation
 */
export interface BulkUpdateItemResult {
  /** Content item ID */
  id: string;
  /** Whether the update was successful */
  success: boolean;
  /** Updated content item (on success) */
  item?: ContentItem;
  /** Error details (on failure) */
  error?: {
    code: string;
    message: string;
    field?: string;
  };
}

/**
 * Response interface for bulk update operations
 */
export interface BulkUpdateContentResponse {
  /** Results for each item */
  results: BulkUpdateItemResult[];
  /** Overall operation success status */
  success: boolean;
  /** Number of successfully updated items */
  successCount: number;
  /** Number of failed updates */
  errorCount: number;
  /** Total processing time in milliseconds */
  processingTime: number;
  /** Whether operation was rolled back (atomic mode only) */
  rolledBack?: boolean;
}

/**
 * Request interface for reverting content to a previous version
 */
export interface RevertContentVersionRequest {
  /** Target version number to revert to */
  targetVersion: number;
  /** Reason for the reversion */
  reason?: string;
  /** Optional change log for the reversion */
  changeLog?: string;
  /** Whether to force revert even if conflicts exist */
  forceRevert?: boolean;
  /** Expected current version for optimistic locking */
  expectedCurrentVersion?: number;
}

/**
 * Response interface for content version reversion
 */
export interface RevertContentVersionResponse {
  /** The updated content item after reversion */
  contentItem: ContentItem;
  /** The new version created by the revert operation */
  newVersion: ContentVersion;
  /** Version that was reverted to */
  revertedToVersion: number;
  /** Previous version before revert */
  previousVersion: number;
  /** Metadata about the revert operation */
  revertMetadata: {
    /** User who performed the revert */
    revertedBy: string;
    /** Timestamp of revert operation */
    revertedAt: string;
    /** Reference to the version this was reverted from */
    revertedFrom: number;
    /** Reason provided for the revert */
    reason?: string;
  };
}

/**
 * Query parameters for content item filtering and searching
 */
export interface ContentQueryParams {
  /** Filter by category ID */
  categoryId?: string;
  /** Filter by content type */
  type?: ContentType;
  /** Filter by status */
  status?: ContentStatus;
  /** Search in content values */
  search?: string;
  /** Filter by tags */
  tags?: string[];
  /** Filter by language (for translations) */
  language?: string;
  /** Sort field */
  sortBy?: 'id' | 'updatedAt' | 'createdAt' | 'version';
  /** Sort direction */
  sortOrder?: 'asc' | 'desc';
  /** Pagination page number */
  page?: number;
  /** Items per page */
  limit?: number;
}

/**
 * Response interface for paginated content results
 */
export interface PaginatedContentResponse {
  /** Content items */
  items: ContentItem[];
  /** Pagination metadata */
  pagination: {
    /** Current page */
    page: number;
    /** Items per page */
    limit: number;
    /** Total number of items */
    total: number;
    /** Total number of pages */
    totalPages: number;
    /** Whether there are more pages */
    hasNext: boolean;
    /** Whether there are previous pages */
    hasPrev: boolean;
  };
}

/**
 * Type guard to check if a value is a ContentItem with runtime validation
 * @param value - The value to check
 * @returns true if the value is a valid ContentItem
 */
export function isContentItem(value: unknown): value is ContentItem {
  if (typeof value !== 'object' || value === null) {
    return false;
  }

  const item = value as Partial<ContentItem>;

  return (
    typeof item.id === 'string' &&
    typeof item.key === 'string' &&
    typeof item.categoryId === 'string' &&
    isContentCategory(item.category) &&
    Object.values(ContentType).includes(item.type as ContentType) &&
    isValidContentValue(item.value) &&
    isValidContentMetadata(item.metadata) &&
    Object.values(ContentStatus).includes(item.status as ContentStatus)
  );
}

/**
 * Validates if a value is a valid ContentValue with proper structure
 * @param value - The content value to validate
 * @returns true if the value is a valid ContentValue
 */
export function isValidContentValue(value: unknown): value is ContentValue {
  if (typeof value !== 'object' || value === null) {
    return false;
  }

  const content = value as Partial<ContentValue>;

  switch (content.type) {
    case 'text':
      return isTextContent(content as ContentValue);
    case 'rich_text':
      return isRichTextContent(content as ContentValue);
    case 'image_url':
      return isImageContent(content as ContentValue);
    case 'config':
      return isConfigContent(content as ContentValue);
    case 'translation':
      return isTranslationContent(content as ContentValue);
    default:
      return false;
  }
}

/**
 * Validates if a value is valid ContentMetadata
 * @param value - The metadata to validate
 * @returns true if the value is valid ContentMetadata
 */
export function isValidContentMetadata(value: unknown): value is ContentMetadata {
  if (typeof value !== 'object' || value === null) {
    return false;
  }

  const metadata = value as Partial<ContentMetadata>;

  return (
    typeof metadata.createdAt === 'string' &&
    typeof metadata.createdBy === 'string' &&
    typeof metadata.updatedAt === 'string' &&
    typeof metadata.updatedBy === 'string' &&
    typeof metadata.version === 'number' &&
    metadata.version >= 1 &&
    Array.isArray(metadata.tags) &&
    Array.isArray(metadata.annotations) &&
    Array.isArray(metadata.auditTrail) &&
    (metadata.notes === undefined || typeof metadata.notes === 'string')
  );
}

/**
 * Type guard to check if a value is a ContentCategory
 */
export function isContentCategory(value: unknown): value is ContentCategory {
  return (
    typeof value === 'object' &&
    value !== null &&
    'id' in value &&
    'name' in value &&
    'description' in value &&
    'permissions' in value
  );
}

/**
 * Type guard to check if a content value is TextContent with validation
 * @param value - The content value to check
 * @returns true if the value is valid TextContent
 */
export function isTextContent(value: ContentValue): value is TextContent {
  if (value.type !== 'text') return false;

  const textContent = value as TextContent;
  return (
    typeof textContent.value === 'string' &&
    (textContent.maxLength === undefined ||
     (typeof textContent.maxLength === 'number' && textContent.maxLength > 0))
  );
}

/**
 * Type guard to check if a content value is RichTextContent with validation
 * @param value - The content value to check
 * @returns true if the value is valid RichTextContent
 */
export function isRichTextContent(value: ContentValue): value is RichTextContent {
  if (value.type !== 'rich_text') return false;

  const richTextContent = value as RichTextContent;
  return (
    typeof richTextContent.value === 'string' &&
    ['html', 'markdown'].includes(richTextContent.format) &&
    (richTextContent.allowedTags === undefined ||
     Array.isArray(richTextContent.allowedTags))
  );
}

/**
 * Type guard to check if a content value is ImageContent with validation
 * @param value - The content value to check
 * @returns true if the value is valid ImageContent
 */
export function isImageContent(value: ContentValue): value is ImageContent {
  if (value.type !== 'image_url') return false;

  const imageContent = value as ImageContent;
  return (
    typeof imageContent.url === 'string' &&
    typeof imageContent.alt === 'string' &&
    (imageContent.width === undefined || typeof imageContent.width === 'number') &&
    (imageContent.height === undefined || typeof imageContent.height === 'number') &&
    (imageContent.format === undefined || typeof imageContent.format === 'string')
  );
}

/**
 * Type guard to check if a content value is ConfigContent with validation
 * @param value - The content value to check
 * @returns true if the value is valid ConfigContent
 */
export function isConfigContent(value: ContentValue): value is ConfigContent {
  if (value.type !== 'config') return false;

  const configContent = value as ConfigContent;
  const validValueTypes = ['boolean', 'number', 'string', 'object'];
  return (
    validValueTypes.includes(typeof configContent.value) &&
    (configContent.schema === undefined ||
     typeof configContent.schema === 'object')
  );
}

/**
 * Type guard to check if a content value is TranslationContent with validation
 * @param value - The content value to check
 * @returns true if the value is valid TranslationContent
 */
export function isTranslationContent(value: ContentValue): value is TranslationContent {
  if (value.type !== 'translation') return false;

  const translationContent = value as TranslationContent;
  return (
    typeof translationContent.values === 'object' &&
    translationContent.values !== null &&
    typeof translationContent.defaultLocale === 'string' &&
    translationContent.defaultLocale in translationContent.values &&
    Object.values(translationContent.values).every(v => typeof v === 'string')
  );
}

/**
 * Utility type for content ID pattern validation
 * Pattern: {category}.{subcategory}.{key}
 * Examples: "navigation.header.title", "error.network.timeout"
 */
export type ContentIdPattern = `${string}.${string}.${string}`;

/**
 * Utility type for extracting category from content ID
 */
export type ExtractCategory<T extends string> = T extends `${infer Category}.${string}.${string}`
  ? Category
  : never;

/**
 * Utility type for extracting subcategory from content ID
 */
export type ExtractSubcategory<T extends string> = T extends `${string}.${infer Subcategory}.${string}`
  ? Subcategory
  : never;

/**
 * Utility type for extracting key from content ID
 */
export type ExtractKey<T extends string> = T extends `${string}.${string}.${infer Key}`
  ? Key
  : never;

/**
 * Type guard to check if a value is a ContentVersion
 */
export function isContentVersion(value: unknown): value is ContentVersion {
  return (
    typeof value === 'object' &&
    value !== null &&
    'id' in value &&
    'contentId' in value &&
    'version' in value &&
    'value' in value &&
    'metadata' in value &&
    'changeLog' in value &&
    'changeType' in value
  );
}

/**
 * Type guard to check if a value is a VersionDiff
 */
export function isVersionDiff(value: unknown): value is VersionDiff {
  return (
    typeof value === 'object' &&
    value !== null &&
    'fromVersion' in value &&
    'toVersion' in value &&
    'changes' in value &&
    'hasBreakingChanges' in value &&
    'summary' in value
  );
}

/**
 * Type guard to check if a value is a VersionChange
 */
export function isVersionChange(value: unknown): value is VersionChange {
  return (
    typeof value === 'object' &&
    value !== null &&
    'type' in value &&
    'description' in value &&
    Object.values(VersionChangeType).includes((value as any).type)
  );
}

// =============================================================================
// VERSION MANAGEMENT UTILITIES
// =============================================================================

/**
 * Version comparison utility function
 * Compares two version numbers and returns comparison result
 */
export function compareVersions(versionA: number, versionB: number): VersionComparison {
  if (versionA < versionB) return VersionComparison.OLDER;
  if (versionA > versionB) return VersionComparison.NEWER;
  return VersionComparison.EQUAL;
}

/**
 * Sort versions in ascending or descending order
 */
export function sortVersions(
  versions: ContentVersion[],
  order: 'asc' | 'desc' = 'desc'
): ContentVersion[] {
  return [...versions].sort((a, b) => {
    const comparison = compareVersions(a.version, b.version);
    return order === 'asc' ? comparison : -comparison;
  });
}

/**
 * Get the latest version from an array of versions
 */
export function getLatestVersion(versions: ContentVersion[]): ContentVersion | null {
  if (versions.length === 0) return null;
  return sortVersions(versions, 'desc')[0];
}

/**
 * Get a specific version by version number
 */
export function getVersionByNumber(
  versions: ContentVersion[],
  versionNumber: number
): ContentVersion | null {
  return versions.find(v => v.version === versionNumber) || null;
}

/**
 * Validate version integrity using checksum
 */
export function validateVersionChecksum(version: ContentVersion): boolean {
  if (!version.checksum) return true; // No checksum to validate

  // Simple checksum calculation (in real implementation, use proper hash)
  const content = JSON.stringify({
    value: version.value,
    metadata: version.metadata,
    changeLog: version.changeLog
  });

  // Basic string hash (replace with proper crypto hash in production)
  let hash = 0;
  for (let i = 0; i < content.length; i++) {
    const char = content.charCodeAt(i);
    hash = ((hash << 5) - hash) + char;
    hash = hash & hash; // Convert to 32-bit integer
  }

  return version.checksum === hash.toString();
}

/**
 * Calculate diff between two content versions
 */
export function calculateVersionDiff(
  fromVersion: ContentVersion,
  toVersion: ContentVersion
): VersionDiff {
  const changes: VersionChange[] = [];
  let hasBreakingChanges = false;

  // Compare content values
  if (JSON.stringify(fromVersion.value) !== JSON.stringify(toVersion.value)) {
    changes.push({
      type: VersionChangeType.MODIFIED,
      field: 'value',
      oldValue: fromVersion.value,
      newValue: toVersion.value,
      path: 'value',
      description: `Content value changed from version ${fromVersion.version} to ${toVersion.version}`
    });

    // Check if content type changed (breaking change)
    if (fromVersion.value.type !== toVersion.value.type) {
      hasBreakingChanges = true;
    }
  }

  // Compare metadata
  const metadataFields: (keyof ContentMetadata)[] = [
    'createdAt', 'createdBy', 'updatedAt', 'updatedBy', 'version', 'tags', 'notes'
  ];

  for (const field of metadataFields) {
    const oldValue = fromVersion.metadata[field];
    const newValue = toVersion.metadata[field];

    if (JSON.stringify(oldValue) !== JSON.stringify(newValue)) {
      changes.push({
        type: VersionChangeType.METADATA,
        field: `metadata.${field}`,
        oldValue,
        newValue,
        path: `metadata.${field}`,
        description: `Metadata field '${field}' changed`
      });
    }
  }

  // Compare change logs
  if (fromVersion.changeLog !== toVersion.changeLog) {
    changes.push({
      type: VersionChangeType.MODIFIED,
      field: 'changeLog',
      oldValue: fromVersion.changeLog,
      newValue: toVersion.changeLog,
      path: 'changeLog',
      description: 'Change log updated'
    });
  }

  // Calculate summary statistics
  const summary = {
    totalChanges: changes.length,
    additions: changes.filter(c => c.type === VersionChangeType.ADDED).length,
    modifications: changes.filter(c => c.type === VersionChangeType.MODIFIED).length,
    deletions: changes.filter(c => c.type === VersionChangeType.DELETED).length
  };

  return {
    fromVersion: fromVersion.version,
    toVersion: toVersion.version,
    changes,
    hasBreakingChanges,
    summary,
    calculatedAt: new Date().toISOString()
  };
}

/**
 * Calculate diff between multiple versions (version history)
 */
export function calculateVersionHistory(
  versions: ContentVersion[]
): VersionDiff[] {
  const sortedVersions = sortVersions(versions, 'asc');
  const diffs: VersionDiff[] = [];

  for (let i = 1; i < sortedVersions.length; i++) {
    const diff = calculateVersionDiff(sortedVersions[i - 1], sortedVersions[i]);
    diffs.push(diff);
  }

  return diffs;
}

/**
 * Check if a version can be safely reverted to
 */
export function canRevertToVersion(
  targetVersion: ContentVersion,
  currentVersion: ContentVersion,
  allVersions: ContentVersion[]
): { canRevert: boolean; reasons: string[] } {
  const reasons: string[] = [];
  let canRevert = true;

  // Check if target version exists
  if (!targetVersion) {
    reasons.push('Target version does not exist');
    canRevert = false;
  }

  // Check if target version is not the current version
  if (targetVersion.version === currentVersion.version) {
    reasons.push('Target version is the same as current version');
    canRevert = false;
  }

  // Check if target version is not in the future
  if (targetVersion.version > currentVersion.version) {
    reasons.push('Cannot revert to a future version');
    canRevert = false;
  }

  // Check for breaking changes in intermediate versions
  const intermediateVersions = allVersions.filter(
    v => v.version > targetVersion.version && v.version <= currentVersion.version
  );

  for (const version of intermediateVersions) {
    const diff = calculateVersionDiff(targetVersion, version);
    if (diff.hasBreakingChanges) {
      reasons.push(`Breaking changes detected in version ${version.version}`);
      // Don't set canRevert to false for breaking changes, just warn
    }
  }

  // Check if version has expired
  if (targetVersion.expiresAt && new Date(targetVersion.expiresAt) < new Date()) {
    reasons.push('Target version has expired');
    canRevert = false;
  }

  return { canRevert, reasons };
}

/**
 * Create a revert version from a target version
 */
export function createRevertVersion(
  contentId: string,
  targetVersion: ContentVersion,
  currentVersion: ContentVersion,
  revertReason?: string
): Omit<ContentVersion, 'id'> {
  const now = new Date().toISOString();
  const newVersionNumber = currentVersion.version + 1;

  return {
    contentId,
    version: newVersionNumber,
    value: targetVersion.value, // Use target version's content
    metadata: {
      ...targetVersion.metadata,
      version: newVersionNumber,
      updatedAt: now,
      // Preserve original creator info from target version
      createdAt: targetVersion.metadata.createdAt,
      createdBy: targetVersion.metadata.createdBy
    },
    changeLog: `Reverted to version ${targetVersion.version}${revertReason ? `: ${revertReason}` : ''}`,
    changeType: VersionChangeType.MODIFIED,
    previousVersion: currentVersion.version,
    rollbackInfo: {
      rolledBackFrom: currentVersion.version,
      rolledBackBy: '', // Should be filled in by the calling code
      rolledBackAt: now,
      rollbackReason: revertReason
    }
  };
}

/**
 * Generate version statistics for content history analysis
 */
export function generateVersionStatistics(
  versions: ContentVersion[]
): {
  totalVersions: number;
  versionsByChangeType: Record<VersionChangeType, number>;
  averageTimeBetweenVersions: number; // in milliseconds
  mostActiveContributor: string;
  contributorStats: Record<string, number>;
  sizeGrowth: {
    initialSize: number;
    currentSize: number;
    growthPercentage: number;
  };
  oldestVersion: ContentVersion | null;
  newestVersion: ContentVersion | null;
} {
  if (versions.length === 0) {
    return {
      totalVersions: 0,
      versionsByChangeType: {
        [VersionChangeType.ADDED]: 0,
        [VersionChangeType.MODIFIED]: 0,
        [VersionChangeType.DELETED]: 0,
        [VersionChangeType.MOVED]: 0,
        [VersionChangeType.METADATA]: 0
      },
      averageTimeBetweenVersions: 0,
      mostActiveContributor: '',
      contributorStats: {},
      sizeGrowth: {
        initialSize: 0,
        currentSize: 0,
        growthPercentage: 0
      },
      oldestVersion: null,
      newestVersion: null
    };
  }

  const sortedVersions = sortVersions(versions, 'asc');
  const oldestVersion = sortedVersions[0];
  const newestVersion = sortedVersions[sortedVersions.length - 1];

  // Count versions by change type
  const versionsByChangeType: Record<VersionChangeType, number> = {
    [VersionChangeType.ADDED]: 0,
    [VersionChangeType.MODIFIED]: 0,
    [VersionChangeType.DELETED]: 0,
    [VersionChangeType.MOVED]: 0,
    [VersionChangeType.METADATA]: 0
  };

  // Count contributor activity
  const contributorStats: Record<string, number> = {};

  for (const version of versions) {
    versionsByChangeType[version.changeType]++;

    const contributor = version.metadata.updatedBy;
    contributorStats[contributor] = (contributorStats[contributor] || 0) + 1;
  }

  // Find most active contributor
  const mostActiveContributor = Object.entries(contributorStats)
    .sort(([, a], [, b]) => b - a)[0]?.[0] || '';

  // Calculate average time between versions
  let totalTimeDiff = 0;
  for (let i = 1; i < sortedVersions.length; i++) {
    const prevTime = new Date(sortedVersions[i - 1].metadata.updatedAt).getTime();
    const currentTime = new Date(sortedVersions[i].metadata.updatedAt).getTime();
    totalTimeDiff += currentTime - prevTime;
  }
  const averageTimeBetweenVersions = sortedVersions.length > 1
    ? totalTimeDiff / (sortedVersions.length - 1)
    : 0;

  // Calculate size growth
  const initialSize = oldestVersion.size || 0;
  const currentSize = newestVersion.size || 0;
  const growthPercentage = initialSize > 0
    ? ((currentSize - initialSize) / initialSize) * 100
    : 0;

  return {
    totalVersions: versions.length,
    versionsByChangeType,
    averageTimeBetweenVersions,
    mostActiveContributor,
    contributorStats,
    sizeGrowth: {
      initialSize,
      currentSize,
      growthPercentage
    },
    oldestVersion,
    newestVersion
  };
}

/**
 * Filter versions by date range
 */
export function filterVersionsByDateRange(
  versions: ContentVersion[],
  startDate: string,
  endDate: string
): ContentVersion[] {
  const start = new Date(startDate).getTime();
  const end = new Date(endDate).getTime();

  return versions.filter(version => {
    const versionTime = new Date(version.metadata.updatedAt).getTime();
    return versionTime >= start && versionTime <= end;
  });
}

/**
 * Filter versions by contributor
 */
export function filterVersionsByContributor(
  versions: ContentVersion[],
  contributorId: string
): ContentVersion[] {
  return versions.filter(version =>
    version.metadata.createdBy === contributorId ||
    version.metadata.updatedBy === contributorId
  );
}

/**
 * Filter versions by change type
 */
export function filterVersionsByChangeType(
  versions: ContentVersion[],
  changeType: VersionChangeType
): ContentVersion[] {
  return versions.filter(version => version.changeType === changeType);
}

/**
 * Get version timeline for visualization
 */
export function getVersionTimeline(
  versions: ContentVersion[]
): Array<{
  version: number;
  timestamp: string;
  contributor: string;
  changeType: VersionChangeType;
  description: string;
}> {
  return sortVersions(versions, 'asc').map(version => ({
    version: version.version,
    timestamp: version.metadata.updatedAt,
    contributor: version.metadata.updatedBy,
    changeType: version.changeType,
    description: version.changeLog
  }));
}

/**
 * Validate version sequence integrity
 */
export function validateVersionSequence(
  versions: ContentVersion[]
): { isValid: boolean; issues: string[] } {
  const issues: string[] = [];
  const sortedVersions = sortVersions(versions, 'asc');

  // Check for missing versions
  for (let i = 1; i < sortedVersions.length; i++) {
    const current = sortedVersions[i];
    const previous = sortedVersions[i - 1];

    if (current.version !== previous.version + 1) {
      issues.push(
        `Gap in version sequence: ${previous.version} → ${current.version}`
      );
    }

    // Check if previous version reference is correct
    if (current.previousVersion && current.previousVersion !== previous.version) {
      issues.push(
        `Incorrect previous version reference in version ${current.version}: ` +
        `expected ${previous.version}, got ${current.previousVersion}`
      );
    }
  }

  // Check for duplicate versions
  const versionNumbers = versions.map(v => v.version);
  const uniqueVersions = new Set(versionNumbers);
  if (versionNumbers.length !== uniqueVersions.size) {
    issues.push('Duplicate version numbers detected');
  }

  // Validate checksums
  for (const version of versions) {
    if (!validateVersionChecksum(version)) {
      issues.push(`Checksum validation failed for version ${version.version}`);
    }
  }

  return {
    isValid: issues.length === 0,
    issues
  };
}

// ====================================
// CONTENT METADATA UTILITIES
// ====================================

/**
 * Type guard to check if a value is a ContentAnnotation
 */
export function isContentAnnotation(value: unknown): value is ContentAnnotation {
  if (typeof value !== 'object' || value === null) {
    return false;
  }

  const annotation = value as Partial<ContentAnnotation>;

  return (
    typeof annotation.id === 'string' &&
    ['note', 'review', 'approval', 'warning', 'todo'].includes(annotation.type as string) &&
    typeof annotation.content === 'string' &&
    typeof annotation.createdBy === 'string' &&
    typeof annotation.createdAt === 'string' &&
    typeof annotation.resolved === 'boolean' &&
    (annotation.tags === undefined || Array.isArray(annotation.tags))
  );
}

/**
 * Type guard to check if a value is an AuditEntry
 */
export function isAuditEntry(value: unknown): value is AuditEntry {
  if (typeof value !== 'object' || value === null) {
    return false;
  }

  const entry = value as Partial<AuditEntry>;

  return (
    typeof entry.id === 'string' &&
    typeof entry.contentId === 'string' &&
    Object.values(AuditAction).includes(entry.action as AuditAction) &&
    typeof entry.performedBy === 'string' &&
    typeof entry.performedAt === 'string' &&
    typeof entry.version === 'number' &&
    typeof entry.description === 'string' &&
    Object.values(ChangeSeverity).includes(entry.severity as ChangeSeverity)
  );
}

/**
 * Factory function to create initial ContentMetadata
 */
export function createContentMetadata(
  createdBy: string,
  options: {
    tags?: string[];
    notes?: string;
    governance?: ContentMetadata['governance'];
    workflow?: ContentMetadata['workflow'];
    localization?: ContentMetadata['localization'];
    custom?: Record<string, any>;
  } = {}
): ContentMetadata {
  const now = new Date().toISOString();
  const initialVersion = 1;

  return {
    createdAt: now,
    createdBy,
    updatedAt: now,
    updatedBy: createdBy,
    version: initialVersion,
    tags: options.tags || [],
    notes: options.notes,
    annotations: [],
    auditTrail: [{
      id: generateAuditId(),
      contentId: '', // Will be set when content is created
      action: AuditAction.CREATE,
      performedBy: createdBy,
      performedAt: now,
      version: initialVersion,
      description: 'Content item created',
      severity: ChangeSeverity.MAJOR
    }],
    workflow: options.workflow,
    governance: options.governance,
    localization: options.localization,
    custom: options.custom
  };
}

/**
 * Factory function to create a new annotation
 */
export function createAnnotation(
  type: ContentAnnotation['type'],
  content: string,
  createdBy: string,
  options: {
    tags?: string[];
    priority?: ContentAnnotation['priority'];
  } = {}
): ContentAnnotation {
  return {
    id: generateAnnotationId(),
    type,
    content,
    createdBy,
    createdAt: new Date().toISOString(),
    resolved: false,
    tags: options.tags,
    priority: options.priority || 'medium'
  };
}

/**
 * Factory function to create a new audit entry
 */
export function createAuditEntry(
  contentId: string,
  action: AuditAction,
  performedBy: string,
  version: number,
  description: string,
  options: {
    severity?: ChangeSeverity;
    previousVersion?: number;
    ipAddress?: string;
    userAgent?: string;
    context?: Record<string, any>;
    changes?: Record<string, { before: any; after: any }>;
  } = {}
): AuditEntry {
  return {
    id: generateAuditId(),
    contentId,
    action,
    performedBy,
    performedAt: new Date().toISOString(),
    version,
    previousVersion: options.previousVersion,
    description,
    severity: options.severity || ChangeSeverity.MINOR,
    ipAddress: options.ipAddress,
    userAgent: options.userAgent,
    context: options.context,
    changes: options.changes
  };
}

/**
 * Validates ContentMetadata against validation rules
 */
export function validateContentMetadata(
  metadata: ContentMetadata,
  rules: MetadataValidationRules = {}
): ValidationResult {
  const errors: ValidationResult['errors'] = [];
  const warnings: ValidationResult['warnings'] = [];

  // Validate tags
  if (rules.maxTags && metadata.tags.length > rules.maxTags) {
    errors.push({
      field: 'tags',
      message: `Too many tags. Maximum allowed: ${rules.maxTags}`,
      code: 'MAX_TAGS_EXCEEDED'
    });
  }

  if (rules.requiredTags) {
    const missingTags = rules.requiredTags.filter(tag => !metadata.tags.includes(tag));
    if (missingTags.length > 0) {
      errors.push({
        field: 'tags',
        message: `Missing required tags: ${missingTags.join(', ')}`,
        code: 'MISSING_REQUIRED_TAGS'
      });
    }
  }

  if (rules.allowedTagPatterns) {
    const invalidTags = metadata.tags.filter(tag =>
      !rules.allowedTagPatterns!.some(pattern => new RegExp(pattern).test(tag))
    );
    if (invalidTags.length > 0) {
      errors.push({
        field: 'tags',
        message: `Invalid tags: ${invalidTags.join(', ')}`,
        code: 'INVALID_TAG_PATTERN'
      });
    }
  }

  // Validate annotations
  if (rules.maxAnnotationLength) {
    const longAnnotations = metadata.annotations.filter(
      annotation => annotation.content.length > rules.maxAnnotationLength!
    );
    if (longAnnotations.length > 0) {
      warnings.push({
        field: 'annotations',
        message: `Some annotations exceed maximum length of ${rules.maxAnnotationLength} characters`,
        code: 'ANNOTATION_TOO_LONG'
      });
    }
  }

  // Run custom validators
  if (rules.customValidators) {
    for (const validator of rules.customValidators) {
      const result = validator(metadata);
      errors.push(...result.errors);
      warnings.push(...result.warnings);
    }
  }

  return {
    valid: errors.length === 0,
    errors,
    warnings
  };
}

/**
 * Normalizes ContentMetadata by cleaning and standardizing fields
 */
export function normalizeContentMetadata(metadata: ContentMetadata): ContentMetadata {
  return {
    ...metadata,
    tags: [...new Set(metadata.tags.map(tag => tag.toLowerCase().trim()))].sort(),
    annotations: metadata.annotations.map(annotation => ({
      ...annotation,
      content: annotation.content.trim(),
      tags: annotation.tags ? [...new Set(annotation.tags.map(tag => tag.toLowerCase().trim()))] : undefined
    })),
    auditTrail: metadata.auditTrail.sort((a, b) =>
      new Date(b.performedAt).getTime() - new Date(a.performedAt).getTime()
    )
  };
}

/**
 * Compares two ContentMetadata objects and returns differences
 */
export function compareContentMetadata(
  previous: ContentMetadata,
  current: ContentMetadata
): {
  hasChanges: boolean;
  changes: Record<string, { before: any; after: any }>;
  fieldChanges: string[];
} {
  const changes: Record<string, { before: any; after: any }> = {};
  const fieldChanges: string[] = [];

  // Compare basic fields
  const basicFields = ['createdAt', 'createdBy', 'updatedAt', 'updatedBy', 'version', 'notes'] as const;
  for (const field of basicFields) {
    if (previous[field] !== current[field]) {
      changes[field] = { before: previous[field], after: current[field] };
      fieldChanges.push(field);
    }
  }

  // Compare tags
  if (JSON.stringify(previous.tags.sort()) !== JSON.stringify(current.tags.sort())) {
    changes.tags = { before: previous.tags, after: current.tags };
    fieldChanges.push('tags');
  }

  // Compare annotations (by count and resolved status)
  if (previous.annotations.length !== current.annotations.length ||
      previous.annotations.some((ann, i) => ann.resolved !== current.annotations[i]?.resolved)) {
    changes.annotations = { before: previous.annotations, after: current.annotations };
    fieldChanges.push('annotations');
  }

  // Compare workflow
  if (JSON.stringify(previous.workflow) !== JSON.stringify(current.workflow)) {
    changes.workflow = { before: previous.workflow, after: current.workflow };
    fieldChanges.push('workflow');
  }

  // Compare governance
  if (JSON.stringify(previous.governance) !== JSON.stringify(current.governance)) {
    changes.governance = { before: previous.governance, after: current.governance };
    fieldChanges.push('governance');
  }

  return {
    hasChanges: fieldChanges.length > 0,
    changes,
    fieldChanges
  };
}

/**
 * Updates ContentMetadata with new changes and creates audit entry
 */
export function updateContentMetadata(
  metadata: ContentMetadata,
  updates: Partial<Pick<ContentMetadata, 'tags' | 'notes' | 'workflow' | 'governance' | 'quality' | 'analytics' | 'relationships' | 'localization' | 'custom'>>,
  updatedBy: string,
  auditInfo: {
    action: AuditAction;
    description: string;
    severity?: ChangeSeverity;
    context?: Record<string, any>;
    ipAddress?: string;
    userAgent?: string;
  }
): ContentMetadata {
  const previousMetadata = { ...metadata };
  const newVersion = metadata.version + 1;
  const now = new Date().toISOString();

  // Apply updates
  const updatedMetadata: ContentMetadata = {
    ...metadata,
    ...updates,
    updatedAt: now,
    updatedBy,
    version: newVersion
  };

  // Calculate changes
  const comparison = compareContentMetadata(previousMetadata, updatedMetadata);

  // Create audit entry
  const auditEntry = createAuditEntry(
    '', // contentId will be set by caller
    auditInfo.action,
    updatedBy,
    newVersion,
    auditInfo.description,
    {
      severity: auditInfo.severity,
      previousVersion: metadata.version,
      ipAddress: auditInfo.ipAddress,
      userAgent: auditInfo.userAgent,
      context: auditInfo.context,
      changes: comparison.changes
    }
  );

  // Add audit entry to trail
  updatedMetadata.auditTrail = [...metadata.auditTrail, auditEntry];

  return normalizeContentMetadata(updatedMetadata);
}

/**
 * Adds an annotation to ContentMetadata
 */
export function addAnnotation(
  metadata: ContentMetadata,
  annotation: ContentAnnotation,
  addedBy: string
): ContentMetadata {
  const auditEntry = createAuditEntry(
    '', // contentId will be set by caller
    AuditAction.METADATA_UPDATE,
    addedBy,
    metadata.version + 1,
    `Added ${annotation.type} annotation: ${annotation.content.substring(0, 50)}...`,
    {
      severity: ChangeSeverity.MINOR,
      context: { annotationId: annotation.id, annotationType: annotation.type }
    }
  );

  return {
    ...metadata,
    annotations: [...metadata.annotations, annotation],
    auditTrail: [...metadata.auditTrail, auditEntry],
    updatedAt: new Date().toISOString(),
    updatedBy: addedBy,
    version: metadata.version + 1
  };
}

/**
 * Resolves an annotation in ContentMetadata
 */
export function resolveAnnotation(
  metadata: ContentMetadata,
  annotationId: string,
  resolvedBy: string,
  resolution?: string
): ContentMetadata {
  const annotationIndex = metadata.annotations.findIndex(ann => ann.id === annotationId);
  if (annotationIndex === -1) {
    throw new Error(`Annotation with ID ${annotationId} not found`);
  }

  const updatedAnnotations = [...metadata.annotations];
  updatedAnnotations[annotationIndex] = {
    ...updatedAnnotations[annotationIndex],
    resolved: true,
    resolvedAt: new Date().toISOString(),
    resolvedBy
  };

  const auditEntry = createAuditEntry(
    '', // contentId will be set by caller
    AuditAction.METADATA_UPDATE,
    resolvedBy,
    metadata.version + 1,
    `Resolved annotation: ${resolution || 'No resolution provided'}`,
    {
      severity: ChangeSeverity.MINOR,
      context: { annotationId, resolution }
    }
  );

  return {
    ...metadata,
    annotations: updatedAnnotations,
    auditTrail: [...metadata.auditTrail, auditEntry],
    updatedAt: new Date().toISOString(),
    updatedBy: resolvedBy,
    version: metadata.version + 1
  };
}

/**
 * Adds tags to ContentMetadata
 */
export function addTags(
  metadata: ContentMetadata,
  newTags: string[],
  addedBy: string
): ContentMetadata {
  const normalizedNewTags = newTags.map(tag => tag.toLowerCase().trim());
  const uniqueNewTags = normalizedNewTags.filter(tag => !metadata.tags.includes(tag));

  if (uniqueNewTags.length === 0) {
    return metadata; // No new tags to add
  }

  const updatedTags = [...metadata.tags, ...uniqueNewTags].sort();

  const auditEntry = createAuditEntry(
    '', // contentId will be set by caller
    AuditAction.TAG_ADDED,
    addedBy,
    metadata.version + 1,
    `Added tags: ${uniqueNewTags.join(', ')}`,
    {
      severity: ChangeSeverity.MINOR,
      context: { addedTags: uniqueNewTags }
    }
  );

  return {
    ...metadata,
    tags: updatedTags,
    auditTrail: [...metadata.auditTrail, auditEntry],
    updatedAt: new Date().toISOString(),
    updatedBy: addedBy,
    version: metadata.version + 1
  };
}

/**
 * Removes tags from ContentMetadata
 */
export function removeTags(
  metadata: ContentMetadata,
  tagsToRemove: string[],
  removedBy: string
): ContentMetadata {
  const normalizedTagsToRemove = tagsToRemove.map(tag => tag.toLowerCase().trim());
  const updatedTags = metadata.tags.filter(tag => !normalizedTagsToRemove.includes(tag));
  const actuallyRemoved = metadata.tags.filter(tag => normalizedTagsToRemove.includes(tag));

  if (actuallyRemoved.length === 0) {
    return metadata; // No tags were actually removed
  }

  const auditEntry = createAuditEntry(
    '', // contentId will be set by caller
    AuditAction.TAG_REMOVED,
    removedBy,
    metadata.version + 1,
    `Removed tags: ${actuallyRemoved.join(', ')}`,
    {
      severity: ChangeSeverity.MINOR,
      context: { removedTags: actuallyRemoved }
    }
  );

  return {
    ...metadata,
    tags: updatedTags,
    auditTrail: [...metadata.auditTrail, auditEntry],
    updatedAt: new Date().toISOString(),
    updatedBy: removedBy,
    version: metadata.version + 1
  };
}

/**
 * Gets the most recent audit entries for a specific action type
 */
export function getRecentAuditEntries(
  metadata: ContentMetadata,
  action?: AuditAction,
  limit: number = 10
): AuditEntry[] {
  let entries = [...metadata.auditTrail];

  if (action) {
    entries = entries.filter(entry => entry.action === action);
  }

  return entries
    .sort((a, b) => new Date(b.performedAt).getTime() - new Date(a.performedAt).getTime())
    .slice(0, limit);
}

/**
 * Gets audit trail summary with statistics
 */
export function getAuditTrailSummary(metadata: ContentMetadata): {
  totalEntries: number;
  actionCounts: Record<AuditAction, number>;
  severityCounts: Record<ChangeSeverity, number>;
  contributors: string[];
  dateRange: { earliest: string; latest: string };
} {
  const actionCounts = Object.values(AuditAction).reduce((acc, action) => {
    acc[action] = 0;
    return acc;
  }, {} as Record<AuditAction, number>);

  const severityCounts = Object.values(ChangeSeverity).reduce((acc, severity) => {
    acc[severity] = 0;
    return acc;
  }, {} as Record<ChangeSeverity, number>);

  const contributors = new Set<string>();
  const dates: string[] = [];

  for (const entry of metadata.auditTrail) {
    actionCounts[entry.action]++;
    severityCounts[entry.severity]++;
    contributors.add(entry.performedBy);
    dates.push(entry.performedAt);
  }

  const sortedDates = dates.sort();

  return {
    totalEntries: metadata.auditTrail.length,
    actionCounts,
    severityCounts,
    contributors: Array.from(contributors),
    dateRange: {
      earliest: sortedDates[0] || '',
      latest: sortedDates[sortedDates.length - 1] || ''
    }
  };
}

/**
 * Generates a unique ID for annotations
 */
function generateAnnotationId(): string {
  return `ann_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`;
}

/**
 * Generates a unique ID for audit entries
 */
function generateAuditId(): string {
  return `audit_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`;
}

// ============================================================================
// CONTENT ITEM FACTORY FUNCTIONS
// ============================================================================

/**
 * Options for creating a ContentItem
 */
export interface CreateContentItemOptions {
  /** Content item ID */
  id: string;
  /** Semantic key for content item */
  key: string;
  /** Category the content belongs to */
  category: ContentCategory;
  /** Content type */
  type: ContentType;
  /** Content value */
  value: ContentValue;
  /** Current status (defaults to DRAFT) */
  status?: ContentStatus;
  /** Optional tags */
  tags?: string[];
  /** Optional editorial notes */
  notes?: string;
  /** Optional change log */
  changeLog?: string;
  /** User ID creating the content */
  createdBy: string;
  /** Additional metadata options */
  metadataOptions?: {
    governance?: ContentMetadata['governance'];
    workflow?: ContentMetadata['workflow'];
    localization?: ContentMetadata['localization'];
    custom?: Record<string, any>;
  };
}

/**
 * Factory function to create a new ContentItem with proper metadata
 * @param options - Configuration options for the content item
 * @returns A new ContentItem with generated metadata
 */
export function createContentItem(options: CreateContentItemOptions): ContentItem {
  if (!isValidContentValue(options.value)) {
    throw new Error(`Invalid content value for type ${options.type}`);
  }

  if (!isContentCategory(options.category)) {
    throw new Error('Invalid content category provided');
  }

  // Create initial metadata
  const metadata = createContentMetadata(options.createdBy, {
    tags: options.tags,
    notes: options.notes,
    governance: options.metadataOptions?.governance,
    workflow: options.metadataOptions?.workflow,
    localization: options.metadataOptions?.localization,
    custom: options.metadataOptions?.custom,
  });

  // Set the content ID in the initial audit entry
  metadata.auditTrail[0].contentId = options.id;

  return {
    id: options.id,
    key: options.key,
    categoryId: options.category.id,
    category: options.category,
    type: options.type,
    value: options.value,
    metadata,
    status: options.status || ContentStatus.DRAFT,
    tags: options.tags,
    notes: options.notes,
    changeLog: options.changeLog,
  };
}

/**
 * Factory function to create TextContent
 * @param value - The text content
 * @param maxLength - Optional maximum length constraint
 * @returns TextContent object
 */
export function createTextContent(value: string, maxLength?: number): TextContent {
  if (maxLength && value.length > maxLength) {
    throw new Error(`Text content exceeds maximum length of ${maxLength} characters`);
  }

  return {
    type: 'text',
    value,
    maxLength,
  };
}

/**
 * Factory function to create RichTextContent
 * @param value - The formatted content
 * @param format - Content format (html or markdown)
 * @param allowedTags - Optional allowed HTML tags for security
 * @returns RichTextContent object
 */
export function createRichTextContent(
  value: string,
  format: 'html' | 'markdown',
  allowedTags?: string[]
): RichTextContent {
  return {
    type: 'rich_text',
    value,
    format,
    allowedTags,
  };
}

/**
 * Factory function to create ImageContent
 * @param url - Image URL
 * @param alt - Alternative text for accessibility
 * @param metadata - Optional image metadata
 * @returns ImageContent object
 */
export function createImageContent(
  url: string,
  alt: string,
  metadata?: {
    width?: number;
    height?: number;
    format?: string;
  }
): ImageContent {
  return {
    type: 'image_url',
    url,
    alt,
    ...metadata,
  };
}

/**
 * Factory function to create ConfigContent
 * @param value - Configuration value
 * @param schema - Optional JSON schema for validation
 * @returns ConfigContent object
 */
export function createConfigContent(
  value: boolean | number | string | object,
  schema?: Record<string, any>
): ConfigContent {
  return {
    type: 'config',
    value,
    schema,
  };
}

/**
 * Factory function to create TranslationContent
 * @param values - Mapping of locale codes to translated text
 * @param defaultLocale - Default locale for fallback
 * @returns TranslationContent object
 */
export function createTranslationContent(
  values: Record<string, string>,
  defaultLocale: string
): TranslationContent {
  if (!(defaultLocale in values)) {
    throw new Error(`Default locale '${defaultLocale}' not found in translation values`);
  }

  return {
    type: 'translation',
    values,
    defaultLocale,
  };
}

// ============================================================================
// CONTENT ITEM UTILITY FUNCTIONS
// ============================================================================

/**
 * Deep clones a ContentItem with all its nested properties
 * @param item - The ContentItem to clone
 * @returns A deep copy of the ContentItem
 */
export function cloneContentItem(item: ContentItem): ContentItem {
  return {
    ...item,
    category: { ...item.category, permissions: { ...item.category.permissions } },
    value: cloneContentValue(item.value),
    metadata: cloneContentMetadata(item.metadata),
    tags: item.tags ? [...item.tags] : undefined,
  };
}

/**
 * Deep clones ContentMetadata
 * @param metadata - The ContentMetadata to clone
 * @returns A deep copy of the ContentMetadata
 */
function cloneContentMetadata(metadata: ContentMetadata): ContentMetadata {
  return {
    ...metadata,
    tags: [...metadata.tags],
    annotations: metadata.annotations.map(annotation => ({
      ...annotation,
      tags: annotation.tags ? [...annotation.tags] : undefined,
    })),
    auditTrail: metadata.auditTrail.map(entry => ({
      ...entry,
      context: entry.context ? { ...entry.context } : undefined,
      changes: entry.changes ? { ...entry.changes } : undefined,
    })),
    workflow: metadata.workflow ? { ...metadata.workflow } : undefined,
    governance: metadata.governance ? {
      ...metadata.governance,
      retentionPolicy: metadata.governance.retentionPolicy ? { ...metadata.governance.retentionPolicy } : undefined,
      complianceTags: metadata.governance.complianceTags ? [...metadata.governance.complianceTags] : undefined,
    } : undefined,
    quality: metadata.quality ? {
      ...metadata.quality,
      issues: metadata.quality.issues ? [...metadata.quality.issues] : undefined,
    } : undefined,
    analytics: metadata.analytics ? {
      ...metadata.analytics,
      performance: metadata.analytics.performance ? { ...metadata.analytics.performance } : undefined,
    } : undefined,
    relationships: metadata.relationships ? {
      ...metadata.relationships,
      childIds: metadata.relationships.childIds ? [...metadata.relationships.childIds] : undefined,
      relatedIds: metadata.relationships.relatedIds ? [...metadata.relationships.relatedIds] : undefined,
      dependencies: metadata.relationships.dependencies ? [...metadata.relationships.dependencies] : undefined,
    } : undefined,
    localization: metadata.localization ? {
      ...metadata.localization,
      translations: { ...metadata.localization.translations },
      tmReferences: metadata.localization.tmReferences ? [...metadata.localization.tmReferences] : undefined,
    } : undefined,
    custom: metadata.custom ? JSON.parse(JSON.stringify(metadata.custom)) : undefined,
  };
}

/**
 * Deep clones a ContentValue based on its type
 * @param value - The ContentValue to clone
 * @returns A deep copy of the ContentValue
 */
export function cloneContentValue(value: ContentValue): ContentValue {
  switch (value.type) {
    case 'text':
      return { ...value };
    case 'rich_text':
      return {
        ...value,
        allowedTags: value.allowedTags ? [...value.allowedTags] : undefined,
      };
    case 'image_url':
      return { ...value };
    case 'config':
      return {
        ...value,
        value: typeof value.value === 'object' && value.value !== null
          ? JSON.parse(JSON.stringify(value.value))
          : value.value,
        schema: value.schema ? JSON.parse(JSON.stringify(value.schema)) : undefined,
      };
    case 'translation':
      return {
        ...value,
        values: { ...value.values },
      };
    default:
      return value;
  }
}

/**
 * Serializes a ContentItem to JSON string with proper error handling
 * @param item - The ContentItem to serialize
 * @returns JSON string representation
 */
export function serializeContentItem(item: ContentItem): string {
  try {
    return JSON.stringify(item, null, 2);
  } catch (error) {
    throw new Error(`Failed to serialize ContentItem: ${error instanceof Error ? error.message : 'Unknown error'}`);
  }
}

/**
 * Deserializes a JSON string to ContentItem with validation
 * @param json - JSON string to deserialize
 * @returns Parsed and validated ContentItem
 */
export function deserializeContentItem(json: string): ContentItem {
  try {
    const parsed = JSON.parse(json);

    if (!isContentItem(parsed)) {
      throw new Error('Invalid ContentItem structure in JSON');
    }

    return parsed;
  } catch (error) {
    if (error instanceof SyntaxError) {
      throw new Error(`Invalid JSON format: ${error.message}`);
    }
    throw error;
  }
}

/**
 * Compares two ContentItems for equality (deep comparison)
 * @param item1 - First ContentItem
 * @param item2 - Second ContentItem
 * @returns true if items are equal, false otherwise
 */
export function compareContentItems(item1: ContentItem, item2: ContentItem): boolean {
  try {
    return JSON.stringify(normalizeContentItemForComparison(item1)) ===
           JSON.stringify(normalizeContentItemForComparison(item2));
  } catch {
    return false;
  }
}

/**
 * Normalizes a ContentItem for comparison by sorting arrays and objects
 * @param item - ContentItem to normalize
 * @returns Normalized ContentItem
 */
function normalizeContentItemForComparison(item: ContentItem): ContentItem {
  const normalized = cloneContentItem(item);

  // Sort tags arrays for consistent comparison
  if (normalized.tags) {
    normalized.tags.sort();
  }
  normalized.metadata.tags.sort();

  // Sort annotations by creation date
  normalized.metadata.annotations.sort((a, b) => a.createdAt.localeCompare(b.createdAt));

  // Sort audit trail by performed date
  normalized.metadata.auditTrail.sort((a, b) => a.performedAt.localeCompare(b.performedAt));

  return normalized;
}

/**
 * Extracts the display value from a ContentValue based on locale
 * @param value - The ContentValue to extract from
 * @param locale - Preferred locale (for translations)
 * @returns String representation of the content
 */
export function extractDisplayValue(value: ContentValue, locale = 'en'): string {
  switch (value.type) {
    case 'text':
    case 'rich_text':
      return value.value;
    case 'image_url':
      return value.alt || value.url;
    case 'config':
      return String(value.value);
    case 'translation':
      return value.values[locale] || value.values[value.defaultLocale] || '';
    default:
      return '';
  }
}

/**
 * Updates a ContentItem with new value and metadata
 * @param item - Original ContentItem
 * @param newValue - New content value
 * @param updatedBy - User ID performing the update
 * @param options - Optional update options
 * @returns Updated ContentItem
 */
export function updateContentItem(
  item: ContentItem,
  newValue: ContentValue,
  updatedBy: string,
  options?: {
    status?: ContentStatus;
    tags?: string[];
    notes?: string;
    changeLog?: string;
    annotation?: ContentAnnotation;
    auditContext?: Record<string, any>;
    ipAddress?: string;
    userAgent?: string;
  }
): ContentItem {
  if (!isValidContentValue(newValue)) {
    throw new Error(`Invalid content value for type ${newValue.type}`);
  }

  // Update metadata with audit trail
  const updatedMetadata = updateContentMetadata(
    item.metadata,
    {
      tags: options?.tags || item.metadata.tags,
      notes: options?.notes || item.metadata.notes,
    },
    updatedBy,
    {
      action: AuditAction.UPDATE,
      description: options?.changeLog || 'Content updated',
      severity: ChangeSeverity.MINOR,
      context: options?.auditContext,
      ipAddress: options?.ipAddress,
      userAgent: options?.userAgent,
    }
  );

  // Set content ID in the new audit entry
  updatedMetadata.auditTrail[updatedMetadata.auditTrail.length - 1].contentId = item.id;

  // Add annotation if provided
  let finalMetadata = updatedMetadata;
  if (options?.annotation) {
    finalMetadata = addAnnotation(finalMetadata, options.annotation, updatedBy);
    finalMetadata.auditTrail[finalMetadata.auditTrail.length - 1].contentId = item.id;
  }

  return {
    ...item,
    value: newValue,
    status: options?.status || item.status,
    tags: options?.tags || item.tags,
    notes: options?.notes || item.notes,
    changeLog: options?.changeLog || item.changeLog,
    metadata: finalMetadata,
  };
}

// ============================================================================
// CONTENT ITEM VALIDATION FUNCTIONS
// ============================================================================

/**
 * Validates a ContentItem with detailed error reporting
 * @param item - ContentItem to validate
 * @param rules - Optional validation rules
 * @returns Detailed validation result
 */
export function validateContentItem(
  item: unknown,
  rules?: MetadataValidationRules
): ValidationResult {
  const result: ValidationResult = {
    valid: true,
    errors: [],
    warnings: [],
  };

  if (!item || typeof item !== 'object') {
    result.valid = false;
    result.errors.push({
      field: 'root',
      message: 'ContentItem must be a non-null object',
      code: 'INVALID_TYPE',
    });
    return result;
  }

  const contentItem = item as Partial<ContentItem>;

  // Validate required fields
  if (!contentItem.id || typeof contentItem.id !== 'string') {
    result.valid = false;
    result.errors.push({
      field: 'id',
      message: 'ContentItem.id is required and must be a string',
      code: 'REQUIRED_FIELD',
    });
  }

  if (!contentItem.key || typeof contentItem.key !== 'string') {
    result.valid = false;
    result.errors.push({
      field: 'key',
      message: 'ContentItem.key is required and must be a string',
      code: 'REQUIRED_FIELD',
    });
  }

  if (!contentItem.categoryId || typeof contentItem.categoryId !== 'string') {
    result.valid = false;
    result.errors.push({
      field: 'categoryId',
      message: 'ContentItem.categoryId is required and must be a string',
      code: 'REQUIRED_FIELD',
    });
  }

  if (!isContentCategory(contentItem.category)) {
    result.valid = false;
    result.errors.push({
      field: 'category',
      message: 'ContentItem.category is invalid',
      code: 'INVALID_STRUCTURE',
    });
  }

  if (!Object.values(ContentType).includes(contentItem.type as ContentType)) {
    result.valid = false;
    result.errors.push({
      field: 'type',
      message: 'ContentItem.type must be a valid ContentType',
      code: 'INVALID_ENUM',
    });
  }

  if (!isValidContentValue(contentItem.value)) {
    result.valid = false;
    result.errors.push({
      field: 'value',
      message: 'ContentItem.value is invalid for the specified type',
      code: 'INVALID_STRUCTURE',
    });
  }

  if (!isValidContentMetadata(contentItem.metadata)) {
    result.valid = false;
    result.errors.push({
      field: 'metadata',
      message: 'ContentItem.metadata is invalid',
      code: 'INVALID_STRUCTURE',
    });
  }

  if (!Object.values(ContentStatus).includes(contentItem.status as ContentStatus)) {
    result.valid = false;
    result.errors.push({
      field: 'status',
      message: 'ContentItem.status must be a valid ContentStatus',
      code: 'INVALID_ENUM',
    });
  }

  // Validate optional fields
  if (contentItem.tags && !Array.isArray(contentItem.tags)) {
    result.valid = false;
    result.errors.push({
      field: 'tags',
      message: 'ContentItem.tags must be an array if provided',
      code: 'INVALID_TYPE',
    });
  }

  if (contentItem.notes && typeof contentItem.notes !== 'string') {
    result.valid = false;
    result.errors.push({
      field: 'notes',
      message: 'ContentItem.notes must be a string if provided',
      code: 'INVALID_TYPE',
    });
  }

  // Validate content ID pattern
  if (contentItem.id && !isValidContentIdPattern(contentItem.id)) {
    result.warnings.push({
      field: 'id',
      message: 'ContentItem.id should follow the pattern: category.subcategory.key',
      code: 'PATTERN_RECOMMENDATION',
    });
  }

  // Validate content value constraints
  if (contentItem.value) {
    const constraintValidation = validateContentValueConstraints(contentItem.value);
    if (!constraintValidation.valid) {
      result.valid = false;
      result.errors.push(...constraintValidation.errors);
    }
    result.warnings.push(...constraintValidation.warnings);
  }

  // Validate metadata with rules
  if (contentItem.metadata && rules) {
    const metadataValidation = validateContentMetadata(contentItem.metadata, rules);
    if (!metadataValidation.valid) {
      result.valid = false;
      result.errors.push(...metadataValidation.errors);
    }
    result.warnings.push(...metadataValidation.warnings);
  }

  return result;
}

/**
 * Validates if a string follows the content ID pattern
 * @param id - Content ID to validate
 * @returns true if ID follows the pattern
 */
export function isValidContentIdPattern(id: string): boolean {
  const pattern = /^[a-zA-Z0-9_]+\.[a-zA-Z0-9_]+\.[a-zA-Z0-9_]+$/;
  return pattern.test(id);
}

/**
 * Validates content value constraints (e.g., text length, required fields)
 * @param value - ContentValue to validate
 * @returns Validation result with specific constraint violations
 */
export function validateContentValueConstraints(value: ContentValue): ValidationResult {
  const result: ValidationResult = {
    valid: true,
    errors: [],
    warnings: [],
  };

  switch (value.type) {
    case 'text':
      const textContent = value as TextContent;
      if (textContent.maxLength && textContent.value.length > textContent.maxLength) {
        result.valid = false;
        result.errors.push({
          field: 'value.value',
          message: `Text content exceeds maximum length of ${textContent.maxLength} characters`,
          code: 'LENGTH_EXCEEDED',
        });
      }
      if (textContent.value.length === 0) {
        result.warnings.push({
          field: 'value.value',
          message: 'Text content is empty',
          code: 'EMPTY_CONTENT',
        });
      }
      break;

    case 'rich_text':
      const richTextContent = value as RichTextContent;
      if (richTextContent.value.length === 0) {
        result.warnings.push({
          field: 'value.value',
          message: 'Rich text content is empty',
          code: 'EMPTY_CONTENT',
        });
      }
      break;

    case 'image_url':
      const imageContent = value as ImageContent;
      try {
        new URL(imageContent.url);
      } catch {
        result.errors.push({
          field: 'value.url',
          message: 'Image URL is not a valid URL',
          code: 'INVALID_URL',
        });
        result.valid = false;
      }
      if (!imageContent.alt || imageContent.alt.length === 0) {
        result.warnings.push({
          field: 'value.alt',
          message: 'Image alt text is empty - this affects accessibility',
          code: 'ACCESSIBILITY_CONCERN',
        });
      }
      break;

    case 'translation':
      const translationContent = value as TranslationContent;
      if (Object.keys(translationContent.values).length === 0) {
        result.valid = false;
        result.errors.push({
          field: 'value.values',
          message: 'Translation content must have at least one locale',
          code: 'EMPTY_TRANSLATIONS',
        });
      }
      Object.entries(translationContent.values).forEach(([locale, text]) => {
        if (!text || text.length === 0) {
          result.warnings.push({
            field: `value.values.${locale}`,
            message: `Translation for locale '${locale}' is empty`,
            code: 'EMPTY_TRANSLATION',
          });
        }
      });
      break;
  }

  return result;
}

// ============================================================================
// CONTENT ITEM ANALYSIS UTILITIES
// ============================================================================

/**
 * Gets content quality score based on validation and metadata
 * @param item - ContentItem to evaluate
 * @returns Quality score from 0 to 100
 */
export function calculateContentQualityScore(item: ContentItem): number {
  let score = 100;

  // Validate the item
  const validation = validateContentItem(item);

  // Deduct points for errors and warnings
  score -= validation.errors.length * 20;
  score -= validation.warnings.length * 5;

  // Deduct points for missing metadata
  if (!item.tags || item.tags.length === 0) score -= 10;
  if (!item.notes) score -= 5;
  if (item.metadata.annotations.length === 0) score -= 5;

  // Add points for comprehensive metadata
  if (item.metadata.governance) score += 5;
  if (item.metadata.workflow) score += 5;
  if (item.metadata.localization) score += 5;

  // Ensure score is within bounds
  return Math.max(0, Math.min(100, score));
}

/**
 * Finds related content items based on tags and category
 * @param item - ContentItem to find relations for
 * @param allItems - Array of all available ContentItems
 * @param maxResults - Maximum number of related items to return
 * @returns Array of related ContentItems
 */
export function findRelatedContentItems(
  item: ContentItem,
  allItems: ContentItem[],
  maxResults = 10
): ContentItem[] {
  const scores = allItems
    .filter(other => other.id !== item.id) // Exclude the item itself
    .map(other => ({
      item: other,
      score: calculateRelationScore(item, other),
    }))
    .filter(({ score }) => score > 0) // Only include items with some relation
    .sort((a, b) => b.score - a.score) // Sort by highest score first
    .slice(0, maxResults); // Limit results

  return scores.map(({ item }) => item);
}

/**
 * Calculates relation score between two content items
 * @param item1 - First ContentItem
 * @param item2 - Second ContentItem
 * @returns Relation score (higher = more related)
 */
function calculateRelationScore(item1: ContentItem, item2: ContentItem): number {
  let score = 0;

  // Same category = high relevance
  if (item1.categoryId === item2.categoryId) {
    score += 30;
  }

  // Same type = medium relevance
  if (item1.type === item2.type) {
    score += 20;
  }

  // Shared tags = relevance per shared tag
  const sharedTags = item1.tags?.filter(tag => item2.tags?.includes(tag)) || [];
  score += sharedTags.length * 10;

  // Similar keys (substring match) = low relevance
  if (item1.key.includes(item2.key) || item2.key.includes(item1.key)) {
    score += 5;
  }

  return score;
}

/**
 * Gets content statistics for analytics
 * @param items - Array of ContentItems to analyze
 * @returns Content statistics object
 */
export function getContentStatistics(items: ContentItem[]): {
  totalItems: number;
  itemsByType: Record<ContentType, number>;
  itemsByStatus: Record<ContentStatus, number>;
  averageQualityScore: number;
  totalAnnotations: number;
  mostUsedTags: Array<{ tag: string; count: number }>;
  recentActivity: Array<{ date: string; count: number }>;
} {
  const itemsByType = Object.values(ContentType).reduce((acc, type) => {
    acc[type] = 0;
    return acc;
  }, {} as Record<ContentType, number>);

  const itemsByStatus = Object.values(ContentStatus).reduce((acc, status) => {
    acc[status] = 0;
    return acc;
  }, {} as Record<ContentStatus, number>);

  const tagCounts: Record<string, number> = {};
  const activityByDate: Record<string, number> = {};
  let totalQualityScore = 0;
  let totalAnnotations = 0;

  for (const item of items) {
    // Count by type and status
    itemsByType[item.type]++;
    itemsByStatus[item.status]++;

    // Calculate quality scores
    totalQualityScore += calculateContentQualityScore(item);

    // Count annotations
    totalAnnotations += item.metadata.annotations.length;

    // Count tags
    if (item.tags) {
      for (const tag of item.tags) {
        tagCounts[tag] = (tagCounts[tag] || 0) + 1;
      }
    }

    // Count activity by date
    const date = item.metadata.updatedAt.split('T')[0];
    activityByDate[date] = (activityByDate[date] || 0) + 1;
  }

  // Get most used tags
  const mostUsedTags = Object.entries(tagCounts)
    .map(([tag, count]) => ({ tag, count }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 10);

  // Get recent activity (last 30 days)
  const recentActivity = Object.entries(activityByDate)
    .map(([date, count]) => ({ date, count }))
    .sort((a, b) => b.date.localeCompare(a.date))
    .slice(0, 30);

  return {
    totalItems: items.length,
    itemsByType,
    itemsByStatus,
    averageQualityScore: items.length > 0 ? totalQualityScore / items.length : 0,
    totalAnnotations,
    mostUsedTags,
    recentActivity,
  };
}// ============================================================================
// ENHANCED CONTENTVALUE TYPES AND UTILITIES (T022)
// ============================================================================

/**
 * Specific value interfaces for type-safe content value operations
 */

/**
 * Text value interface with validation constraints
 */
export interface TextValue {
  /** The text content */
  value: string;
  /** Maximum allowed character length */
  maxLength?: number;
  /** Minimum required character length */
  minLength?: number;
  /** Regex pattern for validation */
  pattern?: string;
  /** Whether the text is required (non-empty) */
  required?: boolean;
}

/**
 * Rich text value interface with format-specific constraints
 */
export interface RichTextValue {
  /** The formatted content (HTML or Markdown) */
  value: string;
  /** Content format specification */
  format: 'html' | 'markdown';
  /** Allowed HTML tags for security (whitelist) */
  allowedTags?: string[];
  /** Forbidden HTML tags (blacklist) */
  forbiddenTags?: string[];
  /** Maximum allowed character length */
  maxLength?: number;
  /** Whether to sanitize content automatically */
  autoSanitize?: boolean;
}

/**
 * Image value interface with validation and metadata
 */
export interface ImageValue {
  /** Image URL */
  url: string;
  /** Alternative text for accessibility */
  alt: string;
  /** Image width in pixels */
  width?: number;
  /** Image height in pixels */
  height?: number;
  /** Image format (e.g., 'png', 'jpg', 'webp') */
  format?: string;
  /** File size in bytes */
  fileSize?: number;
  /** Maximum allowed file size in bytes */
  maxFileSize?: number;
  /** Allowed image formats */
  allowedFormats?: string[];
  /** Whether the image is required */
  required?: boolean;
}

/**
 * Config value interface with type-safe configuration options
 */
export interface ConfigValue {
  /** Configuration value (boolean, number, string, or object) */
  value: boolean | number | string | object;
  /** JSON schema for validation (optional) */
  schema?: Record<string, any>;
  /** Default value if not specified */
  defaultValue?: boolean | number | string | object;
  /** Whether the config is required */
  required?: boolean;
  /** Environment-specific overrides */
  envOverrides?: Record<string, boolean | number | string | object>;
}

/**
 * Translation value interface with locale management
 */
export interface TranslationValue {
  /** Mapping of locale codes to translated text */
  values: Record<string, string>;
  /** Default locale for fallback */
  defaultLocale: string;
  /** Required locales that must have values */
  requiredLocales?: string[];
  /** Supported locales */
  supportedLocales?: string[];
  /** Whether translations are required for all supported locales */
  requireAllLocales?: boolean;
  /** Interpolation variables available in translations */
  variables?: Record<string, string>;
}

// ============================================================================
// ENHANCED VALIDATION RESULT INTERFACE (T022)
// ============================================================================

/**
 * Enhanced validation result interface with consistent structure
 */
export interface ContentValidationResult {
  /** Whether the validation passed */
  isValid: boolean;
  /** Array of validation error messages */
  errors: string[];
  /** Array of validation warnings */
  warnings?: string[];
}

// ============================================================================
// CONTENT VALUE VALIDATION FUNCTIONS
// ============================================================================

/**
 * Validates a text content value
 */
export function validateTextValue(value: TextValue): ContentValidationResult {
  const errors: string[] = [];
  const warnings: string[] = [];

  // Required validation
  if (value.required && (!value.value || value.value.trim() === '')) {
    errors.push('Text value is required');
  }

  // Length validation
  if (value.minLength && value.value.length < value.minLength) {
    errors.push(`Text must be at least ${value.minLength} characters long`);
  }

  if (value.maxLength && value.value.length > value.maxLength) {
    errors.push(`Text must not exceed ${value.maxLength} characters`);
  }

  // Pattern validation
  if (value.pattern && value.value) {
    try {
      const regex = new RegExp(value.pattern);
      if (!regex.test(value.value)) {
        errors.push('Text does not match required pattern');
      }
    } catch (error) {
      errors.push('Invalid pattern specified');
    }
  }

  return {
    isValid: errors.length === 0,
    errors,
    warnings
  };
}

/**
 * Validates a rich text content value
 */
export function validateRichTextValue(value: RichTextValue): ContentValidationResult {
  const errors: string[] = [];
  const warnings: string[] = [];

  // Basic validation
  if (!value.value) {
    errors.push('Rich text value is required');
  }

  // Format validation
  if (!['html', 'markdown'].includes(value.format)) {
    errors.push('Rich text format must be "html" or "markdown"');
  }

  // Length validation
  if (value.maxLength && value.value.length > value.maxLength) {
    errors.push(`Rich text must not exceed ${value.maxLength} characters`);
  }

  // HTML tag validation for HTML format
  if (value.format === 'html' && value.value) {
    if (value.forbiddenTags?.length) {
      const forbiddenRegex = new RegExp(`<(${value.forbiddenTags.join('|')})`, 'gi');
      if (forbiddenRegex.test(value.value)) {
        errors.push('Rich text contains forbidden HTML tags');
      }
    }

    if (value.allowedTags?.length) {
      // Extract all HTML tags from content
      const tagMatches = value.value.match(/<\/?([a-zA-Z][a-zA-Z0-9]*)/g);
      if (tagMatches) {
        const usedTags = tagMatches.map(tag => tag.replace(/<\/?/, '').toLowerCase());
        const disallowedTags = usedTags.filter(tag => !value.allowedTags!.includes(tag));
        if (disallowedTags.length > 0) {
          errors.push(`Rich text contains disallowed HTML tags: ${disallowedTags.join(', ')}`);
        }
      }
    }
  }

  return {
    isValid: errors.length === 0,
    errors,
    warnings
  };
}

/**
 * Validates an image content value
 */
export function validateImageValue(value: ImageValue): ContentValidationResult {
  const errors: string[] = [];
  const warnings: string[] = [];

  // Required validation
  if (value.required && !value.url) {
    errors.push('Image URL is required');
  }

  // URL validation
  if (value.url) {
    try {
      new URL(value.url);
    } catch {
      errors.push('Invalid image URL format');
    }
  }

  // Alt text validation
  if (value.url && !value.alt) {
    warnings.push('Alt text is recommended for accessibility');
  }

  // Dimension validation
  if (value.width && value.width <= 0) {
    errors.push('Image width must be greater than 0');
  }

  if (value.height && value.height <= 0) {
    errors.push('Image height must be greater than 0');
  }

  // File size validation
  if (value.maxFileSize && value.fileSize && value.fileSize > value.maxFileSize) {
    errors.push(`Image file size (${value.fileSize} bytes) exceeds maximum (${value.maxFileSize} bytes)`);
  }

  // Format validation
  if (value.allowedFormats?.length && value.format && !value.allowedFormats.includes(value.format)) {
    errors.push(`Image format "${value.format}" is not allowed. Allowed formats: ${value.allowedFormats.join(', ')}`);
  }

  return {
    isValid: errors.length === 0,
    errors,
    warnings
  };
}

/**
 * Validates a config content value
 */
export function validateConfigValue(value: ConfigValue): ContentValidationResult {
  const errors: string[] = [];
  const warnings: string[] = [];

  // Required validation
  if (value.required && (value.value === undefined || value.value === null)) {
    errors.push('Config value is required');
  }

  // Schema validation (basic JSON Schema support)
  if (value.schema && value.value !== undefined && value.value !== null) {
    try {
      // Basic type validation
      if (value.schema.type) {
        const actualType = typeof value.value;
        const expectedType = value.schema.type;

        if (expectedType === 'array' && !Array.isArray(value.value)) {
          errors.push(`Config value must be an array, got ${actualType}`);
        } else if (expectedType !== 'array' && actualType !== expectedType) {
          errors.push(`Config value must be of type ${expectedType}, got ${actualType}`);
        }
      }

      // Enum validation
      if (value.schema.enum && !value.schema.enum.includes(value.value)) {
        errors.push(`Config value must be one of: ${value.schema.enum.join(', ')}`);
      }

      // Number range validation
      if (typeof value.value === 'number') {
        if (value.schema.minimum !== undefined && value.value < value.schema.minimum) {
          errors.push(`Config value must be at least ${value.schema.minimum}`);
        }
        if (value.schema.maximum !== undefined && value.value > value.schema.maximum) {
          errors.push(`Config value must not exceed ${value.schema.maximum}`);
        }
      }

      // String length validation
      if (typeof value.value === 'string') {
        if (value.schema.minLength !== undefined && value.value.length < value.schema.minLength) {
          errors.push(`Config value must be at least ${value.schema.minLength} characters long`);
        }
        if (value.schema.maxLength !== undefined && value.value.length > value.schema.maxLength) {
          errors.push(`Config value must not exceed ${value.schema.maxLength} characters`);
        }
      }
    } catch (error) {
      errors.push('Invalid schema or value format');
    }
  }

  return {
    isValid: errors.length === 0,
    errors,
    warnings
  };
}

/**
 * Validates a translation content value
 */
export function validateTranslationValue(value: TranslationValue): ContentValidationResult {
  const errors: string[] = [];
  const warnings: string[] = [];

  // Default locale validation
  if (!value.defaultLocale) {
    errors.push('Default locale is required');
  }

  if (value.defaultLocale && !value.values[value.defaultLocale]) {
    errors.push(`Default locale "${value.defaultLocale}" must have a translation`);
  }

  // Required locales validation
  if (value.requiredLocales?.length) {
    const missingLocales = value.requiredLocales.filter(locale => !value.values[locale]);
    if (missingLocales.length > 0) {
      errors.push(`Missing required translations for locales: ${missingLocales.join(', ')}`);
    }
  }

  // Supported locales validation
  if (value.supportedLocales?.length) {
    const unsupportedLocales = Object.keys(value.values).filter(
      locale => !value.supportedLocales!.includes(locale)
    );
    if (unsupportedLocales.length > 0) {
      warnings.push(`Unsupported locales found: ${unsupportedLocales.join(', ')}`);
    }

    // Check if all supported locales are required
    if (value.requireAllLocales) {
      const missingLocales = value.supportedLocales.filter(locale => !value.values[locale]);
      if (missingLocales.length > 0) {
        errors.push(`Missing translations for supported locales: ${missingLocales.join(', ')}`);
      }
    }
  }

  // Empty translation validation
  const emptyTranslations = Object.entries(value.values)
    .filter(([_, translation]) => !translation || translation.trim() === '')
    .map(([locale]) => locale);

  if (emptyTranslations.length > 0) {
    warnings.push(`Empty translations found for locales: ${emptyTranslations.join(', ')}`);
  }

  // Variable interpolation validation
  if (value.variables) {
    const variablePattern = /\{\{(\w+)\}\}/g;
    Object.entries(value.values).forEach(([locale, translation]) => {
      const matches = [...translation.matchAll(variablePattern)];
      const usedVariables = matches.map(match => match[1]);
      const undefinedVariables = usedVariables.filter(variable => !value.variables![variable]);

      if (undefinedVariables.length > 0) {
        warnings.push(`Undefined variables in ${locale}: ${undefinedVariables.join(', ')}`);
      }
    });
  }

  return {
    isValid: errors.length === 0,
    errors,
    warnings
  };
}

// ============================================================================
// ENHANCED TYPE GUARDS FOR RUNTIME TYPE CHECKING
// ============================================================================

/**
 * Type guard to check if a value is a TextValue
 */
export function isTextValue(value: unknown): value is TextValue {
  return (
    typeof value === 'object' &&
    value !== null &&
    'value' in value &&
    typeof (value as TextValue).value === 'string'
  );
}

/**
 * Type guard to check if a value is a RichTextValue
 */
export function isRichTextValue(value: unknown): value is RichTextValue {
  return (
    typeof value === 'object' &&
    value !== null &&
    'value' in value &&
    'format' in value &&
    typeof (value as RichTextValue).value === 'string' &&
    ['html', 'markdown'].includes((value as RichTextValue).format)
  );
}

/**
 * Type guard to check if a value is an ImageValue
 */
export function isImageValue(value: unknown): value is ImageValue {
  return (
    typeof value === 'object' &&
    value !== null &&
    'url' in value &&
    'alt' in value &&
    typeof (value as ImageValue).url === 'string' &&
    typeof (value as ImageValue).alt === 'string'
  );
}

/**
 * Type guard to check if a value is a ConfigValue
 */
export function isConfigValue(value: unknown): value is ConfigValue {
  return (
    typeof value === 'object' &&
    value !== null &&
    'value' in value &&
    (value as ConfigValue).value !== undefined
  );
}

/**
 * Type guard to check if a value is a TranslationValue
 */
export function isTranslationValue(value: unknown): value is TranslationValue {
  return (
    typeof value === 'object' &&
    value !== null &&
    'values' in value &&
    'defaultLocale' in value &&
    typeof (value as TranslationValue).values === 'object' &&
    typeof (value as TranslationValue).defaultLocale === 'string'
  );
}

// ============================================================================
// SERIALIZATION/DESERIALIZATION HELPERS
// ============================================================================

/**
 * Serializes a ContentValue to a JSON-safe string
 */
export function serializeContentValue(value: ContentValue): string {
  try {
    return JSON.stringify(value);
  } catch (error) {
    throw new Error(`Failed to serialize content value: ${error instanceof Error ? error.message : 'Unknown error'}`);
  }
}

/**
 * Deserializes a JSON string to a ContentValue
 */
export function deserializeContentValue(serialized: string): ContentValue {
  try {
    const parsed = JSON.parse(serialized);

    // Validate that the deserialized object is a valid ContentValue
    if (!isValidContentValue(parsed)) {
      throw new Error('Deserialized object is not a valid ContentValue');
    }

    return parsed;
  } catch (error) {
    throw new Error(`Failed to deserialize content value: ${error instanceof Error ? error.message : 'Unknown error'}`);
  }
}

/**
 * Safely serializes a ContentValue with error handling
 */
export function safeSerializeContentValue(value: ContentValue): { success: boolean; data?: string; error?: string } {
  try {
    const serialized = serializeContentValue(value);
    return { success: true, data: serialized };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown serialization error'
    };
  }
}

/**
 * Safely deserializes a JSON string to ContentValue with error handling
 */
export function safeDeserializeContentValue(serialized: string): { success: boolean; data?: ContentValue; error?: string } {
  try {
    const deserialized = deserializeContentValue(serialized);
    return { success: true, data: deserialized };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown deserialization error'
    };
  }
}

// ============================================================================
// CONTENT VALUE TRANSFORMATION UTILITIES
// ============================================================================

/**
 * Transforms a TextContent to TextValue
 */
export function toTextValue(content: TextContent): TextValue {
  return {
    value: content.value,
    maxLength: content.maxLength,
    required: true
  };
}

/**
 * Transforms a RichTextContent to RichTextValue
 */
export function toRichTextValue(content: RichTextContent): RichTextValue {
  return {
    value: content.value,
    format: content.format,
    allowedTags: content.allowedTags,
    autoSanitize: true
  };
}

/**
 * Transforms an ImageContent to ImageValue
 */
export function toImageValue(content: ImageContent): ImageValue {
  return {
    url: content.url,
    alt: content.alt,
    width: content.width,
    height: content.height,
    format: content.format,
    required: true
  };
}

/**
 * Transforms a ConfigContent to ConfigValue
 */
export function toConfigValue(content: ConfigContent): ConfigValue {
  return {
    value: content.value,
    schema: content.schema,
    required: true
  };
}

/**
 * Transforms a TranslationContent to TranslationValue
 */
export function toTranslationValue(content: TranslationContent): TranslationValue {
  return {
    values: content.values,
    defaultLocale: content.defaultLocale,
    requireAllLocales: false
  };
}

/**
 * Transforms a ContentValue back to its specific content interface
 */
export function fromContentValue(value: ContentValue): TextContent | RichTextContent | ImageContent | ConfigContent | TranslationContent {
  switch (value.type) {
    case 'text':
      return {
        type: 'text',
        value: value.value,
        maxLength: value.maxLength
      };
    case 'rich_text':
      return {
        type: 'rich_text',
        value: value.value,
        format: value.format,
        allowedTags: value.allowedTags
      };
    case 'image_url':
      return {
        type: 'image_url',
        url: value.url,
        alt: value.alt,
        width: value.width,
        height: value.height,
        format: value.format
      };
    case 'config':
      return {
        type: 'config',
        value: value.value,
        schema: value.schema
      };
    case 'translation':
      return {
        type: 'translation',
        values: value.values,
        defaultLocale: value.defaultLocale
      };
  }
}


/**
 * Merges two ContentValues of the same type
 */
export function mergeContentValues(base: ContentValue, override: Partial<ContentValue>): ContentValue {
  if (base.type !== override.type && override.type !== undefined) {
    throw new Error('Cannot merge ContentValues of different types');
  }

  return {
    ...base,
    ...override,
    type: base.type // Ensure type is preserved
  } as ContentValue;
}

/**
 * Gets the default value for a specific content type
 */
export function getDefaultContentValue(type: ContentType): ContentValue {
  switch (type) {
    case ContentType.TEXT:
      return {
        type: 'text',
        value: '',
        maxLength: 1000
      };
    case ContentType.RICH_TEXT:
      return {
        type: 'rich_text',
        value: '',
        format: 'markdown',
        allowedTags: ['p', 'br', 'strong', 'em', 'ul', 'ol', 'li']
      };
    case ContentType.IMAGE_URL:
      return {
        type: 'image_url',
        url: '',
        alt: ''
      };
    case ContentType.CONFIG:
      return {
        type: 'config',
        value: ''
      };
    case ContentType.TRANSLATION:
      return {
        type: 'translation',
        values: {},
        defaultLocale: 'en'
      };
    default:
      throw new Error(`Unknown content type: ${type}`);
  }
}

/**
 * Validates any ContentValue using the appropriate validation function
 */
export function validateContentValue(value: ContentValue): ContentValidationResult {
  switch (value.type) {
    case 'text':
      return validateTextValue(toTextValue(value));
    case 'rich_text':
      return validateRichTextValue(toRichTextValue(value));
    case 'image_url':
      return validateImageValue(toImageValue(value));
    case 'config':
      return validateConfigValue(toConfigValue(value));
    case 'translation':
      return validateTranslationValue(toTranslationValue(value));
    default:
      return {
        isValid: false,
        errors: [`Unknown content type: ${(value as any).type}`]
      };
  }
}
