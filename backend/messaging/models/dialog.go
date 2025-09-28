package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Dialog represents a conversation between users (direct message, group, channel, or business)
type Dialog struct {
	ID               uuid.UUID   `json:"id" db:"id"`
	Type             DialogType  `json:"type" db:"type"`
	Name             *string     `json:"name,omitempty" db:"name"`
	Title            *string     `json:"title,omitempty" db:"title"` // Alias for name for compatibility
	Avatar           *string     `json:"avatar,omitempty" db:"avatar"`
	Description      *string     `json:"description,omitempty" db:"description"`
	Participants     UUIDSlice   `json:"participants" gorm:"type:json"`
	ParticipantCount int         `json:"participant_count" db:"participant_count"`
	AdminIDs         UUIDSlice   `json:"admin_ids,omitempty" gorm:"type:json"`
	LastMessageID    *uuid.UUID  `json:"last_message_id,omitempty" db:"last_message_id"`
	UnreadCount      int         `json:"unread_count" db:"unread_count"`
	IsPinned         bool        `json:"is_pinned" db:"is_pinned"`
	IsArchived       bool        `json:"is_archived" db:"is_archived"`
	IsMuted          bool        `json:"is_muted" db:"is_muted"`
	Settings         DialogSettings `json:"settings" gorm:"type:json"`
	Metadata         JSON `json:"metadata,omitempty" gorm:"type:json"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at" db:"updated_at"`
}

// DialogType represents the type of dialog/conversation
type DialogType string

const (
	DialogTypeUser     DialogType = "user"     // Direct message between two users
	DialogTypeGroup    DialogType = "group"    // Group chat with multiple users
	DialogTypeChannel  DialogType = "channel"  // Broadcast channel
	DialogTypeBusiness DialogType = "business" // Business/customer service chat
)

// ParticipantRole represents the role of a participant in a dialog
type ParticipantRole string

const (
	ParticipantRoleMember ParticipantRole = "member"
	ParticipantRoleAdmin  ParticipantRole = "admin"
	ParticipantRoleOwner  ParticipantRole = "owner"
	ParticipantRoleGuest  ParticipantRole = "guest"
)

// ParticipantStatus represents the status of a participant in a dialog
type ParticipantStatus string

const (
	ParticipantStatusActive   ParticipantStatus = "active"
	ParticipantStatusInactive ParticipantStatus = "inactive"
	ParticipantStatusBanned   ParticipantStatus = "banned"
	ParticipantStatusLeft     ParticipantStatus = "left"
)

// DialogParticipant represents a participant in a dialog
type DialogParticipant struct {
	ID        uuid.UUID         `json:"id" db:"id"`
	DialogID  uuid.UUID         `json:"dialog_id" db:"dialog_id"`
	UserID    uuid.UUID         `json:"user_id" db:"user_id"`
	Role      ParticipantRole   `json:"role" db:"role"`
	Status    ParticipantStatus `json:"status" db:"status"`
	JoinedAt  time.Time         `json:"joined_at" db:"joined_at"`
	LeftAt    *time.Time        `json:"left_at,omitempty" db:"left_at"`
	UpdatedAt time.Time         `json:"updated_at" db:"updated_at"`
	IsActive  bool              `json:"is_active" db:"is_active"`
	IsMuted   bool              `json:"is_muted" db:"is_muted"`
}

// DialogSettings represents dialog-specific settings and permissions
type DialogSettings struct {
	MaxParticipants    int                    `json:"max_participants,omitempty"`
	IsPublic          bool                   `json:"is_public"`
	JoinByLink        bool                   `json:"join_by_link"`
	MessageHistory    MessageHistoryAccess   `json:"message_history"`
	WhoCanInvite      InvitePermission       `json:"who_can_invite"`
	WhoCanMessage     MessagePermission      `json:"who_can_message"`
	AutoDeleteAfter   *time.Duration         `json:"auto_delete_after,omitempty"`
	CustomFields      JSON `json:"custom_fields,omitempty"`
}

// MessageHistoryAccess controls who can see message history
type MessageHistoryAccess string

const (
	HistoryAccessEveryone MessageHistoryAccess = "everyone"
	HistoryAccessMembers  MessageHistoryAccess = "members"
	HistoryAccessAdmins   MessageHistoryAccess = "admins"
)

// InvitePermission controls who can invite new members
type InvitePermission string

const (
	InvitePermissionEveryone InvitePermission = "everyone"
	InvitePermissionMembers  InvitePermission = "members"
	InvitePermissionAdmins   InvitePermission = "admins"
)

// MessagePermission controls who can send messages
type MessagePermission string

const (
	MessagePermissionEveryone MessagePermission = "everyone"
	MessagePermissionMembers  MessagePermission = "members"
	MessagePermissionAdmins   MessagePermission = "admins"
)

// JSON is a custom type for handling JSON fields in the database
type JSON map[string]interface{}

// Value implements the driver.Valuer interface for JSON
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for JSON
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into JSON", value)
	}

	return json.Unmarshal(jsonData, j)
}

// GormDataType returns the data type for GORM migration
func (JSON) GormDataType() string {
	return "json"
}

// UUIDSlice is a custom type for handling UUID arrays in the database
type UUIDSlice []uuid.UUID

// Value implements the driver.Valuer interface for database storage
func (us UUIDSlice) Value() (driver.Value, error) {
	if us == nil {
		return nil, nil
	}

	strSlice := make([]string, len(us))
	for i, u := range us {
		strSlice[i] = u.String()
	}

	return json.Marshal(strSlice)
}

// Scan implements the sql.Scanner interface for database retrieval
func (us *UUIDSlice) Scan(value interface{}) error {
	if value == nil {
		*us = nil
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into UUIDSlice", value)
	}

	var strSlice []string
	if err := json.Unmarshal(jsonData, &strSlice); err != nil {
		return err
	}

	uuidSlice := make([]uuid.UUID, len(strSlice))
	for i, str := range strSlice {
		u, err := uuid.Parse(str)
		if err != nil {
			return fmt.Errorf("invalid UUID in slice: %s", str)
		}
		uuidSlice[i] = u
	}

	*us = uuidSlice
	return nil
}

// GormDataType returns the data type for GORM migration
func (UUIDSlice) GormDataType() string {
	return "json"
}

// Value implements the driver.Valuer interface for DialogSettings
func (ds DialogSettings) Value() (driver.Value, error) {
	return json.Marshal(ds)
}

// Scan implements the sql.Scanner interface for DialogSettings
func (ds *DialogSettings) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into DialogSettings", value)
	}

	return json.Unmarshal(jsonData, ds)
}

// GormDataType returns the data type for GORM migration
func (DialogSettings) GormDataType() string {
	return "json"
}

// ValidDialogTypes returns all supported dialog types
func ValidDialogTypes() []DialogType {
	return []DialogType{
		DialogTypeUser,
		DialogTypeGroup,
		DialogTypeChannel,
		DialogTypeBusiness,
	}
}

// IsValid validates if the dialog type is supported
func (dt DialogType) IsValid() bool {
	for _, valid := range ValidDialogTypes() {
		if dt == valid {
			return true
		}
	}
	return false
}

// String returns the string representation of DialogType
func (dt DialogType) String() string {
	return string(dt)
}

// ValidParticipantRoles returns all supported participant roles
func ValidParticipantRoles() []ParticipantRole {
	return []ParticipantRole{
		ParticipantRoleMember,
		ParticipantRoleAdmin,
		ParticipantRoleOwner,
		ParticipantRoleGuest,
	}
}

// IsValid validates if the participant role is supported
func (pr ParticipantRole) IsValid() bool {
	for _, valid := range ValidParticipantRoles() {
		if pr == valid {
			return true
		}
	}
	return false
}

// String returns the string representation of ParticipantRole
func (pr ParticipantRole) String() string {
	return string(pr)
}

// BeforeCreate sets up the DialogParticipant before database creation
func (dp *DialogParticipant) BeforeCreate(tx *gorm.DB) error {
	if dp.ID == uuid.Nil {
		dp.ID = uuid.New()
	}
	dp.JoinedAt = time.Now().UTC()
	dp.IsActive = true
	return nil
}

// Validate performs validation on the DialogParticipant model
func (dp *DialogParticipant) Validate() error {
	if dp.DialogID == uuid.Nil {
		return errors.New("dialog_id is required")
	}
	if dp.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}
	if !dp.Role.IsValid() {
		return fmt.Errorf("invalid participant role: %s", dp.Role)
	}
	return nil
}

// GetMaxParticipants returns the maximum participants allowed for this dialog type
func (dt DialogType) GetMaxParticipants() int {
	switch dt {
	case DialogTypeUser:
		return 2
	case DialogTypeGroup:
		return 5000 // Large group support for SEA markets
	case DialogTypeChannel:
		return 200000 // Broadcast channel support
	case DialogTypeBusiness:
		return 100 // Business chat limit
	default:
		return 2
	}
}

// RequiresName checks if this dialog type requires a name
func (dt DialogType) RequiresName() bool {
	return dt == DialogTypeGroup || dt == DialogTypeChannel || dt == DialogTypeBusiness
}

// SupportsAdmins checks if this dialog type supports admin roles
func (dt DialogType) SupportsAdmins() bool {
	return dt == DialogTypeGroup || dt == DialogTypeChannel || dt == DialogTypeBusiness
}

// Validate performs comprehensive validation on the Dialog model
func (d *Dialog) Validate() error {
	var errs []string

	// Dialog type validation
	if !d.Type.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid dialog type: %s", d.Type))
	}

	// Name validation for types that require it
	if d.Type.RequiresName() {
		if d.Name == nil || strings.TrimSpace(*d.Name) == "" {
			errs = append(errs, fmt.Sprintf("name is required for dialog type: %s", d.Type))
		} else if len(*d.Name) > 255 {
			errs = append(errs, "name must not exceed 255 characters")
		}
	}

	// Participants validation
	if len(d.Participants) == 0 {
		errs = append(errs, "at least one participant is required")
	}

	maxParticipants := d.Type.GetMaxParticipants()
	if len(d.Participants) > maxParticipants {
		errs = append(errs, fmt.Sprintf("too many participants: %d, max allowed: %d", len(d.Participants), maxParticipants))
	}

	// Validate unique participants
	if err := d.validateUniqueParticipants(); err != nil {
		errs = append(errs, err.Error())
	}

	// Admin validation for supported types
	if d.Type.SupportsAdmins() && len(d.AdminIDs) > 0 {
		if err := d.validateAdmins(); err != nil {
			errs = append(errs, err.Error())
		}
	}

	// Unread count validation
	if d.UnreadCount < 0 {
		errs = append(errs, "unread_count cannot be negative")
	}

	// Avatar URL validation if provided
	if d.Avatar != nil && *d.Avatar != "" {
		if err := d.validateAvatarURL(*d.Avatar); err != nil {
			errs = append(errs, fmt.Sprintf("invalid avatar URL: %v", err))
		}
	}

	// Settings validation
	if err := d.validateSettings(); err != nil {
		errs = append(errs, fmt.Sprintf("invalid settings: %v", err))
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

// validateUniqueParticipants ensures all participants are unique
func (d *Dialog) validateUniqueParticipants() error {
	seen := make(map[uuid.UUID]bool)
	for _, participant := range d.Participants {
		if participant == uuid.Nil {
			return errors.New("participant ID cannot be nil")
		}
		if seen[participant] {
			return fmt.Errorf("duplicate participant: %s", participant)
		}
		seen[participant] = true
	}
	return nil
}

// validateAdmins ensures all admins are also participants
func (d *Dialog) validateAdmins() error {
	participantSet := make(map[uuid.UUID]bool)
	for _, participant := range d.Participants {
		participantSet[participant] = true
	}

	for _, adminID := range d.AdminIDs {
		if adminID == uuid.Nil {
			return errors.New("admin ID cannot be nil")
		}
		if !participantSet[adminID] {
			return fmt.Errorf("admin %s is not a participant", adminID)
		}
	}

	return nil
}

// validateAvatarURL validates avatar URL format
func (d *Dialog) validateAvatarURL(url string) error {
	// Basic URL validation - could use net/url for more robust validation
	if len(url) > 512 {
		return errors.New("avatar URL too long")
	}
	if !strings.HasPrefix(url, "https://") {
		return errors.New("avatar URL must use HTTPS")
	}
	return nil
}

// validateSettings validates dialog settings
func (d *Dialog) validateSettings() error {
	settings := d.Settings

	// Max participants validation
	if settings.MaxParticipants > 0 {
		typeMax := d.Type.GetMaxParticipants()
		if settings.MaxParticipants > typeMax {
			return fmt.Errorf("max_participants (%d) exceeds type limit (%d)", settings.MaxParticipants, typeMax)
		}
	}

	// Auto delete validation
	if settings.AutoDeleteAfter != nil {
		if *settings.AutoDeleteAfter < time.Hour {
			return errors.New("auto_delete_after must be at least 1 hour")
		}
		if *settings.AutoDeleteAfter > 365*24*time.Hour {
			return errors.New("auto_delete_after cannot exceed 1 year")
		}
	}

	return nil
}

// BeforeCreate sets up the dialog before database creation
func (d *Dialog) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now

	// Set default settings based on dialog type
	d.setDefaultSettings()

	// Validate before creation
	return d.Validate()
}

// BeforeUpdate sets up the dialog before database update
func (d *Dialog) BeforeUpdate(tx *gorm.DB) error {
	// Update timestamp
	d.UpdatedAt = time.Now().UTC()

	// Validate before update
	return d.Validate()
}

// setDefaultSettings sets default settings based on dialog type
func (d *Dialog) setDefaultSettings() {
	if d.Settings.MaxParticipants == 0 {
		d.Settings.MaxParticipants = d.Type.GetMaxParticipants()
	}

	switch d.Type {
	case DialogTypeUser:
		d.Settings.IsPublic = false
		d.Settings.JoinByLink = false
		d.Settings.MessageHistory = HistoryAccessMembers
		d.Settings.WhoCanInvite = InvitePermissionMembers
		d.Settings.WhoCanMessage = MessagePermissionMembers
	case DialogTypeGroup:
		d.Settings.IsPublic = false
		d.Settings.JoinByLink = false
		d.Settings.MessageHistory = HistoryAccessMembers
		d.Settings.WhoCanInvite = InvitePermissionAdmins
		d.Settings.WhoCanMessage = MessagePermissionMembers
	case DialogTypeChannel:
		d.Settings.IsPublic = true
		d.Settings.JoinByLink = true
		d.Settings.MessageHistory = HistoryAccessEveryone
		d.Settings.WhoCanInvite = InvitePermissionAdmins
		d.Settings.WhoCanMessage = MessagePermissionAdmins
	case DialogTypeBusiness:
		d.Settings.IsPublic = false
		d.Settings.JoinByLink = false
		d.Settings.MessageHistory = HistoryAccessMembers
		d.Settings.WhoCanInvite = InvitePermissionAdmins
		d.Settings.WhoCanMessage = MessagePermissionMembers
	}
}

// AddParticipant adds a new participant to the dialog
func (d *Dialog) AddParticipant(userID uuid.UUID) error {
	if userID == uuid.Nil {
		return errors.New("user ID cannot be nil")
	}

	// Check if already a participant
	for _, participant := range d.Participants {
		if participant == userID {
			return errors.New("user is already a participant")
		}
	}

	// Check participant limit
	if len(d.Participants) >= d.Settings.MaxParticipants {
		return fmt.Errorf("dialog has reached maximum participants (%d)", d.Settings.MaxParticipants)
	}

	d.Participants = append(d.Participants, userID)
	d.UpdatedAt = time.Now().UTC()

	return d.Validate()
}

// RemoveParticipant removes a participant from the dialog
func (d *Dialog) RemoveParticipant(userID uuid.UUID) error {
	if userID == uuid.Nil {
		return errors.New("user ID cannot be nil")
	}

	// Find and remove participant
	for i, participant := range d.Participants {
		if participant == userID {
			d.Participants = append(d.Participants[:i], d.Participants[i+1:]...)

			// Also remove from admins if present
			d.RemoveAdmin(userID)

			d.UpdatedAt = time.Now().UTC()
			return d.Validate()
		}
	}

	return errors.New("user is not a participant")
}

// AddAdmin adds a user as an admin
func (d *Dialog) AddAdmin(userID uuid.UUID) error {
	if !d.Type.SupportsAdmins() {
		return fmt.Errorf("dialog type %s does not support admins", d.Type)
	}

	if userID == uuid.Nil {
		return errors.New("user ID cannot be nil")
	}

	// Check if user is a participant
	isParticipant := false
	for _, participant := range d.Participants {
		if participant == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return errors.New("user must be a participant to become admin")
	}

	// Check if already an admin
	for _, adminID := range d.AdminIDs {
		if adminID == userID {
			return errors.New("user is already an admin")
		}
	}

	d.AdminIDs = append(d.AdminIDs, userID)
	d.UpdatedAt = time.Now().UTC()

	return nil
}

// RemoveAdmin removes a user from admin role
func (d *Dialog) RemoveAdmin(userID uuid.UUID) error {
	for i, adminID := range d.AdminIDs {
		if adminID == userID {
			d.AdminIDs = append(d.AdminIDs[:i], d.AdminIDs[i+1:]...)
			d.UpdatedAt = time.Now().UTC()
			return nil
		}
	}

	return errors.New("user is not an admin")
}

// IsParticipant checks if a user is a participant
func (d *Dialog) IsParticipant(userID uuid.UUID) bool {
	for _, participant := range d.Participants {
		if participant == userID {
			return true
		}
	}
	return false
}

// IsAdmin checks if a user is an admin
func (d *Dialog) IsAdmin(userID uuid.UUID) bool {
	for _, adminID := range d.AdminIDs {
		if adminID == userID {
			return true
		}
	}
	return false
}

// CanUserInvite checks if a user can invite others to this dialog
func (d *Dialog) CanUserInvite(userID uuid.UUID) bool {
	switch d.Settings.WhoCanInvite {
	case InvitePermissionEveryone:
		return true
	case InvitePermissionMembers:
		return d.IsParticipant(userID)
	case InvitePermissionAdmins:
		return d.IsAdmin(userID)
	default:
		return false
	}
}

// CanUserMessage checks if a user can send messages to this dialog
func (d *Dialog) CanUserMessage(userID uuid.UUID) bool {
	switch d.Settings.WhoCanMessage {
	case MessagePermissionEveryone:
		return true
	case MessagePermissionMembers:
		return d.IsParticipant(userID)
	case MessagePermissionAdmins:
		return d.IsAdmin(userID)
	default:
		return false
	}
}

// UpdateUnreadCount updates the unread message count
func (d *Dialog) UpdateUnreadCount(count int) error {
	if count < 0 {
		return errors.New("unread count cannot be negative")
	}
	d.UnreadCount = count
	d.UpdatedAt = time.Now().UTC()
	return nil
}

// SetLastMessage updates the last message reference
func (d *Dialog) SetLastMessage(messageID uuid.UUID) {
	d.LastMessageID = &messageID
	d.UpdatedAt = time.Now().UTC()
}

// ToPublicDialog returns a sanitized version for public API responses
func (d *Dialog) ToPublicDialog(forUserID uuid.UUID) map[string]interface{} {
	response := map[string]interface{}{
		"id":              d.ID,
		"type":            d.Type,
		"name":            d.Name,
		"avatar":          d.Avatar,
		"unread_count":    d.UnreadCount,
		"is_pinned":       d.IsPinned,
		"is_archived":     d.IsArchived,
		"is_muted":        d.IsMuted,
		"last_message_id": d.LastMessageID,
		"created_at":      d.CreatedAt,
		"updated_at":      d.UpdatedAt,
	}

	// Include participant info if user is a participant
	if d.IsParticipant(forUserID) {
		response["participants"] = d.Participants
		response["is_admin"] = d.IsAdmin(forUserID)

		// Include settings for admins or appropriate permission levels
		if d.IsAdmin(forUserID) || d.Settings.IsPublic {
			response["settings"] = d.Settings
		}
	}

	return response
}

// GetMaxParticipants returns the maximum participants allowed for this dialog
func (d *Dialog) GetMaxParticipants() int {
	return d.Type.GetMaxParticipants()
}