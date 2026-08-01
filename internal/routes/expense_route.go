package routes

import (
	"visualfinance/internal/controllers"
	"visualfinance/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupExpenseRoutes(router *gin.Engine) {
	expenseGroup := router.Group("/api/expenses")
	expenseGroup.Use(middleware.RequireAuth()) // Bắt buộc đăng nhập
	{
		expenseGroup.POST("/", controllers.CreateExpense)
		expenseGroup.GET("/", controllers.GetExpenses)
		expenseGroup.PUT("/:id", controllers.UpdateExpense)
		expenseGroup.DELETE("/:id", controllers.DeleteExpense)
	}
}
