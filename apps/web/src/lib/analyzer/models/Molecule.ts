/**
 * Molecule entity representing a combination of atoms
 */

import { Component, ComponentType, ComponentOptions } from './Component';

export enum LayoutType {
  HORIZONTAL = 'horizontal',
  VERTICAL = 'vertical',
  GRID = 'grid',
  ABSOLUTE = 'absolute',
  FLEXIBLE = 'flexible'
}

export interface Composition {
  atomId: string;
  quantity: number;
  required: boolean;
  role: string;
}

export interface Interaction {
  trigger: string;
  source: string;
  target: string;
  action: string;
}

export interface SlotDefinition {
  name: string;
  description: string;
  accepts: ComponentType[];
  required: boolean;
  defaultContent: string | null;
}

export interface MoleculeOptions extends ComponentOptions {
  composition: Composition[];
  layout: LayoutType;
  interactions: Interaction[];
  slots: SlotDefinition[];
}

export class Molecule extends Component {
  composition: Composition[];
  layout: LayoutType;
  interactions: Interaction[];
  slots: SlotDefinition[];

  constructor(options: MoleculeOptions) {
    super({
      ...options,
      type: ComponentType.MOLECULE
    });

    this.composition = options.composition || [];
    this.layout = options.layout || LayoutType.HORIZONTAL;
    this.interactions = options.interactions || [];
    this.slots = options.slots || [];
  }

  /**
   * Add an atom to the composition
   */
  addComposition(composition: Composition): void {
    this.composition.push(composition);
    this.updatedAt = new Date();
  }

  /**
   * Get composition by atom ID
   */
  getCompositionByAtomId(atomId: string): Composition | undefined {
    return this.composition.find(c => c.atomId === atomId);
  }

  /**
   * Calculate total atom count
   */
  getTotalAtomCount(): number {
    return this.composition.reduce((total, comp) => total + comp.quantity, 0);
  }

  /**
   * Get list of required atoms
   */
  getRequiredAtoms(): string[] {
    return this.composition
      .filter(c => c.required)
      .map(c => c.atomId);
  }

  /**
   * Set layout type
   */
  setLayout(layout: LayoutType): void {
    this.layout = layout;
    this.updatedAt = new Date();
  }

  /**
   * Add an interaction
   */
  addInteraction(interaction: Interaction): void {
    this.interactions.push(interaction);
    this.updatedAt = new Date();
  }

  /**
   * Get interactions by source
   */
  getInteractionsBySource(source: string): Interaction[] {
    return this.interactions.filter(i => i.source === source);
  }

  /**
   * Get interactions by target
   */
  getInteractionsByTarget(target: string): Interaction[] {
    return this.interactions.filter(i => i.target === target);
  }

  /**
   * Add a slot definition
   */
  addSlot(slot: SlotDefinition): void {
    this.slots.push(slot);
    this.updatedAt = new Date();
  }

  /**
   * Get slot by name
   */
  getSlot(name: string): SlotDefinition | undefined {
    return this.slots.find(s => s.name === name);
  }

  /**
   * Get list of required slots
   */
  getRequiredSlots(): SlotDefinition[] {
    return this.slots.filter(s => s.required);
  }

  /**
   * Check if molecule has minimum required composition
   */
  isValidMolecule(): boolean {
    // A molecule must have at least 2 atoms
    return this.composition.length >= 2 && this.isValid();
  }

  /**
   * Check if composition references are valid
   */
  hasValidComposition(): boolean {
    return this.composition.every(c => c.atomId && c.atomId.trim() !== '');
  }

  /**
   * Calculate complexity score (0-10)
   */
  calculateComplexity(): number {
    const atomCount = this.getTotalAtomCount();
    const interactionCount = this.interactions.length;
    const slotCount = this.slots.length;

    // Simple formula: weighted sum normalized to 0-10
    const complexity = Math.min(10,
      (atomCount * 0.5) +
      (interactionCount * 0.3) +
      (slotCount * 0.2)
    );

    return Math.round(complexity * 10) / 10; // Round to 1 decimal
  }

  /**
   * Serialize to JSON
   */
  toJSON(): any {
    const baseJson = super.toJSON();
    return {
      ...baseJson,
      composition: this.composition,
      layout: this.layout,
      interactions: this.interactions,
      slots: this.slots
    };
  }

  /**
   * Create from JSON
   */
  static fromJSON(json: any): Molecule {
    const molecule = new Molecule({
      id: json.id,
      name: json.name,
      type: ComponentType.MOLECULE,
      filePath: json.filePath,
      category: json.category,
      description: json.description,
      props: json.props || [],
      dependencies: json.dependencies || [],
      usageCount: json.usageCount || 0,
      deprecated: json.deprecated || false,
      version: json.version,
      composition: json.composition || [],
      layout: json.layout || LayoutType.HORIZONTAL,
      interactions: json.interactions || [],
      slots: json.slots || []
    });

    // Restore additional properties
    if (json.deprecationMessage) {
      molecule.deprecationMessage = json.deprecationMessage;
    }
    if (json.createdAt) {
      molecule.createdAt = new Date(json.createdAt);
    }
    if (json.updatedAt) {
      molecule.updatedAt = new Date(json.updatedAt);
    }
    if (json.composedOf) {
      molecule.composedOf = json.composedOf;
    }
    if (json.usesCompoundPattern !== undefined) {
      molecule.usesCompoundPattern = json.usesCompoundPattern;
    }
    if (json.usesRenderProps !== undefined) {
      molecule.usesRenderProps = json.usesRenderProps;
    }
    if (json.isHOC !== undefined) {
      molecule.isHOC = json.isHOC;
    }
    if (json.isForwardRef !== undefined) {
      molecule.isForwardRef = json.isForwardRef;
    }

    return molecule;
  }
}