package services

import (
	"time"

	"github.com/google/uuid"
	"tchat.dev/calling/models"
	"tchat.dev/calling/repositories"
)

// PresenceService handles user presence and availability management
type PresenceService struct {
	presenceRepo repositories.UserPresenceRepository
}

// NewPresenceService creates a new PresenceService instance
func NewPresenceService(presenceRepo repositories.UserPresenceRepository) *PresenceService {
	return &PresenceService{
		presenceRepo: presenceRepo,
	}
}

// UpdatePresenceRequest represents a request to update user presence
type UpdatePresenceRequest struct {
	UserID uuid.UUID             `json:"user_id" validate:"required"`
	Status models.PresenceStatus `json:"status" validate:"required,oneof=online busy offline"`
}

// GetPresenceResponse represents presence information for a user
type GetPresenceResponse struct {
	UserID    uuid.UUID             `json:"user_id"`
	Status    models.PresenceStatus `json:"status"`
	LastSeen  time.Time             `json:"last_seen"`
	UpdatedAt time.Time             `json:"updated_at"`
}

// CheckAvailabilityResponse represents availability information for a user
type CheckAvailabilityResponse struct {
	UserID    uuid.UUID             `json:"user_id"`
	Available bool                  `json:"available"`
	Status    models.PresenceStatus `json:"status"`
	LastSeen  time.Time             `json:"last_seen"`
}

// GetPresence retrieves the current presence status for a user
func (s *PresenceService) GetPresence(userID uuid.UUID) (*GetPresenceResponse, error) {
	presence, err := s.presenceRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &GetPresenceResponse{
		UserID:    presence.UserID,
		Status:    presence.Status,
		LastSeen:  presence.LastSeen,
		UpdatedAt: presence.UpdatedAt,
	}, nil
}

// UpdatePresence updates a user's presence status
func (s *PresenceService) UpdatePresence(req UpdatePresenceRequest) (*GetPresenceResponse, error) {
	// Get current presence or create new one
	presence, err := s.presenceRepo.GetByUserID(req.UserID)
	if err != nil {
		return nil, err
	}

	// Update the status
	if err := presence.UpdateStatus(req.Status); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.presenceRepo.Update(presence); err != nil {
		return nil, err
	}

	return &GetPresenceResponse{
		UserID:    presence.UserID,
		Status:    presence.Status,
		LastSeen:  presence.LastSeen,
		UpdatedAt: presence.UpdatedAt,
	}, nil
}

// CheckAvailability checks if a user is available for calls
func (s *PresenceService) CheckAvailability(userID uuid.UUID) (*CheckAvailabilityResponse, error) {
	presence, err := s.presenceRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &CheckAvailabilityResponse{
		UserID:    presence.UserID,
		Available: presence.IsAvailable(),
		Status:    presence.Status,
		LastSeen:  presence.LastSeen,
	}, nil
}

// SetUserOnline sets a user as online and available
func (s *PresenceService) SetUserOnline(userID uuid.UUID) error {
	return s.presenceRepo.SetOnline(userID)
}

// SetUserOffline sets a user as offline
func (s *PresenceService) SetUserOffline(userID uuid.UUID) error {
	return s.presenceRepo.SetOffline(userID)
}

// SetUserBusy sets a user as busy (not available for calls)
func (s *PresenceService) SetUserBusy(userID uuid.UUID) error {
	return s.presenceRepo.SetBusy(userID)
}

// SetUserInCall sets a user as currently in a call
func (s *PresenceService) SetUserInCall(userID uuid.UUID) error {
	return s.presenceRepo.SetInCall(userID)
}

// GetOnlineUsers retrieves all currently online users
func (s *PresenceService) GetOnlineUsers() ([]GetPresenceResponse, error) {
	presences, err := s.presenceRepo.GetOnlineUsers()
	if err != nil {
		return nil, err
	}

	responses := make([]GetPresenceResponse, len(presences))
	for i, presence := range presences {
		responses[i] = GetPresenceResponse{
			UserID:    presence.UserID,
			Status:    presence.Status,
			LastSeen:  presence.LastSeen,
			UpdatedAt: presence.UpdatedAt,
		}
	}

	return responses, nil
}

// GetUsersByStatus retrieves all users with a specific presence status
func (s *PresenceService) GetUsersByStatus(status models.PresenceStatus) ([]GetPresenceResponse, error) {
	presences, err := s.presenceRepo.GetUsersByStatus(status)
	if err != nil {
		return nil, err
	}

	responses := make([]GetPresenceResponse, len(presences))
	for i, presence := range presences {
		responses[i] = GetPresenceResponse{
			UserID:    presence.UserID,
			Status:    presence.Status,
			LastSeen:  presence.LastSeen,
			UpdatedAt: presence.UpdatedAt,
		}
	}

	return responses, nil
}

// GetBulkPresence retrieves presence for multiple users efficiently
func (s *PresenceService) GetBulkPresence(userIDs []uuid.UUID) (map[uuid.UUID]*GetPresenceResponse, error) {
	presences, err := s.presenceRepo.GetBulkPresence(userIDs)
	if err != nil {
		return nil, err
	}

	responses := make(map[uuid.UUID]*GetPresenceResponse)
	for userID, presence := range presences {
		if presence != nil {
			responses[userID] = &GetPresenceResponse{
				UserID:    presence.UserID,
				Status:    presence.Status,
				LastSeen:  presence.LastSeen,
				UpdatedAt: presence.UpdatedAt,
			}
		}
	}

	return responses, nil
}

// IsUserOnline checks if a specific user is currently online
func (s *PresenceService) IsUserOnline(userID uuid.UUID) (bool, error) {
	return s.presenceRepo.IsUserOnline(userID)
}

// GetPresenceStats returns presence statistics
func (s *PresenceService) GetPresenceStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get count for each status
	statuses := []models.PresenceStatus{
		models.PresenceStatusOnline,
		models.PresenceStatusBusy,
		models.PresenceStatusInCall,
		models.PresenceStatusOffline,
	}

	for _, status := range statuses {
		users, err := s.presenceRepo.GetUsersByStatus(status)
		if err != nil {
			continue // Skip errored statuses
		}
		stats[string(status)] = len(users)
	}

	stats["total_tracked_users"] = stats["online"].(int) +
		stats["busy"].(int) +
		stats["in_call"].(int) +
		stats["offline"].(int)

	return stats, nil
}

// HandleUserDisconnect handles cleanup when a user disconnects
func (s *PresenceService) HandleUserDisconnect(userID uuid.UUID) error {
	// Set user offline when they disconnect
	return s.presenceRepo.SetOffline(userID)
}

// HandleUserConnect handles user connection events
func (s *PresenceService) HandleUserConnect(userID uuid.UUID) error {
	// Set user online when they connect
	return s.presenceRepo.SetOnline(userID)
}

// CleanupExpiredPresence removes expired presence entries
func (s *PresenceService) CleanupExpiredPresence() error {
	// Get reference to the Redis repository for cleanup
	if redisRepo, ok := s.presenceRepo.(*repositories.RedisUserPresenceRepository); ok {
		return redisRepo.CleanupExpiredPresence()
	}
	return nil // No cleanup needed for other repository types
}

// GetRecentlyActiveUsers returns users who were active within the specified duration
func (s *PresenceService) GetRecentlyActiveUsers(duration time.Duration) ([]GetPresenceResponse, error) {
	// Get all non-offline users first
	statuses := []models.PresenceStatus{
		models.PresenceStatusOnline,
		models.PresenceStatusBusy,
		models.PresenceStatusInCall,
	}

	var allActiveUsers []GetPresenceResponse

	for _, status := range statuses {
		users, err := s.GetUsersByStatus(status)
		if err != nil {
			continue
		}

		// Filter by recent activity
		for _, user := range users {
			if time.Since(user.LastSeen) <= duration {
				allActiveUsers = append(allActiveUsers, user)
			}
		}
	}

	return allActiveUsers, nil
}
