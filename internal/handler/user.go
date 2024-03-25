package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MidnightHelix/MyGram/internal/model"
	"github.com/MidnightHelix/MyGram/internal/service"
	"github.com/MidnightHelix/MyGram/pkg"
	"github.com/MidnightHelix/MyGram/pkg/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler interface {
	GetUsers(ctx *gin.Context)
	GetUsersById(ctx *gin.Context)
	UserSignUp(ctx *gin.Context)
	UserLogin(ctx *gin.Context)
	EditUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type userHandlerImpl struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) UserHandler {
	return &userHandlerImpl{
		svc: svc,
	}
}

// ShowUsers godoc
//
//	@Summary		Show users list
//	@Description	will fetch 3rd party server to get users data
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]model.User
//	@Failure		400	{object}	pkg.ErrorResponse
//	@Failure		404	{object}	pkg.ErrorResponse
//	@Failure		500	{object}	pkg.ErrorResponse
//	@Router			/users [get]
func (u *userHandlerImpl) GetUsers(ctx *gin.Context) {
	users, err := u.svc.GetUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// ShowUsersById godoc
//
//	@Summary		Show users detail
//	@Description	will fetch 3rd party server to get users data to get detail user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	model.User
//	@Failure		400	{object}	pkg.ErrorResponse
//	@Failure		404	{object}	pkg.ErrorResponse
//	@Failure		500	{object}	pkg.ErrorResponse
//	@Router			/users/{id} [get]
func (u *userHandlerImpl) GetUsersById(ctx *gin.Context) {
	// get id user
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}
	user, err := u.svc.GetUsersById(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (u *userHandlerImpl) UserSignUp(ctx *gin.Context) {
	// binding sign-up body
	userSignUp := model.UserSignUp{}
	if err := ctx.Bind(&userSignUp); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	if err := userSignUp.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := u.svc.SignUp(ctx, userSignUp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	token, err := u.svc.GenerateUserAccessToken(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
	}

	data := map[string]any{
		"token": token,
	}
	ctx.JSON(http.StatusCreated, pkg.SuccessResponse{Data: data})
}

func (u *userHandlerImpl) UserLogin(ctx *gin.Context) {
	userLogin := model.UserLogin{}
	if err := ctx.Bind(&userLogin); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	// if err := userLogin.Validate(); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
	// 	return
	// }

	user, err := u.svc.Login(ctx, userLogin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	token, err := u.svc.GenerateUserAccessToken(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
	}

	data := map[string]any{
		"token": token,
	}
	ctx.JSON(http.StatusOK, pkg.SuccessResponse{Data: data})
}

func (u *userHandlerImpl) EditUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}
	user, err := u.svc.GetUsersById(ctx, uint64(id))
	fmt.Println("user ID : ", user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	if user.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "User Not Found"})
		return
	}

	req := model.User{}
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	// if err := userLogin.Validate(); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
	// 	return
	// }

	user, err = u.svc.EditUser(ctx, req, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	data := dto.User{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Age:       &user.Age,
		UpdatedAt: &user.UpdatedAt,
	}
	ctx.JSON(http.StatusOK, pkg.SuccessResponse{Data: data})
}

func (u *userHandlerImpl) DeleteUser(ctx *gin.Context) {
	claims, ok := ctx.Get("claims")
	if !ok {
		fmt.Println("Failed to Retrieve Claims")
		return
	}

	userID := claims.(jwt.MapClaims)["user_id"].(float64)
	id := int(userID)
	if id == 0 {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}
	user, err := u.svc.GetUsersById(ctx, uint64(id))
	fmt.Println("user ID : ", user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	if user.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "User Not Found"})
		return
	}

	err = u.svc.DeleteUser(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pkg.SuccessResponse{Message: "Your account has been successfully deleted"})
}
