import { store } from '../store';
import { api } from './api';
import { authApi } from './auth';
import { usersApi } from './users';
import { messagesApi } from './messages';
import { chatsApi } from './chats';

/**
 * Prefetch Service
 *
 * Proactively loads data that users are likely to need,
 * improving perceived performance and user experience.
 */

export class PrefetchService {
  private static instance: PrefetchService;
  private prefetchCache = new Set<string>();
  private prefetchTimeouts = new Map<string, number>();

  static getInstance(): PrefetchService {
    if (!PrefetchService.instance) {
      PrefetchService.instance = new PrefetchService();
    }
    return PrefetchService.instance;
  }

  /**
   * Prefetch user profile data after successful authentication
   */
  async prefetchAuthenticatedUserData(userId: string) {
    const cacheKey = `auth-user-${userId}`;
    if (this.prefetchCache.has(cacheKey)) return;

    try {
      // Prefetch current user profile
      store.dispatch(authApi.endpoints.getCurrentUser.initiate());

      // Prefetch user's chats
      store.dispatch(chatsApi.endpoints.listChats.initiate({ limit: 20 }));

      // Mark as prefetched
      this.prefetchCache.add(cacheKey);

      // Clear cache after 5 minutes
      setTimeout(() => {
        this.prefetchCache.delete(cacheKey);
      }, 5 * 60 * 1000);
    } catch (error) {
      console.warn('Failed to prefetch authenticated user data:', error);
    }
  }

  /**
   * Prefetch chat messages when user hovers over or focuses on a chat
   */
  async prefetchChatMessages(chatId: string, delay = 500) {
    const cacheKey = `chat-messages-${chatId}`;
    if (this.prefetchCache.has(cacheKey)) return;

    // Clear any existing timeout for this chat
    const existingTimeout = this.prefetchTimeouts.get(cacheKey);
    if (existingTimeout) {
      clearTimeout(existingTimeout);
    }

    // Set a delay to avoid excessive prefetching on rapid hovering
    const timeoutId = setTimeout(async () => {
      try {
        store.dispatch(messagesApi.endpoints.listMessages.initiate({
          chatId,
          limit: 20,
        }));

        this.prefetchCache.add(cacheKey);
        this.prefetchTimeouts.delete(cacheKey);

        // Clear cache after 2 minutes
        setTimeout(() => {
          this.prefetchCache.delete(cacheKey);
        }, 2 * 60 * 1000);
      } catch (error) {
        console.warn(`Failed to prefetch messages for chat ${chatId}:`, error);
      }
    }, delay);

    this.prefetchTimeouts.set(cacheKey, timeoutId);
  }

  /**
   * Prefetch user profile when mentioned or hovered
   */
  async prefetchUserProfile(userId: string, delay = 300) {
    const cacheKey = `user-profile-${userId}`;
    if (this.prefetchCache.has(cacheKey)) return;

    const existingTimeout = this.prefetchTimeouts.get(cacheKey);
    if (existingTimeout) {
      clearTimeout(existingTimeout);
    }

    const timeoutId = setTimeout(async () => {
      try {
        store.dispatch(usersApi.endpoints.getUserById.initiate(userId));

        this.prefetchCache.add(cacheKey);
        this.prefetchTimeouts.delete(cacheKey);

        // Clear cache after 1 minute
        setTimeout(() => {
          this.prefetchCache.delete(cacheKey);
        }, 60 * 1000);
      } catch (error) {
        console.warn(`Failed to prefetch user profile ${userId}:`, error);
      }
    }, delay);

    this.prefetchTimeouts.set(cacheKey, timeoutId);
  }

  /**
   * Prefetch common app data on startup
   */
  async prefetchCommonData() {
    const cacheKey = 'common-startup-data';
    if (this.prefetchCache.has(cacheKey)) return;

    try {
      // Get current user state
      const state = store.getState();
      const isAuthenticated = state.auth.isAuthenticated;

      if (isAuthenticated) {
        // Prefetch user's chat list
        store.dispatch(chatsApi.endpoints.listChats.initiate({ limit: 10 }));

        // Prefetch recent users for mentions/DMs
        store.dispatch(usersApi.endpoints.listUsers.initiate({ limit: 10 }));
      }

      this.prefetchCache.add(cacheKey);

      // Clear cache after 10 minutes
      setTimeout(() => {
        this.prefetchCache.delete(cacheKey);
      }, 10 * 60 * 1000);
    } catch (error) {
      console.warn('Failed to prefetch common data:', error);
    }
  }

  /**
   * Cancel prefetch for a specific resource
   */
  cancelPrefetch(cacheKey: string) {
    const timeoutId = this.prefetchTimeouts.get(cacheKey);
    if (timeoutId) {
      clearTimeout(timeoutId);
      this.prefetchTimeouts.delete(cacheKey);
    }
  }

  /**
   * Clear all prefetch caches and timeouts
   */
  clearAll() {
    this.prefetchCache.clear();
    this.prefetchTimeouts.forEach(timeoutId => clearTimeout(timeoutId));
    this.prefetchTimeouts.clear();
  }

  /**
   * Get prefetch statistics for debugging
   */
  getStats() {
    return {
      cachedItems: this.prefetchCache.size,
      pendingTimeouts: this.prefetchTimeouts.size,
      cachedKeys: Array.from(this.prefetchCache),
    };
  }
}

// Export singleton instance
export const prefetchService = PrefetchService.getInstance();

// React hooks for easy prefetching
export const usePrefetch = () => {
  return {
    prefetchChatMessages: (chatId: string, delay?: number) =>
      prefetchService.prefetchChatMessages(chatId, delay),
    prefetchUserProfile: (userId: string, delay?: number) =>
      prefetchService.prefetchUserProfile(userId, delay),
    cancelPrefetch: (cacheKey: string) =>
      prefetchService.cancelPrefetch(cacheKey),
  };
};