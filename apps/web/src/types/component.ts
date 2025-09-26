/**
 * Component Management System Types
 *
 * Type definitions for cross-platform component management,
 * design token validation, and accessibility compliance tracking.
 */

import { ValidationIssue } from './designToken';

// Component entity types
export interface Component {
  id: string;
  name: string;
  category: ComponentCategory;
  description: string;
  designSpecs: ComponentDesignSpecs;
  implementations: ComponentImplementation[];
  status: ComponentStatus;
  version: string;
  createdAt: string;
  updatedAt: string;
  createdBy: string;
  metadata?: ComponentMetadata;
}

export enum ComponentCategory {
  BUTTON = 'button',
  INPUT = 'input',
  CARD = 'card',
  MODAL = 'modal',
  NAVIGATION = 'navigation',
  DISPLAY = 'display',
  FEEDBACK = 'feedback',
  LAYOUT = 'layout'
}

export enum ComponentStatus {
  DRAFT = 'draft',
  IN_REVIEW = 'in_review',
  APPROVED = 'approved',
  IMPLEMENTED = 'implemented',
  DEPRECATED = 'deprecated'
}

export interface ComponentDesignSpecs {
  variants: ComponentVariant[];
  states: ComponentState[];
  sizes: ComponentSize[];
  designTokens: string[]; // References to design token IDs
  accessibility: AccessibilityRequirements;
  interactions: InteractionSpec[];
}

export interface ComponentVariant {
  name: string;
  description: string;
  properties: Record<string, any>;
  previewUrl?: string;
}

export interface ComponentState {
  name: string;
  description: string;
  trigger: string;
  visualChanges: string[];
}

export interface ComponentSize {
  name: string;
  dimensions: {
    width?: number;
    height?: number;
    minWidth?: number;
    minHeight?: number;
  };
  touchTarget?: number; // For mobile platforms
}

export interface AccessibilityRequirements {
  wcagLevel: 'A' | 'AA' | 'AAA';
  contrastRatio: number;
  keyboardNavigation: boolean;
  screenReaderSupport: boolean;
  ariaLabels: string[];
  focusManagement: boolean;
}

export interface InteractionSpec {
  type: 'hover' | 'focus' | 'active' | 'disabled' | 'loading';
  description: string;
  animation?: AnimationSpec;
  haptic?: boolean; // For mobile platforms
}

export interface AnimationSpec {
  duration: number;
  easing: string;
  properties: string[];
}

// Component implementation types
export interface ComponentImplementation {
  id: string;
  componentId: string;
  platform: Platform;
  status: ImplementationStatus;
  codeLocation: string;
  version: string;
  testCoverage?: number;
  lastValidated?: string;
  validationResults?: ValidationResult[];
  maintainer: string;
  dependencies: string[];
  buildStatus?: BuildStatus;
}

export enum Platform {
  WEB = 'web',
  IOS = 'ios',
  ANDROID = 'android'
}

export enum ImplementationStatus {
  NOT_STARTED = 'not_started',
  IN_PROGRESS = 'in_progress',
  COMPLETED = 'completed',
  TESTING = 'testing',
  VALIDATED = 'validated',
  NEEDS_UPDATE = 'needs_update'
}

export interface ValidationResult {
  type: 'design_token' | 'accessibility' | 'functionality' | 'performance';
  isValid: boolean;
  score: number;
  issues: ValidationIssue[];
  validatedAt: string;
}

export interface BuildStatus {
  success: boolean;
  buildTime: string;
  errors: string[];
  warnings: string[];
}

export interface ComponentMetadata {
  tags: string[];
  figmaUrl?: string;
  storybook?: {
    url: string;
    stories: string[];
  };
  documentation?: {
    url: string;
    sections: string[];
  };
  changelog: ChangelogEntry[];
}

export interface ChangelogEntry {
  version: string;
  changes: string[];
  breaking: boolean;
  date: string;
  author: string;
}

// API request/response types
export interface CreateComponentRequest {
  name: string;
  category: ComponentCategory;
  description: string;
  designSpecs: Partial<ComponentDesignSpecs>;
  metadata?: Partial<ComponentMetadata>;
}

export interface UpdateComponentRequest {
  name?: string;
  description?: string;
  designSpecs?: Partial<ComponentDesignSpecs>;
  status?: ComponentStatus;
  metadata?: Partial<ComponentMetadata>;
}

export interface ComponentListQuery {
  category?: ComponentCategory;
  status?: ComponentStatus;
  search?: string;
  platform?: Platform;
  page?: number;
  limit?: number;
  sort?: 'name' | 'created' | 'updated' | 'status';
  order?: 'asc' | 'desc';
}

export interface ComponentImplementationQuery {
  componentId?: string;
  platform?: Platform;
  status?: ImplementationStatus;
  search?: string;
  includeValidation?: boolean;
  page?: number;
  limit?: number;
}

export interface CreateComponentImplementationRequest {
  componentId: string;
  platform: Platform;
  codeLocation: string;
  maintainer: string;
  dependencies?: string[];
}

export interface UpdateComponentImplementationRequest {
  status?: ImplementationStatus;
  codeLocation?: string;
  testCoverage?: number;
  maintainer?: string;
  dependencies?: string[];
}

// Component validation types
export interface ComponentValidationRequest {
  componentId: string;
  platforms?: Platform[];
  validationTypes?: ValidationResult['type'][];
}

export interface ComponentValidationResponse {
  componentId: string;
  overallScore: number;
  platformResults: Record<Platform, PlatformValidationResult>;
  summary: ValidationSummary;
  recommendations: string[];
}

export interface PlatformValidationResult {
  platform: Platform;
  score: number;
  results: ValidationResult[];
  issues: ValidationIssue[];
  lastValidated: string;
}

export interface ValidationSummary {
  totalChecks: number;
  passedChecks: number;
  criticalIssues: number;
  warningIssues: number;
  infoIssues: number;
  consistencyScore: number;
}

// Component sync types
export interface ComponentSyncRequest {
  componentIds?: string[];
  platforms?: Platform[];
  forceSync?: boolean;
}

export interface ComponentSyncResponse {
  syncedComponents: number;
  updatedImplementations: number;
  validationResults: ComponentValidationResponse[];
  errors: SyncError[];
  syncTime: string;
}

export interface SyncError {
  componentId: string;
  platform: Platform;
  error: string;
  severity: 'error' | 'warning';
}

// Bulk operations
export interface BulkComponentUpdateRequest {
  componentIds: string[];
  updates: Partial<UpdateComponentRequest>;
  validateAfterUpdate?: boolean;
}

export interface BulkComponentUpdateResponse {
  updated: string[];
  failed: BulkUpdateError[];
  validationResults?: ComponentValidationResponse[];
}

export interface BulkUpdateError {
  componentId: string;
  error: string;
  details?: string;
}

// Component analysis types
export interface ComponentAnalyticsQuery {
  dateFrom?: string;
  dateTo?: string;
  platforms?: Platform[];
  categories?: ComponentCategory[];
  groupBy?: 'platform' | 'category' | 'status' | 'day' | 'week' | 'month';
}

export interface ComponentAnalyticsResponse {
  totalComponents: number;
  implementationCoverage: Record<Platform, number>;
  statusDistribution: Record<ComponentStatus, number>;
  validationScores: Record<Platform, number>;
  trends: AnalyticsTrend[];
  topIssues: IssueFrequency[];
}

export interface AnalyticsTrend {
  period: string;
  value: number;
  metric: 'components_created' | 'implementations_completed' | 'validation_score';
}

export interface IssueFrequency {
  issue: string;
  count: number;
  platforms: Platform[];
  severity: ValidationIssue['severity'];
}

// Type guards
export function isComponent(data: unknown): data is Component {
  return (
    typeof data === 'object' &&
    data !== null &&
    'id' in data &&
    'name' in data &&
    'category' in data &&
    'status' in data
  );
}

export function isComponentImplementation(data: unknown): data is ComponentImplementation {
  return (
    typeof data === 'object' &&
    data !== null &&
    'id' in data &&
    'componentId' in data &&
    'platform' in data &&
    'status' in data
  );
}

export function isValidationResult(data: unknown): data is ValidationResult {
  return (
    typeof data === 'object' &&
    data !== null &&
    'type' in data &&
    'isValid' in data &&
    'score' in data
  );
}