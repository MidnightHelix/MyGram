package main

import (
	"github.com/MidnightHelix/MyGram/internal/handler"
	"github.com/MidnightHelix/MyGram/internal/infrastructure"
	"github.com/MidnightHelix/MyGram/internal/middleware"
	"github.com/MidnightHelix/MyGram/internal/repository"
	"github.com/MidnightHelix/MyGram/internal/router"
	"github.com/MidnightHelix/MyGram/internal/service"
	"github.com/MidnightHelix/MyGram/pkg/validator"
	"github.com/gin-gonic/gin"

	_ "github.com/MidnightHelix/MyGram/cmd/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			GO DTS MYGRAM DOCUMENTATION
// @version		1.0
// @description	golang kominfo 006 MyGram api documentation
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.url	http://www.swagger.io/support
// @contact.email	support@swagger.io
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:3000
// @BasePath		/api/v1
// @schemes		http
func main() {
	g := gin.Default()
	v1 := g.Group("/api/v1")
	usersGroup := v1.Group("/users")
	photosGroup := v1.Group("/photos")
	commentsGroup := v1.Group("/comments")
	socialMediasGroup := v1.Group("/socialmedias")

	// dependency injection
	// dig by uber
	// wire

	// https://s8sg.medium.com/solid-principle-in-go-e1a624290346
	gorm := infrastructure.NewGormPostgres()
	userRepo := repository.NewUserQuery(gorm)
	photoRepo := repository.NewPhotoQuery(gorm)
	commentRepo := repository.NewCommentQuery(gorm)
	socialMediaRepo := repository.NewSocialMediaQuery(gorm)
	authMiddleware := middleware.NewAuthMiddleware(userRepo, photoRepo, commentRepo, socialMediaRepo)
	customValidator := validator.NewCustomValidator()

	userSvc := service.NewUserService(userRepo)
	userHdl := handler.NewUserHandler(userSvc, customValidator)
	userRouter := router.NewUserRouter(usersGroup, userHdl, *authMiddleware)

	photoSvc := service.NewPhotoService(photoRepo)
	photoHdl := handler.NewPhotoHandler(photoSvc, customValidator)
	photoRouter := router.NewPhotoRouter(photosGroup, photoHdl, *authMiddleware)

	commentSvc := service.NewCommentService(commentRepo)
	commentHdl := handler.NewCommentHandler(commentSvc, customValidator)
	commentRouter := router.NewCommentRouter(commentsGroup, commentHdl, *authMiddleware)

	socialMediaSvc := service.NewSocialMediaService(socialMediaRepo)
	socialMediaHdl := handler.NewSocialMediaHandler(socialMediaSvc, customValidator)
	socialMediaRouter := router.NewSocialMediaRouter(socialMediasGroup, socialMediaHdl, *authMiddleware)

	// mount
	userRouter.Mount()
	photoRouter.Mount()
	commentRouter.Mount()
	socialMediaRouter.Mount()
	// swagger
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	g.Run(":3000")
}
