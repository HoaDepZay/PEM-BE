package controllers

import (
	"net/http"

	"time"

	"visualfinance/internal/models"
	"visualfinance/internal/pkg/db"
	"visualfinance/internal/pkg/minio"
	"visualfinance/internal/pkg/response"
	"visualfinance/internal/pkg/textparser"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateExpense godoc
// @Summary      Create a new expense with Image
// @Description  Uploads an image, parses the note for amount, and saves to DB.
// @Tags         expenses
// @Accept       multipart/form-data
// @Produce      json
// @Param        image formData file true "Image file"
// @Param        note formData string true "Note containing amount"
// @Param        user_id formData string true "User ID"
// @Success      201      {object}  response.Response
// @Failure      400      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       / [post]
func CreateExpense(c *gin.Context) {
	note := c.PostForm("note")
	userIDStr := c.PostForm("user_id")

	// Validate user ID
	if userIDStr == "" {
		// Temporary fallback for testing if no userID is provided
		userIDStr = "00000000-0000-0000-0000-000000000000"
	}

	// 1. Get the file from form
	fileHeader, err := c.FormFile("image")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Image file is required")
		return
	}

	// 2. Upload to MinIO
	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to read image file")
		return
	}
	defer file.Close()

	objectName := uuid.New().String() + "-" + fileHeader.Filename
	imageURL, err := minio.UploadImage(c.Request.Context(), objectName, file, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to upload image to MinIO: "+err.Error())
		return
	}

	// 3. Parse Amount using Regex
	amount, err := textparser.ParseAmount(note)
	if err != nil || amount <= 0 {
		// Nếu không tìm thấy số tiền, để tạm 0 hoặc báo lỗi
		// Ta ưu tiên cho phép lưu và sửa sau (theo flow Snap & Sync)
	}

	// 4. Save to Database
	expense := models.Expense{
		UserID:      userIDStr,
		Amount:      int64(amount),
		Note:        note,
		ImageURL:    imageURL,
		ExpenseDate: time.Now(),
	}

	if err := db.DB.Create(&expense).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create expense in DB: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Expense created successfully", expense)
}

// GetExpenses godoc
// @Summary      Get all expenses
// @Description  Get a list of expenses (limited to 50 for now)
// @Tags         expenses
// @Produce      json
// @Success      200  {array}   models.Expense
// @Failure      500  {object}  response.Response
// @Router       / [get]
func GetExpenses(c *gin.Context) {
	var expenses []models.Expense
	// For testing, just fetch top 50, optionally add user filtering later
	if err := db.DB.Limit(50).Order("ExpenseDate DESC").Find(&expenses).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch expenses: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Expenses retrieved successfully", expenses)
}
