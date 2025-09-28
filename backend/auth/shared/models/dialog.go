package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DialogType represents the type of dialog/conversation
type DialogType string

const (
	DialogTypeDirect  DialogType = "direct"  // One-on-one conversation
	DialogTypeGroup   DialogType = "group"   // Group conversation
	DialogTypeChannel DialogType = "channel" // Channel/broadcast
)

// IsValid checks if the dialog type is valid
func (dt DialogType) IsValid() bool {
	switch dt {
	case DialogTypeDirect, DialogTypeGroup, DialogTypeChannel:
		return true
	default:
		return false
	}
}

// DialogStatus represents the status of a dialog
type DialogStatus string

const (
	DialogStatusActive   DialogStatus = "active"
	DialogStatusArchived DialogStatus = "archived"
	DialogStatusDeleted  DialogStatus = "deleted"
	DialogStatusMuted    DialogStatus = "muted"
)

// IsValid checks if the dialog status is valid
func (ds DialogStatus) IsValid() bool {
	switch ds {
	case DialogStatusActive, DialogStatusArchived, DialogStatusDeleted, DialogStatusMuted:
		return true
	default:
		return false
	}
}

// DialogPrivacy represents the privacy level of a dialog
type DialogPrivacy string

const (
	DialogPrivacyPublic   DialogPrivacy = "public"   // Anyone can find and join
	DialogPrivacyPrivate  DialogPrivacy = "private"  // Invite only
	DialogPrivacySecret   DialogPrivacy = "secret"   // Hidden from search
)

// IsValid checks if the dialog privacy is valid
func (dp DialogPrivacy) IsValid() bool {
	switch dp {
	case DialogPrivacyPublic, DialogPrivacyPrivate, DialogPrivacySecret:
		return true
	default:
		return false
	}
}

// DialogParticipant represents a participant in a dialog
type DialogParticipant struct {
	UserID      uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid;not null"`
	Role        string    `json:"role" gorm:"column:role;size:20;not null;default:'member'"`
	JoinedAt    time.Time `json:"joined_at" gorm:"column:joined_at;not null"`
	LastReadAt  *time.Time `json:"last_read_at,omitempty" gorm:"column:last_read_at"`
	IsActive    bool      `json:"is_active" gorm:"column:is_active;default:true"`
	IsMuted     bool      `json:"is_muted" gorm:"column:is_muted;default:false"`
	Permissions []string  `json:"permissions,omitempty" gorm:"column:permissions;type:jsonb"`
}

// DialogSettings represents dialog configuration and settings
type DialogSettings struct {
	IsEncrypted          bool     `json:"is_encrypted" gorm:"column:is_encrypted;default:false"`
	AutoDeleteMessages   bool     `json:"auto_delete_messages" gorm:"column:auto_delete_messages;default:false"`
	MessageRetentionDays int      `json:"message_retention_days" gorm:"column:message_retention_days;default:0"`
	AllowedMessageTypes  []string `json:"allowed_message_types" gorm:"column:allowed_message_types;type:jsonb"`
	MaxParticipants      int      `json:"max_participants" gorm:"column:max_participants;default:0"`
	RequireApproval      bool     `json:"require_approval" gorm:"column:require_approval;default:false"`
	AllowInvites         bool     `json:"allow_invites" gorm:"column:allow_invites;default:true"`
	NotificationsEnabled bool     `json:"notifications_enabled" gorm:"column:notifications_enabled;default:true"`
	ReadReceiptsEnabled  bool     `json:"read_receipts_enabled" gorm:"column:read_receipts_enabled;default:true"`
	TypingIndicators     bool     `json:"typing_indicators" gorm:"column:typing_indicators;default:true"`
}

// DialogModeration represents moderation settings and rules
type DialogModeration struct {
	IsModerated       bool     `json:"is_moderated" gorm:"column:is_moderated;default:false"`
	AutoModeration    bool     `json:"auto_moderation" gorm:"column:auto_moderation;default:false"`
	BannedWords       []string `json:"banned_words,omitempty" gorm:"column:banned_words;type:jsonb"`
	RequiredTags      []string `json:"required_tags,omitempty" gorm:"column:required_tags;type:jsonb"`
	SlowModeInterval  int      `json:"slow_mode_interval" gorm:"column:slow_mode_interval;default:0"` // seconds
	MaxMessageLength  int      `json:"max_message_length" gorm:"column:max_message_length;default:4000"`
	AllowExternalLinks bool    `json:"allow_external_links" gorm:"column:allow_external_links;default:true"`
	AllowMedia        bool     `json:"allow_media" gorm:"column:allow_media;default:true"`
	ModerationRules   map[string]interface{} `json:"moderation_rules,omitempty" gorm:"column:moderation_rules;type:jsonb"`
}

// Dialog represents a conversation/chat dialog
type Dialog struct {
	ID   uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name string    `json:"name" gorm:"column:name;size:100;not null"`

	// Dialog properties
	Type        DialogType    `json:"type" gorm:"column:type;type:varchar(20);not null"`
	Status      DialogStatus  `json:"status" gorm:"column:status;type:varchar(20);not null;default:'active'"`
	Privacy     DialogPrivacy `json:"privacy" gorm:"column:privacy;type:varchar(20);not null;default:'private'"`
	Description string        `json:"description,omitempty" gorm:"column:description;size:500"`
	AvatarURL   string        `json:"avatar_url,omitempty" gorm:"column:avatar_url;size:500"`

	// Creator and ownership
	CreatedByID uuid.UUID `json:"created_by_id" gorm:"type:uuid;not null;index"`
	OwnerID     *uuid.UUID `json:"owner_id,omitempty" gorm:"type:uuid;index"`

	// Participants (stored as JSONB for flexibility)
	Participants []DialogParticipant `json:"participants" gorm:"column:participants;type:jsonb"`
	ParticipantCount int             `json:"participant_count" gorm:"column:participant_count;default:0"`

	// Last message information
	LastMessageID      *uuid.UUID `json:"last_message_id,omitempty" gorm:"column:last_message_id;type:uuid"`
	LastMessageContent string     `json:"last_message_content,omitempty" gorm:"column:last_message_content;size:500"`
	LastMessageAt      *time.Time `json:"last_message_at,omitempty" gorm:"column:last_message_at"`
	LastMessageSenderID *uuid.UUID `json:"last_message_sender_id,omitempty" gorm:"column:last_message_sender_id;type:uuid"`

	// Activity metrics
	MessageCount     int64     `json:"message_count" gorm:"column:message_count;default:0"`
	LastActivityAt   time.Time `json:"last_activity_at" gorm:"column:last_activity_at;not null"`
	UnreadCount      int       `json:"unread_count" gorm:"column:unread_count;default:0"`
	IsActive         bool      `json:"is_active" gorm:"column:is_active;default:true"`

	// Configuration and moderation
	Settings   DialogSettings   `json:"settings" gorm:"embedded;embeddedPrefix:settings_"`
	Moderation DialogModeration `json:"moderation" gorm:"embedded;embeddedPrefix:moderation_"`

	// Regional compliance
	DataRegion     string `json:"data_region" gorm:"column:data_region;size:20"`
	ComplianceData map[string]interface{} `json:"compliance_data,omitempty" gorm:"column:compliance_data;type:jsonb"`

	// Search and discovery
	Tags           []string `json:"tags,omitempty" gorm:"column:tags;type:jsonb"`
	SearchKeywords []string `json:"search_keywords,omitempty" gorm:"column:search_keywords;type:jsonb"`
	Language       string   `json:"language,omitempty" gorm:"column:language;size:5;default:'en'"`
	Category       string   `json:"category,omitempty" gorm:"column:category;size:50"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Relationships
	CreatedBy   *User     `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID;references:ID"`
	Owner       *User     `json:"owner,omitempty" gorm:"foreignKey:OwnerID;references:ID"`
	LastMessage *Message  `json:"last_message,omitempty" gorm:"foreignKey:LastMessageID;references:ID"`
	Messages    []Message `json:"messages,omitempty" gorm:"foreignKey:DialogID;references:ID"`
}

// TableName returns the table name for the Dialog model
func (Dialog) TableName() string {
	return "dialogs"
}

// BeforeCreate sets up the dialog before creation
func (d *Dialog) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}

	// Set initial activity timestamp
	if d.LastActivityAt.IsZero() {
		d.LastActivityAt = time.Now()
	}

	// Set default data region
	if d.DataRegion == "" {
		d.DataRegion = "sea-central" // Default region
	}

	// Initialize default settings
	if len(d.Settings.AllowedMessageTypes) == 0 {
		d.Settings.AllowedMessageTypes = []string{"text", "image", "video", "audio", "file", "location", "sticker"}
	}

	// Set max participants based on dialog type
	if d.Settings.MaxParticipants == 0 {
		switch d.Type {
		case DialogTypeDirect:
			d.Settings.MaxParticipants = 2
		case DialogTypeGroup:
			d.Settings.MaxParticipants = 1000
		case DialogTypeChannel:
			d.Settings.MaxParticipants = 10000
		}
	}

	// Add creator as first participant
	if len(d.Participants) == 0 {
		creator := DialogParticipant{
			UserID:      d.CreatedByID,
			Role:        "owner",
			JoinedAt:    time.Now(),
			IsActive:    true,
			IsMuted:     false,
			Permissions: []string{"read", "write", "invite", "manage", "admin"},
		}
		d.Participants = append(d.Participants, creator)
		d.ParticipantCount = 1
	}

	// Set owner if not set
	if d.OwnerID == nil {
		d.OwnerID = &d.CreatedByID
	}

	// Generate search keywords
	d.SearchKeywords = d.GenerateSearchKeywords()

	// Validate the dialog
	if err := d.Validate(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate validates the dialog before updating
func (d *Dialog) BeforeUpdate(tx *gorm.DB) error {
	// Update search keywords
	d.SearchKeywords = d.GenerateSearchKeywords()

	return d.Validate()
}

// Validate validates the dialog data
func (d *Dialog) Validate() error {
	// Validate UUIDs
	if d.ID == uuid.Nil {
		return fmt.Errorf("dialog ID cannot be nil")
	}
	if d.CreatedByID == uuid.Nil {
		return fmt.Errorf("created by ID cannot be nil")
	}

	// Validate type, status, and privacy
	if !d.Type.IsValid() {
		return fmt.Errorf("invalid dialog type: %s", d.Type)
	}
	if !d.Status.IsValid() {
		return fmt.Errorf("invalid dialog status: %s", d.Status)
	}
	if !d.Privacy.IsValid() {
		return fmt.Errorf("invalid dialog privacy: %s", d.Privacy)
	}

	// Validate name
	if len(d.Name) == 0 || len(d.Name) > 100 {
		return fmt.Errorf("dialog name must be between 1 and 100 characters")
	}

	// Validate participant count constraints
	if d.Type == DialogTypeDirect && d.ParticipantCount > 2 {
		return fmt.Errorf("direct dialog cannot have more than 2 participants")
	}

	if d.ParticipantCount > d.Settings.MaxParticipants && d.Settings.MaxParticipants > 0 {
		return fmt.Errorf("participant count (%d) exceeds maximum allowed (%d)",
			d.ParticipantCount, d.Settings.MaxParticipants)
	}

	// Validate participants
	if err := d.validateParticipants(); err != nil {
		return err
	}

	return nil
}

// validateParticipants validates the participants list
func (d *Dialog) validateParticipants() error {
	if len(d.Participants) != d.ParticipantCount {
		return fmt.Errorf("participant count mismatch: expected %d, got %d",
			d.ParticipantCount, len(d.Participants))
	}

	userIDs := make(map[uuid.UUID]bool)
	ownerCount := 0

	for i, participant := range d.Participants {
		// Check for duplicate users
		if userIDs[participant.UserID] {
			return fmt.Errorf("duplicate participant at index %d: user %s", i, participant.UserID)
		}
		userIDs[participant.UserID] = true

		// Validate role
		if !isValidParticipantRole(participant.Role) {
			return fmt.Errorf("invalid participant role at index %d: %s", i, participant.Role)
		}

		// Count owners
		if participant.Role == "owner" {
			ownerCount++
		}
	}

	// Direct dialogs should have exactly 0 or 1 owner
	// Groups and channels should have exactly 1 owner
	if d.Type == DialogTypeDirect && ownerCount > 1 {
		return fmt.Errorf("direct dialog cannot have more than 1 owner")
	} else if (d.Type == DialogTypeGroup || d.Type == DialogTypeChannel) && ownerCount != 1 {
		return fmt.Errorf("%s dialog must have exactly 1 owner", d.Type)
	}

	return nil
}

// isValidParticipantRole checks if a participant role is valid
func isValidParticipantRole(role string) bool {
	validRoles := map[string]bool{
		"owner":     true,
		"admin":     true,
		"moderator": true,
		"member":    true,
		"guest":     true,
		"banned":    true,
	}
	return validRoles[role]
}

// AddParticipant adds a new participant to the dialog
func (d *Dialog) AddParticipant(userID uuid.UUID, role string) error {
	// Check if user is already a participant
	for _, participant := range d.Participants {
		if participant.UserID == userID {
			return fmt.Errorf("user %s is already a participant", userID)
		}
	}

	// Check participant limit
	if d.ParticipantCount >= d.Settings.MaxParticipants && d.Settings.MaxParticipants > 0 {
		return fmt.Errorf("dialog has reached maximum participant limit (%d)", d.Settings.MaxParticipants)
	}

	// Validate role
	if !isValidParticipantRole(role) {
		return fmt.Errorf("invalid participant role: %s", role)
	}

	// Default permissions based on role
	permissions := getDefaultPermissions(role)

	// Create new participant
	participant := DialogParticipant{
		UserID:      userID,
		Role:        role,
		JoinedAt:    time.Now(),
		IsActive:    true,
		IsMuted:     false,
		Permissions: permissions,
	}

	d.Participants = append(d.Participants, participant)
	d.ParticipantCount++
	d.LastActivityAt = time.Now()

	return nil
}

// RemoveParticipant removes a participant from the dialog
func (d *Dialog) RemoveParticipant(userID uuid.UUID) error {
	for i, participant := range d.Participants {
		if participant.UserID == userID {
			// Don't allow removing the last owner
			if participant.Role == "owner" && d.countOwners() == 1 {
				return fmt.Errorf("cannot remove the last owner from dialog")
			}

			// Remove participant
			d.Participants = append(d.Participants[:i], d.Participants[i+1:]...)
			d.ParticipantCount--
			d.LastActivityAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("user %s is not a participant in this dialog", userID)
}

// UpdateParticipantRole updates a participant's role
func (d *Dialog) UpdateParticipantRole(userID uuid.UUID, newRole string) error {
	if !isValidParticipantRole(newRole) {
		return fmt.Errorf("invalid participant role: %s", newRole)
	}

	for i, participant := range d.Participants {
		if participant.UserID == userID {
			// Check if removing the last owner
			if participant.Role == "owner" && newRole != "owner" && d.countOwners() == 1 {
				return fmt.Errorf("cannot change role of the last owner")
			}

			d.Participants[i].Role = newRole
			d.Participants[i].Permissions = getDefaultPermissions(newRole)
			d.LastActivityAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("user %s is not a participant in this dialog", userID)
}

// countOwners counts the number of owners in the dialog
func (d *Dialog) countOwners() int {
	count := 0
	for _, participant := range d.Participants {
		if participant.Role == "owner" {
			count++
		}
	}
	return count
}

// getDefaultPermissions returns default permissions for a role
func getDefaultPermissions(role string) []string {
	permissions := map[string][]string{
		"owner":     {"read", "write", "invite", "manage", "admin", "delete"},
		"admin":     {"read", "write", "invite", "manage", "moderate"},
		"moderator": {"read", "write", "invite", "moderate"},
		"member":    {"read", "write"},
		"guest":     {"read"},
		"banned":    {},
	}
	return permissions[role]
}

// IsParticipant checks if a user is a participant in the dialog
func (d *Dialog) IsParticipant(userID uuid.UUID) bool {
	for _, participant := range d.Participants {
		if participant.UserID == userID && participant.IsActive {
			return true
		}
	}
	return false
}

// GetParticipant gets a participant by user ID
func (d *Dialog) GetParticipant(userID uuid.UUID) (*DialogParticipant, error) {
	for i, participant := range d.Participants {
		if participant.UserID == userID {
			return &d.Participants[i], nil
		}
	}
	return nil, fmt.Errorf("user %s is not a participant in this dialog", userID)
}

// HasPermission checks if a user has a specific permission
func (d *Dialog) HasPermission(userID uuid.UUID, permission string) bool {
	participant, err := d.GetParticipant(userID)
	if err != nil {
		return false
	}

	for _, perm := range participant.Permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

// CanUserWrite checks if a user can write messages
func (d *Dialog) CanUserWrite(userID uuid.UUID) bool {
	return d.IsParticipant(userID) && d.HasPermission(userID, "write")
}

// CanUserInvite checks if a user can invite others
func (d *Dialog) CanUserInvite(userID uuid.UUID) bool {
	return d.Settings.AllowInvites && d.HasPermission(userID, "invite")
}

// UpdateLastMessage updates the last message information
func (d *Dialog) UpdateLastMessage(messageID uuid.UUID, content string, senderID uuid.UUID) {
	d.LastMessageID = &messageID
	d.LastMessageContent = content
	d.LastMessageSenderID = &senderID
	now := time.Now()
	d.LastMessageAt = &now
	d.LastActivityAt = now
	d.MessageCount++
}

// UpdateUnreadCount updates the unread count for the dialog
func (d *Dialog) UpdateUnreadCount(count int) {
	d.UnreadCount = count
	if count < 0 {
		d.UnreadCount = 0
	}
}

// MarkAsRead marks the dialog as read for a specific user
func (d *Dialog) MarkAsRead(userID uuid.UUID) error {
	participant, err := d.GetParticipant(userID)
	if err != nil {
		return err
	}

	now := time.Now()
	participant.LastReadAt = &now
	return nil
}

// IsPublic checks if the dialog is publicly accessible
func (d *Dialog) IsPublic() bool {
	return d.Privacy == DialogPrivacyPublic
}

// IsPrivate checks if the dialog is private
func (d *Dialog) IsPrivate() bool {
	return d.Privacy == DialogPrivacyPrivate
}

// IsSecret checks if the dialog is secret
func (d *Dialog) IsSecret() bool {
	return d.Privacy == DialogPrivacySecret
}

// IsDirect checks if the dialog is a direct conversation
func (d *Dialog) IsDirect() bool {
	return d.Type == DialogTypeDirect
}

// IsGroup checks if the dialog is a group conversation
func (d *Dialog) IsGroup() bool {
	return d.Type == DialogTypeGroup
}

// IsChannel checks if the dialog is a channel
func (d *Dialog) IsChannel() bool {
	return d.Type == DialogTypeChannel
}

// GenerateSearchKeywords generates search keywords for the dialog
func (d *Dialog) GenerateSearchKeywords() []string {
	keywords := []string{
		d.Name,
		string(d.Type),
		string(d.Privacy),
		d.Category,
		d.Language,
	}

	// Add description keywords
	if d.Description != "" {
		descWords := strings.Fields(strings.ToLower(d.Description))
		keywords = append(keywords, descWords...)
	}

	// Add tags
	keywords = append(keywords, d.Tags...)

	// Remove duplicates and empty strings
	seen := make(map[string]bool)
	var unique []string
	for _, keyword := range keywords {
		if keyword != "" && !seen[strings.ToLower(keyword)] {
			seen[strings.ToLower(keyword)] = true
			unique = append(unique, strings.ToLower(keyword))
		}
	}

	return unique
}

// GetDialogSummary returns a summary of dialog information
func (d *Dialog) GetDialogSummary() map[string]interface{} {
	return map[string]interface{}{
		"id":               d.ID,
		"name":             d.Name,
		"type":             d.Type,
		"status":           d.Status,
		"privacy":          d.Privacy,
		"participant_count": d.ParticipantCount,
		"message_count":    d.MessageCount,
		"unread_count":     d.UnreadCount,
		"is_active":        d.IsActive,
		"is_public":        d.IsPublic(),
		"is_private":       d.IsPrivate(),
		"is_secret":        d.IsSecret(),
		"is_direct":        d.IsDirect(),
		"is_group":         d.IsGroup(),
		"is_channel":       d.IsChannel(),
		"created_at":       d.CreatedAt,
		"last_activity":    d.LastActivityAt,
		"last_message_at":  d.LastMessageAt,
	}
}

// MarshalJSON customizes JSON serialization
func (d *Dialog) MarshalJSON() ([]byte, error) {
	type Alias Dialog
	return json.Marshal(&struct {
		*Alias
		DialogSummary  map[string]interface{} `json:"dialog_summary"`
		SearchKeywords []string               `json:"search_keywords,omitempty"`
	}{
		Alias:          (*Alias)(d),
		DialogSummary:  d.GetDialogSummary(),
		SearchKeywords: d.SearchKeywords,
	})
}

// Helper functions for dialog management

// CreateDirectDialog creates a direct dialog between two users
func CreateDirectDialog(user1ID, user2ID uuid.UUID, createdByID uuid.UUID) *Dialog {
	name := fmt.Sprintf("Direct: %s-%s", user1ID.String()[:8], user2ID.String()[:8])

	dialog := &Dialog{
		Name:        name,
		Type:        DialogTypeDirect,
		Status:      DialogStatusActive,
		Privacy:     DialogPrivacyPrivate,
		CreatedByID: createdByID,
		OwnerID:     &createdByID,
		Settings: DialogSettings{
			MaxParticipants:      2,
			AllowInvites:         false,
			NotificationsEnabled: true,
			ReadReceiptsEnabled:  true,
			TypingIndicators:     true,
		},
	}

	// Add both participants
	now := time.Now()
	dialog.Participants = []DialogParticipant{
		{
			UserID:      user1ID,
			Role:        "member",
			JoinedAt:    now,
			IsActive:    true,
			Permissions: []string{"read", "write"},
		},
		{
			UserID:      user2ID,
			Role:        "member",
			JoinedAt:    now,
			IsActive:    true,
			Permissions: []string{"read", "write"},
		},
	}
	dialog.ParticipantCount = 2

	return dialog
}

// CreateGroupDialog creates a group dialog
func CreateGroupDialog(name, description string, createdByID uuid.UUID, privacy DialogPrivacy) *Dialog {
	return &Dialog{
		Name:        name,
		Description: description,
		Type:        DialogTypeGroup,
		Status:      DialogStatusActive,
		Privacy:     privacy,
		CreatedByID: createdByID,
		OwnerID:     &createdByID,
		Settings: DialogSettings{
			MaxParticipants:      1000,
			AllowInvites:         true,
			RequireApproval:      privacy == DialogPrivacyPrivate,
			NotificationsEnabled: true,
			ReadReceiptsEnabled:  true,
			TypingIndicators:     true,
		},
		Moderation: DialogModeration{
			MaxMessageLength:   4000,
			AllowExternalLinks: true,
			AllowMedia:         true,
		},
	}
}

// CreateChannelDialog creates a channel dialog
func CreateChannelDialog(name, description string, createdByID uuid.UUID, privacy DialogPrivacy) *Dialog {
	return &Dialog{
		Name:        name,
		Description: description,
		Type:        DialogTypeChannel,
		Status:      DialogStatusActive,
		Privacy:     privacy,
		CreatedByID: createdByID,
		OwnerID:     &createdByID,
		Settings: DialogSettings{
			MaxParticipants:      10000,
			AllowInvites:         privacy == DialogPrivacyPublic,
			RequireApproval:      privacy != DialogPrivacyPublic,
			NotificationsEnabled: true,
			ReadReceiptsEnabled:  false, // Typically disabled for channels
			TypingIndicators:     false, // Typically disabled for channels
		},
		Moderation: DialogModeration{
			IsModerated:        true,
			AutoModeration:     true,
			MaxMessageLength:   4000,
			AllowExternalLinks: true,
			AllowMedia:         true,
		},
	}
}