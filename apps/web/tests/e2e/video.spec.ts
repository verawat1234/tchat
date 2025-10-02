import { test, expect, Page } from '@playwright/test';

// Test configuration
const BASE_URL = process.env.VITE_API_URL || 'http://localhost:3000';
const API_URL = process.env.VITE_API_URL || 'http://localhost:8080';

// Mock video data
const mockVideo = {
  id: 'e2e-test-video-1',
  title: 'E2E Test Video',
  description: 'Video for end-to-end testing',
  videoUrl: 'https://test-videos.co.uk/vids/bigbuckbunny/mp4/h264/360/Big_Buck_Bunny_360_10s_1MB.mp4',
  thumbnailUrl: 'https://via.placeholder.com/640x360',
  durationSeconds: 10,
  creatorId: 'test-creator',
  uploadStatus: 'available',
};

// Helper functions
async function navigateToVideoPage(page: Page, videoId: string) {
  await page.goto(`${BASE_URL}/videos/${videoId}`);
  await page.waitForLoadState('networkidle');
}

async function waitForVideoPlayer(page: Page) {
  await page.waitForSelector('[data-testid="video-player"]', { timeout: 10000 });
}

async function mockVideoAPI(page: Page) {
  await page.route(`${API_URL}/api/v1/videos/*`, (route) => {
    if (route.request().method() === 'GET') {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockVideo),
      });
    } else {
      route.continue();
    }
  });
}

test.describe('Video Player E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    await mockVideoAPI(page);
  });

  test('should load video player page', async ({ page }) => {
    await navigateToVideoPage(page, mockVideo.id);

    // Check page title
    await expect(page).toHaveTitle(/Tchat - Video/);

    // Video player should be visible
    await waitForVideoPlayer(page);
    const videoPlayer = page.locator('[data-testid="video-player"]');
    await expect(videoPlayer).toBeVisible();
  });

  test('should display video metadata', async ({ page }) => {
    await navigateToVideoPage(page, mockVideo.id);
    await waitForVideoPlayer(page);

    // Video title should be displayed
    const title = page.locator('h1', { hasText: mockVideo.title });
    await expect(title).toBeVisible();

    // Video description should be displayed
    const description = page.locator('text=' + mockVideo.description);
    await expect(description).toBeVisible();
  });

  test('should play video when play button is clicked', async ({ page }) => {
    await navigateToVideoPage(page, mockVideo.id);
    await waitForVideoPlayer(page);

    // Click play button
    const playButton = page.locator('[data-testid="play-pause-button"]');
    await playButton.click();

    // Video should start playing (check for pause button icon)
    await expect(playButton).toHaveAttribute('aria-label', /pause/i);
  });

  test('should pause video when pause button is clicked', async ({ page }) => {
    await navigateToVideoPage(page, mockVideo.id);
    await waitForVideoPlayer(page);

    // Start playing
    const playButton = page.locator('[data-testid="play-pause-button"]');
    await playButton.click();
    await page.waitForTimeout(1000);

    // Pause video
    await playButton.click();

    // Video should be paused (check for play button icon)
    await expect(playButton).toHaveAttribute('aria-label', /play/i);
  });

  test('should update progress bar during playback', async ({ page }) => {
    await navigateToVideoPage(page, mockVideo.id);
    await waitForVideoPlayer(page);

    // Start playing
    const playButton = page.locator('[data-testid="play-pause-button"]');
    await playButton.click();

    // Get initial progress
    const progressBar = page.locator('[data-testid="progress-bar"]');
    const initialValue = await progressBar.getAttribute('aria-valuenow');

    // Wait for 2 seconds
    await page.waitForTimeout(2000);

    // Progress should have advanced
    const newValue = await progressBar.getAttribute('aria-valuenow');
    expect(Number(newValue)).toBeGreaterThan(Number(initialValue));
  });

  test('should seek video when progress bar is clicked', async ({ page }) => {
    await navigateToVideoPage(page, mockVideo.id);
    await waitForVideoPlayer(page);

    // Start playing
    const playButton = page.locator('[data-testid="play-pause-button"]');
    await playButton.click();

    // Click on middle of progress bar
    const progressBar = page.locator('[data-testid="progress-bar"]');
    const box = await progressBar.boundingBox();
    if (box) {
      await page.mouse.click(box.x + box.width / 2, box.y + box.height / 2);

      // Wait for seek to complete
      await page.waitForTimeout(500);

      // Position should be around middle (5 seconds for 10 second video)
      const currentTime = await progressBar.getAttribute('aria-valuenow');
      expect(Number(currentTime)).toBeGreaterThan(3);
      expect(Number(currentTime)).toBeLessThan(7);
    }
  });

  test('should adjust volume', async ({ page }) => {
    await navigateToVideoPage(page, mockVideo.id);
    await waitForVideoPlayer(page);

    // Click volume button to show slider
    const volumeButton = page.locator('[data-testid="volume-button"]');
    await volumeButton.click();

    // Adjust volume to 50%
    const volumeSlider = page.locator('[data-testid="volume-slider"]');
    await volumeSlider.fill('0.5');

    // Volume should be updated
    const volume = await volumeSlider.getAttribute('value');
    expect(Number(volume)).toBe(0.5);
  });

  test('should mute and unmute video', async ({ page }) => {
    await navigateToVideoPage(page, mockVideo.id);
    await waitForVideoPlayer(page);

    // Click mute button
    const volumeButton = page.locator('[data-testid="volume-button"]');
    await volumeButton.click();

    // Should show muted icon
    await expect(volumeButton).toHaveAttribute('aria-label', /unmute/i);

    // Click again to unmute
    await volumeButton.click();

    // Should show unmuted icon
    await expect(volumeButton).toHaveAttribute('aria-label', /mute/i);
  });

  test('should enter and exit fullscreen', async ({ page }) => {
    await navigateToVideoPage(page, mockVideo.id);
    await waitForVideoPlayer(page);

    // Click fullscreen button
    const fullscreenButton = page.locator('[data-testid="fullscreen-button"]');
    await fullscreenButton.click();

    // Wait for fullscreen transition
    await page.waitForTimeout(500);

    // Video player should be in fullscreen
    const videoPlayer = page.locator('[data-testid="video-player"]');
    const isFullscreen = await videoPlayer.evaluate((el) => document.fullscreenElement === el);
    expect(isFullscreen).toBeTruthy();

    // Exit fullscreen
    await fullscreenButton.click();
    await page.waitForTimeout(500);
  });

  test('should change video quality', async ({ page }) => {
    await navigateToVideoPage(page, mockVideo.id);
    await waitForVideoPlayer(page);

    // Click settings button
    const settingsButton = page.locator('[data-testid="settings-button"]');
    await settingsButton.click();

    // Select quality option
    const qualityOption = page.locator('text=Quality');
    await qualityOption.click();

    // Select 720p
    const quality720p = page.locator('text=720p');
    await quality720p.click();

    // Quality setting should be applied
    await expect(page.locator('text=720p')).toBeVisible();
  });

  test('should change playback speed', async ({ page }) => {
    await navigateToVideoPage(page, mockVideo.id);
    await waitForVideoPlayer(page);

    // Click settings button
    const settingsButton = page.locator('[data-testid="settings-button"]');
    await settingsButton.click();

    // Select speed option
    const speedOption = page.locator('text=Speed');
    await speedOption.click();

    // Select 1.5x speed
    const speed15x = page.locator('text=1.5x');
    await speed15x.click();

    // Speed setting should be applied
    await expect(page.locator('text=1.5x')).toBeVisible();
  });

  test('should display video loading state', async ({ page }) => {
    // Mock slow network
    await page.route(`${API_URL}/api/v1/videos/*`, async (route) => {
      await page.waitForTimeout(2000); // 2 second delay
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockVideo),
      });
    });

    await navigateToVideoPage(page, mockVideo.id);

    // Loading indicator should be visible
    const loadingSpinner = page.locator('[data-testid="video-loading"]');
    await expect(loadingSpinner).toBeVisible();

    // Wait for video to load
    await waitForVideoPlayer(page);

    // Loading should be gone
    await expect(loadingSpinner).not.toBeVisible();
  });

  test('should handle video error gracefully', async ({ page }) => {
    // Mock API error
    await page.route(`${API_URL}/api/v1/videos/*`, (route) => {
      route.fulfill({
        status: 404,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Video not found' }),
      });
    });

    await navigateToVideoPage(page, 'non-existent-video');

    // Error message should be displayed
    const errorMessage = page.locator('text=Video not found');
    await expect(errorMessage).toBeVisible();
  });
});

test.describe('Video Upload E2E Tests', () => {
  test('should navigate to upload page', async ({ page }) => {
    await page.goto(`${BASE_URL}/videos/upload`);
    await page.waitForLoadState('networkidle');

    // Upload form should be visible
    const uploadForm = page.locator('[data-testid="upload-form"]');
    await expect(uploadForm).toBeVisible();
  });

  test('should display upload form fields', async ({ page }) => {
    await page.goto(`${BASE_URL}/videos/upload`);

    // Check all form fields are present
    await expect(page.locator('[data-testid="title-input"]')).toBeVisible();
    await expect(page.locator('[data-testid="description-input"]')).toBeVisible();
    await expect(page.locator('[data-testid="file-picker"]')).toBeVisible();
    await expect(page.locator('[data-testid="upload-button"]')).toBeVisible();
  });

  test('should validate required fields', async ({ page }) => {
    await page.goto(`${BASE_URL}/videos/upload`);

    // Try to submit without filling fields
    const uploadButton = page.locator('[data-testid="upload-button"]');
    await uploadButton.click();

    // Validation errors should be shown
    await expect(page.locator('text=Title is required')).toBeVisible();
    await expect(page.locator('text=Description is required')).toBeVisible();
  });

  test('should fill upload form', async ({ page }) => {
    await page.goto(`${BASE_URL}/videos/upload`);

    // Fill title
    const titleInput = page.locator('[data-testid="title-input"]');
    await titleInput.fill('Test Upload Video');

    // Fill description
    const descriptionInput = page.locator('[data-testid="description-input"]');
    await descriptionInput.fill('This is a test video upload');

    // Verify values
    await expect(titleInput).toHaveValue('Test Upload Video');
    await expect(descriptionInput).toHaveValue('This is a test video upload');
  });
});

test.describe('Video List E2E Tests', () => {
  const mockVideos = [
    { ...mockVideo, id: 'video-1', title: 'Video 1' },
    { ...mockVideo, id: 'video-2', title: 'Video 2' },
    { ...mockVideo, id: 'video-3', title: 'Video 3' },
  ];

  test('should display video list', async ({ page }) => {
    // Mock videos API
    await page.route(`${API_URL}/api/v1/videos`, (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ videos: mockVideos, total: 3, page: 1 }),
      });
    });

    await page.goto(`${BASE_URL}/videos`);
    await page.waitForLoadState('networkidle');

    // All videos should be displayed
    for (const video of mockVideos) {
      await expect(page.locator(`text=${video.title}`)).toBeVisible();
    }
  });

  test('should navigate to video when clicked', async ({ page }) => {
    await page.route(`${API_URL}/api/v1/videos`, (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ videos: mockVideos, total: 3, page: 1 }),
      });
    });

    await page.goto(`${BASE_URL}/videos`);
    await page.waitForLoadState('networkidle');

    // Click on first video
    const firstVideo = page.locator(`text=${mockVideos[0].title}`);
    await firstVideo.click();

    // Should navigate to video page
    await page.waitForURL(`${BASE_URL}/videos/${mockVideos[0].id}`);
  });

  test('should filter videos by category', async ({ page }) => {
    await page.route(`${API_URL}/api/v1/videos*`, (route) => {
      const url = route.request().url();
      const filteredVideos = url.includes('category=tutorial')
        ? [mockVideos[0]]
        : mockVideos;

      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ videos: filteredVideos, total: filteredVideos.length }),
      });
    });

    await page.goto(`${BASE_URL}/videos`);

    // Select category filter
    const categoryFilter = page.locator('[data-testid="category-filter"]');
    await categoryFilter.click();
    await page.locator('text=Tutorial').click();

    // Only tutorial videos should be shown
    await expect(page.locator(`text=${mockVideos[0].title}`)).toBeVisible();
  });

  test('should search videos', async ({ page }) => {
    await page.route(`${API_URL}/api/v1/videos*`, (route) => {
      const url = route.request().url();
      const searchTerm = new URL(url).searchParams.get('query') || '';
      const filteredVideos = mockVideos.filter((v) =>
        v.title.toLowerCase().includes(searchTerm.toLowerCase())
      );

      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ videos: filteredVideos, total: filteredVideos.length }),
      });
    });

    await page.goto(`${BASE_URL}/videos`);

    // Enter search query
    const searchInput = page.locator('[data-testid="search-input"]');
    await searchInput.fill('Video 1');
    await searchInput.press('Enter');

    // Wait for search results
    await page.waitForTimeout(500);

    // Only matching video should be shown
    await expect(page.locator('text=Video 1')).toBeVisible();
    await expect(page.locator('text=Video 2')).not.toBeVisible();
  });

  test('should load more videos on scroll', async ({ page }) => {
    const page1Videos = mockVideos.slice(0, 2);
    const page2Videos = [mockVideos[2]];

    let requestCount = 0;
    await page.route(`${API_URL}/api/v1/videos*`, (route) => {
      requestCount++;
      const videos = requestCount === 1 ? page1Videos : page2Videos;
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          videos,
          total: 3,
          page: requestCount,
          hasMore: requestCount === 1,
        }),
      });
    });

    await page.goto(`${BASE_URL}/videos`);
    await page.waitForLoadState('networkidle');

    // Scroll to bottom
    await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight));

    // Wait for more videos to load
    await page.waitForTimeout(1000);

    // Third video should be loaded
    await expect(page.locator('text=Video 3')).toBeVisible();
  });
});

test.describe('Video Social Interactions E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    await mockVideoAPI(page);
  });

  test('should like video', async ({ page }) => {
    await page.route(`${API_URL}/api/v1/videos/*/like`, (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ success: true, liked: true }),
      });
    });

    await navigateToVideoPage(page, mockVideo.id);
    await waitForVideoPlayer(page);

    // Click like button
    const likeButton = page.locator('[data-testid="like-button"]');
    await likeButton.click();

    // Like button should show active state
    await expect(likeButton).toHaveClass(/active/);
  });

  test('should comment on video', async ({ page }) => {
    await page.route(`${API_URL}/api/v1/videos/*/comments`, (route) => {
      if (route.request().method() === 'POST') {
        route.fulfill({
          status: 201,
          contentType: 'application/json',
          body: JSON.stringify({
            id: 'comment-1',
            text: 'Great video!',
            createdAt: new Date().toISOString(),
          }),
        });
      } else {
        route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ comments: [], total: 0 }),
        });
      }
    });

    await navigateToVideoPage(page, mockVideo.id);
    await waitForVideoPlayer(page);

    // Enter comment
    const commentInput = page.locator('[data-testid="comment-input"]');
    await commentInput.fill('Great video!');

    // Submit comment
    const submitButton = page.locator('[data-testid="submit-comment"]');
    await submitButton.click();

    // Comment should appear
    await expect(page.locator('text=Great video!')).toBeVisible();
  });

  test('should share video', async ({ page }) => {
    await navigateToVideoPage(page, mockVideo.id);
    await waitForVideoPlayer(page);

    // Click share button
    const shareButton = page.locator('[data-testid="share-button"]');
    await shareButton.click();

    // Share modal should open
    const shareModal = page.locator('[data-testid="share-modal"]');
    await expect(shareModal).toBeVisible();

    // Should show share options
    await expect(page.locator('text=Copy link')).toBeVisible();
  });
});

test.describe('Performance Tests', () => {
  test('should load video page within performance budget', async ({ page }) => {
    const startTime = Date.now();

    await navigateToVideoPage(page, mockVideo.id);
    await waitForVideoPlayer(page);

    const loadTime = Date.now() - startTime;

    // Should load within 3 seconds
    expect(loadTime).toBeLessThan(3000);
  });

  test('should maintain 60fps during playback', async ({ page }) => {
    await navigateToVideoPage(page, mockVideo.id);
    await waitForVideoPlayer(page);

    // Start playing
    const playButton = page.locator('[data-testid="play-pause-button"]');
    await playButton.click();

    // Measure frame rate
    const fps = await page.evaluate(() => {
      return new Promise<number>((resolve) => {
        let frames = 0;
        let lastTime = performance.now();

        const measureFPS = () => {
          frames++;
          const currentTime = performance.now();
          if (currentTime >= lastTime + 1000) {
            resolve(frames);
          } else {
            requestAnimationFrame(measureFPS);
          }
        };

        requestAnimationFrame(measureFPS);
      });
    });

    // Should maintain close to 60fps (allow some variance)
    expect(fps).toBeGreaterThan(55);
  });

  test('should use reasonable memory', async ({ page }) => {
    await navigateToVideoPage(page, mockVideo.id);
    await waitForVideoPlayer(page);

    // Get memory usage
    const metrics = await page.evaluate(() => {
      if ('memory' in performance) {
        return (performance as any).memory;
      }
      return null;
    });

    if (metrics) {
      // Memory usage should be under 500MB
      const memoryMB = metrics.usedJSHeapSize / 1024 / 1024;
      expect(memoryMB).toBeLessThan(500);
    }
  });
});