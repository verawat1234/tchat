/**
 * FollowButton Component
 *
 * Implements follow/unfollow functionality following Claude Agent SDK patterns:
 * 1. Gather Context: Check current follow status
 * 2. Take Action: Execute follow/unfollow with optimistic updates
 * 3. Verify Work: Validate operation success and provide feedback
 *
 * Features:
 * - Optimistic UI updates for instant feedback
 * - Automatic error recovery with rollback
 * - Loading states and success/error messaging
 * - Accessibility support
 * - Analytics tracking
 */

import React, { useState, useCallback } from 'react';
import { useFollowUserMutation, useUnfollowUserMutation } from '../../services/socialApi';
import type { SocialProfile } from '../../types/social';

export interface FollowButtonProps {
  /** User profile to follow/unfollow */
  user: SocialProfile;
  /** Current user's ID */
  currentUserId: string;
  /** Initial follow status */
  initialIsFollowing?: boolean;
  /** Button variant */
  variant?: 'primary' | 'secondary' | 'ghost';
  /** Button size */
  size?: 'sm' | 'md' | 'lg';
  /** Optional callback on follow */
  onFollow?: (userId: string) => void;
  /** Optional callback on unfollow */
  onUnfollow?: (userId: string) => void;
  /** Optional analytics tracking */
  trackingSource?: string;
}

export const FollowButton: React.FC<FollowButtonProps> = ({
  user,
  currentUserId,
  initialIsFollowing = false,
  variant = 'primary',
  size = 'md',
  onFollow,
  onUnfollow,
  trackingSource = 'profile',
}) => {
  // ========================================================================
  // Phase 1: Gather Context (Agent Loop - Context Gathering)
  // ========================================================================

  const [isFollowing, setIsFollowing] = useState(initialIsFollowing);
  const [optimisticUpdate, setOptimisticUpdate] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // RTK Query mutations with optimistic updates
  const [followUser, { isLoading: isFollowLoading }] = useFollowUserMutation();
  const [unfollowUser, { isLoading: isUnfollowLoading }] = useUnfollowUserMutation();

  const isLoading = isFollowLoading || isUnfollowLoading || optimisticUpdate;

  // ========================================================================
  // Phase 2: Take Action (Agent Loop - Action Execution)
  // ========================================================================

  const handleFollow = useCallback(async () => {
    try {
      // Clear any previous errors
      setError(null);

      // Optimistic update - immediate UI feedback
      setOptimisticUpdate(true);
      const previousState = isFollowing;
      setIsFollowing(true);

      // Execute follow action
      const result = await followUser({
        followerId: currentUserId,
        followingId: user.id,
        source: trackingSource,
      }).unwrap();

      // Phase 3: Verify Work (Agent Loop - Work Verification)
      if (result && result.message) {
        console.log('[FollowButton] Follow successful:', result.message);
        onFollow?.(user.id);

        // Analytics tracking
        if (typeof window !== 'undefined' && (window as any).gtag) {
          (window as any).gtag('event', 'follow_user', {
            event_category: 'social',
            event_label: user.username,
            value: 1,
          });
        }
      }
    } catch (err: any) {
      // Error recovery - rollback optimistic update
      console.error('[FollowButton] Follow failed:', err);
      setIsFollowing(false);
      setError(err?.data?.message || 'Failed to follow user. Please try again.');

      // Show error notification (could integrate with toast notification system)
      setTimeout(() => setError(null), 5000);
    } finally {
      setOptimisticUpdate(false);
    }
  }, [currentUserId, user.id, user.username, followUser, onFollow, trackingSource, isFollowing]);

  const handleUnfollow = useCallback(async () => {
    try {
      // Clear any previous errors
      setError(null);

      // Optimistic update - immediate UI feedback
      setOptimisticUpdate(true);
      const previousState = isFollowing;
      setIsFollowing(false);

      // Execute unfollow action
      const result = await unfollowUser({
        followerId: currentUserId,
        followingId: user.id,
      }).unwrap();

      // Phase 3: Verify Work (Agent Loop - Work Verification)
      if (result && result.message) {
        console.log('[FollowButton] Unfollow successful:', result.message);
        onUnfollow?.(user.id);

        // Analytics tracking
        if (typeof window !== 'undefined' && (window as any).gtag) {
          (window as any).gtag('event', 'unfollow_user', {
            event_category: 'social',
            event_label: user.username,
            value: -1,
          });
        }
      }
    } catch (err: any) {
      // Error recovery - rollback optimistic update
      console.error('[FollowButton] Unfollow failed:', err);
      setIsFollowing(true);
      setError(err?.data?.message || 'Failed to unfollow user. Please try again.');

      // Show error notification
      setTimeout(() => setError(null), 5000);
    } finally {
      setOptimisticUpdate(false);
    }
  }, [currentUserId, user.id, user.username, unfollowUser, onUnfollow, isFollowing]);

  const handleClick = useCallback(() => {
    if (isLoading) return;

    if (isFollowing) {
      handleUnfollow();
    } else {
      handleFollow();
    }
  }, [isFollowing, isLoading, handleFollow, handleUnfollow]);

  // ========================================================================
  // Render UI with Accessibility
  // ========================================================================

  // Variant styles
  const variantStyles = {
    primary: isFollowing
      ? 'bg-gray-200 text-gray-800 hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200'
      : 'bg-blue-500 text-white hover:bg-blue-600 dark:bg-blue-600 dark:hover:bg-blue-700',
    secondary: isFollowing
      ? 'bg-white text-gray-800 border border-gray-300 hover:bg-gray-50 dark:bg-gray-800 dark:text-gray-200 dark:border-gray-600'
      : 'bg-white text-blue-500 border border-blue-500 hover:bg-blue-50 dark:bg-gray-800 dark:text-blue-400 dark:border-blue-400',
    ghost: isFollowing
      ? 'bg-transparent text-gray-600 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800'
      : 'bg-transparent text-blue-500 hover:bg-blue-50 dark:text-blue-400 dark:hover:bg-gray-800',
  };

  // Size styles
  const sizeStyles = {
    sm: 'px-3 py-1 text-sm',
    md: 'px-4 py-2 text-base',
    lg: 'px-6 py-3 text-lg',
  };

  return (
    <div className="relative inline-block">
      <button
        onClick={handleClick}
        disabled={isLoading || currentUserId === user.id}
        className={`
          ${variantStyles[variant]}
          ${sizeStyles[size]}
          rounded-lg font-medium
          transition-all duration-200
          disabled:opacity-50 disabled:cursor-not-allowed
          focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
          ${isLoading ? 'cursor-wait' : 'cursor-pointer'}
        `}
        aria-label={isFollowing ? `Unfollow ${user.username}` : `Follow ${user.username}`}
        aria-pressed={isFollowing}
      >
        {isLoading ? (
          <span className="flex items-center gap-2">
            <svg
              className="animate-spin h-4 w-4"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                className="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                strokeWidth="4"
              />
              <path
                className="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              />
            </svg>
            {isFollowing ? 'Unfollowing...' : 'Following...'}
          </span>
        ) : (
          <span>{isFollowing ? 'Following' : 'Follow'}</span>
        )}
      </button>

      {/* Error message */}
      {error && (
        <div
          className="absolute top-full left-0 right-0 mt-2 p-2 bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200 text-sm rounded-md shadow-lg z-10"
          role="alert"
        >
          {error}
        </div>
      )}
    </div>
  );
};

export default FollowButton;