/**
 * Configuration schema for the component analyzer
 */

export interface ComponentAnalyzerConfig {
  /** Paths configuration */
  paths: PathsConfig;

  /** Analysis configuration */
  analysis: AnalysisConfig;

  /** Validation configuration */
  validation: ValidationConfig;

  /** Duplicate detection configuration */
  duplicates: DuplicatesConfig;

  /** Output configuration */
  output: OutputConfig;
}

export interface PathsConfig {
  /** Path to components directory */
  components: string;

  /** Path to output directory */
  output: string;

  /** Path to registry file */
  registry: string;
}

export interface AnalysisConfig {
  /** Include subdirectories in analysis */
  recursive: boolean;

  /** Maximum directory depth to analyze */
  maxDepth: number;

  /** Glob patterns to exclude from analysis */
  exclude: string[];

  /** Include test files in analysis */
  includeTests: boolean;

  /** Component categorization patterns */
  categorization?: CategorizationPatterns;
}

export interface CategorizationPatterns {
  /** Patterns for identifying atom components */
  atomPatterns?: string[];

  /** Patterns for identifying molecule components */
  moleculePatterns?: string[];

  /** Patterns for identifying organism components */
  organismPatterns?: string[];
}

export interface ValidationConfig {
  /** Validation rules to apply */
  rules: ValidationRules;

  /** Automatically fix issues when possible */
  autoFix: boolean;

  /** Minimum severity level to report */
  severity: 'error' | 'warning' | 'info';
}

export interface ValidationRules {
  /** Naming convention validation */
  naming?: {
    enabled: boolean;
    pattern: string;
  };

  /** Accessibility validation */
  accessibility?: {
    enabled: boolean;
    wcagLevel: 'A' | 'AA' | 'AAA';
  };

  /** Documentation validation */
  documentation?: {
    enabled: boolean;
    requireDescription: boolean;
    requireExamples: boolean;
  };

  /** Structure validation */
  structure?: {
    enabled: boolean;
    requireProps: boolean;
    requireTypes: boolean;
  };

  /** Styling validation */
  styling?: {
    enabled: boolean;
    requireConsistentStyles: boolean;
  };

  /** Performance validation */
  performance?: {
    enabled: boolean;
    maxComplexity: number;
  };
}

export interface DuplicatesConfig {
  /** Similarity threshold percentage (0-100) */
  threshold: number;

  /** Whether to automatically suggest merges */
  autoMerge: boolean;

  /** Factors to consider in duplicate detection */
  factors?: {
    /** Weight for structural similarity */
    structural?: number;

    /** Weight for visual similarity */
    visual?: number;

    /** Weight for functional similarity */
    functional?: number;
  };
}

export interface OutputConfig {
  /** Output format */
  format: 'json' | 'markdown' | 'both' | 'html';

  /** Include visual references in output */
  includeVisuals: boolean;

  /** Include usage examples in output */
  includeExamples: boolean;

  /** Use colors in terminal output */
  colors: boolean;

  /** Verbose output */
  verbose: boolean;
}

/** Default configuration */
export const DEFAULT_CONFIG: ComponentAnalyzerConfig = {
  paths: {
    components: 'src/components',
    output: 'docs/components',
    registry: 'docs/components/registry.json'
  },
  analysis: {
    recursive: true,
    maxDepth: 10,
    exclude: ['*.test.tsx', '*.test.ts', '*.stories.tsx', '*.stories.ts', '__mocks__/*'],
    includeTests: false
  },
  validation: {
    rules: {
      naming: {
        enabled: true,
        pattern: '^[A-Z][a-zA-Z]*$'
      },
      accessibility: {
        enabled: true,
        wcagLevel: 'AA'
      },
      documentation: {
        enabled: true,
        requireDescription: true,
        requireExamples: false
      },
      structure: {
        enabled: true,
        requireProps: true,
        requireTypes: true
      },
      styling: {
        enabled: true,
        requireConsistentStyles: true
      },
      performance: {
        enabled: false,
        maxComplexity: 10
      }
    },
    autoFix: false,
    severity: 'warning'
  },
  duplicates: {
    threshold: 75,
    autoMerge: false,
    factors: {
      structural: 0.4,
      visual: 0.3,
      functional: 0.3
    }
  },
  output: {
    format: 'both',
    includeVisuals: true,
    includeExamples: true,
    colors: true,
    verbose: false
  }
};