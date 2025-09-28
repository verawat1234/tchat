package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"tchat/social/models"
	"tchat/social/repository"
)

// MobileSyncService provides synchronization patterns for KMP mobile clients
type MobileSyncService interface {
	// Incremental sync operations
	GetUserProfileChanges(ctx context.Context, userID uuid.UUID, since time.Time) (*MobileSyncResponse, error)
	GetPostChanges(ctx context.Context, userID uuid.UUID, since time.Time) (*MobileSyncResponse, error)
	GetFollowChanges(ctx context.Context, userID uuid.UUID, since time.Time) (*MobileSyncResponse, error)

	// Bulk sync operations for initial data load
	GetInitialUserData(ctx context.Context, userID uuid.UUID) (*InitialSyncResponse, error)
	GetUserFeed(ctx context.Context, userID uuid.UUID, limit, offset int) (*FeedSyncResponse, error)

	// Conflict resolution for offline changes
	ResolveProfileConflicts(ctx context.Context, userID uuid.UUID, clientData *models.SocialProfile, lastSync time.Time) (*ConflictResolution, error)
	ApplyClientChanges(ctx context.Context, userID uuid.UUID, changes *ClientChanges) (*SyncResult, error)

	// Mobile-optimized discovery
	GetDiscoveryFeed(ctx context.Context, userID uuid.UUID, region string, limit int) ([]*models.SocialProfile, error)
	GetTrendingContent(ctx context.Context, region string, limit int) ([]*models.Post, error)
}

// mobileSyncService implements mobile synchronization patterns
type mobileSyncService struct {
	repo repository.RepositoryManager
}

// NewMobileSyncService creates a new mobile sync service
func NewMobileSyncService(repo repository.RepositoryManager) MobileSyncService {
	return &mobileSyncService{
		repo: repo,
	}
}

// GetUserProfileChanges returns incremental profile changes for mobile sync
func (s *mobileSyncService) GetUserProfileChanges(ctx context.Context, userID uuid.UUID, since time.Time) (*MobileSyncResponse, error) {
	// Get current profile
	profile, err := s.repo.Users().GetSocialProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	response := &MobileSyncResponse{
		UserID:      userID,
		SyncTime:    time.Now(),
		LastSync:    since,
		HasChanges:  false,
		ProfileData: profile,
	}

	// Check if profile was modified since last sync
	if profile != nil && profile.SocialUpdatedAt.After(since) {
		response.HasChanges = true
		response.ChangeType = "profile_updated"
		response.Changes = []ChangeItem{
			{
				Type:      "profile",
				Action:    "update",
				EntityID:  profile.ID.String(),
				Timestamp: profile.SocialUpdatedAt,
				Data:      profile,
			},
		}
	}

	return response, nil
}

// GetPostChanges returns incremental post changes for mobile sync
func (s *mobileSyncService) GetPostChanges(ctx context.Context, userID uuid.UUID, since time.Time) (*MobileSyncResponse, error) {
	// This would integrate with the post repository when available
	// For now, return structure for KMP compatibility testing
	response := &MobileSyncResponse{
		UserID:     userID,
		SyncTime:   time.Now(),
		LastSync:   since,
		HasChanges: false,
		Changes:    []ChangeItem{},
	}

	return response, nil
}

// GetFollowChanges returns incremental follow changes for mobile sync
func (s *mobileSyncService) GetFollowChanges(ctx context.Context, userID uuid.UUID, since time.Time) (*MobileSyncResponse, error) {
	// Get recent follow changes
	response := &MobileSyncResponse{
		UserID:     userID,
		SyncTime:   time.Now(),
		LastSync:   since,
		HasChanges: false,
		Changes:    []ChangeItem{},
	}

	// This would query for follows/unfollows since the last sync
	// Implementation depends on repository patterns

	return response, nil
}

// GetInitialUserData provides complete user data for initial app load
func (s *mobileSyncService) GetInitialUserData(ctx context.Context, userID uuid.UUID) (*InitialSyncResponse, error) {
	// Get user profile
	profile, err := s.repo.Users().GetSocialProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Get followers and following lists (limited for mobile)
	followers, err := s.repo.Users().GetUserFollowers(ctx, userID, 50, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	following, err := s.repo.Users().GetUserFollowing(ctx, userID, 50, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}

	response := &InitialSyncResponse{
		UserID:      userID,
		SyncTime:    time.Now(),
		Profile:     profile,
		Followers:   followers,
		Following:   following,
		MetaData: InitialSyncMetadata{
			TotalFollowers: profile.FollowersCount,
			TotalFollowing: profile.FollowingCount,
			TotalPosts:     profile.PostsCount,
			LastActivity:   profile.SocialUpdatedAt,
		},
	}

	return response, nil
}

// GetUserFeed provides personalized user feed for mobile
func (s *mobileSyncService) GetUserFeed(ctx context.Context, userID uuid.UUID, limit, offset int) (*FeedSyncResponse, error) {
	// Mobile-optimized feed with limited data
	response := &FeedSyncResponse{
		UserID:   userID,
		SyncTime: time.Now(),
		Posts:    []*models.Post{}, // Would be populated from post repository
		HasMore:  false,
		NextOffset: offset + limit,
		FeedMetadata: FeedMetadata{
			Algorithm:     "chronological",
			Region:        "", // Would be derived from user profile
			PersonalizedScore: 0.8,
		},
	}

	return response, nil
}

// ResolveProfileConflicts handles conflicts between server and client data
func (s *mobileSyncService) ResolveProfileConflicts(ctx context.Context, userID uuid.UUID, clientData *models.SocialProfile, lastSync time.Time) (*ConflictResolution, error) {
	// Get current server data
	serverData, err := s.repo.Users().GetSocialProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get server profile: %w", err)
	}

	resolution := &ConflictResolution{
		UserID:       userID,
		ResolvedAt:   time.Now(),
		HasConflicts: false,
		Resolution:   "no_conflict",
		MergedData:   serverData,
	}

	// Check for conflicts (server data modified after client's last sync)
	if serverData != nil && serverData.SocialUpdatedAt.After(lastSync) {
		resolution.HasConflicts = true
		resolution.Resolution = "server_wins" // Default conflict resolution strategy

		// Identify specific conflicts
		conflicts := []FieldConflict{}

		if clientData.DisplayName != serverData.DisplayName {
			conflicts = append(conflicts, FieldConflict{
				Field:       "displayName",
				ClientValue: clientData.DisplayName,
				ServerValue: serverData.DisplayName,
				Resolution:  "server_wins",
			})
		}

		if clientData.Bio != serverData.Bio {
			conflicts = append(conflicts, FieldConflict{
				Field:       "bio",
				ClientValue: clientData.Bio,
				ServerValue: serverData.Bio,
				Resolution:  "server_wins",
			})
		}

		resolution.Conflicts = conflicts
	}

	return resolution, nil
}

// ApplyClientChanges applies validated client changes to server data
func (s *mobileSyncService) ApplyClientChanges(ctx context.Context, userID uuid.UUID, changes *ClientChanges) (*SyncResult, error) {
	result := &SyncResult{
		UserID:         userID,
		ProcessedAt:    time.Now(),
		Success:        true,
		AppliedChanges: 0,
		FailedChanges:  0,
		Errors:         []string{},
	}

	// Apply profile changes if any
	if changes.ProfileChanges != nil {
		updateReq := &models.UpdateSocialProfileRequest{
			DisplayName: changes.ProfileChanges.DisplayName,
			Bio:         changes.ProfileChanges.Bio,
			Interests:   changes.ProfileChanges.Interests,
			SocialLinks: changes.ProfileChanges.SocialLinks,
		}

		if err := s.repo.Users().UpdateSocialProfile(ctx, userID, updateReq); err != nil {
			result.Success = false
			result.FailedChanges++
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to update profile: %v", err))
		} else {
			result.AppliedChanges++
		}
	}

	// Apply follow changes
	for _, followChange := range changes.FollowChanges {
		switch followChange.Action {
		case "follow":
			follow := &models.Follow{
				FollowerID:  userID,
				FollowingID: followChange.TargetUserID,
			}
			if err := s.repo.Users().CreateFollow(ctx, follow); err != nil {
				result.FailedChanges++
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to create follow: %v", err))
			} else {
				result.AppliedChanges++
			}
		case "unfollow":
			if err := s.repo.Users().DeleteFollow(ctx, userID, followChange.TargetUserID); err != nil {
				result.FailedChanges++
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to delete follow: %v", err))
			} else {
				result.AppliedChanges++
			}
		default:
			result.FailedChanges++
			result.Errors = append(result.Errors, fmt.Sprintf("Unknown follow action: %s", followChange.Action))
		}
	}

	return result, nil
}

// GetDiscoveryFeed provides mobile-optimized user discovery
func (s *mobileSyncService) GetDiscoveryFeed(ctx context.Context, userID uuid.UUID, region string, limit int) ([]*models.SocialProfile, error) {
	// Get current user profile for personalization
	currentUser, err := s.repo.Users().GetSocialProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	// Build discovery request based on user's interests and region
	discoveryReq := &models.UserDiscoveryRequest{
		Region:    region,
		Interests: currentUser.Interests,
		Limit:     limit,
		Offset:    0,
	}

	// Get discovery results
	profiles, err := s.repo.Users().DiscoverUsers(ctx, discoveryReq)
	if err != nil {
		return nil, fmt.Errorf("failed to discover users: %w", err)
	}

	return profiles, nil
}

// GetTrendingContent provides regional trending content for mobile
func (s *mobileSyncService) GetTrendingContent(ctx context.Context, region string, limit int) ([]*models.Post, error) {
	// This would integrate with post repository for trending content
	// For now, return empty slice for KMP compatibility
	return []*models.Post{}, nil
}

// Data structures for mobile sync operations

// MobileSyncResponse represents incremental sync data for mobile clients
type MobileSyncResponse struct {
	UserID      uuid.UUID   `json:"userId"`
	SyncTime    time.Time   `json:"syncTime"`
	LastSync    time.Time   `json:"lastSync"`
	HasChanges  bool        `json:"hasChanges"`
	ChangeType  string      `json:"changeType,omitempty"`
	Changes     []ChangeItem `json:"changes"`
	ProfileData *models.SocialProfile `json:"profileData,omitempty"`
}

// ChangeItem represents a single change in sync data
type ChangeItem struct {
	Type      string      `json:"type"`      // "profile", "post", "follow", etc.
	Action    string      `json:"action"`    // "create", "update", "delete"
	EntityID  string      `json:"entityId"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// InitialSyncResponse provides complete user data for app initialization
type InitialSyncResponse struct {
	UserID    uuid.UUID               `json:"userId"`
	SyncTime  time.Time               `json:"syncTime"`
	Profile   *models.SocialProfile   `json:"profile"`
	Followers []*models.Follow        `json:"followers"`
	Following []*models.Follow        `json:"following"`
	MetaData  InitialSyncMetadata     `json:"metadata"`
}

// InitialSyncMetadata provides context for initial sync
type InitialSyncMetadata struct {
	TotalFollowers int       `json:"totalFollowers"`
	TotalFollowing int       `json:"totalFollowing"`
	TotalPosts     int       `json:"totalPosts"`
	LastActivity   time.Time `json:"lastActivity"`
}

// FeedSyncResponse provides paginated feed data for mobile
type FeedSyncResponse struct {
	UserID       uuid.UUID    `json:"userId"`
	SyncTime     time.Time    `json:"syncTime"`
	Posts        []*models.Post `json:"posts"`
	HasMore      bool         `json:"hasMore"`
	NextOffset   int          `json:"nextOffset"`
	FeedMetadata FeedMetadata `json:"metadata"`
}

// FeedMetadata provides context for feed algorithms
type FeedMetadata struct {
	Algorithm         string  `json:"algorithm"`
	Region            string  `json:"region"`
	PersonalizedScore float64 `json:"personalizedScore"`
}

// ConflictResolution handles data conflicts between client and server
type ConflictResolution struct {
	UserID       uuid.UUID        `json:"userId"`
	ResolvedAt   time.Time        `json:"resolvedAt"`
	HasConflicts bool             `json:"hasConflicts"`
	Resolution   string           `json:"resolution"` // "server_wins", "client_wins", "merged"
	MergedData   *models.SocialProfile `json:"mergedData"`
	Conflicts    []FieldConflict  `json:"conflicts"`
}

// FieldConflict represents a specific field conflict
type FieldConflict struct {
	Field       string      `json:"field"`
	ClientValue interface{} `json:"clientValue"`
	ServerValue interface{} `json:"serverValue"`
	Resolution  string      `json:"resolution"`
}

// ClientChanges represents changes from mobile client to be applied
type ClientChanges struct {
	UserID         uuid.UUID              `json:"userId"`
	LastSync       time.Time              `json:"lastSync"`
	ProfileChanges *ClientProfileChanges  `json:"profileChanges,omitempty"`
	FollowChanges  []ClientFollowChange   `json:"followChanges,omitempty"`
}

// ClientProfileChanges represents profile modifications from client
type ClientProfileChanges struct {
	DisplayName *string                    `json:"displayName,omitempty"`
	Bio         *string                    `json:"bio,omitempty"`
	Interests   []string                   `json:"interests,omitempty"`
	SocialLinks *map[string]interface{}    `json:"socialLinks,omitempty"`
}

// ClientFollowChange represents follow/unfollow actions from client
type ClientFollowChange struct {
	Action       string    `json:"action"`       // "follow" or "unfollow"
	TargetUserID uuid.UUID `json:"targetUserId"`
	Timestamp    time.Time `json:"timestamp"`
}

// SyncResult represents the result of applying client changes
type SyncResult struct {
	UserID         uuid.UUID `json:"userId"`
	ProcessedAt    time.Time `json:"processedAt"`
	Success        bool      `json:"success"`
	AppliedChanges int       `json:"appliedChanges"`
	FailedChanges  int       `json:"failedChanges"`
	Errors         []string  `json:"errors"`
}