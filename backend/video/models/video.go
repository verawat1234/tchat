package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Video struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string    `json:"title" gorm:"not null" validate:"required"`
	Description string    `json:"description"`
	ThumbnailURL string   `json:"thumbnail" gorm:"column:thumbnail_url"`
	VideoURL    string    `json:"videoUrl" gorm:"column:video_url;not null" validate:"required"`
	Duration    string    `json:"duration"`
	Views       int64     `json:"views" gorm:"default:0"`
	Likes       int64     `json:"likes" gorm:"default:0"`
	Category    string    `json:"category" validate:"required"`
	Tags        []string  `json:"tags" gorm:"type:text[]"`
	Type        string    `json:"type" gorm:"default:'short'" validate:"required"`
	Status      string    `json:"status" gorm:"default:'active'" validate:"required"`
	ChannelID   uuid.UUID `json:"channelId" gorm:"type:uuid;not null"`
	Channel     Channel   `json:"channel" gorm:"foreignKey:ChannelID"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type Channel struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null" validate:"required"`
	Avatar      string    `json:"avatar"`
	Subscribers int64     `json:"subscribers" gorm:"default:0"`
	Verified    bool      `json:"verified" gorm:"default:false"`
	UserID      uuid.UUID `json:"userId" gorm:"type:uuid"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type VideoInteraction struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	VideoID   uuid.UUID `json:"videoId" gorm:"type:uuid;not null"`
	UserID    uuid.UUID `json:"userId" gorm:"type:uuid;not null"`
	Type      string    `json:"type" validate:"required,oneof=like dislike view share"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

// TableName specifies the table name for Video
func (Video) TableName() string {
	return "videos"
}

// TableName specifies the table name for Channel
func (Channel) TableName() string {
	return "channels"
}

// TableName specifies the table name for VideoInteraction
func (VideoInteraction) TableName() string {
	return "video_interactions"
}

// BeforeCreate hook to generate UUID if not provided
func (v *Video) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook to generate UUID if not provided
func (c *Channel) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// VideoComment represents a comment on a video
type VideoComment struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	VideoID   uuid.UUID  `json:"videoId" gorm:"type:uuid;not null"`
	Video     Video      `json:"video" gorm:"foreignKey:VideoID"`
	UserID    uuid.UUID  `json:"userId" gorm:"type:uuid;not null"`
	UserName  string     `json:"userName" gorm:"not null"`
	UserAvatar string    `json:"userAvatar"`
	Content   string     `json:"content" gorm:"not null" validate:"required"`
	ParentID  *uuid.UUID `json:"parentId,omitempty" gorm:"type:uuid"`
	Parent    *VideoComment `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies   []VideoComment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
	Likes     int64      `json:"likes" gorm:"default:0"`
	IsEdited  bool       `json:"isEdited" gorm:"default:false"`
	CreatedAt time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
}

// VideoShare represents a video share action
type VideoShare struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	VideoID   uuid.UUID `json:"videoId" gorm:"type:uuid;not null"`
	Video     Video     `json:"video" gorm:"foreignKey:VideoID"`
	UserID    uuid.UUID `json:"userId" gorm:"type:uuid;not null"`
	Platform  string    `json:"platform" gorm:"not null" validate:"required"`
	Message   string    `json:"message"`
	ShareURL  string    `json:"shareUrl"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

// TableName specifies the table name for VideoComment
func (VideoComment) TableName() string {
	return "video_comments"
}

// TableName specifies the table name for VideoShare
func (VideoShare) TableName() string {
	return "video_shares"
}

// BeforeCreate hook to generate UUID if not provided
func (vi *VideoInteraction) BeforeCreate(tx *gorm.DB) error {
	if vi.ID == uuid.Nil {
		vi.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook to generate UUID if not provided
func (vc *VideoComment) BeforeCreate(tx *gorm.DB) error {
	if vc.ID == uuid.Nil {
		vc.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook to generate UUID if not provided
func (vs *VideoShare) BeforeCreate(tx *gorm.DB) error {
	if vs.ID == uuid.Nil {
		vs.ID = uuid.New()
	}
	return nil
}