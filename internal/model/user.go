package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint64         `json:"id,omitempty" gorm:"primaryKey"`
	Username     string         `json:"username,omitempty" gorm:"not null;unique;uniqueIndex" binding:"required" validate:"required,min=3,max=50"`
	Email        string         `json:"email,omitempty" gorm:"not null;unique;uniqueIndex" binding:"required" validate:"required,email"`
	Password     string         `json:"password,omitempty" gorm:"not null"`
	DoB          time.Time      `json:"dob,omitempty" gorm:"not null"`
	Age          uint8          `json:"age,omitempty" gorm:"not null" binding:"required" validate:"required,min=9"`
	CreatedAt    time.Time      `json:"created_at,omitempty"`
	UpdatedAt    time.Time      `json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty"`
	Photos       []Photo        `json:"photos,omitempty"`
	SocialMedias []SocialMedia  `json:"social_medias,omitempty"`
	Comments     []Comment      `json:"comments,omitempty"`
}
