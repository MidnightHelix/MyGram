package repository

import (
	"context"

	"github.com/MidnightHelix/MyGram/internal/infrastructure"
	"github.com/MidnightHelix/MyGram/internal/model"
	"gorm.io/gorm"
)

type CommentQuery interface {
	GetComments(ctx context.Context, userID uint64) ([]model.Comment, error)
	GetCommentsByID(ctx context.Context, id uint64) (model.Comment, error)

	CreateComment(ctx context.Context, comment model.Comment) (model.Comment, error)
	EditComment(ctx context.Context, comment model.Comment, id uint64) (model.Comment, error)
	DeleteComment(ctx context.Context, id uint64) error
}

type CommentCommand interface {
	CreateComment(ctx context.Context, comment model.Comment) (model.Comment, error)
}

type commentQueryImpl struct {
	db infrastructure.GormPostgres
}

func NewCommentQuery(db infrastructure.GormPostgres) CommentQuery {
	return &commentQueryImpl{db: db}
}

func (u *commentQueryImpl) GetComments(ctx context.Context, userID uint64) ([]model.Comment, error) {
	db := u.db.GetConnection()
	photos := []model.Comment{}
	if err := db.
		WithContext(ctx).
		Table("comments").
		Where("user_id = ?", userID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Username", "Email")
		}).
		Preload("Photo", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Title", "Caption", "Url", "UserID")
		}).
		//Preload("User").
		Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}

func (u *commentQueryImpl) GetCommentsByID(ctx context.Context, id uint64) (model.Comment, error) {
	db := u.db.GetConnection()
	comment := model.Comment{}
	if err := db.
		WithContext(ctx).
		Table("comments").
		Where("id = ?", id).
		Find(&comment).Error; err != nil {
		return model.Comment{}, err
	}
	return comment, nil
}

func (u *commentQueryImpl) CreateComment(ctx context.Context, comment model.Comment) (model.Comment, error) {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("comments").
		Save(&comment).Error; err != nil {
		return model.Comment{}, err
	}
	return comment, nil
}

func (u *commentQueryImpl) EditComment(ctx context.Context, comment model.Comment, id uint64) (model.Comment, error) {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("comments").
		Where("id = ?", id).
		Select("message").
		Updates(model.Comment{Message: comment.Message}).Error; err != nil {
		return model.Comment{}, err
	}
	return comment, nil
}

func (u *commentQueryImpl) DeleteComment(ctx context.Context, id uint64) error {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("comments").
		Where("id = ?", id).Delete(&model.Comment{}).Error; err != nil {
		return err
	}
	return nil
}
