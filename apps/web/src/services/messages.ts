import { api } from './api';
import type {
  Message,
  CreateMessageRequest,
  PaginatedResponse,
  ApiResponse,
} from '../types/api';

interface ListMessagesParams {
  chatId: string;
  cursor?: string;
  limit?: number;
}

export const messagesApi = api.injectEndpoints({
  endpoints: (builder) => ({
    listMessages: builder.query<PaginatedResponse<Message>, ListMessagesParams>({
      query: ({ chatId, cursor, limit = 20 }) => {
        const params = new URLSearchParams();
        params.append('chatId', chatId);
        params.append('limit', limit.toString());
        if (cursor) params.append('cursor', cursor);
        
        return `/messages?${params.toString()}`;
      },
      providesTags: (result, error, { chatId }) =>
        result
          ? [
              ...result.items.map(({ id }) => ({ type: 'Message' as const, id })),
              { type: 'Message', id: `CHAT-${chatId}` },
            ]
          : [{ type: 'Message', id: `CHAT-${chatId}` }],
      serializeQueryArgs: ({ queryArgs }) => {
        const { chatId } = queryArgs;
        return { chatId }; // Omit cursor to merge paginated results
      },
      merge: (currentCache, newItems, { arg }) => {
        if (!arg.cursor) {
          // First page, replace everything
          return newItems;
        }
        // Append new items for pagination
        return {
          ...newItems,
          items: [...currentCache.items, ...newItems.items],
        };
      },
      forceRefetch({ currentArg, previousArg }) {
        return currentArg?.cursor !== previousArg?.cursor;
      },
    }),

    sendMessage: builder.mutation<Message, CreateMessageRequest>({
      query: (message) => ({
        url: '/messages',
        method: 'POST',
        body: message,
      }),
      invalidatesTags: (result, error, { chatId }) => [
        { type: 'Message', id: `CHAT-${chatId}` },
        { type: 'Chat', id: chatId },
      ],
      transformResponse: (response: ApiResponse<Message>) => response.data,
      async onQueryStarted(arg, { dispatch, queryFulfilled }) {
        // Optimistic update
        const tempId = `temp-${Date.now()}`;
        const patchResult = dispatch(
          api.util.updateQueryData('listMessages', { chatId: arg.chatId }, (draft) => {
            const tempMessage: Message = {
              id: tempId,
              chatId: arg.chatId,
              userId: 'current-user', // Will be replaced with actual user ID
              content: arg.content,
              type: arg.type,
              attachments: arg.attachments,
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            };
            draft.items.push(tempMessage);
          })
        );

        try {
          const { data } = await queryFulfilled;
          // Replace temp message with real one
          dispatch(
            api.util.updateQueryData('listMessages', { chatId: arg.chatId }, (draft) => {
              const index = draft.items.findIndex(msg => msg.id === tempId);
              if (index !== -1) {
                draft.items[index] = data;
              }
            })
          );
        } catch {
          patchResult.undo();
        }
      },
    }),

    updateMessage: builder.mutation<Message, { id: string; content: string }>({
      query: ({ id, content }) => ({
        url: `/messages/${id}`,
        method: 'PATCH',
        body: { content },
      }),
      invalidatesTags: (result, error, { id }) => [{ type: 'Message', id }],
      transformResponse: (response: ApiResponse<Message>) => response.data,
      async onQueryStarted({ id, content }, { dispatch, queryFulfilled }) {
        // Find and update the message optimistically
        const patches: any[] = [];
        
        // We don't know which chat this message belongs to, so we need to search
        // In a real app, you'd pass chatId as well
        const updateAllCaches = () => {
          // This is a simplified version - in production you'd track which caches to update
        };
        
        try {
          await queryFulfilled;
        } catch {
          patches.forEach(patch => patch.undo());
        }
      },
    }),

    deleteMessage: builder.mutation<void, string>({
      query: (id) => ({
        url: `/messages/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, id) => [{ type: 'Message', id }],
    }),
  }),
});

export const {
  useListMessagesQuery,
  useLazyListMessagesQuery,
  useSendMessageMutation,
  useUpdateMessageMutation,
  useDeleteMessageMutation,
} = messagesApi;

// Utility function for infinite scroll
export const useInfiniteMessages = (chatId: string) => {
  const [trigger, result] = messagesApi.useLazyListMessagesQuery();
  
  const loadMore = async (cursor?: string) => {
    await trigger({ chatId, cursor });
  };
  
  return { loadMore, ...result };
};