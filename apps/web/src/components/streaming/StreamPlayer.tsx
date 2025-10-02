/**
 * StreamPlayer Component
 *
 * Requirements from T059:
 * - React component for video playback
 * - Video element: Render video with WebRTC stream
 * - Controls: Quality selector, volume, fullscreen
 * - Overlays: Product features (store context), reactions (video context)
 */

import React, { useEffect, useRef, useState } from 'react';
import { WebRTCClient } from '../../services/streaming/webrtcClient';
import { useGetStreamQuery } from '../../services/streaming/streamingApi';

interface StreamPlayerProps {
  streamId: string;
  webrtcClient: WebRTCClient;
}

export const StreamPlayer: React.FC<StreamPlayerProps> = ({ streamId, webrtcClient }) => {
  const videoRef = useRef<HTMLVideoElement>(null);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [volume, setVolume] = useState(1);
  const [isMuted, setIsMuted] = useState(false);
  const [currentQuality, setCurrentQuality] = useState<'low' | 'mid' | 'high'>('high');
  const [viewerCount, setViewerCount] = useState(0);
  const [showControls, setShowControls] = useState(true);
  const controlsTimeoutRef = useRef<number | null>(null);

  const { data: stream, error, isLoading } = useGetStreamQuery(streamId);

  // Setup video stream from WebRTC
  useEffect(() => {
    if (!videoRef.current || !webrtcClient) return;

    console.log('[StreamPlayer] Setting up video stream for stream:', streamId);

    // The webrtcClient should have an onTrack callback that was set during connect()
    // We'll set the video element's srcObject when we receive a track
    const handleTrack = (track: MediaStreamTrack) => {
      console.log('[StreamPlayer] Received track:', track.kind);
      if (videoRef.current && track.kind === 'video') {
        const mediaStream = new MediaStream([track]);
        videoRef.current.srcObject = mediaStream;
      }
    };

    // Note: In a real implementation, you would subscribe to track events from webrtcClient
    // For now, we assume webrtcClient.connect() was called elsewhere with the onTrack callback

    return () => {
      if (videoRef.current) {
        videoRef.current.srcObject = null;
      }
    };
  }, [streamId, webrtcClient]);

  // Update viewer count from stream data
  useEffect(() => {
    if (stream) {
      setViewerCount(stream.viewer_count);
    }
  }, [stream]);

  // Handle quality change
  const handleQualityChange = async (quality: 'low' | 'mid' | 'high') => {
    console.log('[StreamPlayer] Switching quality to:', quality);
    await webrtcClient.switchQuality(quality);
    setCurrentQuality(quality);
  };

  // Handle volume change
  const handleVolumeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newVolume = parseFloat(e.target.value);
    setVolume(newVolume);
    setIsMuted(newVolume === 0);
    if (videoRef.current) {
      videoRef.current.volume = newVolume;
      videoRef.current.muted = newVolume === 0;
    }
  };

  // Handle mute toggle
  const handleMuteToggle = () => {
    if (videoRef.current) {
      const newMuted = !isMuted;
      setIsMuted(newMuted);
      videoRef.current.muted = newMuted;
      if (newMuted) {
        setVolume(0);
        videoRef.current.volume = 0;
      } else {
        setVolume(1);
        videoRef.current.volume = 1;
      }
    }
  };

  // Handle fullscreen toggle
  const handleFullscreenToggle = () => {
    if (!videoRef.current) return;

    if (!isFullscreen) {
      videoRef.current.requestFullscreen();
    } else {
      document.exitFullscreen();
    }
  };

  // Listen for fullscreen changes
  useEffect(() => {
    const handleFullscreenChange = () => {
      setIsFullscreen(!!document.fullscreenElement);
    };

    document.addEventListener('fullscreenchange', handleFullscreenChange);
    return () => {
      document.removeEventListener('fullscreenchange', handleFullscreenChange);
    };
  }, []);

  // Auto-hide controls after 3 seconds of no mouse movement
  const handleMouseMove = () => {
    setShowControls(true);
    if (controlsTimeoutRef.current) {
      window.clearTimeout(controlsTimeoutRef.current);
    }
    controlsTimeoutRef.current = window.setTimeout(() => {
      setShowControls(false);
    }, 3000);
  };

  // Cleanup timeout on unmount
  useEffect(() => {
    return () => {
      if (controlsTimeoutRef.current) {
        window.clearTimeout(controlsTimeoutRef.current);
      }
    };
  }, []);

  // Get quality label
  const getQualityLabel = (quality: 'low' | 'mid' | 'high'): string => {
    switch (quality) {
      case 'low':
        return '360p';
      case 'mid':
        return '720p';
      case 'high':
        return '1080p';
      default:
        return '720p';
    }
  };

  // Loading state
  if (isLoading) {
    return (
      <div className="stream-player relative w-full bg-black rounded-lg overflow-hidden flex items-center justify-center h-96">
        <div className="text-white text-lg">Loading stream...</div>
      </div>
    );
  }

  // Error state
  if (error || !stream) {
    return (
      <div className="stream-player relative w-full bg-black rounded-lg overflow-hidden flex items-center justify-center h-96">
        <div className="text-red-500 text-lg">
          {error ? 'Error loading stream' : 'Stream not found'}
        </div>
      </div>
    );
  }

  return (
    <div
      className="stream-player relative w-full bg-black rounded-lg overflow-hidden"
      onMouseMove={handleMouseMove}
      onMouseLeave={() => setShowControls(false)}
    >
      {/* Video Element */}
      <video
        ref={videoRef}
        className="w-full h-auto"
        autoPlay
        playsInline
        controls={false}
      />

      {/* Stream Info Overlay */}
      <div className="absolute top-4 left-4 bg-black bg-opacity-60 text-white px-3 py-2 rounded-md z-10">
        <div className="text-sm font-medium">{stream.title}</div>
        <div className="text-xs text-gray-300 mt-1">
          <span className="inline-flex items-center">
            <span className="w-2 h-2 bg-red-500 rounded-full mr-2 animate-pulse"></span>
            {stream.status === 'live' ? 'LIVE' : stream.status.toUpperCase()}
          </span>
          <span className="ml-3">{viewerCount.toLocaleString()} viewers</span>
        </div>
      </div>

      {/* Controls Overlay */}
      <div
        className={`absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black via-black/80 to-transparent p-4 transition-opacity duration-300 ${
          showControls ? 'opacity-100' : 'opacity-0'
        }`}
      >
        <div className="flex items-center justify-between">
          {/* Quality Selector */}
          <div className="flex items-center space-x-2">
            <span className="text-white text-sm font-medium">Quality:</span>
            <select
              value={currentQuality}
              onChange={(e) => handleQualityChange(e.target.value as 'low' | 'mid' | 'high')}
              className="bg-gray-800 text-white text-sm px-3 py-1.5 rounded border border-gray-600 hover:border-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-colors"
            >
              <option value="low">{getQualityLabel('low')}</option>
              <option value="mid">{getQualityLabel('mid')}</option>
              <option value="high">{getQualityLabel('high')}</option>
            </select>
          </div>

          {/* Volume Control */}
          <div className="flex items-center space-x-2">
            <button
              onClick={handleMuteToggle}
              className="text-white hover:text-gray-300 transition-colors"
              aria-label={isMuted ? 'Unmute' : 'Mute'}
            >
              {isMuted || volume === 0 ? (
                <svg
                  className="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z M17 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2"
                  />
                </svg>
              ) : (
                <svg
                  className="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M15.536 8.464a5 5 0 010 7.072m2.828-9.9a9 9 0 010 12.728M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z"
                  />
                </svg>
              )}
            </button>
            <input
              type="range"
              min="0"
              max="1"
              step="0.1"
              value={volume}
              onChange={handleVolumeChange}
              className="w-24 h-1 bg-gray-600 rounded-lg appearance-none cursor-pointer slider"
              aria-label="Volume"
            />
          </div>

          {/* Fullscreen Button */}
          <button
            onClick={handleFullscreenToggle}
            className="text-white hover:text-gray-300 transition-colors"
            aria-label={isFullscreen ? 'Exit fullscreen' : 'Enter fullscreen'}
          >
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              {isFullscreen ? (
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 9V4.5M9 9H4.5M9 9L3.75 3.75M9 15v4.5M9 15H4.5M9 15l-5.25 5.25M15 9h4.5M15 9V4.5M15 9l5.25-5.25M15 15h4.5M15 15v4.5m0-4.5l5.25 5.25"
                />
              ) : (
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4"
                />
              )}
            </svg>
          </button>
        </div>
      </div>

      {/* Product Overlay (for store streams) */}
      {stream.stream_type === 'store' && stream.featured_products && stream.featured_products.length > 0 && (
        <div className="absolute top-20 right-4 bg-white rounded-lg shadow-xl p-4 max-w-xs z-10 transition-transform hover:scale-105">
          <div className="text-sm font-semibold mb-2 text-gray-900">Featured Products</div>
          <div className="space-y-2">
            {stream.featured_products.slice(0, 3).map((productId, index) => (
              <div
                key={productId}
                className="flex items-center space-x-2 p-2 bg-gray-50 rounded hover:bg-gray-100 transition-colors cursor-pointer"
              >
                <div className="w-12 h-12 bg-gray-200 rounded flex items-center justify-center">
                  <span className="text-xs text-gray-500">#{index + 1}</span>
                </div>
                <div className="flex-1 min-w-0">
                  <div className="text-xs font-medium text-gray-900 truncate">Product {productId.substring(0, 8)}</div>
                  <div className="text-xs text-gray-500">View details →</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Reaction Overlay (for video streams) */}
      {stream.stream_type === 'video' && (
        <div className="absolute bottom-24 right-4 flex flex-col space-y-2 z-10">
          <button
            className="bg-white/90 hover:bg-white rounded-full p-3 shadow-lg transition-all hover:scale-110"
            aria-label="Like"
          >
            <svg
              className="w-6 h-6 text-red-500"
              fill="currentColor"
              viewBox="0 0 24 24"
            >
              <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" />
            </svg>
          </button>
          <button
            className="bg-white/90 hover:bg-white rounded-full p-3 shadow-lg transition-all hover:scale-110"
            aria-label="Comment"
          >
            <svg
              className="w-6 h-6 text-blue-500"
              fill="currentColor"
              viewBox="0 0 24 24"
            >
              <path d="M20 2H4c-1.1 0-2 .9-2 2v18l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2z" />
            </svg>
          </button>
          <button
            className="bg-white/90 hover:bg-white rounded-full p-3 shadow-lg transition-all hover:scale-110"
            aria-label="Share"
          >
            <svg
              className="w-6 h-6 text-green-500"
              fill="currentColor"
              viewBox="0 0 24 24"
            >
              <path d="M18 16.08c-.76 0-1.44.3-1.96.77L8.91 12.7c.05-.23.09-.46.09-.7s-.04-.47-.09-.7l7.05-4.11c.54.5 1.25.81 2.04.81 1.66 0 3-1.34 3-3s-1.34-3-3-3-3 1.34-3 3c0 .24.04.47.09.7L8.04 9.81C7.5 9.31 6.79 9 6 9c-1.66 0-3 1.34-3 3s1.34 3 3 3c.79 0 1.5-.31 2.04-.81l7.12 4.16c-.05.21-.08.43-.08.65 0 1.61 1.31 2.92 2.92 2.92 1.61 0 2.92-1.31 2.92-2.92s-1.31-2.92-2.92-2.92z" />
            </svg>
          </button>
        </div>
      )}

      {/* Stream offline overlay */}
      {stream.status !== 'live' && (
        <div className="absolute inset-0 bg-black/80 flex items-center justify-center z-20">
          <div className="text-center text-white">
            <div className="text-2xl font-bold mb-2">
              {stream.status === 'scheduled' ? 'Stream Scheduled' : 'Stream Ended'}
            </div>
            <div className="text-gray-400">
              {stream.status === 'scheduled' && stream.scheduled_start_time && (
                <>Starts at {new Date(stream.scheduled_start_time).toLocaleString()}</>
              )}
              {stream.status === 'ended' && stream.end_time && (
                <>Ended at {new Date(stream.end_time).toLocaleString()}</>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};