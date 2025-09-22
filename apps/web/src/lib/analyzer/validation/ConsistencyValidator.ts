/**
 * ConsistencyValidator - Validates components against consistency rules
 */

import { Component, ComponentType } from '../models/Component';
import { ConsistencyRule, RuleCategory, RuleSeverity } from '../models/ConsistencyRule';
import { ComponentRegistry } from '../models/ComponentRegistry';

export interface ValidationResult {
  componentId: string;
  valid: boolean;
  violations: Violation[];
  warnings: Warning[];
  suggestions: string[];
}

export interface Violation {
  ruleId: string;
  ruleName: string;
  category: RuleCategory;
  severity: RuleSeverity;
  message: string;
  fix?: string;
}

export interface Warning {
  message: string;
  suggestion?: string;
}

export interface ValidationOptions {
  rules?: string[];
  severity?: RuleSeverity;
  autoFix?: boolean;
}

export class ConsistencyValidator {
  private rules: ConsistencyRule[];

  constructor() {
    this.rules = this.getDefaultRules();
  }

  /**
   * Validate a single component
   */
  validateComponent(
    component: Component,
    options?: ValidationOptions
  ): ValidationResult {
    const violations: Violation[] = [];
    const warnings: Warning[] = [];
    const suggestions: string[] = [];

    // Filter rules based on options
    const rulesToApply = this.filterRules(component.type, options);

    // Apply each rule
    for (const rule of rulesToApply) {
      const result = this.applyRule(component, rule);

      if (result.violation) {
        violations.push(result.violation);
      }

      if (result.warning) {
        warnings.push(result.warning);
      }

      if (result.suggestion) {
        suggestions.push(result.suggestion);
      }
    }

    // Auto-fix if requested
    if (options?.autoFix) {
      this.applyAutoFixes(component, violations);
    }

    return {
      componentId: component.id,
      valid: violations.length === 0,
      violations,
      warnings,
      suggestions
    };
  }

  /**
   * Validate all components in registry
   */
  validateRegistry(
    registry: ComponentRegistry,
    options?: ValidationOptions
  ): {
    results: ValidationResult[];
    summary: ValidationSummary;
  } {
    const results: ValidationResult[] = [];
    const components = registry.getAllComponents();

    for (const component of components) {
      const result = this.validateComponent(component, options);
      results.push(result);
    }

    const summary = this.generateSummary(results);

    return {
      results,
      summary
    };
  }

  /**
   * Apply a single rule to a component
   */
  private applyRule(
    component: Component,
    rule: ConsistencyRule
  ): {
    violation?: Violation;
    warning?: Warning;
    suggestion?: string;
  } {
    const result: any = {};

    switch (rule.category) {
      case RuleCategory.NAMING:
        this.validateNaming(component, rule, result);
        break;

      case RuleCategory.STRUCTURE:
        this.validateStructure(component, rule, result);
        break;

      case RuleCategory.STYLING:
        this.validateStyling(component, rule, result);
        break;

      case RuleCategory.ACCESSIBILITY:
        this.validateAccessibility(component, rule, result);
        break;

      case RuleCategory.PERFORMANCE:
        this.validatePerformance(component, rule, result);
        break;

      case RuleCategory.DOCUMENTATION:
        this.validateDocumentation(component, rule, result);
        break;
    }

    return result;
  }

  /**
   * Validate naming conventions
   */
  private validateNaming(
    component: Component,
    rule: ConsistencyRule,
    result: any
  ): void {
    const pattern = /^[A-Z][a-zA-Z]*$/; // PascalCase

    if (!pattern.test(component.name)) {
      result.violation = {
        ruleId: rule.id,
        ruleName: rule.name,
        category: rule.category,
        severity: rule.severity,
        message: `Component name "${component.name}" does not follow PascalCase convention`,
        fix: this.toPascalCase(component.name)
      };
    }

    // Check for meaningful names
    if (component.name.length < 3) {
      result.warning = {
        message: `Component name "${component.name}" is too short`,
        suggestion: 'Use more descriptive component names'
      };
    }

    // Check for generic names
    const genericNames = ['Component', 'Container', 'Wrapper', 'Box'];
    if (genericNames.includes(component.name)) {
      result.warning = {
        message: `Component name "${component.name}" is too generic`,
        suggestion: 'Use a more specific name that describes the component\'s purpose'
      };
    }
  }

  /**
   * Validate component structure
   */
  private validateStructure(
    component: Component,
    rule: ConsistencyRule,
    result: any
  ): void {
    // Check for props validation
    if (component.props.length > 10) {
      result.violation = {
        ruleId: rule.id,
        ruleName: rule.name,
        category: rule.category,
        severity: RuleSeverity.WARNING,
        message: `Component has ${component.props.length} props, which is too many`,
        fix: 'Consider breaking down the component or using composition'
      };
    }

    // Check for required props
    const requiredProps = component.props.filter(p => p.required);
    if (requiredProps.length > 5) {
      result.warning = {
        message: `Component has ${requiredProps.length} required props`,
        suggestion: 'Consider using default values or optional props where appropriate'
      };
    }
  }

  /**
   * Validate styling consistency
   */
  private validateStyling(
    component: Component,
    rule: ConsistencyRule,
    result: any
  ): void {
    // Check for consistent styling approach
    const hasStyledComponents = component.dependencies.some(d =>
      d.includes('styled-components') || d.includes('@emotion')
    );
    const hasCSSModules = component.filePath.includes('.module.');
    const hasTailwind = component.dependencies.some(d => d.includes('tailwind'));

    const stylingApproaches = [hasStyledComponents, hasCSSModules, hasTailwind].filter(Boolean);

    if (stylingApproaches.length > 1) {
      result.warning = {
        message: 'Component uses multiple styling approaches',
        suggestion: 'Stick to a single styling methodology for consistency'
      };
    }
  }

  /**
   * Validate accessibility
   */
  private validateAccessibility(
    component: Component,
    rule: ConsistencyRule,
    result: any
  ): void {
    // Check for accessibility props
    const hasAriaProps = component.props.some(p =>
      p.name.startsWith('aria') || p.name === 'role'
    );

    if (component.type === ComponentType.ATOM && !hasAriaProps) {
      result.warning = {
        message: 'Atom component lacks accessibility props',
        suggestion: 'Consider adding ARIA attributes for better accessibility'
      };
    }

    // Check for onClick without keyboard support
    const hasOnClick = component.props.some(p => p.name === 'onClick');
    const hasOnKeyDown = component.props.some(p =>
      p.name === 'onKeyDown' || p.name === 'onKeyPress'
    );

    if (hasOnClick && !hasOnKeyDown) {
      result.violation = {
        ruleId: rule.id,
        ruleName: rule.name,
        category: rule.category,
        severity: RuleSeverity.WARNING,
        message: 'Component has onClick without keyboard event handler',
        fix: 'Add onKeyDown or onKeyPress handler for keyboard accessibility'
      };
    }
  }

  /**
   * Validate performance considerations
   */
  private validatePerformance(
    component: Component,
    rule: ConsistencyRule,
    result: any
  ): void {
    // Check for too many dependencies
    if (component.dependencies.length > 10) {
      result.violation = {
        ruleId: rule.id,
        ruleName: rule.name,
        category: rule.category,
        severity: RuleSeverity.WARNING,
        message: `Component has ${component.dependencies.length} dependencies`,
        fix: 'Consider reducing dependencies or lazy loading'
      };
    }

    // Check for heavy components based on usage
    if (component.type === ComponentType.ORGANISM && component.usageCount > 20) {
      result.suggestion = 'Consider optimizing this heavily used organism component';
    }
  }

  /**
   * Validate documentation
   */
  private validateDocumentation(
    component: Component,
    rule: ConsistencyRule,
    result: any
  ): void {
    // Check for description
    if (!component.description || component.description.trim() === '') {
      result.violation = {
        ruleId: rule.id,
        ruleName: rule.name,
        category: rule.category,
        severity: RuleSeverity.INFO,
        message: 'Component lacks description',
        fix: 'Add a meaningful description explaining the component\'s purpose'
      };
    }

    // Check for prop documentation
    const undocumentedProps = component.props.filter(p =>
      !p.description || p.description.trim() === ''
    );

    if (undocumentedProps.length > 0) {
      result.warning = {
        message: `${undocumentedProps.length} props lack documentation`,
        suggestion: `Document props: ${undocumentedProps.map(p => p.name).join(', ')}`
      };
    }

    // Check for examples
    const propsWithoutExamples = component.props.filter(p =>
      !p.examples || p.examples.length === 0
    );

    if (propsWithoutExamples.length > 0 && component.type !== ComponentType.ATOM) {
      result.suggestion = 'Consider adding usage examples for complex props';
    }
  }

  /**
   * Filter rules based on component type and options
   */
  private filterRules(
    componentType: ComponentType,
    options?: ValidationOptions
  ): ConsistencyRule[] {
    let filtered = this.rules.filter(rule =>
      rule.enabled && rule.appliesTo.includes(componentType)
    );

    // Filter by specific rules if provided
    if (options?.rules && options.rules.length > 0) {
      filtered = filtered.filter(rule =>
        options.rules!.includes(rule.id) || options.rules!.includes(rule.category)
      );
    }

    // Filter by severity
    if (options?.severity) {
      const severityOrder = {
        [RuleSeverity.ERROR]: 0,
        [RuleSeverity.WARNING]: 1,
        [RuleSeverity.INFO]: 2
      };

      const minSeverity = severityOrder[options.severity];
      filtered = filtered.filter(rule =>
        severityOrder[rule.severity] <= minSeverity
      );
    }

    return filtered;
  }

  /**
   * Apply auto-fixes to component
   */
  private applyAutoFixes(component: Component, violations: Violation[]): void {
    for (const violation of violations) {
      if (violation.fix) {
        switch (violation.category) {
          case RuleCategory.NAMING:
            if (violation.fix) {
              component.name = violation.fix;
            }
            break;
          // Add more auto-fix implementations as needed
        }
      }
    }
  }

  /**
   * Convert string to PascalCase
   */
  private toPascalCase(str: string): string {
    return str
      .replace(/[-_\s]+(.)?/g, (_, char) => char ? char.toUpperCase() : '')
      .replace(/^(.)/, (_, char) => char.toUpperCase());
  }

  /**
   * Generate validation summary
   */
  private generateSummary(results: ValidationResult[]): ValidationSummary {
    const totalChecked = results.length;
    const passed = results.filter(r => r.valid).length;
    const failed = totalChecked - passed;

    const allViolations = results.flatMap(r => r.violations);
    const errorCount = allViolations.filter(v => v.severity === RuleSeverity.ERROR).length;
    const warningCount = allViolations.filter(v => v.severity === RuleSeverity.WARNING).length;
    const infoCount = allViolations.filter(v => v.severity === RuleSeverity.INFO).length;

    const violationsByCategory = new Map<RuleCategory, number>();
    for (const violation of allViolations) {
      const count = violationsByCategory.get(violation.category) || 0;
      violationsByCategory.set(violation.category, count + 1);
    }

    return {
      totalChecked,
      passed,
      failed,
      errorCount,
      warningCount,
      infoCount,
      violationsByCategory: Object.fromEntries(violationsByCategory)
    };
  }

  /**
   * Get default validation rules
   */
  private getDefaultRules(): ConsistencyRule[] {
    return [
      new ConsistencyRule({
        id: 'naming-pascal-case',
        name: 'Component Naming Convention',
        description: 'Components should use PascalCase naming',
        category: RuleCategory.NAMING,
        severity: RuleSeverity.ERROR,
        validator: 'validateNaming',
        appliesTo: [ComponentType.ATOM, ComponentType.MOLECULE, ComponentType.ORGANISM],
        enabled: true
      }),
      new ConsistencyRule({
        id: 'structure-prop-limit',
        name: 'Prop Limit',
        description: 'Components should not have too many props',
        category: RuleCategory.STRUCTURE,
        severity: RuleSeverity.WARNING,
        validator: 'validateStructure',
        appliesTo: [ComponentType.ATOM, ComponentType.MOLECULE, ComponentType.ORGANISM],
        enabled: true
      }),
      new ConsistencyRule({
        id: 'accessibility-keyboard',
        name: 'Keyboard Accessibility',
        description: 'Interactive components must support keyboard navigation',
        category: RuleCategory.ACCESSIBILITY,
        severity: RuleSeverity.WARNING,
        validator: 'validateAccessibility',
        appliesTo: [ComponentType.ATOM, ComponentType.MOLECULE],
        enabled: true
      }),
      new ConsistencyRule({
        id: 'documentation-required',
        name: 'Documentation Required',
        description: 'Components must have descriptions',
        category: RuleCategory.DOCUMENTATION,
        severity: RuleSeverity.INFO,
        validator: 'validateDocumentation',
        appliesTo: [ComponentType.ATOM, ComponentType.MOLECULE, ComponentType.ORGANISM],
        enabled: true
      }),
      new ConsistencyRule({
        id: 'performance-dependencies',
        name: 'Dependency Limit',
        description: 'Components should minimize dependencies',
        category: RuleCategory.PERFORMANCE,
        severity: RuleSeverity.WARNING,
        validator: 'validatePerformance',
        appliesTo: [ComponentType.MOLECULE, ComponentType.ORGANISM],
        enabled: true
      })
    ];
  }
}

export interface ValidationSummary {
  totalChecked: number;
  passed: number;
  failed: number;
  errorCount: number;
  warningCount: number;
  infoCount: number;
  violationsByCategory: Record<string, number>;
}