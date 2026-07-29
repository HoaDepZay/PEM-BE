package routes

import (
	"visualfinance/internal/controllers"
	"visualfinance/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/register", controllers.Register)
		authGroup.POST("/login", controllers.Login)
		authGroup.GET("/verify", controllers.VerifyEmail)
		
		// Protected route example
		authGroup.GET("/me", middleware.RequireAuth(), controllers.GetMe)
	}
}
