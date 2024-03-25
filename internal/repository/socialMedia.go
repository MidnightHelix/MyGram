package repository

import (
	"context"

	"github.com/MidnightHelix/MyGram/internal/infrastructure"
	"github.com/MidnightHelix/MyGram/internal/model"
	"gorm.io/gorm"
)

type SocialMediaQuery interface {
	GetSocialMedias(ctx context.Context, userID uint64) ([]model.SocialMedia, error)
	GetSocialMediaByID(ctx context.Context, id uint64) (model.SocialMedia, error)

	CreateSocialMedia(ctx context.Context, socialMedia model.SocialMedia) (model.SocialMedia, error)
	EditSocialMedia(ctx context.Context, socialMedia model.SocialMedia, id uint64) (model.SocialMedia, error)
	DeleteSocialMedia(ctx context.Context, id uint64) error
}

type SocialMediaCommand interface {
	CreateComment(ctx context.Context, comment model.Comment) (model.Comment, error)
}

type socialMediaQueryImpl struct {
	db infrastructure.GormPostgres
}

func NewSocialMediaQuery(db infrastructure.GormPostgres) SocialMediaQuery {
	return &socialMediaQueryImpl{db: db}
}

func (u *socialMediaQueryImpl) GetSocialMedias(ctx context.Context, userID uint64) ([]model.SocialMedia, error) {
	db := u.db.GetConnection()
	socialMedia := []model.SocialMedia{}
	if err := db.
		WithContext(ctx).
		Table("social_media").
		Where("user_id = ?", userID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Username", "Email")
		}).
		//Preload("User").
		Find(&socialMedia).Error; err != nil {
		return nil, err
	}
	return socialMedia, nil
}

func (u *socialMediaQueryImpl) GetSocialMediaByID(ctx context.Context, id uint64) (model.SocialMedia, error) {
	db := u.db.GetConnection()
	socialMedia := model.SocialMedia{}
	if err := db.
		WithContext(ctx).
		Table("social_media").
		Where("id = ?", id).
		Find(&socialMedia).Error; err != nil {
		return model.SocialMedia{}, err
	}
	return socialMedia, nil
}

func (u *socialMediaQueryImpl) CreateSocialMedia(ctx context.Context, socialMedia model.SocialMedia) (model.SocialMedia, error) {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("social_media").
		Save(&socialMedia).Error; err != nil {
		return model.SocialMedia{}, err
	}
	return socialMedia, nil
}

func (u *socialMediaQueryImpl) EditSocialMedia(ctx context.Context, socialMedia model.SocialMedia, id uint64) (model.SocialMedia, error) {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("social_media").
		Where("id = ?", id).
		Select("name", "url").
		Updates(model.SocialMedia{Name: socialMedia.Name, Url: socialMedia.Url}).Error; err != nil {
		return model.SocialMedia{}, err
	}
	return socialMedia, nil
}

func (u *socialMediaQueryImpl) DeleteSocialMedia(ctx context.Context, id uint64) error {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("social_media").
		Where("id = ?", id).Delete(&model.SocialMedia{}).Error; err != nil {
		return err
	}
	return nil
}
