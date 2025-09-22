import { describe, it, expect, beforeEach } from 'vitest';
import { ComponentRegistry, RegistryStats } from '../../../src/lib/analyzer/models/ComponentRegistry';
import { Component, ComponentType } from '../../../src/lib/analyzer/models/Component';

describe('ComponentRegistry', () => {
  let registry: ComponentRegistry;
  let testComponent: Component;

  beforeEach(() => {
    registry = new ComponentRegistry({
      id: 'test-registry',
      projectName: 'test-project'
    });

    testComponent = new Component({
      id: 'test-component',
      name: 'TestComponent',
      type: ComponentType.ATOM,
      filePath: 'src/components/TestComponent.tsx',
      category: 'test',
      description: 'Test component',
      props: [],
      dependencies: [],
      usageCount: 0,
      deprecated: false,
      version: '1.0.0'
    });
  });

  describe('Component Management', () => {
    it('should add component to registry', () => {
      registry.addComponent(testComponent);
      expect(registry.getComponent('test-component')).toEqual(testComponent);
      expect(registry.getComponentCount()).toBe(1);
    });

    it('should prevent duplicate component IDs', () => {
      registry.addComponent(testComponent);
      expect(() => {
        registry.addComponent(testComponent);
      }).toThrow('Component with ID test-component already exists');
    });

    it('should update existing component', () => {
      registry.addComponent(testComponent);
      
      const updated = testComponent.clone();
      updated.description = 'Updated description';
      
      registry.updateComponent(updated);
      const retrieved = registry.getComponent('test-component');
      expect(retrieved?.description).toBe('Updated description');
    });

    it('should remove component from registry', () => {
      registry.addComponent(testComponent);
      expect(registry.getComponentCount()).toBe(1);
      
      registry.removeComponent('test-component');
      expect(registry.getComponent('test-component')).toBeUndefined();
      expect(registry.getComponentCount()).toBe(0);
    });

    it('should check if component exists', () => {
      expect(registry.hasComponent('test-component')).toBe(false);
      registry.addComponent(testComponent);
      expect(registry.hasComponent('test-component')).toBe(true);
    });
  });

  describe('Component Queries', () => {
    beforeEach(() => {
      // Add various components for testing
      registry.addComponent(new Component({
        id: 'atom-1',
        name: 'Atom1',
        type: ComponentType.ATOM,
        filePath: 'src/atoms/Atom1.tsx',
        category: 'input',
        description: '',
        props: [],
        dependencies: [],
        usageCount: 5,
        deprecated: false,
        version: '1.0.0'
      }));

      registry.addComponent(new Component({
        id: 'molecule-1',
        name: 'Molecule1',
        type: ComponentType.MOLECULE,
        filePath: 'src/molecules/Molecule1.tsx',
        category: 'form',
        description: '',
        props: [],
        dependencies: [],
        usageCount: 3,
        deprecated: false,
        version: '1.0.0'
      }));

      registry.addComponent(new Component({
        id: 'organism-1',
        name: 'Organism1',
        type: ComponentType.ORGANISM,
        filePath: 'src/organisms/Organism1.tsx',
        category: 'layout',
        description: '',
        props: [],
        dependencies: [],
        usageCount: 1,
        deprecated: true,
        version: '1.0.0'
      }));
    });

    it('should get components by type', () => {
      const atoms = registry.getComponentsByType(ComponentType.ATOM);
      expect(atoms).toHaveLength(1);
      expect(atoms[0].id).toBe('atom-1');

      const molecules = registry.getComponentsByType(ComponentType.MOLECULE);
      expect(molecules).toHaveLength(1);
      expect(molecules[0].id).toBe('molecule-1');
    });

    it('should get components by category', () => {
      const inputComponents = registry.getComponentsByCategory('input');
      expect(inputComponents).toHaveLength(1);
      expect(inputComponents[0].id).toBe('atom-1');
    });

    it('should get deprecated components', () => {
      const deprecated = registry.getDeprecatedComponents();
      expect(deprecated).toHaveLength(1);
      expect(deprecated[0].id).toBe('organism-1');
    });

    it('should get most used components', () => {
      const mostUsed = registry.getMostUsedComponents(2);
      expect(mostUsed).toHaveLength(2);
      expect(mostUsed[0].id).toBe('atom-1');
      expect(mostUsed[1].id).toBe('molecule-1');
    });

    it('should search components by name', () => {
      const results = registry.searchByName('Molecule');
      expect(results).toHaveLength(1);
      expect(results[0].id).toBe('molecule-1');
    });

    it('should search components by file path', () => {
      const results = registry.searchByFilePath('molecules/');
      expect(results).toHaveLength(1);
      expect(results[0].id).toBe('molecule-1');
    });

    it('should get all components', () => {
      const all = registry.getAllComponents();
      expect(all).toHaveLength(3);
    });
  });

  describe('Statistics', () => {
    beforeEach(() => {
      // Add components for statistics
      registry.addComponent(new Component({
        id: 'atom-1',
        name: 'Atom1',
        type: ComponentType.ATOM,
        filePath: 'src/atoms/Atom1.tsx',
        category: 'input',
        description: '',
        props: [],
        dependencies: [],
        usageCount: 10,
        deprecated: false,
        version: '1.0.0'
      }));

      registry.addComponent(new Component({
        id: 'atom-2',
        name: 'Atom2',
        type: ComponentType.ATOM,
        filePath: 'src/atoms/Atom2.tsx',
        category: 'input',
        description: '',
        props: [],
        dependencies: [],
        usageCount: 5,
        deprecated: false,
        version: '1.0.0'
      }));

      registry.addComponent(new Component({
        id: 'molecule-1',
        name: 'Molecule1',
        type: ComponentType.MOLECULE,
        filePath: 'src/molecules/Molecule1.tsx',
        category: 'form',
        description: '',
        props: [],
        dependencies: [],
        usageCount: 3,
        deprecated: false,
        version: '1.0.0'
      }));
    });

    it('should calculate registry statistics', () => {
      const stats = registry.getStatistics();
      
      expect(stats.totalComponents).toBe(3);
      expect(stats.atomCount).toBe(2);
      expect(stats.moleculeCount).toBe(1);
      expect(stats.organismCount).toBe(0);
      expect(stats.averageUsageCount).toBeCloseTo(6, 1);
    });

    it('should identify most used components in statistics', () => {
      const stats = registry.getStatistics();
      expect(stats.mostUsedComponents).toContain('atom-1');
      expect(stats.mostUsedComponents).toHaveLength(Math.min(10, 3));
    });

    it('should count duplicates in statistics', () => {
      // This would be set by duplicate detection
      registry.setDuplicatesCount(5);
      const stats = registry.getStatistics();
      expect(stats.duplicatesFound).toBe(5);
    });

    it('should count inconsistencies in statistics', () => {
      // This would be set by validation
      registry.setInconsistenciesCount(3);
      const stats = registry.getStatistics();
      expect(stats.inconsistenciesFound).toBe(3);
    });
  });

  describe('Persistence', () => {
    it('should serialize registry to JSON', () => {
      registry.addComponent(testComponent);
      const json = registry.toJSON();
      
      expect(json).toHaveProperty('id');
      expect(json).toHaveProperty('projectName');
      expect(json).toHaveProperty('components');
      expect(json).toHaveProperty('lastUpdated');
      expect(json).toHaveProperty('statistics');
      expect(json).toHaveProperty('version');
    });

    it('should deserialize registry from JSON', () => {
      registry.addComponent(testComponent);
      const json = registry.toJSON();
      
      const restored = ComponentRegistry.fromJSON(json);
      expect(restored.id).toBe(registry.id);
      expect(restored.projectName).toBe(registry.projectName);
      expect(restored.getComponentCount()).toBe(1);
      expect(restored.getComponent('test-component')).toBeDefined();
    });

    it('should save registry to file', async () => {
      registry.addComponent(testComponent);
      const filePath = '/tmp/test-registry.json';
      
      await registry.saveToFile(filePath);
      const loaded = await ComponentRegistry.loadFromFile(filePath);
      
      expect(loaded.getComponentCount()).toBe(1);
      expect(loaded.getComponent('test-component')).toBeDefined();
    });
  });

  describe('Validation', () => {
    it('should validate registry integrity', () => {
      registry.addComponent(testComponent);
      const validation = registry.validate();
      
      expect(validation.isValid).toBe(true);
      expect(validation.errors).toHaveLength(0);
    });

    it('should detect orphaned references', () => {
      const component = new Component({
        id: 'component-with-ref',
        name: 'ComponentWithRef',
        type: ComponentType.MOLECULE,
        filePath: 'src/ComponentWithRef.tsx',
        category: 'test',
        description: '',
        props: [],
        dependencies: ['nonexistent-component'],
        usageCount: 0,
        deprecated: false,
        version: '1.0.0'
      });

      registry.addComponent(component);
      const validation = registry.validate();
      
      expect(validation.isValid).toBe(false);
      expect(validation.errors).toContain('Orphaned reference: nonexistent-component');
    });
  });

  describe('Batch Operations', () => {
    it('should add multiple components in batch', () => {
      const components = [
        new Component({
          id: 'batch-1',
          name: 'Batch1',
          type: ComponentType.ATOM,
          filePath: 'src/Batch1.tsx',
          category: 'test',
          description: '',
          props: [],
          dependencies: [],
          usageCount: 0,
          deprecated: false,
          version: '1.0.0'
        }),
        new Component({
          id: 'batch-2',
          name: 'Batch2',
          type: ComponentType.ATOM,
          filePath: 'src/Batch2.tsx',
          category: 'test',
          description: '',
          props: [],
          dependencies: [],
          usageCount: 0,
          deprecated: false,
          version: '1.0.0'
        })
      ];

      registry.addComponents(components);
      expect(registry.getComponentCount()).toBe(2);
    });

    it('should clear all components', () => {
      registry.addComponent(testComponent);
      expect(registry.getComponentCount()).toBe(1);
      
      registry.clear();
      expect(registry.getComponentCount()).toBe(0);
    });
  });

  describe('Event Handling', () => {
    it('should emit events on component changes', () => {
      let eventFired = false;
      registry.on('component:added', () => {
        eventFired = true;
      });

      registry.addComponent(testComponent);
      expect(eventFired).toBe(true);
    });

    it('should update lastUpdated timestamp on changes', () => {
      const initialTimestamp = registry.lastUpdated;
      
      setTimeout(() => {
        registry.addComponent(testComponent);
        expect(registry.lastUpdated.getTime()).toBeGreaterThan(initialTimestamp.getTime());
      }, 10);
    });
  });
});