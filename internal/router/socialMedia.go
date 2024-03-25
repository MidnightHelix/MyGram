package router

import (
	"github.com/MidnightHelix/MyGram/internal/handler"
	"github.com/MidnightHelix/MyGram/internal/middleware"
	"github.com/gin-gonic/gin"
)

type SocialMediaRouter interface {
	Mount()
}

type socialMediaRouterImpl struct {
	v              *gin.RouterGroup
	handler        handler.SocialMediaHandler
	authMiddleware middleware.AuthorizationMiddleware
}

func NewSocialMediaRouter(v *gin.RouterGroup, handler handler.SocialMediaHandler, authMiddleware middleware.AuthorizationMiddleware) SocialMediaRouter {
	return &socialMediaRouterImpl{v: v, handler: handler, authMiddleware: authMiddleware}
}

func (u *socialMediaRouterImpl) Mount() {

	u.v.Use(u.authMiddleware.Authentication)

	u.v.POST("", u.handler.CreateSocialMedia)

	u.v.GET("", u.handler.GetSocialMedias)

	u.v.PUT("/:id", u.authMiddleware.SocialMediaAuthorization, u.handler.EditSocialMedia)

	u.v.DELETE("/:id", u.authMiddleware.SocialMediaAuthorization, u.handler.DeleteSocialMedia)
}
