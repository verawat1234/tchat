/**
 * List command - display components from registry
 */

import chalk from 'chalk';
import { ComponentRegistry } from '../../src/core/ComponentRegistry';
import { ComponentType } from '../../src/core/Component';

export interface ListOptions {
  category?: string;
  sort?: string;
  limit?: number;
  format?: string;
  deprecated?: boolean;
}

export class ListCommand {
  async execute(type?: string, options: ListOptions = {}): Promise<void> {
    // Load registry
    const registry = await this.loadRegistry();

    // Get components based on type
    let components = this.getComponentsByType(registry, type);

    // Filter by category if specified
    if (options.category) {
      components = components.filter(c => c.category === options.category);
    }

    // Include/exclude deprecated
    if (!options.deprecated) {
      components = components.filter(c => !c.deprecated);
    }

    // Sort components
    components = this.sortComponents(components, options.sort);

    // Limit results
    if (options.limit) {
      components = components.slice(0, options.limit);
    }

    // Display results
    this.displayResults(components, options.format || 'table');
  }

  private async loadRegistry(): Promise<ComponentRegistry> {
    // Load from default location or config
    return ComponentRegistry.loadFromFile('docs/components/registry.json');
  }

  private getComponentsByType(registry: ComponentRegistry, type?: string): any[] {
    switch (type) {
      case 'atoms':
        return registry.getComponentsByType(ComponentType.ATOM);
      case 'molecules':
        return registry.getComponentsByType(ComponentType.MOLECULE);
      case 'organisms':
        return registry.getComponentsByType(ComponentType.ORGANISM);
      default:
        return registry.getAllComponents();
    }
  }

  private sortComponents(components: any[], sortBy?: string): any[] {
    switch (sortBy) {
      case 'usage':
        return components.sort((a, b) => b.usageCount - a.usageCount);
      case 'created':
        return components.sort((a, b) =>
          new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
        );
      case 'name':
      default:
        return components.sort((a, b) => a.name.localeCompare(b.name));
    }
  }

  private displayResults(components: any[], format: string): void {
    if (components.length === 0) {
      console.log(chalk.yellow('No components found'));
      return;
    }

    switch (format) {
      case 'json':
        console.log(JSON.stringify(components, null, 2));
        break;

      case 'list':
        components.forEach(c => {
          console.log(`${c.name} (${c.type})`);
        });
        break;

      case 'table':
      default:
        console.table(
          components.map(c => ({
            ID: c.id,
            Name: c.name,
            Type: c.type,
            Usage: c.usageCount,
            Status: c.deprecated ? 'deprecated' : 'active'
          }))
        );
        break;
    }

    console.log(chalk.gray(`\nTotal: ${components.length} components`));
  }
}

// Stub imports
class ComponentRegistry {
  static async loadFromFile(path: string): Promise<ComponentRegistry> {
    return new ComponentRegistry({ id: 'test', projectName: 'test' });
  }

  constructor(options: any) {}

  getComponentsByType(type: any): any[] {
    return [];
  }

  getAllComponents(): any[] {
    return [];
  }
}