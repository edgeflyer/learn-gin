package router

import (
	"database/sql"
	"learn-gin/internal/data"
	"learn-gin/internal/handler"
	"learn-gin/internal/middleware"
	repo "learn-gin/internal/repository"
	"learn-gin/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB) *gin.Engine {
	tm := data.NewTransactionManager(db)

	userRepo := repo.NewUserRepo()
	walletRepo := repo.NewWalletRepo()

	userService := service.NewUserService(tm, userRepo, walletRepo)

	userHandler := handler.NewUserHandler(userService)
	r := gin.New()

	r.Use(middleware.GinLogger())
	r.Use(gin.Recovery())

	api := r.Group("")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		authGroup := api.Group("/user")
		authGroup.Use(middleware.Auth())
		{
			authGroup.GET("/me", userHandler.GetMe)
		}
	}
	return r
}