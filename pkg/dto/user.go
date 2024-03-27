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
type UserSignUp struct {
	Username string `json:"username" binding:"required" validate:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=6"`
	Age      uint8  `json:"age" binding:"required" validate:"required,min=9"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"  validate:"required,email"`
	Password string `json:"password" binding:"required"`
}
