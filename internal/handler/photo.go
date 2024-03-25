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

type PhotoHandler interface {
	GetPhotos(ctx *gin.Context)
	EditPhoto(ctx *gin.Context)
	DeletePhoto(ctx *gin.Context)

	PostPhoto(ctx *gin.Context)
}

type photoHandlerImpl struct {
	svc service.PhotoService
}

func NewPhotoHandler(svc service.PhotoService) PhotoHandler {
	return &photoHandlerImpl{
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
func (u *photoHandlerImpl) GetPhotos(ctx *gin.Context) {
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

	photos, err := u.svc.GetPhotos(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	var data []dto.Photo
	for _, item := range photos {
		photo := dto.Photo{
			ID:        item.ID,
			Title:     item.Title,
			Caption:   item.Caption,
			Url:       item.Url,
			UserID:    item.UserID,
			CreatedAt: &item.CreatedAt,
			UpdatedAt: &item.UpdatedAt,
			User: &dto.UserDefault{
				Email:    item.User.Email,
				Username: item.User.Username,
			},
		}
		data = append(data, photo)
	}

	ctx.JSON(http.StatusOK, pkg.SuccessResponse{Data: data})
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
// func (u *photoHandlerImpl) GetUsersById(ctx *gin.Context) {
// 	// get id user
// 	id, err := strconv.Atoi(ctx.Param("id"))
// 	if id == 0 || err != nil {
// 		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
// 		return
// 	}
// 	user, err := u.svc.GetUsersById(ctx, uint64(id))
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, user)
// }

func (u *photoHandlerImpl) PostPhoto(ctx *gin.Context) {
	claims, ok := ctx.Get("claims")
	if !ok {
		fmt.Println("Failed to Retrieve Claims")
		return
	}

	userID := claims.(jwt.MapClaims)["user_id"].(float64)
	id := int(userID)
	if id == 0 {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid Id"})
		return
	}

	// binding sign-up body
	photo := model.Photo{}
	if err := ctx.Bind(&photo); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	// if err := userSignUp.Validate(); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
	// 	return
	// }

	photo, err := u.svc.PostPhoto(ctx, photo, uint64(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	data := dto.Photo{
		ID:        photo.ID,
		Title:     photo.Title,
		Caption:   photo.Caption,
		Url:       photo.Url,
		UserID:    photo.UserID,
		CreatedAt: &photo.CreatedAt,
	}

	ctx.JSON(http.StatusCreated, pkg.SuccessResponse{Data: data})
}

func (u *photoHandlerImpl) EditPhoto(ctx *gin.Context) {
	claims, ok := ctx.Get("claims")
	if !ok {
		fmt.Println("Failed to Retrieve Claims")
		return
	}

	userID := claims.(jwt.MapClaims)["user_id"].(float64)
	userId := int(userID)
	if userId == 0 {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}
	// photo, err := u.svc.GetPhotosById(ctx, uint64(id))
	// fmt.Println("photo ID : ", photo.ID)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
	// 	return
	// }
	// if photo.ID == 0 || photo.UserID != uint64(userId) {
	// 	ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Photo Not Found"})
	// 	return
	// }

	req := model.Photo{}
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	// if err := userLogin.Validate(); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
	// 	return
	// }

	photo, err := u.svc.EditPhoto(ctx, req, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	//ctx.JSON(http.StatusOK, photo)
	data := dto.Photo{
		ID:        photo.ID,
		Title:     photo.Title,
		Caption:   photo.Caption,
		Url:       photo.Url,
		UserID:    photo.UserID,
		UpdatedAt: &photo.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, pkg.SuccessResponse{Data: data})
}

func (u *photoHandlerImpl) DeletePhoto(ctx *gin.Context) {
	claims, ok := ctx.Get("claims")
	if !ok {
		fmt.Println("Failed to Retrieve Claims")
		return
	}

	userID := claims.(jwt.MapClaims)["user_id"].(float64)
	userId := int(userID)
	if userId == 0 {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}
	// photo, err := u.svc.GetPhotosById(ctx, uint64(id))
	// fmt.Println("photo ID : ", photo.ID)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
	// 	return
	// }
	// if photo.ID == 0 || photo.UserID != uint64(userId) {
	// 	ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Photo Not Found"})
	// 	return
	// }

	err = u.svc.DeletePhoto(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pkg.SuccessResponse{Message: "Your photo has been successfully deleted"})
}
