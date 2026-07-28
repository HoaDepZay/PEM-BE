package routes

import (
	"visualfinance/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupExpenseRoutes(router *gin.Engine) {
	expenseGroup := router.Group("/api/v1/expenses")
	{
		expenseGroup.POST("/", controllers.CreateExpense)
		expenseGroup.GET("/", controllers.GetExpenses)
	}
}
