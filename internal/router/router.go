package router

import (
	"learn-gin/internal/data"
	"learn-gin/internal/handler"
	"learn-gin/internal/middleware"
	repo "learn-gin/internal/repository"
	"learn-gin/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	tm := data.NewTransactionManager(db)

	userRepo := repo.NewUserRepo()
	walletRepo := repo.NewWalletRepo()
	refreshTokenRepo := repo.NewRefreshTokenRepo()

	sessionService := service.NewSessionService(userRepo, refreshTokenRepo, tm)
	userService := service.NewUserService(tm, userRepo, walletRepo)

	sessionHandler := handler.NewSessionHandler(sessionService)
	userHandler := handler.NewUserHandler(userService)
	r := gin.New()

	r.Use(middleware.GinLogger())
	r.Use(gin.Recovery())

	api := r.Group("")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", middleware.RateLimitMiddleware(5, 60*time.Second), sessionHandler.Login)
		api.DELETE("/logout", sessionHandler.Logout)
		api.POST("/refresh", sessionHandler.Refresh)
		api.POST("/send-code", middleware.RateLimitMiddleware(1, 60*time.Second), userHandler.SendVerificationCode)
		api.POST("/verify-code", userHandler.VerifyCode)

		authGroup := api.Group("/user")
		authGroup.Use(middleware.Auth())
		{
			authGroup.GET("/me", userHandler.GetMe)
		}
	}
	return r
}