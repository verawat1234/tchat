# Data Model: Dynamic Content Management

## Core Entities

### ContentItem
**Purpose**: Represents any piece of information displayed to users
```typescript
interface ContentItem {
  id: string;                    // Unique identifier (e.g., "header.title", "error.network")
  category: ContentCategory;     // Semantic grouping
  type: ContentType;            // Data type and rendering info
  value: ContentValue;          // The actual content
  metadata: ContentMetadata;    // Management information
  status: ContentStatus;        // Lifecycle state
}

enum ContentType {
  TEXT = 'text',              // Plain text content
  RICH_TEXT = 'rich_text',    // HTML/Markdown content
  IMAGE_URL = 'image_url',    // Image URLs
  CONFIG = 'config',          // Configuration values
  TRANSLATION = 'translation' // Localized content
}

enum ContentStatus {
  DRAFT = 'draft',           // Work in progress
  PUBLISHED = 'published',   // Live content
  ARCHIVED = 'archived'      // Deprecated content
}
```

### ContentCategory
**Purpose**: Groups related content items for organization and management
```typescript
interface ContentCategory {
  id: string;                    // Unique identifier (e.g., "navigation", "errors", "help")
  name: string;                 // Human-readable name
  description: string;          // Purpose and usage notes
  parentId?: string;           // Hierarchical organization
  permissions: CategoryPermissions; // Access control
}

interface CategoryPermissions {
  read: string[];              // User roles that can view
  write: string[];             // User roles that can edit
  publish: string[];           // User roles that can publish
}
```

### ContentValue
**Purpose**: Flexible content storage supporting different data types
```typescript
type ContentValue =
  | TextContent
  | RichTextContent
  | ImageContent
  | ConfigContent
  | TranslationContent;

interface TextContent {
  type: 'text';
  value: string;
  maxLength?: number;
}

interface RichTextContent {
  type: 'rich_text';
  value: string;               // HTML or Markdown
  format: 'html' | 'markdown';
  allowedTags?: string[];     // Security whitelist
}

interface ImageContent {
  type: 'image_url';
  url: string;
  alt: string;
  width?: number;
  height?: number;
  format?: string;
}

interface ConfigContent {
  type: 'config';
  value: boolean | number | string | object;
  schema?: JSONSchema;        // Validation schema
}

interface TranslationContent {
  type: 'translation';
  values: Record<string, string>; // locale -> translated text
  defaultLocale: string;
}
```

### ContentMetadata
**Purpose**: Tracks changes and management information
```typescript
interface ContentMetadata {
  createdAt: string;           // ISO timestamp
  createdBy: string;           // User ID
  updatedAt: string;           // ISO timestamp
  updatedBy: string;           // User ID
  version: number;             // Incremental version
  tags?: string[];             // Searchable tags
  notes?: string;              // Editorial notes
}
```

### ContentVersion
**Purpose**: Tracks changes to content items over time
```typescript
interface ContentVersion {
  id: string;                  // Unique version identifier
  contentId: string;           // Parent content item
  version: number;             // Version number
  value: ContentValue;         // Content at this version
  metadata: ContentMetadata;   // Version-specific metadata
  changeLog: string;           // Description of changes
}
```

## State Management

### Content Slice State
```typescript
interface ContentState {
  // UI state
  selectedLanguage: string;
  contentPreferences: {
    showDrafts: boolean;
    compactView: boolean;
  };

  // Content system state
  lastSyncTime: string;
  syncStatus: 'idle' | 'syncing' | 'error';
  fallbackMode: boolean;       // Using local fallbacks

  // Local fallback cache
  fallbackContent: Record<string, ContentValue>;
}
```

### RTK Query Endpoints State
```typescript
// Managed automatically by RTK Query
interface ContentApiState {
  queries: {
    getContentItems: QueryState<ContentItem[]>;
    getContentItem: QueryState<ContentItem>;
    getContentByCategory: QueryState<ContentItem[]>;
    // ... other endpoint states
  };
  mutations: {
    updateContentItem: MutationState<ContentItem>;
    publishContent: MutationState<ContentItem>;
    // ... other mutation states
  };
}
```

## Relationships

### Content Hierarchy
- ContentCategory (1) → ContentItem (many)
- ContentItem (1) → ContentVersion (many)
- ContentItem references other ContentItems via IDs

### User Permissions
- User roles determine ContentCategory access
- Content operations validated against CategoryPermissions
- Version history preserves permission context

### Localization
- TranslationContent contains multiple locale values
- ContentItem can reference translation tables
- Fallback chain: specific locale → default locale → fallback content

## Validation Rules

### Content ID Format
- Pattern: `{category}.{subcategory}.{key}`
- Examples: `navigation.header.title`, `error.network.timeout`, `help.onboarding.step1`
- Must be unique across all content items
- Cannot contain special characters except dots and underscores

### Content Value Validation
- Text content: Length limits, encoding validation
- Rich text: HTML sanitization, allowed tags enforcement
- Image URLs: URL format validation, domain whitelist
- Config values: JSON schema validation if provided
- Translations: Required default locale, consistent keys

### Version Management
- Version numbers must be sequential
- Cannot delete published versions
- Draft versions can be modified, published versions cannot
- Maximum version history configurable per category

### Category Hierarchy
- Maximum nesting depth: 3 levels
- Circular references prevented
- Parent category deletion requires child migration
- Permission inheritance from parent categories

## Performance Considerations

### Indexing Strategy
- Primary index: ContentItem.id
- Secondary indexes: category, status, updatedAt
- Text search index: value content (for searchable content)
- Composite index: (category, status) for filtered queries

### Caching Layers
- RTK Query cache: Recent content items (5-minute TTL)
- Browser localStorage: Fallback content (persistent)
- Memory cache: Frequently accessed content (session-based)
- CDN cache: Static media content (long TTL)

### Data Loading Patterns
- Critical content: Eager loading on app initialization
- Page content: Lazy loading on route change
- Optional content: On-demand loading when accessed
- Media content: Progressive loading with placeholders

This data model supports all functional requirements while maintaining performance, security, and maintainability for the dynamic content management system.