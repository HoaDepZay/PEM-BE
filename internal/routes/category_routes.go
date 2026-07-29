package routes

import (
	"visualfinance/internal/controllers"
	"visualfinance/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(router *gin.Engine) {
	categoryController := controllers.NewCategoryController()

	api := router.Group("/api/v1")
	{
		categories := api.Group("/categories")
		categories.Use(middleware.RequireAuth())
		{
			categories.GET("/", categoryController.GetCategories)
			categories.POST("/", categoryController.CreateCategory)
			categories.PUT("/:id", categoryController.UpdateCategory)
			categories.DELETE("/:id", categoryController.DeleteCategory)
		}
	}
}
