/**
 * Design Token System Types
 *
 * Comprehensive type definitions for design token management,
 * cross-platform validation, and OKLCH color accuracy.
 */

// Design token core types
export interface DesignToken {
  id: string;
  name: string;
  category: TokenCategory;
  type: TokenType;
  value: TokenValue;
  platformMappings: Record<Platform, string>;
  description: string;
  tags: string[];
  status: TokenStatus;
  version: string;
  createdAt: string;
  updatedAt: string;
  createdBy: string;
  metadata?: TokenMetadata;
}

export enum TokenCategory {
  COLOR = 'color',
  SPACING = 'spacing',
  TYPOGRAPHY = 'typography',
  BORDER_RADIUS = 'border_radius',
  ELEVATION = 'elevation',
  ANIMATION = 'animation',
  BREAKPOINT = 'breakpoint'
}

export enum TokenType {
  // Color types
  PRIMARY_COLOR = 'primary_color',
  SECONDARY_COLOR = 'secondary_color',
  SUCCESS_COLOR = 'success_color',
  WARNING_COLOR = 'warning_color',
  ERROR_COLOR = 'error_color',
  NEUTRAL_COLOR = 'neutral_color',
  TEXT_COLOR = 'text_color',
  BACKGROUND_COLOR = 'background_color',
  BORDER_COLOR = 'border_color',

  // Spacing types
  MARGIN = 'margin',
  PADDING = 'padding',
  GAP = 'gap',

  // Typography types
  FONT_SIZE = 'font_size',
  FONT_WEIGHT = 'font_weight',
  LINE_HEIGHT = 'line_height',
  LETTER_SPACING = 'letter_spacing',
  FONT_FAMILY = 'font_family',

  // Other types
  CORNER_RADIUS = 'corner_radius',
  SHADOW = 'shadow',
  DURATION = 'duration',
  EASING = 'easing',
  VIEWPORT = 'viewport'
}

export enum TokenStatus {
  DRAFT = 'draft',
  ACTIVE = 'active',
  DEPRECATED = 'deprecated',
  ARCHIVED = 'archived'
}

export enum Platform {
  WEB = 'web',
  IOS = 'ios',
  ANDROID = 'android'
}

// Token value types
export type TokenValue =
  | ColorValue
  | SpacingValue
  | TypographyValue
  | RadiusValue
  | ElevationValue
  | AnimationValue
  | BreakpointValue;

export interface ColorValue {
  hex: string;
  oklch?: {
    l: number;  // Lightness (0-100)
    c: number;  // Chroma (0+)
    h: number;  // Hue (0-360)
    alpha?: number; // Alpha (0-1)
  };
  rgb?: {
    r: number;
    g: number;
    b: number;
    a?: number;
  };
  hsl?: {
    h: number;
    s: number;
    l: number;
    a?: number;
  };
}

export interface SpacingValue {
  base: number; // Base value in pixels
  rem?: number; // Rem equivalent
  em?: number;  // Em equivalent
  units: 'px' | 'rem' | 'em' | 'dp' | 'pt';
}

export interface TypographyValue {
  size: number;
  weight: number | string;
  lineHeight?: number | string;
  letterSpacing?: number;
  fontFamily?: string;
  units: 'px' | 'rem' | 'em' | 'sp' | 'pt';
}

export interface RadiusValue {
  value: number;
  units: 'px' | 'rem' | 'dp' | 'pt';
  isInfinite?: boolean; // For fully rounded corners
}

export interface ElevationValue {
  level: number; // 0-24 Material Design elevation
  shadow: {
    offsetX: number;
    offsetY: number;
    blurRadius: number;
    spreadRadius?: number;
    color: string;
  };
}

export interface AnimationValue {
  duration: number; // in milliseconds
  easing: string;   // CSS easing function or platform equivalent
}

export interface BreakpointValue {
  minWidth?: number;
  maxWidth?: number;
  units: 'px' | 'rem' | 'em';
}

export interface TokenMetadata {
  figmaId?: string;
  usage: string[];
  relatedTokens: string[];
  deprecationNote?: string;
  migrationPath?: string;
  designDecision?: string;
  accessibility?: {
    contrastRatio?: number;
    wcagLevel?: 'A' | 'AA' | 'AAA';
    requirements?: string[];
  };
}

// Token validation types
export interface TokenValidationRequest {
  tokenId?: string;
  tokenIds?: string[];
  platforms?: Platform[];
  validationType?: ValidationType[];
  consistencyThreshold?: number; // Default: 0.97 (97%)
}

export enum ValidationType {
  CROSS_PLATFORM_CONSISTENCY = 'cross_platform_consistency',
  OKLCH_ACCURACY = 'oklch_accuracy',
  BASE_UNIT_SYSTEM = 'base_unit_system',
  ACCESSIBILITY_COMPLIANCE = 'accessibility_compliance',
  DESIGN_SYSTEM_RULES = 'design_system_rules'
}

export interface TokenValidationResponse {
  tokenId: string;
  isValid: boolean;
  overallScore: number;
  results: TokenValidationResult[];
  issues: ValidationIssue[];
  recommendations: string[];
  validatedAt: string;
}

export interface TokenValidationResult {
  type: ValidationType;
  isValid: boolean;
  score: number;
  details: ValidationDetails;
  issues: ValidationIssue[];
}

export type ValidationDetails =
  | ConsistencyDetails
  | OKLCHAccuracyDetails
  | BaseUnitDetails
  | AccessibilityDetails
  | DesignSystemDetails;

export interface ConsistencyDetails {
  platformValues: Record<Platform, string>;
  normalizedValues: Record<Platform, string>;
  consistencyScore: number;
  consistentPlatforms: Platform[];
  inconsistentPlatforms: Platform[];
}

export interface OKLCHAccuracyDetails {
  oklchValue: string;
  hexValue: string;
  conversionAccuracy: number;
  colorDifference: number; // Delta E value
  tolerance: number;
}

export interface BaseUnitDetails {
  baseUnit: number;
  actualValue: number;
  isMultiple: boolean;
  nearestValidValue: number;
  deviation: number;
}

export interface AccessibilityDetails {
  contrastRatio: number;
  wcagLevel: 'A' | 'AA' | 'AAA' | 'FAIL';
  meetsMinimum: boolean;
  targetRatio: number;
  testedAgainst: string[]; // Background colors tested against
}

export interface DesignSystemDetails {
  rule: string;
  ruleDescription: string;
  compliance: boolean;
  actualValue: any;
  expectedValue: any;
  suggestion: string;
}

export interface ValidationIssue {
  id: string;
  severity: 'error' | 'warning' | 'info';
  type: ValidationType;
  platform?: Platform;
  title: string;
  message: string;
  description: string;
  fix?: string;
  impactScore: number; // 0-10 scale
  effort: 'low' | 'medium' | 'high';
}

// Token API request/response types
export interface CreateTokenRequest {
  name: string;
  category: TokenCategory;
  type: TokenType;
  value: TokenValue;
  platformMappings: Record<Platform, string>;
  description: string;
  tags?: string[];
  metadata?: Partial<TokenMetadata>;
}

export interface UpdateTokenRequest {
  name?: string;
  value?: TokenValue;
  platformMappings?: Partial<Record<Platform, string>>;
  description?: string;
  tags?: string[];
  status?: TokenStatus;
  metadata?: Partial<TokenMetadata>;
}

export interface TokenListQuery {
  category?: TokenCategory;
  type?: TokenType;
  status?: TokenStatus;
  platform?: Platform;
  search?: string;
  tags?: string[];
  page?: number;
  limit?: number;
  sort?: 'name' | 'category' | 'created' | 'updated';
  order?: 'asc' | 'desc';
  includeDeprecated?: boolean;
}

export interface BulkTokenValidationRequest {
  tokenIds: string[];
  platforms?: Platform[];
  validationTypes?: ValidationType[];
  consistencyThreshold?: number;
  includeRecommendations?: boolean;
}

export interface BulkTokenValidationResponse {
  totalTokens: number;
  validTokens: number;
  invalidTokens: number;
  overallScore: number;
  results: TokenValidationResponse[];
  summary: ValidationSummary;
  prioritizedIssues: ValidationIssue[];
}

export interface ValidationSummary {
  criticalIssues: number;
  warnings: number;
  infoItems: number;
  platformConsistency: Record<Platform, number>;
  categoryScores: Record<TokenCategory, number>;
  constitutionalCompliance: boolean; // 97%+ consistency requirement
}

// Token sync types
export interface TokenSyncRequest {
  tokenIds?: string[];
  platforms?: Platform[];
  sourceOfTruth?: Platform | 'figma' | 'manual';
  validateAfterSync?: boolean;
}

export interface TokenSyncResponse {
  syncedTokens: number;
  updatedPlatforms: Platform[];
  conflicts: TokenConflict[];
  validationResults: TokenValidationResponse[];
  syncTime: string;
}

export interface TokenConflict {
  tokenId: string;
  platform: Platform;
  currentValue: string;
  proposedValue: string;
  confidence: number;
  resolution: 'auto' | 'manual' | 'skip';
}

// Token analytics types
export interface TokenAnalyticsQuery {
  dateFrom?: string;
  dateTo?: string;
  categories?: TokenCategory[];
  platforms?: Platform[];
  groupBy?: 'category' | 'platform' | 'status' | 'day' | 'week' | 'month';
}

export interface TokenAnalyticsResponse {
  totalTokens: number;
  activeTokens: number;
  consistencyScore: number;
  platformCoverage: Record<Platform, number>;
  categoryDistribution: Record<TokenCategory, number>;
  validationTrends: AnalyticsTrend[];
  topIssues: IssueFrequency[];
  migrationProgress?: MigrationProgress;
}

export interface AnalyticsTrend {
  period: string;
  value: number;
  metric: 'tokens_created' | 'validation_score' | 'issues_resolved' | 'consistency_score';
}

export interface IssueFrequency {
  issue: string;
  count: number;
  severity: ValidationIssue['severity'];
  platforms: Platform[];
  trend: 'increasing' | 'stable' | 'decreasing';
}

export interface MigrationProgress {
  totalTokens: number;
  migratedTokens: number;
  pendingTokens: number;
  deprecatedTokens: number;
  progress: number; // 0-100 percentage
  estimatedCompletion?: string;
}

// Export validation configuration
export interface TokenValidationConfig {
  consistencyThreshold: number;
  oklchTolerance: number;
  baseUnitSize: number;
  accessibilityStandard: 'WCAG2_A' | 'WCAG2_AA' | 'WCAG2_AAA';
  platforms: Platform[];
  validationRules: ValidationRule[];
}

export interface ValidationRule {
  id: string;
  name: string;
  description: string;
  category: TokenCategory;
  severity: ValidationIssue['severity'];
  enabled: boolean;
  validator: (token: DesignToken) => TokenValidationResult;
}

// Type guards
export function isDesignToken(data: unknown): data is DesignToken {
  return (
    typeof data === 'object' &&
    data !== null &&
    'id' in data &&
    'name' in data &&
    'category' in data &&
    'value' in data &&
    'platformMappings' in data
  );
}

export function isColorValue(value: TokenValue): value is ColorValue {
  return typeof value === 'object' && 'hex' in value;
}

export function isSpacingValue(value: TokenValue): value is SpacingValue {
  return typeof value === 'object' && 'base' in value && 'units' in value;
}

export function isValidationIssue(data: unknown): data is ValidationIssue {
  return (
    typeof data === 'object' &&
    data !== null &&
    'severity' in data &&
    'message' in data &&
    'type' in data
  );
}