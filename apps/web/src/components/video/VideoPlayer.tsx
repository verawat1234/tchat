import React, { useRef, useEffect, useState, useCallback } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import type { RootState } from '../../store';
import type {
  PlaybackPositionUpdate,
  PlaybackQualityUpdate,
  PlaybackStateUpdate,
  BufferHealthUpdate,
} from '../../types/video';
import {
  updatePlaybackPosition,
  updatePlaybackQuality,
  setPlaybackState,
  updateBufferHealth,
} from '../../store/slices/videoSlice';

export interface VideoPlayerProps {
  videoId: string;
  sessionId?: string;
  autoPlay?: boolean;
  controls?: boolean;
  muted?: boolean;
  loop?: boolean;
  className?: string;
  onEnded?: () => void;
  onError?: (error: Error) => void;
  onPlay?: () => void;
  onPause?: () => void;
  onTimeUpdate?: (currentTime: number) => void;
}

interface VideoQualityOption {
  quality: string;
  resolution: string;
  bitrate: number;
  url: string;
}

export const VideoPlayer: React.FC<VideoPlayerProps> = ({
  videoId,
  sessionId,
  autoPlay = false,
  controls = true,
  muted = false,
  loop = false,
  className = '',
  onEnded,
  onError,
  onPlay,
  onPause,
  onTimeUpdate,
}) => {
  const dispatch = useDispatch();
  const videoRef = useRef<HTMLVideoElement>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [volume, setVolume] = useState(1);
  const [quality, setQuality] = useState('auto');
  const [availableQualities, setAvailableQualities] = useState<VideoQualityOption[]>([]);
  const [showControls, setShowControls] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [bufferProgress, setBufferProgress] = useState(0);

  // Get video data from Redux store
  const videoData = useSelector((state: RootState) =>
    state.video.videos.find(v => v.id === videoId)
  );

  const playbackState = useSelector((state: RootState) =>
    state.video.playbackState[videoId]
  );

  // Load video metadata and stream URL
  useEffect(() => {
    const loadVideo = async () => {
      try {
        const response = await fetch(`/api/v1/videos/${videoId}/stream?quality=${quality}`);
        if (!response.ok) throw new Error('Failed to load video');

        const data = await response.json();

        if (videoRef.current) {
          videoRef.current.src = data.stream_url;
        }

        // Load available qualities
        setAvailableQualities(data.available_qualities?.map((q: string) => ({
          quality: q,
          resolution: q,
          bitrate: 0,
          url: `/api/v1/videos/${videoId}/stream?quality=${q}`,
        })) || []);

      } catch (error) {
        console.error('Error loading video:', error);
        onError?.(error as Error);
      }
    };

    loadVideo();
  }, [videoId, quality, onError]);

  // Sync playback position with Redux
  useEffect(() => {
    if (playbackState?.position !== undefined && videoRef.current) {
      const currentPos = videoRef.current.currentTime;
      const targetPos = playbackState.position;

      // Only sync if difference is significant (>2 seconds)
      if (Math.abs(currentPos - targetPos) > 2) {
        videoRef.current.currentTime = targetPos;
      }
    }
  }, [playbackState?.position]);

  // Handle play
  const handlePlay = useCallback(() => {
    setIsPlaying(true);
    dispatch(setPlaybackState({ videoId, state: 'playing' }));
    onPlay?.();
  }, [videoId, dispatch, onPlay]);

  // Handle pause
  const handlePause = useCallback(() => {
    setIsPlaying(false);
    dispatch(setPlaybackState({ videoId, state: 'paused' }));
    onPause?.();
  }, [videoId, dispatch, onPause]);

  // Handle time update
  const handleTimeUpdate = useCallback(() => {
    if (!videoRef.current) return;

    const time = Math.floor(videoRef.current.currentTime);
    setCurrentTime(time);

    // Update Redux state every second
    if (time !== currentTime) {
      dispatch(updatePlaybackPosition({
        videoId,
        position: time,
        sessionId: sessionId || '',
      }));
      onTimeUpdate?.(time);
    }

    // Update buffer progress
    if (videoRef.current.buffered.length > 0) {
      const bufferedEnd = videoRef.current.buffered.end(videoRef.current.buffered.length - 1);
      const progress = (bufferedEnd / videoRef.current.duration) * 100;
      setBufferProgress(progress);

      dispatch(updateBufferHealth({
        videoId,
        bufferedSeconds: Math.floor(bufferedEnd - time),
        bufferPercentage: progress,
      }));
    }
  }, [videoId, sessionId, currentTime, dispatch, onTimeUpdate]);

  // Handle loaded metadata
  const handleLoadedMetadata = useCallback(() => {
    if (videoRef.current) {
      setDuration(videoRef.current.duration);
    }
  }, []);

  // Handle ended
  const handleEnded = useCallback(() => {
    setIsPlaying(false);
    dispatch(setPlaybackState({ videoId, state: 'ended' }));
    onEnded?.();
  }, [videoId, dispatch, onEnded]);

  // Handle error
  const handleError = useCallback(() => {
    const error = videoRef.current?.error;
    const errorMessage = error?.message || 'Video playback error';
    console.error('Video error:', errorMessage);
    onError?.(new Error(errorMessage));
  }, [onError]);

  // Toggle play/pause
  const togglePlayPause = useCallback(() => {
    if (!videoRef.current) return;

    if (isPlaying) {
      videoRef.current.pause();
    } else {
      videoRef.current.play();
    }
  }, [isPlaying]);

  // Handle seek
  const handleSeek = useCallback((time: number) => {
    if (videoRef.current) {
      videoRef.current.currentTime = time;
      dispatch(updatePlaybackPosition({
        videoId,
        position: Math.floor(time),
        sessionId: sessionId || '',
      }));
    }
  }, [videoId, sessionId, dispatch]);

  // Handle volume change
  const handleVolumeChange = useCallback((newVolume: number) => {
    setVolume(newVolume);
    if (videoRef.current) {
      videoRef.current.volume = newVolume;
    }
  }, []);

  // Handle quality change
  const handleQualityChange = useCallback((newQuality: string) => {
    const currentPos = videoRef.current?.currentTime || 0;
    setQuality(newQuality);
    dispatch(updatePlaybackQuality({ videoId, quality: newQuality }));

    // Restore playback position after quality change
    setTimeout(() => {
      if (videoRef.current) {
        videoRef.current.currentTime = currentPos;
        if (isPlaying) {
          videoRef.current.play();
        }
      }
    }, 100);
  }, [videoId, isPlaying, dispatch]);

  // Toggle fullscreen
  const toggleFullscreen = useCallback(() => {
    if (!videoRef.current?.parentElement) return;

    if (!isFullscreen) {
      if (videoRef.current.parentElement.requestFullscreen) {
        videoRef.current.parentElement.requestFullscreen();
      }
    } else {
      if (document.exitFullscreen) {
        document.exitFullscreen();
      }
    }
    setIsFullscreen(!isFullscreen);
  }, [isFullscreen]);

  // Format time as MM:SS
  const formatTime = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  return (
    <div
      className={`video-player-container relative ${className}`}
      onMouseEnter={() => setShowControls(true)}
      onMouseLeave={() => setShowControls(false)}
    >
      <video
        ref={videoRef}
        className="w-full h-full bg-black"
        autoPlay={autoPlay}
        muted={muted}
        loop={loop}
        playsInline
        onPlay={handlePlay}
        onPause={handlePause}
        onTimeUpdate={handleTimeUpdate}
        onLoadedMetadata={handleLoadedMetadata}
        onEnded={handleEnded}
        onError={handleError}
      />

      {/* Custom Controls */}
      {controls && (
        <div
          className={`absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 to-transparent p-4 transition-opacity ${
            showControls ? 'opacity-100' : 'opacity-0'
          }`}
        >
          {/* Progress Bar */}
          <div className="relative w-full h-1 bg-gray-600 rounded-full mb-4 cursor-pointer group">
            {/* Buffer Progress */}
            <div
              className="absolute h-full bg-gray-400 rounded-full"
              style={{ width: `${bufferProgress}%` }}
            />

            {/* Playback Progress */}
            <div
              className="absolute h-full bg-blue-500 rounded-full"
              style={{ width: `${(currentTime / duration) * 100}%` }}
            />

            {/* Seek Handle */}
            <input
              type="range"
              min="0"
              max={duration || 0}
              value={currentTime}
              onChange={(e) => handleSeek(parseFloat(e.target.value))}
              className="absolute inset-0 w-full opacity-0 cursor-pointer"
            />
          </div>

          <div className="flex items-center justify-between">
            {/* Left Controls */}
            <div className="flex items-center gap-4">
              {/* Play/Pause Button */}
              <button
                onClick={togglePlayPause}
                className="text-white hover:text-blue-400 transition"
                aria-label={isPlaying ? 'Pause' : 'Play'}
              >
                {isPlaying ? (
                  <svg className="w-8 h-8" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z" />
                  </svg>
                ) : (
                  <svg className="w-8 h-8" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M8 5v14l11-7z" />
                  </svg>
                )}
              </button>

              {/* Volume Control */}
              <div className="flex items-center gap-2">
                <button
                  onClick={() => handleVolumeChange(volume > 0 ? 0 : 1)}
                  className="text-white hover:text-blue-400 transition"
                  aria-label={volume > 0 ? 'Mute' : 'Unmute'}
                >
                  {volume > 0 ? (
                    <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02z" />
                    </svg>
                  ) : (
                    <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M16.5 12c0-1.77-1.02-3.29-2.5-4.03v2.21l2.45 2.45c.03-.2.05-.41.05-.63zm2.5 0c0 .94-.2 1.82-.54 2.64l1.51 1.51C20.63 14.91 21 13.5 21 12c0-4.28-2.99-7.86-7-8.77v2.06c2.89.86 5 3.54 5 6.71zM4.27 3L3 4.27 7.73 9H3v6h4l5 5v-6.73l4.25 4.25c-.67.52-1.42.93-2.25 1.18v2.06c1.38-.31 2.63-.95 3.69-1.81L19.73 21 21 19.73l-9-9L4.27 3zM12 4L9.91 6.09 12 8.18V4z" />
                    </svg>
                  )}
                </button>
                <input
                  type="range"
                  min="0"
                  max="1"
                  step="0.1"
                  value={volume}
                  onChange={(e) => handleVolumeChange(parseFloat(e.target.value))}
                  className="w-20"
                />
              </div>

              {/* Time Display */}
              <div className="text-white text-sm">
                {formatTime(currentTime)} / {formatTime(duration)}
              </div>
            </div>

            {/* Right Controls */}
            <div className="flex items-center gap-4">
              {/* Quality Selector */}
              <select
                value={quality}
                onChange={(e) => handleQualityChange(e.target.value)}
                className="bg-black/50 text-white px-2 py-1 rounded border border-gray-600 text-sm"
              >
                <option value="auto">Auto</option>
                {availableQualities.map((q) => (
                  <option key={q.quality} value={q.quality}>
                    {q.quality}
                  </option>
                ))}
              </select>

              {/* Fullscreen Button */}
              <button
                onClick={toggleFullscreen}
                className="text-white hover:text-blue-400 transition"
                aria-label="Fullscreen"
              >
                <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                  {isFullscreen ? (
                    <path d="M5 16h3v3h2v-5H5v2zm3-8H5v2h5V5H8v3zm6 11h2v-3h3v-2h-5v5zm2-11V5h-2v5h5V8h-3z" />
                  ) : (
                    <path d="M7 14H5v5h5v-2H7v-3zm-2-4h2V7h3V5H5v5zm12 7h-3v2h5v-5h-2v3zM14 5v2h3v3h2V5h-5z" />
                  )}
                </svg>
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Loading Indicator */}
      {!videoData && (
        <div className="absolute inset-0 flex items-center justify-center bg-black/50">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
        </div>
      )}
    </div>
  );
};

export default VideoPlayer;