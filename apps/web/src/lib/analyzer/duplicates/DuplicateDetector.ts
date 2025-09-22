/**
 * DuplicateDetector - Finds similar components that could be merged
 */

import { Component, ComponentType } from '../models/Component';
import { ComponentRegistry } from '../models/ComponentRegistry';
import { ParseResult } from '../parser/ASTParser';

export interface DuplicateGroup {
  similarity: number;
  components: ComponentSummary[];
  suggestedMerge: string;
  reasoning: string[];
}

export interface ComponentSummary {
  id: string;
  name: string;
  filePath: string;
  type: ComponentType;
}

export interface DuplicateDetectionOptions {
  threshold: number;
  factors?: {
    structural?: number;
    visual?: number;
    functional?: number;
  };
}

export class DuplicateDetector {
  private defaultOptions: DuplicateDetectionOptions = {
    threshold: 65,  // Lowered threshold to catch more duplicates
    factors: {
      structural: 0.3,
      visual: 0.4,    // Increased visual weight since name/file similarity is important
      functional: 0.3
    }
  };

  /**
   * Find duplicate components in registry
   */
  async findDuplicates(
    registry: ComponentRegistry,
    options?: Partial<DuplicateDetectionOptions>
  ): Promise<DuplicateGroup[]> {
    const opts = { ...this.defaultOptions, ...options };
    const components = registry.getAllComponents();
    const duplicateGroups: DuplicateGroup[] = [];
    const processed = new Set<string>();

    for (let i = 0; i < components.length; i++) {
      if (processed.has(components[i].id)) continue;

      const group: ComponentSummary[] = [{
        id: components[i].id,
        name: components[i].name,
        filePath: components[i].filePath,
        type: components[i].type
      }];

      const reasoning: string[] = [];
      let highestSimilarity = 0;

      for (let j = i + 1; j < components.length; j++) {
        if (processed.has(components[j].id)) continue;

        const similarity = this.calculateSimilarity(
          components[i],
          components[j],
          opts.factors!
        );

        if (similarity >= opts.threshold) {
          group.push({
            id: components[j].id,
            name: components[j].name,
            filePath: components[j].filePath,
            type: components[j].type
          });
          processed.add(components[j].id);
          highestSimilarity = Math.max(highestSimilarity, similarity);

          // Add reasoning
          if (this.haveSimilarNames(components[i].name, components[j].name)) {
            reasoning.push(`Similar names: ${components[i].name} and ${components[j].name}`);
          }
          if (components[i].type === components[j].type) {
            reasoning.push(`Same component type: ${components[i].type}`);
          }
          if (this.haveSimilarProps(components[i], components[j])) {
            reasoning.push('Similar prop structures');
          }
        }
      }

      if (group.length > 1) {
        processed.add(components[i].id);

        // Suggest which component to keep (usually the most used one)
        const suggestedMerge = this.suggestMergeTarget(
          group.map(g => components.find(c => c.id === g.id)!)
        );

        duplicateGroups.push({
          similarity: highestSimilarity,
          components: group,
          suggestedMerge,
          reasoning
        });
      }
    }

    return duplicateGroups.sort((a, b) => b.similarity - a.similarity);
  }

  /**
   * Calculate similarity between two components
   */
  private calculateSimilarity(
    comp1: Component,
    comp2: Component,
    factors: { structural?: number; visual?: number; functional?: number }
  ): number {
    const structuralWeight = factors.structural || 0.4;
    const visualWeight = factors.visual || 0.3;
    const functionalWeight = factors.functional || 0.3;

    const structuralSimilarity = this.calculateStructuralSimilarity(comp1, comp2);
    const visualSimilarity = this.calculateVisualSimilarity(comp1, comp2);
    const functionalSimilarity = this.calculateFunctionalSimilarity(comp1, comp2);

    const totalSimilarity =
      structuralSimilarity * structuralWeight +
      visualSimilarity * visualWeight +
      functionalSimilarity * functionalWeight;

    return Math.round(totalSimilarity * 100);
  }

  /**
   * Calculate structural similarity based on component structure
   */
  private calculateStructuralSimilarity(comp1: Component, comp2: Component): number {
    let similarity = 0;

    // Same type
    if (comp1.type === comp2.type) {
      similarity += 0.3;
    }

    // Similar prop count
    const propDiff = Math.abs(comp1.props.length - comp2.props.length);
    if (propDiff === 0) {
      similarity += 0.3;
    } else if (propDiff <= 2) {
      similarity += 0.2;
    } else if (propDiff <= 4) {
      similarity += 0.1;
    }

    // Similar prop names
    if (this.haveSimilarProps(comp1, comp2)) {
      similarity += 0.2;
    }

    // Similar dependencies
    const sharedDeps = this.getSharedDependencies(comp1, comp2);
    if (sharedDeps.length > 0) {
      similarity += Math.min(0.2, sharedDeps.length * 0.05);
    }

    return Math.min(similarity, 1.0);
  }

  /**
   * Calculate visual similarity (simplified - would need actual rendering analysis)
   */
  private calculateVisualSimilarity(comp1: Component, comp2: Component): number {
    let similarity = 0;

    // Similar names often indicate visual similarity
    if (this.haveSimilarNames(comp1.name, comp2.name)) {
      similarity += 0.5;
    }

    // Check if filenames are variants of the same component
    const file1 = comp1.filePath.split('/').pop() || '';
    const file2 = comp2.filePath.split('/').pop() || '';
    const baseFile1 = file1.replace(/(_new|_old|_working|_Fixed|_replacement|_temp|_backup|_copy|\d+)\.(tsx?|jsx?)$/i, '');
    const baseFile2 = file2.replace(/(_new|_old|_working|_Fixed|_replacement|_temp|_backup|_copy|\d+)\.(tsx?|jsx?)$/i, '');

    if (baseFile1 === baseFile2 && baseFile1 !== '') {
      similarity += 0.4; // Strong indicator of duplicate
    }

    // Same category often means similar visual purpose
    if (comp1.category === comp2.category) {
      similarity += 0.2;
    }

    // Similar file paths might indicate visual similarity
    const path1Parts = comp1.filePath.split('/');
    const path2Parts = comp2.filePath.split('/');
    const commonPath = path1Parts.filter(p => path2Parts.includes(p));
    if (commonPath.length > 1) {
      similarity += 0.1;
    }

    return Math.min(similarity, 1.0);
  }

  /**
   * Calculate functional similarity
   */
  private calculateFunctionalSimilarity(comp1: Component, comp2: Component): number {
    let similarity = 0;

    // Similar categories indicate similar function
    if (comp1.category === comp2.category) {
      similarity += 0.4;
    }

    // Similar prop types indicate similar functionality
    const propTypes1 = comp1.props.map(p => p.type).sort();
    const propTypes2 = comp2.props.map(p => p.type).sort();
    const commonTypes = propTypes1.filter(t => propTypes2.includes(t));
    if (commonTypes.length > 0) {
      similarity += Math.min(0.3, commonTypes.length * 0.1);
    }

    // Both deprecated might indicate they serve similar obsolete function
    if (comp1.deprecated && comp2.deprecated) {
      similarity += 0.1;
    }

    // Similar usage count might indicate similar function importance
    const usageDiff = Math.abs(comp1.usageCount - comp2.usageCount);
    if (usageDiff < 5) {
      similarity += 0.2;
    } else if (usageDiff < 10) {
      similarity += 0.1;
    }

    return Math.min(similarity, 1.0);
  }

  /**
   * Check if components have similar names
   */
  private haveSimilarNames(name1: string, name2: string): boolean {
    // Normalize names by removing common suffixes/prefixes
    const normalize = (name: string) => {
      return name
        .replace(/(_new|_old|_working|_Fixed|_replacement|_temp|_backup|_copy|\d+)$/i, '')
        .replace(/^(new|old|temp|backup|copy)_/i, '');
    };

    const norm1 = normalize(name1);
    const norm2 = normalize(name2);

    // Exact match after normalization
    if (norm1 === norm2) return true;

    // Case-insensitive match
    if (norm1.toLowerCase() === norm2.toLowerCase()) return true;

    // Original exact match
    if (name1 === name2) return true;

    // Case-insensitive match of originals
    if (name1.toLowerCase() === name2.toLowerCase()) return true;

    // Check if one contains the other
    if (name1.includes(name2) || name2.includes(name1)) return true;

    // Levenshtein distance for fuzzy matching
    const distance = this.levenshteinDistance(norm1, norm2);
    const maxLength = Math.max(norm1.length, norm2.length);
    const similarity = 1 - distance / maxLength;

    return similarity > 0.7;
  }

  /**
   * Calculate Levenshtein distance between two strings
   */
  private levenshteinDistance(str1: string, str2: string): number {
    const matrix: number[][] = [];

    for (let i = 0; i <= str2.length; i++) {
      matrix[i] = [i];
    }

    for (let j = 0; j <= str1.length; j++) {
      matrix[0][j] = j;
    }

    for (let i = 1; i <= str2.length; i++) {
      for (let j = 1; j <= str1.length; j++) {
        if (str2.charAt(i - 1) === str1.charAt(j - 1)) {
          matrix[i][j] = matrix[i - 1][j - 1];
        } else {
          matrix[i][j] = Math.min(
            matrix[i - 1][j - 1] + 1, // substitution
            matrix[i][j - 1] + 1,     // insertion
            matrix[i - 1][j] + 1      // deletion
          );
        }
      }
    }

    return matrix[str2.length][str1.length];
  }

  /**
   * Check if components have similar props
   */
  private haveSimilarProps(comp1: Component, comp2: Component): boolean {
    if (comp1.props.length === 0 && comp2.props.length === 0) return true;
    if (comp1.props.length === 0 || comp2.props.length === 0) return false;

    const props1Names = comp1.props.map(p => p.name);
    const props2Names = comp2.props.map(p => p.name);

    const commonProps = props1Names.filter(p => props2Names.includes(p));
    const similarityRatio = (commonProps.length * 2) / (props1Names.length + props2Names.length);

    return similarityRatio > 0.5;
  }

  /**
   * Get shared dependencies between components
   */
  private getSharedDependencies(comp1: Component, comp2: Component): string[] {
    return comp1.dependencies.filter(dep => comp2.dependencies.includes(dep));
  }

  /**
   * Suggest which component to keep when merging
   */
  private suggestMergeTarget(components: Component[]): string {
    // Sort by priority criteria
    const sorted = components.sort((a, b) => {
      // Prefer non-deprecated
      if (a.deprecated !== b.deprecated) {
        return a.deprecated ? 1 : -1;
      }

      // Prefer most used
      if (a.usageCount !== b.usageCount) {
        return b.usageCount - a.usageCount;
      }

      // Prefer most recent
      return b.updatedAt.getTime() - a.updatedAt.getTime();
    });

    return sorted[0].id;
  }

  /**
   * Generate merge suggestions for duplicate groups
   */
  generateMergeSuggestions(group: DuplicateGroup): string[] {
    const suggestions: string[] = [];

    suggestions.push(`Keep component: ${group.suggestedMerge}`);
    suggestions.push(`Remove duplicates: ${group.components
      .filter(c => c.id !== group.suggestedMerge)
      .map(c => c.id)
      .join(', ')}`);

    if (group.similarity === 100) {
      suggestions.push('Components are identical - safe to merge immediately');
    } else if (group.similarity >= 90) {
      suggestions.push('Components are very similar - review minor differences before merging');
    } else {
      suggestions.push('Components are similar - careful review recommended before merging');
    }

    return suggestions;
  }
}