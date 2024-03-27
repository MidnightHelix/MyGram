package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MidnightHelix/MyGram/internal/model"
	"github.com/MidnightHelix/MyGram/internal/service"
	"github.com/MidnightHelix/MyGram/pkg"
	"github.com/MidnightHelix/MyGram/pkg/dto"
	"github.com/MidnightHelix/MyGram/pkg/validator"
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
	svc       service.PhotoService
	validator *validator.CustomValidator
}

func NewPhotoHandler(svc service.PhotoService, validator *validator.CustomValidator) PhotoHandler {
	return &photoHandlerImpl{
		svc:       svc,
		validator: validator,
	}
}

// ShowPhotos godoc
//
//	@Summary		Show photo list
//	@Description	Get photos of user
//	@Tags			photos
//	@Accept			json
//	@Produce		json
//
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
//
//	@Success		200	{object}	[]dto.Photo
//	@Failure		400	{object}	pkg.ErrorResponse
//	@Failure		404	{object}	pkg.ErrorResponse
//	@Failure		500	{object}	pkg.ErrorResponse
//	@Router			/photos [get]
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

//	 PostPhoto godoc
//
//		@Summary		Post a photo
//		@Description	Create photo with input payload
//		@Tags			photos
//		@Accept			json
//		@Produce		json
//
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
//
//	@Param photo body Photo true "Create Photo"
//	@Success		201	{object}	dto.Photo
//	@Failure		400	{object}	pkg.ErrorResponse
//	@Failure		404	{object}	pkg.ErrorResponse
//	@Failure		500	{object}	pkg.ErrorResponse
//	@Router			/photos [post]
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

	if err := u.validator.ValidateStruct(photo); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

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

//	 UpdatePhoto godoc
//
//		@Summary		Update a photo
//		@Description	Update photo with input payload
//		@Tags			photos
//		@Accept			json
//		@Produce		json
//
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
//
//	@Param photo body Photo true "Update Photo"
//	@Param        id   path      int  true  "Photo ID"
//	@Success		200	{object}	dto.Photo
//	@Failure		400	{object}	pkg.ErrorResponse
//	@Failure		404	{object}	pkg.ErrorResponse
//	@Failure		500	{object}	pkg.ErrorResponse
//	@Router			/photos/{id} [put]
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
	if err := u.validator.ValidateStruct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	photo, err := u.svc.EditPhoto(ctx, req, uint64(id))
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
		UpdatedAt: &photo.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, pkg.SuccessResponse{Data: data})
}

// DeletePhoto godoc
//
// @Summary		Delete a photo
// @Description	Delete photo with id param
// @Tags			photos
// @Accept			json
// @Produce		json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "Photo ID"
// @Success		200	{object}	pkg.SuccessResponse
// @Failure		400	{object}	pkg.ErrorResponse
// @Failure		404	{object}	pkg.ErrorResponse
// @Failure		500	{object}	pkg.ErrorResponse
// @Router			/photos/{id} [delete]
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
