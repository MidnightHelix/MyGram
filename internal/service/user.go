package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MidnightHelix/MyGram/internal/model"
	"github.com/MidnightHelix/MyGram/internal/repository"
	"github.com/MidnightHelix/MyGram/pkg/helper"
)

type UserService interface {
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUsersById(ctx context.Context, id uint64) (model.User, error)
	SignUp(ctx context.Context, userSignUp model.UserSignUp) (model.User, error)
	Login(ctx context.Context, userLogin model.UserLogin) (model.User, error)
	EditUser(ctx context.Context, editUser model.User, id uint64) (model.User, error)
	DeleteUser(ctx context.Context, id uint64) error
	// misc
	GenerateUserAccessToken(ctx context.Context, user model.User) (token string, err error)
}

type userServiceImpl struct {
	repo repository.UserQuery
}

func NewUserService(repo repository.UserQuery) UserService {
	return &userServiceImpl{repo: repo}
}

func (u *userServiceImpl) GetUsers(ctx context.Context) ([]model.User, error) {
	users, err := u.repo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, err
}

func (u *userServiceImpl) GetUsersById(ctx context.Context, id uint64) (model.User, error) {
	user, err := u.repo.GetUsersByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	return user, err
}

func (u *userServiceImpl) SignUp(ctx context.Context, userSignUp model.UserSignUp) (model.User, error) {
	// assumption: semua user adalah user baru
	user := model.User{
		Username: userSignUp.Username,
		Email:    userSignUp.Email,
	}

	// encryption password
	// hashing
	pass, err := helper.GenerateHash(userSignUp.Password)
	if err != nil {
		return model.User{}, err
	}
	user.Password = pass

	// store to db
	res, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return model.User{}, err
	}
	return res, err
}

func (u *userServiceImpl) Login(ctx context.Context, userLogin model.UserLogin) (model.User, error) {
	user, err := u.repo.FindByEmail(ctx, userLogin.Email)
	if err != nil {
		return model.User{}, err
	}

	err = helper.CompareHashAndPassword(user.Password, userLogin.Password)
	if err != nil {
		return model.User{}, err
	}

	return user, err
}

func (u *userServiceImpl) GenerateUserAccessToken(ctx context.Context, user model.User) (token string, err error) {
	// generate claim
	now := time.Now()

	claim := model.StandardClaim{
		Jti: fmt.Sprintf("%v", time.Now().UnixNano()),
		Iss: "go-middleware",
		Aud: "golang-006",
		Sub: "access-token",
		Exp: uint64(now.Add(time.Hour).Unix()),
		Iat: uint64(now.Unix()),
		Nbf: uint64(now.Unix()),
	}

	userClaim := model.AccessClaim{
		StandardClaim: claim,
		UserID:        user.ID,
		Username:      user.Username,
		Dob:           user.DoB,
	}

	// userClaim := model.UserClaim{
	// 	UserID:   user.ID,
	// 	Username: user.Username,
	// 	Dob:      user.DoB,
	// }

	token, err = helper.GenerateToken(userClaim)
	return
}

func (u *userServiceImpl) EditUser(ctx context.Context, editUser model.User, id uint64) (model.User, error) {
	res, err := u.repo.EditUser(ctx, editUser, id)
	if err != nil {
		return model.User{}, err
	}
	return res, err
}

func (u *userServiceImpl) DeleteUser(ctx context.Context, id uint64) error {
	err := u.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
