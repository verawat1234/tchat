// T053 - MessageRegistry with lifecycle management
/**
 * MessageRegistry
 * Centralized registry for message component management and lifecycle
 * Provides registration, validation, and component metadata services
 */

import React from 'react';
import { MessageType, MessageData } from '../../types/MessageData';
import { BaseMessageProps } from './MessageComponentFactory';

// Registry entry interface
export interface MessageComponentEntry {
  component: React.ComponentType<any>;
  displayName: string;
  category: MessageComponentCategory;
  priority: number; // For preloading priority (1-10, higher = more important)
  supportLevel: MessageSupportLevel;
  metadata: ComponentMetadata;
  validation?: ComponentValidation;
}

// Component categories for organization
export enum MessageComponentCategory {
  COMMUNICATION = 'communication',
  INTERACTIVE = 'interactive',
  CONTENT = 'content',
  MEDIA = 'media',
  BUSINESS = 'business',
  UTILITY = 'utility'
}

// Support levels for components
export enum MessageSupportLevel {
  FULL = 'full',         // Complete implementation
  BETA = 'beta',         // Feature complete but testing
  PARTIAL = 'partial',   // Limited functionality
  DEPRECATED = 'deprecated' // Legacy support only
}

// Component metadata for registry management
export interface ComponentMetadata {
  version: string;
  author: string;
  description: string;
  dependencies: string[];
  features: string[];
  accessibility: AccessibilityInfo;
  performance: PerformanceMetrics;
  compatibility: BrowserCompatibility;
}

export interface AccessibilityInfo {
  wcagLevel: 'A' | 'AA' | 'AAA';
  screenReader: boolean;
  keyboardNavigation: boolean;
  colorContrast: boolean;
  focusManagement: boolean;
}

export interface PerformanceMetrics {
  bundleSize: number; // KB
  loadTime: number;   // ms
  renderTime: number; // ms
  memoryUsage: number; // MB
}

export interface BrowserCompatibility {
  chrome: string;
  firefox: string;
  safari: string;
  edge: string;
  mobile: boolean;
}

// Component validation interface
export interface ComponentValidation {
  validateProps: (props: any) => ValidationResult;
  validateContent: (content: any) => ValidationResult;
  sanitizeContent?: (content: any) => any;
}

export interface ValidationResult {
  isValid: boolean;
  errors: ValidationError[];
  warnings: ValidationWarning[];
}

export interface ValidationError {
  field: string;
  message: string;
  severity: 'critical' | 'error';
}

export interface ValidationWarning {
  field: string;
  message: string;
  suggestion?: string;
}

// Registry statistics and analytics
export interface RegistryStats {
  totalComponents: number;
  componentsByCategory: Record<MessageComponentCategory, number>;
  componentsBySupport: Record<MessageSupportLevel, number>;
  averageLoadTime: number;
  totalBundleSize: number;
  compatibilityScore: number;
}

// Component usage analytics
export interface ComponentUsageAnalytics {
  componentType: MessageType;
  usageCount: number;
  averageRenderTime: number;
  errorRate: number;
  userSatisfactionScore: number;
  lastUsed: Date;
}

// Main MessageRegistry class
export class MessageRegistry {
  private static instance: MessageRegistry;
  private components: Map<MessageType, MessageComponentEntry> = new Map();
  private analytics: Map<MessageType, ComponentUsageAnalytics> = new Map();
  private preloadQueue: Set<MessageType> = new Set();

  // Singleton pattern for global registry access
  public static getInstance(): MessageRegistry {
    if (!MessageRegistry.instance) {
      MessageRegistry.instance = new MessageRegistry();
    }
    return MessageRegistry.instance;
  }

  // Register a new message component
  public register(
    messageType: MessageType,
    entry: MessageComponentEntry
  ): void {
    if (this.components.has(messageType)) {
      console.warn(`Component ${messageType} is already registered. Overriding.`);
    }

    // Validate entry before registration
    this.validateEntry(entry);

    // Register component
    this.components.set(messageType, entry);

    // Initialize analytics
    this.analytics.set(messageType, {
      componentType: messageType,
      usageCount: 0,
      averageRenderTime: 0,
      errorRate: 0,
      userSatisfactionScore: 0,
      lastUsed: new Date()
    });

    // Add to preload queue based on priority
    if (entry.priority >= 7) {
      this.preloadQueue.add(messageType);
    }

    console.log(`Registered component: ${messageType} (${entry.displayName})`);
  }

  // Unregister a component
  public unregister(messageType: MessageType): boolean {
    const removed = this.components.delete(messageType);
    if (removed) {
      this.analytics.delete(messageType);
      this.preloadQueue.delete(messageType);
      console.log(`Unregistered component: ${messageType}`);
    }
    return removed;
  }

  // Get component entry
  public getComponent(messageType: MessageType): MessageComponentEntry | null {
    return this.components.get(messageType) || null;
  }

  // Check if component is registered
  public isRegistered(messageType: MessageType): boolean {
    return this.components.has(messageType);
  }

  // Get all registered components
  public getAllComponents(): Map<MessageType, MessageComponentEntry> {
    return new Map(this.components);
  }

  // Get components by category
  public getComponentsByCategory(
    category: MessageComponentCategory
  ): Map<MessageType, MessageComponentEntry> {
    const result = new Map<MessageType, MessageComponentEntry>();
    this.components.forEach((entry, type) => {
      if (entry.category === category) {
        result.set(type, entry);
      }
    });
    return result;
  }

  // Get components by support level
  public getComponentsBySupportLevel(
    supportLevel: MessageSupportLevel
  ): Map<MessageType, MessageComponentEntry> {
    const result = new Map<MessageType, MessageComponentEntry>();
    this.components.forEach((entry, type) => {
      if (entry.supportLevel === supportLevel) {
        result.set(type, entry);
      }
    });
    return result;
  }

  // Get high priority components for preloading
  public getPreloadComponents(): MessageType[] {
    return Array.from(this.preloadQueue);
  }

  // Validate message content against registered component
  public validateMessage(message: MessageData): ValidationResult {
    const entry = this.getComponent(message.type);

    if (!entry) {
      return {
        isValid: false,
        errors: [{
          field: 'type',
          message: `Unknown message type: ${message.type}`,
          severity: 'critical'
        }],
        warnings: []
      };
    }

    if (!entry.validation) {
      return {
        isValid: true,
        errors: [],
        warnings: [{
          field: 'validation',
          message: 'No validation rules defined for this component',
          suggestion: 'Consider adding validation for better error handling'
        }]
      };
    }

    // Validate content using component's validation rules
    const contentValidation = entry.validation.validateContent(message.content);

    // Additional registry-level validations
    const registryValidation = this.performRegistryValidation(message, entry);

    return {
      isValid: contentValidation.isValid && registryValidation.isValid,
      errors: [...contentValidation.errors, ...registryValidation.errors],
      warnings: [...contentValidation.warnings, ...registryValidation.warnings]
    };
  }

  // Update component analytics
  public recordUsage(
    messageType: MessageType,
    renderTime: number,
    hadError: boolean = false
  ): void {
    const analytics = this.analytics.get(messageType);
    if (!analytics) return;

    analytics.usageCount++;
    analytics.lastUsed = new Date();

    // Update average render time
    analytics.averageRenderTime =
      (analytics.averageRenderTime * (analytics.usageCount - 1) + renderTime) /
      analytics.usageCount;

    // Update error rate
    if (hadError) {
      analytics.errorRate =
        (analytics.errorRate * (analytics.usageCount - 1) + 100) /
        analytics.usageCount;
    } else {
      analytics.errorRate =
        (analytics.errorRate * (analytics.usageCount - 1)) /
        analytics.usageCount;
    }

    this.analytics.set(messageType, analytics);
  }

  // Get registry statistics
  public getStats(): RegistryStats {
    const componentsByCategory: Record<MessageComponentCategory, number> = {
      [MessageComponentCategory.COMMUNICATION]: 0,
      [MessageComponentCategory.INTERACTIVE]: 0,
      [MessageComponentCategory.CONTENT]: 0,
      [MessageComponentCategory.MEDIA]: 0,
      [MessageComponentCategory.BUSINESS]: 0,
      [MessageComponentCategory.UTILITY]: 0
    };

    const componentsBySupport: Record<MessageSupportLevel, number> = {
      [MessageSupportLevel.FULL]: 0,
      [MessageSupportLevel.BETA]: 0,
      [MessageSupportLevel.PARTIAL]: 0,
      [MessageSupportLevel.DEPRECATED]: 0
    };

    let totalLoadTime = 0;
    let totalBundleSize = 0;
    let compatibilityScoreSum = 0;

    this.components.forEach((entry) => {
      componentsByCategory[entry.category]++;
      componentsBySupport[entry.supportLevel]++;
      totalLoadTime += entry.metadata.performance.loadTime;
      totalBundleSize += entry.metadata.performance.bundleSize;

      // Calculate compatibility score (0-100)
      const compatibility = entry.metadata.compatibility;
      const mobileBonus = compatibility.mobile ? 20 : 0;
      compatibilityScoreSum += 80 + mobileBonus; // Base score for modern browsers
    });

    const totalComponents = this.components.size;

    return {
      totalComponents,
      componentsByCategory,
      componentsBySupport,
      averageLoadTime: totalComponents > 0 ? totalLoadTime / totalComponents : 0,
      totalBundleSize,
      compatibilityScore: totalComponents > 0 ? compatibilityScoreSum / totalComponents : 0
    };
  }

  // Get component analytics
  public getAnalytics(messageType?: MessageType): ComponentUsageAnalytics[] {
    if (messageType) {
      const analytics = this.analytics.get(messageType);
      return analytics ? [analytics] : [];
    }
    return Array.from(this.analytics.values());
  }

  // Export registry configuration for debugging
  public exportConfig(): any {
    const config: any = {
      timestamp: new Date().toISOString(),
      components: {},
      analytics: {},
      preloadQueue: Array.from(this.preloadQueue)
    };

    this.components.forEach((entry, type) => {
      config.components[type] = {
        displayName: entry.displayName,
        category: entry.category,
        priority: entry.priority,
        supportLevel: entry.supportLevel,
        metadata: entry.metadata
      };
    });

    this.analytics.forEach((analytics, type) => {
      config.analytics[type] = analytics;
    });

    return config;
  }

  // Import registry configuration
  public importConfig(config: any): void {
    console.warn('importConfig is for debugging only. Use register() for production.');

    if (config.preloadQueue) {
      this.preloadQueue = new Set(config.preloadQueue);
    }

    if (config.analytics) {
      Object.entries(config.analytics).forEach(([type, analytics]) => {
        this.analytics.set(type as MessageType, analytics as ComponentUsageAnalytics);
      });
    }
  }

  // Private validation methods
  private validateEntry(entry: MessageComponentEntry): void {
    if (!entry.component) {
      throw new Error('Component is required');
    }

    if (!entry.displayName) {
      throw new Error('Display name is required');
    }

    if (entry.priority < 1 || entry.priority > 10) {
      throw new Error('Priority must be between 1 and 10');
    }

    if (!entry.metadata) {
      throw new Error('Metadata is required');
    }

    // Validate metadata completeness
    const { metadata } = entry;
    if (!metadata.version || !metadata.author || !metadata.description) {
      throw new Error('Metadata must include version, author, and description');
    }

    if (!metadata.accessibility || !metadata.performance || !metadata.compatibility) {
      throw new Error('Metadata must include accessibility, performance, and compatibility info');
    }
  }

  private performRegistryValidation(
    message: MessageData,
    entry: MessageComponentEntry
  ): ValidationResult {
    const errors: ValidationError[] = [];
    const warnings: ValidationWarning[] = [];

    // Check support level warnings
    if (entry.supportLevel === MessageSupportLevel.DEPRECATED) {
      warnings.push({
        field: 'supportLevel',
        message: 'This message type is deprecated',
        suggestion: 'Consider migrating to a newer message type'
      });
    } else if (entry.supportLevel === MessageSupportLevel.BETA) {
      warnings.push({
        field: 'supportLevel',
        message: 'This message type is in beta',
        suggestion: 'Test thoroughly before using in production'
      });
    } else if (entry.supportLevel === MessageSupportLevel.PARTIAL) {
      warnings.push({
        field: 'supportLevel',
        message: 'This message type has limited functionality',
        suggestion: 'Some features may not be available'
      });
    }

    // Check accessibility compliance
    if (entry.metadata.accessibility.wcagLevel !== 'AA' &&
        entry.metadata.accessibility.wcagLevel !== 'AAA') {
      warnings.push({
        field: 'accessibility',
        message: 'Component does not meet WCAG AA standards',
        suggestion: 'Improve accessibility compliance for better user experience'
      });
    }

    return {
      isValid: errors.length === 0,
      errors,
      warnings
    };
  }
}

// Default registry instance
export const messageRegistry = MessageRegistry.getInstance();

// React hook for registry access
export const useMessageRegistry = () => {
  return React.useMemo(() => messageRegistry, []);
};

// Registry provider component for React context
const MessageRegistryContext = React.createContext<MessageRegistry>(messageRegistry);

export const MessageRegistryProvider: React.FC<{ children: React.ReactNode }> = ({
  children
}) => {
  return (
    <MessageRegistryContext.Provider value={messageRegistry}>
      {children}
    </MessageRegistryContext.Provider>
  );
};

export const useMessageRegistryContext = () => {
  const context = React.useContext(MessageRegistryContext);
  if (!context) {
    throw new Error('useMessageRegistryContext must be used within MessageRegistryProvider');
  }
  return context;
};

export default MessageRegistry;