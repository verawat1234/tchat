// T069 - E2E tests for message types
/**
 * Message Types E2E Tests
 * Comprehensive Playwright tests for all 13 message components
 * Tests user interactions, accessibility, performance, and cross-browser compatibility
 */

import { test, expect, Page } from '@playwright/test';

// Test data fixtures
const testData = {
  replyMessage: {
    id: 'msg-reply-123',
    type: 'reply',
    content: {
      originalMessageId: 'msg-456',
      replyText: 'This is a comprehensive reply to test thread functionality',
      threadDepth: 2,
      isThreadStart: false,
    },
  },
  quizMessage: {
    id: 'msg-quiz-123',
    type: 'quiz',
    content: {
      title: 'React Knowledge Quiz',
      description: 'Test your React expertise with these questions',
      questions: [
        {
          id: 'q1',
          type: 'multiple_choice',
          question: 'What is the Virtual DOM?',
          options: [
            'A real DOM element',
            'A JavaScript representation of the DOM',
            'A CSS framework',
            'A browser API',
          ],
          correctAnswer: 1,
          explanation: 'The Virtual DOM is a programming concept where a virtual representation of UI is kept in memory.',
        },
        {
          id: 'q2',
          type: 'true_false',
          question: 'React components can only return JSX elements.',
          correctAnswer: false,
          explanation: 'React components can return JSX, strings, numbers, arrays, and more.',
        },
      ],
      timeLimit: 300,
      showResults: true,
      allowRetake: false,
    },
  },
  eventMessage: {
    id: 'msg-event-123',
    type: 'event',
    content: {
      title: 'Team Building Workshop',
      description: 'Join us for an interactive team building session',
      startDate: '2024-02-15T10:00:00Z',
      endDate: '2024-02-15T17:00:00Z',
      location: 'Conference Room A',
      maxAttendees: 20,
      attendees: [
        { id: 'user1', name: 'Alice Johnson', status: 'attending' },
        { id: 'user2', name: 'Bob Smith', status: 'maybe' },
      ],
      rsvpRequired: true,
      reminderSet: false,
    },
  },
  productMessage: {
    id: 'msg-product-123',
    type: 'product',
    content: {
      name: 'Premium Wireless Headphones',
      description: 'High-quality noise-canceling headphones',
      price: 199.99,
      currency: 'USD',
      images: ['/api/images/headphones-1.jpg', '/api/images/headphones-2.jpg'],
      inStock: true,
      variants: [
        { id: 'color-black', name: 'Black', type: 'color', value: '#000000' },
        { id: 'color-white', name: 'White', type: 'color', value: '#FFFFFF' },
      ],
      rating: 4.5,
      reviewCount: 128,
    },
  },
};

// Helper functions
async function navigateToChat(page: Page) {
  await page.goto('/chat/test-room');
  await page.waitForSelector('[data-testid="chat-container"]');
}

async function createMessage(page: Page, messageData: any) {
  // Mock API response for message creation
  await page.route('**/api/v1/messages', async (route) => {
    await route.fulfill({
      status: 201,
      contentType: 'application/json',
      body: JSON.stringify({
        message: {
          ...messageData,
          senderId: 'test-user',
          senderName: 'Test User',
          timestamp: new Date().toISOString(),
          isOwn: true,
        },
      }),
    });
  });

  // Simulate message creation through UI
  await page.getByTestId('message-composer').click();
  await page.getByTestId(`message-type-${messageData.type}`).click();

  return messageData.id;
}

async function waitForMessageToLoad(page: Page, messageId: string) {
  await page.waitForSelector(`[data-testid="message-${messageId}"]`);
}

test.describe('Message Types E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    await navigateToChat(page);
  });

  test.describe('Reply Message Component', () => {
    test('should display reply message with thread visualization', async ({ page }) => {
      // Arrange
      const messageId = await createMessage(page, testData.replyMessage);
      await waitForMessageToLoad(page, messageId);

      // Assert
      const replyMessage = page.getByTestId(`message-${messageId}`);
      await expect(replyMessage).toBeVisible();

      // Check reply text
      await expect(replyMessage.getByText(testData.replyMessage.content.replyText)).toBeVisible();

      // Check thread line visualization
      await expect(replyMessage.locator('.thread-line')).toBeVisible();

      // Check thread depth indication
      const threadIndicator = replyMessage.locator(`[data-depth="${testData.replyMessage.content.threadDepth}"]`);
      await expect(threadIndicator).toBeVisible();
    });

    test('should support nested replies', async ({ page }) => {
      // Create initial reply
      const messageId = await createMessage(page, testData.replyMessage);
      await waitForMessageToLoad(page, messageId);

      // Reply to the reply
      const replyButton = page.getByTestId(`message-${messageId}`).getByTestId('reply-button');
      await replyButton.click();

      const nestedReplyInput = page.getByTestId('nested-reply-input');
      await expect(nestedReplyInput).toBeVisible();
      await nestedReplyInput.fill('This is a nested reply');

      await page.getByTestId('send-nested-reply').click();

      // Verify nested reply appears with increased thread depth
      await page.waitForSelector('[data-depth="3"]');
      const nestedReply = page.locator('[data-depth="3"]').first();
      await expect(nestedReply.getByText('This is a nested reply')).toBeVisible();
    });

    test('should handle thread navigation', async ({ page }) => {
      const messageId = await createMessage(page, testData.replyMessage);
      await waitForMessageToLoad(page, messageId);

      // Click on original message link
      const originalMessageLink = page.getByTestId(`message-${messageId}`).getByTestId('original-message-link');
      await originalMessageLink.click();

      // Should scroll to and highlight original message
      await page.waitForSelector('[data-testid="message-msg-456"][data-highlighted="true"]');
      const originalMessage = page.getByTestId('message-msg-456');
      await expect(originalMessage).toHaveAttribute('data-highlighted', 'true');
    });
  });

  test.describe('Quiz Message Component', () => {
    test('should display quiz with all questions and options', async ({ page }) => {
      const messageId = await createMessage(page, testData.quizMessage);
      await waitForMessageToLoad(page, messageId);

      const quizMessage = page.getByTestId(`message-${messageId}`);

      // Check quiz title and description
      await expect(quizMessage.getByText(testData.quizMessage.content.title)).toBeVisible();
      await expect(quizMessage.getByText(testData.quizMessage.content.description)).toBeVisible();

      // Check first question
      const question1 = quizMessage.getByTestId('question-q1');
      await expect(question1).toBeVisible();
      await expect(question1.getByText('What is the Virtual DOM?')).toBeVisible();

      // Check all options for first question
      for (let i = 0; i < 4; i++) {
        const option = question1.getByTestId(`option-${i}`);
        await expect(option).toBeVisible();
      }

      // Check time limit display
      await expect(quizMessage.getByText('Time Limit: 5:00')).toBeVisible();
    });

    test('should handle quiz interaction and scoring', async ({ page }) => {
      // Mock quiz submission API
      await page.route('**/api/v1/messages/*/interactions', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            result: {
              score: 50,
              correctAnswers: 1,
              totalQuestions: 2,
              answers: [
                { questionId: 'q1', answer: 1, isCorrect: true, timeSpent: 30 },
                { questionId: 'q2', answer: true, isCorrect: false, timeSpent: 45 },
              ],
              completedAt: new Date().toISOString(),
            },
          }),
        });
      });

      const messageId = await createMessage(page, testData.quizMessage);
      await waitForMessageToLoad(page, messageId);

      const quizMessage = page.getByTestId(`message-${messageId}`);

      // Start quiz
      await quizMessage.getByTestId('start-quiz-button').click();

      // Answer first question (correct)
      const question1 = quizMessage.getByTestId('question-q1');
      await question1.getByTestId('option-1').click();

      // Answer second question (incorrect)
      const question2 = quizMessage.getByTestId('question-q2');
      await question2.getByTestId('option-true').click();

      // Submit quiz
      await quizMessage.getByTestId('submit-quiz-button').click();

      // Wait for results
      await page.waitForSelector('[data-testid="quiz-results"]');
      const results = quizMessage.getByTestId('quiz-results');

      // Verify results display
      await expect(results.getByText('Score: 50%')).toBeVisible();
      await expect(results.getByText('1 out of 2 correct')).toBeVisible();

      // Check explanations are shown
      await expect(results.getByText('The Virtual DOM is a programming concept')).toBeVisible();
    });

    test('should enforce time limits', async ({ page }) => {
      // Create quiz with short time limit
      const shortQuizData = {
        ...testData.quizMessage,
        content: { ...testData.quizMessage.content, timeLimit: 5 }, // 5 seconds
      };

      const messageId = await createMessage(page, shortQuizData);
      await waitForMessageToLoad(page, messageId);

      const quizMessage = page.getByTestId(`message-${messageId}`);
      await quizMessage.getByTestId('start-quiz-button').click();

      // Wait for timer to run out
      await page.waitForSelector('[data-testid="quiz-timeout"]', { timeout: 7000 });

      // Verify timeout message
      await expect(quizMessage.getByText('Time\'s up!')).toBeVisible();

      // Verify submit button is disabled
      const submitButton = quizMessage.getByTestId('submit-quiz-button');
      await expect(submitButton).toBeDisabled();
    });
  });

  test.describe('Event Message Component', () => {
    test('should display event details and RSVP options', async ({ page }) => {
      const messageId = await createMessage(page, testData.eventMessage);
      await waitForMessageToLoad(page, messageId);

      const eventMessage = page.getByTestId(`message-${messageId}`);

      // Check event details
      await expect(eventMessage.getByText(testData.eventMessage.content.title)).toBeVisible();
      await expect(eventMessage.getByText(testData.eventMessage.content.description)).toBeVisible();
      await expect(eventMessage.getByText('Conference Room A')).toBeVisible();

      // Check RSVP options
      await expect(eventMessage.getByTestId('rsvp-attending')).toBeVisible();
      await expect(eventMessage.getByTestId('rsvp-maybe')).toBeVisible();
      await expect(eventMessage.getByTestId('rsvp-not-attending')).toBeVisible();

      // Check attendee list
      await expect(eventMessage.getByText('Alice Johnson')).toBeVisible();
      await expect(eventMessage.getByText('Bob Smith')).toBeVisible();

      // Check availability counter
      await expect(eventMessage.getByText('18 spots remaining')).toBeVisible();
    });

    test('should handle RSVP interactions', async ({ page }) => {
      // Mock RSVP API
      await page.route('**/api/v1/messages/*/interactions', async (route) => {
        const request = await route.request().postDataJSON();
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            message: {
              ...testData.eventMessage,
              content: {
                ...testData.eventMessage.content,
                userRsvpStatus: request.data.status,
                attendees: [
                  ...testData.eventMessage.content.attendees,
                  { id: 'current-user', name: 'Current User', status: request.data.status },
                ],
              },
            },
          }),
        });
      });

      const messageId = await createMessage(page, testData.eventMessage);
      await waitForMessageToLoad(page, messageId);

      const eventMessage = page.getByTestId(`message-${messageId}`);

      // Click attending RSVP
      await eventMessage.getByTestId('rsvp-attending').click();

      // Verify RSVP status update
      await page.waitForSelector('[data-rsvp-status="attending"]');
      await expect(eventMessage.getByTestId('current-rsvp-status')).toHaveText('You are attending');

      // Verify updated attendee count
      await expect(eventMessage.getByText('17 spots remaining')).toBeVisible();
    });

    test('should integrate with calendar', async ({ page }) => {
      const messageId = await createMessage(page, testData.eventMessage);
      await waitForMessageToLoad(page, messageId);

      const eventMessage = page.getByTestId(`message-${messageId}`);

      // Test add to calendar functionality
      await eventMessage.getByTestId('add-to-calendar').click();

      // Should trigger calendar integration
      await expect(page.getByTestId('calendar-integration-modal')).toBeVisible();

      // Test reminder setting
      await eventMessage.getByTestId('set-reminder').click();
      await expect(eventMessage.getByTestId('reminder-set')).toBeVisible();
    });
  });

  test.describe('Product Message Component', () => {
    test('should display product information and purchase options', async ({ page }) => {
      const messageId = await createMessage(page, testData.productMessage);
      await waitForMessageToLoad(page, messageId);

      const productMessage = page.getByTestId(`message-${messageId}`);

      // Check product details
      await expect(productMessage.getByText(testData.productMessage.content.name)).toBeVisible();
      await expect(productMessage.getByText(testData.productMessage.content.description)).toBeVisible();
      await expect(productMessage.getByText('$199.99')).toBeVisible();

      // Check product images
      await expect(productMessage.getByTestId('product-image-carousel')).toBeVisible();

      // Check variants
      await expect(productMessage.getByTestId('variant-color-black')).toBeVisible();
      await expect(productMessage.getByTestId('variant-color-white')).toBeVisible();

      // Check ratings
      await expect(productMessage.getByText('4.5')).toBeVisible();
      await expect(productMessage.getByText('128 reviews')).toBeVisible();

      // Check stock status
      await expect(productMessage.getByText('In Stock')).toBeVisible();
    });

    test('should handle variant selection and cart operations', async ({ page }) => {
      // Mock cart API
      await page.route('**/api/v1/messages/*/interactions', async (route) => {
        const request = await route.request().postDataJSON();
        if (request.interactionType === 'cart_add') {
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({
              success: true,
              result: {
                cartItemId: 'cart-item-123',
                quantity: request.data.quantity,
                selectedVariants: request.data.variants,
              },
            }),
          });
        }
      });

      const messageId = await createMessage(page, testData.productMessage);
      await waitForMessageToLoad(page, messageId);

      const productMessage = page.getByTestId(`message-${messageId}`);

      // Select variant
      await productMessage.getByTestId('variant-color-black').click();

      // Verify variant selection
      await expect(productMessage.getByTestId('variant-color-black')).toHaveAttribute('data-selected', 'true');

      // Add to cart
      await productMessage.getByTestId('add-to-cart-button').click();

      // Verify cart addition
      await page.waitForSelector('[data-testid="cart-success-message"]');
      await expect(productMessage.getByTestId('cart-success-message')).toBeVisible();

      // Check quantity controls
      await productMessage.getByTestId('quantity-increase').click();
      await expect(productMessage.getByTestId('quantity-input')).toHaveValue('2');
    });
  });

  test.describe('Performance Tests', () => {
    test('should load message components within performance budgets', async ({ page }) => {
      // Set up performance monitoring
      await page.addInitScript(() => {
        window.performanceMetrics = {
          componentLoadTimes: {},
          renderTimes: {},
        };

        // Override React render method to measure render time
        const originalCreateElement = React.createElement;
        React.createElement = function(...args) {
          const start = performance.now();
          const element = originalCreateElement.apply(this, args);
          const end = performance.now();

          if (args[0] && args[0].displayName && args[0].displayName.includes('Message')) {
            const componentName = args[0].displayName;
            window.performanceMetrics.renderTimes[componentName] =
              (window.performanceMetrics.renderTimes[componentName] || []);
            window.performanceMetrics.renderTimes[componentName].push(end - start);
          }

          return element;
        };
      });

      // Create multiple message types
      const messageIds = [];
      for (const messageData of Object.values(testData)) {
        const messageId = await createMessage(page, messageData);
        messageIds.push(messageId);
      }

      // Wait for all messages to load
      for (const messageId of messageIds) {
        await waitForMessageToLoad(page, messageId);
      }

      // Check performance metrics
      const metrics = await page.evaluate(() => window.performanceMetrics);

      // Verify render times are within budget (< 16ms for 60fps)
      Object.entries(metrics.renderTimes).forEach(([componentName, times]) => {
        const averageTime = times.reduce((a, b) => a + b, 0) / times.length;
        expect(averageTime).toBeLessThan(16);
        console.log(`${componentName} average render time: ${averageTime.toFixed(2)}ms`);
      });
    });

    test('should handle large datasets efficiently', async ({ page }) => {
      // Create quiz with many questions
      const largeQuizData = {
        ...testData.quizMessage,
        content: {
          ...testData.quizMessage.content,
          questions: Array.from({ length: 50 }, (_, i) => ({
            id: `q${i + 1}`,
            type: 'multiple_choice',
            question: `Question ${i + 1}: What is ${i + 1} + 1?`,
            options: [`${i}`, `${i + 1}`, `${i + 2}`, `${i + 3}`],
            correctAnswer: 2,
          })),
        },
      };

      const startTime = performance.now();
      const messageId = await createMessage(page, largeQuizData);
      await waitForMessageToLoad(page, messageId);
      const loadTime = performance.now() - startTime;

      // Should load within 2 seconds even with large dataset
      expect(loadTime).toBeLessThan(2000);

      // Test virtual scrolling or pagination
      const quizMessage = page.getByTestId(`message-${messageId}`);
      await quizMessage.getByTestId('start-quiz-button').click();

      // Should only render visible questions initially
      const visibleQuestions = await quizMessage.locator('[data-testid^="question-"]').count();
      expect(visibleQuestions).toBeLessThanOrEqual(10); // Assuming pagination of 10 questions
    });
  });

  test.describe('Accessibility Tests', () => {
    test('should meet WCAG 2.1 AA standards', async ({ page }) => {
      // Install axe-core
      await page.addScriptTag({ url: 'https://unpkg.com/axe-core@4.7.0/axe.min.js' });

      // Create messages of each type
      for (const [messageType, messageData] of Object.entries(testData)) {
        const messageId = await createMessage(page, messageData);
        await waitForMessageToLoad(page, messageId);

        // Run accessibility tests
        const results = await page.evaluate(async () => {
          const results = await axe.run();
          return results.violations;
        });

        // Should have no accessibility violations
        expect(results).toHaveLength(0);

        if (results.length > 0) {
          console.log(`Accessibility violations in ${messageType}:`, results);
        }
      }
    });

    test('should support keyboard navigation', async ({ page }) => {
      const messageId = await createMessage(page, testData.quizMessage);
      await waitForMessageToLoad(page, messageId);

      const quizMessage = page.getByTestId(`message-${messageId}`);

      // Test keyboard navigation
      await quizMessage.focus();
      await page.keyboard.press('Tab');

      // Should focus on start quiz button
      await expect(quizMessage.getByTestId('start-quiz-button')).toBeFocused();

      await page.keyboard.press('Enter');

      // Should start the quiz
      await expect(quizMessage.getByTestId('question-q1')).toBeVisible();

      // Navigate through options with arrow keys
      await page.keyboard.press('Tab');
      await page.keyboard.press('ArrowDown');
      await page.keyboard.press('ArrowDown');

      // Select with Enter
      await page.keyboard.press('Enter');

      // Verify option is selected
      await expect(quizMessage.getByTestId('option-2')).toHaveAttribute('aria-checked', 'true');
    });

    test('should provide proper screen reader support', async ({ page }) => {
      const messageId = await createMessage(page, testData.replyMessage);
      await waitForMessageToLoad(page, messageId);

      const replyMessage = page.getByTestId(`message-${messageId}`);

      // Check ARIA labels and descriptions
      await expect(replyMessage).toHaveAttribute('role', 'article');
      await expect(replyMessage).toHaveAttribute('aria-label');

      const replyText = replyMessage.getByText(testData.replyMessage.content.replyText);
      await expect(replyText).toHaveAttribute('aria-describedby');

      // Check thread navigation accessibility
      const originalMessageLink = replyMessage.getByTestId('original-message-link');
      await expect(originalMessageLink).toHaveAttribute('aria-label');
      await expect(originalMessageLink).toHaveAttribute('role', 'button');
    });
  });

  test.describe('Cross-Browser Compatibility', () => {
    ['chromium', 'firefox', 'webkit'].forEach(browserName => {
      test(`should work correctly in ${browserName}`, async ({ page, browserName: currentBrowser }) => {
        test.skip(currentBrowser !== browserName, `Skipping ${browserName} test in ${currentBrowser}`);

        // Test basic functionality in each browser
        const messageId = await createMessage(page, testData.quizMessage);
        await waitForMessageToLoad(page, messageId);

        const quizMessage = page.getByTestId(`message-${messageId}`);

        // Test quiz interaction
        await quizMessage.getByTestId('start-quiz-button').click();
        await quizMessage.getByTestId('question-q1').getByTestId('option-1').click();
        await quizMessage.getByTestId('submit-quiz-button').click();

        // Verify results are displayed
        await page.waitForSelector('[data-testid="quiz-results"]');
        await expect(quizMessage.getByTestId('quiz-results')).toBeVisible();
      });
    });
  });

  test.describe('Error Handling', () => {
    test('should gracefully handle component loading errors', async ({ page }) => {
      // Simulate component loading failure
      await page.route('**/MessageTypes/QuizMessage.js', async (route) => {
        await route.abort();
      });

      const messageId = await createMessage(page, testData.quizMessage);

      // Should show fallback UI
      await page.waitForSelector('[data-testid="message-fallback"]');
      const fallbackMessage = page.getByTestId('message-fallback');
      await expect(fallbackMessage).toBeVisible();
      await expect(fallbackMessage).toContainText('Unable to load message');
    });

    test('should handle API errors gracefully', async ({ page }) => {
      // Mock API error
      await page.route('**/api/v1/messages/*/interactions', async (route) => {
        await route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Internal server error' }),
        });
      });

      const messageId = await createMessage(page, testData.quizMessage);
      await waitForMessageToLoad(page, messageId);

      const quizMessage = page.getByTestId(`message-${messageId}`);
      await quizMessage.getByTestId('start-quiz-button').click();
      await quizMessage.getByTestId('question-q1').getByTestId('option-1').click();
      await quizMessage.getByTestId('submit-quiz-button').click();

      // Should show error message
      await page.waitForSelector('[data-testid="quiz-error"]');
      await expect(quizMessage.getByTestId('quiz-error')).toBeVisible();
      await expect(quizMessage.getByTestId('quiz-error')).toContainText('Unable to submit quiz');
    });
  });
});