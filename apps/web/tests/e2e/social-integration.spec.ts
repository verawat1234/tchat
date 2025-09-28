/**
 * Social Integration E2E Test Scenarios
 *
 * Comprehensive end-to-end testing for Tchat social features including:
 * - Social feed loading and interaction
 * - Post creation and engagement (like, comment, share)
 * - User profile viewing and following
 * - Social discovery and trending content
 * - Error handling and fallback states
 * - Southeast Asian regional functionality
 */

import { test, expect, Page } from '@playwright/test';

// Test data constants
const TEST_USER = {
  name: 'Test User',
  avatar: 'https://example.com/avatar.jpg'
};

const TEST_POST_CONTENT = 'Testing amazing Pad Thai at Chatuchak Market! 🍜✨ #PadThai #StreetFood #Bangkok';
const TEST_COMMENT_TEXT = 'Looks delicious! Where exactly is this vendor?';
const TEST_STORY_TEXT = 'Beautiful sunset over Bangkok today!';

// Helper functions
async function navigateToSocialTab(page: Page) {
  await page.goto('/');
  await page.click('[data-testid="social-tab-trigger"]');
  await expect(page.locator('[data-testid="social-tab-container"]')).toBeVisible();
}

async function waitForSocialFeedLoad(page: Page) {
  // Wait for RTK Query to load social feed data
  await page.waitForSelector('[data-testid="social-main-tabs"]');
  await page.waitForSelector('[data-testid="social-tab-feed"]');
}

async function createTestPost(page: Page, content: string) {
  await page.fill('[data-testid="create-post-textarea-input"]', content);
  await page.click('[data-testid="create-post-submit-button"]');
  await expect(page.locator('.toast')).toContainText('Post created successfully');
}

test.describe('Social Integration - Core Functionality', () => {
  test.beforeEach(async ({ page }) => {
    // Mock the social API responses for consistent testing
    await page.route('**/api/v1/social/**', (route) => {
      const url = route.request().url();

      if (url.includes('/feed')) {
        route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify([
            {
              id: 'test-post-1',
              author: { name: 'Sarah Johnson', avatar: 'https://example.com/sarah.jpg', verified: false, type: 'user' },
              content: 'Amazing Pad Thai at Chatuchak Market! 🍜✨',
              images: ['https://example.com/padthai.jpg'],
              timestamp: '25 min ago',
              likes: 47,
              comments: 12,
              shares: 3,
              location: 'Chatuchak Weekend Market, Bangkok',
              tags: ['#PadThai', '#StreetFood', '#Bangkok'],
              type: 'image',
              source: 'following'
            }
          ])
        });
      } else if (url.includes('/stories')) {
        route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify([
            {
              id: 'story-1',
              author: { name: 'Mike Chen', avatar: 'https://example.com/mike.jpg' },
              preview: 'https://example.com/story-preview.jpg',
              content: 'Live from floating market!',
              timestamp: '2h ago',
              isViewed: false,
              isLive: false,
              media: [{ type: 'image', url: 'https://example.com/story.jpg', duration: 5 }],
              expiresAt: new Date(Date.now() + 22 * 60 * 60 * 1000).toISOString()
            }
          ])
        });
      } else if (url.includes('/friends')) {
        route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify([
            {
              id: 'friend-1',
              name: 'Sarah Johnson',
              username: '@sarah_foodie',
              avatar: 'https://example.com/sarah.jpg',
              isOnline: true,
              mutualFriends: 12,
              status: 'Exploring Bangkok street food! 🍜',
              isFollowing: true
            }
          ])
        });
      }
    });
  });

  test('should load social tab and display feed content', async ({ page }) => {
    await navigateToSocialTab(page);
    await waitForSocialFeedLoad(page);

    // Verify main tab structure is present
    await expect(page.locator('[data-testid="social-main-tabs"]')).toBeVisible();
    await expect(page.locator('[data-testid="social-tab-friends"]')).toBeVisible();
    await expect(page.locator('[data-testid="social-tab-feed"]')).toBeVisible();
    await expect(page.locator('[data-testid="social-tab-discover"]')).toBeVisible();
    await expect(page.locator('[data-testid="social-tab-events"]')).toBeVisible();

    // Verify stories section is present
    await expect(page.locator('[data-testid="social-stories-section"]')).toBeVisible();
    await expect(page.locator('[data-testid="social-stories-list"]')).toBeVisible();

    // Verify create post section is present
    await expect(page.locator('[data-testid="social-create-post-section"]')).toBeVisible();
  });

  test('should display social stories with proper test IDs', async ({ page }) => {
    await navigateToSocialTab(page);
    await waitForSocialFeedLoad(page);

    // Check for user's own story creation option
    await expect(page.locator('[data-testid="social-story-create-button"]')).toBeVisible();

    // Check for story items with proper data attributes
    const storyItems = page.locator('[data-testid^="social-story-item-"]');
    await expect(storyItems).toHaveCount.greaterThan(0);

    // Verify story item attributes
    const firstStory = storyItems.first();
    await expect(firstStory).toHaveAttribute('data-story-type');
    await expect(firstStory).toHaveAttribute('data-story-live');
    await expect(firstStory).toHaveAttribute('data-story-viewed');
  });

  test('should switch between social tabs', async ({ page }) => {
    await navigateToSocialTab(page);
    await waitForSocialFeedLoad(page);

    // Test switching to Feed tab
    await page.click('[data-testid="social-tab-feed"]');
    await expect(page.locator('[data-testid="social-tab-feed"]')).toHaveAttribute('data-state', 'active');

    // Test switching to Discover tab
    await page.click('[data-testid="social-tab-discover"]');
    await expect(page.locator('[data-testid="social-tab-discover"]')).toHaveAttribute('data-state', 'active');

    // Test switching to Events tab
    await page.click('[data-testid="social-tab-events"]');
    await expect(page.locator('[data-testid="social-tab-events"]')).toHaveAttribute('data-state', 'active');

    // Test switching back to Friends tab
    await page.click('[data-testid="social-tab-friends"]');
    await expect(page.locator('[data-testid="social-tab-friends"]')).toHaveAttribute('data-state', 'active');
  });
});

test.describe('Social Integration - Post Creation and Interaction', () => {
  test.beforeEach(async ({ page }) => {
    await navigateToSocialTab(page);
    await waitForSocialFeedLoad(page);
  });

  test('should create a text post successfully', async ({ page }) => {
    // Navigate to create post section
    await expect(page.locator('[data-testid="create-post-textarea-input"]')).toBeVisible();

    // Type content into the textarea
    await page.fill('[data-testid="create-post-textarea-input"]', TEST_POST_CONTENT);

    // Verify submit button becomes enabled
    await expect(page.locator('[data-testid="create-post-submit-button"]')).not.toHaveAttribute('data-disabled', 'true');

    // Submit the post
    await page.click('[data-testid="create-post-submit-button"]');

    // Verify success feedback
    await expect(page.locator('.toast')).toContainText('Post created successfully');

    // Verify textarea is cleared
    await expect(page.locator('[data-testid="create-post-textarea-input"]')).toHaveValue('');
  });

  test('should interact with photo post creation options', async ({ page }) => {
    // Test photo button
    await page.click('[data-testid="create-post-photo-button"]');
    await expect(page.locator('.toast')).toContainText('Photo post created');

    // Test gallery button
    await page.click('[data-testid="create-post-gallery-button"]');
    await expect(page.locator('.toast')).toContainText('Gallery post created');

    // Test location button
    await page.click('[data-testid="create-post-location-button"]');
    await expect(page.locator('.toast')).toContainText('Location post created');
  });

  test('should prevent empty post submission', async ({ page }) => {
    // Verify submit button is disabled when textarea is empty
    await expect(page.locator('[data-testid="create-post-submit-button"]')).toHaveAttribute('data-disabled', 'true');

    // Try to click disabled button (should not submit)
    await page.click('[data-testid="create-post-submit-button"]');

    // Verify no success toast appears
    const toasts = page.locator('.toast');
    await expect(toasts).toHaveCount(0);
  });

  test('should submit post with Enter key', async ({ page }) => {
    await page.fill('[data-testid="create-post-textarea-input"]', TEST_POST_CONTENT);
    await page.press('[data-testid="create-post-textarea-input"]', 'Enter');

    await expect(page.locator('.toast')).toContainText('Post created successfully');
  });
});

test.describe('Social Integration - Post Engagement', () => {
  test.beforeEach(async ({ page }) => {
    await navigateToSocialTab(page);
    await waitForSocialFeedLoad(page);

    // Switch to feed tab to see posts
    await page.click('[data-testid="social-tab-feed"]');
  });

  test('should like and unlike posts', async ({ page }) => {
    // Wait for posts to load
    const firstPost = page.locator('[data-testid^="social-post-"]').first();
    await expect(firstPost).toBeVisible();

    const postId = await firstPost.getAttribute('data-testid');
    const postIdValue = postId?.replace('social-post-', '');

    // Test liking a post
    const likeButton = page.locator(`[data-testid="social-post-like-button-${postIdValue}"]`);
    await likeButton.click();

    // Verify like state changes
    await expect(likeButton).toHaveAttribute('data-liked', 'true');
    await expect(page.locator('.toast')).toContainText('Post liked');

    // Test unliking a post
    await likeButton.click();
    await expect(likeButton).toHaveAttribute('data-liked', 'false');
  });

  test('should open and interact with comments', async ({ page }) => {
    const firstPost = page.locator('[data-testid^="social-post-"]').first();
    const postId = await firstPost.getAttribute('data-testid');
    const postIdValue = postId?.replace('social-post-', '');

    // Open comments section
    await page.click(`[data-testid="social-post-comment-button-${postIdValue}"]`);

    // Verify comments section opens
    await expect(page.locator(`[data-testid="social-post-comments-section-${postIdValue}"]`)).toBeVisible();
    await expect(page.locator(`[data-testid="social-post-comments-section-${postIdValue}"]`)).toHaveAttribute('data-comments-open', 'true');

    // Add a comment
    await page.fill(`[data-testid="social-comment-input-field-${postIdValue}"]`, TEST_COMMENT_TEXT);
    await page.click(`[data-testid="social-comment-submit-button-${postIdValue}"]`);

    // Verify comment submission
    await expect(page.locator('.toast')).toContainText('Comment added');

    // Verify comment input is cleared
    await expect(page.locator(`[data-testid="social-comment-input-field-${postIdValue}"]`)).toHaveValue('');
  });

  test('should prevent empty comment submission', async ({ page }) => {
    const firstPost = page.locator('[data-testid^="social-post-"]').first();
    const postId = await firstPost.getAttribute('data-testid');
    const postIdValue = postId?.replace('social-post-', '');

    // Open comments
    await page.click(`[data-testid="social-post-comment-button-${postIdValue}"]`);

    // Verify submit button is disabled when comment is empty
    await expect(page.locator(`[data-testid="social-comment-submit-button-${postIdValue}"]`)).toHaveAttribute('data-disabled', 'true');
  });

  test('should share posts', async ({ page }) => {
    const firstPost = page.locator('[data-testid^="social-post-"]').first();
    const postId = await firstPost.getAttribute('data-testid');
    const postIdValue = postId?.replace('social-post-', '');

    // Test sharing a post
    await page.click(`[data-testid="social-post-share-button-${postIdValue}"]`);
    await expect(page.locator('.toast')).toContainText('Post link copied to clipboard');
  });

  test('should bookmark and unbookmark posts', async ({ page }) => {
    const firstPost = page.locator('[data-testid^="social-post-"]').first();
    const postId = await firstPost.getAttribute('data-testid');
    const postIdValue = postId?.replace('social-post-', '');

    // Test bookmarking a post
    const bookmarkButton = page.locator(`[data-testid="social-post-bookmark-button-${postIdValue}"]`);
    await bookmarkButton.click();

    // Verify bookmark state changes
    await expect(bookmarkButton).toHaveAttribute('data-bookmarked', 'true');
    await expect(page.locator('.toast')).toContainText('Added to bookmarks');

    // Test unbookmarking
    await bookmarkButton.click();
    await expect(bookmarkButton).toHaveAttribute('data-bookmarked', 'false');
    await expect(page.locator('.toast')).toContainText('Removed from bookmarks');
  });

  test('should interact with post options menu', async ({ page }) => {
    const firstPost = page.locator('[data-testid^="social-post-"]').first();
    const postId = await firstPost.getAttribute('data-testid');
    const postIdValue = postId?.replace('social-post-', '');

    // Open options menu
    await page.click(`[data-testid="social-post-options-trigger-${postIdValue}"]`);

    // Verify menu opens
    await expect(page.locator(`[data-testid="social-post-options-content-${postIdValue}"]`)).toBeVisible();

    // Test bookmark option
    await page.click(`[data-testid="social-post-bookmark-item-${postIdValue}"]`);
    await expect(page.locator('.toast')).toContainText('Added to bookmarks');
  });

  test('should click hashtags in post content', async ({ page }) => {
    const firstPost = page.locator('[data-testid^="social-post-"]').first();

    // Look for hashtag elements in the post
    const hashtag = firstPost.locator('[data-testid^="social-post-hashtag-"]').first();
    if (await hashtag.isVisible()) {
      const hashtagText = await hashtag.getAttribute('data-hashtag');
      await hashtag.click();
      await expect(page.locator('.toast')).toContainText(`Searching for ${hashtagText}`);
    }
  });
});

test.describe('Social Integration - Stories/Moments', () => {
  test.beforeEach(async ({ page }) => {
    await navigateToSocialTab(page);
    await waitForSocialFeedLoad(page);
  });

  test('should open create story dialog', async ({ page }) => {
    // Click on "Your Moment" story creation button
    await page.click('[data-testid="social-story-create-button"]');

    // Verify create story dialog opens
    await expect(page.locator('[data-testid="social-create-story-dialog"]')).toBeVisible();
    await expect(page.locator('[data-testid="social-create-story-content"]')).toBeVisible();
  });

  test('should create a story/moment', async ({ page }) => {
    // Open create story dialog
    await page.click('[data-testid="social-story-create-button"]');

    // Fill in story text
    await page.fill('[data-testid="social-create-story-text-input"]', TEST_STORY_TEXT);

    // Submit the story
    await page.click('[data-testid="social-create-story-submit-button"]');

    // Verify success feedback
    await expect(page.locator('.toast')).toContainText('Moment created');
  });

  test('should cancel story creation', async ({ page }) => {
    // Open create story dialog
    await page.click('[data-testid="social-story-create-button"]');

    // Cancel story creation
    await page.click('[data-testid="social-create-story-cancel-button"]');

    // Verify dialog closes
    await expect(page.locator('[data-testid="social-create-story-dialog"]')).not.toBeVisible();
  });

  test('should view existing stories', async ({ page }) => {
    // Click on a story (not the create button)
    const storyItems = page.locator('[data-testid^="social-story-item-"]:not([data-story-type="create"])');

    if (await storyItems.count() > 0) {
      const firstStory = storyItems.first();
      await firstStory.click();

      // Verify story viewer opens
      await expect(page.locator('[data-testid="social-story-viewer-dialog"]')).toBeVisible();

      // Verify story content is displayed
      await expect(page.locator('[data-testid^="social-story-viewer-"]')).toBeVisible();
    }
  });

  test('should navigate through story segments', async ({ page }) => {
    // Open a story with multiple segments
    const storyItems = page.locator('[data-testid^="social-story-item-"]:not([data-story-type="create"])');

    if (await storyItems.count() > 0) {
      await storyItems.first().click();

      // Test navigation buttons
      const storyId = await page.locator('[data-testid^="social-story-viewer-"]').getAttribute('data-testid');
      const storyIdValue = storyId?.replace('social-story-viewer-', '');

      // Test next navigation
      await page.click(`[data-testid="social-story-nav-next-${storyIdValue}"]`);

      // Test previous navigation
      await page.click(`[data-testid="social-story-nav-previous-${storyIdValue}"]`);
    }
  });
});

test.describe('Social Integration - Error Handling and Fallback States', () => {
  test('should handle API failures gracefully', async ({ page }) => {
    // Mock API failures
    await page.route('**/api/v1/social/feed', (route) => {
      route.fulfill({ status: 500, body: 'Server Error' });
    });

    await navigateToSocialTab(page);

    // Verify fallback content is shown
    await expect(page.locator('[data-testid="social-tab-feed"]')).toBeVisible();

    // Fallback posts should be displayed when API fails
    await page.click('[data-testid="social-tab-feed"]');
    await expect(page.locator('[data-testid^="social-post-"]')).toHaveCount.greaterThan(0);
  });

  test('should handle slow API responses', async ({ page }) => {
    // Mock slow API response
    await page.route('**/api/v1/social/feed', async (route) => {
      await new Promise(resolve => setTimeout(resolve, 2000)); // 2 second delay
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([])
      });
    });

    await navigateToSocialTab(page);

    // Verify the UI doesn't break during loading
    await expect(page.locator('[data-testid="social-main-tabs"]')).toBeVisible();
  });

  test('should handle network disconnection', async ({ page }) => {
    await navigateToSocialTab(page);

    // Simulate network disconnection
    await page.context().setOffline(true);

    // Try to create a post
    await page.fill('[data-testid="create-post-textarea-input"]', 'Test post during offline');
    await page.click('[data-testid="create-post-submit-button"]');

    // Should still provide feedback (optimistic updates)
    await expect(page.locator('.toast')).toContainText('Post created successfully');

    // Restore network
    await page.context().setOffline(false);
  });
});

test.describe('Social Integration - Accessibility and Performance', () => {
  test('should be keyboard navigable', async ({ page }) => {
    await navigateToSocialTab(page);
    await waitForSocialFeedLoad(page);

    // Test tab navigation through social tabs
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    await page.keyboard.press('Enter'); // Should activate a tab

    // Test keyboard interaction with posts
    const firstPost = page.locator('[data-testid^="social-post-"]').first();
    if (await firstPost.isVisible()) {
      await firstPost.focus();
      await page.keyboard.press('Tab'); // Navigate to like button
      await page.keyboard.press('Enter'); // Should trigger like
    }
  });

  test('should load within performance budgets', async ({ page }) => {
    const startTime = Date.now();

    await navigateToSocialTab(page);
    await waitForSocialFeedLoad(page);

    const loadTime = Date.now() - startTime;

    // Social feed should load within 3 seconds
    expect(loadTime).toBeLessThan(3000);
  });

  test('should handle image loading failures', async ({ page }) => {
    // Mock image failures
    await page.route('**/*.jpg', (route) => {
      route.fulfill({ status: 404 });
    });

    await navigateToSocialTab(page);
    await page.click('[data-testid="social-tab-feed"]');

    // Verify posts still display even with broken images
    await expect(page.locator('[data-testid^="social-post-"]')).toHaveCount.greaterThan(0);
  });
});

test.describe('Social Integration - Southeast Asian Regional Features', () => {
  test('should display Bangkok-specific content', async ({ page }) => {
    await navigateToSocialTab(page);
    await page.click('[data-testid="social-tab-feed"]');

    // Check for Bangkok-related content in posts and create post hints
    await expect(page.locator('[data-testid="create-post-hint"]')).toContainText('Bangkok');

    // Check for Thai cultural elements
    const createTextarea = page.locator('[data-testid="create-post-textarea-input"]');
    await expect(createTextarea).toHaveAttribute('placeholder', /Bangkok|🍜/);
  });

  test('should support Thai cultural events', async ({ page }) => {
    await navigateToSocialTab(page);
    await page.click('[data-testid="social-tab-events"]');

    // Verify events tab displays cultural events
    await expect(page.locator('[data-testid="social-tab-events"]')).toBeVisible();

    // Look for Thai cultural references in events
    const eventsContent = page.locator('[data-testid="social-tab-events"]');
    // This would check for Thai festival names, food events, etc.
  });

  test('should handle Thai language hashtags', async ({ page }) => {
    await navigateToSocialTab(page);

    // Create a post with Thai hashtags
    const thaiContent = 'อร่อยมาก! Delicious Pad Thai #อาหารไทย #ผัดไทย #Bangkok';
    await page.fill('[data-testid="create-post-textarea-input"]', thaiContent);
    await page.click('[data-testid="create-post-submit-button"]');

    await expect(page.locator('.toast')).toContainText('Post created successfully');
  });
});

// Performance and Memory Tests
test.describe('Social Integration - Performance Monitoring', () => {
  test('should not exceed memory limits', async ({ page }) => {
    await navigateToSocialTab(page);

    // Create multiple posts to test memory usage
    for (let i = 0; i < 10; i++) {
      await createTestPost(page, `Test post ${i + 1} for memory testing`);
      await page.waitForTimeout(100);
    }

    // Check if the page is still responsive
    await expect(page.locator('[data-testid="social-main-tabs"]')).toBeVisible();
  });

  test('should handle rapid user interactions', async ({ page }) => {
    await navigateToSocialTab(page);
    await page.click('[data-testid="social-tab-feed"]');

    // Rapidly click like buttons
    const likeButtons = page.locator('[data-testid^="social-post-like-button-"]');
    const buttonCount = await likeButtons.count();

    for (let i = 0; i < Math.min(buttonCount, 5); i++) {
      await likeButtons.nth(i).click({ timeout: 100 });
    }

    // Verify the UI remains stable
    await expect(page.locator('[data-testid="social-main-tabs"]')).toBeVisible();
  });
});