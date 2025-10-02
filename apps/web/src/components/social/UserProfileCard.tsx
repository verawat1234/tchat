/**
 * UserProfileCard Component
 *
 * Comprehensive user profile card with follow/unfollow functionality.
 * Demonstrates Claude Agent SDK patterns with full context awareness.
 *
 * Agent Loop Implementation:
 * 1. Gather Context: Fetch user profile, followers, following counts, relationship status
 * 2. Take Action: Execute follow/unfollow operations with real-time updates
 * 3. Verify Work: Validate changes and update all related UI components
 */

import React from 'react';
import { FollowButton } from './FollowButton';
import { useGetSocialProfileQuery, useGetFollowersQuery, useGetFollowingQuery } from '../../services/socialApi';
import type { SocialProfile } from '../../types/social';

export interface UserProfileCardProps {
  /** User ID to display profile for */
  userId: string;
  /** Current authenticated user's ID */
  currentUserId: string;
  /** Display variant */
  variant?: 'compact' | 'full';
  /** Show follow statistics */
  showStats?: boolean;
  /** Optional click handler */
  onClick?: () => void;
}

export const UserProfileCard: React.FC<UserProfileCardProps> = ({
  userId,
  currentUserId,
  variant = 'full',
  showStats = true,
  onClick,
}) => {
  // ========================================================================
  // Phase 1: Gather Context (Agent Loop - Context Gathering)
  // ========================================================================

  const {
    data: profile,
    isLoading: profileLoading,
    error: profileError,
  } = useGetSocialProfileQuery(userId, {
    skip: !userId,
  });

  const {
    data: followersData,
    isLoading: followersLoading,
  } = useGetFollowersQuery(
    { userId, limit: 100 },
    { skip: !userId || !showStats }
  );

  const {
    data: followingData,
    isLoading: followingLoading,
  } = useGetFollowingQuery(
    { userId, limit: 100 },
    { skip: !userId || !showStats }
  );

  // Loading state
  if (profileLoading) {
    return (
      <div className="animate-pulse bg-gray-200 dark:bg-gray-700 rounded-lg h-48" />
    );
  }

  // Error state
  if (profileError || !profile) {
    return (
      <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
        <p className="text-red-800 dark:text-red-200">Failed to load profile</p>
      </div>
    );
  }

  // ========================================================================
  // Phase 2: Display Context with Interactive Elements
  // ========================================================================

  const followerCount = followersData?.followers?.length || profile.followersCount || 0;
  const followingCount = followingData?.following?.length || profile.followingCount || 0;
  const isOwnProfile = currentUserId === userId;

  return (
    <div
      className={`
        bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden
        ${onClick ? 'cursor-pointer hover:shadow-lg transition-shadow' : ''}
        ${variant === 'compact' ? 'p-4' : 'p-6'}
      `}
      onClick={onClick}
    >
      {/* Profile Header */}
      <div className="flex items-start gap-4">
        {/* Avatar */}
        <div className="flex-shrink-0">
          <img
            src={profile.avatar || `https://ui-avatars.com/api/?name=${encodeURIComponent(profile.username)}&size=128`}
            alt={profile.username}
            className={`
              rounded-full object-cover
              ${variant === 'compact' ? 'w-16 h-16' : 'w-24 h-24'}
            `}
          />
          {profile.isSocialVerified && (
            <div className="absolute -bottom-1 -right-1 bg-blue-500 rounded-full p-1">
              <svg className="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
            </div>
          )}
        </div>

        {/* Profile Info */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white truncate">
              {profile.display_name || profile.name || profile.username}
            </h3>
            {profile.isSocialVerified && (
              <svg className="w-5 h-5 text-blue-500 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
            )}
          </div>
          <p className="text-sm text-gray-500 dark:text-gray-400 mb-2">
            @{profile.username}
          </p>

          {/* Bio */}
          {variant === 'full' && profile.bio && (
            <p className="text-sm text-gray-700 dark:text-gray-300 mb-4">
              {profile.bio}
            </p>
          )}

          {/* Statistics */}
          {showStats && variant === 'full' && (
            <div className="flex gap-6 mb-4">
              <div>
                <p className="text-2xl font-bold text-gray-900 dark:text-white">
                  {followersLoading ? '...' : followerCount}
                </p>
                <p className="text-sm text-gray-500 dark:text-gray-400">Followers</p>
              </div>
              <div>
                <p className="text-2xl font-bold text-gray-900 dark:text-white">
                  {followingLoading ? '...' : followingCount}
                </p>
                <p className="text-sm text-gray-500 dark:text-gray-400">Following</p>
              </div>
              {profile.postsCount !== undefined && (
                <div>
                  <p className="text-2xl font-bold text-gray-900 dark:text-white">
                    {profile.postsCount}
                  </p>
                  <p className="text-sm text-gray-500 dark:text-gray-400">Posts</p>
                </div>
              )}
            </div>
          )}

          {/* Follow Button (Phase 2: Take Action) */}
          {!isOwnProfile && (
            <FollowButton
              user={profile}
              currentUserId={currentUserId}
              initialIsFollowing={profile.isFollowing || false}
              variant="primary"
              size={variant === 'compact' ? 'sm' : 'md'}
              trackingSource="profile_card"
            />
          )}
        </div>
      </div>
    </div>
  );
};

export default UserProfileCard;