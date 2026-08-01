package routes

import (
	"visualfinance/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupMediaRoutes(router *gin.Engine) {
	router.GET("/api/images/:bucket/*filename", controllers.ServeMedia)
}
