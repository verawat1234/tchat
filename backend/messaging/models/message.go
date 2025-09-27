package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Message represents a message in a dialog/conversation
type Message struct {
	ID           uuid.UUID         `json:"id" db:"id"`
	DialogID     uuid.UUID         `json:"dialog_id" db:"dialog_id"`
	SenderID     uuid.UUID         `json:"sender_id" db:"sender_id"`
	Type         MessageType       `json:"type" db:"type"`
	Content      MessageContent    `json:"content" db:"content"`
	MediaURL     *string           `json:"media_url,omitempty" db:"media_url"`
	ThumbnailURL *string           `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	Status       MessageStatus     `json:"status" db:"status"`
	Metadata     map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	ReplyToID    *uuid.UUID        `json:"reply_to_id,omitempty" db:"reply_to_id"`
	ReplyTo      *MessageReply     `json:"reply_to,omitempty" db:"reply_to"`
	IsEdited     bool              `json:"is_edited" db:"is_edited"`
	IsPinned     bool              `json:"is_pinned" db:"is_pinned"`
	IsDeleted    bool              `json:"is_deleted" db:"is_deleted"`
	Mentions     UUIDSlice         `json:"mentions,omitempty" db:"mentions"`
	Reactions    MessageReactions  `json:"reactions,omitempty" db:"reactions"`
	SentAt       time.Time         `json:"sent_at" db:"sent_at"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at" db:"updated_at"`
	EditedAt     *time.Time        `json:"edited_at,omitempty" db:"edited_at"`
	DeletedAt    *time.Time        `json:"deleted_at,omitempty" db:"deleted_at"`
}

// MessageType represents the type of message content
type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeVoice    MessageType = "voice"
	MessageTypeFile     MessageType = "file"
	MessageTypeImage    MessageType = "image"
	MessageTypeVideo    MessageType = "video"
	MessageTypePayment  MessageType = "payment"
	MessageTypeLocation MessageType = "location"
	MessageTypeSticker  MessageType = "sticker"
	MessageTypeSystem   MessageType = "system"
)

// MessageStatus represents the status of a message
type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
)

// MessageReply represents a reply reference in a message
type MessageReply struct {
	MessageID   uuid.UUID `json:"message_id"`
	SenderID    uuid.UUID `json:"sender_id"`
	Content     string    `json:"content"`
	MessageType string    `json:"message_type"`
}

// MessageForward represents a forward reference in a message
type MessageForward struct {
	OriginalMessageID uuid.UUID `json:"original_message_id"`
	OriginalSenderID  uuid.UUID `json:"original_sender_id"`
	OriginalDialogID  uuid.UUID `json:"original_dialog_id"`
	ForwardedAt       time.Time `json:"forwarded_at"`
}

// MessageContent represents the content of a message (varies by type)
type MessageContent map[string]interface{}

// MessageReactions represents reactions to a message
type MessageReactions map[string][]uuid.UUID // emoji -> user IDs

// Text message content structure
type TextContent struct {
	Text     string                 `json:"text"`
	Entities []MessageEntity        `json:"entities,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Voice message content structure
type VoiceContent struct {
	URL      string  `json:"url"`
	Duration int     `json:"duration"` // in milliseconds
	Waveform []int   `json:"waveform,omitempty"`
	FileSize int64   `json:"file_size"`
	MimeType string  `json:"mime_type"`
}

// File message content structure
type FileContent struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type"`
	Caption  string `json:"caption,omitempty"`
}

// Image message content structure
type ImageContent struct {
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	FileSize  int64  `json:"file_size"`
	Caption   string `json:"caption,omitempty"`
}

// Video message content structure
type VideoContent struct {
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Duration  int    `json:"duration"` // in milliseconds
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	FileSize  int64  `json:"file_size"`
	Caption   string `json:"caption,omitempty"`
}

// Payment message content structure
type PaymentContent struct {
	Amount      int64  `json:"amount"`      // in cents
	Currency    string `json:"currency"`
	Description string `json:"description"`
	Reference   string `json:"reference,omitempty"`
	Status      string `json:"status"`      // pending, completed, failed, cancelled
}

// Location message content structure
type LocationContent struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address,omitempty"`
	Venue     string  `json:"venue,omitempty"`
}

// Sticker message content structure
type StickerContent struct {
	StickerID string `json:"sticker_id"`
	PackID    string `json:"pack_id"`
	URL       string `json:"url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

// System message content structure
type SystemContent struct {
	Type    string                 `json:"type"`    // user_joined, user_left, name_changed, etc.
	Message string                 `json:"message"` // Human-readable message
	Data    map[string]interface{} `json:"data"`    // Additional system data
}

// MessageEntity represents entities within text messages (mentions, links, etc.)
type MessageEntity struct {
	Type   string `json:"type"`   // mention, hashtag, url, email, phone, bold, italic, code
	Offset int    `json:"offset"` // Start position in text
	Length int    `json:"length"` // Length of entity
	URL    string `json:"url,omitempty"`    // For url entities
	UserID string `json:"user_id,omitempty"` // For mention entities
}

// Value implements the driver.Valuer interface for MessageContent
func (mc MessageContent) Value() (driver.Value, error) {
	return json.Marshal(mc)
}

// Scan implements the sql.Scanner interface for MessageContent
func (mc *MessageContent) Scan(value interface{}) error {
	if value == nil {
		*mc = nil
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into MessageContent", value)
	}

	return json.Unmarshal(jsonData, mc)
}

// Value implements the driver.Valuer interface for MessageReactions
func (mr MessageReactions) Value() (driver.Value, error) {
	if mr == nil {
		return nil, nil
	}
	return json.Marshal(mr)
}

// Scan implements the sql.Scanner interface for MessageReactions
func (mr *MessageReactions) Scan(value interface{}) error {
	if value == nil {
		*mr = nil
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into MessageReactions", value)
	}

	return json.Unmarshal(jsonData, mr)
}

// ValidMessageTypes returns all supported message types
func ValidMessageTypes() []MessageType {
	return []MessageType{
		MessageTypeText,
		MessageTypeVoice,
		MessageTypeFile,
		MessageTypeImage,
		MessageTypeVideo,
		MessageTypePayment,
		MessageTypeLocation,
		MessageTypeSticker,
		MessageTypeSystem,
	}
}

// IsValid validates if the message type is supported
func (mt MessageType) IsValid() bool {
	for _, valid := range ValidMessageTypes() {
		if mt == valid {
			return true
		}
	}
	return false
}

// String returns the string representation of MessageType
func (mt MessageType) String() string {
	return string(mt)
}

// RequiresContent checks if this message type requires content
func (mt MessageType) RequiresContent() bool {
	return mt != MessageTypeSystem // System messages can have empty content
}

// IsMedia checks if this is a media message type
func (mt MessageType) IsMedia() bool {
	return mt == MessageTypeVoice || mt == MessageTypeFile ||
		   mt == MessageTypeImage || mt == MessageTypeVideo
}

// Validate performs comprehensive validation on the Message model
func (m *Message) Validate() error {
	var errs []string

	// Dialog ID validation
	if m.DialogID == uuid.Nil {
		errs = append(errs, "dialog_id is required")
	}

	// Sender ID validation
	if m.SenderID == uuid.Nil {
		errs = append(errs, "sender_id is required")
	}

	// Message type validation
	if !m.Type.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid message type: %s", m.Type))
	}

	// Content validation
	if m.Type.RequiresContent() && len(m.Content) == 0 {
		errs = append(errs, "content is required for message type: "+string(m.Type))
	}

	// Type-specific content validation
	if err := m.validateContentForType(); err != nil {
		errs = append(errs, err.Error())
	}

	// Reply validation
	if m.ReplyToID != nil && *m.ReplyToID == uuid.Nil {
		errs = append(errs, "reply_to_id cannot be nil UUID")
	}

	// Self-reply validation
	if m.ReplyToID != nil && *m.ReplyToID == m.ID {
		errs = append(errs, "message cannot reply to itself")
	}

	// Mentions validation
	if err := m.validateMentions(); err != nil {
		errs = append(errs, err.Error())
	}

	// Reactions validation
	if err := m.validateReactions(); err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

// validateContentForType validates content based on message type
func (m *Message) validateContentForType() error {
	switch m.Type {
	case MessageTypeText:
		return m.validateTextContent()
	case MessageTypeVoice:
		return m.validateVoiceContent()
	case MessageTypeFile:
		return m.validateFileContent()
	case MessageTypeImage:
		return m.validateImageContent()
	case MessageTypeVideo:
		return m.validateVideoContent()
	case MessageTypePayment:
		return m.validatePaymentContent()
	case MessageTypeLocation:
		return m.validateLocationContent()
	case MessageTypeSticker:
		return m.validateStickerContent()
	case MessageTypeSystem:
		return m.validateSystemContent()
	default:
		return fmt.Errorf("unknown message type: %s", m.Type)
	}
}

// validateTextContent validates text message content
func (m *Message) validateTextContent() error {
	text, exists := m.Content["text"]
	if !exists {
		return errors.New("text content must have 'text' field")
	}

	textStr, ok := text.(string)
	if !ok {
		return errors.New("text field must be a string")
	}

	if strings.TrimSpace(textStr) == "" {
		return errors.New("text content cannot be empty")
	}

	if len(textStr) > 4096 {
		return errors.New("text content cannot exceed 4096 characters")
	}

	return nil
}

// validateVoiceContent validates voice message content
func (m *Message) validateVoiceContent() error {
	requiredFields := []string{"url", "duration", "file_size"}
	for _, field := range requiredFields {
		if _, exists := m.Content[field]; !exists {
			return fmt.Errorf("voice content must have '%s' field", field)
		}
	}

	// Validate duration
	if duration, ok := m.Content["duration"].(float64); ok {
		if duration <= 0 || duration > 300000 { // Max 5 minutes
			return errors.New("voice duration must be between 1ms and 5 minutes")
		}
	}

	// Validate file size
	if fileSize, ok := m.Content["file_size"].(float64); ok {
		if fileSize <= 0 || fileSize > 50*1024*1024 { // Max 50MB
			return errors.New("voice file size must be between 1 byte and 50MB")
		}
	}

	return nil
}

// validateFileContent validates file message content
func (m *Message) validateFileContent() error {
	requiredFields := []string{"url", "filename", "file_size", "mime_type"}
	for _, field := range requiredFields {
		if _, exists := m.Content[field]; !exists {
			return fmt.Errorf("file content must have '%s' field", field)
		}
	}

	// Validate file size
	if fileSize, ok := m.Content["file_size"].(float64); ok {
		if fileSize <= 0 || fileSize > 100*1024*1024 { // Max 100MB
			return errors.New("file size must be between 1 byte and 100MB")
		}
	}

	return nil
}

// validateImageContent validates image message content
func (m *Message) validateImageContent() error {
	requiredFields := []string{"url", "width", "height", "file_size"}
	for _, field := range requiredFields {
		if _, exists := m.Content[field]; !exists {
			return fmt.Errorf("image content must have '%s' field", field)
		}
	}

	// Validate dimensions
	if width, ok := m.Content["width"].(float64); ok {
		if width <= 0 || width > 8192 {
			return errors.New("image width must be between 1 and 8192 pixels")
		}
	}

	if height, ok := m.Content["height"].(float64); ok {
		if height <= 0 || height > 8192 {
			return errors.New("image height must be between 1 and 8192 pixels")
		}
	}

	return nil
}

// validateVideoContent validates video message content
func (m *Message) validateVideoContent() error {
	requiredFields := []string{"url", "duration", "width", "height", "file_size"}
	for _, field := range requiredFields {
		if _, exists := m.Content[field]; !exists {
			return fmt.Errorf("video content must have '%s' field", field)
		}
	}

	// Validate duration (max 30 minutes)
	if duration, ok := m.Content["duration"].(float64); ok {
		if duration <= 0 || duration > 1800000 {
			return errors.New("video duration must be between 1ms and 30 minutes")
		}
	}

	// Validate file size (max 500MB)
	if fileSize, ok := m.Content["file_size"].(float64); ok {
		if fileSize <= 0 || fileSize > 500*1024*1024 {
			return errors.New("video file size must be between 1 byte and 500MB")
		}
	}

	return nil
}

// validatePaymentContent validates payment message content
func (m *Message) validatePaymentContent() error {
	requiredFields := []string{"amount", "currency", "description", "status"}
	for _, field := range requiredFields {
		if _, exists := m.Content[field]; !exists {
			return fmt.Errorf("payment content must have '%s' field", field)
		}
	}

	// Validate amount (positive, non-zero)
	if amount, ok := m.Content["amount"].(float64); ok {
		if amount <= 0 {
			return errors.New("payment amount must be positive")
		}
	}

	// Validate currency
	validCurrencies := []string{"THB", "SGD", "IDR", "MYR", "PHP", "VND", "USD"}
	if currency, ok := m.Content["currency"].(string); ok {
		valid := false
		for _, validCurrency := range validCurrencies {
			if currency == validCurrency {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid currency: %s", currency)
		}
	}

	// Validate status
	validStatuses := []string{"pending", "completed", "failed", "cancelled"}
	if status, ok := m.Content["status"].(string); ok {
		valid := false
		for _, validStatus := range validStatuses {
			if status == validStatus {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid payment status: %s", status)
		}
	}

	return nil
}

// validateLocationContent validates location message content
func (m *Message) validateLocationContent() error {
	requiredFields := []string{"latitude", "longitude"}
	for _, field := range requiredFields {
		if _, exists := m.Content[field]; !exists {
			return fmt.Errorf("location content must have '%s' field", field)
		}
	}

	// Validate latitude (-90 to 90)
	if lat, ok := m.Content["latitude"].(float64); ok {
		if lat < -90 || lat > 90 {
			return errors.New("latitude must be between -90 and 90")
		}
	}

	// Validate longitude (-180 to 180)
	if lng, ok := m.Content["longitude"].(float64); ok {
		if lng < -180 || lng > 180 {
			return errors.New("longitude must be between -180 and 180")
		}
	}

	return nil
}

// validateStickerContent validates sticker message content
func (m *Message) validateStickerContent() error {
	requiredFields := []string{"sticker_id", "pack_id", "url"}
	for _, field := range requiredFields {
		if _, exists := m.Content[field]; !exists {
			return fmt.Errorf("sticker content must have '%s' field", field)
		}
	}

	return nil
}

// validateSystemContent validates system message content
func (m *Message) validateSystemContent() error {
	requiredFields := []string{"type", "message"}
	for _, field := range requiredFields {
		if _, exists := m.Content[field]; !exists {
			return fmt.Errorf("system content must have '%s' field", field)
		}
	}

	return nil
}

// validateMentions validates message mentions
func (m *Message) validateMentions() error {
	if len(m.Mentions) == 0 {
		return nil
	}

	// Validate unique mentions
	seen := make(map[uuid.UUID]bool)
	for _, mention := range m.Mentions {
		if mention == uuid.Nil {
			return errors.New("mention ID cannot be nil")
		}
		if seen[mention] {
			return fmt.Errorf("duplicate mention: %s", mention)
		}
		seen[mention] = true
	}

	// Validate mention limit
	if len(m.Mentions) > 50 {
		return errors.New("too many mentions (max 50)")
	}

	return nil
}

// validateReactions validates message reactions
func (m *Message) validateReactions() error {
	if len(m.Reactions) == 0 {
		return nil
	}

	// Validate reaction limit
	if len(m.Reactions) > 20 {
		return errors.New("too many reaction types (max 20)")
	}

	// Validate each reaction
	for emoji, userIDs := range m.Reactions {
		if strings.TrimSpace(emoji) == "" {
			return errors.New("reaction emoji cannot be empty")
		}

		if len(userIDs) == 0 {
			return fmt.Errorf("reaction '%s' has no users", emoji)
		}

		if len(userIDs) > 1000 {
			return fmt.Errorf("too many users for reaction '%s' (max 1000)", emoji)
		}

		// Validate unique user IDs per reaction
		seen := make(map[uuid.UUID]bool)
		for _, userID := range userIDs {
			if userID == uuid.Nil {
				return fmt.Errorf("reaction '%s' contains nil user ID", emoji)
			}
			if seen[userID] {
				return fmt.Errorf("duplicate user in reaction '%s': %s", emoji, userID)
			}
			seen[userID] = true
		}
	}

	return nil
}

// BeforeCreate sets up the message before database creation
func (m *Message) BeforeCreate() error {
	// Generate UUID if not set
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}

	// Set creation timestamp
	m.CreatedAt = time.Now().UTC()

	// Initialize empty maps if nil
	if m.Content == nil {
		m.Content = make(MessageContent)
	}
	if m.Reactions == nil {
		m.Reactions = make(MessageReactions)
	}

	// Validate before creation
	return m.Validate()
}

// BeforeUpdate sets up the message before database update
func (m *Message) BeforeUpdate() error {
	// Set edited timestamp if content changed
	if m.IsEdited {
		now := time.Now().UTC()
		m.EditedAt = &now
	}

	// Validate before update
	return m.Validate()
}

// SoftDelete marks the message as deleted without removing it
func (m *Message) SoftDelete() error {
	if m.IsDeleted {
		return errors.New("message is already deleted")
	}

	now := time.Now().UTC()
	m.IsDeleted = true
	m.DeletedAt = &now

	return nil
}

// AddReaction adds a reaction from a user
func (m *Message) AddReaction(emoji string, userID uuid.UUID) error {
	if strings.TrimSpace(emoji) == "" {
		return errors.New("emoji cannot be empty")
	}

	if userID == uuid.Nil {
		return errors.New("user ID cannot be nil")
	}

	if m.Reactions == nil {
		m.Reactions = make(MessageReactions)
	}

	// Check if user already reacted with this emoji
	for _, existingUserID := range m.Reactions[emoji] {
		if existingUserID == userID {
			return errors.New("user has already reacted with this emoji")
		}
	}

	m.Reactions[emoji] = append(m.Reactions[emoji], userID)
	return nil
}

// RemoveReaction removes a reaction from a user
func (m *Message) RemoveReaction(emoji string, userID uuid.UUID) error {
	if m.Reactions == nil {
		return errors.New("no reactions found")
	}

	userIDs, exists := m.Reactions[emoji]
	if !exists {
		return errors.New("reaction not found")
	}

	// Find and remove user ID
	for i, existingUserID := range userIDs {
		if existingUserID == userID {
			m.Reactions[emoji] = append(userIDs[:i], userIDs[i+1:]...)

			// Remove emoji entry if no users left
			if len(m.Reactions[emoji]) == 0 {
				delete(m.Reactions, emoji)
			}

			return nil
		}
	}

	return errors.New("user has not reacted with this emoji")
}

// GetReactionCount returns the total number of reactions
func (m *Message) GetReactionCount() int {
	count := 0
	for _, userIDs := range m.Reactions {
		count += len(userIDs)
	}
	return count
}

// HasUserReacted checks if a user has reacted to this message
func (m *Message) HasUserReacted(userID uuid.UUID) bool {
	for _, userIDs := range m.Reactions {
		for _, existingUserID := range userIDs {
			if existingUserID == userID {
				return true
			}
		}
	}
	return false
}

// IsMentioned checks if a user is mentioned in this message
func (m *Message) IsMentioned(userID uuid.UUID) bool {
	for _, mention := range m.Mentions {
		if mention == userID {
			return true
		}
	}
	return false
}

// CanEdit checks if the message can be edited
func (m *Message) CanEdit() bool {
	// System messages cannot be edited
	if m.Type == MessageTypeSystem {
		return false
	}

	// Deleted messages cannot be edited
	if m.IsDeleted {
		return false
	}

	// Only text messages can be edited for now
	if m.Type != MessageTypeText {
		return false
	}

	return true
}

// ToPublicMessage returns a sanitized version for public API responses
func (m *Message) ToPublicMessage(forUserID uuid.UUID) map[string]interface{} {
	response := map[string]interface{}{
		"id":         m.ID,
		"dialog_id":  m.DialogID,
		"sender_id":  m.SenderID,
		"type":       m.Type,
		"is_edited":  m.IsEdited,
		"is_pinned":  m.IsPinned,
		"created_at": m.CreatedAt,
	}

	// Include content if not deleted or if user is sender
	if !m.IsDeleted || m.SenderID == forUserID {
		response["content"] = m.Content
		response["mentions"] = m.Mentions
		response["reactions"] = m.Reactions

		if m.ReplyToID != nil {
			response["reply_to_id"] = m.ReplyToID
		}

		if m.IsEdited && m.EditedAt != nil {
			response["edited_at"] = m.EditedAt
		}
	} else {
		// Show deleted message placeholder
		response["content"] = map[string]interface{}{
			"text": "This message was deleted",
		}
		response["is_deleted"] = true
	}

	return response
}