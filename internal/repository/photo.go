package repository

import (
	"context"

	"github.com/MidnightHelix/MyGram/internal/infrastructure"
	"github.com/MidnightHelix/MyGram/internal/model"
)

type PhotoQuery interface {
	GetPhotos(ctx context.Context, userID uint64) ([]model.Photo, error)
	GetPhotosByID(ctx context.Context, id uint64) (model.Photo, error)

	CreatePhoto(ctx context.Context, photo model.Photo) (model.Photo, error)
	EditPhoto(ctx context.Context, photo model.Photo, id uint64) (model.Photo, error)
	DeletePhoto(ctx context.Context, id uint64) error
}

type PhotoCommand interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
}

type photoQueryImpl struct {
	db infrastructure.GormPostgres
}

func NewPhotoQuery(db infrastructure.GormPostgres) PhotoQuery {
	return &photoQueryImpl{db: db}
}

func (u *photoQueryImpl) GetPhotos(ctx context.Context, userID uint64) ([]model.Photo, error) {
	db := u.db.GetConnection()
	photos := []model.Photo{}
	if err := db.
		WithContext(ctx).
		Table("photos").
		Where("user_id = ?", userID).
		// Preload("User", func(db *gorm.DB) *gorm.DB {
		// 	return db.Select("ID", "Username", "Email")
		// }).
		Preload("User").
		Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}

func (u *photoQueryImpl) GetPhotosByID(ctx context.Context, id uint64) (model.Photo, error) {
	db := u.db.GetConnection()
	photo := model.Photo{}
	if err := db.
		WithContext(ctx).
		Table("photos").
		Where("id = ?", id).
		Find(&photo).Error; err != nil {
		return model.Photo{}, err
	}
	return photo, nil
}

func (u *photoQueryImpl) CreatePhoto(ctx context.Context, photo model.Photo) (model.Photo, error) {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("photos").
		Save(&photo).Error; err != nil {
		return model.Photo{}, err
	}
	return photo, nil
}

func (u *photoQueryImpl) EditPhoto(ctx context.Context, photo model.Photo, id uint64) (model.Photo, error) {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("photos").
		Where("id = ?", id).
		Select("title", "caption", "url").
		Updates(model.Photo{Title: photo.Title, Caption: photo.Caption, Url: photo.Url}).Error; err != nil {
		return model.Photo{}, err
	}
	return photo, nil
}

func (u *photoQueryImpl) DeletePhoto(ctx context.Context, id uint64) error {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("photos").
		Where("id = ?", id).Delete(&model.Photo{}).Error; err != nil {
		return err
	}
	return nil
}
