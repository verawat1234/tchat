/**
 * Social Post Creation and Engagement Workflows E2E Tests
 *
 * Comprehensive testing scenarios for:
 * - Text, photo, gallery, and location post creation
 * - Post engagement patterns (like, comment, share sequences)
 * - RTK Query integration with social API endpoints
 * - Optimistic updates and error recovery
 * - Post lifecycle management
 */

import { test, expect, Page } from '@playwright/test';

// Test constants
const POST_CREATION_SCENARIOS = {
  TEXT_POST: {
    content: 'Discovering amazing street food in Bangkok! The flavors here are incredible 🍜✨',
    tags: ['#Bangkok', '#StreetFood', '#Thailand'],
    location: null
  },
  PHOTO_POST: {
    content: 'Beautiful sunset at Wat Pho temple 🌅',
    tags: ['#WatPho', '#Sunset', '#Bangkok'],
    location: 'Wat Pho, Bangkok',
    imageUrl: 'https://images.unsplash.com/photo-1563492065-cd5bab1c2d64'
  },
  GALLERY_POST: {
    content: 'Amazing food tour through Chatuchak Market! Multiple shots of the best dishes 🍲🥘',
    tags: ['#ChatuchakMarket', '#FoodTour', '#Bangkok'],
    location: 'Chatuchak Weekend Market',
    images: [
      'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04',
      'https://images.unsplash.com/photo-1743485753872-3b24372fcd24'
    ]
  },
  LOCATION_POST: {
    content: 'Currently at this amazing floating market! The energy here is incredible 🛶',
    tags: ['#FloatingMarket', '#Thailand'],
    location: 'Damnoen Saduak Floating Market, Thailand',
    imageUrl: 'https://images.unsplash.com/photo-1513475382585-d06e58bcb0e0'
  }
};

const ENGAGEMENT_SCENARIOS = {
  LIKE_SEQUENCES: [
    { action: 'like', expectedFeedback: 'Post liked!' },
    { action: 'unlike', expectedFeedback: null }
  ],
  COMMENT_SEQUENCES: [
    { text: 'This looks amazing! Where exactly is this place?', type: 'question' },
    { text: 'I love Thai food! 😍', type: 'reaction' },
    { text: 'Thanks for sharing! Adding this to my must-visit list 📝', type: 'appreciation' }
  ],
  SHARE_SCENARIOS: [
    { method: 'copy_link', expectedFeedback: 'Post link copied to clipboard!' },
    { method: 'social_share', expectedFeedback: 'Post shared!' }
  ]
};

// Helper functions
async function navigateToSocialTab(page: Page) {
  await page.goto('/');
  await page.click('[data-testid="social-tab-trigger"]');
  await expect(page.locator('[data-testid="social-tab-container"]')).toBeVisible();
  await page.waitForLoadState('networkidle');
}

async function mockSocialAPIEndpoints(page: Page) {
  // Mock successful post creation
  await page.route('**/api/v1/social/posts', (route) => {
    if (route.request().method() === 'POST') {
      const postData = route.request().postDataJSON();
      route.fulfill({
        status: 201,
        contentType: 'application/json',
        body: JSON.stringify({
          id: `post-${Date.now()}`,
          ...postData,
          timestamp: 'just now',
          likes: 0,
          comments: 0,
          shares: 0,
          author: {
            name: 'Test User',
            avatar: 'https://example.com/avatar.jpg',
            verified: false,
            type: 'user'
          }
        })
      });
    }
  });

  // Mock like/unlike actions
  await page.route('**/api/v1/social/posts/*/like', (route) => {
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ success: true })
    });
  });

  // Mock comment creation
  await page.route('**/api/v1/social/posts/*/comments', (route) => {
    const commentData = route.request().postDataJSON();
    route.fulfill({
      status: 201,
      contentType: 'application/json',
      body: JSON.stringify({
        id: `comment-${Date.now()}`,
        ...commentData,
        timestamp: 'just now',
        likes: 0,
        isLiked: false
      })
    });
  });
}

async function fillPostContent(page: Page, content: string) {
  await page.fill('[data-testid="create-post-textarea-input"]', content);
  await expect(page.locator('[data-testid="create-post-submit-button"]')).not.toHaveAttribute('data-disabled', 'true');
}

async function submitPost(page: Page) {
  await page.click('[data-testid="create-post-submit-button"]');
  await expect(page.locator('.toast')).toContainText('Post created successfully');
}

async function verifyPostInFeed(page: Page, content: string) {
  // Switch to feed tab to see the created post
  await page.click('[data-testid="social-tab-feed"]');

  // Wait for the new post to appear
  await expect(page.locator(`[data-testid^="social-post-text-"]:has-text("${content.substring(0, 20)}")`)).toBeVisible({ timeout: 5000 });
}

test.describe('Social Post Creation Workflows', () => {
  test.beforeEach(async ({ page }) => {
    await mockSocialAPIEndpoints(page);
    await navigateToSocialTab(page);
  });

  test('should create a text post with hashtags', async ({ page }) => {
    const scenario = POST_CREATION_SCENARIOS.TEXT_POST;

    // Fill in post content
    await fillPostContent(page, scenario.content);

    // Verify hashtags are detected and styled
    const hashtagElements = page.locator('[data-testid^="social-post-hashtag-"]');
    if (await hashtagElements.count() > 0) {
      await expect(hashtagElements.first()).toHaveAttribute('data-hashtag');
    }

    // Submit the post
    await submitPost(page);

    // Verify post appears in feed with correct content
    await verifyPostInFeed(page, scenario.content);

    // Verify hashtags are clickable in the created post
    const createdPost = page.locator('[data-testid^="social-post-"]').first();
    const hashtags = createdPost.locator('[data-testid^="social-post-hashtag-"]');

    if (await hashtags.count() > 0) {
      const firstHashtag = hashtags.first();
      const hashtagText = await firstHashtag.getAttribute('data-hashtag');
      await firstHashtag.click();
      await expect(page.locator('.toast')).toContainText(`Searching for ${hashtagText}`);
    }
  });

  test('should create a photo post', async ({ page }) => {
    // Click photo creation button
    await page.click('[data-testid="create-post-photo-button"]');

    // Verify photo post creation feedback
    await expect(page.locator('.toast')).toContainText('Photo post created');

    // Verify post appears in feed
    await page.click('[data-testid="social-tab-feed"]');

    // Check for the newly created photo post
    const photoPosts = page.locator('[data-testid^="social-post-"][data-post-type="image"]');
    await expect(photoPosts).toHaveCount.greaterThan(0);

    // Verify image elements are present
    const firstPhotoPost = photoPosts.first();
    await expect(firstPhotoPost.locator('[data-testid^="social-post-main-image-"]')).toBeVisible();
  });

  test('should create a gallery post with multiple images', async ({ page }) => {
    // Click gallery creation button
    await page.click('[data-testid="create-post-gallery-button"]');

    // Verify gallery post creation feedback
    await expect(page.locator('.toast')).toContainText('Gallery post created');

    // Verify post appears in feed
    await page.click('[data-testid="social-tab-feed"]');

    // Check for the newly created gallery post
    const galleryPosts = page.locator('[data-testid^="social-post-"][data-post-type="image"]');
    await expect(galleryPosts).toHaveCount.greaterThan(0);

    // Verify gallery-specific content
    const firstGalleryPost = galleryPosts.first();
    await expect(firstGalleryPost.locator('[data-testid^="social-post-text-"]')).toContainText('Multiple shots');
  });

  test('should create a location post', async ({ page }) => {
    // Click location creation button
    await page.click('[data-testid="create-post-location-button"]');

    // Verify location post creation feedback
    await expect(page.locator('.toast')).toContainText('Location post created');

    // Verify post appears in feed with location
    await page.click('[data-testid="social-tab-feed"]');

    const locationPosts = page.locator('[data-testid^="social-post-"]');
    await expect(locationPosts).toHaveCount.greaterThan(0);

    // Verify location metadata is present
    const firstLocationPost = locationPosts.first();
    const locationElement = firstLocationPost.locator('[data-testid^="social-post-location-"]');
    if (await locationElement.isVisible()) {
      await expect(locationElement).toContainText('Market');
    }
  });

  test('should validate post content length and formatting', async ({ page }) => {
    // Test with very long content
    const longContent = 'A'.repeat(1000) + ' Bangkok street food is amazing! 🍜';
    await fillPostContent(page, longContent);
    await submitPost(page);

    // Verify long content is handled properly
    await verifyPostInFeed(page, longContent.substring(0, 50));

    // Test with special characters and emojis
    const specialContent = '🍜🌶️🥘 Testing Thai cuisine with special chars: àáâãäå ñ öøœþ 测试中文 テスト #Thai';
    await fillPostContent(page, specialContent);
    await submitPost(page);

    await verifyPostInFeed(page, specialContent.substring(0, 20));
  });

  test('should handle rapid post creation', async ({ page }) => {
    // Create multiple posts quickly
    const posts = [
      'Quick post 1 about Bangkok',
      'Quick post 2 about Thai food',
      'Quick post 3 about temples'
    ];

    for (const postContent of posts) {
      await fillPostContent(page, postContent);
      await submitPost(page);
      await page.waitForTimeout(500); // Small delay between posts
    }

    // Verify all posts appear in feed
    await page.click('[data-testid="social-tab-feed"]');

    for (const postContent of posts) {
      await expect(page.locator(`[data-testid^="social-post-text-"]:has-text("${postContent.substring(0, 15)}")`)).toBeVisible();
    }
  });

  test('should persist draft content during navigation', async ({ page }) => {
    const draftContent = 'This is a draft post about Bangkok temples';

    // Start typing a post
    await fillPostContent(page, draftContent);

    // Navigate to different tab
    await page.click('[data-testid="social-tab-discover"]');
    await page.click('[data-testid="social-tab-friends"]');

    // Return to original position and verify content persists
    const textareaValue = await page.locator('[data-testid="create-post-textarea-input"]').inputValue();
    expect(textareaValue).toBe(draftContent);
  });
});

test.describe('Social Post Engagement Workflows', () => {
  test.beforeEach(async ({ page }) => {
    await mockSocialAPIEndpoints(page);
    await navigateToSocialTab(page);

    // Create a test post to interact with
    await fillPostContent(page, 'Test post for engagement testing 🍜 #Bangkok');
    await submitPost(page);
    await page.click('[data-testid="social-tab-feed"]');
  });

  test('should perform complete like workflow', async ({ page }) => {
    const firstPost = page.locator('[data-testid^="social-post-"]').first();
    const postId = await firstPost.getAttribute('data-testid');
    const postIdValue = postId?.replace('social-post-', '');

    const likeButton = page.locator(`[data-testid="social-post-like-button-${postIdValue}"]`);
    const likeCount = page.locator(`[data-testid="social-post-like-count-${postIdValue}"]`);

    // Get initial like count
    const initialCount = await likeCount.textContent();
    const initialCountNum = parseInt(initialCount || '0');

    // Test liking
    await likeButton.click();
    await expect(likeButton).toHaveAttribute('data-liked', 'true');
    await expect(page.locator('.toast')).toContainText('Post liked');

    // Verify like count increased
    await expect(likeCount).toContainText((initialCountNum + 1).toString());

    // Test unliking
    await likeButton.click();
    await expect(likeButton).toHaveAttribute('data-liked', 'false');

    // Verify like count decreased
    await expect(likeCount).toContainText(initialCountNum.toString());
  });

  test('should perform complete comment workflow', async ({ page }) => {
    const firstPost = page.locator('[data-testid^="social-post-"]').first();
    const postId = await firstPost.getAttribute('data-testid');
    const postIdValue = postId?.replace('social-post-', '');

    // Open comments section
    await page.click(`[data-testid="social-post-comment-button-${postIdValue}"]`);
    await expect(page.locator(`[data-testid="social-post-comments-section-${postIdValue}"]`)).toBeVisible();

    // Test each comment type
    for (const commentScenario of ENGAGEMENT_SCENARIOS.COMMENT_SEQUENCES) {
      // Add comment
      await page.fill(`[data-testid="social-comment-input-field-${postIdValue}"]`, commentScenario.text);
      await page.click(`[data-testid="social-comment-submit-button-${postIdValue}"]`);

      // Verify comment added
      await expect(page.locator('.toast')).toContainText('Comment added');
      await expect(page.locator(`[data-testid="social-comment-input-field-${postIdValue}"]`)).toHaveValue('');

      // Verify comment appears in list
      await expect(page.locator(`[data-testid^="social-comment-text-"]:has-text("${commentScenario.text}")`)).toBeVisible();
    }

    // Test comment interaction (like a comment)
    const firstComment = page.locator('[data-testid^="social-comment-item-"]').first();
    if (await firstComment.isVisible()) {
      const commentId = await firstComment.getAttribute('data-comment-id');
      await page.click(`[data-testid="social-comment-like-button-${commentId}"]`);

      // Verify comment like count updates
      const commentLikeCount = page.locator(`[data-testid="social-comment-like-count-${commentId}"]`);
      await expect(commentLikeCount).not.toHaveText('0');
    }
  });

  test('should perform complete share workflow', async ({ page }) => {
    const firstPost = page.locator('[data-testid^="social-post-"]').first();
    const postId = await firstPost.getAttribute('data-testid');
    const postIdValue = postId?.replace('social-post-', '');

    // Test share button click
    await page.click(`[data-testid="social-post-share-button-${postIdValue}"]`);
    await expect(page.locator('.toast')).toContainText('Post link copied to clipboard');

    // Test share via options menu
    await page.click(`[data-testid="social-post-options-trigger-${postIdValue}"]`);
    await expect(page.locator(`[data-testid="social-post-options-content-${postIdValue}"]`)).toBeVisible();

    await page.click(`[data-testid="social-post-share-item-${postIdValue}"]`);
    await expect(page.locator('.toast')).toContainText('Post link copied to clipboard');
  });

  test('should perform complete bookmark workflow', async ({ page }) => {
    const firstPost = page.locator('[data-testid^="social-post-"]').first();
    const postId = await firstPost.getAttribute('data-testid');
    const postIdValue = postId?.replace('social-post-', '');

    const bookmarkButton = page.locator(`[data-testid="social-post-bookmark-button-${postIdValue}"]`);

    // Test bookmarking
    await bookmarkButton.click();
    await expect(bookmarkButton).toHaveAttribute('data-bookmarked', 'true');
    await expect(page.locator('.toast')).toContainText('Added to bookmarks');

    // Test unbookmarking
    await bookmarkButton.click();
    await expect(bookmarkButton).toHaveAttribute('data-bookmarked', 'false');
    await expect(page.locator('.toast')).toContainText('Removed from bookmarks');

    // Test bookmark via options menu
    await page.click(`[data-testid="social-post-options-trigger-${postIdValue}"]`);
    await page.click(`[data-testid="social-post-bookmark-item-${postIdValue}"]`);
    await expect(page.locator('.toast')).toContainText('Added to bookmarks');
  });

  test('should handle engagement on live posts', async ({ page }) => {
    // Mock a live post
    await page.evaluate(() => {
      // Simulate live post data being added to the feed
      window.postMessage({
        type: 'LIVE_POST_ADDED',
        payload: {
          id: 'live-post-1',
          type: 'live',
          liveData: { isLive: true, viewers: 1250 }
        }
      }, '*');
    });

    // Look for live posts in the feed
    const livePosts = page.locator('[data-testid^="social-post-"][data-post-type="live"]');

    if (await livePosts.count() > 0) {
      const firstLivePost = livePosts.first();
      const postId = await firstLivePost.getAttribute('data-testid');
      const postIdValue = postId?.replace('social-post-', '');

      // Verify live badge is present
      await expect(page.locator(`[data-testid="social-post-live-badge-${postIdValue}"]`)).toBeVisible();

      // Test joining live stream
      const joinButton = page.locator(`[data-testid="social-post-join-live-button-${postIdValue}"]`);
      if (await joinButton.isVisible()) {
        await joinButton.click();
        await expect(page.locator('.toast')).toContainText('Joining live stream');
      }

      // Test normal engagement on live posts
      await page.click(`[data-testid="social-post-like-button-${postIdValue}"]`);
      await expect(page.locator('.toast')).toContainText('Post liked');
    }
  });

  test('should handle rapid engagement interactions', async ({ page }) => {
    const firstPost = page.locator('[data-testid^="social-post-"]').first();
    const postId = await firstPost.getAttribute('data-testid');
    const postIdValue = postId?.replace('social-post-', '');

    // Rapidly click like button multiple times
    const likeButton = page.locator(`[data-testid="social-post-like-button-${postIdValue}"]`);

    for (let i = 0; i < 5; i++) {
      await likeButton.click();
      await page.waitForTimeout(100);
    }

    // Verify final state is consistent
    const finalLikedState = await likeButton.getAttribute('data-liked');
    expect(['true', 'false']).toContain(finalLikedState);

    // Test rapid commenting
    await page.click(`[data-testid="social-post-comment-button-${postIdValue}"]`);

    const commentInput = page.locator(`[data-testid="social-comment-input-field-${postIdValue}"]`);
    const submitButton = page.locator(`[data-testid="social-comment-submit-button-${postIdValue}"]`);

    for (let i = 0; i < 3; i++) {
      await commentInput.fill(`Rapid comment ${i + 1}`);
      await submitButton.click();
      await page.waitForTimeout(200);
    }

    // Verify all comments appear
    for (let i = 0; i < 3; i++) {
      await expect(page.locator(`[data-testid^="social-comment-text-"]:has-text("Rapid comment ${i + 1}")`)).toBeVisible();
    }
  });
});

test.describe('Social Post Engagement - Error Scenarios', () => {
  test.beforeEach(async ({ page }) => {
    await navigateToSocialTab(page);
  });

  test('should handle like API failures gracefully', async ({ page }) => {
    // Mock API failure for likes
    await page.route('**/api/v1/social/posts/*/like', (route) => {
      route.fulfill({ status: 500, body: 'Server Error' });
    });

    // Create a post and try to like it
    await fillPostContent(page, 'Test post for like failure');
    await submitPost(page);
    await page.click('[data-testid="social-tab-feed"]');

    const firstPost = page.locator('[data-testid^="social-post-"]').first();
    const postId = await firstPost.getAttribute('data-testid');
    const postIdValue = postId?.replace('social-post-', '');

    // Try to like the post
    await page.click(`[data-testid="social-post-like-button-${postIdValue}"]`);

    // Should show error feedback
    await expect(page.locator('.toast')).toContainText('Failed to like post');

    // Like state should be reverted
    await expect(page.locator(`[data-testid="social-post-like-button-${postIdValue}"]`)).toHaveAttribute('data-liked', 'false');
  });

  test('should handle comment API failures gracefully', async ({ page }) => {
    // Mock API failure for comments
    await page.route('**/api/v1/social/posts/*/comments', (route) => {
      route.fulfill({ status: 500, body: 'Server Error' });
    });

    // Create a post and try to comment
    await fillPostContent(page, 'Test post for comment failure');
    await submitPost(page);
    await page.click('[data-testid="social-tab-feed"]');

    const firstPost = page.locator('[data-testid^="social-post-"]').first();
    const postId = await firstPost.getAttribute('data-testid');
    const postIdValue = postId?.replace('social-post-', '');

    // Open comments and try to add one
    await page.click(`[data-testid="social-post-comment-button-${postIdValue}"]`);
    await page.fill(`[data-testid="social-comment-input-field-${postIdValue}"]`, 'Test comment');
    await page.click(`[data-testid="social-comment-submit-button-${postIdValue}"]`);

    // Should show error feedback
    await expect(page.locator('.toast')).toContainText('Failed to add comment');

    // Comment input should not be cleared
    await expect(page.locator(`[data-testid="social-comment-input-field-${postIdValue}"]`)).toHaveValue('Test comment');
  });

  test('should handle post creation API failures', async ({ page }) => {
    // Mock API failure for post creation
    await page.route('**/api/v1/social/posts', (route) => {
      if (route.request().method() === 'POST') {
        route.fulfill({ status: 500, body: 'Server Error' });
      }
    });

    // Try to create a post
    await fillPostContent(page, 'Test post that will fail');
    await page.click('[data-testid="create-post-submit-button"]');

    // Should show error feedback
    await expect(page.locator('.toast')).toContainText('Failed to create post');

    // Content should remain in textarea for retry
    await expect(page.locator('[data-testid="create-post-textarea-input"]')).toHaveValue('Test post that will fail');
  });

  test('should handle network timeouts', async ({ page }) => {
    // Mock slow API responses
    await page.route('**/api/v1/social/**', async (route) => {
      await new Promise(resolve => setTimeout(resolve, 5000)); // 5 second delay
      route.fulfill({ status: 408, body: 'Request Timeout' });
    });

    // Try to create a post
    await fillPostContent(page, 'Test post with timeout');
    await page.click('[data-testid="create-post-submit-button"]');

    // Should eventually show timeout error
    await expect(page.locator('.toast')).toContainText('Request failed', { timeout: 10000 });
  });
});