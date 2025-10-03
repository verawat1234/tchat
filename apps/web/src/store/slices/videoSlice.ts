/**
 * Video Redux Slice
 *
 * Centralized state management for video functionality including:
 * - Video content management
 * - Playback state tracking
 * - Upload progress monitoring
 * - Cross-platform synchronization status
 * - Error handling and loading states
 */

import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import type {
  VideoContent,
  VideoState,
  VideoError,
  PlaybackPositionUpdate,
  PlaybackQualityUpdate,
  PlaybackStateUpdate,
  BufferHealthUpdate,
  PlaybackState,
  SyncStatus,
} from '../../types/video';
import {
  uploadVideoToBackend,
  getBackendStreamUrl,
  syncBackendPlaybackPosition,
  createBackendSyncSession,
} from '../../services/videoApi';

// =============================================================================
// Initial State
// =============================================================================

const initialState: VideoState = {
  videos: [],
  currentVideo: null,
  playbackState: {},
  uploadProgress: {},
  syncStatus: {},
  loading: false,
  error: null,
};

// =============================================================================
// Async Thunks
// =============================================================================

/**
 * Upload video with progress tracking
 */
export const uploadVideo = createAsyncThunk(
  'video/upload',
  async (
    {
      file,
      metadata,
      onProgress,
    }: {
      file: File;
      metadata: {
        title: string;
        description: string;
        tags: string[];
        content_rating: string;
        thumbnail?: File;
        category?: string;
        is_monetized?: boolean;
        price?: number;
      };
      onProgress?: (progress: number) => void;
    },
    { rejectWithValue }
  ) => {
    try {
      const result = await uploadVideoToBackend(file, metadata, onProgress);
      return result;
    } catch (error) {
      return rejectWithValue({
        code: 'UPLOAD_FAILED',
        message: error instanceof Error ? error.message : 'Upload failed',
      });
    }
  }
);

/**
 * Fetch video details
 */
export const fetchVideo = createAsyncThunk(
  'video/fetchVideo',
  async (videoId: string, { rejectWithValue }) => {
    try {
      const response = await fetch(`/api/v1/video/${videoId}`);
      if (!response.ok) throw new Error('Failed to fetch video');
      return await response.json();
    } catch (error) {
      return rejectWithValue({
        code: 'FETCH_FAILED',
        message: error instanceof Error ? error.message : 'Failed to fetch video',
      });
    }
  }
);

/**
 * Fetch videos list
 */
export const fetchVideos = createAsyncThunk(
  'video/fetchVideos',
  async (
    params: {
      creator_id?: string;
      category?: string;
      page?: number;
      limit?: number;
    },
    { rejectWithValue }
  ) => {
    try {
      const queryParams = new URLSearchParams();
      if (params.creator_id) queryParams.append('creator_id', params.creator_id);
      if (params.category) queryParams.append('category', params.category);
      queryParams.append('page', String(params.page || 1));
      queryParams.append('limit', String(params.limit || 20));

      const response = await fetch(`/api/v1/video?${queryParams.toString()}`);
      if (!response.ok) throw new Error('Failed to fetch videos');

      const data = await response.json();
      return data.videos || [];
    } catch (error) {
      return rejectWithValue({
        code: 'FETCH_FAILED',
        message: error instanceof Error ? error.message : 'Failed to fetch videos',
      });
    }
  }
);

/**
 * Initialize playback session with sync
 */
export const initializePlaybackSession = createAsyncThunk(
  'video/initializePlaybackSession',
  async (
    {
      videoId,
      userId,
      platform,
      initialPosition,
    }: {
      videoId: string;
      userId: string;
      platform: string;
      initialPosition?: number;
    },
    { rejectWithValue }
  ) => {
    try {
      const session = await createBackendSyncSession(
        videoId,
        userId,
        platform,
        initialPosition || 0
      );
      return { videoId, session };
    } catch (error) {
      return rejectWithValue({
        code: 'SESSION_INIT_FAILED',
        message: error instanceof Error ? error.message : 'Failed to initialize session',
      });
    }
  }
);

/**
 * Sync playback position
 */
export const syncPlayback = createAsyncThunk(
  'video/syncPlayback',
  async (
    {
      videoId,
      sessionId,
      position,
      platform,
      playbackState,
    }: {
      videoId: string;
      sessionId: string;
      position: number;
      platform: string;
      playbackState?: string;
    },
    { rejectWithValue }
  ) => {
    try {
      const result = await syncBackendPlaybackPosition(
        videoId,
        sessionId,
        position,
        platform,
        playbackState
      );
      return { videoId, ...result };
    } catch (error) {
      return rejectWithValue({
        code: 'SYNC_FAILED',
        message: error instanceof Error ? error.message : 'Sync failed',
      });
    }
  }
);

// =============================================================================
// Video Slice
// =============================================================================

const videoSlice = createSlice({
  name: 'video',
  initialState,
  reducers: {
    // Set current video
    setCurrentVideo: (state, action: PayloadAction<VideoContent | null>) => {
      state.currentVideo = action.payload;
    },

    // Add videos to list
    addVideos: (state, action: PayloadAction<VideoContent[]>) => {
      const existingIds = new Set(state.videos.map(v => v.id));
      const newVideos = action.payload.filter(v => !existingIds.has(v.id));
      state.videos.push(...newVideos);
    },

    // Update playback position
    updatePlaybackPosition: (state, action: PayloadAction<PlaybackPositionUpdate>) => {
      const { videoId, position, sessionId } = action.payload;
      if (!state.playbackState[videoId]) {
        state.playbackState[videoId] = {
          position: 0,
          state: 'paused',
          quality: 'auto',
        };
      }
      state.playbackState[videoId].position = position;
      if (sessionId) {
        state.playbackState[videoId].sessionId = sessionId;
      }
    },

    // Update playback quality
    updatePlaybackQuality: (state, action: PayloadAction<PlaybackQualityUpdate>) => {
      const { videoId, quality } = action.payload;
      if (!state.playbackState[videoId]) {
        state.playbackState[videoId] = {
          position: 0,
          state: 'paused',
          quality: 'auto',
        };
      }
      state.playbackState[videoId].quality = quality;
    },

    // Set playback state
    setPlaybackState: (state, action: PayloadAction<PlaybackStateUpdate>) => {
      const { videoId, state: playbackState } = action.payload;
      if (!state.playbackState[videoId]) {
        state.playbackState[videoId] = {
          position: 0,
          state: 'paused',
          quality: 'auto',
        };
      }
      state.playbackState[videoId].state = playbackState;
    },

    // Update buffer health
    updateBufferHealth: (state, action: PayloadAction<BufferHealthUpdate>) => {
      const { videoId, bufferedSeconds, bufferPercentage } = action.payload;
      // Store buffer health in sync status for monitoring
      if (!state.syncStatus[videoId]) {
        state.syncStatus[videoId] = {
          last_sync_time: new Date().toISOString(),
          synced_platforms: [],
          conflict_detected: false,
        };
      }
    },

    // Set upload progress
    setUploadProgress: (state, action: PayloadAction<{ fileName: string; progress: number }>) => {
      const { fileName, progress } = action.payload;
      state.uploadProgress[fileName] = progress;
    },

    // Clear upload progress
    clearUploadProgress: (state, action: PayloadAction<string>) => {
      delete state.uploadProgress[action.payload];
    },

    // Update sync status
    updateSyncStatus: (state, action: PayloadAction<{ videoId: string; status: SyncStatus }>) => {
      const { videoId, status } = action.payload;
      state.syncStatus[videoId] = status;
    },

    // Clear error
    clearError: (state) => {
      state.error = null;
    },

    // Reset video state
    resetVideoState: (state) => {
      state.videos = [];
      state.currentVideo = null;
      state.playbackState = {};
      state.uploadProgress = {};
      state.syncStatus = {};
      state.loading = false;
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    // Upload video
    builder
      .addCase(uploadVideo.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(uploadVideo.fulfilled, (state, action) => {
        state.loading = false;
        // Clear upload progress for this file
        // Note: fileName not available in response, handled by setUploadProgress
      })
      .addCase(uploadVideo.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as VideoError;
      });

    // Fetch video
    builder
      .addCase(fetchVideo.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchVideo.fulfilled, (state, action) => {
        state.loading = false;
        state.currentVideo = action.payload;

        // Update or add to videos list
        const index = state.videos.findIndex(v => v.id === action.payload.id);
        if (index >= 0) {
          state.videos[index] = action.payload;
        } else {
          state.videos.push(action.payload);
        }
      })
      .addCase(fetchVideo.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as VideoError;
      });

    // Fetch videos
    builder
      .addCase(fetchVideos.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchVideos.fulfilled, (state, action) => {
        state.loading = false;
        state.videos = action.payload;
      })
      .addCase(fetchVideos.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as VideoError;
      });

    // Initialize playback session
    builder
      .addCase(initializePlaybackSession.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(initializePlaybackSession.fulfilled, (state, action) => {
        state.loading = false;
        const { videoId, session } = action.payload;

        // Initialize playback state with session
        state.playbackState[videoId] = {
          position: session.initial_position,
          state: 'paused',
          quality: 'auto',
          sessionId: session.session_id,
        };

        // Initialize sync status
        state.syncStatus[videoId] = {
          last_sync_time: session.created_at,
          synced_platforms: [session.platform],
          conflict_detected: false,
        };
      })
      .addCase(initializePlaybackSession.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as VideoError;
      });

    // Sync playback
    builder
      .addCase(syncPlayback.fulfilled, (state, action) => {
        const { videoId, synced_platforms } = action.payload;

        // Update sync status
        if (state.syncStatus[videoId]) {
          state.syncStatus[videoId].last_sync_time = new Date().toISOString();
          state.syncStatus[videoId].synced_platforms = synced_platforms;
        }
      })
      .addCase(syncPlayback.rejected, (state, action) => {
        state.error = action.payload as VideoError;
      });
  },
});

// =============================================================================
// Exports
// =============================================================================

export const {
  setCurrentVideo,
  addVideos,
  updatePlaybackPosition,
  updatePlaybackQuality,
  setPlaybackState,
  updateBufferHealth,
  setUploadProgress,
  clearUploadProgress,
  updateSyncStatus,
  clearError,
  resetVideoState,
} = videoSlice.actions;

export default videoSlice.reducer;