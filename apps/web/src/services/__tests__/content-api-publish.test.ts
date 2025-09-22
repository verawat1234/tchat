import { describe, it, expect, beforeEach, vi } from 'vitest';
import { configureStore } from '@reduxjs/toolkit';
import { setupListeners } from '@reduxjs/toolkit/query';
import { api } from '../api';
import { contentApi } from '../content';
import {
  ContentItem,
  ContentType,
  ContentStatus,
  ContentCategory,
  ContentValue,
  ContentMetadata
} from '../../types/content';
import type { ApiError } from '../../types/api';

/**
 * T012: Contract Test for publishContent API
 *
 * CRITICAL TDD REQUIREMENT: This test MUST FAIL because the implementation doesn't exist yet.
 * This is intentional and required for proper test-driven development.
 *
 * This test validates:
 * 1. Content publication workflow and state transitions
 * 2. Status change from draft/archived to published
 * 3. Publication metadata (publishedAt, publishedBy)
 * 4. Change log requirements for publication
 * 5. Validation rules for publishable content
 * 6. Permission-based publication control
 * 7. Error handling for invalid state transitions
 */

// PublishContentRequest interface definition (expected API contract)
interface PublishContentRequest {
  /** Content item ID to publish */
  contentId: string;
  /** Required change log describing what is being published */
  changeLog: string;
  /** Optional publication notes */
  notes?: string;
  /** Force publication even if warnings exist */
  force?: boolean;
}

interface PublishContentResponse {
  /** Updated content item with published status */
  contentItem: ContentItem;
  /** Publication metadata */
  publication: {
    /** ISO timestamp when content was published */
    publishedAt: string;
    /** User ID who published the content */
    publishedBy: string;
    /** Version number that was published */
    version: number;
    /** Publication change log */
    changeLog: string;
  };
  /** Warnings encountered during publication (if any) */
  warnings?: string[];
}

// Mock data for testing
const mockContentCategory: ContentCategory = {
  id: 'marketing',
  name: 'Marketing Content',
  description: 'Marketing materials and promotional content',
  permissions: {
    read: ['user', 'editor', 'admin'],
    write: ['editor', 'admin'],
    publish: ['admin']
  }
};

const mockTextContent: ContentValue = {
  type: 'text',
  value: 'New product launch announcement',
  maxLength: 500
};

const mockMetadata: ContentMetadata = {
  createdAt: '2024-01-15T10:00:00.000Z',
  createdBy: 'user-editor-123',
  updatedAt: '2024-01-15T14:30:00.000Z',
  updatedBy: 'user-editor-123',
  version: 3,
  tags: ['announcement', 'product', 'launch'],
  notes: 'Content ready for review and publication'
};

const mockDraftContentItem: ContentItem = {
  id: 'marketing.announcement.product-launch',
  category: mockContentCategory,
  type: ContentType.TEXT,
  value: mockTextContent,
  metadata: mockMetadata,
  status: ContentStatus.DRAFT
};

const mockArchivedContentItem: ContentItem = {
  ...mockDraftContentItem,
  id: 'marketing.announcement.old-product',
  status: ContentStatus.ARCHIVED,
  metadata: {
    ...mockMetadata,
    version: 5
  }
};

const mockPublishedContentItem: ContentItem = {
  ...mockDraftContentItem,
  status: ContentStatus.PUBLISHED,
  metadata: {
    ...mockMetadata,
    version: 4,
    updatedAt: '2024-01-15T15:00:00.000Z',
    updatedBy: 'user-admin-456'
  }
};

const mockPublishRequest: PublishContentRequest = {
  contentId: 'marketing.announcement.product-launch',
  changeLog: 'Publishing product launch announcement for Q1 2024',
  notes: 'Approved by marketing team and legal review completed'
};

const mockPublishResponse: PublishContentResponse = {
  contentItem: mockPublishedContentItem,
  publication: {
    publishedAt: '2024-01-15T15:00:00.000Z',
    publishedBy: 'user-admin-456',
    version: 4,
    changeLog: 'Publishing product launch announcement for Q1 2024'
  }
};

// Error responses for testing
const mockValidationError: ApiError = {
  success: false,
  error: {
    code: 'VALIDATION_ERROR',
    message: 'Content cannot be published: missing required change log',
    details: {
      contentId: 'marketing.announcement.product-launch',
      missingFields: ['changeLog'],
      validationRules: {
        changeLog: 'Required for publication tracking',
        minLength: 10,
        maxLength: 1000
      }
    },
    timestamp: '2024-01-15T15:00:00.000Z'
  }
};

const mockPermissionError: ApiError = {
  success: false,
  error: {
    code: 'FORBIDDEN',
    message: 'Insufficient permissions to publish content in this category',
    details: {
      contentId: 'marketing.announcement.product-launch',
      category: 'marketing',
      requiredRole: 'admin',
      userRole: 'editor',
      action: 'publish'
    },
    timestamp: '2024-01-15T15:00:00.000Z'
  }
};

const mockStateTransitionError: ApiError = {
  success: false,
  error: {
    code: 'INVALID_STATE_TRANSITION',
    message: 'Content is already published',
    details: {
      contentId: 'marketing.announcement.already-published',
      currentStatus: 'published',
      requestedStatus: 'published',
      validTransitions: ['archived'],
      publishedAt: '2024-01-10T10:00:00.000Z',
      publishedBy: 'user-admin-123'
    },
    timestamp: '2024-01-15T15:00:00.000Z'
  }
};

const mockContentNotFoundError: ApiError = {
  success: false,
  error: {
    code: 'NOT_FOUND',
    message: 'Content item not found',
    details: {
      contentId: 'marketing.announcement.non-existent'
    },
    timestamp: '2024-01-15T15:00:00.000Z'
  }
};

// Helper to create test store with API
const createTestStore = () => {
  const store = configureStore({
    reducer: {
      api: api.reducer,
    },
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware().concat(api.middleware),
  });

  setupListeners(store.dispatch);
  return store;
};

// Mock fetch for testing
const mockFetch = vi.fn();
global.fetch = mockFetch;

describe('Content API - publishContent', () => {
  let store: ReturnType<typeof createTestStore>;

  beforeEach(() => {
    store = createTestStore();
    mockFetch.mockClear();
    vi.clearAllMocks();
  });

  describe('Successful content publication', () => {
    it('should publish draft content with complete workflow validation', async () => {
      // Arrange: Mock successful publication response
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: mockPublishResponse,
          meta: {
            timestamp: '2024-01-15T15:00:00.000Z',
            version: '1.0.0'
          }
        }),
      });

      // Act: Trigger the publication mutation that WILL FAIL (TDD requirement)
      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(mockPublishRequest)
      );

      const response = await result;

      // Assert: Validate successful publication response
      expect(response.data).toBeDefined();
      expect(response.error).toBeUndefined();

      const publishData = response.data as PublishContentResponse;

      // Validate content item state change
      expect(publishData.contentItem.id).toBe('marketing.announcement.product-launch');
      expect(publishData.contentItem.status).toBe(ContentStatus.PUBLISHED);
      expect(publishData.contentItem.metadata.version).toBe(4);

      // Validate publication metadata
      expect(publishData.publication.publishedAt).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/);
      expect(publishData.publication.publishedBy).toBe('user-admin-456');
      expect(publishData.publication.version).toBe(4);
      expect(publishData.publication.changeLog).toBe('Publishing product launch announcement for Q1 2024');

      // Validate API call was made correctly
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/content/marketing.announcement.product-launch/publish'),
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Content-Type': 'application/json'
          }),
          body: JSON.stringify(mockPublishRequest)
        })
      );
    });

    it('should publish archived content back to published status', async () => {
      const republishRequest: PublishContentRequest = {
        contentId: 'marketing.announcement.old-product',
        changeLog: 'Re-publishing archived content due to renewed relevance',
        notes: 'Marketing team requested to bring this content back',
        force: false
      };

      const republishResponse: PublishContentResponse = {
        contentItem: {
          ...mockArchivedContentItem,
          status: ContentStatus.PUBLISHED,
          metadata: {
            ...mockArchivedContentItem.metadata,
            version: 6,
            updatedAt: '2024-01-15T16:00:00.000Z',
            updatedBy: 'user-admin-456'
          }
        },
        publication: {
          publishedAt: '2024-01-15T16:00:00.000Z',
          publishedBy: 'user-admin-456',
          version: 6,
          changeLog: 'Re-publishing archived content due to renewed relevance'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: republishResponse
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(republishRequest)
      );

      const response = await result;

      expect(response.data).toBeDefined();
      const publishData = response.data as PublishContentResponse;

      // Validate archived → published transition
      expect(publishData.contentItem.status).toBe(ContentStatus.PUBLISHED);
      expect(publishData.contentItem.metadata.version).toBe(6);
      expect(publishData.publication.changeLog).toBe('Re-publishing archived content due to renewed relevance');
    });

    it('should handle publication with warnings', async () => {
      const publishResponseWithWarnings: PublishContentResponse = {
        ...mockPublishResponse,
        warnings: [
          'Content contains links that may become outdated',
          'Image optimization recommended for better performance',
          'Consider adding alternative text for accessibility'
        ]
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: publishResponseWithWarnings
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(mockPublishRequest)
      );

      const response = await result;

      expect(response.data).toBeDefined();
      const publishData = response.data as PublishContentResponse;

      // Validate warnings are included
      expect(publishData.warnings).toBeDefined();
      expect(publishData.warnings).toHaveLength(3);
      expect(publishData.warnings).toContain('Content contains links that may become outdated');
      expect(publishData.warnings).toContain('Image optimization recommended for better performance');
    });

    it('should handle different content types during publication', async () => {
      const richTextContent: ContentValue = {
        type: 'rich_text',
        value: '<h1>Product Launch</h1><p>Exciting new features coming soon!</p>',
        format: 'html',
        allowedTags: ['h1', 'h2', 'p', 'strong', 'em', 'ul', 'li']
      };

      const richTextItem: ContentItem = {
        ...mockDraftContentItem,
        id: 'marketing.campaign.rich-content',
        type: ContentType.RICH_TEXT,
        value: richTextContent
      };

      const richTextPublishResponse: PublishContentResponse = {
        contentItem: {
          ...richTextItem,
          status: ContentStatus.PUBLISHED,
          metadata: {
            ...richTextItem.metadata,
            version: richTextItem.metadata.version + 1
          }
        },
        publication: {
          publishedAt: '2024-01-15T15:30:00.000Z',
          publishedBy: 'user-admin-456',
          version: richTextItem.metadata.version + 1,
          changeLog: 'Publishing rich text marketing campaign content'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: richTextPublishResponse
        }),
      });

      const richTextPublishRequest: PublishContentRequest = {
        contentId: 'marketing.campaign.rich-content',
        changeLog: 'Publishing rich text marketing campaign content'
      };

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(richTextPublishRequest)
      );

      const response = await result;

      expect(response.data).toBeDefined();
      const publishData = response.data as PublishContentResponse;
      expect(publishData.contentItem.type).toBe(ContentType.RICH_TEXT);
      expect(publishData.contentItem.status).toBe(ContentStatus.PUBLISHED);
    });
  });

  describe('Validation and error handling', () => {
    it('should reject publication without required change log', async () => {
      const invalidRequest = {
        contentId: 'marketing.announcement.product-launch',
        // Missing required changeLog field
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => mockValidationError,
      });

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(invalidRequest as PublishContentRequest)
      );

      const response = await result;

      expect(response.data).toBeUndefined();
      expect(response.error).toBeDefined();

      const error = response.error as any;
      expect(error.status).toBe(400);
      expect(error.data.error.code).toBe('VALIDATION_ERROR');
      expect(error.data.error.details.missingFields).toContain('changeLog');
      expect(error.data.error.details.validationRules.changeLog).toBe('Required for publication tracking');
    });

    it('should reject publication with insufficient permissions', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 403,
        json: async () => mockPermissionError,
      });

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(mockPublishRequest)
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(403);
      expect(error.data.error.code).toBe('FORBIDDEN');
      expect(error.data.error.details.requiredRole).toBe('admin');
      expect(error.data.error.details.userRole).toBe('editor');
      expect(error.data.error.details.action).toBe('publish');
    });

    it('should reject invalid state transitions', async () => {
      const alreadyPublishedRequest: PublishContentRequest = {
        contentId: 'marketing.announcement.already-published',
        changeLog: 'Attempting to republish already published content'
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 409,
        json: async () => mockStateTransitionError,
      });

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(alreadyPublishedRequest)
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(409);
      expect(error.data.error.code).toBe('INVALID_STATE_TRANSITION');
      expect(error.data.error.details.currentStatus).toBe('published');
      expect(error.data.error.details.requestedStatus).toBe('published');
      expect(error.data.error.details.validTransitions).toContain('archived');
    });

    it('should handle non-existent content gracefully', async () => {
      const nonExistentRequest: PublishContentRequest = {
        contentId: 'marketing.announcement.non-existent',
        changeLog: 'Trying to publish non-existent content'
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: async () => mockContentNotFoundError,
      });

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(nonExistentRequest)
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(404);
      expect(error.data.error.code).toBe('NOT_FOUND');
      expect(error.data.error.details.contentId).toBe('marketing.announcement.non-existent');
    });

    it('should validate change log length requirements', async () => {
      const shortChangeLogRequest: PublishContentRequest = {
        contentId: 'marketing.announcement.product-launch',
        changeLog: 'Short' // Too short
      };

      const changeLogValidationError: ApiError = {
        success: false,
        error: {
          code: 'VALIDATION_ERROR',
          message: 'Change log must be between 10 and 1000 characters',
          details: {
            contentId: 'marketing.announcement.product-launch',
            changeLog: {
              provided: 'Short',
              length: 5,
              minLength: 10,
              maxLength: 1000
            }
          },
          timestamp: '2024-01-15T15:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => changeLogValidationError,
      });

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(shortChangeLogRequest)
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(400);
      expect(error.data.error.details.changeLog.length).toBe(5);
      expect(error.data.error.details.changeLog.minLength).toBe(10);
    });
  });

  describe('Permission-based publication control', () => {
    it('should validate category-specific publication permissions', async () => {
      const restrictedCategoryRequest: PublishContentRequest = {
        contentId: 'admin.system.critical-config',
        changeLog: 'Publishing system critical configuration'
      };

      const restrictedPermissionError: ApiError = {
        success: false,
        error: {
          code: 'FORBIDDEN',
          message: 'Publishing system critical content requires super admin permissions',
          details: {
            contentId: 'admin.system.critical-config',
            category: 'admin',
            subcategory: 'system',
            requiredRole: 'super_admin',
            userRole: 'admin',
            action: 'publish',
            securityLevel: 'critical'
          },
          timestamp: '2024-01-15T15:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 403,
        json: async () => restrictedPermissionError,
      });

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(restrictedCategoryRequest)
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.data.error.details.securityLevel).toBe('critical');
      expect(error.data.error.details.requiredRole).toBe('super_admin');
    });

    it('should handle bulk publication with mixed permissions', async () => {
      const bulkPublishRequest = {
        contentIds: [
          'marketing.announcement.product-launch',
          'admin.system.critical-config', // Should fail due to permissions
          'content.landing.hero'
        ],
        changeLog: 'Bulk publication of Q1 content updates'
      };

      const bulkPermissionError: ApiError = {
        success: false,
        error: {
          code: 'PARTIAL_PERMISSION_DENIED',
          message: 'Some content items cannot be published due to insufficient permissions',
          details: {
            allowedItems: ['marketing.announcement.product-launch', 'content.landing.hero'],
            deniedItems: [
              {
                contentId: 'admin.system.critical-config',
                reason: 'Requires super_admin role',
                requiredRole: 'super_admin',
                userRole: 'admin'
              }
            ],
            totalRequested: 3,
            totalAllowed: 2,
            totalDenied: 1
          },
          timestamp: '2024-01-15T15:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 207, // Multi-status
        json: async () => bulkPermissionError,
      });

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(bulkPublishRequest as any)
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(207);
      expect(error.data.error.code).toBe('PARTIAL_PERMISSION_DENIED');
      expect(error.data.error.details.allowedItems).toHaveLength(2);
      expect(error.data.error.details.deniedItems).toHaveLength(1);
    });
  });

  describe('RTK Query integration and caching', () => {
    it('should provide correct cache tags for published content', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: mockPublishResponse
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(mockPublishRequest)
      );

      await result;

      // Publication should invalidate related cache tags:
      // - Content item cache
      // - Content list caches
      // - Category-specific caches
      // When implemented, this should verify cache invalidation
      expect(mockFetch).toHaveBeenCalledTimes(1);
    });

    it('should handle optimistic updates for publication', async () => {
      // Test would validate optimistic UI updates during publication
      // This ensures users see immediate feedback while API call is in progress

      mockFetch.mockImplementationOnce(() =>
        new Promise(resolve => {
          setTimeout(() => {
            resolve({
              ok: true,
              status: 200,
              json: async () => ({
                success: true,
                data: mockPublishResponse
              }),
            });
          }, 100); // Simulate network delay
        })
      );

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(mockPublishRequest)
      );

      // During this time, optimistic updates should show publishing state
      expect(result).toBeDefined();

      await result;

      // After completion, actual published state should be reflected
      expect(mockFetch).toHaveBeenCalledTimes(1);
    });

    it('should support publication rollback on failure', async () => {
      const rollbackError: ApiError = {
        success: false,
        error: {
          code: 'PUBLICATION_FAILED',
          message: 'Publication failed during final validation',
          details: {
            contentId: 'marketing.announcement.product-launch',
            stage: 'final_validation',
            rollbackRequired: true,
            originalStatus: 'draft',
            failureReason: 'Content validation service unavailable'
          },
          timestamp: '2024-01-15T15:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: async () => rollbackError,
      });

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(mockPublishRequest)
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.data.error.code).toBe('PUBLICATION_FAILED');
      expect(error.data.error.details.rollbackRequired).toBe(true);
      expect(error.data.error.details.originalStatus).toBe('draft');
    });
  });

  describe('Type safety and TypeScript integration', () => {
    it('should provide proper TypeScript typing for publication requests', () => {
      // Compile-time type checking
      const validRequest: PublishContentRequest = {
        contentId: 'marketing.announcement.product-launch',
        changeLog: 'Publishing Q1 content updates'
      };

      const requestWithOptionals: PublishContentRequest = {
        contentId: 'marketing.announcement.product-launch',
        changeLog: 'Publishing Q1 content updates',
        notes: 'Approved by marketing team',
        force: false
      };

      expect(validRequest.contentId).toBe('marketing.announcement.product-launch');
      expect(requestWithOptionals.notes).toBe('Approved by marketing team');
      expect(requestWithOptionals.force).toBe(false);
    });

    it('should provide proper TypeScript typing for publication responses', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: mockPublishResponse
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(mockPublishRequest)
      );

      const response = await result;

      if (response.data) {
        // TypeScript should properly infer these types
        const contentItem: ContentItem = response.data.contentItem;
        const publishedAt: string = response.data.publication.publishedAt;
        const publishedBy: string = response.data.publication.publishedBy;
        const version: number = response.data.publication.version;
        const changeLog: string = response.data.publication.changeLog;

        expect(contentItem.status).toBe(ContentStatus.PUBLISHED);
        expect(typeof publishedAt).toBe('string');
        expect(typeof publishedBy).toBe('string');
        expect(typeof version).toBe('number');
        expect(typeof changeLog).toBe('string');
      }
    });
  });

  describe('Edge cases and boundary conditions', () => {
    it('should handle concurrent publication attempts', async () => {
      const concurrencyError: ApiError = {
        success: false,
        error: {
          code: 'CONCURRENT_MODIFICATION',
          message: 'Content was modified by another user during publication',
          details: {
            contentId: 'marketing.announcement.product-launch',
            expectedVersion: 3,
            actualVersion: 4,
            lastModifiedBy: 'user-editor-789',
            lastModifiedAt: '2024-01-15T14:58:00.000Z'
          },
          timestamp: '2024-01-15T15:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 409,
        json: async () => concurrencyError,
      });

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(mockPublishRequest)
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.data.error.code).toBe('CONCURRENT_MODIFICATION');
      expect(error.data.error.details.expectedVersion).toBe(3);
      expect(error.data.error.details.actualVersion).toBe(4);
    });

    it('should handle publication of content with dependencies', async () => {
      const dependencyError: ApiError = {
        success: false,
        error: {
          code: 'DEPENDENCY_CHECK_FAILED',
          message: 'Content has unpublished dependencies that must be published first',
          details: {
            contentId: 'marketing.campaign.main-landing',
            dependencies: [
              {
                contentId: 'marketing.assets.hero-image',
                status: 'draft',
                required: true
              },
              {
                contentId: 'marketing.copy.tagline',
                status: 'draft',
                required: true
              }
            ],
            suggestedAction: 'Publish dependencies first or use force flag'
          },
          timestamp: '2024-01-15T15:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 412,
        json: async () => dependencyError,
      });

      const dependentContentRequest: PublishContentRequest = {
        contentId: 'marketing.campaign.main-landing',
        changeLog: 'Publishing main landing page campaign'
      };

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(dependentContentRequest)
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(412);
      expect(error.data.error.code).toBe('DEPENDENCY_CHECK_FAILED');
      expect(error.data.error.details.dependencies).toHaveLength(2);
      expect(error.data.error.details.suggestedAction).toContain('force flag');
    });

    it('should handle force publication with warnings override', async () => {
      const forcePublishRequest: PublishContentRequest = {
        contentId: 'marketing.announcement.product-launch',
        changeLog: 'Force publishing despite validation warnings',
        force: true
      };

      const forcePublishResponse: PublishContentResponse = {
        ...mockPublishResponse,
        warnings: [
          'Forced publication bypassed dependency checks',
          'Content validation warnings were ignored',
          'Manual review recommended after publication'
        ]
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: forcePublishResponse
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.publishContent.initiate(forcePublishRequest)
      );

      const response = await result;

      expect(response.data).toBeDefined();
      const publishData = response.data as PublishContentResponse;
      expect(publishData.contentItem.status).toBe(ContentStatus.PUBLISHED);
      expect(publishData.warnings).toContain('Forced publication bypassed dependency checks');
      expect(publishData.warnings).toContain('Manual review recommended after publication');
    });
  });
});