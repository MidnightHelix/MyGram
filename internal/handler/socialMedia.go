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

type SocialMediaHandler interface {
	GetSocialMedias(ctx *gin.Context)
	EditSocialMedia(ctx *gin.Context)
	DeleteSocialMedia(ctx *gin.Context)

	CreateSocialMedia(ctx *gin.Context)
}

type socialMediaHandlerImpl struct {
	svc service.SocialMediaService
}

func NewSocialMediaHandler(svc service.SocialMediaService) SocialMediaHandler {
	return &socialMediaHandlerImpl{
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
func (u *socialMediaHandlerImpl) GetSocialMedias(ctx *gin.Context) {
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

	socialMedias, err := u.svc.GetSocialMedias(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	var data []dto.SocialMedia
	for _, item := range socialMedias {
		socialMedia := dto.SocialMedia{
			ID:        item.ID,
			Name:      item.Name,
			Url:       item.Url,
			UserID:    item.UserID,
			CreatedAt: &item.CreatedAt,
			UpdatedAt: &item.UpdatedAt,
			User: &dto.UserDefault{
				ID:       &item.User.ID,
				Email:    item.User.Email,
				Username: item.User.Username,
			},
		}
		data = append(data, socialMedia)
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

func (u *socialMediaHandlerImpl) CreateSocialMedia(ctx *gin.Context) {
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

	socialMedia := model.SocialMedia{}
	if err := ctx.Bind(&socialMedia); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	// if err := userSignUp.Validate(); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
	// 	return
	// }

	socialMedia, err := u.svc.CreateSocialMedia(ctx, socialMedia, uint64(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	data := dto.SocialMedia{
		ID:        socialMedia.ID,
		Name:      socialMedia.Name,
		Url:       socialMedia.Url,
		UserID:    socialMedia.UserID,
		CreatedAt: &socialMedia.CreatedAt,
	}
	ctx.JSON(http.StatusCreated, pkg.SuccessResponse{Data: data})
}

func (u *socialMediaHandlerImpl) EditSocialMedia(ctx *gin.Context) {
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
	// socialMedia, err := u.svc.GetSocialMediaById(ctx, uint64(id))
	// fmt.Println("socialMedia ID : ", socialMedia.ID)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
	// 	return
	// }
	// if socialMedia.ID == 0 || socialMedia.UserID != uint64(userId) {
	// 	ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Social Media Not Found"})
	// 	return
	// }

	req := model.SocialMedia{}
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	// if err := userLogin.Validate(); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
	// 	return
	// }

	socialMedia, err := u.svc.EditSocialMedia(ctx, req, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	data := dto.SocialMedia{
		ID:        socialMedia.ID,
		Name:      socialMedia.Name,
		Url:       socialMedia.Url,
		UserID:    socialMedia.UserID,
		UpdatedAt: &socialMedia.UpdatedAt,
	}
	ctx.JSON(http.StatusOK, pkg.SuccessResponse{Data: data})
}

func (u *socialMediaHandlerImpl) DeleteSocialMedia(ctx *gin.Context) {
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
	// socialMedia, err := u.svc.GetSocialMediaById(ctx, uint64(id))
	// fmt.Println("socialMedia ID : ", socialMedia.ID)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
	// 	return
	// }
	// if socialMedia.ID == 0 || socialMedia.UserID != uint64(userId) {
	// 	ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Social Media Not Found"})
	// 	return
	// }

	err = u.svc.DeleteSocialMedia(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pkg.SuccessResponse{Message: "Your social media has been successfully deleted"})
}
