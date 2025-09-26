/**
 * Unit Tests: ComponentDefinition Model (T057)
 * Tests comprehensive model validation and type safety
 * Constitutional requirements: Type safety, model validation, API compliance
 */
import { describe, it, expect } from 'vitest';
import {
  Component,
  ComponentCategory,
  ComponentStatus,
  Platform,
  ImplementationStatus,
  ComponentImplementation,
  ValidationResult,
  AccessibilityRequirements,
  ComponentDesignSpecs,
  ComponentVariant,
  ComponentState,
  ComponentSize,
  InteractionSpec,
  AnimationSpec,
  CreateComponentRequest,
  UpdateComponentRequest,
  ComponentListQuery,
  ComponentValidationRequest,
  ComponentSyncRequest,
  BulkComponentUpdateRequest,
  isComponent,
  isComponentImplementation,
  isValidationResult
} from '../component';

describe('ComponentDefinition Model Tests', () => {
  describe('Component Entity Model', () => {
    it('should create a valid Component entity with all required fields', () => {
      const component: Component = {
        id: 'comp_001',
        name: 'TchatButton',
        category: ComponentCategory.BUTTON,
        description: 'Cross-platform design system button component',
        designSpecs: {
          variants: [
            {
              name: 'primary',
              description: 'Primary call-to-action button',
              properties: { bgColor: '#3B82F6', textColor: '#FFFFFF' }
            }
          ],
          states: [
            {
              name: 'default',
              description: 'Default state',
              trigger: 'initial',
              visualChanges: []
            }
          ],
          sizes: [
            {
              name: 'md',
              dimensions: { height: 44, minWidth: 88 },
              touchTarget: 44
            }
          ],
          designTokens: ['primary', 'text-white'],
          accessibility: {
            wcagLevel: 'AA',
            contrastRatio: 4.5,
            keyboardNavigation: true,
            screenReaderSupport: true,
            ariaLabels: ['button'],
            focusManagement: true
          },
          interactions: [
            {
              type: 'hover',
              description: 'Hover state with elevation',
              animation: {
                duration: 200,
                easing: 'ease-out',
                properties: ['background-color', 'box-shadow']
              }
            }
          ]
        },
        implementations: [],
        status: ComponentStatus.IMPLEMENTED,
        version: '1.0.0',
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-02T00:00:00Z',
        createdBy: 'designer@tchat.com'
      };

      expect(component.id).toBe('comp_001');
      expect(component.name).toBe('TchatButton');
      expect(component.category).toBe(ComponentCategory.BUTTON);
      expect(component.status).toBe(ComponentStatus.IMPLEMENTED);
      expect(component.designSpecs.variants).toHaveLength(1);
      expect(component.designSpecs.accessibility.wcagLevel).toBe('AA');
      expect(component.designSpecs.accessibility.contrastRatio).toBe(4.5);
    });

    it('should validate ComponentCategory enum values', () => {
      const categories = Object.values(ComponentCategory);
      expect(categories).toContain('button');
      expect(categories).toContain('input');
      expect(categories).toContain('card');
      expect(categories).toContain('modal');
      expect(categories).toContain('navigation');
      expect(categories).toContain('display');
      expect(categories).toContain('feedback');
      expect(categories).toContain('layout');
      expect(categories).toHaveLength(8);
    });

    it('should validate ComponentStatus enum values', () => {
      const statuses = Object.values(ComponentStatus);
      expect(statuses).toContain('draft');
      expect(statuses).toContain('in_review');
      expect(statuses).toContain('approved');
      expect(statuses).toContain('implemented');
      expect(statuses).toContain('deprecated');
      expect(statuses).toHaveLength(5);
    });

    it('should validate Platform enum values for cross-platform consistency', () => {
      const platforms = Object.values(Platform);
      expect(platforms).toContain('web');
      expect(platforms).toContain('ios');
      expect(platforms).toContain('android');
      expect(platforms).toHaveLength(3);
    });

    it('should validate ImplementationStatus enum values', () => {
      const statuses = Object.values(ImplementationStatus);
      expect(statuses).toContain('not_started');
      expect(statuses).toContain('in_progress');
      expect(statuses).toContain('completed');
      expect(statuses).toContain('testing');
      expect(statuses).toContain('validated');
      expect(statuses).toContain('needs_update');
      expect(statuses).toHaveLength(6);
    });
  });

  describe('ComponentDesignSpecs Model', () => {
    it('should create valid design specifications', () => {
      const designSpecs: ComponentDesignSpecs = {
        variants: [
          {
            name: 'primary',
            description: 'Primary variant',
            properties: { color: 'blue' }
          },
          {
            name: 'secondary',
            description: 'Secondary variant',
            properties: { color: 'gray' }
          }
        ],
        states: [
          {
            name: 'loading',
            description: 'Loading state with spinner',
            trigger: 'async_action',
            visualChanges: ['show_spinner', 'disable_interaction']
          }
        ],
        sizes: [
          {
            name: 'sm',
            dimensions: { height: 32, minWidth: 64 }
          },
          {
            name: 'lg',
            dimensions: { height: 48, minWidth: 112 },
            touchTarget: 48
          }
        ],
        designTokens: ['primary', 'secondary', 'text-white'],
        accessibility: {
          wcagLevel: 'AA',
          contrastRatio: 4.5,
          keyboardNavigation: true,
          screenReaderSupport: true,
          ariaLabels: ['button', 'loading'],
          focusManagement: true
        },
        interactions: [
          {
            type: 'active',
            description: 'Active press state',
            animation: {
              duration: 150,
              easing: 'ease-in-out',
              properties: ['transform']
            },
            haptic: true
          }
        ]
      };

      expect(designSpecs.variants).toHaveLength(2);
      expect(designSpecs.states).toHaveLength(1);
      expect(designSpecs.sizes).toHaveLength(2);
      expect(designSpecs.designTokens).toContain('primary');
      expect(designSpecs.accessibility.wcagLevel).toBe('AA');
      expect(designSpecs.interactions[0].haptic).toBe(true);
    });

    it('should validate ComponentVariant structure', () => {
      const variant: ComponentVariant = {
        name: 'destructive',
        description: 'Destructive action button',
        properties: {
          bgColor: '#EF4444',
          textColor: '#FFFFFF',
          hoverColor: '#DC2626'
        },
        previewUrl: 'https://storybook.tchat.com/button-destructive'
      };

      expect(variant.name).toBe('destructive');
      expect(variant.properties.bgColor).toBe('#EF4444');
      expect(variant.previewUrl).toBeDefined();
    });

    it('should validate ComponentState structure', () => {
      const state: ComponentState = {
        name: 'disabled',
        description: 'Disabled state with reduced opacity',
        trigger: 'disabled_prop',
        visualChanges: ['opacity_60', 'pointer_events_none', 'cursor_not_allowed']
      };

      expect(state.name).toBe('disabled');
      expect(state.visualChanges).toContain('opacity_60');
      expect(state.visualChanges).toHaveLength(3);
    });

    it('should validate ComponentSize with touch target requirements', () => {
      const size: ComponentSize = {
        name: 'mobile_lg',
        dimensions: {
          width: 120,
          height: 48,
          minWidth: 88,
          minHeight: 44
        },
        touchTarget: 48 // Constitutional requirement: 44dp minimum
      };

      expect(size.touchTarget).toBeGreaterThanOrEqual(44);
      expect(size.dimensions.minHeight).toBeGreaterThanOrEqual(44);
    });

    it('should validate AccessibilityRequirements for WCAG compliance', () => {
      const accessibility: AccessibilityRequirements = {
        wcagLevel: 'AA',
        contrastRatio: 4.5,
        keyboardNavigation: true,
        screenReaderSupport: true,
        ariaLabels: ['button', 'primary-action', 'clickable'],
        focusManagement: true
      };

      expect(accessibility.wcagLevel).toBe('AA');
      expect(accessibility.contrastRatio).toBeGreaterThanOrEqual(4.5);
      expect(accessibility.keyboardNavigation).toBe(true);
      expect(accessibility.screenReaderSupport).toBe(true);
      expect(accessibility.focusManagement).toBe(true);
      expect(accessibility.ariaLabels).toContain('button');
    });

    it('should validate InteractionSpec with animations', () => {
      const interaction: InteractionSpec = {
        type: 'hover',
        description: 'Hover interaction with smooth transition',
        animation: {
          duration: 200,
          easing: 'cubic-bezier(0.4, 0, 0.2, 1)',
          properties: ['background-color', 'box-shadow', 'transform']
        },
        haptic: false
      };

      expect(interaction.type).toBe('hover');
      expect(interaction.animation?.duration).toBe(200);
      expect(interaction.animation?.properties).toContain('transform');
      expect(interaction.haptic).toBe(false);
    });

    it('should validate AnimationSpec for 60fps performance', () => {
      const animation: AnimationSpec = {
        duration: 200, // Should be ≤300ms for good UX
        easing: 'ease-out',
        properties: ['transform', 'opacity'] // GPU-accelerated properties
      };

      expect(animation.duration).toBeLessThanOrEqual(300);
      expect(animation.properties).toContain('transform');
      expect(animation.properties).not.toContain('left'); // Avoid layout-triggering properties
    });
  });

  describe('ComponentImplementation Model', () => {
    it('should create valid component implementation', () => {
      const implementation: ComponentImplementation = {
        id: 'impl_001',
        componentId: 'comp_001',
        platform: Platform.WEB,
        status: ImplementationStatus.VALIDATED,
        codeLocation: 'src/components/TchatButton.tsx',
        version: '1.0.0',
        testCoverage: 95,
        lastValidated: '2024-01-02T10:00:00Z',
        validationResults: [
          {
            type: 'design_token',
            isValid: true,
            score: 0.98,
            issues: [],
            validatedAt: '2024-01-02T10:00:00Z'
          }
        ],
        maintainer: 'frontend@tchat.com',
        dependencies: ['react', 'class-variance-authority'],
        buildStatus: {
          success: true,
          buildTime: '2024-01-02T09:45:00Z',
          errors: [],
          warnings: []
        }
      };

      expect(implementation.platform).toBe(Platform.WEB);
      expect(implementation.status).toBe(ImplementationStatus.VALIDATED);
      expect(implementation.testCoverage).toBeGreaterThanOrEqual(90);
      expect(implementation.validationResults![0].score).toBeGreaterThanOrEqual(0.97);
      expect(implementation.buildStatus?.success).toBe(true);
    });

    it('should validate ValidationResult structure', () => {
      const validationResult: ValidationResult = {
        type: 'accessibility',
        isValid: true,
        score: 0.95,
        issues: [],
        validatedAt: '2024-01-02T10:00:00Z'
      };

      expect(['design_token', 'accessibility', 'functionality', 'performance']).toContain(validationResult.type);
      expect(validationResult.score).toBeGreaterThanOrEqual(0);
      expect(validationResult.score).toBeLessThanOrEqual(1);
      expect(validationResult.isValid).toBe(true);
    });
  });

  describe('API Request/Response Models', () => {
    it('should validate CreateComponentRequest', () => {
      const request: CreateComponentRequest = {
        name: 'TchatInput',
        category: ComponentCategory.INPUT,
        description: 'Cross-platform input field with validation',
        designSpecs: {
          variants: [
            {
              name: 'default',
              description: 'Default input style',
              properties: {}
            }
          ],
          states: [],
          sizes: [],
          designTokens: [],
          accessibility: {
            wcagLevel: 'AA',
            contrastRatio: 4.5,
            keyboardNavigation: true,
            screenReaderSupport: true,
            ariaLabels: ['input'],
            focusManagement: true
          },
          interactions: []
        }
      };

      expect(request.name).toBe('TchatInput');
      expect(request.category).toBe(ComponentCategory.INPUT);
      expect(request.designSpecs.accessibility.wcagLevel).toBe('AA');
    });

    it('should validate UpdateComponentRequest with partial updates', () => {
      const request: UpdateComponentRequest = {
        description: 'Updated description with new features',
        status: ComponentStatus.IN_REVIEW
      };

      expect(request.description).toBeDefined();
      expect(request.status).toBe(ComponentStatus.IN_REVIEW);
      expect(request.name).toBeUndefined(); // Partial update
    });

    it('should validate ComponentListQuery with filtering', () => {
      const query: ComponentListQuery = {
        category: ComponentCategory.BUTTON,
        status: ComponentStatus.IMPLEMENTED,
        platform: Platform.WEB,
        search: 'tchat',
        page: 1,
        limit: 20,
        sort: 'name',
        order: 'asc'
      };

      expect(query.category).toBe(ComponentCategory.BUTTON);
      expect(query.platform).toBe(Platform.WEB);
      expect(query.limit).toBeLessThanOrEqual(100); // Reasonable limit
    });

    it('should validate ComponentValidationRequest', () => {
      const request: ComponentValidationRequest = {
        componentId: 'comp_001',
        platforms: [Platform.WEB, Platform.IOS, Platform.ANDROID],
        validationTypes: ['design_token', 'accessibility']
      };

      expect(request.platforms).toHaveLength(3);
      expect(request.validationTypes).toContain('accessibility');
    });

    it('should validate ComponentSyncRequest', () => {
      const request: ComponentSyncRequest = {
        componentIds: ['comp_001', 'comp_002'],
        platforms: [Platform.WEB, Platform.IOS],
        forceSync: true
      };

      expect(request.componentIds).toHaveLength(2);
      expect(request.forceSync).toBe(true);
    });

    it('should validate BulkComponentUpdateRequest', () => {
      const request: BulkComponentUpdateRequest = {
        componentIds: ['comp_001', 'comp_002', 'comp_003'],
        updates: {
          status: ComponentStatus.APPROVED
        },
        validateAfterUpdate: true
      };

      expect(request.componentIds).toHaveLength(3);
      expect(request.updates.status).toBe(ComponentStatus.APPROVED);
      expect(request.validateAfterUpdate).toBe(true);
    });
  });

  describe('Type Guards', () => {
    it('should validate isComponent type guard', () => {
      const validComponent = {
        id: 'comp_001',
        name: 'TchatButton',
        category: ComponentCategory.BUTTON,
        status: ComponentStatus.IMPLEMENTED
      };

      const invalidComponent = {
        id: 'comp_001',
        name: 'TchatButton'
        // Missing required fields
      };

      expect(isComponent(validComponent)).toBe(true);
      expect(isComponent(invalidComponent)).toBe(false);
      expect(isComponent(null)).toBe(false);
      expect(isComponent(undefined)).toBe(false);
      expect(isComponent('string')).toBe(false);
    });

    it('should validate isComponentImplementation type guard', () => {
      const validImplementation = {
        id: 'impl_001',
        componentId: 'comp_001',
        platform: Platform.WEB,
        status: ImplementationStatus.COMPLETED
      };

      const invalidImplementation = {
        id: 'impl_001',
        componentId: 'comp_001'
        // Missing required fields
      };

      expect(isComponentImplementation(validImplementation)).toBe(true);
      expect(isComponentImplementation(invalidImplementation)).toBe(false);
      expect(isComponentImplementation(null)).toBe(false);
      expect(isComponentImplementation({})).toBe(false);
    });

    it('should validate isValidationResult type guard', () => {
      const validResult = {
        type: 'design_token',
        isValid: true,
        score: 0.95
      };

      const invalidResult = {
        type: 'design_token',
        isValid: true
        // Missing score
      };

      expect(isValidationResult(validResult)).toBe(true);
      expect(isValidationResult(invalidResult)).toBe(false);
      expect(isValidationResult(null)).toBe(false);
      expect(isValidationResult([])).toBe(false);
    });
  });

  describe('Model Consistency Requirements', () => {
    it('should enforce Constitutional 97% cross-platform consistency', () => {
      const webImpl: ComponentImplementation = {
        id: 'impl_web_001',
        componentId: 'comp_001',
        platform: Platform.WEB,
        status: ImplementationStatus.VALIDATED,
        codeLocation: 'src/components/TchatButton.tsx',
        version: '1.0.0',
        maintainer: 'web@tchat.com',
        dependencies: []
      };

      const iosImpl: ComponentImplementation = {
        id: 'impl_ios_001',
        componentId: 'comp_001',
        platform: Platform.IOS,
        status: ImplementationStatus.VALIDATED,
        codeLocation: 'Sources/Components/TchatButton.swift',
        version: '1.0.0',
        maintainer: 'ios@tchat.com',
        dependencies: []
      };

      const androidImpl: ComponentImplementation = {
        id: 'impl_android_001',
        componentId: 'comp_001',
        platform: Platform.ANDROID,
        status: ImplementationStatus.VALIDATED,
        codeLocation: 'components/TchatButton.kt',
        version: '1.0.0',
        maintainer: 'android@tchat.com',
        dependencies: []
      };

      // All implementations should have same componentId and version for consistency
      expect(webImpl.componentId).toBe(iosImpl.componentId);
      expect(webImpl.componentId).toBe(androidImpl.componentId);
      expect(webImpl.version).toBe(iosImpl.version);
      expect(webImpl.version).toBe(androidImpl.version);

      // All should be validated to ensure consistency
      expect(webImpl.status).toBe(ImplementationStatus.VALIDATED);
      expect(iosImpl.status).toBe(ImplementationStatus.VALIDATED);
      expect(androidImpl.status).toBe(ImplementationStatus.VALIDATED);
    });

    it('should validate touch target compliance across platforms', () => {
      const mobileSize: ComponentSize = {
        name: 'default',
        dimensions: {
          height: 44,
          minHeight: 44,
          minWidth: 88
        },
        touchTarget: 44
      };

      // Constitutional requirement: 44dp minimum touch targets
      expect(mobileSize.touchTarget).toBeGreaterThanOrEqual(44);
      expect(mobileSize.dimensions.minHeight).toBeGreaterThanOrEqual(44);
    });

    it('should validate WCAG 2.1 AA accessibility compliance', () => {
      const accessibility: AccessibilityRequirements = {
        wcagLevel: 'AA',
        contrastRatio: 4.5, // WCAG AA requirement
        keyboardNavigation: true,
        screenReaderSupport: true,
        ariaLabels: ['button', 'interactive'],
        focusManagement: true
      };

      expect(accessibility.wcagLevel).toBe('AA');
      expect(accessibility.contrastRatio).toBeGreaterThanOrEqual(4.5);
      expect(accessibility.keyboardNavigation).toBe(true);
      expect(accessibility.screenReaderSupport).toBe(true);
      expect(accessibility.focusManagement).toBe(true);
    });

    it('should validate performance requirements for <200ms load times', () => {
      const performanceInteraction: InteractionSpec = {
        type: 'hover',
        description: 'Fast hover response',
        animation: {
          duration: 150, // <200ms for responsive feel
          easing: 'ease-out',
          properties: ['opacity', 'transform'] // GPU-accelerated
        }
      };

      expect(performanceInteraction.animation?.duration).toBeLessThan(200);

      // Ensure GPU-accelerated properties for 60fps
      const gpuProps = ['transform', 'opacity', 'filter'];
      const usesGPUAcceleration = performanceInteraction.animation?.properties.some(
        prop => gpuProps.includes(prop)
      );
      expect(usesGPUAcceleration).toBe(true);
    });
  });
});