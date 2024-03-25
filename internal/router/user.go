package router

import (
	"github.com/MidnightHelix/MyGram/internal/handler"
	"github.com/MidnightHelix/MyGram/internal/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouter interface {
	Mount()
}

type userRouterImpl struct {
	v              *gin.RouterGroup
	handler        handler.UserHandler
	authMiddleware middleware.AuthorizationMiddleware
}

func NewUserRouter(v *gin.RouterGroup, handler handler.UserHandler, authMiddleware middleware.AuthorizationMiddleware) UserRouter {
	return &userRouterImpl{v: v, handler: handler, authMiddleware: authMiddleware}
}

func (u *userRouterImpl) Mount() {
	// activity
	// /users/register
	u.v.POST("/register", u.handler.UserSignUp)
	// /users/login
	u.v.POST("/login", u.handler.UserLogin)

	// users
	u.v.Use(u.authMiddleware.Authentication)
	// /users
	u.v.GET("", u.handler.GetUsers)
	// /users/:id
	u.v.GET("/:id", u.handler.GetUsersById)
	// PUT /users
	u.v.PUT("/:id", u.authMiddleware.UserAuthorization, u.handler.EditUser)
	// DELETE /users
	u.v.DELETE("", u.handler.DeleteUser)
}
