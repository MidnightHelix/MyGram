package service

import (
	"context"

	"github.com/MidnightHelix/MyGram/internal/model"
	"github.com/MidnightHelix/MyGram/internal/repository"
)

type PhotoService interface {
	GetPhotos(ctx context.Context, userID uint64) ([]model.Photo, error)
	GetPhotosById(ctx context.Context, id uint64) (model.Photo, error)
	PostPhoto(ctx context.Context, photo model.Photo, userID uint64) (model.Photo, error)

	EditPhoto(ctx context.Context, photo model.Photo, id uint64) (model.Photo, error)
	DeletePhoto(ctx context.Context, id uint64) error
}

type photoServiceImpl struct {
	repo repository.PhotoQuery
}

func NewPhotoService(repo repository.PhotoQuery) PhotoService {
	return &photoServiceImpl{repo: repo}
}

func (u *photoServiceImpl) GetPhotos(ctx context.Context, userID uint64) ([]model.Photo, error) {
	photos, err := u.repo.GetPhotos(ctx, userID)
	if err != nil {
		return nil, err
	}
	return photos, err
}

func (u *photoServiceImpl) GetPhotosById(ctx context.Context, id uint64) (model.Photo, error) {
	photo, err := u.repo.GetPhotosByID(ctx, id)
	if err != nil {
		return model.Photo{}, err
	}
	return photo, err
}

func (u *photoServiceImpl) PostPhoto(ctx context.Context, photo model.Photo, userID uint64) (model.Photo, error) {

	user := model.Photo{
		Title:   photo.Title,
		Caption: photo.Caption,
		Url:     photo.Url,
		UserID:  userID,
	}

	// store to db
	res, err := u.repo.CreatePhoto(ctx, user)
	if err != nil {
		return model.Photo{}, err
	}
	return res, err
}

func (u *photoServiceImpl) EditPhoto(ctx context.Context, photo model.Photo, id uint64) (model.Photo, error) {
	res, err := u.repo.EditPhoto(ctx, photo, id)
	if err != nil {
		return model.Photo{}, err
	}
	return res, err
}

func (u *photoServiceImpl) DeletePhoto(ctx context.Context, id uint64) error {
	err := u.repo.DeletePhoto(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
