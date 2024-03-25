package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MidnightHelix/MyGram/internal/repository"
	"github.com/MidnightHelix/MyGram/pkg"
	"github.com/MidnightHelix/MyGram/pkg/helper"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

type AuthorizationMiddleware struct {
	UserRepository        repository.UserQuery
	PhotoRepository       repository.PhotoQuery
	CommentRepository     repository.CommentQuery
	SocialMediaRepository repository.SocialMediaQuery
}

func NewAuthMiddleware(userRepository repository.UserQuery,
	photoRepository repository.PhotoQuery,
	commentRepository repository.CommentQuery,
	socialMediaRepository repository.SocialMediaQuery) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		UserRepository:        userRepository,
		PhotoRepository:       photoRepository,
		CommentRepository:     commentRepository,
		SocialMediaRepository: socialMediaRepository,
	}
}

func (m *AuthorizationMiddleware) Authentication(ctx *gin.Context) {
	auth := ctx.GetHeader("Authorization")

	authArr := strings.Split(auth, " ")
	if len(authArr) < 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, pkg.ErrorResponse{
			Message: "unauthorized",
			Errors:  []string{"invalid token"},
		})
		return
	}
	if authArr[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, pkg.ErrorResponse{
			Message: "unauthorized",
			Errors:  []string{"invalid authorization method"},
		})
		return
	}

	token := authArr[1]
	claims, err := helper.ValidateToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, pkg.ErrorResponse{
			Message: "unauthorized",
			Errors:  []string{"invalid token", "failed to decode"},
		})
		return
	}

	ctx.Set("claims", claims)
	ctx.Next()
}

func (m *AuthorizationMiddleware) UserAuthorization(ctx *gin.Context) {

	claims, ok := ctx.Get("claims")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, pkg.ErrorResponse{
			Message: "Unauthorized",
			Errors:  []string{"Missing claims in context"},
		})
		return
	}

	userID := claims.(jwt.MapClaims)["user_id"].(float64)
	userId := int(userID)
	if userId == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}

	user, err := m.UserRepository.GetUsersByID(ctx, uint64(id))
	fmt.Println("user ID : ", user.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, pkg.ErrorResponse{
			Message: "Internal Server Error",
			Errors:  []string{err.Error()},
		})
		return
	}
	if user.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, pkg.ErrorResponse{Message: "User Not Found"})
		return
	}

	if user.ID != uint64(userId) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, pkg.ErrorResponse{
			Message: "Forbidden",
			Errors:  []string{"You are not authorized to modify this user"},
		})
	}

	ctx.Set("claims", claims)
	ctx.Next()
}

func (m *AuthorizationMiddleware) PhotoAuthorization(ctx *gin.Context) {

	claims, ok := ctx.Get("claims")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, pkg.ErrorResponse{
			Message: "Unauthorized",
			Errors:  []string{"Missing claims in context"},
		})
		return
	}

	userID := claims.(jwt.MapClaims)["user_id"].(float64)
	userId := int(userID)
	if userId == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}

	photo, err := m.PhotoRepository.GetPhotosByID(ctx, uint64(id))
	fmt.Println("photo ID : ", photo.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, pkg.ErrorResponse{
			Message: "Internal Server Error",
			Errors:  []string{err.Error()},
		})
		return
	}
	if photo.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Photo Not Found"})
		return
	}

	if photo.UserID != uint64(userId) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, pkg.ErrorResponse{
			Message: "Forbidden",
			Errors:  []string{"You are not authorized to modify this photo"},
		})
	}

	ctx.Next()
}

func (m *AuthorizationMiddleware) CommentAuthorization(ctx *gin.Context) {

	claims, ok := ctx.Get("claims")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, pkg.ErrorResponse{
			Message: "Unauthorized",
			Errors:  []string{"Missing claims in context"},
		})
		return
	}

	userID := claims.(jwt.MapClaims)["user_id"].(float64)
	userId := int(userID)
	if userId == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}

	comment, err := m.CommentRepository.GetCommentsByID(ctx, uint64(id))
	fmt.Println("comment ID : ", comment.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, pkg.ErrorResponse{
			Message: "Internal Server Error",
			Errors:  []string{err.Error()},
		})
		return
	}
	if comment.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Comment Not Found"})
		return
	}

	if comment.UserID != uint64(userId) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, pkg.ErrorResponse{
			Message: "Forbidden",
			Errors:  []string{"You are not authorized to modify this comment"},
		})
	}

	ctx.Next()
}

func (m *AuthorizationMiddleware) SocialMediaAuthorization(ctx *gin.Context) {

	claims, ok := ctx.Get("claims")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, pkg.ErrorResponse{
			Message: "Unauthorized",
			Errors:  []string{"Missing claims in context"},
		})
		return
	}

	userID := claims.(jwt.MapClaims)["user_id"].(float64)
	userId := int(userID)
	if userId == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}

	socialMedia, err := m.SocialMediaRepository.GetSocialMediaByID(ctx, uint64(id))
	fmt.Println("socialMedia ID : ", socialMedia.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, pkg.ErrorResponse{
			Message: "Internal Server Error",
			Errors:  []string{err.Error()},
		})
		return
	}
	if socialMedia.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Social Media Not Found"})
		return
	}

	if socialMedia.UserID != uint64(userId) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, pkg.ErrorResponse{
			Message: "Forbidden",
			Errors:  []string{"You are not authorized to modify this Social Media"},
		})
	}

	ctx.Next()
}
