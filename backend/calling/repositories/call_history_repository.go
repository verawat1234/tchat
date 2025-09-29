package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/calling/models"
)

// CallHistoryRepository defines the interface for call history data access
type CallHistoryRepository interface {
	Create(history *models.CallHistory) error
	GetByID(id uuid.UUID) (*models.CallHistory, error)
	Update(history *models.CallHistory) error
	Delete(id uuid.UUID) error
	GetByUserID(userID uuid.UUID, filter models.CallHistoryFilter) ([]models.CallHistory, error)
	GetByCallSessionID(callSessionID uuid.UUID) ([]models.CallHistory, error)
	GetCallHistory(userID uuid.UUID, limit, offset int) ([]models.CallHistory, error)
	GetCallHistoryByType(userID uuid.UUID, callType models.CallType, limit, offset int) ([]models.CallHistory, error)
	GetCallStats(userID uuid.UUID) (map[string]interface{}, error)
	GetRecentCallHistory(userID uuid.UUID, days int, limit int) ([]models.CallHistory, error)
	DeleteUserHistory(userID uuid.UUID) error
}

// GormCallHistoryRepository implements CallHistoryRepository using GORM
type GormCallHistoryRepository struct {
	db *gorm.DB
}

// NewGormCallHistoryRepository creates a new GORM-based call history repository
func NewGormCallHistoryRepository(db *gorm.DB) CallHistoryRepository {
	return &GormCallHistoryRepository{db: db}
}

// Create creates a new call history record
func (r *GormCallHistoryRepository) Create(history *models.CallHistory) error {
	return r.db.Create(history).Error
}

// GetByID retrieves a call history record by its ID
func (r *GormCallHistoryRepository) GetByID(id uuid.UUID) (*models.CallHistory, error) {
	var history models.CallHistory
	err := r.db.First(&history, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrCallNotFound
		}
		return nil, err
	}
	return &history, nil
}

// Update updates an existing call history record
func (r *GormCallHistoryRepository) Update(history *models.CallHistory) error {
	return r.db.Save(history).Error
}

// Delete soft deletes a call history record
func (r *GormCallHistoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.CallHistory{}, "id = ?", id).Error
}

// GetByUserID retrieves call history for a specific user with filters
func (r *GormCallHistoryRepository) GetByUserID(userID uuid.UUID, filter models.CallHistoryFilter) ([]models.CallHistory, error) {
	var histories []models.CallHistory

	query := r.db.Where("user_id = ?", userID)

	// Apply filters
	if filter.CallType != nil {
		query = query.Where("call_type = ?", *filter.CallType)
	}

	if filter.CallStatus != nil {
		query = query.Where("call_status = ?", *filter.CallStatus)
	}

	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}

	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	// Apply ordering
	query = query.Order("created_at DESC")

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	err := query.Find(&histories).Error
	return histories, err
}

// GetByCallSessionID retrieves all call history records for a specific call session
func (r *GormCallHistoryRepository) GetByCallSessionID(callSessionID uuid.UUID) ([]models.CallHistory, error) {
	var histories []models.CallHistory
	err := r.db.Where("call_session_id = ?", callSessionID).Find(&histories).Error
	return histories, err
}

// GetCallHistory retrieves paginated call history for a user
func (r *GormCallHistoryRepository) GetCallHistory(userID uuid.UUID, limit, offset int) ([]models.CallHistory, error) {
	var histories []models.CallHistory

	query := r.db.Where("user_id = ?", userID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&histories).Error
	return histories, err
}

// GetCallHistoryByType retrieves call history filtered by call type
func (r *GormCallHistoryRepository) GetCallHistoryByType(userID uuid.UUID, callType models.CallType, limit, offset int) ([]models.CallHistory, error) {
	var histories []models.CallHistory

	query := r.db.Where("user_id = ? AND call_type = ?", userID, callType).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&histories).Error
	return histories, err
}

// GetCallStats retrieves call statistics for a user
func (r *GormCallHistoryRepository) GetCallStats(userID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total calls
	var totalCalls int64
	if err := r.db.Model(&models.CallHistory{}).Where("user_id = ?", userID).Count(&totalCalls).Error; err != nil {
		return nil, err
	}
	stats["total_calls"] = totalCalls

	// Outgoing calls
	var outgoingCalls int64
	if err := r.db.Model(&models.CallHistory{}).Where("user_id = ? AND initiated_by_me = ?", userID, true).Count(&outgoingCalls).Error; err != nil {
		return nil, err
	}
	stats["outgoing_calls"] = outgoingCalls

	// Incoming calls
	var incomingCalls int64
	if err := r.db.Model(&models.CallHistory{}).Where("user_id = ? AND initiated_by_me = ?", userID, false).Count(&incomingCalls).Error; err != nil {
		return nil, err
	}
	stats["incoming_calls"] = incomingCalls

	// Successful calls (completed)
	var successfulCalls int64
	if err := r.db.Model(&models.CallHistory{}).
		Where("user_id = ? AND call_status = ? AND duration > 0", userID, models.CallStatusEnded).
		Count(&successfulCalls).Error; err != nil {
		return nil, err
	}
	stats["successful_calls"] = successfulCalls

	// Missed calls (incoming failed calls with 0 duration)
	var missedCalls int64
	if err := r.db.Model(&models.CallHistory{}).
		Where("user_id = ? AND initiated_by_me = ? AND call_status = ? AND duration = 0",
			userID, false, models.CallStatusFailed).
		Count(&missedCalls).Error; err != nil {
		return nil, err
	}
	stats["missed_calls"] = missedCalls

	// Voice calls
	var voiceCalls int64
	if err := r.db.Model(&models.CallHistory{}).
		Where("user_id = ? AND call_type = ?", userID, models.CallTypeVoice).
		Count(&voiceCalls).Error; err != nil {
		return nil, err
	}
	stats["voice_calls"] = voiceCalls

	// Video calls
	var videoCalls int64
	if err := r.db.Model(&models.CallHistory{}).
		Where("user_id = ? AND call_type = ?", userID, models.CallTypeVideo).
		Count(&videoCalls).Error; err != nil {
		return nil, err
	}
	stats["video_calls"] = videoCalls

	// Total talk time (in seconds)
	var totalDuration int64
	if err := r.db.Model(&models.CallHistory{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(duration), 0)").
		Scan(&totalDuration).Error; err != nil {
		return nil, err
	}
	stats["total_duration"] = totalDuration

	// Average call duration for successful calls
	var avgDuration float64
	if err := r.db.Model(&models.CallHistory{}).
		Where("user_id = ? AND duration > 0", userID).
		Select("COALESCE(AVG(duration), 0)").
		Scan(&avgDuration).Error; err != nil {
		return nil, err
	}
	stats["average_duration"] = avgDuration

	return stats, nil
}

// GetRecentCallHistory retrieves recent call history with participant info
func (r *GormCallHistoryRepository) GetRecentCallHistory(userID uuid.UUID, days int, limit int) ([]models.CallHistory, error) {
	var histories []models.CallHistory

	query := r.db.Where("user_id = ?", userID)

	if days > 0 {
		query = query.Where("created_at >= NOW() - INTERVAL ? DAY", days)
	}

	query = query.Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&histories).Error
	return histories, err
}

// GetCallFrequency retrieves call frequency data for analytics
func (r *GormCallHistoryRepository) GetCallFrequency(userID uuid.UUID, days int) (map[string]int64, error) {
	frequency := make(map[string]int64)

	// Get calls per day for the last N days
	rows, err := r.db.Raw(`
		SELECT DATE(created_at) as call_date, COUNT(*) as call_count
		FROM call_histories
		WHERE user_id = ? AND created_at >= NOW() - INTERVAL ? DAY
		GROUP BY DATE(created_at)
		ORDER BY call_date DESC
	`, userID, days).Rows()

	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Log error but don't return as this is in defer
			_ = err // Avoid unused variable warning
		}
	}()

	for rows.Next() {
		var date string
		var count int64
		if err := rows.Scan(&date, &count); err == nil {
			frequency[date] = count
		}
	}

	return frequency, nil
}

// DeleteUserHistory deletes all call history for a user (GDPR compliance)
func (r *GormCallHistoryRepository) DeleteUserHistory(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.CallHistory{}).Error
}

// BulkCreateHistory creates multiple call history records efficiently
func (r *GormCallHistoryRepository) BulkCreateHistory(histories []models.CallHistory) error {
	if len(histories) == 0 {
		return nil
	}

	return r.db.CreateInBatches(histories, 100).Error
}

// GetCallPartners retrieves frequent call partners for a user
func (r *GormCallHistoryRepository) GetCallPartners(userID uuid.UUID, limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	rows, err := r.db.Raw(`
		SELECT
			other_participant_id,
			other_participant_name,
			COUNT(*) as call_count,
			SUM(duration) as total_duration,
			MAX(created_at) as last_call
		FROM call_histories
		WHERE user_id = ?
		GROUP BY other_participant_id, other_participant_name
		ORDER BY call_count DESC, last_call DESC
		LIMIT ?
	`, userID, limit).Rows()

	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Log error but don't return as this is in defer
			_ = err // Avoid unused variable warning
		}
	}()

	for rows.Next() {
		var partnerID string
		var partnerName *string
		var callCount int64
		var totalDuration int64
		var lastCall string

		if err := rows.Scan(&partnerID, &partnerName, &callCount, &totalDuration, &lastCall); err == nil {
			partner := map[string]interface{}{
				"user_id":        partnerID,
				"name":           partnerName,
				"call_count":     callCount,
				"total_duration": totalDuration,
				"last_call":      lastCall,
			}
			results = append(results, partner)
		}
	}

	return results, nil
}
