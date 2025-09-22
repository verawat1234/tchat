import { describe, it, expect, beforeEach } from 'vitest';
import { Molecule, Composition, LayoutType, Interaction } from '../../../src/lib/analyzer/models/Molecule';
import { ComponentType } from '../../../src/lib/analyzer/models/Component';

describe('Molecule Entity', () => {
  let molecule: Molecule;

  beforeEach(() => {
    molecule = new Molecule({
      id: 'search-bar',
      name: 'SearchBar',
      type: ComponentType.MOLECULE,
      filePath: 'src/components/molecules/SearchBar.tsx',
      category: 'input',
      description: 'Search input with button',
      props: [],
      dependencies: ['react'],
      usageCount: 5,
      deprecated: false,
      version: '1.0.0',
      composition: [],
      layout: LayoutType.HORIZONTAL,
      interactions: [],
      slots: []
    });
  });

  describe('Composition Management', () => {
    it('should add atom to composition', () => {
      const composition: Composition = {
        atomId: 'input-text',
        quantity: 1,
        required: true,
        role: 'search-input'
      };

      molecule.addComposition(composition);
      expect(molecule.composition).toHaveLength(1);
      expect(molecule.composition[0]).toEqual(composition);
    });

    it('should validate minimum composition requirement', () => {
      expect(molecule.isValidMolecule()).toBe(false);
      
      molecule.addComposition({
        atomId: 'input-text',
        quantity: 1,
        required: true,
        role: 'input'
      });
      
      expect(molecule.isValidMolecule()).toBe(false);
      
      molecule.addComposition({
        atomId: 'button-icon',
        quantity: 1,
        required: true,
        role: 'submit'
      });
      
      expect(molecule.isValidMolecule()).toBe(true);
    });

    it('should get composition by atom ID', () => {
      const composition: Composition = {
        atomId: 'input-text',
        quantity: 1,
        required: true,
        role: 'search-input'
      };

      molecule.addComposition(composition);
      expect(molecule.getCompositionByAtomId('input-text')).toEqual(composition);
      expect(molecule.getCompositionByAtomId('nonexistent')).toBeUndefined();
    });

    it('should calculate total atom count', () => {
      molecule.addComposition({
        atomId: 'input-text',
        quantity: 1,
        required: true,
        role: 'input'
      });
      
      molecule.addComposition({
        atomId: 'button-icon',
        quantity: 2,
        required: true,
        role: 'actions'
      });

      expect(molecule.getTotalAtomCount()).toBe(3);
    });

    it('should list required atoms', () => {
      molecule.addComposition({
        atomId: 'input-text',
        quantity: 1,
        required: true,
        role: 'input'
      });
      
      molecule.addComposition({
        atomId: 'button-icon',
        quantity: 1,
        required: false,
        role: 'action'
      });

      const required = molecule.getRequiredAtoms();
      expect(required).toHaveLength(1);
      expect(required[0]).toBe('input-text');
    });
  });

  describe('Layout Management', () => {
    it('should set layout type', () => {
      molecule.setLayout(LayoutType.VERTICAL);
      expect(molecule.layout).toBe(LayoutType.VERTICAL);
    });

    it('should validate layout type', () => {
      const validLayouts = Object.values(LayoutType);
      validLayouts.forEach(layout => {
        molecule.setLayout(layout);
        expect(molecule.layout).toBe(layout);
      });
    });
  });

  describe('Interactions', () => {
    it('should add interactions between atoms', () => {
      const interaction: Interaction = {
        trigger: 'onClick',
        source: 'button-submit',
        target: 'input-text',
        action: 'submit-search'
      };

      molecule.addInteraction(interaction);
      expect(molecule.interactions).toHaveLength(1);
      expect(molecule.interactions[0]).toEqual(interaction);
    });

    it('should get interactions by source', () => {
      const interaction1: Interaction = {
        trigger: 'onClick',
        source: 'button-submit',
        target: 'input-text',
        action: 'submit'
      };

      const interaction2: Interaction = {
        trigger: 'onChange',
        source: 'input-text',
        target: 'button-submit',
        action: 'enable'
      };

      molecule.addInteraction(interaction1);
      molecule.addInteraction(interaction2);

      const sourceInteractions = molecule.getInteractionsBySource('button-submit');
      expect(sourceInteractions).toHaveLength(1);
      expect(sourceInteractions[0]).toEqual(interaction1);
    });

    it('should get interactions by target', () => {
      const interaction: Interaction = {
        trigger: 'onClick',
        source: 'button-submit',
        target: 'input-text',
        action: 'clear'
      };

      molecule.addInteraction(interaction);
      const targetInteractions = molecule.getInteractionsByTarget('input-text');
      expect(targetInteractions).toHaveLength(1);
      expect(targetInteractions[0]).toEqual(interaction);
    });
  });

  describe('Slot Management', () => {
    it('should add slot definitions', () => {
      molecule.addSlot({
        name: 'prefix',
        description: 'Content before input',
        accepts: [ComponentType.ATOM],
        required: false,
        defaultContent: null
      });

      expect(molecule.slots).toHaveLength(1);
      expect(molecule.slots[0].name).toBe('prefix');
    });

    it('should get slot by name', () => {
      const slot = {
        name: 'suffix',
        description: 'Content after input',
        accepts: [ComponentType.ATOM],
        required: false,
        defaultContent: null
      };

      molecule.addSlot(slot);
      expect(molecule.getSlot('suffix')).toEqual(slot);
      expect(molecule.getSlot('nonexistent')).toBeUndefined();
    });

    it('should list required slots', () => {
      molecule.addSlot({
        name: 'icon',
        description: 'Icon slot',
        accepts: [ComponentType.ATOM],
        required: true,
        defaultContent: null
      });

      molecule.addSlot({
        name: 'tooltip',
        description: 'Tooltip slot',
        accepts: [ComponentType.ATOM],
        required: false,
        defaultContent: null
      });

      const required = molecule.getRequiredSlots();
      expect(required).toHaveLength(1);
      expect(required[0].name).toBe('icon');
    });
  });

  describe('Validation', () => {
    it('should validate molecule structure', () => {
      // Initially invalid (no composition)
      expect(molecule.isValidMolecule()).toBe(false);

      // Add minimum required atoms
      molecule.addComposition({
        atomId: 'input-text',
        quantity: 1,
        required: true,
        role: 'input'
      });
      
      molecule.addComposition({
        atomId: 'button-icon',
        quantity: 1,
        required: true,
        role: 'action'
      });

      expect(molecule.isValidMolecule()).toBe(true);
    });

    it('should validate composition references', () => {
      molecule.addComposition({
        atomId: '',
        quantity: 1,
        required: true,
        role: 'input'
      });

      expect(molecule.hasValidComposition()).toBe(false);

      molecule.composition[0].atomId = 'input-text';
      expect(molecule.hasValidComposition()).toBe(true);
    });
  });

  describe('Complexity Analysis', () => {
    it('should calculate complexity score', () => {
      // Add atoms
      molecule.addComposition({
        atomId: 'input-text',
        quantity: 1,
        required: true,
        role: 'input'
      });
      
      molecule.addComposition({
        atomId: 'button-icon',
        quantity: 2,
        required: true,
        role: 'actions'
      });

      // Add interactions
      molecule.addInteraction({
        trigger: 'onClick',
        source: 'button-1',
        target: 'input-text',
        action: 'clear'
      });

      molecule.addInteraction({
        trigger: 'onClick',
        source: 'button-2',
        target: 'input-text',
        action: 'submit'
      });

      const complexity = molecule.calculateComplexity();
      expect(complexity).toBeGreaterThan(0);
      expect(complexity).toBeLessThanOrEqual(10);
    });
  });

  describe('Serialization', () => {
    it('should serialize to JSON including molecule-specific properties', () => {
      molecule.addComposition({
        atomId: 'input-text',
        quantity: 1,
        required: true,
        role: 'input'
      });

      const json = molecule.toJSON();
      expect(json).toHaveProperty('composition');
      expect(json).toHaveProperty('layout');
      expect(json).toHaveProperty('interactions');
      expect(json).toHaveProperty('slots');
    });

    it('should deserialize from JSON', () => {
      molecule.addComposition({
        atomId: 'input-text',
        quantity: 1,
        required: true,
        role: 'input'
      });

      molecule.addInteraction({
        trigger: 'onClick',
        source: 'button',
        target: 'input',
        action: 'submit'
      });

      const json = molecule.toJSON();
      const restored = Molecule.fromJSON(json);
      
      expect(restored.composition).toEqual(molecule.composition);
      expect(restored.layout).toBe(molecule.layout);
      expect(restored.interactions).toEqual(molecule.interactions);
    });
  });
});