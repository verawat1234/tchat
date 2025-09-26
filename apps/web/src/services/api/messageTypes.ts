// T054 - RTK Query API endpoints for new message types
/**
 * MessageTypes API Service
 * RTK Query endpoints for advanced message type operations
 * Provides CRUD operations, validation, and real-time updates
 */

import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import { MessageData, MessageType, MessageContent } from '../../types/MessageData';

// Request/Response types for API operations
export interface CreateMessageRequest {
  chatId: string;
  type: MessageType;
  content: MessageContent;
  replyToId?: string;
  threadId?: string;
  metadata?: Record<string, any>;
}

export interface CreateMessageResponse {
  message: MessageData;
  thread?: ThreadInfo;
  notifications?: NotificationInfo[];
}

export interface UpdateMessageRequest {
  messageId: string;
  content: Partial<MessageContent>;
  metadata?: Record<string, any>;
}

export interface UpdateMessageResponse {
  message: MessageData;
  affectedMessages?: MessageData[]; // For thread updates
}

export interface MessageInteractionRequest {
  messageId: string;
  interactionType: MessageInteractionType;
  data: any;
}

export interface MessageInteractionResponse {
  success: boolean;
  message: MessageData;
  result?: any; // Quiz results, poll votes, etc.
}

export interface ThreadInfo {
  id: string;
  rootMessageId: string;
  messageCount: number;
  participants: string[];
  lastActivity: string;
}

export interface NotificationInfo {
  type: 'mention' | 'reply' | 'reaction' | 'assignment';
  userId: string;
  messageId: string;
  content: string;
}

export enum MessageInteractionType {
  REACT = 'react',
  VOTE = 'vote',
  RSVP = 'rsvp',
  QUIZ_ANSWER = 'quiz_answer',
  FORM_SUBMIT = 'form_submit',
  CART_ADD = 'cart_add',
  BOOKMARK = 'bookmark',
  SHARE = 'share',
  REPORT = 'report',
  EDIT_CONTENT = 'edit_content'
}

// Validation types
export interface MessageValidationRequest {
  type: MessageType;
  content: MessageContent;
  chatId: string;
}

export interface MessageValidationResponse {
  isValid: boolean;
  errors: ValidationError[];
  warnings: ValidationWarning[];
  sanitizedContent?: MessageContent;
}

export interface ValidationError {
  field: string;
  code: string;
  message: string;
  severity: 'critical' | 'error';
}

export interface ValidationWarning {
  field: string;
  code: string;
  message: string;
  suggestion?: string;
}

// Search and filtering types
export interface MessageSearchRequest {
  chatId: string;
  query?: string;
  messageTypes?: MessageType[];
  dateRange?: {
    start: string;
    end: string;
  };
  senderIds?: string[];
  tags?: string[];
  hasAttachments?: boolean;
  limit?: number;
  offset?: number;
}

export interface MessageSearchResponse {
  messages: MessageData[];
  totalCount: number;
  facets: SearchFacets;
  suggestions?: string[];
}

export interface SearchFacets {
  messageTypes: Record<MessageType, number>;
  senders: Record<string, number>;
  dateHistogram: Record<string, number>;
  tags: Record<string, number>;
}

// Analytics types
export interface MessageAnalyticsRequest {
  chatId: string;
  messageTypes?: MessageType[];
  timeRange: {
    start: string;
    end: string;
  };
  groupBy?: 'day' | 'week' | 'month';
}

export interface MessageAnalyticsResponse {
  messageStats: MessageTypeStats[];
  engagementMetrics: EngagementMetrics;
  trends: TrendData[];
  topPerformers: TopPerformerData[];
}

export interface MessageTypeStats {
  type: MessageType;
  count: number;
  averageEngagement: number;
  completionRate?: number; // For quizzes, forms, etc.
  errorRate: number;
}

export interface EngagementMetrics {
  totalInteractions: number;
  uniqueUsers: number;
  averageResponseTime: number;
  retentionRate: number;
}

export interface TrendData {
  date: string;
  messageCount: number;
  engagementScore: number;
  completionRate: number;
}

export interface TopPerformerData {
  messageId: string;
  type: MessageType;
  engagementScore: number;
  interactionCount: number;
  completionRate?: number;
}

// Base query configuration
const baseQuery = fetchBaseQuery({
  baseUrl: '/api/v1/messages',
  prepareHeaders: (headers, { getState }) => {
    // Add authentication token
    const token = (getState() as any)?.auth?.token;
    if (token) {
      headers.set('authorization', `Bearer ${token}`);
    }
    headers.set('content-type', 'application/json');
    return headers;
  },
});

// Enhanced base query with retry and error handling
const baseQueryWithRetry = async (args: any, api: any, extraOptions: any) => {
  let result = await baseQuery(args, api, extraOptions);

  // Retry logic for transient errors
  if (result.error && 'status' in result.error) {
    const status = result.error.status;
    if (status === 429 || status >= 500) {
      // Exponential backoff for rate limits and server errors
      const retryDelay = extraOptions?.retryDelay || 1000;
      await new Promise(resolve => setTimeout(resolve, retryDelay));
      result = await baseQuery(args, api, extraOptions);
    }
  }

  return result;
};

// RTK Query API definition
export const messageTypesApi = createApi({
  reducerPath: 'messageTypesApi',
  baseQuery: baseQueryWithRetry,
  tagTypes: [
    'Message',
    'Thread',
    'MessageInteraction',
    'MessageAnalytics',
    'MessageValidation'
  ],
  endpoints: (builder) => ({
    // T055 - Create message with advanced validation
    createMessage: builder.mutation<CreateMessageResponse, CreateMessageRequest>({
      query: (body) => ({
        url: '',
        method: 'POST',
        body,
      }),
      invalidatesTags: (result, error, { chatId }) => [
        { type: 'Message', id: chatId },
        { type: 'Thread', id: 'LIST' },
      ],
      // Optimistic update
      async onQueryStarted(arg, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;
          // Update cache with new message
          dispatch(
            messageTypesApi.util.updateQueryData(
              'getMessages',
              { chatId: arg.chatId },
              (draft) => {
                draft.messages.unshift(data.message);
              }
            )
          );
        } catch {
          // Rollback optimistic update on error
        }
      },
    }),

    // T056 - Update message content
    updateMessage: builder.mutation<UpdateMessageResponse, UpdateMessageRequest>({
      query: ({ messageId, ...body }) => ({
        url: `/${messageId}`,
        method: 'PUT',
        body,
      }),
      invalidatesTags: (result, error, { messageId }) => [
        { type: 'Message', id: messageId },
      ],
    }),

    // T057 - Delete message
    deleteMessage: builder.mutation<void, { messageId: string; chatId: string }>({
      query: ({ messageId }) => ({
        url: `/${messageId}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, { messageId, chatId }) => [
        { type: 'Message', id: messageId },
        { type: 'Message', id: chatId },
      ],
    }),

    // T058 - Message interactions (reactions, votes, etc.)
    interactWithMessage: builder.mutation<MessageInteractionResponse, MessageInteractionRequest>({
      query: ({ messageId, interactionType, data }) => ({
        url: `/${messageId}/interactions`,
        method: 'POST',
        body: { interactionType, data },
      }),
      invalidatesTags: (result, error, { messageId }) => [
        { type: 'MessageInteraction', id: messageId },
        { type: 'Message', id: messageId },
      ],
      // Optimistic update for better UX
      async onQueryStarted(arg, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;
          dispatch(
            messageTypesApi.util.updateQueryData(
              'getMessage',
              { messageId: arg.messageId },
              (draft) => {
                Object.assign(draft, data.message);
              }
            )
          );
        } catch {
          // Rollback on error
        }
      },
    }),

    // T059 - Get messages with filtering and pagination
    getMessages: builder.query<
      { messages: MessageData[]; totalCount: number; nextCursor?: string },
      {
        chatId: string;
        messageTypes?: MessageType[];
        limit?: number;
        cursor?: string;
        includeThreads?: boolean;
      }
    >({
      query: ({ chatId, messageTypes, limit = 50, cursor, includeThreads }) => ({
        url: '',
        params: {
          chatId,
          messageTypes: messageTypes?.join(','),
          limit,
          cursor,
          includeThreads,
        },
      }),
      providesTags: (result, error, { chatId }) => [
        { type: 'Message', id: chatId },
        ...(result?.messages.map((msg) => ({ type: 'Message' as const, id: msg.id })) || []),
      ],
      // Transform response for better caching
      transformResponse: (response: any) => ({
        messages: response.data || [],
        totalCount: response.totalCount || 0,
        nextCursor: response.nextCursor,
      }),
    }),

    // T060 - Get single message with thread context
    getMessage: builder.query<
      MessageData & { threadContext?: MessageData[] },
      { messageId: string; includeThread?: boolean }
    >({
      query: ({ messageId, includeThread }) => ({
        url: `/${messageId}`,
        params: { includeThread },
      }),
      providesTags: (result, error, { messageId }) => [
        { type: 'Message', id: messageId },
      ],
    }),

    // T061 - Validate message before creation
    validateMessage: builder.mutation<MessageValidationResponse, MessageValidationRequest>({
      query: (body) => ({
        url: '/validate',
        method: 'POST',
        body,
      }),
      // Don't cache validation responses
      providesTags: [],
    }),

    // T062 - Search messages with advanced filtering
    searchMessages: builder.query<MessageSearchResponse, MessageSearchRequest>({
      query: (params) => ({
        url: '/search',
        params: {
          ...params,
          messageTypes: params.messageTypes?.join(','),
          senderIds: params.senderIds?.join(','),
          tags: params.tags?.join(','),
        },
      }),
      providesTags: ['Message'],
      // Keep search results cached for 5 minutes
      keepUnusedDataFor: 300,
    }),

    // T063 - Get message analytics
    getMessageAnalytics: builder.query<MessageAnalyticsResponse, MessageAnalyticsRequest>({
      query: (params) => ({
        url: '/analytics',
        params: {
          ...params,
          messageTypes: params.messageTypes?.join(','),
        },
      }),
      providesTags: ['MessageAnalytics'],
      // Cache analytics for 1 hour
      keepUnusedDataFor: 3600,
    }),

    // T064 - Get thread information
    getThread: builder.query<
      { thread: ThreadInfo; messages: MessageData[] },
      { threadId: string; limit?: number; offset?: number }
    >({
      query: ({ threadId, limit = 20, offset = 0 }) => ({
        url: `/threads/${threadId}`,
        params: { limit, offset },
      }),
      providesTags: (result, error, { threadId }) => [
        { type: 'Thread', id: threadId },
      ],
    }),

    // T065 - Bulk operations for message management
    bulkUpdateMessages: builder.mutation<
      { updatedCount: number; errors: Array<{ messageId: string; error: string }> },
      {
        messageIds: string[];
        updates: Partial<MessageContent>;
        operation: 'update' | 'delete' | 'archive';
      }
    >({
      query: (body) => ({
        url: '/bulk',
        method: 'POST',
        body,
      }),
      invalidatesTags: (result, error, { messageIds }) => [
        ...messageIds.map((id) => ({ type: 'Message' as const, id })),
        { type: 'Message', id: 'LIST' },
      ],
    }),
  }),
});

// Export hooks for components
export const {
  useCreateMessageMutation,
  useUpdateMessageMutation,
  useDeleteMessageMutation,
  useInteractWithMessageMutation,
  useGetMessagesQuery,
  useGetMessageQuery,
  useValidateMessageMutation,
  useSearchMessagesQuery,
  useLazySearchMessagesQuery,
  useGetMessageAnalyticsQuery,
  useGetThreadQuery,
  useBulkUpdateMessagesMutation,
} = messageTypesApi;

// Export API for store configuration
export default messageTypesApi;