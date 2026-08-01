package routes

import (
	"visualfinance/internal/controllers"
	"visualfinance/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupProfileRoutes(r *gin.Engine) {
	profileGroup := r.Group("/api/profile")
	profileGroup.Use(middleware.RequireAuth())
	{
		profileGroup.GET("/me", controllers.GetProfile)
		profileGroup.POST("/avatar", controllers.UploadAvatar)
		profileGroup.PUT("/info", controllers.UpdateProfile)
	}
}
