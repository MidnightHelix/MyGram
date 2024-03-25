package router

import (
	"github.com/MidnightHelix/MyGram/internal/handler"
	"github.com/MidnightHelix/MyGram/internal/middleware"
	"github.com/gin-gonic/gin"
)

type PhotoRouter interface {
	Mount()
}

type photoRouterImpl struct {
	v              *gin.RouterGroup
	handler        handler.PhotoHandler
	authMiddleware middleware.AuthorizationMiddleware
}

// func NewPhotoRouter(v *gin.RouterGroup, handler handler.PhotoHandler) PhotoRouter {
// 	return &photoRouterImpl{v: v, handler: handler}
// }

// func (u *photoRouterImpl) Mount() {

// 	u.v.Use(middleware.Authentication)

// 	u.v.POST("", u.handler.PostPhoto)

// 	u.v.GET("", u.handler.GetPhotos)

// 	u.v.PUT("/:id", u.handler.EditPhoto)

// 	u.v.DELETE("/:id", u.handler.DeletePhoto)
// }

func NewPhotoRouter(v *gin.RouterGroup, handler handler.PhotoHandler, authMiddleware middleware.AuthorizationMiddleware) PhotoRouter {
	return &photoRouterImpl{v: v, handler: handler, authMiddleware: authMiddleware}
}

func (u *photoRouterImpl) Mount() {

	u.v.Use(u.authMiddleware.Authentication)

	u.v.POST("", u.handler.PostPhoto)

	u.v.GET("", u.handler.GetPhotos)

	u.v.PUT("/:id", u.authMiddleware.PhotoAuthorization, u.handler.EditPhoto)

	u.v.DELETE("/:id", u.authMiddleware.PhotoAuthorization, u.handler.DeletePhoto)
}
