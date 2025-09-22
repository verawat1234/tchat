/**
 * ComponentAnalyzer - Main analyzer that orchestrates all modules
 */

import { EventEmitter } from 'events';
import * as fs from 'fs';
import * as path from 'path';
import { glob } from 'glob';
import { Component, ComponentType } from './models/Component';
import { Atom } from './models/Atom';
import { Molecule } from './models/Molecule';
import { Organism } from './models/Organism';
import { ComponentRegistry } from './models/ComponentRegistry';
import { ASTParser } from './parser/ASTParser';
import { Categorizer } from './categorization/Categorizer';
import { DuplicateDetector } from './duplicates/DuplicateDetector';
import { ConsistencyValidator } from './validation/ConsistencyValidator';
import { DocumentationGenerator } from './docs/DocumentationGenerator';
import { ComponentAnalyzerConfig, DEFAULT_CONFIG } from './types/config';

export interface AnalysisResult {
  componentsFound: number;
  categorized: {
    atoms: number;
    molecules: number;
    organisms: number;
  };
  uncategorized: number;
  errors: string[];
  duration: number;
  registry?: ComponentRegistry;
  details?: ComponentDetail[];
}

export interface ComponentDetail {
  name: string;
  type: string;
  path: string;
  category?: string;
  usage?: number;
}

export class ComponentAnalyzer extends EventEmitter {
  private config: ComponentAnalyzerConfig;
  private parser: ASTParser;
  private categorizer: Categorizer;
  private duplicateDetector: DuplicateDetector;
  private validator: ConsistencyValidator;
  private docGenerator: DocumentationGenerator;
  private registry: ComponentRegistry;

  constructor(config?: Partial<ComponentAnalyzerConfig>) {
    super();
    this.config = { ...DEFAULT_CONFIG, ...config };
    this.parser = new ASTParser();
    this.categorizer = new Categorizer(this.config);
    this.duplicateDetector = new DuplicateDetector();
    this.validator = new ConsistencyValidator();
    this.docGenerator = new DocumentationGenerator();
    this.registry = new ComponentRegistry({
      id: `analyzer-${Date.now()}`,
      projectName: this.getProjectName()
    });
  }

  /**
   * Analyze components in a directory
   */
  async analyze(targetPath: string, options: any = {}): Promise<AnalysisResult> {
    const startTime = Date.now();
    const errors: string[] = [];
    const mergedOptions = { ...this.config.analysis, ...options };

    try {
      // Find all component files
      const files = await this.findComponentFiles(targetPath, mergedOptions);
      this.emit('analysis:start', { totalFiles: files.length });

      // Process each file
      let processed = 0;
      for (const file of files) {
        try {
          await this.processFile(file);
          processed++;
          this.emit('file:processed', {
            file: path.basename(file),
            index: processed,
            total: files.length
          });
        } catch (error: any) {
          errors.push(`Failed to process ${file}: ${error.message}`);
          this.emit('file:error', { file, error: error.message });
        }
      }

      // Run duplicate detection
      if (this.registry.getComponentCount() > 0) {
        const duplicates = await this.duplicateDetector.findDuplicates(
          this.registry,
          { threshold: this.config.duplicates.threshold }
        );
        this.registry.setDuplicatesCount(duplicates.length);
      }

      // Run validation
      const validationResults = this.validator.validateRegistry(this.registry);
      const inconsistencies = validationResults.results.filter(r => !r.valid).length;
      this.registry.setInconsistenciesCount(inconsistencies);

      const duration = Date.now() - startTime;
      const stats = this.registry.getStatistics();

      const result: AnalysisResult = {
        componentsFound: stats.totalComponents,
        categorized: {
          atoms: stats.atomCount,
          molecules: stats.moleculeCount,
          organisms: stats.organismCount
        },
        uncategorized: 0,
        errors,
        duration,
        registry: this.registry
      };

      // Add details if verbose
      if (options.verbose) {
        result.details = this.registry.getAllComponents().map(c => ({
          name: c.name,
          type: c.type,
          path: c.filePath,
          category: c.category,
          usage: c.usageCount
        }));
      }

      this.emit('analysis:complete', result);
      return result;

    } catch (error: any) {
      this.emit('analysis:error', error);
      throw error;
    }
  }

  /**
   * Find all component files in directory
   */
  private async findComponentFiles(targetPath: string, options: any): Promise<string[]> {
    const patterns = [
      '**/*.tsx',
      '**/*.jsx',
      '**/*.ts',
      '**/*.js'
    ];

    const excludePatterns = options.exclude || this.config.analysis.exclude;

    const files: string[] = [];

    for (const pattern of patterns) {
      const matches = await glob(pattern, {
        cwd: targetPath,
        ignore: excludePatterns,
        absolute: true,
        nodir: true,
        maxDepth: options.maxDepth || this.config.analysis.maxDepth
      });
      files.push(...matches);
    }

    // Filter out non-component files
    return files.filter(file => {
      const basename = path.basename(file);

      // Skip test files unless explicitly included
      if (!options.includeTests && (basename.includes('.test.') || basename.includes('.spec.'))) {
        return false;
      }

      // Skip story files
      if (basename.includes('.stories.')) {
        return false;
      }

      // Skip config files
      if (basename.includes('.config.')) {
        return false;
      }

      return true;
    });
  }

  /**
   * Process a single file
   */
  private async processFile(filePath: string): Promise<void> {
    // Read file content
    const content = fs.readFileSync(filePath, 'utf8');

    // Parse with AST
    const parseResult = this.parser.parse(content, filePath);

    // Skip if not a component
    if (!parseResult.isComponent || !parseResult.componentInfo) {
      return;
    }

    // Categorize component
    const categorizationResult = this.categorizer.categorize(parseResult, filePath);

    // Create component based on type
    const component = this.createComponent(
      parseResult.componentInfo.name,
      categorizationResult.type,
      filePath,
      parseResult
    );

    // Add to registry
    try {
      this.registry.addComponent(component);
    } catch (error: any) {
      // Component might already exist, try updating
      if (error.message.includes('already exists')) {
        this.registry.updateComponent(component);
      } else {
        throw error;
      }
    }
  }

  /**
   * Create component instance based on type
   */
  private createComponent(
    name: string,
    type: ComponentType,
    filePath: string,
    parseResult: any
  ): Component {
    const baseOptions = {
      id: this.generateComponentId(name, filePath),
      name,
      filePath,
      category: this.detectCategory(filePath, name),
      description: '',
      props: parseResult.componentInfo?.props || [],
      dependencies: parseResult.dependencies || [],
      usageCount: 0,
      deprecated: false,
      version: '1.0.0'
    };

    switch (type) {
      case ComponentType.ATOM:
        return new Atom({
          ...baseOptions,
          type: ComponentType.ATOM,
          htmlElement: this.detectHtmlElement(parseResult),
          variants: [],
          accessibility: {
            ariaLabel: false,
            ariaDescribedBy: false,
            role: null,
            keyboardNav: false,
            wcagLevel: 'AA'
          }
        });

      case ComponentType.MOLECULE:
        return new Molecule({
          ...baseOptions,
          type: ComponentType.MOLECULE,
          composition: [],
          layout: 'horizontal' as any,
          interactions: [],
          slots: []
        });

      case ComponentType.ORGANISM:
        return new Organism({
          ...baseOptions,
          type: ComponentType.ORGANISM,
          contains: [],
          standalone: true,
          pageSection: false,
          dataSource: null
        });

      default:
        return new Component({
          ...baseOptions,
          type
        });
    }
  }

  /**
   * Generate unique component ID
   */
  private generateComponentId(name: string, filePath: string): string {
    const pathParts = filePath.split('/');
    const fileName = pathParts[pathParts.length - 1].replace(/\.(tsx?|jsx?)$/, '');
    return `${fileName}-${name}`.toLowerCase().replace(/[^a-z0-9-]/g, '-');
  }

  /**
   * Detect component category from path and name
   */
  private detectCategory(filePath: string, name: string): string {
    const lowerPath = filePath.toLowerCase();
    const lowerName = name.toLowerCase();

    // Check path-based categories
    if (lowerPath.includes('button') || lowerName.includes('button')) return 'action';
    if (lowerPath.includes('input') || lowerName.includes('input')) return 'input';
    if (lowerPath.includes('form') || lowerName.includes('form')) return 'form';
    if (lowerPath.includes('navigation') || lowerName.includes('nav')) return 'navigation';
    if (lowerPath.includes('layout') || lowerName.includes('layout')) return 'layout';
    if (lowerPath.includes('modal') || lowerName.includes('modal')) return 'overlay';
    if (lowerPath.includes('card') || lowerName.includes('card')) return 'display';
    if (lowerPath.includes('list') || lowerName.includes('list')) return 'display';
    if (lowerPath.includes('table') || lowerName.includes('table')) return 'data';
    if (lowerPath.includes('chart') || lowerName.includes('chart')) return 'data';

    return 'general';
  }

  /**
   * Detect HTML element for atoms
   */
  private detectHtmlElement(parseResult: any): string {
    if (parseResult.jsxElements && parseResult.jsxElements.length > 0) {
      // Find the most used HTML element
      const htmlElements = parseResult.jsxElements.filter((el: string) =>
        el[0] === el[0].toLowerCase()
      );

      if (htmlElements.length > 0) {
        return htmlElements[0];
      }
    }

    return 'div';
  }

  /**
   * Save registry to file
   */
  async saveRegistry(filePath?: string): Promise<void> {
    const outputPath = filePath || this.config.paths.registry;
    await this.registry.saveToFile(outputPath);
    this.emit('registry:saved', { path: outputPath });
  }

  /**
   * Load registry from file
   */
  async loadRegistry(filePath?: string): Promise<void> {
    const inputPath = filePath || this.config.paths.registry;
    this.registry = await ComponentRegistry.loadFromFile(inputPath);
    this.emit('registry:loaded', { path: inputPath });
  }

  /**
   * Get current registry
   */
  getRegistry(): ComponentRegistry {
    return this.registry;
  }

  /**
   * Find duplicates
   */
  async findDuplicates(options?: any): Promise<any> {
    return this.duplicateDetector.findDuplicates(
      this.registry,
      { ...this.config.duplicates, ...options }
    );
  }

  /**
   * Validate components
   */
  validateComponents(options?: any): any {
    return this.validator.validateRegistry(this.registry, options);
  }

  /**
   * Generate documentation
   */
  async generateDocumentation(options?: any): Promise<any> {
    return this.docGenerator.generateFromRegistry(
      this.registry,
      { ...this.config.output, ...options }
    );
  }

  /**
   * Get project name from package.json
   */
  private getProjectName(): string {
    try {
      const packageJsonPath = path.join(process.cwd(), 'package.json');
      if (fs.existsSync(packageJsonPath)) {
        const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));
        return packageJson.name || 'unknown-project';
      }
    } catch (error) {
      // Ignore errors
    }
    return 'unknown-project';
  }
}