# Message Types System Documentation

## Overview

The Message Types System provides 13 sophisticated React components for rendering different types of interactive messages in the Tchat application. This system enables rich, engaging conversations beyond simple text messages.

## Architecture

### Core Components

- **MessageComponentFactory**: Dynamic component loading with error boundaries
- **MessageRegistry**: Centralized component management and lifecycle
- **Message Components**: 13 specialized components for different message types
- **API Integration**: RTK Query endpoints for backend communication
- **Redux Store**: State management for UI interactions and caching
- **Custom Hooks**: Optimized hooks for common operations

### Design Principles

1. **Performance First**: Lazy loading, code splitting, and optimized rendering
2. **Accessibility**: WCAG 2.1 AA compliance with full keyboard navigation
3. **Cross-Platform**: >95% visual consistency with mobile apps
4. **Error Resilience**: Comprehensive error boundaries and fallback UI
5. **Type Safety**: Full TypeScript coverage with discriminated unions

## Message Types

### 1. ReplyMessage (`MessageType.REPLY`)

**Purpose**: Thread-based conversations with nested reply support

**Features**:
- Visual thread lines with depth indication
- Jump to original message functionality
- Compact mode for deep threads
- Thread context preservation

**Content Structure**:
```typescript
interface ReplyContent {
  originalMessageId: string;
  replyText: string;
  threadDepth: number;
  isThreadStart: boolean;
}
```

**Usage Example**:
```tsx
<ReplyMessage
  message={replyMessage}
  onInteraction={handleInteraction}
  onReplyToReply={handleNestedReply}
  showThreadLine={true}
  maxDepth={5}
/>
```

### 2. QuizMessage (`MessageType.QUIZ`)

**Purpose**: Interactive quizzes with multiple question types

**Features**:
- Multiple question types (multiple choice, true/false, text input)
- Time limits with countdown timer
- Real-time scoring and feedback
- Detailed explanations for each answer
- Retake functionality

**Content Structure**:
```typescript
interface QuizContent {
  title: string;
  description: string;
  questions: QuizQuestion[];
  timeLimit: number; // seconds
  showResults: boolean;
  allowRetake: boolean;
  userAnswers?: UserAnswer[];
  completedAt?: string;
}
```

**Usage Example**:
```tsx
<QuizMessage
  message={quizMessage}
  onQuizComplete={handleQuizCompletion}
  onAnswerChange={handleAnswerChange}
  showProgress={true}
/>
```

### 3. EventMessage (`MessageType.EVENT`)

**Purpose**: Event management with RSVP functionality

**Features**:
- RSVP with attending/maybe/not attending options
- Calendar integration
- Attendee management
- Reminder setting
- Location and time display

**Content Structure**:
```typescript
interface EventContent {
  title: string;
  description: string;
  startDate: string;
  endDate: string;
  location: string;
  maxAttendees?: number;
  attendees: Attendee[];
  rsvpRequired: boolean;
  userRsvpStatus?: 'attending' | 'maybe' | 'not_attending';
  reminderSet: boolean;
}
```

### 4. SurveyMessage (`MessageType.SURVEY`)

**Purpose**: Polls and surveys with real-time results

**Features**:
- Multiple poll types (single choice, multiple choice, rating)
- Real-time result visualization
- Anonymous voting option
- Vote changing capability
- Results analysis

### 5. ProductMessage (`MessageType.PRODUCT`)

**Purpose**: E-commerce integration with product showcase

**Features**:
- Product image carousel
- Variant selection (color, size, etc.)
- Add to cart functionality
- Reviews and ratings display
- Stock status indication

### 6. FormMessage (`MessageType.FORM`)

**Purpose**: Dynamic forms with validation

**Features**:
- Dynamic field types
- Real-time validation
- Multi-step forms
- Conditional field logic
- File upload support

### 7. RichCardMessage (`MessageType.RICH_CARD`)

**Purpose**: Rich media cards with interactive elements

**Features**:
- Carousel and grid layouts
- Media controls for video/audio
- Interactive buttons and links
- Card stacking and grouping

### 8. StatusUpdateMessage (`MessageType.STATUS_UPDATE`)

**Purpose**: Activity tracking and status updates

**Features**:
- Presence indicators
- Activity timeline
- Social interactions (likes, comments)
- Status change notifications

### 9. VenueMessage (`MessageType.VENUE`)

**Purpose**: Location-based content with business info

**Features**:
- Interactive maps
- Business information display
- Booking integration
- Reviews and photos

### 10. SpreadsheetMessage (`MessageType.SPREADSHEET`)

**Purpose**: Data tables with collaborative editing

**Features**:
- Sorting and filtering
- Cell editing
- Formula support
- Collaborative editing
- Export functionality

### 11. EmbedMessage (`MessageType.EMBED`)

**Purpose**: Rich embeds with link previews

**Features**:
- Link preview generation
- Media embedding (videos, tweets, etc.)
- Social metrics display
- Custom embed types

### 12. ScheduleMessage (`MessageType.SCHEDULE`)

**Purpose**: Calendar and scheduling integration

**Features**:
- Calendar view
- Time slot selection
- Meeting scheduling
- Availability checking

### 13. DocumentMessage (`MessageType.DOCUMENT`)

**Purpose**: Document sharing with preview

**Features**:
- Document preview
- Version control
- Collaborative editing
- Comment system
- Download management

## API Integration

### RTK Query Endpoints

```typescript
// Message CRUD operations
useCreateMessageMutation()
useUpdateMessageMutation()
useDeleteMessageMutation()

// Message interactions
useInteractWithMessageMutation()

// Data fetching
useGetMessagesQuery()
useGetMessageQuery()
useGetThreadQuery()

// Search and analytics
useSearchMessagesQuery()
useGetMessageAnalyticsQuery()

// Validation
useValidateMessageMutation()
```

### Custom Hooks

```typescript
// High-level operations
useCreateMessage() // With validation and error handling
useMessageInteraction() // With optimistic updates
useMessages() // Paginated loading with infinite scroll
useMessageSearch() // Debounced search with filters
useMessageAnalytics() // Performance metrics and usage stats

// Utility hooks
useMessagePerformance() // Performance monitoring
useMessageUpdates() // Real-time WebSocket updates
useBulkMessageOperations() // Bulk operations
```

## Performance Optimization

### Code Splitting

All message components are dynamically imported to reduce initial bundle size:

```typescript
const ReplyMessage = lazy(() => import('./ReplyMessage'));
const QuizMessage = lazy(() => import('./QuizMessage'));
// ... other components
```

### Caching Strategy

- **Component-level caching**: React.memo for pure components
- **API caching**: RTK Query with tag-based invalidation
- **Registry caching**: Component metadata and performance metrics
- **Browser caching**: Optimized asset loading

### Performance Budgets

- **Bundle size**: <50KB per component
- **Load time**: <200ms for component loading
- **Render time**: <16ms for 60fps animations
- **Memory usage**: <10MB per message component

## Error Handling

### Error Boundaries

```typescript
<MessageErrorBoundary messageType={message.type} messageId={message.id}>
  <Suspense fallback={<MessageLoadingFallback />}>
    <MessageComponent message={message} />
  </Suspense>
</MessageErrorBoundary>
```

### Fallback UI

- **Component loading errors**: Generic message with retry option
- **API errors**: User-friendly error messages with retry
- **Validation errors**: Inline field-level error display
- **Network errors**: Offline mode with cached data

## Accessibility

### WCAG 2.1 AA Compliance

- **Keyboard Navigation**: Full tab order and arrow key support
- **Screen Readers**: Proper ARIA labels and live regions
- **Color Contrast**: 4.5:1 minimum contrast ratio
- **Focus Management**: Visible focus indicators and logical tab order

### Keyboard Shortcuts

- `Tab` / `Shift+Tab`: Navigate between interactive elements
- `Enter` / `Space`: Activate buttons and select options
- `Arrow keys`: Navigate within components (quiz options, etc.)
- `Escape`: Close modals and cancel operations

## Testing

### Test Coverage

- **Unit Tests**: 95%+ coverage for all components
- **Integration Tests**: API endpoint testing with MSW
- **E2E Tests**: User workflow testing with Playwright
- **Accessibility Tests**: Automated a11y testing with axe-core
- **Performance Tests**: Load time and render performance validation

### Testing Strategy

1. **Contract Tests**: API schema validation
2. **Component Tests**: Isolated component behavior
3. **Integration Tests**: Component + API interaction
4. **E2E Tests**: Full user workflows
5. **Visual Tests**: Cross-browser screenshot comparison

## Development Guide

### Adding New Message Types

1. **Define TypeScript interfaces**:
```typescript
// Add to MessageType enum
export enum MessageType {
  // ... existing types
  NEW_TYPE = 'new_type',
}

// Define content interface
export interface NewTypeContent {
  title: string;
  customField: any;
}
```

2. **Create React component**:
```typescript
export const NewTypeMessage: React.FC<NewTypeMessageProps> = ({
  message,
  onInteraction,
  ...props
}) => {
  // Component implementation
};
```

3. **Register component**:
```typescript
messageRegistry.register(MessageType.NEW_TYPE, {
  component: NewTypeMessage,
  displayName: 'New Type Message',
  category: MessageComponentCategory.INTERACTIVE,
  priority: 5,
  supportLevel: MessageSupportLevel.FULL,
  metadata: {
    // Component metadata
  },
});
```

4. **Add API support**:
```typescript
// Update API endpoints to handle new type
// Add validation rules
// Update search facets
```

5. **Write tests**:
```typescript
// Unit tests
// Integration tests
// E2E tests
```

### Best Practices

1. **Component Design**:
   - Use compound component pattern for complex UIs
   - Implement proper loading and error states
   - Follow accessibility guidelines
   - Optimize for performance

2. **State Management**:
   - Use local state for UI-only concerns
   - Use Redux for cross-component state
   - Implement optimistic updates for better UX
   - Handle offline scenarios

3. **Error Handling**:
   - Always provide fallback UI
   - Log errors for monitoring
   - Give users recovery options
   - Test error scenarios

4. **Performance**:
   - Lazy load components
   - Memoize expensive calculations
   - Implement virtual scrolling for large lists
   - Monitor and optimize render performance

## Configuration

### Environment Variables

```bash
# API Configuration
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_WS_URL=ws://localhost:8080/ws

# Feature Flags
VITE_ENABLE_MESSAGE_ANALYTICS=true
VITE_ENABLE_REAL_TIME_UPDATES=true

# Performance
VITE_MESSAGE_CACHE_SIZE=1000
VITE_COMPONENT_PRELOAD_COUNT=5
```

### Registry Configuration

```typescript
// Configure message registry
messageRegistry.configure({
  cacheSize: 1000,
  preloadPriority: 7,
  performanceMonitoring: true,
  analyticsEnabled: true,
});
```

## Migration Guide

### From Legacy Message System

1. **Update message data structure**:
   - Add `type` field to all messages
   - Migrate content to typed interfaces
   - Update API responses

2. **Replace message rendering**:
   - Use MessageComponentFactory instead of switch statements
   - Update CSS classes to new design system
   - Add error boundaries

3. **Update state management**:
   - Migrate to RTK Query for API calls
   - Update Redux state structure
   - Implement new caching strategy

## Troubleshooting

### Common Issues

1. **Component not loading**:
   - Check network tab for failed imports
   - Verify component is registered in factory
   - Check error boundary logs

2. **Performance issues**:
   - Monitor performance metrics in registry
   - Check for memory leaks in DevTools
   - Verify lazy loading is working

3. **Accessibility problems**:
   - Run axe-core accessibility audit
   - Test with screen reader
   - Verify keyboard navigation

### Debug Tools

- **React DevTools**: Component state and props
- **Redux DevTools**: State management debugging
- **Network Tab**: API call monitoring
- **Performance Tab**: Render performance analysis

## Roadmap

### Planned Features

1. **AI-Powered Message Types**:
   - Smart form generation
   - Intelligent quiz creation
   - Content recommendations

2. **Enhanced Analytics**:
   - User engagement heatmaps
   - A/B testing framework
   - Performance insights dashboard

3. **Advanced Accessibility**:
   - Voice navigation support
   - High contrast theme
   - Reading level optimization

4. **Mobile Optimizations**:
   - Touch gesture support
   - Haptic feedback
   - Offline-first architecture

## Support

For questions, issues, or contributions:

- **Documentation**: `/docs/message-types/`
- **Examples**: `/examples/message-types/`
- **Tests**: `/tests/message-types/`
- **Issue Tracker**: GitHub Issues
- **Development Chat**: #message-types-dev