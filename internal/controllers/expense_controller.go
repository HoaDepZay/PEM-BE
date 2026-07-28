package controllers

import (
	"net/http"

	"visualfinance/internal/models"
	"visualfinance/internal/pkg/db"
	"visualfinance/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func CreateExpense(c *gin.Context) {
	var expense models.Expense

	if err := c.ShouldBindJSON(&expense); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := db.DB.Create(&expense).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create expense: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Expense created successfully", expense)
}

func GetExpenses(c *gin.Context) {
	var expenses []models.Expense
	// For testing, just fetch top 50, optionally add user filtering later
	if err := db.DB.Limit(50).Order("ExpenseDate DESC").Find(&expenses).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch expenses: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Expenses retrieved successfully", expenses)
}
