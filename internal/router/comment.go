package router

import (
	"github.com/MidnightHelix/MyGram/internal/handler"
	"github.com/MidnightHelix/MyGram/internal/middleware"
	"github.com/gin-gonic/gin"
)

type CommentRouter interface {
	Mount()
}

type commentRouterImpl struct {
	v              *gin.RouterGroup
	handler        handler.CommentHandler
	authMiddleware middleware.AuthorizationMiddleware
}

func NewCommentRouter(v *gin.RouterGroup, handler handler.CommentHandler, authMiddleware middleware.AuthorizationMiddleware) CommentRouter {
	return &commentRouterImpl{v: v, handler: handler, authMiddleware: authMiddleware}
}

func (u *commentRouterImpl) Mount() {

	u.v.Use(u.authMiddleware.Authentication)

	u.v.POST("", u.handler.PostComment)

	u.v.GET("", u.handler.GetComments)

	u.v.PUT("/:id", u.authMiddleware.CommentAuthorization, u.handler.EditComment)

	u.v.DELETE("/:id", u.authMiddleware.CommentAuthorization, u.handler.DeleteComment)
}
