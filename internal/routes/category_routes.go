package routes

import (
	"visualfinance/internal/controllers"
	"visualfinance/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(router *gin.Engine) {
	categoryController := controllers.NewCategoryController()

	categories := router.Group("/api/categories")
	{
		categories.Use(middleware.RequireAuth())
		{
			categories.GET("/", categoryController.GetCategories)
			
			// Groups
			categories.POST("/groups", categoryController.CreateGroup)
			categories.PUT("/groups/:id", categoryController.UpdateGroup)
			categories.DELETE("/groups/:id", categoryController.DeleteGroup)

			// Sub-Categories
			categories.POST("/", categoryController.CreateCategory)
			categories.PUT("/:id", categoryController.UpdateCategory)
			categories.DELETE("/:id", categoryController.DeleteCategory)
		}
	}
}
