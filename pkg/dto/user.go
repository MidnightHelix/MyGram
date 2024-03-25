package dto

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint64          `json:"id,omitempty"`
	Email        string          `json:"email,omitempty"`
	Username     string          `json:"username,omitempty"`
	DoB          *time.Time      `json:"dob,omitempty"`
	Age          *uint8          `json:"age,omitempty"`
	CreatedAt    *time.Time      `json:"created_at,omitempty"`
	UpdatedAt    *time.Time      `json:"updated_at,omitempty"`
	DeletedAt    *gorm.DeletedAt `json:"deleted_at,omitempty"`
	Photos       []Photo         `json:"photos,omitempty"`
	SocialMedias []SocialMedia   `json:"social_medias,omitempty"`
	Comments     []Comment       `json:"comments,omitempty"`
}

type UserDefault struct {
	ID       *uint64 `json:"id,omitempty"`
	Email    string  `json:"email"`
	Username string  `json:"username"`
}
