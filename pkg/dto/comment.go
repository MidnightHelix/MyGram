package dto

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint64          `json:"id"`
	Message   string          `json:"message"`
	PhotoID   uint64          `json:"photo_id"`
	UserID    uint64          `json:"user_id"`
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`
	User      *UserDefault    `json:"user,omitempty"`
	Photo     *Photo          `json:"photo,omitempty"`
}
