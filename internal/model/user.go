package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint64         `json:"id,omitempty" gorm:"primaryKey"`
	Username     string         `json:"username,omitempty" gorm:"not null;unique"`
	Email        string         `json:"email,omitempty" gorm:"not null;unique"`
	Password     string         `json:"password,omitempty"`
	DoB          time.Time      `json:"dob,omitempty"`
	Age          uint8          `json:"age,omitempty"`
	CreatedAt    time.Time      `json:"created_at,omitempty"`
	UpdatedAt    time.Time      `json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty"`
	Photos       []Photo        `json:"photos,omitempty"`
	SocialMedias []SocialMedia  `json:"social_medias,omitempty"`
	Comments     []Comment      `json:"comments,omitempty"`
}

type UserMediaSocial struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"user_id"`
	Title     string `json:"title"`
	Url       string `json:"url"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type UserSignUp struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u UserSignUp) Validate() error {
	// check username
	if u.Username == "" {
		return errors.New("invalid username")
	}
	if len(u.Password) < 6 {
		return errors.New("invalid password")
	}
	return nil
}
