/**
 * Tests for Atom entity
 */

import { describe, it, expect, beforeEach } from 'vitest';
import { Atom } from '../Atom';
import { ComponentType } from '../Component';

describe('Atom', () => {
  let atom: Atom;

  beforeEach(() => {
    atom = new Atom({
      id: 'atom-button-primary',
      name: 'PrimaryButton',
      type: ComponentType.ATOM,
      filePath: '/src/components/atoms/PrimaryButton.tsx',
      htmlElement: 'button',
      variants: ['default', 'large', 'small'],
      accessibility: {
        ariaLabel: true,
        ariaDescribedBy: false,
        role: 'button',
        keyboardNav: true,
        wcagLevel: 'AA'
      }
    });
  });

  describe('constructor', () => {
    it('should create an atom with required properties', () => {
      expect(atom).toBeInstanceOf(Atom);
      expect(atom.type).toBe(ComponentType.ATOM);
      expect(atom.htmlElement).toBe('button');
      expect(atom.name).toBe('PrimaryButton');
    });

    it('should initialize with default values when not provided', () => {
      const minimalAtom = new Atom({
        id: 'minimal-atom',
        name: 'MinimalAtom',
        type: ComponentType.ATOM,
        filePath: '/src/components/atoms/Minimal.tsx'
      });

      expect(minimalAtom.htmlElement).toBe('div');
      expect(minimalAtom.variants).toEqual([]);
      expect(minimalAtom.accessibility.ariaLabel).toBe(false);
      expect(minimalAtom.accessibility.wcagLevel).toBe('AA');
    });

    it('should validate component type is ATOM', () => {
      expect(() => {
        new Atom({
          id: 'wrong-type',
          name: 'WrongType',
          type: ComponentType.MOLECULE as any,
          filePath: '/src/components/wrong.tsx'
        });
      }).toThrow();
    });
  });

  describe('htmlElement property', () => {
    it('should get and set HTML element', () => {
      expect(atom.htmlElement).toBe('button');

      atom.htmlElement = 'a';
      expect(atom.htmlElement).toBe('a');
    });

    it('should handle custom elements', () => {
      atom.htmlElement = 'custom-element';
      expect(atom.htmlElement).toBe('custom-element');
    });
  });

  describe('variants property', () => {
    it('should manage variants array', () => {
      expect(atom.variants).toEqual(['default', 'large', 'small']);

      atom.variants.push('disabled');
      expect(atom.variants).toContain('disabled');
      expect(atom.variants.length).toBe(4);
    });

    it('should allow empty variants', () => {
      atom.variants = [];
      expect(atom.variants).toEqual([]);
    });
  });

  describe('accessibility property', () => {
    it('should manage accessibility settings', () => {
      expect(atom.accessibility.ariaLabel).toBe(true);
      expect(atom.accessibility.keyboardNav).toBe(true);
      expect(atom.accessibility.wcagLevel).toBe('AA');
    });

    it('should update individual accessibility properties', () => {
      atom.accessibility.ariaLabel = false;
      atom.accessibility.wcagLevel = 'AAA';
      atom.accessibility.role = 'navigation';

      expect(atom.accessibility.ariaLabel).toBe(false);
      expect(atom.accessibility.wcagLevel).toBe('AAA');
      expect(atom.accessibility.role).toBe('navigation');
    });

    it('should handle null role', () => {
      atom.accessibility.role = null;
      expect(atom.accessibility.role).toBeNull();
    });
  });

  describe('toJSON', () => {
    it('should serialize atom to JSON', () => {
      const json = atom.toJSON();

      expect(json).toMatchObject({
        id: 'atom-button-primary',
        name: 'PrimaryButton',
        type: ComponentType.ATOM,
        htmlElement: 'button',
        variants: ['default', 'large', 'small'],
        accessibility: {
          ariaLabel: true,
          ariaDescribedBy: false,
          role: 'button',
          keyboardNav: true,
          wcagLevel: 'AA'
        }
      });
    });

    it('should include inherited Component properties', () => {
      const json = atom.toJSON();

      expect(json).toHaveProperty('filePath');
      expect(json).toHaveProperty('category');
      expect(json).toHaveProperty('props');
      expect(json).toHaveProperty('dependencies');
      expect(json).toHaveProperty('version');
    });
  });

  describe('fromJSON static method', () => {
    it('should deserialize atom from JSON', () => {
      const json = atom.toJSON();
      const deserializedAtom = Atom.fromJSON(json);

      expect(deserializedAtom).toBeInstanceOf(Atom);
      expect(deserializedAtom.id).toBe(atom.id);
      expect(deserializedAtom.name).toBe(atom.name);
      expect(deserializedAtom.htmlElement).toBe(atom.htmlElement);
      expect(deserializedAtom.variants).toEqual(atom.variants);
      expect(deserializedAtom.accessibility).toEqual(atom.accessibility);
    });

    it('should handle dates correctly in deserialization', () => {
      const json = atom.toJSON();
      const deserializedAtom = Atom.fromJSON(json);

      expect(deserializedAtom.createdAt).toBeInstanceOf(Date);
      expect(deserializedAtom.updatedAt).toBeInstanceOf(Date);
    });
  });

  describe('atom-specific behaviors', () => {
    it('should represent simple UI elements', () => {
      const iconAtom = new Atom({
        id: 'atom-icon',
        name: 'Icon',
        type: ComponentType.ATOM,
        filePath: '/src/components/atoms/Icon.tsx',
        htmlElement: 'svg',
        accessibility: {
          ariaLabel: true,
          ariaDescribedBy: false,
          role: 'img',
          keyboardNav: false,
          wcagLevel: 'AA'
        }
      });

      expect(iconAtom.htmlElement).toBe('svg');
      expect(iconAtom.accessibility.role).toBe('img');
      expect(iconAtom.accessibility.keyboardNav).toBe(false);
    });

    it('should support form elements', () => {
      const inputAtom = new Atom({
        id: 'atom-input',
        name: 'TextInput',
        type: ComponentType.ATOM,
        filePath: '/src/components/atoms/TextInput.tsx',
        htmlElement: 'input',
        variants: ['text', 'email', 'password'],
        accessibility: {
          ariaLabel: true,
          ariaDescribedBy: true,
          role: null,
          keyboardNav: true,
          wcagLevel: 'AA'
        }
      });

      expect(inputAtom.htmlElement).toBe('input');
      expect(inputAtom.variants).toContain('password');
      expect(inputAtom.accessibility.ariaDescribedBy).toBe(true);
    });
  });
});