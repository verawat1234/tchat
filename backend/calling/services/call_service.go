package services

import (
	"time"

	"github.com/google/uuid"
	"tchat.dev/calling/models"
	"tchat.dev/calling/repositories"
)

// CallService handles call management business logic
type CallService struct {
	callRepo     repositories.CallSessionRepository
	presenceRepo repositories.UserPresenceRepository
	historyRepo  repositories.CallHistoryRepository
}

// NewCallService creates a new CallService instance
func NewCallService(
	callRepo repositories.CallSessionRepository,
	presenceRepo repositories.UserPresenceRepository,
	historyRepo repositories.CallHistoryRepository,
) *CallService {
	return &CallService{
		callRepo:     callRepo,
		presenceRepo: presenceRepo,
		historyRepo:  historyRepo,
	}
}

// InitiateCallRequest represents a request to initiate a call
type InitiateCallRequest struct {
	CallerID uuid.UUID       `json:"caller_id" validate:"required"`
	CalleeID uuid.UUID       `json:"callee_id" validate:"required"`
	CallType models.CallType `json:"call_type" validate:"required,oneof=voice video"`
}

// AnswerCallRequest represents a request to answer a call
type AnswerCallRequest struct {
	CallID uuid.UUID `json:"call_id" validate:"required"`
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Accept bool      `json:"accept"`
}

// EndCallRequest represents a request to end a call
type EndCallRequest struct {
	CallID uuid.UUID `json:"call_id" validate:"required"`
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// InitiateCall creates a new call session
func (s *CallService) InitiateCall(req InitiateCallRequest) (*models.CallSession, error) {
	// Validate caller and callee are different
	if req.CallerID == req.CalleeID {
		return nil, models.ErrInvalidCallState
	}

	// Check if caller is available (not already in a call)
	callerPresence, err := s.presenceRepo.GetByUserID(req.CallerID)
	if err != nil {
		return nil, err
	}
	if callerPresence.IsInCall() {
		return nil, models.ErrUserAlreadyInCall
	}

	// Check if callee is available
	calleePresence, err := s.presenceRepo.GetByUserID(req.CalleeID)
	if err != nil {
		return nil, models.ErrUserNotFound
	}
	if !calleePresence.IsAvailable() {
		return nil, models.ErrUserNotAvailable
	}

	// Create call session
	callSession := &models.CallSession{
		Type:        req.CallType,
		Status:      models.CallStatusConnecting,
		InitiatedBy: req.CallerID,
		StartedAt:   time.Now(),
	}

	// Save call session
	if err := s.callRepo.Create(callSession); err != nil {
		return nil, err
	}

	// Create participants
	callerParticipant := &models.CallParticipant{
		CallSessionID:     callSession.ID,
		UserID:            req.CallerID,
		Role:              models.ParticipantRoleCaller,
		JoinedAt:          time.Now(),
		AudioEnabled:      true,
		VideoEnabled:      req.CallType == models.CallTypeVideo,
		ConnectionQuality: models.ConnectionQualityGood,
	}

	calleeParticipant := &models.CallParticipant{
		CallSessionID:     callSession.ID,
		UserID:            req.CalleeID,
		Role:              models.ParticipantRoleCallee,
		JoinedAt:          time.Now(),
		AudioEnabled:      true,
		VideoEnabled:      req.CallType == models.CallTypeVideo,
		ConnectionQuality: models.ConnectionQualityGood,
	}

	// Add participants to call session
	callSession.Participants = []models.CallParticipant{*callerParticipant, *calleeParticipant}

	// Update call session with participants
	if err := s.callRepo.Update(callSession); err != nil {
		return nil, err
	}

	// Update presence status for both users
	callerPresence.SetInCall()
	calleePresence.SetInCall()
	if err := s.presenceRepo.Update(callerPresence); err != nil {
		// Log error but continue - presence update failure shouldn't break call initiation
		_ = err // Avoid unused variable warning
	}
	if err := s.presenceRepo.Update(calleePresence); err != nil {
		// Log error but continue - presence update failure shouldn't break call initiation
		_ = err // Avoid unused variable warning
	}

	// Create initial history records
	historyRecords := models.CreateFromCallSession(callSession)
	for _, record := range historyRecords {
		if err := s.historyRepo.Create(record); err != nil {
			// Log error but continue - history creation failure shouldn't break call initiation
			_ = err // Avoid unused variable warning
		}
	}

	return callSession, nil
}

// AnswerCall handles answering or declining a call
func (s *CallService) AnswerCall(req AnswerCallRequest) (*models.CallSession, error) {
	// Get call session
	callSession, err := s.callRepo.GetByID(req.CallID)
	if err != nil {
		return nil, models.ErrCallNotFound
	}

	// Validate call state
	if callSession.Status != models.CallStatusConnecting {
		return nil, models.ErrCallAlreadyAnswered
	}

	// Validate user is a participant (callee)
	participant := callSession.GetParticipantByUserID(req.UserID)
	if participant == nil {
		return nil, models.ErrParticipantNotFound
	}
	if participant.Role != models.ParticipantRoleCallee {
		return nil, models.ErrPermissionDenied
	}

	if req.Accept {
		// Accept the call
		callSession.Status = models.CallStatusActive
	} else {
		// Decline the call
		callSession.Status = models.CallStatusFailed
		reason := "declined"
		callSession.FailureReason = &reason

		// Update presence for both participants back to online
		for _, p := range callSession.Participants {
			presence, _ := s.presenceRepo.GetByUserID(p.UserID)
			if presence != nil {
				presence.SetOnline()
				if err := s.presenceRepo.Update(presence); err != nil {
					// Log error but continue
					_ = err // Avoid unused variable warning
				}
			}
		}
	}

	// Update call session
	if err := s.callRepo.Update(callSession); err != nil {
		return nil, err
	}

	// Update history records
	historyRecords, _ := s.historyRepo.GetByCallSessionID(callSession.ID)
	for _, record := range historyRecords {
		record.UpdateFromCallSession(callSession)
		if err := s.historyRepo.Update(&record); err != nil {
			// Log error but continue
			_ = err // Avoid unused variable warning
		}
	}

	return callSession, nil
}

// EndCall terminates an active call
func (s *CallService) EndCall(req EndCallRequest) (*models.CallSession, error) {
	// Get call session
	callSession, err := s.callRepo.GetByID(req.CallID)
	if err != nil {
		return nil, models.ErrCallNotFound
	}

	// Validate call state
	if callSession.IsTerminal() {
		return nil, models.ErrCallAlreadyEnded
	}

	// Validate user is a participant
	if !callSession.HasParticipant(req.UserID) {
		return nil, models.ErrParticipantNotFound
	}

	// End the call
	now := time.Now()
	callSession.Status = models.CallStatusEnded
	callSession.EndedAt = &now
	callSession.CalculateDuration()

	// Mark all participants as having left
	for i := range callSession.Participants {
		if callSession.Participants[i].IsActive() {
			callSession.Participants[i].Leave()
		}
	}

	// Update call session
	if err := s.callRepo.Update(callSession); err != nil {
		return nil, err
	}

	// Update presence for all participants back to online
	for _, participant := range callSession.Participants {
		presence, _ := s.presenceRepo.GetByUserID(participant.UserID)
		if presence != nil {
			presence.SetOnline()
			if err := s.presenceRepo.Update(presence); err != nil {
				// Log error but continue
				_ = err // Avoid unused variable warning
			}
		}
	}

	// Update history records
	historyRecords, _ := s.historyRepo.GetByCallSessionID(callSession.ID)
	for _, record := range historyRecords {
		record.UpdateFromCallSession(callSession)
		if err := s.historyRepo.Update(&record); err != nil {
			// Log error but continue
			_ = err // Avoid unused variable warning
		}
	}

	return callSession, nil
}

// GetCallStatus retrieves the current status of a call
func (s *CallService) GetCallStatus(callID uuid.UUID, userID uuid.UUID) (*models.CallSession, error) {
	callSession, err := s.callRepo.GetByID(callID)
	if err != nil {
		return nil, models.ErrCallNotFound
	}

	// Validate user is a participant
	if !callSession.HasParticipant(userID) {
		return nil, models.ErrPermissionDenied
	}

	return callSession, nil
}

// GetActiveCallForUser returns the active call for a user, if any
func (s *CallService) GetActiveCallForUser(userID uuid.UUID) (*models.CallSession, error) {
	return s.callRepo.GetActiveCallByUserID(userID)
}

// ToggleMedia toggles audio or video for a participant
func (s *CallService) ToggleMedia(callID uuid.UUID, userID uuid.UUID, mediaType string, enabled bool) error {
	callSession, err := s.callRepo.GetByID(callID)
	if err != nil {
		return models.ErrCallNotFound
	}

	participant := callSession.GetParticipantByUserID(userID)
	if participant == nil {
		return models.ErrParticipantNotFound
	}

	// Only allow media changes in active calls
	if callSession.Status != models.CallStatusActive {
		return models.ErrInvalidCallState
	}

	switch mediaType {
	case "audio":
		participant.AudioEnabled = enabled
	case "video":
		participant.VideoEnabled = enabled
	default:
		return models.ErrInvalidCallState
	}

	return s.callRepo.Update(callSession)
}

// UpdateConnectionQuality updates a participant's connection quality
func (s *CallService) UpdateConnectionQuality(callID uuid.UUID, userID uuid.UUID, quality models.ConnectionQuality) error {
	callSession, err := s.callRepo.GetByID(callID)
	if err != nil {
		return models.ErrCallNotFound
	}

	participant := callSession.GetParticipantByUserID(userID)
	if participant == nil {
		return models.ErrParticipantNotFound
	}

	if err := participant.SetConnectionQuality(quality); err != nil {
		return err
	}

	return s.callRepo.Update(callSession)
}
