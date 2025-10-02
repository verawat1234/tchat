// apps/web/src/hooks/useSecureVideo.ts
// React hook for secure video playback with blob URLs
// Automatically handles cleanup and memory management

import { useEffect, useState, useCallback } from 'react';
import { secureVideoService } from '../services/secureVideoService';

interface UseSecureVideoOptions {
  videoId: string;
  quality?: string;
  autoLoad?: boolean;
}

interface UseSecureVideoResult {
  blobUrl: string | null;
  isLoading: boolean;
  error: Error | null;
  load: () => Promise<void>;
  revoke: () => void;
}

/**
 * useSecureVideo - Hook for secure video playback
 *
 * Features:
 * - Automatic blob URL creation and cleanup
 * - Memory leak prevention with automatic revocation
 * - Loading and error states
 * - Prevents video downloads
 *
 * @example
 * ```tsx
 * const VideoPlayer = ({ videoId }) => {
 *   const { blobUrl, isLoading, error } = useSecureVideo({
 *     videoId,
 *     quality: '720p',
 *     autoLoad: true
 *   });
 *
 *   if (isLoading) return <div>Loading...</div>;
 *   if (error) return <div>Error: {error.message}</div>;
 *
 *   return (
 *     <video src={blobUrl} controls />
 *   );
 * };
 * ```
 */
export function useSecureVideo({
  videoId,
  quality = 'auto',
  autoLoad = true,
}: UseSecureVideoOptions): UseSecureVideoResult {
  const [blobUrl, setBlobUrl] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  // Load video blob URL
  const load = useCallback(async () => {
    if (!videoId) return;

    setIsLoading(true);
    setError(null);

    try {
      const url = await secureVideoService.createSecureBlobURL(videoId, quality);
      setBlobUrl(url);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Failed to load video'));
    } finally {
      setIsLoading(false);
    }
  }, [videoId, quality]);

  // Revoke blob URL
  const revoke = useCallback(() => {
    if (videoId && quality) {
      const cacheKey = `${videoId}-${quality}`;
      secureVideoService.revokeBlobURL(cacheKey);
      setBlobUrl(null);
    }
  }, [videoId, quality]);

  // Auto-load on mount if enabled
  useEffect(() => {
    if (autoLoad) {
      load();
    }
  }, [autoLoad, load]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      revoke();
    };
  }, [revoke]);

  return {
    blobUrl,
    isLoading,
    error,
    load,
    revoke,
  };
}

/**
 * useSecureVideoWithMediaSource - Advanced hook using MediaSource API
 * For streaming large videos in chunks without downloading entire file
 */
export function useSecureVideoWithMediaSource({
  videoId,
  quality = 'auto',
}: UseSecureVideoOptions) {
  const [sourceUrl, setSourceUrl] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    let mediaSource: MediaSource | null = null;
    let sourceBuffer: SourceBuffer | null = null;

    const loadVideo = async () => {
      setIsLoading(true);
      setError(null);

      try {
        // Get streaming token
        const tokenData = await secureVideoService.getStreamToken(videoId, quality);

        // Create MediaSource
        mediaSource = new MediaSource();
        const url = URL.createObjectURL(mediaSource);
        setSourceUrl(url);

        // Wait for MediaSource to open
        await new Promise<void>((resolve) => {
          mediaSource!.addEventListener('sourceopen', () => resolve(), { once: true });
        });

        // Add source buffer
        sourceBuffer = mediaSource!.addSourceBuffer('video/mp4; codecs="avc1.42E01E, mp4a.40.2"');

        // Fetch video in chunks
        const response = await fetch(tokenData.signed_url, {
          credentials: 'include',
        });

        if (!response.body) {
          throw new Error('No response body');
        }

        const reader = response.body.getReader();

        // Read and append chunks
        while (true) {
          const { done, value } = await reader.read();

          if (done) {
            mediaSource!.endOfStream();
            break;
          }

          // Wait for buffer to be ready
          await new Promise<void>((resolve) => {
            if (!sourceBuffer!.updating) {
              resolve();
            } else {
              sourceBuffer!.addEventListener('updateend', () => resolve(), { once: true });
            }
          });

          // Append chunk
          sourceBuffer!.appendBuffer(value);
        }

        setIsLoading(false);
      } catch (err) {
        setError(err instanceof Error ? err : new Error('Failed to load video'));
        setIsLoading(false);
      }
    };

    loadVideo();

    // Cleanup
    return () => {
      if (sourceUrl) {
        URL.revokeObjectURL(sourceUrl);
      }
      if (mediaSource) {
        try {
          mediaSource.endOfStream();
        } catch (e) {
          // MediaSource may already be closed
        }
      }
    };
  }, [videoId, quality]);

  return {
    sourceUrl,
    isLoading,
    error,
  };
}