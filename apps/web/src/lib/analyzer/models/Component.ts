/**
 * Component entity representing a UI component
 */

export enum ComponentType {
  ATOM = 'atom',
  MOLECULE = 'molecule',
  ORGANISM = 'organism'
}

export interface PropDefinition {
  name: string;
  type: string;
  required: boolean;
  defaultValue?: any;
  description: string;
  examples: string[];
}

export interface ComponentOptions {
  id: string;
  name: string;
  type: ComponentType;
  filePath: string;
  category: string;
  description: string;
  props: PropDefinition[];
  dependencies: string[];
  usageCount: number;
  deprecated: boolean;
  version: string;
}

export class Component {
  id: string;
  name: string;
  type: ComponentType;
  filePath: string;
  category: string;
  description: string;
  props: PropDefinition[];
  dependencies: string[];
  usageCount: number;
  deprecated: boolean;
  deprecationMessage?: string;
  version: string;
  createdAt: Date;
  updatedAt: Date;
  composedOf?: string[];
  usesCompoundPattern?: boolean;
  usesRenderProps?: boolean;
  isHOC?: boolean;
  isForwardRef?: boolean;

  constructor(options: ComponentOptions) {
    // Validate component type
    const validTypes = Object.values(ComponentType);
    if (!validTypes.includes(options.type)) {
      throw new Error('Invalid component type');
    }

    this.id = options.id;
    this.name = options.name;
    this.type = options.type;
    this.filePath = options.filePath;
    this.category = options.category;
    this.description = options.description;
    this.props = options.props ? [...options.props] : [];
    this.dependencies = options.dependencies ? [...options.dependencies] : [];
    this.usageCount = options.usageCount;
    this.deprecated = options.deprecated;
    this.version = options.version;
    this.createdAt = new Date();
    this.updatedAt = new Date();
  }

  /**
   * Add a prop definition to the component
   */
  addProp(prop: PropDefinition): void {
    if (!prop.name || prop.name.trim() === '') {
      throw new Error('Invalid prop definition');
    }
    this.props.push(prop);
    this.updatedAt = new Date();
  }

  /**
   * Get a prop by name
   */
  getProp(name: string): PropDefinition | undefined {
    return this.props.find(p => p.name === name);
  }

  /**
   * Add a dependency
   */
  addDependency(dependency: string): void {
    if (!this.dependencies.includes(dependency)) {
      this.dependencies.push(dependency);
      this.updatedAt = new Date();
    }
  }

  /**
   * Check if component has a specific dependency
   */
  hasDependency(dependency: string): boolean {
    return this.dependencies.includes(dependency);
  }

  /**
   * Increment usage count
   */
  incrementUsage(): void {
    this.usageCount++;
    this.updatedAt = new Date();
  }

  /**
   * Mark component as deprecated
   */
  markAsDeprecated(message?: string): void {
    this.deprecated = true;
    this.deprecationMessage = message;
    this.updatedAt = new Date();
  }

  /**
   * Validate component structure
   */
  isValid(): boolean {
    return !!(
      this.id &&
      this.id.trim() !== '' &&
      this.name &&
      this.name.trim() !== '' &&
      this.filePath &&
      this.filePath.trim() !== '' &&
      this.type &&
      Object.values(ComponentType).includes(this.type)
    );
  }

  /**
   * Serialize to JSON
   */
  toJSON(): any {
    return {
      id: this.id,
      name: this.name,
      type: this.type,
      filePath: this.filePath,
      category: this.category,
      description: this.description,
      props: this.props,
      dependencies: this.dependencies,
      usageCount: this.usageCount,
      deprecated: this.deprecated,
      deprecationMessage: this.deprecationMessage,
      version: this.version,
      createdAt: this.createdAt.toISOString(),
      updatedAt: this.updatedAt.toISOString(),
      composedOf: this.composedOf,
      usesCompoundPattern: this.usesCompoundPattern,
      usesRenderProps: this.usesRenderProps,
      isHOC: this.isHOC,
      isForwardRef: this.isForwardRef
    };
  }

  /**
   * Create from JSON
   */
  static fromJSON(json: any): Component {
    const component = new Component({
      id: json.id,
      name: json.name,
      type: json.type,
      filePath: json.filePath,
      category: json.category,
      description: json.description,
      props: json.props || [],
      dependencies: json.dependencies || [],
      usageCount: json.usageCount || 0,
      deprecated: json.deprecated || false,
      version: json.version
    });

    if (json.deprecationMessage) {
      component.deprecationMessage = json.deprecationMessage;
    }
    if (json.createdAt) {
      component.createdAt = new Date(json.createdAt);
    }
    if (json.updatedAt) {
      component.updatedAt = new Date(json.updatedAt);
    }
    if (json.composedOf) {
      component.composedOf = json.composedOf;
    }
    if (json.usesCompoundPattern !== undefined) {
      component.usesCompoundPattern = json.usesCompoundPattern;
    }
    if (json.usesRenderProps !== undefined) {
      component.usesRenderProps = json.usesRenderProps;
    }
    if (json.isHOC !== undefined) {
      component.isHOC = json.isHOC;
    }
    if (json.isForwardRef !== undefined) {
      component.isForwardRef = json.isForwardRef;
    }

    return component;
  }

  /**
   * Create a deep clone of the component
   */
  clone(): Component {
    return Component.fromJSON(this.toJSON());
  }
}