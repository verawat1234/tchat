# Quickstart: Dynamic Content Management

This quickstart guide validates the core user scenarios from the feature specification through hands-on testing.

## Prerequisites

- Tchat application running with RTK infrastructure (009-create-rtk-follow)
- Backend content API endpoints implemented
- Content database with sample data

## Test Scenario 1: View Dynamic Content

**User Story**: As a user visiting the application, I want to see current, accurate information on all pages.

### Steps
1. Open the Tchat application in a browser
2. Navigate to the home page
3. Verify page content loads from API (not hardcoded values)
4. Navigate to different pages (Chat, Store, Social, Video, More)
5. Verify each page displays dynamic content

### Expected Results
- All text content comes from RTK Query cache
- Page load time is under 200ms for content
- No hardcoded text visible in UI
- Content appears consistent across page refreshes

### Validation Commands
```bash
# Run the application
cd apps/web
npm run dev

# In browser developer tools, check Network tab:
# - Should see API calls to /content/items
# - Should see content data in responses
# - Redux DevTools should show content state

# Check Redux store state contains content data
console.log(store.getState().api.queries); // Should show content queries
```

## Test Scenario 2: Real-time Content Updates

**User Story**: As a content manager, I want to update application content dynamically so that changes are immediately visible to users.

### Steps
1. Open Tchat application in Browser A
2. Open content management interface in Browser B (or API client)
3. Update a visible text element (e.g., navigation title)
4. Verify update appears in Browser A without refresh

### Expected Results
- Content changes propagate to user interface
- No page refresh required for content updates
- Update occurs within 5 seconds of change
- All instances of the content item update simultaneously

### Validation Commands
```bash
# Update content via API (simulating content manager action)
curl -X PUT http://localhost:3001/api/content/items/navigation.header.title \
  -H "Content-Type: application/json" \
  -d '{
    "value": {
      "type": "text",
      "value": "Updated Header Title"
    },
    "notes": "Quickstart test update"
  }'

# Verify RTK Query cache invalidation
# Check browser Redux DevTools for cache updates
```

## Test Scenario 3: Fallback Content Display

**User Story**: As a user, when there are network issues or data source problems, I want to see appropriate fallback content instead of broken or missing information.

### Steps
1. Open Tchat application with network connectivity
2. Load content normally to populate cache
3. Simulate network disconnection or API failure
4. Navigate to different pages
5. Verify fallback content displays correctly

### Expected Results
- Application continues to function during network issues
- Cached content displays when available
- Hardcoded fallback content displays when cache empty
- No broken UI elements or missing text
- Clear indication when using fallback content

### Validation Commands
```bash
# Test network failure simulation
# In browser DevTools:
# 1. Go to Network tab
# 2. Set throttling to "Offline"
# 3. Refresh application
# 4. Verify fallback content displays

# Check fallback content in Redux state
console.log(store.getState().content.fallbackContent);
```

## Test Scenario 4: Content Loading Performance

**User Story**: As a user viewing content-heavy pages, when dynamic content loads, I want page performance to remain acceptable.

### Steps
1. Clear browser cache
2. Navigate to content-heavy page (e.g., Social tab with multiple sections)
3. Measure content loading time
4. Check for content loading indicators
5. Verify progressive content display

### Expected Results
- Initial content appears within 200ms
- Progressive loading for non-critical content
- Loading indicators show during content fetch
- No layout shifts during content loading
- Smooth user experience throughout

### Validation Commands
```bash
# Performance measurement in browser DevTools
# 1. Go to Performance tab
# 2. Record page load
# 3. Analyze content loading timeline

# Check RTK Query loading states
console.log(store.getState().api.queries); // Look for isLoading flags

# Verify content loading performance
performance.mark('content-start');
// Navigate to page
performance.mark('content-loaded');
performance.measure('content-load-time', 'content-start', 'content-loaded');
console.log(performance.getEntriesByName('content-load-time'));
```

## Integration Test Scenarios

### Test 1: Content Category Organization
```bash
# Verify content categories work correctly
curl http://localhost:3001/api/content/categories
# Should return: navigation, errors, help, social, etc.

curl http://localhost:3001/api/content/categories/navigation/items
# Should return all navigation-related content
```

### Test 2: Content Type Support
```bash
# Test different content types
curl http://localhost:3001/api/content/items/navigation.header.title
# Should return: text content

curl http://localhost:3001/api/content/items/social.profile.avatar
# Should return: image_url content

curl http://localhost:3001/api/content/items/app.feature.darkMode
# Should return: config content (boolean)
```

### Test 3: Content Versioning
```bash
# Update content and check version history
curl -X PUT http://localhost:3001/api/content/items/navigation.header.title \
  -H "Content-Type: application/json" \
  -d '{"value": {"type": "text", "value": "New Title"}, "notes": "Version test"}'

curl http://localhost:3001/api/content/items/navigation.header.title/versions
# Should return: version history with metadata
```

## Error Handling Tests

### Test 1: Invalid Content Updates
```bash
# Test validation errors
curl -X PUT http://localhost:3001/api/content/items/navigation.header.title \
  -H "Content-Type: application/json" \
  -d '{"value": {"type": "invalid_type", "value": ""}}'
# Should return: 400 Bad Request
```

### Test 2: Non-existent Content
```bash
# Test 404 handling
curl http://localhost:3001/api/content/items/non.existent.content
# Should return: 404 Not Found

# Verify graceful handling in UI
# Navigate to page expecting non-existent content
# Should show: fallback content, not error
```

### Test 3: Network Timeout Handling
```bash
# Simulate slow network in DevTools
# Set throttling to "Slow 3G"
# Verify: Loading states show appropriately
# Verify: Timeout errors handled gracefully
```

## Acceptance Criteria Validation

### ✅ Functional Requirements Check
- [ ] FR-001: All page content from centralized sources ✓
- [ ] FR-002: Authorized content updates without code changes ✓
- [ ] FR-003: Updated content visible across all pages ✓
- [ ] FR-004: Fallback content when data unavailable ✓
- [ ] FR-005: Content consistency across pages ✓
- [ ] FR-006: Real-time/near real-time updates ✓
- [ ] FR-007: Multiple content types supported ✓
- [ ] FR-008: Content version history maintained ✓
- [ ] FR-009: Content validation before going live ✓
- [ ] FR-010: Content preview capability ✓

### ✅ Performance Requirements Check
- [ ] Content load time under 200ms ✓
- [ ] No layout shifts during content loading ✓
- [ ] Smooth updates without page refresh ✓
- [ ] Acceptable performance with network issues ✓

### ✅ User Experience Check
- [ ] No broken UI during content failures ✓
- [ ] Clear feedback during content operations ✓
- [ ] Consistent content across all page instances ✓
- [ ] Intuitive content management workflow ✓

## Troubleshooting Common Issues

### Content Not Loading
1. Check Redux DevTools for API query states
2. Verify network requests in DevTools Network tab
3. Check browser console for JavaScript errors
4. Verify content API endpoints are accessible

### Content Not Updating
1. Check RTK Query cache invalidation
2. Verify content API mutation responses
3. Check for race conditions in content updates
4. Verify proper tag-based cache invalidation

### Performance Issues
1. Monitor content loading waterfall in DevTools
2. Check for unnecessary content refetches
3. Verify content caching configuration
4. Analyze bundle size impact of content system

### Fallback Content Issues
1. Verify fallback content configuration
2. Check localStorage for cached content
3. Test offline scenarios thoroughly
4. Verify error boundary implementations

## Success Metrics

After completing this quickstart, you should observe:
- 100% dynamic content across all pages
- <200ms content loading performance
- Zero hardcoded text in production build
- Robust error handling and fallback behavior
- Seamless content update experience

The system is ready for production when all test scenarios pass and acceptance criteria are met.