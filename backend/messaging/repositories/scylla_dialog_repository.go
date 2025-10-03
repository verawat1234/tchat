package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"

	"tchat.dev/messaging/models"
	"tchat.dev/messaging/services"
)

// ScyllaDialogRepository implements services.DialogRepository using ScyllaDB
type ScyllaDialogRepository struct {
	session *gocql.Session
}

// NewScyllaDialogRepository creates a new ScyllaDB dialog repository
func NewScyllaDialogRepository(session *gocql.Session) services.DialogRepository {
	return &ScyllaDialogRepository{session: session}
}

// Create creates a new dialog in ScyllaDB
func (r *ScyllaDialogRepository) Create(ctx context.Context, dialog *models.Dialog) error {
	// Convert Participants (UUIDSlice) to JSON
	participantsJSON, err := json.Marshal(dialog.Participants)
	if err != nil {
		return fmt.Errorf("failed to marshal participants: %w", err)
	}

	// Insert into dialogs table
	query := `INSERT INTO dialogs (id, type, participant_ids, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`

	if err := r.session.Query(query,
		dialog.ID,
		dialog.Type.String(),
		string(participantsJSON),
		dialog.CreatedAt,
		dialog.UpdatedAt,
	).WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("failed to insert dialog: %w", err)
	}

	// Insert into user_dialogs table for each participant
	userDialogQuery := `INSERT INTO user_dialogs
		(user_id, dialog_id, type, last_message_at, is_archived, is_muted, unread_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, participantID := range dialog.Participants {
		if err := r.session.Query(userDialogQuery,
			participantID,
			dialog.ID,
			dialog.Type.String(),
			dialog.CreatedAt, // Use creation time as initial last_message_at
			false,            // is_archived
			false,            // is_muted
			0,                // unread_count
			dialog.CreatedAt,
			dialog.UpdatedAt,
		).WithContext(ctx).Exec(); err != nil {
			return fmt.Errorf("failed to insert user_dialog for participant %s: %w", participantID, err)
		}
	}

	return nil
}

// GetByID retrieves a dialog by its ID
func (r *ScyllaDialogRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Dialog, error) {
	query := `SELECT id, type, participant_ids, created_at, updated_at, last_message_id
		FROM dialogs WHERE id = ?`

	var (
		dialogID         uuid.UUID
		dialogType       string
		participantsJSON string
		createdAt        time.Time
		updatedAt        time.Time
		lastMessageID    *uuid.UUID
	)

	err := r.session.Query(query, id).WithContext(ctx).Scan(
		&dialogID, &dialogType, &participantsJSON, &createdAt, &updatedAt,
		&lastMessageID,
	)
	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, fmt.Errorf("dialog not found")
		}
		return nil, fmt.Errorf("failed to get dialog: %w", err)
	}

	// Parse participants from JSON
	var participants models.UUIDSlice
	if err := json.Unmarshal([]byte(participantsJSON), &participants); err != nil {
		return nil, fmt.Errorf("failed to parse participants: %w", err)
	}

	dialog := &models.Dialog{
		ID:            dialogID,
		Participants:  participants,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		LastMessageID: lastMessageID,
	}

	// Parse dialog type
	dialog.Type = models.DialogType(dialogType)

	return dialog, nil
}

// GetByUserID retrieves dialogs for a specific user with filters and pagination
func (r *ScyllaDialogRepository) GetByUserID(ctx context.Context, userID uuid.UUID, filters services.DialogFilters, pagination services.Pagination) ([]*models.Dialog, int64, error) {
	// Build query for user_dialogs table
	query := `SELECT dialog_id, type, last_message_at, is_archived, is_muted, unread_count, created_at, updated_at
		FROM user_dialogs WHERE user_id = ?`

	// Note: ScyllaDB doesn't support COUNT(*) efficiently, so we'll fetch all and count
	// For production, implement counter tables or use estimates
	var dialogEntries []struct {
		DialogID       uuid.UUID
		Type           string
		LastMessageAt  time.Time
		IsArchived     bool
		IsMuted        bool
		UnreadCount    int
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}

	iter := r.session.Query(query, userID).WithContext(ctx).Iter()

	var entry struct {
		DialogID       uuid.UUID
		Type           string
		LastMessageAt  time.Time
		IsArchived     bool
		IsMuted        bool
		UnreadCount    int
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}

	for iter.Scan(&entry.DialogID, &entry.Type, &entry.LastMessageAt, &entry.IsArchived,
		&entry.IsMuted, &entry.UnreadCount, &entry.CreatedAt, &entry.UpdatedAt) {

		// Apply filters
		if filters.Type != nil && entry.Type != filters.Type.String() {
			continue
		}
		if filters.IsArchived != nil && entry.IsArchived != *filters.IsArchived {
			continue
		}
		if filters.IsMuted != nil && entry.IsMuted != *filters.IsMuted {
			continue
		}

		dialogEntries = append(dialogEntries, entry)
	}

	if err := iter.Close(); err != nil {
		return nil, 0, fmt.Errorf("failed to query user_dialogs: %w", err)
	}

	total := int64(len(dialogEntries))

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.PageSize
	end := offset + pagination.PageSize
	if end > len(dialogEntries) {
		end = len(dialogEntries)
	}
	if offset > len(dialogEntries) {
		return []*models.Dialog{}, total, nil
	}

	paginatedEntries := dialogEntries[offset:end]

	// Fetch full dialog details for paginated results
	dialogs := make([]*models.Dialog, 0, len(paginatedEntries))
	for _, entry := range paginatedEntries {
		dialog, err := r.GetByID(ctx, entry.DialogID)
		if err != nil {
			continue // Skip dialogs that couldn't be fetched
		}
		dialogs = append(dialogs, dialog)
	}

	return dialogs, total, nil
}

// Update updates a dialog
func (r *ScyllaDialogRepository) Update(ctx context.Context, dialog *models.Dialog) error {
	query := `UPDATE dialogs SET type = ?, updated_at = ?, last_message_id = ?
		WHERE id = ?`

	return r.session.Query(query,
		dialog.Type.String(),
		time.Now(),
		dialog.LastMessageID,
		dialog.ID,
	).WithContext(ctx).Exec()
}

// Delete deletes a dialog
func (r *ScyllaDialogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM dialogs WHERE id = ?`
	return r.session.Query(query, id).WithContext(ctx).Exec()
}

// GetParticipants retrieves all participants in a dialog
func (r *ScyllaDialogRepository) GetParticipants(ctx context.Context, dialogID uuid.UUID) ([]*models.DialogParticipant, error) {
	// Would need a separate dialog_participants table
	return nil, fmt.Errorf("GetParticipants requires a dialog_participants table")
}

// AddParticipant adds a participant to a dialog
func (r *ScyllaDialogRepository) AddParticipant(ctx context.Context, participant *models.DialogParticipant) error {
	// Would need a separate dialog_participants table
	return fmt.Errorf("AddParticipant requires a dialog_participants table")
}

// RemoveParticipant removes a participant from a dialog
func (r *ScyllaDialogRepository) RemoveParticipant(ctx context.Context, dialogID, userID uuid.UUID) error {
	// Would need a separate dialog_participants table
	return fmt.Errorf("RemoveParticipant requires a dialog_participants table")
}

// UpdateParticipant updates a dialog participant
func (r *ScyllaDialogRepository) UpdateParticipant(ctx context.Context, participant *models.DialogParticipant) error {
	// Would need a separate dialog_participants table
	return fmt.Errorf("UpdateParticipant requires a dialog_participants table")
}

// GetAdmins retrieves all admin participants in a dialog
func (r *ScyllaDialogRepository) GetAdmins(ctx context.Context, dialogID uuid.UUID) ([]*models.DialogParticipant, error) {
	// Would need a separate dialog_participants table
	return nil, fmt.Errorf("GetAdmins requires a dialog_participants table")
}

// SearchDialogs searches for dialogs by query
func (r *ScyllaDialogRepository) SearchDialogs(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.Dialog, error) {
	// Full-text search not natively supported
	return nil, fmt.Errorf("search not implemented - requires external search engine")
}
