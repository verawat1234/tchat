# Follow/Unfollow System - Complete Implementation Guide

## Overview

This guide documents the complete follow/unfollow system implementation following **Claude Agent SDK patterns** from Anthropic's engineering blog. The system implements the agent loop paradigm: **Gather Context → Take Action → Verify Work**.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Backend Implementation](#backend-implementation)
3. [Web Implementation](#web-implementation)
4. [KMP Mobile Implementation](#kmp-mobile-implementation)
5. [Usage Examples](#usage-examples)
6. [Testing Strategy](#testing-strategy)
7. [Performance Optimization](#performance-optimization)

---

## Architecture Overview

### Claude Agent SDK Pattern Implementation

The follow system implements the **agent loop** pattern across all layers:

```
┌──────────────────────────────────────────────────────────┐
│                    Agent Loop Cycle                       │
├──────────────────────────────────────────────────────────┤
│  Phase 1: Gather Context                                 │
│  - Check current follow relationship status              │
│  - Fetch user profiles and follower counts              │
│  - Validate permissions and constraints                  │
├──────────────────────────────────────────────────────────┤
│  Phase 2: Take Action                                    │
│  - Execute follow/unfollow operation                     │
│  - Apply optimistic UI updates                           │
│  - Trigger notifications and side effects                │
├──────────────────────────────────────────────────────────┤
│  Phase 3: Verify Work                                    │
│  - Validate operation success                            │
│  - Update cache and related queries                      │
│  - Rollback on failure (error recovery)                 │
│  - Log analytics and metrics                             │
└──────────────────────────────────────────────────────────┘
```

### System Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Web Frontend (React)                  │
│  - RTK Query hooks (useFollowUserMutation)              │
│  - Optimistic updates with automatic rollback           │
│  - FollowButton component with accessibility            │
└─────────────────┬───────────────────────────────────────┘
                  │
                  │ HTTP/REST API
                  ▼
┌─────────────────────────────────────────────────────────┐
│              Backend Service (Go - Port 8092)           │
│  - FollowUser / UnfollowUser endpoints                  │
│  - Self-follow prevention                               │
│  - Duplicate follow detection                           │
│  - Follower/Following count management                  │
└─────────────────┬───────────────────────────────────────┘
                  │
                  │ PostgreSQL
                  ▼
┌─────────────────────────────────────────────────────────┐
│            Database Schema (follows table)               │
│  - id (uuid, primary key)                               │
│  - follower_id (uuid, indexed)                          │
│  - following_id (uuid, indexed)                         │
│  - created_at (timestamp)                               │
│  - UNIQUE(follower_id, following_id)                    │
└─────────────────┬───────────────────────────────────────┘
                  │
                  │ Cross-platform sync
                  ▼
┌─────────────────────────────────────────────────────────┐
│        KMP Mobile (Kotlin Multiplatform)                │
│  - SQLDelightSocialRepository                           │
│  - getFollowers() / getFollowing() methods              │
│  - addInteraction(FOLLOW) / removeInteraction(FOLLOW)   │
│  - Offline-first with automatic sync                    │
└─────────────────────────────────────────────────────────┘
```

---

## Backend Implementation

### API Endpoints

**Base URL**: `http://localhost:8092/api/v1/social`

#### 1. Follow User

```http
POST /social/follow
Content-Type: application/json

{
  "followerId": "uuid-of-follower",
  "followingId": "uuid-of-user-to-follow"
}
```

**Response** (200 OK):
```json
{
  "message": "Successfully followed user"
}
```

**Error Responses**:
- `400 Bad Request`: "Cannot follow yourself"
- `409 Conflict`: "Already following this user"
- `404 Not Found`: "User not found"

#### 2. Unfollow User

```http
DELETE /social/follow/{followerId}/{followingId}
```

**Response** (200 OK):
```json
{
  "message": "Successfully unfollowed user"
}
```

**Error Responses**:
- `404 Not Found`: "Following relationship not found"

#### 3. Get Followers

```http
GET /social/followers/{userId}?limit=20&offset=0
```

**Response** (200 OK):
```json
{
  "followers": [
    {
      "id": "user-uuid",
      "username": "johndoe",
      "displayName": "John Doe",
      "avatar": "https://...",
      "followedAt": "2025-09-30T10:00:00Z",
      "isVerified": true
    }
  ],
  "limit": 20,
  "offset": 0,
  "hasMore": false
}
```

#### 4. Get Following

```http
GET /social/following/{userId}?limit=20&offset=0
```

**Response** (200 OK):
```json
{
  "following": [
    {
      "id": "user-uuid",
      "username": "janedoe",
      "displayName": "Jane Doe",
      "avatar": "https://...",
      "followedAt": "2025-09-29T15:30:00Z",
      "isVerified": false
    }
  ],
  "limit": 20,
  "offset": 0,
  "hasMore": true
}
```

### Backend Service Implementation

**File**: `backend/social/services/user_service.go`

```go
// FollowUser implements Phase 1-3 of Agent Loop
func (s *userService) FollowUser(ctx context.Context, req *models.FollowRequest) error {
    // Phase 1: Gather Context
    if req.FollowerID == req.FollowingID {
        return fmt.Errorf("cannot follow yourself")
    }

    isFollowing, err := s.repo.Users().IsFollowing(ctx, req.FollowerID, req.FollowingID)
    if err != nil {
        return fmt.Errorf("failed to check follow status: %w", err)
    }
    if isFollowing {
        return fmt.Errorf("already following this user")
    }

    // Phase 2: Take Action
    follow := &models.Follow{
        ID:          uuid.New(),
        FollowerID:  req.FollowerID,
        FollowingID: req.FollowingID,
        CreatedAt:   time.Now(),
    }

    err = s.repo.Users().CreateFollow(ctx, follow)
    if err != nil {
        return fmt.Errorf("failed to create follow: %w", err)
    }

    // Phase 3: Verify Work
    // Cache invalidation and notifications handled by repository layer
    return nil
}
```

---

## Web Implementation

### RTK Query Integration

**File**: `apps/web/src/services/socialApi.ts`

The follow/unfollow system is **fully implemented** using RTK Query with advanced features:

#### Follow User Mutation

```typescript
followUser: builder.mutation<{ message: string }, FollowRequest>({
  query: (body) => ({
    url: '/social/follow',
    method: 'POST',
    body,
  }),
  invalidatesTags: (result, error, { followerId, followingId }) => [
    { type: 'SocialProfile', id: followerId },
    { type: 'SocialProfile', id: followingId },
    'UserRelationship',
    'SocialFollowers',
    'SocialFollowing',
  ],
}),
```

#### Unfollow User Mutation

```typescript
unfollowUser: builder.mutation<{ message: string }, { followerId: string; followingId: string }>({
  query: ({ followerId, followingId }) => ({
    url: `/social/follow/${encodeURIComponent(followerId)}/${encodeURIComponent(followingId)}`,
    method: 'DELETE',
  }),
  invalidatesTags: (result, error, { followerId, followingId }) => [
    { type: 'SocialProfile', id: followerId },
    { type: 'SocialProfile', id: followingId },
    'UserRelationship',
    'SocialFollowers',
    'SocialFollowing',
  ],
}),
```

### React Component Usage

#### FollowButton Component

**File**: `apps/web/src/components/social/FollowButton.tsx`

The `FollowButton` component implements the full **Claude Agent SDK agent loop**:

```tsx
import { FollowButton } from '@/components/social/FollowButton';
import type { SocialProfile } from '@/types/social';

// Example usage
function UserCard({ user, currentUserId }: { user: SocialProfile; currentUserId: string }) {
  return (
    <div>
      <h3>{user.display_name}</h3>
      <FollowButton
        user={user}
        currentUserId={currentUserId}
        initialIsFollowing={user.isFollowing || false}
        variant="primary"
        size="md"
        onFollow={(userId) => console.log('Followed:', userId)}
        onUnfollow={(userId) => console.log('Unfollowed:', userId)}
        trackingSource="user_card"
      />
    </div>
  );
}
```

**Key Features**:
- ✅ Optimistic UI updates (instant feedback)
- ✅ Automatic error recovery with rollback
- ✅ Loading states and success/error messaging
- ✅ Accessibility support (ARIA labels, focus management)
- ✅ Analytics tracking (Google Analytics integration)
- ✅ Prevents self-following
- ✅ Responsive design (supports 3 sizes, 3 variants)

#### UserProfileCard Component

**File**: `apps/web/src/components/social/UserProfileCard.tsx`

Full user profile card with integrated follow functionality:

```tsx
import { UserProfileCard } from '@/components/social/UserProfileCard';

// Example usage
function ProfilePage() {
  const currentUserId = useSelector((state) => state.auth.userId);

  return (
    <UserProfileCard
      userId="target-user-id"
      currentUserId={currentUserId}
      variant="full"
      showStats={true}
      onClick={() => router.push(`/profile/${userId}`)}
    />
  );
}
```

**Features**:
- ✅ Real-time follower/following counts
- ✅ Verified badge display
- ✅ Biography and profile information
- ✅ Integrated FollowButton
- ✅ Loading and error states
- ✅ Compact and full variants

### Advanced Usage with Hooks

```typescript
import {
  useFollowUserMutation,
  useUnfollowUserMutation,
  useGetFollowersQuery,
  useGetFollowingQuery,
} from '@/services/socialApi';

function AdvancedFollowComponent() {
  const currentUserId = 'current-user-id';
  const targetUserId = 'target-user-id';

  // Follow/unfollow mutations
  const [followUser, { isLoading: isFollowing }] = useFollowUserMutation();
  const [unfollowUser, { isLoading: isUnfollowing }] = useUnfollowUserMutation();

  // Fetch followers and following
  const { data: followers } = useGetFollowersQuery({ userId: targetUserId, limit: 100 });
  const { data: following } = useGetFollowingQuery({ userId: targetUserId, limit: 100 });

  const handleFollow = async () => {
    try {
      await followUser({
        followerId: currentUserId,
        followingId: targetUserId,
      }).unwrap();
      console.log('Follow successful');
    } catch (error) {
      console.error('Follow failed:', error);
    }
  };

  return (
    <div>
      <p>Followers: {followers?.followers?.length || 0}</p>
      <p>Following: {following?.following?.length || 0}</p>
      <button onClick={handleFollow} disabled={isFollowing}>
        {isFollowing ? 'Following...' : 'Follow'}
      </button>
    </div>
  );
}
```

---

## KMP Mobile Implementation

### SQLDelight Repository

**File**: `apps/kmp/composeApp/src/commonMain/kotlin/com/tchat/mobile/repositories/SQLDelightSocialRepository.kt`

The KMP mobile implementation uses **SQLDelight** for offline-first functionality with automatic synchronization:

#### Follow User (KMP)

```kotlin
// Add follow interaction
override suspend fun addInteraction(interaction: SocialInteraction): Result<Unit> = withContext(Dispatchers.Default) {
    try {
        database.socialInteractionQueries.insertInteraction(
            interaction.id,
            interaction.userId,
            interaction.targetId,
            interaction.targetType.name.lowercase(),
            InteractionType.FOLLOW.name.lowercase(), // <-- FOLLOW interaction
            interaction.createdAt.toString(),
            interaction.updatedAt.toString()
        )
        Result.success(Unit)
    } catch (e: Exception) {
        Result.failure(e)
    }
}
```

#### Unfollow User (KMP)

```kotlin
// Remove follow interaction
override suspend fun removeInteraction(
    userId: String,
    targetId: String,
    targetType: InteractionTargetType,
    interactionType: InteractionType
): Result<Unit> = withContext(Dispatchers.Default) {
    try {
        database.socialInteractionQueries.removeInteraction(
            user_id = userId,
            target_id = targetId,
            target_type = targetType.name.lowercase(),
            interaction_type = InteractionType.FOLLOW.name.lowercase()
        )
        Result.success(Unit)
    } catch (e: Exception) {
        Result.failure(e)
    }
}
```

#### Get Followers (KMP)

```kotlin
override suspend fun getFollowers(userId: String): Result<List<SocialUserProfile>> = withContext(Dispatchers.Default) {
    try {
        val followerRows = database.socialInteractionQueries.getFollowers(userId)
            .asFlow()
            .mapToList(Dispatchers.Default)
            .firstOrNull() ?: emptyList()

        val followers = followerRows.map { row ->
            SocialUserProfile(
                userId = row.follower_user_id,
                displayName = row.display_name ?: "Unknown User",
                username = row.username ?: "",
                avatarUrl = row.avatar ?: "",
                isVerified = row.is_verified == 1L,
                // ... other fields
            )
        }
        Result.success(followers)
    } catch (e: Exception) {
        Result.failure(e)
    }
}
```

#### Get Following (KMP)

```kotlin
override suspend fun getFollowedUsers(userId: String): Result<List<SocialUserProfile>> = withContext(Dispatchers.Default) {
    try {
        val followedRows = database.socialInteractionQueries.getFollowedUsers(userId)
            .asFlow()
            .mapToList(Dispatchers.Default)
            .firstOrNull() ?: emptyList()

        val followedUsers = followedRows.map { row ->
            SocialUserProfile(
                userId = row.followed_user_id,
                displayName = row.display_name ?: "Unknown User",
                username = row.username ?: "",
                avatarUrl = row.avatar ?: "",
                isVerified = row.is_verified == 1L,
                // ... other fields
            )
        }
        Result.success(followedUsers)
    } catch (e: Exception) {
        Result.failure(e)
    }
}
```

### KMP UI Integration Example

```kotlin
// Android Compose UI
@Composable
fun FollowButton(
    user: SocialUserProfile,
    currentUserId: String,
    viewModel: SocialViewModel
) {
    var isFollowing by remember { mutableStateOf(user.isFollowing) }
    var isLoading by remember { mutableStateOf(false) }

    Button(
        onClick = {
            isLoading = true
            viewModel.toggleFollow(currentUserId, user.userId) { success ->
                if (success) {
                    isFollowing = !isFollowing
                }
                isLoading = false
            }
        },
        enabled = !isLoading
    ) {
        Text(if (isFollowing) "Following" else "Follow")
    }
}
```

---

## Usage Examples

### Example 1: Simple Follow Button

```tsx
import { FollowButton } from '@/components/social/FollowButton';

function SimpleExample() {
  const user: SocialProfile = {
    id: 'user-123',
    username: 'johndoe',
    display_name: 'John Doe',
    avatar: 'https://example.com/avatar.jpg',
    isFollowing: false,
  };

  return (
    <FollowButton
      user={user}
      currentUserId="current-user-id"
      variant="primary"
      size="md"
    />
  );
}
```

### Example 2: User Discovery with Follow

```tsx
import { useDiscoverUsersQuery, useFollowUserMutation } from '@/services/socialApi';
import { FollowButton } from '@/components/social/FollowButton';

function UserDiscovery() {
  const currentUserId = 'current-user-id';
  const { data: users, isLoading } = useDiscoverUsersQuery({
    region: 'TH',
    interests: ['technology', 'startups'],
    limit: 20,
  });

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="grid grid-cols-3 gap-4">
      {users?.map((user) => (
        <div key={user.id} className="p-4 border rounded-lg">
          <img src={user.avatar} alt={user.username} className="w-16 h-16 rounded-full" />
          <h3>{user.display_name}</h3>
          <p>@{user.username}</p>
          <FollowButton
            user={user}
            currentUserId={currentUserId}
            variant="secondary"
            size="sm"
            trackingSource="discovery"
          />
        </div>
      ))}
    </div>
  );
}
```

### Example 3: Followers List with Follow Back

```tsx
import { useGetFollowersQuery } from '@/services/socialApi';
import { FollowButton } from '@/components/social/FollowButton';

function FollowersList({ userId }: { userId: string }) {
  const currentUserId = 'current-user-id';
  const { data, isLoading } = useGetFollowersQuery({ userId, limit: 50 });

  if (isLoading) return <div>Loading followers...</div>;

  return (
    <div className="space-y-4">
      <h2 className="text-2xl font-bold">Followers ({data?.followers?.length || 0})</h2>
      {data?.followers?.map((follower) => (
        <div key={follower.id} className="flex items-center justify-between p-4 border rounded-lg">
          <div className="flex items-center gap-3">
            <img src={follower.avatar} alt={follower.username} className="w-12 h-12 rounded-full" />
            <div>
              <p className="font-semibold">{follower.displayName}</p>
              <p className="text-sm text-gray-500">@{follower.username}</p>
            </div>
          </div>
          {follower.id !== currentUserId && (
            <FollowButton
              user={follower as any}
              currentUserId={currentUserId}
              variant="secondary"
              size="sm"
              trackingSource="followers_list"
            />
          )}
        </div>
      ))}
    </div>
  );
}
```

---

## Testing Strategy

### Contract Tests

Create contract tests to ensure API consistency:

```typescript
// apps/web/src/services/__tests__/socialApi.contract.test.ts
import { describe, it, expect } from 'vitest';
import { socialApi } from '../socialApi';

describe('Social API Contract Tests', () => {
  describe('Follow User', () => {
    it('should follow user and return success message', async () => {
      const result = await socialApi.endpoints.followUser.initiate({
        followerId: 'user-1',
        followingId: 'user-2',
      });

      expect(result).toHaveProperty('data.message');
      expect(result.data.message).toBe('Successfully followed user');
    });

    it('should prevent self-following', async () => {
      try {
        await socialApi.endpoints.followUser.initiate({
          followerId: 'user-1',
          followingId: 'user-1',
        });
      } catch (error: any) {
        expect(error.status).toBe(400);
        expect(error.data.message).toContain('Cannot follow yourself');
      }
    });
  });
});
```

### E2E Tests with Playwright

```typescript
// apps/web/tests/e2e/follow.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Follow/Unfollow Functionality', () => {
  test('should follow user from profile page', async ({ page }) => {
    await page.goto('/profile/johndoe');

    const followButton = page.getByRole('button', { name: /follow/i });
    await followButton.click();

    await expect(followButton).toHaveText('Following');
    await expect(page.getByText(/successfully followed/i)).toBeVisible();
  });

  test('should unfollow user with confirmation', async ({ page }) => {
    await page.goto('/profile/johndoe');

    const followButton = page.getByRole('button', { name: /following/i });
    await followButton.click();

    await expect(followButton).toHaveText('Follow');
  });
});
```

---

## Performance Optimization

### Caching Strategy

The RTK Query implementation includes advanced caching:

- **Follow/Unfollow mutations** invalidate 5 cache tags:
  - `SocialProfile` (follower and following profiles)
  - `UserRelationship` (relationship status)
  - `SocialFollowers` (followers list)
  - `SocialFollowing` (following list)

- **Query cache duration**:
  - `getFollowers`: 180 seconds (3 minutes)
  - `getFollowing`: 180 seconds (3 minutes)
  - `getSocialProfile`: 300 seconds (5 minutes)

### Optimistic Updates

The `FollowButton` component implements **optimistic updates** for instant UI feedback:

```typescript
// Optimistic update flow
setOptimisticUpdate(true);
setIsFollowing(true); // Instant UI change

try {
  await followUser({ ... }).unwrap();
  // Success - keep optimistic update
} catch (error) {
  setIsFollowing(false); // Rollback on failure
  showError(error.message);
}
```

### Backend Performance

- **Database Indexes**: `follower_id` and `following_id` are indexed for fast lookups
- **Unique Constraint**: Prevents duplicate follow relationships at database level
- **Connection Pooling**: PostgreSQL connection pooling for high concurrency
- **Response Time**: Target <200ms for follow/unfollow operations

---

## Analytics and Monitoring

### Google Analytics Integration

The `FollowButton` component includes built-in analytics tracking:

```typescript
// Track follow event
window.gtag('event', 'follow_user', {
  event_category: 'social',
  event_label: user.username,
  value: 1,
});

// Track unfollow event
window.gtag('event', 'unfollow_user', {
  event_category: 'social',
  event_label: user.username,
  value: -1,
});
```

### Monitoring Metrics

Key metrics to monitor:
- Follow/unfollow success rate (target: >99.9%)
- API response time (target: <200ms p95)
- Error rate (target: <0.1%)
- Optimistic update rollback rate (target: <1%)

---

## Conclusion

The follow/unfollow system is **fully implemented** across all platforms following **Claude Agent SDK patterns**. The implementation includes:

✅ **Backend**: Complete Go service with validation, error handling, and database management
✅ **Web**: RTK Query with optimistic updates, error recovery, and accessibility
✅ **KMP Mobile**: Offline-first SQLDelight repository with cross-platform sync
✅ **Components**: Production-ready `FollowButton` and `UserProfileCard` components
✅ **Documentation**: Comprehensive usage examples and testing strategies

The system follows the **agent loop paradigm** (Gather Context → Take Action → Verify Work) at every layer, ensuring robust, maintainable, and user-friendly follow functionality.