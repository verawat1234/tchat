package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/messaging/models"
	"tchat.dev/messaging/services"
)

// DialogRepository implements services.DialogRepository using GORM
type DialogRepository struct {
	db *gorm.DB
}

// NewDialogRepository creates a new dialog repository
func NewDialogRepository(db *gorm.DB) services.DialogRepository {
	return &DialogRepository{db: db}
}

// Create creates a new dialog in the database
func (r *DialogRepository) Create(ctx context.Context, dialog *models.Dialog) error {
	return r.db.WithContext(ctx).Create(dialog).Error
}

// GetByID retrieves a dialog by its ID
func (r *DialogRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Dialog, error) {
	var dialog models.Dialog
	err := r.db.WithContext(ctx).Preload("Participants").First(&dialog, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &dialog, nil
}

// GetByUserID retrieves dialogs for a specific user with filters and pagination
func (r *DialogRepository) GetByUserID(ctx context.Context, userID uuid.UUID, filters services.DialogFilters, pagination services.Pagination) ([]*models.Dialog, int64, error) {
	var dialogs []*models.Dialog
	var total int64

	// Base query to find dialogs where user is a participant
	query := r.db.WithContext(ctx).Model(&models.Dialog{}).
		Joins("JOIN dialog_participants ON dialogs.id = dialog_participants.dialog_id").
		Where("dialog_participants.user_id = ?", userID)

	// Apply filters
	if filters.Type != nil {
		query = query.Where("dialogs.type = ?", *filters.Type)
	}
	if filters.IsArchived != nil {
		if *filters.IsArchived {
			query = query.Where("dialog_participants.is_archived = true")
		} else {
			query = query.Where("dialog_participants.is_archived = false OR dialog_participants.is_archived IS NULL")
		}
	}
	if filters.IsMuted != nil {
		if *filters.IsMuted {
			query = query.Where("dialog_participants.is_muted = true")
		} else {
			query = query.Where("dialog_participants.is_muted = false OR dialog_participants.is_muted IS NULL")
		}
	}
	if filters.UpdatedFrom != nil {
		query = query.Where("dialogs.updated_at >= ?", *filters.UpdatedFrom)
	}
	if filters.UpdatedTo != nil {
		query = query.Where("dialogs.updated_at <= ?", *filters.UpdatedTo)
	}

	// Count total
	if err := query.Select("DISTINCT dialogs.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	orderBy := "updated_at"
	if pagination.OrderBy != "" {
		orderBy = pagination.OrderBy
	}
	order := "DESC"
	if pagination.Order != "" {
		order = pagination.Order
	}

	offset := (pagination.Page - 1) * pagination.PageSize
	err := query.Select("DISTINCT dialogs.*").
		Preload("Participants").
		Order(fmt.Sprintf("dialogs.%s %s", orderBy, order)).
		Limit(pagination.PageSize).
		Offset(offset).
		Find(&dialogs).Error

	return dialogs, total, err
}

// Update updates an existing dialog
func (r *DialogRepository) Update(ctx context.Context, dialog *models.Dialog) error {
	return r.db.WithContext(ctx).Save(dialog).Error
}

// Delete soft deletes a dialog and its participants
func (r *DialogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete dialog participants
		if err := tx.Where("dialog_id = ?", id).Delete(&models.DialogParticipant{}).Error; err != nil {
			return err
		}

		// Delete dialog
		return tx.Delete(&models.Dialog{}, "id = ?", id).Error
	})
}

// GetParticipants retrieves all participants of a dialog
func (r *DialogRepository) GetParticipants(ctx context.Context, dialogID uuid.UUID) ([]*models.DialogParticipant, error) {
	var participants []*models.DialogParticipant
	err := r.db.WithContext(ctx).
		Where("dialog_id = ?", dialogID).
		Find(&participants).Error
	return participants, err
}

// AddParticipant adds a new participant to a dialog
func (r *DialogRepository) AddParticipant(ctx context.Context, participant *models.DialogParticipant) error {
	return r.db.WithContext(ctx).Create(participant).Error
}

// RemoveParticipant removes a participant from a dialog
func (r *DialogRepository) RemoveParticipant(ctx context.Context, dialogID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("dialog_id = ? AND user_id = ?", dialogID, userID).
		Delete(&models.DialogParticipant{}).Error
}

// UpdateParticipant updates a dialog participant's information
func (r *DialogRepository) UpdateParticipant(ctx context.Context, participant *models.DialogParticipant) error {
	return r.db.WithContext(ctx).Save(participant).Error
}

// GetAdmins retrieves admin participants of a dialog
func (r *DialogRepository) GetAdmins(ctx context.Context, dialogID uuid.UUID) ([]*models.DialogParticipant, error) {
	var admins []*models.DialogParticipant
	err := r.db.WithContext(ctx).
		Where("dialog_id = ? AND role IN ?", dialogID, []string{"admin", "owner"}).
		Find(&admins).Error
	return admins, err
}

// SearchDialogs searches for dialogs by name for a specific user
func (r *DialogRepository) SearchDialogs(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.Dialog, error) {
	var dialogs []*models.Dialog
	err := r.db.WithContext(ctx).
		Joins("JOIN dialog_participants ON dialogs.id = dialog_participants.dialog_id").
		Where("dialog_participants.user_id = ? AND dialogs.name ILIKE ?", userID, "%"+query+"%").
		Preload("Participants").
		Order("dialogs.updated_at DESC").
		Limit(limit).
		Find(&dialogs).Error
	return dialogs, err
}