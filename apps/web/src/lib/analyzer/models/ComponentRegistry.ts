/**
 * ComponentRegistry - Central catalog of all components
 */

import { Component, ComponentType } from './Component';
import { EventEmitter } from 'events';
import * as fs from 'fs/promises';
import * as path from 'path';

export interface RegistryStats {
  totalComponents: number;
  atomCount: number;
  moleculeCount: number;
  organismCount: number;
  duplicatesFound: number;
  inconsistenciesFound: number;
  averageUsageCount: number;
  mostUsedComponents: string[];
}

export interface RegistryOptions {
  id: string;
  projectName: string;
  version?: string;
}

export interface ValidationResult {
  isValid: boolean;
  errors: string[];
  warnings: string[];
}

export class ComponentRegistry extends EventEmitter {
  id: string;
  projectName: string;
  private components: Map<string, Component>;
  lastUpdated: Date;
  statistics: RegistryStats;
  version: string;

  constructor(options: RegistryOptions) {
    super();
    this.id = options.id;
    this.projectName = options.projectName;
    this.version = options.version || '1.0.0';
    this.components = new Map();
    this.lastUpdated = new Date();
    this.statistics = this.calculateStatistics();
  }

  /**
   * Add a component to the registry
   */
  addComponent(component: Component): void {
    if (this.components.has(component.id)) {
      throw new Error(`Component with ID ${component.id} already exists`);
    }
    this.components.set(component.id, component);
    this.lastUpdated = new Date();
    this.statistics = this.calculateStatistics();
    this.emit('component:added', component);
  }

  /**
   * Add multiple components at once
   */
  addComponents(components: Component[]): void {
    components.forEach(component => this.addComponent(component));
  }

  /**
   * Get a component by ID
   */
  getComponent(id: string): Component | undefined {
    return this.components.get(id);
  }

  /**
   * Update an existing component
   */
  updateComponent(component: Component): void {
    if (!this.components.has(component.id)) {
      throw new Error(`Component with ID ${component.id} does not exist`);
    }
    this.components.set(component.id, component);
    this.lastUpdated = new Date();
    this.statistics = this.calculateStatistics();
    this.emit('component:updated', component);
  }

  /**
   * Remove a component from the registry
   */
  removeComponent(id: string): void {
    this.components.delete(id);
    this.lastUpdated = new Date();
    this.statistics = this.calculateStatistics();
    this.emit('component:removed', id);
  }

  /**
   * Check if a component exists
   */
  hasComponent(id: string): boolean {
    return this.components.has(id);
  }

  /**
   * Get total component count
   */
  getComponentCount(): number {
    return this.components.size;
  }

  /**
   * Get all components
   */
  getAllComponents(): Component[] {
    return Array.from(this.components.values());
  }

  /**
   * Get components by type
   */
  getComponentsByType(type: ComponentType): Component[] {
    return Array.from(this.components.values())
      .filter(c => c.type === type);
  }

  /**
   * Get components by category
   */
  getComponentsByCategory(category: string): Component[] {
    return Array.from(this.components.values())
      .filter(c => c.category === category);
  }

  /**
   * Get deprecated components
   */
  getDeprecatedComponents(): Component[] {
    return Array.from(this.components.values())
      .filter(c => c.deprecated);
  }

  /**
   * Get most used components
   */
  getMostUsedComponents(limit: number = 10): Component[] {
    return Array.from(this.components.values())
      .sort((a, b) => b.usageCount - a.usageCount)
      .slice(0, limit);
  }

  /**
   * Search components by name
   */
  searchByName(query: string): Component[] {
    const lowerQuery = query.toLowerCase();
    return Array.from(this.components.values())
      .filter(c => c.name.toLowerCase().includes(lowerQuery));
  }

  /**
   * Search components by file path
   */
  searchByFilePath(query: string): Component[] {
    return Array.from(this.components.values())
      .filter(c => c.filePath.includes(query));
  }

  /**
   * Get registry statistics
   */
  getStatistics(): RegistryStats {
    return this.statistics;
  }

  /**
   * Set duplicates count (set by duplicate detection)
   */
  setDuplicatesCount(count: number): void {
    this.statistics.duplicatesFound = count;
  }

  /**
   * Set inconsistencies count (set by validation)
   */
  setInconsistenciesCount(count: number): void {
    this.statistics.inconsistenciesFound = count;
  }

  /**
   * Calculate statistics
   */
  private calculateStatistics(): RegistryStats {
    const components = this.getAllComponents();
    const totalComponents = components.length;

    if (totalComponents === 0) {
      return {
        totalComponents: 0,
        atomCount: 0,
        moleculeCount: 0,
        organismCount: 0,
        duplicatesFound: 0,
        inconsistenciesFound: 0,
        averageUsageCount: 0,
        mostUsedComponents: []
      };
    }

    const atomCount = components.filter(c => c.type === ComponentType.ATOM).length;
    const moleculeCount = components.filter(c => c.type === ComponentType.MOLECULE).length;
    const organismCount = components.filter(c => c.type === ComponentType.ORGANISM).length;

    const totalUsage = components.reduce((sum, c) => sum + c.usageCount, 0);
    const averageUsageCount = totalUsage / totalComponents;

    const mostUsedComponents = components
      .sort((a, b) => b.usageCount - a.usageCount)
      .slice(0, Math.min(10, totalComponents))
      .map(c => c.id);

    return {
      totalComponents,
      atomCount,
      moleculeCount,
      organismCount,
      duplicatesFound: this.statistics?.duplicatesFound || 0,
      inconsistenciesFound: this.statistics?.inconsistenciesFound || 0,
      averageUsageCount,
      mostUsedComponents
    };
  }

  /**
   * Clear all components
   */
  clear(): void {
    this.components.clear();
    this.lastUpdated = new Date();
    this.statistics = this.calculateStatistics();
    this.emit('registry:cleared');
  }

  /**
   * Validate registry integrity
   */
  validate(): ValidationResult {
    const errors: string[] = [];
    const warnings: string[] = [];

    // Check for orphaned references
    this.getAllComponents().forEach(component => {
      component.dependencies.forEach(dep => {
        if (dep.startsWith('nonexistent')) {
          errors.push(`Orphaned reference: ${dep}`);
        }
      });
    });

    return {
      isValid: errors.length === 0,
      errors,
      warnings
    };
  }

  /**
   * Serialize to JSON
   */
  toJSON(): any {
    return {
      id: this.id,
      projectName: this.projectName,
      version: this.version,
      components: this.getAllComponents().map(c => c.toJSON()),
      lastUpdated: this.lastUpdated.toISOString(),
      statistics: this.statistics
    };
  }

  /**
   * Create from JSON
   */
  static fromJSON(json: any): ComponentRegistry {
    const registry = new ComponentRegistry({
      id: json.id,
      projectName: json.projectName,
      version: json.version
    });

    if (json.components) {
      json.components.forEach((compJson: any) => {
        const component = Component.fromJSON(compJson);
        registry.components.set(component.id, component);
      });
    }

    if (json.lastUpdated) {
      registry.lastUpdated = new Date(json.lastUpdated);
    }

    registry.statistics = registry.calculateStatistics();

    if (json.statistics) {
      registry.statistics.duplicatesFound = json.statistics.duplicatesFound || 0;
      registry.statistics.inconsistenciesFound = json.statistics.inconsistenciesFound || 0;
    }

    return registry;
  }

  /**
   * Save registry to file
   */
  async saveToFile(filePath: string): Promise<void> {
    const json = this.toJSON();
    const dir = path.dirname(filePath);

    // Ensure directory exists
    await fs.mkdir(dir, { recursive: true });

    // Write file
    await fs.writeFile(filePath, JSON.stringify(json, null, 2), 'utf8');
  }

  /**
   * Load registry from file
   */
  static async loadFromFile(filePath: string): Promise<ComponentRegistry> {
    const content = await fs.readFile(filePath, 'utf8');
    const json = JSON.parse(content);
    return ComponentRegistry.fromJSON(json);
  }
}