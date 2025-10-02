// apps/web/src/services/secureVideoService.ts
// Secure video service that prevents downloads using blob URLs
// Implements token-based authentication and memory management

import { StreamToken } from '../types/video';

/**
 * SecureVideoService - Handles secure video streaming with blob URLs
 *
 * Features:
 * - Token-based authentication
 * - Blob URL generation from authenticated fetch
 * - Automatic blob URL revocation to prevent memory leaks
 * - Prevents direct video downloads
 */
export class SecureVideoService {
  private blobCache = new Map<string, string>();
  private revokeTimers = new Map<string, NodeJS.Timeout>();

  /**
   * Get streaming token from backend
   */
  async getStreamToken(videoId: string, quality: string = 'auto'): Promise<StreamToken> {
    const response = await fetch(`/api/v1/videos/${videoId}/token?quality=${quality}`, {
      method: 'GET',
      credentials: 'include',
      headers: {
        'Authorization': `Bearer ${this.getAuthToken()}`,
      },
    });

    if (!response.ok) {
      throw new Error('Failed to get streaming token');
    }

    return response.json();
  }

  /**
   * Fetch video as blob using authenticated request
   */
  async fetchVideoBlob(signedUrl: string): Promise<Blob> {
    const response = await fetch(signedUrl, {
      method: 'GET',
      credentials: 'include',
      headers: {
        'Accept': 'video/*',
      },
    });

    if (!response.ok) {
      throw new Error('Failed to fetch video content');
    }

    return response.blob();
  }

  /**
   * Create blob URL for secure video playback
   * This prevents direct video downloads
   */
  async createSecureBlobURL(videoId: string, quality: string = 'auto'): Promise<string> {
    // Check cache first
    const cacheKey = `${videoId}-${quality}`;
    if (this.blobCache.has(cacheKey)) {
      return this.blobCache.get(cacheKey)!;
    }

    try {
      // Get streaming token
      const tokenData = await this.getStreamToken(videoId, quality);

      // Fetch video as blob using signed URL
      const videoBlob = await this.fetchVideoBlob(tokenData.signed_url);

      // Create blob URL
      const blobUrl = URL.createObjectURL(videoBlob);

      // Cache blob URL
      this.blobCache.set(cacheKey, blobUrl);

      // Auto-revoke after token expires
      const expiresIn = new Date(tokenData.expires_at).getTime() - Date.now();
      this.scheduleRevoke(cacheKey, blobUrl, expiresIn);

      return blobUrl;
    } catch (error) {
      console.error('Failed to create secure blob URL:', error);
      throw error;
    }
  }

  /**
   * Schedule automatic blob URL revocation
   */
  private scheduleRevoke(cacheKey: string, blobUrl: string, delayMs: number): void {
    // Clear existing timer if any
    if (this.revokeTimers.has(cacheKey)) {
      clearTimeout(this.revokeTimers.get(cacheKey)!);
    }

    // Schedule revocation
    const timer = setTimeout(() => {
      this.revokeBlobURL(cacheKey);
    }, delayMs);

    this.revokeTimers.set(cacheKey, timer);
  }

  /**
   * Manually revoke blob URL (cleanup)
   */
  revokeBlobURL(cacheKey: string): void {
    const blobUrl = this.blobCache.get(cacheKey);
    if (blobUrl) {
      URL.revokeObjectURL(blobUrl);
      this.blobCache.delete(cacheKey);
    }

    // Clear timer
    const timer = this.revokeTimers.get(cacheKey);
    if (timer) {
      clearTimeout(timer);
      this.revokeTimers.delete(cacheKey);
    }
  }

  /**
   * Cleanup all blob URLs (call on unmount)
   */
  cleanup(): void {
    // Revoke all blob URLs
    for (const blobUrl of this.blobCache.values()) {
      URL.revokeObjectURL(blobUrl);
    }

    // Clear all timers
    for (const timer of this.revokeTimers.values()) {
      clearTimeout(timer);
    }

    this.blobCache.clear();
    this.revokeTimers.clear();
  }

  /**
   * Get auth token from storage
   */
  private getAuthToken(): string {
    // Get token from localStorage or cookie
    return localStorage.getItem('auth_token') || '';
  }

  /**
   * Validate token before use
   */
  async validateToken(token: StreamToken): Promise<boolean> {
    try {
      const response = await fetch(`/api/v1/videos/${token.video_id}/validate-token`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ token }),
      });

      return response.ok;
    } catch (error) {
      console.error('Token validation failed:', error);
      return false;
    }
  }
}

// Singleton instance
export const secureVideoService = new SecureVideoService();