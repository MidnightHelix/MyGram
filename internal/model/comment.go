package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint64 `json:"id" gorm:"primaryKey"`
	UserID    uint64 `json:"user_id" gorm:"not null"`
	PhotoID   uint64 `json:"photo_id" gorm:"not null"`
	Message   string `json:"message" gorm:"not null" binding:"required" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
	User      *User          `json:"user,omitempty" validate:"-"`
	Photo     *Photo         `json:"photo,omitempty" validate:"-"`
}
