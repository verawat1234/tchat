/**
 * Atom entity - the most basic UI building block
 */

import { Component, ComponentType, ComponentOptions } from './Component';

export interface AccessibilityInfo {
  ariaLabel: boolean;
  ariaDescribedBy: boolean;
  role: string | null;
  keyboardNav: boolean;
  wcagLevel: 'A' | 'AA' | 'AAA';
}

export interface AtomOptions extends ComponentOptions {
  htmlElement?: string;
  variants?: string[];
  accessibility?: AccessibilityInfo;
}

export class Atom extends Component {
  htmlElement: string;
  variants: string[];
  accessibility: AccessibilityInfo;

  constructor(options: AtomOptions) {
    // Validate component type if provided
    if (options.type && options.type !== ComponentType.ATOM) {
      throw new Error('Invalid component type for Atom');
    }

    super({
      ...options,
      type: ComponentType.ATOM
    });

    this.htmlElement = options.htmlElement || 'div';
    this.variants = options.variants || [];
    this.accessibility = options.accessibility || {
      ariaLabel: false,
      ariaDescribedBy: false,
      role: null,
      keyboardNav: false,
      wcagLevel: 'AA'
    };
  }

  /**
   * Serialize to JSON
   */
  toJSON(): any {
    const baseJson = super.toJSON();
    return {
      ...baseJson,
      htmlElement: this.htmlElement,
      variants: this.variants,
      accessibility: this.accessibility
    };
  }

  /**
   * Deserialize from JSON
   */
  static fromJSON(json: any): Atom {
    return new Atom({
      ...json,
      createdAt: new Date(json.createdAt),
      updatedAt: new Date(json.updatedAt)
    });
  }
}