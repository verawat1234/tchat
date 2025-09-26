import { describe, it, expect } from 'vitest';
import { ColorConverter, type OklchColor, type RgbColor, type HexColor } from '../../tools/design-tokens/src/services/ColorConverter';

/**
 * Integration Test: Color Token Conversion (OKLCH ↔ Hex)
 *
 * CRITICAL TDD: This test MUST FAIL initially because:
 * 1. ColorConverter service doesn't exist yet
 * 2. Color type interfaces don't exist yet
 * 3. Mathematical conversion algorithms don't exist yet
 *
 * These tests drive the implementation of mathematically accurate color conversion
 * required for design token consistency between web (OKLCH) and iOS (Hex/RGB).
 */

describe('Color Conversion Integration Tests', () => {
  let colorConverter: ColorConverter;

  beforeEach(() => {
    // EXPECTED TO FAIL: ColorConverter class doesn't exist
    colorConverter = new ColorConverter();
  });

  describe('OKLCH to Hex Conversion', () => {
    it('should convert Tailwind primary-500 (blue) accurately', () => {
      // EXPECTED TO FAIL: oklchToHex method doesn't exist
      const oklchBlue: OklchColor = {
        l: 0.6338, // 63.38% lightness
        c: 0.2078, // 20.78% chroma
        h: 252.57  // 252.57° hue (blue)
      };

      const hexResult = colorConverter.oklchToHex(oklchBlue);

      // Expected: Tailwind's blue-500 #3b82f6
      expect(hexResult).toBe('#3b82f6');
    });

    it('should convert Tailwind green-500 accurately', () => {
      const oklchGreen: OklchColor = {
        l: 0.6977, // 69.77% lightness
        c: 0.1686, // 16.86% chroma
        h: 142.5   // 142.5° hue (green)
      };

      const hexResult = colorConverter.oklchToHex(oklchGreen);

      // Expected: Tailwind's green-500 #22c55e
      expect(hexResult).toBe('#22c55e');
    });

    it('should convert Tailwind red-500 accurately', () => {
      const oklchRed: OklchColor = {
        l: 0.6274, // 62.74% lightness
        c: 0.2583, // 25.83% chroma
        h: 27.33   // 27.33° hue (red)
      };

      const hexResult = colorConverter.oklchToHex(oklchRed);

      // Expected: Tailwind's red-500 #ef4444
      expect(hexResult).toBe('#ef4444');
    });

    it('should handle grayscale colors (zero chroma)', () => {
      const oklchGray: OklchColor = {
        l: 0.5,    // 50% lightness
        c: 0,      // 0% chroma (grayscale)
        h: 0       // Hue irrelevant for grayscale
      };

      const hexResult = colorConverter.oklchToHex(oklchGray);

      // Should be a neutral gray
      expect(hexResult).toMatch(/^#[0-9a-f]{6}$/i);
      // All RGB components should be approximately equal for grayscale
      const r = parseInt(hexResult.slice(1, 3), 16);
      const g = parseInt(hexResult.slice(3, 5), 16);
      const b = parseInt(hexResult.slice(5, 7), 16);

      expect(Math.abs(r - g)).toBeLessThan(5);
      expect(Math.abs(g - b)).toBeLessThan(5);
      expect(Math.abs(r - b)).toBeLessThan(5);
    });

    it('should handle pure white and black', () => {
      // Pure white
      const oklchWhite: OklchColor = { l: 1, c: 0, h: 0 };
      expect(colorConverter.oklchToHex(oklchWhite)).toBe('#ffffff');

      // Pure black
      const oklchBlack: OklchColor = { l: 0, c: 0, h: 0 };
      expect(colorConverter.oklchToHex(oklchBlack)).toBe('#000000');
    });
  });

  describe('Hex to OKLCH Conversion (Reverse)', () => {
    it('should convert hex colors back to OKLCH accurately', () => {
      // EXPECTED TO FAIL: hexToOklch method doesn't exist
      const hexBlue = '#3b82f6';

      const oklchResult = colorConverter.hexToOklch(hexBlue);

      // Should be approximately the original values (within tolerance)
      expect(oklchResult.l).toBeCloseTo(0.6338, 2);
      expect(oklchResult.c).toBeCloseTo(0.2078, 2);
      expect(oklchResult.h).toBeCloseTo(252.57, 1);
    });

    it('should handle round-trip conversion with minimal loss', () => {
      const originalOklch: OklchColor = {
        l: 0.6338,
        c: 0.2078,
        h: 252.57
      };

      // Convert to hex and back
      const hex = colorConverter.oklchToHex(originalOklch);
      const roundTripOklch = colorConverter.hexToOklch(hex);

      // Should be very close to original (accounting for rounding)
      expect(roundTripOklch.l).toBeCloseTo(originalOklch.l, 2);
      expect(roundTripOklch.c).toBeCloseTo(originalOklch.c, 2);
      expect(Math.abs(roundTripOklch.h - originalOklch.h)).toBeLessThan(2);
    });
  });

  describe('RGB Intermediate Conversion', () => {
    it('should provide accurate RGB values for iOS UIColor', () => {
      // EXPECTED TO FAIL: oklchToRgb method doesn't exist
      const oklchBlue: OklchColor = {
        l: 0.6338,
        c: 0.2078,
        h: 252.57
      };

      const rgbResult = colorConverter.oklchToRgb(oklchBlue);

      // Expected RGB values for #3b82f6
      expect(rgbResult.r).toBeCloseTo(59, 2);   // 0x3b = 59
      expect(rgbResult.g).toBeCloseTo(130, 2);  // 0x82 = 130
      expect(rgbResult.b).toBeCloseTo(246, 2);  // 0xf6 = 246
    });

    it('should provide normalized RGB values for iOS Color', () => {
      // EXPECTED TO FAIL: oklchToNormalizedRgb method doesn't exist
      const oklchBlue: OklchColor = {
        l: 0.6338,
        c: 0.2078,
        h: 252.57
      };

      const normalizedRgb = colorConverter.oklchToNormalizedRgb(oklchBlue);

      // Values should be between 0 and 1
      expect(normalizedRgb.r).toBeGreaterThanOrEqual(0);
      expect(normalizedRgb.r).toBeLessThanOrEqual(1);
      expect(normalizedRgb.g).toBeGreaterThanOrEqual(0);
      expect(normalizedRgb.g).toBeLessThanOrEqual(1);
      expect(normalizedRgb.b).toBeGreaterThanOrEqual(0);
      expect(normalizedRgb.b).toBeLessThanOrEqual(1);

      // Should match expected normalized values for #3b82f6
      expect(normalizedRgb.r).toBeCloseTo(59 / 255, 3);
      expect(normalizedRgb.g).toBeCloseTo(130 / 255, 3);
      expect(normalizedRgb.b).toBeCloseTo(246 / 255, 3);
    });
  });

  describe('Error Handling', () => {
    it('should throw error for invalid OKLCH values', () => {
      // EXPECTED TO FAIL: Validation doesn't exist
      const invalidOklch: OklchColor = {
        l: -0.1,  // Invalid: negative lightness
        c: 0.2,
        h: 250
      };

      expect(() => {
        colorConverter.oklchToHex(invalidOklch);
      }).toThrow('Invalid OKLCH lightness value: must be between 0 and 1');

      const invalidChroma: OklchColor = {
        l: 0.6,
        c: -0.1,  // Invalid: negative chroma
        h: 250
      };

      expect(() => {
        colorConverter.oklchToHex(invalidChroma);
      }).toThrow('Invalid OKLCH chroma value: must be non-negative');
    });

    it('should throw error for invalid hex colors', () => {
      // Invalid hex formats
      expect(() => {
        colorConverter.hexToOklch('#invalid');
      }).toThrow('Invalid hex color format');

      expect(() => {
        colorConverter.hexToOklch('not-hex');
      }).toThrow('Invalid hex color format');

      expect(() => {
        colorConverter.hexToOklch('#12345');  // Too short
      }).toThrow('Invalid hex color format');
    });

    it('should handle out-of-gamut colors gracefully', () => {
      // EXPECTED TO FAIL: Gamut clamping doesn't exist
      const outOfGamut: OklchColor = {
        l: 0.5,
        c: 1.0,    // Very high chroma, might be out of sRGB gamut
        h: 120
      };

      // Should not throw, but clamp to valid sRGB range
      expect(() => {
        const hex = colorConverter.oklchToHex(outOfGamut);
        expect(hex).toMatch(/^#[0-9a-f]{6}$/i);
      }).not.toThrow();
    });
  });

  describe('Precision and Accuracy', () => {
    it('should maintain precision for design system critical colors', () => {
      // Test the exact colors we use in our design system
      const designSystemColors = [
        { name: 'primary-500', oklch: { l: 0.6338, c: 0.2078, h: 252.57 }, expectedHex: '#3b82f6' },
        { name: 'success-500', oklch: { l: 0.6977, c: 0.1686, h: 142.5 }, expectedHex: '#22c55e' },
        { name: 'error-500', oklch: { l: 0.6274, c: 0.2583, h: 27.33 }, expectedHex: '#ef4444' },
        { name: 'warning-500', oklch: { l: 0.8047, c: 0.1543, h: 85.87 }, expectedHex: '#f59e0b' }
      ];

      designSystemColors.forEach(({ name, oklch, expectedHex }) => {
        const convertedHex = colorConverter.oklchToHex(oklch);
        expect(convertedHex).toBe(expectedHex, `${name} conversion should be exact`);
      });
    });
  });
});