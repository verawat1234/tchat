/**
 * ConsistencyRule - standards that components must follow
 */

import { ComponentType } from './Component';

export enum RuleCategory {
  NAMING = 'naming',
  STRUCTURE = 'structure',
  STYLING = 'styling',
  ACCESSIBILITY = 'accessibility',
  PERFORMANCE = 'performance',
  DOCUMENTATION = 'documentation'
}

export enum RuleSeverity {
  ERROR = 'error',
  WARNING = 'warning',
  INFO = 'info'
}

export interface ConsistencyRuleOptions {
  id: string;
  name: string;
  description: string;
  category: RuleCategory;
  severity: RuleSeverity;
  validator: string;
  appliesTo?: ComponentType[];
  enabled?: boolean;
  autoFixable?: boolean;
  examples?: {
    good: string[];
    bad: string[];
  };
  documentation?: string;
  createdAt?: Date;
  updatedAt?: Date;
}

export class ConsistencyRule {
  id: string;
  name: string;
  description: string;
  category: RuleCategory;
  severity: RuleSeverity;
  validator: string;
  appliesTo: ComponentType[];
  enabled: boolean;
  autoFixable: boolean;
  examples?: {
    good: string[];
    bad: string[];
  };
  documentation?: string;
  createdAt: Date;
  updatedAt: Date;

  constructor(options: ConsistencyRuleOptions) {
    // Validate required fields
    if (!options.id) {
      throw new Error('Rule ID is required');
    }
    if (!options.name) {
      throw new Error('Rule name is required');
    }
    if (!options.validator) {
      throw new Error('Validator function name is required');
    }

    this.id = options.id;
    this.name = options.name;
    this.description = options.description;
    this.category = options.category;
    this.severity = options.severity;
    this.validator = options.validator;
    this.appliesTo = options.appliesTo || [];
    this.enabled = options.enabled !== false;
    this.autoFixable = options.autoFixable || false;
    this.examples = options.examples;
    this.documentation = options.documentation;
    this.createdAt = options.createdAt || new Date();
    this.updatedAt = options.updatedAt || new Date();
  }

  /**
   * Enable the rule
   */
  enable(): void {
    this.enabled = true;
    this.updatedAt = new Date();
  }

  /**
   * Disable the rule
   */
  disable(): void {
    this.enabled = false;
    this.updatedAt = new Date();
  }

  /**
   * Check if rule applies to a specific component type
   */
  appliesToType(type: ComponentType): boolean {
    if (this.appliesTo.length === 0) {
      return false; // If no specific types specified, doesn't apply
    }
    return this.appliesTo.includes(type);
  }

  /**
   * Serialize to JSON
   */
  toJSON(): any {
    return {
      id: this.id,
      name: this.name,
      description: this.description,
      category: this.category,
      severity: this.severity,
      validator: this.validator,
      appliesTo: this.appliesTo,
      enabled: this.enabled,
      autoFixable: this.autoFixable,
      examples: this.examples,
      documentation: this.documentation,
      createdAt: this.createdAt.toISOString(),
      updatedAt: this.updatedAt.toISOString()
    };
  }

  /**
   * Deserialize from JSON
   */
  static fromJSON(json: any): ConsistencyRule {
    return new ConsistencyRule({
      ...json,
      createdAt: new Date(json.createdAt),
      updatedAt: new Date(json.updatedAt)
    });
  }
}