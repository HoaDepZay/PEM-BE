package routes

import (
	"visualfinance/internal/controllers"
	"visualfinance/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/register", controllers.Register)
		authGroup.POST("/login", controllers.Login)
		authGroup.POST("/refresh", controllers.Refresh)
		authGroup.POST("/logout", controllers.Logout)
		authGroup.GET("/verify", controllers.VerifyEmail)
		
		authGroup.POST("/forgot-password", controllers.ForgotPassword)
		authGroup.POST("/verify-otp", controllers.VerifyOTP)
		authGroup.POST("/reset-password", controllers.ResetPassword)

		// Protected routes
		protected := authGroup.Group("/")
		protected.Use(middleware.RequireAuth())
		{
			protected.GET("/me", controllers.GetMe)
			protected.POST("/change-password", controllers.ChangePassword)
		}
	}
}
