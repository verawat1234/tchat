/**
 * Organism entity - complex, self-contained component sections
 */

import { Component, ComponentType, ComponentOptions } from './Component';

export interface OrganismOptions extends ComponentOptions {
  contains?: string[];
  standalone?: boolean;
  pageSection?: boolean;
  dataSource?: string | null;
}

export class Organism extends Component {
  contains: string[];
  standalone: boolean;
  pageSection: boolean;
  dataSource: string | null;

  constructor(options: OrganismOptions) {
    // Validate component type if provided
    if (options.type && options.type !== ComponentType.ORGANISM) {
      throw new Error('Invalid component type for Organism');
    }

    super({
      ...options,
      type: ComponentType.ORGANISM
    });

    this.contains = options.contains || [];
    this.standalone = options.standalone || false;
    this.pageSection = options.pageSection || false;
    this.dataSource = options.dataSource || null;
  }

  /**
   * Serialize to JSON
   */
  toJSON(): any {
    const baseJson = super.toJSON();
    return {
      ...baseJson,
      contains: this.contains,
      standalone: this.standalone,
      pageSection: this.pageSection,
      dataSource: this.dataSource
    };
  }

  /**
   * Deserialize from JSON
   */
  static fromJSON(json: any): Organism {
    return new Organism({
      ...json,
      createdAt: new Date(json.createdAt),
      updatedAt: new Date(json.updatedAt)
    });
  }
}