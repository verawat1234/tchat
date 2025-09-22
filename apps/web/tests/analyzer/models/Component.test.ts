import { describe, it, expect, beforeEach } from 'vitest';
import { Component, ComponentType, PropDefinition } from '../../../src/lib/analyzer/models/Component';

describe('Component Entity', () => {
  let component: Component;

  beforeEach(() => {
    component = new Component({
      id: 'btn-primary',
      name: 'PrimaryButton',
      type: ComponentType.ATOM,
      filePath: 'src/components/atoms/Button.tsx',
      category: 'action',
      description: 'Primary action button',
      props: [],
      dependencies: ['react'],
      usageCount: 0,
      deprecated: false,
      version: '1.0.0'
    });
  });

  describe('Constructor', () => {
    it('should create a component with required properties', () => {
      expect(component.id).toBe('btn-primary');
      expect(component.name).toBe('PrimaryButton');
      expect(component.type).toBe(ComponentType.ATOM);
      expect(component.filePath).toBe('src/components/atoms/Button.tsx');
    });

    it('should initialize timestamps', () => {
      expect(component.createdAt).toBeInstanceOf(Date);
      expect(component.updatedAt).toBeInstanceOf(Date);
    });

    it('should validate component type', () => {
      expect(() => {
        new Component({
          id: 'test',
          name: 'Test',
          type: 'invalid' as ComponentType,
          filePath: 'test.tsx',
          category: 'test',
          description: '',
          props: [],
          dependencies: [],
          usageCount: 0,
          deprecated: false,
          version: '1.0.0'
        });
      }).toThrow('Invalid component type');
    });
  });

  describe('Props Management', () => {
    it('should add prop definitions', () => {
      const prop: PropDefinition = {
        name: 'onClick',
        type: '() => void',
        required: false,
        defaultValue: undefined,
        description: 'Click handler',
        examples: ['() => console.log("clicked")']
      };

      component.addProp(prop);
      expect(component.props).toHaveLength(1);
      expect(component.props[0]).toEqual(prop);
    });

    it('should validate prop definitions', () => {
      expect(() => {
        component.addProp({
          name: '',
          type: 'string',
          required: false,
          defaultValue: '',
          description: '',
          examples: []
        });
      }).toThrow('Invalid prop definition');
    });

    it('should get prop by name', () => {
      const prop: PropDefinition = {
        name: 'variant',
        type: 'string',
        required: false,
        defaultValue: 'primary',
        description: 'Button variant',
        examples: ['primary', 'secondary']
      };

      component.addProp(prop);
      expect(component.getProp('variant')).toEqual(prop);
      expect(component.getProp('nonexistent')).toBeUndefined();
    });
  });

  describe('Dependencies', () => {
    it('should add dependencies', () => {
      component.addDependency('@radix-ui/react-slot');
      expect(component.dependencies).toContain('@radix-ui/react-slot');
    });

    it('should not add duplicate dependencies', () => {
      component.addDependency('react');
      component.addDependency('react');
      expect(component.dependencies.filter(d => d === 'react')).toHaveLength(1);
    });

    it('should check if has dependency', () => {
      expect(component.hasDependency('react')).toBe(true);
      expect(component.hasDependency('vue')).toBe(false);
    });
  });

  describe('Usage Tracking', () => {
    it('should increment usage count', () => {
      const initialCount = component.usageCount;
      component.incrementUsage();
      expect(component.usageCount).toBe(initialCount + 1);
    });

    it('should update timestamp when incrementing usage', () => {
      const initialTimestamp = component.updatedAt;
      // Wait a bit to ensure timestamp difference
      setTimeout(() => {
        component.incrementUsage();
        expect(component.updatedAt.getTime()).toBeGreaterThan(initialTimestamp.getTime());
      }, 10);
    });
  });

  describe('Deprecation', () => {
    it('should mark component as deprecated', () => {
      expect(component.deprecated).toBe(false);
      component.markAsDeprecated('Use NewButton instead');
      expect(component.deprecated).toBe(true);
      expect(component.deprecationMessage).toBe('Use NewButton instead');
    });

    it('should update timestamp when marking as deprecated', () => {
      const initialTimestamp = component.updatedAt;
      setTimeout(() => {
        component.markAsDeprecated();
        expect(component.updatedAt.getTime()).toBeGreaterThan(initialTimestamp.getTime());
      }, 10);
    });
  });

  describe('Validation', () => {
    it('should validate component structure', () => {
      expect(component.isValid()).toBe(true);
    });

    it('should fail validation for invalid file path', () => {
      component.filePath = '';
      expect(component.isValid()).toBe(false);
    });

    it('should fail validation for invalid ID', () => {
      component.id = '';
      expect(component.isValid()).toBe(false);
    });
  });

  describe('Serialization', () => {
    it('should serialize to JSON', () => {
      const json = component.toJSON();
      expect(json).toHaveProperty('id');
      expect(json).toHaveProperty('name');
      expect(json).toHaveProperty('type');
      expect(json).toHaveProperty('filePath');
      expect(json).not.toHaveProperty('_internalState');
    });

    it('should deserialize from JSON', () => {
      const json = component.toJSON();
      const restored = Component.fromJSON(json);
      expect(restored.id).toBe(component.id);
      expect(restored.name).toBe(component.name);
      expect(restored.type).toBe(component.type);
    });
  });

  describe('Cloning', () => {
    it('should create a deep clone', () => {
      const clone = component.clone();
      expect(clone).not.toBe(component);
      expect(clone.id).toBe(component.id);
      expect(clone.props).not.toBe(component.props);
      expect(clone.props).toEqual(component.props);
    });

    it('should allow modification of clone without affecting original', () => {
      const clone = component.clone();
      clone.name = 'ModifiedButton';
      expect(component.name).toBe('PrimaryButton');
      expect(clone.name).toBe('ModifiedButton');
    });
  });
});