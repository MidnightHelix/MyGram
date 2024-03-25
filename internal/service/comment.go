package service

import (
	"context"

	"github.com/MidnightHelix/MyGram/internal/model"
	"github.com/MidnightHelix/MyGram/internal/repository"
)

type CommentService interface {
	GetComments(ctx context.Context, userID uint64) ([]model.Comment, error)
	GetCommentsById(ctx context.Context, id uint64) (model.Comment, error)
	PostComment(ctx context.Context, comment model.Comment, userID uint64) (model.Comment, error)

	EditComment(ctx context.Context, comment model.Comment, id uint64) (model.Comment, error)
	DeleteComment(ctx context.Context, id uint64) error
}

type commentServiceImpl struct {
	repo repository.CommentQuery
}

func NewCommentService(repo repository.CommentQuery) CommentService {
	return &commentServiceImpl{repo: repo}
}

func (u *commentServiceImpl) GetComments(ctx context.Context, userID uint64) ([]model.Comment, error) {
	comments, err := u.repo.GetComments(ctx, userID)
	if err != nil {
		return nil, err
	}
	return comments, err
}

func (u *commentServiceImpl) GetCommentsById(ctx context.Context, id uint64) (model.Comment, error) {
	comment, err := u.repo.GetCommentsByID(ctx, id)
	if err != nil {
		return model.Comment{}, err
	}
	return comment, err
}

func (u *commentServiceImpl) PostComment(ctx context.Context, comment model.Comment, userID uint64) (model.Comment, error) {

	user := model.Comment{
		Message: comment.Message,
		PhotoID: comment.PhotoID,
		UserID:  userID,
	}

	// store to db
	res, err := u.repo.CreateComment(ctx, user)
	if err != nil {
		return model.Comment{}, err
	}
	return res, err
}

func (u *commentServiceImpl) EditComment(ctx context.Context, comment model.Comment, id uint64) (model.Comment, error) {
	res, err := u.repo.EditComment(ctx, comment, id)
	if err != nil {
		return model.Comment{}, err
	}
	return res, err
}

func (u *commentServiceImpl) DeleteComment(ctx context.Context, id uint64) error {
	err := u.repo.DeleteComment(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
