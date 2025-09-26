/**
 * Components API Service (T037)
 *
 * RTK Query service for cross-platform component management.
 * Provides comprehensive CRUD operations, validation, and synchronization
 * for the design system component library across Web, iOS, and Android platforms.
 *
 * Features:
 * - Full CRUD operations for components and implementations
 * - Cross-platform validation and consistency checking
 * - Bulk operations for efficient management
 * - Real-time sync capabilities
 * - Component analytics and reporting
 * - Type-safe APIs with comprehensive error handling
 */

import { api } from './api';
import type {
  Component,
  ComponentImplementation,
  CreateComponentRequest,
  UpdateComponentRequest,
  ComponentListQuery,
  ComponentImplementationQuery,
  CreateComponentImplementationRequest,
  UpdateComponentImplementationRequest,
  ComponentValidationRequest,
  ComponentValidationResponse,
  ComponentSyncRequest,
  ComponentSyncResponse,
  BulkComponentUpdateRequest,
  BulkComponentUpdateResponse,
  ComponentAnalyticsQuery,
  ComponentAnalyticsResponse,
  Platform,
  ComponentStatus,
  ImplementationStatus,
} from '../types/component';
import { PaginatedResponse } from '../types/api';

/**
 * Components API Service
 *
 * Provides comprehensive component management functionality including:
 * - Component lifecycle management (create, read, update, delete)
 * - Cross-platform implementation tracking
 * - Design token integration and validation
 * - Accessibility compliance monitoring
 * - Real-time synchronization across platforms
 * - Analytics and reporting capabilities
 */
export const componentsApi = api.injectEndpoints({
  endpoints: (builder) => ({
    /**
     * Retrieves a paginated list of components with filtering and sorting
     */
    getComponents: builder.query<PaginatedResponse<Component>, ComponentListQuery | void>({
      query: (params = {}) => ({
        url: '/components',
        method: 'GET',
        params,
      }),
      providesTags: (result) =>
        result?.items
          ? [
              ...result.items.map(({ id }) => ({ type: 'Component' as const, id })),
              { type: 'Component', id: 'LIST' },
            ]
          : [{ type: 'Component', id: 'LIST' }],
      // Transform response to ensure type safety
      transformResponse: (response: any): PaginatedResponse<Component> => ({
        items: response.data?.items || [],
        pagination: response.data?.pagination || {
          hasMore: false,
          total: 0,
        },
      }),
    }),

    /**
     * Retrieves a single component by ID with full details
     */
    getComponent: builder.query<Component, string>({
      query: (componentId) => ({
        url: `/components/${componentId}`,
        method: 'GET',
      }),
      providesTags: (result, error, componentId) => [
        { type: 'Component', id: componentId },
      ],
      transformResponse: (response: any): Component => response.data,
    }),

    /**
     * Creates a new component with design specifications
     */
    createComponent: builder.mutation<Component, CreateComponentRequest>({
      query: (componentData) => ({
        url: '/components',
        method: 'POST',
        body: componentData,
      }),
      invalidatesTags: [{ type: 'Component', id: 'LIST' }],
      transformResponse: (response: any): Component => response.data,
    }),

    /**
     * Updates an existing component
     */
    updateComponent: builder.mutation<Component, { id: string; data: UpdateComponentRequest }>({
      query: ({ id, data }) => ({
        url: `/components/${id}`,
        method: 'PUT',
        body: data,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'Component', id },
        { type: 'Component', id: 'LIST' },
      ],
      transformResponse: (response: any): Component => response.data,
    }),

    /**
     * Deletes a component (soft delete - sets status to archived)
     */
    deleteComponent: builder.mutation<void, string>({
      query: (componentId) => ({
        url: `/components/${componentId}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, componentId) => [
        { type: 'Component', id: componentId },
        { type: 'Component', id: 'LIST' },
      ],
    }),

    /**
     * Bulk update multiple components
     */
    bulkUpdateComponents: builder.mutation<BulkComponentUpdateResponse, BulkComponentUpdateRequest>({
      query: (bulkData) => ({
        url: '/components/bulk',
        method: 'PUT',
        body: bulkData,
      }),
      invalidatesTags: [{ type: 'Component', id: 'LIST' }],
      transformResponse: (response: any): BulkComponentUpdateResponse => response.data,
    }),

    /**
     * Retrieves component implementations across platforms
     */
    getComponentImplementations: builder.query<
      PaginatedResponse<ComponentImplementation>,
      ComponentImplementationQuery | void
    >({
      query: (params = {}) => ({
        url: '/components/implementations',
        method: 'GET',
        params,
      }),
      providesTags: (result) =>
        result?.items
          ? [
              ...result.items.map(({ id }) => ({ type: 'ComponentImplementation' as const, id })),
              { type: 'ComponentImplementation', id: 'LIST' },
            ]
          : [{ type: 'ComponentImplementation', id: 'LIST' }],
      transformResponse: (response: any): PaginatedResponse<ComponentImplementation> => ({
        items: response.data?.items || [],
        pagination: response.data?.pagination || {
          hasMore: false,
          total: 0,
        },
      }),
    }),

    /**
     * Retrieves implementations for a specific component
     */
    getComponentImplementationsByComponent: builder.query<
      ComponentImplementation[],
      { componentId: string; platform?: Platform }
    >({
      query: ({ componentId, platform }) => ({
        url: `/components/${componentId}/implementations`,
        method: 'GET',
        params: platform ? { platform } : undefined,
      }),
      providesTags: (result, error, { componentId }) => [
        { type: 'ComponentImplementation', id: `COMPONENT_${componentId}` },
      ],
      transformResponse: (response: any): ComponentImplementation[] => response.data,
    }),

    /**
     * Creates a new component implementation
     */
    createComponentImplementation: builder.mutation<
      ComponentImplementation,
      CreateComponentImplementationRequest
    >({
      query: (implementationData) => ({
        url: '/components/implementations',
        method: 'POST',
        body: implementationData,
      }),
      invalidatesTags: (result, error, implementationData) => [
        { type: 'ComponentImplementation', id: 'LIST' },
        { type: 'ComponentImplementation', id: `COMPONENT_${implementationData.componentId}` },
      ],
      transformResponse: (response: any): ComponentImplementation => response.data,
    }),

    /**
     * Updates a component implementation
     */
    updateComponentImplementation: builder.mutation<
      ComponentImplementation,
      { id: string; data: UpdateComponentImplementationRequest }
    >({
      query: ({ id, data }) => ({
        url: `/components/implementations/${id}`,
        method: 'PUT',
        body: data,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'ComponentImplementation', id },
        { type: 'ComponentImplementation', id: 'LIST' },
      ],
      transformResponse: (response: any): ComponentImplementation => response.data,
    }),

    /**
     * Updates implementation status (for build pipelines)
     */
    updateImplementationStatus: builder.mutation<
      ComponentImplementation,
      { id: string; status: ImplementationStatus; buildInfo?: any }
    >({
      query: ({ id, status, buildInfo }) => ({
        url: `/components/implementations/${id}/status`,
        method: 'PATCH',
        body: { status, buildInfo },
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'ComponentImplementation', id },
      ],
      transformResponse: (response: any): ComponentImplementation => response.data,
    }),

    /**
     * Validates component consistency across platforms
     */
    validateComponent: builder.mutation<ComponentValidationResponse, ComponentValidationRequest>({
      query: (validationRequest) => ({
        url: '/components/validate',
        method: 'POST',
        body: validationRequest,
      }),
      transformResponse: (response: any): ComponentValidationResponse => response.data,
    }),

    /**
     * Validates all components in bulk
     */
    validateAllComponents: builder.mutation<ComponentValidationResponse[], void>({
      query: () => ({
        url: '/components/validate/all',
        method: 'POST',
      }),
      transformResponse: (response: any): ComponentValidationResponse[] => response.data,
    }),

    /**
     * Synchronizes components across platforms
     */
    syncComponents: builder.mutation<ComponentSyncResponse, ComponentSyncRequest>({
      query: (syncRequest) => ({
        url: '/components/sync',
        method: 'POST',
        body: syncRequest,
      }),
      invalidatesTags: [
        { type: 'Component', id: 'LIST' },
        { type: 'ComponentImplementation', id: 'LIST' },
      ],
      transformResponse: (response: any): ComponentSyncResponse => response.data,
    }),

    /**
     * Updates component status (draft, approved, deprecated, etc.)
     */
    updateComponentStatus: builder.mutation<
      Component,
      { id: string; status: ComponentStatus; reason?: string }
    >({
      query: ({ id, status, reason }) => ({
        url: `/components/${id}/status`,
        method: 'PATCH',
        body: { status, reason },
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'Component', id },
        { type: 'Component', id: 'LIST' },
      ],
      transformResponse: (response: any): Component => response.data,
    }),

    /**
     * Retrieves component analytics and metrics
     */
    getComponentAnalytics: builder.query<ComponentAnalyticsResponse, ComponentAnalyticsQuery | void>({
      query: (params = {}) => ({
        url: '/components/analytics',
        method: 'GET',
        params,
      }),
      // Cache analytics for 5 minutes
      keepUnusedDataFor: 300,
      transformResponse: (response: any): ComponentAnalyticsResponse => response.data,
    }),

    /**
     * Retrieves component usage statistics
     */
    getComponentUsage: builder.query<
      { componentId: string; usageCount: number; platforms: Platform[] }[],
      { componentIds?: string[]; platform?: Platform }
    >({
      query: (params) => ({
        url: '/components/usage',
        method: 'GET',
        params,
      }),
      transformResponse: (response: any) => response.data,
    }),

    /**
     * Searches components with advanced filtering
     */
    searchComponents: builder.query<
      PaginatedResponse<Component>,
      { query: string; filters?: ComponentListQuery }
    >({
      query: ({ query, filters }) => ({
        url: '/components/search',
        method: 'GET',
        params: { q: query, ...filters },
      }),
      providesTags: [{ type: 'Component', id: 'SEARCH' }],
      transformResponse: (response: any): PaginatedResponse<Component> => ({
        items: response.data?.items || [],
        pagination: response.data?.pagination || {
          hasMore: false,
          total: 0,
        },
      }),
    }),

    /**
     * Gets component dependencies and relationships
     */
    getComponentDependencies: builder.query<
      {
        dependencies: { id: string; name: string; version: string }[];
        dependents: { id: string; name: string; version: string }[];
        designTokens: { id: string; name: string; category: string }[];
      },
      string
    >({
      query: (componentId) => ({
        url: `/components/${componentId}/dependencies`,
        method: 'GET',
      }),
      providesTags: (result, error, componentId) => [
        { type: 'Component', id: `DEPS_${componentId}` },
      ],
      transformResponse: (response: any) => response.data,
    }),

    /**
     * Exports component specifications for external tools
     */
    exportComponent: builder.mutation<
      { url: string; format: string },
      { componentId: string; format: 'figma' | 'storybook' | 'json' | 'documentation' }
    >({
      query: ({ componentId, format }) => ({
        url: `/components/${componentId}/export`,
        method: 'POST',
        body: { format },
      }),
      transformResponse: (response: any) => response.data,
    }),

    /**
     * Imports component specifications from external tools
     */
    importComponent: builder.mutation<
      Component,
      { source: 'figma' | 'storybook' | 'json'; data: any; options?: any }
    >({
      query: ({ source, data, options }) => ({
        url: '/components/import',
        method: 'POST',
        body: { source, data, options },
      }),
      invalidatesTags: [{ type: 'Component', id: 'LIST' }],
      transformResponse: (response: any): Component => response.data,
    }),
  }),
  overrideExisting: false,
});

// Export hooks for use in React components
export const {
  // Component queries
  useGetComponentsQuery,
  useGetComponentQuery,
  useLazyGetComponentQuery,
  useSearchComponentsQuery,
  useLazySearchComponentsQuery,
  useGetComponentDependenciesQuery,
  useGetComponentUsageQuery,
  useGetComponentAnalyticsQuery,

  // Implementation queries
  useGetComponentImplementationsQuery,
  useGetComponentImplementationsByComponentQuery,
  useLazyGetComponentImplementationsByComponentQuery,

  // Component mutations
  useCreateComponentMutation,
  useUpdateComponentMutation,
  useDeleteComponentMutation,
  useBulkUpdateComponentsMutation,
  useUpdateComponentStatusMutation,

  // Implementation mutations
  useCreateComponentImplementationMutation,
  useUpdateComponentImplementationMutation,
  useUpdateImplementationStatusMutation,

  // Validation and sync mutations
  useValidateComponentMutation,
  useValidateAllComponentsMutation,
  useSyncComponentsMutation,

  // Import/export mutations
  useExportComponentMutation,
  useImportComponentMutation,
} = componentsApi;

/**
 * Component management utilities
 */
export const componentUtils = {
  /**
   * Calculates implementation coverage across platforms
   */
  calculateImplementationCoverage: (implementations: ComponentImplementation[]): Record<Platform, number> => {
    const total = Object.values(Platform).length;
    const coverage: Record<Platform, number> = {} as Record<Platform, number>;

    Object.values(Platform).forEach(platform => {
      const implementationExists = implementations.some(
        impl => impl.platform === platform && impl.status !== 'not_started'
      );
      coverage[platform] = implementationExists ? 100 : 0;
    });

    return coverage;
  },

  /**
   * Determines if component is ready for release
   */
  isReadyForRelease: (component: Component, implementations: ComponentImplementation[]): boolean => {
    // Component must be approved
    if (component.status !== 'approved' && component.status !== 'implemented') {
      return false;
    }

    // All targeted platforms must have validated implementations
    const requiredPlatforms = Object.values(Platform);
    return requiredPlatforms.every(platform => {
      const impl = implementations.find(impl => impl.platform === platform);
      return impl && impl.status === 'validated';
    });
  },

  /**
   * Gets validation issues count by severity
   */
  getValidationIssuesSummary: (validationResponse: ComponentValidationResponse) => {
    let critical = 0, major = 0, minor = 0;

    Object.values(validationResponse.platformResults).forEach(platformResult => {
      platformResult.issues.forEach(issue => {
        switch (issue.severity) {
          case 'error':
            critical++;
            break;
          case 'warning':
            major++;
            break;
          case 'info':
            minor++;
            break;
        }
      });
    });

    return { critical, major, minor };
  },

  /**
   * Formats component changelog for display
   */
  formatChangelog: (component: Component): string[] => {
    if (!component.metadata?.changelog) return [];

    return component.metadata.changelog
      .sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
      .map(entry =>
        `v${entry.version} (${entry.date}): ${entry.changes.join(', ')}${entry.breaking ? ' [BREAKING]' : ''}`
      );
  },
};

// Export component status helpers
export const componentStatusHelpers = {
  canEdit: (status: ComponentStatus): boolean => {
    return status === 'draft' || status === 'in_review';
  },

  canApprove: (status: ComponentStatus): boolean => {
    return status === 'in_review';
  },

  canImplement: (status: ComponentStatus): boolean => {
    return status === 'approved';
  },

  canDeprecate: (status: ComponentStatus): boolean => {
    return status === 'implemented' || status === 'approved';
  },

  getStatusColor: (status: ComponentStatus): string => {
    switch (status) {
      case 'draft': return 'gray';
      case 'in_review': return 'yellow';
      case 'approved': return 'blue';
      case 'implemented': return 'green';
      case 'deprecated': return 'red';
      default: return 'gray';
    }
  },
};

// Export implementation status helpers
export const implementationStatusHelpers = {
  getStatusColor: (status: ImplementationStatus): string => {
    switch (status) {
      case 'not_started': return 'gray';
      case 'in_progress': return 'blue';
      case 'completed': return 'green';
      case 'testing': return 'yellow';
      case 'validated': return 'emerald';
      case 'needs_update': return 'orange';
      default: return 'gray';
    }
  },

  isBlockingRelease: (status: ImplementationStatus): boolean => {
    return status === 'not_started' || status === 'needs_update';
  },

  getNextStatus: (currentStatus: ImplementationStatus): ImplementationStatus | null => {
    switch (currentStatus) {
      case 'not_started': return 'in_progress';
      case 'in_progress': return 'completed';
      case 'completed': return 'testing';
      case 'testing': return 'validated';
      case 'needs_update': return 'in_progress';
      default: return null;
    }
  },
};