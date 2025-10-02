import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

// Types
export interface LiveStream {
  id: string;
  broadcaster_id: string;
  stream_type: 'store' | 'video';
  title: string;
  description?: string;
  broadcaster_kyc_tier: number;
  status: 'scheduled' | 'live' | 'ended';
  stream_key: string;
  viewer_count: number;
  peak_viewer_count: number;
  max_capacity: number;
  total_view_time_seconds: number;
  average_watch_time_seconds?: number;
  scheduled_start_time?: string;
  actual_start_time?: string;
  end_time?: string;
  simulcast_layers: string[];
  current_bitrate_kbps?: number;
  recording_url?: string;
  thumbnail_url?: string;
  featured_products: string[];
  stream_settings: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}

export interface CreateStreamRequest {
  stream_type: 'store' | 'video';
  title: string;
  description?: string;
  scheduled_start_time?: string;
  max_capacity?: number;
  stream_settings?: Record<string, unknown>;
}

export interface UpdateStreamRequest {
  title?: string;
  description?: string;
  scheduled_start_time?: string;
  max_capacity?: number;
  stream_settings?: Record<string, unknown>;
}

export interface StartStreamRequest {
  offer: RTCSessionDescriptionInit;
  simulcast_enabled: boolean;
}

export interface StartStreamResponse {
  answer: RTCSessionDescriptionInit;
  ice_servers: RTCIceServer[];
  stream_id: string;
  status: string;
  started_at: string;
}

export interface ListStreamsParams {
  stream_type?: 'store' | 'video';
  status?: 'scheduled' | 'live' | 'ended';
  broadcaster_id?: string;
  limit?: number;
  offset?: number;
}

export interface ListStreamsResponse {
  streams: LiveStream[];
  total_count: number;
  limit: number;
  offset: number;
}

// RTK Query API
export const streamingApi = createApi({
  reducerPath: 'streamingApi',
  baseQuery: fetchBaseQuery({
    baseUrl: '/api/v1',
    prepareHeaders: (headers) => {
      // Add JWT token from localStorage or Redux state
      const token = localStorage.getItem('auth_token');
      if (token) {
        headers.set('Authorization', `Bearer ${token}`);
      }
      return headers;
    },
  }),
  tagTypes: ['Stream', 'StreamList'],
  endpoints: (builder) => ({
    // Create stream
    createStream: builder.mutation<LiveStream, CreateStreamRequest>({
      query: (body) => ({
        url: '/streams',
        method: 'POST',
        body,
      }),
      invalidatesTags: ['StreamList'],
    }),

    // List streams with filters
    listStreams: builder.query<ListStreamsResponse, ListStreamsParams | void>({
      query: (params = {}) => ({
        url: '/streams',
        method: 'GET',
        params,
      }),
      providesTags: (result) =>
        result
          ? [
              ...result.streams.map(({ id }) => ({ type: 'Stream' as const, id })),
              { type: 'StreamList' as const },
            ]
          : [{ type: 'StreamList' as const }],
    }),

    // Get stream by ID
    getStream: builder.query<LiveStream, string>({
      query: (streamId) => ({
        url: `/streams/${streamId}`,
        method: 'GET',
      }),
      providesTags: (result, error, streamId) => [{ type: 'Stream', id: streamId }],
    }),

    // Update stream
    updateStream: builder.mutation<LiveStream, { streamId: string; data: UpdateStreamRequest }>({
      query: ({ streamId, data }) => ({
        url: `/streams/${streamId}`,
        method: 'PATCH',
        body: data,
      }),
      invalidatesTags: (result, error, { streamId }) => [
        { type: 'Stream', id: streamId },
        { type: 'StreamList' },
      ],
    }),

    // Start stream (WebRTC negotiation)
    startStream: builder.mutation<StartStreamResponse, { streamId: string; data: StartStreamRequest }>({
      query: ({ streamId, data }) => ({
        url: `/streams/${streamId}/start`,
        method: 'POST',
        body: data,
      }),
      invalidatesTags: (result, error, { streamId }) => [
        { type: 'Stream', id: streamId },
        { type: 'StreamList' },
      ],
    }),

    // End stream
    endStream: builder.mutation<
      {
        stream_id: string;
        status: string;
        ended_at: string;
        duration_seconds: number;
        peak_viewer_count: number;
        recording_url: string;
      },
      string
    >({
      query: (streamId) => ({
        url: `/streams/${streamId}/end`,
        method: 'POST',
      }),
      invalidatesTags: (result, error, streamId) => [
        { type: 'Stream', id: streamId },
        { type: 'StreamList' },
      ],
    }),
  }),
});

// Export hooks
export const {
  useCreateStreamMutation,
  useListStreamsQuery,
  useGetStreamQuery,
  useUpdateStreamMutation,
  useStartStreamMutation,
  useEndStreamMutation,
} = streamingApi;