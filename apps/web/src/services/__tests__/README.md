# Content API Contract Tests

## Overview

This directory contains contract tests for the Content Management System API endpoints. These tests follow Test-Driven Development (TDD) principles and are designed to fail until proper implementation is completed.

## Test ID Implementation Guidelines

When implementing the actual content API service and UI components, use these semantic test IDs:

### Content Categories Components

```typescript
// Category List Component
<div data-testid="content-categories-list">
  <div data-testid="content-category-item" data-category-id="navigation">
    <h3 data-testid="content-category-name">Navigation</h3>
    <p data-testid="content-category-description">Navigation menu items and links</p>
    <div data-testid="content-category-permissions">
      <span data-testid="content-category-read-permissions">Read: user, admin, guest</span>
      <span data-testid="content-category-write-permissions">Write: admin, editor</span>
      <span data-testid="content-category-publish-permissions">Publish: admin</span>
    </div>
  </div>
</div>

// Category Management Form
<form data-testid="content-category-form">
  <input
    data-testid="content-category-id-input"
    name="categoryId"
    placeholder="Category ID"
  />
  <input
    data-testid="content-category-name-input"
    name="categoryName"
    placeholder="Category Name"
  />
  <textarea
    data-testid="content-category-description-input"
    name="categoryDescription"
    placeholder="Category Description"
  />
  <button data-testid="content-category-submit-button" type="submit">
    Create Category
  </button>
</form>

// Empty State
<div data-testid="content-categories-empty-state">
  <p data-testid="content-categories-empty-message">No categories found</p>
  <button data-testid="content-categories-create-first-button">
    Create First Category
  </button>
</div>

// Loading State
<div data-testid="content-categories-loading">
  <div data-testid="content-categories-loading-spinner" />
  <p data-testid="content-categories-loading-message">Loading categories...</p>
</div>

// Error State
<div data-testid="content-categories-error">
  <p data-testid="content-categories-error-message">Failed to load categories</p>
  <button data-testid="content-categories-retry-button">Retry</button>
</div>
```

### API Interaction Test IDs

```typescript
// For E2E testing of API interactions
<div data-testid="content-api-status" data-status="loading|success|error" />
<div data-testid="content-api-response-time" data-time="123ms" />
<div data-testid="content-api-error-code" data-error-code="404|500|401" />
```

## Test Implementation Status

### ✅ Completed
- **T007**: Contract test for getContentCategories API endpoint
- Category structure validation
- Permissions structure validation
- Error scenario handling
- RTK Query integration patterns

### ⏳ Pending Implementation
- Content API service (`src/services/content.ts`)
- Backend API endpoint (`/api/content/categories`)
- Category management UI components
- E2E tests with actual user workflows

## Running Tests

```bash
# Run all content API tests
npm test -- content-api

# Run specific contract test
npm test -- content-api-get-categories.test.ts

# Watch mode for development
npm test -- --watch content-api-get-categories.test.ts
```

## TDD Workflow

1. **Red**: Tests fail (current state) ❌
2. **Green**: Implement minimal code to pass tests ✅
3. **Refactor**: Improve code while keeping tests passing 🔄

### Next Implementation Steps

1. Create `src/services/content.ts` with RTK Query endpoints
2. Implement `getContentCategories` query
3. Add proper TypeScript types and error handling
4. Create backend API endpoint
5. Build UI components with test IDs
6. Write E2E tests for complete user workflows

## Test ID Standards

- **Pattern**: `[component]-[element]-[action]`
- **Examples**:
  - `content-category-name` (display element)
  - `content-category-edit-button` (action button)
  - `content-categories-list` (container)
  - `content-category-permissions-read` (specific data)

### Test ID Naming Conventions

- Use kebab-case for all test IDs
- Include component context (`content-category-`)
- Be specific about element purpose (`-name`, `-button`, `-input`)
- Include state information when relevant (`-loading`, `-error`)
- Use data attributes for dynamic values (`data-category-id="navigation"`)

## Error Scenarios Coverage

- ✅ 404 Not Found responses
- ✅ 500 Internal Server Error responses
- ✅ 401 Unauthorized responses
- ✅ Network errors and timeouts
- ✅ Malformed JSON responses
- ✅ Empty data handling
- ✅ Type validation and structure checking

## Performance Considerations

- RTK Query caching and deduplication
- Optimistic updates for better UX
- Proper error boundaries and fallback states
- Request timeout handling (10 second limit)
- Retry logic with exponential backoff