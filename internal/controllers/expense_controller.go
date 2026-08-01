package controllers

import (
	"fmt"
	"net/http"

	"visualfinance/internal/pkg/response"
	"visualfinance/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var expenseService = services.NewExpenseService()

// CreateExpense godoc
// @Summary      Create a new expense with Image
// @Description  Uploads an image, parses the note for amount, and saves to DB.
// @Tags         expenses
// @Accept       multipart/form-data
// @Produce      json
// @Param        image formData file true "Image file"
// @Param        note formData string true "Note containing amount"
// @Success      201      {object}  response.Response
// @Failure      400      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       / [post]
func CreateExpense(c *gin.Context) {
	note := c.PostForm("note")
	categoryID := c.PostForm("category_id")
	amountStr := c.PostForm("amount")
	
	var amount int64
	if amountStr != "" {
		fmt.Sscanf(amountStr, "%d", &amount)
	}
	
	userIDStrObj, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userIDStr := userIDStrObj.(uuid.UUID).String()

	// 1. Get the file from form
	fileHeader, err := c.FormFile("image")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Image file is required")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to read image file")
		return
	}
	defer file.Close()

	// 2. Call service layer
	expense, err := expenseService.CreateExpense(
		c.Request.Context(),
		userIDStr,
		categoryID,
		note,
		fileHeader.Filename,
		fileHeader.Header.Get("Content-Type"),
		amount,
		fileHeader.Size,
		file,
	)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Expense created successfully", expense)
}

// GetExpenses godoc
// @Summary      Get all expenses
// @Description  Get a list of expenses
// @Tags         expenses
// @Produce      json
// @Success      200  {array}   models.Expense
// @Failure      500  {object}  response.Response
// @Router       / [get]
func GetExpenses(c *gin.Context) {
	userIDStrObj, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userIDStr := userIDStrObj.(uuid.UUID).String()

	expenses, err := expenseService.GetExpenses(userIDStr)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch expenses: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Expenses retrieved successfully", expenses)
}

// UpdateExpense godoc
// @Summary      Update an expense
// @Tags         expenses
// @Accept       multipart/form-data
// @Produce      json
// @Param        id path string true "Expense ID"
// @Param        note formData string false "Note"
// @Param        category_id formData string false "Category ID"
// @Param        amount formData string false "Amount"
// @Success      200      {object}  response.Response
// @Failure      400      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /{id} [put]
func UpdateExpense(c *gin.Context) {
	userIDStrObj, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userIDStr := userIDStrObj.(uuid.UUID).String()
	expenseID := c.Param("id")

	note := c.PostForm("note")
	categoryID := c.PostForm("category_id")
	amountStr := c.PostForm("amount")

	expense, err := expenseService.UpdateExpense(userIDStr, expenseID, categoryID, note, amountStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Expense updated successfully", expense)
}

// DeleteExpense godoc
// @Summary      Delete an expense
// @Tags         expenses
// @Produce      json
// @Param        id path string true "Expense ID"
// @Success      200      {object}  response.Response
// @Failure      400      {object}  response.Response
// @Router       /{id} [delete]
func DeleteExpense(c *gin.Context) {
	userIDStrObj, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userIDStr := userIDStrObj.(uuid.UUID).String()
	expenseID := c.Param("id")

	err := expenseService.DeleteExpense(userIDStr, expenseID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Expense deleted successfully", nil)
}

