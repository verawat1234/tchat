/**
 * Global Setup for Visual Regression Testing
 * Prepares environment for cross-platform visual consistency validation
 * Constitutional requirement: 97% visual consistency across platforms
 */
import { chromium, FullConfig } from '@playwright/test';

async function globalSetup(config: FullConfig) {
  console.log('🎨 Setting up Visual Regression Testing Environment');

  // Launch browser for setup tasks
  const browser = await chromium.launch();
  const page = await browser.newPage();

  try {
    // Wait for Storybook to be ready
    console.log('⏳ Waiting for Storybook server...');
    await page.goto('http://localhost:6006', { waitUntil: 'networkidle' });
    await page.waitForSelector('[data-testid="storybook-preview"]', { timeout: 60000 });

    // Verify design tokens are loaded
    console.log('🎯 Verifying design tokens are loaded...');
    const tokensLoaded = await page.evaluate(() => {
      const rootStyles = getComputedStyle(document.documentElement);
      return (
        rootStyles.getPropertyValue('--color-primary').trim() !== '' &&
        rootStyles.getPropertyValue('--spacing-md').trim() !== '' &&
        rootStyles.getPropertyValue('--text-base').trim() !== ''
      );
    });

    if (!tokensLoaded) {
      throw new Error('Design tokens not loaded properly');
    }

    // Verify component stories are available
    console.log('🧩 Verifying component stories...');
    const requiredStories = [
      'components-tchatbutton--variants',
      'components-tchatinput--types',
      'components-tchatcard--variants'
    ];

    for (const story of requiredStories) {
      try {
        await page.goto(`http://localhost:6006/?path=/story/${story}`);
        await page.waitForSelector('[data-testid]', { timeout: 10000 });
      } catch (error) {
        console.warn(`⚠️  Story not available: ${story} (expected for TDD)`);
        // This is expected in TDD - stories will exist after implementation
      }
    }

    // Verify cross-platform reference data
    console.log('📱 Verifying cross-platform reference setup...');

    // Create reference directories if they don't exist
    const fs = await import('fs').then(m => m.promises);
    const path = await import('path');

    const referenceDir = path.join(process.cwd(), 'test-results', 'visual-references');
    const platformDirs = ['web', 'ios', 'android'];

    for (const platform of platformDirs) {
      const platformDir = path.join(referenceDir, platform);
      try {
        await fs.mkdir(platformDir, { recursive: true });
      } catch (error) {
        // Directory might already exist
      }
    }

    // Set up test environment variables
    process.env.VISUAL_TESTING = 'true';
    process.env.CONSISTENCY_THRESHOLD = '0.97'; // Constitutional requirement
    process.env.RENDER_TIMEOUT = '200'; // <200ms requirement

    console.log('✅ Visual Regression Testing Environment Ready');
    console.log('🎯 Constitutional Requirements:');
    console.log('   - 97% visual consistency across platforms');
    console.log('   - <200ms component load times');
    console.log('   - WCAG 2.1 AA accessibility compliance');
    console.log('   - 60fps animation performance');

  } catch (error) {
    console.error('❌ Visual testing setup failed:', error);
    throw error;
  } finally {
    await browser.close();
  }
}

export default globalSetup;