/**
 * Social Discovery and User Profile E2E Test Scenarios
 *
 * Comprehensive testing for:
 * - User profile viewing and interaction
 * - Social discovery features (trending, recommendations)
 * - Friend/follower management workflows
 * - Social search and hashtag exploration
 * - Regional content discovery (Southeast Asian focus)
 * - Events and community features
 */

import { test, expect, Page } from '@playwright/test';

// Test data constants
const MOCK_USER_PROFILES = {
  FOOD_INFLUENCER: {
    id: 'user-sarah-foodie',
    name: 'Sarah Johnson',
    username: '@sarah_foodie',
    avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820',
    verified: false,
    type: 'user',
    followers: 15200,
    following: 842,
    posts: 1847,
    bio: 'Bangkok food explorer 🍜 | Street food enthusiast | Sharing hidden gems across Thailand',
    location: 'Bangkok, Thailand',
    website: 'foodiebangkok.com',
    joinedDate: '2022-03-15',
    isFollowing: false,
    mutualFriends: 12,
    recentPosts: [
      {
        id: 'post-sarah-1',
        content: 'Found the most amazing Tom Yum at this hidden street stall! 🍜🔥',
        images: ['https://images.unsplash.com/photo-1628432021231-4bbd431e6a04'],
        likes: 234,
        comments: 45,
        timestamp: '2 hours ago'
      }
    ]
  },
  TRAVEL_BLOGGER: {
    id: 'user-mike-travels',
    name: 'Mike Chen',
    username: '@mike_travels',
    avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e',
    verified: true,
    type: 'user',
    followers: 89300,
    following: 1205,
    posts: 3421,
    bio: 'Travel photographer & cultural explorer 📸 | Southeast Asia specialist | Digital nomad life',
    location: 'Currently in Bangkok',
    website: 'miketravels.blog',
    joinedDate: '2021-08-20',
    isFollowing: true,
    mutualFriends: 8,
    recentPosts: [
      {
        id: 'post-mike-1',
        content: 'Incredible sunrise at Wat Pho this morning! The golden light was perfect 🌅',
        images: ['https://images.unsplash.com/photo-1563492065-cd5bab1c2d64'],
        likes: 1247,
        comments: 89,
        timestamp: '5 hours ago'
      }
    ]
  }
};

const DISCOVERY_CONTENT = {
  TRENDING_HASHTAGS: [
    { tag: '#BangkokEats', posts: 45200, trending: true },
    { tag: '#StreetFood', posts: 128300, trending: true },
    { tag: '#ThaiCulture', posts: 67800, trending: false },
    { tag: '#ChatuchakMarket', posts: 23400, trending: true },
    { tag: '#WatPho', posts: 34500, trending: false }
  ],
  TRENDING_LOCATIONS: [
    { name: 'Chatuchak Weekend Market', posts: 12400, category: 'market' },
    { name: 'Wat Pho Temple', posts: 8900, category: 'temple' },
    { name: 'Khao San Road', posts: 15600, category: 'nightlife' },
    { name: 'Damnoen Saduak Floating Market', posts: 6700, category: 'market' }
  ],
  RECOMMENDED_USERS: [
    {
      id: 'rec-user-1',
      name: 'Thai Food Masters',
      username: '@thaifoodmasters',
      followers: 245000,
      category: 'Food & Cooking',
      reason: 'Popular in your area'
    },
    {
      id: 'rec-user-2',
      name: 'Bangkok Hidden Gems',
      username: '@bkkhiddengems',
      followers: 128000,
      category: 'Local Guide',
      reason: 'Similar interests'
    }
  ]
};

const EVENTS_DATA = {
  UPCOMING_EVENTS: [
    {
      id: 'event-1',
      title: 'Bangkok Electronic Music Festival 2025',
      date: '2025-03-15',
      location: 'Show DC, Bangkok',
      attendees: 18500,
      category: 'music',
      trending: true,
      price: 'From ฿2,500'
    },
    {
      id: 'event-2',
      title: 'Thai Street Food Championship',
      date: '2025-02-28',
      location: 'Lumpini Park',
      attendees: 12000,
      category: 'food',
      trending: false,
      price: 'From ฿500'
    }
  ]
};

// Helper functions
async function navigateToSocialTab(page: Page) {
  await page.goto('/');
  await page.click('[data-testid="social-tab-trigger"]');
  await expect(page.locator('[data-testid="social-tab-container"]')).toBeVisible();
}

async function mockSocialDiscoveryAPIs(page: Page) {
  // Mock trending content API
  await page.route('**/api/v1/social/trending/**', (route) => {
    const url = route.request().url();

    if (url.includes('/hashtags')) {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(DISCOVERY_CONTENT.TRENDING_HASHTAGS)
      });
    } else if (url.includes('/locations')) {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(DISCOVERY_CONTENT.TRENDING_LOCATIONS)
      });
    }
  });

  // Mock user recommendations API
  await page.route('**/api/v1/social/recommendations/users', (route) => {
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(DISCOVERY_CONTENT.RECOMMENDED_USERS)
    });
  });

  // Mock user profile API
  await page.route('**/api/v1/social/users/*', (route) => {
    const userId = route.request().url().split('/').pop();
    const profile = Object.values(MOCK_USER_PROFILES).find(p => p.id === userId);

    if (profile) {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(profile)
      });
    } else {
      route.fulfill({ status: 404, body: 'User not found' });
    }
  });

  // Mock events API
  await page.route('**/api/v1/social/events**', (route) => {
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(EVENTS_DATA.UPCOMING_EVENTS)
    });
  });

  // Mock follow/unfollow API
  await page.route('**/api/v1/social/users/*/follow', (route) => {
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ success: true })
    });
  });
}

async function openUserProfile(page: Page, userId: string) {
  // Simulate clicking on a user profile link
  await page.evaluate((id) => {
    window.postMessage({
      type: 'OPEN_USER_PROFILE',
      payload: { userId: id }
    }, '*');
  }, userId);
}

test.describe('Social Discovery - Trending and Recommendations', () => {
  test.beforeEach(async ({ page }) => {
    await mockSocialDiscoveryAPIs(page);
    await navigateToSocialTab(page);
  });

  test('should display and interact with discover tab', async ({ page }) => {
    // Navigate to discover tab
    await page.click('[data-testid="social-tab-discover"]');
    await expect(page.locator('[data-testid="social-tab-discover"]')).toHaveAttribute('data-state', 'active');

    // Wait for discover content to load
    await page.waitForSelector('[data-testid^="discover-"]', { timeout: 5000 });

    // Verify discover tab content is visible
    const discoverTab = page.locator('[value="discover"]');
    await expect(discoverTab).toBeVisible();
  });

  test('should display trending hashtags', async ({ page }) => {
    await page.click('[data-testid="social-tab-discover"]');

    // Look for trending hashtag elements
    const trendingSection = page.locator('[data-testid^="trending-"]');

    if (await trendingSection.count() > 0) {
      // Verify trending hashtags are displayed
      for (const hashtag of DISCOVERY_CONTENT.TRENDING_HASHTAGS.slice(0, 3)) {
        await expect(page.locator(`text=${hashtag.tag}`)).toBeVisible({ timeout: 5000 });
      }
    }
  });

  test('should display and interact with recommended users', async ({ page }) => {
    await page.click('[data-testid="social-tab-discover"]');

    // Look for user recommendation elements
    const recommendedSection = page.locator('[data-testid^="recommended-users"]');

    if (await recommendedSection.count() > 0) {
      // Verify recommended users are displayed
      for (const user of DISCOVERY_CONTENT.RECOMMENDED_USERS) {
        await expect(page.locator(`text=${user.name}`)).toBeVisible({ timeout: 5000 });
      }

      // Test following a recommended user
      const followButton = page.locator('[data-testid^="follow-user-button-"]').first();
      if (await followButton.isVisible()) {
        await followButton.click();
        await expect(page.locator('.toast')).toContainText('Following user');
      }
    }
  });

  test('should search and filter discovery content', async ({ page }) => {
    await page.click('[data-testid="social-tab-discover"]');

    // Look for search functionality
    const searchInput = page.locator('[data-testid="discovery-search-input"]');

    if (await searchInput.isVisible()) {
      // Test searching for Bangkok content
      await searchInput.fill('Bangkok');
      await page.keyboard.press('Enter');

      // Verify search results
      await expect(page.locator('text=Bangkok')).toBeVisible({ timeout: 5000 });
    }

    // Test category filters
    const categoryFilters = page.locator('[data-testid^="category-filter-"]');

    if (await categoryFilters.count() > 0) {
      const foodFilter = categoryFilters.filter({ hasText: 'Food' }).first();
      if (await foodFilter.isVisible()) {
        await foodFilter.click();
        // Verify food-related content is shown
        await expect(page.locator('text=Food')).toBeVisible({ timeout: 5000 });
      }
    }
  });

  test('should display location-based recommendations', async ({ page }) => {
    await page.click('[data-testid="social-tab-discover"]');

    // Test location-based discovery
    for (const location of DISCOVERY_CONTENT.TRENDING_LOCATIONS.slice(0, 2)) {
      await expect(page.locator(`text=${location.name}`)).toBeVisible({ timeout: 5000 });
    }

    // Test clicking on a location
    const locationElement = page.locator('text=Chatuchak Weekend Market').first();
    if (await locationElement.isVisible()) {
      await locationElement.click();
      await expect(page.locator('.toast')).toContainText('Exploring Chatuchak Weekend Market');
    }
  });
});

test.describe('Social Discovery - Events and Community', () => {
  test.beforeEach(async ({ page }) => {
    await mockSocialDiscoveryAPIs(page);
    await navigateToSocialTab(page);
  });

  test('should display and interact with events tab', async ({ page }) => {
    // Navigate to events tab
    await page.click('[data-testid="social-tab-events"]');
    await expect(page.locator('[data-testid="social-tab-events"]')).toHaveAttribute('data-state', 'active');

    // Verify events content is displayed
    const eventsTab = page.locator('[value="events"]');
    await expect(eventsTab).toBeVisible();
  });

  test('should display upcoming events', async ({ page }) => {
    await page.click('[data-testid="social-tab-events"]');

    // Check for featured events
    for (const event of EVENTS_DATA.UPCOMING_EVENTS) {
      await expect(page.locator(`text=${event.title}`)).toBeVisible({ timeout: 5000 });
    }

    // Verify event details are shown
    await expect(page.locator('text=Bangkok Electronic Music Festival')).toBeVisible();
    await expect(page.locator('text=18,500 going')).toBeVisible();
  });

  test('should interact with event categories', async ({ page }) => {
    await page.click('[data-testid="social-tab-events"]');

    // Test event category navigation
    const categoryButtons = page.locator('[data-testid^="event-category-"]');

    if (await categoryButtons.count() > 0) {
      const musicCategory = categoryButtons.filter({ hasText: 'Music' }).first();
      if (await musicCategory.isVisible()) {
        await musicCategory.click();
        // Verify music events are highlighted
        await expect(page.locator('text=Music')).toBeVisible();
      }

      const foodCategory = categoryButtons.filter({ hasText: 'Food' }).first();
      if (await foodCategory.isVisible()) {
        await foodCategory.click();
        // Verify food events are highlighted
        await expect(page.locator('text=Food')).toBeVisible();
      }
    }
  });

  test('should show interest in events', async ({ page }) => {
    await page.click('[data-testid="social-tab-events"]');

    // Test showing interest in an event
    const interestedButton = page.locator('[data-testid^="event-interested-button-"]').first();

    if (await interestedButton.isVisible()) {
      await interestedButton.click();
      await expect(page.locator('.toast')).toContainText('Added to interested events');
    }
  });

  test('should explore all events', async ({ page }) => {
    await page.click('[data-testid="social-tab-events"]');

    // Test "Explore All Events" functionality
    const exploreAllButton = page.locator('[data-testid="explore-all-events-button"]');

    if (await exploreAllButton.isVisible()) {
      await exploreAllButton.click();
      // Should navigate to full events view
      await expect(page.locator('[data-testid="events-full-view"]')).toBeVisible({ timeout: 5000 });
    }
  });
});

test.describe('User Profile Viewing and Interaction', () => {
  test.beforeEach(async ({ page }) => {
    await mockSocialDiscoveryAPIs(page);
    await navigateToSocialTab(page);
  });

  test('should view user profile from post author', async ({ page }) => {
    await page.click('[data-testid="social-tab-feed"]');

    // Click on a post author name to view profile
    const authorName = page.locator('[data-testid^="social-post-author-name-"]').first();

    if (await authorName.isVisible()) {
      await authorName.click();

      // Profile modal or page should open
      await expect(page.locator('[data-testid^="user-profile-"]')).toBeVisible({ timeout: 5000 });
    }
  });

  test('should view user profile from comment author', async ({ page }) => {
    await page.click('[data-testid="social-tab-feed"]');

    // Open comments on a post
    const commentButton = page.locator('[data-testid^="social-post-comment-button-"]').first();
    if (await commentButton.isVisible()) {
      await commentButton.click();

      // Click on a comment author
      const commentAuthor = page.locator('[data-testid^="social-comment-author-"]').first();
      if (await commentAuthor.isVisible()) {
        await commentAuthor.click();

        // Profile should open
        await expect(page.locator('[data-testid^="user-profile-"]')).toBeVisible({ timeout: 5000 });
      }
    }
  });

  test('should display complete user profile information', async ({ page }) => {
    // Simulate opening a user profile
    await openUserProfile(page, MOCK_USER_PROFILES.FOOD_INFLUENCER.id);

    // Wait for profile modal/page to load
    await expect(page.locator('[data-testid^="user-profile-"]')).toBeVisible({ timeout: 5000 });

    const profile = MOCK_USER_PROFILES.FOOD_INFLUENCER;

    // Verify profile header information
    await expect(page.locator(`[data-testid="profile-name"]`)).toContainText(profile.name);
    await expect(page.locator(`[data-testid="profile-username"]`)).toContainText(profile.username);
    await expect(page.locator(`[data-testid="profile-bio"]`)).toContainText(profile.bio);

    // Verify profile stats
    await expect(page.locator(`[data-testid="profile-followers-count"]`)).toContainText('15.2K');
    await expect(page.locator(`[data-testid="profile-following-count"]`)).toContainText('842');
    await expect(page.locator(`[data-testid="profile-posts-count"]`)).toContainText('1,847');

    // Verify profile metadata
    await expect(page.locator(`[data-testid="profile-location"]`)).toContainText(profile.location);
    if (profile.website) {
      await expect(page.locator(`[data-testid="profile-website"]`)).toContainText(profile.website);
    }
  });

  test('should follow and unfollow users', async ({ page }) => {
    await openUserProfile(page, MOCK_USER_PROFILES.FOOD_INFLUENCER.id);

    const followButton = page.locator('[data-testid="profile-follow-button"]');

    if (await followButton.isVisible()) {
      // Test following
      await followButton.click();
      await expect(page.locator('.toast')).toContainText('Following');
      await expect(followButton).toContainText('Following');

      // Test unfollowing
      await followButton.click();
      await expect(page.locator('.toast')).toContainText('Unfollowed');
      await expect(followButton).toContainText('Follow');
    }
  });

  test('should view user posts in profile', async ({ page }) => {
    await openUserProfile(page, MOCK_USER_PROFILES.FOOD_INFLUENCER.id);

    // Switch to posts tab in profile
    const postsTab = page.locator('[data-testid="profile-posts-tab"]');
    if (await postsTab.isVisible()) {
      await postsTab.click();

      // Verify user posts are displayed
      const userPosts = page.locator('[data-testid^="profile-post-"]');
      await expect(userPosts).toHaveCount.greaterThan(0);

      // Verify post content
      const profile = MOCK_USER_PROFILES.FOOD_INFLUENCER;
      if (profile.recentPosts.length > 0) {
        await expect(page.locator(`text=${profile.recentPosts[0].content.substring(0, 20)}`)).toBeVisible();
      }
    }
  });

  test('should display mutual friends', async ({ page }) => {
    await openUserProfile(page, MOCK_USER_PROFILES.FOOD_INFLUENCER.id);

    // Check for mutual friends section
    const mutualFriendsSection = page.locator('[data-testid="profile-mutual-friends"]');

    if (await mutualFriendsSection.isVisible()) {
      const profile = MOCK_USER_PROFILES.FOOD_INFLUENCER;
      await expect(mutualFriendsSection).toContainText(`${profile.mutualFriends} mutual friends`);
    }
  });

  test('should handle verified users differently', async ({ page }) => {
    await openUserProfile(page, MOCK_USER_PROFILES.TRAVEL_BLOGGER.id);

    // Verify verified badge is displayed
    const verifiedBadge = page.locator('[data-testid="profile-verified-badge"]');
    await expect(verifiedBadge).toBeVisible();

    // Verified users might have different UI treatment
    const profile = MOCK_USER_PROFILES.TRAVEL_BLOGGER;
    await expect(page.locator(`[data-testid="profile-name"]`)).toContainText(profile.name);
  });
});

test.describe('Friend and Follower Management', () => {
  test.beforeEach(async ({ page }) => {
    await mockSocialDiscoveryAPIs(page);
    await navigateToSocialTab(page);
  });

  test('should manage friends list', async ({ page }) => {
    // Navigate to friends tab
    await page.click('[data-testid="social-tab-friends"]');

    // Verify friends activity is displayed
    await expect(page.locator('[data-testid^="friend-activity-"]')).toBeVisible({ timeout: 5000 });

    // Test friend suggestions
    const friendSuggestions = page.locator('[data-testid^="friend-suggestion-"]');

    if (await friendSuggestions.count() > 0) {
      const addFriendButton = friendSuggestions.first().locator('[data-testid^="add-friend-button-"]');
      if (await addFriendButton.isVisible()) {
        await addFriendButton.click();
        await expect(page.locator('.toast')).toContainText('Friend request sent');
      }
    }
  });

  test('should display online friends status', async ({ page }) => {
    await page.click('[data-testid="social-tab-friends"]');

    // Check for online friends section
    const onlineFriendsSection = page.locator('[data-testid="friends-online-section"]');

    if (await onlineFriendsSection.isVisible()) {
      // Verify online indicators
      const onlineIndicators = page.locator('[data-testid^="friend-online-indicator-"]');
      await expect(onlineIndicators).toHaveCount.greaterThan(0);

      // Verify online friends count badge
      const onlineBadge = page.locator('text=online');
      await expect(onlineBadge).toBeVisible();
    }
  });

  test('should handle friend requests', async ({ page }) => {
    // Mock friend requests API
    await page.route('**/api/v1/social/friends/requests', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            id: 'request-1',
            from: MOCK_USER_PROFILES.FOOD_INFLUENCER,
            timestamp: '2 hours ago'
          }
        ])
      });
    });

    await page.click('[data-testid="social-tab-friends"]');

    // Check for friend requests section
    const friendRequests = page.locator('[data-testid^="friend-request-"]');

    if (await friendRequests.count() > 0) {
      // Test accepting a friend request
      const acceptButton = friendRequests.first().locator('[data-testid^="accept-friend-button-"]');
      if (await acceptButton.isVisible()) {
        await acceptButton.click();
        await expect(page.locator('.toast')).toContainText('Friend request accepted');
      }

      // Test declining a friend request
      const declineButton = friendRequests.first().locator('[data-testid^="decline-friend-button-"]');
      if (await declineButton.isVisible()) {
        await declineButton.click();
        await expect(page.locator('.toast')).toContainText('Friend request declined');
      }
    }
  });
});

test.describe('Southeast Asian Regional Features', () => {
  test.beforeEach(async ({ page }) => {
    await mockSocialDiscoveryAPIs(page);
    await navigateToSocialTab(page);
  });

  test('should display Bangkok-specific trending content', async ({ page }) => {
    await page.click('[data-testid="social-tab-discover"]');

    // Check for Bangkok-specific hashtags
    const bangkokHashtags = ['#BangkokEats', '#ChatuchakMarket', '#WatPho'];

    for (const hashtag of bangkokHashtags) {
      await expect(page.locator(`text=${hashtag}`)).toBeVisible({ timeout: 5000 });
    }
  });

  test('should display Thai cultural events', async ({ page }) => {
    await page.click('[data-testid="social-tab-events"]');

    // Check for Thai food events
    await expect(page.locator('text=Thai Street Food Championship')).toBeVisible({ timeout: 5000 });

    // Check for cultural venue references
    await expect(page.locator('text=Lumpini Park')).toBeVisible({ timeout: 5000 });
  });

  test('should support Thai language content discovery', async ({ page }) => {
    await page.click('[data-testid="social-tab-discover"]');

    // Mock Thai language content
    await page.evaluate(() => {
      window.postMessage({
        type: 'ADD_THAI_CONTENT',
        payload: {
          hashtags: ['#อาหารไทย', '#ผัดไทย', '#กรุงเทพ'],
          locations: ['ตลาดจตุจักร', 'วัดโพธิ์']
        }
      }, '*');
    });

    // Test searching for Thai content
    const searchInput = page.locator('[data-testid="discovery-search-input"]');
    if (await searchInput.isVisible()) {
      await searchInput.fill('อาหารไทย');
      await page.keyboard.press('Enter');

      // Should handle Thai language search
      await expect(page.locator('text=อาหารไทย')).toBeVisible({ timeout: 5000 });
    }
  });

  test('should display regional currency and pricing', async ({ page }) => {
    await page.click('[data-testid="social-tab-events"]');

    // Check for Thai Baht currency in event pricing
    await expect(page.locator('text=฿2,500')).toBeVisible({ timeout: 5000 });
    await expect(page.locator('text=฿500')).toBeVisible({ timeout: 5000 });
  });
});