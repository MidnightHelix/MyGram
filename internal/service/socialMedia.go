package service

import (
	"context"

	"github.com/MidnightHelix/MyGram/internal/model"
	"github.com/MidnightHelix/MyGram/internal/repository"
)

type SocialMediaService interface {
	GetSocialMedias(ctx context.Context, userID uint64) ([]model.SocialMedia, error)
	GetSocialMediaById(ctx context.Context, id uint64) (model.SocialMedia, error)
	CreateSocialMedia(ctx context.Context, socialMedia model.SocialMedia, userID uint64) (model.SocialMedia, error)

	EditSocialMedia(ctx context.Context, socialMedia model.SocialMedia, id uint64) (model.SocialMedia, error)
	DeleteSocialMedia(ctx context.Context, id uint64) error
}

type socialMediaServiceImpl struct {
	repo repository.SocialMediaQuery
}

func NewSocialMediaService(repo repository.SocialMediaQuery) SocialMediaService {
	return &socialMediaServiceImpl{repo: repo}
}

func (u *socialMediaServiceImpl) GetSocialMedias(ctx context.Context, userID uint64) ([]model.SocialMedia, error) {
	socialMedias, err := u.repo.GetSocialMedias(ctx, userID)
	if err != nil {
		return nil, err
	}
	return socialMedias, err
}

func (u *socialMediaServiceImpl) GetSocialMediaById(ctx context.Context, id uint64) (model.SocialMedia, error) {
	socialMedia, err := u.repo.GetSocialMediaByID(ctx, id)
	if err != nil {
		return model.SocialMedia{}, err
	}
	return socialMedia, err
}

func (u *socialMediaServiceImpl) CreateSocialMedia(ctx context.Context, req model.SocialMedia, userID uint64) (model.SocialMedia, error) {

	socialMedia := model.SocialMedia{
		Name:   req.Name,
		Url:    req.Url,
		UserID: userID,
	}

	// store to db
	res, err := u.repo.CreateSocialMedia(ctx, socialMedia)
	if err != nil {
		return model.SocialMedia{}, err
	}
	return res, err
}

func (u *socialMediaServiceImpl) EditSocialMedia(ctx context.Context, socialMedia model.SocialMedia, id uint64) (model.SocialMedia, error) {
	res, err := u.repo.EditSocialMedia(ctx, socialMedia, id)
	if err != nil {
		return model.SocialMedia{}, err
	}
	return res, err
}

func (u *socialMediaServiceImpl) DeleteSocialMedia(ctx context.Context, id uint64) error {
	err := u.repo.DeleteSocialMedia(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
