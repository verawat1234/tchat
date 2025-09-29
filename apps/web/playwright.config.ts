import { defineConfig, devices } from '@playwright/test';

/**
 * Read environment variables from file.
 * https://github.com/motdotla/dotenv
 */
// import dotenv from 'dotenv';
// import path from 'path';
// dotenv.config({ path: path.resolve(__dirname, '.env') });

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
  testDir: './e2e',
  /* Run tests in files in parallel */
  fullyParallel: true,
  /* Fail the build on CI if you accidentally left test.only in the source code. */
  forbidOnly: !!process.env.CI,
  /* Retry on CI only */
  retries: process.env.CI ? 2 : 0,
  /* Opt out of parallel tests on CI. */
  workers: process.env.CI ? 1 : undefined,
  /* Reporter to use. See https://playwright.dev/docs/test-reporters */
  reporter: [
    ['html'],
    ['json', { outputFile: 'test-results.json' }],
    ['junit', { outputFile: 'junit.xml' }],
    ['list'],
  ],
  /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
  use: {
    /* Base URL to use in actions like `await page.goto('/')`. */
    baseURL: process.env.VITE_APP_URL || 'http://localhost:3000',

    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'on-first-retry',

    /* Screenshot on failure */
    screenshot: 'only-on-failure',

    /* Video on failure */
    video: 'retain-on-failure',

    /* Viewport size */
    viewport: { width: 1280, height: 720 },

    /* Ignore HTTPS errors during navigation */
    ignoreHTTPSErrors: true,

    /* Default timeout for actions */
    actionTimeout: 10000,

    /* Navigation timeout */
    navigationTimeout: 30000,

    /* Test ID attribute for element selection */
    testIdAttribute: 'data-testid',
  },

  /* Configure projects for major browsers */
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
      testDir: './e2e/web',
    },

    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
      testDir: './e2e/web',
    },

    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
      testDir: './e2e/web',
    },

    /* Test against mobile viewports. */
    {
      name: 'Mobile Chrome',
      use: {
        ...devices['Pixel 5'],
        // Commerce-specific mobile settings
        contextOptions: {
          permissions: ['geolocation'],
        },
      },
      testDir: './e2e/web',
    },
    {
      name: 'Mobile Safari',
      use: {
        ...devices['iPhone 12'],
        // Commerce-specific mobile settings
        contextOptions: {
          permissions: ['geolocation'],
        },
      },
      testDir: './e2e/web',
    },

    /* Test against branded browsers. */
    {
      name: 'Microsoft Edge',
      use: { ...devices['Desktop Edge'], channel: 'msedge' },
      testDir: './e2e/web',
    },
    {
      name: 'Google Chrome',
      use: { ...devices['Desktop Chrome'], channel: 'chrome' },
      testDir: './e2e/web',
    },

    /* Cross-platform tests */
    {
      name: 'cross-platform',
      use: { ...devices['Desktop Chrome'] },
      testDir: './e2e/cross-platform',
    },

    /* Performance tests */
    {
      name: 'performance',
      use: {
        ...devices['Desktop Chrome'],
        // Disable video/screenshots for performance tests
        video: 'off',
        screenshot: 'off',
      },
      testDir: './e2e/performance',
    },
  ],

  /* Run your local dev server before starting the tests */
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
    timeout: 120 * 1000,
  },

  /* Global timeout */
  timeout: 60000,

  /* Expect timeout */
  expect: {
    timeout: 5000,
  },

  /* Output folder for test artifacts */
  outputDir: 'test-results/',
});