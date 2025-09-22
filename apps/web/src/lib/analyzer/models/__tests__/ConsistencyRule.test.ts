/**
 * Tests for ConsistencyRule entity
 */

import { describe, it, expect, beforeEach } from 'vitest';
import { ConsistencyRule, RuleCategory, RuleSeverity } from '../ConsistencyRule';
import { ComponentType } from '../Component';

describe('ConsistencyRule', () => {
  let rule: ConsistencyRule;

  beforeEach(() => {
    rule = new ConsistencyRule({
      id: 'naming-pascal-case',
      name: 'Component Naming Convention',
      description: 'Components should use PascalCase naming',
      category: RuleCategory.NAMING,
      severity: RuleSeverity.ERROR,
      validator: 'validatePascalCase',
      appliesTo: [ComponentType.ATOM, ComponentType.MOLECULE, ComponentType.ORGANISM],
      enabled: true,
      autoFixable: true,
      examples: {
        good: ['MyComponent', 'PrimaryButton', 'UserProfile'],
        bad: ['myComponent', 'primary-button', 'user_profile']
      },
      documentation: 'https://docs.example.com/naming-conventions'
    });
  });

  describe('constructor', () => {
    it('should create a rule with required properties', () => {
      expect(rule).toBeInstanceOf(ConsistencyRule);
      expect(rule.id).toBe('naming-pascal-case');
      expect(rule.category).toBe(RuleCategory.NAMING);
      expect(rule.severity).toBe(RuleSeverity.ERROR);
      expect(rule.enabled).toBe(true);
    });

    it('should initialize with default values when not provided', () => {
      const minimalRule = new ConsistencyRule({
        id: 'minimal-rule',
        name: 'Minimal Rule',
        description: 'A minimal test rule',
        category: RuleCategory.STRUCTURE,
        severity: RuleSeverity.WARNING,
        validator: 'validateMinimal'
      });

      expect(minimalRule.appliesTo).toEqual([]);
      expect(minimalRule.enabled).toBe(true);
      expect(minimalRule.autoFixable).toBe(false);
      expect(minimalRule.examples).toBeUndefined();
      expect(minimalRule.documentation).toBeUndefined();
    });

    it('should validate required fields', () => {
      expect(() => {
        new ConsistencyRule({
          id: '',
          name: 'Test',
          description: 'Test',
          category: RuleCategory.NAMING,
          severity: RuleSeverity.ERROR,
          validator: 'test'
        });
      }).toThrow('Rule ID is required');

      expect(() => {
        new ConsistencyRule({
          id: 'test',
          name: '',
          description: 'Test',
          category: RuleCategory.NAMING,
          severity: RuleSeverity.ERROR,
          validator: 'test'
        });
      }).toThrow('Rule name is required');

      expect(() => {
        new ConsistencyRule({
          id: 'test',
          name: 'Test',
          description: 'Test',
          category: RuleCategory.NAMING,
          severity: RuleSeverity.ERROR,
          validator: ''
        });
      }).toThrow('Validator function name is required');
    });
  });

  describe('category property', () => {
    it('should categorize rules appropriately', () => {
      expect(rule.category).toBe(RuleCategory.NAMING);

      const structureRule = new ConsistencyRule({
        id: 'structure-prop-limit',
        name: 'Prop Limit Rule',
        description: 'Components should not have too many props',
        category: RuleCategory.STRUCTURE,
        severity: RuleSeverity.WARNING,
        validator: 'validatePropLimit'
      });
      expect(structureRule.category).toBe(RuleCategory.STRUCTURE);

      const accessibilityRule = new ConsistencyRule({
        id: 'a11y-aria-labels',
        name: 'ARIA Labels Required',
        description: 'Interactive elements must have ARIA labels',
        category: RuleCategory.ACCESSIBILITY,
        severity: RuleSeverity.ERROR,
        validator: 'validateAriaLabels'
      });
      expect(accessibilityRule.category).toBe(RuleCategory.ACCESSIBILITY);
    });

    it('should support all rule categories', () => {
      const categories = [
        RuleCategory.NAMING,
        RuleCategory.STRUCTURE,
        RuleCategory.STYLING,
        RuleCategory.ACCESSIBILITY,
        RuleCategory.PERFORMANCE,
        RuleCategory.DOCUMENTATION
      ];

      categories.forEach(category => {
        const testRule = new ConsistencyRule({
          id: `test-${category}`,
          name: `Test ${category}`,
          description: 'Test rule',
          category,
          severity: RuleSeverity.INFO,
          validator: 'validate'
        });
        expect(testRule.category).toBe(category);
      });
    });
  });

  describe('severity property', () => {
    it('should set appropriate severity levels', () => {
      expect(rule.severity).toBe(RuleSeverity.ERROR);

      rule.severity = RuleSeverity.WARNING;
      expect(rule.severity).toBe(RuleSeverity.WARNING);

      rule.severity = RuleSeverity.INFO;
      expect(rule.severity).toBe(RuleSeverity.INFO);
    });

    it('should prioritize rules by severity', () => {
      const errorRule = new ConsistencyRule({
        id: 'error',
        name: 'Error Rule',
        description: 'Critical issue',
        category: RuleCategory.SECURITY,
        severity: RuleSeverity.ERROR,
        validator: 'validate'
      });

      const warningRule = new ConsistencyRule({
        id: 'warning',
        name: 'Warning Rule',
        description: 'Potential issue',
        category: RuleCategory.PERFORMANCE,
        severity: RuleSeverity.WARNING,
        validator: 'validate'
      });

      const infoRule = new ConsistencyRule({
        id: 'info',
        name: 'Info Rule',
        description: 'Suggestion',
        category: RuleCategory.DOCUMENTATION,
        severity: RuleSeverity.INFO,
        validator: 'validate'
      });

      const rules = [infoRule, errorRule, warningRule];
      const sorted = rules.sort((a, b) => {
        const severityOrder = {
          [RuleSeverity.ERROR]: 0,
          [RuleSeverity.WARNING]: 1,
          [RuleSeverity.INFO]: 2
        };
        return severityOrder[a.severity] - severityOrder[b.severity];
      });

      expect(sorted[0].severity).toBe(RuleSeverity.ERROR);
      expect(sorted[1].severity).toBe(RuleSeverity.WARNING);
      expect(sorted[2].severity).toBe(RuleSeverity.INFO);
    });
  });

  describe('appliesTo property', () => {
    it('should specify applicable component types', () => {
      expect(rule.appliesTo).toContain(ComponentType.ATOM);
      expect(rule.appliesTo).toContain(ComponentType.MOLECULE);
      expect(rule.appliesTo).toContain(ComponentType.ORGANISM);
    });

    it('should filter rules by component type', () => {
      const atomOnlyRule = new ConsistencyRule({
        id: 'atom-only',
        name: 'Atom Only Rule',
        description: 'Applies only to atoms',
        category: RuleCategory.STRUCTURE,
        severity: RuleSeverity.WARNING,
        validator: 'validate',
        appliesTo: [ComponentType.ATOM]
      });

      expect(atomOnlyRule.appliesTo).toHaveLength(1);
      expect(atomOnlyRule.appliesTo[0]).toBe(ComponentType.ATOM);
      expect(atomOnlyRule.appliesToType(ComponentType.ATOM)).toBe(true);
      expect(atomOnlyRule.appliesToType(ComponentType.MOLECULE)).toBe(false);
    });

    it('should handle rules that apply to all types', () => {
      expect(rule.appliesToType(ComponentType.ATOM)).toBe(true);
      expect(rule.appliesToType(ComponentType.MOLECULE)).toBe(true);
      expect(rule.appliesToType(ComponentType.ORGANISM)).toBe(true);
    });
  });

  describe('enabled property', () => {
    it('should toggle rule activation', () => {
      expect(rule.enabled).toBe(true);

      rule.disable();
      expect(rule.enabled).toBe(false);

      rule.enable();
      expect(rule.enabled).toBe(true);
    });

    it('should filter active rules', () => {
      const rules = [
        rule,
        new ConsistencyRule({
          id: 'disabled',
          name: 'Disabled Rule',
          description: 'Inactive rule',
          category: RuleCategory.STYLING,
          severity: RuleSeverity.INFO,
          validator: 'validate',
          enabled: false
        })
      ];

      const activeRules = rules.filter(r => r.enabled);
      expect(activeRules).toHaveLength(1);
      expect(activeRules[0].id).toBe('naming-pascal-case');
    });
  });

  describe('autoFixable property', () => {
    it('should indicate if rule can be auto-fixed', () => {
      expect(rule.autoFixable).toBe(true);

      const nonFixableRule = new ConsistencyRule({
        id: 'non-fixable',
        name: 'Complex Rule',
        description: 'Requires manual intervention',
        category: RuleCategory.STRUCTURE,
        severity: RuleSeverity.ERROR,
        validator: 'validate',
        autoFixable: false
      });

      expect(nonFixableRule.autoFixable).toBe(false);
    });

    it('should identify auto-fixable violations', () => {
      const fixableRules = [rule].filter(r => r.autoFixable);
      expect(fixableRules).toHaveLength(1);
    });
  });

  describe('examples property', () => {
    it('should provide good and bad examples', () => {
      expect(rule.examples?.good).toContain('MyComponent');
      expect(rule.examples?.good).toContain('PrimaryButton');
      expect(rule.examples?.bad).toContain('myComponent');
      expect(rule.examples?.bad).toContain('primary-button');
    });

    it('should be optional', () => {
      const ruleWithoutExamples = new ConsistencyRule({
        id: 'no-examples',
        name: 'Rule without examples',
        description: 'No examples provided',
        category: RuleCategory.PERFORMANCE,
        severity: RuleSeverity.INFO,
        validator: 'validate'
      });

      expect(ruleWithoutExamples.examples).toBeUndefined();
    });
  });

  describe('toJSON', () => {
    it('should serialize rule to JSON', () => {
      const json = rule.toJSON();

      expect(json).toMatchObject({
        id: 'naming-pascal-case',
        name: 'Component Naming Convention',
        description: 'Components should use PascalCase naming',
        category: RuleCategory.NAMING,
        severity: RuleSeverity.ERROR,
        validator: 'validatePascalCase',
        appliesTo: expect.arrayContaining([ComponentType.ATOM]),
        enabled: true,
        autoFixable: true
      });
    });

    it('should include optional properties when present', () => {
      const json = rule.toJSON();
      expect(json.examples).toBeDefined();
      expect(json.documentation).toBe('https://docs.example.com/naming-conventions');
    });
  });

  describe('fromJSON static method', () => {
    it('should deserialize rule from JSON', () => {
      const json = rule.toJSON();
      const deserializedRule = ConsistencyRule.fromJSON(json);

      expect(deserializedRule).toBeInstanceOf(ConsistencyRule);
      expect(deserializedRule.id).toBe(rule.id);
      expect(deserializedRule.name).toBe(rule.name);
      expect(deserializedRule.category).toBe(rule.category);
      expect(deserializedRule.severity).toBe(rule.severity);
      expect(deserializedRule.enabled).toBe(rule.enabled);
      expect(deserializedRule.autoFixable).toBe(rule.autoFixable);
    });

    it('should handle optional properties correctly', () => {
      const json = rule.toJSON();
      const deserializedRule = ConsistencyRule.fromJSON(json);

      expect(deserializedRule.examples).toEqual(rule.examples);
      expect(deserializedRule.documentation).toBe(rule.documentation);
    });
  });

  describe('rule validation', () => {
    it('should validate based on validator function name', () => {
      expect(rule.validator).toBe('validatePascalCase');
    });

    it('should support custom validators', () => {
      const customRule = new ConsistencyRule({
        id: 'custom',
        name: 'Custom Validation',
        description: 'Uses custom validator',
        category: RuleCategory.STRUCTURE,
        severity: RuleSeverity.WARNING,
        validator: 'customValidationFunction'
      });

      expect(customRule.validator).toBe('customValidationFunction');
    });

    it('should check if rule applies to component type', () => {
      expect(rule.appliesToType(ComponentType.ATOM)).toBe(true);
      expect(rule.appliesToType(ComponentType.MOLECULE)).toBe(true);
      expect(rule.appliesToType(ComponentType.ORGANISM)).toBe(true);
      expect(rule.appliesToType('INVALID' as ComponentType)).toBe(false);
    });
  });

  describe('rule management', () => {
    it('should support rule sets', () => {
      const ruleSet = [
        rule,
        new ConsistencyRule({
          id: 'prop-types',
          name: 'Prop Types Required',
          description: 'Components must define prop types',
          category: RuleCategory.STRUCTURE,
          severity: RuleSeverity.ERROR,
          validator: 'validatePropTypes'
        }),
        new ConsistencyRule({
          id: 'docs-required',
          name: 'Documentation Required',
          description: 'Components must have JSDoc comments',
          category: RuleCategory.DOCUMENTATION,
          severity: RuleSeverity.WARNING,
          validator: 'validateDocumentation'
        })
      ];

      const namingRules = ruleSet.filter(r => r.category === RuleCategory.NAMING);
      expect(namingRules).toHaveLength(1);

      const errorRules = ruleSet.filter(r => r.severity === RuleSeverity.ERROR);
      expect(errorRules).toHaveLength(2);

      const enabledRules = ruleSet.filter(r => r.enabled);
      expect(enabledRules).toHaveLength(3);
    });
  });
});