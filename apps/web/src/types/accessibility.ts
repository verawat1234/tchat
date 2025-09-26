/**
 * Accessibility System Types
 *
 * Type definitions for WCAG 2.1 AA compliance validation,
 * cross-platform accessibility testing, and audit reporting.
 */

// Accessibility audit core types
export interface AccessibilityAudit {
  id: string;
  name: string;
  description: string;
  scope: AuditScope;
  standards: AccessibilityStandard[];
  results: AuditResult[];
  status: AuditStatus;
  createdAt: string;
  updatedAt: string;
  completedAt?: string;
  createdBy: string;
  metadata?: AuditMetadata;
}

export enum AuditScope {
  COMPONENT = 'component',
  PAGE = 'page',
  APPLICATION = 'application',
  DESIGN_SYSTEM = 'design_system'
}

export enum AccessibilityStandard {
  WCAG_2_1_A = 'wcag_2_1_a',
  WCAG_2_1_AA = 'wcag_2_1_aa',
  WCAG_2_1_AAA = 'wcag_2_1_aaa',
  SECTION_508 = 'section_508',
  EN_301_549 = 'en_301_549',
  ADA = 'ada'
}

export enum AuditStatus {
  SCHEDULED = 'scheduled',
  IN_PROGRESS = 'in_progress',
  COMPLETED = 'completed',
  FAILED = 'failed',
  CANCELLED = 'cancelled'
}

export interface AuditResult {
  id: string;
  auditId: string;
  platform: Platform;
  component?: string;
  url?: string;
  timestamp: string;
  overallScore: number; // 0-100
  passedChecks: number;
  failedChecks: number;
  totalChecks: number;
  findings: AccessibilityFinding[];
  recommendations: Recommendation[];
  compliance: ComplianceResult;
}

export enum Platform {
  WEB = 'web',
  IOS = 'ios',
  ANDROID = 'android',
  DESKTOP = 'desktop'
}

export interface AccessibilityFinding {
  id: string;
  ruleId: string;
  ruleName: string;
  category: FindingCategory;
  severity: FindingSeverity;
  impact: ImpactLevel;
  wcagCriteria: WCAGCriteria[];
  description: string;
  explanation: string;
  element?: ElementInfo;
  location: LocationInfo;
  evidence: Evidence;
  fixes: Fix[];
  isBlocking: boolean; // Blocks release/deployment
  affectedUsers: AffectedUserGroups;
}

export enum FindingCategory {
  PERCEIVABLE = 'perceivable',
  OPERABLE = 'operable',
  UNDERSTANDABLE = 'understandable',
  ROBUST = 'robust',
  KEYBOARD = 'keyboard',
  FOCUS = 'focus',
  COLOR_CONTRAST = 'color_contrast',
  TEXT_ALTERNATIVES = 'text_alternatives',
  ARIA = 'aria',
  FORMS = 'forms',
  NAVIGATION = 'navigation',
  STRUCTURE = 'structure'
}

export enum FindingSeverity {
  CRITICAL = 'critical',    // WCAG Level A violation
  MAJOR = 'major',         // WCAG Level AA violation
  MODERATE = 'moderate',   // WCAG Level AAA violation
  MINOR = 'minor',         // Best practice violation
  INFO = 'info'            // Informational
}

export enum ImpactLevel {
  BLOCKER = 'blocker',     // Prevents feature use entirely
  CRITICAL = 'critical',   // Severely impacts usability
  MAJOR = 'major',        // Significantly impacts usability
  MODERATE = 'moderate',   // Moderately impacts usability
  MINOR = 'minor',        // Slightly impacts usability
  COSMETIC = 'cosmetic'   // Visual/aesthetic issue only
}

export interface WCAGCriteria {
  criterion: string;       // e.g., "1.4.3"
  level: 'A' | 'AA' | 'AAA';
  title: string;
  description: string;
  techniques: string[];    // WCAG technique IDs
}

export interface ElementInfo {
  tagName: string;
  id?: string;
  className?: string;
  ariaRole?: string;
  ariaLabel?: string;
  text?: string;
  xpath?: string;
  cssSelector?: string;
  bounds?: ElementBounds;
}

export interface ElementBounds {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface LocationInfo {
  platform: Platform;
  url?: string;
  screenName?: string;
  componentName?: string;
  viewHierarchy?: string[];
  coordinates?: {
    x: number;
    y: number;
  };
}

export interface Evidence {
  screenshots?: Screenshot[];
  htmlSnippet?: string;
  computedStyles?: Record<string, string>;
  ariaTree?: AriaNode[];
  contrastRatio?: number;
  expectedContrastRatio?: number;
  colorValues?: {
    foreground: string;
    background: string;
  };
}

export interface Screenshot {
  url: string;
  description: string;
  highlightedElements?: ElementBounds[];
  annotated: boolean;
}

export interface AriaNode {
  role: string;
  name?: string;
  description?: string;
  states: string[];
  properties: Record<string, string>;
  children?: AriaNode[];
}

export interface Fix {
  id: string;
  type: FixType;
  priority: 'immediate' | 'high' | 'medium' | 'low';
  effort: 'low' | 'medium' | 'high';
  description: string;
  implementation: ImplementationGuide;
  validation: ValidationSteps;
  beforeCode?: string;
  afterCode?: string;
  designChanges?: DesignChange[];
}

export enum FixType {
  CODE_CHANGE = 'code_change',
  DESIGN_CHANGE = 'design_change',
  CONTENT_CHANGE = 'content_change',
  CONFIGURATION = 'configuration',
  TESTING_REQUIRED = 'testing_required'
}

export interface ImplementationGuide {
  steps: string[];
  codeExamples?: CodeExample[];
  platformSpecific?: Record<Platform, string[]>;
  dependencies?: string[];
  breakingChanges?: boolean;
}

export interface CodeExample {
  language: string;
  before: string;
  after: string;
  explanation: string;
}

export interface ValidationSteps {
  automated?: string[];
  manual?: string[];
  userTesting?: string[];
  tools?: string[];
}

export interface DesignChange {
  property: string;
  currentValue: string;
  recommendedValue: string;
  reason: string;
  affectedComponents?: string[];
}

export enum AffectedUserGroups {
  BLIND = 'blind',
  LOW_VISION = 'low_vision',
  MOTOR_IMPAIRED = 'motor_impaired',
  COGNITIVE_IMPAIRED = 'cognitive_impaired',
  DEAF = 'deaf',
  HARD_OF_HEARING = 'hard_of_hearing',
  KEYBOARD_ONLY = 'keyboard_only',
  SCREEN_READER = 'screen_reader',
  VOICE_CONTROL = 'voice_control'
}

export interface Recommendation {
  id: string;
  category: RecommendationCategory;
  priority: 'critical' | 'high' | 'medium' | 'low';
  title: string;
  description: string;
  benefits: string[];
  implementation: ImplementationGuide;
  estimatedEffort: string;
  roi: string; // Return on investment description
  relatedFindings: string[]; // Finding IDs this addresses
}

export enum RecommendationCategory {
  QUICK_WIN = 'quick_win',
  SYSTEMATIC_IMPROVEMENT = 'systematic_improvement',
  PROCESS_CHANGE = 'process_change',
  TRAINING = 'training',
  TOOLING = 'tooling',
  DESIGN_SYSTEM = 'design_system'
}

export interface ComplianceResult {
  standard: AccessibilityStandard;
  level: 'A' | 'AA' | 'AAA' | 'FAIL';
  score: number; // 0-100
  passedCriteria: WCAGCriteria[];
  failedCriteria: WCAGCriteria[];
  partialCriteria: WCAGCriteria[];
  notApplicable: WCAGCriteria[];
  compliancePercentage: number;
  blockers: string[]; // Critical issues preventing compliance
}

export interface AuditMetadata {
  tools: AuditTool[];
  environment: TestEnvironment;
  testConfiguration: TestConfiguration;
  coverage: CoverageInfo;
  executionTime: number; // seconds
  dataRetention: string; // ISO date
}

export interface AuditTool {
  name: string;
  version: string;
  type: 'automated' | 'manual' | 'hybrid';
  configuration?: Record<string, any>;
}

export interface TestEnvironment {
  platform: Platform;
  browser?: BrowserInfo;
  device?: DeviceInfo;
  assistiveTechnology?: AssistiveTechInfo[];
  viewport?: ViewportInfo;
}

export interface BrowserInfo {
  name: string;
  version: string;
  userAgent?: string;
}

export interface DeviceInfo {
  type: 'desktop' | 'tablet' | 'mobile' | 'tv';
  name?: string;
  screenSize?: {
    width: number;
    height: number;
  };
  pixelDensity?: number;
}

export interface AssistiveTechInfo {
  type: 'screen_reader' | 'magnifier' | 'voice_control' | 'switch_device';
  name: string;
  version?: string;
}

export interface ViewportInfo {
  width: number;
  height: number;
  pixelRatio: number;
}

export interface TestConfiguration {
  includedRules: string[];
  excludedRules: string[];
  customRules?: CustomRule[];
  thresholds: QualityThresholds;
}

export interface CustomRule {
  id: string;
  name: string;
  description: string;
  selector: string;
  check: string; // Function or expression to evaluate
  severity: FindingSeverity;
}

export interface QualityThresholds {
  minimumScore: number;
  maximumCriticalFindings: number;
  maximumMajorFindings: number;
  contrastRatioThreshold: number;
  keyboardNavigationCoverage: number;
  screenReaderCoverage: number;
}

export interface CoverageInfo {
  totalElements: number;
  testedElements: number;
  interactiveElements: number;
  testedInteractive: number;
  coveragePercentage: number;
  skippedElements: SkippedElement[];
}

export interface SkippedElement {
  selector: string;
  reason: 'hidden' | 'disabled' | 'excluded' | 'error';
  description?: string;
}

// API request/response types
export interface CreateAuditRequest {
  name: string;
  description: string;
  scope: AuditScope;
  standards: AccessibilityStandard[];
  targets: AuditTarget[];
  configuration?: Partial<TestConfiguration>;
  scheduledFor?: string;
}

export interface AuditTarget {
  platform: Platform;
  identifier: string; // URL, component name, screen identifier
  type: 'component' | 'page' | 'flow';
  metadata?: Record<string, any>;
}

export interface AuditListQuery {
  scope?: AuditScope;
  status?: AuditStatus;
  standard?: AccessibilityStandard;
  platform?: Platform;
  search?: string;
  createdBy?: string;
  dateFrom?: string;
  dateTo?: string;
  page?: number;
  limit?: number;
  sort?: 'created' | 'updated' | 'score' | 'name';
  order?: 'asc' | 'desc';
}

export interface BulkAuditRequest {
  auditIds: string[];
  action: 'execute' | 'cancel' | 'delete' | 'archive';
  parameters?: Record<string, any>;
}

export interface BulkAuditResponse {
  processed: number;
  successful: string[];
  failed: AuditError[];
}

export interface AuditError {
  auditId: string;
  error: string;
  details?: string;
}

// Accessibility monitoring types
export interface AccessibilityMonitoringConfig {
  enabled: boolean;
  frequency: 'daily' | 'weekly' | 'monthly';
  standards: AccessibilityStandard[];
  platforms: Platform[];
  thresholds: QualityThresholds;
  notifications: NotificationConfig;
  autoRemediation: boolean;
}

export interface NotificationConfig {
  email: boolean;
  slack?: {
    webhook: string;
    channel: string;
  };
  teams?: {
    webhook: string;
  };
  onlyRegressions: boolean;
  severityThreshold: FindingSeverity;
}

export interface AccessibilityReport {
  id: string;
  reportType: ReportType;
  period: ReportPeriod;
  data: ReportData;
  generatedAt: string;
  format: 'html' | 'pdf' | 'json' | 'csv';
  url?: string;
}

export enum ReportType {
  COMPLIANCE_STATUS = 'compliance_status',
  TREND_ANALYSIS = 'trend_analysis',
  FINDINGS_SUMMARY = 'findings_summary',
  EXECUTIVE_SUMMARY = 'executive_summary',
  DETAILED_AUDIT = 'detailed_audit'
}

export interface ReportPeriod {
  from: string;
  to: string;
  interval: 'daily' | 'weekly' | 'monthly' | 'quarterly' | 'custom';
}

export interface ReportData {
  summary: ReportSummary;
  trends?: TrendData[];
  findings?: AccessibilityFinding[];
  compliance?: ComplianceResult[];
  recommendations?: Recommendation[];
}

export interface ReportSummary {
  totalAudits: number;
  averageScore: number;
  complianceRate: number;
  criticalFindings: number;
  resolvedFindings: number;
  platformBreakdown: Record<Platform, PlatformSummary>;
}

export interface PlatformSummary {
  audits: number;
  score: number;
  compliance: number;
  findings: number;
}

export interface TrendData {
  period: string;
  metric: string;
  value: number;
  change?: number; // Percentage change from previous period
}

// Type guards
export function isAccessibilityAudit(data: unknown): data is AccessibilityAudit {
  return (
    typeof data === 'object' &&
    data !== null &&
    'id' in data &&
    'scope' in data &&
    'standards' in data &&
    'status' in data
  );
}

export function isAccessibilityFinding(data: unknown): data is AccessibilityFinding {
  return (
    typeof data === 'object' &&
    data !== null &&
    'ruleId' in data &&
    'severity' in data &&
    'category' in data &&
    'wcagCriteria' in data
  );
}

export function isFix(data: unknown): data is Fix {
  return (
    typeof data === 'object' &&
    data !== null &&
    'type' in data &&
    'description' in data &&
    'implementation' in data
  );
}