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

type SocialMediaHandler interface {
	GetSocialMedias(ctx *gin.Context)
	EditSocialMedia(ctx *gin.Context)
	DeleteSocialMedia(ctx *gin.Context)

	CreateSocialMedia(ctx *gin.Context)
}

type socialMediaHandlerImpl struct {
	svc       service.SocialMediaService
	validator *validator.CustomValidator
}

func NewSocialMediaHandler(svc service.SocialMediaService, validator *validator.CustomValidator) SocialMediaHandler {
	return &socialMediaHandlerImpl{
		svc:       svc,
		validator: validator,
	}
}

// ShowSocialMedia godoc
//
//	@Summary		Show scial media list
//	@Description	Get social medias of user
//	@Tags			socialmedias
//	@Accept			json
//	@Produce		json
//
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
//
//	@Success		200	{object}	[]dto.SocialMedia
//	@Failure		400	{object}	pkg.ErrorResponse
//	@Failure		404	{object}	pkg.ErrorResponse
//	@Failure		500	{object}	pkg.ErrorResponse
//	@Router			/socialmedias [get]
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

//	 CreateSocialMedia godoc
//
//		@Summary		create social media
//		@Description	create social media with input payload
//		@Tags			socialmedias
//		@Accept			json
//		@Produce		json
//
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
//
//	@Param socialMedia body SocialMedia true "Create Social Media"
//	@Success		200	{object}	dto.SocialMedia
//	@Failure		400	{object}	pkg.ErrorResponse
//	@Failure		404	{object}	pkg.ErrorResponse
//	@Failure		500	{object}	pkg.ErrorResponse
//	@Router			/socialmedias [post]
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

	if err := u.validator.ValidateStruct(socialMedia); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

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

//	 UpdateSocialMedia godoc
//
//		@Summary		Update social media
//		@Description	Update social media with input payload
//		@Tags			socialmedias
//		@Accept			json
//		@Produce		json
//
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
//
//	@Param socialMedia body SocialMedia true "Update Social Media"
//	@Param        id   path      int  true  "SocialMedia ID"
//	@Success		200	{object}	dto.SocialMedia
//	@Failure		400	{object}	pkg.ErrorResponse
//	@Failure		404	{object}	pkg.ErrorResponse
//	@Failure		500	{object}	pkg.ErrorResponse
//	@Router			/socialmedias/{id} [put]
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

	if err := u.validator.ValidateStruct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

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

// DeleteSocialMedia godoc
//
// @Summary		Delete social media
// @Description	Delete social media with id param
// @Tags			socialmedias
// @Accept			json
// @Produce		json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "SocialMedia ID"
// @Success		200	{object}	pkg.SuccessResponse
// @Failure		400	{object}	pkg.ErrorResponse
// @Failure		404	{object}	pkg.ErrorResponse
// @Failure		500	{object}	pkg.ErrorResponse
// @Router			/socialmedias/{id} [delete]
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
