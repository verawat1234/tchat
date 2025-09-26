package fixtures

import (
	"time"

	"github.com/google/uuid"
)

// MessagingFixtures provides test data for Messaging models
type MessagingFixtures struct {
	*BaseFixture
}

// NewMessagingFixtures creates a new messaging fixtures instance
func NewMessagingFixtures(seed ...int64) *MessagingFixtures {
	return &MessagingFixtures{
		BaseFixture: NewBaseFixture(seed...),
	}
}

// MessageType represents different types of messages
type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeVideo    MessageType = "video"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeFile     MessageType = "file"
	MessageTypeSticker  MessageType = "sticker"
	MessageTypeLocation MessageType = "location"
	MessageTypeContact  MessageType = "contact"
	MessageTypePayment  MessageType = "payment"
)

// MessageStatus represents message delivery status
type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
)

// Message represents a test message structure
type Message struct {
	ID          uuid.UUID              `json:"id"`
	ChatID      uuid.UUID              `json:"chat_id"`
	SenderID    uuid.UUID              `json:"sender_id"`
	Type        MessageType            `json:"type"`
	Content     map[string]interface{} `json:"content"`
	Metadata    map[string]interface{} `json:"metadata"`
	Status      MessageStatus          `json:"status"`
	ReplyToID   *uuid.UUID             `json:"reply_to_id,omitempty"`
	EditedAt    *time.Time             `json:"edited_at,omitempty"`
	DeletedAt   *time.Time             `json:"deleted_at,omitempty"`
	DeliveredAt *time.Time             `json:"delivered_at,omitempty"`
	ReadAt      *time.Time             `json:"read_at,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// Chat represents a test chat structure
type Chat struct {
	ID          uuid.UUID              `json:"id"`
	Type        string                 `json:"type"` // direct, group, channel
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Avatar      *string                `json:"avatar,omitempty"`
	Participants []uuid.UUID           `json:"participants"`
	CreatedBy   uuid.UUID              `json:"created_by"`
	Settings    map[string]interface{} `json:"settings"`
	Metadata    map[string]interface{} `json:"metadata"`
	LastMessage *uuid.UUID             `json:"last_message,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// BasicTextMessage creates a basic text message for testing
func (m *MessagingFixtures) BasicTextMessage(chatID, senderID uuid.UUID, country string) *Message {
	content := m.SEAContent(country, "greeting")

	return &Message{
		ID:       m.UUID("text-message-" + senderID.String()),
		ChatID:   chatID,
		SenderID: senderID,
		Type:     MessageTypeText,
		Content: map[string]interface{}{
			"text": content,
		},
		Metadata: map[string]interface{}{
			"language": m.Locale(country),
			"region":   country,
			"length":   len(content),
		},
		Status:      MessageStatusSent,
		ReplyToID:   nil,
		EditedAt:    nil,
		DeletedAt:   nil,
		DeliveredAt: nil,
		ReadAt:      nil,
		CreatedAt:   m.PastTime(30), // Created 30 minutes ago
		UpdatedAt:   m.PastTime(30),
	}
}

// ImageMessage creates an image message for testing
func (m *MessagingFixtures) ImageMessage(chatID, senderID uuid.UUID) *Message {
	return &Message{
		ID:       m.UUID("image-message-" + senderID.String()),
		ChatID:   chatID,
		SenderID: senderID,
		Type:     MessageTypeImage,
		Content: map[string]interface{}{
			"url":        "https://example.com/images/test-image.jpg",
			"thumbnail":  "https://example.com/images/test-image-thumb.jpg",
			"width":      1920,
			"height":     1080,
			"file_size":  1024000,
			"mime_type":  "image/jpeg",
			"alt_text":   "Test image",
		},
		Metadata: map[string]interface{}{
			"upload_source": "mobile",
			"compressed":    true,
			"quality":       85,
		},
		Status:    MessageStatusSent,
		CreatedAt: m.PastTime(45),
		UpdatedAt: m.PastTime(45),
	}
}

// VideoMessage creates a video message for testing
func (m *MessagingFixtures) VideoMessage(chatID, senderID uuid.UUID) *Message {
	return &Message{
		ID:       m.UUID("video-message-" + senderID.String()),
		ChatID:   chatID,
		SenderID: senderID,
		Type:     MessageTypeVideo,
		Content: map[string]interface{}{
			"url":       "https://example.com/videos/test-video.mp4",
			"thumbnail": "https://example.com/videos/test-video-thumb.jpg",
			"duration":  30,    // 30 seconds
			"width":     1280,
			"height":    720,
			"file_size": 5000000, // 5MB
			"mime_type": "video/mp4",
		},
		Metadata: map[string]interface{}{
			"upload_source": "mobile",
			"compressed":    true,
			"quality":       "720p",
		},
		Status:    MessageStatusSent,
		CreatedAt: m.PastTime(60),
		UpdatedAt: m.PastTime(60),
	}
}

// AudioMessage creates an audio message for testing
func (m *MessagingFixtures) AudioMessage(chatID, senderID uuid.UUID) *Message {
	return &Message{
		ID:       m.UUID("audio-message-" + senderID.String()),
		ChatID:   chatID,
		SenderID: senderID,
		Type:     MessageTypeAudio,
		Content: map[string]interface{}{
			"url":       "https://example.com/audio/test-audio.mp3",
			"duration":  15, // 15 seconds
			"file_size": 500000, // 500KB
			"mime_type": "audio/mpeg",
			"waveform":  []int{10, 20, 15, 25, 30, 20, 15, 10}, // Audio waveform data
		},
		Metadata: map[string]interface{}{
			"upload_source": "mobile",
			"bitrate":       128,
			"sample_rate":   44100,
		},
		Status:    MessageStatusSent,
		CreatedAt: m.PastTime(20),
		UpdatedAt: m.PastTime(20),
	}
}

// PaymentMessage creates a payment message for testing
func (m *MessagingFixtures) PaymentMessage(chatID, senderID uuid.UUID, currency string) *Message {
	amount := m.Amount(currency)

	return &Message{
		ID:       m.UUID("payment-message-" + senderID.String()),
		ChatID:   chatID,
		SenderID: senderID,
		Type:     MessageTypePayment,
		Content: map[string]interface{}{
			"amount":       amount,
			"currency":     currency,
			"description":  "Test payment",
			"payment_id":   m.UUID("payment-" + senderID.String()).String(),
			"status":       "completed",
		},
		Metadata: map[string]interface{}{
			"payment_method": "wallet",
			"fee":            amount / 100, // 1% fee
			"net_amount":     amount - (amount / 100),
		},
		Status:    MessageStatusSent,
		CreatedAt: m.PastTime(10),
		UpdatedAt: m.PastTime(10),
	}
}

// ReplyMessage creates a reply message for testing
func (m *MessagingFixtures) ReplyMessage(chatID, senderID, replyToID uuid.UUID, country string) *Message {
	message := m.BasicTextMessage(chatID, senderID, country)
	message.ID = m.UUID("reply-message-" + senderID.String())
	message.ReplyToID = &replyToID
	message.Content["text"] = "Reply: " + message.Content["text"].(string)
	message.Metadata["is_reply"] = true
	return message
}

// EditedMessage creates an edited message for testing
func (m *MessagingFixtures) EditedMessage(chatID, senderID uuid.UUID, country string) *Message {
	message := m.BasicTextMessage(chatID, senderID, country)
	message.ID = m.UUID("edited-message-" + senderID.String())
	message.Content["text"] = "Edited: " + message.Content["text"].(string)
	editedAt := m.PastTime(15)
	message.EditedAt = &editedAt
	message.UpdatedAt = editedAt
	message.Metadata["edit_count"] = 1
	return message
}

// DeletedMessage creates a deleted message for testing
func (m *MessagingFixtures) DeletedMessage(chatID, senderID uuid.UUID) *Message {
	message := m.BasicTextMessage(chatID, senderID, "TH")
	message.ID = m.UUID("deleted-message-" + senderID.String())
	deletedAt := m.PastTime(5)
	message.DeletedAt = &deletedAt
	message.Content = map[string]interface{}{
		"text": "[Message deleted]",
	}
	message.Metadata["deleted_by"] = senderID.String()
	return message
}

// DirectChat creates a direct chat for testing
func (m *MessagingFixtures) DirectChat(user1ID, user2ID uuid.UUID) *Chat {
	return &Chat{
		ID:           m.UUID("direct-chat-" + user1ID.String() + "-" + user2ID.String()),
		Type:         "direct",
		Name:         nil,
		Description:  nil,
		Avatar:       nil,
		Participants: []uuid.UUID{user1ID, user2ID},
		CreatedBy:    user1ID,
		Settings: map[string]interface{}{
			"notifications": true,
			"read_receipts": true,
			"encryption":    true,
		},
		Metadata: map[string]interface{}{
			"chat_type":     "direct",
			"participant_count": 2,
		},
		LastMessage: nil,
		CreatedAt:   m.PastTime(1440), // Created yesterday
		UpdatedAt:   m.PastTime(30),   // Updated 30 minutes ago
	}
}

// GroupChat creates a group chat for testing
func (m *MessagingFixtures) GroupChat(creatorID uuid.UUID, participantIDs []uuid.UUID, country string) *Chat {
	name := "Test Group - " + m.Name(country, "male")
	description := "Test group chat for " + country

	allParticipants := append([]uuid.UUID{creatorID}, participantIDs...)

	return &Chat{
		ID:           m.UUID("group-chat-" + creatorID.String()),
		Type:         "group",
		Name:         &name,
		Description:  &description,
		Avatar:       nil,
		Participants: allParticipants,
		CreatedBy:    creatorID,
		Settings: map[string]interface{}{
			"notifications":    true,
			"read_receipts":    true,
			"encryption":       true,
			"admin_only_post":  false,
			"invite_link":      true,
			"max_participants": 100,
		},
		Metadata: map[string]interface{}{
			"chat_type":        "group",
			"participant_count": len(allParticipants),
			"region":           country,
		},
		LastMessage: nil,
		CreatedAt:   m.PastTime(720), // Created 12 hours ago
		UpdatedAt:   m.PastTime(60),  // Updated 1 hour ago
	}
}

// ChannelChat creates a channel chat for testing
func (m *MessagingFixtures) ChannelChat(creatorID uuid.UUID, country string) *Chat {
	name := "Test Channel - " + country
	description := "Official test channel for " + country + " region"
	avatar := "https://example.com/channels/test-channel-avatar.jpg"

	return &Chat{
		ID:           m.UUID("channel-chat-" + creatorID.String()),
		Type:         "channel",
		Name:         &name,
		Description:  &description,
		Avatar:       &avatar,
		Participants: []uuid.UUID{creatorID}, // Only creator initially
		CreatedBy:    creatorID,
		Settings: map[string]interface{}{
			"notifications":     true,
			"read_receipts":     false, // Channels typically don't show read receipts
			"encryption":        false, // Public channels might not be encrypted
			"admin_only_post":   true,
			"public":            true,
			"invite_link":       true,
			"max_participants":  10000,
		},
		Metadata: map[string]interface{}{
			"chat_type":        "channel",
			"participant_count": 1,
			"region":           country,
			"category":         "official",
		},
		LastMessage: nil,
		CreatedAt:   m.PastTime(2880), // Created 2 days ago
		UpdatedAt:   m.PastTime(120),  // Updated 2 hours ago
	}
}

// MessageThread creates a thread of messages for testing
func (m *MessagingFixtures) MessageThread(chatID uuid.UUID, participantIDs []uuid.UUID, messageCount int, country string) []*Message {
	messages := make([]*Message, 0, messageCount)

	for i := 0; i < messageCount; i++ {
		senderID := participantIDs[i%len(participantIDs)]

		var message *Message

		// Vary message types
		switch i % 5 {
		case 0:
			message = m.BasicTextMessage(chatID, senderID, country)
		case 1:
			message = m.ImageMessage(chatID, senderID)
		case 2:
			message = m.AudioMessage(chatID, senderID)
		case 3:
			message = m.PaymentMessage(chatID, senderID, m.Currency(country))
		case 4:
			if len(messages) > 0 {
				message = m.ReplyMessage(chatID, senderID, messages[len(messages)-1].ID, country)
			} else {
				message = m.BasicTextMessage(chatID, senderID, country)
			}
		}

		// Adjust timestamps for chronological order
		message.ID = m.UUID("thread-message-" + string(rune(i)) + "-" + senderID.String())
		message.CreatedAt = m.PastTime(messageCount*5 - i*5) // 5 minutes apart
		message.UpdatedAt = message.CreatedAt

		// Simulate message status progression
		if i < messageCount-2 {
			message.Status = MessageStatusRead
			readAt := message.CreatedAt.Add(time.Minute)
			message.ReadAt = &readAt
			deliveredAt := message.CreatedAt.Add(30 * time.Second)
			message.DeliveredAt = &deliveredAt
		} else if i < messageCount-1 {
			message.Status = MessageStatusDelivered
			deliveredAt := message.CreatedAt.Add(30 * time.Second)
			message.DeliveredAt = &deliveredAt
		}

		messages = append(messages, message)
	}

	return messages
}

// TestMessagingData creates a comprehensive set of messaging test data
func (m *MessagingFixtures) TestMessagingData(userIDs []uuid.UUID, country string) map[string]interface{} {
	if len(userIDs) < 3 {
		// Add more users if needed
		for len(userIDs) < 3 {
			userIDs = append(userIDs, m.UUID("additional-user-"+string(rune(len(userIDs)))))
		}
	}

	// Create different types of chats
	directChat := m.DirectChat(userIDs[0], userIDs[1])
	groupChat := m.GroupChat(userIDs[0], userIDs[1:], country)
	channelChat := m.ChannelChat(userIDs[0], country)

	chats := []*Chat{directChat, groupChat, channelChat}

	// Create message threads for each chat
	allMessages := make([]*Message, 0)

	// Direct chat messages
	directMessages := m.MessageThread(directChat.ID, []uuid.UUID{userIDs[0], userIDs[1]}, 10, country)
	allMessages = append(allMessages, directMessages...)

	// Group chat messages
	groupMessages := m.MessageThread(groupChat.ID, userIDs, 15, country)
	allMessages = append(allMessages, groupMessages...)

	// Channel messages (admin only)
	channelMessages := m.MessageThread(channelChat.ID, []uuid.UUID{userIDs[0]}, 5, country)
	allMessages = append(allMessages, channelMessages...)

	// Add some special message types
	editedMessage := m.EditedMessage(directChat.ID, userIDs[0], country)
	deletedMessage := m.DeletedMessage(groupChat.ID, userIDs[1])
	allMessages = append(allMessages, editedMessage, deletedMessage)

	return map[string]interface{}{
		"chats":    chats,
		"messages": allMessages,
	}
}

// AllMessagingFixtures creates a complete set of messaging-related test data
func AllMessagingFixtures(seed ...int64) *MessagingFixtures {
	return NewMessagingFixtures(seed...)
}