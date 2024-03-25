package dto

import (
	"time"

	"gorm.io/gorm"
)

type Photo struct {
	ID        uint64          `json:"id"`
	Title     string          `json:"title"`
	Caption   string          `json:"caption"`
	Url       string          `json:"photo_url"`
	UserID    uint64          `json:"user_id"`
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`
	Comments  []Comment       `json:"comments,omitempty"`
	User      *UserDefault    `json:"user,omitempty"`
}
