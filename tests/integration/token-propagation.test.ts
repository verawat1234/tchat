import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { promises as fs } from 'fs';
import path from 'path';
import { DesignTokenPropagator, type PropagationConfig, type TokenPropagationResult } from '../../tools/design-tokens/src/services/DesignTokenPropagator';
import type { DesignToken } from '../../tools/design-tokens/src/types';

/**
 * Integration Test: Design Token Propagation
 *
 * CRITICAL TDD: This test MUST FAIL initially because:
 * 1. DesignTokenPropagator service doesn't exist yet
 * 2. Token propagation types don't exist yet
 * 3. File generation and synchronization logic doesn't exist yet
 *
 * These tests drive the implementation of design token propagation
 * from the central system to web CSS and iOS Swift files.
 */

const TEST_OUTPUT_DIR = path.join(__dirname, '../../../test-output/design-tokens');
const WEB_OUTPUT_PATH = path.join(TEST_OUTPUT_DIR, 'web/tokens.css');
const IOS_OUTPUT_PATH = path.join(TEST_OUTPUT_DIR, 'ios/Sources/DesignSystem/Colors.swift');

describe('Design Token Propagation Integration Tests', () => {
  let propagator: DesignTokenPropagator;
  let config: PropagationConfig;

  beforeEach(async () => {
    // EXPECTED TO FAIL: DesignTokenPropagator class doesn't exist
    propagator = new DesignTokenPropagator();

    // EXPECTED TO FAIL: PropagationConfig type doesn't exist
    config = {
      outputPaths: {
        web: WEB_OUTPUT_PATH,
        ios: IOS_OUTPUT_PATH
      },
      templates: {
        web: {
          header: '/* Auto-generated design tokens - DO NOT EDIT MANUALLY */\\n:root {\\n',
          footer: '}\\n',
          tokenFormat: '  --{name}: {value};\\n'
        },
        ios: {
          header: '// Auto-generated design tokens - DO NOT EDIT MANUALLY\\nimport SwiftUI\\n\\nstruct Colors {\\n',
          footer: '}\\n',
          tokenFormat: '    static let {name} = Color(hex: "{value}")\\n'
        }
      },
      validation: {
        enableConsistencyCheck: true,
        maxColorDelta: 0.05,
        requiredPlatforms: ['web', 'ios']
      }
    };

    // Clean up test output directory
    await fs.rm(TEST_OUTPUT_DIR, { recursive: true, force: true });
    await fs.mkdir(TEST_OUTPUT_DIR, { recursive: true });
  });

  afterEach(async () => {
    // Clean up after each test
    await fs.rm(TEST_OUTPUT_DIR, { recursive: true, force: true });
  });

  describe('Color Token Propagation', () => {
    it('should propagate color tokens to both platforms', async () => {
      // EXPECTED TO FAIL: propagateTokens method doesn't exist
      const colorTokens: DesignToken[] = [
        {
          id: '1',
          name: 'primary-500',
          category: 'color',
          value: {
            oklch: { l: 0.6338, c: 0.2078, h: 252.57 }
          },
          platforms: ['web', 'ios'],
          generatedValues: {
            web: {
              css: '--color-primary-500: oklch(63.38% 0.2078 252.57);',
              hex: '#3b82f6'
            },
            ios: {
              swift: 'static let primary500 = Color(hex: "#3b82f6")',
              uiColor: 'UIColor(red: 0.231, green: 0.510, blue: 0.965, alpha: 1.0)'
            }
          },
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z'
        },
        {
          id: '2',
          name: 'success-500',
          category: 'color',
          value: {
            oklch: { l: 0.6977, c: 0.1686, h: 142.5 }
          },
          platforms: ['web', 'ios'],
          generatedValues: {
            web: {
              css: '--color-success-500: oklch(69.77% 0.1686 142.5);',
              hex: '#22c55e'
            },
            ios: {
              swift: 'static let success500 = Color(hex: "#22c55e")',
              uiColor: 'UIColor(red: 0.133, green: 0.773, blue: 0.369, alpha: 1.0)'
            }
          },
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z'
        }
      ];

      const result: TokenPropagationResult = await propagator.propagateTokens(colorTokens, config);

      // Verify propagation success
      expect(result.success).toBe(true);
      expect(result.generatedFiles).toContain(WEB_OUTPUT_PATH);
      expect(result.generatedFiles).toContain(IOS_OUTPUT_PATH);
      expect(result.errors).toHaveLength(0);

      // Verify web CSS file
      const webContent = await fs.readFile(WEB_OUTPUT_PATH, 'utf-8');
      expect(webContent).toContain('--color-primary-500: oklch(63.38% 0.2078 252.57);');
      expect(webContent).toContain('--color-success-500: oklch(69.77% 0.1686 142.5);');
      expect(webContent).toMatch(/^\/\* Auto-generated design tokens/);
      expect(webContent).toContain(':root {');

      // Verify iOS Swift file
      const iosContent = await fs.readFile(IOS_OUTPUT_PATH, 'utf-8');
      expect(iosContent).toContain('static let primary500 = Color(hex: "#3b82f6")');
      expect(iosContent).toContain('static let success500 = Color(hex: "#22c55e")');
      expect(iosContent).toMatch(/^\/\/ Auto-generated design tokens/);
      expect(iosContent).toContain('struct Colors {');
    });

    it('should handle spacing token propagation', async () => {
      const spacingTokens: DesignToken[] = [
        {
          id: '3',
          name: 'spacing-xs',
          category: 'spacing',
          value: { base: 8 },
          platforms: ['web', 'ios'],
          generatedValues: {
            web: { css: '--spacing-xs: 8px;' },
            ios: { swift: 'static let spacingXs: CGFloat = 8' }
          },
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z'
        },
        {
          id: '4',
          name: 'spacing-md',
          category: 'spacing',
          value: { base: 16 },
          platforms: ['web', 'ios'],
          generatedValues: {
            web: { css: '--spacing-md: 16px;' },
            ios: { swift: 'static let spacingMd: CGFloat = 16' }
          },
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z'
        }
      ];

      // Use spacing-specific configuration
      const spacingConfig: PropagationConfig = {
        ...config,
        templates: {
          web: {
            header: '/* Auto-generated spacing tokens */\\n:root {\\n',
            footer: '}\\n',
            tokenFormat: '  --{name}: {value};\\n'
          },
          ios: {
            header: '// Auto-generated spacing tokens\\nimport SwiftUI\\n\\nstruct Spacing {\\n',
            footer: '}\\n',
            tokenFormat: '    {swiftCode}\\n'
          }
        }
      };

      const result = await propagator.propagateTokens(spacingTokens, spacingConfig);

      expect(result.success).toBe(true);

      const webContent = await fs.readFile(WEB_OUTPUT_PATH, 'utf-8');
      expect(webContent).toContain('--spacing-xs: 8px;');
      expect(webContent).toContain('--spacing-md: 16px;');

      const iosContent = await fs.readFile(IOS_OUTPUT_PATH, 'utf-8');
      expect(iosContent).toContain('static let spacingXs: CGFloat = 8');
      expect(iosContent).toContain('static let spacingMd: CGFloat = 16');
      expect(iosContent).toContain('struct Spacing {');
    });
  });

  describe('Synchronization and Consistency', () => {
    it('should detect and report inconsistencies', async () => {
      // EXPECTED TO FAIL: Consistency checking doesn't exist
      const inconsistentTokens: DesignToken[] = [
        {
          id: '5',
          name: 'inconsistent-color',
          category: 'color',
          value: {
            oklch: { l: 0.5, c: 0.2, h: 250 }
          },
          platforms: ['web', 'ios'],
          generatedValues: {
            web: {
              css: '--color-inconsistent: oklch(50% 0.2 250);',
              hex: '#8080ff'  // Correct value
            },
            ios: {
              swift: 'static let inconsistent = Color(hex: "#ff0000")',  // Wrong value!
              uiColor: 'UIColor(red: 1.0, green: 0.0, blue: 0.0, alpha: 1.0)'
            }
          },
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z'
        }
      ];

      const result = await propagator.propagateTokens(inconsistentTokens, config);

      // Should still generate files but report inconsistencies
      expect(result.success).toBe(false);
      expect(result.errors).toHaveLength(1);
      expect(result.errors[0]).toContain('Color inconsistency detected');
      expect(result.errors[0]).toContain('inconsistent-color');
      expect(result.consistencyReport?.inconsistentTokens).toHaveLength(1);
    });

    it('should update existing files incrementally', async () => {
      // First propagation
      const initialTokens: DesignToken[] = [
        {
          id: '6',
          name: 'initial-color',
          category: 'color',
          value: { oklch: { l: 0.5, c: 0.2, h: 120 } },
          platforms: ['web', 'ios'],
          generatedValues: {
            web: { css: '--color-initial: oklch(50% 0.2 120);', hex: '#00ff80' },
            ios: { swift: 'static let initial = Color(hex: "#00ff80")' }
          },
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z'
        }
      ];

      await propagator.propagateTokens(initialTokens, config);

      // Verify initial files
      let webContent = await fs.readFile(WEB_OUTPUT_PATH, 'utf-8');
      let iosContent = await fs.readFile(IOS_OUTPUT_PATH, 'utf-8');

      expect(webContent).toContain('--color-initial');
      expect(iosContent).toContain('static let initial');

      // Second propagation with additional tokens
      const updatedTokens: DesignToken[] = [
        ...initialTokens,
        {
          id: '7',
          name: 'additional-color',
          category: 'color',
          value: { oklch: { l: 0.7, c: 0.15, h: 200 } },
          platforms: ['web', 'ios'],
          generatedValues: {
            web: { css: '--color-additional: oklch(70% 0.15 200);', hex: '#4dccff' },
            ios: { swift: 'static let additional = Color(hex: "#4dccff")' }
          },
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z'
        }
      ];

      const result = await propagator.propagateTokens(updatedTokens, config);

      expect(result.success).toBe(true);

      // Verify updated files contain both tokens
      webContent = await fs.readFile(WEB_OUTPUT_PATH, 'utf-8');
      iosContent = await fs.readFile(IOS_OUTPUT_PATH, 'utf-8');

      expect(webContent).toContain('--color-initial');
      expect(webContent).toContain('--color-additional');
      expect(iosContent).toContain('static let initial');
      expect(iosContent).toContain('static let additional');
    });
  });

  describe('Error Handling', () => {
    it('should handle file system errors gracefully', async () => {
      // EXPECTED TO FAIL: Error handling doesn't exist
      const invalidConfig: PropagationConfig = {
        ...config,
        outputPaths: {
          web: '/invalid/path/that/does/not/exist/tokens.css',
          ios: '/another/invalid/path/Colors.swift'
        }
      };

      const tokens: DesignToken[] = [{
        id: '8',
        name: 'test-color',
        category: 'color',
        value: { oklch: { l: 0.5, c: 0.2, h: 250 } },
        platforms: ['web', 'ios'],
        generatedValues: {
          web: { css: '--color-test: oklch(50% 0.2 250);' },
          ios: { swift: 'static let test = Color(hex: "#8080ff")' }
        },
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z'
      }];

      const result = await propagator.propagateTokens(tokens, invalidConfig);

      expect(result.success).toBe(false);
      expect(result.errors.length).toBeGreaterThan(0);
      expect(result.errors.some(error => error.includes('Could not write to'))).toBe(true);
    });

    it('should validate token data before propagation', async () => {
      // Invalid token data
      const invalidTokens: any[] = [
        {
          id: '9',
          name: '', // Invalid: empty name
          category: 'color',
          value: { oklch: { l: 0.5, c: 0.2, h: 250 } },
          platforms: ['web', 'ios']
        },
        {
          id: '10',
          name: 'valid-name',
          category: 'invalid-category', // Invalid: unknown category
          value: { oklch: { l: 0.5, c: 0.2, h: 250 } },
          platforms: ['web', 'ios']
        }
      ];

      const result = await propagator.propagateTokens(invalidTokens, config);

      expect(result.success).toBe(false);
      expect(result.errors.some(error => error.includes('Invalid token name'))).toBe(true);
      expect(result.errors.some(error => error.includes('Invalid token category'))).toBe(true);
    });
  });

  describe('Performance and Scalability', () => {
    it('should handle large token sets efficiently', async () => {
      // Generate 100+ tokens to test performance
      const largeTokenSet: DesignToken[] = Array.from({ length: 100 }, (_, i) => ({
        id: `perf-${i}`,
        name: `performance-color-${i}`,
        category: 'color',
        value: {
          oklch: {
            l: 0.3 + (i * 0.003), // Varying lightness
            c: 0.1 + (i * 0.001), // Varying chroma
            h: i * 3.6            // Full hue rotation
          }
        },
        platforms: ['web', 'ios'],
        generatedValues: {
          web: {
            css: `--color-performance-${i}: oklch(${30 + i * 0.3}% ${0.1 + i * 0.001} ${i * 3.6});`,
            hex: `#${((i * 13) % 256).toString(16).padStart(2, '0')}8080`
          },
          ios: {
            swift: `static let performance${i} = Color(hex: "#${((i * 13) % 256).toString(16).padStart(2, '0')}8080")`
          }
        },
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z'
      }));

      const startTime = Date.now();
      const result = await propagator.propagateTokens(largeTokenSet, config);
      const duration = Date.now() - startTime;

      expect(result.success).toBe(true);
      expect(duration).toBeLessThan(5000); // Should complete within 5 seconds

      // Verify all tokens were processed
      const webContent = await fs.readFile(WEB_OUTPUT_PATH, 'utf-8');
      const iosContent = await fs.readFile(IOS_OUTPUT_PATH, 'utf-8');

      expect((webContent.match(/--color-performance-/g) || []).length).toBe(100);
      expect((iosContent.match(/static let performance/g) || []).length).toBe(100);
    });
  });
});