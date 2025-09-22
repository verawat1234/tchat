# Data Model: Common UI Components

**Feature**: Common UI Component Breakdown
**Date**: 2025-09-21
**Status**: Complete

## Component Interfaces & Data Structures

### 1. Core Component Properties

#### BaseComponent Interface
```typescript
interface BaseComponent {
  className?: string;
  children?: React.ReactNode;
  testId?: string;
  variant?: string;
  size?: 'sm' | 'md' | 'lg' | 'xl';
}
```

#### Theme Tokens Structure
```typescript
interface ThemeTokens {
  colors: {
    primary: string;
    secondary: string;
    accent: string;
    muted: string;
    background: string;
    foreground: string;
    border: string;
  };
  spacing: {
    xs: string;
    sm: string;
    md: string;
    lg: string;
    xl: string;
  };
  typography: {
    fontSize: Record<string, string>;
    fontWeight: Record<string, number>;
    lineHeight: Record<string, string>;
  };
}
```

### 2. Pagination Component

#### Data Structure
```typescript
interface PaginationProps extends BaseComponent {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  showFirstLast?: boolean;
  showPrevNext?: boolean;
  maxVisiblePages?: number;
  disabled?: boolean;
}

interface PaginationState {
  visiblePages: number[];
  canGoPrevious: boolean;
  canGoNext: boolean;
}
```

#### Behavior Rules
- Current page must be >= 1 and <= totalPages
- Visible pages calculated based on maxVisiblePages (default: 5)
- Disable navigation when disabled prop is true
- Fire onPageChange only for valid page numbers

### 3. Tabs Component

#### Data Structure
```typescript
interface TabsProps extends BaseComponent {
  defaultValue?: string;
  value?: string;
  onValueChange?: (value: string) => void;
  orientation?: 'horizontal' | 'vertical';
  activationMode?: 'automatic' | 'manual';
}

interface TabsListProps extends BaseComponent {
  loop?: boolean;
}

interface TabsTriggerProps extends BaseComponent {
  value: string;
  disabled?: boolean;
}

interface TabsContentProps extends BaseComponent {
  value: string;
  forceMount?: boolean;
}
```

#### State Transitions
- Controlled mode: value prop controls active tab
- Uncontrolled mode: defaultValue sets initial state
- Manual activation: requires click to activate
- Automatic activation: activates on focus

### 4. Card Component

#### Data Structure
```typescript
interface CardProps extends BaseComponent {
  variant?: 'default' | 'outline' | 'ghost';
  padding?: 'none' | 'sm' | 'md' | 'lg';
  shadow?: 'none' | 'sm' | 'md' | 'lg';
  interactive?: boolean;
  onClick?: () => void;
}

interface CardHeaderProps extends BaseComponent {
  title?: string;
  subtitle?: string;
  actions?: React.ReactNode;
}

interface CardContentProps extends BaseComponent {
  padding?: 'inherit' | 'none' | 'sm' | 'md' | 'lg';
}

interface CardFooterProps extends BaseComponent {
  justify?: 'start' | 'center' | 'end' | 'between';
}
```

#### Composition Pattern
```typescript
// Compound component structure
Card.Root
Card.Header
Card.Content
Card.Footer
```

### 5. ChatMessage Component

#### Data Structure
```typescript
interface ChatMessageProps extends BaseComponent {
  type: 'text' | 'image' | 'file' | 'system' | 'typing';
  content: string | MediaContent | SystemContent;
  timestamp: Date;
  sender?: {
    id: string;
    name: string;
    avatar?: string;
  };
  isOwn?: boolean;
  status?: 'sending' | 'sent' | 'delivered' | 'read' | 'failed';
  reactions?: Reaction[];
  reply?: ChatMessageProps;
}

interface MediaContent {
  url: string;
  type: 'image' | 'video' | 'audio' | 'file';
  filename?: string;
  size?: number;
  thumbnail?: string;
}

interface SystemContent {
  type: 'join' | 'leave' | 'rename' | 'notification';
  message: string;
  metadata?: Record<string, any>;
}

interface Reaction {
  emoji: string;
  count: number;
  users: string[];
  hasUserReacted: boolean;
}
```

#### Message Type Behaviors
- **Text**: Rich text with markdown support
- **Image**: Thumbnail with lightbox on click
- **File**: Download link with file info
- **System**: Centered message with muted styling
- **Typing**: Animated indicator

### 6. Badge Component

#### Data Structure
```typescript
interface BadgeProps extends BaseComponent {
  variant?: 'default' | 'success' | 'warning' | 'danger' | 'info';
  size?: 'sm' | 'md' | 'lg';
  dot?: boolean;
  pulse?: boolean;
  count?: number;
  max?: number;
  showZero?: boolean;
}
```

#### Display Rules
- Count badges show number when count > 0 (unless showZero=true)
- Max count displays "99+" when count > max (default: 99)
- Dot variant shows simple indicator
- Pulse adds animation for attention

### 7. Layout Components

#### Data Structure
```typescript
interface LayoutProps extends BaseComponent {
  direction?: 'row' | 'column';
  gap?: 'none' | 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  align?: 'start' | 'center' | 'end' | 'stretch';
  justify?: 'start' | 'center' | 'end' | 'between' | 'around' | 'evenly';
  wrap?: boolean;
}

interface GridProps extends BaseComponent {
  cols?: number | 'auto';
  rows?: number | 'auto';
  gap?: 'none' | 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  responsive?: {
    sm?: Partial<GridProps>;
    md?: Partial<GridProps>;
    lg?: Partial<GridProps>;
  };
}

interface StackProps extends BaseComponent {
  space?: 'none' | 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  divider?: React.ReactNode;
}
```

### 8. Header Component

#### Data Structure
```typescript
interface HeaderProps extends BaseComponent {
  title: string;
  subtitle?: string;
  level?: 1 | 2 | 3 | 4 | 5 | 6;
  actions?: React.ReactNode;
  breadcrumbs?: BreadcrumbItem[];
  sticky?: boolean;
  border?: boolean;
}

interface BreadcrumbItem {
  label: string;
  href?: string;
  onClick?: () => void;
  disabled?: boolean;
}
```

### 9. Sidebar Component

#### Data Structure
```typescript
interface SidebarProps extends BaseComponent {
  position?: 'left' | 'right';
  width?: number | string;
  collapsible?: boolean;
  collapsed?: boolean;
  onCollapsedChange?: (collapsed: boolean) => void;
  overlay?: boolean;
  persistent?: boolean;
}

interface SidebarContentProps extends BaseComponent {
  padding?: 'none' | 'sm' | 'md' | 'lg';
}

interface SidebarNavProps extends BaseComponent {
  items: SidebarNavItem[];
  value?: string;
  onValueChange?: (value: string) => void;
}

interface SidebarNavItem {
  id: string;
  label: string;
  icon?: React.ReactNode;
  href?: string;
  onClick?: () => void;
  disabled?: boolean;
  badge?: string | number;
  children?: SidebarNavItem[];
}
```

## Component Relationships

### Composition Hierarchy
```
Layout Components (Grid, Stack)
├── Card Components
│   ├── Header
│   ├── Content
│   └── Footer
├── ChatMessage Components
├── Navigation Components
│   ├── Tabs
│   ├── Pagination
│   └── Sidebar
└── Utility Components
    ├── Badge
    └── Header
```

### Shared Dependencies
- All components extend BaseComponent interface
- Theme tokens accessed via CSS custom properties
- Accessibility attributes (ARIA) built into base components
- Event handlers follow React conventions (onEventName)

## Validation Rules

### Type Safety
- All props strongly typed with TypeScript
- Required props marked appropriately
- Union types for enumerated values
- Generic types for flexible content

### Runtime Validation
- Page numbers within valid range (Pagination)
- Tab values exist in tab list (Tabs)
- Message content matches type (ChatMessage)
- Layout values are valid CSS values

### Accessibility Requirements
- Keyboard navigation support
- Screen reader compatibility
- Focus management
- ARIA attributes and roles
- Color contrast compliance

## Performance Considerations

### Bundle Size Optimization
- Tree-shakeable exports
- Dynamic imports for heavy components
- Minimal external dependencies
- Optimized re-renders with React.memo

### Runtime Performance
- Virtualization for large lists (when applicable)
- Debounced event handlers
- Memoized calculations
- Efficient state updates

## Migration Strategy

### Backward Compatibility
- Existing components remain functional
- New components opt-in only
- Gradual migration path
- Clear deprecation warnings

### Component Mapping
- Identify existing patterns → new components
- Provide migration guides
- Automated refactoring tools where possible
- Version management strategy