package repository

import (
	"context"

	"github.com/MidnightHelix/MyGram/internal/infrastructure"
	"github.com/MidnightHelix/MyGram/internal/model"
)

type UserQuery interface {
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUsersByID(ctx context.Context, id uint64) (model.User, error)
	FindByEmail(ctx context.Context, email string) (model.User, error)

	CreateUser(ctx context.Context, user model.User) (model.User, error)
	EditUser(ctx context.Context, user model.User, id uint64) (model.User, error)
	DeleteUser(ctx context.Context, id uint64) error
}

type UserCommand interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
}

type userQueryImpl struct {
	db infrastructure.GormPostgres
}

func NewUserQuery(db infrastructure.GormPostgres) UserQuery {
	return &userQueryImpl{db: db}
}

func (u *userQueryImpl) GetUsers(ctx context.Context) ([]model.User, error) {
	db := u.db.GetConnection()
	users := []model.User{}
	if err := db.
		WithContext(ctx).
		Table("users").
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userQueryImpl) GetUsersByID(ctx context.Context, id uint64) (model.User, error) {
	db := u.db.GetConnection()
	users := model.User{}
	if err := db.
		WithContext(ctx).
		Table("users").
		Where("id = ?", id).
		Find(&users).Error; err != nil {
		return model.User{}, err
	}
	return users, nil
}

func (u *userQueryImpl) FindByEmail(ctx context.Context, email string) (model.User, error) {
	db := u.db.GetConnection()
	user := model.User{}
	if err := db.
		WithContext(ctx).
		Table("users").
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *userQueryImpl) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("users").
		Save(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *userQueryImpl) EditUser(ctx context.Context, user model.User, id uint64) (model.User, error) {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("users").
		// Clauses(clause.Returning{}).
		Where("id = ?", id).
		Select("username", "email").
		Updates(model.User{Username: user.Username, Email: user.Email}).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *userQueryImpl) DeleteUser(ctx context.Context, id uint64) error {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("users").
		Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
		return err
	}
	return nil
}
