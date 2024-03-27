package model

import (
	"time"

	"gorm.io/gorm"
)

type SocialMedia struct {
	ID        uint64 `json:"id" gorm:"primaryKey"`
	UserID    uint64 `json:"user_id"`
	Name      string `json:"name" gorm:"not null" binding:"required" validate:"required"`
	Url       string `json:"social_media_url" gorm:"not null" binding:"required" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
	User      *User          `json:"user,omitempty" validate:"-"`
}
