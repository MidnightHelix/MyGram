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

type CommentHandler interface {
	GetComments(ctx *gin.Context)
	EditComment(ctx *gin.Context)
	DeleteComment(ctx *gin.Context)

	PostComment(ctx *gin.Context)
}

type commentHandlerImpl struct {
	svc       service.CommentService
	validator *validator.CustomValidator
}

func NewCommentHandler(svc service.CommentService, validator *validator.CustomValidator) CommentHandler {
	return &commentHandlerImpl{
		svc:       svc,
		validator: validator,
	}
}

// ShowComments godoc
//
//	@Summary		Show comment list
//	@Description	Get comments of user
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]dto.Comment
//	@Failure		400	{object}	pkg.ErrorResponse
//	@Failure		404	{object}	pkg.ErrorResponse
//	@Failure		500	{object}	pkg.ErrorResponse
//	@Router			/comments [get]
func (u *commentHandlerImpl) GetComments(ctx *gin.Context) {
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

	comments, err := u.svc.GetComments(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	var data []dto.Comment
	for _, item := range comments {
		comment := dto.Comment{
			ID:        item.ID,
			Message:   item.Message,
			PhotoID:   item.PhotoID,
			UserID:    item.UserID,
			CreatedAt: &item.CreatedAt,
			UpdatedAt: &item.UpdatedAt,
			User: &dto.UserDefault{
				ID:       &item.UserID,
				Email:    item.User.Email,
				Username: item.User.Username,
			},
			Photo: &dto.Photo{
				ID:      item.Photo.ID,
				Title:   item.Photo.Title,
				Caption: item.Photo.Caption,
				Url:     item.Photo.Url,
				UserID:  item.Photo.UserID,
			},
		}
		data = append(data, comment)
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

//	 PostComment godoc
//
//		@Summary		Post a comment
//		@Description	Create comment with input payload
//		@Tags			comments
//		@Accept			json
//		@Produce		json
//
// @Param comment body Comment true "Create Comment"
// @Success		201	{object}	dto.Comment
// @Failure		400	{object}	pkg.ErrorResponse
// @Failure		404	{object}	pkg.ErrorResponse
// @Failure		500	{object}	pkg.ErrorResponse
// @Router			/comments [post]
func (u *commentHandlerImpl) PostComment(ctx *gin.Context) {
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

	comment := model.Comment{}
	if err := ctx.Bind(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	// if err := userSignUp.Validate(); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
	// 	return
	// }

	if err := u.validator.ValidateStruct(comment); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	comment, err := u.svc.PostComment(ctx, comment, uint64(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	data := dto.Comment{
		ID:        comment.ID,
		Message:   comment.Message,
		PhotoID:   comment.PhotoID,
		UserID:    comment.UserID,
		CreatedAt: &comment.CreatedAt,
	}
	ctx.JSON(http.StatusCreated, pkg.SuccessResponse{Data: data})
}

//	 UpdateComment godoc
//
//		@Summary		Update a comment
//		@Description	Update comment with input payload
//		@Tags			comments
//		@Accept			json
//		@Produce		json
//
// @Param comment body Comment true "Update Comment"
// @Param        id   path      int  true  "Comment ID"
// @Success		200	{object}	dto.Comment
// @Failure		400	{object}	pkg.ErrorResponse
// @Failure		404	{object}	pkg.ErrorResponse
// @Failure		500	{object}	pkg.ErrorResponse
// @Router			/comments/{id} [put]
func (u *commentHandlerImpl) EditComment(ctx *gin.Context) {
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
	// comment, err := u.svc.GetCommentsById(ctx, uint64(id))
	// fmt.Println("comment ID : ", comment.ID)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
	// 	return
	// }
	// if comment.ID == 0 || comment.UserID != uint64(userId) {
	// 	ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Comment Not Found"})
	// 	return
	// }

	req := model.Comment{}
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

	comment, err := u.svc.EditComment(ctx, req, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	data := dto.Comment{
		ID:        comment.ID,
		Message:   comment.Message,
		PhotoID:   comment.PhotoID,
		UserID:    comment.UserID,
		UpdatedAt: &comment.UpdatedAt,
	}
	ctx.JSON(http.StatusOK, pkg.SuccessResponse{Data: data})
}

// DeleteComment godoc
//
// @Summary		Delete a comment
// @Description	Delete comment with id param
// @Tags			comments
// @Accept			json
// @Produce		json
// @Param        id   path      int  true  "Comment ID"
// @Success		200	{object}	pkg.SuccessResponse
// @Failure		400	{object}	pkg.ErrorResponse
// @Failure		404	{object}	pkg.ErrorResponse
// @Failure		500	{object}	pkg.ErrorResponse
// @Router			/comments/{id} [delete]
func (u *commentHandlerImpl) DeleteComment(ctx *gin.Context) {
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
	// comment, err := u.svc.GetCommentsById(ctx, uint64(id))
	// fmt.Println("comment ID : ", comment.ID)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
	// 	return
	// }
	// if comment.ID == 0 || comment.UserID != uint64(userId) {
	// 	ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Comment Not Found"})
	// 	return
	// }

	err = u.svc.DeleteComment(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pkg.SuccessResponse{Message: "Your comment has been successfully deleted"})
}
