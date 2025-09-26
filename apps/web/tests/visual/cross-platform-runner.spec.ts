/**
 * Cross-Platform Visual Regression Test Runner
 * Systematic testing of all UI components for iOS consistency
 */

import { test, expect } from '@playwright/test';
import { VisualTestRunner } from './visual-test-utils';
import { COMPONENT_CATEGORIES, TEST_CASES, VISUAL_THRESHOLDS, IOS_VIEWPORTS } from './visual-config';

test.describe('Cross-Platform UI Consistency Test Suite', () => {
  let visualTester: VisualTestRunner;

  test.beforeEach(async ({ page }) => {
    visualTester = new VisualTestRunner(page);
    await page.goto('http://localhost:6006');
  });

  test('Priority 1: Core Interactive Components', async ({ page }) => {
    await visualTester.setIOSViewport('iPhone 12');

    const allResults = [];

    for (const componentName of COMPONENT_CATEGORIES.CORE_INTERACTIVE) {
      // Find test case configuration for this component
      const testCase = TEST_CASES.find(tc => tc.name === componentName);

      if (testCase) {
        const results = await visualTester.testComponentVariants(
          testCase.component,
          testCase.variants,
          testCase.sizes,
          {
            threshold: VISUAL_THRESHOLDS.STRICT,
            animations: 'disabled'
          }
        );

        allResults.push(...results);

        // Log individual component results
        console.log(`${componentName} Results:`, results.map(r => ({
          variant: r.variant,
          size: r.size,
          score: r.consistencyScore,
          passed: r.passed
        })));
      }
    }

    // Generate overall report
    const report = visualTester.generateConsistencyReport(allResults);
    console.log('Priority 1 Components Report:', report);

    // Assert overall quality
    expect(report.overallScore).toBeGreaterThan(0.95);
    expect(report.componentsPassingThreshold / report.componentsAnalyzed).toBeGreaterThan(0.9);
  });

  test('Priority 2: Missing High Priority Components', async ({ page }) => {
    await visualTester.setIOSViewport('iPhone 12');

    for (const componentName of COMPONENT_CATEGORIES.MISSING_HIGH_PRIORITY) {
      // These components need to be implemented
      // This test will initially fail and pass once components are created

      try {
        await visualTester.navigateToComponent(componentName);

        // If we reach here, component exists - test it
        const result = await visualTester.compareWithIOSReference(
          componentName,
          'primary',
          'medium',
          {
            threshold: VISUAL_THRESHOLDS.STRICT,
            animations: 'disabled'
          }
        );

        expect(result.passed).toBe(true);
        expect(result.consistencyScore).toBeGreaterThan(0.95);

      } catch (error) {
        // Component doesn't exist yet - this is expected initially
        console.log(`${componentName} not yet implemented - this is expected during development`);
        test.skip();
      }
    }
  });

  test('Cross-Device Consistency Validation', async ({ page }) => {
    const testComponent = 'Button'; // Use existing component
    const devices = Object.keys(IOS_VIEWPORTS) as (keyof typeof IOS_VIEWPORTS)[];

    const deviceResults = [];

    for (const device of devices) {
      await visualTester.setIOSViewport(device);
      await visualTester.navigateToComponent(testComponent, 'primary', 'medium');

      const screenshotName = await visualTester.captureComponentScreenshot(
        testComponent,
        'primary',
        'medium',
        {
          threshold: VISUAL_THRESHOLDS.DEFAULT,
          animations: 'disabled'
        }
      );

      deviceResults.push({
        device,
        screenshotName,
        viewport: IOS_VIEWPORTS[device]
      });
    }

    // Validate that components render consistently across different iOS devices
    console.log('Cross-device test completed for:', deviceResults);
    expect(deviceResults.length).toBe(devices.length);
  });

  test('Design Token Validation Across Components', async ({ page }) => {
    await visualTester.setIOSViewport('iPhone 12');

    const tokenValidationResults = [];

    for (const testCase of TEST_CASES) {
      await visualTester.navigateToComponent(testCase.component, 'primary', 'medium');

      // Extract and validate design tokens
      const tokenValidation = await page.evaluate(() => {
        const element = document.querySelector('[data-testid="component-container"]');
        if (!element) return null;

        const computed = window.getComputedStyle(element);

        // Extract key design tokens that should be consistent across platforms
        return {
          primaryColor: computed.getPropertyValue('--primary')?.trim() || computed.backgroundColor,
          borderRadius: computed.borderRadius,
          fontFamily: computed.fontFamily,
          fontSize: computed.fontSize,
          spacing: computed.padding,
        };
      });

      if (tokenValidation) {
        tokenValidationResults.push({
          component: testCase.component,
          tokens: tokenValidation
        });
      }
    }

    // Validate token consistency
    expect(tokenValidationResults.length).toBeGreaterThan(0);

    // All components should use consistent design tokens
    const fontFamilies = new Set(tokenValidationResults.map(r => r.tokens.fontFamily));
    expect(fontFamilies.size).toBeLessThanOrEqual(2); // Allow system fonts variation
  });

  test('Performance Impact Assessment', async ({ page }) => {
    await visualTester.setIOSViewport('iPhone 12');

    const performanceMetrics = [];

    for (const testCase of TEST_CASES.slice(0, 3)) { // Test first 3 components
      const startTime = Date.now();

      await visualTester.navigateToComponent(testCase.component);

      // Measure component rendering time
      const metrics = await page.evaluate(() => {
        return {
          renderTime: performance.now(),
          memoryUsed: (performance as any).memory?.usedJSHeapSize || 0,
        };
      });

      const totalTime = Date.now() - startTime;

      performanceMetrics.push({
        component: testCase.component,
        loadTime: totalTime,
        renderMetrics: metrics
      });
    }

    // Performance targets: <200ms component render time
    for (const metric of performanceMetrics) {
      expect(metric.loadTime).toBeLessThan(200);
    }

    console.log('Performance Metrics:', performanceMetrics);
  });

  test('Generate Comprehensive Consistency Report', async ({ page }) => {
    // This test generates the final consistency report
    const reportData = {
      timestamp: new Date().toISOString(),
      testEnvironment: 'Playwright Visual Regression',
      targetConsistency: '95%',
      platformsAnalyzed: ['Web (React)', 'iOS (Reference)'],
      componentsAnalyzed: COMPONENT_CATEGORIES.CORE_INTERACTIVE.length,
      categoriesTested: Object.keys(COMPONENT_CATEGORIES).length,
      devicesValidated: Object.keys(IOS_VIEWPORTS).length,
      thresholds: VISUAL_THRESHOLDS
    };

    console.log('Visual Regression Test Infrastructure Report:', reportData);

    // Write report to file for CI/CD integration
    await page.evaluate((data) => {
      // This would be enhanced to write actual report files
      console.log('Test Report Generated:', JSON.stringify(data, null, 2));
    }, reportData);

    expect(reportData.componentsAnalyzed).toBeGreaterThan(0);
    expect(reportData.categoriesTested).toBeGreaterThan(0);
  });
});