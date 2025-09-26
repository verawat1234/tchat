package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// MessageType represents the type of message content
type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeVideo    MessageType = "video"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeFile     MessageType = "file"
	MessageTypeLocation MessageType = "location"
	MessageTypeSticker  MessageType = "sticker"
	MessageTypeSystem   MessageType = "system"
)

// IsValid checks if the message type is valid
func (mt MessageType) IsValid() bool {
	switch mt {
	case MessageTypeText, MessageTypeImage, MessageTypeVideo, MessageTypeAudio,
		 MessageTypeFile, MessageTypeLocation, MessageTypeSticker, MessageTypeSystem:
		return true
	default:
		return false
	}
}

// DeliveryStatus represents the delivery status of a message
type DeliveryStatus string

const (
	DeliveryStatusSent      DeliveryStatus = "sent"
	DeliveryStatusDelivered DeliveryStatus = "delivered"
	DeliveryStatusRead      DeliveryStatus = "read"
	DeliveryStatusFailed    DeliveryStatus = "failed"
)

// IsValid checks if the delivery status is valid
func (ds DeliveryStatus) IsValid() bool {
	switch ds {
	case DeliveryStatusSent, DeliveryStatusDelivered, DeliveryStatusRead, DeliveryStatusFailed:
		return true
	default:
		return false
	}
}

// MessageReaction represents a reaction to a message
type MessageReaction struct {
	Emoji   string      `json:"emoji"`
	UserIDs []uuid.UUID `json:"user_ids"`
}

// ReadReceipt represents a read receipt for a message
type ReadReceipt struct {
	UserID uuid.UUID `json:"user_id"`
	ReadAt time.Time `json:"read_at"`
}

// MessageMetadata represents type-specific metadata for messages
type MessageMetadata struct {
	// File message metadata
	FileName string `json:"file_name,omitempty"`
	FileSize int64  `json:"file_size,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	URL      string `json:"url,omitempty"`

	// Location message metadata
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Address   string  `json:"address,omitempty"`

	// Media message metadata
	Duration     int    `json:"duration,omitempty"`      // For audio/video in seconds
	Width        int    `json:"width,omitempty"`         // For images/videos
	Height       int    `json:"height,omitempty"`        // For images/videos
	ThumbnailURL string `json:"thumbnail_url,omitempty"` // For videos

	// Sticker metadata
	StickerPack string `json:"sticker_pack,omitempty"`
	StickerID   string `json:"sticker_id,omitempty"`

	// System message metadata
	SystemType string                 `json:"system_type,omitempty"` // user_joined, user_left, etc.
	SystemData map[string]interface{} `json:"system_data,omitempty"`

	// Encryption metadata (for end-to-end encryption)
	EncryptionKeyID string `json:"encryption_key_id,omitempty"`
	EncryptionIV    string `json:"encryption_iv,omitempty"`
}

// Message represents a message in the messaging system
// Designed for ScyllaDB storage with optimal performance
type Message struct {
	// Primary key (partition key + clustering key)
	ID       uuid.UUID `json:"id" cql:"id"`
	DialogID uuid.UUID `json:"dialog_id" cql:"dialog_id"` // Partition key

	// Message content
	SenderID    uuid.UUID       `json:"sender_id" cql:"sender_id"`
	MessageType MessageType     `json:"message_type" cql:"message_type"`
	Content     string          `json:"content" cql:"content"`           // Encrypted content
	Metadata    MessageMetadata `json:"metadata" cql:"metadata"`

	// Threading and replies
	ReplyToID   *uuid.UUID `json:"reply_to_id,omitempty" cql:"reply_to_id"`
	ThreadID    *uuid.UUID `json:"thread_id,omitempty" cql:"thread_id"`
	IsThreadRoot bool       `json:"is_thread_root" cql:"is_thread_root"`

	// Reactions and interactions
	Reactions []MessageReaction `json:"reactions" cql:"reactions"`

	// Delivery tracking
	DeliveryStatus DeliveryStatus  `json:"delivery_status" cql:"delivery_status"`
	ReadReceipts   []ReadReceipt   `json:"read_receipts" cql:"read_receipts"`

	// Regional compliance
	DataRegion string `json:"data_region" cql:"data_region"`
	Locale     string `json:"locale,omitempty" cql:"locale"`

	// Message versioning and editing
	EditedAt      *time.Time `json:"edited_at,omitempty" cql:"edited_at"`
	OriginalID    *uuid.UUID `json:"original_id,omitempty" cql:"original_id"`
	EditCount     int        `json:"edit_count" cql:"edit_count"`
	IsDeleted     bool       `json:"is_deleted" cql:"is_deleted"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" cql:"deleted_at"`
	DeletedBy     *uuid.UUID `json:"deleted_by,omitempty" cql:"deleted_by"`

	// Timestamps (clustering key)
	CreatedAt time.Time `json:"created_at" cql:"created_at"` // Clustering key for time-based ordering
	UpdatedAt time.Time `json:"updated_at" cql:"updated_at"`

	// Message processing metadata
	ProcessedAt   *time.Time `json:"processed_at,omitempty" cql:"processed_at"`
	DeliveredAt   *time.Time `json:"delivered_at,omitempty" cql:"delivered_at"`
	FirstReadAt   *time.Time `json:"first_read_at,omitempty" cql:"first_read_at"`

	// Search and indexing
	SearchKeywords []string `json:"search_keywords,omitempty" cql:"search_keywords"`
	ContentHash    string   `json:"content_hash,omitempty" cql:"content_hash"`

	// Message size and limits
	ContentLength int `json:"content_length" cql:"content_length"`
	AttachmentCount int `json:"attachment_count" cql:"attachment_count"`
}

// NewMessage creates a new message with defaults
func NewMessage(dialogID, senderID uuid.UUID, messageType MessageType, content string) *Message {
	now := time.Now()
	return &Message{
		ID:             uuid.New(),
		DialogID:       dialogID,
		SenderID:       senderID,
		MessageType:    messageType,
		Content:        content,
		DeliveryStatus: DeliveryStatusSent,
		CreatedAt:      now,
		UpdatedAt:      now,
		ContentLength:  len(content),
		Reactions:      []MessageReaction{},
		ReadReceipts:   []ReadReceipt{},
		SearchKeywords: []string{},
	}
}

// Validate validates the message data
func (m *Message) Validate() error {
	// Validate UUIDs
	if m.ID == uuid.Nil {
		return fmt.Errorf("message ID cannot be nil")
	}
	if m.DialogID == uuid.Nil {
		return fmt.Errorf("dialog ID cannot be nil")
	}
	if m.SenderID == uuid.Nil {
		return fmt.Errorf("sender ID cannot be nil")
	}

	// Validate message type
	if !m.MessageType.IsValid() {
		return fmt.Errorf("invalid message type: %s", m.MessageType)
	}

	// Validate delivery status
	if !m.DeliveryStatus.IsValid() {
		return fmt.Errorf("invalid delivery status: %s", m.DeliveryStatus)
	}

	// Validate content length based on message type
	if err := m.ValidateContentLength(); err != nil {
		return err
	}

	// Validate metadata based on message type
	if err := m.ValidateMetadata(); err != nil {
		return err
	}

	return nil
}

// ValidateContentLength validates content length based on message type
func (m *Message) ValidateContentLength() error {
	maxLengths := map[MessageType]int{
		MessageTypeText:     4000,  // 4KB for text messages
		MessageTypeSystem:   1000,  // 1KB for system messages
		MessageTypeImage:    200,   // Short description for images
		MessageTypeVideo:    200,   // Short description for videos
		MessageTypeAudio:    200,   // Short description for audio
		MessageTypeFile:     200,   // Short description for files
		MessageTypeLocation: 500,   // Address and description
		MessageTypeSticker:  100,   // Sticker identifier
	}

	if maxLength, exists := maxLengths[m.MessageType]; exists {
		if len(m.Content) > maxLength {
			return fmt.Errorf("content length %d exceeds maximum %d for message type %s",
				len(m.Content), maxLength, m.MessageType)
		}
	}

	return nil
}

// ValidateMetadata validates metadata based on message type
func (m *Message) ValidateMetadata() error {
	switch m.MessageType {
	case MessageTypeFile:
		if m.Metadata.FileName == "" {
			return fmt.Errorf("file name is required for file messages")
		}
		if m.Metadata.FileSize <= 0 {
			return fmt.Errorf("file size must be positive for file messages")
		}
		if m.Metadata.MimeType == "" {
			return fmt.Errorf("mime type is required for file messages")
		}

	case MessageTypeLocation:
		if m.Metadata.Latitude == 0 && m.Metadata.Longitude == 0 {
			return fmt.Errorf("latitude and longitude are required for location messages")
		}
		if m.Metadata.Latitude < -90 || m.Metadata.Latitude > 90 {
			return fmt.Errorf("invalid latitude: %f", m.Metadata.Latitude)
		}
		if m.Metadata.Longitude < -180 || m.Metadata.Longitude > 180 {
			return fmt.Errorf("invalid longitude: %f", m.Metadata.Longitude)
		}

	case MessageTypeSticker:
		if m.Metadata.StickerPack == "" || m.Metadata.StickerID == "" {
			return fmt.Errorf("sticker pack and sticker ID are required for sticker messages")
		}

	case MessageTypeSystem:
		if m.Metadata.SystemType == "" {
			return fmt.Errorf("system type is required for system messages")
		}
	}

	return nil
}

// SetDeliveryStatus updates the delivery status and timestamps
func (m *Message) SetDeliveryStatus(status DeliveryStatus) error {
	if !status.IsValid() {
		return fmt.Errorf("invalid delivery status: %s", status)
	}

	m.DeliveryStatus = status
	m.UpdatedAt = time.Now()

	switch status {
	case DeliveryStatusDelivered:
		if m.DeliveredAt == nil {
			now := time.Now()
			m.DeliveredAt = &now
		}
	case DeliveryStatusRead:
		if m.FirstReadAt == nil {
			now := time.Now()
			m.FirstReadAt = &now
		}
	}

	return nil
}

// AddReaction adds a reaction to the message
func (m *Message) AddReaction(emoji string, userID uuid.UUID) {
	// Find existing reaction for this emoji
	for i, reaction := range m.Reactions {
		if reaction.Emoji == emoji {
			// Check if user already reacted
			for _, existingUserID := range reaction.UserIDs {
				if existingUserID == userID {
					return // User already reacted with this emoji
				}
			}
			// Add user to existing reaction
			m.Reactions[i].UserIDs = append(m.Reactions[i].UserIDs, userID)
			m.UpdatedAt = time.Now()
			return
		}
	}

	// Create new reaction
	m.Reactions = append(m.Reactions, MessageReaction{
		Emoji:   emoji,
		UserIDs: []uuid.UUID{userID},
	})
	m.UpdatedAt = time.Now()
}

// RemoveReaction removes a reaction from the message
func (m *Message) RemoveReaction(emoji string, userID uuid.UUID) {
	for i, reaction := range m.Reactions {
		if reaction.Emoji == emoji {
			// Remove user from reaction
			newUserIDs := []uuid.UUID{}
			for _, existingUserID := range reaction.UserIDs {
				if existingUserID != userID {
					newUserIDs = append(newUserIDs, existingUserID)
				}
			}

			if len(newUserIDs) == 0 {
				// Remove entire reaction if no users left
				m.Reactions = append(m.Reactions[:i], m.Reactions[i+1:]...)
			} else {
				m.Reactions[i].UserIDs = newUserIDs
			}
			m.UpdatedAt = time.Now()
			break
		}
	}
}

// AddReadReceipt adds a read receipt for a user
func (m *Message) AddReadReceipt(userID uuid.UUID) {
	// Check if user already has a read receipt
	for _, receipt := range m.ReadReceipts {
		if receipt.UserID == userID {
			return // Already marked as read
		}
	}

	// Add new read receipt
	m.ReadReceipts = append(m.ReadReceipts, ReadReceipt{
		UserID: userID,
		ReadAt: time.Now(),
	})

	// Update delivery status to read if not already
	if m.DeliveryStatus != DeliveryStatusRead {
		m.SetDeliveryStatus(DeliveryStatusRead)
	}
}

// EditContent edits the message content
func (m *Message) EditContent(newContent string, editorID uuid.UUID) error {
	// Validate that only sender can edit (in most cases)
	if editorID != m.SenderID {
		return fmt.Errorf("only sender can edit message content")
	}

	// System messages cannot be edited
	if m.MessageType == MessageTypeSystem {
		return fmt.Errorf("system messages cannot be edited")
	}

	// Store original content on first edit
	if m.OriginalID == nil {
		originalID := m.ID
		m.OriginalID = &originalID
	}

	m.Content = newContent
	m.ContentLength = len(newContent)
	m.EditCount++
	now := time.Now()
	m.EditedAt = &now
	m.UpdatedAt = now

	return nil
}

// SoftDelete marks the message as deleted
func (m *Message) SoftDelete(deleterID uuid.UUID) {
	now := time.Now()
	m.IsDeleted = true
	m.DeletedAt = &now
	m.DeletedBy = &deleterID
	m.UpdatedAt = now

	// Clear sensitive content but keep metadata for audit
	m.Content = "[deleted]"
	m.ContentLength = 9
}

// IsEdited checks if the message has been edited
func (m *Message) IsEdited() bool {
	return m.EditedAt != nil && m.EditCount > 0
}

// IsRead checks if the message has been read by a specific user
func (m *Message) IsRead(userID uuid.UUID) bool {
	for _, receipt := range m.ReadReceipts {
		if receipt.UserID == userID {
			return true
		}
	}
	return false
}

// GetReadCount returns the number of users who have read the message
func (m *Message) GetReadCount() int {
	return len(m.ReadReceipts)
}

// GetReactionCount returns the total number of reactions
func (m *Message) GetReactionCount() int {
	count := 0
	for _, reaction := range m.Reactions {
		count += len(reaction.UserIDs)
	}
	return count
}

// SetRegionalCompliance sets regional compliance fields
func (m *Message) SetRegionalCompliance(dataRegion, locale string) {
	m.DataRegion = dataRegion
	m.Locale = locale
}

// GenerateSearchKeywords generates search keywords from content
func (m *Message) GenerateSearchKeywords() {
	if m.MessageType != MessageTypeText {
		return
	}

	// Simple keyword extraction (in production, use proper text analysis)
	// This is a simplified version - real implementation would use NLP
	content := strings.ToLower(m.Content)
	words := strings.Fields(content)

	keywords := make(map[string]bool)
	for _, word := range words {
		if len(word) >= 3 { // Only words with 3+ characters
			keywords[word] = true
		}
	}

	m.SearchKeywords = make([]string, 0, len(keywords))
	for keyword := range keywords {
		m.SearchKeywords = append(m.SearchKeywords, keyword)
	}
}

// MarshalJSON customizes JSON serialization
func (m *Message) MarshalJSON() ([]byte, error) {
	type Alias Message
	return json.Marshal(&struct {
		*Alias
		IsEdited      bool `json:"is_edited"`
		ReadCount     int  `json:"read_count"`
		ReactionCount int  `json:"reaction_count"`
	}{
		Alias:         (*Alias)(m),
		IsEdited:      m.IsEdited(),
		ReadCount:     m.GetReadCount(),
		ReactionCount: m.GetReactionCount(),
	})
}

// GetTableName returns the ScyllaDB table name for messages
func (m *Message) GetTableName() string {
	return "messages"
}

// GetCQLCreateTable returns the CQL statement to create the messages table
func GetMessagesCQLCreateTable() string {
	return `
	CREATE TABLE IF NOT EXISTS messages (
		dialog_id UUID,
		created_at TIMESTAMP,
		id UUID,
		sender_id UUID,
		message_type TEXT,
		content TEXT,
		metadata TEXT,
		reply_to_id UUID,
		thread_id UUID,
		is_thread_root BOOLEAN,
		reactions TEXT,
		delivery_status TEXT,
		read_receipts TEXT,
		data_region TEXT,
		locale TEXT,
		edited_at TIMESTAMP,
		original_id UUID,
		edit_count INT,
		is_deleted BOOLEAN,
		deleted_at TIMESTAMP,
		deleted_by UUID,
		updated_at TIMESTAMP,
		processed_at TIMESTAMP,
		delivered_at TIMESTAMP,
		first_read_at TIMESTAMP,
		search_keywords LIST<TEXT>,
		content_hash TEXT,
		content_length INT,
		attachment_count INT,
		PRIMARY KEY (dialog_id, created_at, id)
	) WITH CLUSTERING ORDER BY (created_at DESC)
	  AND compaction = {'class': 'TimeWindowCompactionStrategy', 'compaction_window_size': 1, 'compaction_window_unit': 'DAYS'}
	  AND gc_grace_seconds = 86400;
	`
}

// GetCQLInsert returns the CQL statement to insert a message
func (m *Message) GetCQLInsert() string {
	return `
	INSERT INTO messages (
		dialog_id, created_at, id, sender_id, message_type, content, metadata,
		reply_to_id, thread_id, is_thread_root, reactions, delivery_status,
		read_receipts, data_region, locale, edited_at, original_id, edit_count,
		is_deleted, deleted_at, deleted_by, updated_at, processed_at,
		delivered_at, first_read_at, search_keywords, content_hash,
		content_length, attachment_count
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
}

// GetCQLValues returns the values for CQL insertion
func (m *Message) GetCQLValues() []interface{} {
	// Convert complex types to JSON strings for storage
	metadataJSON, _ := json.Marshal(m.Metadata)
	reactionsJSON, _ := json.Marshal(m.Reactions)
	readReceiptsJSON, _ := json.Marshal(m.ReadReceipts)

	return []interface{}{
		m.DialogID, m.CreatedAt, m.ID, m.SenderID, string(m.MessageType),
		m.Content, string(metadataJSON), m.ReplyToID, m.ThreadID,
		m.IsThreadRoot, string(reactionsJSON), string(m.DeliveryStatus),
		string(readReceiptsJSON), m.DataRegion, m.Locale, m.EditedAt,
		m.OriginalID, m.EditCount, m.IsDeleted, m.DeletedAt, m.DeletedBy,
		m.UpdatedAt, m.ProcessedAt, m.DeliveredAt, m.FirstReadAt,
		m.SearchKeywords, m.ContentHash, m.ContentLength, m.AttachmentCount,
	}
}