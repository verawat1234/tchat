/**
 * Categorizer - Determines component type based on atomic design principles
 */

import { ComponentType } from '../models/Component';
import { ParseResult } from '../parser/ASTParser';
import { ComponentAnalyzerConfig } from '../types/config';

export interface CategorizationResult {
  type: ComponentType;
  confidence: number;
  reasoning: string[];
  composition?: string[];
}

export class Categorizer {
  private config: ComponentAnalyzerConfig;

  constructor(config?: Partial<ComponentAnalyzerConfig>) {
    this.config = {
      ...this.getDefaultConfig(),
      ...config
    };
  }

  /**
   * Categorize a component based on its AST analysis
   */
  categorize(parseResult: ParseResult, filePath: string): CategorizationResult {
    const reasoning: string[] = [];

    // Check patterns first
    const patternType = this.checkPatterns(filePath);
    if (patternType) {
      reasoning.push(`Matched ${patternType} pattern`);
      return {
        type: patternType,
        confidence: 0.9,
        reasoning
      };
    }

    // Analyze component complexity
    const complexity = this.calculateComplexity(parseResult);
    const composition = this.analyzeComposition(parseResult);

    // Decision tree for categorization
    if (this.isAtom(parseResult, complexity)) {
      reasoning.push('Single responsibility component');
      reasoning.push('No composition of other components');
      if (parseResult.jsxElements && parseResult.jsxElements.length <= 3) {
        reasoning.push('Simple JSX structure');
      }

      return {
        type: ComponentType.ATOM,
        confidence: this.calculateConfidence(parseResult, ComponentType.ATOM),
        reasoning
      };
    }

    if (this.isMolecule(parseResult, complexity, composition)) {
      reasoning.push('Composition of multiple atoms');
      reasoning.push(`Contains ${composition.length} child components`);
      if (parseResult.hasConditionalRendering) {
        reasoning.push('Has conditional rendering logic');
      }

      return {
        type: ComponentType.MOLECULE,
        confidence: this.calculateConfidence(parseResult, ComponentType.MOLECULE),
        reasoning,
        composition
      };
    }

    // Default to organism for complex components
    reasoning.push('Complex component structure');
    reasoning.push(`High complexity score: ${complexity}`);
    if (parseResult.hasListRendering) {
      reasoning.push('Contains list rendering');
    }
    if (composition.length > 5) {
      reasoning.push(`Large composition: ${composition.length} components`);
    }

    return {
      type: ComponentType.ORGANISM,
      confidence: this.calculateConfidence(parseResult, ComponentType.ORGANISM),
      reasoning,
      composition
    };
  }

  /**
   * Check file path against configured patterns
   */
  private checkPatterns(filePath: string): ComponentType | null {
    const fileName = filePath.split('/').pop() || '';

    if (this.config.analysis.categorization?.atomPatterns) {
      for (const pattern of this.config.analysis.categorization.atomPatterns) {
        if (this.matchesPattern(fileName, pattern)) {
          return ComponentType.ATOM;
        }
      }
    }

    if (this.config.analysis.categorization?.moleculePatterns) {
      for (const pattern of this.config.analysis.categorization.moleculePatterns) {
        if (this.matchesPattern(fileName, pattern)) {
          return ComponentType.MOLECULE;
        }
      }
    }

    if (this.config.analysis.categorization?.organismPatterns) {
      for (const pattern of this.config.analysis.categorization.organismPatterns) {
        if (this.matchesPattern(fileName, pattern)) {
          return ComponentType.ORGANISM;
        }
      }
    }

    return null;
  }

  /**
   * Check if filename matches a glob-like pattern
   */
  private matchesPattern(fileName: string, pattern: string): boolean {
    // Convert glob pattern to regex
    const regex = pattern
      .replace(/\*/g, '.*')
      .replace(/\?/g, '.');

    return new RegExp(`^${regex}$`).test(fileName);
  }

  /**
   * Calculate component complexity
   */
  private calculateComplexity(parseResult: ParseResult): number {
    let complexity = 0;

    // Factor in JSX elements
    if (parseResult.jsxElements) {
      complexity += parseResult.jsxElements.length * 0.5;
    }

    // Factor in hooks usage
    if (parseResult.hooks) {
      complexity += parseResult.hooks.length * 0.3;
    }

    // Factor in conditional rendering
    if (parseResult.hasConditionalRendering) {
      complexity += 1;
    }

    // Factor in list rendering
    if (parseResult.hasListRendering) {
      complexity += 1.5;
    }

    // Factor in props count
    if (parseResult.componentInfo?.props) {
      complexity += parseResult.componentInfo.props.length * 0.2;
    }

    // Factor in dependencies
    const externalDeps = parseResult.dependencies.filter(d =>
      !d.startsWith('.') && !d.startsWith('react')
    );
    complexity += externalDeps.length * 0.4;

    return Math.min(complexity, 10); // Cap at 10
  }

  /**
   * Analyze component composition
   */
  private analyzeComposition(parseResult: ParseResult): string[] {
    const composition: string[] = [];

    if (parseResult.jsxElements) {
      // Find components (capitalized JSX elements)
      parseResult.jsxElements.forEach(element => {
        if (element[0] === element[0].toUpperCase() && !composition.includes(element)) {
          composition.push(element);
        }
      });
    }

    // Add composed components from component info
    if (parseResult.componentInfo?.composedOf) {
      parseResult.componentInfo.composedOf.forEach(comp => {
        if (!composition.includes(comp)) {
          composition.push(comp);
        }
      });
    }

    return composition;
  }

  /**
   * Check if component is an atom
   */
  private isAtom(parseResult: ParseResult, complexity: number): boolean {
    // Atoms are simple, single-purpose components
    if (complexity > 3) return false;

    // Check for simple structure
    const elementCount = parseResult.jsxElements?.length || 0;
    if (elementCount > 5) return false;

    // Check for minimal composition
    const composition = this.analyzeComposition(parseResult);
    if (composition.length > 1) return false;

    // Check for simple props
    const propCount = parseResult.componentInfo?.props?.length || 0;
    if (propCount > 5) return false;

    // ForwardRef components are often atoms
    if (parseResult.componentInfo?.isForwardRef) {
      return true;
    }

    return true;
  }

  /**
   * Check if component is a molecule
   */
  private isMolecule(
    parseResult: ParseResult,
    complexity: number,
    composition: string[]
  ): boolean {
    // Molecules have moderate complexity
    if (complexity < 2 || complexity > 7) return false;

    // Must compose multiple atoms
    if (composition.length < 2) return false;

    // Should not be too complex
    if (composition.length > 5) return false;

    // Check for interaction patterns
    if (parseResult.hooks && parseResult.hooks.length > 0) {
      // Has state management, likely a molecule
      return true;
    }

    return true;
  }

  /**
   * Calculate confidence score for categorization
   */
  private calculateConfidence(parseResult: ParseResult, type: ComponentType): number {
    let confidence = 0.5; // Base confidence

    // Increase confidence for clear indicators
    switch (type) {
      case ComponentType.ATOM:
        if (parseResult.componentInfo?.isForwardRef) confidence += 0.3;
        if (!parseResult.hasConditionalRendering) confidence += 0.1;
        if (!parseResult.hasListRendering) confidence += 0.1;
        break;

      case ComponentType.MOLECULE:
        if (parseResult.hooks && parseResult.hooks.includes('useState')) confidence += 0.2;
        if (parseResult.hasConditionalRendering) confidence += 0.1;
        if (parseResult.componentInfo?.composedOf) confidence += 0.2;
        break;

      case ComponentType.ORGANISM:
        if (parseResult.hasListRendering) confidence += 0.2;
        if (parseResult.usesCompoundPattern) confidence += 0.2;
        if (parseResult.hooks && parseResult.hooks.length > 3) confidence += 0.1;
        break;
    }

    return Math.min(confidence, 1.0);
  }

  /**
   * Get default configuration
   */
  private getDefaultConfig(): ComponentAnalyzerConfig {
    return {
      paths: {
        components: 'src/components',
        output: 'docs/components',
        registry: 'docs/components/registry.json'
      },
      analysis: {
        recursive: true,
        maxDepth: 10,
        exclude: ['*.test.tsx', '*.stories.tsx'],
        includeTests: false,
        categorization: {
          atomPatterns: ['Button*.tsx', 'Input*.tsx', 'Icon*.tsx', 'Label*.tsx'],
          moleculePatterns: ['*Form.tsx', '*Card.tsx', '*List.tsx', '*Bar.tsx'],
          organismPatterns: ['*Section.tsx', '*Layout.tsx', '*Page.tsx', '*Header.tsx']
        }
      },
      validation: {
        rules: {},
        autoFix: false,
        severity: 'warning'
      },
      duplicates: {
        threshold: 75,
        autoMerge: false
      },
      output: {
        format: 'both',
        includeVisuals: true,
        includeExamples: true,
        colors: true,
        verbose: false
      }
    };
  }
}