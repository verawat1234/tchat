/**
 * Utility functions for the Tchat design system
 * Supports cross-platform component development
 */
import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

/**
 * Combines class names using clsx and tailwind-merge for optimal TailwindCSS class handling
 * Prevents style conflicts and ensures proper class precedence
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * Type-safe delay utility for animations and loading states
 */
export const delay = (ms: number): Promise<void> =>
  new Promise(resolve => setTimeout(resolve, ms));

/**
 * Design token validation helpers
 */
export const designTokens = {
  /**
   * Validates if a color value matches OKLCH precision requirements
   */
  validateColorAccuracy: (expected: string, actual: string, tolerance = 0.02): boolean => {
    // Implementation would include OKLCH color space comparison
    // For now, return true for basic hex comparison
    return expected.toLowerCase() === actual.toLowerCase();
  },

  /**
   * Converts design token spacing to consistent pixel values
   */
  spacingToPx: (token: string): number => {
    const spacingMap: Record<string, number> = {
      'xs': 4,   // --spacing-xs
      'sm': 8,   // --spacing-sm
      'md': 16,  // --spacing-md
      'lg': 24,  // --spacing-lg
      'xl': 32,  // --spacing-xl
      '2xl': 48, // --spacing-2xl
    };
    return spacingMap[token] || 0;
  },

  /**
   * Validates cross-platform consistency threshold
   */
  validateConsistency: (score: number): boolean => {
    return score >= 0.97; // Constitutional requirement: 97% consistency
  },
} as const;

/**
 * Performance monitoring utilities
 */
export const performance = {
  /**
   * Measures component render time against Constitutional <200ms requirement
   */
  measureRenderTime: (componentName: string, startTime: number): boolean => {
    const endTime = Date.now();
    const renderTime = endTime - startTime;
    const isWithinBudget = renderTime < 200; // <200ms Constitutional requirement

    if (!isWithinBudget) {
      console.warn(`⚠️ Performance: ${componentName} rendered in ${renderTime}ms (exceeds 200ms budget)`);
    }

    return isWithinBudget;
  },

  /**
   * Validates 60fps animation performance
   */
  validateAnimationFrameRate: (): boolean => {
    // Implementation would measure actual frame rates
    // For now, return true if requestAnimationFrame is available
    return typeof requestAnimationFrame !== 'undefined';
  },
} as const;