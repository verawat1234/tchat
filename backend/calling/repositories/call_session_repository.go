package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/calling/models"
)

// CallSessionRepository defines the interface for call session data access
type CallSessionRepository interface {
	Create(callSession *models.CallSession) error
	GetByID(id uuid.UUID) (*models.CallSession, error)
	Update(callSession *models.CallSession) error
	Delete(id uuid.UUID) error
	GetActiveCallByUserID(userID uuid.UUID) (*models.CallSession, error)
	GetByStatus(status models.CallStatus) ([]models.CallSession, error)
	GetRecentCalls(userID uuid.UUID, limit int) ([]models.CallSession, error)
}

// GormCallSessionRepository implements CallSessionRepository using GORM
type GormCallSessionRepository struct {
	db *gorm.DB
}

// NewGormCallSessionRepository creates a new GORM-based call session repository
func NewGormCallSessionRepository(db *gorm.DB) CallSessionRepository {
	return &GormCallSessionRepository{db: db}
}

// Create creates a new call session in the database
func (r *GormCallSessionRepository) Create(callSession *models.CallSession) error {
	return r.db.Create(callSession).Error
}

// GetByID retrieves a call session by its ID, including participants
func (r *GormCallSessionRepository) GetByID(id uuid.UUID) (*models.CallSession, error) {
	var callSession models.CallSession
	err := r.db.Preload("Participants").First(&callSession, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrCallNotFound
		}
		return nil, err
	}
	return &callSession, nil
}

// Update updates an existing call session in the database
func (r *GormCallSessionRepository) Update(callSession *models.CallSession) error {
	// Update the main call session
	if err := r.db.Save(callSession).Error; err != nil {
		return err
	}

	// Update participants if they exist
	for _, participant := range callSession.Participants {
		if err := r.db.Save(&participant).Error; err != nil {
			return err
		}
	}

	return nil
}

// Delete soft deletes a call session
func (r *GormCallSessionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.CallSession{}, "id = ?", id).Error
}

// GetActiveCallByUserID finds an active call for a specific user
func (r *GormCallSessionRepository) GetActiveCallByUserID(userID uuid.UUID) (*models.CallSession, error) {
	var callSession models.CallSession

	err := r.db.Preload("Participants").
		Joins("JOIN call_participants ON call_sessions.id = call_participants.call_session_id").
		Where("call_participants.user_id = ?", userID).
		Where("call_sessions.status IN ?", []models.CallStatus{
			models.CallStatusConnecting,
			models.CallStatusActive,
		}).
		Where("call_participants.left_at IS NULL").
		First(&callSession).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No active call found
		}
		return nil, err
	}

	return &callSession, nil
}

// GetByStatus retrieves all call sessions with a specific status
func (r *GormCallSessionRepository) GetByStatus(status models.CallStatus) ([]models.CallSession, error) {
	var callSessions []models.CallSession
	err := r.db.Preload("Participants").
		Where("status = ?", status).
		Find(&callSessions).Error
	return callSessions, err
}

// GetRecentCalls retrieves recent calls for a user
func (r *GormCallSessionRepository) GetRecentCalls(userID uuid.UUID, limit int) ([]models.CallSession, error) {
	var callSessions []models.CallSession

	err := r.db.Preload("Participants").
		Joins("JOIN call_participants ON call_sessions.id = call_participants.call_session_id").
		Where("call_participants.user_id = ?", userID).
		Order("call_sessions.started_at DESC").
		Limit(limit).
		Find(&callSessions).Error

	return callSessions, err
}

// GetCallsInTimeRange retrieves calls within a specific time range
func (r *GormCallSessionRepository) GetCallsInTimeRange(userID uuid.UUID, startTime, endTime interface{}) ([]models.CallSession, error) {
	var callSessions []models.CallSession

	query := r.db.Preload("Participants").
		Joins("JOIN call_participants ON call_sessions.id = call_participants.call_session_id").
		Where("call_participants.user_id = ?", userID).
		Where("call_sessions.started_at BETWEEN ? AND ?", startTime, endTime)

	err := query.Find(&callSessions).Error
	return callSessions, err
}

// GetCallsByType retrieves calls of a specific type for a user
func (r *GormCallSessionRepository) GetCallsByType(userID uuid.UUID, callType models.CallType, limit int) ([]models.CallSession, error) {
	var callSessions []models.CallSession

	err := r.db.Preload("Participants").
		Joins("JOIN call_participants ON call_sessions.id = call_participants.call_session_id").
		Where("call_participants.user_id = ?", userID).
		Where("call_sessions.type = ?", callType).
		Order("call_sessions.started_at DESC").
		Limit(limit).
		Find(&callSessions).Error

	return callSessions, err
}

// CountActiveCalls returns the number of currently active calls
func (r *GormCallSessionRepository) CountActiveCalls() (int64, error) {
	var count int64
	err := r.db.Model(&models.CallSession{}).
		Where("status IN ?", []models.CallStatus{
			models.CallStatusConnecting,
			models.CallStatusActive,
		}).
		Count(&count).Error
	return count, err
}

// GetParticipantsByCallID retrieves all participants for a specific call
func (r *GormCallSessionRepository) GetParticipantsByCallID(callID uuid.UUID) ([]models.CallParticipant, error) {
	var participants []models.CallParticipant
	err := r.db.Where("call_session_id = ?", callID).Find(&participants).Error
	return participants, err
}
