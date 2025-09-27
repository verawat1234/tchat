# Chat System Test ID Implementation Guide

## Overview

This document provides comprehensive guidelines for implementing and using test IDs (ucid - unique component identifiers) in the Tchat chat system. It follows the **component-element-action** pattern for systematic test automation and quality assurance.

## Test ID Naming Convention

### Pattern Structure
```
[component]-[element]-[action]
```

### Components
- `chat-detail`: Main chat detail screen
- `chat-input`: Message input area and controls
- `message`: Message bubble and content rendering
- `attachment`: Attachment menu and file options
- `dialog`: Modal dialogs and confirmations

### Elements
- `screen-container`: Main screen wrapper
- `header-name`: Chat participant name display
- `header-status`: Online/offline status indicator
- `avatar-image`: User avatar display
- `back-button`: Navigation back control
- `video-call-button`: Video call initiation
- `audio-call-button`: Audio call initiation
- `more-options-button`: More options menu trigger
- `messages-list`: Messages container/list
- `message-item-{id}`: Individual message with dynamic ID
- `input-bar`: Input area container
- `message-field`: Text input field
- `send-button`: Send message control
- `attachment-button`: Attachment menu trigger
- `attachment-menu`: Attachment options container
- `option-{type}`: Specific attachment types (photo, video, file, location)

### Actions (Implied by Test Type)
- **Click Actions**: Button presses, menu selections
- **Type Actions**: Text input, form filling
- **Verify Actions**: Content validation, state checking
- **Navigate Actions**: Screen transitions, routing

## Implemented Test IDs

### Chat Detail Screen
```kotlin
// Main container
"chat-detail-screen-container"

// Header elements
"chat-detail-avatar-image"
"chat-detail-header-name"
"chat-detail-header-status"
"chat-detail-back-button"
"chat-detail-video-call-button"
"chat-detail-audio-call-button"
"chat-detail-more-options-button"

// Messages area
"chat-detail-messages-list"
"chat-detail-message-item-{messageId}"
```

### Chat Input Components
```kotlin
// Input bar
"chat-input-bar-container"
"chat-input-message-field"
"chat-input-send-button"
"chat-input-attachment-button"
"chat-input-attachment-menu"

// Attachment options
"attachment-option-photo"
"attachment-option-video"
"attachment-option-file"
"attachment-option-location"
```

### Message Components
```kotlin
// Message structure
"message-bubble-{messageId}"
"message-card-{messageId}"
"message-content-{messageId}"
"message-timestamp-{messageId}"
"message-status-{messageId}"

// Message type-specific content
"message-text-content"
"message-image-content"
"message-video-content"
"message-audio-content"
"message-file-content"
"message-location-content"
"message-payment-content"
"message-poll-content"
"message-form-content"
"message-system-content"
"message-sticker-content"
"message-gif-content"
"message-contact-content"
"message-event-content"
"message-embed-content"
"message-deleted-content"
"message-default-content"
```

### Dialog Components
```kotlin
// Dialog containers
"chat-video-call-dialog"
"chat-audio-call-dialog"
"chat-more-options-dialog"
"chat-confirmation-dialog"

// Dialog actions
"chat-action-dialog-cancel-button"
"chat-action-dialog-confirm-button"
```

## Test Scenarios Coverage

### 1. Critical Path Scenarios (Priority: CRITICAL)
- **Chat Navigation**: Screen loading, header verification, back navigation
- **Message Sending**: Complete send/receive workflow with validation

### 2. Core Functionality (Priority: HIGH)
- **Attachment Menu**: File sharing options and menu interactions
- **Call Actions**: Video/audio call initiation and confirmation dialogs
- **Message Types**: Various message content rendering verification
- **Error Handling**: Edge cases and failure scenarios

### 3. User Experience (Priority: MEDIUM)
- **More Options**: Chat management features (mute, clear, block, export)
- **Performance**: Response times and system responsiveness

### 4. Quality Assurance (Priority: LOW)
- **Accessibility**: Screen reader compatibility, keyboard navigation
- **Visual Regression**: UI consistency across different states

## Test Implementation Examples

### Example 1: Message Sending Workflow
```kotlin
@Test
fun testMessageSendingWorkflow() {
    // Step 1: Locate input field
    val messageField = composeTestRule.onNodeWithTag("chat-input-message-field")
    messageField.assertIsDisplayed()

    // Step 2: Type message
    messageField.performTextInput("Hello, test message!")

    // Step 3: Verify send button is enabled
    val sendButton = composeTestRule.onNodeWithTag("chat-input-send-button")
    sendButton.assertIsEnabled()

    // Step 4: Send message
    sendButton.performClick()

    // Step 5: Verify message appears in list
    composeTestRule.onNodeWithTag("chat-detail-messages-list")
        .onChildren()
        .assertCountEquals(expectedMessageCount + 1)

    // Step 6: Verify message content
    composeTestRule.onNodeWithTag("message-text-content")
        .assertTextContains("Hello, test message!")
}
```

### Example 2: Attachment Menu Testing
```kotlin
@Test
fun testAttachmentMenuFunctionality() {
    // Step 1: Open attachment menu
    val attachmentButton = composeTestRule.onNodeWithTag("chat-input-attachment-button")
    attachmentButton.performClick()

    // Step 2: Verify menu appears
    val attachmentMenu = composeTestRule.onNodeWithTag("chat-input-attachment-menu")
    attachmentMenu.assertIsDisplayed()

    // Step 3: Test each attachment option
    listOf("photo", "video", "file", "location").forEach { type ->
        composeTestRule.onNodeWithTag("attachment-option-$type")
            .assertIsDisplayed()
            .assertHasClickAction()
    }

    // Step 4: Test option selection
    composeTestRule.onNodeWithTag("attachment-option-photo")
        .performClick()

    // Step 5: Verify menu closes and action occurs
    attachmentMenu.assertDoesNotExist()
    // Verify toast or navigation occurs
}
```

### Example 3: Dialog Interaction Testing
```kotlin
@Test
fun testVideoCallDialogFlow() {
    // Step 1: Trigger video call
    val videoCallButton = composeTestRule.onNodeWithTag("chat-detail-video-call-button")
    videoCallButton.performClick()

    // Step 2: Verify dialog appears
    val dialog = composeTestRule.onNodeWithTag("chat-video-call-dialog")
    dialog.assertIsDisplayed()

    // Step 3: Test cancel action
    val cancelButton = composeTestRule.onNodeWithTag("chat-action-dialog-cancel-button")
    cancelButton.performClick()

    // Step 4: Verify dialog closes
    dialog.assertDoesNotExist()

    // Step 5: Test confirm action
    videoCallButton.performClick()
    val confirmButton = composeTestRule.onNodeWithTag("chat-action-dialog-confirm-button")
    confirmButton.performClick()

    // Step 6: Verify call initiation
    // Check for toast notification or navigation
}
```

## Performance Testing Guidelines

### Response Time Expectations
- **Screen Loading**: < 2 seconds
- **Message Sending**: < 500ms for message to appear
- **UI Interactions**: < 100ms response time
- **Animation Duration**: 200-300ms for smooth UX

### Memory Management
- **Message List**: Efficient recycling for large chat histories
- **Image Loading**: Lazy loading and caching
- **Memory Leaks**: No retained references after navigation

### Network Optimization
- **Message Retry**: Automatic retry with exponential backoff
- **Offline Support**: Queue messages for later sending
- **Data Usage**: Compress images and optimize payload sizes

## Accessibility Testing

### Screen Reader Support
```kotlin
// Ensure all interactive elements have content descriptions
IconButton(
    onClick = onBackClick,
    modifier = Modifier.testTag("chat-detail-back-button")
) {
    Icon(
        Icons.Default.ArrowBack,
        contentDescription = "Navigate back to chat list" // Critical for accessibility
    )
}
```

### Keyboard Navigation
- Tab order should follow logical flow
- All interactive elements should be focusable
- Enter/Space should activate buttons
- Escape should close dialogs

### Color Contrast
- Ensure minimum 4.5:1 contrast ratio for text
- Test with color blindness simulators
- Provide alternative indicators beyond color

## Error Scenario Testing

### Network Failures
1. **No Internet**: Messages queue for later sending
2. **Slow Network**: Show loading indicators, timeout handling
3. **Server Errors**: Display user-friendly error messages

### Input Validation
1. **Empty Messages**: Send button remains disabled
2. **Maximum Length**: Input field handles overflow gracefully
3. **Special Characters**: Proper encoding and display

### Edge Cases
1. **Rapid Clicks**: Prevent duplicate actions
2. **Concurrent Users**: Handle real-time updates
3. **Device Rotation**: Maintain state across orientation changes

## Continuous Integration Setup

### Automated Test Execution
```bash
# Run all chat tests
./gradlew connectedAndroidTest -Pandroid.testInstrumentationRunnerArguments.class=com.tchat.mobile.chat.*

# Run specific test scenario
./gradlew connectedAndroidTest -Pandroid.testInstrumentationRunnerArguments.class=com.tchat.mobile.chat.MessageSendingTest

# Generate test report
./gradlew testReport
```

### Test Coverage Requirements
- **Unit Tests**: 90% code coverage minimum
- **Integration Tests**: 80% user workflow coverage
- **E2E Tests**: 100% critical path coverage
- **Performance Tests**: All scenarios under performance budgets

## Maintenance Guidelines

### Test ID Stability
- Never change test IDs without updating all dependent tests
- Use semantic versioning for test suite releases
- Document any breaking changes in test ID structure

### Test Data Management
- Use consistent test data across all scenarios
- Clean up test data after each test run
- Avoid dependencies on external services in tests

### Regular Updates
- Review and update test scenarios quarterly
- Add new scenarios for feature additions
- Remove obsolete tests for deprecated features
- Monitor test execution times and optimize slow tests

## Tools and Frameworks

### Recommended Testing Stack
- **Compose Testing**: `androidx.compose.ui.test`
- **Test Runner**: `androidx.test.runner.AndroidJUnitRunner`
- **Assertions**: `androidx.test.ext.truth`
- **Mocking**: `MockK` for Kotlin
- **Performance**: `androidx.benchmark`

### CI/CD Integration
- **GitHub Actions**: Automated test execution on PR
- **Firebase Test Lab**: Device matrix testing
- **Codecov**: Test coverage reporting
- **Allure**: Rich test reporting and history

## Best Practices Summary

1. **Consistent Naming**: Always follow the component-element-action pattern
2. **Descriptive IDs**: Make test IDs self-documenting
3. **Stable References**: Use IDs that won't change with UI updates
4. **Comprehensive Coverage**: Test all user-facing functionality
5. **Performance Awareness**: Include timing expectations in tests
6. **Accessibility First**: Ensure tests cover accessibility requirements
7. **Error Resilience**: Test failure scenarios extensively
8. **Maintainable Tests**: Write tests that are easy to update and debug

This guide serves as the definitive reference for implementing and maintaining test automation in the Tchat chat system. Regular updates and team reviews ensure the testing strategy remains effective and comprehensive.