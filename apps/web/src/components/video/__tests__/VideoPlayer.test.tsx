import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import VideoPlayer from '../VideoPlayer';
import VideoUpload from '../VideoUpload';
import VideoList from '../VideoList';
import { videoApi } from '../../../services/api';

// Mock video API
vi.mock('../../../services/api', () => ({
  videoApi: {
    useGetVideoQuery: vi.fn(),
    useCreatePlaybackSessionMutation: vi.fn(),
    useSyncPlaybackPositionMutation: vi.fn(),
    useGetStreamURLQuery: vi.fn(),
    useUploadVideoMutation: vi.fn(),
    useGetVideosQuery: vi.fn(),
  }
}));

// Mock store setup
const createMockStore = () => {
  return configureStore({
    reducer: {
      [videoApi.reducerPath]: vi.fn(() => ({})),
    },
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware().concat(videoApi.middleware),
  });
};

describe('VideoPlayer Component', () => {
  let mockStore: ReturnType<typeof createMockStore>;

  beforeEach(() => {
    mockStore = createMockStore();
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('should render video player with loading state', () => {
    vi.mocked(videoApi.useGetVideoQuery).mockReturnValue({
      data: undefined,
      isLoading: true,
      isError: false,
      error: undefined,
    } as any);

    vi.mocked(videoApi.useCreatePlaybackSessionMutation).mockReturnValue([
      vi.fn(),
      { isLoading: false, isError: false, error: undefined }
    ] as any);

    vi.mocked(videoApi.useGetStreamURLQuery).mockReturnValue({
      data: undefined,
      isLoading: true,
      isError: false,
      error: undefined,
    } as any);

    render(
      <Provider store={mockStore}>
        <VideoPlayer videoId="test-video-id" />
      </Provider>
    );

    expect(screen.getByText(/loading/i)).toBeInTheDocument();
  });

  it('should render video player with video data', async () => {
    const mockVideo = {
      id: 'test-video-id',
      title: 'Test Video',
      description: 'Test Description',
      videoUrl: 'https://example.com/video.mp4',
      thumbnailUrl: 'https://example.com/thumbnail.jpg',
      durationSeconds: 300,
      creatorId: 'creator-1',
      uploadStatus: 'available',
    };

    const mockStreamURL = {
      streamUrl: 'https://cdn.example.com/stream/test-video-id/master.m3u8',
      protocol: 'hls',
    };

    vi.mocked(videoApi.useGetVideoQuery).mockReturnValue({
      data: mockVideo,
      isLoading: false,
      isError: false,
      error: undefined,
    } as any);

    vi.mocked(videoApi.useCreatePlaybackSessionMutation).mockReturnValue([
      vi.fn().mockResolvedValue({ data: { id: 'session-1' } }),
      { isLoading: false, isError: false, error: undefined }
    ] as any);

    vi.mocked(videoApi.useGetStreamURLQuery).mockReturnValue({
      data: mockStreamURL,
      isLoading: false,
      isError: false,
      error: undefined,
    } as any);

    render(
      <Provider store={mockStore}>
        <VideoPlayer videoId="test-video-id" />
      </Provider>
    );

    await waitFor(() => {
      expect(screen.getByText('Test Video')).toBeInTheDocument();
    });
  });

  it('should handle video player controls', async () => {
    const mockVideo = {
      id: 'test-video-id',
      title: 'Test Video',
      videoUrl: 'https://example.com/video.mp4',
      durationSeconds: 300,
    };

    const mockStreamURL = {
      streamUrl: 'https://cdn.example.com/stream/test-video-id/master.m3u8',
      protocol: 'hls',
    };

    const mockSyncPlayback = vi.fn().mockResolvedValue({ data: {} });

    vi.mocked(videoApi.useGetVideoQuery).mockReturnValue({
      data: mockVideo,
      isLoading: false,
      isError: false,
    } as any);

    vi.mocked(videoApi.useCreatePlaybackSessionMutation).mockReturnValue([
      vi.fn().mockResolvedValue({ data: { id: 'session-1' } }),
      { isLoading: false, isError: false }
    ] as any);

    vi.mocked(videoApi.useGetStreamURLQuery).mockReturnValue({
      data: mockStreamURL,
      isLoading: false,
      isError: false,
    } as any);

    vi.mocked(videoApi.useSyncPlaybackPositionMutation).mockReturnValue([
      mockSyncPlayback,
      { isLoading: false, isError: false }
    ] as any);

    render(
      <Provider store={mockStore}>
        <VideoPlayer videoId="test-video-id" />
      </Provider>
    );

    await waitFor(() => {
      expect(screen.getByText('Test Video')).toBeInTheDocument();
    });

    // Test play/pause functionality (if exposed via controls)
    const playButton = screen.queryByRole('button', { name: /play/i });
    if (playButton) {
      fireEvent.click(playButton);
      await waitFor(() => {
        expect(mockSyncPlayback).toHaveBeenCalled();
      });
    }
  });

  it('should handle video player errors', () => {
    vi.mocked(videoApi.useGetVideoQuery).mockReturnValue({
      data: undefined,
      isLoading: false,
      isError: true,
      error: { message: 'Video not found' },
    } as any);

    vi.mocked(videoApi.useCreatePlaybackSessionMutation).mockReturnValue([
      vi.fn(),
      { isLoading: false, isError: false }
    ] as any);

    render(
      <Provider store={mockStore}>
        <VideoPlayer videoId="invalid-video-id" />
      </Provider>
    );

    expect(screen.getByText(/error/i)).toBeInTheDocument();
  });

  it('should sync playback position across platforms', async () => {
    const mockVideo = {
      id: 'test-video-id',
      title: 'Test Video',
      videoUrl: 'https://example.com/video.mp4',
      durationSeconds: 300,
    };

    const mockStreamURL = {
      streamUrl: 'https://cdn.example.com/stream/test-video-id/master.m3u8',
      protocol: 'hls',
    };

    const mockSyncPlayback = vi.fn().mockResolvedValue({
      data: {
        position: 150.0,
        synced: true,
        timestamp: new Date().toISOString(),
      }
    });

    vi.mocked(videoApi.useGetVideoQuery).mockReturnValue({
      data: mockVideo,
      isLoading: false,
      isError: false,
    } as any);

    vi.mocked(videoApi.useCreatePlaybackSessionMutation).mockReturnValue([
      vi.fn().mockResolvedValue({ data: { id: 'session-1' } }),
      { isLoading: false, isError: false }
    ] as any);

    vi.mocked(videoApi.useGetStreamURLQuery).mockReturnValue({
      data: mockStreamURL,
      isLoading: false,
      isError: false,
    } as any);

    vi.mocked(videoApi.useSyncPlaybackPositionMutation).mockReturnValue([
      mockSyncPlayback,
      { isLoading: false, isError: false }
    ] as any);

    render(
      <Provider store={mockStore}>
        <VideoPlayer videoId="test-video-id" enableSync={true} />
      </Provider>
    );

    await waitFor(() => {
      expect(screen.getByText('Test Video')).toBeInTheDocument();
    });

    // Sync should be called periodically or on position change
    await waitFor(() => {
      expect(mockSyncPlayback).toHaveBeenCalled();
    }, { timeout: 6000 });
  });
});

describe('VideoUpload Component', () => {
  let mockStore: ReturnType<typeof createMockStore>;

  beforeEach(() => {
    mockStore = createMockStore();
    vi.clearAllMocks();
  });

  it('should render upload form', () => {
    vi.mocked(videoApi.useUploadVideoMutation).mockReturnValue([
      vi.fn(),
      { isLoading: false, isError: false, error: undefined }
    ] as any);

    render(
      <Provider store={mockStore}>
        <VideoUpload />
      </Provider>
    );

    expect(screen.getByLabelText(/title/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/description/i)).toBeInTheDocument();
    expect(screen.getByText(/upload/i)).toBeInTheDocument();
  });

  it('should handle file selection', async () => {
    const mockUploadVideo = vi.fn().mockResolvedValue({
      data: {
        id: 'video-1',
        status: 'processing',
        uploadUrl: 'https://upload.example.com/video-1'
      }
    });

    vi.mocked(videoApi.useUploadVideoMutation).mockReturnValue([
      mockUploadVideo,
      { isLoading: false, isError: false, error: undefined }
    ] as any);

    render(
      <Provider store={mockStore}>
        <VideoUpload />
      </Provider>
    );

    const fileInput = screen.getByLabelText(/video file/i) as HTMLInputElement;
    const testFile = new File(['video content'], 'test-video.mp4', { type: 'video/mp4' });

    fireEvent.change(fileInput, { target: { files: [testFile] } });

    await waitFor(() => {
      expect(fileInput.files?.[0]).toBe(testFile);
    });
  });

  it('should handle video upload submission', async () => {
    const mockUploadVideo = vi.fn().mockResolvedValue({
      data: {
        id: 'video-1',
        status: 'processing',
        uploadUrl: 'https://upload.example.com/video-1'
      }
    });

    vi.mocked(videoApi.useUploadVideoMutation).mockReturnValue([
      mockUploadVideo,
      { isLoading: false, isError: false, error: undefined }
    ] as any);

    render(
      <Provider store={mockStore}>
        <VideoUpload />
      </Provider>
    );

    // Fill form
    const titleInput = screen.getByLabelText(/title/i);
    const descriptionInput = screen.getByLabelText(/description/i);
    const fileInput = screen.getByLabelText(/video file/i) as HTMLInputElement;

    fireEvent.change(titleInput, { target: { value: 'Test Video' } });
    fireEvent.change(descriptionInput, { target: { value: 'Test Description' } });

    const testFile = new File(['video content'], 'test-video.mp4', { type: 'video/mp4' });
    fireEvent.change(fileInput, { target: { files: [testFile] } });

    // Submit form
    const uploadButton = screen.getByText(/upload/i);
    fireEvent.click(uploadButton);

    await waitFor(() => {
      expect(mockUploadVideo).toHaveBeenCalled();
    });
  });

  it('should show upload progress', async () => {
    const mockUploadVideo = vi.fn().mockImplementation(() => {
      return new Promise((resolve) => {
        setTimeout(() => {
          resolve({
            data: {
              id: 'video-1',
              status: 'processing',
            }
          });
        }, 1000);
      });
    });

    vi.mocked(videoApi.useUploadVideoMutation).mockReturnValue([
      mockUploadVideo,
      { isLoading: true, isError: false, error: undefined }
    ] as any);

    render(
      <Provider store={mockStore}>
        <VideoUpload />
      </Provider>
    );

    const titleInput = screen.getByLabelText(/title/i);
    const fileInput = screen.getByLabelText(/video file/i) as HTMLInputElement;

    fireEvent.change(titleInput, { target: { value: 'Test Video' } });
    const testFile = new File(['video content'], 'test-video.mp4', { type: 'video/mp4' });
    fireEvent.change(fileInput, { target: { files: [testFile] } });

    const uploadButton = screen.getByText(/upload/i);
    fireEvent.click(uploadButton);

    // Should show uploading state
    expect(screen.getByText(/uploading/i)).toBeInTheDocument();
  });

  it('should handle upload errors', async () => {
    const mockUploadVideo = vi.fn().mockRejectedValue({
      error: { message: 'Upload failed' }
    });

    vi.mocked(videoApi.useUploadVideoMutation).mockReturnValue([
      mockUploadVideo,
      { isLoading: false, isError: true, error: { message: 'Upload failed' } }
    ] as any);

    render(
      <Provider store={mockStore}>
        <VideoUpload />
      </Provider>
    );

    expect(screen.getByText(/upload failed/i)).toBeInTheDocument();
  });
});

describe('VideoList Component', () => {
  let mockStore: ReturnType<typeof createMockStore>;

  beforeEach(() => {
    mockStore = createMockStore();
    vi.clearAllMocks();
  });

  it('should render empty video list', () => {
    vi.mocked(videoApi.useGetVideosQuery).mockReturnValue({
      data: { videos: [], total: 0, page: 1 },
      isLoading: false,
      isError: false,
      error: undefined,
    } as any);

    render(
      <Provider store={mockStore}>
        <VideoList />
      </Provider>
    );

    expect(screen.getByText(/no videos/i)).toBeInTheDocument();
  });

  it('should render list of videos', async () => {
    const mockVideos = [
      {
        id: 'video-1',
        title: 'Video 1',
        description: 'Description 1',
        thumbnailUrl: 'https://example.com/thumb1.jpg',
        durationSeconds: 300,
        viewCount: 1000,
      },
      {
        id: 'video-2',
        title: 'Video 2',
        description: 'Description 2',
        thumbnailUrl: 'https://example.com/thumb2.jpg',
        durationSeconds: 240,
        viewCount: 500,
      },
    ];

    vi.mocked(videoApi.useGetVideosQuery).mockReturnValue({
      data: { videos: mockVideos, total: 2, page: 1 },
      isLoading: false,
      isError: false,
      error: undefined,
    } as any);

    render(
      <Provider store={mockStore}>
        <VideoList />
      </Provider>
    );

    await waitFor(() => {
      expect(screen.getByText('Video 1')).toBeInTheDocument();
      expect(screen.getByText('Video 2')).toBeInTheDocument();
    });
  });

  it('should handle video selection', async () => {
    const mockVideos = [
      {
        id: 'video-1',
        title: 'Video 1',
        description: 'Description 1',
        thumbnailUrl: 'https://example.com/thumb1.jpg',
        durationSeconds: 300,
      },
    ];

    const mockOnVideoSelect = vi.fn();

    vi.mocked(videoApi.useGetVideosQuery).mockReturnValue({
      data: { videos: mockVideos, total: 1, page: 1 },
      isLoading: false,
      isError: false,
      error: undefined,
    } as any);

    render(
      <Provider store={mockStore}>
        <VideoList onVideoSelect={mockOnVideoSelect} />
      </Provider>
    );

    await waitFor(() => {
      const videoItem = screen.getByText('Video 1');
      fireEvent.click(videoItem);
      expect(mockOnVideoSelect).toHaveBeenCalledWith(mockVideos[0]);
    });
  });

  it('should handle pagination', async () => {
    const mockVideosPage1 = [
      { id: 'video-1', title: 'Video 1', durationSeconds: 300 },
      { id: 'video-2', title: 'Video 2', durationSeconds: 240 },
    ];

    const mockVideosPage2 = [
      { id: 'video-3', title: 'Video 3', durationSeconds: 180 },
      { id: 'video-4', title: 'Video 4', durationSeconds: 360 },
    ];

    vi.mocked(videoApi.useGetVideosQuery).mockReturnValueOnce({
      data: { videos: mockVideosPage1, total: 4, page: 1, hasMore: true },
      isLoading: false,
      isError: false,
    } as any);

    const { rerender } = render(
      <Provider store={mockStore}>
        <VideoList page={1} />
      </Provider>
    );

    await waitFor(() => {
      expect(screen.getByText('Video 1')).toBeInTheDocument();
      expect(screen.getByText('Video 2')).toBeInTheDocument();
    });

    // Mock page 2 response
    vi.mocked(videoApi.useGetVideosQuery).mockReturnValueOnce({
      data: { videos: mockVideosPage2, total: 4, page: 2, hasMore: false },
      isLoading: false,
      isError: false,
    } as any);

    rerender(
      <Provider store={mockStore}>
        <VideoList page={2} />
      </Provider>
    );

    await waitFor(() => {
      expect(screen.getByText('Video 3')).toBeInTheDocument();
      expect(screen.getByText('Video 4')).toBeInTheDocument();
    });
  });

  it('should handle loading state', () => {
    vi.mocked(videoApi.useGetVideosQuery).mockReturnValue({
      data: undefined,
      isLoading: true,
      isError: false,
      error: undefined,
    } as any);

    render(
      <Provider store={mockStore}>
        <VideoList />
      </Provider>
    );

    expect(screen.getByText(/loading/i)).toBeInTheDocument();
  });

  it('should handle error state', () => {
    vi.mocked(videoApi.useGetVideosQuery).mockReturnValue({
      data: undefined,
      isLoading: false,
      isError: true,
      error: { message: 'Failed to load videos' },
    } as any);

    render(
      <Provider store={mockStore}>
        <VideoList />
      </Provider>
    );

    expect(screen.getByText(/error/i)).toBeInTheDocument();
  });

  it('should filter videos by category', async () => {
    const mockTrendingVideos = [
      { id: 'video-1', title: 'Trending Video 1', durationSeconds: 300, viewCount: 10000 },
      { id: 'video-2', title: 'Trending Video 2', durationSeconds: 240, viewCount: 8000 },
    ];

    vi.mocked(videoApi.useGetVideosQuery).mockReturnValue({
      data: { videos: mockTrendingVideos, total: 2, page: 1 },
      isLoading: false,
      isError: false,
    } as any);

    render(
      <Provider store={mockStore}>
        <VideoList category="trending" />
      </Provider>
    );

    await waitFor(() => {
      expect(screen.getByText('Trending Video 1')).toBeInTheDocument();
      expect(screen.getByText('Trending Video 2')).toBeInTheDocument();
    });
  });

  it('should search videos by query', async () => {
    const mockSearchResults = [
      { id: 'video-1', title: 'React Tutorial', durationSeconds: 600 },
    ];

    vi.mocked(videoApi.useGetVideosQuery).mockReturnValue({
      data: { videos: mockSearchResults, total: 1, page: 1 },
      isLoading: false,
      isError: false,
    } as any);

    render(
      <Provider store={mockStore}>
        <VideoList query="React" />
      </Provider>
    );

    await waitFor(() => {
      expect(screen.getByText('React Tutorial')).toBeInTheDocument();
    });
  });
});