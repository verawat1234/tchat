/**
 * Tests for UsagePattern entity
 */

import { describe, it, expect, beforeEach } from 'vitest';
import { UsagePattern } from '../UsagePattern';

describe('UsagePattern', () => {
  let pattern: UsagePattern;

  beforeEach(() => {
    pattern = new UsagePattern({
      id: 'pattern-button-click',
      componentId: 'atom-button-primary',
      usageCount: 25,
      locations: [
        '/src/pages/HomePage.tsx',
        '/src/components/forms/LoginForm.tsx',
        '/src/components/modals/ConfirmModal.tsx'
      ],
      propPatterns: [
        { prop: 'variant', value: 'primary', frequency: 0.8 },
        { prop: 'size', value: 'medium', frequency: 0.6 },
        { prop: 'disabled', value: false, frequency: 0.95 }
      ],
      commonCombinations: [
        {
          components: ['Button', 'Icon'],
          frequency: 0.4,
          description: 'Icon button pattern'
        },
        {
          components: ['Button', 'Form'],
          frequency: 0.7,
          description: 'Form submission pattern'
        }
      ]
    });
  });

  describe('constructor', () => {
    it('should create a usage pattern with required properties', () => {
      expect(pattern).toBeInstanceOf(UsagePattern);
      expect(pattern.id).toBe('pattern-button-click');
      expect(pattern.componentId).toBe('atom-button-primary');
      expect(pattern.usageCount).toBe(25);
    });

    it('should initialize with default values when not provided', () => {
      const minimalPattern = new UsagePattern({
        id: 'minimal-pattern',
        componentId: 'test-component',
        usageCount: 1
      });

      expect(minimalPattern.locations).toEqual([]);
      expect(minimalPattern.propPatterns).toEqual([]);
      expect(minimalPattern.commonCombinations).toEqual([]);
    });

    it('should validate required fields', () => {
      expect(() => {
        new UsagePattern({
          id: '',
          componentId: 'test',
          usageCount: 0
        });
      }).toThrow('Pattern ID is required');

      expect(() => {
        new UsagePattern({
          id: 'test',
          componentId: '',
          usageCount: 0
        });
      }).toThrow('Component ID is required');
    });
  });

  describe('usageCount property', () => {
    it('should track usage frequency', () => {
      expect(pattern.usageCount).toBe(25);

      pattern.incrementUsage();
      expect(pattern.usageCount).toBe(26);

      pattern.incrementUsage();
      pattern.incrementUsage();
      expect(pattern.usageCount).toBe(28);
    });

    it('should not allow negative usage count', () => {
      expect(() => {
        pattern.usageCount = -1;
      }).toThrow('Usage count cannot be negative');
    });
  });

  describe('locations property', () => {
    it('should track component usage locations', () => {
      expect(pattern.locations).toHaveLength(3);
      expect(pattern.locations).toContain('/src/pages/HomePage.tsx');
    });

    it('should add new location', () => {
      pattern.addLocation('/src/pages/AboutPage.tsx');
      expect(pattern.locations).toHaveLength(4);
      expect(pattern.locations).toContain('/src/pages/AboutPage.tsx');
    });

    it('should not add duplicate locations', () => {
      pattern.addLocation('/src/pages/HomePage.tsx');
      expect(pattern.locations).toHaveLength(3);
    });

    it('should clear locations', () => {
      pattern.clearLocations();
      expect(pattern.locations).toEqual([]);
    });
  });

  describe('propPatterns property', () => {
    it('should track prop usage patterns', () => {
      expect(pattern.propPatterns).toHaveLength(3);

      const variantPattern = pattern.propPatterns.find(p => p.prop === 'variant');
      expect(variantPattern).toBeDefined();
      expect(variantPattern?.value).toBe('primary');
      expect(variantPattern?.frequency).toBe(0.8);
    });

    it('should add prop pattern', () => {
      pattern.addPropPattern('color', 'blue', 0.3);

      const colorPattern = pattern.propPatterns.find(p => p.prop === 'color');
      expect(colorPattern).toBeDefined();
      expect(colorPattern?.value).toBe('blue');
      expect(colorPattern?.frequency).toBe(0.3);
    });

    it('should update existing prop pattern', () => {
      pattern.addPropPattern('variant', 'secondary', 0.2);

      const variantPatterns = pattern.propPatterns.filter(p => p.prop === 'variant');
      expect(variantPatterns).toHaveLength(2);
    });

    it('should validate frequency range', () => {
      expect(() => {
        pattern.addPropPattern('test', 'value', 1.5);
      }).toThrow('Frequency must be between 0 and 1');

      expect(() => {
        pattern.addPropPattern('test', 'value', -0.1);
      }).toThrow('Frequency must be between 0 and 1');
    });
  });

  describe('commonCombinations property', () => {
    it('should track common component combinations', () => {
      expect(pattern.commonCombinations).toHaveLength(2);

      const iconCombo = pattern.commonCombinations.find(c => c.description === 'Icon button pattern');
      expect(iconCombo).toBeDefined();
      expect(iconCombo?.components).toContain('Button');
      expect(iconCombo?.components).toContain('Icon');
      expect(iconCombo?.frequency).toBe(0.4);
    });

    it('should add combination pattern', () => {
      pattern.addCombination(
        ['Button', 'Tooltip'],
        0.3,
        'Tooltip button pattern'
      );

      expect(pattern.commonCombinations).toHaveLength(3);
      const tooltipCombo = pattern.commonCombinations.find(c => c.description === 'Tooltip button pattern');
      expect(tooltipCombo).toBeDefined();
    });

    it('should require at least two components in combination', () => {
      expect(() => {
        pattern.addCombination(['Button'], 0.5, 'Single component');
      }).toThrow('Combination must have at least 2 components');
    });
  });

  describe('analysis methods', () => {
    it('should get most common prop values', () => {
      const mostCommon = pattern.getMostCommonPropValue('variant');
      expect(mostCommon).toBe('primary');
    });

    it('should return null for unknown prop', () => {
      const unknown = pattern.getMostCommonPropValue('unknown');
      expect(unknown).toBeNull();
    });

    it('should calculate average prop frequency', () => {
      const avgFreq = pattern.getAveragePropFrequency();
      // (0.8 + 0.6 + 0.95) / 3 = 0.783...
      expect(avgFreq).toBeCloseTo(0.783, 2);
    });

    it('should return 0 for no prop patterns', () => {
      pattern.propPatterns = [];
      const avgFreq = pattern.getAveragePropFrequency();
      expect(avgFreq).toBe(0);
    });

    it('should find most frequent combination', () => {
      const mostFrequent = pattern.getMostFrequentCombination();
      expect(mostFrequent).toBeDefined();
      expect(mostFrequent?.description).toBe('Form submission pattern');
      expect(mostFrequent?.frequency).toBe(0.7);
    });
  });

  describe('toJSON', () => {
    it('should serialize pattern to JSON', () => {
      const json = pattern.toJSON();

      expect(json).toMatchObject({
        id: 'pattern-button-click',
        componentId: 'atom-button-primary',
        usageCount: 25,
        locations: expect.arrayContaining(['/src/pages/HomePage.tsx']),
        propPatterns: expect.arrayContaining([
          expect.objectContaining({ prop: 'variant', value: 'primary' })
        ]),
        commonCombinations: expect.arrayContaining([
          expect.objectContaining({ description: 'Icon button pattern' })
        ])
      });
    });
  });

  describe('fromJSON static method', () => {
    it('should deserialize pattern from JSON', () => {
      const json = pattern.toJSON();
      const deserializedPattern = UsagePattern.fromJSON(json);

      expect(deserializedPattern).toBeInstanceOf(UsagePattern);
      expect(deserializedPattern.id).toBe(pattern.id);
      expect(deserializedPattern.componentId).toBe(pattern.componentId);
      expect(deserializedPattern.usageCount).toBe(pattern.usageCount);
      expect(deserializedPattern.locations).toEqual(pattern.locations);
      expect(deserializedPattern.propPatterns).toEqual(pattern.propPatterns);
      expect(deserializedPattern.commonCombinations).toEqual(pattern.commonCombinations);
    });

    it('should handle dates correctly in deserialization', () => {
      const json = pattern.toJSON();
      const deserializedPattern = UsagePattern.fromJSON(json);

      expect(deserializedPattern.lastUpdated).toBeInstanceOf(Date);
    });
  });

  describe('pattern insights', () => {
    it('should identify high-frequency patterns', () => {
      const highFreqProps = pattern.propPatterns.filter(p => p.frequency > 0.7);
      expect(highFreqProps).toHaveLength(2); // disabled: 0.95, variant: 0.8
    });

    it('should track pattern evolution', () => {
      const initialCount = pattern.usageCount;
      const initialLocations = pattern.locations.length;

      // Simulate usage growth
      for (let i = 0; i < 5; i++) {
        pattern.incrementUsage();
        pattern.addLocation(`/src/pages/Page${i}.tsx`);
      }

      expect(pattern.usageCount).toBe(initialCount + 5);
      expect(pattern.locations.length).toBeGreaterThan(initialLocations);
    });
  });
});