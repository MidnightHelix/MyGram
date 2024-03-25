package dto

import (
	"time"

	"gorm.io/gorm"
)

type SocialMedia struct {
	ID        uint64          `json:"id"`
	Name      string          `json:"name"`
	Url       string          `json:"social_media_url"`
	UserID    uint64          `json:"user_id"`
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`
	User      *UserDefault    `json:"user,omitempty"`
}
