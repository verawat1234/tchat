import { api } from './api';
import type {
  Chat,
  ApiResponse,
  PaginatedResponse,
} from '../types/api';

interface ListChatsParams {
  page?: number;
  limit?: number;
  type?: 'direct' | 'group' | 'channel';
}

interface CreateChatRequest {
  name: string;
  type: 'direct' | 'group' | 'channel';
  participants: string[];
}

interface UpdateChatRequest {
  name?: string;
  participants?: string[];
}

export const chatsApi = api.injectEndpoints({
  endpoints: (builder) => ({
    listChats: builder.query<Chat[], ListChatsParams>({
      query: ({ page = 1, limit = 20, type } = {}) => {
        const params = new URLSearchParams();
        params.append('page', page.toString());
        params.append('limit', limit.toString());
        if (type) params.append('type', type);
        
        return `/chats?${params.toString()}`;
      },
      providesTags: (result) =>
        result
          ? [
              ...result.map(({ id }) => ({ type: 'Chat' as const, id })),
              { type: 'Chat', id: 'LIST' },
            ]
          : [{ type: 'Chat', id: 'LIST' }],
      transformResponse: (response: ApiResponse<Chat[]>) => response.data,
    }),

    getChatById: builder.query<Chat, string>({
      query: (id) => `/chats/${id}`,
      providesTags: (result, error, id) => [{ type: 'Chat', id }],
      transformResponse: (response: ApiResponse<Chat>) => response.data,
    }),

    createChat: builder.mutation<Chat, CreateChatRequest>({
      query: (chat) => ({
        url: '/chats',
        method: 'POST',
        body: chat,
      }),
      invalidatesTags: [{ type: 'Chat', id: 'LIST' }],
      transformResponse: (response: ApiResponse<Chat>) => response.data,
    }),

    updateChat: builder.mutation<Chat, { id: string; data: UpdateChatRequest }>({
      query: ({ id, data }) => ({
        url: `/chats/${id}`,
        method: 'PATCH',
        body: data,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'Chat', id },
        { type: 'Chat', id: 'LIST' },
      ],
      transformResponse: (response: ApiResponse<Chat>) => response.data,
      async onQueryStarted({ id, data }, { dispatch, queryFulfilled }) {
        // Optimistic update
        const patchResult = dispatch(
          api.util.updateQueryData('getChatById', id, (draft) => {
            Object.assign(draft, data);
          })
        );
        try {
          await queryFulfilled;
        } catch {
          patchResult.undo();
        }
      },
    }),

    deleteChat: builder.mutation<void, string>({
      query: (id) => ({
        url: `/chats/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, id) => [
        { type: 'Chat', id },
        { type: 'Chat', id: 'LIST' },
        { type: 'Message', id: `CHAT-${id}` },
      ],
    }),

    markChatAsRead: builder.mutation<void, string>({
      query: (chatId) => ({
        url: `/chats/${chatId}/read`,
        method: 'POST',
      }),
      invalidatesTags: (result, error, chatId) => [{ type: 'Chat', id: chatId }],
      async onQueryStarted(chatId, { dispatch, queryFulfilled }) {
        // Optimistically set unread count to 0
        const patchResult = dispatch(
          api.util.updateQueryData('getChatById', chatId, (draft) => {
            draft.unreadCount = 0;
          })
        );
        try {
          await queryFulfilled;
        } catch {
          patchResult.undo();
        }
      },
    }),
  }),
});

export const {
  useListChatsQuery,
  useLazyListChatsQuery,
  useGetChatByIdQuery,
  useLazyGetChatByIdQuery,
  useCreateChatMutation,
  useUpdateChatMutation,
  useDeleteChatMutation,
  useMarkChatAsReadMutation,
} = chatsApi;