/**
 * Messages API - RTK Query endpoints through Gateway
 *
 * Gateway-integrated messaging API with comprehensive CRUD operations,
 * real-time subscriptions, and cross-platform synchronization.
 *
 * Features:
 * - Complete message CRUD through gateway routing
 * - Real-time message subscriptions via WebSocket
 * - File upload integration with progress tracking
 * - Message threading and reply functionality
 * - Delivery receipts and read status tracking
 * - Message search and filtering
 * - Typing indicators and presence
 */

import { api } from './api';
import type {
  Message,
  SendMessageRequest,
  MessageSearchQuery,
  MessageSearchResult,
  MessageThread,
  TypingIndicator,
  MessageAttachment,
} from './messageService';
import type { PaginatedResponse } from '../types/content';

// =============================================================================
// Request/Response Type Definitions
// =============================================================================

export interface GetMessagesRequest {
  chatId: string;
  limit?: number;
  before?: string; // Message ID for pagination
  after?: string;  // Message ID for reverse pagination
  includeThreads?: boolean;
}

export interface GetMessagesResponse extends PaginatedResponse<Message> {
  threads?: MessageThread[];
}

export interface UpdateMessageRequest {
  content?: {
    text?: string;
    html?: string;
    data?: any;
  };
  status?: 'read' | 'delivered' | 'failed';
}

export interface AddReactionRequest {
  emoji: string;
  action: 'add' | 'remove';
}

export interface UploadAttachmentRequest {
  file: File;
  messageId?: string; // For attaching to existing message
  metadata?: {
    description?: string;
    tags?: string[];
  };
}

export interface UploadAttachmentResponse {
  attachment: MessageAttachment;
  uploadUrl?: string; // Presigned URL if needed
}

export interface MessageDeliveryReceipt {
  messageId: string;
  chatId: string;
  userId: string;
  status: 'delivered' | 'read';
  timestamp: string;
}

export interface CreateThreadRequest {
  parentMessageId: string;
  chatId: string;
  initialMessage?: SendMessageRequest;
}

// =============================================================================
// Messages API Endpoints
// =============================================================================

export const messagesApi = api.injectEndpoints({
  endpoints: (builder) => ({
    // ========================================================================
    // Query Endpoints - Message Retrieval
    // ========================================================================

    /**
     * Get messages for a chat with pagination and threading support
     */
    getMessages: builder.query<GetMessagesResponse, GetMessagesRequest>({
      query: ({ chatId, ...params }) => ({
        url: `/messages/chat/${encodeURIComponent(chatId)}`,
        method: 'GET',
        params,
      }),
      providesTags: (result, error, { chatId }) => [
        { type: 'Message', id: 'LIST' },
        { type: 'ChatMessages', id: chatId },
        'Message',
      ],
      keepUnusedDataFor: 300, // 5 minutes
    }),

    /**
     * Get a single message by ID
     */
    getMessage: builder.query<Message, string>({
      query: (messageId) => `/messages/${encodeURIComponent(messageId)}`,
      providesTags: (result, error, messageId) => [
        { type: 'Message', id: messageId },
        'Message',
      ],
    }),

    /**
     * Search messages across chats
     */
    searchMessages: builder.query<MessageSearchResult, MessageSearchQuery>({
      query: (params) => ({
        url: '/messages/search',
        method: 'GET',
        params: {
          ...params,
          // Convert array params to comma-separated strings
          type: params.type ? [params.type] : undefined,
        },
      }),
      providesTags: ['Message'],
      keepUnusedDataFor: 60, // 1 minute for search results
    }),

    /**
     * Get message threads for a chat
     */
    getMessageThreads: builder.query<MessageThread[], string>({
      query: (chatId) => `/messages/chat/${encodeURIComponent(chatId)}/threads`,
      providesTags: (result, error, chatId) => [
        { type: 'MessageThread', id: 'LIST' },
        { type: 'ChatThreads', id: chatId },
      ],
    }),

    /**
     * Get messages in a specific thread
     */
    getThreadMessages: builder.query<GetMessagesResponse, { threadId: string; limit?: number; before?: string }>({
      query: ({ threadId, ...params }) => ({
        url: `/messages/thread/${encodeURIComponent(threadId)}`,
        method: 'GET',
        params,
      }),
      providesTags: (result, error, { threadId }) => [
        { type: 'ThreadMessages', id: threadId },
        'Message',
      ],
    }),

    /**
     * Get unread message count for user
     */
    getUnreadCount: builder.query<{ total: number; byChatId: Record<string, number> }, void>({
      query: () => '/messages/unread/count',
      providesTags: ['UnreadCount'],
      keepUnusedDataFor: 30, // 30 seconds
    }),

    // ========================================================================
    // Mutation Endpoints - Message Management
    // ========================================================================

    /**
     * Send a new message
     */
    sendMessage: builder.mutation<Message, SendMessageRequest>({
      query: (data) => ({
        url: '/messages',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: (result, error, { chatId, threadId }) => [
        { type: 'ChatMessages', id: chatId },
        { type: 'Message', id: 'LIST' },
        'UnreadCount',
        ...(threadId ? [{ type: 'ThreadMessages', id: threadId }] : []),
      ],
      // Optimistic updates
      async onQueryStarted(request, { dispatch, queryFulfilled }) {
        // Create optimistic message
        const optimisticMessage: Message = {
          id: `temp_${Date.now()}`,
          chatId: request.chatId,
          senderId: 'current-user',
          type: request.type,
          content: request.content,
          timestamp: new Date().toISOString(),
          status: 'pending' as any,
          threadId: request.threadId,
          replyToId: request.replyToId,
          reactions: [],
          metadata: {
            version: 1,
            clientId: 'web-client',
            priority: 'normal',
          },
          mentions: request.mentions || [],
          isEdited: false,
          isDeleted: false,
        };

        // Optimistically update the messages list
        const patchResult = dispatch(
          messagesApi.util.updateQueryData('getMessages', { chatId: request.chatId }, (draft) => {
            draft.items.unshift(optimisticMessage);
            draft.total++;
          })
        );

        try {
          const { data: sentMessage } = await queryFulfilled;

          // Replace optimistic message with real message
          dispatch(
            messagesApi.util.updateQueryData('getMessages', { chatId: request.chatId }, (draft) => {
              const index = draft.items.findIndex(m => m.id === optimisticMessage.id);
              if (index !== -1) {
                draft.items[index] = sentMessage;
              }
            })
          );
        } catch {
          // Remove optimistic message on failure
          patchResult.undo();
        }
      },
    }),

    /**
     * Update an existing message (edit)
     */
    updateMessage: builder.mutation<Message, { messageId: string; updates: UpdateMessageRequest }>({
      query: ({ messageId, updates }) => ({
        url: `/messages/${encodeURIComponent(messageId)}`,
        method: 'PUT',
        body: updates,
      }),
      invalidatesTags: (result, error, { messageId }) => [
        { type: 'Message', id: messageId },
        'Message',
      ],
      // Optimistic updates
      async onQueryStarted({ messageId, updates }, { dispatch, queryFulfilled }) {
        const patches: any[] = [];

        // Update in all relevant queries
        patches.push(
          dispatch(
            messagesApi.util.updateQueryData('getMessage', messageId, (draft) => {
              if (updates.content) {
                Object.assign(draft.content, updates.content);
                draft.isEdited = true;
                draft.editedAt = new Date().toISOString();
              }
            })
          )
        );

        try {
          await queryFulfilled;
        } catch {
          patches.forEach(patch => patch.undo());
        }
      },
    }),

    /**
     * Delete a message
     */
    deleteMessage: builder.mutation<void, string>({
      query: (messageId) => ({
        url: `/messages/${encodeURIComponent(messageId)}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, messageId) => [
        { type: 'Message', id: messageId },
        'Message',
      ],
    }),

    /**
     * Add or remove reaction to/from message
     */
    toggleReaction: builder.mutation<Message, { messageId: string; reaction: AddReactionRequest }>({
      query: ({ messageId, reaction }) => ({
        url: `/messages/${encodeURIComponent(messageId)}/reactions`,
        method: 'POST',
        body: reaction,
      }),
      invalidatesTags: (result, error, { messageId }) => [
        { type: 'Message', id: messageId },
      ],
      // Optimistic updates for reactions
      async onQueryStarted({ messageId, reaction }, { dispatch, queryFulfilled }) {
        const patchResult = dispatch(
          messagesApi.util.updateQueryData('getMessage', messageId, (draft) => {
            const existingReactionIndex = draft.reactions.findIndex(
              r => r.userId === 'current-user' && r.emoji === reaction.emoji
            );

            if (reaction.action === 'add' && existingReactionIndex === -1) {
              draft.reactions.push({
                id: `temp_${Date.now()}`,
                emoji: reaction.emoji,
                userId: 'current-user',
                userName: 'Current User',
                timestamp: new Date().toISOString(),
              });
            } else if (reaction.action === 'remove' && existingReactionIndex !== -1) {
              draft.reactions.splice(existingReactionIndex, 1);
            }
          })
        );

        try {
          await queryFulfilled;
        } catch {
          patchResult.undo();
        }
      },
    }),

    /**
     * Mark messages as read
     */
    markAsRead: builder.mutation<void, { messageIds: string[]; chatId: string }>({
      query: ({ messageIds, chatId }) => ({
        url: `/messages/read`,
        method: 'POST',
        body: { messageIds, chatId },
      }),
      invalidatesTags: (result, error, { chatId }) => [
        { type: 'ChatMessages', id: chatId },
        'UnreadCount',
      ],
    }),

    /**
     * Send typing indicator
     */
    sendTypingIndicator: builder.mutation<void, { chatId: string; isTyping: boolean }>({
      query: (data) => ({
        url: '/messages/typing',
        method: 'POST',
        body: data,
      }),
      // No cache invalidation needed for typing indicators
    }),

    // ========================================================================
    // File and Attachment Endpoints
    // ========================================================================

    /**
     * Upload message attachment
     */
    uploadAttachment: builder.mutation<UploadAttachmentResponse, UploadAttachmentRequest>({
      query: ({ file, messageId, metadata }) => {
        const formData = new FormData();
        formData.append('file', file);
        if (messageId) formData.append('messageId', messageId);
        if (metadata) formData.append('metadata', JSON.stringify(metadata));

        return {
          url: '/messages/attachments',
          method: 'POST',
          body: formData,
          formData: true,
        };
      },
    }),

    /**
     * Get attachment download URL
     */
    getAttachmentUrl: builder.query<{ url: string; expiresAt: string }, string>({
      query: (attachmentId) => `/messages/attachments/${encodeURIComponent(attachmentId)}/download`,
      keepUnusedDataFor: 300, // 5 minutes
    }),

    // ========================================================================
    // Threading Endpoints
    // ========================================================================

    /**
     * Create a new message thread
     */
    createThread: builder.mutation<MessageThread, CreateThreadRequest>({
      query: (data) => ({
        url: '/messages/threads',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: (result, error, { chatId }) => [
        { type: 'ChatThreads', id: chatId },
        { type: 'MessageThread', id: 'LIST' },
      ],
    }),

    /**
     * Close/archive a thread
     */
    closeThread: builder.mutation<void, string>({
      query: (threadId) => ({
        url: `/messages/threads/${encodeURIComponent(threadId)}/close`,
        method: 'POST',
      }),
      invalidatesTags: (result, error, threadId) => [
        { type: 'MessageThread', id: threadId },
        'MessageThread',
      ],
    }),

    // ========================================================================
    // Delivery Receipt Endpoints
    // ========================================================================

    /**
     * Send delivery receipt for a message
     */
    sendDeliveryReceipt: builder.mutation<void, { messageId: string; chatId: string; status: 'delivered' | 'read' }>({
      query: (data) => ({
        url: '/messages/receipts',
        method: 'POST',
        body: data,
      }),
      // No cache invalidation needed for receipts
    }),

    /**
     * Send bulk delivery receipts
     */
    sendBulkDeliveryReceipts: builder.mutation<void, { chatId: string; messageIds: string[]; status: 'delivered' | 'read' }>({
      query: (data) => ({
        url: '/messages/receipts/bulk',
        method: 'POST',
        body: data,
      }),
    }),

    /**
     * Get delivery receipts for a message
     */
    getDeliveryReceipts: builder.query<Array<{ userId: string; userName: string; status: string; timestamp: string }>, string>({
      query: (messageId) => `/messages/${encodeURIComponent(messageId)}/receipts`,
      providesTags: (result, error, messageId) => [
        { type: 'DeliveryReceipt', id: messageId },
      ],
    }),

    /**
     * Get read status for multiple messages
     */
    getBulkReadStatus: builder.query<Record<string, 'sent' | 'delivered' | 'read'>, { messageIds: string[] }>({
      query: ({ messageIds }) => ({
        url: '/messages/receipts/bulk-status',
        method: 'POST',
        body: { messageIds },
      }),
      providesTags: ['DeliveryReceipt'],
    }),
  }),
  overrideExisting: false,
});

// Export hooks for all message API endpoints
export const {
  // Query hooks
  useGetMessagesQuery,
  useGetMessageQuery,
  useSearchMessagesQuery,
  useGetMessageThreadsQuery,
  useGetThreadMessagesQuery,
  useGetUnreadCountQuery,
  useGetAttachmentUrlQuery,
  useGetDeliveryReceiptsQuery,
  useGetBulkReadStatusQuery,

  // Mutation hooks
  useSendMessageMutation,
  useUpdateMessageMutation,
  useDeleteMessageMutation,
  useToggleReactionMutation,
  useMarkAsReadMutation,
  useSendTypingIndicatorMutation,
  useUploadAttachmentMutation,
  useCreateThreadMutation,
  useCloseThreadMutation,
  useSendDeliveryReceiptMutation,
  useSendBulkDeliveryReceiptsMutation,
} = messagesApi;