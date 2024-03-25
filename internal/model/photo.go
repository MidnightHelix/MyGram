package model

import (
	"time"

	"gorm.io/gorm"
)

type Photo struct {
	ID        uint64 `json:"id" gorm:"primaryKey"`
	Title     string `json:"title" gorm:"not null"`
	Caption   string `json:"caption"`
	Url       string `json:"photo_url" gorm:"not null"`
	UserID    uint64 `json:"user_id" gorm:"column:user_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
	Comments  []Comment      `json:"comments,omitempty"`
	User      User
}
