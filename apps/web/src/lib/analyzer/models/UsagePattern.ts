/**
 * UsagePattern - tracks where and how components are used
 */

export interface PropPattern {
  prop: string;
  value: any;
  frequency: number;
}

export interface ComponentCombination {
  components: string[];
  frequency: number;
  description: string;
}

export interface UsagePatternOptions {
  id: string;
  componentId: string;
  usageCount: number;
  locations?: string[];
  propPatterns?: PropPattern[];
  commonCombinations?: ComponentCombination[];
  lastUpdated?: Date;
}

export class UsagePattern {
  id: string;
  componentId: string;
  private _usageCount: number;
  locations: string[];
  propPatterns: PropPattern[];
  commonCombinations: ComponentCombination[];
  lastUpdated: Date;

  constructor(options: UsagePatternOptions) {
    // Validate required fields
    if (!options.id) {
      throw new Error('Pattern ID is required');
    }
    if (!options.componentId) {
      throw new Error('Component ID is required');
    }

    this.id = options.id;
    this.componentId = options.componentId;
    this._usageCount = options.usageCount || 0;
    this.locations = options.locations || [];
    this.propPatterns = options.propPatterns || [];
    this.commonCombinations = options.commonCombinations || [];
    this.lastUpdated = options.lastUpdated || new Date();
  }

  /**
   * Get usage count
   */
  get usageCount(): number {
    return this._usageCount;
  }

  /**
   * Set usage count
   */
  set usageCount(value: number) {
    if (value < 0) {
      throw new Error('Usage count cannot be negative');
    }
    this._usageCount = value;
    this.lastUpdated = new Date();
  }

  /**
   * Increment usage count
   */
  incrementUsage(): void {
    this.usageCount++;
  }

  /**
   * Add a location where the component is used
   */
  addLocation(location: string): void {
    if (!this.locations.includes(location)) {
      this.locations.push(location);
      this.lastUpdated = new Date();
    }
  }

  /**
   * Clear all locations
   */
  clearLocations(): void {
    this.locations = [];
    this.lastUpdated = new Date();
  }

  /**
   * Add a prop pattern
   */
  addPropPattern(prop: string, value: any, frequency: number): void {
    if (frequency < 0 || frequency > 1) {
      throw new Error('Frequency must be between 0 and 1');
    }

    this.propPatterns.push({ prop, value, frequency });
    this.lastUpdated = new Date();
  }

  /**
   * Add a component combination pattern
   */
  addCombination(components: string[], frequency: number, description: string): void {
    if (components.length < 2) {
      throw new Error('Combination must have at least 2 components');
    }
    if (frequency < 0 || frequency > 1) {
      throw new Error('Frequency must be between 0 and 1');
    }

    this.commonCombinations.push({ components, frequency, description });
    this.lastUpdated = new Date();
  }

  /**
   * Get the most common value for a prop
   */
  getMostCommonPropValue(prop: string): any | null {
    const patterns = this.propPatterns
      .filter(p => p.prop === prop)
      .sort((a, b) => b.frequency - a.frequency);

    return patterns.length > 0 ? patterns[0].value : null;
  }

  /**
   * Get average prop frequency
   */
  getAveragePropFrequency(): number {
    if (this.propPatterns.length === 0) {
      return 0;
    }

    const sum = this.propPatterns.reduce((acc, p) => acc + p.frequency, 0);
    return sum / this.propPatterns.length;
  }

  /**
   * Get most frequent component combination
   */
  getMostFrequentCombination(): ComponentCombination | undefined {
    if (this.commonCombinations.length === 0) {
      return undefined;
    }

    return this.commonCombinations.reduce((max, current) =>
      current.frequency > max.frequency ? current : max
    );
  }

  /**
   * Serialize to JSON
   */
  toJSON(): any {
    return {
      id: this.id,
      componentId: this.componentId,
      usageCount: this.usageCount,
      locations: this.locations,
      propPatterns: this.propPatterns,
      commonCombinations: this.commonCombinations,
      lastUpdated: this.lastUpdated.toISOString()
    };
  }

  /**
   * Deserialize from JSON
   */
  static fromJSON(json: any): UsagePattern {
    return new UsagePattern({
      ...json,
      lastUpdated: new Date(json.lastUpdated)
    });
  }
}