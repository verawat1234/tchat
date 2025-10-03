import React, { useEffect, useState, useCallback } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import type { RootState } from '../../store';
import type { VideoContent, ContentRating } from '../../types/video';
import { fetchVideos, setCurrentVideo } from '../../store/slices/videoSlice';

export interface VideoListProps {
  creatorId?: string;
  category?: string;
  contentRating?: ContentRating;
  limit?: number;
  onVideoSelect?: (video: VideoContent) => void;
  className?: string;
}

interface VideoFilters {
  creator_id?: string;
  category?: string;
  content_rating?: ContentRating;
  sort_by?: 'created_at' | 'view_count' | 'like_count' | 'title';
  sort_order?: 'asc' | 'desc';
}

export const VideoList: React.FC<VideoListProps> = ({
  creatorId,
  category,
  contentRating,
  limit = 20,
  onVideoSelect,
  className = '',
}) => {
  const dispatch = useDispatch();
  const [page, setPage] = useState(1);
  const [sortBy, setSortBy] = useState<'created_at' | 'view_count' | 'like_count' | 'title'>('created_at');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');
  const [isLoadingMore, setIsLoadingMore] = useState(false);

  const videos = useSelector((state: RootState) => state.video.videos);
  const loading = useSelector((state: RootState) => state.video.loading);
  const error = useSelector((state: RootState) => state.video.error);

  // Load videos on mount and when filters change
  useEffect(() => {
    loadVideos();
  }, [creatorId, category, contentRating, sortBy, sortOrder]);

  const loadVideos = useCallback(async () => {
    const filters: VideoFilters = {
      creator_id: creatorId,
      category,
      content_rating: contentRating,
      sort_by: sortBy,
      sort_order: sortOrder,
    };

    try {
      const queryParams = new URLSearchParams();
      queryParams.append('page', String(page));
      queryParams.append('limit', String(limit));

      Object.entries(filters).forEach(([key, value]) => {
        if (value) {
          queryParams.append(key, String(value));
        }
      });

      const response = await fetch(`/api/v1/video?${queryParams.toString()}`);
      if (!response.ok) throw new Error('Failed to load videos');

      const data = await response.json();
      // Dispatch to Redux store
      // dispatch(setVideos(data.videos));
    } catch (error) {
      console.error('Error loading videos:', error);
    }
  }, [creatorId, category, contentRating, sortBy, sortOrder, page, limit]);

  // Load more videos
  const handleLoadMore = useCallback(async () => {
    setIsLoadingMore(true);
    setPage(prev => prev + 1);
    await loadVideos();
    setIsLoadingMore(false);
  }, [loadVideos]);

  // Handle video click
  const handleVideoClick = useCallback((video: VideoContent) => {
    dispatch(setCurrentVideo(video));
    onVideoSelect?.(video);
  }, [dispatch, onVideoSelect]);

  // Format duration from seconds to MM:SS
  const formatDuration = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  // Format view count
  const formatViews = (count: number): string => {
    if (count >= 1000000) {
      return `${(count / 1000000).toFixed(1)}M`;
    }
    if (count >= 1000) {
      return `${(count / 1000).toFixed(1)}K`;
    }
    return String(count);
  };

  // Format upload time
  const formatUploadTime = (dateString: string): string => {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffSec = Math.floor(diffMs / 1000);
    const diffMin = Math.floor(diffSec / 60);
    const diffHour = Math.floor(diffMin / 60);
    const diffDay = Math.floor(diffHour / 24);
    const diffMonth = Math.floor(diffDay / 30);
    const diffYear = Math.floor(diffMonth / 12);

    if (diffYear > 0) return `${diffYear} year${diffYear > 1 ? 's' : ''} ago`;
    if (diffMonth > 0) return `${diffMonth} month${diffMonth > 1 ? 's' : ''} ago`;
    if (diffDay > 0) return `${diffDay} day${diffDay > 1 ? 's' : ''} ago`;
    if (diffHour > 0) return `${diffHour} hour${diffHour > 1 ? 's' : ''} ago`;
    if (diffMin > 0) return `${diffMin} minute${diffMin > 1 ? 's' : ''} ago`;
    return 'Just now';
  };

  if (error) {
    return (
      <div className="p-4 bg-red-50 border border-red-200 rounded text-red-700">
        Error loading videos: {error.message}
      </div>
    );
  }

  return (
    <div className={`video-list-container ${className}`}>
      {/* Sort Controls */}
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold">Videos</h2>
        <div className="flex items-center gap-4">
          <select
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value as any)}
            className="px-3 py-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
          >
            <option value="created_at">Upload Date</option>
            <option value="view_count">Views</option>
            <option value="like_count">Likes</option>
            <option value="title">Title</option>
          </select>
          <select
            value={sortOrder}
            onChange={(e) => setSortOrder(e.target.value as any)}
            className="px-3 py-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
          >
            <option value="desc">Descending</option>
            <option value="asc">Ascending</option>
          </select>
        </div>
      </div>

      {/* Loading State */}
      {loading && videos.length === 0 && (
        <div className="flex items-center justify-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
        </div>
      )}

      {/* Video Grid */}
      {videos.length > 0 && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {videos.map((video) => (
            <div
              key={video.id}
              onClick={() => handleVideoClick(video)}
              className="cursor-pointer group"
            >
              {/* Thumbnail */}
              <div className="relative aspect-video bg-gray-200 rounded-lg overflow-hidden mb-3">
                <img
                  src={video.thumbnail_url || '/default-thumbnail.jpg'}
                  alt={video.title}
                  className="w-full h-full object-cover group-hover:scale-105 transition-transform"
                />
                {/* Duration Badge */}
                <div className="absolute bottom-2 right-2 px-2 py-1 bg-black/80 text-white text-xs font-medium rounded">
                  {formatDuration(video.duration_seconds)}
                </div>
                {/* Status Badge */}
                {video.availability_status !== 'available' && (
                  <div className="absolute top-2 left-2 px-2 py-1 bg-yellow-500 text-white text-xs font-medium rounded">
                    {video.availability_status}
                  </div>
                )}
              </div>

              {/* Video Info */}
              <div className="flex gap-3">
                {/* Creator Avatar (if available) */}
                <div className="flex-shrink-0">
                  <div className="w-9 h-9 bg-gray-300 rounded-full overflow-hidden">
                    {/* Creator avatar would go here */}
                    <div className="w-full h-full bg-gradient-to-br from-blue-400 to-blue-600" />
                  </div>
                </div>

                {/* Video Details */}
                <div className="flex-1 min-w-0">
                  {/* Title */}
                  <h3 className="font-semibold text-sm line-clamp-2 group-hover:text-blue-600">
                    {video.title}
                  </h3>

                  {/* Creator Name (from creator_id - would need to fetch) */}
                  <p className="text-xs text-gray-600 mt-1">
                    Creator ID: {video.creator_id.slice(0, 8)}...
                  </p>

                  {/* Stats */}
                  <div className="flex items-center gap-2 text-xs text-gray-600 mt-1">
                    <span>{formatViews(video.social_metrics.view_count)} views</span>
                    <span>•</span>
                    <span>{formatUploadTime(video.created_at)}</span>
                  </div>

                  {/* Engagement Stats */}
                  <div className="flex items-center gap-3 text-xs text-gray-500 mt-1">
                    <span className="flex items-center gap-1">
                      <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                        <path d="M2 10.5a1.5 1.5 0 113 0v6a1.5 1.5 0 01-3 0v-6zM6 10.333v5.43a2 2 0 001.106 1.79l.05.025A4 4 0 008.943 18h5.416a2 2 0 001.962-1.608l1.2-6A2 2 0 0015.56 8H12V4a2 2 0 00-2-2 1 1 0 00-1 1v.667a4 4 0 01-.8 2.4L6.8 7.933a4 4 0 00-.8 2.4z" />
                      </svg>
                      {formatViews(video.social_metrics.like_count)}
                    </span>
                    <span className="flex items-center gap-1">
                      <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M18 10c0 3.866-3.582 7-8 7a8.841 8.841 0 01-4.083-.98L2 17l1.338-3.123C2.493 12.767 2 11.434 2 10c0-3.866 3.582-7 8-7s8 3.134 8 7zM7 9H5v2h2V9zm8 0h-2v2h2V9zM9 9h2v2H9V9z" clipRule="evenodd" />
                      </svg>
                      {formatViews(video.social_metrics.comment_count)}
                    </span>
                  </div>

                  {/* Tags */}
                  {video.tags.length > 0 && (
                    <div className="flex flex-wrap gap-1 mt-2">
                      {video.tags.slice(0, 2).map((tag, index) => (
                        <span
                          key={index}
                          className="px-2 py-0.5 bg-gray-100 text-gray-600 text-xs rounded"
                        >
                          #{tag}
                        </span>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Empty State */}
      {!loading && videos.length === 0 && (
        <div className="text-center py-12">
          <svg
            className="mx-auto h-12 w-12 text-gray-400"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"
            />
          </svg>
          <p className="mt-4 text-gray-600">No videos found</p>
        </div>
      )}

      {/* Load More Button */}
      {videos.length > 0 && videos.length >= limit && (
        <div className="mt-8 text-center">
          <button
            onClick={handleLoadMore}
            disabled={isLoadingMore}
            className="px-6 py-3 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-300 disabled:cursor-not-allowed"
          >
            {isLoadingMore ? 'Loading...' : 'Load More'}
          </button>
        </div>
      )}
    </div>
  );
};

export default VideoList;